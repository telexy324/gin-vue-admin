// d3公共创建组件
const commonCreateModule = {
  // 创建水印组件
  createWatemark(that) {
    // 添加水印
    const width = 92 * 1.4
    const height = 70 * 1.4
    const imageG = that.SVG.append('defs')
      .append('pattern')
      .attr('id', 'watermark')
      .attr('fill', 'url(#watermark)')
      .attr('width', width)
      .attr('height', height)
      .attr('x', -40)
      .attr('y', -40)
      .attr('patternUnits', 'userSpaceOnUse')
      .append('image')
    imageG
      .attr('width', width)
      .attr('height', height)
      .attr('x', 0)
      .attr('y', 0)
      .attr('href', require('@common/assets/image/background.png'))
      .attr('src', require('@common/assets/image/background.png'))
      .attr('class', that.isRotate ? 'imageRotate' : '')
    that.SVG.append('rect')
      .attr('width', that.screenWidth)
      .attr('height', that.screenHeight)
      .attr('fill', 'url(#watermark)')
  },
  // 创建阴影Defs组件
  createShadowDefs(that, config = {}) {
    const id = 'md-shadow',
      deviation = config.deviation || 1,
      offsetX = config.offsetX || 1,
      offsetY = config.offsetY || 3,
      slope = config.slope || 0.15

    const defs = that.SVG.append('defs')

    const filter = defs.append('filter').attr('id', id)

    filter
      .append('feGaussianBlur')
      .attr('in', 'SourceAlpha')
      .attr('stdDeviation', deviation)
      .attr('result', 'ambientBlur')

    filter
      .append('feGaussianBlur')
      .attr('in', 'SourceAlpha')
      .attr('stdDeviation', deviation)
      .attr('result', 'castBlur')

    filter
      .append('feOffset')
      .attr('in', 'castBlur')
      .attr('dx', offsetX)
      .attr('dy', offsetY)
      .attr('result', 'offsetBlur')

    filter
      .append('feComposite')
      .attr('in', 'ambientBlur')
      .attr('in2', 'offsetBlur')
      .attr('result', 'compositeShadow')

    filter
      .append('feComponentTransfer')
      .append('feFuncA')
      .attr('type', 'linear')
      .attr('slope', slope)

    const feMerge = filter.append('feMerge')
    feMerge.append('feMergeNode')
    feMerge.append('feMergeNode').attr('in', 'SourceGraphic')
  },
}

// d3公共事件
const commonIncident = {
  // 退出页面
  onExitPage(that) {
    const exit = {
      functionName: 'Exit',
    }
    that.$onAppPartyContext(exit)
  },
  // 旋转页面
  onRotateSVG(that) {
    that.isRotate = !that.isRotate
    // 计算工具栏旋转样式
    if (that.isRotate) {
      setTimeout(() => {
        const padding = 20
        const toolWidth = document.querySelector('.toolbar-mobile').clientWidth,
          toolHeight = document.querySelector('.toolbar-mobile').clientHeight
        that.toolbarStyle = `
                    right: ${-toolWidth + toolHeight + padding}px;
                    bottom: ${(that.screenHeight - toolWidth) / 2 - toolHeight}px
                `
      }, 0)
    }
    that.onResetSVG()
  },
  // 全屏展示
  onFullScreen(that) {
    let element = document.documentElement
    if (that.isScreenState) {
      if (document.exitFullscreen) {
        document.exitFullscreen()
      } else if (document.msExitFullscreen) {
        document.msExitFullscreen()
      } else if (document.mozCancelFullScreen) {
        document.mozCancelFullScreen()
      } else if (document.webkitExitFullscreen) {
        document.webkitExitFullscreen()
      }
    } else {
      if (element.requestFullscreen) {
        element.requestFullscreen()
      } else if (element.msRequestFullscreen) {
        element.msRequestFullscreen()
      } else if (element.mozRequestFullScreen) {
        element.mozRequestFullScreen()
      } else if (element.webkitRequestFullscreen) {
        element.webkitRequestFullscreen()
      }
    }
    setTimeout(that.onResetSVG, 200)
  },
  // 放大SVG
  onMagnifySVG(that) {
    that.zoomHandler.scaleBy(that.SVG, 1.1) // 执行该方法后 会触发zoom事件
    d3.zoomTransform(that.SVG.node())
  },
  // 缩小SVG
  onShrinkSVG(that) {
    that.zoomHandler.scaleBy(that.SVG, 0.9) // 执行该方法后 会触发zoom事件
    d3.zoomTransform(that.SVG.node())
  },
  // 下载图片
  onDownload(that, name) {
    // svg => canvas => png
    const serializer = new XMLSerializer()
    const source = '<?xml version="1.0" standalone="no"?>\r\n' + serializer.serializeToString(that.SVG.node())
    const image = new Image()
    image.src = 'data:image/svg+xml;charset=utf-8,' + encodeURIComponent(source)
    const canvas = document.createElement('canvas')
    canvas.width = that.screenWidth
    canvas.height = that.screenHeight
    const context = canvas.getContext('2d')
    context.fillStyle = '#fff'
    context.fillRect(0, 0, 10000, 10000)
    image.onload = () => {
      context.drawImage(image, 0, 0)
      const a = document.createElement('a')
      a.download = `${name || 'download'}.png`
      a.href = canvas.toDataURL('image/png')
      a.click()
    }
  },
}

export { commonCreateModule, commonIncident }
