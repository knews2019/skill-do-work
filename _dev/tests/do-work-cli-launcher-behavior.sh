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
mkdir -p "$fixture_root/tools/do-work-cli" "$fixture_root/actions" "$fixture_root/fake-bin" "$fixture_root/no-go-bin"
cp "$source_launcher" "$fixture_root/tools/do-work-cli.sh"
chmod +x "$fixture_root/tools/do-work-cli.sh"
printf '# Version\n\n**Current version**: 9.8.7\n' > "$fixture_root/actions/version.md"

# A fake toolchain that records how it was invoked and behaves like `go tool -n`:
# it "builds" the named tool into a private cache and prints that executable's path.
# The stub executable reports its argv and exits with the status the case asks for.
cat > "$fixture_root/fake-bin/go" <<'FAKE_GO'
#!/usr/bin/env bash
set -euo pipefail
if [ "${1:-}" = version ]; then
  printf 'go version %s darwin/arm64\n' "${FAKE_GO_VERSION:-go1.24.0}"
  exit 0
fi
if [ "${1:-}" = tool ] && [ "${2:-}" = -C ] && [ "${4:-}" = -n ] && [ "$#" -eq 5 ]; then
  printf 'tool %s %s\n' "$3" "$5" >> "$FAKE_GO_LOG"
  if [ "${FAKE_GO_BUILD_FAIL:-0}" = 1 ]; then
    echo "fake compile error" >&2
    exit 1
  fi
  # Go 1.24 shape: the first uncached -n prints a temporary build path that no longer exists.
  if [ -n "${FAKE_GO_STALE_MARKER:-}" ] && [ ! -e "$FAKE_GO_STALE_MARKER" ]; then
    : > "$FAKE_GO_STALE_MARKER"
    printf '/tmp/go-build000000/b001/exe/%s\n' "$5"
    exit 0
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
# A PATH without go still carries what the prebuilt route needs (curl, uname, a checksum tool).
for tool in bash curl uname mktemp chmod mv rm mkdir sha256sum shasum; do
  if command -v "$tool" >/dev/null 2>&1; then ln -s "$(command -v "$tool")" "$fixture_root/no-go-bin/$tool"; fi
done

# Every launcher run points the prebuilt route at a file:// base under the fixture, so a
# refusal case never reaches the network and a fallback case reads the release we stage.
export DO_WORK_CLI_RELEASE_BASE="file://$fixture_root/releases"
export XDG_CACHE_HOME="$fixture_root/xdg-cache"
fake_log="$fixture_root/go.log"
: > "$fake_log"
run_launcher() {
  PATH="$fixture_root/fake-bin:$PATH" FAKE_GO_LOG="$fake_log" FAKE_GO_CACHE="$fixture_root/cache" \
    "$fixture_root/tools/do-work-cli.sh" "$@"
}
run_launcher_without_go() {
  PATH="$fixture_root/no-go-bin" "$fixture_root/tools/do-work-cli.sh" "$@"
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

# Go 1.24 prints a temporary build path from the first uncached `go tool -n`; the launcher
# asks once more instead of exec'ing a path that is not there.
stale_output="$(FAKE_GO_STALE_MARKER="$fixture_root/stale-once" run_launcher inspect 2>&1)"
if [ "$stale_output" != 'ran <inspect>' ] || [ "$(wc -l < "$fake_log" | tr -d ' ')" != 3 ]; then
  echo "FAIL: launcher did not re-ask for the cached path after a stale first -n: $stale_output ($(cat "$fake_log"))" >&2
  exit 1
fi
: > "$fake_log"
printf 'tool %s do-work-cli\n' "$(cd "$fixture_root/tools/do-work-cli" && pwd -P)" > "$fake_log"

# The tool's exact exit status is the launcher's exit status, with nothing added to its output.
set +e
status_output="$(FAKE_GO_TOOL_EXIT=3 run_launcher inspect 2>&1)"
tool_status=$?
set -e
if [ "$tool_status" -ne 3 ] || [ "$status_output" != 'ran <inspect>' ]; then
  echo "FAIL: tool exit status 3 was not propagated cleanly (status $tool_status): $status_output" >&2
  exit 1
fi

# The Go floor refusal is actionable and never reaches the tool; with no release staged
# the prebuilt route refuses too, and both reasons are named.
set +e
old_output="$(FAKE_GO_VERSION=go1.23.99 run_launcher inspect 2>&1)"
old_status=$?
set -e
if [ "$old_status" -ne 2 ] || [[ "$old_output" != *'Go 1.24.0 or newer'* ]] || [[ "$old_output" != *'no prebuilt binary could be used'* ]] || [ "$(wc -l < "$fake_log" | tr -d ' ')" != 2 ]; then
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
missing_output="$(run_launcher_without_go inspect 2>&1)"
missing_status=$?
set -e
if [ "$missing_status" -ne 2 ] || [[ "$missing_output" != *'Go 1.24.0 or newer'* ]] || [[ "$missing_output" != *'could not fetch'* ]]; then
  echo "FAIL: missing Go refusal was not actionable (status $missing_status): $missing_output" >&2
  exit 1
fi
if [ -e "$XDG_CACHE_HOME/do-work-cli/9.8.7" ] && [ -n "$(ls -A "$XDG_CACHE_HOME/do-work-cli/9.8.7")" ]; then
  echo "FAIL: a failed fetch left files in the cache: $(ls -A "$XDG_CACHE_HOME/do-work-cli/9.8.7")" >&2
  exit 1
fi

# Prebuilt route: with a release staged for the installed version, a host without Go runs
# the checksum-verified binary, argv and exit status intact, and caches it for next time.
host_os="$(uname -s | tr '[:upper:]' '[:lower:]')"
host_arch="$(uname -m)"
case "$host_arch" in x86_64) host_arch=amd64 ;; aarch64) host_arch=arm64 ;; esac
asset_name="do-work-cli_${host_os}_${host_arch}"
release_directory="$fixture_root/releases/v9.8.7"
mkdir -p "$release_directory"
cat > "$release_directory/$asset_name" <<'PREBUILT'
#!/usr/bin/env bash
printf 'prebuilt'
printf ' <%s>' "$@"
printf '\n'
exit "${FAKE_PREBUILT_EXIT:-0}"
PREBUILT
chmod -x "$release_directory/$asset_name"
(cd "$release_directory" && { sha256sum "$asset_name" 2>/dev/null || shasum -a 256 "$asset_name"; } > SHA256SUMS)
set +e
fallback_output="$(run_launcher_without_go --format json inspect 'two words' 2>"$fixture_root/fallback.stderr")"
fallback_status=$?
set -e
cached_binary="$XDG_CACHE_HOME/do-work-cli/9.8.7/$asset_name"
if [ "$fallback_status" -ne 0 ] || [ "$fallback_output" != 'prebuilt <--format> <json> <inspect> <two words>' ] || [ ! -x "$cached_binary" ] || ! grep -q 'using the prebuilt do-work-cli 9.8.7' "$fixture_root/fallback.stderr"; then
  echo "FAIL: prebuilt fallback did not run the released binary (status $fallback_status): $fallback_output / $(cat "$fixture_root/fallback.stderr")" >&2
  exit 1
