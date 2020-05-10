package upload

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"task/conf"
	"task/pkg/file"
	"task/pkg/utils"
)

//GetFileName 获取文件md5名称
func GetFileName(name string) string {
	ext := path.Ext(name)
	fileName := strings.TrimSuffix(name, ext)
	fileName = utils.EncodeMD5(fileName)

	return fileName + ext
}

//GetFilePath 获取保存路径
func GetFilePath() string {
	return conf.Config.App.SavePath
}

//GetFileFullPath 获取全路径
func GetFileFullPath() string {
	return filepath.Join(conf.Config.App.RuntimePath, GetFilePath())
}

// func CheckFileExt(fileName string) bool {
// 	ext := file.GetExt(fileName)
// 	for _, allowExt := range setting.AppSetting.ImageAllowExts {
// 		if strings.ToUpper(allowExt) == strings.ToUpper(ext) {
// 			return true
// 		}
// 	}

// 	return false
// }

// func CheckFileSize(f multipart.File) bool {
// 	size, err := file.GetSize(f)
// 	if err != nil {
// 		log.Println(err)
// 		return false
// 	}

// 	return size <= setting.AppSetting.ImageMaxSize
// }

// CheckFile 检查文件
func CheckFile(src string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os.Getwd err: %v", err)
	}

	err = file.IsNotExistMkDir(dir + "/" + src)
	if err != nil {
		return fmt.Errorf("file.IsNotExistMkDir err: %v", err)
	}

	perm := file.CheckPermission(src)
	if perm == true {
		return fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}

	return nil
}
