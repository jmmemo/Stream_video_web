package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func RegisterHandlers() *httprouter.Router {
	//1创建路由
	router := httprouter.New()

	//2给路由绑定方法
	router.POST("/user", CreateUser)

	//3返回此路由return
	return router

}

func main() {
	//1拿到路由器
	r := RegisterHandlers()

	//2开始监听
	http.ListenAndServe(":8000", r)

}
