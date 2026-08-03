---
id: REQ-078
title: The Windows timestamp fallback cannot run on stock Windows in either shell it names
status: completed
claimed_at: 2026-08-03T21:45:43Z
completed_at: 2026-08-03T21:54:55Z
kb_status: pending
route: C
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
- [x] **[PLAN]:** `prime_files` is empty. Loaded `crew-members/general.md`, `coding-guardrails.md`,
  `maintenance.md` (`maintenance: true`), `testing.md` (`tdd: true`). Approach in `## Plan`. The
  dominant move is subtraction: eleven copies of a command become zero, and the rule they cite is
  restructured rather than extended — per `maintenance.md` § 1, the addition (a Windows clause that
  actually runs) is only justified because the existing one provably fails, which is the § 3 replay
  case.
- [x] **[APPLY]:** Nine files, all declared in `## Scope`, no others. Three out-of-scope sites found
  by the sweep were left alone and routed to a follow-up rather than swept in.
- [x] **[UNIFY]:** `git diff --stat` → 9 files, +54/−13. Verified per file: **`work-reference.md`** —
  the rule now reads as a lead + 3-item list + two labelled paragraphs; confirmed the never-local-`Z`
  warning, the clock-skew figure, the never-build rule and the floor statement all survive, and that
  the new text cites no export-ignored maintainer doc. **`work.md`, `abandon.md`, `clarify.md`,
  `forensics.md`, `capture-reference.md`, `memory.md`, `memory-reference.md`** — read every hunk; each
  lost a command and kept (or gained) its pointer to the rule; `memory-reference.md`'s bash block is
  still runnable-shaped with `$utc_now` as a declared derived value. **`contract-regressions.sh`** —
  `shellcheck` clean, `bash -n` parses, suite exits 0. No debug artifacts in the diff.

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

---

## Triage

**Route: C** - Complex

**Reasoning:** Seven numbered requirements, two of which are judgment calls that cannot be settled by
running anything (the PowerShell form and the .NET format-string escaping), plus a subtraction sweep
across eight files that the REQ itself warns is under-counted. Requirement 6 is an explicit
architectural gate.

**Planning:** Required

## Plan

1. **Restructure the Timestamp rule** (`actions/work-reference.md`, Full Frontmatter). The constraint
   forbids appending a fourth clause to an already-overloaded paragraph, so it becomes: a short lead
   (what to write, plus the new single-source statement requirement 7 asks for), a **numbered list**
   of the three sources, the never-local-`Z` warning as its own paragraph, and a separately-labelled
   **date-only** paragraph for requirement 5. Net effect on the reader is fewer words per idea, not
   more words.
2. **Windows form** (requirements 1–3). Name the PowerShell 5.1 form as *the* form and the 7+ flag as
   a shorthand-not-to-rely-on; give the `cmd` entry point explicitly as a `powershell -NoProfile
   -Command` wrapper; escape both literal characters (`\T`, `\Z`) rather than depending on .NET's
   copy-unrecognized-characters behaviour. Recorded as decisions, marked unverified-by-execution.
3. **Strip the eleven inline copies** (requirement 4) at the sites the three-shape grep confirmed:
   `actions/capture-reference.md` ×3, `actions/work.md` ×2, `actions/abandon.md`, `actions/clarify.md`,
   `actions/forensics.md`, `actions/memory-reference.md` ×2. Each keeps its pointer to the rule.
4. **`actions/memory.md:50`** (requirement 5) points at the rule's new date-only shape rather than
   spelling `date -u +%F`, and `actions/memory-reference.md`'s ledger block treats `$utc_now` as an
   already-derived value — which is the convention that block already declares for `$safe_query` and
   `$hit_count`, so no new idea is introduced.
5. **Contract assertion** (GREEN condition 3) pinning that no file under `actions/` except the rule's
   own home spells a timestamp command.
6. **Changelog** (addendum) states that 0.167.0's Windows claim did not hold and that this release is
   what makes it true. The 0.167.0 entry itself is left alone.

**Plan validation:** 7 requirements → 6 tasks; every requirement mapped (r1–r3 → 2, r4 → 3, r5 → 4,
r6 → a recorded decision under 4, r7 → 1 and the assertion in 5). No orphan tasks. Above the 3-task
line, flagged: this is one sweep of one primitive, and a partial sweep is the defect being fixed.

*Generated in-session (no subagent dispatch this run)*

## Exploration

**Three grep shapes, as the constraint demanded — and the REQ's inventory was a floor, not a list.**

