package session

import (
	"sync"
	"time"
	"video_server/api/dbops"
	"video_server/api/defs"
	"video_server/api/utils"
)

//全局的 并发安全 sync.map
var sessionMap *sync.Map

//初始化
func init() {
	sessionMap = &sync.Map{}
}

//当前时间int64
func nowInMilli() int64 {
	return time.Now().UnixNano() / 1000000
}

//删除操作，map中删除，DB里也要删除
func deleteExpiredSession(sid string) {
	sessionMap.Delete(sid)
	dbops.DeleteSession(sid)
}

//读取session 从DB
func LoadSessionsFromDB() {
	r, err := dbops.RetrieveAllSessions() //返回的r是一个sync map
	if err != nil {
		return
	}

	//Range sync.map
	r.Range(func(k, v interface{}) bool { //【k是sessionID，v是SimpleSession对象】
		ss := v.(*defs.SimpleSession) //ss是从SimpleSession map【v】中拿到的其中一个session对象
		sessionMap.Store(k, ss)
		return true //对应的bool
	})
}

//创建session,返回id
func GenerateNewSessionId(un string) string {
	id, _ := utils.NewUUID()
	ct := nowInMilli()
	ttl := ct + 30*60*1000

	ss := &defs.SimpleSession{
		Username: un,
		TTL:      ttl,
	}
	sessionMap.Store(id, ss) //存入全局sessionMap【k是sessionID，v是SimpleSession对象】

	return id
}

//过期判断,未过期的返回用户名Username和false
func IsSessionExpired(sid string) (string, bool) {
	ss, ok := sessionMap.Load(sid) //key是sid，value是ss
	if ok {                        //如果能取到值  进行过期判断
		ct := nowInMilli()
		if ss.(*defs.SimpleSession).TTL < ct { //过期了
			deleteExpiredSession(sid) //删除过期的session
			return "", true
		}
		return ss.(*defs.SimpleSession).Username, false //未过期
	}
	return "", true
}
