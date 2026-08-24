---
id: REQ-358
title: "Review fix: tell a broken closing fence from a missing opening one"
status: claimed
claimed_at: 2026-08-24T13:05:00Z
created_at: 2026-08-24T10:40:00Z
user_request: UR-068
addendum_to: REQ-343
domain: testing
review_generated: true
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
route: A
impact: impact-user-visible
effort_estimate: effort-mechanical
write_set:
  - skills/do-work-board/tools/queue-kanban/verify.go
  - skills/do-work-board/tools/queue-kanban/verify_test.go
---

# Tell a Broken Closing Fence From a Missing Opening One

## What

`verify` reports a REQ whose *closing* frontmatter fence is missing as one that "has no leading
frontmatter fence". The operator is told the opposite of what they can see: the file's first line is
`---`. Distinguish the two shapes in the finding's detail text.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Why

REQ-343 exists so an operator can act on a finding "without opening the tool's source" — its own
stated requirement. A finding that names an intact fence as the broken one fails that test in the
direction that wastes the most time: the reader checks the first line, sees `---`, and concludes the
tool is wrong.

`splitFrontmatter` (`frontmatter.go:62-65`) returns `hasFrontmatter=false` for **both** shapes — no
opening `---`, and an opening fence with no closing one ("No closing fence — treat the file as having
no frontmatter"). `appendStructuralDamageFindings` keys only on `ticket.FrontmatterMarkdown == ""`
and emits one detail text for both.

Reproduced 2026-08-24 against the shipped binary, on a file whose first line is verified `---` by
`cat -A`:

```
! structurally-damaged-req: REQ-804 has no leading frontmatter fence, so id, status, user_request
  and every other field parsed empty (the id named here was recovered from the filename)
```

The remedy string happens to name both fences, so the damage is still repairable — this is a wrong
diagnosis rather than a broken one, which is why it is `impact-user-visible` and not critical.

## Detailed Requirements

- A file with an intact opening fence and no closing fence is described as such, not as one missing
  its leading fence.
- A file genuinely missing its opening fence keeps its current wording.
- The remedy text stays actionable for both shapes.
- `verify_test.go` gains a missing-closing-fence case. Its absence is why this shipped — the existing
  fixture covers only the missing-opening-fence shape.

## Constraints

- `_dev/primes/prime-kanban-board.md` governs. Read it first.
- Keep the one-finding-per-fenceless-file rule REQ-343 established (D-03): a fenceless file reports
  once, not once per emptied field. This REQ changes the *wording*, not the count.
- Do not start rejecting files. Leniency is REQ-343's standing constraint and it still holds.

## Red-Green Proof

**RED prompt/case:** A REQ file whose first line is `---` and which has no closing `---` is reported
as "has no leading frontmatter fence". Reproduced above.

**GREEN when:** that file's finding names the closing fence as the broken one; a file genuinely
missing its opening fence still reports the current wording; and both shapes are pinned in
`verify_test.go`.

**Validation:** Inferred during REQ-343's review, then independently reproduced by the orchestrator.

---
*Source: REQ-343 review finding I1 (UR-068).*
