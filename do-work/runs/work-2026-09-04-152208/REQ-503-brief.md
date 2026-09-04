# Builder Brief — REQ-503

Implement the read-only `advance` lifecycle command on branch `worktree-agent-REQ-503-read-only-advance` in worktree `/Users/t2/Desktop/e1-experimental-repos/worktree-agent-REQ-503-read-only-advance`.

Read the committed request at `do-work/working/REQ-503-add-read-only-advance-lifecycle-command.md`, `skills/do-work/tools/do-work-cli/prime-do-work-cli.md`, and `skills/do-work/crew-members/general.md`, `coding-guardrails.md`, `communication-style.md`, `backend.md`, and `testing.md`.

The write boundary is exactly these six paths:

- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands.go` (new)
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/advance_commands_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (modify)
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model_test.go` (modify)
- `skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go` (modify)
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (modify)

Do not modify any `do-work/` path in the builder worktree. Do not touch action prose, contract-regression scripts, repository/request parsing, or existing lifecycle packages. Follow the committed Plan and Exploration. The current source of truth makes targeted `next` the blocked-probe owner; do not recreate the obsolete `run-blocked-check` composition.

Use strict RED/GREEN TDD: write and run the lifecycle phase/refusal/read-only tests before implementation, capture the actual failing output, then implement. Cover Route A/B/C, queue/working/archive identities, malformed/ambiguous/contradictory/impossible states, exact typed projection in JSON/text, and byte-for-byte no mutation. Keep `advance` read-only and snapshot-driven; it must not invoke mutating sibling handlers.

Run focused new-package and result-model tests, registration tests, `go vet ./...`, the uncached CLI module suite, `git diff --check`, and the canonical gate `bash _dev/tests/maintainer-verify.sh`. The preflight baseline has one unrelated existing `internal/corehelpers` failure: `TestProtectedInventoryPersistsLaterXAndRequiresStartedState`; report whether it remains identical and require all newly affected tests to pass.

Commit the exact six-file implementation on the builder branch. Write the complete handback with `apply_patch` to the main-tree path `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-09-04-152208/REQ-503-handback.md`. Include branch/commit, exact file manifest, RED/GREEN evidence, commands and timings, P-A-U APPLY/UNIFY evidence, required lessons read/missing, decisions, discovered tasks, baseline comparison, and integration seams. Return only a one-line status after the handback exists.
