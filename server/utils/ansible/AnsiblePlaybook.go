package lib

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type AnsiblePlaybook struct {
	TemplateID int
	Logger     Logger
}

func (p AnsiblePlaybook) makeCmd(command, tmpPath string, args []string) *exec.Cmd {
	cmd := exec.Command(command, args...) //nolint: gas
	cmd.Dir = tmpPath

	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf(tmpPath))
	cmd.Env = append(cmd.Env, fmt.Sprintf("PWD=%s", cmd.Dir))
	cmd.Env = append(cmd.Env, fmt.Sprintln("PYTHONUNBUFFERED=1"))

	return cmd
}

func (p AnsiblePlaybook) runCmd(command, tmpPath string, args []string) error {
	cmd := p.makeCmd(command, tmpPath, args)
	p.Logger.LogCmd(cmd)
	return cmd.Run()
}

func (p AnsiblePlaybook) GetHosts(tmpPath string, args []string) (hosts []string, err error) {
	args = append(args, "--list-hosts")
	cmd := p.makeCmd("ansible-playbook", tmpPath, args)

	var errb bytes.Buffer
	cmd.Stderr = &errb

	out, err := cmd.Output()
	if err != nil {
		return
	}

	re := regexp.MustCompile(`(?m)^\\s{6}(.*)$`)
	matches := re.FindAllSubmatch(out, 20)
	hosts = make([]string, len(matches))
	for i := range matches {
		hosts[i] = string(matches[i][1])
	}

	return
}

func (p AnsiblePlaybook) RunPlaybook(tmpPath string, args []string, cb func(*os.Process)) error {
	cmd := p.makeCmd("ansible-playbook", tmpPath, args)
	p.Logger.LogCmd(cmd)
	cmd.Stdin = strings.NewReader("")
	err := cmd.Start()
	if err != nil {
		return err
	}
	cb(cmd.Process)
	return cmd.Wait()
}
