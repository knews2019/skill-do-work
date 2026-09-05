  // ---- view / lens / window switching ------------------------------------

  // The Lens group's third choice, URs only, is the by-UR lens with its REQ
  // cards folded away — not a lens of its own. viewState.lens therefore stays
  // "user-request" for both, so every other test of whether the UR lens is on
  // screen (the shared filters, a testing transition, the recently-done chip)
  // keeps deciding correctly without knowing this fold exists.
  var userRequestCardsFolded = false;

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
      durations: document.getElementById("view-durations"),
      timeline: document.getElementById("view-timeline"),
      activity: document.getElementById("view-activity"),
      testing: document.getElementById("view-testing")
    };
    Object.keys(viewPanels).forEach(function (viewName) {
      var isActiveView = viewState.view === viewName;
      viewPanels[viewName].classList.toggle("is-active", isActiveView);
      viewPanels[viewName].hidden = !isActiveView;
    });

    // The board is the Timeline's scroll surface since REQ-587, and .board-main's
    // scroll position is shared by every view with nothing resetting it between
    // them. That never mattered while the Timeline's chart was capped at 58vh —
    // the page was short, so arriving from a Kanban board scrolled 800px down
    // clamped back to near the top. Now the chart is the tall thing on the page,
    // nothing clamps, and the same arrival would drop the reader 800px into the
    // middle of it with correct rows and no reason on screen for where they are.
    // Resetting is what keeps arrival looking the way it does today. It runs
    // before the first render below, so the first anchor read sees 0.
    if (viewState.view === "timeline") {
      document.getElementById("board-main").scrollTop = 0;
    }

    // The grouping lens and the recently-done window only shape the board view;
    // hide their controls elsewhere so the topbar never advertises dead knobs.
    // Shared filters do not shape Durations, and the date window is only the
    // testing view's knob, so hide those controls where they do nothing too.
    var hideSharedFilters = viewState.view === "durations";
    document.getElementById("filter-search").hidden = hideSharedFilters;
    document.getElementById("filter-domain").hidden = hideSharedFilters;
    document.getElementById("filter-status").hidden = hideSharedFilters;
    document.getElementById("lens-group").hidden = viewState.view !== "board";
    document.getElementById("recent-window-group").hidden = viewState.view !== "board";
    document.getElementById("durations-colour-group").hidden = viewState.view !== "durations";
    document.getElementById("durations-window-group").hidden = viewState.view !== "durations";
    document.getElementById("filter-done-window").hidden = viewState.view !== "testing";
    document.getElementById("filter-clear").hidden = !hasActiveVisibleFilters();

    // The Verify Findings strip lives outside the view panels, so it survives
    // every view switch — except the Activity view, which is a transitions log
    // the strip only pushes down (REQ-578). Whether it has anything to say is
    // read back off what the findings renderer drew, never decided again here:
    // a second copy of that rule is how the two would drift apart.
    var findingsStrip = document.getElementById("board-findings");
    var findingsStripHasContent =
      document.getElementById("board-findings-cards").children.length > 0 ||
      document.getElementById("board-findings-skipped-list").children.length > 0;
    findingsStrip.hidden = viewState.view === "activity" || !findingsStripHasContent;

    if (viewState.view === "calendar" && !renderedOnce.calendar) {
      renderCalendar();
      renderedOnce.calendar = true;
    }
    if (viewState.view === "durations" && !renderedOnce.durations) {
      renderDurationsView();
      renderedOnce.durations = true;
    }
    if (viewState.view === "timeline" && !renderedOnce.timeline) {
      renderTimelineView();
      renderedOnce.timeline = true;
    }
    if (viewState.view === "activity" && !renderedOnce.activity) {
      renderActivity();
      renderedOnce.activity = true;
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

  // Two of the three lens buttons select the same lens, so the active one is
  // keyed on the lens plus the fold rather than on data-lens-target alone.
  // A button without data-ur-cards means the always-open reading.
  function setActiveLensButton() {
    document.querySelectorAll("#lens-group [data-lens-target]").forEach(function (button) {
      var isActive =
        button.getAttribute("data-lens-target") === viewState.lens &&
        (button.getAttribute("data-ur-cards") === "folded") === userRequestCardsFolded;
      button.classList.toggle("is-active", isActive);
      button.setAttribute("aria-pressed", isActive ? "true" : "false");
    });
  }

  // One lens button click: the two readings render different DOM from the same
  // renderer, so the cached lens is always dropped before applyLens re-renders.
  function applyLensSelection(lensName, userRequestCardsFold) {
    viewState.lens = lensName;
    userRequestCardsFolded = userRequestCardsFold === "folded";
    setActiveLensButton();
    renderedOnce.userRequestLens = false;
    if (viewState.view === "board") {
      applyLens();
    }
  }

  function applyLens() {
    var columns = document.getElementById("kanban-columns");
    var lensHost = document.getElementById("user-request-lens");
    var byUserRequest = viewState.lens === "user-request";

    columns.hidden = byUserRequest;
    lensHost.hidden = !byUserRequest;
    lensHost.classList.toggle("is-folded", userRequestCardsFolded);
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

  // ---- top bar identity ---------------------------------------------------
  // The identity is one line — wordmark, project, clock — so the bar keeps its
  // height when the control pills beside it wrap. The visible clock is the
  // minute of the stamp; the whole stamp, with the age board.js keeps ticking
  // beside it, is clipped out of the line and shown as its tooltip.

  // Whatever #board-generated holds right now, joined for a tooltip: the
  // server's stamp text, and the relative node board.js appends after it. Read
  // as child text rather than as one string so the two parts can be separated
  // by a middle dot, and so a change to what board.js appends still reads.
  function generatedStampTooltipText(generatedNode) {
    var stampParts = [];
    Array.prototype.forEach.call(generatedNode.childNodes, function (childNode) {
      var childText = (childNode.textContent || "").trim();
      if (childText) {
        stampParts.push(childText);
      }
    });
    return stampParts.join(" · ");
  }

  function renderTopBarIdentity() {
    var identityLine = document.getElementById("board-identity");
    var generatedNode = document.getElementById("board-generated");
    var clockNode = document.getElementById("board-generated-clock");
    if (!identityLine || !generatedNode || !clockNode) {
      return;
    }
    var generatedMs = Date.parse(boardData.generatedAt);
    // Same instant and same zone the server printed into the full stamp, cut
    // to the minute: "12:17 UTC".
    clockNode.textContent = isNaN(generatedMs)
      ? ""
      : new Date(generatedMs).toISOString().slice(11, 16) + " UTC";

    identityLine.title = generatedStampTooltipText(generatedNode);
    // Rebuilt when the pointer arrives rather than on a timer: a tooltip is
    // only read while it is open, and board.js owns the page's one ticker.
    var refreshTooltip = function () {
      identityLine.title = generatedStampTooltipText(generatedNode);
    };
    identityLine.addEventListener("pointerenter", refreshTooltip);
    identityLine.addEventListener("focusin", refreshTooltip);
  }

  function wireControls() {
    renderTopBarIdentity();

    document.querySelectorAll("[data-view-target]").forEach(function (button) {
      button.addEventListener("click", function () {
        viewState.view = button.getAttribute("data-view-target");
        setActiveButton("[aria-label='Board views and lenses']", "data-view-target", viewState.view);
        applyView();
      });
    });

    document.querySelectorAll("[data-lens-target]").forEach(function (button) {
      button.addEventListener("click", function () {
        applyLensSelection(button.getAttribute("data-lens-target"), button.getAttribute("data-ur-cards"));
      });
    });

    document.querySelectorAll("[data-durations-colour]").forEach(function (button) {
      button.addEventListener("click", function () {
        var colourChannel = button.getAttribute("data-durations-colour");
        setDurationsColourChannel(colourChannel);
        setActiveButton("#durations-colour-group", "data-durations-colour", colourChannel);
        if (viewState.view === "durations") {
          renderDurationsView();
          renderedOnce.durations = true;
        }
      });
    });

    document.querySelectorAll("[data-durations-window]").forEach(function (button) {
      button.addEventListener("click", function () {
        applyDurationsWindowSelection(button.getAttribute("data-durations-window"));
      });
    });

    document.querySelectorAll("[data-activity-window]").forEach(function (button) {
      button.addEventListener("click", function () {
        applyActivityWindowSelection(parseInt(button.getAttribute("data-activity-window"), 10));
      });
    });

    document.querySelectorAll("[data-window-hours]").forEach(function (button) {
      button.addEventListener("click", function () {
        applyRecentWindowSelection(parseInt(button.getAttribute("data-window-hours"), 10));
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

  function applyDurationsWindowSelection(windowName) {
    setDurationsWindow(windowName);
    setActiveButton("#durations-window-group", "data-durations-window", windowName);
    renderedOnce.durations = false;
    if (viewState.view === "durations") {
      renderDurationsView();
      renderedOnce.durations = true;
    }
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
