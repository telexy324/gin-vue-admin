<template>
  <div class="task-log-view">
    <el-dialog v-model="visible" :show-close="false">
      <template #header>
        <div class="el-dialog__header">
          <el-row :gutter="10">
            <el-col :span="6">
              <TaskStatus :status="item.status" />
            </el-col>
            <el-col :span="3">Author:</el-col>
            <el-col :span="3" v-text="user.name" />
            <el-col :span="3">start:</el-col>
            <el-col :span="3" v-text="item.start_time" />
            <el-col :span="3">end:</el-col>
            <el-col :span="3" v-text="item.end_time" />
          </el-row>
        </div>
      </template>
      <div ref="output" class="task-log-records">
        <div v-for="record in output" :key="record.ID" class="task-log-records__record">
          <div class="task-log-records__time">
            {{ record.recordTime }}
          </div>
          <div class="task-log-records__output">{{ record.output }}</div>
        </div>
      </div>
      <el-button
        v-if="item.status === 'running' || item.status === 'waiting'"
        type="danger"
        round
        style="position: absolute; bottom: 10px; right: 10px;"
        @click="stopTask()"
      >
        Stop
      </el-button>
    </el-dialog>
  </div>
</template>

<style lang="scss">

// @import '~vuetify/src/styles/settings/_variables';

.task-log-view {
}

.task-log-records {
  background: black;
  color: white;
  height: calc(100vh - 250px);
  overflow: auto;
  font-family: monospace;
  margin: 0 -24px;
  padding: 5px 10px;
}

//.task-log-view--with-message .task-log-records {
//  height: calc(100vh - 300px);
//}

.task-log-records__record {
  display: flex;
  flex-direction: row;
  justify-content: left;
}

.task-log-records__time {
  width: 120px;
  min-width: 120px;
}

.task-log-records__output {
  width: 100%;
}

//@media #{map-get($display-breakpoints, 'sm-and-down')} {
//  .task-log-records {
//    height: calc(100vh - 340px);
//  }
//
//  .task-log-view--with-message .task-log-records {
//    height: calc(100vh - 370px);
//  }
//}
</style>
<script>
import { getTaskById, getTaskOutputs, stopTask } from '@/api/task'
import { getUserInfo } from '@/api/user'
import TaskStatus from '@/components/task/TaskStatus.vue'
import socket from '@/socket'

export default {
  components: { TaskStatus },
  props: {
    itemId: {
      type: Number,
      default: 0
    },
  },
  data() {
    return {
      item: {},
      output: [],
      user: {},
    }
  },
  // watch: {
  //   async itemId() {
  //     this.reset()
  //     await this.loadData()
  //   },
  //   immediate: true,
  // },
  watch: {
    itemId: {
      immediate: true,
      handler: function() {
        this.reset()
        this.loadData()
      }
    }
  },
  async created() {
    socket.addListener((data) => this.onWebsocketDataReceived(data))
    await this.loadData()
  },
  methods: {
    async stopTask(Id) {
      await stopTask(Id)
    },

    reset() {
      this.item = {}
      this.output = []
      this.user = {}
    },

    onWebsocketDataReceived(data) {
      if (data.project_id !== this.projectId || data.task_id !== this.itemId) {
        return
      }

      switch (data.type) {
        case 'update':
          Object.assign(this.item, {
            ...data,
            type: undefined,
          })
          break
        case 'log':
          this.output.push(data.taskOutputs)
          setTimeout(() => {
            this.$refs.output.scrollTop = this.$refs.output.scrollHeight
          }, 200)
          break
        default:
          break
      }
    },

    async loadData() {
      console.log(this.itemId)
      this.item = (await getTaskById({ ID: this.itemId })).data

      this.output = (await getTaskOutputs({ taskId: this.itemId })).data.taskOutputs
      console.log(this.output)

      this.user = (await getUserInfo).data
    },
  },
}
</script>
