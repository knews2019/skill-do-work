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


def mask_text_range(character_list, range_start, range_end):
    for character_index in range(range_start, range_end):
        if character_list[character_index] != "\n":
            character_list[character_index] = " "


def leading_indentation(line):
    indentation_columns = 0
    content_start = 0
    while content_start < len(line) and line[content_start] in " \t":
        if line[content_start] == "\t":
            indentation_columns += 4 - (indentation_columns % 4)
        else:
            indentation_columns += 1
        content_start += 1
    return content_start, indentation_columns


def markdown_column_width(text):
    column_width = 0
    for character in text:
        if character == "\t":
            column_width += 4 - (column_width % 4)
        else:
            column_width += 1
    return column_width


def fence_info_string_is_valid(line, fence_match, fence_group):
    fence = fence_match.group(fence_group)
    fence_info_string = line[fence_match.end(fence_group) :].rstrip("\r\n")
    return fence[0] == "~" or "`" not in fence_info_string


def line_opens_paragraph(line, content_start):
    content = line[content_start:].rstrip("\r\n")
    list_item_match = re.match(
        r"^(?:[-+*]|[0-9]{1,9}[.)])[ \t]+(.*)$", content
    )
    if list_item_match:
        content = list_item_match.group(1)
    if not content or content.startswith("<!--"):
        return False
    fence_match = re.match(r"^(?:`{3,}|~{3,})", content)
    if fence_match and fence_info_string_is_valid(content, fence_match, 0):
        return False
    if re.match(
        r"^(?:#{1,6}(?:[ \t]|$)|>|(?:[-+*]|[0-9]{1,9}[.)])(?:[ \t]|$)|\[[^\]\n]+\]:)",
        content,
    ):
        return False
    if re.fullmatch(r"(?:\*[ \t]*){3,}|(?:_[ \t]*){3,}|(?:-[ \t]*){3,}", content):
        return False
    if re.fullmatch(r"(?:=+|-+)[ \t]*", content):
        return False
    return True


def inline_link_region(markdown_text, label_start):
    opening_bracket_escaped = punctuation_is_escaped(markdown_text, label_start)
    label_end = label_start + 1
    nested_labels = 0
    escaped_label_end = None
    while label_end < len(markdown_text):
        if markdown_text[label_end] == "[" and not punctuation_is_escaped(markdown_text, label_end):
            nested_labels += 1
        elif markdown_text[label_end] == "]":
            closing_bracket_escaped = punctuation_is_escaped(markdown_text, label_end)
            if nested_labels == 0:
                if not closing_bracket_escaped:
                    break
                candidate_parenthesis = label_end + 1
                while (
                    candidate_parenthesis < len(markdown_text)
                    and markdown_text[candidate_parenthesis] == "\\"
                ):
                    candidate_parenthesis += 1
                if (
                    escaped_label_end is None
                    and candidate_parenthesis < len(markdown_text)
                    and markdown_text[candidate_parenthesis] == "("
                ):
                    escaped_label_end = label_end
            elif not closing_bracket_escaped:
                nested_labels -= 1
        elif markdown_text[label_end] == "\n":
            break
        label_end += 1
    if label_end >= len(markdown_text) or markdown_text[label_end] == "\n":
        if escaped_label_end is None:
            return None
        label_end = escaped_label_end

    opening_parenthesis = label_end + 1
    has_backslash_separator = False
    while (
        opening_parenthesis < len(markdown_text)
        and markdown_text[opening_parenthesis] == "\\"
    ):
        has_backslash_separator = True
        opening_parenthesis += 1
    if (
        opening_parenthesis >= len(markdown_text)
        or markdown_text[opening_parenthesis] != "("
    ):
        return None
    structural_delimiter_escaped = (
        opening_bracket_escaped
        or punctuation_is_escaped(markdown_text, label_end)
        or has_backslash_separator
        or punctuation_is_escaped(markdown_text, opening_parenthesis)
    )

    cursor = opening_parenthesis + 1
    nested_parentheses = 0
    while cursor < len(markdown_text):
        character = markdown_text[cursor]
        punctuation_escaped = punctuation_is_escaped(markdown_text, cursor)
        if character == "(" and not punctuation_escaped:
            nested_parentheses += 1
        elif character == ")" and not punctuation_escaped:
            if nested_parentheses == 0:
                return cursor + 1, structural_delimiter_escaped
            nested_parentheses -= 1
        elif character == "\n":
            return None
        cursor += 1
    return None


