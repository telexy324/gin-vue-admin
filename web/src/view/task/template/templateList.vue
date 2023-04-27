<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="searchForm" :inline="true" :model="searchInfo">
        <el-form-item label="模版名">
          <el-input v-model="searchInfo.name" placeholder="template" />
        </el-form-item>
        <el-form-item>
          <el-button size="mini" type="primary" icon="el-icon-search" @click="onSubmit">查询</el-button>
          <el-button size="mini" icon="el-icon-refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button size="mini" type="primary" icon="el-icon-plus" @click="openDialog('addTemplate')">新增</el-button>
        <el-button size="mini" type="primary" icon="el-icon-plus" @click="openLogDialog('addLogTemplate')">新增日志提取</el-button>
        <el-popover v-model:visible="deleteVisible" placement="top" width="160">
          <p>确定要删除吗？</p>
          <div style="text-align: right; margin-top: 8px;">
            <el-button size="mini" type="text" @click="deleteVisible = false">取消</el-button>
            <el-button size="mini" type="primary" @click="onDelete">确定</el-button>
          </div>
          <template #reference>
            <el-button icon="el-icon-delete" size="mini" :disabled="!templates.length" style="margin-left: 10px;">删除</el-button>
          </template>
        </el-popover>
        <el-button class="excel-btn" size="mini" type="success" icon="el-icon-download" @click="openDrawer()">选择系统</el-button>
      </div>
      <el-table :data="tableData" @sort-change="sortChange" @selection-change="handleSelectionChange">
        <el-table-column
          type="selection"
          width="55"
        />
        <el-table-column align="left" label="模版名" min-width="60" prop="name" sortable="custom" />
        <el-table-column align="left" label="最近执行状态" min-width="150" prop="lastTask.status" sortable="custom">
          <template v-slot="scope">
            <TaskStatus :status="scope.row.lastTask.status" />
          </template>
        </el-table-column>
        <el-table-column align="left" label="最近任务" min-width="200" sortable="custom">
          <template v-slot="scope">
            <el-button
              type="text"
              link
              @click="showTaskLog(scope.row.lastTask)"
            >#{{ scope.row.lastTask.ID }}</el-button>
          </template>
        </el-table-column>
        <el-table-column align="left" label="所属系统" min-width="150" prop="systemId" sortable="custom">
          <template #default="scope">
            <div>{{ filterSystemName(scope.row.systemId) }}</div>
          </template>
        </el-table-column>
        <el-table-column align="left" fixed="right" label="操作" width="250">
          <template #default="scope">
            <el-button
              icon="el-icon-edit"
              size="small"
              type="text"
              @click="editTemplate(scope.row)"
            >编辑</el-button>
            <el-button
              icon="el-icon-delete"
              size="small"
              type="text"
              @click="deleteTemplate(scope.row)"
            >删除</el-button>
            <el-button
              icon="el-icon-caret-right"
              size="small"
              type="text"
              @click="runTask(scope.row)"
            >构建</el-button>
            <el-button
              icon="el-icon-edit"
              size="small"
              type="text"
              @click="uploadScript(scope.row)"
            >上传脚本</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="gva-pagination">
        <el-pagination
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>

    </div>

    <el-dialog v-model="dialogFormVisible" :before-close="closeDialog" :title="dialogTitle">
      <warning-bar title="新增Template" />
      <el-form ref="templateForm" :model="form" :rules="rules" label-width="80px">
        <el-row>
          <el-col :span="12">
            <el-form-item label="模版名" prop="name">
              <el-input v-model="form.name" autocomplete="off" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="执行用户" prop="sysUser">
              <el-input v-model="form.sysUser" autocomplete="off" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="描述" prop="description">
          <el-input v-model="form.description" autocomplete="off" type="textarea"/>
        </el-form-item>
        <el-form-item label="目标" prop="targetServerIds">
          <el-cascader
            v-model="form.targetIds"
            style="width:100%"
            :options="serverOptions"
            :show-all-levels="false"
            :props="{ multiple:true,checkStrictly: true,label:'name',value:'ID',disabled:'disabled',emitPath:false}"
            :clearable="false"
          />
        </el-form-item>
        <el-row>
          <el-col :span="12">
            <el-form-item label="执行方式" prop="mode">
              <el-select v-model="form.mode" style="width:100%" @change="commandChange">
                <el-option :value="1" label="命令" />
                <el-option :value="2" label="脚本" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="所属系统" prop="systemId">
              <el-select v-model="form.systemId">
                <el-option v-for="val in systemOptions" :key="val.ID" :value="val.ID" :label="val.name" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item v-if="isCommand" label="命令" prop="command">
          <el-input v-model="form.command" autocomplete="off" type="textarea" />
        </el-form-item>
        <el-form-item v-if="isScript" label="脚本位置" prop="scriptPath">
          <el-input v-model="form.scriptPath" autocomplete="off" :disabled="true" />
        </el-form-item>
        <el-row>
          <el-col :span="6">
            <el-form-item v-if="isScript">
              <el-button size="small" type="primary" @click="checkScript">检查脚本</el-button>
            </el-form-item>
          </el-col>
          <el-col :span="6">
            <el-form-item v-if="isScript" label="脚本内容">
              <el-switch v-model="form.detail" />
            </el-form-item>
          </el-col>
        </el-row>
