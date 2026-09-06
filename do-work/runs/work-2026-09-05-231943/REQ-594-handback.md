# REQ-594 Hand-back — the quiet-grep guard leaves one probe file and covers the repository

REQ-594 is the request that generalizes the SIGPIPE fix: a writer piped into an early-leaving grep
under `set -o pipefail` reports the writer's SIGPIPE death as grep's verdict, and about 130 such sites
remained across the maintainer test tree after REQ-593 (which fixed eleven of them inside one file and
added a scanner that only read its own file).

**Branch:** `worktree-agent-REQ-594-quiet-grep-pipelines`
**Head:** `82e8845` (12 commits on base `895d573`, itself the pre-flight commit on `ceeea69`)
**Worktree:** `/home/user/skill-do-work-worktrees/worktree-agent-REQ-594-quiet-grep-pipelines`
Working tree clean. Nothing outside the worktree was written except this file. Three builders worked in
sequence on the same branch: steps 1-10, steps 11-13, step 14.

## The number this request exists to move

REQ-593's shipped scanner (`quiet_grep_pipeline_offenders`), run verbatim over
`git ls-files -z -- '*.sh'`:

| point | offending logical lines | files |
| --- | --- | --- |
| branch point `895d573` | 129 (135 grep invocations) | 22 of 93 |
| head `82e8845` | 0 | 0 of 94 |

The 94th tracked shell file is the new audit itself, which scans itself clean without a path exemption.

## What changed

25 files, `+488 / -249`.

**The new guard (3 files).**

- `_dev/tests/quiet-grep-pipeline-audit.sh` — new, executable, 237 lines. Holds four things: the
  scanner moved out of the updater probe with the comment fix applied, the repository walk, the
  two-sided fixture, and the doc comment block that travels with the function.
- `_dev/tests/contracts/probe-lanes.sh` — two lines registering the audit as a fast-tier probe, above
  the heavy block. The guard used to run only under `--heavy`; a violation now blocks an ordinary
  commit.
- `_dev/tests/update-script-behavior.sh` — 94 lines deleted, deletion only: the doc comment, the
  scanner, the two `verify_` wrappers that drove it, the fixture `mktemp`/`rm` pair and the two
  `record_failure` calls. `verify_output_matchers_read_greps_verdict` stays. It drives the two shipped
  matchers over a 256 KiB variable and is the only thing proving the mechanism on real data, which is
  not what the scanner does.

**The conversions (21 files, 129 sites).**

Fast tier: `contracts/core-checks.sh` 13 logical lines carrying 19 grep invocations,
`select-simple-reqs-behavior.sh` 15, `contracts/queue-kanban.sh` 2, `p50-estimator-determinism.sh` 1,
`audit-lockins.sh` 1.

Heavy tier: `staged-skills-contract.sh` 4, `install-suite-behavior.sh` 2, and 90 sites across 14 of the
18 files under `prescribed-shell-cases/` — `qualify.sh` 31, `audit-archive-timestamps.sh` 15,
`repair-req-timestamps.sh` 13, `generate-report-image.sh` 6, `generate-report-image-batch.sh` 6,
`publish-portfolio-summary.sh` 4, `atomic-download.sh` 3, `cleanup-req-reservations.sh` 3,
`protected-inventory.sh` 3, `capture-screenshot.sh` 2, and one each in `show-commit-diff.sh`,
`lexical-memory-recall.sh`, `install-memory-hooks.sh`, `architecture-report-preflight.sh`.

Shipped: `skills/do-work/tools/select-simple-reqs.sh` line 47, one line, and the reason this request is
a release.

## The conversion families, and what each one preserved

The rule the whole sweep is judged by: capturing a producer's output and reading it as a herestring
discards the producer's exit status, which silently narrows the assertion. Every site was read for what
it measures before it was converted.

**F1 — the producer replays an already-captured variable (about 110 sites).** These are
`printf '%s\n' "$captured" | grep -q PATTERN` where `$captured` was produced and status-asserted on an
earlier line. Converted to `grep -q PATTERN <<<"$captured"`. Nothing is lost because the pipeline's
writer was `printf` on a variable, and `printf` cannot fail in a way the assertion was reading. The
herestring appends a newline that `printf '%s'` did not; verified identical for every `grep -q`, `-qx`
and `-qxF` reader in the sweep, and it would not be safe for a byte-exact reader.

