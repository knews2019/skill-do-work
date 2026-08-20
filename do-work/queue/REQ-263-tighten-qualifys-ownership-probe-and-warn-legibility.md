---
id: REQ-263
title: Tighten qualify's ownership probe and make its WARN legible
status: pending
created_at: 2026-08-18T19:52:15Z
status_changed_at: 2026-08-18T20:55:14Z
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
- [x] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
