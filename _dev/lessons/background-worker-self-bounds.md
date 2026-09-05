# Background-worker lifetime audit — 2026-09-05

The review of REQ-581 (making descendant-cleanup tests detect a real process-group leak) added six CPU loops for an optional loaded flake check. The incident report records a parent death before the trailing kill, about 4.5 hours of orphaned load, and roughly 24 CPU-hours consumed. The archived review records the extra `-count=5` experiment; its ordinary focused package run was about 2.6 seconds. The load generators caused the sustained CPU consumption. Durations from their lifetime are loaded measurements, not a normal queue-speed baseline.

The rule now lives in [Testing → Background Work and Synthetic Load](../../skills/do-work/crew-members/testing.md#background-work-and-synthetic-load). Parent cleanup is only an early stop; deliberately persistent background work owns its deadline.

## Sweep results

Searched all tracked-source areas, including every suite's actions, hooks, scripts and tools, root tools, and `_dev/tests`. Searched shell background operators, `$!`, `jobs -p`, kill/pkill and indefinite loops, including shell embedded in Go tests. Reviewed matching spawn sites with their payloads and later cleanup. Archived requests and run artifacts remain historical evidence.

| Site | Lifetime or reason it is safe within this scope |
|---|---|
| `_dev/tests/prescribed-shell-cases/generate-report-image.sh` and `generate-report-image-batch.sh` | Each long-lived fixture shell now sets its own 30-second deadline, including nested descendants and retained inactive examples. The bound exceeds the existing cleanup assertions. Existing wait-helper timers are finite sleep-then-signal operations; none were added. Ordinary success/failure backends write their fixture output and exit. |
| `skills/do-work/tools/do-work-cli/internal/ownedprocess/owned_process_group_unix_test.go` | Both respawning shells now have their own 30-second Bash deadline. The five-second failure assertion still detects broken cleanup before natural expiry. The separate builtin-read fixture receives EOF when its test-owned pipe writer closes, including on test-process death. |
| `_dev/tests/update-script-behavior.sh` archive server | The spawned server now starts its own 120-second exit timer. Parent kill/wait only stops it earlier. |
| `_dev/tests/prescribed-shell-cases/run-blocked-check.sh`; Go fixtures in `nextselection/blocked_probe_test.go`, `gittransaction/git_transaction_cancellation_test.go`, `heavyverification/heavy_run_test.go`, `heavyverification/heavy_reuse_regression_test.go`, `toolboxcommands/report_image_process_test.go` | Background payloads already use finite sleeps (10 or 30 seconds). Leaders wait only for those bounded children or exit sooner. |
| `suiteinstall/suite_commands_test.go` in the same CLI module | The detached reap-marker shell already checks its own 60-second deadline. |
| `_dev/tests/contracts/core-checks.sh` missing-argument probe | Current launcher replaces itself with the Go argument parser, which rejects the missing value before work or waits. The existing two-second regression timer also ends by itself. Documented at the site. |
| `_dev/tests/prescribed-shell-cases/capture-screenshot.sh` race probe | Each child performs one finite publication from a small local file and returns. No persistent workload awaits a kill. Documented at the site. |
| `_dev/tests/prescribed-shell-harness.sh` | Cleanup registry, not a workload. Its registered long-lived fixture payloads are bounded as above. |

## Scope exclusions

`session-start.sh` launches the repository's intentional gate scheduler, which polls HEAD and runs the default gate once per observed new revision. It is a service, not temporary work that its launching shell plans to kill. This explains another source of repeated test runs; the hook does not select `--heavy`. Changing that scheduler is outside this incident's requested fix.

`probe-batch.sh` launches suites to completion and waits for them; it does not schedule a later kill. The updater's self-signal stub is foreground code, not a deliberately backgrounded payload. The atomic-download stub's apparent infinite loop exits on its second iteration. No reaper, watchdog, or general orphan-process mechanism was added.

## Validation

Focused process-group and image-fixture tests exercise the existing cleanup assertions. A separate short replay uses the documented Bash/zsh deadline with sleeping work and a two-second bound, SIGKILLs its parent, and observes the child reach its own deadline. No CPU-load experiment is needed to prove that lifetime contract.

Verified: the focused `internal/ownedprocess` package passed in 0.939s; the single-image and batch-image fixture suites returned zero failures; the Bash and zsh parent-SIGKILL replays both reached the child deadline. The extracted archive-server fixture served the expected bytes and exited with status 124 at a shortened two-second deadline. Shell syntax, ShellCheck with sourced definitions, and Go formatting were checked on the changed files. The complete heavy suite and synthetic CPU load were not run.
