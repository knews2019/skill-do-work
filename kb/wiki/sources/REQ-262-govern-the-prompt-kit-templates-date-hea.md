---
title: "Lessons from REQ-262: Govern the prompt-kit templates' date headers"
type: source-summary
topic_cluster: knowledge-and-memory
sources: [raw/processed/2026-09-01/REQ-262-govern-the-prompt-kit-templates-date-hea.md]
related:
  - page: concept-knowledge-and-memory-systems
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-262: Govern the prompt-kit templates' date headers

Part of the [[concept-knowledge-and-memory-systems]] cluster.

## What the REQ was about

Three prompt-kit templates in do-work-knowledge carry `Date: [today]` headers that no paragraph of the Timestamp rule governs and that sit outside the citation checker's reach (they are template content, not action prose). Decide whether they join the date-only paragraph's consumer list (UTC, cited) or are declared template-content-out-of-scope like the fenced-block exemption.

## Solution summary

Added one condition-keyed sentence to the date-only-stamp paragraph in `skills/do-work/actions/work-reference.md`, declaring that a date placeholder inside a template artifact's fenced block is out of scope for the Timestamp rule. The sentence keys on what the site *is* (a fill-in token addressed to the model that emits a document for the user) rather than on which files currently contain one, and names neither the paths nor how many there are. The three prompt-kit templates were deliberately not edited: the decision is that their `Date: [today]` lines are correct as written.

## What worked

**What worked:** Reading the target paragraph's *existing* carve-outs before writing a new one. The `## HH:MM UTC` exclusion three sentences away was already the exact shape this decision needed — an out-of-scope class stated at the rule's home so a sweep walks past it — so the work was matching an established pattern rather than inventing a policy. Also: running the Restatement Sweep against the REQ's own diff, not only against other files. It is what caught M1.

**What didn't:** The first draft stated "holds three today". It felt like useful precision and it was the opposite — `_dev/primes/prime-kanban-board.md:51` records REQ-261 deciding that a count which does not bear on the argument is clutter, and the changelog shows a count-keyed tripwire being removed from *this same paragraph* for the same reason. Writing a count into a paragraph a prior REQ had just de-counted would have re-opened settled drift. The sweep caught it; reading the prime first would have prevented it.

**Worth knowing:** The three prompt-kit `Date: [today]` lines were never reachable by the citation checker in the first place — `_dev/tests/shipped-package-reference-contract.sh` masks fenced, list-fenced, and indented code before it reads. So the REQ's premise ("outside the citation checker's reach") was correct but understated: they are outside it *structurally*, by a property of every fenced block, not by an accident this decision needed to preserve. That is why the exclusion could be stated as a condition with nothing to keep in sync.

## Back-reference

See `do-work/archive/UR-055/REQ-262-govern-the-prompt-kit-template-date-headers.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `24587e5`.
