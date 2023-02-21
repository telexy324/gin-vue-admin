package response

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/task"
)

type TaskTemplatesResponse struct {
	Templates []task.TaskTemplate `json:"taskTemplates"`
}

type TaskTemplateResponse struct {
	Template task.TaskTemplate `json:"taskTemplate"`
}
