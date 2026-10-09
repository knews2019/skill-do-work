---
id: REQ-653
title: '[impact-rule-change] Fan-out orchestration prose is compressed in place and moved to a companion reference'
status: claimed
route: B
estimate:
  p50_active_minutes: 40
  confidence: medium
  basis:
  - Route B
  - 18-file write set
  - 2 subsystems involved
  - 7 acceptance criteria
  calculated_at: 2026-10-09T16:31:42Z
created_at: 2026-10-09T16:06:43Z
user_request: UR-143
domain: general
prime_files: ["_dev/primes/prime-action-files.md", "_dev/primes/prime-releases.md"]
tdd: false
maintenance: true
impact: impact-rule-change
effort_estimate: effort-substantive
related: [REQ-650, REQ-651, REQ-652]
batch: october-review-triage
depends_on: [REQ-652]
write_set: ["skills/do-work/actions/fan-out-reference.md", "skills/do-work/actions/work.md", "skills/do-work/actions/work-reference.md", "skills/do-work/actions/review-work.md", "skills/do-work/actions/cleanup.md", "skills/do-work/actions/restart-with-parallel-handoff.md", "skills/do-work/actions/capture.md", "skills/do-work/actions/capture-reference.md", "skills/do-work/crew-members/background-agents.md", "skills/do-work-knowledge/crew-members/background-agents.md", "skills/do-work-toolbox/crew-members/background-agents.md", "skills/do-work/docs/work-guide.md", "skills/do-work-board/actions/board.md", "skills/do-work-board/docs/board-guide.md", "skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md", "skills/do-work-board/tools/queue-kanban/model.go", "skills/do-work-board/tools/queue-kanban/verify.go", "decisions/log.md"]
claimed_at: 2026-10-09T16:31:01Z
dispatch_at: 2026-10-09T16:34:28Z
---
# Fan-Out Orchestration Prose Is Compressed in Place and Moved to a Companion Reference
## What
Create `skills/do-work/actions/fan-out-reference.md` as the single home for coordinator-shape detail, with no size limit, and move there: Worktree Dispatch Mode (`actions/work-reference.md:394-488`, including "Delegated integration — the coordinator shape" at `:470` and the `REQ-NNN-integrate.md` run-directory row at `:479`), the Mid-Run Messages routing detail (`actions/work.md:165-174`), the Restatement sweep procedure and wave-end set-aside handling (`work.md:368-370` and the fan-out sentences from REQ-641 and REQ-642), and the landed hand-back consumption and dispatch-instant bookkeeping (`work.md:288` and `:296`). In `work.md` and `work-reference.md` each moved passage becomes one or two sentences stating the rule plus a pointer to the section in the new file. Rewrite every moved passage delete-before-add: keep the rule, the named commands and the one-writer invariant, drop restated context.
## Why
Consumer review (validate-feedback 2026-10-09, F6, accepted in part): "they added a lot of dense, wordy instructions for the Orchestrating Model … written as long prose in files that are already huge". Measured: REQ-639 to REQ-642 and REQ-648 added about 1,115 words to `work.md` (12,711 to 13,826) and 646 to `work-reference.md` (22,109 to 22,755); each paragraph is one physical line of up to 2,280 characters. No size budget exists anywhere; ADR-001 (`decisions/records/adr-001-modular-action-prompts-and-companion-references.md`) records `work.md` outgrowing a budget twice before and being split into a companion each time. The maintainer chose: "compress in place and create structured companion reference where things remain clear without any size limitation."
## Verified Facts (from triage)
- Additions in the window, by place: `work-reference.md:470` Delegated integration (~350 words, REQ-639 plus a set-aside clause from REQ-642); `:479` run-directory row; the wave-end semantic-collision sentence in the fan-out section (REQ-641); the Decision Brief emphasis-note sentence (REQ-640); `:871` HANDLED run-level-call clause (REQ-642); `work.md:165` `### Mid-Run Messages (any step)`, the only new heading, ~575 words (REQ-640); `work.md:288` "Landed hand-back — consume it, never re-dispatch" (~149 words, REQ-639); `:296` dispatch-instant paragraph (REQ-639); `:368` "Restatement sweep (MUST)" (REQ-641); `:523` operator row (REQ-648, now revised by REQ-652).
- None of the additions is a shell sequence written as prose; the closest, `work.md:288` and `:296`, are conditional checks around the `record-timing-event` CLI verb, which already owns the arithmetic. Judgment stays prose; no script is extracted.
- Existing section map: `work-reference.md` has 40 `##` sections; Worktree Dispatch Mode is `:394-488`. `work.md` keeps its ten-step orchestrator shape.
- `_dev/tests/shipped-package-reference-contract.sh` validates same-package and cross-package citations; `_dev/tests/contract-regressions.sh` enforces router word budgets on `SKILL.md`, which this REQ does not touch.
## Detailed Requirements
1. New file `skills/do-work/actions/fan-out-reference.md` with a description blockquote that says why it belongs in core (it completes `actions/work.md` for fan-out runs), per `_dev/primes/prime-action-files.md`. Sections, in this order: Worktree Dispatch Mode (moved whole, including the Isolation ladder and Delegated integration), Landed hand-back and dispatch timing, Mid-Run Messages, Restatement sweep and wave-end set-aside. No size limit: clarity wins over brevity inside this file, but every moved passage is still rewritten delete-before-add.
2. `actions/work.md`: Mid-Run Messages shrinks to its routing rule (what changes decides where a message goes, and the run never stops for it) plus a pointer; the landed hand-back and dispatch-instant paragraphs shrink to their rule plus a pointer; the Restatement sweep keeps its MUST line and points at the procedure.
3. `actions/work-reference.md`: Worktree Dispatch Mode becomes a short section that states the every-run-mode rule, the one-writer invariant and the integrator's entry (`advance REQ-NNN`, never `recover`, `--take-over` or `--assume-sole-authority`), and points at the new file. The run-directory table keeps its row.
4. Repoint every citation of the moved headings: `grep -rn "Worktree Dispatch Mode\|Delegated integration\|Mid-Run Messages\|Restatement sweep\|Isolation ladder\|Fan-Out Dispatch" skills/ _dev/ decisions/ CLAUDE.md` finds only the new file, the pointers, and historical records (changelog, archived REQs, ADRs are never rewritten).
5. Word counts after the change: `work.md` below 12,711 and `work-reference.md` below 22,162 (their 0.305.73 values), measured with `wc -w`.
6. Append to `decisions/log.md` under today's date: ADR-001 reaffirmed a third time; companion split chosen with an explicit no-size-limit for the companion; the two files that regrew and the review that prompted it.
7. `bash _dev/tests/shipped-package-reference-contract.sh` and `bash _dev/tests/contract-regressions.sh` green. Release per `_dev/primes/prime-releases.md`.
## Constraints
- `maintenance: true`: load `crew-members/maintenance.md`; delete before you add. No rule's stated condition changes; this moves and compresses text.
- Depends on REQ-652 so the operator sentences it adds to `work.md:523` and `work-reference.md:788` are compressed with the rest, not merged around.
- Do not touch `SKILL.md`, the router, or any crew member. Do not extract a script.
- Read `_dev/primes/prime-action-files.md` before editing; cross-references use that prime's spelling, and a new action-adjacent file updates the owning `SKILL.md` only if it is routed, which this reference is not.
## Dependencies
REQ-652 must be completed first (same two files).
## Builder Guidance
Certainty is high on what moves (the triage listed each passage by line). Latitude: section titles inside the new file, how much each pointer sentence keeps, and whether the Decision Brief emphasis-note sentence stays in `work-reference.md` (it is small and belongs to the brief, so it may stay).
## Red-Green Proof
**RED prompt/case:** `ls skills/do-work/actions/fan-out-reference.md` and `wc -w skills/do-work/actions/work.md skills/do-work/actions/work-reference.md`.
**Why RED now:** The file does not exist; the counts are 13,826 and 22,755; the fan-out detail is spread over nine passages in the two files.
**GREEN when:** The file exists with the four sections; the counts are below 12,711 and 22,162; the citation grep in requirement 4 finds only the new file, pointers and historical records; both contract tests are green; `decisions/log.md` has the entry.
**Validation:** User adjusted. The maintainer chose both compress-in-place and a structured companion reference with no size limit, replacing the four options offered.
## Required Lessons — Dropped for Budget
- `_dev/primes/lessons-action-files.md` as a whole satellite (6855 tokens, over the 2000 budget; `slugged: partial`, so no targeted form). Matching reason: its owning prime governs the action files this REQ edits; families `condition-preserving-prose-extraction` and `restated-mechanism-unchecked` are exactly the risks of moving and compressing rule text.
- `_dev/primes/lessons-shell-commands.md` family `moved-prose-loses-its-anchor` (8480 tokens, over budget; `slugged: partial`, so no targeted form). Matching reason: moved prose loses the anchors other files cite; requirement 4 is the mitigation.
## Full Context
See `do-work/user-requests/UR-143/input.md` for complete verbatim input. No queued candidate in any UR shares this root cause (the queue was empty at capture).
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: consumer review "they added a lot of dense, wordy instructions for the Orchestrating Model … in files that are already huge", accepted in part by `do-work-toolbox validate-feedback` on 2026-10-09 as F6; maintainer answer "compress in place and create structured companion reference where things remain clear without any size limitation".*

