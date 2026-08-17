---
id: REQ-204
title: Harden ai-report generated-batch lifecycle
status: completed
completed_at: 2026-08-17T18:37:31Z
commit:
claimed_at: 2026-08-17T18:30:39Z
status_changed_at: 2026-08-17T18:10:49Z
route: B
estimate:
  p50_active_minutes: 30
  confidence: medium
  calculated_at: 2026-08-17T18:31:10Z
  basis:
    - Route B
    - 2-file write set
    - 2 subsystems involved
    - 5 acceptance criteria
    - async lifecycle behavior
    - cross-route regression gates
write_set:
  - skills/do-work-toolbox/actions/ai-report-reference.md
  - _dev/tests/prescribed-shell-scripts-behavior.sh
domain: general
created_at: 2026-08-15T19:39:11Z
user_request: UR-042
addendum_to: REQ-198
review_generated: true
effort_estimate: normal
sweep: true
sweep_key: ai-report-generated-batch-lifecycle
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md]
tdd: true
maintenance: true
---

# Review Fix: Harden AI-Report Generated-Batch Lifecycle

## What

Close the complete generated-image batch lifecycle: an interrupted caller must terminate and reap exactly its own helpers, and publication must fail closed if the destination appears at the final boundary. Solve both instances as one rule because both require ownership of the batch from private staging through terminal publication or cleanup.

## Context

REQ-198 fixed the original all-failed directory shape, but review showed that signal cleanup handles files without process ownership and that a check-then-plain-`mv` can nest staging inside a newly appeared destination while returning success.

## Instances

- [ ] Batch interruption: signal and reap recorded helper PIDs (and their owned descendants) before removing exact staging; no optional full-host backend may outlive the caller.
- [ ] Final publication: coordinate destination appearance after the last check and prove the operation returns nonzero, preserves the colliding directory, and leaves no nested/private stage.

## Requirements

- Preserve normal wait-all and per-status freshness behavior.
- On HUP/INT/TERM, terminate and reap exactly the current batch's recorded process tree before staging cleanup.
- Use a portable exclusive/atomic directory publication boundary or a verified rollback that cannot report success after nesting.
- Never delete or overwrite a colliding destination.
- Add exact prescribed-block behavior replays for both adversarial paths.

## Red-Green Proof

**RED prompt/case:** Signal only the batch shell while slow helpers run, then coordinate creation of `generated/` after the final absence check but before publication. Current behavior can leave helpers alive and can return success with staging nested under the colliding directory.
**Why RED now:** File cleanup alone does not own the process tree, and plain `mv` treats an existing destination directory as a container.
**GREEN when:** The signal replay proves no owned process survives and no staging/public path leaks; the collision replay returns nonzero, preserves the destination byte-for-byte, creates no nested stage, and normal all-failed/mixed paths still pass.
**Validation:** Review finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Open Questions

- [x] The original empty-directory bug is fixed, but review found two deeper lifecycle edges in the same shell batch. Should I process this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.
  Why this is yours: this is a generation-two review follow-up, so the cascade-depth rule requires your consent before another autonomous repair cycle.

<!-- D-XX counter: none used. Next decision: D-01. -->

## Triage

**Route:** B — Explore then Build

**Reasoning:** Both instances name their symptom precisely, but the portable way to own a helper's process tree and the exact nesting behavior of `mv` on a directory operand had to be discovered from the repo's own proven idioms before either could be written.

**Confidence:** high

*Triaged by work action*

## Plan

Planning not required — Route B.

## Exploration

