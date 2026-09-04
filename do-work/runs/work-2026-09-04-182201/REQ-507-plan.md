# REQ-507 Plan — Hand Archive and Commit Tails to `advance`

## Boundary and approach

Keep the existing `finalization` engine as the only archive/release/commit/provenance authority. Change lifecycle `advance` from the current final agent-judgment stop into a request-bound mechanical `finalize` phase: without a manifest it reports typed `needs_input`; with one exact action-authored manifest it invokes the existing `finalize` handler in-process and projects that handler's ordered typed result unchanged. The action continues to choose and author judgment inputs, then consumes the result instead of naming or replaying tail mechanics.

## Exact file map

- `skills/do-work/actions/work.md` — collapse Step 8 to Fold-First/follow-up, impact, terminal/release, lessons, and cleanup judgment; collapse Step 9 to manifest authoring plus consumption of finalization evidence returned by `advance`.
- `skills/do-work/actions/work-reference.md` — delete mechanical parts of the Changelog Entry Procedure and Commit & Metadata-Commit Procedure; retain changelog title/voice and release/provenance judgment plus the exact typed-result acceptance rule.
- `_dev/tests/contracts/core-checks.sh` — replace retired sentence predicates with structural guards that require the concise judgment boundary and forbid restoring direct `finalize`, archive, staging, commit, and provenance recipes in the action.
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands.go` — classify a reviewed/oriented working REQ as mechanical phase `finalize`, advertise a request-bound `advance` continuation that requires one manifest path, and route supplied finalization input to the composition seam.
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/finalization_gate.go` (new) — parse one exact `--finalization-manifest`, bind request identity/path before mutation, invoke `finalization.Handlers()[finalize]` in-process, and return the subordinate findings, changes, rollback, ordered `finalizations`, and compatibility singular record alongside the active `advance` identity.
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/finalization_gate_test.go` (new) — public CLI RED/GREEN characterization for serial, supplied-worktree, completed-with-issues, already-green/no-release, missing input, identity mismatch, and hostile/irrelevant input.
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands_test.go` — update the phase matrix from `agent judgment: prepare finalization manifest` to mechanical `finalize` and verify its tokenized manifest continuation.
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go` — prove text/JSON parity and non-null normalization when an `advance` result carries ordered finalization records; no new schema is planned because `FinalizationResult` already carries the consumer contract.
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` — update the phase table and ownership map so the prime no longer teaches finalization as an agent-judgment stop.

No change is planned in `internal/finalization/` or `internal/resultmodel/result_model.go`: the existing handler and result type already express the required transaction and evidence. Add either only if a public advance test exposes a real missing field or unsafe handler seam; record that as a scope decision before editing.

## Ordered tasks

1. **Write RED public behavior tests.** Add the four finalization-path fixtures at the public `do-work-cli ... advance REQ-NNN` seam and update the phase-matrix expectation. Before implementation, they should fail because the oriented phase is agent judgment and finalization arguments are rejected/ignored rather than executed.
2. **Compose finalization behind advance.** Add strict finalization-input parsing and identity binding, call the existing typed finalizer without shell interpolation or duplicated lifecycle logic, and preserve the subordinate outcome/rollback/findings/changes and ordered records. A missing manifest stays a typed action-owned-input finding; malformed, irrelevant, mismatched, or multi-manifest input refuses before mutation.
3. **Reduce prose while preserving judgment.** Replace Step 8/9 mechanics with the remaining Fold-First, sweep/impact, terminal/release-content, failure, deferred-lesson, and worktree-cleanup judgments. Reduce both reference procedures to semantic source/bump/ownership/title/voice/provenance choices and the typed acceptance condition. Update the CLI prime at the same time.
4. **Replace predicates and verify GREEN.** Make core contracts assert the new ownership boundary rather than old sentence spellings, run focused lifecycleadvance/finalization/resultmodel tests, race tests, module tests, contract regressions, vet, and the direct repository gate; inspect the diff for restored alternate archive/commit writers.

