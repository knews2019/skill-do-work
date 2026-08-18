  // ---- durations ----------------------------------------------------------
  //
  // Three panels, one shared calendar axis, one measure and one scale each —
  // never two y-scales on a single plot. Panel A answers "what is the spread and
  // where are the outliers", B answers "is it getting slower", C answers "how
  // often does this run at all".
  //
  // The x-axis is LINEAR REAL TIME, not one column per active day. The idle
  // gaps are the cadence finding; compressing them would destroy the answer to
  // "how often are these executed".
  //
  // Panel B applies the calibration's read-time rule (a span over four hours is
  // an assumed pause, a negative span is a broken stamp) while panel A still
  // plots both, raw and labelled. The Go side owns that rule — the payload
  // arrives with `excludedReason` already set, so this file never re-derives it.

  var DURATIONS_SVG_NS = "http://www.w3.org/2000/svg";
  var DURATIONS_VIEW_WIDTH = 1200;
  var DURATIONS_MARGIN_LEFT = 54;
  var DURATIONS_MARGIN_RIGHT = 18;
  var DURATIONS_PLOT_WIDTH = DURATIONS_VIEW_WIDTH - DURATIONS_MARGIN_LEFT - DURATIONS_MARGIN_RIGHT;

  // Panel A — the overflow lane above a scale break exists so three long spans
  // cannot squash the 90% of samples that live under 35 minutes.
  var DURATIONS_LANE_TOP = 22;
  var DURATIONS_LANE_MARK_Y = 40;
  var DURATIONS_BREAK_Y = 62;
  var DURATIONS_MAIN_TOP = 70;
  var DURATIONS_MAIN_BOTTOM = 272;
  var DURATIONS_CEILING_MINUTES = 60;
  var DURATIONS_BELOW_ZERO_Y = 284;
  // Direct labels. WHICH marks get one is decided in durations.go and arrives in
  // the payload as labelRow/labelAnchor; what lives here is only the geometry
  // that turns a row index into a baseline. The gap and the row count are shared
  // with that decision, so TestDurationLabelGeometryMatchesTheRenderer pins this
  // file's constants against the Go ones.
  var DURATIONS_LABEL_ROW_COUNT = 2;
  var DURATIONS_LABEL_GAP = 9;
  var DURATIONS_LABEL_ROW_HEIGHT = 12;
  var DURATIONS_LANE_LABEL_ROW_Y = 44;
  var DURATIONS_REVERSED_LABEL_ROW_Y = 288;
  // Panel B — median minutes per active day.
  var DURATIONS_MEDIAN_TITLE_Y = 316;
  var DURATIONS_MEDIAN_TOP = 334;
  var DURATIONS_MEDIAN_BOTTOM = 414;
  var DURATIONS_MEDIAN_CEILING = 45;
  // A day over the ceiling is drawn as a full-height bar plus a detached sliver
  // above it, so the break reads as "continues above" rather than as a value.
  // Every over-ceiling day gets one, not only the slowest.
  var DURATIONS_MEDIAN_OVER_CEILING_GAP = 6;
  var DURATIONS_MEDIAN_OVER_CEILING_HEIGHT = 3;
  // Panel C — REQs completed per day.
  var DURATIONS_COUNT_TITLE_Y = 450;
  var DURATIONS_COUNT_TOP = 468;
  var DURATIONS_COUNT_BOTTOM = 538;
  var DURATIONS_AXIS_LABEL_Y = 556;
  var DURATIONS_VIEW_HEIGHT = 570;

  var DURATIONS_MONTHS = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];

  function durationRouteFill(route) {
    if (route === "A" || route === "B" || route === "C") {
      return "var(--route-" + route.toLowerCase() + ")";
    }
    return "var(--route-none)";
  }

  function durationRouteName(route) {
    if (route === "A" || route === "B" || route === "C") {
      return "Route " + route;
    }
    return "No route recorded";
  }

  // Durations are read as magnitudes, so a negative span reads as "−30 min"
  // rather than as a subtraction. Long spans switch to hours so the overflow
  // labels stay short.
  function formatDurationMinutes(minutes) {
    var sign = minutes < 0 ? "−" : "";
    var magnitude = Math.abs(minutes);
    if (magnitude < 60) {
      return sign + magnitude.toFixed(1) + " min";
    }
    return sign + Math.floor(magnitude / 60) + "h " + Math.round(magnitude % 60) + "m";
  }

  function formatDurationDayLabel(epochMs) {
    var instant = new Date(epochMs);
    return instant.getUTCDate() + " " + DURATIONS_MONTHS[instant.getUTCMonth()];
  }

  function formatDurationStamp(epochMs) {
    return new Date(epochMs).toISOString().replace("T", " ").slice(0, 16) + " UTC";
  }

  // The direct-label verdict is the payload's, not this file's: `labelRow` is the
  // text row inside the sample's own band, and -1 when the collision rule could
  // not place a label at all. Returns null for a sample that gets none, so the
  // renderer has exactly one place where "is this labelled" is answered.
  function durationsLabelBaselineY(sample) {
    var labelRow = sample.labelRow;
    if (typeof labelRow !== "number" || labelRow < 0 || labelRow >= DURATIONS_LABEL_ROW_COUNT) {
      return null;
    }
    var bandRowY =
      sample.wallMinutes < 0 ? DURATIONS_REVERSED_LABEL_ROW_Y : DURATIONS_LANE_LABEL_ROW_Y;
    return bandRowY + labelRow * DURATIONS_LABEL_ROW_HEIGHT;
  }

  // The remainder sentence goes on the band's LAST text row. The marks sit level
  // with the first row, so a sentence there is legible only while the band is
  // sparse — which is exactly when there is no remainder to print. Placement
  // reserves this row's right edge to match.
  function durationsRemainderBaselineY(bandRowY) {
    return bandRowY + (DURATIONS_LABEL_ROW_COUNT - 1) * DURATIONS_LABEL_ROW_HEIGHT;
  }

  function makeDurationsSvgNode(svg, name, attributes, textContent) {
    var node = document.createElementNS(DURATIONS_SVG_NS, name);
    Object.keys(attributes).forEach(function (key) {
      node.setAttribute(key, attributes[key]);
    });
    if (textContent !== undefined) {
      node.appendChild(document.createTextNode(textContent));
    }
    svg.appendChild(node);
    return node;
  }

  function renderDurationsView() {
    var chartHost = document.getElementById("durations-chart");
    var summaryNode = document.getElementById("durations-summary");
    var readoutNode = document.getElementById("durations-readout");
    var tableBody = document.getElementById("durations-table-body");
    if (!chartHost || !summaryNode || !tableBody) {
      return;
    }

    var durations = boardData.durations || {};
    var samples = durations.samples || [];
    var days = durations.days || [];

    chartHost.textContent = "";
    tableBody.textContent = "";

    if (samples.length === 0) {
      summaryNode.textContent =
        "No archived REQ carries both a claim and a completion stamp yet, so there is nothing to measure.";
      return;
    }

    var excludedSamples = samples.filter(function (sample) {
      return sample.excludedReason;
    });
    summaryNode.textContent =
      samples.length +
      " archived REQ" +
      (samples.length === 1 ? "" : "s") +
      " with both stamps, across " +
      days.length +
      " active day" +
      (days.length === 1 ? "" : "s") +
      ". Panel B excludes " +
      excludedSamples.length +
      " span" +
      (excludedSamples.length === 1 ? "" : "s") +
      " from its medians (over four hours is an assumed pause, negative is a broken stamp); panel A still plots them.";

    var svg = document.createElementNS(DURATIONS_SVG_NS, "svg");
    svg.setAttribute("viewBox", "0 0 " + DURATIONS_VIEW_WIDTH + " " + DURATIONS_VIEW_HEIGHT);
    svg.setAttribute("class", "durations-svg");
    svg.setAttribute("role", "img");
    chartHost.appendChild(svg);

    var sampleTimes = samples.map(function (sample) {
      return Date.parse(sample.completionTime);
    });
    var timeStart = Math.min.apply(null, sampleTimes);
    var timeEnd = Math.max.apply(null, sampleTimes);
    var timeSpan = timeEnd - timeStart || 1;

    svg.setAttribute(
      "aria-label",
      "Three stacked panels sharing a calendar axis from " +
        formatDurationDayLabel(timeStart) +
        " to " +
        formatDurationDayLabel(timeEnd) +
        ". Panel A plots each archived REQ's duration in minutes coloured by route. Panel B plots the median minutes per active day. Panel C counts REQs completed per day. Every value is also listed in the table below."
    );

    function xOfEpoch(epochMs) {
      return DURATIONS_MARGIN_LEFT + ((epochMs - timeStart) / timeSpan) * DURATIONS_PLOT_WIDTH;
    }
    function yOfMinutes(minutes) {
      var clamped = Math.min(minutes, DURATIONS_CEILING_MINUTES);
      return (
        DURATIONS_MAIN_BOTTOM -
        (clamped / DURATIONS_CEILING_MINUTES) * (DURATIONS_MAIN_BOTTOM - DURATIONS_MAIN_TOP)
      );
    }
    function yOfDayMedian(minutes) {
      var clamped = Math.min(minutes, DURATIONS_MEDIAN_CEILING);
      return (
        DURATIONS_MEDIAN_BOTTOM -
        (clamped / DURATIONS_MEDIAN_CEILING) * (DURATIONS_MEDIAN_BOTTOM - DURATIONS_MEDIAN_TOP)
      );
    }
    function gridRow(y, isBaseline, label) {
      makeDurationsSvgNode(svg, "line", {
        x1: DURATIONS_MARGIN_LEFT,
        y1: y,
        x2: DURATIONS_VIEW_WIDTH - DURATIONS_MARGIN_RIGHT,
        y2: y,
        class: isBaseline ? "durations-axis-line" : "durations-grid-line"
      });
      makeDurationsSvgNode(
        svg,
        "text",
        { x: DURATIONS_MARGIN_LEFT - 8, y: y + 4, class: "durations-tick", "text-anchor": "end" },
        label
      );
    }

    // ---- panel A: one mark per REQ ----
    makeDurationsSvgNode(
      svg,
      "text",
      { x: DURATIONS_MARGIN_LEFT, y: 12, class: "durations-axis-title" },
      "A · Duration per REQ · minutes"
    );
    makeDurationsSvgNode(svg, "rect", {
      x: DURATIONS_MARGIN_LEFT,
      y: DURATIONS_LANE_TOP,
      width: DURATIONS_PLOT_WIDTH,
      height: DURATIONS_BREAK_Y - DURATIONS_LANE_TOP - 4,
      class: "durations-overflow-lane",
      rx: 2
    });
    makeDurationsSvgNode(
      svg,
      "text",
      {
        x: DURATIONS_MARGIN_LEFT - 8,
        y: DURATIONS_LANE_MARK_Y + 4,
        class: "durations-tick",
        "text-anchor": "end"
      },
      "60+"
    );
    makeDurationsSvgNode(svg, "line", {
      x1: DURATIONS_MARGIN_LEFT,
      y1: DURATIONS_BREAK_Y,
      x2: DURATIONS_VIEW_WIDTH - DURATIONS_MARGIN_RIGHT,
      y2: DURATIONS_BREAK_Y,
      class: "durations-grid-line"
    });
    [0, 15, 30, 45, 60].forEach(function (minutes) {
      gridRow(yOfMinutes(minutes), minutes === 0, String(minutes));
    });

    var markIndex = [];
    samples.forEach(function (sample) {
      var epochMs = Date.parse(sample.completionTime);
      var minutes = sample.wallMinutes;
      var isOverflow = minutes > DURATIONS_CEILING_MINUTES;
      var isReversed = minutes < 0;
      var markX = xOfEpoch(epochMs);
      var markY = isOverflow
        ? DURATIONS_LANE_MARK_Y
        : isReversed
        ? DURATIONS_BELOW_ZERO_Y
        : yOfMinutes(minutes);
      makeDurationsSvgNode(svg, "circle", {
        cx: markX.toFixed(1),
        cy: markY.toFixed(1),
        r: isOverflow || isReversed ? 5 : 4,
        fill: isReversed ? "var(--durations-critical)" : durationRouteFill(sample.route),
        class: "durations-mark"
      });
      markIndex.push({ x: markX, y: markY, sample: sample, epochMs: epochMs });
    });

    // Direct labels only where a mark carries a value its y cannot: the overflow
    // lane, where every mark sits at one y, and the reversed band. Which marks
    // get one is the payload's answer — `labelRow` is the text row inside that
    // mark's band and -1 when the collision rule could not place a label, and
    // `labelAnchor` is the side it was placed on. Both bands are packed
    // independently there; nothing here recomputes either.
    markIndex.forEach(function (mark) {
      var baselineY = durationsLabelBaselineY(mark.sample);
      if (baselineY === null) {
        return;
      }
      var anchorsBeforeMark = mark.sample.labelAnchor === "end";
      makeDurationsSvgNode(
        svg,
        "text",
        {
          x: (mark.x + (anchorsBeforeMark ? -DURATIONS_LABEL_GAP : DURATIONS_LABEL_GAP)).toFixed(1),
          y: baselineY.toFixed(1),
          class: "durations-mark-label",
          "text-anchor": anchorsBeforeMark ? "end" : "start"
        },
        mark.sample.id + " " + formatDurationMinutes(mark.sample.wallMinutes)
      );
    });

    // Whatever placement could not fit is stated, never dropped in silence: the
    // count is what stops a reader taking the visible labels for all of them.
    var durationLabelCounts = durations.labels || {};
    function drawDurationsRemainder(hiddenCount, bandRowY, remainderTail) {
      if (!hiddenCount) {
        return;
      }
      makeDurationsSvgNode(
        svg,
        "text",
        {
          x: DURATIONS_VIEW_WIDTH - DURATIONS_MARGIN_RIGHT,
          y: bandRowY,
          class: "durations-tick",
          "text-anchor": "end"
        },
        "+" + hiddenCount + " more " + remainderTail
      );
    }
    drawDurationsRemainder(
      durationLabelCounts.overflowHiddenCount,
      durationsRemainderBaselineY(DURATIONS_LANE_LABEL_ROW_Y),
      "over " + DURATIONS_CEILING_MINUTES + " min"
    );
    drawDurationsRemainder(
      durationLabelCounts.reversedHiddenCount,
      durationsRemainderBaselineY(DURATIONS_REVERSED_LABEL_ROW_Y),
      "reversed"
    );

    // ---- panel B: median minutes per active day ----
    makeDurationsSvgNode(
      svg,
      "text",
      { x: DURATIONS_MARGIN_LEFT, y: DURATIONS_MEDIAN_TITLE_Y, class: "durations-axis-title" },
      "B · Median minutes per active day · paused and broken spans excluded"
    );
    [0, 15, 30, 45].forEach(function (minutes) {
      gridRow(
        yOfDayMedian(minutes),
        minutes === 0,
        minutes === DURATIONS_MEDIAN_CEILING ? minutes + "+" : String(minutes)
      );
    });

    var dayWidth = (DURATIONS_PLOT_WIDTH * 86400000) / timeSpan;
    var barWidth = Math.max(4, Math.min(24, dayWidth - 2));
    var slowestDay = null;
    days.forEach(function (day) {
      if (!day.hasMedian) {
        return;
      }
      var dayEpochMs = Date.parse(day.dayTime);
      var barTop = yOfDayMedian(day.medianMinutes);
      makeDurationsSvgNode(svg, "rect", {
        x: (xOfEpoch(dayEpochMs) - barWidth / 2).toFixed(1),
        y: barTop.toFixed(1),
        width: barWidth.toFixed(1),
        height: Math.max(2, DURATIONS_MEDIAN_BOTTOM - barTop).toFixed(1),
        rx: 3,
        class: "durations-bar"
      });
      if (day.medianMinutes > DURATIONS_MEDIAN_CEILING) {
        makeDurationsSvgNode(svg, "rect", {
          x: (xOfEpoch(dayEpochMs) - barWidth / 2).toFixed(1),
          y: (DURATIONS_MEDIAN_TOP - DURATIONS_MEDIAN_OVER_CEILING_GAP).toFixed(1),
          width: barWidth.toFixed(1),
          height: DURATIONS_MEDIAN_OVER_CEILING_HEIGHT,
          rx: 1,
          class: "durations-bar durations-bar-over-ceiling"
        });
      }
      if (slowestDay === null || day.medianMinutes > slowestDay.medianMinutes) {
        slowestDay = day;
      }
    });
    if (slowestDay) {
      makeDurationsSvgNode(
        svg,
        "text",
        {
          x: xOfEpoch(Date.parse(slowestDay.dayTime)).toFixed(1),
          y: (
            yOfDayMedian(slowestDay.medianMinutes) -
            (slowestDay.medianMinutes > DURATIONS_MEDIAN_CEILING
              ? DURATIONS_MEDIAN_OVER_CEILING_GAP + 7
              : 7)
          ).toFixed(1),
          class: "durations-mark-label",
          "text-anchor": "middle"
        },
        slowestDay.medianMinutes.toFixed(0) + " min"
      );
    }

    // ---- panel C: REQs completed per day ----
    makeDurationsSvgNode(
      svg,
      "text",
      { x: DURATIONS_MARGIN_LEFT, y: DURATIONS_COUNT_TITLE_Y, class: "durations-axis-title" },
      "C · REQs completed per day · every sample counted"
    );
    var peakCount = days.reduce(function (highest, day) {
      return Math.max(highest, day.completedCount);
    }, 1);
    makeDurationsSvgNode(svg, "line", {
      x1: DURATIONS_MARGIN_LEFT,
      y1: DURATIONS_COUNT_BOTTOM,
      x2: DURATIONS_VIEW_WIDTH - DURATIONS_MARGIN_RIGHT,
      y2: DURATIONS_COUNT_BOTTOM,
      class: "durations-axis-line"
    });
    makeDurationsSvgNode(
      svg,
      "text",
      {
        x: DURATIONS_MARGIN_LEFT - 8,
        y: DURATIONS_COUNT_TOP + 4,
        class: "durations-tick",
        "text-anchor": "end"
      },
      String(peakCount)
    );
    days.forEach(function (day) {
      var dayEpochMs = Date.parse(day.dayTime);
      var columnHeight = (day.completedCount / peakCount) * (DURATIONS_COUNT_BOTTOM - DURATIONS_COUNT_TOP);
      makeDurationsSvgNode(svg, "rect", {
        x: (xOfEpoch(dayEpochMs) - barWidth / 2).toFixed(1),
        y: (DURATIONS_COUNT_BOTTOM - columnHeight).toFixed(1),
        width: barWidth.toFixed(1),
        height: Math.max(2, columnHeight).toFixed(1),
        rx: 3,
        class: "durations-bar durations-bar-count"
      });
    });

    // ---- shared calendar axis ----
    // Month ticks, with any tick that would collide with the first or last
    // label dropped rather than overprinted.
    makeDurationsSvgNode(
      svg,
      "text",
      { x: DURATIONS_MARGIN_LEFT, y: DURATIONS_AXIS_LABEL_Y, class: "durations-tick" },
      formatDurationDayLabel(timeStart)
    );
    makeDurationsSvgNode(
      svg,
      "text",
      {
        x: DURATIONS_VIEW_WIDTH - DURATIONS_MARGIN_RIGHT,
        y: DURATIONS_AXIS_LABEL_Y,
        class: "durations-tick",
        "text-anchor": "end"
      },
      formatDurationDayLabel(timeEnd)
    );
    var firstMonth = new Date(timeStart);
    var monthCursor = Date.UTC(firstMonth.getUTCFullYear(), firstMonth.getUTCMonth() + 1, 1);
    while (monthCursor < timeEnd) {
      var monthX = xOfEpoch(monthCursor);
      if (monthX - DURATIONS_MARGIN_LEFT > 42 && DURATIONS_VIEW_WIDTH - DURATIONS_MARGIN_RIGHT - monthX > 42) {
        makeDurationsSvgNode(svg, "line", {
          x1: monthX.toFixed(1),
          y1: DURATIONS_MAIN_TOP,
          x2: monthX.toFixed(1),
          y2: DURATIONS_COUNT_BOTTOM,
          class: "durations-month-line"
        });
        makeDurationsSvgNode(
          svg,
          "text",
          { x: monthX.toFixed(1), y: DURATIONS_AXIS_LABEL_Y, class: "durations-tick", "text-anchor": "middle" },
          formatDurationDayLabel(monthCursor)
        );
      }
      var cursorInstant = new Date(monthCursor);
      monthCursor = Date.UTC(cursorInstant.getUTCFullYear(), cursorInstant.getUTCMonth() + 1, 1);
    }

    // ---- hover readout ----
    // A nearest-mark layer over the whole plot, because 8px dots are not a
    // usable hit target and the reader is pointing at a region, not a pixel.
    var hoverSurface = makeDurationsSvgNode(svg, "rect", {
      x: DURATIONS_MARGIN_LEFT,
      y: 0,
      width: DURATIONS_PLOT_WIDTH,
      height: DURATIONS_VIEW_HEIGHT,
      fill: "transparent",
      class: "durations-hover-surface"
    });
    var dayByKey = {};
    days.forEach(function (day) {
      dayByKey[day.dayKey] = day;
    });

    function describeAtPointer(pointerX, pointerY) {
      if (pointerY <= DURATIONS_MEDIAN_TITLE_Y - 12) {
        var nearestMark = null;
        var nearestDistance = Infinity;
        markIndex.forEach(function (mark) {
          var distance = Math.abs(mark.x - pointerX) + Math.abs(mark.y - pointerY) * 0.35;
          if (distance < nearestDistance) {
            nearestDistance = distance;
            nearestMark = mark;
          }
        });
        if (!nearestMark) {
          return "";
        }
        var sample = nearestMark.sample;
        var note = sample.excludedReason
          ? " · excluded from day medians (" +
            (sample.excludedReason === "paused" ? "assumed paused session" : "reversed stamp") +
            ")"
          : "";
        return (
          sample.id +
          " · " +
          durationRouteName(sample.route) +
          " · " +
          formatDurationMinutes(sample.wallMinutes) +
          " · " +
          formatDurationStamp(nearestMark.epochMs) +
          note
        );
      }

      var nearestDay = null;
      var nearestDayDistance = Infinity;
      days.forEach(function (day) {
        var distance = Math.abs(xOfEpoch(Date.parse(day.dayTime)) - pointerX);
        if (distance < nearestDayDistance) {
          nearestDayDistance = distance;
          nearestDay = day;
        }
      });
      if (!nearestDay) {
        return "";
      }
      var medianText = nearestDay.hasMedian
        ? "median " + nearestDay.medianMinutes.toFixed(1) + " min over " + nearestDay.keptCount + " sample" +
          (nearestDay.keptCount === 1 ? "" : "s")
        : "no median — every span that day was excluded by the read-time rule";
      return (
        nearestDay.dayKey +
        " · " +
        nearestDay.completedCount +
        " REQ" +
        (nearestDay.completedCount === 1 ? "" : "s") +
        " completed · " +
        medianText
      );
    }

    function updateReadout(event) {
      if (!readoutNode) {
        return;
      }
      var bounds = svg.getBoundingClientRect();
      if (bounds.width === 0) {
        return;
      }
      var scale = DURATIONS_VIEW_WIDTH / bounds.width;
      var pointerX = (event.clientX - bounds.left) * scale;
      var pointerY = (event.clientY - bounds.top) * scale;
      readoutNode.textContent = describeAtPointer(pointerX, pointerY);
    }
    hoverSurface.addEventListener("mousemove", updateReadout);
    hoverSurface.addEventListener("mouseleave", function () {
      if (readoutNode) {
        readoutNode.textContent = "";
      }
    });

    // ---- table view ----
    // Every value the chart shows is reachable without a pointer.
    samples
      .slice()
      .reverse()
      .forEach(function (sample) {
        var row = document.createElement("tr");
        [
          sample.id,
          sample.route || "—",
          formatDurationStamp(Date.parse(sample.completionTime)),
          formatDurationMinutes(sample.wallMinutes),
          sample.excludedReason === "paused"
            ? "excluded from day median — assumed paused session"
            : sample.excludedReason === "reversed"
            ? "excluded from day median — reversed stamp"
            : ""
        ].forEach(function (cellText) {
          var cell = document.createElement("td");
          cell.textContent = cellText;
          row.appendChild(cell);
        });
        tableBody.appendChild(row);
      });
  }
