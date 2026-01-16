package feishu

import (
	"testing"
)

// TestBuildTextMessage 测试构建文本消息
func TestBuildTextMessage(t *testing.T) {
	request := &GroupRobotRequest{}
	text := "这是一条测试消息"

	result := request.BuildTextMessage(text)

	if result.MsgType != "text" {
		t.Errorf("期望MsgType为text，实际为%s", result.MsgType)
	}

	if result.Content == nil {
		t.Fatal("Content不应该为nil")
	}

	contentMap, ok := result.Content.(map[string]string)
	if !ok {
		t.Fatal("Content应该是map[string]string类型")
	}

	if contentMap["text"] != text {
		t.Errorf("期望文本为%s，实际为%s", text, contentMap["text"])
	}
}

// TestBuildRichMessage 测试构建富文本消息
func TestBuildRichMessage(t *testing.T) {
	request := &GroupRobotRequest{}
	title := "测试标题"
	text := "测试内容"

	result := request.BuildRichMessage(title, text)

	if result.MsgType != "post" {
		t.Errorf("期望MsgType为post，实际为%s", result.MsgType)
	}

	if result.Content == nil {
		t.Fatal("Content不应该为nil")
	}
}

// TestBuildRichMessageWithUsers 测试构建带@用户的富文本消息
func TestBuildRichMessageWithUsers(t *testing.T) {
	request := &GroupRobotRequest{}
	title := "测试标题"
	text := "测试内容"
	userIDs := []string{"user123", "user456"}

	result := request.BuildRichMessage(title, text, userIDs...)

	if result.MsgType != "post" {
		t.Errorf("期望MsgType为post，实际为%s", result.MsgType)
	}

	if result.Content == nil {
		t.Fatal("Content不应该为nil")
	}

	// 验证Content结构
	contentMap, ok := result.Content.(map[string]any)
	if !ok {
		t.Fatal("Content应该是map[string]any类型")
	}

	post, ok := contentMap["post"].(map[string]any)
	if !ok {
		t.Fatal("post字段应该存在")
	}

	zhCn, ok := post["zh_cn"].(map[string]any)
	if !ok {
		t.Fatal("zh_cn字段应该存在")
	}

	if zhCn["title"] != title {
		t.Errorf("期望标题为%s，实际为%v", title, zhCn["title"])
	}
}

// TestBuildImageMessage 测试构建图片消息
func TestBuildImageMessage(t *testing.T) {
	request := &GroupRobotRequest{}
	imageKey := "img_v2_test_key"

	result := request.BuildImageMessage(imageKey)

	if result.MsgType != "image" {
		t.Errorf("期望MsgType为image，实际为%s", result.MsgType)
	}

	if result.Content == nil {
		t.Fatal("Content不应该为nil")
	}

	contentMap, ok := result.Content.(map[string]string)
	if !ok {
		t.Fatal("Content应该是map[string]string类型")
	}

	if contentMap["image_key"] != imageKey {
		t.Errorf("期望image_key为%s，实际为%s", imageKey, contentMap["image_key"])
	}
}

// TestBuildShareMessage 测试构建分享群名片消息
func TestBuildShareMessage(t *testing.T) {
	request := &GroupRobotRequest{}
	chatID := "oc_test_chat_id"

	result := request.BuildShareMessage(chatID)

	if result.MsgType != "share_chat" {
		t.Errorf("期望MsgType为share_chat，实际为%s", result.MsgType)
	}

	if result.Content == nil {
		t.Fatal("Content不应该为nil")
	}

	contentMap, ok := result.Content.(map[string]string)
	if !ok {
		t.Fatal("Content应该是map[string]string类型")
	}

	if contentMap["share_chat_id"] != chatID {
		t.Errorf("期望share_chat_id为%s，实际为%s", chatID, contentMap["share_chat_id"])
	}
}

// TestBuildCardMessage 测试构建交互式卡片消息
func TestBuildCardMessage(t *testing.T) {
	request := &GroupRobotRequest{}
	card := map[string]interface{}{
		"header": map[string]interface{}{
			"title": map[string]interface{}{
				"tag":     "plain_text",
				"content": "这是卡片标题",
			},
		},
		"elements": []interface{}{
			map[string]interface{}{
				"tag": "div",
				"text": map[string]interface{}{
					"tag":     "plain_text",
					"content": "这是卡片内容",
				},
			},
		},
	}

	result := request.BuildCardMessage(card)

	if result.MsgType != "interactive" {
		t.Errorf("期望MsgType为interactive，实际为%s", result.MsgType)
	}

	if result.Card == nil {
		t.Fatal("Card不应该为nil")
	}

	cardMap, ok := result.Card.(map[string]interface{})
	if !ok {
		t.Fatal("Card应该是map[string]interface{}类型")
	}

	if cardMap["header"] == nil {
		t.Error("Card应该包含header字段")
	}
}

