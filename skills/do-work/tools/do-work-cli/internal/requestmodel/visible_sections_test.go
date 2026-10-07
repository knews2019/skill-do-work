package requestmodel

import (
	"strings"
	"testing"
)

func TestVisibleSectionsRetainOffsetsAndProtectUnclosedRegions(t *testing.T) {
	for _, newline := range []string{"\n", "\r\n"} {
		body := strings.ReplaceAll("# Request\n:::note\n## Example\n:::\n<!-- hidden -->\n## Plan\nreal plan\n### Detail\n  ## Requirements <!-- retained -->\nkeep\n<!--\n## Hidden\n", "\n", newline)
		sections := VisibleSections([]byte(body))
		if len(sections) != 2 || sections[0].Name != "Plan" || sections[1].Name != "Requirements <!-- retained -->" {
			t.Fatalf("wrong visible sections: %#v", sections)
		}
		for index, want := range []string{"## Plan\nreal plan\n### Detail\n", "  ## Requirements <!-- retained -->\nkeep\n"} {
			got := body[sections[index].Start:sections[index].End]
			if got != strings.ReplaceAll(want, "\n", newline) {
				t.Fatalf("section offsets changed bytes: %q", got)
			}
		}
	}
}

func visibleSectionNames(body string) []string {
	names := []string{}
	for _, section := range VisibleSections([]byte(body)) {
		names = append(names, section.Name)
	}
	return names
}

// REQ-635 F1: recover --take-over deleted user text because a heading indented
// by a tab or by four or more spaces was read as a section boundary. CommonMark
// admits an ATX heading or a fence only after zero to three spaces.
func TestVisibleSectionsApplyTheZeroToThreeSpaceIndentRule(t *testing.T) {
	cases := []struct {
		body string
		want []string
	}{
		{"# Request\nExample:\n\n    ## Plan\n    user sample\n\nMUST keep.\n## Timing\nt\n", []string{"Timing"}},
		{"# Request\n\t## Plan\np\n", []string{}},
		{"# Request\n \t## Plan\np\n", []string{}},
		{"# Request\n     ## Plan\np\n", []string{}},
		{"# Request\n ## Plan\np\n", []string{"Plan"}},
		{"# Request\n   ## Plan\np\n", []string{"Plan"}},
		// An indented code block line is not a fence, so it hides nothing.
		{"# Request\n    ```\n## Plan\np\n", []string{"Plan"}},
		{"# Request\n\t```\n## Plan\np\n", []string{"Plan"}},
		// A closer indented four spaces is code inside the fence, not its end.
		{"# Request\n```\n    ```\n## Hidden\n```\n## Plan\np\n", []string{"Plan"}},
	}
	for _, newline := range []string{"\n", "\r\n"} {
		for _, testCase := range cases {
			body := strings.ReplaceAll(testCase.body, "\n", newline)
			if got := visibleSectionNames(body); strings.Join(got, "|") != strings.Join(testCase.want, "|") {
				t.Errorf("%q: got sections %q, want %q", body, got, testCase.want)
			}
		}
	}
}

// REQ-635 F5: a fence opener's trailing text is its info string, so a "<!--"
// there must not open a comment that never closes and hides the Timing section.
// REQ-635 D-01: a "<!--" inside an inline code span is literal text.
func TestVisibleSectionsReadFenceOpenersAndCodeSpansBeforeComments(t *testing.T) {
	cases := []struct {
		body string
		want []string
	}{
		{"# Request\n``` <!--\n## X\n```\n## Plan\np\n## Timing\nt\n", []string{"Plan", "Timing"}},
		{"# Request\nUse `<!--` here.\n## Plan\np\n", []string{"Plan"}},
		{"# Request\nUse ``a ` <!-- b`` here.\n## Plan\np\n", []string{"Plan"}},
		{"# Request\nUse `<!--` then <!-- real\n## Hidden\n-->\n## Plan\np\n", []string{"Plan"}},
		// An unmatched backtick is literal, so the comment still opens.
		{"# Request\nUse ` then <!--\n## Plan\np\n", []string{}},
		// Inside a comment backticks mean nothing: the close is still found.
		{"# Request\n<!--\n`-->`\n## Plan\np\n", []string{"Plan"}},
	}
	for _, newline := range []string{"\n", "\r\n"} {
		for _, testCase := range cases {
			body := strings.ReplaceAll(testCase.body, "\n", newline)
			if got := visibleSectionNames(body); strings.Join(got, "|") != strings.Join(testCase.want, "|") {
				t.Errorf("%q: got sections %q, want %q", body, got, testCase.want)
			}
		}
	}
}

// An unclosed comment after a heading is protected by design: the section ends
// where the comment begins, so a destructive writer removes none of it.
func TestVisibleSectionsKeepAnUnclosedCommentHeadingZeroLength(t *testing.T) {
	for _, newline := range []string{"\n", "\r\n"} {
		body := strings.ReplaceAll("# Request\n## Plan <!-- note\nkept\n", "\n", newline)
		sections := VisibleSections([]byte(body))
		if len(sections) != 1 || sections[0].Name != "Plan <!-- note" || sections[0].Start != sections[0].End {
			t.Fatalf("unclosed comment section changed: %#v", sections)
		}
	}
}
