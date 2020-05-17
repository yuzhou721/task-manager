package models

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"task/conf"
	"task/pkg/utils"
	"task/pkg/yzj"

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
	Attach                 []Attach   //附件
	ParentAttach           []Attach   `gorm:"-"`
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
	// TaskStatusUndo 任务未完成
	TaskStatusUndo = 1
	// TaskTypeMain 主任务类型
	TaskTypeMain = 1
	// TaskTypeDept 部门类型
	TaskTypeDept = 2
	// TaskTypeNomal 普通类型
	TaskTypeNomal = 3
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
		v.sendTodo()
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
		v.sendTodo()
	}
	tx.Commit()
	return err
}

//FindByID 根据ID获取任务信息
func (t *Task) FindByID(id string) (err error) {
	if err = db.Where(id).Preload("Attach").First(&t).Error; err != nil {
		return
	}
	if err = t.findChildren(); err != nil {
		return
	}
	attach, err := t.findParentAttach()
	t.ParentAttach = attach

	return
}

func (t *Task) findChildren() (err error) {
	var children []*Task
	if err = db.Model(t).Where("parent_id = ?", t.ID).Order("created_at DESC").Find(&children).Error; err != nil {
		return err
	}
	t.Children = children
	return nil
}

func (t *Task) findParentAttach() (attach []Attach, err error) {
	if t.ParentID == nil {
		return
	}
	db.Where("task_id = ?", t.ParentID).Find(&attach)
	return
}

func (t *Task) findParent() (parent *Task, err error) {
	parent = &Task{}
	if t.ParentID == nil {
		return
	}
	if err = db.Model(t).Where(*t.ParentID).Preload("Attach").First(parent).Error; err != nil {
		return
	}
	return
}

//FindOne 根据条件查询数据
func (t *Task) FindOne() (r *Task, err error) {
	r = &Task{}
	if err = db.Where(&t).First(r).Error; err != nil && err != gorm.ErrRecordNotFound {
		return
	}
	return
}

//Update 根据ID修改实体类
func (t *Task) Update() (err error) {
	if err = db.Model(t).Update(t).Error; err != nil && err != gorm.ErrRecordNotFound {
		return
	}
	if t.Status == TaskStatusDone {
		t.clearTodo()
		t.sendNotify(notifyTypeStatus)
	}
	return
}

//AfterUpdate 钩子函数 用于更新数据状态以后 更新父类状态
func (t *Task) AfterUpdate(tx *gorm.DB) (err error) {
	if t.Status == TaskStatusDone {
		if err = t.updateParentDone(tx); err != nil {
			return
		}
	}
	return
}

