# REQ-629 Exploration: handoff resume vs recover takeover

Paths below are relative to `skills/do-work/tools/do-work-cli/internal/` unless they start with `skills/` or `do-work/`.

## F1. Where `recover` emits its claim findings

- Handler: `lifecycleadvance/recovery_commands.go:26` `handleRecover`.
- Candidates: `recoverableWorkingRequests` (`recovery_commands.go:232-248`), every `do-work/working/` request with `status: claimed` (or `blocked` with `blocked_by`).
- Per-claim decision loop `recovery_commands.go:80-156`:
  - `:85-89` finalization set-aside, claim preserved (no finding).
  - `:90` **the authority predicate:** `authorized := options.assumeSoleAuthority || options.takeOverRequestID == requestID`. Nothing else.
  - `:91-102` not authorized: decision `"takeover available; claim preserved"`, finding `RECOVERY-TAKEOVER-AVAILABLE` (warning, manual), `AutomationStopReason: "working claim requires explicit authority"`, `NextArgv: ["do-work-cli","recover","--take-over",REQ]` (`:99`), verification argv plain `recover`.
  - `:104-118` authorized but `heldForHeavyLanes` (`:178-188`: `## Heavy Verification Plan` section and `commit:` is an ancestor of HEAD): `RECOVERY-CLAIM-HELD-FOR-HEAVY-LANES`, next argv `advance REQ`. This is the only existing "keep a merged claim" protection, and it only runs after authority is given.
  - `:120-148` otherwise: `requeststate` `TransitionRecover` canonical reset, `RECOVERY-CLAIM-RECOVERED`.
  - `:157-163` no working claims: `RECOVERY-NONE`.
  - Usage/not-found failures: `RECOVERY-USAGE` (`:29`), `RECOVERY-TAKEOVER-NOT-FOUND` (`:77`).
- The reset strips sections via `requeststate/state_apply.go:984-1000` `stripGeneratedRecoverySections`, heading regex at `:984`: Triage, Exploration, Plan, Scope, Pre-Flight, Implementation Summary, Qualification, Testing, Review, Lessons Learned, Orientation, Decisions, Discovered Tasks, Timing. Wired at `state_apply.go:579-606`.

**Same-writer vs foreign: there is no distinction at all.** `recover` never compares the checkpoint writer label with the current session. The writer label is only reported as evidence text (`recoveryEvidenceLabels`, `:261-274`, "checkpoint line N writer: X"). No session id, no process liveness, no hostname check. The writer label that claims write is `hostname + ":" + repositoryRoot` (`lifecycleadvance/queue_commands.go:424-430` `queueWriterLabel`), so a same-machine restart in the same checkout produces an identical label, and `recover` still offers `--take-over`.

Tests that pin it (`lifecycleadvance/recovery_commands_test.go`):
- `:115` `TestRecoverWithoutAuthorityOffersTypedTakeoverAndDoesNotMutateClaim`: checks the code, the `"--take-over"` token, that observe mode changes no bytes, and that `--take-over` moves the REQ back to the queue. Fixture writer is `other:/checkout` (foreign), so no test covers the same-writer case.
- `:15` `TestRecoverPublicCommandRunsFinalizationThenRecoversEveryClaim` (sole-authority path).
- `:191` `TestRecoverLegacyCheckpointClaimsThroughPublicCommand` (take-over with legacy/multi-writer entries).
- `:248` `TestRecoverHoldsAClaimedRequestWithAHeavyVerificationPlanForTheDrain`.
- `recovery_set_aside_test.go:20` `TestRecoverPreservesTheClaimOfARequestFinalizationSetAside`.
- `resultmodel/result_model_test.go:464-477` renders a take-over result in text mode (format only).

## F2. Does per-request `advance` continue a claim without takeover?

Code path: `lifecycleadvance/advance_commands.go:37-71` `handleAdvance`. For `advance REQ-NNN` with the REQ in `working/`, it calls `classifyAdvance` (`:80`) -> `classifyWorkingAdvance` (`:127`). There is **no ownership, writer or checkpoint check**; it only requires `status: claimed|completed-with-issues` (`:130`). So yes, a resuming session can continue the claim with no takeover.

But the argv form matters:
- `advance REQ-NNN` (bare, one argument): read-only classifier, returns the next phase (`:61-63`). Works at every phase.
- `advance REQ-NNN --request-path <path>` with nothing after: goes to `executeAdvanceEvidenceGates` (`:68`). At an agent-judgment phase it is **refused**, exit 1, `ADVANCE-GATE-INPUT-IRRELEVANT` "the current lifecycle phase requires agent judgment, not gate argv" (`lifecycleadvance/evidence_gates.go:96-97`, `:275-276`). The `--request-path` form is only for mechanical gate phases (estimate, gates) and is what the classifier itself returns as `next_argv` there.

The UR-132 restart prompt (commit 1794c268, `do-work/RESTART-PROMPT.md`) told the session to "ask the classifier for its next phase (advance REQ-NNN --request-path <its do-work/working path>)". That form refuses at judgment phases, so the session had no working command from the prompt, and `recover`'s next argv was the only offered step.

### Reproduction (scratch repo `.../scratchpad/req629-repro`, suite copied from this tree)

