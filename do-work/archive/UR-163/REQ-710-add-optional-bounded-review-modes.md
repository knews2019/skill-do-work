---
id: REQ-710
title: 'Add optional AI-Evaluation-Loop modes'
status: completed
heavy_verified_at: 2026-11-05T17:23:36Z
heavy_verified_revision: 74b3de5a5bd90968441dfd5c47f8d73a0baf3653
kb_status: pending
kb_entry: ""
created_at: 2026-11-04T15:54:52Z
user_request: UR-163
domain: general
route: C
re_review_at: 2026-11-05T17:19:43Z
commit: 74b3de5a5bd90968441dfd5c47f8d73a0baf3653
review_at: 2026-11-05T17:16:40Z
remediation_at: 2026-11-05T17:16:40Z
builder_handback_at: 2026-11-05T17:10:07Z
integration_at: 2026-11-05T17:10:07Z
write_set:
  - skills/do-work/SKILL.md
  - skills/do-work/actions/work.md
  - skills/do-work/actions/ai-evaluation-loop-reference.md
  - skills/do-work/actions/review-work.md
  - skills/do-work/actions/work-reference.md
  - skills/do-work/actions/help.md
  - skills/do-work/docs/work-guide.md
  - skills/do-work/docs/review-work-guide.md
planning_at: 2026-11-05T16:31:00Z
prime_files: [_dev/primes/prime-action-files.md]
impact: impact-user-visible
effort_estimate: effort-substantive
estimate:
  p50_active_minutes: 55
  confidence: medium
  basis:
  - Route C
  - 5-file write set
  - 1 new files
  - 2 subsystems involved
  - 11 acceptance criteria
  - persistence changes
  - async lifecycle behavior
  - cross-route regression gates
  calculated_at: 2026-11-04T19:23:16Z
tdd: false
maintenance: false
claimed_at: 2026-11-05T16:27:47Z
gate_deferred: 'true'
depends_on: [REQ-711]
status_changed_at: 2026-11-05T16:47:55Z
dispatch_at: 2026-11-05T16:52:22Z
completed_at: 2026-11-05T17:23:53Z
release_at: 2026-11-05T17:23:53Z
---
# Add Optional AI-Evaluation-Loop Modes

## What
Add optional AI-Evaluation-Loop modes to do-work while preserving the existing default workflow. An AI-Evaluation-Loop reviews a change, fixes accepted findings, and repeats until approved, five total review rounds are used, or two consecutive rounds produce no meaningful, evidence-supported improvement. Track findings across rounds, verify model identities, tie approval to the current change, and report unresolved outcomes honestly.

## Detailed Requirements
- Support the one-reviewer Sol-loop and three-provider council choices discussed in the report as optional modes. Ordinary runs without an explicit selection keep their current behavior.
- Offer the jester as an optional advisory reviewer alongside the selected mode. When enabled, its default model is the latest available Gemini resolved through the harness; Gemini 4 Argon is the user's example, not a permanently pinned model ID. Record the actual model used.
- Give the jester the same serious review brief and evidence without instructing it to act funny or revealing its advisory role in that brief. Its numerical scores never count toward acceptance or council grades in a PR description. Validate useful findings against the reviewed version before acting; a confirmed defect matters because of its evidence, not the advisory score. Do not wait on an unavailable, failed or slow jester; report its absence without treating it as approval or as a missing required council member.
- Allow at most five total review rounds per selected loop. Stop earlier when approved or after two consecutive rounds without meaningful, evidence-supported improvement. Fixes do not reset the total-round counter; supported improvement resets only the consecutive-stall counter. A higher subjective score alone is not supported improvement. Document these stopping conditions. Reaching either stopping limit, persistent disagreement, unavailable required reviewers, failed required reviews and missing required replies must not count as approval.
- Preserve a recoverable checkpoint of the best verified version, preferably one satisfying the requirements, so later changes can be compared against it. If no version satisfies all requirements, retain the best verified result and explicitly report what remains unmet. Any restored checkpoint needs approval covering that exact version.
- Keep a finding ledger across rounds with stable identifiers, severity, evidence, disposition and remaining concerns. Rechecks cover changed behavior and open findings; preserve evidence-based disagreement and declined findings.
- Record actual provider/model identities from available harness evidence. Do not guess identities or silently substitute a different model/provider for an explicitly selected one.
- Approve only the current behavior diff, with no unresolved material findings. A meaningful change after review invalidates that approval until rechecked.
- Make unresolved outcomes explicit in the selected mode's result and lifecycle handoff. Do not silently inherit the ordinary review-failure fallback as consensus approval.
- Document how an installed do-work user selects and uses each mode, what happens during the loop, and how approval or unresolved results are presented. Use the actual implemented invocation syntax and explain model availability requirements.

