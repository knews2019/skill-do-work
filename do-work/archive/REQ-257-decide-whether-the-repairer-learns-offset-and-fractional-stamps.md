---
id: REQ-257
title: Decide whether the timestamp repairer learns offset and fractional stamps
status: completed
completed_at: 2026-08-18T22:58:40Z
commit: 6afcbd5
claimed_at: 2026-08-18T21:16:24Z
route: B
created_at: 2026-08-18T17:49:24Z
status_changed_at: 2026-08-18T20:55:14Z
user_request: UR-056
addendum_to: REQ-246
domain: general
effort_estimate: normal
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec:
depends_on: [REQ-255]
maintenance: false
write_set:
- skills/do-work/scripts/repair-req-timestamps.sh
- _dev/tests/prescribed-shell-scripts-behavior.sh
estimate:
  p50_active_minutes: 25
  confidence: medium
  calculated_at: 2026-08-18T21:17:28Z
  basis:
    - Route B
    - 2-file write set
    - 3 acceptance criteria
    - cross-route regression gates
    - full-suite verification
---

# Decide Whether the Timestamp Repairer Learns Offset and Fractional Stamps

## What

REQ-246's repairer deliberately refuses stamps with numeric UTC offsets (`2093-01-01T00:00:00+02:00`) or fractional seconds — repairing them needs timezone arithmetic, and a wrong guess would rewrite a correct stamp (REQ-246 D-04, documented in the script header). **Corrected during the build, by execution:** the board and forensics do *not* warn on those shapes as such — they warn on them *when the instant is future*. An offset or fractional stamp that is simply correct is read correctly, produces right elapsed-time math, and is never flagged by anything. `detectFutureTimestampFields` appends only when the parsed instant is after the skew horizon; shape never enters the predicate. So the residual is not "offset/fractional stamps", it is "offset/fractional stamps that are **also** fabricated or misordered" — a much smaller and stranger population, and that narrowing is most of why the refusal wins. The original sentence read: *"The board and forensics still detect and warn on those shapes, so they remain a detection-without-repair residual."* **Correction (REQ-255 review, I-2):** these are no longer the *only* such residual — a quoted stamp with padding inside the quotes is also board-parseable and refused here; that one is tracked in REQ-267. This asks whether that residual matters enough to implement offset arithmetic in `comparison_key_for`, or whether the documented refusal is the permanent answer.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `_dev/primes/prime-shell-commands.md` including its REQ-255, REQ-250, REQ-246 and REQ-243 lessons, four crew-member files, the whole repairer, `audit-archive-timestamps.sh`, the board's `parseTimestamp` and `detectFutureTimestampFields` seams, forensics Check 11, and the existing `# repair-req-timestamps:` case group. Plan: establish **empirically** what the read side actually does with each shape — the REQ's premise, which turned out to be wrong — then decide on merits, then either implement or pin. In the pin case the RED has to come from a mutation, because a lock-in for existing behaviour cannot go red on its own.
- [x] **[APPLY]:** Stayed inside the two declared files. Four commits, each individually green.
- [x] **[UNIFY]:** `git diff --stat 662788c..HEAD` → 2 files, +76/−7; both diffs read in full. `bash -n` and `shellcheck -S warning` clean on the repairer (exit 0). The suite file's shellcheck output carries only pre-existing info-level findings (SC2016/SC2086/SC2015/SC2329) in fixture-writing lines from earlier REQs — none on any added line, and the gate's own warning-level lane passes. Grepped the added lines for `TODO|FIXME|XXX|MUTATION|console.log|set -x|echo DEBUG` → none. `git status --porcelain` empty. The RED mutation was reverted with `git checkout --` and the file confirmed byte-identical to its committed state **before** the gate run.

## Context

Builder-discovered on REQ-246 (Discovered Tasks, first item), classified [normal]. Gated behind REQ-255 so the repairer's shape handling settles once, on top of the parity sweep, rather than twice.

## Open Questions

- [x] I discovered this out-of-scope task while working on REQ-246: the repairer refuses offset/fractional stamps that the board still warns about — implementing offset arithmetic in `comparison_key_for` would close the last board-detectable-but-unrepaired timestamp class. Should I process this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it — the documented refusal (D-04) is the permanent answer and the board's warning is disclosure enough.

**Answered [2026-08-18]:** User approved via `do-work clarify` — queued for a future work run.

