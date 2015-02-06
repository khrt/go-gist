package main

import "fmt"

type Clipboard struct {
	copyCmd, pasteCmd string
}

var ClipboardCommands = []Clipboard{
	Clipboard{"xclip", "xclip -o"},
	Clipboard{"xsel -i", "xsel -o"},
	Clipboard{"pbcopy", "pbpaste"},
	Clipboard{"putclip", "getclip"},
}

func NewClipboard() *Clipboard {

	for _, c := range ClipboardCommands {
		fmt.Println(c)
	}

	//if runtime.GOOS == "windows" {
	//} else {
	//}

	return nil
}

func (c *Clipboard) Copy(content string) error {
	return nil
}

func (c *Clipboard) Paste() ([]byte, error) {
	return []byte(""), nil
}
