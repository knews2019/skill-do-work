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
