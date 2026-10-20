# REQ-705 hand-back (user second-opinion presets live only inside the report bundle)

- Branch: `worktree-agent-REQ-705-bundle-local-second-opinion-presets`
- Worktree: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-705-bundle-local-second-opinion-presets
- Base commit: `11780ac0` (checked with `git rev-parse --short HEAD` before the first edit)
- Commits: `31868b40` `[REQ-705] move user second-opinion presets into the report bundle` (one commit, six files, committed by exact path)

## File manifest

- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` (modified): the first preset folder is now `filepath.Join(bundle, secondOpinionPresetsFolder)`; the comment says a bundle preset wins over the shipped one; the constant `secondOpinionPresetsFolder` now lives here, after `secondOpinionUsage`, with the comment "the folder name for presets, both in a report bundle and beside the shipped prompt".
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index.go` (modified): the constant and its comment left the const block; the `reportIndexWalk` skip `&& entry.Name() != secondOpinionPresetsFolder` and the doc-comment clause are gone.
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index_test.go` (modified): `TestAIReportIndexSkipsSecondOpinionPresetsFolder` and its comment deleted. No import became unused (`reflect`, `os`, `filepath` still used elsewhere).
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go` (modified): project-preset block removed from `TestSecondOpinionPresetResolutionAndArgumentRules` and from its `// Pins:` comment; new `TestSecondOpinionPresetResolvesFromBundleOnly` right after it.
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modified): :91 now says only loose files are skipped; :107 last sentence is now "Your own preset goes inside the report bundle as `second-opinion-presets/<name>.md`, wins over a shipped preset of the same name, and needs no commit."
- `skills/do-work-toolbox/actions/ai-report.md` (modified): :73 `preset_file` clause reads "(the file used: `<bundle>/second-opinion-presets/<name>.md`, then the shipped presets)".

## [PLAN]

Follow brief steps 1-6 exactly. Write the new test first and prove RED, then change the lookup slice, remove the catalog skip, move the constant, delete the skip test and the project block, then the two prose sentences.

## [APPLY]

1. Wrote `TestSecondOpinionPresetResolvesFromBundleOnly` (cases a, b, c from the brief; ordered slice so case (a) `mine` always runs first). RED run, see Proof.
2. Changed the `folders` slice and its comment; moved the constant; removed the skip and its doc clause; deleted the skip test; removed the project block and its pin phrase.
3. GREEN run and package selection, see Proof.
4. Prose edits in the guide (:91, :107) and `ai-report.md` (:73).
5. After GREEN, one test-only tidy: the new test writes each bundle preset inside its loop instead of writing both up front with the text repeated. Same files, same order (`mine` is written and checked before `bold` exists), same assertions. The selection was rerun after this tidy (below). The RED run was not repeated after the tidy, because no `git stash` and no revert of production code was allowed; case (a) is unchanged, so its RED failure is the same.

## [UNIFY]

`git diff 11780ac0 --stat`:

```
 skills/do-work-toolbox/actions/ai-report.md        |  2 +-
 skills/do-work-toolbox/docs/ai-report-guide.md     |  4 +--
 .../internal/toolboxcommands/report_index.go       |  8 ++---
 .../internal/toolboxcommands/report_index_test.go  | 25 -------------
 .../toolboxcommands/report_second_opinion.go       |  7 ++--
 .../toolboxcommands/report_second_opinion_test.go  | 42 +++++++++++++++-------
 6 files changed, 40 insertions(+), 48 deletions(-)
```

| Check (from worktree root) | Exit | Wall time |
|---|---|---|
| `bash .../REQ-705-probe.sh` | 0 ("bundle-only preset lookup, catalog skip removal and shared-folder sweep checked") | 1.6 s |
| `gofmt -l skills/do-work/tools/do-work-cli/internal` | 0, prints nothing | under 1 s |
| `go vet -C skills/do-work/tools/do-work-cli ./internal/toolboxcommands` | 0 | a few s |
| `git diff --check` | 0 | under 1 s |
| `git diff 11780ac0 --stat` | exactly the six files | under 1 s |

Files checked by reading the diff: all six above. Probe rerun after the test tidy: exit 0. Working tree clean after the commit.

## Proof record

RED (before any production edit), `go test -C skills/do-work/tools/do-work-cli -count=1 -v -run '^TestSecondOpinionPresetResolvesFromBundleOnly$' ./internal/toolboxcommands`:

```
report_second_opinion_test.go:889: missing typed backend summary: {... "observed_evidence":["unknown preset \"mine\"; presets found: bold, visual"] ...}
--- FAIL: TestSecondOpinionPresetResolvesFromBundleOnly (0.07s)
```

