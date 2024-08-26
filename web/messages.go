/**
  @author: decision
  @date: 2024/8/21
  @note:
**/

package web

type errCode int

const (
	errRequestPathInvalid errCode = 10001
	errIndexNotFound      errCode = 10002
	errParseFailed        errCode = 10101
	errParseIdentFailed   errCode = 10102
	errLoadFSIdentity     errCode = 10201
)

var errMsg = map[errCode]string{
	errRequestPathInvalid: "request path invalid",
	errIndexNotFound:      "local index file not exists",
	errParseFailed:        "parse on-chain failed",
	errParseIdentFailed:   "parse fs identity failed",
	errLoadFSIdentity:     "load on-chain file storage index failed",
}

type Message struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func SimpleMsg(code errCode) Message {
	return Message{
		Code: int(code),
		Msg:  errMsg[code],
	}
}
