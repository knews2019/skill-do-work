#!/usr/bin/env python3
"""Reproduce the phase-loading audit against its pinned pre-change revision.

UTF-8 bytes, LF line counts; each procedure file counted once. This measures
procedure reading, not task artifacts, model tokens, latency, cost or quality.
"""
from pathlib import Path
import subprocess
import re

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
# The walkthrough is static: conditional procedures inside a loaded action are
# counted conservatively when they name a reference section. No extra phase is
# entered just because its pointer appears in work.md.
REFERENCE_PATH = CORE_PATH + 'actions/work-reference.md'
BUILD_CREWS = ['general', 'coding-guardrails', 'shared-principles',
               'communication-style', 'prompt-injection', 'anti-slop']
REVIEW_CREWS = [name for name in BUILD_CREWS if name != 'general']


def reference_sections(before=False):
    """Return real sections, ignoring headings in fenced templates."""
    body = read_source(REFERENCE_PATH, before)
    headings = []
    offset = 0
    fence = None
    for line in body.splitlines(keepends=True):
        text = line.decode().lstrip()
        marker = re.match(r'(`{3,}|~{3,})', text)
        if marker:
            token = marker[1]
            if fence is None:
                fence = token
            elif token[0] == fence[0] and len(token) >= len(fence):
                fence = None
        elif fence is None and re.match(rb'#{2,6} ', line):
            heading = line.decode().strip()
            headings.append((heading.split(' ', 1)[1], offset, len(heading.split(' ', 1)[0])))
        offset += len(line)
    return [(title, begin, next((other_begin for _, other_begin, other_level in headings[index + 1:]
                                       if other_level <= level), len(body)))
            for index, (title, begin, level) in enumerate(headings)]


def reading_units(paths, build, before=False):
    crews = BUILD_CREWS if build else REVIEW_CREWS
    full_paths = list(dict.fromkeys(paths + [CORE_PATH + 'crew-members/' + name + '.md' for name in crews]))
    units = [(path, read_source(path, before)) for path in full_paths]
    corpus = '\n'.join(body.decode() for _, body in units)
    sections = reference_sections(before)
    selected = set()
    # Resolve exact named heading mentions, including shortened step citations,
    # and follow nested section citations until the set stops changing.
    while True:
        found = set()
        for index, (title, _, _) in enumerate(sections):
            short_title = re.sub(r' \(.*\)$', '', title)
            if any(name in corpus for name in (title, short_title)):
                found.add(index)
        added = found - selected
        if not added:
            break
        selected.update(added)
        corpus += '\n' + '\n'.join(read_source(REFERENCE_PATH, before)[sections[index][1]:sections[index][2]].decode()
                                   for index in sorted(added))
    reference_body = read_source(REFERENCE_PATH, before)
    for index in sorted(selected):
        title, begin, end = sections[index]
        if any(other_begin <= begin and end <= other_end and other_index != index
               for other_index in selected
               for _, other_begin, other_end in [sections[other_index]]):
            continue
        units.append((REFERENCE_PATH + ' § ' + title, reference_body[begin:end]))
    # Parent sections subsume their nested sections; repeated references
    # and parent/child overlaps are charged only once.
    return units


print('| Walkthrough | Before files | After files | Before lines | After lines | Before bytes | After bytes |')
print('|---|---:|---:|---:|---:|---:|---:|')
read_sets = []
for name, (before_paths, after_paths) in walkthroughs.items():
    before_units = reading_units(before_paths, name != 'Review only', True)
    after_units = reading_units(after_paths, name != 'Review only')
    def unit_totals(units):
        return (len({path.split(' § ')[0] for path, _ in units}),
                sum(body.count(b'\n') for _, body in units), sum(len(body) for _, body in units))
    before_files, before_lines, before_bytes = unit_totals(before_units)
    after_files, after_lines, after_bytes = unit_totals(after_units)
    print(f'| {name} | {before_files} | {after_files} | {before_lines} | {after_lines} | {before_bytes} | {after_bytes} |')
    read_sets.append((name, before_units, after_units))
    for path in after_paths:
        if '-reference.md' in path:
            assert './' + Path(path).name in read_source(CORE_PATH + 'actions/work.md').decode(), path

for name, before_units, after_units in read_sets:
    for label, units in [('Before', before_units), ('After', after_units)]:
        print(f'\n{name} — {label}:')
        for path, body in units:
            print(f'  {len(body):>7} bytes  {path}')

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
