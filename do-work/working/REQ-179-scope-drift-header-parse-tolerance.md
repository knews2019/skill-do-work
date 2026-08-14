---
id: REQ-179
title: Make scope-drift.sh parse Scope headers with trailing annotations
status: claimed
claimed_at: 2026-08-14T10:27:30Z
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
