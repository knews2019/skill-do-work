# REQ-593 Remediation Hand-back — the two deferred false-verdict routes

**Branch:** `worktree-agent-REQ-593-heavy-tier-verdict`
**Head:** `45c860a9980922d7068b112bc3f1866bf2978cbf`
(fast-forwarded onto `claude/do-work-queue-drain-4ee2xl` at `358588a` first, then one commit on top)
**Worktree:** `/home/user/skill-do-work-worktrees/worktree-agent-REQ-593-heavy-tier-verdict`
Working tree is clean. Nothing outside the worktree was written except this file, which is neither
staged nor committed. Nothing was pushed.

## What changed

Three files, all inside the declared five-file scope. `+90 / -13`.

- `_dev/tests/update-script-behavior.sh` — all nine remaining `writer | grep -q` pipelines converted,
  plus a second self-check, `verify_no_check_reads_a_quiet_grep_pipeline`.
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_run.go` — `RunLanes` treats a zero
  or negative `LaneTimeout` as "use `defaultLaneTimeoutSeconds`", and passes the resolved value to
  `runOneLane`. The struct field now says so in its doc comment.
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_run_test.go` —
  `TestRunLanesWithoutLaneTimeoutUsesTheDefaultBound`.

`heavy_reuse_regression_test.go` was not edited. With the default in place it needs no change, and it
passes.

## R1 — the nine pipelines

All nine were fixed. None had to stay a pipeline.

**Seven `tar tzf <archive> | grep -q <marker>` sites (lines 607, 610, 615, 627, 630, 675, 678).** Each
archive's listing is now captured once into a variable, and every check on it reads that variable as a
herestring. A herestring is a temp-file redirect, so there is no writer process to kill and the status
is grep's verdict alone. Capturing once also replaced three `tar` invocations with one for the fallback
archive and two with one for each of the other two archives.

Discarding `tar`'s own exit status was considered against the prime's "Unchecked Exit Status Reads as
Content" rule. It is safe at all three archives:

- The fallback archive's readability is asserted directly above the block, by an explicit
  `tar tzf … >/dev/null` check that records its own failure.
- The requested-branch and default-branch archives are each read by a matched pair of checks — one
  that requires a marker and one that forbids the other branch's marker. An unreadable archive gives an
  empty listing, which fires the "omitted its marker" failure. The empty case cannot pass silently.

**Two `find … -print -quit | grep -q .` sites (lines 722, 725).** These do not lose the SIGPIPE race —
`find -print -quit` writes one short line and exits, so there is never a second write to fail — but they
were converted anyway, because the pipeline hid a different defect of the same family. Piped into
`grep -q`, a `find` that could not read the tree produced no output and was indistinguishable from a
tree with no leaked scratch, and under `pipefail` the failing status was absorbed by the `if`. Each scan
now captures `find`'s output and records a distinct failure if `find` itself failed, then tests the
captured value. That is one extra line per site and it removes the ambiguity rather than handling it.

**The lock-in is structural, and that is deliberate.** The existing
`verify_output_matchers_read_greps_verdict` proves the mechanism behaviorally on the two shipped
matchers with a 256 KiB output. Reproducing the same behavior per tar site would need a synthetic
archive per site and would still be proving the same mechanism nine more times. The new check instead
scans the script's own source for the shape:

```
verify_no_check_reads_a_quiet_grep_pipeline
```

It flags any non-comment line where a single `|` (not `||`) feeds a `grep` whose flags include `q`. A
`grep` reading a file or a herestring has no writer to kill and is not flagged. It uses `awk` rather
than a `grep | grep` pipeline so that a scan which could not read the file reports its own failure
instead of reading as a clean file.

### RED evidence for R1

Each of the nine sites was reverted **individually**, with the other eight fixes in place, and the
scanner was run against that file. Nine of nine caught, each naming its own line:

