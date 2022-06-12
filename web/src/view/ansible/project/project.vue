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
        <el-button size="mini" type="primary" icon="el-icon-plus" @click="openDialog('addProject')">新增</el-button>
        <el-popover v-model:visible="deleteVisible" placement="top" width="160">
          <p>确定要删除吗？</p>
          <div style="text-align: right; margin-top: 8px;">
            <el-button size="mini" type="text" @click="deleteVisible = false">取消</el-button>
            <el-button size="mini" type="primary" @click="onDelete">确定</el-button>
          </div>
          <template #reference>
            <el-button icon="el-icon-delete" size="mini" :disabled="!projects.length" style="margin-left: 10px;">删除</el-button>
          </template>
        </el-popover>
      </div>
      <el-table :data="tableData" @sort-change="sortChange" @selection-change="handleSelectionChange">
        <el-table-column
          type="selection"
          width="55"
        />
        <el-table-column align="left" label="id" min-width="60" prop="ID" sortable="custom" />
        <el-table-column align="left" label="alert" min-width="100" prop="alert" sortable="custom" />
        <el-table-column align="left" label="alert chat" min-width="200" prop="alertChat" sortable="custom" />
        <el-table-column align="left" label="name" min-width="200" prop="name" sortable="custom" />
        <el-table-column align="left" label="max parallel tasks" min-width="150" prop="maxParallelTasks" sortable="custom" />
        <el-table-column align="left" fixed="right" label="操作" width="200">
          <template #default="scope">
            <el-button
              icon="el-icon-edit"
              size="small"
              type="text"
              @click="editProject(scope.row)"
            >编辑</el-button>
            <el-button
              icon="el-icon-delete"
              size="small"
              type="text"
              @click="deleteProject(scope.row)"
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
      <warning-bar title="新增project" />
      <el-form ref="projectForm" :model="form" :rules="rules" label-width="80px">
        <el-form-item label="alert" prop="alert">
          <el-switch v-model="form.alert" />
        </el-form-item>
        <el-form-item label="alert chat" prop="alertChat">
          <el-input v-model="form.alertChat" autocomplete="off" />
        </el-form-item>
        <el-form-item label="name" prop="name">
          <el-input v-model="form.name" autocomplete="off" />
        </el-form-item>
        <el-form-item label="max parallel tasks" prop="maxParallelTasks">
          <el-input v-model="form.maxParallelTasks" autocomplete="off" />
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
  getProjectList,
  addProject,
  updateProject,
  deleteProject,
  getProjectById
} from '@/api/ansibleProject'
import infoList from '@/mixins/infoList'
import { toSQLLine } from '@/utils/stringFun'
import warningBar from '@/components/warningBar/warningBar.vue'

export default {
  name: 'Project',
  components: {
    warningBar
  },
  mixins: [infoList],
  data() {
    return {
      deleteVisible: false,
      listApi: getProjectList,
      dialogFormVisible: false,
      dialogTitle: '新增project',
      projects: [],
      form: {
        alert: false,
        alertChat: '',
        name: '',
        maxParallelTasks: ''
      },
      type: '',
      rules: {
        name: [{ required: true, message: '请输入project名', trigger: 'blur' }],
        alert: [
          { required: true, message: '请输入架构', trigger: 'blur' }
        ],
        alertChat: [
          { required: true, message: '请输入管理IP', trigger: 'blur' }
        ],
        maxParallelTasks: [
          { required: true, message: '请输入maxParallelTasks', trigger: 'blur' }
        ]
      },
      path: path
    }
  },
  created() {
    this.getTableData()
  },
  methods: {
    //  选中api
    handleSelectionChange(val) {
      this.projects = val
    },
    async onDelete() {
      const ids = this.projects.map(item => item.ID)
      const res = await deleteProject({ ids })
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
      this.$refs.projectForm.resetFields()
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
        case 'addProject':
          this.dialogTitle = '新增Project'
          break
        case 'edit':
          this.dialogTitle = '编辑Project'
          break
        default:
          break
      }
      this.type = type
      this.dialogFormVisible = true
    },
    async editProject(row) {
      const res = await getProjectById({ id: row.ID })
      this.form = res.data.project
      this.openDialog('edit')
    },
    async deleteProject(row) {
      this.$confirm('此操作将永久删除project?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })
        .then(async() => {
          const res = await deleteProject(row)
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
      this.$refs.projectForm.validate(async valid => {
        if (valid) {
          switch (this.type) {
            case 'addProject':
              {
                const res = await addProject(this.form)
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
                const res = await updateProject(this.form)
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
