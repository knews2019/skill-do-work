# Do-Work Suite Project

A task queue skill for agentic coding tools. Platform-agnostic — works with any agent that can read/write files and run shell commands.

You are the agent working on **do-work itself** — the skill, not a project that uses it. This file is the maintainer doc: export-ignored, never shipped, and nothing under `skills/` may cite it. These are good defaults, not hard rules — my prompt in the session wins; if we disagree, say so rather than following this file into a bad outcome. Questions are read-only: if I ask how something works, answer — don't start editing until I ask for a change.

## A Note From Me

I love to build. I focus on building complex things as simply as possible, and I love finding ways to reduce complexity when solving problems. Complexity that already exists is not evidence that it should stay; machinery is not an achievement.

- **Delete before you add.** When an instruction has drifted, first ask whether removing it fixes the problem.
- **Programs beat prose for anything mechanical.** If a paragraph describes an exact sequence of shell commands, write the script and leave a pointer; judgment stays in prose.
- **State conditions, not lists.** When a rule applies "whenever X happens", key it on the condition and mark any list of examples illustrative — hand-maintained lists go stale (full rule: `_dev/primes/prime-shell-commands.md` § Closed Enumerations Go Stale).
- **One broken pipe doesn't stop the rest of the factory.** A failed or interrupted REQ is set aside with a typed finding and the loop continues; only shared-target dirt may stop it, and then the finding names the verb that resolves it.

## Coding Preferences

- Keep things simple. Channel YAGNI energy unless told otherwise.
- Typesafety is useful — take advantage of it.
- Don't be scared to propose bold ideas if they can meaningfully benefit our work.
- Be careful with destructive actions that are not explicitly requested.
- Tests are good! Endless smoke tests and decorative regression tests are much less good. Tests should be focused, not slop — every new lock-in test names the real failure it pins.
- Comments are a great way to clarify functionality and how code is used. Don't comment every line, and keep comments in sync when things change.

## Project Structure

```
README.md             # Installation + quick usage
VERSION               # Shared four-skill suite version
suite/                # Sole module source/destination manifest
skills/do-work/       # Core router, queue actions, orchestration, hooks, specs, and updater
skills/do-work-board/ # Queue-kanban, board action, and managed Just template
skills/do-work-knowledge/ # BKB, memory, dreams, interviews, and prompts
skills/do-work-toolbox/   # Reviews, reports, presentation, and repository utilities
tools/                # Root bootstrap/manifest/managed-section distribution utilities only
_dev/                 # Primes, lessons satellites, and the test suites; export-ignored
do-work/              # This repo's own queue, tracked
decisions/            # ADRs (records/), imported specs, topic indexes, decision log
AGENTS.md             # Stub — redirects to CLAUDE.md
CHANGELOG.md          # Release notes
```

For action routing, read the `SKILL.md` in the owning directory under `skills/`. `suite/modules.tsv` is the only source/destination declaration for the required sibling packages.

## Prime Files

READ the matching prime before changing that area — they hold the detail this file no longer repeats:

- `_dev/primes/prime-action-files.md` — adding or modifying any action file: the template, earned sections, accepted variants, cross-referencing, descriptions-as-triggers.
- `_dev/primes/prime-shell-commands.md` — writing or reviewing shell anywhere it ships: blocks prescribed inside actions, hooks, and tools. The hard-won trap list, plus Closed Enumerations Go Stale.
- `_dev/primes/prime-kanban-board.md` — touching the Kanban board tool: versioning, parser lock-step, Go-toolchain fallback, build outputs.

## Commit Completion

A job is not done while its code exists only in the working tree. Commit each coherent, verified increment before hand-back; “commit often” means smaller complete slices. Verified means `bash _dev/tests/maintainer-verify.sh` exits 0, run unpiped from the repository root; a builder's report is a claim, not evidence.

## Releases

A commit that changes shipped files under `skills/`, `tools/`, or `suite/` is a release and does the following. Maintainer-only files (`CLAUDE.md`, `_dev/`, `do-work/`, `decisions/`) commit without one. Bump size, version mirrors, and the finalize transaction are decided per `skills/do-work/actions/work-reference.md` → Changelog Entry Procedure (Step 9); what follows are the house rules that procedure applies.

