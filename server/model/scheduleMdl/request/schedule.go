package request

import "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"

// Find by id and project_id structure
type GetScheduleByTemplateId struct {
	ID float64 `json:"id" form:"id"` // 主键ID
	//ProjectId  float64 `json:"projectId" form:"projectId"`
	TemplateId float64 `json:"templateId" form:"templateId"`
	request.PageInfo
	SortInverted bool
	SortBy       string
	OrderKey     string `json:"orderKey"` // 排序
	Desc         bool   `json:"desc"`     // 排序方式:升序false(默认)|降序true
}
