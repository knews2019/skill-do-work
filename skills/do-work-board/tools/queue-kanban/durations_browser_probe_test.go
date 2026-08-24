package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"
)

// sliceGeneratedStyleBlock lifts the page's inlined stylesheet for focused
// browser probes that need the shipped board styles without the full page.
func sliceGeneratedStyleBlock(t *testing.T, indexHtml string) string {
	t.Helper()
	styleStart := strings.Index(indexHtml, "<style>")
	styleEnd := strings.Index(indexHtml, "</style>")
	if styleStart < 0 || styleEnd < styleStart {
		t.Fatal("generated page carries no inlined <style> block")
	}
	return indexHtml[styleStart+len("<style>") : styleEnd]
}

type durationsDensePanelProbeResult struct {
	LocationHref          string  `json:"locationHref"`
	RenderedSampleCount   int     `json:"renderedSampleCount"`
	PayloadSampleCount    int     `json:"payloadSampleCount"`
	ActiveDayCount        int     `json:"activeDayCount"`
	DaySlotWidth          float64 `json:"daySlotWidth"`
	BusyDaySpread         float64 `json:"busyDaySpread"`
	Deterministic         bool    `json:"deterministic"`
	EveryMarkInOwnDay     bool    `json:"everyMarkInOwnDay"`
	HoveredRequestId      string  `json:"hoveredRequestId"`
	HoverReadout          string  `json:"hoverReadout"`
	RibbonFiniteBounded   bool    `json:"ribbonFiniteBounded"`
	MedianFiniteBounded   bool    `json:"medianFiniteBounded"`
	RibbonOpacity         float64 `json:"ribbonOpacity"`
	MarkOpacity           float64 `json:"markOpacity"`
	BodyBackground        string  `json:"bodyBackground"`
	LongestSpanCount      int     `json:"longestSpanCount"`
	ExpectedSpanCount     int     `json:"expectedSpanCount"`
	LongestSpanOrder      bool    `json:"longestSpanOrder"`
	EveryListField        bool    `json:"everyListField"`
	CountSentence         string  `json:"countSentence"`
	ListOutsideSVG        bool    `json:"listOutsideSvg"`
	SharedWrapper         bool    `json:"sharedWrapper"`
	SVGRequestLabels      int     `json:"svgRequestLabels"`
	SVGLeaderLines        int     `json:"svgLeaderLines"`
	SVGMoreSentences      int     `json:"svgMoreSentences"`
	OverflowHoverId       string  `json:"overflowHoverId"`
	OverflowHoverDuration string  `json:"overflowHoverDuration"`
	OverflowHoverReadout  string  `json:"overflowHoverReadout"`
	ViewportWidth         float64 `json:"viewportWidth"`
	DocumentWidth         float64 `json:"documentWidth"`
	WrapperRight          float64 `json:"wrapperRight"`
	AsideRight            float64 `json:"asideRight"`
	WrapperColumns        string  `json:"wrapperColumns"`
	AsideOverflowY        string  `json:"asideOverflowY"`
}

