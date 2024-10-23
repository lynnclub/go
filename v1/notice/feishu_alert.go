package notice

import (
	"fmt"
	"sync"
	"time"

	"github.com/lynnclub/go/v1/algorithm"
	"github.com/lynnclub/go/v1/elasticsearch"
	"github.com/lynnclub/go/v1/safe"
)

// todo::支持多组配置，按command、path切换

var (
	lastHashs   []lastHash // 报警摘要
	alertMutex  sync.Mutex
	feiShuGroup *FeiShuGroup
)

type lastHash struct {
	hash string
	time time.Time
}

func FeishuAlert(full map[string]interface{}) {
	alertMutex.Lock()
	defer alertMutex.Unlock()

	newHash := lastHash{
		hash: algorithm.MD5(full["command"].(string) + full["message"].(string)),
		time: time.Now(),
	}

	for _, lastHash := range lastHashs {
		if lastHash.hash == newHash.hash {
			if lastHash.time.Add(10 * time.Minute).After(newHash.time) {
				return
			}
		}
	}

	lastHashs = append(lastHashs, newHash)

	// 超过最大容量时移除头部
	if len(lastHashs) > 10 {
		lastHashs = lastHashs[1:]
	}

	safe.Catch(func() {
		if feiShuGroup == nil {
			feiShuGroup = NewFeiShuGroup(
				l.feishu["webhook"],
				l.feishu["sign_key"],
				"",
			)
		}

		content := map[string]interface{}{
			"tag":  "text",
			"text": FeishuAlertFormat(full),
		}
		feiShuGroup.Send("", content, l.feishu["user_id"])
	}, func(err any) {
		println(err)
	})
}

func FeishuAlertFormat(full map[string]interface{}) string {
	querys := []string{
		elasticsearch.GetKuery("message", full["message"].(string)),
	}

	if trace := full["trace"].(string); trace != "" {
		querys = append(querys, elasticsearch.GetKuery("trace", trace))
	}

	return fmt.Sprintf(`环境：%s
级别：%s
时间：%s
IP：%s
追踪：%s
入口：%s
	
%s

%s

如有问题请尽快处理 []~(￣▽￣)~*`,
		full["env"],
		full["level_name"],
		full["datetime"],
		full["ip"],
		full["trace"],
		full["command"],
		full["message"],
		elasticsearch.GetKibanaUrl(l.feishu["kibana_url"], l.feishu["es_index"], querys),
	)
}
