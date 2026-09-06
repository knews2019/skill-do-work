# Builder brief — REQ-554, move the commit/inspect shared body into the prescribed-shell guide

**Worktree:** `/home/user/skill-do-work-worktrees/worktree-agent-REQ-554-move-commit-inspect-shared-body`
**Branch:** `worktree-agent-REQ-554-move-commit-inspect-shared-body`, based on `2876d96`
**Route:** B. **TDD: yes.** **Impact: negligible.** **`maintenance: true`.** **Estimate P50: 30 active minutes.**

## Read first

1. `do-work/runs/work-2026-09-05-231943/REQ-554-exploration.md` in the MAIN checkout — the authority. It has every
   file, every line number, the exact assertion text, and the three places the request itself is
   wrong.
2. `do-work/working/REQ-554-move-the-commit-and-inspect-shared-body-into-the-prescribed-shell-guide.md`
   in your worktree — the request, its `## Triage`, `## Exploration`, `## Scope` and `## Decisions`.
3. **Both `prime_files`, and they are binding:** `_dev/primes/prime-action-files.md` (the template every
   action file must carry — this is what makes the request's ceiling of 10 unreachable) and
   `_dev/primes/prime-shell-commands.md`.
4. Crew rules in `skills/do-work/crew-members/`: general.md, shared-principles.md,
   coding-guardrails.md, communication-style.md, **maintenance.md** (this REQ is
   `maintenance: true` — it edits shipped instruction files), anti-slop.md.

## What must become true

The inventory tag legend, the secret-shaped matching sentence, the four file-reading bullets, the
association semantics and both manual "do it by hand" fallbacks live **once**, in
`skills/do-work/docs/prescribed-shell-primitives.md`, with `commit.md` and `inspect.md` pointing at
them. Each action keeps the wording that is genuinely its own — the X and XD tag wordings differ
between a writing action and a read-only one, and collapsing them would be a behaviour claim dressed
as deduplication.

## Three things the request gets wrong — do not follow it literally

1. **The lock-in ceiling of "≤ 10 shared lines" is unreachable and would be red forever.** Deleting
   every shared *content* line still leaves 17, all of them scaffold `prime-action-files.md:70-80`
   requires both actions to carry. The prior exploration's suggested 20 is also unreachable — keeping
   the ASCII flow diagram and the Error Handling table row, both structural, the floor is 21. **Set
   the ceiling to the count you actually measure after the move**, and write the reason into the
   assertion's comment so the next reader sees why it is not zero.
2. **`_dev/tests/prescribed-shell-canonicalization.sh` counts nothing.** The request says to
   re-baseline counts there. That file has zero numeric constants: lines 66-83 are `grep -Fqx`
   membership checks over eleven headings, 85-117 are `grep -Fq` over sixteen pointer sites. The only
   edit it needs is your new section's heading added to the membership list at lines 66-77, matching
   the existing two-space indentation.
3. **The prior exploration's paste-ready assertion is green on day one.** Its `--glob '*/actions/*.md'`
   matches zero files, so its pipeline prints 0 today when the true count is 4. Use the assertion in
   THIS run's exploration, and prove it red before you accept it.

One thing the request feared does not happen: the canonicalization stale-pattern scan skips the
canonical guide itself (`prescribed-shell-canonicalization.sh:146`), so prose moved into the guide is
exempt. None of its nine stale patterns or seven old-implementation fragments matches any moved line.

## Scope — exactly five files

- `skills/do-work/docs/prescribed-shell-primitives.md`
- `skills/do-work/actions/commit.md`
- `skills/do-work-toolbox/actions/inspect.md`
- `_dev/tests/prescribed-shell-canonicalization.sh`
- `_dev/tests/audit-lockins.sh`

`_dev/tests/audit-lockins.sh` already carries a Finding 9 block added by REQ-552 earlier in this run.
Add yours beside it in the same shape; do not disturb the existing block.

## RED then GREEN

RED: your new lock-in assertion, run standalone against your base revision, must report the current
shared-line count above the ceiling and fail. Show that output before you accept the assertion.
GREEN: after the move, `bash _dev/tests/audit-lockins.sh` prints "Audit lock-in regressions passed."
and exits 0; `bash _dev/tests/prescribed-shell-canonicalization.sh` passes with the new heading; and
the canonical gate exits 0.

## Testing

```
bash _dev/tests/audit-lockins.sh
bash _dev/tests/prescribed-shell-canonicalization.sh
bash _dev/tests/contract-regressions.sh
env DO_WORK_MAINTAINER_TIER=heavy bash _dev/tests/prescribed-shell-scripts-behavior.sh
bash _dev/tests/action-shell-blocks.sh
```

Record each exit line. `action-shell-blocks.sh` matters because you are editing action files that
carry fenced shell blocks.

## Environment

Fresh cloud container. Wrap every gate and every test:

```
env -u NODE_OPTIONS \
  -u GIT_CONFIG_COUNT -u GIT_CONFIG_KEY_0 -u GIT_CONFIG_KEY_1 -u GIT_CONFIG_KEY_2 \
  -u GIT_CONFIG_VALUE_0 -u GIT_CONFIG_VALUE_1 -u GIT_CONFIG_VALUE_2 \
  GIT_CONFIG_GLOBAL=/tmp/claude-0/-home-user-skill-do-work/213e30ac-5958-56c8-9fd2-faaaaf9c4ea6/scratchpad/gitconfig-gate \
  QUEUE_KANBAN_BROWSER=/opt/pw-browsers/chromium \
  <command>
```

Capture exit status from `$?` directly — never pipe a gate to `tail`. This machine has 4 CPUs and
another builder is running: prefer the targeted probes above, and run the full canonical gate at most
once.

## Hand-back

Commit on your branch in the worktree, message ending
`Co-Authored-By: Claude Opus 5 <noreply@anthropic.com>`. Write your report to the ABSOLUTE path
`/home/user/skill-do-work/do-work/runs/work-2026-09-05-231943/REQ-554-handback.md` in the MAIN checkout. Do not stage
or commit it, do not push, do not touch any other worktree. Report the exact ceiling you measured and
what it is composed of.
