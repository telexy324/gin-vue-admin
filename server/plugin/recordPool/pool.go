package recordPool

import (
	"encoding/json"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/application"
	applicationReq "github.com/flipped-aurora/gin-vue-admin/server/model/application/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/logUploadMdl"
	"github.com/flipped-aurora/gin-vue-admin/server/model/scheduleMdl"
	"github.com/flipped-aurora/gin-vue-admin/server/model/taskMdl"
	taskReq "github.com/flipped-aurora/gin-vue-admin/server/model/taskMdl/request"
	taskRes "github.com/flipped-aurora/gin-vue-admin/server/model/taskMdl/response"
	"go.uber.org/zap"
	"strconv"
	"time"
)

var RPool RecordPool

type RecordPool struct {
	// logger channel used to putting log records to database.
	logger chan application.ApplicationRecord
}

func CreateRecordPool() RecordPool {
	RPool = RecordPool{
		logger: make(chan application.ApplicationRecord, 10000), // store log records to database
	}
	go RPool.Run()
	return RPool
}

func (p *RecordPool) Run() {
	ticker := time.NewTicker(1 * time.Second)

	defer func() {
		ticker.Stop()
	}()

	for {
		select {
		case <-ticker.C: // timer 1 seconds
			if len(p.logger) == 0 {
				continue
			}

			record := <-p.logger // new log message which should be put to database
			err := applicationService.CreateApplicationRecord(record)

			if err != nil {
				global.GVA_LOG.Error(err.Error())
			}
		}
	}
}

