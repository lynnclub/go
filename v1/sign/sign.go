package sign

import (
	"fmt"
	"sort"
	"strings"

	"github.com/lynnclub/go/v1/algorithm"
)

func buildQueryString(params map[string]interface{}, secret string) string {
	keys := make([]string, 0, len(params))
	for key := range params {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	var builder strings.Builder
	for i, key := range keys {
		if i > 0 {
			builder.WriteByte('&')
		}
		builder.WriteString(key)
		builder.WriteByte('=')
		builder.WriteString(fmt.Sprintf("%v", params[key]))
	}
	builder.WriteString(secret)

	return builder.String()
}

// MD5 常规md5，get拼接
func MD5(params map[string]interface{}, secret string) string {
	return algorithm.MD5(buildQueryString(params, secret))
}

// SHA1 常规sha1，get拼接
func SHA1(params map[string]interface{}, secret string) string {
	return algorithm.SHA1(buildQueryString(params, secret))
}
