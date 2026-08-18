---
id: REQ-247
title: Archive timestamp audit tool driven by git commit times
status: completed
created_at: 2026-08-18T12:38:26Z
claimed_at: 2026-08-18T18:25:40Z
completed_at: 2026-08-18T19:11:56Z
commit: 4035ddc
kb_status: pending
route: C
user_request: UR-056
domain: general
prime_files: [_dev/primes/prime-shell-commands.md]
tdd: true
suggested_spec:
depends_on: [REQ-246]
maintenance: false
related: [REQ-246, REQ-244, REQ-245]
batch: timestamp-stamping-integrity
write_set:
- skills/do-work/scripts/audit-archive-timestamps.sh
- skills/do-work/actions/capture.md
- _dev/tests/prescribed-shell-scripts-behavior.sh
- skills/do-work/scripts/repair-req-timestamps.sh
estimate:
  p50_active_minutes: 40
  confidence: medium
  calculated_at: 2026-08-18T18:26:21Z
  basis:
    - Route C
    - 4-file write set
    - 1 new file
    - 2 subsystems involved
    - 5 acceptance criteria
    - cross-route regression gates
    - full-suite verification
---

# Archive Timestamp Audit Tool Driven by Git Commit Times

## What

A deliberately-invoked audit tool that scans `do-work/archive/` for detectably wrong `*_at` stamps and repairs them, deriving every replacement from git commit times — the author time of the commit that introduced the stamp. Never run from a hook: repairing the archive is an exception to the immutability rule and stays a conscious invocation.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read the prime, four crew files, REQ-246's script/hook/suites, capture.md and review-work.md precedents, and the constraining contract suites. Approach fixed before code: library-mode refactor of the repairer, thin sourcing auditor, capture.md co-located amendment, five lock-ins. (Transcribed from builder hand-back.)
- [x] **[APPLY]:** Exactly the REQ's four write-set files; no hooks, no board tool, all fixtures in scratch mktemp space. (Transcribed from builder hand-back.)
- [x] **[UNIFY]:** `git diff --stat ad69e56..HEAD`: 4 files, +272/−4, reviewed file-by-file; ShellCheck warning-level clean; no debug artifacts; maintainer-verify exit 0; extra edge probes (no-archive, unknown flag, symlinked root, quoted/space-separated shapes). (Transcribed from builder hand-back.)

## Detailed Requirements

- **Replacement source is git only** (user-specified): the author time of the commit that introduced the stamp line (`git log -L`-style lookup or equivalent). File mtimes are never consulted here — checkout resets them, so they carry no signal for committed archive content.
- **Same detection predicate as REQ-246:** future beyond the 2-minute skew, and impossible field orderings. Share the detection/derivation logic with REQ-246's script rather than duplicating it — that shared logic is why this REQ depends on REQ-246.
- **Ordering clamp** across the repaired file: `created_at ≤ claimed_at ≤ completed_at`, each no later than its introducing commit's time.
- **Amend the archive-immutability rule in the same commit** that ships the tool: `skills/do-work/actions/capture.md` § Immutability Rule (~line 63) gains a stated mechanical-timestamp-repair exception alongside the existing review-annotation exception (`actions/review-work.md` ~line 448 is the precedent wording). Co-location rule: the exception and the tool land together, and any other statement of archive immutability found in the sweep is amended in that commit too.
- **Audit trail:** print each correction (file, field, old value, new value, sourcing commit hash) to stdout. No new frontmatter fields.
- **Dry-run by default is acceptable builder latitude** (report what would change, `--fix` to write); if implemented, the RED/GREEN below refers to the fixing mode.

## Implementation Summary

**What was done:** A deliberately-invoked archive auditor, `scripts/audit-archive-timestamps.sh`, scanning `do-work/archive/**/REQ-*.md` at any depth for REQ-246's detection predicate, deriving every replacement from the introducing commit's author time (`git blame --line-porcelain`; mtime disabled in this path; unanswerable blame reports and leaves bytes untouched). Report-only by default with exit 1 on findings, `--fix` writes through the repairer's full guard set. Sharing is by **sourcing**: the repairer became a sourceable library (two pre-source switches, report-only bail, return guard), so predicate, shape recognizer, clamp, and atomic-write guards are one code body — REQ-255's widening reaches both tools in one edit. `capture.md` § Immutability Rule gained the mechanical-timestamp-repair exception in the same commit as the tool, with two co-located restatement fixes; the orchestrator applied the builder's offered seam scoping `present-work-guide.md`'s immutability restatement to its own workflow. Never hook-wired.

