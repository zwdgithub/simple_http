# simple_http
go http请求小工具, 简单支持链式调用  
每次使用需要重新调用 NewHttpUtil() 获取新的对象
```go
    h := NewHttpUtil()
    content, err := h.Get("https://www.baidu.com").Do().RContent()
```


