# Hand-back: REQ-701 (a failed second-opinion model check keeps the finished page; a subagent model never makes the backend model ambiguous)

- Branch: `worktree-agent-REQ-701-keep-finished-page-unverified-model`
- Worktree: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-701-keep-finished-page-unverified-model
- Base commit (before first edit): `cc3e72c2`
- Commits: `86617f8e` `[REQ-701] keep the finished second-opinion page when the model check fails` (one commit; tree clean after it)
- Boundary: nothing under `do-work/` touched in the worktree; this file is the only main-tree write, unstaged.

## File manifest

- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go` (modified): sentinel `errSecondOpinionModelUnverified`; `secondOpinionMetadata` gains the `expected` parameter and `seen` return and applies requested-model-wins in the shared gemini/claude branch; `secondOpinionRun` keeps a non-empty staging file on any refusal (removes it only when published or empty) and publishes an unverified page and reply pair when the model check fails after the page is written; handler row literal initializes `SeenModels: []string{}`.
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go` (modified): helper `secondOpinionFakeClaude`; new `TestSecondOpinionSubagentModelIsNotAmbiguous` and `TestSecondOpinionUnverifiedModelKeepsFinishedPage`; `TestSecondOpinionFreshGenerationModelMismatchRefusesPublication` renamed to `...PublishesUnverified` with flipped assertions; one kept-staging assertion inside the missing-backlink case; the six-case identity test passes `""` as `expected` and takes the fourth return; `fmt` import added for `fmt.Sprint` on JSON arrays.
- `skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go` (modified): `SeenModels []string \`json:"seen_models"\`` after `ActualModel` on `AIReportSecondOpinionRun`.
- `skills/do-work-toolbox/actions/ai-report.md` (modified): one sentence after `existing files are never overwritten.` describing the unverified name, state and `seen_models`.
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modified): one sentence in § A Second Opinion after the numeric-suffix sentence.

## P-A-U

- [x] **[PLAN]:** Read general, coding-guardrails, shared-principles, communication-style, testing; primes prime-shell-commands and prime-releases; the `native-backend-identity` lesson bullet (REQ-698). Approach exactly as the brief: tests first against a signature-only scaffold (D-02), then the defer rule, the unverified branch, the requested-model-wins rule, the row field, two doc sentences.
- [x] **[APPLY]:** Code written as planned, confined to the five files. No prepare/finalize work, no flags, no Session-form rewrite (REQ-700 territory untouched).
- [x] **[UNIFY]:** `git diff cc3e72c2 --stat`:
  ```
   skills/do-work-toolbox/actions/ai-report.md        |   2 +-
   skills/do-work-toolbox/docs/ai-report-guide.md     |   2 +-
   .../internal/resultmodel/result_model.go           |   1 +
   .../toolboxcommands/report_second_opinion.go       | 103 +++++++++++++------
   .../toolboxcommands/report_second_opinion_test.go  | 114 ++++++++++++++++++++-
   5 files changed, 182 insertions(+), 40 deletions(-)
  ```
  Checks (all from the worktree root):
  | Check | Exit | Wall |
  |---|---|---|
  | RED three-test run (scaffold) | 1 (expected) | 2 s |
  | GREEN family `-run '^TestSecondOpinion'` | 0, 19 pass, 0 skip | 13 s |
  | `go test ./internal/resultmodel/ ./internal/toolboxcommands/` | 0 | 14 s |
  | `go vet ./internal/toolboxcommands/ ./internal/resultmodel/` | 0 | <1 s |
  | `bash REQ-701-probe.sh` | 0 | 2 s |
  | `gofmt -l skills/do-work/tools/do-work-cli/internal` | prints nothing | |
  | `git diff --check` | 0 | |
  Files reviewed: the full source diff of `report_second_opinion.go` and `result_model.go` (read line by line; no debug prints, comments updated), the test diff, both doc sentences in place (grep confirmed one match each). No run needed a second attempt under load.

## Proof record (tdd: true)

RED, run against the compile-only scaffold (new signature and `seen_models` field, old `len(models) != 1` rule and old `Remove(stageName)` defer intact):
```
--- FAIL: TestSecondOpinionFreshGenerationModelMismatchPublishesUnverified (0.40s)
    model switch not published as unverified: [... actual_model:different-model ... page: reason:generation ran model different-model instead of actual-model-1; substitution refused ... seen_models:[] state:failed]
--- FAIL: TestSecondOpinionSubagentModelIsNotAmbiguous (0.34s)
    subagent model treated as ambiguity: map[actual_model: ... page: reason:claude actual model is absent or ambiguous (2 models) ... state:failed]
--- FAIL: TestSecondOpinionUnverifiedModelKeepsFinishedPage (0.40s)
    unverified model did not keep the finished page: map[actual_model:claude-fable-5-1 ... page: reason:claude actual model is absent or ambiguous (2 models) ... state:failed]
FAIL	github.com/knews2019/skill-do-work/do-work-cli/internal/toolboxcommands	1.730s
```
GREEN: `-run '^TestSecondOpinion'` 19 of 19 `--- PASS`, no skips, `ok ... 11.964s`; the two-package run `ok resultmodel 0.270s`, `ok toolboxcommands 14.093s`; probe script exit 0.

