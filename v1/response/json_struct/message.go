package json_struct

type Message struct {
	Status    int         `json:"status"`
	Msg       string      `json:"message"`
	Data      interface{} `json:"data"`
	Timestamp int64       `json:"timestamp"`
}

func (json *Message) Set(status int, msg string, data interface{}, timestamp int64) {
	json.Status = status
	json.Msg = msg
	json.Data = data
	json.Timestamp = timestamp
}
