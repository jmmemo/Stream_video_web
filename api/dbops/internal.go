package dbops

import (
	"database/sql"
	"log"
	"strconv"
	"sync"
	"video_server/api/defs"
)

//写入session，接收session的id，ttl，用户名ttl
func InserSession(sid string, ttl int64, uname string) error {
	ttlstr := strconv.FormatInt(ttl, 10) //把int64型的ttl转换成string，因为要放到sql语句中执行
	stmtIns, err := dbConn.Prepare("INSERT INTO sessions (session_id, TTL, login_name) VALUES (?,?,?)")
	if err != nil {
		return err
	}
	_, err = stmtIns.Exec(sid, ttlstr, uname)
	if err != nil {
		return err
	}
	defer stmtIns.Close()
	return nil
}

//获取session
func RetrieveSession(sid string) (*defs.SimpleSession, error) {
	ss := &defs.SimpleSession{} //将被返回的session对象
	stmtOut, err := dbConn.Prepare("SELECT TTL, login_name FROM sessions WHERE session_id =?")
	if err != nil {
		return nil, err
	}
	var ttl string
	var uname string
	err = stmtOut.QueryRow(sid).Scan(&ttl, &uname) //读到的TTL，login_name写入【string型的：】ttl，uname
	if err != nil && err != sql.ErrNoRows {        //errNoRows并非错误
		return nil, err
	}

	if res, err := strconv.ParseInt(ttl, 10, 64); err == nil { //string型ttl转换成int64
		ss.TTL = res //转换成int64的ttl（res）被写入Session对象ss
		ss.Username = uname
	} else {
		return nil, err
	}
	defer stmtOut.Close()
	return ss, nil
}

//获得所有session，使用并发安全 sync.map
func RetrieveAllSessions() (*sync.Map, error) {
	m := &sync.Map{} //将被返回的map
	//sessions表里有三个属性： 1 session_id 2 TTL 3 login_name 均为类似string型
	stmtOut, err := dbConn.Prepare("SELECT * FROM sessions")
	if err != nil {
		log.Printf("%s", err)
		return nil, err
	}

	rows, err := stmtOut.Query() //执行sql语句从DB拿到 多行数据 rows
	if err != nil {
		log.Printf("%s", err)
		return nil, err
	}

	for rows.Next() { //遍历rows
		var (
			id         string
			ttlstr     string
			login_name string
		)
		//从rows中读到的属性写入上面定义的三个元素
		if err := rows.Scan(&id, &ttlstr, &login_name); err != nil {
			log.Printf("retrive sessions error: %s", err)
			break
		}

		//把string类型的ttl转换成int64，因SimpleSession使用的TTL是int64
		if ttl, err := strconv.ParseInt(ttlstr, 10, 64); err == nil {
			ss := &defs.SimpleSession{
				Username: login_name,
				TTL:      ttl,
			}
			m.Store(id, ss) //以k-v【k是sessionID，v是SimpleSession对象】 键值对形式存入map中
			log.Printf("session id: %s, ttl: %d", id, ss.TTL)
		}
	}
	return m, nil
}

//删除session，根据sessionID
func DeleteSession(sid string) error {
	stmtOut, err := dbConn.Prepare("DELETE FROM sessions WHERE session_id = ?")
	if err != nil {
		log.Printf("%s", err)
		return err
	}

	if _, err := stmtOut.Query(sid); err != nil {
		return err
	}
	return nil
}
