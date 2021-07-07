package model

import (
	"encoding/json"
	"github.com/wonderivan/logger"
)


type WechatToken struct {
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

var AccessToken string

func ( t *WechatToken ) Init () {
	method := "gettoken?corpid=" + wxConfig.CorpID + "&corpsecret=" + wxConfig.CorpSecret
	var w WechatAPI
	body := w.Get(method)
	err := json.Unmarshal(body, t)
	if err != nil {
		logger.Info(err)
	}
	if t.ErrCode == 0 {
		AccessToken = t.AccessToken
		logger.Info("企业微信Token初始化状态：OK",AccessToken)
	} else {
		logger.Info(t.ErrMsg)
	}
	return
}