**Files changed:**
- `skills/do-work/scripts/audit-archive-timestamps.sh` (new) — the auditor: arg parse, symlink refusal, missing-archive no-op, sources the repairer with git-only + report-only switches, recursive scan, summary/exit semantics
- `skills/do-work/scripts/repair-req-timestamps.sh` (modified) — shared-library switches with hook-preserving defaults, git-only gate, parameterized failure messages, report-only planner, library return guard
- `skills/do-work/actions/capture.md` (modified) — § Immutability Rule second exception ("these are the only exceptions; neither is a precedent"); two restatements aligned
- `_dev/tests/prescribed-shell-scripts-behavior.sh` (modified) — five `audit-archive-timestamps:` lock-ins (fix/report-only/clean-scope/ordering-clamp/blameless)
- `skills/do-work-toolbox/docs/present-work-guide.md` (modified) — integration seam: immutability restatement scoped "to this workflow"

*Integrated by orchestrator from builder hand-back; merge range `6268320..4035ddc`.*

## Decisions

Transcribed from the builder hand-back:

- **D-01 (DECIDE & STATE):** share by sourcing the repairer, not a third lib file — one file owns the shape-set, change stays in the write set; direct execution behavior-identical (hook suite green).
- **D-02 (DECIDE & STATE):** report-only exits 1 when corrections are pending — linter convention, usable as a gate.
- **D-03 (DECIDE & STATE):** unanswerable blame reports through `report_failure` and exits 1, file byte-identical — honors "never invents" while keeping the exit contract's meaning.
- **D-04 (DECIDE & STATE):** shipped headers state the one-shape-set property without citing "REQ-255" by number — a maintainer queue id means nothing in a consumer install; the linkage is recorded here.
- **D-05 (DECIDE & STATE):** sweep classification — three capture.md statements amended; workflow-scoped statements left as true; one borderline handed back as a seam (applied by the orchestrator).
- **D-06 (DECIDE & STATE):** archive scan is `find -name 'REQ-*.md'` at any depth — filename condition, not a directory list (Closed Enumerations Go Stale).
- **D-07 (logged):** report-only mode bypasses the write guards — they protect writes, and report-only writes nothing; `--fix` runs the full set.

## Qualification

Passed — 5 files verified in merge range `6268320..4035ddc` (4 builder + 1 seam), all six acceptance criteria traced (git-only proven by the blameless lock-in — file left byte-identical although a valid mtime existed; sourcing = zero duplicated logic; clamp pinned; same-commit amendment in `a6764fd`; no hooks touched; gate green), P-A-U audited. The known space-separated shape holes were probed and confirmed **pre-existing and shared, not copied** — recorded for REQ-255, whose fix now lands in both tools at once.

## Review

**Overall: 94%** | 2026-08-18T19:09:56Z

| Dimension | Score |
|-----------|-------|
| Requirements | 97% |
| Code Quality | 92% |
| Test Adequacy | 88% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Pass |

**Verdict: Approve** — git-author-time-only replacements, shared predicate by sourcing, deliberate invocation only, dry-run default, and the immutability amendment co-located in the tool's own commit; every core behavior reproduced by execution in scratch git repos.

### Requirements Checklist

- [x] **Git only; mtime never consulted for archives** — *reproduced*: stamp committed in commit A (known author date), body edited in commit B → repair wrote **A's** time; a dirty archive file with a valid mtime available was refused byte-identical, failure message naming git blame alone.
- [x] **Shared predicate by sourcing** — the auditor contains no predicate/shape/clamp/guard code of its own (grep-verified); quoted-shape canonicalization works through the audit route.
- [x] **Ordering clamp** — delivered, with one correctly-resolved internal tension recorded: when the derived commit time precedes the anchor, the floor clamps up to the predecessor, which can exceed the introducing commit's time — the two clauses are unsatisfiable together there; ordering wins, the audit line says `clamped to <anchor>`, and the shipped lock-in pins exactly this shape.
- [x] **Same-commit amendment** — verified from git: `a6764fd` carries `capture.md` and the tool together; the seam landed in the merge commit itself, matching the trail.
- [x] **Never hook-wired** — no reference under `hooks/`; hook path exercised against a scratch project: banner intact, mtime fallback still active for dirty queue files (the git-only switch does not leak), direct execution identical.
- [x] **Dry-run default / `--fix`** — *reproduced*: report-only leaves bytes and exits 1 with a rerun hint; `--fix` writes, rerun clean, exit 0; arg/usage/missing-dir/symlinked-root edges all behave.
- [x] **Audit trail with sourcing hash; no new frontmatter** — every correction line carries the 7-char hash or its annotation.
- [x] **Deep-path + legacy scan** — nested UR folder and top-level legacy REQ both scanned in one run.
- [x] **`maintainer-verify` exit 0** — observed un-piped at review end (52 named script cases).

### Findings

**Important:** None

