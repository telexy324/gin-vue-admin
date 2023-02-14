import service from '@/utils/request'

// @Tags Ssh
// @Summary 提交ssh信息
// @Security ApiKeyAuth
// @accept multipart/form-data
// @Produce  application/json
// @Success 200
// @Router /ssh/run [get]
export const runSsh = (data) => {
  return service({
    url: '/ssh/run',
    method: 'get',
    data
  })
}
