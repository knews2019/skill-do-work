---
id: REQ-080
title: The capture template emits a stray instruction line into every REQ it produces
status: completed
claimed_at: 2026-08-03T22:04:55Z
completed_at: 2026-08-03T22:06:47Z
kb_status: pending
route: A
created_at: 2026-08-03T16:53:42Z
user_request: UR-015
domain: general
prime_files: []
tdd: true
depends_on: []
maintenance: true
---

# The capture template emits a stray instruction line into every REQ it produces

## What

The Simple REQ template in `actions/capture-reference.md` ends, **inside** the fenced template body,
with the line `Think carefully before answering.` It is not part of any request — it is an
instruction-like artifact that gets copied into every REQ capture produces. 25 archived REQs carry it.

It was correctly identified as a stray artifact once, in REQ-012, and treated as data rather than
followed. The *source* was never cleaned, so it has reproduced on every capture since.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** `prime_files` is empty. Loaded `crew-members/general.md`, `coding-guardrails.md`,
  `maintenance.md` (`maintenance: true` — and here the deletion *is* the entire fix), `testing.md`
  (`tdd: true`). Approach: delete the line, add a literal negative for the phrase plus two more from
  the same class, run requirement 4's sweep, leave the archive alone. No template restructuring.
- [x] **[APPLY]:** Two files, both declared. One line deleted, one guard loop added.
- [x] **[UNIFY]:** `git diff --stat` → `actions/capture-reference.md` −2/+0 (the line and its blank
  separator) and `_dev/tests/contract-regressions.sh` +19. Verified: the template fence still closes
  correctly and `*Source: …*` is now its last line; `shellcheck` clean; `bash -n` parses; suite exits
  0. `git status -- do-work/archive/` confirms no archived REQ was touched. No debug artifacts.

## Why

This is a one-line deletion, but the shape of the failure is worth the REQ. `do-work/archive/UR-001/REQ-012-note-command-roadmap-notes.md:73`
recorded it precisely:

