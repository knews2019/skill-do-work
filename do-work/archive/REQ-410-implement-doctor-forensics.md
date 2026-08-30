---
id: REQ-410
title: 'Implement doctor, deterministic forensics, and metadata repairs'
status: completed-with-issues
claimed_at: 2026-08-30T20:45:17Z
completed_at: 2026-08-30T22:01:41Z
commit: 210d1459
created_at: 2026-08-29T20:28:26Z
route: C
estimate:
  p50_active_minutes: 100
  confidence: low
  calculated_at: 2026-08-30T20:45:17Z
  basis:
    - Route C
    - 18-file write set
    - 10 new files
    - 6 subsystems involved
    - 4 acceptance criteria
    - dependency depth 4
    - persistence changes
    - cross-route regression gates
    - full-suite verification
user_request: UR-081
domain: general
prime_files: [_dev/primes/prime-action-files.md]
tdd: true
suggested_spec:
depends_on: [REQ-409]
maintenance: false
impact: impact-user-visible
effort_estimate: effort-substantive
related: [REQ-406, REQ-407, REQ-408, REQ-409, REQ-411, REQ-412, REQ-413, REQ-414, REQ-415, REQ-416, REQ-417, REQ-418, REQ-419, REQ-420]
batch: go-no-llm-command-platform
kb_status: pending
---

# Implement Doctor, Deterministic Forensics, and Metadata Repairs

## What
Create `doctor` and move deterministic forensic checks and safe metadata repairs into Go.

## AI Execution State (P-A-U Loop)
- [x] **[PLAN]:** Read the REQ, UR, action/CLI primes and lessons, prior cleanup implementation, crew rules, forensic action/guide, blank scanner, timestamp repairers, and shared Go seams; validated a four-task exact-scope plan.
- [x] **[APPLY]:** Added the doctor command and tests RED-first, then implemented shared damage/line/recovery evidence, deterministic scans, guarded timestamp repair, registration, and delegation within the declared 18-file boundary.
- [x] **[UNIFY]:** Reviewed all 18 files and ran focused/race/vet/full/Windows/contract/reference/diff/scope checks; no debug artifacts, dependency drift, or unrelated queue edits entered the implementation.

## Detailed Requirements
- Implement `do-work-doctor` with deterministic forensics, blanked-record detection, timestamp audit/repair, collision checks, and remediation output.
- Reuse the typed finding/result schema so an LLM caller receives evidence and exact next commands without rescanning.
- Preserve the distinction between read-only diagnosis, provably safe repair, and explicitly destructive recovery.
- Support text/JSON, dry-run where repairing, optional commit, and shared Git guards.

## Constraints
- Missing or failed canonical tooling must stop the operation with actionable output rather than trigger prose mutation.

## Dependencies
Depends on REQ-409 (cleanup safety and repair classifications).

## Builder Guidance
Certainty level: Firm. Characterize every existing deterministic forensic utility before migration.

## Red-Green Proof
**RED prompt/case:** Run doctor over fixtures containing a blanked REQ, wrong timestamp, ID collision, and clean control.
**Why RED now:** These checks are split between shell helpers and LLM-directed action prose and do not share an actionable result schema.
**GREEN when:** Doctor classifies each fixture deterministically, performs only authorized safe repairs, and emits matching text/JSON evidence plus exact remediation/verification commands.
**Validation:** User confirmed via the supplied implementation plan.

## Full Context
See `do-work/user-requests/UR-081/input.md` for complete verbatim input.

---
*Source: UR-081 (Replace LLM bookkeeping and shipped utility logic with a Go command platform)*

## Triage

**Route: C** - Complex

**Reasoning:** Doctor unifies multiple forensic sources, safe metadata mutation, Git transaction boundaries, stable text/JSON evidence, and compatibility/action delegation across several existing utilities.

**Planning:** Required

## Plan

