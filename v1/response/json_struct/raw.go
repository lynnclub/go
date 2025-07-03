package json_struct

import "encoding/json"

type Raw struct {
	Data interface{} `json:"data"`
}

func (json *Raw) Set(status int, msg string, data interface{}, timestamp int64) {
	json.Data = data
}

func (r *Raw) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Data)
}
