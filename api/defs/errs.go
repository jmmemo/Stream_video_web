package defs

//被嵌套结构体，自定义错误信息，自定义错误状态码
type Err struct {
	Error     string `json:"error"`
	ErrorCode string `json:"error_code"`
}

//定义返回的消息结构体，有状态码，和嵌套结构体【自定义错误信息，自定义错误状态码】
type ErrorResponse struct {
	HttpSC int
	Error  Err
}

var (
	ErrorRequestBodyParseFailed = ErrorResponse{
		HttpSC: 400,
		Error: Err{
			Error:     "Request body is not correct",
			ErrorCode: "001",
		},
	}

	ErrorNotAuthUser = ErrorResponse{
		HttpSC: 401,
		Error: Err{
			Error:     "User authentication failed",
			ErrorCode: "002",
		},
	}
)
