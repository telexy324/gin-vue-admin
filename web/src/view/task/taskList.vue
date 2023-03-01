<template>
  <div class="taskList">
    <el-collapse v-model="activeNames" @change="handleChange">
      <el-collapse-item title="Consistency" name="1">
        <div>
          Consistent with real life: in line with the process and logic of real
          life, and comply with languages and habits that the users are used to;
        </div>
        <div>
          Consistent within interface: all elements should be consistent, such
          as: design style, icons and texts, position of elements, etc.
        </div>
      </el-collapse-item>
      <el-collapse-item title="Feedback" name="2">
        <div>
          Operation feedback: enable the users to clearly perceive their
          operations by style updates and interactive effects;
        </div>
        <div>
          Visual feedback: reflect current state by updating or rearranging
          elements of the page.
        </div>
      </el-collapse-item>
      <el-collapse-item title="Efficiency" name="3">
        <div>
          Simplify the process: keep operating process simple and intuitive;
        </div>
        <div>
          Definite and clear: enunciate your intentions clearly so that the
          users can quickly understand and make decisions;
        </div>
        <div>
          Easy to identify: the interface should be straightforward, which helps
          the users to identify and frees them from memorizing and recalling.
        </div>
      </el-collapse-item>
      <el-collapse-item title="Controllability" name="4">
        <div>
          Decision making: giving advices about operations is acceptable, but do
          not make decisions for the users;
        </div>
        <div>
          Controlled consequences: users should be granted the freedom to
          operate, including canceling, aborting or terminating current
          operation.
        </div>
      </el-collapse-item>
    </el-collapse>
  </div>
</template>

<script>
import { getSystemConfig, setSystemConfig } from '@/api/system'
import { emailTest } from '@/api/email'
export default {
  name: 'Config',
  data() {
    return {
      config: {
        system: {},
        jwt: {},
        casbin: {},
        mysql: {},
        excel: {},
        autoCode: {},
        redis: {},
        qiniu: {},
        tencentCOS: {},
        aliyunOSS: {},
        captcha: {},
        zap: {},
        local: {},
        email: {},
        timer: {
          detail: {}
        }
      }
    }
  },
  async created() {
    await this.initForm()
  },
  methods: {
    async initForm() {
      const res = await getSystemConfig()
      if (res.code === 0) {
        this.config = res.data.config
      }
    },
    reload() {},
    async update() {
      const res = await setSystemConfig({ config: this.config })
      if (res.code === 0) {
        this.$message({
          type: 'success',
          message: '配置文件设置成功'
        })
        await this.initForm()
      }
    },
    async email() {
      const res = await emailTest()
      if (res.code === 0) {
        this.$message({
          type: 'success',
          message: '邮件发送成功'
        })
        await this.initForm()
      } else {
        this.$message({
          type: 'error',
          message: '邮件发送失败'
        })
      }
    }
  }
}
</script>

<style lang="scss">
.taskList {
  background: #fff;
  padding:36px;
  border-radius: 2px;
  h2 {
    padding: 10px;
    margin: 10px 0;
    font-size: 16px;
    box-shadow: -4px 0px 0px 0px #e7e8e8;
  }
  ::v-deep(.el-input-number__increase){
    top:5px !important;
  }
  .gva-btn-list{
    margin-top:16px;
  }
}
</style>
