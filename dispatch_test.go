package main

import "testing"

func TestGetGroupWithMe(t *testing.T) {
	if a := getGroupWithMe("wen"); a != "wen" {
		t.Error("err, expect wen, got ", a)
	}
}

func TestGetGroupExceptMe(t *testing.T) {
	if a := getGroupExceptMe("wen"); a != "" {
		t.Error("err, expect empty, got ", a)
	}
}

func TestIsops(t *testing.T) {
	if !isops("wen") {
		t.Error("expect true, got ", false)
	}
}
