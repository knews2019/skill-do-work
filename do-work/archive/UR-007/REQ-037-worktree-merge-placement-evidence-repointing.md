---
id: REQ-037
title: Place the worktree merge in the step sequence and re-point the evidence-consuming steps
status: completed
claimed_at: 2026-07-29T08:09:11Z
completed_at: 2026-07-29T08:46:38Z
commit: 1348d11
route: C
kb_status: pending
created_at: 2026-07-28T22:41:30Z
user_request: UR-007
addendum_to: REQ-033
depends_on: []
related: [REQ-033, REQ-035, REQ-036, REQ-038]
batch: parallel-dispatch
domain: general
prime_files: []
tdd: false
review_generated: true
write_set:
  - actions/work.md
  - actions/work-reference.md
  - actions/review-work.md
  - tools/checks/qualify.sh
maintenance: false
---

# Place the worktree merge in the step sequence and re-point the evidence-consuming steps

## What

REQ-033's worktree dispatch mode says who merges and how (`git merge --no-ff`, dependency order) but never *when* — and every evidence-consuming step downstream assumes uncommitted main-tree work. State the merge point explicitly, and make the diff-based checks read the right evidence when the mode is on.

## Why (if provided)

Review of REQ-033 (confirmed by 2 independent adversarial verifiers per finding): after a merge the main tree is clean, so (a) `tools/checks/qualify.sh` computes its file list from `git diff` + `git diff --staged` only — Step 6.3's checks degrade to WARNs and the debug-artifact scan sees nothing; (b) `actions/review-work.md` pipeline mode reads `git diff` and exits with "nothing to review" on a clean tree — Step 7 silently skips; (c) Step 9 stages "implementation files from the Implementation Summary" that are already committed — nothing stages and the validation mismatches; the Rules bullet "one commit per request" is false in this mode (builder commits + a --no-ff merge commit).

## Detailed Requirements

