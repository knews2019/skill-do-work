# REQ-256 builder brief

**Route A — direct implementation.** Estimated 5 active minutes.

Two doc sites disclose that the SessionStart hook now *writes* to consumer queue files (REQ-246's repairer plus the pre-existing reservation reaper):
1. `README.md` (~line 188, locate by content "SessionStart hook that injects the installed version and pending REQ count") — add the write behavior in one sentence if it fits.
2. `skills/do-work/actions/capture.md` — one line documenting `scripts/repair-req-timestamps.sh` (SessionStart hook) the way `cleanup-req-reservations.sh` is documented there.

Citation form: literal relative paths from the citing file's directory (REQ-249's rule; the reference contract enforces backticked cross-package citations mechanically). Note README.md sits at repo root — a path like `skills/do-work/scripts/...` from there. capture.md cites same-package paths locally (`scripts/repair-req-timestamps.sh`).

## How this build runs

You are a **worktree builder** for the do-work pipeline. Work only inside `/home/user/skill-do-work-worktrees/worktree-agent-REQ-256-disclose-the-session-hooks-queue-write-surface-in-the-docs` on branch `worktree-agent-REQ-256-disclose-the-session-hooks-queue-write-surface-in-the-docs`, cut from integration tip `7bb73c3` (0.212.17).

- Never write anything under `/home/user/skill-do-work` except your hand-back file (below). Never read/write `do-work/` in your worktree (stale snapshot; your REQ body is inlined below).
- Commit on your branch in small increments. Never touch `VERSION`, `CHANGELOG.md`, `skills/do-work/actions/version.md` (serial-only).
- One-line needs outside your write set are integration seams (hand back the exact line); larger needs stop-and-report. Out-of-scope finds go to `## Discovered Tasks`.
- Crew rules from your own worktree first: `skills/do-work/crew-members/general.md`, `coding-guardrails.md`, `communication-style.md`. Read every `prime_files` path.
- P-A-U phasing mandatory — work the block in your REQ body, record evidence in the hand-back. D-XX decisions with reasoning.

## Environment notes

- `bash _dev/tests/maintainer-verify.sh` exits 0 at your branch point — baseline and gate; exit code only, never piped.
- Toolchain: Go 1.26.1, ShellCheck 0.11.0, `just`, Node 22, Chromium. Never bare `go build` in the board tool dir.
- Clock: `date -u +%Y-%m-%dT%H:%M:%SZ` at the moment of use. Fixtures in scratch space only.

## Hand-back

Write to exactly this absolute path (the one main-tree write; never stage or commit it):

```
/home/user/skill-do-work/do-work/runs/work-2026-08-18-200845/REQ-256-handback.md
```

Structure: `# REQ-NNN hand-back` — **Branch**, **Commits**, `## What I built`, `## File manifest`, `## P-A-U evidence`, `## Testing evidence` (real RED/GREEN output; observed maintainer-verify exit code), `## Decisions (D-XX)`, `## Integration seams`, `## Discovered Tasks`, `## Pushback`.

---

# Your REQ (verbatim copy — the live one lives in the main tree)

---
id: REQ-256
title: Disclose the session hook's queue write surface in the docs
status: claimed
created_at: 2026-08-18T17:48:08Z
claimed_at: 2026-08-18T20:08:45Z
route: A
user_request: UR-056
addendum_to: REQ-246
domain: general
review_generated: true
effort_estimate: normal
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: false
write_set:
- README.md
- skills/do-work/actions/capture.md
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-08-18T20:08:45Z
  basis:
    - trivial short-circuit (effort_estimate would be trivial; gate-stamped normal for the disclosure class)
---

# Disclose the Session Hook's Queue Write Surface in the Docs

## What

REQ-246 made the SessionStart hook a *write* surface on consumer queue files — it mechanically repairs detectably wrong `*_at` stamps in `do-work/queue/` and `do-work/working/` at session start. Two shipped texts still describe the hook as read-only-plus-banner; a consumer auditing "what writes to my repo at session start" is misled.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Context

Instance 1 is REQ-246's review finding I3 (restatement sweep — the omission predates REQ-246 for the reservation cleanup, but a hook that edits files unattended is a different disclosure class than one that prints a banner). Instance 2 is the optional integration seam REQ-246's builder offered and deliberately did not write. Citations follow the literal cross-package path rule REQ-249 establishes.

## Instances

- [ ] **`README.md` (~line 188):** "SessionStart hook that injects the installed version and pending REQ count" — add that the hook also reaps stale REQ-number reservations and mechanically repairs detectably wrong queue/working timestamps (one clause each; keep it one sentence if it fits).
- [ ] **`skills/do-work/actions/capture.md`:** document `scripts/repair-req-timestamps.sh` (SessionStart hook) the way `cleanup-req-reservations.sh` is documented — one line stating that detectably wrong queue/working `*_at` stamps are mechanically corrected at session start.

## Requirements

- Both texts state the hook's write behavior; no behavior change anywhere.
- `bash _dev/tests/maintainer-verify.sh` exits 0.

---

## Triage

**Route: A** - Simple

**Reasoning:** Two named doc sites with the sentence content specified; no behaviour change.

**Planning:** Not required

## Plan

**Planning not required** - Route A

*Skipped by work action*
