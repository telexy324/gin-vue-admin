package response

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/taskMdl"
)

type TaskTemplatesResponse struct {
	Templates []taskMdl.TaskTemplate `json:"taskTemplates"`
}

type TaskTemplateResponse struct {
	Template taskMdl.TaskTemplate `json:"taskTemplate"`
}

type TemplateScriptResponse struct {
	Exist  bool   `json:"exist"`
	Script string `json:"script"`
}

type UploadScriptResponse struct {
	FailedIps []string `json:"failedIps"`
}

type TaskTemplateSetResponse struct {
	taskMdl.TaskTemplateSet
	Templates []taskMdl.TaskTemplateSetTemplate `json:"templates"`
}