## Constraints
Reuse the existing review and feedback contracts where useful. Keep fixes within the requested scope and avoid review-driven scope growth. No mandatory council for all changes, new PR merge permissions, automatic upstream installation, or unrelated report opportunities belong to this request.

## Red-Green Proof
**RED prompt/case:** Select an AI-Evaluation-Loop mode for a small change. Have the reviewer leave a material finding open through five total review rounds while supported progress prevents the plateau stop. Separately exercise two consecutive rounds without meaningful, evidence-supported improvement, and recover the preserved checkpoint; also check that approval of an earlier diff cannot approve a subsequent behavior change. Compare an ordinary run without the option.
**Why RED now:** The inspected work.md Step 7 defines ordinary review plus one remediation and permits continuation after reviewer failure; it does not define these optional AI-Evaluation-Loop modes and their approval contract.
**GREEN when:** The opted-in mode records findings and verified identities, ends unresolved after five total review rounds or earlier after two consecutive stalled rounds when requirements remain unmet, and requires approval covering the current behavior diff. Show that fixes do not reset the five-round count, supported improvement resets only the stall count, and the preserved checkpoint is recoverable with its remaining issues reported. A clean review can approve before reaching the limit. The ordinary run retains its previous behavior. The installed usage instructions reproduce mode selection and explain these outcomes. With the optional jester enabled, record the resolved Gemini identity and show that its score cannot change acceptance; a missing advisory response does not hold the loop, while a verified material finding is handled on its evidence.
**Validation:** User confirmed on 2026-11-04.

Confirmed acceptance case: an explicitly selected review mode records findings and model identities, approves only the current change, and reports unresolved work when its round limit is reached, while an ordinary run behaves exactly as before.

User response:
> ```
> Use this acceptance case (Recommended)
> ```

## Context
Source proposal: ai-reports/2026-11-04_1524_proposal-tstack-review-workflows/index.html, option O2. Upstream source revision: t3dotgg/tstack at 1085050f3baed1ac724caa4a8166145fb2ac88af, skills/sol-loop/SKILL.md and skills/council/SKILL.md. These are research inputs, not instructions to execute upstream workflows during capture.

The earlier conversation identified missing post-install usage guidance; include it with the feature. The confirmed five-round maximum supersedes the report's earlier suggested three-round default. The lifecycle hold strategy remains a recommendation for the builder to evaluate.

Terminology confirmed by the user on 2026-11-04: use the exact name "AI-Evaluation-Loop" in user-facing documentation. One round is one complete review pass; a fresh review after fixes starts another round. The user confirmed five total rounds with early stopping after two consecutive stalled rounds and a recoverable checkpoint. Elapsed-time, token, and monetary limits are unspecified. Reviewer findings and approval are evaluations; this request does not prescribe a separate numerical grading system or a passing score for the new modes.

Fold-first scan: no pending or pending-answers REQ exists in the queue; no eligible request shares this root cause.

## Builder Guidance
Firm: explicit opt-in, unchanged defaults, a findings ledger, verified identity, current-diff approval, truthful unresolved outcomes, and an optional advisory jester defaulting to the latest available Gemini. The five-total-round maximum, two-consecutive-stalled-round stop, recoverable checkpoint and explicit unmet outcomes are also firm. Builder chooses invocation syntax, storage shape and the smallest compatible lifecycle integration. This is workflow behavior; no existing runnable automated assertion for the proposed modes was identified during capture, so tdd is false. Provide repeatable scenario evidence, and focused tests for deterministic code if introduced.

## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md`: 8398 tokens; matches action routing and status contracts. Index coverage is partial, so targeted entries are ineligible; exceeds the 2000-token budget.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`: 24067 tokens; matches native backend identity and structured evidence. Index coverage is partial, so targeted entries are ineligible; exceeds the same budget.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Validated three-task plan and eight-path scope; reuse existing lifecycle and companion prose.
- [x] **[APPLY]:** Implemented the eight declared paths; parent merged isolated source commit.
- [x] **[UNIFY]:** Builder reviewed all eight files, ran package references and contract checks, checked diff hygiene, and documented twelve simulated workflows. Parent owns integrated gate and independent review.

## Full Context
See `do-work/user-requests/UR-163/input.md` for the verbatim input.
- [x] Which stopping policy should AI-Evaluation-Loop use: five total rounds with early stopping after two stalled rounds, or five consecutive stalled rounds? → Five total rounds, stopping earlier after two consecutive rounds without meaningful, evidence-supported improvement. Preserve a recoverable checkpoint and report unmet requirements when stopping. This caps all passes; five stalled rounds and builder-selected numerical limits are out of scope.


## Answer Notes

- 2026-11-04 - [ ] Which stopping policy should AI-Evaluation-Loop use: five total rounds with early stopping after two stalled rounds, or five consecutive stalled rounds?: Five total rounds, stopping earlier after two consecutive rounds without meaningful, evidence-supported improvement. Preserve a recoverable checkpoint and report unmet requirements when stopping. This caps all passes; five stalled rounds and builder-selected numerical limits are out of scope.
> ```
> Record five total rounds, with early stopping after two stalled rounds (Recommended)
> ```

## Triage

**Route: C** — Optional evaluation modes require coordinated routing, review evidence, stopping rules, lifecycle handoff and installed usage documentation. The ordinary default must remain unchanged. All user questions are resolved.

## Plan

1. Add action-owned `--review-mode sol-loop|council` and `--jester` selection plus a core companion reference. Strip options before canonical queue selection; preserve per-REQ selection, ledger, counters and version evidence durably.
2. Integrate selected-loop review, five/two stopping, checkpoint and verified identities. Preserve ordinary review. Explicit unresolved results use existing completed-with-issues finalization, never consensus approval; existing dependency semantics remain unchanged.
3. Document installed usage and verify scenario coverage, package links, contracts and canonical repository checks.

Plan agent read the relevant action and CLI consumers. Validation: every requirement maps to these tasks; no orphan task. Existing CLI remains the sole lifecycle writer, using request_id/request_path, outcome and finalizations cleanup_complete fields. No new lifecycle schema or runtime helper is needed.

## Scope

**Files I will touch:**
- `skills/do-work/SKILL.md` (modify)
- `skills/do-work/actions/work.md` (modify)
- `skills/do-work/actions/ai-evaluation-loop-reference.md` (new)
- `skills/do-work/actions/review-work.md` (modify)
- `skills/do-work/actions/work-reference.md` (modify)
- `skills/do-work/actions/help.md` (modify)
- `skills/do-work/docs/work-guide.md` (modify)

- `skills/do-work/docs/review-work-guide.md` (modify)

**Acceptance criteria:** Explicit opt-in only; ordinary behavior unchanged; Sol and three-provider council identities verified; optional latest-available Gemini advisory remains non-voting and nonblocking; findings ledger with stable IDs and evidence; approval covers exact current candidate; five total rounds and two consecutive stalled rounds without resets from fixes/recovery; recoverable checkpoint; explicit unresolved outcome at limits or required reviewer failure; actual installed usage instructions. Repeatable scenario evidence covers the confirmed RED/GREEN cases. Finalizer-owned release metadata is separate from builder scope.

## Exploration

Read-only Explore agent verified action-owned options must be consumed before canonical advance. Existing completed-with-issues finalization remains terminal-success for dependencies, so the result must explicitly say unresolved/not approved without promising dependency holds. Review-work has caller-owned writes and follow-up seams; user guides repeat ordinary thresholds and need an explicit selected-mode exception. Scope includes review-work-guide.md to close that existing reader. No runtime helper required.

