---
id: REQ-720
title: '[impact-negligible] Count work-reference.md sections and always-loaded crew files in the reading-volume measurement'
status: claimed
created_at: 2026-11-12T09:05:50Z
user_request: UR-171
domain: general
prime_files: []
tdd: false
maintenance: false
impact: impact-negligible
effort_estimate: effort-substantive
related: [REQ-719]
batch: phase-loading-follow-ups
write_set: ["_dev/tests/instruction-reading-volume.py", "_dev/reports/work-phase-loading.md"]
claimed_at: 2026-11-12T09:18:38Z
---
# Count work-reference.md Sections and Always-Loaded Crew Files in the Reading-Volume Measurement
## What
Extend `_dev/tests/instruction-reading-volume.py` and the report `_dev/reports/work-phase-loading.md` so the three walkthroughs also count `work-reference.md`, by cited section, and the crew files that are always loaded. Measurement only: no shipped file changes.
## Why
The current totals leave out the largest instruction file. The report excludes `work-reference.md` (119,671 bytes), yet the selection phase file sends even a simple fix to two of its sections that total 52,760 bytes. The phase split saved 12,693 bytes, so the reported saving looks larger than it is against what an agent really reads.
## Context
- This request comes from Finding 5 of the validate-feedback triage of 2026-11-11 (source: pasted text, no severity given). The triage verdict was Accept, and its surface-cost result was N/A: this is a measurement fix, not a new defensive layer.
- `_dev/reports/work-phase-loading.md:18` states the current boundary: the router, `work.md`, the phase files and `review-work.md`.
- `skills/do-work/actions/work-selection-reference.md:53` points a simple fix at the Request File Schema section (35,527 bytes) and the Schema Read Contract section (17,233 bytes) of `skills/do-work/actions/work-reference.md`.
- The script compares a fixed baseline revision with the files on disk (`_dev/tests/instruction-reading-volume.py:11`), and the same file holds the preservation replay for the phase split. The builder edits that file.
- Fold-first scan on 2026-11-12: the queue held no pending or pending-answers REQ, so no candidate in any UR shares this root cause.
## Detailed Requirements
- Extend `_dev/tests/instruction-reading-volume.py` and the report `_dev/reports/work-phase-loading.md`.
- Count `work-reference.md` by cited section: the sections that each walkthrough's loaded files point to, not the whole file and not zero.
- Count the always-loaded crew files.
- Keep the same three walkthroughs: a simple mechanical fix, a complex multi-file change, and a review-only run.
- Record, for each walkthrough, the exact instruction files loaded and the total bytes read.
## Constraints
- Measurement only. No shipped file changes: nothing under `skills/` changes.
- Byte or token volume is used to spot unnecessary reading in unrelated phases. It is never proof of instruction quality and never a reason to delete a required rule.
## Red-Green Proof
**RED prompt/case:** Run `python3 _dev/tests/instruction-reading-volume.py` from the repository root.
**Why RED now:** The Simple fix row ends at 133281 bytes. No output line names `work-reference.md` or a crew file, so those bytes are not counted.
**GREEN when:** For each of the three walkthroughs, the script lists the cited `work-reference.md` sections and the always-loaded crew files with their bytes, and its totals include them. The table in `_dev/reports/work-phase-loading.md` shows the same totals the script prints. This script run, compared with the report table, is the regression check for the finding.
**Validation:** User confirmed on 2026-11-12: the proof is the script output and the report table together.
## Dependencies
None. REQ-719 (the prime byte budget) belongs to the same batch and does not depend on this request.
## Builder Guidance
The user is certain about what to count and where. How the script finds which `work-reference.md` sections a walkthrough cites, and which crew files count as always loaded for each walkthrough, are the builder's decisions; state both in the report so a reader can check them. The finding's verbatim claim also names checks, recovery rules and a `code-review` path; the remedy the user captured is the counting work above.
## Full Context
See `do-work/user-requests/UR-171/input.md` for complete verbatim input.
## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: Finding 5 of the 2026-11-11 phase-loading triage, as pasted by the user in UR-171: "Record the exact instruction files loaded, total bytes read, and every required check and recovery rule reached."*