1. Queue-mode `advance` claimed REQ-001 (commit `[REQ-001] claim request lifecycle`), checkpoint entry `writer: t2s-Virtual-Machine.local:<repro root>`.
2. Added `route: B` + `## Triage`; `advance REQ-001` -> `estimate-p50` (mechanical, next argv `advance REQ-001 --request-path ... -- --route B --write-set <count> ...`). Running that argv returned a satisfied `estimate-p50` gate record (prints p50 10, writes nothing).
3. Wrote `estimate:`, `## Plan`, `## Exploration`, `## Scope`; committed (simulated handoff).
4. "Fresh session", same host, same root (identical writer label):
   - `recover` -> success, `FINALIZATION-NONE`, `RECOVERY-TAKEOVER-AVAILABLE` warning, next argv `do-work-cli recover --take-over REQ-001`, evidence shows the same writer label. No bytes changed.
   - `advance REQ-001 --request-path do-work/working/REQ-001-sample.md` -> exit 1, `refused`, `ADVANCE-GATE-INPUT-IRRELEVANT` (phase was `agent judgment: scope declaration`).
   - `advance REQ-001` -> success, phase `agent judgment: scope declaration`. Claim continues.
   - queue-mode `advance` -> success, `frozen_members: []`, `claimed: []`, no findings, no commit, REQ file byte-identical.
5. In a clone, `recover --take-over REQ-001` -> `RECOVERY-CLAIM-RECOVERED`, commit `[REQ-001] recover-claim request lifecycle`: REQ moved to `do-work/queue/` as `pending`, `route` removed, Triage/Plan/Exploration/Scope removed (estimate and `claimed_at` kept), checkpoint entry removed. RED case from the REQ is confirmed.

## F3. Next queue-mode `advance` after a plain recover

It ignores the claim. Selection only considers `queue/` files (`nextselection/next_selection.go:371` skips any `TreeSection != "queue"`), so a working claim is neither selected nor listed as excluded. Result: empty `frozen_members`, nothing claimed, no finding, no mutation (repro step 4). Consequence: `do-work run` alone, after a plain recover, sees an empty queue and the claimed REQ is never resumed unless the agent continues it explicitly (or it is held for heavy lanes, `skills/do-work/actions/work.md:429`). The only typed pointer `recover` gives for that REQ is `--take-over`.

## F4. Shipped prose

`skills/do-work/actions/restart-with-parallel-handoff.md`:
- `:49` Step 3 asks the handoff to record "The canonical `recover` result for each working claim, including its structural writer evidence and exact takeover command." So the action itself tells the writer to put the takeover command into the handoff.
- `:61-66` paste block: line one is the resume command (`do-work run --fan-out N`, or `do-work clarify`), then "This command is sufficient; everything below it is context." Nothing about continuing claimed REQs or about takeover resetting them. The "run canonical recover, then continue each claimed REQ" wording was the UR-132 session's own paste block, not action text.

Other shipped mentions of take-over (grep `take-over|takeover` under `skills/`):
- `skills/do-work/actions/work-reference.md:310` Crash Recovery: "plain recovery preserves claims and returns exact takeover argv, while explicit `--take-over REQ-NNN` or `--assume-sole-authority` authorizes the canonical reset". "Follow the result" — which recommends the takeover argv.
- `skills/do-work/docs/work-guide.md:76`, `:130`, `:165` (`:165`: "its first canonical recovery result says whether anything needs takeover authority").
- `skills/do-work/actions/run-with-recovery.md:41` (reports takeovers; rwr uses `--assume-sole-authority`, an intentional reset).
- `skills/do-work/actions/forensics.md:54`, `:100`, `:168` (point to Crash Recovery for takeover judgment).
- `skills/do-work/actions/work.md:126` (run recover, then queue-mode advance).

## F5. Smallest command change

File: `lifecycleadvance/recovery_commands.go`, function `handleRecover`, the not-authorized branch `:91-102`.

Option A (smallest, no ownership inference): keep the code and the decision, but change the finding so it cannot be read as "resume":
- `NextArgv` -> `["do-work-cli","--format","json","advance",REQ]` (the read-only classifier, same shape as the heavy-lane hold at `:115`), i.e. continue the claim.
- Put the reset in the stop reason / evidence, e.g. `AutomationStopReason: "working claim preserved; continue it with advance, or reset it to the queue with recover --take-over REQ (strips orchestrator sections)"`.
This needs no writer comparison, so it cannot misfire on a foreign claim: for a foreign claim, `advance REQ` is read-only and harmless, and the reset stays one explicit command away.

Option B (writer-aware): compare each checkpoint entry's `Writer` with `queueWriterLabel(executionContext.RepositoryRoot)` (`queue_commands.go:424`); when every labelled entry matches, emit a new info code (e.g. `RECOVERY-CLAIM-RESUMABLE`) with next argv `advance REQ`; keep `RECOVERY-TAKEOVER-AVAILABLE` for foreign/unlabelled. Risk: hostname:root is not a session identity — a crashed sibling session in the same checkout has the same label, and this repo's memory records real same-checkout concurrent sessions. Option A avoids that judgment.

Either way, a merged REQ (`commit:` / `integration_at`) should arguably never be reset by take-over without a further flag; today only the heavy-plan case is protected (`:104`).

Tests to change/add:
- Update `recovery_commands_test.go:115` `TestRecoverWithoutAuthorityOffersTypedTakeoverAndDoesNotMutateClaim`: it asserts the output contains `"--take-over"`; under Option A it should assert `next_argv` is `advance REQ-713` and that the stop reason names the reset. Keep its take-over half (the reset still works when asked).
- New focused test: same-writer claim with Triage+Scope, plain `recover` then `advance REQ` -> claim stays in `working/`, bytes unchanged, classifier names the next phase (the REQ's GREEN condition).
- Prose: `restart-with-parallel-handoff.md:49` (drop "exact takeover command"; say continue each claim with bare `advance REQ-NNN`, and that `recover --take-over` resets it), `work-reference.md:310` (same sentence), `docs/work-guide.md:165`.
