---
id: REQ-260
title: Run the Go formatter as part of the canonical verify
status: claimed
claimed_at: 2026-08-18T21:16:24Z
route: A
created_at: 2026-08-18T18:41:26Z
status_changed_at: 2026-08-18T20:59:31Z
user_request: UR-051
addendum_to: REQ-251
domain: general
effort_estimate: normal
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: false
suggested_spec:
depends_on: []
maintenance: false
write_set:
- _dev/tests/maintainer-verify.sh
estimate:
  p50_active_minutes: 15
  confidence: medium
  calculated_at: 2026-08-18T21:17:28Z
  basis:
    - Route A
    - 1-file write set
    - 3 acceptance criteria
    - full-suite verification
---

# Run the Go Formatter as Part of the Canonical Verify

## What

`bash _dev/tests/maintainer-verify.sh` runs `go vet` but never the Go formatter, so a formatting slip lands and survives silently — one did, in the board tool's day-domain truncation expression, and was only caught by a builder reading adjacent code. That specific slip is already fixed (REQ-252 corrected it in passing), so this REQ is the rule rather than the instance: the canonical gate should run the formatter over tracked Go files and fail on a non-empty result, exactly as it already runs ShellCheck over tracked shell files.

## Requirements

- The gate runs the Go formatter over tracked Go files and fails when any is unformatted, selecting files the way the ShellCheck lane does (`git ls-files`, never a hand-maintained path list — Closed Enumerations Go Stale).
- The failure names the offending files.
- `bash _dev/tests/maintainer-verify.sh` exits 0 on the current tree — the package is formatted clean today, so a red result means the check found something real.

## Context

Discovered by REQ-251's builder ([low]); the one-character instance was fixed in passing by REQ-252. The user widened the scope at clarify: the instance is closed, the gap that let it through is not.

## Open Questions

- [x] I discovered this out-of-scope task while working on REQ-251: a gofmt formatting miss in `durations.go` from REQ-248. Should I process this as a new task? → Yes, and widened: add a formatter check to the canonical verify so this class cannot recur silently
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it — or fold it into REQ-252, which already owns the file.

**Answered [2026-08-18]:** User approved via `do-work clarify` and **widened the scope**: the original one-character fix is already satisfied by REQ-252's in-passing correction, so this REQ now covers making the gate run the formatter over tracked Go files. Title and requirements updated to match; `effort_estimate` raised from trivial to normal.

---

## Triage

**Route: A** - Simple

**Reasoning:** One named file, and the change mirrors a pattern that already exists inside it (the ShellCheck lane's `git ls-files` selection and fail-on-non-empty shape). Scope is obvious and bounded.

**Planning:** Not required