---

## Triage

**Route: B** - Medium

**Reasoning:** The outcome is clear (a new companion reference, two files compressed) but every citation of the moved headings across four packages and the tests that pin them had to be discovered first; an Explore agent mapped them, and a Scope declaration keeps the eighteen-file write set honest.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

Citation map for the passages that move (Explore agent, 2026-10-09, read-only; full table in the run directory's `REQ-653-exploration.md`):

- **Moved ranges:** `work-reference.md` 394-488 (Worktree Dispatch Mode: Isolation ladder, Naming, operative name, State stays home, Sole integrator, Merge never rebase, When to merge steps 0-4, Hold both endpoints, cumulative range, Post-merge verification, Cleanup happy and crash paths, Fan-Out Dispatch, Auto-wave, Serial-only, Delegated integration, the run-directory table, brief and hand-back path rule, dispatch mechanism); `work.md` 165-174 (Mid-Run Messages), 288 (Landed hand-back), 296 (dispatch instant), 368-370 (Restatement sweep and how to run it).
- **Test that breaks on a naive move (F1):** `_dev/tests/prescribed-shell-canonicalization.sh:86-99` requires `work-reference.md` to contain the literal `../docs/prescribed-shell-primitives.md`; its only occurrence is the "Hold both endpoints" paragraph (line 437), inside the moved range. Keep one sentence in the shortened `work-reference.md` section that links `../docs/prescribed-shell-primitives.md#state-across-command-blocks`.
- **Citations inside the two files, outside the moved ranges (F2), to re-point at `actions/fan-out-reference.md`:** `work-reference.md` 19, 55, 112, 491, 848; `work.md` 35, 37, 39, 106, 286, 309, 315, 370, 475, 478, 494, 537, 550, 551, 583.
- **Bold citations in other files that the reference test checks (F3), to re-point:** `review-work.md:72, :74`; `cleanup.md:39, :134`; `restart-with-parallel-handoff.md:82`; `capture.md:124` (cites `actions/work.md` → **Mid-Run Messages (any step)**); `crew-members/background-agents.md:187` in do-work, and the identical copies in `do-work-knowledge` and `do-work-toolbox` (the sibling copies wrap the name across a line).
- **Plain-text citations that go stale silently, to re-point:** `capture-reference.md:166`; `docs/work-guide.md:99`; `do-work-board/actions/board.md:94, :120`; `do-work-board/docs/board-guide.md:61`; `do-work-board/tools/queue-kanban/lessons-do-kanban.md:37`; Go comments `model.go:1950` and `verify.go:1025, 1079, 1275, 1392` (comment-only edits).
- **Rule for the new file (F4):** `_dev/tests/contracts/core-checks.sh` fails when a Common Rationalizations row in a file closely copies a row elsewhere (ratio above 0.75); the new file carries no such table.
- **What the reference test checks:** every path-shaped token must resolve in both the source and the installed layout; a `path` → **Name** citation's bold name must appear as whole words in a heading or bold label of the target; Markdown anchors must match heading slugs. Companion files open with `# <Name> — Reference` and a `> **Companion file to \`actions/work.md\`.** …` blockquote (`capture-reference.md:1-3`, `estimate-reference.md:1-3`). Same-package citations are spelled `actions/<file>.md` → **Section**; cross-package ones use the literal relative path (`../../do-work/actions/…`).
- **No registration needed:** `SKILL.md` routes only public actions; `staged-skills-contract.sh` `core_files` is a must-exist list, not a registry; the audit and core-check globs pick up the new file automatically.
- **No test pins the moved heading strings against the live files** (the hits in `shipped-package-reference-contract.sh` are a comment and parser fixtures).

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `skills/do-work/actions/fan-out-reference.md` (new)
- `skills/do-work/actions/work.md`
- `skills/do-work/actions/work-reference.md`
- `skills/do-work/actions/review-work.md`
- `skills/do-work/actions/cleanup.md`
- `skills/do-work/actions/restart-with-parallel-handoff.md`
- `skills/do-work/actions/capture.md`
- `skills/do-work/actions/capture-reference.md`
- `skills/do-work/crew-members/background-agents.md`
- `skills/do-work-knowledge/crew-members/background-agents.md`
- `skills/do-work-toolbox/crew-members/background-agents.md`
- `skills/do-work/docs/work-guide.md`
- `skills/do-work-board/actions/board.md`
- `skills/do-work-board/docs/board-guide.md`
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md`
- `skills/do-work-board/tools/queue-kanban/model.go` (one comment)
- `skills/do-work-board/tools/queue-kanban/verify.go` (four comments)
- `decisions/log.md`

**Files I will NOT touch:** `skills/do-work/SKILL.md`, every other `crew-members/*.md`, `_dev/tests/**`, `CHANGELOG*.md`, `do-work/archive/**`, `decisions/records/**`, `kb/**`, `ai-reports/**` (historical mentions stay as they are).

**Acceptance criteria (restated from the REQ):**
1. `skills/do-work/actions/fan-out-reference.md` exists, opens with the companion header convention and a description blockquote saying why it belongs in core, and holds in order: Worktree Dispatch Mode (moved whole, including the Isolation ladder and Delegated integration), Landed hand-back and dispatch timing, Mid-Run Messages, Restatement sweep and wave-end set-aside. Every moved passage is rewritten delete-before-add; the file itself has no size limit.
2. `work.md`: Mid-Run Messages shrinks to its routing rule plus a pointer; the landed hand-back and dispatch-instant paragraphs shrink to their rule plus a pointer; the Restatement sweep keeps its MUST line and points at the procedure.
3. `work-reference.md`: Worktree Dispatch Mode becomes a short section stating the every-run-mode rule, the one-writer invariant, the integrator's entry and never-run list, one sentence that still links `../docs/prescribed-shell-primitives.md#state-across-command-blocks`, and a pointer to the new file. The run-directory table row stays.
4. Every citation of a moved heading outside historical records points at the new file: the grep in the REQ's requirement 4 finds only the new file, the pointers, and historical records.
5. `wc -w`: `work.md` below 12,711 and `work-reference.md` below 22,162.
6. `decisions/log.md` gains today's entry reaffirming ADR-001 with the explicit no-size-limit for the companion.
7. `bash _dev/tests/shipped-package-reference-contract.sh`, `bash _dev/tests/contract-regressions.sh` and `bash _dev/tests/prescribed-shell-canonicalization.sh` are green.

## Pre-Flight

**Git:** ✓ Clean at c951e8e6 apart from this REQ's own working file and `do-work/working/baseline.json` (rewritten by this pre-flight); the two sibling working REQs (REQ-650, REQ-651) are claimed by this session and committed.
**Tests baseline:** ✓ Repository gate `bash _dev/tests/maintainer-verify.sh` exit 0 at c951e8e6 (124 s wall, both fast stages executed, fingerprint_mismatch); focused baseline `bash _dev/tests/shipped-package-reference-contract.sh` green.
**Dependencies:** ✓ Go 1.26.1 on PATH (needed only for `go vet` on the comment-only Go edits)

*Checked by work action*
