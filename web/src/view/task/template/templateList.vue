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
        <el-button size="mini" type="primary" icon="el-icon-plus" :disabled="!hasCreate" @click="openDialog('addTemplate')">新增</el-button>
        <el-button size="mini" type="primary" icon="el-icon-plus" :disabled="!hasCreate" @click="openLogDialog('addLogTemplate')">新增日志提取</el-button>
        <el-button size="mini" type="primary" icon="el-icon-plus" :disabled="!hasCreate" @click="openDeployDialog('addDeployTemplate')">新增程序上传</el-button>
        <el-popover v-model:visible="deleteVisible" placement="top" width="160">
          <p>确定要删除吗？</p>
          <div style="text-align: right; margin-top: 8px;">
            <el-button size="mini" type="text" @click="deleteVisible = false">取消</el-button>
            <el-button size="mini" type="primary" @click="onDelete">确定</el-button>
          </div>
          <template #reference>
            <el-button icon="el-icon-delete" size="mini" :disabled="!templates.length || !hasDelete" style="margin-left: 10px;">删除</el-button>
          </template>
        </el-popover>
        <el-button class="excel-btn" size="mini" type="success" icon="el-icon-download" @click="openDrawer()">选择系统</el-button>
      </div>
      <el-table :data="tableData" @sort-change="sortChange" @selection-change="handleSelectionChange">
        <el-table-column
          type="selection"
          width="55"
        />
        <el-table-column align="left" label="模版名" min-width="150" prop="name" sortable="custom" />
        <el-table-column align="left" label="最近执行状态" min-width="150" prop="lastTask.status" sortable="custom">
          <template v-slot="scope">
            <TaskStatus :status="scope.row.lastTask.status" />
          </template>
        </el-table-column>
        <el-table-column align="left" label="最近任务" min-width="100" sortable="custom">
          <template v-slot="scope">
            <el-button
              type="text"
              link
              :style="{ display: scope.row.lastTask.ID?'':'none' }"
              @click="showTaskLog(scope.row.lastTask)"
            >#{{ scope.row.lastTask.ID }}</el-button>
          </template>
        </el-table-column>
        <el-table-column align="left" label="所属系统" min-width="150" prop="systemId" sortable="custom">
          <template #default="scope">
            <div>{{ filterSystemName(scope.row.systemId) }}</div>
          </template>
        </el-table-column>
        <el-table-column align="left" label="类型" min-width="150" prop="executeType" sortable="custom">
          <template #default="scope">
            <div>{{ filterExecuteType(scope.row.executeType) }}</div>
          </template>
        </el-table-column>
        <el-table-column align="left" fixed="right" label="操作" width="300">
          <template #default="scope">
            <el-button
              icon="el-icon-edit"
              size="small"
              type="text"
              :disabled="!hasEdit"
              @click="editTemplate(scope.row)"
            >编辑</el-button>
            <el-button
              icon="el-icon-delete"
              size="small"
              type="text"
              :disabled="!hasDelete"
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
              :disabled="scope.row.mode!==2"
              @click="uploadScript(scope.row)"
            >上传脚本</el-button>
            <el-popover :ref="`popover-${scope.$index}`" placement="top" width="160">
              <p>请输入模版名</p>
              <div style="text-align: right; margin-top: 8px;">
                <el-input v-model="newName" autocomplete="off" />
                <div style="text-align: right; margin-top: 10px;">
                  <el-button size="mini" type="text" @click="handleClose(scope.$index)">取消</el-button>
                  <el-button size="mini" type="primary" @click="copyTemplate(scope.row, scope.$index)">确定</el-button>
                </div>
              </div>
              <template #reference>
                <el-button icon="el-icon-copy-document" size="small" type="text">复制</el-button>
              </template>
            </el-popover>
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
      <warning-bar title="任务模板，可以对多个服务器生效" />
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
          <el-input v-model="form.description" autocomplete="off" type="textarea" />
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
              <el-select v-model="form.systemId" @change="changeSystemId">
                <el-option v-for="val in systemOptions" :key="val.ID" :value="val.ID" :label="val.name" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="目标" prop="targetIds">
          <el-cascader
            v-model="form.targetIds"
            style="width:100%"
            :options="serverOptions"
            :show-all-levels="false"
            :props="{ multiple:true,checkStrictly: false,label:'name',value:'ID',disabled:'disabled',emitPath:false}"
            :clearable="true"
          />
        </el-form-item>
        <el-form-item v-if="isCommand" label="命令" prop="command">
          <el-input v-model="form.command" autocomplete="off" type="textarea" :rows="10" />
        </el-form-item>
        <el-form-item v-if="isScript" label="脚本位置" prop="scriptPath">
          <el-input v-model="form.scriptPath" autocomplete="off" />
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
        <el-row>
          <el-col :span="6">
            <el-form-item v-if="isScript" label="shell方式" prop="scriptType">
              <el-select v-model="form.shellType">
                <el-option v-for="val in shellTypeOptions" :key="val.key" :value="val.key" :label="val.value" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="18">
            <el-form-item v-if="isScript" label="脚本参数" prop="scriptVars">
              <el-input v-model="form.shellVars" autocomplete="off" />
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
      <warning-bar title="日志提取模板，仅可用于下载日志，只对单个服务器生效" />
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
        <el-form-item label="所属系统" prop="systemId">
          <el-select v-model="logForm.systemId" @change="changeSystemId">
            <el-option v-for="val in systemOptions" :key="val.ID" :value="val.ID" :label="val.name" />
          </el-select>
        </el-form-item>
        <el-form-item label="目标" prop="targetIds">
          <el-cascader
            v-model="logForm.targetIds"
            style="width:100%"
            :options="serverOptions"
            :show-all-levels="false"
            :props="{ multiple:false,checkStrictly: false,label:'name',value:'ID',disabled:'disabled',emitPath:false}"
            :clearable="false"
          />
        </el-form-item>
        <el-form-item label="日志文件夹位置" prop="logPath">
          <el-input v-model="logForm.logPath" autocomplete="off" />
        </el-form-item>
        <el-form-item label="执行方式" prop="logOutput">
          <el-select v-model="logForm.logOutput" @change="logOutputChange">
            <el-option v-for="val in logOutputOptions" :key="val.ID" :value="val.ID" :label="val.name" />
          </el-select>
        </el-form-item>
        <el-form-item v-if="!downloadDirectly" label="上传位置" prop="logDst">
          <el-input v-model="logForm.logDst" autocomplete="off" />
        </el-form-item>
        <el-row>
          <el-col :span="12">
            <el-form-item v-if="!downloadDirectly" label="日志服务器" prop="dstServerId">
              <el-select v-model="logForm.dstServerId" @change="changeServerId">
                <el-option v-for="val in logServerOptions" :key="val.ID" :value="val.ID" :label="val.hostname" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item v-if="!downloadDirectly" label="上传用户" prop="secretId">
              <el-select v-model="logForm.secretId">
                <el-option v-for="val in logSecretOptionsFiltered" :key="val.ID" :value="val.ID" :label="val.name" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
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
          <el-input v-model="scriptForm.scriptPath" autocomplete="off" :disabled="true" />
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
          <el-button :disabled="uploadingDisable" size="small" type="primary" @click="submitUpload">确 定</el-button>
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

    <el-dialog v-model="dialogDeployFormVisible" :before-close="closeDeployDialog" :title="dialogDeployTitle">
      <warning-bar title="程序上传模板" />
      <el-form ref="templateDeployForm" :model="deployForm" :rules="deployRules" label-width="150px">
        <el-row>
          <el-col :span="12">
            <el-form-item label="模版名" prop="name">
              <el-input v-model="deployForm.name" autocomplete="off" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="执行用户" prop="sysUser">
              <el-input v-model="deployForm.sysUser" autocomplete="off" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="描述" prop="description">
          <el-input v-model="deployForm.description" autocomplete="off" type="textarea" />
        </el-form-item>
        <el-form-item label="所属系统" prop="systemId">
          <el-select v-model="deployForm.systemId" @change="changeSystemId">
            <el-option v-for="val in systemOptions" :key="val.ID" :value="val.ID" :label="val.name" />
          </el-select>
        </el-form-item>
        <el-form-item label="目标" prop="targetIds">
          <el-cascader
            v-model="deployForm.targetIds"
            style="width:100%"
            :options="serverOptions"
            :show-all-levels="false"
            :props="{ multiple:true,checkStrictly: false,label:'name',value:'ID',disabled:'disabled',emitPath:false}"
            :clearable="false"
          />
        </el-form-item>
