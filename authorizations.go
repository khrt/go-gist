package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type AuthorizationRequest struct {
	Scopes  []string `json:"scopes"`
	Note    string   `json:"note"`
	NoteUrl string   `json:"note_url"`
}

type AuthorizationResponse struct {
	Id    int
	URL   string
	Token string
}

func Authorize(username, password string) (string, error) {
	payload := AuthorizationRequest{
		[]string{"gist"},
		fmt.Sprintf("go-gist (%d)", time.Now().Unix()),
		"https://github.com/khrt/go-gist",
	}

	buf := bytes.NewBuffer(nil)
	e := json.NewEncoder(buf)
	if err := e.Encode(payload); err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", GitHubAPIURL+"/authorizations", buf)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-type", "application/json")
	req.SetBasicAuth(username, password)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode == 401 && resp.Header.Get("X-GitHub-OTP") != "" {
		fmt.Print("2-factor auth code: ")

		var code string
		fmt.Scanf("%s", &code)

		req.Header.Add("X-GitHub-OTP", code)

		resp, err = client.Do(req)
		if err != nil {
			return "", err
		}
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 201 {
		var e map[string]json.RawMessage
		if err := json.Unmarshal(body, &e); err != nil {
			return "", err
		}

		message := string(e["message"])

		if e["errors"] != nil {
			var ee []map[string]string
			if err := json.Unmarshal(e["errors"], &ee); err != nil {
				return "", err
			}
			message += fmt.Sprintf(" (%s %s)", ee[0]["resource"], ee[0]["code"])
		}

		return "", errors.New(message)
	}

	var auth AuthorizationResponse

	d := json.NewDecoder(strings.NewReader(string(body)))
	if err := d.Decode(&auth); err != nil {
		return "", err
	}

	return auth.Token, nil
}