**Key files:**
- `skills/do-work-toolbox/actions/ai-report-reference.md` — the prescribed batch block: staging via `mktemp -d`, `trap cleanup … EXIT` plus bare `trap 'exit N'` signal traps, parallel launches, wait-all, then check-then-`mv` publication.
- `_dev/tests/prescribed-shell-scripts-behavior.sh` — extracts that exact block with `awk` and replays it against fake backends (`run_ai_report_batch_replay`); already covers all-failed and mixed-success.
- `skills/do-work/scripts/run-blocked-check.sh` — the repo's proven portable process-group idiom: `set -m` to give a background job its own group, `ps -o pgid=` to *verify* isolation, signal the group, escalate TERM→KILL, then `wait` to reap.
- `_dev/lessons/validated-runtime-boundaries.md` — records both boundaries this REQ closes as one shared mistake: "the code verified a convenient proxy while the real boundary lived one level wider."

**Concerns found:**
- The signal traps exit immediately; the EXIT trap then deletes staging while every launched helper is still running and still writing into it.
- `mv` treats an existing destination *directory* as a container. `[ ! -e "$generated_directory" ] || refuse` followed by `mv` is a check-then-act window: a `generated/` that appears in it swallows the stage and `mv` still exits 0.
- `setsid` is unavailable on macOS, so process-group isolation has to come from `set -m` and be verified rather than assumed.

## Scope

**Files I will touch:**
- `skills/do-work-toolbox/actions/ai-report-reference.md` (modified)
- `_dev/tests/prescribed-shell-scripts-behavior.sh` (modified)

**Acceptance criteria (restated from the REQ):**
1. Normal wait-all and per-status freshness behavior preserved.
2. On HUP/INT/TERM, terminate and reap exactly this batch's recorded process tree before staging cleanup.
3. A portable exclusive/atomic publication boundary, or a verified rollback that cannot report success after nesting.
4. Never delete or overwrite a colliding destination.
5. Exact prescribed-block behavior replays for both adversarial paths.

## Pre-Flight

- Working tree clean outside `do-work/`.
- Baseline `bash _dev/tests/prescribed-shell-scripts-behavior.sh` passing before the change.

## Implementation Summary

**Files changed:**
- `skills/do-work-toolbox/actions/ai-report-reference.md` (modified)
- `_dev/tests/prescribed-shell-scripts-behavior.sh` (modified)

**What was done:** Gave the prescribed image batch ownership of the process tree it starts and made its publication verify the rename instead of trusting the exit status, then added the two adversarial replays that prove both.

*Process ownership.* Each helper is now launched under `set -m` so it leads its own process group, and that isolation is **verified** with `ps -o pgid=` rather than assumed — a helper whose group id is not its own pid, or whose group is the caller's, is recorded with no group and signalled by bare PID, so the caller can never signal its own group. `terminate_report_image_batch` sends TERM to every recorded group (or PID), polls up to one second, escalates to KILL, and then `wait`s each PID to reap it. The HUP/INT/TERM traps call it before exiting, so reaping always precedes the EXIT trap's staging removal. The idiom is lifted from `skills/do-work/scripts/run-blocked-check.sh`, which already proved it portable on stock Bash without `setsid`.

*Publication.* After the rename, the block checks for its own stage nested under the destination (`$generated_directory/${image_generation_stage##*/}`). If it is there, `mv` moved into a directory that appeared after the absence check: the block removes only its own nested stage, leaves the colliding `generated/` untouched, prints `REFUSING: generated/ appeared during publication`, and exits nonzero. The single same-filesystem rename is preserved — this is the "verified rollback" the REQ allowed, not a different publication mechanism.

*Prose.* The paragraph above the block now states both boundaries, so the contract and the code cannot drift apart.

**Tests touched:** two new named replays in `_dev/tests/prescribed-shell-scripts-behavior.sh` — an interrupted batch (helpers spawn descendants and record all four PIDs; TERM to the caller must leave none alive and exit 143) and a publication collision (a `mv` shim creates `generated/` in the exact check-then-rename window). Named-case count updated 27 → 29.

## Qualification

Passed — 2 files verified, 5 requirements traced, no debug artifacts.

