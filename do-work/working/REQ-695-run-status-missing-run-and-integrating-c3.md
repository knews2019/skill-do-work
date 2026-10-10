---
id: REQ-695
title: 'Addendum: run-status names a missing --run folder and stops recommending do-work run while an integrator works'
status: claimed
created_at: 2026-10-10T19:19:53Z
user_request: UR-154
addendum_to: REQ-690
domain: backend
prime_files: ["_dev/primes/prime-releases.md", "skills/do-work/tools/do-work-cli/prime-do-work-cli.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
related: [REQ-693, REQ-694, REQ-696]
batch: review-followups-ur154
required_lessons: ["_dev/primes/lessons-releases.md"]
write_set: ["skills/do-work/tools/do-work-cli/internal/runstatus/run_status.go", "skills/do-work/tools/do-work-cli/internal/runstatus/run_status_render.go", "skills/do-work/tools/do-work-cli/internal/runstatus/run_status_test.go"]
claimed_at: 2026-10-10T19:22:01Z
---
# Addendum: run-status Names a Missing --run Folder and Stops Recommending do-work run While an Integrator Works

## What

Two gaps the REQ-690 (do-work status action and run-status report) review left as report only. A mistyped `--run` silently reads as "no run", so hand-backs count as absent and C3 (hand-back landed) rows become C7 (claim past threshold) rows whose text names `run-with-recovery`, and the text never says which run folder it read (review F3). C3 also recommends `do-work run` while an integrator is working on that REQ right now (review F4), which invites a second orchestrator.

## Prior Implementation

REQ-690 shipped `do-work-cli run-status` in `skills/do-work/tools/do-work-cli/internal/runstatus/`, archived at commit `65820f45`. Decision D-19 (coordinator ruling) removed the `--run` stat and refusal, and D-20 declined the reviewer's optional display line. `classifyRow` sends every claimed REQ with a landed hand-back to C3 with `next_argv ["do-work","run"]`; the row already carries `last_activity_phase` and minutes since last activity.

## Detailed Requirements

- F3, display only (D-19 stays: no refusal, no exit-code change): when the `--run` folder does not exist, the report says so. JSON `run_directory` gains a "not found; run fields absent" note (or an equivalent field), and the text output prints a `run: <dir>` line. The reviewer's exact suggestion is in `REQ-690-review.md` under F3.
- F4: a C3 row whose last activity is in phase `integration` and recent (inside the existing quiet boundary) does not recommend `do-work run`. It says an integrator is working on it and carries no runnable `next_argv`. Other C3 rows are unchanged.

## Red-Green Proof

**RED prompt/case:** (1) `run-status --run do-work/runs/typo` on a fixture with a claimed REQ, a landed hand-back in the real run folder, and a claim older than the threshold. (2) A fixture C3 row with last activity in phase `integration` 5 minutes ago.
**Why RED now:** (1) the row is C7 and the output never names the folder it read. (2) the row says "run `do-work run`".
**GREEN when:** (1) the text and JSON say the run folder was not found. (2) the row says an integrator is working on it and its `next_argv` is not `do-work run`.
**Validation:** Inferred during capture

## Constraints

- Keep the D-19 ruling: a missing `--run` folder is reported, never refused.
- One test per failure in `run_status_test.go`.
- `next_argv` stays read-only or empty; never a destructive command.
- No `depends_on` edges in this batch; no other REQ touches these files.

## Builder Guidance

Certainty: high (F3 reproduced on the live repo, F4 seen live on REQ-690 itself; both re-checked at HEAD `0e8e0ef9`). Latitude: the builder picks whether F4 is a new class or a C3 variant; prefer the smaller change.

## Required Lessons — Dropped for Budget

- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 20802 tokens, bare only (`slugged: partial`); matches the runstatus path (families `destructive-next-argv`, `silent-skip-reads-as-red`).

## Full Context

See `do-work/user-requests/UR-154/input.md` for complete verbatim input. Sources: `do-work/archive/UR-153/REQ-690-run-status-action.md` → `## Review` and D-19, D-20; `do-work/runs/work-2026-10-10-131527/REQ-690-review.md` (F3, F4); `do-work/runs/work-2026-10-10-131527/REQ-690-integration-report.md`.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: UR-154 R4 — "a mistyped `--run` silently turns C3 (hand-back landed) rows into C7 (claim past threshold) rows (review F3); C3 recommends `do-work run` even while an integrator runs (F4)."*