- State the merge point once, in the Worktree Dispatch Mode section (`actions/work-reference.md`), as a pipeline position: the orchestrator merges when the builder hands back (end of Step 6), before Step 6.25's Implementation Summary — or an explicitly better position with reasoning. `actions/work.md`'s Step 8 post-merge-verification paragraph then references it instead of assuming it.
- When the mode is on, evidence-consuming steps read the merge range, not the working diff: specify the range once (e.g. from the integration branch's pre-merge position to HEAD — merge-base form) and have Step 6.3 (`tools/checks/qualify.sh` — add a range parameter or an env/argument fallback), Step 7 (`actions/review-work.md` pipeline mode "Get the Diff"), and Step 9's staged-list validation consume it. Serial mode unchanged — the working-diff path stays the default.
- Reconcile Step 9 for merged work: the commit already exists (builder commits + merge commit); Step 9's job becomes the changelog/version bump + writing the `commit:` hash of the merge commit into the REQ, not staging implementation files. Say so where the Rules bullet claims "one commit per request."
- Keep all four REQ-033 ratchets green; keep the serial default byte-compatible in behavior.

## Constraints

- No new action files. `tools/checks/qualify.sh` changes must keep its exit-code contract and header-comment docs in sync.

---

## Triage

**Route: C** - Complex

**Reasoning:** Touches four coupled files including an executable shell script (`tools/checks/qualify.sh`, which carries an exit-code contract + header-comment docs that must stay in sync) and the review action's diff source. The change introduces a coordinated merge-range contract that three evidence-consuming steps must read identically, and the serial (working-diff) path must stay byte-compatible. Getting the shell interface right without regressing the serial default warrants a plan and an independent adversarial review.

**Planning:** Required

## Orchestrator note (pre-plan)

⚠ **Pre-existing coherence gap to fix here (found at claim time):** REQ-035 changed `actions/work.md`'s Step 8 post-merge-verification paragraph to "per-merge is the default whenever more than one REQ is in flight," but the source paragraph it summarizes — `actions/work-reference.md` **Worktree Dispatch Mode → "Post-merge verification, before archive"** (:239) — still says **"Per-batch is the default."** REQ-035's plan intended to update both; its builder updated only the work.md copy. Since REQ-037 rewrites this exact section (its requirement #1 places the merge point in Worktree Dispatch Mode), it must **align work-reference.md's post-merge default to per-merge**, matching work.md, so the two files tell one story. This is in-scope for REQ-037 (same section, same write_set file), not a separate follow-up.

## Plan

### Design decisions (settled before edits)

**Merge point (pipeline position):** the orchestrator merges each builder branch **at hand-back, end of Step 6, before Step 6.25 (Implementation Summary)** — so every downstream evidence step (6.25/6.3/6.5/7/8/9) observes one integrated tree. Any position after 6.25 leaves qualify/review reading a clean main tree (the bug REQ-037 kills).

**Merge-range representation — `<pre>..HEAD`, `<pre>` captured pre-merge:** immediately before `git merge --no-ff <branch>`, in the main tree, record `pre=$(git rev-parse HEAD)` (integration tip before merge) and, right after the merge, `merge_hash=$(git rev-parse --short HEAD)` (the `--no-ff` merge commit → the REQ's `commit:` field). Range = `<pre>..HEAD`. Verified: `git diff <pre>..HEAD` is the endpoint diff; because `<pre>` is the merge commit's first parent, `merge-base(<pre>,HEAD)==<pre>`, so the "merge-base form" collapses to `<pre>..HEAD` (simplest correct form). Capture is **per REQ, per merge**, held in orchestrator memory through Steps 6.25–9 (same pattern as deferred prime-link tuples / lock heartbeats). Capture `<pre>` explicitly, never recover as `HEAD^1` — HEAD moves once the orchestrator makes Step 9 commits or a second merge.

**qualify.sh interface — env var `DO_WORK_DIFF_RANGE`** (not a positional, so `$1` stays the sole documented positional and the exit-2 usage contract is untouched). Unset/empty ⇒ serial default, byte-identical to today. Step 6.3 sets it only in worktree mode.

**Ratchet safety (verified on disk):** `_dev/tests/contract-regressions.sh:199-208` pins qualify.sh only via the exec bit + the `qualify.sh` reference in work.md — no ratchet pins any string *inside* qualify.sh, so the shell body is free to change. Prose ratchets preserved: `[Pp]ost-merge verification` (work.md — keep the bold lead-in), `worktree-agent-REQ-`/`git branch -d` (work-reference.md Naming/Cleanup — untouched), `pairwise disjoint`/`serial-only`/`claimed_reqs`/`this session's own co-dispatched claims` (work.md Step 1 — untouched), `claimed_reqs`/`including this session's own` (work-reference.md Crash Recovery — untouched).

### Per-file edits

- **A — `actions/work-reference.md` "Merge, never rebase" (:237):** becomes the single source of truth for merge point + range + capture (`pre`/`merge_hash`) + the consumer list (6.3 via `DO_WORK_DIFF_RANGE`, Step 7 review Get-the-Diff, Step 8 verification, Step 9 validation) + serial-unchanged note + an integration-seam note (commit any post-merge seam before Step 6.25 so it lands inside `<pre>..HEAD`).
- **B — `tools/checks/qualify.sh` (header + two diff spots):** header syncs the env-var doc + adds `Exit 2: usage error` (code already exits 2); read `diff_range="${DO_WORK_DIFF_RANGE:-}"` after the request_file guard; branch BOTH `changed_file_list` (~:44-46) and the debug-artifact scan (~:112) on `[ -n "$diff_range" ]` — ranged branch uses `git diff "$diff_range" ...`, the `else` is **byte-for-byte today's line**. Keep `:(exclude)do-work/`, the grep chain, `|| true`, and (edit in place) the executable bit. Extend the :40-43 and :63-65 comments: the range's lower bound `<pre>` preserves the "this-REQ-only" guarantee the working+staged default gives serially.
- **C — `actions/work-reference.md` "Post-merge verification, before archive" (:239):** replace "Per-batch is the default" with the per-merge-default wording matching work.md Step 8 (the folded-in gap); keep "the unit you verify is the unit you roll back" and "a red merged state stops the archive."
- **D — `actions/review-work.md` (3 spots):** Step 4 Get-the-Diff pipeline (:66) reads `git diff <pre>..HEAD` in worktree mode (working tree clean post-merge); Two-Modes table pipeline row (:30) notes the merge-range source; nothing-to-review exit (:46) judges "no changes" from the merge range in worktree mode (prevents silent Step 7 skip).
- **E — `actions/work.md` Step 6.3 (:430-434):** add the `DO_WORK_DIFF_RANGE="<pre>..HEAD" …/qualify.sh <req-file>` worktree variant + one sentence; serial omits it.
- **F — `actions/work.md` Step 8 post-merge para (:542):** reference the merge point ("merged at hand-back, end of Step 6 — see Worktree Dispatch Mode"); keep the `[Pp]ost-merge verification` bold lead-in (ratchet).
- **G — `actions/work.md` Step 9 (:594):** worktree-mode paragraph — implementation commit already exists (builder commits + `--no-ff` merge); Step 9 stages only changelog/version/archived-REQ/follow-ups/UR-moves/prime-links, writes `<merge_hash>` (not the changelog commit) into `commit:`, and validates the Implementation-Summary file list against `git diff --name-only <pre>..HEAD` instead of the stage.
- **H — `actions/work-reference.md` Commit & Metadata Procedure (:804/:806):** staging-list + validation-check worktree-mode branches mirroring G.
- **I — `actions/work.md` Step 7 dispatch (:479-481):** the review reads the merge range in worktree mode; pass `<pre>..HEAD` in the review agent's context.
- **J — `actions/work.md` Rules "one commit per request" (:691):** clause — in worktree mode the implementation's commit boundary is the merge (builder commits + `--no-ff` merge); the single-commit rule governs serial mode.

### Plan validation (orchestrator, session do-work-20260729T065754Z-5724)
Requirement coverage complete (all 4 Detailed Requirements + the folded per-batch/per-merge gap map to edits A–J + C); no orphan edits; ~10 sites but one coherence contract (range defined once in A, consumed everywhere) — **do not split** (splitting half-lands the contract, the exact REQ-035 failure mode). Serial byte-compatibility is structural: every change gated on worktree mode / `diff_range` being set, and qualify.sh's serial `else` branches are unchanged. **Carry Risk #3 forward:** writing a *merge* hash into `commit:` means standalone review's `git show <commit>` (review-work.md:68, out of REQ-037's pipeline-mode scope) would show a combined/often-empty merge diff — flag as a candidate follow-up, do not fix here (scope creep). Reviewer should also watch Risk #2 (post-merge integration seams outside the range — mitigated by A's commit-before-6.25 rule) and Risk #1 (preserve qualify.sh's exec bit).

*Generated by Plan agent; validated by work action (orchestrator)*

## Exploration

Folded into the Plan agent's pass (Route C) — it verified every anchor against disk at commit `4296e11` and confirmed git range semantics against real repo commits. Anchor line numbers are recorded inline in the per-file edit list above. Key confirmations: qualify.sh is pinned by the test only via exec-bit + work.md reference (no internal-string ratchet); `git diff <pre>..HEAD -- . ':(exclude)do-work/'` is valid and pipefail-safe; the four REQ-033 ratchets + REQ-032/035 ratchets all live outside the edited spans.

*Exploration folded into the Route C plan*

## Scope

**Files I will touch:**
- `actions/work.md` (modify) — Step 6.3 invocation (E), Step 7 dispatch (I), Step 8 post-merge reference (F), Step 9 worktree paragraph (G), Rules "one commit per request" clause (J)
- `actions/work-reference.md` (modify) — "Merge, never rebase" single-source-of-truth block (A), "Post-merge verification" per-merge alignment (C), Commit & Metadata Procedure staging/validation branches (H)
- `actions/review-work.md` (modify) — Step 4 Get-the-Diff pipeline (D), Two-Modes table row (D), nothing-to-review exit (D)
- `tools/checks/qualify.sh` (modify) — header docs + `diff_range` read + ranged branches on `changed_file_list` and the debug-artifact scan (B); **preserve the executable bit**

**Files I will NOT touch:** `actions/cleanup.md` (its `worktree-agent-` ratchet + Pass 5 are out of scope); `_dev/tests/contract-regressions.sh` (run for verification, add no ratchet — the existing exec-bit + `qualify.sh`-reference + `[Pp]ost-merge verification` ratchets already anchor this surface); `tools/queue-kanban/model.go` (no schema change)

**Acceptance criteria (restated from REQ):**
- [ ] Merge point stated once in Worktree Dispatch Mode as a pipeline position (end of Step 6, before Step 6.25); work.md Step 8 references it rather than assuming it
- [ ] Merge range defined once (`<pre>..HEAD`, `<pre>` captured pre-merge per REQ) and consumed by Step 6.3 (`qualify.sh` via `DO_WORK_DIFF_RANGE`), Step 7 (review Get-the-Diff), and Step 9's staged-list validation; serial mode unchanged (working-diff default)
- [ ] Step 9 reconciled for merged work: no staging of already-committed implementation files; job = changelog/version + write the merge commit's hash; the "one commit per request" Rules bullet says so
- [ ] Folded gap: work-reference.md "Post-merge verification" default aligned to per-merge, matching work.md
- [ ] All four REQ-033 ratchets green; serial default byte-compatible; qualify.sh exit-code contract + header docs in sync; exec bit preserved
- [ ] `bash _dev/tests/contract-regressions.sh` green

*Scope declared by work action (orchestrator, session do-work-20260729T065754Z-5724)*

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]** Implement edits A–J exactly per the validated Plan. Range `<pre>..HEAD` defined once in work-reference.md "Merge, never rebase" (A) and consumed by qualify.sh via `DO_WORK_DIFF_RANGE` (B/E), review Get-the-Diff (D/I), Step 8 verification (F), and Step 9 staging/validation (G/H). Serial path stays byte-identical: every qualify.sh change branches on `[ -n "$diff_range" ]` with the `else` = today's command; every action-file change gated on "worktree dispatch mode". Fold the per-batch→per-merge gap in work-reference.md (C). Preserve qualify.sh's exec bit (edit in place). Add no ratchet.
- [x] **[APPLY]** Edits A–J applied across the four write_set files: A/C/H in work-reference.md; B in qualify.sh (header + `diff_range` read + two ranged branches + two extended comments); D in review-work.md (3 spots); E/F/G/I/J in work.md.
- [x] **[UNIFY]** `git diff --stat` = exactly the 4 write_set files (review-work.md, work-reference.md, work.md, qualify.sh). `bash _dev/tests/contract-regressions.sh` → exit 0. qualify.sh exec bit intact (`-rwxr-xr-x`, index `100755..100755`, no mode change). Serial `else` branches byte-identical to pre-edit (verified stripped). Both serial and ranged runs exit with the same contract; ranged branch exercised end-to-end (found a range-only file, exit 0). No debug artifacts in the diff.

## Decisions

- **D-01 — qualify.sh serial `else` re-indented, command bytes preserved (DECIDE & STATE).** The non-negotiable "byte-for-byte identical" constraint on the unset-range branches: wrapping each line in `if [ -n "$diff_range" ]; then … else … fi` forces +2 spaces of leading indentation, which is unavoidable and behavior-neutral. I kept the command text (git invocation, pipe chain, `:(exclude)do-work/` pathspec, grep chain, `|| true`) character-identical and verified it by stripping leading whitespace and comparing to `HEAD:tools/checks/qualify.sh` — both `changed_file_list` and `debug_artifact_lines` else-branches compare IDENTICAL. Reversible, no reasonable objection.
- **D-02 — new merge-point/range block is a dedicated paragraph, not folded into "Merge, never rebase" (DECIDE & STATE).** Plan Edit A says that section "becomes the single source of truth." I left the existing "Merge, never rebase." paragraph intact (it is about the merge *mechanic*) and added an adjacent bold-led paragraph ("When to merge, and the range every evidence step reads.") for the merge point + range + capture + consumer list + serial note + integration-seam note. Cleaner than a single overloaded paragraph; keeps the ratcheted `git merge --no-ff`/`git branch -d` prose untouched.

## Discovered Tasks

- **[normal]** Standalone review's `git show <commit>` shows an empty/combined diff for a merge commit (Plan Risk #3, out of REQ-037's pipeline-mode scope). REQ-037 makes worktree-mode Step 9 write the `--no-ff` **merge commit** hash into `commit:`. `actions/review-work.md` **standalone** mode (Step 4, "Get the Diff") runs `git show <commit>` off that field; for a merge commit `git show` prints a combined diff (often empty), so a later standalone re-review of a worktree-merged REQ would see nothing to review. Fix candidate: in standalone mode, detect a merge commit and use `git show --first-parent -m <commit>` (or `git diff <commit>^1..<commit>`). Not fixed here — the REQ is pipeline-mode-scoped and review-work.md's standalone Step 4 is explicitly out of scope.
- **[low]** `tools/checks/qualify.sh`'s debug-artifact scan false-positives on any diff that *adds its own detection tokens* — surfaced by this REQ, because the new ranged branch (Edit B) necessarily duplicates the grep chain, so the added lines literally contain `console.log|debugger|…|TODO|FIXME` as the pattern text. With `[UNIFY]` checked, the scan matched its own regex and emitted a FAIL for zero real artifacts (orchestrator overrode it — see `## Qualification`). The blind spot is general: it recurs whenever a diff edits qualify.sh's pattern, or edits any file that legitimately contains those tokens as data (a linter, a doc about debugging). Fix candidates (all imperfect): exclude qualify.sh itself from its own scan; require the token to appear in added *code* rather than inside a quoted grep pattern; or scope the scan to specific file types. Narrow trigger, and orchestrator judgment catches it, so low priority.
- **[low]** The mid-run blocked-flip guard (`actions/work.md` Step 8, "no substantive implementation edits landed this attempt") tests `git diff -- . ':(exclude)do-work/'` against the **main tree** — which in worktree dispatch mode is clean even when the builder committed edits on its own branch. So a worktree builder that hit a missing external precondition *after* committing real work could be mis-flipped to `blocked` (non-terminal) instead of failing. Surfaced by REQ-037's adversarial review (hunter, explicitly logged as out-of-scope-but-noted). Out of REQ-037's declared scope (that scope is the *success-path* evidence steps — qualify/review/Step 9; the blocked-flip is the failure/blocked path) and pre-existing (predates this REQ). Fix candidate: in worktree mode the blocked-flip's "did edits land" test must consult the builder's branch (or the hand-back manifest), not the main-tree diff.

## Implementation Summary

**Files changed:**
- `actions/work.md` (modified)
- `actions/work-reference.md` (modified)
- `actions/review-work.md` (modified)
- `tools/checks/qualify.sh` (modified)

**What was done:** Placed the worktree merge in the pipeline (orchestrator merges at hand-back — end of Step 6, before Step 6.25) and made every evidence-consuming step read the merged diff in worktree mode. `actions/work-reference.md`'s Worktree Dispatch Mode gained a single source-of-truth block defining the merge point, the `pre`/`merge_hash` capture, and the merge range `<pre>..HEAD` (with its merge-base-collapse rationale, the hold-in-memory-never-`HEAD^1` rule, and the consumer list); its "Post-merge verification" default was aligned to per-merge (matching work.md — the folded REQ-035 gap); and its Commit & Metadata Procedure gained worktree-mode staging + validation branches. `tools/checks/qualify.sh` reads an optional `DO_WORK_DIFF_RANGE` and branches both the file-list check and the debug-artifact scan on it — the unset (serial) branch is byte-identical to before, the exec bit and exit-code contract are unchanged, and the header now documents the env var plus `Exit 2`. `actions/review-work.md` reads `<pre>..HEAD` in worktree pipeline mode (Two-Modes row, nothing-to-review exit, Step 4 Get-the-Diff); standalone `git show <commit>` left unchanged (recorded as a Discovered Task). `actions/work.md` Step 6.3, Step 7, Step 8, Step 9, and the "one commit per request" Rules bullet all reference the merge point/range and reconcile Step 9 for already-committed merged work (stage only changelog/version/metadata; write the merge commit's hash). Serial/floor path behaviorally unchanged; no ratchet added; no schema change.

