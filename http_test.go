package simple_http

import (
	"context"
	"net/http"
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
