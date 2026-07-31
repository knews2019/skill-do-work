#!/usr/bin/env bash
# Updates a project-local do-work install. Invoked by `just run-do-work-update`.
set -euo pipefail

upstream_url='https://github.com/knews2019/skill-do-work/archive/refs/heads/main.tar.gz'
# `justfile` is deliberately NOT in this list — it is appended below for a nested install
# only. See the root-fallback note after the skill/project path resolution.
shipped_paths=(SKILL.md actions crew-members prompts interviews specs docs hooks tools CHANGELOG.md README.md next-steps.md)

fail() {
  printf 'do-work update: %s\n' "$*" >&2
  exit 1
}

version_is_newer() {
  awk -v local_version="$1" -v remote_version="$2" '
    BEGIN {
      split(local_version, local_parts, ".")
      split(remote_version, remote_parts, ".")
      for (part_number = 1; part_number <= 3; part_number++) {
        local_part = (part_number in local_parts) ? local_parts[part_number] + 0 : 0
        remote_part = (part_number in remote_parts) ? remote_parts[part_number] + 0 : 0
        if (remote_part > local_part) exit 0
        if (remote_part < local_part) exit 1
      }
      exit 1
    }
  '
}

# Files that legitimately exist only in an install and are not an upstream deletion: build
# output, OS and editor droppings, merge leftovers. Everything else that is install-only gets
# surfaced rather than dropped — see append_install_diff.
is_transient_local_extra() {
  case "$1" in
    tools/queue-kanban/queue-kanban) return 0 ;;   # the compiled board binary (gitignored, rebuilt per run)
    *.DS_Store|*.sw[a-p]|*.orig|*.rej|*~) return 0 ;;
    *) return 1 ;;
  esac
}

# Splits the installed-versus-upstream comparison into two streams, because the two things
# `diff` reports on one `Only in <install>` line have opposite severities.
#
#   $destination      real divergence — content diffs and files upstream ships that the
#                     install lacks. The post-update caller treats ANY line here as fatal.
#   $extras_dest      paths present only in the install. NEVER fatal: a local-only file must
#                     not abort an update (that regression was fixed once already). But
#                     dropping them wholesale — the previous `grep -vF "Only in <install>"` —
#                     is how a file upstream DELETED survives forever downstream: the tar
#                     extraction only overwrites, it never removes, so nothing else in this
#                     script can see a stale shipped action or check. So classify: known
#                     droppings are dropped, the rest are reported for a human to judge.
#                     Reported, never auto-removed — a consumer's own file dropped into a
#                     shipped directory is indistinguishable from a stale one at this level.
append_install_diff() {
  local fresh_root="$1"
  local installed_root="$2"
  local destination="$3"
  local extras_dest="$4"
  local shipped_path diff_line entry_directory entry_basename relative_entry

  for shipped_path in "${shipped_paths[@]}"; do
    if [ -e "$fresh_root/$shipped_path" ] && [ ! -e "$installed_root/$shipped_path" ]; then
      # Upstream ships this but it's wholly absent from the install: `diff` on a
      # missing argument only errors to stderr, so flag it here instead.
      printf 'Missing from install (shipped upstream): %s\n' "$shipped_path" >> "$destination"
    elif [ ! -e "$fresh_root/$shipped_path" ] && [ -e "$installed_root/$shipped_path" ]; then
      # The whole top-level path is gone upstream — the same "stale shipped file" case as the
      # per-entry one below, just at the coarsest grain. Also not a diff failure.
      printf '%s\n' "$shipped_path" >> "$extras_dest"
    elif [ -e "$fresh_root/$shipped_path" ]; then
      while IFS= read -r diff_line; do
        case "$diff_line" in
          # The install root is inside double quotes so a regex/glob metachar in the path
          # cannot widen this pattern — the same reason the old filter used `grep -vF`.
          "Only in $installed_root"*)
            # "Only in <dir>: <name>" → rebuild the install-relative path to classify on.
            entry_basename="${diff_line##*: }"
            entry_directory="${diff_line#Only in }"
            entry_directory="${entry_directory%: *}"
            entry_directory="${entry_directory#"$installed_root"}"
            entry_directory="${entry_directory#/}"
            relative_entry="$entry_basename"
            [ -n "$entry_directory" ] && relative_entry="$entry_directory/$entry_basename"
            is_transient_local_extra "$relative_entry" \
              || printf '%s\n' "$relative_entry" >> "$extras_dest"
            ;;
          *) printf '%s\n' "$diff_line" >> "$destination" ;;
        esac
      done < <(diff -ru "$fresh_root/$shipped_path" "$installed_root/$shipped_path")
    fi
  done
}

project_root=''
if [ "$#" = 2 ] && [ "$1" = '--project-root' ]; then
  project_root="$2"
