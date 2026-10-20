---
id: REQ-705
title: 'Confirm: a second opinion refuses to run while a project preset file is uncommitted'
status: completed
created_at: 2026-10-20T20:47:28Z
user_request: UR-160
addendum_to: REQ-703
builder_decided: true
domain: backend
prime_files: ["skills/do-work/tools/do-work-cli/prime-do-work-cli.md", "_dev/primes/prime-releases.md"]
required_lessons: ["_dev/primes/lessons-releases.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-703, REQ-704]
write_set: ["skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index_test.go", "skills/do-work-toolbox/docs/ai-report-guide.md", "skills/do-work-toolbox/actions/ai-report.md", "CHANGELOG.md", "skills/do-work/CHANGELOG.md", "VERSION", "skills/do-work/VERSION", "skills/do-work/actions/version.md", "skills/do-work/tools/do-work-cli/lessons-do-work-cli.md", "do-work/lessons-index.md"]
status_changed_at: 2026-10-20T21:19:07Z
claimed_at: 2026-10-20T21:19:46Z
route: B
estimate:
  p50_active_minutes: 35
  confidence: medium
  calculated_at: 2026-10-20T21:24:16Z
  basis:
    - Route B
    - 13-file write set
    - 2 subsystems involved
    - 7 acceptance criteria
dispatch_at: 2026-10-20T21:31:30Z
builder_handback_at: 2026-10-20T21:35:28Z
integration_at: 2026-10-20T21:36:05Z
review_at: 2026-10-20T21:49:36Z
commit: 5cbf1f10f0a0ebc7c2d6eeabe385cff6a22abb41
heavy_verified_at: 2026-10-20T21:54:03Z
heavy_verified_revision: 5cbf1f10f0a0ebc7c2d6eeabe385cff6a22abb41
kb_status: pending
completed_at: 2026-10-20T21:55:30Z
release_at: 2026-10-20T21:55:30Z
---

# Confirm: A Second Opinion Refuses To Run While A Project Preset File Is Uncommitted

## What
During REQ-703 (second opinions accept a named preset and a free-text editor direction), the plan kept one existing safety check unchanged: before and after each second-opinion run, the launcher refuses to continue when Git shows any change outside the report bundle. A project preset file, `ai-reports/second-opinion-presets/<name>.md`, sits outside the bundle. So the first run after you write or edit a project preset is refused until you commit that file. This was escalated because it trades convenience against the one check that catches a backend writing files outside the bundle, and only you can judge which matters more.

## What the Builder Chose
Keep the check whole. A run with an uncommitted project preset fails with `outside-bundle changes: ai-reports/second-opinion-presets/<name>.md`, and `skills/do-work-toolbox/docs/ai-report-guide.md` tells you to commit a new preset before using it. Shipped presets (`bold`, `visual`) are not affected, because they live in the installed toolbox, not in your project.

## What Would Change
If you prefer the exemption: the outside-bundle check in `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` would ignore the one preset file the run actually read, and only when its content is unchanged between the check before the run and the check after it. One test would pin that a backend editing the preset is still refused.

## Open Questions
- [x] Should a second opinion run while the project preset file it uses has uncommitted changes? → Neither: presets live only inside the report bundle (bundle-local presets folder) plus the shipped ones; drop the shared ai-reports/second-opinion-presets/ folder
  Recommended: No, keep refusing until the preset is committed (already implemented).
  Value: the check that catches any backend write outside the report bundle stays whole, with no special case.
  Risk: the first run after you write or edit a preset is refused, and you must commit the file and run again; reversible later by adding the exemption.
  Also: Yes, exempt the preset file the run read, as long as its bytes did not change during the run.


## Answer Notes

- 2026-10-20 - [ ] Should a second opinion run while the project preset file it uses has uncommitted changes?: Neither: presets live only inside the report bundle (bundle-local presets folder) plus the shipped ones; drop the shared ai-reports/second-opinion-presets/ folder
> ```
> User, asked whether a run should refuse while a project preset in ai-reports/second-opinion-presets/ is uncommitted: "why would you put the second opinion presets outside the report folder?"
> User, asked where their own presets should live (options: shared folder + commit first / also inside the bundle / only inside the bundle): "Only inside the bundle".
> Reasoning: a preset kept inside the report bundle is never outside-bundle dirt, so the question disappears instead of needing a carve-out. Out of scope now: the shared project folder ai-reports/second-opinion-presets/ (drop its lookup and the catalog skip that only existed for it); shipped presets beside the prompt stay.
> ```

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Follow brief steps 1-6 exactly. Write the new test first and prove RED, then change the lookup slice, remove the catalog skip, move the constant, delete the skip test and the project block, then the two prose sentences (from the builder hand-back).
- [x] **[APPLY]:** As planned, inside the declared scope: the `folders` slice now starts at `<bundle>/second-opinion-presets`; the constant moved to `report_second_opinion.go`; the `reportIndexWalk` skip, its doc clause and its test are gone; the project-preset block left `TestSecondOpinionPresetResolutionAndArgumentRules`; guide :91, :107 and `ai-report.md` :73 edited. One test-only tidy after GREEN (each bundle preset written inside its loop); RED was not re-run after it, case (a) is unchanged (from the builder hand-back).
- [x] **[UNIFY]:** `git diff 11780ac0 --stat`: 6 files changed, 40 insertions(+), 48 deletions(-); `gofmt -l skills/do-work/tools/do-work-cli/internal` empty; `go vet` on `./internal/toolboxcommands` exit 0; `git diff --check` exit 0; the REQ-705 probe exit 0; package selection `SecondOpinion|AIReportIndex|TestHandlersRegisterCanonicalToolboxCommands` 32 PASS, 0 FAIL, 0 SKIP. All six files read back in the diff (from the builder hand-back).

---

## Triage

**Route: B** - Medium

**Reasoning:** The outcome is fixed by the user's answer and the code sits in two known functions, but the change still crosses the preset lookup, the report catalog walk, two test files and two prose files. Exploration found the exact lines, so no separate plan is needed.

**Planning:** Not required

**What the work is:** this REQ was minted as a builder-decided confirm question (keep refusing an uncommitted `ai-reports/second-opinion-presets/<name>.md`, or exempt it). Its title, What, What the Builder Chose and What Would Change describe that old question and are kept as history only. The user's answer (`## Answer Notes`, 2026-10-20) overrides both choices and is the work: a user preset lives only inside the report bundle, at `<bundle>/second-opinion-presets/<name>.md`, plus the shipped presets beside the prompt; the shared `ai-reports/second-opinion-presets/` lookup and the `reportIndexWalk` skip that existed only for it are removed. A bundle preset is never outside-bundle dirt, so the outside-bundle check stays whole with no exemption and the guide's "commit a new preset first" sentence goes.

**Fields judged here (the follow-up template had none):** `domain: backend` (Go launcher plus toolbox prose); `tdd: true` (behavior-changing, and the lookup is reachable from an existing handler test); `impact: impact-user-visible` (where a user puts a preset changes); `effort_estimate: effort-substantive` (small, but it changes behavior and a test contract).

**Open questions (Step 3.5):** the one item is already `- [x]`, answered by the user. Nothing to resolve. Builder choices below are D-01 onwards.

### Red-Green Proof
**RED prompt/case:** in the second-opinion fixture, write an uncommitted `ai-reports/2026-10-15_1200-example/second-opinion-presets/mine.md` and run `ai-report-second-opinion --prepare --backend claude --model claude-fable-5-1 --preset mine`. Today it is refused with `unknown preset "mine"; presets found: bold, visual`, because the lookup reads only `ai-reports/second-opinion-presets/` and the shipped folder.
**Why RED now:** `handleAIReportSecondOpinion` builds its preset folders as `<repo>/ai-reports/second-opinion-presets` then the shipped folder (`report_second_opinion.go:155`); the bundle is never searched.
**GREEN when:** (1) the same call is `prepared` and `preset_file` is `ai-reports/2026-10-15_1200-example/second-opinion-presets/mine.md`, with the file still uncommitted; (2) a bundle `bold.md` wins over the shipped `bold.md` and its text reaches the prompt; (3) a committed `ai-reports/second-opinion-presets/legacy.md` is no longer found: `--preset legacy` is refused with `unknown preset "legacy"; presets found: bold, mine, visual`; (4) shipped-only presets and plain runs behave as in 0.306.6.
**Validation:** derived from the user's answer (2026-10-20); the coordinator named this RED/GREEN pair.

<!-- D-XX counter: none used in Step 3.5. Next decision: D-01. -->

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 22858 tokens, `slugged: partial` so no targeted form; matches "a session that is its own second-opinion backend" and the catalog walk. Pre-dispatch read its REQ-703 `writer-anchored-marker-position` and `session-is-the-backend` bullets by hand; neither constrains a lookup-folder change.
- `_dev/primes/lessons-action-files.md` — 8398 tokens, `slugged: partial`; weak match (one clause in `skills/do-work-toolbox/actions/ai-report.md`), no routing or status contract changes.

## Exploration

**Lessons consult:** `required_lessons` refreshed to `_dev/primes/lessons-releases.md` (666 tokens, whole file, read). Its `canonical-link-outlives-its-target` family matters at UR close, not here: the REQ-703 lesson bullet in `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md:95` links the flat path `do-work/archive/REQ-703-...md`, which moves when REQ-704's integrator closes UR-160. Prime `prime-do-work-cli.md` § Traps, `guard-on-one-entry-path`: the preset lookup has one entry, in `handleAIReportSecondOpinion`, shared by the Session prepare, Targeted and Batch forms; `--finalize` reads the direction from the prepare record and never looks presets up. So one change covers every form.

**Key code (main at 11780ac0, after REQ-703 / 0.306.6):**
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go:153-188`, in `handleAIReportSecondOpinion`: the `folders` slice at :155 is `<repo>/ai-reports/second-opinion-presets`, then `filepath.Dir(templatePath)/second-opinion-presets`. The comment at :154 says a project preset wins. Lstat regular-file check, unknown-name list (sorted, de-duplicated stems of both folders, `none` when empty), trimmed text, `{{` refusal, and the repository-relative `preset_file` all stay as they are. `bundle` is already resolved (absolute, after the supersession chain) at :144, before this block.
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index.go:28-30`: constant `secondOpinionPresetsFolder` and its comment in the shared const block; :83 doc comment "except the second-opinion presets folder"; :94 the skip `&& entry.Name() != secondOpinionPresetsFolder`. After the change the constant's only user is `report_second_opinion.go`.
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index_test.go:251-276`: `TestAIReportIndexSkipsSecondOpinionPresetsFolder`, which exists only for the skip.
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go:842-876`: `TestSecondOpinionPresetResolutionAndArgumentRules`; its block at :853-862 commits `ai-reports/second-opinion-presets/bold.md` and asserts it wins. `secondOpinionFixture` (:21) builds the bundle `ai-reports/2026-10-15_1200-example` and returns `(repo, bundle)`; `secondOpinionFixturePresets` (:828) writes the shipped `bold`/`visual` beside the temp template. `TestSecondOpinionDirectedRunNamesAndMarksPage` (:881) uses shipped presets only and is unaffected.
- Preservation: `secondOpinionSnapshot` (:413) walks the bundle recursively, and `secondOpinionPreserved` (:447) refuses any change to a pre-existing path. A bundle preset is a pre-existing bundle file, so an editor that edits it is refused with `backend changed existing bundle path`. No new code.
- Outside-bundle check: `secondOpinionOutside` (:475) counts only paths outside the bundle prefix, so an uncommitted bundle preset is never outside dirt.
- Catalog: `pickReportEntry` (`report_index.go:171`) reads only the bundle root, and a second-opinion bundle always has `index.html`, so a `second-opinion-presets/` subfolder never changes a bundle's entry.

**Prose:** `skills/do-work-toolbox/docs/ai-report-guide.md:91` says the catalog skips "the `second-opinion-presets/` folder"; :107 ends with the project-folder sentence and "commit a new one before using it, because the run refuses uncommitted changes outside the bundle". `skills/do-work-toolbox/actions/ai-report.md:73` names `preset_file` (the file used) but no location. `skills/do-work-toolbox/SKILL.md` carries only the trigger words; no change. Changelog history entries for 0.306.6 stay as written.

**Decisions:**
- **D-01** (DECIDE & STATE): lookup order is `<bundle>/second-opinion-presets/<name>.md`, then the shipped folder. `<bundle>` is the resolved bundle the run writes into (after `latest` and the supersession chain). A bundle preset wins over a shipped preset of the same name, keeping 0.306.6's "your preset wins" rule. Reason: smallest change, one slice element.
- **D-02** (DECIDE & STATE): the unknown-name list is the stems of the bundle folder plus the shipped folder, sorted and de-duplicated, `none` when both are empty. This is today's code over the new folder slice.
- **D-03** (DECIDE & STATE): remove the `reportIndexWalk` skip and its doc-comment clause; move the constant `secondOpinionPresetsFolder` out of the `report_index.go` const block into `report_second_opinion.go` (its only user), with a comment naming both homes: the bundle and beside the prompt.
- **D-04** (DECIDE & STATE): delete `TestAIReportIndexSkipsSecondOpinionPresetsFolder`, do not convert it. With the skip gone, `ai-reports/second-opinion-presets/` is an ordinary directory, already covered by `TestAIReportIndexCatalogsEveryBundleNamingStyleOnce`; a test that it is listed pins no real failure.
- **D-05** (DECIDE & STATE): the regression lives in a new test `TestSecondOpinionPresetResolvesFromBundleOnly` (the Red-Green Proof's cases 1-3), and the :853-862 project block leaves `TestSecondOpinionPresetResolutionAndArgumentRules`, which keeps the shipped case and the refusals. A new name makes the GREEN probe fail at base. The bundle presets stay uncommitted in the test, which pins the user's reason for the change (no commit step). The `legacy.md` case pins the real failure of a leftover lookup reading a file outside the bundle.
- **D-06** (DECIDE & STATE): a preset in a superseded revision is not seen; only the resolved bundle is searched. Value: one folder, no chain walk. Risk: after `ai-report revise` the user copies their preset into the new bundle; the unknown-name list shows what is found, so the miss is visible. Reversible.
- **D-07** (DECIDE & STATE): docs. Guide :107: replace the project-folder sentence with one sentence: your own preset goes inside the report bundle as `second-opinion-presets/<name>.md`, wins over a shipped preset of the same name, and needs no commit. Guide :91: drop "and the `second-opinion-presets/` folder". `ai-report.md` :73: the `preset_file` clause names the two places a preset comes from. No `SKILL.md` change.

*Generated by pre-dispatch exploration*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` (modify) — bundle-local preset folder in the lookup; the folder constant moves here
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go` (modify) — new bundle-only test; project block removed
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index.go` (modify) — catalog skip and constant removed
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index_test.go` (modify) — skip test deleted
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modify) — preset location sentence, catalog sentence
- `skills/do-work-toolbox/actions/ai-report.md` (modify) — preset_file location clause
- `CHANGELOG.md`, `skills/do-work/CHANGELOG.md`, `VERSION`, `skills/do-work/VERSION`, `skills/do-work/actions/version.md` (modify, integrator only) — patch release after 0.306.6
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`, `do-work/lessons-index.md` (modify, integrator only, only if a lesson is recorded)

**Files I will NOT touch:** `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (row fields unchanged), `skills/do-work-toolbox/docs/ai-report-second-opinion-prompt.md`, the shipped presets `bold.md` and `visual.md`, `skills/do-work-toolbox/SKILL.md`, `skills/do-work/tools/do-work-cli.sh`, the 0.306.6 changelog entries.

**Acceptance criteria (restated from REQ):**
- [ ] A `--preset <name>` resolves `<bundle>/second-opinion-presets/<name>.md` first, then the shipped preset; `preset_file` names the bundle file repository-relative.
- [ ] An uncommitted bundle preset runs without an outside-bundle refusal (no commit step).
- [ ] `ai-reports/second-opinion-presets/` is no longer read; a name only there is an unknown preset.
- [ ] The unknown-name list is the bundle presets plus the shipped presets.
- [ ] The `reportIndexWalk` skip and its test are gone.
- [ ] The guide's commit-first sentence is gone, and the guide and `ai-report.md` say presets live in the bundle or ship with the toolbox.
- [ ] Shipped-only presets, plain runs and all existing second-opinion and catalog tests behave as in 0.306.6.

## Pre-Flight

**Git:** ✓ clean at 11780ac0 (main after REQ-703's release 0.306.6, the REQ-705 answer and both claims), the builder's base before the run-artifacts commit; only do-work/ run artifacts were dirty when the coordinator recorded pre-flight
**Tests baseline:** ✓ `bash do-work/runs/work-2026-10-20-200259/REQ-705-preflight-probe.sh` exit 0, launched (`do-work/working/baseline.json`). Repository gate `DO_WORK_FAST_STAGE_REUSE=off bash _dev/tests/maintainer-verify.sh` at 11780ac0: run 1 exit 0 in 135 s at 1-minute load 3.7; green-gate record satisfied at revision 11780ac0
**Dependencies:** ✓ Go toolchain and ShellCheck present (the gate ran)

*Checked by the coordinator before dispatch (2026-10-20T21:27:51Z to 2026-10-20T21:30:37Z)*

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index_test.go` (modified)
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modified)
- `skills/do-work-toolbox/actions/ai-report.md` (modified)

**What was done:** `--preset <name>` now looks in `<bundle>/second-opinion-presets/<name>.md`, then in the shipped presets; the shared `ai-reports/second-opinion-presets/` folder is no longer read, and the catalog walk no longer skips it. The folder constant moved to its only user, the catalog skip test was deleted, the project-preset block left the resolution test, and the new `TestSecondOpinionPresetResolvesFromBundleOnly` pins an uncommitted bundle preset, a bundle preset shadowing a shipped one, and a refused shared-folder preset. The guide and `ai-report.md` now say a user preset lives in the bundle and needs no commit.

## Decisions

(from the builder hand-back)

- **D-08** (DECIDE & STATE): the new test loops over an ordered `[][2]string` (`mine`, then `bold`), not a map, so case (a) always runs first and the RED failure is deterministic.
- **D-09** (DECIDE & STATE): the legacy case asserts the refusal through `result.Findings[0].Evidence`, the same shape as the refusal loop in `TestSecondOpinionPresetResolutionAndArgumentRules`. No new helper.

## Discovered Tasks

(from the builder hand-back)

- **impact-user-visible** A consumer that already made `ai-reports/second-opinion-presets/` under 0.306.6 now sees it in the catalog, `index` and `find` as an undated bundle whose entry is its first `.md` file. `latest` is not affected, because undated bundles sort last. This is the intended result of D-04 (the skip is gone), but the CHANGELOG entry should tell those users to move their presets into a bundle and delete the old folder. → report only

## Qualification

**Gate records:** `advance --diff-range bac2288b..d8cb783f` → `qualify` satisfied (no debug artifact, no output-primitive finding, P-A-U boxes ticked). `scope-drift` reported six `SCOPE-DECLARED-NOT-TOUCHED` warnings: `CHANGELOG.md`, `VERSION`, `skills/do-work/CHANGELOG.md`, `skills/do-work/VERSION`, `skills/do-work/actions/version.md`, `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`. All six are release and lesson files the integrator writes at finalization, so they are expected, not drift. No touched-but-undeclared path.

**Requirement trace (read against `git diff bac2288b..d8cb783f` and the merged files):**
- Bundle first, then shipped; `preset_file` repository-relative: the `folders` slice in `handleAIReportSecondOpinion` is now `filepath.Join(bundle, secondOpinionPresetsFolder)`, then the shipped folder. `bundle` is resolved by `secondOpinionBundle` before the preset block (after `latest` and the supersession chain). The existing `filepath.Rel` step makes `preset_file` repository-relative. Pinned by case (a) `mine` and case (b) `bold` (the bundle file shadows the shipped one, and its text reaches the prompt).
- Uncommitted bundle preset runs: the test writes both bundle presets without committing them and gets `state: prepared`. `secondOpinionOutside` counts only paths outside the bundle, so no exemption was needed.
- Shared folder no longer read: case (c) commits `ai-reports/second-opinion-presets/legacy.md` and gets `unknown preset "legacy"`.
- Unknown-name list is bundle plus shipped: the same refusal lists `bold, mine, visual`; the list code is unchanged and runs over the new folder slice.
- Catalog skip and its test gone: the `reportIndexWalk` condition, its doc clause, and `TestAIReportIndexSkipsSecondOpinionPresetsFolder` are deleted; the constant moved to `report_second_opinion.go`, its only user.
- Prose: guide :107's commit-first sentence is replaced by the bundle sentence; guide :91 no longer names the folder; `ai-report.md` :73 names `<bundle>/second-opinion-presets/<name>.md`, then the shipped presets. A `git grep` over `skills/`, `_dev/`, `docs/` and `suite/` (changelogs excluded) finds no other mention of the shared folder outside the new test.
- Shipped-only and plain runs unchanged: the shipped case stays in `TestSecondOpinionPresetResolutionAndArgumentRules`; the repository gate below covers the rest.

**RED re-check after the builder's test tidy:** the builder could not re-run RED after its tidy. I ran the final test file against the pre-change production code in a scratch detached checkout at `bac2288b`: `go test -run '^TestSecondOpinionPresetResolvesFromBundleOnly$' ./internal/toolboxcommands` failed at case (a) with `unknown preset "mine"; presets found: bold, visual`, the same failure the builder recorded. The same test passes at `d8cb783f`, and `REQ-705-probe.sh` exits 0 there.

**Scope comparison:** declared 13 paths; touched 6 project paths, all declared; the 7 untouched are the release, lesson and lessons-index files owned by finalization. The files the Scope said it would not touch are untouched.

**Judgment:** the change is substantive, small (40 insertions, 48 deletions) and matches D-01 to D-07 and the hand-back's D-08 and D-09. Qualified.

## Testing

**Tests run:** `DO_WORK_FAST_STAGE_REUSE=off bash _dev/tests/maintainer-verify.sh` at merge `d8cb783f` (run 1, started 2026-10-20T21:37:40Z at 1-minute load 3.99, wall 127 s); then `advance --gate-arg bash --gate-arg _dev/tests/maintainer-verify.sh --gate-exit-status 0 -- --probe-file do-work/runs/work-2026-10-20-200259/REQ-705-probe.sh` (green-gate satisfied; `run-blocked-check` satisfied, probe status 0: bundle-only preset lookup, catalog skip removal and shared-folder sweep checked). After the review fix (F2, one guide sentence, builder-branch commit `e206c552`, re-merged as `5cbf1f10` with the same `<pre>`): the same gate argv at `5cbf1f10` (run 2, started 2026-10-20T21:46:00Z at 1-minute load 2.24, wall 135 s) exit 0, recorded with `record-green-gate`; the probe exits 0 there.
**Result:** ✓ All passing (repository gate exit 0 on the first run: do-work-cli 947 tests, queue-kanban 427 tests, every per-file budget under 30 s, run 2's slowest file `report_second_opinion_test.go` at 27.21 s; builder's package selection `SecondOpinion|AIReportIndex|TestHandlersRegisterCanonicalToolboxCommands` 32 top-level PASS, 0 FAIL, 0 SKIP)

**Red-green validation:**
- `TestSecondOpinionPresetResolvesFromBundleOnly` (`report_second_opinion_test.go`): ✗ before implementation (builder, at base `11780ac0`: case (a) refused with `unknown preset "mine"; presets found: bold, visual`; re-checked by the integrator with the final, tidied test file against the pre-change code at `bac2288b`: same failure) → ✓ after (builder at `31868b40`; integrator at merge `d8cb783f`, PASS 0.18 s). This is the `## Red-Green Proof` pair: RED case = an uncommitted `<bundle>/second-opinion-presets/mine.md`; GREEN cases (1) prepared with the bundle `preset_file`, (2) bundle `bold` wins over shipped `bold`, (3) a committed `ai-reports/second-opinion-presets/legacy.md` is refused with `presets found: bold, mine, visual`. GREEN case (4) (shipped-only and plain runs as in 0.306.6) is covered by the unchanged shipped case and the gate.

**New tests added:**
- `TestSecondOpinionPresetResolvesFromBundleOnly` in `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go`

**Existing tests updated (cross-REQ impact):**
- `report_second_opinion_test.go` (from REQ-703): the project-preset block left `TestSecondOpinionPresetResolutionAndArgumentRules`, because the shared project folder is no longer read (the user's answer).
- `report_index_test.go` (from REQ-703): `TestAIReportIndexSkipsSecondOpinionPresetsFolder` deleted with the skip it pinned (D-04).

**Heavy verification plan:**
- Range: bac2288b6beaf716889f60cad8228cb6adaaa04e..5cbf1f10f0a0ebc7c2d6eeabe385cff6a22abb41 (re-planned after the re-merge; the same four lanes as at `d8cb783f`)
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — the four `toolboxcommands` Go files matched subtree `skills/do-work/tools/do-work-cli`
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — all six changed files matched subtree `skills`
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — the four Go files matched subtree `skills/do-work/tools/do-work-cli`
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — the four Go files matched subtree `skills/do-work/tools/do-work-cli`

*Verified by work action*

## Review

**Overall: 96%** | 2026-10-20T21:49:36Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:**
- F1 `report_index.go:90`: a 0.306.6 consumer's `ai-reports/second-opinion-presets/` now shows in index/find/catalog as an undated bundle; keep the builder's migration bullet in the patch-release CHANGELOG entry — impact-user-visible → report only
- F2 `ai-report-guide.md:107` (behavior `report_second_opinion.go:158`): the guide does not say a revision does not carry a bundle preset, or that a superseded target searches the current revision (D-06), and the refusal does not name the folder searched; add "a revision does not carry it over; copy it into the new bundle" — impact-user-visible → fixed by the integrator in `e206c552` (guide :107 now says a new revision does not carry the preset over and a run on a superseded report uses its current revision), re-merged `5cbf1f10`, gate re-run green; delta re-review of `d8cb783f..5cbf1f10` by the same reviewer: the sentence is true and closes F2 (optional nit, not taken: "the new bundle" could say "the current revision's bundle" for a chain of more than one hop)
- F3 `do-work/user-requests/UR-160/input.md:16,22`: the Summary row and D2 still state the shared project-folder lookup; leave the UR frozen and tell the REQ-704 reviewer that REQ-705's answer supersedes D2 — impact-negligible → report only
- F4 (Nit) `report_second_opinion_test.go:873,877,885,886`: the new test ignores `os.MkdirAll`/`os.WriteFile` errors (same pattern as the removed block) — impact-negligible → report only

**Anti-bloat:** none. No helper, option, file, constant or test that the REQ did not name (one D-05 test added, D-03 constant moved).
**Acceptance:** Pass. Implementation and integration stages: focused `SecondOpinionPreset|AIReportIndex` run passed 6/6 at `d8cb783f`; the integrator's RED at `bac2288b` is recorded; the gate passed at `d8cb783f`.
**Restatement sweep:** redefined: user preset location (bundle-only plus shipped), `preset_file` meaning, catalog walk no longer skips `second-opinion-presets/`. Grepped `skills/`, `docs/`, `_dev/`, `suite/`, READMEs and `do-work/`, with the 0.306.6 CHANGELOG entries and the archived REQ-703 left as history. One stale restatement found: UR-160 input.md:16,22 (F3). The shipped prose (guide :91, :107, `ai-report.md:73`), `result_model.go:715`, the shipped presets, the prompt, `SKILL.md`, the lessons file and REQ-704 agree with the new meaning.
**Suggested testing:** 3 items
**Follow-ups created:** None (F2 fixed; F1, F3, F4 report only. F1 is answered by the CHANGELOG migration bullet; F3 is passed to the coordinator for REQ-704's reviewer)

*Reviewed by review-work action*

## Heavy Verification Plan

- Base revision: bac2288b6beaf716889f60cad8228cb6adaaa04e
- Target revision: 5cbf1f10f0a0ebc7c2d6eeabe385cff6a22abb41
- do-work-cli-integrations: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane do-work-cli-integrations` — `report_index.go`, `report_index_test.go`, `report_second_opinion.go`, `report_second_opinion_test.go` matched subtree `skills/do-work/tools/do-work-cli`
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — all six changed files matched subtree `skills`
- updater: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane updater` — the four Go files matched subtree `skills/do-work/tools/do-work-cli`
- installer: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane installer` — the four Go files matched subtree `skills/do-work/tools/do-work-cli`

## Heavy Verification Result

- Target revision: 5cbf1f10f0a0ebc7c2d6eeabe385cff6a22abb41
- Execution revision: 5cbf1f10f0a0ebc7c2d6eeabe385cff6a22abb41 (detached drain checkout, `QUEUE_KANBAN_BROWSER` set)
- do-work-cli-integrations: executed, exit 0, 68 s
- staged-skills: executed, exit 0, 36 s
- updater: executed, exit 0, 66 s
- installer: executed, exit 0, 31 s

## Lessons Learned

**What worked:** Answering the confirm question by moving the file instead of adding an exemption: a preset inside the bundle is never outside-bundle dirt, so the outside-bundle check stayed whole and the change deleted more lines (48) than it added (40). One lookup entry (`handleAIReportSecondOpinion`) covers the Session, Targeted and Batch forms, so one slice element changed every form.
**What didn't:** The builder could not re-run RED after a test-only tidy because `git stash` and reverting production code were off limits. Copying the final test file into a scratch detached checkout at the pre-merge commit re-ran RED without touching any branch.
**Worth knowing:** The lookup uses the bundle `secondOpinionBundle` resolved, after `latest` and the supersession chain, so a preset in a superseded revision is never seen; the guide now says to copy it into the new bundle. A bundle preset is a pre-existing bundle file, so preservation refuses an editor that edits it. `report_second_opinion_test.go` ran 27.21 s on gate run 2, close to the 30 s per-file budget; the next test added there may need the budget looked at.

## Orientation

Now your own second-opinion preset lives in the report bundle it steers (`<bundle>/second-opinion-presets/<name>.md`) and needs no commit; the shared `ai-reports/second-opinion-presets/` folder is gone. Lives in the do-work-cli toolbox commands (`prime-do-work-cli.md`, ai-report second opinion). Prime spot-check: `prime-do-work-cli.md` and `prime-releases.md` name no path this change removed.

## Timing

Observed 2026-10-20T21:31:30Z to 2026-10-20T21:54:03Z: 22m 33s total, 21m 52s attributed across 6 events, 41s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| verification-gate | 9m 44s | 3 |
| review | 7m 36s | 1 |
| builder-work | 4m 02s | 1 |
| handback-merge | 30s | 1 |

Slowest stage: review / independent review, 7m 36s, outcome success.
Slowest command: verification-gate / heavy drain, 3m 47s, exit 0, .
