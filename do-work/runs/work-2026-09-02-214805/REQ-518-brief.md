# Builder Brief — REQ-518

Implement `do-work/working/REQ-518-run-the-full-gate-once-per-req.md` against baseline `bab2198d59e891ec23d98a209c3c03187bc1741d`.

Use the saved plan and exploration in this run directory. Preserve the unrelated untracked AI report. Stay within the REQ's `## Scope`; do not edit the REQ, checkpoint, run manifest, release/version files, finalization/request-state packages, or gate scripts.

TDD is mandatory: add the public CLI integration test first, run it to an assertion-level `UNKNOWN-COMMAND` failure, retain the RED output, then implement and rerun GREEN. Run focused tests and report every changed file, exact RED/GREEN commands and outcomes, P-A-U evidence, decisions, lessons consulted, and any discovered tasks to `REQ-518-handback.md`.
