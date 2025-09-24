package feishu

import (
	"errors"
	"time"

	"github.com/lynnclub/go/v1/encoding/json"
	"github.com/lynnclub/go/v1/sign"
	"github.com/parnurzeal/gorequest"
)

// GroupRobot 飞书群机器人客户端
type GroupRobot struct {
	Webhook string // 飞书机器人Webhook地址
	SignKey string // 机器人签名密钥，用于验证消息来源
}

// NewGroupRobot 创建群机器人客户端
func NewGroupRobot(webhook, signKey string) *GroupRobot {
	return &GroupRobot{
		Webhook: webhook,
		SignKey: signKey,
	}
}

// SendRaw 发送原始参数到飞书API
func (robot *GroupRobot) SendRaw(params interface{}) (response GroupRobotResponse, err error) {
	_, body, errs := gorequest.New().Post(robot.Webhook).
		Set("Content-Type", "application/json").
		Send(params).
		Timeout(3 * time.Second).
		End()
	if len(errs) > 0 {
		return response, errs[0]
	}

	if err = json.Decode(body, &response); err != nil {
		return response, err
	}

	if response.Code != 0 {
		return response, errors.New(response.Msg)
	}

	return response, nil
}

// Send 发送消息
func (robot *GroupRobot) Send(request *GroupRobotRequest) (response GroupRobotResponse, err error) {
	if robot.SignKey != "" {
		signValue, err := sign.FeiShu(robot.SignKey, time.Now().Unix())
		if err != nil {
			return response, err
		}

		request.Sign = signValue
	}

	return robot.SendRaw(request)
}

// SendText 发送文本消息（快捷方法）
func (robot *GroupRobot) SendText(text string) (GroupRobotResponse, error) {
	request := &GroupRobotRequest{}
	request.BuildTextMessage(text)
	return robot.Send(request)
}

// SendRich 发送富文本消息（快捷方法）
func (robot *GroupRobot) SendRich(title, text, userId string) (GroupRobotResponse, error) {
	request := &GroupRobotRequest{}
	if userId == "" {
		request.BuildRichMessage(title, text)
	} else {
		request.BuildRichMessage(title, text, userId)
	}
	return robot.Send(request)
}

// SendImage 发送图片消息（快捷方法）
func (robot *GroupRobot) SendImage(imageKey string) (GroupRobotResponse, error) {
	request := &GroupRobotRequest{}
	request.BuildImageMessage(imageKey)
	return robot.Send(request)
}

// SendShare 发送分享群名片消息（快捷方法）
func (robot *GroupRobot) SendShare(chatID string) (GroupRobotResponse, error) {
	request := &GroupRobotRequest{}
	request.BuildShareMessage(chatID)
	return robot.Send(request)
}

// SendCard 发送交互式卡片消息（快捷方法）
func (robot *GroupRobot) SendCard(card any) (GroupRobotResponse, error) {
	request := &GroupRobotRequest{}
	request.BuildCardMessage(card)
	return robot.Send(request)
}
