package routers

import (
	"net/http"
	v1 "task/app/controllers/v1"
	"task/conf"

	"github.com/gin-gonic/gin"
)

//InitRouter 初始化router
func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())

	r.Use(gin.Recovery())

	gin.SetMode(conf.Config.RunMode)

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong!")
	})

	apiv1 := r.Group("/api/v1")

	{
		//任务相关接口
		apiv1.POST("/tasks/list", v1.SaveTasks)
		apiv1.POST("/tasks", v1.SaveTask)
		apiv1.POST("/tasks/master", v1.SaveTaskMasterAndSlave)
		apiv1.DELETE("/tasks/:id", v1.DelTask)
		apiv1.PUT("/tasks/:id", v1.UpdateTask)
		apiv1.GET("/tasks", v1.GetTasks)
		apiv1.GET("/tasks/:id", v1.GetTask)
		apiv1.GET("/count/tasks/:id", v1.CountTask)

		//云之家接口
		apiv1.GET("/person/:id", v1.GetPerson)
		apiv1.GET("/org/:id", v1.GetOrg)
		apiv1.GET("/persons/orgs/:id", v1.GetOrgPersons)
		apiv1.GET("/orgs", v1.GetAllOrgs)
		apiv1.GET("/context", v1.AcquireContext)
	}

	return r
}
