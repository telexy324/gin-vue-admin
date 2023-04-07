package taskSvr

import (
	"encoding/json"
	"errors"
	sockets "github.com/flipped-aurora/gin-vue-admin/server/api/v1/socket"
	"github.com/flipped-aurora/gin-vue-admin/server/common"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/application"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/taskMdl"
	request2 "github.com/flipped-aurora/gin-vue-admin/server/model/taskMdl/request"
	"github.com/pkg/sftp"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"mime/multipart"
	"strings"
	"sync"
)

type TaskTemplatesService struct {
}

var TaskTemplatesServiceApp = new(TaskTemplatesService)

func (templateService *TaskTemplatesService) CreateTaskTemplate(template taskMdl.TaskTemplate) (taskMdl.TaskTemplate, error) {
	targetServersJson, err := json.Marshal(template.TargetIds)
	if err != nil {
		return template, err
	}
	s := string(targetServersJson)
	template.TargetServerIds = s
	err = global.GVA_DB.Create(&template).Error
	return template, err
}

func (templateService *TaskTemplatesService) UpdateTaskTemplate(template taskMdl.TaskTemplate) error {
	var oldTaskTemplate taskMdl.TaskTemplate
	targetServersJson, _ := json.Marshal(template.TargetIds)
	upDateMap := make(map[string]interface{})
	upDateMap["name"] = template.Name
	upDateMap["description"] = template.Description
	upDateMap["target_server_ids"] = targetServersJson
	upDateMap["mode"] = template.Mode
	upDateMap["command"] = template.Command
	upDateMap["script_path"] = template.ScriptPath
	upDateMap["last_task_id"] = template.LastTaskId
	upDateMap["sys_user"] = template.SysUser

	err := global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		db := tx.Where("id = ?", template.ID).Find(&oldTaskTemplate)
		txErr := db.Updates(upDateMap).Error
		if txErr != nil {
			global.GVA_LOG.Debug(txErr.Error())
			return txErr
		}
		return nil
	})
	return err
}

func (templateService *TaskTemplatesService) GetTaskTemplates(info request2.TaskTemplateSearch) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	var templates []taskMdl.TaskTemplate
	db := global.GVA_DB.Model(&taskMdl.TaskTemplate{})
	if info.Name != "" {
		name := strings.Trim(info.Name, " ")
		db = db.Where("`name` LIKE ?", "%"+name+"%")
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	err = db.Limit(limit).Offset(offset).Find(&templates).Error
	if err != nil {
		return
	}
	//err = templateService.FillTaskTemplates(templates)
	return err, templates, total
}

func (templateService *TaskTemplatesService) GetTaskTemplate(templateID float64) (template taskMdl.TaskTemplate, err error) {
	err = global.GVA_DB.Where("id =?", templateID).First(&template).Error
	if err != nil {
		return
	}
	//err = templateService.FillTaskTemplate(&template)
	return
}

func (templateService *TaskTemplatesService) DeleteTaskTemplate(templateID float64) error {
	err := global.GVA_DB.Where("id = ?", templateID).First(&taskMdl.TaskTemplate{}).Error
	if err != nil {
		return err
	}
	var template taskMdl.TaskTemplate
	return global.GVA_DB.Where("id = ?", templateID).First(&template).Delete(&template).Error
}

func (templateService *TaskTemplatesService) DeleteTaskTemplateByIds(ids request.IdsReq) error {
	return global.GVA_DB.Delete(&[]taskMdl.TaskTemplate{}, "id in ?", ids.Ids).Error
}

//func (d *SqlDb) GetTaskTemplateRefs(projectID int, templateID int) (db.ObjectReferrers, error) {
//	return d.getObjectRefs(projectID, db.TaskTemplateProps, templateID)
//}

func (templateService *TaskTemplatesService) Validate(tpl taskMdl.TaskTemplate) error {
	if tpl.Name == "" {
		return errors.New("template name can not be empty")
	}

	if tpl.Command == "" && tpl.ScriptPath == "" {
		return errors.New("template playbook can not be empty")
	}

	return nil
}