fi
for leftover in "$XDG_CACHE_HOME/do-work-cli/9.8.7"/.fetch.*; do
  if [ -e "$leftover" ]; then
    echo "FAIL: fetch staging directory survived a successful fallback: $leftover" >&2
    exit 1
  fi
done

# The cached copy serves later runs with no release reachable and nothing on stderr; the
# binary's exit status is the launcher's exit status.
rm -rf "$release_directory"
set +e
cached_output="$(FAKE_PREBUILT_EXIT=3 run_launcher_without_go inspect 2>"$fixture_root/cached.stderr")"
cached_status=$?
set -e
if [ "$cached_status" -ne 3 ] || [ "$cached_output" != 'prebuilt <inspect>' ] || [ -s "$fixture_root/cached.stderr" ]; then
  echo "FAIL: cached prebuilt binary was not reused cleanly (status $cached_status): $cached_output / $(cat "$fixture_root/cached.stderr")" >&2
  exit 1
fi

# A too-old toolchain takes the same prebuilt route.
set +e
old_fallback_output="$(FAKE_GO_VERSION=go1.23.99 run_launcher inspect 2>/dev/null)"
old_fallback_status=$?
set -e
if [ "$old_fallback_status" -ne 0 ] || [ "$old_fallback_output" != 'prebuilt <inspect>' ]; then
  echo "FAIL: old Go did not fall back to the cached prebuilt binary (status $old_fallback_status): $old_fallback_output" >&2
  exit 1
fi

# A binary whose bytes do not match SHA256SUMS is refused, never cached, never run.
rm -rf "$XDG_CACHE_HOME/do-work-cli"
mkdir -p "$release_directory"
printf '#!/usr/bin/env bash\nprintf tampered\n' > "$release_directory/$asset_name"
printf '%s  %s\n' "$(printf '%064d' 0)" "$asset_name" > "$release_directory/SHA256SUMS"
set +e
tampered_output="$(run_launcher_without_go inspect 2>&1)"
tampered_status=$?
set -e
if [ "$tampered_status" -ne 2 ] || [[ "$tampered_output" != *'checksum mismatch'* ]] || [[ "$tampered_output" == *tampered* ]] || [ -e "$cached_binary" ]; then
  echo "FAIL: tampered prebuilt binary was not refused (status $tampered_status): $tampered_output" >&2
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