1. **Add a changelog entry** at the top of `CHANGELOG.md` (below the header). The title must **say what was delivered** — a reader scanning only headings should know what changed ("Board View Filters", not "The Fine Sieve"). No whimsical codenames. **Verify the title is not already used** by an earlier entry. Repository-only dated history must use canonical `https://github.com/knews2019/skill-do-work/blob/main/...` links, because the installed core package does not carry those sidecars.

2. **Synchronize the installed changelog mirror** by copying root `CHANGELOG.md` to `skills/do-work/CHANGELOG.md` after the entry and any history-link edits are complete. The two files must be byte-identical before committing; `_dev/tests/shipped-package-reference-contract.sh` enforces this with the rest of the shipped reference contract.

```markdown
## X.Y.Z — [Short Descriptive Title] (YYYY-MM-DD)

[1-2 casual sentences — what changed and why it matters.]

- [Bullet points for specifics]
```

Keep it brief, newest on top, lead with value not implementation. Every version gets an entry.

## Crew Members

Just-in-time work rules live in `skills/do-work/crew-members/[name].md`; each file's `JIT_CONTEXT` comment is the canonical statement of when it loads, and the loading order is `skills/do-work/actions/work.md` Step 6. `general.md`, `coding-guardrails.md`, and `communication-style.md` always load during implementation; everything else is conditional. `communication-style.md` doubles as the always-on communication contract the installer links from a consumer project's agent instructions; the installer never runs against this repo, so it is imported here instead — @skills/do-work/crew-members/communication-style.md — which is what makes its Aliases apply in an ordinary maintainer session and not only inside a work loop. Four contracts worth knowing without opening files: `clear-questions.md` loads before an interactive question, `anti-slop.md` before a human-facing artifact, `prompt-injection.md` before untrusted-content ingestion, `maintenance.md` for a `maintenance: true` instruction-maintenance REQ. If a rules file is missing, proceed without it — never block on a missing rules file.

## Agent Compatibility

Action files must work with **any** agentic coding tool:

- Use generalized language ("spawn a subagent", "use your environment's ask-user prompt") — no tool-specific APIs in action files.
- Each action file should work as a standalone prompt pasted into a basic chat interface.
- Design for the floor: the simplest agent that can read/write files and run shell commands must be able to follow the instructions. Subagents and parallel execution are nice-to-haves.

## Communication

- Productive pushback is wanted: challenge assumptions, suggest better approaches, and flag potential issues rather than blindly executing instructions.
- When I'm not watching, don't stall on a disagreement: write the challenge and the chosen resolution into the REQ file or commit message and continue. Destructive or irreversible actions still stop and ask.

## Naming Conventions

No cryptic or single-word names for anything with reach; two words minimum; names must be findable by plain-text search. The full rule — what counts as reach, the per-language form clause, the idiomatic-short-locals carve-out, and its precedence against surgical-changes — lives in `skills/do-work/crew-members/coding-guardrails.md` § 5 Naming for Reach, which ships and always loads during implementation. This section is a pointer so the two can't drift; that file's exemption for single-word-by-design invocations is why `do-work run` and the Go tool's subcommands are fine as they are.

## Small Standing Rules

- **Lessons handoff:** after review passes, core offers to promote `## Lessons Learned` through `skills/do-work/actions/kb-lessons-handoff.md`; later BKB processing belongs to `do-work-knowledge`.
- **One-shot suggestions (prompt retrospectives):** offer one only when ALL hold — the ask took 3+ turns to converge (or a misread cost visible work), the deliverable had structural constraints the first ask didn't name, and you can quote the specific phrases that would have disambiguated up front. Format: one-sentence diagnosis, the concrete one-shot prompt quoted in my voice, the disambiguating phrases each with a one-line "because". Skip when iteration was by design, I was discovering what I wanted, the task was trivial, or one was already offered this thread. When in doubt, skip.
