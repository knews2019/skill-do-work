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
- `internal/publication/` is the sole deterministic owner of typed `capture-files`, `answer`, `release`, and `defer-gate` manifests, planning, containment, repository-root-confined parent handles, and atomic publication. Release admission requires an exhaustive exact normalized partition between action-judged `project_owned_targets` and maintainer-only `required_mirrors`; directory spelling and Git tracking are not ownership evidence. Gate deferral binds exact parent/checkpoint preimages, repair identity, and diagnostic evidence; a fold is authorized by the exact pending repair preimage and tolerates an absent reservation after committed-request cleanup, while a present marker must still match. Actions retain gate execution, fingerprinting, revision attribution, path-drift judgment, scheduling, content, diagnosis, and semantic/release judgment, and once supplied these durable mutations have no prose fallback.
- `internal/finalization/` composes request-state, optional release, exact commit, and provenance authorities behind a Git-private phase journal. `finalize --manifest` creates and advances one journal; `recover-finalization --discover` replays journals oldest-first, preserves unstaged protected rows without reading them, and admits only lifecycle/release-semantic unjournaled tails into the same phase engine. Legacy release discovery admits root project metadata, declared workspace members, and declared suite-maintainer topology affirmatively; it has no path-spelling exclusion fallback.
- `internal/gateevidence/` owns exact-argv, repository-bound green-gate records in Git-private state and validates revision ancestry before a baseline may be reused. `check-green-gate` targets `HEAD` by default and an explicit commit with `--at-revision`; the explicit target is what lets an already-green repair no-op keep verifying its own recorded revision after an unrelated commit moves `HEAD`.
- `internal/heavyverification/` owns the versioned heavy-lane manifest contract, rename-aware changed-path discovery, typed coverage rules, selected-lane explanations, and fail-closed force-all fallback. The strict planner reads target-committed manifest bytes. The separate historical-revalidation planner unions already-landed source ranges while reading the manifest from one exact descendant execution revision. The separate lane runner owns execution: it runs named manifest lanes at HEAD inside owned process groups, bounds each with a timeout, and records exit status, skip state, and wall seconds per lane. Per-lane result evidence still belongs to publication and selector resume validation.
- `internal/repairvalidation/` owns the read-only `validate-already-green-repair` authority shared by work and review. It strictly joins repair intake to no-op evidence, reuses `gateevidence` at the recorded past revision, observes NUL-safe Git state, and derives review staging from exact canonical-completion dry-run paths. Callers supply only the exact working request path plus the writer/time reused for completion; they never supply fingerprints, argv, revisions, or allowed paths.
- `internal/cleanup/` plans safe Passes 0–4, consent-gated repairs, link repointing, and worktree evidence.
- `internal/doctor/` owns read-only mechanical forensics and guarded blame-derived timestamp repair; recurring lesson judgment and board verification remain outside it.
- `internal/atomicfile/` owns safe existing-file replacement and exclusive marker creation.
- `internal/gittransaction/` owns dirty-target checks, rollback, staging, and commit guards.
- `internal/ownedprocess/` is the sole owner of owned-process-group launch and teardown for `gittransaction`, `toolboxcommands`, and `heavyverification`. `ConfigureGroup` reports whether group ownership was established rather than deciding what that means, so one API serves a caller that must fail closed and a caller that must degrade. `TerminateGroup` blocks until the group is gone, signalling descendants level-by-level before their parents. `nextselection`'s blocked probe keeps its own runner for its signal-forwarding and `128+signal` status contract.

## Package direction

`requestmodel` may import `schemanormalization`. `repositorymodel` may import `requestmodel` and `atomicfile`. `dependencygraph` consumes repository records. `gittransaction`, `toolboxcommands`, and `heavyverification` may import `ownedprocess`, which imports only the standard library and must never import a command package. Keep imports acyclic and do not import the separate `queue-kanban` command.

## Traps

