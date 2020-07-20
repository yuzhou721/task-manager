package models

type Manager struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

//查询是否是任务中心管理员
//return String  "1" 是管理员   "0" 不是管理员
func (m *Manager) IsAdmin(phone string) (err error) {
	if err = db.Where("phone=?", phone).First(m).Error; err != nil {
		return
	}
	return
}
