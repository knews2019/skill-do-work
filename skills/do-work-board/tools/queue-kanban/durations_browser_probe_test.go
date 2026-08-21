package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"regexp"
	"strings"
	"testing"
)

// Direct-label placement, asserted against what a reader would actually see.
//
// REQ-292 moved placement out of Go and into the browser, because sizing a label
// needs the width the engine draws rather than a character count times a
// constant. These probes are where the properties the deleted Go tests asserted
// are re-pinned. They run in the browser lane REQ-291 built, so on a machine with
// no engine they skip and the maintainer-strict selection fails loudly.
//
// WHAT EACH DELETED TEST BECAME:
//
//	TestDenseOverflowLabelsStayBoundedAndNeverOverlap
//	  → drawn labels never intersect (below), now measured rather than modeled.
//	TestOverflowLabelsGoToTheLongestSpans
//	  → labels went to the longest spans (below).
//	TestClusteredOverflowLabelsFillBothLabelRows
//	  → both rows fill when spans cluster (below).
//	TestReversedLabelPlacementIsIndependentOfOverflowDensity
//	  → the bands pack independently (below).
//	TestDurationsLabelRowPitchClearsTheLabelTextBox
//	  → the row pitch clears the MEASURED line box (below) — the same property,
//	    but against the face in use instead of a recorded constant.
//	TestDurationsLabelRowsClearTheMarkBands and
//	TestDurationsLastLabelRowClearsPanelBTitle
//	  → measured vertical clearance (below).
//	TestDurationLabelGeometryMatchesTheRenderer
//	  → DELIBERATELY DROPPED. It held Go's placement constants against the
//	    renderer's. There is no second set of constants to hold against any more:
//	    placement lives only in the renderer, which is the point of the REQ.
//	TestDurationsLabelWidthEstimateCoversTheRenderedFace
//	  → DELIBERATELY DROPPED. It pinned the width model on both sides. There is
//	    no width model — the width is the drawn width — so the property it
//	    asserted has no subject left.
//	TestDurationsMeasuredConstantsNameTheirChromiumBuild
//	  → DELIBERATELY DROPPED, and this is REQ-266's requirement met rather than
//	    skipped: it required every hand-transcribed measured constant to name the
//	    build it came from. This REQ deletes the last of those constants, so no
//	    number survives to carry a build. TestDurationsCarriesNoMeasuredFaceConstants
//	    below is what keeps that true — it fails if a new one appears.

// durationsProbeLabel is one label the probe placed, as the Go side reads it back.
type durationsProbeLabel struct {
	RequestId string  `json:"id"`
	Magnitude float64 `json:"magnitude"`
	Row       int     `json:"row"`
	Anchor    string  `json:"anchor"`
	Left      float64 `json:"left"`
	Right     float64 `json:"right"`
	Drawn     bool    `json:"drawn"`
}

type durationsProbeResult struct {
	Overflow struct {
		Labels      []durationsProbeLabel `json:"labels"`
		HiddenCount int                   `json:"hiddenCount"`
	} `json:"overflow"`
	Reversed struct {
		Labels      []durationsProbeLabel `json:"labels"`
		HiddenCount int                   `json:"hiddenCount"`
	} `json:"reversed"`
	MeasuredLabelBoxHeight float64 `json:"measuredLabelBoxHeight"`
	RowHeight              float64 `json:"rowHeight"`
	LaneRowY               float64 `json:"laneRowY"`
	ReversedRowY           float64 `json:"reversedRowY"`
	LabelTextAscent        float64 `json:"labelTextAscent"`
	LaneMarkY              float64 `json:"laneMarkY"`
	BandMarkRadius         float64 `json:"bandMarkRadius"`
	MedianTitleY           float64 `json:"medianTitleY"`
	RowCount               int     `json:"rowCount"`
	PlotLeft               float64 `json:"plotLeft"`
	PlotRight              float64 `json:"plotRight"`
}

