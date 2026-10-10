# REQ-681 hand-back

- Branch: worktree-agent-REQ-681-append-section-entry-visible-sections
- Worktree: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-681-append-section-entry-visible-sections
- Base: b629e5cd. Commit: 8994fd85 (one commit, subject starts `[REQ-681]`).

## File manifest
- skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go (modified): `appendSectionEntry` loops over `requestmodel.VisibleSections` (skips `HeadingIndent != 0`); `sectionLineBounds` deleted. Net -14 lines in this file.
- skills/do-work/tools/do-work-cli/internal/requeststate/append_section_entry_test.go (new): 3 RED cases.

## P-A-U
- [PLAN]: Write the three failing cases, then swap the section lookup in `appendSectionEntry` for `VisibleSections` using the reference patch's shape, keeping the "section is last", "section in the middle" and "section missing" outputs byte-identical to before for LF input.
- [APPLY]: Test file first (RED at base), then the function rewrite and the deletion. The reference patch's extra LF control subtest was left out so the file holds exactly the three named cases.
- [UNIFY]: Files checked: state_apply.go (diff read, no debug artifacts; callers at :622, :679, :885 unchanged, signature unchanged), append_section_entry_test.go.
  - `git diff b629e5cd --stat`: append_section_entry_test.go 40 insertions; state_apply.go 10 insertions, 24 deletions (34 changed lines); total 2 files, +50 -24.
  - `go test -count=1 ./internal/requeststate/`: exit 0, 6.9 s (7.1 s wall).
  - `grep -rn sectionLineBounds skills/do-work/tools/do-work-cli`: no match (grep exit 1).
  - REQ-681-probe.sh from worktree root: exit 0.
  - `gofmt -l internal/requeststate`: empty, exit 0. `go vet ./...`: exit 0, 1.0 s.
  - `git diff --cached --check`: exit 0.

## Red-green record
Base (b629e5cd), `go test -run AppendSectionEntry`:
- TestAppendSectionEntryReusesTheExistingHeading/CRLF_file: FAIL, "want exactly one Blocked heading, got 2" (new `\n## Blocked` appended after the CRLF document).
- .../heading_with_trailing_spaces: FAIL, same message, got 2.
- TestAppendSectionEntryIgnoresAHeadingInsideACodeFence: FAIL, entry was added under the fenced `## Blocked` and no real heading was created.
After the change: all 4 PASS lines (2 top-level, 2 subtests) green.

## Decisions
- D-01 DECIDE & STATE: dropped the patch's "LF" control subtest. It does not fail at the base, and the REQ and brief say exactly three RED cases.
- D-02 DECIDE & STATE: kept the patch's `HeadingIndent != 0` skip with its one-line comment, because the REQ names it (Detailed Requirements 1).
- D-03 DECIDE & STATE: mixed line endings left out of scope as decided. Behavior: for a CRLF file the inserted entry and blank lines use LF.

## Discovered Tasks
- The "section is last" check `visible.End >= len(text)` treats a section ending in an unclosed fence/comment (End < len(text)) as a middle section; the entry then goes before the hidden region and the region is kept as "following text". That is safe (nothing is lost) but untested. → report only

## Lessons read
`_dev/primes/lessons-releases.md` (whole file). Skipped the Go satellite `lessons-do-work-cli.md` per brief; no bullet it names was touched.

## Anti-bloat check
`git diff --stat` above. Added that the REQ did not name: none. (Test file and the rewritten function are named by the REQ; no new helper, flag, constant or option.)

## Proposed CHANGELOG entry
Title: Appending a section entry now finds the section with the shared visible-section reader

Fixes a duplicate heading when the request file has CRLF line endings, a heading with trailing spaces, or an example `## Blocked` inside a code fence. The block, cancel and interrupted-claim writers all use this function.
- `appendSectionEntry` uses `requestmodel.VisibleSections`, the same reader `markdownSectionBytes` uses.
- `sectionLineBounds` is deleted (net -14 lines in `state_apply.go`).
- New `append_section_entry_test.go` pins the CRLF, trailing-space and fenced-heading cases.

## Proposed lesson bullet
none (a hand-rolled line scanner next to a shared reader is the whole story; the code comment on the function already says why).

## Integration seams and test wall times
No seam. Callers `:622`, `:679`, `:885` unchanged. Full package: 6.9 s (single run, no retry needed). Probe exit 0.
