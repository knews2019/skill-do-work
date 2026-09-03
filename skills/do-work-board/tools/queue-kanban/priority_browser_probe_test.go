package main

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

type priorityBadgeMeasurement struct {
	RequestID     string
	Text          string
	Color         string
	Background    string
	Width         float64
	Height        float64
	Contained     bool
	OverlapsTitle bool
}

type priorityBoardMeasurement struct {
	Href    string
	Browser string
	Scheme  string
	Ready   []string
	Waiting []string
	Badges  []priorityBadgeMeasurement
}

func TestBrowserBehaviorPriorityOrderAndBadges(t *testing.T) {
	lookupBrowserForBehaviorProbe(t)
	repoRoot := writeVerifyFixture(t, []verifyFixtureFile{
		{"do-work/queue/REQ-401-next.md", priorityBrowserFixture("REQ-401", "Default next card", "", "")},
		{"do-work/queue/REQ-402-later.md", priorityBrowserFixture("REQ-402", "Later card with a long title that exercises wrapping next to badges", "priority: later\n", "")},
		{"do-work/queue/REQ-403-now.md", priorityBrowserFixture("REQ-403", "Now card", "priority: now\n", "")},
		{"do-work/queue/REQ-404-waiting-later.md", priorityBrowserFixture("REQ-404", "Waiting later", "priority: later\n", "depends_on: [REQ-999]\n")},
		{"do-work/queue/REQ-405-waiting-now.md", priorityBrowserFixture("REQ-405", "Waiting now", "priority: now\n", "depends_on: [REQ-998]\n")},
	})
	board, buildError := buildBoard(repoRoot, time.Now().UTC(), defaultRecentWindow, nil)
	if buildError != nil {
		t.Fatal(buildError)
	}
	siteDirectory := t.TempDir()
	if generateError := generateStaticSite(siteDirectory, board); generateError != nil {
		t.Fatal(generateError)
	}
	indexBytes, readError := os.ReadFile(filepath.Join(siteDirectory, "index.html"))
	if readError != nil {
		t.Fatal(readError)
	}
	session := startTrustedInputBrowserSession(t, "priority order and badges", siteDirectory, string(indexBytes))
	defer session.closeBrowserSession()
	session.waitForPageCondition(t, "priority cards", `document.querySelectorAll('[data-cards="pending"] .req-card').length === 5`)
	measurePriorityPage(t, session)
}

func priorityBrowserFixture(requestID, title, priorityLine, dependencyLine string) string {
	return "---\nid: " + requestID + "\ntitle: " + title + "\nstatus: pending\n" + priorityLine + dependencyLine + "---\n"
}

func measurePriorityPage(t *testing.T, session *trustedInputBrowserSession) {
	t.Helper()
	for _, scheme := range []string{"light", "dark"} {
		t.Run(scheme, func(t *testing.T) {
			session.callDevToolsMethod(t, "Emulation.setEmulatedMedia", map[string]any{
				"features": []map[string]string{{"name": "prefers-color-scheme", "value": scheme}},
			}, true)
			var measured priorityBoardMeasurement
			session.decodeResult(t, "priority board measurement", session.evaluateInPage(t, `(function () {
  var groups = Array.from(document.querySelectorAll('[data-cards="pending"] .pending-group'));
  function ids(group) {
    return Array.from(group.querySelectorAll('.req-card')).map(function (card) { return card.dataset.detailId; });
  }
  var badges = Array.from(document.querySelectorAll('[data-cards="pending"] .badge-priority')).map(function (badge) {
    var card = badge.closest('.req-card');
    var title = card.querySelector('.req-card-title');
    var badgeRect = badge.getBoundingClientRect();
    var cardRect = card.getBoundingClientRect();
    var titleRect = title.getBoundingClientRect();
    var style = getComputedStyle(badge);
    return {requestID: card.dataset.detailId, text: badge.textContent.trim(), color: style.color,
      background: style.backgroundColor, width: badgeRect.width, height: badgeRect.height,
      contained: badgeRect.left >= cardRect.left && badgeRect.right <= cardRect.right,
      overlapsTitle: !(badgeRect.right <= titleRect.left || badgeRect.left >= titleRect.right ||
        badgeRect.bottom <= titleRect.top || badgeRect.top >= titleRect.bottom)};
  });
  return {href: location.href, browser: navigator.userAgent,
    scheme: matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light',
    ready: ids(groups[0]), waiting: ids(groups[1]), badges: badges};
})()`), &measured)
			if measured.Scheme != scheme {
				t.Fatalf("scheme = %q, want %q", measured.Scheme, scheme)
			}
			if !strings.HasPrefix(measured.Href, "file://") {
				t.Fatalf("static URL = %q", measured.Href)
			}
			if !reflect.DeepEqual(measured.Ready, []string{"REQ-403", "REQ-401", "REQ-402"}) || !reflect.DeepEqual(measured.Waiting, []string{"REQ-405", "REQ-404"}) {
				t.Fatalf("DOM priority order = ready %v waiting %v", measured.Ready, measured.Waiting)
			}
			wantBadges := map[string]string{"REQ-403": "now", "REQ-402": "later", "REQ-405": "now", "REQ-404": "later"}
			if len(measured.Badges) != len(wantBadges) {
				t.Fatalf("priority badges = %+v, want four and no next badge", measured.Badges)
			}
			for _, badge := range measured.Badges {
				if wantBadges[badge.RequestID] != badge.Text {
					t.Errorf("%s badge = %q, want %q", badge.RequestID, badge.Text, wantBadges[badge.RequestID])
				}
				delete(wantBadges, badge.RequestID)
				if badge.Width <= 0 || badge.Height <= 0 || badge.Color == "" || badge.Background == "" || !badge.Contained || badge.OverlapsTitle {
					t.Errorf("%s badge is not visibly contained: %+v", badge.RequestID, badge)
				}
			}
			if len(wantBadges) != 0 {
				t.Errorf("missing priority badges: %v", wantBadges)
			}
			t.Logf("static %s URL=%s browser=%s ready=%v waiting=%v badges=%+v", scheme, measured.Href, measured.Browser, measured.Ready, measured.Waiting, measured.Badges)
		})
	}
}
