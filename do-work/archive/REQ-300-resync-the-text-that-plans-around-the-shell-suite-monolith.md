---
id: REQ-300
title: "Resync the text that still plans around the pre-split shell behavior suite"
status: completed
created_at: 2026-08-20T08:37:00Z
status_changed_at: 2026-08-20T08:37:00Z
claimed_at: 2026-08-20T11:42:38Z
completed_at: 2026-08-20T11:51:20Z
route: B
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-08-20T11:42:38Z
  basis:
    - trivial short-circuit
user_request: UR-056
addendum_to: REQ-258
domain: general
review_generated: true
sweep: true
impact: impact-user-visible
effort_estimate: effort-mechanical
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: false
write_set:
- do-work/queue/REQ-263-tighten-qualifys-ownership-probe-and-warn-legibility.md
- do-work/queue/REQ-264-make-a-disarmed-p-a-u-audit-visible.md
- do-work/queue/REQ-271-make-the-layout-lock-step-see-every-spelling.md
---

# Resync the Text That Still Plans Around the Pre-Split Shell Behavior Suite

## What

REQ-258 split `_dev/tests/prescribed-shell-scripts-behavior.sh` into one case file per script under `_dev/tests/prescribed-shell-cases/`. Every *code* consumer was fine — there is one caller and it reads only the exit status. What went stale is the text that **plans around** the old single-file layout, and it is all the same root cause, so it is one sweep rather than four REQs.

The live consequence: the restart prompt tells the next session that this file is a scheduling bottleneck forcing four more waves. It is not one any more. REQ-263, REQ-264 and REQ-271 now write three different case files (`qualify.sh`, `qualify.sh`, `repair-req-timestamps.sh` respectively — so 263 and 264 still overlap each other, but neither overlaps 271), which changes how they can be batched.

## Instances

- [x] **`do-work/RESTART-PROMPT.md:33`** — **already closed, no edit.** Commit `6500808` rewrote the file before this REQ ran; the stale framing is gone and `:32`/`:37` now state the post-split reality, with `:49` already naming `prescribed-shell-cases/qualify.sh`. Verified by `grep -n "prescribed-shell" do-work/RESTART-PROMPT.md`. Original instance text: states `_dev/tests/prescribed-shell-scripts-behavior.sh` "is written by REQ-258, 263, 264, 268 and 271 — at most one per wave" and that REQ-258 will dissolve it. REQ-258 has shipped. Rewrite to the post-split reality: the runner is no longer written by case-adding REQs, REQ-263/264 share `prescribed-shell-cases/qualify.sh`, REQ-271 writes `prescribed-shell-cases/repair-req-timestamps.sh`, and 271 is therefore disjoint from both.
- [x] **`do-work/queue/REQ-263` `write_set`** — **edited.** Now names `_dev/tests/prescribed-shell-cases/qualify.sh` instead of the runner.
- [x] **`do-work/queue/REQ-264` `write_set`** — **edited.** Same change.
- [x] **`do-work/queue/REQ-271` `write_set`** — **edited.** Now names `_dev/tests/prescribed-shell-cases/repair-req-timestamps.sh`, verified as the file holding the `board_timestamp_layouts` guard at `:232-241`. The `## Red-Green Proof` command is unchanged, as directed. The stale observed count was corrected 66 -> 76 and re-framed so the count reads as context and `exit 0` as the finding; a pointer to the guard's post-split location was added beside it.
- [x] **`decisions/audits/2026-08-11-defensive-surface.md`, eleven Coverage rows** — **decided: leave them, no edit.** The file settles its own case at `:5` — "This is a historical snapshot, not a living registry" — reinforced by its `**Frozen:** after REQ-171 on 2026-08-12` stamp. Editing eleven rows inside a file that declares itself frozen would break the property that makes it usable as a record. REQ-234's precedent agrees: its restatement sweep found four stale restatements of a suite count and left all four as history. And the pointers are not wrong, only less precise -- the runner still exists, still executes all 76 cases, and still owns the exit status, so each cited case resolves by running it. Original instance text: each cites "`_dev/tests/prescribed-shell-scripts-behavior.sh` <name> case", and the named case now lives in `_dev/tests/prescribed-shell-cases/<script>.sh`. **Decide before editing:** this is a dated decision record, and REQ-234 set the precedent that dated history is left alone. The pointer is one hop from the case either way. If the call is to leave it, tick this box with that reasoning rather than editing.

