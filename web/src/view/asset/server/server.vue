<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="searchForm" :inline="true" :model="searchInfo">
        <el-form-item label="服务器名">
          <el-input v-model="searchInfo.hostname" placeholder="hostname" />
        </el-form-item>
        <el-form-item label="管理IP">
          <el-input v-model="searchInfo.manageIp" placeholder="manageIp" />
        </el-form-item>
        <el-form-item>
          <el-button size="mini" type="primary" icon="el-icon-search" @click="onSubmit">查询</el-button>
          <el-button size="mini" icon="el-icon-refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button size="mini" type="primary" icon="el-icon-plus" :disabled="!hasCreate" @click="openDialog('addServer')">新增</el-button>
        <el-popover v-model:visible="deleteVisible" placement="top" width="160">
          <p>确定要删除吗？</p>
          <div style="text-align: right; margin-top: 8px;">
            <el-button size="mini" type="text" @click="deleteVisible = false">取消</el-button>
            <el-button size="mini" type="primary" @click="onDelete">确定</el-button>
          </div>
          <template #reference>
            <el-button icon="el-icon-delete" size="mini" :disabled="!servers.length || !hasDelete" style="margin-left: 10px;">删除</el-button>
          </template>
        </el-popover>
        <el-upload
          class="excel-btn"
          :action="`${path}/cmdb/importExcel`"
          :headers="{'x-token':token}"
          :show-file-list="false"
        >
<!--          <el-button :disabled="disabledTemp" size="mini" type="primary" icon="el-icon-upload2">导入</el-button>-->
        </el-upload>
<!--        <el-button :disabled="disabledTemp" class="excel-btn" size="mini" type="primary" icon="el-icon-download" @click="handleExcelExport('ExcelExport.xlsx')">导出</el-button>-->
<!--        <el-button :disabled="disabledTemp" class="excel-btn" size="mini" type="success" icon="el-icon-download" @click="downloadExcelTemplate()">下载模板</el-button>-->
        <el-button class="excel-btn" size="mini" type="success" icon="el-icon-download" @click="openDrawer()">选择系统</el-button>
      </div>
      <el-table :data="tableData" @sort-change="sortChange" @selection-change="handleSelectionChange">
        <el-table-column
          type="selection"
          width="55"
        />
<!--        <el-table-column align="left" label="id" min-width="60" prop="ID" sortable="custom" />-->
        <el-table-column align="left" label="服务器名" min-width="150" prop="hostname" sortable="custom" />
        <el-table-column align="left" label="架构" min-width="150" prop="architecture" sortable="custom">
          <template #default="scope">
            <div>{{ filterDict(scope.row.architecture, 'architecture') }}</div>
          </template>
        </el-table-column>
        <el-table-column align="left" label="管理IP" min-width="200" prop="manageIp" sortable="custom" />
        <el-table-column align="left" label="操作系统" min-width="150" prop="os" sortable="custom">
          <template #default="scope">
            <div>{{ filterDict(scope.row.os, 'osVersion') }}</div>
          </template>
        </el-table-column>
        <el-table-column align="left" label="版本" min-width="100" prop="osVersion" sortable="custom" />
        <el-table-column align="left" label="所属系统" min-width="150" prop="systemId" sortable="custom">
          <template #default="scope">
            <div>{{ filterSystemName(scope.row.systemId) }}</div>
          </template>
        </el-table-column>
        <el-table-column align="left" fixed="right" label="操作" width="200">
          <template #default="scope">
            <el-button
              icon="el-icon-edit"
              size="small"
              type="text"
              :disabled="!hasEdit"
              @click="editServer(scope.row)"
            >编辑</el-button>
            <el-button
              icon="el-icon-delete"
              size="small"
              type="text"
              :disabled="!hasDelete"
              @click="deleteServer(scope.row)"
            >删除</el-button>
            <el-button
              icon="el-icon-orange"
              size="small"
              type="text"
              :disabled="!hasSsh"
              @click="runSsh(scope.row)"
            >执行</el-button>
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
      <warning-bar title="新增服务器" />
      <el-form ref="serverForm" :model="form" :rules="rules" label-width="80px">
        <el-form-item label="服务器名" prop="hostname">
          <el-input v-model="form.hostname" autocomplete="off" />
        </el-form-item>
        <el-form-item label="架构" prop="architecture">
