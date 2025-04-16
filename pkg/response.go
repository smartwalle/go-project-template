package pkg

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func NewResponse(code int, data interface{}) *Response {
	var r = &Response{}
	r.Code = code
	r.Data = data
	return r
}
