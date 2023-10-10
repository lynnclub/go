package response

import (
	"github.com/gin-gonic/gin"
	"github.com/lynnclub/go/v1/response/json_struct"
	"net/http"
	"time"
)

// JsonStruct Json结构接口
type JsonStruct interface {
	Set(status int, msg string, data interface{}, timestamp int64)
}

// JsonContext 内容
var JsonContext JsonStruct = &json_struct.Default{}

// Json 响应json
func Json(c *gin.Context, status int, msg string, data ...interface{}) {
	//if msg != "" {
	//	//多语言、语料编码
	//	if lang, ok := c.Get("Lang"); ok {
	//		msg = MsgI18n(msg, lang.(string))
	//	}
	//}

	var newData interface{}
	if len(data) == 1 {
		newData = data[0]
	} else {
		newData = data
	}

	JsonContext.Set(status, msg, newData, time.Now().Unix())

	// jsonp键名为callback
	c.JSONP(http.StatusOK, JsonContext)
	c.Abort()
}