<!--          <el-input v-model="form.architecture" autocomplete="off" />-->
          <el-select v-model="form.architecture">
            <el-option v-for="val in archs" :key="val.value" :value="val.value" :label="val.label" />
          </el-select>
        </el-form-item>
        <el-row>
          <el-col :span="12">
            <el-form-item label="管理IP" prop="manageIp">
              <el-input v-model="form.manageIp" autocomplete="off" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="ssh端口" prop="sshPort">
              <el-input v-model="form.sshPort" autocomplete="off" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="操作系统" prop="os">
          <el-select v-model="form.os">
            <el-option v-for="val in oss" :key="val.value" :value="val.value" :label="val.label" />
          </el-select>
        </el-form-item>
        <el-form-item label="版本" prop="osVersion">
          <el-input v-model="form.osVersion" autocomplete="off" />
        </el-form-item>
        <el-form-item label="所属系统" prop="systemId">
          <el-select v-model="form.systemId">
            <el-option v-for="val in systemOptions" :key="val.ID" :value="val.ID" :label="val.name" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button size="small" @click="closeDialog">取 消</el-button>
          <el-button size="small" type="primary" @click="enterDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>

    <el-dialog v-model="dialogSSHFormVisible" :before-close="closeDialog" :title="dialogSSHTitle">
      <warning-bar title="连接信息，必须使用密码连接" />
      <el-form ref="sshForm" :model="sshForm" :rules="sshRules" label-width="80px">
        <el-form-item label="管理IP" prop="manageIp">
          <el-input v-model="sshForm.server.manageIp" autocomplete="off" :disabled="true" />
        </el-form-item>
        <el-form-item label="SSH端口" prop="sshPort">
          <el-input v-model="sshForm.server.sshPort" autocomplete="off" :disabled="true" />
        </el-form-item>
        <el-form-item label="用户名" prop="username">
          <el-input v-model="sshForm.username" autocomplete="off" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="sshForm.password" autocomplete="off" />
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button size="small" @click="closeSSHDialog">取 消</el-button>
          <el-button size="small" type="primary" @click="enterSSHDialog">连 接</el-button>
        </div>
      </template>
    </el-dialog>
    <div class="term1">
      <div ref="terminalBox" style="height: 60vh;"></div>
    </div>
    <el-drawer v-if="drawer" v-model="drawer" :with-header="false" size="40%" title="请选择系统">
      <Systems ref="systems" @checked="getCheckedServers" />
    </el-drawer>
  </div>
</template>

<script>
const path = import.meta.env.VITE_BASE_API
// 获取列表内容封装在mixins内部  getTableData方法 初始化已封装完成 条件搜索时候 请把条件安好后台定制的结构体字段 放到 this.searchInfo 中即可实现条件搜索

import {
  getServerList,
  addServer,
  updateServer,
  deleteServer,
  getServerById,
  getAllServerIds,
  deleteServerByIds
} from '@/api/cmdb'
import { getPolicyPathByAuthorityId } from '@/api/casbin'
import infoList from '@/mixins/infoList'
import { toSQLLine } from '@/utils/stringFun'
import warningBar from '@/components/warningBar/warningBar.vue'
import { exportExcel, downloadTemplate } from '@/api/cmdb'
import { mapGetters } from 'vuex'
import Systems from '@/components/task/systems.vue'

