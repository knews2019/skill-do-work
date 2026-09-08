---
id: REQ-623
title: 'Reduce orchestration instruction duplication'
status: completed
route: B
review_at: 2026-09-08T22:44:02Z
kb_status: pending
commit: 3298409c7b623a32f567acfec6d09bbd3c46c6e6
heavy_verified_at: 2026-09-08T22:44:02Z
heavy_verified_revision: 3298409c7b623a32f567acfec6d09bbd3c46c6e6
integration_at: 2026-09-08T22:38:27Z
builder_handback_at: 2026-09-08T22:37:56Z
dispatch_at: 2026-09-08T22:31:03Z
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
completed_at: 2026-09-08T22:44:22Z
release_at: 2026-09-08T22:44:22Z
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
- [x] **[APPLY]:** Integrated the seven mapped subtractions in three action files, preserving SKILL.md.
- [x] **[UNIFY]:** Read the complete changed-source diff and durable hand-back; focused checks, owner equality and six final-text walkthroughs passed. Merged-tree gate follows.


## Triage

**Route: B** — Medium.

**Reasoning:** Four known instruction files, but each deletion needs proof of a surviving owner and unchanged normal/recovery decisions.

**Planning:** Not required; exploration-guided implementation.

## Plan

**Planning not required** — Route B: exploration-guided implementation.

## Exploration

Read all four candidate files. Seven proven restatements have surviving owners: hand-back merge, saved-range recovery, review follow-ups, diff selection, guardrails review, pipeline summary, and Scope mirroring. Leave SKILL.md unchanged: its compact routing and screenshot staging are unique. Preserve the direct prescribed-shell link required by existing contracts. The complete surviving-owner mapping, measurements and final walkthrough results are preserved below.

## Scope

**Files I will touch:**
- `skills/do-work/actions/work.md` (modify) — remove mapped restatements while retaining active owner pointers
- `skills/do-work/actions/work-reference.md` (modify) — remove mapped restatements while retaining active owner pointers
- `skills/do-work/actions/review-work.md` (modify) — remove mapped restatements while retaining active owner pointers

**Acceptance criteria (restated from REQ):**
- [x] Every removed responsibility has an unchanged surviving owner.
- [x] Invocation, inputs, typed-result interpretation, authored evidence, judgment and recovery decisions are preserved.
- [x] Record before/after source counts across the four examined files.
- [x] Verify simple, complex and interrupted/deferred paths against final text.
- [x] Existing focused contract checks pass; no runtime-savings claim is made.

## Pre-Flight

Canonical advance returned satisfied preflight and green-gate records. The direct baseline maintainer gate passed in 98 seconds (402 board and 815 CLI tests, all per-file budgets green), and the shipped-reference probe passed. Source tree is clean; only this REQ's authored metadata/run evidence is dirty.

## Implementation Summary

**Files changed:**
- `skills/do-work/actions/work.md` (modified) — delegate the repeated hand-back merge and saved-range recovery procedures; keep review follow-up policy at the review action's owner.
- `skills/do-work/actions/work-reference.md` (modified) — remove the duplicate pipeline/route gloss and delegate Scope mirroring to its existing action step.
- `skills/do-work/actions/review-work.md` (modified) — keep diff commands at Step 4, preserve the standalone commit-field input, and read canonical coding guardrails instead of repeating the five glosses and table.

**What was done:** Replaced repeated prose with active owner pointers at the same decision boundaries. Preserved every heading, the direct prescribed-shell link, the full existing mechanical procedures, schema/evidence rules, recovery branches, review scoring and follow-up policy. No runtime implementation changed.

**Integration range:** `2c327607f6025d54f7555ac00582e5b9bb8de583..3298409c7b623a32f567acfec6d09bbd3c46c6e6`.

## Text counts

Measured with `wc -l -w -c`. Before is the base revision above; after is the hand-back commit. All four requested candidate files are included.

| File | Lines before → after | Words before → after | Bytes before → after | Bytes removed |
|---|---:|---:|---:|---:|
| `skills/do-work/SKILL.md` | 63 → 63 | 698 → 698 | 5,328 → 5,328 | 0 |
| `skills/do-work/actions/work.md` | 593 → 586 | 13,484 → 12,711 | 92,875 → 87,769 | 5,106 |
| `skills/do-work/actions/work-reference.md` | 869 → 867 | 22,176 → 22,030 | 150,156 → 149,241 | 915 |
| `skills/do-work/actions/review-work.md` | 512 → 496 | 6,379 → 6,026 | 43,270 → 40,922 | 2,348 |
| Total | 2,037 → 2,012 | 42,737 → 41,465 | 291,629 → 283,260 | 8,369 |

Net: 25 lines, 1,272 words, and 8,369 bytes removed; 2.87% fewer source bytes. This does not measure the subset a route reads or its runtime cost.

## Surviving-owner evidence

The IDs below match the exploration. Line numbers are from the hand-back commit. The complete surviving merge/worktree section, saved-range proof, follow-up creation step, diff-acquisition step, complexity-triage section, and Scope declaration step were compared against their base text and found byte-identical. The canonical coding-guardrails file and `SKILL.md` were also byte-identical. All Markdown heading sequences in the three edited files were unchanged.

