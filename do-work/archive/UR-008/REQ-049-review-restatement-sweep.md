---
id: REQ-049
title: "Bake a restatement-sweep step into the review instructions: grep every restatement of a changed contract token before verdict"
status: completed
route: B
created_at: 2026-07-29T09:30:45Z
claimed_at: 2026-07-29T13:36:51Z
completed_at: 2026-07-29T13:44:00Z
commit: 6164a10
user_request: UR-008
domain: general
prime_files: []
tdd: false
depends_on: []
related: [REQ-043, REQ-044, REQ-045]
batch: deep-review-followups
write_set: [actions/work.md, actions/review-work.md]
maintenance: false
---

# Review Restatement Sweep (Calibration Fix from the REQ-035–040 Deep Review)

## What

The REQ-035–040 batch passed adversarial review at 86–98%, yet an independent deep pass found that **every top confirmed defect was the same class: a contract token changed in its canonical home while a restatement or consumer elsewhere kept the old semantics** (the proceed-anyway gate phrasing vs the rewritten Crash Recovery gate; the commit-procedure `HASH` block vs the worktree `commit:` paragraph; `git show <commit>` consumers vs the new merge-hash meaning of `commit:`; the "each REQ firms at Step 5.5" claim vs Route A's skip). The reviews tracked requirement coverage well but had no step forcing a repo-wide sweep of restatements. Bake that step in structurally — the durable-behavior route — rather than relying on reviewer judgment.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Canonical home = `actions/review-work.md` Step 6 (Code Review), as a new bolded check **Restatement Sweep** placed immediately after **Risk Assessment** — Risk Assessment already covers *code* consumers of a changed interface ("identify callers/dependents of changed code, flag interfaces whose contract changed"), so the sweep is its documented-consumer twin and reads as continuous with the surrounding dimensions. Both modes reach Step 6, so standalone `do-work review` inherits it for free. Write the trigger as a question about the diff ("does anything here change the meaning of something stated in more than one place?") with illustrative examples only, plus an explicit skip clause for diffs that redefine nothing. Add one Common Rationalizations row (the "it's in a file the REQ didn't declare, so it's out of scope" dismissal, which is the exact mechanism by which the four REQ-035–040 defects were skippable) and one Verification Checklist item so the step is required, not optional. Cross-reference from `actions/work.md` Step 7 as a single **MUST** sentence pointing at the review-work.md home. Restate the prescribed-command grep rule inline (generalized from commands to contracts) — never cite the maintainer doc, which is export-ignored.
- [x] **[APPLY]:** Two files touched, exactly as planned: `actions/review-work.md` (Step 6 sweep block + 1 rationalization row + 1 checklist item), `actions/work.md` (1-sentence MUST cross-reference in Step 7).
- [x] **[UNIFY]:** `git diff --stat` = 2 files, 15 insertions / 0 deletions (`actions/review-work.md` +13, `actions/work.md` +2) — both in `write_set`, nothing else touched; no `SKILL.md`, no frontmatter `status` edits, no git operations. Verified `actions/review-work.md`: sweep block sits inside Step 6 between Risk Assessment and Directive Alignment Check, trigger is condition-based (no hard-coded token list), proportionality skip clause present, severity left to reviewer judgment, new rationalization row carries do-work nouns (REQ, Scope, follow-up) so the Common-Rationalizations ratchet stays satisfied. Verified `actions/work.md`: single added line in Step 7, points at `actions/review-work.md` Step 6 by path. `bash _dev/tests/contract-regressions.sh` → `✓ Contract regression checks passed` (exit 0), which includes the shipped-files-must-not-cite-CLAUDE.md/AGENTS.md grep; independently confirmed with `grep -n 'CLAUDE\.md\|AGENTS\.md'` over both changed files → no hits. No debug artifacts, no TODO/FIXME, no stray files (`git status --porcelain -uall` shows only the two tracked edits plus the untracked-by-design `do-work/`).

## Requirements

1. Add to the review instructions (the Step 7 review dispatch in `actions/work.md` and/or the shared review workflow in `actions/review-work.md` — builder decides the canonical home and cross-references the other) a required sweep: for each contract token, schema field, phrasing, or command primitive the diff **changes the meaning of**, grep the repo for every other statement/consumer of it and verify each still agrees; an un-updated restatement is a finding, not out of scope.
2. State the trigger as a condition, not an enumerated token list (Closed Enumerations Go Stale) — illustrative examples are fine (a gate phrasing, a frontmatter field's semantics, what a hash field holds, a prescribed command's output shape).
3. Keep it proportionate: the sweep applies when the diff redefines something other text restates — not a blanket "grep everything" tax on every one-line fix. Align with the existing CLAUDE.md-origin rule about prescribed-command primitives ("grep the same primitive across all actions before calling it fixed") — this generalizes that rule from commands to contracts, in a shipped home rather than the export-ignored maintainer doc.

## Constraints

- Shipped files must not cite this repo's CLAUDE.md — restate the rule inline in the shipped home (the contract-regressions suite greps for dangling citation idioms).
- Watch SKILL.md's word budget is untouched (this change lives in action files, not SKILL.md).

## Red-Green Proof
**RED prompt/case:** Read `actions/work.md` Step 7's review instructions and `actions/review-work.md` — no step requires sweeping restatements of changed tokens; the four missed findings above are the evidence it was skippable.
**Why RED now:** Review prompts enumerate coverage/severity duties but not cross-file restatement verification.
**GREEN when:** The sweep step exists in the review instructions with a condition-based trigger and proportionality note; a reviewer following the instructions on REQ-035's diff would have been forced to grep the gate phrasing and find the proceed-anyway restatement.
**Validation:** User confirmed (approved capture of the calibration lesson as a structural fix)

## Full Context
See `do-work/user-requests/UR-008/input.md`.

## Triage

**Route:** B (Explore then Build)
**Reasoning:** Clear outcome (a restatement-sweep step must exist in the review instructions), but the builder must decide the canonical home (work.md Step 7 dispatch vs review-work.md shared workflow), place the step so both pipeline and standalone review inherit it, cross-reference the other home, and align with the existing prescribed-command rule — an exploration/placement pass, not a blind insertion. Route B.
**Complexity indicators:** 3 requirements; a placement decision (canonical home + cross-ref); condition-not-enumeration framing (Closed Enumerations); a proportionality guard; a hard constraint (no CLAUDE.md citation in shipped files — restate inline, the contract-regressions suite greps for dangling citations).
**Rigor:** Standard independent review (main-context) + confirm the GREEN counterfactual (a reviewer following the new step on REQ-035's diff would have been forced to grep the gate phrasing) and that no CLAUDE.md citation was introduced.

*Triaged 2026-07-29 by orchestrator (session do-work-20260729T100657Z-34626).*

## Exploration

**Where the review dimensions live.** `actions/review-work.md` is the shared workflow for *both* modes — its "Two Modes" table says explicitly that pipeline and standalone "follow the same workflow. The only difference is where the REQ lives and how you obtain the diff." The evaluative dimensions are all inside **Step 6: Code Review**, as a series of bolded blocks: Code Quality, Test Adequacy, Scope Discipline, Risk Assessment, Directive Alignment Check, Coding-Guardrails Principle Check, Domain-Specific Review. Step 5 is the requirements checklist; Step 7 is acceptance testing; the Scoring Guidelines (which average Requirements / Code Quality / Test Adequacy / Scope) sit after Step 8.

**The natural insertion point.** Risk Assessment's last bullet already reads "Regression risk — identify callers/dependents of changed code, flag interfaces whose contract changed, note shared utilities that other features rely on." That is the *code*-consumer version of exactly this sweep. Placing the sweep immediately after Risk Assessment (and before Directive Alignment Check) makes it read as the documented-consumer twin of a check the file already performs, rather than a bolted-on section. It also deliberately avoids the four averaged score dimensions — the sweep produces findings with judgment-set severity, so the Scoring Guidelines formula is untouched.

**Where work.md Step 7 dispatches.** `actions/work.md:499` "Step 7: Review" spawns `actions/review-work.md` in pipeline mode; the step's body is the dispatch plus the acceptance/score gate. The precedent for a required review sub-check being *stated* in work.md is line 356 (Step 5.5), which says the review step "**MUST** run the scope-drift comparison." I followed that idiom but kept the body in review-work.md — one MUST sentence with the trigger condition and a path pointer, no duplication of the four-part step.

**Why not work.md as the canonical home.** Standalone `do-work review` never reads work.md, so a work.md-only home would leave every manual review without the sweep — and the standalone mode is precisely where a reviewer is auditing an already-committed diff, the case with the most restatement drift to find.

**Counterfactual verified against real history.** At REQ-035's commit (`fd56267`) the diff deleted the "session other than this one" clause from the Crash Recovery per-file claim gate, making freshness alone the exemption. Simulating step 2 of the new sweep — `git grep -inE "concurrency gate skips|gate skips (only|any|every)|claimed by another (live )?session" fd56267 -- actions/` — returns exactly one hit, `actions/work-reference.md:375`: "the per-file concurrency gate skips **only files actively claimed by another live session**." That is the stale proceed-anyway restatement, surfaced by a single grep, and it was not fixed until REQ-044/045.

## Scope

**Files I will touch:**
- `actions/review-work.md` — canonical home: the **Restatement Sweep** block in Step 6, one Common Rationalizations row, one Verification Checklist item.
- `actions/work.md` — one-sentence **MUST** cross-reference in Step 7 pointing at the review-work.md home.

**Acceptance criteria (restated from Requirements):**
1. The review instructions require a sweep: for each contract token / schema field / phrasing / command primitive whose meaning the diff changes, grep the repo for every other statement or consumer and verify agreement; an un-updated restatement is a finding (severity by judgment), not out of scope.
2. The trigger is stated as a condition, with examples marked illustrative — no enumerated token list.
3. The sweep is proportionate: it fires only when the diff redefines something other text restates, with an explicit skip clause; it generalizes the prescribed-command grep rule from commands to contracts, restated inline in the shipped file.
4. No citation of this repo's maintainer doc (`CLAUDE.md`/`AGENTS.md`) in shipped files; `SKILL.md` untouched; `bash _dev/tests/contract-regressions.sh` green.

## Implementation Summary

**Files changed:**
- `actions/review-work.md` (modified)
- `actions/work.md` (modified)

**What was done:** Added a **Restatement Sweep** check to `actions/review-work.md` Step 6 (Code Review), immediately after **Risk Assessment** and before **Directive Alignment Check** — the shared workflow both pipeline mode and standalone `do-work review` run, so both inherit it. The block is four numbered parts plus a one-line origin note: (1) the trigger is a *question asked of the diff* — "does this change the meaning of something that is stated in more than one place?" — with a contract token, a schema field's semantics, a gate's wording, what a stored value holds, and a prescribed command's consumed output shape given explicitly as illustrative examples, and a stated reason not to keep a token list ("a hand-maintained token list goes stale the moment the contract grows"); (2) sweep each redefined element by grepping every other statement or consumer — the token, the phrasings that gloss it, and the tests/tooling/templates that parse it — restating the prescribed-command grep rule inline and generalizing it from commands to contracts; (3) every stale restatement is a finding with severity by judgment (Important when the stale text would make a reader or agent act on the old contract, Minor when cosmetic), explicitly including stale text in files the REQ never declared, routed to a follow-up REQ via Step 10 and never scored as the builder's scope drift; (4) an explicit skip clause — no sweep when nothing was redefined, "the trigger is redefinition, not diff size." Deliberately not a new averaged score dimension, so the Scoring Guidelines formula is unchanged. Reinforced by one Common Rationalizations row (the "that stale restatement is in a file this REQ never declared — out of scope" dismissal, the exact mechanism that made the REQ-035–040 defects skippable) and one Verification Checklist item, which is what makes the sweep required rather than optional. `actions/work.md` Step 7 gained a single **MUST** sentence — the same idiom Step 5.5 uses for the scope-drift comparison — naming the trigger condition and pointing at `actions/review-work.md` Step 6 by path, with no duplication of the four-part body.

## Decisions

**D-01: Canonical home is `actions/review-work.md` Step 6, not `actions/work.md` Step 7.** Review-work.md's own "Two Modes" table states that pipeline and standalone follow the same workflow, so a step defined there is inherited by both; a work.md-only home would leave standalone `do-work review` — the mode most likely to be auditing an already-committed diff for exactly this kind of drift — without the sweep. Work.md gets a one-sentence MUST pointer instead, matching the idiom Step 5.5 already uses to mandate the scope-drift comparison from outside review-work.md.

**D-02: Placed after Risk Assessment, not as a new top-level Step.** Risk Assessment's regression-risk bullet already sweeps *code* callers/dependents of a changed interface; the restatement sweep is its documented-consumer twin, and siting it there makes it continuous with a check the file already performs rather than a bolted-on section. A new `### Step 6.5` would also have forced renumbering pressure on a file other actions reference by step number (`actions/work.md` cites review-work.md Step 4 and Step 10).

**D-03: Findings-only, not a fifth averaged score dimension.** The Scoring Guidelines average Requirements / Code Quality / Test Adequacy / Scope; adding a percentage dimension would have changed every existing score's arithmetic and the `## Review` table shape that `actions/present-work.md` parses. Severity-by-judgment findings deliver the requirement ("a finding, not out of scope") with zero blast radius on the scoring contract.

**D-04: The origin note names the REQ-035–040 batch but no maintainer-doc path.** Shipped files must not cite this repo's export-ignored maintainer doc, so the prescribed-command grep rule is restated inline in review-work.md rather than pointed at. The origin sentence follows existing shipped precedent (`actions/work-reference.md` carries an "Origin incident" note and a "Design note (REQ-018 remediation)"), which gives the check the traceable failure mode the earned-section test requires.

## Review

**Acceptance: Pass — overall ~96%.** Main-context review against the full 2-file diff + the GREEN counterfactual.

**Requirements (all 3 met):**
1. **Restatement Sweep** added canonically to `review-work.md` Step 6, right after Risk Assessment (its documented-consumer twin), inherited by both pipeline and standalone modes; `work.md` Step 7 cross-references it with a MUST sentence. An un-updated restatement is a finding — including in files outside the REQ's Scope → routed to a follow-up, not scored as the builder's scope drift.
2. Trigger is condition-based ("does this change the meaning of something stated in more than one place?"), illustrative examples only, explicitly NOT a token list — honors Closed Enumerations Go Stale.
3. Proportionality guard present (skip when nothing is redefined; the trigger is redefinition, not diff size); generalizes the prescribed-command grep rule from commands to contracts, RESTATED INLINE.

**Constraints honored:** no CLAUDE.md/AGENTS.md citation (verified by grep + contract-regressions' citation check); SKILL.md untouched. Added one Common Rationalizations row (do-work nouns, ratchet-safe) + one Verification Checklist item.
**GREEN counterfactual confirmed:** a reviewer applying the step to REQ-035's diff (which changed the Crash Recovery gate semantics) is forced to grep the gate phrasing → finds the proceed-anyway restatement.

qualify + scope-drift + contract-regressions pass. No Important/Critical findings. No follow-ups.

## Lessons Learned
**What worked:** Placing the sweep as Risk Assessment's "documented-consumer twin" made it read as continuous with the existing review dimensions rather than a bolt-on; putting it in `review-work.md` (not work.md Step 7) means standalone `do-work review` inherits it for free.
**Worth knowing:** This step institutionalizes the exact discipline the manual reviews in this batch applied — and REQ-044/048's own discovered tasks (the work.md Step 8 `failed` contradiction; forensics Check 4) are textbook instances it would catch. The "a stale restatement outside the REQ's Scope is still a finding → follow-up, not scope-drift" clause is what keeps the sweep from being suppressed by scope discipline.

## Orientation
The pipeline's review step now has a required **Restatement Sweep**: when a diff redefines something other text restates (a contract token, a field's semantics, what a hash holds, a command's output shape), the reviewer greps every other statement/consumer and flags stale ones as findings — including outside the REQ's declared Scope (routed to follow-ups). Canonical home: `actions/review-work.md` Step 6 (inherited by pipeline + standalone); cross-referenced from `actions/work.md` Step 7. No map change — hardens review coverage against the REQ-035–040 restatement-drift failure class.
