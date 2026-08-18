# REQ-245 Hand-Back — Name Fabricated Stamps in the Board's Future-Stamp Warnings

## Branch

`worktree-agent-REQ-245-name-fabricated-stamps-in-the-boards-future-stamp-warnings`

Implementation commit: `b74d6c8` — `[REQ-245] name fabricated stamps in the board's future-stamp diagnoses`
Merge base: `2432f450913d01e56f3685a6ca6a53c7424bd25d`

## P-A-U

- [x] **[PLAN]:** Read `_dev/primes/prime-kanban-board.md`, `CLAUDE.md`, and the crew members (`general.md`, `coding-guardrails.md`, `communication-style.md`). The three diagnosis strings are one message rendered by three call sites, so the approach is to hoist the cause clause into a single package-level constant in `model.go` (`futureStampCauseClause`) naming both causes, splice it into all three `fmt.Sprintf` sites, and leave every fix instruction byte-identical. No new checks, no threshold change, `futureTimestampSkewAllowance` untouched. RED first: one test that reads all three *rendered* messages (board warning via `futureStampSyntheticBoard`, anomaly reason via `detectCompletionAnomaly`, verify finding via `runVerifyProbes`) and asserts each names `fabricated`, `wall-clock`, and `Z suffix` — so it fails against the current strings on `fabricated` alone, not on a missing symbol. A second test pins the fix instructions so the rewording cannot quietly eat them. Single-cause comments at `model.go:42`, `model.go:202`, `model.go:1216`, `model.go:1338`, `timestamp_test.go:42`, `completion_anomaly_test.go:227` get brought in line.

- [x] **[APPLY]:** Written exactly as planned, in exactly the four files the brief allows. `generate_test.go`, `future_timestamp_test.go`, `verify_test.go` and everything under `web/` were left alone. No new checks; `futureTimestampSkewAllowance` is byte-identical; the three fix instructions are byte-identical.

