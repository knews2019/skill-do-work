---
id: REQ-105
title: Implement capture-side seeding of `assigned_to` — the schema's advertised behavior doesn't exist
status: pending
created_at: 2026-08-05T09:43:47Z
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
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

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

---
*Source: downstream sync-review finding 1, verified at triage — see UR-019*
