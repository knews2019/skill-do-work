---
id: REQ-087
title: The board and verify hand the user the POSIX-only timestamp command the rule just stopped prescribing
status: completed
claimed_at: 2026-08-04T00:04:10Z
completed_at: 2026-08-04T00:11:37Z
commit: 5cfe1b5
kb_status: pending
created_at: 2026-08-03T21:45:43Z
user_request: UR-015
domain: general
prime_files: [tools/queue-kanban/prime-do-kanban.md]
tdd: false
depends_on: [REQ-078]
maintenance: false
addendum_to: REQ-078
review_generated: true
write_set: [tools/queue-kanban/verify.go, tools/queue-kanban/web/board.js, tools/queue-kanban/model.go, tools/queue-kanban/future_timestamp_test.go]
---

# The board and verify hand the user the POSIX-only timestamp command the rule just stopped prescribing

## What

REQ-078 made `actions/work-reference.md`'s Timestamp rule the only place in `actions/` that spells a
command for obtaining a stamp, and gave that rule a Windows form that actually runs. Three sites in
`tools/queue-kanban/` still hand a user the bare POSIX command:

- `tools/queue-kanban/verify.go:287` — the `Remedy:` string on a future-dated-timestamp finding:
  "re-stamp it with `date -u +%Y-%m-%dT%H:%M:%SZ` (the Timestamp rule in actions/work-reference.md)".
- `tools/queue-kanban/web/board.js:154` — the claim-stopwatch tooltip.
- `tools/queue-kanban/web/board.js:553` — the future-stamp data-warning text.

A Windows user who follows any of the three gets a command that does not exist on their box — the
exact failure REQ-078 fixed one layer up.

## Why

These were left out of REQ-078 on purpose, not by oversight: they are a **different surface with a
different tradeoff**. An action file can say "see the Timestamp rule" because the agent reading it can
open the rule. A board tooltip cannot — the person reading it is looking at a web page, and replacing
a usable command with a file reference makes the tooltip worse. So the fix here is a judgment call
about what a UI should say, not a continuation of REQ-078's sweep.

Low severity: nothing is corrupted, and the finding these strings decorate is itself correct. The
cost is a Windows user copying a dead command out of a UI that looks authoritative.

## Detailed Requirements

