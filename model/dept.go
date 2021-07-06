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
	Multiple map[int]map[string]string
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
	m := make(map[int]map[string]string)
	for _ , v := range d.Department {
		m[v.ID] = make(map[string]string)
		m[v.ID]["name"] = v.Name
		m[v.ID]["parentid"] = strconv.Itoa(v.ParentID)
	}
	for _ , v := range d.Department {
		id := v.ID
		name := v.Name
		dn := "ou=" + name
		var pdn string
		for i := 1;i <= 5; i++ {
			parentId := m[id]["parentid"]
			if parentId != "" && parentId != "0" && parentId != "1" {
				pID, _ := strconv.Atoi(parentId)
				if ( i == 1 ) {
					pdn = pdn + "ou=" + m[pID]["name"]
				} else {
					pdn = pdn + ",ou=" + m[pID]["name"]
				}
				id, _ = strconv.Atoi(parentId)
			}
		}
		if v.ParentID > 1 {
			pdn = pdn + "," + ldapConfig.BaseDn
		} else {
			pdn = ldapConfig.BaseDn
		}
		m[v.ID]["dn"] = dn
		m[v.ID]["pdn"] = pdn
	}
	t.Multiple = m
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
		d := DeptInfo{}
		d.ID = v.ID
		d.Name = v.Name
		d.DN = Dmap.Multiple[v.ID]["dn"]
		d.ParentID = v.ParentID
		d.ParentDN = Dmap.Multiple[v.ID]["pdn"]
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

func ( t *DeptInfo ) ModifyDn ( dn string, rdn string,newSup string  ) interface{} {
	req := ldap.NewModifyDNRequest(dn, rdn,true , newSup )
	err := ldapConn.ModifyDN(req)
	if err != nil {
		fmt.Println(err)
		return t
	}
	t.DN = rdn + "," + newSup
	return t
}

func InitDmap () {
	d := new (DeptsMap )
	d.Init()
	Dmap.Multiple = d.Multiple
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
		d.AddToLdap(v.Name, v.DN + "," + v.ParentDN )
		u := new (UserList)
		u.Get(v.ID)
		SyncAllUser(u.Userlist,v.DN + "," + v.ParentDN)
	}
}