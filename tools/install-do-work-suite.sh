#!/usr/bin/env bash
# Install and reconcile the complete four-skill do-work suite in one recoverable transaction.
set -euo pipefail

upstream_url='https://github.com/knews2019/skill-do-work/archive/refs/heads/main.tar.gz'
manual_settings_instruction='MANUAL STEP: merge .claude/skills/do-work/hooks/hooks.json into .claude/settings.json; preserve every existing entry.'
cancel_exit_status="${DO_WORK_INSTALL_CANCEL_EXIT_STATUS:-0}"
case "$cancel_exit_status" in
  ''|*[!0-9]*) printf 'do-work suite install: invalid cancellation status override\n' >&2; exit 1 ;;
esac

fail() {
  printf 'do-work suite install: %s\n' "$*" >&2
  exit 1
}

print_bootstrap_command() {
  cat <<'BOOTSTRAP'
(
  set -e
  project_root="$(git rev-parse --show-toplevel 2>/dev/null)" || { printf 'do-work bootstrap: run this from inside the target Git repository\n' >&2; exit 1; }
  bootstrap_tmp="$(mktemp -d "${TMPDIR:-/tmp}/do-work-suite-bootstrap.XXXXXX")"
  trap 'rm -rf "$bootstrap_tmp"' EXIT
  archive_file="$bootstrap_tmp/do-work-suite.tar.gz"
  curl -fsSL -o "$archive_file.download" https://github.com/knews2019/skill-do-work/archive/refs/heads/main.tar.gz
  mv "$archive_file.download" "$archive_file"
  mkdir -p "$bootstrap_tmp/source"
  tar xzf "$archive_file" -C "$bootstrap_tmp/source" --strip-components=1
  bash "$bootstrap_tmp/source/tools/install-do-work-suite.sh" --project-root "$project_root" --archive "$archive_file"
)
BOOTSTRAP
}

if [ "$#" -eq 1 ] && [ "$1" = '--print-bootstrap-command' ]; then
  print_bootstrap_command
  exit 0
fi

project_root=''
supplied_archive=''
while [ "$#" -gt 0 ]; do
  case "$1" in
    --project-root)
      [ "$#" -ge 2 ] || fail 'usage: install-do-work-suite.sh --project-root <Git-root> [--archive <source.tar.gz>]'
      project_root="$2"
      shift 2
      ;;
    --archive)
      [ "$#" -ge 2 ] || fail 'usage: install-do-work-suite.sh --project-root <Git-root> [--archive <source.tar.gz>]'
      supplied_archive="$2"
      shift 2
      ;;
    *)
      fail 'usage: install-do-work-suite.sh --project-root <Git-root> [--archive <source.tar.gz>]'
      ;;
  esac
done

[ -n "$project_root" ] || fail '--project-root is required'
[ -d "$project_root" ] || fail "project root does not exist: $project_root"
project_root="$(cd "$project_root" && pwd -P)"
git_root="$(git -C "$project_root" rev-parse --show-toplevel 2>/dev/null || true)"
[ -n "$git_root" ] || fail '--project-root must be a Git repository so recovery is deterministic'
git_root="$(cd "$git_root" && pwd -P)"
[ "$git_root" = "$project_root" ] \
  || fail "--project-root must name the Git worktree root ($git_root)"

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd -P)"
manifest_validator="$script_dir/validate-suite-manifest.sh"
section_replacer="$script_dir/replace-text-section.sh"
[ -f "$manifest_validator" ] || fail 'validate-suite-manifest.sh is missing beside the installer'
[ -f "$section_replacer" ] || fail 'replace-text-section.sh is missing beside the installer'

install_tmp="$(mktemp -d "${TMPDIR:-/tmp}/do-work-suite-install.XXXXXX")"
source_root="$install_tmp/source"
backup_root="$install_tmp/originals"
mkdir -p "$source_root" "$backup_root"

module_sources=()
module_destinations=()
module_relatives=()
module_existed=()
write_started=''
install_verified=''
just_target=''
just_existed=''
settings_target="$project_root/.claude/settings.json"
settings_existed=''

