import service from '@/utils/request'
import download from '@/utils/download'

// @Tags excel
// @Summary 导出Excel
export const exportExcel = (tableData: any[], fileName: string) => {
  return service({
    url: '/excel/exportExcel',
    method: 'post',
    data: {
      fileName,
      infoList: tableData
    },
    responseType: 'blob'
  }).then((res) => {
    download(res as any, fileName)
  })
}

// @Tags excel
// @Summary 导入Excel文件
export const loadExcelData = () => {
  return service({
    url: '/excel/loadExcel',
    method: 'get'
  })
}

// @Tags excel
// @Summary 下载模板
export const downloadTemplate = (fileName: string) => {
  return service({
    url: '/excel/downloadTemplate',
    method: 'get',
    params: { fileName },
    responseType: 'blob'
  }).then((res) => {
    download(res as any, fileName)
  })
}
