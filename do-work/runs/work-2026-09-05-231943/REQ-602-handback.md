# REQ-602 Hand-back — repoint the lesson-satellite links whose archived targets moved, and check satellite links

Worktree: `/home/user/skill-do-work-worktrees/worktree-agent-REQ-602-satellite-links`
Branch: `worktree-agent-REQ-602-satellite-links` (base `claude/do-work-queue-drain-4ee2xl` at `9a5cc08`)
Head: `543564c6a2ddf0d2755558f42cea953c778f1ed2`
Working tree after the last commit: clean (`git status --short` prints nothing).

## Commits

Two commits, in this order:

1. `a70a04f` `[REQ-602] check that every relative link in a lesson satellite resolves` — only `_dev/tests/audit-lockins.sh` (+51 lines).
2. `543564c` `[REQ-602] repoint the fifteen satellite links that moved under UR directories` — the three satellites (15 lines changed, 15 insertions / 15 deletions) plus `do-work/lessons-index.md` (3 rows).

Why two: the request asks for the check to be proved red on the unfixed tree and green after. With the check in its own commit, `git checkout a70a04f && bash _dev/tests/audit-lockins.sh` reproduces the red run below on demand; one commit would have no revision where the check and the broken links coexist.

`git diff 9a5cc08 543564c --stat`:

```
 _dev/primes/lessons-action-files.md   | 24 ++++++++---------
 _dev/primes/lessons-kanban-board.md   |  2 +-
 _dev/primes/lessons-shell-commands.md |  4 +--
 _dev/tests/audit-lockins.sh           | 51 +++++++++++++++++++++++++++++++++++
 do-work/lessons-index.md              |  6 ++---
 5 files changed, 69 insertions(+), 18 deletions(-)
```

## The check

Location: `_dev/tests/audit-lockins.sh`, block headed `# Lesson-satellite links resolve (REQ-602)`, placed after the secret-pattern drift check and before the final exit block. It is run by the gate through `_dev/tests/contracts/probe-lanes.sh` (`register_probe audit_lockins_probe`), the same lane every other lock-in in that file uses.

Condition, not list: a satellite is any `_dev/primes/lessons-*.md` (a glob, with a FAIL when the glob matches nothing); a relative link is any `](target)` whose target is non-empty, does not start with `#`, and does not start with a URL scheme (`[A-Za-z][A-Za-z0-9+.-]*:`); the target with its `#fragment` stripped must exist (`-e`) relative to the satellite's own directory. awk's exit status is read separately from its output. No pipeline feeds a quiet reader; `_dev/tests/quiet-grep-pipeline-audit.sh` passes (95 tracked shell files). ShellCheck reports nothing on the new lines; its seven pre-existing style/info notes on that file (lines 28, 48, 118, 229, 230, 706, 709) are untouched and the gate's warning-level lint accepts them as before.

FAIL line shape: `FAIL: <satellite>:<line> links to <target as written>, which does not resolve from _dev/primes/.`

## Red proof (check committed, links not yet repointed — tree `a70a04f`)

`bash _dev/tests/audit-lockins.sh` printed exactly this, 15 FAIL lines, exit 1:

