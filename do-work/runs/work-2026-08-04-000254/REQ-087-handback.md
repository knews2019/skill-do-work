# Hand-back — REQ-087

**Branch:** `worktree-agent-REQ-087-posix-only-timestamp-command`
**Commit:** `202ff3e`

## File manifest

- `tools/queue-kanban/verify.go` (modified) — the future-dated-`claimed_at` remedy keeps a command,
  but it is now `queue-kanban now`.
- `tools/queue-kanban/web/board.js` (modified) — the clock-skew tooltip (line ~154) and the
  future-stamp badge title (line ~553) drop to the target shape plus a citation.
- `tools/queue-kanban/model.go` (modified) — **outside the declared `write_set`; stop-and-report.**
  A fourth site of the same class: the server-side future-timestamp board warning carries the same
  POSIX-only command in a user-facing string. Same treatment applied.
- `tools/queue-kanban/future_timestamp_test.go` (modified) — **outside the declared `write_set`;
  stop-and-report.** Its assertion pins the literal `date -u +…` that `model.go` emits, so the
  wording change breaks it. Repointed, intent preserved.

## Integration seams

**None.** All four files are inside `tools/queue-kanban/`; nothing shared with the sibling builder.

## Per-surface decisions (requirement 1)

The three named surfaces were **not** given matching wording, per the requirement:

- **`verify.go` — a command survives.** This is CLI output read next to a shell, so dropping to a
  file reference makes it worse. But the command is `queue-kanban now`, which is the Timestamp rule's
  own **option 1**, not its POSIX floor. Anyone reading a `queue-kanban verify` report has the binary
  built by definition, so option 1's "only if already built" precondition is satisfied at the point of
  use, and the string is platform-independent — which is the actual defect being fixed.
- **`board.js` ×2 and `model.go` — shape plus citation.** A tooltip's job is to explain the badge, not
  to be a manual. All three now say "the current UTC instant — YYYY-MM-DDTHH:MM:SSZ, per the Timestamp
  rule in actions/work-reference.md".

Requirement 2 honored: the Windows one-liner is pasted **nowhere**. Requirement 3 is satisfied
vacuously for the two-branch command (it appears in no display string) and directly for the one
surviving command, which is the rule's own first-choice source.

## Requirement 4 — no new contract assertion, and why

A ratchet here would have to separate a **user-facing display string** from a **rationale comment or
test fixture**, and no grep can see that difference. After this change the tool's remaining `date -u`
mentions are `timestamp.go:13,21` (rationale prose explaining why the subcommand exists),
`future_timestamp_test.go:103` (a fixture), and `verify.go:295` (a comment explaining why the POSIX
floor is deliberately *not* spelled there) — every one correct as it stands, and every one a
file-scoped assertion would flag. Scoping the assertion to a hand-listed file set instead just moves
the failure: the list goes stale the day a new display surface is added, which is the closed-enumeration
pattern `CLAUDE.md` names explicitly. The honest control here is the per-surface judgment requirement 1
asked for, not a ratchet that cannot see the distinction it would need to make.

## Notes for the owner

- **The REQ's site inventory was a floor, again.** It named three sites; a fourth (`model.go:321`) is
  the same class and same failure. Found by re-running the grep across `tools/` rather than trusting
  the list — the pattern this session's checkpoint flagged twice already.
- **Cross-REQ test impact, one test, intentional.** `TestFutureTimestampWarningNamesFieldAndFix`
  asserted the literal `date -u +%Y-%m-%dT%H:%M:%SZ` appears in the board warning. That is exactly the
  behavior this REQ changes. The assertion's intent — the warning names the field *and* the fix — is
  preserved: it now pins the target shape and the rule citation instead. Flagged rather than quietly
  edited.
- No version bump, no `CHANGELOG.md` entry — the integrator's, per `CLAUDE.md`'s scope clause.
- `gofmt` clean, `go vet` clean, `go test ./...` green on this branch.
