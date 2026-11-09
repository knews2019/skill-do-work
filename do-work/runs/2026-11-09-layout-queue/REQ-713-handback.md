# REQ-713 Hand-back

Branch `worktree-agent-REQ-713-folder-structure-check-default-param-recipes`: RED 06761c55, GREEN (tip of the branch).

## Files
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/folder_structure.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/toolboxcommands/folder_structure_test.go` (modified)

## What was done
`justRecipeHeader` (a regex that stopped at the first `=`) is replaced by `justRecipeName`: it reads the leading name, then scans the rest of the line outside quotes for the first colon. A colon followed by `=` is a variable assignment; any other colon makes the line a recipe header. Directives (set, alias, export, import, mod) are still excluded.

## Red-green evidence
- RED at 06761c55: `TestFolderStructureCheckReportsEachGapClass` fails, the wanted `FOLDER-STRUCTURE-JUST-LOGIC justfile: run-kanban` is missing. The fixture recipe is `run-kanban $port="8090" url='http://localhost' +rest: lint` with three body lines.
- GREEN at the branch tip: `go test ./internal/toolboxcommands/ -count=1` ok in 22s.
- No false positive: the passing-layout test now also holds `url:="http://localhost"` (assignment with no spaces and a colon inside the value) and a one-line recipe `serve $port="8090" *args:` under a `[positional-arguments]` attribute, and still reports zero gaps.
- On this repository the check now reports `justfile: run-kanban` (3 body lines) beside `justfile: kanban-static`.

## Decisions
- D-01 (DECIDE & STATE): A small quote-aware scan replaces the regex. A regex cannot skip a colon inside a quoted default (`url='http://localhost'`) without lookarounds, which Go's regexp does not have.

## Discovered Tasks
- None.

## Lessons read
- No `required_lessons` entry on this REQ. The matching satellite `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` was dropped for budget at capture (24067 tokens).

## Addendum (review F1)

The reviewer found that `internal/managedsection/just_definitions.go` already ships `justDefinitionName`, a state machine over Just's five literal forms. The first fix duplicated a weaker version of it. Third commit on the branch: `JustLineDefinitionName` (an exported one-line wrapper) is added there, and `folder_structure.go` calls it. The new `justRecipeName`, `justLeadingName` and the `justDirective` list are deleted; the scanner already returns nothing for `set`, `export`, `import` and `mod` lines and for assignments. An `alias` line returns the alias name, which has no body lines and so can never be a gap.

- Files: `skills/do-work/tools/do-work-cli/internal/managedsection/just_definitions.go` (modified), `skills/do-work/tools/do-work-cli/internal/toolboxcommands/folder_structure.go` (modified).
- D-02 (DECIDE & STATE): Reuse over a second scanner, although the REQ names only `folder_structure.go`. One scanner means one place to fix Just syntax.
- Evidence: `gofmt -l` and `go vet` clean on both packages; `go test -run FolderStructure ./internal/toolboxcommands/` ok.