else
  fail 'usage: do-work-update.sh --project-root <project-root>'
fi

[ -d "$project_root" ] || fail "project root does not exist: $project_root"
project_root="$(cd "$project_root" && pwd -P)"
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
skill_root="$(cd "$script_dir/.." && pwd -P)"

case "$skill_root" in
  "$project_root"|"$project_root"/*) ;;
  *) fail "skill is outside this project ($skill_root is not within $project_root); refusing to update a shared install" ;;
esac

# `justfile` is project-owned: `actions/install.md` → Workflow: `just-kanban` states that
# `do-work update` never touches it. In a NESTED install the tarball's justfile extracts
# inside the skill directory — the skill's own copy, refreshed and reviewed like any other
# shipped file. In the ROOT FALLBACK (skill_root IS project_root — the second branch of the
# `just run-do-work-update` recipe) that very same path is the PROJECT's justfile, holding
# the project's own recipes. So it stays out of the reviewed-and-verified shipped set here,
# and is preserved byte-for-byte across the extraction below.
root_fallback_install=''
project_justfile_name=''
if [ "$skill_root" = "$project_root" ]; then
  root_fallback_install=1
  # Which name this project's justfile actually uses — same candidate order as
  # `actions/install.md` → Workflow: `just-kanban` Phase 1. The name matters as much as the
  # bytes: on a case-INSENSITIVE filesystem (APFS/HFS+ by default) the tarball's lowercase
  # `justfile` resolves to an existing `Justfile`, so extracting clobbers the project's file
  # AND renames it; restoring to the recorded name puts both back. For the same reason the
  # candidate test reads the directory's REAL entry names instead of
  # `[ -f "$project_root/$candidate" ]`, which resolves case-insensitively and would always
  # pick the first candidate — renaming the very file this is meant to leave alone.
  project_root_entries="$(ls -A "$project_root" 2>/dev/null || true)"
  for justfile_candidate in justfile Justfile .justfile; do
    if printf '%s\n' "$project_root_entries" | grep -qxF "$justfile_candidate" \
       && [ -f "$project_root/$justfile_candidate" ]; then
      project_justfile_name="$justfile_candidate"
      break
    fi
  done
else
  shipped_paths+=(justfile)
fi

[ -s "$skill_root/SKILL.md" ] || fail "SKILL.md is missing at $skill_root"
[ -f "$skill_root/actions/version.md" ] || fail "actions/version.md is missing at $skill_root"

local_version="$(sed -n 's/^\*\*Current version\*\*: *\([0-9][0-9.]*\).*/\1/p' "$skill_root/actions/version.md" | head -n 1)"
[ -n "$local_version" ] || fail 'could not read the local version'

update_tmp="$(mktemp -d "${TMPDIR:-/tmp}/do-work-update.XXXXXX")"
install_modified=''
update_verified=''
shipped_files_tracked=''
# Declared here, before the EXIT trap is installed, because print_recovery_instructions reads
# it: under `set -u` a trap firing before the git probe below would abort on an unset variable.
dirty_files=''

# This updater keeps no rollback copy: version control is the undo, and duplicating a tracked
# tree on every run buys nothing git does not already hold. The cost is that recovery from a
# mid-update failure is the operator's call, so a partial install has to say exactly what to
# run rather than exit quietly and let a half-old, half-new tree pass for a clean cancel.
print_recovery_instructions() {
  local candidate_path
  local candidate_paths=("${shipped_paths[@]}")
  local tracked_paths=()
  printf 'The install at %s may be partially updated (some files v%s, some v%s).\n' \
    "$skill_root" "$local_version" "$remote_version" >&2
  if [ -z "$shipped_files_tracked" ]; then
    printf 'These files are not tracked in git here, so there is nothing to restore from.\n' >&2
    printf 'Re-run the update once the cause is fixed — the extraction overwrites in place, so repeating it is safe.\n' >&2
    return
  fi
  # Beyond the shipped set: the two paths a failure can damage that the shipped set does not
  # name. Root fallback — the project's own justfile, deliberately excluded from the reviewed
  # set, which a failure between the extraction and its restore leaves holding the skill's
  # recipes. Nested — the stale vendored maintainer docs the update deletes on purpose, worth
  # offering back on a failed run since a re-run deletes them again.
  if [ -n "$root_fallback_install" ]; then
    candidate_paths+=(justfile)
    if [ -n "$project_justfile_name" ] && [ "$project_justfile_name" != 'justfile' ]; then
      candidate_paths+=("$project_justfile_name")
    fi
  else
    candidate_paths+=(CLAUDE.md AGENTS.md)
  fi
  # Only the paths git actually has, so the printed command runs instead of dying on an
  # unmatched pathspec (an install may legitimately be missing a shipped directory, and the
  # extra candidates above are often untracked or absent).
  for candidate_path in "${candidate_paths[@]}"; do
    if [ -n "$(git -C "$skill_root" ls-files -- "$candidate_path")" ]; then
      tracked_paths+=("$candidate_path")
    fi
  done
  printf 'Restore the tracked skill files from git:\n  git -C %s checkout -- %s\n' \
    "$skill_root" "${tracked_paths[*]}" >&2
  # What "git is the undo" does NOT cover. git restores COMMITTED content, so any edit that was
  # uncommitted when this run started is already gone — the extraction overwrote it and there is
  # no copy. Say so here rather than let the checkout above read as a full undo: the operator is
  # about to run it, and would otherwise expect their customizations back.
  if [ -n "$dirty_files" ]; then
    printf 'NOTE: these shipped files had uncommitted edits when the update started:\n%s\n' "$dirty_files" >&2
    printf 'The extraction has already overwritten them and no copy was kept, so the checkout above restores the COMMITTED content — not those edits. Recover them from your editor history or re-apply by hand.\n' >&2
  fi
  printf 'Then review what the extraction added and delete what you do not want (-nd lists, -fd deletes):\n  git -C %s clean -nd -- %s\n' \
    "$skill_root" "${tracked_paths[*]}" >&2
  if [ -n "$root_fallback_install" ]; then
    printf 'This is a root install (the skill IS the project root), so those paths also hold files your project owns — read the -nd list before running -fd.\n' >&2
  fi
}

