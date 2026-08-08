# Do-Work Suite Project

A task queue skill for agentic coding tools. Platform-agnostic — works with any agent that can read/write files and run shell commands.

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

## Before Every Commit

**Scope: the integrating commit only.** This ritual belongs to whoever commits the change into the integration branch — in the work pipeline, the queue owner at Step 9. A builder committing on its own `worktree-agent-*` branch **skips it entirely**: `skills/do-work/actions/version.md` and `CHANGELOG.md` are serial-only files owned by the integrator (`skills/do-work/actions/work-reference.md` → Worktree Dispatch Mode → Fan-Out Dispatch), and a builder bumping either would race every sibling.

1. **Bump the shared version** in `VERSION`, `skills/do-work/VERSION`, and `skills/do-work/actions/version.md` (line starting with `**Current version**:`). Use semver — patch for fixes, minor for features, major for breaking changes. When in doubt, patch. **Verify the new version number is strictly greater than the first existing entry in `CHANGELOG.md`** — duplicate version numbers have occurred before.

2. **Add a changelog entry** at the top of `CHANGELOG.md` (below the header). The title must **say what was delivered** — a reader scanning only headings should know what changed ("Board View Filters", not "The Fine Sieve"). No whimsical codenames. **Verify the title is not already used** by an earlier entry. (Historical entries were retroactively retitled to this convention in 0.117.1.)

```markdown
## X.Y.Z — [Short Descriptive Title] (YYYY-MM-DD)

[1-2 casual sentences — what changed and why it matters.]

- [Bullet points for specifics]
```

Keep it brief, newest on top, lead with value not implementation. Every version gets an entry.

## Action File Conventions

**Every NEW action must justify its package.** Before an action file is added, state — in its description blockquote or an accompanying ADR — why it belongs in core, board, knowledge, or toolbox: what package machinery it needs or which existing action it completes. Reviewers reject additions without this justification. Every new action also updates the owning `skills/do-work*/SKILL.md`, help, and any next-step surface; router budgets are enforced by `_dev/tests/contract-regressions.sh`.

Action files follow a consistent structure. When adding or modifying actions, use this template:

```markdown
# [Action Name] Action

> **Part of the do-work skill.** [1 sentence: what it does and when it's invoked.]

[Optional: read-only flag, philosophy, or key principles — 1-2 paragraphs max]

## When to Use

**Use when:** [2-4 bullets — positive triggers]
**Do NOT use when:** [2-3 bullets — explicit exclusions, with redirect to correct action]

## Input

[What parameters drive behavior: $ARGUMENTS, target REQ/UR, modes]

## Steps

### Step 1: [First action]

### Step 2: [...]

### Step N: [Final action]

## Output Format

[What gets produced — report structure, file changes, or user-facing output]

## Rules

[Include only if earned — see below. Constraints specific to this action, not restated engineering hygiene.]

## Common Rationalizations

[Include only if earned — see below.]

| If you're thinking...              | STOP. Instead...     | Because...               |
| ---------------------------------- | -------------------- | ------------------------ |
| [Shortcut the agent might attempt] | [What to do instead] | [Why the shortcut fails] |

## Red Flags

[Include only if earned — see below.]

- [Observable symptom that something went wrong — helps reviewers detect problems after the fact]

## Verification Checklist

[Include only if earned — see below.]

- [ ] [Concrete exit criterion with evidence requirement]
```

**Required:** Description blockquote, Steps (numbered). **Common:** Input, Output Format, When to Use.

**Earned, not mandatory: Rules, Common Rationalizations, Red Flags, Verification Checklist.** Add one only when the file has something a capable model would otherwise get wrong — do-work machinery (a queue/pipeline mechanic, a frontmatter or schema contract) or a hard-won failure mode with a traceable origin (a real REQ or incident this stops from recurring). "This is generic engineering advice a capable model already follows" is an explicit *non*-reason — true or not, it doesn't earn a section.

