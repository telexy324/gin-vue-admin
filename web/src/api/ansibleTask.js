import service from '@/utils/request'

// @Summary 获取task列表
// @Produce  application/json
// @Param {
//  page     int
//	pageSize int
// }
// @Router /ansible/task/getTaskList [post]
export const getTaskList = (data) => {
  return service({
    url: '/ansible/task/getTaskList',
    method: 'post',
    data
  })
}

// @Summary 新增task
// @Produce  application/json
// @Param task Object
// @Router /ansible/task/addTask [post]
export const addTask = (data) => {
  return service({
    url: '/ansible/task/addTask',
    method: 'post',
    data
  })
}

// @Summary 删除task
// @Produce  application/json
// @Param ID float64
// @Router /ansible/task/deleteTask [post]
export const deleteTask = (data) => {
  return service({
    url: '/ansible/task/deleteTask',
    method: 'post',
    data
  })
}

// @Summary 修改task
// @Produce  application/json
// @Param task Object
// @Router /ansible/task/updateTask [post]
export const updateTask = (data) => {
  return service({
    url: '/ansible/task/updateTask',
    method: 'post',
    data
  })
}

// @Summary 根据id获取task
// @Produce  application/json
// @Param ID float64
// @Router /ansible/task/getTaskById [post]
export const getTaskById = (data) => {
  return service({
    url: '/ansible/task/getTaskById',
    method: 'post',
    data
  })
}

// @Summary 根据id获取taskOutput
// @Produce  application/json
// @Param ID float64
// @Router //ansible/task/getTaskOutputs [post]
export const getTaskOutputs = (data) => {
  return service({
    url: '/ansible/task/getTaskOutputs',
    method: 'post',
    data
  })
}

// @Summary 根据id停止task
// @Produce  application/json
// @Param ID float64
// @Router /ansible/task/stopTask [post]
export const stopTask = (data) => {
  return service({
    url: '/ansible/task/stopTask',
    method: 'post',
    data
  })
}