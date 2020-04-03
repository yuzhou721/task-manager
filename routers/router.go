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
		apiv1.POST("/tasks", v1.SaveTask)
		apiv1.DELETE("/tasks/:id", v1.DelTask)
		apiv1.PUT("/tasks", v1.UpdateTask)
		apiv1.GET("/tasks", v1.GetTask)
		apiv1.GET("/tasks/list", v1.GetTasks)
	}

	return r
}
