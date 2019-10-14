package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

// doc: https://work.weixin.qq.com/api/doc#10514

var (
	port           = flag.String("p", ":8002", "listening address with port")
	CorpId         = flag.String("coprid", "", "corp id") // 企业微信 corpid
	Token          = flag.String("token", "", "app Token")
	EncodingAESKey = flag.String("aeskey", "", "app EncodingAESKey")

	// commanderAgentID = flag.Int("commanderAgentID", 1000005, "commanderAgentID")
	// devAgentID       = flag.Int("devAgentID", 1000006, "devAgentID")

	// for send
	wechatNotifyURL = flag.String("w", "http://localhost:8001", "wechat notify service url")
	receiver        = flag.String("r", "", "default wechat receiver")

	receiverParty = flag.String("party", "10", "default receiver party ( eg. 3 )")
	agentid       = flag.String("agentid", "", "commander agentid ( eg. 1000003 )")
	secret        = flag.String("secret", "", "commander secret ( eg. G5h7CTEqkBw-Fe3luf2JM8UNNJAcYTpbXvpveY7M3lg )")

	receiverPartyDev = flag.String("partyDev", "10", "default receiver party for dev( eg. 3 )")
	agentidDev       = flag.String("agentidDev", "", "agentid for dev ( eg. 1000003 )")
	secretDev        = flag.String("secretDev", "", "secret for dev ( eg. G5h7CTEqkBw-Fe3luf2JM8UNNJAcYTpbXvpveY7M3lg )")

	expire = flag.String("e", "10m", "default expire time duration")
)

const (
	commanderApp = "commander"
	devApp       = "dev"
)

func validateMain() (err error) {
	if *CorpId == "" {
		return fmt.Errorf("CorpId is empty")
	}
	if *Token == "" {
		return fmt.Errorf("app Token is empty")
	}
	if *EncodingAESKey == "" {
		return fmt.Errorf("app EncodingAESKey is empty")
	}
	return
}
func main() {
	flag.Parse()
	log.Println("starting...")

	if err := validateMain(); err != nil {
		log.Fatal(err)
	}
	if err := validate(); err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/", receiveHandler)
	http.HandleFunc("/send", sendHandler)
	http.HandleFunc("/dev", sendFromDevHandler)
	err := http.ListenAndServe(*port, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
