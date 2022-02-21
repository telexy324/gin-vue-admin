package application

import "github.com/flipped-aurora/gin-vue-admin/server/global"

type Admins struct {
	global.GVA_MODEL
	Name         string `json:"name" gorm:"column:name"`                  // 姓名
	Mobile       int    `json:"mobile" gorm:"column:mobile"`              // 电话
	Email        int    `json:"email" gorm:"column:email"`                // 邮件
	DepartmentId int    `json:"departmentId" gorm:"column:department_id"` // 部门id
}

type Departments struct {
	global.GVA_MODEL
	Name     string `json:"name" gorm:"column:name"`          // 姓名
	ParentId int    `json:"parentId" gorm:"column:parent_id"` // 上级部门id
	LeaderId int    `json:"leaderId" gorm:"column:leader_id"` // 领导id
}
