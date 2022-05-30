<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="searchForm" :inline="true" :model="searchInfo">
        <el-form-item label="environment name">
          <el-input v-model="searchInfo.name" placeholder="name" />
        </el-form-item>
        <el-form-item label="project id">
          <el-input v-model="searchInfo.projectIp" placeholder="projectIp" />
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
            <el-button icon="el-icon-delete" size="mini" :disabled="!servers.length" style="margin-left: 10px;">删除</el-button>
          </template>
        </el-popover>
        <el-dropdown>
          <div class="dp-flex justify-content-center align-items height-full width-full">
            <span style="cursor: pointer">
              <span style="margin-left: 5px">{{ currentProject.Name }}</span>
              <i class="el-icon-arrow-down" />
            </span>
          </div>
          <template #dropdown>
            <el-dropdown-menu class="dropdown-group">
              <el-dropdown-item v-for="(item,index) in projects" :key="index" @click="setCurrentProject" v-text="item.Name"></el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
      <el-table :data="tableData" @sort-change="sortChange" @selection-change="handleSelectionChange">
        <el-table-column
          type="selection"
          width="55"
        />
        <el-table-column align="left" label="id" min-width="60" prop="ID" sortable="custom" />
        <el-table-column align="left" label="name" min-width="150" prop="name" sortable="custom" />
        <el-table-column align="left" label="project_id" min-width="150" prop="projectIp" sortable="custom" />
        <el-table-column align="left" fixed="right" label="操作" width="200">
          <template #default="scope">
            <el-button
              icon="el-icon-edit"
              size="small"
              type="text"
              @click="editEnvironment(scope.row)"
            >编辑</el-button>
            <el-button
              icon="el-icon-delete"
              size="small"
              type="text"
              @click="deleteEnvironment(scope.row)"
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
        <el-form-item label="name" prop="name">
          <el-input v-model="form.hostname" autocomplete="off" />
        </el-form-item>
        <el-form-item label="projectId" prop="projectId">
          <el-input v-model="form.architecture" autocomplete="off" />
        </el-form-item>
        <el-form-item label="json" prop="json">
          <el-input v-model="form.json" type="textarea" />
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
  getEnvironmentList,
  addEnvironment,
  updateEnvironment,
  deleteEnvironment,
  getEnvironmentById
} from '@/api/ansibleEnvironment'
import { getProjectList } from '@/api/ansibleProject'
import infoList from '@/mixins/infoList'
import { toSQLLine } from '@/utils/stringFun'
import warningBar from '@/components/warningBar/warningBar.vue'

export default {
  name: 'Environment',
  components: {
    warningBar
  },
  mixins: [infoList],
  data() {
    return {
      deleteVisible: false,
      listApi: getEnvironmentList,
      dialogFormVisible: false,
      dialogTitle: '新增environment',
      environments: [],
      projects: [],
      currentProject: {},
      form: {
        name: '',
        projectIp: '',
        json: ''
      },
      type: '',
      rules: {
        name: [{ required: true, message: '请输入environment name', trigger: 'blur' }],
        projectIp: [
          { required: true, message: '请输入project ip', trigger: 'blur' }
        ]
      },
      path: path
    }
  },
  created() {
    this.getProjects()
    this.getTableData()
  },
  methods: {
    //  选中api
    handleSelectionChange(val) {
      this.environments = val
    },
    async onDelete() {
      const ids = this.environments.map(item => item.ID)
      const res = await deleteEnvironment({ ids })
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
        name: '',
        json: '',
        projectIp: ''
      }
    },
    closeDialog() {
      this.initForm()
      this.dialogFormVisible = false
    },
    openDialog(type) {
      switch (type) {
        case 'addEnvironment':
          this.dialogTitle = '新增Environment'
          break
        case 'edit':
          this.dialogTitle = '编辑Environment'
          break
        default:
          break
      }
      this.type = type
      this.dialogFormVisible = true
    },
    async editEnvironment(row) {
      const res = await getEnvironmentById({ id: row.ID })
      this.form = res.data.server
      this.openDialog('edit')
    },
    async deleteEnvironment(row) {
      this.$confirm('此操作将永久删除Environment?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })
        .then(async() => {
          const res = await deleteEnvironment(row)
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
            case 'addEnvironment':
              {
                const res = await addEnvironment(this.form)
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
                const res = await updateEnvironment(this.form)
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
    async getProjects() {
      const table = await getProjectList({ page: 1, pageSize: 10 })
      if (table.code === 0) {
        this.projects = table.data.list
        this.currentProject = this.projects[0]
        this.searchInfo = this.currentProject.ID
      }
    },
    async setCurrentProject(val) {
      this.currentProject = val
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
