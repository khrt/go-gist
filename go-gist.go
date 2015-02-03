package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strings"
)

var config = ConfigNew()

func main() {
	err := config.Load()
	if err != nil {
		panic(err)
	}

	anonymousFlag := flag.Bool("a", false, "Create an anonymous gist.")
	description := flag.String("d", "", "A description of the gist.")
	gistType := flag.String("t", "", "Sets the file extension and syntax type.")
	loginFlag := flag.Bool("login", false, "Authenticate gist on this computer.")
	privateFlag := flag.Bool("p", false, "Indicates whether the gist is private.")
	update := flag.String("u", "", "Update an existing gist. Takes ID as an argument.")
	user := flag.String("l", "", "List gists for user.")

	flag.Parse()

	if *loginFlag {
		login()
	} else if *user != "" {
		list(*user)
	} else if flag.NArg() > 0 {
		createOrUpdate(*update, *anonymousFlag, !*privateFlag,
			*description, *gistType, flag.Args())
	}
}

func list(user string) {
	resp, err := GistList(user)
	if err != nil {
		fmt.Println(err)
	}

	for _, r := range resp {
		var files []string
		for f := range r.Files {
			files = append(files, f)
		}

		var secret string
		if r.Public == false {
			secret = "(secret)"
		}

		fmt.Printf("%s %s %s\n", r.HtmlUrl, strings.Join(files, " "), secret)
	}
}

func createOrUpdate(uid string, anonymous bool, public bool, desc string, gistType string, args []string) {
	gist := &Gist{
		make(map[string]*File),
		desc,
		public,
	}

	for _, name := range flag.Args() {
		content, err := ioutil.ReadFile(name)
		if err != nil {
			panic(err)
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
		panic(err)
	}

	fmt.Println(url)
}

func login() {

}
