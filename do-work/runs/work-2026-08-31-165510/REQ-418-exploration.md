# REQ-418 Exploration — Toolbox Commands and Absorbed Audit Metrics

## Decision

Use Route C and add a dedicated `internal/toolboxcommands` family to the shared CLI. Keep the seven public command spellings aligned with the retained utility paths so REQ-420 can later replace each script with a build-and-exec shim without translating arguments:

```text
do-work-note
architecture-report-preflight
generate-report-image
generate-report-image-batch
publish-portfolio-summary
install-last30days
audit-metrics
```

`architecture-report-preflight` retains `--scan` / `--publish`; `install-last30days` retains positional `check|install`; `audit-metrics` retains positional `inventory|folders|churn|hotspots` and all current flags. The other script-backed commands retain their current argument order. This gives the compatibility work one command per retained executable while preserving audit-metrics' established subcommand interface.

REQ-418 supplies the canonical implementations, exact characterization, stable text/JSON, mutation guards, and process ownership. It does **not** edit natural-language actions, Just recipes, help/guides, retained toolbox scripts, or the standalone audit-metrics tree. REQ-419 owns delegation/recipes/documentation; REQ-420 owns shim conversion, separate-module removal, and replacement of the separate audit lane in the canonical gate.

## Baseline and executable RED

- Repository baseline: `b713fa8b`.
- `do-work-cli --format json do-work-note test` exits 2 with `UNKNOWN-COMMAND`.
- `do-work-cli --format json audit-metrics inventory` exits 2 with `UNKNOWN-COMMAND`.
- The other five public command names are likewise absent from `cmd/do-work-cli/main.go`.
- Deterministic toolbox mechanics still live in five shell utilities plus the prose-only note action; audit measurement still lives in a separate Go module. There is no shared typed JSON projection for any of them.

That is the runnable RED required by the REQ. The behavioral RED expands it by running each retained fixture against a missing Go command, including report-image interruption/process-tree cases and last30days with no qualifying target Python.

## Sources inspected

Read in full or at the relevant operative sections:

