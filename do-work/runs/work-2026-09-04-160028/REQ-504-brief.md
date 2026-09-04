# Builder Brief — REQ-504

Implement the recovery/checkpoint prose collapse on branch `worktree-agent-REQ-504-recovery-collapse` in worktree `/Users/t2/Desktop/e1-experimental-repos/worktree-agent-REQ-504-recovery-collapse`.

Read the committed request at `do-work/working/REQ-504-collapse-step-10-and-crash-recovery-prose-into-recovery.md`, both listed primes, the CLI prime, and `skills/do-work/crew-members/general.md`, `coding-guardrails.md`, `communication-style.md`, `backend.md`, and `testing.md`. Follow the committed Plan, Exploration, Scope, and inherited review instances.

The write boundary is exactly the request frontmatter's 26 paths. Do not modify any other project path and never modify a `do-work/` path in the builder worktree. Retain the exact `## Crash Recovery (Step 1)` heading so doctor/forensics/cleanup/board consumers remain valid. Preserve the concurrent release commit `be1b67d9` changes already present at the base, including the new single-line machine-hold wording.

Use strict RED/GREEN TDD in the plan's sequence. The capture-time shell-predicate RED is stale; do not manufacture it. First add failing Go tests for the live gaps: coherent claim-only finalization topology, complete structural checkpoint evidence, atomic all-entry recovery, public recovery order, hostile labels, authority/takeover rules, checkpoint-only advance mutation, and unchanged ordinary advance. Capture actual failure output before implementation.

Implement the recommended boundary:

- canonical `recover [--assume-sole-authority] [--take-over REQ-NNN]` composition in `lifecycleadvance`;
- ordinary `advance REQ-NNN` stays byte-for-byte read-only;
- explicit `advance --checkpoint` is its sole mutation and may change only `do-work/CHECKPOINT.md`;
- `run-with-recovery` invokes constant argv and never interpolates an observed writer;
- selection/claim stay with typed `next`/`claim` for this REQ because REQ-505 owns moving them behind `advance`;
- explicit authority atomically removes every same-REQ checkpoint entry, including multiple labels and unlabelled legacy entries, while preserving unrelated entries and project dirt;
- the shell request-state lane and aggregator registration are removed only after equivalent Go/public-binary coverage passes.

Keep action prose at principles/judgment altitude. Rename Step 8 to preparation, correct the orchestrator/finalization ownership restatements, collapse Step 10 to one paragraph plus the short context-wipe principle, collapse the reference recovery/checkpoint algorithms, and update commit/handoff/guide/prime restatements. Do not move judgment into typed prose output.

Run focused RED/GREEN packages, public binary seams, `go test -race` for mutated transaction packages where practical, `go vet ./...`, the uncached CLI module suite, `bash _dev/tests/contract-regressions.sh`, `git diff --check`, and `bash _dev/tests/maintainer-verify.sh`. Respect the per-test-file 30-second budget and report any baseline distinction.

Commit exactly the 26-path implementation. Write a complete handback with `apply_patch` to `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-09-04-160028/REQ-504-handback.md`: branch/commit, exact file manifest, RED/GREEN evidence and timings, public-order evidence, P-A-U APPLY/UNIFY, lessons, decisions, discovered tasks, deleted-shell-coverage mapping, prose-size result, and integration seams. Return only one status line after both commit and handback exist.
