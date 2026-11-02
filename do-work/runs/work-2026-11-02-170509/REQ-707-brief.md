# Builder brief: REQ-707 (ai-report organize files finished months into counted month folders and fixes every link)

- Worktree (your only writable tree): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-707-ai-report-organize-counted-month-folders
- Branch: worktree-agent-REQ-707-ai-report-organize-counted-month-folders, created by the orchestrator from main HEAD right before dispatch (do not create it; commit here; never merge, rebase, push or check out another branch).
- REQ (read-only, main tree): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/working/REQ-707-ai-report-organize-counted-month-folders.md. Read it fully: Detailed Requirements, Constraints, Red-Green Proof, `## Plan` (the Plan agent's plan with decisions D-01 to D-13 and the orchestrator's validation), `## Exploration` (verified function locations, D-14, D-15) and `## Scope` (the nine files).
- UR (read-only): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/user-requests/UR-161/input.md
- Prior REQ (read-only, just released as 0.309.1): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/archive/REQ-706-ai-report-tools-read-month-folders.md. It added `reportMonthFolder` (names `YYYY-MM` and `YYYY-MM_<N>-report(s)` for any N) and the nested catalog walk. Its Discovered Task: duplicate bundle names make supersedes ambiguous, so organize must never create a duplicate name (plan D-04).
- Hand-back file (the ONE main-tree path you may write; never stage or commit it): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-02-170509/REQ-707-handback.md
- Route C, tdd: true, impact-user-visible, effort-substantive, domain backend.

## Rules to load first
/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/general.md, coding-guardrails.md, shared-principles.md, communication-style.md, testing.md, backend.md. Primes: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/_dev/primes/prime-action-files.md, prime-releases.md (read; the release is the orchestrator's), /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/tools/do-work-cli/prime-do-work-cli.md (§ Traps). Lessons: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/_dev/primes/lessons-releases.md (whole file; `canonical-link-outlives-its-target` is this REQ's core risk). From the budget-dropped do-work-cli satellite read only bullets carrying `[family: final-boundary-identity]`, `[family: exact-basename-authority]` and the newest `[family: alternate-writer-contract-drift]` bullet (0.309.1, REQ-706): `grep -n 'final-boundary-identity\|exact-basename-authority\|0.309.1' skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`.

