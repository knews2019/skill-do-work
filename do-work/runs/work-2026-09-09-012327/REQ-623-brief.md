# REQ-623 Builder Brief

---
id: REQ-623
title: 'Reduce orchestration instruction duplication'
status: claimed
route: B
estimate:
  p50_active_minutes: 25
  confidence: medium
  basis:
  - Route B
  - 4-file write set
  - 5 acceptance criteria
  - cross-route regression gates
  calculated_at: 2026-09-08T22:24:16Z
created_at: 2026-09-08T21:49:05Z
user_request: UR-130
domain: general
prime_files: ["_dev/primes/prime-action-files.md", "_dev/primes/prime-shell-commands.md"]
tdd: false
maintenance: true
impact: impact-user-visible
effort_estimate: effort-substantive
depends_on: ["REQ-622"]
related: ["REQ-622", "REQ-624"]
batch: instruction-history-cleanup
write_set: ["skills/do-work/actions/work.md", "skills/do-work/actions/work-reference.md", "skills/do-work/actions/review-work.md"]
claimed_at: 2026-09-08T22:23:11Z
---

# Reduce Orchestration Instruction Duplication

## What
Remove proven duplication from the core SKILL.md, work.md, work-reference.md, and review-work.md. Preserve command contracts, judgment, authored evidence, recovery guidance, and behavior. Verify representative normal and recovery paths and report the text reduction; make no unmeasured runtime-savings claim.

## Detailed Requirements
- Reduce proven restatements in the four declared core files around their canonical owners. Preserve when to invoke commands, required inputs, typed-result interpretation, authored evidence, judgment, and recovery guidance.
- Record before/after text counts and show which surviving owner covers each removed responsibility.

## Constraints
Preserve behavior. Read instrumentation is optional for this cleanup; whole-file counts and a CLI finding code alone do not prove either runtime savings or that prose is dispensable.

## Dependencies
Depends on REQ-622 (Correct verified prose drift), so the cleanup starts from reconciled wording. REQ-624 (Addendum: Trim the shipped changelog to 50 entries) follows this request by the approved serial order.

## Builder Guidance
Use the smallest justified subtraction. No pending or pending-answers candidate in any UR shares this request's root cause; the seven existing queued requests address different release, commit-integrity, and diagnostics defects. The prior guard dedupe in REQ-023 (Dedupe intra-file guard restatements in note, scan-ideas, commit, quick-wins) touched different actions; this request targets the four named core files.

## Completion Proof
Use existing focused contract checks and representative walkthroughs of a simple request, a complex request, and an interrupted/recovery path. Each removed restatement must have a surviving owner with the same responsibility; restore guidance if its absence changes a decision. Report measured text reduction without a runtime-savings claim unless that outcome is actually measured.

## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` — 5879 tokens; condition-preserving prose extraction in the core action files. Partial slug coverage prevents targeted selection.
- `_dev/primes/lessons-shell-commands.md` — 8480 tokens; prescribed command blocks and moved prose anchors. Partial slug coverage prevents targeted selection.

## Full Context
See `do-work/user-requests/UR-130/input.md` for the user instruction, adopted report, and batch decisions. The approved prompt excerpts are in `do-work/user-requests/UR-130/assets/rev-v2-approved-scope.md`.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Apply the seven mapped subtractions in three action files; retain SKILL.md and every unique contract or judgment.
- [ ] **[APPLY]:** Make the scoped changes.
- [ ] **[UNIFY]:** Review every changed file and record the focused checks and outcomes.


## Triage

**Route: B** — Medium.

**Reasoning:** Four known instruction files, but each deletion needs proof of a surviving owner and unchanged normal/recovery decisions.

**Planning:** Not required; exploration-guided implementation.

## Plan

**Planning not required** — Route B: exploration-guided implementation.

## Exploration

Read all four candidate files. Seven proven restatements have surviving owners: hand-back merge, saved-range recovery, review follow-ups, diff selection, guardrails review, pipeline summary, and Scope mirroring. Leave SKILL.md unchanged: its compact routing and screenshot staging are unique. Preserve the direct prescribed-shell link required by existing contracts. The durable exploration is in `do-work/runs/work-2026-09-09-012327/REQ-623-exploration.md`; its mappings and walkthrough results will be promoted into this REQ.