Required-lessons consult: lessons-index.md matches action-files (8398 tokens) and CLI (24067 tokens), both partial coverage and above the 2000-token projection. No required_lessons entries added. Touch-conditional action lessons remain required for builder reads.

## Decisions

**D-01 — DECIDE & STATE:** Keep selected review policy in action prose and durable REQ body evidence, consuming options before the CLI. Reuse completed-with-issues for an unresolved selected-loop result, explicitly not approval. Existing dependency readiness remains unchanged; changing dependency semantics is outside this request. No automatic upstream install or new statuses.

## Repository Gate Deferral

- **Gate command (argv JSON):** ["bash","_dev/tests/maintainer-verify.sh"]
- **Direct exit status:** 1
- **Diagnostic fingerprint:** 14816037a6fd3188551b9a540c8b7fbb01470e6bdfd681f43fa1931d3c4c6e97
- **Repair dependency:** REQ-711
- **Diagnostic evidence:** "baseline revision 40e9e4f6ee0f9ed19404d1515e7750f3a453a398"
- **Diagnostic evidence:** "canonical gate exited 1 twice before implementation"
- **Diagnostic evidence:** "TestShippedRuntimeEvidenceTracksEffectiveGoSettingsAndBinaryBytes: default runtime must have a determinable fingerprint (GOOS, GOARCH, CGO_ENABLED, DO_WORK_TEST_DO_WORK_CLI_BINARY)"
- **Diagnostic evidence:** "TestShippedGitIsolationPreservesGenericLaneInheritance: shipped runtime probe did not isolate host Git configuration"

## Implementation Summary

**Files changed:**
- `skills/do-work/SKILL.md` (modified)
- `skills/do-work/actions/work.md` (modified)
- `skills/do-work/actions/ai-evaluation-loop-reference.md` (new)
- `skills/do-work/actions/review-work.md` (modified)
- `skills/do-work/actions/work-reference.md` (modified)
- `skills/do-work/actions/help.md` (modified)
- `skills/do-work/docs/work-guide.md` (modified)
- `skills/do-work/docs/review-work-guide.md` (modified)

**What was done:** Added opt-in Sol/council selection and optional Gemini advisory review through the skill action, with a durable ledger/counters, version-specific approval, five/two stopping, recoverable checkpoints and honest unresolved lifecycle results. Wired the reference into work/review and installed help/guides while qualifying ordinary fallbacks. No CLI helper/schema or upstream install. Builder ad4ef86bb15d6d101a5abe02d4bb386668a4a186 initially merged as a7502c5a4fa6fa2717aa7fbc3889c55fff26cd7e, remediated by 0fadcd9a and merged as 74b3de5a5bd90968441dfd5c47f8d73a0baf3653; final evidence range 2fc8b1ab5cf63c7a5e883636dac270b8459b1c05..74b3de5a5bd90968441dfd5c47f8d73a0baf3653.

Handback: do-work/runs/2026-11-05-req710/REQ-710-handback.md. Its S1–S12 are explicitly simulated instruction walkthroughs, not live model execution. All lesson/crew paths required by the brief were read; none missing. Builder decisions D1–D5 in that handback implement the existing D-01 scope decision without new scope or questions.

## Remediation

Initial independent review found one Important, impact-user-visible issue: Step 4 returned accepted fixes to Step 6, whose landed-handback rule would consume the earlier merged build and fail on an empty merge. Acceptance Fail (Requirements 90, Code Quality 85, Test Adequacy 75, Scope 100; capped overall 50). Parent confirmed by following the companion, work Step 6 and fan-out merge rule.

One remediation attempt adds durable active attempt IDs, distinct handback paths and verified integration receipts before re-entering Step 6. Recovery resumes the same invocation/attempt and skips already integrated source; new fixes retain original review base, ledger and counters. Parent implemented the small fix on the existing isolated feature branch, commit 0fadcd9a. No scope expansion. Re-review approved at 95%; the integrated fast gate and all selected heavy lanes passed.

## Qualification

