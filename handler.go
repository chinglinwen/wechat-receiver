package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func sendFromDevHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name") // name can be group, default to 10 the commander
	content := r.FormValue("content")
	// withops := r.FormValue("ops") // send to ops too

	name = convert(name)

	if content == "" {
		err := fmt.Errorf("content is empty")
		E(w, err)
		return
	}

	var reply string
	var err error
	if name != "" {
		reply, err = Send(content, SetApp(devApp), SetReceiver(getGroupWithMe(name)))
	} else {
		reply, err = Send(content, SetApp(devApp), SetReceiver(getGroupOps()))
	}
	if err != nil {
		err = fmt.Errorf("send to %v, err: %v, reply: %v\n", name, err, reply)
		E(w, err)
		return
	}
	log.Printf("send to dev %v ok, reply: %q\n", name, reply)
	fmt.Fprintf(w, "send to dev %v ok, reply: %q\n", name, reply)
}

// curl -v localhost:8080/send?content=hello
// curl -v "localhost:8080/send?content=hello&name=10"  // send to group 10
func sendHandler(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name") // name can be group, default to 10 the commander
	content := r.FormValue("content")

	if content == "" {
		err := fmt.Errorf("content is empty")
		E(w, err)
		return
	}

	var reply string
	var err error
	if name != "" {
		reply, err = Send(content, SetReceiver(name))
	} else {
		reply, err = Send(content)
	}
	if err != nil {
		err = fmt.Errorf("send to %v, err: %v, reply: %v\n", name, err, reply)
		E(w, err)
		return
	}
	log.Printf("send to %v ok, reply: %q\n", name, reply)
	fmt.Fprintf(w, "send to %v ok, reply: %q\n", name, reply)
}

func E(w http.ResponseWriter, err error) {
	log.Println(err)
	fmt.Fprintf(w, "%v\n", err)
}

func receiveHandler(w http.ResponseWriter, r *http.Request) {
	// log.Println(formatRequest(r))
	body, _ := ioutil.ReadAll(r.Body)
	// fmt.Println("body:", string(body))

	msg, err := decodeURI(r.RequestURI)
	if err != nil {
		err = fmt.Errorf("decodeuri err: %v, uri: %v", err, r.RequestURI)
		E(w, err)
		return
	}
	// spew.Dump(msg)

	if msg.echostr != "" {
		text, err := msg.verifymsg()
		if err != nil {
			err = fmt.Errorf("decrypt err %v", err)
			E(w, err)
			return
		}
		fmt.Println("text:", text)
		w.Write([]byte(text))
		return
	}

	c, err := msg.decodeBody([]byte(body))
	if err != nil {
		err = fmt.Errorf("decodeBody err: %v, body: %v", err, body)
		E(w, err)
		return
	}
	fmt.Printf("\n\nnewmsg from: %v, text: %v\n", c.FromUsername, c.Content)

	// https://work.weixin.qq.com/api/doc#10514
	// need to return in 5 seconds, so to avoid later resend
	go func() {
		if c.Agentid == commanderAgentID {
			log.Println("it's for commanderApp")
			runCommander(w, c)
		}
		if c.Agentid == devAgentID {
			log.Println("it's for devApp")
			runDev(w, c)
		}
	}()

	fmt.Fprintf(w, "received ok\n")
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
