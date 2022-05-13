package request

// Find by id and project_id structure
type AddUserByProjectId struct {
	ProjectId float64 `json:"project_id" form:"project_id"`
	UserId    float64 `json:"user_id" form:"user_id"`
	Admin     int     `json:"admin" form:"admin"`
}

// Find by id and project_id structure
type DeleteUserByProjectId struct {
	ProjectId float64 `json:"project_id" form:"project_id"`
	UserId    float64 `json:"user_id" form:"user_id"`
}