---
id: REQ-501
title: '[impact-rule-change] Add do-work run-with-recovery and record the one-broken-pipe principle'
status: pending
created_at: 2026-09-02T13:31:12Z
user_request: UR-097
domain: general
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec:
depends_on: [REQ-499]
related: [REQ-499, REQ-500]
batch: run-with-recovery
maintenance: false
impact: impact-rule-change
effort_estimate: effort-substantive
write_set: [skills/do-work/actions/run-with-recovery.md, skills/do-work/SKILL.md, skills/do-work/actions/help.md, skills/do-work/next-steps.md, skills/do-work/actions/work.md, skills/do-work/actions/work-reference.md, decisions/records/adr-022-one-broken-pipe-does-not-stop-the-factory.md, decisions/log.md, decisions/topics/_index_workflow-orchestration.md, CLAUDE.md, _dev/tests/contract-regressions.sh]
---

# Add do-work run-with-recovery and Record the One-Broken-Pipe Principle

## What
Add the `do-work run-with-recovery` verb: `run` under the user's assertion that this checkout is the queue's only writer, so every ownership refusal `run` makes is answered "mine", then hand off to `run` with all arguments passed through. Record the maintainer's running principle, "one broken pipe doesn't stop the rest of the factory from running", as ADR-022, as a paragraph in `work-reference.md`'s Execution Model, and as one bullet in `CLAUDE.md`.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why
Authority is already `run`'s default (`actions/work.md:124`, `actions/work-reference.md:65`). What `run` never does is answer an ownership question with "mine": Crash Recovery leaves foreign, unlabeled, or unnamed `working/` claims untouched and gates takeover on a human after three hours (`work-reference.md:341-365`), and REQ-498 stops on ambiguous shared metadata. After a context compaction the user is often certain that no other writer exists and wants one command that resumes in a fresh context without a prompt. The principle needs a citable home so future orchestration work is judged against it.

## Context
Names considered by the user: run with authority, run with recovery, run with authority and recovery, run-all-here. Chosen: `run-with-recovery`. Not folded into `run` because the assertion must be a deliberate invocation and must never leak into scripted `run` on a shared queue. Not a flag for the same reason and because `work.md` rejects unrecognized arguments.

## Detailed Requirements
- `skills/do-work/SKILL.md`: a new router row **above** the `run` row (first-match-wins; `run-simple-reqs` is the precedent). Triggers: `run-with-recovery`, `recover and run`, `run with authority`. `actions/help.md` gains one line.
- `skills/do-work/actions/run-with-recovery.md`, modeled on `run-simple-reqs.md` as a thin verb that completes `actions/work.md`: When to Use / Do NOT use (use `run` when other checkouts may be live; use `cleanup` for archive repair); Input is `work.md`'s input passed through verbatim, residue rejected.
  - Step 0.1: invoke `recover-finalization --discover --assume-sole-releaser` (REQ-499) through the canonical launcher with the no-fallback wording the launcher lane in `_dev/tests/contract-regressions.sh` checks.
  - Step 0.2: a tail whose implementation diff is uncommitted or whose release was never written resumes at `work.md` Step 9 in prose, the recovery result's exact lifecycle paths standing in for the `complete` result. Judgment stays prose; no new mechanism.
  - Step 0.3: authority crash recovery. Every REQ in `do-work/working/` classifies as this checkout's own crash; Crash Recovery substeps 1 to 3 run with no prompt and no three-hour ladder. Each takeover is reported with the label the checkpoint had.
  - Step 1: hand off to `do-work run $ARGUMENTS`.
  - One Common Rationalizations row: "the checkpoint says another writer, so I'll leave that one" → the user asserted sole authority by choosing this verb.
- `actions/work.md` Step 1 Crash Recovery: one sentence naming the verb. `actions/work-reference.md` Crash Recovery: the authority classification sentence. `actions/work-reference.md` Execution Model: the principle paragraph, keyed on the condition with examples marked illustrative: whenever a failure is local to one REQ, set that REQ aside with a typed finding and keep draining; only shared-target dirt stops the loop, and then the finding names the verb that resolves it.
- `skills/do-work/next-steps.md`: one row, `run` stopped on REQ-498's ambiguous-shared-state code → suggest `run-with-recovery`.
- `decisions/records/adr-022-one-broken-pipe-does-not-stop-the-factory.md`, frontmatter like ADR-021, `related: adr-018 (extends)`: motivation (the REQ-494 incident), the principle, why a verb and not a flag, the shared-target caveat (a dirty checkpoint is a target every claim writes, so it cannot be set aside, only resolved or resolved by assertion), the declined alternative (letting `claim` tolerate a dirty checkpoint weakens a guard that also protects against claiming into a half-resolved merge), and instances: REQ-468 to REQ-472, REQ-491/492, REQ-498, REQ-499, REQ-500, this REQ. Index lines in `decisions/log.md` and `decisions/topics/_index_workflow-orchestration.md`.
- `CLAUDE.md` § A Note From Me, one bullet: **One broken pipe doesn't stop the rest of the factory.** A failed or interrupted REQ is set aside with a typed finding and the loop continues; only shared-target dirt may stop it, and then the finding names the verb that resolves it.
- `_dev/tests/contract-regressions.sh` lanes: the canonical-launcher lane for the new action; a router row-order check (`run-with-recovery` row index below `run`'s); a mutation-tested predicate over the Crash Recovery and Execution Model sentences so deleting "every `working/` REQ classifies as own crash" or the set-aside sentence trips the suite.

## Constraints
- Preserve the single-releaser model; the assertion is a deliberate verb invocation, never a flag that leaks into scripted `run`.
- Never widen recovery to secret-classified or project paths; the action inherits REQ-499's limits.
- Action files stay agent-agnostic and work as a standalone prompt (CLAUDE.md § Agent Compatibility).

## Dependencies
Depends on REQ-499 (the flag the action invokes). Does not depend on REQ-500.

## Builder Guidance
Firm on the verb name and the three homes for the principle. Latitude on exact prose. Read `_dev/primes/prime-action-files.md` before writing the action file.

## Red-Green Proof
**RED prompt/case:** Run `bash _dev/tests/contract-regressions.sh` after adding the three new lanes; separately, feed the sentence "do-work run-with-recovery" to the router table in `SKILL.md`.
**Why RED now:** The lanes fail because the action file, router row, ADR, and predicate sentences do not exist; the router falls through to the `run` row or the unknown-single-word branch.
**GREEN when:** All three lanes pass, mutating either predicate sentence trips the suite, and `SKILL.md` routes `run-with-recovery` to `actions/run-with-recovery.md` before the `run` row.
**Validation:** User confirmed (verb chosen during capture).

## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` — 3436 tokens, over the 2000-token budget; `slugged: partial` so no targeted form. Matched on "changing action routing" and family `cross-action-exception-closure`.

## Open Questions
None.

## Full Context
See `do-work/user-requests/UR-097/input.md` for complete verbatim input.

---
*Source: capture of the run-with-recovery request (UR-097).*
