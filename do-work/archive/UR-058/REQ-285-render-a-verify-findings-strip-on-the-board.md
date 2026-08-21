---
id: REQ-285
title: Render a verify-findings strip on the board
status: completed
created_at: 2026-08-19T13:47:06Z
claimed_at: 2026-08-21T00:22:39Z
completed_at: 2026-08-21T00:31:24Z
kb_status: pending
commit: fed89c9
route: B
user_request: UR-058
domain: frontend
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: false
suggested_spec:
depends_on: [REQ-284]
maintenance: false
related: [REQ-284, REQ-286]
batch: verify-findings-on-board
write_set: [skills/do-work-board/tools/queue-kanban/web/template.html, skills/do-work-board/tools/queue-kanban/web/board-cards.js, skills/do-work-board/tools/queue-kanban/web/board.css, skills/do-work-board/tools/queue-kanban/web/board.js, skills/do-work-board/tools/queue-kanban/generate.go, skills/do-work-board/tools/queue-kanban/generate_test.go]
---

# Render a Verify-Findings Strip on the Board

## What

Add a second always-visible strip to the board client, modeled on the existing completion-anomalies
strip, that renders `verifyFindings` and `verifySkipped` from the board payload.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [x] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [x] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

Seeing the board should mean seeing the problems. Today the thirteen categories the data-warnings banner
does not cover are invisible to anyone who does not run `verify` from a shell.

## Context

`web/template.html:137-146` and `web/board-cards.js:412-431` are the completion-anomalies strip: a titled
section with a count, a hint line, and per-item cards, sitting outside the view panels so it survives
every view switch, the recently-done window, and the shared filters. That is the shape to copy, and its
visual weight is already calibrated.

## Detailed Requirements

- One new strip: title, count, one card per finding, the category as the badge, the remedy as the muted
  second line.
- Skipped probes render in a collapsed footer on the same strip. They must render — a skipped probe shown
  as nothing reads as "checked and clean".
- Reuse the `.board-anomalies-*` rules rather than writing a parallel palette.
- Hidden when there are no findings and no skips, exactly as the anomalies strip hides itself.
- Exempt from the shared filters and from the 24h/48h/7d window, for the same reason the anomalies strip
  is: a finding must not be hideable by a filter combination.
- Nothing in this REQ parses or re-derives a finding. The client renders the producer's list blindly;
  category suppression already happened in Go (REQ-284).

## Constraints

- Framework-free client, consistent with the rest of `web/`.
- Read-only. No write surface, no repair affordance, no link that mutates queue state.

## Dependencies

REQ-284 — the `verifyFindings` / `verifySkipped` payload fields must exist first.

## Builder Guidance

Firm on placement and shape, latitude on the exact wording of the title and hint line. Keep it visually
subordinate to the anomalies strip when both are showing; anomalies are broken bookkeeping, findings are
process drift.

## Red-Green Proof

**RED prompt/case:** Generate a board from a fixture queue carrying one REQ claimed more than 3h ago and
one `worktree-agent-*` leftover branch, open the generated `index.html`, and look at it.

**Why RED now:** the board renders neither. `board.Warnings` carries no claim-age or worktree entry, and
there is no strip that could show one.

**GREEN when:** the rendered page shows a findings strip with a count of 2, one card reading the
`claim-needs-attention` detail with its remedy underneath, one card for the worktree leftover, and the
strip stays on screen after switching views and after applying a filter that hides the claimed column.

**Validation:** User confirmed (accepted as F1, F7, F8 in the `do-work validate-feedback` triage).

## Full Context

See `do-work/user-requests/UR-058/input.md` for the complete verbatim input and the triage verdicts.

---
*Source: upstream suggestion for `knews2019/skill-do-work`, observed against v0.212.25 — "Suggested shape" item 3.*

---

## Triage

**Route: B** - Medium

**Reasoning:** Builder Guidance is "firm on placement and shape", and the REQ points at the exact lines to model on. Discovery was where the render function is called from and what the view/filter controls are actually named, both of which the browser check needed.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

