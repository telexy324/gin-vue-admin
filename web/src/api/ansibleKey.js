import service from '@/utils/request'

// @Summary 获取key列表
// @Produce  application/json
// @Param {
//  page     int
//	pageSize int
// }
// @Router /ansible/key/getKeyList [post]
export const getKeyList = (data) => {
  return service({
    url: '/ansible/key/getKeyList',
    method: 'post',
    data
  })
}

// @Summary 新增key
// @Produce  application/json
// @Param key Object
// @Router /ansible/key/addKey [post]
export const addKey = (data) => {
  return service({
    url: '/ansible/key/addKey',
    method: 'post',
    data
  })
}

// @Summary 删除key
// @Produce  application/json
// @Param ID float64
// @Router /ansible/key/deleteKey [post]
export const deleteKey = (data) => {
  return service({
    url: '/ansible/key/deleteKey',
    method: 'post',
    data
  })
}

// @Summary 修改key
// @Produce  application/json
// @Param key Object
// @Router /ansible/key/updateKey [post]
export const updateKey = (data) => {
  return service({
    url: '/ansible/key/updateKey',
    method: 'post',
    data
  })
}

// @Summary 根据id获取key
// @Produce  application/json
// @Param ID float64
// @Router /ansible/key/getKeyById [post]
export const getKeyById = (data) => {
  return service({
    url: '/ansible/key/getKeyById',
    method: 'post',
    data
  })
}
