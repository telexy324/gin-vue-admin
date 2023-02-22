package response

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/task"
)

type TaskTemplatesResponse struct {
	Templates []taskMdl.TaskTemplate `json:"taskTemplates"`
}

type TaskTemplateResponse struct {
	Template taskMdl.TaskTemplate `json:"taskTemplate"`
}
