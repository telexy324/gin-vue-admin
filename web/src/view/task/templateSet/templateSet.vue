<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="searchForm" :inline="true" :model="searchInfo">
        <el-form-item label="模版集名">
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
        <el-button size="mini" type="primary" icon="el-icon-plus" @click="openDialog('addSet')">新增</el-button>
        <el-popover v-model:visible="deleteVisible" placement="top" width="160">
          <p>确定要删除吗？</p>
          <div style="text-align: right; margin-top: 8px;">
            <el-button size="mini" type="text" @click="deleteVisible = false">取消</el-button>
            <el-button size="mini" type="primary" @click="onDelete">确定</el-button>
          </div>
          <template #reference>
            <el-button icon="el-icon-delete" size="mini" :disabled="!sets.length" style="margin-left: 10px;">删除</el-button>
          </template>
        </el-popover>
        <el-button class="excel-btn" size="mini" type="success" icon="el-icon-download" @click="openDrawer()">选择系统</el-button>
      </div>
      <el-table :data="tableData" @sort-change="sortChange" @selection-change="handleSelectionChange">
        <el-table-column type="expand">
          <template #default="scope">
            <TemplateSetTask :set-id="scope.row.ID" />
          </template>
        </el-table-column>
        <el-table-column
          type="selection"
          width="55"
        />
        <el-table-column align="left" label="name" min-width="60" prop="name" sortable="custom" />
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
              @click="editSet(scope.row)"
            >编辑</el-button>
            <el-button
              icon="el-icon-delete"
              size="small"
              type="text"
              @click="deleteSet(scope.row)"
            >删除</el-button>
            <el-button
              icon="el-icon-caret-right"
              size="small"
              type="text"
              @click="createSetTask(scope.row)"
            >新增</el-button>
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
      <warning-bar title="新增模板集" />
      <el-form ref="setForm" :model="form" :rules="rules" label-width="80px">
        <el-row>
          <el-col :span="12">
            <el-form-item label="模版集名" prop="name">
              <el-input v-model="form.name" autocomplete="off" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="所属系统" prop="systemId">
          <el-select v-model="form.systemId" @change="setTemplate">
            <el-option v-for="val in systemOptions" :key="val.ID" :value="val.ID" :label="val.name" />
          </el-select>
        </el-form-item>
        <div v-for="(item, index) in form.templates" :key="index">
          <el-row>
            <el-col :span="10">
              <el-form-item
                label="模板名"
                :prop="'templates.' + index + '.ID'"
                :rules="{
                  required: true, message: '请选择模板', trigger: 'blur'
                }"
              >
                <el-select v-model="item.templateId">
                  <el-option v-for="val in systemTemplateOptions" :key="val.ID" :value="val.ID" :label="val.name" />
                </el-select>
              </el-form-item>
            </el-col>
            <el-col :span="10">
              <el-form-item
                label="序号"
                :prop="'templates.' + index + '.seq'"
              >
                <el-input v-model.number="item.seq" />
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
          <el-button size="small" type="warning" :disabled="setDisabled" @click="addItem">增加</el-button>
          <el-button size="small" @click="closeDialog">取 消</el-button>
          <el-button size="small" type="primary" @click="enterDialog">确 定</el-button>
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
  getSetList,
  deleteSet,
  getSetById,
  addSet,
  updateSet,
  deleteSetByIds,
  addSetTask,
} from '@/api/template'
import { addTask } from '@/api/task'
import { getAllServerIds } from '@/api/cmdb'
import infoList from '@/mixins/infoList'
import { toSQLLine } from '@/utils/stringFun'
import warningBar from '@/components/warningBar/warningBar.vue'
import { emitter } from '@/utils/bus'
import Systems from '@/components/task/systems.vue'
import TemplateSetTask from '@/view/task/templateSet/components/templateSetTask.vue'

