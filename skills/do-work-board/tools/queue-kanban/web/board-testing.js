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
  // In-progress feedback text, mirrored on every keystroke. The testing view is
  // rebuilt from scratch on each render (filter changes, failed posts), which
  // repaints the textarea — without this stash a re-render would silently
  // discard everything the tester typed.
  var feedbackDraftText = "";
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
    // Prefer the digits of the REQ segment: a compound id like UR-002-REQ-031
    // must sort by 31, not by the first digit run (2).
    var reqSegmentMatch = /REQ-0*(\d+)/i.exec(requestId || "");
    if (reqSegmentMatch) {
      return parseInt(reqSegmentMatch[1], 10);
    }
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
    // Wall clock, not generatedAtMs — see recentlyDoneIds for why.
    var nowMs = Date.now();
    if (filterState.doneWindow === "old") {
      return recencyMs !== 0 && recencyMs <= nowMs - thirtyDaysMs;
    }
    var windowHours = parseInt(filterState.doneWindow, 10);
    return recencyMs > nowMs - windowHours * 3600 * 1000;
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
      container.appendChild(createElement("p", "column-empty", columnEmptyText(hasActiveVisibleFilters())));
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
      var testingUpdatedChip = createElement("span", "testing-meta-chip");
      testingUpdatedChip.appendChild(
        makeInstantWithRelativeNode(request.testingUpdatedAt) ||
          document.createTextNode(request.testingUpdatedAt)
      );
      meta.appendChild(testingUpdatedChip);
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

    function addActionButton(labelText, onActivate, extraClassName, allowWithoutProfile) {
      var actionButton = createElement(
        "button",
        "control-button testing-action" + (extraClassName ? " " + extraClassName : ""),
        labelText
      );
      actionButton.type = "button";
      // Clear is deliberately exempt from the profile gate: the server accepts
      // profile-free clears (a clear only removes fields), and requiring a
      // tester identity to delete a stale record just blocks the cleanup.
      if (!selectedTesterProfile && !allowWithoutProfile) {
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
        feedbackDraftText = "";
        renderTestingView();
      });
      addActionButton("Clear", function () {
        postTestingStatus(requestId, "clear");
      }, "testing-action-clear", true);
    } else if (bucketKey === "testingReturned") {
      addActionButton("Restart testing", function () {
        postTestingStatus(requestId, "in-testing");
      });
      addActionButton("Clear", function () {
        postTestingStatus(requestId, "clear");
      }, "testing-action-clear", true);
    } else if (bucketKey === "testingTested") {
      addActionButton("Re-test", function () {
        postTestingStatus(requestId, "in-testing");
      });
      addActionButton("Clear", function () {
        postTestingStatus(requestId, "clear");
      }, "testing-action-clear", true);
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
    feedbackInput.value = feedbackDraftText;
    feedbackInput.addEventListener("input", function () {
      feedbackDraftText = feedbackInput.value;
    });
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
      // The form stays open (and the draft stays stashed) until the server
      // confirms — postTestingStatus closes it on success. Closing it here
      // would throw away the typed feedback whenever the POST fails.
      feedbackDraftText = feedbackInput.value;
      postTestingStatus(requestId, "returned", feedbackText);
    });
    var cancelButton = createElement("button", "control-button testing-action", "Cancel");
    cancelButton.type = "button";
    cancelButton.addEventListener("click", function () {
      feedbackFormRequestId = null;
      feedbackDraftText = "";
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

  function applyConfirmedTestingTransition(requestId, testingState, feedbackText, payload) {
    var request = requestsById[requestId];
    if (request) {
      request.testingStatus = payload.testingStatus || "";
      request.testedBy = payload.testedBy || "";
      request.testingUpdatedAt = payload.testingUpdatedAt || "";
      request.testingFeedback = testingState === "returned" ? feedbackText || "" : "";
      request.testingStatusUnrecognized = false;
      request.originalTestingStatus = payload.testingStatus || "";
    }
    if (feedbackFormRequestId === requestId) {
      feedbackFormRequestId = null;
      feedbackDraftText = "";
    }
    renderTestingView();
    renderColumns(); // the main board's testing badge tracks the same record
    renderedOnce.userRequestLens = false;
    if (viewState.view === "board" && viewState.lens === "user-request") {
      renderUserRequestLens();
      renderedOnce.userRequestLens = true;
    }
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
        applyConfirmedTestingTransition(requestId, testingState, feedbackText, payload);
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
      addToggleButton.title = "Adding testers needs the live board (do-work-board board)";
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
