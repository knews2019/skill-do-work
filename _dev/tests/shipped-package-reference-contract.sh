#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

python3 - "$repo_root" <<'PY'
import csv
import os
import pathlib
import re
import subprocess
import sys
import urllib.parse

repo_root = pathlib.Path(sys.argv[1]).resolve()
manifest_path = repo_root / "suite/modules.tsv"
repository_slug = "knews2019/skill-do-work"


def fail(message):
    print(f"FAIL: {message}", file=sys.stderr)


def git_output(*arguments):
    return subprocess.check_output(
        ["git", "-C", os.fspath(repo_root), *arguments],
        stderr=subprocess.DEVNULL,
    )


def read_manifest():
    try:
        with manifest_path.open(newline="", encoding="utf-8") as manifest_file:
            rows = list(csv.DictReader(manifest_file, delimiter="\t"))
    except (OSError, csv.Error) as error:
        fail(f"cannot read suite/modules.tsv: {error}")
        raise SystemExit(1)

    if not rows or set(rows[0]) != {"source", "destination"}:
        fail("suite/modules.tsv must declare source and destination columns")
        raise SystemExit(1)

    modules = []
    for row in rows:
        source = pathlib.PurePosixPath(row["source"])
        destination = pathlib.PurePosixPath(row["destination"])
        if source.is_absolute() or destination.is_absolute() or ".." in source.parts or ".." in destination.parts:
            fail(f"unsafe module mapping: {source} -> {destination}")
            raise SystemExit(1)
        modules.append((source, destination))
    return modules


def strip_markdown_code(markdown_text):
    output_lines = []
    fence_character = None
    fence_length = 0

    for line in markdown_text.splitlines(keepends=True):
        fence_match = re.match(r"^[ ]{0,3}(`{3,}|~{3,})", line)
        if fence_match:
            fence = fence_match.group(1)
            if fence_character is None:
                fence_character = fence[0]
                fence_length = len(fence)
            elif fence[0] == fence_character and len(fence) >= fence_length:
                fence_character = None
                fence_length = 0
            output_lines.append("\n" if line.endswith("\n") else "")
            continue
        if fence_character is not None:
            output_lines.append("\n" if line.endswith("\n") else "")
            continue
        output_lines.append(line)

    rendered_text = "".join(output_lines)
    code_free_text = list(rendered_text)
    index = 0
    while index < len(code_free_text):
        if code_free_text[index] != "`":
            index += 1
            continue
        run_end = index
        while run_end < len(code_free_text) and code_free_text[run_end] == "`":
            run_end += 1
        marker = "`" * (run_end - index)
        closing = rendered_text.find(marker, run_end)
        if closing < 0:
            index = run_end
            continue
        for code_index in range(index, closing + len(marker)):
            if code_free_text[code_index] != "\n":
                code_free_text[code_index] = " "
        index = closing + len(marker)
    return "".join(code_free_text)


def inline_link_targets(markdown_text):
    targets = []
    for match in re.finditer(r"\]\(", markdown_text):
        target_start = match.end()
        while target_start < len(markdown_text) and markdown_text[target_start] in " \t":
            target_start += 1
        if target_start >= len(markdown_text):
            continue

        if markdown_text[target_start] == "<":
            target_end = markdown_text.find(">", target_start + 1)
            if target_end >= 0:
                targets.append((target_start + 1, target_end, markdown_text[target_start + 1 : target_end]))
            continue

        cursor = target_start
        nested_parentheses = 0
        escaped = False
        while cursor < len(markdown_text):
            character = markdown_text[cursor]
            if escaped:
                escaped = False
            elif character == "\\":
                escaped = True
            elif character == "(":
                nested_parentheses += 1
            elif character == ")":
                if nested_parentheses == 0:
                    break
                nested_parentheses -= 1
            elif character in " \t\n" and nested_parentheses == 0:
                break
            cursor += 1
        if cursor > target_start:
            targets.append((target_start, cursor, markdown_text[target_start:cursor]))
    return targets


def markdown_targets(markdown_text):
    targets = inline_link_targets(markdown_text)
    for match in re.finditer(r"(?m)^[ ]{0,3}\[[^]\n]+\]:[ \t]*(?:<([^>]+)>|([^\s]+))", markdown_text):
        target = match.group(1) or match.group(2)
        targets.append((match.start(), match.end(), target))

    occupied_ranges = [(start, end) for start, end, _ in targets]
    first_party_url = re.compile(
        r"https://(?:raw\.githubusercontent\.com/knews2019/skill-do-work/|"
        r"github\.com/knews2019/skill-do-work/)[^\s<>]+"
    )
    for match in first_party_url.finditer(markdown_text):
        if any(start <= match.start() < end for start, end in occupied_ranges):
            continue
        target = match.group(0).rstrip("),;:")
        targets.append((match.start(), match.start() + len(target), target))
    return sorted(targets)


def is_dynamic_target(target):
    return any(marker in target for marker in ("{", "}", "<", ">", "$"))


