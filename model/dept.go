package model

import (
	"encoding/json"
	"fmt"
	"github.com/go-ldap/ldap"
	"github.com/wonderivan/logger"
	"strconv"
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
	ParentDN string
	LdapExist int
}

type DeptsMap struct  {
	SimpleMap   map[int]string
	MultipleMap map[int]map[string]string
}

var Dmap DeptsMap

func ( t *DeptsMap ) Init () interface{} {
	method := "department/list?access_token=" + AccessToken
	body := CallWechatApi(method)
	d := new(DeptList)
	err := json.Unmarshal(body, d)
	if err != nil {
		logger.Info(err)
	}
	s := make(map[int]string)
	for _ , v := range d.Department {
		s[v.ID] = v.Name
	}
	m := make(map[int]map[string]string)
	for _, v := range d.Department {
		var (
			myDN string
			parentDN string
		)
		myDN = "ou=" + v.Name
		if v.ParentID <= 1 {
			parentDN = ldapConfig.BaseDn
		} else {
			parentName := s[v.ID]
			parentDN = "ou=" + parentName + "," + ldapConfig.BaseDn
		}
		m[v.ID] = make(map[string]string)
		m[v.ID]["mydn"]  = myDN
		m[v.ID]["parentdn"] = parentDN
		m[v.ID]["parentid"] = strconv.Itoa(v.ParentID)
		m[v.ID]["parentname"] = s[v.ParentID]
	}
	t.SimpleMap = s
	t.MultipleMap = m
	return t
}

func (t *DeptList ) Get () interface{} {
	var departments []DeptInfo
	method := "department/list?access_token=" + AccessToken
	body := CallWechatApi(method)
	err := json.Unmarshal(body, t)
	if err != nil {
		logger.Info(err)
	}
	for _, v := range t.Department {
		var pdn string
		d := DeptInfo{}
		id := v.ID
		name := v.Name
		dn := "ou=" + name
		//循环5次组合DN，支持LDAP里面5层结构
		for i := 1;i <= 5; i++ {
			parentid := Dmap.MultipleMap[id]["parentid"]
			if parentid != "" && parentid != "0" && parentid != "1" {
				parentname := Dmap.MultipleMap[id]["parentname"]
				dn = dn + ",ou=" + parentname
				if ( i == 1 ) {
					pdn = pdn + "ou=" + parentname
				} else {
					pdn = pdn + ",ou=" + parentname
				}
				id,_ = strconv.Atoi(parentid)
			}
		}
		dn = dn + "," + ldapConfig.BaseDn
		pdn = pdn + "," + ldapConfig.BaseDn
		d.ID = v.ID
		d.Name = v.Name
		d.DN = dn
		d.ParentID = v.ParentID
		d.ParentDN = pdn
		departments = append(departments,d)
	}
	t.Department =  departments
	return t
}

func ( t *DeptInfo ) CheckExist (deptName string ) interface{} {
	filter := "(&(" + "ou=" + deptName + "))"
	sql := ldap.NewSearchRequest(ldapConfig.BaseDn,
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
		logger.Info(dn)
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

func ( t *DeptInfo ) ChangeDn ( oldDn string, newDn string ) interface{} {
	req := ldap.NewModifyDNRequest(oldDn, newDn,true , "" )
	err := ldapConn.ModifyDN(req)
	if err != nil {
		fmt.Println(err)
		return t
	}
	t.DN = newDn
	return t
}

func InitDmap () {
	d := new (DeptsMap )
	d.Init()
	Dmap.SimpleMap = d.SimpleMap
	Dmap.MultipleMap = d.MultipleMap
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
		//logger.Info("当前dn",v.DN)
		//logger.Info("上级dn",v.ParentDN)
		d := DeptInfo{}
		d.AddToLdap(v.Name, v.DN)
		u := new (UserList)
		u.Get(v.ID)
		SyncAllUser(u.Userlist,v.DN)
	}
}