**Minor (report only, 3):**
1. *(reproduced)* A **symlinked archive REQ file is silently skipped and counted scanned-clean** — safe direction holds (never write through a link), but silence-reads-as-clean is what the script's header refuses for the root case, and the "worker refuses symlinked files" comment overstates a silent `return 0`.
2. *(reproduced)* The `${var:-default}` switches are an environment override surface: an exported `timestamp_repair_apply_mode=0` flips the hook repairer to report-only, and in that leaked mode direct execution exits 0 with a defect unrepaired. Requires exporting those exact names — realistically inert.
3. *(reproduced)* Mixed `--fix` run (one repairable + one blameless) writes the repairable file and exits 1 correctly per the header contract, but no lock-in pins that branch.

**Pre-existing, not re-reported:** the space-separated / duplicate-key / calendar-impossible shapes behave through the audit path exactly as queued REQ-255's instances describe — inherited via the shared sourced code, not worsened; REQ-255's fix lands in both tools in one edit, as D-01 claims.

**Restatement sweep:** capture.md ×3 amended with the "only exceptions; neither is a precedent" close; every other immutability statement is workflow-scoped and remains true; the orchestrator's seam closes the one borderline. Nothing stale.

**Scope:** `scope-drift.sh` names only the orchestrator-applied seam file — not builder drift; the builder touched exactly its declared set.

**Acceptance: Pass.** **Suggested testing:** lock in the mixed-run branch; voice the symlinked-archive skip; when REQ-255 lands, add one archive-path fixture per widened shape so shared-fix-reaches-both-tools is pinned, not assumed.