- [x] **[UNIFY]:** `git diff --stat` against the merge base shows the four in-scope files and nothing else. Verified each:
  - `model.go` — read every hunk: the new constant + its doc comment, three comment rewrites, two `Sprintf` call sites now passing `futureStampCauseClause`. No behavior change beyond the string; no new branch, no new field, no threshold edit.
  - `verify.go` — one `Sprintf` Detail line; the adjacent `Remedy` comment and value untouched.
  - `timestamp_test.go` — one comment rewrite plus a new appended section (3 helpers, 2 tests). Existing tests unmodified.
  - `completion_anomaly_test.go` — comment rewrite only; zero test-body edits (checked in the diff — this file's assertions are untouched).
  - `gofmt -l skills/do-work-board/tools/queue-kanban/` → no output. `go vet ./...` clean (inside maintainer-verify).
  - No debug artifacts: no `fmt.Println`, no `t.Skip`, no commented-out code. Scratch fixture, generated board and built binary lived in `/tmp` and were deleted; `git status --porcelain` is empty and no build output was committed.
  - `git diff --name-only <base>...HEAD -- do-work/` → empty. Same for `VERSION`, `skills/do-work/VERSION`, `skills/do-work/actions/version.md`, `CHANGELOG.md` → empty.

## Files Changed

```
 .../tools/queue-kanban/completion_anomaly_test.go  |  10 +-
 skills/do-work-board/tools/queue-kanban/model.go   |  51 ++++++---
 .../tools/queue-kanban/timestamp_test.go           | 127 ++++++++++++++++++++-
 skills/do-work-board/tools/queue-kanban/verify.go  |   4 +-
 4 files changed, 165 insertions(+), 27 deletions(-)
```

- **`skills/do-work-board/tools/queue-kanban/model.go`** — added `futureStampCauseClause`, the one shared diagnosis naming both causes, next to `futureTimestampSkewAllowance`. Spliced it into the generate-time data warning (was line 379) and the reversed-span completion-anomaly reason (was line 1232). Rewrote the four single-cause comments the new wording would contradict: the skew-allowance doc (`:42`), the `FutureTimestampFields` field doc (`:202`), the `detectCompletionAnomaly` doc (`:1216`), and the `detectFutureTimestampFields` doc (`:1338`) — each now points at the constant instead of restating one cause, so there is one place to edit if a third cause ever shows up.
- **`skills/do-work-board/tools/queue-kanban/verify.go`** — the future-`claimed_at` finding's `Detail` (was line 371) now splices the same constant. Its `Remedy` (`queue-kanban now`) and the comment explaining why a command is allowed there are untouched.
- **`skills/do-work-board/tools/queue-kanban/timestamp_test.go`** — reworded the comment on `TestFormatCanonicalTimestampConvertsNonUTCZones` (was line 42), which claimed the Z-suffix bug was "the specific corruption the Timestamp rule names"; it is now one of two, and the comment says why this file still tests only that one (it is the corruption a correct *writer* can rule out — no writer can prevent fabrication). Appended the wording lock-in: three fixture helpers plus `TestFutureStampDiagnosesNameBothCauses` and `TestFutureStampDiagnosesKeepTheirFixInstruction`.
- **`skills/do-work-board/tools/queue-kanban/completion_anomaly_test.go`** — comment-only (was line 227). It cited archived REQ-091 as "the real-data case" for the Z-suffix cause; it now names both real-data cases, REQ-091 and the REQ-244/REQ-245 extrapolated stamps. No assertion or fixture changed.

**Why one constant rather than three edited strings:** the REQ asks that the three be updated together "so they cannot drift". Three hand-maintained copies is the mechanism by which they drift; one constant removes it. The clause is only the *cause* — each site keeps its own fix instruction, because their grammar and their remedies genuinely differ (the board warning cites the target shape, the anomaly reason says which stamp to rewrite, verify hands over `queue-kanban now`).

## Red-Green Evidence

**RED** — assertions written against the *current* strings, before any production edit. All three fail, and every failure is on `fabricated`; `wall-clock` and `Z suffix` pass, which is the proof the test fails because the strings name one cause and not because a symbol was missing:

```
$ go test -count=1 -run 'TestFutureStampDiagnoses' ./...
--- FAIL: TestFutureStampDiagnosesNameBothCauses (0.02s)
    timestamp_test.go:163: board future-stamp data warning (model.go) does not name "fabricated" as a possible cause — a reader hitting the other cause is sent to the wrong fix.
        got: REQ-9401 has future-dated timestamp(s): claimed_at 2026-06-30T14:00:00Z — later than the board's generation time (2min clock-skew allowance); likely local wall-clock time stamped with a Z suffix; fix: rewrite with the current UTC instant — YYYY-MM-DDTHH:MM:SSZ, per the Timestamp rule in actions/work-reference.md
    timestamp_test.go:163: reversed-span completion-anomaly reason (model.go) does not name "fabricated" as a possible cause — a reader hitting the other cause is sent to the wrong fix.
        got: completed_at "2026-01-01T10:00:00Z" is earlier than claimed_at "2026-01-02T10:00:00Z" — a reversed span cannot be real; one stamp is usually local wall-clock time written with a Z suffix; rewrite the wrong stamp with the true UTC instant
    timestamp_test.go:163: verify future-claimed_at finding (verify.go) does not name "fabricated" as a possible cause — a reader hitting the other cause is sent to the wrong fix.
        got: REQ-9341 has a future-dated claimed_at (2026-08-03T15:00:00Z) — usually local wall-clock time written with a Z suffix
FAIL
FAIL	github.com/knews2019/skill-do-work/queue-kanban	0.456s
```

Note `TestFutureStampDiagnosesKeepTheirFixInstruction` passed in that same RED run — correct, since the fix instructions were already right and this REQ must not disturb them. It is a guard, not a RED.

**GREEN** — after the constant landed:

```
$ go test -count=1 -run 'TestFutureStampDiagnoses' -v .
=== RUN   TestFutureStampDiagnosesNameBothCauses
--- PASS: TestFutureStampDiagnosesNameBothCauses (0.02s)
=== RUN   TestFutureStampDiagnosesKeepTheirFixInstruction
--- PASS: TestFutureStampDiagnosesKeepTheirFixInstruction (0.02s)
PASS
ok  	github.com/knews2019/skill-do-work/queue-kanban	0.455s

$ go test -count=1 ./...
ok  	github.com/knews2019/skill-do-work/queue-kanban	16.163s
```

**Rendered evidence** (prime: "read the rendered text when the question is 'what does this say'"). Built the binary into `/tmp`, seeded a `/tmp` fixture repo with a future `claimed_at` and a reversed span, and read the real CLI output and the real generated board — not just the assertions:

`queue-kanban summary`:
```
    ! REQ-9401 has future-dated timestamp(s): claimed_at 2026-08-18T15:51:14Z — later than the board's generation time (2min clock-skew allowance); likely a fabricated value (guessed or extrapolated instead of read from the clock) or local wall-clock time written with a Z suffix; fix: rewrite with the current UTC instant — YYYY-MM-DDTHH:MM:SSZ, per the Timestamp rule in actions/work-reference.md
    ! REQ-9402 — completed_at "2026-08-18T09:00:00Z" is earlier than claimed_at "2026-08-18T10:00:00Z" — a reversed span cannot be real; one stamp is usually a fabricated value (guessed or extrapolated instead of read from the clock) or local wall-clock time written with a Z suffix; rewrite the wrong stamp with the true UTC instant
```

`queue-kanban verify`:
```
  ! claim-needs-attention: REQ-9401 has a future-dated claimed_at (2026-08-18T15:51:14Z) — usually a fabricated value (guessed or extrapolated instead of read from the clock) or local wall-clock time written with a Z suffix
  ! completion-anomaly: REQ-9402 (status completed): completed_at "2026-08-18T09:00:00Z" is earlier than claimed_at "2026-08-18T10:00:00Z" — a reversed span cannot be real; one stamp is usually a fabricated value (guessed or extrapolated instead of read from the clock) or local wall-clock time written with a Z suffix; rewrite the wrong stamp with the true UTC instant
```

`queue-kanban generate` → parsed `board-data.js` out of the generated static board: both warnings carry the new clause, so what reaches the browser's data-warnings panel is the reworded text.

Layout check on the longer string: `.board-warnings-list` (`web/board.css:412`) is a plain padded `<ul>` with no `white-space: nowrap`, no `text-overflow`, and no fixed height, so the extra ~90 characters wrap inside the banner. No render regression is possible from the length alone; no browser run was needed. All `/tmp` scratch (fixture tree, generated board, binary) was deleted.

## Verification

```
maintainer-verify: checking Go go1.26.1
maintainer-verify: checking ShellCheck 0.11.0
maintainer-verify: ShellCheck warning-level lint (50 tracked files)
maintainer-verify: aggregate contract suite
Maintainer verification self-test passed.
Suite manifest contract probes passed.
shipped package reference contract: PASS
Shell-block lint self-test passed.
Shell-block lint passed: 74 fenced blocks and 31 shipped shell files; ShellCheck enabled.
SessionStart hook behavior probes passed.
Prescribed shell primitive canonicalization checks passed.
Defensive-surface exact deletion regressions passed.
record-commit-hash and blanked-req-scan guard probes passed.
update-script behavior probes passed.
Prescribed shell script behavior probes passed (42 named script cases).
staged skills contract: PASS
suite installer behavior probes passed.
p50 estimator suite: all probes passed.
Contract regression checks passed.
maintainer-verify: queue-kanban go vet
maintainer-verify: queue-kanban uncached ordinary tests
ok  	github.com/knews2019/skill-do-work/queue-kanban	16.143s
maintainer-verify: queue-kanban strict JavaScript behavior lane
=== RUN   TestMaintainerStrictJavaScriptBehaviorLane
--- PASS: TestMaintainerStrictJavaScriptBehaviorLane (5.09s)
PASS
ok  	github.com/knews2019/skill-do-work/queue-kanban	5.291s
maintainer-verify: audit-metrics go vet
maintainer-verify: audit-metrics uncached tests
ok  	github.com/knews2019/skill-do-work/audit-metrics	1.539s
Maintainer verification passed.
0
```

Not piped. Run from the worktree root, `echo $?` on its own line → `0`.

## Integration Seams

**Nothing to apply by hand for this REQ to merge cleanly** — no shared registry, no cross-REQ text, no version or changelog edit. Three follow-ups the orchestrator should route, all *outside* my write set (I did not touch them):

1. **`web/board-cards.js:193` still says one cause — this is the user-facing copy of the message I changed.** It is the `⚠ future stamp` badge's tooltip, and its text is a hand-written duplicate of `model.go`'s warning: `"…Likely local wall-clock time stamped with a Z suffix; fix: rewrite with the current UTC instant — YYYY-MM-DDTHH:MM:SSZ, per the Timestamp rule…"`. I confirmed against a freshly generated board that the shipped `index.html` now contains the two-cause clause in the Go warnings **and** the one-cause clause in the JS tooltip. A reader who hovers the badge — the more likely path than reading the warnings panel — still gets the old, wrong diagnosis. This is the single most important loose end.
2. **`web/board-core.js:109`** — the stopwatch's `⚠ clock skew` explanation, same single-cause wording, same reader. `web/board.css:936`'s comment and `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md:60` restate it too.
3. **`verify_test.go:1186` and `:1230`** hold a hand-typed copy of the old reversed-span reason as a fixture literal. The tests still pass (the literal is input, not an assertion against `detectCompletionAnomaly`), so nothing is broken — but the fixture now reads as a copy of a message that no longer exists, and the next person to grep for the message text will find a stale one. Worth a one-line sync when someone owns that file.

If you want (1) and (2) folded into this REQ rather than a follow-up, say so and I will do it — it is a two-string edit, and the `futureStampCauseClause` constant gives the JS side an obvious thing to mirror in the same way `futureInstantSkewAllowanceMs` already mirrors `futureTimestampSkewAllowance`.

## Decisions

- **D-01 — The cause clause is now one constant, not three strings.** `futureStampCauseClause` in `model.go` is the single source for the diagnosis; `model.go`'s two sites and `verify.go`'s one splice it in. Reach: anyone adding a fourth renderer, or a third cause, edits one line. This is slightly more than "reword three strings", but the REQ's own requirement is that the three cannot drift, and a constant is the only version of that which survives the next editor. Comments in the file now point at the constant rather than restating a cause, for the same reason.
- **D-02 — Each call site keeps its own fix instruction.** Only the cause is shared. The three remedies genuinely differ (target shape + citation; which stamp to rewrite; `queue-kanban now`), and `TestFutureStampDiagnosesKeepTheirFixInstruction` pins all three so a future consolidation cannot silently flatten them.
- **D-03 — The wording lock-in lives in `timestamp_test.go`, not in each renderer's test file.** The contract is that the three messages *agree*; asserted separately in `future_timestamp_test.go`, `completion_anomaly_test.go` and `verify_test.go`, a fourth cause added to one of them would pass all three suites. A single test that reads all three rendered messages is the only shape that catches it. Reach: if the orchestrator later prefers `future_timestamp_test.go` as the home, the two tests and three helpers move as a block with no edits — see Pushback.

## Lessons Learned

- **A message that grows a second cause makes every nearby comment a half-truth, and the comments outnumber the strings.** The REQ named three strings and three comments; the real count in `model.go` alone was two strings and four comments, two of which (`:42`, `:202`) the REQ had not listed. The tell is any comment containing "the usual cause", "the signature of", or "the specific" — those phrasings are single-cause assertions wearing a description's clothes. Grep for the *claim shape*, not just the string being edited. (This is REQ-231's lesson — "a change that adds a second cause makes every sentence naming the first one a half-truth" — reappearing in a different file; it is worth treating as a standing check rather than a per-REQ discovery.)
- **The Go string and the JS string are two renderings of one message, and only one of them is in the write set.** `web/board-cards.js` hand-duplicates `model.go`'s warning for the badge tooltip, and nothing in the suite compares them — I only found the drift by generating a board and grepping the output `index.html`, which showed the new clause and the old clause in the same file. A `write_set` scoped to the Go files silently scopes out the copy the user actually reads. Any future REQ touching a board *message* should have `web/` in its write set by default, or the pairing should be made mechanical (the way `futureInstantSkewAllowanceMs` mirrors `futureTimestampSkewAllowance` with a "keep the two in lock-step" comment — the comment exists for the constant but not for the message).
- **A comment can contradict itself in four lines and still read fine.** My first draft of the constant's doc said the old wording sent readers "to a fix that cannot help them" and then, two lines later, that "the remedy is the same either way" — both cannot be true. The old message's failure is the *diagnosis*, not the remedy. Worth re-reading a freshly written doc comment as a whole rather than as sentences.

