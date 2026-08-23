  // ---- timeline -----------------------------------------------------------
  //
  // A Gantt of the queue: one row per REQ the visible time window covers, one
  // bar, two segments. The
  // first spans created_at→claimed_at (waiting), the second
  // claimed_at→completed_at (working). A REQ still running draws an open
  // segment to the now-line; a REQ nobody has claimed yet draws an open wait and
  // no work segment at all.
  //
  // Two deliberate departures from the other chart views:
  //
  // 1. NO viewBox. Durations draws into a fixed user-unit space and converts
  //    pointer coordinates by DURATIONS_VIEW_WIDTH / bounds.width — an
  //    assumption a zoom invalidates the moment it exists. This view draws in
  //    CSS pixels, so a pointer's x IS a plot x at every zoom level and there is
  //    no conversion to keep correct.
  //
  // 2. Rows are virtualized. The reporting board carries 560 archived REQs and
  //    the queue only grows; at a readable row height that is thousands of
  //    pixels of extent. Only the rows inside the scrolled window get SVG nodes,
  //    so node count is bounded by the viewport rather than by the archive.
  //
  // The payload decides the numbers. Both spans arrive already measured against
  // one now, signed, with the board's own reversed-stamp verdict attached; this
  // file positions them and never re-measures.

  var TIMELINE_SVG_NS = "http://www.w3.org/2000/svg";
  var TIMELINE_ROW_HEIGHT = 18;
  var TIMELINE_BAR_HEIGHT = 10;
  var TIMELINE_AXIS_HEIGHT = 26;
  // Ticks per axis, and therefore the gap between them: the label format is
  // keyed to that gap, so the two have to be read from one number.
  var TIMELINE_AXIS_TICK_COUNT = 6;
  var TIMELINE_LABEL_WIDTH = 104;
  var TIMELINE_OVERSCAN_ROWS = 4;
  var TIMELINE_MIN_SPAN_MS = 3600000; // one hour in ms — as far in as zoom goes
  var TIMELINE_DAY_MS = 86400000;
  var TIMELINE_YEAR_MS = 365 * TIMELINE_DAY_MS;
  var TIMELINE_ZOOM_STEP = 1.6;
  var TIMELINE_PAN_FRACTION = 0.15;
  var TIMELINE_NOW_JUMP_MARGIN_FRACTION = 0.1;
  // Shortest first, so a window that satisfies more than one level reports the
  // tightest one it actually is.
  var TIMELINE_PERIOD_LEVEL_NAMES = ["day", "week", "month"];
  var TIMELINE_HATCH_PATTERN_ID = "timeline-projected-hatch";
  var TIMELINE_MONTHS = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];

  // The visible time window, held OUTSIDE renderedOnce so switching tabs and
  // coming back preserves the reader's zoom instead of snapping back to fit.
  var timelineViewState = { windowStartMs: 0, windowEndMs: 0, fitted: false };

  // Durations can attach its listeners and forget them: it binds to nodes it
  // rebuilds every render, so the old handlers die with the old DOM. This view
  // binds to the scroll container and to window, which BOTH outlive a render —
  // and a filter change re-renders it. Without an explicit teardown each filter
  // change would leave another live scroll handler behind, and every later
  // scroll would re-render the rows once per stacked handler.
  var timelineListenerTeardowns = [];

  function addTimelineListener(target, eventName, handler, options) {
    target.addEventListener(eventName, handler, options);
    timelineListenerTeardowns.push(function () {
      target.removeEventListener(eventName, handler, options);
    });
  }

  function releaseTimelineListeners() {
    timelineListenerTeardowns.forEach(function (teardown) {
      teardown();
    });
    timelineListenerTeardowns = [];
  }

  // A table rebuild scheduled by the PREVIOUS render holds that render's closure
  // over its own `rows`, and the tbody is shared. It happens to come out right
  // today because the new render schedules its own rebuild and frame callbacks
  // run in order — but nothing enforces that ordering, and a stale frame writing
  // the previous window's rows into the current table is not a bug worth leaving
  // to luck. Cancelled at the same point the listeners are.
  var timelineTableRebuildFrame = null;

  function releaseTimelineTableFrame() {
    if (timelineTableRebuildFrame !== null && window.cancelAnimationFrame) {
      window.cancelAnimationFrame(timelineTableRebuildFrame);
    }
    timelineTableRebuildFrame = null;
  }

  // Each unit is derived from the already-rounded value of the unit below it, so
  // a remainder that rounds up carries instead of overflowing its own field.
  // Rounding per-unit gave "1h 60m" for 119.5 min and "1d 24h" for 2879 min, and
  // rounding inside the sub-hour branch gave "60 min" for 59.96.
  function timelineFormatSpanMinutes(minutes) {
    var sign = minutes < 0 ? "−" : "";
    var wholeMinutes = Math.round(Math.abs(minutes));
    if (wholeMinutes < 60) {
      return sign + wholeMinutes + " min";
    }
    if (wholeMinutes < 60 * 24) {
      return sign + Math.floor(wholeMinutes / 60) + "h " + (wholeMinutes % 60) + "m";
    }
    var wholeHours = Math.round(wholeMinutes / 60);
    return sign + Math.floor(wholeHours / 24) + "d " + (wholeHours % 24) + "h";
  }

  // Every part of a tick's label comes from the tick's own instant. The minute
  // used to be the literal ":00", which made every tick inside one hour read the
  // same — seven ticks, two labels, on the window the Now button lands in.
  //
  // WHICH parts appear is keyed to the gap between ticks rather than to a span
  // threshold of its own, so the label always carries whatever separates one
  // tick from the next: the time once ticks sit less than a day apart, the year
  // once the window is long enough for one day-and-month to come round twice.
  function timelineFormatAxisTick(epochMs, spanMs) {
    var instant = new Date(epochMs);
    var calendarDate = instant.getUTCDate() + " " + TIMELINE_MONTHS[instant.getUTCMonth()];
    if (spanMs / TIMELINE_AXIS_TICK_COUNT < TIMELINE_DAY_MS) {
      return (
        calendarDate +
        " " +
        String(instant.getUTCHours()).padStart(2, "0") +
        ":" +
        String(instant.getUTCMinutes()).padStart(2, "0")
      );
    }
    if (spanMs >= TIMELINE_YEAR_MS) {
      return calendarDate + " " + instant.getUTCFullYear();
    }
    return calendarDate;
  }

  function timelineFormatStamp(epochMs) {
    return new Date(epochMs).toISOString().replace("T", " ").slice(0, 16) + " UTC";
  }

  // Zoom is a pure transform on the visible window, anchored at a fraction of
  // the plot so the instant under the pointer stays under it. Written as its own
  // function because that invariant is the thing worth testing, and testing it
  // through the DOM would test the DOM instead.
  function timelineZoomedWindow(windowStartMs, windowEndMs, zoomFactor, anchorFraction, boundStartMs, boundEndMs) {
    var anchorMs = windowStartMs + (windowEndMs - windowStartMs) * anchorFraction;
    var zoomedSpanMs = (windowEndMs - windowStartMs) / zoomFactor;
    var boundSpanMs = Math.max(boundEndMs - boundStartMs, TIMELINE_MIN_SPAN_MS);
    zoomedSpanMs = Math.min(Math.max(zoomedSpanMs, TIMELINE_MIN_SPAN_MS), boundSpanMs);

    var nextStartMs = anchorMs - zoomedSpanMs * anchorFraction;
    var nextEndMs = nextStartMs + zoomedSpanMs;
    if (nextStartMs < boundStartMs) {
      nextStartMs = boundStartMs;
      nextEndMs = nextStartMs + zoomedSpanMs;
    }
    if (nextEndMs > boundEndMs) {
      nextEndMs = boundEndMs;
      nextStartMs = nextEndMs - zoomedSpanMs;
    }
    return { windowStartMs: nextStartMs, windowEndMs: nextEndMs };
  }

  // Panning slides the window without resizing it, and stops at the bounds — the
  // same clamp the drag path applies, written once so a held arrow key and a long
  // drag cannot stop in different places. The step is a FRACTION of what is on
  // screen rather than a fixed number of milliseconds: a fixed step is
  // imperceptible zoomed out and a jump zoomed in.
  function timelinePannedWindow(windowStartMs, windowEndMs, panFraction, boundStartMs, boundEndMs) {
    var windowSpanMs = windowEndMs - windowStartMs;
    var nextStartMs = Math.min(
      Math.max(windowStartMs + windowSpanMs * panFraction, boundStartMs),
      Math.max(boundEndMs - windowSpanMs, boundStartMs)
    );
    return { windowStartMs: nextStartMs, windowEndMs: nextStartMs + windowSpanMs };
  }

  // The whole keyboard path, as one pure decision: which keys move the window,
  // and where to. It routes zoom through timelineZoomedWindow — the function the
  // wheel and the zoom buttons call — so the keyboard cannot acquire its own
  // floor, ceiling or clamp. Returns null for every key the view does not own,
  // which is what leaves Enter and Space to row activation and Up/Down to
  // scrolling the queue.
  function timelineKeyboardWindow(keyName, windowStartMs, windowEndMs, boundStartMs, boundEndMs) {
    if (keyName === "ArrowLeft" || keyName === "ArrowRight") {
      return timelinePannedWindow(
        windowStartMs,
        windowEndMs,
        keyName === "ArrowRight" ? TIMELINE_PAN_FRACTION : -TIMELINE_PAN_FRACTION,
        boundStartMs,
        boundEndMs
      );
    }
    // "=" and "_" are the unshifted faces of the "+" and "−" keys, so a reader
    // who does not reach for Shift still zooms.
    if (keyName === "+" || keyName === "=") {
      return timelineZoomedWindow(windowStartMs, windowEndMs, TIMELINE_ZOOM_STEP, 0.5, boundStartMs, boundEndMs);
    }
    if (keyName === "-" || keyName === "_") {
      return timelineZoomedWindow(windowStartMs, windowEndMs, 1 / TIMELINE_ZOOM_STEP, 0.5, boundStartMs, boundEndMs);
    }
    return null;
  }

  // Period navigation is the THIRD way to move the window, after the pointer and
  // the keyboard, and it obeys the same rule they do: it computes a CANDIDATE
  // window and hands it to timelineZoomedWindow to settle. Factor 1 at anchor 0
  // means "keep this window, apply the model's floor, ceiling and edge clamp", so
  // the period controls cannot acquire a floor or a clamp of their own.
  //
  // All of it is UTC calendar arithmetic and all of it is O(1) in the board's
  // range: a step is one Date.UTC call, not a walk across the months in between,
  // which is why stepping does not get slower as the archive grows.

  // The start of the calendar period containing an instant.
  function timelinePeriodStart(epochMs, levelName) {
    var instant = new Date(epochMs);
    if (levelName === "month") {
      return Date.UTC(instant.getUTCFullYear(), instant.getUTCMonth(), 1);
    }
    var dayStartMs = Date.UTC(instant.getUTCFullYear(), instant.getUTCMonth(), instant.getUTCDate());
    if (levelName === "week") {
      // getUTCDay is 0 on Sunday; the week runs Monday to Sunday.
      return dayStartMs - ((instant.getUTCDay() + 6) % 7) * TIMELINE_DAY_MS;
    }
    return dayStartMs;
  }

  // The start of the period stepCount periods away. Months step as months rather
  // than as 30 days, so February and March both land on the 1st.
  function timelineSteppedPeriodStart(periodStartMs, levelName, stepCount) {
    if (levelName === "month") {
      var instant = new Date(periodStartMs);
      return Date.UTC(instant.getUTCFullYear(), instant.getUTCMonth() + stepCount, 1);
    }
    return periodStartMs + stepCount * (levelName === "week" ? 7 : 1) * TIMELINE_DAY_MS;
  }

  // One calendar period, positioned by an anchor instant and stepped by
  // stepCount. The step is clamped on the PERIOD rather than on milliseconds:
  // prev and next stop at the first and last period the range reaches, so holding
  // next cannot walk the window off the data and leave a screen with no bars.
  function timelinePeriodWindow(anchorMs, levelName, stepCount, boundStartMs, boundEndMs) {
    var periodStartMs = timelineSteppedPeriodStart(
      timelinePeriodStart(anchorMs, levelName),
      levelName,
      stepCount
    );
    return timelineZoomedWindow(
      periodStartMs,
      timelineSteppedPeriodStart(periodStartMs, levelName, 1),
      1,
      0,
      boundStartMs,
      boundEndMs
    );
  }

  // ---- typed dates ----------------------------------------------------------
  //
  // The FOURTH way to move the window, after the pointer, the keyboard and the
  // period controls, and it obeys the same rule they do: build a candidate and
  // hand it to timelineZoomedWindow to settle. Factor 1 at anchor 0 means "keep
  // this window, apply the model's floor, ceiling and edge clamp", so a typed
  // date cannot acquire a floor or a clamp of its own — the trap the period
  // controls were written to avoid and this inherits for free.
  //
  // Whole UTC days, because the control is a date picker: a start means that
  // day's midnight and an end means the last instant of that day, so typing the
  // same date in both fields gives you that one day rather than an empty window.

  // "YYYY-MM-DD" → the UTC midnight opening that day, or NaN. Deliberately not
  // Date.parse: that accepts far more than a date field emits, and a partial
  // match here would silently move the window somewhere the reader did not type.
  function timelineDateFieldToEpoch(dateText) {
    var match = /^(\d{4})-(\d{2})-(\d{2})$/.exec(String(dateText || "").trim());
    if (!match) {
      return NaN;
    }
    var year = Number(match[1]);
    var month = Number(match[2]);
    var day = Number(match[3]);
    if (month < 1 || month > 12 || day < 1 || day > 31) {
      return NaN;
    }
    var epochMs = Date.UTC(year, month - 1, day);
    // Round-trip guard: Date.UTC rolls 31 February into March rather than
    // rejecting it, and a rolled date is not the one that was typed.
    var rolled = new Date(epochMs);
    if (rolled.getUTCFullYear() !== year || rolled.getUTCMonth() !== month - 1 || rolled.getUTCDate() !== day) {
      return NaN;
    }
    return epochMs;
  }

  // A window START, as its date field shows it.
  function timelineStartEpochToDateField(epochMs) {
    if (isNaN(epochMs)) {
      return "";
    }
    return new Date(epochMs).toISOString().slice(0, 10);
  }

  // A window END, as its date field shows it — and the exact inverse of the
  // parse above, which is the whole reason it is a separate function.
  //
  // The field names the last day the reader wants INCLUDED, while the window's
  // end instant is EXCLUSIVE: a July window ends at 1 August 00:00 and the field
  // must read 31 July. Rendering the end instant's own date instead put every
  // period window a day out — press Month, nudge only the start field, and the
  // end silently moved from 1 August to 2 August, pulling in a REQ the reader
  // never touched and dropping the lit Month chip. Stepping back 1 ms lands
  // inside the last included day whatever the window's shape.
  function timelineEndEpochToDateField(epochMs) {
    if (isNaN(epochMs)) {
      return "";
    }
    return new Date(epochMs - 1).toISOString().slice(0, 10);
  }

  // One or both fields typed, resolved against the window already on screen so a
  // reader can set one end and leave the other alone. Returns null when neither
  // field parses — nothing to apply, and the caller leaves the window untouched
  // rather than jumping somewhere arbitrary.
  //
  // A reversed pair is CLAMPED, not rejected: an end before a start is a partly
  // typed range as often as a mistake, and swapping silently would move the
  // window somewhere the reader did not ask for. The end is pushed to the start's
  // day instead, and the readout then says what the window actually became.
  function timelineTypedWindow(startText, endText, windowStartMs, windowEndMs, boundStartMs, boundEndMs) {
    var typedStartMs = timelineDateFieldToEpoch(startText);
    var typedEndDayMs = timelineDateFieldToEpoch(endText);
    if (isNaN(typedStartMs) && isNaN(typedEndDayMs)) {
      return null;
    }
    var nextStartMs = isNaN(typedStartMs) ? windowStartMs : typedStartMs;
    // The FOLLOWING midnight, because the window's end is exclusive: a field
    // reading 31 July means a window ending at 1 August 00:00, which is exactly
    // what the Month chip produces. That equality is what makes the fields and
    // the period controls describe the same windows rather than windows 1 ms
    // apart — and it is what lets a typed pair light a chip at all.
    var nextEndMs = isNaN(typedEndDayMs) ? windowEndMs : typedEndDayMs + TIMELINE_DAY_MS;
    if (nextEndMs <= nextStartMs) {
      nextEndMs = nextStartMs + TIMELINE_DAY_MS;
    }
    // CLAMP EACH ENDPOINT before settling, because timelineZoomedWindow preserves
    // the span and slides. That is right for a zoom — the reader asked for a
    // width, not a position — and wrong for a typed date, which is a position.
    // Handing it an out-of-range pair made it pin the end to the bound and drag
    // the START back to keep the width: typing 1 July while the end field still
    // read the board's last day moved the start to 30 June, and the field then
    // redrew itself with a date nobody typed. Clamping first means the pair
    // handed over is already inside the range, so the settle below applies the
    // shared floor and edge rules and changes nothing else.
    nextStartMs = Math.min(Math.max(nextStartMs, boundStartMs), boundEndMs);
    nextEndMs = Math.max(Math.min(nextEndMs, boundEndMs), boundStartMs);
    return timelineZoomedWindow(nextStartMs, nextEndMs, 1, 0, boundStartMs, boundEndMs);
  }

  // Which level the window is EXACTLY showing, or null when it is showing a span
  // of the reader's own. A free zoom, a drag, and a period cut short by the end of
  // the range all land here, and the control then states no level instead of
  // claiming one the window no longer has.
  function timelinePeriodLevelOfWindow(windowStartMs, windowEndMs) {
    for (var levelIndex = 0; levelIndex < TIMELINE_PERIOD_LEVEL_NAMES.length; levelIndex++) {
      var levelName = TIMELINE_PERIOD_LEVEL_NAMES[levelIndex];
      if (
        timelinePeriodStart(windowStartMs, levelName) === windowStartMs &&
        timelineSteppedPeriodStart(windowStartMs, levelName, 1) === windowEndMs
      ) {
        return levelName;
      }
    }
    return null;
  }

  // Which level prev and next step by when the window is not exactly a period:
  // the one closest to what is already on screen, so a step after a free zoom is
  // about the size the reader was looking at and still lands on a boundary.
  function timelineNearestPeriodLevel(windowStartMs, windowEndMs) {
    var windowSpanMs = windowEndMs - windowStartMs;
    var nearestLevelName = TIMELINE_PERIOD_LEVEL_NAMES[0];
    var nearestDistanceMs = Infinity;
    TIMELINE_PERIOD_LEVEL_NAMES.forEach(function (levelName) {
      var periodStartMs = timelinePeriodStart(windowStartMs, levelName);
      var periodSpanMs = timelineSteppedPeriodStart(periodStartMs, levelName, 1) - periodStartMs;
      var distanceMs = Math.abs(periodSpanMs - windowSpanMs);
      if (distanceMs < nearestDistanceMs) {
        nearestDistanceMs = distanceMs;
        nearestLevelName = levelName;
      }
    });
    return nearestLevelName;
  }

  // The first row still waiting or still running — where "what is left" sits.
  // Under newest-first order (REQ-318) that is usually near the top rather than
  // hundreds of rows down, so this scan normally returns early; it stays a scan
  // because "newest" and "open" are different questions and a long-open old REQ
  // can sit anywhere in the list.
  function timelineFirstOpenRowIndex(rows) {
    for (var rowIndex = 0; rowIndex < rows.length; rowIndex++) {
      if (rows[rowIndex].waitOpen || rows[rowIndex].workOpen) {
        return rowIndex;
      }
    }
    return -1;
  }

  // Jumping to now is still two movements — recentre the time window, then say
  // where the ROW list goes — but only the first is decided here.
  //
  // The window carries the now-line and the forecast's queue-empty instant
  // together, because "what is left" is the span between them. The row-list
  // movement USED to be decided here too, from the row set the caller was
  // holding. Window-scoped rows (REQ-319) made that the wrong set: the jump
  // changes the window, the window changes which rows exist, and a scrollTop
  // computed before the change indexes into a list that is about to be replaced.
  // So the caller now applies this window, refreshes its rows, and only then
  // asks timelineFirstOpenRowIndex where to scroll — three steps in the one
  // order that can be right.
  function timelineNowJump(nowMs, queueEndMs, boundStartMs, boundEndMs) {
    var earliestMs = isNaN(queueEndMs) ? nowMs : Math.min(nowMs, queueEndMs);
    var latestMs = isNaN(queueEndMs) ? nowMs : Math.max(nowMs, queueEndMs);
    // Margin so both lines sit inside the window rather than on its frame.
    var marginMs = Math.max(
      (latestMs - earliestMs) * TIMELINE_NOW_JUMP_MARGIN_FRACTION,
      TIMELINE_MIN_SPAN_MS / 2
    );
    return timelineZoomedWindow(
      earliestMs - marginMs,
      latestMs + marginMs,
      1,
      0,
      boundStartMs,
      boundEndMs
    );
  }

  // ---- window-scoped rows ---------------------------------------------------
  //
  // A row belongs on screen when the bar it would draw overlaps the visible time
  // window. Before REQ-319 every row was listed at every zoom level and the ones
  // outside the window simply drew nothing — on a 309-REQ board zoomed to one
  // day that is 305 labels above empty space and four bars somewhere inside them.
  //
  // Overlap, not containment: a REQ captured last month and still running today
  // is part of what was happening this week, and hiding it would answer a
  // different question than the reader asked.

  // The SEGMENTS a row draws, as [startMs, endMs] pairs — not one hull over them.
  //
  // The hull was the first shape of this and it was wrong for exactly one row
  // type. A pending REQ draws two DISJOINT marks: an open wait ending at the
  // now-line, and a forecast bar starting after whatever is in flight finishes.
  // A window between the two overlaps their hull and contains no mark at all, so
  // the row was listed with nothing on it — the very thing window-scoping exists
  // to stop. Testing the segments separately is what makes the REQ's GREEN true:
  // every listed row has something drawn on it.
  //
  // Each segment takes a genuine min/max of its own endpoints. Spans are signed,
  // so a reversed stamp can put claimed_at before created_at; the board draws
  // that as a break marker and it still has to be findable in the window it
  // sits in.
  function timelineRowSegments(row, nowMs, projectedRow) {
    var segments = [];
    function addSegment(fromMs, toMs) {
      if (isNaN(fromMs) || isNaN(toMs)) {
        return;
      }
      segments.push({ startMs: Math.min(fromMs, toMs), endMs: Math.max(fromMs, toMs) });
    }
    var createdMs = Date.parse(row.createdTime);
    var claimedMs = row.claimedTime ? Date.parse(row.claimedTime) : NaN;
    var completedMs = row.completedTime ? Date.parse(row.completedTime) : NaN;

    // The wait, drawn created → claimed, or created → now while nobody has
    // claimed it.
    addSegment(createdMs, isNaN(claimedMs) ? nowMs : claimedMs);
    if (row.hasWork) {
      // The work, drawn claimed → completed, or claimed → now while it runs.
      addSegment(claimedMs, isNaN(completedMs) ? nowMs : completedMs);
    }
    if (projectedRow) {
      addSegment(Date.parse(projectedRow.startTime), Date.parse(projectedRow.endTime));
    }
    if (segments.length === 0) {
      // Nothing parsed. The payload guarantees a readable created_at
      // (timeline.go skips rows without one), so this is unreachable today — but
      // the failure has to be "listed everywhere and visibly broken", never
      // "silently deleted from every window including Fit all".
      segments.push({ startMs: -Infinity, endMs: Infinity });
    }
    return segments;
  }

  // The rows with at least one drawn segment overlapping the window, in the
  // order they arrived.
  //
  // Segments are computed ONCE per render and passed in, not derived here: this
  // runs on every window move, a drag moves the window once per frame, and
  // re-parsing three ISO strings per row per frame is the difference between a
  // pointer that keeps up and one that does not.
  function timelineRowsInWindow(rows, rowSegments, windowStartMs, windowEndMs) {
    var inWindow = [];
    for (var rowIndex = 0; rowIndex < rows.length; rowIndex++) {
      var segments = rowSegments[rowIndex];
      for (var segmentIndex = 0; segmentIndex < segments.length; segmentIndex++) {
        var segment = segments[segmentIndex];
        if (segment.startMs <= windowEndMs && segment.endMs >= windowStartMs) {
          inWindow.push(rows[rowIndex]);
          break;
        }
      }
    }
    return inWindow;
  }

  // Which rows have SVG nodes. Everything above and below the scrolled window is
  // absent from the DOM, which is what keeps 560 rows and 5600 rows the same
  // cost. The overscan is what stops a fast scroll showing blank strips.
  function timelineVisibleRowRange(scrollTop, viewportHeight, rowCount) {
    var firstRow = Math.max(0, Math.floor(scrollTop / TIMELINE_ROW_HEIGHT) - TIMELINE_OVERSCAN_ROWS);
    var visibleCount = Math.ceil(viewportHeight / TIMELINE_ROW_HEIGHT) + TIMELINE_OVERSCAN_ROWS * 2;
    return { firstRow: firstRow, lastRow: Math.min(rowCount, firstRow + visibleCount) };
  }

  // A projected bar must never read as a measured one, so it is hatched rather
  // than merely tinted: opacity alone is what a disabled control looks like, and
  // the dashed outline is already spoken for by an OPEN measured bar. The hatch
  // lines carry a class instead of a literal stroke so both themes style them.
  function appendTimelineHatchPattern(svg) {
    var defs = document.createElementNS(TIMELINE_SVG_NS, "defs");
    var pattern = document.createElementNS(TIMELINE_SVG_NS, "pattern");
    pattern.setAttribute("id", TIMELINE_HATCH_PATTERN_ID);
    pattern.setAttribute("width", "6");
    pattern.setAttribute("height", "6");
    pattern.setAttribute("patternUnits", "userSpaceOnUse");
    pattern.setAttribute("patternTransform", "rotate(45)");
    var hatchLine = document.createElementNS(TIMELINE_SVG_NS, "line");
    hatchLine.setAttribute("x1", "0");
    hatchLine.setAttribute("y1", "0");
    hatchLine.setAttribute("x2", "0");
    hatchLine.setAttribute("y2", "6");
    hatchLine.setAttribute("class", "timeline-hatch-line");
    pattern.appendChild(hatchLine);
    defs.appendChild(pattern);
    svg.appendChild(defs);
  }

  // A <g role="button"> takes focus but is not a native button, so Enter and
  // Space never reach the delegated click handler that opens the drawer. Rows
  // advertise the role, so they owe the behavior. Returns the id it activated,
  // or null when the event was not an activation on a detail trigger.
  function timelineKeyboardActivationTarget(keyEvent) {
    if (keyEvent.key !== "Enter" && keyEvent.key !== " " && keyEvent.key !== "Spacebar") {
      return null;
    }
    var trigger = keyEvent.target && keyEvent.target.closest
      ? keyEvent.target.closest("[data-detail-kind]")
      : null;
    if (!trigger) {
      return null;
    }
    return {
      detailKind: trigger.getAttribute("data-detail-kind"),
      detailId: trigger.getAttribute("data-detail-id")
    };
  }

  function makeTimelineSvgNode(parent, name, attributes, textContent) {
    var node = document.createElementNS(TIMELINE_SVG_NS, name);
    Object.keys(attributes).forEach(function (key) {
      node.setAttribute(key, attributes[key]);
    });
    if (textContent !== undefined) {
      node.appendChild(document.createTextNode(textContent));
    }
    parent.appendChild(node);
    return node;
  }

  // The queue-end line and its assumptions, side by side. A forecast that states
  // a date without stating what it assumed is the artifact people screenshot and
  // quote; the assumptions are not a footnote here, they are the other half of
  // the sentence.
  // Emptying both nodes is its own function because the no-rows path has to do
  // it without rendering anything: a forecast left standing beside "no REQ
  // matches" describes rows that are not on screen.
  function clearTimelineForecast() {
    var forecastNode = document.getElementById("timeline-forecast");
    var excludedNode = document.getElementById("timeline-excluded");
    if (forecastNode) {
      forecastNode.textContent = "";
      forecastNode.classList.remove("is-declined");
    }
    if (excludedNode) {
      excludedNode.textContent = "";
    }
    return { forecastNode: forecastNode, excludedNode: excludedNode };
  }

  // showingSubset, not the row list: the rows are a subset and the projection
  // never is, so what this function needs is the one bit the caller can answer
  // and it cannot. It took the filtered rows and ignored them, which is how the
  // view came to say "3 REQs" above a forecast scheduling all 25.
  //
  // The note NAMES NO CAUSE. It used to open "Filters are on", which was true
  // while a filter chip was the only thing that could shrink the row set. Since
  // REQ-319 the visible time window shrinks it too — usually far harder — and
  // both can be on at once, so a sentence naming one of them is wrong about the
  // other and picking between them at render time buys nothing the reader needs.
  // What they need is that the figures below describe more than what is drawn.
  //
  // The note leads the paragraph rather than trailing it, because this sentence
  // is the one people screenshot and quote, and it has to read correctly alone.
  function renderTimelineForecast(projection, showingSubset) {
    var forecastNodes = clearTimelineForecast();
    var forecastNode = forecastNodes.forecastNode;
    var excludedNode = forecastNodes.excludedNode;
    if (!forecastNode || !excludedNode) {
      return;
    }
    var wholeQueueNote = showingSubset
      ? "This covers the whole queue, not the rows shown. "
      : "";

    if (!projection.confident) {
      forecastNode.textContent = wholeQueueNote +
        "No end estimate: " + (projection.declinedReason || "not enough completed work to forecast from") + ".";
      forecastNode.classList.add("is-declined");
    } else {
      forecastNode.classList.remove("is-declined");
      var chainCount = (projection.rows || []).length;
      if (chainCount === 0) {
        forecastNode.textContent = wholeQueueNote + "Nothing left to schedule — every remaining REQ is listed below.";
      } else {
        forecastNode.textContent = wholeQueueNote +
          "Queue empties around " +
          timelineFormatStamp(Date.parse(projection.queueEnd)) +
          " — " +
          chainCount +
          " REQ" +
          (chainCount === 1 ? "" : "s") +
          " run one at a time from " +
          timelineFormatStamp(Date.parse(projection.chainStart)) +
          ". Assumes the median of the last " +
          projection.windowSamples +
          " completed REQs (" +
          projection.normalSamples +
          " substantive at " +
          timelineFormatSpanMinutes(projection.normalMinutes) +
          ", " +
          projection.trivialSamples +
          " mechanical at " +
          timelineFormatSpanMinutes(projection.trivialMinutes) +
          "), one REQ at a time, no parallel builders, and a queue that stops growing. Paused and reversed spans are excluded from both medians." +
          (projection.trivialSamples < projection.minimumSamples ||
          projection.normalSamples < projection.minimumSamples
            ? " A bucket with fewer than " +
              projection.minimumSamples +
              " samples of its own borrows the overall median."
            : "");
      }
    }

    var excluded = projection.excluded || [];
    if (excluded.length === 0) {
      return;
    }
    var excludedHeading = document.createElement("p");
    excludedHeading.className = "timeline-excluded-heading";
    excludedHeading.textContent =
      excluded.length +
      " REQ" +
      (excluded.length === 1 ? "" : "s") +
      (showingSubset ? " from the whole queue" : "") +
      (excluded.length === 1 ? " is" : " are") +
      " not in that estimate, because " +
      (excluded.length === 1 ? "it cannot" : "they cannot") +
      " be given an honest start time:";
    excludedNode.appendChild(excludedHeading);
    var excludedList = document.createElement("ul");
    excludedList.className = "timeline-excluded-list";
    excluded.forEach(function (exclusion) {
      var item = document.createElement("li");
      var trigger = document.createElement("button");
      trigger.type = "button";
      trigger.className = "timeline-excluded-id";
      trigger.setAttribute("data-detail-kind", "request");
      trigger.setAttribute("data-detail-id", exclusion.id);
      trigger.textContent = exclusion.id;
      item.appendChild(trigger);
      item.appendChild(document.createTextNode(" — " + exclusion.reason));
      excludedList.appendChild(item);
    });
    excludedNode.appendChild(excludedList);
  }

  function renderTimelineView() {
    var summaryNode = document.getElementById("timeline-summary");
    var axisHost = document.getElementById("timeline-axis");
    var scrollHost = document.getElementById("timeline-scroll");
    var readoutNode = document.getElementById("timeline-readout");
    var tableBody = document.getElementById("timeline-table-body");
    if (!summaryNode || !axisHost || !scrollHost || !tableBody) {
      return;
    }

    releaseTimelineListeners();
    releaseTimelineTableFrame();

    var timeline = boardData.timeline || {};
    // Filters apply here, unlike Durations. A Gantt filtered to one domain is a
    // straightforward question to ask of a queue; a durations distribution
    // filtered to one domain is a different statistic wearing the same axes.
    //
    // TWO row sets, and the difference is the whole of REQ-319. filterMatchedRows
    // is what the shared filter chips leave, fixed for this render. `rows` is
    // what the visible time window leaves of THAT, and it is re-derived by
    // refreshWindowRows on every window move. Everything that describes what is
    // on screen — the count, the scroll extent, the table, the readout's row
    // index — reads `rows`; everything that describes the population reads
    // filterMatchedRows or the payload.
    var filterMatchedRows = (timeline.rows || []).filter(function (row) {
      return requestMatchesFilters(row.id);
    });
    var rows = filterMatchedRows;
    var nowMs = Date.parse(timeline.now);
    if (isNaN(nowMs)) {
      nowMs = generatedAtMs;
    }

    // The forward half. Keyed by id onto the rows above: a pending REQ already
    // has a row carrying its open wait, and its projected work attaches there
    // rather than becoming a second row for the same REQ.
    var projection = timeline.projection || {};
    var projectedById = {};
    (projection.rows || []).forEach(function (projectedRow) {
      projectedById[projectedRow.id] = projectedRow;
    });

    axisHost.textContent = "";
    scrollHost.textContent = "";
    tableBody.textContent = "";

    // Nothing survives the filters: there is no chart to build and no window to
    // move, so this stays an early return. An empty WINDOW is a different case
    // entirely — the reader zoomed somewhere quiet and needs the axis and the
    // controls to zoom back out — and it is handled inside renderAll, not here.
    if (filterMatchedRows.length === 0) {
      // The forecast describes the rows; with none on screen it must go too.
      clearTimelineForecast();
      summaryNode.textContent = (timeline.rows || []).length
        ? "No REQ matches the current filters."
        : "No REQ carries a readable created_at yet, so there is nothing to place on a timeline.";
      return;
    }

    // Bounds come from the PAYLOAD's range, never from the windowed set: a
    // window that narrowed the rows must not then narrow the bounds it is
    // clamped against, or zooming in would drag the floor up behind it and the
    // reader could never zoom back out. The fallback reads filterMatchedRows for
    // the same reason.
    var boundStartMs = Date.parse(timeline.rangeStart);
    var boundEndMs = Date.parse(timeline.rangeEnd);
    if (isNaN(boundStartMs) || isNaN(boundEndMs) || boundEndMs <= boundStartMs) {
      boundStartMs = Date.parse(filterMatchedRows[0].createdTime);
      boundEndMs = boundStartMs + TIMELINE_MIN_SPAN_MS;
    }
    var queueEndMs = Date.parse(projection.queueEnd);
    if (!isNaN(queueEndMs) && queueEndMs > boundEndMs) {
      boundEndMs = queueEndMs;
    }
    // A little breathing room so a bar that ends exactly at the range edge is
    // not drawn flush against the frame.
    var boundPaddingMs = Math.max((boundEndMs - boundStartMs) * 0.02, 60 * 1000);
    boundStartMs -= boundPaddingMs;
    boundEndMs += boundPaddingMs;

    if (!timelineViewState.fitted || timelineViewState.windowEndMs <= timelineViewState.windowStartMs) {
      timelineViewState.windowStartMs = boundStartMs;
      timelineViewState.windowEndMs = boundEndMs;
      timelineViewState.fitted = true;
    }

    // One segment list per filter-matched row, computed once. refreshWindowRows
    // runs on every window move and compares numbers only.
    var filterMatchedSegments = filterMatchedRows.map(function (row) {
      return timelineRowSegments(row, nowMs, projectedById[row.id]);
    });

    // The reader's place in the list, preserved across a window move. Rows
    // vanishing above the viewport would otherwise slide whatever they were
    // reading off under the axis. Anchored on the id at the top of the viewport:
    // if it survives the move it keeps its screen position, and if it does not,
    // the scroll simply stays where it was and the browser clamps it.
    function topVisibleRowId() {
      if (rows.length === 0) {
        return null;
      }
      var topIndex = Math.min(
        rows.length - 1,
        Math.max(0, Math.floor(scrollHost.scrollTop / TIMELINE_ROW_HEIGHT))
      );
      return rows[topIndex].id;
    }

    function refreshWindowRows() {
      var anchorId = topVisibleRowId();
      var anchorOffset = scrollHost.scrollTop % TIMELINE_ROW_HEIGHT;
      rows = timelineRowsInWindow(
        filterMatchedRows,
        filterMatchedSegments,
        timelineViewState.windowStartMs,
        timelineViewState.windowEndMs
      );
      if (anchorId === null) {
        return;
      }
      for (var rowIndex = 0; rowIndex < rows.length; rowIndex++) {
        if (rows[rowIndex].id === anchorId) {
          // GROW THE EXTENT FIRST. scrollTop is clamped to the scrollable height
          // at the instant it is assigned, and on a widening move the rows SVG
          // is still the narrow window's height here — renderVisibleRows only
          // resizes it afterwards. Writing the anchor first therefore clamped it
          // to the OLD maximum and dropped the reader at the bottom of the
          // window they were leaving: from a month window at scrollTop 400 the
          // Fit-all button landed on 465, the old extent's maximum, when the
          // anchor needed 4900.
          rowsSvg.setAttribute("height", rows.length * TIMELINE_ROW_HEIGHT);
          scrollHost.scrollTop = rowIndex * TIMELINE_ROW_HEIGHT + anchorOffset;
          return;
        }
      }
    }

    function renderSummary() {
      var openCount = rows.filter(function (row) {
        return row.waitOpen || row.workOpen;
      }).length;
      var brokenRowCount = rows.filter(function (row) {
        return row.anomaly || row.waitMinutes < 0 || (row.hasWork && row.workMinutes < 0);
      }).length;
      // The empty window is a state the reader can zoom themselves into, so it
      // says what to do about it rather than reading like a broken board. It
      // must not borrow the "no REQ matches the current filters" wording: the
      // fix for this one is a wider window, not a cleared filter.
      if (rows.length === 0) {
        summaryNode.textContent =
          "Nothing was drawn between " +
          timelineFormatStamp(timelineViewState.windowStartMs) +
          " and " +
          timelineFormatStamp(timelineViewState.windowEndMs) +
          ". Widen the window, step to another period, or press Fit all — " +
          filterMatchedRows.length +
          " REQ" +
          (filterMatchedRows.length === 1 ? " is" : "s are") +
          " outside it.";
        return;
      }
      summaryNode.textContent =
        rows.length +
        " REQ" +
        (rows.length === 1 ? "" : "s") +
        " in the window " +
        timelineFormatStamp(timelineViewState.windowStartMs) +
        " → " +
        timelineFormatStamp(timelineViewState.windowEndMs) +
        ", newest at the top. " +
        openCount +
        " still open, measured to the now-line at " +
        timelineFormatStamp(nowMs) +
        (brokenRowCount ? ". " + brokenRowCount + " with broken stamps, drawn as breaks." : ".") +
        (rows.length < filterMatchedRows.length
          ? " " + (filterMatchedRows.length - rows.length) + " outside the window, not listed."
          : "");
    }

    // The projection is the whole queue's and is never re-derived client-side
    // (see the filter note above); all the forecast needs to know is whether the
    // rows on screen are a subset of it. Under windowing that answer changes as
    // the reader zooms, so it is recomputed — but the forecast's DOM is only
    // rebuilt when the answer actually flips, because rebuilding an exclusion
    // list of N buttons once per drag frame is not free.
    var lastForecastShowedSubset = null;
    function renderForecastIfSubsetChanged() {
      var showingSubset = rows.length < (timeline.rows || []).length;
      if (showingSubset === lastForecastShowedSubset) {
        return;
      }
      lastForecastShowedSubset = showingSubset;
      renderTimelineForecast(projection, showingSubset);
    }

    // ---- table view ----
    // Every value the chart shows is reachable without a pointer, and it lists
    // the same window the chart draws.
    //
    // Built on demand rather than once per render: it is one <tr> per row with
    // no virtualization, and renderAll now runs on every window move — once per
    // frame during a drag. A closed disclosure is the common case and now costs
    // nothing at all; an open one is rebuilt at most once per frame.
    var tableDisclosure = document.querySelector("#view-timeline .timeline-table");

    function renderTimelineTable() {
      timelineTableRebuildFrame = null;
      tableBody.textContent = "";
      rows.forEach(function (row) {
        var request = requestsById[row.id] || {};
        var tableRow = document.createElement("tr");
        [
          row.id,
          request.title || "",
          request.status || "",
          timelineFormatSpanMinutes(row.waitMinutes) + (row.waitOpen ? " (open)" : ""),
          row.hasWork
            ? timelineFormatSpanMinutes(row.workMinutes) + (row.workOpen ? " (open)" : "")
            : "not started",
          row.anomaly ? row.anomalyReason : ""
        ].forEach(function (cellText) {
          var cell = document.createElement("td");
          cell.textContent = cellText;
          tableRow.appendChild(cell);
        });
        tableBody.appendChild(tableRow);
      });
    }

    function markTimelineTableStale() {
      if (!tableDisclosure || !tableDisclosure.open) {
        // Closed: the toggle handler rebuilds it from the current rows on open.
        return;
      }
      if (timelineTableRebuildFrame !== null) {
        return;
      }
      timelineTableRebuildFrame = window.requestAnimationFrame
        ? window.requestAnimationFrame(renderTimelineTable)
        : (renderTimelineTable(), null);
    }

    var axisSvg = makeTimelineSvgNode(axisHost, "svg", {
      class: "timeline-axis-svg",
      height: TIMELINE_AXIS_HEIGHT,
      width: "100%"
    });
    var rowsSvg = makeTimelineSvgNode(scrollHost, "svg", {
      class: "timeline-rows-svg",
      height: rows.length * TIMELINE_ROW_HEIGHT,
      width: "100%",
      role: "img",
      "aria-label":
        "One horizontal bar per REQ drawing anything inside the visible time window, in capture order, newest first. The first segment is the wait from capture to claim, the second is the work from claim to completion. Every value is also listed in the table below."
    });

    function plotWidth() {
      var hostWidth = scrollHost.clientWidth || scrollHost.getBoundingClientRect().width;
      return Math.max(120, hostWidth - TIMELINE_LABEL_WIDTH - 12);
    }

    // 1 SVG unit = 1 CSS pixel here, so this is the whole pointer-to-data story.
    function xOfEpoch(epochMs) {
      var windowSpanMs = timelineViewState.windowEndMs - timelineViewState.windowStartMs || 1;
      return (
        TIMELINE_LABEL_WIDTH +
        ((epochMs - timelineViewState.windowStartMs) / windowSpanMs) * plotWidth()
      );
    }
    function epochOfX(plotX) {
      var windowSpanMs = timelineViewState.windowEndMs - timelineViewState.windowStartMs || 1;
      return (
        timelineViewState.windowStartMs +
        ((plotX - TIMELINE_LABEL_WIDTH) / plotWidth()) * windowSpanMs
      );
    }

    function drawSegment(rowGroup, rowTopY, startMs, endMs, segmentClass) {
      var leftX = xOfEpoch(Math.min(startMs, endMs));
      var rightX = xOfEpoch(Math.max(startMs, endMs));
      var clampedLeft = Math.max(leftX, TIMELINE_LABEL_WIDTH);
      var clampedRight = Math.min(rightX, TIMELINE_LABEL_WIDTH + plotWidth());
      if (clampedRight < TIMELINE_LABEL_WIDTH || clampedLeft > TIMELINE_LABEL_WIDTH + plotWidth()) {
        return;
      }
      makeTimelineSvgNode(rowGroup, "rect", {
        x: clampedLeft.toFixed(1),
        y: (rowTopY + (TIMELINE_ROW_HEIGHT - TIMELINE_BAR_HEIGHT) / 2).toFixed(1),
        width: Math.max(1.5, clampedRight - clampedLeft).toFixed(1),
        height: TIMELINE_BAR_HEIGHT,
        rx: 2,
        class: segmentClass
      });
    }

    function renderVisibleRows() {
      rowsSvg.textContent = "";
      // The scroll extent follows the windowed set, not the filtered one — a
      // window holding four rows must not leave 312 rows of empty scroll below
      // them.
      rowsSvg.setAttribute("height", rows.length * TIMELINE_ROW_HEIGHT);
      appendTimelineHatchPattern(rowsSvg);
      var visible = timelineVisibleRowRange(scrollHost.scrollTop, scrollHost.clientHeight, rows.length);
      for (var rowIndex = visible.firstRow; rowIndex < visible.lastRow; rowIndex++) {
        var row = rows[rowIndex];
        var request = requestsById[row.id] || {};
        var rowTopY = rowIndex * TIMELINE_ROW_HEIGHT;
        var rowGroup = makeTimelineSvgNode(rowsSvg, "g", {
          class: "timeline-row" + (row.anomaly ? " is-broken" : "") + (request.status === "cancelled" ? " is-cancelled" : ""),
          "data-detail-kind": "request",
          "data-detail-id": row.id,
          "data-row-index": String(rowIndex),
          tabindex: "0",
          role: "button",
          "aria-label": timelineRowDescription(row, request)
        });
        makeTimelineSvgNode(rowGroup, "rect", {
          x: 0,
          y: rowTopY,
          width: "100%",
          height: TIMELINE_ROW_HEIGHT,
          class: "timeline-row-hit"
        });
        makeTimelineSvgNode(
          rowGroup,
          "text",
          { x: 6, y: rowTopY + TIMELINE_ROW_HEIGHT - 5, class: "timeline-row-label" },
          row.id
        );

        var createdMs = Date.parse(row.createdTime);
        var claimedMs = row.claimedTime ? Date.parse(row.claimedTime) : nowMs;
        // Same reasoning as the work segment below: a reversed span has no width
        // to draw honestly, so it becomes a break marker at the wait's own start
        // instant rather than a bar drawn left-to-right by drawSegment's
        // min/max sort. An open wait is measured to the now-line and is never
        // reversed, so it keeps its bar.
        if (row.waitMinutes < 0) {
          makeTimelineSvgNode(rowGroup, "rect", {
            x: (xOfEpoch(createdMs) - 3).toFixed(1),
            y: rowTopY + 2,
            width: 6,
            height: TIMELINE_ROW_HEIGHT - 4,
            class: "timeline-segment-broken"
          });
        } else {
          drawSegment(rowGroup, rowTopY, createdMs, claimedMs,
            "timeline-segment timeline-segment-wait" + (row.waitOpen ? " is-open" : ""));
        }

        var projectedRow = projectedById[row.id];
        if (projectedRow) {
          drawProjectedSegment(rowGroup, rowTopY, projectedRow);
        }
        if (row.hasWork) {
          var workStartMs = Date.parse(row.claimedTime);
          var workEndMs = row.completedTime ? Date.parse(row.completedTime) : nowMs;
          // A reversed span has no width to draw honestly, so it is drawn as a
          // break marker at the claim instant rather than as a bar running
          // backwards or clamped forwards to nothing.
          if (row.workMinutes < 0) {
            makeTimelineSvgNode(rowGroup, "rect", {
              x: (xOfEpoch(workStartMs) - 3).toFixed(1),
              y: rowTopY + 2,
              width: 6,
              height: TIMELINE_ROW_HEIGHT - 4,
              class: "timeline-segment-broken"
            });
          } else {
            drawSegment(rowGroup, rowTopY, workStartMs, workEndMs,
              "timeline-segment timeline-segment-work" + (row.workOpen ? " is-open" : ""));
          }
        }
      }
      drawNowRule();
    }

    function drawProjectedSegment(rowGroup, rowTopY, projectedRow) {
      var startMs = Date.parse(projectedRow.startTime);
      var endMs = Date.parse(projectedRow.endTime);
      if (isNaN(startMs) || isNaN(endMs)) {
        return;
      }
      drawSegment(rowGroup, rowTopY, startMs, endMs, "timeline-segment timeline-segment-projected");
    }

    // One node for the whole column, appended after the rows so it reads as an
    // overlay. It lives in the rows SVG rather than in a positioned element
    // because that is what guarantees it uses the same x scale as the bars it
    // crosses — a separate element would need the container's padding folded in
    // by hand, and would drift the first time that padding changed.
    function drawNowRule() {
      if (nowMs < timelineViewState.windowStartMs || nowMs > timelineViewState.windowEndMs) {
        return;
      }
      var nowX = xOfEpoch(nowMs);
      makeTimelineSvgNode(rowsSvg, "line", {
        x1: nowX.toFixed(1),
        y1: 0,
        x2: nowX.toFixed(1),
        y2: rows.length * TIMELINE_ROW_HEIGHT,
        class: "timeline-now-rule"
      });
    }

    function renderAxis() {
      axisSvg.textContent = "";
      var windowSpanMs = timelineViewState.windowEndMs - timelineViewState.windowStartMs || 1;
      var tickCount = TIMELINE_AXIS_TICK_COUNT;
      for (var tickIndex = 0; tickIndex <= tickCount; tickIndex++) {
        var tickMs = timelineViewState.windowStartMs + (windowSpanMs * tickIndex) / tickCount;
        var tickX = xOfEpoch(tickMs);
        makeTimelineSvgNode(axisSvg, "line", {
          x1: tickX.toFixed(1),
          y1: TIMELINE_AXIS_HEIGHT - 6,
          x2: tickX.toFixed(1),
          y2: TIMELINE_AXIS_HEIGHT,
          class: "timeline-axis-tick"
        });
        makeTimelineSvgNode(
          axisSvg,
          "text",
          {
            x: tickX.toFixed(1),
            y: TIMELINE_AXIS_HEIGHT - 10,
            class: "timeline-axis-label",
            "text-anchor": tickIndex === 0 ? "start" : tickIndex === tickCount ? "end" : "middle"
          },
          timelineFormatAxisTick(tickMs, windowSpanMs)
        );
      }
      if (nowMs >= timelineViewState.windowStartMs && nowMs <= timelineViewState.windowEndMs) {
        var nowX = xOfEpoch(nowMs);
        makeTimelineSvgNode(axisSvg, "line", {
          x1: nowX.toFixed(1),
          y1: 0,
          x2: nowX.toFixed(1),
          y2: TIMELINE_AXIS_HEIGHT,
          class: "timeline-now-line"
        });
      }
    }

    // The level a step would use: the one the window is exactly showing, or the
    // nearest one when a free zoom or drag has left it showing a span of its own.
    // Either way prev/next lands on a calendar boundary.
    function steppingLevelName() {
      return (
        timelinePeriodLevelOfWindow(timelineViewState.windowStartMs, timelineViewState.windowEndMs) ||
        timelineNearestPeriodLevel(timelineViewState.windowStartMs, timelineViewState.windowEndMs)
      );
    }

    // The period control has to say when it is no longer exact. After a free zoom
    // or drag the window is a span of the reader's own choosing, and a level
    // button left highlighted would be claiming otherwise. Called from renderAll,
    // so every path that moves the window — buttons, keys, wheel, drag — refreshes
    // it and none of them can leave a stale level on screen.
    function renderPeriodControls() {
      var exactLevelName = timelinePeriodLevelOfWindow(
        timelineViewState.windowStartMs,
        timelineViewState.windowEndMs
      );
      setActiveButton("#view-timeline .timeline-periods", "data-timeline-period", exactLevelName || "");
      var periodStateNode = document.getElementById("timeline-period-state");
      if (periodStateNode) {
        periodStateNode.textContent = exactLevelName ? "one " + exactLevelName : "custom span";
      }
    }

    // The window in text and in the two fields. Called from renderAll, so the
    // pointer, the keyboard, the period chips, Now and Fit all all refresh it and
    // none of them can leave a stale window on screen.
    //
    // The readout carries MINUTES while the fields carry days, and that gap is
    // the point rather than an inconsistency: zoom reaches an hour, the fields
    // cannot express that, and a reader who has zoomed past a day needs to see
    // the window they are actually looking at rather than the nearest date pair.
    var rangeStartField = document.getElementById("timeline-range-start");
    var rangeEndField = document.getElementById("timeline-range-end");
    var rangeReadoutNode = document.getElementById("timeline-range-readout");

    // A field is written back whenever it still holds the value this code last
    // put there — including while it has focus. Skipping every focused field was
    // the obvious rule and the wrong one: a reader who clicks into a field and
    // then zooms the chart leaves that field showing a window the chart is no
    // longer drawn at, and committing it later silently undoes the zoom.
    // Comparing against what we last wrote distinguishes "focused" from "being
    // edited", and only the second is a reason to keep hands off.
    function syncRangeField(field, text) {
      if (!field) {
        return;
      }
      var readerHasEdited =
        document.activeElement === field && field.value !== field.getAttribute("data-synced-value");
      if (readerHasEdited) {
        return;
      }
      field.value = text;
      field.setAttribute("data-synced-value", text);
    }

    function renderRangeControls() {
      if (rangeReadoutNode) {
        rangeReadoutNode.textContent =
          timelineFormatStamp(timelineViewState.windowStartMs) +
          " → " +
          timelineFormatStamp(timelineViewState.windowEndMs);
      }
      syncRangeField(rangeStartField, timelineStartEpochToDateField(timelineViewState.windowStartMs));
      syncRangeField(rangeEndField, timelineEndEpochToDateField(timelineViewState.windowEndMs));
    }

    // Only the field that changed is applied; the other endpoint keeps its exact
    // instant. Re-reading BOTH fields on every commit meant editing one of them
    // also re-applied the other's day-truncated value, so nudging a start could
    // move an end the reader never touched — by a whole day on a period window,
    // and off any sub-day zoom entirely.
    function applyTypedRange(changedField) {
      var startText = changedField === rangeStartField && rangeStartField ? rangeStartField.value : "";
      var endText = changedField === rangeEndField && rangeEndField ? rangeEndField.value : "";
      var typedWindow = timelineTypedWindow(
        startText,
        endText,
        timelineViewState.windowStartMs,
        timelineViewState.windowEndMs,
        boundStartMs,
        boundEndMs
      );
      if (!typedWindow) {
        // Cleared or unparseable. A cleared field is not a request to move, so
        // the window stands — but the field cannot be left blank, or it states a
        // window that does not exist. Restore it unconditionally: the reader
        // emptying it is exactly the case where "leave the focused field alone"
        // would strand it.
        if (changedField) {
          changedField.removeAttribute("data-synced-value");
        }
        renderRangeControls();
        return;
      }
      timelineViewState.windowStartMs = typedWindow.windowStartMs;
      timelineViewState.windowEndMs = typedWindow.windowEndMs;
      renderAll();
    }

    [rangeStartField, rangeEndField].forEach(function (field) {
      if (field) {
        addTimelineListener(field, "change", function () {
          applyTypedRange(field);
        });
      }
    });

    // Every path that moves the window ends here, which is what makes the row
    // set, the counts, the scroll extent and the table describe the same window.
    // The order is load-bearing: rows first, because everything below reads them.
    function renderAll() {
      refreshWindowRows();
      renderSummary();
      renderForecastIfSubsetChanged();
      renderAxis();
      renderVisibleRows();
      renderPeriodControls();
      renderRangeControls();
      markTimelineTableStale();
    }

    function timelineRowDescription(row, request) {
      var parts = [row.id];
      if (request.title) {
        parts.push(request.title);
      }
      parts.push(request.route ? "Route " + request.route : "no route recorded");
      parts.push(request.status || "unknown status");
      parts.push("waited " + timelineFormatSpanMinutes(row.waitMinutes) + (row.waitOpen ? " so far" : ""));
      if (row.hasWork) {
        parts.push("worked " + timelineFormatSpanMinutes(row.workMinutes) + (row.workOpen ? " so far" : ""));
      } else {
        parts.push("not started");
      }
      var projectedRow = projectedById[row.id];
      if (projectedRow) {
        parts.push(
          "projected #" +
            projectedRow.position +
            " in the queue, " +
            timelineFormatSpanMinutes((Date.parse(projectedRow.endTime) - Date.parse(projectedRow.startTime)) / 60000) +
            " (forecast, " +
            projectedRow.bucket +
            ")"
        );
      }
      if (row.anomaly) {
        parts.push("broken stamps: " + row.anomalyReason);
      }
      return parts.join(" · ");
    }

    // ---- interaction ----
    addTimelineListener(scrollHost, "scroll", renderVisibleRows);

    // One keydown listener for the whole chart, registered through the teardown
    // registry like every other listener this view binds to a node that outlives
    // a render. Activation is asked first so Enter and Space keep meaning "open
    // this row" even though the same listener now also moves the window.
    addTimelineListener(scrollHost, "keydown", function (keyEvent) {
      var activation = timelineKeyboardActivationTarget(keyEvent);
      if (activation) {
        keyEvent.preventDefault();
        openDetail(activation.detailKind, activation.detailId);
        return;
      }

      var movedWindow = timelineKeyboardWindow(
        keyEvent.key,
        timelineViewState.windowStartMs,
        timelineViewState.windowEndMs,
        boundStartMs,
        boundEndMs
      );
      if (!movedWindow) {
        return;
      }
      keyEvent.preventDefault();
      timelineViewState.windowStartMs = movedWindow.windowStartMs;
      timelineViewState.windowEndMs = movedWindow.windowEndMs;

      // Rendering rebuilds every row node, so a row that had focus is gone by the
      // time the next key arrives and focus would fall to the body — leaving one
      // arrow press followed by dead keys.
      //
      // This used to assume the row was always still there ("moving the window
      // never moves the vertical scroll, so the same row is still in the
      // virtualized slice"). Window-scoped rows (REQ-319) falsified both halves:
      // a zoom can drop the focused REQ out of the window entirely, and the
      // anchor can move the scroll. So the fallback is the point — focus the
      // chart container, which is what keeps the NEXT key working. Without it,
      // zooming past the focused row killed the keyboard path outright.
      var focusedRow = document.activeElement && document.activeElement.closest
        ? document.activeElement.closest("[data-row-index]")
        : null;
      var focusedRowId = focusedRow ? focusedRow.getAttribute("data-detail-id") : null;
      renderAll();
      if (focusedRowId) {
        var rebuiltRow = rowsSvg.querySelector('[data-detail-id="' + focusedRowId + '"]');
        if (rebuiltRow) {
          rebuiltRow.focus();
        } else {
          scrollHost.focus();
        }
      }
    });

    addTimelineListener(scrollHost, "mousemove", function (moveEvent) {
      if (!readoutNode) {
        return;
      }
      var rowGroup = moveEvent.target.closest ? moveEvent.target.closest("[data-row-index]") : null;
      if (!rowGroup) {
        readoutNode.textContent = "";
        return;
      }
      var row = rows[Number(rowGroup.getAttribute("data-row-index"))];
      readoutNode.textContent = row ? timelineRowDescription(row, requestsById[row.id] || {}) : "";
    });
    addTimelineListener(scrollHost, "mouseleave", function () {
      if (readoutNode) {
        readoutNode.textContent = "";
      }
    });

    // Zoom is modifier-gated so a plain wheel keeps scrolling the rows, which is
    // the motion a 560-row list needs most.
    addTimelineListener(
      scrollHost,
      "wheel",
      function (wheelEvent) {
        if (!wheelEvent.ctrlKey && !wheelEvent.metaKey) {
          return;
        }
        wheelEvent.preventDefault();
        var bounds = scrollHost.getBoundingClientRect();
        var anchorFraction =
          (wheelEvent.clientX - bounds.left - TIMELINE_LABEL_WIDTH) / plotWidth();
        anchorFraction = Math.min(Math.max(anchorFraction, 0), 1);
        var zoomed = timelineZoomedWindow(
          timelineViewState.windowStartMs,
          timelineViewState.windowEndMs,
          wheelEvent.deltaY < 0 ? TIMELINE_ZOOM_STEP : 1 / TIMELINE_ZOOM_STEP,
          anchorFraction,
          boundStartMs,
          boundEndMs
        );
        timelineViewState.windowStartMs = zoomed.windowStartMs;
        timelineViewState.windowEndMs = zoomed.windowEndMs;
        renderAll();
      },
      { passive: false }
    );

    var panState = null;
    addTimelineListener(scrollHost, "pointerdown", function (downEvent) {
      if (downEvent.button !== 0) {
        return;
      }
      panState = { pointerX: downEvent.clientX, windowStartMs: timelineViewState.windowStartMs };
      scrollHost.classList.add("is-panning");
    });
    addTimelineListener(scrollHost, "pointermove", function (moveEvent) {
      if (!panState) {
        return;
      }
      var windowSpanMs = timelineViewState.windowEndMs - timelineViewState.windowStartMs;
      var shiftMs = ((panState.pointerX - moveEvent.clientX) / plotWidth()) * windowSpanMs;
      var nextStartMs = Math.min(
        Math.max(panState.windowStartMs + shiftMs, boundStartMs),
        boundEndMs - windowSpanMs
      );
      timelineViewState.windowStartMs = nextStartMs;
      timelineViewState.windowEndMs = nextStartMs + windowSpanMs;
      renderAll();
    });
    ["pointerup", "pointercancel", "pointerleave"].forEach(function (eventName) {
      addTimelineListener(scrollHost, eventName, function () {
        panState = null;
        scrollHost.classList.remove("is-panning");
      });
    });

    function wireToolbarButton(buttonId, apply) {
      var button = document.getElementById(buttonId);
      if (button) {
        button.onclick = function () {
          apply();
          renderAll();
        };
      }
    }
    wireToolbarButton("timeline-zoom-in", function () {
      var zoomed = timelineZoomedWindow(
        timelineViewState.windowStartMs, timelineViewState.windowEndMs,
        TIMELINE_ZOOM_STEP, 0.5, boundStartMs, boundEndMs);
      timelineViewState.windowStartMs = zoomed.windowStartMs;
      timelineViewState.windowEndMs = zoomed.windowEndMs;
    });
    wireToolbarButton("timeline-zoom-out", function () {
      var zoomed = timelineZoomedWindow(
        timelineViewState.windowStartMs, timelineViewState.windowEndMs,
        1 / TIMELINE_ZOOM_STEP, 0.5, boundStartMs, boundEndMs);
      timelineViewState.windowStartMs = zoomed.windowStartMs;
      timelineViewState.windowEndMs = zoomed.windowEndMs;
    });
    wireToolbarButton("timeline-zoom-fit", function () {
      timelineViewState.windowStartMs = boundStartMs;
      timelineViewState.windowEndMs = boundEndMs;
    });
    // Zooming anchors at the centre, so the forecast at the far right takes a
    // long drag to reach once you have zoomed in far enough to read it. A
    // forecast you have to hunt for does not answer "when does the queue empty".
    // Three steps in the one order that can be right: the window first, then the
    // rows that window admits, and only then where to scroll among them.
    // Deciding the scroll from the pre-jump rows would index into a list the
    // jump is about to replace — see timelineNowJump's own note.
    var nowButton = document.getElementById("timeline-zoom-now");
    if (nowButton) {
      nowButton.onclick = function () {
        var nowWindow = timelineNowJump(
          nowMs, Date.parse(projection.queueEnd), boundStartMs, boundEndMs);
        timelineViewState.windowStartMs = nowWindow.windowStartMs;
        timelineViewState.windowEndMs = nowWindow.windowEndMs;
        renderAll();
        var openRowIndex = timelineFirstOpenRowIndex(rows);
        if (openRowIndex >= 0) {
          scrollHost.scrollTop = openRowIndex * TIMELINE_ROW_HEIGHT;
          renderVisibleRows();
        }
      };
    }

    // The period controls: three levels and a step either way, all of them going
    // through timelinePeriodWindow and therefore through timelineZoomedWindow.
    // The window's own midpoint is the anchor, so choosing a level keeps the
    // reader near what they were looking at and a step moves off the period they
    // are on.
    function applyPeriodWindow(levelName, stepCount) {
      var periodWindow = timelinePeriodWindow(
        (timelineViewState.windowStartMs + timelineViewState.windowEndMs) / 2,
        levelName,
        stepCount,
        boundStartMs,
        boundEndMs
      );
      timelineViewState.windowStartMs = periodWindow.windowStartMs;
      timelineViewState.windowEndMs = periodWindow.windowEndMs;
    }
    document.querySelectorAll("#view-timeline [data-timeline-period]").forEach(function (levelButton) {
      levelButton.onclick = function () {
        applyPeriodWindow(levelButton.getAttribute("data-timeline-period"), 0);
        renderAll();
      };
    });
    wireToolbarButton("timeline-period-prev", function () {
      applyPeriodWindow(steppingLevelName(), -1);
    });
    wireToolbarButton("timeline-period-next", function () {
      applyPeriodWindow(steppingLevelName(), 1);
    });

    addTimelineListener(window, "resize", renderAll);

    if (tableDisclosure) {
      addTimelineListener(tableDisclosure, "toggle", function () {
        if (tableDisclosure.open) {
          renderTimelineTable();
        }
      });
    }
    renderAll();
    // renderAll's mark only SCHEDULES a rebuild, a frame away. A disclosure left
    // open across a filter change would show the previous render's rows for that
    // frame, so the first paint is taken synchronously here and the scheduled
    // frame then finds nothing to do.
    if (tableDisclosure && tableDisclosure.open) {
      releaseTimelineTableFrame();
      renderTimelineTable();
    }
  }
