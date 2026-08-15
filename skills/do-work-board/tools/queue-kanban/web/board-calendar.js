  // ---- calendar -----------------------------------------------------------

  function renderCalendar() {
    var scroll = document.getElementById("calendar-scroll");
    var summary = document.getElementById("calendar-summary");
    scroll.textContent = "";

    var calendar = boardData.calendar || [];
    var shownEntries = calendar.filter(function (entry) {
      return requestMatchesFilters(entry.id);
    });
    summary.textContent = hasActiveFilters()
      ? shownEntries.length +
        " of " +
        calendar.length +
        " completed REQ" +
        (calendar.length === 1 ? "" : "s") +
        " match the current filters"
      : calendar.length + " completed REQ" + (calendar.length === 1 ? "" : "s") + " across the archive";

    // The calendar is sorted most-recent-first, so equal day keys are
    // contiguous — group by walking the list. Days whose entries are all
    // filtered out never flush, so they disappear entirely.
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

    var label = createElement("div", "calendar-day-label");
    var dayDate = new Date(dayKey + "T00:00:00Z");
    if (!isNaN(dayDate.getTime())) {
      label.appendChild(createElement("span", "calendar-day-weekday", calendarWeekdayFormatter.format(dayDate)));
      label.appendChild(createElement("span", "calendar-day-date", calendarDateFormatter.format(dayDate)));
    } else {
      label.appendChild(createElement("span", "calendar-day-date", dayKey));
    }
    label.appendChild(
      createElement("span", "calendar-day-count", entries.length + " done")
    );
    section.appendChild(label);

    var list = createElement("div", "calendar-day-entries");
    entries.forEach(function (entry) {
      var request = requestsById[entry.id];
      var chip = createElement("button", "calendar-chip");
      chip.type = "button";
      chip.dataset.detailKind = "req";
      chip.dataset.detailId = entry.id;
      chip.setAttribute("aria-label", entry.id + (request ? ": " + (request.title || "") : ""));
      chip.appendChild(createElement("span", "calendar-chip-id", entry.id));
      if (request && request.title) {
        chip.appendChild(createElement("span", "calendar-chip-title", request.title));
      }
      list.appendChild(chip);
    });
    section.appendChild(list);

    return section;
  }
