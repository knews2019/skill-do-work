# REQ-593 Exploration — the heavy tier's verdict

*Three independent diagnoses, run against the heavy-lane log at revision 6646ba51. Each
developed and proved its fix in a scratch copy without modifying the repository.*

## Lane: do-work-cli-integrations / internal/requeststate — TestRecoveryRefusesFalseLegacyCheckpointAbsence

**Environmental: True — caused by this run's changes: False**

### Root cause

The test commits into a fixture repository that has no git identity, and this container cannot auto-detect one.

Chain: `newStateRepository` (skills/do-work/tools/do-work-cli/internal/requeststate/state_apply_test.go:706) only runs `git init`. The package has a separate helper `configureStateGit` (same file:716-720) that sets user.name/user.email in the fixture. TestRecoveryRefusesFalseLegacyCheckpointAbsence starts at line 1047, calls newStateRepository at 1048, and commits at 1055 WITHOUT ever calling configureStateGit. The heavy lane argv forces GIT_CONFIG_NOSYSTEM=1 and GIT_CONFIG_GLOBAL=/dev/null, so no system or global identity is visible either. Git then falls back to auto-detection: the container hostname is `vm` with no domain, so git builds `root@vm.(none)`, judges it not a real address, and refuses with `fatal: unable to auto-detect email address`. On a normal workstation with an FQDN hostname (or a global user.email) the same test passes, which is why it was never seen before.

This test is the ONLY outlier in the package. Every other test that commits calls configureStateGit first: state_apply_test.go lines 134, 168, 197, 269, 542, 561, 582, 601, 621, 643, 673 and state_preimage_dirt_test.go:24 all do; line 1048 does not. So there is a real (small) gap in the repository's test fixtures on top of the environmental trigger.

### Verified fix

Environment-only fix, no repository change: add `EMAIL=heavy-lane@example.invalid` to the environment the heavy lanes run in (i.e. add it to the gate wrapper at /tmp/claude-0/-home-user-skill-do-work/213e30ac-5958-56c8-9fd2-faaaaf9c4ea6/scratchpad/gate.sh next to the existing `env -u NODE_OPTIONS ...` line). EMAIL supplies git's default author/committer email only when nothing else is configured, so fixtures that DO set their own identity keep it.

Equivalent verified alternative that does not depend on /etc/passwd gecos: GIT_AUTHOR_NAME="Heavy Lane" GIT_AUTHOR_EMAIL=heavy-lane@example.invalid GIT_COMMITTER_NAME="Heavy Lane" GIT_COMMITTER_EMAIL=heavy-lane@example.invalid. No test in do-work-cli asserts a commit author or committer, so the override is safe. A hostname change was not needed and was not used.

Neither variable is refused by _dev/tests/heavy-runtime-fingerprint.py: that script refuses only BASH_ENV, ENV, PYTHONPATH, PYTHONHOME, NODE_OPTIONS, LD_PRELOAD, LD_LIBRARY_PATH, DYLD_* and any GIT_CONFIG_* other than GIT_CONFIG_NOSYSTEM/GIT_CONFIG_GLOBAL. Neither changes the sealed fingerprint, which hashes tool bytes, go env and the module's git config, never environment values. Proven by the fingerprint-running heavyverification package passing in the green lane run below.

Should the repository be fixed instead? Yes, and it is the better long-term fix: insert `configureStateGit(t, root)` as a new line 1049, directly after `root := newStateRepository(t)` in TestRecoveryRefusesFalseLegacyCheckpointAbsence. That makes the test self-sufficient like its eleven siblings and removes the dependence on host hostname/global config. I proved this works in a scratch copy (see evidence, proof 3), but did not apply it to /home/user/skill-do-work.

Is that repository fix in scope for any request in flight? No. The seven requests in do-work/working are REQ-486 (collapsible UR progress summaries in the board UI), REQ-552 (replace two coreutils exec sites with pure Go), REQ-554 (move the shared commit/inspect body into the prescribed-shell guide), REQ-583 (pin evidence-gate remedy redirection), REQ-587 (Timeline view scroll surface), REQ-591 (reduce repeated setup in the fast gate) and REQ-592 (seal the do-work tree into both fast gate stages). None of them names the requeststate test fixtures, and all of them carry an explicit "no test files beyond the lock-in / scope strictly limited to planned files" clause. Folding the fixture fix into any of them would be scope creep. It needs its own small request.

### Evidence

Repository was not modified: `git status --porcelain` in /home/user/skill-do-work is empty after all work. All patched runs were in scratch copies under /tmp/.../scratchpad.

