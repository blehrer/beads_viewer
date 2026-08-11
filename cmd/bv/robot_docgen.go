package main

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"github.com/Dicklesworthstone/beads_viewer/pkg/agents"
	"github.com/Dicklesworthstone/beads_viewer/pkg/ui"
	"github.com/Dicklesworthstone/beads_viewer/pkg/version"
)

func robotDocgenSections() []string {
	return []string{"agent-blurb", "robot-commands", "keybindings", "all"}
}

func isValidRobotDocgenSection(section string) bool {
	for _, name := range robotDocgenSections() {
		if section == name {
			return true
		}
	}
	return false
}

func generateRobotDocgen(section string) map[string]interface{} {
	section = strings.ToLower(strings.TrimSpace(section))
	if section == "" {
		section = "all"
	}

	result := map[string]interface{}{
		"generated_at": time.Now().UTC().Format(time.RFC3339),
		"version":      version.Version,
		"section":      section,
	}

	if !isValidRobotDocgenSection(section) {
		result["error"] = "Unknown section: " + section
		result["available_sections"] = robotDocgenSections()
		if suggestion := suggestClosest(section, robotDocgenSections()); suggestion != "" {
			result["did_you_mean"] = suggestion
			result["suggested_action"] = "Run `bv --robot-docgen " + suggestion + "`"
		} else {
			result["suggested_action"] = "Run `bv --robot-docgen all` or `bv --robot-docgen agent-blurb`"
		}
		return result
	}

	result["sections"] = robotDocgenSectionContent(section)
	return result
}

func robotDocgenSectionContent(section string) map[string]string {
	content := make(map[string]string)
	switch section {
	case "agent-blurb":
		content["agent-blurb"] = docgenAgentBlurbMarkdown()
	case "robot-commands":
		content["robot-commands"] = docgenRobotCommandsMarkdown()
	case "keybindings":
		content["keybindings"] = docgenKeybindingsMarkdown()
	case "all":
		content["agent-blurb"] = docgenAgentBlurbMarkdown()
		content["robot-commands"] = docgenRobotCommandsMarkdown()
		content["keybindings"] = docgenKeybindingsMarkdown()
	}
	return content
}

func docgenAgentBlurbMarkdown() string {
	return strings.TrimRight(agents.READMEAgentInstructions(), "\n") + "\n"
}

func docgenRobotCommandsMarkdown() string {
	docs := robotCommandDocs()
	names := make([]string, 0, len(docs))
	for name := range docs {
		names = append(names, name)
	}
	sort.Strings(names)

	var b strings.Builder
	b.WriteString("| Command | Output | Use Case |\n")
	b.WriteString("|---------|--------|----------|\n")
	for _, name := range names {
		doc := docs[name]
		flag := robotFlagExampleFormForCommand(name, doc.Flag)
		description := strings.TrimSpace(doc.Description)
		useCase := robotDocgenUseCase(name, doc)
		fmt.Fprintf(&b, "| `%s` | %s | %s |\n", flag, description, useCase)
	}
	return b.String()
}

func robotDocgenUseCase(name string, doc robotCommandDoc) string {
	switch {
	case doc.MutatesState:
		return "State mutation / feedback"
	case doc.NeedsBaseline:
		return "Baseline comparison"
	case doc.NeedsSprint:
		return "Sprint planning"
	case doc.NeedsGit && doc.NeedsIssues:
		return "Git + issue graph analysis"
	case doc.NeedsGit:
		return "Git history analysis"
	case !doc.NeedsIssues:
		return "Agent onboarding / metadata"
	case strings.Contains(name, "triage") || name == "robot-next":
		return "Work selection / triage"
	case strings.Contains(name, "label"):
		return "Label/domain analysis"
	case strings.Contains(name, "search"):
		return "Issue discovery"
	case strings.Contains(name, "graph"):
		return "Graph visualization & export"
	case strings.Contains(name, "correlation"):
		return "Commit correlation"
	case strings.Contains(name, "file"):
		return "File/bead linkage"
	case strings.Contains(name, "sprint") || strings.Contains(name, "forecast") || strings.Contains(name, "burndown") || strings.Contains(name, "capacity"):
		return "Sprint planning"
	default:
		return "Graph-aware analysis"
	}
}

func docgenKeybindingsMarkdown() string {
	docs := ui.GetKeyBindingDocs()
	var b strings.Builder
	b.WriteString("| Key | Description | Category | Context |\n")
	b.WriteString("|-----|-------------|----------|----------|\n")
	for _, doc := range docs {
		fmt.Fprintf(&b, "| `%s` | %s | %s | %s |\n", doc.Key, doc.Desc, doc.Category, doc.Context)
	}
	return b.String()
}

func writeRobotDocgen(out io.Writer, section string, asJSON bool) error {
	payload := generateRobotDocgen(section)
	if _, hasErr := payload["error"]; hasErr {
		if asJSON {
			if err := newRobotEncoder(out).Encode(payload); err != nil {
				return fmt.Errorf("encoding robot-docgen: %w", err)
			}
			return newReportedRobotHandlerExit(2)
		}
		if msg, _ := payload["error"].(string); msg != "" {
			return fmt.Errorf("%s", msg)
		}
		return newReportedRobotHandlerExit(2)
	}

	if asJSON {
		if err := newRobotEncoder(out).Encode(payload); err != nil {
			return fmt.Errorf("encoding robot-docgen: %w", err)
		}
		return nil
	}

	sections, _ := payload["sections"].(map[string]string)
	markdown := renderRobotDocgenMarkdown(payload["section"].(string), sections)
	if _, err := io.WriteString(out, markdown); err != nil {
		return fmt.Errorf("writing robot-docgen: %w", err)
	}
	return nil
}

func renderRobotDocgenMarkdown(section string, sections map[string]string) string {
	if section != "all" {
		return sections[section]
	}

	order := []string{"agent-blurb", "robot-commands", "keybindings"}
	var b strings.Builder
	for i, name := range order {
		if i > 0 {
			b.WriteByte('\n')
		}
		fmt.Fprintf(&b, "<!-- bv-docgen:%s -->\n", name)
		b.WriteString(strings.TrimRight(sections[name], "\n"))
		b.WriteByte('\n')
		fmt.Fprintf(&b, "<!-- /bv-docgen:%s -->\n", name)
	}
	return b.String()
}