def escaped_reference_definition_end(markdown_text, line_start, line_end):
    line = markdown_text[line_start:line_end]
    indentation_length = len(line) - len(line.lstrip(" "))
    if indentation_length > 3:
        return None

    label_start = line_start + indentation_length
    while label_start < line_end and markdown_text[label_start] == "\\":
        label_start += 1
    if label_start >= line_end or markdown_text[label_start] != "[":
        return None

    opening_bracket_escaped = punctuation_is_escaped(markdown_text, label_start)
    label_end = label_start + 1
    while label_end + 1 < line_end:
        if markdown_text[label_end] == "]" and markdown_text[label_end + 1] == ":":
            closing_bracket_escaped = punctuation_is_escaped(markdown_text, label_end)
            if opening_bracket_escaped or closing_bracket_escaped:
                return line_end
            return None
        label_end += 1
    return None


def strip_markdown_code(markdown_text):
    code_free_text = list(markdown_text)
    fence_character = None
    fence_length = 0
    list_fence_indentation = None
    paragraph_active = False
    line_start = 0

    for line in markdown_text.splitlines(keepends=True):
        line_end = line_start + len(line)
        fence_match = re.match(r"^[ ]{0,3}(`{3,}|~{3,})", line)
        if fence_match and not fence_info_string_is_valid(line, fence_match, 1):
            fence_match = None
        content_start, indentation_columns = leading_indentation(line)
        list_fence_match = re.match(
            r"^([ \t]*)(?:[-+*]|[0-9]{1,9}[.)])([ \t]+)(`{3,}|~{3,})",
            line,
        )
        if list_fence_match:
            list_fence_is_indented_code = (
                indentation_columns >= 4 and not paragraph_active
            )
            if list_fence_is_indented_code or not fence_info_string_is_valid(
                line, list_fence_match, 3
            ):
                list_fence_match = None
        blank_line = not line[content_start:].strip("\r\n")
        line_is_code = False

        if (
            fence_character is not None
            and list_fence_indentation is not None
            and not blank_line
            and indentation_columns < list_fence_indentation
        ):
            fence_character = None
            fence_length = 0
            list_fence_indentation = None
            paragraph_active = False

        if fence_character is not None:
            line_is_code = True
            closes_fence = False
            if list_fence_indentation is None:
                root_closing_fence_match = re.match(
                    r"^[ ]{0,3}(`{3,}|~{3,})[ \t]*(?:\r?\n)?$", line
                )
                if root_closing_fence_match:
                    fence = root_closing_fence_match.group(1)
                    closes_fence = (
                        fence[0] == fence_character and len(fence) >= fence_length
                    )
            else:
                closing_fence_match = re.match(
                    r"^([ \t]*)(`{3,}|~{3,})[ \t]*(?:\r?\n)?$", line
                )
                if closing_fence_match:
                    closing_indentation = markdown_column_width(
                        closing_fence_match.group(1)
                    )
                    closing_fence = closing_fence_match.group(2)
                    closes_fence = (
                        list_fence_indentation
                        <= closing_indentation
                        <= list_fence_indentation + 3
                        and closing_fence[0] == fence_character
                        and len(closing_fence) >= fence_length
                    )
            if closes_fence:
                fence_character = None
                fence_length = 0
                paragraph_active = list_fence_indentation is not None
                list_fence_indentation = None
            else:
                paragraph_active = False
        elif fence_match:
            line_is_code = True
            fence = fence_match.group(1)
            fence_character = fence[0]
            fence_length = len(fence)
            paragraph_active = False
        elif list_fence_match:
            line_is_code = True
            fence = list_fence_match.group(3)
            fence_character = fence[0]
            fence_length = len(fence)
            list_fence_indentation = markdown_column_width(
                line[: list_fence_match.start(3)]
            )
            paragraph_active = False
        elif indentation_columns >= 4 and not paragraph_active:
            line_is_code = True
            paragraph_active = False
        elif blank_line:
            paragraph_active = False
        elif indentation_columns >= 4:
            paragraph_active = True
        else:
            paragraph_active = line_opens_paragraph(line, content_start)

        if line_is_code:
            mask_text_range(code_free_text, line_start, line_end)
        line_start = line_end

    rendered_text = "".join(code_free_text)
    index = 0
    while index < len(code_free_text):
        if rendered_text.startswith("<!--", index):
            comment_end = rendered_text.find("-->", index + 4)
            masked_end = len(rendered_text) if comment_end < 0 else comment_end + 3
            mask_text_range(code_free_text, index, masked_end)
            index = masked_end
            continue
        if rendered_text[index] != "`" or punctuation_is_escaped(rendered_text, index):
            index += 1
            continue
        run_end = index
        while run_end < len(rendered_text) and rendered_text[run_end] == "`":
            run_end += 1
        marker_length = run_end - index
        closing_start = None
        search_index = run_end
        while search_index < len(rendered_text):
            if (
                rendered_text[search_index] != "`"
                or punctuation_is_escaped(rendered_text, search_index)
            ):
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
        mask_text_range(code_free_text, index, masked_end)
        index = masked_end

    rendered_text = "".join(code_free_text)
    label_start = 0
    while label_start < len(rendered_text):
        label_start = rendered_text.find("[", label_start)
        if label_start < 0:
            break
        link_region = inline_link_region(rendered_text, label_start)
        if link_region is None:
            label_start += 1
            continue
        region_end, structural_delimiter_escaped = link_region
        if structural_delimiter_escaped:
            mask_text_range(code_free_text, label_start, region_end)
        label_start = region_end

    rendered_text = "".join(code_free_text)
    line_start = 0
    for line in rendered_text.splitlines(keepends=True):
        line_end = line_start + len(line)
        masked_end = escaped_reference_definition_end(
            rendered_text, line_start, line_end
        )
        if masked_end is not None:
            mask_text_range(code_free_text, line_start, masked_end)
        line_start = line_end
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


