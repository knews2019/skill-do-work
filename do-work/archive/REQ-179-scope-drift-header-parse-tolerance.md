---
id: REQ-179
title: Make scope-drift.sh parse Scope headers with trailing annotations
status: completed
claimed_at: 2026-08-14T10:27:30Z
completed_at: 2026-08-14T10:33:56Z
commit: e530dde
kb_status: pending
created_at: 2026-08-14T09:49:31Z
user_request: UR-040
addendum_to: REQ-178
domain: general
review_generated: true
sweep: true
sweep_key: scope-drift-header-parse-brittle
gate: rule-change
effort_estimate: normal
tdd: true
prime_files: [_dev/primes/prime-shell-commands.md]
write_set: [skills/do-work/tools/checks/scope-drift.sh]
---

# Make scope-drift.sh Parse Scope Headers With Trailing Annotations

## What

`skills/do-work/tools/checks/scope-drift.sh` anchors on the literal header `**Files I will touch:**` — both its path extraction (line 47) and its "header present but unparseable" FAIL guard (line 70). A Scope section whose header carries any annotation (REQ-178 wrote `**Files I will touch (all new, all inside …):**`) matches neither, so the check reports `SKIP` exit 2 as if the section were absent — silently disabling the exact comparison the script's own header comment says must never be silently disabled. Root cause: the rule "recognize the header, not one literal spelling of it" is enforced at zero of its two match sites.

## Instances

- [ ] `skills/do-work/tools/checks/scope-drift.sh:47` — awk path-extraction pattern requires the bare literal header; a parenthetical yields zero parsed paths
- [ ] `skills/do-work/tools/checks/scope-drift.sh:70` — the unparseable-header FAIL guard uses the same literal, so the zero-path case above degrades to SKIP instead of FAIL

## Done means the class cannot recur

The header match tolerates trailing annotations before the colon (e.g. `^\*\*Files I will touch\b[^:]*:` semantics) at BOTH sites, so an annotated header either parses or FAILs loudly — never SKIPs. A lock-in test pins the RED case: a REQ fixture with the parenthetical header must produce a real comparison (or FAIL), never SKIP.

## Red-Green Proof
**RED prompt/case:** Run `scope-drift.sh` against a REQ fixture whose Scope header is `**Files I will touch (all new):**` with one declared file and a matching Implementation Summary → currently prints `SKIP` and exits 2.
**Why RED now:** Both match sites require the bare literal; discovered live in REQ-178's review (2026-08-14).
**GREEN when:** The same fixture yields a real set-difference comparison (exit 0 on match), and a header-present-but-genuinely-unparseable fixture yields FAIL, not SKIP; lock-in test added to the suite that covers this script.
**Validation:** Inferred during review (review-generated follow-up, gate: rule-change).

## Full Context

Finding recorded in `REQ-178`'s Review section (gate token audit trail). See `do-work/user-requests/UR-040/input.md` for the parent UR.

---
*Source: review of REQ-178 — review-generated follow-up (sweep)*

---

## Triage

**Route: A** - Simple

**Reasoning:** Single named file, two identified match sites, an exact RED case, and a defined done-state — a focused mechanical fix with a lock-in test. tdd: true.

**Planning:** Not required

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*

## Decisions

- **D-01 — Lock-in test lives inline in `_dev/tests/contract-regressions.sh`.** That suite is where sibling `tools/checks/*` scripts are already exercised behaviorally with mktemp fixtures (`uncommitted-inventory.sh`, `associate-files.sh`); `prescribed-shell-scripts-behavior.sh` covers `skills/*/scripts/` only, and a dedicated file (the `record-commit-hash-guards.sh` pattern) would be ceremony for three probes. Placed immediately after the preflight assertions — before the suite's `extract_kanban_shutdown_line Justfile` call, which aborts the run in case-sensitive-filesystem sandboxes (see Discovered Tasks), so the probes execute even there. Three probes, each pinning a named failure: annotated header must produce a real comparison (exit 0 on matching sets, never SKIP), annotated-header-with-zero-parseable-paths must FAIL exit 1 (never SKIP), and a genuinely absent Scope list must keep the SKIP exit 2 contract Route A relies on.
- **D-02 — Match pattern is `\*\*Files I will touch[^:]*:` at both sites, identical by construction.** This is the REQ's suggested `[^:]*:` semantics adapted to the script's awk/grep idiom: no closing `:**` requirement (so an annotation containing a colon still parses — verified by fixture), left unanchored to match the surgical footprint of the old literal, and deliberately the SAME regex in the extraction awk and the FAIL-guard grep so no input can ever match one site but not the other — that asymmetry was exactly the SKIP degradation. The extraction `sub()` now strips through the first colon, leaving any closing `**` outside the backtick split where it is ignored. `**Files I will NOT touch:**` cannot match ("touch" is not a contiguous substring there). The FAIL message no longer quotes the bare literal since annotated spellings now reach it.

