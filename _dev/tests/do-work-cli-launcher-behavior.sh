#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "$0")/../.." && pwd)"
source_launcher="$repo_root/skills/do-work/tools/do-work-cli.sh"
source_module="$repo_root/skills/do-work/tools/do-work-cli"

if [ ! -x "$source_launcher" ]; then
  echo "FAIL: expected executable do-work-cli launcher at $source_launcher" >&2
  exit 1
fi

# The launcher runs `go tool do-work-cli`, so the module must declare the tool.
if ! grep -q '^tool github.com/knews2019/skill-do-work/do-work-cli/cmd/do-work-cli$' "$source_module/go.mod"; then
  echo "FAIL: go.mod does not declare the do-work-cli tool directive the launcher runs" >&2
  exit 1
fi

fixture_root="$(mktemp -d "${TMPDIR:-/tmp}/do-work-cli-launcher.XXXXXX")"
trap 'rm -rf "$fixture_root"' EXIT
mkdir -p "$fixture_root/tools/do-work-cli" "$fixture_root/fake-bin" "$fixture_root/no-go-bin"
cp "$source_launcher" "$fixture_root/tools/do-work-cli.sh"
chmod +x "$fixture_root/tools/do-work-cli.sh"

# A fake toolchain that records how it was invoked and behaves like `go tool -n`:
# it "builds" the named tool into a private cache and prints that executable's path.
# The stub executable reports its argv and exits with the status the case asks for.
cat > "$fixture_root/fake-bin/go" <<'FAKE_GO'
#!/usr/bin/env bash
set -euo pipefail
if [ "${1:-}" = version ]; then
  printf 'go version %s darwin/arm64\n' "${FAKE_GO_VERSION:-go1.25.0}"
  exit 0
fi
if [ "${1:-}" = tool ] && [ "${2:-}" = -C ] && [ "${4:-}" = -n ] && [ "$#" -eq 5 ]; then
  printf 'tool %s %s\n' "$3" "$5" >> "$FAKE_GO_LOG"
  if [ "${FAKE_GO_BUILD_FAIL:-0}" = 1 ]; then
    echo "fake compile error" >&2
    exit 1
  fi
  mkdir -p "$FAKE_GO_CACHE"
  cat > "$FAKE_GO_CACHE/$5" <<'CACHED_TOOL'
#!/usr/bin/env bash
printf 'ran'
printf ' <%s>' "$@"
printf '\n'
exit "${FAKE_GO_TOOL_EXIT:-0}"
CACHED_TOOL
  chmod +x "$FAKE_GO_CACHE/$5"
  printf '%s\n' "$FAKE_GO_CACHE/$5"
  exit 0
fi
exit 72
FAKE_GO
chmod +x "$fixture_root/fake-bin/go"
for command_name in bash; do
  ln -s "$(command -v "$command_name")" "$fixture_root/no-go-bin/$command_name"
done

fake_log="$fixture_root/go.log"
: > "$fake_log"
run_launcher() {
  PATH="$fixture_root/fake-bin:$PATH" FAKE_GO_LOG="$fake_log" FAKE_GO_CACHE="$fixture_root/cache" \
    "$fixture_root/tools/do-work-cli.sh" "$@"
}

# Argv reaches the tool byte-for-byte, and the module directory is the one beside the launcher.
output="$(run_launcher --format json inspect 'two words')"
if [ "$output" != 'ran <--format> <json> <inspect> <two words>' ]; then
  echo "FAIL: launcher did not preserve argv: $output" >&2
  exit 1
fi
if [ "$(cat "$fake_log")" != "tool $(cd "$fixture_root/tools/do-work-cli" && pwd -P) do-work-cli" ]; then
  echo "FAIL: launcher did not run the module's do-work-cli tool: $(cat "$fake_log")" >&2
  exit 1
fi

# The tool's exact exit status is the launcher's exit status, with nothing added to its output.
set +e
status_output="$(FAKE_GO_TOOL_EXIT=3 run_launcher inspect 2>&1)"
tool_status=$?
set -e
if [ "$tool_status" -ne 3 ] || [ "$status_output" != 'ran <inspect>' ]; then
  echo "FAIL: tool exit status 3 was not propagated cleanly (status $tool_status): $status_output" >&2
  exit 1
fi

# The Go floor refusal is actionable and never reaches the tool.
set +e
old_output="$(FAKE_GO_VERSION=go1.24.99 run_launcher inspect 2>&1)"
old_status=$?
set -e
if [ "$old_status" -ne 2 ] || [[ "$old_output" != *'Go 1.25.0 or newer'* ]] || [ "$(wc -l < "$fake_log" | tr -d ' ')" != 2 ]; then
  echo "FAIL: old Go refusal was not actionable (status $old_status): $old_output" >&2
  exit 1
fi

# A toolchain that cannot build the command refuses actionably and runs nothing.
set +e
build_output="$(FAKE_GO_BUILD_FAIL=1 run_launcher inspect 2>&1)"
build_status=$?
set -e
if [ "$build_status" -ne 2 ] || [[ "$build_output" != *'could not build the command'* ]] || [[ "$build_output" == *ran* ]]; then
  echo "FAIL: build failure was not refused actionably (status $build_status): $build_output" >&2
  exit 1
fi

# A missing toolchain is refused the same way.
set +e
missing_output="$(PATH="$fixture_root/no-go-bin" "$fixture_root/tools/do-work-cli.sh" inspect 2>&1)"
missing_status=$?
set -e
if [ "$missing_status" -ne 2 ] || [[ "$missing_output" != *'Go 1.25.0 or newer'* ]]; then
  echo "FAIL: missing Go refusal was not actionable (status $missing_status): $missing_output" >&2
  exit 1
fi

# The launcher writes nothing into the module tree; the cached executable belongs to GOCACHE.
# It also runs the tool from the caller's directory, not the module's.
if [ -n "$(ls -A "$fixture_root/tools/do-work-cli")" ]; then
  echo "FAIL: launcher wrote into the module directory: $(ls -A "$fixture_root/tools/do-work-cli")" >&2
  exit 1
fi

cwd_output="$(cd "$fixture_root" && run_launcher pwd-probe)"
if [ "$cwd_output" != 'ran <pwd-probe>' ]; then
  echo "FAIL: launcher altered the tool invocation: $cwd_output" >&2
  exit 1
fi

echo "do-work-cli launcher behavior tests passed"