Proof 1 — package-level A/B, run in place (read-only), in skills/do-work/tools/do-work-cli:
  env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null go test ./internal/requeststate/ -count=1
    -> --- FAIL: TestRecoveryRefusesFalseLegacyCheckpointAbsence, "state_apply_test.go:1055: git [commit -qm legacy fixture]: exit status 128 ... fatal: unable to auto-detect email address (got 'root@vm.(none)')". Exactly one failing test in the package, byte-identical to heavy-run-1.log lines 3244-3260.
  same command + EMAIL=heavy-lane@example.invalid
    -> ok  github.com/knews2019/skill-do-work/do-work-cli/internal/requeststate  1.296s
  same command + GIT_AUTHOR_NAME/GIT_AUTHOR_EMAIL/GIT_COMMITTER_NAME/GIT_COMMITTER_EMAIL
    -> ok  github.com/knews2019/skill-do-work/do-work-cli/internal/requeststate  1.242s

Proof 2 — full lane A/B in a scratch clone (/tmp/.../scratchpad/repo-clone, byte copy of the repo), harness variables sanitized the same way gate.sh does (env -u NODE_OPTIONS -u GIT_CONFIG_COUNT -u GIT_CONFIG_KEY_0..2 -u GIT_CONFIG_VALUE_0..2), running the manifest command `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations`:
  without EMAIL -> exit 1, single failure: --- FAIL: TestRecoveryRefusesFalseLegacyCheckpointAbsence, FAIL .../internal/requeststate 6.211s   (log: /tmp/.../scratchpad/lane-noemail-clean.log)
  with EMAIL=heavy-lane@example.invalid -> exit 0, zero FAIL lines, 794 tests, wall 31s   (log: /tmp/.../scratchpad/lane-email-clean.log)
An earlier pair of runs that left NODE_OPTIONS set (/tmp/.../scratchpad/lane-baseline.log, lane-email.log) additionally failed internal/heavyverification with "default runtime must have a determinable fingerprint"; that is the known NODE_OPTIONS injection, confirmed by running the probe directly in the clone -> "fingerprint uncertain: opaque runtime extension: NODE_OPTIONS". It is unrelated to this lane failure and disappears once NODE_OPTIONS is unset.

Proof 3 — the repository-side fix works too, proved without touching the repo: copied skills/do-work/tools/do-work-cli to /tmp/.../scratchpad/cli-copy, inserted `configureStateGit(t, root)` after line 1048, then
  env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null go test ./internal/requeststate/ -count=1
    -> ok  github.com/knews2019/skill-do-work/do-work-cli/internal/requeststate  1.349s   (no EMAIL, no other env change)

Blame — not this session's work: both the failing test and the configureStateGit helper entered in commit 7bd3464 "[work run] work-2026-09-05-124800: REQ-585 qualified, tested, reviewed 94%, held for heavy lanes" (2026-09-05 16:24:19 +0300), which is the only commit that has ever touched state_apply_test.go. That commit predates every commit of REQ-592, REQ-486 increment 1, REQ-552 and REQ-554 (all of which are in 7bd3464..HEAD). None of those four requests touched internal/requeststate test files. The failure was latent from REQ-585 onward and only surfaced now because this heavy run is the first one executed on a container whose hostname has no domain part.

Environment facts confirmed directly: `hostname` -> vm; `getent passwd root` -> gecos field "root" (non-empty, which is why EMAIL alone suffices — git auto-detects the name but not the email); git version 2.43.0; go version go1.26.1 linux/amd64.

Adjacent observation, outside my lane and not investigated: in heavy-run-1.log the do-work-cli-integrations lane was recorded as "skipped in 29s ... SKIP: no browser is available" (lines 4172, 4187) even though its Go tests failed (line 3442). The SKIP text at log line 1911 belongs to the browser lane. A red lane being reported as skipped is worth a separate look by whoever owns the lane-runner result classification.

## Lane: do-work-cli-integrations

**Environmental: False — caused by this run's changes: False**

### Root cause

Two separate real defects, neither about a browser.

