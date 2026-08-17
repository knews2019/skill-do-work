---
id: REQ-221
title: Extract the ai-report image batch into a shipped script
status: pending
status_changed_at: 2026-08-17T19:58:23Z
domain: general
created_at: 2026-08-17T18:37:31Z
user_request: UR-042
addendum_to: REQ-204
review_generated: true
effort_estimate: normal
prime_files: [_dev/primes/prime-shell-commands.md, _dev/primes/prime-action-files.md]
tdd: true
maintenance: true
---

# Discovered Task: Extract the AI-Report Image Batch Into a Shipped Script

## What

Move the prescribed image-batch block out of `ai-report-reference.md` and into `<skill-root>/scripts/generate-report-image-batch.sh`, leaving the action file with a pointer plus the per-report prompts.

## Context

Found while implementing REQ-204. After adding process-tree ownership and publication verification, the fenced block is roughly 110 lines of purely mechanical shell — staging, job control, group verification, signal handling, wait-all, freshness evaluation, publication, rollback — living inside a markdown action file. The only per-report content in it is the two `launch_report_image … "<prompt N>"` lines.

The test suite already treats it as a program: `_dev/tests/prescribed-shell-scripts-behavior.sh` extracts the block with `awk`, rewrites `<skill-root>`, and executes it. That extraction step exists only because the code is embedded in prose.

## Requirements

- The script owns the mechanics: staging, launch, isolation verification, signal ownership, wait-all, per-status freshness, publication, and rollback.
- The action file keeps only what is per-report — the style brief, the prompts, and how the results are referenced from the HTML — plus a pointer to the script.
- Preserve every behavior REQ-198 and REQ-204 locked in; the four existing replays must pass against the script with their assertions unchanged in meaning.
- The replays call the script directly; the `awk` block-extraction harness goes away with the block it was extracting.
- Keep the shipped-package manifest, `suite/modules.tsv`, and the shipped reference contract consistent with the new script.

## Red-Green Proof

**RED prompt/case:** Point the four existing batch replays at a shipped `generate-report-image-batch.sh`; there is no such script, so they cannot run.
**Why RED now:** The mechanics have no executable home, so the only way to test them is to reconstruct them from prose at test time.
**GREEN when:** All four replays invoke the script directly and pass, the action file no longer contains a fenced batch implementation, and the full maintainer baseline is green.
**Validation:** Discovered task; apply `actions/work-reference.md` → **Finding-Closure Ratchet (Step 6.5)**.

## Open Questions

- [x] The image-batch mechanics are now ~110 lines of shell embedded in a markdown action file, and the tests already have to extract them to run them. Should I process this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.
  Why this is yours: this is a structural move, not a defect fix — nothing is broken today, so it is a judgment call about where the mechanics should live rather than something the builder should decide alone. It also touches the shipped package inventory, which is a user-visible surface.
