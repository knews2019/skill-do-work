## Review: REQ-653

**Approve** — the fan-out prose moved to `actions/fan-out-reference.md` with every checked rule intact, every citation lands, and the size targets are met. One number in the decisions log is wrong.
Route B | ad74c131..35803b1e

### What's built
- New companion `skills/do-work/actions/fan-out-reference.md` (5,721 words) holds Worktree Dispatch Mode with its bold labels turned into headings, the landed hand-back and dispatch-instant rules, the Mid-Run Messages branches, and the wave-end sweep procedure.
- `work.md` is 12,706 words (limit 12,711) and `work-reference.md` is 17,046 words (limit 22,162). Each moved passage left a rule sentence and a pointer.
- Citations in four packages, including the Go comments and all three `background-agents.md` copies, now point at the companion.

### Decisions / risks for you
- None. Release (requirement 7, second half) is the integrator's Step 9 job and was not part of this review.

### Findings

**Important:**
- None.

**Minor:**
- F1. `decisions/log.md:126` says `work-reference.md` regrew "from 22,109 to 22,815". At 0.305.73 (`bd6c9798`) `git show bd6c9798:skills/do-work/actions/work-reference.md | wc -w` gives 22,162. That is the baseline requirement 5 names, and no 0.305.73 revision has 22,109. The other numbers in the entry are correct (12,711, 13,864, 22,815, 12,706, 17,046, 5,721). Fix: replace `22,109` with `22,162`. — impact-negligible → report only
- F2. `skills/do-work/actions/fan-out-reference.md:95` (Cleanup — happy path) dropped the clause "and an unmerged one can pass". Pre-move text: "from anywhere else a perfectly-merged branch refuses and an unmerged one can pass". The rule (run `branch -d` from the integration branch, never `-D`/`--force`) is unchanged. Only the reason for the more dangerous direction, a false pass that deletes unmerged work, is gone. — impact-negligible → report only

**Nit:**
- F3. Two small cuts are not listed in D-04, and both are stated elsewhere. One is `work.md:33` "Step 1 advances one claim, Step 6 waits for its builder before Step 6.25 reads the output…". The other is Sole integrator's "A builder that edits the seam itself writes the main tree the orchestrator alone owns" (covered by the section's first sentence). — impact-negligible → report only

### Requirements Checklist

- [x] 1. Companion exists with the `— Reference` header and a blockquote saying why it belongs in core. Sections are in the required order: Worktree Dispatch Mode, Landed hand-back and dispatch timing, Mid-Run Messages, Restatement sweep and wave-end set-aside. — delivered
- [x] 2. `work.md`: the Mid-Run Messages stub keeps the routing rule plus a pointer. The landed hand-back and dispatch-instant paragraphs keep their rule plus a pointer. The Restatement sweep MUST line is byte-unchanged and How to run it points at the procedure. — delivered
- [x] 3. `work-reference.md` Worktree Dispatch Mode is one paragraph. It has the every-run-mode rule, the one-writer rule, the integrator entry `advance REQ-NNN`, and the never-run list (`recover`, `--take-over`, `--assume-sole-authority`, Step 10). It keeps the `../docs/prescribed-shell-primitives.md#state-across-command-blocks` link and a pointer. The run-directory table, with its `REQ-NNN-integrate.md` row, moved whole (D-06). — delivered
- [x] 4. The requirement-4 grep on HEAD finds only the companion, the kept headings and pointers, review-work.md's own sweep text, test fixtures in `_dev/tests/shipped-package-reference-contract.sh`, and the new decisions entry. — delivered
- [x] 5. `wc -w`: 12,706 < 12,711 and 17,046 < 22,162. — delivered
- [~] 6. The `decisions/log.md` entry is present and reaffirms ADR-001 with the no-size-limit. One starting size is wrong (F1). — delivered with a factual error
- [x] 7. Re-run here: `shipped-package-reference-contract.sh` exit 0, `contract-regressions.sh` exit 0 ("Contract regression checks passed."), `prescribed-shell-canonicalization.sh` exit 0, `_dev/tests/contracts/core-checks.sh` exit 0. — delivered (release is Step 9)

### Coordinator checks

1. **Moved rules survived.** I compared the companion against `git show ad74c131:skills/do-work/actions/work-reference.md` lines 394-488 and the pre-merge `work.md` passages:
   - Isolation ladder: all three rungs survived, including the branch-rung "reuse that name", "isolates committed work only" and the no-rung progress line verbatim.
   - The operative-name rule survived, including the re-derive failure.
   - State stays home survived: "examples, not the extent", the stale-snapshot rule and the branch-rung role boundary.
   - Sole integrator survived with its one exception: own hand-back, absolute path, never staged, committed or merged.
   - Merge steps 0-4 survived:
     - Step 0 keeps the three categories, the index guard, `MM` loss, the `git commit -- do-work/` ban and the order below `<pre>`.
     - Step 1 still captures `<pre>` once.
     - Step 2 keeps the three-dot queue guard before the merge and the `Already up to date.` stop.
     - Step 3 keeps the seam inside the merge commit.
     - Step 4 keeps `<merge_hash>`.
   - The cumulative range survived (`<pre₁>..<merge_hash₂>`, Step 9 records the latest).
   - Cleanup keeps "Never `-D`, never `--force`", report and stop.
   - Serial-only is byte-identical.
   - Delegated integration keeps the never-run list (run-with-recovery, `--assume-sole-authority`, `--take-over`) and the one-writer rule. Its brief contents are complete.
   - The run-directory table still has the integrate row.
   - Mid-Run Messages keeps all five branches. The four per-branch "coordinator writes in its next gap" clauses became one **Who writes** rule that covers every branch (D-04). Same meaning.
   - The Restatement sweep MUST is unchanged. The set-aside handling (HANDLED, re-run none) is unchanged.

   No condition changed and no rule was dropped. Only reasons were cut (F2, F3).
