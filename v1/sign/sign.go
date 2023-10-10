package sign

import (
	"fmt"
	"github.com/lynnclub/go/v1/algorithm"
	"sort"
	"strings"
)

// MD5 常规md5，get拼接
func MD5(params map[string]interface{}, secret string) string {
	var keys []string
	for key := range params {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	var list []string
	for _, key := range keys {
		list = append(list, key+"="+fmt.Sprintf("%v", params[key]))
	}

	return algorithm.MD5(strings.Join(list, "&") + secret)
}

// SHA1 常规sha1，get拼接
func SHA1(params map[string]interface{}, secret string) string {
	var keys []string
	for key := range params {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	var list []string
	for _, key := range keys {
		list = append(list, key+"="+fmt.Sprintf("%v", params[key]))
	}

	return algorithm.SHA1(strings.Join(list, "&") + secret)
}
