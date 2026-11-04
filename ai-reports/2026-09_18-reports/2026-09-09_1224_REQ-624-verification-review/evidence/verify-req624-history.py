from pathlib import Path
from collections import Counter
from urllib.parse import urlparse
import io,json,re,subprocess,tarfile
baseline_revision='49c61ba76932abd5d6a40f7509ac3f56193dfea8'
baseline_paths=subprocess.check_output(['git','ls-tree','--name-only',baseline_revision],text=True).splitlines()
baseline_paths=['CHANGELOG.md']+sorted(p for p in baseline_paths if re.fullmatch(r'CHANGELOG-20.*\.md',p))
before=[]
for path in baseline_paths:
    body=subprocess.check_output(['git','show',baseline_revision+':'+path])
    matches=list(re.finditer(rb'^## ([0-9]+\.[0-9]+\.[0-9]+)(?=\s)',body,re.M))
    for index,match in enumerate(matches):
        end=matches[index+1].start() if index+1<len(matches) else len(body)
        before.append({'version':match.group(1).decode(),'file':path,'body':body[match.start():end].decode()})
before_map={r['version']:r['body'] for r in before}
paths=['CHANGELOG.md']+sorted(str(p) for p in Path('.').glob('CHANGELOG-20*.md'))
after=[]; counts={}; headers={}
for path in paths:
    body=Path(path).read_bytes()
    matches=list(re.finditer(rb'^## ([0-9]+\.[0-9]+\.[0-9]+)(?=\s)',body,re.M))
    headers[path]=body[:matches[0].start()].decode()
    versions=[]
    for index,match in enumerate(matches):
        version=match.group(1).decode(); versions.append(version)
        end=matches[index+1].start() if index+1<len(matches) else len(body)
        entry=body[match.start():end].decode()
        assert version in before_map, ('unexpected version',version)
        assert entry==before_map[version], ('changed release bytes',path,version)
        after.append(version)
    numeric=[tuple(map(int,v.split('.'))) for v in versions]
    assert numeric==sorted(numeric,reverse=True), ('order changed',path)
    counts[path]={'count':len(versions),'newest':versions[0],'oldest':versions[-1],'bytes':len(body)}
assert Counter(after)==Counter(r['version'] for r in before)
live=Path('CHANGELOG.md').read_bytes()
assert live==Path('skills/do-work/CHANGELOG.md').read_bytes()
assert counts['CHANGELOG.md']['count']==50
assert re.findall(rb'^## ([0-9]+\.[0-9]+\.[0-9]+)(?=\s)',live,re.M)==[r['version'].encode() for r in before[:50]]
assert set(paths).issuperset(r['file'] for r in before)
links_checked=[]
for path,header in headers.items():
    for target in re.findall(r'\]\(([^)]+)\)',header):
        parsed=urlparse(target)
        if parsed.scheme:
            assert target.startswith('https://github.com/knews2019/skill-do-work/blob/main/'),(path,target)
            destination=parsed.path.split('/blob/main/',1)[1]
        else:
            destination=str(Path(path).parent/parsed.path)
        assert Path(destination).is_file(),(path,target)
        links_checked.append([path,target])
new_archives=set(paths)-{r['file'] for r in before}
assert len(new_archives)==1
assert 'https://github.com/knews2019/skill-do-work/blob/main/'+next(iter(new_archives)) in re.findall(r'\]\(([^)]+)\)',headers['CHANGELOG.md'])
archive_bytes=subprocess.check_output(['git','archive','HEAD'])
with tarfile.open(fileobj=io.BytesIO(archive_bytes)) as archive:
    names=set(archive.getnames())
    assert 'CHANGELOG.md' in names and 'skills/do-work/CHANGELOG.md' in names
    assert not any(path in names for path in paths[1:])
    assert archive.extractfile('CHANGELOG.md').read()==archive.extractfile('skills/do-work/CHANGELOG.md').read()==live
# Execute the published version action's extraction convention: first five releases, newest last.
recent=re.findall(rb'^## ([0-9]+\.[0-9]+\.[0-9]+)(?=\s)',b'\n'.join(live.splitlines()[:80]),re.M)[:5]
assert recent==[r['version'].encode() for r in before[:5]]
assert Path('VERSION').read_text().strip()=='0.305.37'
report={'result':'PASS','total_unique_releases':len(after),'live':counts['CHANGELOG.md'],'archives':{p:c for p,c in counts.items() if p!='CHANGELOG.md'},'all_release_blocks_byte_identical':True,'mirror_byte_identical':True,'archive_export_exclusion':True,'archive_header_links_valid':len(links_checked),'recent_version_order':[v.decode() for v in reversed(recent)]}
print(json.dumps(report,indent=2))
