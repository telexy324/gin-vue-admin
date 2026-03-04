export const downloadImage = (imgsrc: string, name?: string) => {
  const image = new Image()
  image.setAttribute('crossOrigin', 'anonymous')
  image.onload = function () {
    const canvas = document.createElement('canvas')
    canvas.width = image.width
    canvas.height = image.height
    const context = canvas.getContext('2d')
    context?.drawImage(image, 0, 0, image.width, image.height)
    const url = canvas.toDataURL('image/png')

    const a = document.createElement('a')
    const event = new MouseEvent('click')
    a.download = name || 'photo'
    a.href = url
    a.dispatchEvent(event)
  }
  image.src = imgsrc
}

