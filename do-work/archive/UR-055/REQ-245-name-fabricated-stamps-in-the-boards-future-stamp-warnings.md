---
id: REQ-245
title: Name fabricated stamps in the board's future-stamp warnings
status: completed
created_at: 2026-08-18T12:28:33Z
user_request: UR-055
domain: general
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec:
depends_on: []
maintenance: false
related: [REQ-244]
batch: timestamp-stamping-integrity
effort_estimate: trivial
write_set:
- skills/do-work-board/tools/queue-kanban/model.go
- skills/do-work-board/tools/queue-kanban/verify.go
- skills/do-work-board/tools/queue-kanban/timestamp_test.go
- skills/do-work-board/tools/queue-kanban/completion_anomaly_test.go
- skills/do-work-board/tools/queue-kanban/web/board-cards.js
- skills/do-work-board/tools/queue-kanban/web/board-core.js
- skills/do-work-board/tools/queue-kanban/web/board.css
- skills/do-work-board/tools/queue-kanban/prime-do-kanban.md
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-08-18T12:43:06Z
  basis:
    - trivial short-circuit
claimed_at: 2026-08-18T12:43:06Z
route: A
status_changed_at: 2026-08-18T13:52:03Z
completed_at: 2026-08-18T13:52:03Z
commit: 23bad9d
---

# Name Fabricated Stamps in the Board's Future-Stamp Warnings

## What

The board's future-stamp diagnosis messages name exactly one cause — "likely local wall-clock time stamped with a Z suffix" — but a fully fabricated value is a second, now-observed cause, and the current wording sends that reader to the wrong fix. Reword the diagnosis clauses to name both causes; keep the fix instruction (rewrite with the current UTC instant per the Timestamp rule) unchanged, since it is correct for both.

## AI Execution State (P-A-U Loop)

<!-- Filled from the builder's hand-back. In worktree dispatch mode a builder may not write
     do-work/, so the orchestrator carries the completed loop in. Source of record:
     do-work/runs/work-2026-08-18-124358/REQ-245-handback.md -->


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

## Detailed Requirements

Sibling messages found at capture — update together so they don't drift:

- `skills/do-work-board/tools/queue-kanban/model.go:379` — generate-time data warning: "…likely local wall-clock time stamped with a Z suffix; fix: rewrite with the current UTC instant…"
- `skills/do-work-board/tools/queue-kanban/model.go:1232` — reversed-span message: "…one stamp is usually local wall-clock time written with a Z suffix…"
- `skills/do-work-board/tools/queue-kanban/verify.go:371` — verify-time future `claimed_at`: "…usually local wall-clock time written with a Z suffix"

Comments asserting the single-cause story (e.g. `timestamp_test.go:42`, `completion_anomaly_test.go:227`, `model.go:1338`) should be brought in line where they would otherwise contradict the new wording.

## Constraints

- `_dev/primes/prime-kanban-board.md` governs this change — versioning, parser lock-step, build outputs. Read it before touching the tool.
- Message-text change only: no new checks, no threshold changes, the 2-minute skew allowance stays as is.
- Finding provenance (validate-feedback triage, this session): verdict Accept; Surface-cost N/A — text accuracy fix to an existing warning, no new surface.

## Red-Green Proof

**RED prompt/case:** A Go test asserting the future-stamp warning message names fabrication as a possible cause (alongside the wall-clock/Z-suffix cause) fails against the current strings.
**Why RED now:** All three diagnosis messages assert the timezone cause alone; a fabricated stamp — the observed incident — is misdiagnosed by the rendered warning.
**GREEN when:** The three messages name both causes with the fix instruction unchanged, the new assertion passes, and `go test ./...` in the tool directory exits 0.
**Validation:** Inferred during capture

## Full Context

See `do-work/user-requests/UR-055/input.md` for complete verbatim input.

