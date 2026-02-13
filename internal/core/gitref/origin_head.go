package gitref

import "strings"

func ParseOriginHeadRef(ref string) (string, bool) {
	trimmed := strings.TrimSpace(ref)
	if !strings.HasPrefix(trimmed, "refs/remotes/origin/") {
		return "", false
	}
	branch := strings.TrimPrefix(trimmed, "refs/remotes/origin/")
	if branch == "" {
		return "", false
	}
	return branch, true
}

func FormatOriginTarget(branch string) (string, bool) {
	trimmed := strings.TrimSpace(branch)
	if trimmed == "" {
		return "", false
	}
	return "origin/" + trimmed, true
}
