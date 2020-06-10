package main

import (
	"io"
	"net/http"
)

//发送错误响应信息
func sendErrorResponse(w http.ResponseWriter, sc int, errMsg string) {
	w.WriteHeader(sc)
	io.WriteString(w, errMsg)
}
