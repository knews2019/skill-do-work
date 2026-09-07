package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// A character limit still overflows narrow cards with wide glyphs. Measure the
// generated card and retain the full condition for readers using its tooltip.
func TestBrowserBehaviorBlockedBadgeFitsItsCard(t *testing.T) {
	lookupBrowserForBehaviorProbe(t)
	condition := strings.Repeat("Waiting WWWWW for external approval ", 8)
	repoRoot := writeVerifyFixture(t, []verifyFixtureFile{{
		RelativePath: "do-work/queue/REQ-701-blocked-badge.md",
		Content:      "---\nid: REQ-701\ntitle: Blocked badge fixture\nstatus: blocked\nblocked_by: " + condition + "\n---\n",
	}})
	board, err := buildBoard(repoRoot, time.Now().UTC(), defaultRecentWindow, nil)
	if err != nil {
		t.Fatal(err)
	}
	siteDirectory := t.TempDir()
	if err := generateStaticSite(siteDirectory, board); err != nil {
		t.Fatal(err)
	}
	indexBytes, err := os.ReadFile(filepath.Join(siteDirectory, "index.html"))
	if err != nil {
		t.Fatal(err)
	}
	session := startTrustedInputBrowserSession(t, "blocked badge fit", siteDirectory, string(indexBytes))
	defer session.closeBrowserSession()
	session.waitForPageCondition(t, "a rendered blocked card", `document.querySelector('.req-card .badge-blocked')`)
	for _, width := range []int{200, 320} {
		var measured struct {
			CardWidth  float64 `json:"cardWidth"`
			CardRight  float64 `json:"cardRight"`
			BadgeRight float64 `json:"badgeRight"`
			Truncated  bool    `json:"truncated"`
			Tooltip    string  `json:"tooltip"`
			Value      string  `json:"value"`
		}
		expression := fmt.Sprintf(`(function () {
  var badge = document.querySelector('.req-card .badge-blocked');
  var card = badge.closest('.req-card');
  card.style.width = '%dpx';
  var value = badge.querySelector('.badge-blocked-value');
  return {
    cardWidth: card.getBoundingClientRect().width,
    cardRight: card.getBoundingClientRect().right,
    badgeRight: badge.getBoundingClientRect().right,
    truncated: !!value && value.clientWidth > 0 && value.scrollWidth > value.clientWidth,
    tooltip: badge.title,
    value: value ? value.textContent : ''
  };
})()`, width)
		session.decodeResult(t, "blocked badge geometry", session.evaluateInPage(t, expression), &measured)
		if measured.CardWidth <= 0 || measured.BadgeRight > measured.CardRight || !measured.Truncated {
			t.Errorf("blocked badge does not fit a %dpx card: %+v", width, measured)
		}
		if measured.Value != strings.TrimSpace(condition) || !strings.Contains(measured.Tooltip, strings.TrimSpace(condition)) {
			t.Errorf("blocked condition was lost at %dpx: %+v", width, measured)
		}
	}
}
