<template>
  <div>
    <div class="gva-search-box">
      <div class="gva-btn-list">
        <el-button size="mini" type="primary" icon="el-icon-plus" style="margin-bottom: 12px;" @click="runTask">构建</el-button>
      </div>
      <el-steps :active="active" finish-status="success" process-status="error">
        <el-step v-for="item in steps" :key="item.ID" :title="item.templateName" />
      </el-steps>
    </div>
  </div>
</template>

<script>

const path = import.meta.env.VITE_BASE_API
// 获取列表内容封装在mixins内部  getTableData方法 初始化已封装完成 条件搜索时候 请把条件安好后台定制的结构体字段 放到 this.searchInfo 中即可实现条件搜索

import {
  getSetById,
  addSetTask,
  processSetTask,
  getSetTaskById
} from '@/api/template'
import { emitter } from "@/utils/bus";
// import { emitter } from '@/utils/bus'

export default {
  name: 'TemplateSetDetail',
  data() {
    return {
      setId: '',
      path: path,
      steps: [],
      active: 1,
    }
  },
  async created() {
    this.setId = this.$route.params.setId
    await this.initSteps()
  },
  methods: {
    async initSteps() {
      this.steps = (await getSetById({ 'ID': Number(this.setId) })).data.templates
    },
    async runTask() {
      const task = (await addTask({
        templateId: row.ID
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
