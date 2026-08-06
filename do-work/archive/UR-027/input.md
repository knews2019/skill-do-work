---
id: UR-027
title: Stop the trivial follow-up REQ runaway — gate, effort label, cascade reroute, sweep consolidation
created_at: 2026-08-06T15:48:11Z
requests: [REQ-125, REQ-126, REQ-127]
word_count: 5
---

# Stop the trivial follow-up REQ runaway

## Summary

The review step's one-REQ-per-Important-finding rule has no relevance check and no depth limit, so reviews of review-generated REQs cascade. UR-489 (a one-hour pill feature, REQ-1305) spawned sixteen REQs over two days — fifteen of them facets of one root cause — and the user had to spend their own time discovering they were all trivial. This capture ships the three priorities agreed in the design discussion: (1) a disposition gate plus an `effort_estimate` label and board chip, (2) a cascade depth stop that reroutes generation-≥2 follow-ups to `pending-answers`, (3) sweep-REQ consolidation for same-root-cause findings. Inline-fix-at-review-resolution (priority 4) was deliberately deferred and is NOT captured here.

Guiding principle, stated by the user: keep "see something, say something" — nothing reduces what gets recorded; the changes affect labeling, routing, and consolidation only. The label is the user's stated most-important fix ("that way I can easily decide if I want to stop or not the process").

## Extracted Requests

| REQ | Title | Depends on |
|---|---|---|
| REQ-125 | Disposition gate + `effort_estimate` label + board chip | — |
| REQ-126 | Cascade depth stop: generation-≥2 follow-ups reroute to `pending-answers` | REQ-125 |
| REQ-127 | Sweep-REQ consolidation for same-root-cause findings | REQ-125 |

## Batch Constraints

- **Nothing is suppressed.** Every finding still becomes a REQ (or a sweep checklist item). Report-only outcomes were explicitly rejected by the user ("I still want all the REQs created, I just need to know their impact").
- **Severity vocabulary is untouched.** Important/Minor/Nit judgment stays as-is; the gate routes findings, it never re-scores them. State this explicitly in the shipped text so agents don't resolve tension by downgrading severities.
- **Same-commit lock-step** for any field the board parses: `tools/queue-kanban/model.go`, the board-parsed-fields enumeration in this repo's CLAUDE.md, and the Schema Read Contract in `actions/work-reference.md`.
- **Restatement sweep before any REQ here is called done:** the one-REQ-per-Important rule is restated at `actions/work.md` ~:495, ~:501, ~:505 and in `actions/review-work.md` (~:335 Step 10, ~:450, ~:466, ~:493) — all restatements move together. Line numbers refreshed after the 2026-08-06 merge of main (PR #135); re-grep, don't trust them.
- **Inline fixes are out of scope** for all three REQs — deferred by user decision until the shipped result has been lived with.
- All three REQs are `maintenance: true` — each rewrites/narrows the skill's own operating instructions; `crew-members/maintenance.md` (delete-before-you-add) loads at work time.

## Full Verbatim Input

do-work capture-request Ship priorities 1 through 3

## Referenced Conversation Context

The invocation refers to the priority ordering agreed in the preceding design discussion. Preserved verbatim below: the user's original proposal, and the decision record.

### Original proposal (user's opening message, verbatim)

**The story (why this matters):** UR-489 was a one-sentence request — "add a
has-alignment-note pill like the validation pill." The actual work (REQ-1305) was done
in an hour. But the skill's review step creates one follow-up REQ per Important
finding, with no relevance check and no depth limit, so reviews of review-generated
REQs kept spawning more: 1305 → 1307 → 1308/1309 → 1310/1311 → 1312–1318 → 1320/1321.
Sixteen REQs over two days, fifteen of them facets of ONE root cause (hardcoded colors
not tokenized + a guard blind to them). I had to invest my own time to discover they
were all trivial. That cost is the problem to fix.

**Proposed changes to the skill's follow-up-creation rules:**

1. **Good-enough-already gate (runs BEFORE any automatic REQ is created).** When a
   review or discovered-tasks step is about to create a REQ automatically, it must
   first ask: is the current state realistically good enough already? Concretely:
   (a) would any user or developer actually notice this issue in real use, and
   (b) does fixing it establish or change a rule that applies in several places in the
   codebase (a genuine maintainability rule, not a one-spot patch)?
   If NO to both → the finding is trivial: mark the REQ it lands in with the trivial
   flag (see 2), or leave it as a report-only line. Only a yes on (a) or (b) justifies
   normal-effort follow-up work.

2. **Triviality flag, visible on the board.** Add an `effort_estimate` frontmatter
   field (`trivial | normal`, default `normal`) to the REQ schema. Automatic
   follow-ups MUST set it at creation using the gate above; capture MAY set it. The
   queue-kanban board renders `effort_estimate: trivial` as a visible chip so I can
   tell at a glance which queued REQs are small mechanical fixes versus real work.

3. **Fix-inline + one-sweep-REQ rule (replaces one-REQ-per-finding).**
   - A trivial finding in a file the current REQ already touches is fixed inline
     during review resolution and recorded in the report — no new REQ at all.
   - Trivial or same-theme findings elsewhere never get individual REQs: append to an
     existing sweep REQ for that theme under the same UR if one is queued, otherwise
     create ONE consolidated sweep REQ named for the ROOT CAUSE (e.g. "tokenize all
     remaining hardcoded colors and make the guard catch every notation") with a
     checklist of instances. Solving the sweep means the class of finding cannot
     recur — the rule is changed everywhere it applies — not that N spots got patched
     one drop at a time.
   - Only a genuinely non-trivial, thematically unrelated finding still earns its own
     REQ, and it must state in one line why it couldn't fold into a sweep.

4. **Cascade depth stop.** A `review_generated: true` REQ is generation ≥2: its own
   review may fix inline and append to existing sweeps, but may NOT create brand-new
   REQs — leftovers go in the review report for me to triage. This hard-stops the
   UR-489 chain shape at depth 2.

Keep the existing severity language (Important/Minor/Nit); this changes where findings
land, not how they're judged.

### Decision record (agreed during the discussion)

1. **Inline fixer:** orchestrator applies inline fixes at review resolution, not the review agent — review stays read-only. (Then priority 4 / inline fixes as a whole was deferred out of this capture entirely.)
2. **Critical pierce:** yes — critical-grade findings (security, data loss, broken production path) create REQs at any depth. User's words: "yes, do it, I still want all the req's created, I just need to know their impact, so categorization of critical is definitely useful."
3. **Sweep lookup:** frontmatter marker (`sweep: true`), not judgment-matching on titles.
4. **Depth-stop amendment:** generation-≥2 leftovers become `pending-answers` REQs (visible, chipped, human-gated via `do-work clarify`) instead of report-only lines — reconciling the depth stop with "I still want all the REQs created."
5. **Reprioritization (user's final framing):** "the most important fix is the label, that way I can easily decide if I want to stop or not the process. I do like the fact that the current system is working on the 'see something, say something' principle, this way we can discover and fix issues, I just have a problem with the trivial run away."
6. **Scope:** ship priorities 1–3 as three chained REQs (label+gate → reroute → sweep); priority 4 (inline fixes) deferred, decision on record.

---
*Captured: 2026-08-06T15:48:11Z*
