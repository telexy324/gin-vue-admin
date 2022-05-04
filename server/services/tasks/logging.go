package tasks

import (
	"bufio"
	"encoding/json"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"go.uber.org/zap"
	"os/exec"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/sockets"
)

func (t *TaskRunner) Log(msg string) {
	now := time.Now()

	for _, user := range t.users {
		b, err := json.Marshal(&map[string]interface{}{
			"type":       "log",
			"output":     msg,
			"time":       now,
			"task_id":    t.task.ID,
			"project_id": t.task.ProjectID,
		})

		global.GVA_LOG.Panic(err.Error(), zap.Any("level", "Panic"))

		sockets.Message(user, b)
	}

	t.pool.logger <- logRecord{
		task:   t,
		output: msg,
		time:   now,
	}
}

// Readln reads from the pipe
func Readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix = true
		err      error
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}

func (t *TaskRunner) logPipe(reader *bufio.Reader) {

	line, err := Readln(reader)
	for err == nil {
		t.Log(line)
		line, err = Readln(reader)
	}

	if err != nil && err.Error() != "EOF" {
		//don't panic on this errors, sometimes it throw not dangerous "read |0: file already closed" error
		global.GVA_LOG.Panic(err.Error(), zap.Any("level", "Warning"), zap.Any("error", "Failed to read TaskRunner output"))
	}

}

func (t *TaskRunner) LogCmd(cmd *exec.Cmd) {
	stderr, _ := cmd.StderrPipe()
	stdout, _ := cmd.StdoutPipe()

	go t.logPipe(bufio.NewReader(stderr))
	go t.logPipe(bufio.NewReader(stdout))
}

func (t *TaskRunner) panicOnError(err error, msg string) {
	if err != nil {
		t.Log(msg)
		global.GVA_LOG.Panic(err.Error(), zap.Any("level", "Warning"))
	}
}
