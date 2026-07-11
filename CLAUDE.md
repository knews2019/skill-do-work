# Do-Work Skill Project

A task queue skill for agentic coding tools. Platform-agnostic — works with any agent that can read/write files and run shell commands.

## Project Structure

```
SKILL.md              # Entry point — routing logic, action dispatch; authoritative action-name → file-path index
next-steps.md         # Per-action next-step suggestions (referenced by SKILL.md)
README.md             # Installation + quick usage guide
actions/              # Action files — each is a standalone prompt; heavy actions ship a *-reference.md companion
specs/                # Reusable specification templates for common task types (see specs/README.md)
prompts/              # Reusable prompt library (see prompts/README.md for the authoritative index)
interviews/           # Prescriptive templates loaded by the interview action
crew-members/         # Agent rules loaded just-in-time — each file's JIT_CONTEXT comment states when it loads
hooks/                # Optional hook scripts (platform-specific, installable; hooks.json + shell scripts)
tools/                # Shipped compiled tooling — queue-kanban/ (standalone Go module) renders the do-work queue as a Kanban board; built on demand via `do-work board`
docs/                 # User guides for the most commonly used actions — not every action has one
decisions/            # Architecture decisions — ADRs (records/), imported specs, topic indexes, decision log
AGENTS.md             # Stub — redirects to CLAUDE.md
CHANGELOG.md          # Release notes (newest on top)
```

For the per-action file list with descriptions, read `SKILL.md` — it is the canonical name→path mapping. This tree deliberately stops at directories so it cannot drift from the repo.

## Before Every Commit

1. **Bump the version** in `actions/version.md` (line starting with `**Current version**:`). Use semver — patch for fixes, minor for features, major for breaking changes. When in doubt, patch. **Verify the new version number is strictly greater than the first existing entry in `CHANGELOG.md`** — duplicate version numbers have occurred before.

2. **Add a changelog entry** at the top of `CHANGELOG.md` (below the header). The title must **say what was delivered** — a reader scanning only headings should know what changed ("Board View Filters", not "The Fine Sieve"). No whimsical codenames. **Verify the title is not already used** by an earlier entry. (Entries before 0.117.0 used fun codenames; leave them as-is.)

```markdown
## X.Y.Z — [Short Descriptive Title] (YYYY-MM-DD)

[1-2 casual sentences — what changed and why it matters.]

- [Bullet points for specifics]
```

Keep it brief, newest on top, lead with value not implementation. Every version gets an entry.

## Action File Conventions

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

[Constraints, common mistakes, what NOT to do]

## Common Rationalizations

| If you're thinking...              | STOP. Instead...     | Because...               |
| ---------------------------------- | -------------------- | ------------------------ |
| [Shortcut the agent might attempt] | [What to do instead] | [Why the shortcut fails] |

## Red Flags

- [Observable symptom that something went wrong — helps reviewers detect problems after the fact]

## Verification Checklist

- [ ] [Concrete exit criterion with evidence requirement]
```

**Required elements:** Description blockquote, Steps (numbered). **Common elements:** Input, Output Format, Rules, When to Use. **Encouraged elements:** Common Rationalizations, Red Flags, Verification Checklist. **Section order matters:** always Philosophy → When to Use → Input → Steps → Output → Rules → Common Rationalizations → Red Flags → Verification Checklist.

**Accepted variants:**

- **Sub-command dispatchers** (`prime.md`, `bkb.md`) — Use a Sub-Commands table instead of flat steps. Each sub-command has its own workflow section.
- **Multi-mode actions** (`present-work.md`, `review-work.md`, `tutorial.md`) — Use a Modes table, then separate workflow sections per mode. A single `Step 1: Mode Selection` dispatcher at the top is acceptable.
- **State-based actions** (`version.md`, `pipeline.md`) — Response sections keyed by input type instead of sequential steps.
- **Checklist-based diagnostics** (`forensics.md`) — Use a `## Checks` section with independently-runnable items instead of ordered `## Steps`. Each check is a diagnostic probe, not a sequential step.

