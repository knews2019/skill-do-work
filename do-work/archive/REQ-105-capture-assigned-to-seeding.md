---
id: REQ-105
title: Implement capture-side seeding of `assigned_to` — the schema's advertised behavior doesn't exist
status: completed
created_at: 2026-08-05T09:43:47Z
claimed_at: 2026-08-05T10:34:42Z
completed_at: 2026-08-05T10:39:00Z
route: A
user_request: UR-019
domain: general
prime_files: []
tdd: false
suggested_spec:
depends_on: []
maintenance: false
write_set: [actions/capture.md, actions/capture-reference.md]
related: [REQ-097, REQ-106, REQ-107]
batch: sync-review-0174
---

# Implement Capture-Side Seeding of `assigned_to`

## What

`actions/work-reference.md`'s `assigned_to` schema line (Request File Schema — Full Frontmatter) says the field is "Seeded by capture when the user earmarks work" — but `actions/capture.md` and `actions/capture-reference.md` contain zero mentions of `assigned_to`. An agent following the capture action can never seed it, so the advertised behavior doesn't exist. Add earmark detection and seeding instructions to capture.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** One passive earmark-assessment bullet in `actions/capture.md` Step 1 (seed only when the user names a session, verbatim, never invented) + one optional commented `assigned_to` line in `actions/capture-reference.md`'s Simple REQ template, matching the existing external-condition comment idiom. No work-reference edit: its schema-line claim becomes true as-written (verified by re-reading line 130).
- [x] **[APPLY]:** Both edits landed, scope limited to the two declared files.
- [x] **[UNIFY]:** `git diff --stat` shows exactly `actions/capture.md` (+1) and `actions/capture-reference.md` (+1) beyond do-work/ bookkeeping. `bash _dev/tests/contract-regressions.sh` passes. No debug artifacts — both changes are prose/template lines. Verified: capture.md (bullet sits above External-condition assessment, cross-refs resolve), capture-reference.md (comment line sits in the frontmatter template block, style matches neighboring comments).

## Detailed Requirements

