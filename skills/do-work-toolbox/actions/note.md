# Note Action

> **Part of the do-work-toolbox skill.** Invoked when the user wants to jot a lightweight, dated next-step note without going through capture. Appends one line to `do-work/notes.md`; `do-work roadmap` surfaces those lines at the top of its survey, and `do-work-board board` renders them as a Notes strip above the columns.

A note is **not** a REQ. It has no frontmatter, no schema, no RED/GREEN proof, no domain, and triggers no implementation. It is a lightweight hint — "look at X next", "check Y before running", "revisit after Z lands" — that the user deletes directly from `do-work/notes.md` when it's no longer relevant. There is no delete command and no archival: the file is plain text the user edits by hand.

## When to Use

**Use when:**
- The user wants to record an informal next-step thought that doesn't warrant a REQ (`do-work-toolbox note "investigate prototype xyz.html"`).
- Capturing planning context for the next `do-work roadmap` without committing to implementation work.

**Do NOT use when:**
- The thought is an actual task to build → use `do-work capture-request: [describe]` (it creates the UR + REQ pairing).
- The user wants to survey or remove existing notes → that's `do-work roadmap` (display) and a manual edit of `do-work/notes.md` (removal).

## Input

`$ARGUMENTS` is the note text — everything after the `note` keyword.

- `do-work-toolbox note investigate xyz` → text = `investigate xyz`
- `do-work-toolbox note "investigate xyz"` → strip the surrounding quotes → text = `investigate xyz`
- `do-work-toolbox note add investigate xyz` → strip a single leading `add ` → text = `investigate xyz`
- `do-work-toolbox note investigate xyz` → routing already stripped `add note`; `$ARGUMENTS` = `investigate xyz`

If `$ARGUMENTS` is empty after stripping, do not write an empty note — print the one-line usage (`do-work-toolbox note "<text>"`) and stop.

## Steps

### Step 1: Delegate the mechanical note write

Resolve `<project-root>` and the installed core `<skill-root>`, then invoke the canonical command once:

```bash
<skill-root>/tools/do-work-cli.sh --repo-root <project-root> --format json do-work-note "$ARGUMENTS"
```

The command owns normalization, the date, directory creation, exact append bytes, Git target guards, rollback, and the empty-input refusal. Treat its typed changes and exact text output as the complete write result. Missing, failed, or malformed canonical tooling stops actionably; do not fall back to direct prose or shell mutation.

### Step 2: Preserve the action boundary

Do not deduplicate, sort, reformat, capture, or commit. The action retains only the judgment that the user's input is a lightweight note rather than a task; the canonical command performs the deterministic append.

### Step 3: Report

Confirm what was added and where, in one or two lines:

```
Noted → do-work/notes.md
  - [2026-06-01] investigate prototype xyz.html

(Surfaces at the top of `do-work roadmap` and in the Notes strip of `do-work-board board`. Delete the line from do-work/notes.md when it's no longer relevant.)
```

Do **not** create a UR or REQ, do not move into the work loop, and do not run a commit from this action — appending the line is the whole job. `do-work/notes.md` is itself part of the committable Trail of Intent (like URs and REQs); the user commits it in their normal flow, whenever they choose.

## Output Format

- One appended line in `do-work/notes.md` of the form `- [YYYY-MM-DD] <text>`.
- A short confirmation to the user (Step 3).

## Rules

- **A note is not a task.** Never let a note kick off capture, work, or a commit.
- **Append-only.** Don't rewrite, sort, or dedupe existing lines — the user owns the file's contents.
- **Write bullets, never frontmatter.** `do-work-toolbox note` appends a `- [YYYY-MM-DD] <text>` bullet and nothing else — it never adds a heading, a preamble, or frontmatter of its own.
- **The bullet is what makes a line a note.** Users do add a `#` heading, a prose preamble, and `<!-- ... -->` blocks parking pruned entries, so every *reader* of this file (`../../do-work/actions/roadmap.md` Step 0, the `do-work-board board` Notes strip) strips HTML comments first, then renders only the bullet lines. A reader that treats every non-blank line as a note renders the boilerplate — and the bullets buried inside the pruned-entries comment — as notes.
- **The action never commits; the file is committable.** `do-work-toolbox note` only appends — it runs no git command. `do-work/notes.md` is committed alongside the rest of `do-work/` (the Trail of Intent). `do-work/runs/` is committable while a run is live and is deleted once its findings are consumed (see `crew-members/background-agents.md`). On a merge conflict it's append-only, so keep both sides.
- **Empty input is a no-op** with usage, not an empty `- [date]` line.
- **One canonical writer.** Only `do-work-note` may append from this action; missing, failed, or malformed command output never permits a manual fallback.

(The Rules above are the complete contract — every guard this action needs is stated there once.)