# A heading anchor is generated from the heading's *rendered* text: lowercase it, drop
# every character that is not a word character, a hyphen, or a space, then turn spaces
# into hyphens. Repeated slugs take a -1, -2 … suffix in document order.
inline_link_label_pattern = re.compile(r"!?\[([^\[\]]*)\](?:\([^()]*\)|\[[^\[\]]*\])")
atx_heading_pattern = re.compile(r"^[ ]{0,3}(#{1,6})(?:[ \t]+(.*?))?[ \t]*$")
heading_slug_cache = {}


def heading_rendered_text(heading_text):
    previous_text = None
    while previous_text != heading_text:
        previous_text = heading_text
        heading_text = inline_link_label_pattern.sub(r"\1", heading_text)
    return heading_text


def heading_anchor_slug(heading_text):
    slug = heading_rendered_text(heading_text).strip().lower()
    slug = re.sub(r"[^\w\- ]", "", slug)
    return slug.replace(" ", "-")


def heading_anchor_slugs_from_text(markdown_text):
    # strip_markdown_code preserves byte and line offsets, so a line whose masked form no
    # longer opens with # sat inside a fence, an indented block, or an HTML comment and is
    # not a heading. Slugs come from the raw line, because masking blanks inline code that
    # the rendered heading text keeps.
    masked_lines = strip_markdown_code(markdown_text).splitlines()
    anchor_slugs = set()
    slug_occurrences = {}
    for line_index, raw_line in enumerate(markdown_text.splitlines()):
        if line_index >= len(masked_lines):
            break
        if not masked_lines[line_index].lstrip(" \t").startswith("#"):
            continue
        heading_match = atx_heading_pattern.match(raw_line)
        if heading_match is None:
            continue
        heading_text = re.sub(r"[ \t]+#+$", "", heading_match.group(2) or "")
        anchor_slug = heading_anchor_slug(heading_text)
        if not anchor_slug:
            continue
        occurrence_index = slug_occurrences.get(anchor_slug, 0)
        slug_occurrences[anchor_slug] = occurrence_index + 1
        anchor_slugs.add(
            anchor_slug if occurrence_index == 0 else f"{anchor_slug}-{occurrence_index}"
        )
    return anchor_slugs


def heading_anchor_slugs(markdown_file):
    cache_key = os.fspath(markdown_file)
    if cache_key not in heading_slug_cache:
        try:
            markdown_text = markdown_file.read_text(encoding="utf-8")
        except (OSError, UnicodeError):
            heading_slug_cache[cache_key] = None
        else:
            heading_slug_cache[cache_key] = heading_anchor_slugs_from_text(markdown_text)
    return heading_slug_cache[cache_key]


