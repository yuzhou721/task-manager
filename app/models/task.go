package models

import (
	"time"

	"task/pkg/utils"

	"github.com/jinzhu/gorm"
)

//Task 任务数据
type Task struct {
	gorm.Model
	Assigner               string     //派遣人
	AssignerID             string     //派遣人openId
	Phone                  string     //创建人手机号
	DueDate                *time.Time `gorm:"type:date" time_format:"2006-01-02"` //截至时间
	StartTime              *time.Time `gorm:"type:date" time_format:"2006-01-02"` //开始时间
	EndTime                *time.Time `gorm:"type:date" time_format:"2006-01-02"` //结束时间
	Status                 uint       //状态 1 未完成 2 已完成
	Title                  string     //标题
	Content                string     //内容
	ParentID               *uint      //主任务id
	Attach                 string     //附件id
	Type                   uint       // 任务类型 1.主任务 2.部门任务 3.人员任务
	DesignatedPerson       string     //指定人
	DesignatedPersonID     string     //指定人OpenId
	DesignatedDepartment   string     //指定部门
	DesignatedDepartmentID string     //指定部门Id
	Children               []*Task    `gorm:"-"`
}

const (
	//RoleAdmin 管理员
	RoleAdmin = 1
	//RoleDept 部门领导
	RoleDept = 2
	//RoleNormal 普通人员
	RoleNormal = 3
	//TaskStatusDone 任务完成状态
	TaskStatusDone = 2
	TaskStatusUndo = 1
	TaskTypeMain   = 1
	TaskTypeDept   = 2
	TaskTypeNomal  = 3
)

//Save 保存
func (t *Task) Save() (err error) {
	if err = db.Create(t).Error; err != nil {
		return err
	}
	return
}

//SaveOrUpdateList 列表保存
func (t *Task) SaveOrUpdateList(tasks *[]Task) (err error) {
	tx := db.Begin()
	for _, v := range *tasks {
		if v.ID != 0 {
			if err = tx.Update(&v).Error; err != nil {
				tx.Rollback()
			}
		} else {

			if err = tx.Create(&v).Error; err != nil {
				tx.Rollback()
			}
		}
	}
	tx.Commit()
	return err
}

//SaveMasterAndSlave 创建时候保存主从结构 如果数据长度>1则 下标为0的数据是主数据 其他是从数据
func (t *Task) SaveMasterAndSlave(tasks []Task) (err error) {
	tx := db.Begin()
	main := tasks[0]
	if err = tx.Create(&main).Error; err != nil {
		return
	}
	parentID := main.ID
	for _, v := range tasks[1:] {
		v.ParentID = &parentID
		if err = tx.Create(&v).Error; err != nil {
			tx.Rollback()
		}
	}
	tx.Commit()
	return err
}

//FindByID 根据ID获取任务信息
func (t *Task) FindByID(id string) (err error) {
	if err = db.Find(t, id).Error; err != nil {
		return
	}
	if err = t.findChildren(); err != nil {
		return
	}
	return
}

func (t *Task) findChildren() (err error) {
	var children []*Task
	if err = db.Model(t).Where("parent_id = ?", t.ID).Find(&children).Error; err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	t.Children = children
	return nil
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
func (t *Task) FindList(role int, openID, orgID, search string, page, pageSize int) (tasks []Task, count int, err error) {

	// 查询条件
	//TODO:根据角色查询
	searchDb := db.Model(t)
	if role == RoleAdmin {
		searchDb = searchDb.Where("parent_id is null")
	} else if role == RoleDept {
		searchDb = searchDb.Where("assigner_id = ? and type = '2'", openID)
		searchDb = searchDb.Or("designated_department_id = ? ", orgID)
	} else {
		searchDb = searchDb.Where("designated_person_id = ?", openID)
	}
	// 模糊查询
	if search != "" {
		search = "%" + search + "%"
		searchDb.Where("title like ?", search)
	}

	// 统计总数
	searchDb = searchDb.Count(&count)
	// 分页获取
	searchDb = searchDb.Offset(utils.GetPageOffset(page, pageSize)).Limit(pageSize)
	// TODO: 排序

	err = searchDb.Find(&tasks).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}

	return
}

//CountTaskByParentID 统计任务完成数
func (t *Task) CountTaskByParentID(pID string) (total, complate, undo int, err error) {
	var tasks []*Task
	// 查询总数
	err = db.Model(&t).Where("parent_id = ?", pID).Find(&tasks).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		return
	}
	total = len(tasks)
	for _, v := range tasks {
		if v.Status == TaskStatusDone {
			complate++
		} else if v.Status == TaskStatusUndo {
			undo++
		}
	}
	return
}
