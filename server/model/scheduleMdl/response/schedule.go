package response

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/scheduleMdl"
)

type SchedulesResponse struct {
	Schedules []scheduleMdl.Schedule `json:"schedules"`
}

type ScheduleResponse struct {
	Schedule scheduleMdl.Schedule `json:"schedule"`
	SystemId int                  `json:"systemId"`
}