---
*Source: validate-feedback Finding 3 — "Broaden the board's future-stamp warning text: 'local wall-clock time with a Z suffix' is one cause; a fully fabricated value is a second, now-observed one, and the current message sends the reader to the wrong fix."*

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/model.go` (modify) — the shared cause constant and its two call sites
- `skills/do-work-board/tools/queue-kanban/verify.go` (modify) — the verify-time future `claimed_at` Detail line
- `skills/do-work-board/tools/queue-kanban/timestamp_test.go` (modify) — the wording lock-in tests
- `skills/do-work-board/tools/queue-kanban/completion_anomaly_test.go` (modify) — single-cause comment brought in line
- `skills/do-work-board/tools/queue-kanban/web/board-cards.js` (modify) — **added during integration**, the badge tooltip
- `skills/do-work-board/tools/queue-kanban/web/board-core.js` (modify) — **added during integration**, the shared client constant and stopwatch explanation
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modify) — **added during integration**, restating comment
- `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md` (modify) — **added during integration**, restating prose

**Files I will NOT touch:** `generate_test.go` (a sibling builder owns it this wave), `future_timestamp_test.go`, `verify_test.go`, anything under `do-work/`, the version files, `CHANGELOG.md`.

**Acceptance criteria (restated from REQ):**
- [x] The diagnosis clauses name both causes — local wall-clock-with-Z **and** a fabricated value
- [x] The fix instruction is unchanged, since it is correct for both
- [x] The three sibling messages are updated together and cannot drift
- [x] No new checks, no threshold changes, the 2-minute skew allowance unchanged
- [x] Comments asserting the single-cause story brought in line

**Scope widening, logged rather than silent.** The declared `write_set` covered the Go renderers only. During integration the builder reported that `web/board-cards.js` and `web/board-core.js` hand-duplicate the same message for the badge tooltip and the stopwatch explanation, and that both still named one cause. The orchestrator folded them in rather than filing a follow-up — see **Orchestrator Decision** below. `write_set` above is mirrored from this section.

## Pre-Flight

**Git:** ✓ clean at claim (`2432f45`); builder worked in an isolated worktree on its own branch
**Tests baseline:** ✓ `maintainer-verify.sh` exit 0 before dispatch
**Dependencies:** ✓ Go 1.26.1; Node present for the strict JavaScript behavior lane

*Checked by work action*

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/model.go` (modified)
- `skills/do-work-board/tools/queue-kanban/verify.go` (modified)
- `skills/do-work-board/tools/queue-kanban/timestamp_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/completion_anomaly_test.go` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board-cards.js` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board-core.js` (modified)
- `skills/do-work-board/tools/queue-kanban/web/board.css` (modified)
- `skills/do-work-board/tools/queue-kanban/prime-do-kanban.md` (modified)

