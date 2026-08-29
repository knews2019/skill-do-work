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
  var TIMELINE_GROUP_HEADER_HEIGHT = 34;
  var TIMELINE_BAR_HEIGHT = 10;
  var TIMELINE_AXIS_HEIGHT = 26;
  // The TARGET number of gaps between axis ticks, not an exact count.
  //
  // It used to be exact: the axis divided the window into this many equal parts.
  // That put ticks at arbitrary instants on any window that is not a multiple of
  // six of something — 28 hours apart across a week, 5.17 days across a month —
  // and the label formatter then printed a bare DATE for each one, claiming
  // midnights the ticks did not sit at and skipping a calendar day outright.
  //
  // Ticks now sit on calendar boundaries and this picks which boundary: the entry
  // of TIMELINE_AXIS_TICK_STEPS whose gap divides the window into closest to this
  // many parts. The real count therefore varies — 4 for a month of Mondays, 8 for
  // a week of midnights — and that is the point.
  var TIMELINE_AXIS_TICK_COUNT = 6;

  // The gaps an axis is allowed to put between ticks, shortest first, each with
  // the boundary it aligns to. Sub-day and whole-day gaps divide a UTC day, so
  // they align to midnight for free; week gaps align to Monday through
  // TIMELINE_WEEK_ALIGNMENT_MS below; month gaps align to the 1st.
  //
  // The ladder is CLOSED and does not go stale, because it enumerates a
  // mathematical range rather than a set of cases: it spans from the view's zoom
  // floor (one hour) to any archive age, and the selection rule is "closest to
  // TIMELINE_AXIS_TICK_COUNT gaps". A longer archive picks a longer rung; it never
  // needs a new one.
  // The epoch's first Monday, so a week-long gap lands on Mondays.
  var TIMELINE_WEEK_ALIGNMENT_MS = Date.UTC(1970, 0, 5);
  var TIMELINE_AXIS_TICK_STEPS = [
    { ms: 5 * 60000 },
    { ms: 10 * 60000 },
    { ms: 15 * 60000 },
    { ms: 30 * 60000 },
    { ms: 3600000 },
    { ms: 2 * 3600000 },
    { ms: 3 * 3600000 },
    { ms: 4 * 3600000 },
    { ms: 6 * 3600000 },
    { ms: 12 * 3600000 },
    { ms: 86400000 },
    { ms: 2 * 86400000 },
    { ms: 7 * 86400000, alignMs: TIMELINE_WEEK_ALIGNMENT_MS },
    { ms: 14 * 86400000, alignMs: TIMELINE_WEEK_ALIGNMENT_MS },
    { months: 1 },
    { months: 2 },
    { months: 3 },
    { months: 6 },
    { months: 12 },
    { months: 24 },
    { months: 60 },
    { months: 120 }
  ];
  // The most ticks an axis will draw, whatever the window. A backstop against a
  // pathological window rather than a design limit: every rung of the ladder above
  // yields between three and nine gaps for the spans it wins.
  var TIMELINE_AXIS_TICK_LIMIT = 64;
  // Wide enough for the id plus a title worth reading, and no wider: every pixel
  // here comes out of the plot. At the 10px monospace label face measured on
  // Chromium 1194 (6.0219 px/cell) that is 28 cells — "REQ-042" plus a two-space
  // separator leaves 19 for the title, 18 plus an ellipsis once it is cut.
  // Picked from the render at 1400px and 760px (REQ-322), and PINNED: the label
  // probe reads this constant and fails below 20 cells, so it cannot drift back
  // down without someone deciding to.
  var TIMELINE_LABEL_WIDTH = 184;
  var TIMELINE_OVERSCAN_ROWS = 4;
  var TIMELINE_UNKNOWN_USER_REQUEST_NAME = "No UR recorded";
  var TIMELINE_MIN_SPAN_MS = 3600000; // one hour in ms — as far in as zoom goes
  var TIMELINE_DAY_MS = 86400000;
  var TIMELINE_ZOOM_STEP = 1.6;
  var TIMELINE_PAN_FRACTION = 0.15;
  var TIMELINE_NOW_JUMP_MARGIN_FRACTION = 0.1;
  // The narrowest window the Now button may land on.
  //
  // The margin used to floor on TIMELINE_MIN_SPAN_MS / 2, which put Now on a
  // one-hour window whenever the now-line and the forecast's queue-empty instant
  // were close together — the ordinary case on a queue that is nearly drained, and
  // on this repo's board the state Now landed in every time. One hour IS the zoom
  // floor, so the obvious next move was dead: the + button, ctrl+wheel and the +
  // key were all silent no-ops. Half a day gives the now-line a day's work either
  // side of it to be read against, and leaves zoom somewhere to go.
  var TIMELINE_NOW_JUMP_MINIMUM_SPAN_MS = 43200000; // twelve hours in ms
  var TIMELINE_HATCH_PATTERN_ID = "timeline-projected-hatch";
  // The narrowest a drawn segment may be. 1.5px was technically present and
  // practically invisible: at Fit all over three months every completed REQ was
  // a sliver in which the wait and the work occupied the same pixel. Picked from
  // the render (REQ-323).
  var TIMELINE_MIN_SEGMENT_WIDTH = 3;
  // Below this total width a bar cannot show two segments honestly — 3px each
  // plus a boundary — so the row draws ONE marker instead of two slivers
  // claiming a split the pixels cannot carry.
  var TIMELINE_MIN_SPLIT_WIDTH = 7;
  // How far a pointer must travel before a press becomes a pan. Under it the
  // press is a click and the drawer opens on release.
  //
  // The first pointermove used to engage the pan and RE-RENDER, and the click the
  // engine then synthesized had no surviving [data-detail-kind] ancestor to find.
  // The window did not have to move for that — at Fit all it clamps to the window
  // it started in and the click was lost anyway — so the threshold gates the
  // render, not the shift.
  //
  // That is ONE of two ways the chart has eaten a click. The other was pointer
  // capture taken on pointerdown, which retargets the synthesized click whether or
  // not anything re-renders; the threshold cannot protect against it, which is why
  // the capture now waits for the engage too (see capturePanPointer below).
  //
  // 4px: an ordinary hand tremor on a trackpad stays under it, an intended drag
  // clears it in the first frame.
  var TIMELINE_PAN_THRESHOLD_PX = 4;
  var TIMELINE_MONTHS = ["Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"];

  // The visible time window, held OUTSIDE renderedOnce so switching tabs and
  // coming back preserves the reader's zoom instead of snapping back to fit.
  // rovingRowIndex is the row list's single Tab stop. Held here rather than derived
  // because it is the one thing about the list that cannot be recomputed: which row the
  // reader was last on. Every row used to carry tabindex="0", so escaping the chart cost
  // one Tab press per row — 29 on the board that reported it.
  var timelineViewState = { windowStartMs: 0, windowEndMs: 0, fitted: false, rovingRowIndex: 0 };

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

  // At most one whole-view render per frame during a drag: a trackpad delivers
  // pointermoves faster than the compositor draws, and the intermediate windows
  // cannot all be shown.
  //
  // A frame scheduled by the PREVIOUS render holds that render's closure, so a
  // filter change that re-enters this module would let it re-render the old rows
  // one frame later. Same hazard as the table frame above, cancelled at the same
  // point — which is why the handle lives out here rather than inside the render,
  // where the listeners that schedule it do.
  var timelineFrameRender = null;

  function releaseTimelineFrameRender() {
    if (timelineFrameRender !== null && window.cancelAnimationFrame) {
      window.cancelAnimationFrame(timelineFrameRender);
    }
    timelineFrameRender = null;
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
  // WHICH parts appear is keyed to the GAP the axis actually chose, so the label
  // always carries whatever separates one tick from the next.
  //
  // The gap used to be derived here as spanMs / TIMELINE_AXIS_TICK_COUNT, back
  // when the axis divided the window into that many equal parts. It no longer
  // does — the ticks sit on calendar boundaries — and a derived gap would be a
  // second, disagreeing opinion about the spacing. So the gap is passed in, from
  // the same step that positioned the ticks.
  //
  // A DATE-ONLY LABEL IS A CLAIM OF MIDNIGHT, and that is why this pairs with
  // calendar-aligned ticks rather than merely reading nicer. Six equal parts
  // across a week put ticks 28 hours apart and this branch printed a bare date
  // for every one of them: the week of 6 July read "6, 7, 8, 9, 10, 11, 13 Jul"
  // — 12 July missing, and the tick labelled "9 Jul" sitting at 9 Jul 12:00.
  //
  // The year is keyed on whether the window SPANS more than one calendar year,
  // not on whether it is longer than 365 days. A window from December into
  // January needs the year on both ends and is nine days long.
  function timelineFormatAxisTick(epochMs, gapMs, windowStartMs, windowEndMs) {
    var instant = new Date(epochMs);
    var calendarDate = instant.getUTCDate() + " " + TIMELINE_MONTHS[instant.getUTCMonth()];
    if (gapMs < TIMELINE_DAY_MS) {
      return (
        calendarDate +
        " " +
        String(instant.getUTCHours()).padStart(2, "0") +
        ":" +
        String(instant.getUTCMinutes()).padStart(2, "0")
      );
    }
    if (new Date(windowStartMs).getUTCFullYear() !== new Date(windowEndMs).getUTCFullYear()) {
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

  // The whole WINDOW-MOVING keyboard path, as one pure decision: which keys move the
  // window, and where to. It routes zoom through timelineZoomedWindow — the function the
  // wheel and the zoom buttons call — so the keyboard cannot acquire its own floor,
  // ceiling or clamp. Returns null for every key the window does not own, which is what
  // leaves Enter and Space to row activation and Up/Down to the row list's roving stop
  // (and, with focus on the chart rather than a row, to the browser's own scrolling).
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

  // Window navigation is the THIRD way to move the window, after the pointer and
  // the keyboard, and it obeys the same rule they do: it computes a CANDIDATE
  // window and hands it to timelineZoomedWindow to settle. Factor 1 at anchor 0
  // means "keep this window, apply the model's floor, ceiling and edge clamp", so
  // the toolbar cannot acquire a floor or a clamp of its own.

  // The UTC midnight opening the day an instant falls in. All that survives of the
  // calendar arithmetic the toolbar used to run on, because the date FIELDS are
  // still whole UTC days even though the chips beside them no longer are.
  function timelineUtcDayStart(epochMs) {
    var instant = new Date(epochMs);
    return Date.UTC(instant.getUTCFullYear(), instant.getUTCMonth(), instant.getUTCDate());
  }

  // One trailing window, ending at NOW and reaching back the span the chip names.
  // "all" — and anything else that is not a positive number of days — means the
  // whole recorded range, which is the window the view opens on.
  //
  // CLAMP EACH ENDPOINT BEFORE SETTLING, because a trailing window is a POSITION
  // and timelineZoomedWindow preserves a WIDTH. Handed an out-of-range candidate
  // it pins the offending edge and drags the other one along to keep the span: on
  // a three-day archive, Last 90 days would pin the start to the range start and
  // then push the END ninety days forward, past the forecast and off now — so the
  // one thing every chip on this toolbar promises would be false on exactly the
  // boards most likely to press it. Clamped first, an over-long window is CUT
  // SHORT at the range start and still ends at now; the state readout then says
  // "part of" rather than claiming a span the board does not have. This is the
  // same distinction timelineTypedWindow documents at its own clamp.
  function timelineTrailingWindow(trailingWindowValue, nowMs, boundStartMs, boundEndMs) {
    var trailingDayCount = Number(trailingWindowValue);
    var candidateStartMs = boundStartMs;
    var candidateEndMs = boundEndMs;
    if (isFinite(trailingDayCount) && trailingDayCount > 0) {
      candidateStartMs = nowMs - trailingDayCount * TIMELINE_DAY_MS;
      candidateEndMs = nowMs;
    }
    candidateStartMs = Math.min(Math.max(candidateStartMs, boundStartMs), boundEndMs);
    candidateEndMs = Math.max(Math.min(candidateEndMs, boundEndMs), boundStartMs);
    return timelineZoomedWindow(candidateStartMs, candidateEndMs, 1, 0, boundStartMs, boundEndMs);
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
  // whole-day window a day out — set a window, nudge only the start field, and
  // the end silently moved from 1 August to 2 August, pulling in a REQ the reader
  // never touched. Stepping back 1 ms lands inside the last included day whatever
  // the window's shape.
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
    // A DATE OUTSIDE THE RANGE BECOMES THE NEAREST DAY THE BOARD HAS, clamped
    // here at DAY granularity before any of the arithmetic below runs.
    //
    // Clamping only at the endpoint stage was not enough. Typing 2026-09-30 into
    // From, on a board whose range ends 2026-08-25 04:23, put the typed start on
    // boundEndMs; the end field was untouched so it stayed on boundEndMs too; the
    // reversed-pair guard pushed the end a day out, the endpoint clamp pulled it
    // straight back, and the settle then raised the zero span to the zoom floor
    // and slid it — an empty one-hour window tucked behind the right edge, with
    // the field still showing the rejected date. A date past the end means "the
    // end", and the last day the board has is what that is.
    var firstDayMs = timelineUtcDayStart(boundStartMs);
    var lastDayMs = timelineUtcDayStart(boundEndMs - 1);
    if (!isNaN(typedStartMs)) {
      typedStartMs = Math.min(Math.max(typedStartMs, firstDayMs), lastDayMs);
    }
    if (!isNaN(typedEndDayMs)) {
      typedEndDayMs = Math.min(Math.max(typedEndDayMs, firstDayMs), lastDayMs);
    }
    var nextStartMs = isNaN(typedStartMs) ? windowStartMs : typedStartMs;
    // The FOLLOWING midnight, because the window's end is exclusive: a field
    // reading 31 July means a window ending at 1 August 00:00 rather than at
    // 23:59:59.999. That is what makes rendering a window into the fields and
    // parsing them back inverses, so editing one field cannot mangle the other.
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
    // Margin so both lines sit inside the window rather than on its frame, floored
    // so the landing window is a window rather than the zoom floor. See
    // TIMELINE_NOW_JUMP_MINIMUM_SPAN_MS.
    var marginMs = Math.max(
      (latestMs - earliestMs) * TIMELINE_NOW_JUMP_MARGIN_FRACTION,
      TIMELINE_NOW_JUMP_MINIMUM_SPAN_MS / 2
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

  // ---- where a row's two spans end -----------------------------------------
  //
  // ONE ANSWER, read by the segment model and by the renderer, because they have
  // twice drifted into disagreeing about which rows are drawable and the summary
  // then described a chart nobody was looking at. NaN means "there is no honest
  // end", and every reader turns that into the break marker rather than a bar.
  //
  // A span runs to the now-line ONLY when the payload says it is open, and the
  // payload says that only for a REQ that has not stopped (timeline.go). Falling
  // back to now whenever a stamp was missing is what drew 25 finished REQs as
  // work in flight, one of them 24.8 days long.

  // Where the wait ends: at the claim, at the now-line while the REQ is still
  // waiting, at the instant it stopped when it stopped without ever being
  // claimed, and nowhere at all when it stopped and that instant is unknown.
  function timelineWaitEndMs(row, nowMs, claimedMs, completedMs) {
    if (!isNaN(claimedMs)) {
      return claimedMs;
    }
    if (row.waitOpen) {
      return nowMs;
    }
    return completedMs;
  }

  // Where the work ends: at the completion, at the now-line while it is running,
  // and nowhere at all when it stopped and that instant is unknown.
  function timelineWorkEndMs(row, nowMs, completedMs) {
    if (!isNaN(completedMs)) {
      return completedMs;
    }
    return row.workOpen ? nowMs : NaN;
  }

  // Whether this row draws a BREAK rather than a bar — the one verdict the segment
  // model, the renderer and the summary all read. Two causes, one mark: a reversed
  // span (an end before its start), and a REQ that stopped with no resolvable end
  // instant. What the reader needs from either is "the bookkeeping for this row is
  // broken", so they get the same marker and are counted together.
  function timelineRowDrawsABreak(row, nowMs) {
    var claimedMs = row.claimedTime ? Date.parse(row.claimedTime) : NaN;
    var completedMs = row.completedTime ? Date.parse(row.completedTime) : NaN;
    if (row.waitMinutes < 0 || isNaN(timelineWaitEndMs(row, nowMs, claimedMs, completedMs))) {
      return true;
    }
    return !!row.hasWork &&
      (row.workMinutes < 0 || isNaN(timelineWorkEndMs(row, nowMs, completedMs)));
  }

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
  // Every segment is where the renderer DRAWS something, at the size it draws it.
  // A forward span is its own min/max interval. A reversed one — claimed_at
  // before created_at — is a POINT, because a break marker is all the renderer
  // puts on screen for it; taking the min/max there would claim a bar across an
  // interval nothing is drawn across, and the row would be listed in windows
  // showing none of it. The point sits exactly where the marker does, so a broken
  // row stays findable in the window it actually appears in.
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
    //
    // A REVERSED wait is not drawn as a bar at all: renderVisibleRows puts a 6px
    // break marker at the CREATED instant and nothing between. Modelling it as
    // the min/max interval is the forecast-gap defect in another costume — a
    // window inside the interval overlaps a hull containing no mark. Keyed on
    // row.waitMinutes, the same field the renderer branches on, so the two cannot
    // drift into disagreeing about which rows are reversed.
    // A wait with no honest end is the same shape as a reversed one: a point at
    // the created instant, where the renderer puts the break marker.
    var waitEndMs = timelineWaitEndMs(row, nowMs, claimedMs, completedMs);
    if (row.waitMinutes < 0 || isNaN(waitEndMs)) {
      addSegment(createdMs, createdMs);
    } else {
      addSegment(createdMs, waitEndMs);
    }
    if (row.hasWork) {
      var workEndMs = timelineWorkEndMs(row, nowMs, completedMs);
      if (row.workMinutes < 0 || isNaN(workEndMs)) {
        // Same again, at the instant renderVisibleRows puts this row's marker.
        // A REQ that stopped with no resolvable completion instant lands here:
        // there is no width to claim, and the marker says so.
        addSegment(claimedMs, claimedMs);
      } else {
        // The work, drawn claimed → completed, or claimed → now while it runs.
        addSegment(claimedMs, workEndMs);
      }
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

  // Whether a row is too narrow to draw a wait/work split, and if so the extent
  // the single marker covers. Pure so a probe can drive it: the render passes
  // the plot width it already measured rather than this reaching for the DOM.
  //
  // At Fit all over three months a completed REQ was two adjacent slivers of
  // different hues occupying the same pixel — a split the pixels cannot carry.
  // Below TIMELINE_MIN_SPLIT_WIDTH there is no room for two floored segments and
  // a boundary between them, so the row draws one marker for the whole span.
  function timelineCollapsedRowMark(segments, windowStartMs, windowEndMs, plotWidthPx) {
    // Already one drawn segment: there is no split to withdraw. This is also what
    // spares the unparseable row, whose -Infinity → Infinity sentinel segment
    // timelineRowSegments only ever emits ALONE — collapsing that would draw one
    // marker across the whole chart.
    if (segments.length < 2) {
      return null;
    }
    var earliestMs = Infinity;
    var latestMs = -Infinity;
    for (var segmentIndex = 0; segmentIndex < segments.length; segmentIndex++) {
      earliestMs = Math.min(earliestMs, segments[segmentIndex].startMs);
      latestMs = Math.max(latestMs, segments[segmentIndex].endMs);
    }
    var windowSpanMs = windowEndMs - windowStartMs || 1;
    var clampedStartMs = Math.max(earliestMs, windowStartMs);
    var clampedEndMs = Math.min(latestMs, windowEndMs);
    var spanWidthPx = ((clampedEndMs - clampedStartMs) / windowSpanMs) * plotWidthPx;
    if (spanWidthPx >= TIMELINE_MIN_SPLIT_WIDTH) {
      return null;
    }
    return { startMs: earliestMs, endMs: latestMs };
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
        // ASYMMETRIC ON PURPOSE. windowEndMs is the next window's first instant
        // everywhere else here — timelineEndEpochToDateField renders windowEndMs
        // minus 1 to say so — so a segment beginning exactly there belongs to the
        // next window, and drawSegment would put its floored rectangle at the
        // clipped right edge. windowStartMs IS in the window, and a span ending on
        // it draws a mark at the left edge the reader can see, so that side stays
        // inclusive.
        if (segment.startMs < windowEndMs && segment.endMs >= windowStartMs) {
          inWindow.push(rows[rowIndex]);
          break;
        }
      }
    }
    return inWindow;
  }

  // Group AFTER window membership is settled. A header therefore describes
  // exactly the REQs listed beneath it, not hidden members elsewhere in the
  // archive. The row order within each group and the first-seen group order are
  // inherited from the newest-first input. The no-UR bucket is the one
  // exception: it is explicit and always last.
  function timelineGroupWindowRows(
    windowRows,
    requestsById,
    durationSamples,
    nowMs,
    windowStartMs,
    windowEndMs
  ) {
    var durationSampleById = {};
    (durationSamples || []).forEach(function (sample) {
      durationSampleById[sample.id] = sample;
    });
    var groupsById = {};
    var orderedGroups = [];
    var unknownGroup = null;

    windowRows.forEach(function (row, rowIndex) {
      var request = requestsById[row.id] || {};
      var userRequestId = request.userRequestId || "";
      var groupKey = userRequestId || TIMELINE_UNKNOWN_USER_REQUEST_NAME;
      var group = groupsById[groupKey];
      if (!group) {
        group = {
          userRequestId: userRequestId,
          label: userRequestId || TIMELINE_UNKNOWN_USER_REQUEST_NAME,
          members: [],
          acceptedWorkMinutes: 0,
          acceptedWorkCount: 0,
          excludedReasons: {},
          unavailableWorkCount: 0,
          earliestClaimMs: Infinity,
          latestCompletionMs: -Infinity,
          recordedClaimCount: 0,
          overlappingElapsedCount: 0,
          unresolvedCompletionCount: 0,
          running: false,
          elapsedMinutes: null,
          elapsedUnavailableReason: ""
        };
        groupsById[groupKey] = group;
        if (userRequestId) {
          orderedGroups.push(group);
        } else {
          unknownGroup = group;
        }
      }

      group.members.push({ row: row, rowIndex: rowIndex });
      var claimMs = row.claimedTime ? Date.parse(row.claimedTime) : NaN;
      if (!isNaN(claimMs)) {
        group.recordedClaimCount++;
      }
      var completionMs = row.completedTime ? Date.parse(row.completedTime) : NaN;
      var rowRunning = !!(row.waitOpen || row.workOpen);
      if (!isNaN(claimMs) && isNaN(completionMs) && !rowRunning) {
        group.unresolvedCompletionCount++;
      }
      if (rowRunning) {
        group.running = true;
      }

      // A group's elapsed hull is made from each member's CLAIMED interval
      // intersected with this window. A row listed because its WAIT overlaps the
      // window but whose claim is later contributes no claimed interval here.
      // This is deliberately not a clamp of the final group hull: intersecting
      // each member first prevents an off-window member endpoint from expanding
      // the metric beyond what the header says it covers.
      var elapsedEndMs = rowRunning ? Math.min(nowMs, windowEndMs) : completionMs;
      var clippedClaimMs = Math.max(claimMs, windowStartMs);
      var clippedCompletionMs = Math.min(elapsedEndMs, windowEndMs);
      if (
        !isNaN(claimMs) &&
        isFinite(elapsedEndMs) &&
        clippedClaimMs < windowEndMs &&
        clippedCompletionMs >= windowStartMs &&
        clippedCompletionMs >= clippedClaimMs
      ) {
        group.earliestClaimMs = Math.min(group.earliestClaimMs, clippedClaimMs);
        group.latestCompletionMs = Math.max(group.latestCompletionMs, clippedCompletionMs);
        group.overlappingElapsedCount++;
      }

      var sample = durationSampleById[row.id];
      if (!sample) {
        group.unavailableWorkCount++;
      } else if (sample.excludedReason) {
        group.excludedReasons[sample.excludedReason] =
          (group.excludedReasons[sample.excludedReason] || 0) + 1;
      } else if (isFinite(sample.wallMinutes) && !isNaN(claimMs) && !isNaN(completionMs)) {
        // `excludedReason` remains the eligibility verdict. Once accepted, this
        // view contributes only the claim→completion interval that lies inside
        // its own visible window; the sample's full wallMinutes would otherwise
        // make a two-hour window announce four hours of work.
        var clippedWorkStartMs = Math.max(claimMs, windowStartMs);
        var clippedWorkEndMs = Math.min(completionMs, windowEndMs);
        group.acceptedWorkMinutes += Math.max(clippedWorkEndMs - clippedWorkStartMs, 0) / 60000;
        group.acceptedWorkCount++;
      } else {
        group.unavailableWorkCount++;
      }
    });

    if (unknownGroup) {
      orderedGroups.push(unknownGroup);
    }
    orderedGroups.forEach(function (group) {
      // One stopped claimed member without a completion can be the group's true
      // last endpoint. Publishing the latest completion from only its resolved
      // siblings would be a plausible-looking partial span, so uncertainty wins
      // over every computable member.
      if (group.unresolvedCompletionCount) {
        group.elapsedUnavailableReason = "completion endpoint unavailable";
        return;
      }
      if (!group.recordedClaimCount) {
        group.elapsedUnavailableReason = "no recorded claim";
        return;
      }
      if (!group.overlappingElapsedCount) {
        group.elapsedUnavailableReason = "no claimed interval overlaps this window";
        return;
      }
      group.elapsedMinutes = (group.latestCompletionMs - group.earliestClaimMs) / 60000;
    });
    return orderedGroups;
  }

  function timelineGroupDetailText(group) {
    var details = [];
    Object.keys(group.excludedReasons).sort().forEach(function (reason) {
      var count = group.excludedReasons[reason];
      details.push(count + " " + reason + " excluded");
    });
    if (group.unavailableWorkCount) {
      details.push(group.unavailableWorkCount + " without an accepted Durations sample");
    }
    if (group.unresolvedCompletionCount) {
      details.push(group.unresolvedCompletionCount + " unresolved completion");
    }
    return details.length ? details.join("; ") : "all listed work samples accepted";
  }

  function timelineGroupMetricText(group) {
    var elapsedText = group.elapsedMinutes === null
      ? "Elapsed in window unavailable (" + (group.elapsedUnavailableReason || "unknown endpoint") +
        (group.running ? "; running" : "") + ")"
      : "Elapsed in window " + timelineFormatSpanMinutes(group.elapsedMinutes) +
        (group.running ? " (running)" : "");
    var acceptedWorkText = group.acceptedWorkCount
      ? timelineFormatSpanMinutes(group.acceptedWorkMinutes)
      : "0 min";
    return elapsedText + " · Accepted work in window " + acceptedWorkText + " · " +
      group.members.length + " listed REQ" + (group.members.length === 1 ? "" : "s") +
      " · " + timelineGroupDetailText(group);
  }

  function timelineGroupHeaderMetricText(group) {
    var elapsedText = group.elapsedMinutes === null
      ? "Elapsed in window unavailable: " + (group.elapsedUnavailableReason || "unknown endpoint") +
        (group.running ? "; running" : "")
      : "Elapsed in window " + timelineFormatSpanMinutes(group.elapsedMinutes) +
        (group.running ? " running" : "");
    var acceptedWorkText = group.acceptedWorkCount
      ? timelineFormatSpanMinutes(group.acceptedWorkMinutes)
      : "0 min";
    var detailParts = [];
    Object.keys(group.excludedReasons).sort().forEach(function (reason) {
      detailParts.push(reason + " excluded×" + group.excludedReasons[reason]);
    });
    if (group.unavailableWorkCount) {
      detailParts.push("no sample×" + group.unavailableWorkCount);
    }
    if (group.unresolvedCompletionCount) {
      detailParts.push("unresolved completion×" + group.unresolvedCompletionCount);
    }
    return elapsedText + " · Accepted work in window " + acceptedWorkText + " · " +
      (detailParts.length ? detailParts.join("; ") : "no exclusions");
  }

  // Headers and members have fixed heights, but not the same height. Flattening
  // them gives virtualization and focus one display index while rowIndex remains
  // the REQ-only roving index used by the keyboard contract.
  function timelineFlattenGroups(groups) {
    var displayItems = [];
    var topPx = 0;
    groups.forEach(function (group) {
      displayItems.push({ kind: "group", group: group, topPx: topPx, height: TIMELINE_GROUP_HEADER_HEIGHT });
      topPx += TIMELINE_GROUP_HEADER_HEIGHT;
      group.members.forEach(function (member) {
        displayItems.push({
          kind: "request",
          row: member.row,
          rowIndex: member.rowIndex,
          group: group,
          topPx: topPx,
          height: TIMELINE_ROW_HEIGHT
        });
        topPx += TIMELINE_ROW_HEIGHT;
      });
    });
    return { items: displayItems, height: topPx };
  }

  // A deterministic, collision-free encoding of the UTF-16 label gives the
  // table's explicit `headers` relationships an id that survives rebuilds and
  // cannot be confused with a sibling group. The known/unknown prefix keeps a
  // literal UR named like the fallback label distinct from the fallback itself.
  function timelineTableGroupHeaderId(group) {
    var source = group.userRequestId
      ? "ur:" + group.userRequestId
      : "unknown:" + TIMELINE_UNKNOWN_USER_REQUEST_NAME;
    var encoded = [];
    for (var characterIndex = 0; characterIndex < source.length; characterIndex++) {
      encoded.push(source.charCodeAt(characterIndex).toString(16));
    }
    return "timeline-table-group-" + encoded.join("-");
  }

  function timelineVisibleDisplayRange(displayItems, scrollTop, viewportHeight) {
    var overscanPx = TIMELINE_OVERSCAN_ROWS * TIMELINE_ROW_HEIGHT;
    var rangeStartPx = Math.max(0, scrollTop - overscanPx);
    var rangeEndPx = scrollTop + viewportHeight + overscanPx;
    var firstDisplay = 0;
    while (
      firstDisplay < displayItems.length &&
      displayItems[firstDisplay].topPx + displayItems[firstDisplay].height <= rangeStartPx
    ) {
      firstDisplay++;
    }
    var lastDisplay = firstDisplay;
    while (lastDisplay < displayItems.length && displayItems[lastDisplay].topPx < rangeEndPx) {
      lastDisplay++;
    }
    return { firstDisplay: firstDisplay, lastDisplay: lastDisplay };
  }

  function timelineTabbableRequestIndex(rovingRowIndex, displayItems, firstDisplay, lastDisplay) {
    var renderedRequests = [];
    var rovingDisplayIndex = -1;
    for (var itemIndex = 0; itemIndex < displayItems.length; itemIndex++) {
      if (displayItems[itemIndex].kind === "request" && displayItems[itemIndex].rowIndex === rovingRowIndex) {
        rovingDisplayIndex = itemIndex;
        break;
      }
    }
    for (var displayIndex = firstDisplay; displayIndex < lastDisplay; displayIndex++) {
      if (displayItems[displayIndex].kind === "request") {
        renderedRequests.push({ rowIndex: displayItems[displayIndex].rowIndex, displayIndex: displayIndex });
      }
    }
    if (renderedRequests.length === 0) {
      return -1;
    }
    for (var renderedIndex = 0; renderedIndex < renderedRequests.length; renderedIndex++) {
      if (renderedRequests[renderedIndex].rowIndex === rovingRowIndex) {
        return rovingRowIndex;
      }
    }
    if (rovingDisplayIndex >= 0 && rovingDisplayIndex > renderedRequests[renderedRequests.length - 1].displayIndex) {
      return renderedRequests[renderedRequests.length - 1].rowIndex;
    }
    return renderedRequests[0].rowIndex;
  }

  // Which rows have SVG nodes. Everything above and below the scrolled window is
  // absent from the DOM, which is what keeps 560 rows and 5600 rows the same
  // cost. The overscan is what stops a fast scroll showing blank strips.
  function timelineVisibleRowRange(scrollTop, viewportHeight, rowCount) {
    var firstRow = Math.max(0, Math.floor(scrollTop / TIMELINE_ROW_HEIGHT) - TIMELINE_OVERSCAN_ROWS);
    var visibleCount = Math.ceil(viewportHeight / TIMELINE_ROW_HEIGHT) + TIMELINE_OVERSCAN_ROWS * 2;
    return { firstRow: firstRow, lastRow: Math.min(rowCount, firstRow + visibleCount) };
  }

  // Which of the RENDERED rows carries the list's single tabindex="0".
  //
  // The rows are virtualized, so the roving row is often not among them — and a render
  // that then marked nothing tabbable would take the whole list out of the Tab order,
  // which is worse than the 29 stops it replaced. Clamping into the rendered range keeps
  // exactly one stop and keeps it reachable, without touching the stored roving index:
  // that still names the row the reader was on, and scrolling back brings it right back.
  function timelineTabbableRowIndex(rovingRowIndex, firstRenderedRow, lastRenderedRow) {
    var lastSelectable = Math.max(firstRenderedRow, lastRenderedRow - 1);
    return Math.min(Math.max(rovingRowIndex, firstRenderedRow), lastSelectable);
  }

  // Whether a press has become a pan. Latching: once engaged a drag stays
  // engaged however far back toward the press point it travels, or a slow
  // wandering drag would flicker in and out of panning around the threshold.
  //
  // Engagement gates only the FIRST render — the shift itself is still measured
  // from the press point, so the chart tracks the pointer one-for-one and never
  // settles a threshold's worth behind it.
  function timelinePanEngaged(alreadyEngaged, pressX, pointerX) {
    return alreadyEngaged || Math.abs(pointerX - pressX) >= TIMELINE_PAN_THRESHOLD_PX;
  }

  // ---- axis ticks -----------------------------------------------------------
  //
  // Ticks sit on CALENDAR boundaries, chosen from TIMELINE_AXIS_TICK_STEPS. The
  // axis used to divide the window into TIMELINE_AXIS_TICK_COUNT equal parts,
  // which is honest only when the window happens to be a multiple of six of
  // something — and the label formatter printed a bare date for any gap of a day
  // or more, so on a week (28-hour gaps) and on a month (5.17-day gaps) the axis
  // named midnights it was not drawn at and skipped a calendar day outright.

  // How long one step is, as a span, so the ladder can be searched by span. A
  // month is not a fixed length, so a month step reports the mean Gregorian month
  // — good enough to CHOOSE with, and never used to position anything.
  function timelineTickStepSpanMs(step) {
    return step.months ? step.months * 30.436875 * TIMELINE_DAY_MS : step.ms;
  }

  // Which rung of the ladder this window gets: the one dividing it into closest to
  // TIMELINE_AXIS_TICK_COUNT gaps. Ties go to the LONGER step, because the risk at
  // equal distance is a crowded axis rather than a sparse one.
  function timelineAxisTickStep(windowSpanMs) {
    var chosenStep = TIMELINE_AXIS_TICK_STEPS[TIMELINE_AXIS_TICK_STEPS.length - 1];
    var chosenDistance = Infinity;
    for (var stepIndex = 0; stepIndex < TIMELINE_AXIS_TICK_STEPS.length; stepIndex++) {
      var step = TIMELINE_AXIS_TICK_STEPS[stepIndex];
      var distance = Math.abs(
        windowSpanMs / timelineTickStepSpanMs(step) - TIMELINE_AXIS_TICK_COUNT);
      if (distance <= chosenDistance) {
        chosenStep = step;
        chosenDistance = distance;
      }
    }
    return chosenStep;
  }

  // The last boundary of this step at or before an instant.
  function timelineTickAtOrBefore(epochMs, step) {
    if (step.months) {
      var instant = new Date(epochMs);
      var monthIndex = instant.getUTCFullYear() * 12 + instant.getUTCMonth();
      var alignedMonths = monthIndex - (((monthIndex % step.months) + step.months) % step.months);
      // Counted in months from January 1970 rather than passed as a year, because
      // Date.UTC reads a year under 100 as 19xx. Month indices well past 11 roll
      // forward correctly, which is the whole reason this is one call.
      return Date.UTC(1970, alignedMonths - 1970 * 12, 1);
    }
    var alignmentMs = step.alignMs || 0;
    return Math.floor((epochMs - alignmentMs) / step.ms) * step.ms + alignmentMs;
  }

  // One step on from a boundary.
  function timelineSteppedTick(epochMs, step) {
    if (step.months) {
      var instant = new Date(epochMs);
      return Date.UTC(instant.getUTCFullYear(), instant.getUTCMonth() + step.months, 1);
    }
    return epochMs + step.ms;
  }

  // The instants the axis puts a tick at, and the gap it chose between them.
  // Extracted so the GRIDLINES can be drawn from the same list rather than
  // recomputing it: two loops with the same arithmetic is one edit away from an
  // axis whose labels sit beside lines that mean something slightly different.
  //
  // Returns the gap alongside the instants because the LABEL FORMAT depends on it,
  // and deriving it a second time from the instants is how the two came to
  // disagree in the first place.
  function timelineAxisTicks(windowStartMs, windowEndMs) {
    var windowSpanMs = windowEndMs - windowStartMs;
    if (!(windowSpanMs > 0)) {
      // A degenerate window has no interior to space ticks across. One tick at the
      // instant itself is finite and honest; six identical ones were neither.
      return { instants: [windowStartMs], gapMs: TIMELINE_DAY_MS };
    }
    var step = timelineAxisTickStep(windowSpanMs);
    var instants = [];
    // A tick AT windowEndMs is kept, deliberately. The window is half-open, so a
    // week ending 13 July 00:00 draws a tick labelled "13 Jul" while the `to` field
    // beside it calls 12 July the last day included — which reads like a
    // contradiction and is not one. An axis tick names an INSTANT; the fields name
    // DAYS OF COVERAGE. Hiding the boundary tick would make the week look as though
    // it ended at 12 July 00:00, losing a day of the plot to avoid a distinction the
    // readout already states exactly.
    var tickMs = timelineTickAtOrBefore(windowStartMs, step);
    if (tickMs < windowStartMs) {
      tickMs = timelineSteppedTick(tickMs, step);
    }
    while (tickMs <= windowEndMs && instants.length < TIMELINE_AXIS_TICK_LIMIT) {
      instants.push(tickMs);
      tickMs = timelineSteppedTick(tickMs, step);
    }
    if (instants.length === 0) {
      // A window narrower than the shortest rung, which the zoom floor makes
      // unreachable today. It still may not produce an axis with nothing on it.
      instants.push(windowStartMs);
    }
    return { instants: instants, gapMs: timelineTickStepSpanMs(step) };
  }

  // The instants alone, for the gridlines, which need no label and therefore no
  // gap. One function so a gridline can never mean a different instant from the
  // label above it.
  function timelineAxisTickInstants(windowStartMs, windowEndMs) {
    return timelineAxisTicks(windowStartMs, windowEndMs).instants;
  }

  // ---- row labels -----------------------------------------------------------
  //
  // The label column used to hold the id and nothing else, so three hundred rows
  // told the reader nothing about what any of them were. Fitting a title into it
  // needs a truncation boundary, and this module has been burned twice by where
  // that number came from: REQ-292 shipped a width model that returned the same
  // value for every face, so its slots never moved and a wider face drew past
  // them, and REQ-241/REQ-242 measured one 12px face at 12.0372 and 11.2300 on
  // two Chromium builds.
  //
  // So the boundary is MEASURED from the rendered face, and measured ONCE per
  // render rather than per row: the label is a monospace face, so a single
  // advance answers every row. That is a real property of the shipped CSS
  // (.timeline-row-label uses --font-mono) and not an assumption this code is
  // entitled to — timelineMeasureLabelAdvance verifies it and says so.

  // The per-character advance of the label face, measured against the live DOM.
  // Returns 0 when the face is NOT monospace, which is the honest answer rather
  // than a number that would silently mis-cut every title: two strings of equal
  // length that render at different widths mean one advance cannot describe the
  // face, and the caller falls back to showing the id alone.
  function timelineMeasureLabelAdvance(rowsSvg) {
    // Not every host can measure text. The Node behavior lane drives this
    // renderer against a DOM stub, and a stub that cannot lay text out has no
    // advance to report — the same honest answer as a proportional face, and the
    // same consequence: labels fall back to the id alone. Checked up front
    // rather than caught, because a throw here would take the whole render with
    // it for a label that was never going to be measurable.
    if (!rowsSvg || typeof rowsSvg.appendChild !== "function" || typeof rowsSvg.removeChild !== "function") {
      return 0;
    }
    var probeText = document.createElementNS(TIMELINE_SVG_NS, "text");
    if (!probeText || typeof probeText.getComputedTextLength !== "function") {
      return 0;
    }
    probeText.setAttribute("class", "timeline-row-label");
    // Deliberately off-canvas rather than hidden: display:none and
    // visibility:hidden both make getComputedTextLength return 0, which would
    // read as "not monospace" and quietly disable the whole feature.
    probeText.setAttribute("x", "-9999");
    probeText.setAttribute("y", "-9999");
    rowsSvg.appendChild(probeText);

    function widthOf(sample) {
      probeText.textContent = sample;
      return probeText.getComputedTextLength ? probeText.getComputedTextLength() : 0;
    }
    // Two samples of equal length whose glyphs have very different proportional
    // widths. In a monospace face they measure the same; in any proportional one
    // they do not, which is exactly the check REQ-292's lesson asks for — assert
    // the geometry responds to the face rather than trusting a model of it.
    var narrowSample = widthOf("iiiiiiiiii");
    var wideSample = widthOf("MMMMMMMMMM");
    rowsSvg.removeChild(probeText);

    if (!(narrowSample > 0) || Math.abs(narrowSample - wideSample) > 0.5) {
      return 0;
    }
    return narrowSample / 10;
  }

  // How many characters of label fit, given the measured advance. Zero advance
  // means the face could not be measured or is not monospace.
  function timelineLabelCharacterBudget(advancePerCharacter, availableWidth) {
    if (!(advancePerCharacter > 0)) {
      return 0;
    }
    return Math.max(0, Math.floor(availableWidth / advancePerCharacter));
  }

  // How many monospace CELLS a string occupies, and the reason this is not just
  // .length.
  //
  // The measured advance describes the face's Latin cell. A monospace face does
  // not draw every script in one cell: on the shipped 10px face an ASCII glyph
  // is 6.02px, but 中 is 10px and 🙂 is 12.48px. Counting those as one cell each
  // let a CJK title overrun the column by 36px and draw into the plot. So
  // non-ASCII counts as two cells — the East Asian Width convention every
  // monospace terminal uses — which OVER-estimates slightly and therefore cuts
  // early rather than overflowing. Approximate on purpose: the alternative is
  // measuring every row's composed string every frame, which is the cost this
  // whole approach exists to avoid.
  //
  // Iterating code points rather than UTF-16 units also fixes the other half:
  // .slice() on a unit boundary can cut an astral character in two and leave a
  // lone surrogate, which renders as a fallback box.
  function timelineLabelCellCount(text) {
    var cells = 0;
    for (var index = 0; index < text.length; index++) {
      var codePoint = text.codePointAt(index);
      if (codePoint > 0xffff) {
        index++;
      }
      cells += codePoint < 0x0100 ? 1 : 2;
    }
    return cells;
  }

  // The longest prefix of `text` that fits `cellBudget` cells, cut on a code
  // point boundary.
  function timelineLabelPrefixWithinCells(text, cellBudget) {
    var cells = 0;
    var prefix = "";
    for (var index = 0; index < text.length; index++) {
      var codePoint = text.codePointAt(index);
      var character = String.fromCodePoint(codePoint);
      var width = codePoint < 0x0100 ? 1 : 2;
      if (cells + width > cellBudget) {
        break;
      }
      cells += width;
      prefix += character;
      if (codePoint > 0xffff) {
        index++;
      }
    }
    return prefix;
  }

  // "REQ-042  Some title" cut to the budget, with an ellipsis when it was cut.
  // The id is never truncated: it is the row's identity and a half-drawn id is
  // worse than no title at all, so a budget too small for id-plus-title returns
  // the id alone.
  //
  // The title is measured in CELLS, not characters, and a leading
  // "[impact-token] " classification tag is dropped: the tag exists so a human
  // searching the board's title box finds the REQ (actions/capture-reference.md
  // → REQ Title Convention), and in a nineteen-cell column it is the entire
  // budget — the sixteen newest REQs here all read "[impact-user-visib…" and
  // named nothing. The full title, tag included, stays in the tooltip and the
  // table.
  function timelineRowLabelText(requestId, title, characterBudget) {
    if (characterBudget <= 0 || !title) {
      return requestId;
    }
    var labelTitle = title.replace(/^\[[a-z-]+\]\s*/, "");
    if (labelTitle === "") {
      return requestId;
    }
    var separator = "  ";
    var roomForTitle = characterBudget - timelineLabelCellCount(requestId) - separator.length;
    // A title that fits WHOLE always shows, however little room there is: a
    // complete short title is useful at any length.
    if (timelineLabelCellCount(labelTitle) <= roomForTitle) {
      return requestId + separator + labelTitle;
    }
    // One that has to be cut has to be worth cutting. Four characters and an
    // ellipsis names nothing and just makes the column look broken, so below
    // this the id stands alone rather than being followed by a stub.
    if (roomForTitle < 8) {
      return requestId;
    }
    // The ellipsis lives INSIDE the budget, not on top of it, or the label the
    // arithmetic says fits is one cell wider than the column.
    return requestId + separator +
      timelineLabelPrefixWithinCells(labelTitle, roomForTitle - 1).replace(/\s+$/, "") + "…";
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

  // The extent of a row set, from everything its rows are DRAWN across. Used only
  // by the bounds FALLBACK below, so it costs the common path nothing.
  //
  // The now-line is part of that extent for an OPEN row, and leaving it out was a
  // real gap: timelineRowSegments draws an open span to nowMs, so an extent built
  // from stored stamps alone stopped at the row's last recorded instant. Since these
  // bounds are what every control clamps against, the live part of the bar then sat
  // outside every window the reader could reach — Fit all included. (Found by Codex
  // on the pull request. It only bites when the projection is declined or absent,
  // because a confident projection's queue-end already pushes the bound past now.)
  //
  // Returns null when not one row carries a parseable instant, because the honest
  // answer there is "there is nothing to place on a timeline" rather than a
  // fabricated window.
  function timelineRowSetExtent(rows, nowMs) {
    var earliestMs = Infinity;
    var latestMs = -Infinity;
    rows.forEach(function (row) {
      var rowCarriesAnInstant = false;
      [row.createdTime, row.claimedTime, row.completedTime].forEach(function (stamp) {
        if (!stamp) {
          return;
        }
        var instantMs = Date.parse(stamp);
        if (isNaN(instantMs)) {
          return;
        }
        rowCarriesAnInstant = true;
        earliestMs = Math.min(earliestMs, instantMs);
        latestMs = Math.max(latestMs, instantMs);
      });
      // An open span is drawn FROM one of the instants above TO the now-line, so
      // the now-line only extends an extent the row is already in. A row with no
      // readable instant of its own has no segment for it to extend, and counting
      // the now-line there would turn "there is nothing to place on a timeline"
      // into a fabricated hour around the present.
      //
      // Both ends, because timelineRowSegments sorts its own endpoints: a
      // future-dated created_at makes the now-line the EARLIER of the two.
      if (rowCarriesAnInstant && (row.waitOpen || row.workOpen) && !isNaN(nowMs)) {
        earliestMs = Math.min(earliestMs, nowMs);
        latestMs = Math.max(latestMs, nowMs);
      }
    });
    if (!isFinite(earliestMs) || !isFinite(latestMs)) {
      return null;
    }
    return { earliestMs: earliestMs, latestMs: latestMs };
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
        // "every remaining REQ is listed below" is only true when the rows on
        // screen ARE the whole queue — which is exactly what the prefix above
        // denies. Both appeared together under a filter: "This covers the whole
        // queue, not the rows shown. Nothing left to schedule — every remaining REQ
        // is listed below.", over one row, with the excluded paragraph naming a REQ
        // that was not listed anywhere.
        forecastNode.textContent = showingSubset
          ? wholeQueueNote + "Nothing left to schedule."
          : "Nothing left to schedule — every remaining REQ is listed below.";
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

  // Every control this view wires up, so a render that leaves without re-wiring
  // them can take them all out of service together.
  //
  // The condition is the rule: a control this view owns is either wired to the
  // CURRENT render or visibly disabled — never silently wired to the LAST one.
  // The toolbar is bound with `button.onclick =`, which is outside the listener
  // teardown registry, so the no-match early return used to leave every handler
  // alive holding the previous render's rows, its detached rows SVG and its
  // renderAll. One press of Fit all in that state refilled the summary, the
  // forecast and the details table with the REQs the filter had excluded, over a
  // chart that stayed empty.
  function timelineOwnedControls() {
    return [].slice.call(
      document.querySelectorAll(
        "#view-timeline .timeline-periods button," +
          "#view-timeline .timeline-zoom button," +
          "#view-timeline .timeline-range input"
      )
    );
  }

  // Hand every control back with no handler and no way to be pressed. Called on
  // the one path that renders nothing, so nothing the reader can press describes
  // a chart that is not there.
  function retireTimelineControls() {
    timelineOwnedControls().forEach(function (control) {
      control.onclick = null;
      control.disabled = true;
    });
  }

  function enableTimelineControls() {
    timelineOwnedControls().forEach(function (control) {
      control.disabled = false;
    });
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
    releaseTimelineFrameRender();

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
    // The first real display list is built by refreshWindowRows after the view
    // window is fitted. Building metrics here would either need an invented
    // window or briefly measure the full archive before the first render.
    var timelineGroups = [];
    var timelineDisplay = timelineFlattenGroups(timelineGroups);

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
      // And so must every control: this path wires nothing, so anything still
      // wired belongs to a render whose rows the filter has excluded.
      retireTimelineControls();
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
      // THE FALLBACK SPANS THE WHOLE MATCHED SET, not one row of it. It used to
      // read filterMatchedRows[0].createdTime plus an hour — and rows are
      // newest-first (REQ-318), so [0] is the NEWEST capture. The bounds collapsed
      // to a one-hour window around it, and because bounds are what every control
      // clamps against, no control could leave: on this repo's board that would
      // have stranded 287 of 317 REQs permanently out of reach. A degraded
      // fallback may be coarse; it may not be a dead end.
      var matchedExtent = timelineRowSetExtent(filterMatchedRows, nowMs);
      if (!matchedExtent) {
        // Nothing parseable anywhere. Say so rather than invent a window; this is
        // the same message the no-readable-created_at path uses.
        clearTimelineForecast();
        retireTimelineControls();
        summaryNode.textContent =
          "No REQ carries a readable created_at yet, so there is nothing to place on a timeline.";
        return;
      }
      boundStartMs = matchedExtent.earliestMs;
      boundEndMs = Math.max(matchedExtent.latestMs, boundStartMs + TIMELINE_MIN_SPAN_MS);
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
    // The same list keyed by id, so the row loop measures a row's drawn extent
    // from the list that decided the row was in the window rather than from a
    // second walk of the stamps.
    var segmentsById = {};
    filterMatchedRows.forEach(function (row, rowIndex) {
      segmentsById[row.id] = filterMatchedSegments[rowIndex];
    });

    // The reader's place in the list, preserved across a window move. Rows
    // vanishing above the viewport would otherwise slide whatever they were
    // reading off under the axis. Anchored on the id at the top of the viewport:
    // if it survives the move it keeps its screen position, and if it does not,
    // the scroll simply stays where it was and the browser clamps it.
    function topVisibleRowAnchor() {
      if (timelineDisplay.items.length === 0) {
        return null;
      }
      for (var displayIndex = 0; displayIndex < timelineDisplay.items.length; displayIndex++) {
        var displayItem = timelineDisplay.items[displayIndex];
        if (
          displayItem.kind === "request" &&
          displayItem.topPx + displayItem.height > scrollHost.scrollTop
        ) {
          return { id: displayItem.row.id, offsetPx: scrollHost.scrollTop - displayItem.topPx };
        }
      }
      return null;
    }

    function refreshWindowRows() {
      var anchor = topVisibleRowAnchor();
      rows = timelineRowsInWindow(
        filterMatchedRows,
        filterMatchedSegments,
        timelineViewState.windowStartMs,
        timelineViewState.windowEndMs
      );
      timelineGroups = timelineGroupWindowRows(
        rows,
        requestsById,
        ((boardData.durations || {}).samples || []),
        nowMs,
        timelineViewState.windowStartMs,
        timelineViewState.windowEndMs
      );
      timelineDisplay = timelineFlattenGroups(timelineGroups);
      if (anchor === null) {
        return;
      }
      for (var displayIndex = 0; displayIndex < timelineDisplay.items.length; displayIndex++) {
        var displayItem = timelineDisplay.items[displayIndex];
        if (displayItem.kind === "request" && displayItem.row.id === anchor.id) {
          // GROW THE EXTENT FIRST. scrollTop is clamped to the scrollable height
          // at the instant it is assigned, and on a widening move the rows SVG
          // is still the narrow window's height here — renderVisibleRows only
          // resizes it afterwards. Writing the anchor first therefore clamped it
          // to the OLD maximum and dropped the reader at the bottom of the
          // window they were leaving: from a month window at scrollTop 400 the
          // Fit-all button landed on 465, the old extent's maximum, when the
          // anchor needed 4900.
          rowsSvg.setAttribute("height", timelineDisplay.height);
          scrollHost.scrollTop = displayItem.topPx + anchor.offsetPx;
          return;
        }
      }
    }

    function renderSummary() {
      var openCount = rows.filter(function (row) {
        return row.waitOpen || row.workOpen;
      }).length;
      // The SAME predicate the renderer branches on, so the sentence below cannot
      // announce breaks the chart does not draw. It counted row.anomaly too, which
      // is the board's broader bookkeeping verdict and includes rows with a
      // perfectly drawable span: nine such rows on this repo's board produced
      // "9 with broken stamps, drawn as breaks" over a chart with zero break
      // markers on it.
      var brokenRowCount = rows.filter(function (row) {
        return timelineRowDrawsABreak(row, nowMs);
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
        // NAMING THE NOW-LINE IS A CLAIM THAT IT IS ON SCREEN. drawNowRule draws
        // nothing when now falls outside the window, so on any past week this
        // sentence pointed at a rule the reader could not find, beside open bars
        // clipped flush at the frame. The INSTANT still has to be stated — it is
        // what every open span was measured against — so only the pointer goes.
        (nowMs >= timelineViewState.windowStartMs && nowMs <= timelineViewState.windowEndMs
          ? " still open, measured to the now-line at "
          : " still open, measured against ") +
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
      timelineGroups.forEach(function (group) {
        var groupRow = document.createElement("tr");
        groupRow.className = "timeline-table-group";
        groupRow.setAttribute("data-timeline-table-group", group.label);
        groupRow.setAttribute("data-group-count", String(group.members.length));
        var groupHeading = document.createElement("th");
        var groupHeadingId = timelineTableGroupHeaderId(group);
        groupHeading.id = groupHeadingId;
        groupHeading.colSpan = 6;
        groupHeading.textContent = group.label + " · " + timelineGroupMetricText(group);
        groupRow.appendChild(groupHeading);
        tableBody.appendChild(groupRow);

        group.members.forEach(function (member) {
          var row = member.row;
          var request = requestsById[row.id] || {};
          var tableRow = document.createElement("tr");
          tableRow.setAttribute("data-timeline-table-request", row.id);
          [
            { text: row.id, columnHeaderId: "timeline-table-column-req", rowHeader: true },
            { text: request.title || "", columnHeaderId: "timeline-table-column-title" },
            { text: request.status || "", columnHeaderId: "timeline-table-column-status" },
            {
              text: timelineFormatSpanMinutes(row.waitMinutes) + (row.waitOpen ? " (open)" : ""),
              columnHeaderId: "timeline-table-column-waited"
            },
            {
              text: row.hasWork
                ? timelineFormatSpanMinutes(row.workMinutes) + (row.workOpen ? " (open)" : "")
                : "not started",
              columnHeaderId: "timeline-table-column-worked"
            },
            { text: row.anomaly ? row.anomalyReason : "", columnHeaderId: "timeline-table-column-note" }
          ].forEach(function (cellDefinition) {
            var cell = document.createElement(cellDefinition.rowHeader ? "th" : "td");
            if (cellDefinition.rowHeader) {
              cell.scope = "row";
            }
            cell.setAttribute("headers", groupHeadingId + " " + cellDefinition.columnHeaderId);
            cell.textContent = cellDefinition.text;
            tableRow.appendChild(cell);
          });
          tableBody.appendChild(tableRow);
        });
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

    // Past the early return, so this render owns the controls and is about to wire
    // them. Undoes a retirement left by a previous no-match render.
    enableTimelineControls();

    var axisSvg = makeTimelineSvgNode(axisHost, "svg", {
      class: "timeline-axis-svg",
      height: TIMELINE_AXIS_HEIGHT,
      width: "100%"
    });
    var rowsSvg = makeTimelineSvgNode(scrollHost, "svg", {
      class: "timeline-rows-svg",
      height: timelineDisplay.height,
      width: "100%",
      role: "img",
      "aria-label":
        "Window-listed REQs grouped by user request, with the No UR recorded group last. Group elapsed and accepted work are clipped to this visible window. Elapsed spans the earliest overlapping recorded member claim through the latest overlapping resolved completion, or a claimed open member's endpoint bounded by this board's frozen now and the window end. It is unavailable when there is no recorded claim, no claimed interval overlaps the window, or a stopped claimed member has no completion endpoint. Accepted work uses only the in-window claim-to-completion overlap of Durations samples whose read-time verdict did not exclude them. REQ rows remain newest first within each group. The first bar segment is waiting and the second is work. Every value is also listed in the grouped table below."
    });

    // MEASURED ONCE PER RENDER. clientWidth forces a synchronous layout, and
    // xOfEpoch calls this several times per row: one measured render of this
    // repo's board made 171 of them (Chromium 1194), every frame of a drag. The
    // memo is dropped by invalidatePlotWidth() at the top of renderAll.
    //
    // WHAT RE-RENDERS IS NOT A LIST OF CALLERS. This used to say "the resize
    // listener IS renderAll, so the two moments the width can change are both
    // covered", and there were more than two moments. Opening the detail drawer
    // gives it a 620px grid column of its own, narrowing this host by 630px:
    // nothing re-measured, every bar kept the x it had been given against the old
    // 1300px plot, and drawSegment's clamp dropped all fifty of them — clicking a
    // row, which the hint under the chart explicitly invites, blanked the chart.
    // So the trigger is now one shared positive condition: the live host width is
    // non-zero and differs from the width used by the last render. ResizeObserver
    // delivers that check directly in ordinary browser layout; the teardown-owned
    // 50ms fallback below delivers the same check if observer/compositor work is
    // parked. Neither path names a layout caller, so no future cause has to be
    // added to a list.
    //
    // A ZERO-SIZE BOX IS NOT A WIDTH. A window resize fired while the view is
    // hidden measured clientWidth 0, and the floor below turned that into
    // Math.max(120, -196) = 120 — three months of archive crushed into a 120px
    // strip with eight rows in it, and nothing re-rendered afterwards to repair it.
    // renderAll now refuses to render at all against an unmeasurable host, and the
    // next delivery of the shared positive width-change condition brings it back
    // when the box exists. That also repairs re-entry without a second gate in
    // board-controls.js.
    //
    // That guard is the ONLY one. plotWidth deliberately keeps its floor and its
    // memo unconditional: every caller runs inside a render the guard has already
    // admitted, so a defensive branch here would be unreachable, and a guard no
    // mutation can break is dead code rather than safety. Should a call ever slip
    // through, invalidatePlotWidth at the top of the next renderAll drops the memo,
    // so a bad measurement cannot outlive the render that took it.
    var measuredPlotWidthPx = null;
    var renderedHostWidthPx = null;
    function liveHostWidth() {
      return scrollHost.clientWidth || scrollHost.getBoundingClientRect().width || 0;
    }
    function plotIsMeasurable() {
      return liveHostWidth() > 0;
    }
    function plotWidth() {
      if (measuredPlotWidthPx === null) {
        var hostWidth = liveHostWidth();
        renderedHostWidthPx = hostWidth;
        measuredPlotWidthPx = Math.max(120, hostWidth - TIMELINE_LABEL_WIDTH - 12);
      }
      return measuredPlotWidthPx;
    }
    function invalidatePlotWidth() {
      measuredPlotWidthPx = null;
    }

    // The row the hover readout last announced, so an unchanged row is not
    // announced again and a window move can retract one that is gone.
    var announcedHoverRowId = null;
    function clearHoverReadout() {
      announcedHoverRowId = null;
      if (readoutNode) {
        readoutNode.textContent = "";
      }
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
        width: Math.max(TIMELINE_MIN_SEGMENT_WIDTH, clampedRight - clampedLeft).toFixed(1),
        height: TIMELINE_BAR_HEIGHT,
        rx: 2,
        class: segmentClass
      });
    }

    // Measured once per RENDER and not once per row — which is the cost that
    // mattered, since a render draws ~35 rows. It is not once per session:
    // renderVisibleRows is the scroll listener and runs inside renderAll, which
    // the pan handler calls, so this measures on each frame of a drag. Measured
    // cost on Chromium 1194: 0.036–0.046ms alone, about 0.13ms of a 0.61ms
    // 35-row render. It forces one synchronous layout per frame and leaves room
    // for roughly 27 renders inside a 16.7ms frame, so it stays here rather than
    // being hoisted — one place that measures beats two that can disagree.
    var labelAdvance = 0;
    var labelCharacterBudget = 0;

    function renderVisibleRows() {
      // WHICH ROW HELD FOCUS, captured before the rebuild destroys it.
      //
      // A scroll event is ASYNCHRONOUS. refreshWindowRows writes scrollHost.scrollTop
      // to keep the reader's place, and the scroll listener is this function — so the
      // keydown handler's own focus restore ran, and THEN this rebuild wiped the node
      // it had just focused, and focus fell to <body>. One arrow press panned the
      // window and every arrow press after it was dead, because the keydown listener
      // is on the scroll host and <body> is not inside it.
      //
      // Restoring here covers every path that rebuilds rows — the keyboard, a scroll,
      // a drag frame, a filter change — rather than only the one the keyboard takes,
      // which is why the keydown handler no longer does it itself.
      var focusedRowIdBeforeRebuild = null;
      if (
        typeof rowsSvg.contains === "function" &&
        document.activeElement &&
        rowsSvg.contains(document.activeElement) &&
        document.activeElement.closest
      ) {
        var focusedRowBeforeRebuild = document.activeElement.closest("[data-detail-id]");
        focusedRowIdBeforeRebuild = focusedRowBeforeRebuild
          ? focusedRowBeforeRebuild.getAttribute("data-detail-id")
          : null;
      }
      rowsSvg.textContent = "";
      // The scroll extent follows the windowed set, not the filtered one — a
      // window holding four rows must not leave 312 rows of empty scroll below
      // them.
      rowsSvg.setAttribute("height", timelineDisplay.height);
      appendTimelineHatchPattern(rowsSvg);
      // TIMELINE_LABEL_WIDTH is the column; the 6px is the text's own x offset
      // and the trailing 6px keeps the longest label off the first bar.
      labelAdvance = timelineMeasureLabelAdvance(rowsSvg);
      labelCharacterBudget = timelineLabelCharacterBudget(labelAdvance, TIMELINE_LABEL_WIDTH - 12);
      // Before the rows, so SVG paint order puts every gridline behind every bar
      // without a z-index to maintain.
      drawGridlines();
      // Read once here and handed to the collapse decision, which runs per row.
      // drawSegment still calls plotWidth() itself; hoisting THAT is REQ-324's
      // job, and doing half of it here would leave two sources for one number.
      var plotWidthPx = plotWidth();
      var visible = timelineVisibleDisplayRange(
        timelineDisplay.items, scrollHost.scrollTop, scrollHost.clientHeight
      );
      var tabbableRowIndex = timelineTabbableRequestIndex(
        timelineViewState.rovingRowIndex,
        timelineDisplay.items,
        visible.firstDisplay,
        visible.lastDisplay
      );
      for (
        var displayIndex = visible.firstDisplay;
        displayIndex < visible.lastDisplay;
        displayIndex++
      ) {
        var displayItem = timelineDisplay.items[displayIndex];
        if (displayItem.kind === "group") {
          var group = displayItem.group;
          var groupHeader = makeTimelineSvgNode(rowsSvg, "g", {
            class: "timeline-group-header",
            "data-timeline-group": group.label,
            "data-group-count": String(group.members.length),
            "aria-hidden": "true"
          });
          makeTimelineSvgNode(groupHeader, "title", {},
            group.label + " · " + timelineGroupMetricText(group));
          makeTimelineSvgNode(groupHeader, "rect", {
            x: 0,
            y: displayItem.topPx,
            width: "100%",
            height: displayItem.height,
            class: "timeline-group-header-fill"
          });
          makeTimelineSvgNode(groupHeader, "text", {
            x: 8,
            y: displayItem.topPx + 13,
            class: "timeline-group-label"
          }, group.label + " · " + group.members.length + " listed REQ" +
            (group.members.length === 1 ? "" : "s"));
          makeTimelineSvgNode(groupHeader, "text", {
            x: 8,
            y: displayItem.topPx + 27,
            class: "timeline-group-metrics"
          }, timelineGroupHeaderMetricText(group));
          continue;
        }
        var rowIndex = displayItem.rowIndex;
        var row = displayItem.row;
        var request = requestsById[row.id] || {};
        var rowTopY = displayItem.topPx;
        // data-status carries the hue, exactly as the calendar chips do, and the
        // unrecognized verdict is READ FROM THE PAYLOAD rather than re-derived:
        // the board already decides what counts as a real status, and a second
        // reader here would become a second definition (REQ-219's lesson, in
        // this module's own prime). A typo like "blockd-dependency-cycle" is
        // therefore coloured as unrecognized, never quietly as real blocked work.
        var rowGroup = makeTimelineSvgNode(rowsSvg, "g", {
          class:
            "timeline-row" +
            (row.anomaly ? " is-broken" : "") +
            (request.status === "cancelled" ? " is-cancelled" : "") +
            (request.statusUnrecognized ? " is-status-unrecognized" : ""),
          "data-detail-kind": "request",
          "data-detail-id": row.id,
          "data-row-index": String(rowIndex),
          "data-status": request.status || "",
          // Roving: one stop for the whole list. Every other row is EXPLICITLY skipped
          // rather than left without the attribute, because a focusable-by-default
          // element with no tabindex is still a Tab stop.
          tabindex: rowIndex === tabbableRowIndex ? "0" : "-1",
          role: "button"
        });
        // The accessible name comes from the <title> below and NOT from an
        // aria-label. Carrying both put the same 150-character sentence in the
        // name and the description of all three hundred rows, so a screen reader
        // read it twice per row. One source, one announcement — and the <title>
        // has to exist anyway, because it is the pointer tooltip.
        //
        // First child, per SVG's own guidance on <title> placement.
        makeTimelineSvgNode(rowGroup, "title", {}, timelineRowDescription(row, request));
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
          timelineRowLabelText(row.id, request.title, labelCharacterBudget)
        );


        var createdMs = Date.parse(row.createdTime);
        var claimedMs = row.claimedTime ? Date.parse(row.claimedTime) : NaN;
        var completedMs = row.completedTime ? Date.parse(row.completedTime) : NaN;
        // The SAME two answers timelineRowSegments used to decide this row belongs
        // in the window. Reading them from one place is what keeps the row's marks
        // where the segment model said they would be — the claimedMs default used
        // to be nowMs here, which drew a bar to the now-line for a row the model
        // had already decided was a point.
        var waitEndMs = timelineWaitEndMs(row, nowMs, claimedMs, completedMs);
        var projectedRow = projectedById[row.id];
        // A row whose whole span is narrower than a readable two-segment bar
        // draws ONE marker. Two kinds of row are excluded, for one reason: the
        // collapse withdraws the wait/work split, and these carry a SECOND
        // distinction it would withdraw with it. Broken stamps, because the break
        // markers are the point of drawing the row at all. A forecast, because
        // measured-versus-projected is the distinction a reader trusts hardest,
        // and one solid marker over an open wait plus a hatched projection claims
        // work that has not happened.
        var rowHasBrokenStamps = timelineRowDrawsABreak(row, nowMs);
        var collapsedMark = rowHasBrokenStamps || projectedRow
          ? null
          : timelineCollapsedRowMark(
              segmentsById[row.id] || [],
              timelineViewState.windowStartMs,
              timelineViewState.windowEndMs,
              plotWidthPx
            );
        if (collapsedMark) {
          drawSegment(rowGroup, rowTopY, collapsedMark.startMs, collapsedMark.endMs,
            "timeline-segment timeline-segment-whole-row");
          continue;
        }
        // Same reasoning as the work segment below: a span with no honest end has
        // no width to draw, so it becomes a break marker at the wait's own start
        // instant rather than a bar drawn left-to-right by drawSegment's
        // min/max sort. Two causes reach here — a reversed span, and a REQ that
        // stopped with no resolvable end instant — and both draw the same mark,
        // because in both cases what the reader needs to know is "the bookkeeping
        // for this row is broken", not which way.
        if (row.waitMinutes < 0 || isNaN(waitEndMs)) {
          makeTimelineSvgNode(rowGroup, "rect", {
            x: (xOfEpoch(createdMs) - 3).toFixed(1),
            y: rowTopY + 2,
            width: 6,
            height: TIMELINE_ROW_HEIGHT - 4,
            class: "timeline-segment-broken"
          });
        } else {
          drawSegment(rowGroup, rowTopY, createdMs, waitEndMs,
            "timeline-segment timeline-segment-wait" + (row.waitOpen ? " is-open" : ""));
        }

        if (projectedRow) {
          drawProjectedSegment(rowGroup, rowTopY, projectedRow);
        }
        if (row.hasWork) {
          var workStartMs = claimedMs;
          var workEndMs = timelineWorkEndMs(row, nowMs, completedMs);
          // A span with no honest end has no width to draw, so it is drawn as a
          // break marker at the claim instant rather than as a bar running
          // backwards or clamped forwards to nothing.
          if (row.workMinutes < 0 || isNaN(workEndMs)) {
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
      drawQueueEndRule();
      // Last, so the present paints over the forecast rather than under it.
      drawNowRule();

      // And put focus back on the row it was on, or on the chart itself when the
      // window no longer draws that row. Only when focus was INSIDE the rows before
      // the rebuild: a reader who has tabbed away must not be dragged back.
      if (focusedRowIdBeforeRebuild) {
        var rebuiltRow = rowsSvg.querySelector
          ? rowsSvg.querySelector('[data-detail-id="' + focusedRowIdBeforeRebuild + '"]')
          : null;
        var focusTarget = rebuiltRow || scrollHost;
        if (focusTarget && typeof focusTarget.focus === "function") {
          focusTarget.focus();
        }
      }
    }

    // Seven lines per render — one per axis tick — and never one per row, so the
    // virtualization keeps its cost. Read from the same instants renderAxis
    // labels, so a gridline cannot mean a slightly different time than the tick
    // sitting above it.
    function drawGridlines() {
      var tickInstants = timelineAxisTickInstants(
        timelineViewState.windowStartMs, timelineViewState.windowEndMs);
      var rulesHeight = timelineDisplay.height;
      for (var tickIndex = 0; tickIndex < tickInstants.length; tickIndex++) {
        var tickX = xOfEpoch(tickInstants[tickIndex]);
        makeTimelineSvgNode(rowsSvg, "line", {
          x1: tickX.toFixed(1),
          y1: 0,
          x2: tickX.toFixed(1),
          y2: rulesHeight,
          class: "timeline-gridline"
        });
      }
    }

    // The instant the forecast paragraph names, or null when there is nothing
    // honest to draw: a declined projection has no instant, and an instant
    // outside the window would have to be drawn at an edge it does not sit on.
    function queueEndRuleInstant() {
      if (!projection.confident || isNaN(queueEndMs)) {
        return null;
      }
      if (
        queueEndMs < timelineViewState.windowStartMs ||
        queueEndMs > timelineViewState.windowEndMs
      ) {
        return null;
      }
      return queueEndMs;
    }

    function drawQueueEndRule() {
      var instantMs = queueEndRuleInstant();
      if (instantMs === null) {
        return;
      }
      var ruleX = xOfEpoch(instantMs);
      makeTimelineSvgNode(rowsSvg, "line", {
        x1: ruleX.toFixed(1),
        y1: 0,
        x2: ruleX.toFixed(1),
        y2: timelineDisplay.height,
        class: "timeline-queue-end-rule"
      });
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
        y2: timelineDisplay.height,
        class: "timeline-now-rule"
      });
    }

    function renderAxis() {
      axisSvg.textContent = "";
      var axisTicks = timelineAxisTicks(
        timelineViewState.windowStartMs, timelineViewState.windowEndMs);
      var tickInstants = axisTicks.instants;
      var plotLeftEdgeX = TIMELINE_LABEL_WIDTH;
      var plotSpanPx = plotWidth();
      for (var tickIndex = 0; tickIndex < tickInstants.length; tickIndex++) {
        var tickMs = tickInstants[tickIndex];
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
            // Anchored by POSITION, not by index. The first and last ticks used to
            // be the window's own edges, so anchoring them start/end was the same
            // thing; a calendar tick can now fall anywhere, including one pixel
            // from an edge, and an index-keyed anchor would clip it there.
            "text-anchor":
              tickX < plotLeftEdgeX + plotSpanPx * 0.06
                ? "start"
                : tickX > plotLeftEdgeX + plotSpanPx * 0.94
                  ? "end"
                  : "middle"
          },
          timelineFormatAxisTick(
            tickMs,
            axisTicks.gapMs,
            timelineViewState.windowStartMs,
            timelineViewState.windowEndMs
          )
        );
      }
      // The queue-end rule's LABEL lives here rather than in the rows SVG: the
      // rows SVG scrolls, and a caption that scrolls away from the rule it names
      // is worse than none. The rule itself stays in the rows SVG, for the same
      // reason drawNowRule does — that is what guarantees it uses the bars' own
      // x scale.
      var queueEndInstantMs = queueEndRuleInstant();
      if (queueEndInstantMs !== null) {
        var queueEndX = xOfEpoch(queueEndInstantMs);
        makeTimelineSvgNode(axisSvg, "line", {
          x1: queueEndX.toFixed(1),
          y1: 0,
          x2: queueEndX.toFixed(1),
          y2: TIMELINE_AXIS_HEIGHT,
          class: "timeline-queue-end-line"
        });
        // Anchored away from whichever edge it is near, so the caption stays on
        // the canvas instead of being clipped by it. Wider margins than the tick
        // labels above use, because "queue empty" is eleven characters and a tick
        // label is six.
        makeTimelineSvgNode(
          axisSvg,
          "text",
          {
            x: queueEndX.toFixed(1),
            y: 8,
            class: "timeline-queue-end-label",
            "text-anchor":
              queueEndX < plotLeftEdgeX + plotSpanPx * 0.15
                ? "start"
                : queueEndX > plotLeftEdgeX + plotSpanPx * 0.85
                  ? "end"
                  : "middle"
          },
          "queue empty"
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

    // Which chip is lit, and what the state span beside it says. Called from
    // renderAll, so every path that moves the window — chips, arrows, keys, wheel,
    // drag, typed dates — refreshes it and none of them can leave a stale claim on
    // screen.
    //
    // BOTH OUTPUTS ARE DERIVED BY RE-ASKING EACH BUTTON IN THE DOM what window it
    // would produce, using the same function the click uses. The control set is
    // then declared once, in template.html, and "press a chip and that chip lights"
    // is true by construction rather than by a second list in here agreeing with
    // the first one (REQ-235: derive state instead of storing it).
    //
    // The state text carries the one thing the lit chip cannot say. On a board
    // younger than the span asked for — every board in its first three months —
    // Last 90 days gives whatever there is, and lighting the chip unqualified
    // would claim ninety days the board does not have. "part of" is decided by
    // whether the settled window COVERS what the chip asked for, not by comparing
    // spans: the settle applies a one-hour floor and a bound-span ceiling, so a
    // span comparison misreports both a very young board and the whole range.
    //
    // An UNCLIPPED match wins over a clipped one when several chips land on the
    // same window, which is the ordinary case on a young board: All days and Last
    // 90 days are the same three days there, and "all days" is the honest name for
    // them.
    function renderTrailingWindowControls() {
      var firstMatchingChip = null;
      var firstUnclippedChip = null;
      document.querySelectorAll("#view-timeline [data-timeline-period]").forEach(function (chip) {
        var trailingWindowValue = chip.getAttribute("data-timeline-period");
        var candidate = timelineTrailingWindow(trailingWindowValue, nowMs, boundStartMs, boundEndMs);
        if (candidate.windowStartMs !== timelineViewState.windowStartMs ||
          candidate.windowEndMs !== timelineViewState.windowEndMs) {
          return;
        }
        var askedDayCount = Number(trailingWindowValue);
        var match = {
          value: trailingWindowValue,
          label: (chip.textContent || "").trim().toLowerCase(),
          isClipped: isFinite(askedDayCount) && askedDayCount > 0 &&
            (candidate.windowStartMs > nowMs - askedDayCount * TIMELINE_DAY_MS ||
              candidate.windowEndMs < nowMs)
        };
        if (!firstMatchingChip) {
          firstMatchingChip = match;
        }
        if (!match.isClipped && !firstUnclippedChip) {
          firstUnclippedChip = match;
        }
      });
      var matchedChip = firstUnclippedChip || firstMatchingChip;
      setActiveButton("#view-timeline .timeline-periods", "data-timeline-period",
        matchedChip ? matchedChip.value : "");
      var trailingWindowStateNode = document.getElementById("timeline-period-state");
      if (trailingWindowStateNode) {
        trailingWindowStateNode.textContent = !matchedChip
          ? "custom span"
          : matchedChip.isClipped
            ? "part of " + matchedChip.label
            : matchedChip.label;
      }
    }

    // A control that cannot move the window SAYS SO. Silence is what made the Now
    // landing state read as a broken board: pressing +, ctrl-scrolling and holding
    // the + key all did nothing at all, with no disabled state and no message.
    //
    // Every verdict here is derived by asking the shared model what the press would
    // produce and comparing it with the window already on screen — never by
    // restating the model's floor, ceiling or clamp, which is how a second opinion
    // about them would start.
    function renderControlAvailability() {
      function setAvailability(buttonId, wouldMove) {
        var button = document.getElementById(buttonId);
        if (button) {
          button.disabled = !wouldMove;
        }
      }
      function movesTheWindow(candidate) {
        return candidate.windowStartMs !== timelineViewState.windowStartMs ||
          candidate.windowEndMs !== timelineViewState.windowEndMs;
      }
      var zoomIn = timelineZoomedWindow(
        timelineViewState.windowStartMs, timelineViewState.windowEndMs,
        TIMELINE_ZOOM_STEP, 0.5, boundStartMs, boundEndMs);
      var zoomOut = timelineZoomedWindow(
        timelineViewState.windowStartMs, timelineViewState.windowEndMs,
        1 / TIMELINE_ZOOM_STEP, 0.5, boundStartMs, boundEndMs);
      setAvailability("timeline-zoom-in", movesTheWindow(zoomIn));
      setAvailability("timeline-zoom-out", movesTheWindow(zoomOut));
      setAvailability("timeline-zoom-fit", movesTheWindow(timelineFitWindow()));
      [-1, 1].forEach(function (stepCount) {
        var stepped = steppedWindowFor(stepCount);
        setAvailability(
          stepCount < 0 ? "timeline-period-prev" : "timeline-period-next",
          movesTheWindow(stepped) && !stepLandsOffTheData(stepped)
        );
      });
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

    // Which field, if any, the reader is PART-WAY THROUGH typing into. Set on
    // `input`, cleared on the `change` that commits it and on blur.
    //
    // Skipping every FOCUSED field was the first rule and the wrong one: a reader
    // who clicks into a field and then zooms the chart leaves it showing a window
    // the chart is no longer drawn at, and committing it later silently undoes the
    // zoom. Comparing the field's value against the last value this code wrote was
    // the second, and it got both halves wrong:
    //
    //   after a COMMIT the field still holds the reader's own text, which differs
    //   from what we last wrote, so the write-back was skipped and a clamped date
    //   stayed on screen indefinitely — the chart drawn at one window, the field
    //   naming another;
    //
    //   and when the reader CLEARED a field, applyTypedRange removed the attribute
    //   to turn the guard OFF, which turned it on ("" !== null), leaving the field
    //   permanently blank in the one branch whose comment says "Restore it
    //   unconditionally".
    //
    // An event is the honest signal, because editing is a thing the reader does
    // rather than a state of the value. Only one field can be focused, so one
    // variable answers it.
    var rangeFieldBeingEdited = null;

    function syncRangeField(field, text) {
      if (!field || field === rangeFieldBeingEdited) {
        return;
      }
      field.value = text;
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
        // window that does not exist. The restore is unconditional because the
        // change event that got us here has already cleared the editing flag.
        renderRangeControls();
        return;
      }
      timelineViewState.windowStartMs = typedWindow.windowStartMs;
      timelineViewState.windowEndMs = typedWindow.windowEndMs;
      renderAll();
    }

    [rangeStartField, rangeEndField].forEach(function (field) {
      if (!field) {
        return;
      }
      addTimelineListener(field, "input", function () {
        rangeFieldBeingEdited = field;
      });
      addTimelineListener(field, "change", function () {
        // Cleared FIRST: a commit ends the edit, and everything downstream —
        // including the restore of a cleared field — depends on the write-back
        // being allowed again.
        rangeFieldBeingEdited = null;
        applyTypedRange(field);
      });
      // A reader who starts typing and then clicks away without committing is no
      // longer editing either, and the next render owes them the real window.
      addTimelineListener(field, "blur", function () {
        if (rangeFieldBeingEdited === field) {
          rangeFieldBeingEdited = null;
          renderRangeControls();
        }
      });
    });

    // Every path that moves the window ends here, which is what makes the row
    // set, the counts, the scroll extent and the table describe the same window.
    // The order is decided rather than incidental: rows first, because everything
    // below reads them.
    function renderAll() {
      // Before anything reads it: the container may have been resized, and every
      // x in this render has to come from one measurement of one width.
      invalidatePlotWidth();
      // Nothing is visible and nothing is measurable, so a render here would only
      // record wrong numbers — a 120px plot and an eight-row viewport — and leave
      // them on screen when the view comes back. The observer below re-enters as
      // soon as the host has a box.
      if (!plotIsMeasurable()) {
        return;
      }
      refreshWindowRows();
      renderSummary();
      renderForecastIfSubsetChanged();
      renderAxis();
      renderVisibleRows();
      renderTrailingWindowControls();
      renderControlAvailability();
      renderRangeControls();
      // The hover readout describes ONE row, and a window move can take that row
      // off the chart. Every other piece of prose this view owns is refreshed
      // above; this one used to be written only by the pointer, so it went on
      // announcing a REQ that was no longer drawn anywhere.
      clearHoverReadout();
      markTimelineTableStale();
    }

    // The plot's box can change without any window moving and without the browser
    // resizing: the detail drawer takes a grid column, the view is shown or
    // hidden, a container's padding changes. Both delivery paths below call this
    // one positive live-width-change condition instead of enumerating the causes.
    //
    // Guarded because the floor this project designs for includes hosts with no
    // ResizeObserver: the window resize listener below still covers the ordinary
    // case there, which is exactly what it covered before.
    function renderIfPlotWidthChanged() {
      var hostWidth = liveHostWidth();
      if (hostWidth > 0 && hostWidth !== renderedHostWidthPx) {
        renderAll();
      }
    }

    if (typeof ResizeObserver === "function") {
      var plotResizeObserver = new ResizeObserver(function () {
        renderIfPlotWidthChanged();
      });
      plotResizeObserver.observe(scrollHost);
      timelineListenerTeardowns.push(function () {
        plotResizeObserver.disconnect();
      });
    }

    // Chromium's DOM-dump lifecycle can settle layout while compositor callbacks
    // remain parked: the host has its new width, but ResizeObserver is not delivered.
    // This teardown-owned 50ms timer is a second delivery path for the SAME shared
    // condition, not a second list of layout callers. Suppressing or inverting that
    // condition returns the drawer probe to RED; removing only ResizeObserver does
    // not, because this fallback intentionally retains the invariant.
    if (typeof window.setInterval === "function") {
      var plotWidthCheckTimer = window.setInterval(renderIfPlotWidthChanged, 50);
      timelineListenerTeardowns.push(function () {
        window.clearInterval(plotWidthCheckTimer);
      });
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
    // Bring a row inside the scroll viewport. The rows are virtualized, so a roving move
    // to a row outside the rendered range has to scroll BEFORE the rebuild, or the render
    // has no node to focus.
    function displayItemForRowIndex(rowIndex) {
      for (var displayIndex = 0; displayIndex < timelineDisplay.items.length; displayIndex++) {
        var displayItem = timelineDisplay.items[displayIndex];
        if (displayItem.kind === "request" && displayItem.rowIndex === rowIndex) {
          return displayItem;
        }
      }
      return null;
    }

    function scrollTimelineRowIntoView(rowIndex) {
      var displayItem = displayItemForRowIndex(rowIndex);
      if (!displayItem) {
        return;
      }
      var rowTopPx = displayItem.topPx;
      var rowBottomPx = rowTopPx + displayItem.height;
      if (rowTopPx < scrollHost.scrollTop) {
        scrollHost.scrollTop = rowTopPx;
      } else if (rowBottomPx > scrollHost.scrollTop + scrollHost.clientHeight) {
        scrollHost.scrollTop = rowBottomPx - scrollHost.clientHeight;
      }
    }

    // Move the roving stop, and the focus with it. The order matters: scroll, then
    // rebuild SYNCHRONOUSLY, then focus. The scroll write also schedules an asynchronous
    // scroll event whose handler is renderVisibleRows, and that rebuild restores whatever
    // row holds focus when it runs — so focusing before it lands is what survives, and
    // focusing before the synchronous rebuild would not.
    function moveTimelineRowFocus(rowDelta) {
      var currentDisplayIndex = -1;
      for (var displayIndex = 0; displayIndex < timelineDisplay.items.length; displayIndex++) {
        var displayItem = timelineDisplay.items[displayIndex];
        if (
          displayItem.kind === "request" &&
          displayItem.rowIndex === timelineViewState.rovingRowIndex
        ) {
          currentDisplayIndex = displayIndex;
          break;
        }
      }
      if (currentDisplayIndex < 0) {
        return;
      }
      var nextRowIndex = timelineViewState.rovingRowIndex;
      for (
        var nextDisplayIndex = currentDisplayIndex + rowDelta;
        nextDisplayIndex >= 0 && nextDisplayIndex < timelineDisplay.items.length;
        nextDisplayIndex += rowDelta
      ) {
        if (timelineDisplay.items[nextDisplayIndex].kind === "request") {
          nextRowIndex = timelineDisplay.items[nextDisplayIndex].rowIndex;
          break;
        }
      }
      if (nextRowIndex === timelineViewState.rovingRowIndex) {
        return;
      }
      timelineViewState.rovingRowIndex = nextRowIndex;
      scrollTimelineRowIntoView(nextRowIndex);
      renderVisibleRows();
      var nextRowGroup = rowsSvg.querySelector
        ? rowsSvg.querySelector('[data-row-index="' + nextRowIndex + '"]')
        : null;
      if (nextRowGroup && typeof nextRowGroup.focus === "function") {
        nextRowGroup.focus();
      }
    }

    // The roving stop follows focus however focus arrived — a Tab into the list, a click,
    // a rebuild's restore — so the next arrow press moves from where the reader actually
    // is rather than from where the last arrow press left the index.
    //
    // Attributes only, no rebuild: focusin fires on every entry into the list, and a
    // render per focus change would cost a full row sweep for two attribute writes.
    addTimelineListener(scrollHost, "focusin", function (focusEvent) {
      var focusedRowGroup = focusEvent.target && focusEvent.target.closest
        ? focusEvent.target.closest("[data-row-index]")
        : null;
      if (!focusedRowGroup) {
        return;
      }
      var focusedRowIndex = Number(focusedRowGroup.getAttribute("data-row-index"));
      if (!isFinite(focusedRowIndex) || focusedRowIndex === timelineViewState.rovingRowIndex) {
        return;
      }
      timelineViewState.rovingRowIndex = focusedRowIndex;
      var previousStop = rowsSvg.querySelector
        ? rowsSvg.querySelector('[data-row-index][tabindex="0"]')
        : null;
      if (previousStop && previousStop !== focusedRowGroup) {
        previousStop.setAttribute("tabindex", "-1");
      }
      focusedRowGroup.setAttribute("tabindex", "0");
    });

    addTimelineListener(scrollHost, "keydown", function (keyEvent) {
      var activation = timelineKeyboardActivationTarget(keyEvent);
      if (activation) {
        keyEvent.preventDefault();
        openDetail(activation.detailKind, activation.detailId);
        return;
      }

      // Vertical keys rove the row list, but only from a row: with focus on the chart
      // itself they stay the browser's own scroll, which is the only way to move through
      // the queue without entering the list.
      if (keyEvent.key === "ArrowDown" || keyEvent.key === "ArrowUp") {
        var rovingFrom = keyEvent.target && keyEvent.target.closest
          ? keyEvent.target.closest("[data-row-index]")
          : null;
        if (!rovingFrom) {
          return;
        }
        keyEvent.preventDefault();
        moveTimelineRowFocus(keyEvent.key === "ArrowDown" ? 1 : -1);
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

      // No focus restore here. renderVisibleRows owns it, because it owns the
      // rebuild — and it is reached by paths this handler cannot see. Restoring it
      // here as well was worse than redundant: the scroll event that refreshWindowRows
      // triggers is ASYNCHRONOUS, so this restore ran first and the rebuild it
      // scheduled then wiped the node it had focused. One arrow press worked and
      // every press after it was dead.
      renderAll();
    });

    // The readout is role="status" aria-live="polite", so every write to it is an
    // ANNOUNCEMENT. Writing on every mousemove queued one per event — a single
    // sweep across the rows produced dozens for one gesture — so the last id
    // announced is remembered and an unchanged row is not re-announced.
    addTimelineListener(scrollHost, "mousemove", function (moveEvent) {
      if (!readoutNode) {
        return;
      }
      var rowGroup = moveEvent.target.closest ? moveEvent.target.closest("[data-row-index]") : null;
      if (!rowGroup) {
        clearHoverReadout();
        return;
      }
      var row = rows[Number(rowGroup.getAttribute("data-row-index"))];
      if (!row) {
        clearHoverReadout();
        return;
      }
      if (row.id === announcedHoverRowId) {
        return;
      }
      announcedHoverRowId = row.id;
      readoutNode.textContent = timelineRowDescription(row, requestsById[row.id] || {});
    });
    addTimelineListener(scrollHost, "mouseleave", clearHoverReadout);

    // Zoom is modifier-gated so a plain wheel keeps scrolling the rows, which is
    // the motion a 560-row list needs most.
    //
    // Bound to the AXIS as well as the plot. The hint says holding Ctrl and
    // scrolling zooms the time axis, and over the axis strip itself it did nothing —
    // the one place a reader aiming at "the time axis" is most likely to point. The
    // anchor is measured against the SCROLL HOST either way, because both draw
    // against plotWidth() on that host and the two are horizontally aligned; taking
    // the rect from whichever element the wheel landed on would silently introduce a
    // second x scale.
    function handleTimelineWheel(wheelEvent) {
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
    }
    addTimelineListener(scrollHost, "wheel", handleTimelineWheel, { passive: false });
    addTimelineListener(axisHost, "wheel", handleTimelineWheel, { passive: false });

    // A drag issues at most ONE render per frame. Every pointermove used to run a
    // full renderAll — axis plus every visible row — and a trackpad delivers
    // moves faster than the compositor draws them, so most of that work was
    // thrown away before it was seen.
    function requestFrameRender() {
      if (timelineFrameRender !== null) {
        return;
      }
      if (!window.requestAnimationFrame) {
        renderAll();
        return;
      }
      timelineFrameRender = window.requestAnimationFrame(function () {
        timelineFrameRender = null;
        renderAll();
      });
    }

    var panState = null;
    addTimelineListener(scrollHost, "pointerdown", function (downEvent) {
      if (downEvent.button !== 0) {
        return;
      }
      // Armed, not engaged. Nothing moves and nothing re-renders until the
      // pointer has travelled TIMELINE_PAN_THRESHOLD_PX — including the grab
      // cursor, which used to change on the press and told the reader they were
      // dragging when they were clicking.
      panState = {
        pointerX: downEvent.clientX,
        pointerId: downEvent.pointerId,
        windowStartMs: timelineViewState.windowStartMs,
        engaged: false
      };
    });

    // CAPTURE THE POINTER once the pan ENGAGES, so the release is guaranteed to arrive
    // here — and not before, because capture also retargets the synthesized click.
    //
    // Capture routes an outside-host pointerup back to this host, so teardown
    // does not depend on pointerleave being delivered at the boundary. The
    // lostpointercapture handler also ends the pan if the engine revokes capture.
    //
    // Capturing on pointerdown instead cost every mouse click in the chart: capture
    // retargets subsequent pointer events AND the synthesized click to the capturing
    // element, so the delegated handler in board-controls.js found no
    // [data-detail-kind] ancestor and the detail drawer never opened for any press,
    // drag or not. Taking capture at the engage instead makes "a drag is not a click"
    // a property of the capture rather than an accident, and leaves a press that never
    // engages an ordinary, untouched click.
    function capturePanPointer() {
      if (typeof scrollHost.setPointerCapture !== "function" || panState.pointerId === undefined) {
        return;
      }
      try {
        scrollHost.setPointerCapture(panState.pointerId);
      } catch (captureError) {
        // A pointer the engine will not let us capture is not a reason to refuse the
        // drag; the release events below are still the ordinary path, and the teardown
        // asks hasPointerCapture rather than assuming this call succeeded.
      }
    }

    addTimelineListener(scrollHost, "pointermove", function (moveEvent) {
      if (!panState) {
        return;
      }
      var wasEngaged = panState.engaged;
      panState.engaged = timelinePanEngaged(panState.engaged, panState.pointerX, moveEvent.clientX);
      if (!panState.engaged) {
        return;
      }
      if (!wasEngaged) {
        capturePanPointer();
      }
      scrollHost.classList.add("is-panning");
      var windowSpanMs = timelineViewState.windowEndMs - timelineViewState.windowStartMs;
      // Measured from the PRESS point, not from wherever the threshold tripped:
      // anchoring at the trip point would leave the chart four pixels behind the
      // pointer for the rest of the drag.
      var shiftMs = ((panState.pointerX - moveEvent.clientX) / plotWidth()) * windowSpanMs;
      var nextStartMs = Math.min(
        Math.max(panState.windowStartMs + shiftMs, boundStartMs),
        boundEndMs - windowSpanMs
      );
      timelineViewState.windowStartMs = nextStartMs;
      timelineViewState.windowEndMs = nextStartMs + windowSpanMs;
      requestFrameRender();
    });
    ["pointerup", "pointercancel", "pointerleave", "lostpointercapture"].forEach(function (eventName) {
      addTimelineListener(scrollHost, eventName, function () {
        // A drag that ends between frames still has to land on the window it
        // reached, so the pending frame runs rather than being dropped.
        if (panState && panState.engaged && timelineFrameRender !== null) {
          releaseTimelineFrameRender();
          renderAll();
        }
        if (
          panState &&
          panState.pointerId !== undefined &&
          typeof scrollHost.releasePointerCapture === "function" &&
          typeof scrollHost.hasPointerCapture === "function" &&
          scrollHost.hasPointerCapture(panState.pointerId)
        ) {
          scrollHost.releasePointerCapture(panState.pointerId);
        }
        panState = null;
        scrollHost.classList.remove("is-panning");
      });
    });

    // The extent of everything this render draws: every segment of every
    // filter-matched row, plus the projection's own bars. Computed from the segment
    // list the window-scoping already built, so it is the same set of marks the
    // chart puts on screen rather than a second opinion about them.
    var drawnExtent = (function () {
      var earliestMs = Infinity;
      var latestMs = -Infinity;
      filterMatchedSegments.forEach(function (segments) {
        segments.forEach(function (segment) {
          if (!isFinite(segment.startMs) || !isFinite(segment.endMs)) {
            return;
          }
          earliestMs = Math.min(earliestMs, segment.startMs);
          latestMs = Math.max(latestMs, segment.endMs);
        });
      });
      if (!isFinite(earliestMs) || !isFinite(latestMs)) {
        return null;
      }
      return { earliestMs: earliestMs, latestMs: latestMs };
    })();

    // What Fit all lands on: the drawn extent with the same breathing room the
    // payload bounds get, settled through the shared model so it cannot acquire a
    // floor or a clamp of its own.
    function timelineFitWindow() {
      if (!drawnExtent) {
        return { windowStartMs: boundStartMs, windowEndMs: boundEndMs };
      }
      var breathingRoomMs = Math.max(
        (drawnExtent.latestMs - drawnExtent.earliestMs) * 0.02, 60 * 1000);
      return timelineZoomedWindow(
        drawnExtent.earliestMs - breathingRoomMs,
        drawnExtent.latestMs + breathingRoomMs,
        1, 0, boundStartMs, boundEndMs);
    }

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
    // Fit all fits WHAT IS ON SCREEN, which under a filter is the filtered set —
    // not the payload's whole range. The clamp bounds stay the payload's, so the
    // reader can still pan and zoom outside the filtered extent; it is the button
    // that answers "show me everything I am looking at", and filtered to one domain
    // it used to leave most of the plot blank.
    wireToolbarButton("timeline-zoom-fit", function () {
      var fitWindow = timelineFitWindow();
      timelineViewState.windowStartMs = fitWindow.windowStartMs;
      timelineViewState.windowEndMs = fitWindow.windowEndMs;
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
          var openDisplayItem = displayItemForRowIndex(openRowIndex);
          scrollHost.scrollTop = openDisplayItem ? openDisplayItem.topPx : 0;
          renderVisibleRows();
        }
      };
    }

    // The window controls: five trailing windows and a step either way. A chip
    // picks a WIDTH ending at now; the arrows then walk that width back and forth
    // through the archive. Two rules where there used to be three, because a
    // window anchored at now has no calendar grid for a step to follow.
    function applyTrailingWindow(trailingWindowValue) {
      var trailingWindow = timelineTrailingWindow(
        trailingWindowValue, nowMs, boundStartMs, boundEndMs);
      timelineViewState.windowStartMs = trailingWindow.windowStartMs;
      timelineViewState.windowEndMs = trailingWindow.windowEndMs;
    }
    // Whether a step would land entirely PAST everything drawn — off the end of the
    // data in one direction or the other.
    //
    // What this exists for is the cosmetic bound padding: that stretch is there so a
    // bar at the range edge is not flush against the frame, and on a three-month
    // board it is nearly two days wide, so one press of › from a window sitting at
    // the right of the data landed on an empty chart with the arrow left alive.
    //
    // It deliberately does NOT refuse a step into a GAP between drawn things. This
    // repo's own queue has a seventeen-day hole in June; stepping through it is the
    // reader exploring, and the empty-window message already tells them what they
    // are looking at. Only a step off the END of everything is refused, which is why
    // the test is against the extent rather than against occupancy.
    function stepLandsOffTheData(steppedWindow) {
      if (!drawnExtent) {
        return false;
      }
      return steppedWindow.windowStartMs > drawnExtent.latestMs ||
        steppedWindow.windowEndMs < drawnExtent.earliestMs;
    }

    // One screenful, in the direction pressed. Whatever the window is — a chip's
    // trailing span, a free zoom, a drag, a typed pair — the arrows move it by its
    // own width and never resize it.
    function steppedWindowFor(stepCount) {
      return timelinePannedWindow(
        timelineViewState.windowStartMs,
        timelineViewState.windowEndMs,
        stepCount,
        boundStartMs,
        boundEndMs
      );
    }

    // No second guard here. renderControlAvailability below already disables the
    // arrow whose step would land off the data, and a disabled button fires no
    // click — so a check repeated at this call site would be unreachable, and a
    // guard no mutation can break is dead code rather than safety.
    function applyTrailingWindowStep(stepCount) {
      var steppedWindow = steppedWindowFor(stepCount);
      timelineViewState.windowStartMs = steppedWindow.windowStartMs;
      timelineViewState.windowEndMs = steppedWindow.windowEndMs;
    }
    document.querySelectorAll("#view-timeline [data-timeline-period]").forEach(function (chip) {
      chip.onclick = function () {
        applyTrailingWindow(chip.getAttribute("data-timeline-period"));
        renderAll();
      };
    });
    wireToolbarButton("timeline-period-prev", function () {
      applyTrailingWindowStep(-1);
    });
    wireToolbarButton("timeline-period-next", function () {
      applyTrailingWindowStep(1);
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
