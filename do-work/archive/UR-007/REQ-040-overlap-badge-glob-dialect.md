---
id: REQ-040
title: "Board overlap badge: use path.Match and document the glob dialect"
status: completed
claimed_at: 2026-07-29T08:54:22Z
completed_at: 2026-07-29T09:00:22Z
commit: acc4722
route: B
kb_status: promoted
kb_entry: REQ-040-board-overlap-badge-use-path-match-and-d.md
created_at: 2026-07-28T23:18:24Z
user_request: UR-007
addendum_to: REQ-034
depends_on: []
related: [REQ-034]
batch: parallel-dispatch
domain: general
prime_files: []
tdd: false
review_generated: true
write_set:
  - tools/queue-kanban/model.go
  - tools/queue-kanban/model_test.go
  - actions/board.md
  - tools/queue-kanban/prime-do-kanban.md
  - actions/work-reference.md
maintenance: false
---

# Board overlap badge: use path.Match and document the glob dialect

## What

REQ-034's `writeSetPatternsIntersect` uses `filepath.Match`, which is OS-dependent (`*` would cross `/` on Windows because the separator is `\`); `write_set` entries are slash-separated repo-relative paths, so `path.Match` is the correct drop-in primitive. Additionally the glob dialect is documented nowhere: `**` silently under-matches (`filepath.Match("src/**/x.ts", "src/a/b/x.ts")` is false) and a malformed pattern (`ErrBadPattern`) silently degrades to literal-only matching — while the pipeline gate treats an unexpandable glob as overlapping. Fix the primitive and state the dialect where readers meet the field.

## Why (if provided)

REQ-034's review (Important): three classes of silent false negative (glob-vs-glob, `**`, malformed patterns) with only the first documented — the badge quietly missing real contention is what erodes trust in it.

## Detailed Requirements

- Swap `filepath.Match` → `path.Match` in `writeSetPatternsIntersect`; add a test pinning slash semantics (e.g. `web/*` must NOT match `web/a/b.css`) and one for malformed-pattern behavior.
- Document the dialect in: the `writeSetPatternsIntersect` doc comment, `actions/board.md`'s badge paragraph, `tools/queue-kanban/prime-do-kanban.md`'s lesson, and one clause on `actions/work-reference.md`'s `write_set` schema line ("globs use path.Match semantics: `*` never crosses `/`, `**` is not recursive, malformed patterns match nothing on the board — the gate still treats unexpandable globs as overlapping").
- Keep all existing overlap tests green; parser lock-step unchanged (no schema shape change).

## Constraints

- Small, surgical — no new subsections; extend existing sentences/comments.

---

## Triage

**Route: B** - Medium

**Reasoning:** A surgical, well-specified change (swap one primitive, add two tests, document a dialect in four places), but it touches real Go code + tests and needs the Go toolchain, so exploration confirms the exact anchors and the import situation before building.

**Planning:** Not required (Route B)

## Exploration

Anchors verified against disk (post-REQ-039, commit `9c8ef45`):
- `tools/queue-kanban/model.go:1142` `writeSetPatternsIntersect`; the two glob calls are `filepath.Match` at `:1146` and `:1149`; the doc comment at `:1130`–`:1132` names `filepath.Match`.
- **`path/filepath` is imported at `:8` and is used ELSEWHERE in model.go** — `filepath.Base`/`filepath.ToSlash` (`:315`/`:318`) and `filepath.Base`/`filepath.Dir` (`:1248`/`:1259`). **So the import cannot be removed.** Add a `"path"` import alongside `"path/filepath"` and use `path.Match` *only* inside `writeSetPatternsIntersect`; leave the other `filepath.*` calls untouched.
- Baseline green: `go test ./...` ok, `go vet` clean, `gofmt -l` clean (Go 1.26.1).
- `write_set` schema line lives at `actions/work-reference.md:100`; board badge paragraph in `actions/board.md`; the lesson in `tools/queue-kanban/prime-do-kanban.md`.

*Exploration folded into orchestrator scoping (Route B)*

