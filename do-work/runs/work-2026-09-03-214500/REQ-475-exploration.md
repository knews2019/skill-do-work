# REQ-475 exploration — configured Memory reader confinement

## Source basis

- Authoritative request: `do-work/working/REQ-475-confine-all-configured-memory-tree-readers.md`
- Run brief: `do-work/runs/work-2026-09-03-214500/REQ-475-brief.md`
- Governing guidance: `CLAUDE.md`, `skills/do-work/tools/do-work-cli/prime-do-work-cli.md`, and `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`
- Origin: `do-work/archive/REQ-417-implement-interview-memory-commands.md`, especially its remediation and residual-review summary. The archive points to `do-work/runs/work-2026-08-31-165510/REQ-417-rereview.md`, but that file is absent from the current tree and from `git rev-list --all --objects`; the durable archive summary and the remediation diff at `aa5e38bd` are therefore the recoverable re-review evidence.
- Current implementation and tests: `skills/do-work/tools/do-work-cli/internal/knowledgecommands/memory_commands.go`, `memory_commands_test.go`, and `commands_test.go`; behavior contract: `skills/do-work-knowledge/actions/memory.md` and `memory-reference.md`.

## Root cause

REQ-417 added the right confinement shape but applied it only to the mutating forget/bootstrap seams and the status ledger summary. `findMemoryMatches` opens the selected Memory directory once with `os.OpenRoot`, `readRootedMemoryDirectory` rejects a linked/non-directory `logs`, and `readOptionalRootedMemoryFile` rejects non-regular objects and validates the opened file against the pre-open object. Broad recall, status, and `auditMemoryEngine` remained older independent pathname readers (`os.ReadFile`, `os.Stat`, and `filepath.Glob`). A child link is consequently resolved by the OS before those consumers can classify it, and `filepath.Glob(<root>/logs/*.md)` traverses a linked `logs` directory.

This is an authority split, not three isolated missing `Lstat` calls. Adding preflight `Lstat` beside each `ReadFile` would retain the pathname race and the linked-directory escape. Every configured Memory-tree consumer must operate through one already-opened `os.Root`, enumerate through that handle, and open the final file through that handle.

The current shared file reader is no-follow and identity-checked, but it is not actually bounded: `readOptionalRootedMemoryFile` uses `io.ReadAll`, and `readRootedMemoryDirectory` uses `ReadDir(-1)`. The brief calls it bounded, while the checked-in implementation is not. The builder should make the limit explicit in the reusable primitive and at directory enumeration rather than copying the present unbounded call. The documented hard byte bound is 2,500 bytes for `working-memory.md`; logs and ledgers have logical bounds (three recent days / roughly 40 broad-recall lines, five status ledger events, 28-day audit buckets) but no total byte cap. A finite log/ledger/directory ceiling therefore needs a named source constant and matching action/reference wording; it must not be silently invented inside one consumer.

`resolveMemoryRoot` is a related root-boundary seam. It resolves the supplied root physically and only rejects an outside-repository root when `mutating == true`; read-only commands accept an absolute or `..` root outside the repository, despite the Memory action defining the store at `<project-root>/memory`. REQ-475 says repository-rooted. The implementation plan should make read-only configured roots use the same inside-repository check and add root-symlink / `..` / absolute-outside negatives, unless the maintainer deliberately preserves external read-only roots and records that narrower interpretation. Merely opening an already-resolved outside directory with `os.OpenRoot` is no-follow below that directory, but is not repository-rooted.

## Exact reader inventory

