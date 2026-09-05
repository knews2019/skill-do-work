# Builder brief — REQ-592, seal the do-work tree into both fast gate stages

**Worktree:** `/home/user/skill-do-work-worktrees/worktree-agent-REQ-592-seal-the-do-work-tree-into-both-fast-gate-stages`
**Branch:** `worktree-agent-REQ-592-seal-the-do-work-tree-into-both-fast-gate-stages`, based on `15e2ec3`
**Route:** B. **TDD: yes.** **Impact: critical.** **Estimate P50: 30 active minutes, medium confidence.**

## What must become true

A change to any `do-work/` path that a fast gate stage reads must force that stage to execute.
Today both stages skip the whole `do-work/` tree in their seal, so a `do-work/`-only edit reuses
stale evidence and the gate prints `Maintainer verification passed.` and exits 0 while the stage's
own test fails on that same tree.

## The design the exploration settled

Read `do-work/runs/work-2026-09-05-231943/REQ-592-exploration.md` in the MAIN checkout for the full evidence. The
short version, and the reason the obvious fix is wrong:

Deleting the `non_stage_coverage: [do-work]` entry alone makes every `do-work/` path *unclassified*,
and an unclassified path is sealed into **every** stage — so every queue edit would re-run the
do-work-cli stage too, whose only real `do-work/` input is one file. Instead:

1. Declare `do-work` as the **queue-kanban stage's coverage**. That stage builds the board from the
   real tree, so it genuinely reads it; and `fastStageManifestClassifiesPath` then keeps those paths
   out of the do-work-cli stage's seal.
2. Declare `do-work/archive/UR-003/input.md` as an **exact** coverage rule on the do-work-cli stage.
   That is its one real `do-work/` read (`repository_model_test.go:397-417`, a 5608-byte assertion).
3. `non_stage_coverage` becomes `[]`. Verified, not assumed: no `do-work/` subtree is unread by both
   stages once those two rules exist.
4. Add a new `seal_exclusions` list to the manifest — "sealed nowhere, even where a stage covers it"
   — tested **before** the `stageCovered` test, exactly where the `queueStatePrefix` guard sits now.
   It is needed because once `do-work` is queue-kanban's coverage, the untracked-ignored skip stops
   protecting `do-work/test-durations.tsv`, which the stage writes itself and whose recorded
   fingerprint is the pre-run one — so a seal over it can never match.

## Scope — exactly four files

- `_dev/tests/fast-stages.json`
- `skills/do-work/tools/do-work-cli/internal/heavyverification/fast_stage_evidence.go`
- `skills/do-work/tools/do-work-cli/internal/heavyverification/fast_stage_evidence_test.go`
- `_dev/tests/fast-stage-reuse-behavior.sh`

Do not touch `heavy_run.go`, `heavy_evidence.go`, `_dev/tests/heavy-lanes.json`, or
`_dev/tests/maintainer-verify.sh`. The heavy lane has a related gap; it is a discovered task, not
this request.

## Three traps the exploration found

- **The two fixtures' toolchain probe inputs live under `do-work/`.** Once `do-work` is sealed into
  the fixture's stage, the `toolchain probe output changed` case still passes but for the wrong
  reason — the file's own byte seal moved, not the probe output — and `toolchain probe cannot run`
  reaches `fingerprint_uncertain` through a missing seal input instead of a failing probe. Move the
  probe inputs into the excluded subtree. This is not cosmetic; both cases silently stop testing
  what they name otherwise.
- **`do-work/.req-reservations/` is tracked**, 162 marker files the allocator creates and removes
  during ordinary work. Putting it in `non_stage_coverage` does nothing once `do-work` is a stage's
  coverage — the tracked loop seals on `stageCovered || !classified` and `stageCovered` wins. It
  belongs in `seal_exclusions`.
- **`seal_exclusions` is a closed enumeration**, the shape `_dev/primes/prime-shell-commands.md`
  § Closed Enumerations Go Stale warns about. State the admission CONDITION in the Go struct doc — a
  `do-work` path written by the gate or the orchestrator *while a gate runs* and byte-unread by every
  stage — so the next person has a test to apply rather than a list to copy.

## Red-Green proof you must produce

RED first, then GREEN, both recorded in the hand-back with the exact commands and output.

Fast loop, seconds, develop against this:
```
cd /home/user/skill-do-work-worktrees/worktree-agent-REQ-592-seal-the-do-work-tree-into-both-fast-gate-stages/skills/do-work/tools/do-work-cli
go test -count=1 -run TestFastStageReuseDecisionTable ./internal/heavyverification/
bash <worktree>/_dev/tests/fast-stage-reuse-behavior.sh
```

Gate level, the request's own proof — run it in the worktree:
```
bash _dev/tests/maintainer-verify.sh              # run 1: records evidence
printf '\n' >> do-work/archive/UR-003/input.md    # 5608 -> 5609 bytes
bash _dev/tests/maintainer-verify.sh              # run 2
```
RED now: run 2 prints `REUSED (fingerprint_match...)` for `do-work-cli-fast-tests` and
`Maintainer verification passed.`, exit 0, while on the same tree
`go test -short -count=1 -run TestDiscoverRepositoryAcceptsProductionLegacyArchiveInputClass ./internal/repositorymodel/`
fails with `production legacy fixture changed size: got 5609 bytes`.
GREEN after: run 2 prints `EXECUTING (fingerprint_mismatch)` for that stage and the gate exits
non-zero. Then reset the file, re-warm, append one row to `do-work/test-durations.tsv` only, and the
gate must still print REUSED for both stages and exit 0.

## Environment — read this before running any gate

This is a fresh cloud container. Run every gate and every Go test through the sanitized environment,
or you will chase failures that are not yours:

```
env -u NODE_OPTIONS \
  -u GIT_CONFIG_COUNT -u GIT_CONFIG_KEY_0 -u GIT_CONFIG_KEY_1 -u GIT_CONFIG_KEY_2 \
  -u GIT_CONFIG_VALUE_0 -u GIT_CONFIG_VALUE_1 -u GIT_CONFIG_VALUE_2 \
  GIT_CONFIG_GLOBAL=/tmp/claude-0/-home-user-skill-do-work/213e30ac-5958-56c8-9fd2-faaaaf9c4ea6/scratchpad/gitconfig-gate \
  QUEUE_KANBAN_BROWSER=/opt/pw-browsers/chromium \
  bash _dev/tests/maintainer-verify.sh
```

Why each: `NODE_OPTIONS` and the `GIT_CONFIG_*` triple are refused by
`_dev/tests/heavy-runtime-fingerprint.py` as opaque runtime extensions, which makes every fingerprint
uncertain. The global `commit.gpgsign=true` points at an empty signing key, so every `git commit`
inside a test fixture repository fails and lane tests see a dirty tree where they expect a new
revision.

**Capture the gate's exit status directly.** Never pipe it to `tail` — the shell reports the pipe's
status and a red gate reads as green.

## Hand-back

Write your report to the ABSOLUTE path
`/home/user/skill-do-work/do-work/runs/work-2026-09-05-231943/REQ-592-handback.md` (the MAIN checkout, not the
worktree — a relative path lands on your branch and the orchestrator reads nothing). Do not stage or
commit it. Sections: `## Implementation Summary`, `## Decisions` (D-XX, next available is D-01),
`## Red-Green Evidence`, `## Testing`, `## Discovered Tasks`, `## Declared but not touched`.

Commit your work on your branch in the worktree. Do not commit anything in the main checkout.
