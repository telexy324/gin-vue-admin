import service from '@/utils/request'

// @Summary 获取schedule列表
// @Produce  application/json
// @Param {
//  page     int
//	pageSize int
// }
// @Router /ansible/schedule/getTemplateScheduleList [post]
export const getTemplateScheduleList = (data) => {
  return service({
    url: '/ansible/schedule/getTemplateScheduleList',
    method: 'post',
    data
  })
}

// @Summary 新增schedule
// @Produce  application/json
// @Param schedule Object
// @Router /ansible/schedule/addSchedule [post]
export const addSchedule = (data) => {
  return service({
    url: '/ansible/schedule/addSchedule',
    method: 'post',
    data
  })
}

// @Summary 删除schedule
// @Produce  application/json
// @Param ID float64
// @Router /ansible/schedule/deleteSchedule [post]
export const deleteSchedule = (data) => {
  return service({
    url: '/ansible/schedule/deleteSchedule',
    method: 'post',
    data
  })
}

// @Summary 修改schedule
// @Produce  application/json
// @Param schedule Object
// @Router /ansible/schedule/updateSchedule [post]
export const updateSchedule = (data) => {
  return service({
    url: '/ansible/schedule/updateSchedule',
    method: 'post',
    data
  })
}

// @Summary 根据id获取schedule
// @Produce  application/json
// @Param ID float64
// @Router /ansible/schedule/getScheduleById [post]
export const getScheduleById = (data) => {
  return service({
    url: '/ansible/schedule/getScheduleById',
    method: 'post',
    data
  })
}

// @Summary 验证schedule格式
// @Produce  application/json
// @Param schedule Object
// @Router /ansible/schedule/validateScheduleFormat [post]
export const validateScheduleFormat = (data) => {
  return service({
    url: '/ansible/schedule/validateScheduleFormat',
    method: 'post',
    data
  })
}
