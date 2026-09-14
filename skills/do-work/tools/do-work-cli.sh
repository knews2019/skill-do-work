#!/usr/bin/env bash
# do-work-cli compatibility launcher: retained public path
set -euo pipefail

script_dir="${BASH_SOURCE[0]%/*}"
[[ "$script_dir" != "${BASH_SOURCE[0]}" ]] || script_dir=.
script_dir="$(cd "$script_dir" && pwd -P)"
module_dir="$script_dir/do-work-cli"
minimum_go_version="1.24.0"
# Where prebuilt binaries live when no usable Go toolchain is on PATH. Point it at a
# mirror (any URL curl reads, file:// included) on hosts that cannot reach GitHub.
release_base_url="${DO_WORK_CLI_RELEASE_BASE:-https://github.com/knews2019/skill-do-work/releases/download}"

# Pure bash on purpose: this runs on every invocation, and restricted PATHs that
# carry only bash and go must still reach the command. A component keeps its
# leading digits ("26rc1" reads as 26); an absent component reads as 0.
version_at_least() {
  local -a required_parts found_parts
  local component required found
  IFS=. read -r -a required_parts <<<"$minimum_go_version"
  IFS=. read -r -a found_parts <<<"$1"
  for component in 0 1 2; do
    required="${required_parts[component]:-0}"; required="${required%%[^0-9]*}"
    found="${found_parts[component]:-0}"; found="${found%%[^0-9]*}"
    if [ "${found:-0}" -gt "${required:-0}" ]; then return 0; fi
    if [ "${found:-0}" -lt "${required:-0}" ]; then return 1; fi
  done
  return 0
}

# toolchain_refusal says why the Go route is unusable; empty means the toolchain builds
# the command. A refusal is not yet an exit: the prebuilt route below gets its turn.
toolchain_refusal=""
if ! command -v go >/dev/null 2>&1; then
  toolchain_refusal="Go $minimum_go_version or newer is required to run the command and no go was found on PATH"
elif ! go_version_output="$(go version 2>/dev/null)"; then
  toolchain_refusal="could not read the installed Go version; Go $minimum_go_version or newer is required"
else
  read -r _ _ go_version _ <<<"$go_version_output"
  go_version="${go_version#go}"
  if [ -z "$go_version" ] || ! version_at_least "$go_version"; then
    toolchain_refusal="Go $minimum_go_version or newer is required (found ${go_version:-unknown})"
  fi
fi

if [ -z "$toolchain_refusal" ]; then
  # The Go toolchain owns the build. `go tool` compiles the `tool` directive in
  # go.mod, caches the linked executable in GOCACHE under a hash of every input
  # (sources, go.mod, flags, toolchain), and reuses it while nothing changed. No
  # binary lives in this tree and no timestamp decides staleness. `-n` prints the
  # cached executable's path instead of running it, so the command is exec'd from
  # the caller's own directory: `go tool -C` would otherwise run it inside the
  # module, where relative arguments and the default repository root resolve to
  # the wrong tree. `go run` would cache the same way but collapses every exit
  # status to 1, and callers read the command's exact status.
  tool_binary="$(go tool -C "$module_dir" -n do-work-cli)" || {
    echo "do-work-cli: the Go toolchain could not build the command; nothing was run" >&2
    exit 2
  }
  # Go 1.24's first, uncached `go tool -n` builds into the cache but prints the temporary
  # build path instead of the cached one (1.25 prints the cached path straight away), so
  # a path that is not there means "built, ask again", never a failure.
  if [ ! -x "$tool_binary" ]; then
    tool_binary="$(go tool -C "$module_dir" -n do-work-cli)" || {
      echo "do-work-cli: the Go toolchain could not build the command; nothing was run" >&2
      exit 2
    }
  fi
  exec "$tool_binary" "$@"
fi

# Prebuilt route: every release publishes one static binary per supported platform plus
# a SHA256SUMS file (see .github/workflows/release-binaries.yml). The binary for the
# installed suite version is fetched once, verified against SHA256SUMS, and cached under
# the user cache directory; later invocations exec the cached copy without any network.
# The version comes from the installed actions/version.md, so the binary always matches
# the shipped sources beside this launcher.
refuse_prebuilt() {
  echo "do-work-cli: $toolchain_refusal" >&2
  echo "do-work-cli: no prebuilt binary could be used instead: $1" >&2
  exit 2
}

