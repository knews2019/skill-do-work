---
id: REQ-653
title: '[impact-rule-change] Fan-out orchestration prose is compressed in place and moved to a companion reference'
status: completed
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
builder_handback_at: 2026-10-09T17:01:16Z
integration_at: 2026-10-09T17:01:33Z
kb_status: pending
review_at: 2026-10-09T17:16:03Z
commit: ce1a14b69c6ded6028e5697143cf85debf715413
heavy_verified_at: 2026-10-09T17:16:24Z
heavy_verified_revision: ce1a14b69c6ded6028e5697143cf85debf715413
completed_at: 2026-10-09T17:16:58Z
release_at: 2026-10-09T17:16:58Z
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
- [x] **[PLAN]:** Move work-reference.md 394-488 whole into a new companion and restructure its bold labels as headings, so every former `→ **Name**` still resolves. Move work.md Mid-Run Messages, Landed hand-back, dispatch instant and the wave-end half of How to run it into three more sections. Leave rule-plus-pointer stubs. Then re-point every citation the exploration listed and close F1 by keeping the shell-primitives link in work-reference.md. Moving those passages alone left work.md about 250 words over 12,711. The rest came from compressing the other fan-out paragraphs in work.md, two duplicated hand-back sentences and the REQ-652 operator row, which the REQ's Constraints name for compression.
- [x] **[APPLY]:** Done exactly as planned in the 18 Scope files, with no file outside the list. The edits were made with exact-string replacement scripts that assert one match each.
- [x] **[UNIFY]:** `git diff --stat` shows 18 files changed, 225 insertions and 144 deletions. `git diff --check` is clean. Checked files: the two big files and the new file, read in full after the edit. Every re-pointed citation line was confirmed by the citation grep. The two sibling background-agents.md copies match under `cmp`. `go vet ./` and `gofmt -l` are clean on the Go comments.
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

## Implementation Summary

**Files changed:**
- `skills/do-work/actions/fan-out-reference.md` (new)
- `skills/do-work/actions/work.md` (modified)
- `skills/do-work/actions/work-reference.md` (modified)
- `skills/do-work/actions/review-work.md` (modified)
- `skills/do-work/actions/cleanup.md` (modified)
- `skills/do-work/actions/restart-with-parallel-handoff.md` (modified)
- `skills/do-work/actions/capture.md` (modified)
- `skills/do-work/actions/capture-reference.md` (modified)
- `skills/do-work/crew-members/background-agents.md` (modified)
- `skills/do-work-knowledge/crew-members/background-agents.md` (modified)
- `skills/do-work-toolbox/crew-members/background-agents.md` (modified)
- `skills/do-work/docs/work-guide.md` (modified)
- `skills/do-work-board/actions/board.md` (modified)
- `skills/do-work-board/docs/board-guide.md` (modified)
- `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` (modified)
- `skills/do-work-board/tools/queue-kanban/model.go` (modified)
- `skills/do-work-board/tools/queue-kanban/verify.go` (modified)
- `decisions/log.md` (modified)

**What was done:** The Worktree Dispatch Mode section of `work-reference.md` moved whole into the new companion `actions/fan-out-reference.md`, with its bold labels turned into headings, and work.md's Mid-Run Messages detail, landed hand-back and dispatch-instant paragraphs and the wave-end sweep procedure moved there too. Each source passage now states its rule plus a pointer; `work.md` went from 13,864 to 12,706 words and `work-reference.md` from 22,815 to 17,046 words (`wc -w` on the merged tree). Every citation of a moved heading across do-work, do-work-board and the three background-agents.md copies now points at the companion, and `decisions/log.md` records the third reaffirmation of ADR-001. Merge range `ad74c131..ce1a14b6`, cumulative: builder commit d1bf4855 merged at 35803b1e, then review-fix commit dfa8a532 (restores the "an unmerged one can pass" reason in Cleanup, corrects two sizes in the decisions entry) re-merged at ce1a14b6. No integration seams. The companion is 5,732 words after the fix.

## Decisions

*(from the builder hand-back)*

