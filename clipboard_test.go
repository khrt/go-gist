package main

import (
	"testing"
	"time"
)

func TestClipboard(t *testing.T) {
	c, err := NewClipboard()
	if err != nil {
		t.Fatal(err)
	}

	var message string

	message = "GO COPY TEST " + time.Now().String()

	if err := c.Copy(message); err != nil {
		t.Fatal(err)
	}

	message = "GO PASTE TEST " + time.Now().String()

	if err := c.Copy(message); err != nil {
		t.Fatal(err)
	}

	if p, err := c.Paste(); string(p) != message || err != nil {
		t.Fatal(err)
	}
}
