---
source_type: req_lesson
req_id: REQ-419
req_path: do-work/archive/REQ-419-add-flat-just-recipes-action-delegation.md
date: 2026-09-01
domain: general
module: _dev/primes
tags: [general, flat, just, recipes, collision]
---

# Lessons from REQ-419: Add flat Just recipes, collision validation, action delegation, and compatibility aliases

## What the REQ was about

Expose every public command through flat Just recipes and make the existing skill aliases delegate to the canonical CLI.

## Solution summary

Published a 40-definition flat Just interface backed by one managed template, made installer collision/completeness validation derive from that template, and delegated the named deterministic action phases to `do-work-cli` with fail-closed behavior. Review remediation replaced raw `{{args}}` shell-source interpolation in all 33 non-board variadic recipes per surface with Just positional argv and quoted `"$@"`, made recovery recipes include the canonical manifest/time flags, moved hostile tests onto the actual shipped recipe, and aligned action/reference/prime/board restatements with singular CLI ownership.

## What worked

Reproducing the failure through the actual managed template with command-substitution sentinels exposed the real shell boundary; Just's per-recipe `[positional-arguments]` plus `"$@"` preserved every original argument byte.

## What didn't work

POSIX-quoting the outer generated command was insufficient while inner recipes interpolated `{{args}}` as shell source. The first test used an invented safer recipe and therefore proved the wrong seam.

## Worth knowing

Publication recipes must preserve two boundaries together: shell-literal outer arguments and positional inner Just arguments. Named leading parameters need capture-and-shift before canonical CLI flags are rebuilt.

## Back-reference

See `do-work/archive/REQ-419-add-flat-just-recipes-action-delegation.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `cf7b8977`.
