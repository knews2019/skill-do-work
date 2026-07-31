---
id: REQ-034
title: Capture-time decomposition nudge + board write-set overlap badges
status: completed
commit: 9c7e3be
completed_at: 2026-07-28T23:18:24Z
claimed_at: 2026-07-28T22:43:13Z
created_at: 2026-07-28T20:52:05Z
user_request: UR-007
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: [REQ-032]
related: [REQ-032, REQ-033]
batch: parallel-dispatch
write_set:
  - actions/capture-reference.md
  - actions/capture.md
  - tools/queue-kanban/model.go
  - tools/queue-kanban/generate.go
  - tools/queue-kanban/model_test.go
  - tools/queue-kanban/web/board.js
  - tools/queue-kanban/web/board.css
  - actions/board.md
  - tools/queue-kanban/prime-do-kanban.md
  - CLAUDE.md
maintenance: false
---

# Capture-time decomposition nudge + board write-set overlap badges

## What

Two upstream levers so future batches parallelize without runtime machinery: (1) capture's slicing guidance nudges toward REQ boundaries that give each REQ its own files; (2) the Kanban board surfaces write-set overlaps between queued/claimed REQs as a badge.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `crew-members/general.md`, `crew-members/coding-guardrails.md`, `crew-members/ui-design.md`, `crew-members/anti-slop.md`, and the Exploration's line anchors. Approach, in dependency order:
  1. `tools/queue-kanban/model.go` — derived `WriteSetOverlaps []string` on `RequestTicket` + `annotateWriteSetOverlap(tickets)` following `annotateDependencyState`'s compute-in-Go pattern (model.go:1076): pairwise over pending/claimed tickets only, `writeSetPatternsIntersect` = literal equality ∪ `filepath.Match` both directions. Called from `buildBoard` **after** `bucketColumns` so it is structurally impossible for column logic to read it. Update the `WriteSet` field comment (model.go:115-118), which currently asserts "never used for … overlap computation".
  2. `tools/queue-kanban/generate.go` — `WriteSetOverlaps` payload field (`writeSetOverlaps`, omitempty) + assignment in `buildGeneratedBoardData`.
  3. `tools/queue-kanban/model_test.go` — annotation tests: disjoint, literal overlap, glob↔literal both directions, absent/empty set ⇒ no annotation, non-pending/claimed excluded, self never listed.
  4. `tools/queue-kanban/web/board.js` — card badge via the existing `makeBadge`/`truncateBadgeText` helpers (rendered from the derived array; no cross-card computation in JS), plus drawer rows "Write set" and "Overlaps". ES5/`var` style, no `console.log`.
  5. `tools/queue-kanban/web/board.css` — `.badge-write-overlap` beside `.badge-unblocks`, with the same explanatory comment block.
  6. Docs in lock-step: `actions/board.md:115`, `CLAUDE.md:158`, `tools/queue-kanban/prime-do-kanban.md:51` all say "no overlap computation on the board's side" and must become "display-only overlap annotation (badge + drawer), never column logic — co-dispatch decisions stay with the work pipeline's gate". board.md:115 also mis-attaches `write_set` to the blocked badge (grammar fix).
  7. Capture nudge: 1-2 sentences appended to `actions/capture-reference.md:115`'s Slicing convention paragraph (no new section); `actions/capture.md:213` field enumeration gains `write_set`.
