package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"regexp"
	"strings"

	resty "gopkg.in/resty.v1"
)

var (
	backendURL = flag.String("backendurl", "http://localhost:4000", "backend url")
)

func iscmd(cmd string) bool {
	return regexp.MustCompile(`^/`).MatchString(cmd)
}

// send request to wxrobot-backend
// make it compatible with wxrobot
func sendcmd(name, cmd string) (reply string, err error) {
	resp, e := resty. //SetDebug(true).
				R().
				SetQueryParams(map[string]string{
			"from": name,
			"cmd":  cmd,
		}).
		Get(*backendURL)
	if e != nil {
		err = e
		return
	}
	return parseBody(resp.Body())
}

func parseBody(body []byte) (reply string, err error) {
	r := &Result{}
	err = json.Unmarshal(body, r)
	if err != nil {
		err = fmt.Errorf("unmarshal result err: %v, body: %v", err, string(body))
		return
	}
	if r.Error != "" {
		err = fmt.Errorf("cmd err: %v", strings.TrimSuffix(r.Error, "\n"))
		return
	}
	reply = strings.TrimSuffix(r.Data, "\n")
	return
}

type Result struct {
	Type  string `json:"type"`
	Data  string `json:"data"`
	Error string `json:"error"`
}
