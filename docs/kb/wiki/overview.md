# Knowledge Base Overview

Last updated: 2026-09-04 | Articles: 245 | Topic Clusters: 11

The Do-Work Suite Knowledge Base is a compiled, interlinked Markdown wiki capturing architecture, lessons, operational patterns, and historical decisions across the four-skill suite: `do-work`, `do-work-board`, `do-work-knowledge`, and `do-work-toolbox`.

## Core Knowledge Domains

1. **Queue Orchestration & Task Lifecycle**: The `do-work` orchestrator, REQ lifecycle, claim/release semantics, failed-task resolution, lightweight notes, and deterministic state transitions.
2. **Worktree & Parallel Dispatch**: Multi-agent fan-out building, git worktree isolation, write-set disjointness checks, wave dispatch, and multi-checkout coordination.
3. **Checkpoints & Crash Recovery**: Session-start recovery, in-progress record protection, checkpoint writer labels, and defense against state poisoning.
4. **Kanban Board & UI**: The framework-free `queue-kanban` Go tool, web dashboard, drawer inspection, filtering, theme management, and clipboard integration.
5. **Timeline & Metrics**: Timeline duration visualization, period zoom, movement thresholds, and p50 duration estimators.
6. **Prescribed Shell & CLI**: Portable shell patterns, standard Bash trap structures, ShellCheck conformance, and native `do-work-cli` command implementations.
7. **Verification & Testing**: The canonical `maintainer-verify` suite, strict browser probes, contract regression tests, qualify checks, and earned-defense forensics.
8. **Metadata & Timestamps**: Strict UTC ISO-8601 timestamp canonicalization, reservation markers, calibration logs, and cross-package citation rules.
9. **Presentation & Reporting**: Automated stakeholder reports (`ai-report`), portfolio generators (`present-work`), video presentations (`present-video`), and architecture decision records.
10. **Knowledge & Memory**: The Build Knowledge Base (`bkb`) compiler system, usage ledgers, structured interviews, prompt libraries, and agent crew coordination.
11. **Suite & Package Architecture**: Four-skill modular packaging, manifest contracts, bridge updaters, and managed text replacements.
