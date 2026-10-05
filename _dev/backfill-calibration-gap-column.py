#!/usr/bin/env python3
"""One-off, maintainer-only backfill of do-work/calibration-log.tsv (REQ-633).

Adds the max_stamp_gap_minutes column to a log whose header still ends at
completed_at, filling each row from its archived REQ's frontmatter. The value is
the largest gap, in whole minutes rounded down, between consecutive lifecycle
stamps (claimed_at, the optional phase stamps, completed_at; sorted by time;
release_at left out) — the rule the core CLI's requestmodel.FormatCalibrationRow
writes for new rows. wall_minutes and every other column are left untouched.

Idempotent: a log whose header already ends with max_stamp_gap_minutes is left
alone. A row whose REQ cannot be found in do-work/archive/ gets an empty value
and is reported.

Usage: _dev/backfill-calibration-gap-column.py [--repo-root DIR] [--dry-run]
"""
import argparse
from datetime import datetime, timezone
import os
from pathlib import Path
import re
import sys

GAP_COLUMN_NAME = "max_stamp_gap_minutes"
# Same set as requestmodel.CalibrationGapStampFields plus completed_at.
GAP_STAMP_FIELDS = (
    "claimed_at", "planning_at", "dispatch_at", "builder_handback_at",
    "integration_at", "review_at", "remediation_at", "re_review_at", "completed_at",
)
TOP_LEVEL_SCALAR = re.compile(r"^([A-Za-z_][A-Za-z0-9_]*):[ \t]*(.*?)[ \t]*$")


def parse_timestamp(text):
    """Mirror requestmodel.ParseTimestamp: RFC3339, offset-less datetime, bare date."""
    text = text.strip().strip("'\"")
    if not text:
        return None
    if text.endswith("Z"):
        text = text[:-1] + "+00:00"
    try:
        parsed = datetime.fromisoformat(text)
    except ValueError:
        return None
    if parsed.tzinfo is None:
        parsed = parsed.replace(tzinfo=timezone.utc)
    return parsed


def read_frontmatter_scalars(request_path):
    lines = request_path.read_text(encoding="utf-8").splitlines()
    if not lines or lines[0].strip() != "---":
        return {}
    scalars = {}
    for line in lines[1:]:
        if line.strip() == "---":
            break
        match = TOP_LEVEL_SCALAR.match(line)
        if match and match.group(1) not in scalars:
            scalars[match.group(1)] = match.group(2)
    return scalars


def index_archive_by_request_id(archive_root):
    """Map each REQ id to the archived files whose frontmatter declares it."""
    files_by_request_id = {}
    for request_path in sorted(archive_root.rglob("REQ-*.md")):
        request_id = read_frontmatter_scalars(request_path).get("id", "").strip().strip("'\"")
        if request_id:
            files_by_request_id.setdefault(request_id, []).append(request_path)
    return files_by_request_id


def largest_stamp_gap_minutes(scalars):
    instants = sorted(
        instant for instant in (parse_timestamp(scalars.get(field, "")) for field in GAP_STAMP_FIELDS)
        if instant is not None
    )
    if len(instants) < 2:
        return None
    largest_gap_seconds = max((later - earlier).total_seconds() for earlier, later in zip(instants, instants[1:]))
    return int(largest_gap_seconds // 60)


def main():
    argument_parser = argparse.ArgumentParser(description=__doc__.splitlines()[0])
    argument_parser.add_argument("--repo-root", default=".", help="repository holding do-work/ (default: current directory)")
    argument_parser.add_argument("--dry-run", action="store_true", help="report what would change and write nothing")
    arguments = argument_parser.parse_args()

    repo_root = Path(arguments.repo_root)
    log_path = repo_root / "do-work" / "calibration-log.tsv"
    log_lines = log_path.read_text(encoding="utf-8").splitlines()
    if not log_lines:
        print(f"{log_path}: empty, nothing to backfill")
        return 0
    header_columns = log_lines[0].split("\t")
    if header_columns[-1] == GAP_COLUMN_NAME:
        print(f"{log_path}: header already ends with {GAP_COLUMN_NAME}; nothing to do")
        return 0
    if header_columns[0] != "req_id" or header_columns[-1] != "completed_at":
        print(f"{log_path}: unexpected header {log_lines[0]!r}; refusing to guess", file=sys.stderr)
        return 1

    files_by_request_id = index_archive_by_request_id(repo_root / "do-work" / "archive")
    output_lines = [log_lines[0] + "\t" + GAP_COLUMN_NAME]
    filled_count = 0
    over_two_hours_count = 0
    unresolved_rows = []
    for row in log_lines[1:]:
        if not row.strip():
            output_lines.append(row)
            continue
        request_id = row.split("\t")[0]
        archived_files = files_by_request_id.get(request_id, [])
        gap_minutes = None
        if len(archived_files) == 1:
            gap_minutes = largest_stamp_gap_minutes(read_frontmatter_scalars(archived_files[0]))
        if gap_minutes is None:
            unresolved_rows.append(f"{request_id} ({len(archived_files)} archive match(es))")
            output_lines.append(row + "\t")
            continue
        filled_count += 1
        if gap_minutes > 120:
            over_two_hours_count += 1
        output_lines.append(f"{row}\t{gap_minutes}")

    data_row_count = sum(1 for row in log_lines[1:] if row.strip())
    print(f"{log_path}: {data_row_count} rows, {filled_count} filled, {over_two_hours_count} with a gap over 120 min, {len(unresolved_rows)} unresolved")
    for unresolved in unresolved_rows:
        print(f"  unresolved: {unresolved}")
    if arguments.dry_run:
        print("dry run: nothing written")
        return 0
    temporary_path = log_path.with_name(log_path.name + ".backfill-tmp")
    temporary_path.write_text("\n".join(output_lines) + "\n", encoding="utf-8")
    os.replace(temporary_path, log_path)
    print(f"{log_path}: rewritten with {GAP_COLUMN_NAME}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
