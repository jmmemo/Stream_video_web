//单元测试
package dbops

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

var tempvid string

//清空表的信息
func clearTables() {
	dbConn.Exec("truncate users")
	dbConn.Exec("truncate video_info")
	dbConn.Exec("truncate comments")
	dbConn.Exec("truncate sessions")
}

//主测
func TestMain(m *testing.M) {
	clearTables() //测试前清空表
	m.Run()       //调用TestUserWorkFlow进行测试
	clearTables() //测试完后清空表，避免数据残留
}

//【User】的单元测试工作流，由主测调用
func TestUserWorkFlow(t *testing.T) {
	t.Run("Add", testAddUser)
	t.Run("Get", testGetUser)
	t.Run("Del", testDeleteUser)
	t.Run("Reget", testRegetUser)
}

//【Video】的单元测试工作流，由主测调用
func TestVideoWorkFlow(t *testing.T) {
	clearTables()
	t.Run("PrepareUser", testAddUser)
	t.Run("AddVideo", testAddVideoInfo)
	t.Run("GetVideo", testGetVideoInfo)
	t.Run("DelVideo", testDeleteVideoInfo)
	t.Run("RegetVideo", testRegetVideoInfo)
}

//***************************User部分************************
//测试注册功能
func testAddUser(t *testing.T) {
	err := AddUserCredential("avenssi", "123")
	if err != nil {
		t.Errorf("Error of AddUser: %v", err)
	}
}

//测试用户查询功能
func testGetUser(t *testing.T) {
	pwd, err := GetUserCredential("avenssi")
	if pwd != "123" || err != nil {
		t.Errorf("Error of GetUser")
	}
}

//测试删除用户功能
func testDeleteUser(t *testing.T) {
	err := DeleteUser("avenssi", "123")
	if err != nil {
		t.Errorf("Error of DeleteUser: %v", err)
	}
}

//测试删除用户后  是否有残留
func testRegetUser(t *testing.T) {
	pwd, err := GetUserCredential("avenssi")
	if err != nil {
		t.Errorf("Error of RegetUser: %v", err)
	}
	if pwd != "" { //删除用户信息后，密码应该为空，否则删除功能有误
		t.Errorf("Deleting user test failed")
	}
}

//***************************Video部分************************
//测试视频上传
func testAddVideoInfo(t *testing.T) {
	vi, err := AddNewVideo(1, "my-video")
	if err != nil {
		t.Errorf("Error of AddVideoinfo: %v", err)
	}
	tempvid = vi.Id
}

//测试视频信息查询
func testGetVideoInfo(t *testing.T) {
	_, err := GetVideoInfo(tempvid)
	if err != nil {
		t.Errorf("Error of GetVideoInfo: %v", err)
	}
}

//测试视频删除
func testDeleteVideoInfo(t *testing.T) {
	err := DeleteVideoInfo(tempvid)
	if err != nil {
		t.Errorf("Error of DeleteVideoinfo: %v", err)
	}
}

//测试删除视频后  是否有残留
func testRegetVideoInfo(t *testing.T) {
	vi, err := GetVideoInfo(tempvid)
	if err != nil || vi != nil { //【视频信息vi】非空，有误
		t.Errorf("Error of RegetVideoInfo: %v", err)
	}
}

//***************************Comment部分************************
//
func TestComments(t *testing.T) {
	clearTables()
	t.Run("AddUser", testAddUser)
	t.Run("AddComments", testAddComments)
	t.Run("ListComments", testListComments)
}

//
func testAddComments(t *testing.T) {
	vid := "12345"
	aid := 1
	content := "I like this video"

	err := AddNewComments(vid, aid, content)
	if err != nil {
		t.Errorf("Error of AddComments: %v", err)
	}
}

//
func testListComments(t *testing.T) {
	vid := "12345"
	from := 1514764800
	to, _ := strconv.Atoi(strconv.FormatInt(time.Now().UnixNano()/1000000000, 10))

	res, err := ListComments(vid, from, to)
	if err != nil {
		t.Errorf("Error of ListComments: %v", err)
	}

	for i, ele := range res {
		//
		fmt.Println("@%$#%@$%#$@^%#")
		fmt.Printf("comments: %d, %v \n", i, ele)
	}
}
