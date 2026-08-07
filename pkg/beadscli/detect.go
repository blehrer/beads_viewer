package beadscli

import (
	"os"
	"strings"

	"github.com/Dicklesworthstone/beads_viewer/pkg/loader"
)

const envTool = "BV_BEADS_CLI"

// DetectFromRepo selects br vs bd from BV_BEADS_CLI or the workspace under repoPath.
func DetectFromRepo(repoPath string) {
	if v := strings.TrimSpace(os.Getenv(envTool)); v != "" {
		SetTool(v)
		return
	}
	beadsDir, err := loader.GetBeadsDir(repoPath)
	if err != nil {
		SetTool("br")
		return
	}
	detectFromBeadsDir(beadsDir)
}

// DetectFromBeadsDir sets Tool to bd for Dolt-native workspaces, br otherwise.
func DetectFromBeadsDir(beadsDir string) {
	if v := strings.TrimSpace(os.Getenv(envTool)); v != "" {
		SetTool(v)
		return
	}
	detectFromBeadsDir(beadsDir)
}

func detectFromBeadsDir(beadsDir string) {
	if loader.IsBDWorkspace(beadsDir) {
		SetTool("bd")
		return
	}
	SetTool("br")
}

// ExportFlushCommand returns the shell command to flush DB state to JSONL.
func ExportFlushCommand() string {
	switch Tool() {
	case "bd":
		return "bd export --no-memories -o .beads/issues.jsonl"
	default:
		return "br sync --flush-only"
	}
}
