package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
)

const (
	DefaultGistName = "a"
	GistURI         = "gist.github.com"
	GitHubAPIURL    = "https://api.github.com"
)

type File struct {
	Name    string `json:"-"`
	Content string `json:"content"`
	// Response:
	Size     int    `json:"size"`
	Type     string `json:"type"`
	Language string `json:"language"`
}

type Gist struct {
	Files       map[string]*File `json:"files"`
	Description string           `json:"description"`
	Public      bool             `json:"public"`
	// Response:
	HtmlUrl string `json:"html_url"`
	Id      string `json:"id"`
	User    string `json:"user"`
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

	if config.Token != "" && !anonymous {
		req.Header.Add("Authorization", "token "+config.Token)
	}

	body, err := doRequest(req)
	if err != nil {
		return "", err
	}

	if err := json.Unmarshal(body, &gist); err != nil {
		return "", err
	}

	return gist.HtmlUrl, nil
}

func (gist *Gist) Update(uid string) (string, error) {
	re := regexp.MustCompile("http(?:s?)://" + GistURI + "/(.+)")
	if matched := re.MatchString(uid); matched {
		uid = re.ReplaceAllString(uid, "$1")
	}

	buf := bytes.NewBuffer(nil)
	e := json.NewEncoder(buf)
	if err := e.Encode(gist); err != nil {
		return "", err
	}

	req, err := http.NewRequest("PATCH", GitHubAPIURL+"/gists/"+uid, buf)
	if err != nil {
		return "", err
	}

	if config.Token != "" {
		req.Header.Add("Authorization", "token "+config.Token)
	}

	body, err := doRequest(req)
	if err != nil {
		return "", err
	}

	var resp Gist
	if err := json.Unmarshal(body, &resp); err != nil {
		return "", err
	}

	return resp.HtmlUrl, nil
}

func GistList() ([]*Gist, error) {
	req, err := http.NewRequest("GET", GitHubAPIURL+"/gists", nil)
	if err != nil {
		return nil, err
	}

	if config.Token != "" {
		req.Header.Add("Authorization", "token "+config.Token)
	}

	body, err := doRequest(req)
	if err != nil {
		return nil, err
	}

	var resp []*Gist
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

func (gist *Gist) Rawify() (string, error) {
	if gist.HtmlUrl == "" {
		return "", errors.New("HTML URL is empty.")
	}

	req, _ := http.NewRequest("GET", gist.HtmlUrl, nil)
	client := &http.Transport{}
	resp, err := client.RoundTrip(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var url string

	if resp.StatusCode == 200 {
		url = gist.HtmlUrl + "/raw"
	} else if resp.StatusCode == 302 {
		gist.HtmlUrl = resp.Header.Get("location")
		if url, err = gist.Rawify(); err != nil {
			return "", err
		}
	}

	return url, nil
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

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return nil, GistParseError(body)
	}

	return body, nil
}
