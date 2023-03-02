<template>
  <div id="app">
    <router-view />
  </div>
</template>

<script>
export default {
  name: 'App'
}

EventBus.$on('i-show-task', async (e) => {
  if (parseInt(this.$route.query.t || '', 10) !== e.taskId) {
    const query = { ...this.$route.query, t: e.taskId }
    await this.$router.replace({ query })
  }

  this.task = (await axios({
    method: 'get',
    url: `/api/project/${this.projectId}/tasks/${e.taskId}`,
    responseType: 'json',
  })).data

  this.template = (await axios({
    method: 'get',
    url: `/api/project/${this.projectId}/templates/${this.task.template_id}`,
    responseType: 'json',
  })).data

  this.taskLogDialog = true
})

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
