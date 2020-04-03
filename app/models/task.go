package models

import (
	"time"

	"task/pkg/utils"

	"github.com/jinzhu/gorm"
)

//Task 任务数据
type Task struct {
	gorm.Model
	Create       string     //创建人
	CreateOpenID string     //创建人openId
	Phone        string     //创建人手机号
	DueDate      *time.Time //截至时间
	StartTime    *time.Time //开始时间
	EndTime      *time.Time //结束时间
	Status       *int       //状态 1 未完成 2 已完成
	Title        string     //标题
	Content      string     //内容
	ParentID     *uint      //主任务id
	Attach       []uint     //附件id
}

const (
	//RoleAdmin 管理员
	RoleAdmin = 1
	//RoleDept 部门领导
	RoleDept = 2
	//RoleNormal 普通人员
	RoleNormal = 3
)

//Save 保存
func (t *Task) Save() (err error) {
	if err = db.Create(t).Error; err != nil {
		return err
	}
	return
}

//SaveList 列表保存
func (t *Task) SaveList(tasks *[]Task) (err error) {
	tx := db.Begin()
	for _, v := range *tasks {
		if err = tx.Create(&v).Error; err != nil {
			tx.Rollback()
		}
	}
	tx.Commit()
	return err
}

//FindByID 根据ID获取任务信息
func (t *Task) FindByID(id string) (to *Task, err error) {
	if err = db.Find(&to, id).Error; err != nil {
		return
	}
	return
}

//FindOne 根据条件查询数据
func (t *Task) FindOne() (r *Task, err error) {
	if err = db.Where(t).First(&r).Error; err != nil && err != gorm.ErrRecordNotFound {
		return
	}
	return
}

//Update 根据ID修改实体类
func (t *Task) Update() (err error) {
	if err = db.Model(t).Update(t).Error; err != nil && err != gorm.ErrRecordNotFound {
		return
	}
	return
}

//Delete 根据id删除实体类
func (t *Task) Delete() (err error) {
	if err = db.Delete(t).Error; err != nil {
		return
	}
	return
}

//FindList 根据条件以及角色分页查询数据
func (t *Task) FindList(role int, search string, page, pageSize int) (tasks []Task, count int, err error) {

	// 查询条件
	//TODO:根据角色查询
	searchDb := db.Where(t)
	// 统计总数
	searchDb = searchDb.Count(&count)
	// 分页获取
	searchDb = searchDb.Offset(utils.GetPageOffset(page, pageSize)).Limit(pageSize)
	// TODO: 排序

	searchDb.Find(&tasks)

	return
}