## Context

From REQ-258's review, Important finding 1 and Minor finding 2, consolidated by root cause. `write_set` is display-only and Step 5.5 overwrites it from the fresh Scope declaration at claim time, so **no build is at risk** — the stale values mislead the board's overlaps badge and a human planning waves, which is exactly what `RESTART-PROMPT.md` exists for.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Verify each of the five instances against the tree before editing anything, because two of them may already be closed. Then edit only what is still stale: the three `write_set` fields. Instance 1 is verify-and-close; instance 5 is decide-and-close on the audit file's own declared contract.
- [x] **[APPLY]:** Three `write_set` edits plus one stale observed-count figure. No file outside `do-work/queue/` touched.
- [x] **[UNIFY]:** `git diff --stat` reviewed; `bash _dev/tests/maintainer-verify.sh` exits 0; no debug artifacts. Per-file verification listed in `## Implementation Summary`.

---

## Triage

**Route: B** - Medium

**Reasoning:** The five instances name their own targets, but three of them assert a stale state that had to be re-verified against the tree first — two turned out already closed, which changes the work from "edit five things" to "edit three, close two with evidence." That verification is the exploration.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

Each instance was checked against the tree before any edit.

**Instance 1 — `do-work/RESTART-PROMPT.md:33`: already closed, no edit needed.** The stale text this REQ quotes is gone. Commit `6500808` rewrote the file, and the current text states the post-split reality directly:

- `:32` — "`_dev/tests/prescribed-shell-scripts-behavior.sh` keeps its path and exit-status contract but is now a 35-line runner. The 76 cases live one file per script in `_dev/tests/prescribed-shell-cases/` … **A REQ that adds a case now writes that script's case file, not the runner.**"
- `:37` — names the old claim as wrong and hands the rest of the class to this REQ: "This file replaces that text; REQ-300 sweeps the rest."
- `:49` — the ordering table already names `_dev/tests/prescribed-shell-cases/qualify.sh` as REQ-263/264's shared file.

`grep -n "prescribed-shell" do-work/RESTART-PROMPT.md` returns no occurrence of the "at most one per wave" framing or of the five-REQ writer list. This instance closes by verification.

**Instances 2 and 3 — REQ-263 and REQ-264 `write_set`: stale, confirmed.** Both still list `_dev/tests/prescribed-shell-scripts-behavior.sh`. Their cases belong in `_dev/tests/prescribed-shell-cases/qualify.sh` (which exists, and already holds the REQ-254 cases — `qualify: 3 cases, 0 failures`).

**Instance 4 — REQ-271 `write_set`: stale, confirmed, and the target file is verified.** The guard REQ-271 fixes is at `_dev/tests/prescribed-shell-cases/repair-req-timestamps.sh:232-241` — the `board_timestamp_layouts` extraction and its `fail_case`. So `repair-req-timestamps.sh` is the right file, as the instance says. Two sub-findings:

- The `## Red-Green Proof` **command** is correct and stays, as the instance directs.
- The **observed count** in that paragraph reads "66 named cases"; the suite now reports **76**. Updated, because the paragraph is a builder's re-runnable RED observation rather than dated history, and a figure that disagrees with what the builder will see undermines the evidence it is there to supply. The load-bearing half of that observation — `exit 0`, meaning the guard does not notice the mutation — is unchanged.
- Noted for whoever builds REQ-271, not fixed here: its defect 2 (GNU-only `\|` alternation) appears **already resolved** — the line is now `sed -n -E 's/^[[:space:]]*(time\.RFC3339|"2006[^"]*"),$/\1/p'`, which is ERE and portable. Defect 1 (the extraction enumerates spellings instead of keying on the condition) is still live. Out of scope for this REQ, which resyncs text and does not touch that file.

