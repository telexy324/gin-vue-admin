package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/scheduleMdl"
)

// Find by id and project_id structure
type GetScheduleList struct {
	request.PageInfo
	SortInverted bool
	SortBy       string
	OrderKey     string `json:"orderKey"` // 排序
	Desc         bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
	SystemIDs    []int  `json:"systemIds"`
	scheduleMdl.Schedule
}
