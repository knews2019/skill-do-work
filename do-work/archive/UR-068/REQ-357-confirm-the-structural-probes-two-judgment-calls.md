---
id: REQ-357
title: "Finish REQ-343's structural probe: caution wording, no carve-outs, supported shapes, board strip"
status: completed
claimed_at: 2026-08-24T13:55:00Z
status_changed_at: 2026-08-24T14:16:10Z
completed_at: 2026-08-24T14:16:10Z
commit: 5c95f8f
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
write_set: [skills/do-work-board/tools/queue-kanban/verify.go, skills/do-work-board/tools/queue-kanban/verify_test.go, skills/do-work-board/tools/queue-kanban/generate.go, skills/do-work-board/tools/queue-kanban/generate_test.go, do-work/archive/legacy/REQ-060-failed-req-resolution-path.md]
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
- [x] **[PLAN]:** Verify uses the parsed ticket plus retained frontmatter, so keep its lenient parser behavior and make the structural checks schema-driven: distinguish an absent opening fence from a present-but-unclosed fence; retain empty-id detection with caution-level language; remove the legacy-directory predicate; and exempt only affirmative documented shapes (stakeholder, code-review metadata, or `context_ref`). Add focused Go tests that mutate each discriminator and a payload test proving unrecognized statuses remain verifier findings but are omitted from the board strip. The legacy migration note names 11 files, but the tree contains 12 UR-less legacy files: REQ-001 through REQ-011 are valid `source: code-review` plus `review_generated: true` and `scope:` shapes, so no UR can be truthfully invented; REQ-060 has an evidenced source relationship to REQ-059 under UR-008 and will be backfilled. Render the fixture to validate the actual page surfaces.
- [x] **[APPLY]:** Implemented the planned verifier, renderer, regression-test, and evidenced data backfill changes only.
- [x] **[UNIFY]:** Reviewed every changed source and data file, ran `gofmt`, `git diff --check`, focused and full Go tests, the rendered-board fixture, and independent review; no debug artifacts found.

## Instances

- [x] `verify.go` empty-`id` finding: keep the detection, downgrade the wording so it reads as a
  caution rather than damage — the board recovers the id from the filename, so the real hazard is a
  rename silently renumbering the REQ, and the wording should say exactly that severity. Overturns
  D-07. (found by REQ-357 clarify / UR-068)
- [x] `verify.go` legacy carve-out: remove `isLegacyArchiveRequestPath` and its call site, and
  backfill `user_request` into the 11 REQs under `do-work/archive/legacy/` so the check needs no
  special case. Overturns D-08. (found by REQ-357 clarify / UR-068)
- [x] `verify.go` fence finding wording: a file with an intact opening fence and no closing one must
  not be described as "has no leading frontmatter fence" — distinguish the two shapes in the detail
  text, keep the one-finding-per-file rule (D-03) and leniency, and pin the missing-closing-fence
  case in `verify_test.go`. Folded from REQ-358 (review finding I1). (found by REQ-343 review / UR-068)
