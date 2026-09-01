---
title: "Lessons from REQ-269: Draw the cross-package citation class by what a citation is, not by the punctuation around it"
type: source-summary
topic_cluster: suite-and-package-architecture
sources: [raw/processed/2026-09-01/REQ-269-draw-the-cross-package-citation-class-by.md]
related:
  - page: concept-modular-suite-architecture
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-269: Draw the cross-package citation class by what a citation is, not by the punctuation around it

Part of the [[concept-modular-suite-architecture]] cluster.

## What the REQ was about

Eight consecutive REQs in this area have bounded the cross-package citation rule by a **marker** — a leading `../`, a backtick, a fence character — instead of by the **thing** the rule governs: a path in shipped text that a reader is expected to follow from the citing file's own directory. Each fix closes one spelling and leaves the class open, so the next review finds the same defect wearing different punctuation. This REQ ends that by making the checker's condition *be* the rule's condition.

Concretely, three markers are currently doing the bounding, and each has live escapees:

## Solution summary

Redrew the cross-package citation class in the contract checker so its condition is the rule's condition. A token is a citation when its first path segment — after an optional `../` lead — names a sibling package directory. Backticks, the `../` lead, and the fence stopped being part of that test. The scan surface widened from backticked spans only to every path-shaped token in prose, inline code, HTML comment interiors, and in-fence annotations; the fence exemption was split by REQ-249's stated rationale so a fenced payload stays exempt while an annotation beside it does not. Two conditions take a token back out of the class, both about meaning: a path rooted somewhere else (`<skill-root>/…`, `.claude/skills/…`) is not a path from here, and a bare `do-work/…` token is the consumer's queue root — a collision only the core package's name has.

## What worked

**What worked:** Reproducing the captured RED on the untouched tree *before* writing any code. The REQ handed over a mutation and a claim that it passes today; running it first turned that claim into an observation and made every later "it fails now" mean something. The second thing that worked was measuring before designing — a naive widening of the predicate produces 635 non-resolving hits, and knowing that number is what forced the narrowed consumer-state rule instead of a plausible-looking widening that would have been reverted within a day.

**What didn't:** The exploration probe was written twice before it was right, and both failures were the same mistake in miniature as the REQ itself. First `str.strip(".")` ate the leading `..` of every `../` token, so instance 4 looked clean when it was not. Then a mid-token regex `search` matched the tail of `<skill-root>/../do-work/scripts/…` and reported four rooted paths as broken citations. Both are punctuation standing in for meaning — the exact defect being fixed, committed by the tool built to find it. The second one survived into the implementation and was only caught because the real checker reported the four sites; the fix for it (anchor the match at the token start) turned out to be the correct semantic rule and is now a fixture.

**Worth knowing:** The three `crew-members/prompt-injection.md` copies are byte-identical across the packages but nothing enforces that, and they must now differ — from `do-work-toolbox` the correct citation is the same-package `../actions/…`, not a cross-package hop. Anyone tempted to re-mirror those three files will silently reintroduce instance 2. Related: `mask_block_code` now takes an optional `masked_line_ranges` out-parameter so the fence split reuses one fence state machine rather than growing a second walker — if a third caller needs the ranges, that is the seam to extend, not to copy.

## Back-reference

See `do-work/archive/UR-055/REQ-269-draw-the-citation-class-by-what-a-citation-is.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `f71dfee`.
