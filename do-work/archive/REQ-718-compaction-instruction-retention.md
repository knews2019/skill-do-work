---
id: REQ-718
title: 'Reproduce instruction retention across compaction'
write_set:
- _dev/experiments/REQ-718/claude/baseline-prompt.txt
- _dev/experiments/REQ-718/claude/cli-help.txt
- _dev/experiments/REQ-718/claude/compact-prompt.txt
- _dev/experiments/REQ-718/claude/continuation-prompt.txt
- _dev/experiments/REQ-718/claude/experiment-report.md
- _dev/experiments/REQ-718/claude/focused-checks.json
- _dev/experiments/REQ-718/claude/input-fingerprints.json
- _dev/experiments/REQ-718/claude/live-events.jsonl
- _dev/experiments/REQ-718/claude/run-metadata.json
- _dev/experiments/REQ-718/claude/run-probe.py
- _dev/experiments/REQ-718/claude/verify-evidence.py
- _dev/experiments/REQ-718/codex-t3/baseline-observation.md
- _dev/experiments/REQ-718/codex-t3/experiment-report.md
- _dev/experiments/REQ-718/codex-t3/probe-prompts.md
- _dev/experiments/REQ-718/codex-t3/runtime-evidence.json
- _dev/experiments/REQ-718/codex-t3/source-fingerprints.json
- _dev/experiments/REQ-718/codex-t3/transcript-evidence.json
review_at: 2026-11-11T01:28:19Z
kb_status: pending
route: B
estimate:
  p50_active_minutes: 25
  confidence: medium
  calculated_at: 2026-11-11T01:15:53Z
  basis:
  - Route B
  - 2-file write set
  - 2 new files
  - 2 subsystems involved
  - 5 acceptance criteria
status: failed
created_at: 2026-11-10T23:19:33Z
user_request: UR-170
domain: testing
prime_files: []
tdd: false
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
status_changed_at: 2026-11-10T23:23:24Z
claimed_at: 2026-11-11T01:14:02Z
completed_at: 2026-11-11T01:33:19Z
error: 'Actual Codex/T3 live compaction remains Untested: no supported active-provider compaction control or observable boundary was established, and exact backing provider build/settings are unverified. Claude live capture is valid for safe-mode Claude Code 2.1.322 / claude-fable-5-1 only. Combined completion proof is unmet; no runtime retention failure inferred. Missing T3 capability is not expected to self-resolve and cannot be remediated within authorized source scope. Evidence artifacts are preserved in the merged commits.'
error_type: environment
---
# Reproduce Instruction Retention Across Compaction

## What
Run a bounded behavioral experiment that loads a real do-work phase and its required crew rules, compacts mid-phase, and inspects subsequent file reads and compliance with a harmless rule. Record the harness version and configuration so the result can be repeated before proposing instruction-restoration changes.

## Requirements
- Establish the pre-compaction baseline: the selected phase, its applicable procedure and crew files, actual reads, and the harmless rule being observed.
- Trigger compaction mid-phase and record whether the relevant files are re-read and whether the rule is followed afterward. Preserve evidence rather than relying on the model's self-report.
- Record the exact harness/version, model, relevant configuration, prompts, compaction method, and observed results with repeatable steps.
- Distinguish loss of file contents, failure to re-read, and failure to follow an available instruction. Report uncertainty when the observation cannot distinguish them.
- Keep conclusions scoped to the tested harness/configuration. A result showing retention works is a valid outcome; do not presume the video's categorical failure claims are true.

## Constraints
This is an investigative experiment, not a runtime fix. Do not add permanent hooks, reload registries, warning layers, or instruction rules. Keep the experiment within this repository's authorized boundary. Do not change global harness configuration or inspect other repositories without authorization.

## Proof of Completion
There is currently no recorded live compaction reproduction establishing the claimed failure in this repository. Completion requires a repeatable before/after check with observed file-read and harmless-rule evidence, version/configuration details, and a conclusion limited to what was observed. Static source inspection or a model's explanation alone does not satisfy this check. `tdd: false` because this is an exploratory harness experiment rather than a production behavior change with an existing automated RED/GREEN assertion.

