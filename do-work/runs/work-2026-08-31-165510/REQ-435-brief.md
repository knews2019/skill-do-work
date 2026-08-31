# REQ-435 Builder Brief

Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-435-complete-doctor-forensics-delegation-contract`
Branch: `worktree-agent-REQ-435-complete-doctor-forensics-delegation-contract`
REQ: `do-work/working/REQ-435-complete-doctor-forensics-delegation-contract.md`
Exploration: `do-work/runs/work-2026-08-31-165510/REQ-435-exploration.md`
Integration owner: main checkout; do not edit `do-work/`, release metadata, changelogs, or version files.

## Frozen write set

- `skills/do-work/actions/forensics.md`
- `skills/do-work/actions/work-reference.md`
- `skills/do-work/actions/abandon.md`
- `skills/do-work/scripts/repair-req-timestamps.sh`
- `skills/do-work-board/tools/queue-kanban/verify.go`
- `skills/do-work-board/tools/queue-kanban/model.go`
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md`
- `skills/do-work/tools/do-work-cli/internal/doctor/doctor_commands_test.go`

## Required method

- Read the REQ, exploration report, both listed primes and lessons, `_dev/primes/prime-shell-commands.md`, and the queue-kanban prime/lessons plus applicable general/testing/coding guardrails.
- Extend the existing doctor action-delegation contract test first and capture a real RED for missing report authority/stale consumer pointers.
- Apply the default constraint: delete the unused Queue/Archive/Working totals; do not add doctor/result production schema.
- Preserve the authority split: doctor mechanical findings, Crash Recovery takeover judgment, Recurring Corrections judgment, and queue-kanban verify.
- Replace only live deleted-check references in the frozen scope with stable anchors, commands, headings, or finding codes. Historical evidence remains unchanged.
- Stay exactly inside the eight-file scope. Run focused doctor tests, full CLI module/vet, queue-kanban tests/vet as applicable, shipped reference and contract-regression scripts, exact Go 1.25 compatibility, shellcheck for the changed script, and diff checks.
- Commit on the branch, keep the worktree clean, and write handback to main checkout path `do-work/runs/work-2026-08-31-165510/REQ-435-handback.md`.

## Handback

Record full commit, RED/GREEN evidence, every changed file, commands/exits, authority mapping, stable reference inventory, decisions, and discoveries. Do not create or edit queue files.
