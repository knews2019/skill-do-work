---
id: REQ-716
title: 'Delete the root installer copies in tools/; the bootstrap runs the packaged installer'
status: completed
created_at: 2026-11-09T21:47:36Z
user_request: UR-168
domain: general
prime_files: [_dev/primes/prime-shell-commands.md, _dev/primes/prime-releases.md]
tdd: false
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
required_lessons: [skills/do-work/tools/lessons-do-work-update.md, _dev/primes/lessons-releases.md]
estimate:
  p50_active_minutes: 40
  confidence: medium
  calculated_at: 2026-11-09T23:48:09Z
  basis:
    - Route B
    - 18-file write set
    - 3 subsystems involved
    - 7 acceptance criteria
write_set: [skills/do-work/tools/install-do-work-suite.sh, README.md, skills/do-work/tools/do-work-cli-preamble.sh, skills/do-work/tools/fetch-upstream-archive.sh, skills/do-work/tools/replace-text-section.sh, skills/do-work/tools/validate-suite-manifest.sh, skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard.go, skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard_test.go, _dev/tests/heavy-lanes.json, skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_maintainer_tree_test.go, _dev/tests/staged-skills-contract.sh, _dev/tests/install-suite-behavior.sh, _dev/tests/update-script-behavior.sh, _dev/tests/suite-manifest-contract.sh, _dev/tests/contracts/replace-text-section.sh, _dev/tests/audit-lockins.sh, _dev/primes/prime-releases.md]
claimed_at: 2026-11-09T23:39:42Z
route: B
dispatch_at: 2026-11-10T00:58:12Z
builder_handback_at: 2026-11-10T01:05:04Z
integration_at: 2026-11-10T01:07:16Z
review_at: 2026-11-10T01:15:15Z
kb_status: pending
commit: 68e1b00b940fbde618272eade66dfcaba26337ab
heavy_verified_at: 2026-11-10T01:20:06Z
heavy_verified_revision: 68e1b00b940fbde618272eade66dfcaba26337ab
completed_at: 2026-11-10T01:20:07Z
release_at: 2026-11-10T01:20:07Z
---
# Delete the root installer copies in tools/; the bootstrap runs the packaged installer
## What
Move the five root installer scripts in `tools/` (`install-do-work-suite.sh`, `fetch-upstream-archive.sh`, `do-work-cli-preamble.sh`, `replace-text-section.sh`, `validate-suite-manifest.sh`) to `cmd/` with `git mv`, without breaking a fresh install or an update from an older installed version.
## Detailed Requirements
- Map every caller before moving: the README bootstrap (`source/tools/install-do-work-suite.sh`), the Go update path, the shipped skills (about 34 files name these five scripts), and the contract tests.
- Rule out a break for already-installed consumers: if an older installed updater runs `tools/<script>` from inside the newly downloaded archive, keep a thin `tools/` forwarder, or prove no installed version does. Record which one in this REQ.
- Layout rule this applies: scripts live in `cmd/` (rule 5); the maintainer chose `cmd/` over `tools/` in the audit.
## Context
Source: do-work audit-current-folder-structure on this repository (2026-11-09); the maintainer settled every gap through the ask tool and chose to capture the prompts now.
Recorded exceptions from the same audit (do not move them): VERSION, skills/ and suite/ stay at the root because the release tool and the updater archive layout need them there; _dev/ stays as this repository's maintainer-only folder.
## Red-Green Proof
**RED prompt/case:** Run `folder-structure-check` on this repository.
**Why RED now:** It reports `FOLDER-STRUCTURE-ROOT-ENTRY tools/`.
**GREEN when:** `tools/` is no longer reported (or only a recorded forwarder remains), the README bootstrap installs the suite in a temp repository, an update from the previous release works, and the contract suites pass.
**Validation:** User confirmed
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-shell-commands.md` — 8480 tokens, `slugged: partial`; matched because the moved files are shipped shell and six maintainer test scripts change.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 24277 tokens, `slugged: partial`; matched because the release guard in do-work-cli internals changes its shipped roots (`declared-module-is-shipped` family).
## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** RED heavy-lane coverage rows first, then one atomic commit: git rm the five copies, repoint the bootstrap and README, D-05 to D-07, and every maintainer suite to the staged scripts, including the two hostile plants.
- [x] **[APPLY]:** Done as planned in two commits (66ca00c0 RED, 7610175f deletion); the audit-lockins preamble comment and failure text also stopped naming the deleted root copy.
- [x] **[UNIFY]:** git diff --stat 4845a9a7..7610175f: 22 files, +55/-323; each checked: README line equals the heredoc; launchers keep their messages and exit status; no tools/ rows left in heavy-lanes.json; no root tools scan left in audit-lockins; REQ-717's lines untouched; no debug output.
*Source: do-work capture-request: Move the five root installer scripts in tools/ (install-do-work-suite.sh, fetch-upstream-archive.sh, do-work-cli-preamble.sh, replace-text-section.sh, validate-suite-manifest.sh) to cmd/ (layout rule 5: scripts live in cmd/; the maintainer chose cmd over tools). Use git mv. First map every caller: the README bootstrap (source/tools/install-do-work-suite.sh), the Go update path, the shipped skills (80 files mention tools/; only the 34 that name these five scripts matter), and the contract tests. Rule out a break for already-installed consumers: if an older installed updater runs tools/<script> from inside the newly downloaded archive, keep a thin tools/ forwarder (or prove no installed version does) and say which in the REQ. Verify: folder-structure-check no longer reports tools/, the README bootstrap runs in a temp repo, an update from the previous release works, and the contract suites pass.*

---

## Triage

**Route: B** - Medium

**Reasoning:** The outcome is fixed (the five root installer copies in `tools/` go away), but the callers had to be discovered across the README, the Go release guard, the heavy-lane manifest and six maintainer test scripts, and the compatibility question needed evidence from older released code. No design choice is open once the callers are mapped.

**Planning:** Not required

The title and slug still read "to cmd/", but the maintainer chose deletion over a move (D-04): no `cmd/` folder is created.

## Open Questions

- [~] Does an already-installed older updater run `tools/<script>` from inside the newly downloaded archive, so that a thin `tools/` forwarder must stay? → **D-01**: Builder chose: no forwarder; the root `tools/` folder disappears completely. Reasoning: no installed updater of any released version runs a script from the downloaded archive. Every updater generation was read with `git show <commit>:<path>`: (1) the single-skill updaters (root `tools/do-work-update.sh`, 924e6689 to b583c783, 2026-07-29 to 2026-07-31) only extract the archive and read `$fresh_upstream/actions/version.md`; they execute nothing from it, and they already fail closed on every suite archive since the four-module cutover (a3c2612f), because that file is no longer at the archive root. (2) The bridge and shell-suite updaters (bef7334f to 0e8cf0d9, 2026-08-08 to 2026-08-17, and the shell launcher up to 916d2706) run `$script_dir/validate-suite-manifest.sh`, `$script_dir/install-do-work-suite.sh` and `$script_dir/fetch-upstream-archive.sh`, which are the INSTALLED copies beside the updater, and pass the archive only as data (`--root "$fresh_upstream"`, `--archive "$upstream_tarball"`). Every historical installer (8f22cbe4 to eb1a2c57) resolves its helpers the same way (`manifest_validator="$script_dir/validate-suite-manifest.sh"`, `section_replacer="$script_dir/replace-text-section.sh"`, `archive_fetcher="$script_dir/fetch-upstream-archive.sh"`) and reads only `VERSION`, `suite/modules.tsv` and `skills/...` from the archive. REQ-137's archived lesson records this as a rule: "the already-installed bridge must remain the authority for migration validation", after a first draft that trusted the archive's validator was rejected. (3) The Go updater (`update-suite`, from 73283d03, 2026-08-30, to HEAD) runs fetch, validate and install in-process (`skills/do-work/tools/do-work-cli/internal/suiteinstall/update_transaction.go:83-121`); `git log -G '"tools"|tools/(install-do-work-suite|validate-suite-manifest|replace-text-section|fetch-upstream-archive|do-work-cli-preamble)'` over `internal/suiteinstall`, `internal/suitemanifest`, `internal/archivefetch` and `internal/managedsection` finds only a test temp-dir name and two history comments. The one archive-relative caller of root `tools/` is the fresh-install bootstrap (`README.md:22` and the heredoc in `install-do-work-suite.sh --print-bootstrap-command`), which is fetched fresh from the README and changes in the same commit. Value: the root layout gap closes completely, with no forwarder to maintain or to report as a new root entry. Risk: a person who saved an old copy of the README bootstrap snippet (in a CI script or notes) gets `No such file or directory` from `bash "$bootstrap_tmp/source/tools/install-do-work-suite.sh"`; the snippet runs under `set -e`, so it stops before any project write, and the fix is to copy the current README block. Reversible: a later patch can re-add a forwarder file.
- [~] What does the updater's archive-layout validation require at the archive root? → **D-02**: Builder chose: no validator change. Reasoning: `suitemanifest.ValidateSuite` (`skills/do-work/tools/do-work-cli/internal/suitemanifest/suite_manifest.go:52-200`) requires only `VERSION`, `suite/modules.tsv`, the four `skills/<module>/SKILL.md` packages, `skills/do-work/VERSION` and `skills/do-work/actions/version.md`; it never names `tools/` or `cmd/`. Every historical shell validator (f8aecd89 to 6b74ec69) has zero `tools` mentions. So neither old nor new validators care where the root scripts live. Value: no validator edit, no new required path. Risk: none found; the builder's update-from-previous-release check (Scope acceptance 4) confirms it end to end.
- [~] Does layout rule 5 ("scripts live in `cmd/`") apply to a skill package's own `tools/` folder, so the staged copies in `skills/do-work/tools/` must move too? → **D-03**: Builder chose: no; the staged copies stay in `skills/do-work/tools/` and become the only copies. Reasoning: `folder-structure-check` reports root entries only (its root allow-list is `skills/do-work/tools/do-work-cli/internal/toolboxcommands/folder_structure.go:28`); at HEAD it reports `FOLDER-STRUCTURE-ROOT-ENTRY tools/` and nothing under `skills/`. Installed consumers run the staged copies at `.claude/skills/do-work/tools/`, and every installed launcher path, `version.md` and the justfile recipes name them, so moving them would break installed copies. Value: the change stays root-only. Risk: none for consumers.
- [x] Keep five byte-identical root copies (moved to `cmd/`), or delete them and point the bootstrap at the staged copy? (D-04) → Maintainer chose (2026-11-10, through the ask tool): delete the five root copies in `tools/` with git rm; no `cmd/` folder; the README bootstrap and the `--print-bootstrap-command` heredoc run `source/skills/do-work/tools/install-do-work-suite.sh`. The pre-dispatch agent had recorded this as the smaller alternative (D-04) after finding that each root file is a byte-identical copy of the shipped `skills/do-work/tools/<name>.sh` and that the bootstrap is their only archive-relative caller. D-01 still holds: no updater runs a root script from the archive.
- [~] What happens to the release guard's shipped-root list, which names `tools`? → **D-05**: Builder chose: remove `tools`, leaving `suite` as the only extra root (replace, never add; no `cmd`). Reasoning: after the deletion no tracked path is under root `tools/`, so the entry can never match; the installer scripts now live in the declared module `skills/do-work`, which the guard already counts. The test case "the installer tools" (`finalization_release_guard_test.go:53`) is deleted, not repointed: the case "a shipped action file" already proves a change under `skills/do-work/` is shipped. Value: the guard names only roots that exist. Risk: none; a path under a deleted folder cannot appear in a diff.
- [~] The heavy lanes cover the root copies by exact path. Remove those rows, or repoint them? → **D-06**: Builder chose: repoint each root row to its staged path, and drop the one row that would duplicate an existing staged row. Reasoning: every lane rule is an exact path except the `staged-skills` lane's `skills` prefix. The root rows are what made a change to the staged preamble, validator, replacer or fetcher select the `updater` and `installer` lanes, because the mirror check forced every staged edit to touch the root copy too. Removing the rows would silently drop those lanes for exactly the scripts the bootstrap and updater tests run. The `installer` lane already names `skills/do-work/tools/install-do-work-suite.sh`, so its root install row is deleted instead of repointed. Net rows go down by one. This differs from the coordinator's "remove them"; the reason is the lost lane selection. Value: lane selection for these scripts stays as it was. Risk: none; the coverage test pins the new rows.
- [~] The four staged launchers and the preamble still probe `../skills/do-work/tools/` for the preamble and `do-work-cli.sh`, a fallback that existed only for the root copies. Keep or delete? → **D-07**: Builder chose: delete the fallback lines and correct the comments that describe a "root copy" and a "staged mirror". Reasoning: after the deletion that path can only resolve to `skills/do-work/skills/do-work/tools/`, which never exists, so it is dead code and its comments become false in shipped files. A scratch run of the staged-skills, install-suite and update-script suites with the fallback removed passed. Value: shipped launchers describe what they do. Risk: low; a launcher copied away from its preamble already fails with its existing "missing beside this launcher" message.
- D-08 (builder, DECIDE & STATE): the preamble's one-candidate `for` loop became a plain `if [ -f ... ]`. The audit preamble lock-in still recognizes the preamble by its `launcher_arguments=(--format text)` line.
- D-09 (builder, DECIDE & STATE): the retired-runtime list in `_dev/tests/install-suite-behavior.sh` (root `tools/do-work-update.sh` and `tools/queue-kanban`) stays as it is. The paths are absent and the test passes; the staged-skills contract's new root `tools` entry already fails on any file under a root `tools/` folder.
<!-- D-XX counter: last used D-09. Next decision: D-10. -->

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

Required lessons read: `skills/do-work/tools/lessons-do-work-update.md` (REQ-137: the installed validator is the authority for every candidate archive; this is the rule that makes D-01 safe) and `_dev/primes/lessons-releases.md` (REQ-715 `canonical-link-outlives-its-target`: sweep both directions when a path goes away). Also read for this REQ: the `declared-module-is-shipped` bullet of `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` and the `assertion-passes-when-feature-is-dead` bullets of `_dev/primes/lessons-shell-commands.md`.

**What the five root files are.** Each root `tools/<name>.sh` is a byte-identical copy of the shipped `skills/do-work/tools/<name>.sh` (checked with `cmp` at HEAD 792f981c). The root copies exist only so the fresh-install bootstrap can run the installer from an extracted archive. The staged installer works from the extracted archive too: it sources the preamble beside it, which finds `do-work-cli.sh` beside it. So the bootstrap can run `source/skills/do-work/tools/install-do-work-suite.sh` and the root copies can go.

**Callers of the ROOT copies (must change):**
- `README.md:22` — the bootstrap block runs `source/tools/install-do-work-suite.sh`. `_dev/tests/install-suite-behavior.sh:203-215` requires this block to be byte-identical to `--print-bootstrap-command` output.
- `skills/do-work/tools/install-do-work-suite.sh:28` — the heredoc printed by `--print-bootstrap-command`.
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard.go:45` — shipped roots `suite` and `tools`; comment at lines 24-26 says "the installer tools". Pinned by `finalization_release_guard_test.go:53` (D-05).
- `_dev/tests/heavy-lanes.json:58-62` (lane `updater`) and `:78-81` (lane `installer`) — exact paths `tools/<name>.sh`. Pinned by `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_maintainer_tree_test.go:259-263` inside `TestRepositoryManifestNamesEveryLaneScopedMaintainerEntryPoint` (D-06; the file is export-ignored, maintainer-only).
- `_dev/tests/staged-skills-contract.sh` — `:83-96` legacy root runtime list (`tools/checks`, `tools/do-work-update.sh`, `tools/queue-kanban`, `tools/prime-do-work-update.md`), `:119-126` retained bootstrap tools loop (requires four root files), `:861-871` byte-identity mirror loop over `"$repo_root"/tools/*.sh`. Both loops go away with the copies. Trap: a mirror loop left in place globs nothing after the deletion and passes while checking nothing.
- `_dev/tests/install-suite-behavior.sh:10` (installer path), `:120` (message), `:166-174` (synthetic archive: `mkdir -p "$archive_root/tools"` and three copies into it; `cp -R "$repo_root/skills"` already carries the staged scripts).
- `_dev/tests/update-script-behavior.sh:14-17` (four script paths), `:199-208` (`build_suite_tree` creates root `tools/` and copies four scripts into it), `:410-415` and `:450-454` (hostile fixtures plant an `exit 91` installer at root `tools/install-do-work-suite.sh` and an always-pass `tools/validate-suite-manifest.sh`, to prove the installed engine never runs the archive's scripts; this is an existing lock-in for D-01), `:557-562` (fetcher path and message), `:790-796` (delegation loop names `$repo_root/tools/install-do-work-suite.sh`). Trap: once `build_suite_tree` stops creating root `tools/`, the two hostile plants write into a missing folder; bash prints "No such file or directory", the suite still passes, and the lock-in tests nothing. Repoint both plants to `skills/do-work/tools/`, where the archive's scripts now live and where the new bootstrap runs them.
- `_dev/tests/suite-manifest-contract.sh:6`, `:79`; `_dev/tests/contracts/replace-text-section.sh:9`, `:11` — script path and message.
- `_dev/tests/audit-lockins.sh:48` (rg scan roots), `:64` (find roots), `:144-150` (hand-rolled preamble scan; exempts `$repo_root/tools/do-work-cli-preamble.sh` and the staged preamble by full path). Drop the root `tools` scan root and the root exemption; the staged exemption stays.
- `_dev/primes/prime-releases.md:5` — names `tools/` as a shipped root twice.
- Root-copy fallback (D-07): `skills/do-work/tools/install-do-work-suite.sh:62-64`, `fetch-upstream-archive.sh:23-25`, `replace-text-section.sh:20-22`, `validate-suite-manifest.sh:16-18` (comment "probe both" plus the `../skills/do-work/tools/do-work-cli-preamble.sh` line) and `do-work-cli-preamble.sh:5-16` (comment naming the root copy, second candidate `../skills/do-work/tools/do-work-cli.sh`).

**Not callers of the root copies (no change):**
- Shipped files that name a script basename mean the package's own copy: `skills/do-work/docs/prescribed-shell-primitives.md:26`, `skills/do-work/tools/checks/associate-files.sh:6-8`, usage strings at `internal/suiteinstall/suite_commands.go:286,332`, history comments at `internal/suitemanifest/suite_manifest.go:2` and `internal/managedsection/managed_section.go:4`. The capture's "about 34 files" counted `do-work/` records and changelogs.
- `.gitattributes`, `suite/modules.tsv`, the root `justfile`: none names root `tools/`.
- `skills/do-work/actions/version.md` and the Go updater: run the installed `<skill-root>/tools/do-work-cli.sh`, never a root script (D-01, D-02).
- `_dev/primes/lessons-shell-commands.md:119`: a dated lesson about the package's own script; left as written.

**Baseline at HEAD 792f981c (main tree, 2026-11-09):** `folder-structure-check` reports 7 gaps, including `FOLDER-STRUCTURE-ROOT-ENTRY tools/`. Focused suites pass: `suite-manifest-contract.sh` (2 s), `contracts/replace-text-section.sh` (7 s), `audit-lockins.sh` (under 1 s), `DO_WORK_MAINTAINER_TIER=heavy staged-skills-contract.sh` (34 s), `install-suite-behavior.sh` (30 s), `update-script-behavior.sh` (67 s), and the two Go tests above (2 s together).

**End-to-end check of the deletion (scratch clone, 2026-11-10):** a throwaway clone of 792f981c got the whole plan by path substitution: git rm of the five root files (git rm also removed the empty `tools/` folder), the new bootstrap line in the README and the staged installer, D-05, D-06, D-07, and the test edits above including the repointed hostile plants. Results: (1) the GREEN probe passed in 16 s, and it fails at 792f981c; (2) `DO_WORK_MAINTAINER_TIER=heavy` staged-skills, install-suite and update-script suites all passed, the update suite with no "No such file" line; (3) a fresh install through `source/skills/do-work/tools/install-do-work-suite.sh` from an archive of the clone (version set to 0.316.0) exited 0; (4) a project installed from the 792f981c archive through the old bootstrap path (v0.315.0) updated through its installed `do-work-update.sh` to that archive, served from 127.0.0.1, with exit 0 and `Updated to v0.316.0`, and its installed installer then printed the new bootstrap; (5) the old snippet path against the new archive exited 127 before any write.

*Generated by the pre-dispatch agent (explored directly, no Explore subagent)*

## Scope

**Files I will touch:**
- `skills/do-work/tools/install-do-work-suite.sh` (modify) — the bootstrap heredoc runs source/skills/do-work/tools/install-do-work-suite.sh; the root-copy preamble fallback and its comment go (D-07)
- `README.md` (modify) — the installation block names source/skills/do-work/tools/install-do-work-suite.sh
- `skills/do-work/tools/do-work-cli-preamble.sh` (modify) — one candidate beside the preamble; comment no longer names a root copy (D-07)
- `skills/do-work/tools/fetch-upstream-archive.sh` (modify) — root-copy preamble fallback and comment removed (D-07)
- `skills/do-work/tools/replace-text-section.sh` (modify) — root-copy preamble fallback and comment removed (D-07)
- `skills/do-work/tools/validate-suite-manifest.sh` (modify) — root-copy preamble fallback and comment removed (D-07)
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard.go` (modify) — shipped root tools removed, comment follows (D-05)
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard_test.go` (modify) — the installer-tools case is deleted (D-05)
- `_dev/tests/heavy-lanes.json` (modify) — root rows repointed to the staged paths, one duplicate dropped (D-06)
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_maintainer_tree_test.go` (modify) — coverage rows use the staged paths and their lanes (D-06)
- `_dev/tests/staged-skills-contract.sh` (modify) — the retained bootstrap loop and the root mirror loop are deleted; the four legacy root tools entries collapse into one entry for the whole root tools folder
- `_dev/tests/install-suite-behavior.sh` (modify) — installer path is the staged one; the fixture archive gets no root tools folder
- `_dev/tests/update-script-behavior.sh` (modify) — script paths are the staged ones, build_suite_tree gets no root tools folder, the two hostile plants move to the staged folder
- `_dev/tests/suite-manifest-contract.sh` (modify) — validator path
- `_dev/tests/contracts/replace-text-section.sh` (modify) — tool path
- `_dev/tests/audit-lockins.sh` (modify) — root tools scan roots and the root preamble exemption removed
- `_dev/primes/prime-releases.md` (modify) — the shipped roots sentence no longer names tools

Deletion: the five root copies (install-do-work-suite, fetch-upstream-archive, do-work-cli-preamble, replace-text-section and validate-suite-manifest in the root tools folder) are removed with git rm in the same commit as the edits above; the root tools folder disappears and no cmd folder is created.

**Files I will NOT touch:** `.gitattributes`; `suite/modules.tsv`; the root `justfile`, `skills/do-work-board/justfile.template` and `_dev/tests/contracts/queue-kanban.sh` (REQ-717 owns them); `skills/do-work/actions/version.md`; the Go updater, validator and fetcher packages (D-02); history comments in `suite_manifest.go` and `managed_section.go`; `VERSION`, both `CHANGELOG.md` files and the version line in `skills/do-work/actions/version.md` (the finalizer writes the release); no `tools/` forwarder (D-01); no `cmd/` folder (D-04).

**Acceptance criteria (restated from REQ, as changed by the maintainer's D-04 choice):**
- [ ] The five root copies are deleted with git rm; no root `tools/` folder and no `cmd/` copy remain.
- [ ] `folder-structure-check` no longer reports `FOLDER-STRUCTURE-ROOT-ENTRY tools/`.
- [ ] The README bootstrap, byte-identical to `--print-bootstrap-command`, runs `source/skills/do-work/tools/install-do-work-suite.sh` and installs the suite in a fresh temp Git repository from an archive built from the branch.
- [ ] An update from the previous release works: a project installed from the HEAD-of-main archive (v0.315.0) updates to an archive built from the branch through its installed updater.
- [ ] No maintainer check names a root `tools/` path: the mirror check is deleted, the heavy-lane rows and the release guard name only existing paths.
- [ ] The contract suites pass: staged-skills, install-suite, update-script, suite-manifest, replace-text-section, audit-lockins, and the two Go tests.
- [ ] Decision D-01 (no forwarder, with the proof) is recorded in this REQ.

## Pre-Flight

**Git:** ✓ Clean outside `do-work/` at 792f981c; only the claimed REQ files, `baseline.json` and the run directory are changed
**Tests baseline:** ✓ Focused pre-flight probe green (`bash do-work/runs/2026-11-09-layout-queue/REQ-716-preflight-probe.sh`); repository gate `bash _dev/tests/maintainer-verify.sh` exit 0 in 129s at 792f981c, run from a detached checkout
**Dependencies:** ✓ None added
**Base moved:** the builder branches from dc5fd2df (REQ-717's 0.315.1 release), not 792f981c. The repository gate was green at REQ-717's merge 36012976 (162s), and dc5fd2df adds only release files and `do-work/` records.

*Checked by work action*

## Discovered Tasks

- `docs/decisions/records/adr-019-four-skill-suite-contract.md` lines 11, 42 and 67 name the deleted root `tools/validate-suite-manifest.sh`, and line 67 links to it. ADR-016 and ADR-018 already carry similar dead links to a root `tools/queue-kanban/`. Decision records are history and are not shipped. impact-negligible → report only
- `_dev/tests/staged-skills-contract.sh` (about lines 976-981) loops over `git ls-files tools/queue-kanban`, which has been empty since the board moved under `skills/do-work-board/`, so the loop checks nothing. impact-negligible → report only

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/install-do-work-suite.sh` (modified)
- `README.md` (modified)
- `skills/do-work/tools/do-work-cli-preamble.sh` (modified)
- `skills/do-work/tools/fetch-upstream-archive.sh` (modified)
- `skills/do-work/tools/replace-text-section.sh` (modified)
- `skills/do-work/tools/validate-suite-manifest.sh` (modified)
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard_test.go` (modified)
- `_dev/tests/heavy-lanes.json` (modified)
- `skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_maintainer_tree_test.go` (modified)
- `_dev/tests/staged-skills-contract.sh` (modified)
- `_dev/tests/install-suite-behavior.sh` (modified)
- `_dev/tests/update-script-behavior.sh` (modified)
- `_dev/tests/suite-manifest-contract.sh` (modified)
- `_dev/tests/contracts/replace-text-section.sh` (modified)
- `_dev/tests/audit-lockins.sh` (modified)
- `_dev/primes/prime-releases.md` (modified)

Deletion: the five root copies (install-do-work-suite, fetch-upstream-archive, do-work-cli-preamble, replace-text-section and validate-suite-manifest in the root tools folder) were removed with git rm in the same commit as the edits above; the root tools folder is gone and no cmd folder exists.

**What was done:** The README bootstrap and the installer's `--print-bootstrap-command` text now run `source/skills/do-work/tools/install-do-work-suite.sh` from the downloaded archive. The staged launchers lost their fallback to a root copy. The release guard counts the declared modules plus `suite` as shipped roots. Every maintainer suite runs the staged scripts, the heavy-lane rows name the staged paths, and the root mirror check is deleted.

## Qualification

**Range:** `5c74f82a..68e1b00b` (one merge of the REQ branch: 66ca00c0 RED coverage rows, 7610175f deletion and repointing)
**Mechanical gates:** qualify ✓ satisfied, scope-drift ✓ satisfied (17 declared paths plus the five deletions, all touched).

**Requirement trace (orchestrator, read against the diff):**
- Deletion: the five root copies are removed with git rm; no root `tools/` and no `cmd/` folder exist at the merge.
- Bootstrap: `README.md` and the installer's heredoc both run `source/skills/do-work/tools/install-do-work-suite.sh`; the install-suite test checks they are byte-identical.
- Older installs (D-01): no updater of any release runs a script from the downloaded archive; the builder's update from 0.315.1 through the installed `do-work-update.sh` reached the new version with exit 0.
- Release guard (D-05): shipped roots are the declared modules plus `suite`; the case for a root installer path is deleted with the folder.
- Maintainer checks: heavy-lane rows name the staged paths (D-06), the root mirror check and the root scan roots are gone, and the two hostile-archive fixtures plant their scripts in the staged folder, so they still test something.

**Live flow:** fresh install through the new bootstrap path exit 0; update from 0.315.1 exit 0 with `Updated to v0.316.0`; the old README snippet against the new archive exits 127 before any write.

*Qualified by work action*

## Testing

**Tests run:** `bash do-work/runs/2026-11-09-layout-queue/REQ-716-probe.sh` (no root `tools/` path, the bootstrap texts, the release guard and heavy-lane coverage tests), then `bash _dev/tests/maintainer-verify.sh` from a detached checkout at the merge revision. The builder ran an end-to-end install and update with local archives; the reviewer ran the README bootstrap line against an archive of the merge.
**Result:** ✓ Focused probe green. Repository gate exit 0 at the merge 68e1b00b (gate wall 144s). Fresh install through `source/skills/do-work/tools/install-do-work-suite.sh` exit 0 (builder and reviewer). An install of 0.315.1 updated through its installed `do-work-update.sh` to the branch archive with exit 0. The old snippet path against the new archive exits 127 before any write.

**Red-green validation:**
- `TestRepositoryManifestNamesEveryLaneScopedMaintainerEntryPoint`: ✗ at 66ca00c0 (all five staged-path rows, e.g. `skills/do-work/tools/do-work-cli-preamble.sh selected lanes = [staged-skills], want [staged-skills updater installer]`) → ✓ at 7610175f.
- `REQ-716-probe.sh`: ✗ at 4845a9a7 (`a root tools/ path still exists`) → ✓ at 7610175f.

**Existing tests updated (cross-REQ impact):**
- `finalization_release_guard_test.go`: the case for a root installer path is deleted with the folder (D-05).
- `_dev/tests/update-script-behavior.sh`: the two hostile-archive fixtures plant their scripts in `skills/do-work/tools/`, where real files live, so they still assert something.
- `_dev/tests/staged-skills-contract.sh`, `install-suite-behavior.sh`, `suite-manifest-contract.sh`, `contracts/replace-text-section.sh`, `audit-lockins.sh`: run the staged scripts; the root mirror check and the root scan roots are deleted.

**Heavy verification plan:**
- Range: 5c74f82a9b98e2ecc12226e941fed54898b88117..68e1b00b940fbde618272eade66dfcaba26337ab
- All six lanes are selected because coverage is uncertain for the five deleted root `tools/` paths.
- queue-kanban-javascript: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript`
- queue-kanban-browser: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser`
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations`
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills`
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater`
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer`

*Verified by work action*

## Review

**Overall: 96%** | 2026-11-10T01:15:15Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 96% |
| Test Adequacy | 88% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Verdict:** Approve. The five root installer copies in `tools/` are deleted. The README bootstrap and `--print-bootstrap-command` now run `source/skills/do-work/tools/install-do-work-suite.sh`. No installed updater of any release runs a script from the downloaded archive, so already-installed consumers are not affected. Every maintainer suite that used a root copy now runs the staged script and still asserts something.

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:**
- M1: No gate test runs the installer from an extracted archive, which is the exact path the new bootstrap uses. `_dev/tests/install-suite-behavior.sh:200-215` only checks that the README block is byte-identical to the heredoc. Its install cases run the repository's own `skills/do-work/tools/install-do-work-suite.sh` with `--archive`. A future change to how the launcher finds its preamble or `do-work-cli.sh` from an extracted tree would break every fresh install while the gate stays green. This gap existed for the old root path too. This REQ moved the entry point, and only run-local evidence proves the new one (the builder's e2e and this review's fresh install). Fix: after `tar czf "$archive_file"` (about line 172), extract `$archive_file` into `$workdir/bootstrap-source` with `--strip-components=1`. Then run `printf 'y\n' | bash "$workdir/bootstrap-source/skills/do-work/tools/install-do-work-suite.sh" --project-root <fresh git repo> --archive "$archive_file"`, and assert exit 0 and `assert_four_modules`. — impact-negligible → report only
- M2 (stale restatement, inherited from REQ-717, which moved the run-kanban and kanban-static recipe logic into board scripts; it is REQ-717's N1 and is still unfixed at the wave end): `skills/do-work-board/actions/board.md:58` says the git-optional root fallback is "matching `justfile.template`". The template recipes now only pass `{{justfile_directory()}}`, and the git-optional handling lives in `skills/do-work-board/tools/kanban-static.sh:16-17`. Fix: change "matching `justfile.template` and" to "matching `tools/kanban-static.sh` and". — impact-negligible → report only

**Nit findings:**
- N1: The REQ title and slug still say "Move the root installer scripts from tools/ to cmd/". The maintainer chose deletion (D-04), so the archived record and the `[REQ-716]` commit subjects describe a move that did not happen. The Triage note explains this, but a title-only reader will not see the note. Fix: before archiving, set `title: 'Delete the root installer copies in tools/; bootstrap runs the staged installer'`. The slug can stay. — impact-negligible → report only

**Requirements checklist:**
- [x] Map every caller before moving. Delivered: the Exploration section lists the README bootstrap, the heredoc, the release guard, the heavy lanes and six maintainer suites. The "about 34 files" in the capture were `do-work/` records and changelogs, and the reviewer's sweep agrees.
- [x] Rule out a break for already-installed consumers, and record which way (D-01). Delivered and spot-checked against old code. `bef7334f:tools/do-work-update.sh:50,148` (the bridge updater) runs only `$script_dir/validate-suite-manifest.sh`, the installed copy, with `--root "$fresh_upstream"` passed as data. `0e8cf0d9:skills/do-work/tools/do-work-update.sh:43-45,75,82,103` (the shell-suite updater) runs the installed fetcher, validator and installer from `$script_dir`. `0e8cf0d9:skills/do-work/tools/install-do-work-suite.sh:74-77` resolves its helpers from `$script_dir` as well. `7f1603bb:tools/do-work-update.sh` (single-skill) runs nothing from `$fresh_upstream`. `916d2706` and `eb1a2c57` call the installed `do-work-cli.sh update-suite`, and `update_transaction.go:83-121` fetches, validates and installs in-process. No shipped action ever ran the printed bootstrap (`git log -S'print-bootstrap-command' -- skills/*/actions` finds no consumer). D-01 holds.
- [x] Layout rule 5 / no root `tools/` entry. Delivered: the probe's step 3 shows `folder-structure-check` no longer reports `FOLDER-STRUCTURE-ROOT-ENTRY tools/`. Deletion replaces the `cmd/` move by the maintainer's choice (D-04).
- [x] README bootstrap installs the suite in a temp repository. Delivered, and reproduced by this review (see Acceptance).
- [x] Update from the previous release works. Delivered by the builder's e2e (0.315.1 to 0.316.0 through the installed `do-work-update.sh`, exit 0). This review did not re-run it. The code reading under D-01 shows the update path never touches a root archive script.
- [x] Contract suites pass. Delivered and re-run (see Acceptance).

**Specific judgments:**
- (1) D-01: confirmed by `git show` of five updater generations (listed above). No updater executes a script from the downloaded archive. The only archive-relative caller was the README bootstrap, which changes in the same commit.
- (3) Release guard (`finalization_release_guard.go:44-47`): the shipped roots are now `DeclaredModuleSources` plus `suite`. A `_dev/tests` path and a metadata-only release are still refused. A path under any declared module (with or without its own VERSION) or under `suite/` is still accepted. `TestReleaseRefusedWhenTheImplementationShipsNothing` passes (via the probe), with all five remaining cases. A root `tools/...` path now falls under no root and is refused, which is correct now that the folder is retired.
- (4) Maintainer tests: at base, every root-copy reference in `_dev/` and the Go tests (install-suite-behavior.sh 10/120/171-173, update-script-behavior.sh 14-17/204-207/419-420/458-459/560-562/792, suite-manifest-contract.sh 6/79, contracts/replace-text-section.sh 9/11, audit-lockins.sh 48/64/149, heavy-lanes.json 58-62/78-81, staged-skills-contract.sh 83-96/119-126/861-871, finalization_release_guard_test.go 53, heavy_maintainer_tree_test.go 259-263) is now repointed to the staged script or deleted with a reason (D-05, D-06, the mirror loop). The two hostile-archive fixtures (`update-script-behavior.sh:414-415` and `:453-454`) now write into `skills/do-work/tools/`. `build_suite_tree` creates that folder (line 210) and copies the real scripts there, so each plant overwrites a real file rather than failing on a missing folder. The suites then assert `archive-installer-ran` is absent and that the traversal is rejected "despite bundled override". The heavy run logs show zero "No such file" lines.
- D-06 (repoint the heavy-lane rows rather than remove them) is better than the coordinator's "remove them", because the rows are what select the updater and installer lanes for these five scripts. D-07 (drop the dead `../skills/do-work/tools/` fallback) and D-08 are sound. D-09 is acceptable, because `staged-skills-contract.sh:93` now fails on any file under a root `tools/`.

**Acceptance:** Pass. Covered stages: implementation (diff walk of `5c74f82a..68e1b00b`) and integration. The GREEN probe passed from the main tree in 15 s. In a scratch clone at 68e1b00b with `DO_WORK_MAINTAINER_TIER=heavy`: `update-script-behavior.sh` passed (69 s), `install-suite-behavior.sh` (29 s) and `staged-skills-contract.sh` (38 s), with zero "No such file" lines. Deployment, simulated locally: `git archive --prefix=skill-do-work-main/ 68e1b00b` (root entries: CHANGELOG.md, README.md, VERSION, justfile, skills, suite; no `tools/`, no `cmd/`) was extracted, and the README's install line `bash "$bootstrap_tmp/source/skills/do-work/tools/install-do-work-suite.sh" --project-root "$project_root" --archive "$archive_file"` was run against a fresh temp Git repository. It exited 0 and installed four modules at v0.315.1. The installed copy prints the new bootstrap path.
**Restatement sweep:** redefined (1) the fresh-install bootstrap entry point (`source/tools/install-do-work-suite.sh` → `source/skills/do-work/tools/install-do-work-suite.sh`), (2) root `tools/` as a shipped root and home of byte-identical copies of the five installer scripts (now retired: release guard, prime-releases, mirror check, heavy lanes), and (3) the launchers' root-copy preamble fallback. Inherited from REQ-717: where the `run-kanban`/`kanban-static` recipe logic lives, and the REQ-551 delegation lock-in wording. Swept `skills/`, `README.md`, `_dev/`, `CLAUDE.md`, `suite/` and `docs/` (except `docs/changelog/` and `docs/decisions/records/`). No text names a root `tools/` script path, `source/tools/`, a `cmd/` home for these scripts, a "root copy"/"staged mirror"/"preamble pair", or `tools/` as a shipped root. The remaining `tools/<script>.sh` mentions are package-relative (`skills/do-work/docs/prescribed-shell-primitives.md:26`), Go history comments (`suite_manifest.go:2`, `managed_section.go:4`), or dated KB and lesson records. `cmd/` appears only as the generic layout rule (`README.md:126`, `folder-structure.md:47`). The REQ-551 lock-in comment (`audit-lockins.sh:60-62`) matches the merged check. Inline recipe bodies survive only as synthetic fixtures (`update-script-behavior.sh:309-313`). One stale restatement was found: M2 (`board.md:58`). The builder's two discovered tasks (ADR-019 dead link, empty `tools/queue-kanban` loop at `staged-skills-contract.sh:976-981`) still stand and are not duplicated here.
**Suggested testing:** 3 items
- Live acceptance: unassessed. After this release reaches GitHub `main`, paste the README Installation block into a fresh Git repository and confirm it installs. This review ran the install line against a local `git archive`, not the GitHub tarball the block downloads.
- Deployment (update path): unassessed by this review. Re-run the builder's update from an install of the prior release through its installed `do-work-update.sh` against the released archive, and confirm `Updated to v<new>`. The builder ran this against a local archive served from 127.0.0.1.
- Regression scenario: confirm that a saved copy of the old README snippet fails with exit 127 before any project write. The builder observed this, and the changelog entry should tell readers to copy the current README block.
**Follow-ups created:** None (3 findings report only)

### Domain Review
`crew-members/general.md` defines no review rubric. Its discipline rules held: the builder recorded the ADR-019 dead link and the empty `tools/queue-kanban` loop as Discovered Tasks and did not fix them inline. The PLAN, APPLY and UNIFY boxes are all checked. Decisions D-01 to D-09 are in the REQ and the hand-back. The one unplanned wording change (the audit-lockins comment and failure text) is declared under the hand-back's Scope changes.

*Reviewed by review-work action*

## Heavy Verification Plan

- Base revision: 5c74f82a9b98e2ecc12226e941fed54898b88117
- Target revision: 68e1b00b940fbde618272eade66dfcaba26337ab
- queue-kanban-javascript: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — coverage is uncertain for the five deleted root `tools/` paths
- queue-kanban-browser: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — coverage is uncertain for the five deleted root `tools/` paths
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — coverage is uncertain for the five deleted root `tools/` paths
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — coverage is uncertain for the five deleted root `tools/` paths
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — coverage is uncertain for the five deleted root `tools/` paths
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — coverage is uncertain for the five deleted root `tools/` paths

## Heavy Verification Result

- Target revision: 68e1b00b940fbde618272eade66dfcaba26337ab
- Execution revision: 68e1b00b940fbde618272eade66dfcaba26337ab (detached checkout)
- queue-kanban-javascript: exit 0, 8s, executed
- queue-kanban-browser: exit 0, 84s, executed (real browser, not skipped)
- do-work-cli-integrations: exit 0, 71s, executed
- staged-skills: exit 0, 35s, executed
- updater: exit 0, 66s, executed
- installer: exit 0, 27s, executed

## Lessons Learned

**What worked:** Reading every updater generation with `git show` before deciding on a forwarder (D-01): it settled the riskiest question without guessing, and the end-to-end update from 0.315.1 confirmed it. Asking the maintainer whether to move or delete byte-identical copies turned a move into a deletion of 323 lines.
**What didn't:** The first plan moved the copies to `cmd/` because the audit said so; nobody had checked that they duplicated the staged scripts. After the deletion, two hostile-archive fixtures would have written into a folder that no longer exists and kept passing while testing nothing; the pre-dispatch caught it.
**Worth knowing:** The bootstrap line lives twice, in `README.md` and in the installer's `--print-bootstrap-command` heredoc; `install-suite-behavior.sh` checks they are byte-identical. The release guard's shipped roots are now the declared modules plus `suite`; a new root folder of shipped scripts would need adding there. Heavy-lane rules are exact paths, so a moved script must be repointed, not dropped, or it stops selecting its lanes.

## Orientation

Now the suite installs from the packaged installer inside the downloaded archive (`source/skills/do-work/tools/install-do-work-suite.sh`), and the repository root has no `tools/` folder. This lives in the installer and its bootstrap text (`_dev/primes/prime-releases.md`, `_dev/primes/prime-shell-commands.md`). `[MAP CHANGED]` The five root copies are gone, so the staged scripts in `skills/do-work/tools/` are the only copies, and the release guard no longer counts a root `tools/` as shipped. Prime staleness check: `prime-releases.md` was updated in this change (shipped roots); `prime-shell-commands.md` names no root `tools/` path.

## Timing

Observed 2026-11-09T23:40:42Z to 2026-11-10T00:56:32Z: 1h 15m 50s total, 22m 34s attributed across 3 events, 53m 16s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| exploration-preflight | 12m 51s | 1 |
| planning | 9m 43s | 2 |

Slowest stage: exploration-preflight / caller map, release-history proof, baselines, probes, end-to-end simulation, 12m 51s, outcome success.
