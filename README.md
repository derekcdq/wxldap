# wxldap
同步企业微信的组织架构和人员到OPENLDAP


配置文件在conf目录下
ldap.json
修改以下信息，改为你的LDAP的信息

{
  "ldapConfig": {
      "ldapHost": "172.20.11.7:389",
      "baseDn" : "dc=eclincloud,dc=net",
      "rootDn" : "cn=root,dc=example,dc=net",
      "password" : "ldap的管理密码",
      "defaultPassword": "新ldap用户的初始密码"
  }
}

smtp.json
修改以下信息，主要用于通知邮件的发送

{
  "smtpConfig" : {
    "mailHost": "smtp.exmail.qq.com",
    "mailUser"  : "admin@example.com",
    "mailPasswd" : "邮箱密码"
  }
}

wechat.json
修改以下信息，改为你自己在企业微信开放平台的ID和Secret

{
  "wechatConfig": {
    "apiUrl": "https://qyapi.weixin.qq.com/cgi-bin/",
    "corpId": "你在企业微信的企业id",
    "corpSecret": "你的企业微信密钥",
    "interval": 30
  }
}

运行程序：

go run main.go
