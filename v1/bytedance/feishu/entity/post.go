package entity

// Post 富文本
type Post map[string]PostData

// PostData 富文本数据
type PostData struct {
	Title   string                     `json:"title"`
	Content [][]map[string]interface{} `json:"content"`
}
