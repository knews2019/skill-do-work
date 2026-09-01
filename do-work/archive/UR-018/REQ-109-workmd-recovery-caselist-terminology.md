---
id: REQ-109
title: "work.md session-start note still enumerates the recovery case list and calls a label-less entry a foreign claim"
status: completed
status_changed_at: 2026-08-05T11:50:17Z
claimed_at: 2026-08-05T11:51:58Z
completed_at: 2026-08-05T11:59:04Z
commit: 5f50fb7
route: A
kb_status: promoted
kb_entry: REQ-109-work-md-session-start-note-still-enumera.md
created_at: 2026-08-05T11:44:27Z
user_request: UR-018
addendum_to: REQ-108
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: []
maintenance: true
write_set: [actions/work.md]
related: [REQ-104, REQ-108]
batch: parallel-building
---

# work.md Session-Start Note: Recovery Case List Terminology

## What

Discovered during REQ-108 (`[low]`): `actions/work.md`'s Step 10 session-start note (the sentence
listing which `working/` REQs recovery may strip) carries the same closed-enumeration shape REQ-108
just removed from `actions/work-reference.md`'s In-Progress Record — "one that isn't (unlabeled, or
labeled for another checkout) is a foreign claim recovery must not strip." The set is complete and
the behavior is correct, but it calls the label-less case a *foreign claim*, whereas since REQ-104
the canonical term is a *claim of unknown origin* — and the very next sentence in the same note uses
the correct term. Same fix shape as REQ-108: state the condition, defer the list to Crash Recovery.

## Open Questions

- [x] I discovered this out-of-scope task while working on REQ-108: `actions/work.md`'s session-start note restates the recovery case list with pre-REQ-104 terminology (calls a label-less entry a "foreign claim"; behavior is correct, wording predates the drop). Should I process this as a new task? → Confirmed: Yes, add to queue. Trivial fix in a file that already uses the canonical term two lines later; a stale restatement in the pipeline's most-read action file is the drift class REQ-104/108 exist to remove. (User directed common-sense resolution via clarify.)
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.

---
*Source: REQ-108 builder, Discovered Tasks ([low])*

---

## Triage

**Route: A** - Simple

**Reasoning:** Names the specific file (`actions/work.md`), the specific sentence (Step 10 session-start note, item 2), and the exact fix shape already proven by REQ-108: state the classification condition, defer the case list to Crash Recovery, and drop the pre-REQ-104 "foreign claim" label for the label-less case. Terminology-only; behavior is correct as written.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

## Decisions

- **D-01**: Builder chose: adopt REQ-108's "left byte-identical, however it fails to match that label" phrasing rather than paraphrasing "recovery must not strip." Reasoning: `byte-identical` is the term Crash Recovery and the In-Progress Record section both use for that outcome and is strictly stronger (no strip, no frontmatter reset, no move); reusing it makes the two passages read as one contract. Kept `recovery must leave` as the verb so the actor stays explicit.
- **D-02**: Builder chose: cite the enumeration inline as "(the cases are enumerated once, in `actions/work-reference.md` → **Crash Recovery (Step 1)**)" rather than as a trailing sentence. Reasoning: the parenthetical carries the same pointer without growing the item, and keeps the "enumerated once" framing that tells a reader not to restate the list here.
- **D-03**: Builder chose: no bold added to the new clause. Reasoning: item 2 carries no bold today; adding emphasis would be an unrequested style change.

## Implementation Summary

**Files changed:**
- `actions/work.md` (modified)

**What was done:** Rewrote the first sentence of item 2 in Step 10's session-start note: it now states the own-label condition (own `writer:` label → this session's own to recover; any other entry left byte-identical) and defers the case enumeration to `actions/work-reference.md` → **Crash Recovery (Step 1)**, dropping both the hand-enumerated case list `(unlabeled, or labeled for another checkout)` and the pre-REQ-104 "foreign claim" label for the label-less case. The second sentence (entries written at claim time by Step 2) is unchanged. Builder also verified the rest of the note — items 1, 3, line 646, line 658, and the checklist line 668 — uses the canonical terms correctly; no other drift found.

## Qualification

Passed — `tools/checks/qualify.sh` mechanical checks OK; orchestrator read the diff directly: 1 file verified, all 3 REQ requirements traced (condition stated, enumeration deferred to Crash Recovery, "foreign claim" mislabel removed from the label-less case), change is substantive-for-its-kind (terminology contract fix, second sentence untouched). No P-A-U section on this REQ (discovered-task follow-up — capture never wrote one), so no box audit applies.

## Testing

- `bash _dev/tests/contract-regressions.sh` — **pass** ("Contract regression checks passed", including SKILL.md word budget and shipped-path citation checks that cover `actions/work.md`).
- Red-green validation omitted — non-behavioral change (instruction-prose terminology fix); regression evidence used instead.

## Review

**Overall: 95%** | 2026-08-05T11:57:33Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 90% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Findings:** 0 important, 1 minor
**Acceptance:** Pass — the rewritten item 2 states the own-label condition, its citation target (`actions/work-reference.md` → **Crash Recovery (Step 1)**, line 240) resolves, the four canonical cases still cover every classification the sentence defers, and `bash _dev/tests/contract-regressions.sh` was re-run by the reviewer and passed.
**Suggested testing:** 0 items
**Follow-ups created:** None

