package tasks

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ansible"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var AnsibleTaskPool     *TaskPool

type logRecord struct {
	task   *TaskRunner
	output string
	time   time.Time
}

type resourceLock struct {
	lock   bool
	holder *TaskRunner
}

type TaskPool struct {
	// queue contains list of tasks in status TaskWaitingStatus.
	queue []*TaskRunner

	// register channel used to put tasks to queue.
	register chan *TaskRunner

	activeProj map[int]map[int]*TaskRunner

	// runningTasks contains tasks with status TaskRunningStatus.
	runningTasks map[int]*TaskRunner

	// logger channel used to putting log records to database.
	logger chan logRecord

	resourceLocker chan *resourceLock
}

func (p *TaskPool) GetTask(id int) (task *TaskRunner) {

	for _, t := range p.queue {
		if int(t.task.ID) == id {
			task = t
			break
		}
	}

	if task == nil {
		for _, t := range p.runningTasks {
			if int(t.task.ID) == id {
				task = t
				break
			}
		}
	}

	return
}

//nolint: gocyclo
func (p *TaskPool) Run() {
	ticker := time.NewTicker(5 * time.Second)

	defer func() {
		close(p.resourceLocker)
		ticker.Stop()
	}()

	// Lock or unlock resources when running a TaskRunner
	go func(locker <-chan *resourceLock) {
		for l := range locker {
			t := l.holder

			if l.lock {
				if p.blocks(t) {
					panic("Trying to lock an already locked resource!")
				}

				projTasks, ok := p.activeProj[t.task.ProjectID]
				if !ok {
					projTasks = make(map[int]*TaskRunner)
					p.activeProj[t.task.ProjectID] = projTasks
				}
				projTasks[int(t.task.ID)] = t
				p.runningTasks[int(t.task.ID)] = t
				continue
			}

			if p.activeProj[t.task.ProjectID] != nil && p.activeProj[t.task.ProjectID][int(t.task.ID)] != nil {
				delete(p.activeProj[t.task.ProjectID], int(t.task.ID))
				if len(p.activeProj[t.task.ProjectID]) == 0 {
					delete(p.activeProj, t.task.ProjectID)
				}
			}

			delete(p.runningTasks, int(t.task.ID))
		}
	}(p.resourceLocker)

	for {
		select {
		case record := <-p.logger: // new log message which should be put to database
			_, err := taskService.CreateTaskOutput(ansible.TaskOutput{
				TaskID: int(record.task.task.ID),
				Output: record.output,
				Time:   record.time,
			})
			if err != nil {
				global.GVA_LOG.Error(err.Error())
			}
		case task := <-p.register: // new task created by API or schedule
			p.queue = append(p.queue, task)
			msg := "Task " + strconv.Itoa(int(task.task.ID)) + " added to queue"
			task.Log(msg)
			global.GVA_LOG.Info(msg)
			task.updateStatus()

		case <-ticker.C: // timer 5 seconds
			if len(p.queue) == 0 {
				continue
			}

			//get TaskRunner from top of queue
			t := p.queue[0]
			if t.task.Status == ansible.TaskFailStatus {
				//delete failed TaskRunner from queue
				p.queue = p.queue[1:]
				global.GVA_LOG.Info("Task " + strconv.Itoa(int(t.task.ID)) + " removed from queue")
				continue
			}
			if p.blocks(t) {
				//move blocked TaskRunner to end of queue
				p.queue = append(p.queue[1:], t)
				continue
			}
			global.GVA_LOG.Info("Set resource locker with TaskRunner " + strconv.Itoa(int(t.task.ID)))
			p.resourceLocker <- &resourceLock{lock: true, holder: t}
			if !t.prepared {
				go t.prepareRun()
				continue
			}
			go t.run()
			p.queue = p.queue[1:]
			global.GVA_LOG.Info("Task " + strconv.Itoa(int(t.task.ID)) + " removed from queue")
		}
	}
}

func (p *TaskPool) blocks(t *TaskRunner) bool {

	if len(p.runningTasks) >= global.GVA_CONFIG.Ansible.MaxParallelTasks {
		return true
	}

	if p.activeProj[t.task.ProjectID] == nil || len(p.activeProj[t.task.ProjectID]) == 0 {
		return false
	}

	for _, r := range p.activeProj[t.task.ProjectID] {
		if int(r.template.ID) == t.task.TemplateID {
			return true
		}
	}

	proj, err := projectService.GetProject(t.task.ProjectID)

	if err != nil {
		global.GVA_LOG.Error(err.Error())
		return false
	}

	return proj.MaxParallelTasks > 0 && len(p.activeProj[t.task.ProjectID]) >= proj.MaxParallelTasks
}

