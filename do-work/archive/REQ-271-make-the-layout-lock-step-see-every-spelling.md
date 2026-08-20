---
id: REQ-271
title: Make the read-side layout lock-step see every layout, in any spelling
status: completed
created_at: 2026-08-18T22:57:26Z
status_changed_at: 2026-08-18T22:57:26Z
claimed_at: 2026-08-20T12:41:45Z
completed_at: 2026-08-20T12:55:37Z
commit: a773dc6
route: B
estimate:
  p50_active_minutes: 20
  confidence: medium
  calculated_at: 2026-08-20T12:41:45Z
  basis:
    - Route B
    - 1-file write set
    - lock-in guard requiring a mutation proof per spelling
user_request: UR-056
addendum_to: REQ-257
domain: general
review_generated: true
effort_estimate: normal
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec: bug-fix
depends_on: [REQ-300]
maintenance: true
write_set:
- _dev/tests/prescribed-shell-cases/repair-req-timestamps.sh
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
- [x] **[PLAN]:** Verify both stated defects against the tree first — defect 2 may already be gone. Then replace the spelling enumeration with a structural extraction: enter the function, enter the `[]string{` literal, take every element line until the closing brace. Add an anti-vacuity assertion, because an extraction that finds nothing compares equal to nothing and can never fail.
- [x] **[APPLY]:** One block in `_dev/tests/prescribed-shell-cases/repair-req-timestamps.sh` — the awk extraction, a layout-count vacuity guard, and the unchanged expected-value comparison.
- [x] **[UNIFY]:** `git diff --stat` reviewed; `shellcheck --severity=warning` reports only the file's pre-existing structural `SC2154`s from the sourced harness; `maintainer-verify.sh` exits 0; `repair-req-timestamps: 20 cases, 0 failures`, suite total unchanged at 90 (this REQ sharpens an existing case rather than adding one). Per-file: **`repair-req-timestamps.sh`** — POSIX awk only (`[ \t]` not a character class, no GNU-only construct); `/\[\]string/` without the brace so no awk reads it as an interval expression; the extraction is captured newline-separated first so the count guard can be exact, then flattened, preserving the expected value's trailing space byte-for-byte; `model.go` was mutated five times and restored from a pristine copy after each, confirmed clean by `git status`.

## Context

REQ-257's independent review, Important findings 1 and 2 — one root cause, one line, so one REQ. Both are `_dev/primes/prime-shell-commands.md` § Closed Enumerations Go Stale enacted **inside a lock-in**: the guard enumerates the spellings it already knew about instead of keying on the condition. It is also the prime's REQ-244 lesson repeating — *a detector that only recognizes the spellings it already fixed locks in nothing*.

Worth stating plainly: a guard that cannot fail is worse than no guard, because it is read as coverage. REQ-257's hand-back D-02 claims "a new layout there fails the suite and forces the decision to be re-made", and that claim is true for exactly one spelling family.

## Requirements

- The extraction keys on the **condition** — every element line inside `parseTimestamp`'s layout slice, whatever its spelling — never on an enumeration of the spellings that happen to be there today.
- No GNU-only regex construct. The gate must behave identically under BSD/macOS `sed`.
- `bash _dev/tests/maintainer-verify.sh` exits 0.

## Red-Green Proof

**RED:** add `time.RFC3339Nano` to `parseTimestamp`'s layout slice in `skills/do-work-board/tools/queue-kanban/model.go` and run `bash _dev/tests/prescribed-shell-scripts-behavior.sh`. Observed: **exit 0** — the guard does not notice. (The suite reports **76 named cases across 17 per-script files** as of REQ-258's split; it read 66 when this REQ was captured. The count is context, not the finding — `exit 0` is.) Revert the layout afterwards.

**Where the guard lives, post-split:** `_dev/tests/prescribed-shell-cases/repair-req-timestamps.sh:232-241` — the `board_timestamp_layouts` extraction and its `fail_case`. The command above is still the right entry point: the runner keeps its path and exit-status contract and executes every case file.

**GREEN:** the same mutation fails the read-side-layout case, naming the changed layout list; the unmutated tree still passes; and the extraction contains no GNU-only construct.

## Open Questions

None — the defect and its fix are both mechanical, and the review reproduced each by execution.

