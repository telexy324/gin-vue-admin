package scheduleSvr

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ansible"
	"gorm.io/gorm"
)

type SchedulesService struct {
}

var SchedulesServiceApp = new(SchedulesService)

func (schedulesService *SchedulesService) CreateSchedule(schedule ansible.Schedule) (ansible.Schedule, error) {
	err := global.GVA_DB.Create(&schedule).Error
	return schedule, err
}

func (schedulesService *SchedulesService) SetScheduleLastCommitHash(projectID int, scheduleID int, lastCommentHash string) error {
	oldSchedule, err := schedulesService.GetSchedule(float64(projectID), float64(scheduleID))
	if err != nil {
		return err
	}
	upDateMap := make(map[string]interface{})
	upDateMap["last_commit_hash"] = lastCommentHash
	return global.GVA_DB.Model(&oldSchedule).Updates(upDateMap).Error
}

func (schedulesService *SchedulesService) UpdateSchedule(schedule ansible.Schedule) error {
	var oldSchedule ansible.Schedule
	upDateMap := make(map[string]interface{})
	upDateMap["cron_format"] = schedule.CronFormat

	err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		db := tx.Where("id = ? and project_id = ?", schedule.ID, schedule.ProjectID).Find(&oldSchedule)
		txErr := db.Updates(upDateMap).Error
		if txErr != nil {
			global.GVA_LOG.Debug(txErr.Error())
			return txErr
		}
		return nil
	})
	return err
}

func (schedulesService *SchedulesService) GetSchedule(projectID float64, scheduleID float64) (template ansible.Schedule, err error) {
	err = global.GVA_DB.Where("project_id=? and id =?", projectID, scheduleID).First(&template).Error
	return
}

func (schedulesService *SchedulesService) DeleteSchedule(projectID float64, scheduleID float64) error {
	err := global.GVA_DB.Where("id = ? and project_id=? ", scheduleID, projectID).First(&ansible.Schedule{}).Error
	if err != nil {
		return err
	}
	var schedule ansible.Schedule
	return global.GVA_DB.Where("id = ? and project_id=? ", scheduleID, projectID).First(&schedule).Delete(&schedule).Error
}

func (schedulesService *SchedulesService) GetSchedules() (schedules []ansible.Schedule, err error) {
	err = global.GVA_DB.Where("cron_format != ''").Find(&schedules).Error
	return
}

func (schedulesService *SchedulesService) GetTemplateSchedules(projectID float64, templateID float64) (schedules []ansible.Schedule, err error) {
	err = global.GVA_DB.Where("project_id = ? and template_id=?", projectID, templateID).Find(&schedules).Error
	return
}

func (schedulesService *SchedulesService) SetScheduleCommitHash(projectID int, scheduleID int, hash string) error {
	return schedulesService.SetScheduleLastCommitHash(projectID, scheduleID, hash)
}