---

## Triage

**Route: B** - Medium

**Reasoning:** The outcome is clear (decide, then either implement offset/fractional arithmetic or make the refusal permanent and provable) but grounding that decision needs discovery: how `comparison_key_for` and `extract_timestamp_fields` currently recognize shapes after REQ-255's parity sweep, and what the read-side detectors accept.

**Planning:** Not required

---

## Exploration

`skills/do-work/scripts/repair-req-timestamps.sh` holds the whole shape surface in two functions the REQ names:

- `comparison_key_for` (line ~206) — turns a value token into a sortable key; ends by gating on `calendar_components_valid` (line ~179). This is the single place a shape is recognized or refused, and REQ-247's auditor sources this file, so widening it widens the auditor in the same edit.
- `extract_timestamp_fields` (line ~259) — whitespace-token extraction; REQ-255 taught it the quoted and unquoted space-separated spellings.
- The script header (line ~51) states the current refusal: *a numeric UTC offset or fractional seconds is not provably wrong without timezone arithmetic* (REQ-246 D-04).

Lock-ins live in `_dev/tests/prescribed-shell-scripts-behavior.sh` as a `# repair-req-timestamps:` comment-headed case group (lines ~1098–1500), including the two REQ-255 space-separated cases and the skew-constant lock-step case at ~1474 — the pattern any new case follows.

*Explored inline by the orchestrator*

## Scope

**Files I will touch:**
- `skills/do-work/scripts/repair-req-timestamps.sh` (modify) — the decision's implementation: either offset/fractional recognition in `comparison_key_for`, or the refusal made permanent and provable
- `_dev/tests/prescribed-shell-scripts-behavior.sh` (modify) — lock-ins for whichever answer wins

**Files I will NOT touch:** `skills/do-work/scripts/audit-archive-timestamps.sh` (it sources the repairer, so it inherits the change without an edit), the SessionStart hook, the board tool, `CHANGELOG.md`, `VERSION`, `skills/do-work/actions/version.md`.

**Acceptance criteria (restated from REQ):**
- [ ] The offset/fractional residual is decided either way, and the decision is stated where a reader of the script meets it
- [ ] Whichever way it goes, a lock-in pins it — a refusal that is only a comment is not pinned
- [ ] The one code body stays shared: nothing is duplicated into the auditor
- [ ] `bash _dev/tests/maintainer-verify.sh` exits 0

---

## Implementation Summary

**Files changed:**
- `skills/do-work/scripts/repair-req-timestamps.sh` (modified)
- `_dev/tests/prescribed-shell-scripts-behavior.sh` (modified)

**What was done:** The offset/fractional refusal was **kept and made permanent**, and no arithmetic was added — `comparison_key_for` is byte-identical to its pre-REQ state. What changed is the reason and the enforcement. REQ-246 D-04 justified the refusal as *undecidable without timezone arithmetic*; that is wrong, since an RFC3339 offset denotes an exact instant. The header now states the real reason — the arithmetic is the **risk**, not the obstacle — with the concrete hazard spelled out: `2026-08-19T00:29:11+05:00` denotes `2026-08-18T19:29:11Z`, so a repairer that reads the wall clock and ignores the offset sees a value five hours later than the instant and erases a correct stamp as future-dated, unattended, from a SessionStart hook. Refusing can only fail to fix; repairing can destroy. The suite gained two case blocks — a six-shape refusal case and a read-side layout lock-step case that fails if the board's `parseTimestamp` layouts change underneath this decision — plus a numeric-offset fixture in the archive-parity block, taking the named case count from 64 to 66. `audit-archive-timestamps.sh` needed no edit; it sources the repairer and inherits the behaviour.

---

## Discovered Tasks

Transcribed by the orchestrator from `do-work/runs/work-2026-08-18-211613/REQ-257-handback.md` (a worktree builder cannot write this file — REQ-270).

