---
id: REQ-501
title: '[impact-rule-change] Add do-work run-with-recovery and record the one-broken-pipe principle'
status: completed-with-issues
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
estimate:
  p50_active_minutes: 60
  confidence: low
  calculated_at: 2026-09-02T13:36:32Z
  basis:
    - Route C
    - 11-file write set
    - 2 new files
    - 3 subsystems involved
    - 8 acceptance criteria
    - dependency depth 2
    - cross-route regression gates
    - full-suite verification
write_set: [skills/do-work/actions/run-with-recovery.md, skills/do-work/SKILL.md, skills/do-work/actions/help.md, skills/do-work/crew-members/communication-style.md, skills/do-work/next-steps.md, skills/do-work/actions/work.md, skills/do-work/actions/work-reference.md, decisions/records/adr-022-one-broken-pipe-does-not-stop-the-factory.md, decisions/log.md, decisions/topics/_index_workflow-orchestration.md, CLAUDE.md, _dev/tests/contract-regressions.sh]
claimed_at: 2026-09-02T18:10:47Z
planning_at: 2026-09-02T18:16:50Z
dispatch_at: 2026-09-02T18:16:50Z
builder_handback_at: 2026-09-02T18:30:12Z
integration_at: 2026-09-02T18:31:02Z
review_at: 2026-09-02T18:43:42Z
remediation_at: 2026-09-02T19:18:00Z
re_review_at: 2026-09-02T19:25:35Z
completed_at: 2026-09-02T19:25:35Z
commit: 49a23f23
release_at: 2026-09-02T19:26:41Z
---

# Add do-work run-with-recovery and Record the One-Broken-Pipe Principle

## What
Add the `do-work run-with-recovery` verb: `run` under the user's assertion that this checkout is the queue's only writer, so every ownership refusal `run` makes is answered "mine", then hand off to `run` with all arguments passed through. Record the maintainer's running principle, "one broken pipe doesn't stop the rest of the factory from running", as ADR-022, as a paragraph in `work-reference.md`'s Execution Model, and as one bullet in `CLAUDE.md`.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Loaded every requested action, shell, recovery, decision, and agent-compatibility contract; designed a thin recovery wrapper with no new mutation mechanism.
- [x] **[APPLY]:** Added the action, routing/help/alias/next-step integration, exact sole-releaser recovery and crash semantics, ADR-022 and indexes, the one-broken-pipe principle, and mutation-tested regression lanes.
- [x] **[UNIFY]:** Reviewed all 12 files, checked the complete diff and debug artifacts, ran focused and canonical verification, committed the branch, and left the worktree clean.

## Why
Authority is already `run`'s default (`actions/work.md:124`, `actions/work-reference.md:65`). What `run` never does is answer an ownership question with "mine": Crash Recovery leaves foreign, unlabeled, or unnamed `working/` claims untouched and gates takeover on a human after three hours (`work-reference.md:341-365`), and REQ-498 stops on ambiguous shared metadata. After a context compaction the user is often certain that no other writer exists and wants one command that resumes in a fresh context without a prompt. The principle needs a citable home so future orchestration work is judged against it.

## Context
Names considered by the user: run with authority, run with recovery, run with authority and recovery, run-all-here. Chosen: `run-with-recovery`. Not folded into `run` because the assertion must be a deliberate invocation and must never leak into scripted `run` on a shared queue. Not a flag for the same reason and because `work.md` rejects unrecognized arguments.

## Detailed Requirements
- `skills/do-work/SKILL.md`: a new router row **above** the `run` row (first-match-wins; `run-simple-reqs` is the precedent). Triggers: `run-with-recovery`, `rwr`, `run-all-here`, `recover and run`, `run with authority`. `actions/help.md` gains one line.
- `skills/do-work/actions/run-with-recovery.md`, modeled on `run-simple-reqs.md` as a thin verb that completes `actions/work.md`: When to Use / Do NOT use (use `run` when other checkouts may be live; use `cleanup` for archive repair); Input is `work.md`'s input passed through verbatim, residue rejected.
  - Step 0.1: invoke `recover-finalization --discover --assume-sole-releaser` (REQ-499) through the canonical launcher with the no-fallback wording the launcher lane in `_dev/tests/contract-regressions.sh` checks.
  - Step 0.2: a tail whose implementation diff is uncommitted or whose release was never written resumes at `work.md` Step 9 in prose, the recovery result's exact lifecycle paths standing in for the `complete` result. Judgment stays prose; no new mechanism.
  - Step 0.3: authority crash recovery. Every REQ in `do-work/working/` classifies as this checkout's own crash; Crash Recovery substeps 1 to 3 run with no prompt and no three-hour ladder. Each takeover is reported with the label the checkpoint had.
  - Step 1: hand off to `do-work run $ARGUMENTS`.
- Recovery boundary, stated in the action file: a REQ found in `do-work/working/` is taken over and restarted from claim, with its uncommitted implementation preserved in the tree; pre-flight reports it and the builder may reuse it. Mid-step resumption of a build is out of scope for this verb; only the archive/release/commit tail resumes where it stopped (REQ-498, REQ-499).
  - One Common Rationalizations row: "the checkpoint says another writer, so I'll leave that one" → the user asserted sole authority by choosing this verb.
