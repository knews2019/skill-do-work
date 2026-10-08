# Builder brief — REQ-645 (recover --take-over keeps a 1-3 space indented generated-name heading; only a column-0 section is removable)

- Worktree (your only writable tree): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-645-recovery-removal-requires-column-zero-heading
- Branch: worktree-agent-REQ-645-recovery-removal-requires-column-zero-heading, created from main HEAD with `git worktree add -b` (commit here; never merge, never rebase, never push, never check out another branch).
- REQ (read-only, main tree): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/working/REQ-645-recovery-removal-requires-column-zero-heading.md. Read it fully: What, Why, Verified Facts, Detailed Requirements 1-7, Constraints, Builder Guidance, Red-Green Proof. Requirements 6 and 7's release and index refresh belong to the integrator, not you; requirement 7's lesson bullet is yours to PROPOSE in the hand-back.
- UR (read-only): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/user-requests/UR-141/input.md.
- Hand-back file (the ONE main-tree path you may write; never stage or commit it): /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-10-08-141302/REQ-645-handback.md
- Route A, tdd: true, impact-critical (data loss in a recovery path), effort-mechanical, domain backend. Go only.

## Rules to load first
/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/crew-members/general.md, coding-guardrails.md, shared-principles.md, communication-style.md, backend.md (if present), testing.md (tdd: true) — same directory. Primes: /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/tools/do-work-cli/prime-do-work-cli.md (read before editing) and /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/_dev/primes/prime-releases.md (the release step owns VERSION/CHANGELOG, not you). Lessons: in /Users/t2/Desktop/e1-experimental-repos/skill-do-work2/skills/do-work/tools/do-work-cli/lessons-do-work-cli.md read the bullets carrying `[family: closed-enumeration-for-a-condition]` (the 0.305.70 REQ-635 entry is this bug's parent) and `[family: rule-direction-checked-against-callers]` (the shape of this bug: the boundary rule was pinned, the removal rule was not).

## TDD order (mandatory; the integrator verifies red-then-green evidence)
1. In `internal/requeststate/recovery_markdown_test.go`, next to `TestRecoveryKeepsUserTextUnderAFourSpaceIndentedHeading`, add a test for the REQ's Red-Green Proof body: a REQ document whose body is `# Request\n- Example:\n  ## Plan\n  user sample\n\nMUST keep.\n## Timing\ngenerated summary\n` (both `\n` and `\r\n`, like the siblings). Expect `stripGeneratedRecoverySections` to remove only the column-0 `## Timing` section and keep the bullet, `  ## Plan`, `  user sample` and `MUST keep.` byte for byte. Run it, confirm it FAILS (today the indented Plan is deleted), record the exact failure text.
2. Implement; run again (GREEN). Record test name, failure-before, pass-after in the hand-back.

## The change (decided at capture, do not reopen)
- `internal/requeststate/state_apply.go` `stripGeneratedRecoverySections` (line ~983): a section is removable only when its heading line starts at column 0. Carry the indent the way that reads best: a field on `requestmodel.VisibleSection` set by `VisibleSections` in `internal/requestmodel/visible_sections.go` (preferred; it is the parser that already measures the indent via `blockIndentContent`), or a byte check at `body[section.Start]`. Boundary recognition (0-3 spaces) stays exactly as it is: `TestRecoveryPreservesIndentedAndCommentedRequirements` and `TestVisibleSectionsApplyTheZeroToThreeSpaceIndentRule` must stay green and untouched.
- Requirement 4: read the other two callers that remove or replace a section by name, `replaceTimingSection` in `internal/lifecycletiming/lifecycle_timing.go` and `advanceSections` in `internal/lifecycleadvance/advance_commands.go`. If either would overwrite or reject an indented user sample (for example a 2-space `## Timing` under a bullet), apply the same column-0 requirement there with one test each; if already safe, say exactly why in the hand-back. Do not widen beyond that.
- No list-marker or container parsing anywhere. Indent at the heading line is the whole signal. Keep doc comments true (the `VisibleSections` comment and the `stripGeneratedRecoverySections` neighbourhood).
- Do NOT touch CHANGELOG.md, VERSION, any mirror, `do-work/` (except the hand-back file), `_dev/`, or `lessons-do-work-cli.md` (propose the bullet; the integrator writes it).

## Write boundary
The REQ's `write_set` (`state_apply.go`, `visible_sections.go`, `recovery_markdown_test.go`) plus `visible_sections_test.go` if the struct gains a field, plus the two requirement-4 files and their tests only if requirement 4 needs a change. Anything else: stop and say so in the hand-back.

## Verify before hand-back (from the worktree root, wall times recorded)
- `cd skills/do-work/tools/do-work-cli && go build ./... && go vet ./... && go test ./... -count=1` (whole module; ~95 s).
- `gofmt -l .` in that directory prints nothing.
- `bash do-work/runs/work-2026-10-08-141302/REQ-645-probe.sh` does not exist in your worktree; run its body instead: `go test ./internal/requeststate/ ./internal/requestmodel/ -run 'Recovery|VisibleSections|StripGenerated' -count=1` from the module directory (must stay well under 30 s).
- `git diff --stat` and `git diff --check`.

## Hand-back (write the file at the absolute path above)
- Branch name, worktree path and every commit hash (a single commit preferred). Base: b38928de.
- File manifest: each file with (modified)/(new) and one line on what changed.
- P-A-U text for [PLAN], [APPLY], [UNIFY]; [UNIFY] lists git diff --stat, gofmt/vet, and each file checked.
- Red-green evidence: test name, exact failure before, pass after.
- Requirement 4 verdict per caller with the line you read.
- `## Decisions` (D-01 onwards; DECIDE & STATE or ESCALATE with Value and Risk).
- `## Discovered Tasks` (out-of-scope finds, each ending `→ report only` unless impact-critical; do not fix them).
- Lessons read (satellite + family).
- A proposed lesson bullet for lessons-do-work-cli.md in the file's top-entry shape (`- [family: <slug>] <version> (REQ-645, 2026-10-08): **bold one-liner.** body`; pick `rule-direction-checked-against-callers` or `paired-predicate-drift` with a one-line reason), and a proposed CHANGELOG entry (descriptive title saying what shipped, plain words, names this as a data-loss fix in recover --take-over, the residual of 0.305.70). The integrator writes both.
- Integration seams (none expected). Test wall times.