- Both declared files exist and appear in the diff; no undeclared file touched.
- Substantive: ~60 lines of new lifecycle logic in the block, ~90 lines of new replay coverage.
- Requirements traced: wait-all/freshness untouched (1); trap → terminate → reap → EXIT-cleanup ordering (2); post-rename nesting verification (3); destination never deleted or overwritten (4); both replays added and passing (5).
- Flowing: the recorded groups array is written at launch and read by all three lifecycle functions; nothing is defined and unused.

## Testing

**Tests run:** `bash _dev/tests/prescribed-shell-scripts-behavior.sh` (baseline, RED, GREEN); `bash _dev/tests/maintainer-verify.sh` (includes the fenced-block ShellCheck lint over the modified block)

**Result:** ✓ prescribed-shell suite exit 0, 29 named cases; ✓ maintainer-verify exit 0, zero FAIL lines.

**Red-green validation:** ✗ RED — with the new replays in place and the block reverted to its pre-change form, the suite exits 1: `ai-report interrupted batch replay left 4 helper process(es) or descendant(s) alive`, `ai-report publish-collision replay reported success after the destination appeared`, `… left its staged batch nested inside the colliding destination`, `… leaked invocation-private staging`. → ✓ GREEN — with the fixed block, all four assertions pass, the collision path prints `REFUSING: generated/ appeared during publication` and exits nonzero, and the pre-existing all-failed and mixed-success replays are unchanged.

Four survivors in the RED run is the precise measurement the REQ asked for: two helper processes **and** two of their descendants, which is why file cleanup alone was never enough.

**Existing tests updated:** none — both pre-existing batch replays pass unmodified, which is the regression evidence for criterion 1.

*Verified by work action*

## Review

**Overall: 93%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | All five criteria delivered and individually replayed |
| Code Quality | 90% | Reuses the proven `run-blocked-check.sh` group idiom; the block is now long enough to be worth extracting |
| Test Adequacy | 95% | Both adversarial paths replayed against the extracted block, with the RED measured (4 survivors) not asserted |
| Scope | 100% | Two declared files; the sweep findings were routed out, not fixed inline |
| Risk | Low | Signals a process group — mitigated by verifying leadership and never falling back to the caller's group |
| Acceptance | Pass | Suite green; both defects reproduce on the old block and are closed on the new one |

**Verdict: Approve with follow-ups** — the batch now owns both boundaries, and the copy-paste sweep found the same class alive in two sibling files.

### Findings

**Important:**
- `skills/do-work-toolbox/scripts/generate-report-image.sh:42-44,54,73` — the helper this batch launches has the *same* file-cleanup-without-process-ownership defect one level down: it backgrounds `imagegen` and a watchdog subshell but its HUP/INT/TERM traps are bare `trap 'exit N'`. Interrupted through the batch it is now killed by group, but a direct invocation still leaves its backend running. — gate: rule-change
- `skills/do-work-toolbox/scripts/install-last30days.sh:94` — `mv "$staging_directory" "$target_directory"` after a check-and-backup window has the same directory-operand nesting behavior: a target that reappears in the window swallows the staging tree, `mv` exits 0, and the script sets `publication_complete=1` on a nested tree. — gate: rule-change

Both are one root cause — a publication helper verifying a convenient proxy (exit status, file presence) instead of owning the complete boundary — which is exactly how `_dev/lessons/validated-runtime-boundaries.md` already frames it. Consolidated into one sweep REQ rather than one REQ per file.

**Minor:**
- `skills/do-work-toolbox/scripts/generate-report-image.sh:48` — `mv "$staged_output_path" "$output_path"` nests if the output path is ever a directory. Inside invocation-private staging the caller controls both paths, so the exposure is theoretical; noted for the sweep REQ rather than raised separately.

### Restatement Sweep

