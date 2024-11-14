package taskSvr

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	sockets "github.com/flipped-aurora/gin-vue-admin/server/api/v1/socket"
	"github.com/flipped-aurora/gin-vue-admin/server/common"
	"github.com/flipped-aurora/gin-vue-admin/server/consts"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/application"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/taskMdl"
	request2 "github.com/flipped-aurora/gin-vue-admin/server/model/taskMdl/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/taskMdl/response"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"io"
	"mime/multipart"
	"strings"
	"sync"
)

type TaskTemplatesService struct {
}

var TaskTemplatesServiceApp = new(TaskTemplatesService)

func (templateService *TaskTemplatesService) CreateTaskTemplate(template taskMdl.TaskTemplate) (taskMdl.TaskTemplate, error) {
	if !errors.Is(global.GVA_DB.Where("name = ?", template.Name).First(&taskMdl.TaskTemplate{}).Error, gorm.ErrRecordNotFound) {
		return template, errors.New("存在name，请修改name")
	}
	targetServersJson, err := json.Marshal(template.TargetIds)
	if err != nil {
		return template, err
	}
	template.TargetServerIds = string(targetServersJson)
	deployJson, err := json.Marshal(template.TaskDeployInfos)
	if err != nil {
		return template, err
	}
	template.DeployInfos = string(deployJson)
	err = global.GVA_DB.Create(&template).Error
	return template, err
}

