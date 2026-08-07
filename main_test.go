package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileHash(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "a.txt")
	os.WriteFile(p, []byte("hello"), 0644)
	h, err := fileHash(p)
	if err != nil {
		t.Fatal(err)
	}
	// sha256("hello") 的已知值
	want := "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"
	if h != want {
		t.Errorf("哈希不对: %q", h)
	}
}

func TestHashDir(t *testing.T) {
	dir := t.TempDir()
	os.Mkdir(filepath.Join(dir, "sub"), 0755)
	os.WriteFile(filepath.Join(dir, "x.txt"), []byte("a"), 0644)
	os.WriteFile(filepath.Join(dir, "sub", "y.txt"), []byte("b"), 0755)
	m, err := hashDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(m) != 2 {
		t.Errorf("期望 2 个文件, 得到 %d", len(m))
	}
	if _, ok := m["x.txt"]; !ok {
		t.Error("缺少 x.txt 的哈希")
	}
	if _, ok := m["sub/y.txt"]; !ok {
		t.Error("嵌套文件应带相对路径")
	}
}

func TestHashDirMissing(t *testing.T) {
	if _, err := hashDir("不存在的路径"); err == nil {
		t.Error("目录不存在应报错")
	}
}
