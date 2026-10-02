---
id: REQ-629
title: '[impact-rule-change] Resuming a handoff must not reset its own claimed REQs'
status: completed
route: B
estimate:
  p50_active_minutes: 20
  confidence: medium
  basis:
  - Route B
  - 4-file write set
  - 2 subsystems involved
  - 4 acceptance criteria
  calculated_at: 2026-10-02T20:55:30Z
created_at: 2026-10-02T19:49:03Z
user_request: UR-133
domain: general
prime_files: ["_dev/primes/prime-action-files.md"]
tdd: false
maintenance: false
impact: impact-rule-change
effort_estimate: effort-substantive
related: ["REQ-628"]
batch: ur-132-follow-ups
required_lessons: ["skills/do-work/tools/do-work-cli/lessons-do-work-cli.md#rule-direction-checked-against-callers"]
write_set: ["skills/do-work/tools/do-work-cli/internal/lifecycleadvance/recovery_commands.go", "skills/do-work/tools/do-work-cli/internal/lifecycleadvance/recovery_commands_test.go", "skills/do-work/actions/restart-with-parallel-handoff.md", "skills/do-work/actions/work-reference.md", "skills/do-work/docs/work-guide.md"]
dispatch_at: 2026-10-02T20:57:11Z
builder_handback_at: 2026-10-02T21:00:30Z
integration_at: 2026-10-02T21:01:19Z
review_at: 2026-10-02T21:06:50Z
heavy_verified_at: 2026-10-02T21:14:49Z
heavy_verified_revision: 069f922e64c1ed14ca31211642e395a784561f57
claimed_at: 2026-10-02T20:53:37Z
commit: 2cbd685813412cab0b51a783a0fbd6b25a4859a8
completed_at: 2026-10-02T21:15:31Z
release_at: 2026-10-02T21:15:31Z
---
# Resuming a Handoff Must Not Reset Its Own Claimed REQs

## What
Make the handoff (`phandoff`, `actions/restart-with-parallel-handoff.md`) and the crash-recovery takeover agree on what "resume" means, so a session that resumes a handoff continues the handed-off claimed REQs from their written sections instead of resetting them.

## Why
On 2026-10-02 the UR-132 run was handed off with REQ-626 claimed, merged at 505ec74a and gate-green, and REQ-627 claimed with Triage, estimate, Exploration and Scope written. The handoff prompt said "run canonical recover, then continue each claimed REQ". After a machine restart, `recover` reported `RECOVERY-TAKEOVER-AVAILABLE` with the argv `recover --take-over REQ-NNN` ("working claim requires explicit authority"). Running it, as the finding suggested, returned both REQs to `do-work/queue/` as `pending` and stripped every orchestrator-generated section, including route, write_set, Scope, Implementation Summary, Qualification, Testing and Decisions of an already-merged REQ. The run recovered only by re-claiming and restoring the sections byte-for-byte from the handoff commit 1794c268 (decisions REQ-626 D-07 and REQ-627 D-05 in `do-work/archive/UR-132/`). The takeover is documented as a "canonical reset" (`actions/work-reference.md` → Crash Recovery), so the command did what it says; the handoff and the finding's suggested next step led the resuming session into it.

