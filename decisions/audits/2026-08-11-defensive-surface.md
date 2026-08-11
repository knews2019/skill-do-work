# Shipped Defensive Surface: Delete-or-Test Audit

**Ran:** 2026-08-11 · **Scope:** all shipped files under `skills/` · **Rubric:** what incident earned this, and is the fix still cheaper than the surface it added?

This inventory treats a defensive *layer* as an executable guard/fallback family or an explicit prose section (`Rules`, `Common Rationalizations`, `Red Flags`, `Warnings`, recovery). Repeated branches implementing one atomic transaction are one layer, not one row per `if`. The focused probe verifies that every current explicit prose surface and every shipped shell source remains represented here.

## Executable layers

| Location | Incident earned | Disposition | Evidence |
|---|---|---|---|
| `skills/do-work/hooks/session-start.sh` | Missing/reformatted version metadata aborted the startup banner; relative hook paths also resolved from the wrong cwd. | KEEP — minimal fallback after REQ-166. | `_dev/tests/session-start-hook-behavior.sh`; `_dev/tests/action-shell-blocks.sh` |
| `skills/do-work/scripts/show-commit-diff.sh` | Plain `git show` hid real changes on worktree merge commits. | KEEP — canonical merge-aware display. | `_dev/tests/prescribed-shell-scripts-behavior.sh` real-merge case; `_dev/tests/action-shell-blocks.sh` |
| `skills/do-work/scripts/add-local-git-exclude.sh` | Root-built `.git` paths and cwd-relative ignore probes diverged in worktrees and subdirectories. | KEEP — canonical local-exclude publication. | `_dev/tests/prescribed-shell-scripts-behavior.sh` subdirectory case; `_dev/tests/action-shell-blocks.sh` |
| `skills/do-work/scripts/atomic-download.sh` | Failed transfers left non-empty partial final files that presence-gated consumers accepted. | KEEP — private download plus atomic rename. | `_dev/tests/prescribed-shell-scripts-behavior.sh` partial-publication case; `_dev/tests/action-shell-blocks.sh` |
| `skills/do-work/scripts/capture-screenshot.sh` | Concurrent screenshot dispatches could verify one source and publish another dispatch's shared temporary copy. | KEEP — unique private copy plus no-clobber link. | `_dev/tests/prescribed-shell-scripts-behavior.sh` coordinated-race case; `_dev/tests/action-shell-blocks.sh` |
| `skills/do-work/scripts/run-blocked-check.sh` | macOS lacks GNU `timeout`, but blocked probes still require the same 30-second bound and status 124. | KEEP — tested polling fallback. | `_dev/tests/prescribed-shell-scripts-behavior.sh` portable-timeout case; `_dev/tests/action-shell-blocks.sh` |
| `skills/do-work/scripts/protected-inventory.sh` | Re-running inventory after index changes made an earlier secret-quarantined rename destination readable. | KEEP — run-level once-X quarantine around existing checks. | `_dev/tests/prescribed-shell-scripts-behavior.sh` quarantine-association case; `_dev/tests/action-shell-blocks.sh` |
| `skills/do-work/scripts/stage-exact-deletion.sh` | Secret-shaped deletions must be staged without reconstructing content or staging an adjacent path. | KEEP — cached metadata-only exact guard. | `_dev/tests/prescribed-shell-scripts-behavior.sh` pathological-name case; `_dev/tests/action-shell-blocks.sh` |
| `skills/do-work-knowledge/scripts/lexical-memory-recall.sh` | Raw query interpolation made apostrophes syntax and lost deterministic attribution/scoring. | KEEP — data-only sanitization and bounded ranking. | `_dev/tests/prescribed-shell-scripts-behavior.sh` raw-query case; `_dev/tests/action-shell-blocks.sh` |
| `skills/do-work-knowledge/scripts/install-memory-hooks.sh` | A single hook-name gate duplicated one partial install or permanently skipped its missing sibling. | KEEP — independent append gates with rollback. | `_dev/tests/prescribed-shell-scripts-behavior.sh` partial-merge case; `_dev/tests/action-shell-blocks.sh` |
| `skills/do-work-toolbox/scripts/generate-report-image.sh` | Opportunistic backends could return success without the exact output, and agentic fallback escaped its scratch cwd. | KEEP — exact-path verification and opt-in quarantine. | `_dev/tests/prescribed-shell-scripts-behavior.sh` direct-backend case; `_dev/tests/action-shell-blocks.sh` |
| `skills/do-work-toolbox/scripts/install-last30days.sh` | Skill-file-only detection accepted partial installs whose ignore or Python guarantees were missing. | KEEP — install/repair/full verification in one home. | `_dev/tests/prescribed-shell-scripts-behavior.sh` fixture-source case; `_dev/tests/action-shell-blocks.sh` |
| `skills/do-work-knowledge/hooks/memory-session-start.sh` | Raw session captures could re-enter trusted context before the injection guard loaded; missing memory/log state must not break startup. | KEEP. | `_dev/tests/contract-regressions.sh` startup-injection fixtures; `_dev/tests/action-shell-blocks.sh` |
| `skills/do-work-knowledge/hooks/memory-stop-capture.sh` | Stop-hook recursion, credential-shaped transcript text, invalid JSON, absent jq/iconv/hash tools, and capture-size overflow occurred at an optional telemetry boundary. | KEEP. | `_dev/tests/contract-regressions.sh` hook behavior/redaction/budget fixtures; `_dev/tests/action-shell-blocks.sh` |
| `skills/do-work/tools/checks/archive-collision.sh` | A queue REQ with an archived twin could be reprocessed and corrupt lineage. | KEEP. | `_dev/tests/contract-regressions.sh`; `_dev/tests/action-shell-blocks.sh` |
| `skills/do-work/tools/checks/associate-files.sh` | Rename-safe inventories and terminal-status aliases lost file ownership; missing option values could hang. | KEEP. | `_dev/tests/contract-regressions.sh` association fixtures; `_dev/tests/action-shell-blocks.sh` |
| `skills/do-work/tools/checks/blanked-req-scan.sh` | Six real archived REQs were truncated by unsafe commit-hash write-back and were recoverable only before Git GC. | KEEP. | `_dev/tests/record-commit-hash-guards.sh`; `_dev/tests/action-shell-blocks.sh` |
| `skills/do-work/tools/checks/preflight.sh` | Quoted/spaced/NUL-bearing dirty paths were misread and baseline failures were misattributed to active work. | KEEP. | `_dev/tests/contract-regressions.sh` preflight fixtures; `_dev/tests/action-shell-blocks.sh` |
| `skills/do-work/tools/checks/qualify.sh` | Empty/unresolvable worktree ranges could vacuously pass and incomplete summaries could reach review. | KEEP. | `_dev/tests/contract-regressions.sh`; `_dev/tests/action-shell-blocks.sh` |
| `skills/do-work/tools/checks/record-commit-hash.sh` | Direct shell redirection blanked tracked REQ files while recording implementation hashes. | KEEP. | `_dev/tests/record-commit-hash-guards.sh`; `_dev/tests/action-shell-blocks.sh` |
| `skills/do-work/tools/checks/scope-drift.sh` | Undeclared and declared-but-unused files escaped the Route B/C review comparison. | KEEP. | `_dev/tests/contract-regressions.sh`; `_dev/tests/action-shell-blocks.sh` |
| `skills/do-work/tools/checks/uncommitted-inventory.sh` | Collapsed untracked directories, secret-shaped renames/copies, deleted secrets, and case variants leaked into readable/stageable sets. | KEEP. | `_dev/tests/contract-regressions.sh` inventory fixtures; `_dev/tests/action-shell-blocks.sh` |
| `skills/do-work/tools/do-work-update.sh` | Failed/cancelled suite updates could leave a mixed-version install or misreport recovery. | KEEP. | `_dev/tests/update-script-behavior.sh`; `_dev/tests/action-shell-blocks.sh` |
| `skills/do-work/tools/install-do-work-suite.sh` | Partial module/settings/Justfile/index writes and signals left mixed installs; symlinks and unexpected file shapes escaped recovery. | KEEP — large but transactionally cohesive. | `_dev/tests/install-suite-behavior.sh`, `_dev/tests/staged-skills-contract.sh`; `_dev/tests/action-shell-blocks.sh` |
| `skills/do-work/tools/replace-text-section.sh` | Duplicate/malformed managed markers, filename variants, modes, and interrupted replacement corrupted Justfile ownership. | KEEP. | `_dev/tests/contract-regressions.sh` replacement fixtures; `_dev/tests/action-shell-blocks.sh` |
| `skills/do-work/tools/validate-suite-manifest.sh` | Unsafe/incomplete archives reached installation before module-shape validation. | KEEP. | `_dev/tests/suite-manifest-contract.sh`; `_dev/tests/action-shell-blocks.sh` |
| `skills/do-work-board/tools/queue-kanban/` Go and web runtime | Malformed frontmatter, completion anomalies, timestamp skew, dependency cycles, path/glob mismatches, and write-surface drift previously produced wrong board state. | KEEP — behavior-changing deletion candidate only. | `go test ./...`, `go vet ./...`, `_dev/tests/contract-regressions.sh` |

