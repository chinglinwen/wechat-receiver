package main

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"strings"

	"github.com/sbzhu/weworkapi_golang/wxbizmsgcrypt"
)

// https://work.weixin.qq.com/api/doc#90000/90139/90968/%E5%AF%86%E6%96%87%E8%A7%A3%E5%AF%86%E5%BE%97%E5%88%B0msg%E7%9A%84%E8%BF%87%E7%A8%8B
// https://work.weixin.qq.com/api/devtools/devtool.php
type msg struct {
	echostr      string //":[]string{"epa+q0UdhNUr1B2OPCx3DCssoRy9lHKYzyU/qQwbUpGhebQMqT8cjURj4wBCG68Ra81d8LybbqHWOZrBRYaWMw=="}
	msgSignature string //":[]string{"3d2a223ea8c16734138bee502c35166ed5dd6004"},
	nonce        string //":[]string{"1561628896"}
	timestamp    string //":[]string{"1561432777"}}
}

func decodeURI(uri string) (m msg, err error) {
	maps, _ := url.ParseQuery(uri)
	s := make(map[string]string, len(maps))
	for k, v := range maps {
		if strings.Contains(k, "/?") {
			k = strings.TrimPrefix(k, "/?")
			// k1 := strings.TrimPrefix(k, "/?")
			// maps[k1] = v
			// delete(maps, k)
		}
		var v1 string
		if len(v) == 0 {
			v1 = "emptyvalue"
		} else {
			v1 = v[0]
		}
		s[k] = v1
	}
	m = msg{
		echostr:      s["echostr"],
		msgSignature: s["msg_signature"],
		nonce:        s["nonce"],
		timestamp:    s["timestamp"],
	}
	err = m.validate()
	return
}

// no need validate, normal msg does not have echostr
func (m *msg) validate() error {
	if m.msgSignature == "" {
		return fmt.Errorf("msgSignature empty")
	}
	return nil
}

type MsgContent struct {
	ToUsername   string `xml:"ToUserName"`
	FromUsername string `xml:"FromUserName"`
	CreateTime   uint32 `xml:"CreateTime"`
	MsgType      string `xml:"MsgType"`
	Content      string `xml:"Content"`
	Msgid        string `xml:"MsgId"`
	Agentid      int    `xml:"AgentID"`
}

func getwxcpt() *wxbizmsgcrypt.WXBizMsgCrypt {
	return wxbizmsgcrypt.NewWXBizMsgCrypt(*Token, *EncodingAESKey, *CorpId, wxbizmsgcrypt.XmlType)
}

// https://work.weixin.qq.com/api/doc#90000/90135/90930
// https://sourcegraph.com/github.com/sbzhu/weworkapi_golang/-/blob/sample.go#L23:2
func (m *msg) decodeBody(body []byte) (c MsgContent, err error) {
	wxcpt := getwxcpt()
	msg, crypt_err := wxcpt.DecryptMsg(m.msgSignature, m.timestamp, m.nonce, body)
	if nil != crypt_err {
		err = fmt.Errorf("DecryptMsg err %v", err)
		return
	}
	err = xml.Unmarshal(msg, &c)
	if nil != err {
		err = fmt.Errorf("xmlunmarshal err %v", err)
		return
	}
	return
}

// write result back to verify
func (m *msg) verifymsg() (echo string, err error) {
	wxcpt := getwxcpt()
	echobyte, e := wxcpt.VerifyURL(m.msgSignature, m.timestamp, m.nonce, m.echostr)
	if nil != e {
		err = fmt.Errorf("verifyurl failed, err: %v", e)
		return
	}
	echo = string(echobyte)
	return
}

// func (m *msg) decodemsg() (text string, err error) {
// 	wByte, err := base64.StdEncoding.DecodeString(m.echostr)
// 	if err != nil {
// 		return
// 	}
// 	key, err := base64.StdEncoding.DecodeString(EncodingAESKey + "=")
// 	if err != nil {
// 		log.Fatal("EncodingAESKey invalid")
// 	}

// 	x, err := AESDecrypt(wByte, []byte(key))
// 	if err != nil {
// 		err = fmt.Errorf("aes decrypt fail, err: %v, wbyte: %v, key: %v\n", err, wByte, key)
// 		return
// 	}

// 	buf := bytes.NewBuffer(x[16:20])
// 	var length int32
// 	binary.Read(buf, binary.BigEndian, &length)

// 	fmt.Printf("got x: %v len: %v, length: %v\n", string(x), len(x), length)
// 	text = string(x[20 : 20+length])
// 	return

// 	// // verify id
// 	// appIDstart := 20 + length
// 	// if len(x) < int(appIDstart) {
// 	// 	err = errors.New("获取数据错误, 请检查 EncodingAESKey 配置")
// 	// 	return
// 	// }
// 	// id := x[appIDstart : int(appIDstart)+len(corpId)]
// 	// _ = id

// 	// if string(id) == corpId {
// 	// 	text = x[20 : 20+length]
// 	// 	return
// 	// }
// }

// func AESDecrypt(crypted, key []byte) (origData []byte, err error) {
// 	block, err := aes.NewCipher(key)
// 	if err != nil {
// 		err = fmt.Errorf("newcipher err: %v", err)
// 		return
// 	}
// 	blockSize := block.BlockSize()
// 	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
// 	origData = make([]byte, len(crypted))
// 	blockMode.CryptBlocks(origData, crypted)
// 	origData = PKCS5UnPadding(origData)
// 	return
// }

// func PKCS5UnPadding(origData []byte) []byte {
// 	length := len(origData)
// 	unpadding := int(origData[length-1])
// 	return origData[:(length - unpadding)]
// }
