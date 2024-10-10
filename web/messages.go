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
	errFilePathNotFound   errCode = 10103
	errIdentityNotFound   errCode = 10104
	errLoadFSIdentity     errCode = 10201
	errApplicationInvalid errCode = 10301
)

var errMsg = map[errCode]string{
	errRequestPathInvalid: "request path invalid",
	errIndexNotFound:      "local index file not exists",
	errParseFailed:        "parse on-chain failed",
	errParseIdentFailed:   "parse fs identity failed",
	errFilePathNotFound:   "file path not found",
	errIdentityNotFound:   "identity not found in cache",
	errLoadFSIdentity:     "load on-chain file storage index failed",
	errApplicationInvalid: "application invalid, waiting for reload",
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