D1 - the SKIP label was false. The do-work-cli-integrations lane never asks whether a browser exists: run_heavy_lane() in _dev/tests/maintainer-verify.sh (line 857) runs `DO_WORK_HEAVY_TESTS=1 run_budgeted_go_tests .../do-work-cli ./...` and never calls browser_engine_available. QUEUE_KANBAN_BROWSER was therefore irrelevant to this lane; it is only read by the queue-kanban-browser lane. The string came from inside the lane's own test suite. runOneLane (heavy_run.go:264-311) tees the child's combined output through laneSkipWatcher, which marks the whole lane "did not run" as soon as any output line starts with "SKIP:" (laneSkipPrefix, heavy_run.go:27, matched at line 442). The lane runs the heavyverification package's own tests, and TestRunLanesMarksSkipFromExplicitSkipLine executes a fixture lane whose script literally prints "SKIP: no browser is available" (heavy_run_test.go:66). handleRunHeavyVerification hardcodes LaneOutputWriter: os.Stderr (heavy_commands.go:214), so that fixture line lands on the test process's stderr, which is the real lane's stderr, which the parent watcher reads. The parent cannot tell the fixture's echo from a real announcement. Worse, Skipped was computed without looking at the exit status, so the SKIP label suppressed the red finding for a lane that had actually failed.

D2 - the lane is genuinely red underneath. TestRecoveryRefusesFalseLegacyCheckpointAbsence (internal/requeststate/state_apply_test.go:1046-1055) is the only test in that file that commits in a fixture repository without first calling the configureStateGit helper that sets local user.name/user.email. Every lane argv in _dev/tests/heavy-lanes.json forces GIT_CONFIG_NOSYSTEM=1 and GIT_CONFIG_GLOBAL=/dev/null, so no identity is inherited and `git commit -qm "legacy fixture"` fails with "Author identity unknown", exit 128. This is not container-specific: the lane as the manifest declares it strips system and global config, so the test fails on any machine. D1 hid D2.

### Verified fix

Invocation that makes the lane RUN (no code change needed, run from /home/user/skill-do-work). Real exit status at HEAD c9799b5: 1 (red), one failing test:

env -u NODE_OPTIONS -u GIT_CONFIG_COUNT -u GIT_CONFIG_KEY_0 -u GIT_CONFIG_KEY_1 -u GIT_CONFIG_KEY_2 -u GIT_CONFIG_VALUE_0 -u GIT_CONFIG_VALUE_1 -u GIT_CONFIG_VALUE_2 GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations

Code fix, developed and run in a scratch clone at /tmp/claude-0/-home-user-skill-do-work/213e30ac-5958-56c8-9fd2-faaaaf9c4ea6/scratchpad/scratch-repo (nothing under /home/user/skill-do-work was modified). Patch file: /tmp/claude-0/-home-user-skill-do-work/213e30ac-5958-56c8-9fd2-faaaaf9c4ea6/scratchpad/0001-scratch-fix-lane-skip-misclassification-and-fixture-.patch (4 files, +54/-4).

1. internal/requeststate/state_apply_test.go - add configureStateGit(t, root) to TestRecoveryRefusesFalseLegacyCheckpointAbsence, so the fixture carries its own git identity. Fixes the red.
2. internal/heavyverification/heavy_commands.go - replace the hardcoded os.Stderr tee with a package variable laneOutputTee (io.Writer, default os.Stderr).
3. internal/heavyverification/heavy_run_test.go - runHeavyLanes points laneOutputTee at io.Discard for the test's duration, so fixture lane output never reaches the enclosing lane's stderr. Fixes the false SKIP.
4. internal/heavyverification/heavy_run.go - a non-zero exit status now clears the skip announcement, so a lane that ran and failed is reported red no matter what it printed. New lock-in test TestRunLanesReportsARedLaneThatAlsoPrintedASkipLine with a skip-then-fail-lane fixture (prints the skip line, exits 4) pins exactly the misreport seen in this run.

Verified in the scratch clone at commit 833939a: gofmt clean, go vet clean, `go test -count=1 ./internal/heavyverification/... ./internal/requeststate/...` with DO_WORK_HEAVY_TESTS=1 and the same stripped git config - both ok. Lane rerun exit 0, zero FAIL lines, zero lines starting with "SKIP:". End-to-end `do-work-cli run-heavy-verification --lane do-work-cli-integrations` in the scratch clone: "lane do-work-cli-integrations: exit 0 in 31s [executed: fingerprint_uncertain]", no findings.

### Evidence

Blame. Both defects arrive in commit 7bd3464 ("[work run] work-2026-09-05-124800: REQ-585 qualified, tested, reviewed 94%, held for heavy lanes", 2026-09-05), found with `git log -S "SKIP: no browser is available" -- .../heavyverification/` and `git log -S TestRecoveryRefusesFalseLegacyCheckpointAbsence -- .../requeststate/`. That is before REQ-592, REQ-486 increment 1, REQ-552 and REQ-554. The heavy-run log records the lane as "[executed: no_prior_evidence]", so 6646ba51 was the first time this lane ever executed after 7bd3464 landed; nothing this session touched heavy_run.go, heavy_commands.go or state_apply_test.go.