- **D-01 DECIDE & STATE: F1 is closed in prose, not in the test.** work-reference.md's short section keeps one sentence that links `../docs/prescribed-shell-primitives.md#state-across-command-blocks`. That sentence holds the re-typed-literals rule for `<pre>`/`<merge_hash>`, so it is real content and not padding. The test stays unedited because `_dev/tests` is outside the write boundary.
- **D-02 DECIDE & STATE: these headings stay in the source files.** work.md keeps `### Mid-Run Messages (any step)`, `### Mechanical Evidence-Gate Loop`, every Step heading and the Step 6 "Always load" lines. work-reference.md keeps `## Worktree Dispatch Mode (Step 1)`. The new file repeats both names as headings so old and new citations both resolve.
- **D-03 DECIDE & STATE: the moved bold labels became headings.** Examples include Isolation ladder, Naming, State stays home, Sole integrator, When to merge, Fan-Out Dispatch, Auto-wave, Serial-only and Delegated integration. This meets the "structured companion" ask, and the citation checker accepts a heading or a bold label. Citations now name the sub-heading directly, for example `actions/fan-out-reference.md` → **Fan-Out Dispatch**, and no longer chain through **Worktree Dispatch Mode**.
- **D-04 DECIDE & STATE: these cuts were restatement.** In each case the rule and the named commands stayed.
  - In the moved work-reference text, I cut repeated "above/below" pointers, re-explanations of the operative name inside step 2, and the `board.md` Go-check contrast in the concurrency-degrade paragraph. I also cut a second statement of why the merge is the non-interference proof (Auto-wave repeated Fan-Out Dispatch), the "Reached two ways" restatement of work.md Architecture, and "Step 8 deletes the branch before verify" (it only explained "blind after the merge").
  - In Mid-Run Messages, the four "under delegated integration the coordinator writes in its next gap" clauses became one "Who writes" rule above the branches.
  - In work.md, I deleted the `write_set is not an input to the wave` Architecture paragraph because the Rules bullet states the same rule. I deleted two copies of "in worktree dispatch mode this section goes in your hand-back" and folded them into one sentence in the "Write only on your own branch" bullet. I deleted the `--fan-out` bullet's repeat of the degrade condition because Architecture states it.
- **D-05 DECIDE & STATE: one dead citation was fixed.** work.md's Rules bullet on `primary_commit` cited **Worktree Dispatch Mode (Step 1)** for the "second metadata commit". That section never mentioned metadata commits. It now cites `actions/work-reference.md` → **Commit & Metadata-Commit Procedure (Step 9)**.
- **D-06 DECIDE & STATE: the integrate.md row moved with the table.** The REQ says "the run-directory table keeps its row", and the brief says the table moves. The table now lives in the companion, and the `REQ-NNN-integrate.md` row is kept there.
- **D-07 DECIDE & STATE: "byte-identical" covers two of the three copies.** The do-work copy of background-agents.md differs from the two sibling copies by design (same-package versus `../../do-work/` paths and line wraps, as before this REQ). The knowledge and toolbox copies are still byte-identical, and all three got the same edit.
- **D-08 DECIDE & STATE: the operator row was compressed by one clause only.** "Invisible as waiting on you, and the default run hides it" became "the default run skips it, hiding the wait from the user". The REQ's Constraints ask for the REQ-652 operator sentences to be compressed with the rest. The work-reference.md Environment row (REQ-652) was left as it is because that file was already far under its limit.

## Discovered Tasks

*(from the builder hand-back)*

- Some stale plain-text mentions of `work-reference.md` → Worktree Dispatch Mode were deliberately left as historical or out of scope: kb/wiki and ai-reports (about 150), `do-work/HANDOFF-2026-08-18-queue-241-245.md:54`, CHANGELOG entries and ADR-018. → report only.

## Qualification