## Scope

**Files I will touch:**
- `skills/do-work/actions/work.md` (modify) — remove mapped restatements while retaining active owner pointers
- `skills/do-work/actions/work-reference.md` (modify) — remove mapped restatements while retaining active owner pointers
- `skills/do-work/actions/review-work.md` (modify) — remove mapped restatements while retaining active owner pointers

**Acceptance criteria (restated from REQ):**
- [ ] Every removed responsibility has an unchanged surviving owner.
- [ ] Invocation, inputs, typed-result interpretation, authored evidence, judgment and recovery decisions are preserved.
- [ ] Record before/after source counts across the four examined files.
- [ ] Verify simple, complex and interrupted/deferred paths against final text.
- [ ] Existing focused contract checks pass; no runtime-savings claim is made.

## Pre-Flight

Canonical advance returned satisfied preflight and green-gate records. The direct baseline maintainer gate passed in 98 seconds (402 board and 815 CLI tests, all per-file budgets green), and the shipped-reference probe passed. Source tree is clean; only this REQ's authored metadata/run evidence is dirty.


# REQ-623 — Proven instruction duplication

Explored at source revision `5a5e8f050c530e2e7cc744e35c66594d6486ba33`, after REQ-622 (verified prose reconciliation) completed. Read all four requested files in full. This is a subtraction plan, not implementation or runtime measurement.

The strongest change is to replace the second hand-back merge procedure in `work.md` with an active instruction to follow the existing complete owner in `work-reference.md`. Six smaller cuts follow the same test. Edit three source files; leave `SKILL.md` unchanged. Its routing and screenshot-staging responsibilities are already compact, and the staging operation is not owned elsewhere.

## Baseline text counts

Measured with `wc -l -w -c` on the unchanged files. These are source-text counts, not tokens consumed in a run.

| File | Lines | Words | Bytes |
|---|---:|---:|---:|
| `skills/do-work/SKILL.md` | 63 | 698 | 5,328 |
| `skills/do-work/actions/work.md` | 593 | 13,484 | 92,875 |
| `skills/do-work/actions/work-reference.md` | 869 | 22,176 | 150,156 |
| `skills/do-work/actions/review-work.md` | 512 | 6,379 | 43,270 |
| Total | 2,037 | 42,737 | 291,629 |

The selected source spans below contain about 9 KB before replacement pointers. Report the actual net reduction after editing; do not set a reduction quota or infer runtime savings. The reference header still instructs readers to load only the named section and reuse already-loaded context.

## Deletion-to-owner plan

Line numbers identify this baseline only. Preserve the named headings and condition-bearing entry instructions.

