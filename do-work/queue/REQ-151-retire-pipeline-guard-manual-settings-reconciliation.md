---
id: REQ-151
title: "Review fix: Retire the pipeline guard in manual settings reconciliation"
status: pending
domain: general
created_at: 2026-08-08T17:47:25Z
user_request: UR-031
addendum_to: REQ-145
review_generated: true
effort_estimate: normal
---

# Review Fix: Retire the Pipeline Guard in Manual Settings Reconciliation

## What
Make the suite installer's no-JSON-tool fallback explicitly remove only the retired pipeline-guard Stop hook while preserving every unrelated and custom hook entry.

## Context
Found during review of REQ-145. The jq and Python reconciliation paths correctly remove the retired guard, but the fallback instructions for systems with neither tool still say to preserve every existing entry. That leaves a dangling Stop hook pointing at a script REQ-145 deletes.

This is a standalone user-visible upgrade-path defect rather than part of a broader sweep: both byte-identical installer copies and their one fallback regression close it.

## Requirements
- Update `tools/install-do-work-suite.sh` and `skills/do-work/tools/install-do-work-suite.sh` identically.
- In the no-jq/no-Python fallback, instruct the user to remove only Stop-hook objects whose command targets `.claude/skills/do-work/hooks/pipeline-guard.sh`.
- Preserve every unrelated/custom hook, including custom Stop hooks sharing the same event.
- Update the fallback regression so it rejects instructions that preserve the retired guard and proves custom hooks remain protected.
- Keep fresh-install, jq, Python, idempotence, and installer byte-identity tests passing.
