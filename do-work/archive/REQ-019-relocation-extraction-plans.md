---
id: REQ-019
title: "Extraction plans (plan-only) for prompts library, interview subsystem, bkb + dream"
status: completed
created_at: 2026-07-15T17:33:04Z
user_request: UR-003
claimed_at: 2026-07-15T18:38:00Z
completed_at: 2026-07-15T18:47:00Z
route: B
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: []
related: []
batch: harness-bloat-cleanup
maintenance: false
commit: b8acdbb
---

# RELOCATE extraction plans (no extraction in this pass)

## What
Write `decisions/audits/2026-07-15-relocation-extraction-plans.md` covering, for
each RELOCATE item (audit phase-2 bucket):

- **prompts/ library** (12 off-mission files; recommendation: whole library moves,
  `actions/prompts.md` runner stays)
- **interview subsystem** (interview.md + interview-reference.md + interviews/ +
  crew-members/interviewer.md + docs/interview-guide.md)
- **bkb** (bkb.md + bkb-reference.md + docs/bkb-guide.md) **+ dream** (dream.md +
  docs/dream-guide.md) as one package

Each plan names: target repo/skill name, exact file manifest with word counts, the
pointer that remains in do-work (router row fate, install-target or doc pointer),
coupling seams to cut (kb-lessons-handoff degradation, bkb↔dream, interview→bkb
ingest, SKILL.md routing rows, next-steps blocks, tutorial/README mentions,
contract-test lines), and a migration note for existing users (what `do-work
update` does to the removed files, how to install the sibling).

## Why
User-approved RELOCATE bucket; Message 1 explicitly scopes phase 3 to plans only.

## Acceptance criteria
- [x] One plan per package with complete file manifest (verified against the repo).
- [x] Every inbound reference to each package enumerated by grep, with its fate.
- [x] Migration note covers both git-clone and tarball consumers.

## Open Questions
- [ ] Do the 4 dev-adjacent prompts (ADR log, dark-code kit ×3) move with the
      library or stay? → **D-01**: Builder chose: plan moves the whole library.
      Reasoning: coherence — the runner supports project-local dirs, so dev-adjacent
      prompts remain installable per-project. Value: one clean seam. Risk: low,
      reversible at extraction time; the plan marks these 4 as "may stay" for the
      user's final call.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Grep-enumerated inbound refs per package across SKILL.md, next-steps.md, actions/, crew-members/, docs/, README, _dev/, hooks/ before writing any seam list.
- [x] **[APPLY]:** Wrote decisions/audits/2026-07-15-relocation-extraction-plans.md — three plans (prompt library ~19.8k words, interview ~12.6k, bkb+dream ~14.9k), each with manifest, seams, migration note; recommended sequence A→B→C.
- [x] **[UNIFY]:** No shipped file modified (plan doc is export-ignored under decisions/); kb-lessons-handoff boundary explicitly kept in do-work with kb_status/kb_entry schema fields; D-01 (dev-adjacent prompts move, marked "may stay") recorded for user's final call.
