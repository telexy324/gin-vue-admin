<template>
  <div>
    <div class="gva-search-box">
      <div class="gva-btn-list">
        <el-button size="mini" type="primary" icon="el-icon-plus" style="margin-bottom: 12px;" :disabled="disabled" @click="processSetTask">下一步</el-button>
      </div>
<!--      <el-steps :active="active" finish-status="success" :process-status="taskStatus">-->
      <el-steps>
        <el-step v-for="item in steps" :key="item.seq" :title="item.name" :status="getStatus(item.seq)" @click.enter="show(item.seq)"/>
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
      active: 0,
      setTask: '',
      // taskStatus: '',
      disabled: false,
    }
  },
  async created() {
    this.setTaskId = this.$route.params.setTaskId
    // await this.initSteps()
    await this.$nextTick(() => {
      this.initSteps()
    })
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
      this.active = this.setTask.currentStep
      // this.taskStatus = this.getStepStatus(this.setTask.currentTask.status)
      if (this.setTask.currentStep === this.setTask.totalSteps || this.setTask.tasks[this.active - 1].status !== 'success' && this.setTask.tasks[this.active - 1].status !== '') {
        this.disabled = true
      }
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
    // getStepStatus(status) {
    //   switch (status) {
    //     case 'success':
    //       return 'success'
    //     default:
    //       return 'error'
    //   }
    // },
    getStatus(seq) {
      if (seq + 1 < this.active) {
        return 'success'
      } else if (seq + 1 === this.active) {
        const status = this.setTask.tasks[seq].status
        switch (status) {
          case 'success':
            return 'success'
          default:
            return 'error'
        }
      } else {
        return 'wait'
      }
    },
    show(seq) {
      this.showTaskLog(this.setTask.tasks[seq])
    }
  }
}
</script>

<style lang="scss">
.el-step:hover {
  cursor: pointer;
}
</style>

<style scoped lang="scss">
.button-box {
  padding: 10px 20px;
  .el-button {
    float: right;
  }
}
</style>
