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
        <el-button size="mini" type="primary" icon="el-icon-plus" @click="openDialog('addServer')">新增</el-button>
        <el-popover v-model:visible="deleteVisible" placement="top" width="160">
          <p>确定要删除吗？</p>
          <div style="text-align: right; margin-top: 8px;">
            <el-button size="mini" type="text" @click="deleteVisible = false">取消</el-button>
            <el-button size="mini" type="primary" @click="onDelete">确定</el-button>
          </div>
          <template #reference>
            <el-button icon="el-icon-delete" size="mini" :disabled="!apis.length" style="margin-left: 10px;">删除</el-button>
          </template>
        </el-popover>
      </div>
      <el-table :data="tableData" @sort-change="sortChange" @selection-change="handleSelectionChange">
        <el-table-column
          type="selection"
          width="55"
        />
        <el-table-column align="left" label="id" min-width="60" prop="ID" sortable="custom" />
        <el-table-column align="left" label="服务器名" min-width="150" prop="hostname" sortable="custom" />
        <el-table-column align="left" label="架构" min-width="150" prop="architecture" sortable="custom" />
        <el-table-column align="left" label="管理IP" min-width="150" prop="manageIP" sortable="custom" />
        <el-table-column align="left" label="系统" min-width="150" prop="os" sortable="custom" />
        <el-table-column align="left" label="系统版本" min-width="150" prop="osVersion" sortable="custom" />
        <el-table-column align="left" fixed="right" label="操作" width="200">
          <template #default="scope">
            <el-button
              icon="el-icon-edit"
              size="small"
              type="text"
              @click="editServer(scope.row)"
            >编辑</el-button>
            <el-button
              icon="el-icon-delete"
              size="small"
              type="text"
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
      <warning-bar title="新增服务器" />
      <el-form ref="serverForm" :model="form" :rules="rules" label-width="80px">
        <el-form-item label="服务器名" prop="hostname">
          <el-input v-model="form.hostname" autocomplete="off" />
        </el-form-item>
        <el-form-item label="架构" prop="architecture">
          <el-input v-model="form.architecture" autocomplete="off" />
        </el-form-item>
        <el-form-item label="管理IP" prop="manageIP">
          <el-input v-model="form.manageIp" autocomplete="off" />
        </el-form-item>
        <el-form-item label="系统" prop="os">
          <el-input v-model="form.os" autocomplete="off" />
        </el-form-item>
        <el-form-item label="版本" prop="osVersion">
          <el-input v-model="form.osVersion" autocomplete="off" />
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
// 获取列表内容封装在mixins内部  getTableData方法 初始化已封装完成 条件搜索时候 请把条件安好后台定制的结构体字段 放到 this.searchInfo 中即可实现条件搜索

import {
  getServerList,
  addServer,
  updateServer,
  deleteServer,
  getServerById
} from '@/api/cmdb'
import infoList from '@/mixins/infoList'
import { toSQLLine } from '@/utils/stringFun'
import warningBar from '@/components/warningBar/warningBar.vue'

export default {
  name: 'Server',
  components: {
    warningBar
  },
  mixins: [infoList],
  data() {
    return {
      deleteVisible: false,
      listApi: getServerList,
      dialogFormVisible: false,
      dialogTitle: '新增server',
      servers: [],
      form: {
        hostname: '',
        architecture: '',
        manageIp: '',
        os: '',
        osVersion: ''
      },
      type: '',
      rules: {
        hostname: [{ required: true, message: '请输入机器名', trigger: 'blur' }],
        architecture: [
          { required: true, message: '请输入架构', trigger: 'blur' }
        ],
        manageIP: [
          { required: true, message: '请输入管理IP', trigger: 'blur' }
        ],
        os: [
          { required: true, message: '请输入系统名称', trigger: 'blur' }
        ],
        osVersion: [
          { required: true, message: '请输入系统版本', trigger: 'blur' }
        ]
      }
    }
  },
  created() {
    this.getTableData()
  },
  methods: {
    //  选中api
    handleSelectionChange(val) {
      this.servers = val
    },
    async onDelete() {
      const ids = this.servers.map(item => item.ID)
      const res = await deleteServer({ ids })
      if (res.code === 0) {
        this.$message({
          type: 'success',
          message: res.msg
        })
        if (this.tableData.length === ids.length && this.page > 1) {
          this.page--
        }
        this.deleteVisible = false
        this.getTableData()
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
        osVersion: ''
      }
    },
    closeDialog() {
      this.initForm()
      this.dialogFormVisible = false
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
      this.form = res.data.api
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
</style>
