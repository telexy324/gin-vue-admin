import service from '@/utils/request'

// @Summary 获取project列表
// @Produce  application/json
// @Param {
//  page     int
//	pageSize int
// }
// @Router /ansible/project/getProjectList [post]
export const getProjectList = (data) => {
  return service({
    url: '/ansible/project/getProjectList',
    method: 'post',
    data
  })
}

// @Summary 新增project
// @Produce  application/json
// @Param project Object
// @Router /ansible/project/addProject [post]
export const addProject = (data) => {
  return service({
    url: '/ansible/project/addProject',
    method: 'post',
    data
  })
}

// @Summary 删除project
// @Produce  application/json
// @Param ID float64
// @Router /ansible/project/deleteProject [post]
export const deleteProject = (data) => {
  return service({
    url: '/ansible/project/deleteProject',
    method: 'post',
    data
  })
}

// @Summary 修改project
// @Produce  application/json
// @Param project Object
// @Router /ansible/project/updateProject [post]
export const updateProject = (data) => {
  return service({
    url: '/ansible/project/updateProject',
    method: 'post',
    data
  })
}

// @Summary 根据id获取project
// @Produce  application/json
// @Param ID float64
// @Router /ansible/project/getProjectById [post]
export const getProjectById = (data) => {
  return service({
    url: '/ansible/project/getProjectById',
    method: 'post',
    data
  })
}
