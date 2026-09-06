## 0.305.21 — Use Fast POSIX Checksums for Fixture Integrity Verification (2026-09-06)

Optimized shared fixture integrity checking in `_dev/tests/session-start-hook-behavior.sh` by replacing Perl-based `shasum` with native POSIX `cksum`.

- Eliminates 20 Perl interpreter startup invocations across scenario runs, cutting CPU usage by 11.5% (-0.13s).
- Adds explicit mutation probes in `_dev/tests/session-start-hook-behavior.sh` verifying that byte changes, added files, and deleted files fail closed while preserving expected `actions/version.md` rewrites.
