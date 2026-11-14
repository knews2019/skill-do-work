# Builder Brief — REQ-731 (ai-report prints the catalog rebuild command when a new report is missing from the catalog)

You are the builder for one work item of the do-work skill suite. You build in your own Git worktree and hand a branch back. The integrator merges, tests, reviews and releases.

- Operative name and branch: `worktree-agent-REQ-731-report-catalog-rebuild-reminder`.
- Worktree (the ONLY tree you write): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-731-report-catalog-rebuild-reminder`. The coordinator creates it from main HEAD right before dispatch, with the branch checked out. Do not create it yourself.
- Main tree (READ-ONLY for you): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2`
- The request, with its decisions D-01 to D-03, Exploration and Scope: read `do-work/working/REQ-731-report-prints-catalog-rebuild-reminder.md` from the MAIN tree (your worktree may hold an older copy).
- Hand-back file (the one main-tree file you may write; never stage or commit it): `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-15-012821/REQ-731-handback.md`
- Up to ten builders run at the same time in this run, each in its own worktree from the same main HEAD. Integration is serial, in hand-back order.

## Read first (in your worktree)
1. `skills/do-work/crew-members/general.md`, `coding-guardrails.md`, `shared-principles.md`, `communication-style.md`, `testing.md`.
2. `_dev/primes/prime-action-files.md` (you edit an action file), `_dev/primes/prime-releases.md`, and `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (you edit Go in do-work-cli).
3. `CLAUDE.md` (maintainer preferences: simple, delete before you add, focused tests that name the failure they pin).
4. Required lessons: `_dev/primes/lessons-releases.md`. Dropped only for budget, read the matching bullets: `[family: catalog-tie-order]` and `[family: commit-preflight-before-side-effect]` in `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`, and `[family: alternate-writer-contract-drift]` in `_dev/primes/lessons-action-files.md` (find them with a fixed-string grep).

## Decided (do not reopen)
- **Maintainer, 2026-11-14:** print a reminder; creating a report does not rebuild or commit the catalog. No bundle is edited.
- **D-01 — the check is a new read-only flag `ai-report-index --check <bundle>`** in `report_index.go`. No existing command reads `catalog.json` (`--find` walks folders in memory). It reads only `ai-reports/catalog.json`, never walks and never writes. It prints one line when the bundle path is not in the catalog's `bundles[].path`, when the catalog file is missing, or when it cannot be read or parsed; it prints nothing (empty output, success) when the path is listed. Model it on `ai-report-organize --check` (`report_organize.go:97-104`).
- **D-02 — argument and wording:** the argument is the bundle path as the catalog stores it, `ai-reports/<folder>`; trim one trailing `/`. A missing or empty argument returns the usage result. The line names the bundle and the command, in the organize hint's style (`report_organize.go:35`), for example: `ai-reports/<folder> is not in ai-reports/catalog.json yet: run do-work-toolbox ai-report index to rebuild the catalog.` A hand-made catalog (no generator marker) is decoded the same way; if it lists the bundle, print nothing. Keep JSON decoding tolerant of extra fields: REQ-730 adds a field to each bundle in the same wave.
- **D-03 — only `ai-report.md` Step 8.** `architecture-report.md` and `stakeholder-report.md` also write bundles; list them under Discovered Tasks, do not change them.

## What to build
Two commits.

**Commit 1 (RED, test only).** Append `TestAIReportIndexCheckPrintsRebuildLineOnlyWhenBundleMissing` at the END of `report_index_test.go`. Start the function with a comment naming the failure it pins (a new report missing from the catalog with no reminder). Cases, using the existing helpers (`runReportIndex`, `snapshotReportTree`):
1. Catalog built by `ai-report-index`, then a new bundle folder added: `--check ai-reports/<new>` prints exactly one line, ending in `\n`, naming the bundle and `do-work-toolbox ai-report index`.
2. The same bundle after a rebuild: `--check` prints the empty string. A trailing `/` gives the same result.
3. No `catalog.json` at all (only bundle folders): the line prints.
4. Every case writes nothing (snapshot before and after equal) and returns success.
5. `--check` with no bundle argument returns the usage result.
Run it and record the failure.

**Commit 2 (the feature).**
1. `report_index.go`: one new branch in `handleAIReportIndex` (around line 91) for `--check` with exactly one argument; a new function at the END of the file. Update the usage string at line 96 to `Usage: ai-report-index [--find <topic...> | --check <bundle>]`. Add a short comment that `--check` is the report writers' seam and never writes.
2. `skills/do-work-toolbox/actions/ai-report.md` Step 8 (lines 188-200): after the organize check, run the catalog check with the new bundle's path, in a fenced block like the organize one:
   `<skill-root>/../do-work/tools/do-work-cli.sh --repo-root <project-root> --format text ai-report-index --check ai-reports/<folder>`
   Say it prints one line recommending `ai-report index` when the new bundle is not in `ai-reports/catalog.json` (or the catalog does not exist yet), or nothing, and writes nothing; print its line, if any, just before the final bundle-path line, with the organize line. Do not run the check in the revise form (it rebuilds the catalog). Add the reminder to the summary items or the Verification Checklist only if it reads naturally; keep it to one clause.
3. `skills/do-work-toolbox/docs/ai-report-guide.md` "Report Catalog" (lines 94-102): one new sentence or bullet: a new report enters the catalog at the next rebuild (`ai-report index`), and the report's closing summary reminds you while it is missing. Do not rewrite neighbouring paragraphs.

**Not yours:** `VERSION`, both `CHANGELOG.md` files and the version line in `skills/do-work/actions/version.md` (the integrator writes the release); anything under `do-work/`; `ai-reports/` (no bundle, no catalog file); `report_organize.go`; `help.md` (`--check` is a seam, not a user form).

## Rules
- Write only inside the Scope's file list. If the request's own requirements need another file, proceed and report it under "Scope changes" with the reason; otherwise stop and report.
- Shared files: edit only the lines this REQ needs, add new text as its own sentence, bullet or paragraph, never rewrite a neighbour's paragraph for style, and add the Go test as a new function at the end of the test file.
- Out-of-scope finds go to "Discovered Tasks" in the hand-back. Do not fix them.
- Commit by exact path (`git add <paths>`). Never `git add -A` or `git commit -a`. Commit subjects start with `[REQ-731]`. End each commit message with the line `Co-Authored-By: Claude Opus 5.5 <noreply@anthropic.com>`.
- Do not push. Do not merge. Do not start subagents. Do not edit the main tree except the hand-back file.
- Never run the repository gate `_dev/tests/maintainer-verify.sh`; the integrator runs it. A 30-second per-file time budget that fails once under load is rerun once, and both runs are recorded.

## Verify before you hand back (all from your worktree)
1. In `skills/do-work/tools/do-work-cli`: `gofmt -l ./internal/toolboxcommands` prints nothing; `go vet ./internal/toolboxcommands/`; `go test -count=1 ./internal/toolboxcommands/`.
2. `bash _dev/tests/contract-regressions.sh` and `bash _dev/tests/audit-lockins.sh`.
3. By hand in a scratch `mktemp -d` Git repository: create `ai-reports/2026-11-15_1000_a/index.html`, run `ai-report-index`, add `ai-reports/2026-11-15_1100_b/index.html`, then run `--check` for `b` (one line) and for `a` (nothing). Record the output.
4. The GREEN probe: `bash /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-15-012821/REQ-731-probe.sh` run with your worktree as the current directory prints `REQ-731 GREEN probe: PASS`. It fails at base by design (the test does not exist); the pre-dispatch agent confirmed exit 1 at base.

## Hand-back file content
Use these headings:
- `## Commits` — hash and subject each.
- `## Files` — every file created or modified, with the verb.
- `## Scope changes`
- `## Red-green evidence` — the new test failing after commit 1 and passing after commit 2; the GREEN probe failing at base and passing at your tip.
- `## Verification` — each command above with its result; for a rerun under load, both runs.
- `## P-A-U` — three short lines: PLAN, APPLY, UNIFY; the UNIFY line lists what you checked in `git diff --stat`.
- `## Decisions` — D-04 onward (D-01 to D-03 are taken), each marked DECIDE & STATE or ESCALATE, with Value and Risk for ESCALATE.
- `## Discovered Tasks` — include D-03's two other report writers.
- `## Integration seams` — what you changed in each shared file below, with the function or heading and line ranges after your edit.
- `## Proposed CHANGELOG entry` — a descriptive title (for example "New Reports Remind You to Rebuild the Catalog") and the line "What changes for an installed copy": after an ordinary `ai-report`, the summary prints one `ai-report index` line while the new report is missing from `ai-reports/catalog.json`; `ai-report-index --check <bundle>` is new and writes nothing; creation still never rebuilds or commits the catalog.
- `## Proposed lesson` — one bullet, or "none".
- `## Lessons read`

## Integration seams
All ten REQs of this run build from the same main HEAD and integrate serially:
| Shared file | REQ-731 edits | Siblings editing the same file |
| --- | --- | --- |
| `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index.go` | `handleAIReportIndex` (one new branch, usage string) and one new function at file end | REQ-730 (catalog lists second-opinion pages: `pickReportEntry`, `reportBundle`, catalog rendering) |
| `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index_test.go` | one new test function at file end | REQ-730 (new tests at file end) |
| `skills/do-work-toolbox/actions/ai-report.md` | Step 8, after the organize check | REQ-728 and REQ-729 (judge options), REQ-730 (catalog forms) |
| `skills/do-work-toolbox/docs/ai-report-guide.md` | "Report Catalog" | REQ-729, REQ-730 (likely "Report Catalog" too) |

Your final message: under 2500 characters. State the branch tip hash, the verification results, and anything the integrator must decide.
