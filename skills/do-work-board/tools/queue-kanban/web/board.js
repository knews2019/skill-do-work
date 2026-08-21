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
    view: "board", // "board" | "calendar" | "durations" | "timeline" | "testing"
    lens: "flat", // "flat" | "user-request"
    windowHours: 24
  };

  // Shared filters — applied to whichever view is active. userRequestActivity
  // only affects the by-UR lens; userRequestHasOpenOrRecentWork is the canonical
  // statement of what "active" means there, so don't restate the rule here.
  // doneWindow only affects the testing view (its select is hidden elsewhere):
  // "" | hours-as-string ("24", "168", "720") | "old" (older than 30 days).
  var filterState = {
    searchText: "",
    domain: "",
    status: "",
    doneWindow: "",
    userRequestActivity: "active" // "active" | "all"
  };

  var renderedOnce = {
    userRequestLens: false,
    calendar: false,
    durations: false,
    timeline: false,
    testing: false
  };

/* INLINE_BOARD_FRAGMENTS */
  // ---- boot ---------------------------------------------------------------

  var generatedHeaderNode = document.getElementById("board-generated");
  if (generatedHeaderNode) {
    var generatedRelativeNode = makeRelativeTimeNode(boardData.generatedAt);
    if (generatedRelativeNode) {
      generatedHeaderNode.appendChild(document.createTextNode(" "));
      generatedHeaderNode.appendChild(generatedRelativeNode);
    }
  }
  setInterval(refreshRelativeTimeNodes, 1000);

  wireControls();
  wireTestingControls();
  populateFilterSelects();
  populateTestingProfileSelect();
  renderWarningsBanner();
  renderAnomaliesStrip();
  renderVerifyFindingsStrip();
  renderNotesStrip();
  renderColumns();
  applyView();
})();