**The test, not a vibe:** before adding a Common Rationalizations row, ask *can I name the specific failure this row prevents, and where it happened?* No → don't add the row. If every row in a table fails that test, omit the whole section — a generic table is worse than no table: it teaches the reader the section is decorative, so they stop reading the ones that aren't. Apply the same test to Rules and Red Flags — specific to this action, not restated hygiene ("write tests," "don't skip validation"). When a file has nothing that passes, omit the section entirely; don't ship it empty or generic to satisfy the template.

**State intent, not a directive rule, when a capable model can infer the rest.** "Report drift, don't fix it inline" gives the model this action's boundary in one line — a five-line Rules section re-deriving why inline fixes are bad adds nothing a capable model didn't already know.

`_dev/tests/contract-regressions.sh` ratchets the Common Rationalizations rule: a new action file's table must contain at least one do-work-specific noun (illustrative, not exhaustive, per Closed Enumerations Go Stale below — e.g. REQ, UR, queue, frontmatter, pipeline, archive) or the suite fails, naming the file and the fix.

**Section order when present:** Philosophy → When to Use → Input → Steps → Output → Rules → Common Rationalizations → Red Flags → Verification Checklist.

**Accepted variants:**

- **Sub-command dispatchers** (`prime.md`, `bkb.md`) — Use a Sub-Commands table instead of flat steps. Each sub-command has its own workflow section.
- **Multi-mode actions** (`present-work.md`, `review-work.md`, `tutorial.md`) — Use a Modes table, then separate workflow sections per mode. A single `Step 1: Mode Selection` dispatcher at the top is acceptable.
- **State-based actions** (`version.md`, `pipeline.md`) — Response sections keyed by input type instead of sequential steps.
- **Checklist-based diagnostics** (`forensics.md`) — Use a `## Checks` section with independently-runnable items instead of ordered `## Steps`. Each check is a diagnostic probe, not a sequential step.

Cross-reference same-package actions by their local path (for example `actions/work.md`); cross a package boundary with an explicit sibling path such as `../do-work-knowledge/actions/bkb.md`. Shipped files must never cite this repo's own `CLAUDE.md` or `AGENTS.md` — both are export-ignored maintainer instructions. `_dev/tests/contract-regressions.sh` enforces this across the shipped `skills/` tree.

### Prescribed Shell Commands Must Surface What the Steps Consume

Action files are prose that prescribes shell behavior. When a step's logic iterates over the output of a command, the prescribed command must actually emit the items that logic consumes — a mismatch is invisible in the prose and only shows up when run against a real repo. Two traps that have already bitten this skill:

- **`git status --porcelain` collapses wholly-untracked directories** into a single `?? dir/` row — it does not list the files inside. Any step that enumerates untracked files per-item (read each, check extension/size/name) must use `git status --porcelain --untracked-files=all` (`-uall`) or `git ls-files --others --exclude-standard`. The latter also drops correctly-ignored paths, so it doubles as the untracked ignore filter.
- **A blanket skip/exclude list applied _before_ a check silently neuters any check meant to fire inside the excluded set.** Scope skip-lists to the noise they actually target (untracked/ignored) and run tracked-file checks outside the exclusion — e.g. a committed `__pycache__/*.pyc` is correct-to-ignore when untracked but is exactly what a "committed build artifact" check should flag.
- **`git show --name-only` prints the commit header and message before the file list** — a message line can pass a filename grep and become a phantom path, and merge commits list no files at all. Use `git diff-tree --no-commit-id --name-only -r -m <commit>` (or `git show --name-only --format=`) when the output is consumed as file paths.
- **`git show` on a merge commit prints a combined diff that is usually empty** — so any consumer reading a REQ's `commit:` hash as a diff source silently sees nothing on worktree-merged work (the reviewer reads an empty diff as an empty REQ). Detect the second parent with `git rev-parse --verify -q '<sha>^2'` — quoted, since `^` is a glob operator in zsh and an escape character in cmd.exe — and use `git show --first-parent -m <sha>` when it succeeds.
- **Ignore patterns with an interior slash are root-anchored, while `git check-ignore` tests cwd-relative paths** — a guard that checks then appends can mismatch from a subdirectory (duplicate appends, path never ignored). Prefix with `**/` when the consumer may run below the repo root. Relatedly, never build `.git/`-internal paths from `--show-toplevel`; use `git rev-parse --git-path <name>` (worktree- and submodule-safe).
- **Never interpolate raw user text inside shell quoting.** A prescribed command like `$(echo '<user-slug>' | tr ...)` breaks on an apostrophe and is a command-injection vector. Derive a sanitized token as a text operation first, then substitute the already-safe value.
- **`diff -x PATTERN` matches basenames of files _and directories_.** Excluding a build artifact by bare name (`-x queue-kanban`) also excludes any same-named directory — silently blinding the diff to an entire source tree. Filter the diff's *output* for the specific artifact path instead (`| grep -v 'tools/queue-kanban/queue-kanban'`), or use a pattern that can only match the file.
- **`curl -o` writes the final path incrementally** — a mid-transfer failure leaves a non-empty partial file, so any presence- or size-gated consumer (`test -s` detect checks) reads the broken download as complete. Prescribe download-to-a-temp-name plus rename-on-success (`curl -o x.download … && mv x.download x || { rm -f x.download; false; }`); `--remove-on-error` needs curl ≥ 7.83, the rename works everywhere. The `; false` is not optional — `rm -f` on an absent path exits 0, so the plain `|| rm -f` form cleans up and then reports the failed download as a success.
- **Shell state does not survive between prescribed command blocks.** An action's steps run as separate shell invocations (often with a user-confirmation gate between them); a variable defined in one block — especially a `mktemp` random path — expands empty in the next, and an agent that "recovers" by re-running the earlier download can silently bypass a review the flow depends on. Blocks must re-derive what they need from deterministic paths and guard-check that inherited artifacts actually exist.