| ID | Removed repetition and remaining entry point | Exact surviving owner and preserved responsibility |
|---|---|---|
| D1 | `work.md:303` now actively reads and executes the complete merge owner before Step 6.25, with the operative name and both endpoints retained. | `work-reference.md:428`, **When to merge, and the range every evidence step reads.**, retains owner-artifact/index handling before `<pre>`, the pre-merge queue guard, `--no-ff --no-commit`, empty-handback refusal, seam placement inside the merge, endpoint capture and downstream range use. The same unchanged section retains literal reuse (`:437`), cumulative remediation (`:441`), and non-forced cleanup (`:445`). The direct State across command blocks link remains in `work.md`. |
| D2 | `work.md:257` retains `gate_deferred: true` and both deferred implementation field names, then actively applies the named proof. | `work-reference.md:357`, **Saved-range resume proof**, retains paired resolution, non-empty ancestry, rename-aware protected paths, current/history drift checks, malformed-evidence refusal, rebuild only on proven drift, and fresh qualification/tests/gate/review after reuse. The adjacent ordinary/repair/deferred branch table in `work.md` is unchanged. |
| D3 | `work.md:375` applies review Step 10 to every finding and retains the local addendum-cycle check. | `review-work.md:324`, **Step 10: Create Follow-up REQs**, retains impact-first classification, no re-scoring, exact report-only suffix and no noncritical queue writes, explicit capture promotion, critical fold-first creation, all generated fields, separate effort judgment, and distinct failure/builder/stakeholder contracts. `work.md:448` Step 8 retains finalization preparation and cycle judgment. |
| D4 | The modes table retains triggers and locations; `review-work.md:48` sends the ordinary no-change decision to Step 4 and explicitly supplies the REQ's `commit` field for standalone review. | `review-work.md:70`, **Step 4: Get the Diff**, retains working/staged/range selection and the merge-aware standalone helper. The no-op validator paragraph at `:46` is byte-identical and remains before the ordinary diff pointer; its typed authorization is not replaced by an empty-diff judgment. |
| D5 | `review-work.md:151–153` retains **Coding-Guardrails Principle Check (informational)** and explicitly reads the canonical crew. | `crew-members/coding-guardrails.md` is unchanged and retains every principle, condition, and exception. The local review still reports otherwise-missed issues as Minor and forbids double penalties. Existing Code Quality, Test Adequacy, and Scope Discipline dimensions are unchanged. |
| D6 | `work-reference.md:7` retains Architecture as an ownership map with its existing pointer to the action's numbered steps. | `work.md:41`, **Complexity Triage**, and its numbered steps remain the route/sequence owners. The reference's **Who owns what** bullet remains unchanged. |
| D7 | `work-reference.md:577` actively applies the Scope-to-`write_set` mirror through the named action step. The template is unchanged. | `work.md:229`, **Step 5.5: Scope Declaration (Routes B and C)**, retains immediate one-way Scope-to-frontmatter copying, replacement of capture hints, display-only semantics, and the Route A exclusion. The pointer remains visible to capture-reference's template reader. |

## Final-text walkthroughs

These are semantic walkthroughs of the final instructions, not live executions of additional REQs.

| Scenario | Observed final instruction path and result |
|---|---|
| Simple Route A, ordinary successful build | The unchanged router selects `work.md` and requires reading it fully. Work Step 1 still recovers and canonically claims before triage. Route A skips plan/explore/Scope while Step 5 still consults lessons. The builder remains isolated. D1 reaches the unchanged complete merge procedure before the mandatory Implementation Summary. Qualification and review consume the saved merge range. D4's ordinary no-change pointer reaches review Step 4, whose worktree branch reads that range, so a clean working tree does not falsely end the review. Tests, acceptance, and canonical finalization remain present. Result: same decisions. |
| Complex Route C with scope and a review finding | Planning and coverage validation remain at work Step 4, exploration/lessons at Step 5, and the unchanged Step 5.5 writes Scope then mirrors its exact files to `write_set`. D7 points back to that route-conditioned owner, so it does not impose Scope on Route A. D1 still runs the queue guard before merge and includes seams inside the merge commit. D5 reads canonical guardrails with unchanged informational/Minor/no-double-penalty treatment. D3 sends every finding through the unchanged critical-only policy and preserves the cycle check. Result: same authoring, scope, review and follow-up decisions. |
| Interrupted build and remediation merge | Recovery remains first. Per-REQ set-aside handling, collision variants, branch/worktree ladder, and cleanup restrictions are unchanged in the complete owner. Dispatch timestamps remain append-only while timing consumes the held current dispatch instant. D1's owner keeps the first `<pre>` on remediation and updates only `<merge_hash>`, retaining the cumulative diff. Cleanup remains after typed finalization success and refuses dirty/unmerged work without force. Result: same recovery, attribution and cleanup decisions. |
| Deferred parent with a saved range | D2's trigger reaches Saved-range resume proof before reuse/rebuild is chosen. Unchanged valid source reuses the implementation but reruns qualification, focused tests, the canonical gate and independent review. Proven path drift removes the pair and rebuilds. A missing endpoint, invalid ancestry, unreadable path, ambiguous rename or other unverifiable evidence stops safely. The pointer does not offer rebuild as a fallback for an invalid range. Result: same three-way reuse/rebuild/refuse decision. |
| Standalone review | Step 1 still finds only eligible archived targets. The ordinary diff pointer explicitly takes the REQ's `commit` frontmatter and Step 4 uses `show-commit-diff.sh`, retaining first-parent merge handling. The archived UR input remains context-only. Result: the same source hash, diff helper and archive authority. |
| Already-green repository-gate repair | The exact validator branch is unchanged and precedes the ordinary no-change pointer. It still requires exact request/writer/timestamp inputs and typed `already_green_repair.review_allowed: true`; an empty ordinary diff gains no exception. No new gate run or no-op predicate was added. Result: same typed authorization. |