Cross-reference other actions by their **file path** (e.g., `actions/work.md`, or `actions/work-reference.md`'s Schema Read Contract) so an agent reading the file can open the target directly without resolving a name to a path. Companion reference files take a path too (`actions/interview-reference.md`, `actions/bkb-reference.md`). The one exception is a `do-work <verb>` **command invocation** (`do-work run`, `do-work clarify`) — that's how an action is _run_, not a pointer to its file, so keep it as a command. SKILL.md remains the authoritative name→path mapping and may use short names in its routing prose.

### Prescribed Shell Commands Must Surface What the Steps Consume

Action files are prose that prescribes shell behavior. When a step's logic iterates over the output of a command, the prescribed command must actually emit the items that logic consumes — a mismatch is invisible in the prose and only shows up when run against a real repo. Two traps that have already bitten this skill:

- **`git status --porcelain` collapses wholly-untracked directories** into a single `?? dir/` row — it does not list the files inside. Any step that enumerates untracked files per-item (read each, check extension/size/name) must use `git status --porcelain --untracked-files=all` (`-uall`) or `git ls-files --others --exclude-standard`. The latter also drops correctly-ignored paths, so it doubles as the untracked ignore filter.
- **A blanket skip/exclude list applied _before_ a check silently neuters any check meant to fire inside the excluded set.** Scope skip-lists to the noise they actually target (untracked/ignored) and run tracked-file checks outside the exclusion — e.g. a committed `__pycache__/*.pyc` is correct-to-ignore when untracked but is exactly what a "committed build artifact" check should flag.
- **`git show --name-only` prints the commit header and message before the file list** — a message line can pass a filename grep and become a phantom path, and merge commits list no files at all. Use `git diff-tree --no-commit-id --name-only -r -m <commit>` (or `git show --name-only --format=`) when the output is consumed as file paths.
- **Ignore patterns with an interior slash are root-anchored, while `git check-ignore` tests cwd-relative paths** — a guard that checks then appends can mismatch from a subdirectory (duplicate appends, path never ignored). Prefix with `**/` when the consumer may run below the repo root. Relatedly, never build `.git/`-internal paths from `--show-toplevel`; use `git rev-parse --git-path <name>` (worktree- and submodule-safe).
- **Never interpolate raw user text inside shell quoting.** A prescribed command like `$(echo '<user-slug>' | tr ...)` breaks on an apostrophe and is a command-injection vector. Derive a sanitized token as a text operation first, then substitute the already-safe value.
- **`diff -x PATTERN` matches basenames of files _and directories_.** Excluding a build artifact by bare name (`-x queue-kanban`) also excludes any same-named directory — silently blinding the diff to an entire source tree. Filter the diff's *output* for the specific artifact path instead (`| grep -v 'tools/queue-kanban/queue-kanban'`), or use a pattern that can only match the file.
- **Shell state does not survive between prescribed command blocks.** An action's steps run as separate shell invocations (often with a user-confirmation gate between them); a variable defined in one block — especially a `mktemp` random path — expands empty in the next, and an agent that "recovers" by re-running the earlier download can silently bypass a review the flow depends on. Blocks must re-derive what they need from deterministic paths and guard-check that inherited artifacts actually exist.

When a review finds a bug in prescribed-command logic, **grep the same primitive across all actions before calling it fixed** — these patterns are usually copy-pasted, so the fix is rarely local. (The first trap above had been copy-pasted into four action files; the audit only flagged one of them.)

### Closed Enumerations Go Stale

When a rule applies "whenever X happens" (load a guardrail, honor an enum, keep a guide in sync), state the trigger _condition_ in the rule's canonical home and mark any caller/value list as illustrative, not exhaustive. Hand-enumerated lists silently go stale the moment the set grows — one review traced four independent defects to this pattern (capture's stale domain enum, prompt-injection's five-caller list, the docs-exemption list, security.md's loader claims). When extending a set, grep for every other enumeration of it and update or generalize each one.

## Agent Rules

Just-in-time rules live in `crew-members/[name].md`. Each file's `JIT_CONTEXT` comment is the canonical statement of when it loads — that comment is the contract, not any list duplicated here or elsewhere. The loading order for the work pipeline is specified in `actions/work.md` Step 6.

