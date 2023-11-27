package cloudreve

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"
)

type RespSessionStruct struct {
	Code int     `json:"code"`
	Data Session `json:"data"`
	Msg  string  `json:"msg"`
}
type Session struct {
	SessionID string `json:"sessionID"`
	ChunkSize int64  `json:"chunkSize"`
	Expires   int    `json:"expires"`
}

func (c *CloudreveClient) Upload(file io.Reader, fileName string, fileSize int64) (err error) {
	body, err := json.Marshal(map[string]interface{}{
		"path":          "/",
		"size":          fileSize,
		"name":          fileName,
		"policy_id":     c.Policy,
		"last_modified": time.Now().UnixMilli(),
		//"mime_type":
	})
	if err != nil {
		return
	}
	req, err := http.NewRequest("PUT", global.GVA_CONFIG.Cloudreve.Address+"/file/upload", bytes.NewReader(body))

	if err != nil {
		return
	}
	//req.Header.Add("Content-Type", "application/json")
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

	respSession := &RespSessionStruct{}
	if err = json.Unmarshal(respBody, respSession); err != nil {
		return err
	}
	if respSession.Code != 0 {
		return fmt.Errorf("error response code %d", respSession.Code)
	}

	sessionId := respSession.Data.SessionID

	var i int64
	for i = 0; i < fileSize/respSession.Data.ChunkSize+1; i++ {
		bodyBuffer := &bytes.Buffer{}
		bodyWriter := multipart.NewWriter(bodyBuffer)

		fileWriter, _ := bodyWriter.CreateFormFile("files", fileName)

		n := respSession.Data.ChunkSize
		if i == fileSize/respSession.Data.ChunkSize {
			n = fileSize % respSession.Data.ChunkSize
		}
		_, err = io.CopyN(fileWriter, file, n)
		if err != nil {
			return
		}
		//contentType := bodyWriter.FormDataContentType()
		_ = bodyWriter.Close()

		reqUpload, e := http.NewRequest("POST", global.GVA_CONFIG.Cloudreve.Address+"/file/upload/"+sessionId+"/"+strconv.Itoa(int(i)), bytes.NewReader(body))
		if e != nil {
			return e
		}

		reqUpload.Header.Add("Content-Length", strconv.Itoa(int(n)))
		respUpload, e := c.HttpClient.Do(reqUpload)

		if e != nil {
			return e
		}

		if respUpload.StatusCode != 200 {
			return fmt.Errorf("error http code %d", respUpload.StatusCode)
		}
		respUploadBody, e := ioutil.ReadAll(respUpload.Body)

		if e != nil {
			return e
		}

		respUploadStruct := &RespSessionStruct{}
		if err = json.Unmarshal(respUploadBody, respUploadStruct); err != nil {
			return err
		}
		if respUploadStruct.Code != 0 {
			return fmt.Errorf("error response code %d", respUploadStruct.Code)
		}
		_ = respUpload.Body.Close()
	}
	return nil
}
