# Builder brief: REQ-705 (your own second-opinion presets live only inside the report bundle)

- Worktree (your only writable tree): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-705-bundle-local-second-opinion-presets
- Branch: worktree-agent-REQ-705-bundle-local-second-opinion-presets, created by the coordinator from main HEAD with `git worktree add -b` right before dispatch (do not create it; commit here; never merge, never rebase, never push, never check out another branch).
- REQ (read-only, main tree): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/working/REQ-705-confirm-uncommitted-preset-refusal.md. Read it fully. Its title, What, What the Builder Chose and What Would Change describe an old confirm question and are history only. The work is the user's answer in `## Answer Notes` (2026-10-20), restated in `## Triage`, with the decisions D-01 to D-07 in `## Exploration` and the file list in `## Scope`. The release belongs to the integrator; you PROPOSE the changelog entry and a lesson bullet in the hand-back.
- Original REQ (read-only): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/archive/REQ-703-second-opinion-preset-and-editor-direction.md (added presets in 0.306.6) and its hand-back /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-10-20-200259/REQ-703-handback.md.
- UR (read-only): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/user-requests/UR-160/input.md (decision D2 there is what this REQ changes).
- Hand-back file (the ONE main-tree path you may write; never stage or commit it): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-10-20-200259/REQ-705-handback.md
- Route B, tdd: true, impact-user-visible, effort-substantive, domain backend. Small change: one slice element in the Go lookup, one removed skip, one moved constant, one new test, one deleted test, three prose sentences.

