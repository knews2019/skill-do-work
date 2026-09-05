---
id: UR-127
title: 'Reduce repeated setup and unaffected reruns in the fast gate'
created_at: 2026-09-05T19:43:25Z
requests: [REQ-591]
word_count: 3
---

## Summary

Capture the recommended next request from the CPU-spike and test-speed discussion: make routine verification faster while retaining useful failure coverage. The user requested capture only.

## Extracted Requests

| Request | Intent |
| --- | --- |
| REQ-591 (Reduce repeated setup and unaffected reruns in the fast gate) | Reduce repeated setup first, then avoid verification reruns whose complete inputs are unchanged; prove the speed improvement and preserve correctness checks. |

## Conversation Context

The immediately preceding user question was:

> ````text
> this also means that the tests are very heavy.
> 
> Can we make them faster while maintaining a good enough quality?
> ````

The assistant recommended the named request, beginning with immutable fixture/build reuse and then selective verification with conservative invalidation. The broader discussion also considered cheaper decision-matrix tests and shared machine concurrency limits. Those remain context; this capture implements the specifically recommended fast-gate request, not a cross-project scheduler or a general test-deletion program.

The investigation observed overlapping Go and Jest runs across three projects on the same eight-core machine. This was distinct from the earlier orphan synthetic-load incident. Existing self-bounding worker rules still apply, and timings from an overloaded window are not a normal performance baseline.

## Full Verbatim Input

> ```
> capture the req
> ```
