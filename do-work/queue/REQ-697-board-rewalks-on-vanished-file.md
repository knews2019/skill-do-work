---
id: REQ-697
title: 'Board re-walks once when a listed REQ or UR file is gone at read time, instead of dropping it'
status: pending
created_at: 2026-10-11T17:04:32Z
user_request: UR-155
domain: backend
prime_files: ["_dev/primes/prime-kanban-board.md", "skills/do-work-board/tools/queue-kanban/prime-do-kanban.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
write_set: ["skills/do-work-board/tools/queue-kanban/model.go", "skills/do-work-board/tools/queue-kanban/*_test.go"]
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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: F1 and F3 of the validate-feedback triage of the pasted upstream suggestion "the board must not drop a REQ that moves between its walk and its read"*
