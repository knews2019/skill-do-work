# REQ-724 integration report (moving the maintainer's machine-local standing rules into the crew-member files)

- **Release:** 0.322.0 (minor, impact-rule-change). Title: "Crew Files Carry the Rules for Status Explanations, User Decisions and Pasted Commands".
- **Archive:** `do-work/archive/REQ-724-move-global-standing-rules-into-skill.md` (flat: UR-175 stays open for REQ-725).
- **Merge:** base `60c8a8f8`. First merge `47565ced`. After review fixes, re-merge `0074c2744688f020c37b352325a96a1cbb01b605` (the recorded commit hash). Finalization commit `7c15a2dd`.
- **Review:** first review 88% Approve, with F1 and F2 marked required. F1: the toolbox and knowledge packages' lock-step copies of `clear-questions.md`, `coding-guardrails.md` and `background-agents.md` lacked the new lines. F2: the DECIDE & STATE pointer left out the tier's "clear from context" condition. I fixed F1, F2 and the report-only F4 (the stall-check sentence) on the builder branch as `e8bbcd7c` and re-merged. The GREEN probe now accepts one home per package (D-14). Delta re-review: 92% Approve. Report only: F3 (no precedence over an action's own rule), F5 (G6 partly unverified), F6, N1 (probe nit).
- **F6, for the maintainer:** the numeric-target, redesign-UI and model rules load only during implementation or fan-out. Once the global file is deleted, they stop applying in plain sessions. Decide this before that delete step.
- **Gate:** green at both merges (135 s at load 4.61, 127 s at load 4.29). The green gate is recorded for `0074c274`.
- **Heavy lanes at 0074c274:** all six ran and exited 0. The planner selected all six because the range holds a run-artifact commit and that made coverage uncertain. Times: queue-kanban-javascript 7 s, queue-kanban-browser 78 s, do-work-cli-integrations 64 s, staged-skills 36 s, updater 66 s, installer 27 s. The installer lane includes `install-suite-behavior.sh`, which checks the link to `communication-style.md`.
- **Post-finalize:** `contract-regressions.sh` exit 0. Worktree removed, branch deleted, worktrees pruned.
- **Timing events:** builder-work, handback-merge, review, and verification-gate three times (gate, re-gate, heavy drain).
- **Lessons:** no new bullet. The builder proposed none, and the existing `alternate-writer-contract-drift` family covers the F1 miss.
- **Discovered tasks:** none.
- **Dirty in the main checkout:** nothing. This report and the progress log go in one plain commit.