export default {
  name: 'Server',
  components: {
    warningBar,
    Systems
  },
  mixins: [infoList],
  data() {
    return {
      deleteVisible: false,
      listApi: getServerList,
      dialogFormVisible: false,
      dialogSSHFormVisible: false,
      dialogTitle: '新增server',
      dialogSSHTitle: '服务器信息',
      servers: [],
      form: {
        hostname: '',
        architecture: '',
        manageIp: '',
        os: '',
        osVersion: '',
        systemId: '',
        sshPort: '',
      },
      sshForm: {
        server: {
          manageIp: '',
          sshPort: '',
        },
        username: '',
        password: ''
      },
      type: '',
      rules: {
        hostname: [{ required: true, message: '请输入机器名', trigger: 'blur' }],
        architecture: [
          { required: true, message: '请输入架构', trigger: 'blur' }
        ],
        manageIp: [
          { required: true, message: '请输入管理IP', trigger: 'blur' }
        ],
        sshPort: [
          { required: true, message: '请输入连接端口', trigger: 'blur' },
          { validator: this.isNum, trigger: 'blur' }
        ],
        os: [
          { required: true, message: '请输入系统名称', trigger: 'blur' }
        ],
        osVersion: [
          { required: true, message: '请输入系统版本', trigger: 'blur' }
        ],
        systemId: [
          { required: true, message: '请选择所属系统', trigger: 'blur' }
        ],
      },
      path: path,
      drawer: false,
      systemOptions: [],
      archs: [],
      oss: [],
      hasEdit: true,
      hasCreate: true,
      hasDelete: true,
      hasSsh: true,
      disabledTemp: true,
      sshRules: {
        username: [
          { required: true, message: '请输入用户名', trigger: 'blur' }
        ],
        password: [
          { required: true, message: '请输入密码', trigger: 'blur' }
        ],
      },
    }
  },
  computed: {
    ...mapGetters('user', ['userInfo', 'token'])
  },
  async created() {
    this.archs = await this.getDict('architecture')
    this.oss = await this.getDict('osVersion')
    this.systemOptions = (await getAllServerIds()).data
    await this.getTableData()
  },
  mounted() {
    this.authorities()
  },
  methods: {
    //  选中api
    handleSelectionChange(val) {
      this.servers = val
    },
    async onDelete() {
      const ids = this.servers.map(item => item.ID)
      const res = await deleteServerByIds({ ids })
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
      this.$refs.serverForm.resetFields()
      this.form = {
        hostname: '',
        architecture: '',
        manageIp: '',
        os: '',
        osVersion: '',
        systemId: '',
        sshPort: '',
      }
    },
    closeDialog() {
      this.initForm()
      this.dialogFormVisible = false
    },
    initSSHForm() {
      this.$refs.sshForm.resetFields()
      this.sshForm = {
        server: {
          manageIp: '',
          sshPort: '',
        },
        username: '',
        password: ''
      }
    },
    closeSSHDialog() {
      this.initSSHForm()
      this.dialogSSHFormVisible = false
    },
    openDialog(type) {
      switch (type) {
        case 'addServer':
          this.dialogTitle = '新增Server'
          break
        case 'edit':
          this.dialogTitle = '编辑Server'
          break
        default:
          break
      }
      this.type = type
      this.dialogFormVisible = true
    },
    async editServer(row) {
      const res = await getServerById({ id: row.ID })
      this.form = res.data.server
      this.openDialog('edit')
    },
    async deleteServer(row) {
      this.$confirm('此操作将永久删除服务器?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })
        .then(async() => {
          const res = await deleteServer(row)
          if (res.code === 0) {
            this.$message({
              type: 'success',
              message: '删除成功!'
            })
            if (this.tableData.length === 1 && this.page > 1) {
              this.page--
            }
            this.getTableData()
          }
        })
    },
    async relation(row) {
      const routeData = this.$router.resolve({
        name: 'graph',
        params: { cid: row.id }
      })
      window.open(routeData.href, '_blank')
    },
    async runSsh(row) {
      const res = await getServerById({ id: row.ID })
      this.sshForm.server = res.data.server
      this.openSSHDialog()
    },
    openSSHDialog() {
      this.dialogSSHFormVisible = true
    },
    async enterDialog() {
      this.$refs.serverForm.validate(async valid => {
        if (valid) {
          this.form.sshPort = Number(this.form.sshPort)
          switch (this.type) {
            case 'addServer':
              {
                const res = await addServer(this.form)
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
                const res = await updateServer(this.form)
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
    async enterSSHDialog() {
      await this.$router.push({
        name: 'ssh',
        params: {
          manageIp: this.sshForm.server.manageIp,
          username: this.sshForm.username,
          password: this.sshForm.password,
          sshPort: this.sshForm.server.sshPort
        }
      })
      // console.log(this.sshForm.server.manageIp)
      // window.open(routeData.href, '_self')
      this.closeSSHDialog()
    },
    handleExcelExport(fileName) {
      if (!fileName || typeof fileName !== 'string') {
        fileName = 'ExcelExport.xlsx'
      }
      exportExcel({
        fileName: fileName,
        infoList: this.tableData,
      })
    },
    downloadExcelTemplate() {
      downloadTemplate('ExcelTemplate.xlsx')
    },
    openDrawer() {
      this.drawer = true
    },
    getCheckedServers(checkArr) {
      const systemIDs = []
      checkArr.forEach(item => {
        systemIDs.push(Number(item.ID))
      })
      this.searchInfo.systemIDs = systemIDs
      this.getTableData()
      this.drawer = false
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
        return item.path === '/cmdb/updateServer'
      })
      this.hasCreate = !!res.data.paths.some((item) => {
        return item.path === '/cmdb/addServer'
      })
      this.hasDelete = !!res.data.paths.some((item) => {
        return item.path === '/cmdb/deleteServer'
      })
      this.hasSsh = !!res.data.paths.some((item) => {
        return item.path === '/ssh/run'
      })
    },
    isNum(rule, value, callback) {
      const n = /^[0-9]*$/
      if (!n.test(value)) {
        callback(new Error('只能为数字'))
      } else {
        callback()
      }
    }
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
</style>
