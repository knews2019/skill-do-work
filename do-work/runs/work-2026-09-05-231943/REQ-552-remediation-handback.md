# REQ-552 remediation hand-back

REQ-552 replaces two places where the CLI started the external `find` and `cp` programs
with Go code that does the same job. The review scored it 68% and marked two of its five
acceptance criteria unmet. Both are now closed.

- Branch: `worktree-agent-REQ-552-replace-two-coreutils-exec-sites`
- Branch head: `4abb2aa1f5cee4cdeb4228d5a061f24073bb4006`
- Not pushed. Nothing outside the worktree was staged or committed.

## Files changed

- `_dev/tests/prescribed-shell-cases/architecture-report-preflight.sh`
- `_dev/tests/audit-lockins.sh`
- `do-work/working/REQ-552-replace-two-coreutils-exec-sites-with-the-pure-go-the-package-already-has.md`

The Go source was not touched. The five-file scope is unchanged; two of the five files
needed no further edit.

## F1 — the failed-copy case now runs the line the request changed

**Was:** the case pointed `TMPDIR` at a regular file. That breaks `os.CreateTemp` at
`architecture.go:126`, a guard this request did not touch, so the new `os.WriteFile` at
line 137 was never reached. Its only text assertion grepped for
`ARCHITECTURE-PREFLIGHT-FAILED`, which four sites in the same function emit.

**Now:** the publish call runs in a subshell under `ulimit -f 1` (a 512-byte cap on every
file the subshell creates) with a draft of about 2 KB, and `TMPDIR` points at a real
staging directory. `os.CreateTemp` succeeds, the staged write fails with EFBIG, and the
case asserts on the string `draft copy failed: write ` — the in-process write error, not
the shared finding code.

**Proof, three runs:**

| Tree | Result |
|---|---|
| Pre-change code at `7eadf50`, only this fixture file restored | RED — `architecture-report-preflight: 9 cases, 1 failure` |
| Post-change code with the new error handling ablated to `_ = os.WriteFile(...)` | RED — `9 cases, 4 failures` |
| As shipped | GREEN — `9 cases, 0 failures` |

The pre-change run is the one the review asked for. Against the `cp` version the evidence
string is the bare `draft copy failed: ` (GNU `cp` is killed by SIGXFSZ and writes nothing
to stderr), which is measured, not assumed — that is why the assertion names the write
error and why grepping for `draft copy failed` alone would not have been enough.

One trap is written into the fixture comment: the launcher's `go tool` build must already
be cached when the subshell runs, because a cold build cannot write under `ulimit -f 1`.
The earlier cases in the same file warm it. If it is ever not warm, the launcher exits 2
and the new text assertion fails loudly instead of passing on the wrong failure.

## F2 — the lock-in now sees any context variable

**Was:** `exec\.Command(Context)?\((ctx, )?"(find|cp|...)"` matched only a context argument
spelled exactly `ctx`.

**Now:** `exec\.Command(Context\([^,]+,|\()\s*"(find|cp|...)"`. This accepts any context
expression (`runContext`, `invocationContext`, `context.Background()`) while still
requiring the coreutils name to be the command argument itself — first for `exec.Command`,
straight after the context for `exec.CommandContext`. The command enumeration stays in one
place.

**Deviation from the remedy the review named, stated deliberately:** the review's
`exec\.Command(Context)?\([^)]*"(find|cp|...)"` also matches a legitimate
`exec.Command("git", "rm", "-r", path)`, because it does not require the coreutils word to
be the command. 85 of the 90 exec sites in this module run `git`, so that false positive is
reachable. The form committed here catches the same injection and does not.

**Proof, in the real worktree:**

- Injected `exec.CommandContext(runContext, "cp", a, b)` into shipped `architecture.go`:
  one `FAIL: coreutils spawned where the module already has pure Go: …` line, exit 1.
- Reverted (`git diff` on that file empty): `Audit lock-in regressions passed.`, exit 0.
- In a scratch copy, `exec.CommandContext(context.Background(), "cp", a, b)` also goes RED,
  and `exec.Command("git", "rm", "-r", a)` produces no coreutils FAIL line.
- The old pattern does not match the injected line; the new one does.

## F3 — a renamed module directory now fails loudly

The two hard-coded module directories move into one array checked with `[ -d ]` before the
scan. Previously `rg` exited 2 with its stderr discarded and the guard read the empty
output as a clean scan.

**Proof:** renaming `skills/do-work-board/tools/queue-kanban` in a scratch copy prints
`FAIL: coreutils lock-in cannot scan a missing module directory: skills/do-work-board/tools/queue-kanban`
and exits 1. Restoring the name silences it.

## F4 — the stale constraint lines

The request's `## Constraints` section said no test files would change while its own write
set listed two fixture files. Both false lines — "no test files beyond the lock-in" and
"Tests unchanged; the existing package tests are the safety net" — are struck through with
a pointer to D-01, which records the override and the measurement behind it.

**Deviation:** the brief said one line. I amended two, because both stated the same false
thing and leaving one knowingly wrong was worse than the extra strikethrough. Nothing else
in the request file changed.

## Verification

- Canonical gate `bash _dev/tests/maintainer-verify.sh`: exit 0, 83s wall, both fast stages
  `EXECUTING (no_prior_evidence)`, so the suite really ran against the changed files. Run
  once, as asked.
- Heavy lane `env DO_WORK_MAINTAINER_TIER=heavy bash _dev/tests/prescribed-shell-scripts-behavior.sh`:
  exit 0, `110 named script cases across 18 per-script files`. Per-file counts unchanged at
  9 (`architecture-report-preflight`) and 11 (`audit-archive-timestamps`).
- `shellcheck --severity=warning -x` on both shell files: exit 0.
- Every exit status read from `$?` directly, never through a pipe.

## Left for finalization

- **The release bump is still owed and is not done here**, as instructed. Two of the five
  touched files ship under `skills/`. The version number depends on other work in flight,
  and HEAD has moved past this merge with other unreleased changes under
  `skills/do-work-board`, so the bump size needs judging against the accumulated set.
- Three review findings were report-only and remain untouched: the `WalkDir` versus `find`
  ENAMETOOLONG divergence is not named in the accepted-difference list; no case covers a
  permission-denied archive (the suite runs as root, so `chmod 000` is a no-op); and the
  compatibility-shim block at `architecture.go:125-142` is now a pure round-trip whose
  simplification is not recorded as a discovered task.
