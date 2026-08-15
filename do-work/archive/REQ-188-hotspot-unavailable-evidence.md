---
id: REQ-188
title: Hotspot output silently drops unavailable tracked paths
status: completed
created_at: 2026-08-15T07:13:20Z
claimed_at: 2026-08-15T12:14:36Z
completed_at: 2026-08-15T12:33:49Z
user_request: UR-041
domain: backend
prime_files: []
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
effort_estimate: normal
related: [REQ-181, REQ-182, REQ-183, REQ-184, REQ-185, REQ-186, REQ-187]
batch: audit-findings-2026-08-14
write_set: [skills/do-work-toolbox/tools/audit-metrics/churn.go, skills/do-work-toolbox/tools/audit-metrics/churn_test.go, skills/do-work-toolbox/tools/audit-metrics/main.go]
route: B
kb_status: pending
kb_entry:
---

# Hotspot Output Silently Drops Unavailable Tracked Paths

## What

Keep unreadable or otherwise unavailable tracked paths visible in hotspot output as `NOT-MEASURED`, while preserving valid measured rows and warning that the ranking is incomplete.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Replace the lossy hotspot slice with one focused result carrying measured entries plus sorted unavailable paths, preserve numeric sorting/top-count behavior, render every unavailable churn-bearing path in a compact `NOT-MEASURED` section with an incomplete-ranking warning, and lock the full production sequence with a real-Git missing-worktree fixture.
- [x] **[APPLY]:** Added the focused hotspot result, carried measurement failures as sorted unavailable paths, rendered a visible incomplete-ranking warning plus uncapped `NOT-MEASURED` evidence, updated the sole caller, and locked the production compute/render path with a real-Git fixture. Changes are limited to the three scoped Go files.
- [x] **[UNIFY]:** Reviewed the full diff and caller search. Verified `churn.go` preserves measured score/path ordering and numeric `topCount` while independently sorting all unavailable paths; `churn_test.go` proves a valid numeric row plus two sorted missing rows under `topCount=1`; `main.go` passes the focused result through the sole production call. Focused and full uncached tests, `go vet ./...`, `gofmt -l`, `git diff --check`, and the canonical maintainer gate pass with no debug artifacts.

## Why

`computeHotspotEntries` currently continues after a tracked-path measurement error and emits a successful numeric ranking without the path or a `NOT-MEASURED` marker. The omission makes unavailable evidence look like low risk and contradicts REQ-178's visible-warning decision.

## Context

- Audit priority: P3; impact 2; effort normal.
- Root-cause key: `hotspot-unavailable-evidence-visible`.
- Evidence source: `do-work/audits/audit-2026-08-14.md`, Finding 8.
- Reproduce: `cd skills/do-work-toolbox/tools/audit-metrics && probe_dir="$(mktemp -d /tmp/audit-metrics-unreadable.XXXXXX)" && git -C "$probe_dir" init -q --initial-branch=main && git -C "$probe_dir" config user.email probe@example.test && git -C "$probe_dir" config user.name Probe && ln -s missing-target "$probe_dir/unreadable.md" && git -C "$probe_dir" add unreadable.md && git -C "$probe_dir" commit -qm seed && go run . inventory --repo-root "$probe_dir" --top-count 20 && go run . hotspots --repo-root "$probe_dir" --since-window '10 years' --top-count 20`

## Detailed Requirements

- Carry a sorted collection of unavailable tracked paths through the churn/hotspot join instead of discarding measurement errors.
- Preserve every valid measured hotspot row.
- Render each unavailable path with current lines and score shown as `NOT-MEASURED`.
- Emit a visible warning that the numeric ranking is incomplete when unavailable paths exist.
- Add a real-Git fixture containing a tracked path missing from the worktree; prove valid rows survive while the unavailable path remains visible.

## Constraints

