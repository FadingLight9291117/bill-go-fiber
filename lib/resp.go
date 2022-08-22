package lib

type ResponseBody struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func Resp(data interface{}) *ResponseBody { return &ResponseBody{Code: 0, Data: data} }
func EmptyResp() *ResponseBody            { return &ResponseBody{Code: 0} }

//func ErrorResp(err error) *ResponseBody   { return &ResponseBody{Message: err.Error()} }
