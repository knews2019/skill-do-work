# REQ-653 hand-back (fan-out orchestration prose moved to a companion reference)

- **Branch:** worktree-agent-REQ-653-fan-out-orchestration-prose-companion-reference
- **Worktree:** /Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-653-fan-out-orchestration-prose-companion-reference
- **Base:** e36f92ed
- **Commits:** d1bf4855 (one commit)

## File manifest

| File | Action | What changed |
| --- | --- | --- |
| `skills/do-work/actions/fan-out-reference.md` | new | Companion with the `capture-reference.md` header shape and a blockquote saying why it is core. Four `##` sections in order: Worktree Dispatch Mode (Step 1), Landed hand-back and dispatch timing, Mid-Run Messages (any step), Restatement sweep and wave-end set-aside. Former bold labels became `###`/`####` headings. 5,721 words. No Common Rationalizations table. |
| `skills/do-work/actions/work.md` | modified | 13,864 → 12,706 words. Mid-Run Messages is now its routing rule plus a pointer. Landed hand-back and dispatch-instant paragraphs are now rule plus pointer. The Restatement sweep MUST line is unchanged, and How to run it points at the wave-end procedure. Fan-out paragraphs in Architecture, the `--fan-out` input bullet, the builder bullets, Step 10, the `write_set` rule and the operator row were compressed. Every F2 citation was re-pointed. |
| `skills/do-work/actions/work-reference.md` | modified | 22,815 → 17,046 words. `## Worktree Dispatch Mode (Step 1)` is one paragraph. It keeps the every-run-mode rule, the one-writer invariant, the integrator entry and never-run list, the `../docs/prescribed-shell-primitives.md#state-across-command-blocks` link, and the pointer. Citations at lines 19, 55, 112, Reading a Builder-Authored Section and the Decision Brief were re-pointed. |
| `skills/do-work/actions/review-work.md` | modified | Two bold citations re-pointed. |
| `skills/do-work/actions/cleanup.md` | modified | Two bold citations re-pointed. |
| `skills/do-work/actions/restart-with-parallel-handoff.md` | modified | One bold citation re-pointed. |
| `skills/do-work/actions/capture.md` | modified | Mid-Run Messages citation re-pointed to the companion. |
| `skills/do-work/actions/capture-reference.md` | modified | Plain-text citation re-pointed. |
| `skills/do-work/crew-members/background-agents.md` | modified | Citation re-pointed. |
| `skills/do-work-knowledge/crew-members/background-agents.md` | modified | Citation re-pointed (wrapped line kept). |
| `skills/do-work-toolbox/crew-members/background-agents.md` | modified | Same edit, still byte-identical to the knowledge copy. |
| `skills/do-work/docs/work-guide.md` | modified | Plain-text citation re-pointed. |
| `skills/do-work-board/actions/board.md` | modified | Two plain-text citations re-pointed (`../../do-work/actions/fan-out-reference.md`). |
| `skills/do-work-board/docs/board-guide.md` | modified | One citation re-pointed. |
| `skills/do-work-board/tools/queue-kanban/lessons-do-kanban.md` | modified | One citation re-pointed (`../../../do-work/...`). |
| `skills/do-work-board/tools/queue-kanban/model.go` | modified | One comment re-pointed (two lines became one). |
| `skills/do-work-board/tools/queue-kanban/verify.go` | modified | Four comments re-pointed. |
| `decisions/log.md` | modified | New `## [2026-10-09] Fan-out prose moved to a companion reference` entry. |

## AI Execution State (P-A-U Loop)

- [x] **[PLAN]:** Move work-reference.md 394-488 whole into a new companion and restructure its bold labels as headings, so every former `→ **Name**` still resolves. Move work.md Mid-Run Messages, Landed hand-back, dispatch instant and the wave-end half of How to run it into three more sections. Leave rule-plus-pointer stubs. Then re-point every citation the exploration listed and close F1 by keeping the shell-primitives link in work-reference.md. Moving those passages alone left work.md about 250 words over 12,711. The rest came from compressing the other fan-out paragraphs in work.md, two duplicated hand-back sentences and the REQ-652 operator row, which the REQ's Constraints name for compression.
- [x] **[APPLY]:** Done exactly as planned in the 18 Scope files, with no file outside the list. The edits were made with exact-string replacement scripts that assert one match each.
- [x] **[UNIFY]:** `git diff --stat` shows 18 files changed, 225 insertions and 144 deletions. `git diff --check` is clean. Checked files: the two big files and the new file, read in full after the edit. Every re-pointed citation line was confirmed by the citation grep. The two sibling background-agents.md copies match under `cmp`. `go vet ./` and `gofmt -l` are clean on the Go comments.

