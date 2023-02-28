package taskRunnerSvr

//type TaskWithLogger struct {
//	TemplateID int
//	Logger     Logger
//}
//
//func (p TaskWithLogger) makeCmd(command string, args []string, environmentVars *[]string) *exec.Cmd {
//	cmd := exec.Command(command, args...) //nolint: gas
//
//	cmd.Env = os.Environ()
//	cmd.Env = append(cmd.Env, fmt.Sprintf("HOME=%s", global.GVA_CONFIG.Task.TmpPath))
//	cmd.Env = append(cmd.Env, fmt.Sprintf("PWD=%s", cmd.Dir))
//	cmd.Env = append(cmd.Env, "PYTHONUNBUFFERED=1")
//	cmd.Env = append(cmd.Env, "ANSIBLE_FORCE_COLOR=True")
//	if environmentVars != nil {
//		cmd.Env = append(cmd.Env, *environmentVars...)
//	}
//
//	return cmd
//}
//
//func (p TaskWithLogger) runCmd(command string, args []string) error {
//	cmd := p.makeCmd(command, args, nil)
//	p.Logger.LogCmd(cmd)
//	return cmd.Run()
//}
//
//func (p TaskWithLogger) RunPlaybook(args []string, environmentVars *[]string, cb func(*os.Process)) error {
//	cmd := p.makeCmd("ansible-playbook", args, environmentVars)
//	p.Logger.LogCmd(cmd)
//	cmd.Stdin = strings.NewReader("")
//	err := cmd.Start()
//	if err != nil {
//		return err
//	}
//	cb(cmd.Process)
//	return cmd.Wait()
//}
//
//func (p TaskWithLogger) RunGalaxy(args []string) error {
//	return p.runCmd("ansible-galaxy", args)
//}
//
//type Logger interface {
//	Log(msg string)
//	LogCmd(cmd *exec.Cmd)
//	LogSsh(channel *ssh.Channel)
//}
