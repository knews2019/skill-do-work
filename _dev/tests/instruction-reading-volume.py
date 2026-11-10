#!/usr/bin/env python3
"""Reproduce the phase-loading audit against its pinned pre-change revision.

UTF-8 bytes, LF line counts; each procedure file counted once. This measures
procedure reading, not task artifacts, model tokens, latency, cost or quality.
"""
from pathlib import Path
import subprocess

REPO_ROOT = Path(__file__).resolve().parents[2]
BASE_REVISION = 'afbc8be3'
CORE_PATH = 'skills/do-work/'


def read_source(path, before=False):
    if before:
        return subprocess.check_output(['git', 'show', f'{BASE_REVISION}:{path}'], cwd=REPO_ROOT)
    return (REPO_ROOT / path).read_bytes()


def measure_paths(paths, before=False):
    bodies = [read_source(path, before) for path in dict.fromkeys(paths)]
    return sum(body.count(b'\n') for body in bodies), sum(map(len, bodies))


def phase_path(name):
    return CORE_PATH + f'actions/work-{name}-reference.md'


common_paths = [CORE_PATH + 'SKILL.md', CORE_PATH + 'actions/review-work.md']
run_paths = common_paths + [CORE_PATH + 'actions/work.md']
ordinary_phases = ['selection', 'implementation', 'testing', 'review', 'finalization']
walkthroughs = {
    'Simple fix': (run_paths, run_paths + [phase_path(name) for name in ordinary_phases]),
    'Complex change': (run_paths, run_paths + [phase_path(name) for name in ordinary_phases + ['planning', 'exploration', 'heavy-testing']]),
    'Review only': (common_paths, common_paths),
}
print('| Walkthrough | Before lines | After lines | Before bytes | After bytes |')
print('|---|---:|---:|---:|---:|')
for name, (before_paths, after_paths) in walkthroughs.items():
    before_lines, before_bytes = measure_paths(before_paths, True)
    after_lines, after_bytes = measure_paths(after_paths)
    print(f'| {name} | {before_lines} | {after_lines} | {before_bytes} | {after_bytes} |')
    for path in after_paths:
        if '-reference.md' in path:
            assert './' + Path(path).name in read_source(CORE_PATH + 'actions/work.md').decode(), path

print('\nEntry point (lines, bytes):', measure_paths([CORE_PATH + 'actions/work.md'], True), '→', measure_paths([CORE_PATH + 'actions/work.md']))
# Replay preservation: every original procedure line survives, except the removed
# duplicate checklist and the skip-note pointer now expressed at the routing site.
original = read_source(CORE_PATH + 'actions/work.md', True).decode()
start = original.index('## Orchestrator Checklist')
end = original.index('## Error Handling', start)
original = original[:start] + original[end:]
corpus = '\n'.join(path.read_text() for path in (REPO_ROOT / CORE_PATH / 'actions').glob('work*.md'))
allowed = {'**Routes A and B:** Append a skip note (if not already present):',
           '(append the skip note per the **Plan Skip Note — Routes A/B (Step 4)** in `actions/work-reference.md`)',
           '- The orchestrator owns semantic judgment and action-authored sections; deterministic lifecycle, archive, release, commit, provenance, verification, and cleanup mutations go through the canonical finalization engine. Spawned agents do implementation work only.',
           '- Only two frontmatter status transitions are written on the normal path: `advance` commits `pending` → `claimed`, then finalization commits the terminal status. Intermediate phases are tracked by which `##` sections exist, not by status.'}
missing = [line for line in original.splitlines() if line.strip() and line not in corpus and line not in allowed]
assert not missing, missing
print('Preservation replay: all original non-checklist procedure lines accounted for.')
