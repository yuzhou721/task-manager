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