- `actions/work.md` Step 1 Crash Recovery: one sentence naming the verb. `actions/work-reference.md` Crash Recovery: the authority classification sentence. `actions/work-reference.md` Execution Model: the principle paragraph, keyed on the condition with examples marked illustrative: whenever a failure is local to one REQ, set that REQ aside with a typed finding and keep draining; only shared-target dirt stops the loop, and then the finding names the verb that resolves it.
- `skills/do-work/crew-members/communication-style.md` § Aliases: add the standalone token `rwr` — "Run with recovery: follow `actions/run-with-recovery.md`" — the same shape as the existing `phandoff` entry, so the shorthand works in an ordinary maintainer session and not only through the router.
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
**GREEN when:** All three lanes pass, mutating either predicate sentence trips the suite, `SKILL.md` routes `run-with-recovery` and `rwr` to `actions/run-with-recovery.md` before the `run` row, and the `rwr` alias is present in the communication-style Aliases list.
**Validation:** User confirmed the verb and delivery during capture; authority semantics inferred from the input's "full authority" wording and recommended during capture; recovery boundary confirmed during verify-requests.

## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` — 3436 tokens, over the 2000-token budget; `slugged: partial` so no targeted form. Matched on "changing action routing" and family `cross-action-exception-closure`.

## Implementation Summary

- Added `actions/run-with-recovery.md` as a thin, invocation-scoped sole-writer wrapper around canonical finalization recovery and ordinary `run`.
- Wired first-match routing, help, the `rwr` alias, and next-step guidance without changing plain `run` semantics.
- Recorded the one-broken-pipe principle in work-reference, CLAUDE, ADR-022, the decision log, and the workflow topic index.
- Added canonical-launcher, router-order, handoff, and mutation-tested recovery/continuation contract regressions.

## Decisions

- D-01: validate `work.md`'s argument grammar before recovery, then preserve `$ARGUMENTS` verbatim.
- D-02: treat successful recovery records as Step 9 continuation evidence; do not create a parallel archive/release mechanism.
- D-03: keep sole authority invocation-scoped through a separate verb, preserving plain `run` semantics.

## Testing

- RED: contract regressions failed on the missing action, route, and recovery predicates.
- GREEN: `bash _dev/tests/contract-regressions.sh` passed.
- `bash _dev/tests/maintainer-verify.sh` passed; the optional browser lane had the normal no-browser skip.
- `git diff --check` passed.

## Review — Attempt 1

**Overall: 50%** | **Acceptance: Fail** | **Risk: Critical**

All firm public surfaces are present and the static mutation lanes pass, but an isolated executable takeover fixture cannot reach a fresh claim. The action directly resets and moves a foreign working REQ through Crash Recovery substeps 1–3, leaves the queue/working rename dirty and the foreign checkpoint entry in place, then hands off to strict plain `run`:

- `recover-finalization --discover` refuses `FINALIZATION-DISCOVERY-AMBIGUOUS` on the dirty pair.
- `next` excludes the target as `ALREADY-CLAIMED` because the checkpoint residue remains.
- Direct canonical `claim` refuses `GIT-DIRTY-TARGET`.

The single remediation must establish one clean canonical ownership-transfer boundary, remove the exact prior checkpoint entry, preserve unrelated implementation bytes, and add a behavioral public recovery-to-selection-to-claim fixture. If that cannot close, the residual folds into pending REQ-504, which already owns moving Crash Recovery mechanics behind a canonical command.

## Remediation — Attempt 1

- Added a guarded canonical `recover-claim` request-state transaction requiring the invocation-scoped sole-writer assertion, exactly one checkpoint-evidence mode, and an exact-path lifecycle commit.
- Recovery now resets and moves only the asserted working REQ, removes only its checkpoint entry, preserves unrelated unstaged implementation bytes, and leaves plain `claim`/`run` dirty-target guards unchanged.
- Added executable recovery-to-`next`-to-fresh-`claim` coverage plus negative authority, evidence, commit, and rollback cases.
- Verification passed: focused race tests, full Go tests and vet, Go 1.25 compatibility, contract regressions, and the canonical maintainer gate. The optional browser lane skipped because no browser was available.

## Re-Review — Attempt 1

**Overall: 50%** | **Acceptance: Fail** | **Risk: Critical**

The guarded `recover-claim` transaction closes the primitive-level ownership-transfer defect, including exact-path commit/rollback and unrelated-dirt preservation, but the public action remains unable to guarantee a fresh claim:

- `recover-finalization --discover` runs first and refuses the normal interrupted-claim topology as `FINALIZATION-DISCOVERY-AMBIGUOUS`, so claim recovery is unreachable.
- A supported multiple-writer-label checkpoint leaves a second entry behind, and `next` then excludes the recovered REQ as `ALREADY-CLAIMED`.
- The prescribed shell command interpolates an observed writer label inside single quotes instead of transporting it structurally.

The single remediation allowance is exhausted. REQ-501 is therefore terminal `completed-with-issues`; all three residuals are folded into pending REQ-504, which already owns the canonical Crash Recovery/run-with-recovery boundary and behavior-level replacement tests.

## Open Questions
None.

## Full Context
See `do-work/user-requests/UR-097/input.md` for complete verbatim input.

---
*Source: capture of the run-with-recovery request (UR-097).*
