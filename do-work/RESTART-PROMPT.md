```
do-work run --fan-out 1
This command is sufficient; everything below it is context.
```

---

## Reference

- **State gate and parallelism:** REQ-420 now depends on REQ-478 in addition to completed REQ-419. REQ-478 reached its archived terminal state while this prompt was being written, but its foreign lifecycle/source diff is still uncommitted in the shared checkout. Fan-out 1 remains mandatory until that owner commits and the overlapping `skills/do-work/actions/capture.md` / `skills/do-work/actions/work-reference.md` bytes are reconciled. Critical path: settle REQ-478's foreign commit, resume REQ-420, then dependent REQ-460 may run. No `pending-answers` REQ exists.
- **REQ-420 integration state:** ACTIVE in the main checkout, `status: claimed`, writer `t2s-Virtual-Machine.local:/Users/t2/Desktop/e1-experimental-repos/skill-do-work2`. It is not merged or committed; merge commit and full merge range: none. Main was `59d038c4` before this handoff commit.
- **REQ-420 completed evidence:** RED was captured with thinness exit 1 on legacy bodies and parity exit 1 after thinness became non-vacuous. All 41 retained shell paths are now launchers and thinness is GREEN. Focused `go test ./internal/corehelpers ./internal/doctor ./internal/knowledgecommands ./internal/toolboxcommands` is GREEN. Missing Go compatibility seams, the 11-file audit retirement, maintainer-lane removal, and partial SessionStart/board authority reconciliation are in the uncommitted diff.
- **REQ-420 current failures:** The last complete parity run failed 10/18 owners: audit timestamps (12), capture screenshot (3), image batch (16), image (14), last30days (12), protected inventory (2), portfolio (16), qualify (21), repair timestamps (33), and run-blocked (2). Later adapter changes make that recount stale; thinness itself is green. No canonical maintainer gate, vet, no-Python/no-jq install/update, flat-Just/no-LLM, or staged-skills final result exists yet.
- **REQ-420 remaining order:** (1) preserve and settle REQ-478's terminal-but-uncommitted foreign transition, then reconcile the two mixed authority files; (2) adapt all 18 prescribed fixtures to observable canonical CLI behavior, rerun each failing owner, then full parity; (3) close remaining exact semantics and finish folded differential, mutation, actionability, and BKB committed tests; (4) finish authority/restatement and gate wiring; (5) run gofmt/shellcheck, focused and board tests, vet, thinness/parity, no-Python/no-jq install/update, flat Just/no-LLM/aliases, maintainer self-test, then the unpiped canonical maintainer gate; (6) qualify, independently review, apply at most one remediation, archive, commit, and release.
- **P-A-U:** PLAN is checked. APPLY and UNIFY remain unchecked. The builder made no commit.
- **Foreign terminal transition:** REQ-478 was owned by `t2s-Virtual-Machine:/Users/t2/Desktop/e1-experimental-repos/skill-do-work2`; its checkpoint entry is now removed and its archive path exists, but the transition has not been committed. Its queue deletion/archive, `capture-reference.md`, `request_model_test.go`, lessons-index/prime edits, and detached-fixture lifecycle are foreign; `capture.md` and `work-reference.md` are mixed REQ-420/REQ-478 paths. Leave every REQ-478 byte intact.
- **Main checkout verdict:** ACTIVE. Exact uncommitted paths at survey refresh:

```text
M README.md
 M _dev/primes/lessons-action-files.md
 M _dev/primes/prime-kanban-board.md
 M _dev/tests/contract-regressions.sh
 M _dev/tests/maintainer-verify.sh
 M _dev/tests/prescribed-shell-cases/generate-report-image-batch.sh
 M _dev/tests/prescribed-shell-cases/generate-report-image.sh
 M _dev/tests/prescribed-shell-cases/install-last30days.sh
 M _dev/tests/prescribed-shell-cases/publish-portfolio-summary.sh
 M _dev/tests/prescribed-shell-cases/run-blocked-check.sh
 M do-work/CHECKPOINT.md
 M do-work/calibration-log.tsv
 M do-work/lessons-index.md
 D do-work/queue/REQ-478-capture-stamps-required-lessons-under-token-budget.md
 M do-work/working/REQ-420-replace-shell-implementations-verify-parity.md
 M skills/do-work-board/tools/queue-kanban/frontmatter_cli.go
 M skills/do-work-board/tools/queue-kanban/prime-do-kanban.md
 M skills/do-work-board/tools/queue-kanban/timestamp.go
 M skills/do-work-board/tools/queue-kanban/timestamp_test.go
 M skills/do-work-board/tools/queue-kanban/verify.go
 M skills/do-work-board/tools/queue-kanban/verify_test.go
 M skills/do-work-knowledge/hooks/memory-session-start.sh
 M skills/do-work-knowledge/hooks/memory-stop-capture.sh
 M skills/do-work-knowledge/scripts/install-memory-hooks.sh
 M skills/do-work-knowledge/scripts/lexical-memory-recall.sh
 M skills/do-work-toolbox/scripts/architecture-report-preflight.sh
 M skills/do-work-toolbox/scripts/generate-report-image-batch.sh
 M skills/do-work-toolbox/scripts/generate-report-image.sh
 M skills/do-work-toolbox/scripts/install-last30days.sh
 M skills/do-work-toolbox/scripts/publish-portfolio-summary.sh
 D skills/do-work-toolbox/tools/audit-metrics/.gitignore
 D skills/do-work-toolbox/tools/audit-metrics/churn.go
 D skills/do-work-toolbox/tools/audit-metrics/churn_test.go
 D skills/do-work-toolbox/tools/audit-metrics/distribution.go
 D skills/do-work-toolbox/tools/audit-metrics/distribution_test.go
 D skills/do-work-toolbox/tools/audit-metrics/git_support.go
 D skills/do-work-toolbox/tools/audit-metrics/go.mod
 D skills/do-work-toolbox/tools/audit-metrics/inventory.go
 D skills/do-work-toolbox/tools/audit-metrics/inventory_test.go
 D skills/do-work-toolbox/tools/audit-metrics/main.go
 D skills/do-work-toolbox/tools/audit-metrics/prime-audit-metrics.md
 M skills/do-work/actions/capture-reference.md
 M skills/do-work/actions/capture.md
 M skills/do-work/actions/work-reference.md
 M skills/do-work/hooks/session-start.sh
 M skills/do-work/scripts/add-local-git-exclude.sh
 M skills/do-work/scripts/atomic-download.sh
 M skills/do-work/scripts/audit-archive-timestamps.sh
 M skills/do-work/scripts/capture-screenshot.sh
 M skills/do-work/scripts/cleanup-req-reservations.sh
 M skills/do-work/scripts/handoff-state-survey.sh
 M skills/do-work/scripts/protected-inventory.sh
 M skills/do-work/scripts/repair-req-timestamps.sh
 M skills/do-work/scripts/run-blocked-check.sh
 M skills/do-work/scripts/show-commit-diff.sh
 M skills/do-work/scripts/stage-exact-deletion.sh
 M skills/do-work/tools/checks/archive-collision.sh
 M skills/do-work/tools/checks/associate-files.sh
 M skills/do-work/tools/checks/blanked-req-scan.sh
 M skills/do-work/tools/checks/preflight.sh
 M skills/do-work/tools/checks/qualify.sh
 M skills/do-work/tools/checks/record-commit-hash.sh
 M skills/do-work/tools/checks/scope-drift.sh
 M skills/do-work/tools/checks/uncommitted-inventory.sh
 M skills/do-work/tools/do-work-cli.sh
 M skills/do-work/tools/do-work-cli/internal/corehelpers/commands.go
 M skills/do-work/tools/do-work-cli/internal/corehelpers/commands_test.go
 M skills/do-work/tools/do-work-cli/internal/corehelpers/inventory.go
 M skills/do-work/tools/do-work-cli/internal/corehelpers/reservations.go
 M skills/do-work/tools/do-work-cli/internal/doctor/doctor_commands.go
 M skills/do-work/tools/do-work-cli/internal/knowledgecommands/bkb_init.go
 M skills/do-work/tools/do-work-cli/internal/knowledgecommands/bkb_scan.go
 M skills/do-work/tools/do-work-cli/internal/knowledgecommands/commands.go
 M skills/do-work/tools/do-work-cli/internal/knowledgecommands/commands_test.go
 M skills/do-work/tools/do-work-cli/internal/knowledgecommands/memory_commands.go
 M skills/do-work/tools/do-work-cli/internal/requestmodel/request_model_test.go
 M skills/do-work/tools/do-work-cli/internal/toolboxcommands/architecture.go
 M skills/do-work/tools/do-work-cli/internal/toolboxcommands/architecture_test.go
 M skills/do-work/tools/do-work-cli/internal/toolboxcommands/audit_metrics_test.go
 M skills/do-work/tools/do-work-cli/internal/toolboxcommands/last30days.go
 M skills/do-work/tools/do-work-cli/internal/toolboxcommands/portfolio.go
 M skills/do-work/tools/do-work-cli/internal/toolboxcommands/report_image.go
 M skills/do-work/tools/do-work-update.sh
 M skills/do-work/tools/estimate-p50.sh
 M skills/do-work/tools/fetch-upstream-archive.sh
 M skills/do-work/tools/install-do-work-suite.sh
 M skills/do-work/tools/replace-text-section.sh
 M skills/do-work/tools/select-simple-reqs.sh
 M skills/do-work/tools/validate-suite-manifest.sh
 M tools/fetch-upstream-archive.sh
 M tools/install-do-work-suite.sh
 M tools/replace-text-section.sh
 M tools/validate-suite-manifest.sh
?? _dev/tests/fixtures/shipped-shell-command-map.tsv
?? _dev/tests/shipped-shell-parity.sh
?? _dev/tests/shipped-shell-thinness.sh
?? do-work/archive/REQ-478-capture-stamps-required-lessons-under-token-budget.md
?? do-work/runs/work-2026-08-31-165510/REQ-418-exploration.md
?? do-work/runs/work-2026-08-31-165510/REQ-418-handback.md
?? do-work/runs/work-2026-08-31-165510/REQ-418-plan.md
?? do-work/runs/work-2026-08-31-165510/REQ-418-remediation-handback.md
?? do-work/runs/work-2026-08-31-165510/REQ-418-rereview.md
?? do-work/runs/work-2026-08-31-165510/REQ-418-review.md
```