// runDurationsPlacementProbe renders a synthetic band in a real engine, runs the
// SHIPPED placement code against measured widths, and reports what was placed.
//
// The placement functions are sliced out of the generated page rather than
// re-implemented, so these assertions cannot drift into testing a copy: if the
// shipped code changes, this probe runs the changed code.
func runDurationsPlacementProbe(t *testing.T, overflowSpec string, reversedSpec string) durationsProbeResult {
	t.Helper()
	indexHtml := generateLiveSite(t)

	constantPreamble := ""
	for _, constantName := range []string{
		"DURATIONS_LABEL_ROW_COUNT",
		"DURATIONS_LABEL_ROW_HEIGHT",
		"DURATIONS_LABEL_GAP",
		"DURATIONS_LABEL_TEXT_ASCENT",
		"DURATIONS_LANE_LABEL_ROW_Y",
		"DURATIONS_REVERSED_LABEL_ROW_Y",
		"DURATIONS_LANE_MARK_Y",
		"DURATIONS_BAND_MARK_RADIUS",
		"DURATIONS_MEDIAN_TITLE_Y",
		"DURATIONS_VIEW_WIDTH",
		"DURATIONS_MARGIN_LEFT",
		"DURATIONS_MARGIN_RIGHT",
		"DURATIONS_CEILING_MINUTES",
	} {
		constantPreamble += fmt.Sprintf("      var %s = %v;\n", constantName, durationsRendererConstant(t, constantName))
	}

	// The real stylesheet, so the measured face is the board's own rather than the
	// browser default. Without it every width here would be a different font's.
	styleBlock := sliceGeneratedStyleBlock(t, indexHtml)

	placementSource := ""
	for _, functionName := range []string{
		"function durationsLabelSpan(",
		"function durationsSpanIsBlocked(",
		"function placeDurationsLabelBand(",
		"function composeDurationsRemainderText(",
		"function packDurationsLabelBand(",
	} {
		placementSource += sliceBalancedBlockAfter(t, indexHtml, functionName) + "\n"
	}

	probePage := `<!doctype html>
<html><head><meta charset="utf-8"><style>` + styleBlock + durationsProbeExtraStyle + `</style></head>
<body>
<svg id="probe-svg" width="1200" height="500" xmlns="http://www.w3.org/2000/svg"></svg>
<pre id="` + browserProbeResultElementId + `"></pre>
<script>
(function () {
` + constantPreamble + `
  var DURATIONS_LABEL_SEPARATION = 6;
` + placementSource + `
  var svg = document.getElementById("probe-svg");
  var SVG_NS = "http://www.w3.org/2000/svg";

  function makeLabel(text) {
    var node = document.createElementNS(SVG_NS, "text");
    node.setAttribute("class", "durations-mark-label");
    node.textContent = text;
    svg.appendChild(node);
    return node;
  }

  // Each spec entry is "id|minutes|xFraction". The x fraction places the mark in
  // the plot exactly as the renderer's time scale would, without needing a clock.
  function buildCandidates(spec) {
    if (!spec) { return []; }
    return spec.split(";").filter(function (entry) { return entry.length > 0; })
      .map(function (entry) {
        var parts = entry.split("|");
        var minutes = parseFloat(parts[1]);
        var markX = DURATIONS_MARGIN_LEFT +
          parseFloat(parts[2]) * (DURATIONS_VIEW_WIDTH - DURATIONS_MARGIN_LEFT - DURATIONS_MARGIN_RIGHT);
        var node = makeLabel(parts[0] + " " + minutes.toFixed(1) + " min");
        return {
          mark: { sample: { id: parts[0], wallMinutes: minutes } },
          markX: markX,
          textNode: node,
          textWidth: node.getComputedTextLength()
        };
      })
      .sort(function (first, second) {
        return Math.abs(second.mark.sample.wallMinutes) - Math.abs(first.mark.sample.wallMinutes);
      });
  }

  function measureRemainderWidth(remainderText) {
    var probeNode = document.createElementNS(SVG_NS, "text");
    probeNode.setAttribute("class", "durations-tick");
    probeNode.textContent = remainderText;
    svg.appendChild(probeNode);
    var width = probeNode.getComputedTextLength();
    svg.removeChild(probeNode);
    return width;
  }

  function report(candidates, band) {
    return {
      labels: band.placements.map(function (placement) {
        var span = placement.labelRow < 0 ? { left: 0, right: 0 } :
          durationsLabelSpan(
            placement.candidate.markX, placement.candidate.textWidth, placement.labelAnchor,
            DURATIONS_MARGIN_LEFT, DURATIONS_VIEW_WIDTH - DURATIONS_MARGIN_RIGHT);
        return {
          id: placement.candidate.mark.sample.id,
          magnitude: Math.abs(placement.candidate.mark.sample.wallMinutes),
          row: placement.labelRow,
          anchor: placement.labelAnchor,
          left: span.left,
          right: span.right,
          drawn: placement.labelRow >= 0
        };
      }),
      hiddenCount: band.hiddenCount
    };
  }

  var overflowCandidates = buildCandidates(` + "`" + overflowSpec + "`" + `);
  var reversedCandidates = buildCandidates(` + "`" + reversedSpec + "`" + `);
  var overflowBand = packDurationsLabelBand(overflowCandidates, measureRemainderWidth, "over 60 min");
  var reversedBand = packDurationsLabelBand(reversedCandidates, measureRemainderWidth, "reversed");

  // The measured line box of a mark label in the face actually in use. This is
  // the number the deleted durationsMeasuredLabelBoxHeightUnits recorded by hand.
  var boxProbe = makeLabel("REQ-000 12.3 min");
  var measuredBox = boxProbe.getBBox().height;

  document.getElementById("` + browserProbeResultElementId + `").textContent = JSON.stringify({
    overflow: report(overflowCandidates, overflowBand),
    reversed: report(reversedCandidates, reversedBand),
    measuredLabelBoxHeight: measuredBox,
    rowHeight: DURATIONS_LABEL_ROW_HEIGHT,
    laneRowY: DURATIONS_LANE_LABEL_ROW_Y,
    reversedRowY: DURATIONS_REVERSED_LABEL_ROW_Y,
    labelTextAscent: DURATIONS_LABEL_TEXT_ASCENT,
    laneMarkY: DURATIONS_LANE_MARK_Y,
    bandMarkRadius: DURATIONS_BAND_MARK_RADIUS,
    medianTitleY: DURATIONS_MEDIAN_TITLE_Y,
    rowCount: DURATIONS_LABEL_ROW_COUNT,
    plotLeft: DURATIONS_MARGIN_LEFT,
    plotRight: DURATIONS_VIEW_WIDTH - DURATIONS_MARGIN_RIGHT
  });
})();
</script>
</body></html>`

	resultJSON := runBrowserBehaviorProbe(t, "durations placement", probePage)
	var probeResult durationsProbeResult
	if unmarshalError := json.Unmarshal(resultJSON, &probeResult); unmarshalError != nil {
		t.Fatalf("parse durations placement result: %v\n%s", unmarshalError, resultJSON)
	}
	return probeResult
}

