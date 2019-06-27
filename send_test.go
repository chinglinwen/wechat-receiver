package main

import (
	"fmt"
	"testing"
)

func TestSend(t *testing.T) {
	out, err := Send("hello", SetApp(devApp), SetReceiver(getGroupWithMe("wen")))
	if err != nil {
		t.Error("send err", err)
		return
	}
	fmt.Println(out)
}
