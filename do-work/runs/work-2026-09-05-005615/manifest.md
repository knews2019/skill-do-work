# Queue run

Status: in-progress; waiting for active foreign writers

User scope: do-work run to completion, serial claims; do not finish while agents or background commands remain pending.

- REQ-567: completed; release 0.281.1; archive do-work/archive/REQ-567-repair-shipped-lesson-links-to-archived-ur-paths.md; commit 6aa5776b.

- REQ-503: completed; release 0.282.0; commit cefe971d; independent review Partial (78.75%); two noncritical report-only findings; output REQ-503-review.md.
- User priority: run REQ-509 immediately after REQ-503 settles, then resume full queue.

- REQ-509: claimed; Route B; exploration consumed; preparing baseline evidence and isolated builder.

- Read-only later-step prep: prep_504 → REQ-504-prep.md consumed; independent exact-revision review → REQ-504-review.md pending.
- REQ-509 isolated branch: codex/worktree-agent-REQ-509-shared-principles; checkout .git/work-run-20260905/worktree-agent-REQ-509-shared-principles; expected builder artifact REQ-509-handback.md.

- REQ-509 builder build_509 accepted; source implementation in progress; output REQ-509-handback.md pending.

- Read-only later-step prep: prep_510 → REQ-510-prep.md pending.

- REQ-509 builder handback consumed; c68102502066a65c5be0573b17236b0a1ff58695; RED/GREEN and focused contracts pass; integrating next.
- REQ-504 review complete (Acceptance Fail, critical legacy-checkpoint finding); remediation waits until REQ-509 completes.
- REQ-510 prep complete.

- REQ-509 integrated range: 622a5e55de332984d7e180615a2c5c2b6c7ef2d7..2ba5b432658853690e8e5a6d20bd2dcc147e9ada; qualify and focused evidence green; repository gate running; independent review_509 → REQ-509-review.md pending.
- REQ-504 read-only remediation plan → REQ-504-remediation-plan.md pending.

- REQ-509 review accepted (100%, Pass); selected lane evidence complete after canonical version-mirror repair and installer retry; preparing finalization.
- REQ-504 remediation plan complete and ready after REQ-509.

- REQ-509 completed and delivered; release0.282.1; commit9f19533f0112625e9d3fe2bf4f1ac0eaae8f4b47; manifest validation passed; own checkout/branch removed without force; untracked handback consumed after promotion to archive.

- REQ-504 claimed at 3501aedd; critical review consumed, seven-file remediation planned; preflight and direct maintainer gate passed. Builder build_504 → REQ-504-handback.md pending; isolated checkout .git/work-run-20260905/worktree-agent-REQ-504-legacy-checkpoints, branch codex/worktree-agent-REQ-504-legacy-checkpoints.
- REQ-505 read-only review complete, Partial 71.25%, two noncritical report-only findings; consume after canonical claim.

- REQ-504 builder accepted; remediation in progress.
- Read-only later-step review: review_506 → REQ-506-review.md pending.

- Read-only later-step review: review_507 → REQ-507-review.md pending. Saved staged-lane red was repaired by REQ-547 commit c5dff3db; fresh current verification required when claimed.

- Main read-only REQ-544-prep.md ready; publication defects and live writer/test grammar discrepancy identified; cleanup now delegates to shared removal.
- Concurrent unrelated dirt observed and left alone: ai-reports/2026-09-03_2145_do-work-velocity-and-pending-queue-speed/index.html; skills/do-work/tools/do-work-cli/internal/publication/release_mirrors_test.go.

- REQ-507 independent review complete (Partial82.5%, two noncritical report-only findings); saved public finalization/phase/refusal and result parity checks passed. Read-only reviewer checkout removed.

- REQ-506 independent review complete: Fail50%, critical false-success evidence with real launch failure/timeout. Read-only remediation planning → REQ-506-remediation-plan.md pending.
- Other session completed REQ-569 at b3ff309e (release0.283.0); its publication edits are committed and no longer dirty. Keep that implementation separate from REQ-504 attribution.

- REQ-504 completed builder handback consumed from disk:4e351d172b14b822dd5027d3c13d12874ef5774c; no pending builder processes. Integrating seven-file remediation; original cumulative base773787b74acddfdfc4c16498a89d99a5cc3ab716 remains fixed.

- REQ-504 integrated at6a11b60c83615791769d57b082580f0b69323984; cumulative773787b7..6a11b60c; qualification/scope satisfied, integrated focused/full gates running; independent review_504_remediation → REQ-504-re-review.md pending.
- REQ-506 remediation plan complete (eight paths, three tasks), ready after canonical505/506 progression.

- REQ-504 integrated canonical gate and focused evidence passed; selected six heavy lanes executing at6a11b60c with explicitChromium and no evidence reuse.
- Read-only later-step plan: review_507 → REQ-544-plan.md pending.

- REQ-504 independent re-review complete:100%, Pass, no findings; RED replay on pre-repaircode independently reproducedbothF01failures. Fresh heavy verification still pending.

