---
id: REQ-254
title: Let qualify tell a check's own output from leftover instrumentation
status: claimed
created_at: 2026-08-18T14:04:16Z
claimed_at: 2026-08-18T19:12:47Z
route: B
status_changed_at: 2026-08-18T14:04:16Z
user_request: UR-055
addendum_to: REQ-244
domain: general
review_generated: true
effort_estimate: trivial
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: true
write_set:
- skills/do-work/tools/checks/qualify.sh
- _dev/tests/prescribed-shell-scripts-behavior.sh
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-08-18T19:13:28Z
  basis:
    - trivial short-circuit (effort_estimate: trivial)
---

# Let Qualify Tell a Check's Own Output From Leftover Instrumentation

## What

`qualify.sh`'s debug-artifact scan FAILs on any added `print(` line, which makes it fire on a **check's own success output** — the reporting a checker is supposed to have. It fired exactly that way on REQ-244's remediation and had to be overridden on the record.

## Context

REQ-244's review found that the hand-back's GREEN transcript came from a prototype rather than the shipped checker, because the shipped checker **printed nothing on success and so had no transcript to quote**. The fix was to give it one. That fix then tripped `qualify.sh`:

```
FAIL: [UNIFY] is checked but the diff adds debug artifacts — un-check it and flag:
  155:+print(
```

The flagged line is the check's only success line, placed after its failure raise, and `maintainer-verify.sh` prints it on every run. So the gate that exists to catch stray instrumentation fired on the fix for a defect caused by *missing* instrumentation.

A false positive here is not free. It teaches a reader that a qualify FAIL can be waved through, which is the opposite of what the gate is for — and the next override may be a real artifact.

## Requirements

- A checker's own reporting does not trip the debug-artifact scan, while genuine leftover instrumentation still does.
- The distinguishing rule is **stated as a condition** and not as a list of allowed files or allowed line shapes (CLAUDE.md → Closed Enumerations Go Stale; the detail is in `_dev/primes/prime-shell-commands.md`).
- `maintenance: true`: **ask what can be removed before adding.** It is worth asking whether a bare-`print(` grep earns its place at all, given that the scan cannot see the difference between a debug line and a contract's output, and given that the reviewer and the orchestrator both read the diff anyway. Deleting a heuristic that cries wolf may be a better answer than teaching it a new trick.

## Implementation Summary

**What was done:** `qualify.sh`'s debug-artifact scan (Check 4) now splits its tokens by property. Unfinished-work markers (`debugger`, `TODO`, `FIXME` — vocabulary illustrative) FAIL anywhere, unchanged. Output primitives (`print(`, `console.log` — the class, not the fired token) are judged by process-exit ownership: a file that ends its own process (exit idioms, illustrative vocabulary) has a terminal audience, so an added output line is presumptively its own reporting and surfaces as a legible WARN; a file that never ends its process is library code, so the same line FAILs naming the file and reason. The output half walks changed files per path in both serial and range modes. REQ-244's ready-made case now passes with a WARN; genuine library instrumentation still FAILs. Deletion of the heuristic was weighed and declined on the REQ's own GREEN criterion (D-01). Orchestrator applied one seam: `review-work.md`'s diff-hygiene line no longer re-flags a checker's success output.

**Files changed:**
- `skills/do-work/tools/checks/qualify.sh` (modified) — token split (`unfinished_marker_regex` / `output_primitive_regex`), per-file ownership walk, condition stated in the comment
- `_dev/tests/prescribed-shell-scripts-behavior.sh` (modified) — three `qualify:` lock-ins (reporter passes serial+range; library instrumentation FAILs with named file; TODO in a reporter still FAILs), suite 52 → 55 cases
- `skills/do-work/actions/review-work.md` (modified) — integration seam: diff-hygiene wording aligned with the new boundary

*Integrated by orchestrator from builder hand-back; merge range `8f564b3..116eec6`.*

## Decisions