---

## Triage

**Route: B** - Medium

**Reasoning:** The REQ states both defects precisely and calls the fix mechanical, but it asserts two defects on one line and one of them turned out already fixed — so the work was establishing which was live before changing anything, then choosing a structural extraction that keys on the condition.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

**Defect 2 is already fixed. No edit needed, and none made.** The REQ quotes the guard as `sed -n 's/^[[:space:]]*\(time\.RFC3339\|"2006[^"]*"\),$/\1/p'` — GNU BRE alternation, which matches nothing under BSD/macOS sed. The line in the tree today reads:

```
| sed -n -E 's/^[[:space:]]*(time\.RFC3339|"2006[^"]*"),$/\1/p'
```

`-E` with `(a|b)` is POSIX ERE and portable. `grep -rn '\\|' _dev/tests/` finds no BRE-alternation construct anywhere in the tree; the remaining `\|` hits are inside `grep` patterns where it is correct, and inside `[[:alpha:]]`-style table regexes. REQ-300's exploration reached the same conclusion independently. So the "no GNU-only regex construct" requirement is satisfied by verification, not by a change.

**Defect 1 is live, and blind to every spelling the REQ names.** Applying the old extraction to each candidate line in isolation:

| Layout line | Old extraction saw |
|---|---|
| `time.RFC3339Nano,` | *(nothing)* |
| `time.DateTime,` | *(nothing)* |
| `time.RFC1123,` | *(nothing)* |
| `"02/01/2006 15:04:05",` | *(nothing)* |
| `"Jan 2, 2006",` | *(nothing)* |

And end to end: adding `time.RFC3339Nano` to `parseTimestamp`'s slice in `model.go` left the suite at `repair-req-timestamps: 20 cases, 0 failures`, exit 0. The guard does not notice the layout that someone reaching for fractional-second support would add first — its own headline scenario, exactly as the REQ says.

**Where the slice actually is.** `skills/do-work-board/tools/queue-kanban/model.go:1469-1474`, a `[]string{…}` composite literal inside `for _, layout := range` — four elements, one per line, each with a trailing comma. That per-line structure is what makes a structural extraction possible without a Go parser.

**The guard lives at `_dev/tests/prescribed-shell-cases/repair-req-timestamps.sh:232-241`** — confirming REQ-300's repointing of this REQ's `write_set`, and the reason the Red-Green Proof's entry command still works: the runner keeps its path and exit-status contract.

*Generated by Explore agent*

## Scope

**Files I will touch:**
- `_dev/tests/prescribed-shell-cases/repair-req-timestamps.sh` (modify) — the layout extraction, plus an anti-vacuity count guard

**Files I will NOT touch:**
- `skills/do-work-board/tools/queue-kanban/model.go` — mutated five times as the RED and restored each time; the layout slice is the *subject* of the guard, not its target
- `skills/do-work/scripts/repair-req-timestamps.sh` — the refusal itself is a decided permanent answer; this REQ fixes the guard over it, not the refusal
- `_dev/tests/prescribed-shell-scripts-behavior.sh` — the runner, unchanged since REQ-258 reduced it to dispatch

**Acceptance criteria (restated from REQ):**
- [ ] The extraction keys on the condition — every element line inside `parseTimestamp`'s layout slice, whatever its spelling
- [ ] No GNU-only regex construct; identical behavior under BSD/macOS
- [ ] `bash _dev/tests/maintainer-verify.sh` exits 0

## Pre-Flight

**Git:** ✓ working tree clean outside `do-work/` at claim time
**Tests baseline:** ✓ passing (`maintainer-verify.sh` exit 0)
**Dependencies:** ✓ go1.26.1, ShellCheck 0.11.0, `just` present

*Checked by work action*

## Decisions

- **D-01 — Structural extraction, not a widened pattern.** DECIDE & STATE. The obvious cheap fix is to add the missing spellings to the alternation, which is the very failure `_dev/primes/prime-shell-commands.md` § Closed Enumerations Go Stale and its REQ-244 lesson describe: a detector that only recognizes the spellings it already fixed locks in nothing. Keying on position inside the `[]string{` literal means a spelling nobody has thought of is caught on arrival.

