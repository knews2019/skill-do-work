---
id: UR-061
title: Move Durations label placement into the browser and add a browser test lane
created_at: 2026-08-19T14:36:44Z
requests: [REQ-291, REQ-292]
word_count: 9
---

# Move Durations Label Placement into the Browser and Add a Browser Test Lane

## Summary

The user read the solutions report at
`ai-reports/2026-08-19_1345_durations-label-face-robustness/index.html` and chose **O2** —
the option the report itself recommended against on cost grounds — and answered the
objection directly by asking for the browser tests that O2 needs. The report's argument
against O2 was never that it is wrong; it was that it puts a browser in a test path that
has never had one. The user's instruction removes that objection rather than disputing it.

Two deliverables, in order: a browser behavior probe lane beside the existing Node one
(REQ-291), then the move itself — placement leaves Go, measures real text extents in the
browser, and every Go test that asserted against the deleted model is either re-pinned as a
browser probe or explicitly recorded as dropped (REQ-292).

## Extracted Requests

| REQ | Title |
|---|---|
| REQ-291 | Browser behavior probe lane beside the Node behavior lane |
| REQ-292 | Move Durations label placement into the browser and delete the measured-face constants |

## Batch Constraints

- **Order is fixed.** REQ-292 cannot be verified without REQ-291, so it carries `depends_on: [REQ-291]`.
- **No new package-manager dependency if a browser binary alone will do.** Every measured
  number in `durations_test.go` today came from someone running Playwright by hand; the lane
  exists so that stops being a manual ritual, not so npm becomes a build input. The builder
  decides and records which driver it chose and why.
- **A machine with no browser must still pass.** The Node lane's shape is the precedent: run
  when present, print an explicit SKIP line when absent, and a maintainer-strict selection
  that fails rather than skips. `bash _dev/tests/maintainer-verify.sh` exits 0 either way.
- **No test may be deleted without its property being re-pinned or its loss recorded.** Eight
  Go tests assert against the placement model REQ-292 deletes. Silently dropping one is how
  the overprinting defect that started UR-051 comes back.
- **Scope stays on the Durations mark labels.** The tick labels' `--font-mono` stack and the
  axis-title face carry the same open-generic exposure and are deliberately not in this batch.

## Full Verbatim Input

go with o2 that is the propersolution, and add browser tests

[Context: sent with a screenshot of the O2 card from the solutions report, whose
"LOCKS IN" line reads "Replacing them puts a headless browser in the test path;
maintainer-verify.sh runs Go throughout plus one optional Node lane, and has never needed a
browser." The instruction answers that line rather than contesting it.]

[Prior turn, same session, which set up the report: "this is the wrong direction, need a
more robust HTML solution, I'm good even in reengineering the layout if we have to. The
correct way is to run an ai-report where you give me multiple solutions, use /adhd if needed
and then we'll see which way we go. We definetely do not want to manually measure fonts on
different operating systems."]

---
*Captured: 2026-08-19T14:36:44Z*
