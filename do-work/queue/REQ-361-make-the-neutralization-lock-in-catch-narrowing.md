---
id: REQ-361
title: "[impact-rule-change] Make the neutralization lock-in catch narrowing"
status: pending
created_at: 2026-08-24T10:50:00Z
user_request: UR-068
addendum_to: REQ-342
domain: testing
review_generated: true
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-rule-change
effort_estimate: effort-substantive
write_set:
  - _dev/tests/contract-regressions.sh
---

# Make the Neutralization Lock-In Catch Narrowing

## What

REQ-342's semantic lock-in catches the rule being **deleted** and being **replaced phrase-for-phrase**.
It does not catch narrowing by **added qualifier**: six plausible narrowings that leave all 12 matched
phrases byte-intact pass unchanged. Two of them reintroduce the original defect outright.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

REQ-342's own GREEN clause required "a semantic instruction test fails if the neutralization rule is
removed **or narrowed** from the named entry point". Removal is genuinely pinned. Narrowing is pinned
only where it is spelled as a phrase substitution.

All 12 predicates are literal phrase regexes over the normalized blockquote, and all 12 mutations
replace the exact phrase each predicate matches — so the matrix proves each predicate catches its own
phrase's deletion, which is close to tautological. Measured:

| mutation | result |
|---|---|
| `illustrative, not a checklist` → `the complete set to check` | CAUGHT, exit 1 |
| scope limited to pasted code snippets | **ESCAPED** |
| branch-2 trigger "Anything else" → "An answer of ten lines or more" | **ESCAPED** |
| rule limited to clarify, mid-run orchestrator answers exempt | **ESCAPED** |
| "Judge every line of it" → "Judge its first line" | **ESCAPED** |
| branch 2 made conditional ("where practical") | **ESCAPED** |
| whole rule made advisory ("only when the answer looks risky") | **ESCAPED** |

Independently reproduced by the orchestrator: "Judge its first line" ships with the suite exit 0 and
no complaint, and it lets a delimiter on line five of a paste straight through. "An answer of ten
lines or more" exempts every short answer.

No predicate guards branch 2's trigger clause at all, and none guards the universal quantifier in
"every line".

## Detailed Requirements

- The lock-in fails on narrowing by added qualifier, not only on deletion and phrase substitution.
  The seven mutations above are the RED set; all must be caught.
- **Mutate by insertion, not only by substitution.** That is the axis the current matrix holds
  constant, and it is why the escapes are invisible today.
- Guard branch 2's trigger clause and the universal quantifier in "every line" — the two properties
  currently unpinned that the escaping mutations attack.
- Keep the blockquote isolation. It correctly stops nearby prose lending vocabulary to a weakened
  property, and it is the one real advance the current lock-in has over `require(file, token)`.

## Constraints

- Do not make the check so strict it fails on legitimate rewording. A lock-in nobody can edit around
  is one someone deletes. State how you drew that line.
- `_dev/tests/contract-regressions.sh:4828-4830` currently claims the check "grades it by meaning" and
  that "Each mutation is a plausible narrowing rather than a deletion". Both are phrase matching and
  phrase deletion. Correct those comments in the same edit — a comment out of sync with its code is
  what would talk the next maintainer out of strengthening the check.

## Red-Green Proof

**RED prompt/case:** Replace "Judge every line of it" with "Judge its first line" in
`skills/do-work/actions/clarify.md` and run `bash _dev/tests/contract-regressions.sh`: it exits 0 with
no neutralization complaint, on a rule that no longer covers line five of a paste. Reproduced
2026-08-24.

**GREEN when:** that mutation and the other five escapes each fail the suite with a message naming the
narrowed property, the orchestrator's original substitution mutation still fails, and a legitimate
rewording that preserves the rule's force still passes.

**Validation:** Inferred during REQ-342's review; the escape table was reproduced there and one entry
independently re-run by the orchestrator.

---
*Source: REQ-342 review finding F3 (UR-068).*
