//单元测试
package dbops

import "testing"

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