```
FAIL: _dev/primes/lessons-action-files.md:10 links to ../../do-work/archive/REQ-508-reduce-capture-templates-to-schema-backed-examples.md#lessons-learned, which does not resolve from _dev/primes/.
FAIL: _dev/primes/lessons-action-files.md:45 links to ../../do-work/archive/REQ-459-stage-command-owned-calibration-with-lifecycle-release.md#lessons-learned, which does not resolve from _dev/primes/.
FAIL: _dev/primes/lessons-action-files.md:46 links to ../../do-work/archive/REQ-417-implement-interview-memory-commands.md#lessons-learned, which does not resolve from _dev/primes/.
FAIL: _dev/primes/lessons-action-files.md:47 links to ../../do-work/archive/REQ-409-implement-safe-cleanup.md#lessons-learned, which does not resolve from _dev/primes/.
FAIL: _dev/primes/lessons-action-files.md:48 links to ../../do-work/archive/REQ-410-implement-doctor-forensics.md#lessons-learned, which does not resolve from _dev/primes/.
FAIL: _dev/primes/lessons-action-files.md:49 links to ../../do-work/archive/REQ-435-complete-doctor-forensics-delegation-contract.md#lessons-learned, which does not resolve from _dev/primes/.
FAIL: _dev/primes/lessons-action-files.md:50 links to ../../do-work/archive/REQ-411-implement-queue-selection.md#lessons-learned, which does not resolve from _dev/primes/.
FAIL: _dev/primes/lessons-action-files.md:51 links to ../../do-work/archive/REQ-412-implement-request-state-transactions.md#lessons-learned, which does not resolve from _dev/primes/.
FAIL: _dev/primes/lessons-action-files.md:56 links to ../../do-work/archive/REQ-498-make-orchestrator-finalization-resumable.md#lessons-learned, which does not resolve from _dev/primes/.
FAIL: _dev/primes/lessons-action-files.md:57 links to ../../do-work/archive/REQ-513-commit-the-claim-footprint-in-every-mode.md#lessons-learned, which does not resolve from _dev/primes/.
FAIL: _dev/primes/lessons-action-files.md:58 links to ../../do-work/archive/REQ-461-require-affirmative-project-owned-release-targets.md#lessons-learned, which does not resolve from _dev/primes/.
FAIL: _dev/primes/lessons-action-files.md:63 links to ../../do-work/archive/REQ-570-delete-the-pending-heavy-testing-status-held-requests-stay-claimed.md#lessons-learned, which does not resolve from _dev/primes/.
FAIL: _dev/primes/lessons-kanban-board.md:6 links to ../../do-work/archive/REQ-419-add-flat-just-recipes-action-delegation.md#lessons-learned, which does not resolve from _dev/primes/.
FAIL: _dev/primes/lessons-shell-commands.md:6 links to ../../do-work/archive/REQ-419-add-flat-just-recipes-action-delegation.md#lessons-learned, which does not resolve from _dev/primes/.
FAIL: _dev/primes/lessons-shell-commands.md:53 links to ../../do-work/archive/REQ-420-replace-shell-implementations-verify-parity.md#lessons-learned, which does not resolve from _dev/primes/.
exit=1
```

Every one of the fifteen targets the request lists appears above, once each (REQ-419 twice because two satellites cite it, which the request also lists twice). No other FAIL line: the seven other bare `do-work/archive/REQ-NNN-….md` links in the satellites (REQ-554, 555, 556, 593, 594, 595, 596) resolve because those files still live directly under `do-work/archive/`.

## Green proof (tree `543564c`)

`bash _dev/tests/audit-lockins.sh`:

```
Audit lock-in regressions passed.
exit=0
```

## Mutation evidence

On the repointed working tree (before the second commit, with a byte copy saved first), one repointed link in `_dev/primes/lessons-kanban-board.md:6` was changed back from `UR-081/REQ-419-…` to the old flat path, the check was run, the saved copy was restored, the check was run again, and the restored file was compared:

```
1d7eeef13955a59836c385ed3325610f4458f180cb046a1fe4e32b69e961749a  _dev/primes/lessons-kanban-board.md
== mutated: 
FAIL: _dev/primes/lessons-kanban-board.md:6 links to ../../do-work/archive/REQ-419-add-flat-just-recipes-action-delegation.md#lessons-learned, which does not resolve from _dev/primes/.
exit=1
== restored
Audit lock-in regressions passed.
exit=0
1d7eeef13955a59836c385ed3325610f4458f180cb046a1fe4e32b69e961749a  _dev/primes/lessons-kanban-board.md
cmp: byte-identical to the repointed copy
```