- [x] **[APPLY]:** Applied as planned across all 10 declared files, in two sittings — the first builder died mid-build on a transient API error after finishing the Go/JS half (items 1–5); a continuation builder assessed that half against the plan, found it complete and compiling, and finished the doc half (items 6–7) without redoing it. No plan step was dropped or re-scoped. `git diff --stat` shows exactly the 10 Scope files, no eleventh.
- [x] **[UNIFY]:** `git diff --stat` reviewed file by file; `go test ./...`, `go vet ./...`, `gofmt -l`, and `bash _dev/tests/contract-regressions.sh` all green; `git diff | grep -E 'console\.(log|debug|warn)|debugger|fmt\.Print|TODO|FIXME|XXX|HACK'` returns nothing (the inherited Go/JS half is clean — no debug artifacts, no `console.log`, no scaffolding left behind). Per file:
  - `tools/queue-kanban/model.go` (inherited) — `WriteSetOverlaps` is documented as derived-never-parsed beside the existing derived fields; `annotateWriteSetOverlap` is called from `buildBoard` **after** `bucketColumns` with a comment saying why; the stale "never used for … overlap computation" clause on the `WriteSet` field comment was updated rather than left contradicting the new code; `filepath` was already imported (no orphaned import), `sortRequestIdList` is an existing helper (no duplicate sorter).
  - `tools/queue-kanban/generate.go` (inherited) — one payload field + one assignment; the rest of that hunk is gofmt realigning the struct's tag column, not a semantic change (`gofmt -l` clean confirms). `omitempty` keeps the payload byte-identical for every board with no declared overlaps.
  - `tools/queue-kanban/model_test.go` (inherited) — three tests, all verified passing by name: fixture-table (literal, glob, disjoint, empty, `completed`/`reserved`/`pending-answers` exclusion, self-never-listed), both-direction glob matching + a negative glob case, and an end-to-end `buildBoard` test asserting both contending REQs stay in `PendingReady` (the display-only guarantee, tested rather than asserted in prose).
  - `tools/queue-kanban/web/board.js` (inherited) — badge renders only from `request.writeSetOverlaps`; no cross-card computation in JS. ES5 `var` style matches the file, uses the existing `makeBadge`/`truncateBadgeText`/`appendMetaRow`/`makeTicketLinkList` helpers, no new globals, no `console.*`.
  - `tools/queue-kanban/web/board.css` (inherited) — one class beside `.badge-unblocks`, comment block in the same house style; no existing selector modified.
  - `actions/capture-reference.md` — two sentences appended inside the existing Slicing convention paragraph; no new section, no heading, paragraph still reads as one thought.
  - `actions/capture.md` — the line-213 enumeration now names `write_set` and points at Populating `write_set`; the "when the sliced REQs depend on each other" conditional was re-attached to `depends_on` only (it was governing the whole clause).
  - `actions/board.md` — lock-step sentence rewritten (grammar fix + truthful annotation description); new badge-behavior paragraph after the mode list. Read-only preamble (line 5, "writes exactly three things") re-read and deliberately left alone — the annotation is computed in memory during board build and adds no write surface.
  - `tools/queue-kanban/prime-do-kanban.md` — lesson entry rewritten to describe what the code actually does (after-bucketing call site, `reserved` exclusion, both-direction `filepath.Match`, glob-vs-glob caveat, empty-set reading).
  - `CLAUDE.md` — lock-step claim now says display-only annotation instead of "never computes overlap"; blocked-check clause left intact.
  - Cross-check for stale echoes (Closed Enumerations rule): `grep -rn "overlap" actions/ docs/ tools/queue-kanban/ CLAUDE.md SKILL.md _dev/ crew-members/` plus a repo-wide `write_set|writeSet` sweep over `*.md`. Every remaining hit is either the pipeline dispatch gate (`actions/work.md`, `actions/work-reference.md`, `docs/work-guide.md` — unchanged and still accurate) or an unrelated sense of the word. The three "no overlap computation" echoes were exactly the three declared files.

## Why (if provided)

The observed wave-1 serialization was caused at capture time — three REQs sliced so they all wrote one 16-line CSS block. A boundary nudge at capture buys more parallelism than any dispatch mechanism, and a visible overlap badge would have made the constraint obvious on the board before the run started.

## Detailed Requirements

- `actions/capture.md` Step 1 / `actions/capture-reference.md` Slicing convention: when splitting one surface into multiple REQs, prefer boundaries that give each REQ its own files (per-concern files over shared blocks); when overlap is unavoidable, say so in the REQs' `write_set`/Scope so the dispatcher can see it. A nudge, not a gate — intent-preserving slicing still wins when file boundaries would distort the user's request.
- Board: pending/claimed cards whose `write_set` intersects another pending/claimed REQ's get an overlap badge naming the other REQ id(s). Display only — no column logic, no blocking; same badge treatment as existing `depends_on` display.
- Board change follows the lock-step rule (parser fields from REQ-032; badge logic in the web frontend); no independent versioning — normal skill changelog entry.

## Constraints

