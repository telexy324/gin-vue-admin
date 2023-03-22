<template>
  <div>
    <div class="gva-table-box">
      <el-table :data="tableData" @sort-change="sortChange" @selection-change="handleSelectionChange">
        <el-table-column
          type="selection"
          width="55"
        />
        <el-table-column align="left" label="id" min-width="60" sortable="custom">
          <template v-slot="scope">
            <el-button
              type="text"
              link
              @click="showTaskLog(scope.row)"
            >{{ scope.row.ID }}</el-button>
<!--            <a @click="showTaskLog(scope.row)">{{ scope.row.ID }}</a>-->
          </template>
        </el-table-column>
        <el-table-column align="left" label="模板id" min-width="150" prop="templateId" sortable="custom" />
        <el-table-column align="left" label="状态" min-width="150" sortable="custom">
          <template v-slot="scope">
            <TaskStatus :status="scope.row.status" />
          </template>
        </el-table-column>
        <el-table-column align="left" label="创建人" min-width="200" prop="userId" sortable="custom" />
        <el-table-column align="left" label="开始时间" min-width="150" prop="beginTime.Time" sortable="custom" :formatter="dateFormatter1" />
        <el-table-column align="left" label="结束时间" min-width="150" prop="endTime.Time" sortable="custom" :formatter="dateFormatter2" />
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

const path = import.meta.env.VITE_BASE_API
// 获取列表内容封装在mixins内部  getTableData方法 初始化已封装完成 条件搜索时候 请把条件安好后台定制的结构体字段 放到 this.searchInfo 中即可实现条件搜索

import {
  getTaskList,
  deleteTask,
} from '@/api/task'
import infoList from '@/mixins/infoList'
import { toSQLLine } from '@/utils/stringFun'
import { emitter } from '@/utils/bus'
import TaskStatus from '@/components/task/TaskStatus.vue'
import { formatTimeToStr } from '@/utils/date'

export default {
  name: 'TaskList',
  components: {
    TaskStatus
  },
  mixins: [infoList],
  data() {
    return {
      listApi: getTaskList,
      type: '',
      path: path,
      tasks: [],
    }
  },
  created() {
    this.getTableData()
  },
  methods: {
    async onDelete() {
      const ids = this.tasks.map(item => item.ID)
      const res = await deleteTask({ ids })
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
    showTaskLog(t) {
      emitter.emit('i-show-task', t)
    },
    dateFormatter1(row) {
      if (row.beginTime.Time !== null && row.beginTime.Time !== '') {
        var date = new Date(row.beginTime.Time)
        return formatTimeToStr(date, 'yyyy-MM-dd hh:mm:ss')
      } else {
        return ''
      }
    },
    dateFormatter2(row) {
      if (row.endTime.Time !== null && row.endTime.Time !== '') {
        var date = new Date(row.endTime.Time)
        return formatTimeToStr(date, 'yyyy-MM-dd hh:mm:ss')
      } else {
        return ''
      }
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

.excel-btn + .excel-btn {
  margin-left: 10px;
}
</style>