**F2 — the producer can fail, so its status is now asserted explicitly (7 sites).** These did not get a
bare herestring:

- `contracts/core-checks.sh:291-293`, producer `git diff --cached --name-status --no-renames -- .env`,
  split into three `elif` legs so "git broke" is never reported as "wrong cached metadata".
- `contracts/core-checks.sh:349`, the `printf | associate-files.sh` capture became its own `if` leg
  with its own failure message, with the three text reads following as `elif`.
- `prescribed-shell-cases/qualify.sh:368`, a negative assertion that the fixture index stayed clean.
  Captured with `|| fail_case 'could not read the fixture repo index'`. Mutation-tested both legs:
  pointing it at `git diff --name-only` reports "left the fixture repo index modified", making git
  itself fail reports "could not read the fixture repo index", where the old form was silent.
- `prescribed-shell-cases/audit-archive-timestamps.sh:170` and `:195`, three-stage
  `printf | grep -E | grep -q` chains. The filtered lines are captured into a variable and the deciding
  grep reads that variable, so the middle grep is no longer an early-leaving reader's writer. Finding no
  correction lines is the passing answer at both, so the filter's non-zero status is expected rather
  than asserted, said in a comment in place. Both mutation-tested against ids the correction lines
  really carry.
- `staged-skills-contract.sh:60` and `install-suite-behavior.sh:135`, `git check-attr export-ignore`
  positive tests for a forbidden state. A git that could not answer printed nothing, matched nothing and
  passed. Now the capture fails with "could not read the export-ignore attribute for /<path>".

**F3 — the 21 `find … -print -quit | grep -q .` leak checks.** Each captures into a named variable with
`|| fail_case '<case> could not search <where>'`, then tests emptiness. Asserting find's status is right
at all of them because every search root is created by the fixture before the check; this was checked by
running the five case files unmodified and looking for a `find:` diagnostic in the combined output, of
which there were none. Mutation-tested at `atomic-download.sh:14` in both directions: a planted leftover
file reports "leaked private scratch", a missing search root reports "could not search the fixture
tree", where the old form passed silently. One of the 21, `generate-report-image.sh:404`, is a readiness
poll inside a bounded wait loop rather than an assertion, so it reads only the captured path's
emptiness, with a comment saying why it differs from its twenty neighbours.

**F4 — the range extractions in `staged-skills-contract.sh`.** Lines 269 and 308 extracted the identical
`sed -n '/^## Routing/,/^## Dispatch/p' SKILL.md` range twice. Captured once into `core_routing_section`
above the first reader, guarded on sed's status and on the section being non-empty, then read twice as
herestrings.

## The two guards that were already broken

Both are `find … -print -quit | grep -q .` checks for a retired legacy runtime:
`staged-skills-contract.sh:82` and `install-suite-behavior.sh:142`. The census graded the second one
"None" risk and "structurally immune, the producer stops itself", reasoning that `-print -quit` leaves
no writer to kill. That is true for SIGPIPE and irrelevant to the actual bug: a find that reaches an
unreadable sibling directory exits non-zero **after** printing, and `pipefail` then reports that as "no
leftover files".

Reproduced with a `find` stub that prints a leftover path, writes a permission-denied line to stderr and
exits 1, driving the shipped guard block copied verbatim:

- shipped form: `fail_count=0`. The leftover file is sitting there and the suite calls the legacy
  runtime retired.
- replacement form: `fail_count=2` — "could not scan actions for leftover files" and "legacy root
  runtime must be retired at modular cutover: actions (/…/actions/leftover-runtime-file.md)". Two
  findings on purpose: the scan being unreliable and a leftover being present are different facts, and
  an `elif` form would have thrown away the leftover path in exactly the case the bug is about.

Controls with the real `find` and no stub: empty retired directory `fail_count=0`, directory holding a
leftover `fail_count=1`, path absent entirely `fail_count=0`. The new form is not simply louder.