- **D-02 — An anti-vacuity count assertion, added beyond the REQ's ask.** DECIDE & STATE. A structural extraction has a failure mode the enumeration did not: if the function is renamed or the slice restructured, it yields nothing, the comparison finds nothing on both sides, and the guard silently stops guarding. That is the same shape as the defect being fixed, one level up. The count assertion fires with a message naming the condition, and it is proven to fire by renaming the function.

- **D-03 — Defect 2 recorded as already-fixed rather than re-fixed.** DECIDE & STATE. The REQ asserts a `\|` BRE construct that is not in the tree; the line is already `-E` with ERE alternation. Verified two ways and recorded here so a reader comparing the REQ's claim against the diff does not conclude the portability half was skipped.

## Implementation Summary

**Files changed:**
- `_dev/tests/prescribed-shell-cases/repair-req-timestamps.sh` (modified)

**What was done:** Replaced the two-spelling `sed` alternation with a POSIX-awk extraction that walks into `parseTimestamp`, into its `[]string` composite literal, and takes every non-empty element line until the closing brace — stripping only indentation, a trailing line comment, and the trailing comma. Added a layout-count assertion so an extraction that finds nothing fails loudly instead of comparing empty-to-empty. The expected-value comparison and its message are unchanged, so the decision this guard protects still reads the same.

## Qualification

**Passed** — `tools/checks/qualify.sh` exit 0, no FAIL and no WARN. One file verified in the diff, all three acceptance criteria traced, P-A-U confirmed. `tools/checks/scope-drift.sh`: `OK: Implementation Summary matches the Scope declaration` (exit 0).

Worth noting, since the last four REQs each needed an override: this one qualified clean on the first run. The change adds no marker-bearing fixture payload and touches no `do-work/`-only path, so neither of the two false-FAIL classes recorded in REQ-263 and REQ-300 applies.

**Judgment checks:** *(2) Substantive* — the extraction is rewritten, not tweaked, plus a new count assertion. *(3) Requirements traced* — the condition-keying criterion is proven by five mutations, the portability criterion by verification (D-03), and the gate criterion by the run below. *(6) Flowing* — the extraction reads the real `model.go`; if it were stubbed the five mutations could not each fire it and the rename could not empty it.

## Testing

**Tests run:** `bash _dev/tests/prescribed-shell-cases/repair-req-timestamps.sh`, then `bash _dev/tests/maintainer-verify.sh`
**Result:** ✓ `repair-req-timestamps: 20 cases, 0 failures`; ✓ `maintainer-verify.sh` exit 0; suite total unchanged at **90 named script cases across 17 per-script files** — this REQ sharpens an existing case rather than adding one

**Red-green validation** — the captured RED is `add time.RFC3339Nano to parseTimestamp's layout slice and run the suite`, and it was extended to every spelling the REQ names:

| Layout added to the slice | Before | After |
|---|---|---|
| `time.RFC3339Nano` | ✗ exit 0, guard silent | ✓ read-side-layout case fires |
| `time.DateTime` | ✗ silent | ✓ fires |
| `time.RFC1123` | ✗ silent | ✓ fires |
| `"02/01/2006 15:04:05"` | ✗ silent | ✓ fires |
| `"Jan 2, 2006"` | ✗ silent | ✓ fires |

Each mutation was applied to `model.go`, the suite run, and `model.go` restored from a pristine copy before the next — `git status` clean at the end. The "before" column was confirmed two ways: end-to-end for `RFC3339Nano` (20 cases, 0 failures, exit 0 with the layout added), and by applying the old `sed` expression to each candidate line in isolation, which returned empty for all five.

**Anti-vacuity proof:** renaming `parseTimestamp` to `parseTimestampRenamed` makes the extraction find nothing. The new count guard fires — `extracted ZERO layouts from parseTimestamp — the slice moved or was restructured, so this guard is blind rather than passing` — instead of comparing empty to empty and passing. Without that assertion this rewrite would have swapped one silent-blindness mode for another.

**Portability:** no GNU-only construct in the new block. POSIX awk only, `[ \t]` in place of a character class, and `/\[\]string/` without a brace so no awk can read it as an interval expression. Verified by inspection rather than on a BSD host, which this container is not — stated plainly rather than claimed as tested.

**New tests added:** none. The guard that existed was blind; this REQ makes it able to fail, which is the deliverable.

