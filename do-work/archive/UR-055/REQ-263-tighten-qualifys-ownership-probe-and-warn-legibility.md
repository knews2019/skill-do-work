---
id: REQ-263
title: Tighten qualify's ownership probe and make its WARN legible
status: completed
created_at: 2026-08-18T19:52:15Z
status_changed_at: 2026-08-18T20:55:14Z
claimed_at: 2026-08-20T11:58:07Z
completed_at: 2026-08-20T12:12:48Z
commit: fd9b489
route: B
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-08-20T11:58:07Z
  basis:
    - trivial short-circuit
user_request: UR-055
addendum_to: REQ-254
domain: general
review_generated: true
effort_estimate: trivial
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec: bug-fix
depends_on: [REQ-300]
maintenance: true
write_set:
- skills/do-work/tools/checks/qualify.sh
- _dev/tests/prescribed-shell-cases/qualify.sh
---

# Tighten Qualify's Ownership Probe and Make Its WARN Legible

## What

REQ-254's ownership condition ("printed output belongs to whoever owns the process exit") is implemented as a whole-file grep for exit-idiom text, which is weaker than the condition as stated. Reproduced by its review, three ways: adding `sys.exit(0)` in the same diff as a debug print flips FAIL to WARN; a pre-existing `__main__`-guarded exit makes every debug print in a dual-use module WARN; and a library file whose **docstring merely says "exit 1 on failure"** WARNs. Also: the WARN branch omits the matched lines the FAIL branch prints, so "confirm from the diff" costs a manual dig.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Reproduce all four REDs against pre-change code first (REQ-243's lesson), then take both directions the REQ permits: judge ownership at the base revision, and narrow the bare-`exit` form to a statement shape. Print the WARN's matched lines. Add the untracked walk as a separate loop so Check 1 is untouched.
- [x] **[APPLY]:** `qualify.sh` — ownership base, `file_owns_process_exit`, narrowed `process_exit_regex`, WARN prints its lines, untracked walk. `_dev/tests/prescribed-shell-cases/qualify.sh` — seven cases.
- [x] **[UNIFY]:** `git diff --stat` reviewed; `shellcheck --severity=warning` clean on both files; `maintainer-verify.sh` exits 0; suite 76 → 83 cases. No debug artifacts — the only added `print(`/`console.log` strings are inside fixture heredocs the new cases feed to the checker on purpose, and the checker's own verdict on them is the assertion. Per-file: **`qualify.sh`** — shellcheck clean; `changed_file_list` untouched so Check 1's `(modified)`/`(deleted)` branches are byte-for-byte as before; `grep -c` not `grep -q` at the end of the new pipe, per the SIGPIPE-plus-pipefail trap this file already documents; `local`-declared-then-assigned in the helper, per the prime; helper defined above its first call site (REQ-276). **`prescribed-shell-cases/qualify.sh`** — shellcheck reports only the two pre-existing structural `SC2154`s for `repo_root`/`fixture_root`, which come from the sourced harness and are suppressed by the `# shellcheck source=` directive under the suite's own lint; 3 → 10 cases, 0 failures; each new case cleans up its fixture, and the range case removes the prior case's untracked files before committing so they cannot leak into its range.

## Context

REQ-254 review, Important finding 1 (gate: trivial — the WARN's confirm-from-diff instruction plus Step 6.3's judgment contract mitigate; threat model is forgetfulness, not adversarial builders) plus its folded Minor (WARN legibility). Created `pending-answers` per the generation-≥2 depth stop. The categorical fix (exit-added-in-same-diff ⇒ FAIL) is known-wrong: a legitimately new checker adds its prints and its exit in one diff.

## Requirements

