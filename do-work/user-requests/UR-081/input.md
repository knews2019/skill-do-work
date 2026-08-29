---
id: UR-081
title: 'Replace LLM bookkeeping and shipped utility logic with a Go command platform'
created_at: 2026-08-29T20:28:26Z
requests: [REQ-406, REQ-407, REQ-408, REQ-409, REQ-410, REQ-411, REQ-412, REQ-413, REQ-414, REQ-415, REQ-416, REQ-417, REQ-418, REQ-419, REQ-420]
word_count: 995
---

# Replace LLM Bookkeeping and Shipped Utility Logic with a Go Command Platform

## Summary

Capture and implement a suite-wide, Go-based `do-work-cli` that makes deterministic do-work operations available without an LLM while preserving natural-language aliases as delegating interfaces.

## Extracted Requests

| REQ | Request |
|---|---|
| REQ-406 | Create the shared CLI runtime, typed output, launcher, and Git transaction foundation |
| REQ-407 | Migrate bootstrap, install, update, reconciliation, validation, and fetching into Go |
| REQ-408 | Create the shared repository and request model packages |
| REQ-409 | Implement safe cleanup and explicit destructive repairs |
| REQ-410 | Implement doctor and deterministic forensic repairs |
| REQ-411 | Implement dependency-aware queue selection |
| REQ-412 | Implement request-state transactions |
| REQ-413 | Implement capture, answer, release, and publication transactions |
| REQ-414 | Migrate the remaining core helpers |
| REQ-415 | Migrate core and memory hooks into Go |
| REQ-416 | Implement deterministic BKB and Dream commands |
| REQ-417 | Implement interview and memory store commands |
| REQ-418 | Migrate toolbox commands and absorb audit-metrics |
| REQ-419 | Expose flat Just recipes and delegate natural-language actions |
| REQ-420 | Replace old implementations with shims and prove whole-suite parity |

## Batch Constraints

- Implement the REQs in the stated order where they share the CLI foundation.
- Use Go 1.26.1+ and avoid Python or jq as do-work implementation dependencies.
- Keep `queue-kanban` separate, absorb `audit-metrics`, and preserve every existing `.sh` path as a thin compatibility launcher.
- Every deterministic command must produce human text and stable actionable JSON from one typed result model.
- Mutations must honor the specified Git target, commit, rollback, dry-run, and destructive-operation rules.

## Full Verbatim Input

