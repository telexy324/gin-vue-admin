<template>
  <div>
    <div class="gva-search-box">
      <div class="gva-btn-list">
        <el-button size="mini" type="primary" icon="el-icon-plus" style="margin-bottom: 12px;" :disabled="disabled" @click="enterVars">下一步</el-button>
        <el-button v-if="forceCorrectButton" size="mini" type="danger" icon="el-icon-plus" style="margin-bottom: 12px;" @click="forceCorrect">强制执行</el-button>
      </div>
<!--      <el-steps :active="active" finish-status="success" :process-status="taskStatus">-->
      <el-steps>
        <el-step v-for="(item, index) in steps" :key="index" :title="'步骤 ' + index" :status="getStatus(index)" @click.enter="showTasks(index)"/>
      </el-steps>
    </div>
    <el-dialog v-model="VarListVisible" :before-close="closeVarsDialog" title="参数列表">
<!--      <ul class="file-name">-->
<!--        <li-->
<!--          v-for="item in setTask.templates[setTask.currentStep]"-->
<!--          :key="item.innerSeq"-->
<!--          class="file"-->
<!--          @click="checkVars(item.innerSeq)"-->
<!--        >-->
<!--          <pre>{{ item.innerSeq }} {{ item.name }}</pre>-->
<!--        </li>-->
<!--      </ul>-->
<!--      <template #footer>-->
<!--        <div class="dialog-footer">-->
<!--          <el-button size="small" @click="closeVarsDialog">取 消</el-button>-->
<!--          <el-button size="small" type="primary" @click="enterVarsDialog">确 定</el-button>-->
<!--        </div>-->
<!--      </template>-->
      <el-table :data="tableData" @sort-change="sortChange" @selection-change="handleSelectionChange">
        <el-table-column align="left" label="id" min-width="60" prop="id" sortable="custom">
          <template v-slot="scope">
            <el-button
              type="text"
              link
              @click="showTaskLog(scope.row)"
            >#{{ scope.row.ID }}</el-button>
            <!--            <a @click="showTaskLog(scope.row)">{{ scope.row.ID }}</a>-->
          </template>
        </el-table-column>
        <el-table-column align="left" label="模板名" min-width="150" prop="templateId" sortable="custom">
          <template #default="scope">
            <div>{{ filterTemplateName(scope.row.templateId) }}</div>
          </template>
        </el-table-column>
        <el-table-column align="left" label="状态" min-width="150">
          <template v-slot="scope">
            <TaskStatus :status="scope.row.status" />
          </template>
        </el-table-column>
        <el-table-column align="left" label="开始时间" min-width="150" prop="beginTime.Time" sortable="custom" :formatter="dateFormatter1" />
        <el-table-column align="left" label="结束时间" min-width="150" prop="endTime.Time" sortable="custom" :formatter="dateFormatter2" />
        <el-table-column align="right" label="参数" min-width="60" prop="id" sortable="custom">
          <template v-slot="scope">
            <el-button
              type="text"
              link
              @click="checkVars(scope.row.setTaskInnerSeq)"
            >参数</el-button>
            <!--            <a @click="showTaskLog(scope.row)">{{ scope.row.ID }}</a>-->
          </template>
        </el-table-column>
      </el-table>
      <template #footer>
        <div class="dialog-footer">
          <el-button size="small" @click="closeVarsDialog">取 消</el-button>
          <el-button v-if="canExecute" size="small" type="primary" @click="enterVarsDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>
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
          <el-button size="small" type="primary" @click="enterCheckVars">确 定</el-button>
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
  processSetTask, getSetTaskById, setTaskForceCorrect, getTemplateList
} from '@/api/template'
import { getTaskListBySetTaskId } from '@/api/task'
import { getSystemServerIds } from '@/api/cmdb'
import { emitter } from '@/utils/bus'
import warningBar from '@/components/warningBar/warningBar.vue'
import infoList from '@/mixins/infoList'
import TaskStatus from '@/components/task/TaskStatus.vue'
import { formatTimeToStr } from '@/utils/date'

