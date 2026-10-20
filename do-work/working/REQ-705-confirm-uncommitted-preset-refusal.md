---
id: REQ-705
title: 'Confirm: a second opinion refuses to run while a project preset file is uncommitted'
status: claimed
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