| Surface | Exact symbols and current reads | Escape / reporting behavior | Required seam |
|---|---|---|---|
| Broad `memory-recall` | `handleMemoryRecall` → `broadMemoryRecall`; direct `os.ReadFile(working-memory.md)`, `filepath.Glob(logs/*.md)`, direct log `os.ReadFile` | A linked working file, linked log file, or linked `logs` directory is read. Working read errors are silently treated as absence. Hit evidence contains the target line verbatim, so both renderers disclose it. | Open one Memory root in the handler/scan, use rooted directory and file reads, preserve working-first then newest-three-days ordering and capture suppression, and return a path-bearing scan error before hits are rendered or the recall ledger append occurs. |
| Lexical `memory-recall` adjacent seam | `handleMemoryRecall` → `lexicalMemoryRecall`; `Lstat` protects the working file and each glob result, but `filepath.Glob` enumerates `logs` by pathname | Individual file links are refused, but a linked `logs` directory exposes regular target files. This is the same configured-directory root cause even though the REQ instances call out broad recall. | Reuse the same rooted source inventory so empty-query and lexical recall cannot drift. Preserve scoring, eight-hit bound, attribution, and retained-script differential. |
| `memory-status` working file | `handleMemoryStatus` lines 574–580; direct `os.ReadFile` plus `os.Stat` | Follows a link twice; target frontmatter can reach `updated=...`, while target size/mtime reach typed evidence. | Rooted read must return bytes and the already-opened regular file's metadata together; do not re-stat by pathname. Missing working memory keeps `MEMORY-NOT-SET-UP`; refused/special is a distinct path-bearing failure. |
| `memory-status` logs | `handleMemoryStatus` lines 581–593; `filepath.Glob`, newest pathname, direct `os.ReadFile` | Follows a linked directory or newest linked file; a crafted target heading reaches `last_capture`. Glob/read errors are discarded. | Rooted, validated enumeration; reject any candidate `.md` that is linked/special before choosing newest. Read newest through the same root and preserve lexical filename ordering. |
| `memory-status` ledger | `memoryLedgerSummary` uses `readOptionalRootedMemoryFile` | The file itself is no-follow, but every unsafe/read/oversize error is collapsed to the ordinary string `none`, so the finding does not name the refused configured path. | Change the helper to return `(summary, error)` (or typed optional data) and make status fail with `memory/usage-ledger.jsonl` as the affected path; malformed lines remain the ordinary `malformed` summary and raw query/note bytes remain excluded. |
| `memory-audit --engine memory` working file | `handleMemoryAudit` → `auditMemoryEngine`; direct `os.ReadFile` | Follows a link; target `updated` reaches audit evidence and influences classification. All read errors are interpreted as `working_present=false`. | Pass a rooted store reader into `auditMemoryEngine`; distinguish missing from refused/error, preserve frontmatter/section counts, and return an error carrying only the configured relative path. |
| Memory audit logs | `auditMemoryEngine` lines 1132–1152; `filepath.Glob` and direct `os.ReadFile` | Follows directory/file links; target headings influence capture/note counts and filenames influence freshness. Errors are silently ignored. | Enumerate and read every accepted `.md` via the root, deterministically sort names, reject a linked/special entry or directory, and preserve the existing counts and newest-date rule. |
| Memory audit ledger | `auditMemoryEngine` → `collectLedgerAudit(path, "recall", now)` plus direct `os.Stat` | `collectLedgerAudit` plain-reads and ignores errors; linked target timestamps/events influence classification and all ledger evidence. The second pathname stat follows the same link for `ledger_mtime`. | Refactor ledger collection to consume confined bytes plus metadata (or a rooted file result); return parse stats and the same-file mtime without a second pathname lookup. Keep the BKB caller on a clearly separate path or generalize without weakening BKB behavior. |
| Existing sentinel / mutator controls | `handleMemoryBootstrap` reads `.bootstrap-imported` and target logs through `readOptionalRootedMemoryFile`; `findMemoryMatches` uses both rooted helpers | These are the reusable precedent and already have symlink canaries in `TestForgetAndBootstrapNeverFollowPrivateSymlinksOrDiscloseOutsideBytes`. They lack special-file, size-limit, directory-limit, and text-render assertions. | Preserve them while changing the helper signature; add direct primitive coverage for working/log/ledger/sentinel names so no existing caller regresses. |

The `.claude/settings.json` read inside `auditMemoryEngine` is repository configuration, not a configured Memory-tree object, and is outside this REQ's Memory-root confinement. BKB audit reads are also outside REQ-475; do not accidentally route the BKB ledger through a Memory-only limit or error policy.

## Reusable primitive and data flow

