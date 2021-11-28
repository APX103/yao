package helper

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

// Response 请求响应结果
type Response struct {
	Status  int                    `json:"status"`
	Body    string                 `json:"body"`
	Data    interface{}            `json:"data"`
	Headers map[string]interface{} `json:"headers"`
}

// RequestGet 发送GET请求
func RequestGet(url string, params map[string]interface{}, headers map[string]string) Response {
	return RequestSend("GET", url, params, nil, headers)
}

// RequestPost 发送POST请求
func RequestPost(url string, data interface{}, headers map[string]string) Response {
	return RequestSend("POST", url, map[string]interface{}{}, data, headers)
}

// RequestSend 发送Request请求
func RequestSend(method string, url string, params map[string]interface{}, data interface{}, headers map[string]string) Response {

	var body []byte
	var err error
	if data != nil {
		body, err = jsoniter.Marshal(data)
		if err != nil {
			return Response{
				Status: 500,
				Body:   err.Error(),
				Data:   map[string]interface{}{"code": 500, "message": err.Error()},
				Headers: map[string]interface{}{
					"Content-Type": "application/json;charset=utf8",
				},
			}
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return Response{
			Status: 500,
			Body:   err.Error(),
			Data:   map[string]interface{}{"code": 500, "message": err.Error()},
			Headers: map[string]interface{}{
				"Content-Type": "application/json;charset=utf8",
			},
		}
	}

	// Request Header
	for name, header := range headers {
		req.Header.Set(name, header)
	}

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return Response{
			Status: 0,
			Body:   err.Error(),
			Data:   map[string]interface{}{"code": 500, "message": err.Error()},
			Headers: map[string]interface{}{
				"Content-Type": "application/json;charset=utf8",
			},
		}
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body) // response body is []byte
	if err != nil {
		return Response{
			Status: 500,
			Body:   err.Error(),
			Data:   map[string]interface{}{"code": resp.StatusCode, "message": err.Error()},
			Headers: map[string]interface{}{
				"Content-Type": "application/json;charset=utf8",
			},
		}
	}

	// JSON 解析
	var res interface{}
	if strings.HasPrefix(resp.Header.Get("Content-Type"), "application/json") {
		err = jsoniter.Unmarshal(body, &res)
		if err != nil {
			return Response{
				Status: 500,
				Body:   err.Error(),
				Data:   map[string]interface{}{"code": resp.StatusCode, "message": err.Error()},
				Headers: map[string]interface{}{
					"Content-Type": "application/json;charset=utf8",
				},
			}
		}
	}
	respHeaders := map[string]interface{}{}
	for name := range resp.Header {
		respHeaders[name] = resp.Header.Get(name)
	}
	return Response{
		Status:  resp.StatusCode,
		Body:    string(body),
		Data:    res,
		Headers: respHeaders,
	}
}
