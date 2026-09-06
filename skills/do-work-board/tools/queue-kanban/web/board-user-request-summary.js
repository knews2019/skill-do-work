  // ---- user-request progress rollup ---------------------------------------
  //
  // ONE function answers the whole-UR question, and both surfaces render what it
  // returns. The By UR header and the UR drawer stating different numbers for
  // the same user request is the failure this file exists to make impossible, so
  // there is no second counting rule and no second formatter anywhere.
  //
  // TWO RULES THIS FILE DOES NOT GET TO BREAK.
  //
  // It never measures a finished span. The Go side owns the outlier verdict and
  // the origin rule (durations.go's measureImplementationSpan) and ships its
  // answer as hasImplementationSpan / implementationSpanMinutes /
  // implementationSpanReason. Subtracting a completion stamp from a claim stamp
  // here would produce a plausible figure with no outlier rule and no origin
  // correction — for exactly the REQs whose bookkeeping was worst — and it would
  // contradict the `took …` badge on the card directly below it. `completedAt`
  // is not read anywhere in this file. The one piece of clock arithmetic here is
  // the elapsed time of a claim that is STILL RUNNING, which is the same
  // measurement the card's own stopwatch makes, from the same instant.
  //
  // It never reads the wall clock. nowMs arrives as an argument, so one tick
  // states one instant across the header, the drawer, and every card below them.
  //
  // Everything here is also TOTAL by narrowing, never by try/catch: this runs
  // inside the board's only interval callback (see refreshTickingSurfaces), so a
  // throw would kill every ticking surface on the page — and a swallowed
  // exception would hide that bug instead of the freeze it causes, which is
  // worse.

  // Elapsed minutes of a claim that is still running, or null when the stamp
  // cannot be used. Mirrors formatElapsedDuration's rule exactly: a stamp more
  // than the shared allowance ahead of the viewer's clock is a bad stamp, not a
  // young claim. Returning null rather than zero is what lets the caller
  // disclose it instead of silently adding nothing.
  function liveClaimElapsedMinutes(claimedAtText, nowMs) {
    var claimedMs = Date.parse(claimedAtText);
    if (isNaN(claimedMs) || claimedMs - nowMs > futureInstantSkewAllowanceMs) {
      return null;
    }
    return Math.max(0, (nowMs - claimedMs) / 60000);
  }

  // The rollup. Pure: no DOM, no formatting, no wall clock. It reads the UR's
  // COMPLETE membership from userRequestsById[id].requestIds and never the
  // filtered list, so no search, domain, status, window or activity filter can
  // move a figure or its denominator.
  function summarizeUserRequestProgress(userRequestId, nowMs) {
    var summary = {
      userRequestId: userRequestId,
      totalCount: 0,
      successfulCount: 0,
      resolvedCount: 0,
      activeMinutes: 0,
      excludedSpanCount: 0,
      unmeasuredCount: 0,
      skewedClaimCount: 0,
      liveClaimCount: 0,
      remainingMinutes: 0,
      forecastMemberCount: 0,
      knownForecastCount: 0,
      unknownForecastCount: 0,
      overrunForecastCount: 0,
      remainingIsPartial: false
    };
    var userRequest = userRequestsById[userRequestId];
    var requestIds = (userRequest && userRequest.requestIds) || [];
    summary.totalCount = requestIds.length;

    var projection = (boardData.timeline && boardData.timeline.projection) || {};
    var projectionIsConfident = projection.confident === true;

    for (var memberIndex = 0; memberIndex < requestIds.length; memberIndex++) {
      var request = requestsById[requestIds[memberIndex]];
      if (!request) {
        // A membership id the payload does not carry. It stays in the
        // denominator because it is a real member, and both of its figures are
        // evidence the board does not have.
        summary.unmeasuredCount += 1;
        summary.forecastMemberCount += 1;
        summary.unknownForecastCount += 1;
        continue;
      }
      var status = request.status || "";
      if (isCompletedStatus(status)) {
        summary.successfulCount += 1;
      }
      if (isTerminalResolvedStatus(status)) {
        summary.resolvedCount += 1;
      }

      // Live contribution first, because an over-running claim also spends its
      // own forecast down below.
      var liveMinutes = null;
      if (status === "claimed") {
        liveMinutes = liveClaimElapsedMinutes(request.claimedAt, nowMs);
        if (liveMinutes === null) {
          summary.skewedClaimCount += 1;
        } else {
          summary.liveClaimCount += 1;
          summary.activeMinutes += liveMinutes;
        }
      }

      if (request.hasImplementationSpan === true) {
        if (request.implementationSpanReason) {
          // The Go side measured a span and refused it (an assumed pause, or
          // reversed stamps). Disclosed, never added.
          summary.excludedSpanCount += 1;
        } else if (typeof request.implementationSpanMinutes === "number" &&
          isFinite(request.implementationSpanMinutes)) {
          summary.activeMinutes += request.implementationSpanMinutes;
        } else {
          summary.unmeasuredCount += 1;
        }
      } else if (isTerminalResolvedStatus(status) || status === "failed") {
        // Work that has ENDED with no span the Go side would accept: a
        // cancellation (spans are measured for terminal success only), a
        // failure, or a completion whose stamps could not open a span. Real work
        // may have happened and the board cannot say how much, so it is
        // disclosed as unmeasured rather than counted as zero.
        summary.unmeasuredCount += 1;
      }

      if (isTerminalResolvedStatus(status)) {
        continue;
      }
      summary.forecastMemberCount += 1;
      // Three arms, in order: the REQ's own saved p50_active_minutes; else the
      // Timeline's median for this member's effort class, but only while the
      // Timeline calls itself confident; else unknown. A failed member has no
      // estimated next step at all, so it is unknown even carrying a saved
      // figure — the estimate described work that already stopped.
      var estimateMinutes = null;
      if (status !== "failed") {
        if (request.hasEstimateP50ActiveMinutes === true &&
          typeof request.estimateP50ActiveMinutes === "number" &&
          isFinite(request.estimateP50ActiveMinutes)) {
          estimateMinutes = request.estimateP50ActiveMinutes;
        } else if (projectionIsConfident) {
          var medianMinutes = request.effortEstimate === "effort-mechanical"
            ? projection.trivialMinutes
            : projection.normalMinutes;
          if (typeof medianMinutes === "number" && isFinite(medianMinutes)) {
            estimateMinutes = medianMinutes;
          }
        }
      }
      if (estimateMinutes === null) {
        summary.unknownForecastCount += 1;
        continue;
      }
      if (status === "claimed" && liveMinutes === null) {
        // The claim stamp the Active figure already disclosed as skewed. How
        // much of this estimate has been spent is exactly what the board cannot
        // read, so the member's remainder is UNKNOWN. Charging it the full
        // estimate would state a forecast built on a stamp the same rollup just
        // refused, and it would do it silently.
        summary.unknownForecastCount += 1;
        continue;
      }
      summary.knownForecastCount += 1;
      // A running claim has already spent part of its forecast. Floored at zero:
      // a claim that has outrun its estimate owes no negative time. The floor is
      // COUNTED, not just applied — without it a user request whose every member
      // has blown its estimate renders a clean "~0 min", which reads as nearly
      // done and means the opposite.
      var memberRemainingMinutes = estimateMinutes - (liveMinutes || 0);
      if (liveMinutes !== null && liveMinutes > estimateMinutes) {
        summary.overrunForecastCount += 1;
      }
      summary.remainingMinutes += Math.max(0, memberRemainingMinutes);
    }

    summary.remainingIsPartial = summary.unknownForecastCount > 0;
    return summary;
  }

  // ---- summary formatting -------------------------------------------------
  //
  // Minutes and up, through timelineFormatSpanMinutes, never
  // formatElapsedDuration: a whole-UR budget stated to the second is false
  // precision, and coarse labels make the 1 Hz tick nearly free because the text
  // rarely changes. Disclosure is carried by WORDS and SYMBOLS only — "~" for a
  // forecast, "at least" for a known-partial sum, the literal word "unavailable"
  // for a percentage with no denominator, and explicit N-excluded / N-unmeasured
  // / N-unknown / N-over-estimate / clock-skew qualifiers. A fainter ink tone is
  // not a second channel on this board: --ink-faint against --ink-soft measures
  // 1.29:1 in light and 1.82:1 in dark.

  function userRequestSummaryDurationText(minutes) {
    if (minutes > 0 && Math.round(minutes) < 1) {
      return "under a minute";
    }
    return timelineFormatSpanMinutes(minutes);
  }

  function userRequestSummaryPercentageText(count, total) {
    if (total <= 0) {
      return "unavailable";
    }
    return count + "/" + total + " (" + Math.round((count / total) * 100) + "%)";
  }

  function userRequestSummaryActiveText(summary) {
    var qualifiers = [];
    if (summary.excludedSpanCount > 0) {
      qualifiers.push(summary.excludedSpanCount + " excluded");
    }
    if (summary.unmeasuredCount > 0) {
      qualifiers.push(summary.unmeasuredCount + " unmeasured");
    }
    if (summary.skewedClaimCount > 0) {
      qualifiers.push(clockSkewMarkerText);
    }
    var measuredText;
    if (summary.activeMinutes <= 0 && qualifiers.length > 0) {
      // Nothing the board may add, and evidence it knows is missing. "0 min"
      // here would state a measurement nobody made.
      measuredText = "not measured";
    } else if (qualifiers.length > 0) {
      measuredText = "at least " + userRequestSummaryDurationText(summary.activeMinutes);
    } else {
      measuredText = userRequestSummaryDurationText(summary.activeMinutes);
    }
    if (qualifiers.length === 0) {
      return measuredText;
    }
    return measuredText + " (" + qualifiers.join(", ") + ")";
  }

  // "~0 min" with no qualifier is the one reading this figure must never give
  // when it is really "every member has already outrun its estimate". Both
  // qualifiers use the same N-word grammar as the Active figure's excluded /
  // unmeasured counts, so a reader learns one shape and applies it to both.
  function userRequestSummaryRemainingText(summary) {
    if (summary.forecastMemberCount === 0) {
      return "none";
    }
    if (summary.knownForecastCount === 0) {
      return "unknown";
    }
    var qualifiers = [];
    if (summary.unknownForecastCount > 0) {
      qualifiers.push(summary.unknownForecastCount + " unknown");
    }
    if (summary.overrunForecastCount > 0) {
      qualifiers.push(summary.overrunForecastCount + " over estimate");
    }
    var forecastText = "~" + userRequestSummaryDurationText(summary.remainingMinutes);
    if (summary.unknownForecastCount > 0) {
      forecastText = "at least " + forecastText;
    }
    if (qualifiers.length === 0) {
      return forecastText;
    }
    return forecastText + " (" + qualifiers.join(", ") + ")";
  }

  // The five figures, as text, in one place. This array IS the agreement between
  // the two surfaces: the header renders it as a strip and the drawer renders it
  // as meta rows, and neither one composes a string of its own.
  function userRequestSummaryMetrics(summary) {
    return [
      { key: "total", label: "Grouped REQs", value: String(summary.totalCount) },
      { key: "active", label: "Active", value: userRequestSummaryActiveText(summary) },
      { key: "remaining", label: "Remaining", value: userRequestSummaryRemainingText(summary) },
      { key: "successful", label: "Successful", value: userRequestSummaryPercentageText(summary.successfulCount, summary.totalCount) },
      { key: "resolved", label: "Resolved", value: userRequestSummaryPercentageText(summary.resolvedCount, summary.totalCount) }
    ];
  }

  function userRequestSummaryMetricsByKey(summary) {
    var valuesByKey = {};
    userRequestSummaryMetrics(summary).forEach(function (metric) {
      valuesByKey[metric.key] = metric.value;
    });
    return valuesByKey;
  }

  // Marks a rendered value node as ticking. Only a UR with a live claimed
  // contribution carries the attribute, copying makeImplementationSpanNode's
  // deliberate opt-out: a finished sum does not change between ticks, so
  // selecting it every second would be work with no result.
  function markUserRequestSummaryValueNode(valueNode, summary, metricKey) {
    valueNode.dataset.urSummaryMetric = metricKey;
    if (summary.liveClaimCount > 0) {
      valueNode.dataset.urSummaryId = summary.userRequestId;
    }
  }

  // The By UR header's strip. A SIBLING of the head row, never a child of it:
  // .ur-group-head is a <button>, so anything inside it joins the fold control's
  // accessible name, and five metrics would make it announce a paragraph. Its
  // own row also means it can never squeeze the ellipsised .ur-title.
  function makeUserRequestSummaryStrip(summary) {
    var strip = createElement("div", "ur-summary");
    userRequestSummaryMetrics(summary).forEach(function (metric) {
      var metricNode = createElement("span", "ur-summary-metric");
      metricNode.appendChild(createElement("span", "ur-summary-label", metric.label));
      var valueNode = createElement("span", "ur-summary-value", metric.value);
      markUserRequestSummaryValueNode(valueNode, summary, metric.key);
      metricNode.appendChild(valueNode);
      strip.appendChild(metricNode);
    });
    return strip;
  }

  // The second half of the board's tick (see refreshTickingSurfaces in
  // board-core.js, which captures the instant and passes it here). An ATTRIBUTE
  // pass, not a subscriber registry: renderUserRequestLens rebuilds its host on
  // every render and the drawer rebuilds its meta list on every open, so a
  // registry would accumulate references to detached nodes and would need
  // deregistration on close. A node that is no longer in the document is simply
  // not selected here.
  function refreshUserRequestSummaryNodes(nowMs) {
    var summaryNodes = document.querySelectorAll("[data-ur-summary-id]");
    // One rollup per UR per tick, not one per metric node: a 43-member UR would
    // otherwise be walked five times for the header and five more for the drawer.
    var valuesByUserRequestId = {};
    for (var nodeIndex = 0; nodeIndex < summaryNodes.length; nodeIndex++) {
      var summaryNode = summaryNodes[nodeIndex];
      var userRequestId = summaryNode.dataset.urSummaryId;
      var metricKey = summaryNode.dataset.urSummaryMetric;
      if (!userRequestId || !metricKey) {
        continue;
      }
      if (!Object.prototype.hasOwnProperty.call(valuesByUserRequestId, userRequestId)) {
        valuesByUserRequestId[userRequestId] =
          userRequestSummaryMetricsByKey(summarizeUserRequestProgress(userRequestId, nowMs));
      }
      var nextValue = valuesByUserRequestId[userRequestId][metricKey];
      if (nextValue !== undefined && summaryNode.textContent !== nextValue) {
        summaryNode.textContent = nextValue;
      }
    }
  }
