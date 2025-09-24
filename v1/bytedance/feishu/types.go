package feishu

// 文档
// https://open.feishu.cn/document/client-docs/bot-v3/add-custom-bot

// GroupRobotResponse 飞书群机器人API响应结构
type GroupRobotResponse struct {
	Code int    `json:"code"` // 状态码，0表示成功
	Msg  string `json:"msg"`  // 状态消息，描述请求结果
	Data any    `json:"data"` // 响应数据
}

// GroupRobotRequest 飞书群机器人API请求结构
type GroupRobotRequest struct {
	MsgType string `json:"msg_type"`          // 消息类型
	Sign    string `json:"sign,omitempty"`    // 签名（可选）
	Content any    `json:"content,omitempty"` // 消息内容（普通消息使用）
	Card    any    `json:"card,omitempty"`    // 卡片内容（交互式消息使用）
}

// BuildTextMessage 构建文本消息
func (r *GroupRobotRequest) BuildTextMessage(text string) *GroupRobotRequest {
	r.MsgType = "text"
	r.Content = map[string]string{"text": text}
	return r
}

// BuildRichMessage 构建富文本消息
func (r *GroupRobotRequest) BuildRichMessage(title, text string, userIDs ...string) *GroupRobotRequest {
	r.MsgType = "post"

	elements := make([]any, 0, len(userIDs)+1)
	elements = append(elements, map[string]any{
		"tag":  "text",
		"text": text,
	})

	for _, userID := range userIDs {
		elements = append(elements, map[string]any{
			"tag":     "at",
			"user_id": userID,
		})
	}

	r.Content = map[string]any{
		"post": map[string]any{
			"zh_cn": map[string]any{
				"title":   title,
				"content": [][]any{elements},
			},
		},
	}
	return r
}

// BuildImageMessage 构建图片消息
func (r *GroupRobotRequest) BuildImageMessage(imageKey string) *GroupRobotRequest {
	r.MsgType = "image"
	r.Content = map[string]string{"image_key": imageKey}
	return r
}

// BuildShareMessage 构建分享群名片消息
func (r *GroupRobotRequest) BuildShareMessage(chatID string) *GroupRobotRequest {
	r.MsgType = "share_chat"
	r.Content = map[string]string{"share_chat_id": chatID}
	return r
}

// BuildCardMessage 构建交互式卡片消息
func (r *GroupRobotRequest) BuildCardMessage(card any) *GroupRobotRequest {
	r.MsgType = "interactive"
	r.Card = card
	return r
}

// 为了方便构建复杂的富文本消息，提供辅助结构体

// RichElement 富文本元素辅助结构体
type RichElement struct {
	Type   string `json:"tag"`                 // 元素类型: text, at, img, a (链接)
	Text   string `json:"text,omitempty"`      // 文本内容
	UserID string `json:"user_id,omitempty"`   // @用户ID (用于at类型)
	Href   string `json:"href,omitempty"`      // 链接地址 (用于a类型)
	Key    string `json:"image_key,omitempty"` // 图片key (用于img类型)
}

// BuildAdvancedRichMessage 构建复杂的富文本消息
func (r *GroupRobotRequest) BuildAdvancedRichMessage(title string, elements []RichElement) *GroupRobotRequest {
	r.MsgType = "post"

	content := make([]any, 0, len(elements))
	for _, elem := range elements {
		elemMap := map[string]any{"tag": elem.Type}
		if elem.Text != "" {
			elemMap["text"] = elem.Text
		}
		if elem.UserID != "" {
			elemMap["user_id"] = elem.UserID
		}
		if elem.Href != "" {
			elemMap["href"] = elem.Href
		}
		if elem.Key != "" {
			elemMap["image_key"] = elem.Key
		}
		content = append(content, elemMap)
	}

	r.Content = map[string]any{
		"post": map[string]any{
			"zh_cn": map[string]any{
				"title":   title,
				"content": [][]any{content},
			},
		},
	}
	return r
}
