# Builder Brief — REQ-716 (delete the five root installer copies in tools/)

You are the builder for one work item of the do-work skill suite. You build in your own Git worktree and hand a branch back. The orchestrator merges, tests, reviews and releases.

- Operative name and branch: `worktree-agent-REQ-716-installer-scripts-tools-to-cmd`. The REQ title and slug still say "to cmd/", but the maintainer chose deletion: no `cmd/` folder is created.
- Worktree (the ONLY tree you write): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-716-installer-scripts-tools-to-cmd`. The coordinator creates it from main HEAD right before dispatch, with the branch checked out. Do not create it yourself.
- Main tree (READ-ONLY for you): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2`
- The request, with its decisions D-01 to D-07, Exploration and Scope: read `do-work/working/REQ-716-installer-scripts-tools-to-cmd.md` from the MAIN tree (your worktree may hold an older copy).
- Hand-back file (the one main-tree file you may write; never stage or commit it): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/2026-11-09-layout-queue/REQ-716-handback.md`
- REQ-717 (suite justfile recipes call board scripts in one line) is already merged into main. Your base includes it. Do not rename or move anything it added.

## Read first (in your worktree)
1. `skills/do-work/crew-members/general.md`, `coding-guardrails.md`, `shared-principles.md`, `communication-style.md`, `testing.md`.
2. `_dev/primes/prime-shell-commands.md` before you edit any `.sh` file, and `_dev/primes/prime-releases.md`.
3. `CLAUDE.md` (maintainer preferences: simple, delete before you add, focused tests that name the failure they pin).
4. Required lessons: `skills/do-work/tools/lessons-do-work-update.md` and `_dev/primes/lessons-releases.md`. Dropped only for budget, read the matching bullets: `[family: declared-module-is-shipped]` in `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` and `[family: assertion-passes-when-feature-is-dead]` in `_dev/primes/lessons-shell-commands.md` (find them with a fixed-string grep).

## Decided (do not reopen)
- **D-04 (maintainer, 2026-11-10):** delete the five root copies in `tools/` with git rm. No `cmd/` folder. The README bootstrap and the `--print-bootstrap-command` heredoc run `source/skills/do-work/tools/install-do-work-suite.sh`.
- **D-01 — no `tools/` forwarder.** No installed updater of any released version runs a script from the downloaded archive. The shell-era updaters and installers run their INSTALLED copies and pass the archive only as data; the Go updater works in-process. The REQ holds the commit-by-commit proof.
- **D-02 — no validator change.** The archive-layout validator requires only `VERSION`, `suite/modules.tsv` and the four `skills/` packages at the root.
- **D-03 — the staged copies in `skills/do-work/tools/` stay** and become the only copies.
- **D-05 — release guard:** remove `tools` from the shipped-root list in `finalization_release_guard.go:45`, leaving `suite` (simplify the one-item loop to a plain check if that reads better). Delete the test case "the installer tools" at `finalization_release_guard_test.go:53`; the case "a shipped action file" already covers a change under `skills/do-work/`.
- **D-06 — heavy lanes: repoint, do not just remove.** Lane rules are exact paths except the `staged-skills` lane's `skills` prefix. The root rows are what made a staged preamble, validator, replacer or fetcher change select the `updater` and `installer` lanes. Repoint each root row in `_dev/tests/heavy-lanes.json:58-62` and `:78-81` to `skills/do-work/tools/<name>.sh`. In the `installer` lane, delete the root install row instead, because `skills/do-work/tools/install-do-work-suite.sh` is already listed there. In `heavy_maintainer_tree_test.go:259-263`, the five rows become the staged paths with lanes `staged-skills, updater, installer` (and `staged-skills, updater` for the fetcher).
- **D-07 — delete the dead root-copy fallback** in the four staged launchers (`install-do-work-suite.sh:62-64`, `fetch-upstream-archive.sh:23-25`, `replace-text-section.sh:20-22`, `validate-suite-manifest.sh:16-18`: the "probe both" comment and the `../skills/do-work/tools/do-work-cli-preamble.sh` line) and in `do-work-cli-preamble.sh` (the second candidate `../skills/do-work/tools/do-work-cli.sh`, and the comment lines 5-9 that describe a root copy and a staged mirror). Keep each launcher's existing "missing beside this launcher" message and exit status.

## What to build
Two commits.

**Commit 1 (RED, test only).** In `heavy_maintainer_tree_test.go:259-263`, change the five coverage rows to the staged paths and lanes from D-06. Run `TestRepositoryManifestNamesEveryLaneScopedMaintainerEntryPoint` and record the failure (the manifest still names root paths).

**Commit 2 (the deletion, atomic).** Everything below lands in ONE commit, because the contract suites read the scripts by path:
1. `git rm` the five root files in `tools/`. Git removes the empty folder; confirm with `ls tools` failing.
2. Line 28 of the heredoc in `skills/do-work/tools/install-do-work-suite.sh` and `README.md:22`: `source/tools/install-do-work-suite.sh` becomes `source/skills/do-work/tools/install-do-work-suite.sh`. The README block must stay byte-identical to `--print-bootstrap-command` output.
3. D-05, D-06 (`heavy-lanes.json`) and D-07 as stated above. Update the release-guard comment at lines 24-26 so it no longer names root installer tools.
4. `_dev/tests/staged-skills-contract.sh`: delete the retained bootstrap loop (`:119-126`) and the byte-identity mirror loop with its comment (`:861-871`). In `legacy_runtime_paths` (`:83-96`) replace the four `tools/...` entries with one `tools` entry, so any file that reappears under root `tools/` fails; keep the `tools/queue-kanban/queue-kanban` find exclusion. Do not touch REQ-717's lines near 928-1000.
5. `_dev/tests/install-suite-behavior.sh`: `:10` installer path becomes `$repo_root/skills/do-work/tools/install-do-work-suite.sh`; fix the `:120` message; in `:166-174` delete the `mkdir -p "$archive_root/tools"`, the three copies into it and the chmod (keep a `mkdir -p "$archive_root"` before the first copy). `cp -R "$repo_root/skills"` already carries the staged scripts.
6. `_dev/tests/update-script-behavior.sh`: `:14-17` paths become `skills/do-work/tools/...`; in `build_suite_tree` (`:199-208`) drop root `tools` from the mkdir and delete the four copies and the chmod into it; `:557-562` fetcher path and message; `:790-796` delegation loop path. **Trap:** the hostile fixtures at `:410-415` (an `exit 91` installer planted at root `tools/install-do-work-suite.sh`) and `:450-454` (an always-pass `tools/validate-suite-manifest.sh`) then write into a missing folder; bash prints "No such file or directory", the suite still passes, and the lock-in tests nothing. Repoint both plants to `skills/do-work/tools/`, where the archive's scripts now live. After your change, the suite log must contain no "No such file" line.
7. `_dev/tests/suite-manifest-contract.sh:6,79` and `_dev/tests/contracts/replace-text-section.sh:9,11`: path and message use `skills/do-work/tools/`.
8. `_dev/tests/audit-lockins.sh`: drop `"$repo_root/tools"` from the scan roots at `:48`, `:64` and `:148`, and drop the root exemption line `:149`; keep the staged preamble exemption. Do not touch REQ-717's lines near 60-66 other than removing the root scan root.
9. `_dev/primes/prime-releases.md:5`: the shipped roots sentence no longer names `tools/`.

**Not yours:** `VERSION`, both `CHANGELOG.md` files and the version line in `skills/do-work/actions/version.md` (the finalizer writes the release); anything under `do-work/`; `.gitattributes`, `suite/modules.tsv`, the root `justfile`, `skills/do-work-board/` and `_dev/tests/contracts/queue-kanban.sh` (REQ-717); the Go updater, validator and fetcher packages; the history comments in `suite_manifest.go:2` and `managed_section.go:4`; the dated lesson at `_dev/primes/lessons-shell-commands.md:119`.

## Rules
- Write only inside the Scope's file list plus the deletion. If the request's own requirements need another file, proceed and report it under "Scope changes" with the reason; otherwise stop and report.
- Out-of-scope finds go to "Discovered Tasks" in the hand-back. Do not fix them.
- Commit by exact path (`git add <paths>`; the deletion with `git rm`). Never `git add -A` or `git commit -a`. Commit subjects start with `[REQ-716]`. End each commit message with the line `Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>`.
- Do not push. Do not merge. Do not start subagents. Do not edit the main tree except the hand-back file.
- Never run the repository gate `_dev/tests/maintainer-verify.sh`; the integrator runs it. A 30-second per-file time budget that fails once under load is rerun once, and both runs are recorded.

## Verify before you hand back (all from your worktree)
1. In `skills/do-work/tools/do-work-cli`: `gofmt -l ./internal/finalization ./internal/heavyverification` prints nothing; `go vet ./internal/finalization/ ./internal/heavyverification/`; `go test -count=1 ./internal/finalization/ ./internal/heavyverification/`.
2. `bash _dev/tests/suite-manifest-contract.sh`, `bash _dev/tests/contract-regressions.sh`, `bash _dev/tests/audit-lockins.sh`, then with `DO_WORK_MAINTAINER_TIER=heavy`: `_dev/tests/staged-skills-contract.sh` (about 35 s), `_dev/tests/install-suite-behavior.sh` (about 30 s), `_dev/tests/update-script-behavior.sh` (about 70 s; its log has no "No such file" line). `shellcheck` on every `.sh` you changed. `python3 -m json.tool _dev/tests/heavy-lanes.json > /dev/null`.
3. `bash skills/do-work/tools/do-work-cli.sh --repo-root "$PWD" --format text folder-structure-check` lists no `tools/` root entry. `git ls-files tools` prints nothing.
4. The GREEN probe: `bash /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/2026-11-09-layout-queue/REQ-716-probe.sh` run with your worktree as the current directory prints `REQ-716 probe PASS` (it fails at base by design).
5. End to end, in a fresh `mktemp -d` directory outside both trees (the pre-dispatch agent ran exactly this on a simulated deletion; repeat it on your real commit):
   1. Build the old archive with `git archive --prefix=skill-do-work-main/` from the base commit (`git merge-base HEAD main`), and the new archive from your branch HEAD. In the extracted new tree, set `VERSION`, `skills/do-work/VERSION` and the `**Current version**:` line of `skills/do-work/actions/version.md` to `0.316.0`, then re-tar it.
   2. Fresh install through the new bootstrap path: in a new Git repository with one empty commit, extract the new archive with `--strip-components=1` and run `printf 'y\n' | bash <extracted>/skills/do-work/tools/install-do-work-suite.sh --project-root <repo> --archive <new.tar.gz>`. Expect exit 0.
   3. Update from 0.315.0: in a second new Git repository, install the OLD archive through `<old-extracted>/tools/install-do-work-suite.sh` the same way and commit the result. Serve the new archive as `archive/refs/heads/main.tar.gz` with `python3 -I -m http.server <port> --bind 127.0.0.1` from a directory holding only it, then run `printf 'y\n' | DO_WORK_UPSTREAM_URL=http://127.0.0.1:<port>/archive/refs/heads/main.tar.gz bash <repo>/.claude/skills/do-work/tools/do-work-update.sh --project-root <repo>`. Expect exit 0 and `Updated to v0.316.0`. Stop the server.
   4. Record each exit status and the key output line.

