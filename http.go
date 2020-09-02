package main

import (
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

func (h *HttpUtil) Get(_url string, params url.Values) *HttpUtil {
	if h.err != nil {
		return h
	}
	u, err := url.Parse(_url)
	if err != nil {
		h.err = err
		return h
	}
	if params != nil && len(params) > 0 {
		u.RawQuery = params.Encode()
	}
	h.req, h.err = http.NewRequest(http.MethodGet, u.String(), nil)
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

func (h *HttpUtil) ResultContent() (string, error) {
	if h.err != nil {
		return "", h.err
	}
	if !h.do {
		h.Do()
	}
	defer func() {
		h.err = h.resp.Body.Close()
	}()
	b, err := ioutil.ReadAll(h.resp.Body)
	if err != nil {
		return "", err
	}
	return string(b), err
}

func (h *HttpUtil) ResultMap() (map[string]interface{}, error) {
	if h.err != nil {
		return nil, h.err
	}
	if !h.do {
		h.Do()
	}
	defer func() {
		h.err = h.resp.Body.Close()
	}()
	b, err := ioutil.ReadAll(h.resp.Body)
	if err != nil {
		return nil, err
	}
	var r map[string]interface{}
	err = json.Unmarshal(b, &r)
	if err != nil {
		return nil, err
	}
	return r, err
}

func (h *HttpUtil) SetHeader(header http.Header) *HttpUtil {
	h.req.Header = header
	return h
}

func (h *HttpUtil) SetClient(client *http.Client) *HttpUtil {
	h.client = client
	return h
}
