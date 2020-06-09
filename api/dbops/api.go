package dbops

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
	"video_server/api/defs"
	"video_server/api/utils"
)

//***************************用户部分************************

//【C 增】注册用户
func AddUserCredential(loginName string, pwd string) error {
	stmtIns, err := dbConn.Prepare("INSERT INTO users (login_name,pwd) VALUES (?,?)")
	if err != nil {
		return err
	}
	_, err = stmtIns.Exec(loginName, pwd)
	if err != nil {
		return err
	}
	defer stmtIns.Close()
	return nil
}

//【R 查】查询用户信息，返回其密码
func GetUserCredential(loginName string) (string, error) {
	stmtOut, err := dbConn.Prepare("SELECT pwd FROM users WHERE login_name=?")
	if err != nil {
		log.Printf("%s", err)
		return "", err
	}

	var pwd string
	err = stmtOut.QueryRow(loginName).Scan(&pwd)
	if err != nil && err != sql.ErrNoRows {
		return "", nil
	}
	defer stmtOut.Close()

	return pwd, nil
}

//【D 删】删除用户
func DeleteUser(loginName string, pwd string) error {
	stmtDel, err := dbConn.Prepare("DELETE FROM users WHERE login_name=? AND pwd=?")
	if err != nil {
		log.Printf("DeleteUser error:%s", err)
		return err
	}
	_, err = stmtDel.Exec(loginName, pwd)
	if err != nil {
		return err
	}
	defer stmtDel.Close()
	return nil
}

//***************************视频部分************************

//【C 增】：上传视频，返回视频信息结构体
func AddNewVideo(aid int, name string) (*defs.VideoInfo, error) {
	vid, err := utils.NewUUID() //生成唯一id作为【video_info】的主键
	if err != nil {
		return nil, err
	}
	t := time.Now()
	ctime := t.Format("Jan 02 2006, 15:04:05") //格式化时间，作为【video_info】的display_ctime
	stmtIns, err := dbConn.Prepare("INSERT INTO video_info(id,author_id,name,display_ctime)VALUES (?,?,?,?)")
	if err != nil {
		return nil, err
	}
	_, err = stmtIns.Exec(vid, aid, name, ctime) //执行语句
	if err != nil {
		return nil, err
	}
	res := &defs.VideoInfo{ //返回所插入的视频信息
		Id:           vid,
		AuthorId:     aid,
		Name:         name,
		DisplayCtime: ctime,
	}
	defer stmtIns.Close()
	return res, nil
}

//【R 查】：查视频信息，返回视频信息结构体
func GetVideoInfo(vid string) (*defs.VideoInfo, error) {
	stmtOut, err := dbConn.Prepare("SELECT author_id, name, display_ctime FROM video_info WHERE id=?")

	var ( //定义待会被返回的结构体组成元素
		aid  int
		dct  string
		name string
	)

	err = stmtOut.QueryRow(vid).Scan(&aid, &name, &dct)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if err == sql.ErrNoRows { //如果没读到信息
		return nil, nil
	}

	defer stmtOut.Close()

	res := &defs.VideoInfo{ //被返回的视频信息结构体
		Id:           vid,
		AuthorId:     aid,
		Name:         name,
		DisplayCtime: dct,
	}

	return res, nil
}

//【D 删】：删除视频信息
func DeleteVideoInfo(vid string) error {
	stmtDel, err := dbConn.Prepare("DELETE FROM video_info WHERE id = ?")
	if err != nil {
		return err
	}

	_, err = stmtDel.Exec(vid)
	if err != nil {
		return err
	}
	defer stmtDel.Close()
	return nil
}

//***************************评论部分************************
//【C 增】：给视频添加评论，传入视频vid，作者名字aid，及评论内容content
func AddNewComments(vid string, aid int, content string) error {
	id, err := utils.NewUUID() //uuid作为comments表的主键
	if err != nil {
		return err
	}

	stmtIns, err := dbConn.Prepare("INSERT INTO comments (id, video_id, author_id, content) VALUES(?,?,?,?)")
	if err != nil {
		return err
	}

	_, err = stmtIns.Exec(id, vid, aid, content)
	if err != nil {
		return err
	}

	defer stmtIns.Close()
	return nil
}

//【R 查】列出视频评论，传入视频vid，及要显示的时间范围
func ListComments(vid string, from, to int) ([]*defs.Comment, error) {
	stmtOut, err := dbConn.Prepare(`SELECT comments.id, users.login_name, comments.content FROM comments
					INNER JOIN users ON comments.author_id = users.id
					WHERE comments.video_id = ? AND comments.time > FROM_UNIXTIME(?) AND comments.time <= FROM_UNIXTIME(?)
					ORDER BY comments.time DESC `)

	var res []*defs.Comment //被返回的评论结构体对象 数组

	rows, err := stmtOut.Query(vid, from, to)
	if err != nil {
		return res, err
	}

	for rows.Next() { //遍历 query语句拿到的rows
		//分别是：评论id，作者名name（用户名），评论内容content
		var id, name, content string
		if err := rows.Scan(&id, &name, &content); err != nil { //把读到的内容写入对应的评论结构体元素
			return res, err
		}

		c := &defs.Comment{ //将评论结构体元素写入c中
			Id:      id,
			VideoId: vid,
			Author:  name,
			Content: content,
		}
		res = append(res, c) //加入 被返回的评论结构体对象 数组
	}
	defer stmtOut.Close()
	return res, nil
}

//
func GetUser(loginName string) (*defs.User, error) {
	stmtOut, err := dbConn.Prepare("SELECT id, pwd FROM users WHERE login_name=?")
	if err != nil {
		log.Printf("%s", err)
		return nil, err
	}
	var (
		id  int
		pwd string
	)
	err = stmtOut.QueryRow(loginName).Scan(&id, &pwd)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if err == sql.ErrNoRows {
		return nil, nil
	}
	res := &defs.User{
		Id:        id,
		LoginName: loginName,
		Pwd:       pwd,
	}

	defer stmtOut.Close()
	return res, nil

}

//
func ListVideoInfo(uname string, from, to int) ([]*defs.VideoInfo, error) {
	stmtOut, err := dbConn.Prepare(`SELECT video_info.id, video_info.author_id, video_info.name, video_info.display_ctime FROM video_info
		INNER JOIN users ON video_info.author_id = users.id
		WHERE users.login_name=? AND video_info.create_time > FROM_UNIXTIME(?) AND video_info.create_time<=FROM_UNIXTIME(?)
		OREDER BY video_info.create_time DESC`)
	var res []*defs.VideoInfo
	if err != nil {
		return res, err
	}
	rows, err := stmtOut.Query(uname, from, to)
	if err != nil {
		log.Printf("%s", err)
		return res, err
	}

	for rows.Next() {
		var id, name, ctime string
		var aid int
		if err := rows.Scan(&id, &aid, &name, &ctime); err != nil {
			return res, err
		}
		vi := &defs.VideoInfo{
			Id:           id,
			AuthorId:     aid,
			Name:         name,
			DisplayCtime: ctime,
		}
		res = append(res, vi)
	}
	defer stmtOut.Close()
	return res, nil
}
