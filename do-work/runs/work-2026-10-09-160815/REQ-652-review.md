## Review: REQ-652

**Approve with follow-ups** — The operator rule is in all five files as specified, but clarify's new "Done" flip stamps the wrong field, so the board flags the REQ as a completion anomaly instead of showing it under Done.
Route A | merge range fa288233..92a6951c (builder commit 8f6f6e89, merge 92a6951c)

### What's built
- Capture never makes a REQ for a routine operator act (deploy, publish and the like). A tracked special operator configuration becomes one `blocked` REQ with `depends_on` on the AI-buildable REQs. It shows under Pending → Waiting first, then Needs input · Blocked.
- Clarify's plain-blocked prompt has `4. Done — I did it myself`. It writes `## Operator receipts`, sets `completed` in place and does not run `unblock`. Step 6 reports these REQs.
- Both run rows (work.md:523, work-reference.md:788) complete a REQ whose remaining work is a routine operator act. The blocked flip still applies to preconditions. The work guide states the board path.
- Missing: a REQ completed through the Done option does not reach the board's Done column (F1).

### Decisions / risks for you
- D-02 (an operator-only request writes no UR and no REQ): accepted. It declines non-queue intent the same way capture already skips an exact duplicate (capture.md:123 "If same: tell user, skip"). The absolute "always a UR" wording it bends was already bent by that duplicate skip, so this is noted only as a Nit (F4).
- D-04 (the `assigned_to` schema line and board.md's Needs input paragraph left unchanged): confirmed correct. See the Restatement sweep.

### Findings

**Important:**
- F1 `skills/do-work/actions/clarify.md:194`: the Done bullet stamps `status_changed_at: <now>` and never stamps `completed_at`. The schema says every terminal flip MUST stamp `completed_at` (work-reference.md:185-191). It also says `status_changed_at` is only for flips that have no dedicated stamp, and `completed_at` is a dedicated stamp (work-reference.md:172-176). On the board, a REQ with neither `completed_at` nor `commit` resolves to `CompletionUnresolved` (model.go `resolveCompletionTime`). `detectCompletionAnomaly` then flags "terminal status but no completed_at and no resolvable commit hash", so the card goes to Completion anomalies with a warning. `isWithinRecentWindow` returns false for the zero time, so the card never enters RecentlyDone (Done). This misses the maintainer's statement "when the operator section is done, the REQ should be moved to the DONE column" and the REQ's GREEN condition "which the board shows under Done". The root cause is the REQ's own Detailed Requirement 3, which named `status_changed_at`, and the builder followed it. The fix is one clause in the write set: stamp `completed_at: <now>` in place of `status_changed_at`. — impact-user-visible → report only

**Minor:**
- F2 `skills/do-work/actions/clarify.md:194`: the Done bullet calls itself "the same in-place flip the confirmed `builder_decided: true` path takes in Step 5", which is wrong. That path runs through the `answer` command (clarify.md:140-146), stamps `completed_at` and archives directly ("Builder Was Right" steps 2-3). It is not an in-place flip. The wrong precedent came from the REQ's Verified Facts citing clarify.md:267, and it hides F1. Recommended rewording: name the field that path stamps, not the in-place claim. — impact-negligible → report only

**Nit:**
- F3 A Done-completed REQ stays in `do-work/queue/` until cleanup Pass 0 runs. Until then `doctor` reports `STRANDED-TERMINAL-REQUEST` (doctor_scan.go:310, severity warning). This is expected from the design the REQ chose (no new transaction, and `complete` refuses a non-claimed REQ per state_plan.go:189), and the bullet already names Pass 0. Noted only so nobody reads the warning as a regression. — impact-negligible → report only
- F4 D-02: capture.md:10 ("Every invocation produces exactly two things, always paired") and SKILL.md:20 ("A capture always creates a UR") now have a second unlisted exception, an operator-only request, next to the existing duplicate skip at capture.md:123. The tension already existed and D-02 does not add new risk. — impact-negligible → report only

### Requirements Checklist

- [x] R1 capture External-condition assessment: routine act never captured ("examples, not a list", "own time", summary names the step), tracked configuration `blocked` + `blocked_by` + `blocked_at` + `depends_on` + checklist in body, Pending → Waiting then Needs input · Blocked — delivered
- [x] R2 Earmark assessment "at the keyboard" points at the new rule — delivered
- [~] R3 clarify Step 5.5 Done option, `## Operator receipts` under Outside-text containment, `completed` in place, no `unblock`, Step 6 lists them, precedent cited — delivered as specified, but the specified stamp prevents the Done-column outcome (F1), and the precedent is misdescribed (F2)
- [x] R4 work.md:523 and work-reference.md:788: complete instead of flipping for a routine act, blocked flip kept for preconditions (work-reference also names the tracked configuration) — delivered
- [x] R5 work-guide one sentence on the shape and the board path — delivered
- [x] R6 acceptance greps — delivered (the probe reran them at the gate)
- [x] R7 contract tests green — delivered (maintainer-verify exit 0 at the merge)

Specific checks from the coordinator:
1. Each maintainer statement maps to one sentence: routine act not captured (capture.md:108), tracked configuration Waiting → Needs input · Blocked (capture.md:108, work-guide.md:119), Done → completed (clarify.md:194). Met in the text. The Done outcome on the board fails per F1.
2. The Done bullet does not run `unblock` and cites Step 5 `builder_decided: true` and cleanup Pass 0. Both headings exist. The Step 5 characterization is wrong (F2).
3. The two run rows say complete-not-flip for a routine act and keep the blocked flip for preconditions. Met.
4. No field, status, flag, transaction or Go change: the diff is 5 `.md` files, +8/-6. Met.
5. D-02: accepted as built (see Decisions above and F4).
6. Cross-references: `actions/clarify.md` Step 5.5 (exists, :160), Step 4's Outside-text containment (named sub-contract, :101), Step 5 (:126), `actions/cleanup.md` → **Pass 0: Sweep Finished Queue Items** (:35), `actions/work-reference.md` → **Failure Classification (Step 8)** (:765). All land.
7. Retained `blocked_by`/`blocked_at` on a completed REQ: harmless. `bucketColumns` (model.go) routes by status, terminal first, without reading `blocked_by`. Pass 0 keys on status only. The append-only stamp rule (work-reference.md:73) removes `blocked_at` only on unblock, so keeping it is correct. The real board problem is the missing `completed_at` (F1).

### Acceptance Testing

**Result: Partial**
- Prose change, so acceptance was a code-path trace. The capture, run-row and guide text match the board code: model.go `bucketColumns` puts `blocked` with unmet `depends_on` under PendingWaiting and puts bare `blocked` under NeedsInputOrBlocked. board.md:92 says the same.
- The clarify Done path traced through model.go `resolveCompletionTime` → `detectCompletionAnomaly` → `isWithinRecentWindow` lands under Completion anomalies, not Done (F1).
- Gate evidence accepted as recorded: maintainer-verify.sh exit 0 at 92a6951c, and the probe (four acceptance greps plus shipped-package-reference-contract.sh) green.

### Suggested Additional Testing

- After F1 is fixed, build the board against a fixture REQ with `status: completed`, `completed_at`, retained `blocked_by`/`blocked_at`, and no `commit`. Confirm it shows under Done with no anomaly warning.

### Scores (on the record — not the headline)

**Overall: 82%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 90% | R3 stamp defeats its own board outcome |
| Code Quality | 85% | Wrong stamp against the schema; precedent misstated |
| Test Adequacy | N/A | Prose-only; probe greps plus contract tests are the evidence |
| Scope | 100% | Exactly the five write_set files |
| Risk | Low | Board anomaly on a new path; no data loss |
| Acceptance | Partial | Done option lands under Completion anomalies, not Done |

Computation: (90 + 85 + 100) / 3 = 91.7, minus 10 for Partial = 81.7 → 82%.

### Follow-ups created
None (4 findings report only)

## Re-review of fix delta

**Approve**: fix commit 0a1a8da0 resolves F1 and F2. The cumulative range is fa288233..e8b3687c.

- **Delta.** `git diff 92a6951c..e8b3687c` changes one line, the Done bullet in `skills/do-work/actions/clarify.md:194`.
- **F1 resolved.** The bullet now stamps `completed_at: <now>` and no longer writes `status_changed_at`. This matches the schema: `completed_at` on every terminal flip, and no `status_changed_at` when a flip has its own stamp (work-reference.md:172-191). On the board, `resolveCompletionTime` now gets its time from `completed_at`, and `detectCompletionAnomaly` flags nothing: a REQ captured `blocked` has no `claimed_at`, so the reversed-span check does not apply. `isWithinRecentWindow` is true, so the card shows under Done.
- **F2 resolved.** "The same in-place flip" became "Like the confirmed `builder_decided: true` path in Step 5, it completes a REQ that needs no build". That is accurate, because it no longer claims the two paths do the same file operations.
- **Cross-reference lands.** `actions/work-reference.md` → Request File Schema resolves to the heading `## Request File Schema — Full Frontmatter` (work-reference.md:57). This short form is the house spelling, used 13 times across `actions/*.md`.
- **No regression.**
  - Nothing else in the diff changed.
  - The bullet still does not run `unblock`, and the Pass 0 cross-reference is still there.
  - Step 6's wording ("completed by the operator (now `completed`)") still agrees with the bullet.
  - The integrator reports the repository gate exit 0 at e8b3687c (151 s) and the probe green. I accepted that as recorded.
- **Open.** F3 and F4 stay open as report-only Nits. The fix does not affect them.

Updated scores: Requirements 100%, Code Quality 95%, Test Adequacy N/A, Scope 100%, Risk Low, Acceptance Pass. Computation: (100 + 95 + 100) / 3 = 98.3 → 98%.

## Review

**Overall: 98%** | <timestamp>

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | N/A |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- F1 clarify.md:194 Done flip stamped `status_changed_at` instead of `completed_at`, so the board flagged a completion anomaly and never showed the REQ under Done. RESOLVED by 0a1a8da0 (now stamps `completed_at`, citing Request File Schema). — impact-user-visible → report only

**Minor findings:** F2 clarify.md:194 misdescribed the Step 5 `builder_decided: true` path as "the same in-place flip", RESOLVED by 0a1a8da0 (now "Like the confirmed `builder_decided: true` path in Step 5, it completes a REQ that needs no build") — impact-negligible → report only; F3 a Done-completed REQ shows doctor `STRANDED-TERMINAL-REQUEST` until cleanup Pass 0, which is expected, open — impact-negligible → report only; F4 D-02 adds a second unlisted exception to capture.md:10 / SKILL.md:20 "always a UR" next to the existing duplicate skip, accepted, open — impact-negligible → report only
**Acceptance:** Pass — capture, run rows and guide match board placement, and the clarify Done path stamps `completed_at` and lands under Done (model.go code-path trace at e8b3687c)
**Restatement sweep:** redefined the operator-blocked rule (routine act never captured, tracked configuration blocked with depends_on, clarify Done → completed, run completes instead of flipping); swept skills/: work-reference.md:114 `assigned_to` line agrees, do-work-board/actions/board.md:92 Needs input paragraph agrees, board-cards.js:240 tooltip agrees, work-guide.md:102 and capture-guide.md:63 not stale; the one disagreement (clarify Done stamping vs work-reference.md:172-191) was resolved by 0a1a8da0
**Suggested testing:** 1 items
**Follow-ups created:** None (4 findings report only)

*Reviewed by review-work action*
