package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/howeyc/gopass"
)

const (
	VERSION = "0.1"
)

var config = NewConfig()

func main() {
	config.Load()

	//	usage := `go-gist
	//
	//go-gist lets you upload to https://gist.github.com/. Go clone of the official Gist client.
	//
	//Usage:
	//	go-gist (-o|-c|-e) [-p] [-s] [-R] [-d DESC] [-a] [-u URL] [-P] (-f NAME | -t EXT) FILE...
	//	go-gist --login
	//
	//Options:
	//	--login						Authenticate gist on this computer.
	//	-f NAME --filename=NAME		Sets the filename and syntax type.
	//	-t EXT --type=EXT			Sets the file extension and syntax type.
	//	-p --private				Indicates whether the gist is private.
	//	-d DESC --description=DESC  Adds a description to your gist.
	//	-s --shorten                Shorten the gist URL using git.io.
	//	-u URL --update=URL			Update an existing gist.
	//	-a --anonymous				Create an anonymous gist.
	//	-c --copy					Copy the resulting URL to the clipboard.
	//	-e --embed					Copy the embed code for the gist to the clipboard.
	//	-o --open					Open the resulting URL in a browser.
	//	-P --paste					Paste from the clipboard to gist.
	//	-R --raw					Display raw URL of the new gist.
	//	-l USER --list=USER		    Lists all gists for a user.
	//	-h --help					Show this help message and exit.
	//	--version					Show version and exit.`
	//
	//	arguments, err := docopt.Parse(usage, nil, true, "go-gist "+VERSION, false)
	//	fmt.Println(arguments)

	loginFlag := flag.Bool("login", false, "Authenticate gist on this computer.")
	filename := flag.String("f", "", "Sets the filename and syntax type.")
	filetype := flag.String("t", "", "Sets the file extension and syntax type.")
	privateFlag := flag.Bool("p", false, "Indicates whether the gist is private.")
	desc := flag.String("d", "", "A description of the gist.")
	shortenFlag := flag.Bool("s", false, "Shorten the gist URL using git.io.")
	uid := flag.String("u", "", "Update an existing gist. Takes ID as an argument.")
	anonymousFlag := flag.Bool("a", false, "Create an anonymous gist.")
	copyFlag := flag.Bool("c", false, "Copy the resulting URL to the clipboard.")
	openFlag := flag.Bool("o", false, "Open the resulting URL in a browser.")
	pasteFlag := flag.Bool("P", false, "Paste from the clipboard to gist.")
	rawFlag := flag.Bool("R", false, "Display a raw URL of the new gist.")
	listFlag := flag.Bool("l", false, "List gists.")

	versionFlag := flag.Bool("version", false, "Show version and exit.")

	flag.Parse()

	var err error

	if *loginFlag {
		err = login()
	} else if *listFlag {
		err = list()
	} else if *versionFlag {
		fmt.Println("go-gist " + VERSION)
	} else {
		err = gist(*uid, *desc, *filetype, *filename, !*privateFlag, *shortenFlag, *anonymousFlag, *copyFlag, *openFlag, *pasteFlag, *rawFlag)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func login() error {
	fmt.Println("Obtaining OAuth2 access_token from GitHub.")

	var username, password string

	fmt.Print("GitHub username: ")
	fmt.Scanf("%s", &username)

	fmt.Print("GitHub password: ")
	password = string(gopass.GetPasswd())

	fmt.Println()

	token, err := Authorize(username, password)
	if err != nil {
		return err
	}

	if err := config.Update(token); err != nil {
		return err
	}

	fmt.Println("OK")

	return nil
}

func list() error {
	resp, err := GistList()
	if err != nil {
		return err
	}

	for _, r := range resp {
		description := r.Description
		if description == "" {
			var files []string
			for f := range r.Files {
				files = append(files, f)
			}
			description = strings.Join(files, " ")
		}

		var secret string
		if r.Public == false {
			secret = "(secret)"
		}

		fmt.Printf("%s %s %s\n", r.HtmlUrl, description, secret)
	}

	return nil
}

func gist(uid, desc, filetype, filename string, public, shortenFlag, anonymous, copyFlag, openFlag, pasteFlag, rawFlag bool) error {
	var err error
	var clipboard *Clipboard

	if copyFlag || pasteFlag {
		if clipboard, err = NewClipboard(); err != nil {
			return err
		}
	}

	gist := &Gist{Files: make(map[string]*File), Description: desc, Public: public}

	if flag.NArg() > 0 {
		for _, name := range flag.Args() {
			content, err := ioutil.ReadFile(name)
			if err != nil {
				return err
			}

			if filename != "" {
				name = filename
			} else if filetype != "" {
				name += "." + filetype
			}

			gist.Files[name] = &File{Name: name, Content: string(content)}
		}
	} else {
		var name string
		var content []byte

		if filename != "" {
			name = filename
		} else if filetype != "" {
			name = DefaultGistName + "." + filetype
		} else {
			name = DefaultGistName
		}

		if pasteFlag {
			if content, err = clipboard.Paste(); err != nil {
				return err
			}
		} else {
			fmt.Println("(type a gist. <ctrl-c> to cancel, <ctrl-d> when done)")
			if content, err = ioutil.ReadAll(os.Stdin); err != nil {
				return err
			}
		}

		gist.Files[name] = &File{Name: name, Content: string(content)}
	}

	var url string

	if uid != "" {
		url, err = gist.Update(uid)
	} else {
		url, err = gist.Create(anonymous)
	}

	if err != nil {
		return err
	}

	if rawFlag {
		if url, err = gist.Rawify(); err != nil {
			return err
		}
	}

	if shortenFlag {
		if url, err = Shorten(url); err != nil {
			return err
		}
	}

	if copyFlag {
		if err := clipboard.Copy(url); err != nil {
			return err
		}
	}

	if openFlag {
		OpenBrowser(url)
	}

	fmt.Println(url)
	return nil
}