- Do not fail the entire hotspot command solely because one tracked path cannot be measured.
- Do not add a generic diagnostics subsystem.
- Lock-in limit: zero churn-bearing tracked paths silently absent from both measured and `NOT-MEASURED` output.

## Dependencies

None.

## Builder Guidance

Firm intent. One unavailable-path result field, one compact renderer section, and one real-Git test are the earned surface.

## Open Questions

None.

## Red-Green Proof
**RED prompt/case:** In a temporary Git repository, commit a symlink to a missing target, then run `audit-metrics inventory` and `audit-metrics hotspots`.
**Why RED now:** Inventory exposes the unreadable tracked path, while hotspot output succeeds and silently omits it.
**GREEN when:** Valid hotspot rows still render, the missing-worktree path appears in a sorted `NOT-MEASURED` section with lines/score unavailable, and the output warns that ranking is incomplete.
**Validation:** Confirmed by the user during verification on 2026-08-15.

## Assets

`do-work/user-requests/UR-041/assets/REQ-181-screenshot-1-validated-audit-findings.png`

The screenshot shows this request as row 08, labeled P3, impact 2, normal effort.

## Full Context

See `do-work/user-requests/UR-041/input.md` and Finding 8 in the canonical audit for the complete batch constraints and validated evidence record.

---
*Source: "do-work capture-request for these" — expanded from attached validated audit evidence.*

## Triage

**Route B** — the intended one-result-field repair is firm, but the caller signature, renderer/top-count contract, and portable real-Git fixture needed a focused code trace before Scope could be exact.

## Exploration

- `computeHotspotEntries` is the only lossy seam: every `measureTrackedFile` error currently executes `continue`, and only `runHotspotsCommand` calls it.
- A real-Git replay with a valid file plus a tracked broken symlink confirms the asymmetric output: inventory warns about the unreadable file, while hotspots prints the valid numeric row and silently omits the unavailable path.
- A `hotspotResult` with `MeasuredEntries` and `UnavailablePaths` is the minimal shape. Measured entries retain score-descending/path-ascending order; unavailable paths sort independently by path.
- Numeric `topCount` applies only to measured rankings. Unavailable rows must never consume a rank slot or be capped because every missing churn-bearing path is evidence.
- A portable fixture can commit regular files and remove two from the worktree without staging. That exercises the same `git ls-files` plus file-read boundary while also proving unavailable-path sorting.
- The renderer can still show known commit counts from `TouchesByPath`; only current lines and the derived score become `NOT-MEASURED`. OS error strings and a general diagnostics layer are unnecessary.

## Scope

**Files I will touch:**
- `skills/do-work-toolbox/tools/audit-metrics/churn.go` (modify) — carry sorted unavailable paths beside measured hotspots and render the warning/`NOT-MEASURED` section
- `skills/do-work-toolbox/tools/audit-metrics/churn_test.go` (modify) — add the real-Git missing-worktree RED/GREEN fixture and assertions
- `skills/do-work-toolbox/tools/audit-metrics/main.go` (modify) — pass the focused hotspot result from compute to renderer

**Files I will NOT touch:** Git history/rename normalization, inventory behavior, generic diagnostics, other modules, or the audit-metrics prime.

**Acceptance criteria (restated from REQ):**
- [x] Every churn-bearing tracked path is present either as a measured hotspot or as `NOT-MEASURED`.
- [x] Valid numeric rows, sorting, scores, and numeric `topCount` behavior remain intact.
- [x] Unavailable paths sort deterministically, render known commits with lines/score unavailable, and cannot be capped by `topCount`.
- [x] A visible warning states that the numeric ranking is incomplete whenever unavailable paths exist.
- [x] One unavailable path does not fail the whole command or introduce a generic diagnostics subsystem.

## Implementation Summary