*Summary written by work action (orchestrator)*

## Qualification

**Passed (with a documented mechanical-check override).** `scope-drift.sh` clean (4 declared = 4 touched); `git diff` = exactly the 4 write_set files; qualify.sh exec bit intact (no mode change); `bash -n tools/checks/qualify.sh` valid; both serial and ranged qualify paths exercised and exit with the correct contract.

**qualify.sh mechanical FAIL — verified false positive, overridden.** Run serially against this REQ, `qualify.sh` emitted `FAIL: [UNIFY] is checked but the diff adds debug artifacts` on two lines. The orchestrator enumerated every match of the exact scan (`{ git diff; git diff --staged; } | grep '^\+' | grep -v '^\+\+\+ ' | grep -nE 'console\.log|debugger|print\(|TODO|FIXME'`): **exactly 2 matches, both in `tools/checks/qualify.sh`, both the `debug_artifact_lines=…` grep-chain lines that contain qualify.sh's own detection-pattern literal**; the three action files have **0** matches. There are no real debug artifacts — the scan matched its own regex because Edit B necessarily duplicates the grep chain into the new ranged branch. This is a self-referential blind spot in the check, not a defect in the implementation, so the FAIL is overridden by judgment (Step 6.3's mechanical checks feed judgment; the anti-rationalization rule guards against explaining away *real* problems — this match set was verified to contain none). Recorded as a `[low]` Discovered Task above so the blind spot is on the record. Judgment checks 2/3/6: substantive spec + shell edits; all 4 acceptance criteria + the folded gap trace to diff changes; check 6 N/A (the qualify.sh change is gated logic, not a data path).

