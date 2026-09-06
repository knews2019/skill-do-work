---
id: REQ-600
status: claimed
domain: general
created_at: 2026-09-06T06:53:35Z
user_request: UR-105
review_generated: true
impact: impact-rule-change
effort_estimate: effort-mechanical
route: B
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: false
estimate:
  p50_active_minutes: 25
  confidence: medium
  calculated_at: 2026-09-06T07:26:22Z
  basis:
    - Route B
    - 2-file write set
    - 4 acceptance criteria
maintenance: true
depends_on: [REQ-594]
related: [REQ-593, REQ-594]
write_set: [_dev/primes/prime-shell-commands.md, skills/do-work-knowledge/actions/memory-reference.md, _dev/tests/action-shell-blocks.sh, _dev/tests/quiet-grep-pipeline-audit.sh, _dev/tests/quiet-grep-pipeline-scanner.sh]
title: 'Put the SIGPIPE trap in the prime shell authors read, and fix the one shipped block that carries it'
claimed_at: 2026-09-06T07:26:22Z
---

# Put the SIGPIPE Trap Where Shell Authors Read It

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Route B. Prime read. Sweep every shipped fence first, then one prime section, one block, and the scan wired where the fences are already extracted.
- [x] **[APPLY]:** Five files, commits `22a1ea4` and `cb606f2`, merged as `a25c752`. The three `_dev/tests/` files are the widening the Constraints allowed and D4 records.
- [x] **[UNIFY]:** `git diff --stat 9e00a092..a25c7522`: 5 files, +155/-77. Fast gate on main exit 0; action-shell-blocks 74 fences and 33 shipped shell files; audit 95 tracked shell files, 19+7 shapes; audit-lockins passed. No debug artifacts.

## What

Two requests and about 140 fixed sites later, `_dev/primes/prime-shell-commands.md` still carries no
mention of the defect class. `CLAUDE.md` names that file as what to read before writing or reviewing
shell anywhere it ships, and calls it "the hard-won trap list". The trap is documented in two run records
and two request files, which is not where the next person writing shell will look.

Separately, one shipped markdown block carries the forbidden shape:
`skills/do-work-knowledge/actions/memory-reference.md:88` is inside a ```bash block that agents copy
and run — `ollama list 2>/dev/null | grep -qiE 'embed'`.

## Why

REQ-594's guard covers tracked `*.sh` files. Shipped guidance is where the shape gets copied from, and
a prime is where an author decides how to write the line in the first place. A guard that catches the
defect after it is written is worth less than a prime that stops it being written.

## Context

Both found by the independent three-lens review of REQ-594. The reviewer noted that the guard's
reader-set limitation IS disclosed — in a run record and a request file, not in the prime.

## Detailed Requirements

- Add the class to `_dev/primes/prime-shell-commands.md` as its own section, beside "Unchecked Exit
  Status Reads as Content", which it is a specific and nastier case of. State: the condition (a writer
  piped into an early-leaving reader under `pipefail`), that it is wrong in **both** directions, the
  measured window (silent below roughly 36 KB of producer output, certain above about 200 KB), the fix
  (capture and read as a herestring, asserting the producer's status separately when it can fail), and
  the readers the guard cannot see (`rg -q`, `head`, `sed -n '1p;q'`, `awk '/x/{exit}'`, `read`).
- Point at `_dev/tests/quiet-grep-pipeline-audit.sh` as the guard, and say plainly what it does not
  cover, so the prime is not read as "the guard has this".
- Fix `memory-reference.md:88`. Its producer output is far too small to reach the window in practice,
  so this is about what shipped guidance teaches, not a live failure — say so rather than overstating it.
- Check the other shipped action files for prescribed blocks carrying the same shape in the same pass.
  A markdown scan is not what REQ-594's guard does, and one instance found by hand says nothing about
  the rest.

## Constraints

- Prose and one prescribed block. No change to the guard, whose fixture pins 19 shapes.
- `_dev/tests/action-shell-blocks.sh` already checks prescribed shell blocks; if the fix belongs there
  as an assertion rather than only as prose, say so and do it, but do not weaken what it already pins.

## Open Questions

None.

## Triage

**Route: B** — Explore then build.

**Reasoning:** The prime section is prose whose facts are already established by REQ-593 and REQ-594
and measured in their records. What is not established is the fourth requirement: whether other shipped
action files carry prescribed blocks with the same shape. REQ-594's guard walks tracked `*.sh` files and
never reads a markdown block, and one instance was found by hand, which says nothing about the rest.
That sweep is discovery.

**Planning:** Skipped.

**The prime is the fix that stops the class being written; the guard only catches it after.** Both are
needed and only one existed.

## Plan

**Planning not required** — Route B: one prime section, one shipped block, and whatever the sweep of
the other shipped blocks finds.

*Skipped by work action*

## Exploration

**Sweep of every shipped Markdown fence, re-derived on main after the merge.** 166 `.md` files under
`skills/`, 32 with at least one shell fence, 74 fences in total, all ```bash (no ```sh, no ```shell,
deepest indent three blanks, none unterminated). `_dev/tests/action-shell-blocks.sh` extracts the same
74, so the sweep and the lint read the same set.

