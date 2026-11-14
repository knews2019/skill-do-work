---
id: REQ-722
title: 'Cleanup Pass 6 skips assets folders when scanning for blanked records'
status: pending
created_at: 2026-11-14T21:47:51Z
user_request: UR-173
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
write_set: [skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_git.go, skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_git_test.go]
---
# Cleanup Pass 6 Skips Assets Folders When Scanning for Blanked Records
## What
Cleanup Pass 6 (`AddBlankedRepairs`, `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_git.go:117`) treats every `REQ-*.md` or `UR-*.md` file under `do-work/archive`, `do-work/queue` and `do-work/working` as a request record, including evidence files inside `assets/` folders. Such a file has no frontmatter, so the pass reports `BLANKED-RECORD-UNRECOVERABLE` and `cleanup` exits 1. Skip `assets/` folders in the Pass 6 walk and add one test for it.
## Why
A consumer repo got 49 false warnings and exit 1 after archiving two user requests. The cleanup action stops on a nonzero exit, so the false warnings block a clean cleanup.
## Context
- Finding 1 of a validate-feedback triage, source: pasted text (consumer feedback). Original severity: none given. Verdict: Accept. Surface-cost: N/A (direct fix, no new guard).
- Evidence: the walk at `cleanup_git.go:105-119` filters on the file name only. `isBlankedDocument` (`cleanup_git.go:159`) returns true when `requestmodel.ParseDocument` fails, which it does for any file without frontmatter. `repositorymodel` already skips these folders (`internal/repositorymodel/repository_model.go:192`, `directoryEntry.Name() == "assets"` returns `fs.SkipDir`).
- No queued REQ in any UR shares this root cause (the queue was empty at capture).
- Check whether `tools/checks/blanked-req-scan.sh` (named by `actions/cleanup.md` Pass 6) has the same gap. If it does, the builder decides whether it belongs in this fix.
## Red-Green Proof
**RED prompt/case:** A repo whose only do-work file is `do-work/archive/UR-001/assets/REQ-001-evidence.md` (plain text, no frontmatter), committed. Run `do-work-cli --format json cleanup --dry-run`.
**Why RED now:** The output holds one `BLANKED-RECORD-UNRECOVERABLE` finding for that path and the exit code is 1. Reproduced during triage on 2026-11-14.
**GREEN when:** A new Go test in `internal/cleanup/cleanup_git_test.go` calls `AddBlankedRepairs` on that fixture and gets no finding and no group, and the existing test `TestBlankedRestorationRequiresExactConsentAndUsesGitHistory` still passes (a truly blanked record outside `assets/` is still reported).
**Validation:** Inferred during capture
## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 24277 tokens, over the budget, and its row is `slugged: partial` so no targeted entry is legal. Matched because the satellite's prime governs `internal/cleanup` (families `closed-enumeration-for-a-condition`, `guard-on-one-entry-path`).
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: do-work capture-request: cleanup Pass 6 (AddBlankedRepairs, cleanup_git.go:117) flags assets/REQ-\*.md evidence files as BLANKED-RECORD-UNRECOVERABLE and exits 1.*