## Rules to load first
/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/general.md, coding-guardrails.md, shared-principles.md, communication-style.md, testing.md, backend.md. Primes: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/tools/do-work-cli/prime-do-work-cli.md (§ Traps, especially `guard-on-one-entry-path`), /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/_dev/primes/prime-releases.md (read; the release is the integrator's). Lessons: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/_dev/primes/lessons-releases.md (whole file, small). The do-work-cli satellite was dropped for budget; you need none of its bullets for this change.

## The change (decided, do not reopen)
Paths are relative to your worktree root. Line numbers are at main 11780ac0.
1. `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go:153-188`, in `handleAIReportSecondOpinion`: the first element of `folders` (:155) becomes `filepath.Join(bundle, secondOpinionPresetsFolder)` instead of `filepath.Join(execution.RepositoryRoot, reportsFolder, secondOpinionPresetsFolder)`. `bundle` is already the resolved absolute bundle (:144). Update the comment at :154: a preset in the report bundle wins over the shipped one of the same name, which sits beside the trusted prompt. Everything else in the block stays: the Lstat regular-file check, the unknown-name list (now bundle + shipped, D-02), the trim and `{{` refusal, the repository-relative `preset_file`. No change to `--finalize` (it reads the direction from the prepare record), to the outside-bundle check, or to preservation (a bundle preset is a pre-existing bundle file, so `secondOpinionPreserved` already refuses an editor that changes it).
2. `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index.go`: delete `&& entry.Name() != secondOpinionPresetsFolder` (:94) and the clause ", except the second-opinion presets folder" from the doc comment (:83, keep the sentence grammatical). Move the constant `secondOpinionPresetsFolder` and its comment out of the const block (:28-30) into `report_second_opinion.go` (its only user then, D-03), next to the other second-opinion constants near `secondOpinionUsage` (:29). New comment: the folder name for presets, both in a report bundle and beside the shipped prompt.
3. `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index_test.go:251-276`: delete `TestAIReportIndexSkipsSecondOpinionPresetsFolder` and its comment (D-04). Remove any import it alone used.
4. `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion_test.go`:
   - In `TestSecondOpinionPresetResolutionAndArgumentRules` (:842), delete the project-preset block (:853-862, the `project :=` lines through the "project preset did not win" check) and drop "a project preset that does not shadow the shipped one" from its `// Pins:` comment. The shipped case and the refusal map stay unchanged (its unknown list is still `bold, visual`).
   - Add `TestSecondOpinionPresetResolvesFromBundleOnly` right after it, with a `// Pins:` comment naming the two real failures: a preset kept in the report bundle that is not found, and a leftover lookup that still reads the shared `ai-reports/second-opinion-presets/` folder outside the bundle. Body: `repo, bundle := secondOpinionFixture(t)`; `secondOpinionFixturePresets(t)`; write `<bundle>/second-opinion-presets/mine.md` ("Bundle mine: tell it as a story.\n") and `<bundle>/second-opinion-presets/bold.md` ("Bundle bold: go further.\n") and do NOT commit them. Then, with the same `session` args as the test above: (a) `--preset mine` gives state `prepared`, `preset_file` `ai-reports/2026-10-15_1200-example/second-opinion-presets/mine.md`, and the prompt contains `Direction: Bundle mine: tell it as a story.\n`; (b) `--preset bold` gives the bundle `bold.md` as `preset_file` and its text in the prompt; (c) write `ai-reports/second-opinion-presets/legacy.md`, `git add` that one folder only and commit it, then `--preset legacy` is refused with evidence containing `unknown preset "legacy"; presets found: bold, mine, visual`. Use the existing helpers (`secondOpinionInvoke`, `secondOpinionRows`, `toolboxTestGit`); add no new helper.
5. `skills/do-work-toolbox/docs/ai-report-guide.md`: at :107 replace the last sentence ("A project file `ai-reports/second-opinion-presets/<name>.md` wins ... outside the bundle.") with one sentence: your own preset goes inside the report bundle as `second-opinion-presets/<name>.md`, wins over a shipped preset of the same name, and needs no commit. Keep the literal `second-opinion-presets/<name>.md` (the probe checks it). At :91 drop " and the `second-opinion-presets/` folder" so loose files alone are named as skipped.
6. `skills/do-work-toolbox/actions/ai-report.md:73`: where it says `preset_file` (the file used), add where a preset comes from: `<bundle>/second-opinion-presets/<name>.md`, then the shipped presets. One clause, no new paragraph.

## Anti-bloat (YAGNI); the maintainer asked for this to be watched
- Smallest change that does what the user answered. Prefer deleting or reusing existing code over adding.
- No new helpers, flags, options, config, abstractions, files or environment variables.
- One new test, pinning the named failures; no decorative or "while we're here" tests. No new preservation test: preservation already covers every pre-existing bundle path.
- No exemption or special case in the outside-bundle check. No supersession-chain walk for presets (D-06).
- Do not fix adjacent things you notice; record them as discovered tasks.
In the hand-back, paste `git diff --stat` and list every function, constant, option, file or test you added that this brief did not name (expected: none). Each one gets a one-line reason, or you remove it.

## Write boundary
Exactly the six files in "The change". Anything else: stop and say so in the hand-back. Not touched: `internal/resultmodel/result_model.go`, the shipped prompt, `docs/second-opinion-presets/bold.md` and `visual.md`, `skills/do-work-toolbox/SKILL.md`, `skills/do-work/tools/do-work-cli.sh`, every CHANGELOG.

## Hard rules
- Every commit subject on your branch starts with `[REQ-705]` (for example `[REQ-705] move user presets into the report bundle`). The board credits commits by that prefix only.
- Never create, edit, stage or commit any path under `do-work/`, in the worktree or anywhere else, except the one hand-back file above (which you write but never stage or commit). The integrator's queue guard refuses a branch that touches `do-work/`.
- Never merge, rebase, push, or check out another branch. Never run `recover`, `advance`, any `do-work-cli` lifecycle command, `just do-work-update`, or the repository gate (`_dev/tests/maintainer-verify.sh`); the integrator owns the gate.
- Do not touch `CHANGELOG.md`, `skills/do-work/CHANGELOG.md`, `VERSION`, `skills/do-work/VERSION`, `skills/do-work/actions/version.md` or any version mirror: the release belongs to the integrator.
- Scratch files go in `mktemp -d` directories outside both trees and are never committed. Clean up anything you start.
- Go: run focused `go test -run` and `go vet` on the touched package only, never a whole-module run and never the repository gate.
- Load: one builder runs at a time in this run (REQ-705 now, REQ-704 after REQ-705 is released, because both edit `report_second_opinion.go`). If a wall-time budget fails once under load, rerun it once and record both runs.

## Integration seam
No builder runs beside you. REQ-704 (linked iteration tree: `v1cNN` names, tree list on `index.html`, parent link) branches from main only after REQ-705 is integrated and released, and its brief says presets resolve from the bundle. Keep the lookup inside `handleAIReportSecondOpinion` and keep its use of the resolved `bundle`, so REQ-704 can resolve a page-path target's bundle before the lookup. Do not rename `secondOpinionNames`, `secondOpinionDecorate` or `secondOpinionBundle`. The run artifact `REQ-703-probe.sh` names the test you delete; it is history, nobody reruns it.

## Proof to run and record (from the REQ's Red-Green Proof)
1. RED: write `TestSecondOpinionPresetResolvesFromBundleOnly` first, before any production edit, and run `go test -C skills/do-work/tools/do-work-cli -count=1 -v -run '^TestSecondOpinionPresetResolvesFromBundleOnly$' ./internal/toolboxcommands`. Expected: it fails at case (a), and the failure text contains `unknown preset "mine"; presets found: bold, visual`. Record the exact failing line.
2. GREEN: make changes 1 to 3, rerun the same command: `--- PASS`, no skip. Then run `go test -C skills/do-work/tools/do-work-cli -count=1 -v -run 'SecondOpinion|AIReportIndex|TestHandlersRegisterCanonicalToolboxCommands' ./internal/toolboxcommands` and `go test -C skills/do-work/tools/do-work-cli -count=1 -run '^TestSecondOpinionJSONKeepsEveryBackendAndJudgeEvidence$' ./internal/resultmodel`: record the PASS, FAIL and SKIP counts (expect zero FAIL and zero SKIP).

## Verify before hand-back (from the worktree root, wall times recorded)
- `bash /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-10-20-200259/REQ-705-probe.sh` run from your worktree root (it reads paths from `git rev-parse --show-toplevel`, so it checks your tree): exit 0. At main HEAD it fails on purpose ("report_index.go still names secondOpinionPresetsFolder").
- `gofmt -l skills/do-work/tools/do-work-cli/internal` prints nothing; `go vet -C skills/do-work/tools/do-work-cli ./internal/toolboxcommands` exits 0.
- `git diff <base> --stat` shows exactly the six files; `git diff --check` clean.

## Hand-back (write the file at the absolute path above)
- Branch name, worktree path, base commit (`git -C <worktree> rev-parse --short HEAD` before your first edit) and every commit hash (one commit preferred).
- File manifest: each file with (new)/(modified)/(deleted) and one line on what changed.
- P-A-U text for [PLAN], [APPLY], [UNIFY]; [UNIFY] lists `git diff --stat`, every check's exit code and wall time, and each file checked.
- Proof record: the RED failure text, the GREEN run, and the package selection counts.
- `## Decisions` (D-08 onwards, because the REQ uses D-01 to D-07; each DECIDE & STATE, or ESCALATE with Value and Risk).
- `## Discovered Tasks` (out-of-scope finds, each ending `→ report only` unless impact-critical; do not fix them).
- Lessons read (satellites and families).
- Anti-bloat check (see above): the `git diff --stat` and the list of anything added that this brief did not name.
- A proposed CHANGELOG entry: a descriptive title that says what shipped, in plain words (for example "Your Second-Opinion Presets Live in the Report Bundle"), one or two sentences on why it matters (no commit step before a run, nothing outside the bundle), then specific bullets, including that `ai-reports/second-opinion-presets/` is no longer read and the catalog no longer skips it. The integrator writes it with the version.
- A proposed lesson bullet, naming the satellite you think it belongs in, in that file's bullet shape (`- [family: <slug>] [REQ-705: <one-line lesson>](<relative archive link; the integrator fixes the path>)`). Say "none" if nothing was learned that a later builder would need.
- Integration seams for REQ-704 (anything the lookup change means for a page-path target) and test wall times.
