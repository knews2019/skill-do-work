## 0.305.17 — Establish Test Efficiency Baseline with Descendant CPU and Work Counts (2026-09-06)

Introduced opt-in test efficiency instrumentation and reproducible baseline measurement tooling to record process descendant CPU times, toolchain subprocess counts, and per-file Go test duration attribution across representative verification selections.

- `_dev/tests/test-duration-log.sh` provides opt-in measurement functions capturing child CPU via POSIX getrusage and toolchain subprocess executions via PATH shims.
- `_dev/tests/test-efficiency-baseline.sh` orchestrates multi-run benchmarking across inventory, finalization, session-start, shell-audit, and heavy CLI build cases, emitting a structured evidence table.
- `_dev/tests/test-efficiency-baseline-behavior.sh` verifies opt-in isolation, exit code preservation, and Go test event stream parsing.
- Repointed archived lesson satellite links in `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`, `_dev/primes/lessons-action-files.md`, and `_dev/primes/lessons-shell-commands.md`.