## Pushback

**One, minor, and I complied.** The narrowed write set (`model.go`, `verify.go`, `timestamp_test.go`, `completion_anomaly_test.go`) does not contain the natural home for a future-stamp wording test. `future_timestamp_test.go` already exists, already owns `futureStampSyntheticBoard`, and already has `TestFutureTimestampWarningNamesFieldAndFix` asserting fragments of the very string I changed — a test one line away from mine in subject and one file away in location. I followed the brief exactly and put the new tests in `timestamp_test.go` with a section comment explaining the placement; the block is self-contained and moves without edits if you would rather it sat next to its sibling. Flagging it because a test placed away from its subject is how a suite rots, not because the narrowing was wrong — avoiding a collision with a sibling builder is the better trade in a parallel wave.

**Not pushback, but stated plainly:** the REQ's stated goal is that the reader of a future-stamp warning gets the right diagnosis. For the reader who hovers the board's badge, this REQ as scoped does not yet achieve that — `web/board-cards.js` is where that reader looks, and it was outside the write set. See Integration Seams (1). The REQ is complete as specified; the *outcome* it was written for needs that one further string.

---

# Addendum — JS Rendering

Folding Integration Seams (1) and (2) into REQ-245, per the orchestrator's instruction. Same branch, second commit.

**Commit:** `897edec` — `[REQ-245] name fabricated stamps in the board's two client tooltips`

