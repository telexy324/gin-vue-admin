package recordPool

import (
	"encoding/json"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/application"
	applicationReq "github.com/flipped-aurora/gin-vue-admin/server/model/application/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
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
	pageReq := request.PageInfo{}
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
	case "/cmdb/deleteServerByIds":
		if err = json.Unmarshal(req, &idsReq); err != nil {
			global.GVA_LOG.Error("get delete servers req failed", zap.Any("err", err))
			return
		}
		if err = json.Unmarshal(resp, &commonResp); err != nil {
			global.GVA_LOG.Error("get delete servers resp failed", zap.Any("err", err))
			return
		}
		l.Detail = common + "批量删除服务器 " + "IDs " + idsReq.Ids
	case "/cmdb/updateServer":
		l.Detail = common + "登陆"
	case "/cmdb/addSystem":
		l.Detail = common + "登陆"
	case "/cmdb/deleteSystem":
		l.Detail = common + "登陆"
	case "/cmdb/deleteSystemByIds":
		l.Detail = common + "登陆"
	case "/cmdb/updateSystem":
		l.Detail = common + "登陆"
	case "/cmdb/system/addEditRelation":
		l.Detail = common + "登陆"
	case "/cmdb/system/deleteEditRelation":
		l.Detail = common + "登陆"
	case "/cmdb/system/updateEditRelation":
		l.Detail = common + "登陆"
	case "/task/addTask":
		l.Detail = common + "登陆"
	case "/task/deleteTask":
		l.Detail = common + "登陆"
	case "/task/stopTask":
		l.Detail = common + "登陆"
	case "/task/template/addSet":
		l.Detail = common + "登陆"
	case "/task/template/addSetTask":
		l.Detail = common + "登陆"
	case "/task/template/addTemplate":
		l.Detail = common + "登陆"
	case "/task/template/deleteSet":
		l.Detail = common + "登陆"
	case "/task/template/deleteSetByIds":
		l.Detail = common + "登陆"
	case "/task/template/deleteTemplate":
		l.Detail = common + "登陆"
	case "/task/template/deleteTemplateByIds":
		l.Detail = common + "登陆"
	case "/task/template/processSetTask":
		l.Detail = common + "登陆"
	case "/task/template/updateSet":
		l.Detail = common + "登陆"
	case "/task/template/updateTemplate":
		l.Detail = common + "登陆"
	case "/task/schedule/deleteSchedule":
		l.Detail = common + "登陆"
	case "/task/schedule/deleteScheduleByIds":
		l.Detail = common + "登陆"
	case "/task/schedule/updateSchedule":
		l.Detail = common + "登陆"
	case "/task/template/downloadFile":
		l.Detail = common + "登陆"
	case "/ssh/run":
		l.Detail = common + "登陆"
	case "/logUpload/addServer":
		l.Detail = common + "登陆"
	case "/logUpload/deleteServer":
		l.Detail = common + "登陆"
	case "/logUpload/deleteServerByIds":
		l.Detail = common + "登陆"
	case "/logUpload/updateServer":
		l.Detail = common + "登陆"
	case "/logUpload/addSecret":
		l.Detail = common + "登陆"
	case "/logUpload/deleteSecret":
		l.Detail = common + "登陆"
	case "/logUpload/deleteSecretByIds":
		l.Detail = common + "登陆"
	case "/logUpload/updateSecret":
		l.Detail = common + "登陆"
	case "/task/template/uploadLogServer":
		l.Detail = common + "登陆"
	case "/task/template/deployServer":
		l.Detail = common + "登陆"
	}

	r.logger <- l
}
