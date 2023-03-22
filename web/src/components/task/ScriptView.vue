<template>
  <el-dialog
    v-model="visible"
    :show-close="true"
    title="script"
    custom-class="customClass"
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    :before-close="closeDialog"
  >
    <div ref="output" class="task-log-records">
      {{ script }}
    </div>
  </el-dialog>
</template>

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
<script>

export default {
  props: {
    script: {
      type: String,
      default: ''
    },
  },
  data() {
    return {
      visible: false,
    }
  },
  created() {
    this.visible = true
  },
  methods: {
    closeDialog() {
      this.visible = false
      this.$emit('close')
    },
  }
}
</script>
