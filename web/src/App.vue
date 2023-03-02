<template>
  <div id="app">
    <EditDialog
        v-model="taskLogDialog"
        save-button-text="Delete"
        :max-width="1000"
        :hide-buttons="true"
        @close="onTaskLogDialogClosed()"
    >
      <template v-slot:title={}>
        <div class="text-truncate" style="max-width: calc(100% - 36px);">
          <span class="breadcrumbs__item">Task #{{ task ? task.id : null }}</span>
        </div>

        <v-spacer></v-spacer>
        <v-btn
            icon
            @click="taskLogDialog = false; onTaskLogDialogClosed()"
        >
          <v-icon>mdi-close</v-icon>
        </v-btn>
      </template>
      <template v-slot:form="{}">
        <TaskLogView :project-id="projectId" :item-id="task ? task.id : null"/>
      </template>
    </EditDialog>
    <router-view />
  </div>
</template>

<script>
import { getTaskOutputs, getTask, getTaskTemplate } from '@/api/task'
import { emitter } from '@/utils/bus'

export default {
  name: 'App',
  data() {
    return {
      taskLogDialog: null,
      task: null,
      template: null,
    }
  },
  mounted() {
    emitter.$on('i-show-task', async(e) => {
      if (parseInt(this.$route.query.t || '', 10) !== e.taskId) {
        const query = { ...this.$route.query, t: e.taskId }
        await this.$router.replace({ query })
      }

      this.task = await getTask({ id: row.ID })

      this.template = (await axios({
        method: 'get',
        url: `/api/project/${this.projectId}/templates/${this.task.template_id}`,
        responseType: 'json',
      })).data

      this.taskLogDialog = true
    })
  },
  methods: {
    async onTaskLogDialogClosed() {
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
