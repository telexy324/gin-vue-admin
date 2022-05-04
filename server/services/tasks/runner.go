package tasks

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/ansible-semaphore/semaphore/lib"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ansible"
	request2 "github.com/flipped-aurora/gin-vue-admin/server/model/ansible/request"
	utils "github.com/flipped-aurora/gin-vue-admin/server/utils/ansible"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/sockets"
)

const (
	gitURLFilePrefix = "file://"
)

type TaskRunner struct {
	task        ansible.Task
	template    ansible.Template
	inventory   ansible.Inventory
	environment ansible.Environment

	users     []int
	alert     bool
	alertChat *string
	prepared  bool
	process   *os.Process
	pool      *TaskPool
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

func (t *TaskRunner) setStatus(status ansible.TaskStatus) {
	if t.task.Status == ansible.TaskStoppingStatus {
		switch status {
		case ansible.TaskFailStatus:
			status = ansible.TaskStoppedStatus
		case ansible.TaskStoppedStatus:
		default:
			panic("stopping TaskRunner cannot be " + status)
		}
	}

	t.task.Status = status

	t.updateStatus()

	if status == ansible.TaskFailStatus {
		t.sendMailAlert()
	}

	if status == ansible.TaskSuccessStatus || status == ansible.TaskFailStatus {
		t.sendTelegramAlert()
	}
}

func (t *TaskRunner) updateStatus() {
	for _, user := range t.users {
		b, err := json.Marshal(&map[string]interface{}{
			"type":        "update",
			"start":       t.task.Start,
			"end":         t.task.End,
			"status":      t.task.Status,
			"task_id":     t.task.ID,
			"template_id": t.task.TemplateID,
			"project_id":  t.task.ProjectID,
			"version":     t.task.Version,
		})

		global.GVA_LOG.Panic(err.Error())

		sockets.Message(user, b)
	}

	if err := taskService.UpdateTask(t.task); err != nil {
		t.panicOnError(err, "Failed to update TaskRunner status")
	}
}

func (t *TaskRunner) fail() {
	t.setStatus(ansible.TaskFailStatus)
}

func (t *TaskRunner) destroyKeys() {
	err := keyService.Destroy(&t.inventory.SSHKey)
	if err != nil {
		t.Log("Can't destroy inventory user key, error: " + err.Error())
	}

	err = keyService.Destroy(&t.inventory.BecomeKey)
	if err != nil {
		t.Log("Can't destroy inventory become user key, error: " + err.Error())
	}

	err = keyService.Destroy(&t.template.VaultKey)
	if err != nil {
		t.Log("Can't destroy inventory vault password file, error: " + err.Error())
	}
}

//func (t *TaskRunner) createTaskEvent() {
//	objType := db.EventTask
//	desc := "Task ID " + strconv.Itoa(t.task.ID) + " (" + t.template.Name + ")" + " finished - " + strings.ToUpper(string(t.task.Status))
//
//	_, err := t.pool.store.CreateEvent(db.Event{
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
		global.GVA_LOG.Info("Stopped preparing TaskRunner " + strconv.Itoa(int(t.task.ID)))
		global.GVA_LOG.Info("Release resource locker with TaskRunner " + strconv.Itoa(int(t.task.ID)))
		t.pool.resourceLocker <- &resourceLock{lock: false, holder: t}
	}()

	t.Log("Preparing: " + strconv.Itoa(int(t.task.ID)))

	err := checkTmpDir(global.GVA_CONFIG.Ansible.TmpPath)
	if err != nil {
		t.Log("Creating tmp dir failed: " + err.Error())
		t.fail()
		return
	}

	if err != nil {
		t.Log("Fatal error inserting an event")
		panic(err)
	}

	t.Log("Prepare TaskRunner with template: " + t.template.Name + "\n")

	t.updateStatus()

	if err := t.installInventory(); err != nil {
		t.Log("Failed to install inventory: " + err.Error())
		t.fail()
		return
	}

	if err := t.installVaultKeyFile(); err != nil {
		t.Log("Failed to install vault password file: " + err.Error())
		t.fail()
		return
	}

	t.prepared = true
}

