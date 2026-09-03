# REQ-485 exploration — canonical REQ reservation marker filenames

## Recommendation

Use the stored REQ-id spelling as the canonical reservation basename: `REQ-` plus the numeric request value at the repository's normal minimum width of three digits, with no additional fixed-width padding.

Examples:

| Number | Canonical new marker | Legacy spelling still read/reaped |
|---:|---|---|
| 1 | `REQ-001` | `REQ-000001` (and numeric aliases already present) |
| 42 | `REQ-042` | `REQ-000042` |
| 482 | `REQ-482` | `REQ-000482` |
| 1207 | `REQ-1207` | `REQ-001207` |

This is the least surprising contract because canonical stored IDs are minimum-three-digit (`actions/capture.md` starts at `REQ-001`; `actions/work-reference.md` says stored IDs are zero-padded and gives `REQ-067`), capture and defer-gate already normally carry the stored ID in `reservation_path`, and the incident's desired shared path is `REQ-482`. Fixed six-digit marker names are isolated to the two allocators, cleanup's narrow regex, and two guidance sentences. New writers should never emit the six-digit alias; readers should compare the parsed positive integer and accept every all-digit legacy width so both known spellings remain safe.

Do not rename existing markers in place. Migration is read-compatible and naturally drains them through the existing committed-REQ/48-hour cleanup policy.

## Current defect, reproduced from source

- `skills/do-work-board/tools/queue-kanban/allocate.go` → `requestReservationFileName` returns `fmt.Sprintf("REQ-%06d", requestNumber)`. `nextRequestNumber` exclusively creates that path.
- `skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model.go` → `ReserveNextRequestID` independently formats `REQ-%06d`, so it is a second allocator writer with the same stale convention even though it currently has no non-test call site.
- `skills/do-work/tools/do-work-cli/internal/publication/capture_files.go` → `BuildCapturePlan` requires `reservation_path == "do-work/.req-reservations/" + request.ID` and creates exactly that path. For `REQ-482`, this is `REQ-482`.
- `skills/do-work/tools/do-work-cli/internal/publication/defer_gate.go` → `validateDeferGateManifest` and the create branch of `BuildDeferGatePlan` use the same exact-ID spelling for repair markers. `answer` is also a writer through `BuildAnswerPlan`'s delegation of `override_capture` to `BuildCapturePlan`.
- Therefore a board reservation at `REQ-000482` and capture reservation at `REQ-482` are distinct `O_EXCL` targets and both succeed.
- `skills/do-work/tools/do-work-cli/internal/corehelpers/reservations.go` compounds the split: `reservationNamePattern` is `^REQ-(\d{6})$`, so current unpadded capture/defer markers are reported as malformed and are never reaped.

## Writer seams

### Board allocator

File: `skills/do-work-board/tools/queue-kanban/allocate.go`

- Change `requestReservationFileName(int)` from fixed six digits to the canonical minimum-three-digit spelling (`REQ-%03d`). Keep it the sole filename constructor used by `nextRequestNumber`.
- Keep `nextRequestNumber`'s rooted `os.Root` directory handling and `O_CREATE|O_EXCL`; canonicalizing the basename is sufficient to make a current board allocator and current capture collide at the filesystem boundary.
- Do not change the `next-req` stdout contract (`runNextRequestCommand` in `main.go` prints the decimal number). This REQ concerns marker filenames, and changing command output would be a separate caller-facing contract.

Tests: `skills/do-work-board/tools/queue-kanban/allocate_test.go`

- Add an explicit filename-format table for 1/42/482/1207; current tests derive expected paths through `requestReservationFileName`, so they cannot detect a wrong implementation of that helper.
- Rewrite `TestNextRequestNumberAdvancesPastExistingReservation` as a table (or add siblings) that plants both `REQ-482` and `REQ-000482` directly, not via the formatter, and proves both advance the max.
- Preserve `TestNextRequestNumberConcurrentProcessesReserveDistinctNumbers`, symlink/root-swap coverage, and queue-nonmutation coverage unchanged.

