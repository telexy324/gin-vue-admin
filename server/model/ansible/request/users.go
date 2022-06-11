package request

// Find by id and project_id structure
type AddUserByProjectId struct {
	ProjectId float64 `json:"projectId" form:"projectId"`
	UserId    float64 `json:"userId" form:"userId"`
	Admin     int     `json:"admin" form:"admin"`
}

// Find by id and project_id structure
type DeleteUserByProjectId struct {
	ProjectId float64 `json:"projectId" form:"projectId"`
	UserId    float64 `json:"userId" form:"userId"`
}