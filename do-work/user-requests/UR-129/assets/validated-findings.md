# Accepted Validation Evidence

Source-verified at d4ea51d7; no reproduction probes executed. Eleven comments map to seven accepted defects. None already fixed.

## REQ-615 (Distinguish functional manifest edits from release metadata) — F1, F5 — P2

Verdict: Accept.

Allow a release whose shipped implementation changes functional dependencies, entry points, or build configuration in package.json, Cargo.toml, or pyproject.toml; exclude only edits limited to release metadata.

- skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard.go:46-48 excludes each entire metadata-bearing path before root matching.
- skills/do-work/tools/do-work-cli/internal/releaseownership/release_ownership.go:42-46 defines IsReleaseMetadataPath as whether the file carries release metadata at all.
- skills/do-work/tools/do-work-cli/internal/finalization/finalization_prepare.go:129 applies this guard before release planning.

Surface-cost: N/A — direct correction of the existing overbroad guard. No second validation layer is needed.

History/scope: Introduced by 6ecff1bc7aad6ef8bc80f0ecb515270b3d20c45f. Ordinary consumer repositories bypass the guard; no such manifests currently exist under this checkout's skills tree. This is a supported-input defect verified from source, not an observed failed release here.

Regression target: A focused regression in finalization_release_guard_test.go accepts substantive edits and rejects version-only edits for both provenance modes, including package.json, Cargo.toml, and pyproject.toml.

## REQ-616 (Require pending shipped changes for primary release) — F2, F4, F10 — P2

Verdict: Accept.

Require actual pending shipped implementation changes for primary_commit release provenance instead of treating commit_paths membership as evidence of a change.

- skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard.go:64-65 returns manifest.CommitPaths directly; lines 49-51 accept any listed shipped path.
- skills/do-work/tools/do-work-cli/internal/finalization/finalization_prepare.go:158-174 checks normalization and required-target coverage, not pending implementation changes.
- skills/do-work/tools/do-work-cli/internal/gittransaction/exact_commit.go:45-58 filters clean allowlisted entries before staging.

Surface-cost: N/A — direct correction of the existing guard's evidence source.

History/scope: Introduced by 6ecff1bc7aad6ef8bc80f0ecb515270b3d20c45f. The existing primary-provenance guard test even accepts a shipped path without creating it; no later preparation check compensates.

Regression target: A focused regression in finalization_release_guard_test.go refuses the clean-shipped-path bypass and accepts actual shipped additions, edits, and deletions.

## REQ-617 (Preserve exact Git paths in release guard) — F3, F6 — P2

Verdict: Accept.

Read supplied implementation commit filenames as exact NUL-delimited paths so Git display quoting does not reject valid shipped changes.

- skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard.go:68 omits -z; lines 79-81 split display-formatted output on newline after TrimSpace.
- skills/do-work/tools/do-work-cli/internal/finalization/finalization_release_guard.go:45-50 does not unquote the path before comparing the shipped-root prefix.

Surface-cost: N/A — direct parser repair.

History/scope: Introduced by 6ecff1bc7aad6ef8bc80f0ecb515270b3d20c45f. Current HEAD has no special filenames; validation established the code path, not an executed end-to-end reproduction.

Regression target: Focused finalization_release_guard_test.go cases accept shipped non-ASCII/control-character/quote filenames while preserving existing ordinary, merge, root-commit, and maintainer-only behavior.

## REQ-618 (Preserve committed-risk before finalization rollback) — F7 — P1

Verdict: Accept.

Preserve a commit that already exists and its unresolved verification risk when exact-path commit reports failure; do not restore pre-commit lifecycle state after HEAD advances.

- skills/do-work/tools/do-work-cli/internal/gittransaction/exact_commit.go:81-97 can return a nonempty CommitSHA alongside failure when a hook stages an extra path.
- skills/do-work/tools/do-work-cli/internal/gittransaction/git_transaction.go:1458-1461 preserves CommitSHA and revert argv in committedRisk.
- skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go:134-138 returns on Failure before recording SHA; lines 28-33 then enter pre-primary rollback.
- skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go:214-245 restores lifecycle/release preimages and resets prepared identity; finalizationFailure at line 645 loses the nested risk outcome/revert evidence.
- skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go:530-559 cannot compensate because matchingHeadCommit rejects the outside-allowlist extra path.

