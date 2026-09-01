# REQ-418 Plan — Toolbox Commands and Absorbed Audit Metrics

## Route and estimate

- **Route C:** new shared command family, typed metrics model, transactional publication, asynchronous process ownership, and differential migration proof.
- **P50:** 170 active minutes.
- **Confidence:** low.
- **Estimator basis:** Route C; 30-file write set; 24 new files; 7 subsystems; 10 acceptance groups; dependency depth 13; persistence; async lifecycle; cross-route regression gates; full-suite verification.

Computed with:

```text
tools/estimate-p50.sh --route C --write-set 30 --new-files 24 --subsystems 7 --acceptance 10 --deps-depth 13 --persistence --async-behavior --regression --full-suite
```

The estimate is informational and frozen before implementation. The dominant uncertainty is adversarial media cancellation and multi-path publication, not the audit algorithms themselves.

## Outcome

Deliver seven canonical CLI entries with retained argument/status/effect behavior and actionable JSON:

```text
do-work-note
architecture-report-preflight
generate-report-image
generate-report-image-batch
publish-portfolio-summary
install-last30days
audit-metrics inventory|folders|churn|hotspots
```

Keep every current toolbox script and the standalone audit-metrics module unchanged as migration oracles. REQ-419 later exposes recipes and delegates actions; REQ-420 later turns retained executables into thin shims and removes the standalone module.

## Frozen implementation allowlist

Only these 30 tracked product paths may change. Paths marked `(new)` do not exist at the baseline.

1. `skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go`
2. `skills/do-work/tools/do-work-cli/internal/commandruntime/command_runtime.go`
3. `skills/do-work/tools/do-work-cli/internal/commandruntime/command_runtime_test.go`
4. `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go`
5. `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go`
6. `skills/do-work/tools/do-work-cli/prime-do-work-cli.md`
7. `skills/do-work/tools/do-work-cli/internal/toolboxcommands/commands.go` (new)
8. `skills/do-work/tools/do-work-cli/internal/toolboxcommands/commands_test.go` (new)
9. `skills/do-work/tools/do-work-cli/internal/toolboxcommands/mutation.go` (new)
10. `skills/do-work/tools/do-work-cli/internal/toolboxcommands/mutation_test.go` (new)
11. `skills/do-work/tools/do-work-cli/internal/toolboxcommands/note.go` (new)
12. `skills/do-work/tools/do-work-cli/internal/toolboxcommands/note_test.go` (new)
13. `skills/do-work/tools/do-work-cli/internal/toolboxcommands/architecture.go` (new)
14. `skills/do-work/tools/do-work-cli/internal/toolboxcommands/architecture_test.go` (new)
15. `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_image.go` (new)
16. `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_image_process.go` (new)
17. `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_image_process_unix.go` (new)
18. `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_image_process_windows.go` (new)
19. `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_image_test.go` (new)
20. `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_image_process_test.go` (new)
21. `skills/do-work/tools/do-work-cli/internal/toolboxcommands/portfolio.go` (new)
22. `skills/do-work/tools/do-work-cli/internal/toolboxcommands/portfolio_test.go` (new)
23. `skills/do-work/tools/do-work-cli/internal/toolboxcommands/last30days.go` (new)
24. `skills/do-work/tools/do-work-cli/internal/toolboxcommands/last30days_test.go` (new)
25. `skills/do-work/tools/do-work-cli/internal/toolboxcommands/audit_metrics.go` (new)
26. `skills/do-work/tools/do-work-cli/internal/toolboxcommands/audit_metrics_test.go` (new)
27. `skills/do-work/tools/do-work-cli/internal/toolboxcommands/audit_inventory.go` (new)
28. `skills/do-work/tools/do-work-cli/internal/toolboxcommands/audit_inventory_test.go` (new)
29. `skills/do-work/tools/do-work-cli/internal/toolboxcommands/audit_churn.go` (new)
30. `skills/do-work/tools/do-work-cli/internal/toolboxcommands/audit_churn_test.go` (new)

Explicitly forbidden without stopping and re-planning:

- `skills/do-work/tools/do-work-cli/internal/corehelpers/inventory.go` and `inventory_test.go` (REQ-462 ownership);
- every other `internal/corehelpers` path;
- all retained `skills/do-work-toolbox/scripts/*.sh` files;
- all files under `skills/do-work-toolbox/tools/audit-metrics/`;
- toolbox actions, docs, `SKILL.md`, root/core help, `justfile`, and `skills/do-work-board/justfile.template` (REQ-419 ownership);
- installer/update/managed-section recipe collision code (REQ-419 ownership);
- maintainer gate removal of the separate audit lane and shell shim conversion (REQ-420 ownership);
- lifecycle, release, changelog, version, lesson, and queue files (orchestrator-only).

No broad directory wildcard is authorized. If the builder needs another tracked path, it must stop before editing and hand back the exact requirement that cannot be met inside this set.

## Task 1 — RED first: register the family and establish typed compatibility projections

Before production behavior, write command-level tests proving all seven names currently return `UNKNOWN-COMMAND`. Add the `toolboxcommands.Handlers()` registration loop at the sole command boundary only after RED is captured.

Extend the result model narrowly:

- an optional typed audit payload for report kind, aggregates, distributions, band rows, folder rows, churn rows, hotspot rows, unavailable paths, and history metadata;
- an exact compatibility-text projection rendered from the same typed command observation; and
- an internal interruption-status override accepted only for cleaned-up 129/130/143 media interruption.

Pin JSON normalization (no null collections), stable field names/order, ordinary 0–4 outcomes, invalid override refusal, exact text compatibility, and command registration through real CLI invocations. `commands.go` owns names, option dispatch, stable finding constructors, and shared quoting only; domain algorithms stay in their files.

GREEN for this task: every public name resolves, invalid usage returns 2 with exact next/verification argv, text and JSON originate from one result, and no toolbox command imports or executes a retained implementation.

## Task 2 — Implement transactional notes and publication utilities

Use `gittransaction` and `atomicfile` as dependencies; do not modify them.

- **Note:** exact normalization, injectable date, append-only byte preservation, empty no-op, dry-run, dirty-target/empty-index guards, exact commit and rollback.
- **Architecture scan/publish:** numeric prior selection, HTML/watermark rules, UTC slug, missing/unresolvable distinctions, exclusive suffix allocation, hidden verified copy, occupied failed bundle, dry-run/commit, and final-path reporting.
- **Portfolio:** separate verified source copies, canonical-only replacement, snapshot-first exclusive suffixing, distinct inodes, retained snapshot on later canonical failure, directory/final-boundary collision defense, exact commit only on complete requested success.
- **last30days:** source/current-tree union planning, complete payload validation, clone-to-private-source, transactional exact replacement, restoration on pre-publication failure/interruption, target-reappearance defense, local exclude verification through the existing core helper result, Python 3.12+ as a target check only, dry-run, and actionable `--commit` refusal for the private ignored tree.

RED/GREEN tests must transcribe the retained prescribed-shell fixtures rather than merely test happy paths: missing/partial sources, occupied directories, suffix `-10` ordering, unreadable watermark, failed copies/moves/links, final-boundary destination appearance, existing-tree restoration, interruption, local-exclude states, no Git repository, missing/old Python, and exact changed/committed paths. Differential cases compare normalized status, ordered stdout/stderr facts, final bytes, modes, paths, retained outputs, and cleanup.

The snapshot-first recovery state is an explicit exception to all-or-nothing multi-target rollback because preserving it is a named existing user-visible contract. Tests must prove it rather than accidentally normalize it away.

## Task 3 — Implement and adversarially prove the report-image lifecycle

Build one in-process image generation engine used by both commands. Do not have the batch spawn the CLI recursively.

1. Validate every argument and final target before creating scratch.
2. Allocate adjacent invocation-private files/directories with private modes.
3. Start direct `imagegen` with prompt bytes passed only as argv.
4. When and only when explicitly authorized, fall back to `codex exec` in an exact private working directory and copy only verified `generated.png`.
5. Bind every child to an owned process boundary before backend bytes run.
6. On cancellation, signal/terminate/escalate/reap all owned groups before removing scratch.
7. Publish only nonempty status-backed outputs and verify exact final placement/identity.
8. Batch concurrently, retain each result, wait all workers, remove failed members, publish one complete `generated/`, and return success with no output when all fail.

