package v1

import (
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
	}

	c.JSON(http.StatusOK, &Response{
		Success: true,
		Data:    ct,
	})
}
