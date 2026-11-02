# Builder brief: REQ-708 (install, update and new reports recommend ai-report organize when finished months are still loose)

- Worktree (your only writable tree): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-708-recommend-ai-report-organize-on-install-update
- Branch: worktree-agent-REQ-708-recommend-ai-report-organize-on-install-update, created by the orchestrator from main HEAD right before dispatch (do not create it; commit here; never merge, rebase, push or check out another branch).
- REQ (read-only, main tree): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/working/REQ-708-recommend-ai-report-organize-on-install-update.md. Read it fully: Detailed Requirements, Constraints, Red-Green Proof, `## Triage`, `## Exploration` (key code with paths, decisions D-01 to D-10) and `## Scope` (the ten files).
- UR (read-only): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/user-requests/UR-161/input.md
- Prior REQs (read-only, released): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/archive/REQ-706-ai-report-tools-read-month-folders.md (0.309.1) and /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/archive/REQ-707-ai-report-organize-counted-month-folders.md (0.309.2, the organize verb; its `planReportOrganizeLayout` + `hasWork()` is your check).
- Fuller exploration notes (read-only): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-02-170509/REQ-708-exploration-draft.md (written against an earlier main; line numbers may have moved, and its decision numbers differ from the REQ's: the REQ's `## Exploration` is authoritative).
- Hand-back file (the ONE main-tree path you may write; never stage or commit it): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-02-170509/REQ-708-handback.md
- Route B, tdd: true, impact-user-visible, effort-substantive, domain backend.

