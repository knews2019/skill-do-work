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
