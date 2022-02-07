import service from '@/utils/request'

// @Summary 获取server列表
// @Produce  application/json
// @Param {
//  page     int
//	pageSize int
// }
// @Router /cmdb/getServerList [post]
// 查询企业信息
export const getSystemRelations = (data) => {
  return service({
    url: '/cmdb/system/relations',
    method: 'post',
    data
  })
}
