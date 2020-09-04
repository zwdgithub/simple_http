package simple_http

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// 默认超时时间
const defaultTimeout = time.Second * 15

type httpUtil struct {
	req    *http.Request
	resp   *http.Response
	client *http.Client
	do     bool
	err    error
}

func NewHttpUtil() *httpUtil {
	return &httpUtil{
		client: &http.Client{
			Timeout: defaultTimeout,
		},
	}
}

func (h *httpUtil) defaultClient() {
	h.client = &http.Client{
		Timeout: defaultTimeout,
	}
}

// error
func (h *httpUtil) Error() error {
	return h.err
}

// 获取 response
func (h *httpUtil) Response() *http.Response {
	return h.resp
}

// 构建 get request
func (h *httpUtil) Get(url string, params ...url.Values) *httpUtil {
	if h.err != nil {
		return h
	}
	if len(params) > 0 {
		url, h.err = BuildUrl(url, params[0])
		if h.err != nil {
			return h
		}
	}
	h.req, h.err = http.NewRequest(http.MethodGet, url, nil)
	return h
}

// 构建 post request
func (h *httpUtil) Post(url string, reader io.Reader) *httpUtil {
	if h.err != nil {
		return h
	}
	h.req, h.err = http.NewRequest(http.MethodPost, url, reader)
	return h
}

// post form
func (h *httpUtil) PostForm(url string, params url.Values) *httpUtil {
	if h.err != nil {
		return h
	}
	h.Post(url, strings.NewReader(params.Encode()))
	h.req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return h
}

// post json
func (h *httpUtil) PostJson(url string, params interface{}) *httpUtil {
	if h.err != nil {
		return h
	}
	b, err := json.Marshal(params)
	if err != nil {
		h.err = err
		return h
	}
	h.Post(url, bytes.NewReader(b))
	h.req.Header.Add("Content-Type", "application/json")
	return h
}

//TODO retry

// 执行 http 请求
func (h *httpUtil) Do() *httpUtil {
	if h.err != nil {
		return h
	}
	if h.client == nil {
		h.defaultClient()
	}
	h.resp, h.err = h.client.Do(h.req)
	h.do = true
	return h
}

// 获取 response []byte
func (h *httpUtil) Result() ([]byte, error) {
	if !h.do {
		h.Do()
	}
	if h.err != nil {
		return nil, h.err
	}
	defer func() {
		_ = h.resp.Body.Close()
	}()
	b, err := ioutil.ReadAll(h.resp.Body)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// 获取 response string
func (h *httpUtil) RContent() (string, error) {
	b, err := h.Result()
	if err != nil {
		return "", h.err
	}
	return string(b), err
}

// 获取 response 转 map[string]interface
func (h *httpUtil) RMap() (map[string]interface{}, error) {
	b, err := h.Result()
	if err != nil {
		return nil, h.err
	}
	var r map[string]interface{}
	err = json.Unmarshal(b, &r)
	if err != nil {
		return nil, err
	}
	return r, err
}

// 获取 response , 并把值 marshal到r 中
func (h *httpUtil) RUnmarshal(r interface{}) error {
	b, err := h.Result()
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, &r)
	if err != nil {
		return err
	}
	return nil
}

// set request header
func (h *httpUtil) SetHeader(header map[string]string) *httpUtil {
	for k, v := range header {
		h.req.Header.Add(k, v)
	}
	return h
}

// 自定义 http client 属性
func (h *httpUtil) CustomClient(custom func(client *http.Client) *http.Client) *httpUtil {
	h.client = custom(h.client)
	return h
}

// 自定义 request 属性
func (h *httpUtil) CustomRequest(custom func(request *http.Request) *http.Request) *httpUtil {
	h.req = custom(h.req)
	return h
}