## Decisions

- D-01 (pre-dispatch): the requested-model-wins rule lives in the shared gemini/claude branch of `secondOpinionMetadata`, not behind a `claude` condition. Value: one code path; gemini's `stats.models` has the same shape. Risk: none observed; the gemini empty-map case still errors.
- D-02 DECIDE & STATE: the RED run used a signature-only scaffold (new `secondOpinionMetadata` shape, `SeenModels` field, no behavior change) so the three tests failed on their assertions, not on compilation. A compile failure is not RED-before-GREEN under testing.md. The failure texts match the brief's expectation word for word.
- D-03 DECIDE & STATE: the "staged page kept at" path is repository-relative (`ai-reports/<bundle>/.codex-...-staging-X.html`), the same form as `row.PagePath`, which is why `relative` moved up as the brief said. If the staging file is already gone when the defer runs (stat fails), nothing is removed and nothing is appended, so the reason never names a file that does not exist.
- D-04 DECIDE & STATE: `seen` for codex and agy is derived once after the switch (`if seen == nil && model != ""`) rather than inside each case; same result, one line. The final slug-pattern failure still returns `nil` seen as a non-identity error, per brief outcome (a).
- D-05 DECIDE & STATE: `row.SeenModels = append([]string{}, seen...)` so a `nil` seen from a non-identity error keeps the JSON at `[]`, never `null`.

## Discovered Tasks

- impact-low: `secondOpinionRun` sets `row.ActualModel = resolved` from the probe before generation, so a non-identity generation failure (outcome a) reports the probe's model as `actual_model` while `state` is `failed`. Pre-existing, unchanged here → report only.
- impact-low: a kept staging file (`.<cli>-<model>-staging-*.html`) joins the `before` snapshot of the next run, so later runs tolerate it, but nothing ever cleans it up; the user deletes it by hand after reading it → report only.

## Lessons read

- Families: `native-backend-identity` bullet (REQ-698, lessons-do-work-cli.md line 5). Honored: the unverified page is never named after a model the backend did not confirm; only the bytes survive.
- Satellites: none beyond that bullet (the full `lessons-do-work-cli.md` was dropped for budget and no touched code is named by its other bullets). Primes: prime-shell-commands (fake is POSIX `sh`, single-quoted `printf '%s\n'` literals, no bashisms), prime-releases (release is the integrator's).

## Anti-bloat check

`git diff --stat` above: five files, 182 insertions, 40 deletions. Added beyond the brief's named items: none. Named items present: `errSecondOpinionModelUnverified`, `SeenModels`/`seen_models`, `expected` parameter and `seen` return, `stagePublished`, `secondOpinionFakeClaude`, two new tests, the renamed test. Two unnamed locals inside `secondOpinionRun` (`unverified` bool, `publishedModel` string) and the `fmt` test import are short-lived locals, not helpers or surface.

## Proposed CHANGELOG entry

**Second-Opinion Page Survives a Failed Model Check**

A finished second-opinion page is never deleted because the backend's model could not be confirmed. The page is kept under a name that says so, and a subagent's model in Claude's usage map no longer reads as ambiguity.

- `ai-report second-opinion` keeps the staged page on every refusal after the backend wrote it; the failure reason names the kept file. Only an empty staging file, or one already published, is removed.
- When the model cannot be verified after the page is written, the page and reply are published as `index-<cli>-<model>-unverified[-N].html` and `<cli>-<model>-unverified[-N]-reply.md`, run state `unverified`, no banner, no judge, with typed finding `AI-REPORT-SECOND-OPINION-UNVERIFIED` listing the requested model and every model seen.
- For `claude` and `gemini`, the requested model wins when the native usage map lists it; extra entries land in the new `seen_models` row field.

## Proposed lesson bullet

- [family: native-backend-identity] [REQ-701: an identity refusal stops the NAME, never the bytes: a page the backend finished is published under an explicit `-unverified` name with no banner, and a subagent's model beside the requested one is evidence to record, not a second candidate.](<relative archive link>)

## Integration seams

- REQ-700 (prepare/finalize split) branches from main after this lands. The unverified publication lives in `secondOpinionRun` between the generation metadata call and the judge (the `unverified` branch: name from `secondOpinionNames(bundle, backend, resolved+"-unverified")`, `row.RunState = "unverified"`, early return before `secondOpinionJudge`). The kept-staging rule is the deferred closure right after `os.CreateTemp` keyed on `stagePublished`. `seen_models` is set once from the generation envelope by `row.SeenModels = append([]string{}, seen...)` in `secondOpinionRun`; the requested-model-wins rule is in the shared `gemini`/`claude` branch of `secondOpinionMetadata`. The finalize phase must keep all three.
- Test wall times: family 13 s, two packages 14 s, probe 2 s, RED 2 s.