## Rules to load first
/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/general.md, coding-guardrails.md, shared-principles.md, communication-style.md, testing.md, backend.md. Primes: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/_dev/primes/prime-action-files.md, prime-shell-commands.md (the writers gain a prescribed shell line), prime-releases.md (read; the release is the orchestrator's), /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/tools/do-work-cli/prime-do-work-cli.md (§ Traps), /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/tools/prime-do-work-update.md. Lessons (required): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/tools/lessons-do-work-update.md and /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/_dev/primes/lessons-releases.md (both whole, small). From the budget-dropped do-work-cli satellite read only bullets carrying `[family: cross-module-fact-handoff]` or `[family: closed-enumeration-for-a-condition]`.

## The change (decided, do not reopen; D-01 to D-10 in the REQ's Exploration)
Paths relative to your worktree root.
1. `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_organize.go`: one hint constant with exactly the text in D-03; one exported `ReportOrganizeHint(repositoryRoot string) string` returning that line when `planReportOrganizeLayout(repositoryRoot, toolboxNow().Format("2006-01"))` succeeds and `hasWork()` holds, else "" (an error yields "", D-03); and `ai-report-organize --check`: runs only the layout planner, prints the hint line or nothing, writes nothing, exit 0 (D-09, D-10). Usage text: `ai-report-organize [--dry-run | --check]`.
2. `report_index.go` `handleAIReportIndex` no-arg branch: on success append the hint line to the text output (D-07). Not in `reportIndexWrite`, not in `--find`.
3. `skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands.go`: `handleInstallSuite` and `handleUpdateSuite` write `toolboxcommands.ReportOrganizeHint(<project root>)` as one narration line right after the transaction succeeds (install: not when the install was cancelled/declined; update: every success), D-04/D-05. Use the project root the command operates on, not the skill root. Never change an outcome or exit status because of the hint.
4. Prose, one short addition each, using the same prescribed line shape as the existing `ai-report-index` block in `ai-report.md`:
   ```bash
   <skill-root>/../do-work/tools/do-work-cli.sh --repo-root <project-root> --format text ai-report-organize --check
   ```
   - `skills/do-work-toolbox/actions/ai-report.md` Step 8 (covers default, `--kind`, revise): run the check after the bundle is written and print its line, if any, just before the final bundle-path line (the Output Format's last-line rule stays true). Not second-opinion, judge, find or organize (D-08).
   - `skills/do-work-toolbox/actions/architecture-report.md` Step 7: the same; the hint is an allowed extra line beyond the "at most eight lines".
   - `skills/do-work-toolbox/actions/stakeholder-report.md` Step 6: print the hint line, if any, before the path line, because callers take the last line.
   - `skills/do-work/actions/version.md` step 4: one clause: relay the engine's `ai-report organize` recommendation line when it prints one.
   - `skills/do-work-toolbox/docs/ai-report-guide.md`: one sentence in "Organizing Reports by Month": install, update, `ai-report index` and every new report print a one-line hint while a finished month is still loose; nothing runs organize for you. Mention the D-06 limit only if it fits in one clause.
5. Tests (each with a `// Pins:` comment naming the real failure):
   - `report_organize_test.go`: `TestReportOrganizeHintPrintsOnlyWhenWorkRemains` (freeze `toolboxNow` at 2026-11-02 like the existing organize tests). Fixture A: `ai-reports/2026-09-27_2200_a/index.html` → the `ai-report-index` text output ends with the hint line, and `ai-report-organize --check` prints exactly that line and creates no file. Fixture B: `ai-reports/2026-11-01_0900_current/` plus `ai-reports/2026-09_1-report/2026-09-27_2200_a/` → neither prints it and `--check` prints nothing.
   - `suite_commands_test.go`: `TestInstallAndUpdateRecommendOrganizeWhenReportsAreLoose`: a project holding a loose `ai-reports/2026-09-27_2200_a/index.html`; `install-suite` narration holds `do-work-toolbox ai-report organize` exactly once and stdout holds none; then `update-suite` (reuse the existing local archive server helper from the update tests) holds it once. The update's clock is the real one, so date the loose report in a month that is certainly finished (2026-09 is fine while the real date is after September 2026; if `toolboxNow` cannot be frozen across packages, pick a date far in the past such as 2001-01-15). Add one absence assertion to `TestInstallNarrationNeverReachesStandardOutput` (no `ai-reports/` → no hint).

## Anti-bloat (YAGNI); the maintainer watches this
- No new verb, Just recipe, help line, JSON field, info finding, result struct field, package or shell-launcher edit; no hint in second-opinion, judge, find or organize; no report count in the hint; never run organize automatically.
- Two new tests plus one assertion; no decorative tests.
In the hand-back, paste `git diff --stat` and list every function, constant, option, file or test you added that this brief did not name, each with a one-line reason.

## Write boundary
Exactly the ten files in `## Scope`. Anything else: stop and say so in the hand-back.

## Hard rules
- Every commit subject on your branch starts with `[REQ-708]`.
- Never create, edit, stage or commit any path under `do-work/`, in the worktree or anywhere else, except the one hand-back file above (write it, never stage or commit it).
- Never merge, rebase, push, or check out another branch. Never run `recover`, `advance`, any `do-work-cli` lifecycle command, `_dev/tests/maintainer-verify.sh`, or a heavy-only script (`install-suite-behavior.sh`, `update-script-behavior.sh`, `staged-skills-contract.sh` exit 2 as heavy-only; do not run them).
- Do not touch `CHANGELOG.md`, `skills/do-work/CHANGELOG.md`, `VERSION`, `skills/do-work/VERSION`, or the `**Current version**:` line of `skills/do-work/actions/version.md` (you edit only step 4's prose there).
- Scratch files go in `mktemp -d` directories outside both trees. Never chain an `rm -rf` in front of another command.
- Go: focused `go test -run` and `go vet` on the two touched packages only.

## Integration seam
No builder runs beside you; this is the run's last REQ. None expected.

## Proof to run and record (the REQ's Red-Green Proof)
1. RED: write both tests and the absence assertion first, before any production edit, and run `go test -C skills/do-work/tools/do-work-cli -count=1 -v -run '^(TestReportOrganizeHintPrintsOnlyWhenWorkRemains|TestInstallAndUpdateRecommendOrganizeWhenReportsAreLoose|TestInstallNarrationNeverReachesStandardOutput)$' ./internal/toolboxcommands ./internal/suiteinstall`. If a test cannot compile before the exported function exists, call it through the command handlers instead so RED is a test failure; record the failing lines.
2. GREEN: implement, rerun: every `--- PASS`, no skip. Then `go test -C skills/do-work/tools/do-work-cli -count=1 ./internal/suiteinstall` and `-run 'AIReport|Organize|Handlers' ./internal/toolboxcommands`; record counts.

## Verify before hand-back (from the worktree root)
- `bash /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-02-170509/REQ-708-probe.sh` from your worktree root: exit 0 (it fails at main HEAD on purpose).
- `gofmt -l skills/do-work/tools/do-work-cli/internal` prints nothing; `go vet -C skills/do-work/tools/do-work-cli ./internal/toolboxcommands ./internal/suiteinstall` exits 0.
- `bash _dev/tests/contract-regressions.sh`, `shipped-package-reference-contract.sh`, `action-shell-blocks.sh`, `prescribed-shell-canonicalization.sh` and `suite-manifest-contract.sh` exit 0 in your worktree.
- `git diff <base> --stat` shows exactly the Scope files; `git diff --check` clean.

## Hand-back (write the file at the absolute path above)
- Branch, worktree path, base commit and every commit hash.
- File manifest: each file with (new)/(modified)/(deleted) and one line on what changed.
- P-A-U text for [PLAN], [APPLY], [UNIFY]; [UNIFY] lists `git diff --stat`, every check's exit code and wall time, and each file checked.
- Proof record: RED failure text, GREEN run, counts.
- `## Decisions` (D-11 onwards; DECIDE & STATE, or ESCALATE with Value and Risk).
- `## Discovered Tasks` (each ending `→ report only` unless impact-critical).
- Lessons read; anti-bloat check.
- A proposed CHANGELOG entry: a descriptive title that says what shipped (for example "Install, Update and New Reports Recommend Organizing Finished Months"), one or two sentences on why it matters, then specific bullets, including the D-06 first-update limit.
- A proposed lesson bullet in the fitting satellite's shape, or "none".
