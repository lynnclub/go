package entity

// GroupRobotResponse 群机器人响应体
type GroupRobotResponse struct {
	Extra         interface{} `json:"Extra"`
	StatusCode    int         `json:"StatusCode"`
	StatusMessage string      `json:"StatusMessage"`
}