## Detailed Requirements
- Establish, by a real check, whether a resuming session on the same checkout can continue a handed-off claim with per-request `advance REQ-NNN --request-path <working path>` without any takeover, and what `recover`'s plain result then does on the next queue-mode `advance`.
- Fix the side that is wrong (see Open Questions) so that following the handoff's paste block plus the commands' own suggested next steps never resets a claim the handoff described as in flight.
- The handoff paste block (or the handoff action's rules for writing it) states what to run for each claimed REQ and that `recover --take-over` resets a claim, when that remains true after the fix.
- A REQ that is merged (`commit:` or `integration_at` present) must not lose its evidence sections through the documented resume path.

## Constraints
State is binding; prose is advisory (the handoff action's own rule): prefer a fix that the commands enforce over one that only adds prose, when the command change is small. No change to crash recovery for a genuinely crashed claim from another checkout. Shipped files change, so this is a release.

## Dependencies
None. Independent of REQ-628 (board guide update).

## Builder Guidance
Medium certainty on the fix direction; high certainty on the failure. Reproduce it on a scratch repository first: claim a REQ, write sections, write a handoff, then follow the paste block in a fresh process.

## Open Questions
- [~] Which side changes: the handoff instructions, or the takeover command? → **D-01**: Builder chose: Both, minimally — the handoff action tells the resuming session to continue its own claims with per-request advance and never to run takeover on them, and `recover`'s takeover finding stops presenting `--take-over` as a plain next step for that case or states plainly that it resets the claim. Reasoning: the trap was the suggested argv in the command output; a prose-only fix leaves it there, and a new resume mode is only needed if per-request advance cannot continue a claim, which exploration checks first. Value: the documented resume path and the command's own next step agree, so a resumed run keeps merged work. Risk: a reworded or narrowed finding could hide the takeover from a genuinely crashed claim; reversible by restoring the argv, and the REQ's constraint keeps crash recovery for another checkout unchanged.
  Recommended: Both, minimally. Also: prose only, or a new resume mode in `recover`.

<!-- D-XX counter: last used D-01. Next decision: D-02. -->

## Red-Green Proof
**RED prompt/case:** In a scratch repo, claim a REQ, append Triage and Scope, write and commit a handoff, then in a new session follow the paste block and the `recover` finding's next argv: the REQ ends up in `do-work/queue/` as `pending` with Triage and Scope gone.
**Why RED now:** The handoff says "continue each claimed REQ", and the only next step `recover` offers for that claim is `--take-over`, which is a canonical reset.
**GREEN when:** Following the paste block and the commands' suggested next steps leaves the REQ claimed in `do-work/working/` with its sections intact, and the classifier names its next phase.
**Validation:** Inferred during capture

## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` as a whole satellite (17879 tokens, over the 2000 budget; `slugged: partial`). Matched: semantic recovery completeness, structured evidence projection. Narrowed at claim time to its `rule-direction-checked-against-callers` family (changing a finding's next step must be checked against every shipped caller that relies on it), which is now in `required_lessons`.
- `_dev/primes/lessons-action-files.md` (5879 tokens, over budget; `slugged: partial`). Matched: action routing, downstream readers. No family names the handoff-versus-command contract, so no narrower entry exists.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.) Builder: change only the not-authorized branch's next argv and stop reason, pin it test-first in the existing test plus one UR-132 test, align the two prose files after grepping every caller.
- [x] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.) Builder: exactly the four Scope files, tests seen failing first. Orchestrator then added `docs/work-guide.md` on the same branch (D-04).
- [x] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.) Builder: 4 files, 93 insertions, 10 deletions; gofmt and go vet clean; lifecycleadvance tests ok in 24.2s; recovery-set-aside.sh and core-checks.sh contracts pass; no test pins the handoff wording. Orchestrator: work-guide.md diff re-read, 2 lines.

## Full Context
See `do-work/user-requests/UR-133/input.md` for complete verbatim input.

*Source: capture both as REQs (the handoff-versus-takeover mismatch reported at the end of the UR-132 run)*

---

## Triage

**Route: B** - Medium

**Reasoning:** The failure and the desired outcome are clear and the Open Question has a recommended direction, but where `recover` decides a claim needs takeover, whether per-request advance can continue a claim without it, and what the handoff paste block says need discovery before dispatch. Two surfaces (one action file, one Go finding), not an architectural change.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

Full findings with file:line anchors and a scratch-repo reproduction: `do-work/runs/work-2026-10-02-204757/REQ-629-exploration.md`. The findings that change what the builder writes (paths under `skills/do-work/tools/do-work-cli/internal/`):

