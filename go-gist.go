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
	VERSION         = "0.1"
	DefaultGistName = "a"
)

var config = ConfigNew()

func main() {
	config.Load()

	anonymousFlag := flag.Bool("a", false, "Create an anonymous gist.")
	copyFlag := flag.Bool("c", false, "Copy the resulting URL to the clipboard.")
	listFlag := flag.Bool("l", false, "List gists.")
	loginFlag := flag.Bool("login", false, "Authenticate gist on this computer.")
	pasteFlag := flag.Bool("P", false, "Paste from the clipboard to gist.")
	privateFlag := flag.Bool("p", false, "Indicates whether the gist is private.")

	desc := flag.String("d", "", "A description of the gist.")
	filename := flag.String("f", "", "Sets the filename and syntax type.")
	filetype := flag.String("t", "", "Sets the file extension and syntax type.")
	uid := flag.String("u", "", "Update an existing gist. Takes ID as an argument.")

	flag.Parse()

	var err error

	if *loginFlag {
		err = login()
	} else if *listFlag {
		err = list()
	} else {
		err = gist(*uid, *desc, *filetype, *filename, !*privateFlag, *anonymousFlag, *copyFlag, *pasteFlag)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
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

func gist(uid, desc, filetype, filename string, public, anonymous, copyFlag, pasteFlag bool) error {
	var err error
	var clipboard *Clipboard

	if copyFlag || pasteFlag {
		if clipboard, err = NewClipboard(); err != nil {
			return err
		}
	}

	gist := &Gist{make(map[string]*File), desc, public}

	if flag.NArg() > 0 {
		for _, name := range flag.Args() {
			content, err := ioutil.ReadFile(name)
			if err != nil {
				return err
			}

			if filetype != "" {
				name += "." + filetype
			}

			gist.Files[name] = &File{name, string(content)}
		}
	} else {
		var name string
		var content []byte

		if filename != "" {
			name = filename
		} else if filename == "" && filetype != "" {
			name = DefaultGistName + "." + filetype
		} else {
			name = DefaultGistName
		}

		if pasteFlag {
			if content, err = clipboard.Paste(); err != nil {
				return err
			}
		} else {
			if content, err = ioutil.ReadAll(os.Stdin); err != nil {
				return err
			}
		}

		gist.Files[name] = &File{name, string(content)}
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

	if copyFlag {
		if err := clipboard.Copy(url); err != nil {
			return err
		}
	}

	fmt.Println(url)
	return nil
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
