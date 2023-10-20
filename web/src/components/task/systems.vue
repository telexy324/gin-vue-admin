<template>
  <div>
    <div class="clearflex">
      <el-button class="fl-right" size="mini" type="primary" @click="checked">确 定</el-button>
    </div>
    <div style="overflow: auto; height: calc(100% - 30px);">
      <el-tree
        ref="systemTree"
        style="display: inline-block;"
        :data="systemTreeData"
        :props="systemDefaultProps"
        default-expand-all
        :default-checked-keys="keys"
        highlight-current
        node-key="ID"
        show-checkbox
        @check="nodeChange"
      >
        <template #default="{ node }">
          <span class="custom-tree-node">
            <span>{{ node.label }}</span>
          </span>
        </template>
      </el-tree>
    </div>
  </div>
</template>

<script>
import { getAdminSystems } from '@/api/cmdb'
export default {
  name: 'Systems',
  props: {
    row: {
      default: function() {
        return {}
      },
      type: Object
    },
    keys: {
      type: Array,
      default: function() {
        return []
      },
    }
  },
  data() {
    return {
      systemTreeData: [],
      systemTreeIds: [],
      needConfirm: false,
      systemDefaultProps: {
        children: 'children',
        label: function(data) {
          return data.name
        }
      }
    }
  },
  async created() {
    // 获取所有菜单树
    const res = await getAdminSystems()
    this.systemTreeData = res.data.systems

    const arr = []
    res.data.systems.forEach(item => {
      arr.push(Number(item.ID))
    })
    this.systemTreeIds = arr
  },
  methods: {
    nodeChange() {
      this.needConfirm = true
    },
    async checked() {
      const checkArr = this.$refs.systemTree.getCheckedNodes(false, true)
      this.$emit('checked', checkArr)
    }
  }
}
</script>
