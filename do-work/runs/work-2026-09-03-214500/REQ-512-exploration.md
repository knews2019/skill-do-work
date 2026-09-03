# REQ-512 Exploration — Complete Legacy Finalization Semantic Ownership

## Scope and evidence read

This exploration is limited to REQ-512 and is read-only with respect to source and lifecycle state. It is based on:

- the authoritative working request, `REQ-512-complete-legacy-finalization-semantic-ownership.md`;
- the run brief;
- the do-work CLI prime and lessons;
- the general, coding-guardrails, backend, testing, and communication-style crew guidance;
- the REQ-498 and REQ-499 archived recovery/finalization requests;
- the current finalization discovery, command, and apply implementations; and
- the current REQ-499, recovery, git-transaction, request-state, and public recovery tests.

No active action writes a literal `Review Fold` or `Recovery Fold` in this repository. Those headings are legacy evidence that finalization discovery may accept, so their acceptance grammar must itself carry a closed ownership boundary.

## Root cause

REQ-512 has two independent semantic-ownership gaps in `internal/finalization/finalization_discovery.go`.

### Tracked follow-up ownership is open-ended

`lifecyclePathProves` delegates tracked addendum/re-review files to `followupPathProves`. That function correctly reads the durable HEAD preimage and requires the working-tree file to retain it byte-for-byte, but it treats the entire suffix as owned when the suffix:

1. begins with a blank line and a `## Review Fold — REQ-…` or `## Recovery Fold — REQ-…` heading; and
2. contains no second top-level `##` heading.

The regular expression is anchored at EOF only after an unrestricted `[\\s\\S]*`. Consequently an unheaded paragraph, HTML comment, malformed heading, thematic-break delimiter, or delimiter-shaped addition after the intended fold remains inside the regex match and is attributed to the finalizer. Counting only `^##` cannot prove where the named fold ended. The durable preimage proves the start of the append, but the legacy format has no end boundary.

The untracked-file branch is separate: a newly created, structurally valid addendum is owned as a whole through its declared relationship. REQ-512's tracked-append requirement should not weaken or silently broaden that rule.

### Workspace selection starts from coincidental old-version values

`associateReleaseMetadata` first derives the one `oldVersion -> newVersion` transition from dirty release metadata. It then calls `configuredReleaseMetadataPaths`, which enumerates tracked release files and builds the affirmative ownership graph. The selection stage nevertheless scans every owned manifest and lock and admits files because their HEAD semantic version equals `oldVersion`:

- package, Cargo, and Python manifests can be selected even when they did not change;
- Cargo and uv lock mirrors can be selected from a root declaration before a changed source is identified;
- npm `npmRootVersionCopies` counts top-level `version` and `packages[""].version` solely because the clean root lock happens to contain the release's old version.

`workspaceMirrorReplacement` later enforces an exact replacement over that over-broad selection. For a member-only release, this makes an unchanged root package and unchanged root lock copies mandatory participants simply because the root shares the member's old version. The equivalent ownership error exists for Cargo and uv: topology and value coincidence identify a purported release source before the dirty manifest that actually changed `old -> new` has been established.

The causal direction must be reversed:

`changed, owned source manifest -> declared workspace relation -> that source's exact lock member mirror`

not:

`any owned file containing oldVersion -> presumed release participant`.

## Existing ownership and workspace ecosystem inventory

### Repository and suite ownership

`releaseMetadataPath` recognizes changelogs, `VERSION`, npm, Cargo, uv/Python, and the exact do-work action version path. `affirmativeReleaseOwnership` owns the repository root, then seeds suite roots only from tracked `suite/modules.tsv` rows whose source has a tracked `VERSION`. Nested workspaces become eligible only when reached from an already-owned parent through `findOwnedWorkspaceOwner`. This affirmative, recursive topology is the correct boundary and should be retained.

`enumerateTrackedReleasePaths` is an injectable `git ls-files` seam. Enumeration failure already maps to the typed `FINALIZATION-DISCOVERY-RELEASE-ENUMERATION` refusal and occurs before a finalization journal or commit. New selection must continue to use this complete tracked inventory and fail closed on every enumeration/read/parse ambiguity.

