---
id: REQ-461
title: '[impact-user-visible] Require affirmative project-owned release targets'
status: claimed
created_at: 2026-09-01T00:12:38Z
user_request: UR-081
domain: general
prime_files: [_dev/primes/prime-action-files.md]
tdd: true
suggested_spec: bug-fix
depends_on: []
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-413]
batch: go-no-llm-command-platform
review_generated: true
addendum_to: REQ-413
claimed_at: 2026-09-03T10:32:09Z
route: B
write_set:
  - skills/do-work/tools/do-work-cli/internal/publication/release.go
  - skills/do-work/tools/do-work-cli/internal/publication/release_test.go
estimate:
  p50_active_minutes: 25
  confidence: medium
  calculated_at: 2026-09-03T10:33:01Z
  basis:
    - Route B
    - 2-file write set
    - 1 subsystem involved
    - 5 acceptance criteria
    - cross-route regression gates
---

# Require Affirmative Project-Owned Release Targets

## What

Replace convention-based installed/generated path exclusions with condition-complete evidence that every release target is a project-owned source or declared maintainer mirror.

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)

## Finding Provenance

REQ-413's fresh re-review found that release exclusion recognizes named conventions such as `vendor`, `node_modules`, `.codex/skills`, and literal generated directories, but accepts other dependency or generated locations such as `third_party/do-work`, `dist/skills/do-work`, and cache-owned skill trees.

## Detailed Requirements

- Require affirmative, verifiable project-owned classification for every consumer release target instead of inferring safety from the absence of known bad directory names.
- Permit declared maintainer mirrors only through the existing explicit mirror contract and retain byte-identity validation for changelog mirrors.
- Refuse installed skills, dependencies, vendored packages, caches, distribution outputs, and generated trees regardless of their directory spelling.
- Add fixtures for non-example paths including `third_party/do-work`, `dist/skills/do-work`, and an arbitrarily named cache tree.
- Keep refusal actionable by identifying the target and the ownership evidence that is missing or inconsistent.

## Constraints

- Do not replace one directory-name denylist with a larger denylist.
- Preserve caller-selected semantic version and changelog-content judgment.

## Dependencies

No request prerequisite.

## Builder Guidance

Certainty level: Firm on the invariant, flexible on the proof mechanism. Prefer an explicit manifest or repository-owned declaration that the planner can verify mechanically.

## Red-Green Proof

**RED prompt/case:** Supply otherwise valid release targets beneath `third_party`, `dist`, and an unrecognized cache/generated subtree.
**Why RED now:** The current exclusion is a finite convention list, so targets outside those spellings are accepted.
**GREEN when:** Undeclared or non-project-owned targets refuse independent of path spelling, while verified project sources and declared maintainer mirrors apply atomically.

## Full Context

See `do-work/user-requests/UR-081/input.md` and `do-work/runs/work-2026-08-31-165510/REQ-413-rereview.md` for source context and independent evidence.

---
*Source: REQ-413 fresh re-review finding F7.*

---

## Triage

**Route: B** - Medium

**Reasoning:** The invariant is firm and the single function is obvious, but the REQ deliberately leaves the proof mechanism open ("flexible on the proof mechanism"), and what "project-owned" means *mechanically* is not stated anywhere — that had to be established from what the repository can already verify.

**Planning:** Not required

## Plan

**Planning not required** - Route B: Exploration-guided implementation

*Skipped by work action*

## Exploration

`installedReleasePath` (`internal/publication/release.go:170-189`) is a denylist of directory names: `vendor`, `vendored`, `node_modules`, a `.codex`/`.claude` followed by `skills`, a `generated`/`.generated` with a `skills` or `do-work` descendant, and the literal prefix `skills/do-work/`. Its two callers (`:34`, `:72`) use it to stop a **consumer** release from mutating installed or vendored suite metadata — a consumer releasing their own project must not bump the version of the do-work copy installed inside their repo.

Everything the REQ names slips through, and I confirmed each is only a spelling away from a listed one: `third_party/do-work`, `dist/skills/do-work`, and any cache tree whose directory happens not to be called `generated`. This is the same shape as REQ-460 — a contract that states a condition, implemented as membership in a list of examples — and the same family, `closed-enumeration-for-a-condition`.

Two mechanisms the repository can already verify, either of which is affirmative rather than exclusionary:

1. **Package markers.** An installed suite package carries its own `SKILL.md` at the top of its directory. The board tool's repo-root walk already relies on exactly this to skip directories "merely *named* `do-work` that are skill installs (SKILL.md at their top level)" — so a target with a package marker at or above it is an installed copy regardless of what the enclosing directory is called.
2. **Git tracking.** A project-owned source is tracked in this repository; a distribution output or cache tree generally is not. This catches `dist/` and cache trees by what they *are* rather than by what they are named.

Neither alone is complete — a repository that commits its `vendor/` tree is tracked but not project-owned, and a package marker does not exist for a generated non-skill artifact — which is why the affirmative test likely needs both, and why `maintainer_release: true` must remain the only door to the suite's own metadata.

The constraint is the sharp edge here: *do not replace one directory-name denylist with a larger denylist*. A fix that adds `third_party`, `dist` and `.cache` to the list satisfies the fixtures and fails the REQ.

*Generated in-session (single-pass discovery)*

## Scope

**Files I will touch:**
- `skills/do-work/tools/do-work-cli/internal/publication/release.go` (modify) — replace the convention denylist with an affirmative project-owned test, and make the refusal name the missing evidence
- `skills/do-work/tools/do-work-cli/internal/publication/release_test.go` (modify) — fixtures for the three non-example paths the REQ names plus an arbitrarily-named cache tree, and controls for a genuine project source and a declared maintainer mirror

**Files I will NOT touch:** the semantic-version comparison and changelog-content judgment (the REQ preserves caller-selected judgment), and the changelog byte-identity mirror validation.

**Acceptance criteria (restated from REQ):**
- [ ] Every consumer release target requires affirmative, verifiable project-owned classification rather than the absence of known-bad directory names
- [ ] Declared maintainer mirrors are permitted only through the existing explicit mirror contract, with changelog byte-identity validation retained
- [ ] Installed skills, dependencies, vendored packages, caches, distribution outputs and generated trees refuse regardless of directory spelling
- [ ] Fixtures cover `third_party/do-work`, `dist/skills/do-work`, and an arbitrarily named cache tree
- [ ] Refusals name the target and the ownership evidence that is missing or inconsistent
- [ ] No directory-name denylist, larger or otherwise, replaces the old one
