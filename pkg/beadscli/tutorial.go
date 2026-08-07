package beadscli

import "strings"

// TutorialLine rewrites br-prefixed example commands for the active CLI.
func TutorialLine(s string) string {
	tool := Tool()
	if tool == "bd" {
		s = strings.ReplaceAll(s, "br sync", ExportFlushCommand())
	}
	s = strings.ReplaceAll(s, "br ", tool+" ")
	s = strings.ReplaceAll(s, "`br`", "`"+tool+"`")
	return s
}