### npm

- Source identity: `package.json` `name` and `version`.
- Workspace declaration: `workspaces: []` or `workspaces: { packages: [] }`.
- Mirror: `package-lock.json`.
- Member mirror: `packages["<member-path>"].name/version`, with the path relative to the workspace root.
- Root mirrors: lock top-level `version` and `packages[""].version`.
- Required REQ-512 behavior: a changed member source requires only that member entry. Root lock copies are required only when the root `package.json` itself is one of the proven changed sources. Their coincidental equality to the member's old version is not evidence of participation.

### Cargo

- Source identity: `Cargo.toml` `[package].name/version`.
- Workspace declaration: `[workspace].members`.
- Mirror: `Cargo.lock`.
- Member mirror: the unique local `[[package]]` block with matching `name/version` and no registry/git `source`.
- Required REQ-512 behavior: a changed member requires its local package block. A clean parent manifest is topology, not a release source, even when its own package version equals the old version. Existing ambiguity refusal for duplicate/missing eligible blocks remains fail closed.

### uv / Python

- Source identity: `pyproject.toml` `[project]` or `[tool.poetry]` `name/version`.
- Workspace declaration: `[tool.uv.workspace].members`.
- Mirror: `uv.lock`.
- Member mirror: the matching package with `source.editable` or `source.virtual` equal to the member-relative path; the workspace root is represented by `.`.
- Required REQ-512 behavior: only a changed project source contributes a required lock descriptor. A clean root project is not selected by an old-version match.

### Non-workspace release metadata

Repository/suite `VERSION` files, changelogs, the do-work action version, and other existing configured release mirrors remain governed by their current complete-mirror rules. REQ-512 must not turn changed-source-first workspace selection into a general permission for partial releases or arbitrary release metadata. Dirty recognized metadata that is neither an established configured mirror nor derived from a changed owned source must still be refused.

## Bounded tracked-fold format

Replace the open-ended tracked-suffix regex with a small parser for a closed legacy envelope. The accepted suffix should have:

1. the exact append boundary after the durable preimage;
2. one `## Review Fold — <request-id>…` or `## Recovery Fold — <request-id>…` opening heading whose kind and request ID match the document;
3. body bytes; and
4. one exact terminal line, for example:

```markdown
<!-- do-work:finalization-followup-fold-end kind=review request=REQ-754 -->
```

