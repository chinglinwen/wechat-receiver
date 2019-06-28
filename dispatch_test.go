package main

import "testing"

func TestGetGroupWithMe(t *testing.T) {
	if a := getGroupWithMe("wenzhenglin"); a != "wen" {
		t.Error("err, expect wen, got ", a)
	}
}

func TestGetGroupExceptMe(t *testing.T) {
	if a := getGroupExceptMe("wenzhenglin"); a != "" {
		t.Error("err, expect empty, got ", a)
	}
}

func TestIsops(t *testing.T) {
	if !isops("wen") {
		t.Error("expect true, got ", false)
	}
}
