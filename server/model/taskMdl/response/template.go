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

type UploadScriptResponse struct {
	FailedIps []string `json:"failedIps"`
}

type CheckScriptResponse struct {
	FailedIps []string `json:"failedIps"`
	Script    string   `json:"script"`
}

type TaskTemplateSetResponse struct {
	taskMdl.TaskTemplateSet
	Templates []TaskTemplateSetTemplateResponse `json:"templates"`
}

type TaskTemplateSetTemplateResponse struct {
	taskMdl.TaskTemplateSetTemplate
	TemplateName string `json:"templateName"`
}

type TemplateFileListResponse struct {
	FileNames []string `json:"fileNames"`
}
