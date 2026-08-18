  // ---- card construction --------------------------------------------------

  function truncateBadgeText(fullText, maxLength) {
    var limit = maxLength || 48;
    if (!fullText || fullText.length <= limit) {
      return fullText || "";
    }
    return fullText.slice(0, limit - 1).replace(/\s+$/, "") + "…";
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
      // REQ). Shares the Needs-input/Blocked column but names its condition.
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
    if (request.effortEstimate === "trivial" || request.effortEstimateUnrecognized) {
      // The triage bit from the review gate: this REQ is a small mechanical fix,
      // not real work. Only `trivial` chips — `normal` is the default and would
      // chip every card into noise. Display only: the board never buckets,
      // orders, or schedules on it; it exists so the user can tell at a glance
      // which queued follow-ups are cheap to approve or batch.
      //
      // An UNRECOGNIZED value also chips, even though it resolves to `normal`:
      // the resolved value alone would render nothing, so the card would be the
      // one place a typo'd effort_estimate left no footprint — the same
      // never-silently-drop leg domain and route carry on their badges.
      var effortEstimateBadge = makeBadge("badge-effort-estimate", null, request.effortEstimate || "normal");
      if (request.effortEstimateUnrecognized) {
        effortEstimateBadge.appendChild(createElement("span", "status-invalid-flag", "invalid"));
        effortEstimateBadge.title =
          'Unrecognized effort_estimate "' +
          (request.originalEffortEstimate || request.effortEstimate) +
          '" — expected trivial or normal; treated as normal.';
      } else {
        effortEstimateBadge.title =
          "effort_estimate: trivial — a small mechanical fix, stamped by the review gate " +
          "(or capture). Display only: the board never reorders, blocks, or hides on this.";
      }
      badges.appendChild(effortEstimateBadge);
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
      futureStampBadge.title =
        "Future-dated timestamp(s): " +
        request.futureTimestampFields.join(", ") +
        " — later than the board's generation time (2min skew allowance). Likely " +
        futureStampCauseText +
        "; fix: rewrite with the current UTC " +
        "instant — YYYY-MM-DDTHH:MM:SSZ, per the Timestamp rule in actions/work-reference.md.";
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

  function fillColumn(columnKey, requestIds, options, totalCount) {
    var container = document.querySelector('[data-cards="' + columnKey + '"]');
    var countNode = document.querySelector('[data-count="' + columnKey + '"]');
    container.textContent = "";
    if (countNode) {
      countNode.textContent = formatFilteredCount(requestIds.length, totalCount != null ? totalCount : requestIds.length);
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
  // claim right now, versus what is still waiting on an upstream REQ. When
  // nothing is waiting, the headers are noise — the column renders as a flat
  // list, exactly as it did before dependency readiness was computed.
  function fillPendingColumn(readyIds, waitingIds, totalCount) {
    var container = document.querySelector('[data-cards="pending"]');
    var countNode = document.querySelector('[data-count="pending"]');
    container.textContent = "";
    countNode.textContent = formatFilteredCount(readyIds.length + waitingIds.length, totalCount);

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
    container.appendChild(makePendingGroup("Ready", readyIds, "Nothing ready — every pending REQ is waiting"));
    container.appendChild(makePendingGroup("Waiting on dependencies", waitingIds, ""));
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
      var ms = Date.parse(entry.completionTime);
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

      var head = createElement("button", "ur-group-head");
      head.type = "button";
      if (userRequestCardsFolded) {
        // The row is the fold control here, so it cannot also be the drawer
        // trigger — one element must not mean two things. The drawer moves to
        // its own button beside it (both are real buttons, both keyboard-
        // operable), and the fold state is announced on the row that owns it.
        head.setAttribute("aria-expanded", "false");
        head.appendChild(createElement("span", "ur-fold-marker", "▸"));
      } else {
        head.dataset.detailKind = "ur";
        head.dataset.detailId = userRequestId;
      }
      head.appendChild(createElement("span", "ur-id", userRequestId));
      head.appendChild(createElement("span", "ur-title", userRequest.title || "(no input.md title)"));
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
      if (userRequestCardsFolded) {
        var headRow = createElement("div", "ur-group-row");
        headRow.appendChild(head);
        var detailButton = createElement("button", "ur-group-detail", "Details");
        detailButton.type = "button";
        detailButton.dataset.detailKind = "ur";
        detailButton.dataset.detailId = userRequestId;
        detailButton.setAttribute("aria-label", "Open details for " + userRequestId);
        headRow.appendChild(detailButton);
        group.appendChild(headRow);

        var openCards = null;
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
      } else {
        group.appendChild(head);
        group.appendChild(makeUserRequestCards());
      }

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
