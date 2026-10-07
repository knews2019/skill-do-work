## Review: REQ-644

**Approve** — capture now sends "I need to be at the keyboard" work to `blocked` and keeps `assigned_to` for another session. Three small wording findings, all report only.
Route A | 58168fe7 (range cd91f838..58168fe7)

### What's built
- `skills/do-work/actions/capture.md:107` (Earmark assessment) has one new sentence: `assigned_to` routes work between sessions; operator-present work is captured `blocked` with `blocked_by` and `blocked_at`, and `do-work clarify` releases it.
- `capture.md:108` (External-condition assessment) lists the earmark as a look-alike. The closed count "the three look-alikes" became "every look-alike that is not blocked".
- `skills/do-work/docs/work-guide.md:113` opener now says "leave this one for that session", and line 119 ends with the same split in one sentence.

### Decisions / risks for you
- None. Builder decisions D-01 to D-03 (from the REQ file) are sound. D-01 removes a count that would go stale.

### Findings

**Important:**
- None.

**Minor:**
- F1. Case mismatch with the board's column name. The board heading is `Needs input · Blocked` (`skills/do-work-board/tools/queue-kanban/web/template.html:281`), and `work-guide.md:102`, `skills/do-work-board/docs/board-guide.md:18`, and `abandon.md:153` use the lowercase "input". The new text at `capture.md:107` and `work-guide.md:119` says `Needs Input · Blocked`. So the same file (`work-guide.md`) now spells it two ways, and a case-sensitive search for the board's label misses the new lines. The column itself is the right destination. Fix: lowercase "input" at both sites. — impact-negligible → report only
- F2. "naming that person" can push a capture agent past the never-invent rule. In the RED case the user says "I" and "me", not a name. `capture.md:107` says `blocked_by` names "that person in the user's words", while the External-condition assessment it points to says `blocked_by` holds "the condition, in the words the user used". An agent could resolve "me" to the git user name and write that, which is an invented value. Safer wording: "`blocked_by` stating the operator condition in the user's words". The REQ's own constraint ("no person named means no `blocked_by`") has the same tension with its RED case. — impact-rule-change → report only

**Nit:**
- F3. "lands it in Needs Input · Blocked" is true only when the REQ has no unmet `depends_on`. The board places a `blocked` REQ with an unmet dependency under Pending → Waiting (`skills/do-work-board/tools/queue-kanban/model.go:383,1719`). A new operator-hold capture rarely has dependencies. — impact-negligible → report only

### Requirements Checklist

- [x] 1. Earmark assessment sentence distinguishing `assigned_to` from operator-present `blocked` + `blocked_by` + `blocked_at`, released by clarify, pointing to the External-condition assessment rather than restating mechanics — delivered (`capture.md:107`)
- [x] 2. Earmark added as a look-alike in the External-condition assessment — delivered (`capture.md:108`)
- [x] 3. work-guide Earmarking paragraph sentence — delivered (`work-guide.md:119`), plus opener reworded (D-02)
- [x] Constraint: prose only, no field/status/Go/routing change — met (diff touches 2 .md files, 4+/4-)
- [x] Constraint: never-invent rule kept — "Never invent or infer a session name" untouched; `blocked_by` still from the user's words (see F2 for a wording risk)
- [x] No bold label or heading changed — verified in the word diff; only sentence bodies changed
- [x] Text does not say where the board places an earmarked pending REQ — verified; both new sentences describe only where a `blocked` REQ shows, so they stay true after REQ-643 (the sibling REQ that adds a Pending → Earmarked sub-group)

### Acceptance Testing

**Result: Pass**
- Cold read of the Red-Green Proof against `capture.md:107-108` at 58168fe7. "I need to be at the keyboard for this one, hold it for me" names no session. The Earmark assessment's new sentence quotes this exact case and routes it to `blocked` via the External-condition assessment with `blocked_by` and `blocked_at`. The External-condition list no longer leaves the earmark as the only match.
- "leave this for cloud-alpha" names a session. The Earmark assessment's first sentence still seeds `assigned_to: 'cloud-alpha'`, and the look-alike entry at line 108 says reserved-for-another-session work is not blocked.
- Release path checked: `skills/do-work/actions/clarify.md:13,30` collects `blocked` REQs for Step 5.5 and confirms human-confirmable conditions.
- Did not re-run tests. The REQ records `maintainer-verify.sh` exit 0 at 58168fe7 and the focused probe green.

### Suggested Additional Testing

- None beyond the recorded citation-contract run. This is a prose change with unchanged bold labels.

### Restatement Sweep

Trigger: the diff redefines the earmark boundary (when `assigned_to` applies versus `blocked`) and the External-condition look-alike list. I grepped `skills/` for `look-alike`, `three look-alikes`, `leave this one for me`, and `assigned_to` in .md files. The only look-alike list is `capture.md:108`. `work-reference.md:114` (schema line, "the session this REQ is earmarked for"), `work-reference.md:511` (skip report), and `capture-reference.md:41` (unassigned-only conversion) all still agree. `capture.md:143` (addendum earmark rule) only covers an addendum that names a session. It does not contradict the new rule. The builder already reported the gap as a discovered task, and I agree it is report only. No stale restatement found.

### Scores (on the record — not the headline)

**Overall: 95%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | All three sites delivered within one sentence each |
| Code Quality | 88% | Clear prose. F1 casing, F2 "naming that person" wording |
| Test Adequacy | 92% | Prose change: contract tests plus cold read are the right evidence |
| Scope | 100% | Exactly the declared write_set |
| Risk | Low | Advisory capture guidance, no code path |
| Acceptance | Pass | RED case routes to blocked. Session earmark is unchanged |

### Follow-ups created
- None (3 findings report only)

## Review

**Overall: 95%** | <TIMESTAMP>

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 88% |
| Test Adequacy | 92% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
None

**Minor findings:** F1 `capture.md:107` and `work-guide.md:119` write "Needs Input · Blocked", but the board heading and other docs use "Needs input · Blocked" — impact-negligible → report only. F2 `capture.md:107` "`blocked_by` naming that person" may lead an agent to substitute a name for "me". Prefer "the operator condition in the user's words" — impact-rule-change → report only. F3 (nit) "lands it in Needs Input · Blocked" is true only without an unmet `depends_on` (otherwise Pending → Waiting) — impact-negligible → report only
**Acceptance:** Pass — the cold read routes the keyboard case to `blocked` with `blocked_by`/`blocked_at`, and "leave this for cloud-alpha" still seeds `assigned_to`
**Restatement sweep:** redefined the earmark-vs-blocked boundary and the External-condition look-alike list. All restatements (work-reference.md:114,511, capture-reference.md:41, capture.md:143) still agree
**Suggested testing:** 0 items
**Follow-ups created:** None (3 findings report only)

*Reviewed by review-work action*

## Delta re-review (58168fe7..3602e584)

**Verdict:** Pass. **Overall score:** 97% (Code Quality 88% → 94%; other dimensions unchanged).

- The delta touches only `skills/do-work/actions/capture.md:107` and `skills/do-work/docs/work-guide.md:119`. Each changes one phrase, "Needs Input · Blocked" → "Needs input · Blocked". Nothing else changed.
- The new spelling matches the board column heading in `skills/do-work-board/tools/queue-kanban/web/template.html`, `work-guide.md:102`, and `board-guide.md:18`. No "Needs Input · Blocked" remains under `skills/`.
- The full range `cd91f838..3602e584` still touches only those two files (4 lines). No regression against the first review.

**F1:** Resolved. **F2, F3:** unchanged, report only. **New findings:** None.