// TestBuildAdvancedRichMessage 测试构建复杂富文本消息
func TestBuildAdvancedRichMessage(t *testing.T) {
	request := &GroupRobotRequest{}
	title := "复杂富文本标题"
	elements := []RichElement{
		{Type: "text", Text: "这是文本"},
		{Type: "at", UserID: "user123"},
		{Type: "a", Text: "链接", Href: "https://example.com"},
		{Type: "img", Key: "img_key_123"},
	}

	result := request.BuildAdvancedRichMessage(title, elements)

	if result.MsgType != "post" {
		t.Errorf("期望MsgType为post，实际为%s", result.MsgType)
	}

	if result.Content == nil {
		t.Fatal("Content不应该为nil")
	}

	contentMap, ok := result.Content.(map[string]any)
	if !ok {
		t.Fatal("Content应该是map[string]any类型")
	}

	post, ok := contentMap["post"].(map[string]any)
	if !ok {
		t.Fatal("post字段应该存在")
	}

	zhCn, ok := post["zh_cn"].(map[string]any)
	if !ok {
		t.Fatal("zh_cn字段应该存在")
	}

	if zhCn["title"] != title {
		t.Errorf("期望标题为%s，实际为%v", title, zhCn["title"])
	}

	content, ok := zhCn["content"].([][]any)
	if !ok {
		t.Fatal("content字段应该是[][]any类型")
	}

	if len(content) != 1 {
		t.Errorf("期望content有1个段落，实际有%d个", len(content))
	}

	if len(content[0]) != 4 {
		t.Errorf("期望第一个段落有4个元素，实际有%d个", len(content[0]))
	}
}

// TestRichElementTypes 测试不同类型的富文本元素
func TestRichElementTypes(t *testing.T) {
	testCases := []struct {
		name     string
		element  RichElement
		expected string
	}{
		{"文本元素", RichElement{Type: "text", Text: "测试文本"}, "text"},
		{"@用户元素", RichElement{Type: "at", UserID: "user123"}, "at"},
		{"链接元素", RichElement{Type: "a", Text: "点击", Href: "https://example.com"}, "a"},
		{"图片元素", RichElement{Type: "img", Key: "img_key"}, "img"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.element.Type != tc.expected {
				t.Errorf("期望类型为%s，实际为%s", tc.expected, tc.element.Type)
			}
		})
	}
}

// TestChainedBuilding 测试链式调用
func TestChainedBuilding(t *testing.T) {
	request := &GroupRobotRequest{}

	// 测试链式调用返回自身
	result := request.BuildTextMessage("test")
	if result != request {
		t.Error("BuildTextMessage应该返回自身以支持链式调用")
	}

	request2 := &GroupRobotRequest{}
	result2 := request2.BuildRichMessage("title", "text")
	if result2 != request2 {
		t.Error("BuildRichMessage应该返回自身以支持链式调用")
	}

	request3 := &GroupRobotRequest{}
	result3 := request3.BuildImageMessage("key")
	if result3 != request3 {
		t.Error("BuildImageMessage应该返回自身以支持链式调用")
	}

	request4 := &GroupRobotRequest{}
	result4 := request4.BuildShareMessage("chat_id")
	if result4 != request4 {
		t.Error("BuildShareMessage应该返回自身以支持链式调用")
	}

	request5 := &GroupRobotRequest{}
	result5 := request5.BuildCardMessage(map[string]interface{}{})
	if result5 != request5 {
		t.Error("BuildCardMessage应该返回自身以支持链式调用")
	}
}

// TestEmptyRichElements 测试空元素列表
func TestEmptyRichElements(t *testing.T) {
	request := &GroupRobotRequest{}
	title := "空元素标题"
	elements := []RichElement{}

	result := request.BuildAdvancedRichMessage(title, elements)

	if result.MsgType != "post" {
		t.Errorf("期望MsgType为post，实际为%s", result.MsgType)
	}

	contentMap := result.Content.(map[string]any)
	post := contentMap["post"].(map[string]any)
	zhCn := post["zh_cn"].(map[string]any)
	content := zhCn["content"].([][]any)

	if len(content[0]) != 0 {
		t.Errorf("期望空元素列表，实际有%d个元素", len(content[0]))
	}
}