1. **Preserve diagnostic evidence in shared models.** Extend `requestmodel` with source-line evidence for effective top-level fields and extend `repositorymodel` so blank/malformed REQ and UR files remain typed snapshot evidence instead of collapsing into warning strings. Preserve exact bytes, paths, parse failures, filename/frontmatter claims, and symlink/exclusion boundaries.
2. **Build deterministic doctor scans.** Add a `doctor` package that produces stable `resultmodel.CommandResult` findings from one repository snapshot for the mechanical forensics checks: stuck work, hollow completions, missing qualification, orphaned URs, scope overlap, failed-without-follow-up, stale holds, Git divergence, stranded terminal records, invalid statuses, timestamp defects, blanked records, and collision evidence.
3. **Add read-only recovery evidence and guarded timestamp repair.** Extract cleanup's full-history blank recovery into a reusable evidence function while leaving restoration consent in `cleanup --restore-blanked <path>`. Add `doctor --repair-timestamps`, dry-run, and optional commit using lossless scalar edits, atomic publication, and the shared Git transaction boundary.
4. **Register and delegate.** Register `doctor`, reduce the natural-language forensics action to the canonical deterministic command plus its irreducibly judgment-based recurring-correction phase, and update the guide and CLI prime. Leave flat Just recipes and retained shell shims to REQ-419/REQ-420.

**Requirement coverage:** Tasks 1–2 cover deterministic forensics, blank/collision evidence, and exact remediation commands; task 3 covers read-only/safe/destructive classification, dry-run, commit, rollback, and timestamp repair; task 4 covers the public command and action compatibility boundary. Every task maps to a Detailed Requirement.

**Testing:** RED-before-GREEN fixtures cover unknown command registration, a combined clean/blank/timestamp/collision repository, text/JSON parity, strict timestamp shape preservation, dirty/untracked/commit guards, rollback/risk, full-history provenance, malformed discovery, and non-Git diagnosis. Then run focused/race/vet/full-module, Windows compile-only, qualification, scope drift, contract regressions, and the unpiped canonical maintainer gate.

*Generated by Plan agent; validated by orchestrator*

## Exploration

The existing fourteen-check action mixes three ownership classes. Core mechanical checks can move into `doctor`; recurring-correction theme grouping remains an LLM judgment over untrusted lesson prose; release/UI invariants remain owned by the separate, non-importable `queue-kanban verify` binary. Reimplementing the latter would create two authorities, while parsing its prose cannot satisfy the typed-evidence contract.

Timestamp repair has a strict compatibility boundary beyond what generic `time.Parse` accepts: effective top-level last-key-wins fields only; two-minute future skew; `created_at <= claimed_at <= completed_at`; repairable date/naive-whole-second/Z-whole-second shapes; and byte-identical refusal for offsets, fractions, calendar-invalid values, non-ASCII padding, nesting, symlinks, and unclosed frontmatter. CRLF, BOM, comments, modes, and unrelated bytes must survive. The new global Git guard supersedes the legacy dirty/untracked mtime repair path for this command.

Blank detection scans queue, working, and recursive archive while excluding `assets/` and symlinks. Recovery must distinguish the blob source commit from the later `record commit hash` implementation provenance and use full history. Collision findings consume `CollisionEntries` before canonical ID lookup so filename-only claims never degrade into generic missing evidence.

*Generated by Explore agent; reconciled by orchestrator*

## Scope

**Files I will touch (18; 8 new):**
- `skills/do-work/tools/do-work-cli/internal/doctor/doctor_commands.go` (new)
- `skills/do-work/tools/do-work-cli/internal/doctor/doctor_commands_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/doctor/doctor_scan.go` (new)
- `skills/do-work/tools/do-work-cli/internal/doctor/doctor_scan_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/doctor/doctor_timestamps.go` (new)
- `skills/do-work/tools/do-work-cli/internal/doctor/doctor_timestamps_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/doctor/doctor_repair.go` (new)
- `skills/do-work/tools/do-work-cli/internal/doctor/doctor_repair_test.go` (new)
- `skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/requestmodel/request_model.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/requestmodel/request_model_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_git.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_git_test.go` (modified)
- `skills/do-work/actions/forensics.md` (modified)
- `skills/do-work/docs/forensics-guide.md` (modified)
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (modified)