- repository `AGENTS.md`, `CLAUDE.md`, REQ-418, and UR-081;
- `_dev/primes/prime-shell-commands.md` and `skills/do-work-toolbox/tools/audit-metrics/prime-audit-metrics.md`;
- work planning, exploration, scope, TDD, and estimation instructions;
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` and its lesson satellite;
- shared runtime/result/transaction/atomic-publication command surfaces;
- toolbox note, architecture-report, ai-report image, present-work portfolio, install/last30days, and maintainability-audit action contracts;
- all five retained toolbox scripts and their prescribed-shell fixture files;
- every standalone audit-metrics source and test file;
- REQ-414 and REQ-417 planning/exploration precedent; and
- REQ-462 plus its current exploration/plan.

No tracked product file was modified, staged, or committed by this exploration.

## Existing ownership and reusable seams

- `cmd/do-work-cli/main.go` is the sole registration boundary. It should register `toolboxcommands.Handlers()` once; toolbox packages must not self-register.
- `commandruntime` owns global `--repo-root`, `--format`, real-command dispatch, rendering, and the normal 0–4 outcome mapping.
- `resultmodel.CommandResult` is the single text/JSON source. Its existing findings, changes, skipped work, and rollback fields cover command state, but audit-metrics needs an explicit typed metrics projection; an opaque Markdown string is not useful JSON.
- `gittransaction` already provides exact target dirty checks, empty-index `--commit`, rollback, exact staging/commit verification, created-file and created-directory recording, private-untracked target handling, and committed-risk reporting. Use it instead of building toolbox-specific Git guards.
- `atomicfile` supplies existing-file replacement and exclusive file publication. Directory-bundle publication still needs command-owned exclusive directory allocation and exact identity-bound cleanup.
- Media process handling has no reusable public process owner. `nextselection`'s blocked-probe runner is private and timeout-shaped; report image generation needs backend fallback, multiple concurrent children, signal forwarding, retained per-child statuses, and all-failed success. Keep this implementation inside `toolboxcommands` rather than coupling media to queue selection.
- `corehelpers` already exposes the local-exclude command through its handler map. last30days may consume that canonical result rather than copy the exclude matcher, but tree publication and verification remain toolbox-owned.

## Command ownership and compatibility

### `do-work-note`

Normalize exactly one leading `add `, one matching outer quote pair, and surrounding whitespace; reject an empty result without writing. Append one `- [YYYY-MM-DD] <text>` line without sorting, deduplicating, or rewriting existing bytes. The deterministic clock is injectable in tests. The target is `do-work/notes.md`; `--dry-run` and optional `--commit` use the shared Git transaction. A dirty target is an actionable refusal, as required by UR-081, even though the prose action previously appended without a Git guard.

### `architecture-report-preflight`

`--scan <reports-directory>` remains read-only and emits the same ordered `key=value` text: resolvable short HEAD, UTC minute slug, unsuffixed candidate, newest numeric-suffix HTML baseline, watermark, and resolution state. Missing/unreadable watermarks stay distinct from no prior report.

`--publish <draft> <candidate>` reserves the first free bundle directory exclusively, publishes a verified hidden temporary copy to `index.html`, never exposes partial HTML, never reuses a failed occupied bundle, and supports dry-run/commit around the exact final path. The result records both the occupied candidates and the actual published path.

### Report image lifecycle

`generate-report-image` retains absolute-output validation, direct `imagegen` first, explicit `DO_WORK_AI_REPORT_ALLOW_AGENTIC_BACKEND=1` fallback to `codex exec`, inert prompt argv, adjacent private output, nonempty verification, atomic replacement, and cleanup of its exact scratch.

`generate-report-image-batch` validates all `<bare-name>:<prompt>` pairs before allocation, runs every image concurrently through the same in-process generator, retains every status, waits every child, removes only failed outputs, treats all-failed as a successful empty fallback, and publishes `generated/` only when at least one status-backed nonempty image exists. Existing `generated/` or a final-boundary collision is a refusal, never a merge or clobber.

On Unix, each external backend starts in its own verified process group. HUP/INT/TERM cancels the whole invocation, sends the original/TERM signal, waits the bounded grace period, escalates to KILL, reaps, removes private paths, and returns the conventional interruption status. On platforms where the standard library cannot prove descendant ownership, generation fails closed before launch; a Windows compile test pins that boundary. The runtime needs one narrow internal interruption-status override because ordinary command outcomes remain 0–4, while a cleaned-up OS interruption must preserve 129/130/143 compatibility.

### Portfolio publication

Retain `--canonical-only` and `--with-snapshot`. Both outputs come from separately verified copies so a snapshot never shares an inode with canonical. Snapshot publication is exclusive and advances numeric suffixes. The preservation branch remains deliberately snapshot-first: if canonical refresh later fails, the successfully published immutable snapshot is retained and reported, matching the current user-visible recovery contract. `--commit` is allowed only after both requested outputs succeed and exact committed paths can be verified; it does not erase the snapshot-first recovery exception.

### last30days

Retain `check|install <project-root> [source-repository]`, default upstream, complete-payload predicate (`SKILL.md` plus `scripts/last30days.py`), full-subtree installation, exact target replacement/rollback, local Git exclude verification, and target Python 3.12+ discovery. Python is executed only as the installed target's prerequisite probe; no CLI implementation branch imports or embeds Python or jq.

The installed tree is intentionally machine-local and ignored, so `install --commit` is refused as meaningless rather than staging private content. `--dry-run` validates the source and complete planned target set without publication. Check remains usable without mutation and reports non-Git ignore state as `n/a` exactly as today.

### absorbed `audit-metrics`

Port, do not re-derive, the mature standalone module:

- tracked NUL-safe inventory, binary sniff, line/word counts, extension totals, direct-child folder counts;
- nearest-rank median/p90/p95/max;
- optional flag-only WATCH/FLAG bands with strictly-greater edges;
- empty default excludes and raw repeatable repo-relative prefix semantics;
- newest-to-oldest rename normalization, staged-copy history reassignment, deletion non-touch semantics, live-set vs report-set separation, shallow-clone reporting, and unavailable hotspot rows;
- exact stable sort orders, Markdown headings/tables, default top count 10, and default history window `12 months`.

The canonical text output stays byte-compatible Markdown. JSON carries typed inventory aggregates, distributions, bands, folder rows, churn rows, hotspot rows, shallow/commit/window metadata, and unavailable paths—not only a Markdown blob. Both projections are rendered from the same computed report.

The standalone module and tests remain untouched as the characterization oracle. Port its unit cases into the new package and add differential fixtures that run old and new commands on the same repositories, comparing exit status and exact Markdown. Retirement is forbidden until this parity matrix passes; REQ-420 performs the deletion later.

## Result-model seam

The narrow model extension is:

- an optional internal exact-text projection for compatibility-shaped commands, generated from the same typed observation as JSON rather than independently authored output; and
- an optional typed audit-metrics payload containing the report kind and its relevant rows/metadata.

Normal commands continue using findings/changes/skipped/rollback. Every refusal or incomplete operation has a stable code, exact paths, observed evidence, stop reason, next argv, and verification argv. Clean metric rows are data, not warning findings, so audit success remains exit 0. The optional interruption exit override is internal-only and accepted only for conventional signal statuses after cleanup.

## REQ-462 overlap analysis

The overlap is naming only:

- REQ-462 owns `internal/corehelpers/inventory.go` and `inventory_test.go`: uncommitted Git porcelain `M/A/D/X/XD`, secret provenance, quarantine, and REQ association.
- REQ-418 owns new `internal/toolboxcommands/audit_*` files: tracked-file sizes, folders, history churn, and hotspots.

REQ-418 must not modify either corehelper inventory file, `internal/corehelpers/commands.go`, the retained uncommitted-inventory shell helper, or protected-inventory. Its public audit entry point is `audit-metrics`, so there is no command-name collision with `uncommitted-inventory`. The two builders can run concurrently with disjoint write sets.

## Routes considered

### Route A — wrap or invoke the current scripts/module

Rejected. That would leave shell and the separate Go module as authorities, provide opaque JSON, and fail the no-Python/jq/shell-domain-logic end state.

### Route B — direct source port with ad hoc output strings

Rejected. It can reproduce happy-path text but not the shared transaction contract, structured metrics JSON, or adversarial process/publication ownership.

### Route C — one toolbox family with typed reports and shared transactions

Selected. It makes deterministic behavior canonical in Go while keeping compatibility sources alive as differential oracles until REQ-420.

## Risks and controls

- **Media descendants survive cancellation:** process-group creation precedes user bytes; tests cover children/grandchildren, TERM-obedient and TERM-deaf trees, launch-window interruption, wait-all, and unrelated-process survival.
- **Pathname cleanup removes another writer's object:** bind cleanup to rooted object identity/content and inject final-boundary parent/destination swaps.
- **Snapshot semantics get normalized into rollback:** keep the documented snapshot-first partial-success state explicit and test retained snapshot plus unchanged canonical.
- **last30days partially replaces a tree:** enumerate the source/current union before mutation, publish from a complete staging tree, and inject copy, backup, rename, reappearing-target, interruption, ignore, and Python failures.
- **Audit history looks plausible while wrong:** compare exact old/new output on rename, copy-first/delete-later, excluded-live-source, shallow, outright deletion, and unavailable-current-content fixtures; mutation-test the comparator.
- **JSON and Markdown drift:** generate both from one typed report and table-test every row/order; never parse the Markdown back into JSON.
- **REQ-419 scope leaks backward:** no actions, recipes, help, guides, managed templates, or install/update recipe collision rules in this REQ.
- **REQ-462 collision:** exact prohibition on both `corehelpers/inventory*` files; any need for them stops implementation and forces serial re-planning.

## Exploration confidence

Confidence is **low**. The public ownership and compatibility boundaries are firm, and audit-metrics is mature code with strong tests. Error bars are wide because the media commands own asynchronous process trees, five mutation families need Git/rollback evidence, and the 30-file frozen surface spans seven runtime subsystems.