Live mutation on the real files, with the plants removed by an `EXIT` trap: a planted
`interviews/leftover-runtime-note.md` makes `staged-skills-contract.sh` report "legacy root runtime must
be retired at modular cutover", and a planted `actions/leftover-runtime-note.md` makes
`install-suite-behavior.sh` report "fresh-install source still carries legacy root runtime".

The same class of hole was reproduced at the two `git check-attr` sites (stub git exiting 128: shipped
form 0 failures across all three export paths, replacement form 3) and at the two routing-section sites
(headings renamed so sed prints nothing and exits 0: shipped 0, converted 1; SKILL.md unreadable so sed
exits non-zero: shipped 0, converted 1).

Reproduction scripts:
`/tmp/claude-0/-home-user-skill-do-work/213e30ac-5958-56c8-9fd2-faaaaf9c4ea6/scratchpad/req594b2/repro-find-guard.sh`,
`repro-controls.sh`, `repro-routing-guard.sh` in the same directory.

## The scanner hole, and the one condition that closes it

The shipped scanner could be walked past with a comment line. `producer |`, then `# note`, then
`grep -q marker` is one running pipeline to bash — after a pipe bash keeps looking for the next command
past blank and comment lines — and the `| \` spelling behaves the same way. Confirmed by running it, not
by reading awk: `seq 1 5 |` / `# note` / `grep -c .` prints 5, and so does the backslash form.

The shipped rule was `joined_command == "" && /^[[:space:]]*#/ { next }`, so with a pipeline open it
concatenated the comment text and closed the logical line there, building
`tar tzf a.tgz |# an ordinary commentgrep -q marker`, where `commentgrep` no longer matches the grep
token boundary.

The fix is one condition, not two lines: `/^[[:space:]]*#/ { next }`. A comment line is invisible to
bash in that position whether or not a pipeline is open, so skipping it unconditionally matches bash for
both spellings.

Verified additive-only rather than assumed. Both scanner bodies were extracted verbatim (the shipped one
from `git show 895d573:_dev/tests/update-script-behavior.sh`, the new one from the audit file) and run
over every tracked shell file at the branch point: byte-identical output, 129 findings each. Over `HEAD`
both report 0. So nothing was hiding behind the hole in this tree, and the widened guard landed with no
new findings.

## The new audit, and the fixture that pins it

`_dev/tests/quiet-grep-pipeline-audit.sh` walks `git ls-files -z -- '*.sh'`, the same enumeration the
ShellCheck lane uses at `_dev/tests/maintainer-verify.sh:670`. A condition, not a file list: a shell
file added tomorrow is scanned without anything being registered. It fails loudly if git reports no
tracked shell files at all, and reports "could not scan <path>" rather than a clean verdict when a file
cannot be read — both legs were exercised (a planted offender in `audit-lockins.sh` reports the file and
line and exits 1; a tracked file moved aside reports "could not scan").

The fixture is two-sided and asserted by name, not by count, so a failure says which shape was lost.
Each fixture line carries a unique marker token and the checks read the markers out of the scan output,
so adding a shape does not renumber a table. Both halves are written from heredocs that interpolate the
pipe and hash characters (`pipe_character='|'`, `hash_character='#'`), so the audit's own bytes never
carry the offending shape and it scans itself clean with no path exemption. Both halves also parse under
`bash -n`, so they are real shell rather than text that resembles it.

14 must-flag shapes, 7 must-not-flag shapes. The must-not-flag half is not decoration: the request's own
census regex flags 2 of those 7.

### Mutation mapping

16 mutations of the scanner were applied one at a time to a copy of the audit and run against the
fixture. Every one was caught, and every one of the 21 shapes is covered by at least one mutation.

