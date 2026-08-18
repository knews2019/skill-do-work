  // ---- filtering ------------------------------------------------------------
  // Pure client-side: the data island already carries status, domain, UR id,
  // and titles, so every view filters the same record set with the same rules.

  function hasActiveFilters() {
    return (
      filterState.searchText !== "" ||
      filterState.domain !== "" ||
      filterState.status !== ""
    );
  }

  // doneWindow is a Testing-view filter, so it must not change empty-state
  // decisions in Columns, Calendar, or By UR after the user switches views.
  // It still counts for the visible Clear button while Testing is on screen.
  function hasActiveVisibleFilters() {
    return hasActiveFilters() || (viewState.view === "testing" && filterState.doneWindow !== "");
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

  // The by-UR lens's Active scope. "Active" is deliberately wider than "holds a
  // non-terminal REQ": on a fully-shipped queue every REQ is terminal, so the
  // narrow rule was unsatisfiable and the lens went blank while the Columns lens
  // showed those same REQs under RECENTLY DONE. A UR therefore also counts as
  // active when one of its REQs completed inside the current window.
  //
  // recentlyDoneIdSet is passed in rather than derived here: the caller builds it
  // once per render from the shared recentlyDoneIds(), which both keeps this off
  // the calendar-rescan path and guarantees the two lenses can never disagree
  // about what "recent" means.
  function userRequestHasOpenOrRecentWork(userRequest, recentlyDoneIdSet) {
    return (userRequest.requestIds || []).some(function (requestId) {
      if (recentlyDoneIdSet[requestId]) {
        return true;
      }
      var request = requestsById[requestId];
      return request && !isTerminalResolvedStatus(request.status);
    });
  }

  // Reads the selected RECENTLY DONE chip back as prose for the lens's empty and
  // hidden-count copy, so those lines track the chip instead of baking in a span.
  // Days only from 7d up, matching how the chips themselves are labelled.
  function recentWindowPhrase(windowHours) {
    if (windowHours >= 168 && windowHours % 24 === 0) {
      var windowDays = windowHours / 24;
      return "the last " + windowDays + " day" + (windowDays === 1 ? "" : "s");
    }
    return "the last " + windowHours + " hour" + (windowHours === 1 ? "" : "s");
  }

  function userRequestLensEmptyText(
    filtersActive,
    hiddenResolvedCount,
    hiddenResolvedFilterMatchCount,
    windowPhrase
  ) {
    if (hiddenResolvedFilterMatchCount > 0) {
      return (
        "No Active user requests match the current filters — switch URs to All to see " +
        hiddenResolvedFilterMatchCount +
        " resolved match" +
        (hiddenResolvedFilterMatchCount === 1 ? "." : "es.")
      );
    }
    if (filtersActive) {
      return "No user requests match the current filters.";
    }
    if (hiddenResolvedCount > 0) {
      return (
        "No user requests with open work or activity in " +
        windowPhrase +
        " — widen the RECENTLY DONE window, or switch URs to All to browse the archive."
      );
    }
    return "No user requests in this tree yet.";
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
    document.getElementById("filter-clear").hidden = !hasActiveVisibleFilters();
    // Columns have no renderedOnce guard (they render at boot), so refresh
    // them unconditionally; the lazily-rendered views refresh if visible and
    // go stale otherwise, re-rendering on their next activation.
    renderColumns();
    renderedOnce.userRequestLens = false;
    renderedOnce.calendar = false;
    // Timeline filters, unlike Durations: a Gantt narrowed to one domain is a
    // straightforward question about a queue, where a durations distribution
    // narrowed the same way is a different statistic wearing the same axes.
    renderedOnce.timeline = false;
    renderedOnce.testing = false;
    if (viewState.view === "calendar") {
      renderCalendar();
      renderedOnce.calendar = true;
    } else if (viewState.view === "timeline") {
      renderTimelineView();
      renderedOnce.timeline = true;
    } else if (viewState.view === "testing") {
      renderTestingView();
      renderedOnce.testing = true;
    } else if (viewState.lens === "user-request") {
      renderUserRequestLens();
      renderedOnce.userRequestLens = true;
    }
  }
