---
title: "Lessons from REQ-410: Implement doctor, deterministic forensics, and metadata repairs"
type: source-summary
topic_cluster: verification-and-testing
sources: [raw/processed/2026-09-01/REQ-410-implement-doctor-deterministic-forensics.md]
related:
  - page: concept-contract-verification-gates
    rel: evidence-for
created: 2026-09-01
updated: 2026-09-01
confidence: medium
---

# Lessons from REQ-410: Implement doctor, deterministic forensics, and metadata repairs

Part of the [[concept-contract-verification-gates]] cluster.

## What the REQ was about

Create `doctor` and move deterministic forensic checks and safe metadata repairs into Go.

## Solution summary

Added canonical `doctor` diagnosis and guarded timestamp repair using one typed result, one repository snapshot, explicit full-history recovery provenance, strict timestamp compatibility evidence, and shared Git transactions. Registered the command and made the natural-language forensics action and guide delegate deterministic work to it while retaining judgment-only recurring lessons and board-owned verification.

## What worked

**What worked:** One typed snapshot/result path, exact Git guards, live-corpus acceptance, and named adversarial regressions removed duplicate mechanical authorities and closed the initial false-positive and committed-state defects.

**What didn't:** Migrating the producer without tracing every downstream report field and reference left the action contract half-landed. Treating “parseable” as “eligible for repair ordering” also let a diagnosis-only timestamp influence a mutation.

**Worth knowing:** A canonical-tool migration is complete only when every consumer can produce its required output from that tool and all legacy anchors are swept. For repair logic, comparison eligibility must be the same supported-shape predicate as mutation eligibility; unsupported evidence must remain observational and byte-identical.

## Back-reference

See `do-work/archive/REQ-410-implement-doctor-forensics.md` for the full REQ — plan, exploration, implementation, review, and lessons. Commit `210d1459`.
