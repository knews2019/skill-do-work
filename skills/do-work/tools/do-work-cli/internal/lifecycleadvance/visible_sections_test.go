package lifecycleadvance

import "testing"

func TestAdvanceIgnoresExampleAndCommentHeadings(t *testing.T) {
	examples := "```markdown\n## Plan\n```\n<!--\n## Scope\n-->\n"
	sections, reason := advanceSections([]byte(examples))
	if reason != "" || len(sections) != 0 {
		t.Fatalf("hidden headings became evidence: %#v, %s", sections, reason)
	}
	sections, reason = advanceSections([]byte(examples + "## Plan\nPlanning not required.\n"))
	if reason != "" || !hasSection(sections, "Plan") || hasSection(sections, "Scope") {
		t.Fatalf("example caused duplicate or false evidence: %#v, %s", sections, reason)
	}
}
