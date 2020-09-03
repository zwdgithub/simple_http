package simple_http

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// 默认超时时间
const defaultTimeout = time.Second * 15

type HttpUtil struct {
	req    *http.Request
	resp   *http.Response
	client *http.Client
	do     bool
	err    error
}

func NewHttpUtil() *HttpUtil {
	return &HttpUtil{
		client: &http.Client{
			Timeout: defaultTimeout,
		},
	}
}

func (h *HttpUtil) defaultClient() {
	h.client = &http.Client{
		Timeout: defaultTimeout,
	}
}

// error
func (h *HttpUtil) Error() error {
	return h.err
}

// 获取 response
func (h *HttpUtil) Response() *http.Response {
	return h.resp
}

// 构建 get request
func (h *HttpUtil) Get(url string, params url.Values) *HttpUtil {
	if h.err != nil {
		return h
	}
	url, h.err = BuildUrl(url, params)
	if h.err != nil {
		return h
	}
	h.req, h.err = http.NewRequest(http.MethodGet, url, nil)
	return h
}

// 构建 post request
func (h *HttpUtil) Post(url string, params interface{}) *HttpUtil {
	if h.err != nil {
		return h
	}
	b, err := json.Marshal(params)
	if err != nil {
		h.err = err
		return h
	}
	h.req, h.err = http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	return h
}

// 执行 http 请求
func (h *HttpUtil) Do() *HttpUtil {
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
func (h *HttpUtil) Result() ([]byte, error) {
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
func (h *HttpUtil) RContent() (string, error) {
	b, err := h.Result()
	if err != nil {
		return "", h.err
	}
	return string(b), err
}

// 获取 response 转 map[string]interface
func (h *HttpUtil) RMap() (map[string]interface{}, error) {
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
func (h *HttpUtil) RMarshal(r interface{}) error {
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
func (h *HttpUtil) SetHeader(header http.Header) *HttpUtil {
	h.req.Header = header
	return h
}

// 自定义 http client 属性
func (h *HttpUtil) CustomClient(custom func(client *http.Client) *http.Client) *HttpUtil {
	h.client = custom(h.client)
	return h
}

// 自定义 request 属性
func (h *HttpUtil) CustomRequest(custom func(request *http.Request) *http.Request) *HttpUtil {
	h.req = custom(h.req)
	return h
}
