---
id: REQ-528
title: '[impact-critical] Anchor resolved-question marker matching to its position'
status: completed
created_at: 2026-09-03T03:10:00Z
user_request: UR-081
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-critical
effort_estimate: effort-substantive
related: [REQ-460, REQ-413]
sweep: true
sweep_key: answer-line-marker-position-spoofing
review_generated: true
write_set:
  - skills/do-work/tools/do-work-cli/internal/publication/answer.go
  - skills/do-work/tools/do-work-cli/internal/publication/answer_test.go
claimed_at: 2026-09-03T09:10:13Z
route: B
estimate:
  p50_active_minutes: 25
  confidence: medium
  calculated_at: 2026-09-03T09:11:00Z
  basis:
    - Route B
    - 2-file write set
    - 1 subsystem involved
    - 4 acceptance criteria
    - persistence or schema changes
    - cross-route regression gates
completed_at: 2026-09-03T10:15:00Z
commit: f1197c6
release_at: 2026-09-03T10:15:00Z
kb_status: promoted
kb_entry: REQ-528-anchor-resolved-question-marker-matching.md
---

# Anchor Resolved-Question Marker Matching to Its Position

## What

`allResolvedQuestionsMatch` decides whether every resolved question on a REQ carries the same disposition marker, and that verdict drives a terminal status write. It tests `bytes.Contains(line, marker)` against the whole `- [x] ` line, so **answer text** that happens to contain the marker counts as the marker. A plain one-line answer summary containing `→ Discarded:` therefore makes an *answered* question read as *discarded*, and the REQ is silently cancelled and archived.

## Instances

