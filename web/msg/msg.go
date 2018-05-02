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

type LatestListReq struct {
	Page	int	`json:"page"`
}


type LastetListResp struct{
	Current		int		`json:"current"`
	List		[]*LatestResp 	`json:"list"`

}

type LatestResp struct {
	Fid    		int 		`json:"fid"`
	Name   		string 		`json:"name"`
	Avater  	string  	`json:"avatar"`
	Title		string		`json:"title"`
	Description     string 		`json:"description"`
	Linkid		string 		`json:"linkid"`
	Source  	int 		`json:"source"`
	PubDate		int		`json:"pub_date"`
}


type ArticleReq struct {
	Linkid		string 		`json:"linkid"`
	Source  	int		`json:"source"`
}

type ArticleResp struct {
	Title		string 		`json:"title"`
	Content		string		`json:"content"`
	Url 		string 		`json:"url"`
}