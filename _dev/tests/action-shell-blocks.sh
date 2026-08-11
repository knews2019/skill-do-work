#!/usr/bin/env bash
# Lint every shipped Markdown shell fence and shipped shell source with attributable diagnostics.
set -uo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
skills_root="$repo_root/skills"
temporary_root="$(mktemp -d)" || {
  printf 'FAIL: could not allocate temporary directory for shell-block lint.\n' >&2
  exit 1
}
trap 'rm -rf "$temporary_root"' EXIT

failure_count=0
shell_block_count=0
shell_source_count=0
shellcheck_available=false
if command -v shellcheck >/dev/null 2>&1; then
  shellcheck_available=true
fi

source_display_path() {
  local source_path="$1"

  if [[ "$source_path" == "$repo_root/"* ]]; then
    printf '%s\n' "${source_path:${#repo_root}+1}"
  else
    printf '%s\n' "$source_path"
  fi
}

lint_shell_source() {
  local source_path="$1"
  local source_start_line="$2"
  local lint_path="$3"
  local source_kind="$4"
  local bash_diagnostic=''
  local bash_diagnostic_line=''
  local bash_source_line=''
  local shellcheck_diagnostics=''
  local shellcheck_status=0
  local diagnostic_text=''
  local diagnostic_source_line=''
  local diagnostic_column=''
  local diagnostic_message=''
  local -a shellcheck_arguments=(--format=gcc --shell=bash --severity=warning)

  if ! bash_diagnostic="$(bash -n "$lint_path" 2>&1)"; then
    bash_diagnostic_line="$(printf '%s\n' "$bash_diagnostic" \
      | sed -nE 's/.*: line ([0-9]+):.*/\1/p' \
      | sed -n '1p')"
    if [[ -z "$bash_diagnostic_line" ]]; then
      bash_diagnostic_line=1
    fi
    bash_source_line=$((source_start_line + bash_diagnostic_line - 1))
    bash_diagnostic="$(printf '%s' "$bash_diagnostic" \
      | tr '\n' ' ' \
      | sed 's/[[:space:]][[:space:]]*/ /g')"
    printf 'FAIL: %s:%s: bash -n: %s\n' \
      "$source_path" "$bash_source_line" "$bash_diagnostic" >&2
    failure_count=$((failure_count + 1))
    return
  fi

  if [[ "$shellcheck_available" != true ]]; then
    return
  fi

  if [[ "$source_kind" == fence ]]; then
    shellcheck_arguments+=("--exclude=SC2034,SC2154")
  fi
  shellcheck_diagnostics="$(shellcheck "${shellcheck_arguments[@]}" "$lint_path" 2>&1)" \
    || shellcheck_status=$?
  if [[ "$shellcheck_status" -eq 0 ]]; then
    return
  fi

  while IFS= read -r diagnostic_text; do
    [[ -n "$diagnostic_text" ]] || continue
    if [[ "$diagnostic_text" =~ ^.*:([0-9]+):([0-9]+):[[:space:]](.*)$ ]]; then
      diagnostic_source_line=$((source_start_line + BASH_REMATCH[1] - 1))
      diagnostic_column="${BASH_REMATCH[2]}"
      diagnostic_message="${BASH_REMATCH[3]}"
      printf 'FAIL: %s:%s:%s: shellcheck: %s\n' \
        "$source_path" "$diagnostic_source_line" "$diagnostic_column" "$diagnostic_message" >&2
    else
      printf 'FAIL: %s:%s: shellcheck: %s\n' \
        "$source_path" "$source_start_line" "$diagnostic_text" >&2
    fi
    failure_count=$((failure_count + 1))
  done <<< "$shellcheck_diagnostics"
}

