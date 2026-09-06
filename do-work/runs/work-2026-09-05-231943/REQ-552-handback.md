# Hand-back — REQ-552, replace two coreutils exec sites with pure Go

**Branch:** `worktree-agent-REQ-552-replace-two-coreutils-exec-sites`
**Commit:** `2a29fd3f5e6aa57ebf54f2ea7711b29384174c42` (base `d24c270`)
**Worktree:** `/home/user/skill-do-work-worktrees/worktree-agent-REQ-552-replace-two-coreutils-exec-sites`
**Files changed:** 5, exactly the request's `## Scope`. No sixth file.

```
_dev/tests/audit-lockins.sh                                      | 21 +++++
_dev/tests/prescribed-shell-cases/architecture-report-preflight.sh | 24 ++++---
_dev/tests/prescribed-shell-cases/audit-archive-timestamps.sh    | 21 +++--
skills/do-work/tools/do-work-cli/internal/corehelpers/commands.go | 16 +++--
skills/do-work/tools/do-work-cli/internal/toolboxcommands/architecture.go | 7 ++-
```

---

## RED (at base revision `d24c270`, before any edit)

The request's Reproduce command printed two lines and exited 0:

```
$ rg -n 'exec\.Command(Context)?\((ctx, )?"(find|cp|mkdir|grep|sed|ls|rm|mv|cat|touch|head|tail|wc)"' \
    skills/do-work/tools/do-work-cli skills/do-work-board/tools/queue-kanban --glob '!*_test.go'
skills/do-work/tools/do-work-cli/internal/toolboxcommands/architecture.go:133:		if output, copyError := exec.Command("cp", draftPath, stagedPath).CombinedOutput(); copyError != nil {
skills/do-work/tools/do-work-cli/internal/corehelpers/commands.go:724:	command := exec.Command("find", archiveRoot, "-name", "REQ-*.md", "-print0")
EXIT=0
```

The lock-in block, run standalone against the base revision, printed two FAIL lines:

```
FAIL: coreutils spawned where the module already has pure Go: skills/do-work/tools/do-work-cli/internal/corehelpers/commands.go:724:	command := exec.Command("find", archiveRoot, "-name", "REQ-*.md", "-print0")
FAIL: coreutils spawned where the module already has pure Go: skills/do-work/tools/do-work-cli/internal/toolboxcommands/architecture.go:133:		if output, copyError := exec.Command("cp", draftPath, stagedPath).CombinedOutput(); copyError != nil {
failure_count=2
```

Both heavy-tier fixture files were green at base, so the breakage below is caused by this change and not inherited:

```
audit-archive-timestamps: 11 cases, 0 failures.        EXIT=0
architecture-report-preflight: 9 cases, 0 failures.    EXIT=0
```

With the two Go edits applied and the fixtures still on their PATH shims, all five predicted assertions flipped, exactly as the exploration measured:

```
FAIL: audit-archive-timestamps failed-walk case exited zero after the walk failed
FAIL: audit-archive-timestamps failed-walk case reported clean for an archive it never scanned
audit-archive-timestamps: 11 cases, 2 failures.        EXIT=1
FAIL: architecture-report-preflight reported a failed copy as published
FAIL: architecture-report-preflight exposed partial HTML after a failed copy
FAIL: architecture-report-preflight selected a failed publication as a baseline
architecture-report-preflight: 9 cases, 3 failures.    EXIT=1
```

## GREEN (at commit `2a29fd3`)

```
$ rg -n 'exec\.Command(Context)?\((ctx, )?"(find|cp|mkdir|grep|sed|ls|rm|mv|cat|touch|head|tail|wc)"' \
    skills/do-work/tools/do-work-cli skills/do-work-board/tools/queue-kanban --glob '!*_test.go'
rg-exit=1        (no output)

$ bash _dev/tests/audit-lockins.sh
Audit lock-in regressions passed.
lockins-exit=0

$ go -C skills/do-work/tools/do-work-cli test -count=1 ./internal/corehelpers/ ./internal/toolboxcommands/
ok  	github.com/knews2019/skill-do-work/do-work-cli/internal/corehelpers	1.330s
ok  	github.com/knews2019/skill-do-work/do-work-cli/internal/toolboxcommands	1.130s
EXIT=0

$ env DO_WORK_MAINTAINER_TIER=heavy bash _dev/tests/prescribed-shell-scripts-behavior.sh
Prescribed shell script behavior probes passed (110 named script cases across 18 per-script files).
EXIT=0

$ bash _dev/tests/contract-regressions.sh
Contract regression checks passed.
EXIT=0
```

