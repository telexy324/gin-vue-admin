<template>
  <div>
<!--    <div class="gva-search-box">-->
<!--      <el-form ref="searchForm" :inline="true" :model="searchInfo">-->
<!--        <el-form-item label="templateId">-->
<!--          <el-input v-model="searchInfo.templateId" placeholder="模板ID" />-->
<!--        </el-form-item>-->
<!--        <el-form-item>-->
<!--          <el-button size="mini" type="primary" icon="el-icon-search" @click="onSubmit">查询</el-button>-->
<!--          <el-button size="mini" icon="el-icon-refresh" @click="onReset">重置</el-button>-->
<!--        </el-form-item>-->
<!--      </el-form>-->
<!--    </div>-->
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button size="mini" type="primary" icon="el-icon-plus" :disabled="!hasCreate" @click="openDialog('addSchedule')">新增</el-button>
        <el-popover v-model:visible="deleteVisible" placement="top" width="160">
          <p>确定要删除吗？</p>
          <div style="text-align: right; margin-top: 8px;">
            <el-button size="mini" type="text" @click="deleteVisible = false">取消</el-button>
            <el-button size="mini" type="primary" @click="onDelete">确定</el-button>
          </div>
          <template #reference>
            <el-button icon="el-icon-delete" size="mini" :disabled="!schedules.length || !hasDelete" style="margin-left: 10px;">删除</el-button>
          </template>
        </el-popover>
        <el-button class="excel-btn" size="mini" type="success" icon="el-icon-download" @click="openDrawer()">选择系统</el-button>
      </div>
      <el-table :data="tableData" @sort-change="sortChange" @selection-change="handleSelectionChange">
        <el-table-column
          type="selection"
          width="55"
        />
<!--        <el-table-column align="left" label="id" min-width="60" prop="ID" sortable="custom" />-->
        <el-table-column align="left" label="模版名" min-width="150" prop="templateId" sortable="custom">
          <template #default="scope">
            <div>{{ filterTemplateName(scope.row.templateId) }}</div>
          </template>
        </el-table-column>
        <el-table-column align="left" label="启用" min-width="150" prop="valid" sortable="custom">
          <template #default="scope">
            <el-switch
              v-model="scope.row.valid"
              :active-value="1"
              :inactive-value="0"
              :disabled="!hasEdit"
              @change="changeValid(scope.row)"
            />
          </template>
        </el-table-column>
        <el-table-column align="left" label="cronFormat" min-width="150" prop="cronFormat" sortable="custom" />
        <el-table-column align="left" fixed="right" label="操作" width="200">
          <template #default="scope">
            <el-button
              icon="el-icon-edit"
              size="small"
              type="text"
              :disabled="!hasEdit"
              @click="editSchedule(scope.row)"
            >编辑</el-button>
            <el-button
              icon="el-icon-delete"
              size="small"
              type="text"
              :disabled="!hasDelete"
              @click="deleteSchedule(scope.row)"
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
      <warning-bar title="新增定时任务，可预先检查时间格式，类似crontab" />
      <el-form ref="scheduleForm" :model="form" :rules="rules" label-width="120px">