**Follow-ups created:** None · **Sweeps appended to:** none by the reviewer (the both-tools pinning note threaded into REQ-255's requirements by the orchestrator)

*Reviewed by review-work action (independent adversarial pass, orchestrated mode; merge range `6268320..4035ddc`)*

## Lessons Learned

**What worked:** Sharing by sourcing rather than a copied or third-file library — the reviewer verified by grep that the auditor holds zero predicate code, which is what turns REQ-255's future fix into a single edit that reaches both tools. The blameless lock-in that proves mtime is *never* a fallback (file left byte-identical although a valid mtime existed) is the strongest kind of negative evidence this suite has.

**What didn't:** The ordering-clamp requirement's two clauses ("clamped to the anchor" and "no later than the introducing commit's time") are unsatisfiable together when the derived commit time precedes the anchor — the REQ text carried a latent contradiction nobody caught at capture or verify. Resolved in favor of ordering with a transparent audit annotation; a REQ stating two ceilings should state which one wins.

**Worth knowing:** The library switches (`timestamp_repair_apply_mode`, `timestamp_repair_git_only`) read `${var:-default}`, so exported environment variables of those names override them — inert in practice, worth knowing before renaming. A symlinked archive REQ file is silently skipped and counted clean (review Minor 1). Mixed `--fix` runs write what they can and exit 1.

## Orientation

Now the archive can be audited and mechanically repaired — `scripts/audit-archive-timestamps.sh`, deliberate invocation only, git author times as the only truth source, dry-run by default — and the archive-immutability rule names this as its second and final bounded exception. Lives in core's `scripts/` subsystem as a sourcing consumer of the repairer library. [MAP CHANGED] — `repair-req-timestamps.sh` is now a shared library with two consumers, and the immutability rule has a stated exception set. Prime staleness spot-check: `prime-shell-commands.md` paths still resolve; not stale.

## Constraints

- Not in the board tool (read-only frontmatter decision, pinned write-surface count).
- Never wired into any hook. Manual invocation only — the user runs it as an audit.
- `record-commit-hash.sh` guard style: verify before replace, atomic write, tripped guard leaves the file byte-identical.
- Repairs are ordinary commits: the tool edits the working tree and reports; committing the repaired archive files follows the normal commit flow.
- Provenance: second half of the Finding 2 replacement (UR-055 triage → UR-056); requested verbatim by the user in the ask-tool answer.

## Red-Green Proof

**RED prompt/case:** A `_dev/tests/` lock-in test creates a scratch git repo with an archived fixture REQ whose `completed_at` is future-dated relative to its committing instant, runs the audit tool in fixing mode, and asserts the stamp is rewritten to the introducing commit's author time with the correction logged. Fails today because the tool does not exist.
**Why RED now:** Archived wrong stamps are permanent; the board warns on every render and nothing can correct them.
**GREEN when:** The test passes, a clean archived fixture passes through byte-identical, the immutability-rule amendment ships in the same commit as the tool, and `bash _dev/tests/maintainer-verify.sh` exits 0.
**Validation:** User confirmed ("please make an audit tool that will fix all the archive, but there it needs to take the timestamp of the git commits where it was commited")

## Full Context

See `do-work/user-requests/UR-056/input.md` for complete verbatim input.

---
*Source: ask-tool answer — "so please make an audit tool that will fix all the archive, but there it needs to take the timestamp of the git commits where it was commited"*

---

## Triage

**Route: C** - Complex

**Reasoning:** A new audit tool with git-blame-driven derivation, an ordering clamp, a same-commit immutability-rule amendment across every statement of the rule found in a sweep, and shared logic with REQ-246's script — multiple components and a co-location constraint.

**Planning:** Required

## Plan

**Approach**

1. **Share REQ-246's logic by extraction, not duplication.** `repair-req-timestamps.sh` already owns the detection predicate (future beyond 120s, impossible orderings) and the canonical-stamp comparison machinery. Factor what the auditor needs into functions the repairer keeps (or a sourced lib file if that stays ShellCheck-clean and suite-green); the auditor must not re-implement the predicate — that shared logic is why this REQ depends on REQ-246.
2. **Replacement source is git only.** For each defective stamp line in `do-work/archive/**`, derive the author time of the commit that *introduced* that line (`git blame --line-porcelain` is the proven mechanism from REQ-246's committed-file path). Never mtime — checkout resets archive mtimes, so they carry no signal. A blame that cannot answer skips with a report line, never invents.
3. **Clamp** `created_at ≤ claimed_at ≤ completed_at`, each no later than its introducing commit's time.
4. **Dry-run by default, `--fix` to write** (the REQ names this acceptable latitude; the Red-Green Proof refers to the fixing mode). Reuse the repairer's guard set on the write path: verify-before-replace, atomic rename, tripped guard leaves the file byte-identical, nonzero on failure, one audit line per correction with the sourcing commit hash.
5. **Amend the archive-immutability rule in the same commit.** `skills/do-work/actions/capture.md` § Immutability Rule (~line 63) gains a mechanical-timestamp-repair exception alongside the review-annotation precedent (`actions/review-work.md` ~line 448 wording). Sweep for any *other* statement of archive immutability in shipped text and amend those in the same commit — the co-location rule.
6. **Never hook-wired.** No `hooks/` edit anywhere; manual invocation only.
7. **Lock-in cases** in `_dev/tests/prescribed-shell-scripts-behavior.sh`: future stamp in a committed archive fixture repaired to the introducing commit's author time under `--fix`; dry-run reports and leaves bytes; ordering clamp; unreachable-blame skip. Fixtures are scratch git repos, never this repo's own archive.

**Watch for:** the class-vs-instance shape — REQ-246's review found the repairer's shape recognition narrower than the read-side detectors (REQ-255 pends on that). Do not copy that hole into the auditor: reuse the same extraction code so the two tools stay one shape-set, and say in the header that REQ-255's resolution will widen both at once.

*Written inline by the orchestrator (no separate Plan agent) — Route C.*

## Scope

**Files I will touch:**
- `skills/do-work/scripts/audit-archive-timestamps.sh` (new) — the deliberately-invoked archive auditor: git-author-time derivation, ordering clamp, dry-run default, audit trail
- `skills/do-work/actions/capture.md` (modify) — § Immutability Rule gains the mechanical-timestamp-repair exception, same commit as the tool
- `_dev/tests/prescribed-shell-scripts-behavior.sh` (modify) — lock-in cases from the Red-Green Proof
- `skills/do-work/scripts/repair-req-timestamps.sh` (modify, only if sharing requires extraction) — shared detection/derivation logic factored rather than duplicated

**Files I will NOT touch:**
- Any hook file — this tool is never wired into a hook (REQ constraint).
- `skills/do-work-board/tools/queue-kanban/**` — read-only frontmatter decision, pinned write-surface count.
- `do-work/archive/**` itself — the tool operates on consumer archives at runtime; this repo's archive is not a test fixture (build fixtures in scratch).

**Acceptance criteria (restated from REQ):**
- [ ] Replacement values derive only from git commit author times — file mtimes never consulted.
- [ ] Same detection predicate as REQ-246, shared rather than duplicated.
- [ ] Ordering clamp holds across the repaired file: created ≤ claimed ≤ completed, each no later than its introducing commit's time.
- [ ] The archive-immutability rule is amended in the same commit that ships the tool, at every statement of the rule the sweep finds.
- [ ] Never wired into any hook; dry-run default with --fix is acceptable latitude.
- [ ] `bash _dev/tests/maintainer-verify.sh` exits 0.

## Pre-Flight

**Git:** ✓ clean
**Tests baseline:** ✓ `bash _dev/tests/maintainer-verify.sh` exits 0 at the branch point (wave-1 integration tip)
**Dependencies:** ✓ Go 1.26.1, ShellCheck 0.11.0, `just`, Node, Chromium all present

*Checked by work action*
