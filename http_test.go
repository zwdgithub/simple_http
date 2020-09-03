package simple_http

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"
)

func custom(client *http.Client) *http.Client {
	client.Timeout = time.Nanosecond * 1
	client.Transport = nil
	return client
}

func TestHttp(t *testing.T) {
	h := NewHttpUtil()
	content, err := h.Get("https://www.baidu.com", nil).
		CustomClient(func(client *http.Client) *http.Client {
			client.Timeout = 1111
			return client
		}).
		CustomClient(custom).
		RContent()

	t.Log(content)
	t.Log(err)
}

func TestCustomRequest(t *testing.T) {
	h := NewHttpUtil()
	ctx, cancel := context.WithCancel(context.Background())
	customReq := func(req *http.Request) *http.Request {
		req = req.WithContext(ctx)
		return req
	}
	go func() {
		time.Sleep(time.Millisecond * 10)
		cancel()
	}()
	content, err := h.Get("https://www.baidu.com", nil).
		CustomRequest(customReq).
		RContent()

	t.Log(content, err)

}

func TestResponseError(t *testing.T) {
	h := NewHttpUtil()
	h.Get("https://www.1.com", nil).Do()
	if h.Error() != nil {
		t.Log(h.err)
	}
	t.Log(h.Response())
}

/*
{
resultcode: "101",
reason: "错误的请求KEY",
result: null,
error_code: 10001
}
*/

type JuheResponse struct {
	ResultCode string      `json:"resultcode"`
	Reason     string      `json:"reason"`
	Result     interface{} `json:"result"`
	ErrorCode  int         `json:"error_code"`
}

func (j *JuheResponse) String() string {
	return fmt.Sprintf("resultcode: %s, reason: %s, result: %v, errorcode: %d",
		j.ResultCode, j.Reason, j.Result, j.ErrorCode)
}

func TestResponseMarshal(t *testing.T) {
	h := NewHttpUtil()
	var r *JuheResponse
	err := h.Get("http://apis.juhe.cn/ip/ipNew?ip=112.112.11.11&key=").Do().RUnmarshal(&r)
	t.Log(r)
	t.Log(err)

	h1 := NewHttpUtil()
	var m map[string]interface{}
	err = h1.Get("http://apis.juhe.cn/ip/ipNew?ip=112.112.11.11&key=").Do().RUnmarshal(&m)
	t.Log(r)
	t.Log(err)
}

func TestResponseMarshalMap(t *testing.T) {
	h := NewHttpUtil()
	r, err := h.Get("http://apis.juhe.cn/ip/ipNew?ip=112.112.11.11&key=").RMap()
	t.Log(r)
	t.Log(err)
}

func TestForm(t *testing.T) {
	h := NewHttpUtil()
	p := url.Values{}
	p.Set("word", "是")
	p.Set("key", "b641bfe5999d9d6e4523bb03945041a2")
	p.Set("dtype", "json")
	content, err := h.PostForm("http://v.juhe.cn/xhzd/query", p).Do().RContent()
	t.Log(content)
	t.Log(err)
}
