---
id: REQ-198
title: Publish generated directory only after image success
status: pending
domain: general
created_at: 2026-08-15T16:34:08Z
user_request: UR-042
addendum_to: REQ-189
review_generated: true
effort_estimate: normal
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md]
tdd: true
maintenance: true
---

# Review Fix: Publish Generated Directory Only After Image Success

## What

Align the ai-report image-generation shell block with the conditional bundle contract: an all-failed generation attempt must not leave an empty published `generated/` directory, while successful images still publish there with status-backed freshness.

This is a standalone user-visible artifact-shape contract and cannot fold into a sweep: its root cause is the image helper's publication timing, not target parsing or a repeated prose class.

## Context

Found during review of REQ-189. `ai-report-reference.md` currently creates `ai-reports/<report-slug>/generated/` before any optional backend succeeds, although the action, reference output format, and guide say that directory exists only when current-run generated images succeed.

## Requirements

- Delay publication until at least one generation succeeds, or remove the empty current-run directory after all attempts fail.
- Preserve parallel PID/status tracking, invocation-private targets, absolute helper paths, stale-target rejection, and SVG/Mermaid fallback.
- Add or identify a replayable shell behavior assertion for the all-failed case and retain existing mixed-success behavior.

## Red-Green Proof

**RED prompt/case:** Run or inspect an ai-report generation attempt in which every optional image backend job fails; the shell block creates an empty published `generated/` directory before success is known.
**Why RED now:** The resulting bundle contradicts the documented conditional output shape and looks like missing or stale evidence.
**GREEN when:** An all-failed attempt leaves no published `generated/` directory, successful images still publish there, and the existing per-job status/freshness checks remain intact.
**Validation:** Review finding; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.
