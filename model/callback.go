package model

import (
	"encoding/xml"
	"fmt"
	"github.com/sbzhu/weworkapi_golang/wxbizmsgcrypt"
	"io/ioutil"
	"net/http"
)

type MsgContent struct {
	ToUsername   string `xml:"ToUserName"`
	FromUsername string `xml:"FromUserName"`
	CreateTime   uint32 `xml:"CreateTime"`
	MsgType      string `xml:"MsgType"`
	Content      string `xml:"Content"`
	Msgid        string `xml:"MsgId"`
	Agentid      uint32 `xml:"AgentId"`
	ChangeType   string `xml:"ChangeType"`
}


func IndexHandler( w http.ResponseWriter, r *http.Request) {
	wxcpt := wxbizmsgcrypt.NewWXBizMsgCrypt(callbackConfig.Token, callbackConfig.EncodingAeskey, callbackConfig.ReceiverId, wxbizmsgcrypt.XmlType)
	r.ParseForm()
	msgSignature := r.Form.Get("msg_signature")
	timestamp := r.Form.Get("timestamp")
	nonce     := r.Form.Get("nonce")
	s, _ := ioutil.ReadAll(r.Body)
	fmt.Println(msgSignature,timestamp,nonce)
	msg, cryptErr := wxcpt.DecryptMsg(msgSignature, timestamp, nonce, s)
	if cryptErr != nil {
		fmt.Println(cryptErr)
	}
	var m MsgContent
	fmt.Println(string(msg))
	xml.Unmarshal(msg,&m)
	fmt.Println(m.ToUsername,m.Content,m.MsgType,m.Msgid,m.ChangeType)
	deptInfo := new(DeptInfo)
	deptInfo.ChangeDn("ou=公共账号,dc=eclincloud,dc=net","ou=公共账号1,dc=eclincloud,dc=net")
	fmt.Println(deptInfo.DN)

}