- REQ-504 heavy initial: JS0/6s,browser0/90s,CLI0/65s,staged0/25s,updater1/23s,installer2/5s. Retry updater1/27s installer2/7s. Full retry log identifies disk full; no heavy process remains. Await user permission to inspect/clear configured external Go build cache. Claim remains intact, all source/review evidence preserved; no false completion.
- Read-only later-step plan: review_506 → REQ-562-plan.md pending.

- User chose to free disk space personally; no inspection or deletion outside the repository was authorized or performed. No tests run until sufficient space returns. REQ-562 read-only plan complete; all subagents completed, no background process pending.

- REQ-504 updater/installer passed after user restored47GB; all6 selected lanes greenunskipped, finalization prepared next.

- REQ-504 completed/delivered: release0.284.0, finalizationac2196e172ceb8115c659394781ea95bde205d49, implementation6a11b60c83615791769d57b082580f0b69323984. Exact archive source and lesson promoted; schema/release mirror validation passed; branch/worktree removed withoutforce. Consumed handback deleted. All selected lanes passed afterspace restoration, none skipped.

- REQ-505 claimed375735da, exactsavedrangeandheavyresume re-resolved; independentPartial71.25% accepted forcompletion withtwo noncriticalreportonlyfindings; finalizationpreparing.

- REQ-505 completed/delivered:0.285.0, finalization6eaac6674a8026f7d953a01047508fff608b6d2e, originalimplementation716187b847d1de0402b69587a2fe5cf7e7bd8516 preserved, two noncriticalreport-onlyfindings.
- Canonicalrecover blocked on activeother-session REQ-570 +runwork-2026-09-05-020017 manifest; bytespreserved, awaitingownerbookkeepingcommit before506claim.
- Read-only later-step prep: prep_552 → REQ-552-plan.md pending; mainpreparingREQ565 ownernotes.

- Read-only REQ-552 plan and REQ-565 owner notes complete; no agent, test process, or builder remains pending. REQ-552 requires claim-time judgment on two obsolete executable-fault fixtures. Foreign REQ-570 bookkeeping remains active and preserved.

- All later-step preparation complete: REQ-486-plan.md, REQ-571-prep.md, small-prose-prep.md (549/554/555/556), shell-go-prep.md (553/557/558). No agent, background command, or builder checkout remains pending. No new tests or source changes during shared-state wait. Preparing final canonical recovery verification; other-session REQ-570, working/baseline.json and run work-2026-09-05-020017 remain untouched.

## Resume boundary

Final verification `do-work-cli --format json recover-finalization --discover` refused with `FINALIZATION-DISCOVERY-AMBIGUOUS`; changes were empty. Evidence: `.git/work-run-20260905/recovery-boundary-final.json`. Foreign paths are the active REQ-570 working record, `do-work/working/baseline.json`, and the brief/exploration/plan/manifest in `do-work/runs/work-2026-09-05-020017/`. Every byte is preserved. No claim, takeover, source fix, or foreign commit was attempted.

This run completed REQ-567, REQ-503, prioritized REQ-509, REQ-504 (one successful remediation), and REQ-505. REQ-506 remains pending with the independently reproduced critical focused-execution failure and its eight-file remediation plan ready. REQ-507 remains pending: saved staged lane failed, later REQ-547 repaired the obsolete assertion, fresh current verification still required. Neither is falsely marked complete or still in the removed-by-answer old heavy status. Other pending work remains unclaimed.

Resume normal `do-work run` after the active REQ-570 owner settles its bookkeeping. Start with canonical recover and queue advance, re-read any newly landed work action and lifecycle changes, and consume the saved plans only after claim-time revalidation. Do not use sole-authority recovery while that other owner remains active. All this run's agents and background commands finished; its builder worktrees were already cleaned, and all preparation is durable. No external cache files were inspected or deleted; the user freed the space personally.

- User resumed full run with explicit instruction to remain active through pending work. Canonical recover still refuses REQ-570 lifecycle/run/baseline dirt; a read-only owned watcher is pending (session98623). REQ-570 builder branch has now committed d33d4c52; its owner remains active. Another session claimed REQ-572; do not take it over. Activity prep complete for572/573 and preserved in activity-prep.md. Main characterized a real RLIMIT_FSIZE copy-failure seam for552; fixture cleanup completed, evidence saved, no source/queue mutation.

- User's explicit repeat of run-to-completion after the refusal report overrides the earlier stop. Identified foreign active writes remain preserved; no sole-authority takeover or false recovery success. Canonical next selected506 and advance claimed it at7ad53bff. One remediation now in preflight; original evidence retained under Prior headings. Other owners570/572 and their source/release work remain excluded from staging.

- REQ-506 canonically deferred behind REQ-577 by12c1fa537b80255cc36bbcdc47e2d71185eada2b. REQ-577 claimed dc12d12ce434b39c6471f89601e805028dd7e7f9, RouteA, recorded matching red twice. Builder handback REQ-577-handback.md pending.
