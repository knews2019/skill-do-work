# REQ-716 Hand-back (delete the five root installer copies in tools/)

Branch `worktree-agent-REQ-716-installer-scripts-tools-to-cmd`, tip `7610175f`, base `4845a9a7` (VERSION 0.315.1).
Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-716-installer-scripts-tools-to-cmd`

## Commits
- `66ca00c0` [REQ-716] RED: heavy-lane coverage rows name the staged installer scripts
- `7610175f` [REQ-716] delete the root tools/ installer copies; bootstrap runs the staged installer

## Files
- deleted (git rm): `tools/do-work-cli-preamble.sh`, `tools/fetch-upstream-archive.sh`, `tools/install-do-work-suite.sh`, `tools/replace-text-section.sh`, `tools/validate-suite-manifest.sh` (root `tools/` folder is gone)
- modified: `README.md` (bootstrap line runs `source/skills/do-work/tools/install-do-work-suite.sh`)
- modified: `skills/do-work/tools/install-do-work-suite.sh` (heredoc path; root-copy fallback and "probe both" comment removed)
- modified: `skills/do-work/tools/fetch-upstream-archive.sh`, `replace-text-section.sh`, `validate-suite-manifest.sh` (root-copy fallback and comment removed; messages and exit codes kept)
- modified: `skills/do-work/tools/do-work-cli-preamble.sh` (one candidate beside the preamble; the candidate loop became an `if`; comment no longer names a root copy)
- modified: `skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard.go` (shipped roots = declared modules + `suite`; loop became a plain check; comment updated)
- modified: `skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard_test.go` (case "the installer tools" deleted)
- modified: `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_maintainer_tree_test.go` (five coverage rows use staged paths and lanes)
- modified: `_dev/tests/heavy-lanes.json` (updater rows repointed; installer rows repointed, duplicate install row dropped)
- modified: `_dev/tests/staged-skills-contract.sh`, `install-suite-behavior.sh`, `update-script-behavior.sh`, `suite-manifest-contract.sh`, `contracts/replace-text-section.sh`, `audit-lockins.sh`
- modified: `_dev/primes/prime-releases.md` (shipped roots sentence)

## Scope changes
None outside the declared write set. Inside `_dev/tests/audit-lockins.sh` I also changed the preamble lock-in's comment ("The two exempt paths" -> "The exempt path", "a third copy" -> "a second copy") and its failure text ("outside the preamble pair" -> "outside the shared preamble"), because both described the deleted root copy.

## Red-green evidence
- After commit 1 (`66ca00c0`): `go test -run '^TestRepositoryManifestNamesEveryLaneScopedMaintainerEntryPoint$' ./internal/heavyverification/` FAIL, e.g. `skills/do-work/tools/do-work-cli-preamble.sh selected lanes = [staged-skills], want [staged-skills updater installer]` (all five rows failed).
- After commit 2 (`7610175f`): same test `--- PASS`.
- GREEN probe `REQ-716-probe.sh`: at base `4845a9a7` (detached scratch worktree, since removed) exit 1 `REQ-716 probe FAIL: a root tools/ path still exists`; at tip exit 0 `REQ-716 probe PASS` (17 s).

## Verification
1. `gofmt -l ./internal/finalization ./internal/heavyverification`: empty. `go vet`: clean. `go test -count=1 ./internal/finalization/ ./internal/heavyverification/`: ok (59 s, 32 s).
2. `suite-manifest-contract.sh` exit 0 (3 s); `contract-regressions.sh` exit 0 (28 s); `audit-lockins.sh` exit 0 (1 s); `contracts/replace-text-section.sh` exit 0 (8 s).
   With `DO_WORK_MAINTAINER_TIER=heavy`: `staged-skills-contract.sh` exit 0 (43 s); `install-suite-behavior.sh` exit 0 (32 s); `update-script-behavior.sh` exit 0 (71 s). "No such file" lines in all three logs: 0. No reruns were needed.
   `shellcheck -S warning` on all 11 changed `.sh` files: clean. Default severity prints only info/style notes on lines this REQ did not write (28 such notes at base, fewer now).
   `python3 -m json.tool _dev/tests/heavy-lanes.json`: ok.
3. `folder-structure-check`: 4 gaps (`VERSION`, `_dev/`, `skills/`, `suite/`), no `tools/` entry. `git ls-files tools`: empty. `ls tools`: No such file or directory.
4. GREEN probe: PASS (see above).
5. End to end (script `/private/tmp/claude-501/-Users-t2-Desktop-e1-experimental-repos-skill-do-work2/008bccaf-aa66-4644-8fb6-cac0974a4198/scratchpad/builder-REQ-716/e2e.sh`, run dir `/var/folders/2w/kw8sv6rd1z15yjykl787ryph0000gn/T/req716-e2e.iHXMgk`, left in place):
   - Fresh install via `<new>/skills/do-work/tools/install-do-work-suite.sh --archive new.tar.gz` (0.316.0): exit 0, installed VERSION 0.316.0.
   - Base install via `<old>/tools/install-do-work-suite.sh` (base is 0.315.1, not 0.315.0, because REQ-717 released 0.315.1): exit 0. Update through the installed `do-work-update.sh` with `DO_WORK_UPSTREAM_URL` on 127.0.0.1: exit 0, `Updated to v0.316.0 at ... using the four-module suite.`; the updated installer prints the new bootstrap path. Server stopped.
   - The server ran as `python3 -I -c` with `http.server` plus a 300 s self-exit timer instead of `python3 -I -m http.server`, because macOS has no `timeout(1)` and a backgrounded server must carry its own lifetime bound.
   - Old snippet path against the new archive: `bash <new>/tools/install-do-work-suite.sh` exit 127, "No such file or directory", before any write.

## P-A-U
- PLAN: RED test rows first. Then one atomic commit: git rm the five copies, repoint bootstrap/README, D-05/D-06/D-07, and repoint every maintainer suite to the staged scripts, including the two hostile plants.
- APPLY: Done as planned in two commits. The only addition is the audit-lockins comment/message wording (see Scope changes).
- UNIFY: `git diff --stat 4845a9a7 HEAD` = 22 files (+55/-323). Checked each one: the five deletions; the README line matches the heredoc (probe and install-suite byte-identity check); the four launchers keep their "missing beside" message and exit status; the preamble's `launcher_arguments` line is unchanged (the audit preamble pattern still matches it); there are no `tools/` rows left in heavy-lanes.json; there are no `$repo_root/tools` scan roots left in audit-lockins. REQ-717's lines are untouched: the `./queue-kanban` alternative, the board block and run-kanban checks. No debug artifacts.

## Decisions
- D-08 DECIDE & STATE: the preamble's one-candidate `for` loop became a plain `if [ -f ... ]`, matching D-05's "plain check if that reads better". The audit preamble lock-in still recognizes the preamble by its `launcher_arguments=(--format text)` line.
- D-09 DECIDE & STATE: `install-suite-behavior.sh:149-155` (retired-runtime list `tools/do-work-update.sh tools/queue-kanban`) is left alone. It is outside the brief, and it passes, since the paths are absent. `staged-skills-contract.sh`'s new `tools` entry already fails on any file under root `tools/`.

## Discovered Tasks
- `docs/decisions/records/adr-019-four-skill-suite-contract.md:11,42,67` names root `tools/validate-suite-manifest.sh`, and line 67 is a relative link to the deleted file. ADR-016 and ADR-018 already carry similar dead links to root `tools/queue-kanban/`. impact-minor → report only
- `_dev/tests/staged-skills-contract.sh:976-981` loops over `git ls-files tools/queue-kanban`, which has been empty since the board moved under `skills/do-work-board/`, so the loop checks nothing. It sits in REQ-717's region, so I left it. impact-minor → report only

## Integration seams
- `_dev/tests/staged-skills-contract.sh`: `legacy_runtime_paths` 83-94 (four `tools/...` entries → one `tools`; the `tools/queue-kanban/queue-kanban` exclusion kept at 105); retained bootstrap loop deleted (was 119-126); mirror loop and comment deleted (was 861-871). REQ-717's board block now sits about 25 lines earlier, otherwise unchanged.
- `_dev/tests/audit-lockins.sh`: scan root dropped at 48 and 64 (REQ-717's comment and `-e '\./queue-kanban '` pattern at 60-66 kept); preamble lock-in 143-155 (comment, scan root, single exemption, message).
- `_dev/tests/heavy-lanes.json`: updater 58-62 repointed; installer 78-81 (four staged rows, the old duplicate dropped).
- `README.md`: line 22 only.

## Proposed CHANGELOG entry
**Bootstrap Runs the Packaged Installer; Root tools/ Copies Removed**
The README install snippet now runs the installer from `skills/do-work/tools/` inside the downloaded archive. The five duplicate scripts in the repository's root `tools/` folder are gone.
What changes for an installed copy: `install-do-work-suite.sh --print-bootstrap-command` prints the `source/skills/do-work/tools/` path. The staged launchers lose a fallback lookup that never applied to an installed copy. Updates from older versions are unaffected, because installed updaters never run scripts from the archive. A saved copy of the old README snippet stops with "No such file or directory" before it writes anything; copy the current README block.

## Proposed lesson
- [family: hostile-fixture-writes-into-missing-dir] When a test fixture stops creating a folder, every hostile plant written into that folder fails with "No such file or directory". The suite stays green while the lock-in tests nothing. Grep the suite for writes under the removed folder, and repoint them to where the real files now live.

## Lessons read
`skills/do-work/tools/lessons-do-work-update.md` (full), `_dev/primes/lessons-releases.md` (full), `[family: declared-module-is-shipped]` bullet in `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`, `[family: assertion-passes-when-feature-is-dead]` bullets in `_dev/primes/lessons-shell-commands.md` (index lines; the trap is applied to the hostile plants and to the deleted mirror loop). Also read: `prime-shell-commands.md`, `prime-releases.md`, crew members `general.md`, `testing.md`, `shared-principles.md`, `communication-style.md`. I did not read `coding-guardrails.md` in full.