func (templateService *TaskTemplatesService) UpdateTaskTemplate(template taskMdl.TaskTemplate) error {
	var oldTaskTemplate taskMdl.TaskTemplate
	targetServersJson, err := json.Marshal(template.TargetIds)
	if err != nil {
		return err
	}
	deployJson, err := json.Marshal(template.TaskDeployInfos)
	if err != nil {
		return err
	}
	upDateMap := make(map[string]interface{})
	upDateMap["name"] = template.Name
	upDateMap["description"] = template.Description
	upDateMap["target_server_ids"] = string(targetServersJson)
	upDateMap["mode"] = template.Mode
	upDateMap["command"] = template.Command
	upDateMap["script_path"] = template.ScriptPath
	upDateMap["last_task_id"] = template.LastTaskId
	upDateMap["sys_user"] = template.SysUser
	upDateMap["system_id"] = template.SystemId
	upDateMap["execute_type"] = template.ExecuteType
	upDateMap["log_path"] = template.LogPath
	upDateMap["script_hash"] = template.ScriptHash
	upDateMap["log_output"] = template.LogOutput
	upDateMap["log_dst"] = template.LogDst
	upDateMap["dst_server_id"] = template.DstServerId
	upDateMap["secret_id"] = template.SecretId
	upDateMap["shell_type"] = template.ShellType
	upDateMap["shell_vars"] = template.ShellVars
	upDateMap["deploy_infos"] = string(deployJson)
	upDateMap["interactive"] = template.Interactive
	upDateMap["command_var_numbers"] = template.CommandVarNumbers
	upDateMap["log_select"] = template.LogSelect
	upDateMap["deploy_type"] = template.DeployType
	upDateMap["become_user"] = template.BecomeUser

	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		db := tx.Where("id = ?", template.ID).Find(&oldTaskTemplate)
		if oldTaskTemplate.Name != template.Name {
			if !errors.Is(tx.Where("id <> ? AND name = ?", template.ID, template.Name).First(&taskMdl.TaskTemplate{}).Error, gorm.ErrRecordNotFound) {
				global.GVA_LOG.Debug("存在相同name修改失败")
				return errors.New("存在相同name修改失败")
			}
		}
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
	db = db.Where("`system_id` IN ?", info.SystemIDs)
	if info.ExecuteType > 0 {
		db = db.Where("execute_type = ?", info.ExecuteType)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	//err = db.Limit(limit).Offset(offset).Find(&templates).Error
	db = db.Limit(limit).Offset(offset)
	if info.OrderKey != "" {
		var OrderStr string
		if info.Desc {
			OrderStr = info.OrderKey + " desc"
		} else {
			OrderStr = info.OrderKey
		}
		err = db.Order(OrderStr).Find(&templates).Error
	} else {
		err = db.Order("id").Find(&templates).Error
	}
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

//func (templateService *TaskTemplatesService) CheckScript(s application.ApplicationServer, needDetail bool, sshClient *common.SSHClient, template taskMdl.TaskTemplate) (exist bool, output string, err error) {
//	var command string
//	command = `[ -f ` + template.ScriptPath + ` ] && echo yes || echo no`
//	output, err = sshClient.CommandSingle(command)
//	if err != nil || strings.Trim(output, " ") == "no" || strings.Trim(output, "\n") == "no" {
//		global.GVA_LOG.Error("judge script exist: ", zap.String("server IP: ", s.ManageIp), zap.Any("err", err))
//		return
//	}
//	exist = true
//	if needDetail {
//		command = `cat ` + template.ScriptPath
//		output, err = sshClient.CommandSingle(command)
//		if err != nil {
//			global.GVA_LOG.Error("judge script exist: ", zap.String("server IP: ", s.ManageIp), zap.Any("err", err))
//		}
//	}
//	return
//}

func (templateService *TaskTemplatesService) CheckScript(s application.ApplicationServer, sshClient *common.SSHClient, template taskMdl.TaskTemplate) (exist bool, err error) {
	var command string
	command = `[ -f ` + template.ScriptPath + ` ] && echo yes || echo no`
	output, err := sshClient.CommandSingle(command)
	if err != nil || strings.Trim(output, " ") == "no" || strings.Trim(output, "\n") == "no" {
		global.GVA_LOG.Error("judge script exist: ", zap.String("server IP: ", s.ManageIp), zap.Any("err", err))
		return
	}
	ftp, err := sshClient.NewSftp()
	if err != nil {
		return
	}
	defer ftp.Close()

	remote, err := ftp.Open(template.ScriptPath)
	if err != nil {
		return
	}
	defer remote.Close()

	hash := md5.New()
	_, err = io.Copy(hash, remote)
	if err != nil {
		return
	}
	hashString := hex.EncodeToString(hash.Sum(nil))
	if hashString == template.ScriptHash {
		exist = true
	}
	return
}

func (templateService *TaskTemplatesService) CheckScriptDetail(sshClient *common.SSHClient, template taskMdl.TaskTemplate) (output string, err error) {
	command := `cat ` + template.ScriptPath
	return sshClient.CommandSingle(command)
}

func (templateService *TaskTemplatesService) DownloadScript(ID float64, server application.ApplicationServer) (fio io.ReadCloser, err error) {
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
	fileByte, _ := io.ReadAll(file)
	for _, server := range template.TargetServers {
		wg.Add(1)
		go func(w *sync.WaitGroup, s application.ApplicationServer, f chan string) {
			var er error
			var sshClient common.SSHClient
			defer w.Done()
			defer func() {
				if sshClient.Client != nil {
					sshClient.Client.Close()
				}
			}()
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
			sshClient, er = common.FillSSHClient(s.ManageIp, template.SysUser, "", s.SshPort)
			er = sshClient.GenerateClient()
			if er != nil {
				global.GVA_LOG.Error("upload script failed on create ssh client: ", zap.String("server IP: ", s.ManageIp), zap.Any("err", er))
				f <- s.ManageIp
				return
			}
			copied := bytes.NewBuffer(fileByte)
			er = sshClient.Upload(copied, scriptPath)
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

//@author: [telexy324](https://github.com/telexy324)
//@function: AddSet
//@description: 添加模板集
//@param: set taskMdl.TaskTemplateSet
//@return: error

func (templateService *TaskTemplatesService) AddSet(addSetRequest request2.AddSet) (err error) {
	if !errors.Is(global.GVA_DB.Where("name = ?", addSetRequest.Name).First(&taskMdl.TaskTemplateSet{}).Error, gorm.ErrRecordNotFound) {
		return errors.New("存在重复name，请修改name")
	}
	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		setMdl := &taskMdl.TaskTemplateSet{
			SystemId: addSetRequest.SystemId,
			Name:     addSetRequest.Name,
		}
		if txErr := tx.Create(&setMdl).Error; txErr != nil {
			global.GVA_LOG.Error("添加系统失败", zap.Any("err", err))
			return txErr
		}
		if addSetRequest.Templates != nil && len(addSetRequest.Templates) > 0 {
			for _, t := range addSetRequest.Templates {
				template := &taskMdl.TaskTemplate{}
				template.ID = uint(t.TemplateId)
				if err = global.GVA_DB.Find(template).Error; err != nil {
					global.GVA_LOG.Error("模板不存在", zap.Any("err", err))
					continue
				}
				setTemplate := &taskMdl.TaskTemplateSetTemplate{
					TemplateId: t.TemplateId,
					SetId:      int(setMdl.ID),
					Seq:        t.Seq,
				}
				if err = global.GVA_DB.Create(&setTemplate).Error; err != nil {
					global.GVA_LOG.Error("添加模板集模板失败", zap.Any("err", err))
				}
			}
		}
		return nil
	})
	return
}

//@author: [telexy324](https://github.com/telexy324)
//@function: DeleteSet
//@description: 删除模板集
//@param: id float64
//@return: err error

func (templateService *TaskTemplatesService) DeleteSet(id float64) (err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&taskMdl.TaskTemplateSet{}).Error
	if err != nil {
		return
	}
	var existSet taskMdl.TaskTemplateSet
	if err = global.GVA_DB.Where("id = ?", id).First(&existSet).Delete(&existSet).Error; err != nil {
		return err
	}
	var setTemplates []taskMdl.TaskTemplateSetTemplate
	err = global.GVA_DB.Where("set_id = ?", id).Find(&setTemplates).Delete(&setTemplates).Error
	if err != nil {
		return err
	}
	return err
}

//@author: [telexy324](https://github.com/telexy324)
//@function: DeleteSetByIds
//@description: 批量删除模板集
//@param: taskMdl.TaskTemplateSet []taskMdl.TaskTemplateSet
//@return: err error

func (templateService *TaskTemplatesService) DeleteSetByIds(ids request.IdsReq) (err error) {
	if ids.Ids == nil || len(ids.Ids) <= 0 {
		return
	}
	for _, id := range ids.Ids {
		if err = templateService.DeleteSet(float64(id)); err != nil {
			return
		}
	}
	return
}

//@author: [telexy324](https://github.com/telexy324)
//@function: UpdateSet
//@description: 更新模板集
//@param: system taskMdl.TaskTemplateSet
//@return: err error

func (templateService *TaskTemplatesService) UpdateSet(addSetRequest request2.AddSet) (err error) {
	var oldSet taskMdl.TaskTemplateSet
	upDateMap := make(map[string]interface{})
	upDateMap["name"] = addSetRequest.Name
	upDateMap["system_id"] = addSetRequest.SystemId

	err = global.GVA_DB.Transaction(func(tx *gorm.DB) error {
		db := tx.Where("id = ?", addSetRequest.ID).Find(&oldSet)
		if oldSet.Name != addSetRequest.Name {
			if err = tx.Where("id <> ? AND name = ?", addSetRequest.ID, addSetRequest.Name).First(&taskMdl.TaskTemplateSet{}).Error; err != nil && err != gorm.ErrRecordNotFound {
				global.GVA_LOG.Debug("存在相同name修改失败")
				return errors.New("存在相同name修改失败")
			}
		}
		txErr := db.Updates(upDateMap).Error
		if txErr != nil {
			global.GVA_LOG.Debug(txErr.Error())
			return txErr
		}
		if addSetRequest.Templates != nil && len(addSetRequest.Templates) > 0 {
			existSetTemplates := make([]taskMdl.TaskTemplateSetTemplate, 0)
			if txErr = tx.Where("set_id = ?", addSetRequest.ID).Find(&existSetTemplates).Delete(&existSetTemplates).Error; txErr != nil {
				return txErr
			}

			for _, tmp := range addSetRequest.Templates {
				if txErr = tx.Create(&taskMdl.TaskTemplateSetTemplate{
					TemplateId: tmp.TemplateId,
					SetId:      tmp.SetId,
					Seq:        tmp.Seq,
				}).Error; txErr != nil {
					return txErr
				}
			}
		}
		return nil
	})
	return err
}

//@author: [telexy324](https://github.com/telexy324)
//@function: GetSetById
//@description: 返回当前选中system
//@param: id float64
//@return: err error, set taskMdl.TaskTemplateSet

func (templateService *TaskTemplatesService) GetSetById(id float64) (err error, set taskMdl.TaskTemplateSet, templateRes []response.TaskTemplateSetTemplateResponse) {
	if err = global.GVA_DB.Where("id = ?", id).First(&set).Error; err != nil {
		return
	}
	templates := make([]taskMdl.TaskTemplateSetTemplate, 0)
	if err = global.GVA_DB.Where("set_id = ?", id).Order("seq").Find(&templates).Error; err != nil {
		return
	}
	templateRes = make([]response.TaskTemplateSetTemplateResponse, 0)
	for _, t := range templates {
		var template taskMdl.TaskTemplate
		if err = global.GVA_DB.Where("id = ?", t.TemplateId).Find(&template).Error; err != nil {
			return
		}
		res := response.TaskTemplateSetTemplateResponse{
			TaskTemplateSetTemplate: t,
			TemplateName:            template.Name,
		}
		templateRes = append(templateRes, res)
	}
	return
}

//@author: [telexy324](https://github.com/telexy324)
//@function: GetSystemList
//@description: 获取系统分页
//@return: err error, list interface{}, total int64

func (templateService *TaskTemplatesService) GetSetList(info request2.TaskTemplateSetSearch) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	var setList []taskMdl.TaskTemplateSet
	db := global.GVA_DB.Model(&taskMdl.TaskTemplateSet{})
	if info.Name != "" {
		name := strings.Trim(info.Name, " ")
		db = db.Where("`name` LIKE ?", "%"+name+"%")
	}
	if len(info.SystemIDs) > 0 {
		db = db.Where("`system_id` IN ?", info.SystemIDs)
	}
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	//err = db.Limit(limit).Offset(offset).Find(&setList).Error
	db = db.Limit(limit).Offset(offset)
	if info.OrderKey != "" {
		var OrderStr string
		if info.Desc {
			OrderStr = info.OrderKey + " desc"
		} else {
			OrderStr = info.OrderKey
		}
		err = db.Order(OrderStr).Find(&setList).Error
	} else {
		err = db.Order("id").Find(&setList).Error
	}

	setInfoList := make([]response.TaskTemplateSetResponse, 0, len(setList))
	for _, set := range setList {
		setTemplates := make([]taskMdl.TaskTemplateSetTemplate, 0)
		if err = global.GVA_DB.Where("set_id = ?", set.ID).Find(&setTemplates).Error; err != nil {
			return
		}
		ress := make([]response.TaskTemplateSetTemplateResponse, 0)
		for _, t := range setTemplates {
			var template taskMdl.TaskTemplate
			if err = global.GVA_DB.Where("id = ?", t.TemplateId).Find(&template).Error; err != nil {
				return
			}
			res := response.TaskTemplateSetTemplateResponse{
				TaskTemplateSetTemplate: t,
				TemplateName:            template.Name,
			}
			ress = append(ress, res)
		}
		setInfoList = append(setInfoList, response.TaskTemplateSetResponse{
			TaskTemplateSet: set,
			Templates:       ress,
		})
	}
	return err, setInfoList, total
}

//@author: [telexy324](https://github.com/telexy324)
//@function: AddSetTask
//@description: 添加模板集
//@param: setTask taskMdl.TaskTemplateSetTask
//@return: error

func (templateService *TaskTemplatesService) AddSetTask(addSetTaskRequest taskMdl.SetTask) (err error) {
	templates := make([]taskMdl.TaskTemplateSetTemplate, 0)
	if err = global.GVA_DB.Where("set_id = ?", addSetTaskRequest.SetId).Order("seq, id").Find(&templates).Error; err != nil {
		return
	}
	stepsMap := make(map[int]bool)
	for _, t := range templates {
		stepsMap[t.SetId] = true
	}
	templatesBytes, err := json.Marshal(templates)
	if err != nil {
		return
	}
	addSetTaskRequest.TemplatesString = string(templatesBytes)
	addSetTaskRequest.TotalSteps = len(stepsMap)
	return global.GVA_DB.Create(&addSetTaskRequest).Error
}

//@author: [telexy324](https://github.com/telexy324)
//@function: UpdateSetTask
//@description: 更新模板集任务集
//@param: system taskMdl.SetTask
//@return: error

func (templateService *TaskTemplatesService) UpdateSetTask(setTask taskMdl.SetTask) (err error) {
	var oldSetTask taskMdl.SetTask
	upDateMap := make(map[string]interface{})
	upDateMap["current_task_ids_string"] = setTask.CurrentTaskIdsString
	upDateMap["tasks_string"] = setTask.TasksString
	upDateMap["current_step"] = setTask.CurrentStep
	upDateMap["force_correct"] = setTask.ForceCorrect

	db := global.GVA_DB.Where("id = ?", setTask.ID).Find(&oldSetTask)
	err = db.Updates(upDateMap).Error
	if err != nil {
		global.GVA_LOG.Error(err.Error())
		return err
	}
	return
}

//@author: [telexy324](https://github.com/telexy324)
//@function: GetSetTaskById
//@description: 返回当前选中system
//@param: id float64
//@return: err error, set taskMdl.SetTask

func (templateService *TaskTemplatesService) GetSetTaskById(id float64) (err error, setTask taskMdl.SetTask) {
	if err = global.GVA_DB.Where("id = ?", id).First(&setTask).Error; err != nil {
		return
	}
	return
}

//@author: [telexy324](https://github.com/telexy324)
//@function: GetSystemList
//@description: 获取系统分页
//@return: err error, list interface{}, total int64

func (templateService *TaskTemplatesService) GetSetTaskList(info request2.SetTaskSearch) (err error, list interface{}, total int64) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)
	var setTaskList []taskMdl.SetTask
	db := global.GVA_DB.Model(&taskMdl.SetTask{})
	db = db.Where("set_id = ?", info.SetId)
	err = db.Count(&total).Error
	if err != nil {
		return
	}
	//err = db.Limit(limit).Offset(offset).Find(&setTaskList).Error
	db = db.Limit(limit).Offset(offset)
	if info.OrderKey != "" {
		var OrderStr string
		if info.Desc {
			OrderStr = info.OrderKey + " desc"
		} else {
			OrderStr = info.OrderKey
		}
		err = db.Order(OrderStr).Find(&setTaskList).Error
	} else {
		err = db.Order("id").Find(&setTaskList).Error
	}

	return err, setTaskList, total
}

