package scheduleSvr

import "github.com/flipped-aurora/gin-vue-admin/server/service/task"

type ServiceGroup struct {
	ansible.TaskService
	EnvironmentService
	InventoryService
	KeyService
	ProjectService
	SchedulesService
	ansible.TemplatesService
	UserService
}
