package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

//GetPageOffset 根据页数获取数据库查询偏移量
func GetPageOffset(page, pageSize int) int {
	result := 0
	if page > 0 {
		result = (page - 1) * pageSize
	}
	return result
}

//GetPage 参数里获取page信息
func GetPage(c *gin.Context) (page, pageSize int, err error) {
	page, err = strconv.Atoi(c.Query("page"))
	if err != nil {
		return 0, 0, err
	}
	pageSize, err = strconv.Atoi(c.Query("pageSize"))
	if err != nil {
		return 0, 0, err
	}
	return
}
