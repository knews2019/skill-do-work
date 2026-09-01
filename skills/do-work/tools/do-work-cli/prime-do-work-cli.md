# do-work-cli Prime

## Purpose

The standard-library Go module under this directory is the canonical implementation for deterministic do-work operations. Natural-language actions and retained shell paths delegate to it; `queue-kanban` remains a separate board binary.

## Read first

- `go.mod` defines the module and Go floor.
- `cmd/do-work-cli/main.go` is the command registration boundary. Shared packages must not register commands.
- `internal/resultmodel/` and `internal/commandruntime/` own stable command results and rendering.
- `internal/requestmodel/` owns lossless REQ/UR frontmatter documents and authorized field edits.
- `internal/schemanormalization/` owns schema aliases, defaults, warnings, and terminal predicates.
- `internal/repositorymodel/` owns one-pass do-work discovery, exact paths, collisions, and REQ allocation.
- `internal/dependencygraph/` derives readiness, reverse edges, cycles, and depth from a repository snapshot.
- `internal/nextselection/` owns read-only target expansion, queue readiness, process-tree-owned blocked probes, wave/fan-out bounds, estimates, and typed selected/excluded records.
- `internal/corehelpers/` owns the remaining utility handlers and leaf check, inventory, Git, publication, reservation, and survey mechanics; shared download, timestamp, probe, and provenance primitives remain with their domain owners.
- `internal/hookcommands/` owns core and memory SessionStart plus memory Stop protocols; retained hook scripts only launch these commands.
- `internal/archivefetch/` owns in-process HTTP download/retry/redaction and archive transport fallback.
- `internal/requeststate/` owns deterministic `claim`, `unblock`, `complete`, `fail`, and `cancel` plans and their coupled checkpoint, archive, UR, calibration, and provenance mutations.
- `internal/publication/` is the sole deterministic owner of typed `capture-files`, `answer`, and `release` manifests, planning, containment, repository-root-confined parent handles, and atomic publication. Stakeholder history and overrides are typed, file-backed evidence; actions retain content and semantic/release judgment, and once supplied these durable mutations have no prose fallback.
- `internal/cleanup/` plans safe Passes 0–4, consent-gated repairs, link repointing, and worktree evidence.
- `internal/doctor/` owns read-only mechanical forensics and guarded blame-derived timestamp repair; recurring lesson judgment and board verification remain outside it.
- `internal/atomicfile/` owns safe existing-file replacement and exclusive marker creation.
- `internal/gittransaction/` owns dirty-target checks, rollback, staging, and commit guards.

## Package direction

`requestmodel` may import `schemanormalization`. `repositorymodel` may import `requestmodel` and `atomicfile`. `dependencygraph` consumes repository records. Keep imports acyclic and do not import the separate `queue-kanban` command.

## Traps

- Preserve unknown frontmatter and body bytes unless a command explicitly owns the field being changed.
- Discover the do-work tree once; downstream logic consumes typed records and exact paths instead of rescanning.
- Reservation markers under `do-work/.req-reservations/` are durable coordination metadata and must be observed despite hidden-directory pruning.
- Existing-file replacement refuses symlinks and special files, and detects identity or content changes observed during pre-publish validation. It guarantees complete atomic publication, not portable compare-and-swap against arbitrary non-cooperating writers after the final check.
- Keep command output derived from one typed result model and keep implementation dependencies in the Go standard library.
- Cleanup operation groups preflight independently; a dirty group is reported without blocking unrelated safe groups.
- Entirely untracked `Status: consumed` run scratch is the sole non-rollback deletion; revalidate its exact inventory immediately before removal.
- Doctor diagnosis is byte-for-byte read-only. Only `--repair-timestamps` mutates, and blank recovery remains exact cleanup consent.
- Queue selection is byte-for-byte read-only. It executes a scoped `blocked_check` only through the in-process owned process-group runner; every record retains exact request-path and raw probe/unblock evidence, while successful probes affect that invocation's eligibility but never rewrite the REQ.
- Screenshot and HTTP targets publish private 0600 bytes as the final commit point through rooted no-overwrite operations. No later error path removes a published pathname.
- Standalone `DownloadAtomic` stays rooted and no-overwrite. `FetchArchive` may refresh an initially regular target only after final validation proves its identity and content are unchanged.
- Hook protocol bytes are an optional typed result projection. SessionStart launcher failures propagate actionably; memory Stop domain and launcher failures never block session end.
- Request-state commands validate selector-provided exact paths instead of selecting again. Actions retain confirmation, failure classification, terminal/review/release judgment, follow-up authoring, and dependent disposition; once supplied, deterministic lifecycle bytes have no free-form fallback.

## Verify

- Focused package: `go test ./internal/<package>`
- Static analysis: `go vet ./...`
- Module regression: `go test -count=1 ./...`
- Exact Go 1.25 compatibility: `bash _dev/tests/do-work-cli-go125-compatibility.sh` from the repository root
- Hook consumers: `bash _dev/tests/session-start-hook-behavior.sh && bash _dev/tests/memory-hook-behavior.sh`
- Windows atomic compile: `GOOS=windows GOARCH=amd64 go test -c ./internal/atomicfile -o <temporary-path>`
- Windows blocked-probe compile: `GOOS=windows GOARCH=amd64 go test -c ./internal/nextselection -o <temporary-path>`
- Repository baseline: run the unpiped `_dev/tests/maintainer-verify.sh` from the repository root when the integrating workflow calls for it.

## Stakes

- Req: lossless shared models and atomic mutation primitives.
- Value: later commands receive one typed snapshot with stable evidence and exact paths.
- Risk: parser drift or unsafe publication can corrupt durable request history or allocate duplicate IDs.

## Lessons

See [`lessons-do-work-cli.md`](lessons-do-work-cli.md) before changing the shared repository evidence contracts named above.
