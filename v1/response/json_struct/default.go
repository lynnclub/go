package json_struct

type Default struct {
	Status    int         `json:"status"`
	Msg       string      `json:"msg"`
	Data      interface{} `json:"data"`
	Timestamp int64       `json:"timestamp"`
}

func (json *Default) Set(status int, msg string, data interface{}, timestamp int64) {
	json.Status = status
	json.Msg = msg
	json.Data = data
	json.Timestamp = timestamp
}
