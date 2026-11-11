---
id: REQ-718
title: 'Reproduce instruction retention across compaction'
write_set:
- _dev/experiments/REQ-718/claude/**
- _dev/experiments/REQ-718/codex-t3/**
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
status: claimed
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
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
- `_dev/experiments/REQ-718/claude/**` (new): repeatable probe and observed evidence.
- `_dev/experiments/REQ-718/codex-t3/**` (new): actual harness probe and observed evidence.

**Acceptance criteria:**
- Establish actual phase/procedure/crew reads and harmless-rule baseline.
- Observe a real compaction boundary and subsequent reads and compliance.
- Record exact harness/model/configuration and repeatable prompts.
- Distinguish observed loss, absent rereads, and noncompliance, preserving uncertainty.
- Cover Claude Code and actual Codex/T3 without hooks or global changes.

## Pre-Flight

Git: only execution-lead request/run metadata changed. Focused baseline `git diff --check` exited 0. Canonical `bash _dev/tests/maintainer-verify.sh` exited 0; advance recorded both preflight and green-gate satisfied. Dependencies available.