recover_install() {
  local recovery_failed='' module_index destination_path

  set +e
  for module_index in "${!module_destinations[@]}"; do
    destination_path="${module_destinations[$module_index]}"
    rm -rf -- "$destination_path" || recovery_failed=1
    if [ "${module_existed[$module_index]}" = 1 ]; then
      mkdir -p "$destination_path" || recovery_failed=1
      cp -Rp "$backup_root/modules/$module_index/." "$destination_path/" || recovery_failed=1
    fi
  done

  if [ -n "$just_target" ]; then
    rm -rf -- "$just_target" || recovery_failed=1
    if [ "$just_existed" = 1 ]; then
      mkdir -p "$(dirname "$just_target")" || recovery_failed=1
      cp -p "$backup_root/justfile" "$just_target" || recovery_failed=1
    fi
  fi

  rm -rf -- "$settings_target" || recovery_failed=1
  if [ "$settings_existed" = 1 ]; then
    mkdir -p "$(dirname "$settings_target")" || recovery_failed=1
    cp -p "$backup_root/settings.json" "$settings_target" || recovery_failed=1
  fi
  set -e

  if [ -n "$recovery_failed" ]; then
    printf 'do-work suite install: automatic recovery was incomplete; inspect the four skill directories, %s, and %s\n' \
      "$just_target" "$settings_target" >&2
    return 1
  fi
  printf 'do-work suite install: restored every managed path to its exact pre-install state.\n' >&2
}

cleanup() {
  local exit_status=$?
  trap - EXIT
  trap '' HUP INT TERM
  if [ -n "$write_started" ] && [ -z "$install_verified" ]; then
    printf 'do-work suite install: installation did not complete; recovering managed paths.\n' >&2
    recover_install || exit_status=1
  fi
  rm -rf "$install_tmp"
  exit "$exit_status"
}
trap cleanup EXIT
trap 'exit 130' HUP INT TERM

archive_file=''
if [ -n "$supplied_archive" ]; then
  [ -f "$supplied_archive" ] && [ ! -L "$supplied_archive" ] \
    || fail "--archive must name a regular file: $supplied_archive"
  archive_directory="$(cd "$(dirname "$supplied_archive")" && pwd -P)"
  archive_file="$archive_directory/${supplied_archive##*/}"
else
  archive_file="$install_tmp/upstream.tar.gz"
  curl -fsSL -o "$archive_file.download" "$upstream_url" \
    || fail 'upstream archive download failed; no client files were changed'
  mv "$archive_file.download" "$archive_file"
fi

tar xzf "$archive_file" -C "$source_root" --strip-components=1 \
  || fail 'suite archive extraction failed; no client files were changed'
bash "$manifest_validator" --root "$source_root" \
  || fail 'suite manifest validation failed; no client files were changed'
suite_version="$(sed -n '1p' "$source_root/VERSION")"

while IFS=$'\t' read -r module_source module_destination; do
  [ "$module_source" = source ] && continue
  module_sources+=("$source_root/$module_source")
  module_destinations+=("$project_root/$module_destination")
  module_relatives+=("$module_destination")
done < "$source_root/suite/modules.tsv"
[ "${#module_sources[@]}" -eq 4 ] || fail 'validated suite did not produce exactly four install targets'

