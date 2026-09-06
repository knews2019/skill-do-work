# REQ-592 Exploration — seal the do-work tree into both fast gate stages

*Explore agent, read-only, run against HEAD rather than against the audited commit. Every claim below was re-verified in this checkout.*

## Claims verified

### C1. fast_stage_evidence.go skips queueStatePrefix in BOTH its tracked and untracked seal loops, so a do-work/-only change does not invalidate a fast stage.

**Holds: True**

HOLDS. skills/do-work/tools/do-work-cli/internal/heavyverification/fast_stage_evidence.go:195 (tracked loop, inside workingTreeSeals, `if path == "" || strings.HasPrefix(path, queueStatePrefix) { continue }`) and :223 (untracked loop, same guard). queueStatePrefix = "do-work/" is defined at skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_run.go:31. Both skips run BEFORE the `stageCovered` test, so no coverage rule can override them today. Confirmed the only two uses in that file: grep shows queueStatePrefix at fast_stage_evidence.go:195,223 and nowhere else in the fast path.

### C2. The do-work-cli fast stage genuinely reads do-work/ — TestDiscoverRepositoryAcceptsProductionLegacyArchiveInputClass reads and byte-checks do-work/archive/UR-003/input.md.

**Holds: True**

HOLDS. skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model_test.go:397-417. Line 398 builds the path as filepath.Join("..","..","..","..","..","..","do-work","archive","UR-003","input.md") — six levels up from internal/repositorymodel lands exactly on the repo root. Line 406-408 asserts `len(productionBytes) != 5608` with message `production legacy fixture changed size: got %d bytes`. `wc -c do-work/archive/UR-003/input.md` = 5608. The test carries no testing.Short() guard, and I ran it under the stage's own flags: `go test -short -count=1 -run TestDiscoverRepositoryAcceptsProductionLegacyArchiveInputClass ./internal/repositorymodel/` => ok. Stage argv confirmed at _dev/tests/fast-stages.json (do-work-cli-fast-tests: run-go-tests-with-budget.sh skills/do-work/tools/do-work-cli -short ./...) and wired at _dev/tests/maintainer-verify.sh:777-778.

### C3. The queue-kanban fast stage genuinely reads do-work/ — board_live_test.go, durations_test.go, citations_test.go build the board from the real do-work/ tree.

**Holds: True**

HOLDS, with one correction about citations_test.go. skills/do-work-board/tools/queue-kanban/board_live_test.go:16-45 defines liveBoard/liveRepoRoot/liveBoardAt, which resolve the real repo root (walk.go:46 resolveRepoRoot) and call buildBoard against it; used at board_live_test.go:68,86,98,114,188,219,243. durations_test.go:269-279 TestLiveArchiveDurationsMatchTheCalibratedFigures builds the live board and asserts figures pinned to this archive (durations_test.go:214-253: 195 samples, 2026-07-31 median 2.5 min / 2 completed / 1 kept, 2026-08-15 25 completed / median 19.6 min). generate_test.go:703-727 generateLiveSiteInDir builds the live board and is used by ~20 generate_test.go tests. CORRECTION: of the citations_test.go live tests, only TestGeneratedBoardMarkdownFileShipsTheTicketMentionIndex (citations_test.go:1478) runs in the fast stage; the other live-board citations tests are TestBrowserBehavior*/TestJavaScriptBehavior* (citations_test.go:132,267,861,955,1082,1140,1356), and the fast stage excludes those prefixes via DO_WORK_GO_TEST_EXCLUDE_PREFIXES=TestJavaScriptBehavior,TestBrowserBehavior (_dev/tests/maintainer-verify.sh:725). The claim still holds via board_live_test.go, durations_test.go, generate_test.go and citations_test.go:1478.

### C4. fast_stage_evidence_test.go has a case named `queue state changed` that expects `reused`.

**Holds: True**

HOLDS. skills/do-work/tools/do-work-cli/internal/heavyverification/fast_stage_evidence_test.go:182-190, inside TestFastStageReuseDecisionTable (declared at :164). Case name at :185, mutation writes do-work/queue/REQ-002-fixture.md at :187, expectation `LaneDispositionReused` + `laneReasonFingerprintMatch` at :189. Its comment at :183-184 calls it "The one declared exemption".

### C5. _dev/tests/fast-stage-reuse-behavior.sh has a case `queue state alone still reuses`.