<!--        <el-form-item label="程序包上传位置" prop="deployPath">-->
<!--          <el-input v-model="deployForm.deployPath" autocomplete="off" />-->
<!--        </el-form-item>-->
<!--        <el-form-item label="源文件位置" prop="downloadSource">-->
<!--          <el-input v-model="deployForm.downloadSource" autocomplete="off" />-->
<!--        </el-form-item>-->
        <el-row>
          <el-col :span="12">
            <el-form-item label="ftp/sftp服务器" prop="dstServerId">
              <el-select v-model="deployForm.dstServerId" @change="changeServerId">
                <el-option v-for="val in logServerOptions" :key="val.ID" :value="val.ID" :label="val.hostname" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="上传用户" prop="secretId">
              <el-select v-model="deployForm.secretId">
                <el-option v-for="val in logSecretOptionsFiltered" :key="val.ID" :value="val.ID" :label="val.name" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <div v-for="(item, index) in deployForm.taskDeployInfos" :key="index">
          <el-row>
            <el-col :span="18">
              <el-form-item
                :label="'程序包上传位置' + index"
                :prop="'taskDeployInfos.' + index + '.deployPath'"
                :rules="deployRules.deployPath"
              >
                <el-input v-model="item.deployPath" placeholder="全路径，需要包含文件名" />
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col :span="18">
              <el-form-item
                :label="'源文件位置' + index"
                :prop="'taskDeployInfos.' + index + '.downloadSource'"
                :rules="deployRules.downloadSource"
              >
                <el-input v-model="item.downloadSource" placeholder="只能单文件，全路径，不可为文件夹" />
              </el-form-item>
            </el-col>
            <el-col :span="2">
              <el-form-item>
                <i class="el-icon-delete" @click="deleteItem(item, index)" />
              </el-form-item>
            </el-col>
          </el-row>
        </div>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button size="small" type="warning" @click="addItem">增加上传</el-button>
          <el-button size="small" @click="closeDeployDialog">取 消</el-button>
          <el-button size="small" type="primary" @click="enterDeployDialog">确 定</el-button>
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
  addTemplate,
  checkScript,
  deleteTemplate,
  deleteTemplateByIds,
  deployServer,
  downloadFile,
  getFileList,
  getTemplateById,
  getTemplateList,
  updateTemplate,
  uploadLogServer,
} from '@/api/template'
import { addTask } from '@/api/task'
import { getAdminSystems, getSystemServerIds } from '@/api/cmdb'
import { getPolicyPathByAuthorityId } from '@/api/casbin'
import { getSecretList, getServerList } from '@/api/logUpload'
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
      dialogTitle: '新增任务模板',
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
        targetIds: [],
        detail: false,
        systemId: '',
        executeType: 1,
        shellType: '',
        shellVars: '',
      },
      type: '',
      rules: {
        name: [{ required: true, message: '请输入模板名', trigger: 'blur' }],
        mode: [
          { required: true, message: '请选择执行方式', trigger: 'blur' }
        ],
        sysUser: [{ required: true, message: '请输入执行用户', trigger: 'blur' }],
        systemId: [{ required: true, message: '请选择所属系统', trigger: 'blur' }],
        targetIds: [{ required: true, message: '请选择目标', trigger: 'blur' }],
        command: [{ required: true, message: '请输入命令', trigger: 'blur' }],
        scriptPath: [{ required: true, message: '请输入脚本位置', trigger: 'blur' }],
        scriptType: [{ required: true, message: '请选择shell方式', trigger: 'blur' }],
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
      uploadingDisable: false,
      drawer: false,
      logForm: {
        ID: '',
        name: '',
        description: '',
        logPath: '',
        sysUser: '',
        targetIds: [],
        executeType: 2,
        systemId: '',
        logOutput: '',
        logDst: '',
        dstServerId: '',
        secretId: ''
      },
      logRules: {
        name: [{ required: true, message: '请输入模板名', trigger: 'blur' }],
        logPath: [
          { required: true, message: '请输入日志文件夹路径', trigger: 'blur' }
        ],
        sysUser: [{ required: true, message: '请输入执行用户', trigger: 'blur' }],
        systemId: [{ required: true, message: '请选择所属系统', trigger: 'blur' }],
        targetIds: [{ required: true, message: '请选择目标', trigger: 'blur' }],
        logOutput: [{ required: true, message: '请选择下载方式', trigger: 'blur' }],
        logDst: [{ required: true, message: '请输入上传位置', trigger: 'blur' }],
        dstServerId: [{ required: true, message: '请选择日志服务器', trigger: 'blur' }],
        secretId: [{ required: true, message: '请选择上传用户', trigger: 'blur' }],
      },
      dialogLogFormVisible: false,
      dialogLogTitle: '新增日志提取模板',
      logType: '',
      fileNames: [],
      fNames: [],
      currentTemplate: '',
      dialogFormVisibleDownload: false,
      hasEdit: true,
      hasCreate: true,
      hasDelete: true,
      logOutputOptions: [
        { ID: 1, name: '直接下载' },
        { ID: 2, name: '上传服务器' }
      ],
      downloadDirectly: true,
      logServerOptions: [],
      logSecretOptions: [],
      logSecretOptionsFiltered: [],
      shellTypeOptions: [
        { 'key': 1, 'value': 'sh' },
        { 'key': 2, 'value': 'bash' },
      ],
      deployForm: {
        ID: '',
        name: '',
        description: '',
        sysUser: '',
        targetIds: [],
        executeType: 3,
        systemId: '',
        dstServerId: '',
        secretId: '',
        deployInfos: '',
        taskDeployInfos: [],
      },
      deployRules: {
        name: [{ required: true, message: '请输入模板名', trigger: 'blur' }],
        sysUser: [{ required: true, message: '请输入上传用户', trigger: 'blur' }],
        systemId: [{ required: true, message: '请选择所属系统', trigger: 'blur' }],
        targetIds: [{ required: true, message: '请选择目标', trigger: 'blur' }],
        deployPath: [{ required: true, message: '请输入上传位置', trigger: 'blur' }],
        downloadSource: [{ required: true, message: '请输入源程序包位置', trigger: 'blur' }],
        dstServerId: [{ required: true, message: '请选择日志服务器', trigger: 'blur' }],
        secretId: [{ required: true, message: '请选择下载用户', trigger: 'blur' }],
      },
      dialogDeployFormVisible: false,
      dialogDeployTitle: '新增程序上传模板',
      deployType: '',
      newName: '',
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
    // await this.setServerOptions()
    await this.setSystemOptions()
  },
  mounted() {
    emitter.on('i-close-task', () => {
      this.getTableData()
    })
    this.authorities()
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
        targetIds: [],
        detail: false,
        executeType: 1,
        shellType: '',
        shellVars: '',
      }
    },
    closeDialog() {
      this.initForm()
      this.dialogFormVisible = false
      this.serverOptions = []
    },
    openDialog(type) {
      switch (type) {
        case 'addTemplate':
          this.dialogTitle = '新增任务模板'
          this.canCheck = false
          break
        case 'edit':
          this.dialogTitle = '编辑任务模板'
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
        const temp = res.data.taskTemplate
        this.logForm = temp
        await this.setServerOptions(temp.systemId)
        this.logForm.targetIds = this.logForm.targetIds[0]
        temp.logOutput === 2 ? this.downloadDirectly = false : this.downloadDirectly = true
        await this.openLogDialog('edit')
      } else if (res.data.taskTemplate.executeType === 3) {
        const temp = res.data.taskTemplate
        this.deployForm = temp
        await this.setServerOptions(temp.systemId)
        await this.openDeployDialog('edit')
      } else {
        this.form = res.data.taskTemplate
        await this.setServerOptions(this.form.systemId)
        this.openDialog('edit')
      }
    },
    async deleteTemplate(row) {
      this.$confirm('此操作将永久删除任务模板?', '提示', {
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
                if (this.form.shellType === '') {
                  this.form.shellType = 0
                }
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
    async setSystemOptions() {
      const res = await getAdminSystems()
      this.systemOptions = res.data.systems
    },
    async setServerOptions(systemId) {
      const res = await getSystemServerIds({
        ID: systemId
      })
      this.serverOptions = res.data
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
      } else if (row.executeType === 3) {
        const task = (await deployServer({
          ID: row.ID
        })).data.task
        console.log(task.ID)
        this.showTaskLog(task)
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
        detail: this.form.detail
      }))
      if (res.code !== 0) {
        ElMessage({
          showClose: true,
          message: '检查脚本失败',
          type: 'error'
        })
        return
      }
      if (res.data.failedIps.length > 0) {
        const msg = '以下服务器脚本不存在' + res.data.failedIps
        ElMessage({
          showClose: true,
          message: msg,
          type: 'error'
        })
        return
      }
      if (!this.form.detail) {
        ElMessage({
          showClose: true,
          message: '检查成功',
          type: 'success'
        })
      } else {
        this.closeDialog()
        this.showScript(res.data.script)
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
      this.uploadingDisable = false
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
      this.uploadingDisable = true
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
        timeout: 999999,
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
      if (this.scriptForm.items.every((item) => {
        return item.status === 'success'
      })) {
        ElMessage({
          showClose: true,
          message: '上传成功',
          type: 'success'
        })
        this.closeScriptDialog()
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
      this.serverOptions = []
    },
    initLogForm() {
      this.$refs.templateLogForm.resetFields()
      this.logForm = {
        ID: '',
        name: '',
        description: '',
        logPath: '',
        sysUser: '',
        targetIds: [],
        executeType: 1,
      }
    },
    async openLogDialog(type) {
      this.logServerOptions = (await getServerList({
        page: 1,
        pageSize: 99999
      })).data.list
      this.logSecretOptions = (await getSecretList({
        page: 1,
        pageSize: 99999
      })).data.list
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
      this.changeServerId(this.logForm.dstServerId)
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
    async downLoadFile(index) {
      const item = this.fNames[index]
      const id = this.currentTemplate.ID
      const type = this.currentTemplate.logOutput
      this.closeDownloadDialog()
      if (type === 2) {
        const task = (await uploadLogServer({
          ID: id,
          file: item,
        })).data.task
        console.log(task.ID)
        this.dialogFormVisibleDownload = false
        this.showTaskLog(task)
      } else {
        downloadFile(id, item)
      }
    },
    filterSystemName(value) {
      const rowLabel = this.systemOptions.filter(item => item.ID === value)
      return rowLabel && rowLabel[0] && rowLabel[0].name
    },
    async authorities() {
      const res = await getPolicyPathByAuthorityId({
        authorityId: this.userInfo.authorityId
      })
      this.hasEdit = !!res.data.paths.some((item) => {
        return item.path === '/task/template/updateTemplate'
      })
      this.hasCreate = !!res.data.paths.some((item) => {
        return item.path === '/task/template/addTemplate'
      })
      this.hasDelete = !!res.data.paths.some((item) => {
        return item.path === '/task/template/deleteTemplate'
      })
    },
    filterExecuteType(value) {
      switch (value) {
        case 1:
          return '普通'
        case 2:
          return '日志提取'
        default:
          return ''
      }
    },
    logOutputChange(selectValue) {
      selectValue === 1 ? this.downloadDirectly = true : this.downloadDirectly = false
    },
    changeServerId(selectValue) {
      this.logSecretOptionsFiltered = this.logSecretOptions.filter(item => {
        return item.serverId === selectValue
      })
    },
    changeSystemId(selectValue) {
      this.form.targetIds = []
      this.setServerOptions(selectValue)
    },
    initDeployForm() {
      this.$refs.templateDeployForm.resetFields()
      this.deployForm = {
        ID: '',
        name: '',
        description: '',
        sysUser: '',
        targetIds: [],
        executeType: 3,
        systemId: '',
        dstServerId: '',
        secretId: '',
        taskDeployInfos: [],
      }
    },
    async openDeployDialog(type) {
      this.logServerOptions = (await getServerList({
        page: 1,
        pageSize: 99999
      })).data.list
      this.logSecretOptions = (await getSecretList({
        page: 1,
        pageSize: 99999
      })).data.list
      switch (type) {
        case 'addDeployTemplate':
          this.dialogTitle = '新增程序上传模板'
          break
        case 'edit':
          this.dialogTitle = '编辑程序上传模板'
          break
        default:
          break
      }
      this.deployType = type
      this.changeServerId(this.deployForm.dstServerId)
      this.dialogDeployFormVisible = true
    },
    async enterDeployDialog() {
      this.$refs.templateDeployForm.validate(async valid => {
        if (valid) {
          switch (this.deployType) {
            case 'addDeployTemplate':
              {
                this.deployForm.ID = 0
                const res = await addTemplate(this.deployForm)
                if (res.code === 0) {
                  this.$message({
                    type: 'success',
                    message: '添加成功',
                    showClose: true
                  })
                }
                this.getTableData()
                this.closeDeployDialog()
              }
              break
            case 'edit':
              {
                const res = await updateTemplate(this.deployForm)
                if (res.code === 0) {
                  this.$message({
                    type: 'success',
                    message: '编辑成功',
                    showClose: true
                  })
                }
                this.getTableData()
                this.closeDeployDialog()
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
    closeDeployDialog() {
      this.initDeployForm()
      this.dialogDeployFormVisible = false
      this.logServerOptions = []
      this.logSecretOptions = []
    },
    handleClose(index) {
      this.$refs[`popover-${index}`].hide()
      this.newName = ''
    },
    async copyTemplate(row, index) {
      const res = (await getTemplateById({ id: row.ID }))
      if (res.data.taskTemplate.executeType === 2) {
        this.logForm = res.data.taskTemplate
        this.logForm.name = this.newName
        this.logForm.ID = 0
        this.logForm.lastTaskId = 0
        const res1 = await addTemplate(this.logForm)
        if (res1.code === 0) {
          this.$message({
            type: 'success',
            message: '复制成功',
            showClose: true
          })
        }
        this.handleClose(index)
        this.getTableData()
        this.logForm = {
          ID: '',
          name: '',
          description: '',
          logPath: '',
          sysUser: '',
          targetIds: [],
          executeType: 1,
        }
      } else if (res.data.taskTemplate.executeType === 3) {
        this.deployForm = res.data.taskTemplate
        this.deployForm.name = this.newName
        this.deployForm.ID = 0
        this.deployForm.lastTaskId = 0
        const res1 = await addTemplate(this.deployForm)
        if (res1.code === 0) {
          this.$message({
            type: 'success',
            message: '复制成功',
            showClose: true
          })
        }
        this.handleClose(index)
        this.getTableData()
        this.deployForm = {
          ID: '',
          name: '',
          description: '',
          sysUser: '',
          targetIds: [],
          executeType: 3,
          systemId: '',
          deployPath: '',
          downloadSource: '',
          dstServerId: '',
          secretId: ''
        }
      } else {
        this.form = res.data.taskTemplate
        this.form.name = this.newName
        this.form.ID = 0
        this.form.lastTaskId = 0
        const res1 = await addTemplate(this.form)
        if (res1.code === 0) {
          this.$message({
            type: 'success',
            message: '复制成功',
            showClose: true
          })
        }
        this.handleClose(index)
        this.getTableData()
        this.form = {
          ID: '',
          name: '',
          description: '',
          mode: '',
          command: '',
          scriptPath: '',
          sysUser: '',
          targetIds: [],
          detail: false,
          executeType: 1,
          shellType: '',
          shellVars: '',
        }
      }
    },
    deleteItem(item, index) {
      this.deployForm.taskDeployInfos.splice(index, 1)
    },
    addItem() {
      this.deployForm.taskDeployInfos.push({
        deployPath: '',
        downloadSource: ''
      })
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
