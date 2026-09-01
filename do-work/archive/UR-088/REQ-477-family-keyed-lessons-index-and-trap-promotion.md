---
id: REQ-477
title: '[impact-rule-change] Family-keyed lessons, intelligent index, and mandatory Trap promotion'
status: completed
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
write_set: [skills/do-work/actions/work.md, skills/do-work/actions/review-work.md, skills/do-work/crew-members/general.md, skills/do-work-toolbox/crew-members/general.md, skills/do-work/tools/do-work-cli/lessons-do-work-cli.md, skills/do-work/tools/do-work-cli/prime-do-work-cli.md, do-work/lessons-index.md]
claimed_at: 2026-09-01T15:40:02Z
route: C
estimate:
  p50_active_minutes: 40
  confidence: medium
  calculated_at: 2026-09-01T15:41:12Z
  basis:
    - Route C
    - 7-file write set
    - 1 new file
    - 2 subsystems involved
    - 6 acceptance criteria
    - cross-route regression gates
kb_status: promoted
kb_entry: REQ-477-family-keyed-lessons-intelligent-index-a.md
completed_at: 2026-09-01T15:56:55Z
commit: 74b1fd41
---

# Family-Keyed Lessons, Intelligent Index, and Mandatory Trap Promotion

## What

Make the Lessons-Capture Phase produce transferable context: every appended lesson bullet carries a failure-family slug, an intelligent index over the lesson satellites is created/refreshed in the same edit, and a second same-family occurrence makes Trap promotion mandatory instead of a judgment call.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Define one grep-friendly bullet tag and one mechanically reproducible index-row format, update both lesson writers to maintain them and promote a twice-seen family, then seed only the three named do-work-cli families by tagging their existing examples and generalizing existing prime traps.
- [x] **[APPLY]:** Added the shared family marker/index/promotion contract to both lesson writers and crew guidance; seeded the three named do-work-cli families without wholesale backfill.
- [x] **[UNIFY]:** Reviewed the seven-file REQ diff, recomputed every index row from `wc -c`, checked exact slug sets and `git diff --check`, and ran the canonical maintainer gate on a clean HEAD plus this patch.

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

## Addendum (2026-09-01)

User added (v4 revision, validate-feedback Finding 1 — Accept):

> ```
> a when-it-applies hook that names the failure-family slugs inside (e.g. "rollback/deletion/final-boundary work in do-work-cli internals — final-boundary-identity, opaque-evidence-projection"); a mechanical size estimate; and a slug-coverage flag (`slugged: full` when every bullet carries a family slug, `partial` otherwise). Maintained by work.md's Lessons-Capture Phase: whenever a lesson is appended (Step 8's satellite write), the same edit creates or refreshes that satellite's index line — hook slugs, estimate, and coverage flag included.
> ```

- The index hook names the exact family-slug set present in the satellite, not just prose — capture's targeted stamping greps against it.
- Each index line adds a slug-coverage flag: `slugged: full` when every bullet carries a family slug, `partial` otherwise.
- The Lessons-Capture write maintains hook slugs, size estimate, and coverage flag in the same edit as the lesson append.
- Seeding the three known families makes `lessons-do-work-cli.md`'s coverage flag honest, whatever it ends up being.
- Provenance: validate-feedback 2026-09-01, Finding 1. Surface-cost: Earned — a targeted `path#slug` read against a partially-slugged satellite would silently skip pre-slug same-family bullets (replay: `#final-boundary-identity` misses the REQ-436/447/463 bullets); the flag is one token per line, policed by the audit (REQ-479 addendum).

## Triage

**Route: C** — this changes the lesson-write contract shared by orchestration, standalone review, capture routing, and prime guidance.

## Plan

1. Define the canonical `[family: <slug>]` bullet marker and a one-row-per-satellite Markdown index with `ceil(bytes / 4)` estimates, exact slug sets, and coverage.
2. Update both satellite writers to choose a slug, scan before append, refresh the index in the same edit, and require a generalized prime Trap on the second same-family occurrence.
3. Seed only the three accepted do-work-cli families in existing bullets and prime traps, then verify every indexed path and estimate mechanically.

**Plan validation:** All requirements map to one of the three tasks. `review-work.md` is added to scope because it is the second live satellite writer; leaving it unchanged would violate “every appended lesson bullet.”

## Exploration

