# Work Run — 2026-09-05 00:34:20Z

Resumed from `do-work/RESTART-PROMPT.md` (committed as `1147e1a2`) in a fresh Linux container.
Integration branch: `claude/do-work-ur-115-117-f9406c`. Fan-out: 4.

## Environment repairs made before any REQ work

The handoff was written on a macOS checkout. This container needed five repairs
before `bash _dev/tests/maintainer-verify.sh` could run at all. None of them
changes project source; all are host state.

1. **Go toolchain.** Host Go was 1.24.7; the launcher floor is 1.25.0 and the gate
   floor is `go1.26.1`. Pinned with `go env -w GOTOOLCHAIN=go1.26.1`, which the
   module's own `GOTOOLCHAIN=auto` then downloads.
2. **ShellCheck.** Absent. Installed 0.11.0, which is the gate's declared floor.
3. **`just`.** Absent, so `TestPublicationRecipePreservesHostileManifestArgvAcrossShellBoundary`
   failed with `/bin/sh: 1: just: not found`. Installed 1.42.4.
4. **`gc.auto` in the repository config.** Set by the session bootstrap. The shipped
   heavy-runtime fingerprint allows a closed set of scalar Git keys and refuses
   anything else as `opaque Git configuration`, so this one key made every lane
   fingerprint indeterminable. Removed with `git config --unset gc.auto`.
