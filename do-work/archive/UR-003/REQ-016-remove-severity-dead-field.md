---
id: REQ-016
title: Remove the producer-less `severity` frontmatter field from queue-kanban
status: completed
created_at: 2026-07-01T21:06:45Z
claimed_at: 2026-07-01T21:28:40Z
completed_at: 2026-07-01T21:34:01Z
route: A
commit: 023aa50
user_request: UR-003
domain: backend
prime_files: [tools/queue-kanban/prime-do-kanban.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: false
related: [REQ-015]
kb_status: promoted
kb_entry: REQ-016-remove-the-producer-less-severity-frontm.md
---

# Remove the producer-less `severity` frontmatter field from queue-kanban

## What

Remove the `severity` frontmatter pipeline from the queue-kanban tool. No REQ schema in this repo ever emits a top-level `severity:` key (the Schema Read Contract in `actions/work-reference.md` doesn't define one; discovered-task severity lives as an inline `[critical]`/`[normal]`/`[low]` bullet prefix inside `## Discovered Tasks`, never as frontmatter), yet the tool carries a full parse → JSON → badge pipeline for it.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Sweep confirmed exactly the four known-site files (`model.go`, `generate.go`, `web/board.js`, `web/board.css`) — no test/testdata files reference `severity`. Plan: delete `RequestTicket.Severity` field + its `coerceScalarToString(fields["severity"])` parse line in `model.go`; delete the `Severity` JSON field + its assignment in `generate.go`'s `generatedRequest`/`buildGeneratedBoardData`; delete the card-badge `if (request.severity)` block and the drawer `appendMetaRow("Severity", ...)` block in `web/board.js`; delete the `.badge-severity` + `.badge-severity[data-severity=...]` rules in `web/board.css`. `makeBadge` is a shared generic helper (used by domain/ur/route badges too) — leave it untouched, only remove the severity call site. Pure deletion, no ripple into tests since none reference it.
- [x] **[APPLY]:** Deleted `Severity` field from `model.go`'s `RequestTicket` struct and its `coerceScalarToString(fields["severity"])` parse assignment; deleted `Severity` JSON field from `generate.go`'s `generatedRequest` struct and its assignment in `buildGeneratedBoardData`; deleted the `if (request.severity)` badge block in `web/board.js`'s card renderer and the `if (request.severity) appendMetaRow(...)` block in the drawer renderer; deleted the `.badge-severity` and `.badge-severity[data-severity=...]` rule blocks from `web/board.css`. No test/testdata files touched — sweep found none referencing severity. Post-edit sweep (`grep -rni severity .` excluding the gitignored binary) returns nothing.
- [x] **[UNIFY]:** `git diff --stat` confirms exactly the four planned files changed, pure deletion (4 files changed, 2 insertions(+), 23 deletions(-) — the 2 insertions are the reformatted `Route`/`Batch` struct field lines in `model.go` after `Severity` was removed). Reviewed each file: `model.go` — struct field + parse assignment removed, nothing else touched; `generate.go` — JSON tag field + assignment removed; `web/board.js` — both the card-badge `if (request.severity)` block and the drawer `appendMetaRow("Severity", ...)` block removed, `makeBadge` helper (shared by domain/ur/route badges) left untouched; `web/board.css` — `.badge-severity` + `.badge-severity[data-severity=...]` rule blocks removed, no orphaned selectors. `gofmt -l .` → clean (no output). `go vet ./...` → clean (no output). No debug artifacts (no `fmt.Println`/`console.log`/commented-out code introduced). `git status --porcelain` shows only the 4 expected files modified, nothing untracked or extraneous.

## Why (if provided)

Dead vertical feature: the badge can never render because the field is never populated. Root `CLAUDE.md`'s lock-step rule means the parser shouldn't read fields the Schema Read Contract doesn't define — if severity frontmatter is ever wanted, it gets added to the contract first, then mirrored here.

## Context

- Surfaced as an out-of-scope finding during the 2026-07-01 `do-work validate-feedback` triage (see UR-003). Confirmed against HEAD (post-0.102.1).
- Known sites (sweep for others): `model.go` ~line 60 (struct field) and ~line 323 (parse from `fields["severity"]`); `generate.go` ~line 77 (JSON tag) and ~line 181 (copy); `web/board.js` ~lines 121–123 (severity badge) and ~366–367 (drawer meta row); `web/board.css` ~lines 517–521 (`.badge-severity` styles).
- A consumer repo could in theory hand-edit `severity:` into a REQ, but the skill's stance is schema-first: undeclared fields shouldn't get parser support. Note the removal in the run report so it's visible.

## Builder Guidance

Firm — full-vertical removal **user-confirmed** during the verify-requests pass (2026-07-01). Remove the whole vertical (Go struct/parse, JSON export, JS badge + drawer row, CSS) — don't leave half the pipeline. Sweep `tools/queue-kanban/` for any `severity` reference beyond the known sites (tests included) before calling it done.

## Red-Green Proof
**RED prompt/case:** `grep -rni severity tools/queue-kanban/` returns the four-file pipeline while `grep -rn '^severity:' do-work/ actions/ specs/` (and the Schema Read Contract's field table) show no producer or schema definition.
**Why RED now:** The tool ships parse/export/render support for a field that cannot occur under the documented schema.
**GREEN when:** `grep -rni severity tools/queue-kanban/` returns nothing; `go build` and `go test ./...` in `tools/queue-kanban/` pass; the board renders tickets without errors.
**Validation:** User confirmed (verify-requests pass, 2026-07-01 — full-vertical removal confirmed over Go-only or keep)

---
*Source: UR-003 — "capture the two out-of-scope kanban findings as REQs" (finding 2, restated in the UR Summary)*

Think carefully before answering.

---

## Triage

**Route: A** - Simple

**Reasoning:** The REQ enumerates every known site across the four-file vertical (`model.go` struct+parse, `generate.go` JSON tag+copy, `web/board.js` badge+drawer row, `web/board.css` styles), the direction is user-confirmed (full-vertical removal), and the change is pure deletion of a dead feature with a grep-based sweep as the only discovery step. Well-specified with obvious scope.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Implementation Summary

**Files changed:**
- `tools/queue-kanban/model.go` (modified) — removed the `Severity` field from `RequestTicket` and its `coerceScalarToString(fields["severity"])` parse assignment
- `tools/queue-kanban/generate.go` (modified) — removed the `Severity` JSON-tagged field from `generatedRequest` and its copy in `buildGeneratedBoardData`
- `tools/queue-kanban/web/board.js` (modified) — removed the card-badge `if (request.severity)` block and the drawer `appendMetaRow("Severity", ...)` block
- `tools/queue-kanban/web/board.css` (modified) — removed the `.badge-severity` and `.badge-severity[data-severity="high"/"critical"]` rule blocks

**What was done:** Removed the entire producer-less `severity` frontmatter vertical (Go struct/parse → JSON export → JS badge/drawer render → CSS) from the queue-kanban tool. Sweep confirmed no other sites existed (tests included); the shared `makeBadge` helper stays, since domain/ur/route badges use it.

## Qualification

Passed — 4 files verified on disk (`git diff` shows a pure deletion: 2 insertions / 23 deletions, the insertions being reformatted struct lines), the REQ's single requirement (full-vertical removal, no half-pipeline) traced across all four layers, P-A-U boxes checked with substantive notes, no debug artifacts. Orchestrator independently re-ran the grep (0 hits), build, tests, linters, and a live `summary` render against this repo's real queue (16 tickets, exit 0).

## Testing

**Tests run:** `go test -count=1 ./...` in `tools/queue-kanban/` (orchestrator re-ran independently); `go build -o /dev/null .`; `go run . summary --repo-root <repo>`; `go run . generate --out <scratchpad>` render check
**Result:** ✓ All passing (`ok github.com/knews2019/skill-do-work/queue-kanban`); board renders 16 tickets / 3 URs without errors; generated `board-data.js` contains no `"severity"` JSON key

Red-green validation omitted — non-behavioral removal of a dead vertical (no request-specific tests exist for a field nothing produces); regression evidence used instead, tracing to the REQ's `## Red-Green Proof`: RED grep (four-file pipeline present) confirmed before, GREEN criteria (`grep -rni severity tools/queue-kanban/` empty, build + tests pass, board renders) all met after.

**Linters:** `gofmt -l .` clean; `go vet ./...` clean.

## Review

**Approve** — clean full-vertical deletion of the dead `severity` field, verified end to end.
Route A | reviewed pre-commit

### What's built
- The producer-less `severity` frontmatter field is fully removed across all four layers: Go struct/parse (`model.go`), JSON export (`generate.go`), JS badge + drawer row (`web/board.js`), CSS rules (`web/board.css`). The `makeBadge` shared helper (used by domain/ur/route badges) is untouched.

### Decisions / risks for you
- None. Pure deletion of unreachable code; no behavior change for any real REQ since no schema ever emitted `severity:`.

### Findings

**Important:** None. **Minor:** None. **Nit:** None — diff is minimal and surgical (2 insertions / 23 deletions, insertions being gofmt-realigned adjacent struct lines).

### Requirements Checklist

- [x] Remove `RequestTicket.Severity` field + `coerceScalarToString(fields["severity"])` parse line in `model.go` — delivered
- [x] Remove `Severity` JSON field + assignment in `generate.go` — delivered
- [x] Remove card-badge `if (request.severity)` block in `web/board.js` — delivered
- [x] Remove drawer `appendMetaRow("Severity", ...)` block in `web/board.js` — delivered
- [x] Remove `.badge-severity` + `.badge-severity[data-severity=...]` CSS rules — delivered
- [x] Leave shared `makeBadge` helper untouched — delivered
- [x] Sweep `tools/queue-kanban/` for any other `severity` reference (tests included) — delivered, confirmed independently (grep returns nothing)
- [x] No half-pipeline left behind (full-vertical removal, user-confirmed) — delivered

### Acceptance Testing

**Result: Pass** — `gofmt`/`go vet`/`go build`/`go test -count=1` all clean; grep sweep empty (tracked files, gitignored binary excluded); live `summary` + `generate` renders against this repo's real `do-work/` tree succeed with zero `"severity"` keys in `board-data.js`; no orphaned `data-severity`/`.badge-severity` references; Schema Read Contract field table confirmed to define no `severity:` field (the REQ's premise holds).

### Suggested Additional Testing

- None — pure dead-code deletion with no new runtime surface; render + grep + go test/vet/build cover the observable behavior.

### Scores

**Overall: 100%** — Requirements 100 / Code Quality 100 / Test Adequacy N/A (dead-code removal; regression evidence confirmed independently) / Scope 100 / Risk: none / Acceptance: Pass.

### Follow-up REQs Created
None.

*Generated by review-work agent (pipeline mode)*

## Lessons Learned

**What worked:** Enumerating the full vertical (parse → export → render → style) in the REQ up front made the deletion mechanical; verifying the render against the repo's real `do-work/` tree (not just unit tests) proved the board still works end to end.
**What didn't:** Nothing — no dead ends.
**Worth knowing:** The neighboring `batch` frontmatter field looks similar but is NOT a dead vertical — it has real producers in archived REQs (UR-002's REQ-013/REQ-014 frontmatter), so don't sweep it up in a future "same shape" cleanup without re-checking. When grepping the generated `board-data.js` for leftover fields, match the JSON key (`"severity":`) — the bare word appears legitimately in rendered REQ body prose.

## Orientation

The queue-kanban board no longer carries a parse→export→render pipeline for a `severity` frontmatter field that nothing in the skill ever produces — the parser now reads only fields the Schema Read Contract defines. Lives in the board's ticket model + frontend (`tools/queue-kanban/`, per prime-do-kanban's Stakes). No map change; prime spot-check clean.
