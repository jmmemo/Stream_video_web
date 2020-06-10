package main

import (
	"github.com/julienschmidt/httprouter"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

//上传页面处理 模版
func testPageHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	t, _ := template.ParseFiles("./videos/upload.html")

	t.Execute(w, nil)
}

//视频播放流处理
func streamHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	vid := p.ByName("vid-id")
	vl := VIDEO_DIR + vid //文件链接路径

	video, err := os.Open(vl) //打开文件
	if err != nil {
		log.Printf("Error when try to open file: %v", err)
		sendErrorResponse(w, http.StatusInternalServerError, "Internal Error")
		return
	}
	//设置响应头 提示浏览器这是mp4类型
	w.Header().Set("Content-Type", "video/mp4")
	//···
	http.ServeContent(w, r, "", time.Now(), video)

	defer video.Close()
}

//上传视频
func uploadHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	r.Body = http.MaxBytesReader(w, r.Body, MAX_UPLOAD_SIZE)
	if err := r.ParseMultipartForm(MAX_UPLOAD_SIZE); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "File is too big")
		return
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Internal Error")
		return
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("Read file error: %v", err)
		sendErrorResponse(w, http.StatusInternalServerError, "Internal Error")
	}

	fn := p.ByName("vid-id")
	err = ioutil.WriteFile(VIDEO_DIR+fn, data, 0666)
	if err != nil {
		log.Printf("Write file error: %v", err)
		sendErrorResponse(w, http.StatusInternalServerError, "Internal Error")
		return
	}
	w.WriteHeader(http.StatusCreated)
	io.WriteString(w, "Uploaded successfully")
}