## Diffstat (this commit alone)

```
 skills/do-work-board/tools/queue-kanban/model.go   |   6 +-
 .../tools/queue-kanban/prime-do-kanban.md          |   2 +-
 .../tools/queue-kanban/timestamp_test.go           | 104 +++++++++++++++++++++
 .../tools/queue-kanban/web/board-cards.js          |   9 +-
 .../tools/queue-kanban/web/board-core.js           |  18 +++-
 .../do-work-board/tools/queue-kanban/web/board.css |   7 +-
 6 files changed, 134 insertions(+), 12 deletions(-)
```

Branch total against merge base `2432f45` is now 8 files, +298/−38. `git diff --name-only <base>...HEAD` over `do-work/`, `VERSION`, `skills/do-work/VERSION`, `skills/do-work/actions/version.md`, `CHANGELOG.md` **and `generate_test.go`** → empty. The sibling's file was not touched.

## What changed

- **`web/board-core.js`** — added `futureStampCauseText`, one unbroken string literal holding the same sentence as Go's `futureStampCauseClause`, and rendered `clockSkewExplanationText` from it. Also fixed the skew-allowance comment, which said the frozen stopwatch lasts "until the wall clock catches up" — true of the timezone cause, false of a fabricated one, which never catches up.
- **`web/board-cards.js`** — the `⚠ future stamp` badge tooltip now renders `futureStampCauseText` instead of spelling the cause itself. Same "until the wall clock catches up" fix in its comment.
- **`web/board.css`** — the `.badge-future-timestamp` comment called a future stamp "the signature of local wall-clock time written with a Z suffix"; it now names the shared constant instead of restating one cause.
- **`prime-do-kanban.md:60`** — the 0.133.0 lesson said "the usual cause is local wall-clock time written with a Z suffix". It now names both causes and, more usefully, records that the sentence is one string in two languages, names both constants, and names the test that keeps them equal.
- **`model.go`** — reciprocal lock-step comment on `futureStampCauseClause` pointing at `futureStampCauseText`, in the same shape as the existing `futureTimestampSkewAllowance` ↔ `futureInstantSkewAllowanceMs` pairing.

**On item 4 (make the pairing mechanical):** I did the test rather than only the comment, and wrote the comment on both sides too. But the larger part of the answer is structural — **the client now has one cause string where it had two.** The badge tooltip and the stopwatch tooltip can no longer disagree with each other at all, because there is nothing left to disagree; only the Go↔JS pair needs a guard, and that guard is one `strings.Contains`. No sync mechanism was built.

## New tests

- **`TestJavaScriptBehaviorClockSkewTooltipNamesBothCauses`** — executes the real `syncClockSkewTitle` under Node against a stub node and asserts the `title` **a DOM element is actually given**. It participates in the strict lane by name (`-test.run=^TestJavaScriptBehavior`) and increments the probe count, so it is a real behavior probe, not a parse check.
- **`TestFutureStampCauseClauseMatchesTheShippedClient`** — asserts `web/board-core.js` carries `futureStampCauseClause` verbatim, and that `web/board-cards.js` contains no `Z suffix` literal at all (it must render the shared constant, never its own copy). Deliberately **not** Node-gated: the drift it guards is just as real on a machine without Node, and a `t.Skip` would silently retire the check there.

One helper was needed: `sliceJavaScriptStatementsThrough`, the `var`-statement analogue of the suite's existing `sliceBalancedBlockAfter` (a `var` initialised with string literals has no braces to balance). It is anchored on a region spanning both constants, so the probe's slice list did not have to change when the shared constant was introduced between them — which is why the RED below and the GREEN run the identical extraction.

## Red-Green Evidence (JS)

**RED**, written against the pre-change JS. The probe **executed successfully and printed the real rendered tooltip** — it fails on the assertion, not on a missing symbol:

```
$ go test -count=1 -run 'TestJavaScriptBehaviorClockSkewTooltipNamesBothCauses|TestFutureStampCauseClauseMatchesTheShippedClient' .
--- FAIL: TestJavaScriptBehaviorClockSkewTooltipNamesBothCauses (0.35s)
    timestamp_test.go:264: the clock-skew tooltip does not name "fabricated" as a possible cause — the reader who hovers a frozen stopwatch is sent to the wrong fix.
        got: This timestamp is ahead of your clock by more than the 2-minute skew allowance — likely stamped with local wall-clock time plus a Z suffix. Fix the frontmatter with the current UTC instant — YYYY-MM-DDTHH:MM:SSZ, per the Timestamp rule in actions/work-reference.md. Until then the stopwatch cannot measure real elapsed time.
--- FAIL: TestFutureStampCauseClauseMatchesTheShippedClient (0.00s)
    timestamp_test.go:291: web/board-core.js does not carry futureStampCauseClause verbatim, so the board's Go and JS diagnoses disagree.
        want the literal: "a fabricated value (guessed or extrapolated instead of read from the clock) or local wall-clock time written with a Z suffix"
    timestamp_test.go:302: web/board-cards.js spells the future-stamp cause itself; it must render futureStampCauseText so the badge and stopwatch tooltips cannot disagree
FAIL
FAIL	github.com/knews2019/skill-do-work/queue-kanban	0.838s
```

Every failure is on the missing second cause; the fix-instruction assertions passed in the same run, as they should.

