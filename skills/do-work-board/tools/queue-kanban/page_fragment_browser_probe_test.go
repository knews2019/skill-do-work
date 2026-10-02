package main

import (
	"os"
	"path/filepath"
	"testing"
)

// REQ-626, in a real engine: a shared link such as probe.html#timeline opened the
// Board page, and clicking another page left the address bar on the bare URL. The
// page must open on Timeline, with its lazy first render done, and a click on
// Activity must write #activity so the address can be copied as that page's link.
func TestBrowserBehaviorPageFragmentOpensAndFollowsThePage(t *testing.T) {
	lookupBrowserForBehaviorProbe(t)
	siteDirectory := generateLiveSiteInDir(t)
	indexBytes, readError := os.ReadFile(filepath.Join(siteDirectory, "index.html"))
	if readError != nil {
		t.Fatal(readError)
	}
	session := startTrustedInputBrowserSessionAtFragment(
		t, "page fragment", siteDirectory, string(indexBytes), "#timeline")
	defer session.closeBrowserSession()

	session.waitForPageCondition(t, "the view buttons",
		`document.querySelector('[data-view-target="activity"]')`)
	var opened struct {
		Hash            string `json:"hash"`
		TimelineShown   bool   `json:"timelineShown"`
		BoardShown      bool   `json:"boardShown"`
		TimelinePressed string `json:"timelinePressed"`
		BoardPressed    string `json:"boardPressed"`
		TimelineSummary string `json:"timelineSummary"`
	}
	session.decodeResult(t, "page opened from #timeline", session.evaluateInPage(t, `({
  hash: location.hash,
  timelineShown: !document.getElementById("view-timeline").hidden,
  boardShown: !document.getElementById("view-board").hidden,
  timelinePressed: document.querySelector('[data-view-target="timeline"]').getAttribute("aria-pressed"),
  boardPressed: document.querySelector('[data-view-target="board"]').getAttribute("aria-pressed"),
  timelineSummary: document.getElementById("timeline-summary").textContent
})`), &opened)
	if opened.Hash != "#timeline" || !opened.TimelineShown || opened.BoardShown ||
		opened.TimelinePressed != "true" || opened.BoardPressed != "false" {
		t.Fatalf("probe.html#timeline did not open the Timeline page: %+v", opened)
	}
	if opened.TimelineSummary == "" {
		t.Fatalf("the Timeline page is visible but its lazy first render never ran: %+v", opened)
	}

	var afterClick struct {
		Hash          string `json:"hash"`
		ActivityShown bool   `json:"activityShown"`
		HistoryLength int    `json:"historyLength"`
	}
	session.decodeResult(t, "after the Activity click", session.evaluateInPage(t, `(function () {
  var historyLengthBefore = history.length;
  document.querySelector('[data-view-target="activity"]').click();
  return {
    hash: location.hash,
    activityShown: !document.getElementById("view-activity").hidden,
    historyLength: history.length - historyLengthBefore
  };
})()`), &afterClick)
	if afterClick.Hash != "#activity" || !afterClick.ActivityShown {
		t.Fatalf("clicking Activity did not write #activity: %+v", afterClick)
	}
	if afterClick.HistoryLength != 0 {
		t.Fatalf("the click added %d history entries; Back must still leave the board", afterClick.HistoryLength)
	}
}
