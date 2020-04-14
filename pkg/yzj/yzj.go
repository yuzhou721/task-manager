package yzj

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"task/conf"
	"time"
)

type Yzj struct {
	token  string
	AppID  string
	Secret string
	Scope  string
}

const (
	//YzjScopeApp 授权级别app
	YzjScopeApp = "app"
	//YzjScopeGroup 授权级别Group
	YzjScopeGroup = "resGroupSecret"
	//YzjScopeTeam 授权级别Team
	YzjScopeTeam = "team"
	//ContentTypeJSON json
	ContentTypeJSON = "application/json"
	//ContentTypeForm www form
	ContentTypeForm = "application/x-www-form-urlencoded"
)

//Person 人员信息
type Person struct {
	OpenID     string `json:"openId"`     //人员的openid
	Name       string `json:"name"`       //姓名
	PhotoURL   string `json:"photoUrl"`   //头像URL
	Phone      string `json:"phone"`      //手机号码（未开放）
	Email      string `json:"email"`      //邮箱（未开放）
	Department string `json:"department"` //部门
	OrgID      string `json:"orgId"`      //用户所在部门ID
	JobTitle   string `json:"jobTitle"`   //职位
	Gender     string `json:"gender"`     //性别, 0: 不确定; 1: 男; 2: 女
	IsAdmin    string `json:"isAdmin"`    //是否管理员， 1: 是； 0:不是，拥有工作圈的最高权限（与创建者一致）
	Status     string `json:"status"`     //是否在职， 1：在职；0: 离职
	JobNo      string `json:"jobNo"`      //工号
}

//Org 部门信息
type Org struct {
	Name       string   `json:"name"`       //部门名称
	OrgID      string   `json:"orgId"`      //部门ID
	ParentID   string   `json:"parentId"`   //上级部门ID
	InChargers []Person `json:"inChargers"` //负责人
}

//AppTokenRequest APP级别密钥获取请求
type appTokenRequest struct {
	AppID     string `json:"appId"`
	Secret    string `json:"secret"`
	Timestamp int    `json:"timestamp"`
	Scope     string `json:"scope"`
}

type getPersonRequest struct {
	OpenID string `json:"openId"`
	EID    string `json:"eid"`
}

type tokenData struct {
	AccessToken  string `json:"accessToken"`
	ExpireIn     int    `json:"expireIn"`
	RefreshToken string `json:"refreshToken"`
}

type yzjResponse struct {
	Error     string `json:"error"`
	ErrorCode int    `json:"errorCode"`
	Success   bool   `json:"success"`
}

type getPersonResponse struct {
	yzjResponse
	Data []Person `json:"data"`
}

type getOrgResponse struct {
	yzjResponse
	Data Org `json:"data"`
}

type getTokenResponse struct {
	yzjResponse
	Data tokenData `json:"data"`
}

//GetPerson 根据openID获取人员信息
func (y *Yzj) GetPerson(openID string) (*Person, error) {
	err := y.getToken()
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("%v/gateway/opendata-control/data/getperson?accessToken=%v", conf.Config.Yzj.YZJServer, y.token)
	client := &http.Client{}
	data := &getPersonRequest{
		EID:    conf.Config.Yzj.EID,
		OpenID: openID,
	}
	formData := make(url.Values)
	formData.Add("eid", data.EID)
	formData.Add("openId", data.OpenID)

	response, err := client.PostForm(u, formData)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {

	}

	body, _ := ioutil.ReadAll(response.Body)
	// bodystr := string(body)
	var res getPersonResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}
	if res.Success == false {
		return nil, errors.New(res.Error)
	}
	persons := res.Data

	return &persons[0], nil
}

func (y *Yzj) getToken() error {
	var (
		err error
	)
	client := &http.Client{}
	t, err := strconv.Atoi(fmt.Sprintf("%v", time.Now().UnixNano()/1e6))
	if err != nil {
		return err
	}

	data := &appTokenRequest{
		AppID:     y.AppID,
		Secret:    y.Secret,
		Timestamp: t,
		Scope:     y.Scope,
	}
	j, err := json.Marshal(data)
	if err != nil {
		log.Printf("error :%v", err)
	}
	log.Printf("%v", string(j))
	url := fmt.Sprintf("%v/gateway/oauth2/token/getAccessToken", conf.Config.Yzj.YZJServer)
	res, err := client.Post(url, "application/json", bytes.NewBuffer(j))
	if err != nil {
		return err
	}
	if res.StatusCode == 200 {
		body, _ := ioutil.ReadAll(res.Body)
		// bodystr := string(body)
		var m getTokenResponse
		err := json.Unmarshal(body, &m)
		if err != nil {
			return err
		}
		if m.Success == false {
			return errors.New(m.Error)
		}
		d := m.Data
		y.token = d.AccessToken
	}
	return nil
}

//GetOrg 根据orgID和EID获取部门信息
func (y *Yzj) GetOrg(orgID string) (*Org, error) {
	var (
		org Org
		err error
	)
	err = y.getToken()
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("%v/gateway/opendata-control/data/getorg?accessToken=%v", conf.Config.Yzj.YZJServer, y.token)

	client := &http.Client{}
	formData := make(url.Values)
	formData.Add("orgId", orgID)
	formData.Add("eid", conf.Config.Yzj.EID)

	response, err := client.PostForm(u, formData)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {

	}

	body, _ := ioutil.ReadAll(response.Body)
	// bodystr := string(body)
	var res getOrgResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}
	if res.Success == false {
		return nil, errors.New(res.Error)
	}
	org = res.Data

	return &org, nil
}

//SendTodo 发送待办
func (y *Yzj) SendTodo() (err error) {
	err = y.getToken()
	if err != nil {
		return
	}
	return
}
