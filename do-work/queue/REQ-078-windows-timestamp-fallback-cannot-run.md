---
id: REQ-078
title: The Windows timestamp fallback cannot run on stock Windows in either shell it names
status: pending
created_at: 2026-08-03T16:53:42Z
user_request: UR-015
domain: general
prime_files: []
tdd: true
depends_on: []
maintenance: true
addendum_to: REQ-076
---

# The Windows timestamp fallback cannot run on stock Windows in either shell it names

## What

REQ-076 (v0.167.0) added a Windows fallback to the Timestamp rule at `actions/work-reference.md:91`.
As shipped it fails on a stock Windows box in both shells it mentions: the cmdlet flag it uses requires
PowerShell 7+, and the command is offered as the remedy *for `cmd`*, where a bare cmdlet is not a
command at all. The correct PowerShell 5.1 form exists only in the session's own notes, never in the
shipped file.

Then do the subtraction REQ-076 identified and declined to do unasked: ten sites spell `date -u`
inline, so nobody following them ever reaches the preference order.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

0.167.0's headline claim is that it "fixes a real gap on Windows `cmd`, where the prescribed
`date -u +FORMAT` doesn't exist at all." The diagnosis is correct and worth keeping. The prescription
does not hold, which makes the entry's claim stronger than the code. REQ-076's own review flagged the
Windows form as "reasoned, not run" — that understates it: it is reasoned *and wrong*, and the two
failure modes were both determinable by reading, without executing anything.

The subtraction half matters because without it the whole preference order is decorative: an agent
following `actions/work.md` Step 2 sees only `date -u`, so a Windows agent still gets the broken
command even after this REQ fixes the rule.

## Context

`actions/work-reference.md:91`, option 3 as shipped:

> (3) on Windows `cmd`, where that flag form does not exist, PowerShell's
> `Get-Date -AsUTC -Format "yyyy-MM-ddTHH:mm:ssZ"`

**Defect (a) — the flag is PowerShell 7+ only.** `-AsUTC` was added in PowerShell 7. Every Windows
install ships PowerShell 5.1 as `powershell.exe`; PowerShell 7 is `pwsh.exe` and is not present by
default. On 5.1, `-AsUTC` is not a recognized parameter and the call fails outright. The working 5.1
form is `(Get-Date).ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ssZ")`.

**Defect (b) — a cmdlet is prescribed for `cmd`.** The stated context is Windows `cmd`, where
`Get-Date …` yields `'Get-Date' is not recognized as an internal or external command`. Reaching
PowerShell from `cmd` needs an explicit invocation
(`powershell -NoProfile -Command "…"`), which the rule does not supply.

**Defect (c) — one inline site the "rarely reached" finding missed.** REQ-076's review counted the
sites that spell `date -u` inline and got the count right:
`actions/capture-reference.md:16,141,167`, `actions/work.md:225,538`, `actions/abandon.md:58`,
`actions/clarify.md:141`, `actions/memory-reference.md:124,135`, `actions/forensics.md:158`.
It missed `actions/memory.md:50`, which uses `date -u +%F` to build a log filename. Same POSIX
dependency, different format, and **outside** the Timestamp rule's `*_at` scope — so pointing that
site at the rule is the wrong fix and subtraction alone will not reach it.

## Detailed Requirements

1. **Give a Windows form that runs on a stock box.** Name the 5.1 form
   `(Get-Date).ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ssZ")` as the one that always works, and
   `Get-Date -AsUTC -Format "yyyy-MM-ddTHH:mm:ssZ"` at most as a 7+ shorthand. If only one form is
   kept, keep the 5.1 one — it works on both.
2. **Make the `cmd` case invocable.** Prescribe the wrapper an agent in `cmd` actually needs
   (`powershell -NoProfile -Command "…"`), or restate option 3's context as "in PowerShell" and add
   the `cmd` entry point separately. Either is fine; silently naming a cmdlet as the `cmd` remedy is not.
3. **Verify the format string's `Z` before shipping it.** In a .NET custom format string, `Z` is not a
   standard specifier and is emitted literally, which is why the form appears to work — but the
   documented-safe spelling escapes it (`\Z` or a quoted literal). Confirm which spelling this project
   wants and use it consistently in both forms. **This requirement is a read-and-decide, not a test:**
   nothing in this environment can execute PowerShell, so record the reasoning in `## Decisions` and
   mark the Windows path unverified-by-execution in the review, exactly as REQ-076 did — do not claim
   it as tested.
4. **Strip the inline `date -u` from the ten citing sites** listed in Context so exactly one place —
   the Timestamp rule — states the mechanism. Each site keeps its pointer to the rule and loses the
   command. This is the fix REQ-076 identified as serving both goals and left as a minor finding.
5. **Handle `actions/memory.md:50` separately, not by subtraction.** It needs a UTC *date* for a
   filename, which the Timestamp rule does not cover. Either extend the rule to name a date-only shape
   alongside the instant shape, or give that site its own portable form. Say which and why — do not
   quietly fold a `%F` filename into an ISO-instant rule.
6. **Do not add a `now --date` style subcommand for requirement 5 unless it is genuinely warranted.**
   `tools/queue-kanban/` gained `now` under a narrow carve-out in CLAUDE.md's toolchain exception
   ("preferred source for something an action already obtains a shell-portable way"). A second
   timestamp shape may qualify, but justify it against that carve-out explicitly rather than assuming
   the precedent extends.
