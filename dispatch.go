package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func runCommander(w http.ResponseWriter, c MsgContent) {
	// send normal chat back to other member
	reply, err := Send(fmt.Sprintf("%v says: \n---\n%v", c.FromUsername, c.Content), SetApp(commanderApp), SetExceptMe(c.FromUsername))
	if err != nil {
		err = fmt.Errorf("forward from %v, err: %v, reply: %v\n", c.FromUsername, err, reply)
		log.Println(err)
	} else {
		log.Printf("forward from %v ok, reply: %q\n", c.FromUsername, reply)
	}
	// if it's a command send it to commander service
	if !iscmd(c.Content) {
		return
	}

	// send cmd to backend
	output, err := sendcmd(c.FromUsername, c.Content)
	if err != nil {
		// err = fmt.Errorf("sendcmd from %v, err: %v, reply: %v\n", c.FromUsername, err, output)
		// E(w, err)
		// return
		output = err.Error()
	}
	log.Printf("sendcmd from %v ok, reply: %q\n", c.FromUsername, output)

	// send result back to chat
	reply, err = Send(fmt.Sprintf("results: \n---\n%v", output), SetApp(devApp))
	if err != nil {
		err = fmt.Errorf("sendresult to %v, err: %v, reply: %v\n", c.FromUsername, err, reply)
		E(w, err)
		return
	}
	log.Printf("sendresult to %v ok, reply: %q\n", c.FromUsername, reply)

	return
}

// an virtual group serve as help for dev app
var ops = []string{"wen"}

func isops(name string) bool {
	for _, v := range ops {
		if v == name {
			return true
		}
	}
	return false
}
func getGroupExceptMe(name string) string {
	name = convert(name)
	ops1 := append(ops, name)
	ops2 := []string{}
	for _, v := range ops1 {
		if v == name {
			continue
		}
		ops2 = append(ops2, v)
	}
	return strings.Join(ops2, "|")
}

func getGroupWithMe(name string) string {
	name = convert(name)
	if isops(name) {
		return strings.Join(ops, "|")
	}
	return strings.Join(append(ops, name), "|")
}

func getGroupOps() string {
	return strings.Join(ops, "|")
}

// for devapp send msg to yunwei ( we need yunwei's person list )
// anyone in yuwei can reply to that person
// we need extra except me, and can't use party as target
func runDev(w http.ResponseWriter, c MsgContent) {
	name := convertback(c.FromUsername)

	if c.Content == "" {
		log.Printf("got empty content from %v, skip", name)
		return
	}
	// send normal chat back to other member
	g := getGroupExceptMe(c.FromUsername)
	if g != "" {
		reply, err := Send(fmt.Sprintf("%v says: \n---\n%v", name, c.Content), SetApp(devApp), SetReceiver(g))
		if err != nil {
			err = fmt.Errorf("forward from %v, err: %v, reply: %v\n", name, err, reply)
			log.Println(err)
		} else {
			log.Printf("forward from %v ok, reply: %q\n", name, reply)
		}
	} else {
		log.Println("no others member to send, skip forward")
	}

	// if it's a command send it to commander service
	if !iscmd(c.Content) {
		log.Println("not a command, no further steps")
		return
	}

	// send cmd to backend
	output, err := sendRelease(name, c.Content)
	if err != nil {
		// err = fmt.Errorf("sendcmd from %v, err: %v, reply: %v\n", c.FromUsername, err, output)
		// E(w, err)
		// return
		output = err.Error()
	}
	log.Printf("sendRelease from %v ok, reply: %q\n", name, output)

	// send result back to chat
	reply, err := Send(fmt.Sprintf("results: \n---\n%v", output), SetApp(devApp), SetReceiver(getGroupWithMe(name)))
	if err != nil {
		err = fmt.Errorf("sendresult to %v, err: %v, reply: %v\n", name, err, reply)
		E(w, err)
		return
	}
	log.Printf("sendresult to %v ok, reply: %q\n", name, reply)

	return
}
