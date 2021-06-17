package model

import (
	"encoding/json"
	"github.com/go-ldap/ldap"
	"github.com/wonderivan/logger"
)

type DeptList struct {
	Errcode    int          `json:"errcode"`
	Errmsg     string       `json:"errmsg"`
	Department []DeptInfo   `json:"department"`
}

type DeptInfo struct {
	Name     string `json:"name"`
	ID       int    `json:"id"`
	DN       string `json:"dn"`
	ParentID int    `json:"parentid"`
	LdapExist int
}


func ( t *DeptList ) Read ( body []byte ) interface{} {
	err := json.Unmarshal(body, t)
	if err != nil {
		logger.Info(err)
	}
	return t
}

func (t *DeptList ) Get () interface{} {
	var departments []DeptInfo
	method := "department/list?access_token=" + AccessToken
	body := CallWechatApi(method)
	t.Read(body)
	deptsMap := make(map[int]string)
	for _, value := range t.Department {
		deptsMap[value.ID] = value.Name
	}
	for _, value := range t.Department {
		n := DeptInfo{}
		var dn string
		if value.ParentID <= 1 {
			dn = "ou=" + value.Name + ",dc=eclincloud,dc=net"
		} else {
			parentName := deptsMap[value.ParentID]
			dn = "ou=" + value.Name + ",ou=" + parentName + ",dc=eclincloud,dc=net"
		}
		n.ID = value.ID
		n.Name = value.Name
		n.DN = dn
		departments = append(departments,n)
	}
	t.Department =  departments
	return t
}

func ( t *DeptInfo ) CheckExist (deptName string ) interface{} {
	filter := "(&(" + "ou=" + deptName + "))"
	sql := ldap.NewSearchRequest("dc=eclincloud,dc=net",
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false, filter, []string{"dn", "cn", "uid"}, nil)
	sr, _ := ldapConn.Search(sql)
	status := len(sr.Entries)
	t.LdapExist = status
	return t
}

func ( t *DeptInfo ) AddToLdap ( deptName string ,dn string ) interface{} {
	t.CheckExist(deptName)
	if t.LdapExist == 0 {
		sql := ldap.NewAddRequest(dn, nil)
		sql.Attribute("objectClass", []string{"organizationalUnit", "top"})
		err := ldapConn.Add(sql)
		if err != nil {
			logger.Info(err)
		} else {
			logger.Info("同步" + dn + "成功")
		}
	}
	return t
}

func SyncAllDept () {
	deptList := new (DeptList)
	deptList.Get()
	if len(deptList.Department) < 1 {
		return
	}
	logger.Info("开始同步部门及组...")
	for _, v := range deptList.Department {
		logger.Info("开始同步:",v.Name)
		d := DeptInfo{}
		d.AddToLdap(v.Name, v.DN)
		u := new (UserList)
		u.Get(v.ID)
		SyncAllUser(u.Userlist,v.DN)
	}
}