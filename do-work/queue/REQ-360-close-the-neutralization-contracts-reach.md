---
id: REQ-360
title: "[impact-rule-change] Close the neutralization contract's reach"
status: pending
created_at: 2026-08-24T10:50:00Z
user_request: UR-068
addendum_to: REQ-342
domain: security
review_generated: true
prime_files: [_dev/primes/prime-action-files.md, _dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
sweep: true
sweep_key: neutralization-contract-reach
impact: impact-rule-change
effort_estimate: effort-substantive
write_set:
  - skills/do-work/actions/clarify.md
  - skills/do-work/actions/verify-requests.md
  - _dev/tests/contract-regressions.sh
---

# Close the Neutralization Contract's Reach

## What

REQ-342 landed the neutralization rule at the Canonical answered-question format. Two ways its reach
falls short of the contract it states: the condition quantifies over **lines**, so damage that is not
line-shaped escapes it; and one live caller writes a user's answer into a REQ body without citing the
format at all, so it inherits nothing.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

One root cause: the contract was stated at the entry point without closing either end of its reach —
neither every damaging shape, nor every writer.

## Instances

- [ ] **The condition is line-scoped, so non-line-shaped damage escapes it.** The rule reads "Judge
  every line of it against one condition: could this line be read as one of this file's own
  delimiters?" A byte that changes how the **file** is read rather than how a **line** is read is
  outside it, and neither the `> ` prefix nor the over-long fence neutralizes it. Measured on an
  answer containing a NUL, written exactly per branch 2:

  | scan | control | with NUL |
  |---|---|---|
  | `grep -cE '^[[:space:]]*- \[ \]'` | 3 | 3 |
  | `grep -nE '^[[:space:]]*- \[ \]'` lines printed | 3 | **0** |
  | `grep -nE '^## '` lines printed | 3 | **0** |
  | board renders `## Plan` | yes | yes |

  grep reports "binary file matches" on stderr and prints nothing on stdout. Any reader consuming
  matched *lines* — Step 5's "any remaining unchecked item wins", a `while read` loop — sees an empty
  result on a REQ that has two real open questions, while the board shows them. That is the REQ's own
  second consequence with a different first cause. `_dev/primes/prime-shell-commands.md` §
  *Prescribed Shell Commands Must Surface What the Steps Consume* already records the trap.
  Judged `impact-rule-change`.

- [ ] **`actions/verify-requests.md:168` writes a user answer and cites nothing.** Inside Step 7's
  "the user is here right now — resolve them on the spot": `2. If the user answers → add the resolved
  question to the REQ's `## Open Questions` section as `- [x] [question] → [user's answer]``. It
  obtains typed text interactively and writes it into a REQ body in exactly the shape the format
  governs, and the string "Canonical answered-question format" appears nowhere in that file. Line 169
  is the same site's deferred branch. Judged `impact-rule-change`.

  REQ-342's builder did what its REQ asked — "grep every current citation of the format and confirm
  each inherits" — and a citation grep is structurally blind to a caller that never cites. The UR's
  wording is the stronger one: "so every caller that records an answer obeys it, not just clarify."

## Detailed Requirements

- The neutralization condition covers damage that is not line-shaped, or states explicitly and in one
  place what it does not cover and why that residue is acceptable. Either is a legitimate answer;
  silence is not.
- Every action that writes user-obtained text into a REQ body cites the format. **Find them by their
  behaviour, not by their citations** — the citation grep is what missed this one. Grep for the write
  itself: `- [x]`, `→`, "the user's answer", "resolve them on the spot".
- The lock-in gains a case for whichever instances land, so the reach is pinned rather than argued.

## Constraints

- Do not restate the rule in `verify-requests.md`. It cites; the entry point states. REQ-342's
  exactly-once count check enforces that and must keep passing.
- `_dev/primes/prime-action-files.md` governs both action files. Read it first.

## Red-Green Proof

**RED prompt/case:** A REQ carrying a NUL inside a correctly-quoted, correctly-fenced answer returns
zero lines from `grep -nE '^[[:space:]]*- \[ \]'` while the board renders its two real open questions.
Separately, `grep -c "Canonical answered-question format" skills/do-work/actions/verify-requests.md`
returns 0 while that file writes a user answer into a REQ body at line 168.

**GREEN when:** the non-line-shaped case is either neutralized or documented as out of reach with its
reason, `verify-requests.md` cites the format, and a lock-in covers both.

**Validation:** Inferred during REQ-342's review; the NUL measurement and the uncited caller were both
reproduced there.

---
*Source: REQ-342 review findings F1 and F2 (UR-068), folded into one sweep on the shared root cause.*
