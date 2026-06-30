# Note Action

> **Part of the do-work skill.** Invoked when the user wants to jot a lightweight, dated next-step note without going through capture. Appends one line to `do-work/notes.md`; `do-work roadmap` surfaces those lines at the top of its survey.

A note is **not** a REQ. It has no frontmatter, no schema, no RED/GREEN proof, no domain, and triggers no implementation. It is a lightweight hint — "look at X next", "check Y before running", "revisit after Z lands" — that the user deletes directly from `do-work/notes.md` when it's no longer relevant. There is no delete command and no archival: the file is plain text the user edits by hand.

## When to Use

**Use when:**
- The user wants to record an informal next-step thought that doesn't warrant a REQ (`do-work note "investigate prototype xyz.html"`).
- Capturing planning context for the next `do-work roadmap` without committing to implementation work.

**Do NOT use when:**
- The thought is an actual task to build → use `do-work capture-request: [describe]` (it creates the UR + REQ pairing).
- The user wants to survey or remove existing notes → that's `do-work roadmap` (display) and a manual edit of `do-work/notes.md` (removal).

## Input

`$ARGUMENTS` is the note text — everything after the `note` keyword.

- `do-work note investigate xyz` → text = `investigate xyz`
- `do-work note "investigate xyz"` → strip the surrounding quotes → text = `investigate xyz`
- `do-work note add investigate xyz` → strip a single leading `add ` → text = `investigate xyz`
- `do-work add note investigate xyz` → routing already stripped `add note`; `$ARGUMENTS` = `investigate xyz`

If `$ARGUMENTS` is empty after stripping, do not write an empty note — print the one-line usage (`do-work note "<text>"`) and stop.

## Steps

### Step 1: Normalize the text

1. Take `$ARGUMENTS` verbatim.
2. Strip a single leading `add ` token if present (handles `do-work note add …`).
3. Strip a matching pair of surrounding quotes (`"…"` or `'…'`) if present.
4. Trim leading/trailing whitespace.
5. If the result is empty, print usage and stop (Input section).

### Step 2: Append the dated line

1. Resolve today's date as `YYYY-MM-DD` (e.g. via `date +%F`).
2. Ensure the `do-work/` directory exists first (`mkdir -p do-work/`) — it may be absent if `note` is the first do-work command run in a fresh repo — then create `do-work/notes.md` if it isn't there yet (the file is a plain list — no header, no frontmatter).
3. Append a single line: `- [YYYY-MM-DD] <text>`.

Do not deduplicate, sort, or reformat existing lines — only append. The file stays in chronological (append) order; the user curates it by hand.

### Step 3: Report

Confirm what was added and where, in one or two lines:

```
Noted → do-work/notes.md
  - [2026-06-01] investigate prototype xyz.html

(Surfaces at the top of `do-work roadmap`. Delete the line from do-work/notes.md when it's no longer relevant.)
```

Do **not** create a UR or REQ, do not move into the work loop, and do not run a commit from this action — appending the line is the whole job. `do-work/notes.md` is itself part of the committable Trail of Intent (like URs and REQs); the user commits it in their normal flow, whenever they choose.

## Output Format

- One appended line in `do-work/notes.md` of the form `- [YYYY-MM-DD] <text>`.
- A short confirmation to the user (Step 3).

## Rules

- **A note is not a task.** Never let a note kick off capture, work, or a commit.
- **Append-only.** Don't rewrite, sort, or dedupe existing lines — the user owns the file's contents.
- **No frontmatter, no header.** `do-work/notes.md` is a flat list so `roadmap`'s render is trivial (one `## Notes` line per file line).
- **The action never commits; the file is committable.** `do-work note` only appends — it runs no git command. `do-work/notes.md` is committed alongside the rest of `do-work/` (the Trail of Intent); only `do-work/pipeline.json` and `do-work/runs/` are git-excluded (transient state — kept out of git regardless of install layout, via the shipped `.gitignore` or `.git/info/exclude`). On a merge conflict it's append-only, so keep both sides.
- **Empty input is a no-op** with usage, not an empty `- [date]` line.

## Common Rationalizations

| If you're thinking... | STOP. Instead... | Because... |
|---|---|---|
| "This note sounds like real work — I'll capture it as a REQ too" | Just append the note; suggest `do-work capture` only if the user asks | A note is deliberately lighter than a REQ; auto-promoting it defeats the purpose |
| "I'll tidy up / sort the existing notes while I'm here" | Append only; leave the rest untouched | The user curates `notes.md` by hand; reordering loses their intent |
| "The text is empty, I'll just write `- [date]`" | Print usage and stop | An empty dated bullet is noise in the roadmap |

## Red Flags

- A UR or REQ folder/file was created by a `do-work note` invocation (note must never enter the capture pipeline).
- `do-work/notes.md` gained frontmatter, a header, or reordered lines (should be append-only flat list).
- The note action ran `git commit` (or any git write) itself — it must only append; the user commits `do-work/notes.md` in their own flow.

## Verification Checklist

- [ ] Exactly one `- [YYYY-MM-DD] <text>` line appended to `do-work/notes.md` (file created if absent).
- [ ] Leading `add ` and surrounding quotes stripped; whitespace trimmed.
- [ ] Empty input produced usage output and no file write.
- [ ] No UR, REQ, work-loop transition, or commit was triggered.
