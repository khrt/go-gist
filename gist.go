package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const GitHubAPIURL = "https://api.github.com"

type File struct {
	Name    string `json:"-"`
	Content string `json:"content"`
}

type Gist struct {
	Files       map[string]*File `json:"files"`
	Description string           `json:"description"`
	Public      bool             `json:"public"`
}

func (gist *Gist) Create(anonymous bool) (string, error) {
	buf := bytes.NewBuffer(nil)
	e := json.NewEncoder(buf)
	if err := e.Encode(gist); err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", GitHubAPIURL+"/gists", buf)
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-type", "application/json")
	if !anonymous {
		req.Header.Add("Authorization", "token "+config.APIKey)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	fmt.Println(string(body))

	return "", nil
}

func (gist *Gist) Update(uid string) (string, error) {
	return "", nil
}

func GistList(user string) ([]*Gist, error) {
	return nil, nil
}
