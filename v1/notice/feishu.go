package notice

import (
	"github.com/lynnclub/go/v1/bytedance/feishu"
	"github.com/lynnclub/go/v1/bytedance/feishu/entity"
)

type FeiShuGroup struct {
	Webhook string //接口地址
	SignKey string //签名KEY
	Env     string //环境
	robot   *feishu.GroupRobot
}

// NewFeiShuGroup 飞书群实例化
func NewFeiShuGroup(webhook, signKey, env string) *FeiShuGroup {
	return &FeiShuGroup{
		Webhook: webhook,
		SignKey: signKey,
		Env:     env,
		robot:   feishu.NewGroupRobot(webhook, signKey),
	}
}

// Send 飞书群发送
func (group *FeiShuGroup) Send(title string, content map[string]interface{}, userId string) {
	var data entity.PostData
	data.Title = group.Env + title

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

	_, _ = group.robot.Send(richText)
}
