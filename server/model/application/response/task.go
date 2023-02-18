package response

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/application"
)

type TaskTemplatesResponse struct {
	TaskTemplates []application.TaskTemplate `json:"taskTemplates"`
}

type TaskTemplateResponse struct {
	TaskTemplate application.TaskTemplate `json:"taskTemplate"`
}

