package main

// make this into project config?
func convert(name string) string {
	if name == "wenzhenglin" {
		return "wen"
	}
	return name
}

// make this into project config?
func convertback(name string) string {
	if name == "wen" {
		return "wenzhenglin"
	}
	return name
}
