package taskMdl

import (
	"encoding/json"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/application"
	"github.com/flipped-aurora/gin-vue-admin/server/model/logUploadMdl"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type TaskTemplate struct {
	global.GVA_MODEL
	Name              string                          `json:"name" gorm:"type:varchar(100);not null;default:'';column:name"`                      // task名称
	Description       string                          `json:"description" gorm:"type:text;column:description"`                                    // task描述
	TargetServerIds   string                          `json:"targetServerIds" gorm:"type:text;column:target_server_ids"`                          // 关联服务器id
	Mode              int                             `json:"mode" gorm:"type:tinyint(2);not null;default:0;column:mode"`                         // 执行方式 1 命令 2 脚本
	Command           string                          `json:"command" gorm:"type:text;column:command"`                                            // 命令
	ScriptPath        string                          `json:"scriptPath" gorm:"type:varchar(255);not null;default:'';column:script_path"`         // 脚本位置
	LastTaskId        int                             `json:"lastTaskId" gorm:"type:bigint;not null;default:0;column:last_task_id"`               // 最后一次task id
	SysUser           string                          `json:"sysUser" gorm:"type:varchar(30);not null;default:'';column:sys_user"`                // 执行用户
	SystemId          int                             `json:"systemId" gorm:"type:bigint;not null;default:0;column:system_id"`                    // 所属系统
	ExecuteType       int                             `json:"executeType" gorm:"type:tinyint(2);not null;default:0;column:execute_type"`          // 模板类型 1 普通 2 日志提取 3 程序包上传
	LogPath           string                          `json:"logPath" gorm:"type:varchar(255);not null;default:'';column:log_path"`               // 日志位置
	ScriptHash        string                          `json:"scriptHash" gorm:"type:varchar(32);not null;default:'';column:script_hash"`          // 脚本哈希
	LogOutput         int                             `json:"logOutput" gorm:"type:tinyint(2);not null;default:0;column:log_output"`              // 日志下载方式 1 直接 2 上传服务器
	LogDst            string                          `json:"logDst" gorm:"type:varchar(255);not null;default:'';column:log_dst"`                 // 日志服务器上传位置
	DstServerId       int                             `json:"dstServerId" gorm:"type:bigint;not null;default:0;column:dst_server_id"`             // 日志服务器id
	SecretId          int                             `json:"secretId" gorm:"type:bigint;not null;default:0;column:secret_id"`                    // 日志服务器密码
	ShellType         int                             `json:"shellType" gorm:"type:tinyint(2);not null;default:0;column:shell_type"`              // shell类型
	ShellVars         string                          `json:"shellVars" gorm:"type:varchar(255);not null;default:'';column:shell_vars"`           // shell参数
	DeployInfos       string                          `json:"deployInfos" gorm:"type:text;column:deploy_infos"`                                   // 服务器上传位置
	Interactive       int                             `json:"interactive" gorm:"type:tinyint(2);not null;default:0;column:interactive"`           // 执行方式 1 命令 2 脚本
	CommandVarNumbers int                             `json:"commandVarNumbers" gorm:"type:int(5);not null;default:0;column:command_var_numbers"` // 命令参数个数
	LogSelect         int                             `json:"logSelect" gorm:"type:tinyint(2);not null;default:0;column:log_select"`
	DeployType        int                             `json:"deployType" gorm:"type:tinyint(2);not null;default:0;column:deploy_type"` // 下载方式 1 ftp/sftp 2 网盘
	BecomeUser        string                          `json:"becomeUser" gorm:"type:varchar(30);not null;default:'';column:become_user"`
	TaskDeployInfos   []TaskDeployInfo                `json:"taskDeployInfos" gorm:"-"`
	TargetIds         []int                           `json:"targetIds" gorm:"-"`
	TargetServers     []application.ApplicationServer `json:"targetServers" gorm:"-"`
	LastTask          Task                            `json:"lastTask" gorm:"-"`
	LogUploadServer   logUploadMdl.Server             `json:"logUploadServer" gorm:"-"`
	Secret            logUploadMdl.Secret             `json:"secret" gorm:"-"`
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
	if m.DstServerId > 0 {
		if err = tx.Model(&logUploadMdl.Server{}).Where("id = ?", m.DstServerId).Find(&m.LogUploadServer).Error; err != nil {
			global.GVA_LOG.Error("转换失败", zap.Any("err", err))
			return
		}
	}
	if m.SecretId > 0 {
		if err = tx.Model(&logUploadMdl.Secret{}).Where("id = ?", m.SecretId).Find(&m.Secret).Error; err != nil {
			global.GVA_LOG.Error("转换失败", zap.Any("err", err))
			return
		}
	}
	infos := make([]TaskDeployInfo, 0)
	if m.DeployInfos != "" {
		if err = json.Unmarshal([]byte(m.DeployInfos), &infos); err != nil {
			global.GVA_LOG.Error("转换失败", zap.Any("err", err))
			return
		}
		m.TaskDeployInfos = infos
	}
	return nil
}

type TaskTemplateSet struct {
	global.GVA_MODEL
	SystemId int    `json:"systemId" gorm:"type:bigint;not null;default:0;column:system_id" ` // 系统名称
	Name     string `json:"name" gorm:"type:varchar(100);not null;default:'';column:name"`    // 模板集名称
}

func (m *TaskTemplateSet) TableName() string {
	return "application_task_template_sets"
}

type TaskTemplateSetTemplate struct {
	global.GVA_MODEL
	TemplateId int `json:"templateId" gorm:"type:bigint;not null;default:0;column:template_id"` // task id
	SetId      int `json:"setId" gorm:"type:bigint;not null;default:0;column:set_id"`           // task id
	Seq        int `json:"seq" gorm:"type:int(4);not null;default:0;column:seq"`                // 排序
}

func (m *TaskTemplateSetTemplate) TableName() string {
	return "application_task_template_set_templates"
}

type TaskTemplateWithSeq struct {
	TaskTemplate
	Seq      int `json:"seq"`
	SeqInner int `json:"seqInner"`
}

type SetTask struct {
	global.GVA_MODEL
	SetId                int                     `json:"setId" gorm:"type:bigint;not null;default:0;column:set_id"`                // task id
	SystemUserId         int                     `json:"systemUserId" gorm:"type:bigint;not null;default:0;column:system_user_id"` // 执行人
	CurrentTaskId        int                     `json:"currentTaskId" gorm:"type:bigint;not null;default:0;column:current_task_id"`
	TotalSteps           int                     `json:"totalSteps" gorm:"type:int(4);not null;default:0;column:total_steps"`
	CurrentStep          int                     `json:"currentStep" gorm:"type:int(4);not null;default:0;column:current_step"`
	TemplatesString      string                  `json:"templatesString" gorm:"type:text;column:templates_string"`                    // 关联服务器id
	TasksString          string                  `json:"tasksString" gorm:"type:text;column:tasks_string"`                            // 关联服务器id
	ForceCorrect         int                     `json:"forceCorrect" gorm:"type:tinyint(2);not null;default:0;column:force_correct"` // 关联服务器id
	CurrentTaskIdsString string                  `json:"currentTaskIdsString" gorm:"type:text;column:current_task_ids_string"`
	Templates            [][]TaskTemplateWithSeq `json:"templates" gorm:"-"`
	Tasks                [][]Task                `json:"tasks" gorm:"-"`
	CurrentTaskIds       []int                   `json:"CurrentTaskIds" gorm:"-"`
	NeedRedo             int                     `json:"needRedo" gorm:"-"`
}

func (m *SetTask) TableName() string {
	return "application_set_tasks"
}

func (m *SetTask) AfterFind(tx *gorm.DB) (err error) {
	//templateIds := make([][]int, 0)
	//if m.TemplatesString != "" {
	//	if err = json.Unmarshal([]byte(m.TemplatesString), &templateIds); err != nil {
	//		global.GVA_LOG.Error("转换失败", zap.Any("err", err))
	//		return
	//	}
	//}
	//if len(templateIds) > 0 {
	//	for _, ids := range templateIds {
	//		templateWithSeqInner := make([]TaskTemplateWithSeq, 0)
	//		for index, id := range ids {
	//			template := TaskTemplate{}
	//			if err = tx.Model(&TaskTemplate{}).Where("id = ?", id).Find(&template).Error; err != nil {
	//				global.GVA_LOG.Error("转换失败", zap.Any("err", err))
	//				return
	//			}
	//			templateWithSeq := TaskTemplateWithSeq{
	//				TaskTemplate: template,
	//				Seq:          index,
	//			}
	//			templateWithSeqInner = append(templateWithSeqInner, templateWithSeq)
	//		}
	//		m.Templates = append(m.Templates, templateWithSeqInner)
	//	}
	//}
	setTemplates := make([]TaskTemplateSetTemplate, 0)
	if m.TemplatesString != "" {
		if err = json.Unmarshal([]byte(m.TemplatesString), &setTemplates); err != nil {
			global.GVA_LOG.Error("转换失败", zap.Any("err", err))
			return
		}
	}
	if len(setTemplates) > 0 {
		templateInner := make([]TaskTemplateWithSeq, 0)
		var currentSeq int
		for i, setTemplate := range setTemplates {
			template := TaskTemplate{}
			if err = tx.Model(&TaskTemplate{}).Where("id = ?", setTemplate.TemplateId).Find(&template).Error; err != nil {
				global.GVA_LOG.Error("查找失败", zap.Any("err", err))
				return
			}
			templateWithSeq := TaskTemplateWithSeq{
				TaskTemplate: template,
				Seq:          setTemplate.Seq,
				SeqInner:     int(setTemplate.ID),
			}
			if setTemplate.Seq != currentSeq && setTemplate.Seq > 0 && currentSeq > 0 {
				copySlice := make([]TaskTemplateWithSeq, len(templateInner)) // 创建一个与原切片长度相同的切片
				copy(copySlice, templateInner)
				m.Templates = append(m.Templates, copySlice)
				templateInner = templateInner[:0]
			}
			templateInner = append(templateInner, templateWithSeq)
			if i == len(setTemplates)-1 {
				m.Templates = append(m.Templates, templateInner)
			}
			currentSeq = setTemplate.Seq
		}
	}
	taskIds := make([][]int, 0)
	if m.TasksString != "" {
		if err = json.Unmarshal([]byte(m.TasksString), &taskIds); err != nil {
			global.GVA_LOG.Error("转换失败", zap.Any("err", err))
			return
		}
	}
	if len(taskIds) > 0 {
		for _, ids := range taskIds {
			taskInner := make([]Task, 0)
			for _, id := range ids {
				task := Task{}
				if err = tx.Model(&Task{}).Where("id = ?", id).Find(&task).Error; err != nil {
					global.GVA_LOG.Error("查找失败", zap.Any("err", err))
					return
				}
				taskInner = append(taskInner, task)
			}
			m.Tasks = append(m.Tasks, taskInner)
		}
	}
	if m.CurrentTaskIdsString != "" {
		if err = json.Unmarshal([]byte(m.CurrentTaskIdsString), &m.CurrentTaskIds); err != nil {
			global.GVA_LOG.Error("转换失败", zap.Any("err", err))
			return
		}
	}
	return nil
}

type TaskDeployInfo struct {
	DeployPath     string `json:"deployPath"`     // 服务器上传位置
	DownloadSource string `json:"downloadSource"` // 日志服务器下载位置
}
