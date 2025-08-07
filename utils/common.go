package utils

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// IsDir 是否为目录
func IsDir(path string) bool {

	if info, err := os.Stat(path); err == nil {
		return info.IsDir()
	}
	return false
}

// RunCmdByPath 在指定目录下执行命令
func RunCmdByPath(path string, cmdName string, arg ...string) (string, error) {
	if !IsDir(path) {
		return "", errors.New("Can't find repository path:" + path)
	}

	cmd := exec.Command(cmdName, arg...)
	cmd.Dir = path

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.New(string(out))
	}
	return string(out), nil
}

// Sha256 将字符串进行sha256加密
func Sha256(s string) string {
	h := sha256.Sum256([]byte(s))
	hashString := fmt.Sprintf("%x", h)
	return hashString
}

// 检查文件是否存在
func FileIsExisted(filename string) bool {
	existed := true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		existed = false
	}
	return existed
}

// 移除文件
func RemoveFile(file string) error {
	if FileIsExisted(file) {
		return os.Remove(file)
	}
	return nil
}