**Gate records:** `advance --diff-range ad74c131..35803b1e` returned `qualify` satisfied and `scope-drift` satisfied, both with provenance `merged_range` and no findings (no debug artifacts, P-A-U boxes ticked with the hand-back's text).

**Scope comparison (Route B):** (re-checked on the cumulative range `ad74c131..ce1a14b6`: same 18 files) the 18 files declared under `## Scope` and the 18 files in `git diff --name-only ad74c131..35803b1e` are an exact match; nothing outside the declaration was touched and every declared file changed.

**Requirement trace (read from the merged tree, not the hand-back):**
1. `skills/do-work/actions/fan-out-reference.md` exists with `# Fan-Out Orchestration — Reference` and a `> **Companion file to \`actions/work.md\`.**` blockquote that says it belongs in core because every run mode integrates a builder branch. Its `##` sections are, in order: Worktree Dispatch Mode (Step 1) with Isolation ladder through Run directory as `###` headings, Landed hand-back and dispatch timing, Mid-Run Messages (any step), Restatement sweep and wave-end set-aside. ✓
2. `work.md`: Mid-Run Messages is one routing paragraph plus a pointer; the landed hand-back and dispatch-instant paragraphs state their rule and point at the companion; the Restatement sweep MUST line is unchanged and How to run it points at the wave-end section. ✓
3. `work-reference.md` → **Worktree Dispatch Mode (Step 1)** is one paragraph with the every-run-mode rule, the one-writer rule, the integrator entry (`advance REQ-NNN`) and never-run list (`recover`, `--take-over`, `--assume-sole-authority`, Step 10), the `../docs/prescribed-shell-primitives.md#state-across-command-blocks` link, and the pointer. The run-directory table, with its `REQ-NNN-integrate.md` row, moved whole to the companion (D-06). ✓
4. The requirement-4 grep on the merged tree finds only the companion, the pointers and kept headings, review-work.md's own sweep text, test fixtures, ADR-018, CHANGELOG entries and the new decisions entry. A narrower check for any `work.md`/`work-reference.md` citation of a moved heading finds only correct targets (review-work.md Step 6, work.md Step 7). ✓
5. `wc -w`: `work.md` 12,706 (< 12,711), `work-reference.md` 17,046 (< 22,162). ✓
6. `decisions/log.md` has `## [2026-10-09] Fan-out prose moved to a companion reference` reaffirming ADR-001 with the explicit no-size-limit and the sizes before and after. ✓
7. Contract tests: run below under Testing.

**Judgment:** Substantive change, no rule condition visibly changed in the passages read; the full rule-by-rule survival check is delegated to review.

## Testing

**Tests run:** `DO_WORK_FAST_STAGE_REUSE=off bash _dev/tests/maintainer-verify.sh` at merge 35803b1e (repository gate, both Go fast stages executed with reuse disabled, 120 s wall) and again at the review-fix re-merge ce1a14b6 (exit 0, 130 s wall, with the probe re-run directly, exit 0); then `advance --gate-arg bash --gate-arg _dev/tests/maintainer-verify.sh --gate-exit-status 0 -- --probe-file do-work/runs/work-2026-10-09-160815/REQ-653-probe.sh` (probe: companion file exists, both word counts under their limits, `shipped-package-reference-contract.sh` and `prescribed-shell-canonicalization.sh`).
**Result:** ✓ All passing. Gate exit 0 (contract regressions included; queue-kanban 422 Go tests, do-work-cli 879 Go tests). Advance recorded `green-gate`, `run-blocked-check` (probe status 0) and `scope-drift` satisfied.

**Red-green validation:** `tdd: false`; the captured `## Red-Green Proof` is a file-and-count read, checked on the merged tree:
- RED (hand-back, base e36f92ed): `skills/do-work/actions/fan-out-reference.md` absent; `wc -w` 13,864 and 22,815.
- GREEN (merged tree, 35803b1e): file present with the four sections in order; `wc -w` 12,706 for `work.md` and 17,046 for `work-reference.md`; requirement-4 grep clean (see Qualification); both contract tests green inside the gate and the probe; `decisions/log.md` entry present.

**New tests added:** none (prose move; the probe is run evidence, not a shipped test).

**Heavy verification plan:**
- Range: ad74c131101355f5e9a262859f1b9c3f65a364b9..ce1a14b69c6ded6028e5697143cf85debf715413 (re-planned after the review-fix re-merge; same three lanes as at 35803b1e)
- queue-kanban-javascript: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — `lessons-do-kanban.md`, `model.go`, `verify.go` matched subtree `skills/do-work-board/tools/queue-kanban`
- queue-kanban-browser: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — same three paths matched subtree `skills/do-work-board/tools/queue-kanban`
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — all 17 changed paths under `skills/` matched subtree `skills`

*Verified by work action*

## Review

**Overall: 96%** | 2026-10-09T17:16:03Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:** F1 `decisions/log.md:126` gave the 0.305.73 work-reference.md size as 22,109; fixed in ce1a14b6, which now says 22,162 (verified at `bd6c9798`) and gives the companion as 5,732 (verified by `wc -w`) — impact-negligible → report only. F2 `skills/do-work/actions/fan-out-reference.md:95` Cleanup — happy path had dropped the false-pass reason; fixed in ce1a14b6 ("a merged branch can refuse and an unmerged one can pass"), and the rule is unchanged — impact-negligible → report only. F3 (Nit) two cuts missing from D-04, both restated elsewhere: `work.md:33` serial-loop description and Sole integrator's "builder that edits the seam" sentence — impact-negligible → report only.
**Acceptance:** Pass — the contract tests exit 0 on 35803b1e, and shipped-package-reference-contract exits 0 again on ce1a14b6. The review-fix delta `35803b1e..ce1a14b6` (2 files, 2 lines) is correct and introduces nothing new. Cumulative range ad74c131..ce1a14b6.
**Restatement sweep:** redefined the home of the fan-out contract: Worktree Dispatch Mode and its sub-rules (Isolation ladder, Naming, operative name, State stays home, Sole integrator, When to merge, Cleanup, Fan-Out Dispatch, Auto-wave, Serial-only, Delegated integration, run-directory table), Landed hand-back, Dispatch instant, Mid-Run Messages branches, and the wave-end set-aside procedure all moved from `work-reference.md`/`work.md` to `actions/fan-out-reference.md`. Swept skills/, _dev/, decisions/ and CLAUDE.md; every live citation lands on an existing heading with the expected meaning, and no stale restatement was found. Inherited at the wave end: REQ-650's repo-root-only disk probe is clean, REQ-651 had nothing redefined, REQ-652's operator-blocked rule is clean.
**Suggested testing:** 1 items
**Follow-ups created:** None (3 findings report only)

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Moving the whole Worktree Dispatch Mode section and turning its bold labels into headings let every old `→ **Name**` citation keep resolving while the citations were re-pointed one file at a time. Keeping the two source headings (`### Mid-Run Messages (any step)`, `## Worktree Dispatch Mode (Step 1)`) as short rule-plus-pointer stubs meant no stale citation could dangle during the move.
**What didn't:** Moving the passages alone left `work.md` about 250 words over its 12,711 ceiling; the margin came from deleting restatements the move exposed (the same `write_set` and hand-back rules stated in Architecture, Input and Rules). Compression also dropped one reason clause in Cleanup ("an unmerged one can pass") that review caught and the fix restored.
**Worth knowing:** A size number written into a decisions entry needs the same `wc -w` check as the REQ's own ceiling: the triage figure 22,109 disagreed with the 0.305.73 measurement 22,162, and a review fix changes the companion's own count. `_dev/tests/prescribed-shell-canonicalization.sh` requires `work-reference.md` to keep the `../docs/prescribed-shell-primitives.md` link, so that sentence must stay in the short section.

## Orientation

The fan-out and delegated-integration contract for `do-work run` now lives in `actions/fan-out-reference.md`; `work.md` and `work-reference.md` keep each rule in a sentence or two and point there (action-files prime area). [MAP CHANGED] a new companion file is now the home that citations of the merge sequence, isolation ladder, Mid-Run Messages and the wave-end sweep must name. Touched primes' referenced paths still exist.

## Heavy Verification Plan

- Base: ad74c131101355f5e9a262859f1b9c3f65a364b9
- Target: ce1a14b69c6ded6028e5697143cf85debf715413
- queue-kanban-javascript: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-javascript` — `lessons-do-kanban.md`, `model.go`, `verify.go` matched subtree `skills/do-work-board/tools/queue-kanban`
- queue-kanban-browser: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane queue-kanban-browser` — same three paths matched subtree `skills/do-work-board/tools/queue-kanban`
- staged-skills: `env GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null bash _dev/tests/maintainer-verify.sh --heavy-lane staged-skills` — all 17 changed paths under `skills/` matched subtree `skills`

## Heavy Verification Result

- Target revision: ce1a14b69c6ded6028e5697143cf85debf715413
- Execution revision: ce1a14b69c6ded6028e5697143cf85debf715413 (detached checkout, `QUEUE_KANBAN_BROWSER` set to Google Chrome)
- queue-kanban-javascript: exit 0, reused (executed green at 35803b1e in 8 s; its fingerprint inputs did not change in the review-fix delta)
- queue-kanban-browser: exit 0, executed, 84 s
- staged-skills: exit 0, executed, 37 s
- Earlier drain at the first merge 35803b1e: all three executed green (8 s, 80 s, 35 s).

## Timing

Observed 2026-10-09T17:01:16Z to 2026-10-09T17:16:13Z: 14m 57s total, 18m 04s attributed across 5 events, 0s unattributed.

| Category | Elapsed | Events |
| --- | --- | --- |
| review | 10m 48s | 1 |
| verification-gate | 6m 59s | 3 |
| handback-merge | 17s | 1 |

Slowest stage: review / independent review with wave-end sweep and delta re-review, 10m 48s, outcome success.
