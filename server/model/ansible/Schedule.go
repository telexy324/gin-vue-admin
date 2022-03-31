package ansible

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"gorm.io/gorm"
)

type Schedule struct {
	ID             int     `db:"id" json:"id"`
	ProjectID      int     `db:"project_id" json:"project_id"`
	TemplateID     int     `db:"template_id" json:"template_id"`
	CronFormat     string  `db:"cron_format" json:"cron_format"`
	RepositoryID   *int    `db:"repository_id" json:"repository_id"`
	LastCommitHash *string `db:"last_commit_hash" json:"-"`
}

func CreateSchedule(schedule Schedule) (Schedule, error) {
	err := global.GVA_DB.Create(&schedule).Error
	return schedule, err
}

func (m *Schedule) SetScheduleLastCommitHash(projectID int, scheduleID int, lastCommentHash string) error {
	oldSchedule, err := m.GetSchedule(projectID, scheduleID)
	if err != nil {
		return err
	}
	upDateMap := make(map[string]interface{})
	upDateMap["last_commit_hash"] = lastCommentHash
	return global.GVA_DB.Model(&oldSchedule).Updates(upDateMap).Error
}

func (m *Schedule) UpdateSchedule(schedule Schedule) error {
	var oldSchedule Schedule
	upDateMap := make(map[string]interface{})
	upDateMap["cron_format"] = schedule.CronFormat
	upDateMap["repository_id"] = schedule.RepositoryID

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

func (m *Schedule) GetSchedule(projectID int, scheduleID int) (template Schedule, err error) {
	err = global.GVA_DB.Where("project_id=? and id =?", projectID, scheduleID).First(&template).Error
	return
}

func (m *Schedule) DeleteSchedule(projectID int, scheduleID int) error {
	err := global.GVA_DB.Where("id = ? and project_id=? ", scheduleID, projectID).First(&Schedule{}).Error
	if err != nil {
		return err
	}
	var schedule Schedule
	return global.GVA_DB.Where("id = ? and project_id=? ", scheduleID, projectID).First(&schedule).Delete(&schedule).Error
}

func (m *Schedule) GetSchedules() (schedules []Schedule, err error) {
	err = global.GVA_DB.Where("cron_format != ''").Find(&schedules).Error
	return
}

func (m *Schedule) GetTemplateSchedules(projectID int, templateID int) (schedules []Schedule, err error) {
	err = global.GVA_DB.Where("project_id = ? and template_id=?", projectID, templateID).Find(&schedules).Error
	return
}

func (m *Schedule) SetScheduleCommitHash(projectID int, scheduleID int, hash string) error {
	return m.SetScheduleLastCommitHash(projectID, scheduleID, hash)
}
