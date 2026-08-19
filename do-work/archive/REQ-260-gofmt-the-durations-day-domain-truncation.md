---
id: REQ-260
title: Run the Go formatter as part of the canonical verify
status: completed
completed_at: 2026-08-18T22:36:02Z
commit: 307e146
claimed_at: 2026-08-18T21:16:24Z
route: A
created_at: 2026-08-18T18:41:26Z
status_changed_at: 2026-08-18T20:59:31Z
user_request: UR-051
addendum_to: REQ-251
domain: general
effort_estimate: normal
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: false
write_set:
- _dev/tests/maintainer-verify.sh
estimate:
  p50_active_minutes: 15
  confidence: medium
  calculated_at: 2026-08-18T21:17:28Z
  basis:
    - Route A
    - 1-file write set
    - 3 acceptance criteria
    - full-suite verification
---

# Run the Go Formatter as Part of the Canonical Verify

## What

`bash _dev/tests/maintainer-verify.sh` runs `go vet` but never the Go formatter, so a formatting slip lands and survives silently — one did, in the board tool's day-domain truncation expression, and was only caught by a builder reading adjacent code. That specific slip is already fixed (REQ-252 corrected it in passing), so this REQ is the rule rather than the instance: the canonical gate should run the formatter over tracked Go files and fail on a non-empty result, exactly as it already runs ShellCheck over tracked shell files.

## AI Execution State (P-A-U Loop)

This REQ was captured without the standard P-A-U block, and a worktree builder cannot write this file — so the orchestrator adds the block here and transcribes the builder's evidence from `do-work/runs/work-2026-08-18-211613/REQ-260-handback.md`. Without it `qualify.sh`'s box audit finds nothing to audit and passes vacuously; that silent half-disarming is REQ-264.

- [x] **[PLAN]:** Read `_dev/primes/prime-shell-commands.md` including § Closed Enumerations Go Stale, plus `general.md`, `coding-guardrails.md` and `communication-style.md` from the worktree. Studied the existing ShellCheck lane as the pattern to mirror, and probed gofmt's actual contract before writing anything — which is what surfaced D-01.
- [x] **[APPLY]:** One file touched, `_dev/tests/maintainer-verify.sh`. One committed increment, because the lane and its self-test coverage are not independently green: the lane alone breaks the stage-count assertions, and the shim alone breaks the same assertions from the other side. No adjacent cleanup, no version or changelog bump.
- [x] **[UNIFY]:** `git diff --stat HEAD~1` → one file, 75 insertions / 9 deletions. `shellcheck --severity=warning` on the changed file → exit 0. `git status --porcelain` empty — no stray fixtures, no scratch files, no leftover mis-formatting from the mutation probe. All experiment scratch was written to the session scratchpad, never into the worktree. Files checked: `_dev/tests/maintainer-verify.sh` (the only changed file) and `skills/do-work-board/tools/queue-kanban/durations.go` (mutated and restored, confirmed clean). **Class check:** the lane keys on the condition (a tracked Go file is not gofmt-clean), never on a path list — `git ls-files -- '*.go'` returns 56 files across two module directories and a future third module is covered with no edit here, unlike the `go vet`/`go test` lanes which name their two directories. Grepped for Go source under spellings a `*.go` pathspec would miss (`*.gotmpl`, `*.tmpl`, `*.go.txt`, `testdata`): none tracked.

## Requirements

- The gate runs the Go formatter over tracked Go files and fails when any is unformatted, selecting files the way the ShellCheck lane does (`git ls-files`, never a hand-maintained path list — Closed Enumerations Go Stale).
- The failure names the offending files.
- `bash _dev/tests/maintainer-verify.sh` exits 0 on the current tree — the package is formatted clean today, so a red result means the check found something real.

## Context

Discovered by REQ-251's builder ([low]); the one-character instance was fixed in passing by REQ-252. The user widened the scope at clarify: the instance is closed, the gap that let it through is not.

## Open Questions

- [x] I discovered this out-of-scope task while working on REQ-251: a gofmt formatting miss in `durations.go` from REQ-248. Should I process this as a new task? → Yes, and widened: add a formatter check to the canonical verify so this class cannot recur silently
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it — or fold it into REQ-252, which already owns the file.

**Answered [2026-08-18]:** User approved via `do-work clarify` and **widened the scope**: the original one-character fix is already satisfied by REQ-252's in-passing correction, so this REQ now covers making the gate run the formatter over tracked Go files. Title and requirements updated to match; `effort_estimate` raised from trivial to normal.

