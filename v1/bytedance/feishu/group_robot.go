package feishu

import (
	"errors"
	"time"

	"github.com/lynnclub/go/v1/bytedance/feishu/entity"
	"github.com/lynnclub/go/v1/encoding/json"
	"github.com/lynnclub/go/v1/sign"
	"github.com/parnurzeal/gorequest"
)

// GroupRobot 群机器人
type GroupRobot struct {
	Webhook string //接口地址
	SignKey string //签名KEY
}

// NewGroupRobot 实例化
func NewGroupRobot(webhook, signKey string) *GroupRobot {
	return &GroupRobot{
		Webhook: webhook,
		SignKey: signKey,
	}
}

// Send 发送消息
func (robot *GroupRobot) Send(request interface{}) (response entity.GroupRobotResponse, err error) {
	// 类型检测
	msgType := ""
	switch request.(type) {
	case entity.MsgTypeText:
		msgType = "text"
	case entity.MsgTypePost:
		msgType = "post"
	case entity.MsgTypeShareChat:
		msgType = "share_chat"
	case entity.MsgTypeImage:
		msgType = "image"
	case entity.MsgTypeInteractive:
		msgType = "interactive"
	}
	if msgType == "" {
		return response, errors.New("消息类型有误")
	}

	now := time.Now().Unix()

	// 参数
	params := map[string]interface{}{
		"msg_type": msgType,
	}
	if robot.SignKey != "" {
		params["sign"], err = sign.FeiShu(robot.SignKey, now)
		if err != nil {
			return response, err
		}
	}

	requestStr := json.Encode(request)
	if msgType == "interactive" {
		// 消息卡片
		params["card"] = requestStr
	} else {
		// 默认
		params["content"] = requestStr
	}

	// 请求
	_, body, errs := gorequest.New().Post(robot.Webhook).
		Set("Content-Type", "application/json").
		SendMap(params).
		Timeout(5 * time.Second).
		End()
	if len(errs) > 0 {
		return response, errs[0]
	}

	if err = json.Decode(body, &response); err != nil {
		return response, err
	}

	if response.StatusCode != 0 {
		return response, errors.New(response.StatusMessage)
	}

	return response, nil
}

func (robot *GroupRobot) SendRich(title string, content map[string]interface{}, userId string) {
	var data entity.PostData
	data.Title = title

	// 艾特用户
	if userId == "" {
		data.Content = [][]map[string]interface{}{{content}}
	} else {
		data.Content = [][]map[string]interface{}{{content, map[string]interface{}{
			"tag":     "at",
			"user_id": userId,
		}}}
	}

	var richText entity.MsgTypePost
	richText.Post = map[string]entity.PostData{"zh_cn": data}

	_, _ = robot.Send(richText)
}
