package taskSvr

import (
	"errors"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/taskMdl"
	"github.com/flipped-aurora/gin-vue-admin/server/model/taskMdl/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/taskMdl/response"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sort"
	"strconv"
	"time"
)

type TaskService struct {
}

var TaskServiceApp = new(TaskService)

func (taskService *TaskService) CreateTask(task taskMdl.Task) (taskMdl.Task, error) {
	err := global.GVA_DB.Create(&task).Error
	return task, err
}

func (taskService *TaskService) UpdateTask(t taskMdl.Task) error {
	var oldTask taskMdl.Task
	upDateMap := make(map[string]interface{})
	upDateMap["status"] = t.Status
	upDateMap["begin_time"] = t.BeginTime
	upDateMap["end_time"] = t.EndTime

	err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		db := tx.Where("id = ?", t.ID).Find(&oldTask)
		txErr := db.Updates(upDateMap).Error
		if txErr != nil {
			global.GVA_LOG.Debug(txErr.Error())
			return txErr
		}
		return nil
	})
	return err
}

func (taskService *TaskService) CreateTaskOutput(output taskMdl.TaskOutput) (taskMdl.TaskOutput, error) {
	err := global.GVA_DB.Create(&output).Error
	return output, err
}

func (taskService *TaskService) getTasks(templateID int, info request.GetTaskByTemplateId) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&taskMdl.Task{}) //.Preload("User")
	//if templateID == nil {
	//	db = db.Preload("Template")
	//} else {
	//	db = db.Preload("Template").Where("template_id=?", templateID)
	//}
	//db.Order("created desc, id desc")
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	var Tasks []taskMdl.Task
	//var TaskWithTpls []task.TaskWithTpl
	if templateID > 0 {
		db = db.Find("where template_id = ?", templateID)
	}
	//err = db.Limit(limit).Offset(offset).Find(&Tasks).Error
	db = db.Limit(limit).Offset(offset)
	if info.OrderKey != "" {
		var OrderStr string
		if info.Desc {
			OrderStr = info.OrderKey + " desc"
		} else {
			OrderStr = info.OrderKey
		}
		err = db.Order(OrderStr).Find(&Tasks).Error
	} else {
		err = db.Order("id").Find(&Tasks).Error
	}

	//for _, t := range Tasks {
	//	taskWithTpl := task.TaskWithTpl{
	//		Task:             t,
	//		TemplatePlaybook: t.Template.Playbook,
	//		TemplateAlias:    t.Template.Name,
	//		TemplateType:     t.Template.Type,
	//		UserName:         &t.User.Username,
	//	}
	//	TaskWithTpls = append(TaskWithTpls, taskWithTpl)
	//}
	return err, Tasks, total
}

func (taskService *TaskService) GetTask(taskID int) (task taskMdl.Task, err error) {
	err = global.GVA_DB.Where("id = ?", taskID).First(&task).Error
	return
}

func (taskService *TaskService) GetTemplateTasks(templateID int, info request.GetTaskByTemplateId) (err error, list interface{}, total int64) {
	return taskService.getTasks(templateID, info)
}

func (taskService *TaskService) GetProjectTasks(info request.GetTaskByTemplateId) (err error, list interface{}, total int64) {
	return taskService.getTasks(0, info)
}

func (taskService *TaskService) DeleteTaskWithOutputs(taskID int) (err error) {
	// check if task exists in the project
	_, err = taskService.GetTask(taskID)
	if err != nil {
		return
	}
	var t taskMdl.Task
	var taskOutputs []taskMdl.TaskOutput
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		txErr := tx.Where("task_id = ?", taskID).Find(&taskOutputs).Delete(&taskOutputs).Error
		if txErr != nil {
			return txErr
		}
		txErr = tx.Where("id = ?", taskID).Find(&t).Delete(&t).Error
		if txErr != nil {
			return txErr
		}
		return nil
	})
	return
}

func (taskService *TaskService) GetTaskOutputs(taskID int) (output []taskMdl.TaskOutput, err error) {
	// check if task exists in the project
	_, err = taskService.GetTask(taskID)
	if err != nil {
		return
	}
	err = global.GVA_DB.Where("task_id = ?", taskID).Order("manage_ip").Order("record_time").Find(&output).Error
	return
}