- **The authority check has no writer comparison.** `lifecycleadvance/recovery_commands.go:90` authorizes only `--assume-sole-authority` or `--take-over <this REQ>`; otherwise `:91-102` emits `RECOVERY-TAKEOVER-AVAILABLE` whose next argv is `recover --take-over REQ`. The writer label (`hostname:repoRoot`) is evidence text only, so after a restart on the same machine the takeover is still the only offered step. `recovery_commands_test.go:115` pins `"--take-over"` in the output with a foreign-writer fixture.
- **Bare `advance REQ-NNN` continues a claim; `advance REQ --request-path P` is refused at judgment phases** (`ADVANCE-GATE-INPUT-IRRELEVANT`). Reproduced in a scratch repo: plain `recover` changed nothing, bare `advance REQ-001` named the next phase, and `recover --take-over REQ-001` requeued the REQ and stripped route, Triage, Plan, Exploration, Scope and its checkpoint entry. The UR-132 handoff prompt used the refused `--request-path` form.
- **Queue-mode `advance` ignores working claims** (`nextselection/next_selection.go:371`), so a plain `do-work run` sees an empty queue and never resumes a claim; the resuming session must name each claim.
- **Prose:** `skills/do-work/actions/restart-with-parallel-handoff.md:49` tells the handoff writer to record each claim's "exact takeover command"; its paste-block rules say nothing about claimed REQs. `skills/do-work/actions/work-reference.md:310` (Crash Recovery) says "Follow the result". Other mentions (`docs/work-guide.md`, `run-with-recovery.md`, `forensics.md`) describe genuine crash recovery.
- **Smallest command fix:** in the not-authorized branch, make the next argv the read-only `advance REQ` (the heavy-lanes hold at `:115` already uses that shape) and make the stop reason say that `recover --take-over REQ` resets the claim and strips its sections. No writer comparison, so it cannot misfire on a foreign claim; `--take-over` stays available and documented for a genuinely crashed claim.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/recovery_commands.go` (modify) — the takeover finding's next argv becomes read-only advance REQ; the stop reason names --take-over as a reset
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/recovery_commands_test.go` (modify) — update the pinned argv; add a test that a same-writer claim survives plain recover then advance REQ byte-identical in working/
- `skills/do-work/actions/restart-with-parallel-handoff.md` (modify) — record each claim's continue command (bare advance REQ-NNN), not a takeover; the paste block says what to run per claim and that --take-over resets a claim
- `skills/do-work/docs/work-guide.md` (modify, added by orchestrator decision D-04) — the recovery paragraph and the Context-limits bullet stop steering a resuming session to takeover
- `skills/do-work/actions/work-reference.md` (modify) — Crash Recovery says the takeover finding's next step is read-only and that --take-over is a reset for a genuinely abandoned claim

**Files I will NOT touch:** the reset itself (`requeststate/state_apply.go`), queue selection, `run-with-recovery.md`, `forensics.md` (they describe genuine crash recovery; the builder reports any sentence that tells a resuming session to take over), release paths, anything under `do-work/`.

**Acceptance criteria (restated from REQ):**
- [ ] A real check establishes whether per-request `advance` continues a handed-off claim without takeover, and what plain `recover` then does
- [ ] Following the handoff paste block plus the commands' suggested next steps never resets a claim the handoff described as in flight
- [ ] The handoff states what to run for each claimed REQ and that `recover --take-over` resets a claim
- [ ] A merged REQ (`commit:` or `integration_at`) keeps its evidence sections through the documented resume path

## Pre-Flight

**Git:** ✓ Integration tip 5f3038aa on `main`; the only dirt is this run's untracked `do-work/runs/` trail — no third-party paths
**Repository gate:** ✓ `bash _dev/tests/maintainer-verify.sh` exit 0 at 5f3038aa, gate wall 66s (do-work-cli fast stage reused by fingerprint)
**Baseline:** ✓ `bash do-work/runs/work-2026-10-02-204757/helpers/probe-629.sh` (`go test -count=1 -run Recover ./internal/lifecycleadvance/`) green. A first attempt passed the probe path as the test argv and could not launch (status 127, rerun 127); re-run with `bash <probe>` and satisfied.
**Dependencies:** ✓ Go toolchain present; no new dependency planned

