#!/usr/bin/env python3
import hashlib
import json
import re
import shlex
import shutil
import subprocess
import sys
from html import unescape
from pathlib import Path

root = Path(__file__).resolve().parents[3]
run_dir = Path(__file__).resolve().parent
fixture_root = Path(sys.argv[1]).resolve()
project = fixture_root / 'consumer'
suite = fixture_root / 'suite'
archive = fixture_root / 'source.tar.gz'
gallery_path = 'docs/report-gallery/index.html'
manifest_path = 'docs/report-gallery/preview-manifest.json'
trial_page = 'ai-reports/2026-11-11_1411_commit-22f82fed-work-procedures-load-by-phase/index.html'
revision = 'f3840b087913d8501267dda1c1395580304f2ba4'

def run(args, cwd=root, expected=0, input_text=None):
    result = subprocess.run(args, cwd=cwd, text=True, input=input_text, capture_output=True)
    if expected is not None:
        assert result.returncode == expected, (args, result.returncode, result.stdout, result.stderr)
    return result

def gallery(*args, expected=0, cwd=project):
    launcher = cwd / ('.claude/skills' if cwd == project else 'skills') / 'do-work/tools/do-work-cli.sh'
    return run(['bash', str(launcher), '--repo-root', str(cwd), '--format', 'text',
                'ai-report-gallery', *args], cwd=cwd, expected=expected)

def starters(html):
    return [
        unescape(text) for text in
        re.findall(r'<code class="design-prompt"[^>]*>(.*?)</code>', html)
    ]

assert run(['git', 'rev-parse', 'HEAD']).stdout.strip() == revision
project.mkdir()
run(['git', 'init', '-q', '-b', 'main'], cwd=project)
run(['git', 'archive', '--format=tar.gz', '--prefix=suite/', '-o', str(archive), revision])
run(['tar', '-xf', str(archive), '-C', str(fixture_root)])
installer = run(['bash', str(suite / 'skills/do-work/tools/install-do-work-suite.sh'),
                 '--project-root', str(project), '--archive', str(archive)], input_text='y\n')
assert (project / '.claude/skills/do-work-toolbox/actions/ai-report.md').is_file(), (installer.stdout, installer.stderr)
installed_version = (project / '.claude/skills/do-work/VERSION').read_text().strip()

start = run(['just', '--dry-run', 'run-report-gallery'], expected=None)
assert start.returncode != 0 and 'does not contain recipe' in start.stderr, start.stderr
serve = gallery('--serve', expected=None, cwd=root)
assert serve.returncode != 0 and 'Usage: ai-report-gallery' in (serve.stdout + serve.stderr)
main_html = (root / gallery_path).read_text()
source_name = 'phase-by-phase'
assert f'data-gallery-card id="{source_name}"' in main_html
assert re.search(r'<figure\b', (root / trial_page).read_text())
assert not re.search(r'\b(like|dislike|feedback)\b', main_html, re.I)
assert not (root / 'do-work/report-gallery-feedback.json').exists()
(run_dir / 'REQ-745-red.log').write_text(
    f'revision={revision}\n'
    f'RED just lookup exit={start.returncode}: {start.stderr.strip()}\n'
    f'RED CLI live mode exit={serve.returncode}: {(serve.stdout + serve.stderr).strip()}\n'
    f'CONTROL worked card={source_name}; real figure exists in {trial_page}\n'
    'RED checked-in curated gallery has no feedback controls or client, and the project has no feedback record.\n'
    'LIMIT cannot perform like/dislike/reopen on this base because the recipe and controls are absent.\n'
)

