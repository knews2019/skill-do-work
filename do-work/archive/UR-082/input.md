---
id: UR-082
title: 'Fix Accepted Review Findings'
created_at: 2026-08-29T20:26:10Z
requests: [REQ-421, REQ-422, REQ-423, REQ-424]
word_count: 423
---

# Fix Accepted Review Findings

## Summary

Implement four accepted review findings covering consumer-safe board tests, live timeline freshness, fetch interruption, and fallback branch selection. Ship them together as patch release 0.244.9 with tests, lessons, changelogs, and the repository's two-commit implementation/bookkeeping trail.

## Extracted Requests

| Request | Summary |
|---|---|
| REQ-421 | Keep board invariants universal while limiting suite-sized numeric floors to suite checkouts. |
| REQ-422 | Rebuild time-derived Timeline data on unchanged-tree live cache hits. |
| REQ-423 | Make HUP, INT, and TERM terminate the archive fetcher with conventional statuses. |
| REQ-424 | Make Git fallback clone the branch named by a canonical tarball URL. |

## Batch Constraints

- Deliver all four accepted findings in this checkout as release `0.244.9`, titled “Keep Consumer Tests, Timelines, and Fetches Current.”
- Make no public CLI, JSON schema, or API changes.
- Update the linked prime/lesson pairs, reusing the existing REQ-303 consumer-corpus lesson.
- Keep the two fetcher copies and the two changelog copies byte-identical.
- Commit verified implementation once, then record that hash in all four completed REQs in a separate bookkeeping commit.

## Full Verbatim Input

> ````text
> PLEASE IMPLEMENT THIS PLAN:
> # Fix Accepted Review Findings
>
> ## Summary
>
> - Record the four accepted findings as one UR with four provenance-preserving REQs, then implement them in this checkout.
> - Deliver a patch release, `0.244.9`, titled “Keep Consumer Tests, Timelines, and Fetches Current.”
> - No public CLI, JSON schema, or API changes; behavior and tests only.
>
> ## Implementation Changes
>
> - **Consumer-safe board tests:** Keep citation, fence, and shipped-payload invariants running everywhere, but wrap their three suite-sized numeric floors with `suiteCheckoutSkipReason`. Add a consumer-shaped subprocess regression running those exact tests against a small queue.
> - **Fresh live timeline:** Extract timeline payload construction into one helper taking tickets, duration history, and an explicit instant. Add an injectable live-server clock; on unchanged-tree cache hits, update `GeneratedAt` and rebuild the complete Timeline payload from the cached parsed board without reparsing files or Markdown.
> - **Correct fetch interruption:** In both byte-identical fetcher copies, reserve `EXIT` for cleanup and make HUP/INT/TERM terminate as 129/130/143. Test every signal while a valid Git fallback exists, asserting no published target and no success report.
> - **Correct fallback branch:** Preserve the branch parsed by the existing canonical tarball-URL grammar and pass it to a shallow single-branch clone. A missing requested ref fails instead of falling back to default HEAD; URLs without a derivable branch retain current default-HEAD behavior. Test distinct default/requested branch markers under forced HTTP failure.
> - **Lessons and release trail:** Add concise cache-invalidation guidance to the Kanban prime and signal/ref-preservation guidance to the shell prime, with detailed linked entries in their lesson satellites. Reuse the existing REQ-303 consumer-corpus lesson rather than duplicating it. Update all three version locations plus both byte-identical changelogs.
>
> ## Test Plan
>
> - Run focused uncached Go tests for the three consumer corpus checks and the live cache-hit timeline regression.
> - Run `_dev/tests/update-script-behavior.sh`, covering HUP/INT/TERM statuses and non-default-branch archive contents.
> - Verify the root and shipped fetcher scripts are byte-identical.
> - Run `go test -count=1 ./...` in the queue-kanban module.
> - Run the canonical `bash _dev/tests/maintainer-verify.sh` and require exit zero.
> - Commit the verified implementation, then record its hash in the completed REQs through the separate bookkeeping commit required by the repository workflow.
>
> ## Assumptions
>
> - Scope is this repository only; the cited sibling checkout will receive the fixes through the normal suite update path.
> - Branch handling stays within the URL shape already recognized by the fetcher; query strings, fragments, and new URL grammars are out of scope.
> - Any user changes appearing during implementation are preserved and excluded from staging unless they overlap these fixes.
> ````

---
*Captured: 2026-08-29T20:26:10Z*
