import service from '@/utils/request'

// @Summary 获取Token
// @Produce  application/json
// @Param menu Object
// @Router /jumpServer/getToken [post]
export const getToken = (data) => {
  return service({
    url: '/jumpServer/getToken',
    method: 'post',
    data
  })
}

// @Summary 请求server
// @Produce  application/json
// @Param menu Object
// @Router /jumpServer/getServer [post]
export const getServer = (data) => {
  return service({
    url: '/jumpServer/getServer',
    method: 'post',
    data
  })
}
