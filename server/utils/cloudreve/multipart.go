package cloudreve

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"io"
	"net/http"
	"strconv"
	"time"
)

type Session struct {
	SessionID string `json:"sessionID"`
	ChunkSize int64  `json:"chunkSize"`
	Expires   int    `json:"expires"`
}

type Directory struct {
	Objects []Object `json:"objects"`
	Policy  Policy   `json:"policy"`
}
type Policy struct {
	ID string `json:"id"`
}
type Object struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Path          string    `json:"path"`
	Thumb         bool      `json:"thumb"`
	Size          int       `json:"size"`
	Type          string    `json:"type"`
	Date          time.Time `json:"date"`
	CreateDate    time.Time `json:"create_date"`
	SourceEnabled bool      `json:"source_enabled"`
}

func (c *CloudreveClient) Upload(file io.Reader, fileName string, fileSize int64) (err error) {
	reqPolicy, err := http.NewRequest("GET", global.GVA_CONFIG.Cloudreve.Address+"/directory%2F", bytes.NewReader([]byte{}))

	if err != nil {
		return err
	}

	respPolicy, err := c.HttpClient.Do(reqPolicy)

	if err != nil {
		return err
	}

	defer func() {
		_ = respPolicy.Body.Close()
	}()

	if respPolicy.StatusCode != 200 {
		return fmt.Errorf("error http code %d", respPolicy.StatusCode)
	}

	respPolicyBody, err := io.ReadAll(respPolicy.Body)

	if err != nil {
		return err
	}

	respDirectoryStruct := &RespStruct{
		Data: &Directory{},
	}

	if err = json.Unmarshal(respPolicyBody, &respDirectoryStruct); err != nil {
		return err
	}
	if respDirectoryStruct.Code != 0 {
		return fmt.Errorf("error response code %d", respDirectoryStruct.Code)
	}
	directory, _ := respDirectoryStruct.Data.(*Directory)
	policy := directory.Policy.ID
	body, err := json.Marshal(map[string]interface{}{
		"path":          "/",
		"size":          fileSize,
		"name":          fileName,
		"policy_id":     policy,
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

	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		return
	}

	respSession := &RespStruct{
		Data: &Session{},
	}
	if err = json.Unmarshal(respBody, respSession); err != nil {
		return err
	}
	if respSession.Code != 0 {
		return fmt.Errorf("error response code %d", respSession.Code)
	}

	session, _ := respSession.Data.(*Session)
	sessionId := session.SessionID

	var i int64
	for i = 0; i < fileSize/session.ChunkSize+1; i++ {
		//bodyBuffer := &bytes.Buffer{}
		//bodyWriter := multipart.NewWriter(bodyBuffer)
		//
		//fileWriter, _ := bodyWriter.CreateFormFile("files", fileName)

		n := session.ChunkSize
		if i == fileSize/session.ChunkSize {
			n = fileSize % session.ChunkSize
		}
		bodyBuffer := &bytes.Buffer{}
		_, err = io.CopyN(bodyBuffer, file, n)
		if err != nil {
			return
		}
		//contentType := bodyWriter.FormDataContentType()

		reqUpload, e := http.NewRequest("POST", global.GVA_CONFIG.Cloudreve.Address+"/file/upload/"+sessionId+"/"+strconv.Itoa(int(i)), bodyBuffer)
		if e != nil {
			return e
		}

		//reqUpload.Header.Add("Content-Length", strconv.Itoa(int(n)))
		reqUpload.ContentLength = n
		respUpload, e := c.HttpClient.Do(reqUpload)

		if e != nil {
			return e
		}
		//_ = bodyWriter.Close()

		if respUpload.StatusCode != 200 {
			return fmt.Errorf("error http code %d", respUpload.StatusCode)
		}
		respUploadBody, e := io.ReadAll(respUpload.Body)

		if e != nil {
			return e
		}

		respUploadStruct := &RespStruct{}
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