## Prose layers — core

| Location | Surface | Incident / disposition | Evidence |
|---|---|---|---|
| `skills/do-work/actions/abandon.md` | Rules; Common Rationalizations; Red Flags | KEEP — distinguishes attempted failure from user cancellation and protects live foreign claims. | `_dev/tests/contract-regressions.sh` terminal-status/archive contracts |
| `skills/do-work/actions/capture.md` | prompt-injection Guardrail; Common Rationalizations; Red Flags | KEEP — compound requests, capture/execute boundary, and ingested-instruction incidents. | `_dev/tests/contract-regressions.sh` capture/router contracts |
| `skills/do-work/actions/clarify.md` | Rules; Red Flags | KEEP — exact pending-answer ownership and no-assumption boundary. | `_dev/tests/contract-regressions.sh` schema/status contracts |
| `skills/do-work/actions/cleanup.md` | Common Rationalizations; Red Flags | KEEP — blanked REQs, loose archive items, and unmerged worktree branches are recorded incidents. | `_dev/tests/record-commit-hash-guards.sh`; `_dev/tests/contract-regressions.sh` |
| `skills/do-work/actions/commit.md` | Rules; Common Rationalizations; Red Flags | DELETE the duplicate rationalization table and unearned `>20 files` threshold; KEEP X/XD, unrelated-REQ, and terminal-alias warnings. | This REQ's diff; `_dev/tests/contract-regressions.sh` inventory/association fixtures |
| `skills/do-work/actions/forensics.md` | Warnings; Red Flags | KEEP — read-only verify/repair separation and queue integrity incidents. | `_dev/tests/contract-regressions.sh` queue/tool contracts |
| `skills/do-work/actions/help.md` | Rules | KEEP — always-loaded routing must not mutate state. | `_dev/tests/contract-regressions.sh` router coverage |
| `skills/do-work/actions/kb-lessons-handoff.md` | Rules; Common Rationalizations; Red Flags | KEEP — consent, one-REQ provenance, and untrusted Lessons content boundaries. | `_dev/tests/shipped-package-reference-contract.sh`; `_dev/tests/contract-regressions.sh` |
| `skills/do-work/actions/review-work.md` | Common Rationalizations; Red Flags | KEEP — score/follow-up gates and stale-restatement sweep trace to review misses. | `_dev/tests/contract-regressions.sh` review/follow-up contracts |
| `skills/do-work/actions/roadmap.md` | Rules; Common Rationalizations; Red Flags | KEEP — read-only roadmap mutations previously split queue ownership. | `_dev/tests/contract-regressions.sh` action/router contracts |
| `skills/do-work/actions/verify-requests.md` | Common Rationalizations; Red Flags | DELETE — generic “be thorough / don't round” advice repeats the scoring steps and verification checklist; no named incident found. | This REQ's diff; focused audit probe |
| `skills/do-work/actions/version.md` | Red Flags | KEEP — three version surfaces have drifted independently. | `_dev/tests/contract-regressions.sh`, `_dev/tests/staged-skills-contract.sh` |
| `skills/do-work/actions/work.md` and `skills/do-work/actions/work-reference.md` | Rules; Common Rationalizations; Red Flags; Failure Classification | KEEP — queue claims, crash recovery, merge ranges, follow-up gates, and cleanup assertions all trace to archived pipeline failures. | `_dev/tests/contract-regressions.sh`; queue-kanban Go tests |
| `skills/do-work/crew-members/background-agents.md` | Known Failure Mode & Recovery | KEEP — conversation-only handback disappeared when orchestrator sessions died. | `_dev/tests/shipped-package-reference-contract.sh`; `_dev/tests/action-shell-blocks.sh` |
| `skills/do-work/crew-members/clear-questions.md` | Red Flags | KEEP — approval/input gates prevent builders from inventing user intent. | `_dev/tests/shipped-package-reference-contract.sh` |
| `skills/do-work/crew-members/security.md` | OWASP failure sections | KEEP — behavior-changing security baseline, outside safe prose-restatement deletion. | `_dev/tests/shipped-package-reference-contract.sh` |