*Checked by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/recovery_commands.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/lifecycleadvance/recovery_commands_test.go` (modified)
- `skills/do-work/actions/restart-with-parallel-handoff.md` (modified)
- `skills/do-work/actions/work-reference.md` (modified)
- `skills/do-work/docs/work-guide.md` (modified)

**What was done:** `recover`'s `RECOVERY-TAKEOVER-AVAILABLE` finding now offers the read-only `advance REQ-NNN` as its next argv, and its stop reason says `recover --take-over REQ-NNN` resets the claim (requeues it, strips its orchestrator sections) and is only for a claim no live session owns; `--take-over` itself is unchanged. The handoff action records each claim's continue command instead of a takeover command, and its paste block carries one `advance REQ-NNN` line per working claim with a warning that takeover resets a claim. Crash Recovery in work-reference and two work-guide sentences say the same. A new test follows recover's next step on a merged same-machine claim and proves the file stays byte-identical in `working/`. After review, the handoff's paste block lists the claim lines before the resume command and drops the one-command framing, and the guide's continue-vs-run bullet no longer says the checkpoint resumes claims (D-05). Cumulative merge range 89145e03..2cbd6858 (builder commit 42c37e7e, orchestrator commits on the same branch, merges 730c22cf and 2cbd6858).

## Qualification

**Diff range:** 89145e03..2cbd6858 cumulative (builder commit 42c37e7e, orchestrator commits, merges 730c22cf and 2cbd6858 — the second carries the D-05 review fixes)
**Gate records:** qualify satisfied; scope-drift satisfied (five files, exactly the Scope after D-04). A first scope-drift call reported every Scope entry as untouched because backticked non-path words in the Scope descriptions were read as paths; the backticks were removed from those descriptions and the re-run was satisfied.
**Warnings judged:** none remaining.
**Orchestrator read of the diff:** every Detailed Requirement traces. (1) The real check was run twice in scratch repos (exploration and builder): plain recover changes nothing, bare `advance REQ-NNN` names the next phase, `--take-over` resets. (2) recover's next argv is now that read-only advance, and the handoff paste block carries one `advance REQ-NNN` line per claim, so following the paste block plus each command's next step never resets a claim. (3) The handoff states what to run per claim and that takeover resets a claim. (4) The new test proves a merged claim (`commit:` set) stays byte-identical through recover then its next step. Crash recovery for another checkout is unchanged: `--take-over` and `--assume-sole-authority` behave as before.
**P-A-U honesty:** the boxes were ticked by the orchestrator from the builder's hand-back; APPLY cross-checked against `git diff --stat 89145e03..730c22cf` (five files, nothing under do-work/).

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh` on the merged tree (merge 730c22cf, REQ trail commit on top)
**Result:** ✓ All passing — exit 0, gate wall 129s; stage queue-kanban-fast-tests EXECUTING (413 tests, slowest file 19.63s < 30s); stage do-work-cli-fast-tests EXECUTING (867 tests, slowest file 20.86s < 30s). Green-gate record satisfied by advance.

**Focused tests:** `do-work/runs/work-2026-10-02-204757/helpers/probe-629.sh` (`go test -count=1 -run Recover ./internal/lifecycleadvance/`) → exit 0 (advance probe record satisfied). Builder: whole `./internal/lifecycleadvance/` package ok in 24.2s (< 30s); gofmt and go vet clean; `_dev/tests/contracts/recovery-set-aside.sh` and `core-checks.sh` pass.