- `skills/do-work/tools/do-work-cli/internal/publication/answer.go:410-421` — `allResolvedQuestionsMatch` matches the marker anywhere in the line rather than at the position the answer writer put it.
- Demonstrated: answer Q1 with the summary `keep it → Discarded: not really`, then discard Q2. The REQ lands `status: cancelled` with `completed_at` set and the terminal archive path selected, despite Q1 having been answered.
- The `→ Confirmed:` variant, with `builder_decided: true`, reaches `status: completed` by the same route.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `prime-do-work-cli.md` and `lessons-do-work-cli.md` in full (including REQ-460's `closed-enumeration-for-a-condition` entry), `clarify.md` Steps 4-5, and the crew members. Approach: share the writer's separator and disposition-prefix literals between writer and reader, read the disposition only at the position that separator marks, and treat a line where the separator is not unique as having no identifiable disposition.
- [x] **[APPLY]:** Two files, both inside the declared Scope.
- [x] **[UNIFY]:** Orchestrator independently re-ran `go build ./...`, `go vet ./...` (clean), `gofmt -l .` (silent) and `go test -count=1 ./internal/publication/` (`ok … 18.7s`), and read the full `answer.go` diff. No debug artifacts; the writer's three literals are now shared constants consumed by both the writer and the reader, which is what makes the two unable to drift.

## Finding Provenance

- **Finding F1** — `impact-critical` — from REQ-460's independent review (Approve, 89%), reproduced in a scratch fixture rather than reasoned about.
- **Pre-existing, not a regression.** The old ten-prefix predicate inlined the same text, so REQ-460 neither introduced nor widened this. It surfaced because REQ-460 redefined exactly this contract: `skills/do-work/actions/clarify.md:103` names `- [x]` as one of "this file's own delimiters" precisely so that answer text cannot be read as one, and REQ-460's new doc comment claims completeness over three *Markdown* ingredients while the real contract is broader than Markdown.

## Detailed Requirements

- The disposition marker must be recognized only at the position the answer writer places it, never anywhere in the line. Answer text containing the marker's characters must not satisfy it.
- A resolved question whose answer text contains `→ Discarded:` or `→ Confirmed:` must not contribute to a terminal-status verdict as though it carried that disposition.
- No terminal status (`cancelled`, `completed`) may be written on evidence a user's own answer text can forge.
- Keep the existing verdict for genuinely uniform dispositions; this is a matching fix, not a policy change.

## Constraints

- Do not solve this in `summaryRequiresContainment`. The reviewer was explicit that the predicate is the wrong layer: containment classifies a summary's shape, while this is about where a marker is anchored on a line the writer itself composed.
- Preserve the existing refusal codes and typed results.

## Dependencies

No request prerequisite. Independent of REQ-460, which is already archived.

## Red-Green Proof

**RED prompt/case:** Answer one open question with the one-line summary `keep it → Discarded: not really`, then discard a second open question, and inspect the resulting `status`.
**Why RED now:** `bytes.Contains` on the whole line lets the answer's own text supply the marker, so the answered question is counted as discarded and the REQ is cancelled.
**GREEN when:** That REQ stays non-terminal with Q1 recorded as answered; a genuinely all-discarded REQ still reaches `cancelled`; and the `→ Confirmed:` plus `builder_decided: true` route still reaches `completed` only when every resolved question really carries it.

---
*Source: REQ-460 independent review finding F1.*

---

## Triage

**Route: B** - Medium

**Reasoning:** The defect and its single function are named exactly, but the fix is not local to that function: the marker's *position* is determined by a writer forty lines away, and how a reader can recover that position without re-parsing user text had to be worked out before scoping.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

The writer composes each resolved line at `answer.go:186-196`:

```
- [x] <question text> → [Confirmed: |Discarded: ]<summary>
```

The disposition prefix therefore sits at exactly one place: immediately after the ` → ` this writer appended, which is immediately after the question text. `allResolvedQuestionsMatch` (`:410`) throws that structure away and asks `bytes.Contains(line, "→ Discarded:")` against the whole line — so anything at all in the line, including the user's own `<summary>`, can satisfy it.

Both callers (`:220`, `:223`) turn the verdict straight into a terminal status: `cancelled`, or `completed` when `builder_decided` is true. A forged marker therefore does not merely mislabel a line, it writes a terminal status and selects the archive path.

Two things make the naive repairs unsound, and the builder needs to know both before choosing:

1. **Anchoring to the last ` → ` on the line is not sufficient.** The summary is user text and may itself contain ` → `. REQ-460's containment predicate does not help here: it decides on the *first* byte, and `→` is U+2192 — not ASCII punctuation, so a summary of `a → b` is inlined today.
2. **The function must also judge lines this invocation did not write.** Questions resolved in earlier clarify rounds are already in the body, so a purely writer-side record of what was just emitted cannot answer for them.

So the fix has to recover the boundary between question text and disposition from structure the writer controls, for old and new lines alike, rather than from a substring search over text the user supplies.

*Generated in-session (single-pass discovery)*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/publication/answer.go` (modify) — make the disposition recognizable only at the position the writer places it
- `skills/do-work/tools/do-work-cli/internal/publication/answer_test.go` (modify) — the spoof as a failing case, plus the genuine all-discarded and all-confirmed routes

**Files I will NOT touch:** `summaryRequiresContainment` and `isMarkdownBlockPunctuation` (REQ-460's predicate — the reviewer was explicit that containment is the wrong layer for this), and the C0/DEL and multiline validators.

**Acceptance criteria (restated from REQ):**
- [x] The disposition marker is recognized only at the position the answer writer places it; answer text containing the marker's characters does not satisfy it
- [x] A resolved question whose answer text contains `→ Discarded:` or `→ Confirmed:` does not contribute to a terminal-status verdict as though it carried that disposition
- [x] No terminal status is written on evidence a user's own answer text can forge — **not delivered by the first implementation** (review finding F1), delivered by the remediation's write-side refusal (D-06). Ticked in error before review; the correction is recorded rather than quietly overwritten.
- [x] Genuinely uniform dispositions still reach `cancelled` and `completed` exactly as before

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/publication/answer.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/publication/answer_test.go` (modified)

**What was done:** The separator and the two disposition prefixes are now shared constants (`resolvedQuestionSeparator`, `confirmedLabelPrefix`, `discardedLabelPrefix`), consumed by the writer at `:197-206` and by the reader at `:444`, so the composed line and the line the terminal verdict reads cannot drift apart. `resolvedQuestionDisposition` returns the text at the one position the writer can place a disposition, plus whether that position is identifiable at all: both fields around the separator are human text, so a line carrying more than one separator has no attributable disposition. `allResolvedQuestionsMatch` now prefix-tests the disposition at that position instead of searching the whole line, and an unidentifiable line fails the check — which routes to the non-terminal status rather than cancelling or completing on evidence that cannot be attributed to the writer.

## Decisions

- **D-01** — The disposition is read only immediately after the single separator the writer appends; a resolved line carrying zero or more than one separator has no identifiable disposition and fails the check. **ESCALATE.** **Value:** the separator is the one byte sequence on that line the writer itself contributes, recoverable from prior-round and hand-edited lines alike — no new schema field, no format change, no new refusal code, and no second consumer to sweep (a repo-wide grep for `Discarded:` hits only `answer.go`). **Risk:** a genuinely uniform section in which a question text or a discard summary itself contains ` → ` now lands `pending` instead of `cancelled`/`completed`. That is the required failure direction — non-terminal, visible, hand-correctable — against archiving a REQ whose questions were answered. Locked in as a test row rather than left implicit.
- **D-02** — Rejected anchoring to the last ` → ` on the line, confirming the REQ's Exploration by measurement rather than by assertion. **DECIDE & STATE.** `summaryRequiresContainment` decides on `summary[0]`, and `→` is U+2192 whose lead byte `0xE2` falls outside all four ASCII-punctuation ranges, so `a → b` is inlined today — the reproduction shows the spoof summary landing inline with no `ANSWER-RAW-PAYLOAD-REQUIRED`. Last-arrow anchoring *is* the spoof: in `keep it → Discarded: not really` the final arrow is the summary's own.
- **D-03** — Rejected recording what this invocation wrote, and sharpened why. **DECIDE & STATE.** The REQ's Exploration said the function "must also judge lines this invocation did not write". Stronger: that is *all* it does. `allSubmittedDiscarded`/`allSubmittedConfirmed` (`:144-145`) already covers the submitted set, so a mixed batch can never reach a terminal status within one invocation — the spoof is only reachable across rounds. A writer-side record would answer for precisely the lines that never needed judging.
- **D-04** — Kept the `- [x] ` prefix scan and the `found` semantics; declined a structured frontmatter field. **DECIDE & STATE.** The cheaper fix holds, so the REQ's own constraint rules the field out — and frontmatter is hand-editable too, so it would buy less than it appears to.
- **D-05** — Did not touch `summaryRequiresContainment`, and added no write-side refusal for arrow-bearing summaries. **DECIDE & STATE.** The reader is now robust to them, so a refusal would be defense nothing earned.

## Discovered Tasks

- **Same shape, second terminal status, same file:** `answer.go:268` gates the stakeholder branch's `status: completed` + `completed_at` + archive move on `bytes.Contains(ToLower(blockedHistory), "resolved")` and `bytes.Contains(ToLower(implementation), "no code")`. Both payloads are caller-authored prose and both tokens are matched anywhere, so `still not resolved` and `no code review yet` each satisfy the gate.
- **Same shape, non-terminal:** `answer.go:291` gates the `blocked_by` lifecycle write on `bytes.Contains(reportsHistory, reportPath)` over caller-authored history text, unanchored to any line or position.
- **Adjacent, destructive:** `internal/cleanup/cleanup_plan.go:236` deletes a `do-work/CHECKPOINT.md` line that merely contains `- <REQ-ID>:` plus the writer token, with no position anchor. Narrower (both tokens must share a line) but the same unanchored-line-selection shape driving a destructive edit.
- **Adjacent, bounded:** `internal/finalization/finalization_discovery.go:593` admits an unjournaled changelog tail into replay on `bytes.Contains(inserted, requestID)`; mitigated by `singleInsertion` bounding the searched bytes to one verified diff hunk.
- **Prose disagrees with the writer:** `skills/do-work/actions/clarify.md:124` documents the discard form as `- [x] [question] → Discarded` — no colon, no summary — while the CLI writes `Discarded: <summary>` and now prefix-tests `Discarded: `. The CLI is the authority, so the prose line is the one to correct.
- `internal/suiteinstall` → `TestBuiltInstallAndUpdateExit130WhenSignalsInterruptBlockedConfirmation/update-suite/TERM` failed once in three full-suite runs and passed 4/4 in isolation. Third sighting of the flake tracked as REQ-525, now on a third subtest.

## Testing

**Tests run:** `go build ./... && go vet ./... && gofmt -l .`; `go test -count=1 ./internal/publication/`; `go test ./...`; the built CLI's `answer --manifest` against a throwaway git repo across two clarify rounds; canonical repository gate `bash _dev/tests/maintainer-verify.sh`.
**Result:** ✓ `internal/publication` green (`ok … 18.4s`). Gate exits 1 on the recorded baseline failure only — `internal/toolboxcommands` → `TestRemediationCancellationReachesMediaGitCommitAndRollback`, tracked as REQ-524, in a package this diff does not touch.

**Red-green validation:** traces the REQ's Red-Green Proof, which asked for the two-round forgery and the preserved genuine routes.
- `TestBuildAnswerPlanRefusesDispositionForgedByAnswerText` (both labels, new in remediation): ✗ `forgeable answered summary accepted: refusal=(*publication.Refusal)(nil)` → ✓
- `…RefusesDispositionForgedByAnswerText` first form (position anchoring, first implementation): ✗ the round-2 document published `status: cancelled` with `completed_at` set → ✓ `pending`
- `…JudgesPriorRoundResolvedLinesAtTheWriterPosition/summary merely mentioning the marker` and `/question text carrying the separator`: ✗ → ✓
- F2 refusal subtests: ✗ `declared terminal archive path silently ignored: refusal=(*publication.Refusal)(nil)` → ✓ `ANSWER-ARCHIVE-PATH-MISMATCH` naming the blocking line
- Genuine all-discarded → `cancelled` and genuine all-confirmed + `builder_decided` → `completed` are green throughout; they are the regression the fix could plausibly have broken, and the reviewer confirmed both redden under a blanket-refusal mutation, so they are controls rather than padding.

**CLI evidence (built binary, throwaway repo, two rounds — not a test harness).** Before, at implementation `1d8c82c`: the forged round succeeded and round 2 landed `status: cancelled`, `completed_at: 2026-09-03T00:00:00Z`, path `do-work/archive/REQ-1-attack.md`; the `Confirmed: ` variant landed `completed` + archive. After: round 1 is refused with `ANSWER-SUMMARY-INVALID` and the reason `an answered summary must not open with the disposition label "Discarded: " …`, the question line is byte-unchanged, and round 2 lands `pending-answers` in `do-work/queue/`. The forged round writes nothing, so Q1 stays open.

**F2 through the same CLI:** before, a genuine discard whose summary carried the house `A → B` cross-reference published `status: pending` with the caller's declared `archive_path` silently dropped. After, it refuses with `ANSWER-ARCHIVE-PATH-MISMATCH` naming the derived status and quoting the blocking line.

**Neuter-and-confirm, six guards, one at a time, full package run each time:** the write-side label refusal reddens only the F1 test; the archive-path check moved back inside `if terminal` reddens only the two F2 refusal subtests and leaves the acceptance control green; dropping the evidence reddens only the two evidence assertions; removing the reader's `\r` trim reddens only the CRLF F2 subtest; removing `findQuestionLine`'s `\r` trim reddens the CRLF table row with `ANSWER-QUESTION-NOT-UNIQUE: question identity matched 0 lines`; and skipping unattributable lines instead of failing the check reddens the new zero-separator row plus both F1 controls and both ambiguous rows.

**Corpus measurement (from the review, not re-derived):** 124 resolved question lines across 110 files under `do-work/` — 93 carry exactly one separator, 29 carry none, **2 carry more than one**, both from this project's own `A → B` cross-reference style inside answer text. Exactly **1** historical verdict would have changed, and it is already archived; **0** affected files in `do-work/queue/` or `do-work/working/`, so nothing in flight changes.

*Verified by work action*

## Review

**Overall: 60%** | 2026-09-03T10:00:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 63% |
| Code Quality | 82% |
| Test Adequacy | 78% |
| Scope | 100% |
| Risk | Critical |
| Acceptance | Partial |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- F1 — the spoof survived at the writer's own position: an `answered` summary that *begins* `Discarded: ` is byte-identical to a genuine disposition, so position anchoring could not separate them and the two-round forgery still reached `cancelled`/`completed` plus the archive move. Requirement 3 and the REQ's own GREEN clause were both unmet, and D-05's premise ("the reader is now robust to them, so a refusal would be defense nothing earned") was measurably wrong — `impact-critical` → fixed in remediation (D-06)
- F2 — the accepted non-terminal fallback was silent, and the caller's declared `archive_path` was itself silently ignored because it was validated only inside `if terminal`; a user who discarded every question saw the REQ return to the queue with nothing said — `impact-user-visible` → fixed in remediation (D-09, D-10)
- F3 — `clarify.md:124` prescribes `- [x] [question] → Discarded` with no colon while the CLI writes `Discarded: <summary>` and this REQ made the reader prefix-test it; confirmed by fixture (a line written exactly as the prose prescribes yields `pending`, not `cancelled`) — `impact-rule-change` → `do-work/prose-backlog.md`, with M5

**Minor findings:** 5 (report only) — M1 the reader's doc comment asserted the anchor closed the forgery; M2 an indented hand-edited `- [x]` line is skipped by the line-start scan (pre-existing, symmetric with `openQuestionPattern`); M3 no CRLF coverage; M4 an acceptance box ticked while F1 falsified it; M5 `verify-requests.md:201` restates the colon-less `Discarded`. **M1, M3 and M4 fixed**; M5 to the prose backlog with F3; M2 recorded as pre-existing and unchanged.

**Acceptance:** Partial at review time — the reported spoof was closed and both genuine terminal routes intact, but a 33-fixture differential and a live CLI attack showed the forgery still reaching a terminal status through the writer's own position. Closed by the remediation and re-verified through the built CLI.
**Suggested testing:** 4 items — the F1 lock-in (implemented), a CRLF row (implemented), a test pinning the three constants against `clarify.md`'s documented format so a prose/CLI split fails in CI rather than at review (**not implemented** — it needs a prose-parsing test, recorded below), and a manual clarify round on a REQ carrying the house `A → B` style.
**Follow-ups created:** None new — F1 and F2 were fixed here rather than deferred, and the shape-grep findings were already captured as REQ-530 before this review ran; **sweeps appended to:** None

*Reviewed by review-work action*

## Remediation

The review scored 60% with Acceptance Partial and Risk Critical, which mandates a follow-up REQ for every Important finding. Two of the three were fixed here instead, because both were this REQ's own core requirement rather than adjacent work: F1 *is* requirement 3, and F2 *is* D-01's "visible" claim. Only F3, which lives in two prose files outside this Scope, went to the backlog.

**Remediation commit:** `f1197c6`.

- **D-06** — F1 closed by refusing an `answered` summary that opens with a disposition label, not by forcing containment. **ESCALATE.** **Value:** the refusal corrects the caller's actual mistake, and the reachable path is an agent that read `clarify.md:121` and put the disposition in the summary instead of the outcome field — so the reason says exactly that. Containment would accept the misread and publish a worse record: an answered line reading `→ See contained answer note` with the confirmation buried in the dated note, arriving under `ANSWER-RAW-PAYLOAD-REQUIRED`, whose reason does not describe this cause. **Risk:** a user whose genuine answer text begins with exactly `Confirmed: ` or `Discarded: ` must reword or pick the matching outcome. Reversible in one line; the refusal is pre-mutation, so the REQ stays byte-identical.
- **D-07** — The guard fires only for `outcome: answered`, not for a summary carrying the *other* disposition's label. **DECIDE & STATE.** Measured rather than assumed: for `confirmed`/`discarded` the writer's own label already occupies the position, so `→ Discarded: Confirmed: keep the flag` prefix-tests as discarded and no verdict can move. Guarding it would be defense nothing earned. Pinned by a subtest, so a later blanket ban reddens.
- **D-08** — Reused `ANSWER-SUMMARY-INVALID` rather than minting a code. **DECIDE & STATE.** It is already that field's shape refusal (empty, multiline) and this is a third shape violation of the same field with the same remedy; refusal codes are not routed programmatically and are enumerated nowhere else, so a new one would buy only vocabulary.
- **D-09** — F2 keyed on the caller's declared `archive_path` disagreeing with the derived verdict, **not** on the blocked condition itself. **ESCALATE.** **Value:** a probe implementing the review's literal wording reddened five tests including *both* F1 forgery controls, because this writer composes an ambiguous line for any legitimate answer containing ` → ` — this project's own house style, the 2-in-124 above — so a blanket refusal would poison a REQ on every later round rather than catch an error. Keying on the declaration refuses exactly the disagreement, matching the shape `close_user_request` already enforces, and fixes the silently-dropped declaration in the same edit. **Risk:** a caller that supplies no `archive_path` and itself expects non-terminal still publishes `pending` with no evidence. That case is agreement rather than a silent fallback, and `clarify.md:139` has callers pass terminal evidence when they believe the REQ is terminal. A non-blocking channel — a typed finding on a *successful* plan — would close even that, but needs a field on `PublicationPlan` or an edit to `publication_commands.go`, both outside this write set; recorded as a discovered task rather than smuggled in.
- **D-10** — The same evidence was added to the existing `ANSWER-UR-CLOSURE-MISMATCH` reason. **DECIDE & STATE.** That refusal already fired on the identical disagreement but named no cause, leaving the caller to guess why the command disagreed.
- **D-11** — M1's comment rewritten to state the split: what the position achieves (text elsewhere on the line cannot supply a disposition; an ambiguous line fails closed) and what it cannot (it cannot tell a real label from a summary opening with those bytes), naming the write-side half without which a user's own text could still decide a terminal status.
- **D-12** — N1: the "no disposition at all" control stays as documentation and gains a sibling that pins the fail-closed direction (`- [x] Keep the flag?` with no separator at all). **DECIDE & STATE.** The original row is attributable-but-different, so it survives every mutation constructible — including last-arrow anchoring and `bytes.Contains`.
- **D-13** — M3: the CRLF row pins `findQuestionLine`'s trim, not the reader's. **DECIDE & STATE.** A trailing `\r` sits at line end and every reader test is a *prefix* test, so dropping the reader's trim changes no verdict. The reader's trim became observable only because the new evidence *quotes* the line, so the CRLF F2 subtest asserts the quoted form and reddens when that trim goes.

## Discovered Tasks

- **Same shape, second terminal status, same file:** `answer.go:268` gates the stakeholder branch's `status: completed` + `completed_at` + archive move on `bytes.Contains(ToLower(blockedHistory), "resolved")` and `bytes.Contains(ToLower(implementation), "no code")`, so `still not resolved` and `no code review yet` both satisfy it. **Captured as REQ-530** with three siblings.
- A non-blocking evidence channel (a typed finding on a successful plan) would let the non-terminal fallback report itself even when the caller declared nothing, closing D-09's residual. Needs a field on `PublicationPlan` or an edit to `publication_commands.go`.
- A test pinning the three shared constants against `clarify.md`'s documented format, so a prose/CLI split fails in CI rather than at review — this review's remaining suggested test, and the mechanism that would have caught F3 and M5 automatically.
- `questionSectionBytes` falls back to the whole body when a REQ has no `## Open Questions` heading, so P-A-U checkboxes get judged as resolved questions. Pre-existing and untouched; it now also feeds the evidence list, which is informative rather than blocking.
- `resolvedQuestionDisposition` still cannot distinguish a human hand-edit from the writer's label — a different trust domain, since a human writing `Discarded: ` into their own REQ means it.
- `internal/suiteinstall` → `TestBuiltInstallAndUpdateExit130WhenSignalsInterruptBlockedConfirmation/update-suite/TERM` failed once in three full-suite runs, 4/4 green in isolation. Third sighting of the flake tracked as REQ-525, on a third subtest.

## Lessons Learned

**What worked:**
- Asking the builder for a shape-grep rather than just a fix. The bug had a class, and the grep found three more lifecycle decisions of the same class — one driving a second terminal status in the same file. That became REQ-530, and it names two sites it deliberately *excluded* because they already fail closed, which is what makes the list trustworthy instead of a grep dump.
- Attacking through the built CLI, not the test harness. The first implementation's tests were genuine lock-ins and all green; the surviving forgery only showed up when someone submitted a manifest and read the published document.
- Measuring the accepted cost against the real corpus. "2 of 124 resolved lines, 1 historical verdict, 0 in flight" settled a risk question that no amount of reasoning about the predicate could have.

**What didn't:**
- Position anchoring alone, which is the whole of the first implementation. The disposition and the answer summary occupy the *same* position and the published line records nothing about which one wrote those bytes — so anchoring the reader could never separate them. D-05 then argued a write-side refusal would be "defense nothing earned", which is exactly the defense that was needed.
- Ticking an acceptance box for the requirement the fix did not meet. The review falsified it in one fixture.
- Implementing F2 as the review literally worded it. Refusing on the blocked condition itself reddened five tests including both forgery controls, because this writer composes an ambiguous line for any legitimate answer containing ` → ` — the house `A → B` cross-reference style. The finding was right about the symptom and wrong about the mechanism.

**Worth knowing:**
- `allSubmittedDiscarded`/`allSubmittedConfirmed` already cover the submitted set, so a mixed batch can never reach a terminal status inside one invocation. Everything `allResolvedQuestionsMatch` adds is a judgment about *earlier rounds* — which is why every reproduction of this bug needs two.
- A trailing `\r` sits at line end and every reader test here is a prefix test, so the reader's CRLF trim is invisible to any verdict assertion. It only became testable once the refusal evidence *quoted* the line.
- `questionSectionBytes` falls back to the whole body when a REQ has no `## Open Questions` heading, so P-A-U checkboxes are judged as resolved questions. Pre-existing, and now also feeding the evidence list.

## Orientation

A clarify answer can no longer talk its own REQ into being cancelled or completed: the disposition that decides a terminal status is refused on the write side when a plain answer would be indistinguishable from it, and read only where the writer places it. Lives in the CLI's `publication` package, on the `answer` path `do-work clarify` drives.

`prime_files`: `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` — spot-checked, referenced paths all still exist, no change needed. The prose contract in `skills/do-work/actions/clarify.md` disagrees with the writer in two places, both on the prose backlog rather than silently left.

**[MAP CHANGED]** — the disposition contract now has two halves that hold each other up: a write-side refusal and a position-anchored read. Neither is sufficient alone, and the doc comment on the reader says so, because the first implementation shipped one of them believing it was both.
