---
id: REQ-361
title: "[impact-rule-change] Make the neutralization lock-in catch narrowing"
status: claimed
claimed_at: 2026-08-24T18:57:05Z
status_changed_at: 2026-08-24T18:57:05Z
route: C
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
estimate:
  p50_active_minutes: 35
  confidence: medium
  calculated_at: 2026-08-24T18:57:05Z
  basis:
    - Route C
    - 1-file write set
    - 2 subsystems involved
    - 4 acceptance criteria
    - cross-route regression gates
    - full-suite verification
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

## Triage

**Route: C** — The combined body-neutralization and Frontmatter Quoting mutation contract must derive governed fields, scan only fenced examples, distinguish insertion narrowing from legitimate rewording, and preserve the existing isolation matrix; this requires an explicit plan before the one-file test implementation.

## Addendum (2026-08-24) — the Frontmatter Quoting checker belongs here too

REQ-344's review found the same gap on its sibling contract: `grep -rn "Frontmatter Quoting" _dev/
skills/` returns zero hits in any check, so deleting the whole contract paragraph from
`work-reference.md` breaks nothing that runs. REQ-344's own GREEN required a test that fails if the
rule is removed or a write site reverts, and it is not delivered — its D-04 records that honestly and
names the blocker: `contract-regressions.sh` was being rewritten by REQ-342 concurrently.

**That blocker has cleared.** REQ-342 landed. Both contracts now need the same kind of check, so
build them together rather than twice:

- A write-site checker that fires when a shipped file emits a governed frontmatter field as a
  double-quoted or unquoted scalar. It must **derive the governed field set from the schema block**
  (the fields whose comment cites Frontmatter Quoting, minus those naming the escaping encoder)
  rather than hardcoding names, and must scan **inside fenced blocks only** — the rule's own
  counterexample `title: "Fix: A " # B"` sits in inline prose and would otherwise trip it.
- The same insertion-mutation discipline this REQ establishes applies to it: a checker that only
  catches the rule's deletion is the defect this REQ exists to remove, repeated on the sibling.

## Plan

1. Parse the named neutralization block into branch-level semantic clauses and derive the Frontmatter Quoting field inventory from the canonical schema fence.
2. Add positive semantic properties plus scoped contradiction detectors that catch inserted narrowings without banning legitimate words such as the contract's own `only containment bytes` promise.
3. Scan copyable examples only inside Markdown fences, classifying encoder exceptions from schema annotations rather than a hard-coded field exemption.
4. Replay the seven neutralization insertion mutations and sibling quoting deletion, narrowing, unsafe-example, and future-field mutations with named non-vacuity guards.
5. Run contract regressions and the canonical maintainer gate, preserving all prior REQ-342/360 mutation failures.

**Plan validation:** Every requirement maps to a semantic clause, derived inventory, or mutation trial. One file owns the lock-in logic, so the plan stays within the captured scope.

## Exploration

- REQ-342's matrix is a positive vocabulary bag over the isolated Step 4 block. Its replacements go RED because they remove the phrase being checked; inserted contradictions leave every positive phrase present and stay GREEN.
- Bind the inline branch to both its one-line and delimiter-safety triggers, and bind the containment branch to all three alternatives: own passage, any line break, or any delimiter-shaped line. Preserve blockquote isolation and tolerate synonym/order variation through normalized clause checks.
- Derive governed fields from raw-user annotations in `work-reference.md`'s canonical YAML schema fence. The live inventory is `title`, `assigned_to`, `blocked_by`, `blocked_check`, `stakeholder`, `tested_by`, and `testing_feedback`; the test must discover rather than store that list.
- Fence-only example scanning excludes the contract's inline prose counterexample. Encoder-owned double-quoted examples are recognized by their schema discriminator/comment, never by field name or directory.
- Seven insertion trials cover named-shape checklist narrowing, REQ-answer-only reach, illustrative-to-complete wording, optional `> ` containment, short-answer-only inline handling, fixed triple-backtick fencing, and first-line-only prohibition. Sibling trials cover universal reach, citation/annotation integrity, unsafe fenced examples, and a newly annotated future field.

## Scope

**Files I will touch:**

- `_dev/tests/contract-regressions.sh`

**Acceptance criteria:**

- Semantic clause checks pin both branch triggers and universal containment without overfitting one canonical sentence.
- All seven insertion narrowings fail for named reasons while legitimate equivalent rewording and the existing replacement mutations retain their intended outcomes.
- Frontmatter Quoting derives its governed inventory from the schema, scans fenced examples only, and fails on deletion, narrowing, unsafe examples, or a newly governed field without a second field list.
- Focused contract regressions and canonical maintainer verification pass with non-vacuity guards intact.
