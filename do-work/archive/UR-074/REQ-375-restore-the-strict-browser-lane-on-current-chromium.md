---
id: REQ-375
title: '[impact-rule-change] Restore the strict browser lane on current Chromium'
status: completed
status_changed_at: '2026-08-27T13:27:40Z'
created_at: 2026-08-26T14:40:00Z
user_request: UR-074
addendum_to: REQ-374
domain: testing
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: false
impact: impact-rule-change
effort_estimate: effort-substantive
claimed_at: '2026-08-27T13:27:40Z'
route: B
estimate:
  p50_active_minutes: 35
  confidence: medium
  basis:
  - Route B
  - 5-file write set
  - 5 acceptance criteria
  - browser evidence
  - cross-route regression gates
  - full-suite verification
  calculated_at: '2026-08-27T13:27:40Z'
write_set:
- skills/do-work-board/tools/queue-kanban/browser_probe_test.go
- skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go
- skills/do-work-board/tools/queue-kanban/web/board-timeline.js
- _dev/primes/prime-kanban-board.md
- do-work/prose-backlog.md
- skills/do-work-board/tools/queue-kanban/clipboard_browser_probe_test.go
- skills/do-work-board/tools/queue-kanban/user_request_clipboard_browser_probe_test.go
completed_at: '2026-08-27T13:38:28Z'
commit: 54de194b
kb_status: pending
---

# Restore the Strict Browser Lane on Current Chromium

## What

`TestBrowserBehaviorTimelinePointerCaptureWaitsForThePanEngage` fails on Chromium 141.0.7390.37, tripping its own vacuity guard: *"the capture-swallowed outside-release trial never crossed the host pointerleave boundary; the isolator was not exercised and the mutation pair is vacuous."* Because it fails, `TestMaintainerStrictBrowserBehaviorLane` fails, so the whole strict browser lane is unavailable as a gate.

## Context

Discovered while working on REQ-374. Reproduced at HEAD in a separate worktree with REQ-374's diff absent, so it predates that REQ.

The root cause is already recorded in `do-work/prose-backlog.md` as a stale stated reason: `web/board-timeline.js:2571-2574` and `timeline_browser_probe_test.go:2196-2197` both justify the Timeline's pointer capture with "Chromium suppresses the boundary events while a button is held", which did not reproduce on Chromium 1194 — a buttoned exit delivered `pointerleave` to the host four times. The backlog line covers the prose; this REQ covers the failing test, which is not prose-only.

The prime's own note applies: a measured browser behaviour is per-browser, and this probe was calibrated against an engine that behaved differently.

## Red-Green Proof

**RED prompt/case:** `QUEUE_KANBAN_BROWSER=<chromium-141> go test -count=1 -run '^TestMaintainerStrictBrowserBehaviorLane$' .` in `skills/do-work-board/tools/queue-kanban` fails, naming the timeline pointer-capture probe.
**Captured explanation (unverified):** the probe assumes Chromium suppresses boundary events while a button is held; Chromium 141 delivers `pointerleave` to the host, so the trial the mutation pair depends on never runs. **Verified limit:** the exact141 trial did not exercise the host-pointerleave isolator and the guard correctly refused to pass vacuously. That observation alone does not prove a general presence or absence of boundary events; implementation follows the explicit deprecation decision below.
**GREEN when:** ~~that command exits 0 on Chromium 141~~ — voided by the maintainer's 2026-08-27 decision to deprecate Chromium 141. GREEN is now: the strict lane command exits 0 on current stable Chrome, with the probe still exercising the isolator and the guard untouched. The original RED record stands as history of why 141 failed.
**Validation:** Inferred during capture — the failure is reproduced and understood; the fix is not.