**Holds: True**

HOLDS. _dev/tests/fast-stage-reuse-behavior.sh:143-145: line 143 writes do-work/queue/REQ-002-fixture.md into the fixture, line 145 asserts `expect_case 'queue state alone still reuses' no 0 'REUSED (fingerprint_match, recorded '` — i.e. stage did not run, exit 0, REUSED line present.

### C6. do-work/test-durations.tsv is gitignored and written by the stage itself, so it needs a NARROW exclusion or the stage will invalidate its own evidence every run.

**Holds: True**

HOLDS on the facts, CONDITIONAL on the fix shape. Gitignored: .gitignore:3 `/do-work/test-durations.tsv`; `git status --porcelain --ignored do-work` reports it as `!!`; it is the ONLY untracked-or-ignored path under do-work/. Written by the stage: _dev/tests/test-duration-log.sh:4 defaults DO_WORK_TEST_DURATION_LOG to $repo_root/do-work/test-durations.tsv, _dev/tests/maintainer-verify.sh:19 sets exactly that path, and _dev/tests/run-go-tests-with-budget.sh:99-111 appends one row per test file on every stage run. Self-invalidation is real because the recorded fingerprint is the PRE-run one (maintainer-verify.sh:148-154 decides, :191 records that same fingerprint), so a seal covering the log can never match on the next run. CONDITIONAL PART: today the untracked-ignored branch at fast_stage_evidence.go:227 already skips it (`!stageCovered && ignored => continue`), so if the fix only empties non_stage_coverage the exclusion is not strictly required. Under the fix I recommend below (do-work/ becomes queue-kanban's declared COVERAGE), stageCovered becomes true, that skip stops applying, and the explicit narrow exclusion becomes mandatory. It is also the ONLY do-work path the gate writes itself — grep of maintainer-verify.sh for repo_root/do-work returns only line 19.

### Bonus: the REQ's statement that the heavy lane's exclusion 'is safe there because it refuses a dirty tree' is correct.

**Holds: False**

DOES NOT HOLD. The heavy dirty-tree refusal explicitly EXEMPTS do-work/ (heavy_run.go:221 and :225 both test `!strings.HasPrefix(..., queueStatePrefix)` before recording an offending path), and the untracked seal skips it too (heavy_evidence.go:481). The committed seal does cover do-work/ (heavy_evidence.go:299-308 seals any path not classified by the manifest, and _dev/tests/heavy-lanes.json has non_lane_coverage: null, so do-work/ is unclassified and sealed). Net: heavy is safe against COMMITTED do-work changes but has the same false-green shape as the fast gate for an UNCOMMITTED tracked do-work edit — the one tree the refusal skips is the one tree the committed seal cannot see. Out of REQ-592's write_set; file as a discovered task rather than widening scope.

## Files to change

### `/home/user/skill-do-work/skills/do-work/tools/do-work-cli/internal/heavyverification/fast_stage_evidence.go`

Add `SealExclusions []coverageRule `json:"seal_exclusions"`` to the fastStageManifest struct (currently lines 41-45). Add `func fastStageSealExcludesPath(manifest fastStageManifest, path string) bool { return laneCoversPath(manifest.SealExclusions, path) }` next to fastStageManifestClassifiesPath (line 143). Replace the queueStatePrefix guard at line 195 with `if path == "" || fastStageSealExcludesPath(manifest, path) { continue }` and the identical guard at line 223 the same way — keeping the exclusion test BEFORE the stageCovered test, which is what lets an exclusion beat a coverage rule. In decodeFastStageManifest, validate the new list beside the NonStageCoverage loop (line 126-130): `for _, rule := range manifest.SealExclusions { if err := rule.validate(); err != nil { return fastStageManifest{}, fmt.Errorf("seal exclusion: %w", err) } }`. Note decoder.DisallowUnknownFields() (line 76) means the JSON key and the struct field cannot drift. Rewrite three now-false comments: the struct doc at 38-40 ("the one tree no stage covers"), the tracked-loop comment at 191-194 ("The queue tree is bookkeeping no stage reads"), and the untracked-loop comment at 209-212. queueStatePrefix stays referenced by heavy_run.go and heavy_evidence.go, so it does not become dead.

### `/home/user/skill-do-work/_dev/tests/fast-stages.json`

Three edits. (1) queue-kanban-fast-tests coverage gains {"kind": "subtree", "path": "do-work"} — that stage's board build reads the whole tree, so declaring it as coverage is what seals it AND, via fastStageManifestClassifiesPath, keeps it out of the do-work-cli stage's seal. (2) do-work-cli-fast-tests coverage gains {"kind": "exact", "path": "do-work/archive/UR-003/input.md"} — its single live do-work read. (3) non_stage_coverage becomes [] and gains no do-work entry: verified below, no do-work subtree is unread by both stages. Add "seal_exclusions": [{"kind":"subtree","path":"do-work/runs"}, {"kind":"subtree","path":"do-work/.req-reservations"}, {"kind":"subtree","path":"do-work/deliverables"}, {"kind":"exact","path":"do-work/test-durations.tsv"}] with a comment-free JSON but a matching prose note in the Go struct doc naming the CONDITION each entry satisfies (written by the gate or the orchestrator while a gate runs, and byte-unread by every stage). Rule kinds exact/subtree/suffix-under are the accepted set (heavy_verification.go:288-301).

### `/home/user/skill-do-work/skills/do-work/tools/do-work-cli/internal/heavyverification/fast_stage_evidence_test.go`

Fixture: in fastStageTestManifest (lines 23-50), give beta-stage coverage [module-beta, do-work], set non_stage_coverage to [], add seal_exclusions for do-work/runs and do-work/test-durations.tsv, and rewrite the stale header comment at lines 17-22. MOVE the toolchain probe inputs out of the newly-sealed tree: probes at lines 31 and 40 become ["cat", "do-work/runs/alpha-toolchain.txt"] / beta, with buildFastStageTemplateRepository lines 82-83 and the two probe cases at lines 269 and 276 following — otherwise the "toolchain probe output changed" case mismatches on the file seal instead of the probe output and stops proving what it names. Case rewrites: replace the `queue state changed` case (lines 182-190) with TWO cases — `queue state changed forces the stage that reads it` (stageID beta-stage, same mutation, expect LaneDispositionExecuted + laneReasonFingerprintMismatch, comment naming the failure: one newline appended to do-work/archive/UR-003/input.md made the do-work-cli stage's own test fail while the gate printed `Maintainer verification passed.` and exited 0 with that stage REUSED) and `queue state changed leaves a stage that does not read it reusable` (stageID alpha-stage, same mutation, expect reused, comment naming the failure: sealing the queue tree into every stage would make each REQ claim re-run both stages and kill reuse during a drain). Add `the gate's own duration log still reuses` (stageID beta-stage, append to do-work/test-durations.tsv, expect reused, comment: the stage appends that log itself, so a seal over it never matches its own prior evidence) and `a run-log write still reuses` (stageID beta-stage, write do-work/runs/work-x/notes.md, expect reused, comment: the orchestrator writes run logs while the gate runs).

### `/home/user/skill-do-work/_dev/tests/fast-stage-reuse-behavior.sh`

Fixture manifest (lines 38-53): alpha-stage coverage becomes [module-alpha, do-work], non_stage_coverage becomes [], add the same seal_exclusions pair (do-work/runs subtree, do-work/test-durations.tsv exact). Move the probe input do-work/toolchain.txt to do-work/runs/toolchain.txt (lines 47 and 55) and add do-work/runs to the mkdir -p at lines 30-34, for the same isolation reason as the Go fixture. Rewrite the stale comment at lines 36-37. Case rewrite: line 145 `queue state alone still reuses` becomes `a queue-tree change executes the stage that reads it` with expectation `yes 0 'EXECUTING (fingerprint_mismatch)'`, its comment naming the failure it now catches (a do-work/-only edit reused stale evidence and the whole gate printed Maintainer verification passed. and exited 0 while the stage's own test failed on that same tree). Add a new case straight after it: `the gate's own duration log alone still reuses` — append a row to $project_root/do-work/test-durations.tsv, expect `no 0 'REUSED (fingerprint_match, recorded '`. That case genuinely bites: the fixture has no .gitignore entry for the log, so without the seal_exclusions entry the file is untracked-non-ignored AND stage-covered and would be sealed.

### `/home/user/skill-do-work/skills/do-work-board/tools/queue-kanban/board_live_test.go`

NOT part of the fix — listed because it is the load-bearing evidence for the queue-kanban half and a reviewer will want the exact anchor. liveBoard/liveRepoRoot/liveBoardAt at lines 16-45; suiteCheckoutSkipReason at lines 47-64. No edit needed.

## Red-Green recipe

```
GATE-LEVEL (the REQ's own proof, ~minutes). From a scratch worktree at the merge revision:

  git -C /home/user/skill-do-work worktree add /tmp/req592-red HEAD
  cd /tmp/req592-red
  bash _dev/tests/maintainer-verify.sh                     # run 1: records fast-stage evidence
  printf '\n' >> do-work/archive/UR-003/input.md           # 5608 -> 5609 bytes
  bash _dev/tests/maintainer-verify.sh                     # run 2

RED now: run 2 prints `maintainer-verify: stage do-work-cli-fast-tests: REUSED (fingerprint_match, recorded ...)` and `Maintainer verification passed.` and exits 0, while on the same tree
  (cd skills/do-work/tools/do-work-cli && go test -short -count=1 -run TestDiscoverRepositoryAcceptsProductionLegacyArchiveInputClass ./internal/repositorymodel/)
fails with `production legacy fixture changed size: got 5609 bytes` (assertion at internal/repositorymodel/repository_model_test.go:406-408).

GREEN after the fix: run 2 prints `EXECUTING (fingerprint_mismatch)` for do-work-cli-fast-tests and the gate exits non-zero. Then reset the file, re-warm, and touch ONLY the log:
  git checkout -- do-work/archive/UR-003/input.md && bash _dev/tests/maintainer-verify.sh
  printf 'probe\tprobe\t0.00\t0\n' >> do-work/test-durations.tsv
  bash _dev/tests/maintainer-verify.sh                     # must still print REUSED for both stages, exit 0

Equivalent proof for the queue-kanban half: instead of the input.md edit, append a blank line to any do-work/archive/**/REQ-*.md; run 2 must print EXECUTING for queue-kanban-fast-tests.

UNIT/PROBE LEVEL (seconds, the loop to develop against):
  cd /home/user/skill-do-work/skills/do-work/tools/do-work-cli && go test -count=1 -run TestFastStageReuseDecisionTable ./internal/heavyverification/
  bash /home/user/skill-do-work/_dev/tests/fast-stage-reuse-behavior.sh
Both are RED against the current fast_stage_evidence.go the moment the two assertions are rewritten, and GREEN once the seal_exclusions field and fastStageSealExcludesPath land. Note the evidence store is keyed by working-tree root plus Git common dir (fast_stage_evidence.go:378-380), so a linked worktree has its own key space and cannot inherit the main checkout's warm records.
```

## Risks

- R1. The queue-kanban fast stage stops reusing in practice. With do-work/ as its coverage it seals ~730 tracked files / ~23.4 MB (archive 702 files / 20.5 MB, user-requests 7 / 2.6 MB, queue 7, working 5, audits 2, plus 8 top-level files), and every REQ claim, move, or archive changes one of them. During a drain loop that stage will re-execute on essentially every gate run — the REQ-591 reuse feature becomes dead weight for it. This is CORRECT by REQ-592's own requirement (the stage really does read those bytes), but the REQ does not say it, and it is the whole cost of the change. The do-work-cli stage keeps reuse because its only do-work input is one exact file. Hashing cost itself is negligible (~50 ms of SHA-256); the churn is the cost.
- R2. The seal is byte-level; part of the real dependency is existence-level, and the fix cannot express that. filementions.go:35-56 collectRepoFileMentions stats EVERY repo-relative path mentioned in any REQ/UR body against the repo root, and generate.go:713 runs it on every live-board build. Empirically those mentions reach every do-work subtree: 125 mentions under do-work/runs/, 97 of do-work/CHECKPOINT.md, 51 of do-work/calibration-log.tsv, 37 under do-work/audits/, 8 of do-work/test-durations.tsv, 5 under do-work/deliverables/. So excluding do-work/runs/ and do-work/deliverables/ from the seal is a judgement, not a proof: creating or deleting a mentioned file there flips a bool in the shipped board JSON. No fast-stage assertion currently reads that map (filementions_test.go:41 uses a synthetic tree), so the exposure is a JSON-byte change no test checks. The exclusion comment must say this rather than claim nothing reads those trees.
- R3. The two fixture toolchain probes read files under do-work/ (fast_stage_evidence_test.go:31,40,82-83 and fast-stage-reuse-behavior.sh:47,55). If they are left where they are once do-work/ is sealed, the 'toolchain probe output changed' case still passes but for the wrong reason (the file's own byte seal moved, not the probe output), and 'toolchain probe cannot run' reaches fingerprint_uncertain through a missing tracked seal input instead of a failing probe. Both cases silently stop testing what they name. Moving the probe inputs into the excluded subtree is not cosmetic.
- R4. The heavy lane has the same false-green shape for an UNCOMMITTED do-work edit. heavy_run.go:221,225 exempt do-work/ from the dirty-tree refusal and heavy_evidence.go:481 exempts it from the untracked seal, while the committed seal (heavy_evidence.go:299-308, with _dev/tests/heavy-lanes.json carrying non_lane_coverage: null) only sees HEAD objects. REQ-592 asserts the heavy exclusion is safe because heavy refuses a dirty tree — but do-work/ is precisely the tree that refusal skips. Out of the declared write_set; raise it as a discovered task instead of widening this REQ.
- R5. do-work/.req-reservations/ holds 162 tracked marker files that the allocator creates and removes during normal work. They are genuinely unread by both stages (walk.go:198 prunes hidden directories, and repoFileMentionPattern at filementions.go:25-26 cannot match a dot-leading path segment, so they are never stat'd either). But they are TRACKED, so putting them in non_stage_coverage would not save them once do-work/ is a stage's coverage — the tracked loop seals on `stageCovered || !classified` and stageCovered wins. They must go in seal_exclusions, not non_stage_coverage. Getting that wrong makes every REQ-number reservation invalidate the queue-kanban stage.
- R6. Adding seal_exclusions is a new closed enumeration, exactly the shape _dev/primes/prime-shell-commands.md 'Closed Enumerations Go Stale' warns about. Keep it in the manifest (data a reviewer sees beside coverage) rather than a Go slice, and state the admission CONDITION in the struct doc — a do-work path written by the gate or the orchestrator during a gate run and byte-unread by every stage — so the next person adding a churn directory has a test to apply instead of a list to copy.

## Which do-work subtrees NO fast stage reads — verified, not assumed

Answer: essentially none. Only `do-work/.req-reservations/` (and any other hidden directory under `do-work/`) is provably unread by both stages. `non_stage_coverage` should therefore become `[]`, and the churn trees must be handled by a seal exclusion instead.

How each was determined:

BYTE-READ by a fast stage:
- `do-work/archive/UR-003/input.md` — do-work-cli stage, exact 5608-byte assertion (repository_model_test.go:397-417). It is the ONLY real-tree do-work read in the entire do-work-cli module: a grep for six-level parent traversals finds just repository_model_test.go:398 plus heavy_maintainer_tree_test.go:26,197 (which read `_dev/`, not `do-work/`).
- Every `REQ-*.md` under `do-work/queue/`, `do-work/working/`, `do-work/archive/**` and every `input.md` under `do-work/user-requests/**`, `do-work/archive/**` — parsed by buildBoard via enumerateDoWorkTree (walk.go:113-179, parse at model.go:708). Read by the queue-kanban stage's live tests. The pinned figures in durations_test.go:214-253 are computed directly from archived REQ frontmatter, so an archive edit moves them.
- `do-work/notes.md` and `do-work/testers.md` when present (walk.go:165-174, model.go:546). Neither exists today.

NAME- OR EXISTENCE-READ (so still not excludable safely):
- Directory listings of every non-pruned dir under do-work/: a new `REQ-*.md` anywhere outside queue/working/archive becomes a stray-file board warning (walk.go:150-158, model.go:456-463). `do-work/audits/` is walked for exactly this reason even though its two audit-*.md files are never parsed.
- Any path mentioned in a REQ/UR body, stat'd by collectRepoFileMentions (filementions.go:45). Measured mention counts by subtree are in R2.

PRUNED FROM THE WALK (byte-unread):
- `do-work/runs/` and `do-work/deliverables/` — walk.go:192.
- Any `assets/` directory at any depth — walk.go:195; 23 such directories, 25 tracked files.
- Any hidden directory — walk.go:198; this is what makes `do-work/.req-reservations/` (162 files) unread.

Sizes, for the reuse-cost judgement: 1121 tracked files / 26.3 MB under do-work/, of which runs/ is 229 files / 2.9 MB and .req-reservations/ is 162 files / 768 bytes. Excluding those two leaves ~730 files / ~23.4 MB sealed into the queue-kanban stage.

## The narrowest fix — why coverage, not an emptied classification

REQ-592 proposes "give the fast-stage seal its own exclusion set instead of inheriting queueStatePrefix". Deleting the `non_stage_coverage: [do-work]` entry alone would make do-work paths UNCLASSIFIED, and fast_stage_evidence.go:202 seals every unclassified path into EVERY stage. That would force the do-work-cli stage to re-execute on every queue edit even though its only do-work input is one file. Declaring `do-work` as the queue-kanban stage's COVERAGE instead gets per-stage separation from machinery that already exists (fastStageManifestClassifiesPath, fast_stage_evidence.go:143-153, returns true for a path any stage covers, so the do-work-cli stage skips it), and costs one manifest line rather than a second classification concept.

The one thing the existing machinery cannot express is "sealed nowhere, even where a stage covers it" — needed because once `do-work` is queue-kanban's coverage, `stageCovered` is true and the untracked-ignored skip at fast_stage_evidence.go:227 stops protecting `do-work/test-durations.tsv`. That is the single new concept: `seal_exclusions`, checked before `stageCovered`, exactly where the queueStatePrefix guard sits now.

Names introduced (both multi-word and plain-text findable, per CLAUDE.md § Naming Conventions): manifest field `seal_exclusions` / Go field `SealExclusions`, and predicate `fastStageSealExcludesPath`.

## Two things the REQ gets wrong or leaves out

1. Its Requirements bullet says test-durations.tsv "needs an explicit narrow exclusion rather than the whole-tree one". True under the fix above, but NOT true under a bare non_stage_coverage deletion — the file is gitignored and untracked, so fast_stage_evidence.go:227 already skips it. Whoever builds this should know the exclusion is load-bearing only because of the coverage declaration, or they will write a test that passes for the wrong reason.
2. It says the heavy lane's exclusion "is safe there because it refuses a dirty tree". The refusal skips do-work/ by construction (heavy_run.go:225). See R4.

## Blast radius outside the write_set

None. `_dev/tests/fast-stages.json` is referenced only by heavy_commands.go:102 (the default manifest path constant), and no contract test pins its contents — a grep across `_dev/` and `skills/` for fast-stages.json / non_stage_coverage / NonStageCoverage returns only the four write_set files plus that constant. `queueStatePrefix` remains used by heavy_run.go:221,225 and heavy_evidence.go:481, so removing its two fast-path uses leaves no dead identifier.

## Files, for reference (all absolute)

- /home/user/skill-do-work/do-work/working/REQ-592-seal-the-do-work-tree-into-both-fast-gate-stages.md
- /home/user/skill-do-work/_dev/tests/fast-stages.json
- /home/user/skill-do-work/_dev/tests/fast-stage-reuse-behavior.sh
- /home/user/skill-do-work/_dev/tests/maintainer-verify.sh (stage wiring: 102-199, 716-728, 773-779)
- /home/user/skill-do-work/_dev/tests/run-go-tests-with-budget.sh (duration-log append: 99-111)
- /home/user/skill-do-work/_dev/tests/test-duration-log.sh (log path default: line 4)
- /home/user/skill-do-work/skills/do-work/tools/do-work-cli/internal/heavyverification/fast_stage_evidence.go
- /home/user/skill-do-work/skills/do-work/tools/do-work-cli/internal/heavyverification/fast_stage_evidence_test.go
- /home/user/skill-do-work/skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_run.go (queueStatePrefix: 28-31; dirty-tree refusal: 205-234)
- /home/user/skill-do-work/skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_evidence.go (committed seal: 288-320; untracked seal: 465-495)
- /home/user/skill-do-work/skills/do-work/tools/do-work-cli/internal/heavyverification/heavy_verification.go (coverageRule kinds: 35-40, 288-323)
- /home/user/skill-do-work/skills/do-work/tools/do-work-cli/internal/repositorymodel/repository_model_test.go (397-417)
- /home/user/skill-do-work/skills/do-work-board/tools/queue-kanban/walk.go (103-199)
- /home/user/skill-do-work/skills/do-work-board/tools/queue-kanban/filementions.go (25-56)
- /home/user/skill-do-work/skills/do-work-board/tools/queue-kanban/board_live_test.go, durations_test.go, citations_test.go, generate_test.go
