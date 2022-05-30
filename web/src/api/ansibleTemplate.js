import service from '@/utils/request'

// @Summary 获取template列表
// @Produce  application/json
// @Param {
//  page     int
//	pageSize int
// }
// @Router /ansible/template/getTemplateList [post]
export const getTemplateList = (data) => {
  return service({
    url: '/ansible/template/getTemplateList',
    method: 'post',
    data
  })
}

// @Summary 新增template
// @Produce  application/json
// @Param template Object
// @Router /ansible/template/addTemplate [post]
export const addTemplate = (data) => {
  return service({
    url: '/ansible/template/addTemplate',
    method: 'post',
    data
  })
}

// @Summary 删除template
// @Produce  application/json
// @Param ID float64
// @Router /ansible/template/deleteTemplate [post]
export const deleteTemplate = (data) => {
  return service({
    url: '/ansible/template/deleteTemplate',
    method: 'post',
    data
  })
}

// @Summary 修改template
// @Produce  application/json
// @Param template Object
// @Router /ansible/template/updateTemplate [post]
export const updateTemplate = (data) => {
  return service({
    url: '/ansible/template/updateTemplate',
    method: 'post',
    data
  })
}

// @Summary 根据id获取template
// @Produce  application/json
// @Param ID float64
// @Router /ansible/template/getTemplateById [post]
export const getTemplateById = (data) => {
  return service({
    url: '/ansible/template/getTemplateById',
    method: 'post',
    data
  })
}
