package scheduleSvr

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/scheduleMdl"
	"gorm.io/gorm"
)

type ScheduleService struct {
}

var scheduleServiceApp = new(ScheduleService)

func (scheduleService *ScheduleService) CreateSchedule(schedule scheduleMdl.Schedule) (scheduleMdl.Schedule, error) {
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
	upDateMap["cron_format"] = schedule.CronFormat

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

func (scheduleService *ScheduleService) GetSchedule(scheduleID float64) (template scheduleMdl.Schedule, err error) {
	err = global.GVA_DB.Where("id =?", scheduleID).First(&template).Error
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