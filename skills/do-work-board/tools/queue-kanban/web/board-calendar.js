  // ---- calendar -----------------------------------------------------------

  // calendarDayBreakdown counts a list of calendar entries into ordered,
  // non-zero status groups — what a day label and the subhead actually display,
  // and what board.css colours. Pure and self-contained on purpose: the count
  // line is this view's main at-a-glance signal, and a pure function is what a
  // Node behavior probe can execute directly.
  //
  // Every status is spelled out rather than prefix-matched, for the same reason
  // model.go's isCompletedStatus is an exact match: a typo like
  // "blockd-dependency-cycle" must fall through to the unrecognized group
  // instead of being quietly counted as real blocked work.
  function calendarDayBreakdown(entries) {
    var calendarCountGroups = [
      { group: "done", label: "done", statuses: ["completed"] },
      { group: "with-issues", label: "with issues", statuses: ["completed-with-issues"] },
      { group: "claimed", label: "claimed", statuses: ["claimed"] },
      { group: "pending", label: "pending", statuses: ["pending"] },
      { group: "needs-answers", label: "needs answers", statuses: ["pending-answers"] },
      {
        group: "blocked",
        label: "blocked",
        statuses: ["blocked", "blocked-archive-collision", "blocked-dependency-cycle"]
      },
      { group: "failed", label: "failed", statuses: ["failed"] },
      { group: "cancelled", label: "cancelled", statuses: ["cancelled"] },
      { group: "unrecognized", label: "unrecognized", statuses: [] }
    ];
    var counts = {};
    entries.forEach(function (entry) {
      var matched = null;
      calendarCountGroups.forEach(function (candidate) {
        if (!matched && candidate.statuses.indexOf(entry.status) !== -1) {
          matched = candidate.group;
        }
      });
      matched = matched || "unrecognized";
      counts[matched] = (counts[matched] || 0) + 1;
    });
    var breakdown = [];
    calendarCountGroups.forEach(function (candidate) {
      if (counts[candidate.group]) {
        breakdown.push({
          group: candidate.group,
          label: candidate.label,
          count: counts[candidate.group]
        });
      }
    });
    return breakdown;
  }

  // Renders a breakdown as coloured "12 done  2 cancelled" segments, separated
  // by colour and spacing rather than by punctuation (board.css says why). Used
  // for both the day labels and the subhead, so the subhead doubles as this
  // view's colour key and no separate legend is needed.
  function appendCalendarCounts(target, breakdown) {
    breakdown.forEach(function (part) {
      var partNode = createElement("span", "calendar-count-part", part.count + " " + part.label);
      partNode.dataset.statusGroup = part.group;
      target.appendChild(partNode);
    });
  }

  function renderCalendar() {
    var scroll = document.getElementById("calendar-scroll");
    var summary = document.getElementById("calendar-summary");
    scroll.textContent = "";
    summary.textContent = "";

    var calendar = boardData.calendar || [];
    var shownEntries = calendar.filter(function (entry) {
      return requestMatchesFilters(entry.id);
    });
    if (hasActiveFilters()) {
      summary.textContent =
        shownEntries.length +
        " of " +
        calendar.length +
        " REQ" +
        (calendar.length === 1 ? "" : "s") +
        " match the current filters";
    } else {
      appendCalendarCounts(summary, calendarDayBreakdown(calendar));
    }

    // The calendar arrives in render order — queued band, then dated days
    // newest-first, then undated — so equal day keys are contiguous and grouping
    // is a walk. Days whose entries are all filtered out never flush, so they
    // disappear entirely.
    var currentDayKey = null;
    var currentEntries = null;

    function flushDay() {
      if (!currentDayKey || !currentEntries || currentEntries.length === 0) {
        return;
      }
      scroll.appendChild(makeCalendarDay(currentDayKey, currentEntries));
    }

    shownEntries.forEach(function (entry) {
      if (entry.dayKey !== currentDayKey) {
        flushDay();
        currentDayKey = entry.dayKey;
        currentEntries = [];
      }
      currentEntries.push(entry);
    });
    flushDay();
  }

  function makeCalendarDay(dayKey, entries) {
    var section = createElement("section", "calendar-day");
    section.dataset.dayKey = dayKey;

    var label = createElement("div", "calendar-day-label");
    var dayDate = new Date(dayKey + "T00:00:00Z");
    if (dayKey === "queued") {
      // Not a day: work that has not started. It leads the view because the
      // calendar reads newest-first, so "not yet" belongs above today.
      label.appendChild(createElement("span", "calendar-day-weekday", "up next"));
      label.appendChild(createElement("span", "calendar-day-date", "In the queue"));
    } else if (!isNaN(dayDate.getTime())) {
      label.appendChild(createElement("span", "calendar-day-weekday", calendarWeekdayFormatter.format(dayDate)));
      label.appendChild(createElement("span", "calendar-day-date", calendarDateFormatter.format(dayDate)));
    } else {
      // "undated": a REQ that should carry a day but has no usable stamp.
      label.appendChild(createElement("span", "calendar-day-weekday", "no usable stamp"));
      label.appendChild(createElement("span", "calendar-day-date", "Undated"));
    }
    var count = createElement("span", "calendar-day-count");
    appendCalendarCounts(count, calendarDayBreakdown(entries));
    label.appendChild(count);
    section.appendChild(label);

    var list = createElement("div", "calendar-day-entries");
    entries.forEach(function (entry) {
      var request = requestsById[entry.id];
      var chip = createElement("button", "calendar-chip");
      chip.type = "button";
      chip.dataset.detailKind = "req";
      chip.dataset.detailId = entry.id;
      if (entry.status) {
        chip.dataset.status = entry.status;
      }
      if (request && request.statusUnrecognized) {
        chip.classList.add("is-status-unrecognized");
      }
      // The status is spoken as well as coloured — a chip's hue is the whole
      // signal for a sighted reader, and would be no signal at all without this.
      chip.setAttribute(
        "aria-label",
        entry.id + (entry.status ? ", " + entry.status : "") + (request ? ": " + (request.title || "") : "")
      );
      var instant = formatShortInstant(entry.entryTime);
      if (instant) {
        chip.title = (entry.status === "claimed" ? "claimed " : "resolved ") + instant;
      }
      chip.appendChild(createElement("span", "calendar-chip-id", entry.id));
      if (request && request.title) {
        chip.appendChild(createElement("span", "calendar-chip-title", request.title));
      }
      list.appendChild(chip);
    });
    section.appendChild(list);

    return section;
  }
