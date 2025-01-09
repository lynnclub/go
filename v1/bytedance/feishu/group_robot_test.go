package feishu

import "testing"

func TestNewGroupRobot(t *testing.T) {
	content := map[string]interface{}{
		"tag":  "text",
		"text": "test",
	}
	NewGroupRobot("https://open.feishu.cn/open-apis/bot/v2/hook/xxx", "").SendRich("", content, "")
}
