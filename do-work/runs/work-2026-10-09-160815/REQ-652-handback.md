# Hand-back — REQ-652 (routine operator work is never a REQ; a tracked operator configuration is blocked, dependency-gated, and completed from clarify)

- Branch: `worktree-agent-REQ-652-operator-work-routine-not-captured-tracked-config-blocked-dependency-gated`
- Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-652-operator-work-routine-not-captured-tracked-config-blocked-dependency-gated`
- Base: `2eb14357`
- Commits: `8f6f6e89` (single commit)

## File Manifest

- `skills/do-work/actions/capture.md` (modified): Earmark assessment's "at the keyboard" clause now points at the operator rule instead of flatly capturing `blocked`. External-condition assessment gains two sentences: a routine operator act is never a REQ (user does it on their own time, capture summary names the step), and a tracked special operator configuration is one `blocked` REQ with `depends_on` on the AI-buildable REQs, completed by clarify.
- `skills/do-work/actions/clarify.md` (modified): Step 5.5 plain-blocked prompt gains `4. Done — I did it myself` and a **Done** bullet (append `## Operator receipts` under Outside-text containment, set `completed` + `status_changed_at` in place, no `unblock`, precedent cited). Step 6 report lists REQs completed this way.
- `skills/do-work/actions/work.md` (modified): the "park on the operator" Error Handling row gains one sentence: a routine operator act as the remaining work completes the REQ, no blocked flip, no failure, no abandon.
- `skills/do-work/actions/work-reference.md` (modified): Failure Classification Environment row gains the inverse sentence; the blocked flip stays for a precondition to AI work and a tracked special configuration.
- `skills/do-work/docs/work-guide.md` (modified): the earmark paragraph's operator sentence now states the three shapes and the board path Pending → Waiting, Needs input · Blocked, Done.

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]:** Read brief, REQ, UR-143, crew rules, both primes and every `alternate-writer-contract-drift` lesson. Edit only the five named files, delete before adding: rewrite the existing "at the keyboard" clause in capture and the operator sentence in the guide rather than add beside them. Key the routine-act rule on its condition, with example verbs marked as examples. Swept `skills/` for other restatements of the operator-blocked rule (lesson REQ-566/REQ-648); found the `assigned_to` schema line in work-reference and the board action's Needs-input paragraph, both still accurate (see D-04).
- [x] **[APPLY]:** Five prose edits as planned, one commit. No Go, field, status, flag, transaction, probe or board change. All existing bold labels and headings unchanged.
- [x] **[UNIFY]:** `git diff --stat HEAD~1`: 5 files changed, 8 insertions(+), 6 deletions(-). `git diff --check` clean. `bash _dev/tests/shipped-package-reference-contract.sh` exit 0 (2 s). `bash _dev/tests/contract-regressions.sh` exit 0 (30 s). Each file read back: capture (both assessments, cross-ref to clarify Step 5.5), clarify (prompt block, Done bullet with cross-ref `actions/cleanup.md` → **Pass 0: Sweep Finished Queue Items**, Step 6), work.md row 523, work-reference row 788, work-guide line 119. No debug artifacts.

## Red-Green Evidence

- **RED (before):** capture.md:107 routed "at the keyboard" straight to a `blocked` capture and :108 had no rule for an operator act that is the whole remaining work, so a "deploy and verify on the live host" request minted a `blocked` REQ with no `depends_on`. clarify.md Step 5.5 offered only unblock, leave, abandon, so a finished operator REQ went back to `pending`. work.md:523 and work-reference.md:788 flipped it to `blocked` again.
- **GREEN (after):** capture says the routine act is never captured and names the step in the summary; a tracked configuration is `blocked` + `depends_on` and shows Pending → Waiting, then Needs input · Blocked. Clarify offers Done, which writes `## Operator receipts` and flips to `completed` (Done column, cleanup Pass 0 archives). Both run rows complete the REQ instead of flipping.

```
$ grep -n "own time" skills/do-work/actions/capture.md
108:- **External-condition assessment** — ... is never captured as a REQ: the user does it on their own time, ...
$ grep -n "Done — I did it myself" skills/do-work/actions/clarify.md
184:  4. Done — I did it myself
$ grep -n "routine operator" skills/do-work/actions/work.md skills/do-work/actions/work-reference.md
skills/do-work/actions/work.md:523:| Orchestrator wants to park a queued REQ on the operator ...
skills/do-work/actions/work-reference.md:788:| **Environment** | ...
```

