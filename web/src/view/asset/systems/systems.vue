<template>
  <div>
    <div class="gva-search-box">
      <el-form ref="searchForm" :inline="true" :model="searchInfo">
        <el-form-item label="系统名">
          <el-input v-model="searchInfo.name" placeholder="系统名" />
        </el-form-item>
        <el-form-item>
          <el-button size="mini" type="primary" icon="el-icon-search" @click="onSubmit">查询</el-button>
          <el-button size="mini" icon="el-icon-refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button size="mini" type="primary" icon="el-icon-plus" @click="openDialog('addSystem')">新增</el-button>
        <el-popover v-model:visible="deleteVisible" placement="top" width="160">
          <p>确定要删除吗？</p>
          <div style="text-align: right; margin-top: 8px;">
            <el-button size="mini" type="text" @click="deleteVisible = false">取消</el-button>
            <el-button size="mini" type="primary" @click="onDelete">确定</el-button>
          </div>
          <template #reference>
            <el-button icon="el-icon-delete" size="mini" :disabled="!systems.length" style="margin-left: 10px;">删除</el-button>
          </template>
        </el-popover>
      </div>
      <el-table :data="tableData" @sort-change="sortChange" @selection-change="handleSelectionChange">
        <el-table-column
          type="selection"
          width="55"
        />
<!--        <el-table-column align="left" label="id" min-width="60" prop="system.ID" sortable="custom" />-->
        <el-table-column align="left" label="系统名" min-width="150" prop="system.name" sortable="custom" />
        <el-table-column align="left" fixed="right" label="操作" width="200">
          <template #default="scope">
            <el-button
              icon="el-icon-edit"
              size="small"
              type="text"
              @click="editSystem(scope.row)"
            >编辑</el-button>
            <el-button
              icon="el-icon-delete"
              size="small"
              type="text"
              @click="deleteSystem(scope.row)"
            >删除</el-button>
            <el-button
              icon="el-icon-orange"
              size="small"
              type="text"
              @click="showRelation(scope.row)"
            >关系图</el-button>
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
      <warning-bar title="新增系统" />
      <el-form ref="systemForm" :model="form" :rules="rules" label-width="80px">
        <el-form-item label="系统名" prop="system.name">
          <el-input v-model="form.system.name" autocomplete="off" />
        </el-form-item>
        <el-form-item label="管理员" prop="adminIds">
          <el-select
            v-model="form.adminIds"
            multiple
            placeholder="Select"
            style="width: 240px"
          >
            <el-option
              v-for="item in userOptions"
              :key="item.ID"
              :label="item.userName"
              :value="item.ID"
            />
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
  getSystemList, addSystem, updateSystem, deleteSystem, getSystemById, deleteSystemByIds
} from '@/api/cmdb'
import {
  getUserList
} from '@/api/user'
import infoList from '@/mixins/infoList'
import { toSQLLine } from '@/utils/stringFun'
import warningBar from '@/components/warningBar/warningBar.vue'

export default {
  name: 'Systems',
  components: {
    warningBar
  },
  mixins: [infoList],
  data() {
    return {
      deleteVisible: false,
      listApi: getSystemList,
      dialogFormVisible: false,
      dialogTitle: '新增system',
      systems: [],
      form: {
        system: '',
        adminIds: ''
      },
      type: '',
      rules: {
        name: [{ required: true, message: '请输入系统名', trigger: 'blur' }],
        adminIds: [
          { required: true, message: '请选择管理员', trigger: 'blur' }
        ]
      },
      path: path,
      userOptions: []
    }
  },
  async created() {
    await this.getTableData()
    const res = await getUserList({
      page: 1,
      pageSize: 99999
    })
    this.setOptions(res.data.list)
  },
  methods: {
    //  选中api
    handleSelectionChange(val) {
      this.systems = val
    },
    async onDelete() {
      const ids = this.systems.map(item => item.ID)
      const res = await deleteSystemByIds({ ids })
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
      this.$refs.systemForm.resetFields()
      this.form = {
        system: '',
        adminIds: ''
      }
    },
    closeDialog() {
      this.initForm()
      this.dialogFormVisible = false
    },
    openDialog(type) {
      switch (type) {
        case 'addSystem':
          this.dialogTitle = '新增System'
          break
        case 'edit':
          this.dialogTitle = '编辑System'
          break
        default:
          break
      }
      this.type = type
      this.dialogFormVisible = true
    },
    async editSystem(row) {
      const res = await getSystemById({ id: row.system.ID })
      this.form = res.data
      this.form.adminIds = res.data.adminIds
      this.openDialog('edit')
    },
    async deleteSystem(row) {
      this.$confirm('此操作将永久删除系统?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })
        .then(async() => {
          const res = await deleteSystem(row)
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
      this.$refs.systemForm.validate(async valid => {
        if (valid) {
          switch (this.type) {
            case 'addSystem':
              {
                const res = await addSystem(this.form)
                if (res.code === 0) {
                  this.$message({
                    type: 'success',
                    message: '添加成功',
                    showClose: true
                  })
                }
                await this.getTableData()
                this.closeDialog()
              }

              break
            case 'edit':
              {
                const res = await updateSystem(this.form)
                if (res.code === 0) {
                  this.$message({
                    type: 'success',
                    message: '编辑成功',
                    showClose: true
                  })
                }
                await this.getTableData()
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
      this.userOptions = data
    },
    showRelation(row) {
      this.$router.push({
        name: 'antv',
        params: {
          systemId: row.system.ID,
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
