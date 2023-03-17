<template>
  <el-dialog v-model="visible" :show-close="false" custom-class="customClass">
    <template #title>
      <el-row :gutter="10">
        <el-col :span="3">Task #{{ item.ID }}</el-col>
        <el-col :span="3">
          <TaskStatus :status="item.status" />
        </el-col>
        <el-col :span="3">Author:</el-col>
        <el-col :span="3" v-text="user.userName" />
        <el-col :span="3">start:</el-col>
        <el-col :span="3">{{ formatDate(item.beginTime.Time) }}</el-col>
        <el-col :span="3">end:</el-col>
        <el-col :span="3">{{ formatDate(item.endTime.Time) }}</el-col>
      </el-row>
    </template>
    <div ref="output" class="task-log-records">
      <div v-for="record in output" :key="record.ID" class="task-log-records__record">
        <div class="task-log-records__time">
          {{ formatDate(record.recordTime) }}
        </div>
        <div class="task-log-records__output">{{ record.output }}</div>
      </div>
    </div>
    <el-button
      v-if="item.status === 'running' || item.status === 'waiting'"
      type="danger"
      round
      style="position: absolute; bottom: 10px; right: 10px;"
      @click="stopTask(item.ID)"
    >
      Stop
    </el-button>
  </el-dialog>
</template>

<script>
import { getTaskById, getTaskOutputs, stopTask } from '@/api/task'
import { getUserById } from '@/api/user'
import TaskStatus from '@/components/task/TaskStatus.vue'
import socket from '@/socket'
import { formatTimeToStr } from '@/utils/date'

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
      visible: false,
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
    async itemId() {
      this.reset()
      await this.loadData()
    },
  },
  async created() {
    if (!socket.isRunning()) {
      socket.start()
    }
    socket.addListener((data) => this.onWebsocketDataReceived(data))
    console.log(document.baseURI)
    await this.loadData()
  },
  methods: {
    async stopTask(Id) {
      await stopTask({ ID: Id })
    },

    reset() {
      this.item = {}
      this.output = []
      this.user = {}
    },

    onWebsocketDataReceived(data) {
      if (data.taskId !== this.itemId) {
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
          this.output.push(data)
          setTimeout(() => {
            this.$refs.output.scrollTop = this.$refs.output.scrollHeight
          }, 200)
          break
        default:
          break
      }
    },

    async loadData() {
      this.item = (await getTaskById({ ID: this.itemId })).data.task

      this.output = (await getTaskOutputs({ taskId: this.itemId })).data.taskOutputs
      this.user = (await getUserById({ ID: this.item.systemUserId })).data.user
      console.log(this.user)
      this.visible = true
    },

    formatDate: function(time) {
      if (time !== null && time !== '') {
        var date = new Date(time)
        return formatTimeToStr(date, 'yyyy-MM-dd hh:mm:ss')
      } else {
        return ''
      }
    },
  },
}
</script>

<style lang="scss">
.customClass {
  width: 60%;
}
</style>

<style scoped lang="scss">

// @import '~vuetify/src/styles/settings/_variables';

.task-log-records {
  background: black;
  color: white;
  height: calc(100vh - 250px);
  overflow: auto;
  font-family: monospace;
  margin: 0 -18px;
  padding: 5px 10px;
}

.task-log-view--with-message .task-log-records {
  height: calc(100vh - 300px);
}

.task-log-records__record {
  display: flex;
  flex-direction: row;
  justify-content: left;
}

.task-log-records__time {
  width: 140px;
  min-width: 140px;
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