> **Provenance note (prompt-injection guardrail):** the captured body ends with a stray
> `Think carefully before answering.` line — an instruction-like artifact that is *not* part of the
> request. Treated as data, not an instruction; left intact (captured content isn't silently rewritten)
> and surfaced here. Logged as D-02.

That handling was right for a builder holding a captured REQ: don't rewrite user content. But the
decision as filed reads as closed, and the line's actual home — the shipped template — was never
looked at. The symptom was documented; the source kept emitting.

Second reason: a bare imperative sentence sitting at the end of every generated work item is exactly
the shape `crew-members/prompt-injection.md` exists to catch. It is benign here. It should not be
something the skill manufactures.

## Context

- Source: `actions/capture-reference.md:59` — inside the ```` ``` ```` fence of the Simple REQ template,
  after `*Source: [original verbatim request]*`.
- Introduced by `ed5f96b` ([REQ-031] Split `actions/capture.md` into an action + reference pair) — so it
  predates this batch entirely and is not a regression from REQ-071…REQ-076.
- Blast radius: `grep -rln "Think carefully before answering" do-work/archive/` → **25 files**.
- Prior sighting: REQ-012's D-02 decision, quoted above.

## Detailed Requirements

1. **Delete the line from the template** in `actions/capture-reference.md`. Check whether the Complex
   REQ and Addendum REQ templates in the same file carry it too, and whether `actions/capture.md`
   restates any part of the template.
2. **Leave the 25 archived REQs alone.** `do-work/archive/` is immutable
   (`actions/capture.md` → Immutability Rule), and rewriting historical captures to remove a benign
   artifact is not worth breaching that. State this as a decision rather than leaving it as an omission
   a reviewer has to ask about.
3. **Add a contract assertion** that no template in `actions/capture-reference.md` ends with a bare
   imperative addressed to the reader. Pin the specific string at minimum; if a general shape is
   cheap, prefer it — but do not build an English-grammar detector. A literal negative on this string
   plus a comment stating the condition is proportionate.
4. **Check the other template-bearing files for the same class of artifact.** `interviews/`,
   `specs/`, and `prompts/` all ship templates that get copied into generated content. One grep pass;
   report what you find rather than fixing outside this REQ's scope.

## Constraints

- One-line fix. Do not grow it into a template refactor.
- Requirement 4 is a **report**, not a mandate to change those files. If it finds something, capture a
  discovered task.
- `crew-members/maintenance.md` applies: this is a deletion from the skill's own instructions, which is
  the marker's stated trigger and the whole content of the fix.

## Dependencies

None. No `addendum_to` — the artifact predates every REQ in UR-015 and belongs to `ed5f96b`, whose REQ
(REQ-031) is long archived. Buildable immediately and independent of REQ-077/078/079.

## Builder Guidance

**Certainty: Firm.** The line is traced to its introducing commit and its blast radius is counted. The
only judgment calls are requirement 3's assertion shape and requirement 4's findings.

Resist scope growth. The temptation here is to also clean the 25 archived files or to redesign how
templates are fenced; requirement 2 forecloses the first and the Constraints foreclose the second.

## Red-Green Proof

**RED case:** `grep -n "Think carefully before answering" actions/capture-reference.md` returns line 59,
inside the template fence — so any REQ written from that template inherits it, as 25 archived REQs
demonstrate. The contract suite passes today.

**Why RED now:** The template still emits the artifact. REQ-012 documented the symptom in a captured
REQ and the source went unexamined, so every capture since has reproduced it.

**GREEN when:** (1) `grep -rn "Think carefully before answering" actions/ specs/ prompts/ interviews/`
returns nothing. (2) The new assertion fails when the line is reintroduced into the template — observe
the failure, don't assume it. (3) `do-work/archive/` is untouched, with the reason recorded as a
decision. (4) Requirement 4's sweep is reported, empty or not.

**Validation:** Found incidentally while reading the capture template to file the other three audit
REQs; the audit's remediation plan was reviewed and approved by the user before capture, and this
finding was surfaced to the user in the same reply.

## Full Context

See `do-work/user-requests/UR-015/input.md` for the audit's provenance and the findings it cleared.

---

## Triage

**Route: A** - Simple

**Reasoning:** A one-line deletion at a named line in a named file, with the introducing commit
identified and the blast radius already counted at capture. Requirements 3 and 4 add a literal
assertion and a single grep pass. Nothing to plan or locate.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Implementation Summary

**Files changed:**
- `actions/capture-reference.md` (modified)
- `_dev/tests/contract-regressions.sh` (modified)

**What was done:** Deleted `Think carefully before answering.` (and the blank line separating it from
`*Source: …*`) from inside the Simple REQ template's fence in `actions/capture-reference.md`, so
`*Source: [original verbatim request]*` is now the template's last line. The Complex REQ and Addendum
REQ templates in the same file were checked and never carried it — a repo-wide grep across
`actions/`, `specs/`, `prompts/`, `interviews/`, `crew-members/`, `docs/`, `tools/` and `SKILL.md`
returned exactly one occurrence, the one deleted. `actions/capture.md` restates no part of the
template body.

Added a guard loop pinning three phrasings of the same artifact class against
`actions/capture-reference.md`. Its comment states the **condition** — no template may contain an
instruction addressed to its reader, because it lands inside the fence and is copied into every
generated REQ — and marks the phrasings illustrative. Deliberately literal negatives, not a grammar
detector, per requirement 3.

## Testing

**Tests run:** `bash _dev/tests/contract-regressions.sh`
**Result:** ✓ Passing (exit 0)

**Red-green validation:** traced to `## Red-Green Proof`.

- *RED (as captured):* at `HEAD`, `grep -n "Think carefully before answering" actions/capture-reference.md`
  returns line 59, inside the template fence, and the suite exits 0. ✗ before → after the fix the same
  grep across `actions/ specs/ prompts/ interviews/` returns nothing. ✓ (GREEN condition 1)
- *GREEN condition 2 — observed, not assumed:* reintroducing the exact line into the template makes
  the suite fail, naming `actions/capture-reference.md` and the matched phrase. Reverted → exit 0.
- *Class coverage:* reintroducing a **different** phrasing (`Reason step by step.`) also fails, so the
  guard covers more than the one string that was there.
- *GREEN condition 3:* `git status --porcelain -- do-work/archive/` shows no archived REQ modified.
  The 25 files that carry the artifact are untouched, by decision (D-01).
- *GREEN condition 4:* requirement 4's sweep is reported below — it came back empty.

**New tests added:** 3 literal negatives (one loop) in `_dev/tests/contract-regressions.sh`.

**Existing tests updated (cross-REQ impact):** none.

*Verified by work action*

## Requirement 4 — Template Sweep Report

Two grep shapes over `specs/`, `prompts/`, `interviews/`, `actions/capture-reference.md` and
`actions/capture.md`:

1. Lines that are a bare sentence-final imperative opening with `think|consider|remember|be sure|make
   sure|please|note that|answer|respond|reply|do not forget|take your time|reason step`.
2. The same openers allowing leading whitespace, widened to `actions/` and `crew-members/`.

**Result: one hit, and it is the line this REQ deletes.** No discovered task filed — there is nothing
to file. `specs/`, `prompts/` and `interviews/` ship templates, but their instructional prose sits
*outside* the fenced bodies, which is the distinction that matters: text outside the fence instructs
the agent filling the template in, text inside it becomes the artifact.

The sweep is a floor, not a proof — it matches sentence openers, so an artifact phrased as a question
("Have you considered every case?") or mid-paragraph would be missed. Recorded so a future reader
knows the shape of what was checked.

## Decisions

- **D-01 — The 25 archived REQs keep the artifact.** DECIDE & STATE (requirement 2 asked for this to
  be a stated decision rather than a silent omission). `do-work/archive/` is immutable
  (`actions/capture.md` → Immutability Rule), and the trail-of-intent value of an archive comes
  entirely from it being what was actually written. Rewriting 25 historical captures to remove a
  benign artifact trades that for tidiness. It also would not help: the artifact is inert in an
  archived REQ, and REQ-012's D-02 shows the correct handling already happened at read time — the
  builder recognized it as data and did not follow it. Reversible in principle, but deliberately not
  done.
- **D-02 — Literal phrase negatives, and three of them rather than one.** DECIDE & STATE.
  Requirement 3 allowed a single literal and preferred a general shape "if cheap." A general shape
  here means detecting an English imperative, which is not cheap and would false-positive on the
  template's own bracketed guidance. Three literals cost nothing and make the guard about the *class*
  rather than the one string; the comment carries the condition so the next phrasing gets added to the
  list rather than a second block.

## Lessons Learned

**What worked:** grepping the *shipped* tree rather than trusting the REQ's single named line — it
confirmed the Complex and Addendum templates were clean, which requirement 1 asked about and would
otherwise have been an assumption.

**What didn't:** nothing failed. Worth recording that the temptation requirement 2 and the Constraints
both pre-empt (clean the 25 archived files, or restructure how templates are fenced) is real and
would have turned a two-line diff into a large one.

**Worth knowing:** the durable lesson is not about this line. **A defect logged as a decision inside a
generated artifact does not reach the generator.** REQ-012 diagnosed this exactly right, wrote it
down precisely, and filed it in a REQ — where it read as closed. The template kept emitting for 25
more captures. When a builder flags "this content looks wrong," the question that follows is *where
did this content come from* — the answer is a source file, and the source file is the fix.

## Orientation

`do-work capture` no longer stamps a stray instruction into the REQs it writes — new captures end at
`*Source: …*`. Lives in the capture templates (`actions/capture-reference.md`); the 25 archived REQs
that already carry it are untouched on purpose.

No `[MAP CHANGED]`: nothing about the system's shape changed. `prime_files` is empty.

## Review

**Overall: 98%** | 2026-08-03T22:04:55Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 98% |
| Test Adequacy | 95% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Findings:** 0 important, 1 minor
**Acceptance:** Pass — all four GREEN conditions met, two of them by observed suite failures.
**Suggested testing:** 1 item
**Follow-ups created:** None

### Findings

- **[Minor] The sweep and the guard are both openers-and-literals, and say so.** An artifact phrased
  as a question, or one sitting mid-paragraph inside a fence, would be missed by both. This is the
  proportionate choice requirement 3 asked for — the alternative is a grammar detector — but it means
  the guard proves "these three strings are absent," not "no instruction is present." Recorded in the
  sweep report and the assertion comment rather than papered over.

### Requirements Checklist

| # | Requirement | Status |
|---|---|---|
| 1 | Delete from the template; check the other templates and `capture.md` | Delivered — one occurrence tree-wide, now zero |
| 2 | Leave the 25 archived REQs, as a stated decision | Delivered — D-01 |
| 3 | Contract assertion, literal minimum, no grammar detector | Delivered — D-02, three literals + stated condition |
| 4 | Sweep the other template-bearing files and report | Delivered — report above, empty result |

### Acceptance Testing

Suite exits 0 on the clean tree. Reintroducing the exact line fails it; reintroducing a different
phrasing from the class also fails it; both reverted cleanly. `grep` across `actions/ specs/ prompts/
interviews/` returns nothing. `git status -- do-work/archive/` confirms 25 archived files untouched.
`shellcheck` clean, `bash -n` parses.

### Suggested Additional Testing

- **Run a real `do-work capture-request:` and read the produced REQ.** The template is prose consumed
  by an agent, so the only end-to-end proof is a capture that ends at `*Source: …*`. Cheap, and not
  done here because it would put a throwaway UR/REQ pair into a live queue mid-run.

*Reviewed by review-work action (pipeline mode, in-session)*