After committing, `git diff --exit-code -- _dev/primes/lessons-kanban-board.md` exits 0, so the restored file is byte-identical to the committed file as well (sha256 `1d7eeef1…749a` before and after).

## Shipped satellite URL report — `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`

Not edited (release-controlled). Every `https://github.com/knews2019/skill-do-work/blob/main/<repo path>` URL was extracted (`grep -o '…/blob/main/[^)]*'`), `<repo path>` taken with any `#fragment` stripped, and matched as an exact whole line against `git ls-files` (herestring, not a pipeline).

- URLs: **43** (42 distinct paths; `UR-081/REQ-460-…` appears twice).
- Resolve to a tracked file: **43**.
- Do not resolve: **0**. List: (empty).

So its targets did not move the way the maintainer satellites' did; every path already carries its `UR-0xx/` directory. Nothing to fix there.

## lessons-index.md rows

Computed by `/tmp/claude-0/-home-user-skill-do-work/213e30ac-5958-56c8-9fd2-faaaaf9c4ea6/scratchpad/req602-index-rows.sh` (bytes from `wc -c`; tokens `(bytes + 3) / 4` integer division; families = sorted unique `[family: slug]` markers; coverage `full` only when every `- ` bullet carries a marker). Run on the `9a5cc08` bytes it reproduces the three existing rows exactly, which is the check that the program implements the header's formula.

| Satellite | Bytes before → after | Tokens before → after | Families | Coverage |
| --- | --- | --- | --- | --- |
| `_dev/primes/lessons-action-files.md` | 19392 → 19476 | 4848 → 4869 | unchanged (8 slugs, 55 bullets, 14 marked) | `slugged: partial` (unchanged) |
| `_dev/primes/lessons-kanban-board.md` | 23641 → 23648 | 5911 → 5912 | unchanged (7 slugs, 65 bullets, 7 marked) | `slugged: partial` (unchanged) |
| `_dev/primes/lessons-shell-commands.md` | 17502 → 17516 | 4376 → 4379 | unchanged (9 slugs, 57 bullets, 10 marked) | `slugged: partial` (unchanged) |

Byte deltas are exactly the inserted directory names: action-files 12 links × 7 bytes (`UR-0xx/`) = 84; kanban-board 1 × 7 = 7; shell-commands 2 × 7 = 14. Only the `Tokens` cell of each row changed in `do-work/lessons-index.md`. The other four index rows were recomputed too and already match the index (516, 6711, 10543, 628).

## Guards

- `bash _dev/tests/audit-lockins.sh` on `543564c`: exit 0, `Audit lock-in regressions passed.`
- `bash _dev/tests/quiet-grep-pipeline-audit.sh`: exit 0.
- Full fast gate from the worktree: `DO_WORK_GATE_ROOT=/home/user/skill-do-work-worktrees/worktree-agent-REQ-602-satellite-links bash <scratchpad>/gate.sh` — exit 0, last line `Maintainer verification passed.` (gate wall 121s; go-test budget wall 32s, 797 tests). Log: `<scratchpad>/req602-gate.log`.

## Found and not fixed

- `_dev/primes/lessons-action-files.md:62` links REQ-509 by canonical GitHub URL rather than a relative path, unlike every other bullet in the maintainer satellites. The check skips absolute URLs by design, so that link is unverified locally; it does point at `do-work/archive/REQ-509-…`, which `git ls-files` tracks today, but it will go silently stale if REQ-509 moves under a UR directory. Out of scope (no wording changes; not one of the fifteen).
- The shipped satellite's canonical URLs have no local check at all (a `git ls-files` path match like the one above would cover them). The request only asked for a report; the shipped-package reference contract owns shipped Markdown and would be the home if wanted.
- Seven pre-existing ShellCheck style/info notes in `audit-lockins.sh` (listed above), all on lines this REQ did not touch.
