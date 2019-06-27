package main

import (
	"flag"

	resty "gopkg.in/resty.v1"
)

var (
	releaseURL = flag.String("releaseurl", "http://build.newops.haodai.net/api/wechat", "release url")
)

// send request to self-release
// make it compatible with wxrobot
func sendRelease(name, cmd string) (reply string, err error) {
	resp, e := resty. //SetDebug(true).
				SetRedirectPolicy(resty.FlexibleRedirectPolicy(20)).
				R().
				SetQueryParams(map[string]string{
			"from": name,
			"cmd":  cmd,
		}).
		Get(*releaseURL)
	if e != nil {
		err = e
		return
	}
	return parseBody(resp.Body())
}