func CreateTaskPool() TaskPool {
	return TaskPool{
		queue:          make([]*TaskRunner, 0), // queue of waiting tasks
		register:       make(chan *TaskRunner), // add TaskRunner to queue
		activeProj:     make(map[int]map[int]*TaskRunner),
		runningTasks:   make(map[int]*TaskRunner),   // working tasks
		logger:         make(chan logRecord, 10000), // store log records to database
		resourceLocker: make(chan *resourceLock),
	}
}

func (p *TaskPool) StopTask(targetTask ansible.Task) error {
	tsk := p.GetTask(int(targetTask.ID))
	if tsk == nil { // task not active, but exists in database
		tsk = &TaskRunner{
			task: targetTask,
			pool: p,
		}
		err := tsk.populateDetails()
		if err != nil {
			return err
		}
		tsk.setStatus(ansible.TaskStoppedStatus)
	} else {
		status := tsk.task.Status
		tsk.setStatus(ansible.TaskStoppingStatus)
		if status == ansible.TaskRunningStatus {
			if tsk.process == nil {
				panic("running process can not be nil")
			}
			err := tsk.process.Kill()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func getNextBuildVersion(startVersion string, currentVersion string) string {
	re := regexp.MustCompile(`^(.*[^\d])?(\d+)([^\d].*)?$`)
	m := re.FindStringSubmatch(startVersion)

	if m == nil {
		return startVersion
	}

	var prefix, suffix, body string

	switch len(m) - 1 {
	case 3:
		prefix = m[1]
		body = m[2]
		suffix = m[3]
	case 2:
		if _, err := strconv.Atoi(m[1]); err == nil {
			body = m[1]
			suffix = m[2]
		} else {
			prefix = m[1]
			body = m[2]
		}
	case 1:
		body = m[1]
	default:
		return startVersion
	}

	if !strings.HasPrefix(currentVersion, prefix) ||
		!strings.HasSuffix(currentVersion, suffix) {
		return startVersion
	}

	curr, err := strconv.Atoi(currentVersion[len(prefix) : len(currentVersion)-len(suffix)])
	if err != nil {
		return startVersion
	}

	start, err := strconv.Atoi(body)
	if err != nil {
		panic(err)
	}

	var newVer int
	if start > curr {
		newVer = start
	} else {
		newVer = curr + 1
	}

	return prefix + strconv.Itoa(newVer) + suffix
}

func (p *TaskPool) AddTask(taskObj ansible.Task, userID *int, projectID int) (newTask ansible.Task, err error) {
	taskObj.Created = time.Now()
	taskObj.Status = ansible.TaskWaitingStatus
	taskObj.UserID = userID
	taskObj.ProjectID = projectID

	tpl, err := templateService.GetTemplate(float64(projectID), float64(taskObj.TemplateID))
	if err != nil {
		return
	}

	err = taskService.ValidateNewTask(tpl)
	if err != nil {
		return
	}

	if tpl.Type == ansible.TemplateBuild { // get next version for TaskRunner if it is a Build
		var builds []ansible.TaskWithTpl
		e, buildList, _ := taskService.GetTemplateTasks(tpl.ProjectID, int(tpl.ID), request.PageInfo{
			Page:     1,
			PageSize: 1,
		})
		if e != nil {
			err = e
			return
		}
		builds = buildList.([]ansible.TaskWithTpl)
		if len(builds) == 0 || builds[0].Task.Version == nil {
			taskObj.Version = tpl.StartVersion
		} else {
			v := getNextBuildVersion(*tpl.StartVersion, *builds[0].Task.Version)
			taskObj.Version = &v
		}
	}

	newTask, err = taskService.CreateTask(taskObj)
	if err != nil {
		return
	}

	taskRunner := TaskRunner{
		task: newTask,
		pool: p,
	}

	err = taskRunner.populateDetails()
	if err != nil {
		taskRunner.Log("Error: " + err.Error())
		taskRunner.fail()
		return
	}

	p.register <- &taskRunner

	return
}