// REQ-349's target board is materially denser than this repository. This probe
// renders 705 real samples on 47 active days spread across 139 UTC days, making
// each day slot about eight SVG units wide. It measures the complete generated
// page in both colour schemes, including the hover seam that consumes markIndex.
func TestBrowserBehaviorDurationsDensePanelASpreadStaysBoundedAndInteractive(t *testing.T) {
	const samplesPerDay = 15
	const activeDayCount = 47
	fixtureStart := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	fixtureTickets := make([]*RequestTicket, 0, samplesPerDay*activeDayCount)
	for activeDayIndex := 0; activeDayIndex < activeDayCount; activeDayIndex++ {
		dayStart := fixtureStart.AddDate(0, 0, activeDayIndex*3)
		for sampleIndex := 0; sampleIndex < samplesPerDay; sampleIndex++ {
			completedAt := dayStart.Add(time.Duration(8*60+sampleIndex*37) * time.Minute)
			minutes := time.Duration(2+(sampleIndex*7)%57) * time.Minute
			if sampleIndex < 2 {
				// Ninety-four positive overflow samples force the complete list well
				// past one viewport. Equal values also exercise its REQ-id tie-break.
				minutes = time.Duration(90+(activeDayIndex%5)*15) * time.Minute
			}
			ticket := durationTicket(
				fmt.Sprintf("REQ-%04d", activeDayIndex*samplesPerDay+sampleIndex+1),
				[]string{"A", "B", "C"}[sampleIndex%3],
				completedAt.Add(-minutes).Format(time.RFC3339),
				completedAt.Format(time.RFC3339),
			)
			ticket.UserRequestId = fmt.Sprintf("UR-%03d", activeDayIndex+1)
			ticket.Domain = []string{"frontend", "backend", "testing"}[sampleIndex%3]
			ticket.Title = fmt.Sprintf("Dense sample %d with a wrapping title", activeDayIndex*samplesPerDay+sampleIndex+1)
			fixtureTickets = append(fixtureTickets, ticket)
		}
	}

	fixtureBoard := &Board{
		GeneratedAt: time.Date(2026, 8, 24, 12, 0, 0, 0, time.UTC),
		ProjectName: "REQ-349 dense durations probe",
		AllRequests: fixtureTickets,
	}
	siteDirectory := t.TempDir()
	if generateError := generateStaticSite(siteDirectory, fixtureBoard); generateError != nil {
		t.Fatalf("generate dense fixture board: %v", generateError)
	}
	indexBytes, readError := os.ReadFile(filepath.Join(siteDirectory, "index.html"))
	if readError != nil {
		t.Fatalf("read dense fixture index: %v", readError)
	}

	probeScript := `
  (function () {
  var resultNode = document.createElement("pre");
  resultNode.id = "` + browserProbeResultElementId + `";
  document.body.appendChild(resultNode);
  var durationsPanel = document.getElementById("view-durations");
  durationsPanel.hidden = false;
  durationsPanel.style.display = "block";

  function captureMarks() {
    return Array.from(document.querySelectorAll("#durations-chart circle.durations-mark")).map(function (circle) {
      return { x: Number(circle.getAttribute("cx")), y: Number(circle.getAttribute("cy")) };
    });
  }

  setDurationsWindow("all");
  renderDurationsView();
  var firstMarks = captureMarks();
  setDurationsWindow("all");
  renderDurationsView();
  var secondMarks = captureMarks();
  var samples = boardData.durations.samples;
  var firstSampleMs = Date.parse(samples[0].completionTime);
  var lastSampleMs = Date.parse(samples[samples.length - 1].completionTime);
  var timeStart = Math.floor(firstSampleMs / DURATIONS_DAY_MS) * DURATIONS_DAY_MS;
  var timeEnd = Math.floor(lastSampleMs / DURATIONS_DAY_MS) * DURATIONS_DAY_MS + DURATIONS_DAY_MS;
  var timeSpan = timeEnd - timeStart;
  var xOfEpoch = function (epochMs) {
    return DURATIONS_MARGIN_LEFT + ((epochMs - timeStart) / timeSpan) * DURATIONS_PLOT_WIDTH;
  };
  var everyMarkInOwnDay = secondMarks.every(function (mark, sampleIndex) {
    var sampleMs = Date.parse(samples[sampleIndex].completionTime);
    var dayStartMs = Math.floor(sampleMs / DURATIONS_DAY_MS) * DURATIONS_DAY_MS;
    return mark.x >= xOfEpoch(dayStartMs) - 0.05 && mark.x <= xOfEpoch(dayStartMs + DURATIONS_DAY_MS) + 0.05;
  });
  var firstDayMarks = secondMarks.slice(0, ` + fmt.Sprintf("%d", samplesPerDay) + `);
  var busyDaySpread = Math.max.apply(null, firstDayMarks.map(function (mark) { return mark.x; })) -
    Math.min.apply(null, firstDayMarks.map(function (mark) { return mark.x; }));
  var deterministic = firstMarks.length === secondMarks.length && firstMarks.every(function (mark, markIndex) {
    return mark.x === secondMarks[markIndex].x && mark.y === secondMarks[markIndex].y;
  });

  var ribbon = document.querySelector("#durations-chart .durations-quantile-ribbon");
  var median = document.querySelector("#durations-chart .durations-quantile-median");
  function finiteBounded(pathNode) {
    if (!pathNode || !pathNode.getAttribute("d") || /NaN|Infinity/.test(pathNode.getAttribute("d"))) { return false; }
    var box = pathNode.getBBox();
    return [box.x, box.y, box.width, box.height, pathNode.getTotalLength()].every(Number.isFinite) &&
      box.x >= DURATIONS_MARGIN_LEFT - 0.1 && box.x + box.width <= DURATIONS_VIEW_WIDTH - DURATIONS_MARGIN_RIGHT + 0.1 &&
      box.y >= DURATIONS_MAIN_TOP - 0.1 && box.y + box.height <= DURATIONS_MAIN_BOTTOM + 0.1;
  }

  var hoverMarkIndex = 7;
  var hoveredMark = secondMarks[hoverMarkIndex];
  var svg = document.querySelector("#durations-chart svg");
  var hoverSurface = document.querySelector("#durations-chart .durations-hover-surface");
  var bounds = svg.getBoundingClientRect();
  hoverSurface.dispatchEvent(new MouseEvent("mousemove", {
    bubbles: true,
    clientX: bounds.left + hoveredMark.x * bounds.width / DURATIONS_VIEW_WIDTH,
    clientY: bounds.top + hoveredMark.y * bounds.height / DURATIONS_VIEW_HEIGHT
  }));
  var hoveredRequestId = samples[hoverMarkIndex].id;
  var hoverReadout = document.getElementById("durations-readout").textContent;

  var expectedLongestSpans = samples.filter(function (sample) {
    return sample.wallMinutes > DURATIONS_CEILING_MINUTES;
  }).sort(function (first, second) {
    if (first.wallMinutes !== second.wallMinutes) { return second.wallMinutes - first.wallMinutes; }
    if (first.id !== second.id) { return first.id < second.id ? -1 : 1; }
    return first.completionTime < second.completionTime ? -1 : 1;
  });
  var longestSpanRows = Array.from(document.querySelectorAll("#durations-longest-list > li"));
  var longestSpanOrder = longestSpanRows.every(function (row, rowIndex) {
    return expectedLongestSpans[rowIndex] && row.getAttribute("data-request-id") === expectedLongestSpans[rowIndex].id &&
      Number(row.getAttribute("data-wall-minutes")) === expectedLongestSpans[rowIndex].wallMinutes;
  });
  var everyListField = longestSpanRows.every(function (row, rowIndex) {
    var sample = expectedLongestSpans[rowIndex];
    var request = boardData.requests[sample.id];
    return row.querySelector(".durations-longest-spans-request").textContent === sample.id &&
      row.querySelector(".durations-longest-spans-user-request").textContent === request.userRequestId &&
      row.querySelector(".durations-longest-spans-duration").textContent === formatDurationMinutes(sample.wallMinutes) &&
      row.querySelector(".durations-longest-spans-route").textContent === durationRouteName(sample.route) &&
      row.querySelector(".durations-longest-spans-title").textContent === request.title;
  });
  var overflowSample = expectedLongestSpans[0];
  var overflowSampleIndex = samples.indexOf(overflowSample);
  var overflowMark = secondMarks[overflowSampleIndex];
  hoverSurface.dispatchEvent(new MouseEvent("mousemove", {
    bubbles: true,
    clientX: bounds.left + overflowMark.x * bounds.width / DURATIONS_VIEW_WIDTH,
    clientY: bounds.top + overflowMark.y * bounds.height / DURATIONS_VIEW_HEIGHT
  }));
  var overflowHoverReadout = document.getElementById("durations-readout").textContent;
  var chart = document.getElementById("durations-chart");
  var list = document.getElementById("durations-longest-list");
  var aside = list.closest(".durations-longest-spans");
  var wrapper = chart.parentElement;
  var svgRequestLabels = Array.from(svg.querySelectorAll("text")).filter(function (textNode) {
    return /^REQ-[0-9]+(?:\s|$)/.test(textNode.textContent);
  });
  var svgMoreSentences = Array.from(svg.querySelectorAll("text")).filter(function (textNode) {
    return /^\+[0-9]+ more\b/.test(textNode.textContent);
  });
  var wrapperBounds = wrapper.getBoundingClientRect();
  var asideBounds = aside.getBoundingClientRect();

  document.getElementById("` + browserProbeResultElementId + `").textContent = JSON.stringify({
    locationHref: location.href,
    renderedSampleCount: secondMarks.length,
    payloadSampleCount: samples.length,
    activeDayCount: boardData.durations.days.length,
    daySlotWidth: DURATIONS_PLOT_WIDTH * DURATIONS_DAY_MS / timeSpan,
    busyDaySpread: busyDaySpread,
    deterministic: deterministic,
    everyMarkInOwnDay: everyMarkInOwnDay,
    hoveredRequestId: hoveredRequestId,
    hoverReadout: hoverReadout,
    ribbonFiniteBounded: finiteBounded(ribbon),
    medianFiniteBounded: finiteBounded(median),
    ribbonOpacity: ribbon ? Number(getComputedStyle(ribbon).opacity) : 0,
    markOpacity: Number(getComputedStyle(document.querySelector("#durations-chart circle.durations-mark:not(.durations-mark-critical):not(.durations-mark-unknown)")).opacity),
    bodyBackground: getComputedStyle(document.body).backgroundColor,
    longestSpanCount: longestSpanRows.length,
    expectedSpanCount: expectedLongestSpans.length,
    longestSpanOrder: longestSpanOrder,
    everyListField: everyListField,
    countSentence: document.getElementById("durations-longest-count").textContent,
    listOutsideSvg: !svg.contains(list),
    sharedWrapper: wrapper === aside.parentElement,
    svgRequestLabels: svgRequestLabels.length,
    svgLeaderLines: svg.querySelectorAll(".durations-label-leader").length,
    svgMoreSentences: svgMoreSentences.length,
    overflowHoverId: overflowSample.id,
    overflowHoverDuration: formatDurationMinutes(overflowSample.wallMinutes),
    overflowHoverReadout: overflowHoverReadout,
    viewportWidth: window.innerWidth,
    documentWidth: document.documentElement.scrollWidth,
    wrapperRight: wrapperBounds.right,
    asideRight: asideBounds.right,
    wrapperColumns: getComputedStyle(wrapper).gridTemplateColumns,
    asideOverflowY: getComputedStyle(aside).overflowY
  });
  })();
`
	// The generated client is one IIFE, so the probe has to run before its final
	// close to exercise the private renderer and payload rather than a copied
	// helper. The result node itself remains ordinary HTML outside the script.
	probePage := string(indexBytes)
	clientClose := strings.LastIndex(probePage, "})();")
	if clientClose < 0 {
		t.Fatal("generated page has no client IIFE close for dense durations probe")
	}
	probePage = probePage[:clientClose] + probeScript + probePage[clientClose:]

	for _, viewport := range []struct {
		name          string
		colourFlag    string
		width         int
		stackedLayout bool
	}{
		{name: "light-320", colourFlag: "--force-light-mode", width: 320, stackedLayout: true},
		{name: "dark-320", colourFlag: "--force-dark-mode", width: 320, stackedLayout: true},
		{name: "light-768", colourFlag: "--force-light-mode", width: 768, stackedLayout: true},
		{name: "dark-768", colourFlag: "--force-dark-mode", width: 768, stackedLayout: true},
		{name: "light-1280", colourFlag: "--force-light-mode", width: 1280},
		{name: "dark-1280", colourFlag: "--force-dark-mode", width: 1280},
	} {
		viewport := viewport
		t.Run(viewport.name, func(t *testing.T) {
			resultJSON := runBrowserBehaviorProbeInDirectory(
				t, "dense durations "+viewport.name, siteDirectory, probePage,
				fmt.Sprintf("--window-size=%d,1200", viewport.width), viewport.colourFlag,
			)
			var result durationsDensePanelProbeResult
			if decodeError := json.Unmarshal(resultJSON, &result); decodeError != nil {
				t.Fatalf("decode dense durations probe: %v\n%s", decodeError, resultJSON)
			}
			if result.PayloadSampleCount != len(fixtureTickets) {
				t.Errorf("payload carries %d samples, want the %d-sample fixture",
					result.PayloadSampleCount, len(fixtureTickets))
			}
			if result.RenderedSampleCount != result.PayloadSampleCount {
				t.Errorf("all-history render drew %d .durations-mark circles for %d payload samples",
					result.RenderedSampleCount, result.PayloadSampleCount)
			}
			if result.RenderedSampleCount < 700 || result.ActiveDayCount != activeDayCount {
				t.Errorf("rendered %d samples across %d active days, want at least 700 across %d",
					result.RenderedSampleCount, result.ActiveDayCount, activeDayCount)
			}
			if result.DaySlotWidth < 7.5 || result.DaySlotWidth > 8.5 {
				t.Errorf("day slot width %.2f, want roughly 8 SVG units", result.DaySlotWidth)
			}
			if !result.EveryMarkInOwnDay || result.BusyDaySpread < 5 {
				t.Errorf("own-day=%v, busy-day spread=%.2f; want bounded useful spread", result.EveryMarkInOwnDay, result.BusyDaySpread)
			}
			if !result.Deterministic {
				t.Error("identical payload moved marks across consecutive renders")
			}
			if !strings.HasPrefix(result.HoverReadout, result.HoveredRequestId+" ·") {
				t.Errorf("hover at %s's jittered centre read %q", result.HoveredRequestId, result.HoverReadout)
			}
			if !result.RibbonFiniteBounded || !result.MedianFiniteBounded {
				t.Errorf("distribution geometry bounded: ribbon=%v median=%v", result.RibbonFiniteBounded, result.MedianFiniteBounded)
			}
			if result.RibbonOpacity <= 0 || result.RibbonOpacity >= result.MarkOpacity {
				t.Errorf("ribbon opacity %.2f is not subordinate to mark opacity %.2f", result.RibbonOpacity, result.MarkOpacity)
			}
			if result.LocationHref == "" || !strings.Contains(result.LocationHref, browserProbePageFileName) {
				t.Errorf("probe measured unnamed page %q", result.LocationHref)
			}
			if result.BodyBackground == "" || result.BodyBackground == "rgba(0, 0, 0, 0)" {
				t.Errorf("%s board has no resolved body background: %q", viewport.name, result.BodyBackground)
			}
			if result.ExpectedSpanCount < 60 || result.LongestSpanCount != result.ExpectedSpanCount {
				t.Errorf("longest-spans list rendered %d of %d overflow samples; fixture must carry at least 60",
					result.LongestSpanCount, result.ExpectedSpanCount)
			}
			if !result.LongestSpanOrder || !result.EveryListField {
				t.Errorf("complete list order=%v fields=%v; want descending/tied order with all five fields",
					result.LongestSpanOrder, result.EveryListField)
			}
			wantCountSentence := fmt.Sprintf("%d spans over 60 minutes in this window; all are listed.", result.ExpectedSpanCount)
			if result.CountSentence != wantCountSentence {
				t.Errorf("count sentence = %q, want %q", result.CountSentence, wantCountSentence)
			}
			if !result.ListOutsideSVG || !result.SharedWrapper {
				t.Errorf("list outside SVG=%v shared wrapper=%v, want adjacent HTML", result.ListOutsideSVG, result.SharedWrapper)
			}
			if result.SVGRequestLabels != 0 || result.SVGLeaderLines != 0 || result.SVGMoreSentences != 0 {
				t.Errorf("SVG still carries %d REQ labels, %d leaders, and %d +N-more sentences",
					result.SVGRequestLabels, result.SVGLeaderLines, result.SVGMoreSentences)
			}
			if !strings.HasPrefix(result.OverflowHoverReadout, result.OverflowHoverId+" ·") ||
				!strings.Contains(result.OverflowHoverReadout, " · "+result.OverflowHoverDuration+" ·") {
				t.Errorf("overflow hover for %s read %q", result.OverflowHoverId, result.OverflowHoverReadout)
			}
			if result.AsideOverflowY != "auto" {
				t.Errorf("longest-spans aside overflow-y = %q, want auto", result.AsideOverflowY)
			}
			if result.WrapperRight > result.ViewportWidth+1 || result.AsideRight > result.ViewportWidth+1 {
				t.Errorf("horizontal clipping at viewport %.0f: document %.0f wrapper-right %.0f aside-right %.0f",
					result.ViewportWidth, result.DocumentWidth, result.WrapperRight, result.AsideRight)
			}
			columnCount := len(strings.Fields(result.WrapperColumns))
			if viewport.stackedLayout && columnCount != 1 {
				t.Errorf("%dpx layout has grid columns %q, want one stacked column", viewport.width, result.WrapperColumns)
			}
			if !viewport.stackedLayout && columnCount != 2 {
				t.Errorf("%dpx layout has grid columns %q, want chart plus aside", viewport.width, result.WrapperColumns)
			}
			t.Logf("%s %s: viewport %.0f, %d/%d longest spans, columns %s, body %s",
				viewport.name, result.LocationHref, result.ViewportWidth, result.LongestSpanCount,
				result.ExpectedSpanCount, result.WrapperColumns, result.BodyBackground)
		})
	}
}