Keep the primitive in `memory_commands.go` unless the implementation demonstrates a second package consumer. A compact `memoryStoreReader` (display root plus one `*os.Root`) or equivalent function set should own:

1. root opening only after repository containment is established;
2. `readDirectory(relative, maxEntries)` with `Lstat` → rooted `Open` → `SameFile`, bounded batches, deterministic caller sorting, and a final identity check;
3. `readOptionalRegular(relative, maxBytes)` with missing distinguished from refused/error, `Lstat` regular-file refusal, rooted `Open`, pre/post identity and metadata validation, and a `maxBytes+1` read that refuses oversize without returning partial bytes;
4. safe metadata returned from the opened object, so status/audit never call `os.Stat` on the pathname afterward;
5. errors whose text contains only the configured relative path and reason, never contents or a resolved target path.

Consumers should derive typed findings from that one observation. On refusal, do not keep partial hits/evidence and do not append the recall ledger: the result should carry the command-specific code, the configured path in `affected_paths`, and generic regular-file/directory/limit evidence. Since text and JSON are both rendered from `resultmodel.CommandResult`, proving the typed result contains no canary plus rendering both formats closes the projection requirement without separate sanitizers.

Do not preserve `collectLedgerAudit(path, ...)` as an alternate Memory reader. Split parsing from acquisition, e.g. `collectLedgerAudit(data, ...)` or `parseLedgerAudit(io.Reader, ...)`; Memory supplies confined bytes, while BKB retains its own acquisition until a BKB-specific confinement REQ governs it.

## RED/GREEN test seams

Add the RED matrix in `memory_commands_test.go` before production changes. Use an outside regular fixture whose bytes include a unique canary in every field that a successful old read projects (working `updated`, log heading/body, ledger timestamp/event), then exercise real handlers and `resultmodel.RenderResult` in text and JSON.

Required matrix:

- broad recall: linked `working-memory.md`, linked `logs` directory, and linked newest log;
- lexical recall: linked `logs` directory (working/link-entry refusal already exists, but keep one parity assertion if source enumeration is unified);
- status: linked working file, log directory, newest log, and ledger;
- memory audit: linked working file, log directory, any `.md` log, and ledger;
- shared primitive / bootstrap regression: linked and special `.bootstrap-imported`, working, log, and ledger files; linked/special `logs` directory;
- configured-root boundary: symlinked root and explicit `..` / absolute outside root if repository-rooted is enforced;
- finite limits: exact-limit accepted and limit+1 refused for each configured file class, plus exact/max+1 directory-entry cases. A FIFO test must run under a deadline or use a nonblocking special object so RED cannot hang.

Each route assertion should require: non-success/refused outcome appropriate to the command; the configured path (not resolved target) in `AffectedPaths`; canary absent from the typed result, rendered text, and rendered JSON; no partial ordinary hit/audit/status finding; no changed Memory bytes. Status ledger refusal must not degrade to `ledger_tail=none`. `--engine both` should still keep deterministic finding order while making the Memory-engine refusal visible; whether BKB evidence is retained alongside an overall failure should be pinned explicitly.

GREEN parity should retain and strengthen the existing tests:

- `TestMemoryRecallScoringStatusAndBroadRecall`: same broad order, exclusion of session-capture bodies, scoring and source attribution, `cap=2500`, and caller spelling `--memory-root ./memory` in recovery argv;
- `TestLexicalRecallMatchesRetainedScriptAtRecencyBoundaries`: no scoring/order drift after source inventory moves behind the rooted reader;
- `TestRememberForgetAndStatusNeverPersistOrRediscloseProtectedTextInLedger`: status still exposes only event/timestamp/hit summary fields;
- `TestMemoryAuditReportsRequiredProbesAndOldCitationDoesNotMakeActive` and `TestMemoryAuditUsesExactFourteenDayBoundary`: identical classifications, counts, weeks, and boundary semantics;
- `TestRealRuntimeTextAndJSONProjectSameInterviewAndMemoryFindings`: extend to assert the refusal path/code and canary absence, not merely that both renderers mention existing codes;
- byte preservation: digest the store before and after status/audit; for recall account only for the documented best-effort ledger append, or link/refusal cases assert the append did not occur.