```
S1 scanner_exit=1 offending pipeline 637:   if tar tzf "$fallback_archive" 2>/dev/null | grep -q 'private-path'; then
S2 scanner_exit=1 ... 640:   if ! tar tzf "$fallback_archive" 2>/dev/null | grep -q 'VERSION'; then
S3 scanner_exit=1 ... 645:     || tar tzf "$fallback_archive" 2>/dev/null | grep -qv "^$fallback_top_level/"; then
S4 scanner_exit=1 ... 658:   if ! tar tzf "$requested_branch_archive" 2>/dev/null | grep -q 'requested-branch-marker\.txt'; then
S5 scanner_exit=1 ... 661:   if tar tzf "$requested_branch_archive" 2>/dev/null | grep -q 'default-branch-marker\.txt'; then
S6 scanner_exit=1 ... 707:   if ! tar tzf "$default_branch_archive" 2>/dev/null | grep -q 'default-branch-marker\.txt'; then
S7 scanner_exit=1 ... 710:   if tar tzf "$default_branch_archive" 2>/dev/null | grep -q 'requested-branch-marker\.txt'; then
S8 scanner_exit=1 ... 757:   if find "$fetch_root" -name 'preserved.tar.gz.download.*' -print -quit | grep -q .; then
S9 scanner_exit=1 ... 762:   if find "$fetch_root" -name 'preserved.tar.gz.fetching.*' -print -quit | grep -q .; then
```

It also catches a revert of the previously fixed matchers, so the two self-checks overlap rather than
leaving a gap:

```
offending pipeline 91:   if ! printf '%s' "$probe_output" | grep -Eq -- "$pattern_text"; then
```

Against the fixed file the scanner exits 0.

**Full-script ablation.** With all nine reverted in `_dev/tests/update-script-behavior.sh` itself, the
whole probe script (`DO_WORK_MAINTAINER_TIER=heavy`) exits **1** with exactly one FAIL line:

```
FAIL: this script feeds grep -q from a pipeline: under pipefail the writer's SIGPIPE death is reported as grep's verdict
```

Restored, the same script exits **0** with `update-script behavior probes passed.`

## R2 — `RunLanes` and the zero `LaneTimeout`

The fix is in `RunLanes`, not in the calling test. Nine lines including the reason:

```go
laneTimeout := request.LaneTimeout
if laneTimeout <= 0 {
    laneTimeout = defaultLaneTimeoutSeconds * time.Second
}
```

`runOneLane` now receives `laneTimeout` instead of `request.LaneTimeout`. Negative is folded into the
same branch, because a negative duration arms `time.NewTimer` exactly as badly as zero does.

### RED evidence for R2

With only that hunk reverted and the test kept:

```
--- FAIL: TestRunLanesWithoutLaneTimeoutUsesTheDefaultBound (0.09s)
    heavy_run_test.go:281: lane with no LaneTimeout = resultmodel.HeavyLaneExecution{LaneID:"green-lane",
      CommandArgv:[]string{"sh", "lanes/green.sh"}, ExitStatus:124, Skipped:false, WallSeconds:0, ...},
      want exit 0; exit 124 is the lane-timeout status, meaning the unset field armed an instant bound
```

`ExitStatus:124` is `HeavyLaneTimeoutStatus`: the lane was killed before it finished. The test reuses
the existing `green-lane` fixture, which sleeps a full second — far past an instant bound and far below
the 1800-second default — so this is deterministic rather than a replayed race.

Restored, the whole `heavyverification` package exits **0**.

## GREEN evidence

All commands from the worktree root, under the environment the request named.

1. **`bash _dev/tests/maintainer-verify.sh --heavy-lane updater`** → **exit 0**, 0 `^FAIL` lines, 69s.
   Run twice: once after the R1 fix, once again on the final committed bytes.
   Log: `.../scratchpad/lane-updater-final.log`
2. **`bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations`** → **exit 0**,
   0 `^--- FAIL` lines, 0 `^SKIP:` lines, **796 tests** (795 before, one added).
   Log: `.../scratchpad/lane-integrations-rem.log`
3. **The probe script direct** — `DO_WORK_MAINTAINER_TIER=heavy bash
   _dev/tests/update-script-behavior.sh` → **exit 0**, `update-script behavior probes passed.`
4. **`go test -count=1 ./internal/heavyverification/`** → **exit 0**, 13.4s.
5. **Lint** — `gofmt -l` on the package: empty. `go vet`: exit 0.
   `shellcheck --severity=warning _dev/tests/update-script-behavior.sh`: exit 0. `bash -n`: exit 0.

(`.../scratchpad` is
`/tmp/claude-0/-home-user-skill-do-work/213e30ac-5958-56c8-9fd2-faaaaf9c4ea6/scratchpad`.)

