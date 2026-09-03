# REQ-512 Builder Handback

## Commit

- Branch: `codex/REQ-512-finalization-ownership`
- Commit: `54023d91063fb464eba028c07569a8377dd935a3`
- Message: `[REQ-512] complete finalization semantic ownership`
- Worktree after commit: clean

## Changed files

The commit contains exactly the declared three-file write set:

1. `skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go`
2. `skills/do-work/tools/do-work-cli/internal/finalization/finalization_req499_test.go`
3. `skills/do-work/tools/do-work-cli/internal/finalization/finalization_req512_test.go` (new)

No lifecycle/request state, run manifest, queue record, action/contract document, generated site, version, changelog, release file, or out-of-scope source was edited.

## Implemented behavior

- Tracked legacy follow-up ownership now requires a closed fold envelope rooted at the durable HEAD preimage. The opening `Review`/`Recovery` kind and request ID must match an exact terminal marker at EOF. Missing, duplicate, mismatched, or followed markers refuse, as do foreign bytes before the opening and a second fold.
- Dirty manifest transitions are recorded before workspace mirror selection. Clean equal-version package, Cargo, and uv roots are topology only and are not release participants.
- A proven changed source derives its same-directory root lock mirrors and its declared nearest workspace-member mirror. Multiple descriptors are deduplicated and assembled in sorted source-path order.
- npm root `version` and `packages[""]` copies are required only when the root `package.json` itself changed. A member-only release requires only the member package entry.
- Cargo and uv use the same changed-source-first rule, while preserving exact local/source-qualified lock-block checks.
- Dirty release metadata without a selected changed-source/configured-mirror provenance refuses. Enumeration remains typed and fail closed.
- Existing strict, protected-path, unrelated-work, assume-mode, and public recovery-to-claim behavior remains green.

## Literal RED evidence

Command, before production edits:

```sh
go test -count=1 ./internal/finalization -run 'TestREQ512(TrackedFoldRequiresClosedBoundary|WorkspaceMembersSelectChangedSourcesBeforeEqualVersionRoots)' -v
```

Exit: `1`.

Relevant failing assertions:

```text
=== RUN   TestREQ512TrackedFoldRequiresClosedBoundary/unheaded_foreign_tail_has_no_owned_boundary
    finalization_req512_test.go:36: followupPathProves() = true, want false
```

The changed-source-first matrix failed in every ecosystem with `OutcomeRefused`, `FINALIZATION-DISCOVERY-AMBIGUOUS`, and the clean equal-version root as the blocked path:

```text
npm:   skills/do-work/consumer/package.json
cargo: skills/do-work/consumer-rust/Cargo.toml
uv:    skills/do-work/consumer-python/pyproject.toml
```

The npm assertion began:

```text
finalization_req512_test.go:101: changed-source-first workspace recovery = resultmodel.CommandResult{... Outcome:"refused" ... BlockedPaths:[]string{"skills/do-work/consumer/package.json"} ... ReasonCodes:[]string{"FINALIZATION-DISCOVERY-AMBIGUOUS"} ...}
```

Cargo and uv failed at the same assertion with their corresponding root paths. The package result was `FAIL`.

## GREEN evidence

Focused REQ-512 command:

```sh
go test -count=1 ./internal/finalization -run 'TestREQ512' -v
```

Result:

```text
--- PASS: TestREQ512TrackedFoldRequiresClosedBoundary
--- PASS: TestREQ512TrackedFoldForeignTailRefusesWithoutMutation
--- PASS: TestREQ512ChangedWorkspaceRootsStillRequireRootLockMirrors
--- PASS: TestREQ512WorkspaceMembersSelectChangedSourcesBeforeEqualVersionRoots
PASS
ok github.com/knews2019/skill-do-work/do-work-cli/internal/finalization 9.333s
```

The fold matrix covers bounded review/recovery positives and rejects foreign prefix bytes, unheaded suffix bytes, comments, malformed headings, delimiter-shaped suffixes, mismatched/duplicate markers, and a second fold. The end-to-end refusal test proves follow-up bytes, Git status, HEAD, and journal absence remain unchanged.