- [x] `verify.go` supported UR-less shapes: a REQ minted by `code-review.md`'s template
  (`source: code-review` / `review_generated: true` with `scope:`) and a legacy `context_ref` REQ
  archived to `do-work/archive/` root are documented shapes and must produce no `user_request`
  finding — key each on its schema discriminators, state why damage cannot cheaply wear them, and
  pin both plus the carve-out-removal mutations in `verify_test.go`. A genuinely missing field with
  neither shape is still flagged. Folded from REQ-363 (PR #166 review), reconciled with the D-08
  reversal: recognition by discriminator, never by path, and backfill wherever a real UR exists.
  (found by PR #166 review / UR-068)
- [x] `generate.go` board strip: add `unrecognized-req-status` to `boardRenderedVerifyCategories`
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

## Exploration

The verifier keeps the raw malformed file in `BodyMarkdown` whenever frontmatter splitting fails, so missing-closing-fence wording can inspect that retained structured evidence without a filesystem re-walk or model change. The board payload already has one suppression map and an existing generated-data test. The 11 original legacy review records have the documented code-review discriminators; REQ-060 instead traces to archived REQ-059 under UR-008.

## Scope

**Files I will touch:**

- `skills/do-work-board/tools/queue-kanban/verify.go`
- `skills/do-work-board/tools/queue-kanban/verify_test.go`
- `skills/do-work-board/tools/queue-kanban/generate.go`
- `skills/do-work-board/tools/queue-kanban/generate_test.go`
- `do-work/archive/legacy/REQ-060-failed-req-resolution-path.md`

**Also update:** frontmatter `write_set` to mirror the planned files.

**Out of scope:** parser/model changes, path-based exemptions, changes to the eleven valid code-review records, canceled sibling REQs, and any change to verifier exit semantics.

## Implementation Summary

- `skills/do-work-board/tools/queue-kanban/verify.go` (modified): distinguishes missing opening and closing fences from retained parser evidence, labels an empty id as a filename-recovery caution, and replaces the legacy-directory exemption with documented frontmatter discriminators.
- `skills/do-work-board/tools/queue-kanban/verify_test.go` (modified): pins the two fence shapes, caution wording, the supported UR-less schemas, removal mutations, and the removed path exemption.
- `skills/do-work-board/tools/queue-kanban/generate.go` (modified): suppresses unrecognized-status verifier prose from the board finding strip, where its warning banner, blocked-column placement, and invalid badge already surface it.
- `skills/do-work-board/tools/queue-kanban/generate_test.go` (modified): proves the raw verifier still finds a typo status, the generated payload omits that duplicate, and structural damage remains present.
- `do-work/archive/legacy/REQ-060-failed-req-resolution-path.md` (modified): backfills its evidenced `user_request: UR-008`; REQ-001–011 remain valid code-review-shaped UR-less records.

## Testing

- RED: the original focused tests failed against the prior empty-id wording, legacy-directory exemption, and missing-closing-fence classification.
- GREEN: `GOTOOLCHAIN=go1.26.1 go test -count=1 -run 'TestVerify(FlagsEachStructuralDamageShape|StaysSilentOnLegitimateAbsenceOfUserRequest|StaysSilentOnEveryDocumentedUserRequestlessSchema|UserRequestExemptionsRequireTheirSchemaDiscriminator|DistinguishesMissingClosingFrontmatterFence)|TestGeneratedBoardDataCarriesVerifyFindingsWithoutTheOnesTheBoardAlreadyShows' ./...` (PASS).
- Full package: `GOTOOLCHAIN=go1.26.1 go test -count=1 ./...` (PASS, 41.064s).
- Rendered fixture: `http://127.0.0.1:8093/` showed `REQ-900` as `pnding` with INVALID badge and a data-warning banner, while `#board-findings` held only `REQ-901`'s structural finding; this confirms unrecognized status is suppressed only from the duplicate strip.

## Qualification

- `bash skills/do-work/tools/checks/qualify.sh do-work/working/REQ-357-confirm-the-structural-probes-two-judgment-calls.md` — PASS.
- `bash skills/do-work/tools/checks/scope-drift.sh do-work/working/REQ-357-confirm-the-structural-probes-two-judgment-calls.md` — PASS; declared and modified files match.
- `GOTOOLCHAIN=go1.26.1 bash _dev/tests/maintainer-verify.sh` — PASS; canonical strict-browser lane skipped because no browser binary was configured, while the targeted in-app-browser render check above passed.

## Review

Approved after one P2 correction: the missing-closing helper now exactly mirrors `splitFrontmatter`'s newline-terminated opening-fence contract, with a bare EOF `---` regression test. Review confirmed schema-driven exemptions, no path carve-out, truthful REQ-060-only backfill, and correct payload suppression.

## Decisions

- The requested legacy “11” is the count of valid UR-less code-review templates (REQ-001–011), not a set of links to fabricate. Their affirmative schema shape is now recognized. REQ-060 is the only additional legacy record with an evidenced source UR, so it alone receives `UR-008`.
- Missing-closing classification uses retained `BodyMarkdown`; it does not re-walk the filesystem and therefore preserves the verifier's structured-evidence seam.

---
*Source: builder-decided questions from REQ-343 (UR-068); instances folded from REQ-358, REQ-363 (PR #166 review), and REQ-359.*
