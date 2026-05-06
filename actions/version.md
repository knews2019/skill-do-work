# Version Action

> **Part of the do-work skill.** Handles version reporting, update checks, and work recaps.

**Current version**: 0.71.1

**Upstream**: https://raw.githubusercontent.com/knews2019/skill-do-work/main/actions/version.md

## When to Use

**Use when:**
- The user asks "what version", "release notes", "what's new", or "history" → version + last 5 changelog entries.
- The user asks to "update", "check for updates", or "is there a newer version" → update flow.
- The user asks for a "recap" of recent work → recap flow across archive + active URs.

**Do NOT use when:**
- The user wants to see all changelog entries (more than 5) — point them at `CHANGELOG.md` directly instead of loading the full file.
- The user wants to *install* the skill fresh — that's the README install command, not this action.
- The install is global (under `~/.claude/skills/` etc.) and the user wants an update — refuse the auto-update per the preflight in Step 2 below, and redirect them.

## Responding to Version Requests

When user asks "what version", "version", "what's new", "release notes", "what's changed", "updates", or "history":

1. Report the version shown above
2. **Show last 5 skill releases**:
   - Read the first ~80 lines of `CHANGELOG.md` in the skill's root directory (same level as `SKILL.md`) — do NOT load the full file
   - Extract the 5 most recent version entries (split at `## ` headings, take first 5 blocks)
   - Reverse so newest is at the bottom (right where the user's eyes are)
   - Print them after the version number

## Responding to Update Checks

When user asks "check for updates", "update", or "is there a newer version":

1. **Fetch upstream**: Use your environment's web fetch capability to get the raw version.md from the upstream URL above
2. **Extract remote version**: Look for `**Current version**:` in the fetched content
3. **Compare versions**: Use semantic versioning comparison
4. **Report result** using the format below

### Report Format

**If update available** (remote > local):

1. **Tell the user**: `Update available: v{remote} (you have v{local}).`
2. **Preflight: confirm this is a project-local install, not a global one.** The update must overwrite the copy inside the current project — never a user-wide / global install.
   - Resolve the absolute path of the skill's root directory (where `SKILL.md` lives). Call this `<skill-root>`.
   - **Refuse to auto-update if `<skill-root>` is under any user-wide skills location**, including but not limited to:
     - `~/.claude/skills/...`
     - `~/.gemini/skills/...`
     - `~/.cursor/skills/...`
     - `~/.config/*/skills/...`
     - anything else under `$HOME` that isn't also inside the current project's git repo.
   - Resolve the current project's git root: `git -C <invocation-dir> rev-parse --show-toplevel`. If `<skill-root>` is **not** a descendant of that project root, stop and report:
     ```
     Skill is installed at <skill-root>, which is outside the current project (<project-root>).
     Refusing to update a global/shared install from here. Either:
       - cd into the project that owns <skill-root> and re-run, or
       - install the skill locally inside this project (e.g. <project-root>/.claude/skills/do-work/) and re-run.
     ```
     Do NOT proceed. Do NOT suggest the curl command.
   - Only continue once `<skill-root>` is confirmed to live inside the current project's git root.
3. **Check for local changes** to shipped skill files at `<skill-root>`:
   - **Scope the check to skill-owned files only.** Ignore `do-work/` (queue data, archives, deliverables) — those are generated at runtime and should never block an update.
   - If `<skill-root>` is a git repo, run `git -C <skill-root> status --porcelain -- SKILL.md actions/ crew-members/ prompts/ interviews/ specs/ docs/ decisions/ hooks/ CLAUDE.md AGENTS.md CHANGELOG.md README.md next-steps.md .do-work-upstream-manifest` (listing every shipped editable path) and check for uncommitted changes. Any dirty file in these paths will be clobbered by the tar extraction in step 6 if you proceed.
   - If it's **not** a git repo, check whether shipped skill files (actions/, crew-members/, prompts/, interviews/, specs/, docs/, decisions/, hooks/, SKILL.md, CLAUDE.md, next-steps.md, etc.) differ from a fresh install by looking for user-modified content (custom crew-members, edited action files, local prompt/template additions, ADR edits, etc.).
   - **If any shipped skill files are dirty / have local modifications**: Stop and warn the user. List the modified files and ask for explicit confirmation before proceeding. Do NOT auto-update.
   - **If clean**: Proceed to step 4 (download), step 5 (manifest-driven pre-clean), then step 6 (extract).
4. **Download the tarball to a temp file.** We need to read it twice — once to list its contents (for the manifest), once to extract:
   ```
   TARBALL=$(mktemp -t do-work-update.XXXXXX.tar.gz)
   curl -sLf https://github.com/knews2019/skill-do-work/archive/refs/heads/main.tar.gz -o "$TARBALL"
   ```
   `-f` makes curl exit non-zero on HTTP errors instead of saving an HTML error page. If the download fails, fall through to the "fetch fails" branch below.
5. **Pre-clean discoverable directories using the manifest union.** `prompts/` and `interviews/` are enumerated by `do-work prompts list` and `do-work interview list` (they glob `prompts/*.md` and `interviews/*.md`), so any upstream-removed file that stays on disk would still appear as a live workflow. The pre-clean must (a) remove files upstream owns *now or used to own* — so renames and removals are reflected — and (b) **preserve user-authored files** that upstream has never shipped. Do this by intersecting deletes with the union of two manifests:
   ```
   # New manifest from the incoming tarball (top-level .md under prompts/ and interviews/)
   tar tzf "$TARBALL" | awk -F/ 'NF==3 && ($2=="prompts" || $2=="interviews") && $3 ~ /\.md$/ {print $2"/"$3}' \
       > "$TARBALL.new-manifest"

   # Old manifest from the currently-installed skill — may be missing on a pre-fix install
   if [ -f <skill-root>/.do-work-upstream-manifest ]; then
       cp <skill-root>/.do-work-upstream-manifest "$TARBALL.old-manifest"
   else
       : > "$TARBALL.old-manifest"   # bootstrap: empty old manifest
   fi

   sort -u "$TARBALL.old-manifest" "$TARBALL.new-manifest" > "$TARBALL.union"

   # Delete only files in the union — user-authored files (in neither manifest) survive.
   while IFS= read -r relpath; do
       rm -f "<skill-root>/$relpath"
   done < "$TARBALL.union"
   ```
   User-authored files are preserved by **omission from the manifest**, not by `git status` — so committed custom prompts/interviews are safe even though the dirty check at Step 3 wouldn't flag them. Do NOT delete files in `prompts/` or `interviews/` subdirectories — only the top-level `.md` files are globbed by the listing actions, and the manifest only tracks those. Do NOT touch `do-work/` or any other runtime directory.

   **Bootstrap (pre-fix installs).** If `.do-work-upstream-manifest` is missing, the old-manifest is empty and the union equals the new manifest. That preserves user-authored files whose names don't collide with current upstream filenames; from the next update onwards the manifest is on disk and steady-state behavior applies.
6. **Run the extraction in place at `<skill-root>`** (the project-local path confirmed in step 2). `cd` there first so the extraction cannot land in a global directory by mistake:
   ```
   cd <skill-root> && tar xzf "$TARBALL" --strip-components=1 --exclude='_dev'
   rm -f "$TARBALL" "$TARBALL.new-manifest" "$TARBALL.old-manifest" "$TARBALL.union"
   ```
   The tarball includes a fresh `.do-work-upstream-manifest`, so the manifest on disk after extraction is automatically up to date for the next update. **Note:** tar extraction adds and overwrites files but does not delete files removed upstream. For non-discoverable directories (`actions/`, `crew-members/`, `specs/`, `docs/`, `decisions/`) leftovers are harmless — the skill only loads files it references by name. For `prompts/` and `interviews/`, the manifest-driven pre-clean above is what prevents ghost entries; if you skipped it, run `do-work prompts list` and `do-work interview list` after updating and delete anything that looks obsolete. Never delete `do-work/` (runtime state).
7. **Verify**: Read `<skill-root>/actions/version.md` again and confirm the local version now matches the remote version.
8. **Report result**: `Updated to v{remote} at <skill-root>.`

Do NOT just print the curl command and ask the user to run it. You are the agent — run it yourself.

**If up to date** (local >= remote):

```
You're up to date (v{local})
```

**If fetch fails**:

```
Couldn't check for updates.
```

Attempt the update anyway using the download + manifest-driven pre-clean + extraction sequence in steps 4–6 above (still respecting the preflight location check in step 2 and the dirty-tree check in step 3 — refuse if the install is global). If that also fails, report the error and provide the manual command as a last-resort fallback:

```
To manually update, cd into the **project-local** skill root (where SKILL.md lives inside *this* project — NOT ~/.claude/skills/, ~/.gemini/skills/, or any other global skills directory) and run:

cd <project-root>/path/to/skill-do-work
curl -sL https://github.com/knews2019/skill-do-work/archive/refs/heads/main.tar.gz | tar xz --strip-components=1 --exclude='_dev'

Or visit: https://github.com/knews2019/skill-do-work
```

Note: this last-resort command skips the manifest-driven pre-clean, so any prompt or interview removed by upstream may linger on disk. Run `do-work prompts list` and `do-work interview list` after a manual update and delete anything that looks obsolete.

## Responding to Recap Requests

When user asks "recap":

1. **Archive source** (`do-work/archive/UR-*/`): Read as before — title from `input.md`, REQs from `REQ-*.md` files inside each UR folder.
2. **Active source** (`do-work/user-requests/UR-*/`): Read `input.md` for the title. For REQs, scan `do-work/queue/REQ-*.md` files whose `user_request:` frontmatter field matches the UR id (e.g., `user_request: UR-143`). Also check `do-work/working/` for claimed REQs belonging to the UR.
3. **Merge**: Combine both lists, deduplicate by UR id (archive version wins if both exist), sort by UR number descending, take top 5.
4. **Label each UR**:
   - No label if fully archived
   - `(pending)` if the UR has any pending REQs
   - `(completed, awaiting archive)` if all its REQs are completed/done but the UR isn't archived yet
5. **Format as a "Recent Work" section**:
   ```
   ## Recent Work

   UR-144 — Block-level improved translation for ZH pairs
     REQ-361 — Block-level improved translation
   UR-143 — Model selector thinking variants (completed, awaiting archive)
     REQ-360 — Model selector thinking variants
   UR-142 — Quality-Score-Driven Repair Loop (completed, awaiting archive)
     REQ-359 — Quality-Score-Driven Repair Loop
   UR-011 — Dark mode implementation
     REQ-043 — Theme store setup
     REQ-044 — Settings panel toggle
   ```
   One line per UR, one indented line per REQ. No descriptions, no scores, no file lists.
6. **If no archive exists AND no active URs found**: Print `No completed work yet.` and skip this section.

## Red Flags

- The update flow is about to `cd` into a path under `~/.claude/skills/`, `~/.gemini/skills/`, or anywhere outside the current project's git root — STOP. Global installs must never be auto-updated from here.
- Remote version fetched from the upstream URL is empty, malformed, or older than local — abort; don't "update" backwards.
- The dirty-tree check reported modifications but the update proceeded anyway — user's local customizations will be clobbered.
- Recap lists the same UR twice (once from archive, once from active) — the dedup step was skipped; archive version should win.
- Version reported doesn't match the `**Current version**:` line at the top of this file — caching or path confusion; re-read the file from disk.

## Verification Checklist

- [ ] Version output shows the local version, the last 5 changelog entries, newest at the bottom.
- [ ] Update flow refused to proceed for global installs (Step 2 preflight).
- [ ] Update flow refused to proceed when shipped files had uncommitted changes (Step 3 dirty check), unless user explicitly confirmed.
- [ ] Update flow pre-cleaned `prompts/*.md` and `interviews/*.md` (Step 4) before tar extraction.
- [ ] Post-update verification re-read `actions/version.md` and confirmed the local version matches remote.
- [ ] Recap merged archive + active sources, deduped by UR id, kept the archive version on conflicts.
- [ ] Recap output is one line per UR + one indented line per REQ — no scores, no file lists, no descriptions.
