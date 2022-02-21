package response

import "github.com/flipped-aurora/gin-vue-admin/server/model/application"

type AdminsResponse struct {
	Admins []application.Admin `json:"admins"`
}

type AdminResponse struct {
	Admin application.Admin `json:"admin"`
}

type DepartmentsResponse struct {
	Departments []application.Department `json:"departments"`
}