cleanup() {
  rm -rf "$update_tmp"
  # `update_verified` is set on the fully-verified path only, so EVERY earlier exit inside the
  # destructive region lands here — a mid-extract ENOSPC, a failed post-update version check,
  # `set -e` firing anywhere in between. Nothing is undone automatically; the contract is to
  # report the partial state loudly and hand the operator the exact recovery commands.
  if [ -n "$install_modified" ] && [ -z "$update_verified" ]; then
    printf 'Update did not complete.\n' >&2
    print_recovery_instructions
  fi
}
trap cleanup EXIT
upstream_tarball="$update_tmp/upstream.tar.gz"
fresh_upstream="$update_tmp/fresh"
mkdir -p "$fresh_upstream"

printf 'Checking do-work updates…\n'
curl -fsSL -o "$upstream_tarball.download" "$upstream_url" \
  || fail 'upstream tarball download failed; no files were changed'
mv "$upstream_tarball.download" "$upstream_tarball"
tar xzf "$upstream_tarball" -C "$fresh_upstream" --strip-components=1 \
  --exclude='_dev' --exclude='do-work' --exclude='ai-reports' --exclude='.vscode' --exclude='decisions'

remote_version="$(sed -n 's/^\*\*Current version\*\*: *\([0-9][0-9.]*\).*/\1/p' "$fresh_upstream/actions/version.md" | head -n 1)"
[ -n "$remote_version" ] || fail 'could not read the upstream version'

if ! version_is_newer "$local_version" "$remote_version"; then
  printf "You're up to date (v%s)\n" "$local_version"
  exit 0
fi

printf 'Update available: v%s (you have v%s).\n' "$remote_version" "$local_version"

dirty_files=''
if git -C "$skill_root" rev-parse --git-dir >/dev/null 2>&1; then
  dirty_files="$(git -C "$skill_root" status --porcelain -- "${shipped_paths[@]}")"
  # Whether git can actually serve as the undo. A repo whose ignore rules cover the skill —
  # a project that gitignores `.claude/`, say — reports a git dir while tracking none of
  # these files, so the presence of a repo is not the question; tracked content is.
  if [ -n "$(git -C "$skill_root" ls-files -- "${shipped_paths[@]}")" ]; then
    shipped_files_tracked=1
  fi
fi

diff_file="$update_tmp/install.diff"
extras_file="$update_tmp/install-extras.txt"
append_install_diff "$fresh_upstream" "$skill_root" "$diff_file" "$extras_file"

if [ -n "$dirty_files" ]; then
  printf 'Shipped skill files have uncommitted changes:\n%s\n' "$dirty_files" >&2
  # The warning that has to land BEFORE the prompt, because this is the point of no return for
  # this content. The extraction overwrites these paths and no rollback copy is kept, so an
  # uncommitted edit here is unrecoverable afterwards — git can only give back what was
  # committed. This is the one guarantee the removed `cp -R` snapshot used to provide.
  printf 'Continuing OVERWRITES those edits with the upstream files, and no rollback copy is kept — git can only restore what is committed, so uncommitted work here is gone for good.\n' >&2
  printf 'Keep them by committing first, or stash them:\n  git -C %s stash push -- %s\n' \
    "$skill_root" "${shipped_paths[*]}" >&2
fi
if [ -s "$diff_file" ]; then
  printf 'Reviewing installed-versus-upstream skill changes before overwrite:\n'
  cat "$diff_file"