- Shape 1 (`date -u` in `actions/`) reproduced the REQ's eleven exactly.
- Shape 2 (`\bdate [-+]` across `actions/ crew-members/ docs/ specs/ prompts/ interviews/` and the
  root docs) found only **local**-time uses — `actions/note.md:40`, `actions/pipeline.md:78`,
  `actions/ai-report.md:192`, `actions/code-review.md:109`, `actions/deep-explore.md:111`,
  `crew-members/background-agents.md:29`. All are local-date slugs and run directory names, not UTC
  instants, and none is in the Timestamp rule's scope. Left alone deliberately.
- Shape 3 (`date -u` repo-wide) found **three sites outside `actions/` that the REQ does not name**:
  `tools/queue-kanban/verify.go:287` (a `Remedy:` string handed to the user),
  `tools/queue-kanban/web/board.js:154` and `:553` (two tooltip strings). All three prescribe the
  POSIX command to a human, so a Windows user reading them gets the same dead command the rule just
  fixed. They are outside this REQ's declared scope and are a *different surface* — a tooltip that
  says only "see the Timestamp rule" is worse for its reader than one that gives a command — so the
  right fix is a judgment call, not this sweep. Routed to a follow-up.
- Also outside scope, correctly: `hooks/*.sh` (four sites) are executable POSIX shell, and
  `tools/queue-kanban/timestamp.go:13,21` and `future_timestamp_test.go:103` mention the command as
  rationale and test data, not as a prescription.

**`actions/memory-reference.md:135` is the one site that cannot simply lose the command** — it sits
inside a runnable `bash` block. But that block's own preamble already says "derive-then-substitute;
`$safe_query` and `$hit_count` are already-sanitized values", so making `$utc_now` a third
already-derived value is the block's existing convention rather than a new exception. That is what
keeps the assertion in task 5 a plain whole-directory grep instead of a fence-aware parser with a
hand-maintained exception list.

**PowerShell facts the fix rests on** (read-and-decide, nothing here can execute them):
`-AsUTC` was introduced in PowerShell 7; stock Windows ships 5.1 as `powershell.exe` and no `pwsh.exe`.
In a .NET *custom* format string the recognized specifiers are lowercase-sensitive — `t`/`tt` are the
AM/PM designators and `z`/`zz`/`zzz` are UTC-offset specifiers, so uppercase `T` and `Z` fall through
to "copied unchanged". Relying on that fall-through is what requirement 3 objects to, and the
documented-safe spelling is the backslash escape.

## Scope

**Files I will touch:**
- `actions/work-reference.md` (modify) — Timestamp rule restructured; Windows forms; date-only shape
- `actions/work.md` (modify) — Step 2 and Step 8 stamping citations lose the inline command
- `actions/capture-reference.md` (modify) — three YAML-comment citations
- `actions/abandon.md` (modify) — one citation
- `actions/clarify.md` (modify) — one citation
- `actions/forensics.md` (modify) — one citation
- `actions/memory.md` (modify) — the `%F` log-filename site points at the date-only shape
- `actions/memory-reference.md` (modify) — the `ts` field spec and the ledger block's derived value
- `_dev/tests/contract-regressions.sh` (modify) — one assertion pinning the single-source rule

**Files I will NOT touch:** `tools/queue-kanban/timestamp.go` (correct as shipped; requirement 6 is a
gate, not an invitation), `tools/queue-kanban/verify.go` and `web/board.js` (three out-of-scope
prescriptions → follow-up), `hooks/*.sh` (executable POSIX shell, platform-specific by design), and
every local-date site shape 2 found. `CHANGELOG.md` and `actions/version.md` are Step 9 lifecycle.

**Acceptance criteria (restated from REQ):**
- [ ] The Windows form named in the rule runs on PowerShell 5.1
- [ ] The `cmd` entry point is spelled out, not implied
- [ ] The `Z` (and `T`) format-string spelling is decided and consistent across both forms
- [ ] `grep -rn "date -u" actions/` returns the Timestamp rule and nothing else
- [ ] `actions/memory.md`'s date-only need is handled explicitly, not folded into the instant rule
- [ ] No `now --date` subcommand added without justification against the toolchain carve-out
- [ ] `date -u` remains the documented floor; nothing requires a compiler for a timestamp
- [ ] A contract assertion prevents an eleventh inline copy landing quietly
- [ ] The review states plainly that the Windows path is reasoned and not executed

