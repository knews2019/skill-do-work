---
id: REQ-723
title: '[impact-negligible] Archive-collision helper skips assets folders'
status: claimed
created_at: 2026-11-14T22:31:56Z
user_request: UR-174
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
maintenance: false
impact: impact-negligible
effort_estimate: effort-mechanical
route: A
estimate:
  p50_active_minutes: 10
  confidence: high
  calculated_at: 2026-11-14T23:29:33Z
  basis:
    - Route A
    - 2-file write set
    - 3 acceptance criteria
write_set: [skills/do-work/tools/do-work-cli/internal/corehelpers/commands.go, skills/do-work/tools/do-work-cli/internal/corehelpers/archive_collision_test.go]
claimed_at: 2026-11-14T23:28:11Z
dispatch_at: 2026-11-14T23:32:28Z
---
# Archive-Collision Helper Skips Assets Folders
## What
The standalone `archive-collision REQ-NNN` command (`handleArchiveCollision`, `skills/do-work/tools/do-work-cli/internal/corehelpers/commands.go:158-193`) walks `do-work/archive` and lists every file named `REQ-NNN.md` or `REQ-NNN-*.md`, including evidence files inside `assets/` folders. Skip folders named `assets` in that walk, the same rule REQ-722 (Cleanup Pass 6 skips assets folders when scanning for blanked records) added to cleanup and `repositorymodel` already uses.
## Context
- Source: review finding F1 of REQ-722, recorded there as report only and promoted by the user with this capture. Original verdict: impact-negligible.
- Evidence: the walk at `commands.go:164-181` filters on the file name only and never returns `fs.SkipDir`. The reference rule is `skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model.go:192` (`directoryEntry.Name() == "assets"`) and the REQ-722 fix in `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_git.go:115`.
- Reach: the command is launched by `skills/do-work/tools/checks/archive-collision.sh`. No action prescribes it. The `advance` archive-collision hold reads the `repositorymodel` snapshot and is not affected. That is why impact is negligible.
- No queued REQ in any UR shares this root cause (the queue was empty at capture).
- Existing tests live in `skills/do-work/tools/do-work-cli/internal/corehelpers/archive_collision_test.go`.
## Red-Green Proof
**RED prompt/case:** A repo whose archive holds only `do-work/archive/UR-001/assets/REQ-001-evidence.md`. Run `do-work-cli archive-collision REQ-001`.
**Why RED now:** The command prints that evidence path and reports findings, as if REQ-001 had an archived twin.
**GREEN when:** A new Go test in `archive_collision_test.go` runs the command on that fixture and gets success with empty output, and a real archived record in the same UR folder (`do-work/archive/UR-001/REQ-001-done.md`) is still listed when present. The two existing archive-collision tests still pass.
**Validation:** Inferred during capture
## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 24500 tokens, over the 2000-token budget, and its row is `slugged: partial` so no targeted entry is legal. Matched (claim-time consult, 2026-11-15) on the row clause "a record walker of its own that must copy repositorymodel's folder skips" and family `guard-on-one-entry-path`, where REQ-722 recorded the sibling fix. The pre-dispatch agent read that one bullet (line 5) by hand; the builder brief quotes it.
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: do-work capture-request: the standalone archive-collision helper (internal/corehelpers/commands.go:164-176) does not skip assets/ folders, the same gap REQ-722 fixed in cleanup Pass 6. Source: REQ-722 review finding F1.*

## Triage

**Route: A** - Simple

**Reasoning:** The REQ names both files, the exact walk (`commands.go:164-181`) and the reference rule to copy (`repository_model.go:192`, `cleanup_git.go:115`); the change is one `fs.SkipDir` branch plus one Go test.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Decisions
- **D-01**: Copy the skip rule into the walk (return `fs.SkipDir` for a directory named `assets`), not switch the command to the `repositorymodel` snapshot. Reasoning: the REQ asks for the skip rule; it is one branch and keeps the command's output format and its symlink and missing-archive behavior unchanged. Value: smallest diff, same rule as `cleanup_git.go:115`. Risk: low; a later folder-skip added to `repositorymodel` would again need copying here (the lesson family `guard-on-one-entry-path` already names that drift), reversible.
- **D-02**: Skip only `assets`, not the other `repositorymodel` skips (top-level `deliverables`, dot-folders). Reasoning: the walk starts at `do-work/archive`, so `deliverables` can never appear as a top section there, and the REQ scopes the change to `assets`. Value: no behavior change beyond the REQ. Risk: a dot-folder inside the archive would still be scanned, as today; reversible.
