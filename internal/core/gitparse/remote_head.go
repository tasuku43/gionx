package gitparse

import "strings"

func ParseRemoteHeadSymref(lsRemoteOutput string) (branch string, hash string) {
	lines := strings.Split(strings.TrimSpace(lsRemoteOutput), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "ref: ") && strings.HasSuffix(line, "\tHEAD") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				ref := strings.TrimPrefix(parts[1], "refs/heads/")
				if ref != "" {
					branch = ref
				}
			}
			continue
		}
		if strings.HasSuffix(line, "\tHEAD") {
			fields := strings.Fields(line)
			if len(fields) >= 1 {
				hash = fields[0]
			}
		}
	}
	return branch, hash
}
