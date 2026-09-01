---
source_type: req_lesson
req_id: REQ-150
req_path: do-work/archive/UR-031/REQ-150-repair-live-modular-package-references.md
date: 2026-08-08
domain: general
module: skills/do-work/general
tags: [general, review, repair, live, modular]
---

# Lessons from REQ-150: Review fix: Repair live modular package references

## What the REQ was about

Make every URL and relative link published by the live modular packages resolve after installation and from the source archive. Close the entire shipped-reference class so a later package move cannot leave silent dangling guidance.

## Solution summary

Repaired all ten known live reference failures and established an executable package-publication contract. Repository-only history is linked through canonical GitHub pages; runtime references must exist in the exported modular layout; the installed changelog mirrors the release source.

## What worked

- Deriving source and installed topology from `suite/modules.tsv` caught links that resolve in the checkout but escape into unrelated client paths after installation.
- Separating exported raw runtime targets from tracked repository-history targets preserved browseable lessons and sidecars without shipping repository archives in each skill.

## What didn't work

- A handwritten Markdown scanner covered the current link corpus but missed standards-valid indented-code, HTML-comment, escaped-link, and escaped-destination cases; the release gate needs adversarial parser fixtures before it is fully robust.
- Letting the installed changelog drift across three releases made link policy and installed history diverge until byte identity became an explicit contract.

## Worth knowing

- Root `CHANGELOG.md` is now the release source and `skills/do-work/CHANGELOG.md` its byte-identical installed mirror; synchronize after every release edit.
- REQ-154 is held for user consent because REQ-150 is itself review-generated; current links and all production distribution tests are green.

## Back-reference

See `do-work/archive/UR-031/REQ-150-repair-live-modular-package-references.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `6dbb1cf`.
