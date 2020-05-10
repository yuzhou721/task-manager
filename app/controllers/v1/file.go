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
	_, header, err := c.Request.FormFile("file")
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
	data := &models.FileInfo{
		FileName: fileName,
		ID:       uuidS,
		Ext:      ext,
		Path:     savePath,
	}
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
	})

}

// FileDown 下载文件
func FileDown(c *gin.Context) {
	var fileI models.FileInfo
	err := c.BindQuery(&fileI)
	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Success: false,
			Msg:     err.Error(),
		})
	}

	fullPath := filepath.Join(upload.GetFileFullPath(), fileI.Path, fileI.ID)

	if file.CheckNotExist(fullPath) {
		c.JSON(http.StatusBadRequest, &Response{
			Success: false,
			Msg:     "error Exis file",
		})
		return
	}
	c.FileAttachment(fullPath, fileI.FileName)
}