Surface-cost: Earned — the concrete replay is a hook staging an extra tracked path after preparation. A small extension of existing commit/rollback checks is cheaper than corrupting the lifecycle trail; the GREEN hook/recovery test keeps this condition covered.

History/scope: Failure-before-SHA ordering dates to 761d8e6ab47dd940db6d686ef641ea919ee9bb31; the rollback defer arrived in e191b266b9d501ece52d1c9888f7711aaf72f281. The original wording "now" is not supported as a recent regression.

Regression target: A focused finalization hook/recovery regression retains the committed archive and lifecycle, records SHA and risk in journal/result, and proves recovery neither duplicates the commit nor silently clears verification failure.

## REQ-619 (Verify primary commit content against prepared identity) — F8 — P1

Verdict: Accept.

Verify committed implementation content against the complete prepared identity before accepting the primary commit, including hook-staged rewrites and deletions.

- skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go:133 passes nil as the post-commit verifier.
- skills/do-work/tools/do-work-cli/internal/gittransaction/exact_commit.go:88-99 compares filenames only and skips a nil verifier.
- skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go:583-625 verifies archive identity/status/provenance and release images, not implementation bytes.
- skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go:530-559 compares the prepared digest only during recovery detection before making a new commit.

Surface-cost: Earned — a staged hook rewrite/deletion preserves the path set while violating prepared content. Reuse the existing callback/digest mechanisms rather than adding a new state system; paired negative tests and a normal commit control prove the value.

History/scope: The nil callback is present from foundation commit 761d8e6ab47dd940db6d686ef641ea919ee9bb31. No recent removal is established. Deleting a newly created file can already cause a path mismatch; the verified gap specifically includes tracked files.

Regression target: Focused finalization tests fail on staged hook rewrite and deletion with retained SHA/journal risk, and pass for the prepared tracked-plus-new-file content.

## REQ-620 (Include untracked files in prepared commit digest) — F9 — P1

Verdict: Accept.

Bind prepared and recovered commit identities to the same complete content, including untracked archive destinations and implementation additions.

- skills/do-work/tools/do-work-cli/internal/finalization/finalization_discovery.go:1581-1591 hashes git diff HEAD, omitting untracked files.
- skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go:114-120 persists the digest before committing.
- skills/do-work/tools/do-work-cli/internal/gittransaction/exact_commit.go:37-42 requires an empty index; line 58 stages additions only later.
- skills/do-work/tools/do-work-cli/internal/finalization/finalization_apply.go:540-543 hashes the committed diff including additions when trying to recognize recovery.

Surface-cost: N/A — direct repair of the existing commit identity computation.

History/scope: This implementation dates to e191b266b9d501ece52d1c9888f7711aaf72f281. Available local history does not support the original claim to "restore" an earlier temporary-index implementation.

Regression target: A focused finalization_recovery_test.go regression recognizes the tracked-plus-new-file commit after the interruption, preserves the real index, and recovers without duplicate commit or pre-primary rollback.

## REQ-621 (Emit diagnostics for exact-text command results) — F11 — P2

Verdict: Accept.

Expose existing warning/error findings on stderr for exact-text results while preserving compatibility stdout and intentional success/fallback semantics.

- skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_image.go:232-243 adds per-image warning findings but sets ExactTextOutput to only the directory, or an empty string when all images fail.
- skills/do-work/tools/do-work-cli/internal/resultmodel/result_model.go:1027-1028 returns ExactTextOutput before rendering findings.
- skills/do-work/tools/do-work-cli/internal/commandruntime/command_runtime.go:94-95 writes only that rendered output; main.go wires the runtime to stdout.
- skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_image_test.go:42-56 already asserts all-failed success, empty exact text, and two warning findings.

Surface-cost: Earned — partial/all-failed image generation produces real typed warnings that text callers cannot see. A small stderr projection of existing findings restores useful fallback guidance; tests pin both diagnostic presence and compatibility output.

History/scope: The old shell implementation at a7c975c5^:skills/do-work-toolbox/scripts/generate-report-image-batch.sh prints MISSING fallback warnings to stderr. No recent removal of a Go stderr writer is established; the gap arose during shell-to-Go migration (a7c975c5). The shipped report action uses JSON and retains findings, limiting impact to text compatibility callers.

Regression target: Focused commandruntime and image-batch regressions assert per-image stderr diagnostics alongside unchanged stdout and exit status for partial/all-failed results, with unchanged JSON findings.
