package main

import (
	"flag"
	"fmt"
	"io/ioutil"
)

var config = ConfigNew()

func main() {
	err := config.Load()
	if err != nil {
		panic(err)
	}

	login := flag.Bool("login", false, "Authenticate gist on this computer.")
	private := flag.Bool("p", false, "Indicates whether the gist is private.")
	description := flag.String("d", "", "A description of the gist.")
	update := flag.String("u", "", "Update an existing gist.")
	list := flag.String("l", "", "List gists for user.")
	anonymous := flag.Bool("a", false, "Create an anonymous gist.")
	gistType := flag.String("t", "", "Sets the file extension and syntax type.")
	flag.Parse()

	gist := &Gist{
		make(map[string]*File),
		*description,
		!*private,
	}

	for _, name := range flag.Args() {
		content, err := ioutil.ReadFile(name)
		if err != nil {
			panic(err)
		}
		gist.Files[name] = &File{name, string(content)}
	}

	var gistUrl string

	if *update != "" {
		gistUrl, err = gist.Update(*update)
	} else {
		gistUrl, err = gist.Create(*anonymous)
	}

	if err != nil {
		panic(err)
	}

	fmt.Println(gistUrl)
}
