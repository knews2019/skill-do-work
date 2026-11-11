# REQ-718 compaction investigation execution result

Combined acceptance is incomplete. Claude has a verified live compaction result; actual Codex/T3 live compaction is Untested because no supported active-session control or observable boundary was established. No instruction-retention failure or restoration fix is inferred. Step 9 finalized the request as failed/environment, preserving the partial evidence without claiming completion.

## Integrated identities

- Integration branch: main. Original claimed request retained; no recovery or Step 10 performed.
- Range: 8e7cf62aa4c75ab1e82f6e9912d174cc26fbc547..a157fd164214d36ac5f12b72949c56a1e0facc95.
- Claude: t3/execute-claude-slice-work, /Users/t2/.t3/worktrees/skill-do-work2/t3-491e3e50, worker HEAD 993c3befca2af4eeae8aad2fdba386ab7d16438a, serial merge c120c3f853e28a504aaa82dd95397d55a1db97c6.
- Codex: t3/execute-codex-brief, /Users/t2/.t3/worktrees/skill-do-work2/t3-d6e60aeb, worker HEAD b19b6f693e7f0f2aa694cac87f28c1acbfc4ae05 (includes 8adbb47f), serial merge a157fd164214d36ac5f12b72949c56a1e0facc95.
- Both GPT-6.1 Sol high workers completed. All source/evidence implementation stayed in tool-created worktrees. Main-tree writes were lifecycle/integration records only.

## Observed evidence

Claude Code 2.1.322 / claude-fable-5-1, explicit safe/read-only configuration: eight full file-read results verified against source; manual compact boundary 36,066 → 11,389 tokens, 41,513 ms. Continuation made zero tool calls. Its sumShipmentWeights function complied with the loaded naming rule, also explicitly present in the synthetic summary. This does not establish full verbatim file retention or causal retention from the crew file alone. Ordinary project/global settings, automatic/repeated compaction and other rules remain Untested.

Actual T3 server 0.1.7 / codex / gpt-6.1-sol / high: phase and crew read events plus baseline descriptive filename captured. No supported active-worker compaction operation found among 65 exposed T3 tools. PATH CLI 0.167.0 is not verified as the backing binary; default daemon version socket absent. Backing provider build/settings and post-compaction behavior remain unverified. No standalone substitute ran.

Artifacts: _dev/experiments/REQ-718/claude/experiment-report.md and _dev/experiments/REQ-718/codex-t3/experiment-report.md; 17 committed files total. Full source/JSON validation, ten Codex fingerprints, eight Claude read contents and the Claude committed capture verifier pass. Full merge-range whitespace passes.

## Gates and review

Pre-build canonical bash _dev/tests/maintainer-verify.sh exited 0; advance preflight and green-gate satisfied. Final canonical first run exited 1 in existing strict browser test timeout (~29 seconds); mandatory same-argv retry exited 0; advance scope, focused-test and green-gate all satisfied. Focused advance run-blocked-check and scope-drift satisfied. Literal Scope comparator required expanding the authorized directory globs to exact files; no substantive boundary changed. QUALIFY-LIBRARY-OUTPUT is intentional standalone CLI JSON reporting, judged against both main guards.

Inline orchestrated review: 78%, Acceptance Partial. Terminal judgment is failed/environment because required actual Codex proof remains impossible with exposed capabilities; no self-resolving precondition or source remediation is established. No release (maintainer artifacts only), no hooks installed/global changes/pushes/PRs. The Claude worker reports a command-local hooks bypass for its commit; this procedural deviation is recorded in Review and no setting persisted.

## Finalization and cleanup

Canonical advance returned global success and exactly one matching finalization record with cleanup_complete, terminal failed, empty blocked_paths and reason_codes. Archive: do-work/archive/REQ-718-compaction-instruction-retention.md. Finalization commit: af781caf0cc5d239dd74f44805524cec843af6be. Supplied implementation: a157fd164214d36ac5f12b72949c56a1e0facc95. No recovery call was made.

Both clean merged tool-managed worktrees were removed without force; both operative branches deleted using git branch -d from main. Full-gate lock removed. Both delegated runs completed and their handbacks were consumed from disk. No probe/gate command remains. Existing unrelated Kanban server was left intact.

Coordinator retains Step 10, final session checkpoint and run-directory consumption/cleanup. Run scratch and handbacks remain intentionally unstaged for that handoff; this execution report and manifest are saved separately. No new queue item or follow-up was minted. The incomplete experiment needs supported actual T3 compaction control/build visibility before it can be repeated.

Timing summary records the final retry (218 observed seconds); the first failed gate (130 seconds) was folded earlier and is retained in the first-gate log rather than added to that replacement summary.
