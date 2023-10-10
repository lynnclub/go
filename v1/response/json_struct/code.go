package json_struct

type Code struct {
	Code      int         `json:"code"`
	Msg       string      `json:"msg"`
	Data      interface{} `json:"data"`
	Timestamp int64       `json:"timestamp"`
}

func (json *Code) Set(status int, msg string, data interface{}, timestamp int64) {
	json.Code = status
	json.Msg = msg
	json.Data = data
	json.Timestamp = timestamp
}
