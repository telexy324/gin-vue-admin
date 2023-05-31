<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="searchForm" :inline="true" :model="searchInfo">
        <el-form-item label="密钥名">
          <el-input v-model="searchInfo.hostname" placeholder="hostname" />
        </el-form-item>
<!--        <el-form-item label="管理IP">-->
<!--          <el-input v-model="searchInfo.manageIp" placeholder="manageIp" />-->
<!--        </el-form-item>-->
        <el-form-item>
          <el-button size="mini" type="primary" icon="el-icon-search" @click="onSubmit">查询</el-button>
          <el-button size="mini" icon="el-icon-refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button size="mini" type="primary" icon="el-icon-plus" :disabled="!hasCreate" @click="openDialog('addSecret')">新增</el-button>
        <el-popover v-model:visible="deleteVisible" placement="top" width="160">
          <p>确定要删除吗？</p>
          <div style="text-align: right; margin-top: 8px;">
            <el-button size="mini" type="text" @click="deleteVisible = false">取消</el-button>
            <el-button size="mini" type="primary" @click="onDelete">确定</el-button>
          </div>
          <template #reference>
            <el-button icon="el-icon-delete" size="mini" :disabled="!secrets.length || !hasDelete" style="margin-left: 10px;">删除</el-button>
          </template>
        </el-popover>
      </div>
      <el-table :data="tableData" @sort-change="sortChange" @selection-change="handleSelectionChange">
        <el-table-column
          type="selection"
          width="55"
        />
        <el-table-column align="left" label="密钥名" min-width="150" prop="name" sortable="custom" />
        <el-table-column align="left" label="所属服务器" min-width="150" prop="systemId" sortable="custom">
          <template #default="scope">
            <div>{{ filterServer(scope.row.serverId) }}</div>
          </template>
        </el-table-column>
        <el-table-column align="left" fixed="right" label="操作" width="200">
          <template #default="scope">
            <el-button
              icon="el-icon-edit"
              size="small"
              type="text"
              :disabled="!hasEdit"
              @click="editSecret(scope.row)"
            >编辑</el-button>
            <el-button
              icon="el-icon-delete"
              size="small"
              type="text"
              :disabled="!hasDelete"
              @click="deleteSecret(scope.row)"
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
      <warning-bar title="新增上传日志用密钥" />
      <el-form ref="secretForm" :model="form" :rules="rules" label-width="120px">
        <el-form-item label="密钥名" prop="name">
          <el-input v-model="form.name" autocomplete="off" />
        </el-form-item>
        <el-form-item label="密码" prop="name">
          <el-input v-model="form.password" type="password" show-password autocomplete="off" />
        </el-form-item>
        <el-form-item label="所属服务器" prop="serverId">
          <el-select v-model="form.serverId">
            <el-option v-for="val in serverOptions" :key="val.ID" :value="val.ID" :label="val.hostname" />
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
  getSecretList,
  addSecret,
  updateSecret,
  deleteSecret,
  getSecretById,
  deleteSecretByIds
} from '@/api/logUpload'
import { getPolicyPathByAuthorityId } from '@/api/casbin'
import infoList from '@/mixins/infoList'
import { toSQLLine } from '@/utils/stringFun'
import warningBar from '@/components/warningBar/warningBar.vue'
import { mapGetters } from 'vuex'

export default {
  name: 'LogUploadSecret',
  components: {
    warningBar,
  },
  mixins: [infoList],
  data() {
    return {
      deleteVisible: false,
      listApi: getSecretList,
      dialogFormVisible: false,
      dialogTitle: '新增日志上传用密钥',
      secrets: [],
      form: {
        name: '',
        password: '',
        serverId: '',
      },
      rules: {
        name: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
        password: [
          { required: true, message: '请输入密码', trigger: 'blur' }
        ],
        serverId: [{ required: true, message: '请选择所属服务器', trigger: 'blur' }],
      },
      path: path,
      serverOptions: [],
      hasEdit: true,
      hasCreate: true,
      hasDelete: true,
    }
  },
  computed: {
    ...mapGetters('user', ['userInfo', 'token'])
  },
  async created() {
    this.serverOptions = (await getServerList({
      page: 1,
      pageSize: 99999
    })).data.list
    await this.getTableData()
  },
  mounted() {
    this.authorities()
  },
  methods: {
    //  选中api
    handleSelectionChange(val) {
      this.secrets = val
    },
    async onDelete() {
      const ids = this.secrets.map(item => item.ID)
      const res = await deleteSecretByIds({ ids })
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
      this.$refs.secretForm.resetFields()
      this.form = {
        name: '',
        password: '',
        serverId: '',
      }
    },
    closeDialog() {
      this.initForm()
      this.dialogFormVisible = false
    },
    openDialog(type) {
      switch (type) {
        case 'addSecret':
          this.dialogTitle = '新增日志上传用密钥'
          break
        case 'edit':
          this.dialogTitle = '编辑日志上传用密钥'
          break
        default:
          break
      }
      this.type = type
      this.dialogFormVisible = true
    },
    async editSecret(row) {
      const res = await getSecretById({ id: row.ID })
      this.form = res.data.secret
      this.openDialog('edit')
    },
    async deleteSecret(row) {
      this.$confirm('此操作将永久删除密钥?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })
        .then(async() => {
          const res = await deleteSecret(row)
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
      this.$refs.secretForm.validate(async valid => {
        if (valid) {
          switch (this.type) {
            case 'addSecret':
              {
                const res = await addSecret(this.form)
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
                const res = await updateSecret(this.form)
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
    filterServer(value) {
      const rowLabel = this.serverOptions.filter(item => item.ID === value)
      return rowLabel && rowLabel[0] && rowLabel[0].hostname
    },
    async authorities() {
      const res = await getPolicyPathByAuthorityId({
        authorityId: this.userInfo.authorityId
      })
      this.hasEdit = !!res.data.paths.some((item) => {
        return item.path === '/logUpload/updateSecret'
      })
      this.hasCreate = !!res.data.paths.some((item) => {
        return item.path === '/logUpload/addSecret'
      })
      this.hasDelete = !!res.data.paths.some((item) => {
        return item.path === '/logUpload/deleteSecret'
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
</style>