## Decisions

- **D-01 (DECIDE & STATE, recorded challenge):** The upstream proposal's "no REQ, report-only line, recommend abandon" shape was rejected by the maintainer in favour of the blocked + `depends_on` placement the board already has (model.go Waiting / Needs input · Blocked) and a clarify completion path into Done. Built the maintainer's shape only.
- **D-02 (DECIDE & STATE):** A request whose only content is a routine operator act writes no UR and no REQ ("none means no UR and no REQ"). The REQ's Verified Facts say `capture.md:16` already allows a zero-REQ UR, but that exception covers only a capture whose every request was folded; widening it would change the core contract the REQ says stays untouched. Value: no orphan UR, no change to verify-requests linkage. Risk: the verbatim input of an operator-only request is not stored; reversible by one clause if the maintainer wants the UR kept.
- **D-03 (DECIDE & STATE):** The clarify precedent is cited as "the confirmed `builder_decided: true` path ... in Step 5" (the operative bullet is at line 135; line 267 is its checklist line). No Verification Checklist line added: the existing checklist has no blocked-prompt items to sit beside, so a line would be ceremony.
- **D-04 (DECIDE & STATE):** Left unchanged: the `assigned_to` schema line in work-reference ("work that waits on the user as operator is `status: blocked`") and the board action's Needs input paragraph. Both stay true under the new rule (a routine act no longer waits as a REQ; a tracked configuration still is `blocked`), and both are outside the write set.
- **D-05 (DECIDE & STATE):** work-guide keeps a short sentence for the "at the keyboard" precondition case so the rewrite did not drop it.

## Discovered Tasks

- impact-negligible: REQ-652 Verified Facts misstate `capture.md:16` as allowing a zero-REQ UR for non-fold captures; only matters if a later REQ relies on that claim → report only

## Lessons Read

- `_dev/primes/lessons-action-files.md`, every `[family: alternate-writer-contract-drift]` bullet (REQ-477, 498, 513, 461, 531, 566, 640, 641, 642, 644, 647, 648). Applied: swept restatements beyond the cited lines (D-04), copied the board label "Needs input · Blocked" from the existing source spelling, checked the cross-referenced heading exists (`Pass 0: Sweep Finished Queue Items`).

## Proposed Lesson Bullet (integrator writes)

```
- [family: alternate-writer-contract-drift] [REQ-652: a REQ's Verified Facts can cite an existing exception for a shape it does not cover (capture's empty-UR rule is fold-only, not "nothing AI-buildable"); read the cited line's own condition before building on it, and decide the uncovered case in one clause rather than widening the exception](../../do-work/archive/UR-143/REQ-652-operator-work-routine-not-captured-tracked-config-blocked-dependency-gated.md#lessons-learned)
```

## Proposed CHANGELOG Entry (integrator writes)

```
## Routine Operator Work Stays Out of the Queue; Tracked Operator Setup Completes From Clarify

Deploys and publishes you do yourself no longer become REQs that block, re-block and end in abandon. A special operator configuration you ask to track now waits behind its code, shows up when it is ready, and moves to Done when you say you did it.

- Capture never writes a REQ for a routine operator act after the code ships; it captures only the AI-buildable parts and names the operator step in its summary.
- A tracked special operator configuration is captured `blocked` with `depends_on` on the AI-buildable REQs: Pending → Waiting until they finish, then Needs input · Blocked.
- `do-work clarify` adds "Done — I did it myself" for plain blocked REQs: it records your words under `## Operator receipts` and marks the REQ `completed`.
- A run that finds only a routine operator act left completes the REQ and names the step, instead of flipping it to `blocked`.
```

## Integration Seams

None. REQ-653 depends on this REQ and edits the same two run files (work.md row 523, work-reference.md row 788).

## Test Wall Times

| Test | Exit | Wall |
|---|---|---|
| shipped-package-reference-contract.sh | 0 | 2 s |
| contract-regressions.sh | 0 | 30 s |