It fails at case (a), with the expected text `unknown preset "mine"; presets found: bold, visual` (wall 0.95 s).

GREEN, same command: `--- PASS: TestSecondOpinionPresetResolvesFromBundleOnly (0.12s)`, `ok` (wall 1.43 s), no skip.

Package selection:
- `-run 'SecondOpinion|AIReportIndex|TestHandlersRegisterCanonicalToolboxCommands' ./internal/toolboxcommands`: exit 0. Top-level PASS 32, FAIL 0, SKIP 0 (22 subtests also PASS). Wall 15.3 s on the first run; 13.8 s on the rerun after the test tidy, same counts.
- `-run '^TestSecondOpinionJSONKeepsEveryBackendAndJudgeEvidence$' ./internal/resultmodel`: PASS 1, FAIL 0, SKIP 0 (wall 0.49 s).

No budget failure under load, so no rerun was needed for timing.

## Decisions

- **D-08** (DECIDE & STATE): the new test loops over an ordered `[][2]string` (`mine`, then `bold`), not a map, so case (a) always runs first and the RED failure is deterministic.
- **D-09** (DECIDE & STATE): the legacy case asserts the refusal through `result.Findings[0].Evidence`, the same shape as the refusal loop in `TestSecondOpinionPresetResolutionAndArgumentRules`. No new helper.

## Discovered Tasks

- A consumer that already made `ai-reports/second-opinion-presets/` under 0.306.6 now sees it in the catalog, `index` and `find` as an undated bundle whose entry is its first `.md` file. `latest` is not affected, because undated bundles sort last. This is the intended result of D-04 (the skip is gone), but the CHANGELOG entry should tell those users to move their presets into a bundle and delete the old folder. → report only

## Lessons read

- `_dev/primes/lessons-releases.md` (whole file): family `canonical-link-outlives-its-target` (matters at UR-160 close for the REQ-703 lesson link, not for this change) and `manifest-ownership-vs-edit-content`.
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` § Traps, family `guard-on-one-entry-path`: the lookup has one entry in `handleAIReportSecondOpinion`, shared by session prepare, targeted and batch forms, and `--finalize` never looks presets up; so one change covers every form.
- Crew members: general, coding-guardrails, shared-principles, communication-style, testing, backend.

## Anti-bloat check

Diff stat above: 40 insertions, 48 deletions, six files.

Added that the brief did not name: none. No new helper, flag, option, constant (the constant moved, it was not added), file, environment variable or test besides `TestSecondOpinionPresetResolvesFromBundleOnly`.

## Proposed CHANGELOG entry

### Your Second-Opinion Presets Live in the Report Bundle

Your own second-opinion preset now sits inside the report it steers, so it needs no commit before a run and nothing outside the bundle is read or refused.

- `--preset <name>` looks in `<bundle>/second-opinion-presets/<name>.md` first, then in the presets that ship with the toolbox. A bundle preset wins over a shipped preset of the same name, and `preset_file` names the bundle file.
- The shared folder `ai-reports/second-opinion-presets/` is no longer read. A name found only there is an unknown preset, and the refusal lists the bundle and shipped presets it did find.
- The report catalog no longer skips `ai-reports/second-opinion-presets/`. If you made that folder under 0.306.6, move each preset into the bundle you use it with and delete the folder.
- A preset in a superseded revision is not seen: after `ai-report revise`, copy the preset into the new bundle.

## Proposed lesson bullet

None. The change followed the recorded decisions without surprises, and the one trap (a single lookup entry shared by every form) is already the `guard-on-one-entry-path` family.

## Integration seams for REQ-704

- The lookup reads `filepath.Join(bundle, secondOpinionPresetsFolder)` inside `handleAIReportSecondOpinion`, after `bundle` is resolved by `secondOpinionBundle` (absolute, after `latest` and the supersession chain). A page-path target must resolve its bundle into `bundle` before the `if options.preset != ""` block, or the lookup reads the wrong folder.
- `secondOpinionNames`, `secondOpinionDecorate` and `secondOpinionBundle` are unchanged. The constant `secondOpinionPresetsFolder` now lives in `report_second_opinion.go` right after `secondOpinionUsage`.
- A bundle preset folder is a subfolder of the bundle, so preservation (`secondOpinionSnapshot`) walks it and refuses an editor that changes a preset. If REQ-704's tree listing on `index.html` walks bundle subfolders, it should not list `second-opinion-presets/` as an iteration.

Test wall times: RED 0.95 s, GREEN 1.43 s, toolboxcommands selection 15.3 s and 13.8 s, resultmodel 0.49 s, probe 1.6 s.
