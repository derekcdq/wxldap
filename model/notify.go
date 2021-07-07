package model

import (
	"github.com/wonderivan/logger"
	"io/ioutil"
	"net/smtp"
	"strings"
)

type NotifyMail struct {
	Name     string
	UserID   string
	UserMail string
	Template string
	Content  string
}

type NotifyMsg struct {
	Name     string
	UserID   string
	Template string
	MsgText  string
}

func ( t *NotifyMail ) GetContent ( userID string, userName string ) interface{} {
	template, err := ioutil.ReadFile("../templates/ldapMail.tpl")
	content := string(template)
	content = strings.Replace(content, "{name}", userName, -1)
	content = strings.Replace(content, "{userid}", userID, -1)
	if err != nil {
		logger.Info(err)
	}
	t.UserID = userID
	t.Name   = userName
	t.Content = content
	return t
}

func ( t *NotifyMail ) Send ( userID string, userName string ,userEmail string  ) {
	t.GetContent(userID,userName)
	to := []string{userEmail}
	msg := []byte("From: LDAP账号管理员\r\n" + "Subject:LDAP账号开通提醒\r\n" + "\r\n" + t.Content)
	auth := smtp.PlainAuth("", smtpConfig.MailUser, smtpConfig.Password, smtpConfig.MailHost)
	Loop1:
		for i := 0; i <= 3; i++ {
			err := smtp.SendMail(smtpConfig.MailHost+":587", auth, smtpConfig.MailUser, to, msg)
			if err != nil {
				logger.Info(err)
			} else {
				logger.Info("邮件发送成功")
				break Loop1
			}
		}
}

func ( t *NotifyMsg ) Send ( userID string ) interface{} {
	method := "message/send?access_token="+ AccessToken
	msgContent := "您在易临云的LDAP账号已经自动开通\n账号:"+userID+" \n初始密码:ecc123456 \n为了账号安全,请即时登陆http://password.eclincloud.net修改密码"
	body, _ := ioutil.ReadFile("templates/ldapMsg.tpl")
	msgText := string(body)
	msgText = strings.Replace(msgText, "{userId}", userID, -1)
	msgText = strings.Replace(msgText, "{msgContent}", msgContent, -1)
	t.UserID = userID
	t.MsgText = msgText
	var w WechatAPI
	w.Post(method, msgText)
	return t
}