RED/GREEN matrix:

- direct success and failed partial output with a stale final target;
- prompt metacharacters remain inert;
- direct-to-agentic fallback and opt-in absence;
- interruption before/after launch registration, child and grandchild ownership, TERM-obedient zombie-only group, TERM-deaf escalation, unrelated process survival, and conventional status;
- mixed batch status, wait-all, all-failed fallback, usage-before-allocation, existing/final-boundary `generated/`, early interruption, no stale success, and exact scratch cleanup;
- Unix race runs plus Windows build/fail-closed ownership proof.

Use deterministic test seams at identity and launch boundaries. A test that only checks the final file is insufficient: every case asserts surviving PIDs/groups and private-path inventory.

## Task 4 — Port audit-metrics and prove full differential parity

Port the standalone implementation by behavior and tests into `audit_metrics.go`, `audit_inventory.go`, and `audit_churn.go`. Keep computation separate from rendering. The exact Markdown renderer and typed JSON projection consume the same report structs.

Mandatory characterization rows:

- option parsing, leftover tokens, repeatable empty-default excludes, inside-repo `--repo-root`, non-Git/Git-declined errors, top-count/window defaults;
- tracked/unreadable/binary/unterminated files, extension and direct-folder grouping, empty and tie distributions, nearest-rank values;
- absent/one-sided/equal/overlapping WATCH/FLAG thresholds and strictly-greater precedence;
- ordinary churn, rename chains, copy-first/delete-later migration, deletion as non-touch, outright-deleted drop, nearest surviving copy, excluded-but-live source, exclude-after-normalization, SHA-1/SHA-256 hash lines, deterministic ties;
- shallow clones, current content unavailable, binary hotspot score, and NOT-MEASURED ordering.

Port the mature standalone unit tests first. Then run old standalone binary and new real CLI against identical real-Git fixtures, comparing exact Markdown and exit status for all four subcommands. For JSON, assert every Markdown row has one corresponding typed record with identical numeric/path values and stable order. Mutation-test the differential harness by changing a status, heading, numeric value, path, and row order and proving each mismatch is detected.

The standalone module stays intact and its own tests continue to pass. No audit command imports `corehelpers`, and no symbol or file from REQ-462's uncommitted inventory is touched.

## Requirement coverage

- `do-work-note`: Tasks 1–2.
- architecture preflight/publication: Tasks 1–2.
- report image generation/batch, media cancellation, failure cleanup: Tasks 1 and 3.
- portfolio publication and snapshot recovery: Tasks 1–2.
- last30days check/install with target-only Python: Tasks 1–2.
- audit inventory/folders/churn/hotspots, compatible flags/Markdown, typed JSON: Tasks 1 and 4.
- dirty-target, dry-run, optional meaningful commit, rollback/risk, exact paths: Tasks 1–3.
- retained source and characterization parity before retirement: Tasks 2–4.
- queue-kanban remains separate; no board command/package is imported or modified.

## RED/GREEN and verification gate matrix