Focused validation: `go test -count=1 ./internal/knowledgecommands`, `go test -race ./internal/knowledgecommands`, `go vet ./...`, then `go test -count=1 ./...`. Run `_dev/tests/contract-regressions.sh` if the Memory action/reference contract changes. No browser work is relevant.

## Risks and decisions to pin

- **Limits are underspecified for logs/ledger/directories.** The builder must name and document finite ceilings before tests freeze them. Do not mistake output count (three days, five ledger rows, eight lexical hits) for a byte-read bound.
- **Root containment can be compatibility-affecting.** Current read-only `--memory-root` accepts outside roots. The shipped action says project-root Memory and REQ-475 says repository-rooted, so refusal is the safer reading; preserve the user's original spelling in affected paths/recovery argv.
- **Missing is not unsafe.** Missing `logs` and missing ledger are ordinary (`0` / `none`); a present link, special object, oversize object, or changed identity is a refusal. Do not collapse these states.
- **Best-effort ledger writes do not license best-effort ledger reads.** `appendMemoryLedger` may silently skip an unsafe ledger, but status/audit must report a configured object they were asked to inspect.
- **Do not leak through errors.** Never include read buffers, parsed outside values, resolved link targets, or partial results in error evidence. OS errors should be wrapped against the configured relative path only.
- **Preserve one-root observation.** Reopening the Memory root independently for working/log/ledger reads reintroduces a swap seam and inconsistent snapshots. One handler-scoped root is the clean boundary.
- **Audit API must carry failure.** `auditMemoryEngine` currently cannot report read failure. Change its signature rather than encoding an error into the classification/evidence string.

## Compact requirement-mapped plan and scope

1. **RED — all readers / all link positions / no disclosure.** In `memory_commands_test.go`, add the handler-level symlink matrix, special-object and limit tests, store-byte checks, exact affected paths, and typed/text/JSON canary absence. Extend the runtime projection assertion in `commands_test.go` only if needed for real CLI status parity.
2. **One repository-rooted reader.** In `memory_commands.go`, make `resolveMemoryRoot` enforce the chosen repository boundary, evolve `readOptionalRootedMemoryFile` and `readRootedMemoryDirectory` into bounded identity-checked operations, and return opened-file metadata. Keep missing distinct.
3. **Migrate recall.** Give broad and lexical recall one rooted source inventory; eliminate `filepath.Glob` and direct configured-file reads while preserving broad newest-three ordering, capture filtering, lexical scores/order/eight-hit limit, attribution, and failure-before-ledger behavior.
4. **Migrate status.** Read working bytes/metadata, log inventory/newest log, and ledger through the same reader. Propagate unsafe ledger and directory/file errors as path-bearing typed failures; retain ordinary `none`, cap findings, and recovery argv spelling.
5. **Migrate Memory audit.** Make `auditMemoryEngine` return an error, acquire working/log/ledger evidence through the same reader, parse ledger bytes separately, and compute mtime from confined metadata. Preserve BKB acquisition and finding order.
6. **Protect existing callers.** Update bootstrap/find/ledger-summary callers for the bounded helper and add sentinel/mutator regression rows; inventory `memory_commands.go` afterward to ensure no configured Memory child is still read by `os.ReadFile`, `os.Stat`, `filepath.Glob`, or `os.ReadDir`.
7. **Document and qualify.** If new finite limits or outside-root refusal are introduced, update `skills/do-work-knowledge/actions/memory.md` and `memory-reference.md`; update the CLI prime/lesson only if the confinement family is promoted. Run focused/race/full/vet and contract checks, review every changed file, and leave version/changelog, queue lifecycle, integration, and release stamping to the orchestrator.

Expected minimum builder source scope is `memory_commands.go` plus `memory_commands_test.go`; `commands_test.go` is justified only for real-runtime refusal projection, and Memory action/reference files are required only to close the currently undocumented finite bounds/root policy. No board, hook, git-transaction, BKB implementation, generated artifact, or screenshot change belongs to REQ-475.
