package taskPool

import (
	"bytes"
	"crypto/md5"
	"database/sql"
	"encoding/json"
	"fmt"
	sockets "github.com/flipped-aurora/gin-vue-admin/server/api/v1/socket"
	"github.com/flipped-aurora/gin-vue-admin/server/common"
	"github.com/flipped-aurora/gin-vue-admin/server/consts"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/application"
	appReq "github.com/flipped-aurora/gin-vue-admin/server/model/application/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/model/taskMdl"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/jlaffaye/ftp"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
	"gorm.io/gorm"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type TaskRunner struct {
	task     taskMdl.Task
	template taskMdl.TaskTemplate

	users []int
	//alert     bool
	//alertChat *string
	prepared bool
	clients  []*ssh.Client
	pool     *TaskPool
	ftpConn  []*ftp.ServerConn
}

func getMD5Hash(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func (t *TaskRunner) setStatus(status taskMdl.TaskStatus) {
	if t.task.Status == taskMdl.TaskStoppingStatus {
		switch status {
		case taskMdl.TaskFailStatus:
			status = taskMdl.TaskStoppedStatus
		case taskMdl.TaskStoppedStatus:
		default:
			panic("stopping TaskRunner cannot be " + status)
		}
	}

	t.task.Status = status

	t.updateStatus()

	// t.sendSlackAlert()

	//if status == taskMdl.TaskFailStatus {
	//	t.sendMailAlert()
	//}
	//
	//if status == taskMdl.TaskSuccessStatus || status == taskMdl.TaskFailStatus {
	//	t.sendTelegramAlert()
	//}
}

func (t *TaskRunner) updateStatus() {
	for _, user := range t.users {
		b, err := json.Marshal(&map[string]interface{}{
			"type":       "update",
			"beginTime":  t.task.BeginTime,
			"endTime":    t.task.EndTime,
			"status":     t.task.Status,
			"ID":         t.task.ID,
			"templateId": t.task.TemplateId,
			"taskId":     t.task.ID,
		})

		if err != nil {
			global.GVA_LOG.Fatal(err.Error())
		}

		sockets.Message(user, b)
	}

	if err := taskService.UpdateTask(t.task); err != nil {
		t.panicOnError(err, "Failed to update TaskRunner status")
	}
}

func (t *TaskRunner) fail() {
	t.setStatus(taskMdl.TaskFailStatus)
}

//func (t *TaskRunner) destroyKeys() {
//	err := t.inventory.SSHKey.Destroy()
//	if err != nil {
//		t.Log("Can't destroy inventory user key, error: " + err.Error())
//	}
//
//	err = t.inventory.BecomeKey.Destroy()
//	if err != nil {
//		t.Log("Can't destroy inventory become user key, error: " + err.Error())
//	}
//
//	err = t.template.VaultKey.Destroy()
//	if err != nil {
//		t.Log("Can't destroy inventory vault password file, error: " + err.Error())
//	}
//}
//
//func (t *TaskRunner) createTaskEvent() {
//	objType := taskMdl.EventTask
//	desc := "Task ID " + strconv.Itoa(t.task.ID) + " (" + t.template.Name + ")" + " finished - " + strings.ToUpper(string(t.task.Status))
//
//	_, err := t.pool.store.CreateEvent(taskMdl.Event{
//		UserID:      t.task.UserID,
//		ProjectID:   &t.task.ProjectID,
//		ObjectType:  &objType,
//		ObjectID:    &t.task.ID,
//		Description: &desc,
//	})
//
//	if err != nil {
//		t.panicOnError(err, "Fatal error inserting an event")
//	}
//}

func (t *TaskRunner) prepareRun() {
	t.prepared = false

	defer func() {
		global.GVA_LOG.Info("Stopped preparing TaskRunner ", zap.Uint("task ID ", t.task.ID))
		global.GVA_LOG.Info("Release resource locker with TaskRunner ", zap.Uint("task ID ", t.task.ID))
		t.pool.resourceLocker <- &resourceLock{lock: false, holder: t}

		//t.createTaskEvent()
	}()

	t.Log("Preparing: " + strconv.Itoa(int(t.task.ID)))

	if err := checkTmpDir(global.GVA_CONFIG.Task.TmpPath); err != nil {
		t.Log("Creating tmp dir failed: " + err.Error())
		t.fail()
		return
	}

	//objType := taskMdl.EventTask
	//desc := "Task ID " + strconv.Itoa(t.task.ID) + " (" + t.template.Name + ")" + " is preparing"
	//evt := taskMdl.Event{
	//	UserID:      t.task.UserID,
	//	ProjectID:   &t.task.ProjectID,
	//	ObjectType:  &objType,
	//	ObjectID:    &t.task.ID,
	//	Description: &desc,
	//}
	//
	//if _, err := t.pool.store.CreateEvent(evt); err != nil {
	//	t.Log("Fatal error inserting an event")
	//	panic(err)
	//}

	t.Log("Prepare TaskRunner with template: " + t.template.Name + "\n")

	t.updateStatus()

	//if t.repository.GetType() == taskMdl.RepositoryLocal {
	//	if _, err := os.Stat(t.repository.GitURL); err != nil {
	//		t.Log("Failed in finding static repository at " + t.repository.GitURL + ": " + err.Error())
	//		t.fail()
	//		return
	//	}
	//} else {
	//	if err := t.updateRepository(); err != nil {
	//		t.Log("Failed updating repository: " + err.Error())
	//		t.fail()
	//		return
	//	}
	//	if err := t.checkoutRepository(); err != nil {
	//		t.Log("Failed to checkout repository to required commit: " + err.Error())
	//		t.fail()
	//		return
	//	}
	//}
	//
	//if err := t.installInventory(); err != nil {
	//	t.Log("Failed to install inventory: " + err.Error())
	//	t.fail()
	//	return
	//}
	//
	//if err := t.installRequirements(); err != nil {
	//	t.Log("Running galaxy failed: " + err.Error())
	//	t.fail()
	//	return
	//}
	//
	//if err := t.installVaultKeyFile(); err != nil {
	//	t.Log("Failed to install vault password file: " + err.Error())
	//	t.fail()
	//	return
	//}

	t.prepared = true
}

func (t *TaskRunner) run() {
	defer func() {
		global.GVA_LOG.Info("Stopped running TaskRunner ", zap.Uint("task ID ", t.task.ID))
		global.GVA_LOG.Info("Release resource locker with TaskRunner ", zap.Uint("task ID ", t.task.ID))
		t.pool.resourceLocker <- &resourceLock{lock: false, holder: t}

		t.task.EndTime = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
		t.updateStatus()
		//t.createTaskEvent()
		//t.destroyKeys()
	}()

	// TODO: more details
	if t.task.Status == taskMdl.TaskStoppingStatus {
		t.setStatus(taskMdl.TaskStoppedStatus)
		return
	}

	t.task.BeginTime = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	t.setStatus(taskMdl.TaskRunningStatus)

	//objType := taskMdl.EventTask
	//desc := "Task ID " + strconv.Itoa(t.task.ID) + " (" + t.template.Name + ")" + " is running"
	//
	//_, err := t.pool.store.CreateEvent(taskMdl.Event{
	//	UserID:      t.task.UserID,
	//	ProjectID:   &t.task.ProjectID,
	//	ObjectType:  &objType,
	//	ObjectID:    &t.task.ID,
	//	Description: &desc,
	//})

	//if err != nil {
	//	t.Log("Fatal error inserting an event")
	//	panic(err)
	//}

	t.Log("Started: " + strconv.Itoa(int(t.task.ID)))
	t.Log("Run TaskRunner with template: " + t.template.Name + "\n")

	// TODO: ?????
	if t.task.Status == taskMdl.TaskStoppingStatus {
		t.setStatus(taskMdl.TaskStoppedStatus)
		return
	}

	failedIPs := make([]string, 0)
	if t.task.TemplateId == consts.TaskTemplateDiscoverServers {
		failedIPs = t.runDiscoverTask()
	} else {
		if t.template.ExecuteType == consts.ExecuteTypeDownload {
			failedIPs = t.runUploadTask()
		} else if t.template.ExecuteType == consts.ExecuteTypeDeploy {
			failedIPs = t.runDeployTask()
		} else {
			failedIPs = t.runTask()
		}
	}
	if len(failedIPs) > 0 {
		failed, _ := json.Marshal(failedIPs)
		t.Log("Running task failed: " + string(failed))
		t.fail()
		return
	}

	t.setStatus(taskMdl.TaskSuccessStatus)

	//templates, err := t.pool.store.GetTemplates(t.task.ProjectID, taskMdl.TemplateFilter{
	//	BuildTemplateID: &t.task.TemplateID,
	//	AutorunOnly:     true,
	//}, taskMdl.RetrieveQueryParams{})
	//if err != nil {
	//	t.Log("Running playbook failed: " + err.Error())
	//	return
	//}
	//
	//for _, tpl := range templates {
	//	_, err = t.pool.AddTask(taskMdl.Task{
	//		TemplateID:  tpl.ID,
	//		ProjectID:   tpl.ProjectID,
	//		BuildTaskID: &t.task.ID,
	//	}, nil, tpl.ProjectID)
	//	if err != nil {
	//		t.Log("Running playbook failed: " + err.Error())
	//		continue
	//	}
	//}
}

func (t *TaskRunner) prepareError(err error, errMsg string) error {
	if err == gorm.ErrRecordNotFound {
		t.Log(errMsg)
		return err
	}

	if err != nil {
		t.fail()
		panic(err)
	}

	return nil
}

// nolint: gocyclo
func (t *TaskRunner) populateDetails() error {
	// get template
	var err error

	if t.task.TemplateId < 99999900 {
		t.template, err = taskService.GetTaskTemplate(float64(t.task.TemplateId))
		if err != nil {
			return t.prepareError(err, "Template not found!")
		}
	}

	// get project alert setting
	//project, err := t.pool.store.GetProject(t.template.ProjectID)
	//if err != nil {
	//	return t.prepareError(err, "Project not found!")
	//}
	//
	//t.alert = project.Alert
	//t.alertChat = project.AlertChat

	// get project users
	err, userList, _ := userService.GetUserInfoList(request.PageInfo{
		Page:     1,
		PageSize: 1000,
	})
	if err != nil {
		return t.prepareError(err, "Users not found!")
	}

	t.users = []int{}
	for _, user := range userList.([]system.SysUser) {
		t.users = append(t.users, int(user.ID))
	}

	//// get inventory
	//t.inventory, err = t.pool.store.GetInventory(t.template.ProjectID, t.template.InventoryID)
	//if err != nil {
	//	return t.prepareError(err, "Template Inventory not found!")
	//}
	//
	//// get repository
	//t.repository, err = t.pool.store.GetRepository(t.template.ProjectID, t.template.RepositoryID)
	//
	//if err != nil {
	//	return err
	//}
	//
	//err = t.repository.SSHKey.DeserializeSecret()
	//if err != nil {
	//	return err
	//}
	//
	//// get environment
	//if t.template.EnvironmentID != nil {
	//	t.environment, err = t.pool.store.GetEnvironment(t.template.ProjectID, *t.template.EnvironmentID)
	//	if err != nil {
	//		return err
	//	}
	//}
	//
	//if t.task.Environment != "" {
	//	environment := make(map[string]interface{})
	//	if t.environment.JSON != "" {
	//		err = json.Unmarshal([]byte(t.task.Environment), &environment)
	//		if err != nil {
	//			return err
	//		}
	//	}
	//
	//	taskEnvironment := make(map[string]interface{})
	//	err = json.Unmarshal([]byte(t.environment.JSON), &taskEnvironment)
	//	if err != nil {
	//		return err
	//	}
	//
	//	for k, v := range taskEnvironment {
	//		environment[k] = v
	//	}
	//
	//	var ev []byte
	//	ev, err = json.Marshal(environment)
	//	if err != nil {
	//		return err
	//	}
	//
	//	t.environment.JSON = string(ev)
	//}

	return nil
}

//func (t *TaskRunner) installVaultKeyFile() error {
//	if t.template.VaultKeyID == nil {
//		return nil
//	}
//
//	return t.template.VaultKey.Install(taskMdl.AccessKeyRoleAnsiblePasswordVault)
//}
//
//func (t *TaskRunner) checkoutRepository() error {
//
//	repo := lib.GitRepository{
//		Logger:     t,
//		TemplateID: t.template.ID,
//		Repository: t.repository,
//	}
//
//	err := repo.ValidateRepo()
//
//	if err != nil {
//		return err
//	}
//
//	if t.task.CommitHash != nil {
//		// checkout to commit if it is provided for TaskRunner
//		return repo.Checkout(*t.task.CommitHash)
//	}
//
//	// store commit to TaskRunner table
//
//	commitHash, err := repo.GetLastCommitHash()
//
//	if err != nil {
//		return err
//	}
//
//	commitMessage, _ := repo.GetLastCommitMessage()
//
//	t.task.CommitHash = &commitHash
//	t.task.CommitMessage = commitMessage
//
//	return t.pool.store.UpdateTask(t.task)
//}
//
//func (t *TaskRunner) updateRepository() error {
//	repo := lib.GitRepository{
//		Logger:     t,
//		TemplateID: t.template.ID,
//		Repository: t.repository,
//	}
//
//	err := repo.ValidateRepo()
//
//	if err != nil {
//		if !os.IsNotExist(err) {
//			err = os.RemoveAll(repo.GetFullPath())
//			if err != nil {
//				return err
//			}
//		}
//		return repo.Clone()
//	}
//
//	if repo.CanBePulled() {
//		err = repo.Pull()
//		if err == nil {
//			return nil
//		}
//	}
//
//	err = os.RemoveAll(repo.GetFullPath())
//	if err != nil {
//		return err
//	}
//
//	return repo.Clone()
//}

//func (t *TaskRunner) installCollectionsRequirements() error {
//	requirementsFilePath := fmt.Sprintf("%s/collections/requirements.yml", t.getRepoPath())
//	requirementsHashFilePath := fmt.Sprintf("%s.md5", requirementsFilePath)
//
//	if _, err := os.Stat(requirementsFilePath); err != nil {
//		t.Log("No collections/requirements.yml file found. Skip galaxy install process.\n")
//		return nil
//	}
//
//	if hasRequirementsChanges(requirementsFilePath, requirementsHashFilePath) {
//		if err := t.runGalaxy([]string{
//			"collection",
//			"install",
//			"-r",
//			requirementsFilePath,
//			"--force",
//		}); err != nil {
//			return err
//		}
//		if err := writeMD5Hash(requirementsFilePath, requirementsHashFilePath); err != nil {
//			return err
//		}
//	} else {
//		t.Log("collections/requirements.yml has no changes. Skip galaxy install process.\n")
//	}
//
//	return nil
//}
//
//func (t *TaskRunner) installRolesRequirements() error {
//	requirementsFilePath := fmt.Sprintf("%s/roles/requirements.yml", t.getRepoPath())
//	requirementsHashFilePath := fmt.Sprintf("%s.md5", requirementsFilePath)
//
//	if _, err := os.Stat(requirementsFilePath); err != nil {
//		t.Log("No roles/requirements.yml file found. Skip galaxy install process.\n")
//		return nil
//	}
//
//	if hasRequirementsChanges(requirementsFilePath, requirementsHashFilePath) {
//		if err := t.runGalaxy([]string{
//			"role",
//			"install",
//			"-r",
//			requirementsFilePath,
//			"--force",
//		}); err != nil {
//			return err
//		}
//		if err := writeMD5Hash(requirementsFilePath, requirementsHashFilePath); err != nil {
//			return err
//		}
//	} else {
//		t.Log("roles/requirements.yml has no changes. Skip galaxy install process.\n")
//	}
//
//	return nil
//}
//
//func (t *TaskRunner) installRequirements() error {
//	if err := t.installCollectionsRequirements(); err != nil {
//		return err
//	}
//	if err := t.installRolesRequirements(); err != nil {
//		return err
//	}
//	return nil
//}
//
//func (t *TaskRunner) runGalaxy(args []string) error {
//	return lib.AnsiblePlaybook{
//		Logger:     t,
//		TemplateID: t.template.ID,
//		Repository: t.repository,
//	}.RunGalaxy(args)
//}

func (t *TaskRunner) runTask() (failedIPs []string) {
	servers := t.template.TargetServers
	global.GVA_LOG.Info("run ", zap.Uint("task ID: ", t.task.ID))
	wg := &sync.WaitGroup{}
	failedChan := make(chan string, len(servers))
	t.clients = make([]*ssh.Client, 0, len(servers))
	for _, server := range servers {
		wg.Add(1)
		go func(w *sync.WaitGroup, s application.ApplicationServer, f chan string) {
			defer w.Done()
			sshClient, err := common.FillSSHClient(s.ManageIp, t.template.SysUser, "", s.SshPort)
			if err != nil {
				global.GVA_LOG.Error("run task failed on create ssh client: ", zap.Uint("task ID: ", t.task.ID), zap.String("server IP: ", s.ManageIp), zap.Any("err", err))
				f <- s.ManageIp
				return
			}
			err = sshClient.GenerateClient()
			if err != nil {
				global.GVA_LOG.Error("run task failed on generate ssh client: ", zap.Uint("task ID: ", t.task.ID), zap.String("server IP: ", s.ManageIp), zap.Any("err", err))
				f <- s.ManageIp
				return
			}
			defer sshClient.Client.Close()
			t.clients = append(t.clients, sshClient.Client)
			//ssConn, err := sshService.NewSshConn(sshClient.Client, 0, 0)
			//if err != nil {
			//	global.GVA_LOG.Error("run task failed on create ssh session: ", zap.Uint("task ID: ", t.task.ID), zap.String("server IP: ", s.ManageIp), zap.Any("err", err))
			//	f <- s.ManageIp
			//	return
			//}
			//defer ssConn.Close()

			//var command string
			//if t.template.Mode == consts.Command {
			//	command = t.template.Command
			//} else {
			//	command = "sh "+t.template.ScriptPath
			//}
			//
			//err = sshClient.Command(command, t)
			//if err != nil {
			//	global.GVA_LOG.Error("run task failed on exec command: ", zap.Uint("task ID: ", t.task.ID), zap.String("server IP: ", s.ManageIp), zap.Any("err", err))
			//	f <- s.ManageIp
			//	return
			//}
			if t.template.Mode == consts.Command {
				failed := false
				if t.template.Interactive == consts.Interactive {
					commands := strings.Split(t.template.Command, "\n")
					err = sshClient.CommandBatch(commands, t, s.ManageIp)
					if err != nil {
						global.GVA_LOG.Error("run task failed on exec command: ", zap.Uint("task ID: ", t.task.ID), zap.String("server IP: ", s.ManageIp), zap.Any("err", err))
					} else {
						failed = true
					}
				} else {
					commands := strings.Replace(t.template.Command, "\n", " && ", -1)
					commands = "source /etc/profile && source ~/.bashrc && " + commands
					err = sshClient.Commands(commands, t, s.ManageIp)
					if err != nil {
						global.GVA_LOG.Error("run task failed on exec command: ", zap.Uint("task ID: ", t.task.ID), zap.String("server IP: ", s.ManageIp), zap.Any("err", err))
					} else {
						failed = true
					}
				}
				//commands := strings.Split(t.template.Command, "\n")
				//failed := false
				//for _, command := range commands {
				//	err = sshClient.Commands(command, t, s.ManageIp)
				//	if err != nil {
				//		global.GVA_LOG.Error("run task failed on exec command: ", zap.Uint("task ID: ", t.task.ID), zap.String("server IP: ", s.ManageIp), zap.Any("err", err))
				//	} else {
				//		failed = true
				//	}
				//	time.Sleep(time.Microsecond * 500)
				//}
				//err = sshClient.CommandBatch(commands, t, s.ManageIp)
				//if err != nil {
				//	global.GVA_LOG.Error("run task failed on exec command: ", zap.Uint("task ID: ", t.task.ID), zap.String("server IP: ", s.ManageIp), zap.Any("err", err))
				//} else {
				//	failed = true
				//}
				//
				//failed := false
				//commands := strings.Replace(t.template.Command, "\n", " && ", -1)
				//err = sshClient.Commands(commands, t, s.ManageIp)
				//if err != nil {
				//	global.GVA_LOG.Error("run task failed on exec command: ", zap.Uint("task ID: ", t.task.ID), zap.String("server IP: ", s.ManageIp), zap.Any("err", err))
				//} else {
				//	failed = true
				//}
				if !failed {
					f <- s.ManageIp
				}
			} else {
				//command := "sh " + t.template.ScriptPath
				var shellType = "sh "
				if t.template.ShellType == consts.ShellTypeSh {
					shellType = "sh "
				} else if t.template.ShellType == consts.ShellTypeBash {
					shellType = "bash "
				}
				command := shellType + strings.Trim(t.template.ScriptPath, " ")
				if len(strings.Trim(t.template.ShellVars, " ")) > 0 {
					command = command + " " + strings.Trim(t.template.ShellVars, " ")
				}
				err = sshClient.Commands(command, t, s.ManageIp)
				if err != nil {
					global.GVA_LOG.Error("run task failed on exec command: ", zap.Uint("task ID: ", t.task.ID), zap.String("server IP: ", s.ManageIp), zap.Any("err", err))
					f <- s.ManageIp
					return
				}
			}
			//sshClient.RequestShell()
			//if err = sshClient.ConnectShell(t.template.Command, t); err != nil {
			//	global.GVA_LOG.Error("run task failed on run command: ", zap.Uint("task ID: ", t.task.ID), zap.String("server IP: ", s.ManageIp), zap.Any("err", err))
			//	f <- s.ManageIp
			//	return
			//}
			return
		}(wg, server, failedChan)
	}
	wg.Wait()
	close(failedChan)
	for ip := range failedChan {
		failedIPs = append(failedIPs, ip)
	}
	return
}

func (t *TaskRunner) runUploadTask() (failedIPs []string) {
	if t.template.TargetServers == nil || len(t.template.TargetServers) <= 0 {
		global.GVA_LOG.Error("run task failed on nil target server: ", zap.Uint("task ID: ", t.task.ID))
		return
	}
	global.GVA_LOG.Info("run ", zap.Uint("task ID: ", t.task.ID))
	server := t.template.TargetServers[0]
	sshClient, err := common.FillSSHClient(server.ManageIp, t.template.SysUser, "", server.SshPort)
	if err != nil {
		global.GVA_LOG.Error("run task failed on create ssh client: ", zap.Uint("task ID: ", t.task.ID), zap.String("server IP: ", server.ManageIp), zap.Any("err", err))
		failedIPs = append(failedIPs, server.ManageIp)
		return
	}
	err = sshClient.GenerateClient()
	if err != nil {
		global.GVA_LOG.Error("run task failed on generate ssh client: ", zap.Uint("task ID: ", t.task.ID), zap.String("server IP: ", server.ManageIp), zap.Any("err", err))
		failedIPs = append(failedIPs, server.ManageIp)
		return
	}
	defer sshClient.Client.Close()
	t.clients = append(t.clients, sshClient.Client)
	//filePath := "/" + strings.Trim(t.template.LogPath, "/") + "/"
	fileReader, err := sshClient.DownloadReader(t.task.FileDownload)
	if err != nil {
		global.GVA_LOG.Error("run task failed on download file failed", zap.Any("err", err))
		failedIPs = append(failedIPs, server.ManageIp)
		return
	}
	t.Log("download "+t.task.FileDownload+" success", server.ManageIp)
	filePathUpload := "/" + strings.Trim(t.template.LogDst, "/") + "/"
	paths := strings.Split(t.task.FileDownload, "/")
	fileUpload := filePathUpload + paths[len(paths)-1]
	if t.template.LogUploadServer.Mode == consts.LogServerModeFtp {
		ftpClient, err := common.NewFtpClient(t.template.LogUploadServer.ManageIp, t.template.LogUploadServer.Port, t.template.Secret.Name, t.template.Secret.Password)
		if err != nil {
			global.GVA_LOG.Error("create ftp client failed", zap.Any("err", err))
			failedIPs = append(failedIPs, server.ManageIp)
			return
		}
		defer ftpClient.Conn.Quit()
		t.ftpConn = append(t.ftpConn, ftpClient.Conn)
		if err = ftpClient.Upload(fileUpload, fileReader); err != nil {
			global.GVA_LOG.Error("upload via ftp failed", zap.Any("err", err))
			failedIPs = append(failedIPs, server.ManageIp)
			return
		}
	} else if t.template.LogUploadServer.Mode == consts.LogServerModeSSH {
		sshClientUpload, err := common.FillSSHClient(t.template.LogUploadServer.ManageIp, t.template.Secret.Name, t.template.Secret.Password, t.template.LogUploadServer.Port)
		err = sshClientUpload.GenerateClient()
		if err != nil {
			global.GVA_LOG.Error("create ssh client failed: ", zap.String("server IP: ", server.ManageIp), zap.Any("err", err))
			failedIPs = append(failedIPs, server.ManageIp)
			return
		}
		defer sshClientUpload.Client.Close()
		t.clients = append(t.clients, sshClientUpload.Client)
		if err = sshClientUpload.Upload(fileReader, fileUpload); err != nil {
			global.GVA_LOG.Error("upload via sftp failed", zap.Any("err", err))
			failedIPs = append(failedIPs, server.ManageIp)
			return
		}
	}
	t.Log("upload "+fileUpload+" success", t.template.LogUploadServer.ManageIp)
	return
}

//func (t *TaskRunner) runDeployTask() (failedIPs []string) {
//	var fb []byte
//	manageIps := make([]string, 0, len(t.template.TargetServers))
//	var manageIpString string
//	for _, s := range t.template.TargetServers {
//		manageIps = append(failedIPs, s.ManageIp)
//		manageIpString = manageIpString + " " + s.ManageIp
//	}
//	if t.template.LogUploadServer.Mode == consts.LogServerModeFtp {
//		ftpClient, err := common.NewFtpClient(t.template.LogUploadServer.ManageIp, t.template.LogUploadServer.Port, t.template.Secret.Name, t.template.Secret.Password)
//		if err != nil {
//			global.GVA_LOG.Error("create ftp client failed", zap.Any("err", err))
//			failedIPs = manageIps
//			return
//		}
//		defer ftpClient.Conn.Quit()
//		t.ftpConn = append(t.ftpConn, ftpClient.Conn)
//		if fb, err = ftpClient.Download(t.template.DownloadSource); err != nil {
//			global.GVA_LOG.Error("upload via ftp failed", zap.Any("err", err))
//			failedIPs = manageIps
//			return
//		}
//	} else if t.template.LogUploadServer.Mode == consts.LogServerModeSSH {
//		sshClientUpload, err := common.FillSSHClient(t.template.LogUploadServer.ManageIp, t.template.Secret.Name, t.template.Secret.Password, t.template.LogUploadServer.Port)
//		err = sshClientUpload.GenerateClient()
//		if err != nil {
//			global.GVA_LOG.Error("create ssh client failed: ", zap.Any("err", err))
//			failedIPs = manageIps
//			return
//		}
//		defer sshClientUpload.Client.Close()
//		t.clients = append(t.clients, sshClientUpload.Client)
//		if fb, err = sshClientUpload.Download(t.template.DownloadSource); err != nil {
//			global.GVA_LOG.Error("upload via sftp failed", zap.Any("err", err))
//			failedIPs = manageIps
//			return
//		}
//	}
//	t.Log("download "+t.template.DownloadSource+" success", t.template.LogUploadServer.ManageIp)
//	if t.template.TargetServers == nil || len(t.template.TargetServers) <= 0 {
//		global.GVA_LOG.Error("run task failed on nil target server: ", zap.Uint("task ID: ", t.task.ID))
//		return
//	}
//	global.GVA_LOG.Info("run ", zap.Uint("task ID: ", t.task.ID))
//	servers := t.template.TargetServers
//	wg := &sync.WaitGroup{}
//	failedChan := make(chan string, len(servers))
//	t.clients = make([]*ssh.Client, 0, len(servers))
//	for _, server := range servers {
//		wg.Add(1)
//		go func(w *sync.WaitGroup, s application.ApplicationServer, f chan string) {
//			defer w.Done()
//			sshClient, err := common.FillSSHClient(s.ManageIp, t.template.SysUser, "", s.SshPort)
//			if err != nil {
//				global.GVA_LOG.Error("run task failed on create ssh client: ", zap.Uint("task ID: ", t.task.ID), zap.String("server IP: ", s.ManageIp), zap.Any("err", err))
//				f <- s.ManageIp
//				return
//			}
//			err = sshClient.GenerateClient()
//			if err != nil {
//				global.GVA_LOG.Error("run task failed on generate ssh client: ", zap.Uint("task ID: ", t.task.ID), zap.String("server IP: ", s.ManageIp), zap.Any("err", err))
//				f <- s.ManageIp
//				return
//			}
//			defer sshClient.Client.Close()
//			t.clients = append(t.clients, sshClient.Client)
//			fio := bytes.NewReader(fb)
//			if err = sshClient.Upload(fio, t.template.DeployPath); err != nil {
//				global.GVA_LOG.Error("run task failed on upload deploy file: ", zap.Uint("task ID: ", t.task.ID), zap.String("server IP: ", s.ManageIp), zap.Any("err", err))
//				f <- s.ManageIp
//				return
//			}
//			return
//		}(wg, server, failedChan)
//	}
//	wg.Wait()
//	close(failedChan)
//	for ip := range failedChan {
//		failedIPs = append(failedIPs, ip)
//	}
//	t.Log("deploy "+t.template.DeployPath+" success", manageIpString)
//	return
//}

func (t *TaskRunner) runDeployTask() (failedIPs []string) {
	var fb []byte
	manageIps := make([]string, 0, len(t.template.TargetServers))
	var manageIpString string
	for _, s := range t.template.TargetServers {
		manageIps = append(failedIPs, s.ManageIp)
		manageIpString = manageIpString + " " + s.ManageIp
	}
	if t.template.LogUploadServer.Mode == consts.LogServerModeFtp {
		ftpClient, err := common.NewFtpClient(t.template.LogUploadServer.ManageIp, t.template.LogUploadServer.Port, t.template.Secret.Name, t.template.Secret.Password)
		if err != nil {
			global.GVA_LOG.Error("create ftp client failed", zap.Any("err", err))
			failedIPs = manageIps
			return
		}
		defer ftpClient.Conn.Quit()
		t.ftpConn = append(t.ftpConn, ftpClient.Conn)
		for _, deployInfo := range t.template.TaskDeployInfos {
			if fb, err = ftpClient.Download(deployInfo.DownloadSource); err != nil {
				global.GVA_LOG.Error("upload via ftp failed", zap.Any("err", err))
				failedIPs = manageIps
				return
			}
			t.Log("download "+deployInfo.DownloadSource+" success", t.template.LogUploadServer.ManageIp)
			if t.template.TargetServers == nil || len(t.template.TargetServers) <= 0 {
				global.GVA_LOG.Error("run task failed on nil target server: ", zap.Uint("task ID: ", t.task.ID))
				return
			}
			global.GVA_LOG.Info("run ", zap.Uint("task ID: ", t.task.ID))
			servers := t.template.TargetServers
			wg := &sync.WaitGroup{}
			failedChan := make(chan string, len(servers))
			t.clients = make([]*ssh.Client, 0, len(servers))
			for _, server := range servers {
				wg.Add(1)
				go func(w *sync.WaitGroup, s application.ApplicationServer, f chan string) {
					defer w.Done()
					sshClient, err := common.FillSSHClient(s.ManageIp, t.template.SysUser, "", s.SshPort)
					if err != nil {
						global.GVA_LOG.Error("run task failed on create ssh client: ", zap.Uint("task ID: ", t.task.ID), zap.String("server IP: ", s.ManageIp), zap.Any("err", err))
						f <- s.ManageIp
						return
					}
					err = sshClient.GenerateClient()
					if err != nil {
						global.GVA_LOG.Error("run task failed on generate ssh client: ", zap.Uint("task ID: ", t.task.ID), zap.String("server IP: ", s.ManageIp), zap.Any("err", err))
						f <- s.ManageIp
						return
					}
					defer sshClient.Client.Close()
					t.clients = append(t.clients, sshClient.Client)
					fio := bytes.NewReader(fb)
					if err = sshClient.Upload(fio, deployInfo.DeployPath); err != nil {
						global.GVA_LOG.Error("run task failed on upload deploy file: ", zap.Uint("task ID: ", t.task.ID), zap.String("server IP: ", s.ManageIp), zap.Any("err", err))
						f <- s.ManageIp
						return
					}
					return
				}(wg, server, failedChan)
			}
			wg.Wait()
			close(failedChan)
			for ip := range failedChan {
				failedIPs = append(failedIPs, ip)
			}
			t.Log("deploy "+deployInfo.DeployPath+" success", manageIpString)
		}
	} else if t.template.LogUploadServer.Mode == consts.LogServerModeSSH {
		sshClientUpload, err := common.FillSSHClient(t.template.LogUploadServer.ManageIp, t.template.Secret.Name, t.template.Secret.Password, t.template.LogUploadServer.Port)
		err = sshClientUpload.GenerateClient()
		if err != nil {
			global.GVA_LOG.Error("create ssh client failed: ", zap.Any("err", err))
			failedIPs = manageIps
			return
		}
		defer sshClientUpload.Client.Close()
		t.clients = append(t.clients, sshClientUpload.Client)
		for _, deployInfo := range t.template.TaskDeployInfos {
			if fb, err = sshClientUpload.Download(deployInfo.DownloadSource); err != nil {
				global.GVA_LOG.Error("upload via sftp failed", zap.Any("err", err))
				failedIPs = manageIps
				return
			}
			t.Log("download "+deployInfo.DownloadSource+" success", t.template.LogUploadServer.ManageIp)
			if t.template.TargetServers == nil || len(t.template.TargetServers) <= 0 {
				global.GVA_LOG.Error("run task failed on nil target server: ", zap.Uint("task ID: ", t.task.ID))
				return
			}
			global.GVA_LOG.Info("run ", zap.Uint("task ID: ", t.task.ID))
			servers := t.template.TargetServers
			wg := &sync.WaitGroup{}
			failedChan := make(chan string, len(servers))
			t.clients = make([]*ssh.Client, 0, len(servers))
			for _, server := range servers {
				wg.Add(1)
				go func(w *sync.WaitGroup, s application.ApplicationServer, f chan string) {
					defer w.Done()
					sshClient, err := common.FillSSHClient(s.ManageIp, t.template.SysUser, "", s.SshPort)
					if err != nil {
						global.GVA_LOG.Error("run task failed on create ssh client: ", zap.Uint("task ID: ", t.task.ID), zap.String("server IP: ", s.ManageIp), zap.Any("err", err))
						f <- s.ManageIp
						return
					}
					err = sshClient.GenerateClient()
					if err != nil {
						global.GVA_LOG.Error("run task failed on generate ssh client: ", zap.Uint("task ID: ", t.task.ID), zap.String("server IP: ", s.ManageIp), zap.Any("err", err))
						f <- s.ManageIp
						return
					}
					defer sshClient.Client.Close()
					t.clients = append(t.clients, sshClient.Client)
					fio := bytes.NewReader(fb)
					if err = sshClient.Upload(fio, deployInfo.DeployPath); err != nil {
						global.GVA_LOG.Error("run task failed on upload deploy file: ", zap.Uint("task ID: ", t.task.ID), zap.String("server IP: ", s.ManageIp), zap.Any("err", err))
						f <- s.ManageIp
						return
					}
					return
				}(wg, server, failedChan)
			}
			wg.Wait()
			close(failedChan)
			for ip := range failedChan {
				failedIPs = append(failedIPs, ip)
			}
			t.Log("deploy "+deployInfo.DeployPath+" success", manageIpString)
		}
	}
	return
}

func (t *TaskRunner) runDiscoverTask() (failedIPs []string) {
	//mangeIp, err := utils.GetManageIp()
	//if err != nil {
	//	return []string{"localhost"}
	//}
	//prefix := mangeIp.String()[:strings.LastIndex(mangeIp.String(), ".")+1]
	//servers := make([]string, 0, 253)
	//for i := 1; i < 254; i++ {
	//	servers = append(servers, prefix+strconv.Itoa(i))
	//}
	err, taskSystem, _, _ := applicationService.GetSystemById(float64(t.task.SystemId))
	if err != nil {
		global.GVA_LOG.Error("get task system failed", zap.Any("err", err))
		failedIPs = append(failedIPs, "localhost")
		return
	}
	servers := make([]string, 0, 253)
	manageNets := strings.Split(taskSystem.Network, "\n")
	for _, manageNet := range manageNets {
		//netDetails := strings.Split(manageNet, " ")
		//if len(netDetails) < 2 {
		//	global.GVA_LOG.Error("please check network", zap.Int("system id", t.task.SystemId))
		//	continue
		//}
		ips, err := utils.HostsCIDR(manageNet)
		if err != nil {
			global.GVA_LOG.Error("please check network", zap.Int("system id", t.task.SystemId))
			continue
		}
		servers = append(servers, ips...)
	}
	if len(servers) == 0 {
		global.GVA_LOG.Error("parse 0 server, please check network", zap.Int("system id", t.task.SystemId))
		failedIPs = append(failedIPs, "localhost")
		return
	}
	global.GVA_LOG.Info("run ", zap.Uint("task ID: ", t.task.ID))
	wg := &sync.WaitGroup{}
	successChan := make(chan bool, len(servers))
	t.clients = make([]*ssh.Client, 0, len(servers))
	for _, server := range servers {
		wg.Add(1)
		go func(w *sync.WaitGroup, s string, f chan bool) {
			defer w.Done()
			var sshPort = consts.DiscoverSSHPort
			sshClient, _ := common.FillSSHClient(s, t.task.SshUser, "", consts.DiscoverSSHPort)
			err = sshClient.GenerateClient()
			if err != nil {
				sshClient, _ = common.FillSSHClient(s, t.task.SshUser, "", 22)
				if err = sshClient.GenerateClient(); err != nil {
					return
				}
				sshPort = 22
			}
			defer sshClient.Client.Close()
			t.clients = append(t.clients, sshClient.Client)
			newServer := application.ApplicationServer{
				ManageIp: s,
				SystemId: t.task.SystemId,
				SshPort:  sshPort,
			}
			var output string
			if output, err = sshClient.CommandSingle("hostname"); err != nil {
				global.GVA_LOG.Error("get hostname failed", zap.Any("err", err))
			} else {
				lines := strings.Split(output, "\n") // discard error message from ssh login
				if len(lines) <= 1 {
					global.GVA_LOG.Error("get hostname failed")
				} else {
					newServer.Hostname = lines[len(lines)-2]
				}
			}
			if output, err = sshClient.CommandSingle("uname -a"); err != nil {
				global.GVA_LOG.Error("get architecture failed", zap.Any("err", err))
			} else {
				if strings.Contains(output, "86_64") {
					newServer.Architecture = consts.ArchitectureX86
				} else if strings.Contains(output, "aarch") || strings.Contains(output, "arm64") {
					newServer.Architecture = consts.ArchitectureArm
				}
			}
			if output, err = sshClient.CommandSingle("cat /etc/os-release | grep PRETTY_NAME"); err != nil {
				global.GVA_LOG.Error("get os failed", zap.Any("err", err))
			} else {
				if strings.Contains(output, "Kylin") {
					newServer.Os = consts.OsKylin
				} else if strings.Contains(output, "Red") || strings.Contains(output, "arm64") {
					newServer.Os = consts.OsRedhat
				} else if strings.Contains(output, "UnionTech") || strings.Contains(output, "arm64") {
					newServer.Os = consts.OsUnionTech
				} else if strings.Contains(output, "SUSE") || strings.Contains(output, "arm64") {
					newServer.Os = consts.OsSuse
				} else if strings.Contains(output, "CentOS") || strings.Contains(output, "arm64") {
					newServer.Os = consts.OsCentos
				} else if strings.Contains(output, "NeoKylin") || strings.Contains(output, "arm64") {
					newServer.Os = consts.OsKylin
				}
			}
			if output, err = sshClient.CommandSingle("cat /etc/os-release | grep VERSION_ID"); err != nil {
				global.GVA_LOG.Error("get os version failed", zap.Any("err", err))
			} else {
				if split := strings.Split(output, `"`); len(split) >= 2 {
					newServer.OsVersion = split[1]
				}
			}
			if err = applicationService.CreateOrUpdateServer(newServer); err != nil {
				global.GVA_LOG.Error("create new server failed", zap.Any("err", err))
			}
			t.Log("add or update server name " + newServer.Hostname + " manage ip " + newServer.ManageIp)
			successChan <- true
			return
		}(wg, server, successChan)
	}
	wg.Wait()
	close(successChan)
	var succeed bool
	for ip := range successChan {
		if ip {
			succeed = true
			if e, sys, _, _ := applicationService.CmdbSystemService.GetSystemById(float64(t.task.SystemId)); e != nil {
				global.GVA_LOG.Error("get system error,", zap.Any("", e))
			} else {
				needAppend := true
				for _, u := range sys.SshUsers {
					if u == t.task.SshUser {
						needAppend = false
						break
					}
				}
				if needAppend {
					sys.SshUsers = append(sys.SshUsers, t.task.SshUser)
					if e = applicationService.CmdbSystemService.UpdateSystem(appReq.AddSystem{
						ApplicationSystem: sys,
					}); e != nil {
						global.GVA_LOG.Error("update system error,", zap.Any("", e))
					}
				}
			}
			break
		}
	}
	if !succeed {
		failedIPs = append(failedIPs, "localhost")
		t.Log("find no server in "+taskSystem.Network, "localhost")
	}
	return
}

//func hasRequirementsChanges(requirementsFilePath string, requirementsHashFilePath string) bool {
//	oldFileMD5HashBytes, err := ioutil.ReadFile(requirementsHashFilePath)
//	if err != nil {
//		return true
//	}
//
//	newFileMD5Hash, err := getMD5Hash(requirementsFilePath)
//	if err != nil {
//		return true
//	}
//
//	return string(oldFileMD5HashBytes) != newFileMD5Hash
//}

//func writeMD5Hash(requirementsFile string, requirementsHashFile string) error {
//	newFileMD5Hash, err := getMD5Hash(requirementsFile)
//	if err != nil {
//		return err
//	}
//
//	return ioutil.WriteFile(requirementsHashFile, []byte(newFileMD5Hash), 0644)
//}

// checkTmpDir checks to see if the temporary directory exists
// and if it does not attempts to create it
func checkTmpDir(path string) error {
	var err error
	if _, err = os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return os.MkdirAll(path, 0700)
		}
	}
	return err
}