## Finding Provenance
- Finding: F9, the only Accept verdict in the preceding validate-feedback report.
- Original source: pasted text, Matt Maher, “AI Agent Secrets Nobody Tells You”, dated November 10, 2026; transcript near 19:24–19:32.
- Verbatim source claim: “Investigating the break points that aren't written down anywhere is absolutely worth it.”
- Original severity: not supplied.
- Evidence: `skills/do-work/actions/work.md:21–25` requires procedure loads once per context, introduced by commit `22f82fed`; `skills/do-work/tools/do-work-cli/internal/hookcommands/session_start.go:40–87` performs housekeeping but does not restore active instruction files. Static inspection did not establish whether required reads recur after compaction.
- Supporting documentation: https://code.claude.com/docs/en/context-window documents bounded restoration of file and skill contents; https://learn.chatgpt.com/docs/build-skills documents progressive disclosure, not the transcript's absolute Codex restoration claim.
- Surface-cost: N/A for the bounded experiment, which adds no permanent runtime layer. New restoration hooks or enforcement need separate evidence and authorization.

## Capture Context
No queued request in any UR shares this root cause: the queue and working directories were empty at capture. The user explicitly promoted all accepted findings; only F9 was accepted. The experiment is developer-visible diagnostic work and needs exploration, so impact is user-visible and effort is substantive. See `do-work/user-requests/UR-170/input.md` for the referenced accepted finding and verbatim capture instruction.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Split Claude and actual Codex/T3 observations into isolated evidence directories; verify actual compaction or report exact blockers.
- [x] **[APPLY]:** Integrated the two isolated evidence slices serially; no runtime change.
- [x] **[UNIFY]:** Reviewed the 17-file artifact manifest, probe/verifier and scoped reports. Claude verifier and JSON parsing pass; full-range whitespace clean. Canonical final gate recorded in Testing.

*Source: “capture-request all accepted” — F9 from the preceding validate-feedback report, source: pasted video transcript.*

## Open Questions

User decisions refine the existing experiment. These answers govern its scope alongside the original requirements.
- [x] **F5** — How should the experiment handle potentially lost procedure and crew instructions? → Test the current behavior first (recommended): record which files are re-read and which rules survive; propose a fix only if evidence shows a failure.
- [x] **F6** — Should instruction-restoration hooks be investigated or implemented? → Drop the hook approach: keep the workflow dependent on explicit file loading and durable request records.
- [x] **F7** — Should the compaction experiment cover Claude Code? → Include Claude Code (recommended): test our actual skill and referenced files across compaction, recording the installed version.
- [x] **F8** — Should the compaction experiment test skill restoration in the actual Codex setup? → Include the actual Codex setup (recommended): record the harness and model, then test skill-body and referenced-file retention.


## Answer Notes

- 2026-11-10 F5: Test the current behavior first (recommended): record which files are re-read and which rules survive; propose a fix only if evidence shows a failure.
> ```
> Test the current behavior first (recommended): record which files are re-read and which rules survive; propose a fix only if evidence shows a failure.
> ```
- 2026-11-10 F6: Drop the hook approach: keep the workflow dependent on explicit file loading and durable request records.
> ```
> Drop the hook approach: keep the workflow dependent on explicit file loading and durable request records.
> ```
- 2026-11-10 F7: Include Claude Code (recommended): test our actual skill and referenced files across compaction, recording the installed version.
> ```
> Include Claude Code (recommended): test our actual skill and referenced files across compaction, recording the installed version.
> ```
- 2026-11-10 F8: Include the actual Codex setup (recommended): record the harness and model, then test skill-body and referenced-file retention.
> ```
> Include the actual Codex setup (recommended): record the harness and model, then test skill-body and referenced-file retention.
> ```

## Triage

**Route: B** — Outcome is specified; live harness compaction and observation mechanics need discovery.

## Plan

**Planning not required** — Route B: exploration-guided experiment.

## Exploration

Both Claude and Codex executables are installed. T3 exposes Codex GPT-6.1 Sol with high reasoning effort and tool-managed worktree isolation. Standalone Codex cannot stand in for actual T3. Live compaction control and observable transcripts remain unproven and are the workers’ first investigation. No existing compaction probe was found. Lessons index consulted: no matching implementation area for this artifact-only investigation; no required lessons added.

