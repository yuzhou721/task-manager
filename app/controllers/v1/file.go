package v1

import (
	"net/http"
	"path/filepath"
	"task/app/models"
	"task/pkg/file"
	"task/pkg/upload"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

//FileUpload 上传文件
func FileUpload(c *gin.Context) {
	f, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Success: false,
			Msg:     err.Error(),
		})
		return
	}

	fileName := header.Filename
	ext := filepath.Ext(fileName)
	uuid, err := uuid.NewV4()
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Msg:     "genarete uuid error",
		})
	}

	uuidS := uuid.String()
	timeS := time.Now().Format("2006/01/02")
	savePath := filepath.Join(timeS)
	fullPath := filepath.Join(upload.GetFileFullPath(), savePath)
	if file.CheckNotExist(fullPath) {
		file.MkDir(fullPath)
	}

	err = c.SaveUploadedFile(header, filepath.Join(fullPath, uuidS))
	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Success: false,
			Msg:     err.Error(),
		})
		return
	}
	size, err := file.GetSize(f)
	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Success: false,
			Msg:     err.Error(),
		})
		return
	}
	data := &models.File{
		Name:    fileName,
		CacheID: uuidS,
		Ext:     ext,
		Path:    savePath,
		Size:    size,
	}
	if err = data.Save(); err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Success: false,
			Msg:     err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data: gin.H{
			"id": data.ID,
		},
	})

}

// FileDown 下载文件
func FileDown(c *gin.Context) {
	id := c.Param("id")

	var f models.File

	err := f.Find(id)

	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Success: false,
			Msg:     err.Error(),
		})
		return
	}

	fullPath := filepath.Join(upload.GetFileFullPath(), f.Path, f.CacheID)

	if file.CheckNotExist(fullPath) {
		c.JSON(http.StatusBadRequest, &Response{
			Success: false,
			Msg:     "error Exis file",
		})
		return
	}
	c.FileAttachment(fullPath, f.Name)
}

// FileDelete 删除文件
func FileDelete(c *gin.Context) {
	id := c.Param("id")
	var f models.File
	err := f.Delete(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Success: false,
			Msg:     "error Exis file",
		})
		return
	}
	c.Status(http.StatusNoContent)
}
