package v1

import (
	"net/http"
	"task/app/models"

	"github.com/gin-gonic/gin"
)

func IsAdmin(c *gin.Context) {

	phone := c.Param("phone")
	var m models.Manager
	err := m.IsAdmin(phone)

	//如果没从数据库中查询到数据，则报错
	if err != nil {
		if err != nil {
			c.JSON(http.StatusOK, Response{
				Success: true,
				Data: gin.H{
					"isManager": "0",
				},
			})
		}
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data: gin.H{
			"isManager": "1",
		},
	})

}