func (templateService *TaskTemplatesService) SetTaskForceCorrect(id float64) (err error) {
	var setTask taskMdl.SetTask
	if err = global.GVA_DB.Where("id = ?", id).First(&setTask).Error; err != nil {
		return
	}
	if setTask.ForceCorrect == consts.IsForceCorrect {
		return errors.New("set task is force correct do not correct it again")
	}
	setTask.ForceCorrect = consts.IsForceCorrect
	return templateService.UpdateSetTask(setTask)
}

func (templateService *TaskTemplatesService) GetFileList(sshClient *common.SSHClient, template taskMdl.TaskTemplate, selectedDirectory string) (fileInfos []response.FileInfo, isTop bool, err error) {
	selectedDirectory = strings.TrimRight(selectedDirectory, "/") + "/"
	if !strings.Contains(selectedDirectory, template.LogPath) {
		err = errors.New("directory not in log path")
		return
	}
	commandFile := `ls -lh ` + selectedDirectory + ` | grep ^- | awk '{print $9 " " $5}'`
	outputFile, err := sshClient.CommandSingle(commandFile)
	if err != nil {
		return
	}
	commandDirectory := `ls -lh ` + selectedDirectory + ` | grep ^d | awk '{print $9}'`
	outputDirectory, err := sshClient.CommandSingle(commandDirectory)
	if err != nil {
		return
	}
	if outputDirectory != "" {
		outputDirectory = strings.TrimRight(outputDirectory, "\n")
		directories := strings.Split(outputDirectory, "\n")
		for _, d := range directories {
			fileInfos = append(fileInfos, response.FileInfo{
				FileName:  d,
				Directory: true,
			})
		}
	}
	if outputFile != "" {
		outputFile = strings.TrimRight(outputFile, "\n")
		fileNames := strings.Split(outputFile, "\n")
		for _, f := range fileNames {
			fileInfos = append(fileInfos, response.FileInfo{
				FileName:  f,
				Directory: false,
			})
		}
	}
	if len(selectedDirectory) <= len(strings.TrimRight(template.LogPath, "/")+"/") {
		isTop = true
	}
	//files := strings.Split(outputs, "\n")
	//for _, f := range files {
	//	fields := strings.Split(f, " ")
	//	if len(fields) < 2 {
	//		continue
	//	}
	//	fileNames = append(fileNames, fields[0]+strings.Repeat(" ", 12-len(fields[0]))+fields[1])
	//}
	return
}