When a review finds a bug in prescribed-command logic, **grep the same primitive across all actions before calling it fixed** — these patterns are usually copy-pasted, so the fix is rarely local. (The first trap above had been copy-pasted into four action files; the audit only flagged one of them.)

### Closed Enumerations Go Stale

When a rule applies "whenever X happens" (load a guardrail, honor an enum, keep a guide in sync), state the trigger _condition_ in the rule's canonical home and mark any caller/value list as illustrative, not exhaustive. Hand-enumerated lists silently go stale the moment the set grows — one review traced four independent defects to this pattern (capture's stale domain enum, prompt-injection's five-caller list, the docs-exemption list, security.md's loader claims). When extending a set, grep for every other enumeration of it and update or generalize each one.

## Agent Rules

Just-in-time work rules live in `skills/do-work/crew-members/[name].md`. Each file's `JIT_CONTEXT` comment is the canonical statement of when it loads. The loading order is specified in `skills/do-work/actions/work.md` Step 6.

- `general.md` and `coding-guardrails.md` are always loaded during implementation. Everything else loads conditionally per its `JIT_CONTEXT` (e.g., domain match, `tdd`/`caveman` flags, security surface, fan-out, human-facing artifact production, third-party content ingestion, skill-instruction maintenance passes, debugging retries, interviews).
- Four contracts worth knowing without opening files: `clear-questions.md` loads before an interactive question; `anti-slop.md` loads before a human-facing artifact; `prompt-injection.md` loads before untrusted-content ingestion; `maintenance.md` loads for a `maintenance: true` instruction-maintenance REQ. Each package carries the crew files its actions need.
- If a rules file is missing, proceed without it — never block on a missing rules file.

## Queue Path Convention

Pending REQ files live in `do-work/queue/` — not `do-work/` root.

## Shipped Board Tooling (`skills/do-work-board/tools/`)

`skills/do-work-board/tools/queue-kanban/` is a standalone Go module (its own `go.mod`, embedded `web/` frontend) that renders the `do-work/` queue as a Kanban board. It ships in the board module and is invoked by `do-work-board board` (`skills/do-work-board/actions/board.md`). Conventions:

- **Versioning is folded into the skill.** The tool has no independent changelog — its changes get entries in the root `CHANGELOG.md` and a normal skill version bump, exactly like any action. (It was independently versioned through 1.1.0 before being vendored in; that history lives in `decisions/records/adr-016-*`.)
- **Keep the parser in lock-step with the schema.** The board buckets tickets by the `status` vocabulary in `skills/do-work/actions/work-reference.md`; its parsed display fields must stay aligned with `skills/do-work-board/tools/queue-kanban/model.go`, and the Testing placeholders with `skills/do-work-board/tools/queue-kanban/testing.go`.
- **The tool has exactly three write surfaces, and none touches pipeline state.** (1) The board's Testing view writes only the testing placeholders above plus `do-work/testers.md`. (2) `queue-kanban next-version` rewrites the single `**Current version**: X.Y.Z` line in one version file (`actions/version.md` by default, `--version-file` to point elsewhere). (3) `queue-kanban next-req` atomically creates one durable number marker under `do-work/.req-reservations/`; this is queue coordination metadata, not a REQ or pipeline-field edit. Everything else the tool does is read-only — illustratively `summary`, `open-work`, `generate`, `serve`'s board views, `frontmatter`, `verify`, and `now` (which reads a clock, not even the tree); the rule is the count, not this list. Nothing in the tool writes `status`, any other pipeline field, or **`CHANGELOG.md`** — the changelog stays an owner-only, human-authored write. Adding a fourth write surface means amending this sentence in the same commit; that is the co-location rule applied to itself.
- **`verify` is the mechanical half of "Before Every Commit."** It is wired into `skills/do-work/actions/forensics.md`; it reports and routes, while repairs belong to `skills/do-work/actions/cleanup.md`.
- **Toolchain exception to "design for the floor."** The board is the one capability that needs Go (`skills/do-work-board/tools/queue-kanban/go.mod`); `skills/do-work-board/actions/board.md` degrades gracefully when it is absent. Core may use an already-built sibling binary only where a shell-portable fallback remains documented.
- **Never commit build outputs.** The compiled `queue-kanban` binary is gitignored by `skills/do-work-board/tools/queue-kanban/.gitignore`; the `do-work-board static` artifact lands in `build/` at the repo root.

## Lessons → Knowledge Base Handoff

After review passes, core offers to promote `## Lessons Learned` through `skills/do-work/actions/kb-lessons-handoff.md`; later BKB processing belongs to `do-work-knowledge`.

## Agent Compatibility

Action files must work with **any** agentic coding tool:

- Use generalized language ("spawn a subagent", "use your environment's ask-user prompt") — no tool-specific APIs in action files.
- Each action file should work as a standalone prompt pasted into a basic chat interface.
- Design for the floor: the simplest agent that can read/write files and run shell commands must be able to follow the instructions. Subagents and parallel execution are nice-to-haves.

## One-Shot Suggestions (Prompt Retrospectives)

When ALL of these hold — the ask took 3+ turns to converge (or a misread cost visible work), the final deliverable has structural constraints the first ask didn't name (format, destination, stack, audience, scope), and you can point to specific phrases that would have disambiguated up front — close your reply with a short retrospective:

1. One-sentence diagnosis of the core ambiguity.
2. The concrete one-shot prompt the user could have sent, quoted, in their voice — not a template.
3. The specific disambiguating phrases, each with a one-line "because...".
4. Optionally a one-sentence meta-lesson.

Skip it when the iteration was by design (`scan-ideas`, `deep-explore`, review loops), the user was discovering what they wanted mid-conversation, the task was trivial, or you've already offered one this thread. It's feedback on phrasing, not self-flagellation — and when in doubt, skip it.

## Communication Style

- The user appreciates productive pushback — challenge assumptions, suggest better approaches, and flag potential issues rather than blindly executing instructions

## Naming Conventions

No cryptic or single-word names for anything with reach; two words minimum; names must be
findable by plain-text search. The full rule — what counts as reach, the per-language form
clause, the idiomatic-short-locals carve-out, and its precedence against surgical-changes —
now lives in `skills/do-work/crew-members/coding-guardrails.md` § 5 Naming for Reach, which ships and
always loads during implementation. It applies here like it does in any consumer project;
this section is a pointer so the two can't drift. That file's exemption for
single-word-by-design invocations is why `do-work run` and the Go tool's subcommands are
fine as they are.
