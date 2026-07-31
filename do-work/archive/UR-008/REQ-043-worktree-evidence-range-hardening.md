---
id: REQ-043
title: "Worktree evidence-range hardening: qualify range validation, seam-in-range, remediation re-merge, Step 6 merge instruction, commit-procedure carve-out"
status: completed
route: C
created_at: 2026-07-29T09:30:45Z
claimed_at: 2026-07-29T10:06:57Z
completed_at: 2026-07-29T12:36:22Z
commit: 1ac4fa8
user_request: UR-008
addendum_to: REQ-037
domain: general
prime_files: []
tdd: false
depends_on: []
related: [REQ-044, REQ-045, REQ-046]
batch: deep-review-followups
write_set: [actions/work.md, actions/work-reference.md, tools/checks/qualify.sh]
maintenance: true
---

# Worktree Evidence-Range Hardening (REQ-037 Follow-Up)

## What

REQ-037 placed the worktree merge and defined the evidence range `<pre>..<merge_hash>`, but five confirmed defects remain in how the range is produced, validated, and restated. Fix them as one coherence contract — the same reasoning REQ-037's plan used for not splitting ("splitting half-lands the contract").

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Written below as `## Plan` — six requirements, one re-definition of the two range endpoints carries items 3 and 4 together. Verified both git mechanisms in a scratch repo before prescribing them.
- [x] **[APPLY]:** All three write_set files edited; no file outside `write_set` touched.
- [x] **[UNIFY]:** `git diff --stat` = 3 files, +87/-17, exactly the `write_set`. Verified `tools/checks/qualify.sh` (`bash -n` clean; `shellcheck` reports only the pre-existing SC2016 info on the `sed` extraction at line 116, none of my lines; executed the full RED/GREEN range matrix — 5 unresolvable ranges FAIL with exit 1, valid range and both serial forms still exit 0 with identical output); `actions/work-reference.md` (read all four changed hunks in place; terminology normalized to the single incumbent term "merge range"; no dangling pointer — every cross-reference resolves to a section that exists); `actions/work.md` (read all three changed hunks; new Step 6 block sits at the end of Step 6 immediately before `### Step 6.25`). `_dev/tests/contract-regressions.sh` passes. No debug artifacts: grepped the added lines of all three files for `console.log|debugger|print(|TODO|FIXME` — none.

## Prior Implementation

REQ-037 (archived, commit `1348d11`, v0.146.0) established: merge at hand-back (end of Step 6, before Step 6.25); range `<pre>..<merge_hash>` with both endpoints captured around the merge and held in orchestrator memory; consumers re-pointed in `actions/work.md` (Steps 6.3/7/8/9), `actions/work-reference.md` (Worktree Dispatch Mode + Commit procedure), `actions/review-work.md`, and `tools/checks/qualify.sh` (`DO_WORK_DIFF_RANGE` branch). Serial mode byte-unchanged.

## Detailed Requirements

