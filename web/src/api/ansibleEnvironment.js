import service from '@/utils/request'

// @Summary 获取environment列表
// @Produce  application/json
// @Param {
//  page     int
//	pageSize int
// }
// @Router /ansible/environment/getEnvironmentList [post]
export const getEnvironmentList = (data) => {
  return service({
    url: '/ansible/environment/getEnvironmentList',
    method: 'post',
    data
  })
}

// @Summary 新增environment
// @Produce  application/json
// @Param environment Object
// @Router /ansible/environment/addEnvironment [post]
export const addEnvironment = (data) => {
  return service({
    url: '/ansible/environment/addEnvironment',
    method: 'post',
    data
  })
}

// @Summary 删除environment
// @Produce  application/json
// @Param ID float64
// @Router /ansible/environment/deleteEnvironment [post]
export const deleteEnvironment = (data) => {
  return service({
    url: '/ansible/environment/deleteEnvironment',
    method: 'post',
    data
  })
}

// @Summary 修改environment
// @Produce  application/json
// @Param environment Object
// @Router /ansible/environment/updateEnvironment [post]
export const updateEnvironment = (data) => {
  return service({
    url: '/ansible/environment/updateEnvironment',
    method: 'post',
    data
  })
}

// @Summary 根据id获取environment
// @Produce  application/json
// @Param ID float64
// @Router /ansible/environment/getEnvironmentById [post]
export const getEnvironmentById = (data) => {
  return service({
    url: '/ansible/environment/getEnvironmentById',
    method: 'post',
    data
  })
}