def resolve_installed_target(installed_path, modules):
    for source_root, destination_root in modules:
        try:
            relative_path = installed_path.relative_to(destination_root)
        except ValueError:
            continue
        return source_root / relative_path
    return None


def path_is_tracked(relative_path, tracked_paths):
    normalized = relative_path.as_posix()
    return normalized in tracked_paths


def path_is_exported(relative_path):
    path_parts = relative_path.parts
    for depth in range(1, len(path_parts) + 1):
        candidate = pathlib.PurePosixPath(*path_parts[:depth]).as_posix()
        attribute = git_output("check-attr", "export-ignore", "--", candidate).decode("utf-8", "replace")
        if attribute.rstrip().endswith(": set"):
            return False
    return True


def validate_first_party_url(target, tracked_paths):
    parsed = urllib.parse.urlsplit(target)
    decoded_path = urllib.parse.unquote(parsed.path)
    raw_prefix = f"/{repository_slug}/main/"
    blob_prefix = f"/{repository_slug}/blob/main/"

    if parsed.netloc == "raw.githubusercontent.com" and not decoded_path.startswith(raw_prefix):
        return "raw target must use the canonical repository main branch"
    if parsed.netloc == "raw.githubusercontent.com":
        repository_path = pathlib.PurePosixPath(decoded_path[len(raw_prefix) :])
        if not path_is_tracked(repository_path, tracked_paths):
            return "raw target is not a tracked live file"
        if not path_is_exported(repository_path):
            return "raw target is export-ignored"
    elif parsed.netloc == "github.com" and decoded_path.startswith(f"/{repository_slug}/blob/") and not decoded_path.startswith(blob_prefix):
        return "blob target must use the canonical repository main branch"
    elif parsed.netloc == "github.com" and decoded_path.startswith(blob_prefix):
        repository_path = pathlib.PurePosixPath(decoded_path[len(blob_prefix) :])
        if not path_is_tracked(repository_path, tracked_paths):
            return "blob target is not tracked"
    return None


modules = read_manifest()
tracked_paths = {
    os.fsdecode(path_bytes)
    for path_bytes in git_output("ls-files", "-z").split(b"\0")
    if path_bytes
}

markdown_paths = []
for source_root, _ in modules:
    source_prefix = source_root.as_posix() + "/"
    markdown_paths.extend(
        pathlib.PurePosixPath(path)
        for path in tracked_paths
        if path.startswith(source_prefix) and path.lower().endswith(".md")
    )

broken_references = 0
for markdown_path in sorted(set(markdown_paths), key=lambda path: path.as_posix()):
    source_file = repo_root / markdown_path
    try:
        markdown_text = strip_markdown_code(source_file.read_text(encoding="utf-8"))
    except (OSError, UnicodeError) as error:
        fail(f"cannot inspect {markdown_path}: {error}")
        broken_references += 1
        continue

    owning_module = next(
        (module for module in modules if markdown_path.is_relative_to(module[0])),
        None,
    )
    if owning_module is None:
        continue
    source_root, destination_root = owning_module
    installed_file = destination_root / markdown_path.relative_to(source_root)

    for target_start, _, target in markdown_targets(markdown_text):
        line_number = markdown_text.count("\n", 0, target_start) + 1
        target = target.replace("\\ ", " ")
        if not target or target.startswith("#") or is_dynamic_target(target):
            continue

        parsed = urllib.parse.urlsplit(target)
        if parsed.scheme:
            url_error = validate_first_party_url(target, tracked_paths)
            if url_error:
                fail(f"{markdown_path}:{line_number}: {url_error}: {target}")
                broken_references += 1
            continue
        if target.startswith("/"):
            continue

        decoded_target = urllib.parse.unquote(parsed.path)
        if not decoded_target:
            continue
        relative_target = pathlib.PurePosixPath(decoded_target)
        source_target = pathlib.PurePosixPath(os.path.normpath((markdown_path.parent / relative_target).as_posix()))
        installed_target = pathlib.PurePosixPath(os.path.normpath((installed_file.parent / relative_target).as_posix()))
        installed_source_target = resolve_installed_target(installed_target, modules)

        missing_locations = []
        if not (repo_root / source_target).exists():
            missing_locations.append("source")
        if installed_source_target is None or not (repo_root / installed_source_target).exists():
            missing_locations.append("installed")
        if missing_locations:
            fail(
                f"{markdown_path}:{line_number}: relative target is missing in "
                f"{' and '.join(missing_locations)} topology: {target}"
            )
            broken_references += 1

root_changelog = repo_root / "CHANGELOG.md"
installed_changelog = repo_root / "skills/do-work/CHANGELOG.md"
changelog_mismatch = root_changelog.read_bytes() != installed_changelog.read_bytes()
if changelog_mismatch:
    fail("skills/do-work/CHANGELOG.md is not byte-identical to root CHANGELOG.md")

if broken_references or changelog_mismatch:
    print(
        "shipped package reference contract: FAIL "
        f"({broken_references} broken reference(s), "
        f"changelog mirror {'differs' if changelog_mismatch else 'matches'})",
        file=sys.stderr,
    )
    raise SystemExit(1)

print("shipped package reference contract: PASS")
PY
