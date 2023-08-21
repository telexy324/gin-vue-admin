package scheduleSvr

import (
	"encoding/json"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/scheduleMdl"
	scheduleReq "github.com/flipped-aurora/gin-vue-admin/server/model/scheduleMdl/request"
	"gorm.io/gorm"
)

type ScheduleService struct {
}

var scheduleServiceApp = new(ScheduleService)

func (scheduleService *ScheduleService) CreateSchedule(schedule scheduleMdl.Schedule) (scheduleMdl.Schedule, error) {
	if len(schedule.CommandVars) > 0 {
		vars, _ := json.Marshal(schedule.CommandVars)
		schedule.CommandVar = string(vars)
	}
	err := global.GVA_DB.Create(&schedule).Error
	return schedule, err
}

func (scheduleService *ScheduleService) SetScheduleLastCommitHash(scheduleID int, lastCommentHash string) error {
	oldSchedule, err := scheduleService.GetSchedule(float64(scheduleID))
	if err != nil {
		return err
	}
	upDateMap := make(map[string]interface{})
	upDateMap["last_commit_hash"] = lastCommentHash
	return global.GVA_DB.Model(&oldSchedule).Updates(upDateMap).Error
}

func (scheduleService *ScheduleService) UpdateSchedule(schedule scheduleMdl.Schedule) error {
	var oldSchedule scheduleMdl.Schedule
	upDateMap := make(map[string]interface{})
	upDateMap["template_id"] = schedule.TemplateID
	upDateMap["cron_format"] = schedule.CronFormat
	upDateMap["valid"] = schedule.Valid
	if len(schedule.CommandVars) > 0 {
		vars, _ := json.Marshal(schedule.CommandVars)
		schedule.CommandVar = string(vars)
	}
	upDateMap["command_var"] = schedule.CommandVar

	err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		db := tx.Where("id = ?", schedule.ID).Find(&oldSchedule)
		txErr := db.Updates(upDateMap).Error
		if txErr != nil {
			global.GVA_LOG.Debug(txErr.Error())
			return txErr
		}
		return nil
	})
	return err
}

func (scheduleService *ScheduleService) GetSchedule(scheduleID float64) (schedule scheduleMdl.Schedule, err error) {
	err = global.GVA_DB.Where("id =?", scheduleID).First(&schedule).Error
	return
}

func (scheduleService *ScheduleService) DeleteSchedule(scheduleID float64) error {
	err := global.GVA_DB.Where("id = ?", scheduleID).First(&scheduleMdl.Schedule{}).Error
	if err != nil {
		return err
	}
	var schedule scheduleMdl.Schedule
	return global.GVA_DB.Where("id = ?", scheduleID).First(&schedule).Delete(&schedule).Error
}

func (scheduleService *ScheduleService) DeleteScheduleByIds(ids request.IdsReq) error {
	return global.GVA_DB.Delete(&[]scheduleMdl.Schedule{}, "id in ?", ids.Ids).Error
}

func (scheduleService *ScheduleService) GetSchedules() (schedules []scheduleMdl.Schedule, err error) {
	err = global.GVA_DB.Where("cron_format != ''").Find(&schedules).Error
	return
}

func (scheduleService *ScheduleService) GetTemplateSchedules(templateID float64) (schedules []scheduleMdl.Schedule, err error) {
	err = global.GVA_DB.Where("template_id=?", templateID).Find(&schedules).Error
	return
}

func (scheduleService *ScheduleService) SetScheduleCommitHash(scheduleID int, hash string) error {
	return scheduleService.SetScheduleLastCommitHash(scheduleID, hash)
}

func (scheduleService *ScheduleService) GetScheduleList(info scheduleReq.GetScheduleList, templateIDs []int) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	db := global.GVA_DB.Model(&scheduleMdl.Schedule{})
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	var schedules []scheduleMdl.Schedule
	db = db.Where("template_id in ?", templateIDs)
	//err = db.Limit(limit).Offset(offset).Find(&schedules).Error
	db = db.Limit(limit).Offset(offset)
	if info.OrderKey != "" {
		var OrderStr string
		if info.Desc {
			OrderStr = info.OrderKey + " desc"
		} else {
			OrderStr = info.OrderKey
		}
		err = db.Order(OrderStr).Find(&schedules).Error
	} else {
		err = db.Order("id").Find(&schedules).Error
	}
	return err, schedules, total
}
