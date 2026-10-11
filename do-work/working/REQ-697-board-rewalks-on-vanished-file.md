---
id: REQ-697
title: 'Board re-walks once when a listed REQ or UR file is gone at read time, instead of dropping it'
status: claimed
route: A
estimate:
  p50_active_minutes: 5
  confidence: high
  basis: [trivial short-circuit]
created_at: 2026-10-11T17:04:32Z
user_request: UR-155
domain: backend
prime_files: ["_dev/primes/prime-kanban-board.md", "skills/do-work-board/tools/queue-kanban/prime-do-kanban.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
write_set: ["skills/do-work-board/tools/queue-kanban/model.go", "skills/do-work-board/tools/queue-kanban/*_test.go"]
claimed_at: 2026-10-11T17:38:52Z
dispatch_at: 2026-10-11T17:41:15Z
builder_handback_at: 2026-10-11T17:45:18Z
integration_at: 2026-10-11T17:45:33Z
commit: 8d3516621ffbc5293bf9c761278b21b8317578a4
review_at: 2026-10-11T17:54:27Z
kb_status: skipped
---
# Board Re-Walks Once When a Listed File Is Gone at Read Time
## What
`buildBoard` lists the `do-work/` tree first (`model.go:452`), then reads each listed file. When a claim or archive renames a REQ between those two steps, the read fails with "no such file" and the loop drops the REQ with a bare `continue` (`model.go:469-472`). For that response the REQ is in no column, a sibling's `depends_on` reports it "not in the do-work tree", and verify reports `checkpoint-names-missing-req` with a remedy that invites deleting a live claim from `CHECKPOINT.md`. Make `buildBoard` list and parse the whole tree once more when any REQ or UR read fails with `fs.ErrNotExist`, and keep the second result.
## Detailed Requirements
- In `buildBoard`, when `parseRequestTicket` (`model.go:469`) or `parseUserRequestTicket` (`model.go:517`) fails with `errors.Is(err, fs.ErrNotExist)`, run `enumerateDoWorkTree` and the parse loops once more from scratch and use that result. One retry only. If the second pass also meets a vanished file, keep the existing skip for it.
- Cover UR files too: UR folders also move from `user-requests/` to `archive/` (`do-work-cli/internal/publication/answer.go:950`), and the UR loop has the same silent skip (`model.go:516-519`).
- Any other read or parse error keeps the existing silent skip. Update the skip comment to say it is for a corrupt or unreadable file, never a file that moved.
- Test seam: a package-level `readTicketFile = os.ReadFile` (or equivalent) that the test swaps. No new `readRequestTickets` helper is required.
- No new `Board` field and no new warning. After the retry the response is correct, so nothing should tell the user something went wrong.
- No change to `serve.go`. Its cache key is per absolute path (`serve.go:416-421`, `:451-462`), so a moved file forces a rebuild on the next request.
- No change to `verify.go` wording (`appendCheckpointGhostFindings`, `verify.go:731-747`).
## Red-Green Proof
**RED prompt/case:** A synthetic tree has `queue/REQ-A.md`, a `queue/REQ-C.md` with `depends_on: [REQ-A]`, and a `CHECKPOINT.md` naming REQ-A. The swapped reader, on its first call for REQ-A's queue path, renames the file into `working/` and returns the real `os.ReadFile` error.
**Why RED now:** REQ-A is missing from `board.RequestsById`. REQ-C carries the "depends on REQ-A, which is not in the do-work tree" warning. `appendCheckpointGhostFindings` reports `checkpoint-names-missing-req` for REQ-A.
**GREEN when:** REQ-A is on the board with `TreeSection == "working"`. No "not in the do-work tree" warning names it. `appendCheckpointGhostFindings` reports nothing for it. A second test keeps a corrupt (non-missing) REQ file unreadable and asserts it is still skipped with no retry. Existing tests still show that a checkpoint id with no file anywhere raises the finding and a dependency on a truly absent id still warns. `go test ./...` under `skills/do-work-board/tools/queue-kanban` passes, and `grep -n 'fs.ErrNotExist' skills/do-work-board/tools/queue-kanban/model.go` finds the retry branch.
**Validation:** Inferred during capture (RED/GREEN taken from the source's test plan, trimmed by the triage)
## Constraints
- No change to `ApplyPlan` or the claim commit (`state_apply.go:59-67`, `:111`). Its order is already right.
- No client polling or server push for the board page (Option B in the source).
- Do not add the `TreeMovedDuringRead` flag, the "do-work/ changed while the board was reading it" warning, or the `serve.go` cache skip. The triage pushed these back (UR-155 Summary, F2 and F4).
- The source's acceptance step of 50 `curl` requests during a live claim is not the proof. The injected-rename test is.
## Builder Guidance
Finding provenance: F1 (High, source: pasted text, upstream suggestion from a consumer repo on suite v0.305.114) and F3 (Medium, same source) of the 2026-10-11 validate-feedback triage. Evidence: `model.go:469-472` and `model.go:516-519` skip any read error; `state_apply.go:59` replaces in place, then `state_apply.go:66` renames. Surface-cost: Earned. Incident: Oct 11, 13:41 UTC, a consumer board response showed a claimed REQ in no column, with a false dependency warning and a false checkpoint finding. Cost: about 10 lines in `buildBoard`, cheaper than a remedy that can delete a live claim. Test: the injected rename above.
## Required Lessons — Dropped for Budget
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` — 9045 tokens, over the 2000 budget and `slugged: partial` so no targeted form; matched because the REQ changes the queue-kanban model.
- `_dev/primes/lessons-kanban-board.md` — 6122 tokens, over budget and `slugged: partial`; matched on queue-kanban parsing.
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 21728 tokens, over budget and `slugged: partial`; matched on the `silent-skip-reads-as-red` family, the same failure shape (a silent skip read as a real absence).
## Full Context
See `do-work/user-requests/UR-155/input.md` for complete verbatim input.
## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read board primes and implementation crews. Inject ticket reads, reproduce the rename race, and retry both parse loops from a fresh enumeration once.
- [x] **[APPLY]:** Added the ticket-reader seam and one bounded whole-board retry in model.go; added model_retry_test.go for REQ/UR moves, permission errors, and retry limit.
- [x] **[UNIFY]:** Reviewed model.go and model_retry_test.go; gofmt and diff checks clean, go vet passed, and go test ./... passed (56.978s total; per-file budget checked by repository gate). No debug artifacts or extra source files.
*Source: F1 and F3 of the validate-feedback triage of the pasted upstream suggestion "the board must not drop a REQ that moves between its walk and its read"*

## Triage

**Route:** A — clear reproduction and explicitly named parse loops; focused board-reader fix.

<!-- D-XX counter: none used. Next decision: D-01. -->

## Plan

Planning not required — Route A. Add an injected-rename regression first, then retry the whole board read once for missing REQ or UR files; preserve all other skips.

## Decisions

- **D-01 — DECIDE & STATE:** Numeric fixture IDs replace symbolic REQ-A/REQ-C because checkpoint extraction accepts numeric request IDs.
- **D-02 — DECIDE & STATE:** One bounded restart label resets the entire board without extracting or reindenting its assembly. The shared flag limits both parse loops collectively to one retry.

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/model.go` (modified)
- `skills/do-work-board/tools/queue-kanban/model_retry_test.go` (new)

**What was done:** The walk-then-read race treated a renamed file as real absence. Both ticket parsers now use an injectable reader; a first-pass fs.ErrNotExist restarts enumeration and all board assembly once. Other errors and second-pass missing files keep the existing skip. Four regression tests cover both moves, permission errors, and the shared retry bound.

## Qualification

Passed — advance returned a satisfied merged_range qualification record for 32e4851f..8d351662. Reviewed both files against every Detailed Requirement and the captured reproduction. QUALIFY-NEW-FILE-UNWIRED is a Go test-discovery exception: model_retry_test.go is discovered by go test and its four tests were actually executed. P-A-U is complete; no extra project paths, debug artifacts, new Board fields, cache changes, or verify wording changes.

## Testing

**Tests run:** Builder: go vet ./... and go test ./... (exit 0; full module 56.978s). Merged tree: focused four-test probe (exit 0, 0.370s) and bash _dev/tests/maintainer-verify.sh (exit 0, 128s wall).
**Result:** All passing. Uncached fast gate ran 427 board tests and 919 core tests; slowest files were 18.85s and 21.20s, both below 30s. advance returned satisfied focused-test and green-gate records at 8d3516621ffbc5293bf9c761278b21b8317578a4.

**Red-green validation:**
- TestBuildBoardRecoversRequestMovedDuringRead: RED missing REQ plus false dependency warning and checkpoint ghost → GREEN restored working REQ and no false findings.
- TestBuildBoardRecoversUserRequestArchivedDuringRead: RED missing UR and stale first-pass REQ title → GREEN archived input plus fresh REQ title.
- TestBuildBoardStopsAfterOneMissingFileRetry: RED one read → GREEN exactly two reads for either ticket kind and retained readable sibling.
- TestBuildBoardSkipsNonMissingReadErrorsWithoutRetry: permission errors pass before and after, skipped with one read.
- RED observed before implementation (exit 1, 0.430s); same four-test command GREEN after (exit 0, 0.399s). Numeric fixture IDs are the captured symbolic-case equivalent required by checkpoint extraction.

**New tests added:** the four named tests in model_retry_test.go. Existing true-absence tests TestDanglingDependencyWarnsAndBlocks and TestVerifyFlagsCheckpointNamingAMissingRequest still pass. No existing tests changed.

**Heavy verification plan:**
- Range: 32e4851f86ec95074277993c9b2db6636c2521e1..8d3516621ffbc5293bf9c761278b21b8317578a4
- queue-kanban-javascript: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — changed board source matched the declared subtree coverage.
- queue-kanban-browser: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — changed board source matched the declared subtree coverage.
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — changed board source matched the declared subtree coverage.

*Verified by work action*

## Review

**Overall: 100%** | 2026-10-11T17:50:04Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 100% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:** None
**Acceptance:** Pass — implementation and merged-tree integration; independent uncached four new regression tests plus true-absence dependency/checkpoint tests passed (0.270s), with the already-green maintainer gate/full module/vet evidence reused.
**Restatement sweep:** nothing redefined — bounded missing-read recovery preserves parser, schema, column, cache, warning, and verify contracts; no inherited wave triggers or stale board/prime restatements found.
**Suggested testing:** 2 items — Deployment unassessed: complete recorded release/package verification and rebuild the consumer board. Live acceptance unassessed: inspect the rebuilt served board after a claim/archive.
**Follow-ups created:** None (0 findings report only)

*Reviewed by review-work action*

## Orientation

The queue board rebuilds once from current paths when a REQ or UR moves during a read, keeping tickets visible and dependency/checkpoint diagnostics accurate.

## Heavy Verification Plan

Base revision: `32e4851f86ec95074277993c9b2db6636c2521e1`
Target revision: `8d3516621ffbc5293bf9c761278b21b8317578a4`

- Lane: `queue-kanban-javascript`
  - Argv: `["env", "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null", "bash", "_dev/tests/maintainer-verify.sh", "--heavy-lane", "queue-kanban-javascript"]`
  - Reasons: skills/do-work-board/tools/queue-kanban/model.go matched subtree skills/do-work-board/tools/queue-kanban; skills/do-work-board/tools/queue-kanban/model_retry_test.go matched subtree skills/do-work-board/tools/queue-kanban
- Lane: `queue-kanban-browser`
  - Argv: `["env", "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null", "bash", "_dev/tests/maintainer-verify.sh", "--heavy-lane", "queue-kanban-browser"]`
  - Reasons: skills/do-work-board/tools/queue-kanban/model.go matched subtree skills/do-work-board/tools/queue-kanban; skills/do-work-board/tools/queue-kanban/model_retry_test.go matched subtree skills/do-work-board/tools/queue-kanban
- Lane: `staged-skills`
  - Argv: `["env", "GIT_CONFIG_NOSYSTEM=1", "GIT_CONFIG_GLOBAL=/dev/null", "bash", "_dev/tests/maintainer-verify.sh", "--heavy-lane", "staged-skills"]`
  - Reasons: skills/do-work-board/tools/queue-kanban/model.go matched subtree skills; skills/do-work-board/tools/queue-kanban/model_retry_test.go matched subtree skills
