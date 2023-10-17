import service from '@/utils/request'

// @Tags ApplicationRecord
// @Summary 删除ApplicationRecord
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body model.ApplicationRecord true "删除ApplicationRecord"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /cmdb/deleteApplicationRecord [delete]
export const deleteApplicationRecord = (data) => {
  return service({
    url: '/cmdb/deleteApplicationRecord',
    method: 'post',
    data
  })
}

// @Tags ApplicationRecord
// @Summary 删除ApplicationRecord
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "删除ApplicationRecord"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /cmdb/deleteApplicationRecordByIds [delete]
export const deleteApplicationRecordByIds = (data) => {
  return service({
    url: '/cmdb/deleteApplicationRecordByIds',
    method: 'post',
    data
  })
}

// @Tags ApplicationRecord
// @Summary 分页获取ApplicationRecord列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.PageInfo true "分页获取ApplicationRecord列表"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /cmdb/getApplicationRecordList [get]
export const getApplicationRecordList = (params) => {
  return service({
    url: '/cmdb/getApplicationRecordList',
    method: 'get',
    params
  })
}
