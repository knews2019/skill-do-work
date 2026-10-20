# Review: REQ-705

**Approve**. A user preset now comes only from `<bundle>/second-opinion-presets/<name>.md`, then the shipped presets. The shared `ai-reports/second-opinion-presets/` folder and the catalog skip are gone. Nothing blocks. Four noncritical findings, all report only.
Route B | merge `d8cb783f` (range `bac2288b..d8cb783f`, 6 files, +40 / -48)

REQ-705 here means the user's answer in `## Answer Notes`: user presets live only inside the report bundle. The REQ title, What, What the Builder Chose and What Would Change describe the earlier confirm question and were read as history only.

### What's built
- `--preset <name>` looks in the resolved bundle's `second-opinion-presets/` folder first, then beside the shipped prompt. This covers the Session prepare, Targeted and Batch forms, because they share one lookup in `handleAIReportSecondOpinion`. `--finalize` reads the preset from the prepare record.
- An uncommitted bundle preset is never outside-bundle dirt. The prepare path runs `secondOpinionOutsideCheck` and reaches `state: prepared`.
- The catalog walk lists every non-dot directory under `ai-reports/`. The skip constant moved to `report_second_opinion.go`, its only user.
- Not built, and correctly so: release files and the CHANGELOG entry. The integrator writes them at finalization.

### Decisions / risks for you
- D-06 (a preset in a superseded bundle is not seen) is a user trap the guide does not mention. See F2. Value of the current choice: one folder, no walk along the revision chain. Risk: after `ai-report revise`, or when a superseded bundle is targeted and `secondOpinionBundle` redirects to the current revision, `--preset mine` is refused even though the user's file exists.

### Findings

**Important:**
- None.

**Minor:**
- F1. `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index.go:90`: a consumer who made `ai-reports/second-opinion-presets/` under 0.306.6 now sees it in `index`, `find` and `catalog.json` as an undated bundle whose entry is its first `.md` file. `latest` is not affected, because undated bundles sort last. This is the intended effect of the user's answer and D-04. The only gap is the missing release note, and the CHANGELOG is not written yet. Fix: keep the builder's migration bullet ("move each preset into the bundle you use it with and delete the folder") in the patch-release entry. — impact-user-visible → report only
- F2. `skills/do-work-toolbox/docs/ai-report-guide.md:107` (behavior at `report_second_opinion.go:158`): the guide says a user preset lives in the bundle. It does not say that a revision does not carry the preset over, or that targeting a superseded bundle searches the current revision's folder. The `unknown preset` refusal lists names but not the folder it searched. So the user sees a refusal with no hint of the cause. Fix: add one clause to :107: "a revision does not carry it over; copy it into the new bundle". — impact-user-visible → report only
- F3. Restatement in `do-work/user-requests/UR-160/input.md:16` and `:22`. The Summary table ("shipped + project preset files") and D2 ("a project file `ai-reports/second-opinion-presets/<name>.md` with the same name wins") still state the old lookup. review-work.md Step 3 makes the UR the source of truth, so REQ-704's reviewer could score the missing project lookup against REQ-704. The verbatim block at :54 is a frozen record and is correct as history. Fix: do not edit the UR. The REQ-704 reviewer prompt should say that REQ-705's answer supersedes UR-160 D2 on preset location. — impact-negligible → report only

**Nit:**
- F4. `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go:873,877,885,886`: the new test ignores the errors from `os.MkdirAll` and `os.WriteFile`. A failed fixture write would show up as a confusing `unknown preset` failure. The removed project block used the same pattern, so this is not a new style. — impact-negligible → report only

### Requirements Checklist

- [x] `--preset <name>` resolves `<bundle>/second-opinion-presets/<name>.md` first, then the shipped preset, and `preset_file` names the bundle file repository-relative. Delivered: the `folders` slice at `report_second_opinion.go:158` and the existing `filepath.Rel` step at :188. Pinned by test cases `mine` and `bold`.
- [x] An uncommitted bundle preset runs with no outside-bundle refusal. Delivered: `secondOpinionOutside` counts only paths outside the bundle prefix, and prepare runs that check (:535). The test leaves both presets uncommitted and gets `prepared`.
- [x] `ai-reports/second-opinion-presets/` is no longer read. Delivered: the committed `legacy.md` case is refused.
- [x] The unknown-name list is the bundle presets plus the shipped presets. Delivered: `presets found: bold, mine, visual`.
- [x] The `reportIndexWalk` skip and its test are gone. Delivered: `report_index.go:89-92` and the deleted `TestAIReportIndexSkipsSecondOpinionPresetsFolder`.
- [x] The commit-first sentence is gone, and the guide and `ai-report.md` say presets live in the bundle or ship with the toolbox. Delivered: guide :91 and :107, `ai-report.md:73`.
- [x] Shipped-only presets, plain runs and existing tests behave as in 0.306.6. Delivered: the shipped case stays in `TestSecondOpinionPresetResolutionAndArgumentRules`, `TestSecondOpinionDirectedRunNamesAndMarksPage` is unchanged, and the gate passed at `d8cb783f`.
- [x] UR-160 D3 (the direction never overrides the file BOUNDARY). Still holds: `secondOpinionPreserved` refuses any change to a pre-existing bundle path, so an editor cannot rewrite a bundle preset. It also refuses any new non-root file, so a backend cannot plant a `second-opinion-presets/` file.
- [x] Implementation Summary claims match the diff. All 6 files and every claimed change are present.

