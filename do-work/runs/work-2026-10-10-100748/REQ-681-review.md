# Review: REQ-681 (appendSectionEntry finds its section with VisibleSections)

**Approve** — the fix does what the REQ asks, the three RED cases are real, and the LF outputs for "missing", "last" and "middle" are byte-identical to the old function.
Route A | merge range `0362663f..61a42672` (builder commit `8994fd85`), orchestrated mode, worktree dispatch.

## What's built
- `appendSectionEntry` (`skills/do-work/tools/do-work-cli/internal/requeststate/state_apply.go:1024`) now loops over `requestmodel.VisibleSections` and skips headings with `HeadingIndent != 0`. `sectionLineBounds` is deleted.
- A CRLF file, a heading with trailing spaces, and a `## Blocked` inside a code fence no longer produce a duplicate heading or an entry inside the fence.
- New `append_section_entry_test.go` pins those three cases.

## Decisions / risks for you
- None that need a decision. The one behaviour change for LF input that is worse than before (F1) needs a line that starts with three or more of a punctuation mark such as `...`, and no REQ file in `do-work/` has one today.

## Findings

**Important:** None.

**Minor:**
- F1. `state_apply.go:1026`: the new lookup inherits `VisibleSections`' dialect-fence rule, so for LF input a body line that starts with `...`, `"""`, `+++`, `:::` (any enclosing punctuation run of 3 or more) and is never closed hides every later heading. The old exact-line match found `## Blocked` there; the new code adds a new `## Blocked` on every unblock (reproduced in a scratchpad harness: 3 headings after 2 appends). It is the same rule `markdownSectionBytes` and the Timing writer already follow, and a scan of `do-work/**/*.md` found such lines only in inputs, hand-backs and the inbox, never in a REQ file. — impact-negligible → report only
- F2. `state_apply.go:1031`: the builder-reported edge (a last section ending in an unclosed fence or comment has `End < len(text)` and takes the "middle" path) has no test. My harness shows the new output puts the entry before the hidden region, while the old code wrote it inside the unclosed fence, so the change is an improvement. Already recorded as a discovered task. — impact-negligible → report only

**Nit:**
- F3. `state_apply.go:885` (checkpoint caller): this fallback runs only when `repositorymodel.CheckpointClaimBounds` (exact, CRLF-aware, fence-blind) finds no canonical heading. If the file has `## In Progress (interrupted)` with trailing spaces, the new code now inserts into that heading instead of adding a canonical one. Discovery (`projectCheckpointClaims`) and removal (`filterCheckpointClaims`) both scan the whole file when no canonical heading exists, so the entry is still found and removed. No defect, just two section readers with different dialects in one flow. — impact-negligible → report only

**Anti-bloat count (requested):** new helpers 0, new options/flags/constants 0, new files 1 (`append_section_entry_test.go`, named by the REQ), test cases 3 (2 top-level tests, 2 subtests in the first), matching Detailed Requirement 2. The reference patch's LF control subtest was dropped (D-01). No decorative test. One comment line added, explaining the `HeadingIndent` skip. Net `state_apply.go`: +10 −24.

**Restatement sweep:** redefined how `appendSectionEntry` finds its section (shared visible-section reader instead of an exact-line `## <name>` match). Greps run over `skills/`, `_dev/`, `docs/` and the lessons files (`skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`, `_dev/primes/lessons-*.md`): `appendSectionEntry`, `sectionLineBounds`, "second `## Blocked`", "duplicate heading", "exact-line"/"exact line", `CRLF` in actions/crew/docs/primes, and `## Blocked`/`## Cancelled`/`In Progress (interrupted)` lines that mention append/heading/match. No stale restatement found. `actions/work-reference.md:312` ("Keep the exact `## In Progress (interrupted)` heading as the claim-evidence boundary") describes `CheckpointClaimBounds`, which this diff does not change, so it still agrees.

