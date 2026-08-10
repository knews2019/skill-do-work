#!/usr/bin/env bash
# Atomically create, append, or replace do-work's managed text section.
set -euo pipefail

if ! command -v python3 >/dev/null 2>&1; then
  printf 'replace-text-section: python3 is required\n' >&2
  exit 1
fi

exec python3 - "$@" <<'PY'
import argparse
import os
import re
import stat
import sys
import tempfile

BEGIN = b"# >>> do-work:recipes >>>"
END = b"# <<< do-work:recipes <<<"
JUST_IDENTIFIER = re.compile(rb"[A-Za-z_][A-Za-z0-9_-]*")
JUST_ALIAS = re.compile(rb"alias[ \t]+([A-Za-z_][A-Za-z0-9_-]*)[ \t]*:=")


def die(message: str) -> "None":
    print(f"replace-text-section: {message}", file=sys.stderr)
    raise SystemExit(1)


def read_regular(path: str, label: str) -> bytes:
    if os.path.islink(path):
        die(f"{label} must not be a symlink: {path}")
    try:
        metadata = os.stat(path)
    except OSError as error:
        die(f"cannot stat {label} {path}: {error}")
    if not stat.S_ISREG(metadata.st_mode):
        die(f"{label} must be a regular file: {path}")
    try:
        with open(path, "rb") as handle:
            return handle.read()
    except OSError as error:
        die(f"cannot read {label} {path}: {error}")


def lines_with_offsets(data: bytes):
    lines = data.splitlines(keepends=True)
    offsets = []
    cursor = 0
    for line in lines:
        offsets.append(cursor)
        cursor += len(line)
    return lines, offsets


def line_body(line: bytes) -> bytes:
    return line.rstrip(b"\r\n")


def marker_span(data: bytes, label: str, require_section_only: bool = False):
    lines, offsets = lines_with_offsets(data)
    bodies = [line_body(line) for line in lines]
    begin_indexes = [index for index, body in enumerate(bodies) if body == BEGIN]
    end_indexes = [index for index, body in enumerate(bodies) if body == END]

    if not begin_indexes and not end_indexes:
        return None
    if len(begin_indexes) != 1 or len(end_indexes) != 1:
        die(f"{label} must contain exactly one begin marker and one end marker")
    begin_index = begin_indexes[0]
    end_index = end_indexes[0]
    if begin_index >= end_index:
        die(f"{label} has reversed or nested managed markers")
    if require_section_only and (begin_index != 0 or end_index != len(lines) - 1):
        die(f"{label} must contain only the complete managed section")

    span_start = offsets[begin_index]
    span_end = offsets[end_index] + len(lines[end_index])
    return span_start, span_end


def just_definition_name(line: bytes):
    body = line_body(line)
    if not body or body[:1] in (b" ", b"\t", b"#"):
        return None

    alias_match = JUST_ALIAS.match(body)
    if alias_match:
        return alias_match.group(1)

    name_offset = 1 if body.startswith(b"@") else 0
    name_match = JUST_IDENTIFIER.match(body, name_offset)
    if not name_match:
        return None

    remainder = body[name_match.end() :]
    quote = None
    escaped = False
    for index, character in enumerate(remainder):
        if quote is not None:
            if escaped:
                escaped = False
            elif character == 92 and quote == 34:
                escaped = True
            elif character == quote:
                quote = None
            continue
        if character in (34, 39, 96):
            quote = character
        elif character == 58:
            if index + 1 < len(remainder) and remainder[index + 1] == 61:
                return None
            return name_match.group(0)
    return None


def just_multiline_string_state(line: bytes, active_delimiter):
    body = line_body(line)
    if active_delimiter is None and (
        not body or body[:1] in (b" ", b"\t", b"#")
    ):
        return None

    index = 0
    ordinary_quote = None
    escaped = False
    while index < len(body):
        if active_delimiter is not None:
            delimiter_matches = body.startswith(active_delimiter, index)
            if active_delimiter == b"```" and delimiter_matches:
                delimiter_matches = (
                    (index == 0 or body[index - 1] != 96)
                    and (
                        index + len(active_delimiter) == len(body)
                        or body[index + len(active_delimiter)] != 96
                    )
                )
            if delimiter_matches:
                if active_delimiter in (b'"', b'"""'):
                    backslash_count = 0
                    backslash_index = index - 1
                    while backslash_index >= 0 and body[backslash_index] == 92:
                        backslash_count += 1
                        backslash_index -= 1
                    if backslash_count % 2 == 1:
                        index += 1
                        continue
                delimiter_length = len(active_delimiter)
                active_delimiter = None
                index += delimiter_length
                continue
            index += 1
            continue

        character = body[index]
        if ordinary_quote is not None:
            if escaped:
                escaped = False
            elif character == 92 and ordinary_quote == 34:
                escaped = True
            elif character == ordinary_quote:
                ordinary_quote = None
            index += 1
            continue
        if character == 35:
            break
        if body.startswith(b"'''", index):
            active_delimiter = b"'''"
            index += 3
        elif body.startswith(b'"""', index):
            active_delimiter = b'"""'
            index += 3
        elif body.startswith(b"```", index) and (
            (index == 0 or body[index - 1] != 96)
            and (index + 3 == len(body) or body[index + 3] != 96)
        ):
            active_delimiter = b"```"
            index += 3
        elif character in (34, 39, 96):
            ordinary_quote = character
            index += 1
        else:
            index += 1

    if ordinary_quote in (34, 39, 96):
        return bytes((ordinary_quote,))
    return active_delimiter


