package v1

import (
	"net/http"
	"strconv"
	"task/app/models"
	"task/pkg/utils"

	"github.com/gin-gonic/gin"
)

//SaveTask 保存任务
func SaveTask(c *gin.Context) {
	var t models.Task
	err := c.BindJSON(&t)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "error ",
		})
	}
	c.Status(http.StatusCreated)
}

//DelTask 删除任务
func DelTask(c *gin.Context) {
	var t models.Task
	id := c.Param("id")
	iid, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "error id",
		})
	}
	t.ID = uint(iid)
	t.Delete()
	c.Status(http.StatusNoContent)
}

//UpdateTask 修改任务
func UpdateTask(c *gin.Context) {
	var t models.Task
	err := c.BindJSON(&t)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "error id",
		})
	}
	if err = t.Update(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
	}
	c.Status(http.StatusOK)
}

//GetTasks 获取task分页数据
func GetTasks(c *gin.Context) {
	var t models.Task
	role := c.Query("role")
	ri, err := strconv.Atoi(role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
	}
	search := c.Query("search")
	page, pageSize, err := utils.GetPage(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
	}
	err = c.BindJSON(&t)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
	}
	ts, count, err := t.FindList(ri, search, page, pageSize)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  ts,
		"count": count,
	})
}

//GetTask 根据ID获取单条数据
func GetTask(c *gin.Context) {
	var t models.Task
	id := c.Param("id")
	t.FindByID(id)
}

//SaveTasks 保存任务列表
func SaveTasks(c *gin.Context) {
	var tasks []models.Task
	err := c.BindJSON(&tasks)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
	}
	(&models.Task{}).SaveList(&tasks)
}