**Files I will NOT touch:** `queue-kanban`, timestamp/blank scanner shell implementations, SessionStart/hooks, Just recipes, result/runtime/transaction package contracts, cleanup mutation planning, unrelated queue files, and REQ-427.

**Acceptance criteria:**
- [x] Default `doctor` is byte-for-byte read-only and emits stable, sorted text/JSON evidence for the requested mechanical checks, including filename-only collisions and blank recovery provenance.
- [ ] Timestamp diagnosis preserves the legacy top-level/effective-key/skew/ordering grammar and explicitly refuses unsupported shapes without rewriting them.
- [x] `--repair-timestamps` changes only provable clean tracked targets; dry-run is non-mutating, `--commit` requires an empty index, and rollback/risk output uses the shared transaction result.
- [x] Blanked restoration remains explicit cleanup consent, while collision/manual findings carry exact next and verification argv.
- [ ] Natural-language forensics delegates deterministic phases to canonical tooling and stops on tooling failure; recurring lesson judgment and board-owned invariants remain correctly routed.
- [x] Focused, race, vet, full-module, cross-platform, contract, and canonical repository gates pass.

## Decisions

### D-01: Keep default doctor strictly read-only

**Decision:** DECIDE & STATE — diagnosis never mutates; `--repair-timestamps` is the only doctor mutation mode, and blank restoration remains an exact-target cleanup operation.

**Reasoning:** This preserves the established forensic contract and keeps destructive recovery behind REQ-409's explicit consent token.

### D-02: Preserve existing diagnostic ownership boundaries

**Decision:** DECIDE & STATE — doctor owns core mechanical repository findings; the action retains recurring-correction theme judgment; `queue-kanban verify` retains release/UI invariants.

**Reasoning:** Check 10 is semantic judgment, while duplicating the separate board verifier would create two deterministic authorities and parsing its prose would not yield stable typed evidence.

### D-03: Share evidence seams without sharing cleanup application

**Decision:** DECIDE & STATE — expose pure full-history recovery evidence from cleanup and source-line field evidence from requestmodel, but never call cleanup's plan applier from doctor.

**Reasoning:** Doctor needs authoritative provenance and blame locations while REQ-409's known operation-group defects must remain outside diagnosis and timestamp repair.

### D-04: Use one explicit repair grammar

**Decision:** DECIDE & STATE — support `doctor [--repair-timestamps] [--dry-run] [--commit]`; dry-run and commit are valid only with repair mode and are mutually exclusive.

**Reasoning:** The grammar makes mutation intent reviewable and prevents flags from suggesting a repair occurred during default diagnosis.

### D-05: Treat incomplete inspection as non-clean

**Decision:** DECIDE & STATE — contained read, walk, Git-history, or blame failures produce actionable findings/skips and a non-clean outcome; they never disappear into a clean result.

**Reasoning:** A diagnostic is only trustworthy when every reported clean surface was actually inspected.

### D-06: Honor the new Git guard over legacy mtime mutation

**Decision:** DECIDE & STATE — doctor may report legacy mtime evidence but refuses to repair dirty or untracked targets; only clean tracked targets can use blame-derived timestamps.

**Reasoning:** UR-081 explicitly requires dirty-target refusal and exact rollback, superseding the older hook repairer's broader mutation authority for this new command.

### D-07: Convert path coordinate systems only at boundaries

**Decision:** DECIDE & STATE — retain discovery paths relative to `do-work/`, and add the `do-work/` prefix only when crossing into repository-relative Git, result, or filesystem operations.

**Reasoning:** Repository snapshots deliberately model the contained tree, while Git transactions and user-facing argv must name paths from the repository root; mixing the two silently targeted nonexistent files during RED repair fixtures.

### D-08: Validate provenance and implementation paths before inspection

