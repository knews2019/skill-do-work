---
id: REQ-268
title: Never report clean for a state that was never verified
status: completed
created_at: 2026-08-18T21:03:15Z
claimed_at: 2026-08-19T19:11:45Z
completed_at: 2026-08-19T19:45:13Z
commit:
status_changed_at: 2026-08-18T22:20:09Z
user_request: UR-056
addendum_to: REQ-255
domain: general
review_generated: true
effort_estimate: normal
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
route: B
maintenance: false
estimate:
  p50_active_minutes: 20
  confidence: medium
  calculated_at: 2026-08-19T19:14:00Z
  basis:
    - Route B
    - 3-file write set
    - 5 acceptance criteria
    - cross-route regression gates
write_set:
- skills/do-work/scripts/audit-archive-timestamps.sh
- skills/do-work/scripts/repair-req-timestamps.sh
- _dev/tests/prescribed-shell-scripts-behavior.sh
---

# Never Report Clean for a State That Was Never Verified

## What

**The condition, not the file:** an unchecked exit status turns an inspection that never happened into a clean answer. Three instances are known, in two scripts, and all three are reproduced by execution. The REQ was originally scoped to the archive auditor; an external review found the same root cause in the repairer, which is why it is now keyed on the condition instead.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [x] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [x] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Instances

- [x] **A refused stamp is reported as clean.** An archive holding only a calendar-impossible stamp yields `archive audit clean (1 file(s) scanned)` and exit 0, in both report-only and fixing modes. The refusal itself is right — the malformed value is deliberately preserved as evidence — but the reasoning for preserving it is that a human can then see it, and the one tool a person runs deliberately to inspect the archive is exactly where it disappears. Voicing refusals as informational lines, without changing the exit status, keeps the refusal and the report contract consistent.
- [x] **A failed scan is reported as clean.** The file walk runs inside a process substitution, so a nonzero exit from the walk or the sort never reaches the loop's status: with a failing `find` on PATH the auditor prints `archive audit clean (0 file(s) scanned)` and exits 0 while a defective archive file sits untouched. Materialise and validate the walk's output before entering the loop, so an incomplete scan can never be reported as a clean one.

- [x] **The repairer reports clean when its own field extraction fails (added from an external review on PR #145, reproduced by the orchestrator).** `skills/do-work/scripts/repair-req-timestamps.sh:362` assigns `field_rows="$(extract_timestamp_fields "$request_file")"` and then tests only `[ -n "$field_rows" ]`, so a nonzero `awk` exit is discarded and empty output reads as "this file has no timestamp fields". Observed with a failing `awk` first on PATH: a queue file carrying `created_at: 2093-01-01T00:00:00Z` came back **byte-identical and the script exited 0 with no output at all** — while the same file with a working `awk` is repaired. The SessionStart hook discards the script's stderr, so nothing reaches the banner either. This is the third face of the same condition, in the second script.

## Requirements

- **"Clean" is printed only when the inspection actually completed** — every archive file read, every field extraction succeeded — and nothing needed repair. State this as the condition in each script, so a fourth call site inherits it.
- **Sweep the primitive:** every place either script takes a command substitution or a process substitution and then judges only the *content* of the output. These three were found by two independent reviewers looking at other things, so they are a sample.
- A refused defect is visible in the tool's own output, with the exit contract stated to match whatever the builder chooses.
- A scan that could not complete fails loudly rather than reporting a count of zero as success.
- Lock-in cases for both, and `bash _dev/tests/maintainer-verify.sh` exits 0.

## Context

Instance 1 is REQ-255's independent review, finding I-3 (gate: user-visible — the audit's answer is its product, and it misleads exactly on the class it was told to preserve). Instance 2 is an external automated review on pull request 145, reproduced by the orchestrator against the shipped script before recording. Both live in the same file and share one root cause: the tool answers "clean" for states it never verified. Created `pending-answers` per the generation-≥2 cascade stop.

## Open Questions

- [x] The archive auditor says "clean" both when it refused a malformed stamp it deliberately preserved and when its file walk failed outright and scanned nothing. Should I process this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.

**Answered [2026-08-18]:** User approved via `do-work clarify`. Both instances stand: the refused-stamp report line and the failed-walk exit path. The exit contract for a voiced refusal is builder latitude, as the Requirements already say.

---

## Triage

**Route: B** - Medium