## Prose layers — toolbox

| Location | Surface | Incident / disposition | Evidence |
|---|---|---|---|
| `skills/do-work-toolbox/actions/ai-report.md` | Rules; Common Rationalizations; Red Flags | KEEP — fabricated/file-URL screenshots, hidden generated evidence, prompt injection, and layout regressions are named incidents. | `_dev/tests/contract-regressions.sh`; `_dev/tests/shipped-package-reference-contract.sh` |
| `skills/do-work-toolbox/actions/code-review.md` | Rules; Common Rationalizations; Red Flags | KEEP Rules (scope/output contract); DELETE generic review-advice tables, all duplicated by steps/checklist and with no named incident. | This REQ's diff; focused audit probe |
| `skills/do-work-toolbox/actions/deep-explore.md` | Rules; Red Flags | KEEP — bounded read-only exploration and evidence-location contract. | `_dev/tests/shipped-package-reference-contract.sh` |
| `skills/do-work-toolbox/actions/inspect.md` | Rules; Common Rationalizations; Red Flags | DELETE generic duplicate rationale and hollow/debug restatements; KEEP X/XD, association, and Already Committed state warnings. | This REQ's diff; `_dev/tests/contract-regressions.sh` inventory/association fixtures |
| `skills/do-work-toolbox/actions/install.md` | Rules; Common Rationalizations; Red Flags | KEEP — partial download, global-install, credential, toolchain, and trust-gate incidents. | `_dev/tests/action-shell-blocks.sh`; `_dev/tests/contract-regressions.sh` |
| `skills/do-work-toolbox/actions/note.md` | Rules | KEEP — explicit mutation boundary for a note action. | `_dev/tests/shipped-package-reference-contract.sh` |
| `skills/do-work-toolbox/actions/present-work.md` | Red Flags | KEEP — artifact verification, commit evidence, and audience separation incidents. | `_dev/tests/contract-regressions.sh`; `_dev/tests/shipped-package-reference-contract.sh` |
| `skills/do-work-toolbox/actions/prime.md` | Red Flags | KEEP — prime packaging/overwrite boundaries. | `_dev/tests/shipped-package-reference-contract.sh` |
| `skills/do-work-toolbox/actions/quick-wins.md` | Rules; Common Rationalizations; Red Flags | KEEP Rules (read-only/scope contract); DELETE generic scan-thoroughness tables with no named incident. | This REQ's diff; focused audit probe |
| `skills/do-work-toolbox/actions/scan-ideas.md` | Rules | KEEP — report-only boundary prevents idea scans from mutating the backlog. | `_dev/tests/shipped-package-reference-contract.sh` |
| `skills/do-work-toolbox/actions/slop-check.md` | Rules; Common Rationalizations; Red Flags | KEEP — read-only rewrite confirmation, evidence, fact preservation, and report-size constraints are action-specific. | `_dev/tests/shipped-package-reference-contract.sh` |
| `skills/do-work-toolbox/actions/stray-check.md` | Rules; Warnings; Common Rationalizations; Red Flags; Guarded fixes | KEEP — unrecoverable deletion, tracked secrets/build output, dynamic entrypoints, and per-file inventory incidents. | `_dev/tests/contract-regressions.sh`; `_dev/tests/action-shell-blocks.sh` |
| `skills/do-work-toolbox/actions/tidy-repo.md` | Rules; Red Flags | KEEP — report-only repo hygiene and tracked/untracked boundaries. | `_dev/tests/shipped-package-reference-contract.sh` |
| `skills/do-work-toolbox/actions/tutorial.md` | Rules; Red Flags | KEEP — generated tutorial verification and evidence contract. | `_dev/tests/shipped-package-reference-contract.sh` |
| `skills/do-work-toolbox/actions/ui-review.md` | Rules; Common Rationalizations; Red Flags | KEEP Rules (review scope/output contract); DELETE generic checklist-thoroughness tables with no named incident. | This REQ's diff; focused audit probe |
| `skills/do-work-toolbox/actions/validate-feedback.md` | Rules; Common Rationalizations; Red Flags; Guardrails | KEEP — evidence-before-verdict and capture/execute boundaries; prospective surface-cost strengthening is REQ-169. | `_dev/tests/contract-regressions.sh`; REQ-169 |
| `skills/do-work-toolbox/crew-members/background-agents.md` | Known Failure Mode & Recovery | KEEP — same durable handback incident as core. | `_dev/tests/shipped-package-reference-contract.sh`; `_dev/tests/action-shell-blocks.sh` |
| `skills/do-work-toolbox/crew-members/clear-questions.md` | Red Flags | KEEP — approval/input gate. | `_dev/tests/shipped-package-reference-contract.sh` |
| `skills/do-work-toolbox/crew-members/security.md` | OWASP failure sections | KEEP — behavior-changing security baseline. | `_dev/tests/shipped-package-reference-contract.sh` |