- Preserve unknown frontmatter and body bytes unless a command explicitly owns the field being changed.
- Discover the do-work tree once; downstream logic consumes typed records and exact paths instead of rescanning.
- Reservation markers under `do-work/.req-reservations/` are durable coordination metadata and must be observed despite hidden-directory pruning.
- [family: final-boundary-identity] Every final mutation boundary must act on a filesystem object, never on a pathname that once resolved to one. Replacement must refuse symlinks/special files; creation must treat the successful exclusive create, not the plan, as the ownership event. Bind identity as inode **and** content digest — a delete-then-recreate in the same directory reuses the inode — then revalidate it, plus complete mode metadata, at the destructive syscall. Absence is not a replacement: distinguish "the object is gone" from "a different object stands there", or a completed removal is reported as a preserved replacement.
- [family: silent-skip-reads-as-red] A lane that skips silently inside its own test → the runner records a red exit and the drain routes the REQ to remediation for work that never ran; every engine-gated lane announces `SKIP:` on its own output, and the runner keys `skipped` on that prefix, never on a lane-name list.
- [family: closed-enumeration-for-a-condition] When a contract states a **condition**, a list of examples that happens to agree on those examples is a different contract. Test the condition's ingredients, taken from the governing spec's own definition, not the spellings today's cases use — and never replace one name list with a longer one. A partially-correct enumeration is the dangerous shape: the passing cases make the class look handled. Prefer an affirmative proof the repository itself declares (Git's index, a package marker) over inferring safety from the absence of known-bad names, and check both halves of a two-part proof separately — a differential that reddens on a rename proves nothing about the condition.
- [family: opaque-evidence-projection] A generic fallback or opaque aggregate can discard the exact blocker a caller must act on; derive output, specific typed records, and recovery argv from one observation set, and keep implementation dependencies in the Go standard library.
- Cleanup operation groups preflight independently; a dirty group is reported without blocking unrelated safe groups.
- Entirely untracked `Status: consumed` run scratch is the sole non-rollback deletion; revalidate its exact inventory immediately before removal.
- Doctor diagnosis is byte-for-byte read-only. Only `--repair-timestamps` mutates, and blank recovery remains exact cleanup consent.
- Queue selection is byte-for-byte read-only. It executes a scoped `blocked_check` only through the in-process owned process-group runner; every record retains exact request-path and raw probe/unblock evidence, while successful probes affect that invocation's eligibility but never rewrite the REQ.
- Screenshot and HTTP targets publish private 0600 bytes as the final commit point through rooted no-overwrite operations. No later error path removes a published pathname.
- Standalone `DownloadAtomic` stays rooted and no-overwrite. `FetchArchive` may refresh an initially regular target only after final validation proves its identity and content are unchanged.
- Hook protocol bytes are an optional typed result projection; JSON keeps the same output beside every typed finding and filesystem change. Reservation cleanup preserves every marker unless Git establishes committed request authority, including when Git itself is unavailable. New memory logs and ledgers use ordinary create permissions so the caller's umask retains the legacy filesystem effect. SessionStart launcher failures propagate actionably; memory Stop domain and launcher failures never block session end.
- Request-state commands validate selector-provided exact paths instead of selecting again. Actions retain confirmation, failure classification, terminal/review/release judgment, follow-up authoring, and dependent disposition; once supplied, deterministic lifecycle bytes have no free-form fallback.
- [family: alternate-writer-contract-drift] Changing a stored-format contract in its canonical lifecycle owner without sweeping cleanup, recovery, and other alternate writers leaves repository behavior split; pair the canonical regression with equivalent seam coverage at every remaining writer.
- [family: observed-subset-is-not-semantic-completeness] A coherent subset of dirty recovery paths is not proof of the whole transaction; require the configured semantic member set and exact ownership evidence for each admitted byte before finalizing legacy state.
- Finalization journals bind the manifest digest, lifecycle and release images, copied release payloads, exact commit allowlist, and phase commits before lifecycle mutation. Recovery may converge only bytes matching a recorded preimage or postimage; legacy discovery requires exact whole-file lifecycle/calibration/UR/follow-up/release evidence, and any foreign shared hunk, noncanonical journal identity, phase, or payload directory refuses before replay or cleanup.
- Gate deferral keeps parents `pending` with an explicit repair dependency. It never routes through `blocked` or `pending-answers`; repair creation/folding, claim removal, request movement, evidence, reservation, and rollback share one publication transaction. The work action consumes typed `gate_deferral` evidence to suppress the parent locally, widen the repair closure across URs, attribute late failures at a saved base, and validate rename-aware path drift before reuse. Repair REQs never recursively defer at the final gate: red/unverifiable repairs fail canonically while unrelated work continues; an already-green pre-build repair completes only through the durable reviewed no-op path and produces no release mutation.
- [family: reaped-by-its-own-parent] A process is reaped by its parent and by nobody else, so terminating a tree leader-first orphans the descendants to init and a killed orphan stays a zombie — and a zombie still satisfies `kill(pid, 0)`. Signal descendants before their parents, one level at a time, and wait for each level to be reaped before climbing; a wait on our own child ends in milliseconds, while a wait on init is measured in seconds. A group signal may only be sent once the group is proved isolated; otherwise signal the bare pid.
- [family: interruptible-blocking-io] Reading confirmation synchronously and cancelling only the surrounding context → signal handling can wait forever for a caller that is still blocked on input; select cancellation against a buffered read result before writes, then leave post-write recovery with its transaction owner.
- Knowledge scans are read-only, normalized, sorted typed evidence. `bkb-init` never overwrites; Dream's optional newer-source probe is action-owned and is not an eighth scan.
- Interview templates remain data: the CLI renders only the declared mechanical dialect, publishes every artifact before stamping `last_exported_at`, and treats version archives as immutable. Interactive elicitation, approval, and contradiction resolution remain action-owned.
- Memory's tracked `working-memory.md` and private-untracked logs may share one transaction. Declared private targets are observed even when Git-ignored, rolled back through rooted identity checks, and never staged. Ledger appends remain best-effort; bootstrap has no committable target. Remember section/consolidation, forget selection, transcript summaries, semantic fusion, and the final audit verdict require explicit caller judgment.
- [family: smoke-vs-characterization] Registration and smoke checks can stay green while replacement semantics diverge; compare status, ordered evidence, actions, paths, and effects at authority boundaries before retiring an implementation. Toolbox compatibility Markdown and typed audit JSON come from one result. Image backends receive prompts only as argv, run inside owned Unix process groups, and are reaped before scratch cleanup; platforms without provable descendant ownership fail closed. Batch output appears only as one verified `generated/` publication, while an all-failed batch intentionally succeeds with no output.
- [family: collision-fixture-identity] Reusing a target number in both fixture filenames and frontmatter → filename claims can make an incomplete frontmatter identity rule look correct; use unrelated filename numbers plus a malformed adjacent-value negative control.
- [family: projection-before-bounding] Applying a dispatch bound before projecting a frozen target ledger → newly discovered out-of-ledger records consume selector slots and starve retained work; observe canonically without the scheduling bound, project frozen membership, then bound dispatch.
- [family: publication-target-topology-classification] Treating a tracked path or a manifest's self-description as release ownership → dependency and generated trees can authorize their own metadata; seed ownership at an independent repository or declared maintainer root, then propagate only through proven relationships.