func (r *RecordPool) AddRecord(userID int, ip, action string, req, resp []byte) {
	now := time.Now()

	l := application.ApplicationRecord{
		Ip:      ip,
		Action:  action,
		UserID:  userID,
		LogTime: now,
	}

	err, user := systemService.UserService.FindUserById(userID)
	if err != nil {
		global.GVA_LOG.Error("get user info failed", zap.Any("err", err))
		return
	}

	common := "用户 " + user.Username + " "
	idReq := request.GetById{}
	idsReq := request.IdsReq{}
	//pageReq := request.PageInfo{}
	commonResp := response.Response{}

	switch action {
	case "/base/login":
		l.Detail = common + "登陆"
	case "/cmdb/addServer":
		addServerReq := applicationReq.AddServer{}
		if err = json.Unmarshal(req, &addServerReq); err != nil {
			global.GVA_LOG.Error("get add server req failed", zap.Any("err", err))
			return
		}
		if err = json.Unmarshal(resp, &commonResp); err != nil {
			global.GVA_LOG.Error("get add server resp failed", zap.Any("err", err))
			return
		}
		l.Detail = common + "增加服务器 " + addServerReq.Hostname + " IP " + addServerReq.ManageIp
		l.Status = commonResp.Code
		l.ErrorMessage = commonResp.Msg
	case "/cmdb/deleteServer":
		if err = json.Unmarshal(req, &idReq); err != nil {
			global.GVA_LOG.Error("get delete server req failed", zap.Any("err", err))
			return
		}
		if err = json.Unmarshal(resp, &commonResp); err != nil {
			global.GVA_LOG.Error("get delete server resp failed", zap.Any("err", err))
			return
		}
		l.Detail = common + "删除服务器 " + "ID " + strconv.Itoa(int(idReq.ID))
		l.Status = commonResp.Code
		l.ErrorMessage = commonResp.Msg
	case "/cmdb/deleteServerByIds":
		if err = json.Unmarshal(req, &idsReq); err != nil {
			global.GVA_LOG.Error("get delete servers req failed", zap.Any("err", err))
			return
		}
		if err = json.Unmarshal(resp, &commonResp); err != nil {
			global.GVA_LOG.Error("get delete servers resp failed", zap.Any("err", err))
			return
		}
		ids, err := json.Marshal(idsReq.Ids)
		if err != nil {
			global.GVA_LOG.Error("get ids req failed", zap.Any("err", err))
			return
		}
		l.Detail = common + "批量删除服务器 " + "IDs " + string(ids)
		l.Status = commonResp.Code
		l.ErrorMessage = commonResp.Msg
	case "/cmdb/updateServer":
		updateServerReq := applicationReq.UpdateServer{}
		if err = json.Unmarshal(req, &updateServerReq); err != nil {
			global.GVA_LOG.Error("get update server req failed", zap.Any("err", err))
			return
		}
		if err = json.Unmarshal(resp, &commonResp); err != nil {
			global.GVA_LOG.Error("get update server resp failed", zap.Any("err", err))
			return
		}
		l.Detail = common + "修改服务器 " + "ID " + strconv.Itoa(int(idReq.ID))
		l.Status = commonResp.Code
		l.ErrorMessage = commonResp.Msg
	case "/cmdb/addSystem":
		addSystemReq := applicationReq.AddSystem{}
		if err = json.Unmarshal(req, &addSystemReq); err != nil {
			global.GVA_LOG.Error("get add system req failed", zap.Any("err", err))
			return
		}
		if err = json.Unmarshal(resp, &commonResp); err != nil {
			global.GVA_LOG.Error("get add system resp failed", zap.Any("err", err))
			return
		}
		l.Detail = common + "增加系统 " + addSystemReq.Name
		l.Status = commonResp.Code
		l.ErrorMessage = commonResp.Msg
	case "/cmdb/deleteSystem":
		if err = json.Unmarshal(req, &idReq); err != nil {
			global.GVA_LOG.Error("get delete system req failed", zap.Any("err", err))
			return
		}
		if err = json.Unmarshal(resp, &commonResp); err != nil {
			global.GVA_LOG.Error("get delete system resp failed", zap.Any("err", err))
			return
		}
		l.Detail = common + "删除系统 " + "ID " + strconv.Itoa(int(idReq.ID))
		l.Status = commonResp.Code
		l.ErrorMessage = commonResp.Msg
	case "/cmdb/deleteSystemByIds":
		if err = json.Unmarshal(req, &idsReq); err != nil {
			global.GVA_LOG.Error("get delete systems req failed", zap.Any("err", err))
			return
		}
		if err = json.Unmarshal(resp, &commonResp); err != nil {
			global.GVA_LOG.Error("get delete systems resp failed", zap.Any("err", err))
			return
		}
		ids, err := json.Marshal(idsReq.Ids)
		if err != nil {
			global.GVA_LOG.Error("get ids req failed", zap.Any("err", err))
			return
		}
		l.Detail = common + "批量删除系统 " + "IDs " + string(ids)
		l.Status = commonResp.Code
		l.ErrorMessage = commonResp.Msg
	case "/cmdb/updateSystem":
		updateSystemReq := applicationReq.AddSystem{}
		if err = json.Unmarshal(req, &updateSystemReq); err != nil {
			global.GVA_LOG.Error("get update system req failed", zap.Any("err", err))
			return
		}
		if err = json.Unmarshal(resp, &commonResp); err != nil {
			global.GVA_LOG.Error("get update system resp failed", zap.Any("err", err))
			return
		}
		l.Detail = common + "修改系统 " + "ID " + strconv.Itoa(int(idReq.ID))
		l.Status = commonResp.Code
		l.ErrorMessage = commonResp.Msg
	case "/cmdb/system/addEditRelation":
		addSystemEditRelationReq := application.ApplicationSystemEditRelation{}
		if err = json.Unmarshal(req, &addSystemEditRelationReq); err != nil {
			global.GVA_LOG.Error("get add system relation req failed", zap.Any("err", err))
			return
		}
		if err = json.Unmarshal(resp, &commonResp); err != nil {
			global.GVA_LOG.Error("get add system relation resp failed", zap.Any("err", err))
			return
		}
		l.Detail = common + "增加系统关系图, 系统ID " + strconv.Itoa(addSystemEditRelationReq.SystemId)
		l.Status = commonResp.Code
		l.ErrorMessage = commonResp.Msg
	case "/cmdb/system/deleteEditRelation":
		if err = json.Unmarshal(req, &idReq); err != nil {
			global.GVA_LOG.Error("get delete system relation req failed", zap.Any("err", err))
			return
		}
		if err = json.Unmarshal(resp, &commonResp); err != nil {
			global.GVA_LOG.Error("get delete system relation resp failed", zap.Any("err", err))
			return
		}
		l.Detail = common + "删除系统关系图 " + "系统ID " + strconv.Itoa(int(idReq.ID))
		l.Status = commonResp.Code
		l.ErrorMessage = commonResp.Msg
	case "/cmdb/system/updateEditRelation":
		updateSystemEditRelationReq := application.ApplicationSystemEditRelation{}
		if err = json.Unmarshal(req, &updateSystemEditRelationReq); err != nil {
			global.GVA_LOG.Error("get update system relation req failed", zap.Any("err", err))
			return
		}
		if err = json.Unmarshal(resp, &commonResp); err != nil {
			global.GVA_LOG.Error("get update system relation resp failed", zap.Any("err", err))
			return
		}
		l.Detail = common + "修改系统关系图 " + "系统ID " + strconv.Itoa(updateSystemEditRelationReq.SystemId)
		l.Status = commonResp.Code
		l.ErrorMessage = commonResp.Msg
	case "/task/addTask":
		addTaskReq := taskMdl.Task{}
		if err = json.Unmarshal(req, &addTaskReq); err != nil {
			global.GVA_LOG.Error("get add task req failed", zap.Any("err", err))
			return
		}
		if err = json.Unmarshal(resp, &commonResp); err != nil {
			global.GVA_LOG.Error("get add task resp failed", zap.Any("err", err))
			return
		}
		l.Detail = common + "新增任务 " + "系统ID " + strconv.Itoa(addTaskReq.SystemId) + "模板ID " + strconv.Itoa(addTaskReq.TemplateId)
		l.Status = commonResp.Code
		l.ErrorMessage = commonResp.Msg
	case "/task/deleteTask":
		if err = json.Unmarshal(req, &idReq); err != nil {
			global.GVA_LOG.Error("get delete task req failed", zap.Any("err", err))
			return
		}
		if err = json.Unmarshal(resp, &commonResp); err != nil {
			global.GVA_LOG.Error("get delete task resp failed", zap.Any("err", err))
			return
		}
		l.Detail = common + "删除任务 " + "ID " + strconv.Itoa(int(idReq.ID))
		l.Status = commonResp.Code
		l.ErrorMessage = commonResp.Msg
	case "/task/stopTask":
		if err = json.Unmarshal(req, &idReq); err != nil {
			global.GVA_LOG.Error("get stop task req failed", zap.Any("err", err))
			return
		}
		stopTaskResp := taskRes.StopTaskResponse{}
		if err = json.Unmarshal(resp, &stopTaskResp); err != nil {
			global.GVA_LOG.Error("get stop task resp failed", zap.Any("err", err))
			return
		}
		if len(stopTaskResp.FailedIps) > 0 {
			ips, err := json.Marshal(stopTaskResp.FailedIps)
			if err != nil {
				global.GVA_LOG.Error("get stop task failed ips failed", zap.Any("err", err))
				return
			}
			l.Detail = common + "停止任务 " + "ID " + strconv.Itoa(int(idReq.ID)) + "失败节点 " + string(ips)
		} else {
			l.Detail = common + "停止任务 " + "ID " + strconv.Itoa(int(idReq.ID))
		}
		l.Status = commonResp.Code
		l.ErrorMessage = commonResp.Msg
	case "/task/template/addSet":
		addSetReq := taskReq.AddSet{}
		if err = json.Unmarshal(req, &addSetReq); err != nil {
			global.GVA_LOG.Error("get add set req failed", zap.Any("err", err))
			return
		}
		if err = json.Unmarshal(resp, &commonResp); err != nil {
			global.GVA_LOG.Error("get add set resp failed", zap.Any("err", err))
			return
		}
		l.Detail = common + "增加模板集 " + addSetReq.Name + "系统ID " + strconv.Itoa(addSetReq.SystemId)
		l.Status = commonResp.Code
		l.ErrorMessage = commonResp.Msg
	case "/task/template/addSetTask":
		addSetTaskReq := taskMdl.SetTask{}
		if err = json.Unmarshal(req, &addSetTaskReq); err != nil {
			global.GVA_LOG.Error("get add set task req failed", zap.Any("err", err))
			return
		}
		if err = json.Unmarshal(resp, &commonResp); err != nil {
			global.GVA_LOG.Error("get add set task resp failed", zap.Any("err", err))
			return
		}
		l.Detail = common + "增加模板集任务集 " + "模板集ID " + strconv.Itoa(addSetTaskReq.SetId)
		l.Status = commonResp.Code
		l.ErrorMessage = commonResp.Msg
	case "/task/template/addTemplate":
		addTemplateReq := taskMdl.TaskTemplate{}
		if err = json.Unmarshal(req, &addTemplateReq); err != nil {
			global.GVA_LOG.Error("get add template req failed", zap.Any("err", err))
			return
		}
		if err = json.Unmarshal(resp, &commonResp); err != nil {
			global.GVA_LOG.Error("get add template resp failed", zap.Any("err", err))
			return
		}
		l.Detail = common + "增加模板 " + addTemplateReq.Name + "系统ID " + strconv.Itoa(addTemplateReq.SystemId)
		l.Status = commonResp.Code
		l.ErrorMessage = commonResp.Msg
	case "/task/template/deleteSet":
		if err = json.Unmarshal(req, &idReq); err != nil {
			global.GVA_LOG.Error("get delete set req failed", zap.Any("err", err))
			return
		}
		if err = json.Unmarshal(resp, &commonResp); err != nil {
			global.GVA_LOG.Error("get delete set resp failed", zap.Any("err", err))
			return
		}
		l.Detail = common + "删除模板集 " + "ID " + strconv.Itoa(int(idReq.ID))
		l.Status = commonResp.Code
		l.ErrorMessage = commonResp.Msg
	case "/task/template/deleteSetByIds":
		if err = json.Unmarshal(req, &idsReq); err != nil {
			global.GVA_LOG.Error("get delete sets req failed", zap.Any("err", err))
			return
		}
		if err = json.Unmarshal(resp, &commonResp); err != nil {
			global.GVA_LOG.Error("get delete sets resp failed", zap.Any("err", err))
			return
		}
		ids, err := json.Marshal(idsReq.Ids)
		if err != nil {
			global.GVA_LOG.Error("get ids req failed", zap.Any("err", err))
			return
		}
		l.Detail = common + "批量删除模板集 " + "IDs " + string(ids)
		l.Status = commonResp.Code
		l.ErrorMessage = commonResp.Msg
	case "/task/template/deleteTemplate":
		if err = json.Unmarshal(req, &idReq); err != nil {
			global.GVA_LOG.Error("get delete template req failed", zap.Any("err", err))
			return
		}
		if err = json.Unmarshal(resp, &commonResp); err != nil {
			global.GVA_LOG.Error("get delete template resp failed", zap.Any("err", err))
			return
		}
		l.Detail = common + "删除模板 " + "ID " + strconv.Itoa(int(idReq.ID))
		l.Status = commonResp.Code
		l.ErrorMessage = commonResp.Msg
	case "/task/template/deleteTemplateByIds":
		if err = json.Unmarshal(req, &idsReq); err != nil {
			global.GVA_LOG.Error("get delete templates req failed", zap.Any("err", err))
			return
		}
		if err = json.Unmarshal(resp, &commonResp); err != nil {
			global.GVA_LOG.Error("get delete templates resp failed", zap.Any("err", err))
			return
		}
		ids, err := json.Marshal(idsReq.Ids)
		if err != nil {
			global.GVA_LOG.Error("get ids req failed", zap.Any("err", err))
			return
		}
		l.Detail = common + "批量删除模板 " + "IDs " + string(ids)
		l.Status = commonResp.Code
		l.ErrorMessage = commonResp.Msg
	case "/task/template/processSetTask":
		processTaskReq := taskReq.ProcessTaskRequest{}
		if err = json.Unmarshal(req, &processTaskReq); err != nil {
			global.GVA_LOG.Error("get process task req failed", zap.Any("err", err))
			return
		}
		processTaskResp := taskRes.TaskResponse{}
		if err = json.Unmarshal(resp, &processTaskResp); err != nil {
			global.GVA_LOG.Error("get stop task resp failed", zap.Any("err", err))
			return
		}
		l.Detail = common + "执行模板集下一模板 " + "模板集ID " + strconv.Itoa(int(processTaskReq.ID)) + " 下一模板ID " + strconv.Itoa(processTaskResp.Task.TemplateId)
		l.Status = commonResp.Code
		l.ErrorMessage = commonResp.Msg
	case "/task/template/updateSet":
		updateSetReq := taskReq.AddSet{}
		if err = json.Unmarshal(req, &idReq); err != nil {
			global.GVA_LOG.Error("get update set req failed", zap.Any("err", err))
			return
		}
		if err = json.Unmarshal(resp, &commonResp); err != nil {
			global.GVA_LOG.Error("get update set resp failed", zap.Any("err", err))
			return
		}
		l.Detail = common + "修改模板集 " + "ID " + strconv.Itoa(int(updateSetReq.ID))
		l.Status = commonResp.Code
		l.ErrorMessage = commonResp.Msg
	case "/task/template/updateTemplate":
		updateTemplateReq := taskMdl.TaskTemplate{}
		if err = json.Unmarshal(req, &updateTemplateReq); err != nil {
			global.GVA_LOG.Error("get update template req failed", zap.Any("err", err))
			return
		}
		if err = json.Unmarshal(resp, &commonResp); err != nil {
			global.GVA_LOG.Error("get update template resp failed", zap.Any("err", err))
			return
		}
		l.Detail = common + "修改模板 " + "ID " + strconv.Itoa(int(updateTemplateReq.ID))
		l.Status = commonResp.Code
		l.ErrorMessage = commonResp.Msg
	case "/task/schedule/addSchedule":
		addScheduleReq := scheduleMdl.Schedule{}
		if err = json.Unmarshal(req, &addScheduleReq); err != nil {
			global.GVA_LOG.Error("get add schedule req failed", zap.Any("err", err))
			return
		}
		if err = json.Unmarshal(resp, &commonResp); err != nil {
			global.GVA_LOG.Error("get add schedule resp failed", zap.Any("err", err))
			return
		}
		l.Detail = common + "增加定时任务 " + "模板ID " + strconv.Itoa(addScheduleReq.TemplateID)
		l.Status = commonResp.Code
		l.ErrorMessage = commonResp.Msg
	case "/task/schedule/deleteSchedule":
		if err = json.Unmarshal(req, &idReq); err != nil {
			global.GVA_LOG.Error("get delete schedule req failed", zap.Any("err", err))
			return
		}
		if err = json.Unmarshal(resp, &commonResp); err != nil {
			global.GVA_LOG.Error("get delete schedule resp failed", zap.Any("err", err))
			return
		}
		l.Detail = common + "删除定时任务 " + "ID " + strconv.Itoa(int(idReq.ID))
		l.Status = commonResp.Code
		l.ErrorMessage = commonResp.Msg
	case "/task/schedule/deleteScheduleByIds":
		if err = json.Unmarshal(req, &idsReq); err != nil {
			global.GVA_LOG.Error("get delete schedules req failed", zap.Any("err", err))
			return
		}
		if err = json.Unmarshal(resp, &commonResp); err != nil {
			global.GVA_LOG.Error("get delete schedules resp failed", zap.Any("err", err))
			return
		}
		ids, err := json.Marshal(idsReq.Ids)
		if err != nil {
			global.GVA_LOG.Error("get ids req failed", zap.Any("err", err))
			return
		}
		l.Detail = common + "批量删除定时任务 " + "IDs " + string(ids)
		l.Status = commonResp.Code
		l.ErrorMessage = commonResp.Msg
	case "/task/schedule/updateSchedule":
		updateScheduleReq := scheduleMdl.Schedule{}
		if err = json.Unmarshal(req, &updateScheduleReq); err != nil {
			global.GVA_LOG.Error("get update schedule req failed", zap.Any("err", err))
			return
		}
		if err = json.Unmarshal(resp, &commonResp); err != nil {
			global.GVA_LOG.Error("get update schedule resp failed", zap.Any("err", err))
			return
		}
		l.Detail = common + "修改定时任务 " + "ID " + strconv.Itoa(int(updateScheduleReq.ID))
		l.Status = commonResp.Code
		l.ErrorMessage = commonResp.Msg
	case "/task/template/downloadFile":
		downLoadFileReq := taskReq.DownLoadFileRequest{}
		if err = json.Unmarshal(req, &downLoadFileReq); err != nil {
			global.GVA_LOG.Error("get download file req failed", zap.Any("err", err))
			return
		}
		if err = json.Unmarshal(resp, &commonResp); err != nil {
			global.GVA_LOG.Error("get download file resp failed", zap.Any("err", err))
			return
		}
		l.Detail = common + "下载文件 " + " " + downLoadFileReq.File
		l.Status = commonResp.Code
		l.ErrorMessage = commonResp.Msg
	//case "/ssh/run":
	case "/logUpload/addServer":
		addServerReq := logUploadMdl.Server{}
		if err = json.Unmarshal(req, &addServerReq); err != nil {
			global.GVA_LOG.Error("get add log server req failed", zap.Any("err", err))
			return
		}
		if err = json.Unmarshal(resp, &commonResp); err != nil {
			global.GVA_LOG.Error("get add log server resp failed", zap.Any("err", err))
			return
		}
		l.Detail = common + "增加日志服务器 " + addServerReq.Hostname
		l.Status = commonResp.Code
		l.ErrorMessage = commonResp.Msg
	case "/logUpload/deleteServer":
		if err = json.Unmarshal(req, &idReq); err != nil {
			global.GVA_LOG.Error("get delete log server req failed", zap.Any("err", err))
			return
		}
		if err = json.Unmarshal(resp, &commonResp); err != nil {
			global.GVA_LOG.Error("get delete log server resp failed", zap.Any("err", err))
			return
		}
		l.Detail = common + "删除日志服务器 " + "ID " + strconv.Itoa(int(idReq.ID))
		l.Status = commonResp.Code
		l.ErrorMessage = commonResp.Msg
	case "/logUpload/deleteServerByIds":
		if err = json.Unmarshal(req, &idsReq); err != nil {
			global.GVA_LOG.Error("get delete log servers req failed", zap.Any("err", err))
			return
		}
		if err = json.Unmarshal(resp, &commonResp); err != nil {
			global.GVA_LOG.Error("get delete log servers resp failed", zap.Any("err", err))
			return
		}
		ids, err := json.Marshal(idsReq.Ids)
		if err != nil {
			global.GVA_LOG.Error("get ids req failed", zap.Any("err", err))
			return
		}
		l.Detail = common + "批量删除日志服务器 " + "IDs " + string(ids)
		l.Status = commonResp.Code
		l.ErrorMessage = commonResp.Msg
	case "/logUpload/updateServer":
		updateServerReq := logUploadMdl.Server{}
		if err = json.Unmarshal(req, &updateServerReq); err != nil {
			global.GVA_LOG.Error("get update log server req failed", zap.Any("err", err))
			return
		}
		if err = json.Unmarshal(resp, &commonResp); err != nil {
			global.GVA_LOG.Error("get update log server resp failed", zap.Any("err", err))
			return
		}
		l.Detail = common + "修改日志服务器 " + "ID " + strconv.Itoa(int(updateServerReq.ID))
		l.Status = commonResp.Code
		l.ErrorMessage = commonResp.Msg
	case "/logUpload/addSecret":
		addSecretReq := logUploadMdl.Secret{}
		if err = json.Unmarshal(req, &addSecretReq); err != nil {
			global.GVA_LOG.Error("get add log secret req failed", zap.Any("err", err))
			return
		}
		if err = json.Unmarshal(resp, &commonResp); err != nil {
			global.GVA_LOG.Error("get add log secret resp failed", zap.Any("err", err))
			return
		}
		l.Detail = common + "增加日志服务器密钥 " + "日志服务器ID " + strconv.Itoa(addSecretReq.ServerId)
		l.Status = commonResp.Code
		l.ErrorMessage = commonResp.Msg
	case "/logUpload/deleteSecret":
		if err = json.Unmarshal(req, &idReq); err != nil {
			global.GVA_LOG.Error("get delete log secret req failed", zap.Any("err", err))
			return
		}
		if err = json.Unmarshal(resp, &commonResp); err != nil {
			global.GVA_LOG.Error("get delete log secret resp failed", zap.Any("err", err))
			return
		}
		l.Detail = common + "删除日志服务器密钥 " + "ID " + strconv.Itoa(int(idReq.ID))
		l.Status = commonResp.Code
		l.ErrorMessage = commonResp.Msg
	case "/logUpload/deleteSecretByIds":
		if err = json.Unmarshal(req, &idsReq); err != nil {
			global.GVA_LOG.Error("get delete log secrets req failed", zap.Any("err", err))
			return
		}
		if err = json.Unmarshal(resp, &commonResp); err != nil {
			global.GVA_LOG.Error("get delete log secrets resp failed", zap.Any("err", err))
			return
		}
		ids, err := json.Marshal(idsReq.Ids)
		if err != nil {
			global.GVA_LOG.Error("get ids req failed", zap.Any("err", err))
			return
		}
		l.Detail = common + "批量删除日志服务器密钥 " + "IDs " + string(ids)
		l.Status = commonResp.Code
		l.ErrorMessage = commonResp.Msg
	case "/logUpload/updateSecret":
		updateSecretReq := logUploadMdl.Secret{}
		if err = json.Unmarshal(req, &updateSecretReq); err != nil {
			global.GVA_LOG.Error("get update log secret req failed", zap.Any("err", err))
			return
		}
		if err = json.Unmarshal(resp, &commonResp); err != nil {
			global.GVA_LOG.Error("get update log secret resp failed", zap.Any("err", err))
			return
		}
		l.Detail = common + "修改日志服务器密钥 " + "ID " + strconv.Itoa(int(updateSecretReq.ID))
		l.Status = commonResp.Code
		l.ErrorMessage = commonResp.Msg
	case "/task/template/uploadLogServer":
		uploadLogServerReq := taskReq.DownLoadFileRequest{}
		if err = json.Unmarshal(req, &uploadLogServerReq); err != nil {
			global.GVA_LOG.Error("get upload log server req failed", zap.Any("err", err))
			return
		}
		if err = json.Unmarshal(resp, &commonResp); err != nil {
			global.GVA_LOG.Error("get upload log server resp failed", zap.Any("err", err))
			return
		}
		l.Detail = common + "上传日志服务器 " + " " + uploadLogServerReq.File
		l.Status = commonResp.Code
		l.ErrorMessage = commonResp.Msg
	case "/task/template/deployServer":
		if err = json.Unmarshal(req, &idReq); err != nil {
			global.GVA_LOG.Error("get deploy server req failed", zap.Any("err", err))
			return
		}
		if err = json.Unmarshal(resp, &commonResp); err != nil {
			global.GVA_LOG.Error("get deploy server resp failed", zap.Any("err", err))
			return
		}
		l.Detail = common + "上传文件 " + "模板ID " + strconv.Itoa(int(idReq.ID))
		l.Status = commonResp.Code
		l.ErrorMessage = commonResp.Msg
	default:
		return
	}

	r.logger <- l
}