**Reasoning:** The three instances name exact files and lines, but the Requirements demand a sweep of a primitive ("every place either script takes a command or process substitution and then judges only the content") across two scripts — the "where" for the fourth-and-beyond instances is unknown until explored.

**Planning:** Not required

## Exploration

Swept both scripts for the primitive the Requirements name — a command or process
substitution whose exit status is discarded while only its content is judged. Five
sites, three of them the REQ's known instances:

- `repair-req-timestamps.sh:362` — `field_rows="$(extract_timestamp_fields …)"` then
  `[ -n "$field_rows" ]`. A failed `awk` reads as "no timestamp fields". **Instance 3.**
- `repair-req-timestamps.sh:662` — the post-rewrite verification re-parses the temp file
  with `done <<< "$(extract_timestamp_fields "$temp_file")"`. A failed `awk` makes the
  loop body never run, `guard_verdict` stays `ok`, and the rewrite is accepted
  **unverified**. Same condition, in the guard whose whole job is verifying. Found by
  the sweep, not in the REQ.
- `repair-req-timestamps.sh:384` — `comparison_key_for` returns empty for a refused
  value (calendar-impossible, numeric offset, fractional seconds, non-ASCII padding).
  Nothing voices it, so a file holding only a refused stamp is reported clean.
  **Instance 1.**
- `audit-archive-timestamps.sh:95` — `done < <(find … | sort -z)`. A failed walk yields
  zero iterations and `archive audit clean (0 file(s) scanned)`. **Instance 2.**
- `audit-archive-timestamps.sh:83-85` — `script_directory="$(cd … && pwd)"` feeding an
  unchecked `source`. If the source fails, `repair_request_file` is undefined, every call
  errors, the loop still completes and the audit still prints clean. Found by the sweep.

Not defects: the clock probes (`:166-193`) validate shape and exit 1 loudly; the rewrite
`awk` at `:632` already captures `$?` — that site is the model the four fixes copy.

**Existing test that must change:** `_dev/tests/prescribed-shell-scripts-behavior.sh:1576`
asserts the refusal-parity output names neither REQ-911 nor REQ-913 at all. Its stated
intent is narrower than its assertion — "logged a refused stamp **as a correction**" —
so voicing a refusal line breaks the assertion without violating the intent.

## Scope

**Files I will touch:**
- `skills/do-work/scripts/repair-req-timestamps.sh` (modify) — voice refusals; check the extractor's status at both call sites; state the condition in the header
- `skills/do-work/scripts/audit-archive-timestamps.sh` (modify) — materialize and validate the walk; verify the source landed; stop printing "clean" when refusals were voiced
- `_dev/tests/prescribed-shell-scripts-behavior.sh` (modify) — lock-in cases for all five sites; tighten the REQ-255 refusal-parity assertion to its stated intent

**Files I will NOT touch:** `skills/do-work/hooks/session-start.sh` (its `|| true` is deliberate — REQ-274 owns the framing there), the board's Go timestamp readers.

**Acceptance criteria (restated from REQ):**
- [x] "Clean" is printed only when the inspection completed and nothing needed repair — stated as the condition in each script
- [x] Every command/process substitution in both scripts that judges only content is swept, not just the three known instances
- [x] A refused defect is visible in the tool's own output, with the exit contract stated
- [x] A scan that could not complete fails loudly rather than reporting zero as success
- [x] Lock-in cases for both, and `bash _dev/tests/maintainer-verify.sh` exits 0

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Decisions

- **D-01**: Refusal voicing is opt-in through a new library switch
  (`timestamp_repair_voice_refusals`), set to 1 by `audit-archive-timestamps.sh` and left
  at 0 for the hook-run repairer. Reasoning: the first attempt voiced every refusal from
  the library unconditionally and broke four existing cases that assert the hook path
  prints nothing for a value it deliberately leaves alone. Those cases are right — a
  refusal is permanent, so an unconditional line would print the same unhealable text into
  every session's start banner forever, which is REQ-267's wedge and exactly the live
  complaint REQ-274 records. The auditor is a deliberate human inspection where that line
  is the product, and it is the tool the REQ's instance 1 is actually about ("the one tool
  a person runs deliberately to inspect the archive is exactly where it disappears").
  Value: instance 1 closes without a banner regression, and no existing test's intent had
  to be renegotiated. Risk: low and reversible — a future caller that wants the lines sets
  the switch. DECIDE & STATE.