> ````text
> PLEASE IMPLEMENT THIS PLAN:
> # Capture a Go-Based No-LLM Command Platform
> 
> ## Summary
> 
> Capture one umbrella UR titled **“Replace LLM bookkeeping and shipped utility logic with a Go command platform.”**
> 
> Build one suite-wide `do-work-cli` Go module that:
> 
> - Makes every deterministic operation runnable directly and through flat Just recipes.
> - Returns actionable human text and stable JSON suitable for an LLM caller.
> - Becomes the canonical implementation behind existing natural-language actions.
> - Absorbs `audit-metrics`, while `queue-kanban` remains the separate board/UI binary.
> - Moves domain logic from all 41 shipped shell utilities and hooks into Go.
> - Removes Python and jq as do-work implementation dependencies.
> - Preserves existing `.sh` paths as thin build-and-exec compatibility shims.
> - Requires Go 1.26.1+ for installation, update, and runtime.
> - Requires Git for mutations; read-only commands remain usable outside Git.
> 
> ## Public Interface and Behavior
> 
> Create `do-work-cli` under the installed core package with global `--repo-root` and `--format text|json` options. Build it on demand when absent or older than its Go sources, using only the Go standard library unless an existing dependency is demonstrably necessary.
> 
> Expose flat Just recipes for every public command, including:
> 
> - Core: `do-work-cleanup`, `do-work-doctor`, `do-work-next`, `do-work-claim`, `do-work-complete`, `do-work-fail`, `do-work-cancel`, `do-work-unblock`, `do-work-answer`, `do-work-capture-files`, `do-work-release`, and `do-work-update`.
> - Knowledge: `bkb-init`, `bkb-status`, `bkb-lint-structure`, `dream-scan`, interview list/status/export/ingest/reset/versions, and memory remember/forget/recall/status/bootstrap/audit.
> - Toolbox: `do-work-note`, architecture-report preflight/publication, report-image generation, portfolio publication, last30days install/check, and absorbed audit-metrics commands.
> - Preserve existing board and `run-do-work-update` recipes as compatibility aliases.
> 
> All mutating commands support `--dry-run` where meaningful and optional `--commit`. Mutation rules:
> 
> - Refuse target paths already dirty relative to Git; unrelated dirty paths are allowed.
> - `--commit` additionally requires an empty existing index.
> - Stage and commit only exact touched paths, then verify the committed path set.
> - On pre-commit failure, restore clean tracked targets from Git and remove only paths created by that invocation.
> - After a successful commit, never rewrite history automatically; report an exact `git revert <sha>` command if post-commit verification fails.
> - Cleanup applies only provably safe repairs by default. Blanked-record restoration and unmerged-worktree deletion require explicit destructive flags.
> 
> Text and JSON must be rendered from one typed result model. JSON includes:
> 
> - `schema_version`, command, outcome, repository root, findings, changes, skipped work, and rollback result.
> - Each finding’s stable code, severity, affected IDs/paths, observed evidence, fixability classification, reason automation stopped, exact next argv/Just recipe, and verification command.
> - Exit codes: `0` clean/success, `1` findings or safely refused items, `2` usage/precondition/runtime failure, `3` operation failure with successful rollback, `4` incomplete rollback or committed-state risk.
> 
> Existing `do-work …` skill aliases remain unchanged but delegate every deterministic phase to `do-work-cli`; missing or failed canonical tooling stops that operation with actionable output rather than falling back to free-form mutation.
> 
> ## Capture as 15 Ordered REQs
> 
> 1. Create the shared Go module, typed result schema, text/JSON renderers, build cache launcher, Git target preflight, rollback, and commit guards.
> 2. Move bootstrap/install/update, managed-section replacement, settings reconciliation, suite validation, and archive fetching into Go; eliminate Python/jq branches and document the Go prerequisite.
> 3. Build shared REQ/UR frontmatter, schema normalization, timestamp, ID allocation, dependency graph, atomic-file, and repository-model packages.
> 4. Implement safe cleanup Passes 0–4, documentation-link repointing, merged-worktree cleanup, dry-run, explicit blank restoration, and explicit worktree discard.
> 5. Implement `doctor` and migrate deterministic forensics, blanked-record detection, timestamp auditing/repair, collision checks, and actionable remediation output.
> 6. Implement queue selection: targeted resolution, dependency readiness, cycles, waves, assignments, negligible-impact filtering, blocked probes, estimates, and composed summaries.
> 7. Implement request state transactions: claim, unblock, complete, fail, cancel, checkpoint synchronization, archival/UR closure, and calibration logging.
> 8. Implement structured capture-file publication, answered-question application, outside-text containment, status resolution, and parameterized release/version/changelog transactions.
> 9. Migrate remaining core helpers: preflight, qualification, scope drift, inventories, association, commit-hash recording, screenshot/download publication, Git helpers, and handoff surveys.
> 10. Migrate the core SessionStart hook and both memory hooks into Go subcommands, preserving exact hook stdin/stdout, redaction, deduplication, and repair behavior.
> 11. Implement BKB init/status/structural lint and Dream’s seven deterministic scans; leave contradiction resolution, synthesis, and cluster design to LLM actions.
> 12. Implement interview list/status/export/ingest/reset/versions and memory’s deterministic store operations.
> 13. Migrate toolbox publication, report-image lifecycle, last30days management, note append, and absorb the existing audit-metrics module.
> 14. Add every flat Just recipe, dynamic reserved-recipe collision validation, action delegation, help, guides, upgrade behavior, and compatibility aliases.
> 15. Replace every old utility/hook implementation with a thin shim, run whole-suite parity verification, remove retired audit-metrics sources, and enforce that shell contains no domain logic or embedded Python/jq.
> 
> REQs remain serial where they share the CLI foundation; package-specific migrations may run only after the shared runtime and repository model land.
> 
> ## Test Plan and Acceptance
> 
> - Characterize every existing utility before migration; require old shim and new Go subcommand to produce equivalent exit status, output, and filesystem effects on the existing fixtures.
> - Add Go unit and fixture tests for result-schema stability, text/JSON parity, actionable next steps, cleanup passes, queue state transitions, capture atomicity, knowledge/interview rendering, hooks, and media process handling.
> - Test dirty-target refusal, unrelated dirty paths, empty-index commit guarding, exact committed paths, rollback of tracked targets, cleanup of invocation-created paths, and post-commit risk reporting.
> - Test installer/update on fresh and existing projects, CRLF/BOM/NUL inputs, symlinks, file modes, malformed markers/JSON, custom hooks, recipe collisions, cancellation, and rollback—with Python and jq absent.
> - Keep target-specific Python checks only where the target itself is Python-based, such as Python project preflight and last30days.
> - Add a mechanical contract that every retained `.sh` file is a thin launcher only and that no shipped implementation branch embeds Python or jq.
> - Extend `maintainer-verify.sh` with `go vet` and uncached tests for `do-work-cli`; retain queue-kanban verification and replace the separate audit-metrics lane.
> - Final acceptance: install/update succeeds with Go but without Python/jq, every advertised Just command runs without an LLM, existing skill aliases still work, and every utility finding gives an LLM enough evidence and exact next actions without rescanning the repository.
> 
> After leaving Plan Mode, pass this complete request to `do-work capture-request`, record the resulting UR ID, and run `do-work verify-requests <UR-ID>` before implementation.
> ````

---
*Captured: 2026-08-29T20:28:26Z*
