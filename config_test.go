package main

import "testing"

func TestConfigUpdate(t *testing.T) {
	c := NewConfig()

	origApiKey := c.Token
	newApiKey := "newkey"

	if err := c.Update(newApiKey); err != nil {
		t.Fatal(err)
	}

	c.Load()

	if c.Token != newApiKey {
		t.Fail()
	}

	if err := c.Update(origApiKey); err != nil {
		t.Fatal(err)
	}
}
