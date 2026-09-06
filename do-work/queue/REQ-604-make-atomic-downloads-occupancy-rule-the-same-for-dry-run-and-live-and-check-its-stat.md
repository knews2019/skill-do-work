---
id: REQ-604
status: pending
domain: general
created_at: 2026-09-06T08:19:05Z
user_request: UR-105
review_generated: true
impact: impact-user-visible
effort_estimate: effort-mechanical
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
related: [REQ-597]
write_set: [skills/do-work/tools/do-work-cli/internal/corehelpers/commands.go, skills/do-work/tools/do-work-cli/internal/corehelpers/commands_test.go, skills/do-work/docs/prescribed-shell-primitives.md]
title: 'Make atomic-download refuse an occupied target the same way in dry-run and live, and check its stat'
---

# Make Atomic-Download Refuse an Occupied Target the Same Way in Dry-Run and Live, and Check Its Stat

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## What

Measured by REQ-597's guide builder (fixtures A2/A4 under its scratch directory, evidence in
`do-work/runs/work-2026-09-05-231943/REQ-597-handback.md`): `atomic-download --dry-run` refuses any
existing target with exit 2 and `target already exists` (`commands.go:863-865`), while the live run
refuses only a directory (exit 1, `:869`) and then `os.Rename` silently replaces an occupying regular
file, reporting it as `created` with exit 0. A dry run that says "would refuse" followed by a live run
that overwrites is the opposite of what a dry run is for.

Beside it, `commands.go:891` reads `info, _ := os.Stat(targetPath)` and then `info.Size()`: a stat
failure after a successful rename, which a racing delete can produce, dereferences nil.

## Why

The guide's "Atomic download" section describes one occupancy rule; the command has two. The prose was
left alone by REQ-597 because the code, not the sentence, is what is wrong.

## Detailed Requirements

- One occupancy rule for both modes. Say in the record which one: the safe reading is that the live run
  refuses an existing regular file the way its dry run already does, so the two agree and nothing is
  silently replaced. If a caller depends on replacement, find it before choosing (grep every shipped
  prescribed block that runs `atomic-download`).
- The stat after the rename either handles its error or is removed with the size it reported; no
  discarded error before a dereference.
- Tests: dry-run and live against an occupying regular file give the same exit and the same finding;
  live against a directory unchanged; the stat-failure branch, if kept, does not panic.
- The guide's "Atomic download" sentence is re-derived from the new behaviour, measured.

## Constraints

- Shipped Go: a release. Keep the change to the occupancy check and the stat.

## Open Questions

None.
