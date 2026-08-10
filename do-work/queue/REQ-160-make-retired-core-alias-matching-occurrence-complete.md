---
id: REQ-160
title: "Review fix: Make retired core alias matching occurrence-complete"
status: pending
status_changed_at: 2026-08-10T09:20:51Z
domain: general
created_at: 2026-08-09T20:36:33Z
user_request: UR-031
addendum_to: REQ-157
review_generated: true
effort_estimate: normal
sweep: true
sweep_key: retired-core-alias-match-occurrence-completeness
---

# Review Fix: Make Retired Core Alias Matching Occurrence-Complete

## What

Make the test-only retired-command guard evaluate command occurrences and overlapping trigger candidates completely. An exemption or boundary-invalid longer candidate must not hide another valid retired command. Done means this false-negative class cannot recur while direct-alias suffix exclusions remain intact.

## Context

Found during review of REQ-157. Its 186-row historical inventory is complete, but the matcher can still miss retired invocations when a line also contains an exempt queue-board test reference or when a longer install-target prefix fails its right boundary before a shorter historical install/setup head is considered.

## Requirements

- Scope queue-board branding and test-reference exemptions to the exact matched occurrence; a second retired command on the same line must still fail.
- Continue candidate evaluation after a boundary-invalid longer trigger when a shorter historical command head remains valid.
- Treat the former `install-<target>` routing prefix consistently, including unknown targets formerly routed to shim help.
- Preserve negative behavior for embedded/prefixed words, possessives, and direct-alias suffixes.
- Add adversarial regressions for exempt and genuine occurrences sharing one line, overlapping install-target prefixes, unknown `install-` targets, and current sibling commands.
- Keep all aliases test-only and preserve the 186-row inventory, current sibling ownership/routes, prime fingerprints, repaired live surfaces, and full distribution tests.

## Instances

- [ ] A line containing `<title> text "do-work queue board"` plus a second `do-work queue board` invocation must report the second occurrence.
- [ ] A line containing `strings.Contains(bodyText, "do-work queue board")` plus a second retired invocation must report the second occurrence.
- [ ] `do-work install ui-design2` and `do-work setup ui-design2` must fall back to the valid historical install/setup head.
- [ ] `do-work install-custom-target` must be recognized through the former `install-<target>` prefix without turning direct-alias suffixes into positives.

## Open Questions

- [x] REQ-157 completes the historical alias inventory, but its guard still misses some former install-head commands and additional retired commands placed on exempt test-reference lines. Because REQ-157 was itself created by a review, the cascade-depth rule requires your approval before this non-critical follow-up can enter the claimable queue. Should I process this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.
