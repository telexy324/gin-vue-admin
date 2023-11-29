package cloudreve

//
//import (
//	"bytes"
//	"encoding/json"
//	"fmt"
//	"github.com/flipped-aurora/gin-vue-admin/server/global"
//	"io/ioutil"
//	"net/http"
//	"strings"
//)
//
//func (c *CloudreveClient) Download(filePath string) (err error) {
//	filePath = strings.Trim(filePath, " ")
//	lastIndex:=strings.LastIndex(filePath, "/")
//	path:=filePath.
//	reqPolicy, err := http.NewRequest("GET", global.GVA_CONFIG.Cloudreve.Address+"/directory%2F", bytes.NewReader([]byte{}))
//
//	if err != nil {
//		return err
//	}
//
//	respPolicy, err := c.HttpClient.Do(reqPolicy)
//
//	if err != nil {
//		return err
//	}
//
//	defer func() {
//		_ = respPolicy.Body.Close()
//	}()
//
//	if respPolicy.StatusCode != 200 {
//		return fmt.Errorf("error http code %d", respPolicy.StatusCode)
//	}
//
//	respPolicyBody, err := ioutil.ReadAll(respPolicy.Body)
//
//	if err != nil {
//		return err
//	}
//
//	respDirectoryStruct := &RespStruct{
//		Data: &Directory{},
//	}
//
//	if err = json.Unmarshal(respPolicyBody, &respDirectoryStruct); err != nil {
//		return err
//	}
//	if respDirectoryStruct.Code != 0 {
//		return fmt.Errorf("error response code %d", respDirectoryStruct.Code)
//	}
//	directory, _ := respDirectoryStruct.Data.(*Directory)
//	policy := directory.Policy.ID
//}
