package icons

import "testing"

func TestGlossarySorted(t *testing.T) {
	entries := Glossary()
	if len(entries) == 0 {
		t.Fatal("Glossary() returned no entries")
	}
	for i := 1; i < len(entries); i++ {
		a, b := entries[i-1], entries[i]
		if a.Category > b.Category || (a.Category == b.Category && a.Name > b.Name) {
			t.Fatalf("Glossary not sorted at %d: %+v before %+v", i, a, b)
		}
	}
}

func TestDescribe(t *testing.T) {
	desc, ok := Describe(Bug)
	if !ok || desc == "" {
		t.Fatalf("Describe(Bug) = %q, ok=%v", desc, ok)
	}
	if _, ok := Describe(Name("nonexistent")); ok {
		t.Fatal("Describe unknown name should return false")
	}
}

func TestContextualGlossary(t *testing.T) {
	ctx := ContextualGlossary("open", "bug", 0)
	if len(ctx) < 3 {
		t.Fatalf("ContextualGlossary() = %d entries, want at least 3", len(ctx))
	}
	if len(ContextualGlossary("", "", -1)) != 0 {
		t.Fatal("ContextualGlossary with empty input should be empty")
	}
}

func TestGlossaryByCategory(t *testing.T) {
	grouped := GlossaryByCategory()
	if len(grouped) == 0 {
		t.Fatal("GlossaryByCategory returned empty map")
	}
	total := 0
	for _, entries := range grouped {
		total += len(entries)
	}
	if total != len(Glossary()) {
		t.Fatalf("category totals %d != Glossary len %d", total, len(Glossary()))
	}
}

func TestContextualGlossaryExtended(t *testing.T) {
	ctx := ContextualGlossary("review", "feature", 2)
	if len(ctx) < 3 {
		t.Fatalf("expected status+type+priority entries, got %d", len(ctx))
	}
}

func TestGlossaryGlyphUsesActiveSet(t *testing.T) {
	useSet(t, SetEmoji)
	if got := Get(Bug); got != "🐛" {
		t.Fatalf("emoji bug = %q", got)
	}
	useSet(t, SetNerd)
	if Get(Bug) == "🐛" {
		t.Fatal("nerd set should not return emoji bug")
	}
}
