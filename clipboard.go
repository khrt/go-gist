package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Clipboard struct {
	copyCmd, pasteCmd string
}

var ClipboardCommands = []Clipboard{
	Clipboard{"xclip", "xclip -o"},
	Clipboard{"xsel -i", "xsel -o"},
	Clipboard{"pbcopy", "pbpaste"},
	Clipboard{"putclip", "getclip"},
}

func NewClipboard() (*Clipboard, error) {

	for _, c := range ClipboardCommands {
		fmt.Println(c)
	}

	//if runtime.GOOS == "windows" {
	//} else {
	//}

	return &ClipboardCommands[2], nil
}

func (c *Clipboard) Copy(content string) error {
	cmd := exec.Command(c.copyCmd)
	cmd.Stdin = strings.NewReader(content)
	if err := cmd.Run(); err != nil {
		return err
	}

	if p, err := c.Paste(); string(p) != content || err != nil {
		message := "Copying to clipboard failed."

		if os.Getenv("TMUX") != "" && c.copyCmd == "pbcopy" {
			message = message + "\nIf you're running tmux on a mac, try http://robots.thoughtbot.com/post/19398560514/how-to-copy-and-paste-with-tmux-on-mac-os-x"
		}

		return errors.New(message)
	}

	return nil
}

func (c *Clipboard) Paste() ([]byte, error) {
	var stdout bytes.Buffer

	cmd := exec.Command(c.pasteCmd)
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	return stdout.Bytes(), nil
}
