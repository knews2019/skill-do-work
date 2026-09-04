# Queue run

Status: in-progress

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
