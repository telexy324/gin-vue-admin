package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/taskMdl"
)

//// Find by id and project_id structure
//type AddTaskByProjectId struct {
//	ProjectId float64 `json:"projectId" form:"projectId"`
//	task.Task
//}

// Find by id and project_id structure
type TemplateScriptRequest struct {
	ID       float64 `json:"id" form:"id"` // 主键ID
	ServerId float64 `json:"serverId" form:"serverId"`
	Detail   bool    `json:"detail" form:"detail"`
}

type AddSet struct {
	taskMdl.TaskTemplateSet
	Templates []taskMdl.TaskTemplateSetTemplate `json:"templates"`
}

type DownLoadFileRequest struct {
	ID        float64 `json:"id" form:"id"` // 主键ID
	File      string  `json:"file" form:"file"`
	TargetIds []int   `json:"targetIds" form:"targetIds"`
}

type FileListRequest struct {
	ID        float64 `json:"id" form:"id"` // 主键ID
	Directory string  `json:"directory" form:"directory"`
	TargetIds []int   `json:"targetIds" form:"targetIds"`
}

type ProcessTaskRequest struct {
	ID          float64  `json:"id" form:"id"`                   // 主键ID
	CommandVars []string `json:"commandVars" form:"commandVars"` // 结束时间
}
