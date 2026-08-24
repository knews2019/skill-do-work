# REQ-360 builder brief

Worktree: `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2-worktrees/worktree-agent-REQ-360-close-the-neutralization-contracts-reach`

Planning/exploration is accepted. Implement test-first in exactly the following 13 files:

- `_dev/tests/contract-regressions.sh`
- `skills/do-work-board/tools/queue-kanban/frontmatter_test.go`
- `skills/do-work/actions/{clarify.md,verify-requests.md,capture.md,capture-reference.md,stakeholder-answers.md,abandon.md,work-reference.md,review-work.md,sample-archived-req.md}`
- `skills/do-work/docs/{capture-guide.md,work-guide.md}`

Decision D-01 is fixed: refuse/report C0/DEL controls other than LF/TAB before writing; do not
normalize CR or invent a hand-authored escape table. Preserve byte-identical content apart from
containment bytes. Keep one canonical condition, make all body writers actively inherit it, correct
all five frontmatter examples, and make `|-`/`|`/`|+` preserve zero/one/multiple terminal LF bytes.

Start with semantic/mutation and strict-parser RED cases. Each of the five `## Instances` checklist
items must have a closure and falsifiable lock-in. Preserve state transitions, strict-parser target,
board encoder exception, nested estimates, answer plus dated reasoning, and prompt-injection guard.

Run focused contract/parser tests and the canonical maintainer gate. Commit on the worktree branch
and write handback only to
`/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/work-2026-08-24-174255/REQ-360-handback.md`.
