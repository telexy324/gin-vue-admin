package response

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/ansible"
)

type SchedulesResponse struct {
	Schedules []ansible.Schedule `json:"schedules"`
}

type ScheduleResponse struct {
	Schedule ansible.Schedule `json:"schedule"`
}
