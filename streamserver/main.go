package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

//路由和   嵌套 限流结构体
type middleWareHandler struct {
	r *httprouter.Router
	l *ConnLimiter
}

//路由和最大限制连接数cc
func NewMiddleWareHandler(r *httprouter.Router, cc int) http.Handler {
	m := middleWareHandler{}
	m.r = r
	m.l = NewConnLimiter(cc)

	return m
}

func RegisterHandlers() *httprouter.Router {
	//创建路由
	router := httprouter.New()

	router.GET("/videos/:vid-id", streamHandler) //上传

	router.POST("/upload/:vid-id", uploadHandler) //播放

	router.GET("/testpage", testPageHandler) //上传页面

	return router
}

//
func (m middleWareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !m.l.GetConn() {
		sendErrorResponse(w, http.StatusTooManyRequests, "Too many requests")
		return
	}
	m.r.ServeHTTP(w, r)
	defer m.l.ReleaseConn()
}

func main() {
	//拿到路由
	r := RegisterHandlers()

	mh := NewMiddleWareHandler(r, 2) //最大限制连接数cc

	http.ListenAndServe(":9000", mh) //使用自定义handler：mh
}
