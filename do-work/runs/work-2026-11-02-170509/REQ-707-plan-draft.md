# REQ-707 plan draft: `ai-report organize` (Route C)

Base: main 3edfd2c9. REQ-706 is merged, so `reportMonthFolder` and the nested catalog walk already exist.

## Files, in build order

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

## Code shape

```go
// Steps hold mkdir|move|rename|remove with paths as at execution; Left holds path+reason; Edits hold path, final path, before/after bytes.
func planReportOrganizeLayout(repositoryRoot, currentMonth string) (reportOrganizePlan, error) // reads ai-reports/ only
func (plan reportOrganizePlan) hasWork() bool                    // len(Steps) > 0: REQ-708's check
func (plan reportOrganizePlan) finalPath(relative string) string // composes the steps at segment boundaries
func planReportOrganizeLinks(repositoryRoot string, plan *reportOrganizePlan) error // fills Edits; read-only
func applyReportOrganize(repositoryRoot string, plan reportOrganizePlan) ([]resultmodel.RecordedChange, error)
```

The handler takes no arguments or `--dry-run`; anything else gets `usageResult`. It runs these in order: the layout plan with month `toolboxNow().Format("2006-01")`, the link plan, the dirty check, then either the dry-run output or apply followed by `reportIndexWrite`.

## Layout rules

- A bundle is a non-dot child directory (REQ-706 walk). Its month is its name's `^[0-9]{4}-(0[1-9]|1[0-2])-[0-9]{2}` prefix, else its month folder's month.
- Finished = before the current local month; future months count as current.
- Finished month M with N bundles → T = `M_N-reports` (`M_1-report` for 1). Keeper: M's month folder already named T, else M's first month folder in byte order (renamed to T), else a new T. Other bundles of M move into the keeper; a bundle dated another month moves out first. Order: mkdir, moves, rename; `finalPath` composes them.
- Months not finished: bundles in their month folders move to `ai-reports/<bundle>`.
- A month folder left fully empty (dot files too) is removed; one holding other files is listed; a folder with no bundles is never renamed.
- Skip and list a month when two of its bundles share a name or T exists as a file or symlink. No step ever targets an existing path (POSIX `rename` silently replaces an empty directory).
- Listed every run: top-level loose files, symlinks (a symlinked month folder too), undated top-level folders, skipped months. Not listed: dot entries, the catalog pair. No `ai-reports/` → empty plan; a symlinked `ai-reports` → refuse.

## Link rules

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

## Safety and apply

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

## Action prose

`### Organize form` has one paragraph, then the launcher block (`... --format text ai-report-organize [--dry-run]`), then four steps:
1. Print the output. Dry-run stops here.
2. On a refusal, report it and stop.
3. Commit as Step 8 says, by exactly the `commit paths:`: `git add -A -- <paths>`, then `git commit -m "Organize ai-reports by month" -- <paths>`. Report the hash, or print the command. I checked this in a scratch repo: it records a rename and leaves other staged files out.
4. Never push.

The form is the only thing that moves reports or rewrites links. Steps 1 to 8 do not apply.

## Tests (each well under 5 s)

**`TestAIReportOrganizeFilesFinishedMonthsAndKeepsLinksOpen`** pins the REQ's proof: moves break links, the second run must change nothing, and dry-run must write nothing.

Fixture (`toolboxNow` frozen at 2026-11-02, committed via `toolboxTestRepository`): `ai-reports/2026-09-01_1200_alpha/index.html` linking `../2026-09-02_1200_beta/index.html#r1`, `../../do-work/queue/REQ-001-sample.md` and `shot.png`; September `beta` and `gamma`; `2026-11-01_0900_current` linking `../2026-09-03_1200_gamma/index.html`; a queue file naming the beta path plus a `-rev1` negative control; a binary file holding the beta path and a NUL.

Steps: (1) dry-run lists 3 moves and 3 edits, and the tree snapshot is unchanged; (2) real run: September is in `2026-09_3-reports/`, current stays at the top; the sibling link and `shot.png` are byte-identical; the escape link becomes `../../../do-work/...`; current links `../2026-09_3-reports/...`; the queue file is rewritten, `-rev1` and the binary are not; every rewritten link resolves; the catalog lists 4 bundles; (3) commit exactly the `commit paths:`, then `git status --porcelain --untracked-files=all` is empty (the list is complete); (4) second run prints `nothing to organize` and the snapshot is unchanged.

**`TestAIReportOrganizeRecountsMonthFoldersAndRefusesUnsafeTargets`** pins a stale count, a plain `YYYY-MM` folder, the singular name, a current-month folder, folder-level mentions, a rename onto an existing path, and an uncommitted edit swept into the commit.

Fixture → expected: `2026-08_1-report/a` plus loose `2026-08-20_1000_b` → `2026-08_2-reports/`; `2026-07/c` → `2026-07_1-report/`; `2026-11/d` → top level, `2026-11` removed; `2026-06-01_1000_e` plus a regular file `2026-06_1-report` → June skipped and listed; `notes.md`, `undated-folder/`, symlink `linked` → listed, untouched; `docs/notes.md` naming `ai-reports/2026-08_1-report/a/index.html` → rewritten.

The test first runs with `docs/notes.md` uncommitted. Organize must refuse, name that file and change nothing. The test then commits it and runs again.

The existing `TestHandlersRegisterCanonicalToolboxCommands` is updated for the new name and count.

**RED** (fails "not registered"): `go test -C skills/do-work/tools/do-work-cli -count=1 -run '^(TestAIReportOrganize|TestHandlersRegisterCanonicalToolboxCommands)' ./internal/toolboxcommands`. **GREEN:** the same run passes. Then run `-run 'AIReport|Handlers'`, `gofmt -l`, `go vet`, `contract-regressions.sh`, `shipped-package-reference-contract.sh` and `staged-skills-contract.sh`.

## Decisions

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

## YAGNI: do not build

`--commit`, gittransaction, rollback, a new result struct; a Just recipe, an exported API, a new package; symlink retargeting, `url()`/`srcset`/script parsing, percent-decoding; nested month folders, a configurable folder format; report-writer changes, the hint (REQ-708); progress output, a parallel scan, extra rewriter unit tests.

Proposed changelog title: "Report Organizer Files Finished Months Into Counted Folders and Fixes Every Link".