**Files changed:**
- `skills/do-work-toolbox/tools/audit-metrics/churn.go` (modified) — separates measured hotspots from sorted unavailable churn-bearing paths and renders the warning plus uncapped evidence table
- `skills/do-work-toolbox/tools/audit-metrics/churn_test.go` (modified) — exercises the production churn/join/render sequence against one measurable and two missing tracked files
- `skills/do-work-toolbox/tools/audit-metrics/main.go` (modified) — carries the focused result through the sole hotspots command caller

**Behavior:** Hotspot output no longer silently drops churn-bearing tracked paths whose current worktree contents cannot be measured. Numeric rows retain their existing score ordering and `topCount`; unavailable rows render afterward in deterministic path order with known commits and `NOT-MEASURED` lines/score, alongside a visible incomplete-ranking warning.

## Testing

**RED:** The new real-Git fixture initially exited 1: the valid numeric row survived, but the warning and both removed tracked paths were absent from hotspot output.

**GREEN:**
- `go test -count=1 -run '^TestHotspotsReportKeepsUnavailableTrackedPathsVisible$' .` — PASS
- `go test -count=1 ./...` — PASS
- `go vet ./...` — PASS
- `gofmt -l churn.go churn_test.go main.go` — PASS (no output)
- `git diff --check` — PASS
- `bash _dev/tests/maintainer-verify.sh` — PASS

## Qualification

- **Scope:** PASS — `scope-drift.sh` reports the three-file Implementation Summary exactly matches Scope; foreign queue edits remain excluded.
- **Mechanical checks:** PASS — `qualify.sh` found complete P-A-U evidence, the exact modified file set, and no debug artifacts.
- **Substance and traceability:** PASS — the result type closes the single lossy measurement seam, and the real-Git test covers the lock-in limit through the production compute/render sequence.
- **Wiring/data flow:** PASS — `runHotspotsCommand` is the sole caller; measured and unavailable paths flow independently into one renderer, so numeric capping cannot hide unavailable evidence.

## Review

**Result:** Approve — Acceptance: Pass
**Overall score:** 96%

- **Requirements (100%):** Every churn-bearing path now partitions into measured or unavailable evidence; unavailable rows retain known commits, are sorted and uncapped, and trigger the incomplete-ranking warning without making the command fail.
- **Code quality (98%):** The focused result shape closes the one lossy seam without changing churn calculation, Git normalization, inventory, or exposing OS diagnostics.
- **Test adequacy (84%):** The real-Git production-path fixture mutation-locks unavailable-path carry, sorting, warning, known commits, and independence from `topCount`.
- **Scope (100%):** Exactly the three declared Go files changed; scope-drift passes.

**Important findings:** None.
**Minor findings:** The fixture has only one measured file, so it does not independently mutation-lock score-descending/path-ascending ordering or suppression of a second numeric row beyond `topCount`, even though those pre-existing code paths remain unchanged and an independent CLI fixture confirmed them.
**Explicit remediation:** None required for acceptance; a future test touch may add measured rows with distinct/tied scores and assert both ordering and numeric truncation alongside the uncapped unavailable evidence.

## Lessons Learned

**What worked:**
- Partitioning measured rankings from unavailable evidence keeps arithmetic honest while making completeness visible.
- Removing regular tracked files after commit is a portable real-Git fixture for the same read boundary as a broken symlink.

**What didn't:**
- Treating a per-path measurement failure as a harmless `continue` produced a plausible but incomplete report.
- A one-measured-row fixture cannot mutation-lock numeric ordering or capping even when it proves unavailable rows bypass the cap.

**Worth knowing:** `topCount` applies only to numeric hotspot entries. Every churn-bearing unavailable path remains evidence, keeps its known commit count, and renders uncapped in deterministic path order.

**Knowledge handoff:** Pending human triage. No knowledge-base file was written automatically.

## Orientation

**[MAP CHANGED]** Hotspot reporting now has an explicit completeness channel: numeric churn × size rankings stay capped, while every churn-bearing path unavailable in the current worktree remains visible as uncapped `NOT-MEASURED` evidence.
