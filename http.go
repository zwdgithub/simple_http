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

func (h *HttpUtil) initClient() {
	h.client = &http.Client{
		Timeout: defaultTimeout,
	}
}

func (h *HttpUtil) Error() error {
	return h.err
}

func (h *HttpUtil) Response() *http.Response {
	return h.resp
}

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

func (h *HttpUtil) Do() *HttpUtil {
	if h.err != nil {
		return h
	}
	if h.client == nil {
		h.initClient()
	}
	h.resp, h.err = h.client.Do(h.req)
	h.do = true
	return h
}

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

func (h *HttpUtil) RContent() (string, error) {
	b, err := h.Result()
	if err != nil {
		return "", h.err
	}
	return string(b), err
}

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

func (h *HttpUtil) SetHeader(header http.Header) *HttpUtil {
	h.req.Header = header
	return h
}

func (h *HttpUtil) CustomClient(custom func(client *http.Client)) *HttpUtil {
	custom(h.client)
	return h
}

func (h *HttpUtil) CustomRequest(custom func(request *http.Request)) *HttpUtil {
	custom(h.req)
	return h
}
