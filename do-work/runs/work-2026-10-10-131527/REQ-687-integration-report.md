# REQ-687 integration report

- Release 0.305.110, "ai-report Writes a Decision Brief for Unfinished Work With --kind proposal or root-cause". Archive `do-work/archive/REQ-687-ai-report-kind-proposal-root-cause.md` (flat, UR-153 stays open). Finalization commit `18167d00`, cleanup_complete on the first try.
- Range: pre `3bd6a668`. First merge `a0ebfea9`. Re-merge `8c0cf7ed` after the review fixes (recorded as `commit:`). My fix on the builder branch: `84d23733`.
- Merge: four files conflicted with REQ-655/656/657. I kept every REQ's text, and there is still one ai-report row in SKILL.md. Seam: the `--kind` Input paragraph had auto-merged inside the Revise form steps without a conflict, so I moved it to the top-level Input text.
- Before qualify I fixed two stale paths in the REQ (they now point under `archive/UR-144/`). qualify and scope-drift passed on the first call.
- Review: 89%, Approve. I applied F1 to F5 with the reviewer's exact text. F1 (Important): `--kind` lifted only target resolution, so the shared Archive Evidence Sweep would still stop on an open REQ. Now only Safety Load Order and Collision-Safe Publication are inherited. F2 to F5 are small fixes to the guide, When to Use, the Steps 3-4 rule and a checklist line. No second review, because the delta is the reviewer's text only.
- Report only: F6 (`ai-report.md:31` Input sentence leaves out `--kind`), F7 (`architecture-report.md:135` says "presents completed work"), F8 (no rule for a completed target with `--kind`), and nits F9 to F13.
- Gate: exit 0 at `a0ebfea9` (113 s, load 1.92) and at `8c0cf7ed` (117 s, load 2.19). The GREEN probe passed at both. Green record at `8c0cf7ed`.
- Heavy lanes: staged-skills, executed, exit 0: 44 s and then 33 s. Drain checkouts removed with `--force` (seeded tsv).
- Timing: 7 events (2 merge, 4 gate and drain, 1 review), folded. No builder-work event, because the hand-back landed long before.
- Lessons: the builder proposed none, so no satellite bullet. `## Lessons Learned` notes that "only X is inherited" is the rule, and that a merge with no conflict can still put a paragraph in the wrong subsection.
- Cleanup: builder worktree removed, branch deleted with `-d`, worktrees pruned. `contract-regressions.sh` passes.
- Left dirty: only this report (untracked).
- Report-only builder tasks: the `architecture-report.md` line, and the guide and reference `:3` intros, which describe completed work only.