*Verified by work action (orchestrator)*

## Testing

**Tests run:** `bash _dev/tests/contract-regressions.sh`; plus targeted qualify.sh behavior probes (serial default vs `DO_WORK_DIFF_RANGE` set).
**Result:** ✓ Contract regression checks passed (exit 0) — all ratchets green (exec-bit + `qualify.sh` reference + `[Pp]ost-merge verification` + the REQ-032/033/035 prose anchors), none added. Serial `else` branches of qualify.sh confirmed byte-identical to HEAD (modulo forced indentation). The new ranged branch was exercised end-to-end against a real range and behaves per the exit-code contract.

**Red-green validation:** N/A — non-behavioral change: Markdown instruction prose + a gated shell branch whose serial default is byte-identical. Proof is the green regression suite, the byte-identical-else verification, and the ranged-path exercise.

*Verified by work action*

## Review

**Overall: 92%** | 2026-07-29T08:33:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 95% |
| Code Quality | 88% |
| Test Adequacy | 85% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Findings:** 2 raw Important → **0 confirmed** (both refuted by both refuters); 3 Minor; 1 Nit. Full adversarial rigor (reviewer + independent contradiction-hunter + 2 lens-diverse refuters per Important).
**Acceptance:** Pass — `contract-regressions.sh` green; qualify.sh serial `else` branches verified byte-identical to HEAD; exec bit + exit-code contract intact; all cross-reference anchors resolve; the folded per-batch/per-merge gap is closed.
**Refuted Important (report-only, not folded):** the review Red Flag "Implementation Summary lists files but `git diff` shows no changes" (review-work.md) was flagged as a 4th un-repointed working-diff consumer, but **both refuters killed it**: it is a post-acquisition symptom heuristic consuming the diff already obtained in the re-pointed Step 4 (not an independent acquisition site), and the identical "git diff" wording has coexisted with standalone mode's clean working tree without misfiring, because it means "the diff under review," not a literal command. Respected the refutation.
**Out-of-scope observation (recorded, deferred):** the hunter noted the Step 8 blocked-flip guard also reads the main-tree diff — real, but the failure/blocked path, pre-existing, outside REQ-037's success-path scope. Recorded as a `[low]` Discovered Task → follow-up.
**Follow-ups created:** REQ-042 (`pending-answers`) bundling the three Discovered Tasks (standalone-review merge-hash; qualify.sh self-scan blind spot; blocked-flip worktree guard) for your decision via `do-work clarify`.

