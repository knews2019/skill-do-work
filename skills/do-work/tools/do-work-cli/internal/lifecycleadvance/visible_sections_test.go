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

// REQ-645: a two-space "## Plan" under a bullet is the user's sample, so it is
// neither Plan evidence nor a duplicate of the generated column-0 Plan.
func TestAdvanceIgnoresAListNestedTwoSpaceIndentedHeading(t *testing.T) {
	sample := "- Example:\n  ## Plan\n  user sample\n\n"
	sections, reason := advanceSections([]byte(sample))
	if reason != "" || hasSection(sections, "Plan") {
		t.Fatalf("indented sample became Plan evidence: %#v, %s", sections, reason)
	}
	sections, reason = advanceSections([]byte(sample + "## Plan\nPlanning not required.\n"))
	if reason != "" || !hasSection(sections, "Plan") {
		t.Fatalf("indented sample duplicated the generated Plan: %#v, %s", sections, reason)
	}
}