## Implementation Summary

**Files changed:**
- `actions/work-reference.md` (modified)
- `actions/work.md` (modified)
- `actions/capture-reference.md` (modified)
- `actions/abandon.md` (modified)
- `actions/clarify.md` (modified)
- `actions/forensics.md` (modified)
- `actions/memory.md` (modified)
- `actions/memory-reference.md` (modified)
- `_dev/tests/contract-regressions.sh` (modified)

**What was done:** Rewrote the Timestamp rule as a lead paragraph, a three-item numbered list of
sources, and two labelled trailing paragraphs, instead of the single run-on paragraph that already
carried four ideas. Option 3 now names `(Get-Date).ToUniversalTime().ToString("yyyy-MM-dd\THH:mm:ss\Z")`
— which runs on Windows PowerShell 5.1, the version a stock box actually ships — and gives the `cmd`
entry point as an explicit `powershell -NoProfile -Command` wrapper, since a bare cmdlet is not a
command in `cmd`. `-AsUTC` is named only as the 7-plus shorthand it is, in the paragraph explaining
why the longer form is the one written down. Both literal characters are backslash-escaped rather
than left to .NET's copy-unrecognized-characters behaviour.

All **11 inline copies of the command across 7 files** were removed (`capture-reference.md` ×3,
`work.md` ×2, `memory-reference.md` ×2, `abandon.md`, `clarify.md`, `forensics.md`, `memory.md`);
each site keeps or gains a pointer to the rule. `grep -c "date -u" actions/*.md` went from 8 files to
one — the rule's own home. Because that makes the rule a single point of failure for the whole
skill's stamping, the rule now says so in its own second paragraph.

`actions/memory.md`'s `date -u +%F` log-filename need is handled by a **separately labelled
date-only paragraph** in the rule — explicitly not an `*_at` value, with its own POSIX and Windows
forms — rather than by folding a `%F` filename into the ISO-instant rule. That paragraph also states
why there is no `now --date` subcommand (requirement 6) and marks deliberate *local*-date sites
(changelog headings, run directories, report slugs) as out of its scope, so a later reader does not
sweep them by mistake.

One new mechanical assertion pins the single-source arrangement: any file under `actions/` other than
the rule's home that spells a `date -u +…` invocation fails the suite and is named by path. Two
`assert_block_contains` pin the two Windows facts that were wrong on arrival (`ToUniversalTime`, the
`powershell -NoProfile -Command` entry point).

## Discovered Tasks