version_file="$script_dir/../actions/version.md"
suite_version=""
if [ -r "$version_file" ]; then
  while IFS= read -r version_line || [ -n "$version_line" ]; do
    case "$version_line" in
      '**Current version**: '*)
        suite_version="${version_line#\*\*Current version\*\*: }"
        suite_version="${suite_version%%[[:space:]]*}"
        break
        ;;
    esac
  done < "$version_file"
fi
[ -n "$suite_version" ] || refuse_prebuilt "could not read **Current version** from $version_file"

host_os="$(uname -s 2>/dev/null || true)"
host_arch="$(uname -m 2>/dev/null || true)"
case "$host_os" in
  Linux) host_os=linux ;;
  Darwin) host_os=darwin ;;
  *) refuse_prebuilt "no prebuilt binary is published for $host_os/$host_arch" ;;
esac
case "$host_arch" in
  x86_64 | amd64) host_arch=amd64 ;;
  aarch64 | arm64) host_arch=arm64 ;;
  *) refuse_prebuilt "no prebuilt binary is published for $host_os/$host_arch" ;;
esac
asset_name="do-work-cli_${host_os}_${host_arch}"
cache_directory="${XDG_CACHE_HOME:-${HOME:-/tmp}/.cache}/do-work-cli/$suite_version"
cached_binary="$cache_directory/$asset_name"

if [ -x "$cached_binary" ]; then
  exec "$cached_binary" "$@"
fi

command -v curl >/dev/null 2>&1 || refuse_prebuilt "curl is required to fetch $asset_name $suite_version"
release_url="$release_base_url/v$suite_version"
mkdir -p "$cache_directory" || refuse_prebuilt "could not create $cache_directory"
staging_directory="$(mktemp -d "$cache_directory/.fetch.XXXXXX")" || refuse_prebuilt "could not create a staging directory under $cache_directory"
trap 'rm -rf "$staging_directory"' EXIT
if ! curl -fsSL --retry 3 --retry-delay 2 -o "$staging_directory/$asset_name" "$release_url/$asset_name" \
  || ! curl -fsSL --retry 3 --retry-delay 2 -o "$staging_directory/SHA256SUMS" "$release_url/SHA256SUMS"; then
  refuse_prebuilt "could not fetch $asset_name $suite_version from $release_url"
fi

expected_digest=""
while read -r digest_field name_field; do
  if [ "${name_field#\*}" = "$asset_name" ]; then
    expected_digest="$digest_field"
    break
  fi
done < "$staging_directory/SHA256SUMS"
[ -n "$expected_digest" ] || refuse_prebuilt "SHA256SUMS at $release_url does not list $asset_name"
if command -v sha256sum >/dev/null 2>&1; then
  actual_digest="$(sha256sum "$staging_directory/$asset_name")" || refuse_prebuilt "sha256sum failed on the fetched $asset_name"
elif command -v shasum >/dev/null 2>&1; then
  actual_digest="$(shasum -a 256 "$staging_directory/$asset_name")" || refuse_prebuilt "shasum failed on the fetched $asset_name"
else
  refuse_prebuilt "neither sha256sum nor shasum is available to verify $asset_name"
fi
actual_digest="${actual_digest%% *}"
[ "$actual_digest" = "$expected_digest" ] || refuse_prebuilt "checksum mismatch for $asset_name $suite_version from $release_url; nothing was run"
chmod 0755 "$staging_directory/$asset_name"
# Same directory, so this is one atomic rename; a concurrent launcher lands identical bytes.
mv -f "$staging_directory/$asset_name" "$cached_binary"
rm -rf "$staging_directory"
trap - EXIT
echo "do-work-cli: $toolchain_refusal; using the prebuilt do-work-cli $suite_version for $host_os/$host_arch, cached at $cached_binary" >&2
exec "$cached_binary" "$@"
