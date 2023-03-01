package taskPool

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	sockets "github.com/flipped-aurora/gin-vue-admin/server/api/v1/socket"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
	"github.com/flipped-aurora/gin-vue-admin/server/model/taskMdl"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"io"
	"os"
	"strconv"
	"time"
)

type TaskRunner struct {
	task     taskMdl.Task
	template taskMdl.TaskTemplate

	users []int
	//alert     bool
	//alertChat *string
	prepared bool
	process  *os.Process
	pool     *TaskPool
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
			"type":        "update",
			"start":       t.task.BeginTime,
			"end":         t.task.EndTime,
			"status":      t.task.Status,
			"task_id":     t.task.ID,
			"template_id": t.task.TemplateId,
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

		t.task.EndTime = time.Now()
		t.updateStatus()
		//t.createTaskEvent()
		//t.destroyKeys()
	}()

	// TODO: more details
	if t.task.Status == taskMdl.TaskStoppingStatus {
		t.setStatus(taskMdl.TaskStoppedStatus)
		return
	}

	t.task.BeginTime = time.Now()
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

	err := t.runTask()
	if err != nil {
		t.Log("Running task failed: " + err.Error())
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

	t.template, err = taskService.GetTaskTemplate(float64(t.task.TemplateId))
	if err != nil {
		return t.prepareError(err, "Template not found!")
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

func (t *TaskRunner) runTask() (err error) {
	servers := t.template.TargetServers
	for _, server := range servers {
		sshClient, err := sshService.FillSSHClient(server.ManageIp, t.template.SysUser, "", server.SshPort)
		err = sshClient.GenerateClient()
		if err != nil {
			return err
		}
		sshClient.RequestShell()
		if err = sshClient.ConnectShell(t.template.Command, t); err != nil {
			return err
		}
	}
	return nil
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