func (t *Task) updateParentDone(tx *gorm.DB) (err error) {
	if t.ParentID == nil {
		return
	}
	// 查询父任务
	var pt Task
	if t.ParentID == nil {
		return
	}
	if err = tx.Model(t).Where(*t.ParentID).First(&pt).Error; err != nil {
		return
	}
	// 如果没有父任务 跳出
	if err != nil && err == gorm.ErrRecordNotFound {
		return nil
	}
	// 未知错误 跳出
	if err != nil {
		return
	}
	if pt.Status == TaskStatusDone {
		return
	}

	// 查询子任务 因为在同一个事务，所以可以脏读取
	var brother []*Task
	if err = tx.Model(t).Where("parent_id = ?", t.ParentID).Find(&brother).Error; err != nil && err != gorm.ErrRecordNotFound {
		return
	}

	// 如果子任务查询出来完成任务数量等于子任务数量 说明全部任务已完成 修改父任态
	cl := len(brother)
	if cl == 0 {
		return
	}
	var ci = 0
	for _, v := range brother {
		if v.Status == TaskStatusDone {
			ci++
		}
	}

	if cl == ci {
		pt.Status = TaskStatusDone
		if err = tx.Model(t).Update(&pt).Error; err != nil && err != gorm.ErrRecordNotFound {
			return
		}
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
func (t *Task) FindList(role int, openID, orgID, search string, status, page, pageSize int) (tasks []Task, count int, err error) {

	// 查询条件
	searchDb := db.Model(t)
	if role == RoleAdmin {
		searchDb = searchDb.Where("parent_id is null")
	} else if role == RoleDept {
		searchDb = searchDb.Where("(assigner_id = ? and type = '2' ) or designated_department_id = ? ", openID, orgID)
	} else {
		searchDb = searchDb.Where("designated_person_id = ?", openID)
	}
	// 状态查询
	searchDb = searchDb.Where("status = ?", status)
	// 模糊查询
	if search != "" {
		search = "%" + search + "%"
		searchDb = searchDb.Where("title like ? or content = ?", search, search)
	}

	// 统计总数
	searchDb = searchDb.Count(&count)
	// 分页获取
	searchDb = searchDb.Offset(utils.GetPageOffset(page, pageSize)).Limit(pageSize)
	// TODO: 排序
	searchDb = searchDb.Order("created_at DESC")

	err = searchDb.Find(&tasks).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}

	return
}

//CountTaskPercent 统计任务完成数
func (t *Task) CountTaskPercent(ID string) (percent float32, err error) {
	var mPercent float32
	// 查询数据
	err = t.FindByID(ID)
	if err != nil && err == gorm.ErrRecordNotFound {
		return
	}
	total := len(t.Children)
	if total == 0 {
		if t.Status == TaskStatusDone {
			mPercent++
		}
		return
	}
	for _, v := range t.Children {
		if v.Status == TaskStatusDone {
			mPercent++
		} else if v.Status == TaskStatusUndo {
			cPercent, err := v.CountTaskPercent(strconv.Itoa(int(v.ID)))
			if err != nil {
				continue
			}
			mPercent += cPercent
		}
	}

	percent = mPercent / float32(total)
	return
}

//轻应用名称
const itemTitle = "任务中心"

func (t *Task) sendTodo() (err error) {
	y := &yzj.Yzj{
		AppID:  conf.Config.Yzj.AppID,
		Secret: conf.Config.Yzj.Secret,
		Scope:  yzj.YzjScopeApp,
	}
	var (
		openIDs []string
		title   string
		content string
		url     string
		headImg string
	)

	title = "你有一个新任务，请及时处理！"
	content = fmt.Sprintf("任务名称：%v", t.Title)
	url = fmt.Sprintf("%v/#/detail/%v/edit", conf.Config.App.UIURL, t.ID)
	if t.DesignatedDepartmentID != "" {
		org, err := y.GetOrg(t.DesignatedDepartmentID)
		if err != nil {
			log.Printf("error get org:%v", err.Error())
		}
		for _, v := range org.InChargers {
			openIDs = append(openIDs, v.OpenID)
		}
	} else {
		openIDs = append(openIDs, t.DesignatedPersonID)
	}
	err = y.GenerateTODO(strconv.Itoa(int(t.ID)), openIDs, title, content, itemTitle, url, headImg)
	if err != nil {
		return
	}
	return
}

func (t *Task) clearTodo() (err error) {
	y := &yzj.Yzj{
		AppID:  conf.Config.Yzj.AppID,
		Secret: conf.Config.Yzj.Secret,
		Scope:  yzj.YzjScopeApp,
	}
	if t.Status != TaskStatusDone {
		return
	}
	if t.ParentID != nil {
		pt, err := t.findParent()
		if err != nil {
			return err
		}
		pt.clearTodo()
	}
	err = y.OprateTodo(strconv.Itoa(int(t.ID)), []string{strconv.Itoa(int(t.ID))}, 0, 0, 0)
	if err != nil {
		return
	}
	return
}

const (
	notifyTypeStatus = iota
	notifyTypeTimeout
)

func (t *Task) sendNotify(Type int) (err error) {
	y := &yzj.Yzj{
		EID:       conf.Config.Yzj.EID,
		PubSecret: conf.Config.Yzj.PubSecret,
		PubID:     conf.Config.Yzj.PubID,
	}

	var (
		url     string
		openID  string
		content string
	)

	if Type == notifyTypeStatus {
		openID = t.AssignerID
		if t.ParentID != nil {
			pt, err := t.findParent()
			if err != nil {
				return err
			}
			pt.sendNotify(notifyTypeStatus)
		}
		percent, err := t.CountTaskPercent(strconv.Itoa(int(t.ID)))
		if err != nil {
			return err
		}
		if t.Type == TaskTypeNomal && t.Status == TaskStatusDone {
			percent = 1
		}
		percent *= 100
		content = fmt.Sprintf("任务进度有更新，请知悉！\n任务名称:%v\n任务状态:%.2f%%", t.Title, percent)
	} else {
		content = fmt.Sprintf("您有任务即将逾期，请尽快处理！\n任务名称:%v", t.Title)
	}

	url = fmt.Sprintf("%v/#/detail/%v/view", conf.Config.App.UIURL, t.ID)
	err = y.GenerateNotify(content, url, []string{openID})
	if err != nil {
		return
	}
	return
}

//Remind 提醒逾期任务
func Remind(time time.Time) (err error) {
	var tasks []*Task
	if err = db.Model(&Task{}).Where("date(due_date) = date(?)", time.Local()).Find(&tasks).Error; err != nil && err != gorm.ErrRecordNotFound {
		return
	}

	for _, v := range tasks {
		v.sendNotify(notifyTypeTimeout)
	}
	return
}
