# Shared Principles Crew Member

<!-- JIT_CONTEXT: Load during implementation and review, regardless of domain or review mode. These conditions consolidate portable guidance from the work and review Common Rationalizations tables (REQ-509). Action-specific lifecycle and authority boundaries remain in the invoked action. -->

## Principles

| When... | Apply this principle |
|---|---|
| A completion claim is written or reviewed | List every changed file with its action and a factual summary; verify those claims against the actual diff and REQ requirements. |
| The REQ has `tdd: true` | Show the failing test before implementation. Follow [testing.md](./testing.md); a test written afterward is not RED-before-GREEN evidence. |
| A request requiring a Scope declaration will touch a file | Declare that path in Scope before coding, however small the change. |
| A test failure reaches its second repair attempt or later | Load [debugging.md](./debugging.md) and [testing.md](./testing.md) before retrying. |
| Adjacent work falls outside the declared request | Follow [general.md](./general.md#discovered-tasks-contract): record the discovery instead of widening the implementation. |
| Acceptance is inferred from narrower tests | Exercise the requested behavior end-to-end where applicable; unit tests alone do not establish acceptance. |
| Acceptance cannot be exercised | Record Untested and the exact check that could not run. |
| Requirements and tests pass | Apply the Klarna Test: check the user's intended outcome, including what the checklists and tests did not measure. |
