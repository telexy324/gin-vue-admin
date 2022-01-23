import service from '@/utils/request'

// @Summary 获取server列表
// @Produce  application/json
// @Param {
//  page     int
//	pageSize int
// }
// @Router /menu/getServerList [post]
export const getServerList = (data) => {
  return service({
    url: '/menu/getServerList',
    method: 'post',
    data
  })
}

// @Summary 新增server
// @Produce  application/json
// @Param menu Object
// @Router /menu/addServer [post]
export const addServer = (data) => {
  return service({
    url: '/menu/addServer',
    method: 'post',
    data
  })
}

// @Summary 删除server
// @Produce  application/json
// @Param ID float64
// @Router /menu/deleteServer [post]
export const deleteServer = (data) => {
  return service({
    url: '/menu/deleteServer',
    method: 'post',
    data
  })
}

// @Summary 修改server
// @Produce  application/json
// @Param server Object
// @Router /menu/updateServer [post]
export const updateServer = (data) => {
  return service({
    url: '/menu/updateServer',
    method: 'post',
    data
  })
}

// @Tags server
// @Summary 根据id获取服务器
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body api.GetById true "根据id获取服务器"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /menu/getServerById [post]
export const getServerById = (data) => {
  return service({
    url: '/menu/getServerById',
    method: 'post',
    data
  })
}
