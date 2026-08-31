---
id: REQ-454
title: 'Expose UR source tokens in selection records'
status: pending
created_at: 2026-08-31T20:49:21Z
user_request: UR-085
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
related: [REQ-450, REQ-451, REQ-452, REQ-453, REQ-455, REQ-456, REQ-457]
batch: accepted-validate-feedback-root-causes
---

# Expose UR Source Tokens In Selection Records

## What

Propagate each candidate's user-request source token into selected and excluded public result records so a multi-UR targeted caller can associate every expanded member with the UR that produced it without rescanning the queue.

The fold-first scan found no pending or pending-answers REQ, sweep or otherwise, in any UR that shares this selection-result provenance-loss root cause.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Finding Provenance

- **Finding #13 — P2 — source:** `internal/nextselection/next_targets.go:85-86`

> ````text
> [P2] Expose the UR source token in selection records — [prj].claude/skills/do-work/tools/do-work-
> cli/internal/nextselection/next_targets.go:85-86
> For a run targeting multiple URs, SourceToken records which UR expanded each member but is discarded when selected and excluded
> records are built. Since actions/work.md requires announcing each UR expansion from the returned records and forbids rescanning
> the queue, the caller cannot associate ur-expanded members with their source UR; include this provenance in the result model.
> ````

- **Evidence:** Candidate provenance exists as `SourceToken` at `skills/do-work/tools/do-work-cli/internal/nextselection/next_types.go:28-33` and is populated at `next_targets.go:47-53,80-86`. It is dropped when selection records are built at `next_selection.go:220-228`, and the public selected/excluded models at `internal/resultmodel/result_model.go:84-120` expose no replacement. `skills/do-work/actions/work.md:128-135,185` requires announcements from returned records and forbids a queue rescan.
- **Surface-cost result:** N/A — direct propagation of already-computed provenance into the authoritative result.

## Detailed Requirements

- Add source-token provenance to every selected and excluded record produced from a targeted candidate.
- Preserve the exact UR token that expanded each member in multi-UR runs.
- Define the absent/default behavior for non-UR sources.
- Maintain text and JSON result parity.
- Enable action callers to make the required UR expansion announcements from result records alone.

## Constraints

- Do not make callers rescan queue files to reconstruct provenance.
- Do not conflate the source token with the request's own `user_request` field when one invocation targets multiple tokens.

## Dependencies

No request prerequisite. Shared result-model files with other UR-085 requests do not establish dependency ordering.

## Builder Guidance

Certainty level: Firm. Carry the candidate's existing provenance through all projection branches.

## Red-Green Proof

**RED prompt/case:** Target two UR tokens whose expansions each contain selected and excluded members, then consume only the text and JSON result records.
**Why RED now:** `SourceToken` exists on candidates but disappears before public records are returned.
**GREEN when:** Every returned member identifies its originating UR token consistently in text and JSON, and the caller can announce both expansions without rescanning.
**Validation:** User confirmed after validate-feedback accepted Finding #13.

## Full Context

See `do-work/user-requests/UR-085/input.md` for complete verbatim input.

---
*Source: validate-feedback Finding #13, captured by UR-085.*
