---
id: UR-089
title: 'v4 triage: lessons-routing refinements and fold-gate shape condition'
created_at: 2026-09-01T11:04:54Z
requests: [REQ-481]
word_count: 1827
---

# v4 Triage: Lessons-Routing Refinements and Fold-Gate Shape Condition

## Summary

A v4 revision of the UR-088 plan, run through `do-work validate-feedback`. The triage adjudicated ten deltas: seven accepted, one pushed back (re-slicing the committed REQ-477–479 — churn exceeds gain), two already done. Accepted refinements land as addenda to the four pending UR-088 REQs; one new REQ carries the one-time queue stamping pass. Also durably recorded here: the generation-1 consent-gate change is out of scope (decided earlier, restated so it is not lost).

## Extracted Requests

| Destination | Content | Triage |
|---|---|---|
| REQ-477 (addendum) | Hook names family slugs; `slugged: full|partial` coverage flag | Finding 1, Accept |
| REQ-478 (addendum) | Entry forms, full-only targeting, targeted-entry cost rule, relevance ranking, fixed dropped-heading; required round-trip schema predicate | Findings 2+5, Accept |
| REQ-479 (addendum) | Enforced read is additive for all REQs; miss recorded in hand-back; three further mechanical audit checks | Findings 4+6, Accept |
| REQ-480 (addendum) | Goal-shaped vs defect-shaped condition (b); fold recorded in destination with date/source | Finding 8, Accept |
| REQ-481 (new) | One-time required-lessons stamp of the pending queue | Finding 3, Accept |
| — | Re-slice REQ-477–479 along v4 item 7's edges | Finding 7, Push back: already minted and committed (2ee8bee8); only gain is audit parallelism |
| — | REQ-464/465 recurrence de-count; `git log -1` verification race | Findings 9+10, Already done in the captured artifacts and the session's verification |

## Folded Requests

- REQ-477 — index hook slug list and `slugged: full|partial` coverage flag (v4 item 1)
- REQ-478 — stamping entry forms, full-only targeting, cost rule, relevance ranking, fixed dropped-heading (v4 item 2), and the required schema round-trip predicate (v4 item 2)
- REQ-479 — additive enforced read with hand-back miss record (v4 item 4) and the three added audit checks (v4 item 6)
- REQ-480 — goal-shaped condition (b), re-lettered (c)/(d), and the gated-destination fold record (v4 request 2)

## Batch Constraints

- Out of scope, restated so it is not lost: the generation-1 consent-gate change (decided earlier in the 2026-09-01 session).
- The push-back on re-slicing stands: REQ-477–479 keep their committed boundaries; REQ-481 attaches to the existing DAG via `depends_on: [REQ-478]`.

## Full Verbatim Input

