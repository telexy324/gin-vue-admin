import service from '@/utils/request'
import download from '../utils/download'

// @Summary 获取server列表
// @Produce  application/json
// @Param {
//  page     int
//	pageSize int
// }
// @Router /cmdb/getServerList [post]
export const getServerList = (data) => {
  return service({
    url: '/cmdb/getServerList',
    method: 'post',
    data
  })
}

// @Summary 新增server
// @Produce  application/json
// @Param menu Object
// @Router /cmdb/addServer [post]
export const addServer = (data) => {
  return service({
    url: '/cmdb/addServer',
    method: 'post',
    data
  })
}

// @Summary 删除server
// @Produce  application/json
// @Param ID float64
// @Router /cmdb/deleteServer [post]
export const deleteServer = (data) => {
  return service({
    url: '/cmdb/deleteServer',
    method: 'post',
    data
  })
}

// @Summary 修改server
// @Produce  application/json
// @Param server Object
// @Router /cmdb/updateServer [post]
export const updateServer = (data) => {
  return service({
    url: '/cmdb/updateServer',
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
// @Router /cmdb/getServerById [post]
export const getServerById = (data) => {
  return service({
    url: '/cmdb/getServerById',
    method: 'post',
    data
  })
}

// @Tags server
// @Summary 根据系统id获取服务器
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "系统id"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /cmdb/getSystemServers [post]
export const getSystemServers = (data) => {
  return service({
    url: '/cmdb/getSystemServers',
    method: 'post',
    data
  })
}

// @Tags server
// @Summary 获取所有服务器
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /cmdb/getAllServers [get]
export const getAllServerIds = () => {
  return service({
    url: '/cmdb/getAllServerIds',
    method: 'get',
  })
}

// @Tags Server
// @Summary 导出Excel
// @Security ApiKeyAuth
// @accept application/json
// @Produce  application/octet-stream
// @Param data body request2.ExcelInfo true "导出Excel文件信息"
// @Success 200
// @Router /cmdb/exportExcel [post]
export const exportExcel = (data) => {
  return service({
    url: '/cmdb/exportExcel',
    method: 'post',
    data
  })
}

// @Tags Server
// @Summary 导入Excel文件
// @Security ApiKeyAuth
// @accept multipart/form-data
// @Produce  application/json
// @Param file formData file true "导入Excel文件"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"导入成功"}"
// @Router /cmdb/importExcel [post]
export const importExcel = (data) => {
  return service({
    url: '/cmdb/importExcel',
    method: 'post',
    data
  })
}

// @Tags Server
// @Summary 下载模板
// @Security ApiKeyAuth
// @accept multipart/form-data
// @Produce  application/json
// @Success 200
// @Router /cmdb/downloadTemplate [get]
export const downloadTemplate = (fileName) => {
  return service({
    url: '/cmdb/downloadTemplate',
    method: 'get',
    responseType: 'blob'
  }).then((res) => {
    download(res, fileName)
  })
}

// @Summary 获取system列表
// @Produce  application/json
// @Param {
//  page     int
//	pageSize int
// }
// @Router /cmdb/getSystemList [post]
export const getSystemList = (data) => {
  return service({
    url: '/cmdb/getSystemList',
    method: 'post',
    data
  })
}

// @Summary 新增system
// @Produce  application/json
// @Param menu Object
// @Router /cmdb/addSystem [post]
export const addSystem = (data) => {
  return service({
    url: '/cmdb/addSystem',
    method: 'post',
    data
  })
}

// @Summary 删除system
// @Produce  application/json
// @Param ID float64
// @Router /cmdb/deleteSystem [post]
export const deleteSystem = (data) => {
  return service({
    url: '/cmdb/deleteSystem',
    method: 'post',
    data
  })
}

// @Summary 修改system
// @Produce  application/json
// @Param system Object
// @Router /cmdb/updateSystem [post]
export const updateSystem = (data) => {
  return service({
    url: '/cmdb/updateSystem',
    method: 'post',
    data
  })
}

// @Tags system
// @Summary 根据id获取服务器
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body api.GetById true "根据id获取服务器"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /cmdb/getSystemById [post]
export const getSystemById = (data) => {
  return service({
    url: '/cmdb/getSystemById',
    method: 'post',
    data
  })
}

// @Tags system
// @Summary 获取管理员服务器
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.Empty true
// @Success 200 {string} json "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /cmdb/getAdminSystems [post]
export const getAdminSystems = (data) => {
  return service({
    url: '/cmdb/getAdminSystems',
    method: 'post',
    data
  })
}

// @Summary 新增编辑器关系图
// @Produce  application/json
// @Param menu Object
// @Router /cmdb/system/addEditRelation [post]
export const addEditRelation = (data) => {
  return service({
    url: '/cmdb/system/addEditRelation',
    method: 'post',
    data
  })
}

// @Summary 删除编辑器关系图
// @Produce  application/json
// @Param ID float64
// @Router /cmdb/system/deleteEditRelation [post]
export const deleteEditRelation = (data) => {
  return service({
    url: '/cmdb/system/deleteEditRelation',
    method: 'post',
    data
  })
}

// @Summary 修改编辑器关系图
// @Produce  application/json
// @Param server Object
// @Router /cmdb/system/updateEditRelation [post]
export const updateEditRelation = (data) => {
  return service({
    url: '/cmdb/system/updateEditRelation',
    method: 'post',
    data
  })
}

// @Tags System
// @Summary 根据id获取编辑器关系图
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body api.GetById true "根据id获取编辑器关系图"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /cmdb/system/getSystemEditRelation [post]
export const getSystemEditRelation = (data) => {
  return service({
    url: '/cmdb/system/getSystemEditRelation',
    method: 'post',
    data
  })
}
