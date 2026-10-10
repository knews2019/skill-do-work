---
id: REQ-694
title: 'Addendum: req append-section refuses a body that hides later sections, and frontmatter set refuses status and id'
status: claimed
created_at: 2026-10-10T19:19:53Z
user_request: UR-154
addendum_to: REQ-659
domain: backend
prime_files: ["_dev/primes/prime-releases.md", "skills/do-work/tools/do-work-cli/prime-do-work-cli.md"]
tdd: true
maintenance: false
impact: impact-user-visible
effort_estimate: effort-mechanical
related: [REQ-693, REQ-695, REQ-696]
batch: review-followups-ur154
required_lessons: ["_dev/primes/lessons-releases.md"]
write_set: ["skills/do-work/tools/do-work-cli/internal/corehelpers/request_writers.go", "skills/do-work/tools/do-work-cli/internal/corehelpers/request_writers_test.go"]
claimed_at: 2026-10-10T19:22:00Z
route: A
estimate:
  p50_active_minutes: 10
  confidence: high
  calculated_at: 2026-10-10T19:23:55Z
  basis:
    - Route A
    - 2-file write set
    - 2 acceptance criteria
---
# Addendum: req append-section Refuses a Body That Hides Later Sections, and frontmatter set Refuses status and id

## What

Two schema gaps the REQ-659 (frontmatter set and req append-section writers) reviews left as report only. `req append-section` accepts a body that leaves a code fence or HTML comment open, which hides every later section from `advance` (re-review N1). `frontmatter set` writes any value into `status` or `id`, although those fields have owning commands (review F3).

## Prior Implementation

REQ-659 shipped both writers in `skills/do-work/tools/do-work-cli/internal/corehelpers/request_writers.go`, archived at commit `8a33190c`. `handleRequest` inserts the section in `requestmodel.CanonicalSectionOrder` and then re-checks that exactly one visible copy of the new section exists. It does not check that the sections after the insertion point stay visible. `handleFrontmatterSet` guards `*_at` stamps (append-only, canonical form) and structured fields, and nothing else.

## Detailed Requirements

- N1: after the insertion, count the column-0 visible sections (`requestmodel.VisibleSections`) and refuse with `SECTION-WRITE-FAILED`, writing nothing, unless the count grew by exactly one. This replaces or extends the current single-copy re-check (the reviewer notes it also covers the end-of-file case).
- F3: `frontmatter set` refuses the fields `status` and `id`, and the refusal names the owning command (status: the lifecycle transitions such as `advance`, `unblock`, `finalize`; id: never rewritten). Other scalar fields keep today's behaviour.

## Red-Green Proof

**RED prompt/case:** (1) A REQ with a `## Review` section; `req append-section REQ-N --section Testing --from body.md` where `body.md` is a line of three backticks followed by `## x` (or `text <!-- open`). (2) `frontmatter set REQ-N status bogus-status`.
**Why RED now:** (1) exits 0 and the following `## Review` disappears inside the open fence or comment, so `advance` reads the phase as missing. (2) exits 0 and writes `status: bogus-status`.
**GREEN when:** (1) refuses `SECTION-WRITE-FAILED` and the file is byte-identical. (2) refuses, names the owning command, and the file is byte-identical.
**Validation:** Inferred during capture

## Constraints

- One test per failure in `request_writers_test.go`; each must fail with the new guard removed.
- No new allowed-values table for other fields (YAGNI): only `status` and `id` are refused.
- No `depends_on` edges in this batch; no other REQ touches these files.

## Builder Guidance

Certainty: high (both reproduced by the reviewer, re-checked at HEAD `0e8e0ef9`). The reviewer sized N1 at about four lines.

## Required Lessons — Dropped for Budget

- `skills/do-work/tools/do-work-cli/lessons-do-work-cli.md` — 20802 tokens, bare only (`slugged: partial`); matches the do-work-cli corehelpers path (families `writer-anchored-marker-position`, `lifecycle-section-evidence`).

## Full Context

See `do-work/user-requests/UR-154/input.md` for complete verbatim input. Sources: `do-work/archive/UR-145/REQ-659-frontmatter-set-and-req-append-section.md` → `## Review`, `do-work/runs/work-2026-10-10-131527/REQ-659-review.md` (F3), `do-work/runs/work-2026-10-10-131527/REQ-659-rereview.md` (N1).

## AI Execution State (P-A-U Loop)
- [ ] **[PLAN]:** (Agent: Read listed `prime_files` and agent rules. Write brief technical approach here. Do not write code yet.)
- [ ] **[APPLY]:** (Agent: Code written exactly as planned. Scope strictly limited to planned files.)
- [ ] **[UNIFY]:** (Agent: Run `git diff --stat` and review every changed file. Run native project linters. Verify no debug artifacts in diff. List each file you verified and what you checked.)
*Source: UR-154 R2 and R3 — "`req append-section` accepts a body that leaves a code fence or comment open, which hides every later section from `advance` (review N1; suggested fix: refuse unless the visible section count grows by exactly one)." and "`frontmatter set` accepts any value for `status` (review F3, impact-rule-change)."*

## Triage

**Route: A** - Simple

**Reasoning:** The REQ names the one production file and its test file, and both changes are decided: a visible-section count check after the insert in `handleRequest`, and a refusal for `status` and `id` at the top of `handleFrontmatterSet`. The reviewer reproduced both failures and sized the first fix at about four lines.

**Planning:** Not required

No `## Open Questions` section exists, so Step 3.5 records nothing there. Pre-dispatch decisions:

- **D-01 (DECIDE & STATE):** Route A, not B. The location is known (`request_writers.go:161-171` and `:25-41`), so exploration would only re-read what this triage already read. The builder brief carries the file:line pointers.
- **D-02 (DECIDE & STATE):** N1 replaces the single-copy re-check at `request_writers.go:161-171` with the count check (column-0 visible sections after the insert must equal the count before plus one). The REQ allows "replaces or extends"; the count check also covers the end-of-file case the old check guarded (a hidden new heading does not raise the count), so keeping both would be dead code.
- **D-03 (DECIDE & STATE):** F3 adds one new refusal code, `FRONTMATTER-FIELD-OWNED`, checked right after argument parsing in `handleFrontmatterSet` and before the target is resolved or read, so the file is never touched. The evidence text names the owner: for `status`, the lifecycle commands (`advance`, `unblock`, `finalize`) or the hand write at the transition's defining site (clarify's `pending-answers` to `pending` flip is a hand write, `actions/clarify.md:136`); for `id`, that an id is never rewritten. A shipped-prose search found no action that tells an agent to run `frontmatter set` on `status` or `id`, so no caller breaks.
- **D-04 (DECIDE & STATE):** Because D-02 removes the only guard for an unclosed fence or comment at the end of the file (a hidden new heading), the new N1 test carries that case as a third table row beside the open-fence and open-comment bodies. It pins that the replacement keeps the old guarantee; no test pinned it before. Each row fails with the count check removed.
- **D-05 (DECIDE & STATE):** Test names are fixed so the GREEN probe can name them: `TestRequestAppendSectionRefusesBodyThatHidesLaterSections` and `TestFrontmatterSetRefusesLifecycleOwnedFields`, both in `request_writers_test.go`. The builder's first implementation decision is D-06.

<!-- D-XX counter: last used D-05. Next decision: D-06. -->

## Plan

**Planning not required** - Route A: Direct implementation

*Skipped by work action*
