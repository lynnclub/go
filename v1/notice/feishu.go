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
	options   map[string]Option
	lastHashs []lastHash // 摘要
	mutex     sync.Mutex
}

type lastHash struct {
	hash string
	time time.Time
}

type Option struct {
	Webhook   string `json:"webhook"`
	SignKey   string `json:"sign_key"`
	UserId    string `json:"user_id"`
	KibanaUrl string `json:"kibana_url"`
	EsIndex   string `json:"es_index"`
}

func (f *FeishuAlert) Add(name string, option Option) {
	if option.Webhook == "" {
		panic("Option webhook empty " + name)
	}

	f.options[name] = option
}

func (f *FeishuAlert) AddMap(name string, setting map[string]interface{}) {
	option := Option{
		Webhook:   setting["webhook"].(string),
		SignKey:   setting["sign_key"].(string),
		UserId:    setting["user_id"].(string),
		KibanaUrl: setting["kibana_url"].(string),
		EsIndex:   setting["es_index"].(string),
	}

	f.Add(name, option)
}

func (f *FeishuAlert) AddMapBatch(batch map[string]interface{}) {
	for name, setting := range batch {
		f.AddMap(name, setting.(map[string]interface{}))
	}
}

func (f *FeishuAlert) FindOption(entry, defaultName string) string {
	for name := range f.options {
		if strings.Contains(entry, name) {
			return name
		}
	}

	return defaultName
}

func (f *FeishuAlert) Send(log map[string]interface{}) {
	if log["level"].(int) < 2 {
		return
	}

	name := ""
	if log["url"].(string) == "" {
		name = f.FindOption(log["command"].(string), "default_command")
	} else {
		name = f.FindOption(log["url"].(string), "default_api")
	}

	option, ok := f.options[name]
	if !ok {
		return
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
			"text": f.Format(log, option.KibanaUrl, option.EsIndex),
		}
		feishu.NewGroupRobot(option.Webhook, option.SignKey).SendRich("", content, option.UserId)
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
