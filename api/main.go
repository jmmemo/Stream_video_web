package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"video_server/api/session"
)

//中间件
type middleWareHandler struct {
	r *httprouter.Router
}

//劫持中间件的handler，返回自定义handler
func NewMiddleWareHandler(r *httprouter.Router) http.Handler {
	m := middleWareHandler{}
	m.r = r
	return m
}

//自定义中间件的handler，加入 验证session步骤
func (m middleWareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//check session
	validateUserSession(r)
	m.r.ServeHTTP(w, r)
}

func RegisterHandlers() *httprouter.Router {
	//1创建路由
	router := httprouter.New()

	//2给路由绑定方法
	router.POST("/user", CreateUser)

	router.POST("/user/:username", Login)

	router.GET("/user/:username", GetUserInfo)

	router.POST("/user/:username/videos", AddNewVideo)

	router.GET("/user/:username/videos", ListAllVideos)

	router.DELETE("/user/:username/videos/:vid-id", DeleteVideo)

	router.POST("/videos/:vid-id/comments", PostComment)

	router.GET("/videos/:vid-id/comments", ShowComments)

	//3返回此路由return
	return router

}

func Prepare() {
	session.LoadSessionsFromDB()
}

func main() {
	Prepare()

	//1拿到路由
	r := RegisterHandlers()

	mh := NewMiddleWareHandler(r)

	//2开始监听
	http.ListenAndServe(":8000", mh) //使用自定义handler：mh

}
