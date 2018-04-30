package msg

const (
	HttpStatus_Success  = 200
	HttpStatus_Exception = 201
	HttpStatus_InvalidParam = 202
	HttpStatus_InvalidDomain = 203
	HttpStatus_InvalidMethod = 204
)

type Resp struct {
	Code 	int  		`json:"code"`
	Data    interface{}	`json:"data"`
	Msg     string   	`json:"msg"`
}

type LatestResp struct {
	Current     int   `json:"current"`
}
