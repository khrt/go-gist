package main

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const ConfigFileName = ".gist"

type Config struct {
	Token string
	file  string
}

func ConfigNew() *Config {
	c := &Config{}
	c.Load()
	return c
}

func (c *Config) Update(Token string) error {
	c.Token = strings.TrimSpace(Token)
	return ioutil.WriteFile(c.file, []byte(c.Token), 0644)
}

func (c *Config) Load() {
	c.file, _ = c.resolvePath(ConfigFileName)
	apikey, _ := ioutil.ReadFile(c.file)
	if len(apikey) > 0 {
		c.Token = strings.TrimSpace(string(apikey))
	}
}

func (c *Config) homeDir() (string, error) {
	var dir string

	if runtime.GOOS == "windows" {
		dir = os.Getenv("USERPROFILE")
		if dir == "" {
			dir = os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		}
	} else {
		dir = os.Getenv("HOME")
	}

	if dir == "" {
		return dir, errors.New("Home not found")
	}

	return dir, nil
}

func (c *Config) resolvePath(argPath string) (string, error) {
	homeDir, err := c.homeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, argPath), nil
}