- **Detached worktree verdict:** FOREIGN — `/private/var/folders/2w/kw8sv6rd1z15yjykl787ryph0000gn/T/tmp.2qQUpYmVN0`, detached at `4b52ae4d`, dirty REQ-477-era fixture state. Leave byte-identical; not removable because dirty.

```text
M skills/do-work-toolbox/crew-members/general.md
 M skills/do-work/actions/review-work.md
 M skills/do-work/actions/work.md
 M skills/do-work/crew-members/general.md
 M skills/do-work/tools/do-work-cli/lessons-do-work-cli.md
 M skills/do-work/tools/do-work-cli/prime-do-work-cli.md
?? do-work/lessons-index.md
```

- **Detached worktree verdict:** FOREIGN — `/private/var/folders/2w/kw8sv6rd1z15yjykl787ryph0000gn/T/tmp.gYIGToUyFs`, detached at `4b52ae4d`, dirty REQ-477-era fixture state. Leave byte-identical; not removable because dirty.

```text
M skills/do-work-toolbox/crew-members/general.md
 M skills/do-work/actions/review-work.md
 M skills/do-work/actions/work.md
 M skills/do-work/crew-members/general.md
 M skills/do-work/tools/do-work-cli/lessons-do-work-cli.md
 M skills/do-work/tools/do-work-cli/prime-do-work-cli.md
?? do-work/lessons-index.md
```

- **Removed foreign fixture:** `/private/var/folders/2w/kw8sv6rd1z15yjykl787ryph0000gn/T/tmp.ykJ5cw2N9f` was present and dirty in the first survey, then the REQ-478 owner removed it before this prompt was written. Do not recreate or act on it.
- **First-ten-minutes heads-up:** interrupted board-test PID 20453 has exited; do not infer a board result from it. Existing board serve PIDs 8227/8406 predate this handoff. Preserve the six untracked REQ-418 run artifacts. Do not launch parallel builders in this shared checkout until REQ-478 is terminal and the mixed paths are reconciled.
