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
- `internal/knowledgecommands/` owns BKB scaffold/status/structural scans, Dream's seven deterministic scans, Interview list/status/export/ingest/reset/version mechanics, and Memory's exact store plans, lexical recall, status/bootstrap/audit probes; actions retain semantic judgment, reports, locks, consent, consolidation, transcript summarization, optional semantic recall, and repair/recommendation choice.
- `internal/toolboxcommands/` owns notes, architecture preflight/publication, report-image generation, portfolio publication, last30days installation/checks, and typed maintainability audit metrics. Retained toolbox paths are compatibility launchers; the standalone audit-metrics module is retired.
- `internal/archivefetch/` owns in-process HTTP download/retry/redaction and archive transport fallback.
- `internal/requeststate/` owns deterministic `claim`, `unblock`, `complete`, `fail`, and `cancel` plans and their coupled checkpoint, archive, UR, calibration, and provenance mutations.
- `internal/publication/` is the sole deterministic owner of typed `capture-files`, `answer`, `release`, and `defer-gate` manifests, planning, containment, repository-root-confined parent handles, and atomic publication. Gate deferral binds exact parent/checkpoint preimages, repair identity, and diagnostic evidence; a fold is authorized by the exact pending repair preimage and tolerates an absent reservation after committed-request cleanup, while a present marker must still match. Actions retain gate execution, fingerprinting, revision attribution, path-drift judgment, scheduling, content, diagnosis, and semantic/release judgment, and once supplied these durable mutations have no prose fallback.
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
- [family: final-boundary-identity] Existing-file replacement must refuse symlinks/special files and detect observed identity or content changes; reusing a validated pathname at rollback or another final mutation boundary can target a replacement object, so bind or revalidate identity and complete mode metadata at the destructive syscall.
- [family: opaque-evidence-projection] A generic fallback or opaque aggregate can discard the exact blocker a caller must act on; derive output, specific typed records, and recovery argv from one observation set, and keep implementation dependencies in the Go standard library.
- Cleanup operation groups preflight independently; a dirty group is reported without blocking unrelated safe groups.
- Entirely untracked `Status: consumed` run scratch is the sole non-rollback deletion; revalidate its exact inventory immediately before removal.
- Doctor diagnosis is byte-for-byte read-only. Only `--repair-timestamps` mutates, and blank recovery remains exact cleanup consent.
- Queue selection is byte-for-byte read-only. It executes a scoped `blocked_check` only through the in-process owned process-group runner; every record retains exact request-path and raw probe/unblock evidence, while successful probes affect that invocation's eligibility but never rewrite the REQ.
- Screenshot and HTTP targets publish private 0600 bytes as the final commit point through rooted no-overwrite operations. No later error path removes a published pathname.
- Standalone `DownloadAtomic` stays rooted and no-overwrite. `FetchArchive` may refresh an initially regular target only after final validation proves its identity and content are unchanged.
- Hook protocol bytes are an optional typed result projection; JSON keeps the same output beside every typed finding and filesystem change. Reservation cleanup preserves every marker unless Git establishes committed request authority, including when Git itself is unavailable. New memory logs and ledgers use ordinary create permissions so the caller's umask retains the legacy filesystem effect. SessionStart launcher failures propagate actionably; memory Stop domain and launcher failures never block session end.
- Request-state commands validate selector-provided exact paths instead of selecting again. Actions retain confirmation, failure classification, terminal/review/release judgment, follow-up authoring, and dependent disposition; once supplied, deterministic lifecycle bytes have no free-form fallback.
- Gate deferral keeps parents `pending` with an explicit repair dependency. It never routes through `blocked` or `pending-answers`; repair creation/folding, claim removal, request movement, evidence, reservation, and rollback share one publication transaction. The work action consumes typed `gate_deferral` evidence to suppress the parent locally, widen the repair closure across URs, attribute late failures at a saved base, and validate rename-aware path drift before reuse. Repair REQs never recursively defer at the final gate: red/unverifiable repairs fail canonically while unrelated work continues; an already-green pre-build repair completes only through the durable reviewed no-op path and produces no release mutation.
- [family: interruptible-blocking-io] Reading confirmation synchronously and cancelling only the surrounding context → signal handling can wait forever for a caller that is still blocked on input; select cancellation against a buffered read result before writes, then leave post-write recovery with its transaction owner.
- Knowledge scans are read-only, normalized, sorted typed evidence. `bkb-init` never overwrites; Dream's optional newer-source probe is action-owned and is not an eighth scan.
- Interview templates remain data: the CLI renders only the declared mechanical dialect, publishes every artifact before stamping `last_exported_at`, and treats version archives as immutable. Interactive elicitation, approval, and contradiction resolution remain action-owned.
- Memory's tracked `working-memory.md` and private-untracked logs may share one transaction. Declared private targets are observed even when Git-ignored, rolled back through rooted identity checks, and never staged. Ledger appends remain best-effort; bootstrap has no committable target. Remember section/consolidation, forget selection, transcript summaries, semantic fusion, and the final audit verdict require explicit caller judgment.
- [family: smoke-vs-characterization] Registration and smoke checks can stay green while replacement semantics diverge; compare status, ordered evidence, actions, paths, and effects at authority boundaries before retiring an implementation. Toolbox compatibility Markdown and typed audit JSON come from one result. Image backends receive prompts only as argv, run inside owned Unix process groups, and are reaped before scratch cleanup; platforms without provable descendant ownership fail closed. Batch output appears only as one verified `generated/` publication, while an all-failed batch intentionally succeeds with no output.

## Verify

- Focused package: `go test ./internal/<package>`
- Static analysis: `go vet ./...`
- Module regression: `go test -count=1 ./...`
- Exact Go 1.25 compatibility: `bash _dev/tests/do-work-cli-go125-compatibility.sh` from the repository root
- Hook consumers: `bash _dev/tests/session-start-hook-behavior.sh && bash _dev/tests/memory-hook-behavior.sh`
- Knowledge commands: `go test -race ./internal/knowledgecommands` plus action/recipe contracts in `_dev/tests/contract-regressions.sh`
- Private memory transactions: `go test -race ./internal/gittransaction` and verify a mixed tracked/private `--commit` contains only `working-memory.md`.
- Windows atomic compile: `GOOS=windows GOARCH=amd64 go test -c ./internal/atomicfile -o <temporary-path>`
- Windows blocked-probe compile: `GOOS=windows GOARCH=amd64 go test -c ./internal/nextselection -o <temporary-path>`
- Toolbox migration: `go test -race ./internal/toolboxcommands`, exact old/new Markdown comparisons for all four audit modes, and `GOOS=windows GOARCH=amd64 go test -c ./internal/toolboxcommands -o <temporary-path>`.
- Repository baseline: run the unpiped `_dev/tests/maintainer-verify.sh` from the repository root when the integrating workflow calls for it.

## Stakes

- Req: lossless shared models and atomic mutation primitives.
- Value: later commands receive one typed snapshot with stable evidence and exact paths.
- Risk: parser drift or unsafe publication can corrupt durable request history or allocate duplicate IDs.

## Lessons

See [`lessons-do-work-cli.md`](lessons-do-work-cli.md) before changing the shared repository evidence contracts named above.