**Decision:** DECIDE & STATE — accept only 7–40 hexadecimal recorded commit hashes and only relative, contained implementation-summary paths before issuing Git or filesystem probes.

**Reasoning:** Archived Markdown is evidence, not authority to construct arbitrary revisions or absolute/escaping pathspecs.

### D-09: Make unsupported timestamp repair shapes explicit

**Decision:** DECIDE & STATE — report unsupported effective timestamp shapes as byte-identical repair refusals even when no rewrite is attempted.

**Reasoning:** Silence would imply those fields were audited as safely repairable; the refusal preserves the legacy compatibility boundary and makes partial inspection visible.

### D-10: Derive implementation ownership only from file-list bullets

**Decision:** DECIDE & STATE — accept only path-led backticked bullets inside `## Implementation Summary` and deduplicate each path per REQ before scope ownership aggregation.

**Reasoning:** Prose code spans and repeated handback entries are not independent ownership evidence; treating them as paths produced actionable-looking but false forensic findings on the live archive.

### R-D-01: Preserve legacy archived user-request inputs as historical evidence

**Decision:** REVIEW REMEDIATION — retain frontmatter-less `archive/UR-NNN/input.md` files in discovery, but do not classify that valid legacy input class as a damaged record or incomplete-inspection warning.

**Reasoning:** Those inputs are source history, not standalone REQ/UR records; treating them as destroyed records produced a false recovery emergency.

### R-D-02: Carry commit metadata through post-repair verification

**Decision:** REVIEW REMEDIATION — retain the exact repair commit SHA and revert argv outside the generic command result so a failed rediscovery after commit becomes committed-state risk with an exact recovery command.

**Reasoning:** Once a commit exists, verification failure is materially different from an uncommitted repair failure and must not hide the durable mutation.

### R-D-03: Keep compatibility prose non-executable

**Decision:** REVIEW REMEDIATION — mention the retained blank scanner only as a non-executable compatibility surface; doctor is the sole mechanical authority used by forensics.

**Reasoning:** Executable legacy recipes would preserve two diagnostic authorities after migration.

### R-D-04: Advance manual findings with exact-path history inspection

**Decision:** REVIEW REMEDIATION — when a manual finding has affected paths and no more specific action, replace an identical doctor rerun with `git log --full-history -- <exact paths>` while retaining JSON doctor as verification.

**Reasoning:** A next command must advance investigation and remain directly tied to the finding's evidence.

## Pre-Flight

**Repository state:** Clean outside `do-work/`; the only pre-existing local edit is the user's unrelated `do-work/queue/REQ-427-confirm-go-version-floor.md`, fingerprinted before implementation and excluded from scope. Reservation markers and this session checkpoint are queue metadata, not implementation inputs.

**Baseline:** `git diff --check` and uncached `go test ./...` in `skills/do-work/tools/do-work-cli` pass before implementation.

**Dependencies:** REQ-409 is archived `completed-with-issues`, which satisfies the terminal dependency gate. Doctor will reuse only its pure recovery evidence and explicit cleanup consent boundary; REQ-409's four queued operation-group defects are outside this command's application path.

## Implementation Summary

**Files changed:**
- `skills/do-work/tools/do-work-cli/internal/doctor/doctor_commands.go` (new)
- `skills/do-work/tools/do-work-cli/internal/doctor/doctor_commands_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/doctor/doctor_scan.go` (new)
- `skills/do-work/tools/do-work-cli/internal/doctor/doctor_scan_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/doctor/doctor_timestamps.go` (new)
- `skills/do-work/tools/do-work-cli/internal/doctor/doctor_timestamps_test.go` (new)
- `skills/do-work/tools/do-work-cli/internal/doctor/doctor_repair.go` (new)
- `skills/do-work/tools/do-work-cli/internal/doctor/doctor_repair_test.go` (new)
- `skills/do-work/tools/do-work-cli/cmd/do-work-cli/main.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/requestmodel/request_model.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/requestmodel/request_model_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model_test.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_git.go` (modified)
- `skills/do-work/tools/do-work-cli/internal/cleanup/cleanup_git_test.go` (modified)
- `skills/do-work/actions/forensics.md` (modified)
- `skills/do-work/docs/forensics-guide.md` (modified)
- `skills/do-work/tools/do-work-cli/prime-do-work-cli.md` (modified)

