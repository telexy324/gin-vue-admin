import service from '@/utils/request'

// @Summary 获取server列表
// @Produce  application/json
// @Param {
//  id     int
// }
// @Router /cmdb/getServerList [post]
// 查询企业信息
export const getServerRelations = (data) => {
  return service({
    url: '/cmdb/server/relations',
    method: 'post',
    data
  })
}