## The change (decided, do not reopen)
Build exactly the REQ's `## Plan` (files 1 to 9, code shape, layout rules, link rules A and B, safety and apply order, action prose, tests). The plan's decisions D-01 to D-13 and the orchestrator's D-14/D-15 are settled. Where the plan leaves a detail open, choose the simplest option and record it as a decision (D-16 onwards). Three clarifications:
1. REQ-706 widened the month-folder rule after the plan was written: `reportMonthFolder` also accepts `_<N>-report` and `_<N>-reports` for any N. The plan's keeper rule ("the existing month folder of M named T, else the first existing month folder of M in byte order, renamed to T") already covers a hand-named or stale-count folder; make sure a test fixture includes one non-canonical name (for example `2026-08_1-report` holding two bundles after the move, which the plan already has).
2. The verb is `ai-report-organize [--dry-run]`, text output as the plan's "Text output" lists, ending (real run only) with the `commit paths:` block the action commits by exact path. The layout planner stays unexported (D-14); REQ-708 decides how install and update reach it, so keep `planReportOrganizeLayout(repositoryRoot, currentMonth string)` + `hasWork()` free of Git and file scans as the plan says.
3. Keep the action prose short (the plan's "Action prose" section): the program does the mechanics; the form states it is the only command that moves reports or rewrites links, and that it never runs automatically.

## Anti-bloat (YAGNI); the maintainer watches this
- Build the plan, not more: honor its YAGNI list (no `--commit` flag, no rollback/transaction code, no new resultmodel struct, no Just recipe, no exported API or new package, no CSS url()/srcset/script parsing, no percent-decoding, no symlink retargeting, no deeper nesting, no hint printing).
- Two new tests (plus the registration test update), each pinning named failures; no decorative tests and no unit tests for the rewriter beyond the two.
- Do not fix adjacent things you notice; record them as discovered tasks.
In the hand-back, paste `git diff --stat` and list every function, type, constant, option, file or test you added that the plan did not name, each with a one-line reason (or remove it).

## Write boundary
Exactly the nine files in `## Scope`. Anything else: stop and say so in the hand-back.

## Hard rules
- Every commit subject on your branch starts with `[REQ-707]`.
- Never create, edit, stage or commit any path under `do-work/`, in the worktree or anywhere else, except the one hand-back file above (write it, never stage or commit it). Never run organize against the main repository or your worktree's own `ai-reports/`; tests use temporary fixture repositories only.
- Never merge, rebase, push, or check out another branch. Never run `recover`, `advance`, any `do-work-cli` lifecycle command, or `_dev/tests/maintainer-verify.sh`.
- Do not touch `CHANGELOG.md`, `skills/do-work/CHANGELOG.md`, `VERSION`, `skills/do-work/VERSION`, `skills/do-work/actions/version.md`.
- Scratch files go in `mktemp -d` directories outside both trees and are never committed. Never `rm -rf` anything outside your own `mktemp -d` directory, and do not chain a removal in front of another command (a stalled permission prompt cost the last reviewer its turn).
- Go: focused `go test -run` and `go vet` on the touched package only. Keep the two new tests well under 30 s together.

## Integration seam
No builder runs beside you. REQ-708 (one-line organize hint from install, update, `ai-report index` and new reports) builds after you are released and needs one shared check: the layout planner's `hasWork()`. Keep the planner a plain package-level function with a clear name and no Git dependency.

## Proof to run and record (the REQ's Red-Green Proof)
1. RED: write `TestAIReportOrganizeFilesFinishedMonthsAndKeepsLinksOpen` and `TestAIReportOrganizeRecountsMonthFoldersAndRefusesUnsafeTargets` and the registration-test update first, before any production edit, calling the handler through `Handlers()["ai-report-organize"]` so RED is a test failure ("not registered"), not a compile error. Run `go test -C skills/do-work/tools/do-work-cli -count=1 -v -run '^(TestAIReportOrganize.*|TestHandlersRegisterCanonicalToolboxCommands)$' ./internal/toolboxcommands` and record the failing lines.
2. GREEN: implement, rerun: every `--- PASS`, no skip. Then `go test -C skills/do-work/tools/do-work-cli -count=1 -run 'AIReport|Handlers|Architecture' ./internal/toolboxcommands` and record the counts.

## Verify before hand-back (from the worktree root)
- `bash /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-02-170509/REQ-707-probe.sh` run from your worktree root: exit 0 (it fails at main HEAD on purpose).
- `gofmt -l skills/do-work/tools/do-work-cli/internal` prints nothing; `go vet -C skills/do-work/tools/do-work-cli ./internal/toolboxcommands` exits 0.
- `bash _dev/tests/contract-regressions.sh`, `bash _dev/tests/shipped-package-reference-contract.sh`, `bash _dev/tests/action-shell-blocks.sh` and `bash _dev/tests/prescribed-shell-canonicalization.sh` exit 0 in your worktree (the action gains a prescribed shell block).
- `git diff <base> --stat` shows exactly the Scope files; `git diff --check` clean.

## Hand-back (write the file at the absolute path above)
- Branch name, worktree path, base commit and every commit hash.
- File manifest: each file with (new)/(modified)/(deleted) and one line on what changed.
- P-A-U text for [PLAN], [APPLY], [UNIFY]; [UNIFY] lists `git diff --stat`, every check's exit code and wall time, and each file checked.
- Proof record: the RED failure text, the GREEN run, the package selection counts.
- `## Decisions` (D-16 onwards; each DECIDE & STATE, or ESCALATE with Value and Risk).
- `## Discovered Tasks` (out-of-scope finds, each ending `→ report only` unless impact-critical).
- Lessons read; anti-bloat check.
- A proposed CHANGELOG entry: a descriptive title that says what shipped (the plan proposes "Report Organizer Files Finished Months Into Counted Folders and Fixes Every Link"), one or two sentences on why it matters, then specific bullets, including what the link guarantee does not cover.
- A proposed lesson bullet in the fitting satellite's bullet shape, or "none".
- The exact function name and signature REQ-708 should call for its check.
