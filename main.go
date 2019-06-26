package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

// doc: https://work.weixin.qq.com/api/doc#10514

var (
	port           = flag.String("p", ":8080", "listening address with port")
	CorpId         = "ww89720c104a10253f" // 企业微信 corpid
	Token          = "wjstHpLmVMj"
	EncodingAESKey = "y4r70uH4aRkSXhfNaKXdbien8zmnMa8xmKl5bm9Il6m"

	commanderAgentID = 1000005
)

func main() {
	flag.Parse()
	log.Println("starting...")

	http.HandleFunc("/", receiveHandler)
	http.HandleFunc("/send", sendHandler)
	err := http.ListenAndServe(*port, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
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
		log.Printf("send to %v, err: %v, reply: %v\n", name, reply)
		return
	}
	log.Printf("send to %v ok, reply: %q\n", name, reply)
}

func E(w http.ResponseWriter, err error) {
	log.Println(err)
	fmt.Fprintf(w, "%v\n", err)
}

func receiveHandler(w http.ResponseWriter, r *http.Request) {
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

	if c.Agentid != commanderAgentID {
		return
	}
	reply, err := Send(fmt.Sprintf("%v says: \n---\n%v", c.FromUsername, c.Content), SetExceptMe(c.FromUsername))
	if err != nil {
		log.Printf("forward from %v, err: %v, reply: %v\n", c.FromUsername, reply)
		return
	}
	log.Printf("forward from %v ok, reply: %q\n", c.FromUsername, reply)
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
