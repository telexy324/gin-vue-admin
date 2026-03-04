export default class ImageCompress {
  file: File
  fileSize: number
  maxWH: number

  constructor(file: File, fileSize: number, maxWH = 1920) {
    this.file = file
    this.fileSize = fileSize
    this.maxWH = maxWH
  }

  compress() {
    const fileType = this.file.type
    return new Promise<File>((resolve) => {
      const reader = new FileReader()
      reader.readAsDataURL(this.file)
      reader.onload = () => {
        const canvas = document.createElement('canvas')
        const img = document.createElement('img')
        img.src = String(reader.result || '')
        img.onload = () => {
          const ctx = canvas.getContext('2d')
          const wh = this.dWH(img.width, img.height, this.maxWH)
          canvas.width = wh.width
          canvas.height = wh.height
          ctx?.clearRect(0, 0, canvas.width, canvas.height)
          ctx?.drawImage(img, 0, 0, canvas.width, canvas.height)
          const newImgData = canvas.toDataURL(fileType, 0.9)
          const blob = this.dataURLtoBlob(newImgData, fileType)
          resolve(new File([blob], this.file.name))
        }
      }
    })
  }

  dWH(srcW: number, srcH: number, dMax: number) {
    const defaults = { width: srcW, height: srcH }
    if (Math.max(srcW, srcH) > dMax) {
      if (srcW > srcH) {
        defaults.width = dMax
        defaults.height = Math.round(srcH * (dMax / srcW))
      } else {
        defaults.height = dMax
        defaults.width = Math.round(srcW * (dMax / srcH))
      }
    }
    return defaults
  }

  dataURLtoBlob(dataURL: string, fileType?: string) {
    const byteString = atob(dataURL.split(',')[1])
    let mimeString = dataURL.split(',')[0].split(':')[1].split(';')[0]
    const ab = new ArrayBuffer(byteString.length)
    const ia = new Uint8Array(ab)
    for (let i = 0; i < byteString.length; i++) {
      ia[i] = byteString.charCodeAt(i)
    }
    if (fileType) mimeString = fileType
    return new Blob([ab], { type: mimeString })
  }
}

