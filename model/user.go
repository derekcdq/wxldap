package model

import (
	"encoding/json"
	"github.com/go-ldap/ldap"
	"github.com/wonderivan/logger"
	"sort"
	"strconv"
)

type UserList struct {
	Errcode  int    `json:"errcode"`
	Errmsg   string `json:"errmsg"`
	Userlist []UserInfo
}

type UserInfo struct {
	UserID          string  `json:"userid"`
	Name            string  `json:"name"`
	Department      []int64 `json:"department"`
	Position        string  `json:"position"`
	Email           string  `json:"email"`
	Status          int     `json:"status"`
	Main_department int     `json:"main_department"`
	LdapExist       int
	UidNumber       string
}

func ( t *UserList) Read ( body []byte ) interface{} {
	err := json.Unmarshal(body, t)
	if err != nil {
		logger.Info(err)
	}
	return t
}

func ( t *UserList ) Get (deptID int) []UserInfo {
	method := "user/list?access_token=" + AccessToken + "&department_id=" + strconv.Itoa(deptID) + "&fetch_child=0"
	body := CallWechatApi(method)
	t.Read(body)
	return t.Userlist
}

func  ( t *UserInfo ) Read ( body []byte ) interface{} {
	err := json.Unmarshal(body, t)
	if err != nil {
		logger.Info(err)
	}
	return t
}

func ( t *UserInfo ) CheckExist ( userID string ) interface{} {
	filter := "(&(" + "uid=" + userID + "))"
	sql := ldap.NewSearchRequest("dc=eclincloud,dc=net", ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter, []string{"dn", "cn", "uid"}, nil)
	sr, _ := ldapConn.Search(sql)
	if len(sr.Entries) > 0 {
		t.LdapExist = 1
	} else {
		t.LdapExist = 0
	}
	return t
}

func ( t *UserInfo ) GetUidNumber () interface{} {
	uidSql := ldap.NewSearchRequest("dc=eclincloud,dc=net", 3, 0, 0, 0, false,
		"(&(objectClass=posixAccount))", []string{"uidNumber", "dn"}, nil)
	sr, err := ldapConn.Search(uidSql)
	if err != nil {
		logger.Info(err)
	}
	var array []string
	for _, value := range sr.Entries {
		array = append(array, value.GetAttributeValue("uidNumber"))
	}
	sort.Sort(sort.Reverse(sort.StringSlice(array)))
	maxUidNumber, err := strconv.Atoi(array[0])
	if err != nil {
		logger.Info(err)
	}
	t.UidNumber = strconv.Itoa(maxUidNumber + 1)
	return t
}

func ( t *UserInfo) AddToLdap ( userName string, userID string, userEmail string, userDn string ) interface{} {
	t.CheckExist(userID)
	if t.LdapExist == 1 {
		//logger.Info("用户已存在:" + userDn)
		return t
	}
	t.GetUidNumber()
	logger.Info("用户不存在:" + userDn, "\t开始添加")
	sql := ldap.NewAddRequest(userDn, nil)
	sql.Attribute("sn", []string{userName})
	sql.Attribute("cn", []string{userName})
	sql.Attribute("uid", []string{userID})
	sql.Attribute("userPassword", []string{ldapConfig.DefaultPassword})
	sql.Attribute("uidnumber", []string{t.UidNumber})
	sql.Attribute("gidNumber", []string{"500"})
	sql.Attribute("homedirectory", []string{"/home/users/" + userID})
	sql.Attribute("mail", []string{userEmail})
	sql.Attribute("objectClass", []string{"inetOrgPerson", "posixAccount"})
	err := ldapConn.Add(sql)
	if err != nil {
		logger.Info(err)
	} else {
		m := new(NotifyMail)
		m.Send(userID , userName, userEmail)
		s := new(NotifyMsg)
		s.Send(userID)
	}
	return t
}

func SyncAllUser(userList []UserInfo, dn string) {
	if len(userList) < 1 {
		return
	}
	for _, v := range userList {
		if v.Status == 1 {
			userDn := "cn=" + v.Name + "," + dn
			u := UserInfo{}
			u.AddToLdap(v.Name,v.UserID,v.Email, userDn)
		}
	}
}