> ```
> do-work validate-feedback: (check what we can improve based on the following prompt)
> 
> # Orchestrator Improvement: Lessons Routing + Fold-Gate Fix (capture into queue) — v4
> 
> ## Context
> 
> The 2026-09-01 queue analysis found agents relearning the same failure families across the run: "final-boundary identity/rollback" recurred across REQ-413/436/447/463/416 (REQ-457 still queued for it), "opaque-evidence/generic-fallback projections" across REQ-414/430/446, "smoke-vs-characterization" across REQ-414/415. (REQ-464/465 also showed the latter two families but were duplicates of REQ-420's goals — see request 2 — so they are not counted as independent recurrences.) The lessons were on file in skills/do-work/tools/do-work-cli/lessons-do-work-cli.md the whole time — REQ-415 repeated a family whose lesson REQ-414 had already recorded. Root cause: lessons capture works, transfer doesn't. Today a builder reads a lessons satellite only conditionally ("when your change touches code the prime's Read-first or Traps entries name", skills/do-work/actions/work.md:404), bullets carry no family key so recurrence is invisible, and Trap promotion is an unowned judgment call (work.md:604).
> 
> Chosen design: route lessons at capture time through an intelligent index under a per-REQ reading budget, make the reads mandatory per REQ, and keep twice-seen Trap promotion as redundant coverage. Separately, the Fold-First Rule's dependency-ready gate minted duplicate REQs (464/465, folded and cancelled by hand in commit 593c5145) and gets a scoped, shape-aware fix. Both changes go through the queue via do-work capture-request — an orchestrator session is actively working the queue, and instruction changes ride the pipeline like any other work.
> 
> Out of scope (decided earlier, restated so it isn't lost): the generation-1 consent-gate change.
> 
> Answered along the way, no standing changes: no new prime files are needed (all current, per-touched-prime staleness spot-checks already run each REQ), and prime audit should run at milestones (post-REQ-420, post-UR-087), not on a schedule — request 1 makes it the safety net for the new rules.
> 
> ## Plan
> 
> One step: run the capture pipeline (skills/do-work/actions/capture.md, do-work capture-request) with the two-request input below. Capture owns Step 1 assessments (impact/effort/maintenance), its fold-first scan, UR + REQ minting/slicing, and the bookkeeping commit — do not hand-author REQ files or pre-judge frontmatter the action owns.
> 
> ## Capture input — request 1 (lessons routing and transfer)
> 
> > Make lessons transfer instead of merely accumulating, by routing them at capture time under a reading budget and enforcing the read. Seven coupled changes:
> >
> > 1. **Intelligent lessons index.** A single well-known plain-markdown index (suggested do-work/lessons-index.md, precedent do-work/prose-backlog.md; builder may argue a better home) with one line per lessons satellite carrying: the path; a when-it-applies hook that names the failure-family slugs inside (e.g. "rollback/deletion/final-boundary work in do-work-cli internals — final-boundary-identity, opaque-evidence-projection"); a mechanical size estimate (e.g. bytes/4 — exact formula is the builder's call but must be reproducible by a floor agent with wc); and a slug-coverage flag (`slugged: full` when every bullet carries a family slug, `partial` otherwise). Maintained by work.md's Lessons-Capture Phase: whenever a lesson is appended (Step 8's satellite write), the same edit creates or refreshes that satellite's index line — hook slugs, estimate, and coverage flag included. Plain files only — capture and builders on the floor agent must be able to use it with read/grep alone.
> >
> > 2. **Capture stamps mandatory lessons under a token budget.** capture-request, while authoring REQ payloads (capture.md Step 5), reads the index and stamps relevant lessons in a new frontmatter field (suggested `required_lessons`, a list of strings) on each minted REQ. Two entry forms: a bare `path` (whole satellite) or `path#family-slug` (targeted; only the bullets carrying that slug). Targeted entries are permitted only for satellites whose index line says `slugged: full`; a `partial` satellite is stamped bare or not at all, so an un-slugged bullet can never be grepped past. The stamped set's summed cost must stay within a stated budget (suggested ~2000 tokens per REQ; the REQ decides the number and where it lives so it is one findable constant). Cost rule, stated next to the constant: a bare path costs its index estimate; a targeted entry costs the size of its matching bullets (grep the slug, wc the lines, same formula as the index). Relevance ranking, in order: (a) the request text or its likely touched paths match a family named in the hook; (b) the satellite's owning prime governs a path the request names; (c) most recent same-family recurrence. Over budget, capture prefers narrowing a match to a targeted entry over dropping it; anything still dropped is listed in the REQ body under a fixed heading, never silently. Empty/absent when nothing matches — never invented. Add the field to the Simple REQ template (actions/capture-reference.md) and the Request File Schema (actions/work-reference.md); add a round-trip predicate to _dev/tests/contract-regressions.sh proving `required_lessons` survives internal/requestmodel/internal/schemanormalization unchanged (do not assume unknown-field preservation covers it — prove it); board display optional, parser lock-step rule in _dev/primes/prime-kanban-board.md governs if added.
> >
> > 3. **One-time stamp of the pending queue.** When item 2 lands, run the same stamping decision once over every pending, unassigned REQ already in do-work/queue and stamp them (REQ-457 is the motivating case — it is queued for a family whose lesson already exists). Each retroactively stamped REQ notes that in its body. A single pass, not a standing behavior.
> >
> > 4. **Enforced read.** actions/work.md Step 5 (explore context) and Step 6 (builder instructions) make reading every `required_lessons` entry mandatory before implementation — bare path means the whole satellite, `path#slug` means the matching bullets — unconditional. Today's touch-conditional rule at work.md:404 stays in force for all REQs, stamped or not; the mandatory read is additive, not a replacement, and Step 6 says so in one sentence so a builder never guesses which regime applies. A missing listed file is proceed-without-it, per the existing missing-rules-file convention, and the miss is recorded in the hand-back.
> >
> > 5. **Family-keyed lessons + kept Trap promotion.** Every appended lesson bullet carries a short kebab-case failure-family slug (the index hooks, targeted entries, and recurrence checks key on it; same discriminator move as sweep_key). The Step 8 lesson write scans the satellite first: on the second same-family occurrence (by slug, or same-family judgment for pre-slug entries), promotion stops being optional — the writer adds or amends one generalized Trap line in the owning prime's ## Traps in the same edit, citing the slug (work.md:604's "judgment call" sentence becomes condition-keyed), and notes the promotion in the hand-back. Both mechanisms deliberately coexist: capture-time routing covers stamped REQs, Traps cover everything else.
> >
> > 6. **Audit safety net.** Extend skills/do-work-toolbox/actions/prime.md audit to flag, all mechanically: a satellite with 2+ same-family entries and no corresponding Trap line; a satellite missing from the index; an index line whose path is dead; an index estimate more than ~25% off the recomputed value; an index hook whose slug set does not match the slug set actually present in the satellite (either direction); a `slugged: full` flag on a satellite that still has un-slugged bullets.
> >
> > 7. **Seeding, restatements, slicing.** Do not backfill old lesson entries wholesale; seed slugs, index hooks, and coverage flags for the three known recurring families as worked examples (which makes lessons-do-work-cli.md's coverage flag honest, whatever it ends up being). Update crew-members/general.md § Lessons Discipline restatements in the same change. If capture slices this request, the natural edges are: schema field first; stamping (2) and enforced read (4) depend on it; the one-time queue pass (3) depends on stamping; audit (6) is independent.
> >
> > Evidence: the three families above each recurred 3–6 times across the 2026-08-31 run while their lessons sat unfound in the satellite; REQ-415 repeated a family whose lesson REQ-414 had already recorded.
> 
> ## Capture input — request 2 (fold gate)
> 
> > Amend the Fold-First Rule's destination 2 (skills/do-work/actions/capture-reference.md § Fold-First Rule) so a matching non-sweep REQ that is dependency-gated can still receive a fold when all of: (a) the root cause matches; (b) the finding is goal-shaped, not defect-shaped — it restates, refines, or extends the destination's acceptance goals rather than reporting behavior that is broken in what currently ships; a shipped-behavior defect keeps today's behavior (minted standalone) regardless of impact, because a fix must never wait behind a gate; (c) the destination's dependency chain is alive — every depends_on member is terminal-successful or present in the queue/working set (a failed or cancelled member keeps the destination ineligible); and (d) the finding's judged impact is not impact-critical (critical keeps current behavior; the existing escalation rule still applies to folds). Keep the unassigned (assigned_to absent) requirement unchanged. A fold accepted into a gated destination is recorded in the destination body with the fold date and source, so a stalled chain carrying folds is visible to review-work Step 10 (and to the board, if it already surfaces gated REQs). Update any restatement of the eligibility condition (actions/review-work.md Step 10, actions/work-reference.md follow-up flows — they cite the rule by name; verify none restates the gate), and add a contract-regression predicate pinning conditions (b)–(d) if the builder judges the wording load-bearing. Evidence: REQ-464/REQ-465 were minted solely because REQ-420 was dependency-gated at capture time, then manually folded into REQ-420 and cancelled in commit 593c5145; both were goal-shaped (duplicates of REQ-420's acceptance goals), which is exactly the case (b) admits.
> 
> ## Expected capture outcomes
> 
> Capture's judgment, not directives: request 1 likely maintenance: true, impact-rule-change, effort-substantive, likely sliced along the edges named in item 7; request 2 maintenance: true, impact-rule-change, mechanical-to-substantive. Version bump + changelog belong to the eventual integrating commits, not to capture.
> 
> ## Files involved
> 
> - Read/execute now: skills/do-work/actions/capture.md (+ capture-reference.md templates)
> - Created by capture: do-work/user-requests/UR-0XX/input.md, do-work/queue/REQ-4XX-*.md
> - Eventually modified by the implementing REQs (not this session): skills/do-work/actions/work.md, skills/do-work/actions/capture.md, skills/do-work/actions/capture-reference.md, skills/do-work/actions/work-reference.md, skills/do-work/actions/review-work.md, crew-members/general.md, skills/do-work-toolbox/actions/prime.md, do-work/lessons-index.md (new), _dev/tests/contract-regressions.sh, possibly internal/requestmodel/internal/schemanormalization and the kanban parser
> 
> ## Verification
> 
> - Capture completes with UR id(s) and REQ ids; `git log --oneline -n 5` includes capture's bookkeeping commit and `git show --stat` on that commit touches only queue/user-request paths (do not rely on `-1` — the active orchestrator may commit in between).
> - The board shows the new REQs pending; the active orchestrator builds them through the normal loop.
> - No edits to skills/ land in this session — the pipeline makes them.
> 
> ---
> 
> ## Delta from v3 (for review; not part of the capture input)
> 
> - Index line: hook now lists family slugs; new `slugged: full|partial` coverage flag.
> - Stamping: entry format fixed as `path` / `path#slug`; targeted entries only on fully-slugged satellites; cost rule for targeted entries; explicit relevance ranking; dropped matches listed under a fixed heading.
> - New item 3: one-time stamping pass over the pending queue (REQ-457).
> - Enforced read: touch-conditional rule stays for all REQs (additive), stated in Step 6; missing file recorded in hand-back.
> - Schema: "verify whether" → required round-trip contract-regression predicate.
> - Audit: + stale estimate, + hook/satellite slug-set mismatch, + false `full` flag.
> - Slicing: suggested dependency edges.
> - Fold gate: new goal-shaped vs defect-shaped condition (b), matching the original rationale; folds into gated REQs recorded with date/source; review-work.md added to files.
> - Context: out-of-scope consent-gate line restored; REQ-464/465 no longer double-counted as recurrences.
> - Verification: `git log -1` race removed.
> ```

---
*Captured: 2026-09-01T11:04:54Z*
