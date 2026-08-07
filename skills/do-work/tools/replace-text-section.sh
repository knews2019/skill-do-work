#!/usr/bin/env bash
# Atomically create, append, migrate, or replace do-work's managed text section.
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
LEGACY_HEADER = b"# --- do-work board recipes (installed by `do-work install just-kanban`) ---"
RECIPE_NAMES = (
    b"run-kanban",
    b"run-kanban-cli",
    b"kanban-static",
    b"kanban-summary",
    b"run-do-work-update",
)


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


def recipe_header_matches(body: bytes, recipe_name: bytes) -> bool:
    return re.match(rb"^" + re.escape(recipe_name) + rb"(?:[ \t].*)?:[ \t]*$", body) is not None


def legacy_span(data: bytes):
    lines, offsets = lines_with_offsets(data)
    bodies = [line_body(line) for line in lines]
    header_indexes = [index for index, body in enumerate(bodies) if body == LEGACY_HEADER]
    recipe_indexes = {}
    any_recipe_header = False

    for recipe_name in RECIPE_NAMES:
        matches = [
            index
            for index, body in enumerate(bodies)
            if recipe_header_matches(body, recipe_name)
        ]
        if matches:
            any_recipe_header = True
        recipe_indexes[recipe_name] = matches

    if not header_indexes and not any_recipe_header:
        return None
    if len(header_indexes) != 1:
        die("legacy migration requires exactly one legacy do-work recipe header")
    if any(len(recipe_indexes[name]) != 1 for name in RECIPE_NAMES):
        die("legacy migration requires exactly one of each of the five do-work recipes")

    ordered_indexes = [recipe_indexes[name][0] for name in RECIPE_NAMES]
    if ordered_indexes != sorted(ordered_indexes) or header_indexes[0] >= ordered_indexes[0]:
        die("legacy do-work recipes are reversed, duplicated, or out of order")

    known_headers = set(ordered_indexes)
    for index in range(header_indexes[0] + 1, ordered_indexes[-1]):
        body = bodies[index]
        if (
            not body
            or body.startswith(b"#")
            or lines[index].startswith((b" ", b"\t"))
            or index in known_headers
        ):
            continue
        die("legacy do-work block contains interleaved custom content; refusing to own it")

    last_header = ordered_indexes[-1]
    body_index = last_header + 1
    if body_index >= len(lines) or not lines[body_index].startswith((b" ", b"\t")):
        die("legacy run-do-work-update recipe has no indented body")
    last_body = body_index
    while last_body + 1 < len(lines) and lines[last_body + 1].startswith((b" ", b"\t")):
        last_body += 1

    return offsets[header_indexes[0]], offsets[last_body] + len(lines[last_body])


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
parser.add_argument("--migrate-legacy-do-work", action="store_true")
parser.add_argument("--help", action="store_true")
try:
    arguments, residue = parser.parse_known_args()
except SystemExit:
    die("usage: replace-text-section.sh --target <path> --section-file <path> [--template-file <path>] [--migrate-legacy-do-work]")
if arguments.help:
    print("usage: replace-text-section.sh --target <path> --section-file <path> [--template-file <path>] [--migrate-legacy-do-work]")
    raise SystemExit(0)
if residue or not arguments.target or not arguments.section_file:
    die("usage: replace-text-section.sh --target <path> --section-file <path> [--template-file <path>] [--migrate-legacy-do-work]")

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

if target_span is not None:
    replacement_data = target_data[: target_span[0]] + section_data + target_data[target_span[1] :]
else:
    migration_span = legacy_span(target_data) if arguments.migrate_legacy_do_work else None
    if migration_span is not None:
        replacement_data = target_data[: migration_span[0]] + section_data + target_data[migration_span[1] :]
    else:
        # Refuse to append beside recognizable but incomplete legacy recipes even when the
        # caller forgot the migration flag. Duplication is worse than an explicit hard stop.
        if legacy_span(target_data) is not None:
            die("legacy do-work recipes require --migrate-legacy-do-work")
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
