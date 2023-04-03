<template>
  <div class="all">
    <div class="antv-content">
      <div class="antv-menu">
        <h3> 基础图形列表 </h3>
        <ul class="menu-list">
          <li draggable="true" @drag="menuDrag('defaultOval')"><i class="icon-oval" /> <strong>数据库</strong></li>
          <li draggable="true" @drag="menuDrag('defaultSquare')"><i class="icon-square" /><strong>中间件</strong></li>
          <li draggable="true" @drag="menuDrag('defaultYSquare')"><i class="icon-ySquare" /><strong>消息队列</strong></li>
          <li draggable="true" @drag="menuDrag('defaultRhombus')"><i class="icon-rhombus" /><strong>代理</strong></li>
<!--          <li draggable="true" @drag="menuDrag('defaultRhomboid')"><i class="icon-rhomboid" /><strong>平行四边形</strong></li>-->
          <li draggable="true" @drag="menuDrag('defaultCircle')"><i class="icon-circle" /><strong>缓存</strong></li>
          <li draggable="true" @drag="menuDrag('otherImage')"><i class="el-icon-picture" /><strong>自定义图片</strong></li>
        </ul>
        <div v-if="isChange" class="wrapper-btn">
          <el-button type="success" @click="handlerSend">保存当前方案</el-button>
        </div>
      </div>
      <div class="antv-wrapper">
        <div id="wrapper" class="wrapper-canvas" :style="{height: height}" @drop="drop($event)" @dragover.prevent>
          <div class="wrapper-tips">
            <div class="wrapper-tips-item">
              <el-switch v-model="isPortsShow" @change="changePortsShow" />
              <span>链接桩常显</span>
            </div>
          </div>
        </div>
      </div>
      <div v-if="editDrawer" class="edit-main">
        <div class="edit-main-title">
          <h3>{{ editTitle }} </h3>
          <i class="el-icon-close" @click="closeEditForm"></i>
        </div>
        <div v-if="editTitle === '编辑节点'" class="form-main">
          <el-form ref="nodeForm" :model="form" label-width="80px">
            <el-form-item label="节点文本">
              <el-input v-model="form.labelText" size="small" @input="changeNode('labelText', form.labelText)" />
            </el-form-item>
            <el-form-item label="字体大小">
              <el-input v-model="form.fontSize" size="small" @input="changeNode('fontSize', form.fontSize)" />
            </el-form-item>
            <el-form-item label="字体颜色">
              <el-color-picker v-model="form.fontFill" @change="changeNode('fontFill', form.fontFill)" />
            </el-form-item>
            <el-form-item label="节点背景">
              <el-color-picker v-model="form.fill" @change="changeNode('fill', form.fill)" />
            </el-form-item>
            <el-form-item label="边框颜色">
              <el-color-picker v-model="form.stroke" @change="changeNode('stroke', form.stroke)" />
            </el-form-item>
            <div class="see-box">
              <h5>预览</h5>
              <div class="see-item" :style="{ 'background': form.fill, 'color': form.fontFill, 'border-color': form.stroke, 'font-size': form.fontSize + 'px' }">{{form.labelText}}</div>
            </div>
          </el-form>
        </div>
        <div v-if="editTitle === '编辑图片节点'" class="form-main">
          <el-form ref="imageForm" :model="form" label-width="80px">
            <el-form-item label="节点文本">
              <el-input v-model="form.labelText" size="small" @input="changeImageNode('labelText', form.labelText)" />
            </el-form-item>
            <el-form-item label="字体颜色">
              <el-color-picker v-model="form.labelFill" @change="changeImageNode('labelFill', form.labelFill)" />
            </el-form-item>
            <el-form-item label="节点背景">
              <el-color-picker v-model="form.fill" @change="changeImageNode('fill', form.fill)" />
            </el-form-item>
            <el-form-item label="图片地址">
              <el-input v-model="form.xlinkHref" size="small" placeholder="图片地址" @input="changeImageNode('xlinkHref', form.xlinkHref)" />
              <el-image :src="form.xlinkHref" style="width: 80px; height: 80px; background: #f2f2f2" fit="fill" />
            </el-form-item>
            <el-form-item label="图片尺寸">
              <span style="font-size: 14px; padding-right: 5px; color: #888;">宽</span><el-input-number v-model="form.width" :min="0" label="宽" size="mini" @change="changeImageNode('width', form.width)" />
              <span style="font-size: 14px; padding-right: 5px; color: #888;">高</span><el-input-number v-model="form.height" :min="0" label="高" size="mini" @change="changeImageNode('height', form.height)" />
            </el-form-item>
          </el-form>
        </div>
        <div v-if="editTitle === '编辑连线'" class="form-main">
          <el-form ref="edgeForm" :model="form" label-width="80px">
            <el-form-item label="标签内容">
              <el-input
                v-model="form.label"
                size="small"
                placeholder="标签文字，空则没有"
                @input="changeEdgeLabel(form.label, labelForm.fontColor, labelForm.fill, labelForm.stroke)"
              />
              <div v-if="form.label" class="label-style">
                <p>字体颜色：<el-color-picker v-model="labelForm.fontColor" size="mini" @change="changeEdgeLabel(form.label, labelForm.fontColor, labelForm.fill, labelForm.stroke)" /></p>
                <p>背景颜色：<el-color-picker v-model="labelForm.fill" size="mini" @change="changeEdgeLabel(form.label, labelForm.fontColor, labelForm.fill, labelForm.stroke)" /></p>
                <p>描边颜色：<el-color-picker v-model="labelForm.stroke" size="mini" @change="changeEdgeLabel(form.label, labelForm.fontColor, labelForm.fill, labelForm.stroke)" /></p>
              </div>
            </el-form-item>
            <el-form-item label="线条颜色">
              <el-color-picker v-model="form.stroke" size="small" @change="changeEdgeStroke" />
            </el-form-item>
            <el-form-item label="线条样式">
              <el-select v-model="form.connector" size="small" placeholder="请选择" @change="changeEdgeConnector">
                <el-option label="直角" value="normal" />
                <el-option label="圆角" value="rounded" />
                <el-option label="平滑" value="smooth" />
                <el-option label="跳线(两线交叉)" value="jumpover" />
              </el-select>
            </el-form-item>
            <el-form-item label="线条宽度">
              <el-input-number
                v-model="form.strokeWidth"
                size="small"
                :min="2"
                :step="2"
                :max="6"
                label="线条宽度"
                @change="changeEdgeStrokeWidth"
              />
            </el-form-item>
            <el-form-item label="双向箭头">
              <el-switch v-model="form.isArrows" @change="changeEdgeArrows" />
            </el-form-item>
            <el-form-item label="流动线条">
              <el-switch v-model="form.isAnit" @change="changeEdgeAnit" />
            </el-form-item>
            <el-form-item label="调整线条">
              <el-switch v-model="form.isTools" @change="changeEdgeTools" />
            </el-form-item>
          </el-form>
        </div>
        <div class="edit-btn">
          <el-button
            type="danger"
            style="width:100%"
            @click="handlerDel"
          >删除此{{ editTitle === '编辑节点' ? '节点' : '连线' }}</el-button>
        </div>
      </div>
    </div>
  </div>
