# UR-003 — Harness-bloat cleanup of the do-work skill

Captured: 2026-07-15T17:33:04Z
Source: user request (Claude Code session), three messages, verbatim below.

---

## Message 1

We're doing a harness-bloat cleanup of the do-work skill. A prior audit already
produced findings; treat them as hypotheses to verify, not conclusions to execute
blindly. Work in phases and STOP for my review between phases 2 and 3.

PHASE 1 — MAP THE ACTIVE PATH (read-only)
Build a table of every context-bearing asset (SKILL.md, actions/*, crew-members/*,
prompts/*, docs/*, next-steps.md, CHANGELOG.md) with columns:
  file | words | when it loads (always / per-action / JIT-gated / never) |
  job it does | evidence it still helps | duplication (other homes for same rule)
Then compute the actual context cost of the 5 most common invocations
(e.g. note, capture, work, pipeline run, board) — words loaded before user
content is touched. This map is the deliverable of phase 1.

PHASE 2 — CLASSIFY every asset into exactly one bucket:
  KEEP-ACTIVE   — earns its always-loaded place (justify in one line)
  LAZY-LOAD     — value is real but should load only when its phase of work
                  arrives (name the trigger condition)
  RELOCATE      — belongs in a sibling skill / separate repo, not do-work
                  (candidates to verify: bkb, interview, dream, ai-report,
                  deep-explore, the prompts/ library)
  HARDEN        — a prose rule that is machine-checkable (yes/no, countable,
                  schema-validatable) and should become a check instead of
                  instructions (candidates: work.md fractional steps 2.0, 3.5,
                  3.7, 5.5, 5.75, 6.25, 6.3, 6.5, 7.5 — for each, say whether
                  it's a testable condition or genuinely needs prose)
  DELETE        — no current job, or a duplicate of a rule that has another home
                  (candidates: 3 of SKILL.md's 4 action enumerations; Red Flags /
                  Verification Checklist / Common Rationalizations sections in
                  small mechanical actions like note, scan-ideas — keep them in
                  work.md and other heavy actions)
For each item cite the evidence from the phase-1 map. Where a candidate I named
turns out to be wrong, say so and why — don't classify to please the list.

STOP. Present phases 1–2 as a report. Wait for my approval per bucket.

PHASE 3 — EXECUTE (only approved items), as maintenance REQs via do-work's own
capture flow (maintenance: true), so every change has a REQ, a commit, and a
changelog entry. Ordering:
  1. SKILL.md router diet: routing table stays, the other three enumerations go.
     Target: router under ~1,500 words. Verify every verb still routes.
  2. CHANGELOG: truncate live file to last 20 entries, archive the rest,
     export-ignore the archive, keep a no-git-required pointer (see 0.76.1 for
     the tarball gap), confirm actions/version.md still parses the newest 5.
  3. Boilerplate strip in approved small actions.
  4. HARDEN items: convert each approved prose rule to a check; the prose shrinks
     to a one-line pointer at the check.
  5. RELOCATE items: extraction plan only (target repo/skill name, what pointer
     remains in do-work, migration note for existing users) — do not extract in
     this pass unless I say so.

PHASE 4 — RATCHET + RECEIPT
Add regression guards so this doesn't regrow: a router word-count budget checked
in CI or a pre-commit hook, and a rule in the contribution/maintenance docs that
any new action must state why it isn't a sibling skill. Finish with a receipt:
before/after words on the active path per common invocation, files touched,
rules deleted vs relocated vs hardened, and anything you chose NOT to touch
and why.

## Message 2 (mid-audit addition)

also please integrate the following:

Evidence discipline: label every phase-1/2 claim VERIFIED (observed in file or
run), INFERRED (from structure), or USER_REPORTED. Never convert a declared
loading rule ("always read X first") into a claim that X actually loads —
verify load timing against the dispatch logic or mark it INFERRED.

Add a PROBATION bucket: evidence too weak to justify changing it now; name
the bounded test that would settle it.

Before phase-3 execution, re-verify that each file to be changed still hash-
matches its phase-1 state; abort that item if the repo moved underneath you.

## Message 3 (after phase 1–2 report)

keep committing and pushing, there is no downside to it as long as commit
messages allow for traceability it's best to have frequent commits

this is a PR, you can build out everything, alternatively use multiple
branches/PR's

---

## Capture notes

- Phase 1–2 report: `decisions/audits/2026-07-15-harness-bloat-audit-phase1-2.md`
  (committed on this branch). Hash snapshot for the abort-if-moved rule:
  `decisions/audits/2026-07-15-phase1-hashes.txt`.
- Phase-2 corrections accepted into scope: five SKILL.md enumerations (not four);
  dispatch table stays; help menu is LAZY-LOAD not DELETE; router target revised
  to ≤2,500 words (1,500 unreachable without harming routing); boilerplate strip
  scoped to restatement-dedupe in 4 files; HARDEN limited to steps 2.0, 5.75
  (full) and 5.5, 6.3 (partial); ai-report and deep-explore NOT relocated
  (probation). RELOCATE remains extraction-plan-only per Message 1.
- Open decisions resolved by the builder (user offered no preference):
  hardened checks live in `tools/checks/` (shipped, platform-agnostic bash);
  single branch + single PR with one commit per REQ.