## The full canonical gate: green, and it carries R2's real-world evidence

Run once, as instructed:

```
env ... GIT_CONFIG_GLOBAL=<gate config> QUEUE_KANBAN_BROWSER=/opt/pw-browsers/chromium \
  bash _dev/tests/maintainer-verify.sh
```

**Exit 0.** `maintainer-verify: gate wall 85s`, `Maintainer verification passed.`, zero `^--- FAIL` and
zero `^FAIL:` lines. Log: `.../scratchpad/gate-rem.log`

**`TestLaneMutationCannotPublishOrReuseSuccess` is green in that gate run.** It is not guarded by
`testing.Short()` or by `DO_WORK_HEAVY_TESTS`, and the gate's fast stage runs `-short ./...` over the
whole `do-work-cli` module, so the test executes there. Confirmed directly:

```
=== RUN   TestLaneMutationCannotPublishOrReuseSuccess
=== RUN   TestLaneMutationCannotPublishOrReuseSuccess/commit=false
=== RUN   TestLaneMutationCannotPublishOrReuseSuccess/commit=true
--- PASS: TestLaneMutationCannotPublishOrReuseSuccess (0.21s)
```

That is the same test that failed in the previous session's gate run with `HEAVY-RUN-DIRTY-TREE` where
it wanted `HEAVY-RUN-REVISION-CHANGED`. It was not edited.

## Decisions

**D1 — the two `find` sites were converted even though they cannot lose the SIGPIPE race.** The request
said not to rewrite mechanically. These were rewritten for a different, real reason, stated in the code:
piped into `grep -q`, a failed scan and a clean tree produced the same verdict. Each now records a
distinct failure when `find` itself fails.

**D2 — the new lock-in is a source scan, not nine more behavioral replays.** The requirement was that
reverting any of these sites is caught rather than trusted. A scan for the shape does that for all
eleven sites in the file, including two that will never lose the race and any site added later. The
behavioral self-check that already exists proves the mechanism is real; repeating it per site would
prove the same thing nine more times and add a synthetic archive fixture for each.

**D3 — the scan uses `awk`, not `grep | grep -v`.** A `grep` pipeline that cannot read the file emits
nothing, and an empty result is this check's pass condition, so the scan would have had the exact defect
it exists to find. `awk` gives one process and one exit status, and the call site returns failure on it.

**D4 — zero-means-default folds in negative durations too.** `if laneTimeout <= 0` rather than `== 0`.
A negative `time.Duration` arms `time.NewTimer` the same way zero does, and no caller has a use for one.

## Discovered Tasks

- `_dev/tests/test-duration-log.sh` seeds its header with a `printf` to a candidate file plus
  `ln … || true`, then validates the header and hard-fails the whole lane if it is missing. Nothing
  makes that seeding atomic against a concurrent `record_test_duration` append, so a worktree can end
  up with a headerless duration log that fails every subsequent lane until a human deletes an ignored
  file, and the failure message does not say that deleting it is the remedy. Carried forward unchanged
  from the first hand-back; not hit in this session. → report only
- Changelog, mirror and `VERSION` are still owed. This commit changes a shipped file under `skills/`,
  which `_dev/primes/prime-releases.md` calls a release. All three are outside the declared five-file
  scope and belong to the work action's finalize transaction. Carried forward unchanged. → report only

## Acceptance criteria

- [x] All nine `writer | grep -q` pipelines in `_dev/tests/update-script-behavior.sh` converted; none
      had to stay a pipeline.
- [x] Reverting any one of the nine is caught rather than trusted — nine individual ablations, plus a
      full-script ablation with all nine reverted, exit 1 with a named FAIL line.
- [x] `RunLanes` applies `defaultLaneTimeoutSeconds` when `LaneTimeout` is zero or negative, so no
      caller has to remember.
- [x] `TestRunLanesWithoutLaneTimeoutUsesTheDefaultBound` fails with `ExitStatus:124` when that default
      is removed, and passes with it.
- [x] `heavy_reuse_regression_test.go` untouched and green.
- [x] Both named heavy lanes exit 0, read from `$?`; the full gate exits 0 and its
      `TestLaneMutationCannotPublishOrReuseSuccess` is green.