<!--        <el-form-item v-if="isScript">-->
<!--          <el-upload-->
<!--            ref="upload"-->
<!--            action=""-->
<!--            class="upload-demo"-->
<!--            :http-request="httpRequest"-->
<!--            :multiple="false"-->
<!--            :limit="1"-->
<!--            :auto-upload="false"-->
<!--            :file-list="form.file"-->
<!--          >-->
<!--            <el-button size="small" type="primary">选择脚本</el-button>-->
<!--          </el-upload>-->
<!--        </el-form-item>-->
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button size="small" @click="closeDialog">取 消</el-button>
          <el-button size="small" type="primary" @click="enterDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>

    <el-dialog v-model="dialogLogFormVisible" :before-close="closeLogDialog" :title="dialogLogTitle">
      <warning-bar title="新增日志提取模板" />
      <el-form ref="templateLogForm" :model="logForm" :rules="logRules" label-width="150px">
        <el-row>
          <el-col :span="12">
            <el-form-item label="模版名" prop="name">
              <el-input v-model="logForm.name" autocomplete="off" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="执行用户" prop="sysUser">
              <el-input v-model="logForm.sysUser" autocomplete="off" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="描述" prop="description">
          <el-input v-model="logForm.description" autocomplete="off" type="textarea" />
        </el-form-item>
        <el-form-item label="目标" prop="targetServerIds">
          <el-cascader
            v-model="logForm.targetIds"
            style="width:100%"
            :options="serverOptions"
            :show-all-levels="false"
            :props="{ multiple:false,checkStrictly: true,label:'name',value:'ID',disabled:'disabled',emitPath:false}"
            :clearable="false"
          />
        </el-form-item>
        <el-form-item label="所属系统" prop="systemId">
          <el-select v-model="logForm.systemId">
            <el-option v-for="val in systemOptions" :key="val.ID" :value="val.ID" :label="val.name" />
          </el-select>
        </el-form-item>
        <el-form-item label="日志文件夹位置" prop="logPath">
          <el-input v-model="logForm.logPath" autocomplete="off" />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button size="small" @click="closeLogDialog">取 消</el-button>
          <el-button size="small" type="primary" @click="enterLogDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>

    <el-dialog v-model="dialogFormVisibleScript" :before-close="closeScriptDialog" title="上传模板">
      <el-form ref="scriptForm" :model="scriptForm" :rules="rules" label-width="80px">
        <el-form-item label="脚本位置" prop="scriptPath">
          <el-input v-model="scriptForm.scriptPath" autocomplete="off" />
        </el-form-item>
        <el-form-item>
          <el-upload
            ref="upload"
            action=""
            class="upload-demo"
            :http-request="httpRequest"
            :multiple="false"
            :limit="1"
            :auto-upload="false"
            :file-list="fileList"
          >
            <el-button size="small" type="primary">选择脚本</el-button>
          </el-upload>
          <el-progress v-if="uploading" class="progress" :percentage="progressPercent" />
          <div v-for="(item, index) in scriptForm.items" :key="index">
            {{ item.manageIp }}
            <el-progress v-if="uploadingServer" class="progress-server" :percentage="100" :status="item.status" :indeterminate="item.indeterminate" :duration="2" />
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button size="small" @click="closeScriptDialog">取 消</el-button>
          <el-button size="small" type="primary" @click="submitUpload">确 定</el-button>
        </div>
      </template>
    </el-dialog>

    <el-dialog v-model="dialogFormVisibleDownload" :before-close="closeDownloadDialog" title="文件列表">
      <ul class="file-name">
        <li v-for="(item,index) in fileNames" :key="index" @click="downLoadFile(index)">
          <pre>{{ item }}</pre>
        </li>
      </ul>
      <template #footer>
        <div class="dialog-footer">
          <el-button size="small" @click="closeDownloadDialog">取 消</el-button>
        </div>
      </template>
    </el-dialog>
    <el-drawer v-if="drawer" v-model="drawer" :with-header="false" size="40%" title="请选择系统">
      <Systems ref="systems" :keys="searchInfo.systemIds" @checked="getCheckedTemplates" />
    </el-drawer>
  </div>
