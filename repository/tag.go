package repository

import (
	"fmt"
	"git-helper/utils"
	"regexp"
	"strings"
)

type Tag struct {
	Name    string `json:"name"`
	RefName string `json:"refName"`
	Type    string `json:"type"`
	Message string `json:"message"`
	Hash    string `json:"hash"`
	Time    string `json:"time"`
}

// Tags 获取该仓库所有的tag, git tag --sort=createdate
func (r *Repository) Tags() ([]Tag, error) {
	var tags []Tag
	f := fmt.Sprintf("--format=%s", `%(refname:strip=2)<||>%(refname)<||>%(objecttype)<||>%(subject)<||>%(objectname)<||>%(creatordate:relative)<|n|>`)

	out, err := utils.RunCmdByPath(r.Path, "git", "tag", "--sort=-creatordate", f)
	if err != nil {
		return nil, err
	}
	outs := strings.Split(out, "<|n|>")

	for _, tag := range outs {
		rows := strings.Split(strings.TrimSpace(tag), "<||>")
		if len(rows) != 6 {
			continue
		}
		tags = append(tags, Tag{
			Name:    rows[0],
			RefName: rows[1],
			Type:    rows[2],
			Message: rows[3],
			Hash:    rows[4],
			Time:    rows[5],
		})
	}
	return tags, nil
	//下面这种如果用户标签信息中含有"会导致失败
	//var tags []Tag
	//f := fmt.Sprintf("--format=%s", `{"name":"%(refname:short)","refName":"%(refname)","type":"%(objecttype)","message":"%(subject)","hash":"%(objectname)","time":"%(creatordate:relative)"},`)
	//
	//out, err := utils.RunCmdByPath(r.Path, "git", "tag", f)
	//if err != nil {
	//	return nil, err
	//}
	//jsonStr := fmt.Sprintf("[%s]", strings.TrimRight(strings.TrimSpace(out), ","))
	//fmt.Println(jsonStr)
	//err = json.Unmarshal([]byte(jsonStr), &tags)
	//if err != nil {
	//	return nil, err
	//}
	//
	//return tags, nil
}

// DelTag 删除指定tag, git tag -d / git push origin --delete tag
func (r *Repository) DelTag(tagName string, delRemote bool) (string, error) {

	if delRemote {
		out, err := utils.RunCmdByPath(r.Path, "git", "push", "origin", "--delete", "tag", tagName)
		if err != nil {
			return out, err
		}
	}

	out, err := utils.RunCmdByPath(r.Path, "git", "tag", "-d", tagName)
	if err != nil {
		return out, err
	}

	return out, nil
}

// git tag  v10.0 -m "version 1.0"
// CreateTag 创建标签, git tag
func (r *Repository) CreateTag(tag string, msg string) (string, error) {
	var out string
	var err error
	if msg == "" {
		out, err = utils.RunCmdByPath(r.Path, "git", "tag", tag)
	} else {
		out, err = utils.RunCmdByPath(r.Path, "git", "tag", tag, "-m", msg)
	}
	if err != nil {
		return out, err
	}
	return out, nil
}

// RemoteTags 获取远程标签 git ls-remote --tags --refs origin
func (r *Repository) RemoteTags() ([]string, error) {
	var tags []string
	isRemoteRepo, err := r.IsRemoteRepo()
	if err != nil || !isRemoteRepo {
		return tags, nil
	}
	out, err := utils.RunCmdByPath(r.Path, "git", "ls-remote", "--tags", "--refs", "origin")
	if err != nil {
		return tags, err
	}
	re := regexp.MustCompile(`refs/tags/(.+)\b`)
	matches := re.FindAllStringSubmatch(out, -1)
	for _, match := range matches {
		tags = append(tags, match[1])
	}
	return tags, nil
}

// git push origin v1.0
// PushTag 推送标签, git push origin
func (r *Repository) PushTag(name string) (string, error) {

	out, err := utils.RunCmdByPath(r.Path, "git", "push", "origin", name)
	if err != nil {
		return "", err
	}

	return out, nil
}

// CreateTagByCommitId 使用commit创建标签, git tag
func (r *Repository) CreateTagByCommitId(tag string, commitId string) (string, error) {
	//git tag v1.0 a867b4af
	out, err := utils.RunCmdByPath(r.Path, "git", "tag", tag, commitId)
	if err != nil {
		return out, err
	}
	return out, nil
}
