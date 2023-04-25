package taskMdl

import (
	"encoding/json"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/application"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type TaskTemplate struct {
	global.GVA_MODEL
	Name            string                          `json:"name" gorm:"column:name"`                         // task名称
	Description     string                          `json:"description" gorm:"column:description"`           // task描述
	TargetServerIds string                          `json:"targetServerIds" gorm:"column:target_server_ids"` // 关联服务器id
	Mode            int                             `json:"mode" gorm:"column:mode"`                         // 执行方式 1 命令 2 脚本
	Command         string                          `json:"command" gorm:"column:command"`                   // 命令
	ScriptPath      string                          `json:"scriptPath" gorm:"column:script_path"`            // 脚本位置
	LastTaskId      int                             `json:"lastTaskId" gorm:"column:last_task_id"`           // 最后一次task id
	SysUser         string                          `json:"sysUser" gorm:"column:sys_user"`                  // 执行用户
	SystemId        int                             `json:"systemId" gorm:"column:system_id"`                // 所属系统
	ExecuteType     int                             `json:"executeType" gorm:"column:execute_type"`          // 模板类型 0 普通 1 日志提取
	LogPath         string                          `json:"logPath" gorm:"column:log_path"`                  // 日志位置
	TargetIds       []int                           `json:"targetIds" gorm:"-"`
	TargetServers   []application.ApplicationServer `json:"targetServers" gorm:"-"`
	LastTask        Task                            `json:"lastTask" gorm:"-"`
}

func (m *TaskTemplate) TableName() string {
	return "application_task_templates"
}

func (m *TaskTemplate) AfterFind(tx *gorm.DB) (err error) {
	serverIds := make([]int, 0)
	if m.TargetServerIds != "" {
		if err = json.Unmarshal([]byte(m.TargetServerIds), &serverIds); err != nil {
			global.GVA_LOG.Error("转换失败", zap.Any("err", err))
			return
		}
	}
	if len(serverIds) > 0 {
		for _, id := range serverIds {
			server := application.ApplicationServer{}
			if err = tx.Model(&application.ApplicationServer{}).Where("id = ?", id).Find(&server).Error; err != nil {
				global.GVA_LOG.Error("转换失败", zap.Any("err", err))
				return
			}
			m.TargetServers = append(m.TargetServers, server)
			m.TargetIds = append(m.TargetIds, int(server.ID))
		}
	}
	if m.LastTaskId > 0 {
		if err = tx.Model(&Task{}).Where("id = ?", m.LastTaskId).Find(&m.LastTask).Error; err != nil {
			global.GVA_LOG.Error("转换失败", zap.Any("err", err))
			return
		}
	}
	return nil
}

type TaskTemplateSet struct {
	global.GVA_MODEL
	SystemId int    `json:"systemId" gorm:"column:system_id" ` // 系统名称
	Name     string `json:"name" gorm:"column:name"`           // 模板集名称
}

func (m *TaskTemplateSet) TableName() string {
	return "application_task_template_sets"
}

type TaskTemplateSetTemplate struct {
	global.GVA_MODEL
	TemplateId int `json:"templateId" gorm:"column:template_id"` // task id
	SetId      int `json:"setId" gorm:"column:set_id"`           // task id
	Seq        int `json:"seq" gorm:"column:seq"`                // 排序
}

func (m *TaskTemplateSetTemplate) TableName() string {
	return "application_task_template_set_templates"
}

type TaskTemplateWithSeq struct {
	TaskTemplate
	Seq int `json:"seq"`
}

type SetTask struct {
	global.GVA_MODEL
	SetId           int                   `json:"setId" gorm:"column:set_id"`                 // task id
	SystemUserId    int                   `json:"systemUserId" gorm:"column:system_user_id" ` // 执行人
	CurrentTaskId   int                   `json:"currentTaskId" gorm:"column:current_task_id" `
	TotalSteps      int                   `json:"totalSteps" gorm:"column:total_steps" `
	CurrentStep     int                   `json:"currentStep" gorm:"column:current_step" `
	TemplatesString string                `json:"templatesString" gorm:"column:templates_string"` // 关联服务器id
	TasksString     string                `json:"tasksString" gorm:"column:tasks_string"`         // 关联服务器id
	Templates       []TaskTemplateWithSeq `json:"templates" gorm:"-"`
	Tasks           []Task                `json:"tasks" gorm:"-"`
}

func (m *SetTask) TableName() string {
	return "application_set_tasks"
}

func (m *SetTask) AfterFind(tx *gorm.DB) (err error) {
	templateIds := make([]int, 0)
	if m.TemplatesString != "" {
		if err = json.Unmarshal([]byte(m.TemplatesString), &templateIds); err != nil {
			global.GVA_LOG.Error("转换失败", zap.Any("err", err))
			return
		}
	}
	if len(templateIds) > 0 {
		for index, id := range templateIds {
			template := TaskTemplate{}
			if err = tx.Model(&TaskTemplate{}).Where("id = ?", id).Find(&template).Error; err != nil {
				global.GVA_LOG.Error("转换失败", zap.Any("err", err))
				return
			}
			templateWithSeq := TaskTemplateWithSeq{
				TaskTemplate: template,
				Seq:          index,
			}
			m.Templates = append(m.Templates, templateWithSeq)
		}
	}
	taskIds := make([]int, 0)
	if m.TasksString != "" {
		if err = json.Unmarshal([]byte(m.TasksString), &taskIds); err != nil {
			global.GVA_LOG.Error("转换失败", zap.Any("err", err))
			return
		}
	}
	if len(taskIds) > 0 {
		for _, id := range taskIds {
			task := Task{}
			if err = tx.Model(&Task{}).Where("id = ?", id).Find(&task).Error; err != nil {
				global.GVA_LOG.Error("转换失败", zap.Any("err", err))
				return
			}
			m.Tasks = append(m.Tasks, task)
		}
	}
	return nil
}
