# REQ-412 Builder Brief

Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-412-implement-request-state-transactions`
Branch: `worktree-agent-REQ-412-implement-request-state-transactions`
REQ: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/working/REQ-412-implement-request-state-transactions.md`
Plan: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-08-31-165510/REQ-412-plan.md`
Exploration: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-08-31-165510/REQ-412-exploration.md`

Implement exactly the frozen 20-file `write_set` and Scope in the claimed REQ. Key constraints:

- Start with real-command and focused RED fixtures for all five commands, then implement `internal/requeststate` as a plan-first/apply-once authority over existing repository, graph, request-document, Git-transaction, and result contracts.
- Preserve selection provenance: explicit `REQ-NNN` claims bypass dependency gating and clear assignment; default/UR-expanded claims remain gated. Exact selector-returned paths are validated, never rescanned.
- The shared Git expansion is a narrow explicit opt-in for snapshottable existing-untracked exact targets. Default guards, unrelated dirt, index restrictions, rollback semantics, and cleanup's separate exception must remain unchanged.
- Actions retain confirmation, failure classification, review/release judgment, follow-up creation, and dependent disposition. Go owns deterministic lifecycle bytes and paths only.
- Preserve the serial two-commit provenance rule and return committed-state risk/revert evidence after history exists; never amend/reset.
- Run focused requeststate/gittransaction tests, full CLI tests/vet, exact Go 1.25, contract and commit-hash guards, gofmt, diff/scope hygiene, and the canonical gate if practical in the worktree.

State stays home: do not edit or commit anything under `do-work/`, version/changelog files, primes outside the declared CLI prime, or unrelated files. Commit the 20 implementation files on the branch, keep the worktree clean, and write the handback only to `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-08-31-165510/REQ-412-handback.md`. Include commit hash, RED/GREEN evidence, commands/results, exact files, P-A-U evidence, decisions, discovered tasks, and readiness for integration.
