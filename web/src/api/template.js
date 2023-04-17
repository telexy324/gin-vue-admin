import service from '@/utils/request'

// @Summary 获取template列表
// @Produce  application/json
// @Param {
//  page     int
//	pageSize int
// }
// @Router /task/template/getTemplateList [post]
export const getTemplateList = (data) => {
  return service({
    url: '/task/template/getTemplateList',
    method: 'post',
    data
  })
}

// @Summary 新增template
// @Produce  application/json
// @Param template Object
// @Router /task/template/addTemplate [post]
export const addTemplate = (data) => {
  return service({
    url: '/task/template/addTemplate',
    method: 'post',
    data
  })
}

// @Summary 删除template
// @Produce  application/json
// @Param ID float64
// @Router /task/template/deleteTemplate [post]
export const deleteTemplate = (data) => {
  return service({
    url: '/task/template/deleteTemplate',
    method: 'post',
    data
  })
}

// @Summary 修改template
// @Produce  application/json
// @Param template Object
// @Router /task/template/updateTemplate [post]
export const updateTemplate = (data) => {
  return service({
    url: '/task/template/updateTemplate',
    method: 'post',
    data
  })
}

// @Summary 根据id获取template
// @Produce  application/json
// @Param ID float64
// @Router /task/template/getTemplateById [post]
export const getTemplateById = (data) => {
  return service({
    url: '/task/template/getTemplateById',
    method: 'post',
    data
  })
}

// @Summary 检查脚本
// @Produce  application/json
// @Param checkScript Object
// @Router /task/template/checkScript [post]
export const checkScript = (data) => {
  return service({
    url: '/task/template/checkScript',
    method: 'post',
    data
  })
}

// @Summary 下载脚本
// @Produce  application/json
// @Param checkScript Object
// @Router /task/template/downloadScript [post]
export const downloadScript = (data) => {
  return service({
    url: '/task/template/downloadScript',
    method: 'post',
    data
  })
}

// @Summary 上传脚本
// @Produce  application/json
// @Param checkScript Object
// @Router /task/template/downloadScript [post]
export const uploadScript = (data) => {
  return service({
    url: '/task/template/uploadScript',
    method: 'post',
    data
  })
}

// @Tags Template
// @Summary 删除选中Template
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /task/template/deleteTemplateByIds [post]
export const deleteTemplateByIds = (data) => {
  return service({
    url: '/task/template/deleteTemplateByIds',
    method: 'post',
    data
  })
}

// @Summary 获取set列表
// @Produce  application/json
// @Param {
//  page     int
//	pageSize int
// }
// @Router /task/template/getSetList [post]
export const getSetList = (data) => {
  return service({
    url: '/task/template/getSetList',
    method: 'post',
    data
  })
}

// @Summary 新增set
// @Produce  application/json
// @Param menu Object
// @Router /task/template/addSet [post]
export const addSet = (data) => {
  return service({
    url: '/task/template/addSet',
    method: 'post',
    data
  })
}

// @Summary 删除set
// @Produce  application/json
// @Param ID float64
// @Router /task/template/deleteSet [post]
export const deleteSet = (data) => {
  return service({
    url: '/task/template/deleteSet',
    method: 'post',
    data
  })
}

// @Summary 修改set
// @Produce  application/json
// @Param set Object
// @Router /task/template/updateSet [post]
export const updateSet = (data) => {
  return service({
    url: '/task/template/updateSet',
    method: 'post',
    data
  })
}

// @Tags set
// @Summary 根据id获取服务器
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body api.GetById true "根据id获取服务器"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /task/template/getSetById [post]
export const getSetById = (data) => {
  return service({
    url: '/task/template/getSetById',
    method: 'post',
    data
  })
}

// @Tags Set
// @Summary 删除选中Set
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.IdsReq true "ID"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /task/template/deleteSetByIds [post]
export const deleteSetByIds = (data) => {
  return service({
    url: '/task/template/deleteSetByIds',
    method: 'post',
    data
  })
}
