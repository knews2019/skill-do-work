# Do-Work Skill Project

A task queue skill for agentic coding tools. Platform-agnostic — works with any agent that can read/write files and run shell commands.

## Project Structure

```
SKILL.md              # Entry point — routing logic, action dispatch; authoritative action-name → file-path index
next-steps.md         # Per-action next-step suggestions
README.md             # Installation + quick usage
actions/              # Action files — each a standalone prompt; heavy actions ship a *-reference.md companion
specs/                # Specification templates (see specs/README.md)
prompts/              # Prompt library (see prompts/README.md for the index)
interviews/           # Prescriptive templates loaded by the interview action
crew-members/         # Agent rules loaded just-in-time — each file's JIT_CONTEXT comment states when it loads
hooks/                # Optional hook scripts (platform-specific; hooks.json + shell scripts)
tools/                # Shipped compiled tooling — queue-kanban/ renders the do-work queue as a Kanban board; built via `do-work board`
docs/                 # User guides — not every action has one
decisions/            # ADRs (records/), imported specs, topic indexes, decision log
AGENTS.md             # Stub — redirects to CLAUDE.md
CHANGELOG.md          # Release notes
```

For the per-action file list with descriptions, read `SKILL.md` — it is the canonical name→path mapping. This tree deliberately stops at directories so it cannot drift from the repo.

## Before Every Commit

**Scope: the integrating commit only.** This ritual belongs to whoever commits the change into the integration branch — in the work pipeline, the queue owner at Step 9. A builder committing on its own `worktree-agent-*` branch **skips it entirely**: `actions/version.md` and `CHANGELOG.md` are serial-only files owned by the integrator (`actions/work-reference.md` → Worktree Dispatch Mode → Fan-Out Dispatch), and a builder bumping either would race every sibling. This file auto-loads in any session rooted here, a builder's worktree included, so the exemption has to live in the rule rather than in each builder's brief — a rule that every brief must override will eventually meet a brief that forgets.

1. **Bump the version** in `actions/version.md` (line starting with `**Current version**:`). Use semver — patch for fixes, minor for features, major for breaking changes. When in doubt, patch. **Verify the new version number is strictly greater than the first existing entry in `CHANGELOG.md`** — duplicate version numbers have occurred before.

2. **Add a changelog entry** at the top of `CHANGELOG.md` (below the header). The title must **say what was delivered** — a reader scanning only headings should know what changed ("Board View Filters", not "The Fine Sieve"). No whimsical codenames. **Verify the title is not already used** by an earlier entry. (Historical entries were retroactively retitled to this convention in 0.117.1.)

```markdown
## X.Y.Z — [Short Descriptive Title] (YYYY-MM-DD)

[1-2 casual sentences — what changed and why it matters.]

- [Bullet points for specifics]
```

Keep it brief, newest on top, lead with value not implementation. Every version gets an entry.

## Action File Conventions

**Every NEW action must justify not being a sibling skill.** Before an action file is added, state — in its description blockquote or an accompanying ADR — why it belongs inside do-work rather than in a separate skill: what queue/pipeline machinery it needs, or which existing action it completes. Reviewers reject additions without this justification. (Ratchet from the 2026-07 bloat cleanup: bkb, interview, dream, and the prompt library accreted ~47k words with no such gate — see `decisions/audits/2026-07-15-relocation-extraction-plans.md`.) Every new action also adds a routing row, a dispatch row, a help-menu block, and a next-steps block — SKILL.md's word budget is enforced by `_dev/tests/contract-regressions.sh`, and the answer to hitting it is a merge or lazy-load, not a bigger budget.

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

