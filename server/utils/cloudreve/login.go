package cloudreve

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"io"
	"net/http"
	"net/http/cookiejar"
)

type RespStruct struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

type CloudreveClient struct {
	Username   string
	Password   string
	HttpClient *http.Client
	UserId     int
}

func NewCloudreveClient(username, password string) (c *CloudreveClient, err error) {
	cookie, _ := cookiejar.New(nil)
	httpClient := &http.Client{Jar: cookie}

	c = &CloudreveClient{
		Username:   username,
		Password:   password,
		HttpClient: httpClient,
	}
	body, _ := json.Marshal(map[string]string{"userName": username, "Password": password})
	reqLogin, err := http.NewRequest("POST", global.GVA_CONFIG.Cloudreve.Address+"/user/session", bytes.NewReader(body))

	if err != nil {
		return nil, err
	}

	respLogin, err := c.HttpClient.Do(reqLogin)

	if err != nil {
		return nil, err
	}

	defer func() {
		_ = respLogin.Body.Close()
	}()

	if respLogin.StatusCode != 200 {
		return nil, fmt.Errorf("error http code %d", respLogin.StatusCode)
	}

	respLoginBody, err := io.ReadAll(respLogin.Body)

	if err != nil {
		return nil, err
	}

	respCommon := &RespStruct{}
	if err = json.Unmarshal(respLoginBody, respCommon); err != nil {
		return nil, err
	}
	if respCommon.Code != 0 {
		return nil, fmt.Errorf("error response code %d", respCommon.Code)
	}

	return
}
