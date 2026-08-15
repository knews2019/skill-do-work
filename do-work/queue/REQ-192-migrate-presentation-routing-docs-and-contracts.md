---
id: REQ-192
title: Migrate completed-work presentation routing documentation and contracts
status: pending
created_at: 2026-08-15T09:10:53Z
user_request: UR-042
domain: testing
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec: refactor
depends_on: [REQ-189, REQ-190, REQ-191]
maintenance: true
effort_estimate: normal
related: [REQ-189, REQ-190, REQ-191]
batch: completed-work-presentation-consolidation
---

# Migrate Completed-Work Presentation Routing Documentation and Contracts

## What

Update toolbox routing, discovery, completion-flow recommendations, cross-references, tests, and release notes so the suite presents one unambiguous completed-work choice: detailed report through `ai-report`, portfolio through `present-work`, and animated walkthrough through `present-video`.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

The live router, help, tutorials, guides, next-step recommendations, caller lists, and historical successor cross-references currently point at overlapping command ownership. Without a coordinated migration, users can still enter retired detail and video workflows even after the underlying actions are separated.

## Context

REQ-189 establishes canonical detailed reporting, REQ-190 narrows `present-work`, and REQ-191 adds `present-video`. This dependent REQ makes those boundaries discoverable and adds regression coverage that rejects their regrowth.

## Detailed Requirements

### Routing and discovery

- Update the toolbox router and argument hint for the three command families.
- Route `showcase`, `visual report`, and `proof of work` to `ai-report`.
- Route `present-work`, `portfolio`, and `work portfolio` to portfolio-only `present-work`.
- Route `present-video`, `remotion`, and `video walkthrough` to `present-video`.
- Preserve explicit `ai-report` routing and do not rename it.
- Ensure bare `present-work` is usage-only and item-specific `present-work` is migration-guidance-only; do not silently delegate those invocations in the router.

### Documentation and cross-references

- Update toolbox help, tutorials, action guides, next-step recommendations, caller lists, and current cross-references for the new ownership.
- Replace completion-flow recommendations that currently suggest `present-work UR-NNN` with `ai-report UR-NNN`.
- Keep cross-project portfolio examples on `present-work all|portfolio`.
- Add explicit animated walkthrough examples using `present-video <ID>` only.
- Present one unambiguous choice everywhere current behavior is documented:
  - detailed report → `ai-report`;
  - cross-project portfolio → `present-work`;
  - animated walkthrough → `present-video`.
- Update prompt-injection and anti-slop caller guidance so every archive-reading or human-artifact action is covered by the condition and current examples are not stale.
- Correct current successor statements and cross-references without deleting or rewriting accurate historical records of previously generated artifacts.
- Preserve all existing generated artifacts; do not migrate or delete old briefs, `.single.html` files, reports, or video directories.

### Contract and acceptance coverage

- Lock in that `ai-report` is the only action capable of producing stakeholder-facing detailed HTML.
- Lock in that UI work retains screenshots, SVG callout annotations, responsive layout, authentic before/after evidence, and full-page light/dark render verification.
- Lock in that backend and refactor work succeeds in non-visual evidence mode without fabricated screenshots and states that UI captures were not expected.
- Prove that `present-work all` and `present-work portfolio` produce only `do-work/deliverables/portfolio-summary.md`.
- Prove that bare `present-work` writes nothing and prints compact usage.
- Prove that item-specific `present-work` writes nothing and prints the exact `ai-report <ID>` and `present-video <ID>` replacement commands.
- Prove that `present-video` creates a valid Remotion source tree with `registerRoot` and no MP4.
- Prove that `ai-report` and `present-work` never create video artifacts and never invoke `present-video` automatically.
- Prove that all archive-reading actions load prompt-injection guidance before archived user content.
- Prove that presentation actions accept the terminal-success states `completed` and `completed-with-issues` and reject cancelled or unfinished work.
- Make contract tests reject stale Detail Mode, Interactive Explainer, detail Client Brief, sibling-link/detail-depth, automatic video, and unsafe Remotion preview workflows.
- Update shipped-package inventories, caller lists, action/guide counts, routing fixtures, and other contract-owned enumerations required by the new files and routes.

### Release record

- Add a concise changelog entry explaining the command migration and the three resulting ownership boundaries.
- Apply the repository’s required shared version bump and installed changelog mirror synchronization at the integrating commit.

## Constraints

- Detailed report means `ai-report`; cross-project portfolio means `present-work`; animated walkthrough means `present-video`.
- Do not add `--with-video` or automatic video behavior to `ai-report`.
- Do not change UR/REQ schemas, archive formats, `review-work`, or implementation behavior.
- Do not add publishing, hosting, search, MP4 rendering, or automatic video generation.
- Preserve all existing generated artifacts.
- Keep tests focused on the real ownership and archive-safety failures; do not add decorative snapshots of prose.

## Dependencies

Depends on REQ-189, REQ-190, and REQ-191 so the router, docs, inventory, and tests describe actions that already have their final ownership.

## Builder Guidance

Firm intent. Treat current-behavior references as a condition-based sweep, not a hand-maintained filename checklist; keep accurate history intact while removing stale live guidance. Start with failing contract cases for the current overlapping ownership and unsafe preview rules.

## Open Questions

None.

## Red-Green Proof
**RED prompt/case:** Search current routers, argument hints, help, tutorials, completion-flow examples, guides, caller lists, shipped inventories, and contract fixtures for the three command families and the retired `present-work` detail/video workflows, then run the baseline contract suites.
**Why RED now:** `showcase` still routes to `present-work`; completion guidance still recommends `present-work UR-NNN`; no `present-video` action is discoverable; live docs and contracts still endorse Client Brief, Interactive Explainer, detail-depth, and embedded Remotion behavior while baseline checks pass.
**GREEN when:** Every live surface presents the same three-way command ownership; exact routing aliases dispatch correctly; focused contract tests fail if detailed HTML or video regrows under `present-work`, if non-visual work redirects away from `ai-report`, if unsafe Remotion preview commands or MP4 rendering return, or if terminal-success/prompt-injection rules drift; inventories and baseline suites pass; the changelog records the migration without modifying old generated artifacts.
**Validation:** Inferred during capture from the supplied acceptance tests.

## Full Context

See `do-work/user-requests/UR-042/input.md` for the complete verbatim request and batch constraints.

---
*Source: attached “do-work capture-request: Consolidate completed-work presentation around ai-repo…” specification.*