scan_markdown_file() {
  local markdown_path="$1"
  local source_path=''
  local current_line=0
  local source_text=''
  local shell_fence_pattern='^[[:blank:]]{0,3}```(bash|sh)[[:space:]]*$'
  local closing_fence_pattern='^[[:blank:]]{0,3}```[[:space:]]*$'
  local in_shell_block=false
  local block_start_line=0
  local raw_block_path=''
  local lint_block_path=''

  source_path="$(source_display_path "$markdown_path")"
  while IFS= read -r source_text || [[ -n "$source_text" ]]; do
    current_line=$((current_line + 1))
    if [[ "$in_shell_block" != true ]]; then
      if [[ "$source_text" =~ $shell_fence_pattern ]]; then
        shell_block_count=$((shell_block_count + 1))
        block_start_line=$((current_line + 1))
        raw_block_path="$temporary_root/block-$shell_block_count.raw"
        lint_block_path="$temporary_root/block-$shell_block_count.sh"
        : > "$raw_block_path"
        in_shell_block=true
      fi
      continue
    fi

    if [[ "$source_text" =~ $closing_fence_pattern ]]; then
      # Neutralize prose placeholders without changing the block's line count. Do not skip the block.
      sed -E 's|<[[:alnum:]_][[:alnum:]_.: /-]*>|placeholder|g' \
        "$raw_block_path" > "$lint_block_path"
      lint_shell_source "$source_path" "$block_start_line" "$lint_block_path" fence
      rm -f "$raw_block_path" "$lint_block_path"
      in_shell_block=false
      continue
    fi

    printf '%s\n' "$source_text" >> "$raw_block_path"
  done < "$markdown_path"

  if [[ "$in_shell_block" == true ]]; then
    printf 'FAIL: %s:%s: unterminated fenced shell block.\n' \
      "$source_path" "$((block_start_line - 1))" >&2
    failure_count=$((failure_count + 1))
    rm -f "$raw_block_path" "$lint_block_path"
  fi
}

run_self_test() {
  local fixture_path="$temporary_root/broken-shell-block.md"
  local fixture_diagnostics="$temporary_root/self-test-diagnostics.txt"

  {
    printf '# Broken shell fixture\n\n'
    printf '   ```bash\n'
    printf '   if true; then\n'
    printf '     printf "missing fi\\n"\n'
    printf '   ```\n'
  } > "$fixture_path"

  failure_count=0
  shell_block_count=0
  scan_markdown_file "$fixture_path" 2> "$fixture_diagnostics"
  if [[ "$failure_count" -eq 0 ]]; then
    printf 'FAIL: shell-block lint self-test accepted a malformed fenced block.\n' >&2
    return 1
  fi
  if ! grep -Eq 'broken-shell-block\.md:[0-9]+: bash -n:' "$fixture_diagnostics"; then
    printf 'FAIL: shell-block lint self-test did not report fixture path, line, and diagnostic.\n' >&2
    sed -n '1,20p' "$fixture_diagnostics" >&2
    return 1
  fi

  printf 'Shell-block lint self-test passed.\n'
}

case "${1:-}" in
  --self-test)
    run_self_test
    exit $?
    ;;
  '') ;;
  *)
    printf 'Usage: %s [--self-test]\n' "$0" >&2
    exit 2
    ;;
esac

if [[ ! -d "$skills_root" ]]; then
  printf 'FAIL: shipped skills root is missing: %s\n' "$skills_root" >&2
  exit 1
fi

while IFS= read -r markdown_path; do
  scan_markdown_file "$markdown_path"
done < <(find "$skills_root" -type f -name '*.md' -print | LC_ALL=C sort)

while IFS= read -r shell_path; do
  shell_source_count=$((shell_source_count + 1))
  lint_shell_source "$(source_display_path "$shell_path")" 1 "$shell_path" source
done < <(find "$skills_root" -type f -name '*.sh' -print | LC_ALL=C sort)

if [[ "$shell_block_count" -eq 0 ]]; then
  printf 'FAIL: no fenced bash/sh blocks found under %s; extractor or scan root is broken.\n' \
    "$skills_root" >&2
  failure_count=$((failure_count + 1))
fi

if [[ "$failure_count" -gt 0 ]]; then
  printf 'Shell-block lint found %s finding(s) across %s fenced blocks and %s shipped shell files.\n' \
    "$failure_count" "$shell_block_count" "$shell_source_count" >&2
  exit 1
fi

if [[ "$shellcheck_available" == true ]]; then
  printf 'Shell-block lint passed: %s fenced blocks and %s shipped shell files; ShellCheck enabled.\n' \
    "$shell_block_count" "$shell_source_count"
else
  printf 'NOTE: shellcheck unavailable — ran bash -n only.\n'
  printf 'Shell-block lint passed: %s fenced blocks and %s shipped shell files.\n' \
    "$shell_block_count" "$shell_source_count"
fi
