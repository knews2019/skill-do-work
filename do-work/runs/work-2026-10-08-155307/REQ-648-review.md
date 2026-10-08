## Review: REQ-648

**Approve** — both rule sentences landed as quoted, the dangling pointer is fixed, and nothing else in either file changed.
Route A | merge range 336cbb52..4c50ccde

### What's built
- The `assigned_to` schema line (`skills/do-work/actions/work-reference.md:114`) now says the field is for another session or checkout only, and that work waiting on the operator is `status: blocked` with `blocked_by`.
- The work action's Error Handling table (`skills/do-work/actions/work.md:523`) has a new row: an orchestrator that wants to park a queued REQ on the operator flips it to `blocked` and never writes `assigned_to`.
- The existing missing-precondition row (`work.md:522`) now points at `actions/work-reference.md` → **Failure Classification (Step 8)**, Environment row, a heading that exists (`work-reference.md:765`, row at `:788`).

### Decisions / risks for you
- None. The builder's three decisions (D-01 sentence placement, D-02 pointer spelling, D-03 pointer in parentheses inside the quoted row) are recorded in the hand-back and match the diff.

### Findings

**Important:**
- None.

**Minor:**
- F1: `work-reference.md:788` says "see `actions/work.md`'s mid-run blocked-flip procedure", but `work.md` has no such procedure. After this REQ, `work.md:522-523` point back at row `:788`, so the two citations form a loop. The full rule text is already in row `:788`, so a reader loses one wasted hop and nothing else. This is the builder's discovered task. It is a real finding, but the REQ's Detailed Requirement 3 kept this prose on purpose. A later fix would delete the "(see `actions/work.md`'s ...)" clause. `work-reference.md:165` ("Step 8's blocked-flip procedure") is accurate enough, because `work.md` Step 8 substep 3 (`work.md:468`) names the blocked-flip judgments. — impact-negligible → report only

**Nit:**
- F2: The new row says the default run "hides it from the user by skipping it". The schema line and `docs/work-guide.md:120` say the default scan skips **and reports** an assigned REQ. So "hides" overstates it: the REQ is visible only as one skip line, not as "waiting on you". The wording was quoted in the REQ and accepted by the maintainer as worded, so it was not changed. — impact-negligible → report only

### Requirements Checklist

- [x] DR1: schema sentence appended to `work-reference.md:114`, verbatim, every prior sentence kept (word-diff shows one insertion only) — delivered
- [x] DR2: new Error Handling row beside the missing-precondition row, verbatim apart from the D-03 pointer parenthetical — delivered
- [x] DR3: both rows point at **Failure Classification (Step 8)**, Environment row. `grep -rni "Mid-run blocked flip" skills/` finds only the lowercase prose at `work-reference.md:165` and a historical lesson at `do-work-board/tools/queue-kanban/lessons-do-kanban.md:36`, no heading pointer — delivered
- [x] DR4: `capture.md`, `capture-reference.md`, `docs/work-guide.md`, `model.go` absent from the diff — delivered
- [x] DR5: `grep -n operator` hits `work-reference.md:114` and `work.md:523` — delivered
- [x] DR6: contract tests green (builder ran both; this review re-ran the reference contract) — delivered
- [x] DR7: release deferred to the integrator — N/A for this review
- [x] Constraint: no other table row or bold label changed. Word-diff shows exactly 3 changes: one sentence insert, one pointer swap inside row 522, one new row — delivered
- [x] Constraint: nothing adds a case for writing `assigned_to`. The new sentence narrows the existing writers ("Seeded by capture ... or written by a session claiming from another checkout"); the new row forbids the write — delivered

### Restatement sweep (this diff's own element only, not wave-end)
Element redefined: who `assigned_to` is for (session or checkout, never the operator).
- `skills/do-work/actions/capture.md:107` (Earmark assessment): "`assigned_to` routes work between sessions; when the user instead has to be present as operator ... capture the REQ `blocked`" — agrees.
- `skills/do-work/actions/capture.md:108` (External-condition look-alikes): "work reserved for another session or checkout is an earmark → `assigned_to`" — agrees.
- `skills/do-work/docs/work-guide.md:120`: "The field is for another session or checkout; work that waits on you as operator is captured `blocked` with `blocked_by` naming you" — agrees.
- `skills/do-work/actions/capture-reference.md:41` and `capture.md:143`: earmark mechanics only (fold eligibility, addendum carry), no statement of who the field is for — no conflict.
- `work.md:37` (fan-out ready set: "not `assigned_to` another session") — agrees.
No stale restatement found.

### Acceptance Testing

**Result: Pass**
- `grep -n operator` on both files: two hits, as required.
- `## Failure Classification (Step 8)` exists at `work-reference.md:765`, Environment row at `:788` opens with the blocked-flip test.
- `bash _dev/tests/shipped-package-reference-contract.sh`: PASS (re-run by this review).
- `git diff --check 336cbb52..4c50ccde`: exit 0.
- Builder and orchestrator evidence: repository gate exit 0 at `4c50ccde`, probe satisfied, `contract-regressions.sh` exit 0 (not re-run here).

### Suggested Additional Testing
- None. Prose-only change; the acceptance greps cover it.

### Scores (on the record — not the headline)

**Overall: 98%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | All seven requirements and both constraints met |
| Code Quality | 95% | Clean prose, correct pointer spelling; F2 wording nit |
| Test Adequacy | 95% | Read-check red-green, contract tests and gate green |
| Scope | 100% | Touched files equal `write_set`; decisions recorded |
| Risk | None | Prose only, no parser or routing change |
| Acceptance | Pass | Greps and contract test confirm |

### Follow-ups created
- None (2 findings report only)

## Review

**Overall: 98%** | <timestamp>

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 95% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:** F1: `skills/do-work/actions/work-reference.md:788` cites "`actions/work.md`'s mid-run blocked-flip procedure", which does not exist; `work.md:522-523` now point back at that row, so the citations loop (rule text is complete at `:788`; kept by DR3 on purpose; fix is deleting the "see work.md" clause) — impact-negligible → report only. Nit F2: `work.md:523` says the default run "hides" an earmarked REQ, while the schema line and `docs/work-guide.md:120` say it skips and reports; wording accepted as quoted by the maintainer — impact-negligible → report only
**Acceptance:** Pass — `grep -n operator` hits both files, target heading exists, reference contract PASS, `git diff --check` clean
**Restatement sweep:** redefined `assigned_to` audience (session or checkout only, operator work is `blocked`): `actions/capture.md:107`, `capture.md:108`, `docs/work-guide.md:120`, `work.md:37` agree; `actions/capture-reference.md:41`, `capture.md:143` state mechanics only, no conflict
**Suggested testing:** 0 items
**Follow-ups created:** None (2 findings report only)

*Reviewed by review-work action*
