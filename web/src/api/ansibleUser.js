import service from '@/utils/request'

// @Summary 获取user列表
// @Produce  application/json
// @Param {
//  page     int
//	pageSize int
// }
// @Router /ansible/user/getProjectUsers [post]
export const getProjectUsers = (data) => {
  return service({
    url: '/ansible/user/getProjectUsers',
    method: 'post',
    data
  })
}

// @Summary 新增user
// @Produce  application/json
// @Param user Object
// @Router /ansible/user/addUser [post]
export const addUser = (data) => {
  return service({
    url: '/ansible/user/addUser',
    method: 'post',
    data
  })
}

// @Summary 删除user
// @Produce  application/json
// @Param ID float64
// @Router /ansible/user/deleteUser [post]
export const deleteUser = (data) => {
  return service({
    url: '/ansible/user/deleteUser',
    method: 'post',
    data
  })
}

// @Summary 修改user
// @Produce  application/json
// @Param user Object
// @Router /ansible/user/updateUser [post]
export const updateUser = (data) => {
  return service({
    url: '/ansible/user/updateUser',
    method: 'post',
    data
  })
}