| ID | Duplicate source | Surviving owner | Smallest safe edit |
|---|---|---|---|
| D1 | `work.md:303–311`: the condensed hand-back sequence, including owner artifact staging, index guard, `<pre>`, queue guard, merge, seams, `<merge_hash>`, and literal retention (4,249 bytes) | `work-reference.md` → **Worktree Dispatch Mode (Step 1)** → **When to merge, and the range every evidence step reads.** (lines 429–443), plus **Hold both endpoints as re-typed literals** | Replace the duplicate procedure with an active instruction to read and execute that full sequence on the integration branch at builder hand-back, before Step 6.25. Keep a compact instruction to retain the operative name and both endpoints for downstream evidence. Keep the subsequent append-only `integration_at` stamp. Retain the direct **State across command blocks** link to `../docs/prescribed-shell-primitives.md#state-across-command-blocks` in `work.md`; the canonicalization contract requires that local link. |
| D2 | `work.md:257`: saved-range resume algorithm (899 bytes including separator) | `work-reference.md` → **Repository Gate Deferral and Resumption** → **Saved-range resume proof** (lines 359–363) | Keep the trigger, `gate_deferred: true` with the paired deferred implementation fields, and direct the orchestrator to apply the named proof before deciding whether to dispatch or reuse. The full owner retains commit resolution, ancestry, rename-aware path protection, dirty-path rejection, malformed-evidence stop, drift rebuild, and fresh downstream evidence after reuse. Do not change the adjacent ordinary-parent, repair, or deferred-parent branches. |
| D3 | `work.md:382`: complete automatic review follow-up policy and repeated frontmatter list (946 bytes including separator) | `review-work.md` → **Step 10: Create Follow-up REQs** (lines 340–394); `work.md` Step 8 retains finalization preparation and addendum-cycle judgment | Keep an active instruction to apply review Step 10 to every finding, plus the local requirement to check an existing addendum cycle before creating a critical follow-up. The owner already defines impact-before-routing, critical-only queue work, exact noncritical report suffix, promotion, all generated fields, independently judged effort, and separate failure/builder/stakeholder contracts. Keep Step 7's acceptance/score remediation branches and Review template delegation. |
| D4 | `review-work.md:28–32` and the final paragraph of Step 1: repeated diff-acquisition descriptions | `review-work.md` → **Step 4: Get the Diff** (lines 70–79) | Make the Two Modes table describe trigger/location and point to Step 4 for the diff. At the no-change check in Step 1, instruct the reviewer to obtain the mode-specific diff through Step 4 before judging whether changes exist. Retain the point that a clean post-merge working tree is normal. Step 4 remains the single command owner for orchestrated working/staged/range reads and standalone `show-commit-diff.sh`. Preserve the already-green repair validator branch and its exact typed authorization. |
| D5 | `review-work.md:153–169`: five principle glosses followed by a second five-row mnemonic table (2,167 bytes with heading/separators) | `crew-members/coding-guardrails.md` → the five current principles and their exact scope; existing review **Code Quality**, **Test Adequacy**, and **Scope Discipline** dimensions | Preserve **Coding-Guardrails Principle Check (informational)**. Replace the glosses and table with an explicit instruction to read the canonical crew file and spot-check the diff against its current principles. Retain the local review judgment: this is informational, report otherwise-missed issues as Minor, and do not double-penalize across dimensions. This keeps Naming for Reach's introduced-name and idiomatic-local limits in its owner rather than recreating them here. |
| D6 | `work-reference.md:11–12`: second pipeline sequence and route summary (594 bytes) | `work.md` → **Complexity Triage** and **Steps**; the reference Architecture introduction already names the numbered steps as owner | Delete those two bullets and retain **Who owns what**. The ownership map remains useful; the duplicated sequence is unnecessary beside a sentence explicitly saying this section is not a second sequence. Preserve the Architecture heading because callers cite it. |
| D7 | `work-reference.md:579`: second Scope-to-`write_set` mirror procedure (464 bytes including separator) | `work.md` → **Step 5.5: Scope Declaration (Routes B and C)** (lines 229–247) | Replace the paragraph with a compact active pointer to Step 5.5 for the Scope-to-`write_set` mirror. Keep the template unchanged and preserve the named heading. The caller retains the route condition, immediate one-way copy, captured-value replacement, display-only role, and Route A exclusion. `capture-reference.md` points at this template, so retain an explicit pointer here rather than making the mirrored operation disappear. |

These are prose-owner mappings. No CLI finding code alone makes a paragraph disposable. D1's complete command sequence remains written once; D2's recovery judgment remains written once; D3 and D5 retain review judgment at an actively loaded owner.

## Boundaries worth preserving

- Leave the core `SKILL.md` router, full selected-action read, argument forwarding, unknown-word handling, capture-versus-execute boundary, and exclusive screenshot staging intact. Screenshot source resolution in `capture.md` Step 4 is the consumer of staging, not a replacement dispatcher.
- Keep the complete Schema Read Contract, timestamp rules, state/authority model, command-bound typed-result interpretation, canonical finalization result predicate, and recovery procedures. Their length is not evidence of duplication.
- Keep the builder's required-lesson reads and the orchestrator's claim-time lesson projection. They are different readers with different timing. The same applies to authored P-A-U/Scope/Implementation Summary records and their later independent verification.
- Keep the review report and persisted `## Review` templates separate. They intentionally have different presentation order and consumers.
- Keep the direct gate, retry, heavy hold/drain, release, review severity/score policy, and R4/R5/R6 decisions untouched. Do not add tracing or delete CLI commands.
- Keep all existing headings and bold labels cited by other files. In particular, `docs/standing-preferences.md` cites **Coding-Guardrails Principle Check**, and `capture-reference.md` cites **Scope Declaration Template**.

## Representative walkthroughs after editing

These are semantic walkthroughs to record against the final text, not claims of completed live REQ runs.

