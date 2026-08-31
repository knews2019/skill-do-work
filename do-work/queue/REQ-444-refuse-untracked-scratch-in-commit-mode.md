---
id: REQ-444
title: 'Refuse untracked consumed scratch in cleanup commit mode'
status: pending
created_at: 2026-08-31T14:19:37Z
user_request: UR-083
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
related: [REQ-437, REQ-438, REQ-439, REQ-440, REQ-441, REQ-442, REQ-443]
batch: accepted-feedback-regressions
---

# Refuse Untracked Consumed Scratch in Cleanup Commit Mode

## What

Refuse entirely untracked consumed-scratch deletion whenever cleanup runs with `--commit`. Commit mode must not report success after a scratch-only deletion that leaves HEAD unchanged or silently perform scratch deletion outside an otherwise valid tracked cleanup commit.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Finding Provenance

- **Verbatim claim / severity:** `[P2] Honor --commit before deleting untracked scratch.`
- **Evidence:** `cleanup_apply.go:44-49` admits scratch after preflight failure; when no tracked group is eligible, the transaction is skipped and scratch deletion still runs, so success can leave HEAD unchanged.
- **Origin / earned by:** The explicit non-rollback exception in `a57bf51e` did not define a commit-mode boundary. Scratch-only and mixed tracked/scratch replays show deletion outside the requested exact-path commit.
- **Surface-cost:** Earned. One `options.Commit` refusal and two regressions are cheaper than a successful destructive side effect that the requested commit and rollback cannot cover.
- **Fold-first result:** REQ-432 (Enforce the Commit Guard for Consumed Scratch Cleanup) owns the related nonempty-index incident but is dependency-gated through REQ-431, so the canonical fold-first rule forbids widening it and requires this independently runnable REQ.

## Detailed Requirements

- Refuse every entirely untracked consumed-scratch group when `options.Commit` is true, even with an empty index.
- Cover both a scratch-only run and a mixed run containing otherwise eligible tracked cleanup work.
- Leave the scratch inventory byte-identical and report truthful refusal evidence.
- Preserve the existing exact-inventory, rooted-containment, and consumed-manifest checks.
- Preserve non-commit cleanup's narrow, explicitly labeled non-rollback scratch deletion.
- Leave REQ-432's separate global empty-index invariant intact.

## Constraints

- Add one commit-mode boundary to the existing exception; do not add a new transaction class for untracked scratch.
- Shared files with REQ-432 do not create a dependency; this request must remain independently selectable because REQ-432 is currently gated.

## Red-Green Proof

**RED prompt/case:** With an empty index, run `cleanup --commit` first on an entirely untracked consumed run alone, then on that scratch beside one eligible tracked cleanup group.
**Why RED now:** The scratch exception runs outside the transaction; scratch-only cleanup succeeds without a commit, and a mixed run deletes scratch outside its tracked commit.
**GREEN when:** Both commit-mode fixtures refuse and preserve the scratch bytes, tracked commit behavior remains truthful, and the same scratch is still deleted by an otherwise identical non-commit cleanup run.
**Validation:** User confirmed by requesting capture of every accepted validation finding.

## Builder Guidance

Certainty level: Firm. The accepted remedy is the narrow `options.Commit` refusal, not an attempt to make Git roll back untracked scratch.

## Full Context

See `do-work/user-requests/UR-083/input.md` for the complete capture provenance.

---
*Source: accepted Finding 16 from the validated external feedback.*