## Typed finalization evidence the work action consumes

Use top-level `outcome`, `findings`, `changes`, and `rollback`, then consume ordered `finalizations` as authority; singular `finalization` remains compatibility-only when exactly one record exists. For the active record require and retain:

- identity: `request_id`, original `request_path`, `archive_path`, and `journal_path`;
- terminal state: `phase`, `terminal_status`, `resumed`, and `discovered`;
- exact mutation/provenance: `commit_paths`, settled `primary_commit` / `metadata_commit`, and invocation-local `created_primary_commit` / `created_metadata_commit`;
- refusal/recovery: `blocked_paths`, `reason_codes`, `next_argv`, `verification_argv`, and `collection_argv`.

Work may continue only when the global outcome is success, exactly one record matches the active REQ/path, `phase == cleanup_complete`, and both `blocked_paths` and `reason_codes` are empty. It reports the archive and settled/created hashes from that record and performs retained isolated-worktree cleanup only after this typed success. It must not infer completion from Git status, display text, the singular compatibility projection, or archive presence.

## RED/GREEN path matrix

- **Serial:** `primary_commit`, no supplied implementation hash, project/release/lifecycle paths in the strict manifest. GREEN proves archive + optional release payload + primary commit + metadata provenance commit, with matching created/settled hashes and a clean tree.
- **Worktree:** `supplied_commit` with a resolvable ancestor merge hash. GREEN proves the archive records that supplied hash, the tail creates its exact primary commit, and no metadata commit is created.
- **Completed with issues:** `transition: complete`, `terminal_status: completed-with-issues`. GREEN proves the status survives into the archive/result and otherwise follows the same successful finalization contract.
- **Already green:** `primary_commit`, no implementation hash, no release manifest or `release_at`, and lifecycle-only commit paths. GREEN proves archive/provenance succeeds without changing changelog/version/lockfile bytes or admitting project paths.

Each case also asserts request/path binding, ordered record shape, `cleanup_complete`, empty blockers/reasons, and byte/commit effects. Add a refusal case proving a manifest for another REQ or path causes no mutation.

## Judgment that must remain prose

- Fold-First choice: fold into an existing follow-up/sweep/stakeholder REQ or mint a new one; cold-reader question wording and cycle detection.
- Impact stamping and which findings/discoveries are critical enough to queue; noncritical findings remain report-only.
- Terminal/failure classification, including `completed` versus `completed-with-issues`, intent/spec/code/environment judgment, and the external-precondition blocked flip.
- Release semantics: whether a release exists (including already-green no-release), affirmative project ownership, source and bump size, required mirrors, changelog title and human voice/content, and exact payload bytes supplied to the manifest.
- Deferred lesson relevance/family promotion and worktree identity/cleanup judgment after typed success.

## Stale capture paths and scope risks

- The captured `_dev/tests/contract-regressions.sh` path is now only the dispatcher; `_dev/tests/contracts/core-checks.sh` owns the active work-action predicates.
- The captured `internal/lifecycleadvance/` directory is not an auditable file scope; the exact source/test files above should replace it.
- `prime-do-work-cli.md` was omitted at capture but its advance phase table would immediately become stale; include it as a required scope seam.
- The report measured an older 20-step/850-line action. Current headings and line counts have already changed through REQ-504–506; classify by semantic ownership, not captured line numbers or the old Step labels.
- Avoid widening into `internal/finalization/`: its direct tests already cover journal recovery and transaction semantics. Public advance tests should characterize composition, not duplicate or rewrite the engine.
- A restatement sweep may find readers that name the Commit Phase. Update only a reader whose behavior contract becomes false; references to the named phase or to finalization recovery are not automatically stale.
- Prose deletion is safe only as a four-part move: advance composition, deleted action/reference mechanics, deleted predicates, and public behavior tests land together. Retain all manifest judgments the CLI cannot infer.
