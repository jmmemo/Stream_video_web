package defs

//requests用户信息：用户名，用户密码
type UserCredential struct {
	Username string `json:"user_name"`
	Pwd      string `json:"pwd"`
}

//注册response
type SignedUp struct {
	Success   bool   `json:"success"`
	SessionId string `json:"session_id"`
}

//登录response
type SignedIn struct {
	Success   bool   `json:"success"`
	SessionId string `json:"session_id"`
}

//视频信息结构体:视频vid，作者id，视频名字name，网站上显示的时间dct
type VideoInfo struct {
	Id           string
	AuthorId     int
	Name         string
	DisplayCtime string
}

//评论结构体:评论id，视频id，作者author，评论内容content//comments.id, users.login_name, comments.content
type Comment struct {
	Id      string
	VideoId string
	Author  string
	Content string
}

//Session对象：用户名Username，生存时间TTL【int64】
type SimpleSession struct {
	Username string //login name
	TTL      int64
}

//
type UserInfo struct {
	Id int `json:"id"`
}

//
type User struct {
	Id        int
	LoginName string
	Pwd       string
}

//
type NewVideo struct {
	AuthorId int    `json:"author_id"`
	Name     string `json:"name"`
}

//
type VideosInfo struct {
	Videos []*VideoInfo `json:"videos"`
}

//
type NewComment struct {
	AuthorId int    `json:"author_id"`
	Content  string `json:"content"`
}

//
type Comments struct {
	Comments []*Comment `json:"comments"`
}

//
type UserSession struct {
	Username  string `json:"user_name"`
	SessionId string `json:"session_id"`
}
