package repository

import (
	"errors"
	"git-helper/utils"
	"os/exec"
	sysRuntime "runtime"
	"strings"
)

// OpenTerminal 打开终端
func (r *Repository) OpenTerminal() error {

	var cmd *exec.Cmd
	switch sysRuntime.GOOS {
	// open -b com.apple.Terminal
	case "darwin":
		cmd = exec.Command("open", "-b", "com.apple.Terminal", r.Path)
	// x-terminal-emulator -e cd <path>;bash
	case "linux":
		cmd = exec.Command("x-terminal-emulator", "-e", "cd "+r.Path+";bash")
	case "windows":
		cmd = exec.Command("start", "cmd")
	// start cmd
	default:
		return errors.New("unsupported operating system")
	}
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

// OpenFileManage 打开文件管理器
func (r *Repository) OpenFileManage() error {
	var cmd *exec.Cmd
	switch sysRuntime.GOOS {
	case "darwin":
		cmd = exec.Command("open", r.Path)
	case "linux":
		cmd = exec.Command("xdg-open", r.Path)
	case "windows":
		cmd = exec.Command(`cmd`, `/c`, `explorer`, r.Path)
	default:
		return errors.New("unsupported operating system")
	}
	err := cmd.Start()
	if err != nil {
		return err
	}
	return nil
}

// RunCmdInRepository 在仓库目录中执行命令
func (r *Repository) RunCmdInRepository(cmd string, arg []string) (string, error) {
	return utils.RunCmdByPath(r.Path, cmd, arg...)
}

// IsRemoteRepo 仓库是否有远程地址
func (r *Repository) IsRemoteRepo() (bool, error) {

	out, err := utils.RunCmdByPath(r.Path, "git", "remote", "-v")
	if err != nil {
		return false, err
	}

	if strings.TrimSpace(out) == "" {
		return false, nil
	}
	return true, nil
}
