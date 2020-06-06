package defs

//用户信息：用户名，用户密码
type UserCredential struct {
	Username string `json:"user_name"`
	Pwd      string `json:"pwd"`
}