Transcribed from the builder hand-back:

- **D-01 (DECIDE & STATE):** keep a narrowed scan, decline deletion — the REQ's GREEN requires library instrumentation to keep FAILing; the narrowing removes exactly the crying-wolf half and the scan now claims strictly less.
- **D-02 (DECIDE & STATE):** the condition is process-exit ownership — printed output belongs to whoever owns the process exit; both token and exit-idiom vocabularies marked illustrative.
- **D-03 (DECIDE & STATE):** the exemption covers the output-primitive class, not the fired token — `console.log` in a Node checker is the same defect one language over; pinned in the reporter fixture.
- **D-04 (DECIDE & STATE):** unfinished-work markers stay unconditional — the reporter exemption is not a file-level pardon; pinned by the third lock-in.
- **D-05 (DECIDE & STATE):** reporter hits WARN, not silent pass — WARN is the script's existing judgment channel; accepted limit: a forgotten debug print inside a checker WARNs, the honest boundary of intent-blindness.
- **D-06 (DECIDE & STATE):** no `work.md` seam — the FAIL/WARN prose contract is unchanged; only which side one token class lands on moved.

## Qualification

Passed — 3 files in merge range `8f564b3..116eec6` (2 builder + 1 seam), requirements traced (condition-not-list verified in the shipped comment; both REQ RED cases exercised; deletion genuinely weighed with the tie-breaker named), P-A-U audited. The builder's pushback — the REQ's condition-and-deletion asks pull opposite ways and the Red-Green Proof broke the tie — is accepted as the correct reading; the widest consequence (any script's stray print WARNs rather than FAILs) is recorded with its rationale.

## Builder Guidance

Candidate signals, none prescribed: whether the added line is inside `_dev/tests/`; whether the file is itself a check; whether the print is unreachable after a raise/exit on the failure path; whether the printed text is asserted anywhere. The mechanical definition is the builder's call — the requirement is that it keys on a property, not on a filename.

## Red-Green Proof

**RED prompt/case:** run `qualify.sh` over a diff that adds a success-reporting `print(` to a `_dev/tests/` checker, and over a diff that adds a genuine debug `print(` to implementation code.
**Why RED now:** both FAIL identically today; REQ-244's remediation range is a ready-made case.
**GREEN when:** the first passes, the second still FAILs, and the reason each gets is legible.
**Validation:** Orchestrator override recorded in REQ-244's Review Remediation section.

---

## Triage

**Route: B** - Medium

**Reasoning:** The defect and a ready-made red-green case are stated; the open design question — a condition that separates a check's own output from instrumentation, or deleting the heuristic — is the builder's, inside one script plus its lock-ins.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Scope

**Files I will touch:**
- `skills/do-work/tools/checks/qualify.sh` (modify) — the debug-artifact scan keys on a condition, or the bare-print heuristic is removed (maintenance latitude)
- `_dev/tests/prescribed-shell-scripts-behavior.sh` (modify) — red-green lock-ins: checker success line passes, genuine debug print still FAILs

**Files I will NOT touch:**
- `skills/do-work/actions/work.md` — Step 6.3's prose already describes FAIL/WARN semantics; change it only via integration seam if the resolution alters the contract.

**Acceptance criteria (restated from REQ):**
- [ ] A checker's own reporting does not trip the scan; genuine leftover instrumentation still does.
- [ ] The distinguishing rule is stated as a condition, never a file/line-shape list.
- [ ] The delete-the-heuristic answer was genuinely considered (maintenance: true).
- [ ] `bash _dev/tests/maintainer-verify.sh` exits 0.

## Pre-Flight

**Git:** ✓ clean
**Tests baseline:** ✓ `bash _dev/tests/maintainer-verify.sh` exits 0 at the branch point (0.212.14 tip)
**Dependencies:** ✓ Go 1.26.1, ShellCheck 0.11.0, `just`, Node, Chromium present

*Checked by work action*
