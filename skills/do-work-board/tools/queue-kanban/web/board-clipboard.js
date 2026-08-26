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
  //              (frontmatter fence + body). id and title ride in the fence, so
  //              no heading is synthesized and no H1 is de-duplicated.
  //   fallback — board-markdown.js is stale or missing, so all that exists is the
  //              drawer's rendered text. That has no frontmatter, so it gets the
  //              identifying heading below; pasted elsewhere it would otherwise be
  //              anonymous prose. A fence is NEVER fabricated here from display
  //              state — a reassembled fence that looks verbatim is worse than a
  //              heading, because it pastes as a file whose values were guessed.
  //
  // Both shapes then pass through annotateClipboardPayload, which is where this
  // file's old "copied VERBATIM" rule now lives in a narrower form: the Go
  // payload is verbatim, the clipboard payload is those bytes plus marked
  // annotations, and the frontmatter fence is never touched. The round trip the
  // word "verbatim" was protecting — a paste that saves straight back as a valid
  // REQ or UR file — rests on the fence rule alone.
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

  // ---- ticket annotation for the clipboard payload -------------------------
  // The Go payload stays byte-exact: buildGeneratedBoardMarkdownData ships each
  // ticket's original bytes and this file never asks it for anything else. What
  // the CLIPBOARD carries is those bytes plus two marked additions — the first
  // mention of a REQ/UR id in a BODY gains the ticket's title, and one appendix
  // at the very end names every ticket the payload referenced. The frontmatter
  // fence is never touched, so a paste still saves back as a file whose
  // depends_on, related and user_request parse. That invariant is what this
  // whole section exists to hold; everything below is in service of it.

  var referencedTicketsGlossaryHeading = "## Referenced requests (added by the board — not part of the file)";

  // Mirrors splitFrontmatter in frontmatter.go and must stay in lock-step with
  // it: a document is fenced only when it STARTS with "---" (one UTF-8 BOM may
  // precede it), and the fence closes at the first later line equal to "---"
  // after one optional carriage return. No opening fence, no closing fence, and
  // a lone "---" line all mean everything is body — reading an unterminated
  // fence as frontmatter would silently skip a whole document instead.
  function frontmatterFenceEndOffset(documentText) {
    var scanOffset = documentText.charCodeAt(0) === 0xfeff ? 1 : 0;
    var afterByteOrderMark = documentText.slice(scanOffset);
    if (afterByteOrderMark.indexOf("---\n") === 0) {
      scanOffset += "---\n".length;
    } else if (afterByteOrderMark.indexOf("---\r\n") === 0) {
      scanOffset += "---\r\n".length;
    } else {
      return 0;
    }
    while (true) {
      var newlineIndex = documentText.indexOf("\n", scanOffset);
      var currentLine = newlineIndex < 0
        ? documentText.slice(scanOffset)
        : documentText.slice(scanOffset, newlineIndex);
      if (currentLine.replace(/\r$/, "") === "---") {
        return newlineIndex < 0 ? documentText.length : newlineIndex + 1;
      }
      if (newlineIndex < 0) {
        return 0;
      }
      scanOffset = newlineIndex + 1;
    }
  }

  // A fenced block opens on a line whose leading run — after at most three
  // spaces — is three or more backticks or tildes. Returns null for any other
  // line, so a caller can use it as both the open and the close test.
  // Strips Markdown container prefixes — blockquote markers, and the list
  // markers a fence can sit under — so a fence is recognised wherever it is
  // nested. This is not an edge case: the outside-text containment contract
  // (actions/clarify.md Step 4) writes every UR's Full Verbatim Input as a
  // BLOCKQUOTED fence, so without this the preserved words of the user's own
  // request get annotated. UR-075's verbatim block alone holds 21 ticket ids.
  //
  // Only the prefix is removed; the indent rules below then apply to what is
  // left, which is what CommonMark does inside a container.
  function stripContainerPrefix(lineText) {
    var strippedText = lineText;
    var keepStripping = true;
    while (keepStripping) {
      keepStripping = false;
      var blockquoteMatch = /^ {0,3}> ?/.exec(strippedText);
      if (blockquoteMatch) {
        strippedText = strippedText.slice(blockquoteMatch[0].length);
        keepStripping = true;
      }
    }
    return strippedText;
  }

  function codeFenceRunFor(lineText) {
    var containerText = stripContainerPrefix(lineText);
    var scanOffset = 0;
    while (scanOffset < 3 && containerText.charAt(scanOffset) === " ") {
      scanOffset += 1;
    }
    var fenceCharacter = containerText.charAt(scanOffset);
    if (fenceCharacter !== "`" && fenceCharacter !== "~") {
      return null;
    }
    var runEndOffset = scanOffset;
    while (containerText.charAt(runEndOffset) === fenceCharacter) {
      runEndOffset += 1;
    }
    if (runEndOffset - scanOffset < 3) {
      return null;
    }
    var infoText = containerText.slice(runEndOffset);
    // CommonMark forbids a backtick anywhere in a BACKTICK fence's info string,
    // so ```lang`invalid is prose, not a fence. This repo's own renderer already
    // agrees — TestRenderMarkdownInvalidBacktickInfoRemainsQuestionProse pins it
    // — and treating it as a fence here opened a block that never opens in the
    // rendered body, silently leaving every reference until the next fence or
    // EOF unannotated. A tilde fence has no such rule and may carry backticks.
    if (fenceCharacter === "`" && infoText.indexOf("`") >= 0) {
      return null;
    }
    return {
      fenceCharacter: fenceCharacter,
      runLength: runEndOffset - scanOffset,
      infoText: infoText
    };
  }

  // A closing fence uses the opener's character, is at least as long, and
  // carries no info string. The length rule is what lets a ````-fenced block
  // quote a ```-fenced one, which REQ bodies do whenever they show a template.
  function codeFenceRunCloses(candidateRun, openFenceRun) {
    return (
      candidateRun !== null &&
      candidateRun.fenceCharacter === openFenceRun.fenceCharacter &&
      candidateRun.runLength >= openFenceRun.runLength &&
      candidateRun.infoText.trim() === ""
    );
  }

  // Finds the backtick run that closes an inline code span opened by a run of
  // runLength backticks. Only an EQUAL-length run closes it, which is what lets
  // a body write ``a ` b`` without the span ending at the inner backtick.
  function findMatchingBacktickRun(lineText, fromOffset, runLength) {
    var scanOffset = fromOffset;
    while (scanOffset < lineText.length) {
      var runStartOffset = lineText.indexOf("`", scanOffset);
      if (runStartOffset < 0) {
        return -1;
      }
      var runEndOffset = runStartOffset;
      while (lineText.charAt(runEndOffset) === "`") {
        runEndOffset += 1;
      }
      if (runEndOffset - runStartOffset === runLength) {
        return runStartOffset;
      }
      scanOffset = runEndOffset;
    }
    return -1;
  }

  function recordReferencedTicket(annotationState, referencedTicket) {
    if (Object.prototype.hasOwnProperty.call(annotationState.referencedTicketIds, referencedTicket.id)) {
      return;
    }
    annotationState.referencedTicketIds[referencedTicket.id] = true;
    annotationState.referencedTickets.push(referencedTicket);
  }

  // One pass over a text run. expandTitles is false in any code context — a code
  // run must not be contaminated with prose — and flagMissingIds is false inside
  // a fenced block, where REQ bodies print templates and worked examples, so an
  // id that answers to nothing there is an illustration and not a dead
  // reference. Both mirror board-detail.js's insideCodeSpan / insideFencedBlock
  // split, so the drawer and the paste say the same thing about the same body.
  function annotateMentionRun(runText, expandTitles, flagMissingIds, annotationState) {
    var annotatedText = "";
    var cursorOffset = 0;
    var matchResult;
    // bodyMentionPattern is a shared g-flagged RegExp; the drawer leaves its
    // lastIndex wherever its own scan stopped.
    bodyMentionPattern.lastIndex = 0;
    while ((matchResult = bodyMentionPattern.exec(runText)) !== null) {
      var mentionText = matchResult[0];
      var mentionEndOffset = matchResult.index + mentionText.length;
      bodyMentionPattern.lastIndex = mentionEndOffset;
      if (matchResult[1] || matchResult[2]) {
        // A URL or a repo-relative path is one opaque run. The drawer resumes
        // INSIDE it so a nested id can still link; a clipboard payload must not,
        // or an expansion lands mid-path and the paste carries a path that no
        // longer names a file.
        continue;
      }
      var ticketTarget = resolveTicketMention(mentionText);
      if (!ticketTarget) {
        // Plain text has no red, so the appendix is where a paste's reader
        // learns a reference is dead. Ambiguous is not dead: the board holds
        // records that match and refuses to pick one.
        if (flagMissingIds && !isAmbiguousTicketMention(mentionText)) {
          recordReferencedTicket(annotationState, { kind: "", id: mentionText });
        }
        continue;
      }
      recordReferencedTicket(annotationState, ticketTarget);
      if (!expandTitles ||
        Object.prototype.hasOwnProperty.call(annotationState.expandedTicketIds, ticketTarget.id)) {
        continue;
      }
      // describeTicketTitle never answers empty for a record that resolved — it
      // substitutes "untitled" or names the synthesized-UR case — so there is no
      // empty-title branch to take here.
      var inlineTitle = shortTicketTitle(describeTicketTitle(ticketTarget.kind, ticketTarget.id).text);
      annotationState.expandedTicketIds[ticketTarget.id] = true;
      annotatedText += runText.slice(cursorOffset, mentionEndOffset) + " (" + inlineTitle + ")";
      cursorOffset = mentionEndOffset;
    }
    return annotatedText + runText.slice(cursorOffset);
  }

  // Splits one line into alternating prose and inline-code runs so a backticked
  // id keeps its exact bytes while still earning its appendix line. An unclosed
  // backtick run is ordinary prose, which is what a body writing "the ` glyph"
  // needs.
  function annotateLineOutsideFence(lineText, annotationState) {
    if (lineText.indexOf("`") < 0) {
      return annotateMentionRun(lineText, true, true, annotationState);
    }
    var annotatedLine = "";
    var cursorOffset = 0;
    while (cursorOffset < lineText.length) {
      var runStartOffset = lineText.indexOf("`", cursorOffset);
      if (runStartOffset < 0) {
        break;
      }
      var runEndOffset = runStartOffset;
      while (lineText.charAt(runEndOffset) === "`") {
        runEndOffset += 1;
      }
      var openingRunLength = runEndOffset - runStartOffset;
      annotatedLine += annotateMentionRun(lineText.slice(cursorOffset, runStartOffset), true, true, annotationState);
      var closingRunOffset = findMatchingBacktickRun(lineText, runEndOffset, openingRunLength);
      if (closingRunOffset < 0) {
        annotatedLine += lineText.slice(runStartOffset, runEndOffset);
        cursorOffset = runEndOffset;
        continue;
      }
      var codeSpanEndOffset = closingRunOffset + openingRunLength;
      annotatedLine += annotateMentionRun(lineText.slice(runStartOffset, codeSpanEndOffset), false, true, annotationState);
      cursorOffset = codeSpanEndOffset;
    }
    return annotatedLine + annotateMentionRun(lineText.slice(cursorOffset), true, true, annotationState);
  }

  // Walks a body line by line so a fenced block keeps its exact bytes. The fence
  // lines themselves count as fenced content — an info string is not prose — so
  // opener, contents and closer all take the one suppressed call below.
  function annotateMarkdownBody(bodyText, annotationState) {
    var bodyLines = bodyText.split("\n");
    var openFenceRun = null;
    for (var lineIndex = 0; lineIndex < bodyLines.length; lineIndex += 1) {
      var lineText = bodyLines[lineIndex];
      var fenceRun = codeFenceRunFor(lineText);
      var insideFencedBlock = openFenceRun !== null || fenceRun !== null;
      if (openFenceRun !== null) {
        openFenceRun = codeFenceRunCloses(fenceRun, openFenceRun) ? null : openFenceRun;
      } else if (fenceRun !== null) {
        openFenceRun = fenceRun;
      }
      bodyLines[lineIndex] = insideFencedBlock
        ? annotateMentionRun(lineText, false, false, annotationState)
        : annotateLineOutsideFence(lineText, annotationState);
    }
    return bodyLines.join("\n");
  }

  // One document in; the same document with its BODY annotated out, plus the
  // tickets it referenced in first-mention order. First-mention memory is per
  // document because a payload is a concatenation of whole files, and each file
  // pastes and reads on its own.
  function annotateTicketMentions(documentText) {
    var bodyStartOffset = frontmatterFenceEndOffset(documentText);
    var annotationState = { expandedTicketIds: {}, referencedTicketIds: {}, referencedTickets: [] };
    var annotatedBody = annotateMarkdownBody(documentText.slice(bodyStartOffset), annotationState);
    return {
      text: documentText.slice(0, bodyStartOffset) + annotatedBody,
      referencedTickets: annotationState.referencedTickets
    };
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
  // BEFORE the join, because a joined payload carries N frontmatter fences and a
  // split-off-the-first-fence annotator would rewrite every later ticket's
  // depends_on — the exact corruption this feature must not cause. The join
  // itself keeps cat semantics: exact bytes, no invented separator.
  function annotateClipboardPayload(rawDocuments, excludedIds) {
    var referencedTickets = [];
    var annotatedDocuments = rawDocuments.map(function (rawDocument) {
      var annotatedDocument = annotateTicketMentions(rawDocument);
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
      .then(function (clipboardMarkdown) {
        return annotateClipboardPayload([clipboardMarkdown], [requestedId]);
      })
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
  function rawMarkdownDocumentsForRequests(markdownData, requestIds) {
    return requestIds.map(function (requestId) {
      var rawMarkdown = rawMarkdownForDetail(markdownData, "req", requestId);
      if (rawMarkdown === null) {
        throw new Error("raw Markdown unavailable for " + requestId);
      }
      return rawMarkdown;
    });
  }

  function rawMarkdownDocumentsForUserRequestAndRequests(markdownData, userRequestId, requestIds) {
    var rawUserRequest = rawMarkdownForDetail(markdownData, "ur", userRequestId);
    if (rawUserRequest === null) {
      throw new Error("raw Markdown unavailable for " + userRequestId);
    }
    // The grouped id list is the same all-tree set rendered in the UR drawer.
    return [rawUserRequest].concat(rawMarkdownDocumentsForRequests(markdownData, requestIds));
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
          rawMarkdownDocumentsForUserRequestAndRequests(
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
          return annotateClipboardPayload(rawMarkdownDocumentsForRequests(markdownData, requestIds), requestIds);
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
