package model

import (
	"github.com/go-ldap/ldap"
	"github.com/wonderivan/logger"
	"time"
)

var ldapConn *ldap.Conn

func InitLdap() {
	var err error
	ldapConfig := new(LdapConfig)
	ldapConfig.Init()
	logger.Info("开始初始化Ldap服务器连接......")
	Loop1:
		for {
			ldapConn, err = ldap.Dial("tcp", ldapConfig.LdapHost)
			if err != nil {
				logger.Info(err)
			} else {
				err = ldapConn.Bind(ldapConfig.RootDn, ldapConfig.Password)
				if err != nil {
					logger.Info(err)
				} else {
					logger.Info("Ldap服务器连接状态：OK")
					break Loop1
				}
			}
			time.Sleep(time.Duration(10) * time.Second)
		}
}

func ReConnLdap() {
	if ldapConn.IsClosing() {
		logger.Info("LDAP服务器连接丢失,重新连接...")
		InitLdap()
	}
}