| Path | Check the final instructions still make this decision |
|---|---|
| Simple Route A | `SKILL.md` routes `run REQ-NNN` and reads the action fully. Canonical recover/advance claim first; triage skips plan/explore/Scope but still consults required lessons. Builder remains isolated. D1 reaches the complete merge owner before the summary; `<pre>..<merge_hash>` reaches qualification and review. A clean merged working tree does not cause D4's reviewer to exit with nothing to review. Existing tests, acceptance, and finalization still follow. |
| Complex Route C | Planning coverage validation and append-only planning stamp still run, then exploration. Step 5.5 writes Scope before implementation and mirrors its exact files to `write_set`; the D7 pointer cannot introduce Scope into Route A. D1's queue guard runs before merge; integration seams land inside the merge commit. D5 applies canonical guardrails with the same informational/Minor/no-double-penalty treatment. D3 sends noncritical findings to the report and critical findings through the existing fold/cycle path. |
| Interrupted or recovered build | Recover remains first and consumes per-REQ typed set-aside records. Worktree collision variants retain their operative identity; no cleanup is forced. On re-dispatch, phase stamps remain append-only while timing uses the held current dispatch instant. On a remediation merge, D1's owner retains the first `<pre>` and updates only `<merge_hash>`, so the diff is cumulative. |
| Deferred parent with an existing implementation | D2 actively reaches Saved-range resume proof. Valid unchanged source skips a second builder but reruns qualification, focused tests, canonical gate, and independent review. Proven path drift discards stale pointers and rebuilds. Malformed/non-ancestor/unverifiable evidence stops safely. The pointer must not turn an unverifiable range into a rebuild or a completion. |
| Standalone review and no-op repair | D4 obtains standalone changes through merge-aware `show-commit-diff.sh`; its archived UR read stays context-only. Orchestrated already-green repair still uses the exact validator, request/writer/timestamp inputs and `already_green_repair.review_allowed`; an empty ordinary diff gains no no-op exception. |

For each row, record where the decision is made in the resulting files. If any decision requires reconstructing a deleted rule instead of following a surviving directive, restore that guidance.

## Existing verification

Parent-selected checks are appropriate: `bash _dev/tests/shipped-package-reference-contract.sh`, then the existing full `bash _dev/tests/maintainer-verify.sh` gate. The latter already includes `action-shell-blocks.sh` and `prescribed-shell-canonicalization.sh`; do not add tests or modify tests merely to accommodate prose movement.

The source inspection found a specific mechanical constraint: `_dev/tests/prescribed-shell-canonicalization.sh:88–100` requires `work.md`, `work-reference.md`, and `review-work.md` each to carry their own `../docs/prescribed-shell-primitives.md` pointer. D1 must retain it. The shipped reference check covers both source and installed sibling topology and checks named section citations. An ordinary diff-hygiene check plus the owner matrix and walkthroughs provide the remaining review evidence; syntax/link checks alone do not establish semantic preservation.

No implementation or tests were run by this exploration. Source files and queue metadata were not changed.

## Context read

Read the REQ and approved UR-130 scope; `CLAUDE.md`; the full four source files; the action-files and shell-commands primes; both corresponding lesson satellites; and the relevant general, coding-guardrails, maintenance, shared-principles, communication-style, anti-slop, and prompt-injection crews. The satellites were read for their touch-conditional guidance, despite capture's budget drops. The most relevant lessons were condition-preserving extraction, keeping active delegation at the decision boundary, preserving moved anchors, proving pointer resolution separately from deletion, and naming a successor for every removed responsibility. No required path was missing.


## Builder boundary

Implement in the isolated worktree only, branch worktree-agent-REQ-623-instruction-duplication. The stale do-work tree inside the checkout is absent for your purposes: never read/write/stage any of it. Do not edit version or changelog mirrors. Parent owns all queue evidence, release and finalization. Write your full hand-back only at the exact absolute main-tree REQ-623-handback.md path supplied in dispatch, never commit it. Include actual file manifest, P-A-U, Decisions, before/after counts for all four candidate files, final owner mapping, concrete final-text walkthrough evidence, tests run/outcomes, lessons and impact-stamped discoveries. Preserve all existing heading anchors. No source outside the three declared paths and no new tests. Commit the completed source increment; no rebase, force, merge or main-tree edits.
