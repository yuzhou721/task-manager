package models

import (
	"errors"
	"os"
	"path/filepath"
	"task/pkg/file"
	"task/pkg/upload"

	"github.com/jinzhu/gorm"
)

// File 用户上传文件保存
type File struct {
	gorm.Model
	Name    string
	CacheID string
	Path    string
	Size    int
	Ext     string
}

// Save 保存文件
func (f *File) Save() (err error) {
	if err = db.Create(f).Error; err != nil {
		return err
	}
	return
}

// Find 查询文件
func (f *File) Find(id string) (err error) {
	if err = db.First(f, id).Error; err != nil {
		return
	}
	return
}

// Delete 删除文件
func (f *File) Delete(id string) (err error) {
	myDb := db.Begin()
	err = f.Find(id)
	if err != nil {
		myDb.Rollback()
		return
	}
	if err = myDb.Delete(&f).Error; err != nil {
		myDb.Rollback()
		return
	}
	fullPath := filepath.Join(upload.GetFileFullPath(), f.Path, f.CacheID)
	if file.CheckNotExist(fullPath) {
		myDb.Rollback()
		return errors.New("File Not Exist")
	}
	err = os.Remove(fullPath)
	if err != nil {
		myDb.Rollback()
		return errors.New("File Remove Error")
	}

	myDb.Commit()

	return
}