### Acceptance Testing

**Result: Pass** (stages covered: implementation and integration)
- At `d8cb783f` (clean tree apart from run artifacts), `go test -C skills/do-work/tools/do-work-cli -count=1 -run 'SecondOpinionPreset|AIReportIndex' ./internal/toolboxcommands -v` passed: 6 top-level tests, 0 FAIL, 0.65 s. This includes `TestSecondOpinionPresetResolvesFromBundleOnly` and all four remaining `TestAIReportIndex*` tests.
- The new test drives the real handler end to end on the session prepare form. It uses a real Git fixture, an uncommitted bundle preset, a bundle preset that shadows a shipped one, and a committed legacy preset that is refused.
- Red-green: the integrator reran the final test file against the pre-change code at `bac2288b`. It failed with `unknown preset "mine"; presets found: bold, visual` (recorded in `## Qualification`). It passes at `d8cb783f`. I did not repeat the RED run, because that needs a scratch checkout and this review is read-only.
- Code read for regressions: `latest` resolution (`secondOpinionBundle` uses `catalog[0]`, and undated bundles sort last), `pickReportEntry` (reads only the bundle root), `architecture.go:62` (matches architecture folder names only), and preservation (`secondOpinionSnapshot` walks subfolders, so it covers a preset folder). None is affected.
- The repository gate passed at `d8cb783f` (the integrator's run, not rerun here).

### Suggested Additional Testing

- Deployment: unassessed. After the release, install into a consumer project and confirm that the shipped `bold.md` and `visual.md` still resolve from the installed toolbox path.
- Live acceptance: unassessed. In a consumer repo, run `ai-report second-opinion --backend codex --preset mine latest` with an uncommitted `<bundle>/second-opinion-presets/mine.md`, and confirm a published page with no outside-bundle refusal. This exercises the Targeted form, which the tests reach only through shared code.
- Edge case: target a superseded bundle that holds a preset, and confirm the refusal is understandable (F2).

### Domain Review

Backend checklist: invalid input is handled. The preset name must match `^[a-z0-9][a-z0-9-]*$`, so no path traversal into or out of the bundle is possible, and a non-regular preset gets a usage error. The refusal shows only preset names and no internal paths. No auth surface. Prepare stays idempotent per run: a failed prepare removes its reserved files. Existing tests pass.

### Anti-bloat

Additions the REQ did not name: none. No new helper, option, flag, file, environment variable, or decorative test. The one new test is the one D-05 named. The constant moved, as D-03 named, and no new one was added. The net diff is 8 lines smaller.

### Scores (on the record, not the headline)

**Overall: 96%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | 7/7 acceptance criteria delivered |
| Code Quality | 95% | Minimal and follows the file's patterns. Unchecked fixture writes in the test (F4) |
| Test Adequacy | 90% | Real red-green on the named regression. The uncommitted case runs only on the prepare form |
| Scope | 100% | 6 of 13 declared paths touched. The 7 untouched are release and lesson files for finalization. Nothing undeclared |
| Risk | Low | Pre-existing `ai-reports/second-opinion-presets/` folders appear in the catalog (F1) |
| Acceptance | Pass | Implementation and integration stages |

### Follow-ups created
None (4 findings report only)

## Review

**Overall: 96%** | <timestamp>

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
- F2 `ai-report-guide.md:107` (behavior `report_second_opinion.go:158`): the guide does not say a revision does not carry a bundle preset, or that a superseded target searches the current revision (D-06), and the refusal does not name the folder searched; add "a revision does not carry it over; copy it into the new bundle" — impact-user-visible → report only
- F3 `do-work/user-requests/UR-160/input.md:16,22`: the Summary row and D2 still state the shared project-folder lookup; leave the UR frozen and tell the REQ-704 reviewer that REQ-705's answer supersedes D2 — impact-negligible → report only
- F4 (Nit) `report_second_opinion_test.go:873,877,885,886`: the new test ignores `os.MkdirAll`/`os.WriteFile` errors (same pattern as the removed block) — impact-negligible → report only

**Anti-bloat:** none. No helper, option, file, constant or test that the REQ did not name (one D-05 test added, D-03 constant moved).
**Acceptance:** Pass. Implementation and integration stages: focused `SecondOpinionPreset|AIReportIndex` run passed 6/6 at `d8cb783f`; the integrator's RED at `bac2288b` is recorded; the gate passed at `d8cb783f`.
**Restatement sweep:** redefined: user preset location (bundle-only plus shipped), `preset_file` meaning, catalog walk no longer skips `second-opinion-presets/`. Grepped `skills/`, `docs/`, `_dev/`, `suite/`, READMEs and `do-work/`, with the 0.306.6 CHANGELOG entries and the archived REQ-703 left as history. One stale restatement found: UR-160 input.md:16,22 (F3). The shipped prose (guide :91, :107, `ai-report.md:73`), `result_model.go:715`, the shipped presets, the prompt, `SKILL.md`, the lessons file and REQ-704 agree with the new meaning.
**Suggested testing:** 3 items
**Follow-ups created:** None (4 findings report only)

*Reviewed by review-work action*
