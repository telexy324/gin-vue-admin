<template lang="pug">
  .relation.common
    .box(v-if='!err')
      //- 移动端操作工具栏
      .toolbar-mobile(v-if="isMobile" :class="isRotate ? 'rotate-y' : 'rotate-x'" :style="isRotate ? toolbarStyle : ''")
        ul
          li(@click="onExitPage")
            img(src="@common/assets/image/exit.png")
            span Exit
          li
          li(@click="onReset")
            img(src="@common/assets/image/reset.png")
            span Reset
          li
          li(@click="onRotate")
            img(v-if="isRotate" src="@common/assets/image/vertical.png")
            img(v-else src="@common/assets/image/horizontal.png")
            span {{ isRotate ? 'Vertical' : 'Horizontal' }}
      #chart(ref='chart' :style='{ height: `${screenHeight-10}px` }')
    NoData(v-if='err && err.message' :title='err.message')
    Loading(:loading="loading")
</template>

<script>
import API from '@api'
import echartMixins from '@common/mixins/echartMixins'
import NoData from '@components/NoData'
import Loading from '@components/loading'
import { isMobile } from '@common/utils/util'

export default {
  name: 'relation',
  mixins: [echartMixins('#chart')],
  data() {
    return {
      isMobile: isMobile(), // 设备类型判断
      toolbarStyle: '', // 工具栏样式
      screenHeight: 500,
      err: null,
      loading: false,
      cid: null,
      cname: null,
      // echarts 实例
      myChart: null,
    }
  },
  computed: {},
  created() {
    this.screenHeight = window.innerHeight || document.documentElement.clientHeight || document.body.clientHeight
  },
  mounted() {
    const { companyId, companyName } = this.$route.query
    this.cid = companyId
    this.cname = decodeURIComponent(companyName || '')
    this.getData()
  },
  watch: {},
  methods: {
    getData() {
      const body = {
        // cid: '819430495192022480',
        cid: this.cid,
      }
      this.loading = true
      API.SEARCH_COMPANY_RELATION(body)
          .then(this.handleData)
          .catch(err => {
            this.err = err
          })
          .finally(() => (this.loading = false))
    },
    // 处理数据为echarts需要的格式
    handleData(res) {
      console.time('数据处理耗时')
      const {
        data: { path },
      } = res
      if (path.nodes.length == 0) this.err = { message: 'No data' }
      // 目录
      let categories = []
      // 缓存关系数据，有重合路线的曲线展示
      // let linkCache = {}
      // 关系数组处理
      path.links.forEach(link => {
        link.value = link.vector_str_value
        link.source = link.start_node_id
        link.target = link.end_node_id
        // if (linkCache[link.target + link.source]) {
        //     link.lineStyle = {
        //         curveness: 0.1,
        //     }
        // } else {
        //     linkCache[link.source + link.target] = 1
        // }
        if (link.vector_type == 'gv_shareholder') {
          link.value = link.property.shareholding_name || link.vector_str_value
        } else {
          link.lineStyle = {
            curveness: 0.1,
          }
          link.value = link.property.position || link.vector_str_value
        }
      })
      // 节点数组处理
      path.nodes.forEach((node, index) => {
        node.label = {
          show: true,
        }
        const len = categories.length
        const idx = categories.indexOf(node.name)
        if (idx === -1) {
          categories[len] = node.name
        } else {
          node.category = idx
        }
        for (let i = index + 1; i < path.nodes.length; i++) {
          const element = path.nodes[i]
          // 防止name重复，echarts 需要name唯一
          if (element.name == node.name) node.name = node.name + ++index
        }
      })

      categories = categories.map(name => ({ name }))

      this.path = path
      this.categories = categories
      console.timeEnd('数据处理耗时')
      this.initChart()
    },
    // 初始化
    init(reset) {
      // 创建chart
      this.initChart(reset)
    },
    initChart(reset) {
      const chartData = this.path
      const { myChart, categories } = this
      if (reset) {
        this.myChart.dispose()
        this.myChart = null
      }
      if (this.myChart) {
        myChart.resize()
        return
      }
      this.myChart = echarts.init(this.$refs.chart)
      const cid = this.cid
      // console.log('接口数据data', JSON.stringify(chartData.nodes[0], null, 2), chartData.nodes)
      // console.log('接口数据links', JSON.stringify(chartData.links[0], null, 2), chartData.links)
      // console.log('接口数据分类数据categories', JSON.stringify(categories, null, 2))
      const option = {
        tooltip: {},
        animationDuration: 1500,
        animationEasingUpdate: 'quinticInOut',
        series: [
          {
            type: 'graph',
            layout: 'force',
            // 节点大小
            symbolSize(value, { data }) {
              // 根节点大小
              if (data.id == cid) return 88
              // 其它节点大小
              return 70
            },
            force: {
              // 节点之间的斥力因子。值越大则斥力越大
              repulsion: 400,
              // 边的两个节点之间的距离，这个距离也会受 repulsion 影响。值越小则长度越长
              edgeLength: [100, 500],
              // 节点受到的向中心的引力因子。该值越大节点越往中心点靠拢。
              gravity: 0.1,
              // 这个参数能减缓节点的移动速度。取值范围 0 到 1。
              friction: 0.5,
              layoutAnimation: true,
            },
            label: {
              fontSize: 10,
              width: '100%',
              formatter: this.handleNode,
            },
            tooltip: {
              textStyle: {
                fontSize: 10,
              },
              formatter({ data }) {
                let name = data.name
                // 点击节点
                if (data.id) {
                  // name过长时，换行
                  const len = name.length
                  if (len > 18) name = `<div style='max-width:160px;white-space: normal;word-break:break-word'>${name}</div>`
                  return name
                }
                // 点击线
                if (data.property && data.property.shareholding_ratio) return `${data.value}: ${data.property.shareholding_ratio}%`
                return data.value
              },
            },
            edgeLabel: {
              show: true,
              fontSize: 8,
              formatter: '{c}',
              // formatter(i) {
              //     return i.data.value
              // },
            },
            itemStyle: {
              color({ data }) {
                // 主节点
                if (data.id == cid) return '#288bff'
                // 公司
                if (data.type == 'gn_company') return '#3ea3ff'
                // 人物
                return '#19cc9d'
              },
            },
            edgeSymbol: ['', 'arrow'],
            roam: true,
            // draggable: true,
            nodeScaleRatio: 0.6,
            data: chartData.nodes,
            links: chartData.links,
            // 鼠标移到节点上的时候突出显示节点以及节点的边和邻接节点。
            focusNodeAdjacency: true,
            lineStyle: {
              color: 'source',
            },
            categories,
          },
        ],
      }
      this.myChart.setOption(option)
    },
    // 处理节点名字换行
    handleNode({ data }) {
      // 最多显示3行
      let row = 3,
          // 每行的字符数量
          num = 9,
          name = data.name,
          newName = '',
          index = 0

      while (name.length > num * index && index < row) newName += name.substring(num * index, num * ++index) + '\n'
      newName = newName.replace(/\n$/g, '')
      if (name.length > num * row) newName += '...'
      return newName
    },
  },
  components: { NoData, Loading },
}
</script>

<style scoped lang="stylus">
.relation {
  #chart {
    position fixed!important
    left 0
    top 0
    width 100%
    height 100%
  }
}
</style>