- `general.md` and `karpathy.md` are always loaded during implementation. Everything else loads conditionally per its `JIT_CONTEXT` (e.g., domain match, `tdd`/`caveman` flags, security surface, fan-out, human-facing artifact production, third-party content ingestion, skill-instruction maintenance passes, debugging retries, interviews).
- Four contracts worth knowing without opening files: `clear-questions.md` loads before presenting the user any **interactive question** (ask-tool prompt, clarifying question, option menu) and governs question wording — one decision per question, no unglossed shorthand, options that state their consequence; `anti-slop.md` loads before producing any **human-facing artifact**; `prompt-injection.md` loads before ingesting any **content not authored by the current invocation or the shipped skill files**; `maintenance.md` loads before a **deliberate maintenance pass on the skill's own instructions** (fixing a drifting agent/action/crew/prime file, where removing or narrowing is a candidate fix) and codifies delete-before-you-add — the maintenance-time complement to `karpathy.md`'s implementation-time surgical-changes rule. In the work pipeline the trigger is the `maintenance: true` REQ marker (set by capture for a removal/narrowing finding on the skill's own instructions; loaded by `actions/work.md` Step 6) — marker-only, never a description heuristic, which would misfire on ordinary implementation REQs. New actions that hit any of these triggers must load the corresponding file.
- If a rules file is missing, proceed without it — never block on a missing rules file.

## Queue Path Convention

Pending REQ files live in `do-work/queue/`. When referencing the queue in action files, always use `do-work/queue/` — not `do-work/` root.

## Shipped Tooling (`tools/`)

`tools/queue-kanban/` is a standalone Go module (its own `go.mod`, embedded `web/` frontend) that renders the `do-work/` queue as a Kanban board. It ships in the tarball (it is **not** `export-ignore`'d) so `do-work update` carries it into every consumer; the `do-work board` action (`actions/board.md`) builds and runs it. Conventions:

- **Versioning is folded into the skill.** The tool has no independent changelog — its changes get entries in the root `CHANGELOG.md` and a normal skill version bump, exactly like any action. (It was independently versioned through 1.1.0 before being vendored in; that history lives in `decisions/records/adr-016-*`.)
- **Keep the parser in lock-step with the schema.** The board buckets tickets by the `status` vocabulary defined in `actions/work-reference.md`'s Schema Read Contract; `depends_on` and `domain` are parsed for display only (badges, drawer metadata — no column logic). Any change to that contract must be mirrored in `tools/queue-kanban/model.go` (and vice-versa) in the same commit — co-location is the whole point.
- **Toolchain exception to "design for the floor."** The board is the one capability that needs a compiler (Go, per `tools/queue-kanban/go.mod`). `actions/board.md` precondition-checks `go` and degrades gracefully when it's absent — it never blocks the rest of the skill. Don't reach for a compiled tool in any other action.
- **Never commit build outputs.** The compiled `queue-kanban` binary is gitignored by `tools/queue-kanban/.gitignore` (which ships, keeping it ignored downstream); the `do-work board static` artifact lands in `build/` at the repo root.

## Lessons → Knowledge Base Handoff

After a REQ's review passes, review-work (standalone mode) and work (pipeline mode) both offer to promote `## Lessons Learned` into the project's KB via `actions/kb-lessons-handoff.md` — see that file for the full contract (payload shape, consent flow, the optional `kb_status`/`kb_entry` REQ frontmatter fields). The handoff is pure do-work, never blocks archival, and defers to `pending` when no `kb/` exists.

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

- **No cryptic or single-word variable names.** Every variable and function name should be at least two words
  (e.g., `invoice_total`, `retry_count`, `alignment_score`) so its purpose is immediately obvious.
- **Optimize for grepability.** Names should be unique enough across the codebase that a simple text search
  (ripgrep, fd, sad) locates every usage — no IDE or LSP required to trace where a name has effect.
- **Favor clarity over brevity.** `pending_invoice_items` beats `pii`. `max_retry_attempts` beats `mra`.
  If a name needs a comment to explain it, the name isn't good enough.
