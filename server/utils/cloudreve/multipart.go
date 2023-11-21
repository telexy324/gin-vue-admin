package cloudreve

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"time"
)

type RespSessionStruct struct {
	Code int     `json:"code"`
	Data Session `json:"data"`
	Msg  string  `json:"msg"`
}
type Session struct {
	SessionID string `json:"sessionID"`
	ChunkSize int    `json:"chunkSize"`
	Expires   int    `json:"expires"`
}

func (c *CloudreveClient) Upload(file *multipart.FileHeader) (err error) {
	body, _ := json.Marshal(map[string]interface{}{
		"path":          "/",
		"size":          file.Size,
		"name":          file.Filename,
		"policy_id":     c.Policy,
		"last_modified": time.UnixMilli,
		//"mime_type":
	})
	req, err := http.NewRequest("POST", global.GVA_CONFIG.Cloudreve.Address+"/file/upload", bytes.NewReader(body))

	if err != nil {
		return
	}

	resp, err := c.HttpClient.Do(req)

	if err != nil {
		return
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		return fmt.Errorf("error http code %d", resp.StatusCode)
	}

	respBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return
	}

	respCommon := &RespSessionStruct{}
	if err = json.Unmarshal(respBody, respCommon); err != nil {
		return err
	}
	if respCommon.Code != 0 {
		return fmt.Errorf("error response code %d", respCommon.Code)
	}

	sessionId := respCommon.Data.SessionID
	reqUpload, err := http.NewRequest("POST", global.GVA_CONFIG.Cloudreve.Address+"/"+sessionId+"/0", bytes.NewReader(body))

	if err != nil {
		return
	}

	respUpload, err := c.HttpClient.Do(reqUpload)

	if err != nil {
		return
	}

	defer func() {
		_ = respUpload.Body.Close()
	}()

	if respUpload.StatusCode != 200 {
		return fmt.Errorf("error http code %d", respUpload.StatusCode)
	}

	respUploadBody, err := ioutil.ReadAll(respUpload.Body)

	if err != nil {
		return
	}

	respCommon = &RespSessionStruct{}
	if err = json.Unmarshal(respUploadBody, respCommon); err != nil {
		return err
	}
	if respCommon.Code != 0 {
		return fmt.Errorf("error response code %d", respCommon.Code)
	}
	return nil
}