sample_folder = project / 'docs/report-gallery/project-comparison'
sample_folder.mkdir(parents=True, exist_ok=True)
(sample_folder / 'sample.md').write_text(
    '---\ntitle: "Project comparison"\nuse_when: "Compare two baseline cases."\n'
    'source: "shipped"\nadded: 2026-11-18\n---\n\n## Borrow\n'
    '- Pair two independently measured outcomes.\n- Keep units explicit in each panel.\n'
)
generated_path = project / trial_page
generated_path.parent.mkdir(parents=True)
shutil.copy2(root / trial_page, generated_path)
gallery()
ordinary_html = (project / gallery_path).read_text()
ordinary_starters = starters(ordinary_html)
listed = [row.split('\t') for row in gallery('--list').stdout.splitlines()]
assert len(listed) >= 2 and len(ordinary_starters) == len(listed)
assert sum(row[2].startswith('docs/report-gallery/project-comparison/') for row in listed) == 1
assert len([row for row in listed if row[2].startswith('.claude/skills/')]) >= 1
for starter in ordinary_starters:
    command, action, flag, name = shlex.split(starter)
    assert [command, action, flag] == ['do-work-toolbox', 'ai-report', '--design']
    result = gallery('--resolve', name)
    assert len(result.stdout.strip().split('\t')) == 3
project_control = gallery('--resolve', 'Project comparison').stdout.strip()
assert project_control.startswith('project-comparison\t')
assert 'Copy starter prompt' in ordinary_html
assert 'phase-by-phase' not in '\n'.join(ordinary_starters)
assert 'report:2026-11-11_1411_commit-22f82fed-work-procedures-load-by-phase' not in '\n'.join(ordinary_starters)

manifest = json.loads((root / manifest_path).read_text())
paths = [manifest_path]
for card in manifest['cards']:
    paths.extend((card['report'], card['image']))
for relative in paths:
    source = (root / relative).resolve()
    assert source.is_relative_to(root) and source.is_file()
    destination = project / relative
    destination.parent.mkdir(parents=True, exist_ok=True)
    shutil.copy2(source, destination)
gallery()
curated_html = (project / gallery_path).read_text()
assert curated_html.count('data-gallery-card id=') == len(manifest['cards'])
assert 'data-copy="name-phase-by-phase"' in curated_html
assert len(starters(curated_html)) == 1
assert shlex.split(starters(curated_html)[0])[-1] == 'Project comparison'
assert f'data-gallery-card id="{source_name}"' in curated_html
assert not re.search(r'\b(like|dislike|feedback)\b', curated_html, re.I)

generated_names = (source_name, 'report:' + Path(trial_page).parent.name)
resolve_red = []
for name in generated_names:
    result = gallery('--resolve', name, expected=None)
    diagnostic = (result.stdout + result.stderr).strip()
    assert result.returncode != 0 and 'no gallery design matches' in diagnostic, diagnostic
    resolve_red.append(f'RED installed resolver {name!r}: exit={result.returncode}; {diagnostic}')

(run_dir / 'REQ-746-red.log').write_text(
    f'revision={revision}\ninstalled_version={installed_version}\n'
    f'CONTROL installed ordinary gallery shows {len(ordinary_starters)} starter prompts '
    f'({len(ordinary_starters) - 1} shipped and one project sample); all decode and resolve.\n'
    f'CONTROL project resolver row: {project_control}\n'
    f'CONTROL curated gallery renders {len(manifest["cards"])} worked cards; '
    'phase-by-phase has its existing name-only clipboard control.\n'
    'RED curated route displays only the project starter and suppresses shipped starter cards.\n'
    'RED installed ordinary and curated routes emit no starter for the generated trial report.\n'
    + '\n'.join(resolve_red) + '\n'
    f'CONTROL trial source sha256={hashlib.sha256((root / trial_page).read_bytes()).hexdigest()}\n'
    'LIMIT copied-text resolution was exercised through the installed CLI; writer output and live clipboard acceptance belong to the builders.\n'
)
with (run_dir / 'REQ-745-red.log').open('a') as output:
    output.write('CONTROL fresh installed curated render retains worked cards and still has no feedback controls.\n')
print((run_dir / 'REQ-745-red.log').read_text(), end='')
print((run_dir / 'REQ-746-red.log').read_text(), end='')
