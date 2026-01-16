package file

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"testing"
)

// TestCSVWriterOverwrite 测试覆盖模式写入CSV
func TestCSVWriterOverwrite(t *testing.T) {
	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "test_overwrite")

	// 第一次写入
	writer := CSVWriter(filename, true, "Name", "Age", "City")
	writer.Write([]string{"Alice", "30", "Beijing"})
	writer.Flush()

	if err := writer.Error(); err != nil {
		t.Errorf("写入CSV失败: %v", err)
	}

	// 验证文件存在
	if !Exists(filename + ".csv") {
		t.Error("CSV文件应该已创建")
	}

	// 读取并验证内容
	content, err := os.ReadFile(filename + ".csv")
	if err != nil {
		t.Fatalf("读取CSV文件失败: %v", err)
	}

	expected := "Name,Age,City\nAlice,30,Beijing\n"
	if string(content) != expected {
		t.Errorf("期望内容:\n%s\n实际内容:\n%s", expected, string(content))
	}

	// 第二次写入（覆盖）
	writer2 := CSVWriter(filename, true, "Name", "Age", "City")
	writer2.Write([]string{"Bob", "25", "Shanghai"})
	writer2.Flush()

	// 验证内容被覆盖
	content2, err := os.ReadFile(filename + ".csv")
	if err != nil {
		t.Fatalf("读取CSV文件失败: %v", err)
	}

	expected2 := "Name,Age,City\nBob,25,Shanghai\n"
	if string(content2) != expected2 {
		t.Errorf("期望覆盖后内容:\n%s\n实际内容:\n%s", expected2, string(content2))
	}
}

// TestCSVWriterAppend 测试追加模式写入CSV
func TestCSVWriterAppend(t *testing.T) {
	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "test_append")

	// 第一次写入（创建文件）
	writer1 := CSVWriter(filename, false, "Name", "Age")
	writer1.Write([]string{"Alice", "30"})
	writer1.Flush()

	// 第二次写入（追加）
	writer2 := CSVWriter(filename, false, "Name", "Age")
	writer2.Write([]string{"Bob", "25"})
	writer2.Flush()

	// 验证内容
	content, err := os.ReadFile(filename + ".csv")
	if err != nil {
		t.Fatalf("读取CSV文件失败: %v", err)
	}

	expected := "Name,Age\nAlice,30\nBob,25\n"
	if string(content) != expected {
		t.Errorf("期望内容:\n%s\n实际内容:\n%s", expected, string(content))
	}
}

// TestCSVWriterMultipleRows 测试写入多行
func TestCSVWriterMultipleRows(t *testing.T) {
	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "test_multiple")

	writer := CSVWriter(filename, true, "ID", "Name", "Score")

	rows := [][]string{
		{"1", "Alice", "95"},
		{"2", "Bob", "87"},
		{"3", "Charlie", "92"},
		{"4", "David", "88"},
	}

	for _, row := range rows {
		writer.Write(row)
	}
	writer.Flush()

	if err := writer.Error(); err != nil {
		t.Errorf("写入多行CSV失败: %v", err)
	}

	// 读取并验证
	file, err := os.Open(filename + ".csv")
	if err != nil {
		t.Fatalf("打开CSV文件失败: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("读取CSV记录失败: %v", err)
	}

	// 应该有5行（1个header + 4行数据）
	if len(records) != 5 {
		t.Errorf("期望5行记录，实际%d行", len(records))
	}

	// 验证header
	if records[0][0] != "ID" || records[0][1] != "Name" || records[0][2] != "Score" {
		t.Errorf("Header不正确: %v", records[0])
	}

	// 验证第一行数据
	if records[1][0] != "1" || records[1][1] != "Alice" || records[1][2] != "95" {
		t.Errorf("第一行数据不正确: %v", records[1])
	}
}

// TestCSVWriterEmptyHeaders 测试空header
func TestCSVWriterEmptyHeaders(t *testing.T) {
	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "test_empty_headers")

	writer := CSVWriter(filename, true)
	writer.Write([]string{"data1", "data2"})
	writer.Flush()

	content, err := os.ReadFile(filename + ".csv")
	if err != nil {
		t.Fatalf("读取CSV文件失败: %v", err)
	}

	// 空header应该写入空行
	expected := "\ndata1,data2\n"
	if string(content) != expected {
		t.Errorf("期望内容:\n%s\n实际内容:\n%s", expected, string(content))
	}
}

// TestCSVWriterWithSpecialCharacters 测试特殊字符
func TestCSVWriterWithSpecialCharacters(t *testing.T) {
	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "test_special")

	writer := CSVWriter(filename, true, "姓名", "年龄", "城市")
	writer.Write([]string{"张三", "30", "北京"})
	writer.Write([]string{"李四", "25", "上海"})
	writer.Write([]string{"王五,带逗号", "28", "深圳"})
	writer.Flush()

	// 读取并验证
	file, err := os.Open(filename + ".csv")
	if err != nil {
		t.Fatalf("打开CSV文件失败: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("读取CSV记录失败: %v", err)
	}

	// 验证中文header
	if records[0][0] != "姓名" {
		t.Errorf("中文header不正确: %v", records[0])
	}

	// 验证中文数据
	if records[1][0] != "张三" {
		t.Errorf("中文数据不正确: %v", records[1])
	}

	// 验证包含逗号的数据被正确处理
	if records[3][0] != "王五,带逗号" {
		t.Errorf("包含逗号的数据不正确: %v", records[3])
	}
}

// TestIsEmpty 测试isEmpty函数
func TestIsEmpty(t *testing.T) {
	tmpDir := t.TempDir()

	// 创建空文件
	emptyFile := filepath.Join(tmpDir, "empty.txt")
	file, err := os.Create(emptyFile)
	if err != nil {
		t.Fatalf("创建空文件失败: %v", err)
	}

	if !isEmpty(file) {
		t.Error("空文件应该返回true")
	}
	file.Close()

	// 创建非空文件
	nonEmptyFile := filepath.Join(tmpDir, "non_empty.txt")
	file2, err := os.Create(nonEmptyFile)
	if err != nil {
		t.Fatalf("创建文件失败: %v", err)
	}
	file2.WriteString("content")
	file2.Close()

	file2, _ = os.Open(nonEmptyFile)
	if isEmpty(file2) {
		t.Error("非空文件应该返回false")
	}
	file2.Close()
}

// TestCSVWriterAppendToExistingFile 测试追加到已存在的文件
func TestCSVWriterAppendToExistingFile(t *testing.T) {
	tmpDir := t.TempDir()
	filename := filepath.Join(tmpDir, "test_existing")

	// 先创建一个文件
	os.WriteFile(filename+".csv", []byte("Name,Age\nAlice,30\n"), 0644)

	// 追加模式写入（不应该重复写header）
	writer := CSVWriter(filename, false, "Name", "Age")
	writer.Write([]string{"Bob", "25"})
	writer.Flush()

	content, err := os.ReadFile(filename + ".csv")
	if err != nil {
		t.Fatalf("读取CSV文件失败: %v", err)
	}

	// 应该只有一个header
	expected := "Name,Age\nAlice,30\nBob,25\n"
	if string(content) != expected {
		t.Errorf("期望内容:\n%s\n实际内容:\n%s", expected, string(content))
	}
}
