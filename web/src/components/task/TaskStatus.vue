<template>
  <el-tag
    v-if="status"
    :key="status"
    :color="getStatusColor(status)"
    class="mx-1"
    effect="dark"
    round
  >
    {{ humanizeStatus(status) }}
  </el-tag>
</template>
<script>

const TaskStatus = Object.freeze({
  WAITING: 'waiting',
  RUNNING: 'running',
  SUCCESS: 'success',
  ERROR: 'error',
  STOPPING: 'stopping',
  STOPPED: 'stopped',
})

export default {
  props: {
    status: String,
  },

  methods: {
    humanizeStatus(status) {
      switch (status) {
        case TaskStatus.WAITING:
          return 'Waiting'
        case TaskStatus.RUNNING:
          return 'Running'
        case TaskStatus.SUCCESS:
          return 'Success'
        case TaskStatus.ERROR:
          return 'Failed'
        case TaskStatus.STOPPING:
          return 'Stopping...'
        case TaskStatus.STOPPED:
          return 'Stopped'
        default:
          throw new Error(`Unknown task status ${status}`)
      }
    },

    getStatusColor(status) {
      switch (status) {
        case TaskStatus.WAITING:
          return ''
        case TaskStatus.RUNNING:
          return 'primary'
        case TaskStatus.SUCCESS:
          return 'success'
        case TaskStatus.ERROR:
          return 'error'
        case TaskStatus.STOPPING:
          return ''
        case TaskStatus.STOPPED:
          return ''
        default:
          throw new Error(`Unknown task status ${status}`)
      }
    },
  },
}
</script>
