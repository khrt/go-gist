package main

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"runtime"
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
	var clipboard Clipboard

	for _, c := range ClipboardCommands {
		if runtime.GOOS == "windows" {
			return nil, errors.New("Clipboard doesn't work on Windows.")
		} else {
			cmd := exec.Command("which", strings.Split(c.copyCmd, " ")[0])
			if err := cmd.Run(); err == nil {
				clipboard = c
			}
		}
	}

	return &clipboard, nil
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