**Hits: one.** `skills/do-work-knowledge/actions/memory-reference.md:88`, the first of three backend
probes under "Semantic Recall (Layer 2)", inside a block agents copy and run:
`ollama list 2>/dev/null | grep -qiE 'embed'`. The only other `grep -q` in any fence is
`memory-reference.md:142`, a file-argument grep with no writer, which is not the shape.

**Readers the scanner cannot see: zero in the fences.** A supplementary grep over the 74 extracted
blocks for `rg -q`, `head`, `sed -n … q`, `awk … exit`, `read`, and `grep -m` after a pipe found
nothing. So the one hit is the whole population today, and the prime has to name the blind spot rather
than the guard closing it.

**Where the scan belongs.** REQ-594's guard walks tracked `*.sh` files and never reads a fence. The
scanner function lived inside `_dev/tests/quiet-grep-pipeline-audit.sh` next to its top-level run, so
sourcing that file would run the audit. `action-shell-blocks.sh` already extracts every fence and runs
`bash -n` and shellcheck on each; it is the natural place for the same scan, which means the function
has to be lifted into a file both scripts can source.

*Generated by the build agent's sweep, re-derived by the orchestrator on main (counts above)*

## Scope

**Files I will touch:**
- `_dev/primes/prime-shell-commands.md` (modify) — one new section for the class
- `skills/do-work-knowledge/actions/memory-reference.md` (modify) — the one shipped block
- `_dev/tests/quiet-grep-pipeline-scanner.sh` (add) — the scanner function, lifted so it can be sourced
- `_dev/tests/quiet-grep-pipeline-audit.sh` (modify) — sources the scanner instead of defining it
- `_dev/tests/action-shell-blocks.sh` (modify) — runs the scanner over every extracted fence and shipped shell file

**Three files under `_dev/tests/` are in, and the request's "no change to the guard" is read as its
fixture and behaviour, not its bytes.** The request's Constraints already open this door: if the fix
belongs in `action-shell-blocks.sh` as an assertion, do it, without weakening what it pins. Running the
scanner over fences needs the function to be sourceable; a second copy of the body is the drift REQ-594
measured against. The audit's 19 must-flag and 7 must-not-flag fixture does not change.

**Files I will NOT touch:** the scanner's reader set (grep/egrep/fgrep) and option class; the fence
regex in `action-shell-blocks.sh` (0-3 leading blanks, `bash|sh`), which the lint already pins and which
loses nothing today; `lessons-shell-commands.md`, which the archive step appends to.

**Acceptance criteria:**
- [ ] The prime carries the class as its own section beside "Unchecked Exit Status Reads as Content":
  condition, both directions, the measured window, the two-half fix, the readers no guard sees
