package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	neturl "net/url"
	"strings"
)

const GitIOURL = "http://git.io"

func GistParseError(body []byte) error {
	var err map[string]json.RawMessage
	if err := json.Unmarshal(body, &err); err != nil {
		return err
	}

	message := string(err["message"])

	if err["errors"] != nil {
		var errs []map[string]string
		if err := json.Unmarshal(err["errors"], &errs); err != nil {
			return err
		}

		var messages []string
		for _, m := range errs {
			messages = append(messages, fmt.Sprintf("%s %s", m["resource"], m["code"]))
		}
		message += " (" + strings.Join(messages, ", ") + ")"
	}

	return errors.New(message)
}

func Shorten(url string) (string, error) {
	resp, err := http.PostForm(GitIOURL, neturl.Values{"url": {url}})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 201 {
		url = resp.Header.Get("Location")
	}

	return url, nil
}
