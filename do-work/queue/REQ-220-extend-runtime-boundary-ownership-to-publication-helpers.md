---
id: REQ-220
title: Extend runtime-boundary ownership to the remaining publication helpers
status: pending-answers
domain: general
created_at: 2026-08-17T18:37:31Z
user_request: UR-042
addendum_to: REQ-204
review_generated: true
effort_estimate: normal
sweep: true
sweep_key: runtime-boundary-ownership
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
maintenance: false
---

# Review Fix: Extend Runtime-Boundary Ownership to the Remaining Publication Helpers

## What

Close the same runtime-boundary defect class in the two shipped scripts REQ-204's copy-paste sweep found still carrying it: a helper that cleans files without owning the processes it started, and a directory publication that reads `mv`'s exit status as proof the rename landed where it was aimed.

## Context

REQ-204 fixed both boundaries in the `ai-report` prescribed batch block. `_dev/primes/prime-shell-commands.md` requires grepping a fixed primitive across every action before calling it fixed, because these blocks are copy-pasted. That sweep found the same root cause — verifying a convenient proxy instead of owning the complete boundary, exactly as `_dev/lessons/validated-runtime-boundaries.md` frames it — alive in two more places. Neither was in REQ-204's declared scope, so both were routed here rather than fixed inline.

## Instances

- [ ] `skills/do-work-toolbox/scripts/generate-report-image.sh`: backgrounds `imagegen` (line 54) and a watchdog subshell (line 73), but its HUP/INT/TERM traps are bare `trap 'exit N'`. Interrupted through REQ-204's batch it is killed by process group; interrupted as a direct invocation its backend outlives it.
- [ ] `skills/do-work-toolbox/scripts/install-last30days.sh`: `mv "$staging_directory" "$target_directory"` (line 94) follows a check-and-backup window. A target that reappears inside that window swallows the staging tree, `mv` exits 0, and the script sets `publication_complete=1` over a nested tree.
- [ ] `skills/do-work-toolbox/scripts/generate-report-image.sh` (line 48): `mv "$staged_output_path" "$output_path"` nests if the output path is ever a directory. Minor — both paths are caller-controlled inside private staging — but it belongs to the same audit.

## Requirements

- Reuse the already-proven idioms rather than inventing new ones: the verified-process-group pattern from `skills/do-work/scripts/run-blocked-check.sh`, and REQ-204's post-rename nesting verification.
- A helper that cannot prove process-group isolation must signal only the bare PID, never a group.
- A publication that nests must exit nonzero, leave the colliding destination byte-for-byte, and leave no private stage behind.
- Preserve every existing success path and exit status of both scripts.
- Add fixture replays to `_dev/tests/prescribed-shell-scripts-behavior.sh` for each instance, following the two REQ-204 added.

## Red-Green Proof

**RED prompt/case:** Interrupt `generate-report-image.sh` directly while its backend runs and observe the backend survive; run `install-last30days.sh` with a `mv` shim that recreates the target inside the check-then-rename window and observe it report a complete publication over a nested tree.
**Why RED now:** File cleanup does not own a process tree, and `mv`'s exit status does not prove exact-path publication when the destination is a directory. Both scripts ship to users today.
**GREEN when:** The interruption replay proves no launched process survives a direct invocation; the collision replay returns nonzero, preserves the destination byte-for-byte, and leaves no nested or private stage; and every pre-existing named case in the suite still passes.
**Validation:** Review finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Open Questions

- [ ] REQ-204 fixed both runtime boundaries in the ai-report batch, and the required copy-paste sweep found the same root cause in two more shipped scripts. Should I process this as a new task?
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.
  Why this is yours: REQ-204 was itself a generation-two review follow-up, so the cascade-depth rule requires your consent before another autonomous repair cycle. Worth knowing before you decide: these two scripts ship to users, and neither defect is reachable from the path REQ-204 already fixed — the first only bites a direct invocation, the second only a reappearing target.
