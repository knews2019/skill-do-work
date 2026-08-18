---
id: REQ-261
title: Delete the date-only tripwire and keep the rule
status: claimed
claimed_at: 2026-08-18T22:59:48Z
route: A
created_at: 2026-08-18T19:30:47Z
status_changed_at: 2026-08-18T21:01:24Z
user_request: UR-055
addendum_to: REQ-253
domain: general
effort_estimate: trivial
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: true
write_set:
- skills/do-work/actions/work-reference.md
estimate:
  p50_active_minutes: 15
  confidence: medium
  calculated_at: 2026-08-18T22:59:48Z
  basis:
    - Route A
    - 1-file write set
    - 2 acceptance criteria
    - full-suite verification
---

# Delete the Date-Only Tripwire and Keep the Rule

## What

The Timestamp rule's date-only paragraph ends with "revisit if a second consumer appears". Remove that clause. Keep everything else in the sentence — the shell one-liners, and the reason there is no tool subcommand (adding one would widen the skill's single compiled-dependency exception for something the POSIX floor already covers).

The clause is the only part that does not survive its own argument: it keys on how many consumers exist, and consumer count does not bear on whether a shell one-liner suffices. Leaving it invites a re-litigation the surrounding sentence already settles — a list where a condition belongs (CLAUDE.md → State conditions, not lists).

## Requirements

- The "revisit if a second consumer appears" clause is gone; the rest of the date-only paragraph is unchanged in meaning.
- No date-only subcommand is added to the board tool.
- The paragraph still reads as one coherent sentence after the removal — check the ui-review consumer clause REQ-253 added still sits naturally beside it.
- `bash _dev/tests/maintainer-verify.sh` exits 0 (the Timestamp-rule citation contract counts 54 instant / 17 date-only sites today; a prose-only removal must not move them).

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `crew-members/maintenance.md` (this REQ is `maintenance: true`, so removal is the assigned fix) and the whole Timestamp rule section of `work-reference.md` before cutting, to see what the sentence actually claims. Two greps first, to check whether the tripwire is cited or restated anywhere else: `grep -rn "second consumer" .` and `grep -rn "date-only\|date only" --exclude-dir=.git .` — the tripwire occurs exactly once in the repo; every other `date-only` hit cites *the rule*, not the conditional, so nothing dangles. Then `grep -rn "timestamp_rule_block\|single consumer\|date -u +%F\|revisit if" _dev/` to check for a test pinning the sentence: `contract-regressions.sh:1506-1517` pins only `ToUniversalTime` and `powershell -NoProfile -Command`, both untouched.
- [x] **[APPLY]:** One `python3` in-place replacement with a `count(old) == 1` assertion, so a silent multi-site or zero-site edit was impossible. Scope stayed at the one declared file.
- [x] **[UNIFY]:** `git diff --stat` → one file, 1 insertion / 1 deletion. Full `-U0` diff read by eye: the only textual delta is the two removed fragments — no reflow, no trailing whitespace, no stray characters, no debug artifacts. `git status --short` empty after committing. The canonical gate runs the repo's linters and passed. File verified: `skills/do-work/actions/work-reference.md` — checked that the date-only paragraph is still one coherent sentence, that the two named consumers and REQ-253's deliberate-UTC clause survive unchanged, and that the no-subcommand **reason** still stands on its own without the count.

## Context

Discovered by REQ-253's builder ([low]; the tripwire sentence was left verbatim as a deliberate maintainer call). Note the board's pinned write-surface count is unaffected either way — `now`-style output is read-only.

## Open Questions

- [ ] I discovered this out-of-scope task while working on REQ-253: the date-only paragraph's "revisit if a second consumer appears" condition is now true. Should I process this as a new task? → Yes, and the answer is decided: delete the tripwire clause and keep the rule
  Recommended: Yes, add to queue (will flip to 'pending') — and the builder decides between the subcommand and a re-stated threshold.
  Also: No, discard it — two consumers on the shell one-liner is still fine.

**Answered [2026-08-18]:** User approved via `do-work clarify` **and settled the underlying question**, after asking where the sentence came from. Provenance established during clarify: the clause arrived with the repository import (recorded author "Claude", root commit `8d5c2ab`) from a pass that restructured the Timestamp rule — it is builder prose, not a maintainer decision. The user's reasoning, which is now the REQ's requirement: the rule itself is sound (`date -u +%F` works on the POSIX floor, and the board tool is the skill's only sanctioned compiled dependency, so putting a date behind a Go binary would widen that exception for nothing), but the tripwire keys on **consumer count**, which has no bearing on that argument — a shell one-liner is no worse at two callers than at one. Delete the clause; keep the rule. Do not add a date-only subcommand.

---

## Implementation Summary

**Files changed:**
- `skills/do-work/actions/work-reference.md` (modified)

**What was done:** Deleted the date-only tripwire from the Timestamp rule's date-only paragraph — the clause "revisit if a second consumer appears" — **and, in the same sentence, the `for a single consumer` premise it rested on.** That premise was already false on its own line: the sentence names two consumers (the `memory/logs/` mirror and REQ-253's ui-review report header), so a builder cutting only the trailing clause would have shipped a sentence that names two consumers and then argues from one. The rule survives intact and now stands on its reason rather than on a count: there is no date-only subcommand because adding one would spend the skill's narrow compiled-tooling exception on something `date -u +%F` already covers. No subcommand added; no Go source, test, or board-tool file touched.

---

## Discovered Tasks

Transcribed by the orchestrator from `do-work/runs/work-2026-08-18-230100/REQ-261-handback.md` (a worktree builder cannot write this file — REQ-270).

- **None to queue, and the class was checked rather than assumed.** The builder grepped shipped instructions and docs for the same *shape* — a conditional "revisit/reconsider if…" tripwire sitting in standing prose — with `grep -rniE "revisit (this|if|when)|reconsider (if|when|this)|if a (second|third)|when a second" skills/ _dev/primes/ decisions/ CLAUDE.md`. Nothing else in `skills/` or `_dev/primes/` has this shape. The three remaining hits are in `decisions/records/` (adr-013, adr-016, adr-017), where "Revisit if X" sits in an Alternatives or Consequences section — a decision record stating its own boundary, which is the correct home for a conditional. Left alone deliberately; no follow-up recommended.

---

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh` against the merged tree (range `2fc89fe..210abba`), un-piped with the exit code read directly
**Result:** ✓ `GATE_EXIT=0` — "Maintainer verification passed." This run is both Step 6.5's testing and Step 8's post-merge verification.

**The acceptance criterion that actually matters here, verified:** the gate's Timestamp rule citation contract reports `54 instant write sites cited, 17 date-only sites recognized` — **unmoved**. The REQ required the tripwire to go without disturbing which sites the rule governs, and that contract is the mechanical statement of exactly that. A deletion that had clipped the paragraph's governing scope would have moved one of those counts.

**Red-green validation:** none, and none is owed. This is a deletion of standing prose with no behavioural surface: no subcommand was added or removed, no Go source, test, or board-tool file was touched, and the `timestamp_rule_block` assertions in `_dev/tests/contract-regressions.sh` pin `ToUniversalTime` and `powershell -NoProfile -Command`, neither of which this diff goes near. The regression evidence is the citation contract above plus the full gate.

**New tests added:** none. A test asserting "this sentence does not contain the words 'revisit if'" would pin the instance rather than anything worth keeping, which is the failure shape this session has been correcting all evening.

**Existing tests updated (cross-REQ impact):** none.

*Verified by work action*

