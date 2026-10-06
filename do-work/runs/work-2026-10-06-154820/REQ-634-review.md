## Review

**Overall: 97%** | 2026-10-06T16:03:40Z

| Dimension | Score |
|-----------|-------|
| Requirements | 100% |
| Code Quality | 95% |
| Test Adequacy | 95% |
| Scope | 98% |
| Risk | Low |
| Acceptance | Pass |

**Important findings (each with its recorded impact token — this is the durable audit record the judgment mandates):**
- None

**Minor findings:** None
**Acceptance:** Pass — reviewer reproduced the captured RED (pre-change code: exit 2, `PARSE-FAILED`, no owner row) and GREEN (exit 0, `REQ-503<TAB>good-file.txt` on stdout, `ASSOCIATION-SUMMARY-UNPARSED` naming REQ-502 on stderr, no `PARSE-FAILED`) on a scratch repository; the finding also appears under `--format json`; corehelpers Go tests and core-checks.sh pass.
**Suggested testing:** 3 items
**Follow-ups created:** None (1 findings report only)

*Reviewed by review-work action*

### Review detail (orchestrated report)

## Review: REQ-634

**Approve** — One REQ file with a broken Implementation Summary now claims nothing and is named in a warning, and association for every other REQ completes with exit 0.
Route A | e86b314d merged as b9e04c02 (range 09c5ddf7..b9e04c02)

### What's built
- The association walk records a REQ file whose Summary fails to parse, skips it, and keeps walking. The `PARSE-FAILED` exit-2 branch is gone.
- Each skipped file becomes one `ASSOCIATION-SUMMARY-UNPARSED` warning finding with the repository-relative path, the parser message, and exact next/verification argv. In shim text mode the runtime prints it to stderr; stdout owner rows are unchanged.
- commit.md, inspect.md, and the "What `associate` settles" list state the new rule. The release entry and the lessons entry (requirements 5 and 6) are not in the merge range by design; the hand-back assigns them to the orchestrator's finalization.

### Decisions / risks for you
- **Finalization still owes requirements 5 and 6.** The changelog entry, version bump, lessons satellite entry, and lessons-index token refresh must land in the finalization commit. This review scored them as outside the builder range, not as delivered.
- **A skipped in-flight REQ loses its claims.** If the REQ being committed has an unparseable Summary, its own files go to commit Step 4 grouping. The REQ's Constraints accept this.
- **Builder decisions D-01 to D-04 are sound.** A third return value, a fix-hint entry the commands matrix contract requires, evidence ending "this REQ file claims no paths", and the JSON visibility the maintainer accepted at plan approval.

### Findings

**Important:**
- None

**Minor:**
- None

**Nit:**
- F1. In shim text mode, stderr also carries one pre-existing `ASSOCIATION-UNOWNED` warning per unowned candidate (observed in the scratch run and in core-checks output), so in a real commit the new warning sits among those lines. The prose names the code, so an agent can find it. Pre-existing behaviour, not introduced here. — impact-negligible → report only

### Requirements Checklist

- [x] R1. Parse error on one REQ file means it claims no paths, walk continues, skipped files returned — delivered (`inventory.go` callback appends `UnparsedSummaryRecord`, returns nil; third return value)
- [x] R2. `PARSE-FAILED` branch deleted; one warning per skipped file, code `ASSOCIATION-SUMMARY-UNPARSED`, repo-relative path, parser message as evidence; outcome success; printed by the runtime, not a direct stderr write — delivered
- [x] R3. Both lock-ins flipped (core-checks `associate_unmatched` probe, shim subtest) asserting exit 0 / owner row / no `PARSE-FAILED` / skipped file named; new unit test on `AssociateProjectPaths` — delivered
- [x] R4. commit.md and inspect.md sentence rewritten; new bullet in prescribed-shell-primitives.md — delivered
- [ ] R5. Changelog entry and version bump — not in range; owned by orchestrator finalization (not scored against the builder)
- [ ] R6. Lessons satellite entry and index refresh — not in range; owned by orchestrator finalization (not scored against the builder)
- [x] C1. No parser grammar change: `checks.go` untouched in the range; `allBacktickedPaths`, `firstBacktickedPaths`, `qualificationSummaryEntries` unchanged — verified
- [x] C2. Qualify and scope-drift stay strict — verified (`qualify.sh` unmatched case and core-checks scope-drift probe still pin the errors; both green)
- [x] C3. Multi-path bullet contract unchanged — verified (test and probe untouched; the walk still uses `allBacktickedPaths`)
- [x] C4. In-flight `working/` rule unchanged, no `ASSOCIATION-FOUND-INFLIGHT` marker — verified
- [x] C5. `PARSE-FAILED` gone — verified (remaining hits are the historical changelog line, negative assertions, and an export-ignored dated snapshot report)
- [x] C6. Stdout rows unchanged — verified (`REQ-503<TAB>good-file.txt`, no extra rows)
- [x] UR D2 deviation: warning also visible under `--format json` — verified