7. **Keep the floor intact.** After this REQ, no action may require a compiler for a timestamp, and
   `date -u` stays the documented POSIX fallback. If requirement 4 removes the last inline copy, the
   rule is now a single point of failure for the whole skill's stamping — say so in the rule.

## Constraints

- `actions/work-reference.md:91` is a single long paragraph that already carries the never-build rule,
  the local-time-with-Z warning, and the clock-skew allowance. It is at the edge of readable. Prefer
  restructuring it over appending a fourth clause — `crew-members/maintenance.md`'s
  delete-before-you-add rule applies.
- Requirement 4 touches seven files. Grep the primitive across **all** actions before calling it done;
  the ten sites in Context were found by one grep shape and a second shape may find more, which is the
  exact pattern REQ-075 hit when its declared five sites turned out to be eleven.
- No Go logic changes. `tools/queue-kanban/timestamp.go` is correct as shipped (verified: it converts
  rather than relabels, and truncates rather than rounds) and requirement 6 is a justification gate,
  not an invitation.

## Dependencies

`addendum_to: REQ-076` — amends the rule REQ-076 wrote and completes the subtraction it deferred. No
`depends_on`: buildable immediately, and independent of REQ-077 despite both touching
`actions/work.md` (different sections: Step 2's claim record vs. Step 2's inline `date -u`; note the
overlap is one line apart, so if these two are ever fan-out dispatched together the merge is the proof,
not the overlaps badge).

## Builder Guidance

**Certainty: Firm on (a) and (b), read-and-decide on (c) and the format-string question, open on
requirement 5.**

Defects (a) and (b) are determinable by reading and should not be re-litigated. Requirement 3 is where
to spend care, and the honest outcome is a recorded decision plus an unverified marker — this
environment cannot run PowerShell, and a review that claims otherwise is the failure this requirement
guards against.

Requirement 5 has real latitude. The cheap answer (give `actions/memory.md:50` its own portable form)
is probably right; the expansive answer (a date-only shape in the Timestamp rule) is defensible if
other sites want it. Check for other `%F`-shaped uses before choosing.

## Red-Green Proof

**RED case:** Read `actions/work-reference.md:91` as an agent on a stock Windows box. Option 3 names a
PowerShell 7 parameter for a shell (`cmd`) that cannot invoke cmdlets — both failures are visible
without executing anything, which is what makes this provable by inspection.

Second RED, mechanical: `grep -rn "date -u" actions/ | wc -l` returns 11 (ten Timestamp-rule citations
plus `actions/memory.md:50`), while the rule that is supposed to be the single source is one of them.

**Why RED now:** The Windows gap REQ-076 diagnosed is still open — the prescription changed, the
reachability didn't — and the preference order it added is unreachable from any site an agent actually
follows.

**GREEN when:** (1) The Windows form named in the rule runs on PowerShell 5.1, and the `cmd` entry
point is spelled out. (2) `grep -rn "date -u" actions/` returns exactly the Timestamp rule, plus
whatever requirement 5 deliberately leaves at `actions/memory.md:50` with its reason recorded. (3) A
contract assertion pins the inline-command absence so an eleventh copy cannot land quietly. (4) The
review states plainly that the Windows path is reasoned and not executed.

**Validation:** Inferred during an adversarial audit; remediation plan reviewed and approved by the
user before capture.

## Full Context

See `do-work/user-requests/UR-015/input.md` for the audit's provenance and the findings it cleared.

---

## Addendum (2026-08-03)

An **external audit**, triaged separately via `do-work validate-feedback` and captured as UR-016,
reached the same conclusion independently: "Ten shipped sites still directly prescribe `date -u`,
bypassing the new preferred command. When the binary is absent, the Windows fallback uses `Get-Date
-AsUTC`, which the project's own review acknowledges requires PowerShell 7 and fails on Windows
PowerShell 5.1." Same two defects, same ten-site count, same 5.1 replacement form. The user's
instruction was to fold that audit's evidence into this REQ rather than duplicate it.

**It found nothing this REQ did not already have** — and less: it missed defect (b) (a cmdlet
prescribed for `cmd`), defect (c) (`actions/memory.md:50`), and the `Z` format-string question. Recorded
as independent corroboration, which is worth exactly one thing: the diagnosis is now confirmed by two
readers who did not share notes, so requirements 1 and 2 should not be re-litigated during the build.

**One genuine addition — the shipped changelog entry stays uncorrected.** `CHANGELOG.md`'s 0.167.0
entry claims the release "Fixes a real gap on Windows `cmd`, where the prescribed `date -u +FORMAT`
doesn't exist at all." As shipped it does not, for both reasons above. This REQ's `## Why` notes that
the entry's claim is "stronger than the code," but no requirement acts on it — whereas REQ-077's
requirement 7 does exactly that for its own finding ("Disclose the F1 regression in the changelog
entry, not only the fix"). Add the equivalent here: **this REQ's own changelog entry should state that
0.167.0's Windows claim did not hold and that this release is what makes it true.** Do not rewrite the
0.167.0 entry itself — history stays; the correction belongs in the new entry.

No contradiction with anything above.