- `web/template.html:137-146` is the anomalies strip; `web/board-cards.js:493` is `renderAnomaliesStrip`; `web/board.css:505-546` is the `.board-anomalies-*` palette. All three are the shape to copy, as the REQ says.
- **`renderAnomaliesStrip` is called from `web/board.js:70`, not from `board-cards.js`.** The two files share a scope, so a new render function needs a call site in `board.js` — a file the REQ's write set did not list. See D-01.
- View switching is `[data-view-target="…"]` (board, calendar, durations, timeline, testing), not `data-view`. The filters are `#filter-domain`, `#filter-status`, `#filter-done-window`. Needed for the GREEN check, and the first browser probe found nothing until it used the right selectors.
- `createElement(tag, className, text)` is the existing helper and sets text via `textContent`, which is what keeps producer prose from becoming markup.

*Exploration run inline by the orchestrator*

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/web/template.html` (modify) — the strip's markup, including the collapsed skipped-probes footer
- `skills/do-work-board/tools/queue-kanban/web/board-cards.js` (modify) — `renderVerifyFindingsStrip`
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modify) — the few rules that differ from `.board-anomalies-*`
- `skills/do-work-board/tools/queue-kanban/web/board.js` (modify) — one call site (D-01)
- `skills/do-work-board/tools/queue-kanban/generate.go` (modify) — fix the path-reduction regex this REQ's first render exposed (D-02)
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modify) — the regression test for that fix (D-02)

**Files I will NOT touch:** the Go producer's finding logic — the client renders blindly, and suppression stays in Go where REQ-284 put it.

**Acceptance criteria (restated from the REQ):**
1. One strip: title, count, one card per finding, category as badge, remedy as muted second line.
2. Skipped probes render in a collapsed footer on the same strip.
3. Reuse `.board-anomalies-*` rather than a parallel palette.
4. Hidden when there are no findings and no skips.
5. Exempt from the shared filters and the recently-done window.
6. Nothing parses or re-derives a finding.
7. Framework-free, read-only, no write surface.

## Decisions

- **D-01** (DECIDE & STATE): Extended scope to `web/board.js` for a single call line. Reasoning: `renderAnomaliesStrip` — the function this REQ was told to model on — is itself called from there, so a render function that is never called is not the shape being copied. One line, in the block that already calls every sibling renderer.
- **D-02** (ESCALATE): The first render exposed a real defect in **REQ-284's** `reduceAbsolutePaths`, committed earlier in this same run: `remainingAbsolutePath` matched any `/` inside a token, so the already-relative `do-work/calibration-log.tsv` came out as `do-work<path outside this repository>` in `verifySkipped`. Builder chose to fix it here, with a regression test, rather than queue it. Reasoning: it corrupts the very strings this REQ renders, so shipping the strip over it would ship visibly wrong text; and it is a defect in code committed minutes ago in this run, which makes queueing a REQ against it ceremony rather than tracking. RE2 has no lookbehind, so the fix captures the leading boundary and restores it in the replacement. Value: relative paths — the common case after the repo-root reduction — survive intact. Risk: an absolute path that begins a token not preceded by whitespace or an opening delimiter would now be missed; the three-case table test pins the boundaries.
- **D-03** (DECIDE & STATE): The category badge and the `cleanup can fix` marker moved into a `.board-finding-head` flex row after the first screenshot showed "cleanup can / fix" breaking across two lines. Reasoning: half a phrase under a badge reads as a truncated sentence. Found by looking at the rendered page, which is the only way this class of defect is ever found.

## Implementation Summary

**What was done:** Added a second always-visible strip to the board client rendering `verifyFindings` and `verifySkipped` from the payload, modelled on the completion-anomalies strip and deliberately subordinate to it. Each finding is a card with the category as a badge, the detail as body text, the remedy as a muted second line, and a `cleanup can fix` marker only when the producer set `fixable`. Skipped probes render in a collapsed `<details>` footer on the same strip. Fixed the path-reduction regex the first render exposed.

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/web/template.html` (modified) — `#board-findings` section with head, count, hint, cards host, and the `#board-findings-skipped` footer.
- `skills/do-work-board/tools/queue-kanban/web/board-cards.js` (modified) — `renderVerifyFindingsStrip`, rendering the producer's list blindly and setting every string with `textContent`.
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modified) — `.board-findings*` and `.board-finding*` rules; everything shared comes from `.board-anomalies-*`.
- `skills/do-work-board/tools/queue-kanban/web/board.js` (modified) — one call line beside `renderAnomaliesStrip()`.
- `skills/do-work-board/tools/queue-kanban/generate.go` (modified) — `remainingAbsolutePath` boundary-anchored; replacement restores the captured boundary.
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified) — three-case table pinning strip-absolute / keep-relative / replace-outside.