**What was done:** Added canonical `doctor` diagnosis and guarded timestamp repair using one typed result, one repository snapshot, explicit full-history recovery provenance, strict timestamp compatibility evidence, and shared Git transactions. Registered the command and made the natural-language forensics action and guide delegate deterministic work to it while retaining judgment-only recurring lessons and board-owned verification.

## Discovered Tasks

None.

## Qualification

**Attempt 1: Failed judgment check 6 (data flows).** Mechanical qualification and exact scope drift passed, and all 18 files are substantive and wired. However, the real-repository default-doctor smoke emitted nonsensical `SCOPE-OVERLAP` paths such as `(2,557) +`, `(modified) —`, and `See`, with repeated copies of one REQ ID. `implementationPaths` scans every backticked span in the whole Implementation Summary instead of only the path-led `Files changed` list and does not deduplicate per REQ. The command remained byte-for-byte read-only, but its core diagnostic evidence was not trustworthy enough to qualify. Re-qualification requires a live-shape regression and corrected extraction/deduplication.

**Attempt 2: Passed.** The new live-shape regression proves prose code spans and duplicate file mentions cannot become ownership evidence. Mechanical qualification and scope drift remain green; all 18 files are substantive, every requirement maps to an implemented/tested path, same-package static-reference warnings are expected, and a real-repository doctor run stayed byte-identical while emitting only plausible contained paths with one claim per REQ/path.

**Review remediation: Passed.** Seven RED-to-GREEN regressions cover the initial review's six Important findings and one Minor finding: legacy archived inputs, complete stuck/hollow/stale predicates, quoted ASCII timestamp padding, committed rediscovery risk/revert evidence, doctor-only action delegation, exact-path advancing manual findings, and the dead test command. Independent focused/race/vet/full-module and Windows compile gates pass; qualification and exact scope drift remain green. Live text and JSON doctor runs both returned the expected findings outcome, left the repository diff byte-identical, and emitted no `UR-003/input.md` false positive.

## Testing

**Tests run:** focused doctor/requestmodel/repositorymodel/cleanup/cmd tests; focused race tests; `go vet ./...`; uncached full-module tests; Windows doctor and atomicfile compile-only; qualification; scope drift; `git diff --check`; contract and shipped-reference regressions; real-repository read-only smoke; unpiped `bash _dev/tests/maintainer-verify.sh`.

**Result:** ✓ All required tests and the canonical repository gate pass. The optional strict browser lane remained in its documented default-skipped state because no browser is configured.

**Red-green validation:**
- Captured caller seam: `doctor` initially returned `UNKNOWN-COMMAND`/exit 2 → ✓ registered command returns one shared text/JSON result.
- Combined fixture: collision, blank, and timestamp evidence initially produced no findings → ✓ independently sorted typed findings with exact next/verification argv.
- Timestamp repair: named fixtures initially produced no plans → ✓ strict top-level/effective-key repair and refusal cases pass through shared Git guards.
- Re-qualification: realistic Implementation Summary prose and duplicate path mentions initially produced `[src/doctor.go src/doctor.go src/doctor_test.go (2,557) + See (modified) — not/a/file-list-bullet.go]` → ✓ only the two actual file-list paths remain.
- Review remediation: each initial Important finding has a named failing regression → ✓ all seven remediation regressions pass, including exact committed-state revert and the production legacy archive class.

**New tests added:** four doctor test files plus focused source-line, damaged-record, and reusable full-history recovery coverage in the shared packages.

*Verified by work action*

## Review

**Overall: 50%** | 2026-08-30T21:35:15Z

