<template>
  <div class="dashbord-line-box">
    <div class="dashbord-line-title">
      执行任务量
    </div>
    <div
      ref="echart"
      class="dashbord-line"
    />
  </div>
</template>
<script>
import echarts from 'echarts'
import 'echarts/theme/macarons'
import { getTaskDashboardInfo } from '@/api/task'

export default {
  name: 'Line2',
  data() {
    return {
      chart: null,
      dataAxis: [],
      data: [],
      yMax: 500,
      dataShadow: [],
    }
  },
  async created() {
    const res = (await getTaskDashboardInfo()).data.taskDashboardInfos
    this.fillData(res)
    this.initChart()
  },
  // mounted() {
  //   this.$nextTick(() => {
  //     this.initChart()
  //   })
  // },
  beforeUnmount() {
    if (!this.chart) {
      return
    }
    this.chart.dispose()
    this.chart = null
  },
  methods: {
    initChart() {
      this.chart = echarts.init(this.$refs.echart, 'macarons')
      this.setOptions()
    },
    setOptions() {
      this.chart.setOption({
        grid: {
          left: '40',
          right: '20',
          top: '40',
          bottom: '20',
        },
        xAxis: {
          data: this.dataAxis,
          axisTick: {
            show: false,
          },
          axisLine: {
            show: false,
          },
          z: 10,
        },
        yAxis: {
          axisLine: {
            show: false,
          },
          axisTick: {
            show: false,
          },
          axisLabel: {
            textStyle: {
              color: '#999',
            },
          },
        },
        dataZoom: [
          {
            type: 'inside',
          },
        ],
        series: [
          {
            type: 'bar',
            barWidth: '40%',
            itemStyle: {
              borderRadius: [5, 5, 0, 0],
              color: '#188df0',
            },
            emphasis: {
              itemStyle: {
                color: '#188df0',
              },
            },
            data: this.data,
          },
        ],
      })
    },
    fillData(res) {
      if (res.length) {
        res.map(item => {
          this.data.push(item.number)
          this.dataAxis.push(item.date)
          this.dataShadow.push(this.yMax)
        })
      }
    },
  },
}
</script>
<style lang="scss" scoped>
.dashbord-line-box {
  .dashbord-line {
    background-color: #fff;
    height: 360px;
    width: 100%;
  }
  .dashbord-line-title {
    font-weight: 600;
    margin-bottom: 12px;
  }
}
</style>
