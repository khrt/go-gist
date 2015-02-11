package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

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