The failure was already in the original log, masked. /tmp/.../scratchpad/heavy-run-1.log line 3247: `state_apply_test.go:1055: git [commit -qm legacy fixture]: exit status 128: Author identity unknown`, line 3260 `--- FAIL: TestRecoveryRefusesFalseLegacyCheckpointAbsence`. Line 1911 carries the leaked bare `SKIP: no browser is available` from the fixture lane. Line 4187 then summarizes the lane as "skipped in 29s" and line 4172 emits HEAVY-RUN-LANE-SKIPPED instead of HEAVY-RUN-LANE-RED.

Unpatched direct lane run (log: /tmp/.../scratchpad/lane-direct.log): EXIT=1; `go-test budget: FAIL module=/home/user/skill-do-work/skills/do-work/tools/do-work-cli wall=26s exit=1`; exactly one `--- FAIL` (TestRecoveryRefusesFalseLegacyCheckpointAbsence) and exactly one line matching ^SKIP: at line 1913.

Patched lane run (log: /tmp/.../scratchpad/lane-fixed.log): EXIT=0; 0 lines matching ^--- FAIL; 0 lines matching ^SKIP:; 795 tests, wall 46s.

Patched end-to-end (logs: /tmp/.../scratchpad/heavy-run-fixed.log and .err): exit 0, "run-heavy-verification: success", "lane do-work-cli-integrations: exit 0 in 31s [executed: fingerprint_uncertain]", no findings section.

Uniqueness check for D2: an awk scan over state_apply_test.go for test functions that call runStateGit(..., "commit", ...) without configureStateGit returns exactly one name, TestRecoveryRefusesFalseLegacyCheckpointAbsence.

Note for whoever finalizes: the collision is broader than this one fixture. internal/corehelpers/checks.go lines 161, 169 and 245 print user-facing lines that begin with "SKIP:", and _dev/tests/update-script-behavior.sh line 36 and _dev/tests/contracts/probe-lanes.sh line 60 do too. Any lane whose output reaches stdout/stderr unindented with such a line can still be mislabeled. Change 4 (exit status outranks the skip line) limits the damage to green lanes; a stricter skip channel than "any line starting with SKIP:" is worth a follow-up REQ.

## Lane: updater

**Environmental: False — caused by this run's changes: False**

### Root cause