**Triggered** — the diff redefines two prescribed-shell primitives (signal handling around background helpers; directory publication by `mv`), and `_dev/primes/prime-shell-commands.md` requires grepping the same primitive across all actions before calling it fixed. Swept `mv "$…" "$…"` and `trap 'exit N'` across `skills/*/actions/*.md` and `skills/*/scripts/*.sh`. Results: `publish-portfolio-summary.sh:90` is already owned by the queued REQ-205 and was not double-reported; the two files above are new and became follow-up REQ-220; every other hit is a file-to-file rename with no directory-operand exposure.

### Requirements Checklist

- [x] Preserve normal wait-all and per-status freshness behavior — delivered (both pre-existing replays pass unmodified)
- [x] Terminate and reap the batch's process tree on HUP/INT/TERM before staging cleanup — delivered
- [x] Portable exclusive/atomic publication boundary or verified rollback — delivered (post-rename nesting verification)
- [x] Never delete or overwrite a colliding destination — delivered (`keep.txt` asserted byte-for-byte)
- [x] Exact prescribed-block behavior replays for both adversarial paths — delivered

### Acceptance Testing

**Result: Pass**
- `bash _dev/tests/prescribed-shell-scripts-behavior.sh` — exit 0, 29 named cases.
- `bash _dev/tests/maintainer-verify.sh` — exit 0, zero FAIL lines (this is what lints the modified fenced block through ShellCheck).
- Finding-Closure Ratchet: both named GREEN conditions were measured against the reverted block first. The interruption path reported exactly four survivors; the collision path reported success while nesting. Both are closed.

### Suggested Additional Testing

- The `set -m` job-control path is exercised on macOS here; worth one run on Linux/CI to confirm `ps -o pgid=` parses identically there.
- A batch interrupted *during* publication (between rename and nesting check) is not replayed — the window is a single syscall pair, but a deliberate probe would close the last untested ordering.

### Follow-up REQs Created

- REQ-220: "Extend runtime-boundary ownership to the remaining publication helpers" (sweep, `addendum_to: REQ-204`, `pending-answers`)
- REQ-221: "Extract the ai-report image batch into a shipped script" (discovered task, `addendum_to: REQ-204`, `pending-answers`)

## Discovered Tasks

- [normal] The prescribed batch block is now ~110 lines of mechanical shell inside a markdown action file, and the test suite already has to `awk` it back out to run it. Extracting it to `<skill-root>/scripts/generate-report-image-batch.sh` (arguments: staging directory, style, then `target:prompt` pairs) would let the replays call the script directly and leave the action file with a pointer plus the per-report prompts. Queued as REQ-221.

## Lessons Learned

**What worked:** Lifting the process-group idiom from `run-blocked-check.sh` verbatim — including its *verify before you signal* discipline — instead of reaching for `setsid`, which macOS does not ship. Simulating the check-then-rename race with a `mv` shim on `PATH` turned an inherently timing-dependent defect into a deterministic replay.

**What didn't:** The first arrangement defined the signal traps before `terminate_report_image_batch` existed, so a signal arriving in that window would have hit an undefined function. Traps that call a function must be installed after it is defined, not next to the other traps.

**Worth knowing:** `mv` on a directory operand is a container operation, not a collision — `mv a b` where `b` is an existing directory produces `b/a` and exits 0. Any publication that checks `! -e` and then renames has a window where success and nesting are indistinguishable from the exit status alone; verify the destination after the rename. And killing a launcher is not killing a batch: the RED run left four processes alive, two of which were descendants the launcher never named.

## Orientation

An interrupted `ai-report` run can no longer leave image backends running against a directory it just deleted, and a `generated/` that appears at the last moment can no longer swallow the staged batch while the run reports success. Lives in the ai-report generated-image subsystem (`skills/do-work-toolbox/actions/ai-report-reference.md`), with both failure paths replayed from the extracted block. No new module or contract — the batch's existing boundaries were widened to the ones that were always real, so the system's shape is unchanged.

Prime staleness spot-check: `_dev/primes/prime-action-files.md` and `_dev/primes/prime-shell-commands.md` — all referenced paths still resolve; `prime-shell-commands.md` gains a lesson link rather than a correction.