**Red-green validation:** traced to `## Red-Green Proof`; tests written first on the builder branch, GREEN at 42c37e7e:
- TestRecoverWithoutAuthorityOffersTypedTakeoverAndDoesNotMutateClaim (updated): ✗ `takeover finding next argv = [do-work-cli recover --take-over REQ-713], want read-only advance` → ✓
- TestRecoverNextStepContinuesAClaimWithoutResettingIt (new): ✗ following the offered next argv ran `recover --take-over`, moved the merged claim `working/ -> queue/` and committed (the UR-132 failure) → ✓ claim byte-identical in `working/`, tree digest unchanged, classifier names `preflight`
- Scratch repo `req629-green` (the REQ's GREEN prompt end to end): handoff paste block written to the new rules, plain `recover` → next argv `advance REQ-001`; following it → phase `preflight` in `working/`; queue-mode advance → no mutation; REQ checksum unchanged, `git status` clean.

**Existing tests updated (cross-REQ impact):**
- `recovery_commands_test.go` TestRecoverWithoutAuthorityOffersTypedTakeoverAndDoesNotMutateClaim: pinned `--take-over` as the next argv; now pins read-only advance and the reset wording — intentional

**Heavy verification plan:** *(lanes selected by plan-heavy-verification)*
- Range: 89145e03..730c22cf
- do-work-cli-integrations — do-work-cli source changed
- staged-skills — shipped files under skills/ changed
- updater — do-work-cli source changed
- installer — do-work-cli source changed

**Repository gate after the D-05 review fixes:** `bash _dev/tests/maintainer-verify.sh` exit 0 at 64d4f665 (merge 2cbd6858 plus REQ trail), gate wall 123s, both Go stages EXECUTING; focused probe rerun exit 0.

*Verified by work action*

## Review

**Overall: 95%** | 2026-10-02T21:06:50Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 92% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Verdict: Approve.** `recover`'s takeover finding now offers the read-only `advance REQ-NNN` as its next step and names `--take-over` as a reset. The handoff action writes one `advance REQ-NNN` line per working claim. Following the paste block and the commands' own next steps no longer resets a claim. The remaining findings are wording drift, all report only.

**Requirements walk (merge range 89145e03..730c22cf, five files, matching Scope after D-04):**
- R1 Real check: delivered. The exploration and the builder each ran a scratch repo. This review ran a third one with a bare claim (no sections) under a foreign writer label (`otherhost:/elsewhere/repo`). Plain `recover` gave next_argv `do-work-cli --format json advance REQ-001`, and following it named the triage phase in `working/`. The file checksum was unchanged and `git status` was clean. Bare `advance REQ-NNN` with one argument only classifies and never writes (`advance_commands.go:61-62`), so it is read-only for a same-host claim, a foreign claim, a claim with no sections, and an abandoned claim alike.
- R2 Following the paste block plus the suggested next steps never resets a claim: delivered. `recovery_commands.go:98-105` changes only the not-authorized branch. `--take-over` and `--assume-sole-authority` are unchanged.
- R3 Handoff states what to run per claim and that takeover resets: delivered by the action's rules (`restart-with-parallel-handoff.md:49,66,125`). The REQ allows "or the handoff action's rules".
- R4 Merged REQ keeps its evidence sections: delivered. `TestRecoverNextStepContinuesAClaimWithoutResettingIt` uses a `commit:` fixture and checks byte identity plus the tree digest.
- Constraint (no change to crash recovery for a crashed foreign claim): the reset command and its authority are unchanged. The offered next step for that case did change, and D-01 records that risk.
- Consumers of next_argv: no Go code, `_dev/tests/` script, `run-with-recovery.md`, `forensics.md`, `work.md`, or the text renderer depends on the next_argv being `--take-over`. The renderer prints whatever argv it gets (checked: `next: do-work-cli --format json advance REQ-001`). The only pin was the updated test. `run-with-recovery.md` uses `--assume-sole-authority`, so the change does not affect it.
- Stop reason against `requeststate/state_apply.go:579-606,986-1001`: mostly accurate. The reset requeues the claim, deletes `route` (and `write_set` when Scope exists), and strips the generated sections. "As pending" is imprecise (see M2).
- `nextselection` claim: true for selection. `queueCandidates` and the explicit-target paths filter `TreeSection == "queue"` (`next_targets.go:26,69,122,146`). The exploration cited `next_selection.go:371`, which is `summarizeQueue`, not the filter. That is an anchor slip in a run artifact and nothing ships with it. The prose sentence still overreaches (see M1).
- P-A-U: all three boxes ticked. Decisions D-02..D-04 are recorded and match the diff.

**Restatement Sweep:** the diff redefines the takeover finding's next step and what resuming a claim means. I grepped `RECOVERY-TAKEOVER-AVAILABLE`, `take-over`, `takeover`, `resume` and `resume logic` across `skills/` and `_dev/tests/`. Agreeing: `work-reference.md:310`, `work-guide.md:76,130,165`, `forensics.md:54,100,168` (they delegate to Crash Recovery), `run-with-recovery.md:41` (sole-authority path), and the `queue-kanban/verify.go:75,818` comments (threshold, not next step). Stale: `work-guide.md:163` (M4), and the handoff action's own one-command framing at lines 3, 32, 82 and 124 (M3).

**Important findings:** None

**Minor findings:**
- M1 `restart-with-parallel-handoff.md:66` says "Selection reads only the queue, so the resume command never continues a claim". `work.md:429`'s heavy-lane drain does continue held claimed REQs at queue exhaustion, so the sentence is too broad. The only effect is a redundant read-only `advance` line. — impact-negligible → report only
- M2 The stop reason (`recovery_commands.go:101-103`) and `restart-with-parallel-handoff.md:66` say the reset "requeues it as pending". `state_apply.go:579-588,772-781` writes `pending-answers` when Open Questions has an unchecked `- [ ]` and keeps `blocked` for a blocked claim. It also deletes `route`/`write_set`, which the text does not mention. — impact-negligible → report only
- M3 Stale one-command framing in the same action. Line 3 says the next session "resumes from one pasted command", line 82 says "one-line resume", and the Step 1 test at line 32 plus the checklist at line 124 say "`do-work run` with no other reading must do the right thing". None of these holds once a claim exists, because the claim needs its paste line. Lines 66 and 125 mitigate this, but a writer who trusts Step 1's test can still skip the claim lines. — impact-rule-change → report only
- M4 `docs/work-guide.md:163` says "The checkpoint system handles the actual resume logic". Two lines below, the new line 165 says `do-work run` does not pick up working claims. A reader of `continue`/`resume` can still believe the run resumes a claim. — impact-user-visible → report only
- M5 The paste-block rule (`restart-with-parallel-handoff.md:61-66`) puts `do-work run --fan-out N` first and the `advance REQ-NNN` lines after it, without saying which runs first. A merged, gate-green claim (the UR-132 REQ-626 case) can wait behind a whole queue run, or fan-out builders can start before a claim's phase is done. — impact-negligible → report only

**Nit findings:**
- N1 The finding code `RECOVERY-TAKEOVER-AVAILABLE` and the claim decision "takeover available; claim preserved" (`recovery_commands.go:92,95`) still lead with takeover, while the next step is now continue. Keeping the code for compatibility is reasonable. — impact-negligible → report only

**Risk note:** a claim owned by a live session in another checkout now gets a continue next step. Before, it got a reset. A loop that follows next_argv automatically would work the same REQ as its live owner, which is less destructive than the old reset. The finding stays `fixability: manual`, and D-01 records this risk.

**Acceptance:** Pass. `go test -count=1 -run Recover ./internal/lifecycleadvance/` returned ok in 4.9s. A scratch repo with a bare foreign-writer claim confirmed the advance next_argv, the text-renderer output, the triage phase after following it, and the unchanged file bytes.
**Suggested testing:** 2 items. (1) Run a real `phandoff` with one merged claim and one queued REQ, then paste the block into a fresh session and check that the merged claim reaches finalization before the run fans out. (2) Run a claim with an unchecked Open Question through `--take-over` once to confirm the `pending-answers` wording gap in M2.
**Follow-ups created:** None (6 findings report only)

*Reviewed by review-work action*

**Post-review fixes (orchestrator, D-05, merge 2cbd6858):** M3, M4 and M5 fixed in prose on the builder branch — the handoff no longer promises one pasted command or a one-line resume, its Step 1 test and checklist say "for queued work" and name the claim exception, claim lines go above the resume command, and the guide's continue-vs-run bullet points to Context limits. M2's prose wording now says "returns it to the queue"; the Go stop reason keeps "requeues it as pending" (report only). M1 is answered by the new wording ("does not continue a claim's remaining phases"). Verified by re-reading `git diff 730c22cf..2cbd6858` (two files, 8 lines each way) and a green repository gate at 64d4f665.


## Decisions

Builder decisions, from the hand-back (`do-work/runs/work-2026-10-02-204757/REQ-629-handback.md`). All DECIDE & STATE.

- D-02: The paste block's sufficiency sentence becomes plural ("These commands are sufficient; everything below them is context."), because with claims the block carries more than one command and the singular sentence was not true. No test pins it.
- D-03: The new test's fixture is a merged claim (`commit:` set) with a same-host writer label, modelling the UR-132 case; the finding is writer-agnostic, so the label is realism only.

Orchestrator decision:
- D-04: Scope extension. The builder reported that `skills/do-work/docs/work-guide.md` line 165 told a new session to start with `do-work run` and look for "takeover authority", and line 130 said plain `recover` "returns typed takeover options". Line 165 is the same trap this REQ removes, so leaving it would fail the REQ's second requirement for any reader of the guide. The orchestrator fixed both sentences on the builder branch (two lines) and added the file to Scope and `write_set` instead of queuing a follow-up. DECIDE & STATE.
- D-05: Review fixes before release. The review passed (95%) but found that this change left contradictions in the same two files: the handoff action still promised "one pasted command" and a "one-line resume", its Step 1 test and checklist said `do-work run` alone must do the right thing, the paste block did not say whether the claim lines run before or after the resume command (M3, M5), and the guide's continue-vs-run bullet said the checkpoint handles resume (M4). The orchestrator fixed these on the builder branch (claim lines go above the resume command; the reset is described as "returns it to the queue", answering M2's wording in prose) instead of shipping a release that contradicts itself. M1, N1 and the Go stop-reason wording in M2 stay report only. DECIDE & STATE.

