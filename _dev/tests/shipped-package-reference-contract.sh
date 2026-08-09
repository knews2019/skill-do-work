#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"

python3 - "$repo_root" <<'PY'
import csv
import os
import pathlib
import re
import string
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
    code_free_text = list(markdown_text)
    fence_character = None
    fence_length = 0
    line_start = 0

    for line in markdown_text.splitlines(keepends=True):
        line_end = line_start + len(line)
        fence_match = re.match(r"^[ ]{0,3}(`{3,}|~{3,})", line)
        indented_code = fence_character is None and (
            line.startswith("\t") or line.startswith("    ")
        )
        if indented_code:
            pass
        elif fence_match:
            fence = fence_match.group(1)
            if fence_character is None:
                fence_character = fence[0]
                fence_length = len(fence)
            elif fence[0] == fence_character and len(fence) >= fence_length:
                fence_character = None
                fence_length = 0
        elif fence_character is None:
            line_start = line_end
            continue
        for code_index in range(line_start, line_end):
            if code_free_text[code_index] != "\n":
                code_free_text[code_index] = " "
        line_start = line_end

    rendered_text = "".join(code_free_text)
    index = 0
    while index < len(code_free_text):
        if rendered_text.startswith("<!--", index):
            comment_end = rendered_text.find("-->", index + 4)
            masked_end = len(rendered_text) if comment_end < 0 else comment_end + 3
            for code_index in range(index, masked_end):
                if code_free_text[code_index] != "\n":
                    code_free_text[code_index] = " "
            index = masked_end
            continue
        if rendered_text[index] != "`":
            index += 1
            continue
        run_end = index
        while run_end < len(rendered_text) and rendered_text[run_end] == "`":
            run_end += 1
        marker_length = run_end - index
        closing_start = None
        search_index = run_end
        while search_index < len(rendered_text):
            if rendered_text[search_index] != "`":
                search_index += 1
                continue
            search_end = search_index
            while search_end < len(rendered_text) and rendered_text[search_end] == "`":
                search_end += 1
            if search_end - search_index == marker_length:
                closing_start = search_index
                break
            search_index = search_end
        if closing_start is None:
            index = run_end
            continue
        masked_end = closing_start + marker_length
        for code_index in range(index, masked_end):
            if code_free_text[code_index] != "\n":
                code_free_text[code_index] = " "
        index = masked_end
    return "".join(code_free_text)


def punctuation_is_escaped(markdown_text, punctuation_index):
    preceding_backslashes = 0
    cursor = punctuation_index - 1
    while cursor >= 0 and markdown_text[cursor] == "\\":
        preceding_backslashes += 1
        cursor -= 1
    return preceding_backslashes % 2 == 1


def inline_link_targets(markdown_text):
    targets = []
    label_start = 0
    while label_start < len(markdown_text):
        label_start = markdown_text.find("[", label_start)
        if label_start < 0:
            break
        if punctuation_is_escaped(markdown_text, label_start):
            label_start += 1
            continue

        label_end = label_start + 1
        nested_labels = 0
        while label_end < len(markdown_text):
            if markdown_text[label_end] == "[" and not punctuation_is_escaped(markdown_text, label_end):
                nested_labels += 1
            elif markdown_text[label_end] == "]" and not punctuation_is_escaped(markdown_text, label_end):
                if nested_labels == 0:
                    break
                nested_labels -= 1
            label_end += 1
        if label_end >= len(markdown_text):
            label_start += 1
            continue

        opening_parenthesis = label_end + 1
        if (
            opening_parenthesis >= len(markdown_text)
            or markdown_text[opening_parenthesis] != "("
            or punctuation_is_escaped(markdown_text, opening_parenthesis)
        ):
            label_start = label_end + 1
            continue

        target_start = opening_parenthesis + 1
        while target_start < len(markdown_text) and markdown_text[target_start] in " \t":
            target_start += 1
        if target_start >= len(markdown_text):
            break

        if markdown_text[target_start] == "<" and not punctuation_is_escaped(markdown_text, target_start):
            target_end = target_start + 1
            while target_end < len(markdown_text):
                if markdown_text[target_end] == "\n":
                    break
                if markdown_text[target_end] == ">" and not punctuation_is_escaped(markdown_text, target_end):
                    targets.append(
                        (target_start + 1, target_end, markdown_text[target_start + 1 : target_end])
                    )
                    break
                target_end += 1
            label_start = label_end + 1
            continue

        cursor = target_start
        nested_parentheses = 0
        while cursor < len(markdown_text):
            character = markdown_text[cursor]
            punctuation_escaped = punctuation_is_escaped(markdown_text, cursor)
            if character == "(" and not punctuation_escaped:
                nested_parentheses += 1
            elif character == ")" and not punctuation_escaped:
                if nested_parentheses == 0:
                    break
                nested_parentheses -= 1
            elif (
                character in " \t\n"
                and nested_parentheses == 0
                and not punctuation_escaped
            ):
                break
            cursor += 1
        if cursor > target_start:
            targets.append((target_start, cursor, markdown_text[target_start:cursor]))
        label_start = label_end + 1
    return targets


