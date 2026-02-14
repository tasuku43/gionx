package cli

import (
	"fmt"
	"os"
	"path/filepath"
)

type bootstrapSkillpackFile struct {
	relativePath string
	content      string
}

func defaultAgentSkillpackFiles() []bootstrapSkillpackFile {
	return []bootstrapSkillpackFile{
		{
			relativePath: ".kra-skillpack.yaml",
			content: `version: "v1"
pack: "kra-insight-capture"
skills:
  - flow-insight-capture
`,
		},
		{
			relativePath: "flow-insight-capture/SKILL.md",
			content: `---
name: flow-insight-capture
description: Use this skill only for approved insight capture; it should not prescribe how work is executed.
---

# Flow: Insight Capture

## Goal

Leave reusable "pebbles" only when they are explicitly approved.
Do not control or prescribe the user's working style.

## Trigger

Use this skill only when a high-value reusable insight appears.

## Rules

1. Propose capture in conversation first.
2. Persist only when user explicitly approves.
3. Write with:
   kra ws insight add --id <workspace-id> --ticket <ticket> --session-id <session-id> --what "<summary>" --context "<context>" --why "<why>" --next "<next>" --tag <tag> ... --approved
4. Never write when approval is missing.
`,
		},
	}
}

func ensureBootstrapDefaultSkillpack(skillsRoot string, result *bootstrapAgentSkillsResult) error {
	for _, file := range defaultAgentSkillpackFiles() {
		path := filepath.Join(skillsRoot, filepath.FromSlash(file.relativePath))
		if err := ensureBootstrapSkillpackFile(path, file.content, result); err != nil {
			return err
		}
	}
	return nil
}

func ensureBootstrapSkillpackFile(path string, content string, result *bootstrapAgentSkillsResult) error {
	parent := filepath.Dir(path)
	parentInfo, err := os.Stat(parent)
	if err == nil {
		if !parentInfo.IsDir() {
			appendBootstrapConflict(result, parent, "exists and is not a directory")
			return nil
		}
	} else if os.IsNotExist(err) {
		if err := os.MkdirAll(parent, 0o755); err != nil {
			return fmt.Errorf("create %s: %w", parent, err)
		}
		appendUniquePath(&result.Created, parent)
	} else {
		return fmt.Errorf("stat %s: %w", parent, err)
	}

	info, err := os.Lstat(path)
	if err == nil {
		if info.Mode().IsRegular() {
			// Preserve user/tool-managed edits; bootstrap remains non-destructive.
			appendUniquePath(&result.Skipped, path)
			return nil
		}
		appendBootstrapConflict(result, path, "exists and is not a regular file")
		return nil
	}
	if !os.IsNotExist(err) {
		return fmt.Errorf("lstat %s: %w", path, err)
	}

	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	appendUniquePath(&result.Created, path)
	return nil
}

func appendUniquePath(items *[]string, path string) {
	for _, item := range *items {
		if item == path {
			return
		}
	}
	*items = append(*items, path)
}
