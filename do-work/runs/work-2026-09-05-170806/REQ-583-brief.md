# Builder brief — REQ-583

## Where you work

- **Worktree:** `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-583-pin-the-evidence-gate-remedy-redirection-guard-and-interrupted-path`
- **Branch:** `worktree-agent-REQ-583-pin-the-evidence-gate-remedy-redirection-guard-and-interrupted-path` (already checked out there, cut from the integration tip)
- **Commit on that branch.** Do not touch the main tree at `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2` — with exactly one exception, the hand-back file named below.

## Never touch

- Any path under `do-work/` in the main tree or in your worktree, except the one hand-back file.
- `internal/lifecycleadvance/evidence_gates.go` and `internal/corehelpers/commands.go`. The REQ's own constraint is tests only; the three behaviours do not change. You will *temporarily* mutate them to prove RED, then revert. `git diff` at the end must show `evidence_gates_test.go` and nothing else.
- Any other test file, including `internal/corehelpers/commands_test.go` and `internal/nextselection/blocked_probe_test.go`.

## Hand-back

Write `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-09-05-170806/REQ-583-handback.md` (absolute path, main tree — the one exception). Never stage or commit it.

Sections it must carry:
- `## File manifest` — every file created/modified/deleted with the verb, plus your branch head commit.
- `## Red-Green Evidence` — for each of the three mutations: the exact mutation, the named test that failed under it, the failure message, and the pass state after revert. This is the REQ's whole point; a summary without named tests and messages is not evidence.
- `## Decisions` — numbered `D-01`, `D-02`, … Sort each by the decide-vs-escalate gate in `crew-members/coding-guardrails.md` § Think Before Coding. A reversible low-reach call is DECIDE & STATE (reasoning only); an irreversible, taste-dependent or genuinely contestable one is ESCALATE and additionally carries `Value:` and `Risk:` lines.
- `## Discovered Tasks` — anything out of scope you found. Do not fix inline.
- `## Lesson evidence` — which lesson satellites you read, whole or family-targeted, and any that were missing.
- `## Integration seams` — any line that belongs in a file outside your write set. Hand back the exact line and where it goes; never edit the shared file yourself.
- `## Exploration` — only if you discover something the orchestrator's `## Exploration` section got wrong.

## Timing budget

`internal/lifecycleadvance` is under a 30-second per-test-file budget (REQ-574). The interrupted-probe row costs about 0.5s; the `focusedGateState` table costs microseconds. Report the package wall time before and after.
