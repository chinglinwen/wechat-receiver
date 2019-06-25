package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

// doc: https://work.weixin.qq.com/api/doc#10514

var (
	CorpId         = "ww89720c104a10253f" // 企业微信 corpid
	Token          = "wjstHpLmVMj"
	EncodingAESKey = "y4r70uH4aRkSXhfNaKXdbien8zmnMa8xmKl5bm9Il6m"
)

func main() {
	log.Println("starting...")

	http.Handle("/", http.HandlerFunc(receive))
	err := http.ListenAndServe(":1323", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func receive(w http.ResponseWriter, r *http.Request) {
	log.Println(formatRequest(r))

	body, _ := ioutil.ReadAll(r.Body)
	fmt.Println("body:", string(body))

	msg, err := decodeURI(r.RequestURI)
	if err != nil {
		log.Println("decodeuri err, uri: ", err, r.RequestURI)
		return
	}
	spew.Dump(msg)

	if msg.echostr != "" {
		text, err := msg.verifymsg()
		if err != nil {
			log.Println("decrypt err", err)
			return
		}
		fmt.Println("text:", text)
		w.Write([]byte(text))

		return
	}

	c, err := msg.decodeBody([]byte(body))
	if err != nil {
		log.Println("decodeBody err, uri: ", err, body)
		return
	}
	fmt.Printf("got: %#v\n", c)
}

func pretty(a interface{}) {
	b, _ := json.MarshalIndent(a, "", "  ")
	fmt.Println("pretty", string(b))
}

// formatRequest generates ascii representation of a request
func formatRequest(r *http.Request) string {
	log.Printf("r: %#v", r)
	// Create return string
	var request []string
	// Add the request string
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	request = append(request, url)
	// Add the host
	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	// Loop through headers
	for name, headers := range r.Header {
		name = strings.ToLower(name)
		for _, h := range headers {
			request = append(request, fmt.Sprintf("%v: %v", name, h))
		}
	}

	// If this is a POST, add post data
	if r.Method == "POST" {
		r.ParseForm()
		request = append(request, "\n")
		request = append(request, r.Form.Encode())
	}

	// Return the request as a string
	return strings.Join(request, "\n")
}
