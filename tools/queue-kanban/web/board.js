/* ===========================================================================
   queue-kanban — static board behaviour
   Reads the embedded JSON data island and renders every view client-side, with
   zero network. No framework: plain DOM construction, event delegation, and a
   small focus-managed detail drawer.
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

    var top = createElement("div", "req-card-top");
    top.appendChild(createElement("span", "req-card-id", requestId));
    var status = createElement("span", "req-card-status");
    status.appendChild(createElement("span", "status-dot"));
    status.appendChild(document.createTextNode(request.status || "—"));
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
    if (badges.childNodes.length > 0) {
      card.appendChild(badges);
    }

    if (request.dependsOn && request.dependsOn.length > 0) {
      var deps = createElement("div", "req-card-deps");
      deps.appendChild(createElement("span", "dep-chip-lead", "needs"));
      request.dependsOn.forEach(function (dependencyId) {
        deps.appendChild(createElement("span", "dep-chip", dependencyId));
      });
      card.appendChild(deps);
    }

    if (options && options.showCompleted && request.completionTime) {
      card.appendChild(
        createElement("div", "req-card-completed", "done " + formatShortInstant(request.completionTime))
      );
    }

    return card;
  }

  function fillColumn(columnKey, requestIds, options) {
    var container = document.querySelector('[data-cards="' + columnKey + '"]');
    var countNode = document.querySelector('[data-count="' + columnKey + '"]');
    container.textContent = "";
    if (countNode) {
      countNode.textContent = String(requestIds.length);
    }
    if (requestIds.length === 0) {
      container.appendChild(createElement("p", "column-empty", "Nothing here"));
      return;
    }
    requestIds.forEach(function (requestId) {
      container.appendChild(makeRequestCard(requestId, options));
    });
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

  function renderColumns() {
    var columns = boardData.columns || {};
    fillColumn("pending", columns.pending || []);
    fillColumn("claimed", columns.claimed || []);
    fillColumn("needsInputOrBlocked", columns.needsInputOrBlocked || []);
    fillColumn("recentlyDone", recentlyDoneIds(viewState.windowHours), { showCompleted: true });
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

  // ---- by-UR lens ---------------------------------------------------------

  function renderUserRequestLens() {
    var host = document.getElementById("user-request-lens");
    host.textContent = "";
    (boardData.userRequestOrder || []).forEach(function (userRequestId) {
      var userRequest = userRequestsById[userRequestId];
      if (!userRequest) {
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
      var requestIds = userRequest.requestIds || [];
      head.appendChild(createElement("span", "ur-count", requestIds.length + " REQ"));
      group.appendChild(head);

      var cards = createElement("div", "ur-group-cards");
      requestIds.forEach(function (requestId) {
        cards.appendChild(makeRequestCard(requestId, { showCompleted: true }));
      });
      group.appendChild(cards);

      host.appendChild(group);
    });
  }

  // ---- calendar -----------------------------------------------------------

  function renderCalendar() {
    var scroll = document.getElementById("calendar-scroll");
    var summary = document.getElementById("calendar-summary");
    scroll.textContent = "";

    var calendar = boardData.calendar || [];
    summary.textContent =
      calendar.length + " completed REQ" + (calendar.length === 1 ? "" : "s") + " across the archive";

    // The calendar is sorted most-recent-first, so equal day keys are
    // contiguous — group by walking the list.
    var currentDayKey = null;
    var currentEntries = null;

    function flushDay() {
      if (!currentDayKey || !currentEntries) {
        return;
      }
      scroll.appendChild(makeCalendarDay(currentDayKey, currentEntries));
    }

    calendar.forEach(function (entry) {
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

  // ---- detail drawer ------------------------------------------------------

  var overlay = document.getElementById("detail-overlay");
  var drawer = document.getElementById("detail-drawer");
  var drawerKind = document.getElementById("detail-kind");
  var drawerId = document.getElementById("detail-id");
  var drawerTitle = document.getElementById("detail-drawer-title");
  var drawerMeta = document.getElementById("detail-meta");
  var drawerBody = document.getElementById("detail-body");
  var lastFocusedElement = null;

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

  function openRequestDetail(requestId) {
    var request = requestsById[requestId];
    if (!request) {
      return;
    }
    drawerKind.textContent = "REQ";
    drawerId.textContent = requestId;
    drawerTitle.textContent = request.title || "untitled";

    drawerMeta.textContent = "";
    appendMetaRow("Status", request.originalStatus || request.status || "—");
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
      appendMetaRow("Depends on", request.dependsOn.join(", "));
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
    overlay.hidden = false;
    drawer.hidden = false;
    drawerBody.scrollTop = 0;
    drawer.scrollTop = 0;
    document.body.style.overflow = "hidden";
    drawer.focus();
    document.addEventListener("keydown", onDrawerKeydown, true);
  }

  function closeDrawer() {
    if (drawer.hidden) {
      return;
    }
    overlay.hidden = true;
    drawer.hidden = true;
    document.body.style.overflow = "";
    document.removeEventListener("keydown", onDrawerKeydown, true);
    if (lastFocusedElement && typeof lastFocusedElement.focus === "function") {
      lastFocusedElement.focus();
    }
  }

  function onDrawerKeydown(keyEvent) {
    if (keyEvent.key === "Escape") {
      keyEvent.preventDefault();
      closeDrawer();
      return;
    }
    if (keyEvent.key === "Tab") {
      trapFocus(keyEvent);
    }
  }

  function trapFocus(keyEvent) {
    // Disabled inputs must be excluded like disabled buttons: goldmark renders
    // GFM task-list items ("- [ ]") as disabled checkbox inputs, and a disabled
    // element can never be document.activeElement — if it were selected as
    // first/last, the wrap condition would never fire and Tab would escape the
    // aria-modal drawer.
    var focusable = drawer.querySelectorAll(
      'a[href], button:not([disabled]), input:not([disabled]), [tabindex]:not([tabindex="-1"])'
    );
    if (focusable.length === 0) {
      return;
    }
    var first = focusable[0];
    var last = focusable[focusable.length - 1];
    if (keyEvent.shiftKey && document.activeElement === first) {
      keyEvent.preventDefault();
      last.focus();
    } else if (!keyEvent.shiftKey && document.activeElement === last) {
      keyEvent.preventDefault();
      first.focus();
    }
  }

  function openDetail(kind, id) {
    if (kind === "ur") {
      openUserRequestDetail(id);
    } else {
      openRequestDetail(id);
    }
  }

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
    }
  }

  function applyLens() {
    var columns = document.getElementById("kanban-columns");
    var lensHost = document.getElementById("user-request-lens");
    var byUserRequest = viewState.lens === "user-request";

    columns.hidden = byUserRequest;
    lensHost.hidden = !byUserRequest;

    if (byUserRequest && !renderedOnce.userRequestLens) {
      renderUserRequestLens();
      renderedOnce.userRequestLens = true;
    }
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
  overlay.addEventListener("click", closeDrawer);

  // ---- boot ---------------------------------------------------------------

  wireControls();
  renderWarningsBanner();
  renderColumns();
  applyView();
})();
