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
      </div>
      <el-table :data="tableData" @sort-change="sortChange" @selection-change="handleSelectionChange">
        <el-table-column
          type="selection"
          width="55"
        />
        <el-table-column align="left" label="服务器名" min-width="150" prop="hostname" sortable="custom" />
        <el-table-column align="left" label="管理IP" min-width="200" prop="manageIp" sortable="custom" />
        <el-table-column align="left" label="运行方式" min-width="150" prop="systemId" sortable="custom">
          <template #default="scope">
            <div>{{ filterMode(scope.row.mode) }}</div>
          </template>
        </el-table-column>
        <el-table-column align="left" label="端口" min-width="150" prop="port" sortable="custom" />
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
      <warning-bar title="新增上传日志用服务器" />
      <el-form ref="serverForm" :model="form" :rules="rules" label-width="80px">
        <el-form-item label="服务器名" prop="hostname">
          <el-input v-model="form.hostname" autocomplete="off" />
        </el-form-item>
        <el-row>
          <el-col :span="12">
            <el-form-item label="管理IP" prop="manageIp">
              <el-input v-model="form.manageIp" autocomplete="off" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="端口" prop="port">
              <el-input v-model="form.port" autocomplete="off" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="执行方式" prop="systemId">
          <el-select v-model="form.mode">
            <el-option v-for="val in modeOptions" :key="val.key" :value="val.key" :label="val.value" />
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
  deleteServerByIds
} from '@/api/logUpload'
import { getPolicyPathByAuthorityId } from '@/api/casbin'
import infoList from '@/mixins/infoList'
import { toSQLLine } from '@/utils/stringFun'
import warningBar from '@/components/warningBar/warningBar.vue'
import { mapGetters } from 'vuex'

export default {
  name: 'LogUploadServer',
  components: {
    warningBar,
  },
  mixins: [infoList],
  data() {
    return {
      deleteVisible: false,
      listApi: getServerList,
      dialogFormVisible: false,
      dialogTitle: '新增日志上传用服务器',
      servers: [],
      form: {
        hostname: '',
        manageIp: '',
        port: '',
        mode: '',
      },
      rules: {
        hostname: [{ required: true, message: '请输入机器名', trigger: 'blur' }],
        manageIp: [
          { required: true, message: '请输入管理IP', trigger: 'blur' }
        ],
        port: [
          { required: true, message: '请输入连接端口', trigger: 'blur' },
          { validator: this.isNum, trigger: 'blur' }
        ],
      },
      path: path,
      modeOptions: [
        { 'key': 1, 'value': 'ftp' },
        { 'key': 2, 'value': 'sftp' },
      ],
      hasEdit: true,
      hasCreate: true,
      hasDelete: true,
    }
  },
  computed: {
    ...mapGetters('user', ['userInfo', 'token'])
  },
  async created() {
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
        manageIp: '',
        port: '',
        mode: '',
      }
    },
    closeDialog() {
      this.initForm()
      this.dialogFormVisible = false
    },
    openDialog(type) {
      switch (type) {
        case 'addServer':
          this.dialogTitle = '新增日志上传用服务器'
          break
        case 'edit':
          this.dialogTitle = '编辑日志上传用服务器'
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
    async enterDialog() {
      this.$refs.serverForm.validate(async valid => {
        if (valid) {
          this.form.port = Number(this.form.port)
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
    filterMode(value) {
      const rowLabel = this.modeOptions.filter(item => item.key === value)
      return rowLabel && rowLabel[0] && rowLabel[0].value
    },
    async authorities() {
      const res = await getPolicyPathByAuthorityId({
        authorityId: this.userInfo.authorityId
      })
      this.hasEdit = !!res.data.paths.some((item) => {
        return item.path === '/logUpload/updateServer'
      })
      this.hasCreate = !!res.data.paths.some((item) => {
        return item.path === '/logUpload/addServer'
      })
      this.hasDelete = !!res.data.paths.some((item) => {
        return item.path === '/logUpload/deleteServer'
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
