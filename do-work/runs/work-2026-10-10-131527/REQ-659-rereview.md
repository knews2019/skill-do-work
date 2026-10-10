## Delta re-review: REQ-659 (do-work-cli `frontmatter set` and `req append-section`)

**Approve.** The delta closes F1, F4 and F5. The F1 check refuses only real level-two headings and lets subheadings, fenced code and comments through. One new finding (N1) sits next to F1 but is older than this delta.
Delta `a56cd9c4..8a33190c` (3 files: `internal/corehelpers/request_writers.go`, `internal/corehelpers/request_writers_test.go`, `skills/do-work/actions/work-reference.md`).

### Prior findings

- F1 (`--from` body with a `## ` heading): **closed.** `request_writers.go:101-106` trims the body and refuses `SECTION-BODY-HAS-HEADING` when `requestmodel.VisibleSections` finds any heading. This runs before `resolveActiveRequest`, so no file is read or written. The refusal names the `--from` path, and its reason ("remove its ## heading") replaces the misleading "would not be a visible section" reason for the Testing Section Template case. `work.md:21` already says "a file holding only its body", so the prose and the code now agree.
- F4 (Timestamp rule "only place … spells a command"): **closed.** `work-reference.md:61` now says "for obtaining one to hold". The parenthetical names the stamp writer under **Stamps are append-only** (`:73`). A grep of `actions/` finds only `:66` (`date -u`, inside the rule) and `:73` (`--at now`), so the sentence is true.
- F5 ("keep their owning commands"): **closed.** `:73` now says "because their owning command or the hand write at their defining site keeps them". This covers `abandon.md:62` and `clarify.md:138`, `:226`. The sentence reads correctly in context.

### F1 false-positive check (real binary built from `8a33190c`, fixture REQ with `## What` and `## Review`, `--section Testing`)

| `--from` body | Result |
|---|---|
| `### Sub` subheading | success, inserted before `## Review` |
| backtick fence holding `## Review` | success |
| tilde fence holding `## Testing` | success |
| block HTML comment holding `## Review` | success |
| inline `<!-- ## Review -->` | success |
| 4-space indented `## Review` (code block) | success |
| `# Title` (level one), `##Review` (no space), bare `## ` | success |
| `## Testing` (template verbatim) | refused `SECTION-BODY-HAS-HEADING`, bytes unchanged |
| text then `## Review` | refused `SECTION-BODY-HAS-HEADING` |
| `## Review` with CRLF endings | refused `SECTION-BODY-HAS-HEADING` |
| 2-space indented `## Review` | refused `SECTION-BODY-HAS-HEADING` |

The 2-space case is refused even though section matching counts only column-0 headings. That is correct: CommonMark renders a heading indented by 1 to 3 spaces as a heading, so it does not belong in a section body.

### Test adequacy

- The new assertion in `TestRequestAppendSectionRepeatIsNoOpAndConflictRefuses` (`request_writers_test.go:208-211`) fails when the check is disabled. In a scratch worktree I changed the condition to `if false && …`. The test then failed with `request_writers_test.go:211: outcome=success findings=[], want refused (exit 1)`. With the check restored it passes.
- `go test -count=1 ./internal/corehelpers/`: ok (10.1 s) at `8a33190c`.
- `bash _dev/tests/contracts/core-checks.sh` in ROOT: exit 0, "core-checks contract probes passed." (5.0 s).

### New findings

- N1. `req append-section` accepts a body that leaves a fence or HTML comment open. The open fence or comment then hides every section after the insertion point (`request_writers.go:161-171`). The re-check counts only visible copies of the inserted section. It does not check that the sections already in the file stay visible. Reproduced: a body of ` ``` ` + `## x` (or `text <!-- open`) with `--section Testing` exited 0. The `## Review` that followed became part of the open fence or comment, so the file no longer has a visible `## Review`. `advance` then reads the REQ as missing that phase. This is older than the delta. It is listed here because it has the same cause as F1 (the body's shape is not checked). Fix (about 4 lines): compare the column-0 visible section count before and after the insertion and refuse `SECTION-WRITE-FAILED` unless it grew by exactly one. That also covers the end-of-file case the current re-check handles. impact-user-visible → report only

### Scores

**Overall: 94%** (up from 91%)

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 95% | Unchanged. Release is still pending at finalization |
| Code Quality | 92% | F1 closed with a small, well-placed check. N1 is open |
| Test Adequacy | 90% | The new assertion pins F1 (checked by mutation). There is still no test for the visible-copy re-check |
| Scope | 95% | 3 files, only the three findings |
| Risk | Low | Refusals never write. N1 needs a malformed body |
| Acceptance | Pass | Re-checked on the real binary |

**Acceptance:** Pass.
**Follow-ups created:** None (N1 report only).
**Scratch:** the worktree `scratchpad/int-659/rereview/rr` was removed with `git worktree remove`. The fixture and binary are left in `scratchpad/int-659/rereview/`.

*Delta re-review by review-work action*
