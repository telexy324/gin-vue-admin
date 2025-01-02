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
	TemplatesInner []TaskTemplateSetResponseInner `json:"templates"`
}

type TaskTemplateSetTemplateResponse struct {
	taskMdl.TaskTemplateSetTemplate
	TemplateName string `json:"templateName"`
}

type TaskTemplateSetResponseInner struct {
	Seq         int                               `json:"seq"`
	Templates   []TaskTemplateSetTemplateResponse `json:"templates"`
	TemplateIds []int                             `json:"templateIds"`
}

type TemplateFileListResponse struct {
	IsTop            bool       `json:"isTop"`
	FileInfos        []FileInfo `json:"fileInfos"`
	CurrentDirectory string     `json:"currentDirectory"`
}

type FileInfo struct {
	FileName  string `json:"fileName"`
	Directory bool   `json:"directory"`
}
