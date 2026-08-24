  // ---- linkification (ticket ids, urls, file paths) -------------------------
  // Every mention that can navigate somewhere is a real, obviously-styled link:
  // REQ/UR ids open the matching drawer, URLs open in a new tab, and repo file
  // paths that the Go build verified to exist open through the live server's
  // read-only /file endpoint (serve mode only — a static snapshot has no server
  // to answer them). Paths the Go build verified to NOT exist are flagged as
  // missing in both modes.

  var liveFileApiAvailable = Boolean(boardData.liveFileApi);

  // Build-time existence verdict for every file path mentioned in a body
  // (collectRepoFileMentions in Go stats each one against the repo root):
  // true → the file exists (link it), false → it does not (flag it as
  // missing), absent → the Go scanner never saw it, leave it plain.
  var repoFileMentionExists = boardData.repoFileMentions || {};

  // A body mention like "REQ-031" may name a compound card id ("UR-002-REQ-031").
  // Index each card by its REQ segment so the short form still resolves; an
  // ambiguous segment (two cards sharing it) stays unlinked rather than guessing.
  var requestIdByReqSegment = {};
  Object.keys(requestsById).forEach(function (fullRequestId) {
    var segmentMatch = /REQ-\d+[a-z]?/i.exec(fullRequestId);
    if (!segmentMatch || segmentMatch[0] === fullRequestId) {
      return;
    }
    var segmentKey = segmentMatch[0].toUpperCase();
    requestIdByReqSegment[segmentKey] = Object.prototype.hasOwnProperty.call(requestIdByReqSegment, segmentKey)
      ? null // ambiguous — never guess
      : fullRequestId;
  });

  function resolveTicketMention(mentionText) {
    if (Object.prototype.hasOwnProperty.call(requestsById, mentionText)) {
      return { kind: "req", id: mentionText };
    }
    if (Object.prototype.hasOwnProperty.call(userRequestsById, mentionText)) {
      return { kind: "ur", id: mentionText };
    }
    var segmentTargetId = requestIdByReqSegment[mentionText.toUpperCase()];
    if (segmentTargetId) {
      return { kind: "req", id: segmentTargetId };
    }
    return null;
  }

  // In-board link to a REQ/UR drawer. The document-level [data-detail-kind]
  // delegation handles the click (and prevents the href="#" navigation).
  function makeTicketLink(detailKind, detailId, linkText) {
    var ticketLink = createElement("a", "ticket-link", linkText || detailId);
    ticketLink.href = "#";
    ticketLink.dataset.detailKind = detailKind;
    ticketLink.dataset.detailId = detailId;
    return ticketLink;
  }

  function makeExternalUrlLink(urlText) {
    var urlLink = createElement("a", "external-url-link", urlText);
    urlLink.href = urlText;
    urlLink.target = "_blank";
    urlLink.rel = "noopener";
    return urlLink;
  }

  function makeRepoFileLink(repoRelativePath) {
    var fileLink = createElement("a", "repo-file-link", repoRelativePath);
    fileLink.href = "file?path=" + encodeURIComponent(repoRelativePath);
    fileLink.target = "_blank";
    fileLink.rel = "noopener";
    fileLink.title = "Open " + repoRelativePath + " (read-only)";
    return fileLink;
  }

  // Body mention scanner, alternation order = priority: (1) absolute http(s)
  // URL, (2) repo-relative file path whose final segment carries a dot
  // extension (so "and/or" or a bare directory never match), (3) REQ/UR ticket
  // id — compound form first so "UR-002-REQ-031" is one mention, not two.
  // The file-path alternative MUST stay in lock-step with
  // repoFileMentionPattern in filementions.go — the Go scanner decides which
  // paths exist, and a drift silently downgrades mentions to plain text.
  var bodyMentionPattern = new RegExp(
    "(https?://[^\\s<>\"')\\]]+)" +
      "|((?:[A-Za-z0-9_@-]+(?:\\.[A-Za-z0-9_-]+)*/)+[A-Za-z0-9_@-][A-Za-z0-9_@.-]*\\.[A-Za-z][A-Za-z0-9]{0,7})" +
      "|(\\b(?:UR-\\d+-REQ-\\d+[a-z]?|REQ-\\d+[a-z]?|UR-\\d+)\\b)",
    "g"
  );

  // Turns one text run into a fragment with its linkable mentions wrapped.
  // Returns null when nothing linked, so callers can leave the DOM untouched.
  // File paths are only trusted inside code spans — REQ bodies conventionally
  // backtick real paths, and prose fractions like "TLS1.2/1.3" would otherwise
  // produce dead links.
  function buildLinkifiedFragment(sourceText, insideCodeSpan) {
    var fragment = document.createDocumentFragment();
    var linkedAnything = false;
    var cursorIndex = 0;
    var matchResult;
    bodyMentionPattern.lastIndex = 0;
    while ((matchResult = bodyMentionPattern.exec(sourceText)) !== null) {
      var mentionText = matchResult[0];
      var linkNode = null;
      if (matchResult[1]) {
        // Trailing sentence punctuation belongs to the prose, not the URL.
        var trimmedUrl = mentionText.replace(/[.,;:!?]+$/, "");
        mentionText = trimmedUrl;
        linkNode = makeExternalUrlLink(trimmedUrl);
      } else if (matchResult[2]) {
        if (insideCodeSpan) {
          if (repoFileMentionExists[mentionText] === false) {
            // The Go side checked: this path does not exist in the repo.
            linkNode = createElement("span", "repo-file-missing", mentionText);
            linkNode.title = "Not found in this repository";
          } else if (repoFileMentionExists[mentionText] === true && liveFileApiAvailable) {
            linkNode = makeRepoFileLink(mentionText);
          }
          // Exists but static snapshot (no server to open it), or unknown to
          // the Go scanner: stays plain text.
        }
      } else if (matchResult[3]) {
        var ticketTarget = resolveTicketMention(mentionText);
        if (ticketTarget) {
          linkNode = makeTicketLink(ticketTarget.kind, ticketTarget.id, mentionText);
        }
      }
      if (!linkNode) {
        // Skipped mention: resume just past its start so anything nested in it
        // (a REQ id inside a skipped file path) can still match.
        bodyMentionPattern.lastIndex = matchResult.index + 1;
        continue;
      }
      if (matchResult.index > cursorIndex) {
        fragment.appendChild(document.createTextNode(sourceText.slice(cursorIndex, matchResult.index)));
      }
      fragment.appendChild(linkNode);
      cursorIndex = matchResult.index + mentionText.length;
      bodyMentionPattern.lastIndex = cursorIndex;
      linkedAnything = true;
    }
    if (!linkedAnything) {
      return null;
    }
    if (cursorIndex < sourceText.length) {
      fragment.appendChild(document.createTextNode(sourceText.slice(cursorIndex)));
    }
    return fragment;
  }

  // Post-processes a drawer body after its rendered-Markdown innerHTML lands:
  // retargets the renderer's own autolinks to a new tab (a body link must not
  // navigate the board away), then wraps every linkable mention in text nodes.
  function linkifyDetailBody(bodyRootElement, recordTitle) {
    var firstBodyElement = bodyRootElement.firstElementChild;
    if (
      recordTitle &&
      firstBodyElement &&
      firstBodyElement.tagName === "H1" &&
      normalizeHeadingText(firstBodyElement.textContent) === normalizeHeadingText(recordTitle)
    ) {
      firstBodyElement.remove();
    }
    bodyRootElement.querySelectorAll("a[href]").forEach(function (anchorElement) {
      if (/^https?:/i.test(anchorElement.getAttribute("href"))) {
        anchorElement.target = "_blank";
        anchorElement.rel = "noopener";
      }
    });
    var textWalker = document.createTreeWalker(bodyRootElement, NodeFilter.SHOW_TEXT, null);
    var textNodes = [];
    while (textWalker.nextNode()) {
      textNodes.push(textWalker.currentNode);
    }
    textNodes.forEach(function (textNode) {
      var parentElement = textNode.parentElement;
      if (!parentElement || parentElement.closest("a")) {
        return;
      }
      var replacementFragment = buildLinkifiedFragment(textNode.nodeValue, Boolean(parentElement.closest("code")));
      if (replacementFragment) {
        textNode.parentNode.replaceChild(replacementFragment, textNode);
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
  var drawerCopyAllButton = document.getElementById("detail-copy-all");
  var lastFocusedElement = null;

  // The initial board payload deliberately omits raw Markdown. Remember only
  // the open record's identity; the Copy buttons load board-markdown.js on first
  // use, then look up the exact sources by kind + id.
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

  function schemaFieldDetailValue(originalValue, normalizedValue, isUnrecognized) {
    var displayedValue = originalValue || normalizedValue || "—";
    if (!isUnrecognized) {
      return displayedValue;
    }
    var invalidValue = createElement("span", "detail-status-invalid");
    invalidValue.appendChild(document.createTextNode(displayedValue));
    invalidValue.appendChild(createElement("span", "status-invalid-flag", "invalid"));
    return invalidValue;
  }

  // Each dependency listed with the status that decides whether it is met, so
  // "why is this still waiting?" is answerable without opening the upstream REQ.
  function makeDependencyDetailList(request) {
    var unmetDependencyIds = request.unmetDependencies || [];
    var list = createElement("div", "detail-dep-list");
    request.dependsOn.forEach(function (dependencyId) {
      var isUnmet = unmetDependencyIds.indexOf(dependencyId) !== -1;
      var row = createElement("span", isUnmet ? "detail-dep is-unmet" : "detail-dep is-met");
      var dependencyIdNode = requestsById[dependencyId]
        ? makeTicketLink("req", dependencyId)
        : createElement("span", null, dependencyId);
      dependencyIdNode.classList.add("detail-dep-id");
      row.appendChild(dependencyIdNode);
      row.appendChild(createElement("span", "detail-dep-status", describeRequestStatus(dependencyId)));
      list.appendChild(row);
    });
    return list;
  }

  // Comma-separated REQ id list for a meta row, each known id a drawer link.
  // Ids not on the board (free-text blocked_by entries, archived REQs) stay text.
  function makeTicketLinkList(ticketIds) {
    var listContainer = createElement("span");
    ticketIds.forEach(function (ticketId, ticketIndex) {
      if (ticketIndex > 0) {
        listContainer.appendChild(document.createTextNode(", "));
      }
      listContainer.appendChild(
        requestsById[ticketId] ? makeTicketLink("req", ticketId) : document.createTextNode(ticketId)
      );
    });
    return listContainer;
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
    if (request.error) {
      appendMetaRow("Error", request.error);
    }
    if (request.originalErrorType || request.errorType) {
      appendMetaRow(
        "Error type",
        schemaFieldDetailValue(request.originalErrorType, request.errorType, request.errorTypeUnrecognized)
      );
    }
    if (request.domain || request.originalDomain) {
      appendMetaRow(
        "Domain",
        schemaFieldDetailValue(request.originalDomain, request.domain, request.domainUnrecognized)
      );
    }
    if (request.userRequestId) {
      appendMetaRow("User request", makeTicketLink("ur", request.userRequestId));
    }
    if (request.dependsOn && request.dependsOn.length > 0) {
      appendMetaRow("Depends on", makeDependencyDetailList(request));
    }
    if (request.blockedBy && request.blockedBy.length > 0) {
      appendMetaRow("Blocked by", makeTicketLinkList(request.blockedBy));
      if (request.blockedAt) {
        // While the hold is live the row carries the ticking stopwatch; on any
        // other status (stale leftover field) it degrades to the plain instant.
        var blockedRowValue = isBlockedHoldStatus(request.status)
          ? makeInstantWithStopwatchNode(request.blockedAt)
          : null;
        appendMetaRow(
          "Blocked since",
          blockedRowValue || makeInstantWithRelativeNode(request.blockedAt) || request.blockedAt
        );
      }
      if (request.blockedCheck) {
        appendMetaRow("Blocked check", request.blockedCheck);
      }
    }
    var unblockedRequestIds = activeDependentIds(request);
    if (unblockedRequestIds.length > 0) {
      appendMetaRow("Unblocks", makeTicketLinkList(unblockedRequestIds));
    }
    if (request.writeSet && request.writeSet.length > 0) {
      appendMetaRow("Write set", request.writeSet.join(", "));
    }
    if (request.assignedTo) {
      appendMetaRow("Assigned to", request.assignedTo);
    }
    // The card badge names the contending REQs; the drawer makes them clickable
    // so "what else writes these files?" is one hop away.
    var overlappingRequestIds = request.writeSetOverlaps || [];
    if (overlappingRequestIds.length > 0) {
      appendMetaRow("Overlapping write sets", makeTicketLinkList(overlappingRequestIds));
    }
    if (request.route || request.originalRoute) {
      appendMetaRow(
        "Route",
        schemaFieldDetailValue(request.originalRoute, request.route, request.routeUnrecognized)
      );
    }
    if (request.impact || request.originalImpact) {
      appendMetaRow(
        "Impact",
        schemaFieldDetailValue(request.originalImpact, request.impact, request.impactUnrecognized)
      );
    }
    if (request.effortEstimate || request.originalEffortEstimate) {
      appendMetaRow(
        "Effort estimate",
        schemaFieldDetailValue(
          request.originalEffortEstimate,
          request.effortEstimate,
          request.effortEstimateUnrecognized
        )
      );
    }
    if (request.sweep) {
      var detailSweepOpen = request.sweepInstancesOpen || 0;
      var detailSweepDone = request.sweepInstancesDone || 0;
      appendMetaRow(
        "Sweep",
        detailSweepOpen + " open, " + detailSweepDone + " done of " + (detailSweepOpen + detailSweepDone) + " instances"
      );
    }
    if (request.createdAt) {
      appendMetaRow("Created", makeInstantWithRelativeNode(request.createdAt) || request.createdAt);
    }
    if (request.claimedAt) {
      // While the claim is live the row carries the ticking stopwatch; on any
      // other status (stale leftover field) it degrades to the plain instant.
      var claimedRowValue = request.status === "claimed" ? makeInstantWithStopwatchNode(request.claimedAt) : null;
      appendMetaRow(
        "Claimed",
        claimedRowValue || makeInstantWithRelativeNode(request.claimedAt) || request.claimedAt
      );
    }
    if (request.completionTime) {
      var completedRowValue = createElement("span");
      completedRowValue.appendChild(
        makeInstantWithRelativeNode(request.completionTime) ||
          document.createTextNode(formatShortInstant(request.completionTime))
      );
      completedRowValue.appendChild(document.createTextNode(" (" + request.completionTimeSource + ")"));
      appendMetaRow("Completed", completedRowValue);
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
        appendMetaRow("Testing updated", makeInstantWithRelativeNode(request.testingUpdatedAt) || request.testingUpdatedAt);
      }
      if (request.testingFeedback) {
        appendMetaRow("Testing feedback", request.testingFeedback);
      }
    }
    appendMetaRow("Tree", request.treeSection || "—");

    drawerBody.innerHTML = request.bodyHtml || "<p>(empty body)</p>";
    linkifyDetailBody(drawerBody, request.title);
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
      appendMetaRow("REQ ids", makeTicketLinkList(requestIds));
    }
    appendMetaRow("input.md", userRequest.inputFilePresent ? "present" : "synthesized from REQ pointers");

    drawerBody.innerHTML = userRequest.bodyHtml || "<p>(no input.md body)</p>";
    linkifyDetailBody(drawerBody, userRequest.title);
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
    // on the clipboard — reset both controls on every open. Copy all belongs to
    // UR details only; plain Copy remains the one-record action for both kinds.
    drawerCopyButton.textContent = "Copy";
    drawerCopyButton.classList.remove("is-copied", "is-copy-failed");
    drawerCopyAllButton.textContent = "Copy all";
    drawerCopyAllButton.classList.remove("is-copied", "is-copy-failed");
    drawerCopyAllButton.hidden = currentDetailKind !== "ur";
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
