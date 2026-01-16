package feishu

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestNewGroupRobot 测试创建群机器人客户端
func TestNewGroupRobot(t *testing.T) {
	webhook := "https://open.feishu.cn/open-apis/bot/v2/hook/test"
	signKey := "test_sign_key"

	robot := NewGroupRobot(webhook, signKey)

	if robot == nil {
		t.Error("创建群机器人客户端失败")
	}

	if robot.Webhook != webhook {
		t.Errorf("期望Webhook为%s，实际为%s", webhook, robot.Webhook)
	}

	if robot.SignKey != signKey {
		t.Errorf("期望SignKey为%s，实际为%s", signKey, robot.SignKey)
	}
}

// TestNewGroupRobotWithoutSignKey 测试创建不带签名的群机器人客户端
func TestNewGroupRobotWithoutSignKey(t *testing.T) {
	webhook := "https://open.feishu.cn/open-apis/bot/v2/hook/test"

	robot := NewGroupRobot(webhook, "")

	if robot == nil {
		t.Error("创建群机器人客户端失败")
	}

	if robot.Webhook != webhook {
		t.Errorf("期望Webhook为%s，实际为%s", webhook, robot.Webhook)
	}

	if robot.SignKey != "" {
		t.Errorf("期望SignKey为空，实际为%s", robot.SignKey)
	}
}

// TestGroupRobotSendTextStructure 测试文本消息结构
func TestGroupRobotSendTextStructure(t *testing.T) {
	_ = NewGroupRobot("https://example.com", "")

	// 我们不实际发送，只测试请求构建
	request := &GroupRobotRequest{}
	request.BuildTextMessage("测试消息")

	if request.MsgType != "text" {
		t.Errorf("期望MsgType为text，实际为%s", request.MsgType)
	}

	if request.Content == nil {
		t.Error("Content不应该为nil")
	}
}

// TestGroupRobotSendRichStructure 测试富文本消息结构
func TestGroupRobotSendRichStructure(t *testing.T) {
	_ = NewGroupRobot("https://example.com", "")

	request := &GroupRobotRequest{}
	request.BuildRichMessage("测试标题", "测试内容", "user123")

	if request.MsgType != "post" {
		t.Errorf("期望MsgType为post，实际为%s", request.MsgType)
	}

	if request.Content == nil {
		t.Error("Content不应该为nil")
	}
}

// TestGroupRobotSendImageStructure 测试图片消息结构
func TestGroupRobotSendImageStructure(t *testing.T) {
	_ = NewGroupRobot("https://example.com", "")

	request := &GroupRobotRequest{}
	request.BuildImageMessage("img_test_key")

	if request.MsgType != "image" {
		t.Errorf("期望MsgType为image，实际为%s", request.MsgType)
	}

	if request.Content == nil {
		t.Error("Content不应该为nil")
	}
}

// TestGroupRobotSendShareStructure 测试分享消息结构
func TestGroupRobotSendShareStructure(t *testing.T) {
	_ = NewGroupRobot("https://example.com", "")

	request := &GroupRobotRequest{}
	request.BuildShareMessage("chat_test_id")

	if request.MsgType != "share_chat" {
		t.Errorf("期望MsgType为share_chat，实际为%s", request.MsgType)
	}

	if request.Content == nil {
		t.Error("Content不应该为nil")
	}
}

// TestGroupRobotSendCardStructure 测试卡片消息结构
func TestGroupRobotSendCardStructure(t *testing.T) {
	_ = NewGroupRobot("https://example.com", "")

	card := map[string]interface{}{
		"header": map[string]interface{}{
			"title": map[string]interface{}{
				"tag":     "plain_text",
				"content": "测试卡片",
			},
		},
	}

	request := &GroupRobotRequest{}
	request.BuildCardMessage(card)

	if request.MsgType != "interactive" {
		t.Errorf("期望MsgType为interactive，实际为%s", request.MsgType)
	}

	if request.Card == nil {
		t.Error("Card不应该为nil")
	}
}

// TestSendRawSuccess 测试SendRaw成功的情况
func TestSendRawSuccess(t *testing.T) {
	// 创建mock服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code":0,"msg":"success","data":{}}`))
	}))
	defer server.Close()

	robot := NewGroupRobot(server.URL, "")
	request := &GroupRobotRequest{}
	request.BuildTextMessage("测试消息")

	response, err := robot.SendRaw(request)

	if err != nil {
		t.Errorf("SendRaw不应该返回错误: %v", err)
	}

	if response.Code != 0 {
		t.Errorf("期望Code为0，实际为%d", response.Code)
	}

	if response.Msg != "success" {
		t.Errorf("期望Msg为success，实际为%s", response.Msg)
	}
}

// TestSendRawError 测试SendRaw错误响应
func TestSendRawError(t *testing.T) {
	// 创建mock服务器返回错误
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code":9499,"msg":"Bad Request","data":{}}`))
	}))
	defer server.Close()

	robot := NewGroupRobot(server.URL, "")
	request := &GroupRobotRequest{}
	request.BuildTextMessage("测试消息")

	response, err := robot.SendRaw(request)

	if err == nil {
		t.Error("SendRaw应该返回错误")
	}

	if response.Code != 9499 {
		t.Errorf("期望Code为9499，实际为%d", response.Code)
	}
}