The standalone walkthrough found that the old modes table was the most explicit statement of where its commit argument comes from. That input was retained beside the Step 4 pointer before hand-back. No deleted judgment required reconstruction from memory after this correction.

## Decisions

- **D-01 — DECIDE & STATE:** Keep `SKILL.md` unchanged after full examination. Its routing and exclusive screenshot staging have unique responsibilities; changing it solely to touch all four candidate files would not be proven subtraction. This follows the approved narrowed Scope.
- **D-02 — DECIDE & STATE:** Preserve the standalone `commit` frontmatter input at the ordinary diff entry point while leaving Step 4's exact command owner intact. The table's command restatement can disappear without leaving the helper's argument source implicit.


## Qualification

Passed — canonical advance returned satisfied qualify and scope-drift records without findings at the retained merge range. All three changed files match Scope; SKILL.md is unchanged after deliberate examination. Each removed responsibility has an active pointer to an existing owner. Parent reviewed the complete source diff and retained standalone commit-field input. Source counts and semantic walkthrough evidence are recorded above, with no runtime claim.

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh` on merged main through the timing wrapper; shipped-reference probe through advance. Builder also passed action-shell-blocks and prescribed-shell-canonicalization.
**Result:** All passed. Merged gate took 96 seconds with 402 board and 815 CLI tests; all per-file budgets below 30 seconds. Canonical scope-drift, focused-test and green-gate records were satisfied at `3298409c7b623a32f567acfec6d09bbd3c46c6e6`. The six semantic walkthroughs are recorded above; they are not claimed as extra live REQ executions. No new tests.

**Heavy verification plan:**
- Range: `2c327607f6025d54f7555ac00582e5b9bb8de583..3298409c7b623a32f567acfec6d09bbd3c46c6e6`
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — skills/do-work/actions/review-work.md matched subtree skills; skills/do-work/actions/work-reference.md matched subtree skills; skills/do-work/actions/work.md matched subtree skills

## Review

**Overall: 100%** | 2026-09-08T22:41:29Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 100% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:** None
**Acceptance:** Pass — Seven deletion-to-owner mappings, six final-text paths, exact source counts, focused checks, and merged gates preserve the requested behavior; 8,369 bytes removed without a runtime-savings claim.
**Suggested testing:** 0 items
**Follow-ups created:** None (0 findings report only)

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Comparing complete surviving owner sections byte-for-byte made the deletion proof explicit. The final walkthrough independently checked that each active call site still reaches the owner under the same condition.

**What didn't:** Merely pointing the shortened modes table at Step 4 left the standalone commit argument's source less explicit. The walkthrough caught that and retained the field at the call site before hand-back.

**Worth knowing:** The source contract requires `work.md` itself to retain a direct prescribed-shell guide pointer even when its detailed merge commands live in the reference. The no-op validator must remain distinct from the ordinary empty-diff check. Source-byte reduction is not a runtime measurement.

## Orientation

Work and review now delegate repeated procedures to their existing owners. The action and shell primes still resolve their referenced owners; no subsystem map changed. The existing condition-preserving extraction lessons cover this result; no new satellite rule is needed.

## Heavy Verification Plan

Base revision: `2c327607f6025d54f7555ac00582e5b9bb8de583`
Target revision: `3298409c7b623a32f567acfec6d09bbd3c46c6e6`

- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — skills/do-work/actions/review-work.md matched subtree skills; skills/do-work/actions/work-reference.md matched subtree skills; skills/do-work/actions/work.md matched subtree skills

## Heavy Verification Result

Target and execution revision: `3298409c7b623a32f567acfec6d09bbd3c46c6e6`. Stored plan reproduced exactly.
- staged-skills: exit 0, 33 seconds; executed (fingerprint_mismatch), no skip.

## Timing

Observed 2026-09-08T22:24:39Z to 2026-09-08T22:40:04Z: 15m 25s total, 10m 07s attributed across 3 events, 5m 18s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| builder-work | 6m 53s | 1 |
| exploration-preflight | 1m 38s | 1 |
| verification-gate | 1m 36s | 1 |

Slowest stage: builder-work / implementation, 6m 53s, outcome success.
Slowest command: exploration-preflight / baseline-gate, 1m 38s, exit 0, bash (2 argv tokens).