## Hand-back file content
Use these headings:
- `## Commits` — hash and subject each.
- `## Files` — every file created, modified or deleted, with the verb; one line for the five-file deletion.
- `## Scope changes`
- `## Red-green evidence` — the heavy-lane coverage test failing after commit 1 and passing after commit 2; the GREEN probe failing at base and passing at your tip.
- `## Verification` — each command above with its result; for a rerun under load, both runs.
- `## P-A-U` — three short lines: PLAN, APPLY, UNIFY; the UNIFY line lists what you checked in `git diff --stat`.
- `## Decisions` — D-08 onward (D-01 to D-07 are taken), each marked DECIDE & STATE or ESCALATE, with Value and Risk for ESCALATE.
- `## Discovered Tasks`
- `## Integration seams` — what you changed in each shared file below, with line ranges after your edit.
- `## Proposed CHANGELOG entry` — a descriptive title (for example "Bootstrap Runs the Packaged Installer; Root tools/ Copies Removed") and the line "What changes for an installed copy": the installed `install-do-work-suite.sh --print-bootstrap-command` prints the `source/skills/do-work/tools/` path, and the staged launchers lose a fallback that never applied to an installed copy; updates from older versions are unaffected. Mention that a saved copy of the old README bootstrap stops with "No such file or directory" before writing anything.
- `## Proposed lesson` — one bullet, or "none".
- `## Lessons read`

## Integration seams
REQ-717 is merged into your base, so these are lines that already changed, not a future merge:
| Shared file | REQ-716 edits | REQ-717 already changed |
| --- | --- | --- |
| `_dev/tests/staged-skills-contract.sh` | lines 83-96, 119-126, 861-871 | lines 928-1000 (board template and run-kanban script checks) |
| `_dev/tests/audit-lockins.sh` | scan roots at 48, 64, 148 and exemption 149 | the delegating check comment and pattern at 60-66 (keep its `queue-kanban` pattern) |
| `_dev/tests/heavy-lanes.json` | `updater` 58-62, `installer` 78-81 | no change so far |
| `README.md` | line 22 only | none |

Your final message: under 2500 characters. State the branch tip hash, the verification results, and anything the orchestrator must decide.
