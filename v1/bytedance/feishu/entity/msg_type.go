package entity

// 文档：https://www.feishu.cn/hc/zh-CN/articles/360024984973#lineguid-CdVcCt

// MsgTypeText 文本消息
type MsgTypeText struct {
	Text string `json:"text"`
}

// MsgTypePost 富文本消息
type MsgTypePost struct {
	Post map[string]PostData `json:"post"`
}

// MsgTypeShareChat 分享群名片消息
type MsgTypeShareChat struct {
	ShareChatID string `json:"share_chat_id"`
}

// MsgTypeImage 图片消息
type MsgTypeImage struct {
	ImageKey string `json:"image_key"`
}

// MsgTypeInteractive 消息卡片
type MsgTypeInteractive struct {
	Config   map[string]interface{} `json:"config"`
	Elements interface{}            `json:"elements"`
	Header   interface{}            `json:"header"`
}
