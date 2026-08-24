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
  var DURATIONS_DAY_MS = 86400000;

  // Panel A — the overflow lane above a scale break exists so long spans cannot
  // squash the ordinary samples. Every mark stays in the lane; the complete
  // ranked text record now lives in adjacent HTML rather than competing for SVG
  // width inside the chart.
  var DURATIONS_LANE_TOP = 22;
  var DURATIONS_LANE_MARK_Y = 32;
  var DURATIONS_BAND_MARK_RADIUS = 5;
  var DURATIONS_BREAK_Y = 76;
  var DURATIONS_MAIN_TOP = 84;
  var DURATIONS_MAIN_BOTTOM = 286;
  var DURATIONS_CEILING_MINUTES = 60;
  var DURATIONS_ORDINARY_MARK_OPACITY = 0.62;
  var DURATIONS_BELOW_ZERO_Y = 298;
  // An axis tick's baseline sits this far below the y it labels, which is what
  // optically centres an 11px face on that line. A named number because a tick
  // row is a neighbour some other text has to clear, and a test cannot read a
  // literal buried in an attribute.
  var DURATIONS_TICK_BASELINE_DROP = 4;
  // The UR grouping lane (REQ-346) — one bracket per user request, spanning its
  // samples' completion times, under panel A and on panel A's axis.
  //
  // It carries NO text inside the lane. The overflow lane's own history is why:
  // variable-width text at this density is the thing that stopped paying for
  // itself, and a bracket here is narrower than one there. Identity is
  // positional plus the hover; the only words in this block sit in the y-axis
  // gutter and on the title line, where nothing competes for the space.
  var DURATIONS_UR_LANE_TITLE_Y = 352;
  var DURATIONS_UR_LANE_TOP = 358;
  // Six rows, and the row count was chosen against the DENSE board rather than
  // this repository's. Unplaced URs by row count, measured on two generated
  // boards: this repository (65 URs, 87-day axis) 3 rows 14, 4 rows 7, 6 rows 2;
  // a synthetic board at the stated target of 692 samples over 47 active days
  // (140 URs, 90-day axis) 3 rows 55, 4 rows 36, 6 rows 12, 8 rows 2. Four rows
  // is comfortable here and leaves a quarter of the target board's URs in the
  // remainder, which is the wrong board to be comfortable on.
  //
  // Six rows at this pitch costs 67 units of a 697-unit view — a third more lane
  // for a third of the remainder. Eight would place nearly all of them and take
  // 87 units, at which point the lane is nearly half the height of the plot it
  // annotates and stops being a lane.
  var DURATIONS_UR_LANE_ROW_COUNT = 6;
  var DURATIONS_UR_LANE_ROW_PITCH = 10;
  var DURATIONS_UR_BRACKET_HEIGHT = 7;
  // A UR whose samples all completed inside one day is a bracket of nearly zero
  // width on a linear real-time axis — 46 of this repository's 65 URs span under
  // 3 units. The minimum is what keeps such a UR visible as a mark rather than
  // vanishing, and the separation is what keeps two adjacent ones reading as two.
  var DURATIONS_UR_BRACKET_MIN_WIDTH = 3;
  var DURATIONS_UR_BRACKET_SEPARATION = 2;
  // The samples whose REQ carries no `user_request` get one reserved row below
  // the packed ones rather than competing for a row, because they are not a
  // group: they are the absence of one. Keeping the row reserved is what makes
  // "no UR at all" and "no row was free" two different facts on screen.
  var DURATIONS_UR_UNKNOWN_ROW_TOP = DURATIONS_UR_LANE_TOP + DURATIONS_UR_LANE_ROW_COUNT * DURATIONS_UR_LANE_ROW_PITCH;
  var DURATIONS_UR_LANE_BOTTOM = DURATIONS_UR_UNKNOWN_ROW_TOP + DURATIONS_UR_BRACKET_HEIGHT;
  // Panel B — median minutes per active day.
  var DURATIONS_MEDIAN_TITLE_Y = 443;
  var DURATIONS_MEDIAN_TOP = 461;
  var DURATIONS_MEDIAN_BOTTOM = 541;
  var DURATIONS_MEDIAN_CEILING = 45;
  // A day over the ceiling is drawn as a full-height bar plus a detached sliver
  // above it, so the break reads as "continues above" rather than as a value.
  // Every over-ceiling day gets one, not only the slowest.
  var DURATIONS_MEDIAN_OVER_CEILING_GAP = 6;
  var DURATIONS_MEDIAN_OVER_CEILING_HEIGHT = 3;
  // The slowest day's annotation goes BELOW the panel's baseline, centred under
  // its bar. Above the bar is not available: the tallest bar's top is
  // DURATIONS_MEDIAN_TOP and the strip over it belongs to this panel's title,
  // which the annotation used to print straight through whenever the slowest day
  // fell under the title's text (REQ-242 — it read "209 miB · Median minutes…"
  // on a fixture whose leftmost day was slowest). Inside the plot is not
  // available either: at a dense day count the bars are 4 units wide and a label
  // there would overprint its neighbours.
  //
  // That leaves the strip under the baseline, and it has three occupants of its
  // own. Two are cleared: the "0" tick, whose baseline is
  // DURATIONS_MEDIAN_BOTTOM + DURATIONS_TICK_BASELINE_DROP, and panel C's
  // title. This number centres the annotation's text box between them, so the
  // clearance is the same at every x and every bar height — which is the whole
  // point.
  //
  // The third is the month rule, and it is an ACCEPTED crossing rather than a
  // cleared one: .durations-month-line spans DURATIONS_MAIN_TOP to
  // DURATIONS_COUNT_BOTTOM, so no baseline in this strip can avoid it, and on a
  // fixture whose slowest day falls on a month boundary it passes between the
  // "9" and the " min". A one-unit soft rule through a label is the same
  // crossing panel A's reversed-band labels already take, so it is accepted —
  // and the test asserts it STAYS one unit and soft, because that is the only
  // reason it is acceptable.
  //
  // All four facts are asserted, because the original defect was invisible here
  // for no better reason than where this repository's slowest day happens to
  // fall.
  var DURATIONS_MEDIAN_ANNOTATION_BASELINE_Y = 560;
  // Panel C — REQs completed per day.
  var DURATIONS_COUNT_TITLE_Y = 577;
  var DURATIONS_COUNT_TOP = 595;
  var DURATIONS_COUNT_BOTTOM = 665;
  var DURATIONS_AXIS_LABEL_Y = 683;
  var DURATIONS_VIEW_HEIGHT = 697;

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

  // Panel A's route fills predate the chooser and remain its default. UR and
  // domain are categorical, so their names are sorted before colours are
  // assigned: the same payload always gives a category the same colour. Twelve
  // named URs are the usable palette; every later ID deliberately shares Other
  // URs rather than wrapping silently into an indistinguishable category.
  var durationsColourChannel = "route";
  var durationsWindow = "30";
  var DURATIONS_USER_REQUEST_PALETTE_CAPACITY = 12;
  var DURATIONS_CATEGORICAL_FILLS = [
    "var(--durations-category-1)",
    "var(--durations-category-2)",
    "var(--durations-category-3)",
    "var(--durations-category-4)",
    "var(--durations-category-5)",
    "var(--durations-category-6)",
    "var(--durations-category-7)",
    "var(--durations-category-8)",
    "var(--durations-category-9)",
    "var(--durations-category-10)",
    "var(--durations-category-11)",
    "var(--durations-category-12)"
  ];

  function setDurationsColourChannel(colourChannel) {
    if (colourChannel === "route" || colourChannel === "user-request" || colourChannel === "domain") {
      durationsColourChannel = colourChannel;
    }
  }

  function setDurationsWindow(windowName) {
    if (windowName === "30" || windowName === "90" || windowName === "all") {
      durationsWindow = windowName;
    }
  }

  function durationsWindowName() {
    if (durationsWindow === "90") {
      return "Last 90 days";
    }
    if (durationsWindow === "all") {
      return "All history";
    }
    return "Last 30 days";
  }

  function durationSampleRequest(sample) {
    var requestsById = (typeof boardData === "object" && boardData && boardData.requests) || {};
    return requestsById[sample.id] || null;
  }

  function durationSampleDomain(sample) {
    var request = durationSampleRequest(sample);
    return (request && request.domain) || "";
  }

  // ---- the REQ-to-UR join (REQ-346) ---------------------------------------
  // The join is CLIENT-SIDE and stays that way. `boardData.requests[id]` already
  // carries `userRequestId` for every parsed ticket, including the archived ones
  // panel A plots, so putting a second copy in the durations payload would make
  // one fact have two definitions that can disagree — the thing REQ-219 settled
  // in the other direction only because the rule's VERDICT had no other home.
  // This one does. The lookup is a plain map hit per sample, run once per render.
  //
  // A sample whose REQ carries no `user_request` is a real state, not a gap to
  // paper over: REQ-001 through REQ-011 and REQ-060 pre-date the UR system and
  // buildDurationAggregate measures every one of them. They are NAMED on all
  // three surfaces with this one string, never left blank and never given some
  // arbitrary default UR.
  var DURATIONS_UNKNOWN_USER_REQUEST_NAME = "no UR recorded";

  function durationSampleUserRequestId(sample) {
    var request = durationSampleRequest(sample);
    return (request && request.userRequestId) || "";
  }

  function durationUserRequestName(userRequestId) {
    return userRequestId || DURATIONS_UNKNOWN_USER_REQUEST_NAME;
  }

  function durationColourChannelName() {
    if (durationsColourChannel === "user-request") {
      return "UR";
    }
    if (durationsColourChannel === "domain") {
      return "domain";
    }
    return "route";
  }

  function durationColourContext(samples) {
    var userRequestIds = {};
    var domains = {};
    samples.forEach(function (sample) {
      var userRequestId = durationSampleUserRequestId(sample);
      var domain = durationSampleDomain(sample);
      if (userRequestId) {
        userRequestIds[userRequestId] = true;
      }
      if (domain) {
        domains[domain] = true;
      }
    });
    return {
      userRequestIds: Object.keys(userRequestIds).sort(),
      domains: Object.keys(domains).sort()
    };
  }

  function durationCategoricalFill(categoryName, categoryNames, usesUserRequestOverflow) {
    var categoryIndex = categoryNames.indexOf(categoryName);
    if (usesUserRequestOverflow && categoryIndex >= DURATIONS_USER_REQUEST_PALETTE_CAPACITY) {
      return "var(--durations-category-other)";
    }
    return DURATIONS_CATEGORICAL_FILLS[categoryIndex % DURATIONS_CATEGORICAL_FILLS.length];
  }

  function durationMarkColour(sample, colourContext) {
    if (durationsColourChannel === "user-request") {
      var userRequestId = durationSampleUserRequestId(sample);
      if (!userRequestId) {
        return { fill: "none", className: "durations-mark-unknown", label: DURATIONS_UNKNOWN_USER_REQUEST_NAME };
      }
      var isOtherUserRequest = colourContext.userRequestIds.indexOf(userRequestId) >= DURATIONS_USER_REQUEST_PALETTE_CAPACITY;
      return {
        fill: durationCategoricalFill(userRequestId, colourContext.userRequestIds, true),
        className: "",
        label: isOtherUserRequest ? "Other URs" : userRequestId
      };
    }
    if (durationsColourChannel === "domain") {
      var domain = durationSampleDomain(sample);
      if (!domain) {
        return { fill: "none", className: "durations-mark-unknown", label: "No domain recorded" };
      }
      return {
        fill: durationCategoricalFill(domain, colourContext.domains, false),
        className: "",
        label: domain
      };
    }
    if (!sample.route) {
      return { fill: "none", className: "durations-mark-unknown", label: "No route recorded" };
    }
    return { fill: durationRouteFill(sample.route), className: "", label: durationRouteName(sample.route) };
  }

  function durationsColourLegendText() {
    if (durationsColourChannel === "user-request") {
      return "Colour by UR. The first 12 UR IDs by name use distinct colours; later URs share Other URs. No UR recorded is outlined. Reversed stamps stay critical red.";
    }
    if (durationsColourChannel === "domain") {
      return "Colour by domain. Domains use colours by name. No domain recorded is outlined. Reversed stamps stay critical red.";
    }
    return "Colour by route. Route A, Route B, and Route C keep their route colours. No route recorded is outlined. Reversed stamps stay critical red.";
  }

  // ---- UR lane packing -----------------------------------------------------
  // Left-edge order, first free row, counted remainder. The order is the whole
  // of the packing quality: greedy by left edge is the interval-graph colouring
  // walk, and it beat both alternatives tried on this repository's data at every
  // row count (65 URs, 4 rows: 7 unplaced by left edge, 9 by widest-first, 9 by
  // most-samples-first). Ties break on the UR id so a board with two URs
  // starting in the same minute still packs the same way twice.
  //
  // A bracket that finds no row is COUNTED, never dropped in silence — the same
  // rule the overflow lane's labels follow, for the same reason: the rows a
  // reader can see must not be mistaken for all of them.
  function packDurationsUserRequestLane(brackets) {
    var ordered = brackets.slice().sort(function (first, second) {
      if (first.left !== second.left) {
        return first.left - second.left;
      }
      return first.userRequestId < second.userRequestId ? -1 : 1;
    });
    var occupied = [];
    for (var rowIndex = 0; rowIndex < DURATIONS_UR_LANE_ROW_COUNT; rowIndex += 1) {
      occupied.push([]);
    }
    var placements = [];
    var hiddenCount = 0;
    ordered.forEach(function (bracket) {
      for (var rowIndex = 0; rowIndex < DURATIONS_UR_LANE_ROW_COUNT; rowIndex += 1) {
        if (durationsUserRequestRowIsBlocked(occupied[rowIndex], bracket)) {
          continue;
        }
        occupied[rowIndex].push({ left: bracket.left, right: bracket.right });
        placements.push({ bracket: bracket, laneRow: rowIndex });
        return;
      }
      hiddenCount += 1;
    });
    return { placements: placements, hiddenCount: hiddenCount };
  }

  // Deliberately not durationsSpanIsBlocked: that one belongs to the label
  // planner and carries the label separation, and widening it to take a gap
  // argument would re-point every one of its calls to prove this one's point.
  function durationsUserRequestRowIsBlocked(occupiedRow, bracket) {
    for (var placedIndex = 0; placedIndex < occupiedRow.length; placedIndex += 1) {
      var placed = occupiedRow[placedIndex];
      if (bracket.left < placed.right + DURATIONS_UR_BRACKET_SEPARATION &&
        placed.left < bracket.right + DURATIONS_UR_BRACKET_SEPARATION) {
        return true;
      }
    }
    return false;
  }

  // One bracket per UR, from its first sample's completion to its last, plus the
  // single bucket for samples whose REQ names no UR. The bucket is returned
  // apart from the brackets because it never enters the packer: it holds a
  // reserved row of its own (see DURATIONS_UR_UNKNOWN_ROW_TOP).
  function buildDurationsUserRequestBrackets(samples, plotXOfSample) {
    var samplesByUserRequest = {};
    var unknownSamples = [];
    samples.forEach(function (sample) {
      var completionMs = Date.parse(sample.completionTime);
      var sampleX = plotXOfSample(sample);
      var userRequestId = durationSampleUserRequestId(sample);
      if (!userRequestId) {
        unknownSamples.push({ completionMs: completionMs, sampleX: sampleX });
        return;
      }
      var group = samplesByUserRequest[userRequestId];
      if (!group) {
        samplesByUserRequest[userRequestId] = {
          firstMs: completionMs,
          lastMs: completionMs,
          left: sampleX,
          right: sampleX,
          sampleCount: 1
        };
        return;
      }
      group.firstMs = Math.min(group.firstMs, completionMs);
      group.lastMs = Math.max(group.lastMs, completionMs);
      group.left = Math.min(group.left, sampleX);
      group.right = Math.max(group.right, sampleX);
      group.sampleCount += 1;
    });

    // The minimum width can push a bracket past the plot's right edge — the axis
    // domain ends at the midnight AFTER the last completion, so a UR completing
    // late on the last day sits under 2 units from that edge and the widening
    // carries it over. The render caught it on this repository at 65 URs. It is
    // nudged back inside rather than clipped, which is what durationsBarLeftX
    // already does for the outermost day bars.
    var plotRightEdge = DURATIONS_VIEW_WIDTH - DURATIONS_MARGIN_RIGHT;

    function bracketOf(userRequestId, group) {
      var left = group.left;
      var right = Math.max(group.right, left + DURATIONS_UR_BRACKET_MIN_WIDTH);
      if (right > plotRightEdge) {
        left -= right - plotRightEdge;
        right = plotRightEdge;
      }
      return {
        userRequestId: userRequestId,
        sampleCount: group.sampleCount,
        firstMs: group.firstMs,
        lastMs: group.lastMs,
        left: left,
        right: right
      };
    }

    var brackets = Object.keys(samplesByUserRequest).map(function (userRequestId) {
      return bracketOf(userRequestId, samplesByUserRequest[userRequestId]);
    });
    var unknownBracket = null;
    if (unknownSamples.length > 0) {
      unknownBracket = bracketOf("", {
        firstMs: Math.min.apply(null, unknownSamples.map(function (sample) { return sample.completionMs; })),
        lastMs: Math.max.apply(null, unknownSamples.map(function (sample) { return sample.completionMs; })),
        left: Math.min.apply(null, unknownSamples.map(function (sample) { return sample.sampleX; })),
        right: Math.max.apply(null, unknownSamples.map(function (sample) { return sample.sampleX; })),
        sampleCount: unknownSamples.length
      });
    }
    return { brackets: brackets, unknownBracket: unknownBracket };
  }

  function durationsUserRequestLaneRowTop(laneRow) {
    return DURATIONS_UR_LANE_TOP + laneRow * DURATIONS_UR_LANE_ROW_PITCH;
  }

  // Stated on the lane's own title line, right-aligned, so the reader learns the
  // count in the same glance as the rows it qualifies. Deliberately worded apart
  // from DURATIONS_UNKNOWN_USER_REQUEST_NAME: a UR that found no row and a sample
  // that has no UR at all are different facts and must not read as one.
  function composeDurationsUserRequestRemainderText(hiddenCount) {
    return "+" + hiddenCount + " UR" + (hiddenCount === 1 ? "" : "s") + " with no free row";
  }

  // Durations are read as magnitudes, so a negative span reads as "−30 min"
  // rather than as a subtraction. Long spans switch to hours so the overflow
  // labels stay short.
  //
  // Rounding happens BEFORE the split into units and before the branch, so the
  // remainder can never carry: rounding 119.5 min per-unit gave "1h 60m", and
  // rounding 59.96 min inside the sub-hour branch gave "60.0 min". Both are the
  // same bug — a value displayed at one precision but split at another.
  function formatDurationMinutes(minutes) {
    var sign = minutes < 0 ? "−" : "";
    var displayedMinutes = Math.round(Math.abs(minutes) * 10) / 10;
    if (displayedMinutes < 60) {
      return sign + displayedMinutes.toFixed(1) + " min";
    }
    var wholeMinutes = Math.round(displayedMinutes);
    return sign + Math.floor(wholeMinutes / 60) + "h " + (wholeMinutes % 60) + "m";
  }

  function formatDurationDayLabel(epochMs) {
    var instant = new Date(epochMs);
    return instant.getUTCDate() + " " + DURATIONS_MONTHS[instant.getUTCMonth()];
  }

  function formatDurationStamp(epochMs) {
    return new Date(epochMs).toISOString().replace("T", " ").slice(0, 16) + " UTC";
  }

  function renderDurationsLongestSpans(samples, listNode, countNode) {
    if (!listNode || !countNode) {
      return;
    }
    var longestSpans = samples.filter(function (sample) {
      return sample.wallMinutes > DURATIONS_CEILING_MINUTES;
    }).sort(function (first, second) {
      if (first.wallMinutes !== second.wallMinutes) {
        return second.wallMinutes - first.wallMinutes;
      }
      if (first.id !== second.id) {
        return first.id < second.id ? -1 : 1;
      }
      return first.completionTime < second.completionTime ? -1 : 1;
    });

    listNode.textContent = "";
    countNode.textContent = longestSpans.length + " span" +
      (longestSpans.length === 1 ? "" : "s") + " over " + DURATIONS_CEILING_MINUTES +
      " minutes in this window; all are listed.";

    longestSpans.forEach(function (sample) {
      var request = durationSampleRequest(sample) || {};
      var item = document.createElement("li");
      item.setAttribute("data-request-id", sample.id);
      item.setAttribute("data-wall-minutes", String(sample.wallMinutes));

      var heading = document.createElement("div");
      heading.className = "durations-longest-spans-item-heading";
      var requestNode = document.createElement("span");
      requestNode.className = "durations-longest-spans-request";
      requestNode.textContent = sample.id;
      var durationNode = document.createElement("span");
      durationNode.className = "durations-longest-spans-duration";
      durationNode.textContent = formatDurationMinutes(sample.wallMinutes);
      heading.appendChild(requestNode);
      heading.appendChild(durationNode);

      var context = document.createElement("div");
      context.className = "durations-longest-spans-item-context";
      var userRequestNode = document.createElement("span");
      userRequestNode.className = "durations-longest-spans-user-request";
      userRequestNode.textContent = durationUserRequestName(durationSampleUserRequestId(sample));
      var routeNode = document.createElement("span");
      routeNode.className = "durations-longest-spans-route";
      routeNode.textContent = durationRouteName(sample.route);
      context.appendChild(userRequestNode);
      context.appendChild(routeNode);

      var titleNode = document.createElement("p");
      titleNode.className = "durations-longest-spans-title";
      titleNode.textContent = request.title || "Untitled REQ";

      item.appendChild(heading);
      item.appendChild(context);
      item.appendChild(titleNode);
      listNode.appendChild(item);
    });
  }

  // Panel B's one annotation: the slowest day's median, stated in figures because
  // a clipped bar cannot state its own value. It lives in its own function so a
  // behavior probe can ask where it lands for a day at ANY x and ANY median —
  // its placement was wrong from the day the view shipped and invisible here
  // purely because the slowest day fell clear of the title's text.
  function drawDurationsSlowestDayAnnotation(svg, slowestDay, dayCentreX) {
    makeDurationsSvgNode(
      svg,
      "text",
      {
        x: dayCentreX.toFixed(1),
        y: DURATIONS_MEDIAN_ANNOTATION_BASELINE_Y,
        class: "durations-mark-label",
        "text-anchor": "middle"
      },
      slowestDay.medianMinutes.toFixed(0) + " min"
    );
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
    var colourLegendNode = document.getElementById("durations-colour-legend");
    var readoutNode = document.getElementById("durations-readout");
    var tableBody = document.getElementById("durations-table-body");
    var longestSpansList = document.getElementById("durations-longest-list");
    var longestSpansCount = document.getElementById("durations-longest-count");
    if (!chartHost || !summaryNode || !tableBody) {
      return;
    }

    var durations = boardData.durations || {};
    var allSamples = durations.samples || [];
    var allDays = durations.days || [];

    chartHost.textContent = "";
    tableBody.textContent = "";
    if (readoutNode) {
      readoutNode.textContent = "";
    }
    if (colourLegendNode) {
      colourLegendNode.textContent = durationsColourLegendText();
      colourLegendNode.setAttribute("aria-label", durationsColourLegendText());
    }

    if (allSamples.length === 0) {
      renderDurationsLongestSpans([], longestSpansList, longestSpansCount);
      summaryNode.textContent =
        "No archived REQ carries both a claim and a completion stamp yet, so there is nothing to measure.";
      return;
    }

    var allSampleTimes = allSamples.map(function (sample) {
      return Date.parse(sample.completionTime);
    });
    var firstCompletionMs = Math.min.apply(null, allSampleTimes);
    var lastCompletionMs = Math.max.apply(null, allSampleTimes);
    var timeEnd = Math.floor(lastCompletionMs / DURATIONS_DAY_MS) * DURATIONS_DAY_MS + DURATIONS_DAY_MS;
    var timeStart = durationsWindow === "all"
      ? Math.floor(firstCompletionMs / DURATIONS_DAY_MS) * DURATIONS_DAY_MS
      : timeEnd - Number(durationsWindow) * DURATIONS_DAY_MS;
    var timeSpan = timeEnd - timeStart;
    var samples = allSamples.filter(function (sample) {
      var completionMs = Date.parse(sample.completionTime);
      return completionMs >= timeStart && completionMs < timeEnd;
    });
    var days = allDays.filter(function (day) {
      var dayMs = Date.parse(day.dayTime);
      return dayMs >= timeStart && dayMs < timeEnd;
    });
    var windowName = durationsWindowName();
    var windowStartLabel = formatDurationDayLabel(timeStart);
    var windowEndLabel = formatDurationDayLabel(timeEnd);
    renderDurationsLongestSpans(samples, longestSpansList, longestSpansCount);

    var excludedSamples = samples.filter(function (sample) {
      return sample.excludedReason;
    });
    summaryNode.textContent =
      windowName +
      " · " +
      windowStartLabel +
      " to " +
      windowEndLabel +
      " UTC (end exclusive)" +
      " · " +
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

    // Every selected domain spans whole UTC days and ends at the midnight AFTER
    // the global latest completion (REQ-248). Day buckets sit at their days'
    // midnights, while Panels B and C draw them at noon inside those slots. The
    // all-history behavior test explicitly selects `all` before comparing this
    // renderer with durations.go's full-history label-planning domain.
    svg.setAttribute(
      "aria-label",
      windowName +
        ", " +
        samples.length +
        " archived REQ" +
        (samples.length === 1 ? "" : "s") +
        ". Three stacked panels sharing a calendar axis from " +
        windowStartLabel +
        " to " +
        windowEndLabel +
        " UTC, end exclusive" +
        ". Panel A plots each archived REQ's duration in minutes coloured by " +
        durationColourChannelName() +
        ", over a lane of brackets grouping its marks by user request. The adjacent ranked list names every span over " +
        DURATIONS_CEILING_MINUTES + " minutes. " +
        durationsColourLegendText() +
        " Panel B plots the median minutes per active day. Panel C counts REQs completed per day. Every value is also listed in the table below."
    );

    function xOfEpoch(epochMs) {
      return DURATIONS_MARGIN_LEFT + ((epochMs - timeStart) / timeSpan) * DURATIONS_PLOT_WIDTH;
    }
    // A day bucket's marks (bar, annotation, hover target) centre on the day's
    // noon — the middle of the [midnight, next midnight) slot the domain gives
    // it. A bucket drawn AT its midnight would straddle the previous day, and
    // the first day's would straddle the axis.
    function durationsDayCentreX(dayEpochMs) {
      return xOfEpoch(dayEpochMs + DURATIONS_DAY_MS / 2);
    }
    // A day-centred bar fits its own slot while barWidth <= dayWidth - 2; past
    // ~280 active days the 4-unit minimum width exceeds the slot, so the
    // outermost bars are nudged back inside the plot instead of bleeding into
    // the gutters.
    function durationsBarLeftX(dayCentreX, barWidth) {
      var barLeft = dayCentreX - barWidth / 2;
      var rightmostLeft = DURATIONS_VIEW_WIDTH - DURATIONS_MARGIN_RIGHT - barWidth;
      return Math.min(Math.max(barLeft, DURATIONS_MARGIN_LEFT), rightmostLeft);
    }
    // Stable rank inside each UTC day replaces completion-time stacking without
    // turning x into a second measure. The inset scales down with very narrow
    // slots; every centre therefore remains in [day midnight, next midnight)
    // even on long boards, while the stated dense target (about eight units per
    // day) keeps useful spread. Identity is the primary order, completion time
    // and original payload position are deterministic tie-breakers.
    function durationsJitteredMarkXs() {
      var samplesByDay = {};
      var markXs = [];
      samples.forEach(function (sample, sampleIndex) {
        var completionMs = Date.parse(sample.completionTime);
        var dayStartMs = Math.floor(completionMs / DURATIONS_DAY_MS) * DURATIONS_DAY_MS;
        var dayKey = String(dayStartMs);
        if (!samplesByDay[dayKey]) {
          samplesByDay[dayKey] = [];
        }
        samplesByDay[dayKey].push({ sample: sample, sampleIndex: sampleIndex, dayStartMs: dayStartMs });
      });
      Object.keys(samplesByDay).forEach(function (dayKey) {
        var rankedSamples = samplesByDay[dayKey].sort(function (first, second) {
          if (first.sample.id !== second.sample.id) {
            return first.sample.id < second.sample.id ? -1 : 1;
          }
          if (first.sample.completionTime !== second.sample.completionTime) {
            return first.sample.completionTime < second.sample.completionTime ? -1 : 1;
          }
          return first.sampleIndex - second.sampleIndex;
        });
        var dayLeft = xOfEpoch(rankedSamples[0].dayStartMs);
        var dayRight = xOfEpoch(rankedSamples[0].dayStartMs + DURATIONS_DAY_MS);
        var dayWidth = dayRight - dayLeft;
        var inset = Math.min(1, dayWidth / 4);
        rankedSamples.forEach(function (rankedSample, rankIndex) {
          markXs[rankedSample.sampleIndex] = rankedSamples.length === 1
            ? dayLeft + dayWidth / 2
            : dayLeft + inset + (rankIndex / (rankedSamples.length - 1)) * (dayWidth - 2 * inset);
        });
      });
      return markXs;
    }
    function yOfMinutes(minutes) {
      var clamped = Math.min(minutes, DURATIONS_CEILING_MINUTES);
      return (
        DURATIONS_MAIN_BOTTOM -
        Math.sqrt(clamped / DURATIONS_CEILING_MINUTES) * (DURATIONS_MAIN_BOTTOM - DURATIONS_MAIN_TOP)
      );
    }
    // R-7 / linear interpolation: position (n-1)*p in the sorted samples and
    // interpolate between its neighbours. Naming the rule makes small even
    // days deterministic rather than leaving quartiles to library convention.
    function durationQuantile(sortedMinutes, probability) {
      var position = (sortedMinutes.length - 1) * probability;
      var lowerIndex = Math.floor(position);
      var upperIndex = Math.ceil(position);
      if (lowerIndex === upperIndex) {
        return sortedMinutes[lowerIndex];
      }
      return sortedMinutes[lowerIndex] +
        (sortedMinutes[upperIndex] - sortedMinutes[lowerIndex]) * (position - lowerIndex);
    }
    function durationsDailyQuantiles() {
      var minutesByDay = {};
      samples.forEach(function (sample) {
        if (sample.excludedReason) {
          return;
        }
        var completionMs = Date.parse(sample.completionTime);
        var dayStartMs = Math.floor(completionMs / DURATIONS_DAY_MS) * DURATIONS_DAY_MS;
        var dayKey = String(dayStartMs);
        if (!minutesByDay[dayKey]) {
          minutesByDay[dayKey] = [];
        }
        minutesByDay[dayKey].push(sample.wallMinutes);
      });
      return Object.keys(minutesByDay).map(function (dayKey) {
        var sortedMinutes = minutesByDay[dayKey].sort(function (first, second) { return first - second; });
        return {
          dayStartMs: Number(dayKey),
          p25: durationQuantile(sortedMinutes, 0.25),
          median: durationQuantile(sortedMinutes, 0.5),
          p75: durationQuantile(sortedMinutes, 0.75)
        };
      }).sort(function (first, second) { return first.dayStartMs - second.dayStartMs; });
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
        {
          x: DURATIONS_MARGIN_LEFT - 8,
          y: y + DURATIONS_TICK_BASELINE_DROP,
          class: "durations-tick",
          "text-anchor": "end"
        },
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
        y: DURATIONS_LANE_MARK_Y + DURATIONS_TICK_BASELINE_DROP,
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
    [0, 5, 15, 30, 45, 60].forEach(function (minutes) {
      gridRow(yOfMinutes(minutes), minutes === 0, String(minutes));
    });

    var markIndex = [];
    var colourContext = durationColourContext(samples);
    var jitteredMarkXs = durationsJitteredMarkXs();
    var dailyQuantiles = durationsDailyQuantiles();
    if (dailyQuantiles.length > 0) {
      var upperRibbonPoints = dailyQuantiles.map(function (day) {
        return durationsDayCentreX(day.dayStartMs).toFixed(1) + " " + yOfMinutes(day.p75).toFixed(1);
      });
      var lowerRibbonPoints = dailyQuantiles.slice().reverse().map(function (day) {
        return durationsDayCentreX(day.dayStartMs).toFixed(1) + " " + yOfMinutes(day.p25).toFixed(1);
      });
      var medianPoints = dailyQuantiles.map(function (day) {
        return durationsDayCentreX(day.dayStartMs).toFixed(1) + " " + yOfMinutes(day.median).toFixed(1);
      });
      makeDurationsSvgNode(svg, "path", {
        d: "M " + upperRibbonPoints.join(" L ") + " L " + lowerRibbonPoints.join(" L ") + " Z",
        fill: "var(--ink-faint)",
        opacity: 0.18,
        class: "durations-quantile-ribbon"
      });
      makeDurationsSvgNode(svg, "path", {
        d: "M " + medianPoints.join(" L "),
        fill: "none",
        stroke: "var(--ink-soft)",
        "stroke-width": 1.5,
        opacity: 0.62,
        class: "durations-quantile-median"
      });
    }
    samples.forEach(function (sample, sampleIndex) {
      var epochMs = Date.parse(sample.completionTime);
      var minutes = sample.wallMinutes;
      var isOverflow = minutes > DURATIONS_CEILING_MINUTES;
      var isReversed = minutes < 0;
      var markX = jitteredMarkXs[sampleIndex];
      var markY = isOverflow
        ? DURATIONS_LANE_MARK_Y
        : isReversed
        ? DURATIONS_BELOW_ZERO_Y
        : yOfMinutes(minutes);
      var markColour = durationMarkColour(sample, colourContext);
      var markAttributes = {
        cx: markX.toFixed(1),
        cy: markY.toFixed(1),
        r: isOverflow || isReversed ? DURATIONS_BAND_MARK_RADIUS : 4,
        fill: isReversed ? "var(--durations-critical)" : markColour.fill,
        class: "durations-mark" + (isReversed ? " durations-mark-critical" : markColour.className ? " " + markColour.className : "")
      };
      if (!isReversed && !markColour.className) {
        markAttributes.opacity = DURATIONS_ORDINARY_MARK_OPACITY;
      }
      makeDurationsSvgNode(svg, "circle", markAttributes);
      markIndex.push({ x: markX, y: markY, sample: sample, epochMs: epochMs, colourLabel: isReversed ? "Reversed stamp" : markColour.label });
    });

    // ---- panel A's UR grouping lane ----
    // A bracket says "these marks above are one user request". No text goes
    // inside the lane; the two gutter ticks name the rows, the title line states
    // the remainder, and the hover names the bracket under the pointer.
    var jitteredMarkXById = {};
    markIndex.forEach(function (mark) {
      jitteredMarkXById[mark.sample.id] = mark.x;
    });
    var userRequestLane = buildDurationsUserRequestBrackets(samples, function (sample) {
      return jitteredMarkXById[sample.id];
    });
    var packedUserRequestLane = packDurationsUserRequestLane(userRequestLane.brackets);
    makeDurationsSvgNode(
      svg,
      "text",
      { x: DURATIONS_MARGIN_LEFT, y: DURATIONS_UR_LANE_TITLE_Y, class: "durations-axis-title" },
      "URs · one bracket per user request, spanning its samples above"
    );
    if (packedUserRequestLane.hiddenCount > 0) {
      makeDurationsSvgNode(
        svg,
        "text",
        {
          x: DURATIONS_VIEW_WIDTH - DURATIONS_MARGIN_RIGHT,
          y: DURATIONS_UR_LANE_TITLE_Y,
          class: "durations-tick",
          "text-anchor": "end"
        },
        composeDurationsUserRequestRemainderText(packedUserRequestLane.hiddenCount)
      );
    }
    function drawDurationsUserRequestBracket(bracket, rowTop, bracketClass) {
      return makeDurationsSvgNode(svg, "rect", {
        x: bracket.left.toFixed(1),
        y: rowTop,
        width: (bracket.right - bracket.left).toFixed(1),
        height: DURATIONS_UR_BRACKET_HEIGHT,
        rx: 2,
        class: bracketClass
      });
    }
    var laneHoverIndex = [];
    packedUserRequestLane.placements.forEach(function (placement) {
      var rowTop = durationsUserRequestLaneRowTop(placement.laneRow);
      // One tone for every row. Alternating two of them was tried and measured
      // first: --ink-faint against --ink-soft is 1.29:1 on the surface the lane
      // actually sits on (Chromium 1194), which is a channel that reads as one
      // colour. The 3-unit gap between rows is what separates them, and it does
      // not need help.
      drawDurationsUserRequestBracket(placement.bracket, rowTop, "durations-ur-bracket");
      laneHoverIndex.push({ bracket: placement.bracket, rowTop: rowTop });
    });
    if (packedUserRequestLane.placements.length > 0) {
      makeDurationsSvgNode(
        svg,
        "text",
        {
          x: DURATIONS_MARGIN_LEFT - 8,
          y: durationsUserRequestLaneRowTop(0) + DURATIONS_UR_BRACKET_HEIGHT - 1,
          class: "durations-tick",
          "text-anchor": "end"
        },
        "URs"
      );
    }
    if (userRequestLane.unknownBracket) {
      drawDurationsUserRequestBracket(
        userRequestLane.unknownBracket,
        DURATIONS_UR_UNKNOWN_ROW_TOP,
        "durations-ur-bracket durations-ur-bracket-unknown"
      );
      laneHoverIndex.push({ bracket: userRequestLane.unknownBracket, rowTop: DURATIONS_UR_UNKNOWN_ROW_TOP });
      makeDurationsSvgNode(
        svg,
        "text",
        {
          x: DURATIONS_MARGIN_LEFT - 8,
          y: DURATIONS_UR_UNKNOWN_ROW_TOP + DURATIONS_UR_BRACKET_HEIGHT - 1,
          class: "durations-tick",
          "text-anchor": "end"
        },
        "no UR"
      );
    }

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

    var dayWidth = (DURATIONS_PLOT_WIDTH * DURATIONS_DAY_MS) / timeSpan;
    var barWidth = Math.max(4, Math.min(24, dayWidth - 2));
    var slowestDay = null;
    days.forEach(function (day) {
      if (!day.hasMedian) {
        return;
      }
      var dayEpochMs = Date.parse(day.dayTime);
      var barLeftX = durationsBarLeftX(durationsDayCentreX(dayEpochMs), barWidth);
      var barTop = yOfDayMedian(day.medianMinutes);
      makeDurationsSvgNode(svg, "rect", {
        x: barLeftX.toFixed(1),
        y: barTop.toFixed(1),
        width: barWidth.toFixed(1),
        height: Math.max(2, DURATIONS_MEDIAN_BOTTOM - barTop).toFixed(1),
        rx: 3,
        class: "durations-bar"
      });
      if (day.medianMinutes > DURATIONS_MEDIAN_CEILING) {
        makeDurationsSvgNode(svg, "rect", {
          x: barLeftX.toFixed(1),
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
      drawDurationsSlowestDayAnnotation(svg, slowestDay, durationsDayCentreX(Date.parse(slowestDay.dayTime)));
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
        y: DURATIONS_COUNT_TOP + DURATIONS_TICK_BASELINE_DROP,
        class: "durations-tick",
        "text-anchor": "end"
      },
      String(peakCount)
    );
    days.forEach(function (day) {
      var dayEpochMs = Date.parse(day.dayTime);
      var columnHeight = (day.completedCount / peakCount) * (DURATIONS_COUNT_BOTTOM - DURATIONS_COUNT_TOP);
      makeDurationsSvgNode(svg, "rect", {
        x: durationsBarLeftX(durationsDayCentreX(dayEpochMs), barWidth).toFixed(1),
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
      windowStartLabel
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
      windowEndLabel
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
      // The lane's own strip, checked before panel A's marks: a bracket carries
      // no text, so the readout is the only place its identity is written out.
      if (pointerY >= DURATIONS_UR_LANE_TOP - 4 && pointerY <= DURATIONS_UR_LANE_BOTTOM + 4) {
        var nearestBracket = null;
        var nearestBracketDistance = Infinity;
        laneHoverIndex.forEach(function (laneEntry) {
          var horizontalGap = Math.max(
            laneEntry.bracket.left - pointerX,
            pointerX - laneEntry.bracket.right,
            0
          );
          var distance = horizontalGap + Math.abs(laneEntry.rowTop + DURATIONS_UR_BRACKET_HEIGHT / 2 - pointerY);
          if (distance < nearestBracketDistance) {
            nearestBracketDistance = distance;
            nearestBracket = laneEntry.bracket;
          }
        });
        if (!nearestBracket) {
          return "";
        }
        return (
          durationUserRequestName(nearestBracket.userRequestId) +
          " · " +
          nearestBracket.sampleCount +
          " sample" +
          (nearestBracket.sampleCount === 1 ? "" : "s") +
          " · " +
          formatDurationDayLabel(nearestBracket.firstMs) +
          " to " +
          formatDurationDayLabel(nearestBracket.lastMs)
        );
      }

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
          durationUserRequestName(durationSampleUserRequestId(sample)) +
          " · " +
          durationColourChannelName() +
          " colour: " +
          nearestMark.colourLabel +
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
        var distance = Math.abs(durationsDayCentreX(Date.parse(day.dayTime)) - pointerX);
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
          durationUserRequestName(durationSampleUserRequestId(sample)),
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
