package requestmodel

import (
	"strings"
	"testing"
)

func TestVisibleSectionsRetainOffsetsAndProtectUnclosedRegions(t *testing.T) {
	for _, newline := range []string{"\n", "\r\n"} {
		body := strings.ReplaceAll("# Request\n:::note\n## Example\n:::\n<!-- hidden -->\n## Plan\nreal plan\n### Detail\n## Requirements\nkeep\n<!--\n## Hidden\n", "\n", newline)
		sections := VisibleSections([]byte(body))
		if len(sections) != 2 || sections[0].Name != "Plan" || sections[1].Name != "Requirements" {
			t.Fatalf("wrong visible sections: %#v", sections)
		}
		for index, want := range []string{"## Plan\nreal plan\n### Detail\n", "## Requirements\nkeep\n"} {
			got := body[sections[index].Start:sections[index].End]
			if got != strings.ReplaceAll(want, "\n", newline) {
				t.Fatalf("section offsets changed bytes: %q", got)
			}
		}
	}
}
