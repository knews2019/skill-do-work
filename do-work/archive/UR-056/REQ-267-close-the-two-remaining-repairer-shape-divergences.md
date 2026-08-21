---
id: REQ-267
title: Close the two remaining repairer shape divergences
status: completed
completed_at: 2026-08-18T23:39:58Z
commit: f7441d7
claimed_at: 2026-08-18T22:59:48Z
route: B
created_at: 2026-08-18T21:03:15Z
status_changed_at: 2026-08-18T22:20:09Z
user_request: UR-056
addendum_to: REQ-255
domain: general
review_generated: true
sweep: true
sweep_key: repairer-detector-shape-parity
effort_estimate: trivial
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
write_set:
- skills/do-work/scripts/repair-req-timestamps.sh
- _dev/tests/prescribed-shell-scripts-behavior.sh
estimate:
  p50_active_minutes: 25
  confidence: medium
  calculated_at: 2026-08-18T22:59:48Z
  basis:
    - Route B
    - 2-file write set
    - 4 acceptance criteria
    - cross-route regression gates
    - full-suite verification
---

# Close the Two Remaining Repairer Shape Divergences

## What

REQ-255 closed six shape divergences between the timestamp repairer and the board's readers. Its independent review, fuzzing the whole value space, found two more — both pre-existing at REQ-255's branch point, neither among its declared instances, and neither corrupting.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read the four always-load crew files plus `testing.md` (`tdd: true`), `_dev/primes/prime-shell-commands.md` with its REQ-255 and REQ-257 lessons, all 707 lines of the repairer, the board's `frontmatter.go`, `parseTimestamp`/`coerceScalarToString` in `model.go`, the existing repairer cases in the suite, and **the SessionStart hook's actual invocation of the script** — which is where the REQ's severity framing turned out to be wrong. Approach settled before editing: buffer-and-gate the extractor on the closing fence rather than special-casing the wedge (one refusal closes both faces); trim after unquoting rather than documenting a new refusal, since the header already claims the quoted family is repairable; document the residual the byte-oriented trim leaves; lock in through both scan scopes; correct the repudiated comment.
- [x] **[APPLY]:** Two files, both in the declared write set. 120 insertions, 8 deletions across two commits. No adjacent code improved, no unrelated formatting touched.
- [x] **[UNIFY]:** `git diff --stat 2fc89fe HEAD` → 2 files, +120/−8; full diff read line by line, every changed line traced to an instance in the REQ. `shellcheck --severity=warning` (the gate's threshold) exit 0 on both files; plain `shellcheck` reports five SC2004 *style* notes on untouched lines 448-481, pre-existing and below threshold. `bash -n` on the repairer exit 0. Debug-artifact scan over added lines (`TODO|FIXME|XXX|set -x|echo "DEBUG|console.log`): none. `git status --porcelain --untracked-files=all` empty — the Go probes live in the scratchpad, and **no `go build` was ever run without `-o`**.

## Instances

- [x] **An unterminated frontmatter block is repaired where the board sees only body text (reproduced by execution).** A file whose opening fence never closes is scanned to end-of-file by the extractor and its stamp lines are rewritten by the unattended hook, while the board's `splitFrontmatter` returns *no frontmatter* for exactly that shape. The script's own scope comment states the fence-bounded contract the code does not honour here. **Second face of the same root cause:** when such a file also ends with the defective stamp on its final line with no trailing newline, the changed-line guard expects four diff lines and sees two, so the repair is refused and the repairer exits 1 and prints a FAILED line into the session banner **every session, permanently**, with no self-heal. Refusing the fence-broken shape the way the read side does closes both at once.
- [x] **Quoted stamps with padding inside the quotes are board-parseable but refused, and the refusal is undocumented (reproduced by execution).** A Go probe replicating the board's pipeline (YAML unquote, then trim, then parse) accepts `"2093-01-01 00:00:00 "` and would flag it future; the repairer refuses it byte-identical and the refusal falls into the header's catch-all "anything else unparseable", which is false for this shape. The header's parity rule claims the opposite family-wide, so the documentation is wrong rather than merely silent.

- [x] **The clean-fixture comment still justifies the offset refusal with the reason REQ-257 repudiated (added from REQ-257's review, Important 3).** `_dev/tests/prescribed-shell-scripts-behavior.sh:1150` says the repairer must not touch "a numeric-offset value **it cannot compare without timezone arithmetic**" — verbatim the justification REQ-257's header now states is wrong ("the arithmetic is the risk, not the obstacle"). Same family as instance 2: a statement about what the repairer refuses that is wrong rather than merely silent. This file is already in this REQ's write set.

## Requirements

- Each instance is repaired to canonical form or refused byte-identical **with the refusal documented beside the existing refusal entries** — the never-half-rewrite rule from REQ-255 still governs.
- The permanent hook-failure loop is gone: no input shape may leave the SessionStart hook exiting nonzero every session with no path to self-heal.
- Lock-in cases per instance, through both scan scopes.
- `bash _dev/tests/maintainer-verify.sh` exits 0.

## Context

REQ-255's independent review, findings I-1 and I-2 (gate: trivial each — never-corrupt holds, and both shapes are already malformed or exotic). Created `pending-answers` per the generation-≥2 cascade stop. The review also noted a consequence for a queued sibling: REQ-257's description called offset and fractional seconds "the last" board-detectable-but-unrepaired class, which instance 2 makes inaccurate — corrected in that REQ directly.

## Open Questions

- [x] REQ-255's review found two more shapes where the repairer and the board disagree — one repairs a file the board reads as having no frontmatter (and, in one variant, makes the session hook fail forever), the other refuses a shape the board accepts while the docs claim it is handled. Should I process this as a new task? → Confirmed: Yes, add to queue
  Recommended: Yes, add to queue (will flip to 'pending').
  Also: No, discard it.

**Answered [2026-08-18]:** User approved via `do-work clarify`, presented as the only known live defect in the queue that can wedge a session permanently. Nothing was put out of scope — both instances stand, and the never-half-rewrite rule from REQ-255 still governs the refusal path.

---

## Implementation Summary

**Files changed:**
- `skills/do-work/scripts/repair-req-timestamps.sh` (modified)
- `_dev/tests/prescribed-shell-scripts-behavior.sh` (modified)

**What was done:** Both instances closed **at the primitive rather than per symptom.** `extract_timestamp_fields` now buffers its rows and emits them only once the closing `---` fence is seen, so a file whose opening fence never closes yields nothing — refusing it exactly as the board's `splitFrontmatter` does, and taking out both faces of instance 1 at once including the run that exited 1 forever. `comparison_key_for` trims ASCII whitespace after unquoting, because the read side unquotes *and then* trims (`coerceScalarToString` and `parseTimestamp` both run `strings.TrimSpace`), so `created_at: "2093-01-01 00:00:00 "` is a value the board parses and future-badges and is now repaired rather than refused. The header stopped lying in both places: the quoted-family claim now says "with any ASCII whitespace padding inside them trimmed", and two new refusal bullets were added — the unterminated-fence refusal with the permanent-failure history recorded, and **a residual the fix itself creates**: a value padded with non-ASCII whitespace (U+00A0 and friends) is trimmed read-side by Go but not here, because this file matches bytes under `LC_ALL=C`. The suite went 66 → 69 named cases, and the stale comment REQ-257's review flagged at line 1150 — justifying the offset refusal with the reason REQ-257 repudiated — was corrected in the same commit.

---

## Discovered Tasks

Transcribed by the orchestrator from `do-work/runs/work-2026-08-18-230100/REQ-267-handback.md` (a worktree builder cannot write this file — REQ-270).

- **[normal] `record-commit-hash.sh`'s read helpers scan to EOF on an unterminated fence, so `--verify` can read a body `commit:` as the frontmatter one.** Helpers at `:108-122`, called at `:247` and `:318-338`. Their own header says that confusion is structurally impossible; it is not. The script's *writer* already refuses the fence-broken shape — the readers do not, so the guard is asymmetric. **This matters more than its severity suggests:** `record-commit-hash.sh --verify` is the last check every REQ in this pipeline passes through, and it is the one that exists because free-form frontmatter edits once truncated six archived REQs to zero bytes. Out of this REQ's write set.

---

## Testing

**Tests run:** `bash _dev/tests/maintainer-verify.sh` against the merged tree (range `05c4630..f7441d7`), un-piped with the exit code read directly
**Result:** ✓ `GATE_EXIT=0` — "Maintainer verification passed." The suite reports **69 named script cases**, up from 66: this REQ's three lock-ins are present and passing in the integrated tree. This run is both Step 6.5's testing and Step 8's post-merge verification.

**Red-green validation** — RED run against pre-change code at `2fc89fe`, not a prototype:

- **Instance 1, face one** — a *body* line rewritten under an unterminated fence: the pre-change script reported `repaired … created_at: 2093-01-01T00:00:00Z -> 2026-08-10T12:00:00Z (file mtime)` and exited 0, while the real board binary answered `has no frontmatter block` (exit 2) on the very same file. GREEN: exits 0, file byte-identical.
- **Instance 1, face two — the permanent-failure loop, demonstrated over three consecutive runs.** Each printed `FAILED to repair … the rewrite touched 2 diff lines; expected 4 — the edit was rejected and discarded`, exited 1, and left the file byte-identical with no self-heal. GREEN: exits 0 on repeated runs, file byte-identical.
- **Instance 2** — the board's `frontmatter get` returned the padded quoted stamp and flagged it, while the repairer refused it. GREEN: repairs to `created_at: 2026-08-10T12:00:00Z`, with the audit line reporting the full padded old value.

**Fuzz, because an instance list is a sample.** Two harnesses against real oracles — the board binary's `frontmatter get` for visibility, and a Go probe replicating `parseTimestamp`'s four layouts and the 2-minute skew for the verdict:

- **Structural** (fence / BOM / CRLF / EOF shapes): 204 shapes, **34 divergences before, 0 after**. All 34 were one family — a valid opening fence with no exact `---` close, in all four spellings of "not closed" crossed with BOM/CRLF/EOF. No family beyond instance 1.
- **Value space** (quote style × leading pad × trailing pad × 8 value shapes): 408 shapes, **120 divergences before, 0 after**. The documented refusals — offset, fractional, calendar-impossible, non-ASCII padding — are asserted as *refused* rather than skipped, so the fuzz cannot pass by ignoring them.

Both fuzzes were run against the pre-change script recovered with `git show HEAD~1:` to prove they go honestly red, then against the committed tree to prove they go green.

**The methodological finding worth keeping:** the structural fuzz initially **missed the wedge entirely**, because the wedge *refuses* rather than rewrites — a parity fuzz that compares only file mutation is blind to a shape that fails. Adding the repairer's exit status as a second oracle dimension is what surfaced it.

**Mutation testing of the new lock-ins:** with the fixed script swapped for the pre-change one, the new cases fail with 8 distinct assertions across both scan scopes — so the lock-ins pin behaviour rather than merely accompanying it.

**New tests added:** three named cases (66 → 69), covering both faces of instance 1 and instance 2, through both the repairer's and the archive auditor's scan scopes.

**Existing tests updated (cross-REQ impact):** the stale comment at `_dev/tests/prescribed-shell-scripts-behavior.sh:1150` — which justified the offset refusal with the reason REQ-257 repudiated two hours earlier — was corrected in the same commit. No prior REQ's asserted behaviour changed.

*Verified by work action*

---

## Correction to this REQ's own framing

**The SessionStart hook does not exit nonzero.** `skills/do-work/hooks/session-start.sh:59` runs the repairer under `|| true`, deliberately — its comment says that on a tripped guard the script's failure lines *are* the audit trail and must reach the banner. `report_failure` writes to stdout, the hook captures it into `REPAIR_SUMMARY` and echoes it, and the hook exits **0**.

Verified three ways by this REQ's review: reading the hook, reading `report_failure`'s output stream, and running the real hook against a wedge fixture with the pre-change script — banner prints `do-work: FAILED to repair …`, `HOOK EXIT=0`.

The defect is real and the fix is unchanged: the **script** exits 1 on every run and prints a FAILED line into every session's start banner, permanently, with no self-heal. But this REQ was written — and approved, and prioritised into its wave — on the stronger claim, and the orchestrator repeated it to the maintainer twice before a builder checked the hook. The same false mechanism is still stated in three live docs; that is REQ-274.

---

## Review

**Overall: 96%** | 2026-08-18T23:39:40Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 96% |
| Test Adequacy | 90% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded gate disposition — this is the durable audit record the gate mandates):**
- The repudiated "the SessionStart hook exits nonzero" framing is still stated in three live maintainer docs, one of them a standing decision rationale (`CHECKPOINT.md:76`, REQ-255 D-04 — the decision is right, the argument for it names a mechanism that does not exist). — gate: rule-change → REQ-274 created
- **A third divergence family exists, on an axis neither fuzz covered: the field name.** The repairer repairs any top-level `_at` field by suffix rule; the board's `detectFutureTimestampFields` checks a fixed six-name list. Confirmed by execution — `reviewed_at` set beyond the horizon is rewritten unattended and never badged. Latent, not live: nothing in today's schema is missed and no corruption occurs, but it is the one axis where the unattended writer sits provably outside the read side's envelope. — gate: rule-change → REQ-275 created

**Minor findings:** 3 (report only)
- No positive lock-in for a **closed** fence with no trailing newline — the live path of the changed-line arithmetic. Verified correct by execution and pre-existing, so not opened by this diff.
- The new fence lock-ins pin LF only; 24 of the reviewer's 33 pre-change divergences lived in the BOM and CRLF crossings. The fix handles all of them, but a regression in the BOM strip or the `\r` sub would not be caught by these cases.
- This REQ's own record still stated the wrong mechanism, because the builder's Pushback and D-01–D-05 were never transcribed — root cause REQ-270. **Corrected in place above.**

**Acceptance:** Pass — the reviewer built its own oracle rather than reusing the builder's, copying the `queue-kanban` package to scratch and calling the **real** `splitFrontmatter` / `parseFrontmatterFields` / `coerceScalarToString` / `parseTimestamp`. Structural fuzz of 410 shapes: **33 divergences before, 0 after**, with every edge probed — no frontmatter, empty block, fence at EOF without trailing newline, CRLF, BOM ordering, `---` in body — behaving identically on both sides. Independent 704-shape value fuzz: **456 divergences before, 336 after**, and every one of the 336 falls inside a *documented* refusal family (offset 88, fractional 88, non-ASCII pad 160) — the 120 undocumented ones went to zero. Permanent-failure loop reproduced pre-change over three consecutive runs and absent after. Mutation tested **both directions**: the pre-change script fails the new cases with 8 assertions, and teaching the trim to strip U+00A0 turns the residual fixture red — so the documented residual is genuinely pinned rather than merely described. Regression: extractor output byte-identical across all 209 real `_at` rows in `queue/`, `working/` and `archive/`.
**Suggested testing:** 3 items
**Follow-ups created:** REQ-274, REQ-275, and REQ-276 from the builder's Discovered Task; **sweeps appended to:** None

*Reviewed by review-work action*

---

## Lessons Learned

**What worked:** Fixing at the primitive. Buffering the extractor's rows and emitting only on the closing fence closed both faces of instance 1 in one refusal — the body-prose rewrite and the permanent banner failure — rather than special-casing the wedge. And bidirectional mutation testing: the residual this fix *creates* (non-ASCII whitespace padding, which Go trims and `LC_ALL=C` byte matching does not) is pinned by a fixture that goes red if someone later teaches the trim to strip U+00A0. A documented limitation that can fail is worth more than one that is merely written down.

**What didn't:** The first structural fuzz missed the wedge entirely, because the wedge *refuses* rather than rewrites — a parity fuzz comparing only file mutation is blind to every shape that fails instead of corrupting. Adding the exit status as a second oracle dimension is what surfaced it. Then the same limitation recurred one axis over: both the builder's fuzz and the reviewer's held the **field name** constant, so neither could see that the repairer keys on the `_at` suffix while the board keys on six names. That is REQ-275.

**Worth knowing:** **A fuzz's blind spots are exactly the axes it holds constant, and its oracle decides which failures are visible at all.** Twice in one REQ the limiting factor was the measurement rather than the fix. When the next parity sweep is written here, vary the field name, and give the oracle both "what did it write" and "did it succeed". Separately: the hook runs this script under `|| true` on purpose, so a repairer failure is a banner line and never a broken session — a claim inherited rather than re-derived can be right in its conclusion and false in its mechanism, and this one travelled far enough to set the REQ's own approval framing.

## Orientation

A REQ file whose frontmatter fence never closes is now refused by the unattended timestamp repairer exactly as the board's reader refuses it — closing both the body-prose rewrite and the permanent banner-failure loop — and a quoted stamp padded inside its quotes repairs instead of being refused; lives in the do-work core's SessionStart repair path, shared by sourcing with the archive auditor, governed by `_dev/primes/prime-shell-commands.md`. No map change: one parser gained a gate, one recognizer gained a trim, and the refusal list gained two honest entries including a residual the fix itself creates.

