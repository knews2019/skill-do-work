  // ---- small DOM helpers --------------------------------------------------

  function createElement(tagName, className, textContent) {
    var node = document.createElement(tagName);
    if (className) {
      node.className = className;
    }
    if (textContent != null) {
      node.textContent = textContent;
    }
    return node;
  }

  var columnDayFormatter = new Intl.DateTimeFormat("en-US", {
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
    timeZone: "UTC",
    hour12: false
  });
  var calendarWeekdayFormatter = new Intl.DateTimeFormat("en-US", {
    weekday: "long",
    timeZone: "UTC"
  });
  var calendarDateFormatter = new Intl.DateTimeFormat("en-US", {
    year: "numeric",
    month: "short",
    day: "2-digit",
    timeZone: "UTC"
  });

  function formatShortInstant(isoText) {
    var ms = Date.parse(isoText);
    if (isNaN(ms)) {
      return "";
    }
    return columnDayFormatter.format(new Date(ms)) + " UTC";
  }

  // ---- relative time ("6min ago") ----------------------------------------
  // Every visible instant carries a relative companion. Visible nodes get
  // data-instant-ms and a shared 1s ticker keeps them fresh in a tab left
  // open; title tooltips (which cannot host a ticking node) get a
  // render-time snapshot string instead.

  function formatRelativeTime(instantMs, nowMs) {
    var elapsedSeconds = Math.floor((nowMs - instantMs) / 1000);
    if (elapsedSeconds < 1) {
      // Also covers clock skew (instant slightly in the future).
      return "just now";
    }
    if (elapsedSeconds < 60) {
      return elapsedSeconds + "s ago";
    }
    var elapsedMinutes = Math.floor(elapsedSeconds / 60);
    if (elapsedMinutes < 60) {
      return elapsedMinutes + "min ago";
    }
    var elapsedHours = Math.floor(elapsedMinutes / 60);
    if (elapsedHours < 24) {
      return elapsedHours + "h ago";
    }
    var elapsedDays = Math.floor(elapsedHours / 24);
    return elapsedDays + "d ago";
  }

  function makeRelativeTimeNode(isoText) {
    var instantMs = Date.parse(isoText);
    if (isNaN(instantMs)) {
      return null;
    }
    var relativeNode = createElement("span", "relative-time", formatRelativeTime(instantMs, Date.now()));
    relativeNode.dataset.instantMs = String(instantMs);
    return relativeNode;
  }

  // Absolute instant plus its ticking relative label, as one inline node —
  // for visible text (cards, chips, drawer rows). Null when unparseable so
  // callers can fall back to the raw value.
  function makeInstantWithRelativeNode(isoText) {
    var absoluteText = formatShortInstant(isoText);
    if (!absoluteText) {
      return null;
    }
    var wrapperNode = createElement("span", "instant-with-relative");
    wrapperNode.appendChild(document.createTextNode(absoluteText + " "));
    wrapperNode.appendChild(makeRelativeTimeNode(isoText));
    return wrapperNode;
  }

  function formatShortInstantWithRelative(isoText) {
    var absoluteText = formatShortInstant(isoText);
    if (!absoluteText) {
      return "";
    }
    return absoluteText + " (" + formatRelativeTime(Date.parse(isoText), Date.now()) + ")";
  }

  // A stopwatch instant more than this far ahead of the viewer's clock is a bad
  // stamp (futureStampCauseText names the two ways that happens), not a young
  // claim — without the marker the negative clamp would render a dead-looking
  // "0s" on every tick for as long as the stamp stands. Mirrors model.go's
  // futureTimestampSkewAllowance — keep the two in lock-step.
  var futureInstantSkewAllowanceMs = 2 * 60 * 1000;
  var clockSkewMarkerText = "⚠ clock skew";
  // Why the stamp is wrong, said once for the whole client: the stopwatch
  // tooltip below and board-cards.js's "⚠ future stamp" badge both render this,
  // so the two can never tell a reader different stories. Byte-identical to
  // futureStampCauseClause in model.go — keep the two in lock-step; nothing in
  // the build compares them, so TestFutureStampCauseClauseMatchesTheShippedClient
  // is what catches an edit to one and not the other, and the two behavior
  // probes catch a rename that leaves the literal sitting unused in the file.
  // Deliberately one unbroken literal, not a wrapped concatenation: the lock-step
  // test matches it against the Go constant verbatim, and a split literal would
  // force that check to reassemble JavaScript before it could compare.
  var futureStampCauseText = "local wall-clock time written with a Z suffix, or a fabricated value (guessed or extrapolated instead of read from the clock)";
  var clockSkewExplanationText =
    "This timestamp is ahead of your clock by more than the 2-minute skew allowance — " +
    "likely " + futureStampCauseText + ". Fix the frontmatter with " +
    "the current UTC instant — YYYY-MM-DDTHH:MM:SSZ, per the Timestamp rule in " +
    "actions/work-reference.md. Until then the stopwatch cannot measure real elapsed time.";

  // Stopwatch-style elapsed duration ("47s", "4m 07s", "1h 23m", "3d 04h") for
  // a ticket sitting in a state — second-resolution below an hour because
  // claim spans are short, coarser above so it never reads as a wall of digits,
  // and a day tier so week-old queue waits don't render as "170h".
  function formatElapsedDuration(instantMs, nowMs) {
    if (instantMs - nowMs > futureInstantSkewAllowanceMs) {
      return clockSkewMarkerText;
    }
    var totalSeconds = Math.max(0, Math.floor((nowMs - instantMs) / 1000));
    var daysPart = Math.floor(totalSeconds / 86400);
    if (daysPart > 0) {
      var remainderHours = Math.floor((totalSeconds % 86400) / 3600);
      return daysPart + "d " + (remainderHours < 10 ? "0" : "") + remainderHours + "h";
    }
    var hoursPart = Math.floor(totalSeconds / 3600);
    var minutesPart = Math.floor((totalSeconds % 3600) / 60);
    var secondsPart = totalSeconds % 60;
    if (hoursPart > 0) {
      return hoursPart + "h " + (minutesPart < 10 ? "0" : "") + minutesPart + "m";
    }
    if (minutesPart > 0) {
      return minutesPart + "m " + (secondsPart < 10 ? "0" : "") + secondsPart + "s";
    }
    return secondsPart + "s";
  }

  function makeElapsedDurationNode(isoText) {
    var instantMs = Date.parse(isoText);
    if (isNaN(instantMs)) {
      return null;
    }
    var durationNode = createElement("span", "elapsed-duration", formatElapsedDuration(instantMs, Date.now()));
    durationNode.dataset.instantMs = String(instantMs);
    durationNode.dataset.tickFormat = "duration";
    syncClockSkewTitle(durationNode, durationNode.textContent);
    return durationNode;
  }

  // Absolute instant plus its ticking stopwatch, as one inline node — the
  // live-state variant of makeInstantWithRelativeNode. Null when unparseable
  // so callers can fall back to the plain instant+relative rendering.
  function makeInstantWithStopwatchNode(isoText) {
    var absoluteText = formatShortInstant(isoText);
    var durationNode = makeElapsedDurationNode(isoText);
    if (!absoluteText || !durationNode) {
      return null;
    }
    var wrapperNode = createElement("span", "instant-with-relative");
    wrapperNode.appendChild(document.createTextNode(absoluteText + " "));
    wrapperNode.appendChild(durationNode);
    return wrapperNode;
  }

  function isBlockedHoldStatus(statusText) {
    return statusText === "blocked" || (statusText || "").indexOf("blocked-") === 0;
  }

  // Which instant a card's live state timer counts from: every non-terminal
  // card shows how long it has been in its current state, measured from the
  // timestamp that state transition wrote. Dedicated stamps (claimed_at,
  // blocked_at) are authoritative for their states — file mtime
  // must never outrank them, since the pipeline appends sections to the file
  // all through a claim. States without their own transition instant (pending,
  // pending-answers, failed, unrecognized) resolve, in order:
  //   1. status_changed_at (stamped on flips with no dedicated stamp, e.g.
  //      clarify answered → pending) — verb "updated";
  //   2. the LATER of created_at and the file's mtime — verb "updated" when
  //      mtime wins (the file changed after capture: answers recorded, hand
  //      edits), "queued" when created_at wins (untouched since capture).
  // Terminal cards keep the static done line instead.

  // Capture writes the file moments AFTER stamping created_at, so mtime is
  // always a few seconds ahead on an untouched card. Only treat mtime as "the
  // file was edited later" when it beats created_at by more than this — a real
  // edit (clarify answers, a hand fix) lands minutes-to-days later.
  var mtimeMeaningfulDeltaMs = 5 * 60 * 1000;

  function stateTimerSpecFor(request) {
    if (request.status === "claimed" && request.claimedAt) {
      return { verbText: "claimed", instantIso: request.claimedAt };
    }
    if (isBlockedHoldStatus(request.status) && request.blockedAt) {
      return { verbText: "blocked", instantIso: request.blockedAt };
    }
    if (isTerminalResolvedStatus(request.status)) {
      return null;
    }
    if (request.statusChangedAt && !isNaN(Date.parse(request.statusChangedAt))) {
      return { verbText: "updated", instantIso: request.statusChangedAt };
    }
    var createdAtMs = request.createdAt ? Date.parse(request.createdAt) : NaN;
    var fileModifiedAtMs = request.fileModifiedAt ? Date.parse(request.fileModifiedAt) : NaN;
    if (!isNaN(fileModifiedAtMs) && (isNaN(createdAtMs) || fileModifiedAtMs - createdAtMs > mtimeMeaningfulDeltaMs)) {
      return { verbText: "updated", instantIso: request.fileModifiedAt };
    }
    if (request.createdAt) {
      return { verbText: "queued", instantIso: request.createdAt };
    }
    return null;
  }

  // Keeps a duration node's explanatory tooltip in step with the skew marker:
  // present while the marker shows, removed the moment the stamp is corrected
  // (or the wall clock catches up) and the stopwatch starts ticking for real.
  function syncClockSkewTitle(durationNode, labelText) {
    if (labelText === clockSkewMarkerText) {
      durationNode.title = clockSkewExplanationText;
    } else if (durationNode.title === clockSkewExplanationText) {
      durationNode.removeAttribute("title");
    }
  }

  // Takes its instant as an argument rather than reading the clock, so every
  // surface a single tick touches states the same "now" — see
  // refreshTickingSurfaces below, which is the one caller.
  function refreshRelativeTimeNodes(nowMs) {
    var relativeNodes = document.querySelectorAll("[data-instant-ms]");
    for (var nodeIndex = 0; nodeIndex < relativeNodes.length; nodeIndex++) {
      var relativeNode = relativeNodes[nodeIndex];
      var instantMs = Number(relativeNode.dataset.instantMs);
      var isDurationFormat = relativeNode.dataset.tickFormat === "duration";
      var nextLabel = isDurationFormat
        ? formatElapsedDuration(instantMs, nowMs)
        : formatRelativeTime(instantMs, nowMs);
      if (relativeNode.textContent !== nextLabel) {
        relativeNode.textContent = nextLabel;
      }
      if (isDurationFormat) {
        syncClockSkewTitle(relativeNode, nextLabel);
      }
    }
  }

  // The board's ONE tick. web/board.js's interval calls this and nothing else,
  // so this function is where "what moves every second" is stated.
  //
  // The order is load-bearing. refreshRelativeTimeNodes runs FIRST, with the
  // captured instant, because it drives every claim stopwatch, every relative
  // time, every state timer and the clock-skew tooltip on the page. If the
  // newer UR-summary pass below ever threw, the interval callback would die and
  // the whole board would freeze at its first paint — looking exactly like a
  // queue full of very young claims. Running the older pass first means such a
  // throw could not also cost those surfaces their current tick. The real
  // containment is that the rollup is total by narrowing (see
  // board-user-request-summary.js); this order is the belt beside it. Both are
  // asserted: a probe drives a tick whose summary pass throws on purpose and
  // checks the claim stopwatch still advanced, so swapping these two lines
  // fails the suite.
  function refreshTickingSurfaces() {
    var nowMs = Date.now();
    refreshRelativeTimeNodes(nowMs);
    refreshUserRequestSummaryNodes(nowMs);
  }

  // ---- dependency helpers -------------------------------------------------
  // Mirrors model.go's isTerminalResolvedStatus / isCompletedStatus. The board
  // never re-derives which dependencies are unmet — the Go side annotates that
  // (a dangling id counts as unmet, and `cancelled` never satisfies a
  // dependency) and ships it as request.unmetDependencies.

  // "Successful": the work shipped, with or without recorded issues. Composed
  // into the resolved set below exactly the way model.go composes them, so the
  // client never carries a fourth literal list of statuses.
  function isCompletedStatus(status) {
    return status === "completed" || status === "completed-with-issues";
  }

  function isTerminalResolvedStatus(status) {
    return isCompletedStatus(status) || status === "cancelled";
  }

  // The REQs still waiting on this one. A dependent that already resolved is not
  // "unblocked by" anything anymore, so it drops out of the count.
  function activeDependentIds(request) {
    return (request.dependents || []).filter(function (dependentId) {
      var dependent = requestsById[dependentId];
      return dependent && !isTerminalResolvedStatus(dependent.status);
    });
  }

  function describeRequestStatus(requestId) {
    var request = requestsById[requestId];
    return request && request.status ? request.status : "not in tree";
  }

  // ---- ticket mention resolution ------------------------------------------
  // Turning a REQ/UR id written in prose into the board record it names, and
  // describing that record. Shared rather than local to the drawer because the
  // clipboard surface resolves the same mentions, and two copies would drift.

  // A body mention like "REQ-031" may name a compound card id ("UR-002-REQ-031").
  // Index each card by its REQ segment so the short form still resolves; an
  // ambiguous segment (two cards sharing it) maps to null rather than guessing.
  function buildRequestIdByReqSegment() {
    var segmentIndex = {};
    Object.keys(requestsById).forEach(function (fullRequestId) {
      var segmentMatch = /REQ-\d+[a-z]?/i.exec(fullRequestId);
      if (!segmentMatch || segmentMatch[0] === fullRequestId) {
        return;
      }
      var segmentKey = segmentMatch[0].toUpperCase();
      segmentIndex[segmentKey] = Object.prototype.hasOwnProperty.call(segmentIndex, segmentKey)
        ? null // ambiguous — never guess
        : fullRequestId;
    });
    return segmentIndex;
  }

  var requestIdByReqSegment = buildRequestIdByReqSegment();

  function resolveTicketMention(mentionText) {
    if (Object.prototype.hasOwnProperty.call(requestsById, mentionText)) {
      return { kind: "req", id: mentionText };
    }
    if (Object.prototype.hasOwnProperty.call(userRequestsById, mentionText)) {
      return { kind: "ur", id: mentionText };
    }
    var segmentTargetId = requestIdByReqSegment[mentionText.toUpperCase()];
    if (segmentTargetId) {
      return { kind: "req", id: segmentTargetId };
    }
    return null;
  }

  // An ambiguous segment is not a dead reference: the board holds records that
  // match and refuses to pick one. Callers that flag unresolved ids as broken
  // must leave these alone, or the never-guess rule turns into a false alarm.
  function isAmbiguousTicketMention(mentionText) {
    var segmentKey = mentionText.toUpperCase();
    return (
      Object.prototype.hasOwnProperty.call(requestIdByReqSegment, segmentKey) &&
      requestIdByReqSegment[segmentKey] === null
    );
  }

  function ticketTitleFor(detailKind, detailId) {
    var ticketRecord = detailKind === "ur" ? userRequestsById[detailId] : requestsById[detailId];
    return (ticketRecord && ticketRecord.title) || "";
  }

  // A record can exist and still have no title, and the commonest case is not a
  // defect: linkRequestsToUserRequests SYNTHESIZES a UserRequestTicket whenever a
  // REQ names a UR whose input.md was not found, and a synthesized node carries
  // no Title by design. Falling back to the bare id there would reintroduce the
  // exact cryptic number this whole feature removes, on a supported board state.
  // So say what is known instead: why there is no title, never nothing.
  //
  // isFallback lets the caller render the substitute in a quieter voice — it is
  // a description of the record, not the record's own words.
  function describeTicketTitle(detailKind, detailId) {
    var fullTitle = ticketTitleFor(detailKind, detailId);
    if (fullTitle) {
      return { text: fullTitle, isFallback: false };
    }
    var ticketRecord = detailKind === "ur" ? userRequestsById[detailId] : requestsById[detailId];
    if (!ticketRecord) {
      return { text: "", isFallback: false }; // Not on the board at all — the missing-mention branch owns this.
    }
    if (detailKind === "ur" && ticketRecord.inputFilePresent === false) {
      // Same fact the drawer's own "input.md" meta row states for this record.
      return { text: "no input.md — synthesized from REQ pointers", isFallback: true };
    }
    return { text: "untitled", isFallback: true };
  }

  // Inline titles are cut here because a long one expanded mid-sentence swamps
  // the prose it sits in. Nothing is lost: the untruncated title rides in the
  // link's tooltip and in the drawer glossary.
  var inlineTicketTitleMaxLength = 60;

  function shortTicketTitle(fullTitle) {
    var trimmedTitle = (fullTitle || "").trim();
    if (trimmedTitle.length <= inlineTicketTitleMaxLength) {
      return trimmedTitle;
    }
    var cutTitle = trimmedTitle.slice(0, inlineTicketTitleMaxLength);
    var lastSpaceIndex = cutTitle.lastIndexOf(" ");
    if (lastSpaceIndex > 0) {
      // Cut on a word boundary; a title with no space in its first 60
      // characters is one long token, so the hard cut stands.
      cutTitle = cutTitle.slice(0, lastSpaceIndex);
    }
    return cutTitle.replace(/[\s,;:.\-]+$/, "") + "…";
  }