</template>

<script>

const path = import.meta.env.VITE_BASE_API
// 获取列表内容封装在mixins内部  getTableData方法 初始化已封装完成 条件搜索时候 请把条件安好后台定制的结构体字段 放到 this.searchInfo 中即可实现条件搜索

import {
  getTemplateList,
  deleteTemplate,
  getTemplateById,
  addTemplate,
  updateTemplate,
  checkScript,
  deleteTemplateByIds,
  getFileList,
  downloadFile,
} from '@/api/template'
import { addTask } from '@/api/task'
import { getAdminSystems, getAllServerIds } from '@/api/cmdb'
import infoList from '@/mixins/infoList'
import { toSQLLine } from '@/utils/stringFun'
import warningBar from '@/components/warningBar/warningBar.vue'
import { emitter } from '@/utils/bus'
import TaskStatus from '@/components/task/TaskStatus.vue'
import { ElMessage } from 'element-plus'
import { mapGetters } from 'vuex'
// import service from '@/utils/request'
import Axios from 'axios'
import socket from '@/socket'
import Systems from '@/components/task/systems.vue'

export default {
  name: 'TemplateList',
  components: {
    warningBar,
    TaskStatus,
    Systems
  },
  mixins: [infoList],
  data() {
    return {
      deleteVisible: false,
      listApi: getTemplateList,
      dialogFormVisible: false,
      dialogTitle: '新增template',
      templates: [],
      serverOptions: [],
      systemOptions: [],
      form: {
        ID: '',
        name: '',
        description: '',
        mode: 1,
        command: '',
        scriptPath: '',
        sysUser: '',
        targetIds: '',
        detail: false,
        systemId: '',
      },
      type: '',
      rules: {
        name: [{ required: true, message: '请输入模板名', trigger: 'blur' }],
        mode: [
          { required: true, message: '请选择执行方式', trigger: 'blur' }
        ]
      },
      path: path,
      isCommand: true,
      isScript: false,
      canCheck: false,
      dialogFormVisibleScript: false,
      scriptForm: {
        ID: '',
        scriptPath: '',
        file: '',
        items: []
      },
      hasFile: false,
      fileList: [],
      progressPercent: 0,
      uploading: false,
      uploadingServer: false,
      uploadingStatus: false,
      drawer: false,
      logForm: {
        ID: '',
        name: '',
        description: '',
        logPath: '',
        sysUser: '',
        targetIds: '',
        executeType: 2,
        systemId: '',
      },
      logRules: {
        name: [{ required: true, message: '请输入模板名', trigger: 'blur' }],
        logPath: [
          { required: true, message: '请输入日志文件夹路径', trigger: 'blur' }
        ]
      },
      dialogLogFormVisible: false,
      dialogLogTitle: '新增日志提取模板',
      logType: '',
      fileNames: [],
      fNames: [],
      currentTemplate: '',
      dialogFormVisibleDownload: false,
    }
  },
  computed: {
    ...mapGetters('user', ['userInfo', 'token'])
  },
  async created() {
    socket.addListener((data) => this.onWebsocketDataReceived(data))
    if (this.$route.params.systemIds) {
      this.searchInfo.systemIds = this.formRouterParam(this.$route.params.systemIds)
    }
    await this.getTableData()
    await this.setServerOptions()
    await this.setSystemOptions()
  },
  mounted() {
    emitter.on('i-close-task', () => {
      this.getTableData()
    })
  },
  methods: {
    //  选中api
    handleSelectionChange(val) {
      this.templates = val
    },
    async onDelete() {
      const ids = this.templates.map(item => item.ID)
      const res = await deleteTemplateByIds({ ids })
      if (res.code === 0) {
        this.$message({
          type: 'success',
          message: res.msg
        })
        if (this.tableData.length === ids.length && this.page > 1) {
          this.page--
        }
        this.deleteVisible = false
        await this.getTableData()
      }
    },
    // 排序
    sortChange({ prop, order }) {
      if (prop) {
        this.searchInfo.orderKey = toSQLLine(prop)
        this.searchInfo.desc = order === 'descending'
      }
      this.getTableData()
    },
    onReset() {
      this.searchInfo = {}
    },
    // 条件搜索前端看此方法
    onSubmit() {
      this.page = 1
      this.pageSize = 10
      this.getTableData()
    },
    initForm() {
      this.$refs.templateForm.resetFields()
      this.form = {
        ID: '',
        name: '',
        description: '',
        mode: '',
        command: '',
        scriptPath: '',
        sysUser: '',
        targetIds: '',
        detail: false
      }
    },
    closeDialog() {
      this.initForm()
      this.dialogFormVisible = false
    },
    openDialog(type) {
      switch (type) {
        case 'addTemplate':
          this.dialogTitle = '新增Template'
          this.canCheck = false
          break
        case 'edit':
          this.dialogTitle = '编辑Template'
          this.canCheck = true
          this.commandChange(this.form.mode)
          break
        default:
          break
      }
      this.type = type
      this.dialogFormVisible = true
    },
    async editTemplate(row) {
      const res = await getTemplateById({ id: row.ID })
      if (res.data.taskTemplate.executeType === 2) {
        this.logForm = res.data.taskTemplate
        this.logForm.targetIds = this.logForm.targetIds[0]
        this.openLogDialog('edit')
      } else {
        this.form = res.data.taskTemplate
        this.openDialog('edit')
      }
    },
    async deleteTemplate(row) {
      this.$confirm('此操作将永久删除服务器?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })
        .then(async() => {
          const res = await deleteTemplate(row)
          if (res.code === 0) {
            this.$message({
              type: 'success',
              message: '删除成功!'
            })
            if (this.tableData.length === 1 && this.page > 1) {
              this.page--
            }
            await this.getTableData()
          }
        })
    },
    async enterDialog() {
      this.$refs.templateForm.validate(async valid => {
        if (valid) {
          switch (this.type) {
            case 'addTemplate':
              {
                this.form.ID = 0
                const res = await addTemplate(this.form)
                if (res.code === 0) {
                  this.$message({
                    type: 'success',
                    message: '添加成功',
                    showClose: true
                  })
                }
                this.getTableData()
                this.closeDialog()
              }
              break
            case 'edit':
              {
                const res = await updateTemplate(this.form)
                if (res.code === 0) {
                  this.$message({
                    type: 'success',
                    message: '编辑成功',
                    showClose: true
                  })
                }
                this.getTableData()
                this.closeDialog()
              }
              break
            default:
              // eslint-disable-next-line no-lone-blocks
              {
                this.$message({
                  type: 'error',
                  message: '未知操作',
                  showClose: true
                })
              }
              break
          }
        }
      })
    },
    async setServerOptions() {
      const res = await getAllServerIds()
      this.serverOptions = res.data
    },
    async setSystemOptions() {
      const res = await getAdminSystems()
      this.systemOptions = res.data.systems
    },
    async runTask(row) {
      if (row.executeType === 2) {
        this.currentTemplate = row
        const fileNames = (await getFileList({
          ID: row.ID
        })).data.fileNames
        this.fileNames = fileNames.map((item) => {
          const fields = item.split(' ')
          const replacing = ' '.repeat(12 - fields[0].length)
          this.fNames.push(fields[1])
          return item.replace(' ', replacing)
        })
        console.log(this.fileNames)
        this.showFileList()
      } else {
        const task = (await addTask({
          templateId: row.ID
        })).data.task
        console.log(task.ID)
        this.showTaskLog(task)
      }
    },
    showTaskLog(task) {
      emitter.emit('i-show-task', task)
    },
    commandChange(selectValue) {
      if (selectValue === 1) {
        this.isCommand = true
        this.isScript = false
      } else {
        this.isCommand = false
        this.isScript = true
      }
    },
    async checkScript() {
      const res = (await checkScript({
        ID: this.form.ID,
        serverId: this.form.targetIds[0],
        detail: this.form.detail
      }))
      if (res.code !== 0) {
        ElMessage({
          showClose: true,
          message: '检查脚本失败',
          type: 'error'
        })
      }
      if (!res.data.exist) {
        ElMessage({
          showClose: true,
          message: '脚本不存在',
          type: 'error'
        })
      }
      if (!this.form.detail) {
        ElMessage({
          showClose: true,
          message: '检查成功',
          type: 'info'
        })
        if (this.form.detail) {
          this.closeDialog()
          this.showScript(res.data.script)
        }
      }
    },
    showScript(s) {
      emitter.emit('i-show-script', s)
    },
    async uploadScript(row) {
      const res = await getTemplateById({ id: row.ID })
      this.scriptForm = res.data.taskTemplate
      this.scriptForm.items = []
      for (let i = 0; i < res.data.taskTemplate.targetServers.length; i++) {
        const item = { ID: res.data.taskTemplate.targetServers[i].ID, manageIp: res.data.taskTemplate.targetServers[i].manageIp, status: 'warning', indeterminate: 'true' }
        this.scriptForm.items.push(item)
      }
      this.dialogFormVisibleScript = true
    },
    initScriptForm() {
      this.$refs.scriptForm.resetFields()
      this.scriptForm = {
        ID: '',
        scriptPath: '',
        file: '',
        items: [],
      }
    },
    closeScriptDialog() {
      this.progressPercent = 0
      this.uploading = false
      this.uploadingServer = false
      this.initScriptForm()
      this.dialogFormVisibleScript = false
    },
    handleRemove(file, fileList) {
      if (!fileList.length) {
        this.hasFile = false
      }
      this.scriptForm.file = null
    },
    handleChange(file, fileList) {
      if (fileList.length >= 2) {
        return
      }
      if (fileList.length === 1) {
        this.hasFile = true
      }
      this.scriptForm.file = file
    },
    submitUpload() {
      this.$refs.scriptForm.validate(valid => {
        if (valid) {
          this.$refs.upload.submit()
          this.uploading = true
        }
      })
    },
    httpRequest(param) {
      this.progressPercent = 0
      const fd = new FormData()
      fd.append('file', param.file)
      fd.append('scriptPath', this.scriptForm.scriptPath)
      fd.append('ID', this.scriptForm.ID)
      // console.log(this.token)
      // console.log(this.userInfo)
      // const res = await service({
      //   url: '/task/template/uploadScript',
      //   method: 'post',
      //   // headers: { 'Content-Type': 'multipart/form-data', 'x-token': this.token, 'x-user-id': this.user.ID },
      //   formData: fd
      // })
      // console.log(res.code)
      // if (res.code === 0) {
      //   this.$message({
      //     type: 'success',
      //     message: res.msg
      //   })
      // }
      // console.log('hahaha')
      // const config = {
      //   headers: {
      //     'Content-Type': 'multipart/form-data',
      //     'x-token': this.token,
      //     'x-user-id': this.user.ID,
      //   },
      //   timeout: 99999,
      // }
      Axios.post(import.meta.env.VITE_BASE_API + '/task/template/uploadScript', fd, {
        headers: {
          // 'Content-Type': 'multipart/form-data',
          'x-token': this.token,
          // 'x-user-id': this.user.ID,
        },
        // headers: { 'Content-Type': 'multipart/form-data', 'x-token': this.token, 'x-user-id': this.user.ID },
        timeout: 99999,
        onUploadProgress: (progressEvent) => {
          this.progressPercent = Math.floor((progressEvent.loaded * 100) / progressEvent.total)
          if (this.progressPercent >= 100) {
            this.uploadingServer = true
          }
        },
      }).then(response => {
        if (response.data.code === 0 || response.headers.success === 'true') {
          let message = '上传成功'
          if (response.data.data.failedIps.length > 0) {
            message = '上传部分成功, 失败的服务器: ' + response.data.data.failedIps.toString()
          }
          ElMessage({
            showClose: true,
            message: message,
            type: 'success'
          })
          this.closeScriptDialog()
        } else {
          ElMessage({
            showClose: true,
            message: response.data.msg,
            type: 'error'
          })
        }
      }).catch(err => {
        ElMessage({
          showClose: true,
          message: err,
          type: 'error'
        })
      })
    },
    onWebsocketDataReceived(data) {
      if (data.templateID !== this.scriptForm.ID) {
        return
      }
      if (data.type !== 'uploadScript') {
        return
      }
      console.log(this.scriptForm.items)
      for (let i = 0; i < this.scriptForm.items.length; i++) {
        if (this.scriptForm.items[i].ID === data.ID) {
          switch (data.status) {
            case 'success':
              this.scriptForm.items[i].status = 'success'
              this.scriptForm.items[i].indeterminate = this.uploadingStatus
              break
            case 'exception':
              this.scriptForm.items[i].status = 'exception'
              this.scriptForm.items[i].indeterminate = this.uploadingStatus
              break
            default:
              break
          }
        }
      }
    },
    openDrawer() {
      this.drawer = true
    },
    getCheckedTemplates(checkArr) {
      const systemIDs = []
      checkArr.forEach(item => {
        systemIDs.push(Number(item.ID))
      })
      this.searchInfo.systemIds = systemIDs
      this.getTableData()
      this.drawer = false
    },
    formRouterParam(ids) {
      const systemIDs = []
      ids.forEach(
        id => {
          systemIDs.push(Number(id))
        })
      return systemIDs
    },
    closeLogDialog() {
      this.initLogForm()
      this.dialogLogFormVisible = false
    },
    initLogForm() {
      this.$refs.templateLogForm.resetFields()
      this.form = {
        ID: '',
        name: '',
        description: '',
        logPath: '',
        sysUser: '',
        targetIds: '',
        executeType: 1,
      }
    },
    openLogDialog(type) {
      switch (type) {
        case 'addLogTemplate':
          this.dialogTitle = '新增日志提取模板'
          break
        case 'edit':
          this.dialogTitle = '编辑日志提取模板'
          break
        default:
          break
      }
      this.logType = type
      this.dialogLogFormVisible = true
    },
    async enterLogDialog() {
      this.$refs.templateLogForm.validate(async valid => {
        if (valid) {
          switch (this.logType) {
            case 'addLogTemplate':
              {
                this.logForm.ID = 0
                const _targetId = this.logForm.targetIds
                this.logForm.targetIds = [_targetId]
                const res = await addTemplate(this.logForm)
                if (res.code === 0) {
                  this.$message({
                    type: 'success',
                    message: '添加成功',
                    showClose: true
                  })
                }
                this.getTableData()
                this.closeLogDialog()
              }
              break
            case 'edit':
              {
                const _targetId = this.logForm.targetIds
                this.logForm.targetIds = [_targetId]
                const res = await updateTemplate(this.logForm)
                if (res.code === 0) {
                  this.$message({
                    type: 'success',
                    message: '编辑成功',
                    showClose: true
                  })
                }
                this.getTableData()
                this.closeLogDialog()
              }
              break
            default:
              // eslint-disable-next-line no-lone-blocks
              {
                this.$message({
                  type: 'error',
                  message: '未知操作',
                  showClose: true
                })
              }
              break
          }
        }
      })
    },
    showFileList() {
      this.dialogFormVisibleDownload = true
    },
    closeDownloadDialog() {
      this.fileNames = []
      this.currentTemplate = ''
      this.dialogFormVisibleDownload = false
    },
    downLoadFile(index) {
      const item = this.fNames[index]
      const id = this.currentTemplate.ID
      this.closeDownloadDialog()
      downloadFile(id, item)
    },
    filterSystemName(value) {
      const rowLabel = this.systemOptions.filter(item => item.ID === value)
      return rowLabel && rowLabel[0] && rowLabel[0].name
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
.warning {
  color: #dc143c;
}
.excel-btn+.excel-btn{
  margin-left: 10px;
}
.progress{
  width: 300px;
  margin-top: 8px
}
.progress-server{
  width: 300px
}
.excel-btn+.excel-btn{
  margin-left: 10px;
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
    color: #0154ff;
  }
  li:hover {
    background: #f2f2f2;
    cursor: pointer;
  }
}
</style>