fi
if [ -s "$extras_file" ]; then
  printf 'Present in the install but not shipped upstream — the update will NOT remove these:\n'
  sed 's/^/  /' "$extras_file"
  printf 'Upstream may have deleted them; a stale shipped file keeps being read and never updates again.\n'
  printf 'Delete by hand the ones you recognise as ours.\n'
  if [ -n "$root_fallback_install" ]; then
    printf 'This is a root install (the skill IS the project root), so files your project owns inside these directories are listed here too — leave those alone.\n'
  fi
fi

if [ -z "$shipped_files_tracked" ]; then
  printf 'These skill files are not tracked in git here, and this updater keeps no rollback copy: if the update fails partway there is nothing to restore from. Commit or copy the install yourself first if that matters.\n' >&2
fi

printf 'Continue with the update? Files are overwritten in place and no rollback copy is kept. [y/N] '
# EOF / non-interactive stdin (piped, CI, </dev/null) makes `read` return non-zero; under
# `set -e` a bare `read` would abort here before the case below can default to No. Treat
# EOF as an explicit No so the cancel path runs and the script exits 0.
read -r confirmation || confirmation=''
case "$confirmation" in
  y|Y|yes|YES) ;;
  *) printf 'Update cancelled; no files were changed.\n'; exit 0 ;;
esac

preserved_justfile=''
if [ -n "$project_justfile_name" ]; then
  preserved_justfile="$update_tmp/project-justfile"
  cp -p "$project_root/$project_justfile_name" "$preserved_justfile" \
    || fail "could not preserve the project justfile at $project_root/$project_justfile_name"
fi

# First write into the install: from here to the verification below, a failure leaves a
# partial install and the EXIT trap prints the recovery commands.
install_modified=1
find "$skill_root/prompts" -maxdepth 1 -name '*.md' ! -name 'README.md' -delete 2>/dev/null || true
find "$skill_root/interviews" -maxdepth 1 -name '*.md' -delete 2>/dev/null || true
tar xzf "$upstream_tarball" -C "$skill_root" --strip-components=1 \
  --exclude='_dev' --exclude='do-work' --exclude='ai-reports' --exclude='.vscode' --exclude='decisions'

# Root fallback only: undo the extraction's write to the project-owned justfile. Always
# delete the extracted lowercase `justfile` FIRST, then put the project's own file back under
# the name it had. That order is what makes this correct on a case-insensitive filesystem,
# where the extracted `justfile` and a project `Justfile` are one and the same file — delete
# resolves to the clobbered file, and the restore recreates the original name and bytes. On a
# case-sensitive filesystem the delete removes a file the project never had (leaving one would
# plant a second, do-work-authored justfile beside its `Justfile`/`.justfile`) and the restore
# rewrites an untouched file with identical content. Restoring a saved copy rather than passing
# `tar --exclude` keeps this independent of tar dialect and of the archive's top-level
# directory name — `--exclude=justfile` matches that basename at ANY depth, the `diff -x` trap.
if [ -n "$root_fallback_install" ]; then
  rm -f "$project_root/justfile"
  if [ -n "$preserved_justfile" ]; then
    cp -p "$preserved_justfile" "$project_root/$project_justfile_name" \
      || fail "could not restore the project justfile at $project_root/$project_justfile_name"
  fi
fi

# Remove stale vendored maintainer docs (older installs shipped CLAUDE.md/AGENTS.md into
# the skill dir before they were export-ignored from the tarball). ONLY in a nested
# install: when skill_root IS the project root (the justfile recipe's root fallback),
# these are the project's own instruction files, not the skill's — deleting them would
# destroy project instructions, so leave them in place.
if [ -z "$root_fallback_install" ]; then
  rm -f "$skill_root/CLAUDE.md" "$skill_root/AGENTS.md"
fi

installed_version="$(sed -n 's/^\*\*Current version\*\*: *\([0-9][0-9.]*\).*/\1/p' "$skill_root/actions/version.md" | head -n 1)"
[ "$installed_version" = "$remote_version" ] \
  || fail "post-update verification failed (expected v$remote_version, found v${installed_version:-unknown})"

post_diff="$update_tmp/post-update.diff"
# The extras stream is deliberately discarded here: install-only files are exactly what this
# check must NOT abort on, and they were already reported before the confirmation prompt.
append_install_diff "$fresh_upstream" "$skill_root" "$post_diff" "$update_tmp/post-update-extras.txt"
if [ -s "$post_diff" ]; then
  printf 'Update aborted: the extracted files differ from the reviewed upstream tree:\n' >&2
  cat "$post_diff" >&2
  exit 1
fi

update_verified=1  # install matches the reviewed upstream tree; cleanup's failure report stands down
printf 'Updated to v%s at %s\n' "$remote_version" "$skill_root"