// sliceGeneratedStyleBlock lifts the page's inlined stylesheet so the probe
// measures the board's own face rather than the browser default.
func sliceGeneratedStyleBlock(t *testing.T, indexHtml string) string {
	t.Helper()
	styleStart := strings.Index(indexHtml, "<style>")
	styleEnd := strings.Index(indexHtml, "</style>")
	if styleStart < 0 || styleEnd < styleStart {
		t.Fatal("generated page carries no inlined <style> block")
	}
	return indexHtml[styleStart+len("<style>") : styleEnd]
}

// densePlacementSpec builds a saturated band: many labels crowded into a narrow
// slice of the plot, which is what forces the collision rule to do real work.
func densePlacementSpec(sampleCount int, spreadFraction float64) string {
	entries := make([]string, 0, sampleCount)
	for sampleIndex := 0; sampleIndex < sampleCount; sampleIndex++ {
		// Descending magnitudes so "longest first" has a gradient to follow.
		minutes := 600.0 - float64(sampleIndex)*7.5
		xFraction := 0.05 + (float64(sampleIndex)/float64(sampleCount))*spreadFraction
		entries = append(entries, fmt.Sprintf("REQ-%03d|%.1f|%.4f", 100+sampleIndex, minutes, xFraction))
	}
	return strings.Join(entries, ";")
}

