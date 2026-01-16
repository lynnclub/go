package file

import (
	"os"
	"path/filepath"
	"testing"
)

// TestExists 测试文件存在性检查
func TestExists(t *testing.T) {
	// 创建临时目录
	tmpDir := t.TempDir()

	// 测试不存在的文件
	nonExistentFile := filepath.Join(tmpDir, "non_existent.txt")
	if Exists(nonExistentFile) {
		t.Error("不存在的文件应该返回false")
	}

	// 创建一个文件
	existingFile := filepath.Join(tmpDir, "existing.txt")
	err := os.WriteFile(existingFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 测试存在的文件
	if !Exists(existingFile) {
		t.Error("存在的文件应该返回true")
	}
}

// TestExistsWithDirectory 测试目录存在性检查
func TestExistsWithDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	// 测试存在的目录
	if !Exists(tmpDir) {
		t.Error("存在的目录应该返回true")
	}

	// 测试不存在的目录
	nonExistentDir := filepath.Join(tmpDir, "non_existent_dir")
	if Exists(nonExistentDir) {
		t.Error("不存在的目录应该返回false")
	}
}

// TestExistsWithEmptyFile 测试空文件
func TestExistsWithEmptyFile(t *testing.T) {
	tmpDir := t.TempDir()

	// 创建空文件
	emptyFile := filepath.Join(tmpDir, "empty.txt")
	err := os.WriteFile(emptyFile, []byte(""), 0644)
	if err != nil {
		t.Fatalf("创建空文件失败: %v", err)
	}

	// 空文件也应该返回true
	if !Exists(emptyFile) {
		t.Error("空文件应该返回true")
	}
}

// TestExistsWithSymlink 测试符号链接
func TestExistsWithSymlink(t *testing.T) {
	tmpDir := t.TempDir()

	// 创建原始文件
	originalFile := filepath.Join(tmpDir, "original.txt")
	err := os.WriteFile(originalFile, []byte("content"), 0644)
	if err != nil {
		t.Fatalf("创建原始文件失败: %v", err)
	}

	// 创建符号链接
	symlinkFile := filepath.Join(tmpDir, "symlink.txt")
	err = os.Symlink(originalFile, symlinkFile)
	if err != nil {
		t.Skipf("创建符号链接失败（可能不支持）: %v", err)
	}

	// 符号链接应该返回true
	if !Exists(symlinkFile) {
		t.Error("符号链接应该返回true")
	}
}

// TestExistsWithSpecialCharacters 测试特殊字符文件名
func TestExistsWithSpecialCharacters(t *testing.T) {
	tmpDir := t.TempDir()

	// 测试包含空格的文件名
	fileWithSpace := filepath.Join(tmpDir, "file with space.txt")
	err := os.WriteFile(fileWithSpace, []byte("content"), 0644)
	if err != nil {
		t.Fatalf("创建文件失败: %v", err)
	}

	if !Exists(fileWithSpace) {
		t.Error("包含空格的文件名应该正确处理")
	}

	// 测试包含中文的文件名
	fileWithChinese := filepath.Join(tmpDir, "中文文件.txt")
	err = os.WriteFile(fileWithChinese, []byte("内容"), 0644)
	if err != nil {
		t.Fatalf("创建中文文件失败: %v", err)
	}

	if !Exists(fileWithChinese) {
		t.Error("包含中文的文件名应该正确处理")
	}
}

// TestExistsWithRelativePath 测试相对路径
func TestExistsWithRelativePath(t *testing.T) {
	// 测试当前目录下的文件
	currentFile := "file.go"
	if !Exists(currentFile) {
		t.Error("当前包的file.go应该存在")
	}

	// 测试不存在的相对路径
	if Exists("./non_existent_relative.txt") {
		t.Error("不存在的相对路径文件应该返回false")
	}
}
