package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	ConfigFileName = ".gogist"
	ConfigFileEnv  = "GOGIST_CONFIG_FILE"
)

type Config struct {
	APIKey string `json:"apikey"`
	file   string
}

func ConfigNew() *Config {
	c := &Config{}
	return c
}

func (c *Config) Update(APIKey string) error {
	c.APIKey = strings.TrimSpace(APIKey)

	f, err := os.Create(c.file)
	if err != nil {
		return err
	}
	defer f.Close()

	e := json.NewEncoder(f)
	return e.Encode(c)
}

func (c *Config) Load() error {
	fileName, err := c.resolvePath(ConfigFileName, os.Getenv(ConfigFileEnv))
	if err != nil {
		return err
	}

	c.file = fileName

	if _, err := os.Stat(c.file); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("File does not exist")
		}
		return err
	}

	f, err := os.Open(c.file)
	if err != nil {
		return err
	}
	defer f.Close()

	d := json.NewDecoder(f)
	return d.Decode(c)
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

func (c *Config) resolvePath(argPath, envPath string) (string, error) {
	path := envPath

	if path == "" {
		homeDir, err := c.homeDir()
		if err != nil {
			return "", err
		}

		path = filepath.Join(homeDir, argPath)
	}

	return path, nil
}
