package repository

import (
	"git-helper/utils"
	"strings"
)

// Repository 仓库对象
type Repository struct {
	Path string
}

// New 创建仓库对象
func New() *Repository {
	return &Repository{}
}

// SwitchRepository 切换仓库
func (r *Repository) SwitchRepository(path string) error {
	r.Path = path
	return nil
}

// GitPull 拉取仓库,命令为: git pull
func (r *Repository) GitPull() (string, error) {
	return utils.RunCmdByPath(r.Path, "git", "pull")
}

// GitPush 推仓库,命令为: git push
func (r *Repository) GitPush() (string, error) {
	return utils.RunCmdByPath(r.Path, "git", "push")
}

// GitRemoteUrl 获取仓库远程地址,  git remote get-url origin
func (r *Repository) GitRemoteUrl(path string) (string, error) {
	url, err := utils.RunCmdByPath(path, "git", "remote", "get-url", "origin")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(url), nil
}