**GREEN**:

```
$ go test -count=1 -run 'TestJavaScriptBehaviorClockSkewTooltipNamesBothCauses|TestFutureStampCauseClauseMatchesTheShippedClient|TestFutureStampDiagnoses' -v .
=== RUN   TestFutureStampDiagnosesNameBothCauses
--- PASS: TestFutureStampDiagnosesNameBothCauses (0.02s)
=== RUN   TestFutureStampDiagnosesKeepTheirFixInstruction
--- PASS: TestFutureStampDiagnosesKeepTheirFixInstruction (0.02s)
=== RUN   TestJavaScriptBehaviorClockSkewTooltipNamesBothCauses
--- PASS: TestJavaScriptBehaviorClockSkewTooltipNamesBothCauses (0.29s)
=== RUN   TestFutureStampCauseClauseMatchesTheShippedClient
--- PASS: TestFutureStampCauseClauseMatchesTheShippedClient (0.00s)
PASS
ok  	github.com/knews2019/skill-do-work/queue-kanban	0.766s
```

## Generated-`index.html` grep counts

Built the tool twice against one fixture repo (a claimed REQ with a `claimed_at` three hours ahead) and generated a board from each: **BEFORE** is `HEAD` — the Go half already fixed, the JS half not, i.e. exactly the drift I reported — and **AFTER** is this commit. Counts are over the generated `index.html`:

| string in the generated board | BEFORE | AFTER |
|---|---|---|
| `wall-clock time stamped with a Z suffix` (retired badge clause) | 1 | **0** |
| `likely stamped with local wall-clock time plus a Z suffix` (retired stopwatch clause) | 1 | **0** |
| `fabricated` | 1 | 2 |

Both retired one-cause clauses are at **zero**. The BEFORE `fabricated` count of 1 is an unrelated pre-existing comment about never fabricating a Markdown fence; AFTER it is that same comment plus `futureStampCauseText`. So the shipped client carries exactly one statement of the cause, and it names both. (`board-data.js` separately carries 1, the Go-produced warning.)

## Rendered evidence — the live board in a browser

The grep proves the string is in the bundle; it does not prove `futureStampCauseText` **resolves** for `board-cards.js`, which is a different fragment. That is precisely the class of defect a passing string assertion sails past, so I served the fixture on an isolated loopback port and read both tooltips out of the live DOM. `location.href` and `document.title` were returned from the *same* `evaluate` call as every measurement, per the prime's rule:

```
href:      http://127.0.0.1:8747/
pageTitle: req245-js-fixture — do-work queue board
cardCount: 1
```

Badge tooltip, read from `.badge-future-timestamp`:

```
Future-dated timestamp(s): claimed_at 2026-08-18T16:06:55Z — later than the board's generation time
(2min skew allowance). Likely a fabricated value (guessed or extrapolated instead of read from the
clock) or local wall-clock time written with a Z suffix; fix: rewrite with the current UTC instant —
YYYY-MM-DDTHH:MM:SSZ, per the Timestamp rule in actions/work-reference.md.
```

Stopwatch tooltip, read from the `⚠ clock skew` duration node:

```
This timestamp is ahead of your clock by more than the 2-minute skew allowance — likely a fabricated
value (guessed or extrapolated instead of read from the clock) or local wall-clock time written with
a Z suffix. Fix the frontmatter with the current UTC instant — YYYY-MM-DDTHH:MM:SSZ, per the
Timestamp rule in actions/work-reference.md. Until then the stopwatch cannot measure real elapsed time.
```

Data-warnings panel (the Go-produced half, same page): carries the same clause. One console error on the page, checked rather than assumed — a pre-existing `favicon.ico` 404, unrelated to this change. Server killed and all `/tmp` scratch removed afterwards.

## Verification

```
maintainer-verify: checking Go go1.26.1
maintainer-verify: checking ShellCheck 0.11.0
maintainer-verify: ShellCheck warning-level lint (50 tracked files)
maintainer-verify: aggregate contract suite
Maintainer verification self-test passed.
Suite manifest contract probes passed.
shipped package reference contract: PASS
Shell-block lint self-test passed.
Shell-block lint passed: 74 fenced blocks and 31 shipped shell files; ShellCheck enabled.
SessionStart hook behavior probes passed.
Prescribed shell primitive canonicalization checks passed.
Defensive-surface exact deletion regressions passed.
record-commit-hash and blanked-req-scan guard probes passed.
update-script behavior probes passed.
Prescribed shell script behavior probes passed (42 named script cases).
staged skills contract: PASS
suite installer behavior probes passed.
p50 estimator suite: all probes passed.
Contract regression checks passed.
maintainer-verify: queue-kanban go vet
maintainer-verify: queue-kanban uncached ordinary tests
ok  	github.com/knews2019/skill-do-work/queue-kanban	15.670s
maintainer-verify: queue-kanban strict JavaScript behavior lane
=== RUN   TestMaintainerStrictJavaScriptBehaviorLane
--- PASS: TestMaintainerStrictJavaScriptBehaviorLane (4.95s)
PASS
ok  	github.com/knews2019/skill-do-work/queue-kanban	5.220s
maintainer-verify: audit-metrics go vet
maintainer-verify: audit-metrics uncached tests
ok  	github.com/knews2019/skill-do-work/audit-metrics	1.626s
Maintainer verification passed.
0
```

Not piped. From the worktree root, `echo $?` on its own line → `0`.

## Decisions (continuing from D-03)

