---
id: REQ-443
title: '[impact-critical] Keep Git fallback archive prefixes to one component'
status: pending
created_at: 2026-08-31T14:19:37Z
user_request: UR-083
domain: backend
prime_files: [skills/do-work/tools/do-work-cli/prime-do-work-cli.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-critical
effort_estimate: effort-mechanical
related: [REQ-437, REQ-438, REQ-439, REQ-440, REQ-441, REQ-442, REQ-444]
batch: accepted-feedback-regressions
---

# Keep Git Fallback Archive Prefixes to One Component

## What

Use a constant single-component prefix for Git-fallback archives regardless of the selected branch name. Branch text must select the cloned ref only; it must not alter the extraction depth expected by install and update.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Finding Provenance

- **Verbatim claim / severity:** `[P2] Keep Git archive prefixes to one path component.`
- **Evidence:** `archive_fetch.go:108-110` embeds the branch name in `repo-<branch>/`, while install and update extract with `--strip-components=1`.
- **Origin / earned by:** Shell commit `0e8cf0d9` introduced the prefix shape, `f27f564d` preserved it, and exact branch selection in REQ-424 (Clone the Branch Named by the Tarball URL) made slash-containing refs load-bearing. A `release/2.x` replay extracted `VERSION` below `2.x/`, so root manifest validation failed.
- **Surface-cost:** N/A. A constant prefix deletes the incorrect coupling and adds no defensive apparatus.

## Detailed Requirements

- Keep derived branch names as exact Git clone selectors, including names containing `/`.
- Generate every Git-fallback archive beneath one constant path component.
- Preserve the production extraction contract that strips exactly one component.
- Cover a slash-containing local fixture branch through fetch, archive, production-style extraction, and root manifest-file assertion.

## Constraints

- Do not sanitize branch text into a second naming scheme; remove it from the prefix entirely.
- Preserve missing-branch failure and exact requested-branch selection from REQ-424.

## Red-Green Proof

**RED prompt/case:** Force HTTP failure, fetch a local Git branch named `release/2.x`, extract with production's `--strip-components=1`, and look for `VERSION` at the extraction root.
**Why RED now:** The archive prefix contains `repo-release/2.x/`, so one stripped component leaves the suite nested below `2.x/`.
**GREEN when:** The requested slash branch is still selected exactly and the extracted suite—including `VERSION` and manifest inputs—lands directly at the root.
**Validation:** User confirmed by requesting capture of every accepted validation finding.

## Builder Guidance

Certainty level: Firm. Reuse the existing default one-component prefix rather than adding branch-name escaping.

## Full Context

See `do-work/user-requests/UR-083/input.md` for the complete capture provenance.

---
*Source: accepted Finding 23 from the validated external feedback.*
