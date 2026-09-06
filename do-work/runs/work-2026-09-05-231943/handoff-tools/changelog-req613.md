## 0.305.24 — Reduce Fast-Stage Evidence Computation Cost (2026-09-06)

Profiled and simplified fast-stage evidence computation in `heavy-runtime-fingerprint.py` and `do-work-cli` without weakening invalidation.

- Added in-memory caching of resolved binary seals in `_dev/tests/heavy-runtime-fingerprint.py`, eliminating redundant disk reads and SHA-256 hashing of identical executables.
- Consolidated working tree enumeration in `fast_stage_evidence.go` into a single `git ls-files` subprocess with combined `--cached` and `--others` flags.
- Scoped ignored file queries to stage coverage directory roots using `fastStageCoverageRoots`, eliminating full-tree ignored-file traversals while strictly preserving `laneCoversPath` checks.
- Preallocated internal slices and maps for stage seals to avoid dynamic slice reallocations.
- Added comprehensive unit tests and end-to-end behavior probes verifying that ignored files under stage coverage force execution while ignored files outside coverage preserve reuse.