### Requirements Checklist

- [x] **R1 — item 2 states the own-label condition instead of enumerating the non-matching cases** — delivered. The hand-enumerated `(unlabeled, or labeled for another checkout)` is gone; the sentence now reads "a `working/` REQ named there under this checkout's own `writer:` label is this session's own to recover; any other entry recovery must leave byte-identical, however it fails to match that label." Condition stated once, complement left open.
- [x] **R2 — the case list is deferred to Crash Recovery (Step 1) by citation** — delivered, as the inline parenthetical "(the cases are enumerated once, in `actions/work-reference.md` → **Crash Recovery (Step 1)**)". The target section exists and holds all four bullets.
- [x] **R3 — the label-less case is no longer called a "foreign claim"** — delivered. The term is absent from the rewritten sentence entirely; it no longer applies any label to the label-less case, which is stronger than swapping in "claim of unknown origin" and matches how item 3 two lines below already distinguishes the two.
- [x] **Constraint — second sentence of item 2 unchanged** — verified byte-for-byte in the diff: "Step 2 wrote those entries at claim time, one per claim (`actions/work-reference.md` → **In-Progress Record (Step 2)**), which is why the section survives a crash that never reached this step."

### Findings

**Important:** None.

**Minor:**
- `actions/work.md:774` (Verification Checklist) — "No REQ files remain in `do-work/working/` after the work loop ends — except a reported foreign claim this run deliberately left intact (Step 1 Crash Recovery)". Under the post-REQ-104 vocabulary a label-less checkpoint entry is a *claim of unknown origin*, not a foreign claim, and Crash Recovery leaves its `working/` file byte-identical too — so this exit criterion is under-inclusive and would read a correctly-left label-less claim as a checklist violation. Same drift class REQ-104/108/109 exist to remove, in the same file, but outside this REQ's stated sentence-level scope. Suggested fix is one phrase: "except a claim Step 1 deliberately left intact". Not scored as builder scope drift (review-work.md Step 6 Restatement Sweep, item 3).

**Nit:**
- `actions/work.md:655` — "any other entry recovery must leave byte-identical, however it fails to match that label" uses *however* in its "in whatever way" sense, but the preceding comma invites the contrastive reading on a first pass. "whichever way it fails to match that label" carries the same meaning with no garden path. Zero score impact.
- `actions/work.md:655` — "recovery must leave byte-identical" states the default without the human-approved-takeover exception (`actions/work-reference.md` line 272: recovery also runs on "a foreign claim a human approved for takeover"). The prior wording ("recovery must not strip") had the identical gap, so this REQ neither introduces nor widens it, and the deferral citation carries the reader to the ladder. Noted only so the next editor of this line knows the absolute framing is deliberate shorthand.

### Restatement Sweep

**Result: one stale restatement found (Minor, above); everything else agrees.**

Swept `actions/`, `crew-members/`, `docs/`, `interviews/`, `prompts/`, `tools/**/*.md`, `SKILL.md`, `next-steps.md`, `README.md` for `foreign claim`, `claim of unknown origin`, `unlabeled`, `label-less`, `writer:`, `byte-identical`, `takeover`, and `CHECKPOINT`. Verified as still correct:

- `actions/work.md:125` (Step 1 summary) — condition-shaped, applies no label to the label-less case.
- `actions/work.md:646` and `:656` (Step 10 write / delete gate) — enumerate "a foreign `writer:` label and no label at all" as two distinct cases of *entries this checkout did not write*, and `:656` names the second a "claim of unknown origin". Canonical.
- `actions/work.md:658` — no-checkpoint case called a foreign claim; matches Crash Recovery's fourth bullet. Not a finding.
- `actions/work.md:668` (checklist) — "never take over a labeled foreign claim" matches the second bullet's never-enters-the-ladder rule.
- `actions/work-reference.md:61, 246–253, 297, 349, 458–466`, `actions/cleanup.md:40`, `actions/forensics.md:39, 200` — all already carry the post-REQ-104/108 wording.
- `docs/work-guide.md:66, 68, 116, 118, 122` — user-facing prose is condition-shaped throughout ("records it as that session's own interrupted work — under this checkout's own stamp; anything else is left exactly as it is"); it never enumerates the case list and never applies "foreign" to a label-less entry.

*Reviewed by review-work action*

## Lessons Learned

**What worked:** Reusing REQ-108's exact fix shape (state the condition, defer the enumeration, reuse the canonical "byte-identical" term) made the two passages read as one contract instead of two paraphrases.
**What didn't:** N/A — no dead ends on a one-sentence fix.
**Worth knowing:** The review's restatement sweep found a fourth instance of this same drift class in the same file (`actions/work.md:774`, Verification Checklist — "a reported foreign claim" is under-inclusive for the label-less/unknown-origin case). The pattern: REQ-104 changed the classification vocabulary, and each sweep since has caught one more stale restatement (REQ-108 → work-reference.md, REQ-109 → work.md line 655, now line 774). When a vocabulary changes, grep for the old term across *every* shipped file in the first fix, not one file per follow-up.

## Orientation

Terminology-consistency fix in the work pipeline's session-start instructions (`actions/work.md`, Step 10 note): the crash-recovery classification is now stated once as a condition with the case list deferred to `actions/work-reference.md`'s Crash Recovery — closing UR-018's last live member. No map change.
