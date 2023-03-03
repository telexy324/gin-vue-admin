<template>
  <div id="app">
    <el-dialog v-model="taskLogDialog" :before-close="onTaskLogDialogClosed">
      <template #header="{ titleId, titleClass }">
        <div class="my-header">
          <h4 :id="titleId" :class="titleClass">Task #{{ task ? task.id : null }}</h4>
          <el-button type="danger" @click="onTaskLogDialogClosed()">
            <el-icon class="el-icon--left"><CircleCloseFilled /></el-icon>
            Close
          </el-button>
        </div>
      </template>
      <template>
        <TaskLogView :project-id="projectId" :item-id="task ? task.id : null" />
      </template>
    </el-dialog>
    <router-view />
  </div>
</template>

<script>
import { getTask, getTaskTemplate } from '@/api/task'
import { emitter } from '@/utils/bus'
// import { CircleCloseFilled } from '@element-plus/icons-vue'
import TaskLogView from '@/components/task/TaskLogView.vue'

export default {
  name: 'App',
  components: {
    // CircleCloseFilled,
    TaskLogView
  },
  data() {
    return {
      taskLogDialog: null,
      task: null,
      template: null,
    }
  },
  mounted() {
    emitter.$on('i-show-task', async(e) => {
      if (parseInt(this.$route.query.t || '', 10) !== e.ID) {
        const query = { ...this.$route.query, t: e.ID }
        await this.$router.replace({ query })
      }

      this.task = await getTask({ id: e.ID }).data

      this.template = (await getTaskTemplate({ id: e.TemplateId })).data

      this.taskLogDialog = true
    })
  },
  methods: {
    async onTaskLogDialogClosed() {
      this.taskLogDialog = false
      const query = { ...this.$route.query, t: undefined }
      await this.$router.replace({ query })
    },
  }
}

</script>

<style lang="scss">
// 引入初始化样式
@import '@/style/main.scss';
@import '@/style/base.scss';
@import '@/style/mobile.scss';
#app {
  background: #eee;
  height: 100vh;
  overflow: hidden;
  font-weight: 400 !important;
}
.el-button{
  font-weight: 400 !important;
}
</style>
