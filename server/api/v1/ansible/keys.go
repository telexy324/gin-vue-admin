package ansible

import (
	log "github.com/Sirupsen/logrus"
	"github.com/ansible-semaphore/semaphore/api/helpers"
	"github.com/ansible-semaphore/semaphore/db"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/ansible"
	request2 "github.com/flipped-aurora/gin-vue-admin/server/model/ansible/request"
	ansibleRes "github.com/flipped-aurora/gin-vue-admin/server/model/ansible/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"

	"github.com/gorilla/context"
)

type KeysApi struct {
}

// KeyMiddleware ensures a key exists and loads it to the context
func KeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		project := context.Get(r, "project").(db.Project)
		keyID, err := helpers.GetIntParam("key_id", w, r)
		if err != nil {
			return
		}

		key, err := helpers.Store(r).GetAccessKey(project.ID, keyID)

		if err != nil {
			helpers.WriteError(w, err)
			return
		}

		context.Set(r, "accessKey", key)
		next.ServeHTTP(w, r)
	})
}

func GetKeyRefs(w http.ResponseWriter, r *http.Request) {
	key := context.Get(r, "accessKey").(db.AccessKey)
	refs, err := helpers.Store(r).GetAccessKeyRefs(*key.ProjectID, key.ID)
	if err != nil {
		helpers.WriteError(w, err)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, refs)
}

// GetKeys retrieves sorted keys from the database
func GetKeys(w http.ResponseWriter, r *http.Request) {
	if key := context.Get(r, "accessKey"); key != nil {
		k := key.(db.AccessKey)
		helpers.WriteJSON(w, http.StatusOK, k)
		return
	}

	project := context.Get(r, "project").(db.Project)
	var keys []db.AccessKey

	keys, err := helpers.Store(r).GetAccessKeys(project.ID, helpers.QueryParams(r.URL))

	if err != nil {
		helpers.WriteError(w, err)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, keys)
}

