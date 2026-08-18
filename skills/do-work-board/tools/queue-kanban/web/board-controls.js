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
      durations: document.getElementById("view-durations"),
      timeline: document.getElementById("view-timeline"),
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
    document.getElementById("filter-clear").hidden = !hasActiveVisibleFilters();

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

  // Event delegation: any element carrying data-detail-kind opens the drawer.
  document.addEventListener("click", function (clickEvent) {
    var trigger = clickEvent.target.closest("[data-detail-kind]");
    if (trigger && !trigger.disabled) {
      clickEvent.preventDefault();
      openDetail(trigger.getAttribute("data-detail-kind"), trigger.getAttribute("data-detail-id"));
    }
  });

  document.getElementById("detail-close").addEventListener("click", closeDrawer);
