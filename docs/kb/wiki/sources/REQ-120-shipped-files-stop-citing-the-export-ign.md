---
title: "Lessons from REQ-120: Shipped files stop citing the export-ignored maintainer doc"
type: source-summary
topic_cluster: suite-and-package-architecture
sources: [raw/processed/2026-09-01/REQ-120-shipped-files-stop-citing-the-export-ign.md]
related:
  - page: REQ-116-normalize-route-at-the-board-s-read-site
    rel: complements
  - page: REQ-118-the-normalize-flag-must-stop-calling-voc
    rel: complements
  - page: REQ-119-an-off-vocabulary-route-warns-on-the-boa
    rel: complements
created: 2026-09-01
updated: 2026-09-02
confidence: medium
---

# Lessons from REQ-120: Shipped files stop citing the export-ignored maintainer doc

Part of the [[concept-modular-suite-architecture]] cluster.

## What the REQ was about

Four locations in shipped `tools/queue-kanban/` files cite this repo's own `CLAUDE.md`, which is `export-ignore`d — so the citation dangles in every consumer install and every clone that installs from the tarball. `_dev/tests/contract-regressions.sh`'s maintainer-document probe fails on all four. Restate each rule inline or point at a shipped home.

## Solution summary

All four citations of the export-ignored maintainer doc are gone from shipped paths, with each rule restated rather than dropped. The floor-exception comment now points at `actions/board.md` and `actions/work-reference.md` → Timestamp rule, both of which ship. The write-surface comment states the amend-in-the-same-commit obligation without naming the file that carries it. `hasSchemaFieldContract`'s comment keeps the goes-stale reasoning as its own clause. The prime's REQ-116 lesson keeps its full content and refers to "the maintainer doc's lock-step rule" instead of the filename — and gained the actionable half, that a field joins the list when the board starts parsing it.

## What worked

Running the probe and reading **every** FAIL line rather than the tail. This finding was in the output of a suite I had already run and reported on earlier in the session; I had counted the update-script sub-suite's own summary ("7 failure(s)") and the trailing line, and missed a separate check failing above them. The reviewer found it in the same output I had.

## What didn't work

Two things. First, I wrote two of these citations while holding the rule that forbids them in context — the rule's own file was loaded, and the violation still went in, because "cite where this rule lives" is a strong writing instinct and the rule's whole point is that this particular file cannot be cited. Second, the earlier session's REQ-112 introduced the other two, meaning the probe had been red on `main` and nobody noticed — so a check that fails for an unrelated reason (here, 7 root-runner probes) provides cover for real failures in the same output.

## Worth knowing

The probe deliberately flags **any** mention rather than matching citation idioms, because idiom-matching caught 0 of 8 real occurrences before it was inverted. The consequence for a writer: there is no phrasing that legitimately names this file from a shipped path — restate the rule or point at a shipped home. `maintainer_doc_mention_allowlist` is only for mentions of a *consumer project's* CLAUDE.md, so reaching for it to quiet a hit is nearly always the wrong fix.

## Back-reference

See `do-work/archive/UR-025/REQ-120-shipped-files-stop-citing-maintainer-doc.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `27080ba`.