- Capture guidance must not inflate capture into architecture: the nudge is one or two sentences in the existing slicing convention, not a new section.
- Board work degrades gracefully when `write_set` is absent (no badge — absence means unknown, not conflict).

## Dependencies

Depends on REQ-032 (`write_set` field must exist and be parsed).

## Builder Guidance

Certainty: Firm on both deliverables; latitude on badge visual treatment. Keep the capture wording short — this is a nudge whose whole value is being read.

## Red-Green Proof

**RED prompt/case:** The Slicing convention in `actions/capture-reference.md` says nothing about file boundaries; two queue REQs declaring overlapping `write_set` paths render on the board with no indication they contend.
**Why RED now:** Nothing upstream discourages shared-block slicing, and the board can't show contention it doesn't parse.
**GREEN when:** The slicing convention carries the boundary preference; seeding two queue REQs with overlapping `write_set` values shows an overlap badge on both cards naming each other.
**Validation:** User confirmed (design agreed in conversation; "run the capture").

## Full Context

See `do-work/user-requests/UR-007/input.md` for complete verbatim input.

---
*Source: "can we do something to update this skill to encourage them to run in parallel?" (UR-007)*

Think carefully before answering.

---

## Triage

**Route: B** - Medium

**Reasoning:** Both levers have a clear "what" (a 1-2 sentence slicing nudge; a display-only overlap badge fed by the `writeSet` payload REQ-032 already ships) — the open questions are pattern-level: how the frontend renders existing badges, and where the slicing convention text sits. Exploration suffices; no architectural planning needed.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration


