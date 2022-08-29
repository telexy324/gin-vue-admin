<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="searchForm" :inline="true" :model="searchInfo">
        <el-form-item label="template name">
          <el-input v-model="searchInfo.name" placeholder="name" />
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
        <el-popover v-model:visible="deleteVisible" placement="top" width="160">
          <p>确定要删除吗？</p>
          <div style="text-align: right; margin-top: 8px;">
            <el-button size="small" type="text" @click="deleteVisible = false">取消</el-button>
            <el-button size="small" type="primary" @click="onDelete">确定</el-button>
          </div>
          <template #reference>
            <el-button icon="el-icon-delete" size="mini" :disabled="!templates.length" style="margin-left: 10px;">删除</el-button>
          </template>
        </el-popover>
        <el-dropdown  size="small" split-button type="primary">
          {{ currentProject.name }}
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item v-for="(item,index) in projects" :key="index" @click="setCurrentProject(item)">{{ item.name }}</el-dropdown-item>
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
        <el-table-column align="left" label="projectId" min-width="50" prop="projectId" sortable="custom" />
        <el-table-column align="left" label="inventoryId" min-width="50" prop="inventoryId" sortable="custom" />
        <el-table-column align="left" label="environmentId" min-width="50" prop="environmentId" sortable="custom" />
        <el-table-column align="left" label="playbook" min-width="150" prop="playbook" sortable="custom" />
        <el-table-column align="left" label="arguments" min-width="150" prop="arguments" sortable="custom" />
        <el-table-column align="left" label="description" min-width="150" prop="description" sortable="custom" />
        <el-table-column align="left" label="becomeKeyId" min-width="50" prop="becomeKeyId" sortable="custom" />
        <el-table-column align="left" label="vaultKeyId" min-width="50" prop="vaultKeyId" sortable="custom" />
        <el-table-column align="left" label="surveyVars" min-width="150" prop="surveyVars" sortable="custom" />
        <el-table-column align="left" label="vaultKeyId" min-width="150" prop="vaultKeyId" sortable="custom" />
        <el-table-column align="left" fixed="right" label="操作" width="200">
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
      <warning-bar title="新增template" />
      <el-form ref="templateForm" :model="form" :rules="rules" label-width="80px">
        <el-form-item label="name" prop="name">
          <el-input v-model="form.name" autocomplete="off" />
        </el-form-item>
        <el-form-item label="projectId" prop="projectId">
          <el-input v-model="form.projectId" :disabled="true" autocomplete="off" />
        </el-form-item>
        <el-form-item label="type" prop="type">
          <el-input v-model="form.type" autocomplete="off" />
        </el-form-item>
        <el-form-item label="template" prop="template">
          <el-input v-model="form.template" type="textarea" />
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
  getTemplateList,
  addTemplate,
  updateTemplate,
  deleteTemplate,
  getTemplateById
} from '@/api/ansibleTemplate'
import infoList from '@/mixins/infoList'
import { toSQLLine } from '@/utils/stringFun'
import warningBar from '@/components/warningBar/warningBar.vue'
import ansibleProjects from '@/mixins/ansibleProjects'

export default {
  name: 'Template',
  components: {
    warningBar
  },
  mixins: [infoList, ansibleProjects],
  data() {
    return {
      deleteVisible: false,
      listApi: getTemplateList,
      dialogFormVisible: false,
      dialogTitle: '新增template',
      templates: [],
      form: {
        name: '',
        projectId: '',
        type: '',
        template: ''
      },
      type: '',
      rules: {
        name: [{ required: true, message: '请输入template name', trigger: 'blur' }],
        projectId: [
          { required: true, message: '请输入project ip', trigger: 'blur' }
        ]
      },
      path: path
    }
  },
  async created() {
    await this.getProjects()
    await this.getTableData()
  },
  methods: {
    //  选中api
    handleSelectionChange(val) {
      this.templates = val
    },
    async onDelete() {
      const ids = this.templates.map(item => item.ID)
      const res = await deleteTemplate({ ids })
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
      this.searchInfo.projectId = this.currentProject.ID
      this.getTableData()
    },
    initForm() {
      this.$refs.templateForm.resetFields()
      this.form = {
        name: '',
        projectId: '',
        json: '',
        password: ''
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
          this.form.projectId = this.currentProject.ID
          break
        case 'edit':
          this.dialogTitle = '编辑Template'
          break
        default:
          break
      }
      this.type = type
      this.dialogFormVisible = true
    },
    async editTemplate(row) {
      const res = await getTemplateById({ id: row.ID, projectId: this.currentProject.ID })
      this.form = res.data.template
      this.openDialog('edit')
    },
    async deleteTemplate(row) {
      this.$confirm('此操作将永久删除Template?', '提示', {
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
            this.getTableData()
          }
        })
    },
    async enterDialog() {
      this.$refs.templateForm.validate(async valid => {
        if (valid) {
          switch (this.type) {
            case 'addTemplate':
              {
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
.gva-btn-list :deep(.el-dropdown) {
  float: right;
  height: 32px;
}
</style>
