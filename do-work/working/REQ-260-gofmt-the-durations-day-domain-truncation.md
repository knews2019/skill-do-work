---
id: REQ-260
title: Run the Go formatter as part of the canonical verify
status: claimed
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

- **[low] The self-test's cleanup trap emits `self_test_root: unbound variable` on every failure path.** `run_self_test` sets `trap 'rm -rf -- "$self_test_root"' EXIT` with a single-quoted body, so the variable expands at trap time — but it is a function `local`, already out of scope once a `fail_self_test` path returns. Every self-test failure therefore prints a spurious unbound-variable line after its real FAIL line, and the temp directory is not actually removed. **Pre-existing, not introduced here** — the builder verified this by mutating the ShellCheck lane in the *pre-change* script (`git show HEAD~1:_dev/tests/maintainer-verify.sh`), which produces the same trailing line. Exit codes are correct throughout, so this is cosmetic plus a temp-directory leak on failure paths only.

