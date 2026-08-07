package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// 算单个文件的 sha256。
func fileHash(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// 遍历目录，给每个文件算哈希，返回相对路径->哈希的映射。
// 跳过目录本身，符号链接不跟进（避免环）。
func hashDir(root string) (map[string]string, error) {
	result := map[string]string{}
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			rel = path
		}
		// 统一成斜杠，Windows 下 Rel 会给反斜杠，输出保持一致好看
		rel = filepath.ToSlash(rel)
		h, err := fileHash(path)
		if err != nil {
			return err
		}
		result[rel] = h
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("用法: go-hashdir <目录>")
		return
	}
	if args[0] == "-h" || args[0] == "--help" {
		fmt.Println("go-hashdir 给目录下每个文件算 sha256")
		return
	}
	m, err := hashDir(args[0])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		fmt.Printf("%s  %s\n", m[k][:12], k)
	}
	// 顺手打个总数
	fmt.Fprintf(os.Stderr, "共 %d 个文件\n", len(keys))
	_ = strings.TrimSpace
}