// AddKey adds a new key to the database
func AddKey(w http.ResponseWriter, r *http.Request) {
	project := context.Get(r, "project").(db.Project)
	var key db.AccessKey

	if !helpers.Bind(w, r, &key) {
		return
	}

	if key.ProjectID == nil || *key.ProjectID != project.ID {
		helpers.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Project ID in body and URL must be the same",
		})
		return
	}

	if err := key.Validate(true); err != nil {
		helpers.WriteJSON(w, http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
		return
	}

	newKey, err := helpers.Store(r).CreateAccessKey(key)

	if err != nil {
		helpers.WriteError(w, err)
		return
	}

	user := context.Get(r, "user").(*db.User)

	objType := db.EventKey

	desc := "Access Key " + key.Name + " created"
	_, err = helpers.Store(r).CreateEvent(db.Event{
		UserID:      &user.ID,
		ProjectID:   newKey.ProjectID,
		ObjectType:  &objType,
		ObjectID:    &newKey.ID,
		Description: &desc,
	})

	if err != nil {
		log.Error(err)
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdateKey updates key in database
// nolint: gocyclo
func UpdateKey(w http.ResponseWriter, r *http.Request) {
	var key db.AccessKey
	oldKey := context.Get(r, "accessKey").(db.AccessKey)

	if !helpers.Bind(w, r, &key) {
		return
	}

	repos, err := helpers.Store(r).GetRepositories(*key.ProjectID, db.RetrieveQueryParams{})
	if err != nil {
		helpers.WriteError(w, err)
		return
	}

	for _, repo := range repos {
		if repo.SSHKeyID != key.ID {
			continue
		}
		err = repo.ClearCache()
		if err != nil {
			helpers.WriteError(w, err)
			return
		}
	}

	err = helpers.Store(r).UpdateAccessKey(key)
	if err != nil {
		helpers.WriteError(w, err)
		return
	}

	user := context.Get(r, "user").(*db.User)

	desc := "Access Key " + key.Name + " updated"
	objType := db.EventKey

	_, err = helpers.Store(r).CreateEvent(db.Event{
		UserID:      &user.ID,
		ProjectID:   oldKey.ProjectID,
		Description: &desc,
		ObjectID:    &oldKey.ID,
		ObjectType:  &objType,
	})

	if err != nil {
		log.Error(err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RemoveKey deletes a key from the database
func RemoveKey(w http.ResponseWriter, r *http.Request) {
	key := context.Get(r, "accessKey").(db.AccessKey)

	var err error

	err = helpers.Store(r).DeleteAccessKey(*key.ProjectID, key.ID)
	if err == db.ErrInvalidOperation {
		helpers.WriteJSON(w, http.StatusBadRequest, map[string]interface{}{
			"error": "Access Key is in use by one or more templates",
			"inUse": true,
		})
		return
	}

	if err != nil {
		helpers.WriteError(w, err)
		return
	}

	user := context.Get(r, "user").(*db.User)

	desc := "Access Key " + key.Name + " deleted"

	_, err = helpers.Store(r).CreateEvent(db.Event{
		UserID:      &user.ID,
		ProjectID:   key.ProjectID,
		Description: &desc,
	})

	if err != nil {
		log.Error(err)
	}

	w.WriteHeader(http.StatusNoContent)
}

// @Tags Key
// @Summary 新增Key
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body ansible.Key true ""
// @Success 200 {string} string "{"success":true,"data":{},"msg":"添加成功"}"
// @Router /ansible/key/addKey [post]
func (a *KeysApi) AddKey(c *gin.Context) {
	var environment ansible.Environment
	if err := c.ShouldBindJSON(&environment); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(environment, utils.EnvironmentVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if _, err := environmentService.CreateEnvironment(environment); err != nil {
		global.GVA_LOG.Error("添加失败!", zap.Any("err", err))

		response.FailWithMessage("添加失败", c)
	} else {
		response.OkWithMessage("添加成功", c)
	}
}

// @Tags Key
// @Summary 删除Key
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.GetById true "KeyId"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"删除成功"}"
// @Router /ansible/key/deleteKey [post]
func (a *KeysApi) DeleteKey(c *gin.Context) {
	var environment request2.GetByProjectId
	if err := c.ShouldBindJSON(&environment); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(environment, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := environmentService.DeleteEnvironment(environment.ProjectId, environment.ID); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Any("err", err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// @Tags Key
// @Summary 更新Key
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body ansible.Key true "主机名, 架构, 管理ip, 系统, 系统版本"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"更新成功"}"
// @Router /ansible/key/updateKey [post]
func (a *KeysApi) UpdateKey(c *gin.Context) {
	var environment ansible.Environment
	if err := c.ShouldBindJSON(&environment); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(environment, utils.EnvironmentVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := environmentService.UpdateEnvironment(environment); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Any("err", err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// @Tags Key
// @Summary 根据id获取Key
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request2.GetByProjectId true "KeyId"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /ansible/key/getKeyById [post]
func (a *KeysApi) GetKeyById(c *gin.Context) {
	var idInfo request2.GetByProjectId
	if err := c.ShouldBindJSON(&idInfo); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(idInfo, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if environment, err := environmentService.GetEnvironment(idInfo.ProjectId, idInfo.ID); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(ansibleRes.EnvironmentResponse{
			Environment: environment,
		}, "获取成功", c)
	}
}

// @Tags Key
// @Summary 分页获取基础Key列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request2.GetByProjectId true "页码, 每页大小"
// @Success 200 {string} string "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /ansible/key/getKeyList[post]
func (a *KeysApi) GetKeyList(c *gin.Context) {
	var pageInfo request2.GetByProjectId
	if err := c.ShouldBindJSON(&pageInfo); err != nil {
		global.GVA_LOG.Info("error", zap.Any("err", err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err := utils.Verify(pageInfo.PageInfo, utils.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err, environments, total := environmentService.GetEnvironments(pageInfo); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Any("err", err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     environments,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}
