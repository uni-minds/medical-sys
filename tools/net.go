package tools

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type ConnectData struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"msg"`
}

func HttpGet(url string) (data ConnectData, text string, result []byte, err error) {
	tlscfg := tls.Config{InsecureSkipVerify: true}
	tr := &http.Transport{TLSClientConfig: &tlscfg}
	client := http.Client{Timeout: 5 * time.Second, Transport: tr}
	resp, err := client.Get(url)
	if err != nil {
		return data, text, result, err
	}
	defer resp.Body.Close()

	result, _ = ioutil.ReadAll(resp.Body)

	json.Unmarshal(result, &data)

	return data, string(result), result, nil
}

func HttpPost(url string, data interface{}, contentType string) (rdata ConnectData, err error) {
	jsonStr, _ := json.Marshal(data)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Add("content-type", contentType)
	if err != nil {
		return
	}
	defer req.Body.Close()

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(result, &rdata)
	return
}

func FailReturn(code int, msg interface{}) map[string]interface{} {
	var res = make(map[string]interface{})
	res["data"] = ""
	res["code"] = code
	res["msg"] = msg

	return res
}

// SuccessReturn api正确返回函数
func SuccessReturn(msg interface{}) map[string]interface{} {
	var res = make(map[string]interface{})
	res["data"] = msg
	res["code"] = http.StatusOK
	res["msg"] = "success"

	return res
}
