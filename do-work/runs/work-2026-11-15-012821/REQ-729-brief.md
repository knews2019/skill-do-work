# Builder Brief — REQ-729 (ai-report-judge documents that a relative --out starts at --repo-root)

You are the builder for one work item of the do-work skill suite. You build in your own Git worktree and hand a branch back. The integrator merges, tests, reviews and releases.

- Operative name and branch: `worktree-agent-REQ-729-judge-relative-out-rule`.
- Worktree (the ONLY tree you write): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-729-judge-relative-out-rule`. The coordinator creates it from main HEAD right before dispatch, with the branch checked out. Do not create it yourself.
- Main tree (READ-ONLY for you): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2`
- The request, with Triage (Route A) and decision D-01: read `do-work/working/REQ-729-judge-relative-out-rule.md` from the MAIN tree (your worktree may hold an older copy).
- Hand-back file (the one main-tree file you may write; never stage or commit it): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-15-012821/REQ-729-handback.md`
- Up to ten builders run at the same time in this run, all from the same main HEAD. Integration is serial.

## Read first (in your worktree)
1. `skills/do-work/crew-members/general.md`, `coding-guardrails.md`, `shared-principles.md`, `communication-style.md`.
2. `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` and `_dev/primes/prime-releases.md` (the REQ's prime files).
3. `CLAUDE.md` (maintainer preferences: simple, delete before you add).
4. Required lesson: `_dev/primes/lessons-releases.md`. Dropped only for budget (24500 tokens): `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`; read only the bullets tagged `[family: rule-direction-checked-against-callers]` (fixed-string grep). They are why D-01 keeps the behavior.

## Decided (do not reopen)
- **D-01 — keep today's behavior, document it.** A relative `--out` stays resolved against `--repo-root`. The same handler resolves a relative `<bundle-dir>` against `--repo-root` (`report_judge.go:110`, `:132`, helper `judgeAbsolutePath` at `:285`), and `ai-report-second-opinion` uses the same helper (`report_second_opinion.go:161`, `:173`, `:385`). Changing only `--out` would give one command two rules. Shipped callers (`ai-report.md:181`, `stakeholder-report.md:44`) pass `--repo-root <project-root>` and no `--out`. No Go behavior change, so no new Go test (the REQ says text-only needs none).

## What to build
One commit, two text edits. Do not touch `judgeAbsolutePath` or any other code.

1. `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_judge.go:39`, the `judgeUsage` constant. Append the rule to the same string, so it reads:
   `Usage: ai-report-judge <bundle-dir> [--out <dir>] [--browser <path>]; a relative <bundle-dir> or --out starts at --repo-root, not the current folder`
   Keep it one line in that constant. If REQ-728 has already added a flag to this line in your base, keep that flag and append the same clause. Run `gofmt` on the file.
2. `skills/do-work-toolbox/actions/ai-report.md`, Step 7 (Render and Judge). Add one new one-line paragraph right after the command block that ends at `:182`, before the paragraph "The command serves the bundle ...":
   Pass `--out <dir>` to choose the output directory. A relative `<bundle-dir>` or `--out <dir>` starts at `--repo-root`, not the current folder.
   Do not rewrite the neighbouring paragraphs.

The GREEN probe matches the pattern `relative .*--out.* starts at .?--repo-root.?, not the current folder` on the usage line and anywhere in `ai-report.md`. Keep that exact wording.

**Not yours:** `VERSION`, `CHANGELOG.md`, `skills/do-work/actions/version.md`, anything under `do-work/` except the hand-back file, `skills/do-work-toolbox/docs/ai-report-guide.md` (the REQ does not name it; if you think it needs the sentence, put it under Discovered Tasks), `stakeholder-report.md`.

## Rules
- Write only `report_judge.go` and `ai-report.md`. If the request's own requirements need another file, proceed and report it under "Scope changes" with the reason.
- Out-of-scope finds go to "Discovered Tasks". Do not fix them.
- Commit by exact path (`git add <paths>`). Never `git add -A` or `git commit -a`. Commit subject starts with `[REQ-729]`. End the commit message with the line `Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>`.
- Do not push. Do not merge. Do not start subagents. Do not edit the main tree except the hand-back file.
- Never run the repository gate `_dev/tests/maintainer-verify.sh`; the integrator runs it. A 30-second per-file time budget that fails once under load is rerun once, and both runs are recorded.

## Verify before you hand back (all from your worktree)
1. In `skills/do-work/tools/do-work-cli`: `gofmt -l ./internal/toolboxcommands` prints nothing; `go vet ./internal/toolboxcommands/`; `go test -count=1 -run 'TestAIReportJudge' ./internal/toolboxcommands/` (the browser tests skip without a browser engine; record which ran and which skipped).
2. `bash _dev/tests/action-shell-blocks.sh` (it scans shipped Markdown fences), since `ai-report.md` changed.
3. The GREEN probe: `bash /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-15-012821/REQ-729-probe.sh` with your worktree as the current directory prints `REQ-729 probe PASS`. At base it fails by design with "judgeUsage in report_judge.go does not state the relative --out rule".
4. Manual RED-to-GREEN check: run `bash skills/do-work/tools/do-work-cli.sh --repo-root "$PWD" --format text ai-report-judge` with no bundle and confirm the usage message now contains the rule.

## Hand-back file content
Use these headings:
- `## Commits` — hash and subject.
- `## Files` — each file with the verb.
- `## Scope changes`
- `## Red-green evidence` — the GREEN probe failing at base and passing at your tip; the usage message before and after.
- `## Verification` — each command above with its result; for a rerun under load, both runs.
- `## P-A-U` — three short lines: PLAN, APPLY, UNIFY; the UNIFY line lists what you checked in `git diff --stat`.
- `## Decisions` — D-02 onward (D-01 is taken), each marked DECIDE & STATE or ESCALATE, with Value and Risk for ESCALATE.
- `## Discovered Tasks`
- `## Integration seams` — see below; name the exact lines you changed in each shared file.
- `## Proposed CHANGELOG entry` — a descriptive title (for example "ai-report-judge States Where a Relative --out Starts") and the line "What changes for an installed copy": the judge's usage text and the ai-report judge step now say that a relative bundle path or `--out` starts at `--repo-root`; behavior is unchanged.
- `## Proposed lesson` — one bullet, or "none".
- `## Lessons read`

## Integration seams
| Shared file | REQ-729 edits | Siblings that edit it |
| --- | --- | --- |
| `report_judge.go` | only the `judgeUsage` constant (line 39); no function changes | REQ-728 (judge renders one named page) almost certainly edits the same `judgeUsage` line to add a page flag, plus `handleAIReportJudge` and `judgeRender`. The second integrator keeps both: REQ-728's flag and REQ-729's trailing clause. |
| `report_judge_test.go` | none | REQ-728 adds tests |
| `skills/do-work-toolbox/actions/ai-report.md` | one sentence in Step 7 near line 181-184 | REQ-728 (judge form), REQ-730 (catalog lists second-opinion pages), REQ-731 (catalog rebuild command) |

Your final message: under 2500 characters. State the branch tip hash, the verification results, and anything the integrator must decide.