1. **Decide what each surface should say, per surface** — the three sites do not have to match. A
   plausible split: `verify.go`'s remedy keeps a command but adds the Windows one (it is CLI output,
   read next to a shell); the two `board.js` strings drop to the shape ("the current UTC instant,
   `YYYY-MM-DDTHH:MM:SSZ`") and cite the rule, since a tooltip's job is to explain the badge, not to
   be a manual. Argue for whatever you pick.
2. **Do not paste the Windows one-liner into all three.** Three copies of a two-branch command in
   display strings is the inline-copy problem REQ-078 just removed, relocated.
3. **If a command survives anywhere, it must agree with the rule** — `ToUniversalTime`, the `\T`/`\Z`
   escapes, the `powershell -NoProfile -Command` wrapper for `cmd`.
4. **Consider whether a contract assertion is warranted** and say why if not. REQ-078's assertion is
   scoped to `actions/` deliberately; widening it to `tools/` would flag `timestamp.go`'s rationale
   comments and the test fixture, which are correct as they are.

## Constraints

- `tools/queue-kanban/` changes are folded into the skill's own version and changelog — no
  independent versioning (the tool's conventions are in `tools/queue-kanban/prime-do-kanban.md`).
- No behaviour change. These are display strings; the findings and badges they annotate stay exactly
  as they are.

## Dependencies

`depends_on: [REQ-078]` — the Windows form these strings would have to agree with ships there.

## Builder Guidance

**Certainty: Firm on the inventory, open on the wording.** All three sites were found by a three-shape
grep during REQ-078's review and are listed above with line numbers. What each should say is the
actual work.

## Triage

**Route: A** - Simple

**Reasoning:** Firm inventory with line numbers, no location discovery needed; the work is a wording
judgment per surface. Four display strings, one repointed assertion.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Scope

**Files I will touch:**
- `tools/queue-kanban/verify.go` (modify) — the future-dated-`claimed_at` remedy string
- `tools/queue-kanban/web/board.js` (modify) — the clock-skew tooltip and the future-stamp badge title
- `tools/queue-kanban/model.go` (modify) — **added during the build (D-02)**: a fourth site of the
  same class, the server-side future-timestamp board warning
- `tools/queue-kanban/future_timestamp_test.go` (modify) — **added during the build (D-03)**: its
  assertion pins the literal the `model.go` change replaces

**Files I will NOT touch:**
- `tools/queue-kanban/timestamp.go` — its `date -u` mentions are rationale prose explaining why the
  subcommand exists, and are correct as they stand
- `actions/work-reference.md` — the Timestamp rule is the source these strings cite; unchanged

**Acceptance criteria (restated from REQ):**
- [ ] Each surface's wording decided per surface, with an argument (req 1)
- [ ] The Windows one-liner is not pasted into multiple display strings (req 2)
- [ ] Any surviving command agrees with the rule (req 3)
- [ ] The contract-assertion question answered, with reasoning if declined (req 4)
- [ ] No behaviour change — findings and badges unchanged

## Implementation Summary

**Files changed:**
- `tools/queue-kanban/verify.go` (modified)
- `tools/queue-kanban/web/board.js` (modified)
- `tools/queue-kanban/model.go` (modified)
- `tools/queue-kanban/future_timestamp_test.go` (modified)

**What was done:** Built in a git worktree as Builder B of REQ-085's live fan-out acceptance test —
branch `worktree-agent-REQ-087-posix-only-timestamp-command`, commit `202ff3e`, integrated by the
`--no-ff` merge `5cfe1b5` (range `3ccbf36..5cfe1b5`). `verify.go`'s remedy keeps a command but swaps
the POSIX-only `date -u +…` for `queue-kanban now` — the Timestamp rule's own option 1, which prints
the right shape on every platform and whose "only if already built" precondition is satisfied by the
fact that the reader is looking at that binary's output. The three warning/tooltip strings in
`board.js` (×2) and `model.go` drop to the target shape plus a citation of the rule. The
`future_timestamp_test.go` assertion that pinned the old literal was repointed.

## Decisions

- **D-01 (DECIDE & STATE)** — *`verify.go` keeps a command, and the command is `queue-kanban now`.*
  The REQ's suggested split had `verify.go` "keep a command but add the Windows one". Adding the
  two-branch PowerShell form to CLI output would have been three lines of quoting in a remedy string,
  and requirement 2 warns against relocating the inline-copy problem. `queue-kanban now` is strictly
  better: it is the rule's *first-choice* source, it is one token, and it is platform-independent —
  which is the actual defect. The precondition that normally qualifies option 1 ("only if the binary
  is already built") cannot fail here, because the string is printed *by that binary*.
- **D-02 (DECIDE & STATE, scope extension)** — *`tools/queue-kanban/model.go` added to scope.* The
  REQ's inventory named three sites; a grep across `tools/` during the build found a fourth of exactly
  the same class — the server-side future-timestamp board warning at `model.go:321`, a user-facing
  string carrying the same POSIX-only command. Fixing three of four would have left the defect live in
  the warning that reaches the same user as the badge. Reported by the builder as an out-of-scope find
  and recorded here by the orchestrator, with `## Scope` and `write_set` extended before the write —
  not after.
- **D-03 (DECIDE & STATE, scope extension)** — *`tools/queue-kanban/future_timestamp_test.go` added to
  scope.* `TestFutureTimestampWarningNamesFieldAndFix` asserts the literal `date -u +%Y-%m-%dT%H:%M:%SZ`
  appears in the `model.go` warning, so D-02 breaks it by construction. The assertion's *intent* — the
  warning names the field and the fix — is preserved exactly; only the fix's spelling changed, so it
  now pins the target shape and the rule citation. Cross-REQ impact recorded under `## Testing`.
- **D-04 (DECIDE & STATE)** — *Requirement 4: no new contract assertion.* A ratchet here would have to
  distinguish a **user-facing display string** from a **rationale comment or test fixture**, which no
  grep can see. After this change the tool's remaining `date -u` mentions are `timestamp.go:13,21`
  (rationale prose), `future_timestamp_test.go` (a fixture), and `verify.go`'s new comment explaining
  why the POSIX floor is deliberately absent there — every one correct, and every one a file-scoped
  assertion would flag. Hand-listing the files to exclude just relocates the failure: that list goes
  stale the day a new display surface appears, which is the closed-enumeration pattern the maintainer
  doc names explicitly. The control that fits is the per-surface judgment requirement 1 asked for.

## Qualification

Passed — 4 files verified, 5 requirements traced.

- `tools/checks/qualify.sh` run with `DO_WORK_DIFF_RANGE="3ccbf36..5cfe1b5"` (worktree dispatch mode).
- **Requirements traced:** 1 → four surfaces, two different treatments, argued in D-01. 2 → the Windows
  one-liner appears in **no** display string; `grep -rn "ToUniversalTime" tools/` is empty. 3 →
  satisfied by the one surviving command being the rule's own option 1. 4 → D-04. No behaviour change →
  the findings and badges are untouched; only string contents differ, and the full Go suite passes.
- **Post-merge verification:** re-run against the merged main tree, not the builder branch — this is
  the case that matters here, because `verify.go` was also changed by REQ-083 and REQ-084 earlier in
  the same session. `gofmt`, `go vet` and `go test ./...` are all green *after* the merge, so the three
  REQs' edits to one file compose.

## Testing

**Tests run:** `cd tools/queue-kanban && go test ./...` (+ `go vet`, `gofmt`), then re-run post-merge
**Result:** ✓ All passing

Non-behavioral change (display strings), so red-green validation does not apply — the proof is the
repointed assertion plus the full suite.

**Existing tests updated (cross-REQ impact):** one, intentional.
`TestFutureTimestampWarningNamesFieldAndFix` (`future_timestamp_test.go`, from the REQ that introduced
the future-timestamp board warning) asserted the literal `date -u +%Y-%m-%dT%H:%M:%SZ` inside that
warning — precisely the behavior this REQ changes. It was **repointed, not deleted**: it now asserts
the warning contains `YYYY-MM-DDTHH:MM:SSZ` *and* `Timestamp rule in actions/work-reference.md`, so the
original intent (the warning names the field and the fix) is pinned at least as tightly as before. A
comment at the assertion records why it moved. Caught by the suite, not by inspection — the failure is
what surfaced the fourth site's blast radius.

*Verified by work action*

## Lessons Learned

**What worked:** Re-running the grep across `tools/` instead of trusting the REQ's three-site list
found the fourth site. That is now three sweep REQs in a row where the inventory was a floor.

**What didn't:** Nothing failed, but the first instinct — follow the REQ's suggested split literally
and paste the Windows one-liner into `verify.go` — would have satisfied the requirement while
relocating the exact problem requirement 2 warns about. The better answer was already in the rule:
option 1 exists precisely so callers do not have to branch on platform.

**Worth knowing:** `queue-kanban now` is the right timestamp instruction to put in **the tool's own
output**, and only there. The rule's "only if the binary is already built" caveat is what normally
makes option 1 unsafe to recommend blindly — and it cannot fail in a string printed by that binary.
Anywhere else, cite the rule instead.

## Orientation

**Now a Windows user reading a clock-skew tooltip, a future-stamp badge, a board data warning, or a
`queue-kanban verify` remedy gets an instruction they can actually act on** — the three display strings
name the target shape and cite the rule, and verify's remedy points at `queue-kanban now` instead of a
POSIX-only command. Lives in the board tool's display layer (`tools/queue-kanban/prime-do-kanban.md`).

Leaf change — no new module, contract, or data flow; the findings and badges these strings annotate are
byte-identical in behavior. Prime staleness spot-check: `prime-do-kanban.md`'s referenced paths all
still exist; not made stale.

## Full Context

Found by `actions/review-work.md`'s Restatement Sweep during REQ-078's review — see that REQ's
`## Review` → Findings.
