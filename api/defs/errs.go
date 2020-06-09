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
	//解析错误
	ErrorRequestBodyParseFailed = ErrorResponse{
		HttpSC: 400,
		Error: Err{
			Error:     "Request body is not correct",
			ErrorCode: "001",
		},
	}
	//未授权错误
	ErrorNotAuthUser = ErrorResponse{
		HttpSC: 401,
		Error: Err{
			Error:     "User authentication failed",
			ErrorCode: "002",
		},
	}
	//数据库错误
	ErrorDBError = ErrorResponse{
		HttpSC: 500,
		Error: Err{
			Error:     "DB ops failed",
			ErrorCode: "003",
		},
	}
	//
	ErrorInternalFaults = ErrorResponse{
		HttpSC: 500,
		Error: Err{
			Error:     "Internal service error",
			ErrorCode: "004",
		},
	}
)
