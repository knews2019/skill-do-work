---
id: REQ-539
title: 'Cut the contract file to the incident core and split the aggregate into fast and heavy'
status: pending
priority: now
created_at: 2026-09-03T14:49:02Z
user_request: UR-104
domain: testing
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-537, REQ-538, REQ-539, REQ-540, REQ-541, REQ-542]
batch: two-tier-gate
write_set:
  - _dev/tests/contract-regressions.sh
  - _dev/tests/probe-batch.sh
  - _dev/tests/*.sh
  - _dev/tests/contracts/
status_changed_at: 2026-09-03T23:32:33Z
---

# Cut the Contract File to the Incident Core and Split the Aggregate Into Fast and Heavy

## What

Delete every sentence assertion in `_dev/tests/contract-regressions.sh` whose subject is a Markdown sentence unless its comment names a specific incident REQ and its target is a shipped file under `skills/`; delete the extracted-block pins (`$..._block`) wholesale; delete `python3` heredoc checks a `grep` or Go test already covers. Move survivors into per-owner files sourced by the aggregate. Classify each behavioral probe fast or heavy from measured wall time, build do-work-cli once at the top of the heavy aggregate, and land REQ-519's line-count ratchet on the fast file.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Context

- 8,451 lines, 294 assertion sites, 27 python heredocs, 6 `git init` fixtures. Top targets: `$verify_revalidation_block` 19, `$crash_recovery_block` 10, `$review_archived_input_block` 9, `work-reference.md` 9, `work.md` 8.
- Aggregate wall 112 s, bounded by one lane in `probe-batch.sh` where update-script (about 60 s) and install-suite both `go build -o do-work-cli` into the same source-tree path.
- The maintainer's rule: a sentence pin on prose is not a real failure (D3). The Fold-First prose-only test in `capture-reference.md` already says prose discrepancies are backlog items, not REQs; the contract file is the last place that treats them as failures.

## Detailed Requirements

- Measured alone on 2026-09-03 (seconds, probe): 60s update-script-behavior; 40s install-suite-behavior; 18s staged-skills-contract; 11s session-start-hook-behavior; 5s prescribed-shell-canonicalization; 4s action-shell-blocks; 3s do-work-cli-launcher-behavior; 1s suite-manifest-contract; 1s shipped-package-reference-contract; 1s p50-estimator-determinism; 0s select-simple-reqs-behavior; 0s defensive-surface-audit. Re-measure in the commit body after the split. Fast: under about 10 s and no binary or synthetic install build. Heavy: the rest (expected: update-script, install-suite, launcher, staged-skills with its prescribed-shell cases).
- Survivor rule applied line by line; the deleted list goes in the commit body grouped by target file.
- Split survivors into `_dev/tests/contracts/<owner>.sh`, one per owning action or module, each sourced by the aggregate; the aggregate itself keeps only launch, collect and the ratchet.
- Heavy aggregate builds do-work-cli once into a private path and hands it to both the updater and installer probes so they no longer share a lane.
- Ratchet: the fast contract file's line count is recorded in the file and a run fails when the count grows.

## Constraints

- Land in place, not through `do-work run`; one integrating commit with version bump and changelog entry; prove it with one `bash _dev/tests/gate-runner.sh --once`.
- Delete before you add; every deleted test is listed in the commit body with the failure it pinned and why it no longer earns its cost. No new sentence pins, no new prose that walks a shell sequence.
- Never touch another session's claimed file under `do-work/working/`; stage explicit paths.

## Red-Green Proof
**RED prompt/case:** `wc -l _dev/tests/contract-regressions.sh` and `time bash _dev/tests/contract-regressions.sh`.
**Why RED now:** 8,451 lines and about 112 s.
**GREEN when:** the fast contract file is under 1,500 lines and runs under 30 s wall; the heavy aggregate still exercises the installer, updater, launcher and prescribed-shell cases and exits 0 under `--heavy`.
**Validation:** User confirmed (D3)

## Required Lessons — Dropped for Budget

- `_dev/primes/lessons-shell-commands.md` — 3385 tokens, over the 2000-token budget and `slugged: partial`, so no targeted form is legal. Matched because this REQ changes shipped-shell fixtures and prescribed command probes.

## Full Context
See `do-work/user-requests/UR-104/input.md` for complete verbatim input.


## Addendum (2026-09-03)

User added (2026-09-03 21:29 local, in the test-budget session; 22:05 local, "update the batch REQs per A1-A3 via queued addenda" referring to the velocity report at `ai-reports/2026-09-03_2145_do-work-velocity-and-pending-queue-speed/`, item A2):

> ```
> each test file should finish under 30 seconds (use the 80% value 20% effort principle until this is obtained)
>
> The rest of the test are accessible only when calling them with the --heavy parameter.
>
> the catch to call the --heavy parameter need to ask user for permission, and it should not block anything, meaning that where --heavy is required those tasks go into pending-testing status.
>
> Also because these tests tend to ballon, make sure to always measure the test duration, and adjust when the limits are reached.
> ```

- Per-file budget, shell side: every test file the fast tier runs finishes under 30 s wall. The existing "under about 10 s" classification line stays as the fast/heavy sorting rule; 30 s is the hard ceiling the gate enforces on whatever is classified fast.
- The budget is a program, not a sentence: the aggregate (or `maintainer-verify.sh`, whichever already owns the timing loop) measures each probe file's wall time, prints it, appends one row per file per run to a durations log beside `do-work/calibration-log.tsv` (run id, file, seconds, and the number of other gate processes running on the machine so an over-budget row under load can be read correctly), and fails the fast tier when a fast-classified file exceeds 30 s. `--heavy` has no budget. Same shape as the line-count ratchet this REQ already lands; land both in the same place.
- Measured offenders on 2026-09-03 (`update-script-behavior.sh` about 60 s, `install-suite-behavior.sh` about 40 s) are heavy or cut, never left fast-with-an-exception; the maintainer's preference is honest severity over carve-outs.
- GREEN additionally requires: every fast-tier probe file's recorded duration under 30 s, and the durations log carrying one row per fast-tier file after the proving `gate-runner.sh --once`.
- Coherence check: no contradiction with the original sections; the addendum tightens the GREEN condition and adds the measurement mechanism.

## Addendum (2026-09-03, 23:10 local)

User added (23:00 local, "do 1, 2 and 3, I'll release the cloud claims", item 3 of the velocity report's handoff: "narrow REQ-539 to whatever 21dac2b8 did not do"):

- Release 0.271.0 (21dac2b8) landed part of this REQ: the installer, updater, staged-package and expensive subprocess probes moved behind `--heavy`; the shell probe runner measures each standalone script against the 30 s ceiling and fails the fast tier over it; the installer and updater fixtures were narrowed. Do not redo those.
- Still open on main and now the whole scope of this REQ: (1) the per-owner split of `_dev/tests/contract-regressions.sh` into files sourced by the aggregate, with the sentence-pin deletions the original What describes; (2) the durations log beside `do-work/calibration-log.tsv` with one row per file per run and the concurrent gate count (the 22:04 addendum), which 0.271.0 measures but does not persist; (3) the line-count ratchet on the fast contract file (UR-100's ceiling; the file is 8,479 lines and nothing on main fails when it grows); (4) confirm whether the heavy aggregate builds do-work-cli once, and land it if not.
- `depends_on` changed from `[REQ-538]` to `[]`: REQ-538's code landed in 21dac2b8 and its record is being cancelled as landed in place, and a cancelled dependency never satisfies gating.
- Coherence check: the original GREEN (fast contract file under 1,500 lines and under 30 s, heavy aggregate exits 0 under `--heavy`) still holds; the 30 s half is already true on main, the line-count half is not.