//func (taskService *TaskService) GetIncomingVersion(t task.Task) *string {
//	if t.BuildTaskID == nil {
//		return nil
//	}
//
//	buildTask, err := taskService.GetTask(t.ProjectID, *t.BuildTaskID)
//
//	if err != nil {
//		return nil
//	}
//
//	tpl, err := TemplatesServiceApp.GetTemplate(float64(t.ProjectID), float64(buildTask.TemplateID))
//	if err != nil {
//		return nil
//	}
//
//	if tpl.Type == t.TemplateBuild {
//		return buildTask.Version
//	}
//
//	return taskService.GetIncomingVersion(buildTask)
//}

//func (taskService *TaskService) ValidateNewTask(template task.Template) error {
//	switch template.Type {
//	case task.TemplateBuild:
//	case task.TemplateDeploy:
//	case task.TemplateTask:
//	}
//	return nil
//}

//func (taskService *TaskService) Fill(t task.TaskWithTpl) error {
//	if t.Task.BuildTaskID != nil {
//		build, err := taskService.GetTask(t.Task.ProjectID, *t.Task.BuildTaskID)
//		if err == gorm.ErrRecordNotFound {
//			return nil
//		}
//		if err != nil {
//			return err
//		}
//		t.BuildTask = &build
//	}
//	return nil
//}

func (taskService *TaskService) GetTaskDashboardInfo() (output []response.TaskDashboardInfo) {
	for i := 0; i < 12; i++ {
		now := time.Now().AddDate(0, 0, -i)
		day := strconv.Itoa(int(now.Month())) + "月" + strconv.Itoa(now.Day()) + "日"
		var count int64
		timeEnd := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
		timeStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		if err := global.GVA_DB.Model(&taskMdl.Task{}).Where("created_at >= ? and created_at < ?", timeStart, timeEnd).Count(&count).Error; err != nil {
			global.GVA_LOG.Error("get task dashboard info failed ", zap.String("date ", day), zap.Any("error ", err))
		}
		output = append(output, response.TaskDashboardInfo{
			Date:   strconv.Itoa(int(now.Month())) + "月" + strconv.Itoa(now.Day()) + "日",
			Number: count,
		})
	}
	return
}

func (taskService *TaskService) GetSetTasks(info request.GetTaskBySetTaskIdWithSeq) (err error, list interface{}, total int64) {
	var Tasks []taskMdl.Task
	//var TaskWithTpls []task.TaskWithTpl
	if info.SetTaskId <= 0 {
		err = errors.New("set task 不能为空")
		return
	}
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&taskMdl.Task{}) //.Preload("User")
	db = db.Where("set_task_id = ? and set_task_outer_seq = ?", info.SetTaskId, info.CurrentSeq)
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	if total <= 0 {
		//setTaskTemplates := make([]taskMdl.TaskTemplateSetTemplate, 0)
		//if err = global.GVA_DB.Where("seq = ?", info.CurrentSeq).Find(&setTaskTemplates).Error; err != nil {
		//	return
		//}
		//for _, setTaskTemplate := range setTaskTemplates {
		//
		//}
		var setTask taskMdl.SetTask
		if err = global.GVA_DB.Where("id = ?", info.SetTaskId).First(&setTask).Error; err != nil {
			return
		}
		if setTask.TotalSteps == setTask.CurrentStep {
			err = errors.New("任务已结束")
			return
		}
		sort.Slice(setTask.Templates[info.CurrentSeq], func(i, j int) bool {
			return setTask.Templates[info.CurrentSeq][i].SeqInner < setTask.Templates[info.CurrentSeq][j].SeqInner
		})
		total = int64(len(setTask.Templates[info.CurrentSeq]))
		if int64(offset) <= total {
			var targetTemplates []taskMdl.TaskTemplateWithSeq
			if int64(offset+limit) > total {
				targetTemplates = setTask.Templates[info.CurrentSeq][offset:total]
			} else {
				targetTemplates = setTask.Templates[info.CurrentSeq][offset : offset+limit]
			}
			for _, template := range targetTemplates {
				Tasks = append(Tasks, taskMdl.Task{
					TemplateId:      int(template.ID),
					SetTaskId:       int(info.SetTaskId),
					SetTaskInnerSeq: template.SeqInner,
					SetTaskOuterSeq: template.Seq,
				})
			}
		}
		return nil, Tasks, total
	}
	//err = db.Limit(limit).Offset(offset).Find(&Tasks).Error
	db = db.Limit(limit).Offset(offset)
	err = db.Order("id").Find(&Tasks).Error
	return err, Tasks, total
}