# The manifest validator constrains textual destinations. Resolve the nearest existing parent
# physically as well, so a project-local symlink cannot redirect a managed write outside Git.
for destination_path in "${module_destinations[@]}"; do
  [ ! -L "$destination_path" ] || fail "managed destination must not be a symlink: $destination_path"
  if [ -e "$destination_path" ] && [ ! -d "$destination_path" ]; then
    fail "managed destination must be a directory when it exists: $destination_path"
  fi
  destination_parent="$(dirname "$destination_path")"
  while [ ! -d "$destination_parent" ]; do
    next_parent="$(dirname "$destination_parent")"
    [ "$next_parent" != "$destination_parent" ] \
      || fail "cannot resolve managed destination parent: $destination_path"
    destination_parent="$next_parent"
  done
  destination_parent="$(cd "$destination_parent" && pwd -P)"
  case "$destination_parent" in
    "$project_root"|"$project_root"/*) ;;
    *) fail "managed destination resolves outside the project: $destination_path" ;;
  esac
done

for source_path in "${module_sources[@]}"; do
  [ -d "$source_path" ] && [ ! -L "$source_path" ] \
    || fail "managed module source is not a real directory: $source_path"
  if find "$source_path" -type l -print -quit | grep -q .; then
    fail "managed module source contains a symlink: $source_path"
  fi
done

board_template="$source_root/skills/do-work-board/justfile.template"
core_hooks="$source_root/skills/do-work/hooks/hooks.json"
[ -s "$board_template" ] && [ ! -L "$board_template" ] || fail 'board Justfile template is missing or unsafe'
[ -s "$core_hooks" ] && [ ! -L "$core_hooks" ] || fail 'core hook fragment is missing or unsafe'

for justfile_name in justfile Justfile .justfile; do
  candidate_path="$project_root/$justfile_name"
  if [ -e "$candidate_path" ] || [ -L "$candidate_path" ]; then
    [ -f "$candidate_path" ] && [ ! -L "$candidate_path" ] \
      || fail "Justfile target must be a regular file: $candidate_path"
    just_target="$candidate_path"
    just_existed=1
    break
  fi
done
if [ -z "$just_target" ]; then
  just_target="$project_root/justfile"
  just_existed=0
fi

managed_section="$install_tmp/do-work-recipes.just"
awk '
  $0 == "# >>> do-work:recipes >>>" { inside=1 }
  inside { print }
  $0 == "# <<< do-work:recipes <<<" { found_end=1; exit }
  END { if (!inside || !found_end) exit 1 }
' "$board_template" > "$managed_section" \
  || fail 'board template does not contain one complete managed recipe section'

just_candidate="$install_tmp/justfile.candidate"
if command -v python3 >/dev/null 2>&1; then
  if [ "$just_existed" = 1 ]; then
    cp -p "$just_target" "$just_candidate"
    bash "$section_replacer" --target "$just_candidate" --section-file "$managed_section" \
      || fail 'Justfile ownership validation failed; no client files were changed'
  else
    bash "$section_replacer" --target "$just_candidate" --section-file "$managed_section" \
      --template-file "$board_template" \
      || fail 'complete Justfile candidate creation failed; no client files were changed'
  fi
else
  [ "$just_existed" = 0 ] \
    || fail 'python3 is required to reconcile an existing Justfile safely; no client files were changed'
  cp -p "$board_template" "$just_candidate"
fi

[ "$(grep -c '^# >>> do-work:recipes >>>$' "$just_candidate")" -eq 1 ] \
  && [ "$(grep -c '^# <<< do-work:recipes <<<$' "$just_candidate")" -eq 1 ] \
  || fail 'Justfile candidate has invalid managed markers'
for recipe_name in run-kanban run-kanban-cli kanban-static kanban-summary run-do-work-update; do
  grep -Eq "^${recipe_name}[ :].*:$|^${recipe_name}:$" "$just_candidate" \
    || fail "Justfile candidate is missing $recipe_name"
done
if command -v just >/dev/null 2>&1; then
  just --justfile "$just_candidate" --list >/dev/null 2>&1 \
    || fail 'Justfile candidate does not parse; no client files were changed'
fi

settings_tool='manual'
settings_candidate=''
if command -v jq >/dev/null 2>&1; then
  settings_tool='jq'
elif command -v python3 >/dev/null 2>&1; then
  settings_tool='python3'
fi

if [ -e "$settings_target" ] || [ -L "$settings_target" ]; then
  [ -f "$settings_target" ] && [ ! -L "$settings_target" ] \
    || fail "Claude settings must be a regular file: $settings_target"
  settings_existed=1
else
  settings_existed=0
fi

if [ "$settings_tool" != manual ]; then
  settings_input="$install_tmp/settings.input.json"
  settings_candidate="$install_tmp/settings.candidate.json"
  if [ "$settings_existed" = 1 ]; then
    cp -p "$settings_target" "$settings_input"
    settings_mode="$(stat -f '%Lp' "$settings_target" 2>/dev/null || stat -c '%a' "$settings_target" 2>/dev/null || true)"
  else
    printf '{}\n' > "$settings_input"
    settings_mode=644
  fi

  if [ "$settings_tool" = jq ]; then
    jq --slurpfile fragment "$core_hooks" '
      def append_unique($base; $extra):
        reduce $extra[] as $item ($base; if index($item) == null then . + [$item] else . end);
      if .hooks == null then .hooks = {}
        elif (.hooks | type) != "object" then error("settings hooks must be an object")
        else . end
      | if .hooks.Stop == null then .
        elif (.hooks.Stop | type) != "array" then error("settings Stop hook event must be an array")
        else
          .hooks.Stop |= map(
            if type == "object" and (.hooks | type) == "array" then
              .hooks |= map(select(
                ((.command | type) == "string"
                  and (.command | contains(".claude/skills/do-work/hooks/pipeline-guard.sh")))
                | not
              ))
            else . end
          )
          | .hooks.Stop |= map(select(
              (type != "object") or (.hooks | type) != "array" or (.hooks | length) > 0
            ))
          | if (.hooks.Stop | length) == 0 then del(.hooks.Stop) else . end
        end
      | reduce ($fragment[0].hooks | keys[]) as $event (.;
          if .hooks[$event] == null then .hooks[$event] = []
          elif (.hooks[$event] | type) != "array" then error("settings hook event must be an array")
          else . end
          | .hooks[$event] = append_unique(.hooks[$event]; $fragment[0].hooks[$event])
        )
    ' "$settings_input" > "$settings_candidate" \
      || fail 'Claude settings are invalid or cannot accept composed core hooks; no client files were changed'
    jq -e . "$settings_candidate" >/dev/null \
      || fail 'composed Claude settings failed JSON validation; no client files were changed'
  else
    python3 - "$settings_input" "$core_hooks" "$settings_candidate" <<'PY' \
      || fail 'Claude settings are invalid or cannot accept composed core hooks; no client files were changed'
import json
import pathlib
import sys

source_path, fragment_path, output_path = map(pathlib.Path, sys.argv[1:])
with source_path.open() as handle:
    settings = json.load(handle)
with fragment_path.open() as handle:
    fragment = json.load(handle)
if not isinstance(settings, dict):
    raise TypeError("settings root must be an object")

hooks = settings.setdefault("hooks", {})
if not isinstance(hooks, dict):
    raise TypeError("settings hooks must be an object")
stop_entries = hooks.get("Stop")
if stop_entries is not None:
    if not isinstance(stop_entries, list):
        raise TypeError("settings Stop hook event must be an array")
    retained_stop_entries = []
    for entry in stop_entries:
        if isinstance(entry, dict) and isinstance(entry.get("hooks"), list):
            retained_hooks = [
                hook
                for hook in entry["hooks"]
                if not (
                    isinstance(hook, dict)
                    and isinstance(hook.get("command"), str)
                    and ".claude/skills/do-work/hooks/pipeline-guard.sh" in hook["command"]
                )
            ]
            if retained_hooks:
                entry = dict(entry)
                entry["hooks"] = retained_hooks
                retained_stop_entries.append(entry)
        else:
            retained_stop_entries.append(entry)
    if retained_stop_entries:
        hooks["Stop"] = retained_stop_entries
    else:
        hooks.pop("Stop")
for event, entries in fragment.get("hooks", {}).items():
    installed_entries = hooks.setdefault(event, [])
    if not isinstance(installed_entries, list) or not isinstance(entries, list):
        raise TypeError("settings hook event must be an array")
    for entry in entries:
        if entry not in installed_entries:
            installed_entries.append(entry)
with output_path.open("w") as handle:
    json.dump(settings, handle, indent=2)
    handle.write("\n")
PY
    python3 -m json.tool "$settings_candidate" >/dev/null \
      || fail 'composed Claude settings failed JSON validation; no client files were changed'
  fi
  chmod "$settings_mode" "$settings_candidate"
  grep -q 'do-work/hooks/session-start.sh' "$settings_candidate" \
    || fail 'composed settings omitted the core SessionStart hook'
  if grep -q 'do-work/hooks/pipeline-guard.sh' "$settings_candidate"; then
    fail 'composed settings retained the retired pipeline Stop hook'
  fi
fi

printf 'Ready to install do-work suite v%s into %s:\n' "$suite_version" "$project_root"
printf '  %s\n' "${module_relatives[@]}"
printf '  Justfile: %s\n' "${just_target#"$project_root"/}"
printf '  settings reconciler: %s\n' "$settings_tool"

review_diff="$install_tmp/install.diff"
: > "$review_diff"
for module_index in "${!module_destinations[@]}"; do
  printf '\n--- managed destination: %s ---\n' "${module_relatives[$module_index]}" >> "$review_diff"
  diff_status=0
  diff -ruN "${module_destinations[$module_index]}" "${module_sources[$module_index]}" \
    >> "$review_diff" 2>&1 || diff_status=$?
  [ "$diff_status" -le 1 ] \
    || fail "could not compare managed destination ${module_relatives[$module_index]}"
done
printf '\n--- managed configuration: %s ---\n' "${just_target#"$project_root"/}" >> "$review_diff"
diff_status=0
if [ "$just_existed" = 1 ]; then
  diff -u "$just_target" "$just_candidate" >> "$review_diff" 2>&1 || diff_status=$?
else
  diff -u /dev/null "$just_candidate" >> "$review_diff" 2>&1 || diff_status=$?
fi
[ "$diff_status" -le 1 ] || fail 'could not compare the managed Justfile candidate'
printf '\n--- managed configuration: .claude/settings.json ---\n' >> "$review_diff"
if [ "$settings_tool" = manual ]; then
  printf '%s\n' "$manual_settings_instruction" >> "$review_diff"
else
  diff_status=0
  if [ "$settings_existed" = 1 ]; then
    diff -u "$settings_target" "$settings_candidate" >> "$review_diff" 2>&1 || diff_status=$?
  else
    diff -u /dev/null "$settings_candidate" >> "$review_diff" 2>&1 || diff_status=$?
  fi
  [ "$diff_status" -le 1 ] || fail 'could not compare the Claude settings candidate'
fi
printf 'Reviewing the complete managed install before overwrite:\n'
cat "$review_diff"

dirty_managed="$(git -C "$project_root" status --porcelain -- \
  "${module_relatives[@]}" "${just_target#"$project_root"/}" '.claude/settings.json')"
if [ -n "$dirty_managed" ]; then
  printf 'Managed install paths have uncommitted changes. Continuing discards those changes in managed modules and replaces only the owned configuration bytes:\n%s\n' \
    "$dirty_managed" >&2
fi
printf 'Install this complete four-skill suite? [y/N] '
read -r confirmation || confirmation=''
case "$confirmation" in
  y|Y|yes|YES) ;;
  *) printf 'Installation cancelled; no files were changed.\n'; exit "$cancel_exit_status" ;;
esac

mkdir -p "$backup_root/modules"
for module_index in "${!module_destinations[@]}"; do
  destination_path="${module_destinations[$module_index]}"
  if [ -e "$destination_path" ]; then
    module_existed+=(1)
    mkdir -p "$backup_root/modules/$module_index"
    cp -Rp "$destination_path/." "$backup_root/modules/$module_index/"
  else
    module_existed+=(0)
  fi
done
if [ "$just_existed" = 1 ]; then
  cp -p "$just_target" "$backup_root/justfile"
fi
if [ "$settings_existed" = 1 ]; then
  cp -p "$settings_target" "$backup_root/settings.json"
fi

# Confirmation authorizes discarding dirty managed module content. Clear only module paths
# from the index so a previously staged customization cannot survive beneath installed bytes;
# Just/settings remain project configuration and are never wholesale-reset in the index.
for relative_path in "${module_relatives[@]}"; do
  if [ -n "$(git -C "$project_root" ls-files -- "$relative_path")" ]; then
    git -C "$project_root" restore --staged -- "$relative_path"
  fi
done

write_started=1
for module_index in "${!module_destinations[@]}"; do
  source_path="${module_sources[$module_index]}"
  destination_path="${module_destinations[$module_index]}"
  rm -rf -- "$destination_path"
  mkdir -p "$destination_path"
  cp -Rp "$source_path/." "$destination_path/"
done

mkdir -p "$(dirname "$just_target")"
# The installed helper produced and validated this complete candidate before module writes.
# Replace from that candidate so this transaction never executes a helper it just downloaded.
just_temporary="$(mktemp "$(dirname "$just_target")/.do-work-just.install.XXXXXX")"
cp -p "$just_candidate" "$just_temporary"
mv "$just_temporary" "$just_target"

if [ "$settings_tool" != manual ]; then
  mkdir -p "$(dirname "$settings_target")"
  settings_temporary="$(mktemp "$(dirname "$settings_target")/.settings.json.install.XXXXXX")"
  cp -p "$settings_candidate" "$settings_temporary"
  mv "$settings_temporary" "$settings_target"
fi

for module_index in "${!module_destinations[@]}"; do
  diff -qr "${module_sources[$module_index]}" "${module_destinations[$module_index]}" >/dev/null \
    || fail "installed bytes do not match ${module_relatives[$module_index]}"
done
cmp -s "$just_candidate" "$just_target" || fail 'installed Justfile does not match its validated candidate'
if command -v just >/dev/null 2>&1; then
  just --justfile "$just_target" --list >/dev/null 2>&1 \
    || fail 'installed Justfile failed post-write validation'
fi
if [ "$settings_tool" != manual ]; then
  cmp -s "$settings_candidate" "$settings_target" \
    || fail 'installed Claude settings do not match the validated candidate'
  if [ "$settings_tool" = jq ]; then
    jq -e . "$settings_target" >/dev/null || fail 'installed Claude settings failed post-write validation'
  else
    python3 -m json.tool "$settings_target" >/dev/null || fail 'installed Claude settings failed post-write validation'
  fi
fi
installed_version="$(sed -n '1p' "$project_root/.claude/skills/do-work/VERSION")"
[ "$installed_version" = "$suite_version" ] \
  || fail "installed version mismatch (expected $suite_version, found ${installed_version:-unknown})"

install_verified=1
printf 'Installed do-work suite v%s with four verified modules.\n' "$suite_version"
if [ "$settings_tool" = manual ]; then
  printf '%s\n' "$manual_settings_instruction"
fi