- **[normal] The board's future warning prints a value that is not in the file.** For `created_at: 2093-01-01T00:00:00+02:00` the warning reads `created_at 2092-12-31T22:00:00Z`; for `2093-01-01T00:00:00.500Z` and `2093-01-01 00:00:00.5` it reads `2093-01-01T00:00:00Z`. The YAML layer types these as timestamps and `normalizeFrontmatterValue` re-formats them to RFC3339 UTC before the warning is built. A user who greps the file for the value the board named will not find it, and the warning's own remediation — "rewrite with the current UTC instant" — is aimed at a string that is not there. Observed by execution, not fixed: it lives in the board package, outside this REQ's write set. **This matters more now than before**, because this REQ's decision makes the board's warning the *only* signal for these shapes.
- **[low] `created_at: 2093-01-01 00:00:00Z` is repaired here but is not board-parseable** — a space separator with a `Z` matches no `parseTimestamp` layout, and the board emitted no warning for it. The repairer is *broader* than the read side in this one direction, which is benign (an unreadable value is replaced with a readable one). Recorded as a known asymmetry rather than proposed as a REQ.

---

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh` against the merged tree (range `e5e0232..6afcbd5`), un-piped with the exit code read directly
**Result:** ✓ `GATE_EXIT=0` — "Maintainer verification passed." The suite reports **66 named script cases**, up from 64: this REQ's two new case blocks are present and passing in the integrated tree. This run is both Step 6.5's testing and Step 8's post-merge verification.

**Red-green validation:** a refusal cannot go red on its own, so RED is a mutation — the builder widened `comparison_key_for` exactly as "someone silently starts repairing those shapes" would, and ran the suite against pre-change code. Observed, verbatim:

```
FAIL: repair-req-timestamps refused-shape case logged a correction for an offset or fractional stamp
FAIL: repair-req-timestamps refused-shape case rewrote REQ-817-offset-ahead.md — the offset/fractional
      refusal is a decided permanent answer; changing it means re-deciding it, not widening comparison_key_for
… (six fixtures, one per board-parseable refused shape)
```

Suite exit 1. A second RED, after the archive fixture was added, fired the auditor's assertions too — proving the pin reaches the **sourced** body and not just the repairer:

```
FAIL: audit-archive-timestamps refusal-parity case repaired a numeric-offset stamp in the archive —
      the offset refusal is the sourced library one and must reach every tool built on it
```

GREEN is the 66-case suite passing in the merged tree above.

**Premise check — the REQ's own framing was tested, and it was wrong.** Ten shapes through the shipped repairer: only the canonical one is repaired. The same ten through a board binary built to scratch: **six** produce a future warning. Lowercase `z`, `+0200`, and space-before-offset produce none — the board cannot parse those either, so those three refusals are agreed by both sides. The residual is therefore exactly six shapes, and the six lock-in fixtures are those six — the case pins the residual class rather than a sample of it.

**Class check, not instance check.** `comparison_key_for` was sourced out of the shipped script and fuzzed across the value space — 2 dates × 2 times × 6 fraction spellings × 10 zone spellings × 2 separators × 3 quote spellings:

```
fuzz: 48 accepted, 1392 refused, 1440 total
```

Zero leaks. The 48 accepted are exactly the fraction-less, offset-less-or-`Z` shapes. No offset or fractional value produces a comparison key, and none is half-recognized — which is the property REQ-255's never-half-rewrite rule demands.

**New tests added:** two case blocks — the six-shape refusal, and a read-side layout lock-step case that fails if the board's `parseTimestamp` layouts change underneath this decision — plus a numeric-offset fixture in the archive-parity block whose "logged a correction" assertion now covers both refused shapes. 64 → 66 named cases.

**Existing tests updated (cross-REQ impact):** the archive-parity block gained a fixture and a widened assertion; no prior REQ's asserted behaviour changed.

*Verified by work action*

---

## Review

**Overall: 90%** | 2026-08-18T22:58:27Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 82% |
| Test Adequacy | 80% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Important findings (each with its recorded gate disposition — this is the durable audit record the gate mandates):**
- The read-side layout lock-step is blind to most spellings it exists to catch: its `sed` captures only `time.RFC3339` or a `"2006…"`-prefixed literal, so adding `time.RFC3339Nano`, `time.DateTime`, `time.RFC1123` or a non-`2006` literal to `parseTimestamp` leaves the suite green. `time.RFC3339Nano` is precisely the layout a fractional-support widening would add, so the guard is blind in its own headline scenario. — gate: rule-change → REQ-271 created
- Same line: `\|` is GNU BRE alternation, so on BSD/macOS `sed` the extraction returns empty and the case fails spuriously — **the whole maintainer gate goes red on macOS**. Only `\|`-in-sed construct in the `_dev/tests/` tree. — gate: rule-change → REQ-271 (same root cause, same line, same fix)
- Stale restatement inside a declared and touched file: `_dev/tests/prescribed-shell-scripts-behavior.sh:1150` still justifies the offset refusal as one the repairer "cannot compare without timezone arithmetic" — verbatim the reason this REQ's header repudiates. — gate: trivial → appended to REQ-267 as an instance (same file, already in its write set, same family as its instance 2)
- Four shipped citations point `do-work forensics` **Check 11** at the future-dated-timestamp check, which is Check 12; three of the four are in the header this REQ rewrote, and the decision's disclosure argument rests on that pointer being followable. Pre-existing. — gate: trivial → REQ-272 created (as a sweep, because four sites found by one reviewer reading one diff is a sample)

**Minor findings:** 3 (report only)
- The header's headline reason under-covers half the class it justifies: "normalizing an offset means carrying civil-date arithmetic" is offset-only, while fractional truncation is monotone and can never manufacture a false future or a false ordering defect. The decision still stands on the population argument, but this is the same "right answer, wrong reason" pattern the builder correctly diagnosed in REQ-246 D-04, reproduced at half scale.
- The hand-back's "exactly the six board-parseable shapes" is overstated — it is six of a ten-shape sample; the reviewer found five more board-parseable, future-warned, refused spellings. No behavioural consequence, since `comparison_key_for` is an allowlist and anything unlisted is refused by construction.
- The REQ's `## What` still carried the premise the build disproved; corrected in place above rather than left standing.

