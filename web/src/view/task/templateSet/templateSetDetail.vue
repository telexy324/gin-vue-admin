<template>
  <div>
    <div class="gva-search-box">
      <div class="gva-btn-list">
        <el-button size="mini" type="primary" icon="el-icon-plus" style="margin-bottom: 12px;" @click="processSetTask">下一步</el-button>
      </div>
      <el-steps :active="active" finish-status="success" process-status="error">
        <el-step v-for="item in steps" :key="item.ID" :title="item.name" />
      </el-steps>
    </div>
  </div>
</template>

<script>

const path = import.meta.env.VITE_BASE_API
// 获取列表内容封装在mixins内部  getTableData方法 初始化已封装完成 条件搜索时候 请把条件安好后台定制的结构体字段 放到 this.searchInfo 中即可实现条件搜索

import {
  // getSetById,
  // addSetTask,
  processSetTask,
  getSetTaskById
} from '@/api/template'
import { emitter } from '@/utils/bus'

export default {
  name: 'TemplateSetDetail',
  data() {
    return {
      setTaskId: '',
      path: path,
      steps: [],
      active: 1,
      setTask: '',
    }
  },
  async created() {
    this.setTaskId = this.$route.params.setTaskId
    console.log(this.setTaskId)
    await this.initSteps()
  },
  mounted() {
    emitter.on('i-close-task', () => {
      this.initSteps()
    })
  },
  methods: {
    async initSteps() {
      this.setTask = (await getSetTaskById({ 'ID': Number(this.setTaskId) })).data
      this.steps = this.setTask.templates
    },
    async processSetTask() {
      const task = (await processSetTask({
        ID: this.setTask.ID,
      })).data.task
      console.log(task.ID)
      this.showTaskLog(task)
    },
    showTaskLog(task) {
      emitter.emit('i-show-task', task)
    },
  }
}
</script>

<style scoped lang="scss">
.button-box {
  padding: 10px 20px;
  .el-button {
    float: right;
  }
}
</style>
