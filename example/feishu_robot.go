package main

import (
	"log"
	"time"

	"github.com/lynnclub/go/v1/bytedance/feishu"
)

func main() {
	// 初始化机器人客户端
	webhook := "" // 填入你的webhook地址
	signKey := "" // 可选：签名密钥，没有可留空
	userId := ""

	if webhook == "" {
		log.Fatal("请设置webhook地址")
	}

	robot := feishu.NewGroupRobot(webhook, signKey)

	// 1. 文本消息
	response1, err1 := robot.SendText("Hello 飞书！这是一条测试的文本消息")
	if err1 != nil {
		log.Printf("第1个消息发送失败: %v", err1)
	} else {
		log.Printf("第1个消息发送成功: code=%d, msg=%s", response1.Code, response1.Msg)
	}
	time.Sleep(1 * time.Second) // 延时1秒，避免频率限制

	// 2. 富文本消息（带标题）
	response2, err2 := robot.SendRich("系统通知", "测试完成，一切正常运行", "")
	if err2 != nil {
		log.Printf("第2个消息发送失败: %v", err2)
	} else {
		log.Printf("第2个消息发送成功: code=%d, msg=%s", response2.Code, response2.Msg)
	}
	time.Sleep(1 * time.Second)

	// 3. 富文本消息（@用户）
	response3, err3 := robot.SendRich("紧急通知", "数据库连接测试！", userId)
	if err3 != nil {
		log.Printf("第3个消息发送失败: %v", err3)
	} else {
		log.Printf("第3个消息发送成功: code=%d, msg=%s", response3.Code, response3.Msg)
	}
	time.Sleep(1 * time.Second)

	// 4. 图片消息
	// response4, err4 := robot.SendImage("img_v2_图片key")
	// if err4 != nil {
	// 	log.Printf("第4个消息发送失败: %v", err4)
	// } else {
	// 	log.Printf("第4个消息发送成功: code=%d, msg=%s", response4.Code, response4.Msg)
	// }
	// time.Sleep(1 * time.Second)

	// 5. 分享群名片
	// response5, err5 := robot.SendShare("oc_群聊ID")
	// if err5 != nil {
	// 	log.Printf("第5个消息发送失败: %v", err5)
	// } else {
	// 	log.Printf("第5个消息发送成功: code=%d, msg=%s", response5.Code, response5.Msg)
	// }
	//time.Sleep(1 * time.Second)

	// 6. 复杂富文本（文本+链接+@用户）
	elements := []feishu.RichElement{
		{Type: "text", Text: "🔥 项目部署完成！\n"},
		{Type: "text", Text: "📊 监控面板："},
		{Type: "a", Text: "点击查看", Href: "https://monitor.example.com"},
		{Type: "text", Text: " "},
		{Type: "at", UserID: userId},
	}
	request := &feishu.GroupRobotRequest{}
	request.BuildAdvancedRichMessage("部署通知", elements)
	response6, err6 := robot.Send(request)
	if err6 != nil {
		log.Printf("第6个消息发送失败: %v", err6)
	} else {
		log.Printf("第6个消息发送成功: code=%d, msg=%s", response6.Code, response6.Msg)
	}
	time.Sleep(1 * time.Second)

	// 7. 交互式卡片
	card := map[string]any{
		"header": map[string]any{
			"title": map[string]any{
				"tag":     "plain_text",
				"content": "系统告警测试",
			},
			"template": "red",
		},
		"elements": []map[string]any{
			{
				"tag": "div",
				"text": map[string]any{
					"tag":     "lark_md",
					"content": "**CPU**: 85% | **内存**: 78%",
				},
			},
			{
				"tag": "action",
				"actions": []map[string]any{
					{
						"tag": "button",
						"text": map[string]any{
							"tag":     "plain_text",
							"content": "查看详情",
						},
						"url":  "https://dashboard.example.com",
						"type": "primary",
					},
				},
			},
		},
	}
	response7, err7 := robot.SendCard(card)
	if err7 != nil {
		log.Printf("第7个消息发送失败: %v", err7)
	} else {
		log.Printf("第7个消息发送成功: code=%d, msg=%s", response7.Code, response7.Msg)
	}

	log.Println("所有消息发送完成!")
}
