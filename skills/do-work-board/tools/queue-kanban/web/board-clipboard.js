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

  // Two Copy shapes, one per path, and the difference is deliberate:
  //
  //   primary  — board-markdown.js is available, so the payload is the ticket FILE
  //              (frontmatter fence + body) and is copied VERBATIM. id and title
  //              ride in the fence, so no heading is synthesized and no H1 is
  //              de-duplicated: verbatim has to mean verbatim or the paste stops
  //              round-tripping back into a valid file.
  //   fallback — board-markdown.js is stale or missing, so all that exists is the
  //              drawer's rendered text. That has no frontmatter, so it gets the
  //              identifying heading below; pasted elsewhere it would otherwise be
  //              anonymous prose. A fence is NEVER fabricated here from display
  //              state — a reassembled fence that looks verbatim is worse than a
  //              heading, because it pastes as a file whose values were guessed.
  //
  // copyHeadingForDetail and copyTextWithHeading therefore belong to the fallback
  // path only.
  function copyHeadingForDetail(detailKind, detailId) {
    var record = detailKind === "ur" ? userRequestsById[detailId] : requestsById[detailId];
    var recordTitle = record && record.title ? String(record.title).trim() : "";
    return recordTitle ? "# " + detailId + ": " + recordTitle : "# " + detailId;
  }

  // REQ/UR bodies conventionally open with an H1 restating the frontmatter title.
  // Where they do, the heading built above supersedes it (same words, plus the id)
  // and the duplicate line is dropped; otherwise the body is left untouched.
  function copyTextWithHeading(detailKind, detailId, bodyText) {
    var headingLine = copyHeadingForDetail(detailKind, detailId);
    var record = detailKind === "ur" ? userRequestsById[detailId] : requestsById[detailId];
    var recordTitle = record && record.title ? String(record.title) : "";
    var bodyLines = String(bodyText || "").split("\n");

    var firstContentIndex = 0;
    while (firstContentIndex < bodyLines.length && bodyLines[firstContentIndex].trim() === "") {
      firstContentIndex += 1;
    }
    var firstContentLine = firstContentIndex < bodyLines.length ? bodyLines[firstContentIndex].trim() : "";
    // Strip an H1 marker when present; the rendered-text fallback path carries the
    // same restated title with no "#", so both shapes are compared bare.
    var firstHeadingText = firstContentLine.replace(/^#\s+/, "");
    if (recordTitle && normalizeHeadingText(firstHeadingText) === normalizeHeadingText(recordTitle)) {
      bodyLines = bodyLines.slice(firstContentIndex + 1);
    } else {
      bodyLines = bodyLines.slice(firstContentIndex);
    }

    var remainingBody = bodyLines.join("\n").replace(/^\n+/, "");
    return remainingBody ? headingLine + "\n\n" + remainingBody : headingLine + "\n";
  }

  function normalizeHeadingText(headingText) {
    return String(headingText).replace(/\s+/g, " ").trim().toLowerCase();
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
  var copyFeedbackButton = null;

  function resetCopyButton(copyButton) {
    copyButton.textContent = copyButton.dataset.copyDefaultLabel || "Copy";
    copyButton.classList.remove("is-copied", "is-copy-failed");
  }

  function beginCopyFeedback(copyButton) {
    if (copyFeedbackTimer) {
      clearTimeout(copyFeedbackTimer);
      copyFeedbackTimer = null;
    }
    if (copyFeedbackButton && copyFeedbackButton !== copyButton) {
      resetCopyButton(copyFeedbackButton);
    }
    copyFeedbackButton = copyButton;
    copyButton.textContent = "Copying…";
    copyButton.classList.remove("is-copied", "is-copy-failed");
  }

  function showCopyFeedback(copyButton, labelText, stateClass) {
    beginCopyFeedback(copyButton);
    copyButton.textContent = labelText;
    copyButton.classList.add(stateClass);
    copyFeedbackTimer = setTimeout(function () {
      resetCopyButton(copyButton);
      if (copyFeedbackButton === copyButton) {
        copyFeedbackButton = null;
      }
      copyFeedbackTimer = null;
    }, 1600);
  }

  drawerCopyButton.addEventListener("click", function () {
    var requestedKind = currentDetailKind;
    var requestedId = currentDetailId;
    var renderedTextFallback = drawerBody.innerText || "";

    beginCopyFeedback(drawerCopyButton);

    loadBoardMarkdownData()
      .then(
        function (markdownData) {
          var rawMarkdown = rawMarkdownForDetail(markdownData, requestedKind, requestedId);
          // Primary path: the file's own bytes, copied untouched.
          if (rawMarkdown !== null) {
            return rawMarkdown;
          }
          return copyTextWithHeading(requestedKind, requestedId, renderedTextFallback);
        },
        function () {
          // Keep Copy useful for stale/incomplete generated bundles that lack
          // board-markdown.js. Only this path has no frontmatter text available,
          // so only this path gets the synthesized heading.
          return copyTextWithHeading(requestedKind, requestedId, renderedTextFallback);
        }
      )
      .then(writeTextToClipboard)
      .then(
        function () {
          if (!drawer.hidden && currentDetailKind === requestedKind && currentDetailId === requestedId) {
            showCopyFeedback(drawerCopyButton, "Copied ✓", "is-copied");
          }
        },
        function () {
          if (!drawer.hidden && currentDetailKind === requestedKind && currentDetailId === requestedId) {
            showCopyFeedback(drawerCopyButton, "Copy failed", "is-copy-failed");
          }
        }
      );
  });

  function visibleRequestIdsForColumn(copyButton) {
    var columnNode = copyButton.closest(".kanban-column");
    if (!columnNode) {
      return [];
    }
    return Array.prototype.map.call(
      columnNode.querySelectorAll('.req-card[data-detail-kind="req"][data-detail-id]'),
      function (requestCard) {
        return requestCard.dataset.detailId;
      }
    );
  }

  function rawMarkdownForRequests(markdownData, requestIds) {
    return requestIds
      .map(function (requestId) {
        var rawMarkdown = rawMarkdownForDetail(markdownData, "req", requestId);
        if (rawMarkdown === null) {
          throw new Error("raw Markdown unavailable for " + requestId);
        }
        return rawMarkdown;
      })
      // Cat semantics: preserve each file's exact bytes and invent no separator.
      .join("");
  }

  document.querySelectorAll("[data-copy-column]").forEach(function (copyButton) {
    copyButton.addEventListener("click", function () {
      var requestIds = visibleRequestIdsForColumn(copyButton);
      if (requestIds.length === 0) {
        return; // Empty columns are disabled; retain a defensive no-op for scripted clicks.
      }

      beginCopyFeedback(copyButton);
      loadBoardMarkdownData()
        .then(function (markdownData) {
          return rawMarkdownForRequests(markdownData, requestIds);
        })
        .then(writeTextToClipboard)
        .then(
          function () {
            showCopyFeedback(copyButton, "Copied ✓", "is-copied");
          },
          function () {
            showCopyFeedback(copyButton, "Copy failed", "is-copy-failed");
          }
        );
    });
  });
