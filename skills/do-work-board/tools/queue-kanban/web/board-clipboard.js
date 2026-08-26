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

  // Where that document's ticket ids sit, computed by the Go walk in
  // citations.go and shipped in the same payload as the text it indexes. An
  // older generated bundle carries no index at all, so an absent entry yields
  // an empty list and the document copies exactly as it is on disk — the one
  // behaviour a missing index must never break.
  function ticketMentionsForDetail(markdownData, detailKind, detailId) {
    var mentionsById = detailKind === "ur" ? markdownData.userRequestMentions : markdownData.requestMentions;
    if (mentionsById && Object.prototype.hasOwnProperty.call(mentionsById, detailId)) {
      return mentionsById[detailId];
    }
    return [];
  }

  // Text and offsets travel together from here on. They are only meaningful as
  // a pair — an offset measured on one document spliced into another lands
  // mid-word — so no seam below takes the two separately.
  function clipboardDocumentFor(markdownData, detailKind, detailId) {
    var rawMarkdown = rawMarkdownForDetail(markdownData, detailKind, detailId);
    if (rawMarkdown === null) {
      return null;
    }
    return { text: rawMarkdown, ticketMentions: ticketMentionsForDetail(markdownData, detailKind, detailId) };
  }

  // Two Copy shapes, one per path, and the difference is deliberate:
  //
  //   primary  — board-markdown.js is available, so the payload is the ticket FILE
  //              (frontmatter fence + body). id and title ride in the fence, so
  //              no heading is synthesized and no H1 is de-duplicated.
  //   fallback — board-markdown.js is stale or missing, so all that exists is the
  //              drawer's rendered text. That has no frontmatter, so it gets the
  //              identifying heading below; pasted elsewhere it would otherwise be
  //              anonymous prose. A fence is NEVER fabricated here from display
  //              state — a reassembled fence that looks verbatim is worse than a
  //              heading, because it pastes as a file whose values were guessed.
  //
  // ONLY the raw shape passes through annotateClipboardPayload, which is where this
  // file's old "copied VERBATIM" rule now lives in a narrower form: the Go
  // payload is verbatim, the clipboard payload is those bytes plus marked
  // annotations, and the frontmatter fence is never touched. The round trip the
  // word "verbatim" was protecting — a paste that saves straight back as a valid
  // REQ or UR file — rests on the fence rule alone.
  //
  // copyHeadingForDetail and copyTextWithHeading therefore belong to the fallback
  // path only, and that path is deliberately NOT annotated. Its input is
  // drawerBody.innerText, which the drawer has already expanded — annotating it
  // again produced "REQ-1679 (Short one) Short one", duplicating every title.
  // The titles are already there, so the fallback needs nothing added.
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

  // ---- ticket annotation for the clipboard payload -------------------------
  // The Go payload stays byte-exact: buildGeneratedBoardMarkdownData ships each
  // ticket's original bytes and this file never asks it for anything else. What
  // the CLIPBOARD carries is those bytes plus two marked additions — the first
  // mention of a REQ/UR id in a BODY gains the ticket's title, and one appendix
  // at the very end names every ticket the payload referenced. The frontmatter
  // fence is never touched, so a paste still saves back as a file whose
  // depends_on, related and user_request parse. That invariant is what this
  // whole section exists to hold; everything below is in service of it.
  //
  // WHICH mentions may be annotated is a Markdown question, and this file does
  // not answer it. citations.go walks each body's goldmark AST at build time and
  // ships the positions — see its header for the division. What is left here is
  // splicing and the appendix, both of which need the board's titles rather than
  // its Markdown. The scanner this file used to carry disagreed with CommonMark
  // on six constructs; there is nothing left here that can.

  var referencedTicketsGlossaryHeading = "## Referenced requests (added by the board — not part of the file)";

  function recordReferencedTicket(annotationState, referencedTicket) {
    if (Object.prototype.hasOwnProperty.call(annotationState.referencedTicketIds, referencedTicket.id)) {
      return;
    }
    annotationState.referencedTicketIds[referencedTicket.id] = true;
    annotationState.referencedTickets.push(referencedTicket);
  }

  // One document in; the same document with its BODY annotated out, plus the
  // tickets it referenced in first-mention order. First-mention memory is per
  // document because a payload is a concatenation of whole files, and each file
  // pastes and reads on its own — which is why the Go side computes `expand`
  // per document too.
  //
  // Splicing runs BACKWARDS. Each inserted title lengthens the text, so every
  // offset after the insertion point would shift; taking them in descending
  // order means each remaining offset still measures the same character.
  function annotateTicketMentions(documentText, ticketMentions) {
    var annotationState = { referencedTicketIds: {}, referencedTickets: [] };
    var titleInsertions = [];
    (ticketMentions || []).forEach(function (ticketMention) {
      // A mention with no kind is a dead reference. Plain text has no red, so
      // the appendix is where a paste's reader learns it points at nothing.
      recordReferencedTicket(annotationState, { kind: ticketMention.kind || "", id: ticketMention.id });
      if (!ticketMention.expand) {
        return;
      }
      // describeTicketTitle never answers empty for a record that resolved — it
      // substitutes "untitled" or names the synthesized-UR case — so there is no
      // empty-title branch to take here.
      titleInsertions.push({
        offset: ticketMention.offset + ticketMention.length,
        text: " (" + shortTicketTitle(describeTicketTitle(ticketMention.kind, ticketMention.id).text) + ")"
      });
    });

    var annotatedText = documentText;
    for (var insertionIndex = titleInsertions.length - 1; insertionIndex >= 0; insertionIndex -= 1) {
      var titleInsertion = titleInsertions[insertionIndex];
      annotatedText =
        annotatedText.slice(0, titleInsertion.offset) + titleInsertion.text + annotatedText.slice(titleInsertion.offset);
    }
    return { text: annotatedText, referencedTickets: annotationState.referencedTickets };
  }

  function describeReferencedTicket(referencedTicket) {
    if (!referencedTicket.kind) {
      return "not found in this queue";
    }
    // describeRequestStatus only knows requests; a UR has no pipeline status, so
    // its line says what kind of record it is instead — the same substitution
    // the drawer glossary makes.
    var statusText = referencedTicket.kind === "ur" ? "user request" : describeRequestStatus(referencedTicket.id);
    return describeTicketTitle(referencedTicket.kind, referencedTicket.id).text + " (" + statusText + ")";
  }

  // The appendix, with full untruncated titles — the inline cut exists to keep
  // prose readable, and this is the place a reader looks up what was cut.
  // excludedIds drops tickets whose own files are already in the payload: a Copy
  // all must not gloss what it already contains, because those titles are right
  // there in their own fences.
  function buildReferencedTicketsGlossary(referencedTickets, excludedIds) {
    var excludedIdSet = {};
    (excludedIds || []).forEach(function (excludedId) {
      excludedIdSet[excludedId] = true;
    });
    var listedIds = {};
    var glossaryLines = [];
    referencedTickets.forEach(function (referencedTicket) {
      if (Object.prototype.hasOwnProperty.call(excludedIdSet, referencedTicket.id) ||
        Object.prototype.hasOwnProperty.call(listedIds, referencedTicket.id)) {
        return;
      }
      listedIds[referencedTicket.id] = true;
      glossaryLines.push("- " + referencedTicket.id + " — " + describeReferencedTicket(referencedTicket));
    });
    if (glossaryLines.length === 0) {
      return "";
    }
    return "\n---\n\n" + referencedTicketsGlossaryHeading + "\n\n" + glossaryLines.join("\n") + "\n";
  }

  // The one seam every Copy handler goes through. Each document is annotated
  // BEFORE the join, because a joined payload carries N frontmatter fences and
  // offsets measured against one document mean nothing in the concatenation —
  // splicing after the join would rewrite a later ticket's depends_on, the exact
  // corruption this feature must not cause. The join itself keeps cat semantics:
  // exact bytes, no invented separator.
  function annotateClipboardPayload(clipboardDocuments, excludedIds) {
    var referencedTickets = [];
    var annotatedDocuments = clipboardDocuments.map(function (clipboardDocument) {
      var annotatedDocument = annotateTicketMentions(clipboardDocument.text, clipboardDocument.ticketMentions);
      referencedTickets = referencedTickets.concat(annotatedDocument.referencedTickets);
      return annotatedDocument.text;
    });
    var joinedPayload = annotatedDocuments.join("");
    var glossaryText = buildReferencedTicketsGlossary(referencedTickets, excludedIds);
    if (glossaryText === "") {
      return joinedPayload;
    }
    // The blank line the glossary opens with matters: a "---" directly under a
    // text line is a setext H2, which would swallow the payload's last line.
    var missingNewline =
      joinedPayload !== "" && joinedPayload.charAt(joinedPayload.length - 1) !== "\n" ? "\n" : "";
    return joinedPayload + missingNewline + glossaryText;
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
          var clipboardDocument = clipboardDocumentFor(markdownData, requestedKind, requestedId);
          // Primary path: the file's own bytes, annotated in the body only.
          if (clipboardDocument !== null) {
            return annotateClipboardPayload([clipboardDocument], [requestedId]);
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

  // One entry per ticket, in display order, so annotateClipboardPayload can
  // annotate each file before they are concatenated. The throw fails the whole
  // operation rather than silently publishing an incomplete clipboard payload.
  function clipboardDocumentsForRequests(markdownData, requestIds) {
    return requestIds.map(function (requestId) {
      var clipboardDocument = clipboardDocumentFor(markdownData, "req", requestId);
      if (clipboardDocument === null) {
        throw new Error("raw Markdown unavailable for " + requestId);
      }
      return clipboardDocument;
    });
  }

  function clipboardDocumentsForUserRequestAndRequests(markdownData, userRequestId, requestIds) {
    var userRequestDocument = clipboardDocumentFor(markdownData, "ur", userRequestId);
    if (userRequestDocument === null) {
      throw new Error("raw Markdown unavailable for " + userRequestId);
    }
    // The grouped id list is the same all-tree set rendered in the UR drawer.
    return [userRequestDocument].concat(clipboardDocumentsForRequests(markdownData, requestIds));
  }

  drawerCopyAllButton.addEventListener("click", function () {
    if (currentDetailKind !== "ur") {
      return; // Hidden on REQ details; retain a defensive no-op for scripted clicks.
    }
    var requestedUserRequestId = currentDetailId;
    var requestedUserRequest = userRequestsById[requestedUserRequestId];
    var requestedRequestIds = requestedUserRequest && requestedUserRequest.requestIds
      ? requestedUserRequest.requestIds.slice()
      : [];

    beginCopyFeedback(drawerCopyAllButton);
    loadBoardMarkdownData()
      .then(function (markdownData) {
        // The UR and its grouped REQs are all in the payload, so none of them
        // earns an appendix line — their titles ride in their own fences.
        return annotateClipboardPayload(
          clipboardDocumentsForUserRequestAndRequests(
            markdownData, requestedUserRequestId, requestedRequestIds
          ),
          [requestedUserRequestId].concat(requestedRequestIds)
        );
      })
      .then(writeTextToClipboard)
      .then(
        function () {
          if (!drawer.hidden && currentDetailKind === "ur" && currentDetailId === requestedUserRequestId) {
            showCopyFeedback(drawerCopyAllButton, "Copied ✓", "is-copied");
          }
        },
        function () {
          if (!drawer.hidden && currentDetailKind === "ur" && currentDetailId === requestedUserRequestId) {
            showCopyFeedback(drawerCopyAllButton, "Copy failed", "is-copy-failed");
          }
        }
      );
  });

  document.querySelectorAll("[data-copy-column]").forEach(function (copyButton) {
    copyButton.addEventListener("click", function () {
      var requestIds = visibleRequestIdsForColumn(copyButton);
      if (requestIds.length === 0) {
        return; // Empty columns are disabled; retain a defensive no-op for scripted clicks.
      }

      beginCopyFeedback(copyButton);
      loadBoardMarkdownData()
        .then(function (markdownData) {
          return annotateClipboardPayload(clipboardDocumentsForRequests(markdownData, requestIds), requestIds);
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