2. **core-checks pins.** `### Mechanical Evidence-Gate Loop`, `advance.gate_records`, `## Qualification Anti-Rationalization Table (Step 6.3)`, `## Finding-Closure Ratchet (Step 6.5)`, the Step 8/9 tokens and the Changelog/Commit procedure anchors are all present. `core-checks.sh` exits 0.
3. **D-05.** `## Commit & Metadata-Commit Procedure (Step 9)` exists at `work-reference.md:808`. It is where `primary_commit` versus `supplied_commit` provenance is defined, so it fits the "second metadata commit" sentence better than the old target, which never mentioned it.
4. **background-agents.md copies.** The knowledge and toolbox copies are byte-identical (`cmp`). The do-work copy differs only by same-package path and line wrap. The diff is 23 lines both before and after the merge. All three cite **Worktree Dispatch Mode (Step 1)** in the companion with the same meaning.
5. **D-08 operator row.** Before: "the earmark is for another session, it is invisible as "waiting on you", and the default run hides it from the user by skipping it". After: "the earmark is for another session, and the default run skips it, hiding the wait from the user". The rule (never write `assigned_to`, flip to `blocked` with `blocked_by`, routine operator act completes) is unchanged. It agrees with the unchanged Environment row at `work-reference.md:697`.
6. **Decisions sizes.** Pre-merge (ad74c131): 13,864 / 22,815. Merged: 12,706 / 17,046. Companion: 5,721. All match the entry except the 0.305.73 starting value for `work-reference.md` (F1).

### Acceptance Testing

**Result: Pass**
- The three contract tests named by the REQ plus `core-checks.sh` all ran green on HEAD (35803b1e).
- A cold read of the companion as an orchestrator found every in-file italic cross-reference resolving to a heading. Examples: *Naming*, *Sole integrator*, *Cleanup — happy path*, *Hold both endpoints as re-typed literals*, *Delegated integration — the coordinator shape*, *Mid-Run Messages (any step)*.
- Every outward target exists: Crash Recovery (Step 1), In-Progress Record (Step 1), Execution Model, Schema Read Contract → Dependency-source-ready status set, Decision Brief, Stamps are append-only, clarify's Outside-text containment, capture's Immutability Rule, cleanup Pass 5, and background-agents' *Write a manifest per wave*.

### Wave-end Restatement Sweep (inherited, checked against HEAD 35803b1e)

- **REQ-650 (the disk probe now measures only the repo root):** clean. A grep for "per device", "each builder worktree", "deviceIdentity", "measuredDevices" and "disk-space probe for worktrees" outside CHANGELOGs found nothing. `prime-do-kanban.md:37`, `board-guide.md:78` and `verify.go:129/267` all say repo root. REQ-653's `verify.go` edits are comment-only and do not touch the probe.
- **REQ-651 (AI report on board activity correlation):** recorded "nothing redefined". Nothing to sweep.
- **REQ-652 (operator-blocked rule):** clean. The compressed `work.md:515` row, the unchanged `work-reference.md:697` Environment row, the `assigned_to` schema line `work-reference.md:114`, `capture.md:107-108`, `clarify.md:200` and `board.md:92` all agree.

### Own sweep (REQ-653 moved where the fan-out contract lives)

- I ran the requirement-4 grep, plus citations of Sole integrator, State stays home, Naming, Serial-only, Auto-wave, When to merge, Cleanup, Landed hand-back, Dispatch instant, Hold both endpoints, operative name, Claim conflicts and the run-directory table, over `skills/`, `_dev/`, `decisions/` and `CLAUDE.md`. I also checked every `work-reference` citation near worktree, fan-out, merge range, coordinator, integrator, isolation, dispatch or mid-run wording.
- Every live citation lands on an existing heading with the expected meaning. Examples:
  - `review-work.md:72,74`, `cleanup.md:134` and the three `background-agents.md` copies → **Worktree Dispatch Mode (Step 1)**.
  - `capture.md:124` → **Mid-Run Messages (any step)**.
  - `work.md:307` → **When to merge, and the range every evidence step reads**.
  - `work.md:362` → **Restatement sweep and wave-end set-aside**.
  - `work.md:470` → **Delegated integration — the coordinator shape**.
  - The board docs and Go comments → Fan-Out Dispatch / "Naming" / "state stays home" / "Cleanup — happy path".
  - `work-reference.md:400` → **State stays home**.
  - `work-reference.md:757` → **Mid-Run Messages (any step)**.
- No pointers to moved content are left dangling inside `work-reference.md`.
- Out of scope, as the brief says: kb/, ai-reports/, CHANGELOGs, ADR-018 and the HANDOFF file.

### Suggested Additional Testing

- In the next real `--fan-out` run with delegated integration, have the integrator follow only the companion's **Delegated integration** and **When to merge** sections. Confirm it needs nothing from the old `work-reference.md` text.

### Scores (on the record — not the headline)

**Overall: 94%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 95% | All seven delivered; one wrong starting size in the log entry (F1) |
| Code Quality | 92% | Rules intact; two reason clauses cut without a D-04 entry (F2, F3) |
| Test Adequacy | 90% | Prose move; contract tests and probe are the right proof, all green |
| Scope | 100% | 18 declared files = 18 changed files |
| Risk | Low | Citation drift is the main risk; the reference contract test plus the sweep cover it |
| Acceptance | Pass | Contract tests green; companion reads cold and resolves |

### Follow-ups created
None (3 findings report only)

## Review

**Overall: 96%** | <timestamp>

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
