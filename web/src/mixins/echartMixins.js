/**
 * echart基础功能
 */
import { commonCreateModule } from '@common/d3/common'

export default element => {
  return {
    data() {
      return {
        isRotate: false, // 是否旋转
        screenWidth: 0, // 屏幕宽
        screenHeight: 0, // 屏幕高
        // 基础容器组件
        SVG: null, // svg实例
      }
    },
    created() {
    },
    mounted() {
      // 初始化创建svg画布 + 创建水印
      this.createWatemark()
    },
    methods: {
      // 创建svg画布 + 创建水印
      createWatemark() {
        // 删除原svg
        d3.select('svg').remove()
        // 获取屏幕宽高
        this.screenWidth = window.innerWidth || document.documentElement.clientWidth || document.body.clientWidth
        this.screenHeight = window.innerHeight || document.documentElement.clientHeight || document.body.clientHeight
        // 创建svg画布
        this.SVG = d3
          .select('.common') // d3.select() 是选择所有指定元素的第一个
          .append('svg') // 添加画布
          .attr('width', this.screenWidth) // 设置宽度
          .attr('height', this.screenHeight) // 设置高度
        // 创建水印组件
        commonCreateModule.createWatemark(this)
      },
      // 退出页面
      onExitPage() {
        const exit = {
          functionName: 'Exit',
        }
        this.$onAppPartyContext(exit)
      },
      // 放大
      onMagnify() {
        this.option.series[0].zoom += 0.1
        // 渲染
        this.myChart.setOption(this.option)
      },
      // 缩小
      onShrink() {
        this.option.series[0].zoom -= 0.1
        // 渲染
        this.myChart.setOption(this.option)
      },
      // 重置
      onReset() {
        this.init('reset')
      },
      // 旋转页面
      onRotate() {
        this.isRotate = !this.isRotate
        // 计算工具栏旋转样式
        if (this.isRotate) {
          setTimeout(() => {
            const padding = 20
            const toolWidth = document.querySelector('.toolbar-mobile').clientWidth,
              toolHeight = document.querySelector('.toolbar-mobile').clientHeight
            this.toolbarStyle = `
                        right: ${-toolWidth + toolHeight + padding}px;
                        bottom: ${(this.screenHeight - toolWidth) / 2 - toolHeight}px
                    `
          }, 0)
        }
        let chartEle = element || '#myChart'
        const myChart = document.querySelector(chartEle)
        const realWidth = this.isRotate ? this.screenHeight : this.screenWidth
        const realHeight = this.isRotate ? this.screenWidth : this.screenHeight
        this.isRotate
          ? (myChart.style = `transform: rotate(-90deg); transform-origin: ${realWidth / 2}px ${realWidth /
          2}px;width: ${realWidth}px; height: ${realHeight}px;`)
          : (myChart.style = `transform: rotate(0deg);width: ${realWidth}px; height: ${realHeight}px;`)
        // 初始化
        this.createWatemark()
        this.init()
        this.onFlutterRotate()
      },
      // 通知原生进行旋转
      onFlutterRotate() {
        const exit = {
          functionName: 'FullScreen',
          postMessage: this.isRotate ? '1' : '0',
        }
        this.$onAppPartyContext(exit)
      },
    },
  }
}
