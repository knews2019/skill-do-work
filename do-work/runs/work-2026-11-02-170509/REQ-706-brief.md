# Builder brief: REQ-706 (ai-report tools read report bundles inside month folders)

- Worktree (your only writable tree): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-706-ai-report-tools-read-month-folders
- Branch: worktree-agent-REQ-706-ai-report-tools-read-month-folders, created by the orchestrator from main HEAD right before dispatch (do not create it; commit here; never merge, never rebase, never push, never check out another branch).
- REQ (read-only, main tree): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/working/REQ-706-ai-report-tools-read-month-folders.md. Read it fully: Detailed Requirements, Constraints, Red-Green Proof, `## Exploration` (decisions D-01 to D-06, exact lines) and `## Scope` (the file list).
- UR (read-only): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/user-requests/UR-161/input.md
- Original REQ (read-only): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/archive/UR-144/REQ-655-*.md (added `ai-report-index`; use `ls` to find the exact name).
- Hand-back file (the ONE main-tree path you may write; never stage or commit it): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-02-170509/REQ-706-handback.md
- Route B, tdd: true, impact-user-visible, effort-substantive, domain backend.

## Rules to load first
/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/general.md, coding-guardrails.md, shared-principles.md, communication-style.md, testing.md, backend.md. Primes: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/_dev/primes/prime-action-files.md and prime-releases.md (read; the release is the orchestrator's), /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/tools/do-work-cli/prime-do-work-cli.md (§ Traps). Lessons: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/_dev/primes/lessons-releases.md (whole file, small). From the budget-dropped do-work-cli satellite read only the bullets carrying `[family: catalog-tie-order]` (`grep -n 'catalog-tie-order' skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`).

## The change (decided, do not reopen)
Paths are relative to your worktree root. Line numbers are at main d231e2fa.
1. `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index.go`:
   - Add one package-level helper (D-01), for example `reportMonthFolder(name string) bool`, backed by one regexp `^[0-9]{4}-(0[1-9]|1[0-2])(_1-report|_[0-9]+-reports)?$`. Doc comment: a month folder is a direct child of ai-reports/ that holds bundles; plain YYYY-MM, YYYY-MM_1-report or YYYY-MM_<N>-reports; whether N matches the bundle count is not checked here.
   - `reportIndexWalk` (:79-92): for each non-dot directory entry, when the name is a month folder, read its children with `os.ReadDir` and append every non-dot child directory as a bundle with `Path` = `ai-reports/<month>/<child>`; otherwise keep today's behavior. Loose files and symbolic links are never bundles, at either level (`DirEntry.IsDir()` is false for a symlink, so a symlinked month folder is skipped). One level only (D-02). Rewrite the doc comment (:79-80) to say this.
   - `readReportBundle` (:119): it takes the folder name and builds the path from it. Change it so it can read a nested bundle (for example pass the bundle's relative directory under ai-reports, `<folder>` or `<month>/<folder>`, and use `path.Base` of it wherever the folder NAME is read: date, `reportKind`, linked ids, default title). The catalog fields and the flat output must stay byte-identical.
   - Sort, supersedes resolution, `reportPageLink`, `reportIndexFind`: no change expected (they key on `path.Base` or trim the `ai-reports/` prefix). Verify.
2. `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_second_opinion.go:387-388`, `secondOpinionBundle`: accept a bundle whose parent is `ai-reports` OR `ai-reports/<name>` where `reportMonthFolder(<name>)` holds (D-04). New refusal text: `bundle must be a physical report bundle directly under ai-reports or under a month folder`. Keep both `validateNoLinkedAncestors` calls (they refuse a symlinked month folder) and the dot-name refusal.
3. Tests (D-06), each with a `// Pins:` comment naming the real failure:
   - `report_index_test.go`: `TestAIReportIndexReadsBundlesInsideMonthFolders`. Fixture: `ai-reports/2026-09_1-report/2026-09-27_2200_a/index.html` and `ai-reports/2026-11-01_0900_b/index.html` (the REQ's RED case). Run the index handler: exactly two bundles, `b` first, paths `ai-reports/2026-11-01_0900_b` and `ai-reports/2026-09_1-report/2026-09-27_2200_a`, no bundle named after the month folder, and the `index.html` page links `2026-09_1-report/2026-09-27_2200_a/index.html`. Repeat (subtest or loop) with a plain `2026-09/` folder. Also in the same test: a loose file inside the month folder is skipped; a symlinked month folder (`os.Symlink`) adds no bundle; `--find a` (or the bundle's folder text) lists the nested bundle; a flat `-rev1` bundle whose `ai-report-supersedes` meta names the nested folder resolves `supersedes`/`superseded_by` across flat and nested. Reuse the existing fixture helpers in that file.
   - `report_second_opinion_test.go`: `TestSecondOpinionBundleInsideMonthFolder`. Reuse `secondOpinionFixture` and its helpers: move or create the bundle under `ai-reports/2026-09_3-reports/<bundle>` (commit it, so the outside-bundle check is clean), then (a) a `--prepare` run on that nested bundle path passes the bundle check (state `prepared`); (b) `latest` resolves the nested bundle when it is the newest; (c) the same bundle behind a symlinked month folder is refused. Keep it to what these three cases need.
   - The existing catalog and second-opinion tests must pass unchanged: they are the byte-identical proof for a repository with no month folders.
4. `skills/do-work-toolbox/actions/ai-report.md`:
   - :42 catalog sentence: say it catalogs every bundle directly under `ai-reports/` and inside a month folder (`ai-reports/<YYYY-MM...>/`, named in one short clause), whatever its folder naming style.
   - :54 revise step 1: replace "`<dir>` is `ai-reports/<folder>` or a bare `<folder>` and must be a directory directly under `ai-reports/`" with: `<dir>` is `ai-reports/<folder>`, `ai-reports/<month folder>/<folder>`, or a bare `<folder>`, and must be a bundle the catalog lists, directly under `ai-reports/` or inside a month folder; anything else stops with one line. Add one clause that the new revision is still written directly under `ai-reports/` (D-05).
   - :72 second-opinion: replace "It must be a physical direct child of `ai-reports/` with a regular `index.html`" with: It must be a physical bundle directory directly under `ai-reports/` or inside a month folder, with a regular `index.html`.
   - Use the literal words `month folder` (the GREEN probe checks for them). Define it once where it first appears in the file: a direct child of `ai-reports/` named `YYYY-MM`, `YYYY-MM_1-report` or `YYYY-MM_<N>-reports`, which holds the reports of one finished month.
5. `skills/do-work-toolbox/docs/ai-report-guide.md`: :95 (catalog) now says month folders are read as containers; :103 (revise) says the new revision is written directly under `ai-reports/` even when the old one sits in a month folder; the second-opinion paragraph near :115 accepts a bundle in a month folder. Add at most one short paragraph describing the layout: the current month's reports sit directly under `ai-reports/`, a finished month's reports may sit in `ai-reports/YYYY-MM_N-reports/` (`_1-report` for one), and the tools read both. Do NOT describe an organize command: it does not exist yet (REQ-707 adds it). Use the literal words `month folder`.

## Anti-bloat (YAGNI); the maintainer watches this
- Smallest change that does what the REQ asks. Prefer reusing existing code.
- One new helper (the month-folder test) and one regexp; no new flags, options, files, config or environment variables.
- Two new tests, each pinning named failures; no decorative tests.
- Do not touch report writers or anything about organizing/moving reports.
In the hand-back, paste `git diff --stat` and list every function, constant, option, file or test you added that this brief did not name (expected: none beyond the helper and the two tests), each with a one-line reason.

## Write boundary
Exactly the six files in `## Scope`. Anything else: stop and say so in the hand-back.

## Hard rules
- Every commit subject on your branch starts with `[REQ-706]`.
- Never create, edit, stage or commit any path under `do-work/`, in the worktree or anywhere else, except the one hand-back file above (write it, never stage or commit it).
- Never merge, rebase, push, or check out another branch. Never run `recover`, `advance`, any `do-work-cli` lifecycle command, or `_dev/tests/maintainer-verify.sh`.
- Do not touch `CHANGELOG.md`, `skills/do-work/CHANGELOG.md`, `VERSION`, `skills/do-work/VERSION`, `skills/do-work/actions/version.md`.
- Scratch files go in `mktemp -d` directories outside both trees and are never committed.
- Go: focused `go test -run` and `go vet` on the touched package only.

## Integration seam
No builder runs beside you. REQ-707 (an `ai-report organize` command that moves finished months into `YYYY-MM_N-reports` folders and rewrites links) builds after you are released and will reuse your month-folder helper: keep it a plain package-level function in `report_index.go` with a clear name. REQ-708 (a one-line organize hint) follows it.

## Proof to run and record (the REQ's Red-Green Proof)
1. RED: write both new tests first, before any production edit, and run `go test -C skills/do-work/tools/do-work-cli -count=1 -v -run '^(TestAIReportIndexReadsBundlesInsideMonthFolders|TestSecondOpinionBundleInsideMonthFolder)$' ./internal/toolboxcommands`. Expected: the catalog test fails because the month folder is listed as a bundle and `a` is missing; the second-opinion test fails with `bundle must be a physical direct child of ai-reports`. Record the exact failing lines.
2. GREEN: make the production changes, rerun: both `--- PASS`, no skip. Then `go test -C skills/do-work/tools/do-work-cli -count=1 -run 'SecondOpinion|AIReportIndex' ./internal/toolboxcommands` and record PASS/FAIL counts.

## Verify before hand-back (from the worktree root)
- `bash /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-11-02-170509/REQ-706-probe.sh` run from your worktree root (it reads `git rev-parse --show-toplevel`): exit 0. It fails at main HEAD on purpose.
- `gofmt -l skills/do-work/tools/do-work-cli/internal` prints nothing; `go vet -C skills/do-work/tools/do-work-cli ./internal/toolboxcommands` exits 0.
- `bash _dev/tests/contract-regressions.sh` and `bash _dev/tests/shipped-package-reference-contract.sh` exit 0 in your worktree (the prose edits are shipped files).
- `git diff <base> --stat` shows exactly the six files; `git diff --check` clean.

## Hand-back (write the file at the absolute path above)
- Branch name, worktree path, base commit and every commit hash (one commit preferred).
- File manifest: each file with (new)/(modified)/(deleted) and one line on what changed.
- P-A-U text for [PLAN], [APPLY], [UNIFY]; [UNIFY] lists `git diff --stat`, every check's exit code and wall time, and each file checked.
- Proof record: the RED failure text, the GREEN run, and the package selection counts.
- `## Decisions` (D-07 onwards; each DECIDE & STATE, or ESCALATE with Value and Risk).
- `## Discovered Tasks` (out-of-scope finds, each ending `→ report only` unless impact-critical; do not fix them).
- Lessons read.
- Anti-bloat check.
- A proposed CHANGELOG entry: a descriptive title that says what shipped (for example "Report Tools Find Reports Filed in Month Folders"), one or two sentences on why it matters, then specific bullets.
- A proposed lesson bullet for the satellite you think fits, in that file's bullet shape, or "none".