- **D-02**: A voiced refusal does not change the exit status; it changes the summary line.
  The Requirements left the exit contract to the builder. Exiting nonzero would make the
  archive audit permanently red for an archive holding one refused stamp that nothing can
  repair, which is the same unhealable-signal failure mode. Instead the auditor stops
  printing "clean" and reports `archive audit complete (N scanned) — M value(s) refused …
  Not clean: those values were never inspected for defects.` DECIDE & STATE.
- **D-03**: The post-rewrite re-parse now also counts how many planned fields it actually
  confirmed and trips the guard unless that equals `planned_count`. Checking only the
  extractor's exit status would still let a re-parse that returned *some* rows verify
  nothing about a planned line missing from them. Not in the REQ; found while fixing the
  site. DECIDE & STATE.

## Implementation Summary

**Files changed:**
- `skills/do-work/scripts/repair-req-timestamps.sh` (modified)
- `skills/do-work/scripts/audit-archive-timestamps.sh` (modified)
- `_dev/tests/prescribed-shell-scripts-behavior.sh` (modified)

**What was done:** Both scripts now judge every command and process substitution by its
exit status before its content, with that condition stated in each header rather than as a
list of the sites that broke. The repairer checks the frontmatter extractor's status at
both call sites — the planning read (which returned "no timestamp fields" for a failed
`awk`) and the verify-before-replace re-parse (which passed vacuously on zero rows, so the
rewrite was accepted unverified) — and counts every refused value through a new
`report_refusal`, printed only when the caller sets `timestamp_repair_voice_refusals`. The
auditor sets it, materializes and status-checks its `find | sort` walk instead of reading a
process substitution, refuses to run at all when its shared library did not load, and
reports `archive audit complete … Not clean` rather than `clean` when any value was refused.
Five lock-in cases pin all five sites; two existing assertions were narrowed from "names the
file at all" to "logs it as a correction", which is what their own messages already said.

## Qualification

Passed — 3 files verified in the diff (47/78/119 changed lines, no placeholders), 5
requirements traced, P-A-U confirmed against the diff (no debug artifacts, no TODO
markers). Scope-drift check: Implementation Summary matches the Scope declaration exactly.

## Testing

**Tests run:** `bash _dev/tests/prescribed-shell-scripts-behavior.sh`, then `bash _dev/tests/maintainer-verify.sh`
**Result:** ✓ All passing (74 named script cases; maintainer-verify exit 0)

**Red-green validation:**
- `audit-archive-timestamps voiced-refusal case` (3 assertions): ✗ before → ✓ after
- `audit-archive-timestamps failed-walk case` (2 assertions): ✗ before → ✓ after
- `audit-archive-timestamps missing-library case` (2 assertions): ✗ before → ✓ after
- `repair-req-timestamps failed-extraction case` (2 assertions): ✗ before → ✓ after
- `repair-req-timestamps unverified-rewrite case` (3 assertions): ✗ before → ✓ after

All twelve assertions were run RED against the unmodified scripts before any fix was
written, and each was independently reproduced by hand first: a calendar-impossible stamp
reported `archive audit clean (1 file(s) scanned)`; a failing `find` on PATH reported
`archive audit clean (0 file(s) scanned)`; a lone copy of the auditor printed two
`command not found` errors and still reported clean; a failing `awk` left a future-stamped
queue file byte-identical, silent, exit 0; and an `awk` failing only on its second
frontmatter extraction let the rewrite land while its verification never ran.

**New tests added:**
- Five named cases in `_dev/tests/prescribed-shell-scripts-behavior.sh` (74 total, up from 69)

**Existing tests updated (cross-REQ impact):**
- `_dev/tests/prescribed-shell-scripts-behavior.sh` refusal-parity case (from REQ-255) and
  shape-parity case (from REQ-267): both asserted their output never *names* REQ-911 /
  REQ-913 / REQ-914. Their own failure messages say "logged a refused stamp as a
  correction", which is the narrower claim, and it is the one that survives the auditor
  now voicing refusals. Narrowed to grep only `would repair` / `repaired` lines. The four
  repairer-side cases that assert the hook path stays silent were left untouched — D-01
  exists so their intent holds.

*Verified by work action*

## Review

**Reviewer:** independent agent, orchestrated mode against the working-tree diff.

