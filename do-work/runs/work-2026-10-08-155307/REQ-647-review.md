## Review: REQ-647

**Approve**: clarify now tells the user when a REQ it just released is still earmarked, asks before clearing the field, and names it again in the report.
Route A | ffc1a7bb..be1f0ea4 (merge be1f0ea4, builder commit fe2c1878)

### What's built
- `skills/do-work/actions/clarify.md` Step 5.5 "Yes → unblock": after the unchanged transaction, a REQ that still carries `assigned_to` triggers a spoken warning and one two-option question (clear the earmark by hand edit, or keep it and run by name). The closing sentence now carries the `assigned_to` exception.
- Step 6 names each released-but-earmarked REQ with the verbatim session name and `do-work run REQ-NNN`. One matching Verification Checklist line was added.
- `skills/do-work/docs/work-guide.md` Earmarking paragraph gains one closing sentence. No other file changed.

### Decisions / risks for you
- None. Builder decisions D-01 to D-04 (commit routing through `do-work commit`, checklist line, keep option names the claim side effect, `<assigned_to, verbatim>` placeholder matching exit-summary row 6) are sound.

### Findings

**Important:**
- None.

**Minor:**
- None.

**Nit:**
- N1: `clarify.md:189` cites `actions/work-reference.md` → **Request File Schema** while the heading is "Request File Schema — Full Frontmatter". This short form is the house spelling: about 15 other sites use it (`clarify.md:30` and `:162`, `capture.md:107-108`, `work.md:249`, `capture-reference.md:20`, and others), only `work.md:89` uses the full form, and the shipped-package reference contract passes. No change needed. — impact-negligible → report only
- N2: `clarify.md:188-190` the two-option question names no recommended option. For the incident shape (`assigned_to: 'user-interactive'`, which the REQ-648 schema line now says is never a valid earmark), clearing is the evident answer, and naming it would make confirming the fast path, matching clarify's Rules section. — impact-negligible → report only
- N3: the builder stamped both Discovered Tasks `impact-low`, which is not a Step 10 token. Re-stamped here: DT1 (Step 5.5 `unblock` runs without `--commit`, so a clarify session that only unblocks leaves an uncommitted lifecycle change; pre-existing) is impact-user-visible; DT2 ("Builder Was Right / Discarded" describes hand edits while Step 5 says `answer` has no hand-edit fallback) is impact-rule-change. Neither is critical; both stay report-only. — impact-negligible → report only

### Requirements Checklist

- [x] R1 Step 5.5: transaction sentence and "do not reproduce the mutation free-form" byte-identical (diff shows them as unchanged prefix); `assigned_to` condition, verbatim session name, default-run-skips warning, one two-option question; closing sentence matches the REQ text verbatim — delivered
- [x] R2 commit discipline: "Clarify states no commit step for its hand edits" is true of clarify.md as a whole (the only `--commit` is Step 5's `answer` transaction at line 143; the unblock call, the Step 4 reclaim line and the Builder Was Right / Discarded fast paths name none). `actions/commit.md` exists and its scope ("manual edits ... work done between do-work runs", leftover files from `git status`) fits a hand edit — delivered
- [x] R3 Step 6: one condition-keyed sentence naming session name verbatim and `do-work run REQ-NNN` — delivered
- [x] R4 work-guide: existing sentences kept, one sentence appended — delivered
- [x] R5 checklist line beside the existing blocked-REQ line — delivered
- [x] R6 contract scripts green (builder worktree and gate at be1f0ea4, exit 0) — delivered
- [x] Constraints: prose only, no Go/CLI flag/transaction/probe; `work-reference.md` and `capture.md` untouched (`git diff --stat` lists only clarify.md and work-guide.md); never cleared silently ("Never clear the field without that answer, and leave it untouched on keep.") — delivered
- [x] Citations truthful: `work-reference.md:114` says "an assignment persists until an explicit run or a hand-edit clears it" and carries REQ-648's "**For another session or checkout only** ..." sentence; `:511` row 6 says "To drop an assignment without running it, remove the field by hand." Both support the hand-edit clear path — delivered

### Acceptance Testing

**Result: Pass**
- Cold read of the Red-Green case (blocked, `blocked_by: 'the operator'`, `assigned_to: 'user-interactive'`, answered "Yes, unblock it"): Step 5.5 runs the unchanged `unblock` transaction, then the new sentence routes to the warning naming `user-interactive` and the two-option question, and the closing sentence states the default run skips it. Step 6 then names it with the session name and `do-work run REQ-NNN`. On "keep" the field stays; on "clear" the line is removed by hand and left for `do-work commit`.
- `grep -n assigned_to skills/do-work/actions/clarify.md` hits lines 188, 189, 192, 198, 275 at HEAD; none at ffc1a7bb (per Qualification).

### Suggested Additional Testing

- None. Prose-only change, and the REQ forbids a probe.

### Scores (on the record — not the headline)

**Overall: 98%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | All seven requirements and every constraint met |
| Code Quality | 95% | Clear and condition-keyed; no recommended option (N2) |
| Test Adequacy | N/A | Prose-only; gate and reference contract green |
| Scope | 100% | Touched files equal `write_set` |
| Risk | None | — |
| Acceptance | Pass | Red-Green case reaches the question and the exception sentence |

### Follow-ups created
- None (3 findings report only)

## Review

**Overall: 98%** | <integrator stamps>

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | N/A |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- None

**Minor findings:** N1 `clarify.md:189` cites the schema by the short heading "Request File Schema", which is the house spelling at about 15 other sites; no change needed — impact-negligible → report only. N2 `clarify.md:188-190` two-option question names no recommended option although clearing is the evident answer for the incident shape — impact-negligible → report only. N3 builder Discovered Tasks carried non-canonical `impact-low`; re-stamped DT1 (unblock without `--commit`) impact-user-visible and DT2 (Builder Was Right hand edits vs Step 5 no-hand-edit wording) impact-rule-change, both noncritical — impact-negligible → report only
**Acceptance:** Pass — cold read of Step 5.5 routes the blocked + `assigned_to: 'user-interactive'` release to the two-option question and the closing exception; Step 6 names it.
**Restatement sweep:** redefined the Step 5.5 closing "re-enters the queue" sentence (now conditional on `assigned_to`); restatements read and consistent: `clarify.md:13`, `work-guide.md:102`, `work-guide.md:120`, `work-reference.md` exit-summary rows 3 and 6
**Suggested testing:** 0 items
**Follow-ups created:** None (3 findings report only)

*Reviewed by review-work action*
