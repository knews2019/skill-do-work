---
id: REQ-080
title: The capture template emits a stray instruction line into every REQ it produces
status: pending
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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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