The matching terminal marker must be the final semantic line (allowing only the parser's explicitly chosen canonical final newline), occur exactly once, and match both fold kind and request ID. The parser must reject bytes before the opening heading, bytes after the terminal marker, a missing/duplicate/mismatched marker, another top-level fold heading, and delimiter-like material outside the envelope. Review and recovery are enumerated values; unknown kinds fail closed.

This is an equivalently bounded format: the durable HEAD preimage fixes the append's start and the authenticated-by-grammar terminal line fixes its end. Content deliberately placed between the two markers remains part of the declared fold; this format cannot attribute concurrent writers within one owner-declared envelope, so the complete envelope must be appended as one finalization operation. Existing unbounded tracked folds will safely refuse. `--assume-sole-releaser` must not bypass this proof.

## Requirement-mapped implementation plan

### R1 — Prove the exact whole tracked append

- Add a focused helper, such as `boundedFollowupFoldProves(appended, kind, requestID)`, and have `followupPathProves` use it after the durable prefix comparison.
- Parse exact line boundaries rather than using a permissive whole-suffix regex.
- Keep the document relationship checks already performed by `parseRequestDocument`.
- Leave whole-document ownership for a structurally valid untracked follow-up unchanged.

### R2 — Reject foreign bytes around the fold

- Require the first suffix bytes to be the canonical append newline plus the matching fold heading.
- Require exactly one matching terminal marker at EOF and exactly one opening fold.
- Reject unheaded paragraphs, HTML comments, malformed headings, thematic breaks, duplicate/end-marker-shaped additions, and a second fold when any occur outside the one closed envelope.
- Preserve refusal as `FINALIZATION-DISCOVERY-AMBIGUOUS`, with the follow-up path blocked and with HEAD, index, working-tree bytes, and journal state unchanged.

### R3 — Identify changed workspace sources before mirrors

- While `associateReleaseMetadata` reads dirty release paths, retain structured manifest transitions, not just the global sets of old and new version strings. A proposed `releaseVersionChange` value should carry at least `path`, `oldVersion`, and `newVersion`.
- Admit a workspace source only when it is dirty, affirmatively owned, parseable in its ecosystem, and shows the exact common `old -> new` transition.
- Pass this changed-source set into configured-path construction. Do not search clean manifests for `oldVersion` as a selection rule.

### R4 — Derive only the changed source's declared mirrors

- Separate topology discovery from semantic admission. `affirmativeReleaseOwnership`, workspace membership parsing, and nearest-owner resolution remain topology-only.
- For each proven changed source:
  - include the source manifest;
  - derive its same-root lock descriptor when applicable;
  - find its nearest affirmatively owned workspace owner;
  - add only that source's member descriptor to the owner's lock;
  - merge descriptors deterministically when several changed members share a lock.
- Preserve exact byte-level mirror replacement checks after the narrower semantic set is constructed.

### R5 — Make root lock copies conditional on a changed root source

- Replace the current `workspaceReleaseMirror.rootCopies` value inferred from lock coincidence with explicit changed-source intent, for example `rootSourceChanged bool` plus a count derived only for that changed root.
- npm member-only selection must expect `packages["<member-path>"].version` and leave top-level `version` and `packages[""].version` byte-identical.
- If the root `package.json` changed, continue to require all structurally present root lock copies to change exactly.
- Apply the same source-first rule to Cargo root/local blocks and uv root `.` descriptors.

### R6 — Remain typed and fail closed

- A dirty workspace manifest or lock not admitted by the changed-source derivation must not pass merely because generic semantic replacement succeeds.
- Preserve typed ownership, enumeration, read, parse, ambiguity, and exact-mirror refusals. No fallback may silently skip an unreadable path or unknown ecosystem.
- Keep sorting/deduplication deterministic through the existing unique/sorted path and mirror mechanisms.
- Keep strict discovery strict; `--assume-sole-releaser` may widen only the existing named shared-metadata cases and must not widen follow-up or release ownership.

## Exact proposed builder write set

Only these files should be writable for the implementation:

1. `skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go`
2. `skills/do-work/tools/do-work-cli/internal/finalization/finalization_req499_test.go`
3. `skills/do-work/tools/do-work-cli/internal/finalization/finalization_req512_test.go` (new)

The REQ-499 test file is included only because it owns the current positive tracked-fold fixtures and the direct `followupPathProves` table; those fixtures must move to the bounded format rather than preserving the vulnerable grammar. REQ-512-specific rejection and three-ecosystem workspace matrices belong in the new test file.

No action/contract documents, command/result schema, other Go packages, generated sites, screenshots, lifecycle/request/queue records, run manifests, release/version/changelog files, or unrelated tests are in the builder write set. If implementation proves that an active fold producer exists outside the inspected code or another production file must change, stop and return scope drift for an orchestrator decision instead of widening this set.

## RED and GREEN evidence plan

### Tracked folds

First preserve RED evidence showing current unsafe acceptance:

- a tracked addendum whose suffix has an otherwise accepted review fold followed by an unheaded foreign paragraph finalizes today; assert the desired result is `OutcomeRefused`, code `FINALIZATION-DISCOVERY-AMBIGUOUS`, and the follow-up path is blocked;
- repeat at the predicate/table level for foreign bytes before the heading and for material after the intended fold: an HTML comment, malformed `#`/`###` heading, `---`, `***`, duplicate or mismatched terminal marker, terminal-marker-shaped addition, second fold, and missing marker.

Update the existing exact tracked-fold recovery positive and direct predicate positives to include the canonical terminal marker. Add matching bounded Review and Recovery positives. GREEN must also assert refusal atomicity: unchanged working-tree bytes, index/status, HEAD, and absence of a newly created finalization journal.

### Workspace matrices

Create table-driven npm, Cargo, and uv fixtures with no unrelated suite noise. Each fixture should contain:

- one terminal legacy finalization candidate with release metadata sufficient for discovery;
- an owned workspace root whose own version is `1.0.0` and whose declaration includes one member;
- a member source changed `1.0.0 -> 1.0.1`;
- an unchanged root source that remains `1.0.0`; and
- a lock changed only in the member's exact mirror. For npm, top-level `version` and `packages[""].version` remain `1.0.0`.

The primary RED is that strict `--discover` currently refuses because the unchanged root is selected. The GREEN expectation is cleanup completion, with the committed release paths including the member and lock, excluding the unchanged root, and with root bytes exact.

For each ecosystem add:

- a stale member lock mirror negative, which must refuse and leave all bytes/state unchanged;
- a changed-root control, which proves root lock copies/root descriptors are still mandatory when the root source actually changed; and
- where the fixture permits it, multiple changed members sharing one lock to prove deterministic descriptor aggregation without selecting clean siblings.

The existing release-enumeration error test must continue to return `FINALIZATION-DISCOVERY-RELEASE-ENUMERATION` before mutation. Retain or extend tests for an unowned workspace path, ambiguous Cargo local block, invalid uv source identity, and dirty source-less lock metadata as typed fail-closed cases.

### Regression and verification commands

The builder should record exact RED commands/output before production edits, then run:

```sh
go test -count=1 ./internal/finalization
go test -race -count=1 ./internal/finalization ./internal/gittransaction ./internal/requeststate ./internal/publication
go vet ./...
go test -count=1 ./...
DO_WORK_HEAVY_TESTS=1 go test -count=1 ./internal/finalization -run 'TestPublicRecoverFinalizationMovesURThenAllowsRealClaim|TestRecoverFinalization'
```

Run these from `skills/do-work/tools/do-work-cli`. Include focused runs for every new REQ-512 test name in the handback, retain the existing protected-path and unrelated-change preservation coverage, and use the repository's Go 1.25 compatibility/toolchain check if it is separately configured. This is backend discovery behavior with no browser surface, so browser verification is not applicable.

## Risks and required review focus

- **Legacy compatibility:** previously accepted unbounded tracked folds will refuse. This is an intentional safety tightening required by REQ-512; there is no active repository producer to migrate in the proposed write set.
- **Envelope assumption:** bytes deliberately inserted inside a complete owner-declared envelope cannot be separated by syntax alone. The terminal marker proves the append boundary, so producers must append the complete envelope atomically.
- **Nested workspaces:** source-first selection must still use the existing nearest affirmatively owned workspace relation and must not jump to a merely adjacent manifest.
- **Shared locks:** several changed members may map to one lock. Descriptor merging must remain stable and exact, with no map-order-dependent output.
- **Root/member overlap:** a workspace root can also be a package. Root mirror requirements must arise from that root source's actual dirty transition, not from its role as workspace owner.
- **False permission through generic metadata:** narrowing configured paths must not let an independently dirty lock or recognized manifest evade ownership checking. Every dirty release path still needs an explicit semantic provenance.
- **Atomic refusal:** all new ambiguity and mirror failures must occur before journal creation, staging, or commits.
- **Assume mode:** `--assume-sole-releaser` is not a recovery hatch for either gap.

## Integration and orchestrator notes

REQ-512 changes the proof model in finalization discovery but does not change CLI arguments or public result fields. It should integrate as one coherent finalization commit after its focused/full checks and diff review. The orchestrator retains ownership of the lifecycle transition, run manifest/write-set approval, queue stamps, cancellation/supersession decisions, merge/cherry-pick, cross-REQ conflict handling, release notes/versioning, and post-merge public recovery-to-claim verification.

The public acceptance seam is the existing real recovery flow: a qualifying UR is moved, finalization completes, and a subsequent real claim succeeds. Run that heavy test after the focused semantic matrix so the narrower path selection is proven not to strand the queue or weaken protected-path refusal.