| Dimension | Score |
|---|---|
| Requirements Compliance | 80% → 100% after the I1 fix |
| Code Quality | 90% → 100% after the I2 fix |
| Test Adequacy | 95% |
| Scope Discipline | 100% |
| Risk | Low |
| **Acceptance** | **Pass** |
| **Overall** | **91%** (at review time, before the two Important findings were fixed) |

The reviewer reproduced all five red→green pairs independently, two of them by reverting
the actual hunks in a scratch copy and re-running the suite, and confirmed portability on
bash 3.2 / BSD userland. It also verified D-01's factual premise rather than taking it on
trust — that the four repairer-side cases really do assert hook-path silence.

**Important findings — both fixed in this REQ rather than deferred:**

- **I1 — the sweep was incomplete.** Three more substitutions in
  `repair-req-timestamps.sh` judged only content: the post-rename `git diff --numstat`
  guard (`|| true` plus a `0`/`0` fallback, so `[ 0 -gt threshold ]` is false for every
  threshold — the last guard between a bad rewrite and a `repaired` line silently passed,
  *after* the rename), the truncation-floor `git cat-file -s … || echo 0` (which folded
  "no blob in HEAD" together with "the size would not read", skipping the floor for both),
  and the pre-edit numstat baseline. Fixed here, with two new lock-in cases, because the
  REQ's own Requirement says *every* place in either script — deferring the sites inside
  the declared write set would be the one-spelling-at-a-time pattern this REQ exists to
  break. The reviewer's other half — the same primitive inherited from
  `tools/checks/record-commit-hash.sh` and spread across other shipped scripts — is
  genuinely outside this write set and went to **REQ-298** (sweep, `impact-rule-change`).
- **I2 — restatement sweep.** The auditor's own `# Exit 0 —` contract line still said exit
  0 means "clean or fully repaired", contradicting the refusal path added ten lines above
  it. Rewritten to name all three exit-0 outcomes and say only the first is "clean".

Fixing I1 surfaced a regression the reviewer could not have seen: gating the numstat
baseline on its exit status broke a repo with **no HEAD commit at all**, where the failure
is a real absence rather than a git that could not answer. That is the same conflation I1
described, one level up, and it is now a named condition (`head_commit_exists`) with its
own lock-in — a staged-but-never-committed REQ still repairs.

**Minor (report only):** `report_refusal` does not distinguish a recognized-but-rejected
timestamp shape from a non-date value in an `_at`-suffixed key. Pre-existing detection
behavior, only its voicing is new, and this repo's archive has zero such values (262/262
files, 0 refusals). Noted for consumer repos with different `_at` conventions.

## Lessons Learned

**What worked:** Reproducing all three named instances by execution *before* writing any
fix. Each reproduction became the fixture for its lock-in, so the RED run was free and the
test could not be written to match the implementation. The `awk` wrapper that fails only on
its Nth matching invocation is the technique that made the vacuous-verification site
testable at all — an all-or-nothing stub fails too early to reach it.

**What didn't:** Voicing refusals unconditionally from the library. It broke four existing
cases, and the right reading was that those cases were correct and the change was wrong —
a permanent refusal printed from an unattended hook is an unhealable banner line, which is
the exact wedge REQ-267 closed. The switch-gated version cost ten lines and renegotiated no
existing intent. **When a change breaks tests whose intent is sound, the change is the
thing to narrow.**

**Worth knowing:** `read -r a b <<< "$(cmd)"` hides `cmd`'s exit status completely — `$?`
belongs to `read`, and `read` returns nonzero on the ordinary no-output case, which is
exactly why a `|| true` gets attached and the real status is lost for good. Every such site
needs the substitution run as its own statement first. Also: gating a guard on a command's
exit status is only half the fix, because "the command failed" and "there is nothing to
compare against" are different states — HEAD-relative guards in a repo with no commits hit
this immediately.

## Orientation

`do-work forensics`' archive audit and the SessionStart timestamp repairer now refuse to
report a clean answer for anything they did not actually inspect — a failed walk, an
unloadable library, an extractor that could not run, a verification that never ran, and a
value they deliberately refuse each say so in the tool's own output instead of being
absorbed into "clean". Lives in the timestamp-repair subsystem
(`_dev/primes/prime-shell-commands.md`). No map change: same two scripts, same shared-by-
sourcing structure, one new library switch alongside the two that were already there.
Prime spot-check: `prime-shell-commands.md`'s referenced paths all still exist, and its
trap list gains a candidate through REQ-298 rather than through this REQ.
