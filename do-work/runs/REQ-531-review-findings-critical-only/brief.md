# REQ-531 Builder Brief

Implement the contract reversal: review and builder findings below `impact-critical` stay in the report/current archived REQ and never create, append to, or otherwise mutate a follow-up REQ, sweep, pending-answers item, or prose backlog. Only `impact-critical` findings auto-queue.

Requirements:

- In review output, every noncritical finding line ends with `→ report only`.
- In builder discovery handling, noncritical discoveries remain recorded in the current REQ; remove the test-hygiene auto-queue carve-out.
- `Follow-ups created` lists critical findings only. If all N findings are noncritical, write `None (N findings report only)`.
- Tell maintainers how to promote a report-only finding manually with `do-work capture`, quoting the finding line as the capture source.
- Update the Fold-First destination ladder so destination 4 is report-only; no automatic REQ/prose-backlog mutation below critical.
- Preserve unrelated failure classification, builder-decided follow-ups, and stakeholder-requested follow-ups.
- Replace/delete stale contract-regression pins rather than merely adding more; `_dev/tests/contract-regressions.sh` must not grow in line count.

Authorized project write set:

- `skills/do-work/actions/review-work.md`
- `skills/do-work/actions/work-reference.md`
- `skills/do-work/actions/work.md`
- `skills/do-work/actions/capture-reference.md`
- `skills/do-work/actions/capture.md`
- `skills/do-work/crew-members/general.md`
- `skills/do-work-toolbox/crew-members/general.md`
- `skills/do-work/docs/review-work-guide.md`
- `skills/do-work/docs/work-guide.md`
- `skills/do-work/docs/standing-preferences.md`
- `skills/do-work-toolbox/actions/code-review.md`
- `skills/do-work/next-steps.md`
- `_dev/tests/contract-regressions.sh`

Do not modify anything under `do-work/` in the builder worktree. Do not change `VERSION`, changelogs, or release metadata; the orchestrator owns release finalization.

Before editing, read `CLAUDE.md`, `_dev/primes/prime-action-files.md`, `_dev/primes/prime-shell-commands.md`, `_dev/primes/lessons-action-files.md`, `skills/do-work/crew-members/general.md`, `skills/do-work/crew-members/coding-guardrails.md`, `skills/do-work/crew-members/communication-style.md`, and `skills/do-work/crew-members/maintenance.md` completely.

TDD/ratchet proof: first make a named `_dev/tests/contract-regressions.sh` contract check fail against the current text (or demonstrate a named obsolete check fails once its stale surface is deleted), record the exact RED command/output, then implement and record GREEN. Run focused contract regressions and any relevant lint. Keep the contract test file at or below its original 8,479 lines.

Commit all project changes on branch `worktree-agent-REQ-531-review-findings-critical-only`. Then write a handback to `/Users/t2/Desktop/e1-experimental-repos/skill-do-work2/do-work/runs/REQ-531-review-findings-critical-only/handback.md` containing: commit hash; changed files; requirements trace; exact RED/GREEN evidence; tests; restatement-sweep results; remaining risks. The handback is the sole permitted builder write into the main tree.
