import service from '@/utils/request'

// @Summary 获取inventory列表
// @Produce  application/json
// @Param {
//  page     int
//	pageSize int
// }
// @Router /ansible/inventory/getInventoryList [post]
export const getInventoryList = (data) => {
  return service({
    url: '/ansible/inventory/getInventoryList',
    method: 'post',
    data
  })
}

// @Summary 新增inventory
// @Produce  application/json
// @Param inventory Object
// @Router /ansible/inventory/addInventory [post]
export const addInventory = (data) => {
  return service({
    url: '/ansible/inventory/addInventory',
    method: 'post',
    data
  })
}

// @Summary 删除inventory
// @Produce  application/json
// @Param ID float64
// @Router /ansible/inventory/deleteInventory [post]
export const deleteInventory = (data) => {
  return service({
    url: '/ansible/inventory/deleteInventory',
    method: 'post',
    data
  })
}

// @Summary 修改inventory
// @Produce  application/json
// @Param inventory Object
// @Router /ansible/inventory/updateInventory [post]
export const updateInventory = (data) => {
  return service({
    url: '/ansible/inventory/updateInventory',
    method: 'post',
    data
  })
}

// @Summary 根据id获取inventory
// @Produce  application/json
// @Param ID float64
// @Router /ansible/inventory/getInventoryById [post]
export const getInventoryById = (data) => {
  return service({
    url: '/ansible/inventory/getInventoryById',
    method: 'post',
    data
  })
}
