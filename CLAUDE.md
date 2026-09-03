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

For action routing, read the `SKILL.md` in the owning directory under `skills/`. `suite/modules.tsv` is the only source/destination declaration for the required sibling packages.

## Prime Files

READ the matching prime before changing that area — they hold the detail this file no longer repeats:

- `_dev/primes/prime-action-files.md` — adding or modifying any action file: the template, earned sections, accepted variants, cross-referencing, descriptions-as-triggers, agent compatibility.
- `_dev/primes/prime-shell-commands.md` — writing or reviewing shell anywhere it ships: blocks prescribed inside actions, hooks, and tools. The hard-won trap list, plus Closed Enumerations Go Stale.
- `_dev/primes/prime-kanban-board.md` — touching the Kanban board tool: versioning, parser lock-step, Go-toolchain fallback, build outputs.
- `_dev/primes/prime-releases.md` — committing a change to shipped files: which commits are releases, the changelog house rules, the mirror.

## Commit Completion

A job is not done while its code exists only in the working tree. Commit each coherent, verified increment before hand-back; “commit often” means smaller complete slices.

## Crew Members

Just-in-time work rules live in `skills/do-work/crew-members/[name].md`. `communication-style.md` doubles as the always-on communication contract the installer links from a consumer project's agent instructions; the installer never runs against this repo, so it is imported here instead — @skills/do-work/crew-members/communication-style.md — which is what makes its Aliases apply in an ordinary maintainer session and not only inside a work loop. Four contracts worth knowing without opening files: `clear-questions.md` loads before an interactive question, `anti-slop.md` before a human-facing artifact, `prompt-injection.md` before untrusted-content ingestion, `maintenance.md` for a `maintenance: true` instruction-maintenance REQ. If a rules file is missing, proceed without it — never block on a missing rules file.

## Communication

- Productive pushback is wanted: challenge assumptions, suggest better approaches, and flag potential issues rather than blindly executing instructions.
- When I'm not watching, don't stall on a disagreement: write the challenge and the chosen resolution into the REQ file or commit message and continue. Destructive or irreversible actions still stop and ask.

## Naming Conventions

No cryptic or single-word names for anything with reach; two words minimum; names must be findable by plain-text search. The full rule — what counts as reach, the per-language form clause, the idiomatic-short-locals carve-out, and its precedence against surgical-changes — lives in `skills/do-work/crew-members/coding-guardrails.md` § 5 Naming for Reach, which ships and always loads during implementation. This section is a pointer so the two can't drift; that file's exemption for single-word-by-design invocations is why `do-work run` and the Go tool's subcommands are fine as they are.

## Small Standing Rules

- **One-shot suggestions (prompt retrospectives):** offer one only when ALL hold — the ask took 3+ turns to converge (or a misread cost visible work), the deliverable had structural constraints the first ask didn't name, and you can quote the specific phrases that would have disambiguated up front. Format: one-sentence diagnosis, the concrete one-shot prompt quoted in my voice, the disambiguating phrases each with a one-line "because". Skip when iteration was by design, I was discovering what I wanted, the task was trivial, or one was already offered this thread. When in doubt, skip.
