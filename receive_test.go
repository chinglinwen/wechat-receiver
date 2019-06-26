package main

import (
	"fmt"
	"testing"
)

func TestDecodeBody(t *testing.T) {
	uri := "/?msg_signature=586e9e56ed7061e48690e7a918732138c935b4dc&timestamp=1561451603&nonce=1561597366"
	msg, _ := decodeURI(uri)

	body := "<xml><ToUserName><![CDATA[ww89720c104a10253f]]></ToUserName><Encrypt><![CDATA[u7hzodaBGXQY25hFBHAHZcYaP9ReHZnQhy9JooOAMrxQ1KT3Fk7L4u41kfRL5WUVbNpiLEY63Twgqk0ZbjUnupQuGO6r8VojI7zb6hBUtyoPYkuAfbMlg/9vLUVmxfKjq/zXq9f536eDMtgGn5GJtOs+L+zgfo1ezn1ViwfRhPwTOumpZWy0mU0rHlrfgokMNjhYtYIxLaCd/YTqUSSsXGeQQlkgmpShfTJ0/c1Td0mfvpZ4mtctKXWmYo+zRHpN7cC9XjnV+VTEiYRCGTtfEjBOTSLAIM49fqDt272q2G2Za0OFdg13WjccD9jB1fuI3YhkOy/zqlME3QG3b2OElznWE7/XjZcPp6dO1zrAM7WkTYrbtNwmOtgu783qWNIwQwx+sG646+3cpVIwkPbxk+0dN5jw5KsTfTMmhoSS6gc=]]></Encrypt><AgentID><![CDATA[1000003]]></AgentID></xml>"

	c, err := msg.decodeBody([]byte(body))
	if err != nil {
		t.Error("decodeBody err", err)
	}
	if c.Content != "he1" {
		t.Error("content err, expect he1, got", c.Content)
		return
	}
}

func TestDecodeURI(t *testing.T) {

	// baduri := "?msg_signature=fe45e0a29166411e1998e292c5498dcb327eb7fa&timestamp=1561444015&nonce=1561603193"
	// _, err := decodeURI(baduri)
	// if err == nil {
	// 	t.Errorf("uri should be err, got nil")
	// 	return
	// }

	// uri := "/?msg_signature=3d2a223ea8c16734138bee502c35166ed5dd6004&timestamp=1561432777&nonce=1561628896&echostr=epa%2Bq0UdhNUr1B2OPCx3DCssoRy9lHKYzyU%2FqQwbUpGhebQMqT8cjURj4wBCG68Ra81d8LybbqHWOZrBRYaWMw%3D%3D"

	uri := "/?msg_signature=732f718eb2b8ac89c6349034219b1377ac8fa5b0&timestamp=1561443978&nonce=1562067257&echostr=g15JwZFxmBp9I1qM30FQqVDYVzQyMJjaXQQRJOyR%2F3%2BkVm4i0%2B9UKJ66xryguAlOIkF%2FE0winfMP2OWgjegtsQ%3D%3D"

	// uri := "/?msg_signature=3afa129a9c27bfb544c8fcace070acfc90d69b10&timestamp=1561445432&nonce=1561461235"
	msg, _ := decodeURI(uri)

	text, err := msg.verifymsg()
	if err != nil {
		t.Error("decrypt err", err)
		return
	}
	fmt.Println("text:", text)
}