## Open Questions
- [x] I discovered this out-of-scope task while working on REQ-374: the strict browser behavior lane fails at HEAD on Chromium 141 because the timeline pointer-capture probe's vacuity guard fires. Should I process this as a new task? → Confirmed: Deprecate Chrome 141 and proceed against current Chromium; no Chrome 141 compatibility repair.
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.
  - **Deferred 2026-08-27:** The maintainer requested an AI report showing the actual interaction, its origin, and why the test exists before deciding. The concern was that the Kanban board does not implement card dragging. Explain the distinction from Timeline panning with authentic UI evidence. Leave this question unchecked and this REQ `pending-answers`; neither investigation nor report generation authorizes a fix. Date obtained under the Timestamp rule's date-only paragraph in `skills/do-work/actions/work-reference.md`.
  - **Decided 2026-08-27:** After reading the report, the maintainer decided: deprecate Chromium 141 and move forward with the newer standard. No compatibility repair; the vacuity guard and both mutation pairs stay unchanged. The REQ resolves by recording the deprecation, correcting the stale "suppresses boundary events" reasoning the investigation exposed, and proving the strict lane green on current Chrome. Date obtained under the Timestamp rule's date-only paragraph in `skills/do-work/actions/work-reference.md`.
  - **Confirmed in clarification 2026-08-27:** The maintainer explicitly reaffirmed the deprecation in this session. The earlier consent hold is resolved, so this REQ returns to `pending` for implementation and validation, not `completed`. The choice is to move to the newer browser baseline instead of repairing Chrome 141 compatibility. Preserve the vacuity guard and both mutation pairs. The report's focused Chrome 151 pass does not establish a whole-lane pass, and deprecating Chrome 141 does not establish historical behavior for every older release. Date obtained under the Timestamp rule's date-only paragraph in `skills/do-work/actions/work-reference.md`.

    > ```
    > I already said to deprecate Chrome 141
    > ```

## Full Context
See `do-work/user-requests/UR-074/input.md` and REQ-374's `## Discovered Tasks`.

## Investigation Before Clarification

2026-08-27: Work paused at the maintainer’s more specific instruction: remain undecided until the concurrent clarification task provides a visual report. No source files changed. The unchanged focused pointer-capture probe passed eleven runs on Chrome 151.0.7922.174. Exact Chrome for Testing 141.0.7390.37 was downloaded into ignored `build/chrome-141` for reproduction. Full investigation is preserved in the current run directory until the report task consumes it. No approval is inferred from these test results.

The clarification session independently reproduced the focused probe's guard failure on Chrome for Testing 141.0.7390.37 (exit 1) and a passing run on Chrome 151.0.7922.174 (exit 0), without changing Timeline code. Both full strict-browser attempts stalled in the separate text-extent probe and were stopped after about 90 seconds; neither is a completed lane verdict. The default canonical maintainer check passed with its optional browser lane explicitly skipped. Authentic captures show Timeline panning, not moving Kanban cards between columns. The in-app capture tool left the pan marker active after its drag call, so the report distinguishes visual movement evidence from the independent Chrome 151 release proof. These findings refine the decision context without approving implementation.

## Reports

- [2026-08-27] [Timeline panning, not card dragging — AI report for the REQ-375 decision](../../../ai-reports/2026-08-27_1428_req-341-timeline-drag-evidence/index.html). Presents completed REQ-341's test, its origin, current screenshots, raw browser-test results, and verification limits. At publication, REQ-375 remained `pending-answers`; the confirmed deprecation decision above now resolves that hold. Date obtained under the Timestamp rule's date-only paragraph in `skills/do-work/actions/work-reference.md`.

## Triage

**Route: B** — Maintainer approved deprecating Chrome141, without a compatibility repair. Existing Timeline assertions stay intact. CurrentChrome dump-DOM shutdown hang is reproduced independently; reuse the existing protocol transport if needed to let the full strict lane finish.

## Plan

