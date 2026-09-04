  // ---- activity -----------------------------------------------------------
  // "What changed on the queue in the last N hours, and why." One row per
  // lifecycle STAMP, not per REQ: a REQ captured, claimed, dispatched, merged
  // and completed inside the window carries five rows, interleaved with every
  // other REQ by time, newest first. Status does not filter it — a held,
  // blocked, claimed, completed, cancelled or failed REQ inside the window all
  // belong here, which is the whole point: every other time surface on this
  // board is status-shaped, so three REQs claimed, built, merged and held in
  // one afternoon showed up on none of them and the question needed `git log`.
  //
  // The rows arrive already ordered and already phrased. Go decided which
  // stamps exist and what each records (activity.go); this file windows them
  // and draws them, and must not re-derive either — a second definition of
  // "what this stamp means" is the thing the payload shape exists to prevent.

  function activityRowsWithin(windowHours) {
    // Anchored to the wall clock at render time rather than to generatedAt, for
    // the same reason recentlyDoneIds is (board-cards.js): a tab left open past
    // the snapshot would otherwise keep counting the window back from
    // page-generation, so "last 24 hours" slowly stops meaning that.
    var cutoffMs = Date.now() - windowHours * 3600 * 1000;
    return (boardData.activity || []).filter(function (row) {
      var stampMs = Date.parse(row.stampAt);
      return !isNaN(stampMs) && stampMs > cutoffMs;
    });
  }

  function applyActivityWindowSelection(selectedWindowHours) {
    viewState.activityWindowHours = selectedWindowHours || 24;
    setActiveButton("#activity-window-group", "data-activity-window", String(viewState.activityWindowHours));
    renderActivity();
  }

  function activityWindowPhrase(windowHours) {
    if (windowHours < 24) {
      return "last " + windowHours + " hours";
    }
    if (windowHours === 24) {
      return "last 24 hours";
    }
    return "last " + Math.round(windowHours / 24) + " days";
  }

  function renderActivity() {
    var summaryNode = document.getElementById("activity-summary");
    var tableBody = document.getElementById("activity-table-body");
    var emptyNode = document.getElementById("activity-empty");
    if (!summaryNode || !tableBody || !emptyNode) {
      return;
    }

    var windowHours = viewState.activityWindowHours || 24;
    // The shared filter chips apply here, as they do on the Timeline: "what
    // changed lately in this domain" is a straightforward question to ask of a
    // queue. Three numbers, because they answer different questions — how many
    // transitions the window holds, how many distinct REQs those belong to
    // (one busy REQ can own six of them), and, when the chips hide something,
    // how many transitions there were before filtering.
    var windowRows = activityRowsWithin(windowHours);
    var rows = windowRows.filter(function (row) {
      return requestMatchesFilters(row.id);
    });

    var requestIdsSeen = {};
    var distinctRequestCount = 0;
    rows.forEach(function (row) {
      if (!requestIdsSeen[row.id]) {
        requestIdsSeen[row.id] = true;
        distinctRequestCount += 1;
      }
    });

    var summaryText = rows.length + (rows.length === 1 ? " transition" : " transitions")
      + " across " + distinctRequestCount + (distinctRequestCount === 1 ? " REQ" : " REQs")
      + " in the " + activityWindowPhrase(windowHours);
    if (rows.length !== windowRows.length) {
      // Transitions before filtering, matching the leading count rather than
      // introducing a fourth unit no reader asked for.
      summaryText += " (" + windowRows.length + " before filters)";
    }
    summaryNode.textContent = summaryText;

    tableBody.textContent = "";
    rows.forEach(function (row) {
      var request = requestsById[row.id] || {};
      var tableRow = document.createElement("tr");
      // Deliberately NOT unique: several rows of one REQ carry the same id, so
      // this attribute selects a REQ's whole set of transitions rather than a
      // single row. A selector written as "the node for REQ-570" is wrong here.
      tableRow.setAttribute("data-activity-request", row.id);
      [
        { text: row.id, columnHeaderId: "activity-table-column-req", rowHeader: true },
        { text: request.title || "", columnHeaderId: "activity-table-column-title" },
        { text: request.status || "", columnHeaderId: "activity-table-column-status" },
        { text: row.transition || "", columnHeaderId: "activity-table-column-transition" },
        { instant: row.stampAt, columnHeaderId: "activity-table-column-when" },
        { text: row.stampField || "", columnHeaderId: "activity-table-column-stamp" }
      ].forEach(function (cellDefinition) {
        var cell = document.createElement(cellDefinition.rowHeader ? "th" : "td");
        if (cellDefinition.rowHeader) {
          cell.scope = "row";
        }
        cell.setAttribute("headers", cellDefinition.columnHeaderId);
        if (cellDefinition.instant) {
          var instantNode = makeInstantWithRelativeNode(cellDefinition.instant);
          if (instantNode) {
            cell.appendChild(instantNode);
          } else {
            cell.textContent = cellDefinition.instant;
          }
        } else {
          cell.textContent = cellDefinition.text;
        }
        tableRow.appendChild(cell);
      });
      tableBody.appendChild(tableRow);
    });

    // An empty surface must say which of the two empties it is: nothing moved,
    // or the filters hid what did. Rendering the same blank table for both is
    // how "the queue is idle" and "you have a domain chip on" become the same
    // picture, which is the misreading this whole view was raised to fix.
    if (rows.length > 0) {
      emptyNode.hidden = true;
      return;
    }
    emptyNode.hidden = false;
    emptyNode.textContent = windowRows.length === 0
      ? "No lifecycle transition falls inside the " + activityWindowPhrase(windowHours) + "."
      : windowRows.length + (windowRows.length === 1 ? " transition" : " transitions")
        + " happened in this window, but the active filters hide "
        + (windowRows.length === 1 ? "it." : "all of them.");
  }
