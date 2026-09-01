---
title: "Metadata and Timestamps"
type: concept
topic_cluster: metadata-and-timestamps
sources:
  - raw/processed/2026-09-01/REQ-074-recovered-req-loses-its-status-change-ti.md
  - raw/processed/2026-09-01/REQ-076-go-utility-emits-the-canonical-utc-tim.md
  - raw/processed/2026-09-01/REQ-078-the-windows-timestamp-fallback-cannot-ru.md
  - raw/processed/2026-09-01/REQ-100-live-auto-wave-acceptance-run-prove-real.md
  - raw/processed/2026-09-01/REQ-118-the-normalize-flag-must-stop-calling-voc.md
  - raw/processed/2026-09-01/REQ-128-secret-rename-quarantine-survives-re-inv.md
  - raw/processed/2026-09-01/REQ-148-addendum-preserve-association-candidates.md
  - raw/processed/2026-09-01/REQ-247-archive-timestamp-audit-tool-driven-by-g.md
  - raw/processed/2026-09-01/REQ-251-retire-the-stale-copies-of-the-future-st.md
  - raw/processed/2026-09-01/REQ-253-decide-the-timestamp-rule-s-two-uncovere.md
  - raw/processed/2026-09-01/REQ-255-give-the-timestamp-repairer-shape-parity.md
  - raw/processed/2026-09-01/REQ-280-probe-timestamp-ordering-and-point-check.md
  - raw/processed/2026-09-01/REQ-281-reconcile-the-calibration-log-against-th.md
  - raw/processed/2026-09-01/REQ-308-judge-effort-estimate-on-every-capture-a.md
  - raw/processed/2026-09-01/REQ-310-check-a-template-payload-s-citations-aga.md
  - raw/processed/2026-09-01/REQ-314-judge-effort-estimate-on-review-minted-f.md
  - raw/processed/2026-09-01/REQ-316-audit-the-calibration-log-write-step-for.md
related:
  - page: concept-queue-task-lifecycle
    rel: complements
created: 2026-09-01
updated: 2026-09-02
confidence: high
---

# Metadata and Timestamps

Architectural overview and synthesis for the Metadata and Timestamps subsystem in the do-work suite.

## Key Principles & Synthesized Lessons

This cluster synthesizes evidence from 17 source documents:

- [[REQ-074-recovered-req-loses-its-status-change-ti]] — A recovered REQ loses the timestamp that says when it was reset
- [[REQ-076-go-utility-emits-the-canonical-utc-tim]] — Go utility emits the canonical UTC timestamp, preferred over date -u when built
- [[REQ-078-the-windows-timestamp-fallback-cannot-ru]] — The Windows timestamp fallback cannot run on stock Windows in either shell it names
- [[REQ-100-live-auto-wave-acceptance-run-prove-real]] — Live auto-wave acceptance run — prove real wall-clock concurrency
- [[REQ-118-the-normalize-flag-must-stop-calling-voc]] — The normalize flag must stop calling vocabulary-less field values unrecognized
- [[REQ-128-secret-rename-quarantine-survives-re-inv]] — Secret rename quarantine survives re-inventory
- [[REQ-148-addendum-preserve-association-candidates]] — Addendum: preserve association candidates with empty quarantine
- [[REQ-247-archive-timestamp-audit-tool-driven-by-g]] — Archive timestamp audit tool driven by git commit times
- [[REQ-251-retire-the-stale-copies-of-the-future-st]] — Retire the stale copies of the future-stamp message
- [[REQ-253-decide-the-timestamp-rule-s-two-uncovere]] — Decide the Timestamp rule's two uncovered stamp shapes
- [[REQ-255-give-the-timestamp-repairer-shape-parity]] — Give the timestamp repairer shape parity with the read-side detectors
- [[REQ-280-probe-timestamp-ordering-and-point-check]] — Probe timestamp ordering, and point Check 12 at the archive repair that already ships
- [[REQ-281-reconcile-the-calibration-log-against-th]] — Reconcile the calibration log against the frontmatter it was derived from
- [[REQ-308-judge-effort-estimate-on-every-capture-a]] — Judge effort_estimate on every capture, as impact already is
- [[REQ-310-check-a-template-payload-s-citations-aga]] — Check a template payload's citations against where the payload lands
- [[REQ-314-judge-effort-estimate-on-review-minted-f]] — Judge effort_estimate on review-minted follow-ups too
- [[REQ-316-audit-the-calibration-log-write-step-for]] — Audit the calibration-log write step for the REQ-274 stale-stamp bug class

## Cross-References

See related system components and verification gates.
