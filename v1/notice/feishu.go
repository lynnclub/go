package notice

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/lynnclub/go/v1/algorithm"
	"github.com/lynnclub/go/v1/bytedance/feishu"
	"github.com/lynnclub/go/v1/elasticsearch"
	"github.com/lynnclub/go/v1/safe"
)

type FeishuAlert struct {
	settings map[string]struct {
		webhook   string
		signKey   string
		userId    string
		kibanaUrl string
		esIndex   string
	}
	lastHashs []lastHash // 摘要
	mutex     sync.Mutex
}

type lastHash struct {
	hash string
	time time.Time
}

func (f *FeishuAlert) find(entry string) string {
	for name := range f.settings {
		if strings.Contains(entry, name) {
			return name
		}
	}

	return "default"
}

func (f *FeishuAlert) Send(log map[string]interface{}) {
	if log["level"].(int) < 2 {
		return
	}

	name := ""
	if log["url"].(string) == "" {
		name = f.find(log["command"].(string))
	} else {
		name = f.find(log["url"].(string))
	}

	newHash := lastHash{
		hash: algorithm.MD5(log["command"].(string) + log["message"].(string)),
		time: time.Now(),
	}

	for _, lastHash := range f.lastHashs {
		if lastHash.hash == newHash.hash {
			if lastHash.time.Add(10 * time.Minute).After(newHash.time) {
				return
			}
		}
	}

	f.mutex.Lock()

	f.lastHashs = append(f.lastHashs, newHash)
	if len(f.lastHashs) > 10 {
		f.lastHashs = f.lastHashs[1:]
	}

	f.mutex.Unlock()

	safe.Catch(func() {
		content := map[string]interface{}{
			"tag":  "text",
			"text": f.Format(log, f.settings[name].kibanaUrl, f.settings[name].esIndex),
		}
		feishu.NewGroupRobot(f.settings[name].webhook, f.settings[name].signKey).
			SendRich("", content, f.settings[name].userId)
	}, func(err any) {
		println(err)
	})
}

func (f *FeishuAlert) Format(log map[string]interface{}, kibanaUrl, esIndex string) string {
	querys := []string{
		elasticsearch.GetKuery("message", log["message"].(string)),
	}

	if trace := log["trace"].(string); trace != "" {
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
		log["env"],
		log["level_name"],
		log["datetime"],
		log["ip"],
		log["trace"],
		log["command"],
		log["message"],
		elasticsearch.GetKibanaUrl(kibanaUrl, esIndex, querys),
	)
}
