package main

import (
	"bytes"
	"encoding/json"
	"errors"
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

type FileResponse struct {
	Size     int    `json:"size"`
	Type     string `json:"type"`
	Language string `json:"language"`
}

type GistResponse struct {
	Description string                   `json:"description"`
	Files       map[string]*FileResponse `json:"files"`
	HtmlUrl     string                   `json:"html_url"`
	Id          string                   `json:"id"`
	Public      bool                     `json:"public"`
	User        string                   `json:"user"`
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

	if config.APIKey != "" && !anonymous {
		req.Header.Add("Authorization", "token "+config.APIKey)
	}

	body, err := doRequest(req)
	if err != nil {
		return "", err
	}

	var gistResp GistResponse
	err = json.Unmarshal(body, &gistResp)
	if err != nil {
		return "", err
	}

	return gistResp.HtmlUrl, nil
}

func (gist *Gist) Update(uid string) (string, error) {
	buf := bytes.NewBuffer(nil)
	e := json.NewEncoder(buf)
	if err := e.Encode(gist); err != nil {
		return "", err
	}

	req, err := http.NewRequest("PATCH", GitHubAPIURL+"/gists/"+uid, buf)
	if err != nil {
		return "", err
	}

	if config.APIKey != "" {
		req.Header.Add("Authorization", "token "+config.APIKey)
	}

	body, err := doRequest(req)
	if err != nil {
		return "", err
	}

	var gistResp GistResponse
	err = json.Unmarshal(body, &gistResp)
	if err != nil {
		return "", err
	}

	return gistResp.HtmlUrl, nil
}

func GistList(user string) ([]*GistResponse, error) {
	var uri string
	if config.APIKey != "" && user != "" {
		uri = "/users/" + user + "/gists"
	} else {
		uri = "/gists"
	}

	req, err := http.NewRequest("GET", GitHubAPIURL+uri, nil)
	if err != nil {
		return nil, err
	}

	if config.APIKey != "" {
		req.Header.Add("Authorization", "token "+config.APIKey)
	}

	body, err := doRequest(req)
	if err != nil {
		return nil, err
	}

	var gistResp []*GistResponse
	err = json.Unmarshal(body, &gistResp)
	if err != nil {
		return nil, err
	}

	return gistResp, nil
}

func doRequest(req *http.Request) ([]byte, error) {
	req.Header.Add("Content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		var f interface{}
		if err := json.Unmarshal(body, &f); err != nil {
			return nil, err
		}
		return nil, errors.New(fmt.Sprintf("%v", f.(map[string]interface{})["message"]))
	}

	return body, nil
}
