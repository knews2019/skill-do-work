# Review: REQ-507 — hand the archive and commit tails to `finalize`

**Approve** — the finalization tail is genuinely owned by the CLI now, the four terminal paths are proven through the public command, and every acceptance criterion still holds against current `main`.
Route C | merge range `8e3dbf01e0660424965d79acb2e386b6604e4780..ad8bceb7aa0d0c63c230048b6a1f2dae1ef7ccb9` | re-verified at the current tree, no remediation commits

## What's built

Step 8 of the work action is now five judgment items (fold-first follow-ups, sweep and impact stamping, terminal judgment, release and lesson judgment, finalization intent), and Step 9 is one sentence that runs the `advance` continuation and reads its typed result. The archive move, release publication, staging, commit, provenance recording, verification and rollback all run inside the existing finalizer, reached through a new request-bound seam in `advance` that parses exactly one manifest and compares the manifest's request identity inside the finalizer's single decode, before any observable effect. What is still missing is a shipped description of the manifest's own JSON shape — see F-1; that gap predates this REQ.

## Decisions / risks for you

- **The implementation was accepted from a 4 September merge range that later REQs partly reworked.** The remediation builder re-derived all four acceptance criteria from live code and prose and returned zero commits. I checked the same four criteria independently against the current tree and reached the same verdict, so this is not a single-agent claim. The residual risk the orchestrator named in D-07 stands: a criterion weaker than its wording in a way neither tests nor predicates detect. My sweep did not find one.
- **The heavy-lane red that caused this pass was the assertion's fault.** `_dev/tests/staged-skills-contract.sh` no longer calls `assert_core_sibling_reference actions/work.md do-work-board`; the removal carries an attributing comment and is attributed to `c5dff3db`. `grep -c 'do-work-board/' skills/do-work/actions/work.md` returns 0, so the retired assertion would still fire today. The lane itself still needs to be re-run in this run's drain; that is not verified here.

## Findings

**Important:**

- The finalization manifest's JSON schema exists only in Go source (`internal/finalization/finalization_types.go`) and tests. Step 8 item 5 names the manifest's contents in English but never the JSON key names, the `transition` vocabulary, or the `provenance_mode` values; `finalize`'s own usage finding says only `--manifest requires one JSON file`, and no shipped action, doc, prime or help text carries an example. UR-098's constraint says a floor agent must complete a run from `advance` output plus the remaining prose, and that is not currently possible for this one artifact. **Not introduced by REQ-507** — the pre-REQ Step 9 prose described the manifest in the same English terms, so the gap arrived with REQ-498. — impact-user-visible → report only

**Minor:**

- Restatement Sweep hit, introduced by this diff: Step 9 lost the sentence `Check for git with git rev-parse --git-dir 2>/dev/null. If not a git repo, skip.` Only the heading suffix "(Git repos only)" still carries the condition, while `skills/do-work/actions/work-reference.md:383` still restates "actions/work.md Step 9 already skips its Commit Phase here". An agent reading Step 9 alone in a non-git project would run the advance continuation unconditionally. — impact-user-visible → report only
- The REQ's four-part write set names "deleted predicates" as a mandatory part, and that part is empty. At the base revision `8e3dbf01` no predicate pinned Step 8 or Step 9 prose: `_dev/tests/contracts/core-checks.sh`, `_dev/tests/contract-regressions.sh` and `_dev/tests/action-shell-blocks.sh` all return zero matches for `Step 8`, `Step 9`, `finalize --manifest` and `record-commit-hash`. The captured RED case claimed 19 such predicates; earlier REQs in the chain had already retired them, the `## Exploration` section recorded that before implementation, and Plan step 4 substituted structural ownership guards (+40 lines in `core-checks.sh`, 0 deletions). The constraint's intent holds. Recorded so the literal reading is not silently treated as satisfied. — impact-negligible → report only
- Pre-existing stale cross-references to Step 8 as an archiving, UR-closure or branch-deleting step, in files this REQ never declared: `actions/cleanup.md:39` ("Step 8 has already moved out every REQ that run finished"), `actions/cleanup.md:54` and `actions/work-reference.md:262` ("Step 8's UR-final check"), `actions/work-reference.md:414` ("Step 8 deletes the branch before Step 9 runs verify"), and `actions/stakeholder-answers.md:66` ("the standalone shape in `actions/work.md` Step 8 substep 6's table" — Step 8 has neither a substep 6 nor a table). Every one of these was already stale at `8e3dbf01`, so REQ-498's earlier move is the origin, not this REQ. — impact-negligible → report only

**Nit:**

- `skills/do-work/actions/work.md:466` and `:471` — the diff consumed the blank line before the `### Step 9` and `### Step 10` headings. CommonMark still parses both as headings and no linter enforces it here, but it breaks the file's own spacing convention. — impact-negligible → report only

## Requirements Checklist