## Discovered Tasks

From the builder's hand-back and the review, impact-stamped per review-work Step 10; none is impact-critical, so nothing was queued.

- Builder: `docs/work-guide.md` lines 130 and 165 steered a resuming session to takeover. — folded into this REQ by D-04
- M1: "the resume command never continues a claim" was too broad given the heavy-lane drain. — reworded by D-05
- M2: the Go stop reason says the reset "requeues it as pending", but the reset can also produce `pending-answers` or keep `blocked`, and it also drops `route` and `write_set`. — impact-negligible → report only
- M3: one-command framing in the handoff action. — fixed by D-05
- M4: the guide's continue-vs-run bullet said the checkpoint handles resume. — fixed by D-05
- M5: claim lines vs resume command order unstated. — fixed by D-05
- N1: the finding code `RECOVERY-TAKEOVER-AVAILABLE` and its decision text still lead with "takeover" while the next step is continue; kept for compatibility. — impact-negligible → report only

## Lessons Learned

**What worked:** Reproducing the failure in a scratch repo before choosing a fix showed that bare `advance REQ-NNN` already continues a claim, so the fix was one next_argv and its wording rather than a new resume mode.
**What didn't:** The UR-132 handoff told the session to run `advance REQ --request-path P`, which is refused at judgment phases; only the bare form classifies. And the first fix left the handoff action's own "one pasted command" framing contradicting the new claim lines, which only the review caught.
**Worth knowing:** A finding's `next_argv` is an instruction that sessions follow literally; never put a destructive command there, even behind a warning. Queue-mode `advance` ignores working claims, so a resumed session must name each claim. The Scope checker reads every backticked word in a Scope list item as a path; keep commands out of backticks there.

