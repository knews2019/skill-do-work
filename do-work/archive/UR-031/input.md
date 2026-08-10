---
id: UR-031
title: Split do-work into a four-skill suite
created_at: 2026-08-07T18:58:02Z
requests: [REQ-135, REQ-136, REQ-137, REQ-138, REQ-139, REQ-140, REQ-141, REQ-142, REQ-143, REQ-144, REQ-145, REQ-146]
word_count: 9
---

# Split do-work into a Four-Skill Suite

## Summary

Capture the approved, dependency-ordered program for separating the existing all-in-one do-work skill into four required skills in the same repository while preserving the current orchestrator and both update entry points.

## Verification Confirmation

The detailed planning context in this file is a synthesized record of the previously approved plan rather than part of the nine-word verbatim capture command below. On 2026-08-07, after reviewing the UR-031 verification findings, the user instructed `apply the UR-031 verification fixes`. That instruction confirms the Summary, Batch Constraints, Resolved Decisions From Planning, and Full-Cycle Prompt in this file as the authoritative implementation context for REQ-135 through REQ-146. The original Full Verbatim Input remains unchanged for provenance.

On 2026-08-08, during `do-work clarify`, the user confirmed that every known client updater reports exactly `suite-layout-v2` and chose their attestation as sufficient evidence for cutover. No stored client-by-client inventory is required. The user also confirmed that one modular suite release has shipped and every client has completed the four-skill, managed Just section, and previously enabled memory-hook migration required before compatibility cleanup.

## Extracted Requests

| REQ | Title | Role in the program |
|---|---|---|
| REQ-135 | Restore the Warning-Clean ShellCheck Baseline | Clean prerequisite |
| REQ-136 | Define the Four-Skill Suite Contract | Distribution and boundary contract |
| REQ-137 | Ship the Suite-Aware Bridge Updater | Safe migration bridge |
| REQ-138 | Add Managed Text-Section Replacement | Deterministic Justfile ownership |
| REQ-139 | Stage the Modular Core Skill | Core package |
| REQ-140 | Stage the Modular Board Skill | Board package |
| REQ-141 | Stage the Modular Knowledge Skill | Knowledge package |
| REQ-142 | Stage the Modular Toolbox Skill | Toolbox package |
| REQ-143 | Build the Full-Suite Installer and Reconciler | Fresh install and client configuration |
| REQ-144 | Activate the Four-Skill Distribution | Cutover and migration |
| REQ-145 | Remove the Stateful Pipeline | Replace pipeline state with a copyable prompt |
| REQ-146 | Remove Modular-Migration Compatibility Shims | Post-migration cleanup |

## Batch Constraints

- The suite has four required skills: `do-work`, `do-work-board`, `do-work-knowledge`, and `do-work-toolbox`.
- Fresh installs install the full suite. Modules share one repository version and update as one all-or-recover suite transaction: validate every module before writing, report success only after every installed module verifies, and restore every changed managed path after failure so no partial suite remains. This is a transactional guarantee, not a claim that four sibling directories can be replaced by one filesystem-atomic rename.
- Preserve the current feature-rich `do-work run` orchestrator. Remove only the separate stateful pipeline after modular cutover.
- Preserve both public update paths: `do-work update` and `just run-do-work-update`; both must use the same update engine.
- Use a bridge release before changing the live archive layout. Before cutover, require user attestation that every known client updater's `--capabilities` command prints exactly `suite-layout-v2`; no stored client-by-client inventory is required.
- Keep all toolbox functionality; this program separates it but does not prune it.
- Fresh installation uses one documented, copy-paste bootstrap command, creates or updates a managed Justfile section, enables core hooks, and leaves memory capture disabled.
- The exact managed Justfile sentinels are `# >>> do-work:recipes >>>` and `# <<< do-work:recipes <<<`.
- Existing known memory-hook paths should migrate automatically to the knowledge skill without enabling hooks that were previously disabled.
- `actions/kb-lessons-handoff.md` remains in core because the work and review orchestrators invoke it; knowledge storage and follow-on commands resolve through `do-work-knowledge`.
- Client repositories use Git. Dirty managed skill changes may be overwritten only after an explicit warning and confirmation.
- One-release command shims explain the new skill invocation and stop; they do not permanently forward or reimplement extension routing.

## Resolved Decisions From Planning

- **Default install:** full four-skill suite.
- **Compatibility:** one-release shims.
- **Versioning:** one suite version, not per-module versions, installed through the all-or-recover transaction defined above.
- **Migration safety:** bridge release first; existing clients update to the bridge before modular cutover, with user attestation serving as sufficient evidence that every known client reports `suite-layout-v2`.
- **Justfile management:** deterministic managed sentinels; if no Justfile exists, install a complete configured one; otherwise preserve everything outside the managed block.
- **Hooks:** enable core hooks on fresh install, but do not enable memory capture automatically.
- **Core behavior:** retain current planning, testing, review, lessons, archival, and commit behavior.
- **Pipeline replacement:** provide a prefilled sequence that captures, verifies, runs with built-in tests/review, then presents through the toolbox.

## Full-Cycle Prompt Required After Pipeline Removal

```text
Use the installed do-work suite to complete this request end to end:

1. Use do-work to capture the request below and record the resulting UR ID.
2. Run do-work verify-requests for that UR. Stop and report if verification fails.
3. Run the UR's REQs through do-work run. Require its built-in tests and review to pass.
4. Use do-work-toolbox present-work for the same UR.
5. Report the implementation, tests, decisions, and deliverable paths.

Request:
<paste request here>
```

## Full Verbatim Input

run do-work capture-requests , so we have all the REQ's

---
*Captured: 2026-08-07T18:58:02Z*