## Requirements Checklist
- [x] DR1: `appendSectionEntry` uses `VisibleSections` with the `HeadingIndent != 0` skip; `sectionLineBounds` deleted — delivered (`grep -rn sectionLineBounds skills _dev`: no match, exit 1).
- [x] DR2: the three RED cases and nothing more — delivered (CRLF, trailing spaces, fenced heading).
- [x] DR3: no line-ending matching added — delivered (inserted text stays LF, D-03).
- [x] Constraint: smallest change, no new helpers/flags/options/files beyond the named test — delivered (count above).
- [x] Constraint: tests pin only the named failures — delivered.
- [x] Constraint: adjacent issues recorded, not fixed — delivered (unclosed-region edge is a discovered task).
- [x] Constraint: REQ-659 untouched; no consumer copies edited — delivered (diff touches only the two declared files).
- [ ] Constraint: release per `prime-releases.md` — not in this merge range; the hand-back carries a proposed CHANGELOG entry, so the release belongs to finalization. Not scored as missing.
- [x] UR-152 F7: CRLF, trailing-space and fenced-heading failures fixed; mixed line endings left alone — delivered.

## Acceptance Testing

**Result: Pass** (stages covered: implementation and integration)
- `go test -count=1 ./internal/requeststate/ -run AppendSectionEntry -v`: all 3 cases PASS. `go vet ./internal/requeststate/` clean, `gofmt -l internal/requeststate` empty.
- RED confirmed independently: the three assertions run against a copy of the old `appendSectionEntry`/`sectionLineBounds` fail (CRLF 2 headings, trailing 2 headings, fence 1 heading), and pass against the new code.
- Behaviour preservation (LF): a scratchpad harness ran old and new side by side on 18 documents. Byte-identical for section missing, section last, section last with extra trailing newlines, section in the middle (with and without blank lines), empty section last/middle, `Cancelled` in the middle, checkpoint missing/last (`In Progress (interrupted)`), closed fence inside the section, two `## Blocked` headings (first wins in both), `###` inside the section, and empty input. Differences only where the old code was wrong (a fenced or commented `## X` inside the section, an unclosed fence at the end: old wrote the entry inside the fence/comment) or F1.
- The full repository gate passed on the merge per the REQ's `## Testing`; not rerun.

## Suggested Additional Testing
- Deployment: unassessed. A release has not been cut yet; confirm the finalization release carries this fix.
- Live acceptance: unassessed. On a Windows or `core.autocrlf=true` checkout, block and unblock a REQ twice and confirm one `## Blocked` heading.

## Scores (on the record — not the headline)

**Overall: 97%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | All three DRs and constraints met; release is finalization's job |
| Code Quality | 95% | Same shape as `markdownSectionBytes`; F1 is an inherited dialect trade-off |
| Test Adequacy | 95% | Three focused RED cases, red-green confirmed; F2 edge untested by design |
| Scope | 100% | Exactly the two declared files, net deletion |
| Risk | Low | Three callers unchanged; only rare dialect-fence input behaves differently |
| Acceptance | Pass | Implementation and integration stages |

## Follow-ups created
None (3 findings report only)

## Review

**Overall: 97%** | <TIMESTAMP>

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 95% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:**
- F1. `state_apply.go:1026`: an unclosed dialect-fence line (for example one starting `...`) before the section now hides `## Blocked`, so each unblock adds another heading where the old exact-line match found it; same rule as `markdownSectionBytes`, no REQ file in `do-work/` has such a line. — impact-negligible → report only
- F2. `state_apply.go:1031`: a last section ending in an unclosed fence or comment takes the "middle" path (entry goes before the hidden region, better than the old in-fence write) and is untested; already a discovered task. — impact-negligible → report only
- F3 (nit). `state_apply.go:885`: the checkpoint fallback can now adopt a non-canonical `## In Progress (interrupted)  ` heading instead of adding a canonical one; discovery and removal scan the whole file in that case, so claims stay found. — impact-negligible → report only

**Acceptance:** Pass — implementation and integration: 3 RED cases pass (RED reconfirmed against the old code), vet and gofmt clean, 18-document old-vs-new harness byte-identical for every LF missing/last/middle case.
**Restatement sweep:** redefined `appendSectionEntry` section lookup (visible-section reader instead of exact-line `## <name>` match); greps for `appendSectionEntry`, `sectionLineBounds`, second/duplicate heading, exact-line, CRLF and Blocked/Cancelled/In Progress append prose over skills/, _dev/, docs/ and lessons files found no stale restatement.
**Suggested testing:** 2 items
**Follow-ups created:** None (3 findings report only)

*Reviewed by review-work action*