**Acceptance:** Pass — the reviewer re-derived everything rather than accepting it. The hazard reproduced from scratch: with the naive widening, a stamp whose true instant was two hours in the *past* was silently rewritten two hours later than the truth, exit 0, unattended. The refusal side proven inert in the opposite direction (both defect passes skip a field with an empty comparison key), which is the whole asymmetry the decision rests on. The premise re-derived against a self-built board binary — past offset/fractional stamps produce zero warnings, their 2093 twins warn. RED-by-mutation reproduced independently (exit 1, six refusal assertions plus both auditor-parity assertions). Auditor inheritance proven by execution in both directions: with the repairer mutated, the same auditor invocation repaired the offset stamp.
**Suggested testing:** 3 items
**Follow-ups created:** REQ-271, REQ-272; **sweeps appended to:** REQ-267

*Reviewed by review-work action*

---

## Lessons Learned

**What worked:** Testing the REQ's own premise before deciding on it. The REQ asserted the board warns on offset and fractional shapes; ten fixtures through a board binary showed it warns on six, and only when the instant is future. That correction shrank the residual from "these shapes" to "these shapes *and* fabricated", which is most of the case for refusing — a decision made on the REQ's stated premise would have been argued on a false population. Second: fuzzing the recognizer instead of listing cases. 1,440 shapes through `comparison_key_for`, zero leaks, and the 48 accepted are exactly the intended set.

**What didn't:** The secondary lock-in — the one guarding against the board's layouts changing underneath the decision — enumerates the spellings it already knew about instead of keying on the condition, so it is green against `time.RFC3339Nano`, the very layout a fractional-support widening would add. The same line uses GNU-only `sed` alternation, which would turn the whole gate red on macOS. A guard that cannot fail is worse than no guard, because it is read as coverage; both are REQ-271.

**Worth knowing:** REQ-246 D-04 was **right for the wrong reason**, and the wrong reason survived a build, a review, and a follow-up REQ before anyone checked it. "We can't decide this without timezone arithmetic" invites the next reader to go build the arithmetic; "we can, and we must not" closes the question. When a refusal is inherited rather than re-derived, read its stated reason as sceptically as its answer — REQ-267 carries a refusal justified the same way and is queued behind this. Also: this REQ's own header now states a reason that covers offsets but not fractional seconds, so the pattern is not fully escaped yet.

## Orientation

The timestamp repairer's refusal of offset and fractional stamps is now a settled, argued, and tested boundary rather than an open gap — a reader of the script meets the reason where the refusal lives, and the archive auditor inherits it through the shared sourced body with no edit of its own; lives in the do-work core's unattended SessionStart repair path, governed by `_dev/primes/prime-shell-commands.md`. No map change — no executable line changed anywhere; `comparison_key_for` is byte-identical to its pre-REQ state.