| Gate | RED / purpose | GREEN command or evidence |
|---|---|---|
| Registration | all seven names return `UNKNOWN-COMMAND` | focused real-CLI registration tests |
| Note | no canonical append command | exact normalization/date/append, dirty/dry-run/commit/rollback tests |
| Architecture | Go command absent; partial/collision fixtures have no target | old/new scan and publish differential fixtures |
| Portfolio | Go command absent; snapshot/canonical failure states unrepresented | exact bytes, inode separation, suffix, retained-snapshot, commit tests |
| last30days | Go command absent; target Python absence has no typed result | complete-tree, rollback, ignore, Python-present/absent/old, dry-run tests |
| Media | Go commands absent; cancellation cannot be proven | focused process-tree and publication matrix, then race run |
| Audit unit | algorithms exist only in separate module | ported standalone unit corpus in `toolboxcommands` |
| Audit differential | new command is unknown | exact old/new status + Markdown for all four subcommands; typed JSON row parity |
| Result/runtime | no typed metrics or interruption seam | `go test -count=1 ./internal/resultmodel ./internal/commandruntime` |
| Toolbox focused | family absent | `go test -count=1 ./internal/toolboxcommands` |
| Toolbox race | async cancellation/batch races | `go test -race ./internal/toolboxcommands` |
| Shared module | new package could regress prior commands | `go test -count=1 ./...` and `go vet ./...` in do-work-cli module |
| Go floor | platform code may use newer APIs accidentally | `bash _dev/tests/do-work-cli-go125-compatibility.sh` |
| Windows | process owner must compile and fail closed | `GOOS=windows GOARCH=amd64 go test -c ./internal/toolboxcommands -o <temp>` |
| Standalone oracle | migration must not weaken old source/tests | `go test -count=1 ./...` and `go vet ./...` in standalone audit-metrics module |
| Retained scripts | sources remain behavioral oracles | each five focused prescribed-shell case file, then `bash _dev/tests/prescribed-shell-scripts-behavior.sh` |
| Cross-contract | no action/recipe/shim ownership leakage | `bash _dev/tests/contract-regressions.sh` and `bash _dev/tests/staged-skills-contract.sh` |
| Scope | parallel writers stay disjoint | exact changed-path set is a subset of the 30-file allowlist; neither REQ-462 file appears |
| Hygiene | no debug/generated artifacts | `gofmt`, `git diff --check`, full untracked inventory review; no built binary staged |
| Canonical | repository-wide acceptance | unpiped `bash _dev/tests/maintainer-verify.sh` exits 0 |

All checks run unpiped. Focused success does not waive a later red gate. The builder reports the exact RED test name/failure before production edits and the corresponding GREEN result after.

## Decisions and defaults

- **D-01 — legacy executable names are the CLI names.** This makes REQ-420 shims mechanical and preserves argument interfaces.
- **D-02 — actions/recipes wait for REQ-419.** REQ-418 proves direct commands; it does not partially migrate callers.
- **D-03 — retained sources are read-only oracles.** No shell or standalone audit source edit is needed to build the replacement.
- **D-04 — typed audit rows, exact Markdown.** JSON is structured; compatibility text remains pasteable and byte-stable.
- **D-05 — OS interruption is a narrow exit exception.** Normal command outcomes remain 0–4; cleaned media cancellation preserves conventional signal status.
- **D-06 — snapshot-first partial success remains visible.** A retained immutable snapshot is recovery evidence, not rollback debris.
- **D-07 — last30days Python is target evidence.** The Go command probes it but never depends on it to implement publication.
- **D-08 — no core inventory sharing.** Tracked audit metrics and uncommitted porcelain inventory are separate contracts and separate packages.
- **D-09 — no generic process framework.** Media owns its process lifecycle locally; queue selection remains untouched.
- **D-10 — 30 paths are frozen.** Any expansion requires a pre-edit re-plan, not an after-the-fact scope note.

## Plan validation

- Every REQ requirement maps to one implementation task and a named RED/GREEN gate.
- Every planned edit supports command registration, typed output, one requested behavior family, or its focused tests/prime map; there are no orphan action, recipe, guide, or shim edits.
- The work is grouped into four reviewable implementation tasks despite the broad file count: result boundary, transactional utilities, media lifecycle, and audit port.
- The only parallel-wave semantic overlap with REQ-462 is the word `inventory`; the exact file sets and command names are disjoint.
- This plan deliberately leaves separate-module retirement and canonical gate-lane replacement to REQ-420, after parity evidence exists.

*Prepared from main `b713fa8b`; scratch only, not staged or committed.*

## Sole remediation scope note — 2026-09-01

The independent review reproduced concurrent tracked-replacement loss during a failed commit rollback. The frozen toolbox package cannot close that boundary through the current `gittransaction` API: tracked targets are restored from HEAD without exposing an invocation-publication identity/content comparison seam to callers. The sole remediation therefore activates the plan's conditional expansion to exactly:

- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go`
- `skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction_test.go`

The expanded transaction change is limited to identity/content-bound tracked rollback and cancellation propagation required by REQ-418. RED/GREEN tests must protect another writer's replacement and replay all affected transaction callers through the full do-work-cli suite. The original 30 paths remain the primary scope; no other expansion, lifecycle edit, retained oracle edit, REQ-419/420 surface, or REQ-462 inventory path is authorized.
