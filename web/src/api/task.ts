import service from '@/utils/request'

// @Summary 获取task列表
// @Produce  application/json
// @Param {
//  page     int
//	pageSize int
// }
// @Router /task/getTaskList [post]
export const getTaskList = (data) => {
  return service({
    url: '/task/getTaskList',
    method: 'post',
    data
  })
}

// @Summary 新增task
// @Produce  application/json
// @Param task Object
// @Router /task/addTask [post]
export const addTask = (data) => {
  return service({
    url: '/task/addTask',
    method: 'post',
    data
  })
}

// @Summary 删除task
// @Produce  application/json
// @Param ID float64
// @Router /task/deleteTask [post]
export const deleteTask = (data) => {
  return service({
    url: '/task/deleteTask',
    method: 'post',
    data
  })
}

// @Summary 修改task
// @Produce  application/json
// @Param task Object
// @Router /task/updateTask [post]
// export const updateTask = (data) => {
//   return service({
//     url: '/task/updateTask',
//     method: 'post',
//     data
//   })
// }

// @Summary 根据id获取task
// @Produce  application/json
// @Param ID float64
// @Router /task/getTaskById [post]
export const getTaskById = (data) => {
  return service({
    url: '/task/getTaskById',
    method: 'post',
    data
  })
}

// @Summary 根据id获取taskOutput
// @Produce  application/json
// @Param ID float64
// @Router //task/getTaskOutputs [post]
export const getTaskOutputs = (data) => {
  return service({
    url: '/task/getTaskOutputs',
    method: 'post',
    data
  })
}

// @Summary 根据id停止task
// @Produce  application/json
// @Param ID float64
// @Router /task/stopTask [post]
export const stopTask = (data) => {
  return service({
    url: '/task/stopTask',
    method: 'post',
    data
  })
}

// @Summary 获取task面板信息
// @Produce  application/json
// @Param empty
// @Router /task/getTaskDashboardInfo [post]
export const getTaskDashboardInfo = (data) => {
  return service({
    url: '/task/getTaskDashboardInfo',
    method: 'post',
    data
  })
}

// @Summary 获取task列表
// @Produce  application/json
// @Param {
//  page     int
//	pageSize int
// }
// @Router /task/getTaskListBySetTaskId [post]
export const getTaskListBySetTaskId = (data) => {
  return service({
    url: '/task/getTaskListBySetTaskId',
    method: 'post',
    data
  })
}