- **D-04 — The client gets one shared cause string, not two corrected ones.** `futureStampCauseText` lives in `web/board-core.js` (the first fragment in `boardJavaScriptFragmentPaths`, so it is in scope for every later fragment) and both tooltips render it. This is the same delete-before-you-add move as D-01 and it shrinks the problem: the JS↔JS drift risk is gone structurally, leaving exactly one pair to guard.
- **D-05 — The Go↔JS guard is a verbatim string comparison, and the JS literal is written unbroken to make that possible.** `web/board-core.js` holds the clause on one long line rather than as a wrapped concatenation, because a split literal would force the test to reassemble JavaScript before it could compare. A comment on the literal says so, since "why is this line not wrapped like its neighbours" is otherwise a fair question. Reach: whoever edits that string must keep it on one line.
- **D-06 — `TestFutureStampCauseClauseMatchesTheShippedClient` is deliberately not in the Node-gated behavior lane.** A `t.Skip` on a Node-less machine would silently retire the only guard against the two copies drifting. The Node probe and the equality check answer different questions — "does the rendered tooltip say it" and "do the two languages say the same thing" — and only the first needs Node.
- **D-07 — The JS behavior probe went in `timestamp_test.go`, not the lane's own file.** Every existing `TestJavaScriptBehavior*` lives in `generate_test.go`, which is the sibling builder's file. The lane selects by **name pattern**, not by file or registry, so a probe named `TestJavaScriptBehavior…` participates from anywhere in the package — verified by the strict lane passing above with the new probe included. Placing it beside the Go wording tests also keeps all five assertions about this one message in one place.

## Lessons Learned (addendum)

- **The lane is a name pattern, not a registry — which is why a "you may not touch that file" constraint did not block adding a probe to it.** `TestMaintainerStrictJavaScriptBehaviorLane` re-execs the binary with `-test.run=^TestJavaScriptBehavior`, so membership is a naming convention. Worth knowing before anyone concludes that JS behavior coverage has to live in `generate_test.go`.
- **A shared browser writes into whichever repo root it considers its working directory, and that was the main tree.** The Playwright MCP dropped a console log and a page snapshot into `skill-do-work2/.playwright-mcp/` — the main tree, which builders may not write. It is gitignored and already held 36 files from sibling sessions, so this is not a git-hygiene failure, but it does mean **a builder can write outside its worktree without ever issuing a write**. I removed only my own two files (matched by timestamp) and left the siblings' alone; deleting the directory would have destroyed other agents' evidence. Any brief that says "never write in the main tree" should treat browser tooling as an exception to state explicitly rather than a rule to assume.
- **Prefer the pre-change artifact over memory when reporting a before/after count.** I had the BEFORE numbers from earlier in the session, but rebuilt the tool from `HEAD`'s blobs into `/tmp` and regenerated the board rather than quoting them — the earlier figures came from a differently-seeded fixture, and a table comparing two different fixtures would have been quietly wrong while looking authoritative.

## Pushback (addendum)

None on this instruction — the call to fold it in was right, and the badge tooltip turned out to be the reader-facing half of the message. One thing worth recording rather than acting on: `verify_test.go`'s stale fixture literals (original Integration Seams item 3) are now **more** stale, since the reversed-span message they copy has moved twice. Still input fixtures, still passing, still outside my write set, and I agree with routing them as their own small REQ — but if that REQ is not written soon, the next person greping for the reversed-span text will find two copies of a message that no longer exists.

---

# Addendum 2 — Guard Holes and the Fourth Renderer

Review findings 1–4. **Findings 2 and 3 are correct and they are my errors** — I shipped a test whose comment claimed it prevented badge drift when it only forbade one literal, which is worse than shipping no test, because it tells the next reader the case is covered. Finding 4 is my miss too: my own lesson about sweeping for the *claim shape* rather than the changed string, landing one line short in the file I had just edited.

**Commit:** `14103cc` — `[REQ-245] close the badge-tooltip guard hole and the fourth renderer`

Branch was merged with `main` first (`git merge main`, clean) so `forensics.md` was edited at its post-REQ-244 text.

## Diffstat

```
 skills/do-work-board/tools/queue-kanban/model.go   | 30 +++++++---
 .../tools/queue-kanban/prime-do-kanban.md          |  2 +-
 .../tools/queue-kanban/timestamp_test.go           | 66 +++++++++++++++++-----
 .../tools/queue-kanban/web/board-cards.js          | 24 +++++---
 .../tools/queue-kanban/web/board-core.js           |  9 +--
 skills/do-work/actions/forensics.md                |  4 +-
 6 files changed, 97 insertions(+), 38 deletions(-)
```

Merge base with `main` is now `a8ef062`, so `git diff --name-only <base>...HEAD` lists exactly those six files — `do-work/`, `VERSION`, `version.md`, `CHANGELOG.md` and `generate_test.go` are all absent. (Worth noting for your merge check: because the branch now contains a merge of `main`, a **two**-dot diff would show every `do-work/` path main carries. The three-dot form you already use stays correct.)

## Finding 2 + 3 — the badge guard now guards

The tooltip was an inline concatenation buried in `makeRequestCard`, which is ~226 lines and needs `requestsById`, `createElement` and a pile of formatters — that is why I reached for a source-level check the first time, and it was the wrong trade. Extracted `futureStampTooltipText(futureTimestampFields)` in `web/board-cards.js`: a pure function taking the fields and returning the title string, with `makeRequestCard` reduced to `futureStampBadge.title = futureStampTooltipText(request.futureTimestampFields)`. A pure function needs no DOM stub at all, so the probe drives the real thing.

`TestJavaScriptBehaviorFutureStampBadgeTooltipNamesBothCauses` builds the tooltip under Node and asserts the rendered text names all three markers plus every fix fragment. The negative `Z suffix` check it subsumes is **deleted**, not kept alongside — a guard that cannot fail is worse than absent.

**The mutation that used to pass and now fails.** I reproduced the reviewer's edit exactly, replacing `futureStampCauseText +` with `"a misconfigured clock" +`, and ran the old guard and the new probe together:

```
$ go test -count=1 -run 'TestFutureStampCauseClauseMatchesTheShippedClient|TestJavaScriptBehaviorFutureStampBadgeTooltipNamesBothCauses' -v .
=== RUN   TestJavaScriptBehaviorFutureStampBadgeTooltipNamesBothCauses
    timestamp_test.go:308: the badge tooltip does not name "fabricated" as a possible cause — the reader who hovers the card is sent to the wrong fix.
        got: Future-dated timestamp(s): claimed_at 2026-06-30T14:00:00Z — later than the board's generation time (2min skew allowance). Likely a misconfigured clock; fix: rewrite with the current UTC instant — YYYY-MM-DDTHH:MM:SSZ, per the Timestamp rule in actions/work-reference.md.
    timestamp_test.go:308: the badge tooltip does not name "wall-clock" as a possible cause — …
    timestamp_test.go:308: the badge tooltip does not name "Z suffix" as a possible cause — …
--- FAIL: TestJavaScriptBehaviorFutureStampBadgeTooltipNamesBothCauses (0.30s)
=== RUN   TestFutureStampCauseClauseMatchesTheShippedClient
--- PASS: TestFutureStampCauseClauseMatchesTheShippedClient (0.00s)
FAIL
```

That single run is the whole finding: **the source guard passes the mutation, the rendered-text probe fails it**, and the probe's failure prints the actual wrong tooltip a user would have read. Mutation reverted from a pristine copy afterwards; `git diff` on `board-cards.js` confirms only the intended refactor remains.

**GREEN**, unmutated:

```
$ go test -count=1 -run 'TestFutureStamp|TestJavaScriptBehaviorFutureStamp|TestJavaScriptBehaviorClockSkew' -v .
--- PASS: TestFutureStampDiagnosesNameBothCauses (0.02s)
--- PASS: TestFutureStampDiagnosesKeepTheirFixInstruction (0.02s)
--- PASS: TestJavaScriptBehaviorClockSkewTooltipNamesBothCauses (0.37s)
--- PASS: TestJavaScriptBehaviorFutureStampBadgeTooltipNamesBothCauses (0.28s)
--- PASS: TestFutureStampCauseClauseMatchesTheShippedClient (0.00s)
PASS
```

**Finding 3 accepted as stated.** I did not harden the string search — the reviewer is right that the binding trick is contrived and the realistic edit-in-place case is caught. Instead the surviving check's doc comment now states its real guarantee (present in the file, not bound to the variable), names the probes as what covers the rest, and D-06 is amended below.

## Finding 4 — the surviving phrasing

`web/board-core.js:225`'s comment on `syncClockSkewTitle` said the tooltip is "removed the moment the wall clock catches up". Now: "removed the moment the stamp is corrected (or the wall clock catches up)" — both resolution paths, fabricated stamps included.

I then swept the repo for the whole *class* rather than that one line — `wall clock catches up`, `stamped with a Z suffix`, `stamped with local`, `usual cause`, `specific corruption`, `signature of local`. The only live source hit left is the new `board-core.js` wording above. Everything else is `do-work/` history, `CHANGELOG.md` release notes describing what shipped at the time, an `ai-reports/` narrative, and `record-commit-hash.sh:158`'s "usual cause" about an unrelated subject.

## Finding 1 — the fourth renderer, folded in

`skills/do-work/actions/forensics.md` check 12 now reads:

> "REQ-NNN's `{field}` is `{value}` — {N} in the future. Likely local wall-clock time written with a `Z` suffix, or a fabricated value — guessed or extrapolated instead of read from the clock (the Timestamp rule in `actions/work-reference.md` requires the current UTC instant). **For as long as the stamp stands**, elapsed-time math built on it is wrong: …"

The "Until the wall clock catches up" framing is gone here too. The suggested fix was already right for both causes and is untouched.

Because it is a different skill package, nothing can hold it in step mechanically, so the pointer goes in both directions: `futureStampCauseClause`'s doc comment names it and says to change it in the same commit, and the `prime-do-kanban.md` bullet does the same. I deliberately did **not** number the copies ("a fourth", "a fifth") in those pointers — a count is a closed enumeration and goes stale the moment a renderer is added or removed.

**The gate caught a real error here.** My first version of the prime pointer cited `actions/forensics.md`, and the staged-skills contract failed: `prime-do-kanban.md` ships inside `do-work-board`, where that path does not resolve. Corrected to the file's existing cross-package convention, `../do-work/actions/forensics.md`.

```
unresolved staged runtime references in do-work-board:
tools/queue-kanban/prime-do-kanban.md:60: actions/forensics.md
FAIL: staged skills contract probes failed
```

## The clause-order call — taking the reviewer's version

Reversed to `local wall-clock time written with a Z suffix, or a fabricated value (guessed or extrapolated instead of read from the clock)`, applied to all five renderers.

I argued myself out of my original order. The instinct against it was "lead with the newly-observed cause, since that is the reader we misdirected" — but that confuses emphasis with reachability. What matters is where a skimmer's eye stops, and in my order the two causes were separated by a 66-character parenthetical, so cause two was only reachable *through* the elaboration. In the reviewer's order both causes land inside the first 67 characters and the elaboration is deferred to the end, where stopping early costs nothing. Same information, strictly earlier. It is the same defect this REQ exists to fix, one notch milder, and the reviewer was right to name it.

All assertions are order-independent (`fabricated`, `wall-clock`, `Z suffix` as separate `Contains` checks), so no test needed changing to accommodate the swap — which is itself a small argument that they were asserting the right thing.

## Rendered evidence

Regenerated the board and re-read every renderer, because the order changed and a table of counts says nothing about whether the new sentence reads. All five, live:

- board data warning and reversed-span anomaly reason — via `summary`
- verify's future-`claimed_at` finding — via `verify`
- badge tooltip, stopwatch tooltip, anomaly badge tooltip — read from a live board at `http://127.0.0.1:8751/`, with `location.href` and `document.title` returned from the same `evaluate` call (`req245b-fix — do-work queue board`, the fixture I created)