//func (templateService *TaskTemplatesService) FillTaskTemplates(templates []task.TaskTemplate) (err error) {
//	for i := range templates {
//		tpl := &templates[i]
//		var tasks []task.TaskWithTpl
//		e, iTasks, _ := TaskServiceApp.GetTaskTemplateTasks(tpl.ProjectID, int(tpl.ID), request.PageInfo{
//			Page:     1,
//			PageSize: 1,
//		})
//		tasks = iTasks.([]task.TaskWithTpl)
//		if e != nil {
//			return e
//		}
//		if len(tasks) > 0 {
//			tpl.LastTask = &tasks[0]
//		}
//	}
//	return
//}
//
//func (templateService *TaskTemplatesService) FillTaskTemplate(template *task.TaskTemplate) (err error) {
//	if template.VaultKeyID != nil {
//		template.VaultKey, err = KeyServiceApp.GetAccessKey(float64(template.ProjectID), float64(*template.VaultKeyID))
//	}
//
//	if err != nil {
//		return
//	}
//
//	err = templateService.FillTaskTemplates([]task.TaskTemplate{*template})
//
//	if err != nil {
//		return
//	}
//
//	if template.SurveyVarsJSON != nil {
//		err = json.Unmarshal([]byte(*template.SurveyVarsJSON), &template.SurveyVars)
//	}
//
//	return
//}

func (templateService *TaskTemplatesService) CheckScript(s application.ApplicationServer, needDetail bool, sshClient *common.SSHClient, template taskMdl.TaskTemplate) (exist bool, output string, err error) {
	var command string
	command = `[ -f ` + template.ScriptPath + ` ] && echo yes || echo no`
	output, err = sshClient.CommandSingle(command)
	global.GVA_LOG.Info(strings.Trim(output, " "))
	if err != nil || strings.Trim(output, " ") == "no" {
		global.GVA_LOG.Error("judge script exist: ", zap.String("server IP: ", s.ManageIp), zap.Any("err", err))
		return
	}
	exist = true
	if needDetail {
		command = `cat ` + template.ScriptPath
		output, err = sshClient.CommandSingle(command)
		if err != nil {
			global.GVA_LOG.Error("judge script exist: ", zap.String("server IP: ", s.ManageIp), zap.Any("err", err))
		}
	}
	return
}

func (templateService *TaskTemplatesService) DownloadScript(ID float64, server application.ApplicationServer) (file *sftp.File, err error) {
	template, err := templateService.GetTaskTemplate(ID)
	if err != nil {
		return
	}
	sshClient, err := common.FillSSHClient(server.ManageIp, template.SysUser, "", server.SshPort)
	err = sshClient.GenerateClient()
	if err != nil {
		global.GVA_LOG.Error("create ssh client failed: ", zap.String("server IP: ", server.ManageIp), zap.Any("err", err))
		return
	}
	return sshClient.Download(template.ScriptPath)
}

func (templateService *TaskTemplatesService) UploadScript(ID int, file multipart.File, scriptPath string, userID uint) (failedIPs []string, err error) {
	template, err := templateService.GetTaskTemplate(float64(ID))
	if err != nil {
		return
	}
	if scriptPath != template.ScriptPath {
		template.ScriptPath = scriptPath
		if err = templateService.UpdateTaskTemplate(template); err != nil {
			return
		}
	}
	wg := &sync.WaitGroup{}
	failedChan := make(chan string, len(template.TargetServers))
	for _, server := range template.TargetServers {
		wg.Add(1)
		go func(w *sync.WaitGroup, s application.ApplicationServer, f chan string) {
			var er error
			defer w.Done()
			defer func() {
				var status string
				if er == nil {
					status = "success"
				} else {
					status = "exception"
				}
				b, e := json.Marshal(&map[string]interface{}{
					"type":       "uploadScript",
					"manageIp":   s.ManageIp,
					"ID":         s.ID,
					"status":     status,
					"templateID": ID,
				})
				if e != nil {
					global.GVA_LOG.Error(err.Error())
					return
				}
				sockets.Message(int(userID), b)
			}()
			sshClient, er := common.FillSSHClient(s.ManageIp, template.SysUser, "", s.SshPort)
			er = sshClient.GenerateClient()
			if er != nil {
				global.GVA_LOG.Error("upload script failed on create ssh client: ", zap.String("server IP: ", s.ManageIp), zap.Any("err", er))
				f <- s.ManageIp
				return
			}

			er = sshClient.Upload(file, scriptPath)
			if er != nil {
				global.GVA_LOG.Error("upload script failed on upload: ", zap.String("server IP: ", s.ManageIp), zap.Any("err", er))
				f <- s.ManageIp
				return
			}
			return
		}(wg, server, failedChan)
	}
	wg.Wait()
	close(failedChan)
	failedIPs = make([]string, 0)
	for ip := range failedChan {
		failedIPs = append(failedIPs, ip)
	}
	return
}
