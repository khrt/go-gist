package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var config = ConfigNew()

func main() {
	config.Load()

	anonymousFlag := flag.Bool("a", false, "Create an anonymous gist.")
	description := flag.String("d", "", "A description of the gist.")
	gistType := flag.String("t", "", "Sets the file extension and syntax type.")
	loginFlag := flag.Bool("login", false, "Authenticate gist on this computer.")
	privateFlag := flag.Bool("p", false, "Indicates whether the gist is private.")
	update := flag.String("u", "", "Update an existing gist. Takes ID as an argument.")
	listFlag := flag.Bool("l", false, "List gists.")

	flag.Parse()

	var err error

	if *loginFlag {
		err = login()
	} else if *listFlag {
		err = list()
	} else if flag.NArg() > 0 {
		err = createOrUpdate(*update, *anonymousFlag, !*privateFlag, *description, *gistType, flag.Args())
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
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

func createOrUpdate(uid string, anonymous bool, public bool, desc string, gistType string, args []string) error {
	gist := &Gist{make(map[string]*File), desc, public}

	for _, name := range flag.Args() {
		content, err := ioutil.ReadFile(name)
		if err != nil {
			return err
		}

		if gistType != "" {
			name += "." + gistType
		}

		gist.Files[name] = &File{name, string(content)}
	}

	var url string
	var err error

	if uid != "" {
		url, err = gist.Update(uid)
	} else {
		url, err = gist.Create(anonymous)
	}

	if err != nil {
		return err
	}

	fmt.Println(url)
	return nil
}

func login() error {
	fmt.Println("Obtaining OAuth2 access_token from GitHub.")

	var username, password string

	fmt.Print("GitHub usename: ")
	fmt.Scanf("%s", &username)

	fmt.Print("GitHub password: ")
	//if err := exec.Command("/bin/stty", "-echo").Run(); err != nil {
	//	return err
	//}
	fmt.Scanf("%s", &password)
	//if err := exec.Command("/bin/stty", "echo").Run(); err != nil {
	//	return err
	//}

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
