package common

import (
	"golang.org/x/crypto/ssh"
	"os/exec"
)

type Logger interface {
	Log(msg string)
	LogCmd(cmd *exec.Cmd)
	LogSsh(channel *ssh.Channel)
}
