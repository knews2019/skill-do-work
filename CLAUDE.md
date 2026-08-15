# Do-Work Suite Project

A task queue skill for agentic coding tools. Platform-agnostic — works with any agent that can read/write files and run shell commands.

You are the agent working on **do-work itself** — the skill, not a project that uses it. This file is the maintainer doc: export-ignored, never shipped, and nothing under `skills/` may cite it. These are good defaults, not hard rules — my prompt in the session wins; if we disagree, say so rather than following this file into a bad outcome. Questions are read-only: if I ask how something works, answer — don't start editing until I ask for a change.

## A Note From Me

I love to build. I focus on building complex things as simply as possible, and I love finding ways to reduce complexity when solving problems. Complexity that already exists is not evidence that it should stay; machinery is not an achievement.

- **Delete before you add.** When an instruction has drifted, first ask whether removing it fixes the problem.
- **Programs beat prose for anything mechanical.** If a paragraph describes an exact sequence of shell commands, write the script and leave a pointer; judgment stays in prose.
- **State conditions, not lists.** When a rule applies "whenever X happens", key it on the condition and mark any list of examples illustrative — hand-maintained lists go stale (full rule: `_dev/primes/prime-shell-commands.md` § Closed Enumerations Go Stale).

## Coding Preferences

- Keep things simple. Channel YAGNI energy unless told otherwise.
- Typesafety is useful — take advantage of it.
- Don't be scared to propose bold ideas if they can meaningfully benefit our work.
- Be careful with destructive actions that are not explicitly requested.
- Tests are good! Endless smoke tests and decorative regression tests are much less good. Tests should be focused, not slop — every new lock-in test names the real failure it pins.
- Comments are a great way to clarify functionality and how code is used. Don't comment every line, and keep comments in sync when things change.

## Match Ceremony to the Task

- Don't spawn subagents or a multi-agent panel for work a single agent finishes in one pass. Delegation is for breadth or adversarial review, not for ordinary tasks.
- When several agents do work in parallel, state file ownership up front so they don't collide — that is exactly what a REQ's write set is for.

## Glossary

Plain words for this repo's terms — describe things back to me in this vocabulary.

| Term | Meaning |
|---|---|
| I / me | the maintainer |
| user | someone who installed do-work into their own repo — never me |
| UR | user request — the durable statement of intent |
| REQ | one implementable request under a UR; the unit of work through the pipeline |
| action | a `do-work <verb>` entry point; one file under `skills/do-work*/actions/`, routed by SKILL.md |
| crew member | a rules file loaded just-in-time during a build (`skills/do-work/crew-members/`) |
| prime file | a lazy-loaded context doc read only when working in its domain (`_dev/primes/` here) |
| the queue | pending REQ files in `do-work/queue/` — not `do-work/` root |
| the Kanban board | `skills/do-work-board/tools/queue-kanban` — the human view of the queue |
| lock-in test | a test that keeps a fixed bug fixed; the `_dev/tests/` suites are full of them |
| builder's report | what a worktree builder hands back — a claim, not evidence; judge from git state |
| write set | the files a REQ may touch — the collision guard for parallel work; display-only on the board, never a safety guarantee |

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

## Kanban Board Write Surfaces

The Kanban board tool has exactly three write surfaces, and none touches pipeline state: (1) the board's Testing view writes only the testing placeholders plus `do-work/testers.md`; (2) `queue-kanban next-version` rewrites the single `**Current version**: X.Y.Z` line in one version file (`actions/version.md` by default, `--version-file` to point elsewhere); (3) `queue-kanban next-req` atomically creates one durable number marker under `do-work/.req-reservations/` — queue coordination metadata, not a REQ or pipeline-field edit. Everything else the tool does is read-only; the rule is the count, not any list of subcommands. Nothing in the tool writes `status`, any other pipeline field, or **`CHANGELOG.md`** — the changelog stays an owner-only, human-authored write. Adding a fourth write surface means amending this sentence in the same commit; that is the co-location rule applied to itself. Everything else about the board lives in `_dev/primes/prime-kanban-board.md`.

## Before Every Commit

**Scope: the integrating commit only.** This ritual belongs to whoever commits the change into the integration branch — in the work pipeline, the queue owner at Step 9. A builder committing on its own `worktree-agent-*` branch **skips it entirely**: `skills/do-work/actions/version.md` and `CHANGELOG.md` are serial-only files owned by the integrator (`skills/do-work/actions/work-reference.md` → Worktree Dispatch Mode → Fan-Out Dispatch), and a builder bumping either would race every sibling.

1. **Bump the shared version** in `VERSION`, `skills/do-work/VERSION`, and `skills/do-work/actions/version.md` (line starting with `**Current version**:`). Use semver — patch for fixes, minor for features, major for breaking changes. When in doubt, patch. **Verify the new version number is strictly greater than the first existing entry in `CHANGELOG.md`** — duplicate version numbers have occurred before.

2. **Add a changelog entry** at the top of `CHANGELOG.md` (below the header). The title must **say what was delivered** — a reader scanning only headings should know what changed ("Board View Filters", not "The Fine Sieve"). No whimsical codenames. **Verify the title is not already used** by an earlier entry. Repository-only dated history must use canonical `https://github.com/knews2019/skill-do-work/blob/main/...` links, because the installed core package does not carry those sidecars.

3. **Synchronize the installed changelog mirror** by copying root `CHANGELOG.md` to `skills/do-work/CHANGELOG.md` after the entry and any history-link edits are complete. The two files must be byte-identical before committing; `_dev/tests/shipped-package-reference-contract.sh` enforces this with the rest of the shipped reference contract.

```markdown
## X.Y.Z — [Short Descriptive Title] (YYYY-MM-DD)

[1-2 casual sentences — what changed and why it matters.]

- [Bullet points for specifics]
```

Keep it brief, newest on top, lead with value not implementation. Every version gets an entry.

## Verify

`bash _dev/tests/contract-regressions.sh` is the baseline pass/fail check before any hand-back; run the other `_dev/tests/*.sh` suites when you touched what they cover. Exit code zero is the only proof — never accept a summary or a builder's report as evidence that a check passed. Never pipe a check through `| tail` or similar: the pipeline's exit status hides the failure.

## Crew Members

Just-in-time work rules live in `skills/do-work/crew-members/[name].md`; each file's `JIT_CONTEXT` comment is the canonical statement of when it loads, and the loading order is `skills/do-work/actions/work.md` Step 6. `general.md` and `coding-guardrails.md` always load during implementation; everything else is conditional. Four contracts worth knowing without opening files: `clear-questions.md` loads before an interactive question, `anti-slop.md` before a human-facing artifact, `prompt-injection.md` before untrusted-content ingestion, `maintenance.md` for a `maintenance: true` instruction-maintenance REQ. If a rules file is missing, proceed without it — never block on a missing rules file.

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
