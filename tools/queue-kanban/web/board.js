/* ===========================================================================
   queue-kanban — static board behaviour
   Reads the sibling board-data.js payload and renders every view client-side.
   Raw Markdown stays in board-markdown.js until the first Copy click. No
   framework: plain DOM construction, event delegation, and a docked detail
   panel with a drag-to-resize divider.
   =========================================================================== */
(function () {
  "use strict";

  var boardData = window.queueKanbanBoardData;
  if (!boardData || typeof boardData !== "object") {
    document.getElementById("board-main").innerHTML =
      '<p style="color:#d97a59;padding:24px">Board data not loaded. Ensure board-data.js is present beside index.html.</p>';
    return;
  }

  var requestsById = boardData.requests || {};
  var userRequestsById = boardData.userRequests || {};
  var generatedAtMs = Date.parse(boardData.generatedAt);
  if (isNaN(generatedAtMs)) {
    generatedAtMs = Date.now();
  }

  var viewState = {
    view: "board", // "board" | "calendar" | "testing"
    lens: "flat", // "flat" | "user-request"
    windowHours: 24
  };

  // Shared filters — applied to whichever view is active. userRequestActivity
  // only affects the by-UR lens ("active" hides URs whose REQs are all resolved).
  // doneWindow only affects the testing view (its select is hidden elsewhere):
  // "" | hours-as-string ("24", "168", "720") | "old" (older than 30 days).
  var filterState = {
    searchText: "",
    domain: "",
    status: "",
    doneWindow: "",
    userRequestActivity: "active" // "active" | "all"
  };

  var renderedOnce = { userRequestLens: false, calendar: false, testing: false };

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

  // ---- dependency helpers -------------------------------------------------
  // Mirrors model.go's isTerminalResolvedStatus / isCompletedStatus. The board
  // never re-derives which dependencies are unmet — the Go side annotates that
  // (a dangling id counts as unmet, and `cancelled` never satisfies a
  // dependency) and ships it as request.unmetDependencies.

  function isTerminalResolvedStatus(status) {
    return status === "completed" || status === "completed-with-issues" || status === "cancelled";
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

  // ---- filtering ------------------------------------------------------------
  // Pure client-side: the data island already carries status, domain, UR id,
  // and titles, so every view filters the same record set with the same rules.

  function hasActiveFilters() {
    return (
      filterState.searchText !== "" ||
      filterState.domain !== "" ||
      filterState.status !== "" ||
      filterState.doneWindow !== ""
    );
  }

  function searchMatchesRequest(request, requestId, searchNeedle) {
    if (requestId.toLowerCase().indexOf(searchNeedle) !== -1) {
      return true;
    }
    if (request.title && request.title.toLowerCase().indexOf(searchNeedle) !== -1) {
      return true;
    }
    if (request.userRequestId && request.userRequestId.toLowerCase().indexOf(searchNeedle) !== -1) {
      return true;
    }
    return false;
  }

  function searchMatchesUserRequest(userRequest, userRequestId, searchNeedle) {
    if (userRequestId.toLowerCase().indexOf(searchNeedle) !== -1) {
      return true;
    }
    return Boolean(userRequest.title && userRequest.title.toLowerCase().indexOf(searchNeedle) !== -1);
  }

  // options.skipSearch: the by-UR lens sets it when the search already matched
  // the UR header — every card in a matched group stays visible (domain/status
  // still apply).
  function requestMatchesFilters(requestId, options) {
    var request = requestsById[requestId];
    if (!request) {
      // Ids outside the current tree carry no fields to filter on — hide them
      // whenever any filter is set, show them otherwise.
      return !hasActiveFilters();
    }
    if (filterState.domain !== "" && request.domain !== filterState.domain) {
      return false;
    }
    if (filterState.status !== "" && request.status !== filterState.status) {
      return false;
    }
    if (filterState.searchText !== "" && !(options && options.skipSearch)) {
      return searchMatchesRequest(request, requestId, filterState.searchText);
    }
    return true;
  }

  function userRequestIsActive(userRequest) {
    return (userRequest.requestIds || []).some(function (requestId) {
      var request = requestsById[requestId];
      return request && !isTerminalResolvedStatus(request.status);
    });
  }

  function populateFilterSelects() {
    var domainSet = {};
    var statusSet = {};
    Object.keys(requestsById).forEach(function (requestId) {
      var request = requestsById[requestId];
      if (request.domain) {
        domainSet[request.domain] = true;
      }
      if (request.status) {
        statusSet[request.status] = true;
      }
    });
    fillSelectOptions(document.getElementById("filter-domain"), Object.keys(domainSet).sort());
    fillSelectOptions(document.getElementById("filter-status"), Object.keys(statusSet).sort());
  }

  function fillSelectOptions(selectNode, values) {
    values.forEach(function (value) {
      var option = document.createElement("option");
      option.value = value;
      option.textContent = value;
      selectNode.appendChild(option);
    });
  }

  // A filter change re-renders whatever is on screen; the other views are
  // marked stale so they re-render with the new filters when switched to.
  function onFiltersChanged() {
    document.getElementById("filter-clear").hidden = !hasActiveFilters();
    // Columns have no renderedOnce guard (they render at boot), so refresh
    // them unconditionally; the lazily-rendered views refresh if visible and
    // go stale otherwise, re-rendering on their next activation.
    renderColumns();
    renderedOnce.userRequestLens = false;
    renderedOnce.calendar = false;
    renderedOnce.testing = false;
    if (viewState.view === "calendar") {
      renderCalendar();
      renderedOnce.calendar = true;
    } else if (viewState.view === "testing") {
      renderTestingView();
      renderedOnce.testing = true;
    } else if (viewState.lens === "user-request") {
      renderUserRequestLens();
      renderedOnce.userRequestLens = true;
    }
  }

  // ---- card construction --------------------------------------------------

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
      card.dataset.status = request.status; // lets CSS restyle terminal-but-unsuccessful cards (cancelled)
    }

    var top = createElement("div", "req-card-top");
    top.appendChild(createElement("span", "req-card-id", requestId));
    var status = createElement("span", "req-card-status");
    status.appendChild(createElement("span", "status-dot"));
    status.appendChild(document.createTextNode(request.status || "—"));
    if (request.statusUnrecognized) {
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
      badges.appendChild(makeBadge("badge-domain", null, request.domain));
    }
    if (request.userRequestId) {
      badges.appendChild(makeBadge("badge-ur", null, request.userRequestId));
    }
    if (request.route) {
      badges.appendChild(makeBadge("badge-route", "route", request.route));
    }
    if (request.status === "reserved") {
      // Allocated to a DIFFERENT worktree/cloud session (do-work reserve) — the
      // card is grayed out via [data-status="reserved"] CSS; the badge names the
      // owner and the stale flag (>24h) carries the recategorize suggestion.
      var reservedBadge = makeBadge("badge-reserved", "reserved for", request.reservedFor || "unknown session");
      reservedBadge.title =
        "Allocated to a different worktree/cloud session" +
        (request.reservedAt ? " since " + formatShortInstant(request.reservedAt) : "");
      badges.appendChild(reservedBadge);
      if (request.reservationStale) {
        var staleBadge = makeBadge("badge-reservation-stale", null, "stale >24h");
        staleBadge.title =
          "Reserved for more than 24h — the owning session may be dead. Recategorize: do-work release " +
          requestId +
          " (back to queue), do-work run " +
          requestId +
          " (claim here), or leave it if that session is still active.";
        badges.appendChild(staleBadge);
      }
    }
    var unblockedRequestIds = activeDependentIds(request);
    if (unblockedRequestIds.length > 0 && !isTerminalResolvedStatus(request.status)) {
      var unblocksBadge = makeBadge("badge-unblocks", "unblocks", String(unblockedRequestIds.length));
      unblocksBadge.title = "Unblocks " + unblockedRequestIds.join(", ") + " when this lands";
      badges.appendChild(unblocksBadge);
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

    if (options && options.showCompleted && request.completionTime) {
      var completionVerb = request.status === "cancelled" ? "cancelled" : "done";
      card.appendChild(
        createElement("div", "req-card-completed", completionVerb + " " + formatShortInstant(request.completionTime))
      );
    }

    return card;
  }

  // Column counts read "shown / total" while a filter hides cards, so a
  // filtered column is never mistaken for an empty one.
  function formatFilteredCount(shownCount, totalCount) {
    return shownCount < totalCount ? shownCount + " / " + totalCount : String(shownCount);
  }

  function columnEmptyText() {
    return hasActiveFilters() ? "No matches" : "Nothing here";
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
    var cutoffMs = generatedAtMs - windowHours * 3600 * 1000;
    var ids = [];
    (boardData.calendar || []).forEach(function (entry) {
      var ms = Date.parse(entry.completionTime);
      if (!isNaN(ms) && ms > cutoffMs) {
        ids.push(entry.id);
      }
    });
    return ids; // calendar is already most-recent-first
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
  // They carry no honest completion instant, so they are listed here as data
  // bugs to fix — never sorted into Recently done as if completed "now",
  // never aged out by the 24h/48h/7d window, and visible from every view.
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
  // `do-work note` and rendered with textContent, so a stray `<` or a pasted
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

    (boardData.userRequestOrder || []).forEach(function (userRequestId) {
      var userRequest = userRequestsById[userRequestId];
      if (!userRequest) {
        return;
      }
      if (filterState.userRequestActivity === "active" && !userRequestIsActive(userRequest)) {
        hiddenResolvedCount += 1;
        return;
      }

      // A search hit on the UR header keeps the whole group; domain/status
      // still filter the cards inside it.
      var groupMatchesSearch =
        filterState.searchText !== "" &&
        searchMatchesUserRequest(userRequest, userRequestId, filterState.searchText);
      var requestIds = userRequest.requestIds || [];
      var shownRequestIds = requestIds.filter(function (requestId) {
        return requestMatchesFilters(requestId, { skipSearch: groupMatchesSearch });
      });
      if (hasActiveFilters() && shownRequestIds.length === 0) {
        return;
      }

      var group = createElement("section", "ur-group");

      var head = createElement("button", "ur-group-head");
      head.type = "button";
      head.dataset.detailKind = "ur";
      head.dataset.detailId = userRequestId;
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
      group.appendChild(head);

      var cards = createElement("div", "ur-group-cards");
      shownRequestIds.forEach(function (requestId) {
        cards.appendChild(makeRequestCard(requestId, { showCompleted: true }));
      });
      group.appendChild(cards);

      host.appendChild(group);
    });

    if (host.childNodes.length === 0) {
      var emptyText = hasActiveFilters()
        ? "No user requests match the current filters."
        : "No active user requests — every UR is fully resolved. Switch URs to All to browse the archive.";
      host.appendChild(createElement("p", "ur-lens-empty", emptyText));
      return;
    }
    if (hiddenResolvedCount > 0) {
      host.appendChild(
        createElement(
          "p",
          "ur-lens-hidden-note",
          hiddenResolvedCount +
            " fully resolved UR" +
            (hiddenResolvedCount === 1 ? "" : "s") +
            " hidden — switch URs to All to see them."
        )
      );
    }
  }

  // ---- calendar -----------------------------------------------------------

  function renderCalendar() {
    var scroll = document.getElementById("calendar-scroll");
    var summary = document.getElementById("calendar-summary");
    scroll.textContent = "";

    var calendar = boardData.calendar || [];
    var shownEntries = calendar.filter(function (entry) {
      return requestMatchesFilters(entry.id);
    });
    summary.textContent = hasActiveFilters()
      ? shownEntries.length +
        " of " +
        calendar.length +
        " completed REQ" +
        (calendar.length === 1 ? "" : "s") +
        " match the current filters"
      : calendar.length + " completed REQ" + (calendar.length === 1 ? "" : "s") + " across the archive";

    // The calendar is sorted most-recent-first, so equal day keys are
    // contiguous — group by walking the list. Days whose entries are all
    // filtered out never flush, so they disappear entirely.
    var currentDayKey = null;
    var currentEntries = null;

    function flushDay() {
      if (!currentDayKey || !currentEntries || currentEntries.length === 0) {
        return;
      }
      scroll.appendChild(makeCalendarDay(currentDayKey, currentEntries));
    }

    shownEntries.forEach(function (entry) {
      if (entry.dayKey !== currentDayKey) {
        flushDay();
        currentDayKey = entry.dayKey;
        currentEntries = [];
      }
      currentEntries.push(entry);
    });
    flushDay();
  }

  function makeCalendarDay(dayKey, entries) {
    var section = createElement("section", "calendar-day");

    var label = createElement("div", "calendar-day-label");
    var dayDate = new Date(dayKey + "T00:00:00Z");
    if (!isNaN(dayDate.getTime())) {
      label.appendChild(createElement("span", "calendar-day-weekday", calendarWeekdayFormatter.format(dayDate)));
      label.appendChild(createElement("span", "calendar-day-date", calendarDateFormatter.format(dayDate)));
    } else {
      label.appendChild(createElement("span", "calendar-day-date", dayKey));
    }
    label.appendChild(
      createElement("span", "calendar-day-count", entries.length + " done")
    );
    section.appendChild(label);

    var list = createElement("div", "calendar-day-entries");
    entries.forEach(function (entry) {
      var request = requestsById[entry.id];
      var chip = createElement("button", "calendar-chip");
      chip.type = "button";
      chip.dataset.detailKind = "req";
      chip.dataset.detailId = entry.id;
      chip.setAttribute("aria-label", entry.id + (request ? ": " + (request.title || "") : ""));
      chip.appendChild(createElement("span", "calendar-chip-id", entry.id));
      if (request && request.title) {
        chip.appendChild(createElement("span", "calendar-chip-title", request.title));
      }
      list.appendChild(chip);
    });
    section.appendChild(list);

    return section;
  }

  // ---- testing view --------------------------------------------------------
  // Tracks who tested which finished REQ. The REQ Markdown files are the
  // database: actions POST to the live server's /api/testing/* endpoints, which
  // upsert the testing_* placeholder frontmatter fields (and append tester
  // profiles to do-work/testers.md). A static snapshot has no server, so the
  // view renders read-only there (boardData.liveTestingApi is unset).

  var testingLiveApiAvailable = Boolean(boardData.liveTestingApi);
  var testerProfileStorageKey = "queueKanbanTesterProfile";
  var selectedTesterProfile = "";
  var feedbackFormRequestId = null; // REQ id whose card is showing the inline feedback form
  var testingErrorTimer = null;

  function isTerminalSuccessStatus(status) {
    return status === "completed" || status === "completed-with-issues";
  }

  // The instant a testing card is sorted and date-filtered by: the last testing
  // activity when there is one, else the REQ's resolved completion instant.
  // 0 means neither is known.
  function testingRecencyMs(request) {
    var activityMs = Date.parse(request.testingUpdatedAt || "");
    if (isNaN(activityMs)) {
      activityMs = Date.parse(request.completionTime || "");
    }
    return isNaN(activityMs) ? 0 : activityMs;
  }

  function requestIdNumber(requestId) {
    var digitsMatch = /(\d+)/.exec(requestId || "");
    return digitsMatch ? parseInt(digitsMatch[1], 10) : 0;
  }

  // Testing columns read newest-first — with hundreds of finished REQs, the
  // ones just done are the ones a tester is looking for. Cards with no known
  // instant sink to the bottom; the numeric REQ id (higher = newer) breaks ties.
  function sortMostRecentFirst(requestIds) {
    requestIds.sort(function (leftId, rightId) {
      var recencyDelta = testingRecencyMs(requestsById[rightId]) - testingRecencyMs(requestsById[leftId]);
      if (recencyDelta !== 0) {
        return recencyDelta;
      }
      return requestIdNumber(rightId) - requestIdNumber(leftId);
    });
    return requestIds;
  }

  // Testing-view-only date window (the select is hidden on other views and the
  // filter is applied only here, so it can never blank the board's pending
  // columns). Cards with no known instant only show under "Any date".
  function matchesDoneWindow(requestId) {
    if (filterState.doneWindow === "") {
      return true;
    }
    var recencyMs = testingRecencyMs(requestsById[requestId]);
    var thirtyDaysMs = 720 * 3600 * 1000;
    if (filterState.doneWindow === "old") {
      return recencyMs !== 0 && recencyMs <= generatedAtMs - thirtyDaysMs;
    }
    var windowHours = parseInt(filterState.doneWindow, 10);
    return recencyMs > generatedAtMs - windowHours * 3600 * 1000;
  }

  // A REQ belongs on the testing view when it finished successfully (testable)
  // or already carries a testing record (which must never disappear, even if
  // its pipeline status later changed — e.g. returned work re-queued for a fix).
  function computeTestingBuckets() {
    var buckets = { testingReady: [], testingInTesting: [], testingReturned: [], testingTested: [] };
    (boardData.requestOrder || []).forEach(function (requestId) {
      var request = requestsById[requestId];
      if (!request) {
        return;
      }
      // An unrecognized testing_status is still a testing record — it must
      // stay visible (in Ready to test, with the invalid flag) even when the
      // REQ's pipeline status is no longer terminal-success.
      var hasTestingRecord = Boolean(request.testingStatus) || Boolean(request.testingStatusUnrecognized);
      if (!isTerminalSuccessStatus(request.status) && !hasTestingRecord) {
        return;
      }
      if (request.testingStatus === "in-testing") {
        buckets.testingInTesting.push(requestId);
      } else if (request.testingStatus === "returned") {
        buckets.testingReturned.push(requestId);
      } else if (request.testingStatus === "tested") {
        buckets.testingTested.push(requestId);
      } else {
        buckets.testingReady.push(requestId);
      }
    });
    Object.keys(buckets).forEach(function (bucketKey) {
      sortMostRecentFirst(buckets[bucketKey]);
    });
    return buckets;
  }

  function renderTestingView() {
    document.getElementById("testing-readonly-note").hidden = testingLiveApiAvailable;
    var buckets = computeTestingBuckets();
    Object.keys(buckets).forEach(function (bucketKey) {
      var shownIds = filterRequestIds(buckets[bucketKey]).filter(matchesDoneWindow);
      fillTestingColumn(bucketKey, shownIds, buckets[bucketKey].length);
    });
  }

  function fillTestingColumn(columnKey, requestIds, totalCount) {
    var container = document.querySelector('[data-cards="' + columnKey + '"]');
    var countNode = document.querySelector('[data-count="' + columnKey + '"]');
    container.textContent = "";
    countNode.textContent = formatFilteredCount(requestIds.length, totalCount);
    if (requestIds.length === 0) {
      container.appendChild(createElement("p", "column-empty", columnEmptyText()));
      return;
    }
    requestIds.forEach(function (requestId) {
      container.appendChild(makeTestingCard(requestId, columnKey));
    });
  }

  // A testing card wraps the normal REQ card (still opens the detail drawer)
  // with a testing-meta line and an action row. The wrapper is a div — the REQ
  // card itself is a <button>, and buttons must not nest.
  function makeTestingCard(requestId, bucketKey) {
    var request = requestsById[requestId];
    var wrapper = createElement("div", "testing-card");
    wrapper.appendChild(makeRequestCard(requestId, { showCompleted: true }));

    var meta = createElement("div", "testing-card-meta");
    if (request.testedBy) {
      meta.appendChild(createElement("span", "testing-meta-chip", "tester: " + request.testedBy));
    }
    if (request.testingUpdatedAt) {
      meta.appendChild(createElement("span", "testing-meta-chip", formatShortInstant(request.testingUpdatedAt)));
    }
    if (request.testingStatusUnrecognized) {
      var invalidChip = createElement("span", "testing-meta-chip is-invalid", "invalid testing_status");
      invalidChip.title =
        'Unrecognized testing_status "' +
        (request.originalTestingStatus || "") +
        '" — expected in-testing, tested, or returned. Shown as not tested.';
      meta.appendChild(invalidChip);
    }
    if (meta.childNodes.length > 0) {
      wrapper.appendChild(meta);
    }

    if (request.testingStatus === "returned" && request.testingFeedback) {
      wrapper.appendChild(createElement("div", "testing-feedback", request.testingFeedback));
    }

    if (testingLiveApiAvailable) {
      wrapper.appendChild(
        feedbackFormRequestId === requestId
          ? makeTestingFeedbackForm(requestId)
          : makeTestingActionsRow(requestId, bucketKey)
      );
    }
    return wrapper;
  }

  function makeTestingActionsRow(requestId, bucketKey) {
    var actionsRow = createElement("div", "testing-actions");

    function addActionButton(labelText, onActivate, extraClassName) {
      var actionButton = createElement(
        "button",
        "control-button testing-action" + (extraClassName ? " " + extraClassName : ""),
        labelText
      );
      actionButton.type = "button";
      if (!selectedTesterProfile) {
        actionButton.disabled = true;
        actionButton.title = "Select a tester profile first";
      } else {
        actionButton.addEventListener("click", onActivate);
      }
      actionsRow.appendChild(actionButton);
    }

    if (bucketKey === "testingReady") {
      addActionButton("Start testing", function () {
        postTestingStatus(requestId, "in-testing");
      });
    } else if (bucketKey === "testingInTesting") {
      addActionButton("Mark tested", function () {
        postTestingStatus(requestId, "tested");
      });
      addActionButton("Return with feedback", function () {
        feedbackFormRequestId = requestId;
        renderTestingView();
      });
      addActionButton("Clear", function () {
        postTestingStatus(requestId, "clear");
      }, "testing-action-clear");
    } else if (bucketKey === "testingReturned") {
      addActionButton("Restart testing", function () {
        postTestingStatus(requestId, "in-testing");
      });
      addActionButton("Clear", function () {
        postTestingStatus(requestId, "clear");
      }, "testing-action-clear");
    } else if (bucketKey === "testingTested") {
      addActionButton("Re-test", function () {
        postTestingStatus(requestId, "in-testing");
      });
      addActionButton("Clear", function () {
        postTestingStatus(requestId, "clear");
      }, "testing-action-clear");
    }
    return actionsRow;
  }

  function makeTestingFeedbackForm(requestId) {
    var form = createElement("div", "testing-feedback-form");
    var feedbackInput = document.createElement("textarea");
    feedbackInput.className = "testing-feedback-input";
    feedbackInput.rows = 3;
    feedbackInput.placeholder = "What needs fixing?";
    feedbackInput.setAttribute("aria-label", "Feedback for " + requestId);
    form.appendChild(feedbackInput);

    var formActions = createElement("div", "testing-actions");
    var confirmButton = createElement("button", "control-button testing-action", "Return");
    confirmButton.type = "button";
    confirmButton.addEventListener("click", function () {
      var feedbackText = feedbackInput.value.trim();
      if (feedbackText === "") {
        showTestingError("Feedback must not be empty — describe what to fix.");
        return;
      }
      feedbackFormRequestId = null;
      postTestingStatus(requestId, "returned", feedbackText);
    });
    var cancelButton = createElement("button", "control-button testing-action", "Cancel");
    cancelButton.type = "button";
    cancelButton.addEventListener("click", function () {
      feedbackFormRequestId = null;
      renderTestingView();
    });
    formActions.appendChild(confirmButton);
    formActions.appendChild(cancelButton);
    form.appendChild(formActions);

    setTimeout(function () {
      feedbackInput.focus();
    }, 0);
    return form;
  }

  function showTestingError(errorText) {
    var errorNode = document.getElementById("testing-error");
    errorNode.textContent = errorText;
    errorNode.hidden = false;
    if (testingErrorTimer) {
      clearTimeout(testingErrorTimer);
    }
    testingErrorTimer = setTimeout(function () {
      errorNode.hidden = true;
    }, 8000);
  }

  function decodeTestingApiResponse(httpResponse) {
    return httpResponse
      .json()
      .catch(function () {
        return { ok: false, error: "HTTP " + httpResponse.status };
      })
      .then(function (payload) {
        if (!httpResponse.ok || !payload.ok) {
          throw new Error(payload.error || "HTTP " + httpResponse.status);
        }
        return payload;
      });
  }

  // On success the server's confirmed transition is applied to the local data
  // island and the view re-renders — no page reload, so the active view and
  // filters survive. The next full reload re-reads the files themselves.
  function postTestingStatus(requestId, testingState, feedbackText) {
    fetch("/api/testing/status", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        requestId: requestId,
        testingStatus: testingState,
        testedBy: selectedTesterProfile,
        feedback: feedbackText || ""
      })
    })
      .then(decodeTestingApiResponse)
      .then(function (payload) {
        var request = requestsById[requestId];
        if (request) {
          request.testingStatus = payload.testingStatus || "";
          request.testedBy = payload.testedBy || "";
          request.testingUpdatedAt = payload.testingUpdatedAt || "";
          request.testingFeedback = testingState === "returned" ? feedbackText || "" : "";
          request.testingStatusUnrecognized = false;
          request.originalTestingStatus = payload.testingStatus || "";
        }
        renderTestingView();
        renderColumns(); // the main board's testing badge tracks the same record
      })
      .catch(function (postError) {
        showTestingError("Update failed: " + postError.message);
        renderTestingView();
      });
  }

  function populateTestingProfileSelect() {
    var profileSelect = document.getElementById("testing-profile-select");
    while (profileSelect.options.length > 1) {
      profileSelect.remove(1);
    }
    fillSelectOptions(profileSelect, boardData.testingProfiles || []);

    var storedProfile = "";
    try {
      storedProfile = localStorage.getItem(testerProfileStorageKey) || "";
    } catch (storageError) {
      // Persistence is best-effort.
    }
    if (storedProfile && (boardData.testingProfiles || []).indexOf(storedProfile) !== -1) {
      selectedTesterProfile = storedProfile;
      profileSelect.value = storedProfile;
    }
  }

  function wireTestingControls() {
    var profileSelect = document.getElementById("testing-profile-select");
    profileSelect.addEventListener("change", function () {
      selectedTesterProfile = profileSelect.value;
      try {
        localStorage.setItem(testerProfileStorageKey, selectedTesterProfile);
      } catch (storageError) {
        // Persistence is best-effort.
      }
      renderTestingView();
    });

    var addToggleButton = document.getElementById("testing-profile-add-toggle");
    var addForm = document.getElementById("testing-profile-add-form");
    var addNameInput = document.getElementById("testing-profile-add-name");

    if (!testingLiveApiAvailable) {
      addToggleButton.disabled = true;
      addToggleButton.title = "Adding testers needs the live board (do-work board)";
    }

    addToggleButton.addEventListener("click", function () {
      addForm.hidden = false;
      addToggleButton.hidden = true;
      addNameInput.focus();
    });

    function closeAddForm() {
      addForm.hidden = true;
      addToggleButton.hidden = false;
      addNameInput.value = "";
    }

    document.getElementById("testing-profile-add-cancel").addEventListener("click", closeAddForm);

    function submitNewProfile() {
      var profileName = addNameInput.value.trim();
      if (profileName === "") {
        showTestingError("Tester name must not be empty.");
        return;
      }
      fetch("/api/testing/profile", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ name: profileName })
      })
        .then(decodeTestingApiResponse)
        .then(function (payload) {
          boardData.testingProfiles = payload.profiles || [];
          populateTestingProfileSelect();
          // Select the just-added profile (the server may return the existing
          // spelling when the name was already known, so match case-insensitively).
          var matchedProfile = (boardData.testingProfiles || []).filter(function (knownProfile) {
            return knownProfile.toLowerCase() === profileName.toLowerCase();
          })[0];
          if (matchedProfile) {
            selectedTesterProfile = matchedProfile;
            profileSelect.value = matchedProfile;
            try {
              localStorage.setItem(testerProfileStorageKey, matchedProfile);
            } catch (storageError) {
              // Persistence is best-effort.
            }
          }
          closeAddForm();
          renderTestingView();
        })
        .catch(function (postError) {
          showTestingError("Could not add tester: " + postError.message);
        });
    }

    document.getElementById("testing-profile-add-confirm").addEventListener("click", submitNewProfile);
    addNameInput.addEventListener("keydown", function (keyEvent) {
      if (keyEvent.key === "Enter") {
        keyEvent.preventDefault();
        submitNewProfile();
      } else if (keyEvent.key === "Escape") {
        keyEvent.preventDefault();
        closeAddForm();
      }
    });
  }

  // ---- detail panel (docked beside the board, non-modal) -------------------

  var detailResizer = document.getElementById("detail-resizer");
  var drawer = document.getElementById("detail-drawer");
  var drawerKind = document.getElementById("detail-kind");
  var drawerId = document.getElementById("detail-id");
  var drawerTitle = document.getElementById("detail-drawer-title");
  var drawerMeta = document.getElementById("detail-meta");
  var drawerBody = document.getElementById("detail-body");
  var drawerCopyButton = document.getElementById("detail-copy");
  var lastFocusedElement = null;

  // The initial board payload deliberately omits raw Markdown. Remember only
  // the open record's identity; the Copy button loads board-markdown.js on its
  // first use, then looks up the exact source by kind + id.
  var currentDetailKind = "";
  var currentDetailId = "";

  function appendMetaRow(label, valueNode) {
    var dt = createElement("dt", null, label);
    var dd = createElement("dd");
    if (typeof valueNode === "string") {
      dd.textContent = valueNode;
    } else {
      dd.appendChild(valueNode);
    }
    drawerMeta.appendChild(dt);
    drawerMeta.appendChild(dd);
  }

  // Each dependency listed with the status that decides whether it is met, so
  // "why is this still waiting?" is answerable without opening the upstream REQ.
  function makeDependencyDetailList(request) {
    var unmetDependencyIds = request.unmetDependencies || [];
    var list = createElement("div", "detail-dep-list");
    request.dependsOn.forEach(function (dependencyId) {
      var isUnmet = unmetDependencyIds.indexOf(dependencyId) !== -1;
      var row = createElement("span", isUnmet ? "detail-dep is-unmet" : "detail-dep is-met");
      row.appendChild(createElement("span", "detail-dep-id", dependencyId));
      row.appendChild(createElement("span", "detail-dep-status", describeRequestStatus(dependencyId)));
      list.appendChild(row);
    });
    return list;
  }

  function openRequestDetail(requestId) {
    var request = requestsById[requestId];
    if (!request) {
      return;
    }
    drawerKind.textContent = "REQ";
    drawerId.textContent = requestId;
    drawerTitle.textContent = request.title || "untitled";

    drawerMeta.textContent = "";
    if (request.statusUnrecognized) {
      var invalidStatus = createElement("span", "detail-status-invalid");
      invalidStatus.appendChild(document.createTextNode(request.originalStatus || request.status || "—"));
      invalidStatus.appendChild(createElement("span", "status-invalid-flag", "invalid"));
      appendMetaRow("Status", invalidStatus);
      appendMetaRow(
        "Fix",
        "This status is not in the schema vocabulary, so the ticket is parked under Needs input / Blocked. " +
          "Edit the REQ's status: field to a recognized value (actions/work-reference.md → Schema Read Contract) " +
          "or run do-work forensics to sweep the tree for invalid statuses."
      );
    } else {
      appendMetaRow("Status", request.originalStatus || request.status || "—");
    }
    if (request.domain) {
      appendMetaRow("Domain", request.domain);
    }
    if (request.userRequestId) {
      var urLink = createElement("button", "control-button", request.userRequestId);
      urLink.type = "button";
      urLink.dataset.detailKind = "ur";
      urLink.dataset.detailId = request.userRequestId;
      urLink.style.padding = "2px 10px";
      appendMetaRow("User request", urLink);
    }
    if (request.dependsOn && request.dependsOn.length > 0) {
      appendMetaRow("Depends on", makeDependencyDetailList(request));
    }
    var unblockedRequestIds = activeDependentIds(request);
    if (unblockedRequestIds.length > 0) {
      appendMetaRow("Unblocks", unblockedRequestIds.join(", "));
    }
    if (request.route) {
      appendMetaRow("Route", request.route);
    }
    if (request.createdAt) {
      appendMetaRow("Created", request.createdAt);
    }
    if (request.completionTime) {
      appendMetaRow("Completed", formatShortInstant(request.completionTime) + " (" + request.completionTimeSource + ")");
    }
    if (request.completionAnomaly) {
      var anomalyValue = createElement("span", "detail-status-invalid");
      anomalyValue.appendChild(
        document.createTextNode(request.completionAnomalyReason || "completion instant unresolved")
      );
      anomalyValue.appendChild(createElement("span", "status-invalid-flag", "anomaly"));
      appendMetaRow("Completion anomaly", anomalyValue);
      appendMetaRow(
        "Fix",
        "Add completed_at: <ISO instant> (e.g. 2026-07-18T12:00:00Z) and/or a commit: field holding the " +
          "implementation commit hash to this REQ's frontmatter."
      );
    }
    if (request.testingStatus || request.testingStatusUnrecognized) {
      var testingSummary = request.testingStatusUnrecognized
        ? (request.originalTestingStatus || "?") + " (invalid — expected in-testing, tested, or returned)"
        : request.testingStatus;
      if (request.testedBy) {
        testingSummary += " — " + request.testedBy;
      }
      appendMetaRow("Testing", testingSummary);
      if (request.testingUpdatedAt) {
        appendMetaRow("Testing updated", formatShortInstant(request.testingUpdatedAt) || request.testingUpdatedAt);
      }
      if (request.testingFeedback) {
        appendMetaRow("Testing feedback", request.testingFeedback);
      }
    }
    appendMetaRow("Tree", request.treeSection || "—");

    drawerBody.innerHTML = request.bodyHtml || "<p>(empty body)</p>";
    currentDetailKind = "req";
    currentDetailId = requestId;
    showDrawer();
  }

  function openUserRequestDetail(userRequestId) {
    var userRequest = userRequestsById[userRequestId];
    if (!userRequest) {
      return;
    }
    drawerKind.textContent = "UR";
    drawerId.textContent = userRequestId;
    drawerTitle.textContent = userRequest.title || "(no input.md title)";

    drawerMeta.textContent = "";
    var requestIds = userRequest.requestIds || [];
    appendMetaRow("Grouped REQs", String(requestIds.length));
    if (requestIds.length > 0) {
      appendMetaRow("REQ ids", requestIds.join(", "));
    }
    appendMetaRow("input.md", userRequest.inputFilePresent ? "present" : "synthesized from REQ pointers");

    drawerBody.innerHTML = userRequest.bodyHtml || "<p>(no input.md body)</p>";
    currentDetailKind = "ur";
    currentDetailId = userRequestId;
    showDrawer();
  }

  function showDrawer() {
    // Only capture the return-focus target on a TRUE first open (drawer currently
    // hidden). A REQ drawer can re-enter showDrawer() while already open — e.g. its
    // inner "User request" button navigates to the UR detail — and overwriting
    // lastFocusedElement there would lose the originating card, so closing would fail
    // to restore focus to the trigger that opened the drawer.
    if (drawer.hidden) {
      lastFocusedElement = document.activeElement;
    }
    drawer.hidden = false;
    detailResizer.hidden = false;
    // A lingering "Copied ✓" from the previous ticket would misreport what is
    // on the clipboard — reset the button on every open.
    drawerCopyButton.textContent = "Copy";
    drawerCopyButton.classList.remove("is-copied", "is-copy-failed");
    drawerBody.scrollTop = 0;
    drawer.scrollTop = 0;
    drawer.focus();
    document.addEventListener("keydown", onDetailPanelKeydown, true);
  }

  function closeDrawer() {
    if (drawer.hidden) {
      return;
    }
    drawer.hidden = true;
    detailResizer.hidden = true;
    currentDetailKind = "";
    currentDetailId = "";
    document.removeEventListener("keydown", onDetailPanelKeydown, true);
    if (lastFocusedElement && typeof lastFocusedElement.focus === "function") {
      lastFocusedElement.focus();
    }
  }

  // The panel is docked, not modal — the board stays interactive, so there is
  // no focus trap and no scrim. Escape still dismisses it from anywhere.
  function onDetailPanelKeydown(keyEvent) {
    if (keyEvent.key === "Escape") {
      keyEvent.preventDefault();
      closeDrawer();
    }
  }

  function openDetail(kind, id) {
    if (kind === "ur") {
      openUserRequestDetail(id);
    } else {
      openRequestDetail(id);
    }
  }

  // ---- detail panel resizing ------------------------------------------------
  // The divider drags like Jira's issue split view: pointer capture for
  // mouse/touch, arrow keys while focused, double-click to reset. The width
  // lives in the --detail-panel-width custom property so CSS grid does the
  // layout, and it persists across reloads via localStorage (best-effort —
  // a denied storage context only loses persistence, never the resize).

  var detailPanelDefaultWidthPx = 620;
  var detailPanelMinWidthPx = 360; // mirrored by the clamp() in board.css
  var boardMinVisibleWidthPx = 340; // never let the panel push the board below this
  var detailPanelWidthStorageKey = "queueKanbanDetailPanelWidthPx";
  var detailResizeState = null;

  function applyDetailPanelWidth(candidateWidthPx) {
    var maxWidthPx = Math.max(detailPanelMinWidthPx, window.innerWidth - boardMinVisibleWidthPx);
    var clampedWidthPx = Math.min(Math.max(candidateWidthPx, detailPanelMinWidthPx), maxWidthPx);
    document.documentElement.style.setProperty("--detail-panel-width", clampedWidthPx + "px");
    detailResizer.setAttribute("aria-valuenow", String(Math.round(clampedWidthPx)));
    detailResizer.setAttribute("aria-valuemax", String(Math.round(maxWidthPx)));
    return clampedWidthPx;
  }

  function persistDetailPanelWidth(widthPx) {
    try {
      localStorage.setItem(detailPanelWidthStorageKey, String(Math.round(widthPx)));
    } catch (storageError) {
      // Persistence is best-effort; the in-page resize already applied.
    }
  }

  (function restoreDetailPanelWidth() {
    var storedWidthPx = NaN;
    try {
      storedWidthPx = parseFloat(localStorage.getItem(detailPanelWidthStorageKey));
    } catch (storageError) {
      // Fall through to the stylesheet default.
    }
    if (!isNaN(storedWidthPx)) {
      applyDetailPanelWidth(storedWidthPx);
    }
  })();

  detailResizer.addEventListener("pointerdown", function (pointerEvent) {
    detailResizeState = {
      pointerId: pointerEvent.pointerId,
      startClientX: pointerEvent.clientX,
      startWidthPx: drawer.getBoundingClientRect().width
    };
    detailResizer.setPointerCapture(pointerEvent.pointerId);
    document.body.classList.add("is-resizing-detail");
    pointerEvent.preventDefault();
  });

  detailResizer.addEventListener("pointermove", function (pointerEvent) {
    if (!detailResizeState || pointerEvent.pointerId !== detailResizeState.pointerId) {
      return;
    }
    // The panel sits on the right, so dragging the divider left grows it.
    applyDetailPanelWidth(
      detailResizeState.startWidthPx + (detailResizeState.startClientX - pointerEvent.clientX)
    );
  });

  function endDetailPanelResize(pointerEvent) {
    if (!detailResizeState || pointerEvent.pointerId !== detailResizeState.pointerId) {
      return;
    }
    detailResizeState = null;
    document.body.classList.remove("is-resizing-detail");
    persistDetailPanelWidth(drawer.getBoundingClientRect().width);
  }
  detailResizer.addEventListener("pointerup", endDetailPanelResize);
  detailResizer.addEventListener("pointercancel", endDetailPanelResize);

  detailResizer.addEventListener("dblclick", function () {
    persistDetailPanelWidth(applyDetailPanelWidth(detailPanelDefaultWidthPx));
  });

  detailResizer.addEventListener("keydown", function (keyEvent) {
    var stepPx = keyEvent.shiftKey ? 64 : 16;
    var currentWidthPx = drawer.getBoundingClientRect().width;
    var nextWidthPx = null;
    if (keyEvent.key === "ArrowLeft") {
      nextWidthPx = currentWidthPx + stepPx; // divider moves left → panel grows
    } else if (keyEvent.key === "ArrowRight") {
      nextWidthPx = currentWidthPx - stepPx;
    } else if (keyEvent.key === "Home") {
      nextWidthPx = detailPanelMinWidthPx;
    } else if (keyEvent.key === "End") {
      nextWidthPx = window.innerWidth; // applyDetailPanelWidth clamps to the max
    }
    if (nextWidthPx !== null) {
      keyEvent.preventDefault();
      persistDetailPanelWidth(applyDetailPanelWidth(nextWidthPx));
    }
  });

  // ---- view / lens / window switching ------------------------------------

  function setActiveButton(groupSelector, attributeName, value) {
    var buttons = document.querySelectorAll(groupSelector + " [" + attributeName + "]");
    buttons.forEach(function (button) {
      var isActive = button.getAttribute(attributeName) === value;
      button.classList.toggle("is-active", isActive);
      button.setAttribute("aria-pressed", isActive ? "true" : "false");
    });
  }

  function applyView() {
    var viewPanels = {
      board: document.getElementById("view-board"),
      calendar: document.getElementById("view-calendar"),
      testing: document.getElementById("view-testing")
    };
    Object.keys(viewPanels).forEach(function (viewName) {
      var isActiveView = viewState.view === viewName;
      viewPanels[viewName].classList.toggle("is-active", isActiveView);
      viewPanels[viewName].hidden = !isActiveView;
    });

    // The grouping lens and the recently-done window only shape the board view;
    // hide their controls elsewhere so the topbar never advertises dead knobs.
    // The date window is the testing view's knob for the same reason.
    document.getElementById("lens-group").hidden = viewState.view !== "board";
    document.getElementById("recent-window-group").hidden = viewState.view !== "board";
    document.getElementById("filter-done-window").hidden = viewState.view !== "testing";

    if (viewState.view === "calendar" && !renderedOnce.calendar) {
      renderCalendar();
      renderedOnce.calendar = true;
    }
    if (viewState.view === "testing" && !renderedOnce.testing) {
      renderTestingView();
      renderedOnce.testing = true;
    }
    if (viewState.view === "board") {
      applyLens();
    } else {
      updateUserRequestActivityVisibility();
    }
  }

  function applyLens() {
    var columns = document.getElementById("kanban-columns");
    var lensHost = document.getElementById("user-request-lens");
    var byUserRequest = viewState.lens === "user-request";

    columns.hidden = byUserRequest;
    lensHost.hidden = !byUserRequest;
    updateUserRequestActivityVisibility();

    if (byUserRequest && !renderedOnce.userRequestLens) {
      renderUserRequestLens();
      renderedOnce.userRequestLens = true;
    }
  }

  // The Active/All toggle only means something on the by-UR lens — hide it
  // everywhere else so the topbar doesn't advertise a dead control.
  function updateUserRequestActivityVisibility() {
    document.getElementById("ur-activity-group").hidden =
      viewState.view !== "board" || viewState.lens !== "user-request";
  }

  function wireControls() {
    document.querySelectorAll("[data-view-target]").forEach(function (button) {
      button.addEventListener("click", function () {
        viewState.view = button.getAttribute("data-view-target");
        setActiveButton("[aria-label='Board views and lenses']", "data-view-target", viewState.view);
        applyView();
      });
    });

    document.querySelectorAll("[data-lens-target]").forEach(function (button) {
      button.addEventListener("click", function () {
        viewState.lens = button.getAttribute("data-lens-target");
        setActiveButton("#lens-group", "data-lens-target", viewState.lens);
        if (viewState.view === "board") {
          applyLens();
        }
      });
    });

    document.querySelectorAll("[data-window-hours]").forEach(function (button) {
      button.addEventListener("click", function () {
        viewState.windowHours = parseInt(button.getAttribute("data-window-hours"), 10) || 24;
        setActiveButton("#recent-window-group", "data-window-hours", String(viewState.windowHours));
        renderColumns();
      });
    });

    var searchInput = document.getElementById("filter-search");
    searchInput.addEventListener("input", function () {
      filterState.searchText = searchInput.value.trim().toLowerCase();
      onFiltersChanged();
    });

    var domainSelect = document.getElementById("filter-domain");
    domainSelect.addEventListener("change", function () {
      filterState.domain = domainSelect.value;
      onFiltersChanged();
    });

    var statusSelect = document.getElementById("filter-status");
    statusSelect.addEventListener("change", function () {
      filterState.status = statusSelect.value;
      onFiltersChanged();
    });

    var doneWindowSelect = document.getElementById("filter-done-window");
    doneWindowSelect.addEventListener("change", function () {
      filterState.doneWindow = doneWindowSelect.value;
      onFiltersChanged();
    });

    document.getElementById("filter-clear").addEventListener("click", function () {
      filterState.searchText = "";
      filterState.domain = "";
      filterState.status = "";
      filterState.doneWindow = "";
      searchInput.value = "";
      domainSelect.value = "";
      statusSelect.value = "";
      doneWindowSelect.value = "";
      onFiltersChanged();
    });

    document.querySelectorAll("[data-ur-activity]").forEach(function (button) {
      button.addEventListener("click", function () {
        filterState.userRequestActivity = button.getAttribute("data-ur-activity");
        setActiveButton("#ur-activity-group", "data-ur-activity", filterState.userRequestActivity);
        renderedOnce.userRequestLens = false;
        if (viewState.view === "board" && viewState.lens === "user-request") {
          renderUserRequestLens();
          renderedOnce.userRequestLens = true;
        }
      });
    });
  }

  // Event delegation: any element carrying data-detail-kind opens the drawer.
  document.addEventListener("click", function (clickEvent) {
    var trigger = clickEvent.target.closest("[data-detail-kind]");
    if (trigger && !trigger.disabled) {
      clickEvent.preventDefault();
      openDetail(trigger.getAttribute("data-detail-kind"), trigger.getAttribute("data-detail-id"));
    }
  });

  document.getElementById("detail-close").addEventListener("click", closeDrawer);

  // ---- copy-to-clipboard for the open ticket --------------------------------
  // Raw source is a separate sibling script so the initial board does not pay
  // for both rendered HTML and Markdown. Loading a script instead of fetch()
  // preserves direct file:// use. Clipboard writes prefer the async API and
  // fall back to a hidden textarea when permission or protocol blocks it.

  var boardMarkdownLoadPromise = null;

  function loadBoardMarkdownData() {
    if (window.queueKanbanBoardMarkdownData) {
      return Promise.resolve(window.queueKanbanBoardMarkdownData);
    }
    if (boardMarkdownLoadPromise) {
      return boardMarkdownLoadPromise;
    }

    boardMarkdownLoadPromise = new Promise(function (resolve, reject) {
      var markdownScript = document.createElement("script");
      markdownScript.src = "board-markdown.js";
      markdownScript.onload = function () {
        var markdownData = window.queueKanbanBoardMarkdownData;
        if (markdownData && typeof markdownData === "object") {
          resolve(markdownData);
        } else {
          reject(new Error("board-markdown.js did not define Markdown data"));
        }
      };
      markdownScript.onerror = function () {
        reject(new Error("board-markdown.js could not be loaded"));
      };
      document.head.appendChild(markdownScript);
    }).catch(function (loadError) {
      // A generated bundle may have been copied without its lazy sibling.
      // Clear the promise so a later click can retry after the file appears.
      boardMarkdownLoadPromise = null;
      throw loadError;
    });

    return boardMarkdownLoadPromise;
  }

  function rawMarkdownForDetail(markdownData, detailKind, detailId) {
    var markdownById = detailKind === "ur" ? markdownData.userRequests : markdownData.requests;
    if (markdownById && Object.prototype.hasOwnProperty.call(markdownById, detailId)) {
      return markdownById[detailId];
    }
    return null;
  }

  function writeTextToClipboard(clipboardText) {
    if (navigator.clipboard && typeof navigator.clipboard.writeText === "function") {
      return navigator.clipboard.writeText(clipboardText).catch(function () {
        return writeTextViaHiddenTextarea(clipboardText);
      });
    }
    return writeTextViaHiddenTextarea(clipboardText);
  }

  function writeTextViaHiddenTextarea(clipboardText) {
    return new Promise(function (resolve, reject) {
      var scratchTextarea = document.createElement("textarea");
      scratchTextarea.value = clipboardText;
      scratchTextarea.setAttribute("readonly", "");
      scratchTextarea.style.position = "fixed";
      scratchTextarea.style.opacity = "0";
      document.body.appendChild(scratchTextarea);
      scratchTextarea.select();
      var copySucceeded = false;
      try {
        copySucceeded = document.execCommand("copy");
      } catch (execError) {
        copySucceeded = false;
      }
      document.body.removeChild(scratchTextarea);
      if (copySucceeded) {
        resolve();
      } else {
        reject(new Error("execCommand copy failed"));
      }
    });
  }

  var copyFeedbackTimer = null;

  function showCopyFeedback(labelText, stateClass) {
    drawerCopyButton.textContent = labelText;
    drawerCopyButton.classList.remove("is-copied", "is-copy-failed");
    drawerCopyButton.classList.add(stateClass);
    if (copyFeedbackTimer) {
      clearTimeout(copyFeedbackTimer);
    }
    copyFeedbackTimer = setTimeout(function () {
      drawerCopyButton.textContent = "Copy";
      drawerCopyButton.classList.remove("is-copied", "is-copy-failed");
    }, 1600);
  }

  drawerCopyButton.addEventListener("click", function () {
    var requestedKind = currentDetailKind;
    var requestedId = currentDetailId;
    var renderedTextFallback = drawerBody.innerText || "";

    drawerCopyButton.textContent = "Copying…";
    drawerCopyButton.classList.remove("is-copied", "is-copy-failed");

    loadBoardMarkdownData()
      .then(
        function (markdownData) {
          var rawMarkdown = rawMarkdownForDetail(markdownData, requestedKind, requestedId);
          return rawMarkdown === null ? renderedTextFallback : rawMarkdown;
        },
        function () {
          // Keep Copy useful for stale/incomplete generated bundles that lack
          // board-markdown.js, while current bundles retain exact source text.
          return renderedTextFallback;
        }
      )
      .then(writeTextToClipboard)
      .then(
        function () {
          if (!drawer.hidden && currentDetailKind === requestedKind && currentDetailId === requestedId) {
            showCopyFeedback("Copied ✓", "is-copied");
          }
        },
        function () {
          if (!drawer.hidden && currentDetailKind === requestedKind && currentDetailId === requestedId) {
            showCopyFeedback("Copy failed", "is-copy-failed");
          }
        }
      );
  });

  // ---- boot ---------------------------------------------------------------

  wireControls();
  wireTestingControls();
  populateFilterSelects();
  populateTestingProfileSelect();
  renderWarningsBanner();
  renderAnomaliesStrip();
  renderNotesStrip();
  renderColumns();
  applyView();
})();