| mutation | fixture verdict |
| --- | --- |
| M01 drop `--quiet` from the option class | no longer caught: grep --quiet |
| M02 drop `--silent` | no longer caught: grep --silent |
| M03 drop `--max-count` | no longer caught: grep --max-count=1 |
| M04 narrow `[qm]` to `[q]` | no longer caught: grep -m 1 |
| M05 drop the `(e\|f)?` alternation | no longer caught: egrep -q |
| M06 anchor the reader to the start of its stage | no longer caught: LC_ALL=C grep -q, command grep -q, /usr/bin/grep -q |
| M07 revert the comment fix to the shipped condition | no longer caught: comment line after the pipe, comment line after the continuation |
| M08 drop the backslash-continuation branch | no longer caught: backslash continuation before the reader, comment line after the continuation |
| M09 drop the pipe-at-end-of-line branch | no longer caught: pipe at end of line with the reader on the next line, comment line after the pipe |
| M10 drop the logical-or neutralization | wrongly flagged: a logical or, which is not a pipe |
| M11 drop the reader word boundary | wrongly flagged: a helper command whose name merely ends in grep |
| M12 require the option directly after the reader name | no longer caught: grep -F -q, where -q is not the first option |
| M13 treat any line containing a hash as a comment | no longer caught: a trailing # note on the offending line itself |
| M14 delete the comment rule altogether | no longer caught: both comment shapes; wrongly flagged: a whole-line comment carrying the offending text |
| M15 scan the first pipe stage as well | wrongly flagged: file-argument reader, herestring reader, logical or, quoted pipe in the pattern |
| M16 widen the option class to `[qmc]` | wrongly flagged: grep -c, a reader that runs to EOF |

Mutation harnesses:
`/tmp/claude-0/-home-user-skill-do-work/213e30ac-5958-56c8-9fd2-faaaaf9c4ea6/scratchpad/req594b3/mutate.py`
and `mutate2.py`.

## The shipped file, and why this is a release

`skills/do-work/tools/select-simple-reqs.sh:47` changed from
`if printf '%s\n' "$selector_output" | grep -qx 'run_set: '; then` to
`if grep -qx 'run_set: ' <<<"$selector_output"; then`. That single line is what makes REQ-594 a release
rather than a maintainer-only change: everything else lives under `_dev/`, which `.gitattributes:48`
marks `export-ignore`.

Drawing the guard's scope around `_dev/tests/` to avoid the version bump would have been an exemption
list with exactly one entry, aimed at the one known instance of the defect the guard exists to forbid.

The conversion is lossless, verified rather than assumed. `do-work-cli.sh`'s status is captured into
`$selector_status` at lines 36-40 and asserted at 41-44, which prints the output and exits when it is
non-zero, so the matcher at line 47 was only ever replaying a variable and the pipeline had no producer
status left to lose. `printf '%s\n' "$v"` and `<<<"$v"` produce identical `od -c` bytes over five output
shapes (empty string, exactly `run_set: `, `run_set: REQ-1`, multi-line output ending in `run_set: `,
and output that already ends with a newline), and the `grep -qx` verdict agrees in all five.

