package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/application"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

type AddTaskTemplate struct {
	TaskTemplate *application.TaskTemplate `json:"taskTemplate"`
	ServerIds    []int                     `json:"serverIds"`
	AuthorityId  string                    `json:"authorityId"` // 角色ID
}

type UpdateTaskTemplate struct {
	TaskTemplate *application.TaskTemplate `json:"taskTemplate"`
	ServerIds    []int                     `json:"serverIds"`
	AuthorityId  string                    `json:"authorityId"` // 角色ID
}

type TaskTemplateSearch struct {
	application.TaskTemplate
	request.PageInfo
}