_Raw Explore-agent session transcript (833 KB of JSONL) was pasted here instead of the agent's summary. Removed when this repo began tracking `do-work/` (0.157.0) rather than committing verbatim session capture — session UUIDs, prompt text, and local paths — into permanent git history. The REQ's Scope, Implementation, and Review sections below are intact and carry the actual decision trail._

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `actions/capture-reference.md` (modify) — 1-2 sentence slicing nudge appended to the Slicing convention paragraph (line 115): prefer boundaries giving each REQ its own files; unavoidable overlap gets said in `write_set`/Scope
- `actions/capture.md` (modify) — line 213 complex-batch enumeration gains `write_set` (closes REQ-032 review Minor 5's stale enumeration)
- `tools/queue-kanban/model.go` (modify) — `annotateWriteSetOverlap` derived annotation (pending/claimed tickets; pairwise intersection with `filepath.Match` for glob-vs-literal; follows the `annotateDependencyState` compute-in-Go pattern), derived field on `RequestTicket`
- `tools/queue-kanban/generate.go` (modify) — `writeSetOverlaps` payload field + assignment
- `tools/queue-kanban/model_test.go` (modify) — overlap annotation tests (disjoint, literal overlap, glob-vs-literal, absent set ⇒ no badge, non-pending/claimed excluded)
- `tools/queue-kanban/web/board.js` (modify) — overlap badge on pending/claimed cards naming the other REQ id(s) (render from the derived array, `makeBadge` pattern); drawer rows for Write set / Overlaps
- `tools/queue-kanban/web/board.css` (modify) — `.badge-write-overlap` class in existing badge style
- `actions/board.md` (modify) — rewrite the line-115 lock-step sentence (fixes REQ-032 review Minor 4's blocked-badge mis-attachment; "no overlap computation" becomes "display-only overlap annotation, never column logic — the dispatch gate stays the pipeline's") + one badge-behavior sentence
- `tools/queue-kanban/prime-do-kanban.md` (modify) — line ~51 echo of "no overlap computation" updated
- `CLAUDE.md` (modify) — line ~158 lock-step claim updated (display annotation vs dispatch gate)

**Files I will NOT touch:** `SKILL.md`, `actions/work.md` / `actions/work-reference.md` (dispatch gate untouched — REQ-035/036/037/038 own that surface), `template.html` (drawer rows are JS-generated), `_dev/tests/contract-regressions.sh` (existing `fields["write_set"]` ratchet already covers the parse; badge is optional UI — conscious skip unless the builder finds a cheap anchor).

**Acceptance criteria (restated from REQ):**
- [x] Slicing convention carries the boundary preference in 1-2 sentences (a nudge, not a gate; intent-preserving slicing still wins) — `actions/capture-reference.md`, appended to the existing paragraph
- [x] Unavoidable overlap is declared in `write_set`/Scope so the dispatcher can see it — same two sentences; `actions/capture.md:213` now names `write_set` in the complex-batch enumeration
- [x] Board: pending/claimed cards whose `write_set` intersects another pending/claimed REQ's show an overlap badge naming the other REQ id(s) — `annotateWriteSetOverlap` → `writeSetOverlaps` payload → `.badge-write-overlap`
- [x] Display only — no column logic, no blocking; absence of `write_set` means no badge (unknown ≠ conflict) — enforced structurally (annotation runs after `bucketColumns`) and asserted by `TestWriteSetOverlapNeverAffectsColumnPlacement`
- [x] Same badge treatment as existing `depends_on`-style display; parser/payload lock-step in the same commit — model.go/generate.go/board.js/board.css plus `actions/board.md`, `CLAUDE.md`, `prime-do-kanban.md` all in one diff
- [x] `go test ./...` green including new overlap tests; `bash _dev/tests/contract-regressions.sh` green; GREEN check: two queue REQs with overlapping write_sets render badges naming each other — see **Verification** below

## Verification

```
$ cd tools/queue-kanban && go test ./...
ok  	github.com/knews2019/skill-do-work/queue-kanban	0.843s

$ go test -run WriteSetOverlap -v ./...
--- PASS: TestAnnotateWriteSetOverlapPairsContendingRequests (0.00s)
--- PASS: TestAnnotateWriteSetOverlapMatchesGlobsInBothDirections (0.00s)
--- PASS: TestWriteSetOverlapNeverAffectsColumnPlacement (0.00s)
PASS

$ go vet ./...            # no output
$ gofmt -l tools/queue-kanban/   # no output

$ bash _dev/tests/contract-regressions.sh
Contract regression checks passed.
```

**GREEN check (scratch queue, deleted afterward).** Three pending REQs seeded into a throwaway `do-work/queue/`: REQ-901 `[…/board.css, …/board.js]`, REQ-902 `[…/web/*.css]` (glob), REQ-903 `[docs/some-guide.md]` (disjoint). Ran `queue-kanban generate` against it and read `board-data.js`:

```
REQ-901 | writeSet= ['tools/queue-kanban/web/board.css', 'tools/queue-kanban/web/board.js'] | writeSetOverlaps= ['REQ-902']
REQ-902 | writeSet= ['tools/queue-kanban/web/*.css']                                        | writeSetOverlaps= ['REQ-901']
REQ-903 | writeSet= ['docs/some-guide.md']                                                  | writeSetOverlaps= None
```

Each contending REQ names the other (glob against literal, both directions); the disjoint one is annotated with nothing. The generated `index.html` inlines both `badge-write-overlap` and `writeSetOverlaps`, so the rendering path ships in the static artifact too. Scratch dir and both built binaries deleted; `git status --porcelain -uall` shows only the 10 declared files.

## Decisions

**D-01 — Pairwise overlap is computed in Go, not in `board.js`. (DECIDE & STATE)**
Detailed Requirements say "badge logic in the web frontend". Read literally that puts the first-ever cross-card computation into JS, where nothing can test it — there is no JS test runner, and `go test` only asserts substrings of the inlined `board.js`. Every other cross-card relationship already works the other way: `annotateDependencyState` computes `UnmetDependencies`/`Dependents` in Go and JS renders the derived array for the `unblocks` badge. Read "badge logic in the frontend" as *rendering* being the frontend's job and put the intersection in `annotateWriteSetOverlap`. Buys three real tests; costs nothing the REQ asked for.

**D-02 — `annotateWriteSetOverlap` is called AFTER `bucketColumns`, deliberately. (DECIDE & STATE)**
`annotateDependencyState` runs *before* bucketing because bucketing consumes it. This one is the opposite requirement — it must never reach column logic — so it runs after, and column code physically cannot read a field that isn't populated yet. Makes "display only" structural instead of a comment someone can violate later. Both the call site and the field carry a comment saying so.

**D-03 — Glob semantics: literal equality ∪ `filepath.Match` in both directions; glob-vs-glob compared as literals only. (DECIDE & STATE)**
Nothing in the tool matched globs before (`write_set` was read verbatim). Exact equality alone would miss `web/*.css` vs `web/board.css` — precisely the contention the badge exists to surface — so each pair is tried as pattern-against-text both ways, which also makes badge presence independent of REQ ordering. Two globs are **not** intersected: the stdlib has no pattern-intersection primitive and hand-rolling one for a display hint is over-build. Consequence, stated in the `writeSetPatternsIntersect` doc comment and in `prime-do-kanban.md`: a glob-vs-glob pair that shares files renders no badge. Acceptable because the safety-relevant reader is the pipeline's gate, which treats an unexpandable glob as overlapping (`actions/work.md`), not this annotation.

**D-04 — `reserved` is excluded from the compared tier even though it shares the Claimed column. (DECIDE & STATE)**
The REQ says pending/claimed. `reserved` means allocated to another worktree/cloud session (0.125.0) — the board cannot see what that session is actually writing, and its declared set may already be stale there. A badge against it asserts contention the operator can neither verify nor act on. Only `pending` and `claimed` — the tier a dispatcher could still put in flight together — are compared; terminal and needs-input tiers are out for the same reason. Encoded in `isWriteSetOverlapCandidateStatus`, pinned by the fixture-table test.

**D-05 — An empty or absent `write_set` never overlaps on the board — the OPPOSITE of the pipeline gate's reading. (DECIDE & STATE)**
The gate reads absent as "overlaps everything ⇒ serialize" because a false *disjoint* there corrupts files. On the board a false *badge* on every card that never declared a set would cry wolf and train the operator to ignore the badge entirely. Both readings are correct for their consumer, so both are documented where they live (`writeSetsIntersect` doc comment; `actions/board.md`'s badge paragraph spells out that no badge means no *declared* overlap — unknown, not safe — and names the gate's opposite reading so the two can't be conflated).

**D-06 — Badge visual: pending accent, not the blocked/alarm treatment. (DECIDE & STATE — Builder Guidance granted latitude here.)**
`.badge-write-overlap` reuses `--accent-pending`/`--tint-pending` plus the mono font and the `max-width: 100%` truncation partner from `.badge-blocked`. Contention is something to schedule around, not a failure; the blocked treatment would over-signal for a display-only hint. Label `overlaps`, value = the id list truncated to 24 chars, with the full list and the "display only — `do-work run`'s dispatch gate decides" caveat in the `title` tooltip.

**D-07 — The drawer gains two rows: "Write set" and "Overlapping write sets". (DECIDE & STATE)**
`writeSet` has been parsed since REQ-032 and rendered nowhere. Since the badge is truncated by design, the drawer is where the full answer belongs; the overlap row uses `makeTicketLinkList` so "what else writes these files?" is one click. No `template.html` change — drawer rows are JS-generated.

**D-08 — `prime-do-kanban.md`'s REQ-032 lesson was converted from a link into an inline entry. (DECIDE & STATE)**
The line the docs half had to rewrite was that file's only *linked* lesson, pointing into `do-work/archive/`. The file's own header comment mandates inlined lessons for exactly this reason: this repo's `do-work/` tree is git-excluded and export-ignored, so the link was dead in every other clone and every consumer install. Fixing it was the same edit, not an adjacent improvement.

**D-09 — `actions/board.md`'s "writes exactly three things" preamble left unchanged. (SILENT — verified, recorded because it was explicitly checked.)**
The annotation is computed in memory during board build and rendered; it adds no write surface. The read-only contract is still accurate as written.

**D-10 — `_dev/tests/contract-regressions.sh` was not touched. (DECIDE & STATE)**
It is outside the declared `write_set` and the Scope called it a conscious skip. The existing `fields["write_set"]` ratchet still pins the parse, and the new lock-step claim now has three Go tests behind it rather than a grep. Whether the "annotation runs after bucketing" invariant also deserves a shell ratchet is the orchestrator's call — filed under Discovered Tasks rather than silently dropped.

## Discovered Tasks

- **No contract-regression ratchet for the display-only overlap invariant.** `TestWriteSetOverlapNeverAffectsColumnPlacement` covers the behavior, but nothing pins the *instruction-side* claim the way the three REQ-032 ratchets pin the gate text. A cheap anchor exists: grep `tools/queue-kanban/model.go` for `annotateWriteSetOverlap` appearing after `bucketColumns` in `buildBoard`, and/or grep `actions/board.md` + `CLAUDE.md` for the display-only wording. Out of this REQ's `write_set` (D-10).
- **`board.js` badge rendering has no automated coverage.** There is no JS test runner; `generate_test.go` asserts substrings of the inlined `board.js` for other behaviors, so a `badge-write-overlap` / `writeSetOverlaps` substring assertion there would be the cheapest ratchet against a silent frontend regression. `generate_test.go` was outside the declared Scope.
- **No `docs/board-guide.md` exists.** Every other major action has a user guide; the board's features (badges included) are documented only inside `actions/board.md`, which is agent-facing instruction rather than user documentation. Pre-existing gap, noticed while looking for other places the overlap claim might live.

## Implementation Summary

**Files changed:**
- `tools/queue-kanban/model.go` (modified) — `WriteSetOverlaps []string` derived field; `annotateWriteSetOverlap` + intersection helpers (literal equality ∪ `filepath.Match` both directions; glob-vs-glob as literals, documented caveat); called from `buildBoard` **after** `bucketColumns` so "display only" is structural; pending/claimed only (`reserved` excluded, D-04)
- `tools/queue-kanban/generate.go` (modified) — `writeSetOverlaps` payload field + assignment
- `tools/queue-kanban/model_test.go` (modified) — three tests: pairing fixture table (literal/glob/disjoint/empty/terminal statuses/self), both-direction glob matching, and a `buildBoard` end-to-end asserting overlap never affects column placement
- `tools/queue-kanban/web/board.js` (modified) — `overlaps REQ-…` card badge rendering the derived array (no JS computation); drawer rows "Write set" + linked "Overlapping write sets"
- `tools/queue-kanban/web/board.css` (modified) — `.badge-write-overlap` (pending accent — contention is a scheduling hint, not a failure)
- `actions/capture-reference.md` (modified) — two-sentence slicing nudge inside the existing Slicing convention paragraph (per-concern file boundaries; declare unavoidable overlap in `write_set`)
- `actions/capture.md` (modified) — complex-batch field enumeration gains `write_set` (closes REQ-032 review Minor 5); `depends_on` conditional re-scoped
- `actions/board.md` (modified) — lock-step sentence rewritten (fixes REQ-032 review Minor 4; display-only overlap annotation described truthfully; co-dispatch stays with the pipeline gate) + badge-behavior paragraph
- `tools/queue-kanban/prime-do-kanban.md` (modified) — stale "no overlap computation" lesson rewritten to match the code; REQ-032 lesson inlined per the file's own header (dead link in consumer installs, D-08)
- `CLAUDE.md` (modified) — lock-step claim updated (parse + display-only overlap annotation; dispatch decisions stay in the pipeline)

**What was done:** Capture's slicing convention now nudges toward per-REQ file boundaries with unavoidable overlap declared in `write_set`, and the Kanban board surfaces write-set contention: a Go-side display-only annotation (`annotateWriteSetOverlap`, computed after column bucketing) feeds an `overlaps` badge and drawer rows on pending/claimed cards, glob-aware in both directions, with the three stale "no overlap computation" claims (board.md, prime, maintainer doc) updated in lock-step.

## Qualification

Passed — 10 files in `git diff --stat`, exact Scope match (`tools/checks/scope-drift.sh`: OK); all six acceptance criteria traced (slicing nudge inside the existing paragraph, overlap declared in write_set, badge naming contending REQ ids on pending/claimed only, display-only made structural by the after-bucketing call site + `TestWriteSetOverlapNeverAffectsColumnPlacement`, badge follows existing treatment with parser/payload lock-step, all three stale "no overlap computation" claims updated). Continuation-build note: the first builder died mid-run on a transient API error after finishing the Go/JS half; the continuation builder verified that half against the plan (unchanged), completed the five doc files, and reconciled the P-A-U/Decisions trail — `go build` was confirmed green before relaunch, and the final diff is a single coherent implementation. `tools/checks/qualify.sh`: OK. No debug artifacts (grep for console.log/debugger/fmt.Print/TODO clean).

## Testing

**Tests run:** `cd tools/queue-kanban && go test ./...`; `go vet ./...`; `gofmt -l`; `bash _dev/tests/contract-regressions.sh` (all re-run by the orchestrator, independent of the builder)
**Result:** ✓ All passing (go: ok incl. 3 new overlap tests; vet clean; gofmt clean; contract regressions passed)

**Red-green validation:**
- GREEN condition from Red-Green Proof: scratch queue with overlapping `write_set`s → `REQ-901 overlaps ['REQ-902']`, `REQ-902 overlaps ['REQ-901']` (glob vs literal, both directions), disjoint `REQ-903` gets none; badge markup + `writeSetOverlaps` confirmed inlined in the generated static artifact (builder transcript, scratch cleaned)
- Before: `grep writeSet tools/queue-kanban/web/board.js` had no hits and the Slicing convention said nothing about file boundaries → after: badge renderer + drawer rows present; nudge shipped

**New tests added:**
- `model_test.go`: `TestAnnotateWriteSetOverlapPairsContendingRequests`, `TestAnnotateWriteSetOverlapMatchesGlobsInBothDirections`, `TestWriteSetOverlapNeverAffectsColumnPlacement`

*Verified by work action*

## Review' section exactly per review-work.md's 'Append to REQ File' template (verdict, scores table, acceptance result, findings each tagged Critical/Important/Minor/Nit, requirements checklist, suggested testing). The orchestrator appends it verbatim — no preamble around it."},"uuid":"aa185d7d-7592-4203-aa53-17af9dc5b60a","timestamp":"2026-07-28T23:10:45.884Z","userType":"external","entrypoint":"cli","cwd":"/Users/t2/Desktop/e1-experimental-repos/skill-do-work2","sessionId":"2a6fd1c7-fc1c-4828-a98a-a644a710ffd3","version":"2.1.220","gitBranch":"main"}

**Orchestrator addendum:** Standard single-pass review (calibrated depth — Route B). Important finding → follow-up REQ-040 (glob dialect: `path.Match` + document `*`/`**`/malformed-pattern behavior). The review's second would-be follow-up (prime inline-lesson recurrence) is Minor and goes to the report plus REQ-041's consent list. Builder's three Discovered Tasks classified [normal]/[low], none meeting the test-hygiene carve-out → consolidated into REQ-041 (`pending-answers`, reviewed via `do-work clarify`).

## Lessons Learned

**What worked:** Continuation-build recovery — the first builder died mid-run on a transient API error, and because its half compiled and the REQ file carried the full trail (Scope, Exploration anchors), a fresh builder completed the remaining doc files without redoing anything, then reconciled the P-A-U/Decisions record to match reality. Compute-in-Go for the overlap annotation (D-01) bought three real tests where a frontend computation would have shipped untested; calling it after `bucketColumns` (D-02) made "display only" structural — a pattern worth reusing for any future display-only annotation.
**What didn't:** Choosing the glob primitive without a dialect audit — `filepath.Match` is OS-dependent (on Windows `*` crosses `/`), `**` silently under-matches, and malformed patterns silently degrade to literal matching. None of that was documented; REQ-040 fixes primitive + docs.
**Worth knowing:** Absent `write_set` reads as "no badge" on the board but "serialize" at the dispatch gate — same semantics (unknown ≠ safe), opposite rendering; the three places documenting this use clashing adverbs (review Minor) — pick one framing when next touching them. The prime `tools/queue-kanban/prime-do-kanban.md` mandates *inline* lessons (its `do-work/` links die in consumer installs), but the pipeline's Lessons-capture step still prescribes links — recurrence noted in REQ-041.

## Orientation

Capture now nudges slicing toward per-REQ file boundaries, and the board makes write-set contention visible: an `overlaps` badge plus drawer rows on pending/claimed cards, fed by a Go-side display-only annotation that runs after column bucketing (structurally unable to affect placement). Lives in capture guidance (`actions/capture-reference.md`/`actions/capture.md`) and the queue-kanban tool (model/generate/web). Completes UR-007's original three-REQ batch — the upstream lever (slicing) and the visibility lever (badge) around REQ-032's dispatch gate. One display-fidelity follow-up (REQ-040) queued.
