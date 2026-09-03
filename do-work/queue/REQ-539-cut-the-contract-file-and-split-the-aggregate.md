---
id: REQ-539
title: 'Cut the contract file to the incident core and split the aggregate into fast and heavy'
status: pending
created_at: 2026-09-03T14:49:02Z
user_request: UR-104
domain: testing
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: false
suggested_spec:
depends_on: [REQ-538]
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

