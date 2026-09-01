---
id: REQ-441
title: '[impact-critical] Validate HTTP archives before publication'
status: pending
created_at: 2026-08-31T14:19:37Z
user_request: UR-083
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: [REQ-414]
maintenance: false
impact: impact-critical
effort_estimate: effort-substantive
related: [REQ-437, REQ-438, REQ-439, REQ-440, REQ-442, REQ-443, REQ-444]
batch: accepted-feedback-regressions
---

# Validate HTTP Archives Before Publication

## What

Download HTTP archive candidates to a private sibling, validate them there, and only then publish onto an absent or regular non-symlink target. A failed HTTP validation followed by failed Git fallback must leave every pre-existing target byte-identical.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Finding Provenance

- **Verbatim claim / severity:** `[P1] Validate HTTP downloads before replacing the target.`
- **Evidence:** `archive_fetch.go:117-124` passes `ArchiveTargetPath` directly to `atomic-download.sh` and validates only after the helper has published it.
- **REQ-414 seam reconciliation:** REQ-414 removes the shell-publication call and validates the HTTP body inside Go before no-overwrite publication. Once REQ-414 is terminal, re-triage this REQ against the remaining cross-route target-shape/replacement contract; do not replay the superseded `atomic-download.sh` evidence.
- **Origin / earned by:** The Go port `f27f564d` inherited the behavior from shell commit `0e8cf0d9`. A successful HTTP transfer containing an unreadable tar followed by failed Git fallback replaced a pre-existing archive despite `FetchArchive`'s preservation contract.
- **Surface-cost:** Earned. The destructive replay and symlink/special-target class justify one private stage lifecycle and target-mode guard; tests must preserve an old regular target and refuse non-regular targets.

## Detailed Requirements

- Allocate HTTP staging beside the public target and pass only that private path to the download primitive.
- Validate the staged candidate as a readable archive before publication.
- Before publication, require the public target to be absent or an existing regular non-symlink file.
- On HTTP invalidity and Git fallback failure, preserve a pre-existing target byte-for-byte and remove only invocation-owned scratch.
- Apply a consistent target-shape rule to both HTTP and Git publication paths.

## Constraints

- Preserve the two-route fallback contract and truthful route evidence.
- Do not reimplement archive parsing; retain the existing readability check.
- Keep private staging same-directory so final publication can remain atomic on supported platforms.

## Red-Green Proof

**RED prompt/case:** Seed a regular archive target, make HTTP return success with a non-tar body, and make Git fallback fail; separately use a symlink or special file as the public target.
**Why RED now:** HTTP publishes onto the public target before Go validates the candidate.
**GREEN when:** Total failure preserves the old regular target byte-for-byte with no scratch, non-regular targets are refused unchanged, and a valid HTTP archive still publishes successfully.
**Validation:** User confirmed by requesting capture of every accepted validation finding.

## Builder Guidance

Certainty level: Firm on staging and target preservation. Reuse existing publication mechanics where that reduces, rather than expands, surface.

## Full Context

See `do-work/user-requests/UR-083/input.md` for the complete capture provenance.

---
*Source: accepted Finding 20 from the validated external feedback.*
