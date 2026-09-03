# REQ-490 builder brief

Fix `queueDependencyDepth` so a dependency already resolved as met by the graph's duplicate-status aggregation contributes no depth. Add the exact two-completed-duplicates plus pending-dependent `--wave 0` RED/GREEN lock-in while preserving single-record, missing, cyclic, and ambiguous-unsatisfied behavior. Builder source only; lifecycle, run, release, and queue files remain orchestrator-owned. Hand back branch, exact commit, changed files, tests, risks, and merge guidance.