**Restatement Sweep.** The diff redefines the associate exit contract. Grep for `PARSE-FAILED`, "unmatched backtick", and associate near "exit 2" across the repository (excluding `do-work/`): no stale restatement remains. `skills/do-work/CHANGELOG.md:390` is a dated history line. The "unmatched backtick" hits in `checks.go`, `finalization_discovery.go`, `qualify.sh`, and the scope-drift probe belong to the strict parsers the REQ keeps. `ai-reports/2026-09-06_2113_do-work-cli-guide/index.html` embeds source at a fixed past commit and is export-ignored. prescribed-shell-primitives.md line 47 (stdout rows, exit 2 on quarantine resolution failure) still agrees.

### Acceptance Testing

**Result: Pass**
- `go test -count=1 ./internal/corehelpers/` from the CLI directory at HEAD b9e04c02: ok, 8.0s.
- `bash _dev/tests/contracts/core-checks.sh`: exit 0, "core-checks contract probes passed."
- Scratch repository outside the repo (`…/scratchpad/review-fixture`): committed archive REQ-502 (bullet with an unclosed second span) and REQ-503 (claims `good-file.txt`), uncommitted `good-file.txt` and `unowned-file.txt`.
  - RED, pre-change tree extracted with `git archive 09c5ddf7 skills/do-work`: `start` exit 0, `associate` exit 2, stdout `PARSE-FAILED: unmatched backtick in Implementation Summary`, no owner row.
  - GREEN, current scripts/protected-inventory.sh: `start` exit 0, `associate` exit 0, stdout exactly `REQ-503<TAB>good-file.txt`, stderr includes `finding ASSOCIATION-SUMMARY-UNPARSED [warning]: unmatched backtick in Implementation Summary; this REQ file claims no paths` with `paths: do-work/archive/UR-301/REQ-502-unmatched-summary.md`, zero `PARSE-FAILED`.
  - JSON, after a fresh `start`, `do-work-cli.sh --format json protected-inventory associate`: exit 0, outcome success, finding `ASSOCIATION-SUMMARY-UNPARSED` severity warning with the REQ-502 path, evidence, and non-empty `next_argv`/`verification_argv`; stderr empty.
- The captured Red-Green Proof matches exactly.

### Suggested Additional Testing

- Two unparseable REQ files in one walk: confirm two warnings, one per file, and that both files' paths go unclaimed.
- An archived non-terminal REQ with a broken Summary: confirm no warning (the status gate runs before the parse), so only counted REQs are reported.
- One real `do-work commit` run in a repository with a broken archived REQ: confirm the agent reports the skipped file in its summary as commit.md now instructs.

### Scores (on the record — not the headline)

**Overall: 97%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 100% | R1-R4 and all constraints delivered; R5-R6 belong to finalization, outside the range |
| Code Quality | 95% | Small, condition-keyed change; old message-substring branch removed; comments explain why |
| Test Adequacy | 95% | RED/GREEN at unit, shim, and shell-probe layers; cross-REQ test changes traced to REQ-603 and REQ-539 |
| Scope | 98% | commands.go and commands_test.go outside the declared list, justified in D-02 and signature updates |
| Risk | Low | Exit contract loosens from 2 to 0 for one case; unreadable REQ files still fail loudly |
| Acceptance | Pass | RED and GREEN reproduced end to end; JSON finding verified |

### Follow-ups created
- None (1 findings report only)
