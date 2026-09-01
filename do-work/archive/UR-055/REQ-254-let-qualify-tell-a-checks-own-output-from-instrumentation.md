---
id: REQ-254
title: Let qualify tell a check's own output from leftover instrumentation
status: completed
created_at: 2026-08-18T14:04:16Z
claimed_at: 2026-08-18T19:12:47Z
completed_at: 2026-08-18T19:53:41Z
commit: 116eec6
kb_status: promoted
kb_entry: REQ-254-let-qualify-tell-a-check-s-own-output-fr.md
route: B
status_changed_at: 2026-08-18T14:04:16Z
user_request: UR-055
addendum_to: REQ-244
domain: general
review_generated: true
effort_estimate: trivial
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: true
write_set:
- skills/do-work/tools/checks/qualify.sh
- _dev/tests/prescribed-shell-scripts-behavior.sh
estimate:
  p50_active_minutes: 5
  confidence: high
  calculated_at: 2026-08-18T19:13:28Z
  basis:
    - trivial short-circuit (effort_estimate: trivial)
---

# Let Qualify Tell a Check's Own Output From Leftover Instrumentation

## What

`qualify.sh`'s debug-artifact scan FAILs on any added `print(` line, which makes it fire on a **check's own success output** — the reporting a checker is supposed to have. It fired exactly that way on REQ-244's remediation and had to be overridden on the record.

## Context

REQ-244's review found that the hand-back's GREEN transcript came from a prototype rather than the shipped checker, because the shipped checker **printed nothing on success and so had no transcript to quote**. The fix was to give it one. That fix then tripped `qualify.sh`:

```
FAIL: [UNIFY] is checked but the diff adds debug artifacts — un-check it and flag:
  155:+print(
```

The flagged line is the check's only success line, placed after its failure raise, and `maintainer-verify.sh` prints it on every run. So the gate that exists to catch stray instrumentation fired on the fix for a defect caused by *missing* instrumentation.

A false positive here is not free. It teaches a reader that a qualify FAIL can be waved through, which is the opposite of what the gate is for — and the next override may be a real artifact.

## Requirements

- A checker's own reporting does not trip the debug-artifact scan, while genuine leftover instrumentation still does.
- The distinguishing rule is **stated as a condition** and not as a list of allowed files or allowed line shapes (CLAUDE.md → Closed Enumerations Go Stale; the detail is in `_dev/primes/prime-shell-commands.md`).
- `maintenance: true`: **ask what can be removed before adding.** It is worth asking whether a bare-`print(` grep earns its place at all, given that the scan cannot see the difference between a debug line and a contract's output, and given that the reviewer and the orchestrator both read the diff anyway. Deleting a heuristic that cries wolf may be a better answer than teaching it a new trick.

## Implementation Summary

**What was done:** `qualify.sh`'s debug-artifact scan (Check 4) now splits its tokens by property. Unfinished-work markers (`debugger`, `TODO`, `FIXME` — vocabulary illustrative) FAIL anywhere, unchanged. Output primitives (`print(`, `console.log` — the class, not the fired token) are judged by process-exit ownership: a file that ends its own process (exit idioms, illustrative vocabulary) has a terminal audience, so an added output line is presumptively its own reporting and surfaces as a legible WARN; a file that never ends its process is library code, so the same line FAILs naming the file and reason. The output half walks changed files per path in both serial and range modes. REQ-244's ready-made case now passes with a WARN; genuine library instrumentation still FAILs. Deletion of the heuristic was weighed and declined on the REQ's own GREEN criterion (D-01). Orchestrator applied one seam: `review-work.md`'s diff-hygiene line no longer re-flags a checker's success output.

**Files changed:**
- `skills/do-work/tools/checks/qualify.sh` (modified) — token split (`unfinished_marker_regex` / `output_primitive_regex`), per-file ownership walk, condition stated in the comment
- `_dev/tests/prescribed-shell-scripts-behavior.sh` (modified) — three `qualify:` lock-ins (reporter passes serial+range; library instrumentation FAILs with named file; TODO in a reporter still FAILs), suite 52 → 55 cases
- `skills/do-work/actions/review-work.md` (modified) — integration seam: diff-hygiene wording aligned with the new boundary