- `actions/capture.md` Step 1 gains an earmark assessment: when the request text names a session the work is reserved for (e.g. "leave this for cloud-alpha", "this one's for the laptop checkout"), seed `assigned_to: "<session>"` in the REQ frontmatter — YAML-quoted, verbatim (no normalization; the field is verbatim-read class per `actions/work-reference.md`'s Schema Read Contract). Never invent a session name — seed only when the user names one.
- `actions/capture-reference.md` templates gain the optional `assigned_to` frontmatter line with a comment matching the schema line's semantics (advisory, not a lock; cleared on explicit claim).
- Keep the lock-step: `actions/work-reference.md`'s schema line and `tools/queue-kanban/model.go`'s parser comments must still agree with what capture now does. The schema line's claim becomes true as-written, so it likely needs no edit — verify rather than assume.

## Context

Found by a downstream consumer's review of the 0.170.1 → 0.174.3 sync; verified here at triage — `grep -rn assigned_to actions/capture.md actions/capture-reference.md` returns nothing while `actions/work-reference.md:130` makes the seeding claim. Triage correction to the reviewer's scope: the 0.172.0 CHANGELOG entry does **not** carry the "seeded by capture" claim (no changelog rewrite needed); the claim lives only in the schema line. The field itself shipped in REQ-097 / ADR-018 (`decisions/records/adr-018-regrain-session-ownership-to-claim-anywhere-one-releaser.md`).

## Builder Guidance

Certainty: Firm on direction — the reviewer offered "drop the claim" as the alternative, but seeding at capture is the natural home for the earmark ("leave this one, it's mine" is said at capture time; the only other path is hand-editing frontmatter) and the schema line already advertises it. Keep the addition small: one Step 1 assessment bullet, one template line, no new interactive question — earmark detection is passive (seed when named, otherwise omit).

## Red-Green Proof
**RED prompt/case:** `grep -c 'assigned_to' actions/capture.md actions/capture-reference.md` → 0 matches in both files, while `actions/work-reference.md`'s schema line claims "Seeded by capture when the user earmarks work". A capture of "add retry logic to the sync client — leave this for cloud-alpha" produces a REQ with no `assigned_to`.
**Why RED now:** Capture has no instruction to detect an earmark or emit the field, so the documented behavior is unreachable.
**GREEN when:** Capture's Step 1 documents earmark detection, the capture-reference template carries the optional `assigned_to` line, and the same example capture would emit `assigned_to: "cloud-alpha"` — with the schema line still true as-written.
**Validation:** Inferred during capture (triage-verified against the repo; direction is the reviewer's stated preference)

## Full Context
See `do-work/user-requests/UR-019/input.md` for complete verbatim input.

## Triage

**Route: A (Direct to Builder)** — the REQ names the exact files (`actions/capture.md`, `actions/capture-reference.md`) and the change is well-specified: one Step 1 assessment bullet plus one optional template line. No exploration or planning needed; the triage-time file list matches the captured `write_set`.

## Implementation Summary

**What was done:** Added earmark detection to capture so the schema's "Seeded by capture when the user earmarks work" claim is now implementable: a Step 1 **Earmark assessment** bullet in `actions/capture.md` (seed `assigned_to: "<session>"` verbatim and YAML-quoted only when the user names a session; never invent; advisory semantics restated with pointers to the schema), and an optional commented `assigned_to` template line in `actions/capture-reference.md`'s Simple REQ frontmatter block.

**Files changed:**
- `actions/capture.md` (modified) — Step 1 gains the Earmark assessment bullet, placed above the External-condition assessment
- `actions/capture-reference.md` (modified) — Simple REQ template gains the optional `# assigned_to:` comment line after `maintenance:`

**Lock-step verification:** `actions/work-reference.md:130` ("Seeded by capture when the user earmarks work") is now true as-written — no edit needed. `tools/queue-kanban/model.go`'s display-only parse is unaffected (no semantic change to the field). CHANGELOG carries no seeding claim (triage-corrected scope), so no changelog rewrite.

**Tests:** `bash _dev/tests/contract-regressions.sh` passes.

## Qualification

Passed — 2 files verified (`tools/checks/qualify.sh` OK; both changes substantive instruction lines, not placeholders), all 3 requirements traced (Step 1 bullet → capture.md; template line → capture-reference.md; lock-step → verified no-edit-needed), P-A-U confirmed against the diff.

## Testing

- `bash _dev/tests/contract-regressions.sh` — passes (word budgets, CLAUDE.md-citation guard, rationalization-table ratchet all green).
- **Red-green validation:** RED (captured): `grep -c 'assigned_to' actions/capture.md actions/capture-reference.md` → 0 matches in both. Confirmed RED before the change (triage grep, exit 1). GREEN: same grep now returns 1 match per file — the Step 1 earmark-assessment bullet and the template comment line. The example capture ("leave this for cloud-alpha") is now instructed to emit `assigned_to: "cloud-alpha"`, with the schema line true as-written.

## Orientation

Now capture can actually earmark work for a named session: the `assigned_to` seeding the schema always advertised is an instruction the capture action carries, not just a claim. Lives in the capture action (do-work's intake surface); no map change.

## Review

**Acceptance: Pass** (Route A quick scan, calibrated depth). Requirements: all three detailed requirements delivered; the third resolved as verify-no-edit, which the Implementation Summary documents. Code review: both lines match their neighbors' idiom (assessment-bullet prose; frontmatter comment style), cross-references point at real sections (`actions/work-reference.md` → Schema Read Contract / Request File Schema), and the never-invent guard is stated in the imperative where the seeding instruction lives. Restatement sweep: the only other statements of capture-seeding are `actions/work-reference.md:130` (now true) and ADR-018 (silent on seeding) — no stale restatements. Scope: touched files exactly match `write_set`. No findings.

---
*Source: downstream sync-review finding 1, verified at triage — see UR-019*
