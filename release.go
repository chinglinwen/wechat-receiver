package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"strings"

	resty "gopkg.in/resty.v1"
)

var (
	releaseURL = flag.String("releaseurl", "http://release.haodai.net/api/wechat", "release url")
)

// send request to self-release
// make it compatible with wxrobot
func sendRelease(name, cmd string) (reply string, err error) {
	resp, e := resty. //SetDebug(true).
				SetRedirectPolicy(resty.FlexibleRedirectPolicy(20)).
				R().
				SetQueryParams(map[string]string{
			"from": convertback(name),
			"cmd":  cmd,
		}).
		Get(*releaseURL)
	if e != nil {
		err = e
		return
	}
	return parseReleaseBody(resp.Body())
}
func parseReleaseBody(body []byte) (reply string, err error) {
	r := &Result{}
	err = json.Unmarshal(body, r)
	if err != nil {
		err = fmt.Errorf("unmarshal result err: %v, body: %v", err, string(body))
		return
	}
	if r.Error != "" {
		err = fmt.Errorf("err: %v", strings.TrimSuffix(r.Error, "\n"))
		return
	}
	reply = strings.TrimSuffix(r.Data, "\n")
	return
}
