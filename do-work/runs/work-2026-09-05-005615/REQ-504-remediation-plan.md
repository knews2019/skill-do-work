# REQ-504 Remediation Plan

Fix the disagreement between checkpoint readers and writers while retaining the already-supported legacy layout. A format migration is unnecessary: refresh can update its frontmatter and preserve the existing body, and authorized removal can consume the same claim range that discovery already recognizes.

Read-only planning at current revision `2ba5b432658853690e8e5a6d20bd2dcc147e9ada`. No source, tests, queue records, or lifecycle state changed; no tests ran in this planning pass. The exact failed public fixture and outputs are in `REQ-504-review.md` beside this file.

## Exact proposed scope

All paths below are relative to `skills/do-work/tools/do-work-cli/`:

- `internal/repositorymodel/repository_model.go`
- `internal/repositorymodel/repository_model_test.go`
- `internal/requeststate/state_apply.go`
- `internal/requeststate/state_apply_test.go`
- `internal/lifecycleadvance/checkpoint_commands.go`
- `internal/lifecycleadvance/checkpoint_commands_test.go`
- `internal/lifecycleadvance/recovery_commands_test.go`

No changes should be needed in recovery orchestration, state-plan authority rules, publication, cleanup, result schema, action prose, or the contract aggregator. If a discovered caller requires widening, record the concrete reason before doing so.

## Three coherent tasks

1. **Make the public failures RED first.** Add `TestRecoverLegacyCheckpointClaimsThroughPublicCommand` to the recovery acceptance file. Use the reviewed committed working-REQ fixture and a heading-less checkpoint, first prove observation leaves bytes intact and returns the actual writer, then require explicit-authority recovery to succeed, remove every matching legacy record and continuation, preserve unrelated records, and allow a fresh claim. Add `TestAdvanceCheckpointPreservesLegacyClaimDiscovery` to checkpoint acceptance tests: capture public recovery evidence, refresh, capture again, and compare semantic identity/writer/header evidence while allowing source-line offsets caused by frontmatter. Assert the retained entry and continuation bytes are exact and `preserved_claims` agrees with post-write discovery. Before production edits, the first test must fail on `RECOVER-CLAIM-CHECKPOINT-EVIDENCE` and the second on disappearance of checkpoint evidence. Include multiple labels, an unlabelled record, and a numeric alias as subcases of these failures rather than adding redundant smoke tests.

2. **Use one structural claim range and preserve the existing layout.** Export or narrowly expose the repository model's existing checkpoint-section helper as a descriptive two-word-or-longer identifier, such as `CheckpointClaimBounds`. Return the claim range and whether a real canonical heading exists: exact canonical heading takes precedence through the next level-two heading, otherwise the already-supported whole-document legacy range applies. Keep CRLF handling consistent with current discovery. Have repository projection, request-state removal, and the `checkpoint-absent` predicate consume that same range; this closes both the observed refusal and the false assertion of absence for a real legacy entry. Keep request identity matching and writer authorization unchanged. In `checkpointSessionBytes`, delete the branch that appends an empty canonical heading to an existing nonempty body; continue creating the normal canonical body for a genuinely empty/new checkpoint. Preserve all non-owned body bytes. Update the discovery comment that currently promises a writer-driven upgrade, since preserving a supported layout is the chosen fix.

   **Account for the fresh-claim caller in the same file.** `checkpointWithClaim` currently appends a canonical section too, which can hide another request's retained legacy evidence immediately after successful recovery. When an existing checkpoint is in legacy mode and contains observed legacy claims, append the new canonical entry line without introducing a competing section; continue using the normal section writer for canonical or empty documents. This is part of the public recovery → fresh claim acceptance control, not a separate feature. Prefer this small branch over relocating foreign records or inventing a migration. The shared range helper can distinguish structure; count/identify actual legacy claims through existing structural discovery logic rather than a loose substring check.

3. **Prove the boundaries and complete the ordinary hand-back.** Run the two named public RED→GREEN tests, then the focused owner packages. Add bounded unit controls in repository-model and request-state tests: a real canonical section ignores matching-looking entries in Completed/Notes sections; a heading token inside ordinary prose is not a heading; CRLF canonical records and LF legacy records retain their contract; absent evidence is refused when a legacy record really exists; unrelated request blocks stay byte-identical. Assert writer-specific removal retains another writer and unlabelled records, while authorized all-entry removal deletes only the requested identity. The fresh-claim control must retain the other legacy request's public evidence. Run the canonical required integration/repository verification only through the orchestrator's normal phase commands after integration; do not reuse the old target's heavy results for changed source.

## Scope and authority cautions

`RemoveAllCheckpointClaims` now also serves terminal removal in `state_plan.go`, and `RemoveOwnedCheckpointClaim` is called by cleanup. The shared range change therefore affects those callers even without editing them. Preserve their existing identity/authority predicates and cover the helpers' canonical negative controls; the accepted legacy layout should be consistently discoverable and removable, not treated as new removal authority.

REQ-544 remains a separate pending caller-authored-text anchoring sweep. Its named publication and cleanup instances are outside this seven-file remediation. Its earlier cleanup snippet has already become a call to `RemoveOwnedCheckpointClaim` at this revision, so the orchestrator should reassess that instance when REQ-544 is claimed. This plan does not fold or retire it. Exact heading recognition here is needed for the same checkpoint interpretation across the affected owners; broad substring-gate cleanup would be scope expansion.

Do not backfill the reviewed request with a passing verdict until the public reproduction is green at the integrated remediation revision. The existing five-package and six-lane green evidence remains accurate historical evidence, but it did not cover this legacy composition failure.
