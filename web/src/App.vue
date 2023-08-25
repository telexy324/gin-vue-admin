<template>
  <div id="app">
<!--    <el-dialog v-model="taskLogDialog" :before-close="onTaskLogDialogClosed" :title="dialogTitle">-->
<!--      <template #header="{ titleId, titleClass }">-->
<!--        <div class="el-dialog__header">-->
<!--          <h4 :id="titleId" :class="titleClass">Task #{{ task ? task.ID : null }}</h4>-->
<!--          <el-button type="danger" @click="onTaskLogDialogClosed()">-->
<!--            <el-icon class="el-icon&#45;&#45;left"></el-icon>-->
<!--            Close-->
<!--          </el-button>-->
<!--        </div>-->
<!--      </template>-->
    <TaskLogView
      v-if="isGetData"
      :item-id="task ? task.ID : null"
      :visiable="!!task"
      @close="onTaskLogDialogClosed"
    />
<!--    <ScriptView-->
<!--      v-if="scriptDialog"-->
<!--      :script="script"-->
<!--      @close="onScriptLogDialogClosed"-->
<!--    />-->
<!--    </el-dialog>-->
    <router-view />
  </div>
</template>

<script>

const discoverServers = 99999999
// const gatherInformation = 99999998

import { getTaskById } from '@/api/task'
import { getTemplateById } from '@/api/template'
import { emitter } from '@/utils/bus'
// import { CircleCloseFilled } from '@element-plus/icons-vue'
import TaskLogView from '@/components/task/TaskLogView.vue'
// import ScriptView from '@/components/task/ScriptView.vue'
import socket from '@/socket'

export default {
  name: 'App',
  components: {
    // CircleCloseFilled,
    TaskLogView,
    // ScriptView
  },
  data() {
    return {
      taskLogDialog: null,
      task: null,
      template: null,
      taskID: 0,
      // dialogTitle: '',
      isGetData: false,
      script: '',
      scriptDialog: null,
    }
  },
  mounted() {
    if (!socket.isRunning()) {
      socket.start()
    }
    emitter.on('i-show-task', async(e) => {
      // if (parseInt(this.$route.query.t || '', 10) !== e.ID) {
      //   const query = { ...this.$route.query, t: e.ID }
      //   await this.$router.replace({ query })
      // }
      if (!socket.isRunning()) {
        socket.start()
      }
      this.task = (await getTaskById({ ID: e.ID })).data.task
      if (!(e.templateId >= discoverServers - 100)) {
        this.template = (await getTemplateById({ ID: e.templateId })).data.template
      }
      // this.dialogTitle = 'Task #' + this.task.ID
      this.isGetData = true
      // this.taskLogDialog = true
    })
    // emitter.on('i-show-script', (e) => {
    //   if (!socket.isRunning()) {
    //     socket.start()
    //   }
    //   this.script = e
    //   console.log(e)
    //   this.scriptDialog = true
    // })
  },
  methods: {
    async onTaskLogDialogClosed() {
      // this.taskLogDialog = false
      // const query = { ...this.$route.query, t: undefined }
      // await this.$router.replace({ query })
      // await this.$router.push({ name: this.$route.name })
      this.isGetData = false
      emitter.emit('i-close-task')
    },
    async onScriptLogDialogClosed() {
      this.scriptDialog = false
      // const query = { ...this.$route.query, t: undefined }
      // await this.$router.replace({ query })
      await this.$router.push({ name: this.$route.name })
    },
  }
}

</script>

<style lang="scss">
// 引入初始化样式
@import '@/style/main.scss';
@import '@/style/base.scss';
@import '@/style/mobile.scss';
@import '@/style/antv.scss';

#app {
  background: #eee;
  height: 100vh;
  overflow: hidden;
  font-weight: 400 !important;
}

.el-button {
  font-weight: 400 !important;
}
</style>
