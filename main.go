package main

import (
	"flag"
	"log"
	"net/http"
)

// doc: https://work.weixin.qq.com/api/doc#10514

var (
	port           = flag.String("p", ":8002", "listening address with port")
	CorpId         = "ww89720c104a10253f" // 企业微信 corpid
	Token          = "wjstHpLmVMj"
	EncodingAESKey = "y4r70uH4aRkSXhfNaKXdbien8zmnMa8xmKl5bm9Il6m"
)

const (
	commanderAgentID = 1000005
	devAgentID       = 1000006
	commanderApp     = "commander"
	devApp           = "dev"
)

func main() {
	flag.Parse()
	log.Println("starting...")

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