**Tests touched:** one new Go case (D-02). The strip's own behavior is verified by driving the rendered page in a browser, recorded below.

## Qualification

Passed — 6 files verified, 7 acceptance criteria traced, P-A-U confirmed.

- **[UNIFY] audit:** `node --check` clean on both JS files; `gofmt -l .` clean; `go vet ./...` clean; `go test ./...` ok. `maintainer-verify.sh` exits 0, including its **strict JavaScript behavior lane** (`TestMaintainerStrictJavaScriptBehaviorLane`, PASS, 16.7s). Diff grepped for debug artifacts — no `console.log`, no `debugger`. Zero page errors and zero console errors in every browser run.
- **Substantive:** the strip renders real findings from a real generated board; screenshots attached to the evidence below.
- **Requirements traced:** AC1 → the rendered cards; AC2 → the `<details>` footer, visible in every run; AC3 → only six rules differ from `.board-anomalies-*`, each with a stated reason; AC4 → the early return when both lists are empty; AC5 → verified in the browser across five view switches and a filter that hides the claimed REQ; AC6 → the renderer reads `finding.category/detail/remedy/fixable` and derives nothing; AC7 → no framework, no fetch, no mutation affordance.
- **Flowing:** not a stub — two and then three real findings rendered from a real payload.

## Testing

- `bash _dev/tests/maintainer-verify.sh` — exit 0, including the strict JavaScript behavior lane.
- `go test ./...` — ok. `node --check` on both changed JS files — clean.

**Red-green validation** — the REQ's captured RED was "open the generated page and look at it", so it was checked that way, driving Chromium against a generated fixture board (one REQ claimed 5h ago, one `worktree-agent-*` leftover):

| Captured GREEN condition | Result |
|---|---|
| findings strip present with a count of 2 | ✅ count badge reads `2` |
| one card reading the `claim-needs-attention` detail with its remedy underneath | ✅ detail and remedy both present |
| one card for the worktree leftover | ✅ `merged-worktree-leftover`, with `cleanup can fix` |
| strip stays on screen after switching views | ✅ visible in all five: board, calendar, durations, timeline, testing |
| strip stays after a filter that hides the claimed column | ✅ `filter-status=pending` hides REQ-901's card; strip and count unchanged |
| no page errors | ✅ zero `pageerror`, zero console errors |

**Suppression verified in the rendered page, not just the payload.** A second fixture added a completion anomaly (REQ-904, reversed span). Both strips then render together, and the findings strip shows `timestamp-ordering` for REQ-904 while `completion-anomaly` for the *same REQ* appears only in the anomalies strip and the warnings banner. That is the suppression doing exactly its job: one REQ, two finding classes, each surfacing once in the right place.

**Visual review of the rendered page** (the only way two of these were findable):
- First screenshot showed `cleanup can fix` breaking across two lines mid-phrase → fixed (D-03).
- First render showed `do-work<path outside this repository>` in the skipped footer → a real REQ-284 regex defect, fixed with a test (D-02).
- Both strips together confirm the required subordination: findings uses a muted surface and soft border against the anomalies strip's accent tint.

## Review

**Overall: 93%** — Acceptance: Pass

### Requirements Check

| Requirement | Status |
|---|---|
| Title, count, one card per finding, category badge, remedy as muted second line | ✅ Confirmed in the rendered page |
| Skipped probes in a collapsed footer on the same strip | ✅ `<details>`, visible whenever any probe was skipped |
| Reuse `.board-anomalies-*` rather than a parallel palette | ✅ Six new rules, each a deliberate difference |
| Hidden when there are no findings and no skips | ✅ early return |
| Exempt from shared filters and the recently-done window | ✅ Verified across five views and a hiding filter |
| Nothing parses or re-derives a finding | ✅ reads four fields, derives nothing |
| Framework-free, read-only, no write surface | ✅ Board tool stays at three write surfaces |

### Findings

**Important — none.**

**Minor:**

