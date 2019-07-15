package main

import "log"

var converts = map[string]string{
	"wenzhenglin": "wen",
	"xiaoli":      "robot",
}

func getback(name string) string {
	for k, v := range converts {
		if v == name {
			log.Printf("convertback name %v -> %v\n", name, k)
			return k
		}
	}
	return ""
}

// make this into project config?
func convert(name string) string {
	// if name == "wenzhenglin" {
	// 	return "wen"
	// }
	if v, ok := converts[name]; ok {
		log.Printf("convert name %v -> %v\n", name, v)
		return v
	}
	return name
}

// make this into project config?
func convertback(name string) string {
	// if name == "wen" {
	// 	return "wenzhenglin"
	// }
	if v := getback(name); v != "" {
		return v
	}
	return name
}
