---
id: REQ-559
title: '[impact-rule-change] Retry a red repository gate once before deferring or minting a repair REQ'
status: claimed
priority: now
created_at: 2026-09-03T20:05:46Z
user_request: UR-106
domain: general
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: false
impact: impact-rule-change
effort_estimate: effort-substantive
related: [REQ-548, REQ-531, REQ-560]
batch: lifecycle-overhead
write_set:
  - skills/do-work/actions/work.md
  - skills/do-work/actions/work-reference.md
  - skills/do-work/tools/checks/preflight.sh
claimed_at: 2026-09-04T18:15:54Z
route: B
estimate:
  p50_active_minutes: 25
  confidence: medium
  calculated_at: 2026-09-04T18:17:56Z
  basis:
    - Route B
    - 3-file write set
    - 2 subsystems involved
    - 5 acceptance criteria
    - async lifecycle behavior
dispatch_at: 2026-09-04T18:23:51Z
builder_handback_at: 2026-09-04T18:30:54Z
integration_at: 2026-09-04T18:30:54Z
---

# Retry a Red Repository Gate Once Before Deferring or Minting a Repair REQ

## What

When the repository gate exits non-zero, at the baseline (Step 5 pre-flight) or after integration (Step 6.5), rerun the exact same argv once, immediately, before any classification. A green rerun is recorded as green and the run continues; the retry is written to the run's progress output as one line. A red rerun enters the existing path unchanged: fingerprint, diagnostic worktree, `defer-gate`, repair REQ. The retry happens in the program that launches the gate where one exists (`preflight.sh` for the baseline), and as one rule in `work-reference.md` cited from Step 6.5 for the direct post-merge run.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read both primes plus the crew rules; located the real gate launch site (`internal/corehelpers/checks.go` `handlePreflight`, which `preflight.sh` only execs into) and the two prose lanes that treat a first non-zero exit as final.
- [x] **[APPLY]:** Retry mechanics in `checks.go` with two behavior probes in `checks_test.go`; one condition-keyed prose rule in `work-reference.md` cited from `work.md` Step 5, Step 5.75 and Step 6.5 item 4.
- [x] **[UNIFY]:** `git diff --stat` reviewed for all four files; `gofmt -l` empty, `go vet` clean, `go build ./...` clean, `action-shell-blocks.sh` and `contract-regressions.sh` pass; `git status --porcelain -uall` showed exactly the four modified files, no debug artifacts.

## Why

"that is the lifecycle around it that we need to improve a lot": REQ-548 was an already-green repair. The gate passed when rerun, and the only thing the 28-minute detour proved was that the first run had flaked.

## Context