Real defect in the test harness: `_dev/tests/update-script-behavior.sh` lines 84-96. The file runs under `set -uo pipefail`, and both output matchers are written as a pipeline:

  assert_output_matches() { ... if ! printf '%s' "$probe_output" | grep -Eq -- "$pattern_text"; then record_failure ...
  assert_output_lacks()   { ... if   printf '%s' "$probe_output" | grep -Eq -- "$pattern_text"; then record_failure ...

`grep -q` exits at the FIRST match. The pattern `four-module suite` sits on line 2 of a ~36 KB probe output (the updater prints the whole reviewed install diff after it). Bash's builtin `printf` writes that 36 KB to the pipe in 4096-byte stdio chunks. When `grep -q` is scheduled early enough to find the match and exit while `printf` still has chunks left, `printf` dies on SIGPIPE (141). Under `pipefail` the pipeline's status becomes 141, `! pipeline` is true, and the probe is reported as "output did not match" even though the pattern was matched. That is exactly why the FAIL line quotes an output that plainly contains "archive layout: four-module suite." — the assertion and the output do not really disagree; the assertion is reading the writer's death, not grep's verdict.

It is a race, so it is load-dependent: the lane passes in isolation and fails on a loaded machine. `assert_output_lacks` has the same bug in the more dangerous direction — a pattern that SHOULD have been flagged can be silently swallowed when grep matches early and the writer dies.

The matchers came in with the repository import commit 7bd3464 and have never been edited. Nothing about the failure is specific to the updater's own code: `skills/do-work/tools/do-work-update.sh`, `tools/*`, the templates and `suite/modules.tsv` are byte-identical to the import.

### Verified fix

Feed grep a herestring instead of a pipeline, in `/home/user/skill-do-work/_dev/tests/update-script-behavior.sh` (lines 84-96). A herestring is a temp-file redirect, so there is no writer process and no pipefail interaction:

  # Both matchers feed grep a herestring rather than `printf ... | grep`. `grep -q` exits at
  # the first match, and under `set -o pipefail` the writer's SIGPIPE became the pipeline's
  # status, so a probe output long enough for the writer to still be writing was reported as
  # a miss whenever the reader won that race.
  assert_output_matches() {
    local pattern_text="$1" probe_name="$2"
    if ! grep -Eq -- "$pattern_text" <<<"$probe_output"; then
      record_failure "$probe_name — output did not match /$pattern_text/. Output: $probe_output"
    fi
  }

  assert_output_lacks() {
    local pattern_text="$1" probe_name="$2"
    if grep -Eq -- "$pattern_text" <<<"$probe_output"; then
      record_failure "$probe_name — output unexpectedly matched /$pattern_text/. Output: $probe_output"
    fi
  }

`<<<` appends a trailing newline that `printf '%s'` did not; grep is line-based, so no assertion changes meaning. `<<<` is already the house pattern elsewhere (`_dev/tests/session-start-hook-behavior.sh:140`, `_dev/tests/maintainer-verify.sh:660`).

Nothing in /home/user/skill-do-work was modified. The patched file lives at /tmp/claude-0/-home-user-skill-do-work/213e30ac-5958-56c8-9fd2-faaaaf9c4ea6/scratchpad/repo/_dev/tests/update-script-behavior.sh; `diff -u` against the real file shows only the block above.

Same construct exists in `_dev/tests/select-simple-reqs-behavior.sh` and `_dev/tests/p50-estimator-determinism.sh`, but those probe outputs are a few hundred bytes, far below the size where the race can open. Worth converting eventually, not required to make this lane green.

### Evidence

R1. Isolation runs are green at HEAD (c9799b5, one docs-only commit past 6646ba5), proving no static regression:
  - `DO_WORK_MAINTAINER_TIER=heavy bash _dev/tests/update-script-behavior.sh` -> exit 0, "update-script behavior probes passed."
  - `bash _dev/tests/maintainer-verify.sh --heavy-lane updater` -> exit 0.

R2. Replaying the exact bytes the failing run printed shows the pattern DOES match. I extracted lines 3629-4165 of the heavy log into /tmp/.../scratchpad/captured-output.txt (36,160 bytes) and ran the assertion body on it: `pipeline=0`.

R3. The race is real at that exact size. /tmp/.../scratchpad/race.sh runs the ORIGINAL matcher 500 times on those 36,160 bytes under `taskset -c 0` (single CPU, i.e. contention like a loaded heavy run):
  original  -> "false failures: 3 / 500", then "2 / 500" on a repeat
  herestring -> "false failures with herestring: 0 / 500"

R4. Deterministic reproduction of the identical FAIL line in the real script. In a scratch copy of the repo I enlarged the fixture payload so the reviewed diff exceeds the pipe buffer (writer guaranteed to still be writing when grep exits):
  FAIL: suite update: identifies layout — output did not match /four-module suite/. Output: Checking do-work updates…
  ...and the dumped output contains "four-module suite" 3 times. Exit 1, 2 failures. (/tmp/.../scratchpad/repro3.log)

R5. Same enlarged fixture with the herestring fix applied -> exit 0, "update-script behavior probes passed.", 0 FAIL lines. (/tmp/.../scratchpad/repro4.log)

R6. Fix-only scratch copy (diagnostic payload reverted), run through the real lane entry point: `bash .../repo/_dev/tests/maintainer-verify.sh --heavy-lane updater` -> lane exit=0, "update-script behavior probes passed.", 71s. (/tmp/.../scratchpad/repro5.log)

R7. Lint: `shellcheck --severity=warning` (the gate's severity, from _dev/tests/action-shell-blocks.sh:45) on the patched file -> status 0.

R8. Attribution. `git log -L84,96:_dev/tests/update-script-behavior.sh` returns exactly one commit, the import 7bd3464 — the matchers were never edited since. `git diff --name-only 7bd3464..HEAD` over the lane's declared coverage (do-work-update.sh, tools/, justfile.template, agent-instructions.template.md) shows only the test file changed, and the one commit that changed it, 4b2abdd (REQ-592, background workers bound their own lifetime), only added `time.AfterFunc(120*time.Second, ...)` to the fixture archive server — it does not touch output size or the matchers. REQ-486 increment 1, REQ-552 and REQ-554 touched none of the lane's covered files.

Mechanism confirmation (why the writer's status leaks): /tmp/.../scratchpad/pipe_size_test.py forces the pipe capacity to 4096 and feeds the same 36,161 bytes to `grep -Eq`: grep=0 while the writer exits non-zero with "printf: write error". At the default 65536 capacity the writer exits 0. Same shape as the scheduling race, just made deterministic.