- Orchestrated writes live in `actions/work.md` Step 8; standalone review has a separate writer in `actions/review-work.md` Step 9.5.
- Six tracked lesson satellites are genuine source satellites. Build-directory copies and ordinary files whose names merely contain “lessons” are excluded.
- The do-work-cli prime already has three traps that can be generalized and slugged without growing the routing index.
- Existing satellites are mostly un-slugged, so their initial coverage is honestly `partial`; no wholesale backfill is needed.

## Scope

**Files I will touch:**
- `skills/do-work/actions/work.md` (modify) — orchestrated lesson-write contract
- `skills/do-work/actions/review-work.md` (modify) — standalone lesson-write contract
- `skills/do-work/crew-members/general.md` (modify) — builder-facing Lessons Discipline
- `skills/do-work-toolbox/crew-members/general.md` (modify) — semantic mirror used by toolbox prime workflows
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` (modify) — three seeded family examples
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (modify) — three generalized family traps
- `do-work/lessons-index.md` (new) — one mechanically maintained row per tracked source satellite

**Acceptance criteria:** Every future bullet is slugged; both writers refresh path, hook slugs, estimate, and coverage in the same edit; second occurrence mandates Trap promotion; all source satellites are indexed; only the three named families are seeded.

## Decisions

- **D-01 — DECIDE & STATE:** Include `review-work.md` even though capture's initial write set omitted it. It is a live standalone satellite writer, so excluding it would leave an alternate path producing un-slugged, unindexed lessons.
- **D-02 — DECIDE & STATE:** Keep the toolbox copy of `general.md` semantically aligned; its package-relative prime citation remains the only intentional difference.

## Implementation Summary

**Files changed:**
- `skills/do-work/actions/work.md` (modified) — slugged orchestrated lesson writes, same-edit index refresh, mandatory recurrence promotion
- `skills/do-work/actions/review-work.md` (modified) — the same contract for standalone review writes
- `skills/do-work/crew-members/general.md` (modified) — canonical Lessons Discipline
- `skills/do-work-toolbox/crew-members/general.md` (modified) — toolbox semantic mirror
- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` (modified) — seeded the three accepted family markers only
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (modified) — generalized traps for the three recurring families
- `do-work/lessons-index.md` (new) — six tracked source satellites with hooks, exact family sets, mechanical estimates, and coverage

**What was done:** Lesson writers now emit a literal family discriminator, refresh one reproducible index row in the same edit, and promote a generalized trap by the second occurrence. The initial index is complete for tracked source satellites and honest about partial legacy coverage.

## Qualification

Passed — seven declared files are substantive and wired through both lesson-writing paths; all detailed and addendum requirements trace to the diff; P-A-U is complete.

## Testing

- Custom index contract: PASS — exact six-path inventory, live paths, `ceil(bytes/4)` estimates, seeded lesson/trap slugs.
- `git diff --check` on the REQ scope: PASS.
- Clean isolated `bash _dev/tests/contract-regressions.sh`: PASS.
- Clean isolated `bash _dev/tests/shipped-package-reference-contract.sh`: PASS.
- Clean isolated `bash _dev/tests/maintainer-verify.sh`: PASS (browser lane explicitly skipped because no browser was available, per the gate's own contract).
- The shared dirty tree's first contract run failed only in paused REQ-420 shim surfaces; the same REQ-477 patch passed from clean HEAD, establishing attribution.

## Review

**Verdict: Approve.** The implementation covers both live writers, preserves legacy lesson prose, avoids a backfill, keeps every index field mechanically reproducible, and turns the three named recurrences into generalized traps. No Important findings.

**Acceptance:** Pass

## Lessons Learned

**What worked:** Enumerating live writers before editing prevented the standalone review path from silently bypassing the new contract; testing the patch on a clean detached worktree separated REQ-477 evidence from paused REQ-420 changes.

**What didn't:** Running the repository contract suite directly in the shared dirty tree produced many unrelated REQ-420 failures and could not qualify this change.

**Worth knowing:** An output-format rule is incomplete until every writer is swept, even when the REQ names only the primary writer.

## Orientation

[MAP CHANGED] Lesson satellites now have a family-keyed routing index and a mechanical recurrence-to-Trap path, shared by orchestrated work and standalone review. This is the base that capture-time and claim-time lesson routing can consume in REQ-478 and REQ-479.
