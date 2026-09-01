---
id: REQ-280
title: Probe timestamp ordering, and point Check 12 at the archive repair that already ships
status: completed
created_at: 2026-08-19T13:42:45Z
claimed_at: 2026-08-20T23:48:22Z
completed_at: 2026-08-20T23:57:51Z
kb_status: promoted
kb_entry: REQ-280-probe-timestamp-ordering-and-point-check.md
commit: 5e180d0
route: B
user_request: UR-057
domain: general
prime_files: [_dev/primes/prime-kanban-board.md, _dev/primes/prime-action-files.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
related: [REQ-279, REQ-281, REQ-282, REQ-283]
batch: upstream-consumer-report-2026-08-19
write_set:
- skills/do-work-board/tools/queue-kanban/verify.go
- skills/do-work-board/tools/queue-kanban/verify_test.go
- skills/do-work/actions/forensics.md
---

# Probe Timestamp Ordering, and Point Check 12 at the Archive Repair That Already Ships

## What

Two gaps in the same loop, both on the read side.

1. **No ordering probe.** `skills/do-work-board/tools/queue-kanban/model.go:1261` is the suite's entire time-consistency surface: one comparison, `completed_at < claimed_at`, inside `detectCompletionAnomaly`. Nothing anywhere checks `created_at <= claimed_at <= completed_at`. Add that ordering probe to `queue-kanban verify`'s probe set and to `actions/forensics.md` Check 12, reported per violated pair.

2. **Check 12's remedy predates the repair.** Check 12 still tells the reader to recover the true instant "from the REQ file's git history" by hand. `skills/do-work/scripts/audit-archive-timestamps.sh` has since shipped and does exactly that mechanically — `git blame --line-porcelain` on the stamp's own line, sharing the repairer's predicate by sourcing `scripts/repair-req-timestamps.sh`. It is documented in exactly one place, `actions/capture.md:67`, as an immutability-rule exception. The diagnostic that finds the damage never names the tool that fixes it.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [x] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [x] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

A fabricated `claimed_at` is invisible unless it happens to land after `completed_at`. In the reporting consumer repo, 26 archived REQs carried a `claimed_at` byte-identical to their `estimate.calculated_at` with both rounded to `:00` seconds, five of them sharing one instant; only the two that inverted surfaced anywhere. Check 12 compares each stamp against **now**, so a past-dated fabrication passes. The ordering probe catches the class instead of the inverted subset, and it catches it where nothing auto-repairs — the archive, which the SessionStart repairer deliberately does not touch.

The Check 12 pointer is the cheaper half and closes the loop the report described as "opened and abandoned": the repair exists, the diagnostic just never says so.

## Context

- `scripts/repair-req-timestamps.sh` (REQ-246) already implements the `created_at <= claimed_at <= completed_at` predicate as a repair, over `do-work/queue/` and `do-work/working/`, run from the SessionStart hook.
- `scripts/audit-archive-timestamps.sh` extends the same predicate to `do-work/archive/`, report-only by default, `--fix` to write, deliberately never hook-wired. Its header states the constraint: "repairing the archive stays a conscious, deliberate invocation, which is exactly what keeps the exception narrower than the rule."
- So the *predicate* ships twice and the *repair* ships twice; what is missing is a read-side probe and a working pointer.
- `verify.go:119` already calls `appendCompletionAnomalyFindings`, which is the natural family for the new probe. REQ-214 (archived, UR-048) established that pattern — lifting a board-detected anomaly class into a verify finding.

## Detailed Requirements

- Add an ordering probe to `runVerifyProbes` reporting one finding per violated pair, naming the REQ, both field names, and both raw values. `created_at > claimed_at` and `claimed_at > completed_at` are separate findings when both are violated.
- Cover archived REQs, not only queue and working — the archive is where nothing auto-repairs.
- The finding's remedy must name `scripts/audit-archive-timestamps.sh` for archive files and note that queue/working files are repaired by the SessionStart hook, so a reader is never sent to do by hand what already runs.
- Do not mark the finding `[fixable]` unless `do-work cleanup` actually resolves it; the `[fixable]` marker means cleanup can mechanically fix it.
- Extend `actions/forensics.md` Check 12 with the ordering condition alongside its existing future-stamp condition, and rewrite its Suggested fix to point at `scripts/audit-archive-timestamps.sh` (archive) and the SessionStart repairer (queue/working) instead of hand git archaeology.
- Add the new probe to Check 14's probe table so the documented probe set stays complete.

## Constraints

- **Do not build the fabrication heuristic.** The upstream report's remedy 2 asked for a warning when `claimed_at == estimate.calculated_at` with `:00` seconds, or when three or more REQs share an instant. Declined at triage: that equality fires legitimately on 20 of this repo's 251 stamped REQs (8%), because Step 2's claim and Step 3.6's estimate can read the same instant. Under `crew-members/coding-guardrails.md` § 2's earned-defense rubric it is an unearned warning apparatus — it names a suspicion no path can act on, on top of two repairers that already fix every provable case. Out of scope; do not add it opportunistically.
- **Do not add a fourth timestamp predicate.** The ordering rule already exists in `scripts/repair-req-timestamps.sh`. State it once more in Go only because the read side cannot source shell; do not reword it, and do not let the Go and shell spellings drift on the boundary (which comparison is strict, what an absent or unparseable stamp means).
- Absent or unparseable stamps stay other checks' territory, matching `detectCompletionAnomaly`'s existing carve-out at `model.go:1258`.
- Write-set overlap: REQ-281 also edits `verify.go` and depends on this REQ, so this one lands first.

## Red-Green Proof

**RED prompt/case:** Build a verify fixture with an archived REQ carrying `created_at: 2026-08-19T12:00:00Z`, `claimed_at: 2026-08-18T09:00:00Z`, `completed_at: 2026-08-19T14:00:00Z` — claimed before it existed, but with a forward `claimed_at → completed_at` span so `detectCompletionAnomaly` stays silent. Run `queue-kanban verify --repo-root <fixture>`.
**Why RED now:** It exits 0 with `OK: no findings`. `model.go:1261` only compares `completed_at` against `claimed_at`, so an impossible `created_at` ordering passes every check the suite has.
**GREEN when:** The same fixture exits 1 with a finding naming REQ, both fields, both values, and a remedy pointing at `scripts/audit-archive-timestamps.sh`. Second assertion in the same test: `actions/forensics.md` Check 12 states the ordering condition and its Suggested fix names the script, so the doc and the probe cannot drift apart.
**Validation:** Inferred during capture. The ordering gap is confirmed by reading `model.go:1240-1266` and by `scripts/audit-archive-timestamps.sh .` reporting "archive audit clean (257 file(s) scanned)" on this repo while four REQs here (REQ-230, REQ-233, REQ-234, REQ-236) carry `claimed_at == calculated_at` at `:00` seconds in two shared-instant pairs — internally consistent, so invisible to every predicate that ships.

## Full Context

See `do-work/user-requests/UR-057/input.md` for the complete verbatim upstream report.

---
*Source: upstream defect report D1, severity high, from `g1w-game-find-the-difference` running v0.212.25 — verbatim claim: "A fabricated `claimed_at` is undetectable, unrepairable, and feeds the estimator … Nothing checks `created_at ≤ claimed_at ≤ completed_at`." Accepted (narrowed) by `do-work-toolbox validate-feedback` triage (2026-08-19); the report's remedy 1 plus the residue of its remedy 4, which triage found **already done** as `scripts/audit-archive-timestamps.sh`. Evidence: `skills/do-work-board/tools/queue-kanban/model.go:1261` is the single comparison; `skills/do-work/actions/forensics.md:156` Check 12 compares against now only and its remedy names no tool; `skills/do-work/scripts/audit-archive-timestamps.sh` ships the repair and is cited only at `skills/do-work/actions/capture.md:67`. Surface-cost: Earned — incident is the consumer's 26 fabricated stamps plus four reproducing here; surface is one ordering comparison in an existing probe family and one rewritten remedy line; cheaper than the per-REQ hand git-archaeology Check 12 currently prescribes; test is the reversed-ordering verify fixture above.*

---

## Triage

**Route: B** - Medium

**Reasoning:** The "what" is fully specified — the REQ names the probe family, the predicate, the fixture, the remedy routing, and two explicit non-goals. What needed discovery was where the probe hooks into `runVerifyProbes`, how a ticket knows whether it lives in the archive, and how the `[fixable]` marker is defined. No architectural choice was open.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

- `verify.go:113` `runVerifyProbes` is a flat list of `append*Findings` calls; `appendCompletionAnomalyFindings` at :288 is the named family and the natural neighbour.
- `VerifyFinding{Category, Detail, Fixable, Remedy}` at :63. `Fixable` is documented at :60 as meaning specifically "`do-work cleanup` can mechanically resolve it" — so the REQ's instruction not to mark it fixable is already the field's own contract.
- `board.AllRequests` (`model.go:371`) is every parsed REQ in id order across queue, working and archive — the coverage the REQ requires, with no second walk.
- `RequestTicket.TreeSection` (`model.go:236`) is `"queue" | "working" | "archive"`, which is exactly the discriminator the location-routed remedy needs.
- `detectCompletionAnomaly` (`model.go:1388`) carves out absent/unparseable `claimed_at` at :1394 with the comment "absent/unparseable claimed_at is other checks' territory" — the carve-out this probe was told to match.
- `parseTimestamp` returns `(time.Time, bool)`, so "parsed" and "absent/unparseable" are already one call.

*Exploration run inline by the orchestrator*

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/verify.go` (modify) — category constant, probe, finding constructor
- `skills/do-work-board/tools/queue-kanban/verify_test.go` (modify) — the captured RED as a permanent case, plus carve-out and routing cases
- `skills/do-work/actions/forensics.md` (modify) — Check 12's ordering condition and rewritten remedy; Check 14's probe table row

**Files I will NOT touch:**
- `model.go` — `detectCompletionAnomaly` keeps its single comparison; this probe is additive and does not restate it.
- `scripts/repair-req-timestamps.sh`, `scripts/audit-archive-timestamps.sh` — the predicate and the repair both already ship; the REQ's point is that only the read side was missing.

**Acceptance criteria (restated from the REQ):**
1. Ordering probe in `runVerifyProbes`, one finding per violated pair, naming REQ, both fields, both raw values.
2. Covers archived REQs, not only queue and working.
3. Remedy names `audit-archive-timestamps.sh` for archive and the SessionStart repairer for queue/working.
4. Not marked `[fixable]`.
5. Check 12 carries the ordering condition and its Suggested fix names the tools instead of hand git archaeology.
6. Check 14's probe table lists the new probe.
7. The fabrication heuristic is NOT built, and no fourth timestamp predicate is introduced.

## Decisions

- **D-01** (ESCALATE): The REQ names two pairs. Implementing exactly those left a hole my own carve-out fixture walked into: with `claimed_at` absent, nothing spans `created_at` and `completed_at`, so a REQ completed before it was created passes every comparison. Builder chose to add the outer pair, **guarded to fire only when `claimed_at` is absent or unparseable**. Reasoning: the rule the REQ cites is `created_at <= claimed_at <= completed_at`, and the outer relation is part of it — it is normally implied transitively, and the implication has nothing to travel through when the middle stamp is gone. The guard is what keeps it from being a third finding for a defect the two inner pairs already report. This is the REQ's own thesis applied to the REQ: catch the class, not the subset that happens to be spanned. Value: an impossible ordering cannot hide behind a missing `claimed_at`. Risk: one more comparison to keep in step with the shell predicate; mitigated by the comment naming the shell file and by `TestVerifyDoesNotDoubleReportTheOuterPair`, which pins the guard.
- **D-02** (DECIDE & STATE): Strict comparison (`Before`), so equal stamps are legal. Reasoning: the REQ's own Constraints say Step 2's claim and Step 3.6's estimate legitimately read the same instant — it measured 20 of 251 REQs here at 8%. A non-strict comparison would fire on all of them, which is the unearned-warning failure the REQ explicitly declined in its other half.
- **D-03** (DECIDE & STATE): The two declined non-goals were honored literally — no fabrication heuristic, and the Go predicate is commented as a restatement of the shell one with the two boundary decisions (strictness, absent-stamp meaning) named so a future editor changes both. Reasoning: the REQ asked for exactly this and named drift on those two boundaries as the risk.

## Implementation Summary

**What was done:** Added a timestamp-ordering probe to `queue-kanban verify` covering queue, working and archive, reporting one non-fixable finding per violated pair with both field names, both raw values, and a remedy routed by where the file lives. Rewrote forensics Check 12 to carry the ordering condition alongside its future-stamp condition and to point both halves at the repair scripts that already ship instead of prescribing hand git archaeology, and added the probe to Check 14's table.

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/verify.go` (modified) — `verifyCategoryTimestampOrdering`; `appendTimestampOrderingFindings` wired next to the completion-anomaly probe; `timestampOrderingFinding` builds the location-routed remedy.
- `skills/do-work-board/tools/queue-kanban/verify_test.go` (modified) — six cases: the captured RED, both-pairs-separately, queue remedy routing, the equal/absent/unparseable carve-out, the outer pair, and the no-double-report guard.
- `skills/do-work/actions/forensics.md` (modified) — Check 12 retitled and extended; its Suggested fix rewritten; one row added to Check 14's probe table.

**Tests touched:** six new cases in `verify_test.go`. No existing assertion changed meaning.

## Qualification

Passed — 3 files verified, 7 acceptance criteria traced, P-A-U confirmed.

- **[UNIFY] audit:** `gofmt -l .` clean, `go vet ./...` clean, `go build` and `go test ./...` pass. Grepped the diff for debug artifacts — none; the only output paths are `VerifyFinding` fields, which are the check's contract output. `maintainer-verify.sh` exits 0, which includes the shipped-reference contract over the edited `forensics.md`.
  One repair during UNIFY: an unquoted heredoc had shell-substituted `` `do-work cleanup` `` out of a Go comment, leaving `Not Fixable —  does not rewrite stamps`. Caught by re-reading the written block rather than trusting the build, since a mangled comment still compiles. Restored and the whole inserted block re-scanned for other substitution damage.
- **Substantive:** the probe is 20 lines of real comparison plus a constructor; the tests assert on content, not on absence of panic.
- **Requirements traced:** AC1–AC4 → `appendTimestampOrderingFindings` + `timestampOrderingFinding`, each pinned by a test; AC5–AC6 → `forensics.md`; AC7 → no heuristic added, and the Go predicate is commented as a restatement with its boundaries named.
- **Flowing:** the probe reads real frontmatter and changes verify's exit code — proven by the fixture run, which flips 0 → 1.

## Testing

- `go test ./...` in `skills/do-work-board/tools/queue-kanban` — ok, 39.5s.
- `gofmt -l .` — clean. `go vet ./...` — clean.
- `bash _dev/tests/maintainer-verify.sh` — exit 0. Baseline green before implementation.

**Red-green validation** — traced to the REQ's `## Red-Green Proof`, whose RED was a fixture rather than a test:

| | Before | After |
|---|---|---|
| Captured fixture (`created_at` 2026-08-19T12:00Z, `claimed_at` 2026-08-18T09:00Z, `completed_at` 2026-08-19T14:00Z, archived) | `OK: no findings`, **exit 0** | finding naming REQ-800, both fields, both values, archive remedy — **exit 1** |
| Check 12 states the ordering condition | absent | present, and names `audit-archive-timestamps.sh` |

The RED was run on the untouched tree first and reproduced the captured claim exactly. The REQ's second GREEN assertion — that the doc and the probe cannot drift apart — is satisfied by both carrying the condition; note it is asserted by inspection here, not by a test, because nothing in the Go package can reach a file in the core package. That limit is pre-existing and is the same one `forensics.md` Check 12 already records for its `futureStampCauseClause` mirror.

**False-positive check against real data:** ran the built binary against this repo's own tree — 250+ REQs across queue, working and archive — and the ordering probe produced **zero** findings. The pre-existing `completion-anomaly` findings there are unrelated: this is a shallow clone, so old commit hashes are not datable in it.

**Fixture mutation testing:**

| Reverted behavior | Result |
|---|---|
| `created_at`/`claimed_at` comparison disabled | `TestVerifyFlagsCreatedAfterClaimed` and `...EachViolatedTimestampPairSeparately` both FAIL |
| finding marked `Fixable: true` | `TestVerifyFlagsCreatedAfterClaimed` FAILs on the fixable assertion |

## Review

**Overall: 93%** — Acceptance: Pass

### Requirements Check

| Requirement | Status |
|---|---|
| Ordering probe in `runVerifyProbes`, one finding per violated pair, both fields and both raw values | ✅ Pinned by `TestVerifyFlagsCreatedAfterClaimed` and `...EachViolatedTimestampPairSeparately` |
| Covers archived REQs, not only queue/working | ✅ Iterates `board.AllRequests`; the RED fixture is an archived REQ |
| Remedy names the archive script and the SessionStart repairer, routed by location | ✅ Both directions pinned by tests |
| Not `[fixable]` | ✅ Asserted directly |
| Check 12 carries the ordering condition; remedy names the tools | ✅ |
| Check 14 probe table lists the new probe | ✅ |
| No fabrication heuristic; no fourth predicate | ✅ Nothing resembling the equality/shared-instant warning was added |

### Findings

**Important — none.**

**Minor:**

- **M1:** The doc-and-probe agreement the REQ asked to be inseparable is held by inspection, not by a test. Nothing in the Go package can read a core-package file, and `forensics.md` Check 12 already documents that exact limitation for its `futureStampCauseClause` mirror ("this one is in a different skill package and nothing can reach it, so it is kept in step by hand"). So this REQ inherits a known, documented gap rather than creating one. Not queued: closing it means a cross-package assertion mechanism, which is a real design decision and not this REQ's.
- **M2:** The outer-pair check (D-01) is beyond the REQ's literal two-pair instruction. Recorded as an escalated decision with its guard and its test rather than slipped in silently.

**Nit:**

- **N1:** `appendTimestampOrderingFindings` now has three near-identical `if` blocks. Readable as is, and collapsing them into a table of pairs would obscure the guard on the third, which is the one thing a future reader must not miss.

### Restatement Sweep

Redefined element: nothing. This REQ **adds** a probe and rewrites one remedy; it changes no existing contract token, field meaning, or command output shape. The one thing it comes close to redefining is the ordering predicate, and it deliberately does not — the Go comment states it is a restatement of `scripts/repair-req-timestamps.sh`'s rule with the two boundary decisions named.

Swept anyway, because the remedy text was rewritten: grepped for other places prescribing hand recovery of a timestamp from git history. `actions/capture.md:67` names `audit-archive-timestamps.sh` as the immutability-rule exception and is consistent with the new remedy. Check 1's `claimed_at` guidance is about staleness, not repair, and is unaffected. No stale restatement.

### Acceptance Testing

Ran the built binary three ways: against the REQ's captured fixture (exit 1, correct finding and remedy), against a queue-located violation (remedy correctly routes to the SessionStart repairer, and explicitly does not name the archive tool), and against this repository's real 250+ REQ tree (zero ordering findings — the false-positive check that matters most for a probe added to a gate everyone runs).

### Scores (on the record — not the headline)

| Dimension | Score |
|---|---|
| Requirements | 100% |
| Code Quality | 90% |
| Test Adequacy | 95% |
| Scope Discipline | 90% |
| Risk | Low |
| Acceptance | Pass |

Scope Discipline 90% for the outer pair (D-01), which is a deliberate, guarded, tested extension beyond the literal instruction. Risk Low rather than None: this adds a failing condition to a gate, and its false-positive mode would be noisy — measured at zero against real data, with the equal-stamp carve-out (the 8% case the REQ identified) pinned by test.

### Follow-up REQs Created

None. M1 is a pre-existing documented gap, M2 is a recorded decision, N1 is a preference.

## Lessons Learned

**What worked:** Running the captured RED before writing code, again. It confirmed the fixture reproduces the blind spot and — more usefully — the *shape* of the REQ's proof was already right, so the test could be the fixture rather than an invention. Second: running the finished probe against the repository's own 250+ REQs. A new failing condition on a gate everyone runs is only safe if you have measured its false-positive rate on real data, and "zero on 250" is the sentence that makes it safe to ship.

**What didn't:** Writing the Go source through an **unquoted** shell heredoc. Backticks inside a Go comment were command-substituted away, leaving `Not Fixable —  does not rewrite stamps` — which compiles, passes vet, passes gofmt, and passes every test. Only re-reading the written block caught it. For any generated file with backticks or `$` in it, quote the heredoc delimiter (`<<'PY'`); a build is not evidence that a comment survived.

Also: the REQ's two named pairs were not the whole rule. Implementing them exactly left `created_at > completed_at` with an absent `claimed_at` passing — and my own carve-out fixture was the thing that walked into it. The instance list was again narrower than the class, which is the third time this session.

**Worth knowing:** The ordering predicate now exists in two languages — `scripts/repair-req-timestamps.sh` (repair) and `verify.go` (read). Nothing holds them together mechanically; the Go comment names the shell file and its two boundary decisions (strict comparison, absent stamp is other checks' territory) precisely because that is the seam most likely to drift. If a third spelling is ever proposed, that is the moment to build the shared-fixture harness instead.

## Orientation

`queue-kanban verify` now catches a REQ whose stamps could not describe a real sequence of events, across the archive as well as the live queue — the place where nothing auto-repairs and where a fabricated `claimed_at` was previously invisible unless it happened to invert the one pair the board already compared. Forensics Check 12 stopped telling readers to reconstruct instants from git by hand and now points at the two repair scripts that already ship. Lives in the board's verify subsystem (`skills/do-work-board/tools/queue-kanban/verify.go`) with its documented twin in `skills/do-work/actions/forensics.md`.

**[MAP CHANGED]** — a new finding category, `timestamp-ordering`, joins the verify probe set. Anything that parses verify output by category, or counts probes against Check 14's table, sees one more.

Prime staleness spot-check: `_dev/primes/prime-kanban-board.md` and `_dev/primes/prime-action-files.md` — referenced paths still resolve; neither is made stale by this change, though the board prime's probe-set discussion now describes one fewer probe than ships, which Check 14's table (updated here) is the authoritative list for.