<!--        <el-form-item label="模版名" prop="templateId">-->
<!--          <el-select-->
<!--            v-model="form.templateId"-->
<!--            placeholder="Select"-->
<!--            style="width: 240px"-->
<!--          >-->
<!--            <el-option-->
<!--              v-for="item in templateOptions"-->
<!--              :key="item.ID"-->
<!--              :label="item.name"-->
<!--              :value="item.ID"-->
<!--            />-->
<!--          </el-select>-->
<!--        </el-form-item>-->
        <el-row>
          <el-col :span="12">
            <el-form-item label="所属系统" prop="systemId">
              <el-select v-model="form.systemId" @change="changeSystemId">
                <el-option v-for="val in systemOptions" :key="val.ID" :value="val.ID" :label="val.name" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="模版名" prop="templateId">
              <el-select
                v-model="form.templateId"
                placeholder="Select"
                @change="changeTemplateId"
              >
                <el-option
                  v-for="item in templateTempOptions"
                  :key="item.ID"
                  :label="item.name"
                  :value="item.ID"
                />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="cronFormat" prop="cronFormat">
          <el-input v-model="form.cronFormat" autocomplete="off" />
        </el-form-item>
        <el-row>
          <el-col :span="12">
            <el-form-item>
              <el-button size="small" type="primary" @click="checkCronFormat">检查Schedule格式</el-button>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="启用">
              <el-switch
                v-model="form.valid"
                :active-value="1"
                :inactive-value="0"
              />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="目标" prop="targetIds">
          <el-cascader
            v-model="form.targetIds"
            style="width:100%"
            :options="checkedServerOptions"
            :show-all-levels="false"
            :props="{ multiple:true,checkStrictly: false,label:'name',value:'ID',disabled:'disabled',emitPath:false}"
            :clearable="true"
          />
        </el-form-item>
        <div v-for="(item, index) in form.commandVars" :key="index">
          <el-form-item
            :label="'参数' + index"
            :prop="'commandVars.' + index"
            :rules="rules.commandVars"
          >
            <el-input v-model="form.commandVars[index]" />
          </el-form-item>
        </div>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button size="small" @click="closeDialog">取 消</el-button>
          <el-button size="small" type="primary" @click="enterDialog">确 定</el-button>
        </div>
      </template>
    </el-dialog>
    <el-drawer v-if="drawer" v-model="drawer" :with-header="false" size="40%" title="请选择系统">
      <Systems ref="systems" :keys="searchInfo.systemIds" @checked="getCheckedSchedules" />
    </el-drawer>
  </div>
</template>

<script>

const path = import.meta.env.VITE_BASE_API
// 获取列表内容封装在mixins内部  getTableData方法 初始化已封装完成 条件搜索时候 请把条件安好后台定制的结构体字段 放到 this.searchInfo 中即可实现条件搜索

import {
  getTemplateScheduleList,
  addSchedule,
  updateSchedule,
  deleteSchedule,
  getScheduleById,
  validateScheduleCronFormat,
  deleteScheduleByIds
} from '@/api/schedule'
import {
  getTemplateList
} from '@/api/template'
import { getAdminSystems, getSystemServerIds } from '@/api/cmdb'
import { getPolicyPathByAuthorityId } from '@/api/casbin'
import infoList from '@/mixins/infoList'
import { toSQLLine } from '@/utils/stringFun'
import warningBar from '@/components/warningBar/warningBar.vue'
import { mapGetters } from 'vuex'
import Systems from '@/components/task/systems.vue'