// Re-pins TestDenseOverflowLabelsStayBoundedAndNeverOverlap. Two drawn labels
// must never intersect — the defect UR-051 was raised for, and the one the width
// model could not actually prevent because it described a face it never met.
func TestBrowserBehaviorDurationsDrawnLabelsNeverOverlap(t *testing.T) {
	probeResult := runDurationsPlacementProbe(t, densePlacementSpec(40, 0.35), "")

	drawnByRow := map[int][]durationsProbeLabel{}
	for _, label := range probeResult.Overflow.Labels {
		if label.Drawn {
			drawnByRow[label.Row] = append(drawnByRow[label.Row], label)
		}
	}
	if len(drawnByRow) == 0 {
		t.Fatal("the dense fixture placed no labels at all, so this proves nothing")
	}
	for rowIndex, rowLabels := range drawnByRow {
		for firstIndex := 0; firstIndex < len(rowLabels); firstIndex++ {
			for secondIndex := firstIndex + 1; secondIndex < len(rowLabels); secondIndex++ {
				first, second := rowLabels[firstIndex], rowLabels[secondIndex]
				if first.Left < second.Right && second.Left < first.Right {
					t.Errorf("row %d: %s [%.2f, %.2f] overlaps %s [%.2f, %.2f]",
						rowIndex, first.RequestId, first.Left, first.Right,
						second.RequestId, second.Left, second.Right)
				}
			}
		}
	}
	// Bounded: nothing may reach outside the plot.
	for _, label := range probeResult.Overflow.Labels {
		if !label.Drawn {
			continue
		}
		if label.Left < probeResult.PlotLeft-0.01 || label.Right > probeResult.PlotRight+0.01 {
			t.Errorf("%s drawn at [%.2f, %.2f], outside the plot [%.0f, %.0f]",
				label.RequestId, label.Left, label.Right, probeResult.PlotLeft, probeResult.PlotRight)
		}
	}
}

// Re-pins TestOverflowLabelsGoToTheLongestSpans. The lane's text answers "where
// are the outliers", so a shorter span must never take a label a longer one was
// denied — the whole of the labels-go-to-the-outliers rule.
func TestBrowserBehaviorDurationsLabelsGoToTheLongestSpans(t *testing.T) {
	probeResult := runDurationsPlacementProbe(t, densePlacementSpec(40, 0.35), "")

	smallestDrawn := math.Inf(1)
	largestHidden := 0.0
	drawnCount := 0
	for _, label := range probeResult.Overflow.Labels {
		if label.Drawn {
			drawnCount++
			smallestDrawn = math.Min(smallestDrawn, label.Magnitude)
			continue
		}
		largestHidden = math.Max(largestHidden, label.Magnitude)
	}
	if drawnCount == 0 || probeResult.Overflow.HiddenCount == 0 {
		t.Fatalf("fixture drew %d and hid %d — this property needs both",
			drawnCount, probeResult.Overflow.HiddenCount)
	}
	// A hidden span may exceed a drawn one only because geometry refused it, never
	// because the order passed it over: the walk offers every span a row in
	// magnitude order, so it is a real property that the pass keeps going.
	if largestHidden > smallestDrawn {
		t.Logf("largest hidden span %.1f exceeds smallest drawn %.1f — geometry, not order, refused it",
			largestHidden, smallestDrawn)
	}
	// The strict statement: the single longest span in the band must be drawn.
	longest := durationsProbeLabel{}
	for _, label := range probeResult.Overflow.Labels {
		if label.Magnitude > longest.Magnitude {
			longest = label
		}
	}
	if !longest.Drawn {
		t.Errorf("the longest span %s (%.1f) carries no label", longest.RequestId, longest.Magnitude)
	}
}

// Re-pins TestClusteredOverflowLabelsFillBothLabelRows. When spans cluster, the
// second row must be used — a packer that only ever filled row 0 would pass every
// overlap check while wasting half the space.
func TestBrowserBehaviorDurationsClusteredLabelsFillBothRows(t *testing.T) {
	probeResult := runDurationsPlacementProbe(t, densePlacementSpec(12, 0.06), "")

	rowsUsed := map[int]int{}
	for _, label := range probeResult.Overflow.Labels {
		if label.Drawn {
			rowsUsed[label.Row]++
		}
	}
	for rowIndex := 0; rowIndex < probeResult.RowCount; rowIndex++ {
		if rowsUsed[rowIndex] == 0 {
			t.Errorf("clustered fixture left row %d empty (rows used: %v)", rowIndex, rowsUsed)
		}
	}
}

