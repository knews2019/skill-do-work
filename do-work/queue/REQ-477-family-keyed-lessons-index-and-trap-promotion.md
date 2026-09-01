---
id: REQ-477
title: '[impact-rule-change] Family-keyed lessons, intelligent index, and mandatory Trap promotion'
status: pending
created_at: 2026-09-01T10:47:44Z
user_request: UR-088
domain: general
prime_files: [_dev/primes/prime-action-files.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: false
impact: impact-rule-change
effort_estimate: effort-substantive
related: [REQ-478, REQ-479]
batch: lessons-transfer-routing
write_set: [skills/do-work/actions/work.md, skills/do-work/crew-members/general.md, do-work/lessons-index.md]
---

# Family-Keyed Lessons, Intelligent Index, and Mandatory Trap Promotion

## What

Make the Lessons-Capture Phase produce transferable context: every appended lesson bullet carries a failure-family slug, an intelligent index over the lesson satellites is created/refreshed in the same edit, and a second same-family occurrence makes Trap promotion mandatory instead of a judgment call.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Detailed Requirements

- **Intelligent lessons index.** A single well-known plain-markdown index with one line per lessons satellite — path, a when-it-applies hook naming the failure families inside (e.g. "rollback/deletion/final-boundary work in do-work-cli internals"), and an approximate token estimate (mechanical, e.g. bytes/4 — the exact formula is the builder's call but must be reproducible by a floor agent with `wc`). Suggested location `do-work/lessons-index.md` (precedent: `do-work/prose-backlog.md`); the builder may argue a better home.
- **Maintained by the writer, in the same edit.** Whenever `skills/do-work/actions/work.md`'s Lessons-Capture Phase appends a lesson (the Step 8 satellite write), the same edit creates or refreshes that satellite's index line, token estimate included.
- **Family-keyed bullets.** Every appended lesson bullet carries a short kebab-case failure-family slug — the index hooks and recurrence checks key on it; same discriminator move as `sweep_key`. The slug must be literal-matchable (greppable).
- **Mandatory twice-seen Trap promotion.** The Step 8 lesson write scans the satellite first: on the second same-family occurrence (by slug, or same-family judgment for pre-slug entries), promotion stops being optional — the writer adds or amends one generalized Trap line in the owning prime's `## Traps` in the same edit, citing the family slug, and notes the promotion in the hand-back. The "judgment call" sentence at `skills/do-work/actions/work.md:604` becomes condition-keyed.
- **Seed worked examples only.** Do not backfill old lesson entries wholesale; seed slugs and index hooks for the three known recurring families: final-boundary identity/rollback (REQ-413/436/447/463/416), opaque-evidence/generic-fallback projections (REQ-414/430/446), smoke-vs-characterization gates (REQ-414/415).
- Update `skills/do-work/crew-members/general.md` § Lessons Discipline restatements in the same change.

## Constraints

- Plain files only — capture and builders on the floor agent (read/write files + shell) must be able to use the index with read/grep alone; no new tooling dependency.
- Both routing mechanisms deliberately coexist: capture-time stamping (REQ-478) covers stamped REQs; Trap promotion covers everything else. Do not trade one away while implementing the other.
- The prime stays a routing index: promotion writes one generalized Trap line, never the per-REQ narrative (that stays in the satellite).

## Dependencies

None. REQ-478 and REQ-479 build on this REQ's index format and slugs.

## Builder Guidance

Certainty level: Firm — the design was decided interactively with the maintainer (2026-09-01 session). Latitude: index home, line format, and estimate formula are the builder's call within the stated constraints.

- [~] Exact index location → builder decides; `do-work/lessons-index.md` recommended (precedent `do-work/prose-backlog.md`).

## Red-Green Proof

**RED prompt/case:** Follow today's Lessons-Capture Phase instructions to append a second lesson for an already-recorded failure family (e.g. another final-boundary identity finding): no slug is required, no index exists to refresh, and no instruction compels a Trap line — `work.md:604` leaves promotion as an unowned judgment call.
**Why RED now:** Three failure families each recurred 3–6 times across the 2026-08-31 run while their lessons sat in `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md`; REQ-415 repeated a family whose lesson REQ-414 had already recorded.
**GREEN when:** The shipped Lessons-Capture instructions require slugged bullets and the same-edit index refresh with token estimate; a twice-seen family mandates a generalized Trap line in the owning prime in the same edit; the three known families are seeded as worked examples in the satellite, index, and prime Traps.
**Validation:** User confirmed (approved plan, 2026-09-01 session).

## Full Context

See `do-work/user-requests/UR-088/input.md` for complete verbatim input.

---
*Source: UR-088 (Lessons routing with token-budgeted mandatory reads and a fold-gate fix)*