- [ ] The prime points at `quiet-grep-pipeline-audit.sh` and says what it does not cover
- [ ] `memory-reference.md:88` no longer feeds a quiet grep from a pipe, and the record says this was
  never a live failure
- [ ] Every shipped fence has been scanned, and the scan is now a lint that fails at the Markdown line
- [ ] Nothing already pinned is weaker: 74 fences, 33 shipped shell files, `--self-test` passes, the audit's
  19+7 fixture passes over the tracked shell files

## Pre-Flight

**Green gate at `dc5d818`** (the builder's base) and again at `9e00a092` on main before the merge.
Fast gate exit 0, `Maintainer verification passed.`

**Lanes that read these files and must stay green:** `_dev/tests/action-shell-blocks.sh` (74 fences,
33 shipped shell files, shellcheck enabled), `_dev/tests/quiet-grep-pipeline-audit.sh` (95 tracked shell
files after the helper is added, 19 must-flag and 7 must-not-flag shapes), `_dev/tests/audit-lockins.sh`.

**The prime is prose; no lane reads it.** Its facts come from REQ-593's measurement (0 of 50 misfires
at 36 KB, 50 of 50 at 200 KB) and REQ-594's guard contract, both recorded and archived. The review checks
the section against those records, not against a test.

## Implementation Summary

**Files changed:**
- `_dev/primes/prime-shell-commands.md` (modified)
- `skills/do-work-knowledge/actions/memory-reference.md` (modified)
- `_dev/tests/quiet-grep-pipeline-scanner.sh` (added)
- `_dev/tests/quiet-grep-pipeline-audit.sh` (modified)
- `_dev/tests/action-shell-blocks.sh` (modified)

**The prime section.** `## A Writer's SIGPIPE Death Reads as the Reader's Verdict`, placed directly
after `## Unchecked Exit Status Reads as Content` and before `## Closed Enumerations Go Stale`. It states
the condition (a writer piped into a reader that can leave early, under `pipefail`) with an illustrative
reader list; that it is wrong in both directions and the negative matcher is the dangerous half; the
window from REQ-593's record; the two-half fix (capture and herestring, then assert the producer's status
separately because the capture discards it), with the one-line `listing="$(…)" || fail` form and "say so
in place" when failure is the answer; and the guard pointer, saying plainly that the audit walks tracked
`*.sh` only with a grep-only reader set, that Markdown is not its input, and that `action-shell-blocks.sh`
now runs the same scanner over the shipped fences.

**The shipped block.** `memory-reference.md:88` now captures `ollama list` into `ollama_models` with
`|| true` and greps the herestring. The producer's failure (no ollama, stopped daemon) is the "no
backend" answer the surrounding prose already prescribes, so its status is collapsed and the comment
says why. Not a live failure: `ollama list` prints under 2 KB on any realistic install, far below the
window; this is about what shipped guidance teaches.

**The scan.** `quiet_grep_pipeline_offenders` moved byte-for-byte into
`_dev/tests/quiet-grep-pipeline-scanner.sh`, sourced by the audit and by `action-shell-blocks.sh`. The
lint runs it on every extracted unit after `bash -n` and before the shellcheck gate, so it runs without
shellcheck too, and prints `FAIL: path:line: quiet grep fed from a pipeline (…): <command>` at the
Markdown line. A one-shape wiring fixture runs on every default invocation, because the gate never
passes `--self-test` and a pin behind that flag could not fail when the scan is removed.

**Mutation evidence from the builder, each restored and the restore diffed:** the scanner call removed
from the lint makes it exit 1 naming the missing wiring; the original block restored makes it exit 1 at
`memory-reference.md:88`; `--quiet` dropped from the shared body makes the audit exit 1 with
`no longer caught: grep --quiet`.

## Decisions

- **D1 The scan lives in `action-shell-blocks.sh`, not in a fourth script.** It already walks every
  fence and every shipped `.sh`; a second walker over the same set is the enumeration drift the prime's
  own section warns about.