**Instance 5 — the eleven Coverage rows: decided, leave them.** The audit file settles its own case in its header, `decisions/audits/2026-08-11-defensive-surface.md:5`: "**This is a historical snapshot, not a living registry.** It records the defensive surfaces reviewed through REQ-171; new shell files and prose sections do not update it." Its `**Frozen:** after REQ-171 on 2026-08-12` stamp says the same thing. Editing eleven rows inside a file that declares itself frozen would break the contract that makes it trustworthy as a record.

REQ-234's precedent, which the instance asks to be weighed, points the same way — its restatement sweep found four stale restatements of a suite count and left all four: "all four are history, correctly left alone, and the builder explicitly declined to rewrite them."

The pointers are also not *wrong*. `prescribed-shell-scripts-behavior.sh` still exists, still runs all 76 cases, and still owns the exit status — so "the `prescribed-shell-scripts-behavior.sh` real-merge case" resolves by running that file. It is less precise than it could be, one hop from the case, exactly as the instance says.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `do-work/queue/REQ-263-tighten-qualifys-ownership-probe-and-warn-legibility.md` (modify) — `write_set` names the case file
- `do-work/queue/REQ-264-make-a-disarmed-p-a-u-audit-visible.md` (modify) — same
- `do-work/queue/REQ-271-make-the-layout-lock-step-see-every-spelling.md` (modify) — `write_set` names the case file; stale observed count corrected

**Files I will NOT touch:**
- `do-work/RESTART-PROMPT.md` — instance 1 is already closed by `6500808`; re-editing correct text is churn
- `decisions/audits/2026-08-11-defensive-surface.md` — decided to leave, per its own frozen-snapshot header and REQ-234's precedent
- `_dev/tests/prescribed-shell-cases/repair-req-timestamps.sh` — REQ-271's fix, not this REQ's; this REQ resyncs text that *plans around* the suite and changes no test
- `skills/do-work/tools/checks/qualify.sh` — REQ-263/264's target

**Acceptance criteria (restated from REQ):**
- [ ] Every instance is closed, by edit or by recorded decision — none left ambiguous
- [ ] No `write_set` in the queue still names `_dev/tests/prescribed-shell-scripts-behavior.sh` as a file a case-adding REQ writes
- [ ] The wave-planning text a next session reads matches the post-split layout
- [ ] `bash _dev/tests/maintainer-verify.sh` exits 0

## Pre-Flight

**Git:** ✓ working tree clean outside `do-work/`
**Tests baseline:** ✓ passing (`bash _dev/tests/maintainer-verify.sh`, exit 0, `launched: true`)
**Dependencies:** ✓ toolchain present — go1.26.1 (pinned via `go env -w GOTOOLCHAIN`), ShellCheck 0.11.0, `just` 1.21.0

*Checked by work action*

## Implementation Summary

**Files changed:**
- `do-work/queue/REQ-263-tighten-qualifys-ownership-probe-and-warn-legibility.md` (modified)
- `do-work/queue/REQ-264-make-a-disarmed-p-a-u-audit-visible.md` (modified)
- `do-work/queue/REQ-271-make-the-layout-lock-step-see-every-spelling.md` (modified)

**What was done:** Repointed three queue REQs' `write_set` from the dissolved monolith `_dev/tests/prescribed-shell-scripts-behavior.sh` to the per-script case files their cases actually belong in — `prescribed-shell-cases/qualify.sh` for REQ-263 and REQ-264, `prescribed-shell-cases/repair-req-timestamps.sh` for REQ-271, the latter verified against the `board_timestamp_layouts` guard at that file's `:232-241`. Corrected REQ-271's stale observed case count (66 → 76) and re-framed the paragraph so `exit 0` reads as the finding and the count as context, plus a pointer to the guard's post-split location. Closed the other two instances without editing: `RESTART-PROMPT.md` was already fixed by `6500808`, and the defensive-surface audit declares itself a frozen historical snapshot, so its eleven Coverage rows are left as history per REQ-234's precedent.

