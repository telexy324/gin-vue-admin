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
)

const (
	LogOutputTypeDirect = iota + 1
	LogOutputTypeUpload
)

const (
	LogServerModeFtp = iota + 1
	LogServerModeSSH
)