- The ownership probe moves toward code-shaped exit occurrences (not docstring/comment text) and/or base-revision ownership for pre-existing files — direction is builder latitude; the boundary that ships is pinned by a lock-in either way (REQ-250's lesson: pin the documented limitation with a fixture that can fail, e.g. the docstring-"exit 1" case).
- The WARN branch prints the matched lines exactly as the FAIL branch does.
- `bash _dev/tests/maintainer-verify.sh` exits 0.

## Open Questions

- [ ] REQ-254's review found the ownership probe satisfiable by non-semantic bytes (same-diff exit, `__main__` guard, docstring prose) and the WARN branch less legible than FAIL. Should I process this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it — the WARN channel plus orchestrator judgment is mitigation enough.

**Answered [2026-08-18]:** User approved via `do-work clarify` — queued for a future work run.

## Addendum (2026-08-20)

From the 2026-08-20 consumer review (UR-062), filed P2: "Scan untracked implementation files for
output artifacts — `.claude/skills/do-work/tools/checks/qualify.sh:210`". Folded in here rather than
captured separately because its remedy writes exactly this REQ's write set.

The gap, verified against the script: in serial mode `changed_file_list` is built from
`git diff --name-only` plus `--staged` (`qualify.sh:104-111`), neither of which lists untracked
files. Step 6.3 runs **before** this REQ's commit, so a new source file the builder never staged is
invisible to the debug-artifact gate — and Check 1's `(new)` branch tests only `[ -f "$file_path" ]`
on disk (`:122-124`), so the file passes that check and is then never scanned. A checked `[UNIFY]`
can ship with leftover instrumentation in it.

The reported gap is **wider than the report states**: the `debugger|TODO|FIXME` scan at `:196-200`
reads the same diff stream, so it misses untracked files too — not only the
`print(`/`console.log` half at `:210-224`. Both halves need the fix.

Worktree dispatch mode is unaffected: `$diff_range` reads committed work.

### Added requirements

- Both artifact scans see untracked, non-ignored files in serial mode —
  `git ls-files --others --exclude-standard` is the intended source. The `do-work/` path exclusion
  applies to them exactly as it does to the diffed paths.
- An untracked file has no diff, so it is scanned whole-file rather than through
  `git diff -- "$path"`. Every line is "added" for a file that has never been committed, so this does
  not weaken the added-lines-only contract.
- The ownership condition (`:220-224`) still decides FAIL vs WARN for an untracked file, on the same
  terms as a tracked one.
- Adding untracked paths must not change Check 1's `(modified)` / `(deleted)` WARN behavior — if a
  single shared list would, build the artifact scan's set separately.
- Worktree dispatch mode (`$diff_range` set) is unchanged.

### Added Red-Green Proof

**RED prompt/case:** In a repo with a REQ whose `[UNIFY]` box is checked, create a new **unstaged,
untracked** library file containing a `console.log` line and a `TODO` comment, list it as `(new)` in
the Implementation Summary, then run `qualify.sh <req-file>`. Today it exits 0 with no FAIL for
either artifact class.
**Why RED now:** Both scans read only working and staged diffs, and untracked files appear in
neither.
**GREEN when:** That same run FAILs, naming the untracked file and its matched lines — while a
`.gitignore`d file in the same tree is still not scanned, and worktree-dispatch mode
(`DO_WORK_DIFF_RANGE` set) behaves exactly as before.
**Validation:** Inferred during capture — read from the script's list construction and its two scan
sites.

---

## Triage

**Route: B** - Medium

**Reasoning:** The defects are named and reproduced, but the REQ grants latitude on direction ("code-shaped exit occurrences **and/or** base-revision ownership") and rules one obvious fix out as known-wrong, so which boundary ships had to come from reading the script and running the failures. The addendum adds a second surface (untracked files) whose entry point also needed locating.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

**Every stated RED was reproduced against pre-change code before any edit** — `_dev/primes/prime-shell-commands.md` § Lessons, REQ-243: a "RED" that was already green means the premise moved. All four were red, in a scratch fixture repo with a library module, a docstring module, a dual-use module, and a real shell checker:

| Case | Pre-change behavior | Verdict |
|---|---|---|
| `sys.exit(0)` added in the same diff as `print(raw)` | `WARN`, exit 0 | red as stated |
| `__main__`-guarded `sys.exit(0)`, library-level `print()` added | `WARN`, exit 0 | red as stated |
| docstring reading "the caller should exit 1 on failure", `print()` added | `WARN`, exit 0 | red as stated |
| untracked `.js` with `console.log` **and** `TODO`, `[UNIFY]` checked | exit 0, no FAIL at all | red as stated |

The WARN in all three of the first cases printed **no matched lines**, confirming the folded Minor by observation rather than by reading.

**Where the two defects live.** Ownership: `qualify.sh:219`, `grep -qE "$process_exit_regex" "$changed_path"` — a whole-file grep against the **post-change working copy**, which is what lets an exit added in this diff change the verdict for this diff. Untracked blindness: `changed_file_list` at `:104-111` is built from `git diff --name-only` plus `--staged`, neither of which lists untracked files, and both artifact scans read from that. Check 1's `(new)` branch tests only `[ -f "$file_path" ]`, so an untracked file passes that check and is then never scanned — the addendum's reading of the report, confirmed.

**The addendum's claim that the gap is wider than reported is correct.** The `debugger|TODO|FIXME` scan at `:198-200` reads the same diff stream as the `print(`/`console.log` half, so both were blind to untracked files, not just the half the consumer review named.

**The right untracked source is already documented.** `_dev/primes/prime-shell-commands.md` trap 1: `git status --porcelain` collapses a wholly-untracked directory into one `?? dir/` row, so a per-file consumer needs `-uall` or `git ls-files --others --exclude-standard` — and the latter "also drops correctly-ignored paths, so it doubles as the untracked ignore filter." That is exactly the addendum's requirement, so no design decision was needed here.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work/tools/checks/qualify.sh` (modify) — ownership base, code-shaped exit probe, WARN line printing, untracked walk
- `_dev/tests/prescribed-shell-cases/qualify.sh` (modify) — lock-in cases for each boundary that ships

**Files I will NOT touch:**
- `skills/do-work/actions/work.md` Step 6.3 — the prose describes the ownership *condition*, which is unchanged; only its implementation got sharper
- `_dev/tests/prescribed-shell-scripts-behavior.sh` — the runner, which REQ-258 reduced to dispatch; a case-adding REQ writes its script's case file
- `tools/checks/scope-drift.sh` — a sibling check, out of scope here

**Acceptance criteria (restated from REQ):**
- [ ] The ownership probe moves toward code-shaped exit occurrences and/or base-revision ownership, and the boundary that ships is pinned by a lock-in either way
- [ ] The WARN branch prints the matched lines exactly as the FAIL branch does
- [ ] Both artifact scans see untracked, non-ignored files in serial mode; a `.gitignore`d file in the same tree is still not scanned
- [ ] An untracked file is scanned whole (no diff exists for it)
- [ ] The ownership condition still decides FAIL vs WARN for an untracked file
- [ ] Check 1's `(modified)`/`(deleted)` WARN behavior is unchanged
- [ ] Worktree dispatch mode (`DO_WORK_DIFF_RANGE` set) is unchanged
- [ ] `bash _dev/tests/maintainer-verify.sh` exits 0

## Pre-Flight

**Git:** ✓ working tree clean outside `do-work/`
**Tests baseline:** ✓ passing (`bash _dev/tests/maintainer-verify.sh`, exit 0, `launched: true`)
**Dependencies:** ✓ go1.26.1, ShellCheck 0.11.0, `just` present

*Checked by work action*

## Decisions

- **D-01 — Ownership is judged at the base revision, not at the post-change working copy.** DECIDE & STATE. The REQ offers this as one of two permitted directions and rules the categorical alternative out as known-wrong. Base is `HEAD` serially and the range's lower bound in worktree mode; a path absent there is new and is judged on its own content, which is exactly what keeps a legitimately new checker — prints and exit arriving together in one new file — a WARN. Pinned by two cases: `REQ-953` (same-diff exit must FAIL) and `REQ-954` (new checker must WARN). Reasoning: it closes the defect without needing to distinguish "added an exit because this is a checker" from "added an exit that happens to sit near a debug print", because the question it asks is about the file's identity *before* the change.

- **D-02 — A `__main__`-guarded exit keeps its WARN verdict; what changes is that the WARN is now legible.** ESCALATE. The REQ lists the dual-use-module case as one of three reproduced defects, and this decision declines to change its verdict. A module with `if __name__ == "__main__": sys.exit(...)` genuinely *is* runnable and genuinely does own its process exit when run that way, so a print inside it is genuinely ambiguous — which is what WARN is for; FAIL would fire on a legitimate CLI's own output. On inspection the practical harm was not the verdict but that the WARN said "confirm from the diff" while withholding the lines it had already found, so the fix for this case is the folded Minor. **Value:** no false FAIL on any dual-use CLI, and the ambiguity is surfaced with the evidence attached instead of sending the reader to dig. **Risk:** a debug print buried in a dual-use module's library half still reads as WARN rather than FAIL, so it relies on the orchestrator reading the now-printed lines; reversible in one branch if that proves too weak (medium reach, fully reversible).

- **D-03 — Code-shaped means statement-shaped, and the residual is pinned rather than chased.** DECIDE & STATE. Full comment-and-string awareness needs a per-language parser, so the narrowing is two cheap rules: drop full-line comments, and require the bare `exit N` form to begin a statement *and* terminate right after its status. That rejects every prose shape the review found ("the caller should exit 1 on failure", "On a bad row, exit 1", "returns exit 1 if missing") while keeping real shell (`exit 1`, `exit 1; fi`, `if …; then exit 2; fi`, `exit $?`). What survives is a docstring line consisting of nothing but the idiom and its status, which is textually identical to a shell statement. Recorded as a documented residual in the script and pinned by `REQ-956`, whose failure message says to update the note — REQ-250's lesson, pin every documented limitation with a fixture that can fail.

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/checks/qualify.sh` (modified)
- `_dev/tests/prescribed-shell-cases/qualify.sh` (modified)

**What was done:** Ownership of a file's process exit is now resolved against an `ownership_base` (`HEAD` serially, the merge range's lower bound in worktree mode) through a new `file_owns_process_exit` helper, so an exit idiom added in the same diff no longer re-labels library code as a reporter; a path absent at the base is judged on its own content. `process_exit_regex`'s bare-`exit` clause was narrowed to a statement shape and the helper drops full-line comments before matching, so docstring and comment prose no longer counts as ownership. The output-primitive WARN now prints its matched lines with the same `head -10 | sed` treatment the FAIL branch uses. A separate serial-mode walk over `git ls-files --others --exclude-standard` runs both artifact scans whole-file over untracked, non-ignored paths, skipping `do-work/` per path and skipped entirely when `DO_WORK_DIFF_RANGE` is set. Seven lock-in cases added, each mutation-proven able to fail.

## Qualification

`tools/checks/qualify.sh` on this REQ: **one FAIL, overridden with evidence; two WARNs, both correct.**

**FAIL — `[UNIFY]` is checked but the diff adds debug artifacts.** Five flagged lines, all in `_dev/tests/prescribed-shell-cases/qualify.sh`, none in shipped code:

```
106:+  '  // TODO: handle the empty case' '  return raw;' '}' \
108:+printf '%s\n' 'console.log("ignored debug");' '// TODO: ignored marker' \
111:+  && fail_case 'qualify untracked case passed an untracked file carrying a console.log and a TODO'
115:+  || fail_case 'qualify untracked case missed the TODO — the unfinished-marker scan still does not read untracked files'
135:+printf '%s\n' 'console.log("stray");' '// TODO: stray marker' > "$qualify_repo/src/stray_helper.js"
```

Three are **fixture payloads** written into a throwaway repo for the checker to detect — the marker *is* the test input, and the assertion is that qualify finds it. Two are `fail_case` message strings that merely contain the word. `git diff --name-only -- . ':(exclude)do-work/'` returns only the two declared files, and `skills/do-work/tools/checks/qualify.sh` contributes zero flagged markers.

**This is the established convention, not new sloppiness.** `git show HEAD:_dev/tests/prescribed-shell-cases/qualify.sh` already carries the identical pattern at `:67` — `printf '# TODO: tighten the site regex\n' >> "$qualify_repo/site-checker.sh"` — written by REQ-254 for exactly this purpose. `[UNIFY]` stays checked; the box is not un-checked for a marker that is the fixture.

**Worth stating: this false FAIL is guaranteed to recur, and it is a different class from REQ-301's.** REQ-301 covers a *moved* line read as an added one. This is a *genuinely new* line whose content is the test's payload, in the case file of the very checker that detects it. Any REQ that adds a case here trips it — and `_dev/tests/prescribed-shell-cases/qualify.sh` is in the write set of **REQ-264 and REQ-301** as well, so the next two builders will both hit it. Not filed as a follow-up (the 2026-08-20 stopping policy defaults generation-two findings to no new task) and **not folded into REQ-301 either**, because REQ-301 was approved at a stated scope and quietly widening it is not this REQ's call. Surfaced in the hand-back with a recommendation instead. Note that a blanket `_dev/tests/` exclusion would be the wrong fix: `_dev/primes/prime-shell-commands.md` names that exact anti-pattern — "a blanket skip/exclude list applied *before* a check silently neuters any check meant to fire inside the excluded set."

**WARN ×2 — both correct, and both are this REQ's own feature working.** `qualify.sh` and the case file each gained `print(`/`console.log` strings and each owns its process exit, so both read as the file's own reporting. The lines are now printed under each WARN, which is the folded Minor demonstrating itself on the very commit that adds it.

**`scope-drift.sh`: `OK: Implementation Summary matches the Scope declaration`** (exit 0). No undeclared touch, no unused declaration.

**Judgment checks:** *(2) Substantive* — `qualify.sh` gains a helper, a rewritten regex, an ownership base, and a 30-line scan block; the case file goes 3 → 10 cases. *(3) Requirements traced* — every acceptance criterion has a named case or a named run below. *(6) Flowing* — the ownership probe reads real content from `git show`/`cat`, and the untracked walk reads real files; both are proven by the mutation tests, which could not fail if either were stubbed.

## Testing

**Tests run:** `bash _dev/tests/prescribed-shell-cases/qualify.sh`, then `bash _dev/tests/maintainer-verify.sh`
**Result:** ✓ `qualify: 10 cases, 0 failures`; ✓ `maintainer-verify.sh` exit 0; suite total 76 → **83 named script cases across 17 per-script files**

**Red-green validation** — every RED was observed on pre-change code first, then re-run:

- `REQ-953` same-diff exit: ✗ `WARN`, exit 0 before → ✓ `FAIL: … never ends its own process`, exit 1 after
- `REQ-955` docstring "the caller should exit 1 on failure": ✗ `WARN`, exit 0 before → ✓ `FAIL` naming the file, exit 1 after
- `REQ-957` WARN legibility: ✗ WARN printed no matched lines before → ✓ prints them after
- `REQ-958` untracked file with `console.log` + `TODO`: ✗ exit 0 with no FAIL at all before → ✓ two FAILs naming the file and both line sets after; the `.gitignore`d sibling stays unscanned
- `REQ-954` new checker, prints and exit together: ✓ `WARN` both before and after — a **guard against over-fixing**, since the categorical rule the REQ rules out would have turned this into a FAIL
- `REQ-959` worktree dispatch mode: ✓ untracked files unscanned both before and after
- `REQ-956` documented residual: ✓ still WARN, pinned so the stated limit cannot drift silently

**Mutation testing — every new case proven able to fail.** A green case that cannot fail is read as coverage while pinning nothing (`_dev/primes/prime-shell-commands.md` § Lessons, REQ-257), so each boundary was reverted in turn:

| Mutation | Result |
|---|---|
| ownership reads the post-change working copy again | `REQ-953` fails (2 assertions) |
| bare-`exit` clause loosened back to the prose-permissive form, parenthesised forms kept | `REQ-955` fails |
| the same loosening, whole regex | `REQ-954`, `REQ-956`, and both reporter cases fail |
| WARN stops printing its matched lines | `REQ-957` fails |
| the untracked walk removed | `REQ-958` and `REQ-959` fail |

**Boundary table for the narrowed probe**, produced by running the shipped regex over candidate lines — recorded because "code-shaped" is a claim about exactly this set:

```
OWNS  exit 1                              -  The caller should exit 1 on failure.
OWNS      exit 0                          -  On a bad row, exit 1
OWNS    exit 1; fi                        -  Raises and the caller should exit 1
OWNS  if [ -z "$x" ]; then exit 2; fi     -  returns exit 1 if missing
OWNS  exit $?                             -  # exits 1 on failure
OWNS  sys.exit(0) / raise SystemExit(…) / process.exit(1); / os._exit(…)
OWNS      exit 1        <-- the documented residual: line-initial prose, pinned by REQ-956
```

**New tests added:** seven cases in `_dev/tests/prescribed-shell-cases/qualify.sh` — `REQ-953` through `REQ-959` as tabled above.

**Existing tests updated:** none. The three REQ-254 cases pass unchanged, which is the evidence that the reporter exemption itself was not weakened.

*Verified by work action*

## Review

**Overall: 91%** | 2026-08-20T12:11:27Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 90% |
| Test Adequacy | 95% |
| Scope | 95% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- None

**Minor findings:** 2 (report only)
- `skills/do-work/actions/work.md:446` glosses this script as verifying "P-A-U box audit + debug artifacts **in the diff**". After this REQ the serial-mode artifact scans also read untracked files whole, so the gloss is now less than the whole input set. Left unedited on purpose: `work.md` Step 6.3 was explicitly declared out of scope, the authoritative statement of what the scans read belongs in the script and was updated there instead, and the gloss is not *wrong* for the tracked majority it summarizes. Reported rather than fixed so the choice is on the record — a two-word edit outside a declared Scope is still scope drift.
- `D-02` declines to change the verdict on one of the three defects the REQ names (the `__main__`-guarded dual-use module). The reasoning is recorded and the practical harm is closed by the WARN-legibility half, but a reader comparing the REQ's finding list against the diff will see two verdicts changed and one deliberately kept.

**Restatement sweep:** the diff changes what a gate reads and how its WARN reports, so the sweep asked who else states either. Nothing outside the two declared files restates the ownership condition or the exit-idiom set — `grep -rn "owns its process exit\|ends its own process"` over `skills/` and `_dev/` returns no `.md` hit. The one gloss that describes the scans' input is the `work.md:446` Minor above. The script's own header, which is the contract surface, now states the untracked behavior and the base-revision ownership rule.

**Acceptance:** Pass — all eight acceptance criteria met with named evidence: seven new cases, all mutation-proven able to fail; `maintainer-verify.sh` exit 0; suite 76 → 83 cases; Check 1 and worktree mode both proven unchanged by dedicated cases.

**Suggested testing:** 2 items
- No case covers an untracked file inside a wholly-untracked *directory*. `git ls-files --others --exclude-standard` lists those individually (this is exactly the prime's trap 1, and the reason that source was chosen over `git status --porcelain`), so the behavior is expected to hold — but nothing pins it, and a future change to the enumeration source would not be caught.
- No case covers a path with a leading dash or a space. The new `grep` calls pass the path after `--`, so it should hold; unpinned.

**Scope 95%:** the script's header comment was extended beyond the change strictly needed to satisfy the requirements. It is the same file and the right home for the contract, so this is inside the declared Scope, but it is a deliberate addition rather than something the REQ asked for.

**Follow-ups created:** None; **sweeps appended to:** None

## Lessons Learned

**What worked:** Running every stated RED before writing anything, and mutation-testing every new case afterwards. The first pass caught that the residual I had *documented* was not the residual that actually exists — my fixture used "On a bad row, exit 1", which the narrowed probe correctly rejects, so the case failed for the right reason and forced the note to be rewritten to the true boundary (a line-initial prose line ending right after the status). A green suite would have shipped a residual note that described the wrong limit.

**What didn't:** Two fixture bugs of my own, both from state leaking between cases in a shared fixture repo. The range case FAILed because the previous case's untracked offenders were still lying in the tree when `fixture_repo_commit_all` ran, sweeping a `TODO` into the range under test. In a case file where every block shares one repo, a case that *creates* untracked files owes the next case a cleanup before any commit — the fixture root being per-file is not per-case isolation.

**Worth knowing:** Adding a case to this particular case file is guaranteed to make `qualify.sh` FAIL the REQ that adds it, because the fixtures must contain the very markers the checker detects. It happened here, REQ-254 hit it before, and `_dev/tests/prescribed-shell-cases/qualify.sh` is in the write set of REQ-264 and REQ-301 too. Override it with the `git diff --name-only` evidence and leave `[UNIFY]` checked; do not reach for a `_dev/tests/` exclusion, which is the anti-pattern the prime's trap 2 names. Also: `grep -q` at the end of a pipe under `pipefail` reports "no match" when it actually matched and SIGPIPE'd the upstream — this file documents the trap and the new helper had to route around it with `grep -c`.

## Orientation

The debug-artifact gate now asks a sharper question and reads a wider input. It judges whether a file owns its process exit **at the revision the change started from**, so an exit added alongside a debug print no longer excuses it, and it only counts exit idioms that look like statements rather than prose in a docstring or comment. In serial mode it also reads untracked, non-ignored files whole, closing the hole where a new source file the builder never staged was scanned by neither artifact half. Every WARN it raises now prints the lines it found. Lives in the Step 6.3 qualification gate (`_dev/primes/prime-shell-commands.md`).

Not `[MAP CHANGED]` — one gate got more accurate inside its existing contract; no new checklist item, no new field, no caller change. Staleness spot-check on `_dev/primes/prime-shell-commands.md`: every referenced path resolves, and its trap 1 (`git ls-files --others --exclude-standard` as the per-file untracked source that doubles as the ignore filter) is what this change is built on, so the prime is not stale — it is load-bearing here.
