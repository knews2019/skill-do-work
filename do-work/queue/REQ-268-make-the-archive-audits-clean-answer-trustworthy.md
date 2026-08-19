---
id: REQ-268
title: Never report clean for a state that was never verified
status: pending
created_at: 2026-08-18T21:03:15Z
status_changed_at: 2026-08-18T22:20:09Z
user_request: UR-056
addendum_to: REQ-255
domain: general
review_generated: true
effort_estimate: normal
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
write_set:
- skills/do-work/scripts/audit-archive-timestamps.sh
- skills/do-work/scripts/repair-req-timestamps.sh
- _dev/tests/prescribed-shell-scripts-behavior.sh
---

# Never Report Clean for a State That Was Never Verified

## What

**The condition, not the file:** an unchecked exit status turns an inspection that never happened into a clean answer. Three instances are known, in two scripts, and all three are reproduced by execution. The REQ was originally scoped to the archive auditor; an external review found the same root cause in the repairer, which is why it is now keyed on the condition instead.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Instances

- [ ] **A refused stamp is reported as clean.** An archive holding only a calendar-impossible stamp yields `archive audit clean (1 file(s) scanned)` and exit 0, in both report-only and fixing modes. The refusal itself is right — the malformed value is deliberately preserved as evidence — but the reasoning for preserving it is that a human can then see it, and the one tool a person runs deliberately to inspect the archive is exactly where it disappears. Voicing refusals as informational lines, without changing the exit status, keeps the refusal and the report contract consistent.
- [ ] **A failed scan is reported as clean.** The file walk runs inside a process substitution, so a nonzero exit from the walk or the sort never reaches the loop's status: with a failing `find` on PATH the auditor prints `archive audit clean (0 file(s) scanned)` and exits 0 while a defective archive file sits untouched. Materialise and validate the walk's output before entering the loop, so an incomplete scan can never be reported as a clean one.

- [ ] **The repairer reports clean when its own field extraction fails (added from an external review on PR #145, reproduced by the orchestrator).** `skills/do-work/scripts/repair-req-timestamps.sh:362` assigns `field_rows="$(extract_timestamp_fields "$request_file")"` and then tests only `[ -n "$field_rows" ]`, so a nonzero `awk` exit is discarded and empty output reads as "this file has no timestamp fields". Observed with a failing `awk` first on PATH: a queue file carrying `created_at: 2093-01-01T00:00:00Z` came back **byte-identical and the script exited 0 with no output at all** — while the same file with a working `awk` is repaired. The SessionStart hook discards the script's stderr, so nothing reaches the banner either. This is the third face of the same condition, in the second script.

## Requirements

- **"Clean" is printed only when the inspection actually completed** — every archive file read, every field extraction succeeded — and nothing needed repair. State this as the condition in each script, so a fourth call site inherits it.
- **Sweep the primitive:** every place either script takes a command substitution or a process substitution and then judges only the *content* of the output. These three were found by two independent reviewers looking at other things, so they are a sample.
- A refused defect is visible in the tool's own output, with the exit contract stated to match whatever the builder chooses.
- A scan that could not complete fails loudly rather than reporting a count of zero as success.
- Lock-in cases for both, and `bash _dev/tests/maintainer-verify.sh` exits 0.

## Context

Instance 1 is REQ-255's independent review, finding I-3 (gate: user-visible — the audit's answer is its product, and it misleads exactly on the class it was told to preserve). Instance 2 is an external automated review on pull request 145, reproduced by the orchestrator against the shipped script before recording. Both live in the same file and share one root cause: the tool answers "clean" for states it never verified. Created `pending-answers` per the generation-≥2 cascade stop.

## Open Questions

- [x] The archive auditor says "clean" both when it refused a malformed stamp it deliberately preserved and when its file walk failed outright and scanned nothing. Should I process this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.

**Answered [2026-08-18]:** User approved via `do-work clarify`. Both instances stand: the refused-stamp report line and the failed-walk exit path. The exit contract for a voiced refusal is builder latitude, as the Requirements already say.
