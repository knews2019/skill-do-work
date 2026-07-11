/* ===========================================================================
   queue-kanban — static board behaviour
   Reads the embedded JSON data island and renders every view client-side, with
   zero network. No framework: plain DOM construction, event delegation, and a
   docked detail panel with a drag-to-resize divider.
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
    view: "board", // "board" | "calendar"
    lens: "flat", // "flat" | "user-request"
    windowHours: 24
  };

  // Shared filters — applied to whichever view is active. userRequestActivity
  // only affects the by-UR lens ("active" hides URs whose REQs are all resolved).
  var filterState = {
    searchText: "",
    domain: "",
    status: "",
    userRequestActivity: "active" // "active" | "all"
  };

  var renderedOnce = { userRequestLens: false, calendar: false };

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
    return filterState.searchText !== "" || filterState.domain !== "" || filterState.status !== "";
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
    if (viewState.view === "calendar") {
      renderCalendar();
      renderedOnce.calendar = true;
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
    var unblockedRequestIds = activeDependentIds(request);
    if (unblockedRequestIds.length > 0 && !isTerminalResolvedStatus(request.status)) {
      var unblocksBadge = makeBadge("badge-unblocks", "unblocks", String(unblockedRequestIds.length));
      unblocksBadge.title = "Unblocks " + unblockedRequestIds.join(", ") + " when this lands";
      badges.appendChild(unblocksBadge);
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

  // The raw Markdown of whatever the drawer currently shows, set on every
  // open. Copying hands over the ticket's source text, not the rendered HTML,
  // so a paste into chat/email/another REQ keeps headings, checkboxes, links.
  var currentDetailMarkdown = "";

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
    appendMetaRow("Tree", request.treeSection || "—");

    drawerBody.innerHTML = request.bodyHtml || "<p>(empty body)</p>";
    currentDetailMarkdown = request.bodyMarkdown || "";
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
    currentDetailMarkdown = userRequest.bodyMarkdown || "";
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
    var boardPanel = document.getElementById("view-board");
    var calendarPanel = document.getElementById("view-calendar");
    var showCalendar = viewState.view === "calendar";

    boardPanel.classList.toggle("is-active", !showCalendar);
    calendarPanel.classList.toggle("is-active", showCalendar);
    calendarPanel.hidden = !showCalendar;
    boardPanel.hidden = showCalendar;

    if (showCalendar && !renderedOnce.calendar) {
      renderCalendar();
      renderedOnce.calendar = true;
    }
    if (!showCalendar) {
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

    document.getElementById("filter-clear").addEventListener("click", function () {
      filterState.searchText = "";
      filterState.domain = "";
      filterState.status = "";
      searchInput.value = "";
      domainSelect.value = "";
      statusSelect.value = "";
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
  // Prefers the async Clipboard API; falls back to a hidden textarea +
  // execCommand when the API is unavailable (the static board is often opened
  // over plain file:// or http://) or when it rejects (permission denied).

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
    // Stale board-data.js (generated before bodyMarkdown shipped) carries no
    // raw source — degrade to the rendered text so the button still works.
    var clipboardText = currentDetailMarkdown || drawerBody.innerText || "";
    writeTextToClipboard(clipboardText).then(
      function () {
        showCopyFeedback("Copied ✓", "is-copied");
      },
      function () {
        showCopyFeedback("Copy failed", "is-copy-failed");
      }
    );
  });

  // ---- boot ---------------------------------------------------------------

  wireControls();
  populateFilterSelects();
  renderWarningsBanner();
  renderNotesStrip();
  renderColumns();
  applyView();
})();