func (t *TaskRunner) run() {
	defer func() {
		global.GVA_LOG.Info("Stopped running TaskRunner " + strconv.Itoa(int(t.task.ID)))
		global.GVA_LOG.Info("Release resource locker with TaskRunner " + strconv.Itoa(int(t.task.ID)))
		t.pool.resourceLocker <- &resourceLock{lock: false, holder: t}

		now := time.Now()
		t.task.End = &now
		t.updateStatus()
		t.destroyKeys()
	}()

	// TODO: more details
	if t.task.Status == ansible.TaskStoppingStatus {
		t.setStatus(ansible.TaskStoppedStatus)
		return
	}

	now := time.Now()
	t.task.Start = &now
	t.setStatus(ansible.TaskRunningStatus)

	t.Log("Started: " + strconv.Itoa(int(t.task.ID)))
	t.Log("Run TaskRunner with template: " + t.template.Name + "\n")

	// TODO: ?????
	if t.task.Status == ansible.TaskStoppingStatus {
		t.setStatus(ansible.TaskStoppedStatus)
		return
	}

	err := t.runPlaybook()
	if err != nil {
		t.Log("Running playbook failed: " + err.Error())
		t.fail()
		return
	}

	t.setStatus(ansible.TaskSuccessStatus)

	err, list, _ := templateService.GetTemplates(request2.GetByProjectId{
		ProjectId: float64(t.task.ProjectID),
	}, ansible.TemplateFilter{
		BuildTemplateID: &t.task.TemplateID,
		AutorunOnly:     true,
	})
	if err != nil {
		t.Log("Running playbook failed: " + err.Error())
		return
	}

	templates := list.([]ansible.Template)
	for _, tpl := range templates {
		taskID := int(t.task.ID)
		_, err = t.pool.AddTask(ansible.Task{
			TemplateID:  tpl.ID,
			ProjectID:   tpl.ProjectID,
			BuildTaskID: &taskID,
		}, nil, tpl.ProjectID)
		if err != nil {
			t.Log("Running playbook failed: " + err.Error())
			continue
		}
	}
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

//nolint: gocyclo
func (t *TaskRunner) populateDetails() error {
	// get template
	var err error

	t.template, err = templateService.GetTemplate(float64(t.task.ProjectID), float64(t.task.TemplateID))
	if err != nil {
		return t.prepareError(err, "Template not found!")
	}

	// get project alert setting
	project, err := projectService.GetProject(t.template.ProjectID)
	if err != nil {
		return t.prepareError(err, "Project not found!")
	}

	t.alert = project.Alert
	t.alertChat = project.AlertChat

	// get project users
	users, err := userService.GetProjectUsers(t.template.ProjectID)
	if err != nil {
		return t.prepareError(err, "Users not found!")
	}

	t.users = []int{}
	for _, user := range users {
		t.users = append(t.users, int(user.ID))
	}

	// get inventory
	t.inventory, err = inventoryService.GetInventory(float64(t.template.ProjectID), float64(t.template.InventoryID))
	if err != nil {
		return t.prepareError(err, "Template Inventory not found!")
	}

	// get environment
	if t.template.EnvironmentID != nil {
		t.environment, err = environmentService.GetEnvironment(float64(t.template.ProjectID), float64(*t.template.EnvironmentID))
		if err != nil {
			return err
		}
	}

	if t.task.Environment != "" {
		environment := make(map[string]interface{})
		if t.environment.JSON != "" {
			err = json.Unmarshal([]byte(t.task.Environment), &environment)
			if err != nil {
				return err
			}
		}

		taskEnvironment := make(map[string]interface{})
		err = json.Unmarshal([]byte(t.environment.JSON), &taskEnvironment)
		if err != nil {
			return err
		}

		for k, v := range taskEnvironment {
			environment[k] = v
		}

		var ev []byte
		ev, err = json.Marshal(environment)
		if err != nil {
			return err
		}

		t.environment.JSON = string(ev)
	}

	return nil
}

func (t *TaskRunner) installVaultKeyFile() error {
	if t.template.VaultKeyID == nil {
		return nil
	}

	return keyService.Install(&t.template.VaultKey, ansible.AccessKeyRoleAnsiblePasswordVault)
}

func (t *TaskRunner) runPlaybook() (err error) {
	args, err := t.getPlaybookArgs()
	if err != nil {
		return
	}

	return utils.AnsiblePlaybook{
		Logger:     t,
		TemplateID: t.template.ID,
	}.RunPlaybook(args, func(p *os.Process) { t.process = p })
}

func (t *TaskRunner) getEnvironmentExtraVars() (str string, err error) {
	extraVars := make(map[string]interface{})

	if t.environment.JSON != "" {
		err = json.Unmarshal([]byte(t.environment.JSON), &extraVars)
		if err != nil {
			return
		}
	}

	taskDetails := make(map[string]interface{})

	if t.task.Message != "" {
		taskDetails["message"] = t.task.Message
	}

	if t.task.UserID != nil {
		var user ansible.User
		user, err = t.pool.store.GetUser(*t.task.UserID)
		if err == nil {
			taskDetails["username"] = user.Username
		}
	}

	if t.template.Type != ansible.TemplateTask {
		taskDetails["type"] = t.template.Type
		incomingVersion := t.task.GetIncomingVersion(t.pool.store)
		if incomingVersion != nil {
			taskDetails["incoming_version"] = incomingVersion
		}
		if t.template.Type == ansible.TemplateBuild {
			taskDetails["target_version"] = t.task.Version
		}
	}

	vars := make(map[string]interface{})
	vars["task_details"] = taskDetails
	extraVars["semaphore_vars"] = vars

	ev, err := json.Marshal(extraVars)
	if err != nil {
		return
	}

	str = string(ev)

	return
}

//nolint: gocyclo
func (t *TaskRunner) getPlaybookArgs() (args []string, err error) {
	playbookName := t.task.Playbook
	if playbookName == "" {
		playbookName = t.template.Playbook
	}

	var inventory string
	switch t.inventory.Type {
	case ansible.InventoryFile:
		inventory = t.inventory.Inventory
	case ansible.InventoryStatic:
		inventory = global.GVA_CONFIG.Ansible.TmpPath + "/inventory_" + strconv.Itoa(int(t.task.ID))
	default:
		err = fmt.Errorf("invalid invetory type")
		return
	}

	args = []string{
		"-i", inventory,
	}

	if t.inventory.SSHKeyID != nil {
		switch t.inventory.SSHKey.Type {
		case ansible.AccessKeySSH:
			args = append(args, "--private-key="+t.inventory.SSHKey.GetPath())
			//args = append(args, "--extra-vars={\"ansible_ssh_private_key_file\": \""+t.inventory.SSHKey.GetPath()+"\"}")
			if t.inventory.SSHKey.SshKey.Login != "" {
				args = append(args, "--extra-vars={\"ansible_user\": \""+t.inventory.SSHKey.SshKey.Login+"\"}")
			}
		case ansible.AccessKeyLoginPassword:
			args = append(args, "--extra-vars=@"+t.inventory.SSHKey.GetPath())
		case ansible.AccessKeyNone:
		default:
			err = fmt.Errorf("access key does not suite for inventory's user credentials")
			return
		}
	}

	if t.inventory.BecomeKeyID != nil {
		switch t.inventory.BecomeKey.Type {
		case ansible.AccessKeyLoginPassword:
			args = append(args, "--extra-vars=@"+t.inventory.BecomeKey.GetPath())
		case ansible.AccessKeyNone:
		default:
			err = fmt.Errorf("access key does not suite for inventory's sudo user credentials")
			return
		}
	}

	if t.task.Debug {
		args = append(args, "-vvvv")
	}

	if t.task.DryRun {
		args = append(args, "--check")
	}

	if t.template.VaultKeyID != nil {
		args = append(args, "--vault-password-file", t.template.VaultKey.GetPath())
	}

	extraVars, err := t.getEnvironmentExtraVars()
	if err != nil {
		t.Log(err.Error())
		t.Log("Could not remove command environment, if existant it will be passed to --extra-vars. This is not fatal but be aware of side effects")
	} else if extraVars != "" {
		args = append(args, "--extra-vars", extraVars)
	}

	var templateExtraArgs []string
	if t.template.Arguments != nil {
		err = json.Unmarshal([]byte(*t.template.Arguments), &templateExtraArgs)
		if err != nil {
			t.Log("Invalid format of the template extra arguments, must be valid JSON")
			return
		}
	}

	var taskExtraArgs []string
	if t.template.AllowOverrideArgsInTask && t.task.Arguments != nil {
		err = json.Unmarshal([]byte(*t.task.Arguments), &taskExtraArgs)
		if err != nil {
			t.Log("Invalid format of the TaskRunner extra arguments, must be valid JSON")
			return
		}
	}

	args = append(args, templateExtraArgs...)
	args = append(args, taskExtraArgs...)
	args = append(args, playbookName)

	return
}

func hasRequirementsChanges(requirementsFilePath string, requirementsHashFilePath string) bool {
	oldFileMD5HashBytes, err := ioutil.ReadFile(requirementsHashFilePath)
	if err != nil {
		return true
	}

	newFileMD5Hash, err := getMD5Hash(requirementsFilePath)
	if err != nil {
		return true
	}

	return string(oldFileMD5HashBytes) != newFileMD5Hash
}

func writeMD5Hash(requirementsFile string, requirementsHashFile string) error {
	newFileMD5Hash, err := getMD5Hash(requirementsFile)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(requirementsHashFile, []byte(newFileMD5Hash), 0644)
}

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
