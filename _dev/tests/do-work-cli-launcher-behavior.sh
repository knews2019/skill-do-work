#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "$0")/../.." && pwd)"
source_launcher="$repo_root/skills/do-work/tools/do-work-cli.sh"

if [ ! -x "$source_launcher" ]; then
  echo "FAIL: expected executable do-work-cli launcher at $source_launcher" >&2
  exit 1
fi

fixture_root="$(mktemp -d "${TMPDIR:-/tmp}/do-work-cli-launcher.XXXXXX")"
trap 'rm -rf "$fixture_root"' EXIT
mkdir -p "$fixture_root/tools/do-work-cli/cmd/do-work-cli" "$fixture_root/fake-bin"
cp "$source_launcher" "$fixture_root/tools/do-work-cli.sh"
printf 'package main\n' > "$fixture_root/tools/do-work-cli/cmd/do-work-cli/main.go"
printf 'module example.invalid/do-work-cli\n\ngo 1.26.1\n' > "$fixture_root/tools/do-work-cli/go.mod"

cat > "$fixture_root/fake-bin/go" <<'FAKE_GO'
#!/usr/bin/env bash
set -euo pipefail
if [ "${1:-}" = version ]; then
  printf 'go version %s darwin/arm64\n' "${FAKE_GO_VERSION:-go1.26.1}"
  exit 0
fi
if [ "${1:-}" = build ]; then
  printf 'build\n' >> "$FAKE_GO_BUILD_LOG"
  if [ "${FAKE_GO_FAIL_BUILD:-0}" = 1 ]; then
    exit 71
  fi
  output_path=""
  shift
  while [ "$#" -gt 0 ]; do
    if [ "$1" = -o ]; then
      output_path="$2"
      shift 2
      continue
    fi
    shift
  done
  cat > "$output_path" <<'BUILT_BINARY'
#!/usr/bin/env bash
printf 'built'
printf ' <%s>' "$@"
printf '\n'
BUILT_BINARY
  chmod +x "$output_path"
  exit 0
fi
exit 72
FAKE_GO
chmod +x "$fixture_root/fake-bin/go" "$fixture_root/tools/do-work-cli.sh"

build_log="$fixture_root/build.log"
: > "$build_log"
run_launcher() {
  PATH="$fixture_root/fake-bin:$PATH" FAKE_GO_BUILD_LOG="$build_log" \
    "$fixture_root/tools/do-work-cli.sh" "$@"
}

output="$(run_launcher --format json inspect 'two words')"
if [ "$output" != 'built <--format> <json> <inspect> <two words>' ]; then
  echo "FAIL: launcher did not preserve argv: $output" >&2
  exit 1
fi
if [ "$(wc -l < "$build_log" | tr -d ' ')" != 1 ]; then
  echo "FAIL: missing binary did not trigger exactly one build" >&2
  exit 1
fi

run_launcher inspect >/dev/null
if [ "$(wc -l < "$build_log" | tr -d ' ')" != 1 ]; then
  echo "FAIL: fresh binary rebuilt" >&2
  exit 1
fi

sleep 1
touch "$fixture_root/tools/do-work-cli/cmd/do-work-cli/main.go"
run_launcher inspect >/dev/null
if [ "$(wc -l < "$build_log" | tr -d ' ')" != 2 ]; then
  echo "FAIL: newer Go source did not rebuild" >&2
  exit 1
fi

sleep 1
touch "$fixture_root/tools/do-work-cli/cmd/do-work-cli/main.go"
set +e
failed_output="$(PATH="$fixture_root/fake-bin:$PATH" FAKE_GO_BUILD_LOG="$build_log" FAKE_GO_FAIL_BUILD=1 "$fixture_root/tools/do-work-cli.sh" inspect 2>&1)"
failed_status=$?
set -e
if [ "$failed_status" -eq 0 ] || [[ "$failed_output" == built* ]]; then
  echo "FAIL: failed rebuild ran stale output (status $failed_status): $failed_output" >&2
  exit 1
fi

rm -f "$fixture_root/tools/do-work-cli/do-work-cli"
set +e
old_output="$(PATH="$fixture_root/fake-bin:$PATH" FAKE_GO_BUILD_LOG="$build_log" FAKE_GO_VERSION=go1.26.0 "$fixture_root/tools/do-work-cli.sh" inspect 2>&1)"
old_status=$?
set -e
if [ "$old_status" -eq 0 ] || [[ "$old_output" != *'Go 1.26.1 or newer'* ]]; then
  echo "FAIL: old Go refusal was not actionable (status $old_status): $old_output" >&2
  exit 1
fi

if find "$fixture_root/tools/do-work-cli" -maxdepth 1 -name 'do-work-cli.build.*' -print -quit | grep -q .; then
  echo "FAIL: launcher left a temporary build artifact" >&2
  exit 1
fi

echo "do-work-cli launcher behavior tests passed"