The workspace matrices cover npm, Cargo, and uv member-only success with clean equal-version roots excluded byte-identically; existing stale-lock cases remain green for all three; changed-root controls prove root lock copies/blocks are still required.

Existing focused recovery/ownership regression:

```sh
go test -count=1 ./internal/finalization -run 'TestRecoverFinalizationRequiresWorkspaceMemberLockMirrors|TestREQ512' -v
```

Result: `PASS`, package `ok` in `13.354s`, including stale and updated npm/Cargo/uv cases.

Full finalization package:

```sh
go test -count=1 ./internal/finalization
```

Result:

```text
ok github.com/knews2019/skill-do-work/do-work-cli/internal/finalization 40.430s
```

Race gate:

```sh
go test -race -count=1 ./internal/finalization ./internal/gittransaction ./internal/requeststate ./internal/publication
```

Result: all four packages `ok`; finalization completed in `47.334s`.

Vet and full module:

```sh
go vet ./...
go test -count=1 ./...
```

Result: vet exit `0`; every module package passed, including finalization in `55.848s` and publication in `32.872s`.

Go 1.25 compatibility:

```sh
GOTOOLCHAIN=go1.25.0 go test -count=1 ./internal/finalization -run 'TestREQ512'
```

Result:

```text
ok github.com/knews2019/skill-do-work/do-work-cli/internal/finalization 8.618s
```

Heavy public recovery and strict protections:

```sh
DO_WORK_HEAVY_TESTS=1 go test -count=1 ./internal/finalization -run 'TestPublicRecoverFinalizationMovesURThenAllowsRealClaim|TestRecoverFinalizationAssumeSoleReleaserStillRefusesStagedProtectedPath|TestRecoverFinalizationAssumeSoleReleaserAttributesOnlySharedMetadata' -v
```

Result: all three tests passed; package `ok` in `3.366s`.

Diff checks:

```sh
git diff --check
git diff --cached --check
git diff --cached --stat
```

Result: no whitespace errors; staged stat was three files, 315 insertions and 51 deletions. Every changed file was reviewed before commit. No browser check applies to this backend-only CLI discovery change.

## Scope notes and decisions

- The existing REQ-499 npm workspace fixture was adjusted within its already declared test file. Its setup had accidentally left a dirty root `package-lock.json` after committing a versionless root source, which is source-less lock mutation under the new explicit provenance rule. The fixture now commits a clean equal-version root manifest and root lock before exercising the member lock. This is an intentional test-contract correction, not a production exception.
- Existing unbounded tracked folds now refuse. No active repository action writes these legacy headings; any manual/external producer that expects strict recovery must append the complete opening/body/terminal-marker envelope atomically.
- `--assume-sole-releaser` does not widen either follow-up ownership or release ownership.
- No scope drift or unresolved blocker was encountered.

## Risks and merge guidance

- The bounded marker is a deliberate compatibility tightening. A previously accepted tracked fold without the terminal marker is now ambiguous by design.
- The syntax bounds bytes around the declared envelope; as with any owner-declared envelope, it cannot distinguish a concurrent writer who inserts bytes inside a completely formed envelope. Producers must write the complete append as one operation.
- Source-less dirty locks now refuse rather than being admitted by old-version coincidence. This is the fail-closed behavior REQ-512 requires.
- Cherry-pick `54023d91063fb464eba028c07569a8377dd935a3` onto the orchestrator integration branch. Re-run the focused REQ-512 test and the public recovery-to-claim test after merge. Lifecycle completion, queue stamps, cancellation/supersession handling, merge bookkeeping, and any release/version/changelog work remain orchestrator-owned.

## Discovered follow-up work

None required for REQ-512 acceptance. If a future action begins producing tracked review/recovery folds, its contract should reuse the exact terminal-marker grammar added here rather than inventing a second format.
