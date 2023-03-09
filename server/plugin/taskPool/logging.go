package taskPool

import (
	"bufio"
	"encoding/json"
	"github.com/flipped-aurora/gin-vue-admin/server/api/v1/socket"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
	"os/exec"
	"time"
	"unicode/utf8"
)

func (t *TaskRunner) Log(msg string) {
	now := time.Now()

	for _, user := range t.users {
		b, err := json.Marshal(&map[string]interface{}{
			"type":       "log",
			"output":     msg,
			"recordTime": now,
			"taskId":     t.task.ID,
			"ID":         time.Now().UnixNano(),
		})

		if err != nil {
			global.GVA_LOG.Fatal(err.Error())
		}

		socket.Message(user, b)
	}

	t.pool.logger <- logRecord{
		task:   t,
		output: msg,
		time:   now,
	}
}

// Readln reads from the pipe
func Readln(r *bufio.Reader) (string, error) {
	//var (
	//	isPrefix = true
	//	err      error
	//	line, ln []byte
	//)
	//for isPrefix && err == nil {
	//	line, isPrefix, err = r.ReadLine()
	//	ln = append(ln, line...)
	//}
	//return string(ln), err

	var (
		err  error
		line []byte
	)
	for err == nil {
		global.GVA_LOG.Info("starting read line")
		x, size, e := r.ReadRune()
		err = e
		global.GVA_LOG.Info("start read line")
		if size > 0 && e == nil {
			if x != utf8.RuneError {
				p := make([]byte, utf8.RuneLen(x))
				utf8.EncodeRune(p, x)
				line = append(line, p...)
			} else {
				line = append(line, []byte("@")...)
			}
		}
	}
	return string(line), err
}

func (t *TaskRunner) logPipe(reader *bufio.Reader) {

	line, err := Readln(reader)
	for err == nil {
		global.GVA_LOG.Info("read output", zap.String(" ", line))
		t.Log(line)
		line, err = Readln(reader)
	}

	if err != nil && err.Error() != "EOF" {
		//don't panic on these errors, sometimes it throws not dangerous "read |0: file already closed" error
		global.GVA_LOG.Warn("Failed to read TaskRunner output", zap.Any("err", err))
	}

}

func (t *TaskRunner) LogCmd(cmd *exec.Cmd) {
	stderr, _ := cmd.StderrPipe()
	stdout, _ := cmd.StdoutPipe()

	go t.logPipe(bufio.NewReader(stderr))
	go t.logPipe(bufio.NewReader(stdout))
}

func (t *TaskRunner) LogSsh(channel ssh.Channel) {
	go t.logPipe(bufio.NewReader(channel))
}

func (t *TaskRunner) panicOnError(err error, msg string) {
	if err != nil {
		t.Log(msg)
		global.GVA_LOG.Fatal(msg, zap.Any("err", err))
	}
}