## Scope

**Files I will touch:**
- `tools/queue-kanban/model.go` (modify) — add `"path"` import; swap the two `filepath.Match` calls in `writeSetPatternsIntersect` to `path.Match`; update the doc comment to name `path.Match` and state the dialect (`*` never crosses `/`, `**` not recursive, malformed ⇒ `ErrBadPattern` ⇒ no match). Leave the other `filepath.*` uses.
- `tools/queue-kanban/model_test.go` (modify) — add a slash-semantics test (`web/*` must NOT match `web/a/b.css`) and a malformed-pattern test; keep all existing overlap tests green.
- `actions/board.md` (modify) — the badge paragraph: document the glob dialect.
- `tools/queue-kanban/prime-do-kanban.md` (modify) — the lesson: document the dialect.
- `actions/work-reference.md` (modify) — one clause on the `write_set` schema line (`:100`): "globs use `path.Match` semantics: `*` never crosses `/`, `**` is not recursive, malformed patterns match nothing on the board — the gate still treats unexpandable globs as overlapping."

**Files I will NOT touch:** the schema field shape (no parser lock-step needed — this is display-only annotation logic, not a new/changed frontmatter field); `_dev/tests/contract-regressions.sh` (the `fields["write_set"]` parser ratchet stays green — the field isn't changing shape).

**Acceptance criteria (restated from REQ):**
- [ ] `filepath.Match` → `path.Match` in `writeSetPatternsIntersect` (add `"path"` import; keep `filepath` for its other uses)
- [ ] A test pins slash semantics (`web/*` does NOT match `web/a/b.css`) and one covers malformed-pattern behavior
- [ ] The dialect is documented in all four places: the `writeSetPatternsIntersect` doc comment, `actions/board.md`'s badge paragraph, `prime-do-kanban.md`'s lesson, and the `actions/work-reference.md` `write_set` schema line
- [ ] All existing overlap tests stay green; no schema shape change (parser lock-step unaffected)
- [ ] `go test ./...`, `go vet ./...`, `gofmt -l` all clean

*Scope declared by work action (orchestrator, session do-work-20260729T065754Z-5724)*

## AI Execution State (P-A-U Loop)

- [x] PLAN — success criteria transformed to verifiable goals; anchors + import situation confirmed against disk
- [x] ACT — swap `filepath.Match` → `path.Match` in `writeSetPatternsIntersect` (add `"path"` import, keep `"path/filepath"`); document the dialect in the doc comment, `actions/board.md`, `prime-do-kanban.md`, and `actions/work-reference.md`; add the two new tests
- [x] UNIFY — `git diff --stat` reviewed (5 files, +58/-11); `go test`/`go vet` clean, `gofmt -l` empty; `contract-regressions.sh` exit 0; no debug artifacts

## Implementation Summary

**Files changed:**
- `tools/queue-kanban/model.go` (modified)
- `tools/queue-kanban/model_test.go` (modified)
- `actions/board.md` (modified)
- `tools/queue-kanban/prime-do-kanban.md` (modified)
- `actions/work-reference.md` (modified)

**What was done:** Swapped `writeSetPatternsIntersect`'s two `filepath.Match` calls to `path.Match` (added a `"path"` import alongside the retained `"path/filepath"`, which is still used by `filepath.Base`/`ToSlash`/`Dir` elsewhere in the file) so the board's overlap annotation is OS-independent — `write_set` entries are slash-separated repo-relative paths, and `filepath.Match`'s `\` separator on Windows would let `*` wrongly cross `/`. Rewrote the doc comment to name `path.Match`, give the rationale, and state the dialect (`*` never crosses `/`; `**` is not recursive — no superglob; malformed ⇒ `ErrBadPattern` ⇒ no-match), plus a display-only note that the pipeline's dispatch gate still serializes an unexpandable/overlapping glob, so a board false-negative never loosens the gate. Added two tests (`…GlobStarNeverCrossesSlash` with a same-segment positive control, and `…MalformedPatternMatchesNothing`) in the existing fixture style. Documented the same dialect in `actions/board.md`'s badge paragraph, `tools/queue-kanban/prime-do-kanban.md`'s lesson, and one clause on `actions/work-reference.md`'s `write_set` schema line. No frontmatter/schema shape change, so the board parser lock-step is unaffected.

*Summary written by work action (orchestrator)*

## Qualification

**Passed.** `qualify.sh` OK (5 files present + in diff, no debug artifacts); `scope-drift.sh` clean (5 declared = 5 touched). Orchestrator independent re-run: `go vet ./...` clean, `gofmt -l .` empty, `go test -count=1` (no cache) green including the 2 new tests, `contract-regressions.sh` exit 0. Judgment: the Go change is substantive (real primitive swap + guard preserved); the import is handled correctly (`"path"` added, `"path/filepath"` retained and still used by `.Base`/`.ToSlash`/`.Dir`); all 5 acceptance criteria trace to diff changes; check 6 (data flows) N/A — this is a pure/deterministic matcher, not a data path.

*Verified by work action (orchestrator)*

## Testing

**Tests run:** `cd tools/queue-kanban && go test -count=1 ./...` (fresh, no cache); `go vet ./...`; `gofmt -l .`; `bash _dev/tests/contract-regressions.sh` (repo root).
**Result:** ✓ All green. New tests pass fresh: `TestAnnotateWriteSetOverlapGlobStarNeverCrossesSlash`, `TestAnnotateWriteSetOverlapMalformedPatternMatchesNothing`. Existing overlap tests unchanged and green. vet/gofmt clean; contract-regressions exit 0 (the `fields["write_set"]` parser ratchet stays green — no schema shape change).

**Red-green validation:** The two new tests encode behavior that the *old* `filepath.Match` code got wrong on Windows (`*` crossing `/`) — under `path.Match` they pass and pin the corrected slash semantics + malformed-pattern handling. Existing tests provide the non-regression evidence for the unchanged literal/single-segment cases.

*Verified by work action*

## Review

**Overall: 95%** (Route B — single-reviewer pass per the session's calibration) | 2026-07-29T09:00:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Findings:** 0 important, 0 minor. All requirements delivered; the fix is a clean one-primitive swap with the build-breaking import trap avoided; the slash-semantics test includes a positive control (pins the boundary, not mere absence); the dialect is documented in all four required places with the display-only/gate-still-serializes framing consistent across code comment, board doc, prime, and schema line.
**Acceptance:** Pass — fresh Go suite + contract-regressions green.
**Follow-ups created:** None.

*Reviewed by work action (orchestrator, single-reviewer pass)*

## Lessons Learned

**What worked:** Keeping `path/filepath` (still used by `.Base`/`.ToSlash`/`.Dir`) and adding `"path"` alongside avoided the obvious build-break of "swap the import." The positive-control assertion in the slash test (`web/*` DOES match `web/app.css`) makes the test prove the boundary is the slash, not the glob silently failing.

**What didn't:** Nothing surprising — the anchors and the import situation were confirmed before the build, so it landed first try.

**Worth knowing:** `path.Match` has no recursive `**` superglob and returns `ErrBadPattern` on malformed input — so the board's overlap badge is a best-effort *visualizer*, not the dispatch gate. The safety-relevant reader is the work pipeline's gate, which treats an unexpandable/overlapping glob as overlapping (serialize); a board false-negative never loosens it. `write_set` entries must stay slash-separated repo-relative paths for `path.Match` to be correct.

## Orientation

The board's `write_set` overlap badge now matches globs with OS-independent `path.Match` (was `filepath.Match`, which crossed `/` on Windows), and the glob dialect is documented where readers meet the field — the `writeSetPatternsIntersect` doc comment, `actions/board.md`, `tools/queue-kanban/prime-do-kanban.md`, and the `actions/work-reference.md` `write_set` schema line. Display-only annotation fix in the board tool (`tools/queue-kanban/model.go` + test); no schema/parser shape change, no map impact. Closes REQ-034's review finding (three classes of silent false-negative, only one documented).