// TestSendRawInvalidJSON 测试SendRaw无效JSON响应
func TestSendRawInvalidJSON(t *testing.T) {
	// 创建mock服务器返回无效JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`invalid json`))
	}))
	defer server.Close()

	robot := NewGroupRobot(server.URL, "")
	request := &GroupRobotRequest{}
	request.BuildTextMessage("测试消息")

	_, err := robot.SendRaw(request)

	if err == nil {
		t.Error("SendRaw应该返回JSON解析错误")
	}
}

// TestSendWithoutSignKey 测试不带签名发送
func TestSendWithoutSignKey(t *testing.T) {
	// 创建mock服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code":0,"msg":"success","data":{}}`))
	}))
	defer server.Close()

	robot := NewGroupRobot(server.URL, "")
	request := &GroupRobotRequest{}
	request.BuildTextMessage("测试消息")

	response, err := robot.Send(request)

	if err != nil {
		t.Errorf("Send不应该返回错误: %v", err)
	}

	if response.Code != 0 {
		t.Errorf("期望Code为0，实际为%d", response.Code)
	}

	// 验证没有签名
	if request.Sign != "" {
		t.Errorf("期望Sign为空，实际为%s", request.Sign)
	}
}

// TestSendWithSignKey 测试带签名发送
func TestSendWithSignKey(t *testing.T) {
	// 创建mock服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code":0,"msg":"success","data":{}}`))
	}))
	defer server.Close()

	robot := NewGroupRobot(server.URL, "test_sign_key")
	request := &GroupRobotRequest{}
	request.BuildTextMessage("测试消息")

	response, err := robot.Send(request)

	if err != nil {
		t.Errorf("Send不应该返回错误: %v", err)
	}

	if response.Code != 0 {
		t.Errorf("期望Code为0，实际为%d", response.Code)
	}

	// 验证有签名
	if request.Sign == "" {
		t.Error("期望Sign不为空")
	}
}

// TestSendText 测试SendText快捷方法
func TestSendText(t *testing.T) {
	// 创建mock服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code":0,"msg":"success","data":{}}`))
	}))
	defer server.Close()

	robot := NewGroupRobot(server.URL, "")
	response, err := robot.SendText("测试文本消息")

	if err != nil {
		t.Errorf("SendText不应该返回错误: %v", err)
	}

	if response.Code != 0 {
		t.Errorf("期望Code为0，实际为%d", response.Code)
	}
}

// TestSendRichWithoutUserId 测试SendRich不带用户ID
func TestSendRichWithoutUserId(t *testing.T) {
	// 创建mock服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code":0,"msg":"success","data":{}}`))
	}))
	defer server.Close()

	robot := NewGroupRobot(server.URL, "")
	response, err := robot.SendRich("测试标题", "测试内容", "")

	if err != nil {
		t.Errorf("SendRich不应该返回错误: %v", err)
	}

	if response.Code != 0 {
		t.Errorf("期望Code为0，实际为%d", response.Code)
	}
}

// TestSendRichWithUserId 测试SendRich带用户ID
func TestSendRichWithUserId(t *testing.T) {
	// 创建mock服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code":0,"msg":"success","data":{}}`))
	}))
	defer server.Close()

	robot := NewGroupRobot(server.URL, "")
	response, err := robot.SendRich("测试标题", "测试内容", "user123")

	if err != nil {
		t.Errorf("SendRich不应该返回错误: %v", err)
	}

	if response.Code != 0 {
		t.Errorf("期望Code为0，实际为%d", response.Code)
	}
}

// TestSendImage 测试SendImage快捷方法
func TestSendImage(t *testing.T) {
	// 创建mock服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code":0,"msg":"success","data":{}}`))
	}))
	defer server.Close()

	robot := NewGroupRobot(server.URL, "")
	response, err := robot.SendImage("img_test_key")

	if err != nil {
		t.Errorf("SendImage不应该返回错误: %v", err)
	}

	if response.Code != 0 {
		t.Errorf("期望Code为0，实际为%d", response.Code)
	}
}

// TestSendShare 测试SendShare快捷方法
func TestSendShare(t *testing.T) {
	// 创建mock服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code":0,"msg":"success","data":{}}`))
	}))
	defer server.Close()

	robot := NewGroupRobot(server.URL, "")
	response, err := robot.SendShare("chat_test_id")

	if err != nil {
		t.Errorf("SendShare不应该返回错误: %v", err)
	}

	if response.Code != 0 {
		t.Errorf("期望Code为0，实际为%d", response.Code)
	}
}

// TestSendCard 测试SendCard快捷方法
func TestSendCard(t *testing.T) {
	// 创建mock服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code":0,"msg":"success","data":{}}`))
	}))
	defer server.Close()

	robot := NewGroupRobot(server.URL, "")
	card := map[string]interface{}{
		"header": map[string]interface{}{
			"title": map[string]interface{}{
				"tag":     "plain_text",
				"content": "测试卡片",
			},
		},
	}
	response, err := robot.SendCard(card)

	if err != nil {
		t.Errorf("SendCard不应该返回错误: %v", err)
	}

	if response.Code != 0 {
		t.Errorf("期望Code为0，实际为%d", response.Code)
	}
}