Also clean: `gofmt -l internal/corehelpers internal/toolboxcommands` (no output),
`go vet` on both packages, `shellcheck --format=gcc --shell=bash --severity=warning -x`
on all three shell files (exit 0; the three style notes at `audit-lockins.sh:28,48,118`
are pre-existing and below the gate's severity), and
`GOOS=windows GOARCH=amd64 go test -c` for both packages (exit 0 — dropping the two
subprocesses is a portability win the prime asked to check).

Every command ran inside the worktree (never a copy) and inside the sanitized env
wrapper from the brief.

---

## The two code edits

**`internal/corehelpers/commands.go`** — `archiveWalkFailure` now probes with
`filepath.WalkDir` and returns the first traversal error, formatted as
`"<path>: <error>"`, as its evidence. No filter on `REQ-*.md`: an unreadable subtree is
the failure whether or not it holds a matching file, which is what `find` reported too.
The `os.Stat` guard above is byte-identical. No import changes (`os/exec` still serves
the `curl` call at the bottom of the file).

**`internal/toolboxcommands/architecture.go`** — the compatibility-shim copy is
`os.WriteFile(stagedPath, data, 0o600)`, keeping the `draft copy failed: ` evidence
prefix. `os.CreateTemp` already made `stagedPath` at 0600 and `os.WriteFile` truncates an
existing file without re-applying the mode, which reproduces what `cp` did. No import
changes.

## The lock-in

One block appended after the Finding 2 block in `_dev/tests/audit-lockins.sh`, pasted
from the exploration's recipe. Both decided properties are intact: `--glob '!*_test.go'`
is kept (without it `suiteinstall/update_transaction_test.go:25`, a fixture that spawns
`cp -R`, would keep the count at 1 forever), and the command list is byte-identical to
the audit's own Reproduce pattern. The comment in the file states both reasons so the
next reader does not widen it. Registration at `contracts/probe-lanes.sh` untouched.

## The two fixture cases — rewritten, not deleted

Both drove the failure by putting a fake binary on `PATH`, which is inert the moment the
work happens in-process. Each now produces the same failure inside the command.

**`audit-archive-timestamps.sh`, the failed-walk case.** One archive branch is nested
past `PATH_MAX` (30 levels of a 200-character name), so the walk's own `open()` fails
with `ENAMETOOLONG`. A permission bit cannot be used: this suite runs as root in
containers, where `chmod 000` is a no-op, and the case would pass without a walk ever
failing. Verified in both directions before writing it in: with the deep branch the run
exits 1 with `audit-archive-timestamps: the archive walk failed — nothing was
inspected.`; with the branch removed it exits 0 with `do-work: archive audit clean
(1 file(s) scanned).`

**`architecture-report-preflight.sh`, the failed-copy case.** `TMPDIR` points at a
regular file, so the staged copy cannot be created and publication stops before anything
reaches the candidate directory. `GOTMPDIR` is set to a real directory in the same
invocation because the launcher resolves its binary with `go tool -n`, which also needs a
usable temporary directory — without `GOTMPDIR` the toolchain fails first, the run exits
2 with `do-work-cli: the Go toolchain could not build the command`, and the case passes
without ever reaching the publish path. That trap is measured and the reason is written
into the fixture beside the variable.

**Each case gained one assertion that names the failure it expects** (`walk failed` in
the first, the `ARCHITECTURE-PREFLIGHT-FAILED` finding code in the second). Without them
each case would pass on any unrelated nonzero exit its new mechanism happens to provoke,
which is the same way the PATH-shim versions went inert.

**Proof the rewrites are not inert.** With the two guards removed from a scratch build —
`archiveWalkFailure` discarding its evidence, and the `os.CreateTemp` failure swallowed —
the cases fail:

```
audit-archive-timestamps: 11 cases, 3 failures.        EXIT=1
architecture-report-preflight: 9 cases, 4 failures.    EXIT=1
```

The code was restored from a byte copy afterwards and the committed files match the
edits described above. Case counts are unchanged (11 and 9), so the aggregate stays at
110 named cases.

## Challenge recorded, work continued

The request's Constraints ("Tests unchanged; the existing package tests are the safety
net", "no test files beyond the lock-in") were not satisfiable — measured, five
`fail_case` lines across two files. The exploration recommended deleting both cases (its
O1). **I did not take that**, per the brief: deleting the archive case would leave the
archive-walk failure path with no coverage at all, and there is a mechanism that produces
each failure in-process. The write set grew by exactly the two fixture files the request's
`## Scope` already declares. Recorded here and in the commit message rather than raised as
a question, per CLAUDE.md § Communication.

## Discovered tasks

- **D1 — `commands.go:721-722` returns the evidence string `<nil>`.** When `do-work/archive`
  exists but is not a directory, `err` is nil, so `fmt.Sprint(err)` yields `"<nil>"`.
  Non-empty, so the gate still fires, but the operator is told nothing useful. Named
  out-of-scope by the brief; unchanged.
- **D2 — a lesson is owed to `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`.**
  The exploration (§8) says this work should produce one and that no prior lesson covers
  replacing a coreutils call with in-process Go. Three things this run measured belong in
  it: a fixture that controls behaviour with a fake binary on `PATH` dies silently when
  the code stops shelling out; `chmod 000` cannot produce an unreadable path for a suite
  that runs as root, so depth past `PATH_MAX` is the portable substitute; and breaking
  `TMPDIR` to fail an in-process temp file also breaks the `go tool` build inside
  `do-work-cli.sh` unless `GOTMPDIR` is set alongside it.
- **D3 — not a task, a routing note.** This commit changes shipped files under `skills/`,
  so per `_dev/primes/prime-releases.md` it is a release. The version bump, changelog
  entry and mirror are Step 9's (`work-reference.md` → Changelog Entry Procedure), owned by
  the finalizer, and are deliberately absent from this branch: they would have been a
  sixth and seventh file.

## Not blocked

Nothing is outstanding. Do not push; the worktree was left in place.