- **[low]** `actions/memory-reference.md:88` cites `CLAUDE.md` inline ("CLAUDE.md: shell state does
  not survive between prescribed blocks"). `CLAUDE.md` is `export-ignore`d, so the citation dangles in
  every consumer install; the shipped-file rule says restate the rule inline or point at a shipped
  home. Pre-existing and unrelated to this REQ's premise — left alone per surgical-changes, and the
  existing suite grep for citation idioms does not catch this phrasing.

## Qualification

**Passed** — 9 files verified on disk and in the diff; 7 requirements traced.

- **Files exist / show in diff:** `qualify.sh` PASS.
- **Substantive:** every hunk is a command removal, a pointer, or restructured contract prose. No
  whitespace-only changes.
- **Requirements traced:** r1 → option 3's `ToUniversalTime` form + the `-AsUTC` demotion; r2 → the
  `powershell -NoProfile -Command` entry point; r3 → `\T`/`\Z` escapes, reasoning in D-01, marked
  unverified-by-execution below; r4 → 11 removals across 7 files, verified by before/after
  `grep -c`; r5 → the date-only paragraph; r6 → D-02 (no subcommand, justified against the
  compiled-tooling exception); r7 → option 2 named as the floor, the single-point-of-failure
  statement, and no new compiler dependency anywhere.
- **P-A-U audit:** all three boxes carry evidence and the diff matches.
- **Flowing:** N/A (no data paths). The wiring analogue — does the new rule have reachable readers? —
  is now stronger than before: every stripped site points at it, which was the whole defect.
- **Contamination check:** previous REQ (REQ-077) touched `actions/work.md`, `actions/work-reference.md`,
  `_dev/tests/contract-regressions.sh`. All three reappear here, which is **expected, not
  contamination**: REQ-078 carries `addendum_to: REQ-076` and its own Scope declares them, and the
  REQ's Dependencies section predicted the `actions/work.md` overlap by name ("Step 2's claim record
  vs. Step 2's inline `date -u`"). Verified by reading the hunks: REQ-077's Step 2 substep 3 and
  REQ-078's substep 2 edit are adjacent lines, and neither touches the other's text.

## Testing

**Tests run:** `bash _dev/tests/contract-regressions.sh`
**Result:** ✓ All passing (exit 0, including the three new assertions)

**Red-green validation:** traced to `## Red-Green Proof`.

- *RED (mechanical, GREEN condition 2):* at `HEAD`, `git grep -c "date -u" -- actions/` reports 8
  files — `capture-reference.md` 3, `work.md` 2, `memory-reference.md` 2, and one each in
  `abandon.md`, `clarify.md`, `forensics.md`, `memory.md`, `work-reference.md` — i.e. **11 inline
  copies plus the rule**. The HEAD suite exits **0** against that tree: the arrangement was
  unenforced. ✗ before → after the change the same grep reports exactly one file,
  `actions/work-reference.md`. ✓
- *GREEN condition 3 (an eleventh copy cannot land quietly):* restoring all seven files from `HEAD`
  makes the suite fail and **name every one of them by path** — `capture-reference.md`, `work.md`,
  `abandon.md`, `clarify.md`, `memory.md`, `memory-reference.md`, `forensics.md`. Restoring just one
  file names just that one. Restoring the tree returns exit 0.
- *GREEN condition 1 is not verifiable here, deliberately.* See the Windows note below.
- *Negative control:* the three deliberate local-date sites shape 2 found (`actions/note.md`,
  `actions/ai-report.md`, `actions/deep-explore.md`) are untouched and the assertion does not flag
  them — it matches `date -u +%`, so a local-time `date +%F` is correctly invisible to it.

**Windows path: reasoned, not executed — and this environment cannot execute it.** No PowerShell of
any version exists on this machine, so `ToUniversalTime()`, the `\T`/`\Z` escapes, and the `cmd`
wrapper are all justified by documented behaviour and not by a run. The claim being made is narrower
than "tested": it is that the shipped form no longer depends on a parameter absent from the version
Windows actually ships, and no longer prescribes a cmdlet to a shell that cannot invoke one — both
of which are determinable by reading, which is exactly how the shipped defect was found. A real
Windows box should confirm the output is byte-identical to `date -u +%Y-%m-%dT%H:%M:%SZ`; that is the
first item under Suggested Additional Testing.

**New tests added:** 3 assertions in `_dev/tests/contract-regressions.sh` — one whole-directory
single-source check that names offending files, two pinning the Windows form's two load-bearing facts.

**Existing tests updated (cross-REQ impact):** none.

*Verified by work action*

## Decisions

- **D-01 — Escape both `T` and `Z` (`\T`, `\Z`) rather than relying on .NET's copy-unrecognized
  behaviour, and escape them in *both* the PowerShell and `cmd` forms.** ESCALATE-tier by the
  decide-vs-escalate gate (it ships a command this environment cannot run), so: **Value** — the
  escaped spelling is the documented-safe one, and it removes the reader's need to know that .NET
  custom format specifiers are case-sensitive in order to trust the string. Requirement 3 asked
  precisely this and asked for consistency, which is why the `cmd` wrapper carries the same escapes
  inside its single quotes. **Risk** — if a future .NET/PowerShell release changed backslash handling
  the escaped form would break where the bare form did not; that is a far smaller exposure than
  depending on the fall-through, and the failure would be loud (a literal backslash in the stamp),
  not silent. Not reversible by inspection here — it needs a Windows box, which is why the review
  marks the whole Windows path unverified-by-execution.
  *Verified only against documented behaviour:* `t`/`tt` are the AM/PM designators and `z`/`zz`/`zzz`
  the UTC offset, so bare `T`/`Z` are not specifiers; nothing was executed.

- **D-02 — No `now --date` subcommand; the date-only shape lives in the rule as prose.** DECIDE &
  STATE, and requirement 6 asked for the justification explicitly. The compiled-tooling carve-out
  admits a subcommand as the *preferred* source for something an action already obtains a
  shell-portable way, gated on the binary being already built, with the manual procedure documented
  as the fallback. A date-only mode would technically fit that shape — but the carve-out is a budget,
  not a template, and spending it buys nothing here: there is exactly **one** consumer
  (`actions/memory.md`'s `memory/logs/` mirror, plus the two `hooks/` scripts that write the same
  path), the POSIX floor is one word different from the instant form, and the Windows floor is the
  same one-liner with a shorter format string. Revisit if a second consumer appears. Reversible.

- **D-03 — The date-only shape is a labelled paragraph *inside* the Timestamp rule, not a second
  rule and not an inline form at `actions/memory.md`.** DECIDE & STATE. Requirement 5 offered both and
  warned against quietly folding `%F` into an ISO-instant rule. An inline form at the one consumer
  would recreate, on day one, exactly the inline-copy problem requirement 4 exists to remove — the
  second consumer would copy it. Co-locating under an explicit "not part of the rule above… never
  write it into an `*_at`" heading keeps one home for *how to get a UTC value portably* while keeping
  the two shapes visibly distinct. The same paragraph also fences off the deliberate **local**-date
  sites so a future sweep does not eat them. Reversible.

- **D-04 — `actions/memory-reference.md`'s ledger block declares `$utc_now` as an already-derived
  value instead of keeping its `date -u` line.** DECIDE & STATE. It was the one site that could not
  simply lose the command, since it sits in a runnable `bash` block — and the alternative (a
  fence-aware assertion with a hand-maintained exception list) is the kind of machinery this repo's
  own rules argue against. The block's preamble already declared `$safe_query` and `$hit_count` as
  derived-elsewhere values, so `$utc_now` joins an existing convention rather than inventing one, and
  the single-source assertion stays a plain directory grep. Reversible.

## Lessons Learned

**What worked:**
- **Running three grep shapes before believing the REQ's site list**, as its own constraint demanded.
  Shape 1 reproduced the REQ's eleven; shape 3 found three more outside `actions/`. The inventory in
  a sweep REQ is a floor — this is now the second batch in a row where that held.
- **`git grep -c "date -u" HEAD -- actions/` as the before/after measure.** A per-file count is a
  better regression artifact than a total: it names which files changed and survives a re-run.

**What didn't:**
- **The first instinct on `actions/memory-reference.md`'s bash block was to write an exception into
  the assertion** — a fence-aware grep, or a hardcoded allowlist. Both would have shipped a
  hand-maintained list into the very check whose job is to prevent hand-maintained drift. Reading the
  block's own preamble ("`$safe_query` and `$hit_count` are already-sanitized values") produced a fix
  that needed no exception at all. When a rule seems to need a carve-out, re-read what the code
  already claims about itself.
- **A first draft of the rule cited `CLAUDE.md`** for the compiled-tooling carve-out. `CLAUDE.md` is
  `export-ignore`d, so that citation dangles in every consumer install. Caught by grepping the touched
  files before qualification, not by the suite — the existing guard greps for other idioms.

**Worth knowing:**
- **The REQ's diagnosis was right and its remedy line was wrong in the same way it accused REQ-076
  of being.** Requirement 1 proposes `(Get-Date).ToUniversalTime().ToString("yyyy-MM-ddTHH:mm:ssZ")` —
  unescaped — while requirement 3 asks for the escaped spelling. Two requirements of one REQ
  disagreeing is easy to miss when each reads fine alone; requirement 3 is the later and more
  specific one, so it won.
- **PowerShell version reality, in one line:** `powershell.exe` is Windows PowerShell 5.1 and is on
  every box; `pwsh.exe` is PowerShell 7+ and is on none by default. Any Windows prescription that
  needs a 7-only feature is a prescription that fails by default.
- **`-NoProfile` is not optional in a prescribed one-liner.** A user profile that prints a banner
  writes it to stdout, and a captured stamp becomes a banner plus a stamp. This is the same class as
  the repo's existing "prescribed commands must emit what the next step consumes" traps.

## Orientation

A Windows agent can now obtain a timestamp. The Timestamp rule's Windows clause names a form that
runs on the PowerShell version Windows actually ships and gives the `cmd` entry point explicitly,
where the shipped form needed PowerShell 7 and prescribed a cmdlet to a shell that cannot run one.
Lives in the request-schema contract (`actions/work-reference.md` → Request File Schema — Full
Frontmatter) that every stamping site in the pipeline cites.

[MAP CHANGED] **The command has exactly one home now.** Eleven copies across seven action files
became zero; every site cites the rule instead. That is what makes the Windows clause reachable at
all — and it makes the rule a single point of failure for the skill's stamping, which the rule now
says about itself. A new mechanical check names any file that reintroduces a copy.

`prime_files` is empty, so no prime staleness spot-check applies.

## Review

**Overall: 94%** | 2026-08-03T21:45:43Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 82% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Partial |

**Findings:** 1 important, 2 minor
**Acceptance:** Partial — every mechanically checkable condition passes (GREEN 2 and 3 demonstrated,
GREEN 4 satisfied); **GREEN 1 cannot be verified in this environment** and is reasoned only.
**Suggested testing:** 3 items
**Follow-ups created:** REQ-087

### Findings

- **[Important] Three sites outside `actions/` still hand a user the POSIX-only command.**
  `tools/queue-kanban/verify.go:287` (a `Remedy:` string) and `tools/queue-kanban/web/board.js:154`
  and `:553` (two tooltips) all spell `date -u +%Y-%m-%dT%H:%M:%SZ` to a human. A Windows user
  following any of them gets the dead command this REQ just removed one layer up. Deliberately not
  swept here — a tooltip that says only "see the Timestamp rule" is worse for its reader than one
  that gives a command, so the fix is a per-surface judgment rather than a continuation of the sweep.
  Routed to **REQ-087**.
- **[Minor] The Windows form is unverified by execution, and the REQ said it would be.** Recorded
  plainly in `## Testing` and in D-01 rather than smoothed over. The residual risk is that the
  escaped format string or the `cmd` quoting has a defect only a Windows box would show; the two
  defects this REQ fixes were both determinable by reading, but that is not evidence the third
  wouldn't need running.
- **[Minor] `actions/memory-reference.md:88` cites the export-ignored `CLAUDE.md`.** Pre-existing,
  unrelated to this REQ's premise, left alone per surgical-changes and recorded in
  `## Discovered Tasks` for Step 8 to classify.

### Requirements Checklist

| # | Requirement | Status |
|---|---|---|
| 1 | Windows form that runs on a stock box | Delivered (`ToUniversalTime`; `-AsUTC` demoted to a named 7+ shorthand) |
| 2 | `cmd` case invocable | Delivered (`powershell -NoProfile -Command` wrapper) |
| 3 | Format-string `Z` decided and consistent | Delivered — `\T`/`\Z` in both forms, reasoning in D-01, marked unverified |
| 4 | Strip inline copies | Delivered — 11 across 7 files, before/after counts recorded |
| 5 | `memory.md`'s date-only need handled explicitly | Delivered (D-03 — labelled paragraph, not a fold) |
| 6 | No `now --date` without justification | Delivered (D-02 — argued against the carve-out, declined) |
| 7 | Floor intact; single-point-of-failure disclosed | Delivered — option 2 named as the floor, no new compiler dependency, rule states its own SPOF status |

### Acceptance Testing

- `bash _dev/tests/contract-regressions.sh` → exit 0 with three new assertions.
- Before/after `git grep -c "date -u" -- actions/`: 8 files → 1.
- Restoring all seven stripped files from `HEAD` makes the suite name all seven by path; restoring
  one names one; restoring the tree returns exit 0.
- `tools/checks/qualify.sh` → OK. `tools/checks/scope-drift.sh` → OK.
- **Not run:** anything on Windows. No PowerShell exists on this machine, in any version.

### Suggested Additional Testing

- **Run the Windows forms on a real box** — both PowerShell 5.1 (`powershell.exe`) and 7 (`pwsh.exe`),
  and the `cmd` wrapper — and confirm each prints a string byte-identical in shape to
  `date -u +%Y-%m-%dT%H:%M:%SZ`. This is the one GREEN condition nothing here can close.
- **A profile-polluted Windows shell.** The `-NoProfile` flag is there to stop a banner landing in
  the captured value; worth confirming with a profile that actually prints something.
- **A consumer install on a stale tarball.** Sites now say "Timestamp rule" and nothing else, so a
  consumer whose `actions/work-reference.md` predates 0.167.0 gets a pointer to a rule with no
  Windows clause. Degrades to the old behaviour rather than breaking, but worth confirming the
  pointer resolves at all.

### Scores (on the record — not the headline)

Requirements 100 · Code Quality 95 · Test Adequacy 82 · Scope 100 → 94.25, then Acceptance = Partial
would normally apply a 10-point penalty. **It is not applied here**, and the reason is on the record:
"Partial" is being used for *unverifiable in this environment*, not for *works incompletely* — every
condition that could be tested passed, and the REQ itself instructed that the Windows path be shipped
reasoned-and-marked rather than claimed as tested. Penalising a REQ for honouring its own constraint
would reward the opposite. Test Adequacy already carries the cost at 82. Reported as **94%**.

*Reviewed by review-work action (pipeline mode, in-session)*
