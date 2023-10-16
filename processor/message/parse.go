package message

import (
	"regexp"
	"strings"
)

func GetAtMembersFromRawContent(rawContent string, skip int) (members []string) {
	// 如果占用了占位符，那么就不解析
	reg := regexp.MustCompile(`ƒ@!(.*?)∆`)
	if rawMatches := reg.FindAllStringSubmatch(rawContent, -1); len(rawMatches) > 0 {
		return []string{}
	}

	// 替换占位符
	raw := strings.ReplaceAll(rawContent, `\u003c`, "ƒ")
	raw = strings.ReplaceAll(raw, `\u003e`, "∆")

	// 解析 @ 的用户
	matches := reg.FindAllStringSubmatch(raw, -1)
	if matches == nil || len(matches) <= skip {
		return []string{}
	} else {
		// skip 个数的 @ 不解析
		for _, s := range matches[skip:] {
			members = append(members, s[1])
		}
		return members
	}
}
