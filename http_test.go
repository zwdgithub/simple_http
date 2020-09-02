package simple_http

import (
	"net/http"
	"testing"
	"time"
)

func custom(client *http.Client) {
	client.Timeout = time.Nanosecond * 1
	client.Transport = nil
}

func TestHttp(t *testing.T) {
	h := NewHttpUtil()
	content, err := h.Get("https://www.baidu.com", nil).
		CustomClient(func(client *http.Client) {
			client.Timeout = 1111
		}).
		CustomClient(custom).
		Do().
		RContent()

	t.Log(content)
	t.Log(err)
}

func TestHttp1(t *testing.T) {
	t.Log("http1")
}