Cross-reference other actions by their **file path** (e.g., `actions/work.md`, or `actions/work-reference.md`'s Schema Read Contract) so an agent reading the file can open the target directly without resolving a name to a path. Shipped files must never cite this repo's own `CLAUDE.md` or `AGENTS.md` — both are export-ignored (they're the maintainer doc, and nested CLAUDE.md auto-loads into consumer agents' context), so the citation dangles downstream; restate the rule inline or point at a shipped home instead. `_dev/tests/contract-regressions.sh` enforces this by flagging **any** mention of either file in a shipped path — not by matching citation idioms, which caught 0 of 8 real occurrences before being inverted. References to a *consumer project's* CLAUDE.md (prime routing, tidy-repo, KB schema files) are fine, and that exemption is recorded as a **per-file** allowlist in the check itself; a new shipped file mentioning the maintainer doc fails until someone decides which of the two it is. Companion reference files take a path too (`actions/interview-reference.md`, `actions/bkb-reference.md`). The one exception is a `do-work <verb>` **command invocation** (`do-work run`, `do-work clarify`) — that's how an action is _run_, not a pointer to its file, so keep it as a command. SKILL.md remains the authoritative name→path mapping and may use short names in its routing prose.

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

Just-in-time rules live in `crew-members/[name].md`. Each file's `JIT_CONTEXT` comment is the canonical statement of when it loads — that comment is the contract, not any list duplicated here or elsewhere. The loading order for the work pipeline is specified in `actions/work.md` Step 6.

- `general.md` and `coding-guardrails.md` are always loaded during implementation. Everything else loads conditionally per its `JIT_CONTEXT` (e.g., domain match, `tdd`/`caveman` flags, security surface, fan-out, human-facing artifact production, third-party content ingestion, skill-instruction maintenance passes, debugging retries, interviews).
- Four contracts worth knowing without opening files: `clear-questions.md` loads before presenting the user any **interactive question** (ask-tool prompt, clarifying question, option menu) and governs question wording — one decision per question, no unglossed shorthand, options that state their consequence; `anti-slop.md` loads before producing any **human-facing artifact**; `prompt-injection.md` loads before ingesting any **content not authored by the current invocation or the shipped skill files**; `maintenance.md` loads before a **deliberate maintenance pass on the skill's own instructions** (fixing a drifting agent/action/crew/prime file, where removing or narrowing is a candidate fix) and codifies delete-before-you-add — the maintenance-time complement to `coding-guardrails.md`'s implementation-time surgical-changes rule. In the work pipeline the trigger is the `maintenance: true` REQ marker (set by capture for a removal/narrowing finding on the skill's own instructions; loaded by `actions/work.md` Step 6) — marker-only, never a description heuristic, which would misfire on ordinary implementation REQs. New actions that hit any of these triggers must load the corresponding file.
- If a rules file is missing, proceed without it — never block on a missing rules file.

## Queue Path Convention

Pending REQ files live in `do-work/queue/` — not `do-work/` root.

## Shipped Tooling (`tools/`)

`tools/queue-kanban/` is a standalone Go module (its own `go.mod`, embedded `web/` frontend) that renders the `do-work/` queue as a Kanban board. It ships in the tarball (it is **not** `export-ignore`'d) so `do-work update` carries it into every consumer; the `do-work board` action (`actions/board.md`) builds and runs it. Conventions:

- **Versioning is folded into the skill.** The tool has no independent changelog — its changes get entries in the root `CHANGELOG.md` and a normal skill version bump, exactly like any action. (It was independently versioned through 1.1.0 before being vendored in; that history lives in `decisions/records/adr-016-*`.)
- **Keep the parser in lock-step with the schema.** The board buckets tickets by the `status` vocabulary defined in `actions/work-reference.md`'s Schema Read Contract; `depends_on`, `domain`, `route`, `write_set`, `assigned_to`, and the blocked fields (`blocked_by`/`blocked_at`/`blocked_check`) are parsed for display only (badges, drawer metadata — no column logic; a `status: blocked` card is routed by its `status` value alone, and the board never runs `blocked_check` — the work pipeline does). `write_set` feeds one derived, display-only overlap annotation (`annotateWriteSetOverlap`, run after bucketing) behind the `overlaps` badge and its drawer row — never column logic, never scheduling. Nothing schedules on `write_set` at all, and that does not depend on how many REQs are in flight: under fan-out dispatch it is advisory input to a human's pick and the merge is the non-interference proof (`actions/work-reference.md` → Worktree Dispatch Mode → Fan-Out Dispatch). The field is purely a display input. `assigned_to` is the same class one step further out: the board only badges it, and the single reader that *acts* on it is the work pipeline's default scan, which skips-and-reports as a courtesy and is overridden by explicit targeting — never the board. Any change to that contract must be mirrored in `tools/queue-kanban/model.go` (and vice-versa) in the same commit — co-location is the whole point. The same lock-step applies to the testing placeholders (`testing_status` / `tested_by` / `testing_updated_at` / `testing_feedback`, mirrored in `tools/queue-kanban/testing.go`).
- **The tool has exactly two write surfaces, and neither touches pipeline state.** (1) The board's Testing view writes only the testing placeholders above plus `do-work/testers.md`. (2) `queue-kanban next-version` rewrites the single `**Current version**: X.Y.Z` line in one version file (`actions/version.md` by default, `--version-file` to point elsewhere) — outside `do-work/` entirely. Everything else the tool does is read-only: `summary`, `generate`, `serve`'s board views, `next-req`, `verify`, and `now` (which reads a clock, not even the tree). Nothing in the tool writes `status`, any other pipeline field, any file in `do-work/` beyond `testers.md`, or **`CHANGELOG.md`** — the changelog stays an owner-only, human-authored write, because unique version numbers do not make a shared prepend safe. Adding a third write surface means amending this sentence in the same commit; that is the co-location rule applied to itself.
- **`verify` is the mechanical half of "Before Every Commit."** `queue-kanban verify` checks items 1 and 2 of that section (version/changelog agreement, entry-title reuse) plus queue, assignment, UR-closure and worktree invariants, and it is wired into `actions/forensics.md` Check 14. It reports and routes; it never repairs — fixes belong to `actions/cleanup.md`, which asks first.
- **Toolchain exception to "design for the floor."** The board is the one capability that needs a compiler (Go, per `tools/queue-kanban/go.mod`). `actions/board.md` precondition-checks `go` and degrades gracefully when it's absent — it never blocks the rest of the skill. Don't reach for a compiled tool in any other action, with one narrow class of exception: a subcommand may be named as the **preferred** source for something an action already obtains a shell-portable way, provided the fallback stays documented and nothing builds the binary to get it. Three qualify today — `next-req`, `next-version`, and `now` (`actions/work-reference.md` → Timestamp rule) — each gated on the binary being *already built* and each falling back to the manual procedure it accelerates. That gate is the whole exception: an action that would compile the tool, or that has no floor path, is the prohibited shape.
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