### Core repository allocator

File: `skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model.go`

- In `ReserveNextRequestID`, derive `markerName` from the already-canonical `formatRequestID(candidateNumber)` (`REQ-%03d`) instead of a second `REQ-%06d` formatter.
- `DiscoverRepository` is also the reader for this allocator. Give reservation entries a reservation-specific exact parser: the general `requestNumberFromText` deliberately accepts a prefix because it also parses `REQ-NNN-slug.md`; using it for marker basenames currently accepts suffix junk such as `REQ-482-copy`. Reservation parsing must require the whole basename to be `REQ-` plus digits while remaining width-agnostic.
- Keep the rooted directory identity checks, `atomicfile.CreateExclusiveAt`, and post-create current-directory check intact.

Tests: `skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model_test.go`

- Make `TestDiscoverRepositoryCoversLiveArchiveReservationAndExcludedLayouts` include one canonical and one legacy-width marker and assert both numeric identities are retained.
- Update `TestReserveNextRequestIDIsCollisionFreeAcrossEvidenceAndRaces` to assert exact returned basenames (`REQ-013`, `REQ-014`) rather than only numbers/existence.
- Add a malformed adjacent-value control (`REQ-482-copy`) that is not reservation evidence. This follows the `collision-fixture-identity` lesson: exact evidence sources must not accidentally mask one another.
- Preserve `TestReservationDirectorySymlinkIsRefused` and `TestReservationDirectorySwapCannotRedirectMarker`.

### Capture, answer override, and defer-gate publication

Files:

- `skills/do-work/tools/do-work-cli/internal/publication/capture_files.go`
- `skills/do-work/tools/do-work-cli/internal/publication/defer_gate.go`
- Prefer one small package-local canonical reservation-name/path helper shared by both files (a new `internal/publication/reservations.go` is justified if keeping the helper out of either lifecycle-specific file makes the shared rule explicit).

Required behavior:

- Parse the manifest REQ/repair ID numerically and derive the canonical minimum-three-digit marker path. Do not concatenate the raw ID and call that canonical.
- Refuse a manifest whose declared `reservation_path` is not that canonical path. Retain the existing typed codes (`CAPTURE-RESERVATION-MISMATCH`, `DEFER-GATE-RESERVATION-PATH-INVALID`) unless a more specific noncanonical-path code is deliberately introduced; the evidence should name the expected canonical path.
- Before admitting a create, enumerate/inspect the reservation directory by numeric identity and refuse when either canonical or legacy-width spelling for that number already exists. Merely checking the canonical target is insufficient for migration compatibility.
- `BuildAnswerPlan` needs no separate implementation: its `override_capture` already delegates to `BuildCapturePlan`. It does need a regression assertion through `answer_test.go` only if capture's direct test does not prove the delegated refusal is preserved with the `ANSWER-OVERRIDE-CAPTURE-` prefix.
- `BuildDeferGatePlan` fold mode consumes an existing reservation too. It must accept a matching marker under either spelling, still verify exact contents when present, and refuse ambiguity/stale bytes. Create mode must reject either spelling before creating the canonical one.
- Do not weaken `ApplyPlan`/`createRootedFile`: successful exclusive creation remains the ownership event. Once all current writers derive the same canonical basename, the existing final `atomicfile.CreateExclusiveAt` boundary supplies the live cross-writer collision.

Tests:

- `skills/do-work/tools/do-work-cli/internal/publication/capture_files_test.go`: add (1) noncanonical fixed-six manifest path refusal; (2) canonical capture refusal when a legacy six-digit alias for the same number exists; (3) canonical marker mutation path assertion.
- `skills/do-work/tools/do-work-cli/internal/publication/defer_gate_test.go`: cover create collision against a legacy alias and fold acceptance of a legacy alias with matching bytes; keep stale-byte and absent-committed-marker authority cases.
- `skills/do-work/tools/do-work-cli/internal/publication/answer_test.go`: optional narrow delegated noncanonical reservation test; no new answer-side algorithm.
- A true lock-in must encode both current flows, not two spellings produced by a test helper. Best seam: an integration test invokes/builds `queue-kanban next-req` against a temp repository, then runs `capture-files` for that returned ID with the canonical manifest path and asserts typed `CAPTURE-COLLISION`/non-success plus exactly one marker and no published REQ. If cross-module process cost is judged too high, the minimum acceptable substitute is a capture test that plants the board allocator's literal canonical output and a board test that asserts that literal; a helper-derived expectation alone can drift green.

### Documentation writers

Files:

- `skills/do-work/actions/capture.md` — change both `REQ-NNNNNN` marker examples (allocation contract and Step 5 manifest) to the canonical stored-ID spelling, explicitly saying legacy widths are read-only compatibility and manifests must declare the canonical path.
- `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md` — change the shipped `next-req` marker example from `REQ-NNNNNN` to `REQ-NNN`/stored-ID spelling.

No edit is needed in `_dev/primes/prime-kanban-board.md`: it counts the `next-req` write surface but intentionally does not duplicate its filename format. No edit is needed in `skills/do-work/scripts/cleanup-req-reservations.sh`: it is only a compatibility launcher and delegates all matching/reaping to the Go CLI.

## Reader and cleanup seams

### Allocation max scans

- Board: `allocate.go` → the reservation loop inside `nextRequestNumber`. Replace its use of the prefix-oriented `requestNumberFromText` with the exact, width-agnostic reservation parser; both `REQ-482` and `REQ-000482` contribute integer 482.
- Core: `repository_model.go` → `DiscoverRepository` projection of `.req-reservations/*` into `RepositorySnapshot.ReservationFiles`. Apply the same exact-basename/numeric rule locally; this module is separate from queue-kanban and cannot import its helper.
- Duplicate spellings for one number may remain as two evidence records, but max computation is numeric and safe. Do not silently delete or rename either during discovery.

### Cleanup

File: `skills/do-work/tools/do-work-cli/internal/corehelpers/reservations.go`

- Widen `reservationNamePattern` from exactly six digits to an exact all-digit basename and parse the integer. Both canonical and legacy-width names then share the existing committed-request and 48-hour eligibility rules.
- Keep malformed suffixes, directories, symlinks/special files, Git-authority loss, identity races, and removal failures fail-soft exactly as today.
- If both spellings coexist, evaluate and record each concrete file separately; a committed matching REQ or stale age should reap both. The numeric committed-request lookup (`requestFilePattern` + `strconv.Atoi`) already makes request-file padding irrelevant.
- Preserve the REQ-463 final-boundary rule: re-read committed Git authority, object identity, age, and eligibility immediately before each remove. Filename compatibility must not bypass that revalidation.

Tests: `skills/do-work/tools/do-work-cli/internal/corehelpers/reservations_test.go`

- Table committed-request removal and timeout removal over canonical and six-digit legacy spellings.
- Add the coexistence case: both names for one committed request are both removed and both exact paths appear in `Changes`.
- Keep a malformed adjacent filename (`REQ-482-copy` or `bad`) preserved with `RESERVATION-MALFORMED`.
- Preserve Git-unavailable, unborn-repository, dry-run, refreshed-mtime, unsafe-root, and identity-race tests. `internal/hookcommands/session_start_test.go` only needs its canonical fixture updated; it delegates cleanup and is not a second parser.

### Non-readers / no scope