def just_definition_names(data: bytes):
    definition_names = set()
    active_delimiter = None
    for line in data.splitlines(keepends=True):
        line_starts_in_multiline_string = active_delimiter is not None
        active_delimiter = just_multiline_string_state(line, active_delimiter)
        if line_starts_in_multiline_string:
            continue
        definition_name = just_definition_name(line)
        if definition_name is not None:
            definition_names.add(definition_name)
    return definition_names


def atomic_replace(path: str, content: bytes, mode: int) -> None:
    parent = os.path.dirname(os.path.abspath(path))
    if not os.path.isdir(parent):
        die(f"target parent directory does not exist: {parent}")
    descriptor = -1
    temporary_path = ""
    try:
        descriptor, temporary_path = tempfile.mkstemp(
            prefix=f".{os.path.basename(path)}.", suffix=".tmp", dir=parent
        )
        with os.fdopen(descriptor, "wb") as handle:
            descriptor = -1
            handle.write(content)
            handle.flush()
            os.fsync(handle.fileno())
            os.fchmod(handle.fileno(), mode)
        os.replace(temporary_path, path)
        temporary_path = ""
        directory_flags = os.O_RDONLY | getattr(os, "O_DIRECTORY", 0)
        try:
            directory_descriptor = os.open(parent, directory_flags)
            try:
                os.fsync(directory_descriptor)
            finally:
                os.close(directory_descriptor)
        except OSError:
            # The atomic rename has already succeeded. Some filesystems do not support
            # directory fsync; that durability limitation must not be reported as a failed
            # replacement after the visible target has changed.
            pass
    except OSError as error:
        die(f"atomic replacement failed for {path}: {error}")
    finally:
        if descriptor >= 0:
            os.close(descriptor)
        if temporary_path:
            try:
                os.unlink(temporary_path)
            except FileNotFoundError:
                pass


parser = argparse.ArgumentParser(add_help=False)
parser.add_argument("--target")
parser.add_argument("--section-file")
parser.add_argument("--template-file")
parser.add_argument("--reject-recipe-collisions", action="store_true")
parser.add_argument("--help", action="store_true")
try:
    arguments, residue = parser.parse_known_args()
except SystemExit:
    die("usage: replace-text-section.sh --target <path> --section-file <path> [--template-file <path>] [--reject-recipe-collisions]")
if arguments.help:
    sys.stdout.write("usage: replace-text-section.sh --target <path> --section-file <path> [--template-file <path>] [--reject-recipe-collisions]\n")
    raise SystemExit(0)
if residue or not arguments.target or not arguments.section_file:
    die("usage: replace-text-section.sh --target <path> --section-file <path> [--template-file <path>] [--reject-recipe-collisions]")

section_data = read_regular(arguments.section_file, "section file")
section_span = marker_span(section_data, "section file", require_section_only=True)
if section_span is None or not section_data.endswith(b"\n"):
    die("section file must be one newline-terminated managed section")

target_exists = os.path.exists(arguments.target) or os.path.islink(arguments.target)
if not target_exists:
    if not arguments.template_file:
        die("target is absent; --template-file is required")
    template_data = read_regular(arguments.template_file, "template file")
    template_span = marker_span(template_data, "template file")
    if template_span is None or template_data[template_span[0] : template_span[1]] != section_data:
        die("template file must contain the supplied managed section exactly once")
    template_mode = stat.S_IMODE(os.stat(arguments.template_file).st_mode)
    atomic_replace(arguments.target, template_data, template_mode)
    raise SystemExit(0)

if os.path.islink(arguments.target):
    die(f"target must not be a symlink: {arguments.target}")
target_data = read_regular(arguments.target, "target")
target_mode = stat.S_IMODE(os.stat(arguments.target).st_mode)
target_span = marker_span(target_data, "target")

if arguments.reject_recipe_collisions:
    reserved_recipe_names = just_definition_names(section_data)
    if not reserved_recipe_names:
        die("section file defines no Just recipes or aliases for collision validation")
    if target_span is None:
        unmanaged_target_data = target_data
    else:
        unmanaged_target_data = target_data[: target_span[0]] + target_data[target_span[1] :]
    collision_names = sorted(just_definition_names(unmanaged_target_data) & reserved_recipe_names)
    if collision_names:
        die(
            "target defines reserved Just recipe or alias outside managed section: "
            + ", ".join(name.decode("ascii") for name in collision_names)
        )

if target_span is not None:
    replacement_data = target_data[: target_span[0]] + section_data + target_data[target_span[1] :]
else:
    if not target_data:
        separator = b""
    elif target_data.endswith(b"\n\n"):
        separator = b""
    elif target_data.endswith(b"\n"):
        separator = b"\n"
    else:
        separator = b"\n\n"
    replacement_data = target_data + separator + section_data

if replacement_data != target_data:
    atomic_replace(arguments.target, replacement_data, target_mode)
PY