- 2026-09-03: REQ-531 claimed 21:19; its baseline gate failed with `update-script-behavior.sh: printf write error: Broken pipe` (REQ-548's recorded diagnostic evidence); deferred 21:37; REQ-548 claimed 21:37, found the gate green, archived 21:47 as an already-green no-op; REQ-531 re-claimed 21:47. Two full gate runs and one REQ lifecycle to learn that nothing was broken.
- The already-green repair path (0.266.8) made that no-op cheaper, one gate run instead of three, but it still costs a claim, a checkpoint entry, a finalization and an archive. A retry before classification removes the whole detour for every transient failure and changes nothing for a real one: a red retry carries the same fingerprint into the same defer path.
- Exactly one retry. A second failure is a real failure by definition here; bounding it keeps a genuinely red gate from doubling its cost.
- The retry counts as a gate run for the one-gate-per-machine budget rule in `_dev/tests/`; it is the same argv, run alone.

## Detailed Requirements

- Baseline: `preflight.sh` reruns the gate command once when the first run exits non-zero, records only the second result in `baseline.json`, and prints one WARN line naming the retry and both exit codes.
- Post-merge (Step 6.5 item 4): one rule in `work-reference.md`, cited from the step: on a non-zero direct exit, run the exact argv once more, directly and unpiped; zero exit records green and continues; non-zero continues with the existing diagnostic and defer procedure using the second run's output as the fingerprint source.
- The run's progress output shows the retry as one line; the REQ's Testing section records both exit codes when a retry happened.
- No new status, no new flag, no new REQ type. `defer-gate`, `repository_gate_repair`, and the already-green path are untouched; they simply fire less often.
- Delete any sentence that would now contradict the rule (for example, wording that treats the first non-zero exit as final).

## Constraints

- Mechanics in the script, judgment in prose; no new prose walking a shell sequence (CLAUDE.md, prime-shell-commands.md).
- Never touch another session's claimed file under `do-work/working/`; stage explicit paths.
- The repository gate itself, `_dev/tests/maintainer-verify.sh`, is not edited by this REQ.

## Red-Green Proof
**RED prompt/case:** Run a REQ through `do-work run` while the baseline gate fails once transiently (the recorded shape: a broken pipe from a probe's `printf` while its reader has exited) and passes on the next run.
**Why RED now:** the first non-zero exit is final: the REQ is deferred, a `repository_gate_repair` REQ is minted and claimed, and its already-green no-op completion is the only thing that lets the parent resume. REQ-548 on 2026-09-03 is the record: 28 minutes, two gate runs, one archived REQ, zero code changed.
**GREEN when:** the same transient failure produces one immediate rerun; the green rerun is recorded, the REQ continues without deferral, no repair REQ exists, and the progress output carries one retry line. A gate that fails twice with the same fingerprint still defers exactly as today.
**Validation:** Inferred during capture; the maintainer approved the capture ("do 1, 2 and 3").

## Required Lessons — Dropped for Budget

- `_dev/primes/lessons-action-files.md` — 4050 tokens, over the 2000-token budget and `slugged: partial`, so no targeted form is legal. Matched because this REQ changes a pipeline step contract and its downstream readers.
- `_dev/primes/lessons-shell-commands.md` — 3385 tokens, over budget and `slugged: partial`. Matched because this REQ changes a shipped check script.

## Full Context
See `do-work/user-requests/UR-106/input.md` for complete verbatim input.

---

## Triage

**Route: B** - Medium

**Reasoning:** The change is well specified but its two edit sites — the gate launch inside `preflight.sh` and the post-merge rule in `work-reference.md` Step 6.5 — have to be located and read before editing. Outcome clear, location needs discovery.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Pre-Flight

**Git:** ✓ working tree clean outside `do-work/`
**Tests baseline:** ⚠ `bash _dev/tests/maintainer-verify.sh` red BEFORE any change — one pre-existing failure, `_dev/tests/session-start-hook-behavior.sh took 44s; each test file must finish under 30s`. A wall-clock budget miss on a slow container, no assertion failed. Recorded in `do-work/working/baseline-failures.txt` so Step 6.5 separates it from new regressions; not attributable to this REQ and not deferred to a repair REQ.
**Dependencies:** ✓ Go 1.26.1 and ShellCheck 0.11.0 provisioned for this session (container shipped Go 1.24.7 / no ShellCheck)

*Checked by work action*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/corehelpers/checks.go` (modify) — the baseline retry mechanics, at the site that actually launches the gate command
- `skills/do-work/tools/do-work-cli/internal/corehelpers/checks_test.go` (modify) — behavior probes pinning the bounded single retry
- `skills/do-work/actions/work-reference.md` (modify) — the one retry rule, plus the sentences that treated a first non-zero exit as final
- `skills/do-work/actions/work.md` (modify) — Step 5, Step 5.75 and Step 6.5 item 4 cite the rule

**Files I will NOT touch:** `_dev/tests/maintainer-verify.sh` (the gate itself), `skills/do-work/tools/checks/preflight.sh` (a five-line launcher with no logic to change), `CHANGELOG.md` and `VERSION` (Step 9 finalization).

**Acceptance criteria (restated from REQ):**
- [x] A red baseline gate is rerun once, immediately, with the exact same argv
- [x] Only the second result is recorded in `baseline.json`
- [x] One WARN line names the retry and both exit codes
- [x] A second failure enters the existing fingerprint / defer-gate path unchanged
- [x] Exactly one retry, no new status, flag, or REQ type

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/corehelpers/checks.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/corehelpers/checks_test.go` (modified)
- `skills/do-work/actions/work-reference.md` (modified)
- `skills/do-work/actions/work.md` (modified)

**What was done:** A new `runBaselineCommand` helper launches the gate argv, and `handlePreflight` calls it a second time on any non-zero exit, recording only the second run in `baseline.json` and `baseline-failures.txt` and emitting a `PREFLIGHT-BASELINE-RETRIED` warning that renders as one WARN line naming both exit statuses. A matching prose rule, keyed on the condition "any direct non-zero exit of the canonical gate argv" rather than on one call site, was added to `work-reference.md` under Repository Gate Deferral and Resumption and is cited from `work.md` Step 5, Step 5.75 and Step 6.5 item 4; the sentences that treated the first non-zero exit as final were rewritten.

## Decisions

**D-01 — DECIDE & STATE. The baseline retry lives in the Go handler, not in `preflight.sh`.** The write set named the launcher, but `preflight.sh` is five lines that exec `do-work-cli.sh preflight`; the command is executed by `handlePreflight` in `internal/corehelpers/checks.go`. Editing the launcher could not have produced the behavior. The site named in the write set behaves as the REQ specifies without being edited.

**D-02 — DECIDE & STATE. The prose rule is keyed on the condition, not on the Step 6.5 site alone.** Keying it only at the post-merge site would have left the pre-build baseline lane unchanged — the lane that actually deferred REQ-531 (adding the batch field) and minted REQ-548 (the already-green repair), which is the incident in this REQ's Why section. One rule in one place, stated as a condition per CLAUDE.md's "State conditions, not lists"; `work.md` cites it from both lanes.

**D-03 — DECIDE & STATE. Retry on any non-zero exit, including 126 and 127.** One condition, no exception list. A command that cannot launch fails instantly, so the second attempt costs nothing measurable, and the not-launched finding still fires off the second run. An exception list would buy no time and would go stale.

**D-04 — DECIDE & STATE. `baseline.json` keeps its three existing fields.** The REQ says record only the second result. A `first_exit_status` field would add schema surface every reader would then have to be taught to ignore; both statuses reach the Testing section through the WARN line, which is what the REQ asks for.

**D-05 — DECIDE & STATE. `runBaselineCommand` is a named package-level function, not an inline closure.** It is called twice with identical arguments, which is what "the same argv, run again" means; a named function makes that sameness structural instead of a copy-paste that can drift.

**D-06 — DECIDE & STATE. One optional line added to the Testing Section Template.** The REQ requires both exit codes in `## Testing` when a retry happened. Without a template slot that requirement lives only in a rule nobody rereads while filling in the template.

## Discovered Tasks

- **impact-low, report only** — `skills/do-work/docs/command-line-guide.md` line 46 describes `defer-gate`'s manifest as binding the "direct non-zero gate result". Still true, since that result is now the second run's, but a reader comparing it against the new rule may want it said explicitly.
- **impact-low, report only** — `preflightCompatibilityText` renders every other finding as a hard-coded literal but reads `finding.Evidence[0]` for the retry case. That is the only way to carry both exit statuses without a second data channel; worth a consistent evidence-carrying convention if a second finding ever needs the same thing.
- **impact-low, report only** — this repository uses the same script as both the Step 5 pre-flight test command and the Step 5.75 canonical gate, so a genuinely red suite can now cost up to four full runs in one REQ. Bounded and correct; if wall-clock ever becomes the complaint, reusing the pre-flight result for the baseline lane is the lever.

## Qualification

Passed — 4 files verified in the merge range `904b4d3..0e07aac`, 5 requirements traced, P-A-U confirmed. `qualify.sh` returned `OK: mechanical qualification passed`. Judgment checks: the `checks.go` change is substantive real logic (a named launcher called twice, a bounded retry, a new warning finding, a renderer case), not a stub; every acceptance criterion maps to a named file; nothing is hardcoded or stubbed.

## Testing

**Tests run:** `go test ./internal/corehelpers/ -run Preflight -count=1` and `bash _dev/tests/maintainer-verify.sh`
**Result:** ✓ Focused probes passing (3 tests). Repository gate passes everything except one pre-existing failure carried from the baseline.

**Repository gate retry:** first run exited 1, rerun exited 1 — the same failure both times, so the retry correctly did not mask it. Recorded here because this REQ is what added that retry.

**Pre-existing baseline failure (not attributable to this REQ):** `_dev/tests/session-start-hook-behavior.sh` finishes in 39-44s against a 30s per-file budget. It failed identically before any change (`do-work/working/baseline-failures.txt`) and no assertion inside it fails. It is a wall-clock budget miss on a slow container, so it was not deferred to a repository-gate repair REQ.

**Red-green validation:**
- `TestPreflightRerunsRedBaselineOnceAndRecordsOnlyTheRerun`: ✗ before implementation (`baseline command ran 1 times, want exactly 2 (one run plus one retry)`) → ✓ after
- `TestPreflightKeepsRedBaselineWhenTheSingleRetryAlsoFails`: ✗ before implementation (`baseline command ran 1 times, want exactly 2 — the retry is bounded at one`) → ✓ after

Both RED results were reproduced independently by the orchestrator, by restoring `checks.go` to the pre-merge revision with the new probes in place.

**New tests added:**
- `TestPreflightRerunsRedBaselineOnceAndRecordsOnlyTheRerun` — a command exiting 3 then 0 is invoked exactly twice; `baseline.json` holds `exit_status: 0`; no failures file is left; one retry line naming both statuses is printed.
- `TestPreflightKeepsRedBaselineWhenTheSingleRetryAlsoFails` — a command that always exits 3 is still invoked exactly twice, so the retry is bounded at one; `baseline.json` holds `exit_status: 3` and the failures file holds the second run's output.

**Environment note:** this container shipped Go 1.24.7, no ShellCheck, and `just` 1.21.0. The gate requires Go 1.26.1, ShellCheck 0.11.0, and a `just` new enough for `[positional-arguments]`. All three were provisioned before the gate could run at all; none of those gaps was caused by this REQ.
