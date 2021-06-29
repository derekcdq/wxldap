package model

import (
	"encoding/xml"
	"fmt"
	"github.com/sbzhu/weworkapi_golang/wxbizmsgcrypt"
	"github.com/wonderivan/logger"
	"io/ioutil"
	"net/http"
)

type MsgContent struct {
	ToUsername   string `xml:"ToUserName"`
	FromUsername string `xml:"FromUserName"`
	CreateTime   uint32 `xml:"CreateTime"`
	MsgType      string `xml:"MsgType"`
	Content      string `xml:"Content"`
	Id           int
	Name         string
	Msgid        string `xml:"MsgId"`
	Agentid      uint32 `xml:"AgentId"`
	ChangeType   string `xml:"ChangeType"`
}


func IndexHandler ( w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	wxcpt := wxbizmsgcrypt.NewWXBizMsgCrypt(callbackConfig.Token, callbackConfig.EncodingAeskey, callbackConfig.ReceiverId, wxbizmsgcrypt.XmlType)
	msgSignature := r.Form.Get("msg_signature")
	timestamp := r.Form.Get("timestamp")
	nonce     := r.Form.Get("nonce")
	s, _ := ioutil.ReadAll(r.Body)
	msg, cryptErr := wxcpt.DecryptMsg(msgSignature, timestamp, nonce, s)
	fmt.Println(string(msg))
	if cryptErr != nil {
		fmt.Println(cryptErr)
	}
	var m MsgContent
	xml.Unmarshal(msg,&m)
	changeType := m.ChangeType
	switch  changeType {
	case "update_user":
		fmt.Println("updateUserDept")
	case "update_party":
		fmt.Println("updateparty")
		if  m.Name != "" {
			UpdateParty(m.Id,m.Name)
		}
	default:
		fmt.Println("pass")
	}

}

func UpdateParty (deptID int, rDn string) {
	InitDmap()
	dn := Dmap.Multiple[deptID]["dn"] + "," + Dmap.Multiple[deptID]["pdn"]
	rDn = "ou=" + rDn
	newSup :=  Dmap.Multiple[deptID]["pdn"]
	d := new(DeptInfo)
	d.ChangeDn(dn,rDn,newSup)
	logger.Info("部门名称变更成功，新名称为:",d.DN)
}