- **M1:** The strip's behavior is pinned by a browser check recorded here, not by an automated test. The repo has a strict JavaScript behavior lane that this passes, but nothing in it asserts the strip specifically — so a future refactor of `renderVerifyFindingsStrip` would not fail any check. The same is true of the anomalies strip it copies, so this matches the area's existing standard rather than lowering it. Worth a REQ if the board client ever grows a DOM-level test harness; not worth inventing one here.
- **M2:** D-02 fixed a defect in code committed earlier in this run rather than queueing it. A reviewer could reasonably want that as its own REQ; the argument for folding it in is that it corrupted the exact strings this REQ renders.

**Nit:**

- **N1:** The hint line ("queue and process problems `queue-kanban verify` detects — each names what to do about it") is longer than the anomalies strip's. It reads fine at 1400px and wraps acceptably narrower.

### Restatement Sweep

Redefined element: none. This REQ adds a strip and fixes a regex; it changes no contract, field meaning, or output shape. `verifyFindings`/`verifySkipped` were defined by REQ-284 and are consumed here exactly as defined.

Swept for statements about what the board renders: `_dev/primes/prime-kanban-board.md` (defers to the tool and to forensics Check 14 — neither changed), root `CLAUDE.md` § Kanban Board Write Surfaces (**checked: still exactly three; this REQ adds no write surface**), `web/template.html`'s own comments (the anomalies strip's comment still describes the anomalies strip correctly, and the new strip carries its own). No stale restatement.

### Acceptance Testing

Driven in Chromium against generated fixture boards, because the captured RED was "open the page and look at it". Three passes: the two-finding fixture matching the REQ's GREEN exactly; a five-view switch plus a filter that hides the claimed REQ; and a three-finding fixture with a completion anomaly present, which put both strips on screen at once and confirmed the suppression and the visual subordination. Zero page errors throughout.

### Scores (on the record — not the headline)

| Dimension | Score |
|---|---|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 85% |
| Scope Discipline | 90% |
| Risk | Low |
| Acceptance | Pass |

Test Adequacy 85% for M1 — the Go half is tested, the DOM half is verified but not pinned. Scope Discipline 90% for the two scope extensions, both recorded as decisions.

### Follow-up REQs Created

None. D-02's fix carries its own regression test; M1 is a standing property of this area rather than a defect this REQ introduced.

## Lessons Learned

**What worked:** Actually opening the page. Two defects were invisible to every other check and both were obvious on sight — a phrase breaking across two lines, and a mangled path in the skipped footer. `node --check` passed, `go test` passed, the strict JS lane passed, and the page was still wrong. For a rendering REQ the screenshot is not a nicety, it is the test.

**What didn't:** The first browser probe used `[data-view]` and found no view buttons at all, silently reporting an empty list instead of failing. The real attribute is `data-view-target`. A probe that finds nothing and says so calmly is indistinguishable from a probe that found nothing because there was nothing — the check should have asserted it found five buttons before using them.

**Worth knowing:** The path-reduction defect (D-02) is the more interesting one. `reduceAbsolutePaths` ran two passes — replace the repo root with a relative path, then strip anything still absolute — and the second pass ate the first pass's output, because RE2 has no lookbehind and a bare `/` matches mid-token. Any two-pass text reduction where the second pass's pattern can match the first pass's product has this shape. The ordering is load-bearing and the boundary capture is what makes it safe; both are commented at the regex.

## Orientation

The board now shows what `queue-kanban verify` finds. A second strip beside completion anomalies lists every finding the producer emits — category, what is wrong, and what to do about it — plus a collapsed footer for probes that could not run, because a skipped probe rendered as nothing reads as "checked and clean". It sits outside the view panels, so no view switch, recently-done window, or filter combination can hide a finding. Lives in the board client (`skills/do-work-board/tools/queue-kanban/web/`), indexed by `_dev/primes/prime-kanban-board.md`.

Leaf change in contract terms — no new field, no renamed concept, no new write surface (the board tool stays at three). It is a significant change in what a person sees: thirteen categories of queue and process breakage that previously reached nobody without a shell now appear on the page everyone already looks at.

Prime staleness spot-check: `_dev/primes/prime-kanban-board.md` — referenced paths still resolve; its parser-lock-step rule is untouched (no frontmatter field involved), and its build-outputs section is unaffected.
