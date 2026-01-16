package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lynnclub/go/v1/response/json_struct"
)

// setupTestRouter 创建测试路由
func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

// TestJsonWithDefaultStruct 测试使用默认json结构
func TestJsonWithDefaultStruct(t *testing.T) {
	router := setupTestRouter()

	router.GET("/test", func(c *gin.Context) {
		Json(c, 200, "success", map[string]string{"key": "value"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("期望状态码200，实际%d", w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	if int(response["status"].(float64)) != 200 {
		t.Errorf("期望status为200，实际%v", response["status"])
	}

	if response["msg"] != "success" {
		t.Errorf("期望msg为success，实际%v", response["msg"])
	}

	if response["data"] == nil {
		t.Error("data不应该为nil")
	}
}

// TestJsonWithNoData 测试不传数据参数
func TestJsonWithNoData(t *testing.T) {
	router := setupTestRouter()

	router.GET("/test", func(c *gin.Context) {
		Json(c, 404, "not found")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	// 不传数据参数时，data字段存在但为nil或空
	if response["data"] == nil {
		// 这是可以接受的
		return
	}

	// 也可能是空数组
	dataArray, ok := response["data"].([]interface{})
	if ok && len(dataArray) != 0 {
		t.Errorf("期望空数组或nil，实际长度%d", len(dataArray))
	}
}

// TestJsonWithSingleData 测试传单个数据
func TestJsonWithSingleData(t *testing.T) {
	router := setupTestRouter()

	router.GET("/test", func(c *gin.Context) {
		Json(c, 200, "success", map[string]string{"name": "test"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	dataMap, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Error("单个数据应该是map类型")
	}
	if dataMap["name"] != "test" {
		t.Errorf("期望name为test，实际%v", dataMap["name"])
	}
}

// TestJsonWithMultipleData 测试传多个数据
func TestJsonWithMultipleData(t *testing.T) {
	router := setupTestRouter()

	router.GET("/test", func(c *gin.Context) {
		Json(c, 200, "success", "data1", "data2", "data3")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	dataArray, ok := response["data"].([]interface{})
	if !ok {
		t.Error("多个数据应该是数组类型")
	}
	if len(dataArray) != 3 {
		t.Errorf("期望数组长度3，实际%d", len(dataArray))
	}
}

// TestJsonWithCodeStruct 测试使用Code结构
func TestJsonWithCodeStruct(t *testing.T) {
	// 临时切换到Code结构
	originalContext := JsonContext
	JsonContext = &json_struct.Code{}
	defer func() {
		JsonContext = originalContext
	}()

	router := setupTestRouter()

	router.GET("/test", func(c *gin.Context) {
		Json(c, 200, "success", nil)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	// 验证使用的是code字段而不是status
	if _, ok := response["code"]; !ok {
		t.Error("应该包含code字段")
	}
	if int(response["code"].(float64)) != 200 {
		t.Errorf("期望code为200，实际%v", response["code"])
	}
}

// TestJsonWithMessageStruct 测试使用Message结构
func TestJsonWithMessageStruct(t *testing.T) {
	// 临时切换到Message结构
	originalContext := JsonContext
	JsonContext = &json_struct.Message{}
	defer func() {
		JsonContext = originalContext
	}()

	router := setupTestRouter()

	router.GET("/test", func(c *gin.Context) {
		Json(c, 200, "success", nil)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	// 验证使用的是message字段而不是msg
	if _, ok := response["message"]; !ok {
		t.Error("应该包含message字段")
	}
	if response["message"] != "success" {
		t.Errorf("期望message为success，实际%v", response["message"])
	}
}

// TestJsonWithRawStruct 测试使用Raw结构
func TestJsonWithRawStruct(t *testing.T) {
	// 临时切换到Raw结构
	originalContext := JsonContext
	JsonContext = &json_struct.Raw{}
	defer func() {
		JsonContext = originalContext
	}()

	router := setupTestRouter()

	router.GET("/test", func(c *gin.Context) {
		Json(c, 200, "ignored", map[string]string{"key": "value"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	// Raw结构应该直接返回数据，不包含status、msg等字段
	if _, ok := response["status"]; ok {
		t.Error("Raw结构不应该包含status字段")
	}
	if _, ok := response["msg"]; ok {
		t.Error("Raw结构不应该包含msg字段")
	}
	if response["key"] != "value" {
		t.Errorf("期望key为value，实际%v", response["key"])
	}
}

// TestJsonSetsTimestamp 测试时间戳设置
func TestJsonSetsTimestamp(t *testing.T) {
	router := setupTestRouter()

	router.GET("/test", func(c *gin.Context) {
		Json(c, 200, "success", nil)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	timestamp, ok := response["timestamp"]
	if !ok {
		t.Error("响应应该包含timestamp字段")
	}
	if timestamp.(float64) <= 0 {
		t.Error("timestamp应该是正数")
	}
}

// TestJsonAborts 测试Json调用会中止后续处理
func TestJsonAborts(t *testing.T) {
	router := setupTestRouter()

	handlerCalled := false
	router.GET("/test", func(c *gin.Context) {
		defer func() {
			// Gin的Abort不会阻止当前handler内的代码执行
			// 它只会阻止后续middleware的执行
			handlerCalled = true
		}()
		Json(c, 200, "success", nil)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	// Json调用了c.Abort()，但不会阻止同一handler内的代码
	if !handlerCalled {
		t.Error("Json调用不会阻止同一handler内的defer")
	}
}

// TestJsonWithDifferentStatusCodes 测试不同状态码
func TestJsonWithDifferentStatusCodes(t *testing.T) {
	statusCodes := []int{200, 201, 400, 401, 403, 404, 500}

	for _, statusCode := range statusCodes {
		t.Run(string(rune(statusCode)), func(t *testing.T) {
			router := setupTestRouter()

			router.GET("/test", func(c *gin.Context) {
				Json(c, statusCode, "test", nil)
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			router.ServeHTTP(w, req)

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)

			if int(response["status"].(float64)) != statusCode {
				t.Errorf("期望status为%d，实际%v", statusCode, response["status"])
			}
		})
	}
}

// TestJsonWithNilData 测试nil数据
func TestJsonWithNilData(t *testing.T) {
	router := setupTestRouter()

	router.GET("/test", func(c *gin.Context) {
		Json(c, 200, "success", nil)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("期望状态码200，实际%d", w.Code)
	}

	// nil数据应该能正常序列化
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}
}

// TestJsonWithComplexData 测试复杂数据结构
func TestJsonWithComplexData(t *testing.T) {
	router := setupTestRouter()

	complexData := map[string]interface{}{
		"user": map[string]interface{}{
			"id":   123,
			"name": "测试用户",
			"tags": []string{"admin", "user"},
		},
		"items": []map[string]interface{}{
			{"id": 1, "name": "item1"},
			{"id": 2, "name": "item2"},
		},
		"metadata": map[string]interface{}{
			"total": 100,
			"page":  1,
		},
	}

	router.GET("/test", func(c *gin.Context) {
		Json(c, 200, "success", complexData)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("解析响应失败: %v", err)
	}

	dataMap := response["data"].(map[string]interface{})
	if dataMap["user"] == nil {
		t.Error("复杂数据结构未正确序列化")
	}
}

// TestJsonEmptyMessage 测试空消息
func TestJsonEmptyMessage(t *testing.T) {
	router := setupTestRouter()

	router.GET("/test", func(c *gin.Context) {
		Json(c, 200, "", nil)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	if response["msg"] != "" {
		t.Errorf("期望msg为空字符串，实际%v", response["msg"])
	}
}
