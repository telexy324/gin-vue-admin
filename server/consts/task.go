package consts

const (
	Command = iota + 1
	Script
)

const (
	ScheduleValid = iota + 1
)

const (
	ExecuteTypeNormal = iota + 1
	ExecuteTypeDownload
	ExecuteTypeDeploy
)

const (
	LogOutputTypeDirect = iota + 1
	LogOutputTypeUpload
	LogOutputTypeNetDisk
)

const (
	LogServerModeFtp = iota + 1
	LogServerModeSSH
)

const (
	ShellTypeSh = iota + 1
	ShellTypeBash
	ShellTypePython
)

const (
	TaskTemplateDiscoverServers         = 99999999
	TaskTemplateGatherServerInformation = 99999998
)

const (
	NonInteractive = iota
	Interactive
)

const Replacer = "${}"

const (
	DeployTypeFtpSftp = iota + 1
	DeployTypeNetDisk
)
