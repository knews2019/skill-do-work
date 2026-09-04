# REQ-504 Review Preparation

Ready for independent review after REQ-509. This is read-only preparation, not an approval: no tests ran, and no request or lifecycle records changed.

## Authority and exact range

- Request: `do-work/queue/REQ-504-collapse-step-10-and-crash-recovery-prose-into-recovery.md` (resolve its current exact path again after claim).
- Original input: `do-work/user-requests/UR-098/input.md`.
- Route C, maintenance rule change, `tdd: true`, depends on REQ-503. The user wants mechanics owned by commands and judgment retained in prose, with each migration carrying command changes, deleted prose, deleted predicates, and behavioral coverage.
- Base: `773787b74acddfdfc4c16498a89d99a5cc3ab716`.
- Implementation/evidence target: `f412a8411057d0a833df5584657161008f315b84`.
- `git diff --stat BASE..TARGET` reports 27 paths: the declared 26 implementation paths plus the request's own evidence record. Treat the latter as bookkeeping, not unexplained implementation scope drift.
- Loaded the review action, original input, complete REQ, action-files and shell-commands primes, CLI prime, maintenance and anti-slop crew. Review remains responsible for its own full code inspection and guardrail checks.

## Saved verification

The request records a successful canonical answer at `2026-09-04T21:00:44Z`; both `commit` and `heavy_verified_revision` equal the exact target. Saved and recomputed plans matched base, target, lane selection, argv, and reasons.

| Selected lane | Exit | Seconds |
|---|---:|---:|
| queue-kanban-javascript | 0 | 8 |
| queue-kanban-browser | 0 | 97 |
| do-work-cli-integrations | 0 | 60 |
| staged-skills | 0 | 23 |
| updater | 0 | 51 |
| installer | 0 | 23 |

All six passed without skips. Each command was `bash _dev/tests/maintainer-verify.sh --heavy-lane LANE`. Browser used `/Applications/Google Chrome.app/Contents/MacOS/Google Chrome` through `QUEUE_KANBAN_BROWSER`. Three shell lanes initially refused an invalid duration-log header before tests started; the original logs were retained, the canonical duration helper initialized a fresh log, and only those three lanes reran at the same revision. No source changed. Detailed prior-run artifacts remain under `.git/clarify-heavy-_p3jvjz7/`.

The implementation record additionally reports RED before production, focused five-package GREEN, race checks, vet, full module and contract success, and the direct canonical repository gate. Fresh heavy reruns are unnecessary solely to repeat this evidence.

## Focused review cases

The command, prose deletion, shell-lane deletion, and Go tests are all visible in the saved diff. The captured sentence-predicate RED premise had become stale before dispatch; the Plan explicitly substitutes live public-boundary RED cases and deletes `_dev/tests/contracts/request-state.sh` plus its aggregator registration after replacement coverage.

1. Trace public finalization-first recovery through committed and uncommitted claims, hostile writer data, duplicate same-request labels, unlabelled evidence, untouched unrelated bytes, canonical selection, and fresh claim. Existing acceptance case: `TestRecoverPublicCommandRunsFinalizationThenRecoversEveryClaim`.
2. Verify observation never resets claims; explicit takeover alone authorizes mutation. Existing case: `TestRecoverWithoutAuthorityOffersTypedTakeoverAndDoesNotMutateClaim`. Inspect invalid/missing/repeated authority arguments and refusal propagation; two public tests do not establish every branch by themselves.
3. Verify claim-only topology is exempted narrowly rather than hiding an actual interrupted finalization. Existing case: `TestRecoverFinalizationAcceptsCoherentClaimOnlyTopology`; inspect its boundary against the other discovery tests.
4. Check canonical-section versus legacy heading-less checkpoint discovery, all-entry planning, aliases and indented continuation removal, and unrelated-byte preservation. Relevant cases: `TestDiscoverRepositoryProjectsCheckpointClaimsInSourceOrder`, `TestRecoveryPlanAcceptsAllCheckpointEntriesAuthority`, and `TestAuthorizedCheckpointRemovalDropsEverySameRequestEntry`.
5. Check checkpoint mode changes only `do-work/CHECKPOINT.md`, preserving foreign/unlabelled live entries exactly. Existing case: `TestAdvanceCheckpointChangesOnlyCheckpointAndPreservesLiveEntries`. The saved foundation also has `TestOrdinaryAdvanceRemainsReadOnlyAfterCheckpointMode`.
6. Trace the seven inherited REQ-498/501 review instances to their implementation and test evidence. Sweep live consumers of removed checkpoint-template wording and recovery ownership. Step 10 contains one loop paragraph plus a separate context-wipe principle; the staged exception leaving selection in `next` is explicitly documented because REQ-505 owns the next migration.

A proportional fresh acceptance run in a detached checkout at the target can select the named public recovery and checkpoint tests, plus the claim-only finalization test. The prime maps wider validation to `go test ./internal/<package>`; relevant owners are lifecycleadvance, finalization, requeststate, repositorymodel, and resultmodel. Any checkout created for review should be removed after its processes finish.

## Current-version caveats

Judge this saved migration against its range. Later REQ-505 deliberately makes queue-mode `advance REQ` select and claim, REQ-506 adds evidence gates, and REQ-507 adds finalization. Running the saved read-only queue test against current HEAD would test a different contract.

REQ-515 later changes finalization recovery to set aside per-request refusals and preserves those claims through public recovery; current `recovery_commands.go` adds `setAsideRequestIDs` for this purpose. The saved REQ-504 decision to stop after a refused request-state transaction is recorded as D-07 and should not be confused with that later finalization behavior. Current work.md has also received later gate and heavy-evidence changes, including REQ-559/560/564.

Report any verified residual with impact and exact evidence. This preparation has not established a finding or a review score.