## Orientation

Now a session resuming a handoff continues its claimed REQs with `advance REQ-NNN`, and `recover` suggests that read-only step instead of the resetting takeover; lives in the do-work-cli recovery command and the handoff action (`_dev/primes/prime-action-files.md`). Not a map change: one finding's next step and the handoff's paste-block rules. Prime spot-check: `prime-action-files.md` paths still exist and it does not restate recover's next step.

## Heavy Verification Plan

- Base revision: 89145e03f06cb724cb2ec83b6f1357c0ca972997
- Target revision: 2cbd685813412cab0b51a783a0fbd6b25a4859a8 (landed in `commit:`); planned per merge: 89145e03..730c22cf and 8f9f3099..2cbd6858, lanes unioned
- do-work-cli-integrations — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — do-work-cli source changed
- staged-skills — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — shipped files under skills/ changed
- updater — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — do-work-cli source changed
- installer — `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — do-work-cli source changed

## Heavy Verification Result

- Target revision: 2cbd685813412cab0b51a783a0fbd6b25a4859a8
- Execution revision: 069f922e64c1ed14ca31211642e395a784561f57 (detached drain checkout `.git/work-run-2026-10-02-204757/drain-head`, one run for REQ-628 and REQ-629 at the integration tip containing both merges)
- staged-skills: exit 0, executed, 33s
- do-work-cli-integrations: exit 0, executed, 67s
- updater: exit 0, executed, 67s
- installer: exit 0, executed, 29s

Green: every selected lane present, exit 0, none skipped, none reused. An earlier run at 8f9f3099 (before REQ-629's D-05 prose fixes) was also green on all four lanes.
