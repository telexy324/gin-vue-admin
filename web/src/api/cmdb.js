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
export const getAllServers = () => {
  return service({
    url: '/cmdb/getAllServers',
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
