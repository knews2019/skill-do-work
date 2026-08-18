# REQ-254 builder brief

**Route B.** Estimated 5 active minutes (P50, effort_estimate: trivial — but take the time the condition design needs).

The defect: `qualify.sh`'s debug-artifact scan FAILs on any added `print(` line, so it fired on a checker's own success output (REQ-244's remediation — the override is on the record in that archived REQ). The requirement is a **condition, not a list**: a checker's own reporting passes, genuine leftover instrumentation still FAILs, and the rule keys on a property (candidate signals in the REQ body — none prescribed). Genuinely weigh deletion: the scan cannot see intent, reviewers read the diff anyway, and a heuristic that cries wolf teaches people to wave through FAILs — decide, and record the reasoning as a D-XX either way.

Live context: the behavior suite is at 52 named cases; your lock-ins add to `_dev/tests/prescribed-shell-scripts-behavior.sh`. REQ-244's ready-made red-green case is described in your REQ body's Context.

## How this build runs

You are a **worktree builder** dispatched by the do-work work pipeline. Everything binding is in this brief.

**Your tree, your branch.** Work only inside `/home/user/skill-do-work-worktrees/worktree-agent-REQ-254-let-qualify-tell-a-checks-own-output-from-instrumentation` — a full checkout on branch `worktree-agent-REQ-254-let-qualify-tell-a-checks-own-output-from-instrumentation`, cut from integration tip `eae572d` (version 0.212.14).

- Never write anything under `/home/user/skill-do-work` — the one exception is your hand-back file, named below.
- Never read or write `do-work/` in your own worktree (stale snapshot; your REQ body is inlined below).
- Commit on your own branch in small increments. Do not touch `VERSION`, `CHANGELOG.md`, or `skills/do-work/actions/version.md` — serial-only, integrator-owned.
- A needed one-line edit outside your write set is an *integration seam*: hand back the exact line and where it goes. Larger needs: stop and report.
- Out-of-scope finds go in `## Discovered Tasks` — never fixed inline.

**Crew rules** (read from your own worktree first): `skills/do-work/crew-members/general.md`, `coding-guardrails.md`, `communication-style.md`. This REQ is `tdd: true` and `maintenance: true` — also read `crew-members/testing.md` and `crew-members/maintenance.md`; the delete-the-heuristic answer deserves genuine consideration. Read every `prime_files` path too.

**P-A-U phasing is mandatory** — work the [PLAN]/[APPLY]/[UNIFY] block; record evidence in your hand-back (the orchestrator transcribes and audits it against the diff). Log significant choices as D-XX (DECIDE & STATE vs ESCALATE with Value/Risk).

## Environment notes

- `bash _dev/tests/maintainer-verify.sh` exits 0 at your branch point — baseline and gate. Exit code is the only proof; never pipe through `tail`.
- Toolchain: Go 1.26.1, ShellCheck 0.11.0, `just`, Node 22, Chromium (Playwright; `/opt/pw-browsers/chromium`).
- **Never run bare `go build`** in `skills/do-work-board/tools/queue-kanban/` — build to scratch (`go build -o /tmp/<name> .`).
- Read the clock with `date -u +%Y-%m-%dT%H:%M:%SZ` at the moment you stamp; never carry or compute one.
- Build fixtures in scratch space, never in this repo's own `do-work/`.

## Hand-back

Write your report to exactly this absolute path (the one main-tree write allowed; never stage or commit it):

```
/home/user/skill-do-work/do-work/runs/work-2026-08-18-191338/REQ-254-handback.md
```

Structure: `# REQ-NNN hand-back` with **Branch**, **Commits** (oldest first), `## What I built`, `## File manifest` (one full path per line, `(new|modified|deleted)` + one factual line), `## P-A-U evidence`, `## Testing evidence` (real RED and GREEN output — never from a prototype or memory; the observed maintainer-verify exit code), `## Decisions (D-XX)`, `## Integration seams`, `## Discovered Tasks`, `## Pushback`.

**Standing warning:** recent REQs kept shipping mechanisms that closed an instance while claiming a class — hunt the hole before the reviewer does. The one recent exception (REQ-250) got there by greping the same primitive across the whole file before calling the class closed.

---

# Your REQ (verbatim copy — the live one lives in the main tree)

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