export default {
  name: 'TemplateSetDetail',
  components: { TaskStatus, warningBar },
  mixins: [infoList],
  data() {
    return {
      listApi: getTaskListBySetTaskId,
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
        setTaskInnerSeq: 0,
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
      VarListVisible: false,
      varMap: new Map(),
      serverOptionsMap: new Map(),
      netDiskMap: new Map(),
      innerSeq: 0,
      templateOptions: [],
      canExecute: true,
    }
  },
  async created() {
    this.setTaskId = this.$route.params.setTaskId
    // await this.initSteps()
    await this.$nextTick(() => {
      this.initSteps()
    })
    const res = await getTemplateList({
      page: 1,
      pageSize: 99999
    })
    this.setOptions(res.data.list)
  },
  mounted() {
    emitter.on('i-close-task', () => {
      this.initSteps()
    })
    this.serverOptionsMap = new Map()
    this.varMap = new Map()
    this.netDiskMap = new Map()
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
        const tasks = this.setTask.tasks[seq]
        if (tasks.every(task => task.status === 'success')) {
          return 'success'
        }
        return 'error'
      } else {
        return 'wait'
      }
    },
    show(seq) {
      this.showTaskLog(this.setTask.tasks[seq])
    },
    initCommandVarsForm() {
      this.commandVarForm = {
        setTaskInnerSeq: 0,
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
      this.VarListVisible = true
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
    },
    closeVarsDialog() {
      this.VarListVisible = false
      this.initVarsList()
    },
    initVarsList() {
      this.varMap = []
    },
    async enterVars() {
      this.searchInfo.setTaskId = Number(this.setTaskId)
      this.searchInfo.currentSeq = Number(this.setTask.currentStep)
      await this.getTableData()
      this.setTask.templates[this.setTask.currentStep].forEach(template => {
        const innerCommandVarForm = {
          setTaskInnerSeq: template.seqInner,
          vars: [],
          targetIds: [],
        }
        innerCommandVarForm.targetIds = this.setCheckedServerOptionsNew(template)
        for (let i = 0; i < template.commandVarNumbers; i++) {
          innerCommandVarForm.vars.push('')
        }
        // this.runningTemplateId = this.setTask.templates[this.setTask.currentStep].ID
        // this.CommandVarFormVisible = true
        if (template.deployType === 2) {
          this.netDiskMap.set(template.seqInner, true)
        }
        this.varMap.set(template.seqInner, innerCommandVarForm)
      })
      this.VarListVisible = true
      this.canExecute = true
    },
    async setCheckedServerOptionsNew(template) {
      const res = await getSystemServerIds({
        ID: template.systemId
      })
      const serverOptions = res.data
      const innerTargetIds = []
      if (template.executeType !== 2) {
        serverOptions[0].children = serverOptions[0].children.filter((item) => {
          if (template.targetServerIds.includes(item.ID)) {
            innerTargetIds.push(item.ID)
            return true
          }
          return false
        })
      }
      this.serverOptionsMap.set(template.seqInner, serverOptions)
      return innerTargetIds
    },
    checkVars(innerSeq) {
      this.commandVarForm = this.varMap.get(innerSeq)
      this.checkedServerOptions = this.serverOptionsMap.get(innerSeq)
      this.netDisk = this.netDiskMap.get(innerSeq)
      this.VarListVisible = false
      this.CommandVarFormVisible = true
    },
    enterCheckVars() {
      this.$refs.CommandVarForm.validate(async valid => {
        if (valid) {
          const innerSeq = this.commandVarForm.setTaskInnerSeq
          this.varMap.set(innerSeq, this.commandVarForm)
          this.commandVarForm = []
          this.serverOptionsMap.set(innerSeq, this.serverOptions)
          this.serverOptions = []
          this.netDiskMap.set(innerSeq, this.netDisk)
          this.netDisk = []
        }
      })
      this.closeCommandVarsDialog()
    },
    async enterVarsDialog() {
      const data = []
      this.varMap.forEach((value, index) => {
        data.push({
          ID: Number(index),
          commandVars: value.vars,
          targetIds: value.targetIds,
          netDiskUser: value.netDiskUser ? value.netDiskUser : '',
          netDiskPassword: value.netDiskPassword ? value.netDiskPassword : '',
        })
      })
      await processSetTask({
        ID: this.setTask.ID,
        processTaskRequestVars: data
      })
      this.closeVarsDialog()
      await this.initSteps()
    },
    filterTemplateName(value) {
      const rowLabel = this.templateOptions.filter(item => item.ID === value)
      return rowLabel && rowLabel[0] && rowLabel[0].name
    },
    setOptions(data) {
      this.templateOptions = data
    },
    dateFormatter1(row) {
      if (row.beginTime.Time !== null && row.beginTime.Time !== '') {
        var date = new Date(row.beginTime.Time)
        return formatTimeToStr(date, 'yyyy-MM-dd hh:mm:ss')
      } else {
        return ''
      }
    },
    dateFormatter2(row) {
      if (row.endTime.Time !== null && row.endTime.Time !== '') {
        const date = new Date(row.endTime.Time)
        return formatTimeToStr(date, 'yyyy-MM-dd hh:mm:ss')
      } else {
        return ''
      }
    },
    async showTasks(seq) {
      this.searchInfo.setTaskId = Number(this.setTaskId)
      this.searchInfo.currentSeq = seq
      await this.getTableData()
      this.VarListVisible = true
      this.canExecute = false
    },
  },
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
.file-name {
  width: 100%;
  height: 100%;
  text-align: left;
  li {
    width: 100%;
    white-space:nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    i {
      margin-right: 8px;
    }
    padding: 10px 0;
    font-size: 16px;
    font-weight: 400;
    //color: #0154ff;
  }
  .directory {
    color: #01ff80;
  }
  .file {
    color: #0154ff;
  }
  li:hover {
    background: #f2f2f2;
    cursor: pointer;
  }
}
</style>
