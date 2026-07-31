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

# What the automatic rollback restores. Beyond the shipped set: the project justfile that a
# failure *between* the extraction and its restore could have left clobbered (root fallback),
# and the stale vendored maintainer docs the update deletes on purpose (nested install).
restore_paths=("${shipped_paths[@]}")
if [ -n "$root_fallback_install" ]; then
  # Order is load-bearing on a case-insensitive filesystem: clearing the extracted lowercase
  # `justfile` must come BEFORE restoring a `Justfile`, or the clear would delete the file
  # just restored. Appending the project's own name is skipped when it IS `justfile`.
  restore_paths+=(justfile)
  if [ -n "$project_justfile_name" ] && [ "$project_justfile_name" != 'justfile' ]; then
    restore_paths+=("$project_justfile_name")
  fi
else
  restore_paths+=(CLAUDE.md AGENTS.md)
fi

[ -s "$skill_root/SKILL.md" ] || fail "SKILL.md is missing at $skill_root"
[ -f "$skill_root/actions/version.md" ] || fail "actions/version.md is missing at $skill_root"

local_version="$(sed -n 's/^\*\*Current version\*\*: *\([0-9][0-9.]*\).*/\1/p' "$skill_root/actions/version.md" | head -n 1)"
[ -n "$local_version" ] || fail 'could not read the local version'

update_tmp="$(mktemp -d "${TMPDIR:-/tmp}/do-work-update.XXXXXX")"
backup_path=''
backup_ready=''
update_verified=''

restore_from_backup() {
  # Put every restore path back exactly as the rollback copy has it: remove, then copy
  # back. Removing first is what makes this an undo rather than a merge — a path the
  # update ADDED (absent from the backup) is cleared instead of surviving the rollback.
  # Only the enumerated paths are touched; never `rm -rf "$skill_root"`, which in the
  # root-fallback install is the entire project — git history, the do-work/ queue, and all.
  local restore_path
  local restore_failed=''
  for restore_path in "${restore_paths[@]}"; do
    [ -e "$backup_path/$restore_path" ] || [ -e "$skill_root/$restore_path" ] || continue
    # `:?` on both halves so a would-be-empty expansion aborts instead of resolving to `/`.
    rm -rf "${skill_root:?}/${restore_path:?}" || { restore_failed=1; continue; }
    [ -e "$backup_path/$restore_path" ] || continue
    cp -R "$backup_path/$restore_path" "$skill_root/$restore_path" || restore_failed=1
  done
  [ -z "$restore_failed" ]
}

cleanup() {
  rm -rf "$update_tmp"
  # Transactional outcome: either the update completes and verifies, or the install goes
  # back to its pre-update state. `update_verified` is set on the fully-verified path only,
  # so EVERY earlier exit inside the destructive region lands here and undoes the writes —
  # a mid-extract ENOSPC (the cp -R backup just doubled the install's on-disk size), a
  # failed post-update version check, or `set -e` firing anywhere in between. The rollback
  # copy is still kept and reported: it is the audit trail and the manual fallback for the
  # one case this cannot handle, a restore that itself fails partway.
  if [ -n "$backup_ready" ] && [ -z "$update_verified" ]; then
    printf 'Update did not complete — restoring the install from the rollback copy…\n' >&2
    if restore_from_backup; then
      printf 'Restored to v%s; nothing was updated. Rollback copy kept at: %s\n' "$local_version" "$backup_path" >&2
    else
      printf 'AUTOMATIC RESTORE FAILED — the install is partially modified.\nRestore it by hand from the rollback copy: %s\n' "$backup_path" >&2
    fi
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
fi

diff_file="$update_tmp/install.diff"
extras_file="$update_tmp/install-extras.txt"
append_install_diff "$fresh_upstream" "$skill_root" "$diff_file" "$extras_file"

if [ -n "$dirty_files" ]; then
  printf 'Shipped skill files have uncommitted changes:\n%s\n' "$dirty_files" >&2
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

printf 'Continue with the update? This creates a rollback copy first. [y/N] '
# EOF / non-interactive stdin (piped, CI, </dev/null) makes `read` return non-zero; under
# `set -e` a bare `read` would abort here before the case below can default to No. Treat
# EOF as an explicit No so the cancel path runs and the script exits 0.
read -r confirmation || confirmation=''
case "$confirmation" in
  y|Y|yes|YES) ;;
  *) printf 'Update cancelled; no files were changed.\n'; exit 0 ;;
esac

backup_path="$skill_root.preupdate-$(date -u +%Y%m%dT%H%M%SZ).bak"
cp -R "$skill_root" "$backup_path"
backup_ready=1  # rollback copy exists; every exit below rolls back unless the update verifies

preserved_justfile=''
if [ -n "$project_justfile_name" ]; then
  preserved_justfile="$update_tmp/project-justfile"
  cp -p "$project_root/$project_justfile_name" "$preserved_justfile" \
    || fail "could not preserve the project justfile at $project_root/$project_justfile_name"
fi

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

update_verified=1  # install matches the reviewed upstream tree; cleanup's rollback stands down
printf 'Updated to v%s at %s\n' "$remote_version" "$skill_root"
printf 'Rollback copy: %s\n' "$backup_path"