export default {
  name: 'TemplateSet',
  components: {
    warningBar,
    Systems,
    TemplateSetTask
  },
  mixins: [infoList],
  data() {
    return {
      deleteVisible: false,
      listApi: getSetList,
      dialogFormVisible: false,
      dialogTitle: '新增模板集',
      sets: [],
      form: {
        ID: 0,
        name: '',
        systemId: '',
        templates: [],
      },
      type: '',
      rules: {
        name: [{ required: true, message: '请输入模板名', trigger: 'blur' }],
        systemId: [
          { required: true, message: '请选择系统', trigger: 'blur' }
        ],
      },
      path: path,
      drawer: false,
      systemOptions: [],
      systemTemplateOptions: [],
      setDisabled: true,
    }
  },
  // watch: {
  //   'form.systemId': function(newVal, oldval) {
  //     this.setDisabled = !newVal
  //     this.form.templates = []
  //     this.searchInfo.systemIds = []
  //     this.searchInfo.systemIds.push(Number(newVal))
  //     this.setTemplateOptions()
  //   },
  // },
  async created() {
    // socket.addListener((data) => this.onWebsocketDataReceived(data))
    if (this.$route.params.systemIds) {
      this.searchInfo.systemIds = this.formRouterParam(this.$route.params.systemIds)
    }
    await this.getTableData()
    const res = await getAllServerIds()
    this.setSystemOptions(res.data)
    await this.setTemplateOptions()
  },
  methods: {
    //  选中api
    handleSelectionChange(val) {
      this.sets = val
    },
    async onDelete() {
      const ids = this.sets.map(item => item.ID)
      const res = await deleteSetByIds({ ids })
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
      this.$refs.setForm.resetFields()
      this.form = {
        ID: 0,
        name: '',
        systemId: '',
        templates: [],
      }
    },
    closeDialog() {
      this.initForm()
      this.dialogFormVisible = false
    },
    openDialog(type) {
      switch (type) {
        case 'addSet':
          this.dialogTitle = '新增Set'
          break
        case 'edit':
          this.dialogTitle = '编辑Set'
          break
        default:
          break
      }
      this.type = type
      this.dialogFormVisible = true
    },
    async editSet(row) {
      const res = await getSetById({ id: row.ID })
      this.form = res.data
      this.openDialog('edit')
    },
    async deleteSet(row) {
      this.$confirm('此操作将永久删除模板集?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })
        .then(async() => {
          const res = await deleteSet(row)
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
      this.$refs.setForm.validate(async valid => {
        if (valid) {
          switch (this.type) {
            case 'addSet':
              {
                const res = await addSet(this.form)
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
                const res = await updateSet(this.form)
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
    async runTask(row) {
      const task = (await addTask({
        templateId: row.ID
      })).data.task
      this.showTaskLog(task)
    },
    showTaskLog(task) {
      emitter.emit('i-show-task', task)
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
    filterSystemName(value) {
      const rowLabel = this.systemOptions.filter(item => item.ID === value)
      return rowLabel && rowLabel[0] && rowLabel[0].name
    },
    addItem() {
      this.form.templates.push({
        setId: this.form.ID ? this.form.ID : 0,
        templateId: '',
        seq: 99
      })
    },
    deleteItem(item, index) {
      this.form.templates.splice(index, 1)
    },
    filterTemplateName(value) {
      const rowLabel = this.systemTemplateOptions.filter(item => item.ID === value)
      return rowLabel && rowLabel[0] && rowLabel[0].name
    },
    async setTemplateOptions() {
      this.searchInfo.executeType = 1
      const res = await getTemplateList({ page: 1, pageSize: 99999, ...this.searchInfo })
      this.systemTemplateOptions = res.data.list
      console.log(this.systemTemplateOptions)
    },
    setSystemOptions(data) {
      this.systemOptions = data
    },
    // showDetail(row) {
    //   this.$router.push({
    //     name: 'templateSetDetail',
    //     params: {
    //       setId: row.ID,
    //     }
    //   })
    // },
    async createSetTask(row) {
      await addSetTask({ setId: row.ID })
      await this.getTableData()
    },
    setTemplate(val) {
      this.setDisabled = !val
      this.form.templates = []
      this.searchInfo.systemIds = []
      this.searchInfo.systemIds.push(Number(val))
      this.setTemplateOptions()
    }
  }
}
</script>

<style lang="scss">
.el-table__expand-icon {
  height: 35px;
}
</style>

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

.el-row {
  padding: 0 0;
}
</style>