def markdown_targets(markdown_text):
    targets = inline_link_targets(markdown_text)
    line_start = 0
    for line in markdown_text.splitlines(keepends=True):
        definition_start = len(line) - len(line.lstrip(" "))
        if (
            definition_start > 3
            or definition_start >= len(line)
            or line[definition_start] != "["
        ):
            line_start += len(line)
            continue

        label_end = definition_start + 1
        while label_end < len(line):
            if line[label_end] == "]" and not punctuation_is_escaped(line, label_end):
                break
            label_end += 1
        if label_end >= len(line) or label_end + 1 >= len(line) or line[label_end + 1] != ":":
            line_start += len(line)
            continue

        target_start = label_end + 2
        while target_start < len(line) and line[target_start] in " \t":
            target_start += 1
        target_end = target_start
        if target_start < len(line) and line[target_start] == "<":
            target_start += 1
            target_end = target_start
            while target_end < len(line):
                if line[target_end] == ">" and not punctuation_is_escaped(line, target_end):
                    break
                if line[target_end] in "\r\n":
                    break
                target_end += 1
            if target_end >= len(line) or line[target_end] != ">":
                line_start += len(line)
                continue
        else:
            while target_end < len(line) and not line[target_end].isspace():
                target_end += 1
        if target_end > target_start:
            targets.append(
                (
                    line_start + target_start,
                    line_start + target_end,
                    line[target_start:target_end],
                )
            )
        line_start += len(line)

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


def normalize_markdown_target(target):
    normalized_target = []
    target_index = 0
    while target_index < len(target):
        if (
            target[target_index] == "\\"
            and target_index + 1 < len(target)
            and (target[target_index + 1] in string.punctuation or target[target_index + 1] == " ")
        ):
            normalized_target.append(target[target_index + 1])
            target_index += 2
            continue
        normalized_target.append(target[target_index])
        target_index += 1
    return "".join(normalized_target)


def run_parser_fixtures():
    parser_cases = [
        (
            "backtick fenced code",
            "[before](before.md)\n```markdown\n[hidden](missing-fence.md)\n```\n[after](after.md)\n",
            ["before.md", "after.md"],
        ),
        (
            "tilde fenced code",
            "~~~\n[hidden](missing-tilde.md)\n~~~\n[live](live.md)\n",
            ["live.md"],
        ),
        (
            "indented code",
            "    [hidden](missing-indent.md)\n\t[also hidden](missing-tab.md)\n   [live](live.md)\n",
            ["live.md"],
        ),
        (
            "HTML comments",
            "[before](before.md) <!-- [hidden](missing-inline-comment.md) -->\n<!--\n`comment backtick` [hidden](missing-block-comment.md)\n--> [after](after.md)\n",
            ["before.md", "after.md"],
        ),
        (
            "unterminated HTML comment",
            "[before](before.md)\n<!-- [hidden](missing-eof-comment.md)\n",
            ["before.md"],
        ),
        (
            "inline code delimiter runs",
            "`` code `tick` and ```longer``` [hidden](missing-code.md) `` [live](live.md)\n",
            ["live.md"],
        ),
        (
            "comment marker inside inline code",
            "`<!-- [hidden](missing-code-comment.md) -->` [live](live.md)\n",
            ["live.md"],
        ),
        (
            "escaped inline syntax",
            r"\[hidden](missing-open.md) [hidden\](missing-close.md) \\[live](even.md)" + "\n",
            ["even.md"],
        ),
        (
            "escaped reference syntax",
            r"\[hidden]: missing-open.md" + "\n" + r"[hidden\]: missing-close.md" + "\n[live]: live.md\n",
            ["live.md"],
        ),
        (
            "published link forms",
            "[inline](inline.md) ![image](image.png)\n[bare]: bare.md\n[angle]: <angle.md>\n",
            ["inline.md", "image.png", "bare.md", "angle.md"],
        ),
        (
            "balanced and escaped destinations",
            r"[nested](docs/name(v1).md) [escaped](docs/name\(v2\).md) "
            r"[closing](docs/name\).md) [even](docs/name\\(v3).md)"
            + "\n",
            [
                "docs/name(v1).md",
                "docs/name(v2).md",
                "docs/name).md",
                r"docs/name\(v3).md",
            ],
        ),
        (
            "escaped destination punctuation",
            r"[punctuation](docs/file\!\#\[draft\].md) [backslash](docs/path\\leaf.md)" + "\n",
            ["docs/file!#[draft].md", r"docs/path\leaf.md"],
        ),
    ]

    fixture_failures = 0
    for fixture_name, fixture_markdown, expected_targets in parser_cases:
        masked_markdown = strip_markdown_code(fixture_markdown)
        if len(masked_markdown) != len(fixture_markdown):
            fail(f"parser fixture {fixture_name!r}: masking changed source offsets")
            fixture_failures += 1
        if [
            character_index
            for character_index, character in enumerate(masked_markdown)
            if character == "\n"
        ] != [
            character_index
            for character_index, character in enumerate(fixture_markdown)
            if character == "\n"
        ]:
            fail(f"parser fixture {fixture_name!r}: masking changed line offsets")
            fixture_failures += 1
        actual_targets = [
            normalize_markdown_target(target)
            for _, _, target in markdown_targets(masked_markdown)
        ]
        if actual_targets != expected_targets:
            fail(
                f"parser fixture {fixture_name!r}: expected {expected_targets!r}, "
                f"got {actual_targets!r}"
            )
            fixture_failures += 1

    if fixture_failures:
        raise SystemExit(1)


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


run_parser_fixtures()
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
        target = normalize_markdown_target(target)
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
