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
