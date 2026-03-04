import { toast } from 'sonner'

type DownloadLike = {
  data?: Blob | { type?: string } | any
}

const download = (res: DownloadLike, fileName: string) => {
  if (typeof res.data !== 'undefined') {
    if (res.data?.type === 'application/json') {
      const reader = new FileReader()
      reader.onload = function () {
        try {
          const message = JSON.parse(String(reader.result || '{}')).msg
          toast.error(message || '下载失败')
        } catch (_) {
          toast.error('下载失败')
        }
      }
      reader.readAsText(new Blob([res.data]))
      return
    }
  }

  const downloadUrl = window.URL.createObjectURL(new Blob([res as any]))
  const a = document.createElement('a')
  a.style.display = 'none'
  a.href = downloadUrl
  a.download = fileName
  a.dispatchEvent(new MouseEvent('click'))
}

export default download
