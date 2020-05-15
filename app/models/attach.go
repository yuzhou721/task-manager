package models

import (
	"strconv"

	"github.com/jinzhu/gorm"
)

//Attach 附件
type Attach struct {
	gorm.Model
	Name   string `json:"name"`
	Status string `json:"status"`
	UID    string `json:"uid"`
	FileID uint   `json:"fileID"`
	TaskID uint   `json:"-"`
}

// Delete 根据ID删除
func (a *Attach) Delete(id string) (err error) {
	myDB := db.Begin()
	if err = myDB.First(a, id).Error; err != nil {
		myDB.Rollback()
		return
	}
	if err = myDB.Delete(a).Error; err != nil && err != gorm.ErrRecordNotFound {
		myDB.Rollback()
		return
	}
	var f File
	err = f.Delete(strconv.Itoa(int(a.FileID)))
	if err != nil {
		myDB.Rollback()
		return
	}
	myDB.Commit()
	return
}