</template>
<script>
import { Graph, Shape } from '@antv/x6'
import { configSetting, configNodeShape, configNodePorts, configEdgeLabel, graphBindKey } from '@/utils/antvSetting'
import {
  updateEditRelation,
  getSystemEditRelation,
} from '@/api/cmdb'
export default {
  name: 'AntV6X',
  /**
   * 这个是作为子组件分别接受了两个数据一个是高度height，一个是反显图表数据tempGroupJson
   * 作为子组件例子 <AntVXSix v-model="tempGroupJson" height="720px" />
   *
   */
  props: {
    height: {
      type: String,
      default: '100vh' // '720px'
    },
    value: {
      type: String,
      default: ''
    }
  },
  data() {
    return {
      graph: null,
      isChange: false,
      isPortsShow: false,
      menuItem: '',
      selectCell: '',
      editDrawer: false,
      editTitle: '',
      form: {},
      labelForm: {
        fontColor: '#333',
        fill: '#FFF',
        stroke: '#555'
      },
      systemId: 0,
      defaultId: 0,
      defaultValue: null,
    }
  },
  watch: {
    value: {
      handler: function() {
        if (this.graph) {
          this.isChange = false
          this.isPortsShow = false
          this.menuItem = ''
          this.selectCell = ''
          this.editDrawer = false
          this.graph.dispose()
          this.initGraph('watch')
        }
      },
      deep: true,
      immediate: true
    }
  },
  created() {
  },
  async mounted() {
    this.systemId = this.$route.params.systemId
    const res = (await getSystemEditRelation({
      'ID': Number(this.systemId)
    })).data
    this.defaultId = res.ID
    this.defaultValue = res.relation
    this.initGraph('default')
  },
  beforeDestroy() {
    this.graph.dispose()
  },
  methods: {
    // 链接桩的显示与隐藏，主要是照顾菱形
    changePortsShow(val) {
      const container = document.getElementById('wrapper')
      const ports = container.querySelectorAll('.x6-port-body')
      for (let i = 0, len = ports.length; i < len; i = i + 1) {
        ports[i].style.visibility = val ? 'visible' : 'hidden'
      }
    },
    // 初始化渲染画布
    initGraph(f) {
      const graph = new Graph({
        container: document.getElementById('wrapper'),
        ...configSetting(Shape),
      })
      // 画布事件
      graph.on('node:mouseenter', () => {
        this.changePortsShow(true)
      })
      graph.on('node:mouseleave', () => {
        if (this.isPortsShow) return
        this.changePortsShow(false)
      })
      // 点击编辑
      graph.on('cell:click', ({ cell }) => {
        this.editForm(cell)
      })
      // 画布键盘事件
      graphBindKey(graph)
      // 删除
      graph.bindKey(['delete', 'backspace'], () => {
        this.handlerDel()
      })
      // 赋值
      this.graph = graph
      // 返现方法
      if (f === 'watch' && this.value && JSON.parse(this.value).length) {
        const resArr = JSON.parse(this.value)
        // 导出的时候删除了链接桩设置加回来
        const portsGroups = configNodePorts().groups
        if (resArr.length) {
          const jsonTemp = resArr.map(item => {
            if (item.ports) item.ports.groups = portsGroups
            return item
          })
          graph.fromJSON(jsonTemp)
        }
      } else if (f === 'default' && this.defaultValue && JSON.parse(this.defaultValue).length) {
        const resArr = JSON.parse(this.defaultValue)
        // 导出的时候删除了链接桩设置加回来
        const portsGroups = configNodePorts().groups
        if (resArr.length) {
          const jsonTemp = resArr.map(item => {
            if (item.ports) item.ports.groups = portsGroups
            return item
          })
          graph.fromJSON(jsonTemp)
        }
      }
      // 画布有变化
      graph.on('cell:changed', () => {
        this.isChangeValue()
      })
    },
    // 画布是否有变动
    isChangeValue() {
      if (!this.isChange) {
        this.isChange = true
        this.$emit('cellChanged', true)
      }
    },
    menuDrag(type) {
      this.menuItem = configNodeShape(type)
      switch (type) {
        case 'defaultOval':
          this.menuItem.customAttr = '数据库'
          break
        case 'defaultSquare':
          this.menuItem.customAttr = '中间件'
          break
        case 'defaultYSquare':
          this.menuItem.customAttr = '消息队列'
          break
        case 'defaultRhombus':
          this.menuItem.customAttr = '代理'
          break
        case 'defaultCircle':
          this.menuItem.customAttr = '缓存'
          break
        case 'otherImage':
          this.menuItem.customAttr = '自定义图片'
          break
        default:
          this.menuItem.customAttr = ''
      }
    },
    drop(event) {
      const nodeItem = {
        ...this.menuItem,
        x: event.offsetX - (this.menuItem.width / 2),
        y: event.offsetY - (this.menuItem.height / 2),
        ports: configNodePorts()
      }
      // 创建节点
      const node = this.graph.addNode(nodeItem)
      node.attr('label/text', nodeItem.customAttr)
      this.isChangeValue()
    },
    editForm(cell) {
      if (this.selectCell) this.selectCell.removeTools() // 删除修改线的工具
      this.selectCell = cell
      // 编辑node节点
      if (cell.isNode() && cell.data.type && cell.data.type.includes('default')) {
        this.editTitle = '编辑节点'
        const body = cell.attrs.body || cell.attrs.rect || cell.attrs.polygon || cell.attrs.circle
        this.form = {
          labelText: cell.attrs.label.text || '',
          fontSize: cell.attrs.label.fontSize || 14,
          fontFill: cell.attrs.label.fill || '',
          fill: body.fill || '',
          stroke: body.stroke || ''
        }
        this.editDrawer = true
        return
      }
      // 编辑图片节点
      if (cell.isNode() && cell.data.type && cell.data.type === 'otherImage') {
        this.editTitle = '编辑图片节点'
        const attrs = cell.attrs || { body: { fill: '' }, label: { text: '', fill: '' }, image: { xlinkHref: '', height: 80, width: 80 }}
        this.form = {
          fill: attrs.body.fill,
          labelText: attrs.label.text,
          labelFill: attrs.label.fill,
          height: (attrs.image && attrs.image.height) || 80,
          width: (attrs.image && attrs.image.width) || 80,
          xlinkHref: (attrs.xlinkHref && attrs.image.xlinkHref) || 'https://gw.alipayobjects.com/zos/bmw-prod/2010ac9f-40e7-49d4-8c4a-4fcf2f83033b.svg'
        }
        this.editDrawer = true
        return
      }
      // 编辑线
      if (!cell.isNode() && cell.shape === 'edge') {
        this.editTitle = '编辑连线'
        this.form = {
          label: (cell.labels && cell.labels[0]) ? cell.labels[0].attrs.labelText.text : '',
          stroke: cell.attrs.line.stroke || '',
          connector: 'rounded',
          strokeWidth: cell.attrs.line.strokeWidth || '',
          isArrows: !!cell.attrs.line.sourceMarker,
          isAnit: !!cell.attrs.line.strokeDasharray,
          isTools: false
        }
        // 看是否有label
        const edgeCellLabel = cell.labels && cell.labels[0] && cell.labels[0].attrs || false
        if (this.form.label && edgeCellLabel) {
          this.labelForm = {
            fontColor: edgeCellLabel.labelText.fill || '#333',
            fill: edgeCellLabel.labelBody.fill || '#fff',
            stroke: edgeCellLabel.labelBody.stroke || '#555'
          }
        } else {
          this.labelForm = { fontColor: '#333', fill: '#FFF', stroke: '#555' }
        }
        this.editDrawer = true
      }
    },
    closeEditForm() {
      this.editDrawer = false
      if (this.selectCell) this.selectCell.removeTools()
    },
    // 修改一般节点
    changeNode(type, value) {
      switch (type) {
        case 'labelText':
          this.selectCell.attr('label/text', value)
          break
        case 'fontSize':
          this.selectCell.attr('label/fontSize', value)
          break
        case 'fontFill':
          this.selectCell.attr('label/fill', value)
          break
        case 'fill':
          this.selectCell.attr('body/fill', value)
          break
        case 'stroke':
          this.selectCell.attr('body/stroke', value)
          break
      }
    },
    // 修改图片节点
    changeImageNode(type, value) {
      switch (type) {
        case 'labelText':
          this.selectCell.attr('label/text', value)
          break
        case 'labelFill':
          this.selectCell.attr('label/fill', value)
          break
        case 'fill':
          this.selectCell.attr('body/fill', value)
          break
        case 'xlinkHref':
          this.selectCell.attr('image/xlinkHref', value)
          break
        case 'height':
          this.selectCell.attr('image/height', value)
          break
        case 'width':
          this.selectCell.attr('image/width', value)
          break
      }
    },
    // 修改边label属性
    changeEdgeLabel(label, fontColor, fill, stroke) {
      this.selectCell.setLabels([configEdgeLabel(label, fontColor, fill, stroke)])
      if (!label) this.labelForm = { fontColor: '#333', fill: '#FFF', stroke: '#555' }
    },
    // 修改边的颜色
    changeEdgeStroke(val) {
      this.selectCell.attr('line/stroke', val)
    },
    // 边的样式
    changeEdgeConnector(val) {
      switch (val) {
        case 'normal':
          this.selectCell.setConnector(val)
          break
        case 'smooth':
          this.selectCell.setConnector(val)
          break
        case 'rounded':
          this.selectCell.setConnector(val, { radius: 20 })
          break
        case 'jumpover':
          this.selectCell.setConnector(val, { radius: 20 })
          break
      }
    },
    // 边的宽度
    changeEdgeStrokeWidth(val) {
      if (this.form.isArrows) {
        this.selectCell.attr({
          line: {
            strokeWidth: val,
            sourceMarker: {
              width: 12 * (val / 2) || 12,
              height: 8 * (val / 2) || 8
            },
            targetMarker: {
              width: 12 * (val / 2) || 12,
              height: 8 * (val / 2) || 8
            }
          }
        })
      } else {
        this.selectCell.attr({
          line: {
            strokeWidth: val,
            targetMarker: {
              width: 12 * (val / 2) || 12,
              height: 8 * (val / 2) || 8
            }
          }
        })
      }
    },
    // 边的箭头
    changeEdgeArrows(val) {
      if (val) {
        this.selectCell.attr({
          line: {
            sourceMarker: {
              name: 'block',
              width: 12 * (this.form.strokeWidth / 2) || 12,
              height: 8 * (this.form.strokeWidth / 2) || 8
            },
            targetMarker: {
              name: 'block',
              width: 12 * (this.form.strokeWidth / 2) || 12,
              height: 8 * (this.form.strokeWidth / 2) || 8
            },
          }
        })
      } else {
        this.selectCell.attr({
          line: {
            sourceMarker: '',
            targetMarker: {
              name: 'block',
              size: 10 * (this.form.strokeWidth / 2) || 10
            },
          }
        })
      }
    },
    // 边的添加蚂蚁线
    changeEdgeAnit(val) {
      if (val) {
        this.selectCell.attr({
          line: {
            strokeDasharray: 5,
            style: {
              animation: 'ant-line 30s infinite linear',
            }
          }
        })
      } else {
        this.selectCell.attr({
          line: {
            strokeDasharray: 0,
            style: {
              animation: '',
            }
          }
        })
      }
    },
    // 给线添加调节工具
    changeEdgeTools(val) {
      if (val) this.selectCell.addTools(['vertices', 'segments'])
      else this.selectCell.removeTools()
    },
    // 删除节点
    handlerDel() {
      this.$confirm(`此操作将永久删除此${this.editTitle === '编辑节点' ? '节点' : '连线'}, 是否继续?`, '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        const cells = this.graph.getSelectedCells()
        if (cells.length) {
          this.graph.removeCells(cells)
          this.form = {}
          this.editDrawer = false
          this.$message({ type: 'success', message: '删除成功!' })
        }
      }).catch(() => {})
    },
    // 导出
    async handlerSend() {
      // 我在这里删除了链接桩的设置，和工具（为了减少数据），反显的时候要把删除的链接桩加回来
      const { cells: jsonArr } = this.graph.toJSON()
      const tempGroupJson = jsonArr.map(item => {
        if (item.ports && item.ports.groups) delete item.ports.groups
        if (item.tools) delete item.tools
        return item
      })
      if (this.selectCell) {
        this.selectCell.removeTools()
        this.selectCell = ''
      }
      this.$emit('finish', JSON.stringify(tempGroupJson))
      // console.log(JSON.stringify(tempGroupJson))
      const res = await updateEditRelation({
        ID: Number(this.defaultId),
        systemId: Number(this.systemId),
        relation: JSON.stringify(tempGroupJson)
      })
      if (res.code === 0) {
        this.$message({
          type: 'success',
          message: '保存成功',
          showClose: true
        })
      }
    },
  }
}
</script>
<style lang="scss">
@keyframes ant-line {
  to {
      stroke-dashoffset: -1000
  }
}
</style>
<style lang="scss" scoped="scoped">
.all{
  border-radius: 8px;
  overflow: hidden;
}
.antv-content{
  background: #fff;
  display: flex;
  overflow: hidden;
  position: relative;
  .antv-menu{
    width: 200px;
    border-right: 1px solid #d5d5d5;
    padding: 10px;
    h3{
      padding: 10px;
    };
    li{
      padding: 10px;
      border-radius: 8px;
      border: 1px solid #555;
      background: #fff;
      margin: 5px 10px;
      font-size: 12px;
      display: flex;
      align-items: center;
      cursor: pointer;
      transition: all 0.5s ease;
      &:hover{
        box-shadow: 0 0 5px rgba($color: #000000, $alpha: 0.3);
      }
      i{
        font-size: 18px;
        margin-right: 10px;
      }
      strong{
        flex: 1;
      }
    }
  }
  .antv-wrapper{
    flex: 1;
    position: relative;
    .wrapper-canvas{
      position: relative;
      height: 100vh;
      min-height: 720px;
    }
    .wrapper-tips{
      padding: 10px;
      display: flex;
      align-items: center;
      position: absolute;
      top: 0;
      left: 0;
      .wrapper-tips-item{
        span{
          padding-left: 10px;
          font-size: 12px;
        }
      }
    }
  }
}
i.icon-oval{
    display: inline-block;
    width: 16px;
    height: 10px;
    border-radius: 10px;
    border: 2px solid #555;
}
i.icon-square{
    display: inline-block;
    width: 16px;
    height: 10px;
    border: 2px solid #555;
}
i.icon-ySquare{
   display: inline-block;
    width: 16px;
    height: 10px;
    border-radius: 4px;
    border: 2px solid #555;
}
i.icon-rhombus{
   display: inline-block;
    width: 10px;
    height: 10px;
    border: 2px solid #555;
    transform: rotate(45deg);
}
i.icon-rhomboid{
   display: inline-block;
    width: 10px;
    height: 10px;
    border: 2px solid #555;
    transform: skew(-30deg);
}
i.icon-circle{
   display: inline-block;
    width: 16px;
    height: 16px;
    border-radius: 16px;
    border: 2px solid #555;
}
.edit-main{
  position: absolute;
  right: 0;
  top: 0;
  height: 100%;
  width: 280px;
  border-left: 1px solid #f2f2f2;
  box-shadow: 0 -10px 10px rgba($color: #000000, $alpha: 0.3);
  padding: 20px;
  background: #fff;
  box-sizing: border-box;
  .edit-main-title{
    display: flex;
    justify-content: space-between;
    align-items: center;
    h3{
      flex: 1;
    }
    i{
      cursor: pointer;
      font-size: 20px;
      opacity: 0.7;
      &:hover{
        opacity: 1;
      }
    }
  }
  .form-main{
    padding: 20px 0;
    .label-style{
      background: #f2f2f2;
      padding: 0 10px;
      p{
        display: flex;
        align-items: center;
        font-size: 12px;
      }
    }
  }
  .edit-btn{
  }
  .see-box{
  padding: 20px;
    background: #f2f2f2;
    h5{
      padding-bottom: 10px;
    }
    .see-item{
      padding: 10px 30px;
      border: 2px solid #333;
      text-align: center;
    }
  }
}
.wrapper-btn{
  text-align: center;
  padding: 20px;
  button{
    width: 100%;
  }
}
</style>
