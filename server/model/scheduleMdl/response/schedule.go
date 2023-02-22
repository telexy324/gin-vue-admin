package response

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/schedule"
)

type SchedulesResponse struct {
	Schedules []scheduleMdl.Schedule `json:"schedules"`
}

type ScheduleResponse struct {
	Schedule scheduleMdl.Schedule `json:"schedule"`
}
