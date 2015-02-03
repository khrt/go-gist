package main

import (
	"os"
	"testing"
)

func TestConfigNew(t *testing.T) {
	c := ConfigNew()
	t.Log("Config:", c)
}

func TestConfigLoad(t *testing.T) {
	c := ConfigNew()

	if err := c.Load(); err != nil {
		t.Fatal(err)
	}

	t.Log("Config:", c)
}

func TestConfigUpdate(t *testing.T) {
	os.Setenv(ConfigFileEnv, "./env-gogist.conf")
	c := ConfigNew()

	if err := c.Load(); err != nil {
		t.Log("error", err)
	}

	err := c.Update("newkey")
	if err != nil {
		t.Fatal(err)
	}

	t.Log("Config:", c)
}
