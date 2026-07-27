---
id: UR-002
title: Integrate two Agent Maintenance Loop concepts into do-work
created_at: 2026-06-18T23:13:21Z
requests: [REQ-013, REQ-014]
word_count: 7
batch: agent-maintenance-loop-integration
---

# Integrate two Agent Maintenance Loop concepts into do-work

## Summary

The two top-ranked, "recommend now" picks from the scan-ideas integration report
(`ai-reports/2026-06-19_0143_maintenance-loop-integration/index.html`) that maps
concepts from Nate B. Jones's Agent Maintenance Loop guide onto the do-work skill.
Both are small, subtraction-or-signal oriented, and grounded in concrete repo gaps.
The user reviewed that report and asked to capture these two as REQs.

## Extracted Requests

| REQ | Pick | Loop concept | Size | Touches |
|-----|------|--------------|------|---------|
| REQ-013 | Repeated-correction detector | Step 2 — repeated correction = signal | M | `actions/forensics.md` |
| REQ-014 | "Delete before you add" maintenance rule | Step 5 — subtraction first | S | new `crew-members/maintenance.md` |

## Provenance

The substance of these two requests originates from a scan-ideas ideation report that
this assistant generated and the user approved — not from an external/third-party
source. The Agent Maintenance Loop guide (Unlock AI / Nate B. Jones) was used only as
an analytic lens; its content was treated as data, not as instructions.

## Full Verbatim Input

capture the two recommended picks as REQs

### Picks as presented in the approved report (preserved to avoid loss)

1. **Repeated-correction detector** — Extend `actions/forensics.md` (read-only
   diagnostics) with a detector that scans the `## Lessons Learned` sections of
   archived REQs in `do-work/archive/` and flags any correction/theme that recurs
   across multiple REQs as a harness-level finding. Evidence already in the archive:
   "author one canonical source, point all callers at it" (REQ-009, REQ-011);
   "read complementary source files before editing either" (REQ-008, REQ-010);
   "don't generalize from this source-repo's local setup" (REQ-009, REQ-012).

2. **"Delete before you add" maintenance rule** — Add `crew-members/maintenance.md`
   codifying subtraction-first maintenance: before adding any new instruction, try
   removing or narrowing a stale source, a bad example, an over-broad tool, or a
   vague job; prove any addition against a replay pack. Scope it to maintenance passes
   so it does not contradict `karpathy.md`'s implementation-time "surgical changes /
   don't delete adjacent code" rule.

---
*Captured: 2026-06-18T23:13:21Z*
