# REQ-743 workflow handback

Committed the add and browse forms, discovery surfaces, and gallery generation after verified suite writes.

- Branch — t3/implement-agent-forms-hooks
- SHA — a8031e04b9b0cf78a0d1dd7ea39f32112950fb1a
- Worktree clean at handback.

## Behavior

The add form loads the prompt-injection guardrail before reading the source, claims a fresh sample directory, records the agreed schema and presentation moves, and regenerates the page. It uses exact paths when following the project's commit rule. Browse regenerates the page and reports its path.

The README bootstrap and direct installer use the install-suite handler. The direct updater shell script, both Just recipes, and the do-work update agent route use the update-suite handler. They call the registered gallery renderer in-process after successful writes and verification, with the consumer repository root. This reads the newly installed toolbox samples through the trusted running binary. A refusal or failure emits one stderr line and preserves suite success. Cancelled and already-current runs skip generation.

## Files

- Forms and discovery — [ai-report.md](/Users/t2/.t3/worktrees/skill-do-work2/t3-903dfe0f-9c784110da458af1/skills/do-work-toolbox/actions/ai-report.md), [SKILL.md](/Users/t2/.t3/worktrees/skill-do-work2/t3-903dfe0f-9c784110da458af1/skills/do-work-toolbox/SKILL.md), [help.md](/Users/t2/.t3/worktrees/skill-do-work2/t3-903dfe0f-9c784110da458af1/skills/do-work-toolbox/actions/help.md), [ai-report-guide.md](/Users/t2/.t3/worktrees/skill-do-work2/t3-903dfe0f-9c784110da458af1/skills/do-work-toolbox/docs/ai-report-guide.md).
- Update output relay — [version.md](/Users/t2/.t3/worktrees/skill-do-work2/t3-903dfe0f-9c784110da458af1/skills/do-work/actions/version.md), [README.md](/Users/t2/.t3/worktrees/skill-do-work2/t3-903dfe0f-9c784110da458af1/README.md).
- Hook and Go tests — [suite_commands.go](/Users/t2/.t3/worktrees/skill-do-work2/t3-903dfe0f-9c784110da458af1/skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands.go), [suite_commands_test.go](/Users/t2/.t3/worktrees/skill-do-work2/t3-903dfe0f-9c784110da458af1/skills/do-work/tools/do-work-cli/internal/suiteinstall/suite_commands_test.go), [install_transaction_test.go](/Users/t2/.t3/worktrees/skill-do-work2/t3-903dfe0f-9c784110da458af1/skills/do-work/tools/do-work-cli/internal/suiteinstall/install_transaction_test.go).
- Entry-point probes — [install-suite-behavior.sh](/Users/t2/.t3/worktrees/skill-do-work2/t3-903dfe0f-9c784110da458af1/_dev/tests/install-suite-behavior.sh), [update-script-behavior.sh](/Users/t2/.t3/worktrees/skill-do-work2/t3-903dfe0f-9c784110da458af1/_dev/tests/update-script-behavior.sh).

## Tests and phase evidence

- RED before production edits — TestReportGalleryRebuildsFromTheWrittenSuite found the first-install page missing. TestReportGalleryWriteProblemsWarnOnceWithoutRecoveringTheSuite found the gallery warning missing for both handlers and both blocked paths.
- GREEN for the hook — both tests passed on Go 1.24.0 with a temporary Go overlay supplying a gallery command double. It exercised confirmed install, changed-version update, cancellation, already-current update, refusal, filesystem failure, one stderr line, intact JSON, and installed version retained. The overlay was removed.
- Existing suiteinstall tests passed on Go 1.27.4 with the two new renderer-dependent tests excluded. go vet ./internal/suiteinstall, contract-regressions.sh, action-shell-blocks.sh, shipped-package-reference-contract.sh, shell syntax checks, and git diff --check passed.
- The isolated installer lane reached the bootstrap assertion and found all four gallery links missing. The updater lane reached a real version-advancing update through the direct shell and both Just routes and found only the three expected missing-gallery assertions. Full GREEN for these lanes and the Go tests with the real renderer belongs to the combined-slice integration.

## Rollout issue

An updater binary installed before this hook runs its own transaction and therefore cannot generate the gallery while upgrading itself. A fresh bootstrap or subsequent file-writing update run by a gallery-aware version generates the page automatically. A project that has just crossed from a pre-hook version can generate the page with ai-report gallery.