**What was done:** Hoisted the future-stamp cause clause into one Go constant naming both causes and spliced it into all three renderers, then did the same on the client side with one shared JavaScript constant behind the badge tooltip and the stopwatch explanation. A verbatim string comparison guards the Go copy against the JS copy. Every fix instruction is byte-identical to before.

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh`
**Result:** ✓ exit 0 (unpiped; `echo $?` printed `0` on its own line)

**Red-green validation:**
- `TestFutureStampDiagnosesNameBothCauses`: ✗ before implementation → ✓ after. All three failures were on `fabricated` alone while `wall-clock` and `Z suffix` passed — proof the test failed because the strings named one cause, not because a symbol was missing.
- `TestFutureStampDiagnosesKeepTheirFixInstruction`: passed in the same RED run. Correct — it is a guard against the rewording eating the remedies, not a RED.
- JS behavior probe: ✗ against the pre-change client strings for the right reason → ✓ after.

**Rendered evidence (the assertions do not cover this):** the two retired one-cause clauses appear **zero** times in the generated `index.html`, and both tooltips were read out of the live DOM on an isolated loopback origin with `location.href` and `document.title` returned from the same `evaluate` call as every measurement.

**New tests added:**
- `TestFutureStampDiagnosesNameBothCauses` — reads all three *rendered* messages and asserts each names both causes
- `TestFutureStampDiagnosesKeepTheirFixInstruction` — pins all three remedies so a future consolidation cannot flatten them
- `TestFutureStampCauseClauseMatchesTheShippedClient` — verbatim Go↔JS comparison, deliberately outside the Node-gated lane
- a `TestJavaScriptBehavior…` probe reading the rendered tooltips

**Existing tests updated (cross-REQ impact):** none — `completion_anomaly_test.go` took a comment change only, with zero test-body edits.

*Verified by work action*

## Orchestrator Decision — Scope Widened During Integration

The builder's hand-back offered to fold `web/board-cards.js` and `web/board-core.js` in, or to route them as follow-ups. **Folded in.** The REQ's stated purpose is that a reader of a future-stamp warning gets the right diagnosis; the badge tooltip is where that reader looks, and it is a hand-written duplicate of the Go string. Closing the REQ with only the Go half would have satisfied the specification and left its own goal unmet — the same "second cause makes every sentence naming the first a half-truth" pattern the REQ is about, one layer out.

`verify_test.go`'s stale fixture literals were **not** folded in and are routed as their own follow-up: they are input fixtures rather than assertions, nothing is broken, and widening the REQ twice for a cosmetic staleness is a worse trade than one small REQ.

## Decisions


- **D-01 — The cause clause is now one constant, not three strings.** `futureStampCauseClause` in `model.go` is the single source for the diagnosis; `model.go`'s two sites and `verify.go`'s one splice it in. Reach: anyone adding a fourth renderer, or a third cause, edits one line. This is slightly more than "reword three strings", but the REQ's own requirement is that the three cannot drift, and a constant is the only version of that which survives the next editor. Comments in the file now point at the constant rather than restating a cause, for the same reason.
- **D-02 — Each call site keeps its own fix instruction.** Only the cause is shared. The three remedies genuinely differ (target shape + citation; which stamp to rewrite; `queue-kanban now`), and `TestFutureStampDiagnosesKeepTheirFixInstruction` pins all three so a future consolidation cannot silently flatten them.
- **D-03 — The wording lock-in lives in `timestamp_test.go`, not in each renderer's test file.** The contract is that the three messages *agree*; asserted separately in `future_timestamp_test.go`, `completion_anomaly_test.go` and `verify_test.go`, a fourth cause added to one of them would pass all three suites. A single test that reads all three rendered messages is the only shape that catches it. Reach: if the orchestrator later prefers `future_timestamp_test.go` as the home, the two tests and three helpers move as a block with no edits — see Pushback.


- **D-04 — The client gets one shared cause string, not two corrected ones.** `futureStampCauseText` lives in `web/board-core.js` (the first fragment in `boardJavaScriptFragmentPaths`, so it is in scope for every later fragment) and both tooltips render it. This is the same delete-before-you-add move as D-01 and it shrinks the problem: the JS↔JS drift risk is gone structurally, leaving exactly one pair to guard.
- **D-05 — The Go↔JS guard is a verbatim string comparison, and the JS literal is written unbroken to make that possible.** `web/board-core.js` holds the clause on one long line rather than as a wrapped concatenation, because a split literal would force the test to reassemble JavaScript before it could compare. A comment on the literal says so, since "why is this line not wrapped like its neighbours" is otherwise a fair question. Reach: whoever edits that string must keep it on one line.
- **D-06 — `TestFutureStampCauseClauseMatchesTheShippedClient` is deliberately not in the Node-gated behavior lane.** A `t.Skip` on a Node-less machine would silently retire the only guard against the two copies drifting. The Node probe and the equality check answer different questions — "does the rendered tooltip say it" and "do the two languages say the same thing" — and only the first needs Node.
- **D-07 — The JS behavior probe went in `timestamp_test.go`, not the lane's own file.** Every existing `TestJavaScriptBehavior*` lives in `generate_test.go`, which is the sibling builder's file. The lane selects by **name pattern**, not by file or registry, so a probe named `TestJavaScriptBehavior…` participates from anywhere in the package — verified by the strict lane passing above with the new probe included. Placing it beside the Go wording tests also keeps all five assertions about this one message in one place.

## Lessons Learned


- **A message that grows a second cause makes every nearby comment a half-truth, and the comments outnumber the strings.** The REQ named three strings and three comments; the real count in `model.go` alone was two strings and four comments, two of which (`:42`, `:202`) the REQ had not listed. The tell is any comment containing "the usual cause", "the signature of", or "the specific" — those phrasings are single-cause assertions wearing a description's clothes. Grep for the *claim shape*, not just the string being edited. (This is REQ-231's lesson — "a change that adds a second cause makes every sentence naming the first one a half-truth" — reappearing in a different file; it is worth treating as a standing check rather than a per-REQ discovery.)
- **The Go string and the JS string are two renderings of one message, and only one of them is in the write set.** `web/board-cards.js` hand-duplicates `model.go`'s warning for the badge tooltip, and nothing in the suite compares them — I only found the drift by generating a board and grepping the output `index.html`, which showed the new clause and the old clause in the same file. A `write_set` scoped to the Go files silently scopes out the copy the user actually reads. Any future REQ touching a board *message* should have `web/` in its write set by default, or the pairing should be made mechanical (the way `futureInstantSkewAllowanceMs` mirrors `futureTimestampSkewAllowance` with a "keep the two in lock-step" comment — the comment exists for the constant but not for the message).
- **A comment can contradict itself in four lines and still read fine.** My first draft of the constant's doc said the old wording sent readers "to a fix that cannot help them" and then, two lines later, that "the remedy is the same either way" — both cannot be true. The old message's failure is the *diagnosis*, not the remedy. Worth re-reading a freshly written doc comment as a whole rather than as sentences.


- **The lane is a name pattern, not a registry — which is why a "you may not touch that file" constraint did not block adding a probe to it.** `TestMaintainerStrictJavaScriptBehaviorLane` re-execs the binary with `-test.run=^TestJavaScriptBehavior`, so membership is a naming convention. Worth knowing before anyone concludes that JS behavior coverage has to live in `generate_test.go`.
- **A shared browser writes into whichever repo root it considers its working directory, and that was the main tree.** The Playwright MCP dropped a console log and a page snapshot into `skill-do-work2/.playwright-mcp/` — the main tree, which builders may not write. It is gitignored and already held 36 files from sibling sessions, so this is not a git-hygiene failure, but it does mean **a builder can write outside its worktree without ever issuing a write**. I removed only my own two files (matched by timestamp) and left the siblings' alone; deleting the directory would have destroyed other agents' evidence. Any brief that says "never write in the main tree" should treat browser tooling as an exception to state explicitly rather than a rule to assume.
- **Prefer the pre-change artifact over memory when reporting a before/after count.** I had the BEFORE numbers from earlier in the session, but rebuilt the tool from `HEAD`'s blobs into `/tmp` and regenerated the board rather than quoting them — the earlier figures came from a differently-seeded fixture, and a table comparing two different fixtures would have been quietly wrong while looking authoritative.