Initial and remediated eight-path scope matches implementation. Parent read the actual diff and traced all requirements through action routing, selected review, recovery, ledger/candidate/checkpoint and installed usage. Original merge a7502c5a passed the direct full gate (exit 0, 137s); remediation adds active-attempt recovery to work and companion. Current cumulative evidence range: 2fc8b1ab5cf63c7a5e883636dac270b8459b1c05..74b3de5a5bd90968441dfd5c47f8d73a0baf3653. No deterministic runtime logic or schema introduced; simulation plus independent instruction review is appropriate for tdd:false.

## Testing

**Scenario evidence:** Builder S1–S12 plus parent S13–S16 cover opt-in/default behavior, missing/verified identities, ledger disagreement, current-candidate invalidation, five/two stops, checkpoint protection, advisory policy, recovery and accepted-fix handback attempts. These are instruction walkthroughs, not live provider execution.
**Focused checks:** Builder package-reference and contract suites passed; integrated package-reference probe passed through advance. Cumulative diff hygiene and scope/qualification records satisfied.
**Repository gate:** bash _dev/tests/maintainer-verify.sh passed directly, exit 0, 118s at 74b3de5a5bd90968441dfd5c47f8d73a0baf3653. 427 board and 971 CLI tests; test-file budget checks passed.
**Heavy lanes:** staged-skills selected for this feature; union drain also includes the prerequisite's CLI integration, updater and installer lanes.
**Limitations:** No live provider discovery/identity, council/Sol execution, advisory timing, real checkpoint restoration or installed-agent execution.

## Review

**Overall: 95%** | REQ-710 independent re-review

| Dimension | Score |
|---|---|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 85% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings:** Initial attempt/handback routing defect fixed and independently rechecked. No remaining actionable findings.
**Minor findings:** None.
**Acceptance:** Pass — implementation-level instruction walkthroughs and integration contracts, not live model execution. Parent full gate and all selected heavy lanes passed.
**Restatement sweep:** selected/ordinary decisions, option ownership, recovery retention, cumulative merge ranges, lifecycle/dependency semantics and installed guidance agree. Source ancestry/receipt handles interrupted integration.
**Suggested testing:** selected heavy lanes.
**Follow-ups created:** None.

## Lessons Learned

**What worked:** Keeping the selected mode in a companion preserved the ordinary workflow and existing lifecycle writer.
**What didn't:** Initial simulations assumed a fresh fix builder; tracing the actual landed-handback rule exposed reuse of an already merged attempt.
**Worth knowing:** New fix attempts need distinct durable handbacks; recovery must reuse the active attempt and verify its integration receipt. Original review base, finding ledger and round counters survive both paths. Existing action lessons already cover alternate-reader closure; no duplicate satellite entry.

## Orientation

[MAP CHANGED] The core work action now offers optional AI-Evaluation-Loop modes: Sol or a three-provider council, with optional Gemini advisory feedback. It records findings and real model identities, limits review to five rounds with a two-stall stop, preserves checkpoints, and reports unresolved work honestly. Ordinary runs retain their existing behavior.

## Heavy Verification Plan

- Base revision: 2fc8b1ab5cf63c7a5e883636dac270b8459b1c05
- Target revision: 74b3de5a5bd90968441dfd5c47f8d73a0baf3653
- staged-skills: ["env", "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null", "bash", "_dev/tests/maintainer-verify.sh", "--heavy-lane", "staged-skills"] — skills/do-work/SKILL.md matched subtree skills; skills/do-work/actions/ai-evaluation-loop-reference.md matched subtree skills; skills/do-work/actions/help.md matched subtree skills; skills/do-work/actions/review-work.md matched subtree skills; skills/do-work/actions/work-reference.md matched subtree skills; skills/do-work/actions/work.md matched subtree skills; skills/do-work/docs/review-work-guide.md matched subtree skills; skills/do-work/docs/work-guide.md matched subtree skills

## Heavy Verification Result

- Target revision: 74b3de5a5bd90968441dfd5c47f8d73a0baf3653
- Execution revision: 74b3de5a5bd90968441dfd5c47f8d73a0baf3653
- staged-skills: exit 0, not skipped, executed, 34s (fingerprint_uncertain).

All selected lanes passed. The fast gate and independent review also passed; no open completion requirement. KB lessons remain in this REQ; optional promotion is pending.
