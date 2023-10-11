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
    <el-dialog v-model="CommandVarFormVisible" :before-close="closeCommandVarsDialog" :title="dialogTitle">
      <warning-bar title="请输入任务参数" />
      <el-form ref="CommandVarForm" :model="commandVarForm" :rules="commandVarRules" label-width="80px">
        <div v-for="(item, index) in commandVarForm.vars" :key="index">
          <el-form-item
            :label="'参数' + index"
            :prop="'vars.' + index"
            :rules="commandVarRules.vars"
          >
            <el-input v-model="commandVarForm.vars[index]" />
          </el-form-item>
        </div>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button size="small" @click="closeCommandVarsDialog">取 消</el-button>
          <el-button size="small" type="primary" @click="enterCommandVarsDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>
    <el-dialog v-model="confirmVisible" :before-close="closeConfirm" title="请确认">
      即将运行: {{ setTask.templates[setTask.currentStep].name }}
      <template #footer>
        <div class="dialog-footer">
          <el-button size="small" @click="closeConfirm">取 消</el-button>
          <el-button size="small" type="primary" @click="enterConfirm">确 定</el-button>
        </div>
      </template>
    </el-dialog>
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
import warningBar from '@/components/warningBar/warningBar.vue'
import { addTask } from '@/api/task'

export default {
  name: 'TemplateSetDetail',
  components: { warningBar },
  data() {
    return {
      setTaskId: '',
      path: path,
      steps: [],
      active: 0,
      setTask: '',
      // taskStatus: '',
      disabled: false,
      currentSystem: '',
      taskVars: [],
      CommandVarFormVisible: false,
      commandVarForm: {
        vars: [],
      },
      commandVarRules: {
        vars: [
          { required: true, message: '请输入任务参数', trigger: 'blur' },
        ],
      },
      runningTemplateId: '',
      confirmVisible: false,
      confirmed: false,
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
      if (this.setTask.templates[this.setTask.currentStep].commandVarNumbers > 0) {
        for (let i = 0; i < this.setTask.templates[this.setTask.currentStep].commandVarNumbers; i++) {
          this.commandVarForm.vars.push('')
        }
        this.runningTemplateId = this.setTask.templates[this.setTask.currentStep].ID
        this.CommandVarFormVisible = true
      } else {
        if (!this.confirmed) {
          this.confirmVisible = true
          return
        }
        this.confirmed = false
        const task = (await processSetTask({
          ID: this.setTask.ID,
        })).data.task
        this.showTaskLog(task)
      }
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
    },
    initCommandVarsForm() {
      this.commandVarForm = {
        vars: [],
      }
    },
    async enterCommandVarsDialog() {
      this.$refs.CommandVarForm.validate(async valid => {
        if (valid) {
          const task = (await processSetTask({
            ID: this.setTask.ID,
            commandVars: this.commandVarForm.vars
          })).data.task
          this.closeCommandVarsDialog()
          this.showTaskLog(task)
        }
      })
    },
    closeCommandVarsDialog() {
      this.CommandVarFormVisible = false
      this.initCommandVarsForm()
      this.runningTemplateId = ''
    },
    enterConfirm() {
      this.confirmVisible = false
      this.confirmed = true
      this.processSetTask()
    },
    closeConfirm() {
      this.confirmVisible = false
      this.pendingTemplate = ''
    },
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
