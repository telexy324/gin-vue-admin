package request

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
