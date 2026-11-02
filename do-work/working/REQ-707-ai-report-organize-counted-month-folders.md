---
id: REQ-707
title: 'ai-report organize files finished months into counted month folders and fixes every link'
status: claimed
created_at: 2026-11-02T16:55:00Z
user_request: UR-161
domain: backend
prime_files: ["_dev/primes/prime-action-files.md", "_dev/primes/prime-releases.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-706, REQ-708]
batch: ai-report-month-folders
required_lessons: ["_dev/primes/lessons-releases.md"]
write_set: ["skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_organize.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_organize_test.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/commands.go", "skills/do-work/tools/do-work-cli/internal/toolboxcommands/commands_test.go", "skills/do-work-toolbox/actions/ai-report.md", "skills/do-work-toolbox/actions/help.md", "skills/do-work-toolbox/SKILL.md", "skills/do-work-toolbox/docs/ai-report-guide.md", "skills/do-work/docs/command-line-guide.md"]
depends_on: [REQ-706]
claimed_at: 2026-11-02T17:55:27Z
route: C
estimate:
  p50_active_minutes: 40
  confidence: medium
  calculated_at: 2026-11-02T17:55:37Z
  basis:
    - Route C
    - 9-file write set
    - 2 subsystems involved
    - 10 acceptance criteria
planning_at: 2026-11-02T17:55:56Z
---
# ai-report Organize Files Finished Months Into Counted Month Folders and Fixes Every Link
## What
Add `do-work-toolbox ai-report organize` for the repository the suite is installed in. It moves every report from a finished month into `ai-reports/YYYY-MM_N-reports/`, keeps the current month's reports directly under `ai-reports/`, and rewrites every link to a moved report so it still opens.
## Why
The user wants the reports folder readable by month, with the report count visible in the folder name for months that are over, and wants a command to do it instead of a project-local script.
## Detailed Requirements
- A new `organize` form in `skills/do-work-toolbox/actions/ai-report.md` (with its routing and help lines), backed by one canonical do-work-cli verb. The moves and link rewrites are mechanical, so a program does them.
- A report's month comes from the `yyyy-mm-dd` date prefix of its folder name. A finished month is any month before the current local month. Current-month reports stay at `ai-reports/<bundle>/`.
- Each finished-month report moves into `ai-reports/<YYYY-MM>_<N>-reports/` (`<YYYY-MM>_1-report` when N is 1), where N is the number of report bundles in that month after the move.
- Recount on every run: a finished-month folder whose name no longer matches its bundle count, or a plain `YYYY-MM` folder, is renamed to the counted name. A month folder for the current month is emptied back to the top level.
- Fix every link to every moved or renamed report across the repository: repo-relative path mentions such as `ai-reports/<bundle>/index.html` in REQs, URs, docs, `blocked_by` lines and `## Reports` history; relative links inside reports that point at another report, such as `../<sibling>/index.html`; and the catalog. Every link that opened before opens after. Binary files are never edited.
- Regenerate `ai-reports/catalog.json` and `ai-reports/index.html` through the catalog command.
- A second run on an organized repository changes nothing.
- `--dry-run` prints the whole plan (moves, renames, and every file whose links change) and writes nothing.
- A report folder without a date prefix, a loose file, or a symbolic link stays in place and is listed in the output.
- Commit as the project's rules say, by exact paths, as ai-report Step 8 does. Never push.
- The user guide explains the layout and the command.
## Constraints
- The current month's reports stay directly under `ai-reports/`. Only finished months move into a month folder, because the count is final only when the month is over. Report writers (ai-report in every form, architecture-report, stakeholder-report) keep writing new bundles directly under `ai-reports/`: no write-side change.
- A finished-month folder is named `YYYY-MM_N-reports` (`YYYY-MM_1-report` for one report). A plain `YYYY-MM` folder, as a consumer's own script makes, is read as a month folder too.
- The bundle contract is unchanged: one folder holding its entry file, the same catalog fields, and the same `yyyy-mm-dd_hhmm_<slug>` names. A month folder is only a parent directory.
- A repository with no month folders behaves exactly as before.
- Only `ai-report organize` moves reports or rewrites links. A report run never does.
- The inbox note was observed against 0.307.0, so its line numbers may have drifted.
## Dependencies
Depends on REQ-706 (ai-report tools read report bundles inside month folders): organize writes the layout that REQ-706 teaches the readers. REQ-708 depends on this.
## Builder Guidance
The user chose: the count in the folder name for finished months, the current month flat, and this command as the only thing that moves reports or rewrites links.
## Red-Green Proof
**RED prompt/case:** A test project has 3 reports from September 2026 and 1 from the current month. One September report links to another September report, and a queue file links to one September report. Run `ai-report organize`.
**Why RED now:** The command does not exist.
**GREEN when:** The 3 September reports move into `ai-reports/2026-09_3-reports/`, the current-month report stays where it is, both links still open, and a second run changes nothing. `--dry-run` shows the plan and writes nothing.
**Validation:** User confirmed.
## Required Lessons — Dropped for Budget
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` as a whole satellite (23234 tokens per `do-work/lessons-index.md`, over the 2000 budget; `slugged: partial`, so no targeted form). Matching reason: its row covers changing do-work-cli internals; families `final-boundary-identity, exact-basename-authority, catalog-tie-order` match this work.
- `_dev/primes/lessons-action-files.md` as a whole satellite (8398 tokens, over the 2000 budget; `slugged: partial`). Matching reason: family `alternate-writer-contract-drift`: a new ai-report form whose rule must match every reader of the bundle layout.
## Full Context
See `do-work/user-requests/UR-161/input.md` for the complete verbatim input: the user's words, the six answers chosen during capture, and the inbox note it adopts in full. No queued candidate in any UR shares this root cause (the queue was empty at capture).
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: "Um, moved on the folder. And if it's an older month, maybe we can also put account how many reports are in that folder. And make sure to also always fix all of the links to it."*

## Triage

**Route: C** - Complex

**Reasoning:** A new do-work-cli verb with three parts (a month-layout planner, a repository-wide link rewriter, an apply step that moves directories) plus a new action form, routing, help and guide text. The link rules and the safety of moving directories need a plan before code, and REQ-708 must be able to reuse the planner as its check.

**Planning:** Required

**Open questions (Step 3.5):** none in the REQ. The user's decisions are in the REQ's Constraints and Builder Guidance (counted folder names for finished months, current month flat, organize is the only command that moves reports or rewrites links, it never runs automatically). Plan decisions start at D-01.

**Test policy for this run (user decision, 2026-11-02):** focused tests only. The gate argv for every evidence step is `bash do-work/runs/work-2026-11-02-170509/REQ-707-gate.sh`; `_dev/tests/maintainer-verify.sh` and heavy lanes are not run.

**Prior work this builds on:** REQ-706 (released as 0.309.1) added `reportMonthFolder` in `report_index.go` (month-folder names `YYYY-MM`, `YYYY-MM_<N>-report(s)` for any N, D-13 there) and the nested catalog walk; its review fixed the architecture-report pre-flight to read month folders. Its Discovered Task (duplicate bundle names across month folders make supersedes ambiguous) is handled here by refusing to create a duplicate name.

<!-- D-XX counter: none used. Next decision: D-01. -->

## Plan

Base: main 3edfd2c9. REQ-706 is merged, so `reportMonthFolder` and the nested catalog walk already exist.

### Files, in build order

Go files are in `skills/do-work/tools/do-work-cli/internal/toolboxcommands/`.

1. NEW `report_organize_test.go`: two tests that call `Handlers()["ai-report-organize"]`. RED is "not registered", not a compile error.
2. `commands_test.go`: add the name; the handler count goes from 11 to 12.
3. `commands.go`: `CommandAIReportOrganize = "ai-report-organize"` and its handler.
4. NEW `report_organize.go`. It reuses `reportMonthFolder`, `reportIndexWrite`, `rootedPublishFile`, `rootedMkdirExclusive`, `validateNoLinkedAncestors`, `rootedfs.Rename`, `auditGit` and `toolboxNow`. `report_index.go` is unchanged.
5. `skills/do-work-toolbox/actions/ai-report.md`: the Input sentence names `organize`. Add `### Organize form` before "Revise form".
6. `skills/do-work-toolbox/actions/help.md`: add `  ai-report organize [--dry-run] File finished months into month folders; fix links`.
7. `skills/do-work-toolbox/SKILL.md`: add `organize reports` to the ai-report routing row.
8. `skills/do-work-toolbox/docs/ai-report-guide.md`: `## Organizing Reports by Month` after `## Report Catalog`, about 8 lines (layout, recount, link promise and limits, dry-run, dirty refusal, commit). REQ-706's layout sentence gains "and `ai-report organize` puts them there".
9. `skills/do-work/docs/command-line-guide.md:26`: replace the stale verb list with "The `ai-report-*` toolbox CLI verbs have no recipe".

### Code shape

```go
// Steps hold mkdir|move|rename|remove with paths as at execution; Left holds path+reason; Edits hold path, final path, before/after bytes.
func planReportOrganizeLayout(repositoryRoot, currentMonth string) (reportOrganizePlan, error) // reads ai-reports/ only
func (plan reportOrganizePlan) hasWork() bool                    // len(Steps) > 0: REQ-708's check
func (plan reportOrganizePlan) finalPath(relative string) string // composes the steps at segment boundaries
func planReportOrganizeLinks(repositoryRoot string, plan *reportOrganizePlan) error // fills Edits; read-only
func applyReportOrganize(repositoryRoot string, plan reportOrganizePlan) ([]resultmodel.RecordedChange, error)
```

The handler takes no arguments or `--dry-run`; anything else gets `usageResult`. It runs these in order: the layout plan with month `toolboxNow().Format("2006-01")`, the link plan, the dirty check, then either the dry-run output or apply followed by `reportIndexWrite`.

### Layout rules

- A bundle is a non-dot child directory (REQ-706 walk). Its month is its name's `^[0-9]{4}-(0[1-9]|1[0-2])-[0-9]{2}` prefix, else its month folder's month.
- Finished = before the current local month; future months count as current.
- Finished month M with N bundles → T = `M_N-reports` (`M_1-report` for 1). Keeper: M's month folder already named T, else M's first month folder in byte order (renamed to T), else a new T. Other bundles of M move into the keeper; a bundle dated another month moves out first. Order: mkdir, moves, rename; `finalPath` composes them.
- Months not finished: bundles in their month folders move to `ai-reports/<bundle>`.
- A month folder left fully empty (dot files too) is removed; one holding other files is listed; a folder with no bundles is never renamed.
- Skip and list a month when two of its bundles share a name or T exists as a file or symlink. No step ever targets an existing path (POSIX `rename` silently replaces an empty directory).
- Listed every run: top-level loose files, symlinks (a symlinked month folder too), undated top-level folders, skipped months. Not listed: dot entries, the catalog pair. No `ai-reports/` → empty plan; a symlinked `ai-reports` → refuse.

### Link rules

**Files scanned:** `auditGit(root, "ls-files", "-z", "--cached", "--others", "--exclude-standard")`, deduplicated (tracked plus untracked non-ignored); regular files only (Lstat). Binary (a NUL in the first 8000 bytes, Git's own test) is never edited. The catalog pair is skipped (regenerated). Outside `ai-reports/`, only files containing `ai-reports/` are parsed. Outside a Git repository, organize refuses.

**Rule A: repository paths, in every scanned file.**
1. Find each `ai-reports/` whose preceding byte is not in `[A-Za-z0-9_.-]`.
2. Take the run after it. The run stops at whitespace or any of ``"'`<>()[]{}|,;*:#?\``, and trailing dots are trimmed.
3. If an old moved or renamed path is a prefix of the run, ending at `/` or at the end of the run, replace that prefix with `finalPath(prefix)`.

The match is exact: `X-rev1` never matches `X` (exact-basename-authority). The rule covers `:19` suffixes, `/*`, `../ai-reports/`, absolute paths, `file://` links and `blob/main/ai-reports/` URLs. Shipped lessons in this repo link with that last shape, which is the canonical-link-outlives-its-target trap.

**Rule B: relative links, in files inside `ai-reports/` only.**
1. Tokens are HTML `href=` and `src=` values (either quote style) and Markdown `](target)`. A token is skipped when it has a scheme or starts with `/`, `#` or `?`.
2. Split off the query and fragment.
3. Resolve the path from the file's current folder. Skip the link if it climbs above the repository root, or if its target does not exist now (root `Lstat`).
4. If the link, resolved from `finalPath(file)`, still reaches `finalPath(target)`, keep its bytes.
5. Otherwise write `filepath.Rel(dir(finalPath(file)), finalPath(target))`, keeping the trailing `/`, the query and the fragment.

Rule A skips any span Rule B rewrote. Files outside `ai-reports/` need only Rule A, because a relative path from outside can reach a report only by naming `ai-reports/`.

**Promise: every link that opened before opens after.** A link that was already broken is never touched, and a link that still opens keeps its bytes. That covers links within one bundle, siblings that move together, and `../self/x.png`.

**Not handled (the guide says so):** symlinks to a moved report (never followed or retargeted); CSS `url()`, `srcset`, script-built links, encoded paths, Markdown reference links and autolinks, relative paths in prose; ignored or binary files and other clones, branches, PRs, bookmarks; links to a removed current-month folder; another repository's identical report name (rewritten anyway); a count change between an interrupted run and its rerun.

### Safety and apply

**Dirty check.** Before any write, and in dry-run too, organize runs `git --literal-pathspecs status --porcelain=v1 -z --untracked-files=all -- <step Froms, removed folders, edit paths>`. The catalog pair is excluded. Any record refuses with `AI-REPORT-ORGANIZE-DIRTY` and names the paths. The check is needed for two reasons: the exact-path commit must not sweep in someone else's edit, and it makes Git the undo for an interrupted run. Uncommitted changes elsewhere in the tree are not organize's business.

**Apply order.**
1. Edits, at the files' current paths. Re-read each file, compare it with `Before`, then write with `rootedPublishFile(replace=true)` and the file's own mode.
2. `rootedMkdirExclusive`.
3. Moves.
4. Renames. Before each move and rename, run `validateNoLinkedAncestors` on both ends, require the destination to be absent and the source to be a real directory, then call `rootedfs.Rename` (final-boundary-identity).
5. Removes (`os.Root.Remove`, which fails on a folder that is not empty).
6. `reportIndexWrite`.

On error the run stops and lists what was done; a rerun finishes. Edits go first because a rewritten link that does not open from its current place is never touched again; moves first would let a crash strand links. A hand-made catalog refusal keeps the moves and yields outcome `findings` (revise step 6).

**Output:** text lines `moved A -> <final>`, `renamed F -> T`, `removed F (empty month folder)`, `links <final path> (<n>)`, `left <path> (<reason>)`, the catalog line, then `commit paths:` (one sorted unique path per line: every From and final To, removed folders, edits' final paths, the catalog pair). Dry-run: `dry run: nothing written` plus the same lines, no commit block. No work: `nothing to organize` plus the left lines, nothing written (not even the catalog). JSON: `Changes` (From `deleted`, To `created`, edits `modified`, catalog) and `SkippedWork` for left items.

### Action prose

`### Organize form` has one paragraph, then the launcher block (`... --format text ai-report-organize [--dry-run]`), then four steps:
1. Print the output. Dry-run stops here.
2. On a refusal, report it and stop.
3. Commit as Step 8 says, by exactly the `commit paths:`: `git add -A -- <paths>`, then `git commit -m "Organize ai-reports by month" -- <paths>`. Report the hash, or print the command. I checked this in a scratch repo: it records a rename and leaves other staged files out.
4. Never push.

The form is the only thing that moves reports or rewrites links. Steps 1 to 8 do not apply.

### Tests (each well under 5 s)

**`TestAIReportOrganizeFilesFinishedMonthsAndKeepsLinksOpen`** pins the REQ's proof: moves break links, the second run must change nothing, and dry-run must write nothing.

Fixture (`toolboxNow` frozen at 2026-11-02, committed via `toolboxTestRepository`): `ai-reports/2026-09-01_1200_alpha/index.html` linking `../2026-09-02_1200_beta/index.html#r1`, `../../do-work/queue/REQ-001-sample.md` and `shot.png`; September `beta` and `gamma`; `2026-11-01_0900_current` linking `../2026-09-03_1200_gamma/index.html`; a queue file naming the beta path plus a `-rev1` negative control; a binary file holding the beta path and a NUL.

Steps: (1) dry-run lists 3 moves and 3 edits, and the tree snapshot is unchanged; (2) real run: September is in `2026-09_3-reports/`, current stays at the top; the sibling link and `shot.png` are byte-identical; the escape link becomes `../../../do-work/...`; current links `../2026-09_3-reports/...`; the queue file is rewritten, `-rev1` and the binary are not; every rewritten link resolves; the catalog lists 4 bundles; (3) commit exactly the `commit paths:`, then `git status --porcelain --untracked-files=all` is empty (the list is complete); (4) second run prints `nothing to organize` and the snapshot is unchanged.

**`TestAIReportOrganizeRecountsMonthFoldersAndRefusesUnsafeTargets`** pins a stale count, a plain `YYYY-MM` folder, the singular name, a current-month folder, folder-level mentions, a rename onto an existing path, and an uncommitted edit swept into the commit.

Fixture → expected: `2026-08_1-report/a` plus loose `2026-08-20_1000_b` → `2026-08_2-reports/`; `2026-07/c` → `2026-07_1-report/`; `2026-11/d` → top level, `2026-11` removed; `2026-06-01_1000_e` plus a regular file `2026-06_1-report` → June skipped and listed; `notes.md`, `undated-folder/`, symlink `linked` → listed, untouched; `docs/notes.md` naming `ai-reports/2026-08_1-report/a/index.html` → rewritten.

The test first runs with `docs/notes.md` uncommitted. Organize must refuse, name that file and change nothing. The test then commits it and runs again.

The existing `TestHandlersRegisterCanonicalToolboxCommands` is updated for the new name and count.

**RED** (fails "not registered"): `go test -C skills/do-work/tools/do-work-cli -count=1 -run '^(TestAIReportOrganize|TestHandlersRegisterCanonicalToolboxCommands)' ./internal/toolboxcommands`. **GREEN:** the same run passes. Then run `-run 'AIReport|Handlers'`, `gofmt -l`, `go vet`, `contract-regressions.sh`, `shipped-package-reference-contract.sh` and `staged-skills-contract.sh`.

### Decisions

The rules above are DECIDE & STATE. Contestable ones are marked.

- **D-01:** the verb is `ai-report-organize` with only a `--dry-run` flag. There is no `--commit`: the action commits as Step 8 does.
- **D-02:** rename the keeper rather than refill it, so `.DS_Store` and loose files travel and no orphan folder stays. Duplicates or a file-shaped target skip only their month (one broken pipe).
- **D-03:** the file set comes from Git; binary means a NUL in the first 8000 bytes; edits go before moves; no rollback (Git plus a rerun is the recovery); the catalog pair is regenerated, never scanned.
- **D-04 (contestable):** Rule A accepts `/` before `ai-reports/`, so URLs and absolute paths are fixed. The accepted cost: another repository's identical report name is rewritten too.
- **D-05 (contestable):** the dirty refusal covers only the paths organize touches, including untracked files inside a moving bundle. The alternative was to proceed and leave dirty files out of the commit. That would commit a moved report while its link fix stays uncommitted.
- **D-06 (ESCALATE to REQ-708):** `do-work update` lives in `internal/suiteinstall`, so importing `toolboxcommands` there would be a command-to-command import.
  - Value: REQ-708 picks the cheapest seam. It can call `ai-report-organize --dry-run`, export a wrapper, or move the layout planner into a leaf package.
  - Risk now: none. The planner stays unexported and free of command plumbing.
- **D-07:** history files, such as the archive and the changelogs, are rewritten too, because the REQ says "every link".

### YAGNI: do not build

`--commit`, gittransaction, rollback, a new result struct; a Just recipe, an exported API, a new package; symlink retargeting, `url()`/`srcset`/script parsing, percent-decoding; nested month folders, a configurable folder format; report-writer changes, the hint (REQ-708); progress output, a parallel scan, extra rewriter unit tests.

Proposed changelog title: "Report Organizer Files Finished Months Into Counted Folders and Fixes Every Link".

*Generated by Plan agent* (full draft: `do-work/runs/work-2026-11-02-170509/REQ-707-plan-draft.md`, written against main `3edfd2c9`; REQ-706's review fixes landed after it and are compatible: `reportMonthFolder` now also reads `_<N>-report(s)` for any N, which the keeper/rename rules already handle as "an existing month folder of M")

**Plan validation (orchestrator):**
- Requirement coverage: every Detailed Requirement maps to a planned task: the verb and form (files 3-7), finished-month rule and counted names (layout rules), recount, plain `YYYY-MM` rename and current-month folder emptied (layout rules, second test), every link fixed (rules A and B, first test), catalog regenerated (`reportIndexWrite`), idempotent second run and `--dry-run` (first test), undated/loose/symlink listed as left (layout rules, second test), exact-path commit and never push (organize form steps 3-4), guide (file 8). The Constraints (writers unchanged, only organize moves) are in the YAGNI list. REQ-706's discovered task (duplicate bundle names) is covered by D-04.
- Orphan tasks: `skills/do-work/docs/command-line-guide.md:26` is not named by the REQ; it is a stale restatement of the verb list that this REQ would make staler, so it stays (restatement sweep at plan time).
- Scope sanity: ⚠ the plan has 9 files. Quality degrades past 3 tasks; this is not split because the four prose files carry one line or one short section each and the verb is one cohesive unit (planner, rewriter, apply). The builder is told to keep the prose minimal.
- Consumer field contract: the action's commit step consumes the `commit paths:` block (every From, final To, removed folders, edited files, catalog pair), and REQ-708 consumes the layout planner's `hasWork()`; both are named.

## Exploration

**Lessons consult:** `required_lessons` refreshed to `_dev/primes/lessons-releases.md` (666 tokens, whole satellite, read). Its `canonical-link-outlives-its-target` family is this REQ's core risk: a link names a path, and a move renames the path underneath it; the link rewriter must cover shipped GitHub `blob/main/ai-reports/...` links too (plan rule A). The two capture-time drops stay dropped for budget; their matching families were read by hand: `final-boundary-identity` (bind identity before a destructive move; the plan checks destination absence and source type right before each `rootedfs.Rename`), `exact-basename-authority` (rule A matches an exact path prefix at a segment boundary, never a normalized name), `catalog-tie-order` (unchanged sort), and the new REQ-706 `alternate-writer-contract-drift` bullet (grep every code path that lists the reports folder).

**Exploration source:** the Route C Plan agent read the files the plan names (`report_index.go`, `commands.go`, `mutation.go`, `rootedfs`, `audit_metrics.go`, `note.go`, the toolbox prose files and existing report links in this repository) before writing the plan, so no separate Explore agent was spawned; the orchestrator verified the reused functions at main `384d0c3d` (the REQ-707 claim commit, after 0.309.1):
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_index.go`: `reportMonthFolder` (month-folder names, REQ-706 D-01/D-13), `reportIndexWalk` (nested walk), `reportIndexWrite` (catalog regeneration with the hand-made refusal).
- `internal/toolboxcommands/mutation.go:71` `validateNoLinkedAncestors`, `:114` `rootedMkdirExclusive`, `:132` `rootedPublishFile`.
- `internal/rootedfs/rooted_fs.go:32` `Rename(root *os.Root, oldPath, newPath string)`.
- `internal/toolboxcommands/audit_metrics.go:135` `auditGit(root, args...)` (for `ls-files` and `status`).
- `internal/toolboxcommands/note.go:14` `toolboxNow` (the test clock; `note_test.go:36` shows the freeze-and-restore pattern).
- Test helpers: `report_image_test.go:67` `toolboxTestRepository`, `:79` `toolboxTestGit`; registration test `commands_test.go` (`TestHandlersRegisterCanonicalToolboxCommands`).
- Prose anchors: `skills/do-work-toolbox/actions/help.md:15-19` (ai-report lines), `skills/do-work-toolbox/SKILL.md:25` (ai-report routing row), `skills/do-work-toolbox/actions/ai-report.md` (Catalog forms, then Revise form), `skills/do-work-toolbox/docs/ai-report-guide.md` `## Report Catalog` (REQ-706's layout paragraph), `skills/do-work/docs/command-line-guide.md` (toolbox verbs without a recipe).

**Decisions:** the plan's D-01 to D-13 are adopted as written (see `## Plan`). Orchestrator additions:
- **D-14** (DECIDE & STATE): D-12's escalation is settled for REQ-708, not here: REQ-707 keeps the layout planner unexported in `toolboxcommands` and REQ-708 picks how install and update reach it. Reason: no consumer exists yet; YAGNI.
- **D-15** (DECIDE & STATE): this repository's own `ai-reports/` is not organized as part of this REQ. Running organize on a project is the user's call (the REQ: organize never runs automatically).

*Generated by the Plan agent's exploration pass, verified by the orchestrator*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_organize.go` (new) — the organize verb: layout planner, link planner, apply
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_organize_test.go` (new) — two focused tests (Red-Green Proof; recount and refusals)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/commands.go` (modify) — register `ai-report-organize`
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/commands_test.go` (modify) — registration count and name
- `skills/do-work-toolbox/actions/ai-report.md` (modify) — new organize form
- `skills/do-work-toolbox/actions/help.md` (modify) — one help line
- `skills/do-work-toolbox/SKILL.md` (modify) — routing trigger words
- `skills/do-work-toolbox/docs/ai-report-guide.md` (modify) — organizing section
- `skills/do-work/docs/command-line-guide.md` (modify) — toolbox verbs without a recipe

**Files I will NOT touch:** `report_index.go` and every report writer (`ai-report-reference.md`, `architecture-report.md`, `stakeholder-report.md`), the install and update code (REQ-708), the Justfile template, this repository's own `ai-reports/`, and the release files (finalization writes them).

**Acceptance criteria (restated from REQ):**
- [ ] `do-work-toolbox ai-report organize` exists as a form in `ai-report.md` with routing and help lines, backed by one do-work-cli verb.
- [ ] A report's month comes from its folder's `yyyy-mm-dd` prefix; finished = before the current local month; current-month reports stay at `ai-reports/<bundle>/`.
- [ ] Each finished-month report moves into `ai-reports/<YYYY-MM>_<N>-reports/` (`_1-report` when N is 1), N = bundles in that month after the move.
- [ ] Every run recounts: a wrong-count folder or a plain `YYYY-MM` folder is renamed to the counted name; a current-month folder is emptied back to the top level.
- [ ] Every link to a moved or renamed report is fixed across the repository (repo-relative mentions, relative links between reports, the catalog); every link that opened before opens after; binary files are never edited.
- [ ] `ai-reports/catalog.json` and `ai-reports/index.html` are regenerated through the catalog command.
- [ ] A second run on an organized repository changes nothing.
- [ ] `--dry-run` prints the whole plan (moves, renames, every file whose links change) and writes nothing.
- [ ] An undated folder, a loose file or a symbolic link stays in place and is listed.
- [ ] The form commits by exact paths as ai-report Step 8 does, and never pushes; the guide explains the layout and the command.

## Pre-Flight

**Git:** ✓ clean tracked tree at 384d0c3d (the REQ-707 claim commit, after REQ-706's release 0.309.1); only this REQ's orchestrator sections and the run directory's new REQ-707 files were dirty. The user's untracked inbox note stays hidden from discovery by the run's temporary `.git/info/exclude` line.
**Tests baseline:** ✓ `bash do-work/runs/work-2026-11-02-170509/REQ-707-preflight-probe.sh` exit 0, launched (`do-work/working/baseline.json`). Focused repository gate `bash do-work/runs/work-2026-11-02-170509/REQ-707-gate.sh` at 384d0c3d exit 0 (2026-11-02T17:56Z, about 45 s, 1-minute load 3.0); green-gate record satisfied. The GREEN probe `REQ-707-probe.sh` fails at base, as it must (the verb does not exist).
**Dependencies:** ✓ Go toolchain and ShellCheck present (the gate ran both)

*Checked by work action*
