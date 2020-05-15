package v1

import (
	"net/http"
	"strconv"
	"task/app/models"
	"task/pkg/utils"

	"github.com/gin-gonic/gin"
)

//Response 返回值统一封装
type Response struct {
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
	Success bool        `json:"success"`
}

//SaveTask 保存任务
func SaveTask(c *gin.Context) {
	var t models.Task
	err := c.BindJSON(&t)
	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Msg:     err.Error(),
			Success: false,
		})
	}
	err = t.Save()
	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Msg:     err.Error(),
			Success: false,
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
	var ri = 0
	role := c.Query("role")
	ri, err := strconv.Atoi(role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
	}
	openID := c.Query("openId")
	orgID := c.Query("orgId")
	search := c.Query("search")
	page, pageSize, err := utils.GetPage(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}

	status := c.Query("status")
	statusI, err := strconv.Atoi(status)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}
	ts, count, err := t.FindList(ri, openID, orgID, search, statusI, page, pageSize)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
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
	err := t.FindByID(id)
	if err != nil {
		c.JSON(http.StatusOK, &Response{
			Success: false,
			Msg:     err.Error(),
		})
	}
	c.JSON(http.StatusOK, &Response{
		Data:    t,
		Success: true,
	})
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
	(&models.Task{}).SaveOrUpdateList(&tasks)
	c.Status(http.StatusCreated)
}

//SaveTaskMasterAndSlave 创建页面保存主从表
func SaveTaskMasterAndSlave(c *gin.Context) {
	var tasks []models.Task
	err := c.BindJSON(&tasks)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			&Response{
				Msg:     err.Error(),
				Success: false,
			})
		return
	}
	//如果值录入了主任务 就只保存一条
	if len(tasks) == 1 {
		tasks[0].Save()
	} else {
		//如果录入子任务 就按照主从保存
		(&models.Task{}).SaveMasterAndSlave(tasks)
	}
	c.Status(http.StatusCreated)
}

//CountTask 根据id查询任务数量
func CountTask(c *gin.Context) {
	var t models.Task
	id := c.Param("id")
	percent, err := t.CountTaskPercent(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Msg: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, &Response{
		Success: true,
		Data: &gin.H{
			"percent": percent,
		},
	})

}

// DeleteAttach 删除附件
func DeleteAttach(c *gin.Context) {
	var a models.Attach
	id := c.Param("id")
	err := a.Delete(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Msg: err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)

}