// REQ-266's mechanism, chosen as a check rather than a review convention.
//
// The rule REQ-252 established on the Go side: a browser-measured number must
// name the build it was measured on, because a face is per-browser and an undated
// number reads as timeless fact. `go/parser` cannot reach the JS surface, which is
// why the JS comments drifted into carrying three such numbers with no build
// between them.
//
// THE DISCRIMINATOR IS WHAT THE NUMBER CLAIMS, not whether it is a measurement.
// A measured number cited as EVIDENCE FOR A PAST DECISION is already dated — by
// the REQ it names, which is a stronger anchor than a build string, since the REQ
// carries the whole argument. A measured number presented as CURRENT FACT about
// the face in use is the one that needs a build, and after REQ-292 the honest
// answer for those is not to date them but to delete them: the engine answers at
// test time instead.
//
// So this check enforces the rule in the direction that survives: a measured
// number in a JS comment must sit beside either a REQ reference (evidence) or a
// build name (dated fact). A bare one is neither, and is what this catches.
func TestDurationsJavaScriptCommentsDateTheirMeasurements(t *testing.T) {
	rendererBytes, readError := embeddedWebAssets.ReadFile("web/board-durations.js")
	if readError != nil {
		t.Fatalf("read web/board-durations.js: %v", readError)
	}

	// A measurement claim: a decimal with two or more places, next to a word that
	// makes it a statement about drawn text. Plain geometry constants (13, 11, 9)
	// are declared numbers, not measurements, and are deliberately not matched.
	measurementClaim := regexp.MustCompile(
		`(?i)[0-9]+\.[0-9]{2,}[- ]?(?:unit|px|em)?s?\b[^\n]*?\b(?:measur|ascent|descent|line box|per character|draws)`)
	alternate := regexp.MustCompile(
		`(?i)\b(?:measur|ascent|descent|line box|per character|draws)\b[^\n]*?[0-9]+\.[0-9]{2,}`)
	// Either anchor dates the claim: the REQ that measured it, or the build.
	datingAnchor := regexp.MustCompile(`(?i)REQ-[0-9]{3}|chromium|playwright|firefox|webkit|safari`)

	for lineNumber, line := range strings.Split(string(rendererBytes), "\n") {
		commentStart := strings.Index(line, "//")
		if commentStart < 0 {
			continue
		}
		commentText := line[commentStart:]
		if !measurementClaim.MatchString(commentText) && !alternate.MatchString(commentText) {
			continue
		}
		// The dating anchor may sit anywhere in the surrounding comment block, not
		// only on the claim's own line — a paragraph names its REQ once.
		blockText := durationsSurroundingCommentBlock(string(rendererBytes), lineNumber)
		if datingAnchor.MatchString(blockText) {
			continue
		}
		t.Errorf("web/board-durations.js:%d carries a measured number with no REQ or build to date it:\n  %s\n"+
			"A face is per-browser, so an undated measurement reads as timeless fact. Cite the REQ that "+
			"measured it, name the build, or delete the number and let a browser probe answer at test time.",
			lineNumber+1, strings.TrimSpace(commentText))
	}
}

// durationsSurroundingCommentBlock returns the contiguous run of comment lines the
// given line sits in, so a paragraph that names its REQ once dates every claim in it.
func durationsSurroundingCommentBlock(sourceText string, lineNumber int) string {
	lines := strings.Split(sourceText, "\n")
	isComment := func(index int) bool {
		return index >= 0 && index < len(lines) && strings.HasPrefix(strings.TrimSpace(lines[index]), "//")
	}
	blockStart := lineNumber
	for isComment(blockStart - 1) {
		blockStart--
	}
	blockEnd := lineNumber
	for isComment(blockEnd + 1) {
		blockEnd++
	}
	return strings.Join(lines[blockStart:blockEnd+1], "\n")
}