Planning not required — focused implementation guided by the request and existing patterns.

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/user_request_clipboard_browser_probe_test.go` (modify; same strict-lane literal expectation repair)
- `skills/do-work-board/tools/queue-kanban/clipboard_browser_probe_test.go` (modify; strict-lane literal expectation repair)
- `skills/do-work-board/tools/queue-kanban/browser_probe_test.go` (modify)
- `skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go` (modify)
- `skills/do-work-board/tools/queue-kanban/web/board-timeline.js` (modify)
- `_dev/primes/prime-kanban-board.md` (modify)
- `do-work/prose-backlog.md` (modify)

**Acceptance criteria (restated from REQ):**
- All Detailed Requirements and the captured Red-Green Proof are satisfied.

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]:** Read board prime/lessons and testing/debugging/always-on crew rules. Preserve the approved deprecation and both Timeline mutation pairs/vacuity guard. Review and correct pre-existing comment changes without unsupported claims about every older engine. Reproduce current Chrome's separate dump-DOM hang, then reuse the existing protocol transport for deterministic result-node completion if needed. Run focused current-browser probes, the whole strict browser lane, independent review, and canonical verification before completion.
- [x] **[APPLY]:** Implemented the approved deprecation and transport repair within the declared/amended write set; no Timeline runtime or guard changes.
- [x] **[UNIFY]:** Parent inspected all seven diffs and cross-REQ expectations, verified Timeline non-comment byte equality, and ran native syntax, browser and qualification checks. Final gate results are recorded below.

## Exploration

The maintainer's direct confirmation, delivered through clarification, resolves the original consent hold. It authorizes deprecating Chrome141, not historical claims about all prior versions and not weakening the Timeline test. Existing uncommitted comments in the prime, Timeline source/test and prose backlog were not produced by this task or the clarification task; they are preserved as input and now fall within this approved request's review scope. Correct their unsupported assertions before integration.

The installed browser reports Google Chrome151.0.7922.174. A tiny local page with no board scripts prints complete DOM under the exact existing dump-DOM flags but its main process fails to exit within10seconds. Adding no-first-run/no-default-browser-check does not change this. Redirecting stdout/stderr to regular files also times out, ruling out a CombinedOutput pipe-reader-only hang. Each diagnostic used a fresh temporary profile and terminated only its own process group. Earlier actual text-extent and strict-lane attempts stalled90seconds; these are not passing evidence.

The repo already has a persistent DevTools-pipe session with explicit cleanup, page-URL checks, result decoding and bounded condition waits. Reuse that existing mechanism for result-node probes instead of adding a package, a second driver or a kill-and-assume-success workaround. Keep every fixture/assertion, the zero-probe guard and both Timeline mutation pairs. Completion still requires the full strict lane on the newer browser. Browser behavior varies by build; do not infer a successful whole lane from the focused Timeline pass.

### Strict-lane scope amendment

Once the measurement transport could finish, the whole browser suite exposed an old Copy-all literal expectation still using unmarked parentheses. This optional test had been hidden behind the dump-DOM stall during REQ-389. Add clipboard_browser_probe_test.go before editing its two expected body literals to the approved ASCII arrow; keep appendix/exclusions/order assertions intact. The transport also exposed HTML entity serialization in result-node text: serialized outerHTML turns the arrow into an entity. Read the result node's textContent directly over the existing protocol and retain the nonempty/JSON-object guards, removing the now-unneeded serialized-HTML extractor rather than layering another decoder.

The complete initial current-browser suite finished in78.558s with exactly two failures: Board-column Copy and UR Copy retained the old unmarked expected titles, and the serialized-result transport encoded the arrow. All Timeline probes passed, including both unchanged mutation pairs. Add the matching UR clipboard browser test to the declared scope before updating its expected in-body markers. This is the same prior-REQ expectation repair, not a new product behavior.

## Implementation Summary

- `skills/do-work-board/tools/queue-kanban/browser_probe_test.go` (modified). Reuses the existing DevTools-pipe session for measurement probes, waits for a populated result node on the page's real clock, reads literal textContent, and retains object-shape/caller JSON validation and strict probe counting. Removes the dump-DOM command and serialized-HTML extractor. Adds one actual-browser literal-text regression. Explicit per-probe cleanup reuses the existing idempotent close.
- `skills/do-work-board/tools/queue-kanban/clipboard_browser_probe_test.go` (modified). Updates two copied-body expectations to REQ-389's approved ASCII arrow, exposed once the complete browser suite could run. Appendix, ordering and exclusion assertions remain unchanged.
- `skills/do-work-board/tools/queue-kanban/user_request_clipboard_browser_probe_test.go` (modified). Updates the same in-body marker in the UR Copy/Copy-all expectation; full title/status appendix checks remain intact.
- `skills/do-work-board/tools/queue-kanban/timeline_browser_probe_test.go` (modified). Comments only: distinguishes the exact Chrome141 failed isolator trial from general boundary-event behavior, explains synthetic lost-capture coverage, and identifies former dump-DOM limitations as historical. Both mutation pairs and every assertion condition/message are unchanged.
- `skills/do-work-board/tools/queue-kanban/web/board-timeline.js` (modified). Comments only: explains capture's outside-host pointerup routing and lost-capture teardown without an unsupported blanket boundary-event claim. Runtime JavaScript is unchanged.
- `_dev/primes/prime-kanban-board.md` (modified). Records the current Chromium target and confirmed Chrome141 deprecation without inventing a range of historical engine behavior.
- `do-work/prose-backlog.md` (modified). Marks the corrected boundary-event explanation drained by this request; other entries are unchanged.

## Decisions

- **D1 — Deprecate, do not repair Chrome141:** This is the maintainer's confirmed choice, recorded above. No compatibility branch, browser-version skip, altered mutation, relaxed vacuity guard or runtime Timeline change is introduced.
- **D2 — Reuse the protocol connection:** Independent minimal-page replay confirmed that the current Chrome process did not exit after emitting dump-DOM output. First-run flags and regular-file output did not resolve it. The repository already had a protocol session with URL checks, bounded exchanges/readiness and cleanup. Sharing it deletes the stuck completion mechanism without adding a driver/package. Legacy virtual-time flags are filtered at the measurement wrapper; explicit completion uses the existing bounded real-clock condition wait instead.
- **D3 — Read data, not HTML serialization:** A new actual-browser RED showed literal arrow/markup/entity-like text being changed by outerHTML serialization. Read the result node's textContent over the protocol, keep the object-prefix and caller JSON checks, and remove the obsolete HTML extractor. This preserves text such as literal entity spellings without double decoding.

## Testing

### Diagnosis and negative evidence

- Existing report evidence: exact Chrome141.0.7390.37 focused Timeline test failed its unchanged vacuity guard; focused Chrome151.0.7922.174 passed. Both earlier full strict attempts stalled in text-extent testing and were terminated around90seconds. No prior full-browser pass was inferred.
- Parent minimal local HTML under Chrome151 printed complete DOM but the main process did not exit within10seconds. The same failure persisted with no-first-run/no-default-browser-check flags and with stdout/stderr redirected to regular files. Fresh temporary profiles and only each diagnostic's own process group were used. This separated the shutdown problem from board logic and from a pipe-reader-only hang.
- First preflight selection accidentally named a nonexistent extraction test and was not accepted as test evidence. Repeated preflight used the existing zero-probe rejection test and passed. Pre-existing comment edits were explicitly inventoried and reviewed under this approved scope rather than silently staged.
- After transport reuse, the real text-extent test exited0 (2.118s), measuring98.16×13.00 at11px. The initial complete browser run exited1 (78.558s), with exactly two old Copy-all marker/entity expectations failing. All Timeline probes passed, including the unchanged mutation pairs.
- New literal-result test RED: exit1 (1.460s). Markup-like text, literal entity spelling and the ASCII arrow were HTML-encoded by the intermediate serialization. The test now passes with literal textContent.

### Positive evidence

- Four Chrome151 probes (text extent, literal result, Board Copy-all and Timeline pointer capture): exit0 (6.199s). Actual Timeline result: ordinary click opened REQ-066; early-capture mutation did not; normal outside release ended panning; swallowed-capture mutation left panning true with one host-pointerleave event isolated.
- Both Copy-all probes plus literal result after the final marker corrections: exit0 (3.913s). The same full appendix expectations passed.
- Exact required command: `QUEUE_KANBAN_BROWSER='/Applications/Google Chrome.app/Contents/MacOS/Google Chrome' go -C skills/do-work-board/tools/queue-kanban test -count=1 -run '^TestMaintainerStrictBrowserBehaviorLane$' -v -timeout 15m .` completed exit0 (74.775s). This is the whole strict lane, not merely a focused test.
- Independent review ran five focused Chrome151 acceptance tests, exit0 (7.331s), and independently ran the zero-probe rejection test, exit0 (13.053s). The latter verifies the internally launched missing-browser lane still fails with its exact diagnostic.
- Parent and reviewer compared both Timeline files to HEAD after removing only full-line comments: every remaining byte is identical. Runtime source hash: f5348ff6257c4b1a5fa588a2c4f5f8d634961fd0ad4a357102f16624edac6e8a. Test non-comment hash: d6ce79f5699b124d10bf0a608b7b538ade2851a602f53e0f1b64ce8c68a062c6. Rechecked after the final historical-comment corrections.
- JavaScript syntax, gofmt and diff checks passed. A final strict repeat and canonical gate cover the final comment/release state; their actual results are appended below when completed.

## Review Orientation

Inspect the measurement wrapper's shared protocol call, result-node readiness, literal text decoding and cleanup. Probe counting remains in the shared session; missing browsers still skip ordinary checks and fail the strict zero-probe guard. No assertion or mutation was removed to obtain the current-browser pass. Read the two prior clipboard test changes as intentional representation updates, not changes to their order/source/appendix invariants.

## Discovered Tasks

None outstanding. The stale Copy-all expectations and result serialization were prerequisites uncovered while restoring this strict lane and are documented above. The remaining prose backlog entries were not changed.

## Lessons Learned

Browser process exit is not a reliable result-readiness signal: this Chrome build emitted complete dump-DOM output but did not exit, even for a tiny local page. A bounded protocol read of the completed result node restores observability without relaxing product assertions. Read textContent when the contract is JSON text; serializing HTML can silently change literal clipboard content. A failed mutation-isolator trial identifies that trial's missing event, not all event behavior in that engine or every older release.

## Qualification

Mechanical qualification exit0. Parent verified all seven files against the approved scope and amended test dependencies. Data flows retain literal result text, explicit readiness, URL identity, probe accounting and the original caller assertions. Both Timeline files have byte-identical non-comment content against HEAD. New comments state only supported policy and observed trial evidence. The original report bundle is unchanged; its relative link is repointed only as the completed REQ moves into the UR archive.

## Review

**Overall: 100%** | 2026-08-27T13:36:52Z

| Dimension | Score |
|---|---|
| Requirements | 100% |
| Code Quality | 100% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings:** None outstanding.
**Minor findings:** Three current-tense dump-DOM comments corrected to describe historical limitations/current result-node probes. Non-comment equality was independently checked and rechecked after correction.
**Acceptance:** Pass — independent five-probe Chrome151 run exited0 (7.331s), independent zero-probe rejection exited0 (13.053s), and parent exact full strict lane exited0 (74.775s). Reviewer inspected the full strict output and the seven-file diff, confirming unchanged Timeline mutations and assertions.
**Suggested testing:** Final canonical gate, owned by parent.
**Follow-ups created:** None; **sweeps appended to:** None.

*Reviewed by review-work action; independent reviewer.*

### Final browser verification

After the final comment corrections, the same exact strict command completed again with exit0, package76.994s. This repeat used Chrome151.0.7922.174 and covered the final source state, including every browser behavior probe and the nonzero-probe guard. No Chrome141 compatibility success is claimed or required under the confirmed deprecation.

## Repository Verification

`bash _dev/tests/maintainer-verify.sh` completed with exit 0 on the final implementation/release state. Contract suite, Go checks, and strict JavaScript lane passed. The default optional browser lane was explicitly skipped; separate browser evidence is recorded above where applicable.
