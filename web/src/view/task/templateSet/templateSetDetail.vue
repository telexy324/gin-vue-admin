<template>
  <div>
    <div class="gva-search-box">
      <div class="gva-btn-list">
        <el-button size="mini" type="primary" icon="el-icon-plus" style="margin-bottom: 12px;" :disabled="disabled" @click="processSetTask">下一步</el-button>
        <el-button v-if="forceCorrectButton" size="mini" type="danger" icon="el-icon-plus" style="margin-bottom: 12px;" @click="forceCorrect">强制执行</el-button>
      </div>
<!--      <el-steps :active="active" finish-status="success" :process-status="taskStatus">-->
      <el-steps>
        <el-step v-for="item in steps" :key="item.seq" :title="item.name" :status="getStatus(item.seq)" @click.enter="show(item.seq)"/>
      </el-steps>
    </div>
    <el-dialog v-model="CommandVarFormVisible" :before-close="closeCommandVarsDialog" title="参数">
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
        <el-form-item label="目标" prop="targetIds">
          <el-cascader
            v-model="commandVarForm.targetIds"
            style="width:100%"
            :options="checkedServerOptions"
            :show-all-levels="false"
            :props="{ multiple:true,checkStrictly: false,label:'name',value:'ID',disabled:'disabled',emitPath:false}"
            :clearable="true"
          />
        </el-form-item>
        <el-row v-if="netDisk">
          <el-col :span="12">
            <el-form-item label="用户名" prop="netDiskUser">
              <el-input v-model="commandVarForm.netDiskUser" autocomplete="off" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="密码" prop="netDiskPassword">
              <el-input v-model="commandVarForm.netDiskPassword" autocomplete="off" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button size="small" @click="closeCommandVarsDialog">取 消</el-button>
          <el-button size="small" type="primary" @click="enterCommandVarsDialog">确 定</el-button>
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
  getSetTaskById,
  setTaskForceCorrect
} from '@/api/template'
import { getSystemServerIds } from '@/api/cmdb'
import { emitter } from '@/utils/bus'
import warningBar from '@/components/warningBar/warningBar.vue'

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
        targetIds: [],
      },
      commandVarRules: {
        vars: [{ required: true, message: '请输入任务参数', trigger: 'blur' }],
        targetIds: [{ required: true, message: '请选择目标', trigger: 'blur' }],
      },
      runningTemplateId: '',
      confirmVisible: false,
      confirmed: false,
      checkedServerOptions: [],
      netDisk: false,
      forceCorrectButton: false,
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
      this.isForceCorrectButton()
      if (this.setTask.currentStep === this.setTask.totalSteps || this.setTask.tasks[this.active - 1].status !== 'success' && this.setTask.tasks[this.active - 1].status !== '' && this.setTask.forceCorrect === 0) {
        this.disabled = true
      }
    },
    async processSetTask() {
      await this.setCheckedServerOptions(this.setTask.templates[this.setTask.currentStep])
      for (let i = 0; i < this.setTask.templates[this.setTask.currentStep].commandVarNumbers; i++) {
        this.commandVarForm.vars.push('')
      }
      this.runningTemplateId = this.setTask.templates[this.setTask.currentStep].ID
      this.CommandVarFormVisible = true
      this.setTask.templates[this.setTask.currentStep].deployType === 2 ? this.netDisk = true : this.netDisk = false
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
      if (seq + 1 <= this.active) {
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
        targetIds: [],
      }
    },
    async enterCommandVarsDialog() {
      this.$refs.CommandVarForm.validate(async valid => {
        if (valid) {
          const task = (await processSetTask({
            ID: this.setTask.ID,
            commandVars: this.commandVarForm.vars,
            targetIds: this.commandVarForm.targetIds,
            netDiskUser: this.commandVarForm.netDiskUser,
            netDiskPassword: this.commandVarForm.netDiskPassword,
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
    async setCheckedServerOptions(template) {
      const res = await getSystemServerIds({
        ID: template.systemId
      })
      const serverOptions = res.data
      if (template.executeType !== 2) {
        serverOptions[0].children = serverOptions[0].children.filter((item) => {
          if (template.targetServerIds.includes(item.ID)) {
            this.commandVarForm.targetIds.push(item.ID)
            return true
          }
          return false
        })
      }
      this.checkedServerOptions = serverOptions
    },
    async forceCorrect() {
      await setTaskForceCorrect({
        ID: this.setTask.ID,
      })
      this.forceCorrectButton = false
      this.disabled = false
      await this.initSteps()
    },
    isForceCorrectButton() {
      if (this.setTask.tasks[this.active - 1].status !== 'success' && this.setTask.forceCorrect === 0) {
        this.forceCorrectButton = true
      }
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
