<template>
  <div>
    <div class="gva-table-box">
      <el-table :data="tableData" @sort-change="sortChange">
        <el-table-column align="left" label="id" min-width="60" sortable="custom">
          <template v-slot="scope">
            <el-button
              type="text"
              link
              @click="showDetail(scope.row)"
            >#{{ scope.row.ID }}</el-button>
          </template>
        </el-table-column>
        <el-table-column align="left" label="状态" min-width="150" sortable="custom">
          <template v-slot="scope">
            <TaskStatus :status="getCurrentStatus(scope.row)" />
          </template>
        </el-table-column>
        <el-table-column align="left" label="创建人" min-width="200" prop="systemUserId" sortable="custom">
          <template #default="scope">
            <div>{{ filterUserName(scope.row.systemUserId) }}</div>
          </template>
        </el-table-column>
        <el-table-column align="left" label="当前任务ID" min-width="150" sortable="custom">
          <template v-slot="scope">
            <el-button
              type="text"
              link
              @click="showTaskLog(scope.row)"
            >#{{ scope.row.currentTaskId }}</el-button>
            <!--            <a @click="showTaskLog(scope.row)">{{ scope.row.ID }}</a>-->
          </template>
        </el-table-column>
        <el-table-column align="left" label="当前步骤" min-width="150" prop="currentStep" sortable="custom" />
        <el-table-column align="left" label="总步骤" min-width="150" prop="totalSteps" sortable="custom" />
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
  getSetTaskList
} from '@/api/template'
import {
  getUserList
} from '@/api/user'
import infoList from '@/mixins/infoList'
import { toSQLLine } from '@/utils/stringFun'
import { emitter } from '@/utils/bus'
import TaskStatus from '@/components/task/TaskStatus.vue'

export default {
  name: 'TemplateSetTask',
  components: {
    TaskStatus
  },
  props: {
    setId: {
      type: Number,
      default: function() {
        return 0
      },
    }
  },
  mixins: [infoList],
  data() {
    return {
      listApi: getSetTaskList,
      path: path,
      userOptions: []
    }
  },
  async created() {
    this.searchInfo.setId = this.setId
    await this.getTableData()
    const resUser = await getUserList({
      page: 1,
      pageSize: 99999
    })
    this.setUserOptions(resUser.data.list)
  },
  methods: {
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
    setUserOptions(data) {
      this.userOptions = data
    },
    filterUserName(value) {
      if (value === 999999) {
        return '定时任务'
      }
      const rowLabel = this.userOptions.filter(item => item.ID === value)
      return rowLabel && rowLabel[0] && rowLabel[0].userName
    },
    showDetail(row) {
      this.$router.push({
        name: 'templateSetDetail',
        params: {
          setTaskId: row.ID,
        }
      })
    },
    getCurrentStatus(row) {
      if (row.tasks) {
        return row.tasks.find((item) => {
          return row.currentTaskId === item.ID
        }).status
      } else {
        return
      }
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

.excel-btn + .excel-btn {
  margin-left: 10px;
}
</style>