1. **qualify.sh must fail loudly on an unresolvable range** (verified by execution). With `DO_WORK_DIFF_RANGE="bogus..range"`, `git diff` exits 128 with empty stdout; the script (`set -uo pipefail`, no `-e`) continues and prints `OK: mechanical qualification passed`, exit 0 — the debug-artifact check is silently disabled. Validate the range resolves (e.g. `git rev-parse` both endpoints or check `git diff`'s exit status) before any check consumes it; an unresolvable range is a FAIL naming the range, never OK.
2. **Stop prescribing the range as shell variables that die between command blocks.** `actions/work-reference.md`'s capture is written as `pre=$(git rev-parse HEAD)` / `merge_hash=$(...)` while consumers sit in later blocks — the repo's own "shell state does not survive between prescribed blocks" trap. An agent writing `"$pre..$merge_hash"` in a fresh shell produces `".."` → fatal → (per defect 1) a false OK. Prescribe holding both values as re-typed literals (the same pattern the lock's `session_id` uses), and make the hold-in-memory instruction name `<merge_hash>` as well as `<pre>`.
3. **The integration-seam sentence is arithmetically false.** `actions/work-reference.md` → Worktree Dispatch Mode says commit the seam before Step 6.25 "so it lands inside `<pre>..<merge_hash>`" — but `<merge_hash>` is captured immediately after the merge, so a post-merge seam commit is its *child* and outside the range; qualify, review, and Step 9 validation never see the seam. Fix: capture `<merge_hash>` after the seam commit (re-defining it as the evidence tip, and reconciling the "`<pre>` is `<merge_hash>`'s first parent" claim), or an equivalent that provably puts seams inside the evidence range.
4. **Define the range across remediation re-merges.** Re-qualification/remediation loops merge the builder's fix branch again; `pre₂..merge_hash₂` covers only the fix delta, so review reads only the fix and every originally-touched file WARNs. Specify the cumulative range (`pre₁..merge_hash₂` — first pre-merge tip to latest merge) and which hash Step 9 records.
5. **Step 6 must imperatively instruct the orchestrator to merge and capture `<pre>`.** In `actions/work.md`, the merge's first mention is retroactive at Step 6.3 ("already committed and merged (end of Step 6)"); Step 6's only worktree bullet is builder-facing. Add the orchestrator-side hand-back instruction (merge `--no-ff`, capture `<pre>`/`<merge_hash>`, apply+commit seams) at the end of Step 6, pointing at the reference for detail.
6. **Remove the commit-procedure self-contradiction.** `actions/work-reference.md`'s worktree paragraph says `commit:` gets `<merge_hash>`, but the generic "Write commit hash back" paragraph + bash block just below still prescribe `HASH=$(git rev-parse --short HEAD)` after the changelog commit and call it "the real implementation commit hash" with no worktree carve-out. Add the carve-out so the copy-pasteable block cannot write the changelog commit's hash in worktree mode.

## Constraints

- Serial mode must remain byte-compatible: every change gated on worktree mode / `DO_WORK_DIFF_RANGE` being set; qualify.sh's serial branches untouched.
- Grep the same primitives repo-wide before calling any item fixed (per the repo-wide restatement rule) — REQ-037's own review missed restatements this way.
- Do not touch the standalone `git show <commit>` consumers — that is REQ-055 (extended by addendum to cover present-work/pipeline/ai-report).

## Red-Green Proof
**RED prompt/case:** In a worktree-mode run: `DO_WORK_DIFF_RANGE="bogus..range" tools/checks/qualify.sh <req-file>` prints a git fatal yet reports `OK: mechanical qualification passed` (exit 0). Reading `actions/work.md` Step 6 top-to-bottom, an orchestrator gets no merge instruction until Step 6.3 mentions it retroactively. A seam committed after the merge is absent from `git diff <pre>..<merge_hash>` despite the reference claiming it "lands inside" the range.
**Why RED now:** Range validation, block-survival, seam arithmetic, re-merge semantics, the Step 6 imperative, and the commit-procedure carve-out were all outside REQ-037's executed scope or shipped self-contradicting.
**GREEN when:** qualify.sh FAILs (non-zero/explicit FAIL line) on any unresolvable range; the prescribed capture survives block boundaries by design; a post-merge seam commit is provably inside the evidence range; remediation re-merges have a defined cumulative range; Step 6 contains the orchestrator-side merge instruction; the commit-procedure block carries the worktree carve-out.
**Validation:** User confirmed (approved capture of the reviewed finding set)

## Full Context
See `do-work/user-requests/UR-008/input.md`.

## Triage

**Route:** C (Full Pipeline)
**Reasoning:** A single coherence contract spanning six confirmed defects across three files (`actions/work.md`, `actions/work-reference.md`, `tools/checks/qualify.sh`), where the requirements themselves warn that a split half-lands the contract and that restatements were missed by the parent REQ's own review. Multi-file, semantically-coupled, self-referential spec surface → Route C.
**Complexity indicators:** 6 interlocking requirements; a repo-wide restatement-grep constraint; an execution-verified behavior change in `qualify.sh`; `maintenance: true` (deliberate pass on the skill's own instructions — delete/narrow discipline applies).
**Rigor:** Full adversarial review at Step 7 (reviewer + independent contradiction-hunter + lens-diverse refuters) — this is the concurrency/worktree-evidence surface where prior full-adversarial passes confirmed real coherence defects.

*Triaged 2026-07-29 by orchestrator (session do-work-20260729T100657Z-34626).*

## Plan

The six requirements collapse into **one re-definition plus four local fixes.** Requirements 3 and 4 both ask "what exactly are the two endpoints?", so they are answered together by pinning the endpoints' capture points rather than by adding two independent caveats:

- **`<pre>` = the integration tip before this REQ's FIRST merge, captured once and never re-captured.**
- **`<merge_hash>` = the latest `--no-ff` merge commit for this REQ, which itself contains the integration seam.**

Under that definition the range `<pre>..<merge_hash>` is automatically cumulative across remediation re-merges (`pre₁..merge_hash₂`) and automatically contains the seam — and, critically, every existing restatement of the two token names stays true verbatim, including the three in `actions/review-work.md` which is outside this REQ's `write_set`. See D-01 for the seam mechanism and why the REQ's literally-suggested fix was rejected.

Order of work and file mapping:

1. **`tools/checks/qualify.sh` (req 1)** — insert a range-validation gate after the `git_available` probe and before `changed_file_list` is computed, wholly inside `if [ -n "$diff_range" ]`. Structural check (two-dot), then each endpoint non-empty, then each endpoint resolves via `git rev-parse --verify --quiet <x>^{commit}`, then a final `git diff --name-only <range>` probe. Any failure ⇒ one `FAIL:` line naming the range and the specific defect, plus a second line saying why it is not a warning, then `exit 1`. The per-endpoint diagnostic is the direct read-out for requirement 2's failure mode (a `..` from a lost hash).
2. **`actions/work-reference.md` (reqs 2, 3, 4, 6)** — rewrite the "When to merge, and the range every evidence step reads" paragraph into a 4-step hand-back sequence + three short bolded rules (literal-holding, the range and its consumers, cumulative re-merges); update **Sole integrator** so the seam is applied inside the merge commit; add the worktree carve-out to the "Write commit hash back" paragraph and its bash block.
3. **`actions/work.md` (reqs 5, 4, 6-adjacent)** — add the imperative orchestrator hand-back merge block at the end of Step 6; point Step 6.3's range prose at both held literals and at the new hard FAIL; add "latest merge / never `rev-parse HEAD`" to Step 9's worktree paragraph.
4. **Verify by execution**, not by reading: both git mechanisms in a scratch repo *before* prescribing them, then the RED/GREEN range matrix against the edited script, then `_dev/tests/contract-regressions.sh`.

**Seam-arithmetic decision (the load-bearing choice):** fold the seam into the merge commit via `git merge --no-ff --no-commit` → apply seam → `git commit`, rather than committing the seam after the merge and re-pointing `<merge_hash>` at it. Full reasoning in D-01.

## Exploration

Restatement sweep — every grep run, and every site found. Ripgrep over the whole repo excluding `do-work/**` and `.git/**`.

**Tokens grepped:** `merge_hash`, `DO_WORK_DIFF_RANGE`, `<pre>`, `lands inside` / `land inside`, `first parent` / `first-parent`, `HEAD^1`, `rev-parse --short HEAD`, `real implementation commit`, `per-batch`, `per-merge`, `merge range`, `evidence range`, `pre-merge`, `integration tip`, `integration seam`, `seam`, `no-ff`, `worktree dispatch`, `after the merge`.

**Sites found and FIXED (all inside `write_set`):**

| Site | What it restated | Action |
|---|---|---|
| `tools/checks/qualify.sh:9-13` (header `Env:` block) | `DO_WORK_DIFF_RANGE` contract | Added the hard-FAIL clause |
| `tools/checks/qualify.sh:95-101` | `<pre>` = "tip just before this REQ's merge" | Re-stated as FIRST merge; added that `<merge_hash>` carries the seam; added pointer to the canonical definition |
| `tools/checks/qualify.sh:128-131` | same `<pre>` claim, inside the `(modified)` branch | Synced to "FIRST merge" |
| `actions/work-reference.md:235` (**Sole integrator**) | "the orchestrator applies it in the main tree **after the merge**" — the upstream source of the false seam arithmetic | Changed to "inside the merge commit", pointing at hand-back step 3 |
| `actions/work-reference.md:239` (**When to merge…**) | the canonical range definition, the `pre=$(...)`/`merge_hash=$(...)` shell-variable capture, the "first parent" claim, the "lands inside `<pre>..<merge_hash>`" seam sentence | Fully rewritten (reqs 2, 3, 4) |
| `actions/work-reference.md:808` (commit procedure, worktree paragraph) | "`commit:` gets `merge_hash`, captured in Step 6" | Added "the **latest** merge if remediation re-merged"; token written as `<merge_hash>` to match every other site |
| `actions/work-reference.md:812-823` (**Write commit hash back**) | `HASH=$(git rev-parse --short HEAD)` + "the real implementation commit hash", no worktree carve-out | Carve-out added in prose, in the bash block, and in the closing sentence (req 6) |
| `actions/work.md:440` (Step 6.3) | "`<pre>` is this REQ's pre-merge integration tip" | Both endpoints named; re-typed-literal rule; the new hard FAIL surfaced at the call site |
| `actions/work.md:604` (Step 9) | "the `--no-ff` merge commit captured in Step 6, not this changelog commit's hash" | Added "latest one if remediation re-merged" and an explicit "never `git rev-parse --short HEAD`" |
| `actions/work.md` end of Step 6 | *nothing existed* — the gap requirement 5 names | New imperative hand-back block |

**Sites found and deliberately NOT touched:**

- `actions/review-work.md:30`, `:46`, `:66` — three restatements of `<pre>..<merge_hash>` and "the merge range the orchestrator passes." **Outside `write_set`, and correct without change**: they name the two tokens and consume whatever the orchestrator passes. Because the re-definition kept both token *names* and only pinned their capture points, these stay true under both the seam fold and the cumulative re-merge range. This is the main reason D-01 rejected renaming or adding a third endpoint token — it would have forced an edit here.
- `actions/work.md:400` (builder-facing integration-seam bullet) — says hand the seam line back rather than editing the shared file. Unaffected by *where* the orchestrator then commits it.
- `actions/work.md:550` + `actions/work-reference.md:241` (post-merge verification, per-merge vs per-batch) — "merged at hand-back, end of Step 6, before Step 6.25" is still exactly true.
- `actions/work.md:701`, `:703` (`## Rules`) — "the implementation's commit boundary is the merge" becomes *more* true once the seam is inside the merge commit.
- `actions/cleanup.md:119` (Pass 5), `crew-members/background-agents.md:202-203` (worktree isolation as an axis) — reference the mode, not the range.
- `CHANGELOG.md:33-36`, `CHANGELOG-archive.md` — shipped history; not rewritten.
- `actions/pipeline-reference.md:249-250` — `<pre>` there is the HTML tag, an unrelated homograph.
- **All standalone `git show <commit>` consumers** — explicitly REQ-055's surface per this REQ's Constraints. Not read, not edited.

**Terminology check (found during the sweep, not in the requirements):** my first draft coined "evidence range" alongside the incumbent "merge range", which would have left two names for one concept — and the incumbent is the one used in the un-editable `actions/review-work.md`. Normalized every new occurrence back to **merge range**; `rg 'evidence range' actions/ tools/` is now empty.

**Verified in a scratch repo before prescribing (`/private/tmp/.../scratchpad/{mergetest,remergetest}`):**

- `git merge --no-ff --no-commit <branch>` on a fast-forwardable branch stops before committing and sets `MERGE_HEAD`; a subsequent `git commit` produces a real 2-parent merge commit whose **first parent is `<pre>`**; the seam file staged in between **is inside** `git diff --name-only <pre>..<merge_hash>`; and `git branch -d` still recognizes the branch as merged (the free merged-ness assertion survives).
- Re-merge: `merge-base(pre₁, merge_hash₂) == pre₁`, and `pre₁..merge_hash₂` carries both the original and the remediation delta.
- `Already up to date.` exits 0 and sets **no** `MERGE_HEAD` — so a blind `git commit` after it would fail or fabricate a non-merge commit. This is why the prescribed sequence carries an explicit guard for that answer.
- Honest cost of the cumulative range, measured not assumed: an unrelated orchestrator commit landing between the two merges **does** fall inside `pre₁..merge_hash₂`. Documented as over-inclusion with the direction-of-error argument rather than hidden.

## Scope

**Files I will touch:**

- `tools/checks/qualify.sh`
- `actions/work-reference.md`
- `actions/work.md`

Exactly the three `write_set` paths; nothing outside them was needed.

**Acceptance criteria (restated from Detailed Requirements + Red-Green Proof):**

1. `tools/checks/qualify.sh` FAILs — non-zero exit **and** an explicit `FAIL` line naming the range — on any unresolvable `DO_WORK_DIFF_RANGE`, validated before any check consumes the range. Never `OK` on a range git cannot read.
2. The prescribed capture holds `<pre>` **and** `<merge_hash>` as re-typed literals, not shell variables; the hold-in-memory instruction names both.
3. A seam is **provably** inside `<pre>..<merge_hash>` (verified by execution, not asserted), and the "`<pre>` is `<merge_hash>`'s first parent" claim is reconciled with the re-merge case.
4. Remediation re-merges have a defined cumulative range `pre₁..merge_hash₂`, and the hash Step 9 records is stated.
5. `actions/work.md` Step 6 ends with an imperative orchestrator-side merge instruction (merge `--no-ff`, capture both endpoints, apply+commit seams) pointing at the reference for detail — no longer first mentioned retroactively at Step 6.3.
6. The generic "Write commit hash back" block carries a worktree carve-out, so it cannot write the changelog commit's hash into `commit:`.
7. **Serial mode byte-compatible**: every change gated on worktree mode / `DO_WORK_DIFF_RANGE` being set; `qualify.sh`'s serial branches untouched.

## Implementation Summary

**Files changed**

- `actions/work.md` (modified) — added the imperative orchestrator hand-back merge block (capture `<pre>` → `git merge --no-ff --no-commit` → apply seams → `git commit` → capture `<merge_hash>`) at the end of Step 6; Step 6.3's range prose now names both held literals and the script's new hard FAIL; Step 9's worktree paragraph now says "latest merge if remediation re-merged" and explicitly forbids `git rev-parse --short HEAD` there.
- `actions/work-reference.md` (modified) — rewrote the merge/range paragraph in **Worktree Dispatch Mode (Step 1)** as a 4-step hand-back sequence plus three rules (re-typed literals not shell variables; the range and its four consumers; cumulative range across remediation re-merges with the hash Step 9 records); **Sole integrator** now applies seams inside the merge commit; deleted the arithmetically-false "commit the seam before Step 6.25 so it lands inside `<pre>..<merge_hash>`" sentence; generalized "first parent" to "first-parent ancestor"; added the worktree carve-out to **Commit & Metadata-Commit Procedure (Step 9)**'s hash write-back prose, bash block, and closing sentence, replacing the un-survivable `HASH=$(...)` variable with a re-typed literal.
- `tools/checks/qualify.sh` (modified) — added a range-validation gate (two-dot structure, both endpoints non-empty, both resolve as commits, final `git diff` probe) that exits 1 with a `FAIL` line naming the range and the specific defect; synced the three comments restating what `<pre>`/`<merge_hash>` mean; documented the hard FAIL in the header `Env:` block. Entirely inside `if [ -n "$diff_range" ]` — no serial branch touched.

**What was done**

Closed the six-defect coherence contract REQ-037 left open, as one change. The centrepiece is a re-definition of the two range endpoints — `<pre>` captured once before the *first* merge, `<merge_hash>` re-captured after the *latest* merge and containing the integration seam — which makes the seam provably in-range and the range automatically cumulative across remediation re-merges, without renaming a token or adding a third one (so the three restatements in `actions/review-work.md`, outside this REQ's write-set, stay correct untouched). The seam is folded into the merge commit with `git merge --no-ff --no-commit`, verified in a scratch repo to still yield a 2-parent merge commit with `<pre>` as first parent and to keep `git branch -d`'s merged-ness assertion intact. `qualify.sh` no longer prints `OK` when git could not read the range: it exits 1 naming the range and, when an endpoint is empty, says outright that the hash was lost between command blocks — the exact diagnostic for the shell-variable defect the same REQ fixes in the prose. Serial mode is byte-compatible: every script change lives inside the `DO_WORK_DIFF_RANGE`-set branch, confirmed by running the unset and empty-string paths and getting identical output to before.

## Decisions

**D-01 — Fold the integration seam into the merge commit instead of re-pointing `<merge_hash>` at a post-merge seam commit.** DECIDE & STATE.

Requirement 3 literally suggested "capture `<merge_hash>` after the seam commit (re-defining it as the evidence tip)". I rejected that reading, because `<merge_hash>` has a second job: Step 9 writes it into `commit:` as the implementation's provenance record, and `git show` consumers read that field. Pointing `commit:` at a one-line seam commit would make the REQ's recorded commit show the wiring instead of the implementation — trading a range defect for a provenance defect. Keeping both jobs on one hash would otherwise have required a third token (`<evidence_tip>`), which in turn would have forced edits to `actions/review-work.md` — outside this REQ's `write_set`.

Instead the merge is prescribed as `git merge --no-ff --no-commit` → apply seams → `git commit`. The seam is then *part of* the merge commit, so one hash is simultaneously the merge commit (`commit:`), the range's upper bound, and `<pre>`'s child with `<pre>` as its **first parent** — which is why the reconciled first-parent claim is stronger, not weaker, than before. Verified by execution in a scratch repo: 2-parent merge commit, seam present in `git diff <pre>..<merge_hash>`, `git branch -d` still deletes the branch as merged. Reversible: it is prose describing a two-command sequence.

**D-02 — Capture `<pre>` with `git rev-parse --short HEAD`, not `git rev-parse HEAD`.** DECIDE & STATE. The endpoints must now be *hand-re-typed* between command blocks, and a 40-character hash re-typed by hand is a transcription error waiting to happen; git lengthens `--short` automatically when a prefix would be ambiguous, so there is no correctness cost. This also makes both endpoints symmetric — `<merge_hash>` was already `--short` because `commit:` wants the short form.

**D-03 — Delete the `HASH=$(git rev-parse --short HEAD)` variable from the metadata-commit block rather than rename it.** DECIDE & STATE. The frontmatter edit sits *between* the assignment and the `git commit -m` that used `${HASH}`, so the variable could never survive to its use — the same block-boundary defect requirement 2 fixes in the range capture. Renaming it to a two-word name would have preserved a broken pattern; the block now reads the hash, has the agent write it into the file, and re-types it as a literal in the commit message.

**D-04 — Standardize on "merge range", dropping the "evidence range" wording my first draft introduced.** DECIDE & STATE. Two names for one concept is exactly the drift this REQ exists to remove, and the incumbent term is the one used by `actions/review-work.md`, which is outside this REQ's `write_set` and therefore cannot be brought along.

**D-05 — Document the cumulative range's over-inclusion instead of engineering it away.** DECIDE & STATE. Measured in a scratch repo: an orchestrator commit landing on the integration branch between a REQ's two merges falls inside `pre₁..merge_hash₂`. Suppressing it would need per-commit filtering the pipeline has no mechanism for, and the error direction is the safe one — an extra file surfaces as an undeclared touch a human judges, whereas the alternative (`pre₂..merge_hash₂`) silently hides the REQ's own work. Stated as a known cost with that reasoning, next to the definition.

## Review

**Acceptance: Pass — overall ~96%.** Main-context adversarial review: the spec-critical review workflow was blocked by a session usage limit (all 4 lens agents errored before running), so the orchestrator performed the review directly against the full 3-file diff plus an independent repo-wide restatement sweep.

**Requirements (all 6 met):**
1. qualify.sh FAILs loudly on an unresolvable range — verified by execution: `DO_WORK_DIFF_RANGE="bogus..range"` → exit 1 + `FAIL: DO_WORK_DIFF_RANGE does not resolve … not a commit`. Gate covers not-a-git-repo, non-two-dot form, empty bounds (with the "lost between command blocks" diagnostic), both endpoints via `git rev-parse --verify`, and a `git diff --name-only` catch-all.
2. Range held as re-typed literals — work.md Step 6 + work-reference.md prescribe holding both `<pre>` AND `<merge_hash>` as literals (the lock `session_id` pattern).
3. Seam-in-range — `git merge --no-ff --no-commit` → apply seams → `git commit` folds the seam INTO the merge commit; correct, and `commit:` still names the merge (provenance intact). No third token, no review-work.md edit (D-01).
4. Cumulative re-merge range — keep `<pre₁>`, re-capture only `<merge_hash>`; range `<pre₁>..<merge_hash₂>`; Step 9 records the latest; over-inclusion flagged as the safe error direction.
5. Step 6 imperative merge instruction — new orchestrator-side hand-back block (4 steps + `Already up to date` guard).
6. Commit-procedure carve-out — generic hash-writeback block is now SERIAL-ONLY; `HASH=$(git rev-parse --short HEAD)` replaced with a guarded serial-only line + re-typed literal; Step 9 warns against rev-parsing HEAD in worktree mode.

**Coherence:** "first parent" reconciled to "first-parent ancestor"; false "lands inside" sentence removed (0 repo hits); terminology unified to "merge range".
**Serial byte-compat:** new validation gated `if [ -n "$diff_range" ]`; all other qualify.sh changes are comments; verified byte-identical serial output.
**Independent restatement sweep:** grepped `first parent`, `merge range`/`evidence range`, `rev-parse --short HEAD`, `real implementation commit`, `lands inside` repo-wide — every site consistent; review-work.md's 3 range references true unchanged.

No Important/Critical findings. No follow-ups queued.

## Lessons Learned
**What worked:** `--no-ff --no-commit` → seam → commit folds the integration seam into the merge commit itself, resolving "seam inside `<pre>..<merge_hash>`" with no third endpoint token and no review-work.md edit (D-01). Re-defining capture points while keeping both token *names* is what let 5 downstream restatements stay true verbatim.
**What didn't:** The REQ's literal suggestion (capture `merge_hash` after a *separate* seam commit) was rejected — it would point `commit:` at a one-line seam commit instead of the implementation, trading a range defect for a provenance defect.
**Worth knowing:** qualify.sh runs `set -uo pipefail` without `-e`, so a `git diff` exiting 128 on an unresolvable range previously scrolled its fatal off-screen and passed vacuously; the new gate validates before any check consumes the range. Under remediation re-merge `<pre₁>` is a first-parent *ancestor* (not the direct parent) of the latest merge, but the merge-base identity still holds.

## Orientation

Worktree dispatch mode's merge range `<pre>..<merge_hash>` now has a complete, self-consistent production/validation/restatement contract across `actions/work.md` (Step 6/6.3/9), `actions/work-reference.md` (Worktree Dispatch Mode + Commit procedure), and `tools/checks/qualify.sh` — which now fails loudly on an unresolvable range instead of passing vacuously. Serial mode unchanged. No map change — hardens the REQ-037 subsystem.
