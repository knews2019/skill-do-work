# Codex/T3 instruction-retention handback

The actual Codex/T3 baseline and evidence package are committed. **Live compaction and all post-compaction checks remain Untested.** No supported control targeting compaction of this running worker was established. This is a blocked experiment result, not acceptance completion.

## Operative identity and commits

- Branch: `t3/execute-codex-brief`
- Worktree: `/Users/t2/.t3/worktrees/skill-do-work2/t3-d6e60aeb`
- Base: `a867b1a08d38dfb5605602a6a0ac5408c0368eb9`
- Evidence commit: `8adbb47f` (Record actual Codex T3 compaction probe blockers)
- Whitespace correction commit / operative HEAD: `b19b6f693e7f0f2aa694cac87f28c1acbfc4ae05`
- Worker: <t3-thread://v2/344b608c-b133-4724-a50e-9ac6cd780e62/thread:delegated-task:command%3Amcp%3A38ed711b-666b-4b5e-8e14-3e5fa8a3df6c%3Adelegate-task%3Af8b09f80-e2d9-401b-81d3-43384a803871>

Integrator: consume both commits, or merge this branch at the operative HEAD. No push or PR was created.

## Exact file manifest

All six files are new, under `_dev/experiments/REQ-718/codex-t3/`:

| File | Action and reviewed content |
|---|---|
| `experiment-report.md` | Added scoped conclusions, baseline, blockers, coverage, uncertainty and cleanup |
| `baseline-observation.md` | Added observable descriptive-name baseline with explicit limits |
| `probe-prompts.md` | Added original durable brief, actual discovery commands, baseline procedure and unexecuted continuation |
| `runtime-evidence.json` | Added scoped T3 configuration, server version, CLI results and 65-tool inventory |
| `transcript-evidence.json` | Added exact selected actual activity excerpts with IDs, timestamps and truncation flags |
| `source-fingerprints.json` | Added SHA-256 fingerprints for ten read instruction files |

No production code, tests, hooks, runtime instructions, release metadata, or working REQ changed. Main-tree writes are only this handback and the authorized progress log. These two lifecycle files are not committed on the worker branch.

## Focused checks and actual outcomes

- Evidence validation: exit 0. Three JSON files parsed; ten instruction fingerprints matched; actual file-change event 28 and Naming for Reach excerpt verified; model, inventory count and daemon-version blocker matched; report marks unexecuted checks Untested.
- Scope validation: exit 0. Exact six-file staged scope, valid UTF-8/newlines, descriptive artifact names, and limited common credential-marker scan. Retained fields were also manually inspected; this is not an exhaustive secret scanner.
- Initial `git diff --cached --check`: exit 2 for five quoted blank lines with trailing spaces. The first commit landed before correction. The follow-up commit removes the spaces.
- Corrected staged `git diff --cached --check`: exit 0.
- Final `git diff a867b1a08d38dfb5605602a6a0ac5408c0368eb9 HEAD --check`: exit 0.
- Full-range `git diff ... --stat`: six added files, 668 insertions. Every file was reviewed for its purpose, evidence claims, and isolation.
- Final `git status --short`: empty.
- No code tests or native build required for this evidence-only change. The canonical gate belongs to the integrator; pre-build exit 0 was supplied in the dispatch, not independently rerun here.
- Discovery failures retained: combined version command exit 1 at missing `t3` executable; scoped `rg` exit 2 for absent repo-local `.codex`/`.t3`; `codex app-server daemon version` exit 1 for missing default control socket. Help/version calls otherwise completed.

## Acceptance coverage and blockers

1. Actual phase/crew baseline: observed. Step 6 implementation procedure, core skill body and applicable crews were read. A file-change event creates `baseline-observation.md`, consistent with the existing descriptive naming rule.
2. Live compaction and subsequent reads/compliance: **Untested**. None of 65 exposed T3 tools describes compaction. CLI help exposes no control associated with this running worker. No verified boundary exists.
3. Runtime repeatability: T3 server **0.1.7**, provider **codex**, model **gpt-6.1-sol**, effort **high**, full-access/default modes observed through scoped APIs. Installed PATH CLI **0.167.0** observed. The exact backing Codex binary build and merged compaction settings remain unverified. Default daemon version query cannot connect to its socket.
4. Loss/reread/noncompliance distinction: no content loss established; no post-boundary classification possible. Repeated explicit reads and output truncation are not compaction evidence.
5. Actual Codex/T3 slice covered without replacement or hooks. Claude belongs to the separate integrator-owned slice and is not assessed here.

The experiment cannot support a retention success/failure conclusion or restoration changes. Baseline compliance has overlapping instruction sources, and this was not blinded. T3 activity visibility does not establish exhaustive provider-event visibility.

## Decisions

- D-01 — DECIDE & STATE: stop after bounded control discovery and preserve the blocker. A standalone CLI session would violate the requested actual-harness scope, while artificial context flooding lacks a known threshold or reliable event boundary.
- D-02 — DECIDE & STATE: preserve selected exact non-reasoning activity excerpts and source hashes rather than exporting a full conversation. This retains relevant read/compliance evidence while excluding unrelated private material.
- D-03 — DECIDE & STATE: keep P-A-U evidence here and leave the working REQ untouched under the explicit isolated-builder boundary. PLAN: scoped discovery and baseline; APPLY: six evidence artifacts; UNIFY: checks and full-range review above.

## Discovered Tasks

- impact-developer-visible: A repeatable actual T3 compaction run needs a supported active-provider compaction control plus an observable runtime boundary, and a scoped way to verify the backing provider build/settings. No runtime feature was implemented in this slice. → report only

## Lessons consulted

None required by the durable brief or working REQ, and none newly required by the artifact-only touch set. Implementation procedure and required crew files were read; maintainer release prime confirmed that `_dev/` changes do not require release metadata.

## Process cleanup

All command invocations completed. No delegated agents, provider sessions, servers, daemons, background jobs, or persistent command sessions were started. No global logs/configuration were inspected or changed, no hooks installed, no queue lifecycle changed, and no pushes/PRs made. The branch/worktree remain available for integration.