- The board's card parser/model does not read reservation markers. `enumerateDoWorkTree` prunes hidden directories; `nextRequestNumber` separately opens `.req-reservations`. Therefore the schema/parser lock-step rule does **not** require `model.go`, frontend JavaScript, generated board data, or browser probes.
- `internal/corehelpers/inventory.go` preserves hidden-directory entries but does not parse reservation identity; no behavior change is needed.
- Marker file contents differ today (allocators write empty files; publication writes `REQ-ID\n`). Cleanup is name/mtime/Git based. Content canonicalization is outside REQ-485 and should not be coupled into this fix.
- No migration command, bulk rename, board UI change, request-file renumbering, or shell-parser rewrite belongs in scope.

## Requirement-mapped implementation plan

1. **RED — cross-writer collision:** add the process-level or literal-cross-contract test at candidate 482 and prove current board reservation plus capture both succeed under different basenames.
2. **Canonical writers:** change both allocator formatters to stored-ID/minimum-three-digit names; centralize publication path derivation and reject noncanonical capture/defer manifests.
3. **Legacy reads:** make the two allocation discoveries, capture/defer collision/fold checks, and cleanup parse exact numeric marker basenames independent of width.
4. **GREEN — migration:** prove canonical-vs-canonical `O_EXCL` collision, pre-existing six-digit aliases block capture/defer and advance max, and cleanup reaps both spellings without weakening authority/race guards.
5. **Guidance:** update `capture.md` and shipped board prime; leave the maintainer write-surface count and launcher architecture intact.
6. **Verification:** run `go test ./...` independently in both Go modules. No browser harness is warranted because no rendered/parser field changes. Run the repository's normal contract/release gates after integration because the shipped guidance and board prime change.

## Scope and integration risks

- **Separate modules:** queue-kanban and do-work-cli have separate `go.mod` files (Go 1.26 vs 1.25), so they cannot share an internal Go helper. Contract lockstep must be proved across the process/literal boundary.
- **Helper-derived false greens:** the existing allocator tests locate created files through the production formatter; they will pass for either spelling. Exact literal assertions and an end-to-end two-flow fixture are required.
- **Plan-time alias race:** compatibility scans can block a legacy alias that already exists, but cannot make two different filenames one atomic namespace if an old binary creates the alias after planning. The supported current writers must converge on one basename; `O_EXCL` at that canonical path is the concurrency proof. Do not claim atomicity against concurrently running pre-fix software.
- **Fold compatibility:** defer-gate treats absence of a marker as valid only for a clean committed repair. Failing to search the legacy alias would misclassify a present legacy marker as absent and bypass its byte check.
- **Broad regex hazard:** the generic REQ parsers intentionally accept filename prefixes. Reusing them unchanged for markers admits suffix junk. Add exact marker parsers rather than tightening the generic functions and accidentally breaking request-file discovery.
- **Large IDs:** `%03d` is a minimum width, not a maximum; it remains canonical above 999 without truncation. Numeric parsing should reject overflow and non-positive values instead of treating them as reservation evidence.
- **Dirty main tree:** implementation belongs in the assigned isolated builder worktree. Only this exploration artifact is written in main; lifecycle, changelog/version, merge, and release bookkeeping remain orchestrator-owned.

## Files read

- `CLAUDE.md`
- authoritative `do-work/working/REQ-485-canonicalize-req-reservation-marker-filenames.md`
- run brief `do-work/runs/work-2026-09-03-214500/REQ-485-brief.md`
- `_dev/primes/prime-kanban-board.md` and `_dev/primes/lessons-kanban-board.md`
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` and `lessons-do-work-cli.md`
- shipped board `prime-do-kanban.md` / `lessons-do-kanban.md`
- allocator, repository-model, publication, cleanup, hook, launcher, guidance, and named tests cited above

The controlling lessons are `alternate-writer-contract-drift` (sweep every creator/consumer), `collision-fixture-identity` (literal independent evidence, not helper-coupled fixtures), `final-boundary-identity` (keep exclusive-create/removal revalidation), and the cleanup migration lesson requiring parity across authority boundaries rather than a smoke-only test.
