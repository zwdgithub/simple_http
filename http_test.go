package main

import "testing"

func TestHttp(t *testing.T) {
	h := NewHttpUtil()
	content, err := h.Get("http://www.baidu.com", nil).Do().ResultContent()
	t.Log(content)
	t.Log(err)
}
