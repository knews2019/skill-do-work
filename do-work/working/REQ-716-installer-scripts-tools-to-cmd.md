---
id: REQ-716
title: 'Move the root installer scripts from tools/ to cmd/'
status: claimed
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
---
# Move the root installer scripts from tools/ to cmd/
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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
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
<!-- D-XX counter: last used D-07. Next decision: D-08. -->

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
