---
id: REQ-234
title: Stop the shell behavior suite counting its own cases
status: completed
completed_at: 2026-08-18T11:07:00Z
commit: 48ed251
claimed_at: 2026-08-18T10:56:00Z
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-08-18T10:56:00Z
  basis:
    - trivial short-circuit
route: A
status_changed_at: 2026-08-18T10:26:34Z
domain: general
created_at: 2026-08-18T01:44:18Z
user_request: UR-042
addendum_to: REQ-229
effort_estimate: trivial
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: false
maintenance: true
write_set:
- _dev/tests/prescribed-shell-scripts-behavior.sh
---

# Discovered Task: Stop the Shell Behavior Suite Counting Its Own Cases

## What

`_dev/tests/prescribed-shell-scripts-behavior.sh` ends by printing how many cases it ran, from a hand-maintained literal. The number matches no count derivable from the file, so it reports a remembered figure as a measured one. Derive it or drop it.

## Context

Found while adding two cases in REQ-229. The closing line read `(45 named script cases)` while the file carried 40 case-header comments; REQ-229 bumped it to 47 because adding two cases in the existing style makes that correct under whatever convention produced 45, but it did not repair the underlying problem.

`_dev/primes/prime-shell-commands.md` § *Closed Enumerations Go Stale* is this exact pattern, and a summary line reporting a suite's own size is a bad place for it: a reader has every reason to take that number as measured.

## Requirements

- The closing line either reports a count computed from the file at run time, or reports no count at all. It must not carry a literal.
- If a count is computed, the thing being counted has a stated definition in the file — the shape that makes something "one case" — so the number and the file cannot disagree.
- No test case is added, removed, or weakened by this change.
- `bash _dev/tests/maintainer-verify.sh` still exits 0.

## Red-Green Proof

**RED prompt/case:** the literal at the end of `_dev/tests/prescribed-shell-scripts-behavior.sh`.
**Why RED now:** it is a hand-maintained number that already disagrees with every count derivable from the file.
**GREEN when:** that literal is absent — either replaced by a computed count with its counting rule stated, or removed along with the claim.
**Validation:** Discovered during REQ-229; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**, deletion branch.

## Open Questions

- [x] I discovered this out-of-scope task while working on REQ-229: one of the maintainer test scripts finishes by announcing how many test cases it just ran, and that number is typed into the file by hand rather than counted. It is already out of step — it claimed 45 while the file held 40 of the obvious unit — so it reports something remembered as though it were measured. Nothing is broken: every test still runs and still passes, and the number appears only in a success message. The fix is either to count the cases at run time or to stop claiming a number. It is your call rather than mine because working out what the original number was counting means picking a definition of "one case" for this suite, and picking it wrong would silently change what that line reports to every future reader. Should I process this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.
  *Answered 2026-08-18 via `do-work clarify` — user approved deriving the count at run time (with a stated counting rule) or dropping the claim.*

---

## Triage

**Route: A** - Simple

**Reasoning:** One named file, one named line, and the REQ states both acceptable outcomes. The only open question was compute-vs-drop, which is a judgment call inside a single file rather than a discovery problem.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Implementation Summary

**Files changed:**
- `_dev/tests/prescribed-shell-scripts-behavior.sh` (modified)

**What was done:** The closing line's `(47 named script cases)` literal is gone. The number is now grepped out of the file at run time, and the shape being counted — one column-zero header comment of the form `# <script-name>: <what it proves>`, one per fixture block — is stated in a comment immediately above the grep, so the reported number and the file cannot disagree. One `suite_file` variable was added beside the existing `repo_root` derivation so the count reads the file by absolute path rather than depending on the caller's directory. The reported figure changes 47 → 42, deliberately: no rule over the file yielded 47.

## Qualification

Passed — 1 file verified in the merge range `53885e4..48ed251`, 4 acceptance criteria traced, no P-A-U section on this REQ.

Judgment checks, run against the merged tree rather than taken from the builder's report:
- **The counting rule does not self-match.** The explanatory comment it added begins "One case is one fixture block", and the pattern requires a lowercase leading character, so the comment describing the count is not counted by it. Verified: the grep returns 42 with the comment present.
- **No wrapped header is double-counted.** Several case headers wrap onto a second line; a continuation that happened to begin `# word: ` would inflate the count silently. Checked mechanically for matches on consecutive line numbers — none.
- **Every match is a real case header.** Listed all 42; each is `# <script-name>: <claim>` opening a fixture block.
- **Not cwd-sensitive.** Ran the suite from an unrelated directory by absolute path: `Prescribed shell script behavior probes passed (42 named script cases).`, exit 0.

## Testing

