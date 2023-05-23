package utils

var (
	IdVerify                  = Rules{"ID": {NotEmpty()}}
	ApiVerify                 = Rules{"Path": {NotEmpty()}, "Description": {NotEmpty()}, "ApiGroup": {NotEmpty()}, "Method": {NotEmpty()}}
	MenuVerify                = Rules{"Path": {NotEmpty()}, "ParentId": {NotEmpty()}, "Name": {NotEmpty()}, "Component": {NotEmpty()}, "Sort": {Ge("0")}}
	MenuMetaVerify            = Rules{"Title": {NotEmpty()}}
	LoginVerify               = Rules{"CaptchaId": {NotEmpty()}, "Captcha": {NotEmpty()}, "Username": {NotEmpty()}, "Password": {NotEmpty()}}
	RegisterVerify            = Rules{"Username": {NotEmpty()}, "NickName": {NotEmpty()}, "Password": {NotEmpty()}, "AuthorityId": {NotEmpty()}}
	PageInfoVerify            = Rules{"Page": {NotEmpty()}, "PageSize": {NotEmpty()}}
	CustomerVerify            = Rules{"CustomerName": {NotEmpty()}, "CustomerPhoneData": {NotEmpty()}}
	AutoCodeVerify            = Rules{"Abbreviation": {NotEmpty()}, "StructName": {NotEmpty()}, "PackageName": {NotEmpty()}, "Fields": {NotEmpty()}}
	AuthorityVerify           = Rules{"AuthorityId": {NotEmpty()}, "AuthorityName": {NotEmpty()}, "ParentId": {NotEmpty()}}
	AuthorityIdVerify         = Rules{"AuthorityId": {NotEmpty()}}
	OldAuthorityVerify        = Rules{"OldAuthorityId": {NotEmpty()}}
	ChangePasswordVerify      = Rules{"Username": {NotEmpty()}, "Password": {NotEmpty()}, "NewPassword": {NotEmpty()}}
	SetUserAuthorityVerify    = Rules{"AuthorityId": {NotEmpty()}}
	ServerVerify              = Rules{"Hostname": {NotEmpty()}}
	ServerRelationVerify      = Rules{"StartServerId": {NotEmpty()}, "EndServerId": {NotEmpty()}}
	SystemRelationVerify      = Rules{"StartServerId": {NotEmpty()}, "EndServerId": {NotEmpty()}}
	AdminVerify               = Rules{"Name": {NotEmpty()}, "Mobile": {NotEmpty()}, "DepartmentId": {NotEmpty()}}
	SystemVerify              = Rules{"Name": {NotEmpty()}}
	AppVerify                 = Rules{"Name": {NotEmpty()}}
	TaskTemplateVerify        = Rules{"Name": {NotEmpty()}}
	TaskVerify                = Rules{"Name": {NotEmpty()}}
	ScheduleVerify            = Rules{"TemplateId": {NotEmpty()}, "CronFormat": {NotEmpty()}}
	SystemEditRelationVerify  = Rules{"SystemId": {NotEmpty()}}
	TaskTemplateSetVerify     = Rules{"Name": {NotEmpty()}}
	TaskTemplateSetTaskVerify = Rules{"SetId": {NotEmpty()}}
	LogServerVerify           = Rules{"HostName": {NotEmpty()}, "ManageIp": {NotEmpty()}}
)
