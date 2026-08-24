---
id: REQ-343
title: "Let verify see a structurally damaged REQ file"
status: claimed
claimed_at: 2026-08-24T08:55:00Z
created_at: 2026-08-23T22:35:07Z
user_request: UR-068
domain: testing
prime_files: [_dev/primes/prime-kanban-board.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
related: [REQ-342, REQ-344]
maintenance: false
route: B
impact: impact-user-visible
effort_estimate: effort-substantive
write_set:
  - skills/do-work-board/tools/queue-kanban/verify.go
  - skills/do-work-board/tools/queue-kanban/verify_test.go
---

# Let Verify See a Structurally Damaged REQ File

## What

`queue-kanban verify` reports `OK: no findings` and exits 0 on REQ files whose structure is broken.
Give it a structural-anomaly probe, and lift the unrecognized-status warnings the board already
produces into findings the same way the duplicate-id probe does.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read `prime-kanban-board.md`. Followed `appendCompletionAnomalyFindings`' pattern — forward the board's structured evidence, do not re-walk the tree or parse warning prose. Two probe categories rather than one, because the status class is a pre-existing warning class the REQ names separately.
- [x] **[APPLY]:** Two files, both inside the write set. `model.go` would have been the cleaner home for the two fields the probe needs, but it is outside the write set — hence the `requestFrontmatterFields` helper, which re-reads bytes `buildBoard` already sliced onto the ticket, through the production parsers.
- [x] **[UNIFY]:** Audited by the orchestrator against the merged range `2143f2c..793ab23`: built the binary and ran it against a purpose-built damaged fixture (broken opening fence, typo status, healthy control) — both damage shapes named, the healthy REQ silent — and against the real tree, where it adds zero findings despite 11 `archive/legacy/` REQs carrying no `user_request`.

## Why

On the user's fixture, six of seven REQs carried delimiter damage — including one whose opening
frontmatter fence was broken so its `status`, `title` and `user_request` all parsed empty — and verify
printed `OK: no findings` and exited 0. Verify is the mechanical half of the pre-commit ritual, so a
clean exit there is what an operator trusts before committing. It is also the safety net for the
damage REQ-342 and REQ-344 exist to prevent: without this probe, that damage is silent twice over.

The mechanism is the parser's leniency, which is correct on its own terms: recovery means the damage
surfaces as *empty fields* rather than a parse error, so nothing throws. `buildBoard`'s
unrecognized-status warning is the only trace it leaves.

## Context

Verified against the source. `collectVerifyFindings` (`verify.go:141-158`) runs thirteen `append*`
probes. Exactly one of them reads `board.Warnings`: `appendDuplicateRequestIdFindings`
(`verify.go:282-292`), and it filters by `duplicateRequestIdWarningPrefix`, so every other warning
class — unrecognized status included — is never lifted. `ExitCode` (`verify.go:104`) keys on
`len(report.Findings)`, so an unlifted warning cannot affect the exit status.

`appendCompletionAnomalyFindings` is the pattern to follow rather than invent: its own comment
records that "verify was blind to every anomaly class until then" (REQ-214), and it forwards the
board's structured evidence instead of re-walking the tree or parsing warning prose.

**Capture decision — the missing-`user_request` class was narrowed, and the user did not ask for
that.** The request lists "a missing `user_request` pointer" as an anomaly the probe should fail on.
Taken literally that fires on every stakeholder-questions REQ, which carries no `user_request` **by
design** (`actions/work-reference.md` → Stakeholder REQ Template), and on every REQ in
`archive/legacy/`, which predates the field. A probe that flags correct files is a probe someone
turns off, so the requirement below carves those two out. **Value:** the probe stays trustworthy, so
its findings keep meaning something. **Risk:** a genuinely damaged REQ that happens to look like a
stakeholder REQ would slip through — narrow, because the carve-out keys on the documented shape and
not on the field's mere absence. Reversible: delete the carve-out and the literal reading is back.

## Detailed Requirements

- A structural-anomaly probe fails the mechanical check on a REQ file with any of: no leading
  frontmatter fence, an empty or unrecognized `status`, an empty `id`, or a missing `user_request`
  pointer.
- Existing unrecognized-status warnings are lifted into findings the same way
  `appendDuplicateRequestIdFindings` lifts duplicate-id warnings.
- Forward the board's structured evidence; do not parse warning prose and do not re-walk the tree —
  the same rule `appendCompletionAnomalyFindings` states for itself.
- Each finding names the broken field and its remedy, so the operator can act without opening the
  tool's source.
- The probe distinguishes damage from legitimate absence: a stakeholder-questions REQ and a REQ in
  `archive/legacy/` both legitimately carry no `user_request`, and neither is a finding. This narrows
  the request's literal wording — see the capture decision in `## Context` for why, and overturn it
  there if the narrowing is unwanted.

## Constraints

- `_dev/primes/prime-kanban-board.md` governs this tool. Read it first, including the parser
  lock-step convention.
- **Keep the parser's leniency.** The point is to report the damage, not to start rejecting files:
  a REQ with one bad line must still parse and still appear on the board.
- Do not weaken `TestMaintainerStrictBrowserBehaviorLaneRejectsZeroProbes` or any existing verify
  probe, and do not change what the board renders — this REQ adds detection, not display.

## Builder Guidance

**Certainty: firm — the gap, the mechanism and the pattern to copy are all confirmed in the source.**
The one real judgment is the legitimate-absence carve-out above; get that wrong and verify cries wolf
on every stakeholder REQ, which is how a probe gets disabled.

Write the fixture before the probe. The user's shape — several REQs damaged, at least one with a
broken opening fence — is the right RED, and it is also the fixture that proves the carve-out, so
include a stakeholder REQ and a legacy REQ in it that must NOT be flagged.

## Open Questions

None — the user named the four anomaly classes and the lifting pattern.

## Red-Green Proof

**RED prompt/case:** Build a fixture repo whose `do-work/queue/` holds REQ files with each of the four
damage shapes, plus a stakeholder REQ and a legacy REQ that are structurally fine. Run
`queue-kanban verify --repo-root <fixture>`: it prints `OK: no findings` and exits 0.

**Why RED now:** Only duplicate-id warnings are lifted into findings, and `ExitCode` keys on the
finding count, so an unrecognized or empty status cannot fail the check.

**GREEN when:** The same fixture produces one finding per damaged REQ, each naming the broken field,
verify exits nonzero, and the stakeholder and legacy REQs produce no finding. A healthy fixture still
exits 0 with `OK: no findings`.

**Validation:** User confirmed — the fixture, the counts, the mechanism and the lifting pattern are
stated verbatim in the input, and each was re-verified against the source during capture.

## Assets

None. The fixture is the deliverable's own RED.

---
*Source: UR-068 — see `do-work/user-requests/UR-068/input.md` for complete verbatim input.*

---

## Triage

**Route: B** - Medium

**Reasoning:** The defect, the mechanism and the pattern to follow were all established in the REQ's `## Context` down to line numbers. What needed discovery was where the fields live — `RequestTicket` does not carry a declared `id` or the `stakeholder:` marker — and how to reach them without re-walking the tree.

**Planning:** Not required.

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Scope

**Files I will touch:**
- `skills/do-work-board/tools/queue-kanban/verify.go` (modify) — two probe categories and their helpers
- `skills/do-work-board/tools/queue-kanban/verify_test.go` (modify) — fixture, damage-shape table, carve-out and discriminator tests

**Files I will NOT touch:** `model.go` (the cleaner home for the two fields the probe needs, but outside the write set — see D-04) and `generate.go` (the board's findings-strip suppression map — see D-06).

**Acceptance criteria (restated from REQ):**
- [x] A structural probe fails the check on: no leading fence, empty/unrecognized status, empty id, missing user_request
- [x] Unrecognized-status warnings lifted the way duplicate-id warnings are
- [x] Board's structured evidence forwarded; no prose parsing, no second tree walk
- [x] Each finding names the broken field and its remedy
- [x] Stakeholder REQs and `archive/legacy/` REQs produce no finding
- [x] Parser leniency preserved — a damaged REQ still parses and still reaches the board

## Implementation Summary

**Files changed:**
- `skills/do-work-board/tools/queue-kanban/verify.go` (modified)
- `skills/do-work-board/tools/queue-kanban/verify_test.go` (modified)

**What was done:** Added `verifyCategoryStructurallyDamagedRequest` and `verifyCategoryUnrecognizedRequestStatus`, wired into `collectVerifyFindings` after the duplicate-id probe. The structural probe reports a missing leading fence, an empty `id:`, and a missing `user_request:`; a fenceless file reports once rather than once per emptied field, since the fence remedy repairs all of them. The status probe lifts `RequestTicket.StatusUnrecognized` — the structured form of the verdict `bucketColumns` already reaches — rather than matching warning prose. Legitimate absence is carved out by the documented shape: a `stakeholder:` marker, or an `archive/legacy/` path anchored on separators the same way `isArchivedUserRequestPath` is.

## Testing

**Tests run:** `GOTOOLCHAIN=go1.26.1 go test -count=1 ./...`, `GOTOOLCHAIN=go1.26.1 QUEUE_KANBAN_BROWSER=... bash _dev/tests/maintainer-verify.sh`
**Result:** ✓ Module suite ok (122.9s); gate exit 0 with the strict browser lane actually run rather than skipped

**Red-green validation:**
- Baseline binary on the damage fixture: `OK: no findings`, exit 0 → after: five findings naming field and remedy, exit 1
- Orchestrator's independent fixture (broken opening fence + typo status + healthy control): both damage shapes named, healthy REQ silent

**Mutation evidence:** six mutations applied to `verify.go`, each caught — dropping either carve-out, letting a broken fence fall through to the per-field probes, unwiring the status probe, exempting every REQ from the `user_request` check, and ceasing to notice a missing fence.

**False-positive check (REQ-280's rule):** verify against the real tree before and after is byte-identical — zero new findings, with 11 `archive/legacy/` REQs carrying no `user_request`. Independently reproduced by the orchestrator.

**Leniency preserved:** `summary` still reports `total REQ tickets : 8` on the damaged fixture; all eight ids reach `board.RequestsById`.

*Verified by work action*

## Decisions

<!-- D-XX counter: last used D-06. Next decision: D-07. -->

- **D-01 — Two categories, not one. DECIDE.** `structurally-damaged-req` covers fence/`id`/`user_request`; `unrecognized-req-status` covers the status class, which the REQ names separately and which is a pre-existing warning class.
- **D-02 — Lift `ticket.StatusUnrecognized`, not the warning sentence. DECIDE.** The REQ asks for both "the same way `appendDuplicateRequestIdFindings` lifts" and "do not parse warning prose". The flag is the structured form of the same verdict, so lifting it satisfies the pattern without the prose match.
- **D-03 — A fenceless file reports once. DECIDE.** Its `id`, `status` and `user_request` are empty *because* the fence is gone, and the fence remedy repairs all of them — the same reasoning `appendTimestampOrderingFindings` records for its outer pair.
- **D-04 — `requestFrontmatterFields` re-parses the ticket's retained bytes rather than extending `RequestTicket`. DECIDE.** Adding `IdDeclared`/`Stakeholder` to the ticket is the cleaner home, but `model.go` is outside the write set. This reads bytes `buildBoard` already sliced and kept, through the production parsers — not a second walk and not a second parser. It deletes if those fields are ever promoted.
- **D-05 — Both clean-base fixtures gained `user_request: UR-071`. DECIDE.** Without it, "clean tree" meant a shape no captured REQ has. The assertion is unchanged; the fixture is more honest.
- **D-06 — The board's findings strip will render both new categories. ESCALATE.** `attachVerifyFindings` forwards everything not in `boardRenderedVerifyCategories`, which lives in `generate.go`, outside the write set. **Value:** the strip stays the single honest view of what fails the mechanical check, and no `generate.go` change was needed to ship detection. **Risk:** the status class now appears three times on the page (data warning, per-card invalid badge, findings strip); reversible with one line in the suppression map. Filed as a discovered task rather than fixed inline.

## Discovered Tasks

- Suppress `unrecognized-req-status` from the board's findings strip (`generate.go` → `boardRenderedVerifyCategories`) — that class already reaches the page twice, so the strip is a third copy. One line, outside this write set. `impact-negligible`, `effort-mechanical`.
- `do-work cleanup` has no pass for structural damage: every finding here is non-`Fixable`, so the remedies are hand edits. A pass that could restore a missing fence or backfill `user_request` from the containing `archive/UR-NNN/` directory would make the fixable count mean more. Needs its own capture — fence repair is not mechanically safe. `impact-user-visible`, `effort-substantive`.
- The `stakeholder:` marker is not parsed onto `RequestTicket`. Verify now depends on it while the board ignores it entirely, so a stakeholder REQ is invisible as such on the board. `model.go` territory. `impact-rule-change`, `effort-mechanical`.

## Open Questions

- [~] Where should an empty `id:` be reported? → **D-07**: Builder flagged it as damage per the requirement, and said so in the detail text. Reasoning: `deriveRequestIdFromFilename` means the board never loses the REQ over it, so the real exposure is narrower than the other three classes — a file rename silently renumbers the REQ. Value: the rename hazard is named where an operator will see it. Risk: if it reads as noise in practice, deleting the `id` branch leaves the other three classes intact. Carried to REQ-357 for the maintainer to confirm or overturn.
- [~] Should `archive/legacy/` be the carve-out's key rather than a `created_at` cutoff? → **D-08**: Builder used the directory, because that is what the REQ's `## Context` names and it needs no date arithmetic. Value: no clock dependency in a structural probe. Risk: a REQ written today and dropped into `archive/legacy/` would be exempt — narrow, and visible the moment anyone looks at the directory. Carried to REQ-357.