## Scope

**Files I will touch:**
- `_dev/experiments/REQ-718/claude/baseline-prompt.txt` (new)
- `_dev/experiments/REQ-718/claude/cli-help.txt` (new)
- `_dev/experiments/REQ-718/claude/compact-prompt.txt` (new)
- `_dev/experiments/REQ-718/claude/continuation-prompt.txt` (new)
- `_dev/experiments/REQ-718/claude/experiment-report.md` (new)
- `_dev/experiments/REQ-718/claude/focused-checks.json` (new)
- `_dev/experiments/REQ-718/claude/input-fingerprints.json` (new)
- `_dev/experiments/REQ-718/claude/live-events.jsonl` (new)
- `_dev/experiments/REQ-718/claude/run-metadata.json` (new)
- `_dev/experiments/REQ-718/claude/run-probe.py` (new)
- `_dev/experiments/REQ-718/claude/verify-evidence.py` (new)
- `_dev/experiments/REQ-718/codex-t3/baseline-observation.md` (new)
- `_dev/experiments/REQ-718/codex-t3/experiment-report.md` (new)
- `_dev/experiments/REQ-718/codex-t3/probe-prompts.md` (new)
- `_dev/experiments/REQ-718/codex-t3/runtime-evidence.json` (new)
- `_dev/experiments/REQ-718/codex-t3/source-fingerprints.json` (new)
- `_dev/experiments/REQ-718/codex-t3/transcript-evidence.json` (new)

**Acceptance criteria:**
- Establish actual phase/procedure/crew reads and harmless-rule baseline.
- Observe a real compaction boundary and subsequent reads and compliance.
- Record exact harness/model/configuration and repeatable prompts.
- Distinguish observed loss, absent rereads, and noncompliance, preserving uncertainty.
- Cover Claude Code and actual Codex/T3 without hooks or global changes.

## Pre-Flight

Git: only execution-lead request/run metadata changed. Focused baseline `git diff --check` exited 0. Canonical `bash _dev/tests/maintainer-verify.sh` exited 0; advance recorded both preflight and green-gate satisfied. Dependencies available.

## Implementation Summary

**Files changed:**
- `_dev/experiments/REQ-718/claude/baseline-prompt.txt` (new)
- `_dev/experiments/REQ-718/claude/cli-help.txt` (new)
- `_dev/experiments/REQ-718/claude/compact-prompt.txt` (new)
- `_dev/experiments/REQ-718/claude/continuation-prompt.txt` (new)
- `_dev/experiments/REQ-718/claude/experiment-report.md` (new)
- `_dev/experiments/REQ-718/claude/focused-checks.json` (new)
- `_dev/experiments/REQ-718/claude/input-fingerprints.json` (new)
- `_dev/experiments/REQ-718/claude/live-events.jsonl` (new)
- `_dev/experiments/REQ-718/claude/run-metadata.json` (new)
- `_dev/experiments/REQ-718/claude/run-probe.py` (new)
- `_dev/experiments/REQ-718/claude/verify-evidence.py` (new)
- `_dev/experiments/REQ-718/codex-t3/baseline-observation.md` (new)
- `_dev/experiments/REQ-718/codex-t3/experiment-report.md` (new)
- `_dev/experiments/REQ-718/codex-t3/probe-prompts.md` (new)
- `_dev/experiments/REQ-718/codex-t3/runtime-evidence.json` (new)
- `_dev/experiments/REQ-718/codex-t3/source-fingerprints.json` (new)
- `_dev/experiments/REQ-718/codex-t3/transcript-evidence.json` (new)

Claude real compaction evidence and repeatable bounded runner were added. Actual Codex/T3 baseline and capability blockers were recorded; post-compaction checks remain Untested. No retention failure was established. Both branches were serially merged; operative range is 8e7cf62aa4c75ab1e82f6e9912d174cc26fbc547..a157fd164214d36ac5f12b72949c56a1e0facc95.

## Decisions

