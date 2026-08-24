---
id: REQ-357
title: "Finish REQ-343's structural probe: caution wording, no carve-outs, supported shapes, board strip"
status: pending
status_changed_at: 2026-08-24T13:35:36Z
created_at: 2026-08-24T09:50:00Z
user_request: UR-068
addendum_to: REQ-343
builder_decided: true
domain: testing
review_generated: true
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
sweep: true
sweep_key: req-343-structural-probe-remediation
impact: impact-user-visible
effort_estimate: effort-substantive
---

# Finish REQ-343's Structural Probe

## What

One sweep closing everything REQ-343's review and clarify left open on the structural probe: the
maintainer's two overturned builder decisions (D-07, D-08), plus the three sibling review fixes
folded in from REQ-358, REQ-363, and REQ-359 (all cancelled-with-pointer on 2026-08-24; their full
requirements, red-green proofs, and constraints live in their archived files).

The maintainer's governing rationale, recorded at clarify: the probe exists to keep the repo's
health visible — proper warnings are what let AI coders find and fix issues cheaply, and git makes
every resulting change trackable. So: warnings carry honest severity, data is backfilled to conform
rather than exempted by path, and documented UR-less shapes are recognized by their schema
discriminators, never by where a file happens to sit.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Instances

- [ ] `verify.go` empty-`id` finding: keep the detection, downgrade the wording so it reads as a
  caution rather than damage — the board recovers the id from the filename, so the real hazard is a
  rename silently renumbering the REQ, and the wording should say exactly that severity. Overturns
  D-07. (found by REQ-357 clarify / UR-068)
- [ ] `verify.go` legacy carve-out: remove `isLegacyArchiveRequestPath` and its call site, and
  backfill `user_request` into the 11 REQs under `do-work/archive/legacy/` so the check needs no
  special case. Overturns D-08. (found by REQ-357 clarify / UR-068)
- [ ] `verify.go` fence finding wording: a file with an intact opening fence and no closing one must
  not be described as "has no leading frontmatter fence" — distinguish the two shapes in the detail
  text, keep the one-finding-per-file rule (D-03) and leniency, and pin the missing-closing-fence
  case in `verify_test.go`. Folded from REQ-358 (review finding I1). (found by REQ-343 review / UR-068)
- [ ] `verify.go` supported UR-less shapes: a REQ minted by `code-review.md`'s template
  (`source: code-review` / `review_generated: true` with `scope:`) and a legacy `context_ref` REQ
  archived to `do-work/archive/` root are documented shapes and must produce no `user_request`
  finding — key each on its schema discriminators, state why damage cannot cheaply wear them, and
  pin both plus the carve-out-removal mutations in `verify_test.go`. A genuinely missing field with
  neither shape is still flagged. Folded from REQ-363 (PR #166 review), reconciled with the D-08
  reversal: recognition by discriminator, never by path, and backfill wherever a real UR exists.
  (found by PR #166 review / UR-068)
- [ ] `generate.go` board strip: add `unrecognized-req-status` to `boardRenderedVerifyCategories`
  and amend the map's doc comment in the same edit — the class already reaches the page three ways.
  `structurally-damaged-req` keeps forwarding; `verify`'s own exit status and finding list are
  unchanged. Confirm by rendering, not by reading the code path. Folded from REQ-359 (review finding
  I2). (found by REQ-343 review / UR-068)

## Why

Neither original decision was a defect — REQ-343 is complete and green — but a builder's judgment
call nobody confirms becomes a rule by default. Both were escalated, both were overturned. The three
folded siblings are the remaining review findings against the same probe; one builder in one run
closes all of it in the same two Go files plus the suppression map.

## Open Questions

- [x] Should verify report an empty `id:` field at all? → Keep the detection but downgrade the wording to a caution, not damage.
  Recommended (overturned): Keep it as damage — the rename hazard is real and the detail text already scopes the claim.
  Also: Drop the `id` branch (the other three classes are unaffected); or keep it but downgrade the
  wording so it reads as a caution rather than damage.
  *2026-08-24 — answered by the maintainer via clarify. Reasoning: the goal of the structural probe is
  repo health that AI coders can act on — proper warnings are what make fixes easy, and git versioning
  keeps every resulting change transparent and trackable. Since the board recovers the id from the
  filename, an empty `id:` is a real hazard (rename silently renumbers) but not damage; calling it a
  caution keeps the signal honest without overstating severity. Out of scope: dropping the branch —
  the detection stays.*

- [x] Should the `archive/legacy/` carve-out key on the directory, or on a `created_at` cutoff? → Neither: drop the carve-out entirely and backfill `user_request` into the 11 legacy REQs.
  Recommended (overturned): Keep the directory — a structural probe with a clock dependency is worse, and the
  exemption is visible the moment anyone looks at the directory.
  Also: Add a `created_at` cutoff alongside the directory.
  *2026-08-24 — answered by the maintainer via clarify. Reasoning: same repo-health goal — an
  exemption hides a gap the probe exists to surface; backfilling makes every file conform so the
  check needs no special case, and the backfill itself is one transparent, git-tracked change. Out of
  scope: both keying variants — no carve-out survives, so the directory-vs-cutoff question is moot.
  Implementation note: remove `isLegacyArchiveRequestPath` and its call site, backfill `user_request`
  into the legacy REQs.*

## Context

Both original decisions are recorded as D-07 and D-08 in the archived REQ-343, with the builder's
Value and Risk lines. The code they describe is live: `appendStructuralDamageFindings` and
`isLegacyArchiveRequestPath` in `skills/do-work-board/tools/queue-kanban/verify.go`.

Converted to a sweep on 2026-08-24 at the maintainer's instruction (fold trivial siblings into few
well-defined REQs): absorbed REQ-358, REQ-363, and REQ-359, each cancelled with a pointer here.
(The `builder_decided: true` marker was also added at answer time — the creating template should
have stamped it and had not.)

---
*Source: builder-decided questions from REQ-343 (UR-068); instances folded from REQ-358, REQ-363 (PR #166 review), and REQ-359.*