*Integrated by orchestrator from builder hand-back; merge range `8f564b3..116eec6`.*

## Decisions

Transcribed from the builder hand-back:

- **D-01 (DECIDE & STATE):** keep a narrowed scan, decline deletion — the REQ's GREEN requires library instrumentation to keep FAILing; the narrowing removes exactly the crying-wolf half and the scan now claims strictly less.
- **D-02 (DECIDE & STATE):** the condition is process-exit ownership — printed output belongs to whoever owns the process exit; both token and exit-idiom vocabularies marked illustrative.
- **D-03 (DECIDE & STATE):** the exemption covers the output-primitive class, not the fired token — `console.log` in a Node checker is the same defect one language over; pinned in the reporter fixture.
- **D-04 (DECIDE & STATE):** unfinished-work markers stay unconditional — the reporter exemption is not a file-level pardon; pinned by the third lock-in.
- **D-05 (DECIDE & STATE):** reporter hits WARN, not silent pass — WARN is the script's existing judgment channel; accepted limit: a forgotten debug print inside a checker WARNs, the honest boundary of intent-blindness.
- **D-06 (DECIDE & STATE):** no `work.md` seam — the FAIL/WARN prose contract is unchanged; only which side one token class lands on moved.

## Qualification

Passed — 3 files in merge range `8f564b3..116eec6` (2 builder + 1 seam), requirements traced (condition-not-list verified in the shipped comment; both REQ RED cases exercised; deletion genuinely weighed with the tie-breaker named), P-A-U audited. The builder's pushback — the REQ's condition-and-deletion asks pull opposite ways and the Red-Green Proof broke the tie — is accepted as the correct reading; the widest consequence (any script's stray print WARNs rather than FAILs) is recorded with its rationale.

## Builder Guidance

Candidate signals, none prescribed: whether the added line is inside `_dev/tests/`; whether the file is itself a check; whether the print is unreachable after a raise/exit on the failure path; whether the printed text is asserted anywhere. The mechanical definition is the builder's call — the requirement is that it keys on a property, not on a filename.

## Red-Green Proof

**RED prompt/case:** run `qualify.sh` over a diff that adds a success-reporting `print(` to a `_dev/tests/` checker, and over a diff that adds a genuine debug `print(` to implementation code.
**Why RED now:** both FAIL identically today; REQ-244's remediation range is a ready-made case.
**GREEN when:** the first passes, the second still FAILs, and the reason each gets is legible.
**Validation:** Orchestrator override recorded in REQ-244's Review Remediation section.

---

## Triage

**Route: B** - Medium

**Reasoning:** The defect and a ready-made red-green case are stated; the open design question — a condition that separates a check's own output from instrumentation, or deleting the heuristic — is the builder's, inside one script plus its lock-ins.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Scope

**Files I will touch:**
- `skills/do-work/tools/checks/qualify.sh` (modify) — the debug-artifact scan keys on a condition, or the bare-print heuristic is removed (maintenance latitude)
- `_dev/tests/prescribed-shell-scripts-behavior.sh` (modify) — red-green lock-ins: checker success line passes, genuine debug print still FAILs

**Files I will NOT touch:**
- `skills/do-work/actions/work.md` — Step 6.3's prose already describes FAIL/WARN semantics; change it only via integration seam if the resolution alters the contract.

**Acceptance criteria (restated from REQ):**
- [ ] A checker's own reporting does not trip the scan; genuine leftover instrumentation still does.
- [ ] The distinguishing rule is stated as a condition, never a file/line-shape list.
- [ ] The delete-the-heuristic answer was genuinely considered (maintenance: true).
- [ ] `bash _dev/tests/maintainer-verify.sh` exits 0.

## Pre-Flight

**Git:** ✓ clean
**Tests baseline:** ✓ `bash _dev/tests/maintainer-verify.sh` exits 0 at the branch point (0.212.14 tip)
**Dependencies:** ✓ Go 1.26.1, ShellCheck 0.11.0, `just`, Node, Chromium present

*Checked by work action*

## AI Execution State (P-A-U Loop)

Added by the orchestrator at integration — this review-generated REQ predates the block, which is itself review finding I2. Transcribed from the builder hand-back:

- [x] **[PLAN]:** Read brief, crew rules incl. maintenance/testing, the prime, qualify.sh, suite structure, work.md Step 6.3, and REQ-244's actual remediation diff. Candidate signals evaluated and rejected for cause; settled on process-exit ownership.
- [x] **[APPLY]:** TDD order — lock-ins first, RED observed against unmodified qualify.sh, then the split. Two write-set files, +108/−2.
- [x] **[UNIFY]:** Diff reviewed; shellcheck warning-level clean; new print( occurrences are fixture bytes and the suite owns its exit; maintainer-verify exit 0; tree clean.

## Orchestrator correction and armed-run override

Two records set straight, per review finding I2:

1. **The mid-integration commit's "P-A-U transcription" claim was false.** The transcription replace calls targeted the capture template's checkbox text, which this file never contained; they no-op'd silently and the commit message overclaimed. The block above is the real transcription. The same slip affects REQ-250's trail commit message (its archived file also lacks the section); recorded here and in the session checkpoint rather than by editing history.
2. **Armed-run override, on the record (REQ-244 precedent).** With a P-A-U section present, shipped `qualify.sh` over this REQ's own range `8f564b3..116eec6` FAILs twice: the marker scan on the lock-in's TODO *fixture* lines and on qualify.sh's own regex/comment text, and the output walk on `review-work.md`'s seam line ("never ends its own process" — a category error for prose). All three are contract/fixture/documentation bytes, not unfinished work or instrumentation — exactly the protected class. Overridden: the diff was read, each flagged line is quoted here by its role, and the rule fix (WARN on missing section; prose false positive) is REQ-264/REQ-263's territory.

## Review

**Overall: 91%** | 2026-08-18T19:49:39Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 85% |
| Test Adequacy | 90% |
| Scope | 90% |
| Risk | Low |
| Acceptance | Pass |

**Verdict: Approve** — the ownership condition survives adversarial execution at its boundary: both RED cases reproduce GREEN in serial and range modes, all three lock-ins are mutation-falsifiable, and every hole found is either the REQ's accepted limit (D-05) or a pre-existing class this change strictly narrowed.

### Requirements Checklist

- [x] **Checker's own reporting passes; instrumentation still FAILs** — *reproduced by execution* in both modes: reporter WARN + exit 0 naming the reason; library debug print FAIL + exit 1 naming file and reason; browser `console.log` with no exit FAILs (D-03's class claim); TODO in an exit-owning checker still FAILs (D-04, verified live).
- [x] **Condition, never a list** — the shipped comment states the ownership condition and marks both vocabularies illustrative; the exit-idiom enum degrades safely (a missing language's idiom falls back to the OLD behavior — FAIL, overridable — never a silent pass; probed with a Lua-style CLI).
- [x] **Deletion genuinely considered** — D-01's tie-breaker is correct: the REQ's own GREEN requires library instrumentation to keep FAILing; the scan now claims strictly less and nothing that passed before now FAILs.
- [x] **`maintainer-verify` exit 0** — observed un-piped (55-case suite included).

### Findings

**Important (audit record with gates):**
- **I1 — the ownership condition is satisfiable by non-semantic bytes** (same-diff `sys.exit(0)` smuggle; `__main__`-guarded exit in a dual-use module; a docstring that merely *says* "exit 1") — each flips FAIL to WARN; the implemented whole-file grep is weaker than the stated condition. *Reproduced by execution.* Mitigation is structural (WARN names the file and instructs confirm-from-diff; Step 6.3 defines WARN as judgment; D-05 concedes intent-blindness on the WARN side), and the categorical fix would break the REQ's own GREEN case. — gate: **trivial** → REQ-263 (pending-answers, generation ≥2), folded with the WARN-legibility Minor.
- **I2 — this REQ's own qualification "Passed" with Check 4's FAIL half disarmed**: the file carried no P-A-U section, both UNIFY-gated FAIL branches key on a checked box, and an armed run over the range FAILs on protected-class false positives with no override recorded. Not unique to this builder — every review-generated REQ from the previous session lacks the section. — gate: **rule-change** → REQ-264 (pending-answers, generation ≥2); the missing transcription and the armed-run override are corrected in this file, above.

**Minor:** 5 (report only) — WARN omits matched lines (folds into REQ-263); `work.md:750`'s Because-cell now conditionally stale ("a diff containing console.log is a false claim the qualifier will catch" — in an exit-owning file it deliberately WARNs; a reader acting on it only over-complies); the prose-file false positive persists with a confidently wrong reason (docs quoting the tokens FAIL as "never ends its own process" — pre-existing in kind, class strictly narrowed); mechanical scope-drift flag on the orchestrator-applied seam landing in `review-work.md` where the Scope predicted `work.md` (traceability intact, not builder drift); seam sentence style dense (cosmetic).

**Nit:** WARN fires even with `[UNIFY]` unchecked while FAIL requires it checked — harmless asymmetry.

### Acceptance Testing

**Result: Pass** — both RED cases GREEN in both modes in scratch fixtures; boundary probes executed (dual-use, smuggle, docstring, browser JS, comment/string literals, moved lines, shell `echo` unscanned, stale-enum edge); all three lock-ins proven mutation-falsifiable by executing three mutations (restored byte-identical); suite count 52 → 55 verified by the suite's own runtime count; gate exit 0 with a clean tree at review end.

### Suggested Additional Testing

Pin the accepted WARN boundary itself (docstring-"exit 1" case) so a future tightening shows red-green; exercise the missing-P-A-U WARN against REQ-250's and REQ-254's shapes when REQ-264 lands; confirm Go's `fmt.Println`/`os.Exit(` sitting outside both vocabularies is the intended floor.

**Follow-ups created:** REQ-263, REQ-264 (by orchestrator; pending-answers per the depth stop) · **Sweeps appended to:** none

*Reviewed by review-work action (independent adversarial pass, orchestrated mode; merge range `8f564b3..116eec6`)*

## Lessons Learned

**What worked:** Grounding the condition in the real fired case (REQ-244's actual remediation diff) instead of the REQ's abstract description — that is what disqualified three plausible candidate signals before any code. Safe-degradation as a design property: an unknown language's exit idiom falls back to the old FAIL-and-override behavior, never a silent pass.

**What didn't:** The condition-as-implemented (whole-file grep for exit-idiom text) is weaker than the condition-as-stated (file ends its own process) — the ninth instance-vs-class occurrence, this time in the fix for the gate that hunts that shape. And the pipeline's own paperwork had the same hole one level up: review-generated REQs carry no P-A-U block, so the box audit was silently disarmed for exactly the REQs this session processed, and a false "transcription" claim in a commit message went unnoticed until the review re-armed the audit.

**Worth knowing:** WARN is qualify's judgment channel and FAIL its gate; REQ-254 moved one token class between them, and `work.md:750`'s Because-cell is now conditionally stale (over-compliance direction). A forgotten debug print inside any exit-owning script WARNs — the honest boundary of intent-blindness.

## Orientation

Now qualify's debug-artifact scan distinguishes a check's own reporting (WARN, judgment) from leftover instrumentation (FAIL) by process-exit ownership, and unfinished-work markers still FAIL anywhere. Lives in core's `tools/checks/` gate layer. [MAP CHANGED] — Check 4 is now two half-checks with different channels, and the review layer's hygiene wording follows it. Prime staleness spot-check: `prime-shell-commands.md` paths still resolve; not stale.
