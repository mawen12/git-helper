package repository

import (
	"errors"
	"git-helper/utils"
	"regexp"
	"strconv"
	"strings"
)

type MergeKind int

const conflictRegexp = `\+[<>=]{7}\s\.our[\w\W]*?\+[<>=]{7}\s\.their`

const (
	Invalid MergeKind = iota
	Conflicted
	Clean
)

type MergeResult struct {
	Kind  MergeKind `json:"kind"`
	Count int       `json:"count"`
}

// PreMergeResult 模拟合并,基于 git merge-base 找出共同祖先, 基于 git merge-tree 模拟三方合并
func (r *Repository) PreMergeResult(currentHash, MergeHash string) (MergeResult, error) {

	var invalidResult = MergeResult{
		Kind:  Invalid,
		Count: 0,
	}

	// 找出两个分支的共同祖先, h 为共同祖先的 commitHash
	h, err := utils.RunCmdByPath(r.Path, "git", "merge-base", currentHash, MergeHash)
	if err != nil {
		return invalidResult, err
	}
	baseHash := strings.TrimSpace(h)

	// 模拟一次三方合并的结果,不修改工作区,不产生提交 baseHash 两个分支共同祖先, currentHash 分支1, MergeHash 要合并进来的分支
	out, err := utils.RunCmdByPath(r.Path, "git", "merge-tree", baseHash, currentHash, MergeHash)
	if err != nil {
		return invalidResult, err
	}

	reg := regexp.MustCompile(conflictRegexp)
	matches := reg.FindAllString(out, -1)

	if len(matches) != 0 {
		return MergeResult{Kind: Conflicted, Count: len(matches)}, nil
	}
	out, err = utils.RunCmdByPath(r.Path, "git", "rev-list", "--left-right", "--count", currentHash+"..."+MergeHash)
	if err != nil {
		return invalidResult, err
	}

	c := strings.Split(out, "\t")

	if len(c) != 2 {
		return invalidResult, errors.New("rev-list out err")
	}
	i, err := strconv.ParseInt(strings.TrimSpace(c[1]), 10, 64)
	if err != nil {
		return invalidResult, err
	}
	return MergeResult{Kind: Clean, Count: int(i)}, nil
}

// MergeRebase 变基,将其平移到另一个分支的后面,使提交历史变得更线性, 更整洁, git rebase
func (r *Repository) MergeRebase(ourBranch, theirBranch string) (string, error) {

	ok, err := r.SwitchBranch(ourBranch)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", errors.New("switch branch error")
	}
	out, err := utils.RunCmdByPath(r.Path, "git", "rebase", theirBranch)

	if err != nil {
		return "", err
	}
	return out, nil
}

// MergeCommit 合并,保留两个分支的历史,产生一个新的合并提交, git merge
func (r *Repository) MergeCommit(ourBranch, theirBranch string) (string, error) {
	ok, err := r.SwitchBranch(ourBranch)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", errors.New("switch branch error")
	}
	out, err := utils.RunCmdByPath(r.Path, "git", "merge", theirBranch)

	if err != nil {
		return "", err
	}
	return out, nil
}

// MergeSquash 压缩合并,仅保留一个提交,内容等于之前提交总和,不保留原始提交历史, git merge --squash
func (r *Repository) MergeSquash(ourBranch, theirBranch string) (string, error) {
	ok, err := r.SwitchBranch(ourBranch)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", errors.New("switch branch error")
	}
	out, err := utils.RunCmdByPath(r.Path, "git", "merge", "--squash", theirBranch)

	if err != nil {
		return "", err
	}
	return out, nil
}
