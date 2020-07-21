package v1

import (
	"log"
	"net/http"
	"task/conf"
	"task/pkg/yzj"

	"github.com/gin-gonic/gin"
)

//GetPerson 通过接口获取用户信息
func GetPerson(c *gin.Context) {
	openID := c.Param("id")
	y := &yzj.Yzj{
		AppID:  conf.Config.Yzj.AppID,
		Secret: conf.Config.Yzj.Secret,
		Scope:  yzj.YzjScopeApp,
	}
	person, err := y.GetPerson(openID)
	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Success: false,
			Msg:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, &Response{
		Success: true,
		Data:    person,
	})
}

//GetOrg 获取部门数据
func GetOrg(c *gin.Context) {
	orgID := c.Param("id")
	y := &yzj.Yzj{
		AppID:  conf.Config.Yzj.AppID,
		Secret: conf.Config.Yzj.Secret,
		Scope:  yzj.YzjScopeApp,
	}
	org, err := y.GetOrg(orgID)
	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Success: false,
			Msg:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, &Response{
		Success: true,
		Data:    org,
	})

}

// GetAllOrgs 获取所有部门信息
func GetAllOrgs(c *gin.Context) {
	y := &yzj.Yzj{
		AppID:  conf.Config.Yzj.AppID,
		Secret: conf.Config.Yzj.Secret,
		Scope:  yzj.YzjScopeApp,
	}
	orgs, err := y.GetAllOrgs()
	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Success: false,
			Msg:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, &Response{
		Success: true,
		Data:    orgs,
	})
}

// GetOrgPersons 获取部门人员
func GetOrgPersons(c *gin.Context) {
	orgID := c.Param("id")
	y := &yzj.Yzj{
		AppID:  conf.Config.Yzj.AppID,
		Secret: conf.Config.Yzj.Secret,
		Scope:  yzj.YzjScopeApp,
	}
	_, persons, err := y.GetOrgPersons(orgID)
	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Success: false,
			Msg:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, &Response{
		Success: true,
		Data:    persons,
	})
}

// AcquireContext 根据ticket获取上下文
func AcquireContext(c *gin.Context) {
	ticket := c.Param("ticket")
	y := &yzj.Yzj{
		AppID:  conf.Config.Yzj.AppID,
		Secret: conf.Config.Yzj.Secret,
		Scope:  yzj.YzjScopeApp,
	}
	ct, err := y.AcquireContext(ticket)
	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Success: false,
			Msg:     err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &Response{
		Success: true,
		Data:    ct,
	})
}

func IsManager(c *gin.Context) {
	openId := c.Param("openId")
	y := &yzj.Yzj{
		EID:    conf.Config.Yzj.EID,
		Secret: conf.Config.Yzj.OnlyReadSecret,
		Scope:  yzj.YzjScopeGroup,
	}

	roleId := conf.Config.Yzj.TaskRoleId
	arrOpenId, err := y.GetTaskManager(roleId)
	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Success: false,
			Msg:     err.Error(),
		})
		return
	}

	flag := "0"
	msg := ""

	index := -1
	for i := 0; i < len(arrOpenId); i++ {
		if arrOpenId[i] == openId {
			index = i
			break
		}
	}
	if index == -1 {
		log.Printf("用户%v不是管理员", openId)
		msg = "该用户不是管理员"
	} else {
		log.Printf("用户%v是管理员", openId)
		flag = "1"
		msg = "该用户是管理员"
	}
	c.JSON(http.StatusOK, &Response{
		Success: true,
		Data:    flag,
		Msg:     msg,
	})

	return
}
