package response

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/schedule"
)

type SchedulesResponse struct {
	Schedules []schedule.Schedule `json:"schedules"`
}

type ScheduleResponse struct {
	Schedule schedule.Schedule `json:"schedule"`
}