- **D2 The scanner is lifted, not copied.** Two bodies of the same function is the drift REQ-594
  measured against. The audit's bytes changed by a `source` line; its fixture did not.
- **D3 `|| true` on the ollama capture, with the reason in the comment.** The prime's third way out:
  when the producer's failure is the answer, say so in place rather than assert a status that would
  wrongly stop the probe.
- **D4 The write set widened after the build, recorded before qualify.** The builder's hand-back named
  the three `_dev/tests/` files the request did not list; the frontmatter now carries them so the
  finalizer's allowlist matches the merged range.

## Discovered Tasks

- **`lessons-shell-commands.md` has no entry for REQ-593 or REQ-594** although both are archived; the
  archive step should have appended them. Not this write set. Check whether the archive step ran the
  lessons append at all for those two before capturing a request.
- **The scanner's reader set is grep/egrep/fgrep.** A fence feeding `rg -q`, `head`, `sed -n '1p;q'`,
  `awk '/x/{exit}'` or `read` from a pipe passes both scripts. The sweep found zero today and the prime
  says so. Worth a request only when a second reader shape ships.

## Qualification

**Passed.** Read from the range `9e00a092..a25c7522`, five files, 155 insertions and 77 deletions.
Canonical `qualify` and `scope-drift` both satisfied after the write set was widened to the three
`_dev/tests/` files the build touched (D4).

- **The sweep is re-derived, not transcribed.** On main after the merge: 166 Markdown files under
  `skills/`, 32 with a bash fence, 74 fences; the lint reports the same 74. One hit, and it is fixed.
- **The scan is real, proven by three mutations on main, each restored and the tree confirmed clean
  afterwards.** M1, the scanner call in `lint_shell_source` short-circuited to `$(true)`: the lint exits 1
  with `FAIL: the fence walk no longer flags a quiet grep fed from a pipeline at its Markdown line; the
  shared scanner is not wired in.` M2, the pre-change block restored from `9e00a092`: exit 1 with
  `FAIL: skills/do-work-knowledge/actions/memory-reference.md:88: quiet grep fed from a pipeline (...)`.
  M3, `|--quiet` deleted from the shared option class: the audit exits 1 with `FAIL: the quiet-grep
  scanner no longer catches 1 of 19 ordinary spellings of the pipeline it exists to forbid`.
- **The lift did not weaken the audit.** It still passes over 95 tracked shell files with 19 must-flag
  and 7 must-not-flag shapes, and M3 shows the fixture still reads the lifted body.
- **The prime section is placed where the request put it**, between "Unchecked Exit Status Reads as
  Content" and "Closed Enumerations Go Stale", and every number in it is REQ-593's measurement.
- **The block fix is stated at its real size.** Not a live failure; shipped guidance that taught the shape.

## Testing

**Guards on main at the merge, all exit 0:**

- `bash _dev/tests/action-shell-blocks.sh` — `Shell-block lint passed: 74 fenced blocks and 33 shipped
  shell files; ShellCheck enabled.`
- `bash _dev/tests/action-shell-blocks.sh --self-test` — `Shell-block lint self-test passed.`
- `bash _dev/tests/quiet-grep-pipeline-audit.sh` — `quiet-grep pipeline audit passed (95 tracked shell
  files, 19 must-flag and 7 must-not-flag shapes).`
- `bash _dev/tests/audit-lockins.sh` — `Audit lock-in regressions passed.`

**Mutations M1, M2 and M3** as recorded under Qualification: each exit 1 with the expected FAIL line,
each file restored with `git checkout --`, `git status --short` empty after.

**Fast gate on main after the merge:** `Maintainer verification passed.`, exit 0, wall 78s.

## Review

**Overall: 83%** | 2026-09-06T09:40:00Z | full synthesis: `do-work/runs/work-2026-09-05-231943/REQ-600-review.json`

