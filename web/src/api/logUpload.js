import service from '@/utils/request'

// @Summary 获取server列表
// @Produce  application/json
// @Param {
//  page     int
//	pageSize int
// }
// @Router /logUpload/getServerList [post]
export const getServerList = (data) => {
  return service({
    url: '/logUpload/getServerList',
    method: 'post',
    data
  })
}

// @Summary 新增secret
// @Produce  application/json
// @Param menu Object
// @Router /logUpload/addServer [post]
export const addServer = (data) => {
  return service({
    url: '/logUpload/addServer',
    method: 'post',
    data
  })
}

// @Summary 删除secret
// @Produce  application/json
// @Param ID float64
// @Router /logUpload/deleteServer [post]
export const deleteServer = (data) => {
  return service({
    url: '/logUpload/deleteServer',
    method: 'post',
    data
  })
}

// @Summary 修改secret
// @Produce  application/json
// @Param secret Object
// @Router /logUpload/updateServer [post]
export const updateServer = (data) => {
  return service({
    url: '/logUpload/updateServer',
    method: 'post',
    data
  })
}

// @Tags secret
// @Summary 根据id获取服务器
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body api.GetById true "根据id获取服务器"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /logUpload/getServerById [post]
export const getServerById = (data) => {
  return service({
    url: '/logUpload/getServerById',
    method: 'post',
    data
  })
}

// @Tags Server
// @Summary 删除选中Server
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /logUpload/deleteServerByIds [post]
export const deleteServerByIds = (data) => {
  return service({
    url: '/logUpload/deleteServerByIds',
    method: 'post',
    data
  })
}

// @Summary 获取secret列表
// @Produce  application/json
// @Param {
//  page     int
//	pageSize int
// }
// @Router /logUpload/getSecretList [post]
export const getSecretList = (data) => {
  return service({
    url: '/logUpload/getSecretList',
    method: 'post',
    data
  })
}

// @Summary 新增secret
// @Produce  application/json
// @Param menu Object
// @Router /logUpload/addSecret [post]
export const addSecret = (data) => {
  return service({
    url: '/logUpload/addSecret',
    method: 'post',
    data
  })
}

// @Summary 删除secret
// @Produce  application/json
// @Param ID float64
// @Router /logUpload/deleteSecret [post]
export const deleteSecret = (data) => {
  return service({
    url: '/logUpload/deleteSecret',
    method: 'post',
    data
  })
}

// @Summary 修改secret
// @Produce  application/json
// @Param secret Object
// @Router /logUpload/updateSecret [post]
export const updateSecret = (data) => {
  return service({
    url: '/logUpload/updateSecret',
    method: 'post',
    data
  })
}

// @Tags secret
// @Summary 根据id获取服务器
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body api.GetById true "根据id获取服务器"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /logUpload/getSecretById [post]
export const getSecretById = (data) => {
  return service({
    url: '/logUpload/getSecretById',
    method: 'post',
    data
  })
}

// @Tags Secret
// @Summary 删除选中Secret
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /logUpload/deleteSecretByIds [post]
export const deleteSecretByIds = (data) => {
  return service({
    url: '/logUpload/deleteSecretByIds',
    method: 'post',
    data
  })
}