```
badge:     …(2min skew allowance). Likely local wall-clock time written with a Z suffix, or a fabricated
           value (guessed or extrapolated instead of read from the clock); fix: rewrite with the current
           UTC instant — YYYY-MM-DDTHH:MM:SSZ, per the Timestamp rule in actions/work-reference.md.
stopwatch: …by more than the 2-minute skew allowance — likely local wall-clock time written with a Z
           suffix, or a fabricated value (guessed or extrapolated instead of read from the clock). Fix the
           frontmatter with the current UTC instant — …
```

Retired-clause counts in the freshly generated `index.html`, including the mutation string as a control:

| string | count |
|---|---|
| `wall-clock time stamped with a Z suffix` | 0 |
| `likely stamped with local wall-clock time plus a Z suffix` | 0 |
| `Likely local wall-clock time stamped` | 0 |
| `a misconfigured clock` (mutation control) | 0 |

Server killed, `/tmp` scratch removed, and my two `.playwright-mcp/` files removed from the main tree (siblings' left intact, as before).

## Verification

```
maintainer-verify: checking Go go1.26.1
maintainer-verify: checking ShellCheck 0.11.0
maintainer-verify: ShellCheck warning-level lint (50 tracked files)
maintainer-verify: aggregate contract suite
Maintainer verification self-test passed.
Suite manifest contract probes passed.
shipped package reference contract: PASS
Shell-block lint self-test passed.
Shell-block lint passed: 74 fenced blocks and 31 shipped shell files; ShellCheck enabled.
SessionStart hook behavior probes passed.
Prescribed shell primitive canonicalization checks passed.
Defensive-surface exact deletion regressions passed.
record-commit-hash and blanked-req-scan guard probes passed.
update-script behavior probes passed.
Prescribed shell script behavior probes passed (42 named script cases).
staged skills contract: PASS
suite installer behavior probes passed.
p50 estimator suite: all probes passed.
Contract regression checks passed.
maintainer-verify: queue-kanban go vet
maintainer-verify: queue-kanban uncached ordinary tests
ok  	github.com/knews2019/skill-do-work/queue-kanban	15.705s
maintainer-verify: queue-kanban strict JavaScript behavior lane
=== RUN   TestMaintainerStrictJavaScriptBehaviorLane
--- PASS: TestMaintainerStrictJavaScriptBehaviorLane (5.44s)
PASS
ok  	github.com/knews2019/skill-do-work/queue-kanban	5.686s
maintainer-verify: audit-metrics go vet
maintainer-verify: audit-metrics uncached tests
ok  	github.com/knews2019/skill-do-work/audit-metrics	1.410s
Maintainer verification passed.
0
```

Not piped. `echo $?` on its own line → `0`.

## Decisions (continuing from D-07)

- **D-06 — AMENDED.** The original claim, that splitting the guard into a Node-gated probe and an always-on equality check means the pair is covered with or without Node, is **too strong**. What the always-on check actually proves is that the sentence is *present in* `board-core.js`, not that `futureStampCauseText` is *bound to* it; binding something else while the literal survives in a comment passes. The honest statement: without Node, the guard catches the realistic failure (someone edits one language and not the other) and misses the contrived one (someone rebinds the variable); with Node, the two behavior probes close it. The test's doc comment now says exactly this rather than implying total coverage.
- **D-08 — A string that must be asserted gets extracted into a function, even when its call site is untestable.** `futureStampTooltipText` exists because `makeRequestCard` is too large to drive and the tooltip is too important to leave unasserted. The general rule: when the cheap guard you can write does not actually guard, the answer is to make the thing addressable, not to accept the cheap guard. Reach: any future card badge whose text carries a diagnosis should be extracted the same way.
- **D-09 — Cross-package copies get a pointer, never a count.** `forensics.md` cannot be reached from the board module, so the only available mechanism is a sentence telling the next editor it exists. Both `model.go` and `prime-do-kanban.md` carry it, phrased as a condition ("one more copy lives outside this module … change it in the same commit") rather than an enumeration, per CLAUDE.md's Closed Enumerations Go Stale.
- **D-10 — Both causes are named before any elaboration.** Reviewer's ordering, adopted for all five renderers, with the reason recorded in `futureStampCauseClause`'s doc comment and the prime bullet so a future editor does not "improve" the sentence by moving the parenthetical back between the causes.

## Lessons Learned (addendum 2)

- **A negative assertion is the weakest possible way to pin a string, and it reads like the strongest.** `assert board-cards.js contains no "Z suffix"` looks like it constrains the tooltip; it constrains only one spelling of one cause, and every wrong tooltip that avoids those two words passes. The tell is that the assertion never mentions the thing it claims to protect — it named a file and a substring, not a tooltip. **If a test's comment describes an outcome the assertion cannot observe, the assertion is wrong.** That is the shape to grep for in review, and it is exactly what a mutation catches in one run.
- **"Too expensive to test" is usually "not yet extracted".** I rejected the badge probe as needing a 226-line function stubbed, and settled for a guard that did not guard. The actual cost of getting it right was a nine-line pure function. When the reason for a weak test is the size of the enclosing function, that is an argument for extraction, not for the weak test.
- **A doc comment that overstates a guarantee is a defect with the same shape as the one this REQ fixes.** D-06 told a future reader that the Node-less path was covered when it was not — the same class of error as a message naming one cause when there are two: a confident sentence that sends the next person somewhere useless. Guarantees in comments deserve the same "is this literally true" pass as the strings they describe.
- **The sweep that found finding 4 for the reviewer is one I had already written down and then did not run to completion.** My own Addendum-1 lesson said to grep for the claim shape ("the usual cause", "until the wall clock catches up") rather than the changed string; I applied it to `model.go` and not to the JS file I was editing in the same commit. A lesson recorded in a hand-back is not a lesson applied — the sweep has to be a step, not an intention, which is why this time I ran it across the whole repo and pasted the result rather than reasoning about where it was likely to hit.