*Reviewed by review-work action (pipeline mode, full adversarial rigor per the session's calibration)*

## Remediation

The adversarial pass confirmed no Important defects, but two of the Minors were coherence/correctness gaps in REQ-037's *own* new design, in-scope and cheap, so they were folded in (rather than left report-only):

- **Minor (correctness) — merge-range upper bound `HEAD` → captured `<merge_hash>`.** Both finders independently flagged that the range was written `<pre>..HEAD` while the plan had *captured* `merge_hash` precisely because "HEAD moves once the orchestrator commits the changelog or merges the next sibling." In **per-batch** verification (a permitted mode) sibling merges land before this REQ is verified, so `<pre>..HEAD` would sweep in unrelated commits and misattribute them. Replaced every `<pre>..HEAD` with the stable `<pre>..<merge_hash>` across all four files (12 occurrences) and rewrote the merge-base rationale to match (`merge-base(<pre>, <merge_hash>) == <pre>`), keeping one intentional `<pre>..HEAD` mention that *explains the hazard* of the HEAD form. This makes the range definition consistent with its own captured-anchor rationale.
- **Nit — qualify.sh WARN text.** The two WARN strings said "working/staged diff" unconditionally; appended "(or the merge range, in worktree mode)" so the message is accurate in ranged mode (the polish the Plan's Edit B flagged as optional and the builder skipped).

**Report-only (not folded):** the qualify.sh self-scan Minor is the same self-referential false positive already documented in `## Qualification` and recorded as a `[low]` Discovered Task — its fix candidates are non-trivial and out of scope; the refuted Red-Flag Important stands refuted.

**Re-verification:** `contract-regressions.sh` exit 0; qualify.sh `bash -n` clean, exec bit intact, no mode change; `scope-drift.sh` clean (still the 4 declared files); grep confirms no unintended `<pre>..HEAD` remains.

*Remediated by work action (orchestrator)*

## Lessons Learned

**What worked:** Defining the merge range once (single source of truth in Worktree Dispatch Mode) and pointing every consumer at it made the re-pointing auditable — the hunter could grep every worktree/range site and confirm the consumer list matched the re-pointed sites. Gating every change on "worktree dispatch mode" / `[ -n "$diff_range" ]` kept the serial floor path provably byte-identical.

**What didn't:** The plan captured a stable anchor (`merge_hash`) for exactly the "HEAD moves" reason, then defined the range with live `HEAD` anyway — a self-inconsistency that only bites in per-batch mode, which the requirements-walk missed and two independent finders caught. Lesson: when you capture a stable handle *because* a live ref moves, use that handle at every reference, not just the one that motivated capturing it.

**Worth knowing:** A self-check that greps a diff for `console.log|TODO|FIXME|debugger` cannot cleanly qualify a change to *its own* detection pattern (or to any file that carries those tokens as data) — qualify.sh false-positived on this REQ. Recorded as a `[low]` follow-up; until fixed, the orchestrator must eyeball such a FAIL rather than bounce the builder.

## Orientation

Worktree dispatch mode now has a **placed merge and a re-pointed evidence path**: the orchestrator merges each builder branch at hand-back (end of Step 6, before Step 6.25), captures the range `<pre>..<merge_hash>` around that merge, and every evidence-consuming step (Step 6.3 `qualify.sh` via `DO_WORK_DIFF_RANGE`, Step 7 review, Step 8 verification, Step 9 staged-list validation) reads the merge range instead of the post-merge-clean working tree. Step 9 is reconciled for already-committed merged work (stage only changelog/version/metadata; `commit:` gets the merge hash). Lives in `actions/work-reference.md` (Worktree Dispatch Mode + Commit procedure), `actions/work.md` (Steps 6.3/7/8/9 + Rules), `actions/review-work.md` (Step 4 + Two-Modes + nothing-to-review), and `tools/checks/qualify.sh` (the `DO_WORK_DIFF_RANGE` branch). **[MAP CHANGED]** `qualify.sh` gained an input contract (`DO_WORK_DIFF_RANGE`), and worktree-mode REQs now carry a **merge commit** hash in `commit:` — a new downstream coupling that REQ-042 (standalone review) tracks. Serial mode is byte-unchanged; this completes the worktree-mode evidence story that REQ-033 opened and REQ-035/036 hardened.
