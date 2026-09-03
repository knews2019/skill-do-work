#!/usr/bin/env bash
# Byte-preserving managed-section replacement behavior.
set -euo pipefail
repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd)"
fail_count=0
# Managed Just sections are a byte-preserving ownership boundary, not a prose convention.
# Exercise the real utility across replacement, append, creation, malformed
# markers, filename variants, spaces, modes, idempotence, and Just parsing.
replace_section_tool="$repo_root/tools/replace-text-section.sh"
if [ ! -x "$replace_section_tool" ]; then
  printf 'FAIL: tools/replace-text-section.sh is missing or not executable — managed recipe ownership has no implementation.\n' >&2
  fail_count=$((fail_count + 1))
else
  section_workdir="$(mktemp -d)"
  section_file="$section_workdir/managed-section.just"
  template_file="$section_workdir/complete-template.just"
  printf '# >>> do-work:recipes >>>\nmanaged-probe:\n    echo managed\n# <<< do-work:recipes <<<\n' > "$section_file"
  printf 'set shell := ["bash", "-cu"]\n\n# >>> do-work:recipes >>>\nmanaged-probe:\n    echo managed\n# <<< do-work:recipes <<<\n' > "$template_file"
  chmod 750 "$template_file"

  byte_target="$section_workdir/project with spaces/Justfile"
  mkdir -p "$(dirname "$byte_target")"
  printf 'prefix\000byte\n# >>> do-work:recipes >>>\nold:\n    echo old\n# <<< do-work:recipes <<<\nsuffix\n' > "$byte_target"
  chmod 640 "$byte_target"
  expected_target="$section_workdir/expected-byte-target"
  printf 'prefix\000byte\n# >>> do-work:recipes >>>\nmanaged-probe:\n    echo managed\n# <<< do-work:recipes <<<\nsuffix\n' > "$expected_target"
  if ! "$replace_section_tool" --target "$byte_target" --section-file "$section_file"; then
    printf 'FAIL: replace-text-section could not replace one valid managed section.\n' >&2
    fail_count=$((fail_count + 1))
  elif ! cmp -s "$byte_target" "$expected_target"; then
    printf 'FAIL: replace-text-section changed bytes outside the managed section or wrote the wrong replacement.\n' >&2
    fail_count=$((fail_count + 1))
  fi
  target_mode="$(stat -c '%a' "$byte_target" 2>/dev/null || stat -f '%Lp' "$byte_target" 2>/dev/null || true)"
  if [ "$target_mode" != 640 ]; then
    printf 'FAIL: replace-text-section changed the existing target mode (got %s, want 640).\n' "$target_mode" >&2
    fail_count=$((fail_count + 1))
  fi
  cp "$byte_target" "$section_workdir/idempotent-snapshot"
  if ! "$replace_section_tool" --target "$byte_target" --section-file "$section_file" \
     || ! cmp -s "$byte_target" "$section_workdir/idempotent-snapshot"; then
    printf 'FAIL: replace-text-section is not byte-idempotent on repeated execution.\n' >&2
    fail_count=$((fail_count + 1))
  fi

  absent_target="$section_workdir/absent project/.justfile"
  mkdir -p "$(dirname "$absent_target")"
  if ! "$replace_section_tool" --target "$absent_target" --section-file "$section_file" --template-file "$template_file" \
     || ! cmp -s "$absent_target" "$template_file"; then
    printf 'FAIL: replace-text-section did not create an absent target from the complete supplied template.\n' >&2
    fail_count=$((fail_count + 1))
  fi
  absent_mode="$(stat -c '%a' "$absent_target" 2>/dev/null || stat -f '%Lp' "$absent_target" 2>/dev/null || true)"
  if [ "$absent_mode" != 750 ]; then
    printf 'FAIL: replace-text-section did not preserve the complete template mode on create (got %s, want 750).\n' "$absent_mode" >&2
    fail_count=$((fail_count + 1))
  fi

  for justfile_variant in justfile Justfile .justfile; do
    variant_target="$section_workdir/variants/$justfile_variant"
    mkdir -p "$(dirname "$variant_target")"
    printf 'custom-%s:\n    echo untouched\n' "$justfile_variant" > "$variant_target"
    if ! "$replace_section_tool" --target "$variant_target" --section-file "$section_file"; then
      printf 'FAIL: replace-text-section could not append to marker-free %s.\n' "$justfile_variant" >&2
      fail_count=$((fail_count + 1))
    elif [ "$(grep -c '^# >>> do-work:recipes >>>$' "$variant_target")" -ne 1 ] \
      || ! grep -q "^custom-$justfile_variant:" "$variant_target"; then
      printf 'FAIL: replace-text-section duplicated ownership or changed custom content in %s.\n' "$justfile_variant" >&2
      fail_count=$((fail_count + 1))
    fi
  done

  reserved_section_file="$section_workdir/reserved-section.just"
  awk '
    $0 == "# >>> do-work:recipes >>>" { inside=1 }
    inside { print }
    $0 == "# <<< do-work:recipes <<<" { found_end=1; exit }
    END { if (!inside || !found_end) exit 1 }
  ' "$repo_root/skills/do-work-board/justfile.template" > "$reserved_section_file"

  collision_index=0
  while IFS= read -r reserved_recipe_name; do
    collision_index=$((collision_index + 1))
    collision_target="$section_workdir/collision-$collision_index.just"
    printf '%s:\n    echo collision\n' "$reserved_recipe_name" > "$collision_target"
    cp "$collision_target" "$collision_target.before"
    collision_output="$section_workdir/collision-$collision_index.out"
    if "$replace_section_tool" --target "$collision_target" --section-file "$reserved_section_file" \
      --reject-recipe-collisions >"$collision_output" 2>&1; then
      printf 'FAIL: replace-text-section accepted external reserved recipe or alias %s.\n' "$reserved_recipe_name" >&2
      fail_count=$((fail_count + 1))
    elif ! grep -Fq "reserved Just recipe or alias outside managed section: $reserved_recipe_name" "$collision_output"; then
      printf 'FAIL: replace-text-section collision error did not name %s.\n' "$reserved_recipe_name" >&2
      fail_count=$((fail_count + 1))
    elif ! cmp -s "$collision_target" "$collision_target.before"; then
      printf 'FAIL: replace-text-section changed the target after rejecting collision %s.\n' "$reserved_recipe_name" >&2
      fail_count=$((fail_count + 1))
    fi
  done < <(just --justfile "$reserved_section_file" --summary | tr ' ' '\n' | LC_ALL=C sort)

  multiline_header_index=0
  for multiline_header_kind in \
    ordinary-single ordinary-double ordinary-backtick \
    triple-single triple-double triple-backtick; do
    multiline_header_index=$((multiline_header_index + 1))
    multiline_header_payload='payload'
    case "$multiline_header_kind" in
      ordinary-single) multiline_header_delimiter="'" ;;
      ordinary-double) multiline_header_delimiter='"' ;;
      ordinary-backtick) multiline_header_delimiter='`' ;;
      triple-single)
        multiline_header_delimiter="'''"
        multiline_header_payload="payload's"
        ;;
      triple-double)
        multiline_header_delimiter='"""'
        multiline_header_payload='payload"s'
        ;;
      triple-backtick)
        multiline_header_delimiter='```'
        multiline_header_payload='payload`s'
        ;;
    esac
    multiline_header_target="$section_workdir/multiline-header-$multiline_header_index.just"
    printf 'run-kanban value=%s\n%s\n%s:\n    echo collision\n' \
      "$multiline_header_delimiter" "$multiline_header_payload" "$multiline_header_delimiter" \
      > "$multiline_header_target"
    if command -v just >/dev/null 2>&1 \
      && ! just --justfile "$multiline_header_target" --list >/dev/null 2>&1; then
      printf 'FAIL: %s multiline-default recipe-header fixture is not valid Just syntax.\n' \
        "$multiline_header_kind" >&2
      fail_count=$((fail_count + 1))
      continue
    fi
    cp "$multiline_header_target" "$multiline_header_target.before"
    multiline_header_output="$section_workdir/multiline-header-$multiline_header_index.out"
    if "$replace_section_tool" --target "$multiline_header_target" --section-file "$reserved_section_file" \
      --reject-recipe-collisions >"$multiline_header_output" 2>&1; then
      printf 'FAIL: replace-text-section accepted reserved recipe with a %s multiline default.\n' \
        "$multiline_header_kind" >&2
      fail_count=$((fail_count + 1))
    elif ! grep -Fq \
      'reserved Just recipe or alias outside managed section: run-kanban' \
      "$multiline_header_output"; then
      printf 'FAIL: replace-text-section did not name the reserved recipe with a %s multiline default.\n' \
        "$multiline_header_kind" >&2
      fail_count=$((fail_count + 1))
    elif ! cmp -s "$multiline_header_target" "$multiline_header_target.before"; then
      printf 'FAIL: replace-text-section changed the target after rejecting a %s multiline-default collision.\n' \
        "$multiline_header_kind" >&2
      fail_count=$((fail_count + 1))
    fi
  done

  raw_multiline_target="$section_workdir/raw-multiline-string.just"
  {
    printf '%s\n' \
      'custom-before:' \
      '    echo before' \
      "raw_value := '''" \
      'run-kanban:' \
      'alias kanban-summary := ignored' \
      "\\'''" \
      'custom-after:' \
      '    echo after'
  } > "$raw_multiline_target"
  if command -v just >/dev/null 2>&1 \
    && ! just --justfile "$raw_multiline_target" --list >/dev/null 2>&1; then
    printf 'FAIL: raw multiline-string collision fixture is not valid Just syntax.\n' >&2
    fail_count=$((fail_count + 1))
  elif ! "$replace_section_tool" --target "$raw_multiline_target" --section-file "$reserved_section_file" \
    --reject-recipe-collisions; then
    printf 'FAIL: replace-text-section treated reserved-looking raw multiline-string content as a collision.\n' >&2
    fail_count=$((fail_count + 1))
  elif command -v just >/dev/null 2>&1 \
    && ! just --justfile "$raw_multiline_target" --list >/dev/null 2>&1; then
    printf 'FAIL: replace-text-section produced invalid Just after accepting a raw multiline string.\n' >&2
    fail_count=$((fail_count + 1))
  fi

  cooked_multiline_target="$section_workdir/cooked-multiline-string-crlf.just"
  {
    printf '%s\r\n' \
      'custom-before:' \
      '    echo before' \
      'cooked_value := """' \
      'run-kanban-cli:' \
      'odd escaped delimiter remains: \"""' \
      'alias kanban-summary := ignored' \
      'even escaped delimiter closes: \\"""' \
      'joined_value := """closed""" + """open' \
      'run-do-work-update:' \
      '"""' \
      'custom-after:' \
      '    echo after'
  } > "$cooked_multiline_target"
  if command -v just >/dev/null 2>&1 \
    && ! just --justfile "$cooked_multiline_target" --list >/dev/null 2>&1; then
    printf 'FAIL: cooked multiline-string collision fixture is not valid Just syntax.\n' >&2
    fail_count=$((fail_count + 1))
  elif ! "$replace_section_tool" --target "$cooked_multiline_target" --section-file "$reserved_section_file" \
    --reject-recipe-collisions; then
    printf 'FAIL: replace-text-section treated reserved-looking cooked multiline-string content as a collision.\n' >&2
    fail_count=$((fail_count + 1))
  elif command -v just >/dev/null 2>&1 \
    && ! just --justfile "$cooked_multiline_target" --list >/dev/null 2>&1; then
    printf 'FAIL: replace-text-section produced invalid Just after accepting a cooked multiline string.\n' >&2
    fail_count=$((fail_count + 1))
  fi

  ordinary_and_command_multiline_target="$section_workdir/ordinary-and-command-multiline-literals.just"
  {
    printf '%s\n' \
      'custom-before:' \
      '    echo before' \
      "raw_value := '" \
      'run-kanban:' \
      'alias kanban-summary := ignored' \
      "raw backslash is literal: \\'" \
      'cooked_value := "' \
      'run-kanban-cli:' \
      'escaped double quote remains: \"' \
      'alias kanban-summary := ignored' \
      'even backslashes close: \\"' \
      'command_value := ```' \
      '  printf "%s\n" safe' \
      'run-do-work-update:' \
      'alias kanban-static := ignored' \
      '```' \
      'custom-after:' \
      '    echo after'
  } > "$ordinary_and_command_multiline_target"
  if command -v just >/dev/null 2>&1 \
    && ! just --justfile "$ordinary_and_command_multiline_target" --list >/dev/null 2>&1; then
    printf 'FAIL: ordinary-quote and triple-backtick collision fixture is not valid Just syntax.\n' >&2
    fail_count=$((fail_count + 1))
  elif ! "$replace_section_tool" --target "$ordinary_and_command_multiline_target" --section-file "$reserved_section_file" \
    --reject-recipe-collisions; then
    printf 'FAIL: replace-text-section treated reserved-looking ordinary-quote or triple-backtick content as a collision.\n' >&2
    fail_count=$((fail_count + 1))
  elif command -v just >/dev/null 2>&1 \
    && ! just --justfile "$ordinary_and_command_multiline_target" --list >/dev/null 2>&1; then
    printf 'FAIL: replace-text-section produced invalid Just after accepting ordinary-quote and triple-backtick literals.\n' >&2
    fail_count=$((fail_count + 1))
  fi

  ordinary_backtick_multiline_target="$section_workdir/ordinary-backtick-multiline-command.just"
  {
    printf '%s\n' \
      'custom-before:' \
      '    echo before' \
      'command_value := `' \
      'run-kanban:' \
      'alias kanban-summary := ignored' \
      'raw backslash closes: \` + `' \
      'printf safe' \
      '`' \
      'custom-after:' \
      '    echo after'
  } > "$ordinary_backtick_multiline_target"
  cp "$ordinary_backtick_multiline_target" "$ordinary_backtick_multiline_target.before"
  ordinary_backtick_multiline_expected="$section_workdir/ordinary-backtick-multiline-command.expected"
  cp "$ordinary_backtick_multiline_target" "$ordinary_backtick_multiline_expected"
  printf '\n' >> "$ordinary_backtick_multiline_expected"
  cat "$reserved_section_file" >> "$ordinary_backtick_multiline_expected"
  if command -v just >/dev/null 2>&1 \
    && ! just --justfile "$ordinary_backtick_multiline_target" --list >/dev/null 2>&1; then
    printf 'FAIL: ordinary multiline-backtick collision fixture is not valid Just syntax.\n' >&2
    fail_count=$((fail_count + 1))
  elif ! "$replace_section_tool" --target "$ordinary_backtick_multiline_target" --section-file "$reserved_section_file" \
    --reject-recipe-collisions; then
    printf 'FAIL: replace-text-section treated reserved-looking ordinary multiline-backtick content as a collision.\n' >&2
    fail_count=$((fail_count + 1))
  elif ! cmp -s "$ordinary_backtick_multiline_target" "$ordinary_backtick_multiline_expected"; then
    printf 'FAIL: replace-text-section wrote unexpected bytes after accepting an ordinary multiline-backtick command.\n' >&2
    fail_count=$((fail_count + 1))
  elif command -v just >/dev/null 2>&1 \
    && ! just --justfile "$ordinary_backtick_multiline_target" --list >/dev/null 2>&1; then
    printf 'FAIL: replace-text-section produced invalid Just after accepting an ordinary multiline-backtick command.\n' >&2
    fail_count=$((fail_count + 1))
  fi

  ordinary_backtick_collision_target="$section_workdir/ordinary-backtick-nearby-collisions.just"
  {
    printf '%s\n' \
      '# unmatched ordinary backtick in comment: `' \
      'one_line_command := `printf safe`' \
      'custom-body:' \
      '    echo unmatched-backtick `' \
      'custom-summary:' \
      '    echo custom' \
      'run-kanban:' \
      '    echo real collision before command' \
      'command_value := `' \
      'kanban-static:' \
      'raw backslash closes: \` + `' \
      'alias kanban-static := ignored' \
      '`' \
      'alias kanban-summary := custom-summary' \
      '@run-kanban-cli:' \
      '    echo real collision after command' \
      'run-do-work-update:' \
      '    echo real collision after inactive forms'
  } > "$ordinary_backtick_collision_target"
  if command -v just >/dev/null 2>&1 \
    && ! just --justfile "$ordinary_backtick_collision_target" --list >/dev/null 2>&1; then
    printf 'FAIL: ordinary multiline-backtick nearby-collision control is not valid Just syntax.\n' >&2
    fail_count=$((fail_count + 1))
  else
    cp "$ordinary_backtick_collision_target" "$ordinary_backtick_collision_target.before"
    ordinary_backtick_collision_output="$section_workdir/ordinary-backtick-nearby-collisions.out"
    expected_ordinary_backtick_collision='finding SECTION-RESERVED-RECIPE-COLLISION [error]: target defines reserved Just recipe or alias outside managed section: kanban-summary, run-do-work-update, run-kanban, run-kanban-cli'
    if "$replace_section_tool" --target "$ordinary_backtick_collision_target" --section-file "$reserved_section_file" \
      --reject-recipe-collisions >"$ordinary_backtick_collision_output" 2>&1; then
      printf 'FAIL: replace-text-section ignored real reserved definitions around an ordinary multiline-backtick command.\n' >&2
      fail_count=$((fail_count + 1))
    elif ! grep -Fxq -- "$expected_ordinary_backtick_collision" "$ordinary_backtick_collision_output"; then
      printf 'FAIL: replace-text-section did not report only exact sorted collisions around an ordinary multiline-backtick command.\n' >&2
      fail_count=$((fail_count + 1))
    elif ! cmp -s "$ordinary_backtick_collision_target" "$ordinary_backtick_collision_target.before"; then
      printf 'FAIL: replace-text-section changed the target after rejecting collisions around an ordinary multiline-backtick command.\n' >&2
      fail_count=$((fail_count + 1))
    fi
  fi

  ordinary_and_command_collision_target="$section_workdir/ordinary-and-command-nearby-collisions.just"
  {
    printf '%s\n' \
      'custom-summary:' \
      '    echo custom' \
      'run-kanban:' \
      '    echo real collision before raw string' \
      "raw_value := '" \
      'alias kanban-static := ignored' \
      "raw backslash is literal: \\'" \
      'alias kanban-summary := custom-summary' \
      'cooked_value := "' \
      'kanban-static:' \
      'escaped double quote remains: \"' \
      'alias kanban-static := ignored' \
      'even backslashes close: \\"' \
      '@run-kanban-cli:' \
      '    echo real collision before command literal' \
      'command_value := ```' \
      '  printf "%s\n" safe' \
      'kanban-static:' \
      '```' \
      'run-do-work-update:' \
      '    echo real collision after command literal'
  } > "$ordinary_and_command_collision_target"
  if command -v just >/dev/null 2>&1 \
    && ! just --justfile "$ordinary_and_command_collision_target" --list >/dev/null 2>&1; then
    printf 'FAIL: ordinary-quote and triple-backtick nearby-collision control is not valid Just syntax.\n' >&2
    fail_count=$((fail_count + 1))
  else
    cp "$ordinary_and_command_collision_target" "$ordinary_and_command_collision_target.before"
    ordinary_and_command_collision_output="$section_workdir/ordinary-and-command-nearby-collisions.out"
    expected_ordinary_and_command_collision='finding SECTION-RESERVED-RECIPE-COLLISION [error]: target defines reserved Just recipe or alias outside managed section: kanban-summary, run-do-work-update, run-kanban, run-kanban-cli'
    if "$replace_section_tool" --target "$ordinary_and_command_collision_target" --section-file "$reserved_section_file" \
      --reject-recipe-collisions >"$ordinary_and_command_collision_output" 2>&1; then
      printf 'FAIL: replace-text-section ignored real reserved definitions around ordinary-quote or triple-backtick literals.\n' >&2
      fail_count=$((fail_count + 1))
    elif ! grep -Fxq -- "$expected_ordinary_and_command_collision" "$ordinary_and_command_collision_output"; then
      printf 'FAIL: replace-text-section did not report only exact sorted collisions around ordinary-quote and triple-backtick literals.\n' >&2
      fail_count=$((fail_count + 1))
    elif ! cmp -s "$ordinary_and_command_collision_target" "$ordinary_and_command_collision_target.before"; then
      printf 'FAIL: replace-text-section changed the target after rejecting collisions around ordinary-quote or triple-backtick literals.\n' >&2
      fail_count=$((fail_count + 1))
    fi
  fi

  inactive_literal_forms_target="$section_workdir/inactive-literal-forms-nearby-collision.just"
  {
    printf '%s\n' \
      '# delimiter-looking comment: ```' \
      'one_line_command := `printf safe`' \
      'custom-body:' \
      "    printf '%s\\n' '\`\`\`'" \
      'kanban-static:' \
      '    echo real collision after inactive forms'
  } > "$inactive_literal_forms_target"
  if command -v just >/dev/null 2>&1 \
    && ! just --justfile "$inactive_literal_forms_target" --list >/dev/null 2>&1; then
    printf 'FAIL: inactive multiline-literal opener control is not valid Just syntax.\n' >&2
    fail_count=$((fail_count + 1))
  else
    cp "$inactive_literal_forms_target" "$inactive_literal_forms_target.before"
    inactive_literal_forms_output="$section_workdir/inactive-literal-forms-nearby-collision.out"
    expected_inactive_literal_forms_collision='finding SECTION-RESERVED-RECIPE-COLLISION [error]: target defines reserved Just recipe or alias outside managed section: kanban-static'
    if "$replace_section_tool" --target "$inactive_literal_forms_target" --section-file "$reserved_section_file" \
      --reject-recipe-collisions >"$inactive_literal_forms_output" 2>&1; then
      printf 'FAIL: replace-text-section let an inactive comment, recipe body, or one-line backtick hide a real collision.\n' >&2
      fail_count=$((fail_count + 1))
    elif ! grep -Fxq -- "$expected_inactive_literal_forms_collision" "$inactive_literal_forms_output"; then
      printf 'FAIL: replace-text-section did not report the exact collision after inactive multiline-literal opener forms.\n' >&2
      fail_count=$((fail_count + 1))
    elif ! cmp -s "$inactive_literal_forms_target" "$inactive_literal_forms_target.before"; then
      printf 'FAIL: replace-text-section changed the target after rejecting a collision after inactive literal forms.\n' >&2
      fail_count=$((fail_count + 1))
    fi
  fi

  multiline_collision_target="$section_workdir/multiline-nearby-collisions.just"
  {
    printf '%s\n' \
      '# delimiter-looking comment: """' \
      "single_token := '\"\"\"'" \
      "double_token := \"'''\"" \
      'backtick_token := `printf "\"\"\""`' \
      'custom-shell:' \
      "    echo '\"\"\"'" \
      'custom-summary:' \
      '    echo custom' \
      'alias kanban-summary := custom-summary' \
      'payload := """' \
      'run-kanban:' \
      'alias run-kanban-cli := ignored' \
      'even escaped delimiter closes: \\"""' \
      'run-do-work-update:' \
      '    echo real collision'
  } > "$multiline_collision_target"
  cp "$multiline_collision_target" "$multiline_collision_target.before"
  multiline_collision_output="$section_workdir/multiline-nearby-collisions.out"
  expected_multiline_collision='finding SECTION-RESERVED-RECIPE-COLLISION [error]: target defines reserved Just recipe or alias outside managed section: kanban-summary, run-do-work-update'
  if "$replace_section_tool" --target "$multiline_collision_target" --section-file "$reserved_section_file" \
    --reject-recipe-collisions >"$multiline_collision_output" 2>&1; then
    printf 'FAIL: replace-text-section ignored real reserved definitions around a multiline string.\n' >&2
    fail_count=$((fail_count + 1))
  elif ! grep -Fxq -- "$expected_multiline_collision" "$multiline_collision_output"; then
    printf 'FAIL: replace-text-section did not report only exact sorted real collisions around a multiline string.\n' >&2
    fail_count=$((fail_count + 1))
  elif ! cmp -s "$multiline_collision_target" "$multiline_collision_target.before"; then
    printf 'FAIL: replace-text-section changed the target after rejecting collisions around a multiline string.\n' >&2
    fail_count=$((fail_count + 1))
  fi

  noncollision_target="$section_workdir/noncollisions.just"
  {
    printf '# run-kanban:\n'
    printf 'reserved_value := "run-kanban-cli:"\n'
    printf "[doc('kanban-static: is reserved')]\n"
    printf 'custom-recipe: kanban-summary\n'
    printf '    echo run-do-work-update:\n'
    printf 'run-kanban-extra:\n    echo prefix\n'
    printf 'alias custom-summary := kanban-summary\n\n'
    cat "$reserved_section_file"
  } > "$noncollision_target"
  if ! "$replace_section_tool" --target "$noncollision_target" --section-file "$reserved_section_file" \
    --reject-recipe-collisions; then
    printf 'FAIL: replace-text-section treated comments, variables, attributes, dependencies, bodies, prefixes, aliases to reserved recipes, or managed definitions as collisions.\n' >&2
    fail_count=$((fail_count + 1))
  fi

  external_collision_target="$section_workdir/external-after-managed.just"
  cat "$reserved_section_file" > "$external_collision_target"
  printf '\nrun-kanban:\n    echo external collision\n' >> "$external_collision_target"
  cp "$external_collision_target" "$external_collision_target.before"
  if "$replace_section_tool" --target "$external_collision_target" --section-file "$reserved_section_file" \
    --reject-recipe-collisions >/dev/null 2>&1; then
    printf 'FAIL: replace-text-section ignored a reserved recipe outside an existing managed span.\n' >&2
    fail_count=$((fail_count + 1))
  elif ! cmp -s "$external_collision_target" "$external_collision_target.before"; then
    printf 'FAIL: replace-text-section mutated a managed target before rejecting its external collision.\n' >&2
    fail_count=$((fail_count + 1))
  fi

  retired_flag_target="$section_workdir/retired-flag.just"
  printf 'custom-only:\n    echo untouched\n' > "$retired_flag_target"
  cp "$retired_flag_target" "$retired_flag_target.before"
  if "$replace_section_tool" --target "$retired_flag_target" --section-file "$section_file" --migrate-legacy-do-work >/dev/null 2>&1; then
    printf 'FAIL: replace-text-section still accepts the retired legacy-migration flag.\n' >&2
    fail_count=$((fail_count + 1))
  elif ! cmp -s "$retired_flag_target" "$retired_flag_target.before"; then
    printf 'FAIL: replace-text-section changed the target after rejecting the retired flag.\n' >&2
    fail_count=$((fail_count + 1))
  fi

  malformed_index=0
  for malformed_content in \
    '# >>> do-work:recipes >>>|one|# >>> do-work:recipes >>>|# <<< do-work:recipes <<<' \
    '# <<< do-work:recipes <<<|one|# >>> do-work:recipes >>>' \
    '# >>> do-work:recipes >>>|one' \
    'one|# <<< do-work:recipes <<<' ; do
    malformed_index=$((malformed_index + 1))
    malformed_target="$section_workdir/malformed-$malformed_index.just"
    printf '%s\n' "$malformed_content" | tr '|' '\n' > "$malformed_target"
    cp "$malformed_target" "$malformed_target.before"
    if "$replace_section_tool" --target "$malformed_target" --section-file "$section_file" >/dev/null 2>&1; then
      printf 'FAIL: replace-text-section accepted malformed marker case %s.\n' "$malformed_index" >&2
      fail_count=$((fail_count + 1))
    elif ! cmp -s "$malformed_target" "$malformed_target.before"; then
      printf 'FAIL: replace-text-section changed the target after rejecting malformed marker case %s.\n' "$malformed_index" >&2
      fail_count=$((fail_count + 1))
    fi
  done

  if command -v just >/dev/null 2>&1 \
    && ! just --justfile "$section_workdir/variants/justfile" --list >/dev/null 2>&1; then
    printf 'FAIL: replace-text-section produced a Justfile that does not parse.\n' >&2
    fail_count=$((fail_count + 1))
  fi
  rm -rf "$section_workdir"
fi

[ "$fail_count" -eq 0 ] || exit 1
printf 'replace-text-section contract probes passed.\n'