| Check (from worktree root) | Exit | Wall |
| --- | --- | --- |
| `bash _dev/tests/shipped-package-reference-contract.sh` | 0 | 1 s |
| `bash _dev/tests/prescribed-shell-canonicalization.sh` | 0 | <1 s |
| `bash _dev/tests/contract-regressions.sh` | 0 | 27 s |
| `bash _dev/tests/contracts/core-checks.sh` (runs standalone; `near_identical_cross_file_pairs 0`) | 0 | 6 s |
| `cd skills/do-work-board/tools/queue-kanban && go vet ./` | 0 | <1 s |
| `git diff --check` | 0 | — |

## Red-green evidence

- **RED:** `skills/do-work/actions/fan-out-reference.md` did not exist. `wc -w` gave 13,864 for work.md and 22,815 for work-reference.md. The fan-out detail sat in nine passages across the two files.
- **GREEN:** The new file exists with the four sections in the required order. `wc -w` now gives 12,706 for work.md (limit 12,711) and 17,046 for work-reference.md (limit 22,162). The citation grep finds only the following:
  - the new file;
  - the pointers;
  - the two kept headings (`### Mid-Run Messages (any step)` in work.md and `## Worktree Dispatch Mode (Step 1)` in work-reference.md);
  - work.md's own Restatement sweep MUST line and its checklist line;
  - review-work.md's own **Restatement sweep:** report label;
  - the test-parser fixtures and comment in shipped-package-reference-contract.sh;
  - ADR-018, which is a historical record;
  - the new decisions entry.

  All four test scripts pass, and decisions/log.md has the entry.

## Decisions

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

- Some stale plain-text mentions of `work-reference.md` → Worktree Dispatch Mode were deliberately left as historical or out of scope: kb/wiki and ai-reports (about 150), `do-work/HANDOFF-2026-08-18-queue-241-245.md:54`, CHANGELOG entries and ADR-018. → report only.

## Lessons read

- `_dev/primes/lessons-action-files.md`: family-targeted. I read `[family: condition-preserving-prose-extraction]` (REQ-509) and `[family: restated-mechanism-unchecked]` (REQ-639).
- `_dev/primes/lessons-shell-commands.md`: family-targeted. I read `[family: moved-prose-loses-its-anchor]` (REQ-554), which is why every "above/below" in the moved text was rewritten to a named section or a file-qualified citation.
- I read the primes `prime-action-files.md` and `prime-releases.md`. No required entry was missing.

## Proposed lesson bullet (integrator writes it)

```
- [family: condition-preserving-prose-extraction] [REQ-653: a word-count ceiling on a source file is met by the moved passages only if their pointers are tiny; the remaining margin came from deleting restatements that the move exposed (the same rule stated in Architecture, Input and Rules), so search for duplicates of each moved rule before compressing anything else](../../do-work/archive/UR-143/REQ-653-fan-out-orchestration-prose-companion-reference.md#lessons-learned)
```

## Proposed CHANGELOG entry (integrator writes it)

```
## Fan-Out Orchestration Moves to Its Own Reference

The orchestrator files had regrown with dense fan-out prose, and a consumer review called it hard to follow. The maintainer chose to compress in place and to give the coordinator detail a structured companion with no size limit.

- New `actions/fan-out-reference.md` holds worktree dispatch (isolation ladder, hand-back merge sequence, cleanup, fan-out and auto-wave, delegated integration, run-directory table), landed hand-back and dispatch timing, mid-run message routing, and the wave-end restatement sweep.
- `actions/work.md` drops from 13,864 to 12,706 words, and `actions/work-reference.md` drops from 22,815 to 17,046 words. Each keeps the rule and points to the companion.
- Every citation of the moved sections is re-pointed across the core, board, knowledge and toolbox packages. No rule's condition changed.
```

## Integration seams

None.
