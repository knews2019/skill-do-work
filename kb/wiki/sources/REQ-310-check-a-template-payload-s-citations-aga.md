---
title: "Lessons from REQ-310: Check a template payload's citations against where the payload lands"
type: source-summary
topic_cluster: metadata-and-timestamps
sources: [raw/processed/2026-09-01/REQ-310-check-a-template-payload-s-citations-aga.md]
related:
  - page: concept-timestamp-and-metadata-governance
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-310: Check a template payload's citations against where the payload lands

Part of the [[concept-timestamp-and-metadata-governance]] cluster.

## What the REQ was about

REQ-269 split the fence exemption by its reason: a fenced payload is exempt because its text lands in *another* file, so it is not a citation from the file it is written in. That reasoning is sound and the exemption survives on it — but it only establishes that the payload is not checkable **from here**. It says nothing about whether the payload resolves **there**, and nothing checks that today.

A live example, found while confirming REQ-269's scope-drift was intended:

## Solution summary

The two toolbox payloads that point into core actions now emit consumer-root installed citations,
and the existing review-generated producer contract rejects their retired source-relative form.
This is a leaf repair to the action/template reference contract; no new destination syntax was added.

## What worked

**What worked:** Counting the affected payloads before designing selected the direct repair, while
running both citation gates exposed the existing assertion that encoded the wrong source-relative
form.

**What didn't:** The first replacement used `actions/...`, matching core-generated REQs but failing
the staged toolbox reference gate. A consumer-root installed path is stable across queue, working,
and archive moves and satisfies both readers.

**Worth knowing:** Template citations have two locations: the producer and the emitted artifact.
Source-package correctness says nothing about the destination; tests must judge the emitted form.

## Back-reference

See `do-work/archive/UR-055/REQ-310-check-template-payload-citations-at-their-destination.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `61cfc28`.