def run_anchor_slug_fixtures():
    anchor_cases = [
        (
            "case folding and punctuation removal",
            "## Portfolio summary publication\n## Raw text before shell quoting\n",
            {"portfolio-summary-publication", "raw-text-before-shell-quoting"},
        ),
        (
            "inline code contributes its content, not its backticks",
            "## What `run` does (and does not) do\n",
            {"what-run-does-and-does-not-do"},
        ),
        (
            "link labels replace their destinations",
            "## See [the guide](../docs/guide.md#top) first\n",
            {"see-the-guide-first"},
        ),
        (
            "code blocks and comments hold no headings",
            "```sh\n# fenced comment\n```\n"
            "    # indented comment\n"
            "<!-- # commented heading -->\n"
            "# Real Heading\n",
            {"real-heading"},
        ),
        (
            "repeated headings take numbered suffixes",
            "# Same\n\n# Same\n\n# Same\n",
            {"same", "same-1", "same-2"},
        ),
        (
            "closing hash sequences are not part of the anchor",
            "## Trailing Hashes ##\n",
            {"trailing-hashes"},
        ),
    ]

    fixture_failures = 0
    for fixture_name, fixture_markdown, expected_slugs in anchor_cases:
        actual_slugs = heading_anchor_slugs_from_text(fixture_markdown)
        if actual_slugs != expected_slugs:
            fail(
                f"anchor fixture {fixture_name!r}: expected {sorted(expected_slugs)!r}, "
                f"got {sorted(actual_slugs)!r}"
            )
            fixture_failures += 1

    if fixture_failures:
        raise SystemExit(1)


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
            "root backtick info string Goldmark differential",
            "```lang`invalid\n[live](visible.md)\n",
            ["visible.md"],
        ),
        (
            "invalid root backtick info string starts a paragraph",
            "```lang`invalid\n    [live](live-invalid-root-continuation.md)\n",
            ["live-invalid-root-continuation.md"],
        ),
        (
            "root tilde info strings may contain backticks",
            "~~~lang`valid\n"
            "[hidden](missing-tilde-info.md)\n"
            "~~~\n"
            "[live](live-after-tilde-info.md)\n",
            ["live-after-tilde-info.md"],
        ),
        (
            "root fence lookalike lines",
            "```markdown\n"
            "[hidden](missing-before-lookalike.md)\n"
            "``` trailing text\n"
            "[hidden](missing-after-lookalike.md)\n"
            "```\n"
            "[live](live-after-root-fence.md)\n",
            ["live-after-root-fence.md"],
        ),
        (
            "indented code",
            "    [hidden](missing-indent.md)\n\t[also hidden](missing-tab.md)\n   [live](live.md)\n",
            ["live.md"],
        ),
        (
            "effective-column indented code and paragraph continuation",
            " \t[hidden](missing-one-space-tab.md)\n"
            "  \t[hidden](missing-two-space-tab.md)\n"
            "   \t[hidden](missing-three-space-tab.md)\n\n"
            "active paragraph\n"
            "    [live](live-continuation.md)\n\n"
            "    [hidden](missing-code-block.md)\n",
            ["live-continuation.md"],
        ),
        (
            "bullet-list paragraph continuations",
            "- dash paragraph\n"
            "    [live](live-dash-continuation.md)\n"
            "+ plus paragraph\n"
            "    [live](live-plus-continuation.md)\n"
            "* star paragraph\n"
            "    [live](live-star-continuation.md)\n\n"
            "-\n"
            "    [hidden](missing-empty-bullet-code.md)\n\n"
            "- ```\n"
            "    [hidden](missing-list-fence-code.md)\n\n"
            "- paragraph before blank\n\n"
            "    [hidden](missing-blank-separated-bullet-code.md)\n",
            [
                "live-dash-continuation.md",
                "live-plus-continuation.md",
                "live-star-continuation.md",
            ],
        ),
        (
            "ordered-list paragraph continuations",
            "1. period paragraph\n"
            "    [live](live-period-continuation.md)\n"
            "2) parenthesis paragraph\n"
            "    [live](live-parenthesis-continuation.md)\n"
            "123456789. nine-digit paragraph\n"
            "    [live](live-nine-digit-continuation.md)\n\n"
            "1.\n"
            "    [hidden](missing-empty-ordered-code.md)\n\n"
            "3) paragraph before blank\n\n"
            "    [hidden](missing-blank-separated-ordered-code.md)\n",
            [
                "live-period-continuation.md",
                "live-parenthesis-continuation.md",
                "live-nine-digit-continuation.md",
            ],
        ),
        (
            "list-item fences with attached info strings",
            "- ```markdown\n"
            "  [hidden](missing-bullet-backtick-fence.md)\n"
            "  ```\n"
            "  [live](live-bullet-after-fence.md)\n"
            "1. ~~~text\n"
            "   [hidden](missing-ordered-tilde-fence.md)\n"
            "   ~~~\n"
            "   [live](live-ordered-after-fence.md)\n"
            "123456789) ```text\n"
            "           [hidden](missing-nine-digit-backtick-fence.md)\n"
            "           ```\n"
            "           [live](live-nine-digit-after-fence.md)\n"
            "  - ~~~markdown\n"
            "    [hidden](missing-nested-tilde-fence.md)\n"
            "    ~~~\n"
            "    [live](live-nested-after-fence.md)\n"
            "- ```lang`invalid\n"
            "    [live](live-invalid-list-fence-info.md)\n"
            "- paragraph\n"
            "    [live](live-list-paragraph-continuation.md)\n\n"
            "    - ```not-a-list-fence\n"
            "    [hidden](missing-indented-list-shaped-code.md)\n\n"
            "[live](live-after-indented-code.md)\n",
            [
                "live-bullet-after-fence.md",
                "live-ordered-after-fence.md",
                "live-nine-digit-after-fence.md",
                "live-nested-after-fence.md",
                "live-invalid-list-fence-info.md",
                "live-list-paragraph-continuation.md",
                "live-after-indented-code.md",
            ],
        ),
        (
            "unclosed list-item fences end at their container boundary",
            "- ```markdown\n"
            "  [hidden](missing-unclosed-bullet-fence.md)\n"
            "[live](live-after-unclosed-bullet-fence.md)\n"
            "1. ~~~text\n"
            "   [hidden](missing-unclosed-ordered-fence.md)\n"
            "[live](live-after-unclosed-ordered-fence.md)\n"
            "  - ```nested\n"
            "    [hidden](missing-unclosed-nested-fence.md)\n"
            "[live](live-after-unclosed-nested-fence.md)\n"
            "- ```markdown\n"
            "  [hidden](missing-over-indented-closer-content.md)\n"
            "      ```\n"
            "[live](live-after-over-indented-closer.md)\n",
            [
                "live-after-unclosed-bullet-fence.md",
                "live-after-unclosed-ordered-fence.md",
                "live-after-unclosed-nested-fence.md",
                "live-after-over-indented-closer.md",
            ],
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
            "escaped inline code delimiters",
            r"\`[live](live-escaped-backtick.md)\` "
            r"\\`[hidden](missing-even-backslash.md)\\` "
            "`[hidden](missing-ordinary-code.md)` "
            "``[hidden](missing-exact-run-code.md)``\n",
            ["live-escaped-backtick.md"],
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
            "escaped link regions",
            r"\[hidden](missing-escaped-relative.md) "
            r"\[hidden](https://raw.githubusercontent.com/knews2019/skill-do-work/main/missing-escaped-url.md) "
            r"\\[live](live-even-relative.md) "
            r"\\[live](https://raw.githubusercontent.com/knews2019/skill-do-work/main/VERSION)"
            + "\n",
            [
                "live-even-relative.md",
                "https://raw.githubusercontent.com/knews2019/skill-do-work/main/VERSION",
            ],
        ),
        (
            "odd-parity escaped closing-bracket inline links",
            r"[hidden\](https://raw.githubusercontent.com/knews2019/skill-do-work/main/missing-escaped-closing-bracket.md) "
            r"[live\\](https://raw.githubusercontent.com/knews2019/skill-do-work/main/VERSION) "
            r"[live \] label](https://raw.githubusercontent.com/knews2019/skill-do-work/main/VERSION) "
            r"[live [nested]](https://raw.githubusercontent.com/knews2019/skill-do-work/main/VERSION)"
            + "\n"
            + r"[incomplete\](https://raw.githubusercontent.com/knews2019/skill-do-work/main/VERSION"
            + "\n",
            [
                "https://raw.githubusercontent.com/knews2019/skill-do-work/main/VERSION",
                "https://raw.githubusercontent.com/knews2019/skill-do-work/main/VERSION",
                "https://raw.githubusercontent.com/knews2019/skill-do-work/main/VERSION",
                "https://raw.githubusercontent.com/knews2019/skill-do-work/main/VERSION",
            ],
        ),
        (
            "backslash-separated destination openers are not inline links",
            r"[hidden]\(https://raw.githubusercontent.com/knews2019/skill-do-work/main/missing-escaped-opening-parenthesis.md) "
            r"[hidden]\\(https://raw.githubusercontent.com/knews2019/skill-do-work/main/missing-even-opening-parenthesis.md) "
            r"[ordinary](https://raw.githubusercontent.com/knews2019/skill-do-work/main/VERSION)"
            + "\n",
            [
                "https://raw.githubusercontent.com/knews2019/skill-do-work/main/VERSION",
            ],
        ),
        (
            "destination-opening parenthesis adjacency",
            r"[zero](zero-relative.md) "
            r"[one]\(missing-one-backslash-relative.md) "
            r"[two]\\(missing-two-backslash-relative.md) "
            r"[three]\\\(missing-three-backslash-relative.md) "
            r"[four]\\\\(missing-four-backslash-relative.md)"
            + "\n",
            [
                "zero-relative.md",
            ],
        ),
        (
            "escaped opening brackets inside live labels",
            r"\[hidden](missing-escaped-outer.md) "
            r"[\[live](live-escaped-opening-label.md) "
            r"[prefix \[content](live-escaped-opening-content.md)"
            + "\n",
            [
                "live-escaped-opening-label.md",
                "live-escaped-opening-content.md",
            ],
        ),
        (
            "escaped reference syntax",
            r"\[hidden]: missing-open.md" + "\n" + r"[hidden\]: missing-close.md" + "\n[live]: live.md\n",
            ["live.md"],
        ),
        (
            "escaped reference-definition regions",
            r"\[hidden]: https://raw.githubusercontent.com/knews2019/skill-do-work/main/missing-open-reference.md"
            + "\n"
            + r"[hidden\]: https://raw.githubusercontent.com/knews2019/skill-do-work/main/missing-close-reference.md"
            + "\n"
            + r"\\[live]: https://raw.githubusercontent.com/knews2019/skill-do-work/main/VERSION"
            + "\n"
            + r"[live\\]: https://raw.githubusercontent.com/knews2019/skill-do-work/main/VERSION"
            + "\n[live]: https://raw.githubusercontent.com/knews2019/skill-do-work/main/VERSION\n",
            [
                "https://raw.githubusercontent.com/knews2019/skill-do-work/main/VERSION",
                "https://raw.githubusercontent.com/knews2019/skill-do-work/main/VERSION",
                "https://raw.githubusercontent.com/knews2019/skill-do-work/main/VERSION",
            ],
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
run_anchor_slug_fixtures()
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
            continue

        # Anchor resolution is keyed on a condition, never on a list of files or link
        # forms: a link is anchor-checked exactly when it survived the skips above and
        # carries both a fragment and a Markdown target. Everything skipped is skipped
        # because this checker cannot resolve it without leaving the repository or
        # guessing — a scheme'd URL (http(s), mailto, anything else) is not read; a
        # bare fragment names no file to resolve against; a root-absolute or templated
        # target has no single on-disk meaning; and a heading only exists in Markdown,
        # so a fragment on any other suffix is an application's own addressing scheme.
        # Anchors resolve against ATX headings, which is what this repository writes; a
        # heading declared any other way reports here as a missing anchor.
        anchor = urllib.parse.unquote(parsed.fragment)
        if not anchor or relative_target.suffix.lower() != ".md":
            continue

        # The two topologies usually resolve back to the same source file; check each
        # distinct one so a link that lands somewhere else once installed is still read.
        anchor_targets = []
        for topology_name, resolved_target in (
            ("source", source_target),
            ("installed", installed_source_target),
        ):
            if resolved_target in anchor_targets:
                continue
            anchor_targets.append(resolved_target)
            target_slugs = heading_anchor_slugs(repo_root / resolved_target)
            if target_slugs is None:
                fail(
                    f"{markdown_path}:{line_number}: cannot read {topology_name} "
                    f"topology target {resolved_target}: {target}"
                )
                broken_references += 1
            elif anchor not in target_slugs:
                fail(
                    f"{markdown_path}:{line_number}: anchor #{anchor} is not a heading "
                    f"in {topology_name} topology target {resolved_target}: {target}"
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