Behavioral cover beyond the scanner: `_dev/tests/select-simple-reqs-behavior.sh` probe T5 ("an empty
selection is a normal answer") drives this exact branch. It passes after the conversion, and mutating
the converted pattern makes it fail with "T5 empty selection: must state that nothing qualifies".

**The release apparatus is not in this branch and is the orchestrator's to make**: no `VERSION` bump, no
`CHANGELOG.md` entry, no mirror to `skills/do-work/CHANGELOG.md`. No gate check ties a shipped-file
change to a version bump, so the branch is green as it stands; the only changelog rule the gate enforces
is that the two changelogs stay byte-identical, and neither was touched.

## Evidence

**Fast gate, run from the worktree against the worktree, at head `82e8845`:** exit 0.

```
DO_WORK_GATE_ROOT=/home/user/skill-do-work-worktrees/worktree-agent-REQ-594-quiet-grep-pipelines \
  bash /tmp/claude-0/-home-user-skill-do-work/213e30ac-5958-56c8-9fd2-faaaaf9c4ea6/scratchpad/gate.sh
```

Last lines: `maintainer-verify: gate wall 99s` / `Maintainer verification passed.` The audit reports
`quiet-grep pipeline audit passed (94 tracked shell files, 14 must-flag and 7 must-not-flag shapes)` and
`test-file duration: quiet-grep-pipeline-audit.sh 2s (limit <30s)`. Standalone it runs in 0.20s; the 1-2s
the batch reports is whole-seconds rounding inside the parallel probe launcher. Log:
`/tmp/claude-0/-home-user-skill-do-work/213e30ac-5958-56c8-9fd2-faaaaf9c4ea6/scratchpad/req594b3/gate-fast-final.log`.

**Heavy lanes**, each `--heavy-lane <id>` through the same wrapper, all at head `82e8845`:

| lane | exit | last line |
| --- | --- | --- |
| `staged-skills` | 0 | `test-file duration: _dev/tests/staged-skills-contract.sh 23s (limit none (heavy))`, after `Prescribed shell script behavior probes passed (110 named script cases across 18 per-script files).` and `staged skills contract: PASS` |
| `updater` | 0 | `test-file duration: _dev/tests/update-script-behavior.sh 72s (limit none (heavy))`, after `update-script behavior probes passed.` |
| `installer` | 0 | `test-file duration: _dev/tests/install-suite-behavior.sh 16s (limit none (heavy))`, after `suite installer behavior probes passed.` |

Logs: `gate-staged-skills.log`, `gate-updater.log`, `gate-installer.log` under
`/tmp/claude-0/-home-user-skill-do-work/213e30ac-5958-56c8-9fd2-faaaaf9c4ea6/scratchpad/req594b3/`.
The `staged-skills` lane is the one that reaches all 90 prescribed-shell-case conversions, which a fast
run never touches.

**ShellCheck** in the lane's own form (every tracked `*.sh` in one invocation from the repository root,
94 files, `--severity=warning`): exit 0. Per-file linting is not the right check for the
`prescribed-shell-cases` files, which emit SC2154 for harness variables their sourced preamble sets.

**Scanner over the head:** 0 offending logical lines, with both the shipped scanner and the
comment-fixed one.

## Found and not fixed

1. **The reader set is grep, egrep and fgrep only.** `rg -q`, `head`, `sed -n '/x/{p;q;}'`,
   `awk '/x/{exit}'` and `read` are the same defect with a different early-leaving reader and the guard
   does not see them. No instance exists in the tree today; a wider reader set is a separate request.
2. **The scanner has no heredoc awareness.** A heredoc body whose text happens to spell the shape would
   be reported as an offender even though it is data. No instance today, and the audit's own fixtures
   avoid it by interpolating the pipe and hash characters, which is the workaround a future case would
   copy.
3. **The scanner splits on a bare pipe**, so a `|` inside a quoted grep pattern creates a phantom stage.
   Zero false positives across all 94 files, and the shape is pinned as a must-not-flag fixture line, but
   the parser is textual and a sufficiently strange line could still fool it.
4. **`find`'s status is asserted at all 21 leak checks because every search root exists at check time.**
   If a future case searches a root that may legitimately be absent, that assertion turns a correct pass
   into a failure. The pattern to copy is at `atomic-download.sh:14`, not blindly.
5. **Three converted sites sit inside `if false; then` blocks** (`generate-report-image.sh:404` and
   `:424`, `generate-report-image-batch.sh:425`) — disabled fixtures whose behaviour is now covered by Go
   tests. They are scanned as text, so they had to be converted, but no run exercises them; their only
   evidence is `bash -n`, ShellCheck and the scanner.
6. **A duplicate producer was deleted at `staged-skills-contract.sh:770`**, outside the defect class and
   outside the brief. It extracted the same routing range a third time, unguarded, feeding `grep -cF`
   where an empty section reads as the passing count 0 — the same vacuous pass, one classification away.
   It now reads the guarded capture. Flagged here because it was not in the census and nothing else will
   report it.
7. **The audit walks index paths and reads working-tree files.** A new shell file is not scanned until it
   is added to the index, and a tracked file deleted from the working tree is reported as "could not
   scan" rather than skipped. Both are the intended behaviours, stated so a future reader does not
   rediscover them as bugs.
8. **The request file's own checkboxes and acceptance criteria are untouched**, as is the release
   apparatus in the section above. Both belong to finalization, not to the builders.
