package cloudreve

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"io"
	"net/http"
	"strings"
)

func (c *CloudreveClient) Download(filePath string) (err error, file io.ReadCloser) {
	filePath = strings.Trim(filePath, " ")
	lastIndex := strings.LastIndex(filePath, "/")
	if len(filePath) == 0 || lastIndex == len(filePath) {
		return errors.New("invalid file path"), nil
	}
	fileName := filePath[lastIndex+1:]
	path := filePath[:lastIndex+1]
	pathUrl := "/directory%2F"
	if lastIndex != -1 {
		pathUrl = "/directory" + strings.ReplaceAll(path, "/", "%2F")
	}
	reqPolicy, err := http.NewRequest("GET", global.GVA_CONFIG.Cloudreve.Address+pathUrl, bytes.NewReader([]byte{}))

	if err != nil {
		return err, nil
	}

	respPolicy, err := c.HttpClient.Do(reqPolicy)

	if err != nil {
		return err, nil
	}

	defer func() {
		_ = respPolicy.Body.Close()
	}()

	if respPolicy.StatusCode != 200 {
		return fmt.Errorf("error http code %d", respPolicy.StatusCode), nil
	}

	respPolicyBody, err := io.ReadAll(respPolicy.Body)

	if err != nil {
		return err, nil
	}

	respDirectoryStruct := &RespStruct{
		Data: &Directory{},
	}

	if err = json.Unmarshal(respPolicyBody, &respDirectoryStruct); err != nil {
		return err, nil
	}
	if respDirectoryStruct.Code != 0 {
		return fmt.Errorf("error response code %d", respDirectoryStruct.Code), nil
	}
	directory, _ := respDirectoryStruct.Data.(*Directory)

	objectID := ""
	for _, object := range directory.Objects {
		if object.Name == fileName {
			objectID = object.ID
			break
		}
	}
	if objectID == "" {
		return errors.New("file not found"), nil
	}
	req, err := http.NewRequest("PUT", global.GVA_CONFIG.Cloudreve.Address+"/file/download/"+objectID, bytes.NewReader([]byte{}))

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
		return fmt.Errorf("error http code %d", resp.StatusCode), nil
	}

	respBody, err := io.ReadAll(resp.Body)

	if err != nil {
		return
	}

	respUrl := &RespStruct{}
	if err = json.Unmarshal(respBody, respUrl); err != nil {
		return err, nil
	}
	if respUrl.Code != 0 {
		return fmt.Errorf("error response code %d", respUrl.Code), nil
	}

	url, _ := respUrl.Data.(string)
	reqDownload, e := http.NewRequest("GET", strings.ReplaceAll(global.GVA_CONFIG.Cloudreve.Address, "/api/v3", "")+url, bytes.NewReader([]byte{}))
	if e != nil {
		return e, nil
	}
	respDownload, e := c.HttpClient.Do(reqDownload)

	if e != nil {
		return e, nil
	}
	//_ = bodyWriter.Close()

	if respDownload.StatusCode != 200 {
		return fmt.Errorf("error http code %d", respDownload.StatusCode), nil
	}
	return nil, respDownload.Body
}
