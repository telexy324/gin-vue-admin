import service from '@/utils/request'

// @Summary 获取schedule列表
// @Produce  application/json
// @Param {
//  page     int
//	pageSize int
// }
// @Router /task/schedule/getTemplateScheduleList [post]
export const getTemplateScheduleList = (data) => {
  return service({
    url: '/task/schedule/getTemplateScheduleList',
    method: 'post',
    data
  })
}

// @Summary 新增schedule
// @Produce  application/json
// @Param schedule Object
// @Router /task/schedule/addSchedule [post]
export const addSchedule = (data) => {
  return service({
    url: '/task/schedule/addSchedule',
    method: 'post',
    data
  })
}

// @Summary 删除schedule
// @Produce  application/json
// @Param ID float64
// @Router /task/schedule/deleteSchedule [post]
export const deleteSchedule = (data) => {
  return service({
    url: '/task/schedule/deleteSchedule',
    method: 'post',
    data
  })
}

// @Summary 修改schedule
// @Produce  application/json
// @Param schedule Object
// @Router /task/schedule/updateSchedule [post]
export const updateSchedule = (data) => {
  return service({
    url: '/task/schedule/updateSchedule',
    method: 'post',
    data
  })
}

// @Summary 根据id获取schedule
// @Produce  application/json
// @Param ID float64
// @Router /task/schedule/getScheduleById [post]
export const getScheduleById = (data) => {
  return service({
    url: '/task/schedule/getScheduleById',
    method: 'post',
    data
  })
}

// @Summary 验证schedule格式
// @Produce  application/json
// @Param schedule Object
// @Router /task/schedule/validateScheduleFormat [post]
export const validateScheduleCronFormat = (data) => {
  return service({
    url: '/task/schedule/validateScheduleCronFormat',
    method: 'post',
    data
  })
}

// @Tags Schedule
// @Summary 删除选中Schedule
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /task/schedule/deleteScheduleByIds [post]
export const deleteScheduleByIds = (data) => {
  return service({
    url: '/task/schedule/deleteScheduleByIds',
    method: 'post',
    data
  })
}

