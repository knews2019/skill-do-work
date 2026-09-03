# Lessons Index

Plain routing index for lesson satellites. Each satellite has exactly one table row. `tokens` is `ceil(bytes / 4)`: run `wc -c` on the satellite, add three, then integer-divide by four. `families` is the exact sorted set of `[family: <slug>]` markers present in lesson bullets, or `none`; `slugged` is `full` only when every lesson bullet has at least one marker.

| Satellite | When it applies | Families | Tokens | Coverage |
| --- | --- | --- | ---: | --- |
| `_dev/primes/lessons-action-files.md` | Changing action routing, pipeline fields, status contracts, downstream readers, alternate artifact writers, or budgeted context routing | `alternate-writer-contract-drift, budgeted-context-routing, cross-action-exception-closure` | 3636 | `slugged: partial` |
| `_dev/primes/lessons-kanban-board.md` | Changing queue-kanban parsing, views, static output, timeline behavior, or board publication | `none` | 4707 | `slugged: partial` |
| `_dev/primes/lessons-shell-commands.md` | Changing shipped shell, argv/quoting, prescribed command blocks, publication scripts, or migration parity fixtures | `legacy-fixture-implementation-shape` | 3385 | `slugged: partial` |
| `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` | Changing queue-kanban model, parser, UI, timeline, testing, or browser behavior | `unknown-reads-as-clean` | 5562 | `slugged: partial` |
| `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` | Changing alternate stored-format writers, semantic recovery completeness, collision identity fixtures, frozen target projection and fan-out bounds, rollback/deletion/final-boundary identity, interruptible blocking input, publication topology classification, structured evidence projection, cross-action exception closure, or migration parity in do-work-cli internals | `alternate-writer-contract-drift, collision-fixture-identity, cross-action-exception-closure, final-boundary-identity, interruptible-blocking-io, observed-subset-is-not-semantic-completeness, opaque-evidence-projection, projection-before-bounding, publication-target-topology-classification, smoke-vs-characterization` | 2786 | `slugged: partial` |
| `skills/do-work/tools/lessons-do-work-update.md` | Changing updater managed paths, preservation rules, receipts, or source/install parity | `none` | 628 | `slugged: partial` |