## Discovered Tasks

- `_dev/tests/contract-regressions.sh` reads `Justfile` (capital J) in `extract_kanban_shutdown_line`, but the tracked file is lowercase `justfile` — on a case-sensitive filesystem the awk open fails and `set -e` aborts the suite with exit 2 at that line, skipping every check after it (~1500 lines never run). Passes only on case-insensitive filesystems. Pre-existing; observed in the baseline run before any REQ-179 change.
- Known environmental flake confirmed flaky: `prescribed-shell-scripts-behavior.sh`'s "run-blocked-check process-tree case left the descendant alive" failed in the baseline and RED runs but passed in the GREEN run of this session — same sandbox, no related change.
- `scope-drift.sh` run against `do-work/archive/REQ-178-audit-metrics-mechanical-tool.md` now reports symmetric drift because that REQ's Scope declared bare filenames (`churn.go`) while its Implementation Summary reported full paths (`skills/do-work-toolbox/tools/audit-metrics/churn.go`). That is REQ-178's REQ-file formatting, not a tool defect — noting in case the orchestrator wants a Scope-path-spelling convention or normalization follow-up.

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/checks/scope-drift.sh` (modified — annotation-tolerant header pattern `\*\*Files I will touch[^:]*:` at both match sites; exit-code contract preserved)
- `_dev/tests/contract-regressions.sh` (modified — three inline behavioral lock-in probes, sibling-style mktemp fixtures, placed before the suite's pre-existing abort point)

**What was done:** Fixed the silent self-disable: an annotated `**Files I will touch (…):**` header now either parses (real comparison) or FAILs loudly — never SKIPs; genuinely-absent Scope still SKIPs exit 2. Lock-in probes pin the annotated-comparison, zero-paths-FAIL, and absent-SKIP contracts.

## Qualification

Passed — both files verified on disk, `bash -n` clean, requirements traced to both instance sites (the site asymmetry was the bug; identical pattern at both by construction, D-02). Orchestrator independently re-ran: REQ-178's archived REQ (the live RED case) now yields a real DRIFT comparison at exit 1, not SKIP; a Route A REQ without Scope still SKIPs exit 2.

## Testing

**Tests run:** new lock-in probes in `contract-regressions.sh` + fixture matrix + `shipped-package-reference-contract.sh` (exit 0)
**Result:** ✓ All REQ-179 probes pass; suite's remaining failures are pre-existing environmental (process-tree probe — confirmed flaky, failed RED run, passed GREEN run) and the pre-existing Justfile case-mismatch abort (→ REQ-180)

**Red-green validation:**
- Annotated-header fixture: `SKIP … exit=2` before (verbatim captured) → `OK: Implementation Summary matches the Scope declaration` exit 0 after
- New suite probes: 2 FAIL lines against the unfixed script (verbatim captured) → zero after
- 9-case matrix green: bare-literal regression, inline paths, zero-parseable→FAIL-not-SKIP, absent→SKIP, real drift→exit 1, NOT-touch exclusion, colon-in-annotation
- Finding-closure: the captured GREEN's named regression check exists and fails-before/passes-after ✓

*Verified by work action*

## Review

**Overall: 96%** | 2026-08-14 (Route A quick scan, orchestrated mode)

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 98% |
| Scope | 100% |
| Risk | None |
| Acceptance | Pass |

**Findings:** none Important. Minor (report only): REQ-178's archived Scope declared bare filenames vs full-path summary — now visible as symmetric DRIFT output when re-checked; historical formatting, archived REQs are immutable, no follow-up (a path-spelling convention would be new surface for a one-off). Discovered [normal]: the Justfile case mismatch → REQ-180 (pending-answers).

**Acceptance:** Pass — orchestrator re-ran the RED case (REQ-178 archive: real comparison, exit 1) and the SKIP-contract case (Route A REQ: SKIP exit 2) independently; TDD evidence credible and verbatim; both instance checkboxes covered by the single pattern change.

*Reviewed by work action (orchestrated quick scan)*

## Lessons Learned

**What worked:** Reproducing RED against the real defect artifact (REQ-178's archive) instead of only synthetic fixtures — it also surfaced the bare-filename Scope formatting as a bonus signal.
**What didn't:** N/A — first pass green.
**Worth knowing:** `contract-regressions.sh` aborts at line ~1797 on case-sensitive filesystems (`Justfile` vs `justfile`, REQ-180) — everything after never runs there, so a green-looking late-suite check may simply be unreached. The process-tree probe is confirmed flaky in this sandbox (failed at RED, passed at GREEN, no related change).

## Orientation

Scope-drift checking no longer silently self-disables on annotated headers: `tools/checks/scope-drift.sh` (work pipeline, Step 5.5/7 checks subsystem) now parses or fails loudly, with lock-in probes pinning the contract. Leaf fix — no map change.
