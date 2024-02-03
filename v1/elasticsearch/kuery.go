package elasticsearch

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

func GetKuery(field, value string) string {
	value = strings.ReplaceAll(value, `"`, "")

	if len(value) > 90 {
		value = value[:90]
		lastPos := checkIncompleteness(value)
		if lastPos > -1 {
			value = value[:lastPos]
		}
	}

	return fmt.Sprintf(`%s:"%s"`, field, value)
}

func GetKibanaUrl(url, index string, querys []string) string {
	risonA := map[string]interface{}{
		"index": index,
		"query": map[string]interface{}{
			"language": "kuery",
			"query":    strings.Join(querys, " AND "),
		},
	}

	risonG := map[string]interface{}{
		"time": map[string]string{
			"from": time.Now().AddDate(0, 0, -1).Format(time.RFC3339),
			"to":   time.Now().AddDate(0, 0, 1).Format(time.RFC3339),
		},
	}

	return fmt.Sprintf(url+"#/?_a=(%s)&_g=(%s)", toUrlRison(risonA), toUrlRison(risonG))
}

func checkIncompleteness(text string) int {
	lastPos := []int{
		strings.LastIndex(text, " "),
		strings.LastIndex(text, "\\"),
		strings.LastIndex(text, "::"),
	}

	max := lastPos[0]
	for _, num := range lastPos[1:] {
		if num > max {
			max = num
		}
	}

	return max
}

func toUrlRison(m map[string]interface{}) string {
	var rison string
	for key, value := range m {
		key = url.QueryEscape(key)
		switch v := value.(type) {
		case map[string]interface{}:
			rison += fmt.Sprintf("%s:(%s),", key, toUrlRison(v))
		case string:
			valueStr := strings.ReplaceAll(v, "'", "!'")
			rison += fmt.Sprintf("%s:'%s',", key, url.QueryEscape(fmt.Sprintf("%v", valueStr)))
		}
	}

	return strings.TrimRight(rison, ",")
}