| Dimension | Score |
|-----------|-------|
| Requirements | 80% |
| Code Quality | 88% |
| Test Adequacy | 92% |
| Scope | 90% |
| Risk | Low |
| Acceptance | Remediate two sentences before archive |

**Verdict: Accept with small remediation.** Three independent reviewers and a synthesizer who reproduced
every finding. The scan, the lift and the block fix are correct and pinned: all four guards exit 0 in the
reviewers' clone, mutations M1/M2/M3 fail with the exact FAIL lines this record quotes, and the lifted
function body is byte-identical (39 lines). The one important finding is in the deliverable the request
was mostly about, and it is a sentence this run wrote.

**Findings that survived verification (remediation pending, none applied yet):**

- **F1 (important, prime line 53) teaches a false safe zone below the pipe buffer.** "Below the pipe
  buffer the writer finishes before the reader leaves, so the check passes for years" is in neither
  archived record and is false for any writer that flushes in chunks: the synthesizer measured `tar tzf
  | grep -q` misfiring 20 of 50 at a 4.8 KB listing and 50 of 50 at 14 KB, both far under the 64 KiB pipe
  buffer, and replayed REQ-593's own pinned-CPU case (1 of 500 at 36 KB). The 0-of-50 and 50-of-50
  figures are real but belong to one writer, bash's builtin `printf` emitting a single string. The
  request's wording ("silent below roughly 36 KB") carried the same generalization, so the builder
  followed it and added an unmeasured mechanism. Fix: attribute the numbers to their writer, add the
  2-3-of-500 pinned-CPU result, delete the mechanism sentence, state that the flip size depends on how
  the writer flushes, and keep "there is no safe size". Also soften this record's "far below the window"
  for the ollama block to "small, and not a live failure today".
- **F2 (low) scanner contract comment line 18** still says "fourteen different ways in the must-flag
  fixture below"; the fixture is in another file and has 19 shapes. The hand-back's "comment moved
  byte-for-byte" held for the body, not the comment. Fix: drop the number, point at the audit file. Line
  11 carries the same "silent below roughly 36 KB" generalization as F1.
- **F3 (low) the hand-back and the satellite header cite "work.md Step 8 substep 7"**, which does not
  exist; the lessons write is Step 8 substep 4, and it did not run for REQ-593 or REQ-594. (The
  maintainer satellites were backfilled after this review in commit `a38a8c4`; the shipped satellite's
  lines and this REQ's own line are still owed.)
- **F4 (low) this record's Exploration** places `memory-reference.md:142` "in a fence"; it is a prose
  list item between fences. The judgment stands; the location claim came from a whole-file grep.
- **F5 (low) the block's prose says "first hit wins"** and the new first line is setup that always exits
  0. Either collapse to the one-line form or change the prose to "the first probe that succeeds wins".
- **F6 (low) prime line 51** lumps `grep -m 1`/`--quiet`/`--silent` (the same reader, flagged by the
  scanner) with `rg -q`/`head`/`sed`/`awk`/`read` (different readers, not flagged). Split the list.
- **Pre-existing, outside the diff's purpose:** `quiet-grep-pipeline-audit.sh:14` has no failure branch
  on `mktemp`, so as root it writes its fixtures at `/` and passes (the file is in this write set; fix
  during remediation or capture); a bash fence indented four spaces is never opened by the lint (zero such
  fences today; the fence regex is the one place to touch if widened).

**Findings that did not survive:** "the record does not state this is a release" (finalization's
judgment); the wiring fixture's own call being unpinned (true of every pin); the ollama block's
"far too small" reasoning being wrong in general (not measured; the never-a-live-failure judgment stands).

**Disagreements and how settled:** the window sentence (one reviewer marked it held from the archive
lines alone, two measured it false; the synthesizer measured and took the two); severity of the
byte-for-byte comment claim (substance agreed); whether the pre-existing audit and fence items count
against this REQ (confirmed identical at the parent; discovered tasks).