export default {
  name: 'Schedule',
  components: {
    warningBar,
    Systems
  },
  mixins: [infoList],
  data() {
    return {
      deleteVisible: false,
      listApi: getTemplateScheduleList,
      dialogFormVisible: false,
      dialogTitle: '新增schedule',
      schedules: [],
      form: {
        templateId: '',
        cronFormat: '',
        valid: 0,
        systemId: '',
        commandVars: [],
        targetIds: [],
      },
      type: '',
      rules: {
        templateId: [{ required: true, message: '请关联模板', trigger: 'blur' }],
        cronFormat: [{ required: true, message: '请输入定时任务', trigger: 'blur' }],
        commandVars: [
          { required: true, message: '请输入任务参数', trigger: 'blur' },
        ],
      },
      path: path,
      templateOptions: [],
      hasEdit: true,
      hasCreate: true,
      hasDelete: true,
      drawer: false,
      systemOptions: [],
      templateTempOptions: [],
      checkedServerOptions: [],
    }
  },
  computed: {
    ...mapGetters('user', ['userInfo'])
  },
  async created() {
    await this.getTableData()
    await this.setSystemOptions()
    await this.setAllTemplateOptions()
  },
  mounted() {
    this.authorities()
  },
  methods: {
    //  选中api
    handleSelectionChange(val) {
      this.schedules = val
    },
    async onDelete() {
      const ids = this.schedules.map(item => item.ID)
      const res = await deleteScheduleByIds({ ids })
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
      this.$refs.scheduleForm.resetFields()
      this.form = {
        templateId: '',
        cronFormat: '',
        valid: 0,
        systemId: '',
        commandVars: [],
      }
    },
    closeDialog() {
      this.initForm()
      this.dialogFormVisible = false
      this.templateTempOptions = []
    },
    openDialog(type) {
      switch (type) {
        case 'addSchedule':
          this.dialogTitle = '新增定时任务'
          break
        case 'edit':
          this.dialogTitle = '编辑定时任务'
          break
        default:
          break
      }
      this.type = type
      this.dialogFormVisible = true
    },
    async editSchedule(row) {
      const res = await getScheduleById({ id: row.ID })
      this.form = res.data.schedule
      this.form.systemId = res.data.systemId
      await this.setTemplateOptions(this.form.systemId)
      this.openDialog('edit')
    },
    async deleteSchedule(row) {
      this.$confirm('此操作将永久删除定时任务?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })
        .then(async() => {
          const res = await deleteSchedule(row)
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
      await validateScheduleCronFormat(this.form)
      this.$refs.scheduleForm.validate(async valid => {
        if (valid) {
          switch (this.type) {
            case 'addSchedule':
              {
                const res = await addSchedule(this.form)
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
                const res = await updateSchedule(this.form)
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
    setOptions(data) {
      this.templateOptions = data
    },
    async checkCronFormat() {
      const _form = Object.assign({}, this.form)
      _form.templateId = _form.templateId ? _form.templateId : 0
      const res = await validateScheduleCronFormat(_form)
      if (res.code === 0) {
        this.$message({
          type: 'success',
          message: res.msg,
          showClose: true
        })
      }
    },
    async changeValid(row) {
      await updateSchedule(row)
      await this.getTableData()
    },
    filterTemplateName(value) {
      const rowLabel = this.templateOptions.filter(item => item.ID === value)
      return rowLabel && rowLabel[0] && rowLabel[0].name
    },
    async authorities() {
      const res = await getPolicyPathByAuthorityId({
        authorityId: this.userInfo.authorityId
      })
      this.hasEdit = !!res.data.paths.some((item) => {
        return item.path === '/task/schedule/updateSchedule'
      })
      this.hasCreate = !!res.data.paths.some((item) => {
        return item.path === '/task/schedule/addSchedule'
      })
      this.hasDelete = !!res.data.paths.some((item) => {
        return item.path === '/task/schedule/deleteSchedule'
      })
    },
    openDrawer() {
      this.drawer = true
    },
    getCheckedSchedules(checkArr) {
      const systemIDs = []
      checkArr.forEach(item => {
        systemIDs.push(Number(item.ID))
      })
      this.searchInfo.systemIds = systemIDs
      this.getTableData()
      this.drawer = false
    },
    async setSystemOptions() {
      const res = await getAdminSystems()
      this.systemOptions = res.data.systems
    },
    async setAllTemplateOptions() {
      const res = await getTemplateList({
        page: 1,
        pageSize: 99999,
      })
      this.setOptions(res.data.list)
    },
    changeSystemId(selectValue) {
      this.setTemplateOptions(selectValue)
    },
    async setTemplateOptions(selectValue) {
      const systemIDs = []
      systemIDs.push(selectValue)
      const res = await getTemplateList({
        page: 1,
        pageSize: 99999,
        systemIds: systemIDs,
      })
      this.templateTempOptions = res.data.list
    },
    changeTemplateId(selectValue) {
      this.form.commandVars = []
      for (let i = 0; i < this.templateTempOptions.find(item => item.ID === selectValue).commandVarNumbers; i++) {
        this.form.commandVars.push('')
      }
    },
    async setCheckedServerOptions(template) {
      const res = await getSystemServerIds({
        ID: template.systemId
      })
      const serverOptions = res.data
      serverOptions[0].children = serverOptions[0].children.filter((item) => {
        if (template.targetServerIds.includes(item.ID)) {
          if (template.executeType !== 2) {
            this.commandVarForm.targetIds.push(item.ID)
          }
          return true
        }
        return false
      })
      this.checkedServerOptions = serverOptions
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
