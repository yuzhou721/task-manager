package test

import (
	"task/conf"
	"task/pkg/yzj"
	"testing"
)

func TestGetPerson(t *testing.T) {
	y := &yzj.Yzj{
		AppID:  conf.Config.Yzj.AppID,
		Secret: conf.Config.Yzj.Secret,
		Scope:  yzj.YzjScopeApp,
	}
	person, err := y.GetPerson("5df05b6bd08e4390b7f7306f")
	if err != nil {
		t.Error("error:", err)
	}
	t.Log("person:", person)
}

func TestGetOrg(t *testing.T) {
	y := &yzj.Yzj{
		AppID:  conf.Config.Yzj.AppID,
		Secret: conf.Config.Yzj.Secret,
		Scope:  yzj.YzjScopeApp,
	}
	org, err := y.GetOrg("273413ce-1a56-11ea-9751-0050569293b2")
	if err != nil {
		t.Error("error:", err)
	}
	t.Log("org:", org)
}

func TestSendNotify(t *testing.T) {
	y := &yzj.Yzj{
		EID:       conf.Config.Yzj.EID,
		PubSecret: conf.Config.Yzj.PubSecret,
		PubID:     conf.Config.Yzj.PubID,
	}
	err := y.GenerateNotify("测试信息", "http://www.baidu.com", []string{"5de5fd59d08e886badfeb8d8"})
	if err != nil {
		t.Errorf("err send Notify:%v", err)
	}
}

func TestSendTodo(t *testing.T) {
	y := &yzj.Yzj{
		AppID:  conf.Config.Yzj.AppID,
		Secret: conf.Config.Yzj.Secret,
		Scope:  yzj.YzjScopeApp,
	}
	err := y.GenerateTODO("2555", []string{"5de5fd59d08e886badfeb8d8"}, "测试title", "测试content", "测试itemTitle", "http://www.baidu.com", "http://yunzhi.cyats.com/docrest/file/downloadfile/5dfb1ae34f1c47119704b1f7")
	if err != nil {
		t.Errorf("err send Todo:%v", err.Error())
	}
}

func TestOperationTodo(t *testing.T) {
	y := &yzj.Yzj{
		AppID:  conf.Config.Yzj.AppID,
		Secret: conf.Config.Yzj.Secret,
		Scope:  yzj.YzjScopeApp,
	}
	err := y.OprateTodo("2555", []string{"5de5fd59d08e886badfeb8d8"}, 1, 1, 1)
	if err != nil {
		t.Errorf("err change todo:%v", err.Error())
	}
}

func TestGetOrgPersons(t *testing.T) {
	y := &yzj.Yzj{
		AppID:  conf.Config.Yzj.AppID,
		Secret: conf.Config.Yzj.Secret,
		Scope:  yzj.YzjScopeApp,
	}
	a, b, err := y.GetOrgPersons("273413ce-1a56-11ea-9751-0050569293b2")
	if err != nil {
		t.Errorf("err getOrgPersons :%v", err)
	}
	t.Logf("inChargers:%v , member : %v", a, b)
}

func TestGetAllOrgs(t *testing.T) {
	y := &yzj.Yzj{
		AppID:  conf.Config.Yzj.AppID,
		Secret: conf.Config.Yzj.Secret,
		Scope:  yzj.YzjScopeApp,
	}
	a, err := y.GetAllOrgs()
	if err != nil {
		t.Errorf("err getAllOrgs :%v", err)
	}
	t.Logf("allorgs:%v", a)
}

func TestAcquireContext(t *testing.T) {
	y := &yzj.Yzj{
		AppID:  conf.Config.Yzj.AppID,
		Secret: conf.Config.Yzj.Secret,
		Scope:  yzj.YzjScopeApp,
	}
	c, e := y.AcquireContext("APPURLWITHTICKET6be84cd6011995696306406a4835dcf4")
	if e != nil {
		t.Errorf("err AcquireContext:%v", e)
	}
	t.Logf("Context:%v", c)
}
