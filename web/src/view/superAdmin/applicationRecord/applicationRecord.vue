<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo">
        <el-form-item label="请求路径">
          <el-input v-model="searchInfo.action" placeholder="搜索条件" />
        </el-form-item>
        <el-form-item>
          <el-button size="mini" type="primary" icon="el-icon-search" @click="onSubmit">查询</el-button>
          <el-button size="mini" icon="el-icon-refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
      <div class="gva-btn-list">

        <el-popover v-model:visible="deleteVisible" placement="top" width="160">
          <p>确定要删除吗？</p>
          <div style="text-align: right; margin-top: 8px;">
            <el-button size="mini" type="text" @click="deleteVisible = false">取消</el-button>
            <el-button size="mini" type="primary" @click="onDelete">确定</el-button>
          </div>
          <template #reference>
            <el-button icon="el-icon-delete" size="mini" style="margin-left: 10px;" :disabled="!multipleSelection.length">删除</el-button>
          </template>
        </el-popover>
        <el-button size="mini" style="margin-left: 10px;" type="success" icon="el-icon-download" @click="downLoadFile">导出</el-button>
      </div>
      <el-table
        ref="multipleTable"
        :data="tableData"
        style="width: 100%"
        tooltip-effect="dark"
        row-key="ID"
        @selection-change="handleSelectionChange"
      >
        <el-table-column align="left" type="selection" width="55" />
        <el-table-column align="left" label="操作人" width="140">
          <template #default="scope">
            <div>{{ scope.row.user.userName }}({{ scope.row.user.nickName }})</div>
          </template>
        </el-table-column>
        <el-table-column align="left" label="日期" width="180">
          <template #default="scope">{{ formatDate(scope.row.logTime) }}</template>
        </el-table-column>
        <el-table-column align="left" label="状态码" prop="status" width="120">
          <template #default="scope">
            <div>
              <el-tag :type="getStatusColor(scope.row.status)">{{ getStatusMessage(scope.row.status) }}</el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column align="left" label="请求IP" prop="ip" width="120" />
        <el-table-column align="left" label="请求路径" prop="action" width="240" />
        <el-table-column align="left" label="详情" prop="detail" width="360" />
        <el-table-column align="left" label="错误信息" prop="errorMessage" width="360" />
        <el-table-column align="left" label="按钮组">
          <template #default="scope">
            <el-popover :visible="scope.row.visible" placement="top" width="160">
              <p>确定要删除吗？</p>
              <div style="text-align: right; margin-top: 8px;">
                <el-button size="mini" type="text" @click="scope.row.visible = false">取消</el-button>
                <el-button size="mini" type="primary" @click="deleteApplicationRecord(scope.row)">确定</el-button>
              </div>
              <template #reference>
                <el-button icon="el-icon-delete" size="mini" type="text">删除</el-button>
              </template>
            </el-popover>
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
  </div>
</template>

<script>
import {
  deleteApplicationRecord,
  getApplicationRecordList,
  deleteApplicationRecordByIds,
  exportApplicationRecord
} from '@/api/applicationRecord' // 此处请自行替换地址
import infoList from '@/mixins/infoList'

export default {
  name: 'ApplicationRecord',
  mixins: [infoList],
  data() {
    return {
      listApi: getApplicationRecordList,
      dialogFormVisible: false,
      type: '',
      deleteVisible: false,
      multipleSelection: [],
      formData: {
        ip: null,
        detail: null,
        action: null,
        status: null,
        errorMessage: null,
        userId: null,
        logTime: null,
      }
    }
  },
  created() {
    this.getTableData()
  },
  methods: {
    onReset() {
      this.searchInfo = {}
    },
    // 条件搜索前端看此方法
    onSubmit() {
      this.page = 1
      this.pageSize = 10
      this.getTableData()
    },
    handleSelectionChange(val) {
      this.multipleSelection = val
    },
    async onDelete() {
      const ids = []
      this.multipleSelection &&
        this.multipleSelection.forEach(item => {
          ids.push(item.ID)
        })
      const res = await deleteApplicationRecordByIds({ ids })
      if (res.code === 0) {
        this.$message({
          type: 'success',
          message: '删除成功'
        })
        if (this.tableData.length === ids.length && this.page > 1) {
          this.page--
        }
        this.deleteVisible = false
        this.getTableData()
      }
    },
    async deleteApplicationRecord(row) {
      row.visible = false
      const res = await deleteApplicationRecord({ ID: row.ID })
      if (res.code === 0) {
        this.$message({
          type: 'success',
          message: '删除成功'
        })
        if (this.tableData.length === 1 && this.page > 1) {
          this.page--
        }
        this.getTableData()
      }
    },
    getStatusColor(status) {
      switch (status) {
        case 0:
          return 'success'
        case 7:
          return 'danger'
        default:
          throw new Error(`Unknown code status ${status}`)
      }
    },
    getStatusMessage(status) {
      switch (status) {
        case 0:
          return '成功'
        case 7:
          return '失败'
        default:
          throw new Error(`Unknown code status ${status}`)
      }
    },
    downLoadFile() {
      const ids = []
      this.multipleSelection &&
      this.multipleSelection.forEach(item => {
        ids.push(item.ID)
      })
      const fileName = 'logRecord' + Math.round(new Date()) + '.xlsx'
      exportApplicationRecord({ ids }, fileName)
    },
  }
}
</script>

<style lang="scss">
.table-expand {
  padding-left: 60px;
  font-size: 0;
  label {
    width: 90px;
    color: #99a9bf;
    .el-form-item {
      margin-right: 0;
      margin-bottom: 0;
      width: 50%;
    }
  }
}
</style>