---

## Triage

**Route: A** - Simple

**Reasoning:** One named file, and the change mirrors a pattern that already exists inside it (the ShellCheck lane's `git ls-files` selection and fail-on-non-empty shape). Scope is obvious and bounded.

**Planning:** Not required

---

## Implementation Summary

**Files changed:**
- `_dev/tests/maintainer-verify.sh` (modified)

**What was done:** Added a Go-formatter lane to the canonical gate, immediately after the ShellCheck lane and mirroring its file selection — `git ls-files -z -- '*.go'`, never a hand-maintained path list — with a refusal to run if git reports zero tracked Go files, the count printed before checking, and a failure that names every offending path plus the `gofmt -w` remedy. Two properties of gofmt shaped it: it **exits zero even while listing unformatted files**, so the lane's verdict is the emptiness of its output rather than its exit status (a lane copied literally from the ShellCheck lane's exit-status shape would have looked correct and never fired); and it has no version flag, so it is resolved as `"$(go env GOROOT)/bin/gofmt"` to inherit the gate's existing exact `go1.26.1` pin instead of taking whatever `gofmt` comes first on PATH. Self-test coverage was extended in the same commit — a `gofmt` shim case through a fixture GOROOT, a `go env GOROOT` shim case, a dual-pathspec `git ls-files` shim case, the new stage in the success assertions, and **two** failure shapes rather than one: formatter-errors and formatter-succeeds-but-lists-a-file, the latter asserting the failure output names the offending file.

---

## Discovered Tasks

Transcribed by the orchestrator from the hand-back — a worktree builder cannot write this file, which is REQ-270.

- **[low] The self-test's cleanup trap emits `self_test_root: unbound variable` on every failure path.** `run_self_test` sets `trap 'rm -rf -- "$self_test_root"' EXIT` with a single-quoted body, so the variable expands at trap time — but it is a function `local`, already out of scope once a `fail_self_test` path returns. Every self-test failure therefore prints a spurious unbound-variable line after its real FAIL line, and the temp directory is not actually removed. **Pre-existing, not introduced here** — the builder verified this by mutating the ShellCheck lane in the *pre-change* script (`git show HEAD~1:_dev/tests/maintainer-verify.sh`), which produces the same trailing line. Exit codes are correct throughout. **Corrected by this REQ's review, which measured it:** the leak is not confined to failure paths. The self-test's own mutation sub-run is *designed* to fail, so **every `--self-test` run leaks one temp directory, green runs included** — and since `contract-regressions.sh` invokes `--self-test` and the gate invokes that suite, every full gate run leaks one. Measured at 37→38 directories in `/tmp` across a single green self-test, with 37 already accumulated, most predating this change. Scope the follow-up to that behaviour, not to the narrower one the hand-back described.

---

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh` against the merged tree (range `79d594c..307e146`), run un-piped with the exit code read directly
**Result:** ✓ `GATE_EXIT=0` — "Maintainer verification passed." The new lane is present and running in the integrated gate: `maintainer-verify: gofmt formatting check (56 tracked files)`. This run is both Step 6.5's testing and Step 8's post-merge verification.

**Red-green validation:** the acceptance criterion is that the lane *bites*, which cannot be shown by a green gate. Evidence comes from two independent mutation runs rather than from assertion:

- **Builder, on its branch:** reintroduced the exact slip from this REQ's story — a stray space in the day-domain truncation expression in `skills/do-work-board/tools/queue-kanban/durations.go` (`Truncate(24 * time.Hour)` → `Truncate(24  * time.Hour)`) — and observed `GATE_EXIT_MISFORMATTED=1` with the failure naming `skills/do-work-board/tools/queue-kanban/durations.go` and printing the `gofmt -w` remedy. Reverted, `git status --porcelain` empty, `GATE_EXIT_AFTER_REVERT=0`.
- **Independent reviewer, on the merged tree:** asked to reproduce that mutation itself rather than accept it, plus to attempt to *disarm* the lane and check whether `--self-test` catches the disarming. Findings recorded in `## Review`.

The orchestrator deliberately did not run a third mutation of the same file concurrently with the reviewer's — two agents mutating and reverting one path in one tree is how mutation evidence stops meaning anything.

