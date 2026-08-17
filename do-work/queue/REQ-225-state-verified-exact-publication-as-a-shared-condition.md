---
id: REQ-225
title: State verified-exact-publication once as a condition in the shipped shell guide
status: pending
status_changed_at: 2026-08-17T21:09:46Z
domain: general
created_at: 2026-08-17T21:02:00Z
user_request: UR-042
addendum_to: REQ-220
review_generated: true
effort_estimate: normal
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: false
maintenance: true
write_set:
- skills/do-work/docs/prescribed-shell-primitives.md
---

# Discovered Task: State Verified-Exact-Publication Once as a Condition

## What

`skills/do-work/docs/prescribed-shell-primitives.md` is the canonical shipped guide for shell used across do-work actions. It states the "a rename onto an occupied destination nests instead of colliding" rule **only** inside its `## Portfolio summary publication` section, phrased as a property of one script. Restate it once as a condition that applies to any publication, and have the per-script sections point at that statement instead of each carrying (or omitting) their own copy.

## Context

Found while implementing REQ-220. The same defect has now been fixed in four separate places — `publish-portfolio-summary.sh` (REQ-199/205), the `ai-report` prescribed batch block (REQ-204), and `generate-report-image.sh` plus `install-last30days.sh` (REQ-220). Each fix was local, and each of the last three was found by a review sweep rather than by reading the guide, because the guide never says the rule in a form that would make a reader check their own publication against it.

`CLAUDE.md` § *State conditions, not lists* names this exact failure: when a rule applies "whenever X happens", key it on the condition, because a hand-maintained per-script list goes stale as the set grows. `_dev/primes/prime-shell-commands.md` § *Closed Enumerations Go Stale* records the same lesson from four independent defects.

## Requirements

- State the condition once — something to the effect of *any publication whose destination could be occupied verifies the path it actually wrote, rather than reading the rename's exit status as proof* — in a location that applies to every shipped publication helper, not inside one script's section.
- Have the existing per-script sections reference that statement rather than restating it. The `## Portfolio summary publication` section's script-specific consequences (snapshot candidates advance by numeric suffix, the canonical path fails closed) stay where they are — those are policy, not the shared primitive.
- Do not change any script behavior. This REQ is documentation only; all four scripts already implement the rule.
- Keep the shipped reference contract green (`_dev/tests/shipped-package-reference-contract.sh`) — this is a shipped file, so it may not cite `_dev/` paths.

## Open Questions

- [x] Should the shared shell guide state the publication-nesting rule once as a general condition, instead of describing it inside one script's section? → Confirmed: Yes, add to queue
  *[2026-08-17] User confirmed via `do-work clarify`. Consent given for this cascade-depth-two follow-up to run another autonomous cycle. Scope stays documentation-only: no script behavior changes, and the portfolio-summary section's script-specific policy stays where it is.*
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it — leave the rule stated per script.
  Why this is yours: nothing is broken. All four affected scripts already verify their publications, so this changes no behavior and fixes no defect — it is a judgment call about how the shipped guide is organized, and reorganizing a canonical shipped document is a taste decision rather than a repair. It is also worth knowing that REQ-220 was itself a review follow-up two generations deep, so the cascade-depth rule requires your consent before another autonomous cycle. The argument for doing it: this defect class has been found four times by review sweep and zero times by someone reading the guide, which is the concrete cost of the rule living inside one script's section. The argument against: the guide is deliberately organized by executable home, and a cross-cutting section cuts against that organizing principle.