| Dimension | Score |
|-----------|-------|
| Requirements | 62% |
| Code Quality | 60% |
| Test Adequacy | 55% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Fail |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- Valid legacy frontmatter-less archived `input.md` is falsely classified as a destroyed/unrecoverable blanked record — impact-user-visible → remediation required
- Stuck-work, hollow-completion, and stale-hold checks only partially preserve their canonical predicates and evidence bands — impact-user-visible → remediation required
- Quoted timestamps with supported ASCII padding are incorrectly refused by the repairability predicate — impact-user-visible → remediation required
- Post-commit rediscovery failure is downgraded to ordinary failure without committed-state risk or exact revert evidence — impact-user-visible → remediation required
- Delegated forensics prose still instructs readers to run retired timestamp and blank-scanner mechanics — impact-user-visible → remediation required
- Manual/inspection findings can lack affected paths, duplicate invalid-status evidence, and offer only an identical doctor rerun instead of actionable next steps — impact-user-visible → remediation required

**Minor findings:** 1 (a dead overwritten `git add` test command)
**Acceptance:** Fail — live doctor falsely diagnoses a valid legacy UR source as destroyed, and multiple required forensic/result contracts remain partial.
**Suggested testing:** 5 grouped items
**Follow-ups created:** None; **sweeps appended to:** None

*Reviewed by review-work action — remediation pending*

## Re-Review

**Overall: 50%** | 2026-08-30T22:01:41Z

| Dimension | Score |
|-----------|-------|
| Requirements | 80% |
| Code Quality | 78% |
| Test Adequacy | 82% |
| Scope | 100% |
| Risk | Low |
| Acceptance | Fail |

**Important findings (terminal, after the single allowed remediation):**
- Unsupported but parseable timestamp shapes can still become ordering anchors and clamp a supported successor, violating byte-identical refusal and provably safe repair — impact-user-visible → follow-up REQ-434
- The delegated forensics action requires counts and remedies absent from doctor while forbidding rescans, and retained consumers point at deleted numbered checks — impact-user-visible → follow-up REQ-435

**Acceptance:** Fail — all seven initial remediation regressions pass, but the mixed timestamp-anchor repair and end-to-end forensics delegation contracts remain incomplete. The raw dimensional average is 85%; Acceptance Fail caps the terminal score at 50%.

**Follow-ups created:** REQ-434, REQ-435; **sweeps appended to:** None (fold-first scan found no matching pending sweep).

*Re-reviewed by a fresh review-work agent; no further review remediation is permitted*

## Remediation

The single Route C review remediation corrected all six initial Important findings and the Minor test defect, with seven named RED-to-GREEN regressions. Independent focused, race, vet, full-module, Windows compile, qualification, scope-drift, live text/JSON, and canonical maintainer gates pass. The fresh re-review nevertheless found two remaining user-visible acceptance defects outside those initial regressions. Per the one-remediation limit, REQ-410 closes as `completed-with-issues`; REQ-434 owns unsupported timestamp ordering anchors and REQ-435 owns the incomplete doctor/forensics delegation contract.

## Lessons Learned

**What worked:** One typed snapshot/result path, exact Git guards, live-corpus acceptance, and named adversarial regressions removed duplicate mechanical authorities and closed the initial false-positive and committed-state defects.

**What didn't:** Migrating the producer without tracing every downstream report field and reference left the action contract half-landed. Treating “parseable” as “eligible for repair ordering” also let a diagnosis-only timestamp influence a mutation.

**Worth knowing:** A canonical-tool migration is complete only when every consumer can produce its required output from that tool and all legacy anchors are swept. For repair logic, comparison eligibility must be the same supported-shape predicate as mutation eligibility; unsupported evidence must remain observational and byte-identical.

## Orientation

[MAP CHANGED] The do-work CLI now has a canonical `doctor` command for typed forensic diagnosis, blank-record recovery evidence, collision checks, and guarded timestamp repair, with the forensics action delegating its mechanical phase. The feature is archived with two known contract gaps tracked by REQ-434 and REQ-435; the action-file and CLI primes still resolve to the canonical subsystem paths.