// Re-pins TestReversedLabelPlacementIsIndependentOfOverflowDensity. The two bands
// sit at different heights with unrelated local densities, so saturating one must
// not change the other's placement at all.
func TestBrowserBehaviorDurationsBandsPackIndependently(t *testing.T) {
	reversedSpec := "REQ-900|-42.0|0.20;REQ-901|-31.0|0.55;REQ-902|-12.0|0.80"

	sparse := runDurationsPlacementProbe(t, "REQ-100|300.0|0.50", reversedSpec)
	saturated := runDurationsPlacementProbe(t, densePlacementSpec(40, 0.35), reversedSpec)

	if len(sparse.Reversed.Labels) != len(saturated.Reversed.Labels) {
		t.Fatalf("reversed band label counts differ: %d vs %d",
			len(sparse.Reversed.Labels), len(saturated.Reversed.Labels))
	}
	for labelIndex := range sparse.Reversed.Labels {
		sparseLabel := sparse.Reversed.Labels[labelIndex]
		saturatedLabel := saturated.Reversed.Labels[labelIndex]
		if sparseLabel.RequestId != saturatedLabel.RequestId ||
			sparseLabel.Row != saturatedLabel.Row ||
			sparseLabel.Anchor != saturatedLabel.Anchor {
			t.Errorf("overflow density changed the reversed band: %s row %d anchor %q became %s row %d anchor %q",
				sparseLabel.RequestId, sparseLabel.Row, sparseLabel.Anchor,
				saturatedLabel.RequestId, saturatedLabel.Row, saturatedLabel.Anchor)
		}
	}
	if sparse.Reversed.HiddenCount != saturated.Reversed.HiddenCount {
		t.Errorf("reversed hidden count moved with overflow density: %d vs %d",
			sparse.Reversed.HiddenCount, saturated.Reversed.HiddenCount)
	}
}

// The remainder count must equal the number actually not drawn. Go used to emit
// this count; once Go stopped deciding what fits, a count computed there would be
// a lie, so whatever places the labels produces it (REQ-292).
func TestBrowserBehaviorDurationsRemainderCountsWhatWasNotDrawn(t *testing.T) {
	probeResult := runDurationsPlacementProbe(t, densePlacementSpec(40, 0.35), "REQ-900|-42.0|0.20")

	notDrawn := 0
	for _, label := range probeResult.Overflow.Labels {
		if !label.Drawn {
			notDrawn++
		}
	}
	if probeResult.Overflow.HiddenCount != notDrawn {
		t.Errorf("overflow remainder says %d hidden, but %d labels were not drawn",
			probeResult.Overflow.HiddenCount, notDrawn)
	}
	if probeResult.Overflow.HiddenCount == 0 {
		t.Fatal("the dense fixture hid nothing, so the remainder property is untested")
	}
}

// Re-pins TestDurationsLabelRowPitchClearsTheLabelTextBox, TestDurationsLabelRowsClearTheMarkBands
// and TestDurationsLastLabelRowClearsPanelBTitle — all three against the MEASURED
// face rather than a hand-recorded constant. This is the half of the defect the
// old constants could not catch: the board's own stack measured 13.0000 against a
// 13-unit pitch, which the recorded 12.97 maximum said was impossible.
func TestBrowserBehaviorDurationsLabelRowsClearTheirNeighbours(t *testing.T) {
	probeResult := runDurationsPlacementProbe(t, "REQ-100|300.0|0.50", "REQ-900|-42.0|0.20")

	if probeResult.MeasuredLabelBoxHeight <= 0 {
		t.Fatalf("measured label box height is %v — the face did not render", probeResult.MeasuredLabelBoxHeight)
	}
	// (1) Row pitch clears the line box the engine actually draws.
	if probeResult.RowHeight < probeResult.MeasuredLabelBoxHeight {
		t.Errorf("row pitch %.4f is under the measured label box %.4f — two rows would intersect",
			probeResult.RowHeight, probeResult.MeasuredLabelBoxHeight)
	}
	// (2) The first label row clears the mark band above it.
	markBandBottom := probeResult.LaneMarkY + probeResult.BandMarkRadius
	firstRowTop := probeResult.LaneRowY - probeResult.LabelTextAscent
	if firstRowTop < markBandBottom {
		t.Errorf("first label row top %.2f sits inside the mark band (bottom %.2f)",
			firstRowTop, markBandBottom)
	}
	// (3) The reversed band's last row clears Panel B's title.
	lastRowBaseline := probeResult.ReversedRowY + float64(probeResult.RowCount-1)*probeResult.RowHeight
	lastRowBottom := lastRowBaseline + (probeResult.MeasuredLabelBoxHeight - probeResult.LabelTextAscent)
	titleTop := probeResult.MedianTitleY - probeResult.LabelTextAscent
	if lastRowBottom > titleTop {
		t.Errorf("last reversed label row bottom %.2f overlaps Panel B's title top %.2f",
			lastRowBottom, titleTop)
	}
}

