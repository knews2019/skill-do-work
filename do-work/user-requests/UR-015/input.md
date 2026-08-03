---
id: UR-015
title: Adversarial audit of the 0.164.0-0.167.0 batch
created_at: 2026-08-03T16:53:42Z
requests: [REQ-077, REQ-078, REQ-079, REQ-080]
word_count: 6
---

# Adversarial audit of the 0.164.0-0.167.0 batch

## Summary

The user asked for an adversarial audit of the six REQs that shipped 2026-08-03 (REQ-071 … REQ-076,
versions 0.164.0 → 0.167.0). Unlike a normal UR, the verbatim input is a six-word instruction; the
substance of this UR is the audit's own findings, verified against the repo rather than against the
session's completion summary.

**Provenance.** The findings below are the auditing agent's, not the user's. They were produced by
re-deriving each shipped claim from the working tree: running `_dev/tests/contract-regressions.sh`,
`gofmt`/`go vet`/`go test` in `tools/queue-kanban/`, `queue-kanban verify`, and `queue-kanban now`
against `date -u`; and by sweeping both fingerprints of every premise the batch retired. The user
reviewed and approved the resulting remediation plan before capture, which is the authorization for
these four REQs. Anything the audit could not execute is marked unverified rather than asserted.

**What the audit confirmed as sound** — recorded so a later reader does not re-litigate it: every
mechanical check passes; version 0.167.0 agrees with the top changelog entry with strictly ascending
versions and no reused titles; REQ-075's eleven-site sweep is genuinely complete for both the strong
and the weak fingerprint; and `tools/queue-kanban/timestamp.go` truncates rather than rounds and
converts rather than relabels, which were its two available ways to be wrong.

## Extracted Requests

| REQ | Title | Finding |
| --- | --- | --- |
| REQ-077 | Crash recovery's own-crash branch is unreachable, and its retired premise survives in the same file | F1 + F3 (High) |
| REQ-078 | The Windows timestamp fallback cannot run on stock Windows in either shell it names | F5 (Medium) |
| REQ-079 | Two guards pin the weaker fingerprint of the premise they exist to retire | F4 + F7 (Medium/Low) |
| REQ-080 | The capture template emits a stray instruction line into every REQ it produces | New finding, predates this batch |

## Batch Constraints

- **REQ-077 is the only one that changes runtime behavior.** The other three are instruction and
  guard corrections. All four are `maintenance: true`: each one's candidate fix is a removal or a
  narrowing of the skill's own operating instructions, which is the marker's stated trigger.
- **REQ-077 and REQ-079 both touch `_dev/tests/contract-regressions.sh`.** Different assertion
  blocks (REQ-071's at ~line 364, REQ-075's at ~line 137), so they do not conflict textually, but
  the merge is the non-interference proof if they are ever built concurrently — not the `write_set`
  overlap badge (`actions/work-reference.md` → Worktree Dispatch Mode → Fan-Out Dispatch).
- **Two audit items are deliberately not REQs.** The stray `do-work/working/baseline.json` and the
  three root `HANDOFF*.md` files (F8) belong to the next `do-work cleanup` run. The KB handoff (F6) is
  already done — see the correction below.
- **One audit finding was withdrawn.** F6 was filed as "the KB handoff ran and was left uncommitted,
  with `kb_entry` pointing at untracked files." That was a stale read: the audit's opening `git status`
  ran ~1 second before commit `336692f` landed, which committed both halves correctly. A **second
  session was writing to this checkout during the audit** — it has since stopped (no `worktree-agent-*`
  branches, no `claimed` REQ in `do-work/working/`, checkpoint reports nothing interrupted), so there
  was no live race, but it is the one-owner-per-checkout hazard showing up in practice, which is what
  REQ-077 and REQ-079 are both circling. Surviving sub-issues: REQ-071/072/073 are still
  `kb_status: pending`, and REQ-075's `kb_entry` filename enshrines the pre-build count of five that
  the build itself corrected to eleven. Neither is a REQ; both are one-line follow-ups for the user.
- **One gap stays open and is nobody's REQ yet:** REQ-073's live two-builder acceptance test has
  still never run. Everything built since was serial, so grep proves the prose and nothing proves two
  builders compose. The procedure is in `do-work/archive/UR-013/REQ-073-fan-out-dispatch-n-builders-one-owner.md`
  under `## Review` → Suggested additional testing. It needs a human running a real fan-out, not a REQ.

## Full Verbatim Input

audit adversarially the work that was done

---
*Captured: 2026-08-03T16:53:42Z*
