package yzj

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"task/conf"
	"task/pkg/utils"
	"time"
)

type Yzj struct {
	token        string
	refreshToken string
	expireIn     int
	AppID        string
	Secret       string
	Scope        string
	EID          string
	PubID        string
	PubSecret    string
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

//resGroupSecret 级别密钥获取请求
type resGroupSecretTokenRequest struct {
	Eid       string `json:"eid"`
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

type pubResponse struct {
	PubID       string `json:"pubId"`
	SourceMsgID string `json:"sourceMsgId"`
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
	// log.Printf("%v", string(j))
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

func (y *Yzj) getTokenOfResGroupSecret() error {
	var (
		err error
	)
	client := &http.Client{}
	t, err := strconv.Atoi(fmt.Sprintf("%v", time.Now().UnixNano()/1e6))
	if err != nil {
		return err
	}

	data := &resGroupSecretTokenRequest{
		Eid:       y.EID,
		Secret:    y.Secret,
		Timestamp: t,
		Scope:     y.Scope,
	}
	j, err := json.Marshal(data)
	if err != nil {
		log.Printf("error :%v", err)
	}
	// log.Printf("%v", string(j))
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

type todoRequest struct {
	URL       string      `json:"url"`       //	是 	点击待办跳转的URL
	SourceID  string      `json:"sourceId"`  //	是 	生成的待办所关联的第三方服务业务记录的ID，待办的批次号，对应待办处理中的sourceitemid
	Content   string      `json:"content"`   //	否 	待办内容
	Title     string      `json:"title"`     //	是 	来自字段的内容显示
	Itemtitle string      `json:"itemtitle"` //	否 	待办项标题内容显示,选填，如不填，则默认为title值
	HeadImg   string      `json:"headImg"`   //	是 	待办在客户端显示的图URL
	AppID     string      `json:"appId"`     //	是 	生成的待办所关联的第三方服务类型ID,appId和sourceId组合在一起标识云之家唯一的一条待办
	SenderID  string      `json:"senderId"`  //	否 	待办的发送人的openId
	Params    []todoParam `json:"params"`
}

type todoParam struct {
	OpenID string     `json:"openId"` //	是 	待办接受人ID，可填多人
	Status todoStatus `json:"status"`
}

type todoStatus struct {
	DO   int `json:"DO"`   //	否 	目标处理状态，0表示未办，1表示已办，默认为0
	READ int `json:"READ"` //	否 	目标读状态，0表示未读，1表示已读，默认为0
}

type todoResponse struct {
	Success   string      `json:"success"`
	ErrorCode int         `json:"errorCode"`
	Error     string      `json:"error"`
	Data      interface{} `json:"data"`
}

//GenerateTODO 发送待办
func (y *Yzj) GenerateTODO(sourceID string, openIDs []string, title, content, itemTitle, URL, headImgURL string) (err error) {
	err = y.getToken()
	if err != nil {
		return
	}
	params := []todoParam{}
	status := todoStatus{0, 0}
	for _, v := range openIDs {
		p := todoParam{v, status}
		params = append(params, p)
	}

	request := &todoRequest{
		URL:       URL,
		SourceID:  sourceID,
		Content:   content,
		Title:     title,
		Itemtitle: itemTitle,
		HeadImg:   headImgURL,
		AppID:     y.AppID,
		Params:    params,
	}
	j, err := json.Marshal(request)
	if err != nil {
		log.Printf("error :%v", err)
		return
	}
	u := fmt.Sprintf("%v/gateway/newtodo/open/generatetodo.json?accessToken=%v", conf.Config.Yzj.YZJServer, y.token)
	client := &http.Client{}
	response, err := client.Post(u, "application/json", bytes.NewBuffer(j))

	log.Printf("发送待办的url:" + u)
	log.Printf("参数:")
	fmt.Println(request)
	if err != nil {
		return err
	}

	var res todoResponse
	err = marshal(response, &res)
	if err != nil {
		return
	}

	log.Printf("发送待办return:")
	fmt.Println(res)
	if res.Success == "false" {
		return errors.New(res.Error)
	}
	return
}

type oprateTodoRequest struct {
	Sourcetype   string     `json:"sourcetype"`   //	是 	应用ID，即appId
	Sourceitemid string     `json:"sourceitemid"` //	是 	即发送待办的sourceId,生成的待办所关联的第三方服务业务记录的ID，是待办的批次号
	Openids      []string   `json:"openids"`      // 	否 	可填多人，不填则更改sourceitemid下所有人员的待办状态
	ActionType   actionType `json:"actiontype"`
	Sync         bool       `json:"sync"`
}

type actionType struct {
	Deal   int `json:"deal"`   //	否 	目标处理状态，0表示未办，1表示已办，默认为0
	Read   int `json:"read"`   //	否 	目标读状态，0表示未读，1表示已读，默认为0
	Delete int `json:"delete"` //	否 	目标删除状态，0表示未删除，1表示已删除
}

//OprateTodo 修改TODO状态
func (y *Yzj) OprateTodo(sourceID string, openIDs []string, deal int, read, delete int) (err error) {
	err = y.getToken()
	if err != nil {
		return
	}
	at := actionType{deal, read, delete}
	request := &oprateTodoRequest{
		Sourceitemid: sourceID,
		Sourcetype:   y.AppID,
		Openids:      openIDs,
		ActionType:   at,
		Sync:         true,
	}
	j, err := json.Marshal(request)
	if err != nil {
		log.Printf("error :%v", err)
		return
	}

	log.Printf("清除待办......")
	log.Printf("参数：")
	fmt.Println(request)
	url := fmt.Sprintf("%v/gateway/newtodo/open/action.json?accessToken=%v", conf.Config.Yzj.YZJServer, y.token)
	client := &http.Client{}
	response, err := client.Post(url, "application/json", bytes.NewBuffer(j))
	if err != nil {
		return
	}
	log.Printf("url：" + url)
	var res todoResponse
	err = marshal(response, &res)
	if err != nil {
		return
	}

	if res.Success == "false" {
		return errors.New(res.Error)
	}
	return
}

func marshal(res *http.Response, responseData interface{}) (err error) {
	body, _ := ioutil.ReadAll(res.Body)
	log.Printf("body = %v", string(body))
	if res.StatusCode != 200 {
		return errors.New(string(body))
	}
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		return
	}
	return
}

type pubRequest struct {
	From from `json:"from"` //"from":"发送方信息，格式为JSON对象",
	To   []to `json:"to"`   //"to":"接收方信息，格式为包含一至多个接收方信息JSON对象的JSON数组",
	Type int  `json:"type"` //"type":"消息类型，格式为整型",(取值 2：单文本,5：文本链接,6：图文链接)
	Msg  msg  `json:"msg"`  //"msg":"发布到讯通的消息内容，格式为JSON对象"
}

type from struct {
	No       string `json:"no"`       //"no":"发送方企业的企业注册号(eid)，格式为字符串",
	Pub      string `json:"pub"`      //"pub":"发送使用的公共号ID，格式为字符串",
	Time     string `json:"time"`     //"time":"发送时间，为'currentTimeMillis()以毫秒为单位的当前时间'的字符串或数字",
	Nonce    int    `json:"nonce"`    //"nonce":"随机数，格式为字符串或数字",
	PubToken string `json:"pubtoken"` //"pubtoken":"公共号加密串，格式为字符串。"
}

type to struct {
	No   string   `json:"no"`   //"no":"接收方企业的企业注册号(eID)，格式为字符串",
	User []string `json:"user"` //"user":"接收方的用户ID，格式为包含OPENID的JSON数组"
}

type msg struct {
	Text  string `json:"text"`  //	"text":"文本消息内容，String",
	URL   string `json:"url"`   //	"url":"文本链接地址，格式为经过URLENCODE编码的字符串",
	APPID string `json:"appid"` //	"appid": "如果打开的链接是轻应用,必须传入轻应用号讯通才能传入参数ticket,参考<轻应用框架>开发",
	Todo  int    `json:"todo"`  //	"todo":"int，必填,暂时只能为0，表示推送原公共号消息",
}

//GenerateNotify 构建通知
func (y *Yzj) GenerateNotify(text, url string, openIDs []string) (err error) {
	msg := msg{
		URL:   url,
		Text:  text,
		APPID: conf.Config.Yzj.AppID,
		Todo:  0,
	}
	t := fmt.Sprintf("%v", time.Now().UnixNano()/1e6)
	eid := y.EID
	pubID := y.PubID
	pubSecret := y.PubSecret
	nonce := rand.Int()

	pubToken := utils.Sha([]string{eid, pubID, pubSecret, strconv.Itoa(nonce), t})
	f := from{
		No:       eid,
		Pub:      pubID,
		Time:     t,
		Nonce:    nonce,
		PubToken: pubToken,
	}

	to1 := to{
		No:   eid,
		User: openIDs,
	}

	request := &pubRequest{
		From: f,
		To:   []to{to1},
		Type: 5,
		Msg:  msg,
	}

	j, err := json.Marshal(request)
	if err != nil {
		log.Printf("error :%v", err)
	}
	log.Println("发送公众号消息")
	log.Println("参数为:")
	fmt.Println(request)
	u := fmt.Sprintf("%v/pubacc/pubsend", conf.Config.Yzj.YZJServer)
	log.Println("请求url:" + u)
	client := &http.Client{}
	response, err := client.Post(u, "application/json", bytes.NewBuffer(j))

	if err != nil {
		return
	}
	var res pubResponse
	err = marshal(response, res)
	if err != nil {
		return
	}
	return
}

type getOrgPersonsData struct {
	InChargers []Person `json:"inChargers"`
	Members    []Person `json:"members"`
}

type getOrgPersonsResponse struct {
	yzjResponse
	Data getOrgPersonsData `json:"data"`
}

//GetOrgPersons 根据orgID获取部门人员信息
func (y *Yzj) GetOrgPersons(orgID string) (inChargers, Members []Person, err error) {
	err = y.getToken()
	if err != nil {
		return
	}
	formData := make(url.Values)
	formData.Add("orgId", orgID)
	formData.Add("eid", conf.Config.Yzj.EID)
	url := fmt.Sprintf("%v/gateway/opendata-control/data/getorgpersons?accessToken=%v", conf.Config.Yzj.YZJServer, y.token)
	client := &http.Client{}
	response, err := client.PostForm(url, formData)
	var responseData getOrgPersonsResponse
	err = marshal(response, &responseData)
	if err != nil {
		return
	}
	if responseData.Success == false {
		err = errors.New(responseData.Error)
		return
	}
	inChargers = responseData.Data.InChargers
	Members = responseData.Data.Members

	return
}

type getAllOrgsResponse struct {
	yzjResponse
	Data []Org `json:"data"`
}

// GetAllOrgs 获取所有部门信息
func (y *Yzj) GetAllOrgs() (orgs []Org, err error) {
	err = y.getToken()
	if err != nil {
		return
	}
	formData := make(url.Values)
	formData.Add("eid", conf.Config.Yzj.EID)
	url := fmt.Sprintf("%v/gateway/opendata-control/data/getallorgs?accessToken=%v", conf.Config.Yzj.YZJServer, y.token)
	client := &http.Client{}
	response, err := client.PostForm(url, formData)
	if err != nil {
		return
	}
	var responseData getAllOrgsResponse
	err = marshal(response, &responseData)
	if err != nil {
		return
	}
	if responseData.Success == false {
		err = errors.New(responseData.Error)
		return
	}
	orgs = responseData.Data
	return
}

// YzjContext 云之家上下文
type YzjContext struct {
	AppID     string `json:"appid"`
	XTID      string `json:"xtid"`
	OID       string `json:"oid"`
	EID       string `json:"eid"`
	UserName  string `json:"username"`
	UserID    string `json:"userid"`
	UID       string `json:"uid"`
	TID       string `json:"tid"`
	JobNo     string `json:"jobNo"`
	NetworkID string `json:"networkid"`
	DeviceID  string `json:"devceId"`
	OpenID    string `json:"openid"`
	OrgId     string `json:"orgId"`
	Mobile    string `json:"mobile"`
}

type acquireContextResponse struct {
	yzjResponse
	Data YzjContext `json:"data"`
}

type acquireContextRequest struct {
	AppID  string `json:"appid"`
	Ticket string `json:"ticket"`
}

// AcquireContext 获取用户上下文信息
func (y *Yzj) AcquireContext(ticket string) (context YzjContext, err error) {
	err = y.getToken()
	if err != nil {
		return
	}
	request := &acquireContextRequest{
		AppID:  conf.Config.Yzj.AppID,
		Ticket: ticket,
	}
	j, err := json.Marshal(request)
	if err != nil {
		log.Printf("error :%v", err)
	}

	u := fmt.Sprintf("%v/gateway/ticket/user/acquirecontext?accessToken=%v", conf.Config.Yzj.YZJServer, y.token)
	client := &http.Client{}
	response, err := client.Post(u, "application/json", bytes.NewBuffer(j))
	if err != nil {
		return
	}
	var responseData acquireContextResponse
	err = marshal(response, &responseData)
	if err != nil {
		return
	}
	if responseData.Success == false {
		err = errors.New(responseData.Error)
		return
	}
	context = responseData.Data
	return

}

type role struct {
	RoleId string `json:"roleId"`
}

type roleRequest struct {
	Nonce string `json:"nonce"`
	EId   string `json:"eid"`
	Data  string `json:"data"`
}

type YzjRole struct {
	OrgIds string `json:"orgIds"`
	OpenId string `json:"openId"`
}

type roleResponse struct {
	yzjResponse
	Data []YzjRole `json:"data"`
}

//根据角色Id获取任务中心人员列表
func (y *Yzj) GetTaskManager(roleId string) (arrOpenId []string, err error) {
	log.Printf("根据openId获取任务中心管理员")
	r := role{RoleId: roleId}
	ro, err := json.Marshal(r)

	if err != nil {
		return
	}

	err = y.getTokenOfResGroupSecret()
	if err != nil {
		return
	}
	if err != nil {
		log.Printf("error :%v", err)
	}

	u := fmt.Sprintf("%v/gateway/openimport/open/roletag/getPersonsByRoleAndPage?accessToken=%v", conf.Config.Yzj.YZJServer, y.token)

	DataUrlVal := url.Values{}
	client := &http.Client{}

	DataUrlVal.Add("nonce", string(time.Now().Unix()))
	DataUrlVal.Add("eid", y.EID)
	DataUrlVal.Add("data", string(ro))

	req, err := http.NewRequest("POST", u, strings.NewReader(DataUrlVal.Encode()))

	log.Printf("请求url%v", u)
	if err != nil {
		return
	}
	//伪装头部
	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Accept-Language", "zh-CN,zh;q=0.8,en-US;q=0.6,en;q=0.4")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-Length", "25")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Cookie", "user_trace_token=20170425200852-dfbddc2c21fd492caac33936c08aef7e; LGUID=20170425200852-f2e56fe3-29af-11e7-b359-5254005c3644; showExpriedIndex=1; showExpriedCompanyHome=1; showExpriedMyPublish=1; hasDeliver=22; index_location_city=%E5%85%A8%E5%9B%BD; JSESSIONID=CEB4F9FAD55FDA93B8B43DC64F6D3DB8; TG-TRACK-CODE=search_code; SEARCH_ID=b642e683bb424e7f8622b0c6a17ffeeb; Hm_lvt_4233e74dff0ae5bd0a3d81c6ccf756e6=1493122129,1493380366; Hm_lpvt_4233e74dff0ae5bd0a3d81c6ccf756e6=1493383810; _ga=GA1.2.1167865619.1493122129; LGSID=20170428195247-32c086bf-2c09-11e7-871f-525400f775ce; LGRID=20170428205011-376bf3ce-2c11-11e7-8724-525400f775ce; _putrc=AFBE3C2EAEBB8730")
	req.Header.Add("X-Anit-Forge-Code", "0")
	req.Header.Add("X-Anit-Forge-Token", "None")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/61.0.3163.100 Safari/537.36")
	req.Header.Add("X-Requested-With", "XMLHttpRequest")

	//提交请求
	resp, err := client.Do(req)

	var responseData roleResponse

	defer resp.Body.Close()
	if err != nil {
		return
	}
	//读取返回值
	response, err := ioutil.ReadAll(resp.Body)
	log.Printf("返回结果：%v", string(response))

	//角色下没有人员时，返回data为空字符串，有人员信息时，返回值为数组，没做差异化解析，可能会报错。
	err = json.Unmarshal(response, &responseData)
	if err != nil {
		return
	}

	//遍历返回值
	for _, val := range responseData.Data {
		arrOpenId = append(arrOpenId, val.OpenId)
	}
	return
}
