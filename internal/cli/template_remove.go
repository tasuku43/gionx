package cli

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/tasuku43/kra/internal/app/templatecmd"
	"github.com/tasuku43/kra/internal/infra/paths"
)

func (c *CLI) runTemplateRemove(args []string) int {
	templateName := ""
	promptedName := false
	for len(args) > 0 && strings.HasPrefix(args[0], "-") {
		switch args[0] {
		case "-h", "--help", "help":
			c.printTemplateRemoveUsage(c.Out)
			return exitOK
		case "--name":
			if len(args) < 2 {
				fmt.Fprintln(c.Err, "--name requires a value")
				c.printTemplateRemoveUsage(c.Err)
				return exitUsage
			}
			templateName = strings.TrimSpace(args[1])
			args = args[2:]
		default:
			fmt.Fprintf(c.Err, "unknown flag for template remove: %q\n", args[0])
			c.printTemplateRemoveUsage(c.Err)
			return exitUsage
		}
	}
	if len(args) > 1 {
		fmt.Fprintf(c.Err, "unexpected args for template remove: %q\n", strings.Join(args, " "))
		c.printTemplateRemoveUsage(c.Err)
		return exitUsage
	}
	if len(args) == 1 {
		if templateName != "" {
			fmt.Fprintln(c.Err, "--name cannot be combined with positional <template>")
			c.printTemplateRemoveUsage(c.Err)
			return exitUsage
		}
		templateName = strings.TrimSpace(args[0])
	}
	if templateName == "" {
		useColorErr := writerSupportsColor(c.Err)
		line, err := c.promptLine(renderTemplateNamePrompt(useColorErr))
		if err != nil {
			fmt.Fprintf(c.Err, "read template name: %v\n", err)
			return exitError
		}
		templateName = strings.TrimSpace(line)
		promptedName = true
	}
	if err := validateWorkspaceTemplateName(templateName); err != nil {
		fmt.Fprintln(c.Err, err.Error())
		c.printTemplateRemoveUsage(c.Err)
		return exitUsage
	}

	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(c.Err, "get working dir: %v\n", err)
		return exitError
	}
	root, err := paths.ResolveExistingRoot(wd)
	if err != nil {
		fmt.Fprintf(c.Err, "resolve KRA_ROOT: %v\n", err)
		return exitError
	}
	if err := c.ensureDebugLog(root, "template-remove"); err != nil {
		fmt.Fprintf(c.Err, "enable debug logging: %v\n", err)
	}
	c.debugf("run template remove name=%q", templateName)

	templatePath, err := removeWorkspaceTemplate(root, templateName)
	if err != nil {
		printTemplateCreateError(c.Err, writerSupportsColor(c.Err), templateName, "", err)
		return exitError
	}
	commitSHA, err := templatecmd.CommitRemove(context.Background(), root, templateName)
	if err != nil {
		printTemplateCreateError(c.Err, writerSupportsColor(c.Err), templateName, "", fmt.Errorf("commit remove change: %w", err))
		return exitError
	}

	if promptedName {
		printTemplateCreateInputs(c.Err, writerSupportsColor(c.Err), templateName, "")
	}
	useColorOut := writerSupportsColor(c.Out)
	lines := []string{
		styleSuccess("Removed 1 / 1", useColorOut),
		fmt.Sprintf("%s %s", styleSuccess("✔", useColorOut), templateName),
		styleMuted(fmt.Sprintf("path: %s", templatePath), useColorOut),
	}
	if strings.TrimSpace(commitSHA) != "" {
		lines = append(lines, styleMuted(fmt.Sprintf("commit: %s", shortCommitSHA(commitSHA)), useColorOut))
	}
	if promptedName {
		body := make([]string, 0, len(lines))
		for _, line := range lines {
			body = append(body, fmt.Sprintf("%s%s", uiIndent, line))
		}
		printSection(c.Out, renderResultTitle(useColorOut), body, sectionRenderOptions{
			blankAfterHeading: false,
			trailingBlank:     true,
		})
	} else {
		printResultSection(c.Out, useColorOut, lines...)
	}
	return exitOK
}

func removeWorkspaceTemplate(root string, name string) (string, error) {
	templatePath := workspaceTemplatePath(root, name)
	st, err := os.Stat(templatePath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("template %q not found: %s", name, templatePath)
		}
		return "", fmt.Errorf("stat template: %w", err)
	}
	if !st.IsDir() {
		return "", fmt.Errorf("template path is not a directory: %s", templatePath)
	}
	if err := os.RemoveAll(templatePath); err != nil {
		return "", fmt.Errorf("remove template path: %w", err)
	}
	return templatePath, nil
}