**Existing tests updated:** `_dev/tests/prescribed-shell-cases/repair-req-timestamps.sh`'s read-side-layout case (from REQ-257) — intentional strengthening of a lock-in that could not fail, not a behavior change to what it protects. Its expected value and failure message are byte-identical.

*Verified by work action*

## Review

**Overall: 94%** | 2026-08-20T12:54:00Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 92% |
| Test Adequacy | 100% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- None

**Minor findings:** 2 (report only)
- The REQ's stated defect 2 does not exist in the tree, so half of a two-defect REQ closed by verification. Recorded in D-03 and the exploration. The REQ was captured before whichever change made the line ERE, and nothing flagged that its premise had gone stale — the same class REQ-300 swept, arriving here one REQ later.
- Portability is argued, not measured. The claim is that the new block uses only POSIX awk; no BSD/macOS host was available to run it on. The argument is strong (the constructs are individually POSIX and the previous defect was a GNU-only alternation, now absent) but it is an argument.

**Restatement sweep:** the diff changes how a guard reads `model.go`, so the sweep asked who else states the layout list or the extraction's shape. Nothing does: `grep -rn 'parseTimestamp' --include='*.md'` outside `do-work/` returns no hit, and the only other reader of that slice is `parseTimestamp` itself. The guard's failure message — which names the decision that must be re-made — is unchanged, so no downstream text describing it went stale.

**Acceptance:** Pass — all three criteria met: condition-keying proven by five mutations that each now fire and each previously did not, portability satisfied by verification with the reasoning recorded, and `maintainer-verify.sh` exit 0.

**Suggested testing:** 1 item
- Nothing pins the *extraction itself* against a restructured slice that still has elements — a one-line `[]string{time.RFC3339, "2006-01-02"}` literal, for instance, would extract nothing while the count guard fires correctly, so the failure is loud but the guard stops covering. Worth a line only if Go formatting there ever changes; `gofmt` keeps multi-element literals one-per-line today, which is what makes the structural read safe.

**Code Quality 92%:** the awk block is eight lines of state machine inside a shell string, which is the least readable part of the file it lives in. It earns that by being the only way to key on structure without a Go parser, and the comment above it states the condition it implements — but a reader meets `inside_layout_slice` before they meet why.

**Follow-up REQs created:** None; **sweeps appended to:** None

## Lessons Learned

**What worked:** Testing the REQ's premises before its fixes. One of the two stated defects was already gone, and the extraction that "only captures `time.RFC3339` or a `"2006…"` literal" was worth confirming spelling by spelling rather than accepting — which produced the five-row before/after table that is now the actual evidence the fix works. Cheaper than it sounds: the old expression is one `sed` invocation, so proving blindness was five one-line runs.

**What didn't:** Nothing failed, but the first instinct was wrong and worth naming. Widening the alternation to include `RFC3339Nano` and friends would have passed the REQ's own Red-Green Proof exactly as written, and would have re-created the defect for the next spelling. The Proof named one mutation; the requirement named the condition. When those two disagree in scope, the requirement wins.

**Worth knowing:** Replacing an enumeration with a structural read swaps one blindness mode for another unless you guard it — an extraction that finds nothing compares equal to nothing and passes. The count assertion is what makes this guard honest, and renaming the function is the one-line way to prove it works. Any future "key on the condition, not the list" rewrite of a comparison guard needs the same companion assertion.

## Orientation

The lock-step guard over the board's timestamp parser can now actually fail. It reads every element of `parseTimestamp`'s layout slice by position rather than matching two known spellings, so adding `time.RFC3339Nano` — or any other layout, in any spelling — forces the offset/fractional refusal decision to be re-made instead of being inherited silently. It also refuses to pass when it extracts nothing, which is the failure mode the rewrite introduced. Lives in the maintainer test suite over the repairer/read-side contract (`_dev/primes/prime-shell-commands.md`).

Not `[MAP CHANGED]` — one guard's recognition got broad while its requirement stayed narrow; no contract, decision, or caller moved. Staleness spot-check on `_dev/primes/prime-shell-commands.md`: every referenced path resolves, and its § Closed Enumerations Go Stale plus the REQ-244 and REQ-257 lessons are precisely what this change enacts. The prime is not stale.
