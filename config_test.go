package main

import "testing"

func TestConfigUpdate(t *testing.T) {
	c := ConfigNew()

	origApiKey := c.APIKey // 08ec43864e8131ab4f5778041eb663c67119b786
	newApiKey := "newkey"

	if err := c.Update(newApiKey); err != nil {
		t.Fatal(err)
	}

	c.Load()

	if c.APIKey != newApiKey {
		t.Fail()
	}

	if err := c.Update(origApiKey); err != nil {
		t.Fatal(err)
	}
}
