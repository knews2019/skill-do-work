---
id: REQ-710
title: 'Add optional AI-Evaluation-Loop modes'
status: pending
created_at: 2026-11-04T15:54:52Z
user_request: UR-163
domain: general
route: C
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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
- skills/do-work/SKILL.md
- skills/do-work/actions/work.md
- skills/do-work/actions/ai-evaluation-loop-reference.md (new)
- skills/do-work/actions/review-work.md
- skills/do-work/actions/work-reference.md
- skills/do-work/actions/help.md
- skills/do-work/docs/work-guide.md

- skills/do-work/docs/review-work-guide.md

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
