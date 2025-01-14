package elasticsearch

import (
	"sync"
)

var (
	pool    = &sync.Map{}             //实例池
	mutex   sync.Mutex                //互斥锁
	options = make(map[string]Option) //配置池
)

type Option struct {
	Address  []string `json:"address"`  //地址，字符串数组
	Username string   `json:"username"` //密码，默认空
	Password string   `json:"password"` //密码，默认空
}

func Add(name string, option Option) {
	if len(option.Address) == 0 {
		panic("Option address array empty " + name)
	}

	options[name] = option
}

func AddMap(name string, setting map[string]interface{}) {
	addressStrings := make([]string, 0)
	for _, v := range setting["address"].([]interface{}) {
		addressStrings = append(addressStrings, v.(string))
	}

	option := Option{
		Address: addressStrings,
	}

	if username, ok := setting["username"]; ok {
		option.Username = username.(string)
	}
	if password, ok := setting["password"]; ok {
		option.Password = password.(string)
	}

	Add(name, option)
}

func AddMapBatch(batch map[string]interface{}) {
	for name, setting := range batch {
		AddMap(name, setting.(map[string]interface{}))
	}
}