**Every file changed by this REQ is a `do-work/` queue path, by design** — see `## Qualification` for why that trips two mechanical checks and what the evidence is.

## Qualification

**Two mechanical checks FAIL, both from one cause, both overridden with evidence. The cause is this REQ's shape, not its implementation.**

`tools/checks/qualify.sh` → `FAIL: Implementation Summary lists only do-work/ paths — the REQ was not implemented` (exit 1).
`tools/checks/scope-drift.sh` → `FAIL: a '**Files I will touch**' header is present in ## Scope but no backticked paths parse from it` (exit 1).

Both exclude `do-work/` at the pathspec level — `qualify.sh:160-162` and `scope-drift.sh`'s `grep -v '^do-work/'` on the declared and reported lists. That exclusion is right for the case it was written for: a builder that touched only queue metadata did no work. **This REQ's entire deliverable *is* queue metadata** — it resyncs three queued REQs' `write_set`, which is the whole point of the sweep — so both checks empty their input and report the empty result as failure. Neither is evidence about the implementation.

**Evidence the work landed** (file system and git, not the summary's claim — Anti-Rationalization row 1):

```
$ git diff --numstat -- <the three declared files>
1  1   do-work/queue/REQ-263-tighten-qualifys-ownership-probe-and-warn-legibility.md
1  1   do-work/queue/REQ-264-make-a-disarmed-p-a-u-audit-visible.md
4  2   do-work/queue/REQ-271-make-the-layout-lock-step-see-every-spelling.md
```

**Evidence the acceptance criterion is met** — no queue `write_set` still names the dissolved runner:

```
$ grep -rn "prescribed-shell-scripts-behavior" do-work/queue/
REQ-301:27  ... byte-identical in `git show HEAD:_dev/tests/prescribed-shell-scripts-behavior.sh`
REQ-271:53  **RED:** ... run `bash _dev/tests/prescribed-shell-scripts-behavior.sh`
```

Two hits remain and **both are correct, neither is a `write_set`**: REQ-301's is a historical citation to a past tree state, and REQ-271's is the Red-Green Proof entry command that this REQ's own instance 4 directs must not change. The runner still exists and still owns the exit status, so both references resolve.

**Scope-drift comparison, done by hand** because the script cannot parse a `do-work/`-only declaration (work.md Step 5.5 provides for the by-hand fallback). Declared: the three queue REQ files. Reported: the same three. Set-difference both directions: empty. **No drift** — no undeclared touch, no unused declaration.

**Judgment checks:** *(2) Substantive* — each diff is a real content change, not whitespace: two one-line `write_set` repoints and REQ-271's four-line change (`write_set`, the corrected count with its reframe, and the guard-location pointer). *(3) Requirements traced* — all five `## Instances` carry a verdict, three by edit and two by recorded decision, and each of the four acceptance criteria has evidence above or in `## Testing`. *(6) Flowing* — not applicable; no data path, no new file.

**Not filed as a follow-up, deliberately.** The two-check blind spot for a queue-metadata REQ is real and now has a worked example, but this is the first such REQ, and the 2026-08-20 remaining-work report's stopping policy is explicit that generation-two findings default to no new task. Recorded here and in `## Lessons Learned` so a second occurrence has a precedent to cite rather than re-deriving it. Surfaced in the hand-back for the user's call.

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh`
**Result:** ✓ exit 0 (unchanged from the pre-flight baseline)

**Red-green validation:** omitted — `tdd: false`, and this is a non-behavioral change. No shipped code, script, or test was touched; the diff is three queued REQ files' prose and frontmatter. Regression evidence is the unchanged green baseline plus the two greps above.

**New tests added:** none. Nothing here is mechanically checkable without a rule that queue-metadata text must match the test-file layout, which would be a new contract rather than this sweep's job.

*Verified by work action*

## Review

**Overall: 93%** | 2026-08-20T11:50:31Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 90% |
| Test Adequacy | N/A |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- None

**Minor findings:** 2 (report only)
- Two of five instances closed without an edit. That is the right outcome on the evidence, but it means the sweep's realized value is smaller than its capture implied: the durable change is three `write_set` lines and one corrected figure. Recorded so the next reader does not expect five edits from the commit.
- REQ-258's archived trail (`do-work/archive/REQ-258-…:231`) recorded "No action" on REQ-271's stale count, and this REQ changed it. Not a contradiction on inspection — REQ-258's reason was scope ("stale before this REQ"), not that the figure should stay wrong, and REQ-300's instance 4 puts it explicitly in scope. Noted because a reader comparing the two trails would otherwise see a reversal with no explanation.

**Restatement sweep:** the diff changes one stated figure and three pointer values, so the sweep asked who else states them. `66 named cases` appears in three archived files (REQ-257 twice as its own delivery record, REQ-258 once as the deferral above) — all history, correctly left alone. The new pointer values are restated in exactly two live places, `do-work/RESTART-PROMPT.md:49` and REQ-301's `write_set`/body, and **both already name the case files correctly** — REQ-301 was captured after the split. `skills/do-work/CHANGELOG.md:10` cites `prescribed-shell-cases/repair-req-timestamps.sh` correctly. No live restatement is stale.

**Acceptance:** Pass — every instance carries a verdict; `grep -rn "prescribed-shell-scripts-behavior" do-work/queue/` leaves only two references that are correct by design (a historical citation and the Red-Green entry command instance 4 protects); `maintainer-verify.sh` exits 0.

**Suggested testing:** 1 item
- Nothing mechanically prevents a future `write_set` from naming a test path that does not exist. A queue-lint asserting every `write_set` path resolves in the tree would have caught all three of these staleness instances at capture time rather than three REQs later. Worth a line if the class recurs; not filed, per the stopping policy.

**Code Quality 90%, not 100%:** the two no-edit verdicts are recorded inside the `## Instances` bullets by prefixing the original text with the verdict, which leaves each bullet carrying both the new decision and the superseded instruction. Legible, and it preserves what was asked, but a reader meets the verdict and the stale ask in one paragraph.

**Follow-ups created:** None; **sweeps appended to:** None

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Verifying every instance against the tree before editing any of them. Two of five were already closed — one by a commit that landed after this REQ was captured (`6500808`), one by the target file's own header — and a sweep that had trusted its instance list would have rewritten correct text in `RESTART-PROMPT.md` and broken a file that declares itself frozen. On a `sweep: true` REQ the instance list is a capture-time hypothesis, not a work order.

**What didn't:** Both mechanical gates failed, and neither failure was about the code. `qualify.sh` and `scope-drift.sh` each exclude `do-work/` at the pathspec level, so a REQ whose deliverable *is* queue metadata empties their input and reads as unimplemented. The by-hand fallback that work.md Step 5.5 provides for a missing script turned out to be the right tool for a script that runs but cannot see the work — which is not the case that fallback was written for.

**Worth knowing:** `decisions/audits/2026-08-11-defensive-surface.md` settles questions about itself in its own header — "This is a historical snapshot, not a living registry" — so a sweep that reaches it should read line 5 before planning eleven edits. The same shape is worth checking on any dated file under `decisions/`. And REQ-234's archived restatement sweep is the citable precedent for leaving dated history alone; it is more useful as a reference than re-deriving the argument each time.

## Orientation

Wave planning for the queued shell-suite REQs now reads the post-split layout: REQ-263 and REQ-264 declare the `qualify.sh` case file they share, REQ-271 declares the `repair-req-timestamps.sh` case file that actually holds the guard it fixes, so the board's overlaps badge and a human batching work both see the real collision set — 263 with 264, and 271 disjoint from both. Lives in the queue metadata that `_dev/primes/prime-shell-commands.md` governs the test side of.

Not `[MAP CHANGED]` — three pointer values and one figure; no contract, structure, or concept altered. Staleness spot-check on `_dev/primes/prime-shell-commands.md`: every referenced path resolves, and nothing in it names the pre-split single-file layout. The prime is not stale.
