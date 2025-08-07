package main

import (
	"errors"
	"git-helper/utils"
	"github.com/atotto/clipboard"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"io/ioutil"
)

// Sha256 加密
func (a *App) Sha256(s string) string {
	return utils.Sha256(s)
}

// MessageDialog 原生消息弹窗
func (a *App) MessageDialog(title string, message string) {
	_, _ = runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:    runtime.InfoDialog,
		Title:   title,
		Message: message,
	})
}

// SaveJsonFile 保存 JSON 文件
func (a *App) SaveJsonFile(t string) error {

	if err := utils.RemoveFile(a.dataSaveJson); err != nil {
		return err
	}
	err := ioutil.WriteFile(a.dataSaveJson, []byte(t), 0666)

	if err != nil {
		return err
	}
	return nil
}

// ReadJsonFile 读取 JSON 文件
func (a *App) ReadJsonFile() (string, error) {

	b, err := ioutil.ReadFile(a.dataSaveJson)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Clipboard 写入剪切版
func (a *App) Clipboard(t string) error {
	return clipboard.WriteAll(t)
}

// IsGitRepository 是否为git仓库
func (a *App) IsGitRepository(path string) (bool, error) {
	if !utils.IsDir(path + "/.git") {
		return false, errors.New("not a git repository")
	}
	return true, nil
}
