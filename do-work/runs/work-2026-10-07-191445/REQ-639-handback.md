# REQ-639 hand-back (Delegated integration: the coordinator shape)

- Branch: `worktree-agent-REQ-639-delegated-integration-coordinator-shape` (no name collision)
- Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-639-delegated-integration-coordinator-shape`
- Base: `0efba253`. Commit: `a3ddde35` ([REQ-639] delegated integration: the coordinator shape for fan-out runs)

## File manifest

- `skills/do-work/actions/work.md` (modified): Step 6 gains the **Landed hand-back — consume it, never re-dispatch.** condition before the builder spawn. The dispatch paragraph writes the held instant into the manifest row. Step 10 gains the integrator-never-checkpoints sentences. Orchestrator Checklist lines for Step 6 and Step 10 name the new conditions.
- `skills/do-work/actions/work-reference.md` (modified): Fan-Out Dispatch gains the **Delegated integration — the coordinator shape.** paragraph after **Serial-only**. The guardrail table gains the `per-integrator input` row, and the `manifest.md` row names the held dispatch instant. Hand-back sequence step 0 stages `REQ-NNN-integrate.md` with the other owner-written run artifacts.
- `skills/do-work/docs/work-guide.md` (modified): two sentences in "Building several REQs at once".

No tests touched. No Go, no flag, no heading changed, no new file.

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]:** Read the brief, the REQ, prime-action-files.md, the six alternate-writer-contract-drift bullets, and the general, coding-guardrails, shared-principles and communication-style crew rules. Approach: five prose insertions at the REQ's named anchors, 4(b) rewritten per the coordinator note, then a restatement sweep across work.md, work-reference.md, work-guide.md and background-agents.md, changing a restatement only where it would otherwise be wrong. Checked recovery_commands.go before wording 4(b), which showed `--take-over REQ-NNN` resets only that one claim.
- [x] **[APPLY]:** Edits stayed inside the three declared files. Requirements 1-5 applied in place. Sweep-driven changes inside scope: the `manifest.md` table row, the two checklist lines, and the hand-back merge step 0 stage set (D-04). Committed as `a3ddde35`.
- [x] **[UNIFY]:** `git diff --stat 0efba253..a3ddde35`: work-reference.md 7 (+5 -2), work.md 10 (+6 -4), work-guide.md 2 (+1 -1); 3 files, 12 insertions, 7 deletions. `git diff --check` clean. `shipped-package-reference-contract.sh` PASS 1s, `contract-regressions.sh` PASS 21s (first run 23s, before D-04). Files checked by cold read: work.md Step 6 landed paragraph, dispatch paragraph, Step 10, checklist; work-reference.md new paragraph, guardrail table, step 0; work-guide.md paragraph. Every citation names a real heading or bold label. No debug artifacts.

## Decisions

- **D-01 (DECIDE & STATE): 4(b) entry is `advance REQ-NNN`, not `recover`.** Followed the coordinator note. The paragraph says the integrator does not run `recover`, that the coordinator ran it at the queue boundary, and that a mid-wave `recover` reads the wave's own in-flight artifacts (for example the untracked run directory) as dirt it cannot attribute. The integrator enters with the read-only `advance REQ-NNN`. The never-rwr and never-take-over sentence is kept. Reversible prose, evidence from this run.
- **D-02 (DECIDE & STATE): `--take-over` is worded by what it actually resets.** The REQ says all three verbs "reset every sibling claim". `recovery_commands.go` filters `--take-over REQ-NNN` to that one claim. The text now says run-with-recovery and `--assume-sole-authority` reset every sibling claim and requeue the wave, and `--take-over` resets the integrator's own claim and strips the sections the coordinator wrote. The prohibition is unchanged.
- **D-03 (DECIDE & STATE): checklist lines for Step 6 and Step 10 changed.** Left alone, the Step 6 line says "spawn agent" for a landed hand-back, and the Step 10 line would lead an integrator into `advance --checkpoint`. Both now name the condition in a short clause.
- **D-04 (DECIDE & STATE): hand-back merge step 0 stages `REQ-NNN-integrate.md`.** Step 0 enumerates the owner-written run artifacts to stage. Without the new brief there, an integrator classifies it as "leave alone and name", so it is never committed and stays untracked dirt that a later `recover` cannot attribute. This is inside the declared file. No Go code enumerates run-artifact names.
- **D-05 (DECIDE & STATE): Rules bullet "the orchestrator is the sole integrator" unchanged.** The integrator plays the orchestrator role for its span, and work.md already says builder and orchestrator are roles, not processes. Also unchanged for the same reason: the builder bullet at Step 6, **Sole integrator**, and the "one integrator applies every seam" sentence, because integrators run one at a time.
- **D-06 (DECIDE & STATE): work-guide **What isn't specified** and the Execution Model sentence unchanged.** The new guide sentence keeps the one-writer rule visible ("writes nothing in the working tree while that agent runs"), so a coordinator that only waits is not the two-writers case. Rewording the unspecified case would widen this REQ.
- **D-07 (DECIDE & STATE): the manifest-row instant is conditioned on "When the run has a manifest".** The serial loop does not always have a run directory, so the sentence states the condition instead of assuming one.

## Discovered Tasks

- impact-negligible: The landed-hand-back condition does not say what to do when the manifest row lacks the dispatch instant, for example a run manifest written before this release. An integrator would then have no legal instant for record-timing-event. → report only
- impact-negligible: `crew-members/background-agents.md` step 3 lists the manifest's columns (agent, slice, output filename, landed status) without the held dispatch instant. It is the generic pattern, so it was left as is and is outside the write boundary. → report only

## Lessons read

- `_dev/primes/lessons-action-files.md` lines 52-60: the six `[family: alternate-writer-contract-drift]` bullets (REQ-477, 498, 513, 461, 531, 566), family-targeted. The prime's Traps line for the same family.
- No `required_lessons` stamped; none missing.

## Proposed lesson entry (for _dev/primes/lessons-action-files.md)

- [family: alternate-writer-contract-drift] [REQ-639: adding a new run artifact (the integrator brief) is a writer change too; the hand-back merge's stage-set enumeration was a reader that would have left it untracked dirt, found only by sweeping restatements past the REQ's named anchors. A captured claim about a recovery verb's reach was also wrong (`--take-over` resets one claim), so check the command before restating it.]

## Proposed CHANGELOG entry

**Fan-out runs can hand each REQ's integration to one agent at a time.** In `do-work run --fan-out`, the session that dispatches builders may stay a coordinator and pass each REQ's merge-to-commit span to one integrator agent at a time. Step 6 now consumes a hand-back that has already landed instead of building the REQ again, and the run manifest carries the dispatch time so an integrator can record builder timing. The integrator enters with `advance REQ-NNN`, never runs recovery resets or the session checkpoint, and the coordinator writes nothing while it runs. This keeps the user's conversation free during long integrations. The time saved is still in the build phase only.

## Integration seams

None.

## Tests

| Command (from worktree root) | Result | Wall time |
| --- | --- | --- |
| `bash _dev/tests/shipped-package-reference-contract.sh` | PASS | 1s |
| `bash _dev/tests/contract-regressions.sh` | PASS | 21s (23s on the run before D-04) |
| `git diff --check` | clean | |

## Restatement grep accounting

Run from `skills/do-work`, on the four files `actions/work.md actions/work-reference.md docs/work-guide.md crew-members/background-agents.md`:

```
grep -nE '<pattern>' actions/work.md actions/work-reference.md docs/work-guide.md crew-members/background-agents.md
```

- `landed|re-dispatch|re-run them`: work.md 278 and 489 are new. background-agents 52 and 131 agree (landed status column, findings file means done). work-reference 480 changed. Other hits (heavy-hold "landed" at 193/271/429/441/443/788/810, guide 103, work.md 511) are a different sense of "landed implementation"; unchanged.
- `dispatch instant|dispatch_at`: work.md 278 and 286 changed. work-reference 157 (schema example) and 263 (stamp read contract) describe the field itself; unchanged. 480 changed.
- `advance --checkpoint|session-end writer`: work.md 468 and 498 changed. work-guide 76, 90, 165 and work-reference 522 describe the command for a normal session; still true, unchanged.
- `integrator|integration stays serial|Integration is serial`: work.md 37 and work-reference 457 agree with series integrators. work.md 299 and 540, work-reference 409/424/434/456 unchanged (D-05). background-agents 181 agrees. New text at work.md 278/286/468/498 and work-reference 470/479.
- `run-with-recovery|assume-sole-authority|take-over`: work-guide 130 and 165, work-reference 23, 297, 304, 310, 516 describe the verbs for a whole session; unchanged. New prohibition at work-reference 470.
- `two sessions|same working tree|one working tree`: work-reference 17 and work-guide 138 unchanged (D-06). work-reference 458 is about builders sharing a tree; unchanged.
- `per-builder|per-integrator|manifest.md`: work-reference 430 changed (D-04), 477/478 unchanged, 479 new, 480 changed. background-agents 51 and 192 generic; noted as a discovered task.
- Extra: `grep -rn 'brief\.md\|handback\.md\|integrate\.md' skills --include='*.go'` found no Go enumeration of run artifacts.
