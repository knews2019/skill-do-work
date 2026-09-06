  // ---- card construction --------------------------------------------------

  function truncateBadgeText(fullText, maxLength) {
    var limit = maxLength || 48;
    if (!fullText || fullText.length <= limit) {
      return fullText || "";
    }
    return fullText.slice(0, limit - 1).replace(/\s+$/, "") + "…";
  }

  // The "⚠ future stamp" badge's tooltip. A named function rather than an inline
  // concatenation so it can be executed on its own: the tooltip is the diagnosis
  // a reader is most likely to actually read, and the only way to assert what it
  // says is to build one. Shares futureStampCauseText with the stopwatch tooltip
  // (board-core.js) so the two cannot tell different stories.
  function futureStampTooltipText(futureTimestampFields) {
    return (
      "Future-dated timestamp(s): " +
      futureTimestampFields.join(", ") +
      " — later than the board's generation time (2min skew allowance). Likely " +
      futureStampCauseText +
      "; fix: rewrite with the current UTC " +
      "instant — YYYY-MM-DDTHH:MM:SSZ, per the Timestamp rule in actions/work-reference.md."
    );
  }

  function makeBadge(className, labelText, valueText, datasetName, datasetValue) {
    var badge = createElement("span", "badge " + className);
    if (labelText) {
      badge.appendChild(createElement("span", "badge-label", labelText));
    }
    badge.appendChild(document.createTextNode(valueText));
    if (datasetName) {
      badge.dataset[datasetName] = datasetValue;
    }
    return badge;
  }

  // Wall time from the EARLIEST recorded lifecycle stamp to completion, stated
  // on the done line beside the completion instant it already carries. The
  // origin is not necessarily claimed_at: a claim stamp rewritten after the work
  // happened would otherwise erase every phase that ran before it. This is the
  // calibration interval, not an assertion about active implementation time; the
  // detail drawer's observed phase breakdown shows where the recorded wall span
  // went, and the timeline's work bar is a DIFFERENT reading — it still splits
  // at claimed_at, because that bar is a statement about the claim itself.
  //
  // Both the number and its verdict arrive decided from Go (durations.go's
  // measureImplementationSpan). Go also supplies the completed pause-badge text,
  // so the read-time ceiling separating "that much work" from "an overnight
  // pause" keeps exactly one definition and this renderer never becomes a
  // second one.
  //
  // The node is a PLAIN span.elapsed-duration: it reuses the state timer's
  // vocabulary and styling so the card's two time lines read alike, but it
  // carries no data-instant-ms / data-tickFormat, which is what the 1s ticker
  // (board-core.js's refreshRelativeTimeNodes) selects on. A finished span must
  // not be rewritten every second, so makeElapsedDurationNode is deliberately not
  // used here.
  //
  // The verdict is branched on BEFORE the formatter runs. formatElapsedDuration
  // returns the clock-skew marker for a reversed pair, and that branch must stay
  // unreachable here: a reversed span prints the inline flag INSTEAD of a
  // magnitude, because a negative duration cannot be true and the card's
  // `anomaly` badge already carries the full diagnosis and the fix.
  function makeImplementationSpanNode(request) {
    if (!request.hasImplementationSpan) {
      return null;
    }
    if (request.implementationSpanReason === "reversed") {
      var reversedNode = createElement("span", "elapsed-duration implementation-span");
      var reversedFlag = createElement("span", "status-invalid-flag", "reversed stamps");
      reversedFlag.title =
        "completed_at is earlier than claimed_at, so this REQ's claim-to-completion wall span cannot be stated — " +
        "see the card's anomaly badge for which stamp to rewrite.";
      reversedNode.appendChild(reversedFlag);
      return reversedNode;
    }
    // formatElapsedDuration measures nowMs − instantMs, so a zero origin with the
    // Go-measured span as the "now" formats that span. The two frontmatter stamps
    // are deliberately NOT re-parsed here: Go reads a space-separated stamp as
    // UTC while Date.parse reads it as local time, which would silently shift the
    // span the card states away from the one the Durations view plots.
    var spanOriginMs = 0;
    var spanEndMs = Math.round(request.implementationSpanMinutes * 60000);
    var spanNode = createElement("span", "elapsed-duration implementation-span");
    spanNode.appendChild(
      createElement(
        "span",
        "implementation-span-value",
        "wall time " + formatElapsedDuration(spanOriginMs, spanEndMs)
      )
    );
    if (request.implementationSpanReason === "paused") {
      spanNode.appendChild(document.createTextNode(" "));
      var pausedFlag = createElement(
        "span",
        "status-invalid-flag",
        boardData.implementationSpanPausedBadgeText || "long span · assumed pause"
      );
      pausedFlag.title =
        "Duration-quality marker only: this claim-to-completion wall span is longer than the board's " +
        "single-session ceiling, so it is assumed to include a pause and excluded from duration medians. " +
        "The REQ remains completed.";
      spanNode.appendChild(pausedFlag);
    }
    return spanNode;
  }

  function makeRequestCard(requestId, options) {
    var request = requestsById[requestId];
    var card = createElement("button", "req-card");
    card.type = "button";
    card.dataset.detailKind = "req";
    card.dataset.detailId = requestId;

    if (!request) {
      // A dependency target outside the current tree — render the bare id.
      card.appendChild(createElement("span", "req-card-id", requestId));
      card.disabled = true;
      return card;
    }
    card.setAttribute("aria-label", requestId + ": " + (request.title || "untitled"));
    if (request.status) {
      card.dataset.status = request.status; // lets CSS apply the semantic rail and status-pill treatment
    }

    var top = createElement("div", "req-card-top");
    top.appendChild(createElement("span", "req-card-id", requestId));
    var status = createElement("span", "req-card-status");
    status.appendChild(createElement("span", "status-dot"));
    status.appendChild(document.createTextNode(request.status || "—"));
    if (request.statusUnrecognized) {
      card.className += " is-status-unrecognized";
      status.className += " is-status-unrecognized";
      status.appendChild(createElement("span", "status-invalid-flag", "invalid"));
      status.title =
        'Unrecognized status "' +
        (request.originalStatus || request.status) +
        "\" — edit the REQ's status: to a Schema Read Contract value or run do-work forensics";
    }
    top.appendChild(status);
    card.appendChild(top);

    card.appendChild(createElement("h3", "req-card-title", request.title || "untitled"));

    var badges = createElement("div", "req-card-badges");
    if (filterState.searchText && !searchMatchesRequest(request, requestId, filterState.searchText, true)) {
      var citedTicketId = citationMatchedTicketId(request, filterState.searchText);
      if (citedTicketId) {
        badges.appendChild(makeBadge("citation-match", null, "cites " + citedTicketId));
        card.setAttribute("aria-label", card.getAttribute("aria-label") + "; cites " + citedTicketId);
      }
    }
    if (request.domain) {
      var domainBadge = makeBadge("badge-domain", null, request.domain);
      if (request.domainUnrecognized) {
        domainBadge.appendChild(createElement("span", "status-invalid-flag", "invalid"));
        domainBadge.title = 'Unrecognized domain "' + (request.originalDomain || request.domain) + '"';
      } else if (request.originalDomain && request.originalDomain !== request.domain) {
        domainBadge.title = 'Declared as "' + request.originalDomain + '"; normalized to "' + request.domain + '"';
      }
      badges.appendChild(domainBadge);
    }
    if (request.userRequestId) {
      badges.appendChild(makeBadge("badge-ur", null, request.userRequestId));
    }
    if (request.priority === "now" || request.priority === "later") {
      var priorityBadge = makeBadge(
        "badge-priority badge-priority-" + request.priority,
        null,
        request.priority
      );
      priorityBadge.title =
        "priority: " + request.priority + " — authored queue order inside this dependency group";
      badges.appendChild(priorityBadge);
    }
    if (request.route) {
      var routeBadge = makeBadge("badge-route", "route", request.route);
      if (request.routeUnrecognized) {
        routeBadge.appendChild(createElement("span", "status-invalid-flag", "invalid"));
        routeBadge.title = 'Unrecognized route "' + (request.originalRoute || request.route) + '"';
      } else if (request.originalRoute && request.originalRoute !== request.route) {
        routeBadge.title = 'Declared as "' + request.originalRoute + '"; normalized to "' + request.route + '"';
      }
      badges.appendChild(routeBadge);
    }
    if (request.status === "blocked" && request.blockedBy && request.blockedBy.length > 0) {
      // Waiting on an external condition (a service being up, a person answering)
      // — distinct from pending-answers (user questions) and depends_on (another
      // REQ). The badge follows blocked status across either active-work column.
      var blockedCondition = request.blockedBy.join(", ");
      var blockedBadge = makeBadge("badge-blocked", "blocked by", truncateBadgeText(blockedCondition));
      var blockedTitle = blockedCondition;
      if (request.blockedAt) {
        blockedTitle += " — since " + formatShortInstantWithRelative(request.blockedAt);
      }
      blockedTitle += request.blockedCheck
        ? " — auto-probe set; `do-work run` re-checks it and unblocks on exit 0"
        : " — clear via `do-work clarify` or by editing the REQ once the condition is met";
      blockedBadge.title = blockedTitle;
      badges.appendChild(blockedBadge);
    }
    var unblockedRequestIds = activeDependentIds(request);
    if (unblockedRequestIds.length > 0 && !isTerminalResolvedStatus(request.status)) {
      var unblocksBadge = makeBadge("badge-unblocks", "unblocks", String(unblockedRequestIds.length));
      unblocksBadge.title = "Unblocks " + unblockedRequestIds.join(", ") + " when this lands";
      badges.appendChild(unblocksBadge);
    }
    var writeSetOverlapIds = request.writeSetOverlaps || [];
    if (writeSetOverlapIds.length > 0) {
      // Another pending/claimed REQ declares a write_set that could touch the
      // same files. The Go side (annotateWriteSetOverlap) does the pairwise
      // comparison and gates it to the pending/claimed tier — this only renders
      // the derived list. Display only: nothing here blocks, sorts, or moves a
      // card, and nothing schedules on write_set at any builder count.
      var overlapBadge = makeBadge(
        "badge-write-overlap",
        "overlaps",
        truncateBadgeText(writeSetOverlapIds.join(", "), 24)
      );
      overlapBadge.title =
        "Declared write_set could touch the same files as " +
        writeSetOverlapIds.join(", ") +
        " — an informational heads-up about declared file contention. Display only: " +
        "the board never blocks or reorders on this, and nothing schedules on write_set — " +
        "when several builders run at once it is advisory input to your pick, and the merge is what proves they did not collide.";
      badges.appendChild(overlapBadge);
    }
    if (request.assignedTo) {
      // The advisory cooperative claim marker: this REQ is earmarked for another
      // session. The value is never normalized — session names have no canonical
      // vocabulary, so nothing folds case, maps aliases, or rewrites it. Only the
      // badge's visible text is truncated for layout; the title tooltip below and
      // the drawer row carry the full value. Display only: the board never
      // buckets, orders, or hides a card on it. The one reader that acts on it is
      // the work pipeline's default scan, which skips and reports an assigned REQ
      // and is overridden by explicitly targeting it.
      var assignedBadge = makeBadge(
        "badge-assigned",
        "assigned",
        truncateBadgeText(request.assignedTo, 18)
      );
      assignedBadge.title =
        "Earmarked for " +
        request.assignedTo +
        " — an advisory claim marker, not a lock. Another session's default run skips and reports it; " +
        "naming it explicitly overrides that and clears the field. Display only: the board never " +
        "reorders, blocks, or hides on this.";
      badges.appendChild(assignedBadge);
    }
    if (request.effortEstimate === "effort-mechanical" || request.effortEstimateUnrecognized) {
      // The size triage bit: this REQ is a small mechanical fix, not real work.
      // Only `effort-mechanical` chips — `effort-substantive` is the default and
      // would chip every card into noise. Display only: the board never buckets,
      // orders, or schedules on it; it exists so the user can tell at a glance
      // which queued follow-ups are cheap to approve or batch.
      //
      // An UNRECOGNIZED value also chips, even though it resolves to
      // `effort-substantive`: the resolved value alone would render nothing, so
      // the card would be the one place a typo'd effort_estimate left no
      // footprint — the same never-silently-drop leg domain and route carry on
      // their badges.
      var effortEstimateBadge = makeBadge(
        "badge-effort-estimate",
        null,
        request.effortEstimate || "effort-substantive"
      );
      if (request.effortEstimateUnrecognized) {
        effortEstimateBadge.appendChild(createElement("span", "status-invalid-flag", "invalid"));
        effortEstimateBadge.title =
          'Unrecognized effort_estimate "' +
          (request.originalEffortEstimate || request.effortEstimate) +
          '" — expected effort-mechanical or effort-substantive; treated as effort-substantive.';
      } else {
        effortEstimateBadge.title =
          "effort_estimate: effort-mechanical — a small mechanical fix. This is the SIZE axis, " +
          "judged by whoever wrote the REQ and never derived from the impact verdict. " +
          "Display only: the board never reorders, blocks, or hides on this.";
      }
      badges.appendChild(effortEstimateBadge);
    }
    if (
      (request.impact && request.impact !== "impact-user-visible") ||
      request.impactUnrecognized
    ) {
      // Whether anyone would ever notice the work — the axis effort_estimate was
      // never able to answer. Every value chips except `impact-user-visible`,
      // which is the contract default and would chip every card into noise.
      // Display only: the board never buckets, orders, or schedules on it.
      //
      // An UNRECOGNIZED value also chips, for the same never-silently-drop
      // reason the effort chip above carries: it resolves to the default, which
      // renders nothing, so without this leg a typo'd impact would leave no
      // footprint on the card.
      var impactBadge = makeBadge("badge-impact", null, request.impact || "impact-user-visible");
      if (request.impactUnrecognized) {
        impactBadge.appendChild(createElement("span", "status-invalid-flag", "invalid"));
        impactBadge.title =
          'Unrecognized impact "' +
          (request.originalImpact || request.impact) +
          '" — expected impact-critical, impact-user-visible, impact-rule-change, or ' +
          "impact-negligible; treated as impact-user-visible.";
      } else {
        impactBadge.title =
          "impact: " +
          request.impact +
          " — whether anyone would ever notice this work, judged at capture or at follow-up " +
          "creation. Display only: the board never reorders, blocks, or hides on this.";
      }
      badges.appendChild(impactBadge);
    }
    if (request.sweep) {
      // A consolidation sweep holds many instances behind one card, so the raw
      // column count understates real work — the chip's open/total count is
      // what keeps queue depth honest. Display only: the board never buckets,
      // orders, or schedules on the marker.
      var sweepInstancesOpen = request.sweepInstancesOpen || 0;
      var sweepInstancesDone = request.sweepInstancesDone || 0;
      var sweepInstancesTotal = sweepInstancesOpen + sweepInstancesDone;
      var sweepBadge = makeBadge(
        "badge-sweep",
        "sweep",
        sweepInstancesOpen + " open of " + sweepInstancesTotal
      );
      sweepBadge.title =
        "sweep: true — one REQ per root cause, holding an ## Instances checklist. " +
        sweepInstancesOpen +
        " unticked instance(s) of " +
        sweepInstancesTotal +
        ". Display only: the board never reorders, blocks, or hides on this.";
      badges.appendChild(sweepBadge);
    }
    if (request.completionAnomaly) {
      // Broken completion bookkeeping (flagged by the Go side) — mark the card
      // wherever it renders, not just inside the anomalies strip.
      card.classList.add("is-completion-anomaly");
      var anomalyBadge = makeBadge("badge-completion-anomaly", null, "anomaly");
      anomalyBadge.title =
        "Completion anomaly: " +
        (request.completionAnomalyReason || "completion instant unresolved") +
        " — fix: add completed_at: <ISO instant> and/or a valid commit hash field to the REQ frontmatter.";
      badges.appendChild(anomalyBadge);
    }
    if (request.futureTimestampFields && request.futureTimestampFields.length > 0) {
      // A frontmatter stamp later than the board's generation time (+2min skew)
      // — flagged instead of rendered silently, since every elapsed-time
      // reading derived from it is wrong for as long as the stamp stands. The
      // cause text is shared with the stopwatch tooltip (board-core.js's
      // futureStampCauseText) so the two can never disagree.
      var futureStampBadge = makeBadge("badge-future-timestamp", null, "⚠ future stamp");
      futureStampBadge.title = futureStampTooltipText(request.futureTimestampFields);
      badges.appendChild(futureStampBadge);
    }
    if (request.testingStatus) {
      // The testing track (see the Testing view) surfaces on the main board too,
      // so a finished card's tested/returned state is visible without switching.
      var testingBadge = makeBadge(
        "badge-testing badge-testing-" + request.testingStatus,
        "testing",
        request.testingStatus
      );
      if (request.testedBy) {
        testingBadge.title = request.testingStatus + " by " + request.testedBy;
      }
      badges.appendChild(testingBadge);
    }
    if (badges.childNodes.length > 0) {
      card.appendChild(badges);
    }

    if (request.dependsOn && request.dependsOn.length > 0) {
      var unmetDependencyIds = request.unmetDependencies || [];
      var deps = createElement("div", "req-card-deps");
      deps.appendChild(createElement("span", "dep-chip-lead", "needs"));
      request.dependsOn.forEach(function (dependencyId) {
        var isUnmet = unmetDependencyIds.indexOf(dependencyId) !== -1;
        var chip = createElement("span", isUnmet ? "dep-chip is-unmet" : "dep-chip is-met", dependencyId);
        chip.title = dependencyId + " — " + describeRequestStatus(dependencyId);
        deps.appendChild(chip);
      });
      card.appendChild(deps);
    }

    var stateTimerSpec = stateTimerSpecFor(request);
    if (stateTimerSpec) {
      var stateTimerNode = makeInstantWithStopwatchNode(stateTimerSpec.instantIso);
      if (stateTimerNode) {
        var stateTimerLine = createElement("div", "req-card-completed", stateTimerSpec.verbText + " ");
        stateTimerLine.appendChild(stateTimerNode);
        card.appendChild(stateTimerLine);
      }
    }

    if (options && options.showCompleted && request.completionTime) {
      var completionVerb = request.status === "cancelled" ? "cancelled" : "done";
      var completionLine = createElement("div", "req-card-completed", completionVerb + " ");
      var completionInstantNode = makeInstantWithRelativeNode(request.completionTime);
      if (completionInstantNode) {
        completionLine.appendChild(completionInstantNode);
      }
      var implementationSpanNode = makeImplementationSpanNode(request);
      if (implementationSpanNode) {
        completionLine.appendChild(implementationSpanNode);
      }
      card.appendChild(completionLine);
    }

    return card;
  }

  // Column counts read "shown / total" while a filter hides cards, so a
  // filtered column is never mistaken for an empty one.
  function formatFilteredCount(shownCount, totalCount) {
    return shownCount < totalCount ? shownCount + " / " + totalCount : String(shownCount);
  }

  function columnEmptyText(filtersActive) {
    var resolvedFiltersActive =
      typeof filtersActive === "boolean"
        ? filtersActive
        : hasActiveFilters();
    return resolvedFiltersActive ? "No matches" : "Nothing here";
  }

  // Copy availability is derived from the same filtered/windowed id slice that
  // fills the column. The click path reads the rendered cards back in DOM order,
  // so this count is only the accessible/disabled state — never a second copy of
  // which REQs the button will write.
  function updateColumnCopyButton(columnKey, shownCount) {
    var copyButton = document.querySelector('[data-copy-column="' + columnKey + '"]');
    if (!copyButton) {
      return;
    }
    var columnLabel = copyButton.dataset.copyLabel || columnKey;
    copyButton.disabled = shownCount === 0;
    copyButton.dataset.visibleCount = String(shownCount);
    copyButton.setAttribute(
      "aria-label",
      shownCount === 0
        ? "No " + columnLabel + " REQs to copy"
        : "Copy all " + shownCount + " " + columnLabel + " REQ" + (shownCount === 1 ? "" : "s")
    );
  }

  function fillColumn(columnKey, requestIds, options, totalCount) {
    var container = document.querySelector('[data-cards="' + columnKey + '"]');
    var countNode = document.querySelector('[data-count="' + columnKey + '"]');
    container.textContent = "";
    if (countNode) {
      countNode.textContent = formatFilteredCount(requestIds.length, totalCount != null ? totalCount : requestIds.length);
    }
    // Some focused behavior probes execute fillColumn in isolation. The real
    // assembled client always carries the helper; the guard keeps that older
    // function-level seam focused on its empty-copy decision.
    if (typeof updateColumnCopyButton === "function") {
      updateColumnCopyButton(columnKey, requestIds.length);
    }
    if (requestIds.length === 0) {
      container.appendChild(createElement("p", "column-empty", columnEmptyText()));
      return;
    }
    requestIds.forEach(function (requestId) {
      container.appendChild(makeRequestCard(requestId, options));
    });
  }

  // The Pending column is the only one that sub-groups: what the work loop could
  // claim right now, versus what is still waiting — on an upstream REQ, or for
  // the heavy-lane drain the loop runs at queue exhaustion. When
  // nothing is waiting, the headers are noise — the column renders as a flat
  // list, exactly as it did before dependency readiness was computed.
  function fillPendingColumn(readyIds, waitingIds, totalCount) {
    var container = document.querySelector('[data-cards="pending"]');
    var countNode = document.querySelector('[data-count="pending"]');
    container.textContent = "";
    countNode.textContent = formatFilteredCount(readyIds.length + waitingIds.length, totalCount);
    updateColumnCopyButton("pending", readyIds.length + waitingIds.length);

    if (readyIds.length === 0 && waitingIds.length === 0) {
      container.appendChild(createElement("p", "column-empty", columnEmptyText()));
      return;
    }
    if (waitingIds.length === 0) {
      readyIds.forEach(function (requestId) {
        container.appendChild(makeRequestCard(requestId));
      });
      return;
    }
    container.appendChild(makePendingGroup("Ready", readyIds, "Nothing ready — everything here is waiting"));
    container.appendChild(makePendingGroup("Waiting", waitingIds, ""));
  }

  function makePendingGroup(labelText, requestIds, emptyText) {
    var group = createElement("section", "pending-group");
    var header = createElement("h3", "pending-group-label");
    header.appendChild(createElement("span", "pending-group-name", labelText));
    header.appendChild(createElement("span", "pending-group-count", String(requestIds.length)));
    group.appendChild(header);

    if (requestIds.length === 0) {
      group.appendChild(createElement("p", "column-empty", emptyText));
      return group;
    }
    requestIds.forEach(function (requestId) {
      group.appendChild(makeRequestCard(requestId));
    });
    return group;
  }

  // ---- recently-done window (recomputed client-side) ----------------------

  function recentlyDoneIds(windowHours) {
    // Anchored to the wall clock at render time, not generatedAtMs: a tab left
    // open past the snapshot would otherwise keep counting the window back
    // from page-generation, so "last 24 hours" slowly stops meaning that.
    var cutoffMs = Date.now() - windowHours * 3600 * 1000;
    var ids = [];
    (boardData.calendar || []).forEach(function (entry) {
      // The calendar carries every REQ — queued, claimed, and failed included —
      // so this cannot treat an entry's presence as "done". The gate mirrors
      // Go's bucketColumns RecentlyDone exactly: terminal-RESOLVED only, which
      // keeps cancelled in (it shares that column) and keeps failed out (it
      // belongs to Needs-input/Blocked). Without it, a REQ claimed an hour ago
      // shows up in "Recently done".
      if (!isTerminalResolvedStatus(entry.status)) {
        return;
      }
      var ms = Date.parse(entry.entryTime);
      if (!isNaN(ms) && ms > cutoffMs) {
        ids.push(entry.id);
      }
    });
    return ids; // calendar is already most-recent-first
  }

  function applyRecentWindowSelection(selectedWindowHours) {
    viewState.windowHours = selectedWindowHours || 24;
    setActiveButton("#recent-window-group", "data-window-hours", String(viewState.windowHours));
    // The window scopes both lenses. Columns have no renderedOnce guard, while
    // the by-UR lens must either re-render now or stay stale until selected.
    renderColumns();
    renderedOnce.userRequestLens = false;
    if (viewState.view === "board" && viewState.lens === "user-request") {
      renderUserRequestLens();
      renderedOnce.userRequestLens = true;
    }
  }

  function filterRequestIds(requestIds) {
    return requestIds.filter(function (requestId) {
      return requestMatchesFilters(requestId);
    });
  }

  function renderColumns() {
    var columns = boardData.columns || {};
    var pendingReadyIds = columns.pendingReady || [];
    var pendingWaitingIds = columns.pendingWaiting || [];
    fillPendingColumn(
      filterRequestIds(pendingReadyIds),
      filterRequestIds(pendingWaitingIds),
      pendingReadyIds.length + pendingWaitingIds.length
    );
    var claimedIds = columns.claimed || [];
    fillColumn("claimed", filterRequestIds(claimedIds), null, claimedIds.length);
    var needsInputIds = columns.needsInputOrBlocked || [];
    fillColumn("needsInputOrBlocked", filterRequestIds(needsInputIds), null, needsInputIds.length);
    var recentIds = recentlyDoneIds(viewState.windowHours);
    fillColumn("recentlyDone", filterRequestIds(recentIds), { showCompleted: true }, recentIds.length);
  }

  // ---- data warnings banner ------------------------------------------------

  function renderWarningsBanner() {
    var warnings = boardData.warnings || [];
    if (warnings.length === 0) {
      return;
    }
    var banner = createElement("aside", "board-warnings");
    banner.setAttribute("role", "note");
    banner.appendChild(
      createElement(
        "strong",
        "board-warnings-title",
        warnings.length === 1 ? "1 data warning" : warnings.length + " data warnings"
      )
    );
    var list = createElement("ul", "board-warnings-list");
    warnings.forEach(function (warningText) {
      list.appendChild(createElement("li", null, warningText));
    });
    banner.appendChild(list);
    var main = document.getElementById("board-main");
    main.insertBefore(banner, main.firstChild);
  }

  // ---- completion anomalies strip -----------------------------------------
  // Terminal REQs whose completion bookkeeping is broken (columns
  // .completionAnomalies, flagged by detectCompletionAnomaly in model.go).
  // Some classes carry no honest completion instant; a reversed span carries
  // one that cannot be real. Either way they are listed here as data bugs to
  // fix — this strip is never aged out by the 24h/48h/7d window and stays
  // visible from every view (a ticket with a resolvable in-window instant may
  // additionally appear under Recently done).
  // Deliberately exempt from the shared filters: an anomaly must not be
  // hideable by a filter combination.

  function renderAnomaliesStrip() {
    var anomalyIds = (boardData.columns || {}).completionAnomalies || [];
    var strip = document.getElementById("board-anomalies");
    if (anomalyIds.length === 0) {
      strip.hidden = true;
      return;
    }
    strip.hidden = false;
    document.getElementById("board-anomalies-count").textContent = String(anomalyIds.length);
    var cardsHost = document.getElementById("board-anomalies-cards");
    cardsHost.textContent = "";
    anomalyIds.forEach(function (requestId) {
      var entry = createElement("div", "board-anomaly-entry");
      entry.appendChild(makeRequestCard(requestId));
      var request = requestsById[requestId];
      if (request && request.completionAnomalyReason) {
        entry.appendChild(createElement("p", "board-anomaly-reason", request.completionAnomalyReason));
      }
      cardsHost.appendChild(entry);
    });
  }

  // ---- verify findings strip (REQ-285) -----------------------------------
  // `queue-kanban verify` detects sixteen categories of queue and process
  // breakage; the board rendered three, so the rest reached nobody who did not
  // run verify from a shell. This strip renders whatever the Go producer put in
  // boardData.verifyFindings — blindly, on purpose: the suppression of the
  // categories the board already shows another way happened in Go (REQ-284), and
  // a second copy of that judgment here is how the two would drift apart.
  //
  // Nearly the same exemptions as the anomalies strip above, for the same
  // reason: it lives outside the view panels, ignores the recently-done window,
  // and ignores the shared filters, because a finding must not be hideable by a
  // filter combination. The one exception is the Activity view, which hides the
  // strip from applyView (board-controls.js, REQ-578) — this renderer decides
  // emptiness and nothing else. Every string is set with textContent — a detail
  // or remedy is producer text that can carry any punctuation and must never
  // become markup.
  //
  // One row per finding, each behind its own click (REQ-589): a finding and a
  // probe that could not run are the same kind of thing to the reader, so they
  // share one row shape in one list, and the remedy — read only after deciding to
  // act — lives inside the row's own disclosure. Weight comes from the producer
  // alone: `fixable` (a command resolves it) and "this probe never ran" are the
  // only two things that colour a dot. No severity, colour scale or re-ordering
  // is invented here.
  //
  // The shell is static in web/template.html — the band, the outer disclosure,
  // the label, the counts, the Show/Hide labels and the two hosts. This renderer
  // fills the closed line's subject list and the rows and decides nothing else,
  // which is what keeps the strip's shape readable in one file.

  var verifyFindingsOpenStorageKey = "queueKanbanVerifyFindingsOpen";

  function renderVerifyFindingsStrip() {
    var findings = boardData.verifyFindings || [];
    var skipped = boardData.verifySkipped || [];
    var strip = document.getElementById("board-findings");
    if (!strip) {
      return; // an older template; the payload is still valid without this strip
    }
    // Two hosts, one list: they are `display: contents`, so every row below is
    // laid out by the list element that wraps them. The split is per payload
    // array, and applyView reads these two ids to decide whether the strip has
    // anything to say (board-controls.js, REQ-578).
    //
    // All three are cleared BEFORE the empty check, so a re-render with nothing
    // to report leaves nothing behind. Hiding the strip is not enough on its own:
    // applyView asks the two row hosts whether there is content, so stale rows
    // under a hidden strip would put it back on screen at the next view switch.
    var findingsHost = document.getElementById("board-findings-cards");
    var skippedHost = document.getElementById("board-findings-skipped-list");
    var subjectsHost = document.getElementById("board-findings-subjects");
    findingsHost.textContent = "";
    skippedHost.textContent = "";
    subjectsHost.textContent = "";

    if (findings.length === 0 && skipped.length === 0) {
      strip.hidden = true;
      return;
    }
    strip.hidden = false;
    document.getElementById("board-findings-count").textContent =
      formatFindingsSummary(findings.length, skipped.length);

    groupFindingsBySubject(findings).forEach(function (findingGroup) {
      findingGroup.findings.forEach(function (finding) {
        // The closed line names every FINDING once, not every group once: two
        // findings on one worktree are two things to answer for. A finding the
        // producer gave no subject is named by its category, the only identifier
        // it has — the subject is never parsed back out of the detail sentence.
        subjectsHost.appendChild(makeSubjectListItem(
          finding.subject || categoryWords(finding.category),
          finding.fixable ? "board-findings-dot board-findings-dot-fixable" : "board-findings-dot"
        ));
        findingsHost.appendChild(makeFindingRow(finding));
      });
    });

    // The probes that never ran collapse to one grey entry on the closed line and
    // stay one row each below it: the count is what a reader needs before opening.
    if (skipped.length > 0) {
      subjectsHost.appendChild(makeSubjectListItem(
        formatSkippedProbeCount(skipped.length),
        "board-findings-dot board-findings-dot-skipped"
      ));
    }
    skipped.forEach(function (skippedProbe) {
      skippedHost.appendChild(makeSkippedProbeRow(skippedProbe));
    });
  }

  // "2 findings · 1 probe not checked". Each half joins only when it has a count,
  // so the header never advertises a non-event ("0 probes not checked").
  function formatFindingsSummary(findingCount, skippedCount) {
    var summaryParts = [];
    if (findingCount > 0) {
      summaryParts.push(findingCount + (findingCount === 1 ? " finding" : " findings"));
    }
    if (skippedCount > 0) {
      summaryParts.push(formatSkippedProbeCount(skippedCount));
    }
    return summaryParts.join(" · ");
  }

  // The same phrase in the header and on the closed line's grey entry, written
  // once so the two can never disagree about the plural.
  function formatSkippedProbeCount(skippedCount) {
    return skippedCount + (skippedCount === 1 ? " probe not checked" : " probes not checked");
  }

  // The producer's category token as words: lowercased, hyphens shown as spaces
  // (REQ-589 D5). A mechanical transform and never a lookup table — a category
  // the board has never heard of has to read as words too, and a table would go
  // stale the first time verify grows a probe.
  function categoryWords(category) {
    return String(category || "finding").toLowerCase().replace(/-/g, " ");
  }

  // Findings sharing one non-empty `subject` become one group, in first-seen
  // order; findings the producer gave no subject stay in producer order in a
  // trailing group. The match is exact string equality on the payload field —
  // the subject is never parsed back out of the detail sentence, which is prose
  // and free to change. The lookup map has a null prototype so a subject spelled
  // like an Object member ("constructor") cannot collide with one.
  function groupFindingsBySubject(findings) {
    var groupsBySubject = Object.create(null);
    var orderedGroups = [];
    var unsubjectedGroup = { subject: "", findings: [] };
    findings.forEach(function (finding) {
      var subject = finding.subject || "";
      if (subject === "") {
        unsubjectedGroup.findings.push(finding);
        return;
      }
      if (!groupsBySubject[subject]) {
        groupsBySubject[subject] = { subject: subject, findings: [] };
        orderedGroups.push(groupsBySubject[subject]);
      }
      groupsBySubject[subject].findings.push(finding);
    });
    if (unsubjectedGroup.findings.length > 0) {
      orderedGroups.push(unsubjectedGroup);
    }
    return orderedGroups;
  }

  // One entry on the closed line: the weight dot, then the name in the mono face.
  function makeSubjectListItem(labelText, dotClassName) {
    var subjectItem = createElement("span", "board-findings-subject-item");
    subjectItem.appendChild(createElement("span", dotClassName));
    subjectItem.appendChild(createElement("span", "board-findings-subject", labelText));
    return subjectItem;
  }

  // One row: a disclosure whose summary is the whole scannable line — dot,
  // subject, category, the detail clipped to the line's end, the "cleanup can
  // fix" pill, the chevron — and whose content is the remedy. Splitting them this
  // way is the point of the REQ: the summary says what is wrong, and the reader
  // asks for what to do about it.
  function makeFindingRow(finding) {
    var row = createElement("details", "board-findings-row");
    var rowSummary = createElement("summary");
    rowSummary.appendChild(createElement(
      "span",
      // Exactly verify's meaning: `do-work cleanup` can resolve this one
      // mechanically. Never inferred here — the producer sets the flag.
      finding.fixable ? "board-findings-dot board-findings-dot-fixable" : "board-findings-dot"
    ));
    if (finding.subject) {
      rowSummary.appendChild(createElement("span", "board-findings-subject", finding.subject));
    }
    rowSummary.appendChild(createElement("span", "board-findings-category", categoryWords(finding.category)));
    rowSummary.appendChild(createElement("span", "board-findings-detail", finding.detail || ""));
    if (finding.fixable) {
      rowSummary.appendChild(createElement("span", "board-findings-fixable", "cleanup can fix"));
    }
    rowSummary.appendChild(createFindingChevronIcon());
    row.appendChild(rowSummary);
    if (finding.remedy) {
      var remedyBlock = createElement("div", "board-findings-remedy");
      remedyBlock.appendChild(createElement("span", "board-findings-remedy-label", "What to do:"));
      remedyBlock.appendChild(createElement("span", "board-findings-remedy-text", finding.remedy));
      row.appendChild(remedyBlock);
    }
    return row;
  }

  // A probe that could not run is a row of the same shape: the grey dot, the
  // category "not checked", and the producer's whole sentence as the detail.
  // That sentence is NOT split into a subject and a reason here — verify ships a
  // skipped probe as one string with no subject field, so carving one out would
  // mean reading a subject back out of prose, which is the mistake REQ-588
  // removed everywhere else. Opening the row unclips the sentence instead of
  // showing a remedy: a probe that never ran has nothing to recommend.
  function makeSkippedProbeRow(skippedProbe) {
    var row = createElement("details", "board-findings-row");
    var rowSummary = createElement("summary");
    rowSummary.appendChild(createElement("span", "board-findings-dot board-findings-dot-skipped"));
    rowSummary.appendChild(createElement("span", "board-findings-category", "not checked"));
    rowSummary.appendChild(createElement("span", "board-findings-detail", skippedProbe));
    rowSummary.appendChild(createFindingChevronIcon());
    row.appendChild(rowSummary);
    return row;
  }

  // Inline SVG must be created in the SVG namespace: document.createElement
  // returns an HTML element of the same name that never paints. An SVG element's
  // `className` is read-only, so its class goes through setAttribute. Every value
  // here is a literal from this file — no payload text reaches an SVG attribute.
  function createFindingChevronIcon() {
    var svgNamespaceUri = "http://www.w3.org/2000/svg";
    var chevronIcon = document.createElementNS(svgNamespaceUri, "svg");
    chevronIcon.setAttribute("class", "board-findings-chevron");
    chevronIcon.setAttribute("viewBox", "0 0 16 16");
    chevronIcon.setAttribute("aria-hidden", "true");
    var chevronPath = document.createElementNS(svgNamespaceUri, "path");
    chevronPath.setAttribute("d", "M6 3.5 10.5 8 6 12.5");
    chevronPath.setAttribute("fill", "none");
    chevronPath.setAttribute("stroke", "currentColor");
    chevronPath.setAttribute("stroke-width", "1.6");
    chevronPath.setAttribute("stroke-linecap", "round");
    chevronPath.setAttribute("stroke-linejoin", "round");
    chevronIcon.appendChild(chevronPath);
    return chevronIcon;
  }

  // Whether the band is open is remembered per browser under one key (REQ-589
  // D4), best-effort exactly like the detail panel's width in board-detail.js: a
  // browser that denies storage loses the memory, never the strip. Default
  // closed, so a reader who has never touched it gets the one-line band. Row
  // state is deliberately not remembered — an open remedy is about the finding
  // being read right now, not a preference.
  function readStoredVerifyFindingsOpenState() {
    try {
      return localStorage.getItem(verifyFindingsOpenStorageKey) === "open";
    } catch (storageError) {
      return false; // fall through to the closed default
    }
  }

  function persistVerifyFindingsOpenState(isStripOpen) {
    try {
      localStorage.setItem(verifyFindingsOpenStorageKey, isStripOpen ? "open" : "closed");
    } catch (storageError) {
      // Persistence is best-effort; the strip already opened or closed.
    }
  }

  // Wired once at load rather than inside the renderer: the disclosure is static
  // in the template, so one listener outlives every re-render, and the restore
  // does not have to wait for a payload it never reads. The handler asks the
  // element for its own state instead of tracking a direction — `toggle` fires
  // for a click, for a keypress, and for anything else that opens a <details>.
  (function restoreVerifyFindingsOpenState() {
    var stripDisclosure = document.getElementById("board-findings-strip");
    if (!stripDisclosure) {
      return; // an older template; the rest of the strip still renders
    }
    stripDisclosure.open = readStoredVerifyFindingsOpenState();
    stripDisclosure.addEventListener("toggle", function () {
      persistVerifyFindingsOpenState(stripDisclosure.open);
    });
  })();

  // ---- notes strip (do-work/notes.md) -------------------------------------
  // Notes are plain text, never Markdown: they are appended verbatim by
  // `do-work-toolbox note` and rendered with textContent, so a stray `<` or a pasted
  // tag in a hint can never become markup.

  function renderNotesStrip() {
    var notes = boardData.notes || [];
    var strip = document.getElementById("board-notes");
    if (notes.length === 0) {
      strip.hidden = true;
      return;
    }
    strip.hidden = false;
    document.getElementById("board-notes-count").textContent = String(notes.length);

    var list = document.getElementById("board-notes-list");
    list.textContent = "";
    notes.forEach(function (note) {
      var item = createElement("li", "board-note");
      if (note.date) {
        var dateNode = createElement("time", "board-note-date", note.date);
        dateNode.setAttribute("datetime", note.date);
        item.appendChild(dateNode);
      }
      item.appendChild(createElement("span", "board-note-text", note.text || ""));
      list.appendChild(item);
    });
  }

  // ---- by-UR lens ---------------------------------------------------------

  function renderUserRequestLens() {
    var host = document.getElementById("user-request-lens");
    host.textContent = "";
    var hiddenResolvedCount = 0;
    var hiddenResolvedFilterMatchCount = 0;

    // Built once per render, not once per UR: recentlyDoneIds walks the whole
    // calendar, and a real tree has hundreds of URs to test against it.
    var recentlyDoneIdSet = {};
    recentlyDoneIds(viewState.windowHours).forEach(function (requestId) {
      recentlyDoneIdSet[requestId] = true;
    });

    (boardData.userRequestOrder || []).forEach(function (userRequestId) {
      var userRequest = userRequestsById[userRequestId];
      if (!userRequest) {
        return;
      }

      // Compute non-scope filters before the Active gate. If a search matches a
      // resolved UR that Active hides, switching to All is a real escape and the
      // empty state must say so instead of claiming the filters matched nothing.
      var groupMatchesSearch =
        filterState.searchText !== "" &&
        searchMatchesUserRequest(userRequest, userRequestId, filterState.searchText);
      var requestIds = userRequest.requestIds || [];
      var shownRequestIds = requestIds.filter(function (requestId) {
        return requestMatchesFilters(requestId, { skipSearch: groupMatchesSearch });
      });
      if (
        filterState.userRequestActivity === "active" &&
        !userRequestHasOpenOrRecentWork(userRequest, recentlyDoneIdSet)
      ) {
        hiddenResolvedCount += 1;
        if (hasActiveFilters() && shownRequestIds.length > 0) {
          hiddenResolvedFilterMatchCount += 1;
        }
        return;
      }
      if (hasActiveFilters() && shownRequestIds.length === 0) {
        return;
      }

      var group = createElement("section", "ur-group");

      // Both readings of this lens build the same grid; the folded one just
      // builds it late, when the reader opens the row, and drops it again on
      // collapse. Nothing outlives the nodes — the next render rebuilds them,
      // so the click listener below needs no teardown.
      function makeUserRequestCards() {
        var cardsNode = createElement("div", "ur-group-cards");
        shownRequestIds.forEach(function (requestId) {
          cardsNode.appendChild(makeRequestCard(requestId, { showCompleted: true }));
        });
        return cardsNode;
      }

      // One shape of head for both readings. The head is the fold control, so
      // it cannot also be the drawer trigger — one element must not mean two
      // things. The drawer lives on its own button beside it inside
      // .ur-group-row (both are real buttons, both keyboard-operable), and the
      // fold state is announced on the element that owns it. The reading
      // decides exactly one thing: where the fold starts.
      var head = createElement("button", "ur-group-head");
      head.type = "button";
      head.setAttribute("aria-expanded", userRequestCardsFolded ? "false" : "true");
      head.appendChild(createElement("span", "ur-fold-marker", "▸"));
      head.appendChild(createElement("span", "ur-id", userRequestId));
      head.appendChild(createElement("span", "ur-title", userRequest.title || "(no input.md title)"));
      if (groupMatchesSearch && !searchMatchesUserRequest(userRequest, userRequestId, filterState.searchText, true)) {
        var citationBadge = createElement("span", "badge citation-match", "cites");
        citationBadge.title = "Cites " + citationMatchedTicketId(userRequest, filterState.searchText);
        citationBadge.setAttribute("aria-label", citationBadge.title);
        head.appendChild(citationBadge);
      }
      if (!userRequest.inputFilePresent) {
        head.appendChild(createElement("span", "ur-synthetic", "no input.md"));
      }
      head.appendChild(
        createElement(
          "span",
          "ur-count",
          shownRequestIds.length < requestIds.length
            ? shownRequestIds.length + " / " + requestIds.length + " REQ"
            : requestIds.length + " REQ"
        )
      );
      var headRow = createElement("div", "ur-group-row");
      headRow.appendChild(head);
      var detailButton = createElement("button", "ur-group-detail", "Details");
      detailButton.type = "button";
      detailButton.dataset.detailKind = "ur";
      detailButton.dataset.detailId = userRequestId;
      detailButton.setAttribute("aria-label", "Open details for " + userRequestId);
      headRow.appendChild(detailButton);
      group.appendChild(headRow);

      // By UR arrives with its grid already built; URs only builds it on first
      // open. From the first click on they are the same machine.
      var openCards = userRequestCardsFolded ? null : makeUserRequestCards();
      if (openCards) {
        group.appendChild(openCards);
      }
      head.addEventListener("click", function () {
        if (openCards) {
          group.removeChild(openCards);
          openCards = null;
          head.setAttribute("aria-expanded", "false");
          return;
        }
        openCards = makeUserRequestCards();
        group.appendChild(openCards);
        head.setAttribute("aria-expanded", "true");
      });

      host.appendChild(group);
    });

    var windowPhrase = recentWindowPhrase(viewState.windowHours);
    var listIsEmpty = host.childNodes.length === 0;

    if (listIsEmpty) {
      // Three distinct reasons for an empty list, and the reader can only act on
      // the right one: a filter matched nothing, the Active scope hid everything
      // outside the window, or the tree genuinely has no URs. The middle branch
      // names both escapes — widen the window, or drop the scope.
      var emptyText = userRequestLensEmptyText(
        hasActiveFilters(),
        hiddenResolvedCount,
        hiddenResolvedFilterMatchCount,
        windowPhrase
      );
      host.appendChild(createElement("p", "ur-lens-empty", emptyText));
    }

    // The note used to sit behind an early return, so it never fired in the one
    // case it exists for: every UR hidden. It stays silent when filters emptied
    // the list; if those filters matched only scope-hidden URs, the empty-state
    // decision above already gives the precise All-scope escape and match count.
    if (hiddenResolvedCount > 0 && !(listIsEmpty && hasActiveFilters())) {
      host.appendChild(
        createElement(
          "p",
          "ur-lens-hidden-note",
          hiddenResolvedCount +
            " UR" +
            (hiddenResolvedCount === 1 ? "" : "s") +
            " with no open work or activity in " +
            windowPhrase +
            " hidden — switch URs to All to see them."
        )
      );
    }
  }
