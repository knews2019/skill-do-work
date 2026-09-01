---
id: UR-088
title: 'Lessons routing with token-budgeted mandatory reads and a fold-gate fix'
created_at: 2026-09-01T10:47:44Z
requests: [REQ-477, REQ-478, REQ-479, REQ-480]
word_count: 753
---

# Lessons Routing with Token-Budgeted Mandatory Reads and a Fold-Gate Fix

## Summary

Two orchestrator improvements decided interactively during the 2026-09-01 queue-analysis session. (1) Make lessons transfer instead of merely accumulating: an intelligent index over the lesson satellites (path + when-it-applies hook + token estimate), capture-time stamping of mandatory lessons reads on new REQs under a token budget, enforced reads in the work pipeline, family-slugged lesson bullets with mandatory twice-seen Trap promotion, and a prime-audit safety net. (2) Widen the Fold-First Rule's destination 2 so non-critical findings can fold into a dependency-gated REQ whose chain is alive, instead of minting duplicates.

## Extracted Requests

| REQ | Request | Depends on |
|---|---|---|
| REQ-477 | Family-keyed lessons, intelligent index, and mandatory twice-seen Trap promotion | — |
| REQ-478 | Capture stamps required lessons under a token budget | REQ-477 |
| REQ-479 | Enforce required-lessons reads and audit un-promoted families | REQ-478 |
| REQ-480 | Allow fold-first conversion into live dependency-gated destinations | — |

## Batch Constraints

- Plain files only for the index and its consumers — capture and builders on the floor agent (read/write files + shell) must be able to use the mechanism with read/grep alone; no new tooling dependency.
- The token-budget constant is stated once in one findable place, never scattered.
- Both lesson-routing mechanisms deliberately coexist: capture-time stamping covers stamped REQs; twice-seen Trap promotion covers everything else. Do not trade one away while implementing the other.
- Do not backfill old lesson entries wholesale; seed slugs and index hooks only for the three known recurring families as worked examples.
- The maintainer chose "Keep both" (index + promotion) and "capture into the queue" explicitly; do not relitigate.
- Evidence baseline: failure families final-boundary identity/rollback (REQ-413/436/447/463/416), opaque-evidence/generic-fallback projections (REQ-414/430/446), and smoke-vs-characterization gates (REQ-414/415) each recurred while their lessons sat in skills/do-work/tools/do-work-cli/lessons-do-work-cli.md; REQ-464/REQ-465 duplication and manual fold recorded in commit 593c5145.

## Full Verbatim Input

> ```
> Request 1 — lessons routing and transfer:
> 
> Make lessons transfer instead of merely accumulating, by routing them at capture time and enforcing the read. Five coupled changes:
> 
> 1. Intelligent lessons index. A single well-known plain-markdown index (suggested location do-work/lessons-index.md, precedent do-work/prose-backlog.md; builder may argue a better home) with one line per lessons satellite — path, a when-it-applies hook naming the failure families inside (e.g. "rollback/deletion/final-boundary work in do-work-cli internals"), and an approximate token estimate (mechanical, e.g. bytes/4 — the exact formula is the builder's call but must be reproducible by a floor agent with wc). Maintained by work.md's Lessons-Capture Phase: whenever a lesson is appended (Step 8's satellite write), the same edit creates or refreshes that satellite's index line, estimate included. Plain files only — capture and builders on the floor agent must be able to use it with read/grep alone.
> 
> 2. Capture stamps mandatory lessons under a token budget. capture-request, while authoring REQ payloads (capture.md Step 5), reads the index and decides which lessons files are relevant to the request being captured, stamping them in a new frontmatter field (suggested required_lessons: [paths]) on each minted REQ — but the stamped set's summed index estimates must stay within a stated budget (suggested ~2000 tokens per REQ; the REQ decides the number and where it is stated so it is one findable constant, not scattered). Over budget, capture ranks by relevance and stamps the best-fitting subset; because lesson bullets carry family slugs (item 4), capture may stamp a targeted reference (path plus family slug) so the builder greps only the relevant bullets instead of reading the whole satellite — the cheapest way to stay in budget without dropping a match. What was considered and dropped is noted in the REQ body, never silently. Empty/absent when nothing matches — never invented. Add the field to the Simple REQ template (actions/capture-reference.md) and the Request File Schema (actions/work-reference.md); verify whether internal/requestmodel and internal/schemanormalization must learn the field for lossless preservation (unknown-field preservation may already cover it), and whether the board needs anything (display optional; parser lock-step rule in _dev/primes/prime-kanban-board.md governs if it does).
> 
> 3. Enforced read. actions/work.md Step 5 (explore context) and Step 6 (builder instructions) make reading every required_lessons entry mandatory before implementation — a bare path means the whole satellite, a targeted path-plus-slug reference means the matching bullets — unconditional, unlike today's touch-conditional rule at work.md:404, which stays for unstamped REQs. A missing listed file is proceed-without-it, per the existing missing-rules-file convention.
> 
> 4. Family-keyed lessons plus kept Trap promotion. Every appended lesson bullet carries a short kebab-case failure-family slug (the index hooks and recurrence checks key on it; same discriminator move as sweep_key). The Step 8 lesson write scans the satellite first: on the second same-family occurrence, promotion stops being optional — the writer adds or amends one generalized Trap line in the owning prime's Traps section in the same edit (work.md:604's "judgment call" sentence becomes condition-keyed). Both mechanisms deliberately coexist: capture-time routing covers stamped REQs, Traps cover everything else.
> 
> 5. Audit safety net. Extend skills/do-work-toolbox/actions/prime.md audit to flag: a satellite with 2+ same-family entries and no corresponding Trap line; a satellite missing from the index; an index line whose path is dead.
> 
> Do not backfill old lesson entries wholesale; seed slugs and index hooks for the three known recurring families as worked examples. Update crew-members/general.md Lessons Discipline restatements in the same change. Evidence: three failure families each recurred 3-6 times across the 2026-08-31 run while their lessons sat unfound in the satellite — final-boundary identity/rollback (REQ-413/436/447/463/416), opaque-evidence/generic-fallback projections (REQ-414/430/446), smoke-vs-characterization gates (REQ-414/415); REQ-415 repeated a family whose lesson REQ-414 had already recorded.
> 
> Request 2 — fold gate:
> 
> Amend the Fold-First Rule's destination 2 (skills/do-work/actions/capture-reference.md, Fold-First Rule) so a matching non-sweep REQ that is dependency-gated can still receive a fold when all of: (a) the root cause matches, (b) the destination's dependency chain is alive — every depends_on member is terminal-successful or present in the queue/working set (a failed or cancelled member keeps the destination ineligible), and (c) the finding's judged impact is not impact-critical (critical keeps current behavior; the existing escalation rule still applies to folds). Keep the unassigned (assigned_to absent) requirement unchanged. Update any restatement of the eligibility condition (actions/review-work.md Step 10, actions/work-reference.md follow-up flows — they cite the rule by name; verify none restates the gate), and add a contract-regression predicate if the builder judges the wording load-bearing. Evidence: REQ-464 and REQ-465 were minted solely because REQ-420 was dependency-gated at capture time, then manually folded into REQ-420 and cancelled in commit 593c5145.
> ```

---
*Captured: 2026-09-01T10:47:44Z*