- [x] Fold-First minting, sweep consolidation and impact stamping stay prose — delivered (`work.md:457-458`, Step 8 items 1-2)
- [x] Archive, release payload validation, staging, commit, provenance and verification run inside `finalize` driven by `advance` — delivered (`finalization_gate.go` → `finalization.FinalizeBound`; `advance_commands.go:55-56, 250-256`)
- [x] Mechanical prose of Steps 8 and 9 deleted — delivered (93 lines removed from `work.md`, 93 from `work-reference.md`)
- [x] Mechanical parts of both reference procedures deleted — delivered (`work-reference.md:775-789`; the changelog procedure keeps only release judgment, the commit procedure only manifest authorship and ordered-record consumption)
- [~] Predicates deleted — N/A in effect; none existed at the base revision. Structural ownership guards added instead (F-3)
- [x] Go tests for serial, worktree, completed-with-issues and already-green paths — delivered (`TestAdvanceFinalizationRunsTerminalPathMatrix`, four subtests)
- [x] Judgment stays prose; `advance` emits typed findings, never paragraphs — delivered (every refusal path returns a coded finding)
- [x] One step per REQ, no rewrite of `work.md` — delivered (the diff touches only the Step 8, Step 9 and two reference-procedure regions)
- [ ] The floor agent completes a run from `advance` output plus remaining prose — not established for the manifest artifact (F-1), and not a regression from this REQ

## Acceptance Testing

**Result: Pass**

- `GOMAXPROCS=2 go test -count=1 ./internal/lifecycleadvance ./internal/resultmodel` — exit 0, lifecycleadvance 18.520s, resultmodel 0.375s.
- `bash _dev/tests/contracts/core-checks.sh` — exit 0, final lines `shared principles loads and mutations passed; near_identical_cross_file_pairs 0` and `core-checks contract probes passed.`
- Live wiring, read-only: `do-work-cli.sh --repo-root . --format json advance REQ-507` returns `phase: agent judgment: review`, `phase_kind: agent_judgment`, and one missing-evidence entry naming the `Review` section of this REQ's own working file. That is the phase immediately upstream of `finalize`, and it confirms the classifier reaches the new terminal branch through the real repository rather than only through fixtures.
- I read the assertions rather than trusting the test names. The refusal matrix asserts an unchanged whole-tree digest, unchanged `HEAD`, and the absence of a `hostile-marker` file on all seven cases. The terminal matrix asserts the archived terminal status, that `implementation.txt` never enters the finalization commit allowlist on the supplied-commit path, that the supplied hash actually owns the implementation bytes, and that `VERSION` and `CHANGELOG.md` are byte-unchanged with no `release_at:` on the no-release path.
- Not run: `_dev/tests/maintainer-verify.sh` and the heavy lanes, per the review brief. The staged-skills lane in particular is still owed by this run's drain.

## Suggested Additional Testing

1. Re-run the `staged-skills` heavy lane at the current tree. The 4 September red is explained by a since-deleted assertion, but the explanation is a code reading, not an execution.
2. Exercise the finalize path in a non-git project, or confirm the intended behavior there. Step 9's explicit git-repo skip is gone (F-2) and the finalizer's own refusal is what a floor agent would hit instead.
3. Have a fresh agent with no session context author a finalization manifest from Step 8 item 5 plus `advance` output alone. That is the direct test of UR-098's floor-agent constraint and the fastest way to size F-1.
4. Run the `do-work-cli-integrations` heavy lane. It is the lane that exercises the public command seam this REQ introduced.

## Scores (on the record — not the headline)

**Overall: 96%**

| Dimension | Score | Notes |
|-----------|-------|-------|
| Requirements | 96% | All four acceptance criteria hold at the current tree; the "deleted predicates" part was vacuous with a recorded reason |
| Code Quality | 93% | Binding placed before any effect, names read as two words and grep cleanly, `FinalizeBound`'s comment states why; two heading-spacing nits |
| Test Adequacy | 96% | Four terminal paths plus seven no-mutation refusals through the public command, with tree-digest, HEAD, hostile-marker and release-byte assertions; text/JSON parity and record ordering covered |
| Scope | 100% | 12 files touched, 12 declared; qualification found no undeclared touch and no unused declaration |
| Risk | None | The new dispatch branch is keyed on the `finalize` phase string, which no other phase uses, and the phase matrix covers every phase |
| Acceptance | Pass | Both focused suites green, contract suite green, live advance projection correct |

## Follow-ups created

None (5 findings report only)

## Self-Validation

Re-examining my own pass: I verified the two Go test files by reading their assertions, not their names, which is the failure mode `lessons-do-work-cli.md` records under `smoke-vs-characterization`. I did not run the `finalization` package (67s, already green in the builder's run) or the heavy lanes, and I said so rather than scoring them. I checked each stale restatement against the base revision `8e3dbf01` before attributing it, which is what separated F-2 (introduced here) from F-5 (pre-existing). No new issue surfaced in this pass.

*Reviewed by review-work action, orchestrated mode, 2026-09-05T10:23:13Z*