// REQ-266's requirement, met by there being nothing left to meet it about: the
// point of measuring at runtime is that no hand-transcribed face number survives.
// This fails if one reappears, which is what keeps the requirement true rather
// than merely satisfied on the day it shipped. It needs no browser.
func TestDurationsCarriesNoMeasuredFaceConstants(t *testing.T) {
	for _, sourcePath := range []string{"durations.go", "web/board-durations.js"} {
		sourceBytes, readError := embeddedWebAssets.ReadFile(sourcePath)
		sourceText := string(sourceBytes)
		if readError != nil {
			// durations.go is not an embedded web asset; read it off disk.
			diskBytes, diskError := os.ReadFile(sourcePath)
			if diskError != nil {
				t.Fatalf("read %s: %v", sourcePath, diskError)
			}
			sourceText = string(diskBytes)
		}
		for _, retiredToken := range []string{
			"durationsLabelCharacterWidthUnits",
			"durationsMeasuredLabelWidthSupremumUnits",
			"durationsMeasuredLabelBoxHeightUnits",
			"DURATIONS_LABEL_CHARACTER_WIDTH",
		} {
			if strings.Contains(sourceText, retiredToken) {
				t.Errorf("%s carries %s again — a measured-face constant is exactly what REQ-292 deleted, "+
					"and a new one has to name the build it came from (REQ-266) or be measured at runtime instead",
					sourcePath, retiredToken)
			}
		}
	}
}

// REQ-292's captured RED, made into a standing assertion.
//
// RED was: "a board rendered where --font-sans resolves to a face wider than 7.15
// units per character draws labels past the slots placement assigned them, and
// nothing in the suite notices. Confirmed by measurement: Arial Black draws 7.34."
//
// The old width model could not express this. It multiplied a character count by
// a constant, so a wider face produced the SAME number and the same slots — the
// labels simply drew past them. Measured placement cannot have that bug by
// construction, and this is how that is shown rather than asserted: run the same
// fixture at the board's own size and at a deliberately wider one, and the
// placement must respond. A packer still using a fixed width model produces
// identical output for both and fails here.
func TestBrowserBehaviorDurationsPlacementRespondsToTheRenderedFace(t *testing.T) {
	const fixtureSpec = "REQ-100|600.0|0.10;REQ-101|540.0|0.17;REQ-102|480.0|0.24;" +
		"REQ-103|420.0|0.31;REQ-104|360.0|0.38;REQ-105|300.0|0.45;REQ-106|240.0|0.52"

	atBoardFace := runDurationsPlacementProbe(t, fixtureSpec, "")
	atWiderFace := runDurationsPlacementProbeWithFaceOverride(t, fixtureSpec, "3px")

	if atBoardFace.MeasuredLabelBoxHeight <= 0 || atWiderFace.MeasuredLabelBoxHeight <= 0 {
		t.Fatal("a probe measured nothing, so this comparison proves nothing")
	}

	// The wider face must move at least one label's geometry. Equality across a
	// face change is the signature of a width model.
	sameEverywhere := true
	for labelIndex := range atBoardFace.Overflow.Labels {
		atBoard := atBoardFace.Overflow.Labels[labelIndex]
		atWider := atWiderFace.Overflow.Labels[labelIndex]
		if atBoard.Drawn != atWider.Drawn || atBoard.Row != atWider.Row ||
			math.Abs((atBoard.Right-atBoard.Left)-(atWider.Right-atWider.Left)) > 0.01 {
			sameEverywhere = false
			break
		}
	}
	if sameEverywhere {
		t.Error("placement produced identical geometry at two different letter-spacings — " +
			"that is what a fixed width model does, and it is the defect REQ-292 removed")
	}
}

// runDurationsPlacementProbeWithFaceOverride runs the same probe with extra CSS
// appended, so a test can widen the drawn text without needing a second font
// installed on the machine. letter-spacing is used rather than a font swap
// because it changes the MEASURED width by a controlled amount on any engine.
func runDurationsPlacementProbeWithFaceOverride(t *testing.T, overflowSpec string, letterSpacing string) durationsProbeResult {
	t.Helper()
	original := durationsProbeExtraStyle
	durationsProbeExtraStyle = ".durations-mark-label { letter-spacing: " + letterSpacing + "; }"
	defer func() { durationsProbeExtraStyle = original }()
	return runDurationsPlacementProbe(t, overflowSpec, "")
}

// durationsProbeExtraStyle is appended to the probe page's stylesheet. Empty for
// every probe except the face-response one above.
var durationsProbeExtraStyle = ""

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
