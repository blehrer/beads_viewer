package ui

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// Markdown builder functions must use *IconMD helpers, not lipgloss-styled icons.
var (
	markdownBuilderNameRE = regexp.MustCompile(`^build.*Markdown$`)

	forbiddenStyledIconCalls = map[string]struct{}{
		"GetPriorityIcon":       {},
		"GetStatusIcon":         {},
		"RenderPriorityIcon":    {},
		"RenderStatusDot":       {},
		"RenderTriageScoreIcon": {},
	}

	skipMarkdownPolicyScan = map[string]struct{}{
		"styles.go":  {},
		"helpers.go": {},
	}
)

func TestMarkdownBuilders_NoStyledIcons(t *testing.T) {
	fset := token.NewFileSet()

	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatalf("read pkg/ui: %v", err)
	}

	var violations []string

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".go") || strings.HasSuffix(name, "_test.go") {
			continue
		}
		if _, skip := skipMarkdownPolicyScan[name]; skip {
			continue
		}

		path := filepath.Join(".", name)
		file, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			t.Fatalf("parse %s: %v", path, err)
		}

		for _, decl := range file.Decls {
			fn, ok := decl.(*ast.FuncDecl)
			if !ok || fn.Body == nil || !isMarkdownBuilderFunc(fn.Name.Name) {
				continue
			}

			ast.Inspect(fn.Body, func(n ast.Node) bool {
				call, ok := n.(*ast.CallExpr)
				if !ok {
					return true
				}
				callee, ok := styledIconCallee(call.Fun)
				if !ok {
					return true
				}
				pos := fset.Position(call.Pos())
				violations = append(violations, strings.Join([]string{
					pos.Filename,
					fn.Name.Name,
					callee,
					pos.String(),
				}, "\t"))
				return true
			})
		}
	}

	if len(violations) > 0 {
		t.Fatalf("markdown builders must not call lipgloss-styled icon helpers (use *IconMD instead):\n  %s",
			strings.Join(violations, "\n  "))
	}
}

func isMarkdownBuilderFunc(name string) bool {
	return name == "renderBeadHistoryMD" || markdownBuilderNameRE.MatchString(name)
}

func styledIconCallee(expr ast.Expr) (string, bool) {
	switch e := expr.(type) {
	case *ast.Ident:
		if _, forbidden := forbiddenStyledIconCalls[e.Name]; forbidden {
			return e.Name, true
		}
	case *ast.SelectorExpr:
		if _, forbidden := forbiddenStyledIconCalls[e.Sel.Name]; forbidden {
			return e.Sel.Name, true
		}
	}
	return "", false
}