**Prior-art check on the tool's contract (the reason the lane is not a copy of the ShellCheck lane):** `gofmt -l` on a deliberately mis-formatted file lists that file and **still exits 0** (observed, `exit-bad=0`). A lane trusting the exit status would have been decorative — green forever, catching nothing.

**New tests added:** self-test coverage in the same file — a `gofmt` shim case through a fixture GOROOT, a `go env GOROOT` shim case, a dual-pathspec `git ls-files` shim case, the new stage in the success-stage assertions (expected totals 8→9 and 9→10), and two failure shapes: `gofmt-lint` (formatter errors) and `gofmt-unformatted` (formatter exits 0 but lists a file), the second additionally asserting the failure output names the offending file. Self-test observed: `SELFTEST_EXIT=0`.

**Existing tests updated (cross-REQ impact):** none — the stage-count assertions changed as a direct consequence of adding the lane, not as a behavior change to any prior REQ's test.

*Verified by work action*

---

## Review

**Overall: 96%** | 2026-08-18T22:35:51Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded gate disposition — this is the durable audit record the gate mandates):**
- None.

**Minor findings:** 3 (report only, no follow-up REQ under Step 10's threshold)
- The `[ ! -x "$gofmt_command" ]` guard is the one element of the lane `--self-test` does not pin — deleting it leaves the self-test green. Fail-closed still holds without it (a missing formatter path exits 127 under `set -euo pipefail`, confirmed directly), so the guard is error-message quality rather than the safety mechanism. Recorded because it is exactly the shape a later simplify pass deletes. — gate: trivial
- The gofmt exit-status trap was not recorded in `_dev/primes/prime-shell-commands.md`, the REQ's declared prime and the stated home for this class. Routed through this action's Lessons-Capture Phase rather than a new REQ, and **written into the prime in this REQ's commit** — generalized past gofmt, since the class is any checking tool that reports findings on stdout while exiting zero. — gate: rule-change
- The transcribed Discovered Task understated the temp-directory leak; corrected in place above from the reviewer's measurement. — gate: trivial (pre-existing)

**Acceptance:** Pass — the lane bites (gate exit 1 on the reintroduced slip, naming the file and printing the remedy), reads output-emptiness rather than gofmt's zero exit status, and fails closed at exit 2 on a genuine formatter error. Six of seven attempted disarm mutations turn `--self-test` red, including deleting the lane, neutering its verdict, reverting to PATH resolution, and replacing the `git ls-files` pathspec with a hardcoded list. Restatement Sweep found nothing stale — no file outside the script enumerates the gate's lanes, which is a design property worth preserving.
**Suggested testing:** 3 items
**Follow-ups created:** None; **sweeps appended to:** None

*Reviewed by review-work action*

---

## Lessons Learned

**What worked:** Probing the tool's real contract before writing the lane. `gofmt -l` lists unformatted files and still exits 0, so the instruction "mirror the ShellCheck lane" — which was the right pointer — would have produced a decorative check if followed literally. The builder mirrored the file *selection* and the failure *shape* while deliberately differing on how the verdict is read, and said so. Second: writing two failure cases into the self-test rather than one. The generic error-injection case only proves the gate dies when the formatter *errors*; the real-world failure is a formatter that succeeds and prints a filename, and only the second case pins that.

**What didn't:** Nothing failed, but the review found the one unpinned element by attacking the lane rather than exercising it — seven disarm mutations, six caught. The `-x` formatter guard survives deletion with the self-test green. That asymmetry is the lesson: a lane can be well tested for *biting* and still be quietly removable, and only the disarm direction shows which parts are load-bearing.

**Worth knowing:** `local name="$(...)"` takes `local`'s own exit status, so a command substitution that crashes is masked from `set -e` — the shipped code splits the declaration from the assignment for exactly this reason, and reverting that "simplification" would let the gate pass green on a broken formatter run. Also: this gate's lanes are enumerated nowhere but the script itself. The justfile and `CLAUDE.md` both point at it and explicitly say the script owns the command inventory, which is why adding a lane needed no companion edits anywhere — worth preserving deliberately rather than by luck.

## Orientation

The canonical verify gate now fails on unformatted Go, so a formatting slip can no longer land and survive silently the way the day-domain truncation one did; lives in the maintainer verification suite, alongside the ShellCheck lane it sits beside and mirrors. No map change — one lane added to an existing script, no module, contract, or data flow altered. The declared prime `_dev/primes/prime-shell-commands.md` gained the exit-status-versus-stdout trap this REQ earned; its other referenced paths still exist.