D-01 — Expanded the two pre-authorized directory globs to exact artifact paths for the literal Scope comparator; no boundary expanded.
D-02 — QUALIFY-LIBRARY-OUTPUT flags are false positives: both Python files are standalone CLI reporters with main guards, and their JSON is the intended terminal contract. Retain outputs.
D-03 — Preserve Codex/T3 external capability blocker; do not claim combined acceptance completion or substitute another harness.

## Qualification

The 17-file merge range is substantive evidence-only work. Full diff whitespace passes; P-A-U complete. Literal Scope expanded from the two authorized globs to the exact manifest. Python output findings judged intentional standalone reporter output, not debug/library leakage. Verified eight Claude read results, boundary and rule summary using the committed verifier. Codex JSON parses and ten source hashes match. No production semantics changed. Actual Codex compaction acceptance remains Untested.

## Testing

Focused Claude capture verifier exited 0 after serial merges; all eight full Read results, real manual compact boundary, retained naming rule, compliant output and zero continuation tool calls verified. Codex three JSON artifacts parsed and ten instruction hashes match. Seventeen artifact UTF-8/JSON checks and limited credential-pattern scan pass. `git diff --check 8e7cf62a..a157fd16` passes. Scope and focused-test advance gates satisfied.

Canonical `bash _dev/tests/maintainer-verify.sh`: pre-build exit 0; final first run exit 1 from existing strict browser ~29-second timeout; mandatory identical-argv retry exit 0. The retry is the authoritative final result. No tests modified. Advance green-gate satisfied on retry status 0. Live actual Codex compaction/continuation remains Untested for exact external capability reasons in Review; passing artifact and repository checks do not fulfill that proof.

## Review

**Overall: 78%** | 2026-11-11T01:28:19Z

| Dimension | Score |
|---|---|
| Requirements | 60% |
| Code Quality | 90% |
| Test Adequacy | 60% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Partial |

**Important findings:**
- Actual Codex/T3 has no established supported active-session compaction control or verified backing build; live before/after proof is Untested. Combined completion proof is unmet. This external tooling limitation is not expected to self-resolve; terminal disposition is failed/environment, preserving the delivered evidence. A source remediation cannot supply a missing host capability within authorized scope. — impact-user-visible → report only

**Minor findings:**
- Claude safe-mode run excludes ordinary project/global configuration; evidence is limited to that explicit read-only configuration. — impact-user-visible → report only
- Claude replay writes an absolute progress-log path in this run directory; reuse after moving/removing that directory requires adjusting the path. Existing capture verification remains repeatable. — impact-negligible → report only
- Claude worker used command-local core.hooksPath=/dev/null for its commit, contrary to the skill commit rule. No setting persisted and no hook was installed; no further bypass used by integrator. — impact-negligible → report only

**Acceptance:** Partial — Claude runtime events and full read contents independently checked; Codex source hashes, parsed evidence and scoped capability discovery checked. Actual Codex boundary/continuation Untested. No retention failure established.
**Restatement sweep:** nothing redefined; this is the only wave member and changes no instruction/runtime contract.
**Suggested testing:** 1 item — repeat actual T3 with supported compaction control, observable boundary and verified backing build.
**Follow-ups created:** None (4 findings report only)

Self-validation: retained summary is evidence of naming-rule availability, not every source file; naming has overlapping causes. No claim of general retention success. Original UR and F5–F8 checked. Review performed inline per orchestrated-mode option.

## Lessons Learned

**What worked:** Explicit Claude manual compact boundary and raw read/tool results distinguish absent rereads from lost instructions.
**What didn't:** Actual Codex/T3 discovery did not expose a supported active-session compaction operation; installed PATH CLI version could not identify the backing provider build.
**Worth knowing:** A visible retained summary plus compliant output supports only the observed rule; it does not establish full file retention or causation. No relevant existing utility prime requires a satellite addition.

## Orientation

The repository now has verifiable Claude compaction evidence and an actual Codex/T3 baseline with explicit unresolved capability blockers. The combined experiment remains incomplete.

## Timing

Observed 2026-11-11T01:28:47Z to 2026-11-11T01:32:24Z: 3m 37s total, 3m 38s attributed across 1 events, 0s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| verification-gate | 3m 38s | 1 |

Slowest command: verification-gate / canonical-gate-retry, 3m 38s, exit 0, bash (2 argv tokens).
