package requeststate

import (
	"strings"
	"testing"
)

func TestAppendSectionEntryReusesTheExistingHeading(t *testing.T) {
	entry := "- [2026-10-01] blocked on \"x\" — cleared by y"
	for _, testCase := range []struct {
		name, newline, headingSuffix string
	}{
		{"CRLF file", "\r\n", ""},
		{"heading with trailing spaces", "\n", "  "},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			document := strings.Join([]string{
				"---", "id: REQ-1", "---", "# Title", "",
				"## Blocked" + testCase.headingSuffix, "",
				"- earlier entry", "",
				"## Plan", "", "text", "",
			}, testCase.newline)
			updated := string(appendSectionEntry([]byte(document), "Blocked", entry))
			if got := strings.Count(updated, "## Blocked"); got != 1 {
				t.Fatalf("want exactly one Blocked heading, got %d in:\n%q", got, updated)
			}
			if entryAt, planAt := strings.Index(updated, entry), strings.Index(updated, "## Plan"); entryAt < strings.Index(updated, "## Blocked") || entryAt > planAt {
				t.Fatalf("entry must land inside the existing Blocked section, got:\n%q", updated)
			}
		})
	}
}

func TestAppendSectionEntryIgnoresAHeadingInsideACodeFence(t *testing.T) {
	document := "---\nid: REQ-1\n---\n# Title\n\n```\n## Blocked\n```\n"
	updated := string(appendSectionEntry([]byte(document), "Blocked", "- new entry"))
	if got := strings.Count(updated, "\n## Blocked"); got != 2 {
		t.Fatalf("a fenced example heading is not a section, so a real one must be added after it; got:\n%q", updated)
	}
}
