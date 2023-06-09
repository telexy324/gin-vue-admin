package taskPool

import (
	"database/sql"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/taskMdl"
	"go.uber.org/zap"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var TPool TaskPool

type logRecord struct {
	task     *TaskRunner
	output   string
	time     time.Time
	manageIp string
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

	activeTask map[int]*TaskRunner

	// runningTasks contains tasks with status TaskRunningStatus.
	runningTasks map[int]*TaskRunner

	// logger channel used to putting log records to database.
	logger chan logRecord

	//store db.Store

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

// nolint: gocyclo
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

				p.activeTask[int(t.task.ID)] = t
				p.runningTasks[int(t.task.ID)] = t
				continue
			}

			if p.activeTask != nil && p.activeTask[int(t.task.ID)] != nil {
				delete(p.activeTask, int(t.task.ID))
				//if len(p.activeTask) == 0 {
				//	delete(p.activeProj, t.task.ProjectID)
				//}
			}

			delete(p.runningTasks, int(t.task.ID))
		}
	}(p.resourceLocker)

	for {
		select {
		case record := <-p.logger: // new log message which should be put to database
			_, err := taskService.CreateTaskOutput(taskMdl.TaskOutput{
				TaskId:     int(record.task.task.ID),
				Output:     record.output,
				RecordTime: record.time,
				ManageIp:   record.manageIp,
			})

			if err != nil {
				global.GVA_LOG.Error(err.Error())
			}
		case task := <-p.register: // new task created by API or schedule
			p.queue = append(p.queue, task)
			//log.Debug(task)
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
			if t.task.Status == taskMdl.TaskFailStatus {
				//delete failed TaskRunner from queue
				p.queue = p.queue[1:]
				global.GVA_LOG.Info("Task removed from queue ", zap.Uint("Task ID ", t.task.ID))
				continue
			}
			if p.blocks(t) {
				//move blocked TaskRunner to end of queue
				p.queue = append(p.queue[1:], t)
				continue
			}
			global.GVA_LOG.Info("Set resource locker with ", zap.Uint("TaskRunner ", t.task.ID))
			p.resourceLocker <- &resourceLock{lock: true, holder: t}
			if !t.prepared {
				go t.prepareRun()
				continue
			}
			go t.run()
			p.queue = p.queue[1:]
			global.GVA_LOG.Info("Task removed from queue ", zap.Uint("Task ID ", t.task.ID))
		}
	}
}

func (p *TaskPool) blocks(t *TaskRunner) bool {

	if len(p.runningTasks) >= global.GVA_CONFIG.Task.MaxParallelTasks {
		return true
	}

	if p.activeTask == nil || len(p.activeTask) == 0 {
		return false
	}

	for _, r := range p.activeTask {
		if int(r.template.ID) == t.task.TemplateId {
			return true
		}
	}

	//proj, err := p.store.GetProject(t.task.ProjectID)
	//
	//if err != nil {
	//	log.Error(err)
	//	return false
	//}
	//
	//return proj.MaxParallelTasks > 0 && len(p.activeProj[t.task.ProjectID]) >= proj.MaxParallelTasks
	return false
}

func CreateTaskPool() TaskPool {
	TPool = TaskPool{
		queue:          make([]*TaskRunner, 0), // queue of waiting tasks
		register:       make(chan *TaskRunner), // add TaskRunner to queue
		activeTask:     make(map[int]*TaskRunner),
		runningTasks:   make(map[int]*TaskRunner),   // working tasks
		logger:         make(chan logRecord, 10000), // store log records to database
		resourceLocker: make(chan *resourceLock),
	}
	go TPool.Run()
	return TPool
}

func (p *TaskPool) StopTask(targetTask taskMdl.Task) (failedIps []string, err error) {
	tsk := p.GetTask(int(targetTask.ID))
	failedIps = make([]string, 0)
	if tsk == nil { // task not active, but exists in database
		tsk = &TaskRunner{
			task: targetTask,
			pool: p,
		}
		err = tsk.populateDetails()
		if err != nil {
			return
		}
		tsk.setStatus(taskMdl.TaskStoppedStatus)
		//tsk.createTaskEvent()
	} else {
		status := tsk.task.Status
		tsk.setStatus(taskMdl.TaskStoppingStatus)
		if status == taskMdl.TaskRunningStatus {
			if tsk.clients == nil || len(tsk.clients) == 0 {
				panic("running process can not be nil")
			}
			for _, client := range tsk.clients {
				if err = client.Close(); err != nil {
					global.GVA_LOG.Error("close client failed", zap.Uint("Task ID ", targetTask.ID), zap.String("client ip ", client.RemoteAddr().String()))
					failedIps = append(failedIps, client.RemoteAddr().String())
				}
			}
			for _, client := range tsk.ftpConn {
				if err = client.Quit(); err != nil {
					global.GVA_LOG.Error("close ftp client failed", zap.Uint("Task ID ", targetTask.ID))
					//failedIps = append(failedIps, client.RemoteAddr().String())
				}
			}
		}
	}
	return
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

func (p *TaskPool) AddTask(taskObj taskMdl.Task, userID int, setTaskId int) (newTask taskMdl.Task, err error) {
	taskObj.Status = taskMdl.TaskWaitingStatus
	taskObj.SystemUserId = userID
	taskObj.BeginTime = sql.NullTime{
		Time:  time.Unix(0, 0),
		Valid: true,
	}
	taskObj.EndTime = sql.NullTime{
		Time:  time.Unix(0, 0),
		Valid: true,
	}
	taskObj.SetTaskId = setTaskId

	//tpl, err := taskService.GetTaskTemplate(float64(taskObj.TemplateId))
	//if err != nil {
	//	return
	//}

	//err = taskObj.ValidateNewTask(tpl)
	//if err != nil {
	//	return
	//}

	//if tpl.Type == db.TemplateBuild { // get next version for TaskRunner if it is a Build
	//	var builds []db.TaskWithTpl
	//	builds, err = p.store.GetTemplateTasks(tpl.ProjectID, tpl.ID, db.RetrieveQueryParams{Count: 1})
	//	if err != nil {
	//		return
	//	}
	//	if len(builds) == 0 || builds[0].Version == nil {
	//		taskObj.Version = tpl.StartVersion
	//	} else {
	//		v := getNextBuildVersion(*tpl.StartVersion, *builds[0].Version)
	//		taskObj.Version = &v
	//	}
	//}

	newTask, err = taskService.CreateTask(taskObj)
	newTask.FileDownload, newTask.SystemId, newTask.SshUser = taskObj.FileDownload, taskObj.SystemId, taskObj.SshUser
	if err != nil {
		return
	}
	if newTask.TemplateId < 99999900 {
		template, er := taskService.GetTaskTemplate(float64(taskObj.TemplateId))
		if er != nil {
			err = er
			return
		}
		template.LastTaskId = int(newTask.ID)
		err = taskService.UpdateTaskTemplate(template)
		if err != nil {
			return
		}
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

	//objType := db.EventTask
	//desc := "Task ID " + strconv.Itoa(newTask.ID) + " queued for running"
	//_, err = p.store.CreateEvent(db.Event{
	//	UserID:      userID,
	//	ProjectID:   &projectID,
	//	ObjectType:  &objType,
	//	ObjectID:    &newTask.ID,
	//	Description: &desc,
	//})

	return
}
