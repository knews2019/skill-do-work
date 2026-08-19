---
id: REQ-280
title: Probe timestamp ordering, and point Check 12 at the archive repair that already ships
status: pending
created_at: 2026-08-19T13:42:45Z
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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