**Tests run:** `bash _dev/tests/prescribed-shell-scripts-behavior.sh` (from the repo root and from an unrelated cwd), then `bash _dev/tests/maintainer-verify.sh`
**Result:** ✓ Exit 0, run unpiped. The verify output carries the suite's new closing line at line 15: `Prescribed shell script behavior probes passed (42 named script cases).`

**Red-green validation:** this is the **deletion branch** of the Finding-Closure Ratchet — the named surface is the literal, and closure is that it is gone.
- RED, before: `printf 'Prescribed shell script behavior probes passed (47 named script cases).\n'`, with no count derivable from the file equal to 47 (header comments 42, column-zero comments 115, `fail_case` invocations 179, distinct "… case" phrases 46).
- GREEN, after: no literal on the closing line and no hardcoded `47` anywhere in the file; the number printed is the grep's output.

**New tests added:** none — the REQ forbids adding, removing, or weakening a case, and this changes only what the suite reports about itself.

*Verified by work action*

## Discovered Tasks

- [low] The suite names its assertion units two ways. Most read `<script> <name> case …`, but 23 distinct messages across 8 families (`ai-report` batch replays, `generate-report-image-batch` usage-error, `install-last30days` checks, `run-blocked-check` process-tree cleanup) say `replay`, `check`, `cleanup`, or the plural `cases`. That split vocabulary is why the original 45/47 was unrecoverable, and it is why the count now keys on the header comment rather than the assertion label. Renaming 23 assertion messages is test text the REQ forbids weakening, so it is recorded rather than done.

## Review

**Overall: 97%** | 2026-08-18T11:05:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | N/A |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded gate disposition):** None

**Minor findings:** 1 (report only)
- The count measures **header comments**, not executed blocks. Delete a fixture block and leave its header — or the reverse — and the number is wrong again, just wrong from a stated rule instead of from memory. That is a real residual, and it is the honest ceiling of a textual count; the alternative (counting at execution time) would mean threading a counter through 42 blocks, which is more machinery than the line is worth. Recorded so nobody reads the number as stronger than it is.

**Restatement sweep:** the diff changes what a suite reports about itself, so the sweep asked who else states that number. `CHANGELOG.md` carries it twice as dated release history (`41 → 44`, `27 → 29`) and two archived hand-backs carry it once each — all four are history, correctly left alone, and the builder explicitly declined to rewrite them. Nothing parses the closing line: `_dev/tests/staged-skills-contract.sh:183` invokes the suite and consumes only its exit status. No live restatement remains.

**Acceptance:** Pass — the literal is gone, the rule is stated in the file, the count is derived from it, and `maintainer-verify.sh` exits 0 with the new line visible in its output.

**Suggested testing:** 1 item
- Nothing pins the counting rule itself. If someone later adds a fixture block with a header that does not match the stated shape, the number quietly under-reports and no check notices. A one-line assertion that the count is non-zero and equals the number of `fail_case`-bearing blocks would close that, if it is judged worth the line.

**Follow-ups created:** None; **sweeps appended to:** None

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Testing the delete branch against `maintenance.md`'s own standard instead of applying delete-before-you-add reflexively. The repo's history showed the count reached for four separate times as evidence of suite growth — twice in hand-backs, twice in shipped release notes — which is a standing habit, not decoration. Deleting it would have removed something demonstrably wanted; computing it serves the same want and makes it true for the first time. The reflex was still honored, because what got deleted is the hand-maintained number and the obligation to remember it.

**What didn't:** Defining a case by its `fail_case` label. It looked like the natural unit and gave 46, temptingly close to the literal 47 — but 23 assertions across 8 families say `replay`, `check`, or `cleanup` instead of `case`, so that rule would silently disagree with a human count in exactly the way a derived number must not. The near-miss was the trap: a rule that almost reproduces the remembered figure is more dangerous than one that obviously does not, because it invites you to stop looking.

**Worth knowing:** No rule over this file yields 47, which independently confirms REQ-229's finding that the original convention is unrecoverable. Counts in hand-backs and changelog entries written before this commit are on that lost convention and are not comparable with counts after it; they were left unedited because they are history.

## Orientation

`_dev/tests/prescribed-shell-scripts-behavior.sh` now derives the case count it reports from its own text, with the counting rule stated beside the grep, instead of printing a number someone remembered. Lives in the maintainer test suite (`_dev/primes/prime-shell-commands.md`).

Not `[MAP CHANGED]` — one closing line, no contract and no structure altered. Staleness spot-check on `_dev/primes/prime-shell-commands.md`: every referenced path resolves, and its § *Closed Enumerations Go Stale* now has one fewer counter-example in the tree it describes. The prime is not stale.