5. **Host environment variables the fingerprint treats as opaque.** `NODE_OPTIONS`
   (an explicit member of the probe's opaque-runtime-extension list) and the agent
   proxy's `GIT_CONFIG_COUNT`/`GIT_CONFIG_KEY_*`/`GIT_CONFIG_VALUE_*` ssh-to-https
   rewrite. Stripped per-invocation by a scratch wrapper rather than unset globally,
   so ordinary git operations keep the proxy rewrite.

Repairs 4 and 5 are what `TestShippedRuntimeEvidenceTracksEffectiveGoSettingsAndBinaryBytes`
and `TestShippedGitIsolationPreservesGenericLaneInheritance` were reporting. Both are
host state, not project state, so the judgment was to clear the obstacle rather than
touch the shipped contract.

The clone also arrived shallow (71 commits), which left `24ed2fdd` and `06367337` —
REQ-506's original base and implementation — unresolvable. `git fetch --unshallow`
restored full history (3193 commits).

## Baseline gate

`bash _dev/tests/maintainer-verify.sh` at `e4d78d81`, direct and unpiped, exited **0**
with zero failures. queue-kanban: 382 tests, 28s wall, slowest file `generate_test.go`
15.22s. do-work-cli: 728 tests, 42s wall, slowest file
`internal/publication/defer_gate_test.go` 16.96s. Every test file is under the 30s
budget. Log: `.git/work-run-20260905/REQ-577/gate-baseline.log`.

## Claims

- **REQ-577** — taken over with `recover --take-over REQ-577`; the command classified it
  `held for heavy lanes; claim preserved`, which is the intended outcome. Its six lanes
  drain at queue exhaustion.
- **REQ-574** — left untouched. It carries the other run's writer label and this run
  asserts no authority over it.
- Wave claimed at `--fan-out 4`: **REQ-506, REQ-510, REQ-544, REQ-562**.

## Dispatch

Three builders run concurrently, each in its own worktree under
`.git/work-run-20260905/` on a branch of the same name:

- `worktree-agent-REQ-506-run-evidence-gates-from-advance`
- `worktree-agent-REQ-544-anchor-every-lifecycle-gate-over-caller-authored-text`
- `worktree-agent-REQ-510-sweep-work-reference-sections-owned-by-cli-tests`

**REQ-562 is deliberately held from dispatch** even though it is claimed. It extends
REQ-448's timing prose in `work.md` and `work-reference.md`, and REQ-510 is at the same
moment deleting sections from `work-reference.md` to bring it under 700 lines. That is a
semantic collision, not a textual one, so a clean git merge would prove nothing. Its
builder goes out once REQ-510 is merged.

## Heavy-lane feasibility probe (run ahead of REQ-577's drain)

All six lanes were probed early rather than at queue exhaustion, because a red
lane there withdraws REQ-577's `commit:` and un-readies REQ-506, which is the
critical path for REQ-507.

| Lane | Result |
|---|---|
| queue-kanban-javascript | exit 0 |
| staged-skills | exit 0 |
| updater | exit 0 |
| installer | exit 0 |
| queue-kanban-browser | exit 1 — one test, under diagnosis |
| do-work-cli-integrations | exit 1 — cause found, see below |

Nothing was skipped. Chromium at `/opt/pw-browsers/chromium-1194/chrome-linux/chrome`
drives the browser lane fine via `QUEUE_KANBAN_BROWSER`, so `HEAVY-RUN-LANE-SKIPPED`
is not a risk here.

### F-01 — `do-work-cli-integrations` cannot pass on a host whose hostname git rejects

`TestRecoveryRefusesFalseLegacyCheckpointAbsence` builds a git fixture and commits
into it. The lane argv sets `GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null` on
purpose, so the fixture has no configured identity and git falls back to
auto-detecting one from the hostname. This container's hostname is `vm`, git derives
`root@vm.(none)`, rejects it, and the commit fails with "Author identity unknown".

This is not container-specific in the way it first looks. The fixture depends on the
host having a hostname git considers valid — true on a developer laptop, not true in
most containers or CI. Reproduced directly: the identical commit under the same two
variables fails bare and succeeds with `EMAIL` set.

Worked around for this run by exporting `EMAIL=do-work@localhost` in the run wrapper.
The durable fix is for the fixture to set its own identity in the repository it
creates, which is project source outside any currently claimed request. Recorded here
as discovered work.

### F-02 — the per-file 30s budget makes the gate unreliable under a hot builder fleet

REQ-580's post-merge gate exited 1 with **zero** failing test assertions. Two shell
fixtures blew the per-file budget instead: `session-start-hook-behavior.sh` at 106s
against a 30s limit, and `prescribed-shell-canonicalization.sh` at 67s. In the quiet
baseline run the same two files took 26s and 12s.

The cause is CPU contention: nine builders, a reviewer and the lane probes were all
running. These are shell fixtures, so `GOMAXPROCS` does not help them — they are
simply competing for four cores.

The consequence for this run is a real scheduling constraint, not a defect to fix:
**the canonical gate and a large builder fleet cannot run at the same time.** Since
integration is serial anyway and every merge invalidates the previous REQ's post-merge
verification, the answer is to stop adding builders, let the fleet drain, and run the
per-REQ gates in a quiet window. A gate result that has to be re-run is slower than
waiting for one that can be trusted.

REQ-574, held by the other run, is "bring do-work-cli test files under the 30s
budget" — the same budget, a different cause.

### F-03 — `queue-kanban-browser` cannot produce trustworthy evidence in this container

Diagnosed with 23 single-test runs across four configurations. Verdict: environmental,
and already known to this repository.

The only browser here is `/opt/pw-browsers/chromium-1194/chrome-linux/chrome`, which
reports **Chromium 141.0.7390.37**. `timeline_browser_probe_test.go:3243-3246` says in
its own comment that "Chrome 141.0.7390.37 failed this guard and is no longer a
compatibility target (REQ-375)", and `_dev/primes/prime-kanban-board.md` states the
strict browser lane targets current stable Chromium with Chrome 141 deprecated. So the
lane is being asked to run on a build the project has already ruled out.

The mechanism, measured through a scratch overlay rather than inferred: on the
capture-swallowed trial Chromium 141 dispatches no boundary event at all — no
`pointerleave`, no `pointerout` — because it does not update the hover chain while a
mouse button is held. The release `pointermove` is delivered and the pointer really is
outside the host, so the product behaviour is correct; the only failing assertion is
the test's own guard that its isolator was exercised. **Every product assertion passed
in all 23 runs.**

Pre-existing: neither `89ade961` nor `cd179d58` touches the board module, the test file
is byte-identical between them, and the base revision failed 3 of 4 runs in a throwaway
worktree. Headed mode under Xvfb does not fix it. No Chromium flag does either — it is
renderer hit-testing behaviour.

**Consequence for REQ-577.** Without `QUEUE_KANBAN_BROWSER` the lane prints
`SKIP: no browser is available` and exits 0, which is the designed behaviour for a host
with no supported browser. Under work Step 7.7 a skipped lane leaves the request
claimed and held with a `HEAVY-RUN-LANE-SKIPPED` finding for the next drain to retry.
That is the honest outcome and this run takes it. Setting the variable would instead
convert a documented deprecated-browser incompatibility into a red lane belonging to no
change in this run, and a red lane withdraws REQ-577's `commit:` and un-readies REQ-506.

The guard is also not deterministic even on the deprecated build — 2 passes in 23 runs —
so a green from this browser would not be worth trusting either.

**REQ-577 therefore cannot be finalized in this container.** It stays claimed and held.
This blocks nothing else: its landed `commit:` already makes it source-ready, which is
why the canonical selector offers REQ-506 as gate-deferred rather than dependency-unmet.