## Verify

- Focused package: `go test ./internal/<package>`
- Green-gate evidence: `go test ./internal/gateevidence`
- Already-green repair authority: `go test ./internal/repairvalidation ./internal/gateevidence ./internal/resultmodel`
- Static analysis: `go vet ./...`
- Module regression: `go test -count=1 ./...`
- Hook consumers: `bash _dev/tests/session-start-hook-behavior.sh`
- Knowledge commands: `go test -race ./internal/knowledgecommands` plus action/recipe contracts in `_dev/tests/contract-regressions.sh`
- Private memory transactions: `go test -race ./internal/gittransaction` and verify a mixed tracked/private `--commit` contains only `working-memory.md`.
- Windows atomic compile: `GOOS=windows GOARCH=amd64 go test -c ./internal/atomicfile -o <temporary-path>`
- Windows owned-process compile: `GOOS=windows GOARCH=amd64 go build ./...` plus `GOOS=windows GOARCH=amd64 go vet ./internal/ownedprocess ./internal/gittransaction ./internal/toolboxcommands ./internal/heavyverification`. `ownedprocess` is the only build-tagged split these four packages depend on, and `gittransaction` must keep degrading to default cancellation there rather than failing closed.
- Windows blocked-probe compile: `GOOS=windows GOARCH=amd64 go test -c ./internal/nextselection -o <temporary-path>`
- Toolbox migration: `go test -race ./internal/toolboxcommands`, exact old/new Markdown comparisons for all four audit modes, and `GOOS=windows GOARCH=amd64 go test -c ./internal/toolboxcommands -o <temporary-path>`.
- Repository baseline: run the unpiped `_dev/tests/maintainer-verify.sh` from the repository root when the integrating workflow calls for it.

## Stakes

- Req: lossless shared models and atomic mutation primitives.
- Value: later commands receive one typed snapshot with stable evidence and exact paths.
- Risk: parser drift or unsafe publication can corrupt durable request history or allocate duplicate IDs.

## Lessons

See [`lessons-do-work-cli.md`](lessons-do-work-cli.md) before changing the shared repository evidence contracts named above.
