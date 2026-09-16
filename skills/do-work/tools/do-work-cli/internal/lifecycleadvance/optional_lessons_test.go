package lifecycleadvance

import (
	"fmt"
	"strings"
	"testing"
)

func TestAdvanceAllowsOnlyRouteAToOmitLessons(t *testing.T) {
	for _, route := range []string{"A", "B", "C"} {
		for _, orientation := range []bool{false, true} {
			t.Run(fmt.Sprintf("%s/orientation=%t", route, orientation), func(t *testing.T) {
				root := t.TempDir()
				body := strings.ReplaceAll(routeCBodyThrough("Review"), "Route C.", "Route "+route+".")
				if route != "C" {
					body = strings.ReplaceAll(body, "Specific plan.", "Planning not required.")
				}
				if route == "A" {
					body = strings.Replace(body, "## Exploration\n\nFound patterns.\n\n", "", 1)
					body = strings.Replace(body, "## Scope\n\n**Files I will touch:**\n- `owned.go`\n\n", "", 1)
				}
				if orientation {
					body += "## Orientation\n\nStraightforward correction.\n"
				}
				writeAdvanceRequest(t, root, "working", "REQ-703", "claimed", "route: "+route+"\nplanning_at: 2026-09-04T12:00:00Z\nwrite_set: [owned.go]\nestimate:\n  p50_active_minutes: 5\n", body)
				result, _ := runAdvanceJSON(t, root, "REQ-703")
				if orientation && route != "A" {
					if result.Outcome != "refused" || len(result.Findings) == 0 || result.Findings[0].Code != "ADVANCE-EVIDENCE-MISSING" {
						t.Fatalf("required lessons accepted as absent: %+v", result)
					}
					return
				}
				if result.Outcome != "success" || result.Advance == nil {
					t.Fatalf("valid route evidence refused: %+v", result)
				}
				if orientation {
					if result.Advance.Phase != "finalize" {
						t.Fatalf("Route A cannot finalize without lessons: %+v", result.Advance)
					}
				} else {
					want := "Lessons Learned"
					if route == "A" {
						want = "Orientation"
					}
					if len(result.Advance.MissingEvidence) != 1 || result.Advance.MissingEvidence[0].Section != want {
						t.Fatalf("missing evidence=%+v, want %s", result.Advance.MissingEvidence, want)
					}
				}
			})
		}
	}
}