## Prose layers — knowledge and board

| Location | Surface | Incident / disposition | Evidence |
|---|---|---|---|
| `skills/do-work-board/actions/board.md` | Rules; Common Rationalizations; Red Flags | KEEP — stale binary, wrong-root discovery, absent Go, and local artifact incidents. | `go test ./...`; `_dev/tests/contract-regressions.sh` |
| `skills/do-work-knowledge/actions/bkb.md` | Red Flags | KEEP — KB schema/ingest/destructive-operation boundaries. | `_dev/tests/shipped-package-reference-contract.sh` |
| `skills/do-work-knowledge/actions/dream.md` | Rules; Common Rationalizations; Red Flags | KEEP — destructive consolidation, inbound links, source immutability, consent, and prompt-injection incidents. | `_dev/tests/shipped-package-reference-contract.sh` |
| `skills/do-work-knowledge/actions/interview.md` | Rules; Common Rationalizations; Red Flags | KEEP — confirmation checkpoints, archive-before-reset, template/export gates. | `_dev/tests/shipped-package-reference-contract.sh` |
| `skills/do-work-knowledge/actions/memory-value.md` | Rules; Red Flags | KEEP — read-only experiment reporting and ledger boundaries. | `_dev/tests/contract-regressions.sh` memory contracts |
| `skills/do-work-knowledge/actions/memory.md` | Rules; Common Rationalizations; Red Flags | KEEP — bounded working memory, bootstrap sentinel, and best-effort ledger behavior. | `_dev/tests/contract-regressions.sh` memory fixtures |
| `skills/do-work-knowledge/actions/prompts.md` | Rules; Common Rationalizations; Red Flags | KEEP — preview/execute split, shipped-vs-project trust, and injected prompt bodies. | `_dev/tests/shipped-package-reference-contract.sh` |
| `skills/do-work-knowledge/actions/setup-memory.md` | Rules | KEEP — setup ownership and local-ignore boundary. | `_dev/tests/contract-regressions.sh` memory contracts |
| `skills/do-work-knowledge/crew-members/background-agents.md` | Known Failure Mode & Recovery | KEEP — same durable handback incident as core. | `_dev/tests/shipped-package-reference-contract.sh`; `_dev/tests/action-shell-blocks.sh` |
| `skills/do-work-knowledge/crew-members/clear-questions.md` | Red Flags | KEEP — approval/input gate. | `_dev/tests/shipped-package-reference-contract.sh` |
| `skills/do-work-knowledge/crew-members/security.md` | OWASP failure sections | KEEP — behavior-changing security baseline. | `_dev/tests/shipped-package-reference-contract.sh` |
| `skills/do-work-knowledge/prompts/prompt-kit-step0-pen-and-paper-exercises-to-prepare-prompt.md` | Rules; Common Rationalizations; Red Flags | KEEP — offline authorship/contamination and no-inference contract. | `_dev/tests/shipped-package-reference-contract.sh` |
| `skills/do-work-knowledge/prompts/prompt-kit-step1-four-discipline-diagnostic.md`; `skills/do-work-knowledge/prompts/prompt-kit-step2-personal-context-doc.md`; `skills/do-work-knowledge/prompts/prompt-kit-step3-spec-engineer.md`; `skills/do-work-knowledge/prompts/prompt-kit-step4-intent-and-delegation-framework.md`; `skills/do-work-knowledge/prompts/prompt-kit-step5-eval-harness.md`; `skills/do-work-knowledge/prompts/prompt-kit-step6-constraint-architecture.md` | Rules; Red Flags; failure-mode extraction | KEEP — each prompt's handoff/schema/constraint contract; behavior-changing template content. | `_dev/tests/shipped-package-reference-contract.sh` |

## Deletes applied

- Removed four complete decorative Common Rationalizations + Red Flags pairs from general review actions (`code-review`, `quick-wins`, `ui-review`, `verify-requests`). Their steps, scoring/category schemas, output formats, and verification checklists remain unchanged.
- Removed `commit`'s duplicate two-row rationalization section and its arbitrary “single commit with >20 files” warning. Semantic purpose—not file count—remains the grouping rule.
- Removed `inspect`'s duplicate rationalization section and three generic Red Flags; retained the X/XD, association, and Already Committed state-specific warnings.

## Behavior-changing candidates retained

No observable action behavior was deleted. The installer transaction, queue/worktree recovery, memory capture safety, security crew rules, destructive-operation confirmations, and board parser tolerances are all expensive surfaces, but removing any would change outcomes rather than merely remove robustness theater. This audit records them as KEEP; a future request must name the desired behavior change and replay the cited incident before narrowing one.
