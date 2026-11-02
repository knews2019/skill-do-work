---
id: REQ-709
title: 'recover ignores loose do-work files that no request record can own, so do-work commit can commit inbox notes'
status: pending
created_at: 2026-11-02T19:41:34Z
user_request: UR-162
domain: backend
prime_files: []
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
---
# Recover Ignores Loose do-work Files That No Request Record Can Own
## What
The recovery check (`recover`, which `do-work commit` and the work loop run first) refuses with `FINALIZATION-DISCOVERY-AMBIGUOUS` when the only uncommitted file is a note under `do-work/` that no request record can own, such as an untracked `do-work/inbox/<note>.md`. It should leave such a file alone, so `do-work commit` can commit it.
## Why
`do-work commit` stops on the refusal, so it can never commit an inbox note, which is exactly the leftover file it exists for.
## Detailed Requirements
- A file under `do-work/` outside the folders that hold request records (the queue, working, archive and user-request folders, REQ reservation markers, and run records) is not finalization evidence. Recovery neither refuses on it nor associates it with a request.
- Key the rule on that condition, stated once in code: where request records live. Do not list note folders such as `inbox/`. The user chose this scope so that every note (inbox notes, handoff notes such as `do-work/HANDDOWN-*.md`) is covered.
- Real unfinished finalization work next to such a note still refuses or settles exactly as today.
- `do-work commit` then groups the note as an ordinary change with no owning request.
## Context
Seen twice on 2026-11-02. During the UR-161 run, `recover` refused on the untracked note `do-work/inbox/2026-10-28_do-work-upstream-suggestion-month-folder-report-layer.md`; the run agent hid it with a temporary `.git/info/exclude` line. Then `do-work commit` refused on the same note and committed it by judgment (commit `14878ab6`). The refusal is raised in `skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go`, where the note lands in the unresolved shared-state list. No queued candidate in any UR shares this root cause (the queue was empty at capture).
## Red-Green Proof
**RED prompt/case:** A test repository whose only uncommitted file is an untracked `do-work/inbox/note.md`. Run `recover`, then the `do-work commit` preflight.
**Why RED now:** `recover` refuses with `FINALIZATION-DISCOVERY-AMBIGUOUS` and names the note, so `do-work commit` stops.
**GREEN when:** `recover` returns success, and `do-work commit` commits the note as a group with no owning request. A real unfinished request next to the note is still handled exactly as today.
**Validation:** User confirmed.
## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` as a whole satellite (23234 tokens per `do-work/lessons-index.md`, over the 2000 budget; `slugged: partial`, so no targeted form). Matching reason: its row covers semantic recovery completeness and condition-based classifiers; families `closed-enumeration-for-a-condition`, `observed-subset-is-not-semantic-completeness` match this work.
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: "recover should not treat untracked do-work/inbox notes as unfinished finalization work, so do-work commit can commit them"*
