---
id: REQ-271
title: Make the read-side layout lock-step see every layout, in any spelling
status: pending
created_at: 2026-08-18T22:57:26Z
status_changed_at: 2026-08-18T22:57:26Z
user_request: UR-056
addendum_to: REQ-257
domain: general
review_generated: true
effort_estimate: normal
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: true
write_set:
- _dev/tests/prescribed-shell-scripts-behavior.sh
---

# Make the Read-Side Layout Lock-Step See Every Layout, in Any Spelling

## What

REQ-257 decided the offset/fractional refusal is permanent, and pinned it with a lock-in that is supposed to fail if the board's `parseTimestamp` layouts change underneath the decision. **That guard is blind to most spellings it exists to catch, and it is not portable.** One line carries both defects:

```
| sed -n 's/^[[:space:]]*\(time\.RFC3339\|"2006[^"]*"\),$/\1/p'
```

1. **It only captures `time.RFC3339` or a `"2006…"`-prefixed literal.** REQ-257's review added `time.RFC3339Nano` and `time.DateTime` to `parseTimestamp`'s layout slice and **the suite stayed green** — same for `time.RFC1123`, `"02/01/2006 15:04:05"` and `"Jan 2, 2006"`. Only a `"2006…"` literal fires it. `time.RFC3339Nano` is precisely the layout someone reaching for fractional-second support would add, so the guard is blind in its own headline scenario.
2. **`\|` is GNU BRE alternation.** On BSD/macOS `sed` the pattern matches nothing, the extracted list comes back empty, and the case `fail_case`s spuriously — **the whole maintainer gate goes red on macOS**. This is the only `\|`-in-sed construct in the entire `_dev/tests/` tree, and the repo already carries a macOS-portability lesson (bash 3.2, REQ-216).

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Context

REQ-257's independent review, Important findings 1 and 2 — one root cause, one line, so one REQ. Both are `_dev/primes/prime-shell-commands.md` § Closed Enumerations Go Stale enacted **inside a lock-in**: the guard enumerates the spellings it already knew about instead of keying on the condition. It is also the prime's REQ-244 lesson repeating — *a detector that only recognizes the spellings it already fixed locks in nothing*.

Worth stating plainly: a guard that cannot fail is worse than no guard, because it is read as coverage. REQ-257's hand-back D-02 claims "a new layout there fails the suite and forces the decision to be re-made", and that claim is true for exactly one spelling family.

## Requirements

- The extraction keys on the **condition** — every element line inside `parseTimestamp`'s layout slice, whatever its spelling — never on an enumeration of the spellings that happen to be there today.
- No GNU-only regex construct. The gate must behave identically under BSD/macOS `sed`.
- `bash _dev/tests/maintainer-verify.sh` exits 0.

## Red-Green Proof

**RED:** add `time.RFC3339Nano` to `parseTimestamp`'s layout slice in `skills/do-work-board/tools/queue-kanban/model.go` and run `bash _dev/tests/prescribed-shell-scripts-behavior.sh`. Observed today: **66 named cases, exit 0** — the guard does not notice. Revert the layout afterwards.

**GREEN:** the same mutation fails the read-side-layout case, naming the changed layout list; the unmutated tree still passes; and the extraction contains no GNU-only construct.

## Open Questions

None — the defect and its fix are both mechanical, and the review reproduced each by execution.
