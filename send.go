// send to wechat-notify service
package main

import (
	"flag"
	"fmt"
	"log"
	"regexp"
	"strings"

	resty "gopkg.in/resty.v1"
)

var (
	wechatNotifyURL = flag.String("w", "http://localhost:8001", "wechat notify service url")
	receiver        = flag.String("r", "", "default wechat receiver")
	receiverParty   = flag.String("party", "10", "default receiver party ( eg. 3 )")
	agentid         = flag.String("agentid", "", "commander agentid ( eg. 1000003 )")
	secret          = flag.String("secret", "", "commander secret ( eg. G5h7CTEqkBw-Fe3luf2JM8UNNJAcYTpbXvpveY7M3lg )")
	agentidDev      = flag.String("appAgentid", "", "agentid for dev ( eg. 1000003 )")
	secretDev       = flag.String("appSecret", "", "secret for dev ( eg. G5h7CTEqkBw-Fe3luf2JM8UNNJAcYTpbXvpveY7M3lg )")

	expire = flag.String("e", "10m", "default expire time duration")
)

func validate() (err error) {
	if *agentid == "" {
		return fmt.Errorf("agentid is empty")
	}
	if *secret == "" {
		return fmt.Errorf("secret is empty")
	}
	if *agentidDev == "" {
		return fmt.Errorf("agentidDev is empty")
	}
	if *secretDev == "" {
		return fmt.Errorf("secretDev is empty")
	}
	return nil
}

type sendconfig struct {
	touser   string
	toparty  string
	exceptme string
	app      string // commmander or dev
}

type sendoption func(*sendconfig)

// both touser and toparty
func SetReceiver(receiver string) sendoption {
	return func(c *sendconfig) {
		if regexp.MustCompile(`^[0-9]+$`).MatchString(receiver) {
			c.toparty = receiver
			c.touser = ""
			return
		}
		c.touser = receiver
		c.toparty = ""
	}
}

func SetApp(app string) sendoption {
	return func(c *sendconfig) {
		c.app = app
	}
}

func SetExceptMe(exceptme string) sendoption {
	return func(c *sendconfig) {
		c.exceptme = exceptme
	}
}
func SendPerson(message, person string) (reply string, err error) {
	return Send(message, SetReceiver(person))
}

func Send(message string, options ...sendoption) (reply string, err error) {
	c := &sendconfig{
		touser:  *receiver,
		toparty: *receiverParty,
	}
	for _, option := range options {
		option(c)
	}
	if c.touser == "" && c.toparty == "" {
		err = fmt.Errorf("empty user and party, will not send")
		return
	}

	// default to commander
	id := *agentid
	sec := *secret
	if c.app == devApp {
		log.Println("send to devapp")
		id = *agentidDev
		sec = *secretDev
	}

	if id == "" {
		err = fmt.Errorf("agentid is empty")
		return
	}
	if sec == "" {
		err = fmt.Errorf("secret is empty")
		return
	}

	// now := time.Now().Format("2006-1-2 15:04:05")
	// precontent := fmt.Sprintf("时间: %v\n", now)

	r := strings.NewReplacer("\"", " ", "{", "", "}", "")
	message = r.Replace(message)

	resp, e := resty. //SetDebug(true).
				R().
				SetQueryParams(map[string]string{
			"user":    c.touser,
			"toparty": c.toparty,
			"agentid": id,
			"secret":  sec,
			// "precontent": precontent,
			"content":  message,
			"expire":   "3s", // set to short
			"exceptme": c.exceptme,
		}).
		Get(*wechatNotifyURL)

	if e != nil {
		err = e
		return
	}
	reply = string(resp.Body())
	return
}
