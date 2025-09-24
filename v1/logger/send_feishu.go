package logger

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/lynnclub/go/v1/algorithm"
	"github.com/lynnclub/go/v1/array"
	"github.com/lynnclub/go/v1/bytedance/feishu"
	"github.com/lynnclub/go/v1/elasticsearch"
)

var Feishu *FeishuAlert

func NewFeishu(options map[string]interface{}) *FeishuAlert {
	instance := &FeishuAlert{}
	instance.AddMapBatch(options)

	return instance
}

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
	Levels    []string `json:"levels"`
	Webhook   string   `json:"webhook"`
	SignKey   string   `json:"sign_key"`
	UserId    string   `json:"user_id"`
	KibanaUrl string   `json:"kibana_url"`
	EsIndex   string   `json:"es_index"`
}

func (f *FeishuAlert) Add(name string, option Option) {
	if option.Webhook == "" {
		panic("Option webhook empty " + name)
	}

	if f.options == nil {
		f.options = make(map[string]Option)
	}

	f.options[name] = option
}

func (f *FeishuAlert) AddMap(name string, setting map[string]interface{}) {
	levels := []string{}
	if tmps, ok := setting["levels"].([]string); ok {
		for _, level := range tmps {
			levels = append(levels, strings.ToUpper(level))
		}
	}

	option := Option{
		Levels:    levels,
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

func (f *FeishuAlert) FindOption(levelName string, entry, defaultName string) string {
	for name, option := range f.options {
		if strings.Contains(entry, name) && (len(option.Levels) == 0 || array.In(option.Levels, levelName)) {
			return name
		}
	}

	return defaultName
}

func (f *FeishuAlert) Send(log LogEntry) {
	if log.Level <= 200 {
		return
	}

	name := ""
	if log.URL == "" {
		name = f.FindOption(log.LevelName, log.Command, "default_command")
	} else {
		name = f.FindOption(log.LevelName, log.URL, "default_api")
	}

	option, ok := f.options[name]
	if !ok {
		return
	}

	keyword := ""
	if traces, ok := log.Extra.([]string); ok && len(traces) > 0 {
		keyword = traces[0]
	} else {
		keyword = log.Command + log.Message
	}

	if keyword != "" {
		newHash := lastHash{
			hash: algorithm.MD5(keyword),
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
	}

	content := f.Format(log, option.KibanaUrl, option.EsIndex)
	feishu.NewGroupRobot(option.Webhook, option.SignKey).SendRich("", content, option.UserId)
}

func (f *FeishuAlert) Format(log LogEntry, kibanaUrl, esIndex string) string {
	querys := []string{
		elasticsearch.GetKuery("message", log.Message),
	}

	traceParam := ""
	if log.Trace == "" {
		traceParam = elasticsearch.GetKuery("command", log.Command)
		log.Trace = log.Command
	} else {
		traceParam = elasticsearch.GetKuery("trace", log.Trace)
	}

	querys = append(querys, traceParam)

	return fmt.Sprintf(`环境：%s
级别：%s
时间：%s
IP：%s
追踪：%s
入口：%s

%s

详情
%s
链路
%s

如有问题请尽快处理 []~(￣▽￣)~*`,
		log.Env,
		log.LevelName,
		log.Datetime,
		log.IP,
		log.Trace,
		log.Command,
		log.Message,
		elasticsearch.GetKibanaUrl(kibanaUrl, esIndex, querys),
		elasticsearch.GetKibanaUrl(kibanaUrl, esIndex, []string{traceParam}),
	)
}
